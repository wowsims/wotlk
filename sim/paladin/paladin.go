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

var TalentTreeSizes = [3]int{26, 26, 26}

type Paladin struct {
	core.Character

	PaladinAura proto.PaladinAura

	Talents *proto.PaladinTalents

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
	HandOfReckoning       *core.Spell
	ShieldOfRighteousness *core.Spell
	AvengersShield        *core.Spell
	JudgementOfWisdom     *core.Spell
	JudgementOfLight      *core.Spell
	HammerOfWrath         *core.Spell
	SealOfVengeance       *core.Spell
	SealOfRighteousness   *core.Spell
	SealOfCommand         *core.Spell
	AvengingWrath         *core.Spell
	DivineProtection      *core.Spell
	SovDotSpell           *core.Spell
	// SealOfWisdom        *core.Spell
	// SealOfLight         *core.Spell

	HolyShieldAura *core.Aura
	// RighteousFuryAura       *core.Aura
	DivinePleaAura          *core.Aura
	JudgementOfWisdomAura   *core.Aura
	JudgementOfLightAura    *core.Aura
	SealOfVengeanceAura     *core.Aura
	SealOfCommandAura       *core.Aura
	SealOfRighteousnessAura *core.Aura
	AvengingWrathAura       *core.Aura
	DivineProtectionAura    *core.Aura
	ForbearanceAura         *core.Aura

	// SealOfWisdomAura        *core.Aura
	// SealOfLightAura         *core.Aura

	ArtOfWarInstantCast *core.Aura

	SpiritualAttunementMetrics *core.ResourceMetrics

	HasTuralyonsOrLiadrinsBattlegear2Pc bool

	DemonAndUndeadTargetCount int32

	AvoidClippingConsecration           bool
	HoldLastAvengingWrathUntilExecution bool

	mutualLockoutDPAW *core.Timer
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

	// TODO: Figure out a way to just start with 1 DG cooldown available without making a redundant Spell
	//if paladin.Talents.DivineGuardian == 2 {
	//	raidBuffs.divineGuardians++
	//}
}

func (paladin *Paladin) AddPartyBuffs(_ *proto.PartyBuffs) {
}

func (paladin *Paladin) Initialize() {
	// Update auto crit multipliers now that we have the targets.
	paladin.AutoAttacks.MHConfig.CritMultiplier = paladin.MeleeCritMultiplier()

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
	paladin.registerHandOfReckoningSpell()
	paladin.registerShieldOfRighteousnessSpell()
	paladin.registerAvengersShieldSpell()
	paladin.registerJudgements()

	paladin.registerSpiritualAttunement()
	paladin.registerDivinePleaSpell()
	paladin.registerDivineProtectionSpell()
	paladin.registerForbearanceDebuff()

	for i := int32(0); i < paladin.Env.GetNumTargets(); i++ {
		unit := paladin.Env.GetTargetUnit(i)
		if unit.MobType == proto.MobType_MobTypeDemon || unit.MobType == proto.MobType_MobTypeUndead {
			paladin.DemonAndUndeadTargetCount += 1
		}
	}
}

func (paladin *Paladin) Reset(_ *core.Simulation) {
	paladin.CurrentSeal = nil
	paladin.CurrentJudgement = nil
}

// maybe need to add stat dependencies
func NewPaladin(character core.Character, talentsStr string) *Paladin {
	paladin := &Paladin{
		Character: character,
		Talents:   &proto.PaladinTalents{},
	}
	core.FillTalentsProto(paladin.Talents.ProtoReflect(), talentsStr, TalentTreeSizes)

	// This is used to cache its effect in talents.go
	paladin.HasTuralyonsOrLiadrinsBattlegear2Pc = paladin.HasSetBonus(ItemSetTuralyonsBattlegear, 2)

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

	// Paladins get 1 block value per 2 str
	paladin.AddStatDependency(stats.Strength, stats.BlockValue, .5)

	// Base dodge is unaffected by Diminishing Returns
	paladin.PseudoStats.BaseDodge += 0.0327
	paladin.PseudoStats.BaseParry += 0.05

	return paladin
}

// Shared 30sec cooldown for Divine Protection and Avenging Wrath
func (paladin *Paladin) GetMutualLockoutDPAW() *core.Timer {
	return paladin.Character.GetOrInitTimer(&paladin.mutualLockoutDPAW)
}

func init() {
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceSindorei, Class: proto.Class_ClassPaladin}] = stats.Stats{
		stats.Health:      6754,
		stats.Stamina:     143,
		stats.Intellect:   101,
		stats.Mana:        4394,
		stats.Spirit:      103,
		stats.Strength:    148,
		stats.AttackPower: 240,
		stats.Agility:     92,
		stats.MeleeCrit:   3.269 * core.CritRatingPerCritChance,
		stats.SpellCrit:   3.269 * core.CritRatingPerCritChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceDraenei, Class: proto.Class_ClassPaladin}] = stats.Stats{
		stats.Health:      6754,
		stats.Stamina:     143,
		stats.Intellect:   98,
		stats.Mana:        4394,
		stats.Spirit:      107,
		stats.Strength:    152,
		stats.AttackPower: 240,
		stats.Agility:     87,
		stats.MeleeCrit:   3.269 * core.CritRatingPerCritChance,
		stats.SpellCrit:   3.269 * core.CritRatingPerCritChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceHuman, Class: proto.Class_ClassPaladin}] = stats.Stats{
		stats.Health:      6754,
		stats.Stamina:     143,
		stats.Intellect:   98,
		stats.Mana:        4394,
		stats.Spirit:      105,
		stats.Strength:    151,
		stats.AttackPower: 240,
		stats.Agility:     90,
		stats.MeleeCrit:   3.269 * core.CritRatingPerCritChance,
		stats.SpellCrit:   3.269 * core.CritRatingPerCritChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceDwarf, Class: proto.Class_ClassPaladin}] = stats.Stats{
		stats.Health:      6754,
		stats.Stamina:     144,
		stats.Intellect:   97,
		stats.Mana:        4394,
		stats.Spirit:      104,
		stats.Strength:    156,
		stats.AttackPower: 240,
		stats.Agility:     86,
		stats.MeleeCrit:   3.269 * core.CritRatingPerCritChance,
		stats.SpellCrit:   3.269 * core.CritRatingPerCritChance,
	}
}
