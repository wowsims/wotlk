package rogue

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func RegisterRogue() {
	core.RegisterAgentFactory(
		proto.Player_Rogue{},
		proto.Spec_SpecRogue,
		func(character *core.Character, options *proto.Player) core.Agent {
			return NewRogue(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_Rogue)
			if !ok {
				panic("Invalid spec value for Rogue!")
			}
			player.Spec = playerSpec
		},
	)
}

const (
	SpellFlagBuilder     = core.SpellFlagAgentReserved2
	SpellFlagFinisher    = core.SpellFlagAgentReserved3
	SpellFlagColdBlooded = core.SpellFlagAgentReserved4
)

var TalentTreeSizes = [3]int{27, 28, 28}

const RogueBleedTag = "RogueBleed"

type Rogue struct {
	core.Character

	Talents *proto.RogueTalents
	Options *proto.Rogue_Options

	bleedCategory *core.ExclusiveCategory

	sliceAndDiceDurations [6]time.Duration
	exposeArmorDurations  [6]time.Duration

	Backstab         *core.Spell
	BladeFlurry      *core.Spell
	DeadlyPoison     *core.Spell
	FanOfKnives      *core.Spell
	Feint            *core.Spell
	Garrote          *core.Spell
	Ambush           *core.Spell
	Hemorrhage       *core.Spell
	GhostlyStrike    *core.Spell
	HungerForBlood   *core.Spell
	InstantPoison    [3]*core.Spell
	WoundPoison      [3]*core.Spell
	Mutilate         *core.Spell
	MutilateMH       *core.Spell
	MutilateOH       *core.Spell
	Shiv             *core.Spell
	SinisterStrike   *core.Spell
	TricksOfTheTrade *core.Spell
	Shadowstep       *core.Spell
	Preparation      *core.Spell
	Premeditation    *core.Spell
	ShadowDance      *core.Spell
	ColdBlood        *core.Spell
	Vanish           *core.Spell

	Envenom      *core.Spell
	Eviscerate   *core.Spell
	ExposeArmor  *core.Spell
	Rupture      *core.Spell
	SliceAndDice *core.Spell

	lastDeadlyPoisonProcMask core.ProcMask

	deadlyPoisonProcChanceBonus float64
	instantPoisonPPMM           core.PPMManager
	woundPoisonPPMM             core.PPMManager

	AdrenalineRushAura   *core.Aura
	BladeFlurryAura      *core.Aura
	EnvenomAura          *core.Aura
	ExposeArmorAuras     core.AuraArray
	HungerForBloodAura   *core.Aura
	KillingSpreeAura     *core.Aura
	OverkillAura         *core.Aura
	SliceAndDiceAura     *core.Aura
	MasterOfSubtletyAura *core.Aura
	ShadowstepAura       *core.Aura
	ShadowDanceAura      *core.Aura
	DirtyDeedsAura       *core.Aura
	HonorAmongThieves    *core.Aura
	StealthAura          *core.Aura

	masterPoisonerDebuffAuras core.AuraArray
	savageCombatDebuffAuras   core.AuraArray
	woundPoisonDebuffAuras    core.AuraArray

	QuickRecoveryMetrics *core.ResourceMetrics

	costModifier               func(float64) float64
	finishingMoveEffectApplier func(sim *core.Simulation, numPoints int32)
}

func (rogue *Rogue) GetCharacter() *core.Character {
	return &rogue.Character
}

func (rogue *Rogue) GetRogue() *Rogue {
	return rogue
}

func (rogue *Rogue) AddRaidBuffs(_ *proto.RaidBuffs)   {}
func (rogue *Rogue) AddPartyBuffs(_ *proto.PartyBuffs) {}

func (rogue *Rogue) finisherFlags() core.SpellFlag {
	flags := SpellFlagFinisher
	if rogue.Talents.SurpriseAttacks {
		flags |= core.SpellFlagCannotBeDodged
	}
	return flags
}

// Apply the effect of successfully casting a finisher to combo points
func (rogue *Rogue) ApplyFinisher(sim *core.Simulation, spell *core.Spell) {
	numPoints := rogue.ComboPoints()
	rogue.SpendComboPoints(sim, spell.ComboPointMetrics())
	rogue.finishingMoveEffectApplier(sim, numPoints)
}

func (rogue *Rogue) HasMajorGlyph(glyph proto.RogueMajorGlyph) bool {
	return rogue.HasGlyph(int32(glyph))
}

func (rogue *Rogue) HasMinorGlyph(glyph proto.RogueMinorGlyph) bool {
	return rogue.HasGlyph(int32(glyph))
}

func (rogue *Rogue) Initialize() {
	// Update auto crit multipliers now that we have the targets.
	rogue.AutoAttacks.MHConfig().CritMultiplier = rogue.MeleeCritMultiplier(false)
	rogue.AutoAttacks.OHConfig().CritMultiplier = rogue.MeleeCritMultiplier(false)

	if rogue.Talents.QuickRecovery > 0 {
		rogue.QuickRecoveryMetrics = rogue.NewEnergyMetrics(core.ActionID{SpellID: 31245})
	}

	rogue.costModifier = rogue.makeCostModifier()

	rogue.registerStealthAura()
	rogue.registerBackstabSpell()
	rogue.registerDeadlyPoisonSpell()
	rogue.registerPoisonAuras()
	rogue.registerEviscerate()
	rogue.registerExposeArmorSpell()
	rogue.registerFanOfKnives()
	rogue.registerFeintSpell()
	rogue.registerGarrote()
	rogue.registerHemorrhageSpell()
	rogue.registerInstantPoisonSpell()
	rogue.registerWoundPoisonSpell()
	rogue.registerMutilateSpell()
	rogue.registerRupture()
	rogue.registerShivSpell()
	rogue.registerSinisterStrikeSpell()
	rogue.registerSliceAndDice()
	rogue.registerThistleTeaCD()
	rogue.registerTricksOfTheTradeSpell()
	rogue.registerAmbushSpell()
	rogue.registerEnvenom()
	rogue.registerVanishSpell()

	rogue.finishingMoveEffectApplier = rogue.makeFinishingMoveEffectApplier()
}

func (rogue *Rogue) ApplyEnergyTickMultiplier(multiplier float64) {
	rogue.EnergyTickMultiplier += multiplier
}

func (rogue *Rogue) Reset(sim *core.Simulation) {
	for _, mcd := range rogue.GetMajorCooldowns() {
		mcd.Disable()
	}
}

func (rogue *Rogue) MeleeCritMultiplier(applyLethality bool) float64 {
	primaryModifier := rogue.preyOnTheWeakMultiplier(rogue.CurrentTarget)
	var secondaryModifier float64
	if applyLethality {
		secondaryModifier += 0.06 * float64(rogue.Talents.Lethality)
	}
	return rogue.Character.MeleeCritMultiplier(primaryModifier, secondaryModifier)
}
func (rogue *Rogue) SpellCritMultiplier() float64 {
	primaryModifier := rogue.preyOnTheWeakMultiplier(rogue.CurrentTarget)
	return rogue.Character.SpellCritMultiplier(primaryModifier, 0)
}

func NewRogue(character *core.Character, options *proto.Player) *Rogue {
	rogueOptions := options.GetRogue()

	rogue := &Rogue{
		Character: *character,
		Talents:   &proto.RogueTalents{},
		Options:   rogueOptions.Options,
	}
	core.FillTalentsProto(rogue.Talents.ProtoReflect(), options.TalentsString, TalentTreeSizes)

	// Passive rogue threat reduction: https://wotlk.wowhead.com/spell=21184/rogue-passive-dnd
	rogue.PseudoStats.ThreatMultiplier *= 0.71
	rogue.PseudoStats.CanParry = true

	maxEnergy := 100.0
	if rogue.Talents.Vigor {
		maxEnergy += 10
	}
	if rogue.HasMajorGlyph(proto.RogueMajorGlyph_GlyphOfVigor) {
		maxEnergy += 10
	}
	if rogue.HasSetBonus(Arena, 4) {
		maxEnergy += 10
	}
	rogue.EnableEnergyBar(maxEnergy)
	rogue.ApplyEnergyTickMultiplier([]float64{0, 0.08, 0.16, 0.25}[rogue.Talents.Vitality])

	rogue.EnableAutoAttacks(rogue, core.AutoAttackOptions{
		MainHand:       rogue.WeaponFromMainHand(0), // Set crit multiplier later when we have targets.
		OffHand:        rogue.WeaponFromOffHand(0),  // Set crit multiplier later when we have targets.
		AutoSwingMelee: true,
	})
	rogue.applyPoisons()

	rogue.AddStatDependency(stats.Strength, stats.AttackPower, 1)
	rogue.AddStatDependency(stats.Agility, stats.AttackPower, 1)
	rogue.AddStatDependency(stats.Agility, stats.MeleeCrit, core.CritPerAgiMaxLevel[character.Class]*core.CritRatingPerCritChance)

	return rogue
}

// Apply the effects of the Cut to the Chase talent
// TODO: Put a fresh instance of SnD rather than use the original as per client
func (rogue *Rogue) ApplyCutToTheChase(sim *core.Simulation) {
	if rogue.Talents.CutToTheChase > 0 && rogue.SliceAndDiceAura.IsActive() {
		procChance := float64(rogue.Talents.CutToTheChase) * 0.2
		if sim.Proc(procChance, "Cut to the Chase") {
			rogue.SliceAndDiceAura.Duration = rogue.sliceAndDiceDurations[5]
			rogue.SliceAndDiceAura.Activate(sim)
		}
	}
}

// Deactivate Stealth if it is active. This must be added to all abilities that cause Stealth to fade.
func (rogue *Rogue) BreakStealth(sim *core.Simulation) {
	if rogue.StealthAura.IsActive() {
		rogue.StealthAura.Deactivate(sim)
		rogue.AutoAttacks.EnableAutoSwing(sim)
	}
}

// Does the rogue have a dagger equipped in the specified hand (main or offhand)?
func (rogue *Rogue) HasDagger(hand core.Hand) bool {
	if hand == core.MainHand {
		return rogue.MainHand().WeaponType == proto.WeaponType_WeaponTypeDagger
	}
	return rogue.OffHand().WeaponType == proto.WeaponType_WeaponTypeDagger
}

// Check if the rogue is considered in "stealth" for the purpose of casting abilities
func (rogue *Rogue) IsStealthed() bool {
	if rogue.StealthAura.IsActive() {
		return true
	}
	if rogue.Talents.ShadowDance && rogue.ShadowDanceAura.IsActive() {
		return true
	}
	return false
}

// Agent is a generic way to access underlying rogue on any of the agents.
type RogueAgent interface {
	GetRogue() *Rogue
}
