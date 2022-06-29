package core

import (
	"fmt"
	"time"

	"github.com/wowsims/tbc/sim/core/items"
	"github.com/wowsims/tbc/sim/core/proto"
)

// ReplaceMHSwing is called right before an auto attack fires
//  If it returns nil, the attack takes place as normal. If it returns a Spell,
//  that Spell is used in place of the attack.
//  This allows for abilities that convert a white attack into yellow attack.
type ReplaceMHSwing func(sim *Simulation, mhSwingSpell *Spell) *Spell

// Represents a generic weapon. Pets / unarmed / various other cases dont use
// actual weapon items so this is an abstraction of a Weapon.
type Weapon struct {
	BaseDamageMin        float64
	BaseDamageMax        float64
	SwingSpeed           float64
	NormalizedSwingSpeed float64
	SwingDuration        time.Duration // Duration between 2 swings.
	CritMultiplier       float64
	SpellSchool          SpellSchool
}

func newWeaponFromUnarmed(critMultiplier float64) Weapon {
	// These numbers are probably wrong but nobody cares.
	return Weapon{
		BaseDamageMin:        0,
		BaseDamageMax:        0,
		SwingSpeed:           1,
		NormalizedSwingSpeed: 1,
		SwingDuration:        time.Second,
		CritMultiplier:       critMultiplier,
	}
}

func newWeaponFromItem(item items.Item, critMultiplier float64) Weapon {
	normalizedWeaponSpeed := 2.4
	if item.WeaponType == proto.WeaponType_WeaponTypeDagger {
		normalizedWeaponSpeed = 1.7
	} else if item.HandType == proto.HandType_HandTypeTwoHand {
		normalizedWeaponSpeed = 3.3
	} else if item.RangedWeaponType != proto.RangedWeaponType_RangedWeaponTypeUnknown {
		normalizedWeaponSpeed = 2.8
	}

	return Weapon{
		BaseDamageMin:        item.WeaponDamageMin,
		BaseDamageMax:        item.WeaponDamageMax,
		SwingSpeed:           item.SwingSpeed,
		NormalizedSwingSpeed: normalizedWeaponSpeed,
		SwingDuration:        time.Duration(item.SwingSpeed * float64(time.Second)),
		CritMultiplier:       critMultiplier,
	}
}

// Returns weapon stats using the main hand equipped weapon.
func (character *Character) WeaponFromMainHand(critMultiplier float64) Weapon {
	if weapon := character.GetMHWeapon(); weapon != nil {
		return newWeaponFromItem(*weapon, critMultiplier)
	} else {
		return newWeaponFromUnarmed(critMultiplier)
	}
}

// Returns weapon stats using the off hand equipped weapon.
func (character *Character) WeaponFromOffHand(critMultiplier float64) Weapon {
	if weapon := character.GetOHWeapon(); weapon != nil {
		return newWeaponFromItem(*weapon, critMultiplier)
	} else {
		return Weapon{}
	}
}

// Returns weapon stats using the off hand equipped weapon.
func (character *Character) WeaponFromRanged(critMultiplier float64) Weapon {
	if weapon := character.GetRangedWeapon(); weapon != nil {
		return newWeaponFromItem(*weapon, critMultiplier)
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

func (weapon Weapon) EnemyWeaponDamage(sim *Simulation, attackPower float64) float64 {
	rand := 1 + 0.5*sim.RandomFloat("Enemy Weapon Damage")
	return weapon.BaseDamageMin * (rand + attackPower*EnemyAutoAttackAPCoefficient)
}

func (weapon Weapon) BaseDamage(sim *Simulation) float64 {
	return weapon.BaseDamageMin + (weapon.BaseDamageMax-weapon.BaseDamageMin)*sim.RandomFloat("Weapon Base Damage")
}

func (weapon Weapon) AverageDamage() float64 {
	return (weapon.BaseDamageMin + weapon.BaseDamageMax) / 2
}

func (weapon Weapon) CalculateWeaponDamage(sim *Simulation, attackPower float64) float64 {
	return weapon.BaseDamage(sim) + (weapon.SwingSpeed*attackPower)/MeleeAttackRatingPerDamage
}

func (weapon Weapon) CalculateAverageWeaponDamage(attackPower float64) float64 {
	return weapon.AverageDamage() + (weapon.SwingSpeed*attackPower)/MeleeAttackRatingPerDamage
}

func (weapon Weapon) CalculateNormalizedWeaponDamage(sim *Simulation, attackPower float64) float64 {
	return weapon.BaseDamage(sim) + (weapon.NormalizedSwingSpeed*attackPower)/MeleeAttackRatingPerDamage
}

type MeleeDamageCalculator func(attackPower float64, bonusWeaponDamage float64) float64

// Returns whether this hit effect is associated with the main-hand weapon.
func (ahe *SpellEffect) IsMH() bool {
	const mhmask = ProcMaskMeleeMH
	return ahe.ProcMask.Matches(mhmask)
}

// Returns whether this hit effect is associated with the off-hand weapon.
func (ahe *SpellEffect) IsOH() bool {
	return ahe.ProcMask.Matches(ProcMaskMeleeOH)
}

// Returns whether this hit effect is associated with either melee weapon.
func (ahe *SpellEffect) IsMelee() bool {
	return ahe.ProcMask.Matches(ProcMaskMelee)
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

	// Set this to true to use the OH delay macro, mostly used by enhance shamans.
	// This will intentionally delay OH swings to that they always fall within the
	// 0.5s window following a MH swing.
	DelayOHSwings bool

	MainhandSwingAt time.Duration
	OffhandSwingAt  time.Duration
	RangedSwingAt   time.Duration

	MHEffect     SpellEffect
	OHEffect     SpellEffect
	RangedEffect SpellEffect

	MHAuto     *Spell
	OHAuto     *Spell
	RangedAuto *Spell

	RangedSwingInProgress bool

	ReplaceMHSwing ReplaceMHSwing

	// The time at which the last MH swing occurred.
	previousMHSwingAt time.Duration

	// PendingAction which handles auto attacks.
	autoSwingAction    *PendingAction
	autoSwingCancelled bool
}

// Options for initializing auto attacks.
type AutoAttackOptions struct {
	MainHand       Weapon
	OffHand        Weapon
	Ranged         Weapon
	AutoSwingMelee bool // If true, core engine will handle calling SwingMelee() for you.
	DelayOHSwings  bool
	ReplaceMHSwing ReplaceMHSwing
}

func (unit *Unit) EnableAutoAttacks(agent Agent, options AutoAttackOptions) {
	unit.AutoAttacks = AutoAttacks{
		agent:          agent,
		unit:           unit,
		MH:             options.MainHand,
		OH:             options.OffHand,
		Ranged:         options.Ranged,
		AutoSwingMelee: options.AutoSwingMelee,
		DelayOHSwings:  options.DelayOHSwings,
		ReplaceMHSwing: options.ReplaceMHSwing,
		IsDualWielding: options.MainHand.SwingSpeed != 0 && options.OffHand.SwingSpeed != 0,
	}

	if unit.Type == EnemyUnit {
		unit.AutoAttacks.MHEffect = SpellEffect{
			ProcMask:         ProcMaskMeleeMHAuto,
			DamageMultiplier: 1,
			ThreatMultiplier: 1,
			BaseDamage:       BaseDamageConfigEnemyWeapon(MainHand),
			OutcomeApplier:   unit.OutcomeFuncEnemyMeleeWhite(),
		}
		unit.AutoAttacks.OHEffect = SpellEffect{
			ProcMask:         ProcMaskMeleeOHAuto,
			DamageMultiplier: 1,
			ThreatMultiplier: 1,
			BaseDamage:       BaseDamageConfigEnemyWeapon(OffHand),
			OutcomeApplier:   unit.OutcomeFuncEnemyMeleeWhite(),
		}
	} else {
		unit.AutoAttacks.MHEffect = SpellEffect{
			ProcMask:         ProcMaskMeleeMHAuto,
			DamageMultiplier: 1,
			ThreatMultiplier: 1,
			BaseDamage:       BaseDamageConfigMeleeWeapon(MainHand, false, 0, 1, true),
			OutcomeApplier:   unit.OutcomeFuncMeleeWhite(options.MainHand.CritMultiplier),
		}
		unit.AutoAttacks.OHEffect = SpellEffect{
			ProcMask:         ProcMaskMeleeOHAuto,
			DamageMultiplier: 1,
			ThreatMultiplier: 1,
			BaseDamage:       BaseDamageConfigMeleeWeapon(OffHand, false, 0, 1, true),
			OutcomeApplier:   unit.OutcomeFuncMeleeWhite(options.OffHand.CritMultiplier),
		}
		unit.AutoAttacks.RangedEffect = SpellEffect{
			ProcMask:         ProcMaskRangedAuto,
			DamageMultiplier: 1,
			ThreatMultiplier: 1,
			BaseDamage:       BaseDamageConfigRangedWeapon(0),
			OutcomeApplier:   unit.OutcomeFuncRangedHitAndCrit(options.Ranged.CritMultiplier),
		}
	}
}

func (aa *AutoAttacks) IsEnabled() bool {
	return aa.MH.SwingSpeed != 0
}

// Empty handler so Agents don't have to provide one if they have no logic to add.
func (unit *Unit) OnAutoAttack(sim *Simulation, spell *Spell) {}

func (aa *AutoAttacks) reset(sim *Simulation) {
	if !aa.IsEnabled() {
		return
	}

	aa.MHAuto = aa.unit.GetOrRegisterSpell(SpellConfig{
		ActionID:    ActionID{OtherID: proto.OtherAction_OtherActionAttack, Tag: 1},
		SpellSchool: aa.MH.GetSpellSchool(),
		Flags:       SpellFlagMeleeMetrics,

		ApplyEffects: ApplyEffectFuncDirectDamage(aa.MHEffect),
	})

	aa.OHAuto = aa.unit.GetOrRegisterSpell(SpellConfig{
		ActionID:    ActionID{OtherID: proto.OtherAction_OtherActionAttack, Tag: 2},
		SpellSchool: aa.OH.GetSpellSchool(),
		Flags:       SpellFlagMeleeMetrics,

		ApplyEffects: ApplyEffectFuncDirectDamage(aa.OHEffect),
	})

	if aa.RangedEffect.ProcMask != ProcMaskUnknown {
		aa.RangedAuto = aa.unit.GetOrRegisterSpell(SpellConfig{
			ActionID:    ActionID{OtherID: proto.OtherAction_OtherActionShoot},
			SpellSchool: SpellSchoolPhysical,
			Flags:       SpellFlagMeleeMetrics,

			Cast: CastConfig{
				DefaultCast: Cast{
					CastTime: 1, // Dummy non-zero value so the optimization doesnt remove the cast time.
				},
				ModifyCast: func(_ *Simulation, _ *Spell, cast *Cast) {
					cast.CastTime = aa.RangedSwingWindup()
				},
				IgnoreHaste: true,
				AfterCast: func(sim *Simulation, spell *Spell) {
					aa.RangedSwingInProgress = false
					aa.agent.OnAutoAttack(sim, aa.RangedAuto)
				},
			},

			ApplyEffects: ApplyEffectFuncDirectDamage(aa.RangedEffect),
		})
	}

	aa.MainhandSwingAt = 0
	aa.OffhandSwingAt = 0

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
	aa.resetAutoSwing(sim)

	aa.RangedSwingAt = 0
	aa.RangedSwingInProgress = false
}

func (aa *AutoAttacks) resetAutoSwing(sim *Simulation) {
	if aa.autoSwingCancelled || !aa.AutoSwingMelee {
		return
	}

	if aa.autoSwingAction != nil {
		aa.autoSwingAction.Cancel(sim)
	}

	pa := &PendingAction{
		NextActionAt: aa.NextAttackAt(),
		Priority:     ActionPriorityAuto,
	}

	pa.OnAction = func(sim *Simulation) {
		aa.SwingMelee(sim, aa.unit.CurrentTarget)
		pa.NextActionAt = aa.NextAttackAt()

		// Cancelled means we made a new one because of a swing speed change.
		if !pa.cancelled {
			sim.AddPendingAction(pa)
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
		aa.autoSwingCancelled = true
	}
}

// Renables the auto swing action for the iteration
func (aa *AutoAttacks) EnableAutoSwing(sim *Simulation) {
	// Already enabled so nothing to do
	if aa.autoSwingAction != nil {
		return
	}

	if aa.MainhandSwingAt < sim.CurrentTime {
		aa.MainhandSwingAt = sim.CurrentTime
	}
	if aa.OffhandSwingAt < sim.CurrentTime {
		aa.OffhandSwingAt = sim.CurrentTime
	}
	if aa.RangedSwingAt < sim.CurrentTime {
		if aa.RangedSwingInProgress {
			panic("Ranged swing already in progress!")
		}
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

// Ranged swings have a 0.5s 'windup' time before they can fire, affected by haste.
// This function computes the amount of windup time based on the current haste.
func (aa *AutoAttacks) RangedSwingWindup() time.Duration {
	return time.Duration(float64(time.Millisecond*500) / aa.unit.RangedSwingSpeed())
}

// Time between a ranged auto finishes casting and the next one becomes available.
func (aa *AutoAttacks) RangedSwingGap() time.Duration {
	return time.Duration(float64(aa.Ranged.SwingDuration-time.Millisecond*500) / aa.unit.RangedSwingSpeed())
}

// Returns the amount of time available before ranged auto will be clipped.
func (aa *AutoAttacks) TimeBeforeClippingRanged(sim *Simulation) time.Duration {
	return aa.RangedSwingAt - aa.RangedSwingWindup() - sim.CurrentTime
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
	aa.agent.OnAutoAttack(sim, attackSpell)
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

	if aa.DelayOHSwings && (sim.CurrentTime-aa.previousMHSwingAt) > time.Millisecond*500 {
		// Delay the OH swing for later, so it follows the MH swing.
		aa.OffhandSwingAt = aa.MainhandSwingAt + time.Millisecond*100
		if sim.Log != nil {
			aa.unit.Log(sim, "Delaying OH swing by %s", aa.OffhandSwingAt-sim.CurrentTime)
		}
		return
	}

	aa.OHAuto.Cast(sim, target)
	aa.OffhandSwingAt = sim.CurrentTime + aa.OffhandSwingSpeed()
	aa.agent.OnAutoAttack(sim, aa.OHAuto)
}

// Performs an autoattack using the ranged weapon, if the ranged CD is ready.
func (aa *AutoAttacks) TrySwingRanged(sim *Simulation, target *Unit) {
	if aa.RangedSwingAt > sim.CurrentTime {
		return
	}

	aa.RangedAuto.Cast(sim, target)
	aa.RangedSwingAt = sim.CurrentTime + aa.RangedSwingSpeed()
	aa.RangedSwingInProgress = true

	// It's important that we update the GCD timer AFTER starting the ranged auto.
	// Otherwise the hardcast action won't be created separately.
	nextGCD := sim.CurrentTime + aa.RangedAuto.CurCast.CastTime
	if nextGCD > aa.unit.NextGCDAt() {
		aa.unit.SetGCDTimer(sim, nextGCD)
	}
}

func (aa *AutoAttacks) ModifySwingTime(sim *Simulation, amount float64) {
	if !aa.IsEnabled() {
		return
	}

	mhSwingTime := aa.MainhandSwingAt - sim.CurrentTime
	if mhSwingTime > 1 { // If its 1 we end up rounding down to 0 and causing a panic.
		aa.MainhandSwingAt = sim.CurrentTime + time.Duration(float64(mhSwingTime)/amount)
	}

	if aa.OH.SwingSpeed != 0 {
		ohSwingTime := aa.OffhandSwingAt - sim.CurrentTime
		if ohSwingTime > 1 {
			newTime := time.Duration(float64(ohSwingTime) / amount)
			if newTime > 0 {
				aa.OffhandSwingAt = sim.CurrentTime + newTime
			}
		}
	}

	aa.resetAutoSwing(sim)
}

// Delays all swing timers until the specified time.
func (aa *AutoAttacks) DelayAllUntil(sim *Simulation, readyAt time.Duration) {
	autoChanged := false

	if readyAt > aa.MainhandSwingAt {
		aa.MainhandSwingAt = readyAt
		if aa.AutoSwingMelee {
			autoChanged = true
		}
	}
	if readyAt > aa.OffhandSwingAt {
		aa.OffhandSwingAt = readyAt
		if aa.AutoSwingMelee {
			autoChanged = true
		}
	}
	if readyAt > aa.RangedSwingAt {
		if aa.RangedSwingInProgress {
			panic("Ranged swing already in progress!")
		}
		aa.RangedSwingAt = readyAt
	}

	if autoChanged {
		aa.resetAutoSwing(sim)
	}
}

func (aa *AutoAttacks) DelayRangedUntil(sim *Simulation, readyAt time.Duration) {
	if aa.RangedSwingInProgress {
		panic("Ranged swing already in progress!")
	}
	if readyAt > aa.RangedSwingAt {
		aa.RangedSwingAt = readyAt
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

// Returns the time at which all melee swings will be ready.
func (aa *AutoAttacks) MeleeSwingsReadyAt() time.Duration {
	return MaxDuration(aa.MainhandSwingAt, aa.OffhandSwingAt)
}

// Returns true if all melee weapons are ready for a swing.
func (aa *AutoAttacks) MeleeSwingsReady(sim *Simulation) bool {
	return aa.MainhandSwingAt <= sim.CurrentTime &&
		(aa.OH.SwingSpeed == 0 || aa.OffhandSwingAt <= sim.CurrentTime)
}

// Returns the time at which the next event will occur, considering both autos and the gcd.
func (aa *AutoAttacks) NextEventAt(sim *Simulation) time.Duration {
	if aa.NextAttackAt() == sim.CurrentTime {
		panic(fmt.Sprintf("Returned 0 from next attack at %s, mh: %s, oh: %s", sim.CurrentTime, aa.MainhandSwingAt, aa.OffhandSwingAt))
	}
	return MinDuration(
		sim.CurrentTime+aa.unit.GCD.TimeToReady(sim),
		aa.NextAttackAt())
}

type PPMManager struct {
	mhProcChance     float64
	ohProcChance     float64
	rangedProcChance float64

	// For feral druids, certain PPM effects use their equipped weapon speed
	// instead of their paw attack speed.
	mhSpecialProcChance float64
	ohSpecialProcChance float64
}

// Returns whether the effect procced.
func (ppmm *PPMManager) Proc(sim *Simulation, procMask ProcMask, label string) bool {
	if procMask.Matches(ProcMaskMeleeMH) {
		return ppmm.mhProcChance > 0 && sim.RandomFloat(label) < ppmm.mhProcChance
	} else if procMask.Matches(ProcMaskMeleeOH) {
		return ppmm.ohProcChance > 0 && sim.RandomFloat(label) < ppmm.ohProcChance
	} else if procMask.Matches(ProcMaskRanged) {
		return ppmm.rangedProcChance > 0 && sim.RandomFloat(label) < ppmm.rangedProcChance
	}
	return false
}

// Returns whether the effect procced.
// This is different from Proc() in that yellow melee hits use a proc chance based on the equipped
// weapon speed rather than the base attack speed. This distinction matters for feral druids.
func (ppmm *PPMManager) ProcWithWeaponSpecials(sim *Simulation, procMask ProcMask, label string) bool {
	if procMask.Matches(ProcMaskMeleeMHAuto) {
		return ppmm.mhProcChance > 0 && sim.RandomFloat(label) < ppmm.mhProcChance
	} else if procMask.Matches(ProcMaskMeleeMHSpecial) {
		return ppmm.mhSpecialProcChance > 0 && sim.RandomFloat(label) < ppmm.mhSpecialProcChance
	} else if procMask.Matches(ProcMaskMeleeOHAuto) {
		return ppmm.ohProcChance > 0 && sim.RandomFloat(label) < ppmm.ohProcChance
	} else if procMask.Matches(ProcMaskMeleeOHSpecial) {
		return ppmm.ohSpecialProcChance > 0 && sim.RandomFloat(label) < ppmm.ohSpecialProcChance
	} else if procMask.Matches(ProcMaskRanged) {
		return ppmm.rangedProcChance > 0 && sim.RandomFloat(label) < ppmm.rangedProcChance
	}
	return false
}

func (aa *AutoAttacks) NewPPMManager(ppm float64, procMask ProcMask) PPMManager {
	if !aa.IsEnabled() {
		return PPMManager{}
	}

	character := aa.agent.GetCharacter()

	ppmm := PPMManager{}
	if procMask.Matches(ProcMaskMeleeMH) {
		ppmm.mhProcChance = ppm * aa.MH.SwingSpeed / 60.0
		ppmm.mhSpecialProcChance = ppmm.mhProcChance
		if character != nil {
			if mhWeapon := character.GetMHWeapon(); mhWeapon != nil {
				ppmm.mhSpecialProcChance = ppm * mhWeapon.SwingSpeed / 60.0
			}
		}
	}
	if procMask.Matches(ProcMaskMeleeOH) {
		ppmm.ohProcChance = ppm * aa.OH.SwingSpeed / 60.0
		ppmm.ohSpecialProcChance = ppmm.ohProcChance
		if character != nil {
			if ohWeapon := character.GetOHWeapon(); ohWeapon != nil {
				ppmm.ohSpecialProcChance = ppm * ohWeapon.SwingSpeed / 60.0
			}
		}
	}
	if procMask.Matches(ProcMaskRanged) {
		ppmm.rangedProcChance = ppm * aa.Ranged.SwingSpeed / 60.0
	}

	return ppmm
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
		OnSpellHitTaken: func(aura *Aura, sim *Simulation, spell *Spell, spellEffect *SpellEffect) {
			if !spellEffect.Outcome.Matches(OutcomeParry) {
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
