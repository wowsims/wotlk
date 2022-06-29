package feral

import (
	"time"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/items"
	"github.com/wowsims/tbc/sim/core/proto"
	"github.com/wowsims/tbc/sim/druid"
)

func (cat *FeralDruid) OnGCDReady(sim *core.Simulation) {
	cat.doRotation(sim)
}

// Ported from https://github.com/NerdEgghead/TBC_cat_sim

func (cat *FeralDruid) shift(sim *core.Simulation) bool {
	cat.waitingForTick = false

	// If we have just now decided to shift, then we do not execute the
	// shift immediately, but instead trigger an input delay for realism.
	if !cat.readyToShift {
		cat.readyToShift = true
		return false
	}

	cat.readyToShift = false
	return cat.PowerShiftCat(sim)
}

func (cat *FeralDruid) doRotation(sim *core.Simulation) bool {
	// On gcd do nothing
	if !cat.GCD.IsReady(sim) {
		return false
	}

	// If we're out of form always shift back in
	if !cat.InForm(druid.Cat) {
		return cat.CatForm.Cast(sim, nil)
	}

	// If we previously decided to shift, then execute the shift now once
	// the input delay is over.
	if cat.readyToShift {
		return cat.shift(sim)
	}

	rotation := &cat.Rotation

	if cat.Rotation.MaintainFaerieFire && cat.ShouldFaerieFire(sim) {
		return cat.FaerieFire.Cast(sim, cat.CurrentTarget)
	}

	energy := cat.CurrentEnergy()
	cp := cat.ComboPoints()
	ripDebuff := cat.RipDot.IsActive()
	ripEnd := cat.RipDot.ExpiresAt()
	mangleDebuff := cat.MangleAura.IsActive()
	mangleEnd := cat.MangleAura.ExpiresAt()
	rakeDebuff := cat.RakeDot.IsActive()
	nextTick := cat.NextEnergyTickAt()
	shiftCost := cat.CatForm.DefaultCast.Cost
	omenProc := cat.ClearcastingActive()

	// 10/6/21 - Added logic to not cast Rip if we're near the end of the
	// fight.

	ripNow := cp >= rotation.RipCP && !ripDebuff
	ripweaveNow := rotation.UseRipTrick &&
		cp >= rotation.RipTrickCP &&
		!ripDebuff &&
		energy >= RipTrickMin &&
		!cat.PseudoStats.NoCost

	remainingDuration := sim.GetRemainingDuration()
	ripNow = (ripNow || ripweaveNow) && remainingDuration >= RipEndThresh

	// TODO: Can we use fight % completion as an estimate for this instead of exact time calculation?
	biteAtEnd := (cp >= rotation.BiteCP &&
		(remainingDuration < RipEndThresh ||
			(ripDebuff && sim.Duration-ripEnd < RipEndThresh)))

	mangleNow := !ripNow && !mangleDebuff
	mangleCost := cat.Mangle.DefaultCast.Cost

	biteBeforeRip := ripDebuff && rotation.UseBite &&
		ripEnd-sim.CurrentTime >= BiteTime

	biteNow := (biteBeforeRip || rotation.BiteOverRip) &&
		cp >= rotation.BiteCP

	// TODO: Can we use fight % completion as an estimate for this instead of exact time calculation?
	ripNext := (ripNow || (cp >= rotation.RipCP && ripEnd <= nextTick)) &&
		sim.Duration-nextTick >= RipEndThresh

	mangleNext := !ripNext && (mangleNow || mangleEnd <= nextTick)

	// 12/2/21 - Added waitToMangle parameter that tells us whether we
	// should wait for the next Energy tick and cast Mangle, assuming we
	// are less than a tick's worth of Energy from being able to cast it. In
	// a standard Wolfshead rotation, waitForMangle is identical to
	// mangleNext, i.e. we only wait for the tick if Mangle will have
	// fallen off before the next tick. In a no-Wolfshead rotation, however,
	// it is preferable to Mangle rather than Shred as the second special in
	// a standard cycle, provided a bonus like 2pT6 is present to bring the
	// Mangle Energy cost down to 38 or below so that it can be fit in
	// alongside a Shred.
	waitToMangle := mangleNext || (!rotation.Wolfshead && mangleCost <= 38)

	biteBeforeRipNext := biteBeforeRip && ripEnd-nextTick >= BiteTime

	prioBiteOverMangle := rotation.BiteOverRip || !mangleNow

	timeToNextTick := nextTick - sim.CurrentTime
	cat.waitingForTick = true
	markOOM := false

	if cat.CurrentMana() < shiftCost {
		// No-shift rotation
		if ripNow && (energy >= 30 || omenProc) {
			cat.Metrics.MarkOOM(&cat.Unit, time.Second)
			return cat.Rip.Cast(sim, cat.CurrentTarget)
		} else if mangleNow && (energy >= mangleCost || omenProc) {
			cat.Metrics.MarkOOM(&cat.Unit, time.Second)
			return cat.Mangle.Cast(sim, cat.CurrentTarget)
		} else if biteNow && (energy >= 35 || omenProc) {
			cat.Metrics.MarkOOM(&cat.Unit, time.Second)
			return cat.FerociousBite.Cast(sim, cat.CurrentTarget)
		} else if energy >= 42 || omenProc {
			cat.Metrics.MarkOOM(&cat.Unit, time.Second)
			return cat.Shred.Cast(sim, cat.CurrentTarget)
		} else {
			markOOM = true
		}
	} else if energy < 10 {
		cat.shift(sim)
	} else if ripNow {
		if energy >= 30 || omenProc {
			cat.Rip.Cast(sim, cat.CurrentTarget)
			cat.waitingForTick = false
		} else if timeToNextTick > MaxWaitTime {
			cat.shift(sim)
		}
	} else if (biteNow || biteAtEnd) && prioBiteOverMangle {
		// Decision tree for Bite usage is more complicated, so there is
		// some duplicated logic with the main tree.

		// Shred versus Bite decision is the same as vanilla criteria.

		// Bite immediately if we'd have to wait for the following cast.
		cutoffMod := 20.0
		if timeToNextTick <= time.Second {
			cutoffMod = 0.0
		}
		if energy >= 57.0+cutoffMod || (energy >= 15+cutoffMod && omenProc) {
			return cat.Shred.Cast(sim, cat.CurrentTarget)
		}
		if energy >= 35 {
			return cat.FerociousBite.Cast(sim, cat.CurrentTarget)
		}
		// If we are doing the Rip rotation with Bite filler, then there is
		// a case where we would Bite now if we had enough energy, but once
		// we gain enough energy to do so, it's too late to Bite relative to
		// Rip falling off. In this case, we wait for the tick only if we
		// can Shred or Mangle afterward, and otherwise shift and won't Bite
		// at all this cycle. Returning 0.0 is the same thing as waiting for
		// the next tick, so this logic could be written differently if
		// desired to match the rest of the rotation code, where waiting for
		// tick is handled implicitly instead.
		wait := false
		if energy >= 22 && biteBeforeRip && !biteBeforeRipNext {
			wait = true
		} else if energy >= 15 && (!biteBeforeRip || biteBeforeRipNext || biteAtEnd) {
			wait = true
		} else if !ripNext && (energy < 20 || !mangleNext) {
			wait = false
			cat.shift(sim)
		} else {
			wait = true
		}
		if wait && timeToNextTick > MaxWaitTime {
			cat.shift(sim)
		}
	} else if energy >= 35 && energy <= BiteTrickMax &&
		rotation.UseRakeTrick &&
		timeToNextTick > cat.latency &&
		!omenProc &&
		cp >= BiteTrickCP {
		return cat.FerociousBite.Cast(sim, cat.CurrentTarget)
	} else if energy >= 35 && energy < mangleCost &&
		rotation.UseRakeTrick &&
		timeToNextTick > time.Second+cat.latency &&
		!rakeDebuff &&
		!omenProc {
		return cat.Rake.Cast(sim, cat.CurrentTarget)
	} else if mangleNow {
		if energy < mangleCost-20 && !ripNext {
			cat.shift(sim)
		} else if energy >= mangleCost || omenProc {
			return cat.Mangle.Cast(sim, cat.CurrentTarget)
		} else if timeToNextTick > MaxWaitTime {
			cat.shift(sim)
		}
	} else if energy >= 22 {
		if omenProc {
			return cat.Shred.Cast(sim, cat.CurrentTarget)
		}
		// If our energy value is between 50-56 with 2pT6, or 60-61 without,
		// and we are within 1 second of an Energy tick, then Shredding now
		// forces us to shift afterwards, whereas we can instead cast two
		// Mangles instead for higher cpm. This scenario is most relevant
		// when using a no-Wolfshead rotation with 2pT6, and it will
		// occur whenever the initial Shred on a cycle misses.
		if energy >= 2*mangleCost-20 && energy < 22+mangleCost &&
			timeToNextTick <= 1.0*time.Second &&
			rotation.UseMangleTrick &&
			(!rotation.UseRakeTrick || mangleCost == 35) {
			return cat.Mangle.Cast(sim, cat.CurrentTarget)
		}
		if energy >= 42 {
			return cat.Shred.Cast(sim, cat.CurrentTarget)
		}
		if energy >= mangleCost && timeToNextTick > time.Second+cat.latency {
			return cat.Mangle.Cast(sim, cat.CurrentTarget)
		}
		if timeToNextTick > MaxWaitTime {
			cat.shift(sim)
		}
	} else if !ripNext && (energy < mangleCost-20 || !waitToMangle) {
		cat.shift(sim)
	} else if timeToNextTick > MaxWaitTime {
		cat.shift(sim)
	}
	// Model two types of input latency: (1) When waiting for an energy tick
	// to execute the next special ability, the special will in practice be
	// slightly delayed after the tick arrives. (2) When executing a
	// powershift without clipping the GCD, the shift will in practice be
	// slightly delayed after the GCD ends.

	if cat.readyToShift {
		cat.SetGCDTimer(sim, sim.CurrentTime+cat.latency)
	} else if cat.waitingForTick {
		cat.SetGCDTimer(sim, sim.CurrentTime+timeToNextTick+cat.latency)
		if markOOM {
			cat.Metrics.MarkOOM(&cat.Unit, timeToNextTick+cat.latency)
		}
	}

	return false
}

const BiteTrickCP = int32(2)
const BiteTrickMax = 39.0
const BiteTime = time.Second * 0.0
const RipTrickMin = 52.0
const RipEndThresh = time.Second * 10
const MaxWaitTime = time.Second * 1.0

type FeralDruidRotation struct {
	RipCP          int32
	BiteCP         int32
	RipTrickCP     int32
	UseBite        bool
	BiteOverRip    bool
	UseMangleTrick bool
	UseRipTrick    bool
	UseRakeTrick   bool
	Wolfshead      bool

	MaintainFaerieFire bool
}

func (cat *FeralDruid) setupRotation(rotation *proto.FeralDruid_Rotation) {

	UseBite := (rotation.Biteweave && rotation.FinishingMove == proto.FeralDruid_Rotation_Rip) ||
		rotation.FinishingMove == proto.FeralDruid_Rotation_Bite
	RipCP := rotation.RipMinComboPoints

	if rotation.FinishingMove != proto.FeralDruid_Rotation_Rip {
		RipCP = 6
	}

	cat.Rotation = FeralDruidRotation{
		RipCP:          RipCP,
		BiteCP:         rotation.BiteMinComboPoints,
		RipTrickCP:     rotation.RipMinComboPoints,
		UseBite:        UseBite,
		BiteOverRip:    UseBite && rotation.FinishingMove != proto.FeralDruid_Rotation_Rip,
		UseMangleTrick: rotation.MangleTrick,
		UseRipTrick:    rotation.Ripweave,
		UseRakeTrick:   rotation.RakeTrick && !druid.ItemSetThunderheartHarness.CharacterHasSetBonus(&cat.Character, 2),
		Wolfshead:      cat.Equip[items.ItemSlotHead].ID == 8345,

		MaintainFaerieFire: rotation.MaintainFaerieFire,
	}

}
