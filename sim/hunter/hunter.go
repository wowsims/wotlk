package hunter

import (
	"time"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/proto"
	"github.com/wowsims/tbc/sim/core/stats"
)

const ThoridalTheStarsFuryItemID = 34334

func RegisterHunter() {
	core.RegisterAgentFactory(
		proto.Player_Hunter{},
		proto.Spec_SpecHunter,
		func(character core.Character, options proto.Player) core.Agent {
			return NewHunter(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_Hunter)
			if !ok {
				panic("Invalid spec value for Hunter!")
			}
			player.Spec = playerSpec
		},
	)
}

type Hunter struct {
	core.Character

	Talents  proto.HunterTalents
	Options  proto.Hunter_Options
	Rotation proto.Hunter_Rotation

	pet *HunterPet

	AmmoDPS         float64
	AmmoDamageBonus float64

	hasGronnstalker2Pc bool
	currentAspect      *core.Aura

	killCommandEnabledUntil time.Duration // Time that KC enablement expires.
	killCommandBlocked      bool          // True while Steady Shot is casting, to prevent KC.

	latency     time.Duration
	timeToWeave time.Duration

	nextAction   int
	nextActionAt time.Duration

	// Expected single-cast damage values calculated by the presim, used for adaptive logic.
	avgShootDmg  float64
	avgWeaveDmg  float64
	avgSteadyDmg float64
	avgMultiDmg  float64
	avgArcaneDmg float64

	// Used for deciding when we can use hawk for the rest of the fight.
	manaSpentPerSecondAtFirstAspectSwap float64
	permaHawk                           bool

	// Cached values for adaptive rotation calcs.
	rangedSwingSpeed   float64
	rangedWindup       float64
	shootDPS           float64
	weaveDPS           float64
	steadyDPS          float64
	steadyShotCastTime float64
	multiShotCastTime  float64
	arcaneShotCastTime float64
	useMultiForCatchup bool

	AspectOfTheHawk  *core.Spell
	AspectOfTheViper *core.Spell

	AimedShot    *core.Spell
	ArcaneShot   *core.Spell
	KillCommand  *core.Spell
	MultiShot    *core.Spell
	RapidFire    *core.Spell
	RaptorStrike *core.Spell
	ScorpidSting *core.Spell
	SerpentSting *core.Spell
	SteadyShot   *core.Spell

	SerpentStingDot *core.Dot

	AspectOfTheHawkAura  *core.Aura
	AspectOfTheViperAura *core.Aura
	ScorpidStingAura     *core.Aura
	TalonOfAlarAura      *core.Aura

	hardcastOnComplete core.CastFunc
}

func (hunter *Hunter) GetCharacter() *core.Character {
	return &hunter.Character
}

func (hunter *Hunter) GetHunter() *Hunter {
	return hunter
}

func (hunter *Hunter) AddRaidBuffs(raidBuffs *proto.RaidBuffs) {
}
func (hunter *Hunter) AddPartyBuffs(partyBuffs *proto.PartyBuffs) {
	if hunter.Talents.TrueshotAura {
		partyBuffs.TrueshotAura = true
	}
}

func (hunter *Hunter) Initialize() {
	// Update auto crit multipliers now that we have the targets.
	hunter.AutoAttacks.MHEffect.OutcomeApplier = hunter.OutcomeFuncMeleeWhite(hunter.critMultiplier(false, hunter.CurrentTarget))
	hunter.AutoAttacks.OHEffect.OutcomeApplier = hunter.OutcomeFuncMeleeWhite(hunter.critMultiplier(false, hunter.CurrentTarget))
	hunter.AutoAttacks.RangedEffect.OutcomeApplier = hunter.OutcomeFuncRangedHitAndCrit(hunter.critMultiplier(true, hunter.CurrentTarget))

	hunter.registerAspectOfTheHawkSpell()
	hunter.registerAspectOfTheViperSpell()

	hunter.registerAimedShotSpell()
	hunter.registerArcaneShotSpell()
	hunter.registerKillCommandSpell()
	hunter.registerMultiShotSpell()
	hunter.registerRaptorStrikeSpell()
	hunter.registerScorpidStingSpell()
	hunter.registerSerpentStingSpell()
	hunter.registerSteadyShotSpell()

	hunter.hardcastOnComplete = func(sim *core.Simulation, _ *core.Unit) {
		hunter.rotation(sim, false)
	}

	hunter.DelayDPSCooldownsForArmorDebuffs()
}

func (hunter *Hunter) Reset(sim *core.Simulation) {
	hunter.killCommandEnabledUntil = 0
	hunter.killCommandBlocked = false
	hunter.nextAction = OptionNone
	hunter.nextActionAt = 0
	hunter.rangedSwingSpeed = 0
	hunter.manaSpentPerSecondAtFirstAspectSwap = 0
	hunter.permaHawk = false

	huntersMarkAura := core.HuntersMarkAura(hunter.CurrentTarget, hunter.Talents.ImprovedHuntersMark, false)
	huntersMarkAura.Activate(sim)

	if sim.Log != nil && !hunter.Rotation.LazyRotation {
		hunter.Log(sim, "Average damage values for adaptive rotation: shoot=%0.02f, weave=%0.02f, steady=%0.02f, multi=%0.02f, arcane=%0.02f", hunter.avgShootDmg, hunter.avgWeaveDmg, hunter.avgSteadyDmg, hunter.avgMultiDmg, hunter.avgArcaneDmg)
	}
}

func NewHunter(character core.Character, options proto.Player) *Hunter {
	hunterOptions := options.GetHunter()

	hunter := &Hunter{
		Character: character,
		Talents:   *hunterOptions.Talents,
		Options:   *hunterOptions.Options,
		Rotation:  *hunterOptions.Rotation,

		latency:     time.Millisecond * time.Duration(hunterOptions.Options.LatencyMs),
		timeToWeave: time.Millisecond * time.Duration(hunterOptions.Rotation.TimeToWeaveMs+hunterOptions.Options.LatencyMs),

		hasGronnstalker2Pc: ItemSetGronnstalker.CharacterHasSetBonus(&character, 2),
	}
	hunter.EnableManaBar()

	if hunter.Rotation.PercentWeaved <= 0 {
		hunter.Rotation.Weave = proto.Hunter_Rotation_WeaveNone
	}
	if hunter.Rotation.Weave == proto.Hunter_Rotation_WeaveNone {
		// Forces override of WF. When not weaving we'll be standing far back so weapon
		// stone can be used.
		hunter.HasMHWeaponImbue = true
	}

	hunter.PseudoStats.CanParry = true

	rangedWeapon := hunter.WeaponFromRanged(0)
	hunter.PseudoStats.RangedSpeedMultiplier = 1
	if hunter.HasRangedWeapon() && hunter.GetRangedWeapon().ID == ThoridalTheStarsFuryItemID {
		hunter.PseudoStats.RangedSpeedMultiplier *= 1.15
	} else {
		switch hunter.Options.Ammo {
		case proto.Hunter_Options_TimelessArrow:
			hunter.AmmoDPS = 53
		case proto.Hunter_Options_MysteriousArrow:
			hunter.AmmoDPS = 46.5
		case proto.Hunter_Options_AdamantiteStinger:
			hunter.AmmoDPS = 43
		case proto.Hunter_Options_WardensArrow:
			hunter.AmmoDPS = 37
		case proto.Hunter_Options_HalaaniRazorshaft:
			hunter.AmmoDPS = 34
		case proto.Hunter_Options_BlackflightArrow:
			hunter.AmmoDPS = 32
		}
		hunter.AmmoDamageBonus = hunter.AmmoDPS * rangedWeapon.SwingSpeed

		switch hunter.Options.QuiverBonus {
		case proto.Hunter_Options_Speed10:
			hunter.PseudoStats.RangedSpeedMultiplier *= 1.1
		case proto.Hunter_Options_Speed11:
			hunter.PseudoStats.RangedSpeedMultiplier *= 1.11
		case proto.Hunter_Options_Speed12:
			hunter.PseudoStats.RangedSpeedMultiplier *= 1.12
		case proto.Hunter_Options_Speed13:
			hunter.PseudoStats.RangedSpeedMultiplier *= 1.13
		case proto.Hunter_Options_Speed14:
			hunter.PseudoStats.RangedSpeedMultiplier *= 1.14
		case proto.Hunter_Options_Speed15:
			hunter.PseudoStats.RangedSpeedMultiplier *= 1.15
		}
	}

	hunter.EnableAutoAttacks(hunter, core.AutoAttackOptions{
		// We don't know crit multiplier until later when we see the target so just
		// use 0 for now.
		MainHand: hunter.WeaponFromMainHand(0),
		OffHand:  hunter.WeaponFromOffHand(0),
		Ranged:   rangedWeapon,
		ReplaceMHSwing: func(sim *core.Simulation, _ *core.Spell) *core.Spell {
			return hunter.TryRaptorStrike(sim)
		},
	})
	hunter.AutoAttacks.RangedEffect.BaseDamage.Calculator = core.BaseDamageFuncRangedWeapon(hunter.AmmoDamageBonus)

	if hunter.Options.RemoveRandomness {
		weaponAvg := (hunter.AutoAttacks.Ranged.BaseDamageMin + hunter.AutoAttacks.Ranged.BaseDamageMax) / 2
		hunter.AutoAttacks.Ranged.BaseDamageMin = weaponAvg
		hunter.AutoAttacks.Ranged.BaseDamageMax = weaponAvg

		hunter.AddStat(stats.MeleeHit, core.MeleeHitRatingPerHitChance*100)
		hunter.AddStat(stats.MeleeCrit, core.MeleeCritRatingPerCritChance*-100)
	}

	hunter.pet = hunter.NewHunterPet()

	hunter.AddStatDependency(stats.StatDependency{
		SourceStat:   stats.Intellect,
		ModifiedStat: stats.SpellCrit,
		Modifier: func(intellect float64, spellCrit float64) float64 {
			return spellCrit + (intellect/55)*core.SpellCritRatingPerCritChance
		},
	})

	hunter.AddStatDependency(stats.StatDependency{
		SourceStat:   stats.Strength,
		ModifiedStat: stats.AttackPower,
		Modifier: func(strength float64, attackPower float64) float64 {
			return attackPower + strength*1
		},
	})

	hunter.AddStatDependency(stats.StatDependency{
		SourceStat:   stats.Agility,
		ModifiedStat: stats.RangedAttackPower,
		Modifier: func(agility float64, rap float64) float64 {
			return rap + agility*1
		},
	})

	hunter.AddStatDependency(stats.StatDependency{
		SourceStat:   stats.Agility,
		ModifiedStat: stats.MeleeCrit,
		Modifier: func(agility float64, meleeCrit float64) float64 {
			return meleeCrit + (agility/40)*core.MeleeCritRatingPerCritChance
		},
	})

	return hunter
}

func init() {
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceBloodElf, Class: proto.Class_ClassHunter}] = stats.Stats{
		stats.Health:    3388,
		stats.Strength:  61,
		stats.Agility:   153,
		stats.Stamina:   106,
		stats.Intellect: 81,
		stats.Spirit:    82,
		stats.Mana:      3383,

		stats.AttackPower:       120,
		stats.RangedAttackPower: 130,
		stats.MeleeCrit:         -1.53 * core.MeleeCritRatingPerCritChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceDraenei, Class: proto.Class_ClassHunter}] = stats.Stats{
		stats.Health:    3388,
		stats.Strength:  65,
		stats.Agility:   148,
		stats.Stamina:   107,
		stats.Intellect: 78,
		stats.Spirit:    85,
		stats.Mana:      3383,

		stats.AttackPower:       120,
		stats.RangedAttackPower: 130,
		stats.MeleeCrit:         -1.53 * core.MeleeCritRatingPerCritChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceDwarf, Class: proto.Class_ClassHunter}] = stats.Stats{
		stats.Health:    3388,
		stats.Strength:  66,
		stats.Agility:   147,
		stats.Stamina:   111,
		stats.Intellect: 76,
		stats.Spirit:    82,
		stats.Mana:      3383,

		stats.AttackPower:       120,
		stats.RangedAttackPower: 130,
		stats.MeleeCrit:         -1.53 * core.MeleeCritRatingPerCritChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceNightElf, Class: proto.Class_ClassHunter}] = stats.Stats{
		stats.Health:    3388,
		stats.Strength:  61,
		stats.Agility:   156,
		stats.Stamina:   107,
		stats.Intellect: 77,
		stats.Spirit:    83,
		stats.Mana:      3383,

		stats.AttackPower:       120,
		stats.RangedAttackPower: 130,
		stats.MeleeCrit:         -1.53 * core.MeleeCritRatingPerCritChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceOrc, Class: proto.Class_ClassHunter}] = stats.Stats{
		stats.Health:    3388,
		stats.Strength:  67,
		stats.Agility:   148,
		stats.Stamina:   110,
		stats.Intellect: 74,
		stats.Spirit:    86,
		stats.Mana:      3383,

		stats.AttackPower:       120,
		stats.RangedAttackPower: 130,
		stats.MeleeCrit:         -1.53 * core.MeleeCritRatingPerCritChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceTauren, Class: proto.Class_ClassHunter}] = stats.Stats{
		stats.Health:    3388,
		stats.Strength:  69,
		stats.Agility:   146,
		stats.Stamina:   110,
		stats.Intellect: 72,
		stats.Spirit:    85,
		stats.Mana:      3383,

		stats.AttackPower:       120,
		stats.RangedAttackPower: 130,
		stats.MeleeCrit:         -1.53 * core.MeleeCritRatingPerCritChance,
	}
	trollStats := stats.Stats{
		stats.Health:    3388,
		stats.Strength:  65,
		stats.Agility:   153,
		stats.Stamina:   109,
		stats.Intellect: 73,
		stats.Spirit:    84,
		stats.Mana:      3383,

		stats.AttackPower:       120,
		stats.RangedAttackPower: 130,
		stats.MeleeCrit:         -1.53 * core.MeleeCritRatingPerCritChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceTroll10, Class: proto.Class_ClassHunter}] = trollStats
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceTroll30, Class: proto.Class_ClassHunter}] = trollStats
}

// Agent is a generic way to access underlying hunter on any of the agents.
type HunterAgent interface {
	GetHunter() *Hunter
}
