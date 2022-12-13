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
		Hps:  swr.Hps.ToProto(),
		Tps:  swr.Tps.ToProto(),
		Dtps: swr.Dtps.ToProto(),
	}
}

func CalcStatWeight(swr *proto.StatWeightsRequest, statsToWeigh []stats.Stat, referenceStat stats.Stat, progress chan *proto.ProgressMetrics) StatWeightsResult {
	if swr.Player.BonusStats == nil {
		swr.Player.BonusStats = make([]float64, stats.Len)
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
	resultsLow := make([]*proto.RaidSimResult, stats.Len)
	resultsHigh := make([]*proto.RaidSimResult, stats.Len)

	var iterationsTotal int32
	var iterationsDone int32
	var simsTotal int32
	var simsCompleted int32

	doStat := func(stat stats.Stat, value float64, isLow bool) {
		defer waitGroup.Done()

		simRequest := googleProto.Clone(baseSimRequest).(*proto.RaidSimRequest)
		simRequest.Raid.Parties[0].Players[0].BonusStats[stat] += value

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
	statModsLow := stats.Stats{}
	statModsHigh := stats.Stats{}

	// Make sure reference stat is included.
	statModsLow[referenceStat] = defaultStatMod
	statModsHigh[referenceStat] = defaultStatMod

	for _, stat := range statsToWeigh {
		statMod := defaultStatMod
		if stat == stats.SpellHit {
			statMod = spellHitStatMod
			if baseStats[stat] < spellHitCap && baseStats[stat]+statMod > spellHitCap {
				// Check that newMod is atleast half of the previous mod, or we introduce a lot of deviation in the weight calc
				newMod := baseStats[stat] - spellHitCap
				if newMod > 0.5*statMod {
					statModsHigh[stat] = newMod
					statModsLow[stat] = -newMod
				} else {
					// Otherwise we go the opposite way of cap
					statModsHigh[stat] = -statMod
					statModsLow[stat] = -statMod
				}

				continue
			}
		} else if stat == stats.MeleeHit {
			statMod = meleeHitStatMod
			if baseStats[stat] < melee2HHitCap && baseStats[stat]+statMod > melee2HHitCap {
				// Check that newMod is atleast half of the previous mod, or we introduce a lot of deviation in the weight calc
				newMod := baseStats[stat] - melee2HHitCap
				if newMod > 0.5*statMod {
					statModsHigh[stat] = newMod
					statModsLow[stat] = -newMod
				} else {
					// Otherwise we go the opposite way of cap
					statModsHigh[stat] = -statMod
					statModsLow[stat] = -statMod
				}
				continue
			}
		} else if stat == stats.Expertise {
			// Expertise is non-linear, so adjust in increments that match the stepwise reduction.
			statMod = ExpertisePerQuarterPercentReduction
		} else if stat == stats.Armor {
			statMod = defaultStatMod * 10
		}
		statModsHigh[stat] = statMod
		statModsLow[stat] = -statMod
	}

	// Start all the threads.
	for stat := range statModsLow {
		if statModsLow[stat] == 0 {
			continue
		}
		waitGroup.Add(2)
		atomic.AddInt32(&iterationsTotal, swr.SimOptions.Iterations*2)
		atomic.AddInt32(&simsTotal, 2)

		go doStat(stats.Stat(stat), statModsLow[stat], true)
		go doStat(stats.Stat(stat), statModsHigh[stat], false)
	}

	// Wait for thread results.
	waitGroup.Wait()

	// Compute weight results.
	result := StatWeightsResult{}
	for _, stat := range statsToWeigh {
		if resultsLow[stat] == nil && resultsHigh[stat] == nil {
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
		if stat == stats.SpellHit {
			if baseStats[stat] >= spellHitCap {
				resultsLow[stat] = nil
			}
		} else if stat == stats.MeleeHit {
			if baseStats[stat] >= melee2HHitCap {
				resultsLow[stat] = nil
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

			sum := 0.0
			sumSq := 0.0
			for i := 0; i < len(sample); i++ {
				sum += sample[i]
				sumSq += sample[i] * sample[i]
			}
			iters := float64(len(sample))
			avg := sum / iters
			weightResults.Weights[stat] = avg
			weightResults.WeightsStdev[stat] = math.Abs(math.Sqrt((sumSq / iters)) - (avg * avg))
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
			if weightResults.Weights[refStat] == 0 {
				return
			}
			weightResults.EpValues[stat] = weightResults.Weights[stat] / weightResults.Weights[refStat]
			weightResults.EpValuesStdev[stat] = weightResults.WeightsStdev[stat] / math.Abs(weightResults.Weights[refStat])
		}

		calcEpResults(&result.Dps, referenceStat)
		calcEpResults(&result.Hps, referenceStat)
		calcEpResults(&result.Tps, referenceStat)
		calcEpResults(&result.Dtps, DTPSReferenceStat)
	}

	return result
}
