package core

import (
	"time"

	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

// ReplaceMHSwing is called right before an auto attack fires
//
//	If it returns nil, the attack takes place as normal. If it returns a Spell,
//	that Spell is used in place of the attack.
//	This allows for abilities that convert a white attack into yellow attack.
type ReplaceMHSwing func(sim *Simulation, mhSwingSpell *Spell) *Spell

// Represents a generic weapon. Pets / unarmed / various other cases dont use
// actual weapon items so this is an abstraction of a Weapon.
type Weapon struct {
	BaseDamageMin              float64
	BaseDamageMax              float64
	MeleeAttackRatingPerDamage float64
	SwingSpeed                 float64
	NormalizedSwingSpeed       float64
	SwingDuration              time.Duration // Duration between 2 swings.
	CritMultiplier             float64
	SpellSchool                SpellSchool
}

func (w Weapon) DPS() float64 {
	if w.SwingSpeed == 0 {
		return 0
	} else {
		return (w.BaseDamageMin + w.BaseDamageMax) / 2.0 / w.SwingSpeed
	}
}

func (w Weapon) WithBonusDPS(bonusDps float64) Weapon {
	newWeapon := w
	bonusSwingDamage := bonusDps * w.SwingSpeed
	newWeapon.BaseDamageMin += bonusSwingDamage
	newWeapon.BaseDamageMax += bonusSwingDamage
	return newWeapon
}

func newWeaponFromUnarmed(critMultiplier float64) Weapon {
	// These numbers are probably wrong but nobody cares.
	return Weapon{
		BaseDamageMin:              0,
		BaseDamageMax:              0,
		SwingSpeed:                 1,
		NormalizedSwingSpeed:       1,
		SwingDuration:              time.Second,
		CritMultiplier:             critMultiplier,
		MeleeAttackRatingPerDamage: MeleeAttackRatingPerDamage,
	}
}

func newWeaponFromItem(item Item, critMultiplier float64) Weapon {
	normalizedWeaponSpeed := 2.4
	if item.WeaponType == proto.WeaponType_WeaponTypeDagger {
		normalizedWeaponSpeed = 1.7
	} else if item.HandType == proto.HandType_HandTypeTwoHand {
		normalizedWeaponSpeed = 3.3
	} else if item.RangedWeaponType != proto.RangedWeaponType_RangedWeaponTypeUnknown {
		normalizedWeaponSpeed = 2.8
	}

	return Weapon{
		BaseDamageMin:              item.WeaponDamageMin,
		BaseDamageMax:              item.WeaponDamageMax,
		SwingSpeed:                 item.SwingSpeed,
		NormalizedSwingSpeed:       normalizedWeaponSpeed,
		SwingDuration:              time.Duration(item.SwingSpeed * float64(time.Second)),
		CritMultiplier:             critMultiplier,
		MeleeAttackRatingPerDamage: MeleeAttackRatingPerDamage,
	}
}

// Returns weapon stats using the main hand equipped weapon.
func (character *Character) WeaponFromMainHand(critMultiplier float64) Weapon {
	if weapon := character.GetMHWeapon(); weapon != nil {
		return newWeaponFromItem(*weapon, critMultiplier).WithBonusDPS(character.PseudoStats.BonusMHDps)
	} else {
		return newWeaponFromUnarmed(critMultiplier).WithBonusDPS(character.PseudoStats.BonusMHDps)
	}
}

// Returns weapon stats using the off hand equipped weapon.
func (character *Character) WeaponFromOffHand(critMultiplier float64) Weapon {
	if weapon := character.GetOHWeapon(); weapon != nil {
		return newWeaponFromItem(*weapon, critMultiplier).WithBonusDPS(character.PseudoStats.BonusOHDps)
	} else {
		return Weapon{}
	}
}

// Returns weapon stats using the ranged equipped weapon.
func (character *Character) WeaponFromRanged(critMultiplier float64) Weapon {
	if weapon := character.GetRangedWeapon(); weapon != nil {
		return newWeaponFromItem(*weapon, critMultiplier).WithBonusDPS(character.PseudoStats.BonusRangedDps)
	} else {
		return Weapon{}
	}
}

func (weapon Weapon) GetSpellSchool() SpellSchool {
	if weapon.SpellSchool == SpellSchoolNone {
		return SpellSchoolPhysical
	} else {
		return weapon.SpellSchool
	}
}

func (weapon Weapon) EnemyWeaponDamage(sim *Simulation, attackPower float64, damageSpread float64) float64 {
	// Maximum damage range is 133% of minimum damage; AP contribution is % of minimum damage roll
	// Patchwerk follows special damage range rules.
	// TODO: Scrape more logs to determine these values more accurately. AP defined in constants.go

	rand := 1 + damageSpread*sim.RandomFloat("Enemy Weapon Damage")

	return weapon.BaseDamageMin * (rand + attackPower*EnemyAutoAttackAPCoefficient)
}

func (weapon Weapon) BaseDamage(sim *Simulation) float64 {
	return weapon.BaseDamageMin + (weapon.BaseDamageMax-weapon.BaseDamageMin)*sim.RandomFloat("Weapon Base Damage")
}

func (weapon Weapon) AverageDamage() float64 {
	return (weapon.BaseDamageMin + weapon.BaseDamageMax) / 2
}

func (weapon Weapon) CalculateWeaponDamage(sim *Simulation, attackPower float64) float64 {
	return weapon.BaseDamage(sim) + (weapon.SwingSpeed*attackPower)/weapon.MeleeAttackRatingPerDamage
}

func (weapon Weapon) CalculateAverageWeaponDamage(attackPower float64) float64 {
	return weapon.AverageDamage() + (weapon.SwingSpeed*attackPower)/weapon.MeleeAttackRatingPerDamage
}

func (weapon Weapon) CalculateNormalizedWeaponDamage(sim *Simulation, attackPower float64) float64 {
	return weapon.BaseDamage(sim) + (weapon.NormalizedSwingSpeed*attackPower)/weapon.MeleeAttackRatingPerDamage
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

	// Set this to 1 to sync your auto attacks together, or 2 to use the OH delay macro, mostly used by enhance shamans.
	// This will intentionally perfectly sync or delay OH swings to that they always fall within the
	// 0.5s window following a MH swing.
	SyncType int32

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

	// The time at which the last MH swing occurred.
	previousMHSwingAt time.Duration
	PreviousSwingAt   time.Duration

	// Current melee swing speed, based on haste stat and melee swing multiplier pseudostat.
	curSwingSpeed float64

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
	if options.MainHand.MeleeAttackRatingPerDamage == 0 {
		options.MainHand.MeleeAttackRatingPerDamage = MeleeAttackRatingPerDamage
	}
	if options.OffHand.MeleeAttackRatingPerDamage == 0 {
		options.OffHand.MeleeAttackRatingPerDamage = MeleeAttackRatingPerDamage
	}
	unit.AutoAttacks = AutoAttacks{
		agent:           agent,
		unit:            unit,
		MH:              options.MainHand,
		OH:              options.OffHand,
		Ranged:          options.Ranged,
		AutoSwingMelee:  options.AutoSwingMelee,
		AutoSwingRanged: options.AutoSwingRanged,
		SyncType:        options.SyncType,
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
			ap := MaxFloat(0, spell.Unit.stats[stats.AttackPower])
			baseDamage := spell.Unit.AutoAttacks.MH.EnemyWeaponDamage(sim, ap, spell.Unit.PseudoStats.DamageSpread)

			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeEnemyMeleeWhite)
		}
		unit.AutoAttacks.OHConfig.ApplyEffects = func(sim *Simulation, target *Unit, spell *Spell) {
			ap := MaxFloat(0, spell.Unit.stats[stats.AttackPower])
			baseDamage := spell.Unit.AutoAttacks.MH.EnemyWeaponDamage(sim, ap, spell.Unit.PseudoStats.DamageSpread) * 0.5

			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeEnemyMeleeWhite)
		}
	}

	// Will be un-cancelled in Reset(), this is just to prevent any swing logic
	// from being triggered during initialization.
	unit.AutoAttacks.autoSwingCancelled = true
}

func (aa *AutoAttacks) IsEnabled() bool {
	return aa.MH.SwingSpeed != 0
}

// Empty handler so Agents don't have to provide one if they have no logic to add.
func (unit *Unit) OnAutoAttack(sim *Simulation, spell *Spell) {}

func (aa *AutoAttacks) finalize() {
	if !aa.IsEnabled() {
		return
	}

	aa.MHAuto = aa.unit.GetOrRegisterSpell(aa.MHConfig)
	aa.OHAuto = aa.unit.GetOrRegisterSpell(aa.OHConfig)

	if aa.RangedConfig.ProcMask != ProcMaskUnknown {
		aa.RangedAuto = aa.unit.GetOrRegisterSpell(aa.RangedConfig)
	}
}

func (aa *AutoAttacks) reset(sim *Simulation) {
	if !aa.IsEnabled() {
		return
	}

	aa.curSwingSpeed = aa.unit.SwingSpeed()

	aa.MainhandSwingAt = 0
	aa.OffhandSwingAt = 0
	aa.RangedSwingAt = 0
	aa.PreviousSwingAt = 0

	// Apply random delay of 0 - 50% swing time, to one of the weapons if dual wielding
	if aa.IsDualWielding {
		// Set a fake value for previousMHSwing so that offhand swing delay works
		// properly at the start.
		aa.previousMHSwingAt = time.Second * -1

		var delay time.Duration
		var isMHDelay bool
		if aa.unit.Type == EnemyUnit {
			delay = time.Duration(float64(aa.MH.SwingDuration / 2))
			isMHDelay = false
		} else {
			delay = time.Duration(sim.RandomFloat("SwingResetDelay") * float64(aa.MH.SwingDuration/2))
			isMHDelay = sim.RandomFloat("SwingResetWeapon") < 0.5
		}

		if isMHDelay {
			aa.MainhandSwingAt = delay
		} else {
			aa.OffhandSwingAt = delay
		}
	}

	aa.autoSwingAction = nil
	aa.autoSwingCancelled = false
}

func (aa *AutoAttacks) startPull(sim *Simulation) {
	if aa.IsEnabled() && aa.unit.IsEnabled() {
		aa.resetAutoSwing(sim)
	}
}

func (aa *AutoAttacks) resetAutoSwing(sim *Simulation) {
	if aa.autoSwingCancelled || (!aa.AutoSwingMelee && !aa.AutoSwingRanged) || sim.CurrentTime < 0 {
		return
	}

	if aa.autoSwingAction != nil {
		aa.autoSwingAction.Cancel(sim)
	}

	pa := &PendingAction{
		NextActionAt: TernaryDuration(aa.AutoSwingMelee, aa.NextAttackAt(), aa.RangedSwingAt),
		Priority:     ActionPriorityAuto,
	}

	if aa.AutoSwingMelee {
		pa.OnAction = func(sim *Simulation) {
			aa.SwingMelee(sim, aa.unit.CurrentTarget)
			pa.NextActionAt = aa.NextAttackAt()

			// Cancelled means we made a new one because of a swing speed change.
			if !pa.cancelled {
				sim.AddPendingAction(pa)
			}
		}
	} else { // Ranged
		pa.OnAction = func(sim *Simulation) {
			aa.SwingRanged(sim, aa.unit.CurrentTarget)
			pa.NextActionAt = aa.RangedSwingAt

			// Cancelled means we made a new one because of a swing speed change.
			if !pa.cancelled {
				sim.AddPendingAction(pa)
			}
		}
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

// Renables the auto swing action for the iteration
func (aa *AutoAttacks) EnableAutoSwing(sim *Simulation) {
	// Already enabled so nothing to do
	if aa.autoSwingAction != nil {
		return
	}
	if sim.CurrentTime < 0 {
		return
	}

	if aa.MainhandSwingAt < sim.CurrentTime {
		aa.MainhandSwingAt = sim.CurrentTime
	}
	if aa.OffhandSwingAt < sim.CurrentTime {
		aa.OffhandSwingAt = sim.CurrentTime
	}
	if aa.RangedSwingAt < sim.CurrentTime {
		aa.RangedSwingAt = sim.CurrentTime
	}

	aa.autoSwingCancelled = false
	aa.resetAutoSwing(sim)
}

// The amount of time between two MH swings.
func (aa *AutoAttacks) MainhandSwingSpeed() time.Duration {
	return time.Duration(float64(aa.MH.SwingDuration) / aa.unit.SwingSpeed())
}

// The amount of time between two OH swings.
func (aa *AutoAttacks) OffhandSwingSpeed() time.Duration {
	return time.Duration(float64(aa.OH.SwingDuration) / aa.unit.SwingSpeed())
}

// The amount of time between two ranged swings.
func (aa *AutoAttacks) RangedSwingSpeed() time.Duration {
	return time.Duration(float64(aa.Ranged.SwingDuration) / aa.unit.RangedSwingSpeed())
}

// SwingMelee will check any swing timers if they are up, and if so, swing!
func (aa *AutoAttacks) SwingMelee(sim *Simulation, target *Unit) {
	aa.TrySwingMH(sim, target)
	aa.TrySwingOH(sim, target)
}

func (aa *AutoAttacks) SwingRanged(sim *Simulation, target *Unit) {
	aa.TrySwingRanged(sim, target)
}

// Performs an autoattack using the main hand weapon, if the MH CD is ready.
func (aa *AutoAttacks) TrySwingMH(sim *Simulation, target *Unit) {
	if aa.MainhandSwingAt > sim.CurrentTime {
		return
	}

	attackSpell := aa.MaybeReplaceMHSwing(sim, aa.MHAuto)

	attackSpell.Cast(sim, target)
	aa.MainhandSwingAt = sim.CurrentTime + aa.MainhandSwingSpeed()
	aa.previousMHSwingAt = sim.CurrentTime
	aa.PreviousSwingAt = sim.CurrentTime
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
	replacementSpell := aa.ReplaceMHSwing(sim, mhSwingSpell)
	if replacementSpell == nil {
		return mhSwingSpell
	} else {
		return replacementSpell
	}
}

// Performs an autoattack using the main hand weapon, if the OH CD is ready.
func (aa *AutoAttacks) TrySwingOH(sim *Simulation, target *Unit) {
	if !aa.IsDualWielding || aa.OffhandSwingAt > sim.CurrentTime {
		return
	}

	if (aa.SyncType == 1) && (sim.CurrentTime-aa.previousMHSwingAt) > time.Millisecond*500 {
		// Perfectly Sync MH and OH attacks
		aa.OffhandSwingAt = aa.MainhandSwingAt
		if sim.Log != nil {
			aa.unit.Log(sim, "Resyncing Weapons")
		}
		return
	} else if (aa.SyncType == 2) && (sim.CurrentTime-aa.previousMHSwingAt) > time.Millisecond*500 {
		// Delay the OH swing for later, so it follows the MH swing.
		aa.OffhandSwingAt = aa.MainhandSwingAt + time.Millisecond*100
		if sim.Log != nil {
			aa.unit.Log(sim, "Delaying OH swing by %s", aa.OffhandSwingAt-sim.CurrentTime)
		}
		return
	}

	aa.OHAuto.Cast(sim, target)
	aa.OffhandSwingAt = sim.CurrentTime + aa.OffhandSwingSpeed()
	aa.PreviousSwingAt = sim.CurrentTime
	if !sim.Options.Interactive {
		if aa.unit.IsUsingAPL {
			aa.unit.Rotation.DoNextAction(sim)
		} else {
			aa.agent.OnAutoAttack(sim, aa.OHAuto)
		}
	}
}

// Performs an autoattack using the ranged weapon, if the ranged CD is ready.
func (aa *AutoAttacks) TrySwingRanged(sim *Simulation, target *Unit) {
	if aa.RangedSwingAt > sim.CurrentTime {
		return
	}

	aa.RangedAuto.Cast(sim, target)
	aa.RangedSwingAt = sim.CurrentTime + aa.RangedSwingSpeed()
	aa.PreviousSwingAt = sim.CurrentTime
	if !sim.Options.Interactive {
		if aa.unit.IsUsingAPL {
			aa.unit.Rotation.DoNextAction(sim)
		} else {
			aa.agent.OnAutoAttack(sim, aa.RangedAuto)
		}
	}
}

func (aa *AutoAttacks) UpdateSwingTime(sim *Simulation) {
	if !aa.IsEnabled() || aa.AutoSwingRanged {
		return
	}

	oldSwingSpeed := aa.curSwingSpeed
	aa.curSwingSpeed = aa.unit.SwingSpeed()

	f := oldSwingSpeed / aa.curSwingSpeed

	remainingSwingTime := aa.MainhandSwingAt - sim.CurrentTime
	if remainingSwingTime > 0 {
		aa.MainhandSwingAt = sim.CurrentTime + time.Duration(float64(remainingSwingTime)*f)
	}

	if aa.IsDualWielding {
		remainingSwingTime := aa.OffhandSwingAt - sim.CurrentTime
		if remainingSwingTime > 0 {
			aa.OffhandSwingAt = sim.CurrentTime + time.Duration(float64(remainingSwingTime)*f)
		}
	}

	aa.resetAutoSwing(sim)
}

// StopMeleeUntil should be used whenever a non-melee spell is cast. It stops melee, then restarts it
// at end of cast, but with a reset swing timer (as if swings had just landed).
func (aa *AutoAttacks) StopMeleeUntil(sim *Simulation, readyAt time.Duration, desyncOH bool) {
	if !aa.AutoSwingMelee { // if not auto swinging, don't auto restart.
		return
	}
	aa.CancelAutoSwing(sim)

	// Used by warrior to desync offhand after Shattering Throw.
	if desyncOH {
		// schedule restart action
		sim.AddPendingAction(&PendingAction{
			NextActionAt: readyAt,
			Priority:     ActionPriorityAuto,
			OnAction:     aa.desyncedRestartMelee,
		})
	} else {
		// schedule restart action
		sim.AddPendingAction(&PendingAction{
			NextActionAt: readyAt,
			Priority:     ActionPriorityAuto,
			OnAction:     aa.restartMelee,
		})
	}

}

func (aa *AutoAttacks) restartMelee(sim *Simulation) {
	if !aa.autoSwingCancelled {
		return
	}

	aa.MainhandSwingAt = sim.CurrentTime + aa.MainhandSwingSpeed()
	if aa.IsDualWielding {
		aa.OffhandSwingAt = sim.CurrentTime + aa.OffhandSwingSpeed()
	}
	aa.autoSwingCancelled = false
	aa.resetAutoSwing(sim)
}

// Emulating how desyncing OH works in the game.
// After swing timer has passed half the swing time, Offhand swing timer will be reset.
func (aa *AutoAttacks) desyncedRestartMelee(sim *Simulation) {
	if !aa.autoSwingCancelled {
		return
	}

	aa.MainhandSwingAt = sim.CurrentTime + aa.MainhandSwingSpeed()
	if aa.IsDualWielding {
		aa.OffhandSwingAt = sim.CurrentTime + aa.OffhandSwingSpeed() + aa.OffhandSwingSpeed()/2
	}
	aa.autoSwingCancelled = false
	aa.resetAutoSwing(sim)
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

	aa.resetAutoSwing(sim)
}

func (aa *AutoAttacks) DelayRangedUntil(sim *Simulation, readyAt time.Duration) {
	if readyAt > aa.RangedSwingAt {
		aa.RangedSwingAt = readyAt
		aa.resetAutoSwing(sim)
	}
}

// Returns the time at which the next attack will occur.
func (aa *AutoAttacks) NextAttackAt() time.Duration {
	nextAttack := aa.MainhandSwingAt
	if aa.OH.SwingSpeed != 0 {
		nextAttack = MinDuration(nextAttack, aa.OffhandSwingAt)
	}
	return nextAttack
}

type PPMManager struct {
	mhProcChance     float64
	ohProcChance     float64
	rangedProcChance float64
	procMask         ProcMask
}

// Returns whether the effect procced.
func (ppmm *PPMManager) Proc(sim *Simulation, procMask ProcMask, label string) bool {
	// Without this procs that can proc only from white attacks
	// are still procing from specials
	if !procMask.Matches(ppmm.procMask) {
		return false
	}

	chance := ppmm.Chance(procMask)
	return chance > 0 && sim.RandomFloat(label) < chance
}

func (ppmm *PPMManager) Chance(procMask ProcMask) float64 {
	if procMask.Matches(ProcMaskMeleeMH) {
		return ppmm.mhProcChance
	} else if procMask.Matches(ProcMaskMeleeOH) {
		return ppmm.ohProcChance
	} else if procMask.Matches(ProcMaskRanged) {
		return ppmm.rangedProcChance
	} else if procMask.Matches(ppmm.procMask) {
		return ppmm.mhProcChance // probably a 'proc from proc' so use main hand.
	}

	return 0
}

func (aa *AutoAttacks) NewPPMManager(ppm float64, procMask ProcMask) PPMManager {
	if !aa.IsEnabled() {
		return PPMManager{}
	}

	ppmm := PPMManager{}
	ppmm.procMask = procMask
	if procMask.Matches(ProcMaskMeleeMH) {
		ppmm.mhProcChance = ppm * aa.MH.SwingSpeed / 60.0
	}
	if procMask.Matches(ProcMaskMeleeOH) {
		ppmm.ohProcChance = ppm * aa.OH.SwingSpeed / 60.0
	}
	if procMask.Matches(ProcMaskRanged) {
		ppmm.rangedProcChance = ppm * aa.Ranged.SwingSpeed / 60.0
	}

	return ppmm
}

// Returns whether a PPM-based effect procced.
//
// Using NewPPMManager() is preferred; this function should only be used when
// the attacker is not known at initialization time.
func (aa *AutoAttacks) PPMProc(sim *Simulation, ppm float64, procMask ProcMask, label string) bool {
	if !aa.IsEnabled() {
		return false
	}

	procChance := 0.0
	if procMask.Matches(ProcMaskMeleeMH) {
		procChance = ppm * aa.MH.SwingSpeed / 60.0
	} else if procMask.Matches(ProcMaskMeleeOH) {
		procChance = ppm * aa.OH.SwingSpeed / 60.0
	} else if procMask.Matches(ProcMaskRanged) {
		procChance = ppm * aa.Ranged.SwingSpeed / 60.0
	}

	return procChance > 0 && sim.RandomFloat(label) < procChance
}

func (unit *Unit) applyParryHaste() {
	if !unit.PseudoStats.ParryHaste || !unit.AutoAttacks.IsEnabled() {
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
			swingSpeed := aura.Unit.AutoAttacks.MainhandSwingSpeed()
			minRemainingTime := time.Duration(float64(swingSpeed) * 0.2) // 20% of Swing Speed
			defaultReduction := minRemainingTime * 2                     // 40% of Swing Speed

			if remainingTime <= minRemainingTime {
				return
			}

			parryHasteReduction := MinDuration(defaultReduction, remainingTime-minRemainingTime)
			newReadyAt := aura.Unit.AutoAttacks.MainhandSwingAt - parryHasteReduction
			if sim.Log != nil {
				aura.Unit.Log(sim, "MH Swing reduced by %s due to parry haste, will now occur at %s", parryHasteReduction, newReadyAt)
			}

			aura.Unit.AutoAttacks.MainhandSwingAt = newReadyAt
			aura.Unit.AutoAttacks.resetAutoSwing(sim)
		},
	})
}
