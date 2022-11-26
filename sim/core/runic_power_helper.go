package core

import (
	"fmt"
	"time"
)

type RuneCost uint16

func NewRuneCost(rp, blood, frost, unholy, death uint8) RuneCost {
	value := int16(0)
	if blood == 1 {
		value = 1
	} else if blood == 2 {
		value = 3
	}

	if frost == 1 {
		value += 1 << 2
	} else if frost == 2 {
		value += 3 << 2
	}

	if unholy == 1 {
		value += 1 << 4
	} else if unholy == 2 {
		value += 3 << 4
	}

	if death == 1 {
		value += 1 << 6
	} else if death == 2 {
		value += 3 << 6
	} else if death > 2 {
		value += 3 << 6 // we cant represent more than 2 death runes
	}

	value += int16(rp) << 8

	return RuneCost(value)
}

func (rc RuneCost) String() string {
	return fmt.Sprintf("RP: %d, Blood: %d, Frost: %d, Unholy: %d, Death: %d", rc.RunicPower(), rc.Blood(), rc.Frost(), rc.Unholy(), rc.Death())
}

// HasRune returns if this cost includes a rune portion.
//
//	If any bit is set in the rune bits it means that there is a rune cost.
func (rc RuneCost) HasRune() bool {
	const runebits = int16(0b11111111)
	return runebits&int16(rc) > 0
}

func (rc RuneCost) RunicPower() uint8 {
	const rpbits = uint16(0b1111111100000000)
	return uint8((uint16(rc) & rpbits) >> 8)
}

func (rc RuneCost) Blood() uint8 {
	runes := uint16(rc) & 0b11
	switch runes {
	case 0b00:
		return 0
	case 0b01:
		return 1
	case 0b11:
		return 2
	}
	return 0
}

func (rc RuneCost) Frost() uint8 {
	runes := uint16(rc) & 0b1100
	switch runes {
	case 0:
		return 0
	case 0b0100:
		return 1
	case 0b1100:
		return 2
	}
	return 0
}

func (rc RuneCost) Unholy() uint8 {
	runes := uint16(rc) & 0b110000
	switch runes {
	case 0:
		return 0
	case 0b010000:
		return 1
	case 0b110000:
		return 2
	}
	return 0
}

func (rc RuneCost) Death() uint8 {
	runes := uint16(rc) & 0b11000000
	switch runes {
	case 0:
		return 0
	case 0b01000000:
		return 1
	case 0b11000000:
		return 2
	}
	return 0
}

func (rp *RunicPowerBar) GainDeathRuneMetrics(sim *Simulation, spell *Spell, currRunes int32, newRunes int32) {
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

func (rp *RunicPowerBar) CorrectBloodTapConversion(sim *Simulation, bloodGainMetrics *ResourceMetrics, deathGainMetrics *ResourceMetrics, spell *Spell) {
	// 1. converts a blood rune -> death rune
	// 3. then convert one inactive blood or death rune -> active

	slot := int8(-1)
	if rp.runeStates&isDeaths[0] == 0 {
		slot = 0
	} else if rp.runeStates&isDeaths[1] == 0 {
		slot = 1
	}
	if slot > -1 {
		rp.ConvertToDeath(sim, slot, false, sim.CurrentTime+time.Second*20)
		rp.btslot = slot
	}

	slot = -1
	if rp.runeStates&isSpents[0] == isSpents[0] {
		slot = 0
	} else if rp.runeStates&isSpents[1] == isSpents[1] {
		slot = 1
	}
	if slot > -1 {
		rp.runeStates = ^isSpents[slot] & rp.runeStates // unset spent flag for this rune.
		rp.runeMeta[slot].regenAt = NeverExpires        // no regen timer if there was one
		rp.GainRuneMetrics(sim, rp.deathRuneGainMetrics, 1)
		rp.onDeathRuneGain(sim)
	}

	// if PA isn't running, make it run 20s from now to disable BT
	rp.launchPA(sim, sim.CurrentTime+20.0*time.Second)
}
