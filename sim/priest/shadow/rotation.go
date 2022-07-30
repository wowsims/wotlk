package shadow

import (
	"math"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
	"github.com/wowsims/wotlk/sim/priest"
)

// TODO: probably do something different instead of making it global?
const (
	dudidx int = iota
	mfidx
	TFmod
)

// some global variables used througout the code
var CurrentWait time.Duration
var mb_dmg float64
var dp_dmg float64
var vt_dmg float64
var swd_dmg float64
var mf_dmg float64
var overwriteDPS float64
var numTicks int
var mbidx int
var dpidx int
var vtidx int
var swdidx int
var currDotTickSpeed float64

//var remain_fight float64

func (spriest *ShadowPriest) OnGCDReady(sim *core.Simulation) {
	spriest.tryUseGCD(sim)
}

func (spriest *ShadowPriest) OnManaTick(sim *core.Simulation) {
	if spriest.FinishedWaitingForManaAndGCDReady(sim) {
		spriest.tryUseGCD(sim)
	}
}
func (spriest *ShadowPriest) tryUseGCD(sim *core.Simulation) {

	if spriest.rotation.PrecastVt && sim.CurrentTime == 0 {
		spriest.SpendMana(sim, spriest.VampiricTouch.DefaultCast.Cost, spriest.VampiricTouch.ResourceMetrics)
		spriest.VampiricTouch.SkipCastAndApplyEffects(sim, spriest.CurrentTarget)
	}

	// initialize function specific variables
	var spell *core.Spell
	var TFmod float64
	var SWstacks float64
	var numswptickstime float64
	var cd_dpso float64
	var cd_dps float64
	var chosen_mfs int
	var num_DP_ticks float64
	var num_VT_ticks float64
	var wait1 time.Duration
	var wait2 time.Duration
	var wait time.Duration
	var wait3 time.Duration
	//initialize spell indices
	mbidx = 0
	dpidx = 1
	vtidx = 2
	swdidx = 3

	// initialize helpful variables for calculations later
	vtCastTime := spriest.ApplyCastSpeed(time.Millisecond * 1500)
	gcd := spriest.SpellGCD()
	var mf_reduc_time time.Duration
	if spriest.HasSetBonus(priest.ItemSetCrimsonAcolyte, 4) {
		mf_reduc_time = time.Millisecond * 170
	}
	tickLength := spriest.ApplyCastSpeed(time.Second - mf_reduc_time)
	//if tickLength<gcd/3{
	//	tickLength = gcd/3
	//}

	DotTickSpeed := float64(spriest.ApplyCastSpeed(time.Second * 3))
	critChance := (spriest.GetStat(stats.SpellCrit) / (core.CritRatingPerCritChance * 100))
	remain_fight := float64(sim.GetRemainingDuration())
	//bosshealth := float64(sim.GetRemainingDurationPercent())
	castmf2 := 0 // if SW stacks = 3, and we want to get SWP up at 5 stacks exactly, then we want to hard code a MF2
	bestIdx := -1

	// grab all of the CDs remaining durations to use in the dps calculation
	allCDs := []time.Duration{
		spriest.MindBlast.TimeToReady(sim),
		spriest.DevouringPlagueDot.RemainingDuration(sim),
		spriest.VampiricTouchDot.RemainingDuration(sim) - vtCastTime,
		spriest.ShadowWordDeath.TimeToReady(sim),
		0,
	}
	if allCDs[mbidx] < 0 {
		allCDs[mbidx] = 0
	}
	if allCDs[vtidx] < 0 {
		allCDs[vtidx] = 0
	}
	if allCDs[dpidx] < 0 {
		allCDs[dpidx] = 0
	}
	if allCDs[swdidx] < 0 {
		allCDs[swdidx] = 0
	}
	rottype := spriest.rotation.RotationType

	if spriest.ShadowWeavingAura.IsActive() {
		SWstacks = float64(spriest.ShadowWeavingAura.GetStacks())
	}

	//	if sim.Log != nil {
	//spriest.Log(sim, "BLtimeRemain %d", BLtimeRemain)
	//}
	//}

	if rottype == proto.ShadowPriest_Rotation_Basic || rottype == proto.ShadowPriest_Rotation_Clipping {

		if spriest.DevouringPlagueDot.RemainingDuration(sim) <= 0 {
			bestIdx = 1
		} else if spriest.Talents.VampiricTouch && spriest.VampiricTouchDot.RemainingDuration(sim) <= vtCastTime {
			bestIdx = 2
		} else if !spriest.ShadowWordPainDot.IsActive() && SWstacks >= 5 {
			bestIdx = 5
			//} else if spriest.Starshards.IsReady(sim) {
			//spell = spriest.Starshards
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
		num_DP_ticks = math.Floor(remain_fight / DotTickSpeed)
		if num_DP_ticks > 8 {
			num_DP_ticks = 8
		}

		num_VT_ticks = math.Floor(remain_fight / DotTickSpeed)
		if num_VT_ticks > 5 {
			num_VT_ticks = 5
		}

		var duration time.Duration
		aura := spriest.GetActiveAuraWithTag(core.BloodlustAuraTag)
		if aura != nil {
			duration = aura.RemainingDuration(sim)
			if sim.Log != nil {
				//spriest.Log(sim, "BLtimeRemain %d", duration.Seconds())
			}
		}

		// Spell damage numbers that are updated before each cast in order to determine the most optimal next cast based on dps over a finite window
		// This is needed throughout the code to determine the optimal spell(s) to cast next
		// MB dmg
		mb_dmg = (1025 + spriest.GetStat(stats.SpellPower)*(0.428*(1+float64(spriest.Talents.Misery)*0.05))) * (1 + float64(spriest.Talents.Darkness)*0.02) * (1 + TFmod) *
			core.TernaryFloat64(spriest.Talents.Shadowform, 1.15, 1) * (1 + 1*(critChance+float64(spriest.Talents.MindMelt)*0.02))
		//mb_dmg = 0
		// DP dmg
		dp_init := ((172 + spriest.GetStat(stats.SpellPower)*0.1849) * 8.0 * float64(spriest.Talents.ImprovedDevouringPlague) * 0.1 * (1.0 + (float64(spriest.Talents.Darkness)*0.02 +
			float64(spriest.Talents.TwinDisciplines)*0.01 + float64(spriest.Talents.ImprovedDevouringPlague)*0.05)) * core.TernaryFloat64(spriest.HasSetBonus(priest.ItemSetConquerorSanct, 2), 1.15, 1) * core.TernaryFloat64(spriest.Talents.Shadowform, 1.15, 1) * (1 + 0.5*(critChance+core.TernaryFloat64(spriest.HasSetBonus(priest.ItemSetCrimsonAcolyte, 4), 0.05, 0))))
		dp_dot := ((172 + spriest.GetStat(stats.SpellPower)*0.1849) * num_DP_ticks *
			(1.0 + (float64(spriest.Talents.Darkness)*0.02 + float64(spriest.Talents.TwinDisciplines)*0.01 + float64(spriest.Talents.ImprovedDevouringPlague)*0.05 + core.TernaryFloat64(spriest.HasSetBonus(priest.ItemSetConquerorSanct, 2), 0.15, 0))) * core.TernaryFloat64(spriest.Talents.Shadowform, 1.15, 1) *
			(1 + 1*(critChance+float64(spriest.Talents.MindMelt)*0.03) + core.TernaryFloat64(spriest.HasSetBonus(priest.ItemSetCrimsonAcolyte, 4), 0.05, 0)))
		dp_dmg = dp_init + dp_dot

		// VT dmg
		vt_dmg = (170 + spriest.GetStat(stats.SpellPower)*0.4) * num_VT_ticks *
			(1.0 + float64(spriest.Talents.Darkness)*0.02) * core.TernaryFloat64(spriest.Talents.Shadowform, 1.15, 1) * (1 + 1*(critChance+float64(spriest.Talents.MindMelt)*0.03+core.TernaryFloat64(spriest.HasSetBonus(priest.ItemSetCrimsonAcolyte, 4), 0.05, 0)))

		// SWD dmg
		swd_dmg = (618 + spriest.GetStat(stats.SpellPower)*0.429) * (1 + 0.5*(critChance+float64(spriest.Talents.MindMelt)*0.02+core.TernaryFloat64(spriest.HasSetBonus(priest.ItemSetValorous, 4), 0.1, 0))*float64(spriest.Talents.ShadowPower)*0.2) *
			(1.0 + (float64(spriest.Talents.Darkness)*0.02 + float64(spriest.Talents.TwinDisciplines)*0.01)) * core.TernaryFloat64(spriest.Talents.Shadowform, 1.15, 1) * swdmfglyphMod
		//swd_dmg = 0
		// MF dmg 3 ticks
		mf_dmg = (588 + spriest.GetStat(stats.SpellPower)*(0.2570*3*(1+float64(spriest.Talents.Misery)*0.05))) * core.TernaryFloat64(spriest.Talents.Shadowform, 1.15, 1) * (1.0 + (float64(spriest.Talents.Darkness)*0.02 +
			float64(spriest.Talents.TwinDisciplines)*0.01)) * (1 + TFmod + mfglyphMod) * (1 + 1*(critChance+float64(spriest.Talents.MindMelt)*0.02))

		// SWP is seperate because it doesnt follow the same logic for casting as the other spells
		swp_Tdmg := ((230 + spriest.GetStat(stats.SpellPower)*0.1829) *
			(1.0 + float64(spriest.Talents.Darkness)*0.02 + float64(spriest.Talents.TwinDisciplines)*0.01) * core.TernaryFloat64(spriest.Talents.Shadowform, 1.15, 1) *
			(1 + 1*(critChance+float64(spriest.Talents.MindMelt)*0.03)))

		// this should be cleaned up, but essentially we want to cast SWP either 3rd or 5th in the rotation which is fight length dependent

		wait_for_5 := 0   // if SW stacks = 3, and we want to get SWP up at 5 stacks exactly, then this flag gets set to 1
		cast_SPW_now := 0 // if SW stacks = 3, and we want to get SWP up immediately becaues fight length is low enough, then this flag gets set to 1
		if sim.Log != nil {
			//spriest.Log(sim, "SWstacks %d", SWstacks)
		}
		if SWstacks > 2 && SWstacks < 5 && !spriest.ShadowWordPainDot.IsActive() {
			Added_dmg := mb_dmg*0.12 + mf_dmg*0.22*2/3 + swp_Tdmg*2*float64(gcd.Seconds())/3
			numswptickstime = Added_dmg / (swp_Tdmg * 0.06) * 3 //if the fight lenght is < numswptickstime then use swp 3rd.. if > then use at weaving = 5
			//
			if remain_fight*math.Pow(10, -9) < numswptickstime { //
				cast_SPW_now = 1
			} else {
				wait_for_5 = 1
				castmf2 = 1
			}
		}

		var currDPS float64
		var nextTickWait time.Duration

		if spriest.DevouringPlagueDot.IsActive() {
			New_psuedo_haste := spriest.PseudoStats.CastSpeedMultiplier
			New_hasteR := spriest.GetStat(stats.SpellHaste)

			currDotTickSpeed = 3 / (spriest.DPstatpH * (1 + spriest.DPstatH/32.79/100))
			DotTickSpeednew := 3 / (New_psuedo_haste * (1 + New_hasteR/32.79/100))

			dpRemainTicks := float64(allCDs[dpidx].Seconds()) / currDotTickSpeed
			nextTick := dpRemainTicks - math.Floor(dpRemainTicks)
			nextTickWait = time.Duration(nextTick * currDotTickSpeed * 1e9)

			if sim.Log != nil {
				//spriest.Log(sim, "nextTick[%d]", nextTick)
				//spriest.Log(sim, "nextTickWait[%d]", nextTickWait)
				//spriest.Log(sim, "dpRemainTicks[%d]", dpRemainTicks)
			}

			//potmfdps := math.Floor(nextTick * currDotTickSpeed / float64(tickLength.Seconds()))

			dp_init_curr := ((172 + spriest.DPstatSp*0.1849) * 8.0 * float64(spriest.Talents.ImprovedDevouringPlague) * 0.1 * (1.0 + (float64(spriest.Talents.Darkness)*0.02 +
				float64(spriest.Talents.TwinDisciplines)*0.01 + float64(spriest.Talents.ImprovedDevouringPlague)*0.05)) * core.TernaryFloat64(spriest.HasSetBonus(priest.ItemSetConquerorSanct, 2), 1.15, 1) * core.TernaryFloat64(spriest.Talents.Shadowform, 1.15, 1) * (1 + 0.5*(critChance+core.TernaryFloat64(spriest.HasSetBonus(priest.ItemSetCrimsonAcolyte, 4), 0.05, 0))))
			dp_dot_curr := ((172 + spriest.DPstatSp*0.1849) *
				(1.0 + (float64(spriest.Talents.Darkness)*0.02 + float64(spriest.Talents.TwinDisciplines)*0.01 + float64(spriest.Talents.ImprovedDevouringPlague)*0.05 + core.TernaryFloat64(spriest.HasSetBonus(priest.ItemSetConquerorSanct, 2), 0.15, 0))) * core.TernaryFloat64(spriest.Talents.Shadowform, 1.15, 1) *
				(1 + 1*(critChance+float64(spriest.Talents.MindMelt)*0.03) + core.TernaryFloat64(spriest.HasSetBonus(priest.ItemSetCrimsonAcolyte, 4), 0.05, 0)))

			cd_dmg := mb_dmg
			if spriest.HasSetBonus(priest.ItemSetCrimsonAcolyte, 4) {
				cd_dmg = mf_dmg / 3 * 2
			}

			currDPS = (dp_init_curr + dp_dot_curr*8 + cd_dmg) / (currDotTickSpeed * 8)
			overwriteDPS = (dp_init_curr + dp_init + dp_dot_curr*1 + dp_dot) / (DotTickSpeednew*8 + currDotTickSpeed*1)
			if sim.Log != nil {
			}
		}
		var currDPS2 float64
		var overwriteDPS2 float64
		//var nextTickWait2 time.Duration
		if spriest.DevouringPlagueDot.IsActive() && duration.Seconds() < 3 && duration.Seconds() > 0.1 {

			New_psuedo_haste := spriest.PseudoStats.CastSpeedMultiplier
			New_hasteR := spriest.GetStat(stats.SpellHaste)

			currDotTickSpeed = 3 / (spriest.DPstatpH * (1 + spriest.DPstatH/32.79/100))
			dpRemainTicks := 8 - float64(allCDs[dpidx].Seconds())/currDotTickSpeed
			//nextTick := (dpRemainTicks) - math.Floor(dpRemainTicks)
			//nextTickWait2 = time.Duration(nextTick * currDotTickSpeed * 1e9)

			dp_init_curr := ((172 + spriest.DPstatSp*0.1849) * 8.0 * float64(spriest.Talents.ImprovedDevouringPlague) * 0.1 * (1.0 + (float64(spriest.Talents.Darkness)*0.02 +
				float64(spriest.Talents.TwinDisciplines)*0.01 + float64(spriest.Talents.ImprovedDevouringPlague)*0.05)) * core.TernaryFloat64(spriest.HasSetBonus(priest.ItemSetConquerorSanct, 2), 1.15, 1) * core.TernaryFloat64(spriest.Talents.Shadowform, 1.15, 1) * (1 + 0.5*(critChance+core.TernaryFloat64(spriest.HasSetBonus(priest.ItemSetCrimsonAcolyte, 4), 0.05, 0))))
			dp_dot_next := ((172 + spriest.DPstatSp*0.1849) *
				(1.0 + (float64(spriest.Talents.Darkness)*0.02 + float64(spriest.Talents.TwinDisciplines)*0.01 + float64(spriest.Talents.ImprovedDevouringPlague)*0.05 + core.TernaryFloat64(spriest.HasSetBonus(priest.ItemSetConquerorSanct, 2), 0.15, 0))) * core.TernaryFloat64(spriest.Talents.Shadowform, 1.15, 1) *
				(1 + 1*(critChance+float64(spriest.Talents.MindMelt)*0.03) + core.TernaryFloat64(spriest.HasSetBonus(priest.ItemSetCrimsonAcolyte, 4), 0.05, 0)))

			overwriteDPS2 = dp_init_curr + dpRemainTicks*(dp_dot_next-dp_dot_next/(New_psuedo_haste*(1+New_hasteR/32.79/100)))
			cd_dmg := mb_dmg
			if spriest.HasSetBonus(priest.ItemSetCrimsonAcolyte, 4) {
				cd_dmg = mf_dmg / 3 * 2
			}
			currDPS2 = cd_dmg

			if sim.Log != nil {
				//spriest.Log(sim, "currDPS2[%d]", currDPS2)
				//spriest.Log(sim, "overwriteDPS2[%d]", overwriteDPS2)
				//spriest.Log(sim, "dpRemainTicks[%d]", dpRemainTicks)
			}

		}

		if spriest.VampiricTouchDot.IsActive() {

		}
		if sim.Log != nil {
			//priest.Log(sim, "spriest.PseudoStats.CastSpeedMultiplier %d", spriest.PseudoStats.CastSpeedMultiplier)
		}
		// Make an array of DPCT per spell that will be used to find the optimal spell to cast
		spellDPCT := []float64{
			// MB dps
			mb_dmg / float64((gcd + allCDs[mbidx]).Seconds()),
			// DP dps
			dp_dmg / float64((gcd + allCDs[dpidx]).Seconds()),
			// VT dps
			vt_dmg / float64((gcd + allCDs[vtidx]).Seconds()),
			// SWD dps
			swd_dmg / float64((gcd + allCDs[swdidx]).Seconds()),
			// MF dps 3 ticks
			mf_dmg / float64((tickLength * 3).Seconds()),
		}

		// Find the maximum DPCT spell
		bestDmg := 0.0
		for i, v := range spellDPCT {
			if sim.Log != nil {
				//spriest.Log(sim, "\tSpellDamages[%d]: %01.f", i, v)
				//spriest.Log(sim, "\tcdDiffs[%d]: %0.1f", i, cdDiffs[i].Seconds())
			}
			if v > bestDmg {
				bestIdx = i
				bestDmg = v
			}
		}
		if sim.Log != nil {
			//spriest.Log(sim, "best=next[%d]", bestIdx)
		}
		// Find the minimum CD ability to make sure that shouldnt be cast first
		nextCD := core.NeverExpires
		nextIdx := -1
		for i, v := range allCDs[1 : len(allCDs)-1] {
			if sim.Log != nil {
				//spriest.Log(sim, "\tallCDs[%d]: %01.f", i, v)
				//spriest.Log(sim, "\tcdDiffs[%d]: %0.1f", i, cdDiffs[i].Seconds())
			}
			if v < nextCD {
				nextCD = v
				nextIdx = i + 1
			}
		}
		var Waitmin time.Duration
		Waitmin = nextCD

		// Now it's possible that the wait time for the chosen spell is long, if that's the case, then it might be better to investigate the dps over a 2 spell window to see if casting something else will benefit
		if bestIdx < 4 {
			CurrentWait = allCDs[bestIdx]
		}
		if sim.Log != nil {
			//spriest.Log(sim, "best=next[%d]", bestIdx)
		}
		if nextIdx != 4 && bestIdx != 4 && bestIdx != 5 && CurrentWait > Waitmin && CurrentWait.Seconds() < 3 { // right now 3 might not be correct number, but we can study this to optimize

			if bestIdx == 2 { // MB VT DP SWD
				cd_dpso = vt_dmg / float64((gcd + CurrentWait).Seconds())
			} else if bestIdx == 0 {
				cd_dpso = mb_dmg / float64((gcd + CurrentWait).Seconds())
			} else if bestIdx == 3 {
				cd_dpso = swd_dmg / float64((gcd + CurrentWait).Seconds())
			} else if bestIdx == 1 {
				cd_dpso = dp_dmg / float64((gcd + CurrentWait).Seconds())
			}

			if nextIdx == 2 {
				cd_dps = vt_dmg / float64((gcd + Waitmin).Seconds())
			} else if nextIdx == 0 {
				cd_dps = mb_dmg / float64((gcd + Waitmin).Seconds())
			} else if nextIdx == 3 {
				cd_dps = swd_dmg / float64((gcd + Waitmin).Seconds())
			} else if nextIdx == 1 {
				cd_dps = dp_dmg / float64((gcd + Waitmin).Seconds())
			}

			residual_wait := CurrentWait - gcd
			if residual_wait < 0 {
				residual_wait = 0
			}
			total_dps__poss0 := (cd_dpso * float64((CurrentWait + gcd).Seconds())) / float64((gcd + CurrentWait).Seconds())
			total_dps__poss1 := (cd_dpso*float64((CurrentWait+gcd).Seconds()) + cd_dps*float64((Waitmin+gcd).Seconds())) / float64((Waitmin + gcd + gcd + residual_wait).Seconds())

			total_dps__poss2 := float64(0)
			total_dps__poss3 := float64(0)
			residual_MF := CurrentWait - CurrentWait

			residual_MF = CurrentWait - 2*tickLength
			if residual_MF < 0 {
				residual_MF = 0
			}
			total_dps__poss2 = (cd_dpso*float64((CurrentWait+gcd).Seconds()) + mf_dmg*2/3) / float64((2*tickLength + gcd + residual_MF).Seconds())
			residual_MF = CurrentWait - 3*tickLength
			if residual_MF < 0 {
				residual_MF = 0
			}
			total_dps__poss3 = (cd_dpso*float64((CurrentWait+gcd).Seconds()) + mf_dmg) / float64((3*tickLength + gcd + residual_MF).Seconds())
			if sim.Log != nil {
				//spriest.Log(sim, "total_dps__poss0[%d]", total_dps__poss0)
				//spriest.Log(sim, "total_dps__poss1[%d]", total_dps__poss1)
				//spriest.Log(sim, "total_dps__poss2[%d]", total_dps__poss2)
				//spriest.Log(sim, "total_dps__poss3[%d]", total_dps__poss3)
			}

			if total_dps__poss1 > total_dps__poss0 {
				if total_dps__poss2 > total_dps__poss1 { // check if it's better to cast MF instead of minimum wait time spell
					bestIdx = 4
				} else if total_dps__poss3 > total_dps__poss1 && bestIdx != 4 {
					bestIdx = 4
				} else {
					bestIdx = nextIdx // if choosing the minimum wait time spell first is highest dps, then change the index and current wait
					CurrentWait = Waitmin
					if sim.Log != nil {
						//spriest.Log(sim, "best=next[%d]", bestIdx)
					}
				}
			}

		}
		// Now it's possible that the wait time is > 1 gcd and is the minimum wait time.. this is unlikely in wrath given how good MF is, but still might be worth to check

		if overwriteDPS-currDPS > 200 {
			bestIdx = 1
			CurrentWait = time.Duration(nextTickWait)
			if sim.Log != nil {
				//spriest.Log(sim, "currDPS %d", currDPS)
				//spriest.Log(sim, "overwriteDPS %d", overwriteDPS)
				//spriest.Log(sim, "currentwait %d", float64(CurrentWait.Seconds()))
			}
		} else {
			overwriteDPS = 0
		}
		if sim.Log != nil {
			//spriest.Log(sim, "best=next[%d]", bestIdx)
		}
		// if chosen wait time is > 0.3*GCD (this was optimized in private sim, but might want to reoptimize with procs/ect) then check if it's more dps to to add a mf sequence
		if bestIdx != 4 && float64(CurrentWait.Seconds()) > 0.4*float64(gcd.Seconds()) {

			if bestIdx == 2 { // MB VT DP SWD
				cd_dpso = vt_dmg
			} else if bestIdx == 0 {
				cd_dpso = mb_dmg
			} else if bestIdx == 3 {
				cd_dpso = swd_dmg
			} else if bestIdx == 1 {
				cd_dpso = dp_dmg
			}

			addedgcd := core.MaxDuration(gcd, time.Duration(2)*tickLength)
			addedgcdtime := addedgcd - time.Duration(2)*tickLength

			delta_1mf := CurrentWait - gcd
			if delta_1mf < 0 {
				delta_1mf = 0
			}
			delta_2mf := CurrentWait - (tickLength*2 + addedgcdtime)
			if delta_2mf < 0 {
				delta_2mf = 0
			}
			delta_3mf := CurrentWait - tickLength*3
			if delta_3mf < 0 {
				delta_3mf = 0
			}

			dpsPossibleshort := []float64{
				(cd_dpso) / float64((gcd + CurrentWait).Seconds()),
				(cd_dpso + mf_dmg/3) / float64((delta_1mf + gcd + gcd).Seconds()),
				(cd_dpso + mf_dmg/3*2) / float64((delta_2mf + tickLength*2 + addedgcdtime + gcd).Seconds()),
				(cd_dpso + mf_dmg) / float64((delta_3mf + tickLength*3 + gcd).Seconds()),
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
			MFaddIdx := highestPossibleIdx

			if MFaddIdx == 0 {
				chosen_mfs = 0
			} else if MFaddIdx == 1 {
				chosen_mfs = 1
			} else if MFaddIdx == 2 {
				chosen_mfs = 2
			} else if MFaddIdx == 3 {
				chosen_mfs = 3
			}
			if chosen_mfs > 0 {
				bestIdx = 4
			}
		}

		if bestIdx == 3 && tickLength*2 <= gcd {
			if spellDPCT[3] < spellDPCT[4]*2/3 {
				bestIdx = 4
			}
		}

		if chosen_mfs == 1 && allCDs[swdidx] == 0 && swd_dmg != 0 {
			if tickLength*2 <= gcd {
				bestIdx = 4
			} else {
				bestIdx = 3
				CurrentWait = 0
			}
		}

		if castmf2 > 0 {
			bestIdx = 4
		}
		if SWstacks == 5 && !spriest.ShadowWordPainDot.IsActive() {
			bestIdx = 5
		}
		if cast_SPW_now > 0 {
			bestIdx = 5
		}
		if wait_for_5 == 1 && SWstacks == 5 {
			bestIdx = 5
		}

		if overwriteDPS-currDPS > 200 && (CurrentWait < gcd/2 || float64(CurrentWait) >= currDotTickSpeed*0.9) {
			bestIdx = 1
			CurrentWait = 0
		}

		if overwriteDPS-currDPS > 200 && CurrentWait <= gcd && CurrentWait >= gcd/2 && allCDs[swdidx] == 0 {
			if tickLength*2 <= gcd {
				bestIdx = 4
			} else {
				bestIdx = 3
				CurrentWait = 0
			}
		}

		if overwriteDPS2-currDPS2 > 200 { //Seems to be a dps loss to overwrite a DP to snap shot
			bestIdx = 1
			CurrentWait = 0
		}
		if sim.Log != nil {
		}

		if CurrentWait > 0 && bestIdx != 5 && bestIdx != 4 {
			spriest.WaitUntil(sim, sim.CurrentTime+CurrentWait)
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

		if castmf2 == 0 {
			if spriest.InnerFocus != nil && spriest.InnerFocus.IsReady(sim) {
				spriest.InnerFocus.Cast(sim, nil)
			}
		}

		var numTicks int

		if rottype == proto.ShadowPriest_Rotation_Basic || rottype == proto.ShadowPriest_Rotation_Clipping {

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
				if rottype == proto.ShadowPriest_Rotation_Basic {
					numTicks = spriest.BasicMindflayRotation(sim, allCDs, gcd, tickLength)
				} else if rottype == proto.ShadowPriest_Rotation_Clipping {
					numTicks = spriest.ClippingMindflayRotation(sim, allCDs, gcd, tickLength)
				}
			}
		} else {
			if chosen_mfs == 1 {
				numTicks = 1 // determiend above that it's more dps to add MF1, need if it's not better to enter ideal rotation instead
			} else if castmf2 == 1 {
				numTicks = 2
			} else {
				numTicks = spriest.IdealMindflayRotation(sim, allCDs, gcd, tickLength) //enter the mf optimizaiton routine to optimze mf clips and for next optimal spell
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
func (spriest *ShadowPriest) IdealMindflayRotation(sim *core.Simulation, allCDs []time.Duration, gcd time.Duration, tickLength time.Duration) int {
	nextCD := core.NeverExpires
	nextIdx := -1
	for i, v := range allCDs {
		if v < nextCD {
			nextCD = v
			nextIdx = i
		}
	}

	if CurrentWait != 0 {
		nextCD = CurrentWait
	}

	var numTicks int
	if nextCD < gcd {
		numTicks = 0
	} else {
		numTicks = int(nextCD / tickLength)
	}

	mfDamage := mf_dmg * 0.3333

	if sim.Log != nil {
		//spriest.Log(sim, "nextCD %d", nextCD)
	}

	if numTicks < 100 && overwriteDPS == 0 { // if the code entered this loop because mf is the higest dps spell, and the number of ticks that can fit in the remaining cd time is < 1, then just cast a mf3 as it essentially fits perfectly
		// TODO: Should spriest latency be added to the second option here?
		mfTime := core.MaxDuration(gcd, time.Duration(numTicks)*tickLength)
		if numTicks == 0 {
			mfTime = core.MaxDuration(gcd, time.Duration(3)*tickLength)
		}
		// Amount of gap time after casting mind flay, but before each CD is available.
		//fmt.Println("numTicks_Start", numTicks)
		//fmt.Println("mfTime", mfTime)
		if sim.Log != nil {
			//spriest.Log(sim, "start_ticks %d", numTicks)
		}

		cdDiffs := []time.Duration{
			allCDs[0] - mfTime,
			allCDs[1] - mfTime,
			allCDs[2] - mfTime,
			allCDs[3] - mfTime,
		}
		if cdDiffs[mbidx] < 0 {
			cdDiffs[mbidx] = 0
		}
		if cdDiffs[vtidx] < 0 {
			cdDiffs[vtidx] = 0
		}
		if cdDiffs[dpidx] < 0 {
			cdDiffs[dpidx] = 0
		}
		if cdDiffs[swdidx] < 0 {
			cdDiffs[swdidx] = 0
		}

		spellDamages := []float64{
			// MB dps
			mb_dmg / (gcd + cdDiffs[mbidx]).Seconds(),
			// DP dps
			dp_dmg / (gcd + cdDiffs[dpidx]).Seconds(),
			// VT dps
			vt_dmg / (gcd + cdDiffs[vtidx]).Seconds(),
			// SWD dps
			swd_dmg / (gcd + cdDiffs[swdidx]).Seconds(),

			mf_dmg / ((tickLength * 3).Seconds()),
		}

		bestIdx := 0
		bestDmg := 0.0
		for i, v := range spellDamages {
			if v > bestDmg {
				if sim.Log != nil {
					spriest.Log(sim, "\tspellDamages[%d]: %01.f", i, v)
				}
				bestIdx = i
				bestDmg = v
			}
		}

		if numTicks < 1 && bestIdx == 4 {
			numTicks = 3
			return numTicks
		}

		if bestIdx != nextIdx && spellDamages[nextIdx] < spellDamages[bestIdx] {
			numTicks = int(allCDs[bestIdx] / tickLength)
			mfTime = core.MaxDuration(gcd, time.Duration(numTicks)*tickLength)
			if numTicks > 3 && numTicks < 5 {
				addedgcd := core.MaxDuration(gcd, time.Duration(2)*tickLength)
				addedgcdtime := addedgcd - time.Duration(2)*tickLength
				mfTime = core.MaxDuration(gcd, time.Duration(numTicks)*tickLength+2*addedgcdtime)
			}
			cdDiffs = []time.Duration{
				allCDs[0] - mfTime,
				allCDs[1] - mfTime,
				allCDs[2] - mfTime,
				allCDs[3] - mfTime,
			}
			if sim.Log != nil {
				//spriest.Log(sim, "cdDiffs[bestIdx] %d", cdDiffs[bestIdx])
				//spriest.Log(sim, "mid_ticks2 %d", numTicks)
			}
			if float64(cdDiffs[bestIdx]) < float64(-0.33) {
				numTicks = numTicks - 1
				cdDiffs[bestIdx] = cdDiffs[bestIdx] + tickLength
			}
		}

		chosenWait := cdDiffs[bestIdx]

		var new_ind int
		if chosenWait > gcd {
			check_CDs := cdDiffs
			check_CDs[bestIdx] = time.Second * 15
			// need to find a way to sort the cdDiffs and find the next highest dps cast with lower wait time
			for i, v := range check_CDs {
				if v < nextCD {
					//nextCDc = v
					new_ind = i
				}
			}
		}
		skip_next := 0
		var total_wait_curr float64
		var num_ticks_avail float64
		var remain_time1 float64
		var remain_time2 float64
		var remain_time3 float64
		var add_time1 float64
		var add_time2 float64
		var add_time3 float64

		if float64(chosenWait.Seconds()) > float64(gcd.Seconds()) && bestIdx != new_ind && new_ind > -1 {

			tick_var := float64(numTicks)
			if numTicks == 1 {
				total_wait_curr = float64(chosenWait.Seconds()) - tick_var*float64(gcd.Seconds())
			} else {
				total_wait_curr = float64(chosenWait.Seconds()) - tick_var*float64(tickLength.Seconds())
			}

			if total_wait_curr-float64(gcd.Seconds()) <= float64(gcd.Seconds()) {
				if total_wait_curr > float64(tickLength.Seconds()) {
					num_ticks_avail = math.Floor((total_wait_curr - float64(gcd.Seconds())) / float64(tickLength.Seconds()))
				} else {
					num_ticks_avail = math.Floor((total_wait_curr - float64(gcd.Seconds())) / float64(gcd.Seconds()))
				}
			} else {
				num_ticks_avail = math.Floor((total_wait_curr - float64(gcd.Seconds())) / float64(tickLength.Seconds()))
			}

			if num_ticks_avail < 0 {
				num_ticks_avail = 0
			}

			remain_time1 = total_wait_curr - float64(tickLength.Seconds())*num_ticks_avail - float64(gcd.Seconds())
			remain_time2 = total_wait_curr - 1*float64(tickLength.Seconds())*num_ticks_avail - float64(gcd.Seconds())
			remain_time3 = total_wait_curr - 2*float64(tickLength.Seconds())*num_ticks_avail - float64(gcd.Seconds())

			if remain_time1 > 0 {
				add_time1 = remain_time1
			} else {
				add_time1 = 0
			}

			if remain_time2 > 0 {
				add_time2 = remain_time2
			} else {
				add_time2 = 0
			}

			if remain_time3 > 0 {
				add_time3 = remain_time3
			} else {
				add_time3 = 0
			}

			dpsPossible0 := []float64{
				0,
				0,
				0,
			}

			cd_dpsb := spellDamages[bestIdx]
			cd_dpsn := spellDamages[new_ind]

			dpsPossible0[0] = (num_ticks_avail*mfDamage + cd_dpsb*float64(gcd.Seconds()) + cd_dpsn*float64(gcd.Seconds())) / (num_ticks_avail*float64(tickLength.Seconds()) + 2*float64(gcd.Seconds()) + add_time1)
			dpsPossible0[1] = (tick_var*mfDamage + cd_dpsb*(float64(cdDiffs[bestIdx].Seconds())+float64(gcd.Seconds())) + cd_dpsn*(float64(cdDiffs[new_ind].Seconds()))) / (tick_var*float64(tickLength.Seconds()) + (float64(cdDiffs[bestIdx].Seconds()) + float64(gcd.Seconds())) + (float64(cdDiffs[new_ind].Seconds()) + add_time2))
			dpsPossible0[2] = ((tick_var+1)*mfDamage + cd_dpsb*(float64(cdDiffs[len(cdDiffs)-1-1].Seconds())+float64(gcd.Seconds())) + cd_dpsn*(float64(cdDiffs[len(cdDiffs)-1].Seconds())-float64(tickLength.Seconds()))) / ((tick_var+1)*float64(tickLength.Seconds()) + (float64(cdDiffs[bestIdx].Seconds()) + float64(gcd.Seconds())) + (float64(cdDiffs[new_ind].Seconds()) + add_time3))

			highestPossibleDmg := 0.0
			highestPossibleIdx := -1
			if highestPossibleIdx == 0 {
				for i, v := range dpsPossible0 {
					if sim.Log != nil {
						//spriest.Log(sim, "\tdpsPossible[%d]: %01.f", i, v)
					}
					if v >= highestPossibleDmg {
						highestPossibleIdx = i
						highestPossibleDmg = v
					}
				}
			}
			if highestPossibleIdx > 0 {
				numTicks = highestPossibleIdx + 1
			} else {
				numTicks = int(num_ticks_avail)
				skip_next = 1
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
		if sim.Log != nil {
			spriest.Log(sim, "cdDiffs[bestIdx] %d", cdDiffs[bestIdx].Seconds())
			spriest.Log(sim, "tickLength %d", tickLength)
			spriest.Log(sim, "numTicks %d", numTicks)
		}
		if skip_next == 0 {
			finalMFStart := math.Mod(float64(numTicks), 3) // Base ticks before adding additional
			if sim.Log != nil {
				spriest.Log(sim, "finalMFStart %d", finalMFStart)
			}
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
					dpsPossible[1] = (bestDmg*dpsDuration + mf_dmg*1/3) / float64(gcd+gcd)          // new damage for 1 extra tick
					dpsPossible[2] = (bestDmg*dpsDuration + mf_dmg*2/3) / float64(2*tickLength+gcd) // new damage for 2 extra tick
					dpsPossible[3] = (bestDmg*dpsDuration + mf_dmg) / float64(3*tickLength+gcd)     // new damage for 3 extra tick
				case 1:
					total_check_time := 2 * tickLength

					if total_check_time < gcd {
						newDuration := float64((gcd + gcd).Seconds())
						dpsPossible[1] = (bestDmg*dpsDuration + (mf_dmg * 1 / 3 * float64(finalMFStart+1))) / newDuration
					} else {
						newDuration := float64(((total_check_time - gcd) + gcd).Seconds())
						dpsPossible[1] = (bestDmg*dpsDuration + (mf_dmg * 1 / 3 * float64(finalMFStart+1))) / newDuration
					}
					// % check add 2
					total_check_time2 := 2 * tickLength
					if total_check_time2 < gcd {
						dpsPossible[2] = (bestDmg*dpsDuration + (mf_dmg * 1 / 3 * float64(finalMFStart+2))) / float64(gcd+gcd)
					} else {
						dpsPossible[2] = (bestDmg*dpsDuration + (mf_dmg * 1 / 3 * float64(finalMFStart+2))) / float64(total_check_time2+gcd)
					}
				case 2:
					// % check add 1
					total_check_time := tickLength
					newDuration := float64((total_check_time + gcd).Seconds())
					dpsPossible[1] = (bestDmg*dpsDuration + mf_dmg*1/3) / newDuration

				default:
					dpsPossible[1] = (bestDmg*dpsDuration + mf_dmg*1/3) / float64(gcd+gcd)
					if tickLength*2 > gcd {
						dpsPossible[2] = (bestDmg*dpsDuration + mf_dmg*2/3) / float64(2*tickLength+gcd)
					} else {
						dpsPossible[2] = (bestDmg*dpsDuration + mf_dmg*2/3) / float64(gcd+gcd)
					}
					dpsPossible[3] = (bestDmg*dpsDuration + mf_dmg) / float64(3*tickLength+gcd)
				}
			}

			// Find the highest possible dps and its index
			// highestPossibleIdx := 0
			highestPossibleDmg := 0.0
			if highestPossibleIdx == 0 {
				for i, v := range dpsPossible {
					if sim.Log != nil {
						//	spriest.Log(sim, "\tdpsPossible[%d]: %01.f", i, v)
					}
					if v >= highestPossibleDmg {
						highestPossibleIdx = i
						highestPossibleDmg = v
					}
				}
			}
			numTicks += highestPossibleIdx
			if sim.Log != nil {
				spriest.Log(sim, "final_ticks %d", numTicks)
			}
			if numTicks == 1 && tickLength*3 <= time.Duration(float64(gcd)*1.05) {
				numTicks = numTicks + 2
			}
			if numTicks == 1 && tickLength*2 <= gcd {
				numTicks = numTicks + 1
			}
			if sim.Log != nil {
				spriest.Log(sim, "final_ticks %d", numTicks)
			}
			//  Now that the number of optimal ticks has been determined to optimize dps
			//  Now optimize mf2s and mf3s
			if numTicks == 1 {
				numTicks = 1
			} else if numTicks == 2 || numTicks == 4 {
				numTicks = 2
			} else {
				numTicks = 3
			}
		}
	} else {
		numTicks = int(nextCD / tickLength)
		if nextCD-core.MaxDuration(gcd, time.Duration(2)*tickLength) < 0 {
			numTicks = numTicks - 1
		}
		if sim.Log != nil {
			//spriest.Log(sim, "c_ticks %d", numTicks)
			//spriest.Log(sim, "nextCD %d", nextCD)
			//spriest.Log(sim, "tickLength %d", tickLength)
		}
		if numTicks == 0 {
			if sim.Log != nil {
				//spriest.Log(sim, "zero ticks %d", numTicks)
			}
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

	if sim.Log != nil {
		//spriest.Log(sim, "<spriest> NextCD: %0.2f", nextCD.Seconds())
	}

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
