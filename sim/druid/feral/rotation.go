package feral

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/druid"
)

type FeralDruidRotation struct {
	// RotationType proto.FeralDruid_Rotation_AplType

	MinCombosForRip    int32
	MaxWaitTime        time.Duration
	MaintainFaerieFire bool
	UseShredTrick      bool
	UseRipTrick        bool
}

func (cat *FeralDruid) OnGCDReady(sim *core.Simulation) {
	if !cat.GCD.IsReady(sim) {
		return
	}

	cat.missChance = cat.MissChance()
	cat.bleedAura = cat.CurrentTarget.GetExclusiveEffectCategory(core.BleedEffectCategory).GetActiveAura()

	if cat.preRotationCleanup(sim) {
		waiting := false
		nextAction := time.Duration(0)
		// if cat.Rotation.RotationType == proto.FeralDruid_Rotation_SingleTarget {
		// 	valid, nextAction = cat.doRotation(sim)
		// } else {
		// 	valid, nextAction = cat.doAoeRotation(sim)
		// }
		waiting, nextAction = cat.doRotation(sim)
		if waiting {
			cat.postRotation(sim, nextAction)
		}
	}

	// Replace gcd event with our own if we casted a spell
	if !cat.GCD.IsReady(sim) {
		nextGcd := cat.NextGCDAt()
		cat.DoNothing()
		cat.CancelGCDTimer(sim)

		cat.NextRotationAction(sim, nextGcd)
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

func (cat *FeralDruid) tryPowershift(sim *core.Simulation) {
	if cat.InForm(druid.Cat) {
		cat.ClearForm(sim)
		cat.TryUseCooldowns(sim)
	}

	if cat.Rotation.MaintainFaerieFire && cat.ShouldFaerieFire(sim, cat.CurrentTarget) && (cat.CurrentMana() >= cat.CatForm.DefaultCast.Cost+cat.FaerieFire.DefaultCast.Cost) {
		cat.FaerieFire.Cast(sim, cat.CurrentTarget)
	} else {
		cat.CatForm.Cast(sim, nil)
		cat.readyToShift = false
		cat.lastShift = sim.CurrentTime
	}
}

func (cat *FeralDruid) maxShifts() int32 {
	return int32(cat.MaxMana() / cat.CatForm.DefaultCast.Cost)
}

func (cat *FeralDruid) numShiftsRemaining() int32 {
	return int32(cat.CurrentMana() / cat.CatForm.DefaultCast.Cost)
}

func (cat *FeralDruid) timeToCast(numSpecials int32) time.Duration {
	// Rough calculation, won't be exact! Intended to skew conservatively.
	numPowershiftedSpecials := min(numSpecials, cat.numShiftsRemaining()*2)
	numOomSpecials := numSpecials - numPowershiftedSpecials
	return core.DurationFromSeconds(float64(numPowershiftedSpecials*2 + numOomSpecials*4))
}

func (cat *FeralDruid) canRip(sim *core.Simulation, isTrick bool) bool {
	if cat.Rip.CurDot().IsActive() {
		return false
	}
	// Allow Rip if conservative napkin math estimate says that we can cast the Rip and then build 5 Combo Points in time before the current Savage Roar expires.
	roarDur := cat.SavageRoarAura.RemainingDuration(sim)
	fightDur := sim.GetRemainingDuration()
	remainingFightTimeAfterRoar := fightDur - roarDur

	// solve "remainingFightTimeAfterRoar = 5*roarCP+9" for roarCP
	// add 1 to round up instead of down
	roarCp := int32((remainingFightTimeAfterRoar-time.Second*9)/(time.Second*5)) + 1
	minRoarCp := min(roarCp, 5)

	// Actions to generate minRoarCp, plus cast Roar itself.
	actionsToCastRoar := minRoarCp + 1

	// Don't let roar expire.
	if cat.timeToCast(actionsToCastRoar) >= roarDur {
		return false
	}

	// Don't rip if it won't be able to tick for the full duration.
	if fightDur <= time.Second*10 {
		return false
	}

	if cat.ComboPoints() == 5 {
		return true
	}

	// If we can't get any more combo points before roar expires, then we should rip now.
	// If we can generate another CP and then rip without letting roar expire, then wait.
	if cat.timeToCast(actionsToCastRoar+1) >= roarDur {
		return true
	}

	// Caller decided that we should "rip trick".
	if isTrick {
		return true
	}

	return false
}

/*
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
	shredDpc := cat.Shred.ExpectedInitialDamage(sim, cat.CurrentTarget)
	potentialRakeTicks := min(cat.Rake.CurDot().NumberOfTicks, int32(sim.GetRemainingDuration()/time.Second*3))
	rakeDpc := cat.Rake.ExpectedInitialDamage(sim, cat.CurrentTarget) + cat.Rake.ExpectedTickDamage(sim, cat.CurrentTarget)*float64(potentialRakeTicks)
	return rakeDpc / cat.Rake.DefaultCast.Cost, shredDpc / cat.Shred.DefaultCast.Cost
}
*/

func (cat *FeralDruid) clipRoar(sim *core.Simulation) bool {
	// If existing Roar already covers us to end of fight, then don't clip it
	roarDur := cat.SavageRoarAura.RemainingDuration(sim)
	fightDur := sim.GetRemainingDuration()

	if roarDur >= fightDur {
		return false
	}

	// If a fresh Roar cast now at the current number of Combo Points *would* cover us to end of fight, then clip now for maximum CP efficiency
	newRoarDur := cat.SavageRoarDurationTable[cat.ComboPoints()]

	if newRoarDur >= fightDur {
		return true
	}

	// Roar clips that don't cover us to end of fight should only be done at 5 CP
	if cat.ComboPoints() < 5 {
		return false
	}

	// Calculate the minimum number of Roar casts that will cover us to end of fight if we (a) let the current one expire naturally vs. (b) clip it now.
	minRoarsPossible := (fightDur - roarDur) / newRoarDur
	projectedRoarCasts := fightDur / newRoarDur

	// Allow a clip at the earliest time that doesn't result in an extra Roar cast
	return projectedRoarCasts == minRoarsPossible
}

func (cat *FeralDruid) preRotationCleanup(sim *core.Simulation) bool {
	// If we previously decided to shift, then execute the shift now once
	// the input delay is over.
	if cat.readyToShift {
		cat.tryPowershift(sim)
		return false
	}

	return true
}

func (cat *FeralDruid) postRotation(sim *core.Simulation, nextAction time.Duration) {
	nextAction += cat.latency

	if nextAction <= sim.CurrentTime {
		panic("nextaction in the past")
	} else {
		cat.NextRotationAction(sim, nextAction)
	}
}

func (cat *FeralDruid) shouldPoolMana(sim *core.Simulation, numShiftsToOom int32) bool {
	if cat.Talents.Furor == 0 {
		return true
	}

	effectiveFightDur := sim.GetRemainingDuration() - core.DurationFromSeconds(3.0) - cat.latency
	numShiftsToFightEnd := int32(effectiveFightDur / (time.Second * 4))
	canPoolMana := (numShiftsToOom < cat.maxShifts()) && (numShiftsToOom < numShiftsToFightEnd-1) && (sim.CurrentTime-cat.lastShift > time.Second*5)

	if !canPoolMana {
		return false
	}

	for _, cd := range cat.GetMajorCooldowns() {
		if cd.IsEnabled() && cd.Type.Matches(core.CooldownTypeMana) && cd.IsReady(sim) {
			return false
		}
	}

	return true
}

func (cat *FeralDruid) doRotation(sim *core.Simulation) (bool, time.Duration) {
	// Store state variables that will be used in calculations
	rotation := &cat.Rotation
	curCp := cat.ComboPoints()
	curEnergy := cat.CurrentEnergy()
	nextEnergy := curEnergy + core.EnergyPerTick
	nextTick := cat.NextEnergyTickAt()
	timeToNextTick := nextTick - sim.CurrentTime
	isClearcast := cat.ClearcastingAura.IsActive()
	hasRoar := cat.SavageRoarAura.IsActive()
	numShiftsToOom := cat.numShiftsRemaining()
	fightDur := sim.GetRemainingDuration()
	shredCost := cat.Shred.DefaultCast.Cost
	mangleCost := cat.MangleCat.DefaultCast.Cost
	ripCost := cat.Rip.DefaultCast.Cost

	// First determine the next special ability we want to cast
	poolMana := cat.shouldPoolMana(sim, numShiftsToOom)
	canShredTrick := rotation.UseShredTrick && cat.bleedAura.IsActive() && (curEnergy >= shredCost) && (timeToNextTick > time.Second) && ((nextEnergy-shredCost >= mangleCost) || (timeToNextTick > core.GCDDefault)) && (numShiftsToOom > 1) && !poolMana
	canRipTrick := rotation.UseRipTrick && (curCp >= 1) && ((curEnergy >= ripCost && curEnergy < mangleCost) || (nextEnergy >= ripCost && nextEnergy < mangleCost)) && (numShiftsToOom > 0) && !poolMana

	var nextAbility *druid.DruidSpell

	if curCp >= 1 && !hasRoar {
		nextAbility = cat.SavageRoar
	} else if isClearcast {
		nextAbility = cat.Shred
	} else if curCp >= 1 && cat.clipRoar(sim) {
		nextAbility = cat.SavageRoar
	} else if (curCp >= rotation.MinCombosForRip || canRipTrick) && cat.canRip(sim, canRipTrick) {
		nextAbility = cat.Rip
	} else if canShredTrick {
		nextAbility = cat.Shred
	} else {
		nextAbility = cat.MangleCat
	}

	// Then determine whether to cast vs. wait vs. shift
	waitForWildStrikesProc := (cat.WildStrikesBuffAura != nil) && !cat.WildStrikesBuffAura.IsActive()
	poolEnergy := poolMana && ((curCp == 5) || waitForWildStrikesProc) && (nextEnergy < 100) && (nextAbility == cat.MangleCat)
	nextAction := sim.CurrentTime

	if nextAbility.CanCast(sim, cat.CurrentTarget) && !poolEnergy {
		nextAbility.Cast(sim, cat.CurrentTarget)
		return false, nextAction
	}

	shiftNow := ((nextEnergy < nextAbility.DefaultCast.Cost) || (timeToNextTick > rotation.MaxWaitTime)) && (fightDur > core.GCDDefault+cat.latency)

	if shiftNow && (poolMana || (numShiftsToOom == 0)) {
		shiftNow = false
		cat.Metrics.MarkOOM(sim)

		if !cat.poolingMana {
			cat.poolingMana = true
			cat.poolStartTime = sim.CurrentTime
		}
	}

	if shiftNow {
		cat.readyToShift = true

		if cat.poolingMana {
			cat.poolingMana = false
			cat.Metrics.AddOOMTime(sim, sim.CurrentTime-cat.poolStartTime)
		}
	} else {
		nextAction = nextTick
	}

	return true, nextAction
}

/*
func (cat *FeralDruid) doRotation(sim *core.Simulation) (bool, time.Duration) {
	rotation := &cat.Rotation

	curEnergy := cat.CurrentEnergy()
	curRage := cat.CurrentRage()
	curCp := cat.ComboPoints()
	isClearcast := cat.ClearcastingAura.IsActive()
	simTimeRemain := sim.GetRemainingDuration()
	shiftCost := cat.CatForm.DefaultCast.Cost
	rakeDot := cat.Rake.CurDot()
	ripDot := cat.Rip.CurDot()
	lacerateDot := cat.Lacerate.CurDot()
	isBleedActive := cat.AssumeBleedActive || ripDot.IsActive() || rakeDot.IsActive() || lacerateDot.IsActive()

	// Prioritize using rake/rip with omen procs if bleed isnt active
	// But less priority then mangle aura
	ripCcCheck := core.Ternary(isBleedActive, !isClearcast, true)
	rakeCcCheck := core.Ternary(isBleedActive, !isClearcast, cat.bleedAura.IsActive())

	endThresh := time.Second * 10

	ripNow := (curCp >= rotation.MinCombosForRip) && !ripDot.IsActive() && (simTimeRemain >= endThresh) && ripCcCheck
	biteAtEnd := (curCp >= rotation.MinCombosForBite) && ((simTimeRemain < endThresh) || (ripDot.IsActive() && (simTimeRemain-ripDot.RemainingDuration(sim) < endThresh)))

	// Clip Mangle if it won't change the total number of Mangles we have to
	// cast before the fight ends.
	mangleRefreshNow := !cat.bleedAura.IsActive() && simTimeRemain > time.Second
	mangleRefreshPending := cat.bleedAura.IsActive() && cat.bleedAura.RemainingDuration(sim) < (simTimeRemain-time.Second)
	clipMangle := false

	if mangleRefreshPending {
		numManglesRemaining := int((time.Second + (simTimeRemain - cat.bleedAura.RemainingDuration(sim) - time.Second)).Minutes())
		earliestMangle := (sim.GetRemainingDuration() + sim.CurrentTime) - time.Duration(numManglesRemaining)*time.Minute
		clipMangle = sim.CurrentTime >= earliestMangle
	}

	mangleNow := !ripNow && cat.MangleCat != nil && (mangleRefreshNow || clipMangle)

	biteBeforeRip := (curCp >= rotation.MinCombosForBite) && ripDot.IsActive() && cat.SavageRoarAura.IsActive() && rotation.UseBite && cat.canBite(sim)
	biteNow := (biteBeforeRip || biteAtEnd) && !isClearcast && curEnergy < 67

	// During Berserk, we additionally add an Energy constraint on Bite
	// usage to maximize the total Energy expenditure we can get.
	if biteNow && cat.BerserkAura.IsActive() {
		biteNow = curEnergy <= rotation.BerserkBiteThresh
	}

	rakeNow := rotation.UseRake && !rakeDot.IsActive() && (simTimeRemain > rakeDot.Duration) && rakeCcCheck

	// Additionally, don't Rake if the current Shred DPE is higher due to
	// trinket procs etc.
	if rotation.RakeDpeCheck && rakeNow {
		rakeDpe, shredDpe := cat.calcBuilderDpe(sim)
		rakeNow = (rakeDpe > shredDpe)
	}

	// Additionally, don't Rake if there is insufficient time to max out
	// our available glyph of shred extensions before rip falls off
	if rakeNow && ripDot.IsActive() {
		maxRipDur := time.Duration(cat.maxRipTicks) * ripDot.TickLength
		remainingExt := cat.maxRipTicks - ripDot.NumberOfTicks
		energyForShreds := curEnergy - cat.CurrentRakeCost() - 30 + float64((ripDot.StartedAt()+maxRipDur-sim.CurrentTime)/core.EnergyTickDuration) + core.Ternary(cat.tfExpectedBefore(sim, ripDot.StartedAt()+maxRipDur), 60.0, 0.0)
		maxShredsPossible := min(energyForShreds/cat.Shred.DefaultCast.Cost, (ripDot.ExpiresAt() - (sim.CurrentTime + time.Second)).Seconds())
		rakeNow = remainingExt == 0 || (maxShredsPossible > float64(remainingExt))
	}

	// Disable Energy pooling for Rake in weaving rotations, since these
	// rotations prioritize weave cpm over Rake uptime.
	poolForRake := !(rotation.FlowerWeave || (rotation.BearweaveType != proto.FeralDruid_Rotation_None))

	// Berserk algorithm: time Berserk for just after a Tiger's Fury
	// *unless* we'll lose Berserk uptime by waiting for Tiger's Fury to
	// come off cooldown. The latter exception is necessary for
	// Lacerateweave rotation since TF timings can drift over time.
	waitForTf := cat.Talents.Berserk && (cat.TigersFury.ReadyAt() <= cat.BerserkAura.Duration) && (cat.TigersFury.ReadyAt()+time.Second < simTimeRemain-cat.BerserkAura.Duration)
	berserkNow := cat.Berserk.IsReady(sim) && !waitForTf && ripDot.IsActive() && !isClearcast

	// Additionally, for Lacerateweave rotation, postpone the final Berserk
	// of the fight to as late as possible so as to minimize the impact of
	// dropping Lacerate stacks during the Berserk window. Rationale for the
	// 3 second additional leeway given beyond just berserk_dur in the below
	// expression is to be able to fit in a final TF and dump the Energy
	// from it in cases where Berserk and TF CDs are desynced due to drift.
	if berserkNow && rotation.BearweaveType == proto.FeralDruid_Rotation_Lacerate && cat.berserkUsed && simTimeRemain < cat.Berserk.CD.Duration {
		berserkNow = simTimeRemain < cat.BerserkAura.Duration+(3*time.Second)
	}

	roarNow := curCp >= 1 && (!cat.SavageRoarAura.IsActive() || cat.clipRoar(sim))

	// Faerie Fire on cooldown for Omen procs. Each second of FF delay is
	// worth ~7 Energy, so it is okay to waste up to 7 Energy to cap when
	// determining whether to cast it vs. dump Energy first. That puts the
	// Energy threshold for FF usage as 107 minus 10 for the Clearcasted
	// special minus 10 for the FF GCD = 87 Energy.
	ffThresh := 87.0
	if cat.BerserkAura.IsActive() {
		ffThresh = cat.Rotation.BerserkFfThresh
	}
	ffNow := cat.FaerieFire.CanCast(sim, cat.CurrentTarget) && !isClearcast && curEnergy < ffThresh && (!ripNow || (curEnergy < cat.CurrentRipCost()))

	// Also add an end of fight condition to make sure we can spend down our
	// Energy post-FF before the encounter ends. Time to spend is
	// given by 1 second for FF GCD plus 1 second for Clearcast Shred plus
	// 1 second per 42 Energy that we have after that Clearcast Shred.
	if ffNow {
		simTimeSecs := sim.GetRemainingDuration().Seconds()
		maxShredsWithoutFF := (int)((curEnergy + simTimeSecs*10) / cat.Shred.DefaultCast.Cost)
		numShredsWithoutFF := min(maxShredsWithoutFF, int(simTimeSecs)+1)
		numShredsWithFF := min(maxShredsWithoutFF+1, int(simTimeSecs))
		ffNow = numShredsWithFF > numShredsWithoutFF
	}

	// Additionally, block Shred and Rake casts if FF is coming off CD in
	// less than a second (and we won't Energy cap by pooling).
	nextFfEnergy := curEnergy + float64((cat.FaerieFire.TimeToReady(sim)+cat.latency)/core.EnergyTickDuration)
	waitForFf := (cat.FaerieFire.TimeToReady(sim) < time.Second-cat.Rotation.MaxFfDelay) && (nextFfEnergy < ffThresh) && !isClearcast && (!ripDot.IsActive() || ripDot.RemainingDuration(sim) > time.Second)

	cat.ripRefreshPending = false

	pendingPool := PoolingActions{}
	pendingPool.create(4)

	if ripDot.IsActive() && (ripDot.RemainingDuration(sim) < simTimeRemain-endThresh) && curCp == 5 {
		ripCost := core.Ternary(cat.berserkExpectedAt(sim, ripDot.ExpiresAt()), cat.Rip.DefaultCast.Cost*0.5, cat.Rip.DefaultCast.Cost)
		pendingPool.addAction(ripDot.ExpiresAt(), ripCost)
		cat.ripRefreshPending = true
	}
	if poolForRake && rakeDot.IsActive() && (rakeDot.RemainingDuration(sim) < simTimeRemain-rakeDot.Duration) {
		rakeCost := core.Ternary(cat.berserkExpectedAt(sim, rakeDot.ExpiresAt()), cat.Rake.DefaultCast.Cost*0.5, cat.Rake.DefaultCast.Cost)
		pendingPool.addAction(rakeDot.ExpiresAt(), rakeCost)
	}
	if mangleRefreshPending {
		mangleCost := core.Ternary(cat.berserkExpectedAt(sim, cat.bleedAura.ExpiresAt()), cat.MangleCat.DefaultCast.Cost*0.5, cat.MangleCat.DefaultCast.Cost)
		pendingPool.addAction(cat.bleedAura.ExpiresAt(), mangleCost)
	}
	if cat.SavageRoarAura.IsActive() {
		roarCost := core.Ternary(cat.berserkExpectedAt(sim, cat.SavageRoarAura.ExpiresAt()), cat.SavageRoar.DefaultCast.Cost*0.5, cat.SavageRoar.DefaultCast.Cost)
		pendingPool.addAction(cat.SavageRoarAura.ExpiresAt(), roarCost)
	}

	pendingPool.sort()

	latencySecs := cat.latency.Seconds()
	// Allow for bearweaving if the next pending action is >= 4.5s away
	furorCap := min(20.0*float64(cat.Talents.Furor), 85)
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
	bearweaveNow := rotation.BearweaveType != proto.FeralDruid_Rotation_None && curEnergy <= weaveEnergy && !isClearcast && (!cat.ripRefreshPending || ripDot.ExpiresAt() >= weaveEnd) && !cat.BerserkAura.IsActive()

	if bearweaveNow && rotation.BearweaveType != proto.FeralDruid_Rotation_Lacerate {
		bearweaveNow = !cat.tfExpectedBefore(sim, weaveEnd)
	}

	// Also add an end of fight condition to make sure we can spend down our
	// Energy post-bearweave before the encounter ends. Time to spend is
	// given by weave_end plus 1 second per 42 Energy that we have at
	// weave_end.
	if bearweaveNow {
		energyToDump := curEnergy + ((weaveEnd - sim.CurrentTime).Seconds() * 10)
		bearweaveNow = weaveEnd+time.Duration(math.Floor(energyToDump/42)*float64(time.Second)) < sim.CurrentTime+simTimeRemain
	}

	// If we're maintaining Lacerate, then allow for emergency bearweaves
	// if Lacerate is about to fall off even if the above conditions do not
	// apply.
	lacRemain := core.Ternary(lacerateDot.IsActive(), lacerateDot.RemainingDuration(sim), time.Duration(0))
	emergencyBearweave := rotation.BearweaveType == proto.FeralDruid_Rotation_Lacerate && lacerateDot.IsActive() && (float64(lacRemain) < (2.5+latencySecs)*float64(time.Second)) && (lacRemain < simTimeRemain) && !cat.BerserkAura.IsActive()

	// As an alternative to bearweaving, cast GotW on the raid under
	// analagous conditions to the above. Only difference is that there is
	// more available time/Energy leeway for the technique, since
	// flowershifts take only 3 seconds to execute.
	flowershiftEnergy := min(furorCap, 75) - 10*cat.SpellGCD().Seconds() - 20*latencySecs

	flowerEnd := time.Duration(float64(sim.CurrentTime) + float64(cat.SpellGCD()) + (2.5+2*latencySecs)*float64(time.Second))
	flowerFfDelay := flowerEnd - cat.FaerieFire.ReadyAt()
	flowershiftNow := rotation.FlowerWeave && (curEnergy <= flowershiftEnergy) && !isClearcast && (!cat.ripRefreshPending || ripDot.ExpiresAt() >= flowerEnd) && !cat.BerserkAura.IsActive() && !cat.tfExpectedBefore(sim, flowerEnd) && flowerFfDelay < rotation.MaxFfDelay

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
		energyToDump := curEnergy + ((flowerEnd - sim.CurrentTime).Seconds() * 10)
		flowershiftNow = flowerEnd+time.Duration(math.Floor(energyToDump/42)*float64(time.Second)) < sim.CurrentTime+simTimeRemain
	}

	floatingEnergy := pendingPool.calcFloatingEnergy(cat, sim)
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
		shiftNow := (curEnergy+15.0+(10.0*latencySecs) > furorCap) || (cat.ripRefreshPending && (ripDot.RemainingDuration(sim) < (3.0 * time.Second))) || cat.BerserkAura.IsActive()
		shiftNext := (curEnergy+30.0+(10.0*latencySecs) > furorCap) || (cat.ripRefreshPending && (ripDot.RemainingDuration(sim) < time.Duration(4500*time.Millisecond))) || cat.BerserkAura.IsActive()

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
			timeToDump := (3 * time.Second) + cat.latency + time.Duration(math.Floor(energyToDump/42)*float64(time.Second))
			shiftNow = timeToDump >= simTimeRemain
		}

		nextSwing := cat.AutoAttacks.NextAttackAt()

		if emergencyLacerate && cat.Lacerate.CanCast(sim, cat.CurrentTarget) {
			cat.Lacerate.Cast(sim, cat.CurrentTarget)
			return false, 0
		} else if shiftNow {
			// If we are resetting our swing timer using Albino Snake or a
			// duplicate weapon swap, then do an additional check here to
			// see whether we can delay the shift until the next bear swing
			// goes out in order to maximize the gains from the reset.
			projectedDelay := nextSwing + 2*cat.latency - sim.CurrentTime
			ripConflict := cat.ripRefreshPending && (ripDot.ExpiresAt() < sim.CurrentTime+projectedDelay+(1500*time.Millisecond))
			nextCatSwing := sim.CurrentTime + cat.latency + time.Duration(float64(cat.AutoAttacks.MainhandSwingSpeed())/float64(2500*time.Millisecond))
			canDelayShift := !ripConflict && cat.Rotation.SnekWeave && (curEnergy+10*projectedDelay.Seconds() <= furorCap) && (nextSwing < nextCatSwing)

			if canDelayShift {
				timeToNextAction = nextSwing - sim.CurrentTime
			} else {
				cat.readyToShift = true
			}
		} else if powerbearNow {
			cat.shiftBearCat(sim, true)
		} else if lacerateNow && cat.Lacerate.CanCast(sim, cat.CurrentTarget) {
			cat.Lacerate.Cast(sim, cat.CurrentTarget)
			return false, 0
		} else if cat.MangleBear.CanCast(sim, cat.CurrentTarget) {
			cat.MangleBear.Cast(sim, cat.CurrentTarget)
			return false, 0
		} else if cat.Lacerate.CanCast(sim, cat.CurrentTarget) {
			cat.Lacerate.Cast(sim, cat.CurrentTarget)
			return false, 0
		} else {
			timeToNextAction = nextSwing - sim.CurrentTime
		}
	} else if emergencyBearweave {
		cat.readyToShift = true
	} else if ffNow {
		cat.FaerieFire.Cast(sim, cat.CurrentTarget)
		return false, 0
	} else if berserkNow {
		cat.Berserk.Cast(sim, nil)
		cat.UpdateMajorCooldowns()
		return false, 0
	} else if roarNow {
		if cat.SavageRoar.CanCast(sim, cat.CurrentTarget) {
			cat.SavageRoar.Cast(sim, nil)
			return false, 0
		}
		timeToNextAction = time.Duration((cat.CurrentSavageRoarCost() - curEnergy) * float64(core.EnergyTickDuration))
	} else if ripNow {
		if cat.Rip.CanCast(sim, cat.CurrentTarget) {
			cat.Rip.Cast(sim, cat.CurrentTarget)
			return false, 0
		}
		timeToNextAction = time.Duration((cat.CurrentRipCost() - curEnergy) * float64(core.EnergyTickDuration))
	} else if biteNow {
		if cat.FerociousBite.CanCast(sim, cat.CurrentTarget) {
			cat.FerociousBite.Cast(sim, cat.CurrentTarget)
			return false, 0
		}
		timeToNextAction = time.Duration((cat.CurrentFerociousBiteCost() - curEnergy) * float64(core.EnergyTickDuration))
	} else if mangleNow && !waitForFf {
		if cat.MangleCat.CanCast(sim, cat.CurrentTarget) {
			cat.MangleCat.Cast(sim, cat.CurrentTarget)
			return false, 0
		}
		timeToNextAction = time.Duration((cat.CurrentMangleCatCost() - curEnergy) * float64(core.EnergyTickDuration))
	} else if rakeNow && !waitForFf {
		if cat.Rake.CanCast(sim, cat.CurrentTarget) {
			cat.Rake.Cast(sim, cat.CurrentTarget)
			return false, 0
		}
		timeToNextAction = time.Duration((cat.CurrentRakeCost() - curEnergy) * float64(core.EnergyTickDuration))
	} else if bearweaveNow {
		cat.readyToShift = true
	} else if flowershiftNow && curEnergy < 42 {
		cat.readyToGift = true
	} else if (rotation.MangleSpam && !isClearcast) || cat.PseudoStats.InFrontOfTarget {
		if cat.MangleCat != nil && excessE >= cat.CurrentMangleCatCost() {
			cat.MangleCat.Cast(sim, cat.CurrentTarget)
			return false, 0
		}
		timeToNextAction = time.Duration((cat.CurrentMangleCatCost() - excessE) * float64(core.EnergyTickDuration))
	} else if !waitForFf {
		if excessE >= cat.CurrentShredCost() || isClearcast {
			cat.Shred.Cast(sim, cat.CurrentTarget)
			return false, 0
		}
		// Also Shred if we're about to cap on Energy. Catches some edge
		// cases where floating_energy > 100 due to too many synced timers.
		if curEnergy > 100-(10*latencySecs) {
			cat.Shred.Cast(sim, cat.CurrentTarget)
			return false, 0
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
				return false, 0
			}
			timeToNextAction = time.Duration((cat.CurrentShredCost() - curEnergy) * float64(core.EnergyTickDuration))
		}
	}

	// Model in latency when waiting on Energy for our next action
	nextAction := sim.CurrentTime + timeToNextAction
	paValid, rt := pendingPool.nextRefreshTime()
	if paValid {
		nextAction = min(nextAction, rt)
	}

	// If Lacerateweaving, then also schedule an action just before Lacerate
	// expires to ensure we can save it in time.
	lacRefreshTime := lacerateDot.ExpiresAt() - (1500 * time.Millisecond) - (3 * cat.latency)
	if rotation.BearweaveType == proto.FeralDruid_Rotation_Lacerate && lacerateDot.IsActive() && lacerateDot.RemainingDuration(sim) < simTimeRemain && (sim.CurrentTime < lacRefreshTime) {
		nextAction = min(nextAction, lacRefreshTime)
	}

	return true, nextAction
}

type FeralDruidRotation struct {
	RotationType proto.FeralDruid_Rotation_AplType

	BearweaveType      proto.FeralDruid_Rotation_BearweaveType
	MaintainFaerieFire bool
	MinCombosForRip    int32
	UseRake            bool
	UseBite            bool
	BiteTime           time.Duration
	MinCombosForBite   int32
	MangleSpam         bool
	BerserkBiteThresh  float64
	BerserkFfThresh    float64
	Powerbear          bool
	MinRoarOffset      time.Duration
	RipLeeway          time.Duration
	MaxFfDelay         time.Duration
	RevitFreq          float64
	LacerateTime       time.Duration
	SnekWeave          bool
	FlowerWeave        bool
	RakeDpeCheck       bool

	AoeMangleBuilder bool
}

*/

func (cat *FeralDruid) setupRotation(config *proto.APLActionCatOptimalRotationAction) {
	cat.Rotation = FeralDruidRotation{
		MinCombosForRip:    config.MinCombosForRip,
		MaxWaitTime:        core.DurationFromSeconds(float64(config.MaxWaitTime)),
		MaintainFaerieFire: config.MaintainFaerieFire,
		UseShredTrick:      config.UseShredTrick,
		UseRipTrick:        false,
	}
}
