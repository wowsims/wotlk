package core

import (
	"time"

	"golang.org/x/exp/slices"

	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
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
	CooldownTypeExplosive
	CooldownTypeSurvival
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

	// Whether the cooldown meets all optional conditions for activation. These
	// conditions will be ignored when the user specifies their own activation time.
	// This is for things like mana thresholds, which are optimizations for better
	// automatic timing.
	ShouldActivate CooldownActivationCondition

	// Fixed timings at which to use this cooldown. If these are specified, they
	// are used instead of ShouldActivate.
	timings []time.Duration

	// Number of times this MCD was used so far in the current iteration.
	numUsages int

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

// Roughly how long until the next cast will happen, accounting for both spell CD and user-specified timings.
func (mcd *MajorCooldown) TimeToNextCast(sim *Simulation) time.Duration {
	timeToReady := mcd.TimeToReady(sim)
	if mcd.numUsages < len(mcd.timings) {
		timeToReady = MaxDuration(timeToReady, mcd.timings[mcd.numUsages]-sim.CurrentTime)
	}
	return timeToReady
}

func (mcd *MajorCooldown) IsEnabled() bool {
	return !mcd.disabled
}

func (mcd *MajorCooldown) Enable() {
	if mcd != nil && mcd.disabled {
		mcd.disabled = false
	}
}

func (mcd *MajorCooldown) Disable() {
	if mcd != nil && !mcd.disabled {
		mcd.disabled = true
	}
}

func (mcd *MajorCooldown) GetTimings() []time.Duration {
	return mcd.timings
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

func (mcd *MajorCooldown) shouldActivateHelper(sim *Simulation, character *Character) bool {
	if mcd.Type.Matches(CooldownTypeSurvival) && character.cooldownConfigs.HpPercentForDefensives != 0 {
		if character.CurrentHealthPercent() > character.cooldownConfigs.HpPercentForDefensives {
			return false
		}
	}

	if !mcd.Spell.CanCast(sim, character.CurrentTarget) {
		return false
	}

	if mcd.numUsages < len(mcd.timings) {
		return sim.CurrentTime >= mcd.timings[mcd.numUsages]
	} else {
		return mcd.ShouldActivate(sim, character)
	}
}

// Activates this MCD, if all the conditions pass.
// Returns whether the MCD was activated.
func (mcd *MajorCooldown) tryActivateHelper(sim *Simulation, character *Character) bool {
	shouldActivate := mcd.shouldActivateHelper(sim, character)

	if shouldActivate {
		if mcd.Spell.Flags.Matches(SpellFlagHelpful) {
			mcd.Spell.Cast(sim, &character.Unit)
		} else {
			mcd.Spell.Cast(sim, character.CurrentTarget)
		}

		mcd.numUsages++
		if sim.Log != nil {
			character.Log(sim, "Major cooldown used: %s", mcd.Spell.ActionID)
		}
	}

	return shouldActivate
}

type cooldownConfigs struct {
	Cooldowns                 []*proto.Cooldown
	HpPercentForDefensives    float64
	DesyncProcTrinket1Seconds int32
	DesyncProcTrinket2Seconds int32
}

type majorCooldownManager struct {
	// The Character whose cooldowns are being managed.
	character *Character

	// User-specified cooldown configs.
	cooldownConfigs cooldownConfigs

	// Cached list of major cooldowns sorted by priority, for resetting quickly.
	initialMajorCooldowns []MajorCooldown

	// Major cooldowns, ordered by next available. This should always contain
	// the same cooldowns as initialMajorCooldowns, but the order will change over
	// the course of the sim.
	majorCooldowns []*MajorCooldown
	minReady       time.Duration

	tryUsing bool
	fullSort bool
}

func newMajorCooldownManager(cooldowns *proto.Cooldowns) majorCooldownManager {
	if cooldowns == nil {
		return majorCooldownManager{}
	}

	cooldownConfigs := cooldownConfigs{
		HpPercentForDefensives:    cooldowns.HpPercentForDefensives,
		DesyncProcTrinket1Seconds: cooldowns.DesyncProcTrinket1Seconds,
		DesyncProcTrinket2Seconds: cooldowns.DesyncProcTrinket2Seconds,
	}
	for _, cooldownConfig := range cooldowns.Cooldowns {
		if cooldownConfig.Id != nil {
			cooldownConfigs.Cooldowns = append(cooldownConfigs.Cooldowns, cooldownConfig)
		}
	}

	return majorCooldownManager{
		cooldownConfigs: cooldownConfigs,
	}
}

func (mcdm *majorCooldownManager) initialize(character *Character) {
	mcdm.character = character
}

func (mcdm *majorCooldownManager) finalize() {
	// Match user-specified cooldown configs to existing cooldowns.
	for i := range mcdm.initialMajorCooldowns {
		mcd := &mcdm.initialMajorCooldowns[i]
		mcd.timings = []time.Duration{}

		for _, cooldownConfig := range mcdm.cooldownConfigs.Cooldowns {
			configID := ProtoToActionID(cooldownConfig.Id)
			if configID.SameAction(mcd.Spell.ActionID) {
				mcd.timings = make([]time.Duration, len(cooldownConfig.Timings))
				for t, timing := range cooldownConfig.Timings {
					mcd.timings[t] = DurationFromSeconds(timing)
				}
				break
			}
		}
	}

	mcdm.majorCooldowns = make([]*MajorCooldown, len(mcdm.initialMajorCooldowns))
}

// Adds a delay to the first usage of all CDs so that debuffs have time
// to be applied. MCDs that have a user-specified timing are not delayed.
//
// This function should be called from Agent.Init().
func (mcdm *majorCooldownManager) DelayDPSCooldownsForArmorDebuffs(delay time.Duration) {
	if mcdm.character.IsUsingAPL {
		return
	}
	mcdm.character.Env.RegisterPostFinalizeEffect(func() {
		for i := range mcdm.initialMajorCooldowns {
			mcd := &mcdm.initialMajorCooldowns[i]
			if len(mcd.timings) == 0 && mcd.Type.Matches(CooldownTypeDPS) && !mcd.Type.Matches(CooldownTypeExplosive) {
				oldShouldActivate := mcd.ShouldActivate
				mcd.ShouldActivate = func(sim *Simulation, character *Character) bool {
					if oldShouldActivate != nil && !oldShouldActivate(sim, character) {
						return false
					}
					return sim.CurrentTime >= delay
				}
			}
		}
	})
}

// Adds a delay to the first usage of all CDs overriding shouldActivate for cooldownTypeDPS,
// MCDs that have a user-specified timing are not delayed.
// This function should be called from Agent.Init().
func (mcdm *majorCooldownManager) DelayDPSCooldowns(delay time.Duration) {
	mcdm.character.Env.RegisterPostFinalizeEffect(func() {
		for i := range mcdm.initialMajorCooldowns {
			mcd := &mcdm.initialMajorCooldowns[i]
			if len(mcd.timings) == 0 && mcd.Type.Matches(CooldownTypeDPS) {
				oldShouldActivate := mcd.ShouldActivate
				mcd.ShouldActivate = func(sim *Simulation, character *Character) bool {
					if oldShouldActivate != nil && !oldShouldActivate(sim, character) {
						return false
					}
					return sim.CurrentTime >= delay
				}
			}
		}
	})
}

func desyncTrinketProcAura(aura *Aura, delay time.Duration) {
	if cb := aura.OnSpellHitDealt; cb != nil {
		aura.OnSpellHitDealt = func(aura *Aura, sim *Simulation, spell *Spell, result *SpellResult) {
			if sim.CurrentTime >= delay {
				cb(aura, sim, spell, result)
			}
		}
	}

	if cb := aura.OnSpellHitTaken; cb != nil {
		aura.OnSpellHitTaken = func(aura *Aura, sim *Simulation, spell *Spell, result *SpellResult) {
			if sim.CurrentTime >= delay {
				cb(aura, sim, spell, result)
			}
		}
	}

	if cb := aura.OnPeriodicDamageDealt; cb != nil {
		aura.OnPeriodicDamageDealt = func(aura *Aura, sim *Simulation, spell *Spell, result *SpellResult) {
			if sim.CurrentTime >= delay {
				cb(aura, sim, spell, result)
			}
		}
	}

	if cb := aura.OnHealDealt; cb != nil {
		aura.OnHealDealt = func(aura *Aura, sim *Simulation, spell *Spell, result *SpellResult) {
			if sim.CurrentTime >= delay {
				cb(aura, sim, spell, result)
			}
		}
	}

	if cb := aura.OnPeriodicHealDealt; cb != nil {
		aura.OnPeriodicHealDealt = func(aura *Aura, sim *Simulation, spell *Spell, result *SpellResult) {
			if sim.CurrentTime >= delay {
				cb(aura, sim, spell, result)
			}
		}
	}

	if cb := aura.OnCastComplete; cb != nil {
		aura.OnCastComplete = func(aura *Aura, sim *Simulation, spell *Spell) {
			if sim.CurrentTime >= delay {
				cb(aura, sim, spell)
			}
		}
	}
}

func findTrinketAura(character *Character, trinketID int32) *Aura {
	for _, aura := range character.auras {
		if aura.ActionIDForProc.ItemID == trinketID {
			return aura
		}
	}
	return nil
}

// Desyncs trinket procs per configured user settings.
// Hold the first proc back until some time into the simulation (i.e. because the player
// un-equipped and re-equipped the trinket before pull).
func (mcdm *majorCooldownManager) DesyncTrinketProcs() {
	if delay := time.Duration(mcdm.cooldownConfigs.DesyncProcTrinket1Seconds) * time.Second; delay > 0 {
		if trinket1 := mcdm.character.Equip[ItemSlotTrinket1]; trinket1.ID > 0 && HasItemEffect(trinket1.ID) {
			mcdm.character.Env.RegisterPostFinalizeEffect(func() {
				if aura := findTrinketAura(mcdm.character, trinket1.ID); aura != nil {
					desyncTrinketProcAura(aura, delay)
				}
			})
		}
	}

	if delay := time.Duration(mcdm.cooldownConfigs.DesyncProcTrinket2Seconds) * time.Second; delay > 0 {
		if trinket2 := mcdm.character.Equip[ItemSlotTrinket2]; trinket2.ID > 0 && HasItemEffect(trinket2.ID) {
			mcdm.character.Env.RegisterPostFinalizeEffect(func() {
				if aura := findTrinketAura(mcdm.character, trinket2.ID); aura != nil {
					desyncTrinketProcAura(aura, delay)
				}
			})
		}
	}
}

func (mcdm *majorCooldownManager) reset(sim *Simulation) {
	for i := range mcdm.majorCooldowns {
		newMCD := &MajorCooldown{}
		*newMCD = mcdm.initialMajorCooldowns[i]
		mcdm.majorCooldowns[i] = newMCD
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
	mcd.Spell.Flags |= SpellFlagAPL | SpellFlagMCD

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

func (mcdm *majorCooldownManager) removeInitialMajorCooldown(actionID ActionID) {
	for i, mcd := range mcdm.initialMajorCooldowns {
		if mcd.Spell.SameAction(actionID) {
			mcdm.initialMajorCooldowns = append(mcdm.initialMajorCooldowns[:i], mcdm.initialMajorCooldowns[i+1:]...)
			mcdm.majorCooldowns = mcdm.majorCooldowns[:len(mcdm.majorCooldowns)-1]
			return
		}
	}
}

func (mcdm *majorCooldownManager) GetMajorCooldown(actionID ActionID) *MajorCooldown {
	for _, mcd := range mcdm.majorCooldowns {
		if mcd.Spell.SameAction(actionID) {
			return mcd
		}
	}
	return nil
}
func (mcdm *majorCooldownManager) GetMajorCooldownIgnoreTag(actionID ActionID) *MajorCooldown {
	for _, mcd := range mcdm.majorCooldowns {
		if mcd.Spell.SameActionIgnoreTag(actionID) {
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

func (mcdm *majorCooldownManager) getFirstReadyMCD(sim *Simulation) *MajorCooldown {
	if sim.CurrentTime < mcdm.minReady {
		return nil
	}

	for _, mcd := range mcdm.majorCooldowns {
		if !mcd.IsReady(sim) {
			return nil
		}
		if mcd.shouldActivateHelper(sim, mcdm.character) {
			return mcd
		}
	}

	return nil
}

func (mcdm *majorCooldownManager) TryUseCooldowns(sim *Simulation) {
	if sim.CurrentTime < mcdm.minReady {
		return
	}

	mcdm.tryUsing = true
restart:
	for _, mcd := range mcdm.majorCooldowns {
		if !mcd.IsReady(sim) {
			break
		}

		if mcd.tryActivateInternal(sim, mcdm.character) {
			if mcd.IsReady(sim) {
				continue // activation failed
			}
			mcdm.sort()

			if mcd.Spell.DefaultCast.GCD > 0 {
				// If the GCD was used, don't use any more MCDs until the next cycle so
				// their durations aren't partially wasted.
				break
			}

			// many MCDs are off the GCD, so it makes sense to continue
			goto restart
		}
	}
	mcdm.tryUsing = false

	mcdm.minReady = mcdm.majorCooldowns[0].ReadyAt()
}

// This function should be called if the CD for a major cooldown changes outside
// of the TryActivate() call.
func (mcdm *majorCooldownManager) UpdateMajorCooldowns() {
	if mcdm.tryUsing {
		panic("Do not call UpdateMajorCooldowns while iterating cooldowns in TryUseCooldowns")
	}
	if len(mcdm.majorCooldowns) == 0 {
		mcdm.minReady = NeverExpires
		return
	}
	mcdm.sort()
	mcdm.minReady = mcdm.majorCooldowns[0].ReadyAt()
}

func (mcdm *majorCooldownManager) sort() {
	slices.SortStableFunc(mcdm.majorCooldowns, func(m1, m2 *MajorCooldown) bool {
		// Since we're just comparing and don't actually care about the remaining CD, ok to use 0 instead of sim.CurrentTime.
		return m1.ReadyAt() < m2.ReadyAt() || (m1.ReadyAt() == m2.ReadyAt() && m1.Priority > m2.Priority)
	})
}

// Add a major cooldown to the given agent, which provides a temporary boost to a single stat.
// This is use for effects like Icon of the Silver Crescent and Bloodlust Brooch.
func RegisterTemporaryStatsOnUseCD(character *Character, auraLabel string, tempStats stats.Stats, duration time.Duration, config SpellConfig) {
	aura := character.NewTemporaryStatsAura(auraLabel, config.ActionID, tempStats, duration)

	cdType := CooldownTypeUnknown
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
