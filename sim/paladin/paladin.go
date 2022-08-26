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

type Additive []float64
type Multiplicative []Additive

func (mod *Additive) Get() float64 {
	sum := 1.0
	// Combine additive bonuses.
	for _, value := range *mod {
		sum += value
	}
	return sum
}

func (mod *Multiplicative) Get() float64 {
	multiplier := 1.0
	// Combine multiplicative bonuses.
	for _, additive := range *mod {
		multiplier *= additive.Get()
	}
	return multiplier
}

func (mod *Multiplicative) Clone() Multiplicative {
	return (*mod)[:]
}

type Paladin struct {
	core.Character

	PaladinAura proto.PaladinAura

	Talents proto.PaladinTalents

	CurrentSeal      *core.Aura
	CurrentJudgement *core.Aura

	DivinePlea            *core.Spell
	DivineStorm           *core.Spell
	HolyWrath             *core.Spell
	Consecration          *core.Spell
	CrusaderStrike        *core.Spell
	Exorcism              *core.Spell
	HolyShield            *core.Spell
	HammerOfTheRighteous  *core.Spell
	ShieldOfRighteousness *core.Spell
	AvengersShield        *core.Spell
	JudgementOfWisdom     *core.Spell
	JudgementOfLight      *core.Spell
	HammerOfWrath         *core.Spell
	SealOfVengeance       *core.Spell
	SealOfRighteousness   *core.Spell
	SealOfCommand         *core.Spell
	// SealOfWisdom        *core.Spell
	// SealOfLight         *core.Spell

	ConsecrationDot     *core.Dot
	SealOfVengeanceDots []*core.Dot

	HolyShieldAura *core.Aura
	// RighteousFuryAura       *core.Aura
	DivinePleaAura          *core.Aura
	JudgementOfWisdomAura   *core.Aura
	JudgementOfLightAura    *core.Aura
	SealOfVengeanceAura     *core.Aura
	SealOfCommandAura       *core.Aura
	SealOfRighteousnessAura *core.Aura

	// SealOfWisdomAura        *core.Aura
	// SealOfLightAura         *core.Aura

	RighteousVengeanceSpell  *core.Spell
	RighteousVengeanceDots   []*core.Dot
	RighteousVengeancePools  []float64
	RighteousVengeanceDamage []float64

	ArtOfWarInstantCast *core.Aura

	SpiritualAttunementMetrics *core.ResourceMetrics

	HasTuralyonsOrLiadrinsBattlegear2Pc bool

	DemonAndUndeadTargetCount int32
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
	paladin.registerHammerOfWrathSpell()
	paladin.registerHolyWrathSpell()

	paladin.registerExorcismSpell()
	paladin.registerHolyShieldSpell()
	paladin.registerHammerOfTheRighteousSpell()
	paladin.registerShieldOfRighteousnessSpell()
	paladin.registerAvengersShieldSpell()
	paladin.registerJudgements()

	paladin.registerSpiritualAttunement()
	paladin.registerDivinePleaSpell()
	paladin.registerRighteousVengeanceSpell()

	targets := paladin.Env.GetNumTargets()

	if paladin.Talents.RighteousVengeance > 0 {
		paladin.RighteousVengeanceDots = []*core.Dot{}
		for i := int32(0); i < targets; i++ {
			paladin.RighteousVengeanceDots = append(paladin.RighteousVengeanceDots, paladin.makeRighteousVengeanceDot(paladin.Env.GetTargetUnit(i)))
		}
		paladin.RighteousVengeancePools = []float64{}
		for i := int32(0); i < targets; i++ {
			paladin.RighteousVengeancePools = append(paladin.RighteousVengeancePools, 0.0)
		}
		paladin.RighteousVengeanceDamage = []float64{}
		for i := int32(0); i < targets; i++ {
			paladin.RighteousVengeanceDamage = append(paladin.RighteousVengeanceDamage, 0.0)
		}
	}

	paladin.SealOfVengeanceDots = []*core.Dot{}
	for i := int32(0); i < targets; i++ {
		paladin.SealOfVengeanceDots = append(paladin.SealOfVengeanceDots, paladin.createSealOfVengeanceDot(paladin.Env.GetTargetUnit(i)))
	}

	for i := int32(0); i < paladin.Env.GetNumTargets(); i++ {
		unit := paladin.Env.GetTargetUnit(i)
		if unit.MobType == proto.MobType_MobTypeDemon || unit.MobType == proto.MobType_MobTypeUndead {
			paladin.DemonAndUndeadTargetCount += 1
		}
	}
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

	paladin.HasTuralyonsOrLiadrinsBattlegear2Pc = paladin.HasSetBonus(ItemSetTuralyonsBattlegear, 2) || paladin.HasSetBonus(ItemSetLiadrinsBattlegear, 2)

	paladin.PseudoStats.CanParry = true

	paladin.EnableManaBar()

	// Paladins get 3 times their level in base AP
	// then 2 AP per STR, then lose the first 20 AP
	paladin.AddStatDependency(stats.Strength, stats.AttackPower, 2.0)
	paladin.AddStat(stats.AttackPower, -20)

	// Paladins get 1% crit per 52.08 agil
	paladin.AddStatDependency(stats.Agility, stats.MeleeCrit, (1.0/52.08)*core.CritRatingPerCritChance)

	// Paladins get 1% dodge per 52.08 agil
	paladin.AddStatDependency(stats.Agility, stats.Dodge, (1.0/52.08)*core.DodgeRatingPerDodgeChance)

	// Paladins get more melee haste from haste than other classes, 25.22/1%
	paladin.PseudoStats.MeleeHasteRatingPerHastePercent = 25.22

	return paladin
}

func init() {
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceBloodElf, Class: proto.Class_ClassPaladin}] = stats.Stats{
		stats.Health:      6754,
		stats.Stamina:     141,
		stats.Intellect:   102,
		stats.Mana:        4394,
		stats.Spirit:      104,
		stats.Strength:    148,
		stats.AttackPower: 240,
		stats.Agility:     92,
		stats.MeleeCrit:   3.27 * core.CritRatingPerCritChance,
		stats.SpellCrit:   3.27 * core.CritRatingPerCritChance,
		stats.Dodge:       3.27 * core.DodgeRatingPerDodgeChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceDraenei, Class: proto.Class_ClassPaladin}] = stats.Stats{
		stats.Health:      6754,
		stats.Stamina:     142,
		stats.Intellect:   113,
		stats.Mana:        4394,
		stats.Spirit:      107,
		stats.Strength:    152,
		stats.AttackPower: 240,
		stats.Agility:     87,
		stats.MeleeCrit:   3.27 * core.CritRatingPerCritChance,
		stats.SpellCrit:   3.27 * core.CritRatingPerCritChance,
		stats.Dodge:       3.27 * core.DodgeRatingPerDodgeChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceHuman, Class: proto.Class_ClassPaladin}] = stats.Stats{
		stats.Health:      6754,
		stats.Stamina:     160,
		stats.Intellect:   98,
		stats.Mana:        4394,
		stats.Spirit:      113,
		stats.Strength:    173,
		stats.AttackPower: 240,
		stats.Agility:     90,
		stats.MeleeCrit:   3.27 * core.CritRatingPerCritChance,
		stats.SpellCrit:   3.27 * core.CritRatingPerCritChance,
		stats.Dodge:       3.27 * core.DodgeRatingPerDodgeChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceDwarf, Class: proto.Class_ClassPaladin}] = stats.Stats{
		stats.Health:      6754,
		stats.Stamina:     146,
		stats.Intellect:   97,
		stats.Mana:        4394,
		stats.Spirit:      104,
		stats.Strength:    175,
		stats.AttackPower: 240,
		stats.Agility:     86,
		stats.MeleeCrit:   3.27 * core.CritRatingPerCritChance,
		stats.SpellCrit:   3.27 * core.CritRatingPerCritChance,
		stats.Dodge:       3.27 * core.DodgeRatingPerDodgeChance,
	}
}
