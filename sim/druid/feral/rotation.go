package feral

import (
	"sort"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
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

	panic("bear auto attacks not implemented")
}

// Ported from https://github.com/NerdEgghead/WOTLK_cat_sim

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
	// For now, consider only the case where Rip will expire after Roar
	ripdotRemaining := cat.RipDot.RemainingDuration(sim)
	if !cat.RipDot.IsActive() || (ripdotRemaining <= cat.SavageRoarAura.RemainingDuration(sim)) || (sim.GetRemainingDuration()-ripdotRemaining < 10*time.Second) {
		return false
	}

	// Calculate how much Energy we expect to accumulate after Roar expires
	// but before Rip expires.
	maxRipDur := time.Duration(float64(cat.maxRipTicks) * float64(cat.RipDot.TickLength))

	ripDur := cat.RipDot.Aura.StartedAt() + maxRipDur - sim.CurrentTime
	roarDur := cat.SavageRoarAura.ExpiresAt() - sim.CurrentTime
	availableTime := ripDur - roarDur
	expectedEnergyGain := 10.0 * float64(availableTime/time.Second)

	if cat.tfExpectedBefore(sim, cat.RipDot.ExpiresAt()) {
		expectedEnergyGain += 60.0
	}
	if cat.Talents.OmenOfClarity {
		expectedEnergyGain += float64(availableTime/cat.AutoAttacks.MainhandSwingSpeed()) / float64(time.Second) * (3.5 / 60. * (1.0 - cat.missChance) * 42.0)
	}

	if cat.ClearcastingAura.IsActive() {
		expectedEnergyGain += 42.0
	}

	expectedEnergyGain += float64(availableTime/time.Second) / cat.Rotation.RevitFreq * 0.15 * 8.0

	// Add current Energy minus cost of Roaring now
	roarCost := core.TernaryFloat64(cat.BerserkAura.IsActive(), 12.5, 25.0)
	availableEnergy := cat.CurrentEnergy() - roarCost + expectedEnergyGain

	// Now calculate the effective Energy cost for building back 5 CPs once
	// Roar expires and casting Rip
	ripCost := core.TernaryFloat64(cat.berserkExpectedAt(sim, cat.RipDot.ExpiresAt()), 15.0, 30.0)
	cpPerBuilder := 1 + cat.GetStat(stats.MeleeCrit)
	costPerBuilder := (42. + 42. + 35.) / 3. * (1 + 0.2*cat.missChance)
	ripRefreshCost := 5./cpPerBuilder*costPerBuilder + ripCost

	// If the cost is less than the expected Energy gain in the available
	// time, then there's no reason to clip Roar.
	if availableEnergy >= ripRefreshCost {
		return false
	}

	// On the other hand, if there is a time conflict, then use the
	// empirical parameter for how much we're willing to clip Roar.
	return roarDur <= cat.Rotation.MaxRoarClip
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

func (cat *FeralDruid) doOnAutoAttack(sim *core.Simulation, spell *core.Spell) {

}

func (cat *FeralDruid) doTigersFury(sim *core.Simulation) {
	// Handle tigers fury
	leewayTime := core.MaxFloat(float64(cat.GCD.TimeToReady(sim)/time.Second), float64(cat.latency/time.Second))
	tfEnergyThresh := 40.0 - 10.0*(leewayTime+core.TernaryFloat64(cat.ClearcastingAura.IsActive(), 1.0, 0))
	tfNow := (cat.CurrentEnergy() < tfEnergyThresh) && cat.TigersFury.IsReady(sim) && !cat.BerserkAura.IsActive()

	if tfNow {
		cat.TigersFury.Cast(sim, nil)
	}
}

func (cat *FeralDruid) doRotation(sim *core.Simulation) {
	// If we previously decided to shift, then execute the shift now once
	// the input delay is over.
	if cat.readyToShift {
		cat.shiftBearCat(sim, false)
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

	endThresh := time.Second * 10

	ripNow := (curCp >= rotation.MinCombosForRip) && !cat.RipDot.IsActive() && (simTimeRemain >= endThresh) && !isClearcast
	biteAtEnd := (curCp >= rotation.MinCombosForBite) && ((simTimeRemain < endThresh) || (cat.RipDot.IsActive() && (simTimeRemain-cat.RipDot.RemainingDuration(sim) < endThresh)))
	mangleNow := !ripNow && !cat.MangleAura.IsActive() && !isClearcast

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
	pendingActions := []pendingAction{}

	if cat.RipDot.IsActive() && (cat.RipDot.RemainingDuration(sim) < simTimeRemain-endThresh) {
		ripCost := core.TernaryFloat64(cat.berserkExpectedAt(sim, cat.RipDot.ExpiresAt()), 15.0, 30.0)
		pendingActions = append(pendingActions, pendingAction{cat.RipDot.ExpiresAt(), ripCost})
		ripRefreshPending = true
	}
	if cat.RakeDot.IsActive() && (cat.RakeDot.RemainingDuration(sim) < simTimeRemain-(9*time.Second)) {
		rakeCost := core.TernaryFloat64(cat.berserkExpectedAt(sim, cat.RakeDot.ExpiresAt()), 17.5, 35.0)
		pendingActions = append(pendingActions, pendingAction{cat.RakeDot.ExpiresAt(), rakeCost})
	}
	if cat.MangleAura.IsActive() && (cat.MangleAura.RemainingDuration(sim) < simTimeRemain-time.Second) {
		mangleCost := cat.MangleCat.BaseCost
		if cat.berserkExpectedAt(sim, cat.MangleAura.ExpiresAt()) {
			mangleCost *= 0.5
		}
		pendingActions = append(pendingActions, pendingAction{cat.MangleAura.ExpiresAt(), mangleCost})
	}
	if cat.SavageRoarAura.IsActive() {
		roarCost := core.TernaryFloat64(cat.berserkExpectedAt(sim, cat.SavageRoarAura.ExpiresAt()), 12.5, 25)
		pendingActions = append(pendingActions, pendingAction{cat.SavageRoarAura.ExpiresAt(), roarCost})
	}

	sort.SliceStable(pendingActions, func(i, j int) bool {
		return pendingActions[i].refreshTime < pendingActions[j].refreshTime
	})

	if !cat.SavageRoarAura.IsActive() && curEnergy >= 25 {
		endThresh = time.Second * 10
	}

	latencySecs := cat.latency / time.Second
	// Allow for bearweaving if the next pending action is >= 4.5s away
	furorCap := core.MinFloat(20.0*float64(cat.Talents.Furor), 85)
	weaveEnergy := furorCap - 30 - 20*float64(latencySecs)

	if cat.Talents.Furor > 3 {
		weaveEnergy -= 15.0
	}

	weaveEnd := time.Duration(float64(sim.CurrentTime) + (float64(4.5)+2)*float64(latencySecs))

	bearweaveNow := rotation.BearweaveType != proto.FeralDruid_Rotation_None && curEnergy <= weaveEnergy && !isClearcast && (!ripRefreshPending || cat.RipDot.ExpiresAt() >= weaveEnd) && !cat.BerserkAura.IsActive()

	if bearweaveNow && rotation.BearweaveType != proto.FeralDruid_Rotation_Lacerate {
		bearweaveNow = !cat.tfExpectedBefore(sim, weaveEnd)
	}

	// If we're maintaining Lacerate, then allow for emergency bearweaves
	// if Lacerate is about to fall off even if the above conditions do not
	// apply.
	lacRemain := core.TernaryDuration(cat.LacerateDot.IsActive(), cat.LacerateDot.RemainingDuration(sim), time.Duration(0))
	emergencyBearweave := rotation.BearweaveType == proto.FeralDruid_Rotation_Lacerate && cat.LacerateDot.IsActive() && (float64(lacRemain) < (float64(2.5)*float64(time.Second) + float64(latencySecs))) && (lacRemain < simTimeRemain)

	floatingEnergy := 0.0
	previousTime := sim.CurrentTime

	for _, s := range pendingActions {
		delta_t := s.refreshTime - previousTime
		if float64(delta_t/time.Second) < s.cost/10.0 {
			floatingEnergy += s.cost - 10.0*float64(delta_t/time.Second)
			previousTime = s.refreshTime
		} else {
			previousTime += time.Duration((s.cost / 10.0) * float64(time.Second))
		}
	}

	excessE := curEnergy - floatingEnergy
	timeToNextAction := time.Duration(0)

	if !cat.CatFormAura.IsActive() {
		// Shift back into Cat Form if (a) our first bear auto procced
		// Clearcasting, or (b) our first bear auto didn't generate enough
		// Rage to Mangle or Maul, or (c) we don't have enough time or
		// Energy leeway to spend an additional GCD in Dire Bear Form.
		shiftNow := (curEnergy+15.0+(10.0*float64(latencySecs)) > furorCap) || (ripRefreshPending && (cat.RipDot.ExpiresAt() < sim.CurrentTime+(3.0*time.Second)))
		shiftNext := (curEnergy+30.0+(10.0*float64(latencySecs)) > furorCap) || (ripRefreshPending && (cat.RipDot.ExpiresAt() < sim.CurrentTime+time.Duration(float64(4.5)*float64(time.Second))))

		powerbearNow := true
		if rotation.BearweaveType == proto.FeralDruid_Rotation_None {
			powerbearNow = false
			shiftNow = shiftNow || curRage < 10
		}

		buildLacerate := !cat.LacerateDot.IsActive() || cat.LacerateDot.GetStacks() < 5

		maintainLacerate := !buildLacerate && (lacRemain <= rotation.LacerateTime) && (curRage < 38 || shiftNext) && (lacRemain < simTimeRemain)

		lacerateNow := rotation.BearweaveType == proto.FeralDruid_Rotation_Lacerate && (buildLacerate || maintainLacerate)

		emergencyLacerate := rotation.BearweaveType == proto.FeralDruid_Rotation_Lacerate && cat.LacerateDot.IsActive() && (float64(lacRemain) < (3.0+2)*float64(latencySecs)) && lacRemain < simTimeRemain

		if (rotation.BearweaveType != proto.FeralDruid_Rotation_Lacerate) || !lacerateNow {
			shiftNow = shiftNow || isClearcast
		}

		if emergencyLacerate && cat.CanLacerate(sim) {
			cat.Lacerate.Cast(sim, cat.CurrentTarget)
			return
		} else if shiftNow {
			cat.readyToShift = true
		} else if powerbearNow {
			cat.shiftBearCat(sim, true)
		} else if lacerateNow && cat.CanLacerate(sim) {
			cat.Lacerate.Cast(sim, cat.CurrentTarget)
			return
		} else if cat.CanMangleBear(sim) && !isClearcast {
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
		return
	} else if roarNow {
		if cat.CanSavageRoar() {
			cat.SavageRoar.Cast(sim, nil)
			return
		} else {
			timeToNextAction = time.Duration(((cat.CurrentSavageRoarCost() - curEnergy) / 10.0) * float64(time.Second))
		}
	} else if ripNow {
		if cat.CanRip() {
			cat.Rip.Cast(sim, cat.CurrentTarget)
			return
		} else {
			timeToNextAction = time.Duration(((cat.CurrentRipCost() - curEnergy) / 10.0) * float64(time.Second))
		}
	} else if biteNow {
		if cat.CanFerociousBite() {
			cat.FerociousBite.Cast(sim, cat.CurrentTarget)
			return
		} else {
			timeToNextAction = time.Duration(((cat.CurrentFerociousBiteCost() - curEnergy) / 10.0) * float64(time.Second))
		}
	} else if rakeNow {
		if cat.CanRake() {
			cat.Rake.Cast(sim, cat.CurrentTarget)
			return
		} else {
			timeToNextAction = time.Duration(((cat.CurrentRakeCost() - curEnergy) / 10.0) * float64(time.Second))
		}
	} else if mangleNow {
		if cat.CanMangleCat() {
			cat.MangleCat.Cast(sim, cat.CurrentTarget)
			return
		} else {
			timeToNextAction = time.Duration(((cat.CurrentMangleCatCost() - curEnergy) / 10.0) * float64(time.Second))
		}
	} else if bearweaveNow {
		cat.readyToShift = true
	} else if rotation.MangleSpam && !isClearcast {
		if excessE >= cat.CurrentMangleCatCost() {
			cat.MangleCat.Cast(sim, cat.CurrentTarget)
			return
		} else {
			timeToNextAction = time.Duration(((cat.CurrentMangleCatCost() - excessE) / 10.0) * float64(time.Second))
		}
	} else {
		if excessE >= cat.CurrentShredCost() || isClearcast {
			cat.Shred.Cast(sim, cat.CurrentTarget)
			return
		} else {
			timeToNextAction = time.Duration(((cat.CurrentShredCost() - excessE) / 10.0) * float64(time.Second))
		}
	}

	// Model in latency when waiting on Energy for our next action
	nextAction := sim.CurrentTime + timeToNextAction
	if len(pendingActions) > 0 {
		nextAction = core.MinDuration(nextAction, pendingActions[0].refreshTime)
	}

	nextAction += latencySecs

	// TODO: This probably shouldnt happen
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
	MaxRoarClip        time.Duration
	RevitFreq          float64
	LacerateTime       time.Duration
}

func (cat *FeralDruid) setupRotation(rotation *proto.FeralDruid_Rotation) {
	hotUptime := 0.75
	cat.Rotation = FeralDruidRotation{
		BearweaveType:      rotation.BearWeaveType,
		MaintainFaerieFire: rotation.MaintainFaerieFire,
		MinCombosForRip:    rotation.MinCombosForRip,
		UseRake:            rotation.UseRake,
		UseBite:            false,
		BiteTime:           time.Duration(float64(rotation.BiteTime) * float64(time.Second)),
		MinCombosForBite:   rotation.MinCombosForBite,
		MangleSpam:         rotation.MangleSpam,
		BerserkBiteThresh:  float64(rotation.BerserkBiteThresh),
		Powerbear:          rotation.Powerbear,
		MaxRoarClip:        time.Duration(float64(rotation.MaxRoarClip) * float64(time.Second)),
		RevitFreq:          15.0 / (8 * hotUptime),
		LacerateTime:       10.0 * time.Second,
	}

}
