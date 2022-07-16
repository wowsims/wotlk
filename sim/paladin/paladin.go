package paladin

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

const (
	SpellFlagSeal      = core.SpellFlagAgentReserved1
	SpellFlagJudgement = core.SpellFlagAgentReserved2
)

type hybridScaling struct {
	AP float64
	SP float64
}

type Paladin struct {
	core.Character

	PaladinAura proto.PaladinAura

	Talents proto.PaladinTalents

	CurrentSeal      *core.Aura
	CurrentJudgement *core.Aura

	DivineStorm       *core.Spell
	Consecration      *core.Spell
	CrusaderStrike    *core.Spell
	Exorcism          *core.Spell
	HolyShield        *core.Spell
	JudgementOfWisdom *core.Spell
	JudgementOfLight  *core.Spell
	SealOfVengeance   *core.Spell
	// SealOfWisdom        *core.Spell
	// SealOfLight         *core.Spell
	// SealOfRighteousness *core.Spell

	ConsecrationDot *core.Dot
	// SealOfVengeanceDot *core.Dot

	HolyShieldAura        *core.Aura
	JudgementOfWisdomAura *core.Aura
	JudgementOfLightAura  *core.Aura
	SealOfVengeanceAura   *core.Aura
	// SealOfCommandAura       *core.Aura
	// SealOfWisdomAura        *core.Aura
	// SealOfLightAura         *core.Aura
	// SealOfRighteousnessAura *core.Aura

	ArtOfWarInstantCast *core.Aura

	SpiritualAttunementMetrics *core.ResourceMetrics
}

// Implemented by each Paladin spec.
type PaladinAgent interface {
	GetPaladin() *Paladin
}

func (paladin *Paladin) GetCharacter() *core.Character {
	return &paladin.Character
}

func (paladin *Paladin) HasMajorGlyph(glyph proto.PaladinMajorGlyph) bool {
	return paladin.HasGlyph(int32(glyph))
}
func (paladin *Paladin) HasMinorGlyph(glyph proto.PaladinMinorGlyph) bool {
	return paladin.HasGlyph(int32(glyph))
}

func (paladin *Paladin) GetPaladin() *Paladin {
	return paladin
}

func (paladin *Paladin) AddRaidBuffs(raidBuffs *proto.RaidBuffs) {
	raidBuffs.DevotionAura = core.MaxTristate(raidBuffs.DevotionAura, core.MakeTristateValue(
		paladin.PaladinAura == proto.PaladinAura_DevotionAura,
		paladin.Talents.ImprovedDevotionAura == 5))

	// TODO: Fix
	raidBuffs.RetributionAura = core.MaxTristate(raidBuffs.RetributionAura, core.MakeTristateValue(
		paladin.PaladinAura == proto.PaladinAura_RetributionAura,
		paladin.Talents.SanctifiedRetribution == true))

	//if paladin.Talents.SanctifiedRetribution {
	//	raidBuffs.SanctifiedRetribution = true
	//}
}

func (paladin *Paladin) AddPartyBuffs(partyBuffs *proto.PartyBuffs) {
}

func (paladin *Paladin) Initialize() {
	// Update auto crit multipliers now that we have the targets.
	paladin.AutoAttacks.MHEffect.OutcomeApplier = paladin.OutcomeFuncMeleeWhite(paladin.MeleeCritMultiplier())

	paladin.setupSealOfVengeance()
	// paladin.setupSealOfTheCrusader()
	// paladin.setupSealOfWisdom()
	// paladin.setupSealOfLight()
	// paladin.setupSealOfRighteousness()
	// paladin.setupJudgementRefresh()

	paladin.registerCrusaderStrikeSpell()
	paladin.registerDivineStormSpell()
	paladin.RegisterConsecrationSpell(8)

	paladin.registerExorcismSpell()
	paladin.registerHolyShieldSpell()
	paladin.registerJudgements()

	paladin.registerSpiritualAttunement()
}

func (paladin *Paladin) Reset(sim *core.Simulation) {
	paladin.CurrentSeal = nil
	paladin.CurrentJudgement = nil
}

func (paladin *Paladin) OnAutoAttack(sim *core.Simulation, spell *core.Spell) {
}

// maybe need to add stat dependencies
func NewPaladin(character core.Character, talents proto.PaladinTalents) *Paladin {
	paladin := &Paladin{
		Character: character,
		Talents:   talents,
	}

	paladin.PseudoStats.CanParry = true

	paladin.EnableManaBar()

	// Add paladin stat dependencies
	paladin.AddStatDependency(stats.StatDependency{
		SourceStat:   stats.Strength,
		ModifiedStat: stats.AttackPower,
		Modifier: func(strength float64, attackPower float64) float64 {
			return attackPower + strength*2
		},
	})

	paladin.AddStatDependency(stats.StatDependency{
		SourceStat:   stats.Agility,
		ModifiedStat: stats.MeleeCrit,
		Modifier: func(agility float64, meleeCrit float64) float64 {
			return meleeCrit + (agility/25)*core.CritRatingPerCritChance
		},
	})

	paladin.AddStatDependency(stats.StatDependency{
		SourceStat:   stats.Agility,
		ModifiedStat: stats.Dodge,
		Modifier: func(agility float64, dodge float64) float64 {
			return dodge + (agility/25)*core.DodgeRatingPerDodgeChance
		},
	})

	return paladin
}

func init() {
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceBloodElf, Class: proto.Class_ClassPaladin}] = stats.Stats{
		stats.Health:      3197,
		stats.Stamina:     118,
		stats.Intellect:   87,
		stats.Mana:        2953,
		stats.Spirit:      88,
		stats.Strength:    123,
		stats.AttackPower: 190,
		stats.Agility:     79,
		stats.MeleeCrit:   14.35,
		stats.SpellCrit:   73.69,
		stats.Dodge:       0.65 * core.DodgeRatingPerDodgeChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceDraenei, Class: proto.Class_ClassPaladin}] = stats.Stats{
		stats.Health:      3197,
		stats.Stamina:     119,
		stats.Intellect:   84,
		stats.Mana:        2953,
		stats.Spirit:      91,
		stats.Strength:    127,
		stats.AttackPower: 190,
		stats.Agility:     74,
		stats.MeleeCrit:   14.35,
		stats.SpellCrit:   73.69,
		stats.Dodge:       0.65 * core.DodgeRatingPerDodgeChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceHuman, Class: proto.Class_ClassPaladin}] = stats.Stats{
		stats.Health:      3197,
		stats.Stamina:     120,
		stats.Intellect:   83,
		stats.Mana:        2953,
		stats.Spirit:      97,
		stats.Strength:    126,
		stats.AttackPower: 190,
		stats.Agility:     77,
		stats.MeleeCrit:   14.35,
		stats.SpellCrit:   73.69,
		stats.Dodge:       0.65 * core.DodgeRatingPerDodgeChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceDwarf, Class: proto.Class_ClassPaladin}] = stats.Stats{
		stats.Health:      3197,
		stats.Stamina:     123,
		stats.Intellect:   82,
		stats.Mana:        2953,
		stats.Spirit:      88,
		stats.Strength:    128,
		stats.AttackPower: 190,
		stats.Agility:     73,
		stats.MeleeCrit:   14.35,
		stats.SpellCrit:   73.69,
		stats.Dodge:       0.65 * core.DodgeRatingPerDodgeChance,
	}
}
