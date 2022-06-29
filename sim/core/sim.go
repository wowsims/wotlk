package core

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"strings"
	"time"

	"github.com/wowsims/tbc/sim/core/proto"
)

type Simulation struct {
	*Environment

	Options proto.SimOptions

	rand Rand

	// Used for testing only, see RandomFloat().
	isTest    bool
	testRands map[string]Rand

	// Current Simulation State
	pendingActions []*PendingAction
	CurrentTime    time.Duration // duration that has elapsed in the sim since starting
	Duration       time.Duration // Duration of current iteration

	ProgressReport func(*proto.ProgressMetrics)

	Log  func(string, ...interface{})
	logs []string

	executePhase          bool
	executePhaseCallbacks []func(*Simulation)
}

func RunSim(rsr proto.RaidSimRequest, progress chan *proto.ProgressMetrics) (result *proto.RaidSimResult) {
	return runSim(rsr, progress, false)
}

func runSim(rsr proto.RaidSimRequest, progress chan *proto.ProgressMetrics, skipPresim bool) (result *proto.RaidSimResult) {
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

	sim := NewSim(rsr)

	if !skipPresim {
		if progress != nil {
			progress <- &proto.ProgressMetrics{
				TotalIterations: sim.Options.Iterations,
				PresimRunning:   true,
			}
			runtime.Gosched() // allow time for message to make it back out.
		}
		presimResult := sim.runPresims(rsr)
		if progress != nil {
			progress <- &proto.ProgressMetrics{
				TotalIterations: sim.Options.Iterations,
				PresimRunning:   false,
			}
			sim.ProgressReport = func(progMetric *proto.ProgressMetrics) {
				progress <- progMetric
			}
			runtime.Gosched() // allow time for message to make it back out.
		}

		// Use pre-sim as estimate for length of fight (when using health fight)
		if sim.Encounter.EndFightAtHealth > 0 && presimResult != nil {
			sim.BaseDuration = time.Duration(presimResult.AvgIterationDuration) * time.Second
			sim.Duration = time.Duration(presimResult.AvgIterationDuration) * time.Second
			sim.Encounter.DurationIsEstimate = false // we now have a pretty good value for duration
		}
	}

	// using a variable here allows us to mutate it in the deferred recover, sending out error info
	result = sim.run()

	return result
}

func NewSim(rsr proto.RaidSimRequest) *Simulation {
	simOptions := *rsr.SimOptions
	rseed := simOptions.RandomSeed
	if rseed == 0 {
		rseed = time.Now().UnixNano()
	}

	env, _ := NewEnvironment(*rsr.Raid, *rsr.Encounter)
	return &Simulation{
		Environment: env,
		Options:     simOptions,

		rand: NewSplitMix(uint64(rseed)),

		isTest:    simOptions.IsTest,
		testRands: make(map[string]Rand),
	}
}

// Returns a random float.
//
// In tests, although we can set the initial seed, test results are still very
// sensitive to the exact order of RandomFloat() calls. To mitigate this, when
// testing we use a separate rand object for each RandomFloat callsite,
// distinguished by the label string.
func (sim *Simulation) RandomFloat(label string) float64 {
	if !sim.isTest {
		return sim.rand.NextFloat64()
	}

	labelRand, isPresent := sim.testRands[label]
	if !isPresent {
		labelRand = NewSplitMix(uint64(hash(label)))
		sim.testRands[label] = labelRand
	}
	v := labelRand.NextFloat64()
	// if sim.Log != nil {
	// 	sim.Log("FLOAT64 '%s': %0.5f", label, v)
	// }
	return v
}

func (sim *Simulation) Reset() {
	sim.reset()
}

// Reset will set sim back and erase all current state.
// This is automatically called before every 'Run'.
func (sim *Simulation) reset() {
	if sim.Log != nil {
		sim.Log("SIM RESET")
		sim.Log("----------------------")
	}

	// Reset primary targets damage taken for tracking health fights.
	sim.Encounter.DamageTaken = 0

	if sim.Encounter.DurationIsEstimate && sim.CurrentTime != 0 {
		sim.BaseDuration = sim.CurrentTime
		sim.Encounter.DurationIsEstimate = false
	}
	sim.Duration = sim.BaseDuration
	if sim.DurationVariation != 0 {
		variation := sim.DurationVariation * 2
		sim.Duration += time.Duration((sim.RandomFloat("sim duration") * float64(variation))) - sim.DurationVariation
	}

	sim.CurrentTime = 0.0

	sim.pendingActions = make([]*PendingAction, 0, 64)

	sim.executePhase = false
	sim.executePhaseCallbacks = []func(*Simulation){}

	// Targets need to be reset before the raid, so that players can check for
	// the presence of permanent target auras in their Reset handlers.
	for _, target := range sim.Encounter.Targets {
		target.Reset(sim)
	}

	sim.Raid.reset(sim)

	sim.initManaTickAction()
}

// Run runs the simulation for the configured number of iterations, and
// collects all the metrics together.
func (sim *Simulation) run() *proto.RaidSimResult {
	logsBuffer := &strings.Builder{}
	if sim.Options.Debug || sim.Options.DebugFirstIteration {
		sim.Log = func(message string, vals ...interface{}) {
			logsBuffer.WriteString(fmt.Sprintf("[%0.2f] "+message+"\n", append([]interface{}{sim.CurrentTime.Seconds()}, vals...)...))
		}
	}

	// Uncomment this to print logs directly to console.
	// sim.Options.Debug = true
	// sim.Log = func(message string, vals ...interface{}) {
	// 	fmt.Printf(fmt.Sprintf("[%0.1f] "+message+"\n", append([]interface{}{sim.CurrentTime.Seconds()}, vals...)...))
	// }

	for _, target := range sim.Encounter.Targets {
		target.init(sim)
	}

	for _, party := range sim.Raid.Parties {
		for _, player := range party.Players {
			character := player.GetCharacter()
			character.init(sim, player)

			for _, petAgent := range character.Pets {
				petAgent.GetCharacter().init(sim, petAgent)
			}
		}
	}

	sim.runOnce()
	firstIterationDuration := sim.CurrentTime
	totalDuration := firstIterationDuration

	if !sim.Options.Debug {
		sim.Log = nil
	}

	var st time.Time
	for i := int32(1); i < sim.Options.Iterations; i++ {
		// fmt.Printf("Iteration: %d\n", i)
		if sim.ProgressReport != nil && time.Since(st) > time.Millisecond*100 {
			metrics := sim.Raid.GetMetrics(i + 1)
			sim.ProgressReport(&proto.ProgressMetrics{TotalIterations: sim.Options.Iterations, CompletedIterations: i + 1, Dps: metrics.Dps.Avg})
			runtime.Gosched() // ensure that reporting threads are given time to report, mostly only important in wasm (only 1 thread)
			st = time.Now()
		}
		sim.runOnce()
		totalDuration += sim.CurrentTime
	}
	result := &proto.RaidSimResult{
		RaidMetrics:      sim.Raid.GetMetrics(sim.Options.Iterations),
		EncounterMetrics: sim.Encounter.GetMetricsProto(sim.Options.Iterations),

		Logs:                   logsBuffer.String(),
		FirstIterationDuration: firstIterationDuration.Seconds(),
		AvgIterationDuration:   totalDuration.Seconds() / float64(sim.Options.Iterations),
	}

	// Final progress report
	if sim.ProgressReport != nil {
		sim.ProgressReport(&proto.ProgressMetrics{TotalIterations: sim.Options.Iterations, CompletedIterations: sim.Options.Iterations, Dps: result.RaidMetrics.Dps.Avg, FinalRaidResult: result})
	}

	return result
}

// RunOnce is the main event loop. It will run the simulation for number of seconds.
func (sim *Simulation) runOnce() {
	sim.reset()

	for {
		last := len(sim.pendingActions) - 1
		pa := sim.pendingActions[last]
		sim.pendingActions = sim.pendingActions[:last]
		if pa.cancelled {
			continue
		}

		// Use duration as an end check if not using health.
		if sim.Encounter.EndFightAtHealth == 0 {
			if pa.NextActionAt > sim.Duration {
				break
			}
		} else if sim.Encounter.EndFightAtHealth < sim.Encounter.DamageTaken {
			break
		}

		if pa.NextActionAt > sim.CurrentTime {
			sim.advance(pa.NextActionAt - sim.CurrentTime)
		}

		pa.OnAction(sim)
	}

	for _, pa := range sim.pendingActions {
		if pa.CleanUp != nil {
			pa.CleanUp(sim)
		}
	}

	sim.Raid.doneIteration(sim)
	sim.Encounter.doneIteration(sim)

	for _, unit := range sim.Raid.AllUnits {
		unit.Metrics.doneIteration(sim.CurrentTime.Seconds())
	}
	for _, target := range sim.Encounter.Targets {
		target.Metrics.doneIteration(sim.CurrentTime.Seconds())
	}
}

func (sim *Simulation) AddPendingAction(pa *PendingAction) {
	for index, v := range sim.pendingActions {
		if v.NextActionAt < pa.NextActionAt || (v.NextActionAt == pa.NextActionAt && v.Priority >= pa.Priority) {
			sim.pendingActions = append(sim.pendingActions, pa)
			copy(sim.pendingActions[index+1:], sim.pendingActions[index:])
			sim.pendingActions[index] = pa
			return
		}
	}
	sim.pendingActions = append(sim.pendingActions, pa)
}

// Advance moves time forward counting down auras, CDs, mana regen, etc
func (sim *Simulation) advance(elapsedTime time.Duration) {
	sim.CurrentTime += elapsedTime

	if !sim.executePhase {
		if (sim.Encounter.EndFightAtHealth == 0 && sim.CurrentTime >= sim.Encounter.executePhaseBegins) ||
			(sim.Encounter.EndFightAtHealth > 0 && sim.GetRemainingDurationPercent() <= 0.2) {
			sim.executePhase = true
			for _, callback := range sim.executePhaseCallbacks {
				callback(sim)
			}
		}
	}

	for _, party := range sim.Raid.Parties {
		for _, agent := range party.Players {
			agent.GetCharacter().advance(sim, elapsedTime)
		}
	}

	for _, target := range sim.Encounter.Targets {
		target.Advance(sim, elapsedTime)
	}
}

func (sim *Simulation) RegisterExecutePhaseCallback(callback func(*Simulation)) {
	sim.executePhaseCallbacks = append(sim.executePhaseCallbacks, callback)
}
func (sim *Simulation) IsExecutePhase() bool {
	return sim.executePhase
}

func (sim *Simulation) GetRemainingDuration() time.Duration {
	if sim.Encounter.EndFightAtHealth > 0 {
		if !sim.Encounter.DurationIsEstimate || sim.CurrentTime < time.Second*5 {
			return sim.Duration - sim.CurrentTime
		}

		// Estimate time remaining via avg dps
		dps := sim.Encounter.DamageTaken / sim.CurrentTime.Seconds()
		dur := time.Duration((sim.Encounter.EndFightAtHealth-sim.Encounter.DamageTaken)/dps) * time.Second
		return dur
	}
	return sim.Duration - sim.CurrentTime
}

// Returns the percentage of time remaining in the current iteration, as a value from 0-1.
func (sim *Simulation) GetRemainingDurationPercent() float64 {
	if sim.Encounter.EndFightAtHealth > 0 {
		return 1.0 - sim.Encounter.DamageTaken/sim.Encounter.EndFightAtHealth
	}
	return float64(sim.Duration-sim.CurrentTime) / float64(sim.Duration)
}
