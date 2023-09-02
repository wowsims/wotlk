package core

import (
	"math"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
	googleProto "google.golang.org/protobuf/proto"
)

const DTPSReferenceStat = stats.Armor

type UnitStats struct {
	Stats       stats.Stats
	PseudoStats []float64
}

func NewUnitStats() UnitStats {
	return UnitStats{
		PseudoStats: make([]float64, stats.PseudoStatsLen),
	}
}
func (s *UnitStats) AddStat(stat stats.UnitStat, value float64) {
	if stat.IsStat() {
		s.Stats[stat.StatIdx()] += value
	} else {
		s.PseudoStats[stat.PseudoStatIdx()] += value
	}
}
func (s *UnitStats) Get(stat stats.UnitStat) float64 {
	if stat.IsStat() {
		return s.Stats[stat.StatIdx()]
	} else {
		return s.PseudoStats[stat.PseudoStatIdx()]
	}
}

func (s *UnitStats) ToProto() *proto.UnitStats {
	return &proto.UnitStats{
		Stats:       s.Stats[:],
		PseudoStats: s.PseudoStats,
	}
}

type StatWeightValues struct {
	Weights       UnitStats
	WeightsStdev  UnitStats
	EpValues      UnitStats
	EpValuesStdev UnitStats
}

func NewStatWeightValues() StatWeightValues {
	return StatWeightValues{
		Weights:       NewUnitStats(),
		WeightsStdev:  NewUnitStats(),
		EpValues:      NewUnitStats(),
		EpValuesStdev: NewUnitStats(),
	}
}

func (swv *StatWeightValues) ToProto() *proto.StatWeightValues {
	return &proto.StatWeightValues{
		Weights:       swv.Weights.ToProto(),
		WeightsStdev:  swv.WeightsStdev.ToProto(),
		EpValues:      swv.EpValues.ToProto(),
		EpValuesStdev: swv.EpValuesStdev.ToProto(),
	}
}

type StatWeightsResult struct {
	Dps    StatWeightValues
	Hps    StatWeightValues
	Tps    StatWeightValues
	Dtps   StatWeightValues
	Tmi    StatWeightValues
	PDeath StatWeightValues
}

func NewStatWeightsResult() *StatWeightsResult {
	return &StatWeightsResult{
		Dps:    NewStatWeightValues(),
		Hps:    NewStatWeightValues(),
		Tps:    NewStatWeightValues(),
		Dtps:   NewStatWeightValues(),
		Tmi:    NewStatWeightValues(),
		PDeath: NewStatWeightValues(),
	}
}

func (swr *StatWeightsResult) ToProto() *proto.StatWeightsResult {
	return &proto.StatWeightsResult{
		Dps:    swr.Dps.ToProto(),
		Hps:    swr.Hps.ToProto(),
		Tps:    swr.Tps.ToProto(),
		Dtps:   swr.Dtps.ToProto(),
		Tmi:    swr.Tmi.ToProto(),
		PDeath: swr.PDeath.ToProto(),
	}
}

func CalcStatWeight(swr *proto.StatWeightsRequest, referenceStat stats.Stat, progress chan *proto.ProgressMetrics) *StatWeightsResult {
	if swr.Player.BonusStats == nil {
		swr.Player.BonusStats = &proto.UnitStats{}
	}
	if swr.Player.BonusStats.Stats == nil {
		swr.Player.BonusStats.Stats = make([]float64, stats.Len)
	}
	if swr.Player.BonusStats.PseudoStats == nil {
		swr.Player.BonusStats.PseudoStats = make([]float64, stats.PseudoStatsLen)
	}

	raidProto := SinglePlayerRaidProto(swr.Player, swr.PartyBuffs, swr.RaidBuffs, swr.Debuffs)
	raidProto.Tanks = swr.Tanks

	simOptions := swr.SimOptions
	simOptions.SaveAllValues = true

	// Cut in half since we're doing above and below separately.
	// This number needs to be the same for the baseline sim too, so that RNG lines up perfectly.
	simOptions.Iterations /= 2

	// Make sure an RNG seed is always set because it gives more consistent results.
	// When there is no user-supplied seed it needs to be a randomly-selected seed
	// though, so that run-run differences still exist.
	if simOptions.RandomSeed == 0 {
		simOptions.RandomSeed = time.Now().UnixNano()
	}

	// Reduce variance even more by using test-level RNG controls.
	simOptions.IsTest = true

	//baseStatsResult := ComputeStats(&proto.ComputeStatsRequest{
	//	Raid: raidProto,
	//})
	//baseStats := baseStatsResult.RaidStats.Parties[0].Players[0].FinalStats

	baseSimRequest := &proto.RaidSimRequest{
		Raid:       raidProto,
		Encounter:  swr.Encounter,
		SimOptions: simOptions,
	}
	baselineResult := RunRaidSim(baseSimRequest)
	if baselineResult.ErrorResult != "" {
		// TODO: get stack trace out.
		return &StatWeightsResult{}
	}

	var waitGroup sync.WaitGroup

	// Do half the iterations with a positive, and half with a negative value for better accuracy.
	resultsLow := make([]*proto.RaidSimResult, stats.UnitStatsLen)
	resultsHigh := make([]*proto.RaidSimResult, stats.UnitStatsLen)

	var iterationsTotal int32
	var iterationsDone int32
	var simsTotal int32
	var simsCompleted int32

	concurrency := (runtime.NumCPU() - 1) * 2
	if concurrency <= 0 {
		concurrency = 2
	}

	tickets := make(chan struct{}, concurrency)
	for i := 0; i < concurrency; i++ {
		tickets <- struct{}{}
	}

	doStat := func(stat stats.UnitStat, value float64, isLow bool) {
		defer waitGroup.Done()
		// wait until we have CPU time available.
		<-tickets

		simRequest := googleProto.Clone(baseSimRequest).(*proto.RaidSimRequest)
		stat.AddToStatsProto(simRequest.Raid.Parties[0].Players[0].BonusStats, value)

		reporter := make(chan *proto.ProgressMetrics, 10)
		go RunSim(simRequest, reporter) // RunRaidSim(simRequest)

		var localIterations int32
		var errorStr string
		var simResult *proto.RaidSimResult

		for metrics := range reporter {
			atomic.AddInt32(&iterationsDone, metrics.CompletedIterations-localIterations)
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
				break
			}
		}
		// TODO: get stack trace out if final result error is set.
		if errorStr != "" {
			panic("Stat weights error: " + errorStr)
		}

		if isLow {
			resultsLow[stat] = simResult
		} else {
			resultsHigh[stat] = simResult
		}
		tickets <- struct{}{}
	}

	const defaultStatMod = 20.0
	statModsLow := make([]float64, stats.UnitStatsLen)
	statModsHigh := make([]float64, stats.UnitStatsLen)

	// Make sure reference stat is included.
	statModsLow[referenceStat] = defaultStatMod
	statModsHigh[referenceStat] = defaultStatMod

	statsToWeigh := stats.ProtoArrayToStatsList(swr.StatsToWeigh)
	for _, s := range statsToWeigh {
		stat := stats.UnitStatFromStat(s)
		statMod := defaultStatMod
		if stat.EqualsStat(stats.Expertise) {
			statMod = ExpertisePerQuarterPercentReduction
		} else if stat.EqualsStat(stats.Armor) || stat.EqualsStat(stats.BonusArmor) {
			statMod = defaultStatMod * 10
		}
		statModsHigh[stat] = statMod
		statModsLow[stat] = -statMod
	}
	for _, s := range swr.PseudoStatsToWeigh {
		stat := stats.UnitStatFromPseudoStat(s)
		statMod := 3.0
		statModsHigh[stat] = statMod
		statModsLow[stat] = -statMod
	}

	// Start all the threads.
	for i := range statModsLow {
		stat := stats.UnitStatFromIdx(i)
		if statModsLow[stat] == 0 {
			continue
		}
		waitGroup.Add(2)
		atomic.AddInt32(&iterationsTotal, swr.SimOptions.Iterations*2)
		atomic.AddInt32(&simsTotal, 2)

		go doStat(stat, statModsLow[stat], true)
		go doStat(stat, statModsHigh[stat], false)
	}

	// Wait for thread results.
	waitGroup.Wait()

	// Compute weight results.
	result := NewStatWeightsResult()
	for i := 0; i < stats.UnitStatsLen; i++ {
		stat := stats.UnitStatFromIdx(i)
		if resultsLow[stat] == nil && resultsHigh[stat] == nil {
			continue
		}

		baselinePlayer := baselineResult.RaidMetrics.Parties[0].Players[0]
		modPlayerLow := resultsLow[stat].RaidMetrics.Parties[0].Players[0]
		modPlayerHigh := resultsHigh[stat].RaidMetrics.Parties[0].Players[0]

		// Check for hard caps. Hard caps will have results identical to the baseline because RNG is fixed.
		// When we find a hard-capped stat, just skip it (will return 0).
		if modPlayerHigh.Dps.Avg == baselinePlayer.Dps.Avg && modPlayerHigh.Hps.Avg == baselinePlayer.Hps.Avg && modPlayerHigh.Tmi.Avg == baselinePlayer.Tmi.Avg {
			continue
		}

		calcWeightResults := func(baselineMetrics *proto.DistributionMetrics, modLowMetrics *proto.DistributionMetrics, modHighMetrics *proto.DistributionMetrics, weightResults *StatWeightValues) {
			var lo, hi aggregator
			if resultsLow != nil {
				for i := 0; i < int(simOptions.Iterations); i++ {
					lo.add(modLowMetrics.AllValues[i] - baselineMetrics.AllValues[i])
				}
				lo.scale(1 / statModsLow[stat])
			}
			if resultsHigh != nil {
				for i := 0; i < int(simOptions.Iterations); i++ {
					hi.add(modHighMetrics.AllValues[i] - baselineMetrics.AllValues[i])
				}
				hi.scale(1 / statModsHigh[stat])
			}

			mean, stdev := lo.merge(&hi).meanAndStdDev()
			weightResults.Weights.AddStat(stat, mean)
			weightResults.WeightsStdev.AddStat(stat, stdev)
		}

		calcWeightResults(baselinePlayer.Dps, modPlayerLow.Dps, modPlayerHigh.Dps, &result.Dps)
		calcWeightResults(baselinePlayer.Hps, modPlayerLow.Hps, modPlayerHigh.Hps, &result.Hps)
		calcWeightResults(baselinePlayer.Threat, modPlayerLow.Threat, modPlayerHigh.Threat, &result.Tps)
		calcWeightResults(baselinePlayer.Dtps, modPlayerLow.Dtps, modPlayerHigh.Dtps, &result.Dtps)
		calcWeightResults(baselinePlayer.Tmi, modPlayerLow.Tmi, modPlayerHigh.Tmi, &result.Tmi)
		meanLow := (modPlayerLow.ChanceOfDeath - baselinePlayer.ChanceOfDeath) / statModsLow[stat]
		meanHigh := (modPlayerHigh.ChanceOfDeath - baselinePlayer.ChanceOfDeath) / statModsHigh[stat]
		result.PDeath.Weights.AddStat(stat, (meanLow+meanHigh)/2)
		result.PDeath.WeightsStdev.AddStat(stat, 0)
	}

	// Compute EP results.
	for i := range statModsLow {
		stat := stats.UnitStatFromIdx(i)
		if statModsLow[stat] == 0 {
			continue
		}

		calcEpResults := func(weightResults *StatWeightValues, refStat stats.Stat) {
			if weightResults.Weights.Stats[refStat] == 0 {
				return
			}
			mean := weightResults.Weights.Get(stat) / weightResults.Weights.Stats[refStat]
			stdev := weightResults.WeightsStdev.Get(stat) / math.Abs(weightResults.Weights.Stats[refStat])
			weightResults.EpValues.AddStat(stat, mean)
			weightResults.EpValuesStdev.AddStat(stat, stdev)
		}

		calcEpResults(&result.Dps, referenceStat)
		calcEpResults(&result.Hps, referenceStat)
		calcEpResults(&result.Tps, referenceStat)
		calcEpResults(&result.Dtps, DTPSReferenceStat)
		calcEpResults(&result.Tmi, DTPSReferenceStat)
		calcEpResults(&result.PDeath, DTPSReferenceStat)
	}

	return result
}
