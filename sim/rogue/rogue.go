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
		func(character core.Character, options *proto.Player) core.Agent {
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
	SpellFlagBuilder  = core.SpellFlagAgentReserved2
	SpellFlagFinisher = core.SpellFlagAgentReserved3
)

const RogueBleedTag = "RogueBleed"

type Rogue struct {
	core.Character

	Talents  *proto.RogueTalents
	Options  *proto.Rogue_Options
	Rotation *proto.Rogue_Rotation

	priorityItems      []roguePriorityItem
	rotationItems      []rogueRotationItem
	assassinationPrios []assassinationPrio
	subtletyPrios      []subtletyPrio
	bleedCategory      *core.ExclusiveCategory

	sliceAndDiceDurations [6]time.Duration
	exposeArmorDurations  [6]time.Duration

	allMCDsDisabled bool

	maxEnergy float64

	BuilderPoints     int32
	Builder           *core.Spell
	Backstab          *core.Spell
	BladeFlurry       *core.Spell
	DeadlyPoison      *core.Spell
	FanOfKnives       *core.Spell
	Feint             *core.Spell
	Garrote           *core.Spell
	Ambush            *core.Spell
	Hemorrhage        *core.Spell
	GhostlyStrike     *core.Spell
	HungerForBlood    *core.Spell
	InstantPoison     [3]*core.Spell
	WoundPoison       [3]*core.Spell
	Mutilate          *core.Spell
	Shiv              *core.Spell
	SinisterStrike    *core.Spell
	TricksOfTheTrade  *core.Spell
	Shadowstep        *core.Spell
	Preparation       *core.Spell
	Premeditation     *core.Spell
	ShadowDance       *core.Spell
	ColdBlood         *core.Spell
	MasterOfSubtlety  *core.Spell
	Overkill          *core.Spell
	HonorAmongThieves *core.Spell

	Envenom      [6]*core.Spell
	Eviscerate   [6]*core.Spell
	ExposeArmor  [6]*core.Spell
	Rupture      [6]*core.Spell
	SliceAndDice [6]*core.Spell

	lastDeadlyPoisonProcMask    core.ProcMask
	deadlyPoisonProcChanceBonus float64
	deadlyPoisonDots            []*core.Dot
	garroteDot                  *core.Dot
	instantPoisonPPMM           core.PPMManager
	ruptureDot                  *core.Dot
	woundPoisonPPMM             core.PPMManager
	HonorAmongThievesDot        *core.Dot

	AdrenalineRushAura   *core.Aura
	BladeFlurryAura      *core.Aura
	EnvenomAura          *core.Aura
	ExposeArmorAura      *core.Aura
	HungerForBloodAura   *core.Aura
	KillingSpreeAura     *core.Aura
	OverkillAura         *core.Aura
	SliceAndDiceAura     *core.Aura
	TricksOfTheTradeAura *core.Aura
	MasterOfSubtletyAura *core.Aura
	ShadowstepAura       *core.Aura
	ShadowDanceAura      *core.Aura
	DirtyDeedsAura       *core.Aura

	masterPoisonerDebuffAuras []*core.Aura
	savageCombatDebuffAuras   []*core.Aura
	woundPoisonDebuffAuras    []*core.Aura

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

func (rogue *Rogue) AddRaidBuffs(raidBuffs *proto.RaidBuffs)    {}
func (rogue *Rogue) AddPartyBuffs(partyBuffs *proto.PartyBuffs) {}

func (rogue *Rogue) finisherFlags() core.SpellFlag {
	flags := SpellFlagFinisher
	if rogue.Talents.SurpriseAttacks {
		flags |= core.SpellFlagCannotBeDodged
	}
	return flags
}

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
	rogue.AutoAttacks.MHConfig.CritMultiplier = rogue.MeleeCritMultiplier(false)
	rogue.AutoAttacks.OHConfig.CritMultiplier = rogue.MeleeCritMultiplier(false)

	if rogue.Talents.QuickRecovery > 0 {
		rogue.QuickRecoveryMetrics = rogue.NewEnergyMetrics(core.ActionID{SpellID: 31245})
	}

	rogue.costModifier = rogue.makeCostModifier()

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

	if rogue.Talents.MasterPoisoner > 0 || rogue.Talents.CutToTheChase > 0 || rogue.Talents.Mutilate {
		rogue.registerEnvenom()
	}

	rogue.finishingMoveEffectApplier = rogue.makeFinishingMoveEffectApplier()
	rogue.DelayDPSCooldownsForArmorDebuffs(time.Second * 14)
}

func (rogue *Rogue) getExpectedEnergyPerSecond() float64 {
	const finishersPerSecond = 1.0 / 6
	const averageComboPointsSpendOnFinisher = 4.0
	bonusEnergyPerSecond := float64(rogue.Talents.CombatPotency) * 3 * 0.2 * 1.0 / (rogue.AutoAttacks.OH.SwingSpeed / 1.4)
	bonusEnergyPerSecond += float64(rogue.Talents.FocusedAttacks)
	bonusEnergyPerSecond += float64(rogue.Talents.RelentlessStrikes) * 0.04 * 25 * finishersPerSecond * averageComboPointsSpendOnFinisher
	return (core.EnergyPerTick*rogue.EnergyTickMultiplier)/core.EnergyTickDuration.Seconds() + bonusEnergyPerSecond
}

func (rogue *Rogue) ApplyEnergyTickMultiplier(multiplier float64) {
	rogue.EnergyTickMultiplier += multiplier
}

func (rogue *Rogue) getExpectedComboPointPerSecond() float64 {
	const criticalPerSecond = 1
	honorAmongThievesChance := []float64{0, 0.33, 0.66, 1.0}[rogue.Talents.HonorAmongThieves]
	return criticalPerSecond * honorAmongThievesChance
}

func (rogue *Rogue) Reset(sim *core.Simulation) {
	for _, mcd := range rogue.GetMajorCooldowns() {
		mcd.Disable()
	}
	rogue.allMCDsDisabled = true
	rogue.lastDeadlyPoisonProcMask = core.ProcMaskEmpty
	if rogue.HonorAmongThieves != nil {
		rogue.HonorAmongThievesDot.Deactivate(sim)
		rogue.HonorAmongThievesDot.NumberOfTicks = int32(sim.Duration + sim.DurationVariation)
		rogue.HonorAmongThievesDot.RecomputeAuraDuration()
		rogue.HonorAmongThievesDot.Activate(sim)
	}
	// Vanish triggered effects (Overkill and Master of Subtlety) prepull activation
	if rogue.OverkillAura != nil && rogue.Options.StartingOverkillDuration > 0 {
		rogue.OverkillAura.Activate(sim)
		rogue.OverkillAura.UpdateExpires(sim.CurrentTime + time.Second*time.Duration(rogue.Options.StartingOverkillDuration))
	}
	if rogue.MasterOfSubtletyAura != nil && rogue.Options.StartingOverkillDuration > 0 {
		rogue.MasterOfSubtletyAura.Activate(sim)
		rogue.MasterOfSubtletyAura.UpdateExpires(sim.CurrentTime + time.Second*time.Duration(rogue.Options.StartingOverkillDuration))
	}
	rogue.setPriorityItems(sim)
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

func NewRogue(character core.Character, options *proto.Player) *Rogue {
	rogueOptions := options.GetRogue()

	rogue := &Rogue{
		Character: character,
		Talents:   rogueOptions.Talents,
		Options:   rogueOptions.Options,
		Rotation:  rogueOptions.Rotation,
	}

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
	if rogue.HasSetBonus(ItemSetGladiatorsVestments, 4) {
		maxEnergy += 10
	}
	rogue.maxEnergy = maxEnergy
	rogue.EnableEnergyBar(maxEnergy, rogue.OnEnergyGain)
	rogue.ApplyEnergyTickMultiplier([]float64{0, 0.08, 0.16, 0.25}[rogue.Talents.Vitality])

	rogue.EnableAutoAttacks(rogue, core.AutoAttackOptions{
		MainHand:       rogue.WeaponFromMainHand(0), // Set crit multiplier later when we have targets.
		OffHand:        rogue.WeaponFromOffHand(0),  // Set crit multiplier later when we have targets.
		AutoSwingMelee: true,
	})
	rogue.applyPoisons()

	rogue.AddStatDependency(stats.Strength, stats.AttackPower, 1)
	rogue.AddStatDependency(stats.Agility, stats.AttackPower, 1)
	rogue.AddStatDependency(stats.Agility, stats.MeleeCrit, core.CritRatingPerCritChance/83.15)

	return rogue
}

func (rogue *Rogue) ApplyCutToTheChase(sim *core.Simulation) {
	if rogue.Talents.CutToTheChase > 0 && rogue.SliceAndDiceAura.IsActive() {
		procChance := float64(rogue.Talents.CutToTheChase) * 0.2
		if sim.Proc(procChance, "Cut to the Chase") {
			rogue.SliceAndDiceAura.Duration = rogue.sliceAndDiceDurations[5]
			rogue.SliceAndDiceAura.Activate(sim)
		}
	}
}

func (rogue *Rogue) CanMutilate() bool {
	return rogue.Talents.Mutilate &&
		rogue.HasMHWeapon() && rogue.HasOHWeapon() &&
		rogue.GetMHWeapon().WeaponType == proto.WeaponType_WeaponTypeDagger &&
		rogue.GetOHWeapon().WeaponType == proto.WeaponType_WeaponTypeDagger
}

func init() {
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceBloodElf, Class: proto.Class_ClassRogue}] = stats.Stats{
		stats.Health:    3524,
		stats.Strength:  110,
		stats.Agility:   191,
		stats.Stamina:   105,
		stats.Intellect: 46,
		stats.Spirit:    65,

		stats.AttackPower: 140,
		stats.MeleeCrit:   -0.3 * core.CritRatingPerCritChance,
		stats.SpellCrit:   -0.3 * core.CritRatingPerCritChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceDwarf, Class: proto.Class_ClassRogue}] = stats.Stats{
		stats.Health:    3524,
		stats.Strength:  118,
		stats.Agility:   185,
		stats.Stamina:   106,
		stats.Intellect: 42,
		stats.Spirit:    66,

		stats.AttackPower: 140,
		stats.MeleeCrit:   -0.3 * core.CritRatingPerCritChance,
		stats.SpellCrit:   -0.3 * core.CritRatingPerCritChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceGnome, Class: proto.Class_ClassRogue}] = stats.Stats{
		stats.Health:    3524,
		stats.Strength:  108,
		stats.Agility:   191,
		stats.Stamina:   105,
		stats.Intellect: 48,
		stats.Spirit:    67,

		stats.AttackPower: 140,
		stats.MeleeCrit:   -0.3 * core.CritRatingPerCritChance,
		stats.SpellCrit:   -0.3 * core.CritRatingPerCritChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceHuman, Class: proto.Class_ClassRogue}] = stats.Stats{
		stats.Health:    3524,
		stats.Strength:  113,
		stats.Agility:   189,
		stats.Stamina:   107,
		stats.Intellect: 43,
		stats.Spirit:    69,

		stats.AttackPower: 140,
		stats.MeleeCrit:   -0.3 * core.CritRatingPerCritChance,
		stats.SpellCrit:   -0.3 * core.CritRatingPerCritChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceNightElf, Class: proto.Class_ClassRogue}] = stats.Stats{
		stats.Health:    3524,
		stats.Strength:  109,
		stats.Agility:   193,
		stats.Stamina:   107,
		stats.Intellect: 43,
		stats.Spirit:    67,

		stats.AttackPower: 140,
		stats.MeleeCrit:   -0.3 * core.CritRatingPerCritChance,
		stats.SpellCrit:   -0.3 * core.CritRatingPerCritChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceOrc, Class: proto.Class_ClassRogue}] = stats.Stats{
		stats.Health:    3524,
		stats.Strength:  116,
		stats.Agility:   186,
		stats.Stamina:   106,
		stats.Intellect: 40,
		stats.Spirit:    69,

		stats.AttackPower: 140,
		stats.MeleeCrit:   -0.3 * core.CritRatingPerCritChance,
		stats.SpellCrit:   -0.3 * core.CritRatingPerCritChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceTroll, Class: proto.Class_ClassRogue}] = stats.Stats{
		stats.Health:    3524,
		stats.Strength:  114,
		stats.Agility:   191,
		stats.Stamina:   105,
		stats.Intellect: 39,
		stats.Spirit:    68,

		stats.AttackPower: 140,
		stats.MeleeCrit:   -0.3 * core.CritRatingPerCritChance,
		stats.SpellCrit:   -0.3 * core.CritRatingPerCritChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceUndead, Class: proto.Class_ClassRogue}] = stats.Stats{
		stats.Health:    3524,
		stats.Strength:  112,
		stats.Agility:   187,
		stats.Stamina:   105,
		stats.Intellect: 41,
		stats.Spirit:    72,

		stats.AttackPower: 140,
		stats.MeleeCrit:   -0.3 * core.CritRatingPerCritChance,
		stats.SpellCrit:   -0.3 * core.CritRatingPerCritChance,
	}
}

// Agent is a generic way to access underlying rogue on any of the agents.
type RogueAgent interface {
	GetRogue() *Rogue
}
