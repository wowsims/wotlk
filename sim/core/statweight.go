package core

import (
	"math"
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

func (us UnitStats) ToProto() *proto.UnitStats {
	return &proto.UnitStats{
		Stats:       us.Stats[:],
		PseudoStats: us.PseudoStats,
	}
}

type StatWeightValues struct {
	Weights       UnitStats
	WeightsStdev  UnitStats
	EpValues      UnitStats
	EpValuesStdev UnitStats
}

func (swv StatWeightValues) ToProto() *proto.StatWeightValues {
	return &proto.StatWeightValues{
		Weights:       swv.Weights.ToProto(),
		WeightsStdev:  swv.WeightsStdev.ToProto(),
		EpValues:      swv.EpValues.ToProto(),
		EpValuesStdev: swv.EpValuesStdev.ToProto(),
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
		Hps:  swr.Hps.ToProto(),
		Tps:  swr.Tps.ToProto(),
		Dtps: swr.Dtps.ToProto(),
	}
}

func CalcStatWeight(swr *proto.StatWeightsRequest, referenceStat stats.Stat, progress chan *proto.ProgressMetrics) StatWeightsResult {
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

	// Make sure a RNG seed is always set because it gives more consistent results.
	// When there is no user-supplied seed it needs to be a randomly-selected seed
	// though, so that run-run differences still exist.
	if simOptions.RandomSeed == 0 {
		simOptions.RandomSeed = time.Now().UnixNano()
	}

	// Reduce variance even more by using test-level RNG controls.
	simOptions.IsTest = true

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

	var waitGroup sync.WaitGroup

	// Do half the iterations with a positive, and half with a negative value for better accuracy.
	resultsLow := make([]*proto.RaidSimResult, stats.UnitStatsLen)
	resultsHigh := make([]*proto.RaidSimResult, stats.UnitStatsLen)

	var iterationsTotal int32
	var iterationsDone int32
	var simsTotal int32
	var simsCompleted int32

	doStat := func(stat stats.UnitStat, value float64, isLow bool) {
		defer waitGroup.Done()

		simRequest := googleProto.Clone(baseSimRequest).(*proto.RaidSimRequest)
		if stat.IsStat() {
			simRequest.Raid.Parties[0].Players[0].BonusStats.Stats[stat.StatIdx()] += value
		} else {
			simRequest.Raid.Parties[0].Players[0].BonusStats.PseudoStats[stat.PseudoStatIdx()] += value
		}

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
			resultsLow[stat.Idx()] = simResult
		} else {
			resultsHigh[stat.Idx()] = simResult
		}
	}

	// Melee hit cap is 8% in WoTLK
	melee2HHitCap := 8 * MeleeHitRatingPerHitChance
	// Spell hit cap is 17% in WoTLK
	spellHitCap := 17 * SpellHitRatingPerHitChance
	if swr.Debuffs != nil && (swr.Debuffs.Misery || swr.Debuffs.FaerieFire == proto.TristateEffect_TristateEffectImproved) {
		spellHitCap -= 3 * SpellHitRatingPerHitChance
	}

	const defaultStatMod = 20.0
	const meleeHitStatMod = defaultStatMod
	const spellHitStatMod = defaultStatMod
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
			// Expertise is non-linear, so adjust in increments that match the stepwise reduction.
			statMod = ExpertisePerQuarterPercentReduction
		} else if stat.EqualsStat(stats.Armor) {
			statMod = defaultStatMod * 10
		}
		statModsHigh[stat.Idx()] = statMod
		statModsLow[stat.Idx()] = -statMod
	}
	for _, s := range swr.PseudoStatsToWeigh {
		stat := stats.UnitStatFromPseudoStat(s)
		statMod := 3.0
		statModsHigh[stat.Idx()] = statMod
		statModsLow[stat.Idx()] = -statMod
	}

	// Start all the threads.
	for i := range statModsLow {
		stat := stats.UnitStatFromIdx(i)
		if statModsLow[stat.Idx()] == 0 {
			continue
		}
		waitGroup.Add(2)
		atomic.AddInt32(&iterationsTotal, swr.SimOptions.Iterations*2)
		atomic.AddInt32(&simsTotal, 2)

		go doStat(stat, statModsLow[stat.Idx()], true)
		go doStat(stat, statModsHigh[stat.Idx()], false)
	}

	// Wait for thread results.
	waitGroup.Wait()

	// Compute weight results.
	result := StatWeightsResult{}
	for i := 0; i < stats.UnitStatsLen; i++ {
		stat := stats.UnitStatFromIdx(i)
		if resultsLow[stat.Idx()] == nil && resultsHigh[stat.Idx()] == nil {
			continue
		}

		baselinePlayer := baselineResult.RaidMetrics.Parties[0].Players[0]
		modPlayerLow := resultsLow[stat].RaidMetrics.Parties[0].Players[0]
		modPlayerHigh := resultsHigh[stat].RaidMetrics.Parties[0].Players[0]

		// Check for hard caps. Hard caps will have identical results because RNG is fixed.
		if modPlayerHigh.Dps.Avg == baselinePlayer.Dps.Avg && modPlayerHigh.Hps.Avg == baselinePlayer.Hps.Avg {
			continue
		}

		// For spell/melee hit, only use the direction facing away from the nearest soft/hard cap.
		if stat.EqualsStat(stats.SpellHit) {
			if baseStats.Stats[stat] >= spellHitCap {
				resultsLow[stat] = nil
			} else if baseStats.Stats[stat]+statModsHigh[stat] > spellHitCap {
				resultsHigh[stat] = nil
			}
		} else if stat.EqualsStat(stats.MeleeHit) {
			if baseStats.Stats[stat] >= melee2HHitCap {
				resultsLow[stat] = nil
			} else if baseStats.Stats[stat]+statModsHigh[stat] > melee2HHitCap {
				resultsHigh[stat] = nil
			}
		}

		calcWeightResults := func(baselineMetrics *proto.DistributionMetrics, modLowMetrics *proto.DistributionMetrics, modHighMetrics *proto.DistributionMetrics, weightResults *StatWeightValues) {
			var sample []float64
			if resultsLow != nil {
				for i := 0; i < int(simOptions.Iterations); i++ {
					sample = append(sample, (modLowMetrics.AllValues[i]-baselineMetrics.AllValues[i])/statModsLow[stat])
				}
			}
			if resultsHigh != nil {
				for i := 0; i < int(simOptions.Iterations); i++ {
					sample = append(sample, (modHighMetrics.AllValues[i]-baselineMetrics.AllValues[i])/statModsHigh[stat])
				}
			}

			weightResults.Weights.Stats[stat], weightResults.WeightsStdev.Stats[stat] = calcMeanAndStdev(sample)
		}

		calcWeightResults(baselinePlayer.Dps, modPlayerLow.Dps, modPlayerHigh.Dps, &result.Dps)
		calcWeightResults(baselinePlayer.Hps, modPlayerLow.Hps, modPlayerHigh.Hps, &result.Hps)
		calcWeightResults(baselinePlayer.Threat, modPlayerLow.Threat, modPlayerHigh.Threat, &result.Tps)
		calcWeightResults(baselinePlayer.Dtps, modPlayerLow.Dtps, modPlayerHigh.Dtps, &result.Dtps)
	}

	// Compute EP results.
	for statIdx := range statModsLow {
		stat := stats.Stat(statIdx)
		if statModsLow[stat] == 0 || statModsHigh[stat] == 0 {
			continue
		}

		calcEpResults := func(weightResults *StatWeightValues, refStat stats.Stat) {
			if weightResults.Weights.Stats[refStat] == 0 {
				return
			}
			weightResults.EpValues.Stats[stat] = weightResults.Weights.Stats[stat] / weightResults.Weights.Stats[refStat]
			weightResults.EpValuesStdev.Stats[stat] = weightResults.WeightsStdev.Stats[stat] / math.Abs(weightResults.Weights.Stats[refStat])
		}

		calcEpResults(&result.Dps, referenceStat)
		calcEpResults(&result.Hps, referenceStat)
		calcEpResults(&result.Tps, referenceStat)
		calcEpResults(&result.Dtps, DTPSReferenceStat)
	}

	return result
}
