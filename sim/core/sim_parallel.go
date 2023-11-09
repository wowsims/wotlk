package core

import (
	"runtime"
	"runtime/debug"
	"sync"

	"github.com/wowsims/wotlk/sim/core/proto"
	googleProto "google.golang.org/protobuf/proto"
)

func runSimParallel(rsr *proto.RaidSimRequest, progress chan *proto.ProgressMetrics) (result *proto.RaidSimResult) {
	defer func() {
		if err := recover(); err != nil {
			errStr := ""
			switch errt := err.(type) {
			case string:
				errStr = errt
			case error:
				errStr = errt.Error()
			}

			errStr += "\nStack Trace:\n" + string(debug.Stack())
			result = &proto.RaidSimResult{
				ErrorResult: errStr,
			}
			if progress != nil {
				progress <- &proto.ProgressMetrics{
					FinalRaidResult: result,
				}
			}
		}
		if progress != nil {
			close(progress)
		}
	}()

	numParallel := 4

	totalIterations := rsr.SimOptions.Iterations
	rsr.SimOptions.Iterations = totalIterations / int32(numParallel)
	sims := []*Simulation{}

	for i := 0; i < numParallel; i++ {
		sims[i] = NewSim(googleProto.Clone(rsr).(*proto.RaidSimRequest))
	}

	progress <- &proto.ProgressMetrics{
		TotalIterations: totalIterations,
		PresimRunning:   true,
	}
	runtime.Gosched() // allow time for message to make it back out.

	presimResult := sims[0].runPresims(rsr)
	if presimResult != nil && presimResult.ErrorResult != "" {
		if progress != nil {
			progress <- &proto.ProgressMetrics{
				TotalIterations: totalIterations,
				FinalRaidResult: presimResult, // send error from presim out.
			}
		}
		return presimResult
	}
	progress <- &proto.ProgressMetrics{
		TotalIterations: totalIterations,
		PresimRunning:   false,
	}

	// TODO: setup intermediate progress function for each sim
	for i := 0; i < numParallel; i++ {
		sims[i].ProgressReport = func(progMetric *proto.ProgressMetrics) {
			progress <- progMetric
		}
	}

	runtime.Gosched() // allow time for message to make it back out.

	// Use pre-sim as estimate for length of fight (when using health fight)
	// TODO: setup each sim
	// if sim1.Encounter.EndFightAtHealth > 0 && presimResult != nil {
	// 	sim1.BaseDuration = time.Duration(presimResult.AvgIterationDuration) * time.Second
	// 	sim1.Duration = time.Duration(presimResult.AvgIterationDuration) * time.Second
	// 	sim1.Encounter.DurationIsEstimate = false // we now have a pretty good value for duration
	// }

	wg := &sync.WaitGroup{}
	wg.Add(numParallel)

	results := []*proto.RaidSimResult{}

	// var logsBuffer
	runPartial := func(sim *Simulation, i int) {
		results[i] = sim.run()
		wg.Done()
	}
	for i := 0; i < numParallel; i++ {
		go runPartial(sims[i], i)
	}

	wg.Wait()
	result = &proto.RaidSimResult{}

	for i := 0; i < numParallel; i++ {
		result.AvgIterationDuration += results[i].AvgIterationDuration

		if result.FirstIterationDuration == 0 {
			result.FirstIterationDuration = results[i].FirstIterationDuration
		}

		if result.Logs == "" && results[i].Logs != "" {
			result.Logs = results[i].Logs
		}

		if result.ErrorResult == "" && results[i].ErrorResult != "" {
			result.ErrorResult = results[i].ErrorResult
		}

		if result.RaidMetrics == nil {
			result.RaidMetrics = results[i].RaidMetrics
		} else {
			addDistrib(result.RaidMetrics.Dps, results[i].RaidMetrics.Dps)
			addDistrib(result.RaidMetrics.Hps, results[i].RaidMetrics.Hps)
			addParties(result.RaidMetrics.Parties, results[i].RaidMetrics.Parties)
		}

		if result.EncounterMetrics == nil {
			result.EncounterMetrics = results[i].EncounterMetrics
		} else {
			addUnits(result.EncounterMetrics.Targets, results[i].EncounterMetrics.Targets)
		}

	}

	return result
}

func addParties(a, b []*proto.PartyMetrics) {
	for i, pa := range a {
		pb := b[i]
		addDistrib(pa.Dps, pb.Dps)
		addDistrib(pa.Hps, pb.Hps)
		addUnits(pa.Players, pb.Players)
	}
}

func addUnits(a, b []*proto.UnitMetrics) {
	for i, pa := range a {
		pb := b[i]
		addDistrib(pa.Dps, pb.Dps)
		addDistrib(pa.Dpasp, pb.Dpasp)
		addDistrib(pa.Threat, pb.Threat)
		addDistrib(pa.Dtps, pb.Dtps)
		addDistrib(pa.Tmi, pb.Tmi)
		addDistrib(pa.Hps, pb.Hps)
		addDistrib(pa.Tto, pb.Tto)

		pa.SecondsOomAvg += pb.SecondsOomAvg
		pa.ChanceOfDeath += pb.ChanceOfDeath

		// addDistrib(pa.Actions, pb.Actions)
		// addDistrib(pa.Auras, pb.Auras)
		// addDistrib(pa.Resources, pb.Resources)
		// addDistrib(pa.Pets, pb.Pets)
	}
}

func addDistrib(a *proto.DistributionMetrics, b *proto.DistributionMetrics) {
	if b == nil {
		return
	}

	a.Avg += b.Avg
	a.Stdev += b.Stdev

	a.Max = max(b.Max, a.Max)
	a.MaxSeed = max(b.MaxSeed, a.MaxSeed)
	a.Min = min(a.Min, b.Min)
	a.MinSeed = min(b.MinSeed, a.MinSeed)

	for k, v := range b.Hist {
		a.Hist[k] += v
	}

	// used only for stat weights
	// a.AllValues += b.AllValues
}
