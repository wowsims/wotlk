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
	if cat.InForm(druid.Cat) {
		cat.doTigersFury(sim)
	}
}

func (cat *FeralDruid) OnGCDReady(sim *core.Simulation) {
	cat.TryUseCooldowns(sim)
	cat.doRotation(sim)
}

func (cat *FeralDruid) OnAutoAttack(sim *core.Simulation, spell *core.Spell) {
	if cat.InForm(druid.Cat) {
		return
	}
	if cat.InForm(druid.Humanoid) {
		panic("auto attack out of form?")
	}

	cat.checkQueueMaul(sim)
}

// Ported from https://github.com/NerdEgghead/WOTLK_cat_sim

func (cat *FeralDruid) checkQueueMaul(sim *core.Simulation) {
	// If we will have enough time and Energy leeway to stay in
	// Dire Bear Form once the GCD expires, then only Maul if we
	// will be left with enough Rage to cast Mangle or Lacerate
	// on that global.

	furorCap := core.MinFloat(20.0*float64(cat.Talents.Furor), 85.0)
	ripRefreshPending := cat.RipDot.IsActive() && (cat.RipDot.RemainingDuration(sim) < sim.GetRemainingDuration()-time.Second*10)
	gcdTimeToRdy := cat.GCD.TimeToReady(sim)
	energyLeeway := furorCap - 15.0 - 10.0*(gcdTimeToRdy+cat.latency).Seconds()
	shiftNext := cat.CurrentEnergy() > energyLeeway

	if ripRefreshPending {
		shiftNext = shiftNext || (cat.RipDot.RemainingDuration(sim) < (gcdTimeToRdy + time.Second*3))
	}

	lacerateNext := false
	emergencyLacerateNext := false
	mangleNext := false

	if cat.Rotation.BearweaveType == proto.FeralDruid_Rotation_Lacerate {
		lacerateLeeway := cat.Rotation.LacerateTime + gcdTimeToRdy
		lacerateNext = !cat.LacerateDot.IsActive() || (cat.LacerateDot.GetStacks() < 5) || (cat.LacerateDot.RemainingDuration(sim) <= lacerateLeeway)
		emergencyLeeway := gcdTimeToRdy + (3 * time.Second) + (2 * cat.latency)
		emergencyLacerateNext = cat.LacerateDot.IsActive() && (cat.LacerateDot.RemainingDuration(sim) <= emergencyLeeway)
		mangleNext = cat.MangleBear != nil && !lacerateNext && (!cat.MangleAura.IsActive() || (cat.MangleAura.RemainingDuration(sim) < gcdTimeToRdy+time.Second*3))
	} else {
		mangleNext = cat.MangleBear != nil && cat.MangleBear.TimeToReady(sim) < gcdTimeToRdy
		lacerateNext = cat.LacerateDot.IsActive() && (cat.LacerateDot.GetStacks() < 5 || cat.LacerateDot.RemainingDuration(sim) < gcdTimeToRdy+(time.Second*4))
	}

	maulRageThresh := 10.0
	if emergencyLacerateNext {
		maulRageThresh += cat.Lacerate.BaseCost
	} else if shiftNext {
		maulRageThresh = 10.0
	} else if mangleNext {
		maulRageThresh += cat.MangleBear.BaseCost
	} else if lacerateNext {
		maulRageThresh += cat.Lacerate.BaseCost
	}

	if cat.CurrentRage() >= maulRageThresh {
		cat.QueueMaul(sim)
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

	if !cat.InForm(druid.Cat | druid.Bear) {
		panic("unsupported shift, must be in form")
	}

	toCat := cat.InForm(druid.Bear)
	if powershift {
		toCat = !toCat
	}

	if toCat {
		return cat.PowerShiftCat(sim)
	} else {
		cat.PowerShiftBear(sim)
		// Bundle Enrage if available
		if cat.Enrage.IsReady(sim) {
			cat.Enrage.Cast(sim, nil)
		}
		cat.checkQueueMaul(sim)
		return true
	}
}

func (cat *FeralDruid) canBite(sim *core.Simulation) bool {
	return cat.RipDot.RemainingDuration(sim) >= cat.Rotation.BiteTime &&
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

func (cat *FeralDruid) clipRoar(sim *core.Simulation) bool {
	ripdotRemaining := cat.RipDot.RemainingDuration(sim)
	if !cat.RipDot.IsActive() || (ripdotRemaining < 10*time.Second) {
		return false
	}

	// Project Rip end time assuming full Glyph of Shred extensions
	maxRipDur := time.Duration(cat.maxRipTicks) * cat.RipDot.TickLength
	ripDur := cat.RipDot.Aura.StartedAt() + maxRipDur - sim.CurrentTime
	roarDur := cat.SavageRoarAura.RemainingDuration(sim)

	if roarDur > ripDur {
		return false
	}

	// Calculate when roar would end if casted now
	newRoarDur := cat.SavageRoarDurationTable()[cat.ComboPoints()]

	// Clip as soon as we have enough CPs for the new roar to expire well
	// after the current rip
	return newRoarDur >= (ripDur + cat.Rotation.MaxRoarOffset)
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

	if tfNow {
		cat.TigersFury.Cast(sim, nil)
		// Kick gcd loop, also need to account for any gcd 'left'
		// otherwise it breaks gcd logic
		cat.WaitUntil(sim, sim.CurrentTime+gcdTimeToRdy)
	}
}

func (cat *FeralDruid) doRotation(sim *core.Simulation) {
	// If we previously decided to shift, then execute the shift now once
	// the input delay is over.
	if cat.readyToShift {
		didShift := cat.shiftBearCat(sim, false)
		if !didShift {
			panic("didnt shift?")
		}
		// Reset swing timer from snek (or idol/weapon swap) when going into cat
		if cat.InForm(druid.Cat) && cat.Rotation.SnekWeave {
			cat.AutoAttacks.StopMeleeUntil(sim, sim.CurrentTime)
		}
		return
	}

	rotation := &cat.Rotation

	if rotation.MaintainFaerieFire && cat.ShouldFaerieFire(sim) {
		cat.FaerieFire.Cast(sim, cat.CurrentTarget)
		return
	}

	cat.missChance = cat.MissChance()

	curEnergy := cat.CurrentEnergy()
	curRage := cat.CurrentRage()
	curCp := cat.ComboPoints()
	isClearcast := cat.ClearcastingAura.IsActive()
	simTimeRemain := sim.GetRemainingDuration()
	shiftCost := cat.CatForm.DefaultCast.Cost

	endThresh := time.Second * 10

	ripNow := (curCp >= rotation.MinCombosForRip) && !cat.RipDot.IsActive() && (simTimeRemain >= endThresh) && !isClearcast
	biteAtEnd := (curCp >= rotation.MinCombosForBite) && ((simTimeRemain < endThresh) || (cat.RipDot.IsActive() && (simTimeRemain-cat.RipDot.RemainingDuration(sim) < endThresh)))
	mangleNow := !ripNow && !cat.MangleAura.IsActive() && cat.MangleCat != nil

	biteBeforeRip := (curCp >= rotation.MinCombosForBite) && cat.RipDot.IsActive() && cat.SavageRoarAura.IsActive() && rotation.UseBite && cat.canBite(sim)
	biteNow := (biteBeforeRip || biteAtEnd) && !isClearcast

	// During Berserk, we additionally add an Energy constraint on Bite
	// usage to maximize the total Energy expenditure we can get.
	if biteNow && cat.BerserkAura.IsActive() {
		biteNow = curEnergy <= rotation.BerserkBiteThresh
	}

	rakeNow := rotation.UseRake && !cat.RakeDot.IsActive() && (simTimeRemain > cat.RakeDot.Duration) && !isClearcast

	//berserkEnergyThresh := core.TernaryFloat64(isClearcast, 80.0, 90.0)
	berserkNow := cat.Berserk.IsReady(sim) && (cat.TigersFury.TimeToReady(sim) > cat.BerserkAura.Duration)

	roarNow := curCp >= 1 && (!cat.SavageRoarAura.IsActive() || cat.clipRoar(sim))

	ripRefreshPending := false
	pendingActions := make([]pendingAction, 0, 4)

	if cat.RipDot.IsActive() && (cat.RipDot.RemainingDuration(sim) < simTimeRemain-endThresh) {
		ripCost := core.TernaryFloat64(cat.berserkExpectedAt(sim, cat.RipDot.ExpiresAt()), cat.Rip.BaseCost*0.5, cat.Rip.BaseCost)
		pendingActions = append(pendingActions, pendingAction{cat.RipDot.ExpiresAt(), ripCost})
		ripRefreshPending = true
	}
	if cat.RakeDot.IsActive() && (cat.RakeDot.RemainingDuration(sim) < simTimeRemain-(9*time.Second)) {
		rakeCost := core.TernaryFloat64(cat.berserkExpectedAt(sim, cat.RakeDot.ExpiresAt()), cat.Rake.BaseCost*0.5, cat.Rake.BaseCost)
		pendingActions = append(pendingActions, pendingAction{cat.RakeDot.ExpiresAt(), rakeCost})
	}
	if cat.MangleAura.IsActive() && (cat.MangleAura.RemainingDuration(sim) < simTimeRemain-time.Second) {
		mangleCost := core.TernaryFloat64(cat.berserkExpectedAt(sim, cat.MangleAura.ExpiresAt()), cat.MangleCat.BaseCost*0.5, cat.MangleCat.BaseCost)
		pendingActions = append(pendingActions, pendingAction{cat.MangleAura.ExpiresAt(), mangleCost})
	}
	if cat.SavageRoarAura.IsActive() {
		roarCost := core.TernaryFloat64(cat.berserkExpectedAt(sim, cat.SavageRoarAura.ExpiresAt()), cat.SavageRoar.BaseCost*0.5, cat.SavageRoar.BaseCost)
		pendingActions = append(pendingActions, pendingAction{cat.SavageRoarAura.ExpiresAt(), roarCost})
	}

	slices.SortStableFunc(pendingActions, func(p1, p2 pendingAction) bool {
		return p1.refreshTime < p2.refreshTime
	})

	latencySecs := cat.latency.Seconds()
	// Allow for bearweaving if the next pending action is >= 4.5s away
	furorCap := core.MinFloat(20.0*float64(cat.Talents.Furor), 85)
	weaveEnergy := furorCap - 30 - 20*latencySecs

	if cat.Talents.Furor > 3 {
		weaveEnergy -= 15.0
	}

	weaveEnd := time.Duration(float64(sim.CurrentTime) + (4.5+2*latencySecs)*float64(time.Second))

	bearweaveNow := rotation.BearweaveType != proto.FeralDruid_Rotation_None && curEnergy <= weaveEnergy && !isClearcast && (!ripRefreshPending || cat.RipDot.ExpiresAt() >= weaveEnd) && !cat.BerserkAura.IsActive()

	if bearweaveNow && rotation.BearweaveType != proto.FeralDruid_Rotation_Lacerate {
		bearweaveNow = !cat.tfExpectedBefore(sim, weaveEnd)
	}

	// If we're maintaining Lacerate, then allow for emergency bearweaves
	// if Lacerate is about to fall off even if the above conditions do not
	// apply.
	lacRemain := core.TernaryDuration(cat.LacerateDot.IsActive(), cat.LacerateDot.RemainingDuration(sim), time.Duration(0))
	emergencyBearweave := rotation.BearweaveType == proto.FeralDruid_Rotation_Lacerate && cat.LacerateDot.IsActive() && (float64(lacRemain) < (2.5+latencySecs)*float64(time.Second)) && (lacRemain < simTimeRemain)

	if bearweaveNow || emergencyBearweave {
		// oom check, if we arent able to shift into bear and back
		// then abandon bearweave
		if cat.CurrentMana() < shiftCost*2.0 {
			bearweaveNow = false
			emergencyBearweave = false
			cat.Metrics.MarkOOM(sim)
		}
	}

	floatingEnergy := 0.0
	previousTime := sim.CurrentTime

	for _, s := range pendingActions {
		delta_t := float64((s.refreshTime - previousTime) / core.EnergyTickDuration)
		if delta_t < s.cost {
			floatingEnergy += s.cost - delta_t
			previousTime = s.refreshTime
		} else {
			previousTime += time.Duration(s.cost * float64(core.EnergyTickDuration))
		}
	}

	excessE := curEnergy - floatingEnergy
	timeToNextAction := time.Duration(0)

	if !cat.CatFormAura.IsActive() {
		// Shift back into Cat Form if (a) our first bear auto procced
		// Clearcasting, or (b) our first bear auto didn't generate enough
		// Rage to Mangle or Maul, or (c) we don't have enough time or
		// Energy leeway to spend an additional GCD in Dire Bear Form.
		shiftNow := (curEnergy+15.0+(10.0*latencySecs) > furorCap) || (ripRefreshPending && (cat.RipDot.RemainingDuration(sim) < (3.0 * time.Second)))
		shiftNext := (curEnergy+30.0+(10.0*latencySecs) > furorCap) || (ripRefreshPending && (cat.RipDot.RemainingDuration(sim) < time.Duration(4500*time.Millisecond)))

		var powerbearNow bool
		if rotation.Powerbear {
			powerbearNow = !shiftNow && curRage < 10
		} else {
			powerbearNow = false
			shiftNow = shiftNow || curRage < 10
		}

		buildLacerate := !cat.LacerateDot.IsActive() || cat.LacerateDot.GetStacks() < 5

		maintainLacerate := !buildLacerate && (lacRemain <= rotation.LacerateTime) && (curRage < 38 || shiftNext) && (lacRemain < simTimeRemain)

		lacerateNow := rotation.BearweaveType == proto.FeralDruid_Rotation_Lacerate && (buildLacerate || maintainLacerate)

		emergencyLacerate := rotation.BearweaveType == proto.FeralDruid_Rotation_Lacerate && cat.LacerateDot.IsActive() && (lacRemain < 3*time.Second+2*cat.latency) && lacRemain < simTimeRemain

		if (rotation.BearweaveType != proto.FeralDruid_Rotation_Lacerate) || !lacerateNow {
			shiftNow = shiftNow || isClearcast
		}

		if emergencyLacerate && cat.CanLacerate(sim) {
			cat.Lacerate.Cast(sim, cat.CurrentTarget)
			return
		} else if shiftNow {

			// If we are resetting our swing timer using Albino Snake or a
			// duplicate weapon swap, then do an additional check here to
			// see whether we can delay the shift until the next bear swing
			// goes out in order to maximize the gains from the reset.
			projectedDelay := cat.AutoAttacks.MainhandSwingAt + 2*cat.latency - sim.CurrentTime
			ripConflict := ripRefreshPending && (cat.RipDot.ExpiresAt() < sim.CurrentTime+projectedDelay+(1500*time.Millisecond))
			nextCatSwing := sim.CurrentTime + cat.latency + time.Duration(float64(cat.AutoAttacks.MainhandSwingSpeed())/float64(2500*time.Millisecond))
			canDelayShift := !ripConflict && cat.Rotation.SnekWeave && (curEnergy+10*projectedDelay.Seconds() <= furorCap) && (cat.AutoAttacks.MainhandSwingAt < nextCatSwing)

			if canDelayShift {
				timeToNextAction = cat.AutoAttacks.MainhandSwingAt - sim.CurrentTime
			} else {
				cat.readyToShift = true
			}
		} else if powerbearNow {
			cat.shiftBearCat(sim, true)
		} else if lacerateNow && cat.CanLacerate(sim) {
			cat.Lacerate.Cast(sim, cat.CurrentTarget)
			return
		} else if cat.CanMangleBear(sim) {
			cat.MangleBear.Cast(sim, cat.CurrentTarget)
			return
		} else if cat.CanLacerate(sim) {
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
		if cat.CanSavageRoar() {
			cat.SavageRoar.Cast(sim, nil)
			return
		} else {
			timeToNextAction = time.Duration((cat.CurrentSavageRoarCost() - curEnergy) * float64(core.EnergyTickDuration))
		}
	} else if ripNow {
		if cat.CanRip() {
			cat.Rip.Cast(sim, cat.CurrentTarget)
			return
		} else {
			timeToNextAction = time.Duration((cat.CurrentRipCost() - curEnergy) * float64(core.EnergyTickDuration))
		}
	} else if biteNow {
		if cat.CanFerociousBite() {
			cat.FerociousBite.Cast(sim, cat.CurrentTarget)
			return
		} else {
			timeToNextAction = time.Duration((cat.CurrentFerociousBiteCost() - curEnergy) * float64(core.EnergyTickDuration))
		}
	} else if rakeNow {
		if cat.CanRake() {
			cat.Rake.Cast(sim, cat.CurrentTarget)
			return
		} else {
			timeToNextAction = time.Duration((cat.CurrentRakeCost() - curEnergy) * float64(core.EnergyTickDuration))
		}
	} else if mangleNow {
		if cat.CanMangleCat() {
			cat.MangleCat.Cast(sim, cat.CurrentTarget)
			return
		} else {
			timeToNextAction = time.Duration((cat.CurrentMangleCatCost() - curEnergy) * float64(core.EnergyTickDuration))
		}
	} else if bearweaveNow {
		cat.readyToShift = true
	} else if (rotation.MangleSpam && !isClearcast) || cat.PseudoStats.InFrontOfTarget {
		if cat.MangleCat != nil && excessE >= cat.CurrentMangleCatCost() {
			cat.MangleCat.Cast(sim, cat.CurrentTarget)
			return
		} else {
			timeToNextAction = time.Duration((cat.CurrentMangleCatCost() - excessE) * float64(core.EnergyTickDuration))
		}
	} else {
		if excessE >= cat.CurrentShredCost() || isClearcast {
			cat.Shred.Cast(sim, cat.CurrentTarget)
			return
		} else {
			timeToNextAction = time.Duration((cat.CurrentShredCost() - excessE) * float64(core.EnergyTickDuration))
		}
	}

	// Model in latency when waiting on Energy for our next action
	nextAction := sim.CurrentTime + timeToNextAction
	if len(pendingActions) > 0 {
		nextAction = core.MinDuration(nextAction, pendingActions[0].refreshTime)
	}

	nextAction += cat.latency

	if nextAction <= sim.CurrentTime {
		cat.DoNothing()
	} else {
		cat.WaitUntil(sim, nextAction)
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
	MaxRoarOffset      time.Duration
	RevitFreq          float64
	LacerateTime       time.Duration
	SnekWeave          bool
}

func (cat *FeralDruid) setupRotation(rotation *proto.FeralDruid_Rotation) {
	cat.Rotation = FeralDruidRotation{
		BearweaveType:      rotation.BearWeaveType,
		MaintainFaerieFire: rotation.MaintainFaerieFire,
		MinCombosForRip:    rotation.MinCombosForRip,
		UseRake:            rotation.UseRake,
		UseBite:            rotation.UseBite,
		BiteTime:           time.Duration(float64(rotation.BiteTime) * float64(time.Second)),
		MinCombosForBite:   rotation.MinCombosForBite,
		MangleSpam:         rotation.MangleSpam,
		BerserkBiteThresh:  float64(rotation.BerserkBiteThresh),
		Powerbear:          rotation.Powerbear,
		MaxRoarOffset:      time.Duration(float64(rotation.MaxRoarOffset) * float64(time.Second)),
		RevitFreq:          15.0 / (8 * float64(rotation.HotUptime)),
		LacerateTime:       10.0 * time.Second,
		SnekWeave:          core.Ternary(rotation.BearWeaveType == proto.FeralDruid_Rotation_None, false, rotation.SnekWeave),
	}
}
