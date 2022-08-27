package shadow

import (
	"math"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

var mbIdb = 0
var dpIdx = 1
var vtIdx = 2
var swdIdx = 3

func (spriest *ShadowPriest) OnGCDReady(sim *core.Simulation) {
	spriest.tryUseGCD(sim)
}

func (spriest *ShadowPriest) OnManaTick(sim *core.Simulation) {
	if spriest.FinishedWaitingForManaAndGCDReady(sim) {
		spriest.tryUseGCD(sim)
	}
}
func (spriest *ShadowPriest) tryUseGCD(sim *core.Simulation) {
	// TODO: probably do something different instead of making it global?
	const (
		dudidx int = iota
		mfidx
		// TFmod
	)

	// some global variables used througout the code
	var currentWait time.Duration
	var mbDamage float64
	var dpDamage float64
	var vtDamage float64
	var swdDamage float64
	var mfDamage float64
	var overwriteDPS float64

	var currDotTickSpeed float64

	// ------------------------------------------
	// AoE (Mind Sear)
	// ------------------------------------------

	//var remain_fight float64

	if sim.CurrentTime == 0 && spriest.rotation.PrecastVt {
		spriest.SpendMana(sim, spriest.VampiricTouch.DefaultCast.Cost, spriest.VampiricTouch.ResourceMetrics)
		spriest.VampiricTouch.SkipCastAndApplyEffects(sim, spriest.CurrentTarget)
	}

	// initialize function specific variables
	var spell *core.Spell
	var TFmod float64
	var swStacks float64
	var numswptickstime float64
	var cdDpso float64
	var cdDps float64
	var chosenMfs int
	var num_DP_ticks float64
	var num_VT_ticks float64
	var wait1 time.Duration
	var wait2 time.Duration
	var wait time.Duration
	var wait3 time.Duration

	// initialize helpful variables for calculations later
	vtCastTime := spriest.ApplyCastSpeed(time.Millisecond * 1500)
	gcd := spriest.SpellGCD()
	tickLength := spriest.MindFlayTickDuration()

	dotTickSpeed := float64(spriest.ApplyCastSpeed(time.Second * 3))
	critChance := (spriest.GetStat(stats.SpellCrit) + spriest.CurrentTarget.PseudoStats.BonusCritRatingTaken + spriest.CurrentTarget.PseudoStats.BonusSpellCritRatingTaken) / (core.CritRatingPerCritChance * 100)
	remain_fight := float64(sim.GetRemainingDuration())
	castMf2 := 0 // if SW stacks = 3, and we want to get SWP up at 5 stacks exactly, then we want to hard code a MF2
	bestIdx := -1

	// grab all of the CDs remaining durations to use in the dps calculation
	allCDs := []time.Duration{
		core.MaxDuration(0, spriest.MindBlast.TimeToReady(sim)),
		core.MaxDuration(0, spriest.DevouringPlagueDot.RemainingDuration(sim)),
		core.MaxDuration(0, spriest.VampiricTouchDot.RemainingDuration(sim)-vtCastTime),
		core.MaxDuration(0, spriest.ShadowWordDeath.TimeToReady(sim)),
		0,
	}

	rotType := spriest.rotation.RotationType

	if spriest.ShadowWeavingAura.IsActive() {
		swStacks = float64(spriest.ShadowWeavingAura.GetStacks())
	}

	if rotType == proto.ShadowPriest_Rotation_AoE {
		numTicks := 5
		spell = spriest.MindSear[numTicks]
		if success := spell.Cast(sim, spriest.CurrentTarget); !success {
			spriest.WaitForMana(sim, spell.CurCast.Cost)
		}
		return

	} else if rotType == proto.ShadowPriest_Rotation_Basic || rotType == proto.ShadowPriest_Rotation_Clipping {

		if spriest.DevouringPlagueDot.RemainingDuration(sim) <= 0 {
			bestIdx = 1
		} else if spriest.Talents.VampiricTouch && spriest.VampiricTouchDot.RemainingDuration(sim) <= vtCastTime {
			bestIdx = 2
		} else if !spriest.ShadowWordPainDot.IsActive() && swStacks >= 5 {
			bestIdx = 5
		} else if spriest.Talents.MindFlay {
			bestIdx = 4
		}
	} else {

		// if shadow word pain is active on the target, then increase damage of MB/MF by 10%
		if spriest.ShadowWordPainDot.IsActive() {
			TFmod = float64(spriest.Talents.TwistedFaith) * 0.02
		} else {
			TFmod = 0
		}

		mfglyphMod := 0.0
		if spriest.HasGlyph(int32(proto.PriestMajorGlyph_GlyphOfMindFlay)) {
			mfglyphMod = 0.1
		}

		swdmfglyphMod := 1.0
		if spriest.HasGlyph(int32(proto.PriestMajorGlyph_GlyphOfShadowWordDeath)) && sim.IsExecutePhase35() {
			swdmfglyphMod = 1.1
		}

		// Reduce number of DP/VT ticks based on remaining duration
		num_DP_ticks = math.Floor(remain_fight / dotTickSpeed)
		if num_DP_ticks > 8 {
			num_DP_ticks = 8
		}

		num_VT_ticks = math.Floor(remain_fight / dotTickSpeed)
		if num_VT_ticks > 5 {
			num_VT_ticks = 5
			if spriest.T9TwoSetBonus {
				num_VT_ticks = 7
			}
		}

		var duration time.Duration
		aura := spriest.GetActiveAuraWithTag(core.BloodlustAuraTag)
		if aura != nil {
			duration = aura.RemainingDuration(sim)
		}

		// Spell damage numbers that are updated before each cast in order to determine the most optimal next cast based on dps over a finite window
		// This is needed throughout the code to determine the optimal spell(s) to cast next
		// MB dmg
		mbDamage = (1025 + spriest.GetStat(stats.SpellPower)*(0.428*(1+float64(spriest.Talents.Misery)*0.05))) * (1 + float64(spriest.Talents.Darkness)*0.02) * (1 + TFmod) *
			core.TernaryFloat64(spriest.Talents.Shadowform, 1.15, 1) * (1 + 0.5*(critChance+float64(spriest.Talents.MindMelt)*0.02)*float64(spriest.Talents.ShadowPower)*0.2)
		if !spriest.options.UseMindBlast {
			mbDamage = 0
		}

		// DP dmg
		dpInit := ((172 + spriest.GetStat(stats.SpellPower)*0.1849) * 8.0 * float64(spriest.Talents.ImprovedDevouringPlague) * 0.1 * (1.0 + (float64(spriest.Talents.Darkness)*0.02 +
			float64(spriest.Talents.TwinDisciplines)*0.01 + float64(spriest.Talents.ImprovedDevouringPlague)*0.05)) * core.TernaryFloat64(spriest.T8TwoSetBonus, 1.15, 1) * core.TernaryFloat64(spriest.Talents.Shadowform, 1.15, 1) * (1 + 0.5*(critChance+core.TernaryFloat64(spriest.T10TwoSetBonus, 0.05, 0))))
		dpDot := ((172 + spriest.GetStat(stats.SpellPower)*0.1849) * num_DP_ticks *
			(1.0 + (float64(spriest.Talents.Darkness)*0.02 + float64(spriest.Talents.TwinDisciplines)*0.01 + float64(spriest.Talents.ImprovedDevouringPlague)*0.05 + core.TernaryFloat64(spriest.T8TwoSetBonus, 0.15, 0))) * core.TernaryFloat64(spriest.Talents.Shadowform, 1.15, 1) *
			(1 + 1*(critChance+float64(spriest.Talents.MindMelt)*0.03) + core.TernaryFloat64(spriest.T10TwoSetBonus, 0.05, 0)))
		dpDamage = dpInit + dpDot

		// VT dmg
		vtDamage = (170 + spriest.GetStat(stats.SpellPower)*0.4) * num_VT_ticks *
			(1.0 + float64(spriest.Talents.Darkness)*0.02) * core.TernaryFloat64(spriest.Talents.Shadowform, 1.15, 1) * (1 + 1*(critChance+float64(spriest.Talents.MindMelt)*0.03+core.TernaryFloat64(spriest.T10TwoSetBonus, 0.05, 0)))

		// SWD dmg
		swdDamage = (618 + spriest.GetStat(stats.SpellPower)*0.429) * (1 + 0.5*(critChance+float64(spriest.Talents.MindMelt)*0.02+core.TernaryFloat64(spriest.T7FourSetBonus, 0.1, 0))*float64(spriest.Talents.ShadowPower)*0.2) *
			(1.0 + (float64(spriest.Talents.Darkness)*0.02 + float64(spriest.Talents.TwinDisciplines)*0.01)) * core.TernaryFloat64(spriest.Talents.Shadowform, 1.15, 1) * swdmfglyphMod
		if !spriest.options.UseShadowWordDeath {
			swdDamage = 0
		}

		// MF dmg 3 ticks
		mfDamage = (588 + spriest.GetStat(stats.SpellPower)*(0.2570*3*(1+float64(spriest.Talents.Misery)*0.05))) * core.TernaryFloat64(spriest.Talents.Shadowform, 1.15, 1) * (1.0 + (float64(spriest.Talents.Darkness)*0.02 +
			float64(spriest.Talents.TwinDisciplines)*0.01)) * (1 + TFmod + mfglyphMod) * (1 + 0.5*(critChance+float64(spriest.Talents.MindMelt)*0.02+core.TernaryFloat64(spriest.T9FourSetBonus, 0.05, 0))*float64(spriest.Talents.ShadowPower)*0.2)

		// SWP is seperate because it doesnt follow the same logic for casting as the other spells
		swpTickDamage := ((230 + spriest.GetStat(stats.SpellPower)*0.1829) *
			(1.0 + float64(spriest.Talents.Darkness)*0.02 + float64(spriest.Talents.TwinDisciplines)*0.01) * core.TernaryFloat64(spriest.Talents.Shadowform, 1.15, 1) *
			(1 + 1*(critChance+float64(spriest.Talents.MindMelt)*0.03)))

		// this should be cleaned up, but essentially we want to cast SWP either 3rd or 5th in the rotation which is fight length dependent

		waitFor5 := 0   // if SW stacks = 3, and we want to get SWP up at 5 stacks exactly, then this flag gets set to 1
		castSwpNow := 0 // if SW stacks = 3, and we want to get SWP up immediately becaues fight length is low enough, then this flag gets set to 1

		if swStacks > 2 && swStacks < 5 && !spriest.ShadowWordPainDot.IsActive() {
			addedDmg := mbDamage*0.12 + mfDamage*0.22*2/3 + swpTickDamage*2*float64(gcd.Seconds())/3
			numswptickstime = addedDmg / (swpTickDamage * 0.06) * 3 //if the fight lenght is < numswptickstime then use swp 3rd.. if > then use at weaving = 5
			//
			if remain_fight*math.Pow(10, -9) < numswptickstime { //
				castSwpNow = 1
			} else {
				waitFor5 = 1
				castMf2 = 1
			}
		}

		var currDPS float64
		var nextTickWait time.Duration

		if spriest.DevouringPlagueDot.IsActive() {
			newPsuedoHaste := spriest.PseudoStats.CastSpeedMultiplier
			newHasteRating := spriest.GetStat(stats.SpellHaste)

			currDotTickSpeed = 3 / (spriest.DPstatpH * (1 + spriest.DPstatH/32.79/100))
			dotTickSpeednew := 3 / (newPsuedoHaste * (1 + newHasteRating/32.79/100))

			dpRemainTicks := float64(allCDs[dpIdx].Seconds()) / currDotTickSpeed
			nextTick := dpRemainTicks - math.Floor(dpRemainTicks)
			nextTickWait = time.Duration(nextTick * currDotTickSpeed * 1e9)

			//potmfdps := math.Floor(nextTick * currDotTickSpeed / float64(tickLength.Seconds()))

			dpInitCurr := ((172 + spriest.DPstatSp*0.1849) * 8.0 * float64(spriest.Talents.ImprovedDevouringPlague) * 0.1 * (1.0 + (float64(spriest.Talents.Darkness)*0.02 +
				float64(spriest.Talents.TwinDisciplines)*0.01 + float64(spriest.Talents.ImprovedDevouringPlague)*0.05)) * core.TernaryFloat64(spriest.T8TwoSetBonus, 1.15, 1) * core.TernaryFloat64(spriest.Talents.Shadowform, 1.15, 1) * (1 + 0.5*(critChance+core.TernaryFloat64(spriest.T10TwoSetBonus, 0.05, 0))))
			dpDotCurr := ((172 + spriest.DPstatSp*0.1849) *
				(1.0 + (float64(spriest.Talents.Darkness)*0.02 + float64(spriest.Talents.TwinDisciplines)*0.01 + float64(spriest.Talents.ImprovedDevouringPlague)*0.05 + core.TernaryFloat64(spriest.T8TwoSetBonus, 0.15, 0))) * core.TernaryFloat64(spriest.Talents.Shadowform, 1.15, 1) *
				(1 + 1*(critChance+float64(spriest.Talents.MindMelt)*0.03) + core.TernaryFloat64(spriest.T10TwoSetBonus, 0.05, 0)))

			cdDamage := mbDamage
			if spriest.T10FourSetBonus || cdDamage == 0 {
				cdDamage = mfDamage / 3 * 2
			}

			currDPS = (dpInitCurr + dpDotCurr*8 + cdDamage) / (currDotTickSpeed * 8)
			overwriteDPS = (dpInitCurr + dpInit + dpDotCurr*1 + dpDot) / (dotTickSpeednew*8 + currDotTickSpeed*1)
		}
		var currDPS2 float64
		var overwriteDPS2 float64
		//var nextTickWait2 time.Duration
		if spriest.DevouringPlagueDot.IsActive() && duration.Seconds() < 3 && duration.Seconds() > 0.1 {

			newPsuedoHaste := spriest.PseudoStats.CastSpeedMultiplier
			newHasteRating := spriest.GetStat(stats.SpellHaste)

			currDotTickSpeed = 3 / (spriest.DPstatpH * (1 + spriest.DPstatH/32.79/100))
			dpRemainTicks := 8 - float64(allCDs[dpIdx].Seconds())/currDotTickSpeed
			//nextTick := (dpRemainTicks) - math.Floor(dpRemainTicks)
			//nextTickWait2 = time.Duration(nextTick * currDotTickSpeed * 1e9)

			dpInitCurr := ((172 + spriest.DPstatSp*0.1849) * 8.0 * float64(spriest.Talents.ImprovedDevouringPlague) * 0.1 * (1.0 + (float64(spriest.Talents.Darkness)*0.02 +
				float64(spriest.Talents.TwinDisciplines)*0.01 + float64(spriest.Talents.ImprovedDevouringPlague)*0.05)) * core.TernaryFloat64(spriest.T8TwoSetBonus, 1.15, 1) * core.TernaryFloat64(spriest.Talents.Shadowform, 1.15, 1) * (1 + 0.5*(critChance+core.TernaryFloat64(spriest.T10TwoSetBonus, 0.05, 0))))
			dpDotNext := ((172 + spriest.DPstatSp*0.1849) *
				(1.0 + (float64(spriest.Talents.Darkness)*0.02 + float64(spriest.Talents.TwinDisciplines)*0.01 + float64(spriest.Talents.ImprovedDevouringPlague)*0.05 + core.TernaryFloat64(spriest.T8TwoSetBonus, 0.15, 0))) * core.TernaryFloat64(spriest.Talents.Shadowform, 1.15, 1) *
				(1 + 1*(critChance+float64(spriest.Talents.MindMelt)*0.03) + core.TernaryFloat64(spriest.T10TwoSetBonus, 0.05, 0)))

			overwriteDPS2 = dpInitCurr + dpRemainTicks*(dpDotNext-dpDotNext/(newPsuedoHaste*(1+newHasteRating/32.79/100)))
			cdDamage := mbDamage
			if spriest.T10FourSetBonus || cdDamage == 0 {
				cdDamage = mfDamage / 3 * 2
			}
			currDPS2 = cdDamage

			// if sim.Log != nil {
			// 	spriest.Log(sim, "currDPS2[%d]", currDPS2)
			// 	spriest.Log(sim, "overwriteDPS2[%d]", overwriteDPS2)
			// 	spriest.Log(sim, "dpRemainTicks[%d]", dpRemainTicks)
			// }
		}

		// Make an array of DPCT per spell that will be used to find the optimal spell to cast
		spellDPCT := []float64{
			// MB dps
			mbDamage / float64((gcd + allCDs[mbIdb]).Seconds()),
			// DP dps
			dpDamage / float64((gcd + allCDs[dpIdx]).Seconds()),
			// VT dps
			vtDamage / float64((gcd + allCDs[vtIdx]).Seconds()),
			// SWD dps
			swdDamage / float64((gcd + allCDs[swdIdx]).Seconds()),
			// MF dps 3 ticks
			mfDamage / float64((tickLength * 3).Seconds()),
		}

		if sim.Log != nil {
			//spriest.Log(sim, "mbDamage[%d]", mbDamage)
			//spriest.Log(sim, "mb time[%d]", float64((gcd + allCDs[mbIdb]).Seconds()))
			//spriest.Log(sim, "mftime[%d]", float64((tickLength * 3).Seconds()))
			//spriest.Log(sim, "gcd[%d]", gcd.Seconds())
			//spriest.Log(sim, "CastSpeedMultiplier[%d]", spriest.PseudoStats.CastSpeedMultiplier)
			//spriest.Log(sim, "critChance[%d]", critChance)
		}

		// Find the maximum DPCT spell
		bestDmg := 0.0
		for i, v := range spellDPCT {
			if sim.Log != nil {
				//spriest.Log(sim, "\tspellDPCT[%d]: %01.f", i, v)
				//spriest.Log(sim, "\tcdDiffs[%d]: %0.1f", i, allCDs[i].Seconds())
			}
			if v > bestDmg {
				bestIdx = i
				bestDmg = v
			}
		}
		// Find the minimum CD ability to make sure that shouldnt be cast first
		nextCD := core.NeverExpires
		nextIdx := -1
		for i, v := range allCDs[1 : len(allCDs)-1] {
			// if sim.Log != nil {
			//   spriest.Log(sim, "\tallCDs[%d]: %01.f", i, v)
			// 	 spriest.Log(sim, "\tcdDiffs[%d]: %0.1f", i, cdDiffs[i].Seconds())
			// }
			if v < nextCD {
				nextCD = v
				nextIdx = i + 1
			}
		}
		waitmin := nextCD

		// Now it's possible that the wait time for the chosen spell is long, if that's the case, then it might be better to investigate the dps over a 2 spell window to see if casting something else will benefit
		if bestIdx < 4 {
			currentWait = allCDs[bestIdx]
		}

		if allCDs[0] < gcd && bestIdx == 4 && allCDs[3] == 0 {
			totalDps__poss := (mbDamage + swdDamage) / float64((gcd + gcd).Seconds())
			totalDps__poss3 := (mbDamage + mfDamage*2/3) / float64((2*tickLength + gcd).Seconds())

			if totalDps__poss > totalDps__poss3 {
				bestIdx = 3
				currentWait = allCDs[bestIdx]
			}
		}

		if nextIdx != 4 && bestIdx != 4 && bestIdx != 5 && currentWait > waitmin && currentWait.Seconds() < 3 { // right now 3 might not be correct number, but we can study this to optimize

			if bestIdx == 2 { // MB VT DP SWD
				cdDpso = vtDamage / float64((gcd + currentWait).Seconds())
			} else if bestIdx == 0 {
				cdDpso = mbDamage / float64((gcd + currentWait).Seconds())
			} else if bestIdx == 3 {
				cdDpso = swdDamage / float64((gcd + currentWait).Seconds())
			} else if bestIdx == 1 {
				cdDpso = dpDamage / float64((gcd + currentWait).Seconds())
			}

			if nextIdx == 2 {
				cdDps = vtDamage / float64((gcd + waitmin).Seconds())
			} else if nextIdx == 0 {
				cdDps = mbDamage / float64((gcd + waitmin).Seconds())
			} else if nextIdx == 3 {
				cdDps = swdDamage / float64((gcd + waitmin).Seconds())
			} else if nextIdx == 1 {
				cdDps = dpDamage / float64((gcd + waitmin).Seconds())
			}

			residualWait := currentWait - gcd
			if residualWait < 0 {
				residualWait = 0
			}
			totalDps__poss0 := (cdDpso * float64((currentWait + gcd).Seconds())) / float64((gcd + currentWait).Seconds())
			totalDps__poss1 := (cdDpso*float64((currentWait+gcd).Seconds()) + cdDps*float64((waitmin+gcd).Seconds())) / float64((waitmin + gcd + gcd + residualWait).Seconds())

			totalDps__poss2 := float64(0)
			totalDps__poss3 := float64(0)

			residualMF := currentWait - 2*tickLength
			if residualMF < 0 {
				residualMF = 0
			}
			totalDps__poss2 = (cdDpso*float64((currentWait+gcd).Seconds()) + mfDamage*2/3) / float64((2*tickLength + gcd + residualMF).Seconds())
			residualMF = currentWait - 3*tickLength
			if residualMF < 0 {
				residualMF = 0
			}
			totalDps__poss3 = (cdDpso*float64((currentWait+gcd).Seconds()) + mfDamage) / float64((3*tickLength + gcd + residualMF).Seconds())

			//	if sim.Log != nil {
			//		spriest.Log(sim, "nextIdx[%d]", nextIdx)
			//		spriest.Log(sim, "bestIdx[%d]", bestIdx)
			//		spriest.Log(sim, "residualWait[%d]", residualWait.Seconds())
			//		spriest.Log(sim, "total_dps__poss0[%d]", totalDps__poss0)
			//		spriest.Log(sim, "total_dps__poss1[%d]", totalDps__poss1)
			//		spriest.Log(sim, "total_dps__poss2[%d]", totalDps__poss2)
			//		spriest.Log(sim, "total_dps__poss3[%d]", totalDps__poss3)
			//	}

			if (totalDps__poss1 > totalDps__poss0) || (totalDps__poss2 > totalDps__poss0) || (totalDps__poss3 > totalDps__poss0) {
				if totalDps__poss1 > totalDps__poss0 && totalDps__poss1 > totalDps__poss2 && totalDps__poss1 > totalDps__poss3 {
					bestIdx = nextIdx // if choosing the minimum wait time spell first is highest dps, then change the index and current wait
					currentWait = waitmin
				} else {
					bestIdx = 4
				}
			}

		}
		// Now it's possible that the wait time is > 1 gcd and is the minimum wait time.. this is unlikely in wrath given how good MF is, but still might be worth to check

		if overwriteDPS-currDPS > 200 {
			bestIdx = 1
			currentWait = time.Duration(nextTickWait)
			// if sim.Log != nil {
			// 	spriest.Log(sim, "currDPS %d", currDPS)
			// 	spriest.Log(sim, "overwriteDPS %d", overwriteDPS)
			// 	spriest.Log(sim, "currentwait %d", float64(CurrentWait.Seconds()))
			// }
		} else {
			overwriteDPS = 0
		}
		// if sim.Log != nil {
		// 	spriest.Log(sim, "best=next[%d]", bestIdx)
		// }
		// if chosen wait time is > 0.3*GCD (this was optimized in private sim, but might want to reoptimize with procs/ect) then check if it's more dps to to add a mf sequence
		if bestIdx != 4 && float64(currentWait.Seconds()) > 0.3*float64(gcd.Seconds()) {

			if bestIdx == 2 { // MB VT DP SWD
				cdDpso = vtDamage
			} else if bestIdx == 0 {
				cdDpso = mbDamage
			} else if bestIdx == 3 {
				cdDpso = swdDamage
			} else if bestIdx == 1 {
				cdDpso = dpDamage
			}

			addedgcd := core.MaxDuration(gcd, time.Duration(2)*tickLength)
			addedgcdtime := addedgcd - time.Duration(2)*tickLength

			deltaMf1 := currentWait - gcd
			if deltaMf1 < 0 {
				deltaMf1 = 0
			}
			deltaMf2 := currentWait - (tickLength*2 + addedgcdtime)
			if deltaMf2 < 0 {
				deltaMf2 = 0
			}
			deltaMf3 := currentWait - tickLength*3
			if deltaMf3 < 0 {
				deltaMf3 = 0
			}

			dpsPossibleshort := []float64{
				(cdDpso) / float64((gcd + currentWait).Seconds()),
				(cdDpso + mfDamage/3) / float64((deltaMf1 + gcd + gcd).Seconds()),
				(cdDpso + mfDamage/3*2) / float64((deltaMf2 + tickLength*2 + addedgcdtime + gcd).Seconds()),
				(cdDpso + mfDamage) / float64((deltaMf3 + tickLength*3 + gcd).Seconds()),
			}

			// Find the highest possible dps and its index
			highestPossibleIdx := 0
			highestPossibleDmg := 0.0
			if highestPossibleIdx == 0 {
				for i, v := range dpsPossibleshort {
					if v >= highestPossibleDmg {
						if sim.Log != nil {
							//spriest.Log(sim, "\thighestPossibleDmg[%d]: %01.f", i, v)
						}
						highestPossibleIdx = i
						highestPossibleDmg = v
					}
				}
			}
			mfAddIdx := highestPossibleIdx

			if mfAddIdx == 0 {
				chosenMfs = 0
			} else if mfAddIdx == 1 {
				chosenMfs = 1
			} else if mfAddIdx == 2 {
				chosenMfs = 2
			} else if mfAddIdx == 3 {
				chosenMfs = 3
			}
			if chosenMfs > 0 {
				bestIdx = 4
			}
		}

		if bestIdx == 3 && tickLength*2 <= gcd {
			if spellDPCT[3] < spellDPCT[4]*2/3 {
				bestIdx = 4
			}
		}

		if chosenMfs == 1 && allCDs[swdIdx] == 0 && swdDamage != 0 {
			if tickLength*2 <= gcd {
				bestIdx = 4
			} else {
				bestIdx = 3
				currentWait = 0
			}
		}

		if castMf2 > 0 {
			bestIdx = 4
		}
		if swStacks == 5 && !spriest.ShadowWordPainDot.IsActive() {
			bestIdx = 5
		}
		if castSwpNow > 0 {
			bestIdx = 5
		}
		if waitFor5 == 1 && swStacks == 5 {
			bestIdx = 5
		}

		if overwriteDPS-currDPS > 200 && (currentWait < gcd/2 || float64(currentWait) >= currDotTickSpeed*0.9) {
			bestIdx = 1
			currentWait = 0
		}

		if overwriteDPS-currDPS > 200 && currentWait <= gcd && currentWait >= gcd/2 && allCDs[swdIdx] == 0 {
			if tickLength*2 <= gcd {
				bestIdx = 4
			} else {
				bestIdx = 3
				currentWait = 0
			}
		}

		if overwriteDPS2-currDPS2 > 200 { //Seems to be a dps loss to overwrite a DP to snap shot
			bestIdx = 1
			currentWait = 0
		}

		if currentWait > 0 && bestIdx != 5 && bestIdx != 4 {
			spriest.WaitUntil(sim, sim.CurrentTime+currentWait)
			return
		}

	}

	if bestIdx == 0 {
		spell = spriest.MindBlast
	} else if bestIdx == 1 {
		spell = spriest.DevouringPlague
	} else if bestIdx == 2 {
		spell = spriest.VampiricTouch
	} else if bestIdx == 3 {
		spell = spriest.ShadowWordDeath
	} else if bestIdx == 5 {
		spell = spriest.ShadowWordPain // once swp is cast need a way for talents to refresh the duration
	} else if bestIdx == 4 {

		if castMf2 == 0 {
			if spriest.InnerFocus != nil && spriest.InnerFocus.IsReady(sim) {
				spriest.InnerFocus.Cast(sim, nil)
			}
		}

		var numTicks int

		if rotType == proto.ShadowPriest_Rotation_Basic || rotType == proto.ShadowPriest_Rotation_Clipping {

			if spriest.MindBlast.TimeToReady(sim) == 0 {
				spell = spriest.MindBlast
				if success := spell.Cast(sim, spriest.CurrentTarget); !success {
					spriest.WaitForMana(sim, spell.CurCast.Cost)
				}
				return
			} else if spriest.ShadowWordDeath.TimeToReady(sim) == 0 {
				spell = spriest.ShadowWordDeath
				if success := spell.Cast(sim, spriest.CurrentTarget); !success {
					spriest.WaitForMana(sim, spell.CurCast.Cost)
				}
				return
			} else {
				if rotType == proto.ShadowPriest_Rotation_Basic {
					numTicks = spriest.BasicMindflayRotation(sim, allCDs, gcd, tickLength)
				} else if rotType == proto.ShadowPriest_Rotation_Clipping {
					numTicks = spriest.ClippingMindflayRotation(sim, allCDs, gcd, tickLength)
				}
			}
		} else {
			if chosenMfs == 1 {
				numTicks = 1 // determiend above that it's more dps to add MF1, need if it's not better to enter ideal rotation instead
			} else if castMf2 == 1 {
				numTicks = 2
			} else {
				numTicks = spriest.IdealMindflayRotation(sim, allCDs, gcd, tickLength, currentWait, mfDamage, mbDamage, dpDamage, vtDamage, swdDamage, overwriteDPS) //enter the mf optimizaiton routine to optimze mf clips and for next optimal spell
			}
		}

		if numTicks == 0 {
			// Means we'd rather wait for next CD (swp, vt, etc) than start a MF cast.
			nextCD := core.NeverExpires
			for _, v := range allCDs[1 : len(allCDs)-1] {
				if v < nextCD {
					nextCD = v
				}
			}
			spriest.WaitUntil(sim, sim.CurrentTime+nextCD)
			return
		}

		spell = spriest.MindFlay[numTicks]
	} else {

		mbcd := spriest.MindBlast.TimeToReady(sim)
		swdcd := spriest.ShadowWordDeath.TimeToReady(sim)
		vtidx := spriest.VampiricTouchDot.RemainingDuration(sim) - vtCastTime
		swpidx := spriest.ShadowWordPainDot.RemainingDuration(sim)
		dpidx := spriest.DevouringPlagueDot.RemainingDuration(sim)
		wait1 = core.MinDuration(mbcd, swdcd)
		wait2 = core.MinDuration(dpidx, wait1)
		wait3 = core.MinDuration(vtidx, swpidx)
		wait = core.MinDuration(wait3, wait2)
		spriest.WaitUntil(sim, sim.CurrentTime+wait)
		return
	}
	if success := spell.Cast(sim, spriest.CurrentTarget); !success {
		spriest.WaitForMana(sim, spell.CurCast.Cost)
	} else if spell == spriest.VampiricTouch {
		spriest.VTstatpH = spriest.PseudoStats.CastSpeedMultiplier
		spriest.VTstatH = spriest.GetStat(stats.SpellHaste)
		spriest.VTstatSp = spriest.GetStat(stats.SpellPower) + spriest.GetStat(stats.ShadowSpellPower)
	} else if spell == spriest.DevouringPlague {
		spriest.DPstatpH = spriest.PseudoStats.CastSpeedMultiplier
		spriest.DPstatH = spriest.GetStat(stats.SpellHaste)
		spriest.DPstatSp = spriest.GetStat(stats.SpellPower) + spriest.GetStat(stats.ShadowSpellPower)
	}
}

// Returns the number of MF ticks to use, or 0 to wait for next CD.
func (spriest *ShadowPriest) BasicMindflayRotation(sim *core.Simulation, allCDs []time.Duration, gcd time.Duration, tickLength time.Duration) int {
	// just do MF3, never clipping
	nextCD := core.NeverExpires
	for _, v := range allCDs {
		if v < nextCD {
			nextCD = v
		}
	}
	// But don't start a MF if we can't get a single tick off.
	if nextCD < gcd {
		return 0
	} else {
		return 3
	}
}

// Returns the number of MF ticks to use, or 0 to wait for next CD.
func (spriest *ShadowPriest) IdealMindflayRotation(sim *core.Simulation, allCDs []time.Duration, gcd time.Duration, tickLength time.Duration,
	currentWait time.Duration, mfDamage, mbDamage, dpDamage, vtDamage, swdDamage, overwriteDPS float64) int {
	nextCD := core.NeverExpires
	nextIdx := -1

	newCDs := []time.Duration{
		core.MaxDuration(0, allCDs[0]),
		core.MaxDuration(0, allCDs[1]),
		core.MaxDuration(0, allCDs[2]),
	}

	for i, v := range newCDs {
		if v < nextCD {
			nextCD = v
			nextIdx = i
		}
	}

	if currentWait != 0 {
		nextCD = currentWait
	}

	var numTicks int
	numTicks_Base := 0.0
	numTicks_floored := 0.0
	if nextCD < gcd/2 {
		numTicks = 0
	} else {
		numTicks_Base = float64(nextCD.Seconds()) / float64(tickLength.Seconds())
		numTicks_floored = math.Floor(float64(nextCD.Seconds()) / float64(tickLength.Seconds()))
		numTicks = int(numTicks_Base)
	}

	AlmostAnotherTick := numTicks_Base - numTicks_floored

	if AlmostAnotherTick > 0.75 {
		numTicks += 1
	}

	mfTickDamage := mfDamage * 0.3333

	if float64(tickLength.Seconds()) < float64(gcd.Seconds())/2.9 {
		numTicks = 3
		return numTicks
	}

	//if sim.Log != nil {
	//spriest.Log(sim, "AlmostAnotherTick %d", AlmostAnotherTick)
	//spriest.Log(sim, "numTicks %d", numTicks)
	//spriest.Log(sim, "tickLength %d", tickLength.Seconds())
	//spriest.Log(sim, "nextCD %d", nextCD.Seconds())
	//spriest.Log(sim, "numTicks_Base %d", numTicks_Base)
	//spriest.Log(sim, "numTicks_floored %d", numTicks_floored)
	//}

	if numTicks < 100 && overwriteDPS == 0 { // if the code entered this loop because mf is the higest dps spell, and the number of ticks that can fit in the remaining cd time is < 1, then just cast a mf3 as it essentially fits perfectly
		// TODO: Should spriest latency be added to the second option here?

		mfTime := core.MaxDuration(gcd, time.Duration(numTicks)*tickLength)
		if numTicks == 0 {
			mfTime = core.MaxDuration(gcd, time.Duration(numTicks)*tickLength)
		}

		if sim.Log != nil {
			//spriest.Log(sim, "mfTime %d", mfTime.Seconds())
			//spriest.Log(sim, "allCDs %d", allCDs[0].Seconds())
			//spriest.Log(sim, "mf3Time %d", float64(time.Duration(3*tickLength).Seconds()))
		}
		// Amount of gap time after casting mind flay, but before each CD is available.

		cdDiffs := []time.Duration{
			core.MaxDuration(0, allCDs[0]-mfTime),
			core.MaxDuration(0, allCDs[1]-mfTime),
			core.MaxDuration(0, allCDs[2]-mfTime),
			core.MaxDuration(0, allCDs[3]-mfTime),
			0,
		}

		mfspdmg := 0.0
		if numTicks != 0 {
			mfspdmg = mfTickDamage * float64(numTicks) / float64((time.Duration(numTicks) * tickLength).Seconds())
		} else if numTicks > 3 {
			mfspdmg = mfTickDamage * float64(3) / float64((time.Duration(3) * tickLength).Seconds())
		}
		if sim.Log != nil {
			//spriest.Log(sim, "mfspdmg %d", mfspdmg)
		}
		spellDamages := []float64{
			// MB dps
			mbDamage / (gcd + cdDiffs[mbIdb]).Seconds(),
			// DP dps
			dpDamage / (gcd + cdDiffs[dpIdx]).Seconds(),
			// VT dps
			vtDamage / (gcd + cdDiffs[vtIdx]).Seconds(),
			// SWD dps
			swdDamage / (gcd + cdDiffs[swdIdx]).Seconds(),

			mfspdmg,
		}

		bestIdx := 0
		bestDmg := 0.0
		for i, v := range spellDamages {
			if sim.Log != nil {
				//spriest.Log(sim, "\tspellDamages[%d]: %01.f", i, v)
			}
			if v > bestDmg {
				bestIdx = i
				bestDmg = v
			}
		}

		//if numTicks < 1 && bestIdx == 4 {
		//	numTicks = 3
		//return numTicks
		//}

		if sim.Log != nil {
			//spriest.Log(sim, "bestIdx %d", bestIdx)
			//spriest.Log(sim, "nextIdx %d", nextIdx)
			//spriest.Log(sim, "spellDamages[bestIdx]  %d", spellDamages[bestIdx])
			//spriest.Log(sim, "spellDamages[nextIdx]  %d", spellDamages[nextIdx])
		}

		if bestIdx != nextIdx && spellDamages[nextIdx] < spellDamages[bestIdx] {
			numTicks_Base = float64(allCDs[bestIdx].Seconds()) / float64(tickLength.Seconds())
			numTicks_floored = math.Floor(float64(allCDs[bestIdx].Seconds()) / float64(tickLength.Seconds()))
			numTicks = int(numTicks_Base)

			AlmostAnotherTick := numTicks_Base - numTicks_floored

			if AlmostAnotherTick > 0.75 {
				numTicks += 1
			}

			mfTime = core.MaxDuration(gcd, time.Duration(numTicks)*tickLength)
			if numTicks > 3 && numTicks < 5 {
				addedgcd := core.MaxDuration(gcd, time.Duration(2)*tickLength)
				addedgcdtime := addedgcd - time.Duration(2)*tickLength
				mfTime = core.MaxDuration(gcd, time.Duration(numTicks)*tickLength+2*addedgcdtime)
			}
			cdDiffs = []time.Duration{
				core.MaxDuration(0, allCDs[0]-mfTime),
				core.MaxDuration(0, allCDs[1]-mfTime),
				core.MaxDuration(0, allCDs[2]-mfTime),
				core.MaxDuration(0, allCDs[3]-mfTime),
				0,
			}
			if sim.Log != nil {
				//spriest.Log(sim, "numTicks %d", numTicks)
				//spriest.Log(sim, "cdDiffs[bestIdx] %d", cdDiffs[bestIdx])
				//spriest.Log(sim, "mid_ticks2 %d", numTicks)
			}
			if float64(cdDiffs[bestIdx]) < float64(-0.33) {
				numTicks = numTicks - 1
				cdDiffs[bestIdx] = cdDiffs[bestIdx] + tickLength
			}
		}

		if numTicks < 0 {
			numTicks = 0
		}

		chosenWait := cdDiffs[bestIdx]

		if sim.Log != nil {
			//spriest.Log(sim, "numTicks %d", numTicks)
			//spriest.Log(sim, "mfTime %d", mfTime.Seconds())
			//spriest.Log(sim, "chosenWait %d", chosenWait.Seconds())
		}

		var newInd int
		if chosenWait > gcd {
			check_CDs := cdDiffs
			check_CDs[bestIdx] = time.Second * 15
			// need to find a way to sort the cdDiffs and find the next highest dps cast with lower wait time
			for i, v := range check_CDs {
				if v < nextCD {
					//nextCDc = v
					newInd = i
				}
			}
		}
		skipNext := 0
		var totalWaitCurr float64
		var numTicksAvail float64
		var remainTime1 float64
		var remainTime2 float64
		var remainTime3 float64
		var addTime1 float64
		var addTime2 float64
		var addTime3 float64

		if float64(chosenWait.Seconds()) > float64(gcd.Seconds()) && bestIdx != newInd && newInd > -1 {

			tick_var := float64(numTicks)
			if numTicks == 1 {
				totalWaitCurr = float64(chosenWait.Seconds()) - tick_var*float64(gcd.Seconds())
			} else {
				totalWaitCurr = float64(chosenWait.Seconds()) - tick_var*float64(tickLength.Seconds())
			}

			if totalWaitCurr-float64(gcd.Seconds()) <= float64(gcd.Seconds()) {
				if totalWaitCurr > float64(tickLength.Seconds()) {
					numTicksAvail = math.Floor((totalWaitCurr - float64(gcd.Seconds())) / float64(tickLength.Seconds()))
				} else {
					numTicksAvail = math.Floor((totalWaitCurr - float64(gcd.Seconds())) / float64(gcd.Seconds()))
				}
			} else {
				numTicksAvail = math.Floor((totalWaitCurr - float64(gcd.Seconds())) / float64(tickLength.Seconds()))
			}

			if numTicksAvail < 0 {
				numTicksAvail = 0
			}

			remainTime1 = totalWaitCurr - float64(tickLength.Seconds())*numTicksAvail - float64(gcd.Seconds())
			remainTime2 = totalWaitCurr - 1*float64(tickLength.Seconds())*numTicksAvail - float64(gcd.Seconds())
			remainTime3 = totalWaitCurr - 2*float64(tickLength.Seconds())*numTicksAvail - float64(gcd.Seconds())

			if remainTime1 > 0 {
				addTime1 = remainTime1
			} else {
				addTime1 = 0
			}

			if remainTime2 > 0 {
				addTime2 = remainTime2
			} else {
				addTime2 = 0
			}

			if remainTime3 > 0 {
				addTime3 = remainTime3
			} else {
				addTime3 = 0
			}

			dpsPossible0 := []float64{
				0,
				0,
				0,
			}

			cd_dpsb := spellDamages[bestIdx]
			cd_dpsn := spellDamages[newInd]

			dpsPossible0[0] = (numTicksAvail*mfTickDamage + cd_dpsb*float64(gcd.Seconds()) + cd_dpsn*float64(gcd.Seconds())) / (numTicksAvail*float64(tickLength.Seconds()) + 2*float64(gcd.Seconds()) + addTime1)
			dpsPossible0[1] = (tick_var*mfTickDamage + cd_dpsb*(float64(cdDiffs[bestIdx].Seconds())+float64(gcd.Seconds())) + cd_dpsn*(float64(cdDiffs[newInd].Seconds()))) / (tick_var*float64(tickLength.Seconds()) + (float64(cdDiffs[bestIdx].Seconds()) + float64(gcd.Seconds())) + (float64(cdDiffs[newInd].Seconds()) + addTime2))
			dpsPossible0[2] = ((tick_var+1)*mfTickDamage + cd_dpsb*(float64(cdDiffs[len(cdDiffs)-1-1].Seconds())+float64(gcd.Seconds())) + cd_dpsn*(float64(cdDiffs[len(cdDiffs)-1].Seconds())-float64(tickLength.Seconds()))) / ((tick_var+1)*float64(tickLength.Seconds()) + (float64(cdDiffs[bestIdx].Seconds()) + float64(gcd.Seconds())) + (float64(cdDiffs[newInd].Seconds()) + addTime3))

			highestPossibleDmg := 0.0
			highestPossibleIdx := -1
			if highestPossibleIdx == 0 {
				for i, v := range dpsPossible0 {

					if v >= highestPossibleDmg {
						highestPossibleIdx = i
						highestPossibleDmg = v
					}
				}
			}
			if highestPossibleIdx > 0 {
				numTicks = highestPossibleIdx + 1
			} else {
				numTicks = int(numTicksAvail)
				skipNext = 1
			}
		}

		if numTicks > 3 {
			if (allCDs[bestIdx] - time.Duration(numTicks-1)*tickLength - gcd) >= 0 {
				//if (allCDs[3]-time.Duration(numTicks-1)*tickLength <= 0) || (allCDs[0]-time.Duration(numTicks-1)*tickLength <= 0) { \\might need to readd this for later phases
				if allCDs[3]-time.Duration(numTicks-1)*tickLength <= 0 {
					numTicks = 3
					return numTicks
				}
			}
		}

		if skipNext == 0 {
			finalMFStart := math.Mod(float64(numTicks), 3) // Base ticks before adding additional
			dpsPossible := []float64{
				bestDmg, // dps with no tick and just wait
				0,
				0,
				0,
			}
			dpsDuration := float64((chosenWait + gcd).Seconds())

			highestPossibleIdx := 0
			// TODO: Modified this slightly to expand time window, but it still doesn't change dps for any tests.
			// Probably can remove this entirely (and then also the if highestPossibleIdx == 0 right after)
			if (finalMFStart == 2) && (chosenWait <= tickLength && chosenWait > (tickLength-time.Millisecond*15)) {
				highestPossibleIdx = 1 // if the wait time is equal to an extra mf tick, and there are already 2 ticks, then just add 1
			}

			if highestPossibleIdx == 0 {
				switch finalMFStart {
				case 0:
					// this means that the extra ticks will be relative to starting a new mf cast entirely
					dpsPossible[1] = (bestDmg*dpsDuration + mfDamage*1/3) / float64(gcd.Seconds()+gcd.Seconds())          // new damage for 1 extra tick
					dpsPossible[2] = (bestDmg*dpsDuration + mfDamage*2/3) / float64(2*tickLength.Seconds()+gcd.Seconds()) // new damage for 2 extra tick
					dpsPossible[3] = (bestDmg*dpsDuration + mfDamage) / float64(3*tickLength.Seconds()+gcd.Seconds())     // new damage for 3 extra tick

				case 1:
					total_check_time := 2 * tickLength

					if total_check_time < gcd {
						newDuration := float64((gcd + gcd).Seconds())
						dpsPossible[1] = (bestDmg*dpsDuration + (mfDamage * 1 / 3 * float64(finalMFStart+1))) / newDuration
					} else {
						newDuration := float64(((total_check_time - gcd) + gcd).Seconds())
						dpsPossible[1] = (bestDmg*dpsDuration + (mfDamage * 1 / 3 * float64(finalMFStart+1))) / newDuration
					}
					// % check add 2
					total_check_time2 := 2 * tickLength.Seconds()
					if total_check_time2 < gcd.Seconds() {
						dpsPossible[2] = (bestDmg*dpsDuration + (mfDamage * 1 / 3 * float64(finalMFStart+2))) / float64(gcd.Seconds()+gcd.Seconds())
					} else {
						dpsPossible[2] = (bestDmg*dpsDuration + (mfDamage * 1 / 3 * float64(finalMFStart+2))) / float64(total_check_time2+gcd.Seconds())
					}
				case 2:
					// % check add 1
					total_check_time := tickLength
					newDuration := float64((total_check_time + gcd).Seconds())
					dpsPossible[1] = (bestDmg*dpsDuration + mfDamage*1/3) / newDuration

				default:
					dpsPossible[1] = (bestDmg*dpsDuration + mfDamage*1/3) / float64(gcd.Seconds()+gcd.Seconds())
					if tickLength*2 > gcd {
						dpsPossible[2] = (bestDmg*dpsDuration + mfDamage*2/3) / float64(2*tickLength.Seconds()+gcd.Seconds())
					} else {
						dpsPossible[2] = (bestDmg*dpsDuration + mfDamage*2/3) / float64(gcd.Seconds()+gcd.Seconds())
					}
					dpsPossible[3] = (bestDmg*dpsDuration + mfDamage) / float64(3*tickLength.Seconds()+gcd.Seconds())
				}
			}

			// Find the highest possible dps and its index
			// highestPossibleIdx := 0
			highestPossibleDmg := 0.0
			if highestPossibleIdx == 0 {
				for i, v := range dpsPossible {
					if sim.Log != nil {
						//spriest.Log(sim, "\tdpsPossible[%d]: %01.f", i, v)
					}
					if v >= highestPossibleDmg {
						highestPossibleIdx = i
						highestPossibleDmg = v
					}
				}
			}

			numTicks += highestPossibleIdx
			// if sim.Log != nil {
			// 	spriest.Log(sim, "final_ticks %d", numTicks)
			// }
			if numTicks == 1 && tickLength*3 <= time.Duration(float64(gcd)*1.05) {
				numTicks = numTicks + 2
			}
			if numTicks == 1 && tickLength*2 <= gcd {
				numTicks = numTicks + 1
			}
			//  Now that the number of optimal ticks has been determined to optimize dps
			//  Now optimize mf2s and mf3s

			//if numTicks == 0 {
			//return numTicks
			//}

			if numTicks == 1 {
				numTicks = 1
			} else if numTicks == 0 {
				numTicks = 2
			} else if numTicks == 2 || numTicks == 4 {
				numTicks = 2
			} else {
				numTicks = 3
			}
		}
	} else {
		numTicks = int(nextCD / tickLength)
		if nextCD-core.MaxDuration(gcd, time.Duration(2)*tickLength) < 0 && numTicks != 0 {
			numTicks = numTicks - 1
		}
		// if sim.Log != nil {
		// 	spriest.Log(sim, "c_ticks %d", numTicks)
		// 	spriest.Log(sim, "nextCD %d", nextCD)
		// 	spriest.Log(sim, "tickLength %d", tickLength)
		// }
		if numTicks == 0 {
			// if sim.Log != nil {
			//   spriest.Log(sim, "zero ticks %d", numTicks)
			// }
			numTicks = 2
		}
		if numTicks >= 3 {
			numTicks = 3
		}
	}

	return numTicks
}

func (spriest *ShadowPriest) ClippingMindflayRotation(sim *core.Simulation, allCDs []time.Duration, gcd time.Duration, tickLength time.Duration) int {
	nextCD := core.NeverExpires
	for _, v := range allCDs[1 : len(allCDs)-1] {
		if v < nextCD {
			nextCD = v
		}
	}

	// if sim.Log != nil {
	// 	spriest.Log(sim, "<spriest> NextCD: %0.2f", nextCD.Seconds())
	// }

	// This means a CD is coming up before we could cast a single MF
	if nextCD < gcd {
		return 0
	}

	// How many ticks we have time for.
	numTicks := int((nextCD - time.Duration(spriest.rotation.Latency)) / tickLength)
	if numTicks == 1 {
		return 1
	} else if numTicks == 2 || numTicks == 4 {
		return 2
	} else {
		return 3
	}
}
