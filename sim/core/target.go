package core

import (
	"strconv"
	"time"

	"github.com/wowsims/tbc/sim/core/proto"
	"github.com/wowsims/tbc/sim/core/stats"
)

type Encounter struct {
	Duration           time.Duration
	DurationVariation  time.Duration
	executePhaseBegins time.Duration
	Targets            []*Target

	EndFightAtHealth float64
	// DamgeTaken is used to track health fights instead of duration fights.
	//  Once primary target has taken its health worth of damage, fight ends.
	DamageTaken float64
	// In health fight: set to true until we get something to base on
	DurationIsEstimate bool
}

func NewEncounter(options proto.Encounter) Encounter {
	encounter := Encounter{
		Duration:           DurationFromSeconds(options.Duration),
		DurationVariation:  DurationFromSeconds(options.DurationVariation),
		executePhaseBegins: DurationFromSeconds(options.Duration * (1 - options.ExecuteProportion)),
		Targets:            []*Target{},
	}
	// If UseHealth is set, we use the sum of targets health.
	if options.UseHealth {
		for _, t := range options.Targets {
			encounter.EndFightAtHealth += t.Stats[stats.Health]
		}
		if encounter.EndFightAtHealth == 0 {
			encounter.EndFightAtHealth = 1 // default to something so we don't instantly end without anything.
		}
	}

	for targetIndex, targetOptions := range options.Targets {
		target := NewTarget(*targetOptions, int32(targetIndex))
		encounter.Targets = append(encounter.Targets, target)
	}
	if len(encounter.Targets) == 0 {
		// Add a dummy target. The only case where targets aren't specified is when
		// computing character stats, and targets won't matter there.
		encounter.Targets = append(encounter.Targets, NewTarget(proto.Target{}, 0))
	}

	if encounter.EndFightAtHealth > 0 {
		// Until we pre-sim set duration to 10m
		encounter.Duration = time.Minute * 10
		encounter.DurationIsEstimate = true
	}

	return encounter
}

func (encounter *Encounter) doneIteration(sim *Simulation) {
	for i, _ := range encounter.Targets {
		target := encounter.Targets[i]
		target.doneIteration(sim)
	}
}

func (encounter *Encounter) GetMetricsProto(numIterations int32) *proto.EncounterMetrics {
	metrics := &proto.EncounterMetrics{
		Targets: make([]*proto.UnitMetrics, len(encounter.Targets)),
	}

	i := 0
	for _, target := range encounter.Targets {
		metrics.Targets[i] = target.GetMetricsProto(numIterations)
		i++
	}

	return metrics
}

// Target is an enemy/boss that can be the target of player attacks/spells.
type Target struct {
	Unit

	AI TargetAI
}

func NewTarget(options proto.Target, targetIndex int32) *Target {
	unitStats := stats.Stats{}
	if options.Stats != nil {
		copy(unitStats[:], options.Stats[:])
	}

	target := &Target{
		Unit: Unit{
			Type:        EnemyUnit,
			Index:       targetIndex,
			Label:       "Target " + strconv.Itoa(int(targetIndex)+1),
			Level:       options.Level,
			MobType:     options.MobType,
			auraTracker: newAuraTracker(),
			stats:       unitStats,
			PseudoStats: stats.NewPseudoStats(),
			Metrics:     NewUnitMetrics(),
		},
	}
	target.GCD = target.NewTimer()
	if target.Level == 0 {
		target.Level = 73
	}
	if target.stats[stats.MeleeCrit] == 0 {
		target.stats[stats.MeleeCrit] = UnitLevelFloat64(target.Level, 0.05, 0.052, 0.054, 0.056) * MeleeCritRatingPerCritChance
	}

	if target.Level == 73 && options.SuppressDodge {
		// Sunwell boss Dodge Suppression. -20% dodge and -5% miss chance.
		target.PseudoStats.DodgeReduction += 0.2
		target.PseudoStats.IncreasedMissChance -= 0.05
	}

	target.PseudoStats.CanBlock = true
	target.PseudoStats.CanParry = true
	target.PseudoStats.ParryHaste = options.ParryHaste
	target.PseudoStats.InFrontOfTarget = true
	if target.Level == 73 && options.CanCrush {
		target.PseudoStats.CanCrush = true
	}

	preset := GetPresetTargetWithID(options.Id)
	if preset != nil && preset.AI != nil {
		target.AI = preset.AI()
	}

	return target
}

func (target *Target) finalize() {
	target.Unit.finalize()
}

func (target *Target) init(sim *Simulation) {
	target.Unit.init(sim)
}

func (target *Target) Reset(sim *Simulation) {
	target.Unit.reset(sim, nil)
	//target.SetGCDTimer(sim, 0)
}

func (target *Target) Advance(sim *Simulation, elapsedTime time.Duration) {
	target.Unit.advance(sim, elapsedTime)
}

func (target *Target) doneIteration(sim *Simulation) {
	target.Unit.doneIteration(sim)
}

func (target *Target) NextTarget() *Target {
	nextIndex := target.Index + 1
	if nextIndex >= target.Env.GetNumTargets() {
		nextIndex = 0
	}
	return target.Env.GetTarget(nextIndex)
}

func (target *Target) GetMetricsProto(numIterations int32) *proto.UnitMetrics {
	metrics := target.Metrics.ToProto(numIterations)
	metrics.Name = target.Label
	metrics.Auras = target.auraTracker.GetMetricsProto(numIterations)
	return metrics
}

// Holds cached values for outcome/damage calculations, for a specific attacker+defender pair.
//
// These are updated dynamically when attacker or defender stats change.
type AttackTable struct {
	Attacker *Unit
	Defender *Unit

	BaseMissChance      float64
	BaseSpellMissChance float64
	BaseBlockChance     float64
	BaseDodgeChance     float64
	BaseParryChance     float64
	BaseGlanceChance    float64

	GlanceMultiplier float64
	HitSuppression   float64
	CritSuppression  float64

	PartialResistArcaneRollThreshold00 float64
	PartialResistArcaneRollThreshold25 float64
	PartialResistArcaneRollThreshold50 float64
	PartialResistHolyRollThreshold00   float64
	PartialResistHolyRollThreshold25   float64
	PartialResistHolyRollThreshold50   float64
	PartialResistFireRollThreshold00   float64
	PartialResistFireRollThreshold25   float64
	PartialResistFireRollThreshold50   float64
	PartialResistFrostRollThreshold00  float64
	PartialResistFrostRollThreshold25  float64
	PartialResistFrostRollThreshold50  float64
	PartialResistNatureRollThreshold00 float64
	PartialResistNatureRollThreshold25 float64
	PartialResistNatureRollThreshold50 float64
	PartialResistShadowRollThreshold00 float64
	PartialResistShadowRollThreshold25 float64
	PartialResistShadowRollThreshold50 float64

	BinaryArcaneHitChance float64
	BinaryHolyHitChance   float64
	BinaryFireHitChance   float64
	BinaryFrostHitChance  float64
	BinaryNatureHitChance float64
	BinaryShadowHitChance float64

	ArmorDamageReduction float64
}

func NewAttackTable(attacker *Unit, defender *Unit) *AttackTable {
	table := &AttackTable{
		Attacker: attacker,
		Defender: defender,
	}

	if defender.Type == EnemyUnit {
		// Assumes attacker (the Player) is level 70.
		table.BaseSpellMissChance = UnitLevelFloat64(defender.Level, 0.04, 0.05, 0.06, 0.17)
		table.BaseMissChance = UnitLevelFloat64(defender.Level, 0.05, 0.055, 0.06, 0.08)
		table.BaseBlockChance = 0.05
		table.BaseDodgeChance = UnitLevelFloat64(defender.Level, 0.05, 0.055, 0.06, 0.065)
		table.BaseParryChance = UnitLevelFloat64(defender.Level, 0.05, 0.055, 0.06, 0.14)
		table.BaseGlanceChance = UnitLevelFloat64(defender.Level, 0.06, 0.12, 0.18, 0.24)

		table.GlanceMultiplier = UnitLevelFloat64(defender.Level, 0.95, 0.95, 0.85, 0.75)
		table.HitSuppression = UnitLevelFloat64(defender.Level, 0, 0, 0, 0.01)
		table.CritSuppression = UnitLevelFloat64(defender.Level, 0, 0.01, 0.02, 0.048)
	} else {
		// Assumes defender (the Player) is level 70.
		table.BaseSpellMissChance = 0.05
		table.BaseMissChance = UnitLevelFloat64(attacker.Level, 0.05, 0.048, 0.046, 0.044)
		table.BaseBlockChance = UnitLevelFloat64(attacker.Level, 0.05, 0.048, 0.046, 0.044)
		table.BaseDodgeChance = UnitLevelFloat64(attacker.Level, 0, -0.002, -0.004, -0.006)
		table.BaseParryChance = UnitLevelFloat64(attacker.Level, 0.05, 0.048, 0.046, 0.044)
	}

	table.UpdateArmorDamageReduction()
	table.UpdatePartialResists()

	return table
}
