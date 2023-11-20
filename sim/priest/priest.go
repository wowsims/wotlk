package priest

import (
	"github.com/wowsims/classic/sim/core"
	"github.com/wowsims/classic/sim/core/proto"
)

var TalentTreeSizes = [3]int{15, 16, 16}

const (
	PriestRuneChestTwistedFaith    = 425210
	PriestRuneChestVoidPlague      = 425211
	PriestRuneHandsPenance         = 402844
	PriestRuneHandsShadowWordDeath = 402833
	PriestRuneLegsHomunculi        = 402836
	PriestRuneLegsSharedPain       = 402838
)

type Priest struct {
	core.Character
	Talents *proto.PriestTalents

	Latency float64

	InnerFocusAura    *core.Aura
	ShadowWeavingAura *core.Aura

	BindingHeal       *core.Spell
	CircleOfHealing   *core.Spell
	DevouringPlague   *core.Spell
	FlashHeal         *core.Spell
	GreaterHeal       *core.Spell
	HolyFire          *core.Spell
	InnerFocus        *core.Spell
	ShadowWordPain    *core.Spell
	MindBlast         *core.Spell
	MindFlay          *core.Spell
	MindFlayModifier  float64
	MindBlastModifier float64
	MindSear          *core.Spell
	Penance           *core.Spell
	PenanceHeal       *core.Spell
	PowerWordShield   *core.Spell
	PrayerOfHealing   *core.Spell
	PrayerOfMending   *core.Spell
	Renew             *core.Spell
	EmpoweredRenew    *core.Spell
	ShadowWordDeath   *core.Spell
	Smite             *core.Spell
	VoidPlague        *core.Spell

	WeakenedSouls core.AuraArray

	ProcPrayerOfMending core.ApplySpellResults

	DpInitMultiplier float64
}

func (priest *Priest) GetCharacter() *core.Character {
	return &priest.Character
}

func (priest *Priest) AddRaidBuffs(raidBuffs *proto.RaidBuffs) {
	raidBuffs.ShadowProtection = true
	raidBuffs.DivineSpirit = true

	raidBuffs.PowerWordFortitude = max(raidBuffs.PowerWordFortitude, core.MakeTristateValue(
		true,
		priest.Talents.ImprovedPowerWordFortitude == 2))
}

func (priest *Priest) AddPartyBuffs(_ *proto.PartyBuffs) {
}

func (priest *Priest) Initialize() {
	// Shadow Insight gained from Glyph of Shadow
	// Finalized spirit off gear and not dynamic spirit (e.g. Spirit Tap does not increase this)
	priest.registerSetBonuses()
	priest.registerMindBlast()
	priest.registerMindFlay()
	priest.registerShadowWordPainSpell()
	priest.registerDevouringPlagueSpell()
	priest.RegisterSmiteSpell()
	priest.registerHolyFire()
	// TODO: Classic implement runes
	priest.registerVoidPlagueSpell()
	priest.RegisterPenanceSpell()
	priest.registerShadowWordDeathSpell()

	priest.registerPowerInfusionCD()
}

func (priest *Priest) RegisterHealingSpells() {
	priest.registerPenanceHealSpell()
	priest.registerBindingHealSpell()
	priest.registerCircleOfHealingSpell()
	priest.registerFlashHealSpell()
	priest.registerGreaterHealSpell()
	priest.registerPowerWordShieldSpell()
	priest.registerPrayerOfHealingSpell()
	priest.registerPrayerOfMendingSpell()
	priest.registerRenewSpell()
}

func (priest *Priest) AddShadowWeavingStack(sim *core.Simulation) {
	if priest.ShadowWeavingAura == nil {
		return
	}

	if sim.RollWithLabel(0, 1, "ShadowWeaving") < (0.2 * float64(priest.Talents.ShadowWeaving)) {
		priest.ShadowWeavingAura.Activate(sim)
		priest.ShadowWeavingAura.AddStack(sim)
	}
}

func (priest *Priest) Reset(_ *core.Simulation) {
	priest.MindFlayModifier = 1
	priest.MindBlastModifier = 1
}

func New(char *core.Character, talents string) *Priest {
	priest := &Priest{
		Character: *char,
		Talents:   &proto.PriestTalents{},
	}
	core.FillTalentsProto(priest.Talents.ProtoReflect(), talents, TalentTreeSizes)

	priest.EnableManaBar()

	return priest
}

// Agent is a generic way to access underlying priest on any of the agents.
type PriestAgent interface {
	GetPriest() *Priest
}
