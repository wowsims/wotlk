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
		Hps:  swr.Dps.ToProto(),
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
	baselineDpsMetrics := baselineResult.RaidMetrics.Parties[0].Players[0].Dps
	baselineHpsMetrics := baselineResult.RaidMetrics.Parties[0].Players[0].Hps
	baselineTpsMetrics := baselineResult.RaidMetrics.Parties[0].Players[0].Threat
	baselineDtpsMetrics := baselineResult.RaidMetrics.Parties[0].Players[0].Dtps

	var waitGroup sync.WaitGroup

	// Do half the iterations with a positive, and half with a negative value for better accuracy.
	resultLow := StatWeightsResult{}
	resultHigh := StatWeightsResult{}
	dpsHistsLow := [stats.Len]map[int32]int32{}
	dpsHistsHigh := [stats.Len]map[int32]int32{}
	hpsHistsLow := [stats.Len]map[int32]int32{}
	hpsHistsHigh := [stats.Len]map[int32]int32{}
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
		dpsMetrics := simResult.RaidMetrics.Parties[0].Players[0].Dps
		hpsMetrics := simResult.RaidMetrics.Parties[0].Players[0].Hps
		tpsMetrics := simResult.RaidMetrics.Parties[0].Players[0].Threat
		dtpsMetrics := simResult.RaidMetrics.Parties[0].Players[0].Dtps
		dpsDiff := (dpsMetrics.Avg - baselineDpsMetrics.Avg) / value
		hpsDiff := (hpsMetrics.Avg - baselineHpsMetrics.Avg) / value
		tpsDiff := (tpsMetrics.Avg - baselineTpsMetrics.Avg) / value
		dtpsDiff := -(dtpsMetrics.Avg - baselineDtpsMetrics.Avg) / value

		if isLow {
			resultLow.Dps.Weights[stat] = dpsDiff
			resultLow.Hps.Weights[stat] = hpsDiff
			resultLow.Tps.Weights[stat] = tpsDiff
			resultLow.Dtps.Weights[stat] = dtpsDiff
			resultLow.Dps.WeightsStdev[stat] = dpsMetrics.Stdev / math.Abs(value)
			resultLow.Hps.WeightsStdev[stat] = hpsMetrics.Stdev / math.Abs(value)
			resultLow.Tps.WeightsStdev[stat] = tpsMetrics.Stdev / math.Abs(value)
			resultLow.Dtps.WeightsStdev[stat] = dtpsMetrics.Stdev / math.Abs(value)
			dpsHistsLow[stat] = dpsMetrics.Hist
			hpsHistsLow[stat] = hpsMetrics.Hist
			tpsHistsLow[stat] = tpsMetrics.Hist
			dtpsHistsLow[stat] = dtpsMetrics.Hist
		} else {
			resultHigh.Dps.Weights[stat] = dpsDiff
			resultHigh.Hps.Weights[stat] = tpsDiff
			resultHigh.Tps.Weights[stat] = tpsDiff
			resultHigh.Dtps.Weights[stat] = dtpsDiff
			resultHigh.Dps.WeightsStdev[stat] = dpsMetrics.Stdev / math.Abs(value)
			resultHigh.Hps.WeightsStdev[stat] = hpsMetrics.Stdev / math.Abs(value)
			resultHigh.Tps.WeightsStdev[stat] = tpsMetrics.Stdev / math.Abs(value)
			resultHigh.Dtps.WeightsStdev[stat] = dtpsMetrics.Stdev / math.Abs(value)
			dpsHistsHigh[stat] = dpsMetrics.Hist
			hpsHistsHigh[stat] = hpsMetrics.Hist
			tpsHistsHigh[stat] = tpsMetrics.Hist
			dtpsHistsHigh[stat] = dtpsMetrics.Hist
		}
	}

	// Melee hit cap is 8% in WoTLK
	melee2HHitCap := 8 * MeleeHitRatingPerHitChance
	// Spell hit cap is 17% in WoTLK
	spellHitCap := 17 * SpellHitRatingPerHitChance
	if swr.Debuffs != nil && (swr.Debuffs.Misery || swr.Debuffs.FaerieFire == proto.TristateEffect_TristateEffectImproved) {
		spellHitCap -= 3 * SpellHitRatingPerHitChance
	}

	const defaultStatMod = 10.0
	const meleeHitStatMod = defaultStatMod
	const spellHitStatMod = defaultStatMod
	//const meleeHitStatMod = MeleeHitRatingPerHitChance * 0.5
	//const spellHitStatMod = SpellHitRatingPerHitChance * 0.5
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
		}
		statModsHigh[stat] = statMod
		statModsLow[stat] = -statMod
	}

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

	waitGroup.Wait()

	for _, stat := range statsToWeigh {
		// Check for hard caps. Hard caps will have a high diff of exactly 0 because RNG is fixed.
		if resultHigh.Dps.Weights[stat] == 0 {
			statModsHigh[stat] = 0
			continue
		}

		// For spell/melee hit, only use the direction facing away from the nearest soft/hard cap.
		//
		if stat == stats.SpellHit {
			if baseStats[stat] >= spellHitCap {
				statModsLow[stat] = statModsHigh[stat]
				resultLow.Dps.Weights[stat] = resultHigh.Dps.Weights[stat]
				resultLow.Hps.Weights[stat] = resultHigh.Hps.Weights[stat]
				resultLow.Tps.Weights[stat] = resultHigh.Tps.Weights[stat]
				resultLow.Dtps.Weights[stat] = resultHigh.Dtps.Weights[stat]
			}
		} else if stat == stats.MeleeHit {
			if baseStats[stat] >= melee2HHitCap {
				statModsLow[stat] = statModsHigh[stat]
				resultLow.Dps.Weights[stat] = resultHigh.Dps.Weights[stat]
				resultLow.Hps.Weights[stat] = resultHigh.Hps.Weights[stat]
				resultLow.Tps.Weights[stat] = resultHigh.Tps.Weights[stat]
				resultLow.Dtps.Weights[stat] = resultHigh.Dtps.Weights[stat]
			}
		}
	}

	result := StatWeightsResult{}
	for statIdx := range statModsLow {
		stat := stats.Stat(statIdx)
		if statModsLow[stat] == 0 || statModsHigh[stat] == 0 {
			continue
		}

		result.Dps.Weights[stat] = (resultLow.Dps.Weights[stat] + resultHigh.Dps.Weights[stat]) / 2
		result.Hps.Weights[stat] = (resultLow.Hps.Weights[stat] + resultHigh.Hps.Weights[stat]) / 2
		result.Tps.Weights[stat] = (resultLow.Tps.Weights[stat] + resultHigh.Tps.Weights[stat]) / 2
		result.Dtps.Weights[stat] = (resultLow.Dtps.Weights[stat] + resultHigh.Dtps.Weights[stat]) / 2

		result.Dps.WeightsStdev[stat] = (resultLow.Dps.WeightsStdev[stat] + resultHigh.Dps.WeightsStdev[stat]) / 2
		result.Hps.WeightsStdev[stat] = (resultLow.Hps.WeightsStdev[stat] + resultHigh.Hps.WeightsStdev[stat]) / 2
		result.Tps.WeightsStdev[stat] = (resultLow.Tps.WeightsStdev[stat] + resultHigh.Tps.WeightsStdev[stat]) / 2
		result.Dtps.WeightsStdev[stat] = (resultLow.Dtps.WeightsStdev[stat] + resultHigh.Dtps.WeightsStdev[stat]) / 2
	}

	for statIdx := range statModsLow {
		stat := stats.Stat(statIdx)
		if statModsLow[stat] == 0 || statModsHigh[stat] == 0 {
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
