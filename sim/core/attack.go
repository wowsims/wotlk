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
	return unit.AutoAttacks.mh.CalculateWeaponDamage(sim, attackPower)
}
func (unit *Unit) MHNormalizedWeaponDamage(sim *Simulation, attackPower float64) float64 {
	return unit.AutoAttacks.mh.CalculateNormalizedWeaponDamage(sim, attackPower)
}

func (unit *Unit) OHWeaponDamage(sim *Simulation, attackPower float64) float64 {
	return 0.5 * unit.AutoAttacks.oh.CalculateWeaponDamage(sim, attackPower)
}
func (unit *Unit) OHNormalizedWeaponDamage(sim *Simulation, attackPower float64) float64 {
	return 0.5 * unit.AutoAttacks.oh.CalculateNormalizedWeaponDamage(sim, attackPower)
}

func (unit *Unit) RangedWeaponDamage(sim *Simulation, attackPower float64) float64 {
	return unit.AutoAttacks.ranged.CalculateWeaponDamage(sim, attackPower)
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

func (aa *AutoAttacks) DualWielding() bool {
	return aa.isDualWielding
}

func (aa *AutoAttacks) MH() *Weapon {
	return &aa.mh
}

func (aa *AutoAttacks) SetMH(weapon Weapon) {
	aa.mh = weapon
	aa.mhAuto.CritMultiplier = weapon.CritMultiplier
	aa.curMHSwingDuration = DurationFromSeconds(aa.mh.SwingSpeed / aa.curMeleeSpeed)
}

func (aa *AutoAttacks) OH() *Weapon {
	return &aa.oh
}

func (aa *AutoAttacks) SetOH(weapon Weapon) {
	aa.oh = weapon
	aa.ohAuto.CritMultiplier = weapon.CritMultiplier
	aa.curOHSwingDuration = DurationFromSeconds(aa.oh.SwingSpeed / aa.curMeleeSpeed)
}

func (aa *AutoAttacks) Ranged() *Weapon {
	return &aa.ranged
}

func (aa *AutoAttacks) SetRanged(weapon Weapon) {
	aa.ranged = weapon
	aa.rangedAuto.CritMultiplier = weapon.CritMultiplier
	aa.curRangedSwingDuration = DurationFromSeconds(aa.ranged.SwingSpeed / aa.curRangedSpeed)
}

func (aa *AutoAttacks) AutoSwingMelee() bool {
	return aa.autoSwingMelee
}

func (aa *AutoAttacks) AutoSwingRanged() bool {
	return aa.autoSwingRanged
}

func (aa *AutoAttacks) MHAuto() *Spell {
	return aa.mhAuto
}

func (aa *AutoAttacks) OHAuto() *Spell {
	return aa.ohAuto
}

func (aa *AutoAttacks) RangedAuto() *Spell {
	return aa.rangedAuto
}

func (aa *AutoAttacks) OffhandSwingAt() time.Duration {
	return aa.offhandSwingAt
}

func (aa *AutoAttacks) SetOffhandSwingAt(offhandSwingAt time.Duration) {
	aa.offhandSwingAt = offhandSwingAt
}

type AutoAttacks struct {
	agent Agent
	unit  *Unit

	mh     Weapon
	oh     Weapon
	ranged Weapon

	isDualWielding bool

	// If true, core engine will handle calling SwingMelee(). Set to false to manually manage
	// swings, for example for hunter melee weaving.
	autoSwingMelee bool

	// If true, core engine will handle calling SwingRanged(). Unless you're a hunter, don't
	// use this.
	autoSwingRanged bool

	mainhandSwingAt time.Duration
	offhandSwingAt  time.Duration
	rangedSwingAt   time.Duration

	// These are created in EnableAutoAttacks, and can be safely altered before finalize(), where the related spells are created
	MHConfig     SpellConfig
	OHConfig     SpellConfig
	RangedConfig SpellConfig

	mhAuto     *Spell
	ohAuto     *Spell
	rangedAuto *Spell

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
	AutoSwingRanged bool // If true, core engine will handle calling SwingRanged() for you.
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
		mh:              options.MainHand,
		oh:              options.OffHand,
		ranged:          options.Ranged,
		autoSwingMelee:  options.AutoSwingMelee,
		autoSwingRanged: options.AutoSwingRanged,
		ReplaceMHSwing:  options.ReplaceMHSwing,
		isDualWielding:  options.MainHand.SwingSpeed != 0 && options.OffHand.SwingSpeed != 0,
	}

	unit.AutoAttacks.MHConfig = SpellConfig{
		ActionID:    ActionID{OtherID: proto.OtherAction_OtherActionAttack, Tag: 1},
		SpellSchool: unit.AutoAttacks.mh.GetSpellSchool(),
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
		SpellSchool: unit.AutoAttacks.oh.GetSpellSchool(),
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
			baseDamage := spell.Unit.AutoAttacks.mh.EnemyWeaponDamage(sim, ap, spell.Unit.PseudoStats.DamageSpread)

			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeEnemyMeleeWhite)
		}
		unit.AutoAttacks.OHConfig.ApplyEffects = func(sim *Simulation, target *Unit, spell *Spell) {
			ap := max(0, spell.Unit.stats[stats.AttackPower])
			baseDamage := spell.Unit.AutoAttacks.mh.EnemyWeaponDamage(sim, ap, spell.Unit.PseudoStats.DamageSpread) * 0.5

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
	if aa.autoSwingMelee {
		aa.mhAuto = aa.unit.GetOrRegisterSpell(aa.MHConfig)
		aa.ohAuto = aa.unit.GetOrRegisterSpell(aa.OHConfig)
	}

	if aa.autoSwingRanged {
		aa.rangedAuto = aa.unit.GetOrRegisterSpell(aa.RangedConfig)
	}
}

func (aa *AutoAttacks) reset(sim *Simulation) {
	if !aa.autoSwingMelee && !aa.autoSwingRanged {
		return
	}

	if aa.autoSwingMelee {
		aa.curMeleeSpeed = aa.unit.SwingSpeed()
		aa.updateMeleeDurations()

		aa.mainhandSwingAt = 0
		aa.offhandSwingAt = 0

		// Apply random delay of 0 - 50% swing time, to one of the weapons if dual wielding
		if aa.isDualWielding {
			if aa.unit.Type == EnemyUnit {
				aa.offhandSwingAt = DurationFromSeconds(aa.mh.SwingSpeed / 2)
			} else {
				if sim.RandomFloat("SwingResetWeapon") < 0.5 {
					aa.mainhandSwingAt = DurationFromSeconds(sim.RandomFloat("SwingResetDelay") * aa.mh.SwingSpeed / 2)
				} else {
					aa.offhandSwingAt = DurationFromSeconds(sim.RandomFloat("SwingResetDelay") * aa.mh.SwingSpeed / 2)
				}
			}
		}
	}

	if aa.autoSwingRanged {
		aa.curRangedSpeed = aa.unit.RangedSwingSpeed()
		aa.UpdateRangedDuration()

		aa.rangedSwingAt = 0
	}

	aa.autoSwingAction = nil
	aa.autoSwingCancelled = false
}

func (aa *AutoAttacks) startPull(sim *Simulation) {
	if aa.autoSwingCancelled {
		return
	}

	if aa.autoSwingMelee {
		aa.rescheduleMelee(sim)
	}

	if aa.autoSwingRanged {
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
		NextActionAt: aa.rangedSwingAt,
		Priority:     ActionPriorityAuto,
		OnAction: func(sim *Simulation) {
			aa.SwingRanged(sim, aa.unit.CurrentTarget)
			pa.NextActionAt = aa.rangedSwingAt

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

	if aa.autoSwingMelee {
		if aa.mainhandSwingAt < sim.CurrentTime {
			aa.mainhandSwingAt = sim.CurrentTime
		}
		if aa.offhandSwingAt < sim.CurrentTime {
			aa.offhandSwingAt = sim.CurrentTime
		}

		aa.rescheduleMelee(sim)
	}

	if aa.autoSwingRanged {
		if aa.rangedSwingAt < sim.CurrentTime {
			aa.rangedSwingAt = sim.CurrentTime
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
	if aa.mainhandSwingAt > sim.CurrentTime {
		return
	}

	attackSpell := aa.mhAuto

	if aa.ReplaceMHSwing != nil {
		if aa.unit.IsUsingAPL {
			// Need to check APL here to allow last-moment HS queue casts.
			aa.unit.Rotation.DoNextAction(sim)
		}
		// Allow MH swing to be overridden for abilities like Heroic Strike.
		attackSpell = aa.ReplaceMHSwing(sim, aa.mhAuto)
	}

	// Update swing timer BEFORE the cast, so that APL checks for TimeToNextAuto behave correctly
	// if the attack causes APL evaluations (e.g. from rage gain).
	aa.mainhandSwingAt = sim.CurrentTime + aa.curMHSwingDuration
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
	if !aa.isDualWielding || aa.offhandSwingAt > sim.CurrentTime {
		return
	}

	// Update swing timer BEFORE the cast, so that APL checks for TimeToNextAuto behave correctly
	// if the attack causes APL evaluations (e.g. from rage gain).
	aa.offhandSwingAt = sim.CurrentTime + aa.curOHSwingDuration
	aa.ohAuto.Cast(sim, target)

	if !sim.Options.Interactive {
		if aa.unit.IsUsingAPL {
			aa.unit.Rotation.DoNextAction(sim)
		} else {
			aa.agent.OnAutoAttack(sim, aa.ohAuto)
		}
	}
}

// Performs an auto attack using the ranged weapon, if the Ranged CD is ready.
func (aa *AutoAttacks) TrySwingRanged(sim *Simulation, target *Unit) {
	if aa.rangedSwingAt > sim.CurrentTime {
		return
	}

	// Update swing timer BEFORE the cast, so that APL checks for TimeToNextAuto behave correctly
	// if the attack causes APL evaluations.
	aa.rangedSwingAt = sim.CurrentTime + aa.RangedSwingSpeed()
	aa.rangedAuto.Cast(sim, target)

	if !sim.Options.Interactive {
		if aa.unit.IsUsingAPL {
			aa.unit.Rotation.DoNextAction(sim)
		} else {
			aa.agent.OnAutoAttack(sim, aa.rangedAuto)
		}
	}
}

func (aa *AutoAttacks) updateMeleeDurations() {
	aa.curMHSwingDuration = DurationFromSeconds(aa.mh.SwingSpeed / aa.curMeleeSpeed)
	if aa.isDualWielding {
		aa.curOHSwingDuration = DurationFromSeconds(aa.oh.SwingSpeed / aa.curMeleeSpeed)
	}
}

func (aa *AutoAttacks) UpdateRangedDuration() {
	aa.curRangedSwingDuration = DurationFromSeconds(aa.ranged.SwingSpeed / aa.curRangedSpeed)
}

func (aa *AutoAttacks) UpdateSwingTimers(sim *Simulation) {
	if aa.autoSwingRanged {
		aa.curRangedSpeed = aa.unit.RangedSwingSpeed()
		aa.UpdateRangedDuration()
		// ranged attack speed changes aren't applied mid-"swing"
	}

	if aa.autoSwingMelee {
		oldSwingSpeed := aa.curMeleeSpeed

		aa.curMeleeSpeed = aa.unit.SwingSpeed()
		aa.updateMeleeDurations()

		f := oldSwingSpeed / aa.curMeleeSpeed

		if remainingSwingTime := aa.mainhandSwingAt - sim.CurrentTime; remainingSwingTime > 0 {
			aa.mainhandSwingAt = sim.CurrentTime + time.Duration(float64(remainingSwingTime)*f)
		}

		if aa.isDualWielding {
			if remainingSwingTime := aa.offhandSwingAt - sim.CurrentTime; remainingSwingTime > 0 {
				aa.offhandSwingAt = sim.CurrentTime + time.Duration(float64(remainingSwingTime)*f)
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
	if !aa.autoSwingMelee { // if not auto swinging, don't auto restart.
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

	aa.mainhandSwingAt = sim.CurrentTime + aa.curMHSwingDuration
	if aa.isDualWielding {
		aa.offhandSwingAt = sim.CurrentTime + aa.curOHSwingDuration
		if desyncOH {
			// Used by warrior to desync offhand after unglyphed Shattering Throw.
			aa.offhandSwingAt += aa.curOHSwingDuration / 2
		}
	}

	aa.rescheduleMelee(sim)
}

// Delays all swing timers for the specified amount. Only used by Slam.
func (aa *AutoAttacks) DelayMeleeBy(sim *Simulation, delay time.Duration) {
	if delay <= 0 {
		return
	}

	aa.mainhandSwingAt += delay
	if aa.isDualWielding {
		aa.offhandSwingAt += delay
	}

	aa.rescheduleMelee(sim)
}

func (aa *AutoAttacks) DelayRangedUntil(sim *Simulation, readyAt time.Duration) {
	if readyAt <= aa.rangedSwingAt {
		return
	}

	aa.rangedSwingAt = readyAt

	aa.rescheduleRanged(sim)
}

// Returns the time at which the next attack will occur.
func (aa *AutoAttacks) NextAttackAt() time.Duration {
	if aa.isDualWielding && aa.offhandSwingAt < aa.mainhandSwingAt {
		return aa.offhandSwingAt
	} else {
		return aa.mainhandSwingAt
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
	if !aa.autoSwingMelee && !aa.autoSwingRanged {
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

	mergeOrAppend(aa.mh.SwingSpeed, procMask&^ProcMaskRanged&^ProcMaskMeleeOH) // "everything else", even if not explicitly flagged MH
	mergeOrAppend(aa.oh.SwingSpeed, procMask&ProcMaskMeleeOH)
	mergeOrAppend(aa.ranged.SwingSpeed, procMask&ProcMaskRanged)

	for i := range ppmm.procChances {
		ppmm.procChances[i] *= ppm / 60
	}

	return ppmm
}

// Returns whether a PPM-based effect procced.
// Using NewPPMManager() is preferred; this function should only be used when
// the attacker is not known at initialization time.
func (aa *AutoAttacks) PPMProc(sim *Simulation, ppm float64, procMask ProcMask, label string, spell *Spell) bool {
	if !aa.autoSwingMelee && !aa.autoSwingRanged {
		return false
	}

	switch {
	case spell.ProcMask.Matches(procMask &^ ProcMaskMeleeOH &^ ProcMaskRanged):
		return sim.RandomFloat(label) < ppm*aa.mh.SwingSpeed/60.0
	case spell.ProcMask.Matches(procMask & ProcMaskMeleeOH):
		return sim.RandomFloat(label) < ppm*aa.oh.SwingSpeed/60.0
	case spell.ProcMask.Matches(procMask & ProcMaskRanged):
		return sim.RandomFloat(label) < ppm*aa.ranged.SwingSpeed/60.0
	}
	return false
}

func (unit *Unit) applyParryHaste() {
	if !unit.PseudoStats.ParryHaste || !unit.AutoAttacks.autoSwingMelee {
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

			remainingTime := aura.Unit.AutoAttacks.mainhandSwingAt - sim.CurrentTime
			swingSpeed := aura.Unit.AutoAttacks.curMHSwingDuration
			minRemainingTime := time.Duration(float64(swingSpeed) * 0.2) // 20% of Swing Speed
			defaultReduction := minRemainingTime * 2                     // 40% of Swing Speed

			if remainingTime <= minRemainingTime {
				return
			}

			parryHasteReduction := min(defaultReduction, remainingTime-minRemainingTime)
			newReadyAt := aura.Unit.AutoAttacks.mainhandSwingAt - parryHasteReduction
			if sim.Log != nil {
				aura.Unit.Log(sim, "MH Swing reduced by %s due to parry haste, will now occur at %s", parryHasteReduction, newReadyAt)
			}

			aura.Unit.AutoAttacks.mainhandSwingAt = newReadyAt
			aura.Unit.AutoAttacks.rescheduleMelee(sim)
		},
	})
}
