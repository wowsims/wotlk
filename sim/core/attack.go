package core

import (
	"slices"
	"time"

	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

// ReplaceMHSwing is called right before a main hand auto attack fires.
// It must never return nil, but either a replacement spell or the passed in regular mhSwingSpell.
// This allows for abilities that convert a white attack into a yellow attack.
type ReplaceMHSwing func(sim *Simulation, mhSwingSpell *Spell) *Spell

// Represents a generic weapon. Pets / unarmed / various other cases don't use
// actual weapon items so this is an abstraction of a Weapon.
type Weapon struct {
	BaseDamageMin        float64
	BaseDamageMax        float64
	AttackPowerPerDPS    float64
	SwingSpeed           float64
	NormalizedSwingSpeed float64
	CritMultiplier       float64
	SpellSchool          SpellSchool
}

func (weapon *Weapon) DPS() float64 {
	if weapon.SwingSpeed == 0 {
		return 0
	} else {
		return (weapon.BaseDamageMin + weapon.BaseDamageMax) / 2.0 / weapon.SwingSpeed
	}
}

func newWeaponFromUnarmed(critMultiplier float64) Weapon {
	// These numbers are probably wrong but nobody cares.
	return Weapon{
		BaseDamageMin:        0,
		BaseDamageMax:        0,
		SwingSpeed:           1,
		NormalizedSwingSpeed: 1,
		CritMultiplier:       critMultiplier,
		AttackPowerPerDPS:    DefaultAttackPowerPerDPS,
	}
}

func newWeaponFromItem(item *Item, critMultiplier float64, bonusDps float64) Weapon {
	normalizedWeaponSpeed := 2.4
	if item.WeaponType == proto.WeaponType_WeaponTypeDagger {
		normalizedWeaponSpeed = 1.7
	} else if item.HandType == proto.HandType_HandTypeTwoHand {
		normalizedWeaponSpeed = 3.3
	} else if item.RangedWeaponType != proto.RangedWeaponType_RangedWeaponTypeUnknown {
		normalizedWeaponSpeed = 2.8
	}

	return Weapon{
		BaseDamageMin:        item.WeaponDamageMin + bonusDps*item.SwingSpeed,
		BaseDamageMax:        item.WeaponDamageMax + bonusDps*item.SwingSpeed,
		SwingSpeed:           item.SwingSpeed,
		NormalizedSwingSpeed: normalizedWeaponSpeed,
		CritMultiplier:       critMultiplier,
		AttackPowerPerDPS:    DefaultAttackPowerPerDPS,
	}
}

// Returns weapon stats using the main hand equipped weapon.
func (character *Character) WeaponFromMainHand(critMultiplier float64) Weapon {
	if weapon := character.GetMHWeapon(); weapon != nil {
		return newWeaponFromItem(weapon, critMultiplier, character.PseudoStats.BonusMHDps)
	} else {
		return newWeaponFromUnarmed(critMultiplier)
	}
}

// Returns weapon stats using the off-hand equipped weapon.
func (character *Character) WeaponFromOffHand(critMultiplier float64) Weapon {
	if weapon := character.GetOHWeapon(); weapon != nil {
		return newWeaponFromItem(weapon, critMultiplier, character.PseudoStats.BonusOHDps)
	} else {
		return Weapon{}
	}
}

// Returns weapon stats using the ranged equipped weapon.
func (character *Character) WeaponFromRanged(critMultiplier float64) Weapon {
	if weapon := character.GetRangedWeapon(); weapon != nil {
		return newWeaponFromItem(weapon, critMultiplier, character.PseudoStats.BonusRangedDps)
	} else {
		return Weapon{}
	}
}

func (weapon *Weapon) GetSpellSchool() SpellSchool {
	if weapon.SpellSchool == SpellSchoolNone {
		return SpellSchoolPhysical
	} else {
		return weapon.SpellSchool
	}
}

func (weapon *Weapon) EnemyWeaponDamage(sim *Simulation, attackPower float64, damageSpread float64) float64 {
	// Maximum damage range is 133% of minimum damage; AP contribution is % of minimum damage roll.
	// Patchwerk follows special damage range rules.
	// TODO: Scrape more logs to determine these values more accurately. AP defined in constants.go

	rand := 1 + damageSpread*sim.RandomFloat("Enemy Weapon Damage")

	return weapon.BaseDamageMin * (rand + attackPower*EnemyAutoAttackAPCoefficient)
}

func (weapon *Weapon) BaseDamage(sim *Simulation) float64 {
	return weapon.BaseDamageMin + (weapon.BaseDamageMax-weapon.BaseDamageMin)*sim.RandomFloat("Weapon Base Damage")
}

func (weapon *Weapon) AverageDamage() float64 {
	return (weapon.BaseDamageMin + weapon.BaseDamageMax) / 2
}

func (weapon *Weapon) CalculateWeaponDamage(sim *Simulation, attackPower float64) float64 {
	return weapon.BaseDamage(sim) + (weapon.SwingSpeed*attackPower)/weapon.AttackPowerPerDPS
}

func (weapon *Weapon) CalculateAverageWeaponDamage(attackPower float64) float64 {
	return weapon.AverageDamage() + (weapon.SwingSpeed*attackPower)/weapon.AttackPowerPerDPS
}

func (weapon *Weapon) CalculateNormalizedWeaponDamage(sim *Simulation, attackPower float64) float64 {
	return weapon.BaseDamage(sim) + (weapon.NormalizedSwingSpeed*attackPower)/weapon.AttackPowerPerDPS
}

func (unit *Unit) MHWeaponDamage(sim *Simulation, attackPower float64) float64 {
	return unit.AutoAttacks.MH.CalculateWeaponDamage(sim, attackPower)
}
func (unit *Unit) MHNormalizedWeaponDamage(sim *Simulation, attackPower float64) float64 {
	return unit.AutoAttacks.MH.CalculateNormalizedWeaponDamage(sim, attackPower)
}

func (unit *Unit) OHWeaponDamage(sim *Simulation, attackPower float64) float64 {
	return 0.5 * unit.AutoAttacks.OH.CalculateWeaponDamage(sim, attackPower)
}
func (unit *Unit) OHNormalizedWeaponDamage(sim *Simulation, attackPower float64) float64 {
	return 0.5 * unit.AutoAttacks.OH.CalculateNormalizedWeaponDamage(sim, attackPower)
}

func (unit *Unit) RangedWeaponDamage(sim *Simulation, attackPower float64) float64 {
	return unit.AutoAttacks.Ranged.CalculateWeaponDamage(sim, attackPower)
}

type MeleeDamageCalculator func(attackPower float64, bonusWeaponDamage float64) float64

// Returns whether this hit effect is associated with the main-hand weapon.
func (spell *Spell) IsMH() bool {
	return spell.ProcMask.Matches(ProcMaskMeleeMH)
}

// Returns whether this hit effect is associated with the off-hand weapon.
func (spell *Spell) IsOH() bool {
	return spell.ProcMask.Matches(ProcMaskMeleeOH)
}

// Returns whether this hit effect is associated with either melee weapon.
func (spell *Spell) IsMelee() bool {
	return spell.ProcMask.Matches(ProcMaskMelee)
}

type AutoAttacks struct {
	agent  Agent
	unit   *Unit
	MH     Weapon
	OH     Weapon
	Ranged Weapon

	IsDualWielding bool

	// If true, core engine will handle calling SwingMelee(). Set to false to manually manage
	// swings, for example for hunter melee weaving.
	AutoSwingMelee bool

	// If true, core engine will handle calling SwingRanged(). Unless you're a hunter, don't
	// use this.
	AutoSwingRanged bool

	MainhandSwingAt time.Duration
	OffhandSwingAt  time.Duration
	RangedSwingAt   time.Duration

	MHConfig     SpellConfig
	OHConfig     SpellConfig
	RangedConfig SpellConfig

	MHAuto     *Spell
	OHAuto     *Spell
	RangedAuto *Spell

	ReplaceMHSwing ReplaceMHSwing

	// Current melee and ranged swing speeds, and corresponding swing durations, updated in UpdateSwingTimers.
	curMeleeSpeed      float64
	curMHSwingDuration time.Duration
	curOHSwingDuration time.Duration

	curRangedSpeed         float64
	curRangedSwingDuration time.Duration

	// PendingAction which handles auto attacks.
	autoSwingAction    *PendingAction
	autoSwingCancelled bool
}

// Options for initializing auto attacks.
type AutoAttackOptions struct {
	MainHand        Weapon
	OffHand         Weapon
	Ranged          Weapon
	AutoSwingMelee  bool // If true, core engine will handle calling SwingMelee() for you.
	AutoSwingRanged bool // If true, core engine will handle calling SwingMelee() for you.
	SyncType        int32
	ReplaceMHSwing  ReplaceMHSwing
}

func (unit *Unit) EnableAutoAttacks(agent Agent, options AutoAttackOptions) {
	if options.MainHand.AttackPowerPerDPS == 0 {
		options.MainHand.AttackPowerPerDPS = DefaultAttackPowerPerDPS
	}
	if options.OffHand.AttackPowerPerDPS == 0 {
		options.OffHand.AttackPowerPerDPS = DefaultAttackPowerPerDPS
	}
	unit.AutoAttacks = AutoAttacks{
		agent:           agent,
		unit:            unit,
		MH:              options.MainHand,
		OH:              options.OffHand,
		Ranged:          options.Ranged,
		AutoSwingMelee:  options.AutoSwingMelee,
		AutoSwingRanged: options.AutoSwingRanged,
		ReplaceMHSwing:  options.ReplaceMHSwing,
		IsDualWielding:  options.MainHand.SwingSpeed != 0 && options.OffHand.SwingSpeed != 0,
	}

	unit.AutoAttacks.MHConfig = SpellConfig{
		ActionID:    ActionID{OtherID: proto.OtherAction_OtherActionAttack, Tag: 1},
		SpellSchool: unit.AutoAttacks.MH.GetSpellSchool(),
		ProcMask:    ProcMaskMeleeMHAuto,
		Flags:       SpellFlagMeleeMetrics | SpellFlagIncludeTargetBonusDamage | SpellFlagNoOnCastComplete,

		DamageMultiplier: 1,
		CritMultiplier:   options.MainHand.CritMultiplier,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *Simulation, target *Unit, spell *Spell) {
			baseDamage := spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower()) +
				spell.BonusWeaponDamage()

			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWhite)
		},
	}

	unit.AutoAttacks.OHConfig = SpellConfig{
		ActionID:    ActionID{OtherID: proto.OtherAction_OtherActionAttack, Tag: 2},
		SpellSchool: unit.AutoAttacks.OH.GetSpellSchool(),
		ProcMask:    ProcMaskMeleeOHAuto,
		Flags:       SpellFlagMeleeMetrics | SpellFlagIncludeTargetBonusDamage | SpellFlagNoOnCastComplete,

		DamageMultiplier: 1,
		CritMultiplier:   options.OffHand.CritMultiplier,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *Simulation, target *Unit, spell *Spell) {
			baseDamage := spell.Unit.OHWeaponDamage(sim, spell.MeleeAttackPower()) +
				spell.BonusWeaponDamage()

			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWhite)
		},
	}

	unit.AutoAttacks.RangedConfig = SpellConfig{
		ActionID:    ActionID{OtherID: proto.OtherAction_OtherActionShoot},
		SpellSchool: SpellSchoolPhysical,
		ProcMask:    ProcMaskRangedAuto,
		Flags:       SpellFlagMeleeMetrics | SpellFlagIncludeTargetBonusDamage,

		DamageMultiplier: 1,
		CritMultiplier:   options.Ranged.CritMultiplier,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *Simulation, target *Unit, spell *Spell) {
			baseDamage := spell.Unit.RangedWeaponDamage(sim, spell.RangedAttackPower(target)) +
				spell.BonusWeaponDamage()
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeRangedHitAndCrit)
		},
	}

	if unit.Type == EnemyUnit {
		unit.AutoAttacks.MHConfig.ApplyEffects = func(sim *Simulation, target *Unit, spell *Spell) {
			ap := max(0, spell.Unit.stats[stats.AttackPower])
			baseDamage := spell.Unit.AutoAttacks.MH.EnemyWeaponDamage(sim, ap, spell.Unit.PseudoStats.DamageSpread)

			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeEnemyMeleeWhite)
		}
		unit.AutoAttacks.OHConfig.ApplyEffects = func(sim *Simulation, target *Unit, spell *Spell) {
			ap := max(0, spell.Unit.stats[stats.AttackPower])
			baseDamage := spell.Unit.AutoAttacks.MH.EnemyWeaponDamage(sim, ap, spell.Unit.PseudoStats.DamageSpread) * 0.5

			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeEnemyMeleeWhite)
		}
	}

	// Will be un-cancelled in Reset(), this is just to prevent any swing logic
	// from being triggered during initialization.
	unit.AutoAttacks.autoSwingCancelled = true
}

// Empty handler so Agents don't have to provide one if they have no logic to add.
func (unit *Unit) OnAutoAttack(_ *Simulation, _ *Spell) {}

func (aa *AutoAttacks) finalize() {
	if aa.AutoSwingMelee {
		aa.MHAuto = aa.unit.GetOrRegisterSpell(aa.MHConfig)
		aa.OHAuto = aa.unit.GetOrRegisterSpell(aa.OHConfig)
	}

	if aa.AutoSwingRanged {
		aa.RangedAuto = aa.unit.GetOrRegisterSpell(aa.RangedConfig)
	}
}

func (aa *AutoAttacks) reset(sim *Simulation) {
	if !aa.AutoSwingMelee && !aa.AutoSwingRanged {
		return
	}

	if aa.AutoSwingMelee {
		aa.curMeleeSpeed = aa.unit.SwingSpeed()
		aa.UpdateMeleeDurations()

		aa.MainhandSwingAt = 0
		aa.OffhandSwingAt = 0

		// Apply random delay of 0 - 50% swing time, to one of the weapons if dual wielding
		if aa.IsDualWielding {
			if aa.unit.Type == EnemyUnit {
				aa.OffhandSwingAt = DurationFromSeconds(aa.MH.SwingSpeed / 2)
			} else {
				if sim.RandomFloat("SwingResetWeapon") < 0.5 {
					aa.MainhandSwingAt = DurationFromSeconds(sim.RandomFloat("SwingResetDelay") * aa.MH.SwingSpeed / 2)
				} else {
					aa.OffhandSwingAt = DurationFromSeconds(sim.RandomFloat("SwingResetDelay") * aa.MH.SwingSpeed / 2)
				}
			}
		}
	}

	if aa.AutoSwingRanged {
		aa.curRangedSpeed = aa.unit.RangedSwingSpeed()
		aa.UpdateRangedDuration()

		aa.RangedSwingAt = 0
	}

	aa.autoSwingAction = nil
	aa.autoSwingCancelled = false
}

func (aa *AutoAttacks) startPull(sim *Simulation) {
	if aa.autoSwingCancelled {
		return
	}

	if aa.AutoSwingMelee {
		aa.rescheduleMelee(sim)
	}

	if aa.AutoSwingRanged {
		aa.rescheduleRanged(sim)
	}
}

func (aa *AutoAttacks) rescheduleRanged(sim *Simulation) {
	if aa.autoSwingAction != nil {
		aa.autoSwingAction.Cancel(sim)
	}

	if aa.autoSwingAction != nil {
		aa.autoSwingAction.Cancel(sim)
	}

	var pa *PendingAction

	pa = &PendingAction{
		NextActionAt: aa.RangedSwingAt,
		Priority:     ActionPriorityAuto,
		OnAction: func(sim *Simulation) {
			aa.SwingRanged(sim, aa.unit.CurrentTarget)
			pa.NextActionAt = aa.RangedSwingAt

			// Cancelled means we made a new one because of a swing speed change.
			if !pa.cancelled {
				sim.AddPendingAction(pa)
			}
		},
	}

	aa.autoSwingAction = pa
	sim.AddPendingAction(pa)
}

func (aa *AutoAttacks) rescheduleMelee(sim *Simulation) {
	if aa.autoSwingAction != nil {
		aa.autoSwingAction.Cancel(sim)
	}

	var pa *PendingAction

	pa = &PendingAction{
		NextActionAt: aa.NextAttackAt(),
		Priority:     ActionPriorityAuto,
		OnAction: func(sim *Simulation) {
			aa.SwingMelee(sim, aa.unit.CurrentTarget)
			pa.NextActionAt = aa.NextAttackAt()

			// Cancelled means we made a new one because of a swing speed change.
			if !pa.cancelled {
				sim.AddPendingAction(pa)
			}
		},
	}

	aa.autoSwingAction = pa
	sim.AddPendingAction(pa)
}

// Stops the auto swing action for the rest of the iteration. Used for pets
// after being disabled.
func (aa *AutoAttacks) CancelAutoSwing(sim *Simulation) {
	if aa.autoSwingAction != nil {
		aa.autoSwingAction.Cancel(sim)
		aa.autoSwingAction = nil
	}
	aa.autoSwingCancelled = true
}

// Re-enables the auto swing action for the iteration
func (aa *AutoAttacks) EnableAutoSwing(sim *Simulation) {
	// Already enabled so nothing to do
	if !aa.autoSwingCancelled {
		return
	}

	aa.autoSwingCancelled = false

	if aa.AutoSwingMelee {
		if aa.MainhandSwingAt < sim.CurrentTime {
			aa.MainhandSwingAt = sim.CurrentTime
		}
		if aa.OffhandSwingAt < sim.CurrentTime {
			aa.OffhandSwingAt = sim.CurrentTime
		}

		aa.rescheduleMelee(sim)
	}

	if aa.AutoSwingRanged {
		if aa.RangedSwingAt < sim.CurrentTime {
			aa.RangedSwingAt = sim.CurrentTime
		}

		aa.rescheduleRanged(sim)
	}
}

// The amount of time between two MH swings.
func (aa *AutoAttacks) MainhandSwingSpeed() time.Duration {
	return aa.curMHSwingDuration
}

// The amount of time between two OH swings.
func (aa *AutoAttacks) OffhandSwingSpeed() time.Duration {
	return aa.curOHSwingDuration
}

// The amount of time between two Ranged swings.
func (aa *AutoAttacks) RangedSwingSpeed() time.Duration {
	return aa.curRangedSwingDuration
}

// SwingMelee will check any swing timers if they are up, and if so, swing!
func (aa *AutoAttacks) SwingMelee(sim *Simulation, target *Unit) {
	aa.TrySwingMH(sim, target)
	aa.TrySwingOH(sim, target)
}

func (aa *AutoAttacks) SwingRanged(sim *Simulation, target *Unit) {
	aa.TrySwingRanged(sim, target)
}

// Performs an auto attack using the main hand weapon, if the MH CD is ready.
func (aa *AutoAttacks) TrySwingMH(sim *Simulation, target *Unit) {
	if aa.MainhandSwingAt > sim.CurrentTime {
		return
	}

	if aa.unit.IsUsingAPL {
		// Need to check APL here to allow last-moment HS queue casts.
		aa.unit.Rotation.DoNextAction(sim)
	}

	attackSpell := aa.MaybeReplaceMHSwing(sim, aa.MHAuto)

	// Update swing timer BEFORE the cast, so that APL checks for TimeToNextAuto behave correctly
	// if the attack causes APL evaluations (e.g. from rage gain).
	aa.MainhandSwingAt = sim.CurrentTime + aa.curMHSwingDuration
	attackSpell.Cast(sim, target)

	if !sim.Options.Interactive {
		if aa.unit.IsUsingAPL {
			aa.unit.Rotation.DoNextAction(sim)
		} else {
			aa.agent.OnAutoAttack(sim, attackSpell)
		}
	}
}

// Optionally replaces the given swing spell with an Agent-specified MH Swing replacer.
// This is for effects like Heroic Strike or Raptor Strike.
func (aa *AutoAttacks) MaybeReplaceMHSwing(sim *Simulation, mhSwingSpell *Spell) *Spell {
	if aa.ReplaceMHSwing == nil {
		return mhSwingSpell
	}

	// Allow MH swing to be overridden for abilities like Heroic Strike.
	return aa.ReplaceMHSwing(sim, mhSwingSpell)
}

// Performs an auto attack using the main hand weapon, if the OH CD is ready.
func (aa *AutoAttacks) TrySwingOH(sim *Simulation, target *Unit) {
	if !aa.IsDualWielding || aa.OffhandSwingAt > sim.CurrentTime {
		return
	}

	// Update swing timer BEFORE the cast, so that APL checks for TimeToNextAuto behave correctly
	// if the attack causes APL evaluations (e.g. from rage gain).
	aa.OffhandSwingAt = sim.CurrentTime + aa.curOHSwingDuration
	aa.OHAuto.Cast(sim, target)

	if !sim.Options.Interactive {
		if aa.unit.IsUsingAPL {
			aa.unit.Rotation.DoNextAction(sim)
		} else {
			aa.agent.OnAutoAttack(sim, aa.OHAuto)
		}
	}
}

// Performs an auto attack using the ranged weapon, if the Ranged CD is ready.
func (aa *AutoAttacks) TrySwingRanged(sim *Simulation, target *Unit) {
	if aa.RangedSwingAt > sim.CurrentTime {
		return
	}

	// Update swing timer BEFORE the cast, so that APL checks for TimeToNextAuto behave correctly
	// if the attack causes APL evaluations.
	aa.RangedSwingAt = sim.CurrentTime + aa.RangedSwingSpeed()
	aa.RangedAuto.Cast(sim, target)

	if !sim.Options.Interactive {
		if aa.unit.IsUsingAPL {
			aa.unit.Rotation.DoNextAction(sim)
		} else {
			aa.agent.OnAutoAttack(sim, aa.RangedAuto)
		}
	}
}

// This is used internally, and for druid shifts (where a weapon changes without resetting the swing timer).
func (aa *AutoAttacks) UpdateMeleeDurations() {
	aa.curMHSwingDuration = DurationFromSeconds(aa.MH.SwingSpeed / aa.curMeleeSpeed)
	if aa.IsDualWielding {
		aa.curOHSwingDuration = DurationFromSeconds(aa.OH.SwingSpeed / aa.curMeleeSpeed)
	}
}

func (aa *AutoAttacks) UpdateRangedDuration() {
	aa.curRangedSwingDuration = DurationFromSeconds(aa.Ranged.SwingSpeed / aa.curRangedSpeed)
}

func (aa *AutoAttacks) UpdateSwingTimers(sim *Simulation) {
	if aa.AutoSwingRanged {
		aa.curRangedSpeed = aa.unit.RangedSwingSpeed()
		aa.UpdateRangedDuration()
		// ranged attack speed changes aren't applied mid-"swing"
	}

	if aa.AutoSwingMelee {
		oldSwingSpeed := aa.curMeleeSpeed

		aa.curMeleeSpeed = aa.unit.SwingSpeed()
		aa.UpdateMeleeDurations()

		f := oldSwingSpeed / aa.curMeleeSpeed

		if remainingSwingTime := aa.MainhandSwingAt - sim.CurrentTime; remainingSwingTime > 0 {
			aa.MainhandSwingAt = sim.CurrentTime + time.Duration(float64(remainingSwingTime)*f)
		}

		if aa.IsDualWielding {
			if remainingSwingTime := aa.OffhandSwingAt - sim.CurrentTime; remainingSwingTime > 0 {
				aa.OffhandSwingAt = sim.CurrentTime + time.Duration(float64(remainingSwingTime)*f)
			}
		}

		if aa.autoSwingCancelled {
			return
		}

		if sim.CurrentTime < 0 {
			return
		}

		aa.rescheduleMelee(sim)
	}
}

// StopMeleeUntil should be used whenever a non-melee spell is cast. It stops melee, then restarts it
// at end of cast, but with a reset swing timer (as if swings had just landed).
func (aa *AutoAttacks) StopMeleeUntil(sim *Simulation, readyAt time.Duration, desyncOH bool) {
	if !aa.AutoSwingMelee { // if not auto swinging, don't auto restart.
		return
	}

	aa.CancelAutoSwing(sim)

	// schedule restart action
	sim.AddPendingAction(&PendingAction{
		NextActionAt: readyAt,
		Priority:     ActionPriorityAuto,
		OnAction: func(sim *Simulation) {
			aa.restartMelee(sim, desyncOH)
		},
	})
}

func (aa *AutoAttacks) restartMelee(sim *Simulation, desyncOH bool) {
	if !aa.autoSwingCancelled {
		return
	}

	aa.autoSwingCancelled = false

	aa.MainhandSwingAt = sim.CurrentTime + aa.curMHSwingDuration
	if aa.IsDualWielding {
		aa.OffhandSwingAt = sim.CurrentTime + aa.curOHSwingDuration
		if desyncOH {
			// Used by warrior to desync offhand after unglyphed Shattering Throw.
			aa.OffhandSwingAt += aa.curOHSwingDuration / 2
		}
	}

	aa.rescheduleMelee(sim)
}

// Delays all swing timers for the specified amount. Only used by Slam.
func (aa *AutoAttacks) DelayMeleeBy(sim *Simulation, delay time.Duration) {
	if delay <= 0 {
		return
	}

	aa.MainhandSwingAt += delay
	if aa.IsDualWielding {
		aa.OffhandSwingAt += delay
	}

	aa.rescheduleMelee(sim)
}

func (aa *AutoAttacks) DelayRangedUntil(sim *Simulation, readyAt time.Duration) {
	if readyAt <= aa.RangedSwingAt {
		return
	}

	aa.RangedSwingAt = readyAt

	aa.rescheduleRanged(sim)
}

// Returns the time at which the next attack will occur.
func (aa *AutoAttacks) NextAttackAt() time.Duration {
	if aa.IsDualWielding && aa.OffhandSwingAt < aa.MainhandSwingAt {
		return aa.OffhandSwingAt
	} else {
		return aa.MainhandSwingAt
	}
}

type PPMManager struct {
	procMasks   []ProcMask
	procChances []float64
}

// Returns whether the effect procced.
func (ppmm *PPMManager) Proc(sim *Simulation, procMask ProcMask, label string) bool {
	for i, m := range ppmm.procMasks {
		if m.Matches(procMask) {
			return sim.RandomFloat(label) < ppmm.procChances[i]
		}
	}
	return false
}

func (ppmm *PPMManager) Chance(procMask ProcMask) float64 {
	for i, m := range ppmm.procMasks {
		if m.Matches(procMask) {
			return ppmm.procChances[i]
		}
	}
	return 0
}

func (aa *AutoAttacks) NewPPMManager(ppm float64, procMask ProcMask) PPMManager {
	if !aa.AutoSwingMelee && !aa.AutoSwingRanged {
		return PPMManager{}
	}

	ppmm := PPMManager{procMasks: make([]ProcMask, 0, 2), procChances: make([]float64, 0, 2)}

	mergeOrAppend := func(speed float64, mask ProcMask) {
		if speed == 0 || mask == 0 {
			return
		}

		if i := slices.Index(ppmm.procChances, speed); i != -1 {
			ppmm.procMasks[i] |= mask
			return
		}

		ppmm.procMasks = append(ppmm.procMasks, mask)
		ppmm.procChances = append(ppmm.procChances, speed)
	}

	mergeOrAppend(aa.MH.SwingSpeed, procMask&^ProcMaskRanged&^ProcMaskMeleeOH) // "everything else", even if not explicitly flagged MH
	mergeOrAppend(aa.OH.SwingSpeed, procMask&ProcMaskMeleeOH)
	mergeOrAppend(aa.Ranged.SwingSpeed, procMask&ProcMaskRanged)

	for i := range ppmm.procChances {
		ppmm.procChances[i] *= ppm / 60
	}

	return ppmm
}

// Returns whether a PPM-based effect procced.
// Using NewPPMManager() is preferred; this function should only be used when
// the attacker is not known at initialization time.
func (aa *AutoAttacks) PPMProc(sim *Simulation, ppm float64, procMask ProcMask, label string, spell *Spell) bool {
	if !aa.AutoSwingMelee && !aa.AutoSwingRanged {
		return false
	}

	switch {
	case spell.ProcMask.Matches(procMask &^ ProcMaskMeleeOH &^ ProcMaskRanged):
		return sim.RandomFloat(label) < ppm*aa.MH.SwingSpeed/60.0
	case spell.ProcMask.Matches(procMask & ProcMaskMeleeOH):
		return sim.RandomFloat(label) < ppm*aa.OH.SwingSpeed/60.0
	case spell.ProcMask.Matches(procMask & ProcMaskRanged):
		return sim.RandomFloat(label) < ppm*aa.Ranged.SwingSpeed/60.0
	}
	return false
}

func (unit *Unit) applyParryHaste() {
	if !unit.PseudoStats.ParryHaste || !unit.AutoAttacks.AutoSwingMelee {
		return
	}

	unit.RegisterAura(Aura{
		Label:    "Parry Haste",
		Duration: NeverExpires,
		OnReset: func(aura *Aura, sim *Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitTaken: func(aura *Aura, sim *Simulation, spell *Spell, result *SpellResult) {
			if !result.Outcome.Matches(OutcomeParry) {
				return
			}

			remainingTime := aura.Unit.AutoAttacks.MainhandSwingAt - sim.CurrentTime
			swingSpeed := aura.Unit.AutoAttacks.curMHSwingDuration
			minRemainingTime := time.Duration(float64(swingSpeed) * 0.2) // 20% of Swing Speed
			defaultReduction := minRemainingTime * 2                     // 40% of Swing Speed

			if remainingTime <= minRemainingTime {
				return
			}

			parryHasteReduction := min(defaultReduction, remainingTime-minRemainingTime)
			newReadyAt := aura.Unit.AutoAttacks.MainhandSwingAt - parryHasteReduction
			if sim.Log != nil {
				aura.Unit.Log(sim, "MH Swing reduced by %s due to parry haste, will now occur at %s", parryHasteReduction, newReadyAt)
			}

			aura.Unit.AutoAttacks.MainhandSwingAt = newReadyAt
			aura.Unit.AutoAttacks.rescheduleMelee(sim)
		},
	})
}
