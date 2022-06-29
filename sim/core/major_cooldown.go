package core

import (
	"sort"
	"time"

	"github.com/wowsims/tbc/sim/core/proto"
	"github.com/wowsims/tbc/sim/core/stats"
)

const (
	CooldownPriorityLow       = -1.0
	CooldownPriorityDefault   = 0.0
	CooldownPriorityDrums     = 2.0
	CooldownPriorityBloodlust = 1.0
)

type CooldownType byte

const (
	CooldownTypeUnknown CooldownType = 0
	CooldownTypeMana    CooldownType = 1 << iota
	CooldownTypeDPS
	CooldownTypeSurvival
	CooldownTypeUsableShapeShifted
)

func (ct CooldownType) Matches(other CooldownType) bool {
	return (ct & other) != 0
}

// Condition for whether a cooldown can/should be activated.
// Returning false prevents the cooldown from being activated.
type CooldownActivationCondition func(*Simulation, *Character) bool

// Function for activating a cooldown.
// Returns whether the activation was successful.
type CooldownActivation func(*Simulation, *Character)

// Function for making a CooldownActivation.
//
// We need a function that returns a CooldownActivation rather than a
// CooldownActivation, so captured local variables can be reset on Sim reset.
type CooldownActivationFactory func(*Simulation) CooldownActivation

type MajorCooldown struct {
	// Spell that is cast when this MCD is activated.
	Spell *Spell

	// Cooldowns with higher priority get used first. This is important when some
	// cooldowns have a non-zero cast time. For example, Drums should be used
	// before Bloodlust.
	Priority float64

	// Internal category, used for filtering. For example, mages want to disable
	// all DPS cooldowns during their regen rotation.
	Type CooldownType

	// Whether the cooldown meets all hard requirements for activation (e.g. resource cost).
	// Note chat whether the cooldown is off CD is automatically checked, so it does not
	// need to be checked again by this function.
	CanActivate CooldownActivationCondition

	// Whether the cooldown meets all optional conditions for activation. These
	// conditions will be ignored when the user specifies their own activation time.
	// This is for things like mana thresholds, which are optimizations for better
	// automatic timing.
	ShouldActivate CooldownActivationCondition

	// Factory for creating the activate function on every Sim reset.
	ActivationFactory CooldownActivationFactory

	// Fixed timings at which to use this cooldown. If these are specified, they
	// are used instead of ShouldActivate.
	timings []time.Duration

	// Number of times this MCD was used so far in the current iteration.
	numUsages int

	// Internal lambda function to use the cooldown.
	activate CooldownActivation

	// Whether this MCD is currently disabled.
	disabled bool
}

func (mcd *MajorCooldown) ReadyAt() time.Duration {
	return mcd.Spell.ReadyAt()
}

func (mcd *MajorCooldown) IsReady(sim *Simulation) bool {
	return mcd.Spell.IsReady(sim)
}

func (mcd *MajorCooldown) TimeToReady(sim *Simulation) time.Duration {
	return mcd.Spell.TimeToReady(sim)
}

func (mcd *MajorCooldown) IsEnabled() bool {
	return !mcd.disabled
}

func (mcd *MajorCooldown) GetTimings() []time.Duration {
	return mcd.timings[:]
}

// Public version of TryActivate for manual activation by Agent code.
// Note that this version will work even if the MCD is disabled.
func (mcd *MajorCooldown) TryActivate(sim *Simulation, character *Character) bool {
	return mcd.tryActivateHelper(sim, character)
}

func (mcd *MajorCooldown) tryActivateInternal(sim *Simulation, character *Character) bool {
	if mcd.disabled {
		return false
	}

	return mcd.tryActivateHelper(sim, character)
}

// Activates this MCD, if all the conditions pass.
// Returns whether the MCD was activated.
func (mcd *MajorCooldown) tryActivateHelper(sim *Simulation, character *Character) bool {
	if mcd.Spell.DefaultCast.GCD > 0 && !character.GCD.IsReady(sim) {
		return false
	}

	if !mcd.CanActivate(sim, character) {
		return false
	}

	var shouldActivate bool
	if mcd.numUsages < len(mcd.timings) {
		shouldActivate = sim.CurrentTime >= mcd.timings[mcd.numUsages]
	} else {
		shouldActivate = mcd.ShouldActivate(sim, character)
	}

	if shouldActivate {
		mcd.activate(sim, character)
		mcd.numUsages++
		if sim.Log != nil {
			character.Log(sim, "Major cooldown used: %s", mcd.Spell.ActionID)
		}
	}

	return shouldActivate
}

type majorCooldownManager struct {
	// The Character whose cooldowns are being managed.
	character *Character

	// User-specified cooldown configs.
	cooldownConfigs proto.Cooldowns

	// Cached list of major cooldowns sorted by priority, for resetting quickly.
	initialMajorCooldowns []MajorCooldown

	// Major cooldowns, ordered by next available. This should always contain
	// the same cooldows as initialMajorCooldowns, but the order will change over
	// the course of the sim.
	majorCooldowns []*MajorCooldown

	tryUsing bool
	fullSort bool
}

func newMajorCooldownManager(cooldowns *proto.Cooldowns) majorCooldownManager {
	cds := proto.Cooldowns{}
	if cooldowns != nil {
		cds = *cooldowns
	}

	return majorCooldownManager{
		cooldownConfigs: cds,
	}
}

func (mcdm *majorCooldownManager) initialize(character *Character) {
	mcdm.character = character
}

func (mcdm *majorCooldownManager) finalize(character *Character) {
	if mcdm.initialMajorCooldowns == nil {
		mcdm.initialMajorCooldowns = []MajorCooldown{}
	}

	// Match user-specified cooldown configs to existing cooldowns.
	for i, _ := range mcdm.initialMajorCooldowns {
		mcd := &mcdm.initialMajorCooldowns[i]
		mcd.timings = []time.Duration{}

		if mcdm.cooldownConfigs.Cooldowns != nil {
			for _, cooldownConfig := range mcdm.cooldownConfigs.Cooldowns {
				configID := ProtoToActionID(*cooldownConfig.Id)
				if configID.SameAction(mcd.Spell.ActionID) {
					mcd.timings = make([]time.Duration, len(cooldownConfig.Timings))
					for t, timing := range cooldownConfig.Timings {
						mcd.timings[t] = DurationFromSeconds(timing)
					}
					break
				}
			}
		}
	}

	mcdm.majorCooldowns = make([]*MajorCooldown, len(mcdm.initialMajorCooldowns))
}

// Adds a delay to the first usage of all CDs so that armor debuffs have time
// to be applied. MCDs that have a user-specified timing are not delayed.
//
// This function should be called from Agent.Init().
func (mcdm *majorCooldownManager) DelayDPSCooldownsForArmorDebuffs() {
	if !mcdm.character.CurrentTarget.HasAuraWithTag(SunderExposeAuraTag) {
		return
	}

	mcdm.character.Env.RegisterPostFinalizeEffect(func() {
		const delay = time.Second * 10
		for i, _ := range mcdm.initialMajorCooldowns {
			mcd := &mcdm.initialMajorCooldowns[i]
			if len(mcd.timings) == 0 && mcd.Type.Matches(CooldownTypeDPS) {
				mcd.timings = append(mcd.timings, delay)
			}
		}
	})
}

func (mcdm *majorCooldownManager) reset(sim *Simulation) {
	// Need to create all cooldowns before calling ActivationFactory on any of them,
	// so that any cooldown can do lookups on other cooldowns.
	for i, _ := range mcdm.majorCooldowns {
		newMCD := &MajorCooldown{}
		*newMCD = mcdm.initialMajorCooldowns[i]
		mcdm.majorCooldowns[i] = newMCD
	}

	for i, _ := range mcdm.majorCooldowns {
		mcdm.majorCooldowns[i].activate = mcdm.majorCooldowns[i].ActivationFactory(sim)
		if mcdm.majorCooldowns[i].activate == nil {
			panic("Nil cooldown activation returned!")
		}
	}

	// For initial sorting.
	mcdm.UpdateMajorCooldowns()
}

// Registers a major cooldown to the Character, which will be automatically
// used when available.
func (mcdm *majorCooldownManager) AddMajorCooldown(mcd MajorCooldown) {
	if mcdm.character.Env.IsFinalized() {
		panic("Major cooldowns may not be added once finalized!")
	}
	if mcd.Spell == nil {
		panic("Major cooldown must have a Spell!")
	}

	spell := mcd.Spell
	if mcd.ActivationFactory == nil {
		mcd.ActivationFactory = func(sim *Simulation) CooldownActivation {
			return func(sim *Simulation, character *Character) {
				spell.Cast(sim, character.CurrentTarget)
			}
		}
	}

	if mcd.Type.Matches(CooldownTypeSurvival) && mcdm.cooldownConfigs.HpPercentForDefensives != 0 {
		origCanActivate := mcd.CanActivate
		mcd.CanActivate = func(sim *Simulation, character *Character) bool {
			if character.CurrentHealthPercent() > mcdm.cooldownConfigs.HpPercentForDefensives {
				return false
			}

			return origCanActivate == nil || origCanActivate(sim, character)
		}
	}

	if mcd.CanActivate == nil {
		mcd.CanActivate = func(sim *Simulation, character *Character) bool {
			return true
		}
	}
	if mcd.ShouldActivate == nil {
		mcd.ShouldActivate = func(sim *Simulation, character *Character) bool {
			return true
		}
	}

	mcdm.initialMajorCooldowns = append(mcdm.initialMajorCooldowns, mcd)
}

func (mcdm *majorCooldownManager) GetInitialMajorCooldown(actionID ActionID) MajorCooldown {
	for _, mcd := range mcdm.initialMajorCooldowns {
		if mcd.Spell.SameAction(actionID) {
			return mcd
		}
	}

	return MajorCooldown{}
}

func (mcdm *majorCooldownManager) GetMajorCooldown(actionID ActionID) *MajorCooldown {
	for _, mcd := range mcdm.majorCooldowns {
		if mcd.Spell.SameAction(actionID) {
			return mcd
		}
	}

	return nil
}

// Returns all MCDs.
func (mcdm *majorCooldownManager) GetMajorCooldowns() []*MajorCooldown {
	return mcdm.majorCooldowns
}

func (mcdm *majorCooldownManager) GetMajorCooldownIDs() []*proto.ActionID {
	ids := make([]*proto.ActionID, len(mcdm.initialMajorCooldowns))
	for i, mcd := range mcdm.initialMajorCooldowns {
		ids[i] = mcd.Spell.ActionID.ToProto()
	}
	return ids
}

func (mcdm *majorCooldownManager) HasMajorCooldown(actionID ActionID) bool {
	return mcdm.GetMajorCooldown(actionID) != nil
}

func (mcdm *majorCooldownManager) DisableMajorCooldown(actionID ActionID) {
	mcd := mcdm.GetMajorCooldown(actionID)
	if mcd != nil {
		mcd.disabled = true
	}
}

func (mcdm *majorCooldownManager) EnableMajorCooldown(actionID ActionID) {
	mcd := mcdm.GetMajorCooldown(actionID)
	if mcd != nil {
		mcd.disabled = false
	}
}

// Disabled all MCDs that are currently enabled, and returns a list of the MCDs
// which were disabled by this call.
// If cooldownType is not CooldownTypeUnknown, then will be restricted to cooldowns of that type.
func (mcdm *majorCooldownManager) DisableAllEnabledCooldowns(cooldownType CooldownType) []*MajorCooldown {
	disabledMCDs := []*MajorCooldown{}
	for _, mcd := range mcdm.majorCooldowns {
		if mcd.IsEnabled() && (cooldownType == CooldownTypeUnknown || mcd.Type.Matches(cooldownType)) {
			mcdm.DisableMajorCooldown(mcd.Spell.ActionID)
			disabledMCDs = append(disabledMCDs, mcd)
		}
	}
	return disabledMCDs
}

func (mcdm *majorCooldownManager) EnableAllCooldowns(mcdsToEnable []*MajorCooldown) {
	for _, mcd := range mcdsToEnable {
		mcdm.EnableMajorCooldown(mcd.Spell.ActionID)
	}
}

func (mcdm *majorCooldownManager) TryUseCooldowns(sim *Simulation) {
	mcdm.tryUsing = true
	for curIdx := 0; curIdx < len(mcdm.majorCooldowns) && mcdm.majorCooldowns[curIdx].IsReady(sim); curIdx++ {
		mcd := mcdm.majorCooldowns[curIdx]
		if mcd.tryActivateInternal(sim, mcdm.character) {
			if mcdm.sortOne(mcd, curIdx) {
				if mcdm.fullSort {
					// We need to re-sort the whole array
					mcdm.sort()
					mcdm.fullSort = false
					// Reset back to start because new things could be available to activate now.
					curIdx = 0
				} else {
					// This just means the current MCD was sorted further back and now we need to re-check the current idx.
					curIdx--
				}
			}
			if mcd.Spell.DefaultCast.GCD > 0 {
				// If we used a MCD that uses the GCD (like drums), hold off on using
				// any remaining MCDs so they aren't wasted.
				break
			}
		}
	}

	mcdm.tryUsing = false
}

// sortOne will take the given mcd and attempt to sort it towards the back.
// If it finds a linked CD (like trinkets that share offensive CD) it will sort them backwards first.
//  If while sorting it finds something further back with lower CD than the previous one (for example, after activating cold snap)
//  it will mark that the whole slice needs to be re-sorted "mcdm.fullSort" and returns immediately.
func (mcdm *majorCooldownManager) sortOne(mcd *MajorCooldown, curIdx int) bool {
	newReadyAt := mcd.ReadyAt()
	var lastReadAt time.Duration
	for sortIdx := curIdx + 1; sortIdx < len(mcdm.majorCooldowns); sortIdx++ {
		if mcdm.majorCooldowns[sortIdx].Spell.SharedCD.Timer == mcd.Spell.SharedCD.Timer {
			mcdm.sortOne(mcdm.majorCooldowns[sortIdx], sortIdx)
			if mcdm.fullSort {
				return true
			}
		}
		otherReady := mcdm.majorCooldowns[sortIdx].ReadyAt()
		if otherReady < lastReadAt {
			// This means we had some CDs get changed during last activation. We will need a full re-sort.
			mcdm.fullSort = true
			return true
		}
		if otherReady > newReadyAt || (otherReady == newReadyAt && mcdm.majorCooldowns[sortIdx].Priority < mcd.Priority) {
			// This means that this sortIDX is the first spot that is *after* the new ready time.
			// move all CDs before this one forward,
			if sortIdx-1 > curIdx {
				copy(mcdm.majorCooldowns[curIdx:], mcdm.majorCooldowns[curIdx+1:sortIdx])
				mcdm.majorCooldowns[sortIdx-1] = mcd
				return true
			}
			return false
		}
		lastReadAt = otherReady
	}
	// This means it needs to go to the back
	copy(mcdm.majorCooldowns[curIdx:], mcdm.majorCooldowns[curIdx+1:])
	mcdm.majorCooldowns[len(mcdm.majorCooldowns)-1] = mcd
	return true
}

// This function should be called if the CD for a major cooldown changes outside
// of the TryActivate() call.
func (mcdm *majorCooldownManager) UpdateMajorCooldowns() {
	if mcdm.tryUsing {
		panic("Do not call UpdateMajorCooldowns while iterating cooldowns in TryUseCooldowns")
	}
	mcdm.sort()
}

func (mcdm *majorCooldownManager) sort() {
	sort.Slice(mcdm.majorCooldowns, func(i, j int) bool {
		// Since we're just comparing and don't actually care about the remaining CD, ok to use 0 instead of sim.CurrentTime.
		cdA := mcdm.majorCooldowns[i].ReadyAt()
		cdB := mcdm.majorCooldowns[j].ReadyAt()
		return cdA < cdB || (cdA == cdB && mcdm.majorCooldowns[i].Priority > mcdm.majorCooldowns[j].Priority)
	})
}

// Add a major cooldown to the given agent, which provides a temporary boost to a single stat.
// This is use for effects like Icon of the Silver Crescent and Bloodlust Brooch.
func RegisterTemporaryStatsOnUseCD(character *Character, auraLabel string, tempStats stats.Stats, duration time.Duration, config SpellConfig) {
	aura := character.NewTemporaryStatsAura(auraLabel, config.ActionID, tempStats, duration)

	cdType := CooldownTypeUsableShapeShifted
	if tempStats.DotProduct(stats.Stats{stats.Armor: 1, stats.Block: 1, stats.BlockValue: 1, stats.Dodge: 1, stats.Parry: 1, stats.Health: 1}).Equals(stats.Stats{}) {
		cdType |= CooldownTypeDPS
	} else {
		cdType |= CooldownTypeSurvival
	}

	config.Flags |= SpellFlagNoOnCastComplete
	config.ApplyEffects = func(sim *Simulation, _ *Unit, _ *Spell) {
		aura.Activate(sim)
	}
	spell := character.RegisterSpell(config)

	character.AddMajorCooldown(MajorCooldown{
		Spell: spell,
		Type:  cdType,
	})
}

// Helper function to make an ApplyEffect for a temporary stats on-use cooldown.
func MakeTemporaryStatsOnUseCDRegistration(auraLabel string, tempStats stats.Stats, duration time.Duration, config SpellConfig, cdFunc func(*Character) Cooldown, sharedCDFunc func(*Character) Cooldown) ApplyEffect {
	return func(agent Agent) {
		localConfig := config
		character := agent.GetCharacter()
		if cdFunc != nil {
			localConfig.Cast.CD = cdFunc(character)
		}
		if sharedCDFunc != nil {
			localConfig.Cast.SharedCD = sharedCDFunc(character)
		}
		RegisterTemporaryStatsOnUseCD(character, auraLabel, tempStats, duration, localConfig)
	}
}
