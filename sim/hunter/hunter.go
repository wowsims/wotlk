package hunter

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
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

	AspectOfTheDragonhawk *core.Spell
	AspectOfTheViper      *core.Spell

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

	AspectOfTheDragonhawkAura *core.Aura
	AspectOfTheViperAura      *core.Aura
	ImprovedSteadyShotAura    *core.Aura
	ScorpidStingAura          *core.Aura
	TalonOfAlarAura           *core.Aura

	hardcastOnComplete core.CastFunc
}

func (hunter *Hunter) GetCharacter() *core.Character {
	return &hunter.Character
}

func (hunter *Hunter) GetHunter() *Hunter {
	return hunter
}

func (hunter *Hunter) AddRaidBuffs(raidBuffs *proto.RaidBuffs) {
	if hunter.Talents.FerociousInspiration == 3 && hunter.pet != nil {
		raidBuffs.FerociousInspiration = true
	}
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

	hunter.registerAspectOfTheDragonhawkSpell()
	hunter.registerAspectOfTheViperSpell()

	hunter.registerAimedShotSpell()
	hunter.registerArcaneShotSpell()
	hunter.registerMultiShotSpell()
	hunter.registerRaptorStrikeSpell()
	hunter.registerScorpidStingSpell()
	hunter.registerSerpentStingSpell()
	hunter.registerSteadyShotSpell()

	hunter.registerKillCommandCD()
	hunter.registerRapidFireCD()

	hunter.hardcastOnComplete = func(sim *core.Simulation, _ *core.Unit) {
		hunter.rotation(sim, false)
	}

	hunter.DelayDPSCooldownsForArmorDebuffs()
}

func (hunter *Hunter) Reset(sim *core.Simulation) {
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

	// Passive bonus (used to be from quiver).
	hunter.PseudoStats.RangedSpeedMultiplier *= 1.15

	if hunter.HasRangedWeapon() && hunter.GetRangedWeapon().ID != ThoridalTheStarsFuryItemID {
		switch hunter.Options.Ammo {
		case proto.Hunter_Options_IcebladeArrow:
			hunter.AmmoDPS = 91.5
		case proto.Hunter_Options_SaroniteRazorheads:
			hunter.AmmoDPS = 67.5
		case proto.Hunter_Options_TerrorshaftArrow:
			hunter.AmmoDPS = 46.5
		case proto.Hunter_Options_TimelessArrow:
			hunter.AmmoDPS = 53
		case proto.Hunter_Options_MysteriousArrow:
			hunter.AmmoDPS = 46.5
		case proto.Hunter_Options_AdamantiteStinger:
			hunter.AmmoDPS = 43
		case proto.Hunter_Options_BlackflightArrow:
			hunter.AmmoDPS = 32
		}
		hunter.AmmoDamageBonus = hunter.AmmoDPS * rangedWeapon.SwingSpeed
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
		hunter.AddStat(stats.MeleeCrit, core.CritRatingPerCritChance*-100)
	}

	hunter.pet = hunter.NewHunterPet()

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
			return meleeCrit + (agility/40)*core.CritRatingPerCritChance
		},
	})

	return hunter
}

func init() {
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceBloodElf, Class: proto.Class_ClassHunter}] = stats.Stats{
		stats.Health:    7324,
		stats.Strength:  71,
		stats.Agility:   183,
		stats.Stamina:   126,
		stats.Intellect: 94,
		stats.Spirit:    96,
		stats.Mana:      5046,

		stats.AttackPower:       120,
		stats.RangedAttackPower: 130,
		stats.MeleeCrit:         -1.53 * core.CritRatingPerCritChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceDraenei, Class: proto.Class_ClassHunter}] = stats.Stats{
		stats.Health:    7324,
		stats.Strength:  75,
		stats.Agility:   178,
		stats.Stamina:   127,
		stats.Intellect: 91,
		stats.Spirit:    99,
		stats.Mana:      5046,

		stats.AttackPower:       120,
		stats.RangedAttackPower: 130,
		stats.MeleeCrit:         -1.53 * core.CritRatingPerCritChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceDwarf, Class: proto.Class_ClassHunter}] = stats.Stats{
		stats.Health:    7324,
		stats.Strength:  76,
		stats.Agility:   177,
		stats.Stamina:   131,
		stats.Intellect: 89,
		stats.Spirit:    96,
		stats.Mana:      5046,

		stats.AttackPower:       120,
		stats.RangedAttackPower: 130,
		stats.MeleeCrit:         -1.53 * core.CritRatingPerCritChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceNightElf, Class: proto.Class_ClassHunter}] = stats.Stats{
		stats.Health:    7324,
		stats.Strength:  71,
		stats.Agility:   193,
		stats.Stamina:   127,
		stats.Intellect: 93,
		stats.Spirit:    97,
		stats.Mana:      5046,

		stats.AttackPower:       120,
		stats.RangedAttackPower: 130,
		stats.MeleeCrit:         -1.53 * core.CritRatingPerCritChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceOrc, Class: proto.Class_ClassHunter}] = stats.Stats{
		stats.Health:    7324,
		stats.Strength:  77,
		stats.Agility:   178,
		stats.Stamina:   130,
		stats.Intellect: 87,
		stats.Spirit:    100,
		stats.Mana:      5046,

		stats.AttackPower:       120,
		stats.RangedAttackPower: 130,
		stats.MeleeCrit:         -1.53 * core.CritRatingPerCritChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceTauren, Class: proto.Class_ClassHunter}] = stats.Stats{
		stats.Health:    7324,
		stats.Strength:  79,
		stats.Agility:   183,
		stats.Stamina:   130,
		stats.Intellect: 88,
		stats.Spirit:    99,
		stats.Mana:      5046,

		stats.AttackPower:       120,
		stats.RangedAttackPower: 130,
		stats.MeleeCrit:         -1.53 * core.CritRatingPerCritChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceTroll, Class: proto.Class_ClassHunter}] = stats.Stats{
		stats.Health:    7324,
		stats.Strength:  75,
		stats.Agility:   190,
		stats.Stamina:   129,
		stats.Intellect: 89,
		stats.Spirit:    98,
		stats.Mana:      5046,

		stats.AttackPower:       120,
		stats.RangedAttackPower: 130,
		stats.MeleeCrit:         -1.53 * core.CritRatingPerCritChance,
	}
}

// Agent is a generic way to access underlying hunter on any of the agents.
type HunterAgent interface {
	GetHunter() *Hunter
}
