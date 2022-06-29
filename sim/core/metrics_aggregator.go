package core

import (
	"math"
	"time"

	"github.com/wowsims/tbc/sim/core/proto"
)

type ResourceKey struct {
	ActionID ActionID
	Type     proto.ResourceType
}

type DistributionMetrics struct {
	// Values for the current iteration. These are cleared after each iteration.
	Total float64

	// Aggregate values. These are updated after each iteration.
	sum        float64
	sumSquared float64
	max        float64
	hist       map[int32]int32 // rounded DPS to count
}

func (distMetrics *DistributionMetrics) reset() {
	distMetrics.Total = 0
}

// This should be called when a Sim iteration is complete.
func (distMetrics *DistributionMetrics) doneIteration(encounterDurationSeconds float64) {
	dps := distMetrics.Total / encounterDurationSeconds

	distMetrics.sum += dps
	distMetrics.sumSquared += dps * dps
	distMetrics.max = MaxFloat(distMetrics.max, dps)

	dpsRounded := int32(math.Round(dps/10) * 10)
	distMetrics.hist[dpsRounded]++
}

func (distMetrics *DistributionMetrics) ToProto(numIterations int32) *proto.DistributionMetrics {
	dpsAvg := distMetrics.sum / float64(numIterations)

	return &proto.DistributionMetrics{
		Avg:   dpsAvg,
		Stdev: math.Sqrt((distMetrics.sumSquared / float64(numIterations)) - (dpsAvg * dpsAvg)),
		Max:   distMetrics.max,
		Hist:  distMetrics.hist,
	}
}

func NewDistributionMetrics() DistributionMetrics {
	return DistributionMetrics{
		hist: make(map[int32]int32),
	}
}

type UnitMetrics struct {
	dps    DistributionMetrics
	threat DistributionMetrics
	dtps   DistributionMetrics

	CharacterIterationMetrics

	// Aggregate values. These are updated after each iteration.
	numItersDead int32
	oomTimeSum   float64
	actions      map[ActionID]*ActionMetrics
	resources    []*ResourceMetrics
}

// Metrics for the current iteration, for 1 agent. Keep this as a separate
// struct so its easy to clear.
type CharacterIterationMetrics struct {
	Died    bool // Whether this unit died in the current iteration.
	WentOOM bool // Whether the agent has hit OOM at least once in this iteration.

	ManaSpent       float64
	ManaGained      float64
	BonusManaGained float64 // Only includes amount from mana pots / runes / innervates.

	OOMTime time.Duration // time spent not casting and waiting for regen.
}

type ActionMetrics struct {
	IsMelee bool // True if melee action, false if spell action.

	// Metrics for this action, for each possible target.
	Targets []TargetedActionMetrics
}

func (actionMetrics *ActionMetrics) ToProto(actionID ActionID) *proto.ActionMetrics {
	var targetMetrics []*proto.TargetedActionMetrics
	for _, tam := range actionMetrics.Targets {
		targetMetrics = append(targetMetrics, tam.ToProto())
	}

	return &proto.ActionMetrics{
		Id:      actionID.ToProto(),
		IsMelee: actionMetrics.IsMelee,
		Targets: targetMetrics,
	}
}

type TargetedActionMetrics struct {
	UnitIndex int32

	Casts   int32
	Hits    int32
	Crits   int32
	Crushes int32
	Misses  int32
	Dodges  int32
	Parries int32
	Blocks  int32
	Glances int32

	Damage float64
	Threat float64
}

func (tam *TargetedActionMetrics) ToProto() *proto.TargetedActionMetrics {
	return &proto.TargetedActionMetrics{
		UnitIndex: tam.UnitIndex,

		Casts:   tam.Casts,
		Hits:    tam.Hits,
		Crits:   tam.Crits,
		Crushes: tam.Crushes,
		Misses:  tam.Misses,
		Dodges:  tam.Dodges,
		Parries: tam.Parries,
		Blocks:  tam.Blocks,
		Glances: tam.Glances,
		Damage:  tam.Damage,
		Threat:  tam.Threat,
	}
}

func NewUnitMetrics() UnitMetrics {
	return UnitMetrics{
		dps:     NewDistributionMetrics(),
		threat:  NewDistributionMetrics(),
		dtps:    NewDistributionMetrics(),
		actions: make(map[ActionID]*ActionMetrics),
	}
}

type ResourceMetrics struct {
	ActionID ActionID
	Type     proto.ResourceType

	Events     int32
	Gain       float64
	ActualGain float64

	EventsFromPreviousIterations     int32
	ActualGainFromPreviousIterations float64
}

func (resourceMetrics *ResourceMetrics) ToProto() *proto.ResourceMetrics {
	return &proto.ResourceMetrics{
		Id:   resourceMetrics.ActionID.ToProto(),
		Type: resourceMetrics.Type,

		Events:     resourceMetrics.Events,
		Gain:       resourceMetrics.Gain,
		ActualGain: resourceMetrics.ActualGain,
	}
}

func (resourceMetrics *ResourceMetrics) reset() {
	resourceMetrics.EventsFromPreviousIterations = resourceMetrics.Events
	resourceMetrics.ActualGainFromPreviousIterations = resourceMetrics.ActualGain
}
func (resourceMetrics *ResourceMetrics) EventsForCurrentIteration() int32 {
	return resourceMetrics.Events - resourceMetrics.EventsFromPreviousIterations
}
func (resourceMetrics *ResourceMetrics) ActualGainForCurrentIteration() float64 {
	return resourceMetrics.ActualGain - resourceMetrics.ActualGainFromPreviousIterations
}

func (resourceMetrics *ResourceMetrics) AddEvent(gain float64, actualGain float64) {
	resourceMetrics.Events++
	resourceMetrics.Gain += gain
	resourceMetrics.ActualGain += actualGain
}

func (unitMetrics *UnitMetrics) NewResourceMetrics(actionID ActionID, resourceType proto.ResourceType) *ResourceMetrics {
	newMetrics := &ResourceMetrics{
		ActionID: actionID,
		Type:     resourceType,
	}
	unitMetrics.resources = append(unitMetrics.resources, newMetrics)
	return newMetrics
}

// Convenience helpers for NewResourceMetrics.
func (unit *Unit) NewHealthMetrics(actionID ActionID) *ResourceMetrics {
	return unit.Metrics.NewResourceMetrics(actionID, proto.ResourceType_ResourceTypeHealth)
}
func (unit *Unit) NewManaMetrics(actionID ActionID) *ResourceMetrics {
	return unit.Metrics.NewResourceMetrics(actionID, proto.ResourceType_ResourceTypeMana)
}
func (unit *Unit) NewRageMetrics(actionID ActionID) *ResourceMetrics {
	return unit.Metrics.NewResourceMetrics(actionID, proto.ResourceType_ResourceTypeRage)
}
func (unit *Unit) NewEnergyMetrics(actionID ActionID) *ResourceMetrics {
	return unit.Metrics.NewResourceMetrics(actionID, proto.ResourceType_ResourceTypeEnergy)
}
func (unit *Unit) NewComboPointMetrics(actionID ActionID) *ResourceMetrics {
	return unit.Metrics.NewResourceMetrics(actionID, proto.ResourceType_ResourceTypeComboPoints)
}
func (unit *Unit) NewFocusMetrics(actionID ActionID) *ResourceMetrics {
	return unit.Metrics.NewResourceMetrics(actionID, proto.ResourceType_ResourceTypeFocus)
}

// Adds the results of a spell to the character metrics.
func (unitMetrics *UnitMetrics) addSpell(spell *Spell) {
	actionMetrics, ok := unitMetrics.actions[spell.ActionID]

	if !ok {
		actionMetrics = &ActionMetrics{IsMelee: spell.Flags.Matches(SpellFlagMeleeMetrics)}
		unitMetrics.actions[spell.ActionID] = actionMetrics
	}

	if len(actionMetrics.Targets) == 0 {
		actionMetrics.Targets = make([]TargetedActionMetrics, len(spell.SpellMetrics))
		for i, _ := range actionMetrics.Targets {
			tam := &actionMetrics.Targets[i]
			tam.UnitIndex = spell.Unit.AttackTables[i].Defender.Index
		}
	}

	for i, spellTargetMetrics := range spell.SpellMetrics {
		tam := &actionMetrics.Targets[i]
		tam.Casts += spellTargetMetrics.Casts
		tam.Misses += spellTargetMetrics.Misses
		tam.Hits += spellTargetMetrics.Hits
		tam.Crits += spellTargetMetrics.Crits
		tam.Crushes += spellTargetMetrics.Crushes
		tam.Dodges += spellTargetMetrics.Dodges
		tam.Parries += spellTargetMetrics.Parries
		tam.Blocks += spellTargetMetrics.Blocks
		tam.Glances += spellTargetMetrics.Glances
		tam.Damage += spellTargetMetrics.TotalDamage
		tam.Threat += spellTargetMetrics.TotalThreat
		unitMetrics.dps.Total += spellTargetMetrics.TotalDamage
		unitMetrics.threat.Total += spellTargetMetrics.TotalThreat

		target := spell.Unit.AttackTables[i].Defender
		target.Metrics.dtps.Total += spellTargetMetrics.TotalDamage
	}
}

// This should be called at the end of each iteration, to include metrics from Pets in
// those of their owner.
// Assumes that doneIteration() has already been called on the pet metrics.
func (unitMetrics *UnitMetrics) AddFinalPetMetrics(petMetrics *UnitMetrics) {
	unitMetrics.dps.Total += petMetrics.dps.Total
}

func (unitMetrics *UnitMetrics) MarkOOM(unit *Unit, dur time.Duration) {
	unitMetrics.CharacterIterationMetrics.OOMTime += dur
	unitMetrics.CharacterIterationMetrics.WentOOM = true
}

func (unitMetrics *UnitMetrics) reset() {
	unitMetrics.dps.reset()
	unitMetrics.threat.reset()
	unitMetrics.dtps.reset()
	unitMetrics.CharacterIterationMetrics = CharacterIterationMetrics{}

	for _, resourceMetrics := range unitMetrics.resources {
		resourceMetrics.reset()
	}
}

// This should be called when a Sim iteration is complete.
func (unitMetrics *UnitMetrics) doneIteration(encounterDurationSeconds float64) {
	unitMetrics.dps.doneIteration(encounterDurationSeconds)
	unitMetrics.threat.doneIteration(encounterDurationSeconds)
	unitMetrics.dtps.doneIteration(encounterDurationSeconds)
	unitMetrics.oomTimeSum += float64(unitMetrics.OOMTime.Seconds())
	if unitMetrics.Died {
		unitMetrics.numItersDead++
	}
}

func (unitMetrics *UnitMetrics) ToProto(numIterations int32) *proto.UnitMetrics {
	protoMetrics := &proto.UnitMetrics{
		Dps:           unitMetrics.dps.ToProto(numIterations),
		Threat:        unitMetrics.threat.ToProto(numIterations),
		Dtps:          unitMetrics.dtps.ToProto(numIterations),
		SecondsOomAvg: unitMetrics.oomTimeSum / float64(numIterations),
		ChanceOfDeath: float64(unitMetrics.numItersDead) / float64(numIterations),
	}

	for actionID, action := range unitMetrics.actions {
		protoMetrics.Actions = append(protoMetrics.Actions, action.ToProto(actionID))
	}
	for _, resource := range unitMetrics.resources {
		if resource.Events > 0 {
			protoMetrics.Resources = append(protoMetrics.Resources, resource.ToProto())
		}
	}

	return protoMetrics
}

type AuraMetrics struct {
	ID ActionID

	// Metrics for the current iteration.
	Uptime time.Duration

	// Aggregate values. These are updated after each iteration.
	uptimeSum        time.Duration
	uptimeSumSquared time.Duration
}

func (auraMetrics *AuraMetrics) reset() {
	auraMetrics.Uptime = 0
}

// This should be called when a Sim iteration is complete.
func (auraMetrics *AuraMetrics) doneIteration() {
	auraMetrics.uptimeSum += auraMetrics.Uptime
	auraMetrics.uptimeSumSquared += auraMetrics.Uptime * auraMetrics.Uptime
}

func (auraMetrics *AuraMetrics) ToProto(numIterations int32) *proto.AuraMetrics {
	uptimeAvg := auraMetrics.uptimeSum.Seconds() / float64(numIterations)

	return &proto.AuraMetrics{
		Id: auraMetrics.ID.ToProto(),

		UptimeSecondsAvg:   uptimeAvg,
		UptimeSecondsStdev: math.Sqrt((auraMetrics.uptimeSumSquared.Seconds() / float64(numIterations)) - (uptimeAvg * uptimeAvg)),
	}
}

// Calculates DPS for an action.
func GetActionDPS(playerMetrics proto.UnitMetrics, iterations int32, duration time.Duration, actionID ActionID, ignoreTag bool) float64 {
	totalDPS := 0.0
	for _, action := range playerMetrics.Actions {
		metricsActionID := ProtoToActionID(*action.Id)
		if actionID.SameAction(metricsActionID) || (ignoreTag && actionID.SameActionIgnoreTag(metricsActionID)) {
			for _, tam := range action.Targets {
				totalDPS += tam.Damage / float64(iterations) / duration.Seconds()
			}
		}
	}
	return totalDPS
}

// Calculates average cast damage for an action.
func GetActionAvgCast(playerMetrics proto.UnitMetrics, actionID ActionID) float64 {
	for _, action := range playerMetrics.Actions {
		if actionID.SameAction(ProtoToActionID(*action.Id)) {
			casts := int32(0)
			damage := 0.0
			for _, tam := range action.Targets {
				casts += tam.Casts
				damage += tam.Damage
			}
			if casts == 0 {
				return 0
			} else {
				return damage / float64(casts)
			}
		}
	}
	return 0
}
