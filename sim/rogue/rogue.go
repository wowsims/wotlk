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
		func(character core.Character, options proto.Player) core.Agent {
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

type Rogue struct {
	core.Character

	Talents  proto.RogueTalents
	Options  proto.Rogue_Options
	Rotation proto.Rogue_Rotation

	priorityItems []roguePriorityItem
	rotationItems []rogueRotationItem

	sliceAndDiceDurations [6]time.Duration
	exposeArmorDurations  [6]time.Duration
	disabledMCDs          []*core.MajorCooldown

	initialArmorDebuffAura *core.Aura

	BuilderPoints    int32
	Builder          *core.Spell
	Backstab         *core.Spell
	BladeFlurry      *core.Spell
	DeadlyPoison     *core.Spell
	FanOfKnives      *core.Spell
	Hemorrhage       *core.Spell
	HungerForBlood   *core.Spell
	InstantPoison    [3]*core.Spell
	WoundPoison      [3]*core.Spell
	Mutilate         *core.Spell
	Shiv             *core.Spell
	SinisterStrike   *core.Spell
	TricksOfTheTrade *core.Spell

	Envenom      [6]*core.Spell
	Eviscerate   [6]*core.Spell
	ExposeArmor  [6]*core.Spell
	Rupture      [6]*core.Spell
	SliceAndDice [6]*core.Spell

	lastDeadlyPoisonProcMask    core.ProcMask
	deadlyPoisonProcChanceBonus float64
	instantPoisonPPMM           core.PPMManager
	deadlyPoisonDots            []*core.Dot
	ruptureDot                  *core.Dot

	AdrenalineRushAura   *core.Aura
	BladeFlurryAura      *core.Aura
	DeathmantleProcAura  *core.Aura
	VanCleefsProcAura    *core.Aura
	EnvenomAura          *core.Aura
	ExposeArmorAura      *core.Aura
	HungerForBloodAura   *core.Aura
	KillingSpreeAura     *core.Aura
	OverkillAura         *core.Aura
	SliceAndDiceAura     *core.Aura
	TricksOfTheTradeAura *core.Aura

	masterPoisonerDebuffAuras []*core.Aura
	savageCombatDebuffAuras   []*core.Aura
	woundPoisonDebuffAuras    []*core.Aura

	QuickRecoveryMetrics *core.ResourceMetrics

	CastModifier               func(*core.Simulation, *core.Spell, *core.Cast)
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
	rogue.AutoAttacks.MHEffect.OutcomeApplier = rogue.OutcomeFuncMeleeWhite(rogue.MeleeCritMultiplier(true, false))
	rogue.AutoAttacks.OHEffect.OutcomeApplier = rogue.OutcomeFuncMeleeWhite(rogue.MeleeCritMultiplier(false, false))

	if rogue.Talents.QuickRecovery > 0 {
		rogue.QuickRecoveryMetrics = rogue.NewEnergyMetrics(core.ActionID{SpellID: 31245})
	}

	rogue.CastModifier = rogue.makeCastModifier()

	rogue.registerBackstabSpell()
	rogue.registerDeadlyPoisonSpell()
	rogue.registerPoisonAuras()
	rogue.registerEviscerate()
	rogue.registerExposeArmorSpell()
	rogue.registerFanOfKnives()
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

	if rogue.Talents.Mutilate {
		rogue.registerEnvenom()
	}

	rogue.finishingMoveEffectApplier = rogue.makeFinishingMoveEffectApplier()
	rogue.DelayDPSCooldownsForArmorDebuffs()
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
	rogue.EnergyTickMultiplier *= multiplier
}

func (rogue *Rogue) Reset(sim *core.Simulation) {
	rogue.disabledMCDs = rogue.DisableAllEnabledCooldowns(core.CooldownTypeUnknown)
	rogue.initialArmorDebuffAura = rogue.CurrentTarget.GetActiveAuraWithTag(core.MajorArmorReductionTag)
	rogue.lastDeadlyPoisonProcMask = core.ProcMaskEmpty
	rogue.setPriorityItems(sim)
}

func (rogue *Rogue) MeleeCritMultiplier(isMH bool, applyLethality bool) float64 {
	primaryModifier := 1.0
	secondaryModifier := 0.0
	preyModifier := rogue.preyOnTheWeakMultiplier(rogue.CurrentTarget)
	if applyLethality {
		secondaryModifier += 0.06 * float64(rogue.Talents.Lethality)
	}
	primaryModifier *= preyModifier
	return rogue.Character.MeleeCritMultiplier(primaryModifier, secondaryModifier)
}
func (rogue *Rogue) SpellCritMultiplier() float64 {
	primaryModifier := rogue.preyOnTheWeakMultiplier(rogue.CurrentTarget)
	return rogue.Character.SpellCritMultiplier(primaryModifier, 0)
}

func NewRogue(character core.Character, options proto.Player) *Rogue {
	rogueOptions := options.GetRogue()

	rogue := &Rogue{
		Character: character,
		Talents:   *rogueOptions.Talents,
		Options:   *rogueOptions.Options,
		Rotation:  *rogueOptions.Rotation,
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
	rogue.EnableEnergyBar(maxEnergy, rogue.OnEnergyGain)
	rogue.EnergyTickMultiplier *= (1 + []float64{0, 0.08, 0.16, 0.25}[rogue.Talents.Vitality])

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
		chanceToRefresh := float64(rogue.Talents.CutToTheChase) * 0.2
		if chanceToRefresh == 1 || sim.RandomFloat("Cut to the Chase") < chanceToRefresh {
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
		stats.Strength:  112,
		stats.Agility:   206,
		stats.Stamina:   88,
		stats.Intellect: 43,
		stats.Spirit:    57,

		stats.AttackPower: 140,
		stats.MeleeCrit:   -0.3 * core.CritRatingPerCritChance,
		stats.SpellCrit:   -0.3 * core.CritRatingPerCritChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceDwarf, Class: proto.Class_ClassRogue}] = stats.Stats{
		stats.Health:    3524,
		stats.Strength:  120,
		stats.Agility:   200,
		stats.Stamina:   92,
		stats.Intellect: 38,
		stats.Spirit:    57,

		stats.AttackPower: 140,
		stats.MeleeCrit:   -0.3 * core.CritRatingPerCritChance,
		stats.SpellCrit:   -0.3 * core.CritRatingPerCritChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceGnome, Class: proto.Class_ClassRogue}] = stats.Stats{
		stats.Health:    3524,
		stats.Strength:  110,
		stats.Agility:   206,
		stats.Stamina:   88,
		stats.Intellect: 45,
		stats.Spirit:    58,

		stats.AttackPower: 140,
		stats.MeleeCrit:   -0.3 * core.CritRatingPerCritChance,
		stats.SpellCrit:   -0.3 * core.CritRatingPerCritChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceHuman, Class: proto.Class_ClassRogue}] = stats.Stats{
		stats.Health:    3524,
		stats.Strength:  115,
		stats.Agility:   204,
		stats.Stamina:   89,
		stats.Intellect: 39,
		stats.Spirit:    58,

		stats.AttackPower: 140,
		stats.MeleeCrit:   -0.3 * core.CritRatingPerCritChance,
		stats.SpellCrit:   -0.3 * core.CritRatingPerCritChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceNightElf, Class: proto.Class_ClassRogue}] = stats.Stats{
		stats.Health:    3524,
		stats.Strength:  111,
		stats.Agility:   208,
		stats.Stamina:   88,
		stats.Intellect: 39,
		stats.Spirit:    58,

		stats.AttackPower: 140,
		stats.MeleeCrit:   -0.3 * core.CritRatingPerCritChance,
		stats.SpellCrit:   -0.3 * core.CritRatingPerCritChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceOrc, Class: proto.Class_ClassRogue}] = stats.Stats{
		stats.Health:    3524,
		stats.Strength:  118,
		stats.Agility:   201,
		stats.Stamina:   91,
		stats.Intellect: 36,
		stats.Spirit:    61,

		stats.AttackPower: 140,
		stats.MeleeCrit:   -0.3 * core.CritRatingPerCritChance,
		stats.SpellCrit:   -0.3 * core.CritRatingPerCritChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceTroll, Class: proto.Class_ClassRogue}] = stats.Stats{
		stats.Health:    3524,
		stats.Strength:  116,
		stats.Agility:   206,
		stats.Stamina:   90,
		stats.Intellect: 35,
		stats.Spirit:    59,

		stats.AttackPower: 140,
		stats.MeleeCrit:   -0.3 * core.CritRatingPerCritChance,
		stats.SpellCrit:   -0.3 * core.CritRatingPerCritChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceUndead, Class: proto.Class_ClassRogue}] = stats.Stats{
		stats.Health:    3524,
		stats.Strength:  114,
		stats.Agility:   202,
		stats.Stamina:   90,
		stats.Intellect: 37,
		stats.Spirit:    63,

		stats.AttackPower: 140,
		stats.MeleeCrit:   -0.3 * core.CritRatingPerCritChance,
		stats.SpellCrit:   -0.3 * core.CritRatingPerCritChance,
	}
}

// Agent is a generic way to access underlying rogue on any of the agents.
type RogueAgent interface {
	GetRogue() *Rogue
}
