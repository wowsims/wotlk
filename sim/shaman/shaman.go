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

func NewShaman(character *core.Character, talents string, totems *proto.ShamanTotems, selfBuffs SelfBuffs, thunderstormRange bool) *Shaman {
	shaman := &Shaman{
		Character:           *character,
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
	shaman.AddStatDependency(stats.Agility, stats.MeleeCrit, core.CritPerAgiMaxLevel[character.Class]*core.CritRatingPerCritChance)
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
	Shield  proto.ShamanShield
	ImbueMH proto.ShamanImbue
	ImbueOH proto.ShamanImbue
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

	Talents   *proto.ShamanTalents
	SelfBuffs SelfBuffs

	Totems *proto.ShamanTotems

	// The expiration time of each totem (earth, air, fire, water).
	TotemExpirations [4]time.Duration

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

	FireElemental      *FireElemental
	FireElementalTotem *core.Spell

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
		raidBuffs.ManaSpringTotem = max(raidBuffs.ManaSpringTotem, proto.TristateEffect_TristateEffectRegular)
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
		raidBuffs.WindfuryTotem = max(wfVal, raidBuffs.WindfuryTotem)
	}

	switch shaman.Totems.Earth {
	case proto.EarthTotem_StrengthOfEarthTotem:
		totem := proto.TristateEffect_TristateEffectRegular
		if shaman.Talents.EnhancingTotems == 3 {
			totem = proto.TristateEffect_TristateEffectImproved
		}
		raidBuffs.StrengthOfEarthTotem = max(raidBuffs.StrengthOfEarthTotem, totem)
	case proto.EarthTotem_StoneskinTotem:
		raidBuffs.StoneskinTotem = max(raidBuffs.StoneskinTotem, core.MakeTristateValue(
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

	// This registration must come after all the totems are registered
	shaman.registerCallOfTheElements()

	shaman.registerBloodlustCD()

	shaman.NewTemporaryStatsAura("DC Pre-Pull SP Proc", core.ActionID{SpellID: 60494}, stats.Stats{stats.SpellPower: 765}, time.Second*10)
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
	if shaman.Totems.Fire == proto.FireTotem_TotemOfWrath {
		shaman.applyToWDebuff(sim)
	}
}

func (shaman *Shaman) ElementalCritMultiplier(secondary float64) float64 {
	critBonus := 0.2*float64(shaman.Talents.ElementalFury) + secondary
	return shaman.SpellCritMultiplier(1, critBonus)
}
