package shaman

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

const baseMana = 4396.0

// Start looking to refresh 5 minute totems at 4:55.
const TotemRefreshTime5M = time.Second * 295

const (
	SpellFlagShock    = core.SpellFlagAgentReserved1
	SpellFlagElectric = core.SpellFlagAgentReserved2
	SpellFlagTotem    = core.SpellFlagAgentReserved3
	SpellFlagFireNova = core.SpellFlagAgentReserved4
)

func NewShaman(character core.Character, talents proto.ShamanTalents, totems proto.ShamanTotems, selfBuffs SelfBuffs, thunderstormRange bool) *Shaman {
	if totems.Fire == proto.FireTotem_TotemOfWrath && !talents.TotemOfWrath {
		totems.Fire = proto.FireTotem_NoFireTotem
	}

	shaman := &Shaman{
		Character:           character,
		Talents:             talents,
		Totems:              totems,
		SelfBuffs:           selfBuffs,
		thunderstormInRange: thunderstormRange,
	}
	shaman.EnableManaBar()

	// Add Shaman stat dependencies
	shaman.AddStatDependency(stats.Strength, stats.AttackPower, 1.0+1)
	shaman.AddStatDependency(stats.Agility, stats.AttackPower, 1.0+1)
	shaman.AddStatDependency(stats.Agility, stats.MeleeCrit, 1.0+core.CritRatingPerCritChance/83.3)
	// Set proper Melee Haste scaling
	shaman.PseudoStats.MeleeHasteRatingPerHastePercent /= 1.3

	if selfBuffs.Shield == proto.ShamanShield_WaterShield {
		shaman.AddStat(stats.MP5, 100)
	}

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

	LavaburstWeave bool // flag to enable lava burst weaving for enh

	Talents   proto.ShamanTalents
	SelfBuffs SelfBuffs

	Totems proto.ShamanTotems

	// The type of totem which should be dropped next and time to drop it, for
	// each totem type (earth, air, fire, water).
	NextTotemDropType [4]int32
	NextTotemDrops    [4]time.Duration

	LightningBolt   *core.Spell
	LightningBoltLO *core.Spell

	ChainLightning    *core.Spell
	ChainLightningLOs []*core.Spell

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

	GraceOfAirTotem      *core.Spell
	MagmaTotem           *core.Spell
	ManaSpringTotem      *core.Spell
	SearingTotem         *core.Spell
	StrengthOfEarthTotem *core.Spell
	TotemOfWrath         *core.Spell
	TranquilAirTotem     *core.Spell
	TremorTotem          *core.Spell
	WindfuryTotem        *core.Spell
	WrathOfAirTotem      *core.Spell

	FlameShockDot   *core.Dot
	SearingTotemDot *core.Dot
	MagmaTotemDot   *core.Dot

	ClearcastingAura         *core.Aura
	ElementalMasteryAura     *core.Aura
	ElementalMasteryBuffAura *core.Aura
	NaturesSwiftnessAura     *core.Aura
	MaelstromWeaponAura      *core.Aura
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
	if shaman.Totems.Fire == proto.FireTotem_TotemOfWrath {
		raidBuffs.TotemOfWrath = true
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
}

func (shaman *Shaman) Initialize() {
	// Precompute all the spell templates.
	shaman.registerStormstrikeSpell()
	shaman.LightningBolt = shaman.newLightningBoltSpell(false)
	shaman.LightningBoltLO = shaman.newLightningBoltSpell(true)
	shaman.LavaBurst = shaman.newLavaBurstSpell()
	shaman.FireNova = shaman.newFireNovaSpell()

	if shaman.SelfBuffs.Shield == proto.ShamanShield_LightningShield {
		shaman.registerLightningShieldSpell()
	}

	shaman.ChainLightning = shaman.newChainLightningSpell(false)
	numHits := core.MinInt32(3, shaman.Env.GetNumTargets())
	shaman.ChainLightningLOs = []*core.Spell{}
	for i := int32(0); i < numHits; i++ {
		shaman.ChainLightningLOs = append(shaman.ChainLightningLOs, shaman.newChainLightningSpell(true))
	}

	if shaman.Talents.Thunderstorm {
		shaman.Thunderstorm = shaman.newThunderstormSpell(shaman.thunderstormInRange)
	}

	if shaman.Talents.LavaLash {
		shaman.LavaLash = shaman.newLavaLashSpell()
	}

	shaman.registerFeralSpirit()

	shaman.registerShocks()
	shaman.registerGraceOfAirTotemSpell()
	shaman.registerMagmaTotemSpell()
	shaman.registerManaSpringTotemSpell()
	shaman.registerSearingTotemSpell()
	shaman.registerStrengthOfEarthTotemSpell()
	shaman.registerTotemOfWrathSpell()
	shaman.registerTranquilAirTotemSpell()
	shaman.registerTremorTotemSpell()
	shaman.registerWindfuryTotemSpell()
	shaman.registerWrathOfAirTotemSpell()

	shaman.registerBloodlustCD()
}

func (shaman *Shaman) Reset(sim *core.Simulation) {
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
				shaman.NextTotemDrops[FireTotem] = TotemRefreshTime5M
				if shaman.NextTotemDropType[FireTotem] != int32(proto.FireTotem_TotemOfWrath) {
					shaman.NextTotemDrops[FireTotem] = 0 // attack totems we drop immediately
				} else if shaman.NextTotemDropType[FireTotem] == int32(proto.FireTotem_TotemOfWrath) {
					shaman.applyToWDebuff(sim)
				}
			}
		case WaterTotem:
			if shaman.Totems.Water == proto.WaterTotem_ManaSpringTotem {
				shaman.NextTotemDrops[i] = TotemRefreshTime5M
			}
		}
	}

	shaman.FlameShock.CD.Reset()
}

func (shaman *Shaman) ElementalCritMultiplier(secondary float64) float64 {
	critBonus := (0.2 * float64(shaman.Talents.ElementalFury)) + secondary
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
		stats.Mana:        baseMana,
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
		stats.Mana:        baseMana,
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
		stats.Mana:        baseMana,
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
		stats.Mana:        baseMana,
		stats.SpellCrit:   2.2 * core.CritRatingPerCritChance,
		stats.AttackPower: 95, // TODO: confirm this.
		stats.MeleeCrit:   2.92 * core.CritRatingPerCritChance,
	}
}
