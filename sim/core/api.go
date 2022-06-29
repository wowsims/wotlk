// Proto-based function interface for the simulator
package core

import (
	"github.com/wowsims/tbc/sim/core/items"
	"github.com/wowsims/tbc/sim/core/proto"
	"github.com/wowsims/tbc/sim/core/stats"
)

/**
 * Returns all items, enchants, and gems recognized by the sim.
 */
func GetGearList(request *proto.GearListRequest) *proto.GearListResult {
	result := &proto.GearListResult{
		Encounters: presetEncounters[:],
	}

	for i := range items.Items {
		item := items.Items[i]
		result.Items = append(result.Items, item.ToProto())
	}
	for i := range items.Gems {
		gem := items.Gems[i]
		result.Gems = append(result.Gems, gem.ToProto())
	}
	for i := range items.Enchants {
		enchant := items.Enchants[i]
		result.Enchants = append(result.Enchants, enchant.ToProto())
	}

	return result
}

/**
 * Returns character stats taking into account gear / buffs / consumes / etc
 */
func ComputeStats(csr *proto.ComputeStatsRequest) *proto.ComputeStatsResult {
	_, raidStats := NewEnvironment(*csr.Raid, proto.Encounter{})

	return &proto.ComputeStatsResult{
		RaidStats: raidStats,
	}
}

/**
 * Returns stat weights and EP values, with standard deviations, for all stats.
 */
func StatWeights(request *proto.StatWeightsRequest) *proto.StatWeightsResult {
	statsToWeigh := stats.ProtoArrayToStatsList(request.StatsToWeigh)

	result := CalcStatWeight(*request, statsToWeigh, stats.Stat(request.EpReferenceStat), nil)

	return result.ToProto()
}

func StatWeightsAsync(request *proto.StatWeightsRequest, progress chan *proto.ProgressMetrics) {
	statsToWeigh := stats.ProtoArrayToStatsList(request.StatsToWeigh)
	go func() {
		result := CalcStatWeight(*request, statsToWeigh, stats.Stat(request.EpReferenceStat), progress)
		progress <- &proto.ProgressMetrics{
			FinalWeightResult: result.ToProto(),
		}
	}()
}

/**
 * Runs multiple iterations of the sim with a full raid.
 */
func RunRaidSim(request *proto.RaidSimRequest) *proto.RaidSimResult {
	return RunSim(*request, nil)
}

func RunRaidSimAsync(request *proto.RaidSimRequest, progress chan *proto.ProgressMetrics) {
	go RunSim(*request, progress)
}
