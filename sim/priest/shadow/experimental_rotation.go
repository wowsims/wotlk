package shadow

import (
	//"fmt"
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

const (
	MindBlastIdx int = iota
	ShadowWordDeathIdx
	DevouringPlagueIdx
	VampiricTouchIdx
	ShadowWordPainIdx
	MindFlay1Idx
	MindFlay2Idx
	MindFlay3Idx
	SpellLen
)

func (spriest *ShadowPriest) experimentalRotation(sim *core.Simulation) {
	spell := spriest.chooseSpellExperimental(sim)

	if spell == spriest.MindFlay[3] && spriest.InnerFocus != nil && spriest.InnerFocus.IsReady(sim) && spriest.ShadowWeavingAura.GetStacks() == 5 {
		spriest.InnerFocus.Cast(sim, nil)
	}

	if !spell.IsReady(sim) {
		spriest.WaitUntil(sim, spell.ReadyAt())
	} else if success := spell.Cast(sim, spriest.CurrentTarget); !success {
		spriest.WaitForMana(sim, spell.CurCast.Cost)
	}
}

func (spriest *ShadowPriest) chooseSpellExperimental(sim *core.Simulation) *core.Spell {
	if !spriest.DevouringPlagueDot.IsActive() {
		return spriest.DevouringPlague
	}

	gcd := core.MaxDuration(core.GCDMin, spriest.ApplyCastSpeed(core.GCDDefault))
	vtCastTime := gcd
	if spriest.VampiricTouchDot != nil && (!spriest.VampiricTouchDot.IsActive() || sim.CurrentTime+vtCastTime >= spriest.VampiricTouchDot.ExpiresAt()) {
		return spriest.VampiricTouch
	}

	if !spriest.ShadowWordPainDot.IsActive() {
		if spriest.CanRolloverSWP {
			// At the beginning of the iteration, its better to wait for 5 stacks of weaving
			// before taking the first snapshot.
			stacks := spriest.ShadowWeavingAura.GetStacks()
			if spriest.ShadowWeavingAura == nil || stacks == 5 {
				return spriest.ShadowWordPain
			} else if stacks == 2 && spriest.MindBlast.IsReady(sim) {
				return spriest.MindBlast
			} else if stacks >= 2 {
				return spriest.MindFlay[5-stacks]
			} else {
				if spriest.options.UseMindBlast {
					return spriest.MindBlast
				} else if spriest.options.UseShadowWordDeath {
					return spriest.ShadowWordDeath
				} else {
					return spriest.MindFlay[2]
				}
			}
		} else {
			return spriest.ShadowWordPain
		}
	}

	//if spriest.CanRolloverSWP {
	//	// If we can rollover SWP and we can get a sufficiently better snapshot, take it.
	//	snapshotDmg := spriest.ShadowWordPain.ExpectedDamageFromCurrentSnapshot(sim, spriest.CurrentTarget)
	//	newDmg := spriest.ShadowWordPain.ExpectedDamage(sim, spriest.CurrentTarget)
	//	if newDmg > snapshotDmg + 200 {
	//		if sim.Log != nil {
	//			spriest.Log(sim, "Better SWP snapshot, old: %0.02f, new: %0.02f", snapshotDmg, newDmg)
	//		}
	//		return spriest.ShadowWordPain
	//	}
	//}

	// Time spent casting each spell.
	castTime := [SpellLen]time.Duration{
		gcd, // MB
		gcd, // SWD
		gcd, // DP
		gcd, // VT
		gcd, // SWP
		spriest.ApplyCastSpeed(spriest.MindFlay[1].DefaultCast.EffectiveTime()) + spriest.AverageMindFlayLatencyDelay(1, gcd),
		spriest.ApplyCastSpeed(spriest.MindFlay[2].DefaultCast.EffectiveTime()) + spriest.AverageMindFlayLatencyDelay(2, gcd),
		spriest.ApplyCastSpeed(spriest.MindFlay[3].DefaultCast.EffectiveTime()) + spriest.AverageMindFlayLatencyDelay(3, gcd),
	}
	// When the GCD would become ready again, for each cast.
	castCompleteAt := [SpellLen]time.Duration{
		core.MaxDuration(sim.CurrentTime, spriest.MindBlast.ReadyAt()) + gcd,
		core.MaxDuration(sim.CurrentTime, spriest.ShadowWordDeath.ReadyAt()) + gcd,
		sim.CurrentTime + gcd, // DP
		sim.CurrentTime + gcd, // VT
		sim.CurrentTime + gcd, // SWP
		sim.CurrentTime + spriest.ApplyCastSpeed(spriest.MindFlay[1].DefaultCast.EffectiveTime()) + spriest.AverageMindFlayLatencyDelay(1, gcd),
		sim.CurrentTime + spriest.ApplyCastSpeed(spriest.MindFlay[2].DefaultCast.EffectiveTime()) + spriest.AverageMindFlayLatencyDelay(2, gcd),
		sim.CurrentTime + spriest.ApplyCastSpeed(spriest.MindFlay[3].DefaultCast.EffectiveTime()) + spriest.AverageMindFlayLatencyDelay(3, gcd),
	}

	// Time period over which a spell deals its full damage.
	cadence := [SpellLen]time.Duration{
		spriest.MindBlast.CD.Duration,
		spriest.ShadowWordDeath.CD.Duration,
		spriest.DevouringPlagueDot.TickPeriod() * time.Duration(spriest.DevouringPlagueDot.NumberOfTicks),
		spriest.VampiricTouchDot.TickPeriod() * time.Duration(spriest.VampiricTouchDot.NumberOfTicks),
		spriest.ShadowWordPainDot.TickPeriod() * time.Duration(spriest.ShadowWordPainDot.NumberOfTicks),
		castCompleteAt[MindFlay1Idx] - sim.CurrentTime,
		castCompleteAt[MindFlay2Idx] - sim.CurrentTime,
		castCompleteAt[MindFlay3Idx] - sim.CurrentTime,
	}

	// Time after which casting the corresponding spell would be a delay of
	// its cadence.
	spellDelayStart := [SpellLen]time.Duration{
		core.MaxDuration(sim.CurrentTime, spriest.MindBlast.ReadyAt()),
		core.MaxDuration(sim.CurrentTime, spriest.ShadowWordDeath.ReadyAt()),
		core.MaxDuration(sim.CurrentTime, spriest.DevouringPlagueDot.ExpiresAt()),
		core.MaxDuration(sim.CurrentTime, spriest.VampiricTouchDot.ExpiresAt()-vtCastTime),
		core.NeverExpires, // SWP rolls over, so it never gets delayed.
		sim.CurrentTime,   // MF1
		sim.CurrentTime,   // MF2
		sim.CurrentTime,   // MF3
	}

	dpTickDamage := spriest.DevouringPlague.ExpectedDamage(sim, spriest.CurrentTarget)
	dpInitDamage := dpTickDamage * spriest.DpInitMultiplier

	// Total damage done by a cast of each spell.
	spellDamage := [SpellLen]float64{
		0,
		0,
		dpInitDamage + 8*dpTickDamage,
		spriest.VampiricTouch.ExpectedDamage(sim, spriest.CurrentTarget) * float64(spriest.VampiricTouchDot.NumberOfTicks),
		spriest.ShadowWordPain.ExpectedDamage(sim, spriest.CurrentTarget) * float64(spriest.ShadowWordPainDot.NumberOfTicks),
		0,
		0,
		0,
	}
	if spriest.options.UseMindBlast {
		spellDamage[MindBlastIdx] = spriest.MindBlast.ExpectedDamage(sim, spriest.CurrentTarget)
		// Account for Glyph of Shadow and imp spirit tap
		spellDamage[MindBlastIdx] += 100.0 * 30
	}
	if spriest.options.UseShadowWordDeath {
		spellDamage[ShadowWordDeathIdx] = spriest.ShadowWordDeath.ExpectedDamage(sim, spriest.CurrentTarget)
		spellDamage[ShadowWordDeathIdx] += 100.0 * 0
	}
	if spriest.Talents.MindFlay {
		spellDamage[MindFlay1Idx] = spriest.MindFlay[1].ExpectedDamage(sim, spriest.CurrentTarget)
		spellDamage[MindFlay2Idx] = spellDamage[MindFlay1Idx] * 2
		spellDamage[MindFlay3Idx] = spellDamage[MindFlay1Idx] * 3
	}

	spellDPS := [SpellLen]float64{
		spellDamage[MindBlastIdx] / cadence[MindBlastIdx].Seconds(),
		spellDamage[ShadowWordDeathIdx] / cadence[ShadowWordDeathIdx].Seconds(),
		spellDamage[DevouringPlagueIdx] / cadence[DevouringPlagueIdx].Seconds(),
		spellDamage[VampiricTouchIdx] / cadence[VampiricTouchIdx].Seconds(),
		spellDamage[ShadowWordPainIdx] / cadence[ShadowWordPainIdx].Seconds(),
		spellDamage[MindFlay1Idx] / cadence[MindFlay1Idx].Seconds(),
		spellDamage[MindFlay2Idx] / cadence[MindFlay2Idx].Seconds(),
		spellDamage[MindFlay3Idx] / cadence[MindFlay3Idx].Seconds(),
	}

	// These are the only efficient options worth considering.
	mbCastTime := castTime[MindBlastIdx]
	mf2CastTime := castTime[MindFlay2Idx]
	mf3CastTime := castTime[MindFlay3Idx]

	//// For each spell option, find the combination of following spells which
	//// causes the LEAST amount of delay on DP and VT, and record that delay.
	//dpCastAt := spellDelayStart[DevouringPlagueIdx]
	//vtCastAt := spellDelayStart[VampiricTouchIdx]
	//scanUntil := core.MaxDuration(dpCastAt, vtCastAt) + gcd*3

	//type PathOption struct {
	//	Spell   *core.Spell
	//	Damage  float64
	//	DoneAt  time.Duration
	//	CanMBAt time.Duration // When Mind Blast will be ready on this path.

	//	PrevOption  *PathOption
	//	DPRefreshed bool
	//	VTRefreshed bool
	//}

	//pathOptions := []PathOption{
	//	PathOption{
	//		Spell:  spriest.MindBlast,
	//		Damage: spellDamage[MindBlastIdx],
	//		DoneAt: castCompleteAt[MindBlastIdx],
	//		CanMBAt: castCompleteAt[MindBlastIdx] + spriest.MindBlast.CD.Duration,
	//	},
	//	PathOption{
	//		Spell:  spriest.MindFlay[2],
	//		Damage: spellDamage[MindFlay2Idx],
	//		DoneAt: castCompleteAt[MindFlay2Idx],
	//		CanMBAt: spellDelayStart[MindBlastIdx],
	//	},
	//	PathOption{
	//		Spell:  spriest.MindFlay[3],
	//		Damage: spellDamage[MindFlay3Idx],
	//		DoneAt: castCompleteAt[MindFlay3Idx],
	//		CanMBAt: spellDelayStart[MindBlastIdx],
	//	},
	//}
	//cur := &pathOptions[0]
	//if pathOptions[1].DoneAt < pathOptions[0].DoneAt {
	//	cur = &pathOptions[1]
	//}

	////var dotDelays [SpellLen]time.Duration
	////var maxDotDelay time.Duration
	//for cur.DoneAt <= scanUntil {
	//	// Process the curTime, by adding new nodes for each cast option.
	//	if !cur.DPRefreshed && cur.DoneAt >= dpCastAt {
	//		pathOptions = append(pathOptions, PathOption{
	//			Damage:  cur.Damage - spellDPS[DevouringPlagueIdx] * (cur.DoneAt - dpCastAt).Seconds(),
	//			DoneAt: cur.DoneAt + gcd,
	//			CanMBAt: cur.CanMBAt,
	//			PrevOption:  cur,
	//			DPRefreshed: true,
	//			VTRefreshed: cur.VTRefreshed,
	//		})
	//	} else if !cur.VTRefreshed && cur.DoneAt >= vtCastAt {
	//		pathOptions = append(pathOptions, PathOption{
	//			Damage:  cur.Damage - spellDPS[VampiricTouchIdx] * (cur.DoneAt - vtCastAt).Seconds(),
	//			DoneAt: cur.DoneAt + gcd,
	//			CanMBAt: cur.CanMBAt,
	//			PrevOption:  cur,
	//			DPRefreshed: cur.DPRefreshed,
	//			VTRefreshed: true,
	//		})
	//	} else {
	//		pathOptions = append(pathOptions, PathOption{
	//			Damage:  cur.Damage + spellDamage[MindFlay3Idx],
	//			DoneAt:  cur.DoneAt + mf3CastTime,
	//			CanMBAt: cur.CanMBAt,
	//			PrevOption:  cur,
	//			DPRefreshed: cur.DPRefreshed,
	//			VTRefreshed: cur.VTRefreshed,
	//		})
	//		pathOptions = append(pathOptions, PathOption{
	//			Damage:  cur.Damage + spellDamage[MindFlay2Idx],
	//			DoneAt:  cur.DoneAt + mf2CastTime,
	//			CanMBAt: cur.CanMBAt,
	//			PrevOption:  cur,
	//			DPRefreshed: cur.DPRefreshed,
	//			VTRefreshed: cur.VTRefreshed,
	//		})
	//		if cur.DoneAt >= cur.CanMBAt {
	//			pathOptions = append(pathOptions, PathOption{
	//				Damage:  cur.Damage + spellDamage[MindBlastIdx],
	//				DoneAt: cur.DoneAt + mbCastTime,
	//				CanMBAt: cur.DoneAt + mbCastTime + spriest.MindBlast.CD.Duration,
	//				PrevOption:  cur,
	//				DPRefreshed: cur.DPRefreshed,
	//				VTRefreshed: cur.VTRefreshed,
	//			})
	//		}
	//	}

	//	// Find the next cur, which is the smallest option greater than curTime.
	//	var bestOption *PathOption
	//	for i := 0; i < len(pathOptions); i++ {
	//		if pathOptions[i].DoneAt > cur.DoneAt {
	//			if bestOption == nil || pathOptions[i].DoneAt < bestOption.DoneAt || (pathOptions[i].DoneAt == bestOption.DoneAt && pathOptions[i].Damage > bestOption.Damage) {
	//				bestOption = &pathOptions[i]
	//			}
	//		}
	//	}
	//	cur = bestOption
	//}

	//var bestOption *PathOption
	//for i, _ := range pathOptions {
	//	option := &pathOptions[i]
	//	if option.DPRefreshed && option.VTRefreshed {
	//		if bestOption == nil || option.Damage > bestOption.Damage {
	//			bestOption = option
	//		}
	//	}
	//}
	//if bestOption == nil {
	//	panic(fmt.Sprintf("No best option, %d, %s\n", len(pathOptions), scanUntil - sim.CurrentTime))
	//}

	//for bestOption.PrevOption != nil {
	//	bestOption = bestOption.PrevOption
	//}
	//return bestOption.Spell

	// For each spell option, find the combination of following spells which
	// causes the LEAST amount of delay on DP and VT, and record that delay.
	dot1DelayAt := spellDelayStart[DevouringPlagueIdx]
	dot2DelayAt := spellDelayStart[VampiricTouchIdx]
	if dot1DelayAt > dot2DelayAt {
		// Swap so that dot1DelayAt always comes before dot2DelayAt.
		dot1DelayAt = spellDelayStart[VampiricTouchIdx]
		dot2DelayAt = spellDelayStart[DevouringPlagueIdx]
	}

	type PathOption struct {
		DoneAt  time.Duration
		CanMBAt time.Duration // When Mind Blast will be ready on this path.
	}
	startingPathOptions := []PathOption{
		PathOption{
			DoneAt:  castCompleteAt[MindBlastIdx],
			CanMBAt: castCompleteAt[MindBlastIdx] + spriest.MindBlast.CD.Duration,
		},
		PathOption{
			DoneAt:  castCompleteAt[MindFlay2Idx],
			CanMBAt: spellDelayStart[MindBlastIdx],
		},
		PathOption{
			DoneAt:  castCompleteAt[MindFlay3Idx],
			CanMBAt: spellDelayStart[MindBlastIdx],
		},
	}

	var dotDelays [SpellLen]time.Duration
	var maxDotDelay time.Duration
	for startIdx, startingPathOption := range startingPathOptions {
		cur := startingPathOption
		pathOptions := []PathOption{
			cur,
		}

		for cur.DoneAt <= dot1DelayAt {
			// Process the curTime, by adding new nodes for each cast option.
			pathOptions = append(pathOptions, PathOption{
				DoneAt:  cur.DoneAt + mf2CastTime,
				CanMBAt: cur.CanMBAt,
			})
			pathOptions = append(pathOptions, PathOption{
				DoneAt:  cur.DoneAt + mf3CastTime,
				CanMBAt: cur.CanMBAt,
			})
			if cur.DoneAt >= cur.CanMBAt {
				pathOptions = append(pathOptions, PathOption{
					DoneAt:  cur.DoneAt + mbCastTime,
					CanMBAt: cur.DoneAt + mbCastTime + spriest.MindBlast.CD.Duration,
				})
			}

			// Find the next cur, which is the smallest option greater than curTime.
			var bestOption PathOption
			for i := 0; i < len(pathOptions); i++ {
				if pathOptions[i].DoneAt > cur.DoneAt && (bestOption.DoneAt == 0 || pathOptions[i].DoneAt < bestOption.DoneAt) {
					bestOption = pathOptions[i]
				}
			}
			cur = bestOption
		}

		delay := cur.DoneAt - dot1DelayAt
		maxDotDelay = core.MaxDuration(maxDotDelay, delay)
		if startIdx == 0 {
			dotDelays[MindBlastIdx] = delay
		} else if startIdx == 1 {
			dotDelays[MindFlay2Idx] = delay
		} else {
			dotDelays[MindFlay3Idx] = delay
		}
	}

	// Resulting net damage for each spell, with the opportunity cost of casting
	// other spells subtracted.
	netSpellDamage := [SpellLen]float64{
		spellDamage[0],
		spellDamage[1],
		spellDamage[2],
		spellDamage[3],
		spellDamage[4],
		spellDamage[5],
		spellDamage[6],
		spellDamage[7],
	}

	// Subtract opportunity cost of not choosing each spell.
	for chosenSpellIdx, _ := range spellDamage {
		if spellDamage[chosenSpellIdx] == 0 {
			continue
		}

		for otherSpellIdx, _ := range spellDamage {
			if chosenSpellIdx == otherSpellIdx { //|| otherSpellIdx == MindFlay1Idx || otherSpellIdx == MindFlay3Idx {
				continue
			}

			var delay time.Duration
			if otherSpellIdx == DevouringPlagueIdx || otherSpellIdx == VampiricTouchIdx {
				if chosenSpellIdx == MindBlastIdx || chosenSpellIdx == MindFlay2Idx || chosenSpellIdx == MindFlay3Idx {
					delay = dotDelays[chosenSpellIdx] * 12 / 6
				} else {
					delay = maxDotDelay * 12 / 6
				}
			} else {
				delay = core.MaxDuration(0, castCompleteAt[chosenSpellIdx]-spellDelayStart[otherSpellIdx])
			}
			opportunityCostDmg := spellDPS[otherSpellIdx] * delay.Seconds()
			//opportunityCostDmg := spellDPCT[otherSpellIdx] * delay.Seconds()
			//if opportunityCostDmg < -9999 {
			//	panic(fmt.Sprintf("Opp cost: %0.01f, delay, %s, other dmg: %0.01f, other dpct: %0.01f\n", opportunityCostDmg, delay, spellDamage[otherSpellIdx], spellDPCT[otherSpellIdx]))
			//}
			netSpellDamage[chosenSpellIdx] -= opportunityCostDmg
		}
	}

	bestIdx := 0
	bestDamage := -999999.9
	for i, dmg := range netSpellDamage {
		if dmg >= bestDamage && dmg != 0 {
			if i == DevouringPlagueIdx || i == VampiricTouchIdx || i == ShadowWordPainIdx {
				continue
			}
			bestIdx = i
			bestDamage = dmg
		}
	}

	spells := []*core.Spell{
		spriest.MindBlast,
		spriest.ShadowWordDeath,
		spriest.DevouringPlague,
		spriest.VampiricTouch,
		spriest.ShadowWordPain,
		spriest.MindFlay[1],
		spriest.MindFlay[2],
		spriest.MindFlay[3],
	}
	//if !spells[bestIdx].IsReady(sim) {
	//	fmt.Printf("MB dmg: %0.01f, net: %0.01f\n", spellDamage[MindBlastIdx], netSpellDamage[MindBlastIdx])
	//	fmt.Printf("MF2 dmg: %0.01f, net: %0.01f\n", spellDamage[MindFlay2Idx], netSpellDamage[MindFlay2Idx])
	//}
	return spells[bestIdx]
}
