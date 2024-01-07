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

	HolyShieldAura          *core.Aura
	RighteousFuryAura       *core.Aura
	DivinePleaAura          *core.Aura
	JudgementOfWisdomAura   *core.Aura
	JudgementOfLightAura    *core.Aura
	SealOfVengeanceAura     *core.Aura
	SealOfCommandAura       *core.Aura
	SealOfRighteousnessAura *core.Aura
	AvengingWrathAura       *core.Aura
	DivineProtectionAura    *core.Aura
	ForbearanceAura         *core.Aura
	VengeanceAura           *core.Aura

	// SealOfWisdomAura        *core.Aura
	// SealOfLightAura         *core.Aura

	ArtOfWarInstantCast *core.Aura

	SpiritualAttunementMetrics *core.ResourceMetrics

	HasTuralyonsOrLiadrinsBattlegear2Pc bool

	DemonAndUndeadTargetCount int32

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
	raidBuffs.DevotionAura = max(raidBuffs.DevotionAura, core.MakeTristateValue(
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
	paladin.AutoAttacks.MHConfig().CritMultiplier = paladin.MeleeCritMultiplier()

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
func NewPaladin(character *core.Character, talentsStr string) *Paladin {
	paladin := &Paladin{
		Character: *character,
		Talents:   &proto.PaladinTalents{},
	}
	core.FillTalentsProto(paladin.Talents.ProtoReflect(), talentsStr, TalentTreeSizes)

	// This is used to cache its effect in talents.go
	paladin.HasTuralyonsOrLiadrinsBattlegear2Pc = paladin.HasSetBonus(ItemSetTuralyonsBattlegear, 2)

	paladin.PseudoStats.CanParry = true

	paladin.EnableManaBar()
	paladin.AddStatDependency(stats.Strength, stats.AttackPower, 2.0)
	paladin.AddStatDependency(stats.Agility, stats.MeleeCrit, core.CritPerAgiMaxLevel[character.Class]*core.CritRatingPerCritChance)

	// Paladins get 0.0167 dodge per agi. ~1% per 59.88
	paladin.AddStatDependency(stats.Agility, stats.Dodge, (1.0/59.88)*core.DodgeRatingPerDodgeChance)

	// Paladins get more melee haste from haste than other classes
	paladin.PseudoStats.MeleeHasteRatingPerHastePercent /= 1.3

	// Paladins get 1 block value per 2 str
	paladin.AddStatDependency(stats.Strength, stats.BlockValue, .5)

	// Bonus Armor and Armor are treated identically for Paladins
	paladin.AddStatDependency(stats.BonusArmor, stats.Armor, 1)

	// Base dodge is unaffected by Diminishing Returns
	paladin.PseudoStats.BaseDodge += 0.034943
	paladin.PseudoStats.BaseParry += 0.05

	return paladin
}

// Shared 30sec cooldown for Divine Protection and Avenging Wrath
func (paladin *Paladin) GetMutualLockoutDPAW() *core.Timer {
	return paladin.Character.GetOrInitTimer(&paladin.mutualLockoutDPAW)
}
