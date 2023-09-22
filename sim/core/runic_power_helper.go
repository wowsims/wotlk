package core

import (
	"fmt"
	"time"
)

// RuneCost's bit layout is: <rrrr_rrrr_dduu_ffbb>. Each part is just a count now (0..3 for runes).
type RuneCost uint16

func NewRuneCost(rp, blood, frost, unholy, death int8) RuneCost {
	return RuneCost(rp)<<8 | RuneCost((death&0b11)<<6|(unholy&0b11)<<4|(frost&0b11)<<2|blood&0b11)
}

func (rc RuneCost) String() string {
	return fmt.Sprintf("RP: %d, Blood: %d, Frost: %d, Unholy: %d, Death: %d", rc.RunicPower(), rc.Blood(), rc.Frost(), rc.Unholy(), rc.Death())
}

// HasRune returns if this cost includes a rune portion.
func (rc RuneCost) HasRune() bool {
	return rc&0b1111_1111 > 0
}

func (rc RuneCost) RunicPower() int8 {
	return int8(rc >> 8)
}

func (rc RuneCost) Blood() int8 {
	return int8(rc & 0b11)
}

func (rc RuneCost) Frost() int8 {
	return int8((rc >> 2) & 0b11)
}

func (rc RuneCost) Unholy() int8 {
	return int8((rc >> 4) & 0b11)
}

func (rc RuneCost) Death() int8 {
	return int8((rc >> 6) & 0b11)
}

func (rp *RunicPowerBar) GainDeathRuneMetrics(sim *Simulation, _ *Spell, currRunes int32, newRunes int32) {
	if !rp.isACopy {
		metrics := rp.deathRuneGainMetrics
		metrics.AddEvent(1, float64(newRunes)-float64(currRunes))

		if sim.Log != nil {
			rp.unit.Log(sim, "Gained 1.000 death rune from %s (%d --> %d).", metrics.ActionID, currRunes, newRunes)
		}
	}
}

func (rp *RunicPowerBar) CancelBloodTap(sim *Simulation) {
	if rp.btslot == -1 {
		return
	}
	rp.ConvertFromDeath(sim, rp.btslot)
	bloodTapAura := rp.unit.GetAura("Blood Tap")
	bloodTapAura.Deactivate(sim)
	rp.btslot = -1
}

func (rp *RunicPowerBar) CorrectBloodTapConversion(sim *Simulation) {
	// 1. converts a blood rune -> death rune
	// 2. then convert one inactive blood or death rune -> active
	slot := int8(-1)
	if rp.runeStates&isDeaths[0] == 0 {
		slot = 0
	} else if rp.runeStates&isDeaths[1] == 0 {
		slot = 1
	}
	if slot > -1 {
		rp.btslot = slot
		rp.ConvertToDeath(sim, slot, sim.CurrentTime+time.Second*20)
	}

	slot = -1
	if rp.runeStates&isSpents[0] == isSpents[0] {
		slot = 0
	} else if rp.runeStates&isSpents[1] == isSpents[1] {
		slot = 1
	}
	if slot > -1 {
		rp.regenRune(sim, sim.CurrentTime, slot)
	}

	// if PA isn't running, make it run 20s from now to disable BT
	rp.launchPA(sim, sim.CurrentTime+20.0*time.Second)
}
