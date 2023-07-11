package shaman

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

var TalentTreeSizes = [3]int{25, 29, 26}

// Start looking to refresh 5 minute totems at 4:55.
const TotemRefreshTime5M = time.Second * 295

const (
	SpellFlagShock     = core.SpellFlagAgentReserved1
	SpellFlagElectric  = core.SpellFlagAgentReserved2
	SpellFlagTotem     = core.SpellFlagAgentReserved3
	SpellFlagFocusable = core.SpellFlagAgentReserved4
)

func NewShaman(character core.Character, talents string, totems *proto.ShamanTotems, selfBuffs SelfBuffs, thunderstormRange bool) *Shaman {
	shaman := &Shaman{
		Character:           character,
		Talents:             &proto.ShamanTalents{},
		Totems:              totems,
		SelfBuffs:           selfBuffs,
		thunderstormInRange: thunderstormRange,
	}
	shaman.waterShieldManaMetrics = shaman.NewManaMetrics(core.ActionID{SpellID: 57960})

	core.FillTalentsProto(shaman.Talents.ProtoReflect(), talents, TalentTreeSizes)
	shaman.EnableManaBar()

	if shaman.Totems.Fire == proto.FireTotem_TotemOfWrath && !shaman.Talents.TotemOfWrath {
		shaman.Totems.Fire = proto.FireTotem_FlametongueTotem
	}

	// Add Shaman stat dependencies
	shaman.AddStatDependency(stats.Strength, stats.AttackPower, 1)
	shaman.AddStatDependency(stats.Agility, stats.AttackPower, 1)
	shaman.AddStatDependency(stats.Agility, stats.MeleeCrit, core.CritRatingPerCritChance/83.3)
	shaman.AddStatDependency(stats.BonusArmor, stats.Armor, 1)
	// Set proper Melee Haste scaling
	shaman.PseudoStats.MeleeHasteRatingPerHastePercent /= 1.3

	if selfBuffs.Shield == proto.ShamanShield_WaterShield {
		shaman.AddStat(stats.MP5, 100)
	}

	// When using the tier bonus for snapshotting we do not use the bonus spell
	if totems.EnhTierTenBonus {
		totems.BonusSpellpower = 0
	}

	shaman.FireElemental = shaman.NewFireElemental(float64(totems.BonusSpellpower))
	return shaman
}

// Which buffs this shaman is using.
type SelfBuffs struct {
	Bloodlust bool
	Shield    proto.ShamanShield
	ImbueMH   proto.ShamanImbue
	ImbueOH   proto.ShamanImbue
}

// Indexes into NextTotemDrops for self buffs
const (
	AirTotem int = iota
	EarthTotem
	FireTotem
	WaterTotem
)

// Shaman represents a shaman character.
type Shaman struct {
	core.Character

	thunderstormInRange bool // flag if thunderstorm will be in range.

	ShamanisticRageManaThreshold float64 //% of mana to use sham. rage at

	Talents   *proto.ShamanTalents
	SelfBuffs SelfBuffs

	Totems *proto.ShamanTotems

	// The type of totem which should be dropped next and time to drop it, for
	// each totem type (earth, air, fire, water).
	NextTotemDropType [4]int32
	NextTotemDrops    [4]time.Duration

	LightningBolt   *core.Spell
	LightningBoltLO *core.Spell

	ChainLightning     *core.Spell
	ChainLightningHits []*core.Spell
	ChainLightningLOs  []*core.Spell

	LavaBurst   *core.Spell
	FireNova    *core.Spell
	LavaLash    *core.Spell
	Stormstrike *core.Spell

	LightningShield     *core.Spell
	LightningShieldAura *core.Aura

	Thunderstorm *core.Spell

	EarthShock *core.Spell
	FlameShock *core.Spell
	FrostShock *core.Spell

	FeralSpirit  *core.Spell
	SpiritWolves *SpiritWolves

	castFireElemental     bool
	FireElemental         *FireElemental
	FireElementalTotem    *core.Spell
	fireElementalSnapShot *core.SnapshotManager

	MagmaTotem           *core.Spell
	ManaSpringTotem      *core.Spell
	HealingStreamTotem   *core.Spell
	SearingTotem         *core.Spell
	StrengthOfEarthTotem *core.Spell
	TotemOfWrath         *core.Spell
	TremorTotem          *core.Spell
	StoneskinTotem       *core.Spell
	WindfuryTotem        *core.Spell
	WrathOfAirTotem      *core.Spell
	FlametongueTotem     *core.Spell

	MaelstromWeaponAura *core.Aura

	// Healing Spells
	tidalWaveProc          *core.Aura
	ancestralHealingAmount float64
	AncestralAwakening     *core.Spell
	LesserHealingWave      *core.Spell
	HealingWave            *core.Spell
	ChainHeal              *core.Spell
	Riptide                *core.Spell
	EarthShield            *core.Spell

	waterShieldManaMetrics *core.ResourceMetrics

	hasHeroicPresence bool
}

// Implemented by each Shaman spec.
type ShamanAgent interface {
	core.Agent

	// The Shaman controlled by this Agent.
	GetShaman() *Shaman
}

func (shaman *Shaman) GetCharacter() *core.Character {
	return &shaman.Character
}

func (shaman *Shaman) HasMajorGlyph(glyph proto.ShamanMajorGlyph) bool {
	return shaman.HasGlyph(int32(glyph))
}
func (shaman *Shaman) HasMinorGlyph(glyph proto.ShamanMinorGlyph) bool {
	return shaman.HasGlyph(int32(glyph))
}

func (shaman *Shaman) AddRaidBuffs(raidBuffs *proto.RaidBuffs) {
	switch shaman.Totems.Fire {
	case proto.FireTotem_TotemOfWrath:
		raidBuffs.TotemOfWrath = true
	case proto.FireTotem_FlametongueTotem:
		raidBuffs.FlametongueTotem = true
	}

	switch shaman.Totems.Water {
	case proto.WaterTotem_ManaSpringTotem:
		raidBuffs.ManaSpringTotem = core.MaxTristate(raidBuffs.ManaSpringTotem, proto.TristateEffect_TristateEffectRegular)
		if shaman.Talents.RestorativeTotems == 5 {
			raidBuffs.ManaSpringTotem = proto.TristateEffect_TristateEffectImproved
		}
	}

	switch shaman.Totems.Air {
	case proto.AirTotem_WrathOfAirTotem:
		raidBuffs.WrathOfAirTotem = true
	case proto.AirTotem_WindfuryTotem:
		wfVal := proto.TristateEffect_TristateEffectRegular
		if shaman.Talents.ImprovedWindfuryTotem > 0 {
			wfVal = proto.TristateEffect_TristateEffectImproved
		}
		raidBuffs.WindfuryTotem = core.MaxTristate(wfVal, raidBuffs.WindfuryTotem)
	}

	switch shaman.Totems.Earth {
	case proto.EarthTotem_StrengthOfEarthTotem:
		totem := proto.TristateEffect_TristateEffectRegular
		if shaman.Talents.EnhancingTotems == 3 {
			totem = proto.TristateEffect_TristateEffectImproved
		}
		raidBuffs.StrengthOfEarthTotem = core.MaxTristate(raidBuffs.StrengthOfEarthTotem, totem)
	case proto.EarthTotem_StoneskinTotem:
		raidBuffs.StoneskinTotem = core.MaxTristate(raidBuffs.StoneskinTotem, core.MakeTristateValue(
			true,
			shaman.Talents.GuardianTotems == 2,
		))
	}

	if shaman.Talents.UnleashedRage > 0 {
		raidBuffs.UnleashedRage = true
	}

	if shaman.Talents.ElementalOath > 0 {
		raidBuffs.ElementalOath = true
	}
}
func (shaman *Shaman) AddPartyBuffs(partyBuffs *proto.PartyBuffs) {
	if shaman.Talents.ManaTideTotem {
		partyBuffs.ManaTideTotems++
	}

	shaman.hasHeroicPresence = partyBuffs.HeroicPresence
}

func (shaman *Shaman) Initialize() {
	enableSnapshot := shaman.Totems.BonusSpellpower == 0

	shaman.registerChainLightningSpell()
	shaman.registerFeralSpirit()
	shaman.registerFireElementalTotem()
	shaman.registerFireNovaSpell()
	shaman.registerLavaBurstSpell()
	shaman.registerLavaLashSpell()
	shaman.registerLightningBoltSpell()
	shaman.registerLightningShieldSpell()
	shaman.registerMagmaTotemSpell()
	shaman.registerManaSpringTotemSpell()
	shaman.registerHealingStreamTotemSpell()
	shaman.registerSearingTotemSpell()
	shaman.registerShocks()
	shaman.registerStormstrikeSpell()
	shaman.registerStrengthOfEarthTotemSpell()
	shaman.registerThunderstormSpell()
	shaman.registerTotemOfWrathSpell()
	shaman.registerFlametongueTotemSpell()
	shaman.registerTremorTotemSpell()
	shaman.registerStoneskinTotemSpell()
	shaman.registerWindfuryTotemSpell()
	shaman.registerWrathOfAirTotemSpell()

	shaman.registerBloodlustCD()

	if shaman.Totems.UseFireElemental && enableSnapshot && !shaman.IsUsingAPL {
		shaman.fireElementalSnapShot = core.NewSnapshotManager(shaman.GetCharacter())
		shaman.setupProcTrackers()
	}

	if shaman.Talents.SpiritWeapons {
		shaman.PseudoStats.ThreatMultiplier -= 0.3
	}

	// Healing stream totem applies a HoT (aura) and so needs to be handled as a pre-pull action
	// instead of during init/reset.
	if shaman.Totems.Water == proto.WaterTotem_HealingStreamTotem {
		shaman.RegisterPrepullAction(0, func(sim *core.Simulation) {
			shaman.HealingStreamTotem.Cast(sim, &shaman.Unit)
		})
	}
}

func (shaman *Shaman) RegisterHealingSpells() {
	shaman.registerAncestralHealingSpell()
	shaman.registerLesserHealingWaveSpell()
	shaman.registerHealingWaveSpell()
	shaman.registerRiptideSpell()
	shaman.registerEarthShieldSpell()
	shaman.registerChainHealSpell()

	if shaman.Talents.TidalWaves > 0 {
		shaman.tidalWaveProc = shaman.GetOrRegisterAura(core.Aura{
			Label:    "Tidal Wave Proc",
			ActionID: core.ActionID{SpellID: 53390},
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Deactivate(sim)
			},
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				shaman.HealingWave.CastTimeMultiplier *= 0.7
				shaman.LesserHealingWave.BonusCritRating += core.CritRatingPerCritChance * 25
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				shaman.HealingWave.CastTimeMultiplier /= 0.7
				shaman.LesserHealingWave.BonusCritRating -= core.CritRatingPerCritChance * 25
			},
			MaxStacks: 2,
		})
	}
}

func (shaman *Shaman) Reset(sim *core.Simulation) {
	if shaman.Totems.UseFireElemental {
		shaman.setupFireElementalCooldowns()
		shaman.castFireElemental = false
	}

	// Check to see if we are casting a totem to set its expire time.
	for i := range shaman.NextTotemDrops {
		shaman.NextTotemDrops[i] = core.NeverExpires
		switch i {
		case AirTotem:
			if shaman.Totems.Air != proto.AirTotem_NoAirTotem {
				shaman.NextTotemDrops[i] = TotemRefreshTime5M
				shaman.NextTotemDropType[i] = int32(shaman.Totems.Air)
			}
		case EarthTotem:
			if shaman.Totems.Earth != proto.EarthTotem_NoEarthTotem {
				shaman.NextTotemDrops[i] = TotemRefreshTime5M
				shaman.NextTotemDropType[i] = int32(shaman.Totems.Earth)
			}
		case FireTotem:
			shaman.NextTotemDropType[FireTotem] = int32(shaman.Totems.Fire)
			if shaman.NextTotemDropType[FireTotem] != int32(proto.FireTotem_NoFireTotem) {
				if shaman.NextTotemDropType[FireTotem] != int32(proto.FireTotem_TotemOfWrath) &&
					shaman.NextTotemDropType[FireTotem] != int32(proto.FireTotem_FlametongueTotem) {
					if !shaman.Totems.UseFireMcd {
						shaman.NextTotemDrops[FireTotem] = 0
					}
				} else {
					shaman.NextTotemDrops[FireTotem] = TotemRefreshTime5M
					if shaman.NextTotemDropType[FireTotem] == int32(proto.FireTotem_TotemOfWrath) {
						shaman.applyToWDebuff(sim)
					}
				}
			}
		case WaterTotem:
			shaman.NextTotemDropType[i] = int32(shaman.Totems.Water)
			shaman.NextTotemDrops[i] = TotemRefreshTime5M
		}
	}

	shaman.FlameShock.CD.Reset()
}

func (shaman *Shaman) setupProcTrackers() {
	snapshotManager := shaman.fireElementalSnapShot

	snapshotManager.AddProc(40212, "Potion of Wild Magic", true)
	snapshotManager.AddProc(33697, "Blood Fury", true)
	snapshotManager.AddProc(59620, "Berserking MH Proc", false)
	snapshotManager.AddProc(59620, "Berserking OH Proc", false)

	//AP Ring Procs
	snapshotManager.AddProc(44308, "Signet of Edward the Odd Proc", false)
	snapshotManager.AddProc(50401, "Ashen Band of Unmatched Vengeance Proc", false)
	snapshotManager.AddProc(50402, "Ashen Band of Endless Vengeance Proc", false)
	snapshotManager.AddProc(52571, "Ashen Band of Unmatched Might Proc", false)
	snapshotManager.AddProc(52572, "Ashen Band of Endless Might Proc", false)

	//SP Trinket Procs
	snapshotManager.AddProc(40255, "Dying Curse Proc", false)
	snapshotManager.AddProc(40682, "Sundial of the Exiled Proc", false)
	snapshotManager.AddProc(37660, "Forge Ember Proc", false)
	snapshotManager.AddProc(45518, "Flare of the Heavens Proc", false)
	snapshotManager.AddProc(54572, "Charred Twilight Scale Proc", false)
	snapshotManager.AddProc(54588, "Charred Twilight Scale H Proc", false)
	snapshotManager.AddProc(47213, "Abyssal Rune Proc", false)
	snapshotManager.AddProc(45490, "Pandora's Plea Proc", false)
	snapshotManager.AddProc(50348, "Dislodged Foreign Object H", false)
	snapshotManager.AddProc(50353, "Dislodged Foreign Object", false)
	snapshotManager.AddProc(50360, "Phylactery of the Nameless Lich Proc", false)
	snapshotManager.AddProc(50365, "Phylactery of the Nameless Lich H Proc", false)
	snapshotManager.AddProc(50345, "Muradin's Spyglass H Proc", false)
	snapshotManager.AddProc(50340, "Muradin's Spyglass Proc", false)

	// SP Ring Procs
	snapshotManager.AddProc(50398, "Ashen Band of Endless Destruction", false)

	//AP Trinket Procs
	snapshotManager.AddProc(40684, "Mirror of Truth Proc", false)
	snapshotManager.AddProc(45522, "Blood of the Old God Proc", false)
	snapshotManager.AddProc(40767, "Sonic Booster Proc", false)
	snapshotManager.AddProc(44914, "Anvil of Titans Proc", false)
	snapshotManager.AddProc(45286, "Pyrite Infuser Proc", false)
	snapshotManager.AddProc(47214, "Banner of Victory Proc", false)
	snapshotManager.AddProc(49074, "Coren's Chromium Coaster Proc", false)
	snapshotManager.AddProc(50342, "Whispering Fanged Skull Proc", false)
	snapshotManager.AddProc(50343, "Whispering Fanged Skull H Proc", false)
	snapshotManager.AddProc(54569, "Sharpened Twilight Scale Proc", false)
	snapshotManager.AddProc(54590, "Sharpened Twilight Scale H Proc", false)
	snapshotManager.AddProc(47115, "Deaths Verdict Agility Proc", false)
	snapshotManager.AddProc(47131, "Deaths Verdict H Agility Proc", false)
	snapshotManager.AddProc(47303, "Deaths Choice Agility Proc", false)
	snapshotManager.AddProc(47464, "Deaths Choice H Agility Proc", false)
	snapshotManager.AddProc(71492, "Deathbringer's Will Strength Proc", false)
	snapshotManager.AddProc(71561, "Deathbringer's Will H Strength Proc", false)
	snapshotManager.AddProc(71492, "Deathbringer's Will Agility Proc", false)
	snapshotManager.AddProc(71561, "Deathbringer's Will H Agility Proc", false)
	snapshotManager.AddProc(71492, "Deathbringer's Will AP Proc", false)
	snapshotManager.AddProc(71561, "Deathbringer's Will H AP Proc", false)

	// Tier Bonus
	snapshotManager.AddProc(70831, "Maelstrom Power", false)
}

func (shaman *Shaman) setupFireElementalCooldowns() {
	if shaman.fireElementalSnapShot == nil {
		return
	}

	shaman.fireElementalSnapShot.ClearMajorCooldowns()

	// blood fury (orc)
	shaman.fireElementalCooldownSync(core.ActionID{SpellID: 33697}, false)

	// potion of Wild Magic
	shaman.fireElementalCooldownSync(core.ActionID{ItemID: 40212}, true)

	//active sp trinkets
	shaman.fireElementalCooldownSync(core.ActionID{ItemID: 37873}, false)
	shaman.fireElementalCooldownSync(core.ActionID{ItemID: 45148}, false)
	shaman.fireElementalCooldownSync(core.ActionID{ItemID: 48724}, false)
	shaman.fireElementalCooldownSync(core.ActionID{ItemID: 50357}, false)

	// active ap trinkets
	shaman.fireElementalCooldownSync(core.ActionID{ItemID: 35937}, false)
	shaman.fireElementalCooldownSync(core.ActionID{ItemID: 36871}, false)
	shaman.fireElementalCooldownSync(core.ActionID{ItemID: 37166}, false)
	shaman.fireElementalCooldownSync(core.ActionID{ItemID: 37556}, false)
	shaman.fireElementalCooldownSync(core.ActionID{ItemID: 37557}, false)
	shaman.fireElementalCooldownSync(core.ActionID{ItemID: 38080}, false)
	shaman.fireElementalCooldownSync(core.ActionID{ItemID: 38081}, false)
	shaman.fireElementalCooldownSync(core.ActionID{ItemID: 38761}, false)
	shaman.fireElementalCooldownSync(core.ActionID{ItemID: 39257}, false)
	shaman.fireElementalCooldownSync(core.ActionID{ItemID: 45263}, false)
	shaman.fireElementalCooldownSync(core.ActionID{ItemID: 46086}, false)
	shaman.fireElementalCooldownSync(core.ActionID{ItemID: 47734}, false)
}

func (shaman *Shaman) fireElementalCooldownSync(actionID core.ActionID, isPotion bool) {
	if majorCd := shaman.Character.GetMajorCooldown(actionID); majorCd != nil {
		majorCd.ShouldActivate = func(sim *core.Simulation, character *core.Character) bool {
			return shaman.castFireElemental || (shaman.FireElementalTotem.CD.TimeToReady(sim) > majorCd.Spell.CD.Duration && !isPotion) || shaman.FireElementalTotem.CD.ReadyAt() > shaman.Env.Encounter.Duration
		}

		shaman.fireElementalSnapShot.AddMajorCooldown(majorCd)
	}
}

func (shaman *Shaman) ElementalCritMultiplier(secondary float64) float64 {
	critBonus := 0.2*float64(shaman.Talents.ElementalFury) + secondary
	return shaman.SpellCritMultiplier(1, critBonus)
}

func init() {
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceDraenei, Class: proto.Class_ClassShaman}] = stats.Stats{
		stats.Health:      6759,
		stats.Strength:    121,
		stats.Agility:     71,
		stats.Stamina:     135,
		stats.Intellect:   126,
		stats.Spirit:      145,
		stats.Mana:        4396,
		stats.SpellCrit:   2.2 * core.CritRatingPerCritChance,
		stats.AttackPower: 95, // TODO: confirm this.
		stats.MeleeCrit:   2.92 * core.CritRatingPerCritChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceOrc, Class: proto.Class_ClassShaman}] = stats.Stats{
		stats.Health:      6759,
		stats.Strength:    123,
		stats.Agility:     71,
		stats.Stamina:     138,
		stats.Intellect:   122,
		stats.Spirit:      146,
		stats.Mana:        4396,
		stats.SpellCrit:   2.2 * core.CritRatingPerCritChance,
		stats.AttackPower: 95, // TODO: confirm this.
		stats.MeleeCrit:   2.92 * core.CritRatingPerCritChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceTauren, Class: proto.Class_ClassShaman}] = stats.Stats{
		stats.Health:      6759,
		stats.Strength:    125,
		stats.Agility:     69,
		stats.Stamina:     138,
		stats.Intellect:   120,
		stats.Spirit:      145,
		stats.Mana:        4396,
		stats.SpellCrit:   2.2 * core.CritRatingPerCritChance,
		stats.AttackPower: 95, // TODO: confirm this.
		stats.MeleeCrit:   2.92 * core.CritRatingPerCritChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceTroll, Class: proto.Class_ClassShaman}] = stats.Stats{
		stats.Health:      6759,
		stats.Strength:    121,
		stats.Agility:     76,
		stats.Stamina:     137,
		stats.Intellect:   122,
		stats.Spirit:      144,
		stats.Mana:        4396,
		stats.SpellCrit:   2.2 * core.CritRatingPerCritChance,
		stats.AttackPower: 95, // TODO: confirm this.
		stats.MeleeCrit:   2.92 * core.CritRatingPerCritChance,
	}
}
