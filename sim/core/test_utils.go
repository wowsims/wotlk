package core

import (
	"log"
	"testing"

	"github.com/wowsims/tbc/sim/core/proto"
	"github.com/wowsims/tbc/sim/core/stats"
	googleProto "google.golang.org/protobuf/proto"
)

var DefaultSimTestOptions = &proto.SimOptions{
	Iterations: 50,
	IsTest:     true,
	Debug:      false,
	RandomSeed: 101,
}
var StatWeightsDefaultSimTestOptions = &proto.SimOptions{
	Iterations: 1000,
	IsTest:     true,
	Debug:      false,
	RandomSeed: 101,
}
var AverageDefaultSimTestOptions = &proto.SimOptions{
	Iterations: 5000,
	IsTest:     true,
	Debug:      false,
	RandomSeed: 101,
}

const ShortDuration = 60
const LongDuration = 300

var DefaultTargetProto = proto.Target{
	Level: 73,
	Stats: stats.Stats{
		stats.Armor:       7684,
		stats.AttackPower: 320,
		stats.BlockValue:  54,
	}.ToFloatArray(),
	MobType: proto.MobType_MobTypeDemon,

	SwingSpeed:    2,
	MinBaseDamage: 4192.05,
	CanCrush:      true,
	ParryHaste:    true,
}

func NewDefaultTarget() *proto.Target {
	var target = &proto.Target{}
	*target = DefaultTargetProto
	return target
}

func MakeDefaultEncounterCombos(debuffs *proto.Debuffs) []EncounterCombo {
	var DefaultTarget = NewDefaultTarget()

	multipleTargets := []*proto.Target{}
	for i := 0; i < 20; i++ {
		multipleTargets = append(multipleTargets, DefaultTarget)
	}

	return []EncounterCombo{
		EncounterCombo{
			Label: "ShortSingleTarget",
			Encounter: &proto.Encounter{
				Duration:          ShortDuration,
				ExecuteProportion: 0.2,
				Targets: []*proto.Target{
					DefaultTarget,
				},
			},
		},
		EncounterCombo{
			Label: "LongSingleTarget",
			Encounter: &proto.Encounter{
				Duration:          LongDuration,
				ExecuteProportion: 0.2,
				Targets: []*proto.Target{
					DefaultTarget,
				},
			},
		},
		EncounterCombo{
			Label: "LongMultiTarget",
			Encounter: &proto.Encounter{
				Duration:          LongDuration,
				ExecuteProportion: 0.2,
				Targets:           multipleTargets,
			},
		},
	}
}

func MakeSingleTargetEncounter(variation float64) *proto.Encounter {
	return &proto.Encounter{
		Duration:          LongDuration,
		DurationVariation: variation,
		ExecuteProportion: 0.2,
		Targets: []*proto.Target{
			NewDefaultTarget(),
		},
	}
}

func CharacterStatsTest(label string, t *testing.T, raid *proto.Raid, expectedStats stats.Stats) {
	csr := &proto.ComputeStatsRequest{
		Raid: raid,
	}

	result := ComputeStats(csr)
	finalStats := stats.FromFloatArray(result.RaidStats.Parties[0].Players[0].FinalStats)

	const tolerance = 0.5
	if !finalStats.EqualsWithTolerance(expectedStats, tolerance) {
		t.Fatalf("%s failed: CharacterStats() = %v, expected %v", label, finalStats, expectedStats)
	}
}

func StatWeightsTest(label string, t *testing.T, _swr *proto.StatWeightsRequest, expectedStatWeights stats.Stats) {
	// Make a copy so we can safely change fields.
	swr := googleProto.Clone(_swr).(*proto.StatWeightsRequest)

	swr.Encounter.Duration = LongDuration
	swr.SimOptions.Iterations = 5000

	result := StatWeights(swr)
	resultWeights := stats.FromFloatArray(result.Dps.Weights)

	const tolerance = 0.05
	if !resultWeights.EqualsWithTolerance(expectedStatWeights, tolerance) {
		t.Fatalf("%s failed: CalcStatWeight() = %v, expected %v", label, resultWeights, expectedStatWeights)
	}
}

func RaidSimTest(label string, t *testing.T, rsr *proto.RaidSimRequest, expectedDps float64) {
	result := RunRaidSim(rsr)
	if result.ErrorResult != "" {
		t.Fatalf("Sim failed with error: %s", result.ErrorResult)
	}
	tolerance := 0.5
	if result.RaidMetrics.Dps.Avg < expectedDps-tolerance || result.RaidMetrics.Dps.Avg > expectedDps+tolerance {
		// Automatically print output if we had debugging enabled.
		if rsr.SimOptions.Debug {
			log.Printf("LOGS:\n%s\n", result.Logs)
		}
		t.Fatalf("%s failed: expected %0f dps from sim but was %0f", label, expectedDps, result.RaidMetrics.Dps.Avg)
	}
}

func RaidBenchmark(b *testing.B, rsr *proto.RaidSimRequest) {
	rsr.Encounter.Duration = LongDuration
	rsr.SimOptions.Iterations = 1

	// Set to false because IsTest adds a lot of computation.
	rsr.SimOptions.IsTest = false

	for i := 0; i < b.N; i++ {
		RunRaidSim(rsr)
	}
}
