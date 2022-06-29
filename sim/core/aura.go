package core

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/wowsims/tbc/sim/core/proto"
	"github.com/wowsims/tbc/sim/core/stats"
)

const NeverExpires = time.Duration(math.MaxInt64)

type OnInit func(aura *Aura, sim *Simulation)
type OnReset func(aura *Aura, sim *Simulation)
type OnDoneIteration func(aura *Aura, sim *Simulation)
type OnGain func(aura *Aura, sim *Simulation)
type OnExpire func(aura *Aura, sim *Simulation)
type OnStacksChange func(aura *Aura, sim *Simulation, oldStacks int32, newStacks int32)

const Inactive = -1

// Aura lifecycle:
//
// myAura := unit.RegisterAura(myAuraConfig)
// myAura.Activate(sim)
// myAura.SetStacks(sim, 3)
// myAura.Refresh(sim)
// myAura.Deactivate(sim)
type Aura struct {
	// String label for this Aura. Gauranteed to be unique among the Auras for a single Unit.
	Label string

	// For easily grouping auras.
	Tag string

	ActionID ActionID // If set, metrics will be tracked for this aura.

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

	// The number of stacks, or charges, of this aura. If this aura doesn't care
	// about charges, is just 0.
	stacks    int32
	MaxStacks int32

	// If nonzero, activation of this aura will deactivate other auras with the
	// same Tag and equal or lower Priority.
	Priority float64

	// Lifecycle callbacks.
	OnInit          OnInit
	OnReset         OnReset
	OnDoneIteration OnDoneIteration
	OnGain          OnGain
	OnExpire        OnExpire
	OnStacksChange  OnStacksChange // Invoked when the number of stacks of this aura changes.

	OnCastComplete        OnCastComplete   // Invoked when a spell cast completes casting, before results are calculated.
	OnSpellHitDealt       OnSpellHit       // Invoked when a spell hits and this unit is the caster.
	OnSpellHitTaken       OnSpellHit       // Invoked when a spell hits and this unit is the target.
	OnPeriodicDamageDealt OnPeriodicDamage // Invoked when a dot tick occurs and this unit is the caster.
	OnPeriodicDamageTaken OnPeriodicDamage // Invoked when a dot tick occurs and this unit is the target.

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
	return aura.active
}

func (aura *Aura) Refresh(sim *Simulation) {
	if aura.Duration == NeverExpires {
		aura.expires = NeverExpires
	} else {
		aura.expires = sim.CurrentTime + aura.Duration
		aura.Unit.minExpires = 0
	}
}

func (aura *Aura) GetStacks() int32 {
	return aura.stacks
}

func (aura *Aura) SetStacks(sim *Simulation, newStacks int32) {
	if newStacks < 0 {
		panic("SetStacks newStacks cannot be negative but is " + strconv.Itoa(int(newStacks)))
	}
	if aura.MaxStacks == 0 {
		panic("MaxStacks required to set Aura stacks: " + aura.Label)
	}
	oldStacks := aura.stacks
	newStacks = MinInt32(newStacks, aura.MaxStacks)

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

func (aura *Aura) RemainingDuration(sim *Simulation) time.Duration {
	if aura.expires == NeverExpires {
		return NeverExpires
	} else {
		return aura.expires - sim.CurrentTime
	}
}

func (aura *Aura) ExpiresAt() time.Duration {
	return aura.expires
}

type AuraFactory func(*Simulation) *Aura

// Callback for doing something on reset.
type ResetEffect func(*Simulation)

// auraTracker is a centralized implementation of CD and Aura tracking.
//  This is used by all Units.
type auraTracker struct {
	// Effects to invoke on every sim reset.
	resetEffects []ResetEffect

	// Maps MagicIDs to sim duration at which CD is done. Using array for perf.
	cooldowns []time.Duration

	// All registered auras, both active and inactive.
	auras []*Aura

	aurasByTag map[string][]*Aura

	// IDs of Auras that may expire and are currently active, in no particular order.
	activeAuras []*Aura

	// caches the minimum expires time of all active auras; reset to 0 on Activate(), Deactivate(), and Refresh()
	minExpires time.Duration

	// Auras that have a non-nil XXX function set and are currently active.
	onCastCompleteAuras        []*Aura
	onSpellHitDealtAuras       []*Aura
	onSpellHitTakenAuras       []*Aura
	onPeriodicDamageDealtAuras []*Aura
	onPeriodicDamageTakenAuras []*Aura
}

func newAuraTracker() auraTracker {
	return auraTracker{
		resetEffects:               []ResetEffect{},
		activeAuras:                make([]*Aura, 0, 16),
		onCastCompleteAuras:        make([]*Aura, 0, 16),
		onSpellHitDealtAuras:       make([]*Aura, 0, 16),
		onSpellHitTakenAuras:       make([]*Aura, 0, 16),
		onPeriodicDamageDealtAuras: make([]*Aura, 0, 16),
		onPeriodicDamageTakenAuras: make([]*Aura, 0, 16),
		auras:                      make([]*Aura, 0, 16),
		aurasByTag:                 make(map[string][]*Aura),
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
	if aura.Priority != 0 && aura.Tag == "" {
		panic("Aura.Priority requires Aura.Tag also be set")
	}
	if at.GetAura(aura.Label) != nil {
		panic(fmt.Sprintf("Aura %s already registered!", aura.Label))
	}
	if len(at.auras) > 100 {
		panic(fmt.Sprintf("Over 100 registered auras when registering %s! There is probably an aura being registered every iteration.", aura.Label))
	}

	newAura := &Aura{}
	*newAura = aura
	newAura.Unit = unit
	newAura.metrics.ID = aura.ActionID
	newAura.activeIndex = Inactive
	newAura.onCastCompleteIndex = Inactive
	newAura.onSpellHitDealtIndex = Inactive
	newAura.onSpellHitTakenIndex = Inactive
	newAura.onPeriodicDamageDealtIndex = Inactive
	newAura.onPeriodicDamageTakenIndex = Inactive

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
		curAura.OnCastComplete = aura.OnCastComplete
		curAura.OnSpellHitDealt = aura.OnSpellHitDealt
		curAura.OnSpellHitTaken = aura.OnSpellHitTaken
		curAura.OnPeriodicDamageDealt = aura.OnPeriodicDamageDealt
		curAura.OnPeriodicDamageTaken = aura.OnPeriodicDamageTaken
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
func (at *auraTracker) HasActiveAuraWithTag(tag string) bool {
	for _, aura := range at.aurasByTag[tag] {
		if aura.active {
			return true
		}
	}
	return false
}

// Returns if an aura should be refreshed at a specific priority, i.e. the aura
// is about to expire AND the replacement aura has at least as high priority.
//
// This is used to decide whether to refresh effects with multiple strengths,
// like Thunder Clap/Deathfrost or Faerie Fire ranks.
func (at *auraTracker) ShouldRefreshAuraWithTagAtPriority(sim *Simulation, tag string, priority float64, refreshWindow time.Duration) bool {
	activeAura := at.GetActiveAuraWithTag(tag)

	return activeAura == nil ||
		priority > activeAura.Priority ||
		(priority == activeAura.Priority && activeAura.RemainingDuration(sim) <= refreshWindow)
}

// Registers a callback to this Character which will be invoked on
// every Sim reset.
func (at *auraTracker) RegisterResetEffect(resetEffect ResetEffect) {
	at.resetEffects = append(at.resetEffects, resetEffect)
}

func (at *auraTracker) init(sim *Simulation) {
	// Auras are initialized later, on their first reset().
}

func (at *auraTracker) reset(sim *Simulation) {
	at.activeAuras = at.activeAuras[:0]
	at.onCastCompleteAuras = at.onCastCompleteAuras[:0]
	at.onSpellHitDealtAuras = at.onSpellHitDealtAuras[:0]
	at.onSpellHitTakenAuras = at.onSpellHitTakenAuras[:0]
	at.onPeriodicDamageDealtAuras = at.onPeriodicDamageDealtAuras[:0]
	at.onPeriodicDamageTakenAuras = at.onPeriodicDamageTakenAuras[:0]

	for _, resetEffect := range at.resetEffects {
		resetEffect(sim)
	}

	for _, aura := range at.auras {
		aura.reset(sim)
	}
}

func (at *auraTracker) advance(sim *Simulation) {
	if at.minExpires > sim.CurrentTime {
		return
	}

restart:
	minExpires := NeverExpires
	for _, aura := range at.activeAuras {
		if aura.expires <= sim.CurrentTime && aura.expires != 0 {
			aura.Deactivate(sim)
			goto restart // activeAuras have changed
		}
		if aura.expires < minExpires {
			minExpires = aura.expires
		}
	}
	at.minExpires = minExpires
}

func (at *auraTracker) doneIteration(sim *Simulation) {
	// Expire all the remaining auras. Need to keep looping because sometimes
	// expiring auras can trigger other auras.
	foundUnexpired := true
	for foundUnexpired {
		foundUnexpired = false
		for _, aura := range at.auras {
			if aura.IsActive() {
				foundUnexpired = true
				aura.Deactivate(sim)
			}
		}
	}

	for _, aura := range at.auras {
		aura.doneIteration(sim)
	}

	// Add metrics for any auras that are still active.
	for _, aura := range at.auras {
		aura.metrics.doneIteration()
	}
}

// Adds a new aura to the simulation. If an aura with the same ID already
// exists it will be replaced with the new one.
func (aura *Aura) Activate(sim *Simulation) {
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

	// If there is already an active aura stronger than this one, then this one
	// is blocked.
	if aura.Tag != "" {
		for _, otherAura := range aura.Unit.GetAurasWithTag(aura.Tag) {
			if otherAura.Priority > aura.Priority && otherAura.active {
				return
			}
		}
	}

	// Remove weaker versions of the same aura.
	if aura.Priority != 0 {
		for _, otherAura := range aura.Unit.GetAurasWithTag(aura.Tag) {
			if otherAura.Priority <= aura.Priority && otherAura != aura {
				otherAura.Deactivate(sim)
			}
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

	if sim.Log != nil && !aura.ActionID.IsEmptyAction() {
		aura.Unit.Log(sim, "Aura gained: %s", aura.ActionID)
	}

	if aura.OnGain != nil {
		aura.OnGain(aura, sim)
	}
}

// Moves an Aura to the front of the list of active Auras, so its callbacks are invoked first.
func (aura *Aura) Prioritize() {
	if aura.onCastCompleteIndex > 0 {
		otherAura := aura.Unit.onCastCompleteAuras[0]
		aura.Unit.onCastCompleteAuras[0] = aura
		aura.Unit.onCastCompleteAuras[len(aura.Unit.onCastCompleteAuras)-1] = otherAura
		otherAura.onCastCompleteIndex = aura.onCastCompleteIndex
		aura.onCastCompleteIndex = 0
	}

	if aura.onSpellHitDealtIndex > 0 {
		otherAura := aura.Unit.onSpellHitDealtAuras[0]
		aura.Unit.onSpellHitDealtAuras[0] = aura
		aura.Unit.onSpellHitDealtAuras[len(aura.Unit.onSpellHitDealtAuras)-1] = otherAura
		otherAura.onSpellHitDealtIndex = aura.onSpellHitDealtIndex
		aura.onSpellHitDealtIndex = 0
	}

	if aura.onSpellHitTakenIndex > 0 {
		otherAura := aura.Unit.onSpellHitTakenAuras[0]
		aura.Unit.onSpellHitTakenAuras[0] = aura
		aura.Unit.onSpellHitTakenAuras[len(aura.Unit.onSpellHitTakenAuras)-1] = otherAura
		otherAura.onSpellHitTakenIndex = aura.onSpellHitTakenIndex
		aura.onSpellHitTakenIndex = 0
	}

	if aura.onPeriodicDamageDealtIndex > 0 {
		otherAura := aura.Unit.onPeriodicDamageDealtAuras[0]
		aura.Unit.onPeriodicDamageDealtAuras[0] = aura
		aura.Unit.onPeriodicDamageDealtAuras[len(aura.Unit.onPeriodicDamageDealtAuras)-1] = otherAura
		otherAura.onPeriodicDamageDealtIndex = aura.onPeriodicDamageDealtIndex
		aura.onPeriodicDamageDealtIndex = 0
	}

	if aura.onPeriodicDamageTakenIndex > 0 {
		otherAura := aura.Unit.onPeriodicDamageTakenAuras[0]
		aura.Unit.onPeriodicDamageTakenAuras[0] = aura
		aura.Unit.onPeriodicDamageTakenAuras[len(aura.Unit.onPeriodicDamageTakenAuras)-1] = otherAura
		otherAura.onPeriodicDamageTakenIndex = aura.onPeriodicDamageTakenIndex
		aura.onPeriodicDamageTakenIndex = 0
	}
}

// Remove an aura by its ID
func (aura *Aura) Deactivate(sim *Simulation) {
	if !aura.active {
		return
	}
	aura.active = false

	if aura.stacks != 0 {
		aura.SetStacks(sim, 0)
	}
	if aura.OnExpire != nil {
		aura.OnExpire(aura, sim)
	}

	if !aura.ActionID.IsEmptyAction() {
		if sim.CurrentTime > aura.expires {
			aura.metrics.Uptime += aura.expires - aura.startTime
		} else {
			aura.metrics.Uptime += sim.CurrentTime - aura.startTime
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

		aura.Unit.minExpires = 0
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
}

// Constant-time removal from slice by swapping with the last element before removing.
func removeBySwappingToBack(arr []*Aura, removeIdx int32) []*Aura {
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
func (at *auraTracker) OnSpellHitDealt(sim *Simulation, spell *Spell, spellEffect *SpellEffect) {
	for _, aura := range at.onSpellHitDealtAuras {
		// this check is to handle a case where auras are deactivated during iteration.
		if !aura.active {
			continue
		}
		aura.OnSpellHitDealt(aura, sim, spell, spellEffect)
	}
}
func (at *auraTracker) OnSpellHitTaken(sim *Simulation, spell *Spell, spellEffect *SpellEffect) {
	for _, aura := range at.onSpellHitTakenAuras {
		// this check is to handle a case where auras are deactivated during iteration.
		if !aura.active {
			continue
		}
		aura.OnSpellHitTaken(aura, sim, spell, spellEffect)
	}
}

// Invokes the OnPeriodicDamage
//   As a debuff when target is being hit by dot.
//   As a buff when caster's dots are ticking.
func (at *auraTracker) OnPeriodicDamageDealt(sim *Simulation, spell *Spell, spellEffect *SpellEffect) {
	for _, aura := range at.onPeriodicDamageDealtAuras {
		aura.OnPeriodicDamageDealt(aura, sim, spell, spellEffect)
	}
}
func (at *auraTracker) OnPeriodicDamageTaken(sim *Simulation, spell *Spell, spellEffect *SpellEffect) {
	for _, aura := range at.onPeriodicDamageTakenAuras {
		aura.OnPeriodicDamageTaken(aura, sim, spell, spellEffect)
	}
}

func (at *auraTracker) GetMetricsProto(numIterations int32) []*proto.AuraMetrics {
	metrics := make([]*proto.AuraMetrics, 0, len(at.auras))

	for _, aura := range at.auras {
		if !aura.metrics.ID.IsEmptyAction() {
			metrics = append(metrics, aura.metrics.ToProto(numIterations))
		}
	}

	return metrics
}

// Returns the same Aura for chaining.
func MakePermanent(aura *Aura) *Aura {
	aura.Duration = NeverExpires
	if aura.OnReset == nil {
		aura.OnReset = func(aura *Aura, sim *Simulation) {
			aura.Activate(sim)
		}
	} else {
		oldOnReset := aura.OnReset
		aura.OnReset = func(aura *Aura, sim *Simulation) {
			oldOnReset(aura, sim)
			aura.Activate(sim)
		}
	}
	return aura
}

// Helper for the common case of making an aura that adds stats.
func (character *Character) NewTemporaryStatsAura(auraLabel string, actionID ActionID, tempStats stats.Stats, duration time.Duration) *Aura {
	return character.NewTemporaryStatsAuraWrapped(auraLabel, actionID, tempStats, duration, nil)
}

// Alternative that allows modifying the Aura config.
func (character *Character) NewTemporaryStatsAuraWrapped(auraLabel string, actionID ActionID, tempStats stats.Stats, duration time.Duration, modConfig func(*Aura)) *Aura {
	var buffs stats.Stats
	var unbuffs stats.Stats

	config := Aura{
		Label:    auraLabel,
		ActionID: actionID,
		Duration: duration,
		OnInit: func(aura *Aura, sim *Simulation) {
			buffs = character.ApplyStatDependencies(tempStats)
			unbuffs = buffs.Multiply(-1)
		},
		OnGain: func(aura *Aura, sim *Simulation) {
			character.AddStatsDynamic(sim, buffs)
			if sim.Log != nil {
				character.Log(sim, "Gained %s from %s.", buffs.FlatString(), actionID)
			}
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			if sim.Log != nil {
				character.Log(sim, "Lost %s from fading %s.", buffs.FlatString(), actionID)
			}
			character.AddStatsDynamic(sim, unbuffs)
		},
	}

	if modConfig != nil {
		modConfig(&config)
	}

	return character.GetOrRegisterAura(config)
}
