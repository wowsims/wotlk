package feral

import (
	"time"

	"golang.org/x/exp/slices"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/druid"
)

func (cat *FeralDruid) OnEnergyGain(sim *core.Simulation) {
	cat.TryUseCooldowns(sim)
	if cat.InForm(druid.Cat) && !cat.readyToShift {
		cat.doTigersFury(sim)
	}
}

func (cat *FeralDruid) OnGCDReady(sim *core.Simulation) {
	cat.TryUseCooldowns(sim)
	if !cat.GCD.IsReady(sim) {
		return
	}

	cat.doRotation(sim)

	// Replace gcd event with our own if we casted a spell
	if !cat.GCD.IsReady(sim) {
		nextGcd := cat.NextGCDAt()
		cat.DoNothing()
		cat.CancelGCDTimer(sim)

		cat.NextRotationAction(sim, nextGcd)
	}
}

func (cat *FeralDruid) OnAutoAttack(sim *core.Simulation, spell *core.Spell) {
	if cat.InForm(druid.Cat) {
		return
	}
	if cat.InForm(druid.Humanoid) {
		panic("auto attack out of form?")
	}

	// If the swing/Maul resulted in an Omen proc, then schedule the
	// next player decision based on latency.

	if cat.Talents.OmenOfClarity && cat.ClearcastingAura.RemainingDuration(sim) == cat.ClearcastingAura.Duration {
		// Kick gcd loop, also need to account for any gcd 'left'
		// otherwise it breaks gcd logic
		kickTime := core.MaxDuration(cat.NextGCDAt(), sim.CurrentTime+cat.latency)
		cat.NextRotationAction(sim, kickTime)
	}
}

func (cat *FeralDruid) NextRotationAction(sim *core.Simulation, kickAt time.Duration) {
	if cat.rotationAction != nil {
		cat.rotationAction.Cancel(sim)
	}

	cat.rotationAction = &core.PendingAction{
		Priority:     core.ActionPriorityGCD,
		OnAction:     cat.OnGCDReady,
		NextActionAt: kickAt,
	}

	sim.AddPendingAction(cat.rotationAction)
}

// Ported from https://github.com/NerdEgghead/WOTLK_cat_sim

func (cat *FeralDruid) checkReplaceMaul(sim *core.Simulation) *core.Spell {
	// If we will have enough time and Energy leeway to stay in
	// Dire Bear Form once the GCD expires, then only Maul if we
	// will be left with enough Rage to cast Mangle or Lacerate
	// on that global.

	ripDot := cat.Rip.CurDot()

	furorCap := core.MinFloat(20.0*float64(cat.Talents.Furor), 85.0)
	ripRefreshPending := ripDot.IsActive() && (ripDot.RemainingDuration(sim) < sim.GetRemainingDuration()-time.Second*10)
	gcdTimeToRdy := cat.GCD.TimeToReady(sim)
	energyLeeway := furorCap - 15.0 - float64((gcdTimeToRdy+cat.latency)/core.EnergyTickDuration)
	shiftNext := cat.CurrentEnergy() > energyLeeway

	if ripRefreshPending {
		shiftNext = shiftNext || (ripDot.RemainingDuration(sim) < (gcdTimeToRdy + time.Second*3))
	}

	lacerateNext := false
	emergencyLacerateNext := false
	mangleNext := false

	lacerateDot := cat.Lacerate.CurDot()
	if cat.Rotation.BearweaveType == proto.FeralDruid_Rotation_Lacerate {
		lacerateLeeway := cat.Rotation.LacerateTime + gcdTimeToRdy
		lacerateNext = !lacerateDot.IsActive() || (lacerateDot.GetStacks() < 5) || (lacerateDot.RemainingDuration(sim) <= lacerateLeeway)
		emergencyLeeway := gcdTimeToRdy + (3 * time.Second) + (2 * cat.latency)
		emergencyLacerateNext = lacerateDot.IsActive() && (lacerateDot.RemainingDuration(sim) <= emergencyLeeway)
		mangleNext = cat.MangleBear != nil && !lacerateNext && (!cat.bleedAura.IsActive() || (cat.bleedAura.RemainingDuration(sim) < gcdTimeToRdy+time.Second*3) || (sim.CurrentTime-cat.lastShift < time.Duration(1500*time.Millisecond)))
	} else {
		mangleNext = cat.MangleBear != nil && cat.MangleBear.TimeToReady(sim) < gcdTimeToRdy
		lacerateNext = lacerateDot.IsActive() && (lacerateDot.GetStacks() < 5 || lacerateDot.RemainingDuration(sim) < gcdTimeToRdy+(time.Second*4))
	}

	maulRageThresh := 10.0
	if emergencyLacerateNext {
		maulRageThresh += cat.Lacerate.DefaultCast.Cost
	} else if shiftNext {
		maulRageThresh = 10.0
	} else if mangleNext {
		maulRageThresh += cat.MangleBear.DefaultCast.Cost
	} else if lacerateNext {
		maulRageThresh += cat.Lacerate.DefaultCast.Cost
	}

	if cat.CurrentRage() >= maulRageThresh {
		return cat.Maul
	} else {
		return nil
	}
}

func (cat *FeralDruid) shiftBearCat(sim *core.Simulation, powershift bool) bool {
	cat.waitingForTick = false

	// If we have just now decided to shift, then we do not execute the
	// shift immediately, but instead trigger an input delay for realism.
	if !cat.readyToShift {
		cat.readyToShift = true
		return false
	}
	cat.readyToShift = false

	toCat := !cat.InForm(druid.Cat)
	if powershift {
		toCat = !toCat
	}

	cat.lastShift = sim.CurrentTime
	if toCat {
		return cat.CatForm.Cast(sim, nil)
	} else {
		cat.BearForm.Cast(sim, nil)
		// Bundle Enrage if available
		if cat.Enrage.IsReady(sim) {
			cat.Enrage.Cast(sim, nil)
		}
		return true
	}
}

func (cat *FeralDruid) flowerCast(sim *core.Simulation) {
	cat.readyToGift = false
	cat.ClearForm(sim)

	if !cat.GiftOfTheWild.Cast(sim, &cat.Unit) {
		panic("gotw cast failed")
	}
	// Purposely just using 'Cancel' here to avoid any caster melee attacks
	// Basically mimicing a '/stopattack' macro on cast
	cat.AutoAttacks.CancelAutoSwing(sim)
}

func (cat *FeralDruid) canBite(sim *core.Simulation) bool {
	return cat.Rip.CurDot().RemainingDuration(sim) >= cat.Rotation.BiteTime &&
		cat.SavageRoarAura.RemainingDuration(sim) >= cat.Rotation.BiteTime
}

func (cat *FeralDruid) berserkExpectedAt(sim *core.Simulation, futureTime time.Duration) bool {
	if cat.BerserkAura.IsActive() {
		return futureTime < cat.BerserkAura.ExpiresAt() || futureTime > cat.Berserk.ReadyAt()
	}
	if cat.Berserk.IsReady(sim) {
		return futureTime > sim.CurrentTime+cat.Berserk.CD.Duration
	}
	if cat.TigersFuryAura.IsActive() && cat.Talents.Berserk {
		return futureTime > cat.TigersFuryAura.ExpiresAt()
	}
	return false
}

func (cat *FeralDruid) calcBuilderDpe(sim *core.Simulation) (float64, float64) {
	// Calculate current damage-per-Energy of Rake vs. Shred. Used to
	// determine whether Rake is worth casting when player stats change upon a
	// dynamic proc occurring
	shredDpc := cat.Shred.ExpectedDamage(sim, cat.CurrentTarget)
	rakeDpc := cat.Rake.ExpectedDamage(sim, cat.CurrentTarget)
	return rakeDpc / 35., shredDpc / 42.
}

func (cat *FeralDruid) clipRoar(sim *core.Simulation) bool {
	ripDot := cat.Rip.CurDot()
	ripdotRemaining := ripDot.RemainingDuration(sim)
	if !ripDot.IsActive() || (sim.GetRemainingDuration()-ripdotRemaining < 10*time.Second) {
		return false
	}

	// Project Rip end time assuming full Glyph of Shred extensions
	maxRipDur := time.Duration(cat.maxRipTicks) * ripDot.TickLength
	ripDur := ripDot.Aura.StartedAt() + maxRipDur - sim.CurrentTime
	roarDur := cat.SavageRoarAura.RemainingDuration(sim)

	if roarDur > ripDur {
		return false
	}

	// Calculate when roar would end if casted now
	newRoarDur := cat.SavageRoarDurationTable[cat.ComboPoints()]

	// Clip as soon as we have enough CPs for the new roar to expire well
	// after the current rip
	return newRoarDur >= (ripDur + cat.Rotation.MinRoarOffset)
}

func (cat *FeralDruid) tfExpectedBefore(sim *core.Simulation, futureTime time.Duration) bool {
	if !cat.TigersFury.IsReady(sim) {
		return cat.TigersFury.ReadyAt() < futureTime
	}
	if cat.BerserkAura.IsActive() {
		return cat.BerserkAura.ExpiresAt() < futureTime
	}
	return true
}

func (cat *FeralDruid) doTigersFury(sim *core.Simulation) {
	// Handle tigers fury
	if !cat.TigersFury.IsReady(sim) {
		return
	}

	gcdTimeToRdy := cat.GCD.TimeToReady(sim)
	leewayTime := core.MaxDuration(gcdTimeToRdy, cat.latency)
	tfEnergyThresh := 40.0 - 10.0*(leewayTime+core.TernaryDuration(cat.ClearcastingAura.IsActive(), 1*time.Second, 0)).Seconds()
	tfNow := (cat.CurrentEnergy() < tfEnergyThresh) && !cat.BerserkAura.IsActive()

	// If Lacerateweaving, then delay Tiger's Fury if Lacerate is due to
	// expire within 3 GCDs (two cat specials + shapeshift), since we
	// won't be able to spend down our Energy fast enough to avoid
	// Energy capping otherwise.
	lacerateDot := cat.Lacerate.CurDot()
	if cat.Rotation.BearweaveType == proto.FeralDruid_Rotation_Lacerate {
		nextPossibleLac := sim.CurrentTime + leewayTime + cat.latency + time.Duration(3.5*float64(time.Second))
		tfNow = tfNow && (!lacerateDot.IsActive() || (lacerateDot.ExpiresAt() > nextPossibleLac) || (lacerateDot.RemainingDuration(sim) > sim.GetRemainingDuration()))
	}

	if tfNow {
		cat.TigersFury.Cast(sim, nil)
		// Kick gcd loop, also need to account for any gcd 'left'
		// otherwise it breaks gcd logic
		cat.NextRotationAction(sim, sim.CurrentTime+leewayTime)
	}
}

func (cat *FeralDruid) doRotation(sim *core.Simulation) {
	if cat.BerserkAura.IsActive() {
		cat.berserkUsed = true
	}

	// If previously decided to gift, then gift
	if cat.readyToGift {
		cat.flowerCast(sim)
		return
	}

	// If we previously decided to shift, then execute the shift now once
	// the input delay is over.
	if cat.readyToShift {
		cat.shiftBearCat(sim, false)
		// Reset swing timer from snek (or idol/weapon swap) when going into cat
		if cat.InForm(druid.Cat) && cat.Rotation.SnekWeave {
			cat.AutoAttacks.StopMeleeUntil(sim, sim.CurrentTime, false)
		}
		return
	}

	rotation := &cat.Rotation

	if rotation.MaintainFaerieFire && cat.ShouldFaerieFire(sim, cat.CurrentTarget) {
		cat.FaerieFire.Cast(sim, cat.CurrentTarget)
		return
	}

	cat.missChance = cat.MissChance()
	cat.bleedAura = cat.CurrentTarget.GetExclusiveEffectCategory(core.BleedEffectCategory).GetActiveAura()

	curEnergy := cat.CurrentEnergy()
	curRage := cat.CurrentRage()
	curCp := cat.ComboPoints()
	isClearcast := cat.ClearcastingAura.IsActive()
	simTimeRemain := sim.GetRemainingDuration()
	shiftCost := cat.CatForm.DefaultCast.Cost
	rakeDot := cat.Rake.CurDot()
	ripDot := cat.Rip.CurDot()
	lacerateDot := cat.Lacerate.CurDot()

	endThresh := time.Second * 10

	ripNow := (curCp >= rotation.MinCombosForRip) && !ripDot.IsActive() && (simTimeRemain >= endThresh) && !isClearcast
	biteAtEnd := (curCp >= rotation.MinCombosForBite) && ((simTimeRemain < endThresh) || (ripDot.IsActive() && (simTimeRemain-ripDot.RemainingDuration(sim) < endThresh)))
	mangleNow := !ripNow && !cat.bleedAura.IsActive() && cat.MangleCat != nil

	biteBeforeRip := (curCp >= rotation.MinCombosForBite) && ripDot.IsActive() && cat.SavageRoarAura.IsActive() && rotation.UseBite && cat.canBite(sim)
	biteNow := (biteBeforeRip || biteAtEnd) && !isClearcast

	// During Berserk, we additionally add an Energy constraint on Bite
	// usage to maximize the total Energy expenditure we can get.
	if biteNow && cat.BerserkAura.IsActive() {
		biteNow = curEnergy <= rotation.BerserkBiteThresh
	}

	rakeNow := rotation.UseRake && !rakeDot.IsActive() && (simTimeRemain > rakeDot.Duration) && !isClearcast

	// Additionally, don't Rake if the current Shred DPE is higher due to
	// trinket procs etc.
	if rakeNow {
		rakeDpe, shredDpe := cat.calcBuilderDpe(sim)
		rakeNow = (rakeDpe > shredDpe)
	}

	// Berserk algorithm: time Berserk for just after a Tiger's Fury
	// *unless* we'll lose Berserk uptime by waiting for Tiger's Fury to
	// come off cooldown. The latter exception is necessary for
	// Lacerateweave rotation since TF timings can drift over time.
	waitForTf := cat.Talents.Berserk && (cat.TigersFury.ReadyAt() <= cat.BerserkAura.Duration) && (cat.TigersFury.ReadyAt()+time.Second < sim.GetRemainingDuration()-cat.BerserkAura.Duration)
	berserkNow := cat.Berserk.IsReady(sim) && !waitForTf

	// Additionally, for Lacerateweave rotation, postpone the final Berserk
	// of the fight to as late as possible so as to minimize the impact of
	// dropping Lacerate stacks during the Berserk window. Rationale for the
	// 3 second additional leeway given beyond just berserk_dur in the below
	// expression is to be able to fit in a final TF and dump the Energy
	// from it in cases where Berserk and TF CDs are desynced due to drift.
	if berserkNow && rotation.BearweaveType == proto.FeralDruid_Rotation_Lacerate && cat.berserkUsed && sim.GetRemainingDuration() < cat.Berserk.CD.Duration {
		berserkNow = sim.GetRemainingDuration() < cat.BerserkAura.Duration+(3*time.Second)
	}

	roarNow := curCp >= 1 && (!cat.SavageRoarAura.IsActive() || cat.clipRoar(sim))

	ripRefreshPending := false
	pendingActions := make([]pendingAction, 0, 4)

	if ripDot.IsActive() && (ripDot.RemainingDuration(sim) < simTimeRemain-endThresh) && curCp == 5 {
		ripCost := core.TernaryFloat64(cat.berserkExpectedAt(sim, ripDot.ExpiresAt()), cat.Rip.DefaultCast.Cost*0.5, cat.Rip.DefaultCast.Cost)
		pendingActions = append(pendingActions, pendingAction{ripDot.ExpiresAt(), ripCost})
		ripRefreshPending = true
	}
	if rakeDot.IsActive() && (rakeDot.RemainingDuration(sim) < simTimeRemain-rakeDot.Duration) {
		rakeCost := core.TernaryFloat64(cat.berserkExpectedAt(sim, rakeDot.ExpiresAt()), cat.Rake.DefaultCast.Cost*0.5, cat.Rake.DefaultCast.Cost)
		pendingActions = append(pendingActions, pendingAction{rakeDot.ExpiresAt(), rakeCost})
	}
	if cat.bleedAura.IsActive() && (cat.bleedAura.RemainingDuration(sim) < simTimeRemain-time.Second) {
		mangleCost := core.TernaryFloat64(cat.berserkExpectedAt(sim, cat.bleedAura.ExpiresAt()), cat.MangleCat.DefaultCast.Cost*0.5, cat.MangleCat.DefaultCast.Cost)
		pendingActions = append(pendingActions, pendingAction{cat.bleedAura.ExpiresAt(), mangleCost})
	}
	if cat.SavageRoarAura.IsActive() {
		roarCost := core.TernaryFloat64(cat.berserkExpectedAt(sim, cat.SavageRoarAura.ExpiresAt()), cat.SavageRoar.DefaultCast.Cost*0.5, cat.SavageRoar.DefaultCast.Cost)
		pendingActions = append(pendingActions, pendingAction{cat.SavageRoarAura.ExpiresAt(), roarCost})
	}

	slices.SortStableFunc(pendingActions, func(p1, p2 pendingAction) bool {
		return p1.refreshTime < p2.refreshTime
	})

	latencySecs := cat.latency.Seconds()
	// Allow for bearweaving if the next pending action is >= 4.5s away
	furorCap := core.MinFloat(20.0*float64(cat.Talents.Furor), 85)
	weaveEnergy := furorCap - 30 - 20*latencySecs

	// With 4/5 or 5/5 Furor, force 2-GCD bearweaves whenever possible
	if cat.Talents.Furor > 3 {
		weaveEnergy -= 15.0

		// Force a 3-GCD weave when stacking Lacerates for the first time
		if rotation.BearweaveType == proto.FeralDruid_Rotation_Lacerate && !lacerateDot.IsActive() {
			weaveEnergy -= 15.0
		}
	}

	weaveEnd := time.Duration(float64(sim.CurrentTime) + (4.5+2*latencySecs)*float64(time.Second))
	bearweaveNow := rotation.BearweaveType != proto.FeralDruid_Rotation_None && curEnergy <= weaveEnergy && !isClearcast && (!ripRefreshPending || ripDot.ExpiresAt() >= weaveEnd) && !cat.BerserkAura.IsActive()

	if bearweaveNow && rotation.BearweaveType != proto.FeralDruid_Rotation_Lacerate {
		bearweaveNow = !cat.tfExpectedBefore(sim, weaveEnd)
	}

	// Also add an end of fight condition to make sure we can spend down our
	// Energy post-bearweave before the encounter ends. Time to spend is
	// given by weave_end plus 1 second per 42 Energy that we have at
	// weave_end.
	if bearweaveNow {
		energyToDump := curEnergy + ((weaveEnd - sim.CurrentTime).Seconds() * 10)
		bearweaveNow = weaveEnd+time.Duration((energyToDump/42)*float64(time.Second)) < sim.CurrentTime+sim.GetRemainingDuration()
	}

	// If we're maintaining Lacerate, then allow for emergency bearweaves
	// if Lacerate is about to fall off even if the above conditions do not
	// apply.
	lacRemain := core.TernaryDuration(lacerateDot.IsActive(), lacerateDot.RemainingDuration(sim), time.Duration(0))
	emergencyBearweave := rotation.BearweaveType == proto.FeralDruid_Rotation_Lacerate && lacerateDot.IsActive() && (float64(lacRemain) < (2.5+latencySecs)*float64(time.Second)) && (lacRemain < simTimeRemain) && !cat.BerserkAura.IsActive()

	// As an alternative to bearweaving, cast GotW on the raid under
	// analagous conditions to the above. Only difference is that there is
	// more available time/Energy leeway for the technique, since
	// flowershifts take only 3 seconds to execute.
	flowershiftEnergy := core.MinFloat(furorCap, 75) - 10*cat.SpellGCD().Seconds() - 20*latencySecs

	flowerEnd := time.Duration(float64(sim.CurrentTime) + (3.0+2*latencySecs)*float64(time.Second))
	flowershiftNow := rotation.FlowerWeave && (curEnergy <= flowershiftEnergy) && !isClearcast && (!ripRefreshPending || ripDot.ExpiresAt() >= flowerEnd) && !cat.BerserkAura.IsActive() && !cat.tfExpectedBefore(sim, flowerEnd)

	if bearweaveNow || emergencyBearweave {
		// oom check, if we arent able to shift into bear and back
		// then abandon bearweave
		if cat.CurrentMana() < shiftCost*2.0 {
			bearweaveNow = false
			emergencyBearweave = false
			cat.Metrics.MarkOOM(sim)
		}
	}

	if flowershiftNow {
		// if we cant cast and get back then abandon flowershift
		if cat.CurrentMana() <= shiftCost+cat.GiftOfTheWild.DefaultCast.Cost {
			flowershiftNow = false
			cat.Metrics.MarkOOM(sim)

		}
	}

	// Also add an end of fight condition to make sure we can spend down our
	// Energy post-flowershift before the encounter ends. Time to spend is
	// given by flower_end plus 1 second for Clearcast Shred plus 1 second
	// per 42 Energy that we have after that Clearcast Shred.
	if flowershiftNow {
		energyToDump := curEnergy + ((flowerEnd + time.Second - sim.CurrentTime).Seconds() * 10)
		flowershiftNow = flowerEnd+time.Second+time.Duration((energyToDump/42)*float64(time.Second)) < sim.CurrentTime+sim.GetRemainingDuration()
	}

	floatingEnergy := 0.0
	previousTime := sim.CurrentTime
	tfPending := false

	for _, s := range pendingActions {
		delta_t := float64((s.refreshTime - previousTime) / core.EnergyTickDuration)
		if !tfPending {
			tfPending = cat.tfExpectedBefore(sim, s.refreshTime)
			if tfPending {
				s.cost -= 60
			}
		}

		if delta_t < s.cost {
			floatingEnergy += s.cost - delta_t
			previousTime = s.refreshTime
		} else {
			previousTime += time.Duration(s.cost * float64(core.EnergyTickDuration))
		}
	}

	excessE := curEnergy - floatingEnergy
	timeToNextAction := time.Duration(0)

	if !cat.CatFormAura.IsActive() && rotation.FlowerWeave {
		// If the previous GotW cast was unsuccessful and we still have
		// leeway available, then try again. Otherwise, shift back into Cat
		// Form.
		if flowershiftNow {
			cat.flowerCast(sim)
		} else {
			cat.readyToShift = true
		}
	} else if !cat.CatFormAura.IsActive() {
		// Shift back into Cat Form if (a) our first bear auto procced
		// Clearcasting, or (b) our first bear auto didn't generate enough
		// Rage to Mangle or Maul, or (c) we don't have enough time or
		// Energy leeway to spend an additional GCD in Dire Bear Form.
		shiftNow := (curEnergy+15.0+(10.0*latencySecs) > furorCap) || (ripRefreshPending && (ripDot.RemainingDuration(sim) < (3.0 * time.Second))) || cat.BerserkAura.IsActive()
		shiftNext := (curEnergy+30.0+(10.0*latencySecs) > furorCap) || (ripRefreshPending && (ripDot.RemainingDuration(sim) < time.Duration(4500*time.Millisecond))) || cat.BerserkAura.IsActive()

		var powerbearNow bool
		if rotation.Powerbear {
			powerbearNow = !shiftNow && curRage < 10
		} else {
			powerbearNow = false
			shiftNow = shiftNow || curRage < 10
		}

		buildLacerate := !lacerateDot.IsActive() || lacerateDot.GetStacks() < 5
		maintainLacerate := !buildLacerate && (lacRemain <= rotation.LacerateTime) && (curRage < 38 || shiftNext) && (lacRemain < simTimeRemain)

		lacerateNow := rotation.BearweaveType == proto.FeralDruid_Rotation_Lacerate && (buildLacerate || maintainLacerate)
		emergencyLacerate := rotation.BearweaveType == proto.FeralDruid_Rotation_Lacerate && lacerateDot.IsActive() && (lacRemain < 3*time.Second+2*cat.latency) && lacRemain < simTimeRemain

		if (rotation.BearweaveType != proto.FeralDruid_Rotation_Lacerate) || !lacerateNow {
			shiftNow = shiftNow || isClearcast
		}

		// Also add an end of fight condition to prevent extending a weave
		// if we don't have enough time to spend the pooled Energy thus far.
		if !shiftNow {
			energyToDump := curEnergy + 30 + 10*latencySecs
			timeToDump := (3 * time.Second) + cat.latency + time.Duration((energyToDump/42)*float64(time.Second))
			shiftNow = timeToDump >= sim.GetRemainingDuration()
		}

		if emergencyLacerate && cat.Lacerate.CanCast(sim, cat.CurrentTarget) {
			cat.Lacerate.Cast(sim, cat.CurrentTarget)
			return
		} else if shiftNow {
			// If we are resetting our swing timer using Albino Snake or a
			// duplicate weapon swap, then do an additional check here to
			// see whether we can delay the shift until the next bear swing
			// goes out in order to maximize the gains from the reset.
			projectedDelay := cat.AutoAttacks.MainhandSwingAt + 2*cat.latency - sim.CurrentTime
			ripConflict := ripRefreshPending && (ripDot.ExpiresAt() < sim.CurrentTime+projectedDelay+(1500*time.Millisecond))
			nextCatSwing := sim.CurrentTime + cat.latency + time.Duration(float64(cat.AutoAttacks.MainhandSwingSpeed())/float64(2500*time.Millisecond))
			canDelayShift := !ripConflict && cat.Rotation.SnekWeave && (curEnergy+10*projectedDelay.Seconds() <= furorCap) && (cat.AutoAttacks.MainhandSwingAt < nextCatSwing)

			if canDelayShift {
				timeToNextAction = cat.AutoAttacks.MainhandSwingAt - sim.CurrentTime
			} else {
				cat.readyToShift = true
			}
		} else if powerbearNow {
			cat.shiftBearCat(sim, true)
		} else if lacerateNow && cat.Lacerate.CanCast(sim, cat.CurrentTarget) {
			cat.Lacerate.Cast(sim, cat.CurrentTarget)
			return
		} else if cat.MangleBear.CanCast(sim, cat.CurrentTarget) {
			cat.MangleBear.Cast(sim, cat.CurrentTarget)
			return
		} else if cat.Lacerate.CanCast(sim, cat.CurrentTarget) {
			cat.Lacerate.Cast(sim, cat.CurrentTarget)
			return
		} else {
			timeToNextAction = cat.AutoAttacks.MainhandSwingAt - sim.CurrentTime
		}
	} else if emergencyBearweave {
		cat.readyToShift = true
	} else if berserkNow {
		cat.Berserk.Cast(sim, nil)
		cat.UpdateMajorCooldowns()
		return
	} else if roarNow {
		if cat.SavageRoar.CanCast(sim, cat.CurrentTarget) {
			cat.SavageRoar.Cast(sim, nil)
			return
		}
		timeToNextAction = time.Duration((cat.CurrentSavageRoarCost() - curEnergy) * float64(core.EnergyTickDuration))
	} else if ripNow {
		if cat.Rip.CanCast(sim, cat.CurrentTarget) {
			cat.Rip.Cast(sim, cat.CurrentTarget)
			return
		}
		timeToNextAction = time.Duration((cat.CurrentRipCost() - curEnergy) * float64(core.EnergyTickDuration))
	} else if biteNow {
		if cat.FerociousBite.CanCast(sim, cat.CurrentTarget) {
			cat.FerociousBite.Cast(sim, cat.CurrentTarget)
			return
		}
		timeToNextAction = time.Duration((cat.CurrentFerociousBiteCost() - curEnergy) * float64(core.EnergyTickDuration))
	} else if rakeNow {
		if cat.Rake.CanCast(sim, cat.CurrentTarget) {
			cat.Rake.Cast(sim, cat.CurrentTarget)
			return
		}
		timeToNextAction = time.Duration((cat.CurrentRakeCost() - curEnergy) * float64(core.EnergyTickDuration))
	} else if mangleNow {
		if cat.MangleCat.CanCast(sim, cat.CurrentTarget) {
			cat.MangleCat.Cast(sim, cat.CurrentTarget)
			return
		}
		timeToNextAction = time.Duration((cat.CurrentMangleCatCost() - curEnergy) * float64(core.EnergyTickDuration))
	} else if bearweaveNow {
		cat.readyToShift = true
	} else if flowershiftNow && curEnergy < 42 {
		cat.readyToGift = true
	} else if (rotation.MangleSpam && !isClearcast) || cat.PseudoStats.InFrontOfTarget {
		if cat.MangleCat != nil && excessE >= cat.CurrentMangleCatCost() {
			cat.MangleCat.Cast(sim, cat.CurrentTarget)
			return
		}
		timeToNextAction = time.Duration((cat.CurrentMangleCatCost() - excessE) * float64(core.EnergyTickDuration))
	} else {
		if excessE >= cat.CurrentShredCost() || isClearcast {
			cat.Shred.Cast(sim, cat.CurrentTarget)
			return
		}
		// Also Shred if we're about to cap on Energy. Catches some edge
		// cases where floating_energy > 100 due to too many synced timers.
		if curEnergy > 100-(10*latencySecs) {
			cat.Shred.Cast(sim, cat.CurrentTarget)
			return
		}

		timeToNextAction = time.Duration((cat.CurrentShredCost() - excessE) * float64(core.EnergyTickDuration))

		// When Lacerateweaving, there are scenarios where Lacerate is
		// synced with other pending actions. When this happens, pooling for
		// the pending action will inevitably lead to capping on Energy,
		// since we will be forced to shift into Dire Bear Form immediately
		// after pooling in order to save the Lacerate. Instead, it is
		// preferable to just Shred and bearweave early.
		nextCastEnd := sim.CurrentTime + timeToNextAction + cat.latency + time.Second*2
		ignorePooling := cat.BerserkAura.IsActive() || (rotation.BearweaveType == proto.FeralDruid_Rotation_Lacerate && lacerateDot.IsActive() && (lacerateDot.ExpiresAt().Seconds()-1.5-latencySecs <= nextCastEnd.Seconds()))

		if ignorePooling {
			if curEnergy >= cat.CurrentShredCost() {
				cat.Shred.Cast(sim, cat.CurrentTarget)
				return
			}
			timeToNextAction = time.Duration((cat.CurrentShredCost() - curEnergy) * float64(core.EnergyTickDuration))
		}
	}

	// Model in latency when waiting on Energy for our next action
	nextAction := sim.CurrentTime + timeToNextAction
	if len(pendingActions) > 0 {
		nextAction = core.MinDuration(nextAction, pendingActions[0].refreshTime)
	}

	// Also schedule an action right at Energy cap to make sure we never
	// accidentally over-cap while waiting on other timers.
	timeToCap := time.Duration(((100.0 - curEnergy) / 10.0) * float64(time.Second))
	nextAction = core.MinDuration(nextAction, sim.CurrentTime+timeToCap)

	// If Lacerateweaving, then also schedule an action just before Lacerate
	// expires to ensure we can save it in time.
	lacRefreshTime := lacerateDot.ExpiresAt() - (1500 * time.Millisecond) - (3 * cat.latency * time.Second)
	if rotation.BearweaveType == proto.FeralDruid_Rotation_Lacerate && lacerateDot.IsActive() && lacerateDot.RemainingDuration(sim) < sim.GetRemainingDuration() && (sim.CurrentTime < lacRefreshTime) {
		nextAction = core.MinDuration(nextAction, lacRefreshTime)
	}
	nextAction += cat.latency

	if nextAction <= sim.CurrentTime {
		panic("nextaction in the past")
	} else {
		cat.NextRotationAction(sim, nextAction)
	}
}

type pendingAction struct {
	refreshTime time.Duration
	cost        float64
}

type FeralDruidRotation struct {
	BearweaveType      proto.FeralDruid_Rotation_BearweaveType
	MaintainFaerieFire bool
	MinCombosForRip    int32
	UseRake            bool
	UseBite            bool
	BiteTime           time.Duration
	MinCombosForBite   int32
	MangleSpam         bool
	BerserkBiteThresh  float64
	Powerbear          bool
	MinRoarOffset      time.Duration
	RevitFreq          float64
	LacerateTime       time.Duration
	SnekWeave          bool
	FlowerWeave        bool
}

func (cat *FeralDruid) setupRotation(rotation *proto.FeralDruid_Rotation) {
	cat.Rotation = FeralDruidRotation{
		BearweaveType:      rotation.BearWeaveType,
		MaintainFaerieFire: rotation.MaintainFaerieFire,
		MinCombosForRip:    core.Ternary(rotation.MinCombosForRip > 0, rotation.MinCombosForRip, 1),
		UseRake:            rotation.UseRake,
		UseBite:            rotation.UseBite,
		BiteTime:           time.Duration(float64(rotation.BiteTime) * float64(time.Second)),
		MinCombosForBite:   core.Ternary(rotation.MinCombosForBite > 0, rotation.MinCombosForBite, 1),
		MangleSpam:         rotation.MangleSpam,
		BerserkBiteThresh:  float64(rotation.BerserkBiteThresh),
		Powerbear:          rotation.Powerbear,
		MinRoarOffset:      time.Duration(float64(rotation.MinRoarOffset) * float64(time.Second)),
		RevitFreq:          15.0 / (8 * float64(rotation.HotUptime)),
		LacerateTime:       8.0 * time.Second,
		SnekWeave:          core.Ternary(rotation.BearWeaveType == proto.FeralDruid_Rotation_None, false, rotation.SnekWeave),
		FlowerWeave:        core.Ternary(rotation.BearWeaveType == proto.FeralDruid_Rotation_None, rotation.FlowerWeave, false),
	}

	// Use automatic values unless specified
	if rotation.ManualParams {
		return
	}

	hasT72P := cat.HasSetBonus(druid.ItemSetDreamwalkerBattlegear, 2)
	hasT84P := cat.HasSetBonus(druid.ItemSetNightsongBattlegear, 4)

	cat.Rotation.UseRake = true
	cat.Rotation.UseBite = true

	if cat.Rotation.FlowerWeave {
		if hasT84P {
			cat.Rotation.MinRoarOffset = 26 * time.Second
		} else {
			cat.Rotation.MinRoarOffset = 13 * time.Second
		}
		cat.Rotation.BiteTime = 4 * time.Second
	} else {
		if hasT72P {
			cat.Rotation.MinRoarOffset = 14 * time.Second
		} else if hasT84P {
			cat.Rotation.MinRoarOffset = 27 * time.Second
		} else {
			cat.Rotation.MinRoarOffset = 12 * time.Second
		}
		cat.Rotation.BiteTime = 10 * time.Second
	}
}
