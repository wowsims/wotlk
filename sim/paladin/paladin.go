package paladin

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

const (
	SpellFlagSecondaryJudgement = core.SpellFlagAgentReserved1
	SpellFlagPrimaryJudgement   = core.SpellFlagAgentReserved2
)

type hybridScaling struct {
	AP float64
	SP float64
}

type Modifier []float64
type Modifiers []Modifier

func (mod *Modifiers) Get() float64 {
	value := 1.0
	for _, m := range *mod {
		sum := 1.0

		// Combine additive bonuses.
		for _, a := range m {
			sum += a
		}

		// Combine multiplicative bonuses.
		value *= sum
	}
	return value
}

func (mod *Modifiers) Clone() Modifiers {
	return (*mod)[:]
}

type Paladin struct {
	core.Character

	PaladinAura proto.PaladinAura

	Talents proto.PaladinTalents

	CurrentSeal      *core.Aura
	CurrentJudgement *core.Aura

	DivinePlea          *core.Spell
	DivineStorm         *core.Spell
	Consecration        *core.Spell
	CrusaderStrike      *core.Spell
	Exorcism            *core.Spell
	HolyShield          *core.Spell
	JudgementOfWisdom   *core.Spell
	JudgementOfLight    *core.Spell
	SealOfVengeance     *core.Spell
	SealOfRighteousness *core.Spell
	SealOfCommand       *core.Spell
	// SealOfWisdom        *core.Spell
	// SealOfLight         *core.Spell

	ConsecrationDot *core.Dot
	// SealOfVengeanceDot *core.Dot

	HolyShieldAura *core.Aura
	// RighteousFuryAura       *core.Aura
	JudgementOfWisdomAura   *core.Aura
	JudgementOfLightAura    *core.Aura
	SealOfVengeanceAura     *core.Aura
	SealOfCommandAura       *core.Aura
	SealOfRighteousnessAura *core.Aura

	// SealOfWisdomAura        *core.Aura
	// SealOfLightAura         *core.Aura

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

	if paladin.PaladinAura == proto.PaladinAura_RetributionAura {
		raidBuffs.RetributionAura = true
	}

	if paladin.Talents.SanctifiedRetribution {
		raidBuffs.SanctifiedRetribution = true
	}

	if paladin.Talents.SwiftRetribution == 3 {
		raidBuffs.SwiftRetribution = paladin.Talents.SwiftRetribution == 3 // TODO: Fix-- though having something between 0/3 and 3/3 is unlikely
	}
}

func (paladin *Paladin) AddPartyBuffs(partyBuffs *proto.PartyBuffs) {
}

func (paladin *Paladin) Initialize() {
	// Update auto crit multipliers now that we have the targets.
	paladin.AutoAttacks.MHEffect.OutcomeApplier = paladin.OutcomeFuncMeleeWhite(paladin.MeleeCritMultiplier())

	paladin.registerSealOfVengeanceSpellAndAura()
	paladin.registerSealOfRighteousnessSpellAndAura()
	paladin.registerSealOfCommandSpellAndAura()
	// paladin.setupSealOfTheCrusader()
	// paladin.setupSealOfWisdom()
	// paladin.setupSealOfLight()
	// paladin.setupSealOfRighteousness()
	// paladin.setupJudgementRefresh()

	paladin.registerCrusaderStrikeSpell()
	paladin.registerDivineStormSpell()
	paladin.registerConsecrationSpell()

	paladin.registerExorcismSpell()
	paladin.registerHolyShieldSpell()
	paladin.registerJudgements()

	paladin.registerSpiritualAttunement()
	paladin.registerDivinePleaSpell()
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
	paladin.AddStatDependency(stats.Strength, stats.AttackPower, 1.0+2)
	paladin.AddStatDependency(stats.Agility, stats.MeleeCrit, 1.0+(core.CritRatingPerCritChance/25))
	paladin.AddStatDependency(stats.Agility, stats.Dodge, 1.0+(core.DodgeRatingPerDodgeChance/25))

	return paladin
}

func init() {
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceBloodElf, Class: proto.Class_ClassPaladin}] = stats.Stats{
		stats.Health:      8164,
		stats.Stamina:     141,
		stats.Intellect:   102,
		stats.Mana:        5644,
		stats.Spirit:      104,
		stats.Strength:    148,
		stats.AttackPower: 516,
		stats.Agility:     92,
		stats.MeleeCrit:   5.03 * core.CritRatingPerCritChance,
		stats.SpellCrit:   3.95 * core.CritRatingPerCritChance,
		stats.Dodge:       5.03 * core.DodgeRatingPerDodgeChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceDraenei, Class: proto.Class_ClassPaladin}] = stats.Stats{
		stats.Health:      8174,
		stats.Stamina:     142,
		stats.Intellect:   113,
		stats.Mana:        5809,
		stats.Spirit:      107,
		stats.Strength:    152,
		stats.AttackPower: 524,
		stats.Agility:     87,
		stats.MeleeCrit:   4.92 * core.CritRatingPerCritChance,
		stats.SpellCrit:   4.01 * core.CritRatingPerCritChance,
		stats.Dodge:       4.95 * core.DodgeRatingPerDodgeChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceHuman, Class: proto.Class_ClassPaladin}] = stats.Stats{
		stats.Health:      8354,
		stats.Stamina:     160,
		stats.Intellect:   98,
		stats.Mana:        5584,
		stats.Spirit:      113,
		stats.Strength:    173,
		stats.AttackPower: 566,
		stats.Agility:     90,
		stats.MeleeCrit:   5 * core.CritRatingPerCritChance,
		stats.SpellCrit:   3.92 * core.CritRatingPerCritChance,
		stats.Dodge:       5 * core.DodgeRatingPerDodgeChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceDwarf, Class: proto.Class_ClassPaladin}] = stats.Stats{
		stats.Health:      8214,
		stats.Stamina:     146,
		stats.Intellect:   97,
		stats.Mana:        5569,
		stats.Spirit:      104,
		stats.Strength:    175,
		stats.AttackPower: 570,
		stats.Agility:     86,
		stats.MeleeCrit:   4.92 * core.CritRatingPerCritChance,
		stats.SpellCrit:   3.92 * core.CritRatingPerCritChance,
		stats.Dodge:       4.93 * core.DodgeRatingPerDodgeChance,
	}
}
