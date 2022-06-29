package core

import (
	"math"
	"math/rand"
	"sync"
	"sync/atomic"

	"github.com/wowsims/tbc/sim/core/proto"
	"github.com/wowsims/tbc/sim/core/stats"
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
	Tps  StatWeightValues
	Dtps StatWeightValues
}

func (swr StatWeightsResult) ToProto() *proto.StatWeightsResult {
	return &proto.StatWeightsResult{
		Dps:  swr.Dps.ToProto(),
		Tps:  swr.Tps.ToProto(),
		Dtps: swr.Dtps.ToProto(),
	}
}

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
	baselineTpsMetrics := baselineResult.RaidMetrics.Parties[0].Players[0].Threat
	baselineDtpsMetrics := baselineResult.RaidMetrics.Parties[0].Players[0].Dtps

	var waitGroup sync.WaitGroup

	// Do half the iterations with a positive, and half with a negative value for better accuracy.
	resultLow := StatWeightsResult{}
	resultHigh := StatWeightsResult{}
	dpsHistsLow := [stats.Len]map[int32]int32{}
	dpsHistsHigh := [stats.Len]map[int32]int32{}
	tpsHistsLow := [stats.Len]map[int32]int32{}
	tpsHistsHigh := [stats.Len]map[int32]int32{}
	dtpsHistsLow := [stats.Len]map[int32]int32{}
	dtpsHistsHigh := [stats.Len]map[int32]int32{}

	var iterationsTotal int32
	var iterationsDone int32
	var simsTotal int32
	var simsCompleted int32

	doStat := func(stat stats.Stat, value float64, isLow bool) {
		defer waitGroup.Done()

		simRequest := googleProto.Clone(baseSimRequest).(*proto.RaidSimRequest)
		simRequest.Raid.Parties[0].Players[0].BonusStats[stat] += value
		simRequest.SimOptions.Iterations /= 2 // Cut in half since we're doing above and below separately.

		reporter := make(chan *proto.ProgressMetrics, 10)
		go RunSim(*simRequest, reporter) // RunRaidSim(simRequest)

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
		// TODO: get stack trace out if final result error is set.
		if errorStr != "" {
			panic("Stat weights error: " + errorStr)
		}
		dpsMetrics := simResult.RaidMetrics.Parties[0].Players[0].Dps
		tpsMetrics := simResult.RaidMetrics.Parties[0].Players[0].Dps
		dtpsMetrics := simResult.RaidMetrics.Parties[0].Players[0].Dtps
		dpsDiff := (dpsMetrics.Avg - baselineDpsMetrics.Avg) / value
		tpsDiff := (tpsMetrics.Avg - baselineTpsMetrics.Avg) / value
		dtpsDiff := (dtpsMetrics.Avg - baselineDtpsMetrics.Avg) / value

		if isLow {
			resultLow.Dps.Weights[stat] = dpsDiff
			resultLow.Tps.Weights[stat] = tpsDiff
			resultLow.Dtps.Weights[stat] = dtpsDiff
			resultLow.Dps.WeightsStdev[stat] = dpsMetrics.Stdev / math.Abs(value)
			resultLow.Tps.WeightsStdev[stat] = tpsMetrics.Stdev / math.Abs(value)
			resultLow.Dtps.WeightsStdev[stat] = dtpsMetrics.Stdev / math.Abs(value)
			dpsHistsLow[stat] = dpsMetrics.Hist
			tpsHistsLow[stat] = tpsMetrics.Hist
			dtpsHistsLow[stat] = dtpsMetrics.Hist
		} else {
			resultHigh.Dps.Weights[stat] = dpsDiff
			resultHigh.Tps.Weights[stat] = tpsDiff
			resultHigh.Dtps.Weights[stat] = dtpsDiff
			resultHigh.Dps.WeightsStdev[stat] = dpsMetrics.Stdev / math.Abs(value)
			resultHigh.Tps.WeightsStdev[stat] = tpsMetrics.Stdev / math.Abs(value)
			resultHigh.Dtps.WeightsStdev[stat] = dtpsMetrics.Stdev / math.Abs(value)
			dpsHistsHigh[stat] = dpsMetrics.Hist
			tpsHistsHigh[stat] = tpsMetrics.Hist
			dtpsHistsHigh[stat] = dtpsMetrics.Hist
		}
	}

	const defaultStatMod = 50.0
	statModsLow := stats.Stats{}
	statModsHigh := stats.Stats{}

	// Make sure reference stat is included.
	statModsLow[referenceStat] = defaultStatMod
	statModsHigh[referenceStat] = defaultStatMod

	for _, stat := range statsToWeigh {
		statMod := defaultStatMod
		if stat == stats.SpellHit || stat == stats.MeleeHit || stat == stats.Expertise {
			statMod = 15
		}
		statModsHigh[stat] = statMod
		statModsLow[stat] = -statMod
	}

	for stat, _ := range statModsLow {
		if statModsLow[stat] == 0 {
			continue
		}
		waitGroup.Add(2)
		atomic.AddInt32(&iterationsTotal, swr.SimOptions.Iterations)
		atomic.AddInt32(&simsTotal, 2)

		go doStat(stats.Stat(stat), statModsLow[stat], true)
		go doStat(stats.Stat(stat), statModsHigh[stat], false)
	}

	waitGroup.Wait()

	melee2HHitCap := 9 * MeleeHitRatingPerHitChance
	if swr.Debuffs != nil && swr.Debuffs.FaerieFire == proto.TristateEffect_TristateEffectImproved {
		melee2HHitCap -= 3 * MeleeHitRatingPerHitChance
	}

	for _, stat := range statsToWeigh {
		// Check for hard caps.
		if stat == stats.SpellHit || stat == stats.MeleeHit || stat == stats.Expertise {
			if resultHigh.Dps.Weights[stat] < 0.1 {
				statModsHigh[stat] = 0
				continue
			}
		}

		// For spell/melee hit, only use the direction facing away from the nearest soft/hard cap.
		if stat == stats.SpellHit {
			if baseStats[stat] > 80 {
				statModsHigh[stat] = statModsLow[stat]
				resultHigh.Dps.Weights[stat] = resultLow.Dps.Weights[stat]
				resultHigh.Tps.Weights[stat] = resultLow.Tps.Weights[stat]
				resultHigh.Dtps.Weights[stat] = resultLow.Dtps.Weights[stat]
			}
		} else if stat == stats.MeleeHit {
			if baseStats[stat] > 30 {
				if baseStats[stat] < melee2HHitCap || baseStats[stat] > melee2HHitCap+30 {
					statModsHigh[stat] = statModsLow[stat]
					resultHigh.Dps.Weights[stat] = resultLow.Dps.Weights[stat]
					resultHigh.Tps.Weights[stat] = resultLow.Tps.Weights[stat]
					resultHigh.Dtps.Weights[stat] = resultLow.Dtps.Weights[stat]
				} else {
					statModsLow[stat] = statModsHigh[stat]
					resultLow.Dps.Weights[stat] = resultHigh.Dps.Weights[stat]
					resultLow.Tps.Weights[stat] = resultHigh.Tps.Weights[stat]
					resultLow.Dtps.Weights[stat] = resultHigh.Dtps.Weights[stat]
				}
			}
			//} else if stat == stats.Expertise {
			//	if baseStats[stat] > 20 {
			//		statModsHigh[stat] = statModsLow[stat]
			//		resultHigh.Dps.Weights[stat] = resultLow.Dps.Weights[stat]
			//		resultHigh.Tps.Weights[stat] = resultLow.Tps.Weights[stat]
			//		resultHigh.Dtps.Weights[stat] = resultLow.Dtps.Weights[stat]
			//	}
		}
	}

	result := StatWeightsResult{}
	for statIdx, _ := range statModsLow {
		stat := stats.Stat(statIdx)
		if statModsLow[stat] == 0 || statModsHigh[stat] == 0 {
			continue
		}

		result.Dps.Weights[stat] = (resultLow.Dps.Weights[stat] + resultHigh.Dps.Weights[stat]) / 2
		result.Tps.Weights[stat] = (resultLow.Tps.Weights[stat] + resultHigh.Tps.Weights[stat]) / 2
		result.Dtps.Weights[stat] = (resultLow.Dtps.Weights[stat] + resultHigh.Dtps.Weights[stat]) / 2

		result.Dps.WeightsStdev[stat] = (resultLow.Dps.WeightsStdev[stat] + resultHigh.Dps.WeightsStdev[stat]) / 2
		result.Tps.WeightsStdev[stat] = (resultLow.Tps.WeightsStdev[stat] + resultHigh.Tps.WeightsStdev[stat]) / 2
		result.Dtps.WeightsStdev[stat] = (resultLow.Dtps.WeightsStdev[stat] + resultHigh.Dtps.WeightsStdev[stat]) / 2
	}

	for statIdx, _ := range statModsLow {
		stat := stats.Stat(statIdx)
		if statModsLow[stat] == 0 || statModsHigh[stat] == 0 {
			continue
		}

		result.Dps.EpValues[stat] = result.Dps.Weights[stat] / result.Dps.Weights[referenceStat]
		result.Tps.EpValues[stat] = result.Tps.Weights[stat] / result.Tps.Weights[referenceStat]
		result.Dps.EpValuesStdev[stat] = result.Dps.WeightsStdev[stat] / math.Abs(result.Dps.Weights[referenceStat])
		result.Tps.EpValuesStdev[stat] = result.Tps.WeightsStdev[stat] / math.Abs(result.Tps.Weights[referenceStat])
		if result.Dtps.Weights[DTPSReferenceStat] != 0 {
			result.Dtps.EpValues[stat] = result.Dtps.Weights[stat] / result.Dtps.Weights[DTPSReferenceStat]
			result.Dtps.EpValuesStdev[stat] = result.Dtps.WeightsStdev[stat] / math.Abs(result.Dps.Weights[DTPSReferenceStat])
		}

		//dpsWeightStdevLow := computeStDevFromHists(swr.SimOptions.Iterations/2, statModsLow[stat], dpsHistsLow[stat], baselineDpsMetrics.Hist, nil, statModsLow[referenceStat])
		//dpsWeightStdevHigh := computeStDevFromHists(swr.SimOptions.Iterations/2, statModsHigh[stat], dpsHistsHigh[stat], baselineDpsMetrics.Hist, nil, statModsHigh[referenceStat])
		//result.Dps.WeightsStdev[stat] = (dpsWeightStdevLow + dpsWeightStdevHigh) / 2

		//dpsEpStdevLow := computeStDevFromHists(swr.SimOptions.Iterations/2, statModsLow[stat], dpsHistsLow[stat], baselineDpsMetrics.Hist, dpsHistsLow[referenceStat], statModsLow[referenceStat])
		//dpsEpStdevHigh := computeStDevFromHists(swr.SimOptions.Iterations/2, statModsHigh[stat], dpsHistsHigh[stat], baselineDpsMetrics.Hist, dpsHistsHigh[referenceStat], statModsHigh[referenceStat])
		//result.Dps.EpValuesStdev[stat] = (dpsEpStdevLow + dpsEpStdevHigh) / 2

		//tpsWeightStdevLow := computeStDevFromHists(swr.SimOptions.Iterations/2, statModsLow[stat], tpsHistsLow[stat], baselineTpsMetrics.Hist, nil, statModsLow[referenceStat])
		//tpsWeightStdevHigh := computeStDevFromHists(swr.SimOptions.Iterations/2, statModsHigh[stat], tpsHistsHigh[stat], baselineTpsMetrics.Hist, nil, statModsHigh[referenceStat])
		//result.Tps.WeightsStdev[stat] = (tpsWeightStdevLow + tpsWeightStdevHigh) / 2

		//tpsEpStdevLow := computeStDevFromHists(swr.SimOptions.Iterations/2, statModsLow[stat], tpsHistsLow[stat], baselineTpsMetrics.Hist, tpsHistsLow[referenceStat], statModsLow[referenceStat])
		//tpsEpStdevHigh := computeStDevFromHists(swr.SimOptions.Iterations/2, statModsHigh[stat], tpsHistsHigh[stat], baselineTpsMetrics.Hist, tpsHistsHigh[referenceStat], statModsHigh[referenceStat])
		//result.Tps.EpValuesStdev[stat] = (tpsEpStdevLow + tpsEpStdevHigh) / 2

		//dtpsWeightStdevLow := computeStDevFromHists(swr.SimOptions.Iterations/2, statModsLow[stat], dtpsHistsLow[stat], baselineDtpsMetrics.Hist, nil, statModsLow[DTPSReferenceStat])
		//dtpsWeightStdevHigh := computeStDevFromHists(swr.SimOptions.Iterations/2, statModsHigh[stat], dtpsHistsHigh[stat], baselineDtpsMetrics.Hist, nil, statModsHigh[DTPSReferenceStat])
		//result.Dtps.WeightsStdev[stat] = (dtpsWeightStdevLow + dtpsWeightStdevHigh) / 2

		//if result.Dtps.Weights[DTPSReferenceStat] != 0 {
		//	dtpsEpStdevLow := computeStDevFromHists(swr.SimOptions.Iterations/2, statModsLow[stat], dtpsHistsLow[stat], baselineDtpsMetrics.Hist, dtpsHistsLow[DTPSReferenceStat], statModsLow[DTPSReferenceStat])
		//	dtpsEpStdevHigh := computeStDevFromHists(swr.SimOptions.Iterations/2, statModsHigh[stat], dtpsHistsHigh[stat], baselineDtpsMetrics.Hist, dtpsHistsHigh[DTPSReferenceStat], statModsHigh[DTPSReferenceStat])
		//	result.Dtps.EpValuesStdev[stat] = (dtpsEpStdevLow + dtpsEpStdevHigh) / 2
		//}
	}

	return result
}

func computeStDevFromHists(iters int32, modValue float64, moddedStatDpsHist map[int32]int32, baselineDpsHist map[int32]int32, referenceDpsHist map[int32]int32, referenceModValue float64) float64 {
	if referenceDpsHist != nil && len(referenceDpsHist) == 1 {
		return 0
	}

	sum := 0.0
	sumSquared := 0.0
	n := iters * 10
	for i := int32(0); i < n; {
		denominator := 1.0
		if referenceDpsHist != nil {
			denominator = float64(sampleFromDpsHist(referenceDpsHist, iters)-sampleFromDpsHist(baselineDpsHist, iters)) / referenceModValue
		}

		if denominator != 0 {
			ep := (float64(sampleFromDpsHist(moddedStatDpsHist, iters)-sampleFromDpsHist(baselineDpsHist, iters)) / modValue) / denominator
			sum += ep
			sumSquared += ep * ep
			i++
		}
	}
	epAvg := sum / float64(n)
	epStDev := math.Sqrt((sumSquared / float64(n)) - (epAvg * epAvg))
	return epStDev
}

// Picks a random value from a histogram, taking into account the bucket sizes.
func sampleFromDpsHist(hist map[int32]int32, histNumSamples int32) int32 {
	if len(hist) == 1 {
		for roundedDps, _ := range hist {
			return roundedDps
		}
	}

	r := rand.Float64()
	sampleIdx := int32(math.Floor(float64(histNumSamples) * r))

	var curSampleIdx int32
	for roundedDps, count := range hist {
		curSampleIdx += count
		if curSampleIdx >= sampleIdx {
			return roundedDps
		}
	}

	panic("Invalid dps histogram")
}
