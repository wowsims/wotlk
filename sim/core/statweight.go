package core

import (
	"math"
	"sync"
	"sync/atomic"

	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
	googleProto "google.golang.org/protobuf/proto"
)

const DTPSReferenceStat = stats.Armor

type StatWeightValues struct {
	Weights       stats.Stats
	WeightsStdev  stats.Stats
	EpValues      stats.Stats
	EpValuesStdev stats.Stats
}

func (swv StatWeightValues) ToProto() *proto.StatWeightValues {
	return &proto.StatWeightValues{
		Weights:       swv.Weights[:],
		WeightsStdev:  swv.WeightsStdev[:],
		EpValues:      swv.EpValues[:],
		EpValuesStdev: swv.EpValuesStdev[:],
	}
}

type StatWeightsResult struct {
	Dps  StatWeightValues
	Hps  StatWeightValues
	Tps  StatWeightValues
	Dtps StatWeightValues
}

func (swr StatWeightsResult) ToProto() *proto.StatWeightsResult {
	return &proto.StatWeightsResult{
		Dps:  swr.Dps.ToProto(),
		Hps:  swr.Dps.ToProto(),
		Tps:  swr.Tps.ToProto(),
		Dtps: swr.Dtps.ToProto(),
	}
}

// Jooper 19/10/2022:
//		Firstly, I changed from a stating stat mod for every stat of 50 to a relative stat mod of 5% for every stat, this is due to how stats generally have
//      different orders of magnitude, which may lead to really bad results for certain stats as 50 might represent a very small delta. This is still not perfect
//      since in cases where we have very little of some stat the delta will also be quite small and therefore "ruin" the calculation, we can fix this later by
//      defining some baseline stat mod for each stat. I believe this change is an improvement still from the previous method, it has made DTPS stat weights a lot
//		more stable as for example a 50 armor increase was negligble and just gave back noise for the weight.
//
//     	Secondly, I changed the calculation from a central difference to a forward difference, this is the default convention used in SimulationCraft for stat weights.
//		Generally you look at stat weights from the point of view of addition, not subtraction, therefore I believe that while central differencing makes more sense from
//		a mathematical point of view it may not always represent the better data from the point of view of stat weights in game. Additionally, the way in which the high/low
//		results were combined (as essentially an average) didn't make particular sense to me, it should be more something like:
//
//			sw = (low*delta_low + high*delta_high)/(delta_low + delta_high)
//
// 		as this represents a weighted average from both ends. This makes sense because when checking for caps as its possible that deltas on either side were different sizes.
// 		I think, like it is done in SimC, adding an option to perform central differences (rather than it being default) may be more useful in the long run.
//
//		Finally, I added expertise cap detection, both dodge / parry caps are now taken into account (we were getting reports of weird expertise weights).

func CalcStatWeight(swr proto.StatWeightsRequest, statsToWeigh []stats.Stat, referenceStat stats.Stat, progress chan *proto.ProgressMetrics) StatWeightsResult {
	if swr.Player.BonusStats == nil {
		swr.Player.BonusStats = make([]float64, stats.Len)
	}

	raidProto := SinglePlayerRaidProto(swr.Player, swr.PartyBuffs, swr.RaidBuffs, swr.Debuffs)
	raidProto.Tanks = swr.Tanks

	simOptions := swr.SimOptions

	baseStatsResult := ComputeStats(&proto.ComputeStatsRequest{
		Raid: raidProto,
	})
	baseStats := baseStatsResult.RaidStats.Parties[0].Players[0].FinalStats

	baseSimRequest := &proto.RaidSimRequest{
		Raid:       raidProto,
		Encounter:  swr.Encounter,
		SimOptions: simOptions,
	}
	baselineResult := RunRaidSim(baseSimRequest)
	if baselineResult.ErrorResult != "" {
		// TODO: get stack trace out.
		return StatWeightsResult{}
	}
	baselineDpsMetrics := baselineResult.RaidMetrics.Parties[0].Players[0].Dps
	baselineHpsMetrics := baselineResult.RaidMetrics.Parties[0].Players[0].Hps
	baselineTpsMetrics := baselineResult.RaidMetrics.Parties[0].Players[0].Threat
	baselineDtpsMetrics := baselineResult.RaidMetrics.Parties[0].Players[0].Dtps

	var waitGroup sync.WaitGroup

	result := StatWeightsResult{}
	dpsHists := [stats.Len]map[int32]int32{}
	hpsHists := [stats.Len]map[int32]int32{}
	tpsHists := [stats.Len]map[int32]int32{}
	dtpsHists := [stats.Len]map[int32]int32{}

	var iterationsTotal int32
	var iterationsDone int32
	var simsTotal int32
	var simsCompleted int32

	doStat := func(stat stats.Stat, value float64) {
		defer waitGroup.Done()

		simRequest := googleProto.Clone(baseSimRequest).(*proto.RaidSimRequest)
		simRequest.Raid.Parties[0].Players[0].BonusStats[stat] += value

		reporter := make(chan *proto.ProgressMetrics, 10)
		go RunSim(*simRequest, reporter)

		var localIterations int32
		var errorStr string
		var simResult *proto.RaidSimResult
	statsim:
		for {
			select {
			case metrics, ok := <-reporter:
				if !ok {
					break statsim
				}
				atomic.AddInt32(&iterationsDone, (metrics.CompletedIterations - localIterations))
				localIterations = metrics.CompletedIterations
				if metrics.FinalRaidResult != nil {
					atomic.AddInt32(&simsCompleted, 1)
					simResult = metrics.FinalRaidResult
				}
				if progress != nil {
					progress <- &proto.ProgressMetrics{
						TotalIterations:     atomic.LoadInt32(&iterationsTotal),
						CompletedIterations: atomic.LoadInt32(&iterationsDone),
						CompletedSims:       atomic.LoadInt32(&simsCompleted),
						TotalSims:           atomic.LoadInt32(&simsTotal),
					}
				}
				if metrics.FinalRaidResult != nil {
					errorStr = metrics.FinalRaidResult.ErrorResult
					break statsim
				}
			}
		}
		if errorStr != "" {
			panic("Stat weights error: " + errorStr)
		}
		dpsMetrics := simResult.RaidMetrics.Parties[0].Players[0].Dps
		hpsMetrics := simResult.RaidMetrics.Parties[0].Players[0].Hps
		tpsMetrics := simResult.RaidMetrics.Parties[0].Players[0].Threat
		dtpsMetrics := simResult.RaidMetrics.Parties[0].Players[0].Dtps
		dpsDiff := (dpsMetrics.Avg - baselineDpsMetrics.Avg) / value
		hpsDiff := (hpsMetrics.Avg - baselineHpsMetrics.Avg) / value
		tpsDiff := (tpsMetrics.Avg - baselineTpsMetrics.Avg) / value
		dtpsDiff := -(dtpsMetrics.Avg - baselineDtpsMetrics.Avg) / value

		result.Dps.Weights[stat] = dpsDiff
		result.Hps.Weights[stat] = hpsDiff
		result.Tps.Weights[stat] = tpsDiff
		result.Dtps.Weights[stat] = dtpsDiff
		result.Dps.WeightsStdev[stat] = dpsMetrics.Stdev / math.Abs(value)
		result.Hps.WeightsStdev[stat] = hpsMetrics.Stdev / math.Abs(value)
		result.Tps.WeightsStdev[stat] = tpsMetrics.Stdev / math.Abs(value)
		result.Dtps.WeightsStdev[stat] = dtpsMetrics.Stdev / math.Abs(value)
		dpsHists[stat] = dpsMetrics.Hist
		hpsHists[stat] = hpsMetrics.Hist
		tpsHists[stat] = tpsMetrics.Hist
		dtpsHists[stat] = dtpsMetrics.Hist
	}

	percMod := 0.05
	statMods := stats.Stats{}
	statMods[referenceStat] = baseStats[referenceStat] * percMod

	statCap := func(stat stats.Stat, mod float64, cap float64) {
		if baseStats[stat] < cap && baseStats[stat]+mod > cap {
			modToCap := cap - baseStats[stat]
			if modToCap >= 0.5*mod {
				statMods[stat] = modToCap
			}
		}
	}

	meleeSoftCap := 8.0 * MeleeHitRatingPerHitChance
	meleeHardCap := 27.0 * MeleeHitRatingPerHitChance
	spellCap := 17.0 * SpellHitRatingPerHitChance
	if swr.Debuffs != nil && (swr.Debuffs.Misery || swr.Debuffs.FaerieFire == proto.TristateEffect_TristateEffectImproved) {
		spellCap -= 3 * SpellHitRatingPerHitChance
	}

	expSoftCap := 26.0
	expHardCap := 56.0

	for _, stat := range statsToWeigh {
		mod := baseStats[stat] * percMod
		if stat == stats.MeleeHit {
			statCap(stats.MeleeHit, mod, meleeSoftCap)
			statCap(stats.MeleeHit, mod, meleeHardCap)
		} else if stat == stats.SpellHit {
			statCap(stats.SpellHit, mod, spellCap)
		} else if stat == stats.Expertise {
			statCap(stats.Expertise, math.Floor(mod), expSoftCap)
			statCap(stats.Expertise, math.Floor(mod), expHardCap)
		}

		if statMods[stat] == 0 {
			statMods[stat] = mod
		}
	}

	for stat, _ := range statMods {
		if statMods[stat] == 0 {
			continue
		}
		waitGroup.Add(1)
		atomic.AddInt32(&iterationsTotal, swr.SimOptions.Iterations)
		atomic.AddInt32(&simsTotal, 1)

		go doStat(stats.Stat(stat), statMods[stat])
	}

	waitGroup.Wait()

	for _, stat := range statsToWeigh {
		// Check for hard caps.
		if stat == stats.SpellHit || stat == stats.MeleeHit || stat == stats.Expertise {
			if result.Dps.Weights[stat] < 0.1 {
				statMods[stat] = 0
				continue
			}
		}
	}

	for statIdx, _ := range statMods {
		stat := stats.Stat(statIdx)
		if statMods[stat] == 0 {
			continue
		}

		result.Dps.EpValues[stat] = result.Dps.Weights[stat] / result.Dps.Weights[referenceStat]
		result.Dps.EpValuesStdev[stat] = result.Dps.WeightsStdev[stat] / math.Abs(result.Dps.Weights[referenceStat])

		result.Hps.EpValues[stat] = result.Hps.Weights[stat] / result.Hps.Weights[referenceStat]
		result.Hps.EpValuesStdev[stat] = result.Hps.WeightsStdev[stat] / math.Abs(result.Hps.Weights[referenceStat])

		result.Tps.EpValues[stat] = result.Tps.Weights[stat] / result.Tps.Weights[referenceStat]
		result.Tps.EpValuesStdev[stat] = result.Tps.WeightsStdev[stat] / math.Abs(result.Tps.Weights[referenceStat])

		if result.Dtps.Weights[DTPSReferenceStat] != 0 {
			result.Dtps.EpValues[stat] = result.Dtps.Weights[stat] / result.Dtps.Weights[DTPSReferenceStat]
			result.Dtps.EpValuesStdev[stat] = result.Dtps.WeightsStdev[stat] / math.Abs(result.Dtps.Weights[DTPSReferenceStat])
		}
	}

	return result
}
