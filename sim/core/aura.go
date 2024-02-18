package core

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"golang.org/x/exp/constraints"

	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

const NeverExpires = time.Duration(math.MaxInt64)

type OnInit func(aura *Aura, sim *Simulation)
type OnReset func(aura *Aura, sim *Simulation)
type OnDoneIteration func(aura *Aura, sim *Simulation)
type OnGain func(aura *Aura, sim *Simulation)
type OnExpire func(aura *Aura, sim *Simulation)
type OnStacksChange func(aura *Aura, sim *Simulation, oldStacks int32, newStacks int32)
type OnStatsChange func(aura *Aura, sim *Simulation, oldStats stats.Stats, newStats stats.Stats)

// Callback for after a spell hits the target and after damage is calculated. Use it for proc effects
// or anything that comes from the final result of the spell.
type OnSpellHit func(aura *Aura, sim *Simulation, spell *Spell, result *SpellResult)

// OnPeriodicDamage is called when dots tick, after damage is calculated. Use it for proc effects
// or anything that comes from the final result of a tick.
type OnPeriodicDamage func(aura *Aura, sim *Simulation, spell *Spell, result *SpellResult)

const Inactive = -1

// Aura lifecycle:
//
// myAura := unit.RegisterAura(myAuraConfig)
// myAura.Activate(sim)
// myAura.SetStacks(sim, 3)
// myAura.Refresh(sim)
// myAura.Deactivate(sim)
type Aura struct {
	// String label for this Aura. Guaranteed to be unique among the Auras for a single Unit.
	Label string

	// For easily grouping auras.
	Tag string

	ActionID        ActionID // If set, metrics will be tracked for this aura.
	ActionIDForProc ActionID // If set, indicates that this aura is a trigger aura for the specified proc.

	Icd *Cooldown // The internal cooldown if any

	Duration time.Duration // Duration of aura, upon being applied.

	startTime time.Duration // Time at which the aura was applied.
	expires   time.Duration // Time at which aura will be removed.

	// The unit this aura is attached to.
	Unit *Unit

	active                     bool
	activeIndex                int32 // Position of this aura's index in the activeAuras array.
	onCastCompleteIndex        int32 // Position of this aura's index in the onCastCompleteAuras array.
	onSpellHitDealtIndex       int32 // Position of this aura's index in the onSpellHitAuras array.
	onSpellHitTakenIndex       int32 // Position of this aura's index in the onSpellHitAuras array.
	onPeriodicDamageDealtIndex int32 // Position of this aura's index in the onPeriodicDamageAuras array.
	onPeriodicDamageTakenIndex int32 // Position of this aura's index in the onPeriodicDamageAuras array.
	onHealDealtIndex           int32 // Position of this aura's index in the onHealAuras array.
	onHealTakenIndex           int32 // Position of this aura's index in the onHealAuras array.
	onPeriodicHealDealtIndex   int32 // Position of this aura's index in the onPeriodicHealAuras array.
	onPeriodicHealTakenIndex   int32 // Position of this aura's index in the onPeriodicHealAuras array.

	// The number of stacks, or charges, of this aura. If this aura doesn't care
	// about charges, is just 0.
	stacks    int32
	MaxStacks int32

	ExclusiveEffects []*ExclusiveEffect

	// Lifecycle callbacks.
	OnInit          OnInit
	OnReset         OnReset
	OnDoneIteration OnDoneIteration
	OnGain          OnGain
	OnExpire        OnExpire
	OnStacksChange  OnStacksChange // Invoked when the number of stacks of this aura changes.
	OnStatsChange   OnStatsChange  // Invoked when the stats of this aura owner changes.

	OnCastComplete        OnCastComplete   // Invoked when a spell cast completes casting, before results are calculated.
	OnSpellHitDealt       OnSpellHit       // Invoked when a spell hits and this unit is the caster.
	OnSpellHitTaken       OnSpellHit       // Invoked when a spell hits and this unit is the target.
	OnPeriodicDamageDealt OnPeriodicDamage // Invoked when a dot tick occurs and this unit is the caster.
	OnPeriodicDamageTaken OnPeriodicDamage // Invoked when a dot tick occurs and this unit is the target.
	OnHealDealt           OnSpellHit       // Invoked when a heal hits and this unit is the caster.
	OnHealTaken           OnSpellHit       // Invoked when a heal hits and this unit is the target.
	OnPeriodicHealDealt   OnPeriodicDamage // Invoked when a hot tick occurs and this unit is the caster.
	OnPeriodicHealTaken   OnPeriodicDamage // Invoked when a hot tick occurs and this unit is the target.

	// If non-default, stat bonuses fron the OnGain callback of this aura will be
	// included in Character Stats in the UI.
	BuildPhase CharacterBuildPhase

	// Metrics for this aura.
	metrics AuraMetrics

	initialized bool
}

func (aura *Aura) init(sim *Simulation) {
	if aura.initialized {
		return
	}
	aura.initialized = true

	if aura.OnInit != nil {
		aura.OnInit(aura, sim)
	}
}

func (aura *Aura) reset(sim *Simulation) {
	aura.init(sim)

	if aura.IsActive() {
		panic("Active aura during reset: " + aura.Label)
	}
	if aura.stacks != 0 {
		panic("Aura nonzero stacks during reset: " + aura.Label)
	}
	aura.metrics.reset()

	if aura.OnReset != nil {
		aura.OnReset(aura, sim)
	}
}

func (aura *Aura) doneIteration(sim *Simulation) {
	if aura.IsActive() {
		panic("Active aura during doneIter: " + aura.Label)
	}
	if aura.stacks != 0 {
		panic("Aura nonzero stacks during doneIter: " + aura.Label)
	}

	aura.startTime = 0
	aura.expires = 0

	if aura.OnDoneIteration != nil {
		aura.OnDoneIteration(aura, sim)
	}
}

func (aura *Aura) IsActive() bool {
	if aura == nil {
		return false
	}
	return aura.active
}

func (aura *Aura) Refresh(sim *Simulation) {
	if aura.Duration == NeverExpires {
		aura.expires = NeverExpires
	} else {
		aura.expires = sim.CurrentTime + aura.Duration
		if aura.expires < aura.Unit.minExpires {
			aura.Unit.minExpires = aura.expires
			sim.rescheduleTracker(aura.expires)
		}
	}
}

func (aura *Aura) GetStacks() int32 {
	if aura == nil {
		return 0
	}
	return aura.stacks
}

func (aura *Aura) SetStacks(sim *Simulation, newStacks int32) {
	if !aura.IsActive() && newStacks != 0 {
		panic("Trying to set non-zero stacks on inactive aura!")
	}
	if newStacks < 0 {
		panic("SetStacks newStacks cannot be negative but is " + strconv.Itoa(int(newStacks)))
	}
	if aura.MaxStacks == 0 {
		panic("MaxStacks required to set Aura stacks: " + aura.Label)
	}
	oldStacks := aura.stacks
	newStacks = min(newStacks, aura.MaxStacks)

	if oldStacks == newStacks {
		return
	}

	if sim.Log != nil {
		aura.Unit.Log(sim, "%s stacks: %d --> %d", aura.ActionID, oldStacks, newStacks)
	}
	aura.stacks = newStacks
	if aura.OnStacksChange != nil {
		aura.OnStacksChange(aura, sim, oldStacks, newStacks)
	}
	if aura.stacks == 0 {
		aura.Deactivate(sim)
	}
}
func (aura *Aura) AddStack(sim *Simulation) {
	aura.SetStacks(sim, aura.stacks+1)
}
func (aura *Aura) RemoveStack(sim *Simulation) {
	aura.SetStacks(sim, aura.stacks-1)
}

func (aura *Aura) UpdateExpires(newExpires time.Duration) {
	aura.expires = newExpires
}

// The amount of time this aura has been active.
func (aura *Aura) TimeActive(sim *Simulation) time.Duration {
	if aura.IsActive() {
		return sim.CurrentTime - aura.startTime
	} else {
		return 0
	}
}

func (aura *Aura) RemainingDuration(sim *Simulation) time.Duration {
	if aura.expires == NeverExpires {
		return NeverExpires
	} else {
		return aura.expires - sim.CurrentTime
	}
}

func (aura *Aura) StartedAt() time.Duration {
	return aura.startTime
}

func (aura *Aura) ExpiresAt() time.Duration {
	return aura.expires
}

// Adds a handler to be called OnGain, in addition to any current handlers.
func (aura *Aura) ApplyOnGain(newOnGain OnGain) {
	oldOnGain := aura.OnGain
	if oldOnGain == nil {
		aura.OnGain = newOnGain
	} else {
		aura.OnGain = func(aura *Aura, sim *Simulation) {
			oldOnGain(aura, sim)
			newOnGain(aura, sim)
		}
	}
}

// Adds a handler to be called OnExpire, in addition to any current handlers.
func (aura *Aura) ApplyOnExpire(newOnExpire OnExpire) {
	oldOnExpire := aura.OnExpire
	if oldOnExpire == nil {
		aura.OnExpire = newOnExpire
	} else {
		aura.OnExpire = func(aura *Aura, sim *Simulation) {
			oldOnExpire(aura, sim)
			newOnExpire(aura, sim)
		}
	}
}

type AuraFactory func(*Simulation) *Aura

// Callback for doing something on reset.
type ResetEffect func(*Simulation)

// auraTracker is a centralized implementation of CD and Aura tracking.
//
//	This is used by all Units.
type auraTracker struct {
	// Effects to invoke on every sim reset.
	resetEffects []ResetEffect

	*ExclusiveEffectManager

	// All registered auras, both active and inactive.
	auras []*Aura

	aurasByTag map[string][]*Aura

	// IDs of Auras that may expire and are currently active, in no particular order.
	activeAuras []*Aura

	// caches the minimum expires time of all active auras; might be stale (too low) after Deactivate().
	minExpires time.Duration

	// Auras that have a non-nil XXX function set and are currently active.
	onCastCompleteAuras        []*Aura
	onSpellHitDealtAuras       []*Aura
	onSpellHitTakenAuras       []*Aura
	onPeriodicDamageDealtAuras []*Aura
	onPeriodicDamageTakenAuras []*Aura
	onHealDealtAuras           []*Aura
	onHealTakenAuras           []*Aura
	onPeriodicHealDealtAuras   []*Aura
	onPeriodicHealTakenAuras   []*Aura
}

func newAuraTracker() auraTracker {
	return auraTracker{
		resetEffects:           []ResetEffect{},
		ExclusiveEffectManager: &ExclusiveEffectManager{},
		aurasByTag:             make(map[string][]*Aura),
	}
}

func (at *auraTracker) GetAura(label string) *Aura {
	for _, aura := range at.auras {
		if aura.Label == label {
			return aura
		}
	}
	return nil
}
func (at *auraTracker) GetAuras() []*Aura {
	return at.auras
}
func (at *auraTracker) GetAuraByID(actionID ActionID) *Aura {
	for _, aura := range at.auras {
		if aura.ActionID.SameAction(actionID) {
			return aura
		}
	}
	return nil
}
func (at *auraTracker) GetIcdAuraByID(actionID ActionID) *Aura {
	for _, aura := range at.auras {
		if (aura.ActionID.SameAction(actionID) || aura.ActionIDForProc.SameAction(actionID)) && aura.Icd != nil {
			return aura
		}
	}
	return nil
}
func (at *auraTracker) HasAura(label string) bool {
	aura := at.GetAura(label)
	return aura != nil
}
func (at *auraTracker) HasActiveAura(label string) bool {
	aura := at.GetAura(label)
	return aura != nil && aura.IsActive()
}

func (at *auraTracker) registerAura(unit *Unit, aura Aura) *Aura {
	if unit == nil {
		panic("Aura unit is required!")
	}
	if aura.Label == "" {
		panic("Aura label is required!")
	}
	if at.GetAura(aura.Label) != nil {
		panic(fmt.Sprintf("Aura %s already registered!", aura.Label))
	}
	if len(at.auras) > 200 {
		panic(fmt.Sprintf("Over 200 registered auras when registering %s! There is probably an aura being registered every iteration.", aura.Label))
	}

	newAura := &Aura{}
	*newAura = aura
	newAura.Unit = unit
	newAura.Icd = aura.Icd
	newAura.metrics.ID = aura.ActionID
	newAura.activeIndex = Inactive
	newAura.onCastCompleteIndex = Inactive
	newAura.onSpellHitDealtIndex = Inactive
	newAura.onSpellHitTakenIndex = Inactive
	newAura.onPeriodicDamageDealtIndex = Inactive
	newAura.onPeriodicDamageTakenIndex = Inactive
	newAura.onHealDealtIndex = Inactive
	newAura.onHealTakenIndex = Inactive
	newAura.onPeriodicHealDealtIndex = Inactive
	newAura.onPeriodicHealTakenIndex = Inactive

	at.auras = append(at.auras, newAura)
	if newAura.Tag != "" {
		at.aurasByTag[newAura.Tag] = append(at.aurasByTag[newAura.Tag], newAura)
	}

	return newAura
}
func (unit *Unit) RegisterAura(aura Aura) *Aura {
	return unit.auraTracker.registerAura(unit, aura)
}

func (unit *Unit) GetOrRegisterAura(aura Aura) *Aura {
	curAura := unit.GetAura(aura.Label)
	if curAura == nil {
		return unit.RegisterAura(aura)
	} else {
		curAura.Icd = aura.Icd
		curAura.OnCastComplete = aura.OnCastComplete
		curAura.OnSpellHitDealt = aura.OnSpellHitDealt
		curAura.OnSpellHitTaken = aura.OnSpellHitTaken
		curAura.OnPeriodicDamageDealt = aura.OnPeriodicDamageDealt
		curAura.OnPeriodicDamageTaken = aura.OnPeriodicDamageTaken
		curAura.OnHealDealt = aura.OnHealDealt
		curAura.OnHealTaken = aura.OnHealTaken
		curAura.OnPeriodicHealDealt = aura.OnPeriodicHealDealt
		curAura.OnPeriodicHealTaken = aura.OnPeriodicHealTaken
		return curAura
	}
}

func (at *auraTracker) GetAurasWithTag(tag string) []*Aura {
	return at.aurasByTag[tag]
}

func (at *auraTracker) HasAuraWithTag(tag string) bool {
	return len(at.aurasByTag[tag]) > 0
}

func (at *auraTracker) GetActiveAuraWithTag(tag string) *Aura {
	for _, aura := range at.aurasByTag[tag] {
		if aura.active {
			return aura
		}
	}
	return nil
}
func (at *auraTracker) NumActiveAurasWithTag(tag string) int32 {
	count := int32(0)
	for _, aura := range at.aurasByTag[tag] {
		if aura.active {
			count++
		}
	}
	return count
}
func (at *auraTracker) HasActiveAuraWithTag(tag string) bool {
	for _, aura := range at.aurasByTag[tag] {
		if aura.active {
			return true
		}
	}
	return false
}
func (at *auraTracker) HasActiveAuraWithTagExcludingAura(tag string, excludeAura *Aura) bool {
	for _, aura := range at.aurasByTag[tag] {
		if aura.active && aura != excludeAura {
			return true
		}
	}
	return false
}

// Registers a callback to this Character which will be invoked on
// every Sim reset.
func (at *auraTracker) RegisterResetEffect(resetEffect ResetEffect) {
	at.resetEffects = append(at.resetEffects, resetEffect)
}

func (at *auraTracker) reset(sim *Simulation) {
	at.activeAuras = at.activeAuras[:0]
	at.onCastCompleteAuras = at.onCastCompleteAuras[:0]
	at.onSpellHitDealtAuras = at.onSpellHitDealtAuras[:0]
	at.onSpellHitTakenAuras = at.onSpellHitTakenAuras[:0]
	at.onPeriodicDamageDealtAuras = at.onPeriodicDamageDealtAuras[:0]
	at.onPeriodicDamageTakenAuras = at.onPeriodicDamageTakenAuras[:0]
	at.onHealDealtAuras = at.onHealDealtAuras[:0]
	at.onHealTakenAuras = at.onHealTakenAuras[:0]
	at.onPeriodicHealDealtAuras = at.onPeriodicHealDealtAuras[:0]
	at.onPeriodicHealTakenAuras = at.onPeriodicHealTakenAuras[:0]

	for _, resetEffect := range at.resetEffects {
		resetEffect(sim)
	}

	at.minExpires = NeverExpires

	for _, aura := range at.auras {
		aura.reset(sim)
	}
}

// inlineable stub for advance
func (at *auraTracker) tryAdvance(sim *Simulation) time.Duration {
	if sim.CurrentTime < at.minExpires {
		return at.minExpires
	}
	return at.advance(sim)
}

func (at *auraTracker) advance(sim *Simulation) time.Duration {
restart:
	at.minExpires = NeverExpires
	for _, aura := range at.activeAuras {
		if aura.expires <= sim.CurrentTime {
			aura.Deactivate(sim)
			goto restart // activeAuras have changed
		}
		at.minExpires = min(at.minExpires, aura.expires)
	}
	return at.minExpires
}

func (at *auraTracker) expireAll(sim *Simulation) {
restart:
	for _, aura := range at.activeAuras {
		aura.Deactivate(sim)
		goto restart
	}
	at.minExpires = NeverExpires
}

func (at *auraTracker) doneIteration(sim *Simulation) {
	// deactivate all auras, even permanent ones
restart:
	for _, aura := range at.auras {
		if aura.active {
			aura.Deactivate(sim)
			goto restart
		}
	}

	for _, aura := range at.auras {
		aura.doneIteration(sim)
	}

	for _, aura := range at.auras {
		aura.metrics.doneIteration()
	}
}

// Adds a new aura to the simulation. If an aura with the same ID already
// exists it will be replaced with the new one.
func (aura *Aura) Activate(sim *Simulation) {
	aura.metrics.Procs++
	if aura.IsActive() {
		if sim.Log != nil && !aura.ActionID.IsEmptyAction() {
			aura.Unit.Log(sim, "Aura refreshed: %s", aura.ActionID)
		}
		aura.Refresh(sim)
		return
	}

	if aura.Duration == 0 {
		panic("Aura with 0 duration")
	}

	// Activate exclusive effects.
	// If there is already an active aura stronger than this one, then this one
	// will be blocked.
	for i, ee := range aura.ExclusiveEffects {
		if !ee.Activate(sim) {
			// Go back and deactivate the effects.
			if i > 0 {
				for j := 0; j < i; j++ {
					aura.ExclusiveEffects[j].Deactivate(sim)
				}
			}
			return
		}
	}

	aura.active = true
	aura.startTime = sim.CurrentTime
	aura.Refresh(sim)

	if aura.Duration != NeverExpires {
		aura.activeIndex = int32(len(aura.Unit.activeAuras))
		aura.Unit.activeAuras = append(aura.Unit.activeAuras, aura)
	}

	if aura.OnCastComplete != nil {
		aura.onCastCompleteIndex = int32(len(aura.Unit.onCastCompleteAuras))
		aura.Unit.onCastCompleteAuras = append(aura.Unit.onCastCompleteAuras, aura)
	}

	if aura.OnSpellHitDealt != nil {
		aura.onSpellHitDealtIndex = int32(len(aura.Unit.onSpellHitDealtAuras))
		aura.Unit.onSpellHitDealtAuras = append(aura.Unit.onSpellHitDealtAuras, aura)
	}

	if aura.OnSpellHitTaken != nil {
		aura.onSpellHitTakenIndex = int32(len(aura.Unit.onSpellHitTakenAuras))
		aura.Unit.onSpellHitTakenAuras = append(aura.Unit.onSpellHitTakenAuras, aura)
	}

	if aura.OnPeriodicDamageDealt != nil {
		aura.onPeriodicDamageDealtIndex = int32(len(aura.Unit.onPeriodicDamageDealtAuras))
		aura.Unit.onPeriodicDamageDealtAuras = append(aura.Unit.onPeriodicDamageDealtAuras, aura)
	}

	if aura.OnPeriodicDamageTaken != nil {
		aura.onPeriodicDamageTakenIndex = int32(len(aura.Unit.onPeriodicDamageTakenAuras))
		aura.Unit.onPeriodicDamageTakenAuras = append(aura.Unit.onPeriodicDamageTakenAuras, aura)
	}

	if aura.OnHealDealt != nil {
		aura.onHealDealtIndex = int32(len(aura.Unit.onHealDealtAuras))
		aura.Unit.onHealDealtAuras = append(aura.Unit.onHealDealtAuras, aura)
	}

	if aura.OnHealTaken != nil {
		aura.onHealTakenIndex = int32(len(aura.Unit.onHealTakenAuras))
		aura.Unit.onHealTakenAuras = append(aura.Unit.onHealTakenAuras, aura)
	}

	if aura.OnPeriodicHealDealt != nil {
		aura.onPeriodicHealDealtIndex = int32(len(aura.Unit.onPeriodicHealDealtAuras))
		aura.Unit.onPeriodicHealDealtAuras = append(aura.Unit.onPeriodicHealDealtAuras, aura)
	}

	if aura.OnPeriodicHealTaken != nil {
		aura.onPeriodicHealTakenIndex = int32(len(aura.Unit.onPeriodicHealTakenAuras))
		aura.Unit.onPeriodicHealTakenAuras = append(aura.Unit.onPeriodicHealTakenAuras, aura)
	}

	if sim.Log != nil && !aura.ActionID.IsEmptyAction() {
		aura.Unit.Log(sim, "Aura gained: %s", aura.ActionID)
	}

	// don't invoke possible callbacks until the internal state is consistent
	if aura.OnGain != nil {
		aura.OnGain(aura, sim)
	}
}

// Remove an aura by its ID
func (aura *Aura) Deactivate(sim *Simulation) {
	if !aura.active {
		return
	}
	aura.active = false

	if !aura.ActionID.IsEmptyAction() {
		if sim.CurrentTime > aura.expires {
			aura.metrics.Uptime += aura.expires - max(aura.startTime, 0)
		} else {
			aura.metrics.Uptime += sim.CurrentTime - max(aura.startTime, 0)
		}
	}

	if sim.Log != nil && !aura.ActionID.IsEmptyAction() {
		aura.Unit.Log(sim, "Aura faded: %s", aura.ActionID)
	}

	aura.expires = 0
	if aura.activeIndex != Inactive {
		removeActiveIndex := aura.activeIndex
		aura.Unit.activeAuras = removeBySwappingToBack(aura.Unit.activeAuras, removeActiveIndex)
		if removeActiveIndex < int32(len(aura.Unit.activeAuras)) {
			aura.Unit.activeAuras[removeActiveIndex].activeIndex = removeActiveIndex
		}
		aura.activeIndex = Inactive
	}

	if aura.onCastCompleteIndex != Inactive {
		removeOnCastCompleteIndex := aura.onCastCompleteIndex
		aura.Unit.onCastCompleteAuras = removeBySwappingToBack(aura.Unit.onCastCompleteAuras, removeOnCastCompleteIndex)
		if removeOnCastCompleteIndex < int32(len(aura.Unit.onCastCompleteAuras)) {
			aura.Unit.onCastCompleteAuras[removeOnCastCompleteIndex].onCastCompleteIndex = removeOnCastCompleteIndex
		}
		aura.onCastCompleteIndex = Inactive
	}

	if aura.onSpellHitDealtIndex != Inactive {
		removeOnSpellHitDealtIndex := aura.onSpellHitDealtIndex
		aura.Unit.onSpellHitDealtAuras = removeBySwappingToBack(aura.Unit.onSpellHitDealtAuras, removeOnSpellHitDealtIndex)
		if removeOnSpellHitDealtIndex < int32(len(aura.Unit.onSpellHitDealtAuras)) {
			aura.Unit.onSpellHitDealtAuras[removeOnSpellHitDealtIndex].onSpellHitDealtIndex = removeOnSpellHitDealtIndex
		}
		aura.onSpellHitDealtIndex = Inactive
	}

	if aura.onSpellHitTakenIndex != Inactive {
		removeOnSpellHitTakenIndex := aura.onSpellHitTakenIndex
		aura.Unit.onSpellHitTakenAuras = removeBySwappingToBack(aura.Unit.onSpellHitTakenAuras, removeOnSpellHitTakenIndex)
		if removeOnSpellHitTakenIndex < int32(len(aura.Unit.onSpellHitTakenAuras)) {
			aura.Unit.onSpellHitTakenAuras[removeOnSpellHitTakenIndex].onSpellHitTakenIndex = removeOnSpellHitTakenIndex
		}
		aura.onSpellHitTakenIndex = Inactive
	}

	if aura.onPeriodicDamageDealtIndex != Inactive {
		removeOnPeriodicDamageDealt := aura.onPeriodicDamageDealtIndex
		aura.Unit.onPeriodicDamageDealtAuras = removeBySwappingToBack(aura.Unit.onPeriodicDamageDealtAuras, removeOnPeriodicDamageDealt)
		if removeOnPeriodicDamageDealt < int32(len(aura.Unit.onPeriodicDamageDealtAuras)) {
			aura.Unit.onPeriodicDamageDealtAuras[removeOnPeriodicDamageDealt].onPeriodicDamageDealtIndex = removeOnPeriodicDamageDealt
		}
		aura.onPeriodicDamageDealtIndex = Inactive
	}

	if aura.onPeriodicDamageTakenIndex != Inactive {
		removeOnPeriodicDamageTaken := aura.onPeriodicDamageTakenIndex
		aura.Unit.onPeriodicDamageTakenAuras = removeBySwappingToBack(aura.Unit.onPeriodicDamageTakenAuras, removeOnPeriodicDamageTaken)
		if removeOnPeriodicDamageTaken < int32(len(aura.Unit.onPeriodicDamageTakenAuras)) {
			aura.Unit.onPeriodicDamageTakenAuras[removeOnPeriodicDamageTaken].onPeriodicDamageTakenIndex = removeOnPeriodicDamageTaken
		}
		aura.onPeriodicDamageTakenIndex = Inactive
	}

	if aura.onHealDealtIndex != Inactive {
		removeOnHealDealtIndex := aura.onHealDealtIndex
		aura.Unit.onHealDealtAuras = removeBySwappingToBack(aura.Unit.onHealDealtAuras, removeOnHealDealtIndex)
		if removeOnHealDealtIndex < int32(len(aura.Unit.onHealDealtAuras)) {
			aura.Unit.onHealDealtAuras[removeOnHealDealtIndex].onHealDealtIndex = removeOnHealDealtIndex
		}
		aura.onHealDealtIndex = Inactive
	}

	if aura.onHealTakenIndex != Inactive {
		removeOnHealTakenIndex := aura.onHealTakenIndex
		aura.Unit.onHealTakenAuras = removeBySwappingToBack(aura.Unit.onHealTakenAuras, removeOnHealTakenIndex)
		if removeOnHealTakenIndex < int32(len(aura.Unit.onHealTakenAuras)) {
			aura.Unit.onHealTakenAuras[removeOnHealTakenIndex].onHealTakenIndex = removeOnHealTakenIndex
		}
		aura.onHealTakenIndex = Inactive
	}

	if aura.onPeriodicHealDealtIndex != Inactive {
		removeOnPeriodicHealDealt := aura.onPeriodicHealDealtIndex
		aura.Unit.onPeriodicHealDealtAuras = removeBySwappingToBack(aura.Unit.onPeriodicHealDealtAuras, removeOnPeriodicHealDealt)
		if removeOnPeriodicHealDealt < int32(len(aura.Unit.onPeriodicHealDealtAuras)) {
			aura.Unit.onPeriodicHealDealtAuras[removeOnPeriodicHealDealt].onPeriodicHealDealtIndex = removeOnPeriodicHealDealt
		}
		aura.onPeriodicHealDealtIndex = Inactive
	}

	if aura.onPeriodicHealTakenIndex != Inactive {
		removeOnPeriodicHealTaken := aura.onPeriodicHealTakenIndex
		aura.Unit.onPeriodicHealTakenAuras = removeBySwappingToBack(aura.Unit.onPeriodicHealTakenAuras, removeOnPeriodicHealTaken)
		if removeOnPeriodicHealTaken < int32(len(aura.Unit.onPeriodicHealTakenAuras)) {
			aura.Unit.onPeriodicHealTakenAuras[removeOnPeriodicHealTaken].onPeriodicHealTakenIndex = removeOnPeriodicHealTaken
		}
		aura.onPeriodicHealTakenIndex = Inactive
	}

	// don't invoke possible callbacks until the internal state is consistent
	if aura.stacks != 0 {
		aura.SetStacks(sim, 0)
	}

	// Deactivate exclusive effects.
	for _, ee := range aura.ExclusiveEffects {
		ee.Deactivate(sim)
	}

	if aura.OnExpire != nil {
		aura.OnExpire(aura, sim)
	}
}

// Constant-time removal from slice by swapping with the last element before removing.
func removeBySwappingToBack[T any, U constraints.Integer](arr []T, removeIdx U) []T {
	arr[removeIdx] = arr[len(arr)-1]
	return arr[:len(arr)-1]
}

// Invokes the OnCastComplete event for all tracked Auras.
func (at *auraTracker) OnCastComplete(sim *Simulation, spell *Spell) {
	for _, aura := range at.onCastCompleteAuras {
		aura.OnCastComplete(aura, sim, spell)
	}
}

// Invokes the OnSpellHit event for all tracked Auras.
func (at *auraTracker) OnSpellHitDealt(sim *Simulation, spell *Spell, result *SpellResult) {
	for _, aura := range at.onSpellHitDealtAuras {
		// this check is to handle a case where auras are deactivated during iteration.
		if !aura.active {
			continue
		}
		aura.OnSpellHitDealt(aura, sim, spell, result)
	}
}
func (at *auraTracker) OnSpellHitTaken(sim *Simulation, spell *Spell, result *SpellResult) {
	for _, aura := range at.onSpellHitTakenAuras {
		// this check is to handle a case where auras are deactivated during iteration.
		if !aura.active {
			continue
		}
		aura.OnSpellHitTaken(aura, sim, spell, result)
	}
}

// Invokes the OnPeriodicDamage
//
//	As a debuff when target is being hit by dot.
//	As a buff when caster's dots are ticking.
func (at *auraTracker) OnPeriodicDamageDealt(sim *Simulation, spell *Spell, result *SpellResult) {
	for _, aura := range at.onPeriodicDamageDealtAuras {
		aura.OnPeriodicDamageDealt(aura, sim, spell, result)
	}
}
func (at *auraTracker) OnPeriodicDamageTaken(sim *Simulation, spell *Spell, result *SpellResult) {
	for _, aura := range at.onPeriodicDamageTakenAuras {
		aura.OnPeriodicDamageTaken(aura, sim, spell, result)
	}
}

// Invokes the OnHeal event for all tracked Auras.
func (at *auraTracker) OnHealDealt(sim *Simulation, spell *Spell, result *SpellResult) {
	for _, aura := range at.onHealDealtAuras {
		// this check is to handle a case where auras are deactivated during iteration.
		if !aura.active {
			continue
		}
		aura.OnHealDealt(aura, sim, spell, result)
	}
}
func (at *auraTracker) OnHealTaken(sim *Simulation, spell *Spell, result *SpellResult) {
	for _, aura := range at.onHealTakenAuras {
		// this check is to handle a case where auras are deactivated during iteration.
		if !aura.active {
			continue
		}
		aura.OnHealTaken(aura, sim, spell, result)
	}
}

// Invokes the OnPeriodicHeal
//
//	As a debuff when target is being hit by dot.
//	As a buff when caster's dots are ticking.
func (at *auraTracker) OnPeriodicHealDealt(sim *Simulation, spell *Spell, result *SpellResult) {
	for _, aura := range at.onPeriodicHealDealtAuras {
		aura.OnPeriodicHealDealt(aura, sim, spell, result)
	}
}
func (at *auraTracker) OnPeriodicHealTaken(sim *Simulation, spell *Spell, result *SpellResult) {
	for _, aura := range at.onPeriodicHealTakenAuras {
		aura.OnPeriodicHealTaken(aura, sim, spell, result)
	}
}

func (at *auraTracker) GetMetricsProto() []*proto.AuraMetrics {
	metrics := make([]*proto.AuraMetrics, 0, len(at.auras))

	for _, aura := range at.auras {
		if !aura.metrics.ID.IsEmptyAction() {
			metrics = append(metrics, aura.metrics.ToProto())
		}
	}

	return metrics
}

type AuraArray []*Aura

func (auras AuraArray) Get(target *Unit) *Aura {
	return auras[target.UnitIndex]
}

func (caster *Unit) NewAllyAuraArray(makeAura func(*Unit) *Aura) AuraArray {
	auras := make([]*Aura, len(caster.Env.AllUnits))
	for _, target := range caster.Env.AllUnits {
		if target.Type != EnemyUnit {
			auras[target.UnitIndex] = makeAura(target)
		}
	}
	return auras
}

func (caster *Unit) NewEnemyAuraArray(makeAura func(*Unit) *Aura) AuraArray {
	auras := make([]*Aura, len(caster.Env.AllUnits))
	for _, target := range caster.Env.AllUnits {
		if target.Type == EnemyUnit {
			auras[target.UnitIndex] = makeAura(target)
		}
	}
	return auras
}
