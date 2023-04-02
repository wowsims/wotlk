package core

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

type Encounter struct {
	Duration          time.Duration
	DurationVariation time.Duration
	Targets           []*Target
	TargetUnits       []*Unit

	ExecuteProportion_20 float64
	ExecuteProportion_25 float64
	ExecuteProportion_35 float64

	EndFightAtHealth float64
	// DamgeTaken is used to track health fights instead of duration fights.
	//  Once primary target has taken its health worth of damage, fight ends.
	DamageTaken float64
	// In health fight: set to true until we get something to base on
	DurationIsEstimate bool

	// Value to multiply by, for damage spells which are subject to the aoe cap.
	aoeCapMultiplier float64
}

func NewEncounter(options *proto.Encounter) Encounter {
	options.ExecuteProportion_25 = MaxFloat(options.ExecuteProportion_25, options.ExecuteProportion_20)
	options.ExecuteProportion_35 = MaxFloat(options.ExecuteProportion_35, options.ExecuteProportion_25)

	encounter := Encounter{
		Duration:             DurationFromSeconds(options.Duration),
		DurationVariation:    DurationFromSeconds(options.DurationVariation),
		ExecuteProportion_20: MaxFloat(options.ExecuteProportion_20, 0),
		ExecuteProportion_25: MaxFloat(options.ExecuteProportion_25, 0),
		ExecuteProportion_35: MaxFloat(options.ExecuteProportion_35, 0),
		Targets:              []*Target{},
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
		target := NewTarget(targetOptions, int32(targetIndex))
		encounter.Targets = append(encounter.Targets, target)
		encounter.TargetUnits = append(encounter.TargetUnits, &target.Unit)
	}
	if len(encounter.Targets) == 0 {
		// Add a dummy target. The only case where targets aren't specified is when
		// computing character stats, and targets won't matter there.
		target := NewTarget(&proto.Target{}, 0)
		encounter.Targets = append(encounter.Targets, target)
		encounter.TargetUnits = append(encounter.TargetUnits, &target.Unit)
	}

	if encounter.EndFightAtHealth > 0 {
		// Until we pre-sim set duration to 10m
		encounter.Duration = time.Minute * 10
		encounter.DurationIsEstimate = true
	}

	encounter.updateAOECapMultiplier()

	return encounter
}

func (encounter *Encounter) AOECapMultiplier() float64 {
	return encounter.aoeCapMultiplier
}
func (encounter *Encounter) updateAOECapMultiplier() {
	encounter.aoeCapMultiplier = MinFloat(10/float64(len(encounter.Targets)), 1)
}

func (encounter *Encounter) doneIteration(sim *Simulation) {
	for i := range encounter.Targets {
		target := encounter.Targets[i]
		target.doneIteration(sim)
	}
}

func (encounter *Encounter) GetMetricsProto() *proto.EncounterMetrics {
	metrics := &proto.EncounterMetrics{
		Targets: make([]*proto.UnitMetrics, len(encounter.Targets)),
	}

	i := 0
	for _, target := range encounter.Targets {
		metrics.Targets[i] = target.GetMetricsProto()
		i++
	}

	return metrics
}

// Target is an enemy/boss that can be the target of player attacks/spells.
type Target struct {
	Unit

	AI TargetAI
}

func NewTarget(options *proto.Target, targetIndex int32) *Target {
	unitStats := stats.Stats{}
	if options.Stats != nil {
		copy(unitStats[:], options.Stats)
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

			StatDependencyManager: stats.NewStatDependencyManager(),
		},
	}
	defaultRaidBossLevel := int32(CharacterLevel + 3)
	target.GCD = target.NewTimer()
	if target.Level == 0 {
		target.Level = defaultRaidBossLevel
	}
	if target.stats[stats.MeleeCrit] == 0 {
		// Treat any % crit buff an enemy would gain as though it was scaled with level 80 ratings
		target.stats[stats.MeleeCrit] = UnitLevelFloat64(target.Level, 5.0, 5.2, 5.4, 5.6) * CritRatingPerCritChance
	}

	if target.Level == defaultRaidBossLevel && options.SuppressDodge {
		// ICC boss Dodge Suppression. -20% dodge only.
		target.PseudoStats.DodgeReduction += 0.2
	}

	target.PseudoStats.CanBlock = true
	target.PseudoStats.CanParry = true
	target.PseudoStats.ParryHaste = options.ParryHaste
	target.PseudoStats.InFrontOfTarget = true
	target.PseudoStats.TightEnemyDamage = options.TightEnemyDamage

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
	target.SetGCDTimer(sim, 0)
	if target.AI != nil {
		target.AI.Reset(sim)
	}
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

func (target *Target) GetMetricsProto() *proto.UnitMetrics {
	metrics := target.Metrics.ToProto()
	metrics.Name = target.Label
	metrics.UnitIndex = target.UnitIndex
	metrics.Auras = target.auraTracker.GetMetricsProto()
	return metrics
}

// Holds cached values for outcome/damage calculations, for a specific attacker+defender pair.
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
	CritSuppression  float64

	DamageDealtMultiplier        float64 // attacker buff, applied in applyAttackerModifiers()
	DamageTakenMultiplier        float64 // defender debuff, applied in applyTargetModifiers()
	NatureDamageTakenMultiplier  float64
	HauntSEDamageTakenMultiplier float64
	HealingDealtMultiplier       float64
}

func NewAttackTable(attacker *Unit, defender *Unit) *AttackTable {
	table := &AttackTable{
		Attacker: attacker,
		Defender: defender,

		DamageDealtMultiplier:        1,
		DamageTakenMultiplier:        1,
		NatureDamageTakenMultiplier:  1,
		HauntSEDamageTakenMultiplier: 1,
		HealingDealtMultiplier:       1,
	}

	if defender.Type == EnemyUnit {
		// Assumes attacker (the Player) is level 80.
		table.BaseSpellMissChance = UnitLevelFloat64(defender.Level, 0.04, 0.05, 0.06, 0.17)
		table.BaseMissChance = UnitLevelFloat64(defender.Level, 0.05, 0.055, 0.06, 0.08)
		table.BaseBlockChance = 0.05
		table.BaseDodgeChance = UnitLevelFloat64(defender.Level, 0.05, 0.055, 0.06, 0.065)
		table.BaseParryChance = UnitLevelFloat64(defender.Level, 0.05, 0.055, 0.06, 0.14)
		table.BaseGlanceChance = UnitLevelFloat64(defender.Level, 0.06, 0.12, 0.18, 0.24)

		table.GlanceMultiplier = UnitLevelFloat64(defender.Level, 0.95, 0.95, 0.85, 0.75)
		table.CritSuppression = UnitLevelFloat64(defender.Level, 0, 0.01, 0.02, 0.048)
	} else {
		// Assumes defender (the Player) is level 80.
		table.BaseSpellMissChance = 0.05
		table.BaseMissChance = UnitLevelFloat64(attacker.Level, 0.05, 0.048, 0.046, 0.044)
		table.BaseBlockChance = UnitLevelFloat64(attacker.Level, 0.05, 0.048, 0.046, 0.044)
		table.BaseDodgeChance = UnitLevelFloat64(attacker.Level, 0, -0.002, -0.004, -0.006)
		table.BaseParryChance = UnitLevelFloat64(attacker.Level, 0, -0.002, -0.004, -0.006)
	}

	return table
}
