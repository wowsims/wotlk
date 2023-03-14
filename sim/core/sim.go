package core

import (
	"fmt"
	"math/rand"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/wowsims/wotlk/sim/core/proto"
)

type Simulation struct {
	*Environment

	Options *proto.SimOptions

	rand  Rand
	rseed int64

	// Used for testing only, see RandomFloat().
	isTest    bool
	testRands map[string]Rand

	// Current Simulation State
	pendingActions []*PendingAction
	CurrentTime    time.Duration // duration that has elapsed in the sim since starting
	Duration       time.Duration // Duration of current iteration

	ProgressReport func(*proto.ProgressMetrics)

	Log func(string, ...interface{})

	executePhase20Begins  time.Duration
	executePhase25Begins  time.Duration
	executePhase35Begins  time.Duration
	executePhase20        bool
	executePhase25        bool
	executePhase35        bool
	executePhaseCallbacks []func(*Simulation, int) // 2nd parameter is 35 for 35%, 25 for 25% and 20 for 20%
}

func RunSim(rsr *proto.RaidSimRequest, progress chan *proto.ProgressMetrics) (result *proto.RaidSimResult) {
	return runSim(rsr, progress, false)
}

func runSim(rsr *proto.RaidSimRequest, progress chan *proto.ProgressMetrics, skipPresim bool) (result *proto.RaidSimResult) {
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
		if presimResult != nil && presimResult.ErrorResult != "" {
			if progress != nil {
				progress <- &proto.ProgressMetrics{
					TotalIterations: sim.Options.Iterations,
					FinalRaidResult: presimResult,
				}
			}
			return presimResult
		}
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

func NewSim(rsr *proto.RaidSimRequest) *Simulation {
	simOptions := rsr.SimOptions
	rseed := simOptions.RandomSeed
	if rseed == 0 {
		rseed = time.Now().UnixNano()
	}

	env, _ := NewEnvironment(rsr.Raid, rsr.Encounter)
	return &Simulation{
		Environment: env,
		Options:     simOptions,

		rand:  NewSplitMix(uint64(rseed)),
		rseed: rseed,

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
	return sim.labelRand(label).NextFloat64()
}

func (sim *Simulation) labelRand(label string) Rand {
	if !sim.isTest {
		return sim.rand
	}

	labelRand, isPresent := sim.testRands[label]
	if !isPresent {
		// Add rseed to the label to we still have run-run variance for stat weights.
		labelRand = NewSplitMix(uint64(makeTestRandSeed(sim.rseed, label)))
		sim.testRands[label] = labelRand
	}
	return labelRand
}

func (sim *Simulation) reseedRands(i int64) {
	rseed := sim.Options.RandomSeed + i
	sim.rand.Seed(rseed)

	if sim.isTest {
		for label, rand := range sim.testRands {
			rand.Seed(makeTestRandSeed(rseed, label))
		}
	}
}

func makeTestRandSeed(rseed int64, label string) int64 {
	return int64(hash(label + strconv.FormatInt(rseed, 16)))
}

func (sim *Simulation) RandomExpFloat(label string) float64 {
	return rand.New(sim.labelRand(label)).ExpFloat64()
}

// Shorthand for commonly-used RNG behavior.
// Returns a random number between min and max.
func (sim *Simulation) Roll(min float64, max float64) float64 {
	return sim.RollWithLabel(min, max, "Damage Roll")
}
func (sim *Simulation) RollWithLabel(min float64, max float64, label string) float64 {
	return min + (max-min)*sim.RandomFloat(label)
}

func (sim *Simulation) Proc(p float64, label string) bool {
	switch {
	case p >= 1:
		return true
	case p <= 0:
		return false
	default:
		return sim.RandomFloat(label) < p
	}
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

	if sim.Encounter.DurationIsEstimate && sim.CurrentTime != 0 {
		sim.BaseDuration = sim.CurrentTime
		sim.Encounter.DurationIsEstimate = false
	}
	sim.Duration = sim.BaseDuration
	if sim.DurationVariation != 0 {
		variation := sim.DurationVariation * 2
		sim.Duration += time.Duration(sim.RandomFloat("sim duration")*float64(variation)) - sim.DurationVariation
	}
	sim.executePhase20Begins = time.Duration(float64(sim.Duration) * (1.0 - sim.Encounter.ExecuteProportion_20))
	sim.executePhase25Begins = time.Duration(float64(sim.Duration) * (1.0 - sim.Encounter.ExecuteProportion_25))
	sim.executePhase35Begins = time.Duration(float64(sim.Duration) * (1.0 - sim.Encounter.ExecuteProportion_35))

	sim.pendingActions = make([]*PendingAction, 0, 64)

	sim.executePhase20 = false
	sim.executePhase25 = false
	sim.executePhase35 = false
	sim.executePhaseCallbacks = []func(*Simulation, int){}

	sim.CurrentTime = 0

	sim.Environment.reset(sim)

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
	firstIterationDuration := sim.Duration
	if sim.Encounter.EndFightAtHealth != 0 {
		firstIterationDuration = sim.CurrentTime
	}
	totalDuration := firstIterationDuration

	if !sim.Options.Debug {
		sim.Log = nil
	}

	var st time.Time
	for i := int32(1); i < sim.Options.Iterations; i++ {
		// fmt.Printf("Iteration: %d\n", i)
		if sim.ProgressReport != nil && time.Since(st) > time.Millisecond*100 {
			metrics := sim.Raid.GetMetrics()
			sim.ProgressReport(&proto.ProgressMetrics{TotalIterations: sim.Options.Iterations, CompletedIterations: i, Dps: metrics.Dps.Avg, Hps: metrics.Hps.Avg})
			runtime.Gosched() // ensure that reporting threads are given time to report, mostly only important in wasm (only 1 thread)
			st = time.Now()
		}

		// Before each iteration, reset state to seed+iterations
		sim.reseedRands(int64(i))

		sim.runOnce()
		iterDuration := sim.Duration
		if sim.Encounter.EndFightAtHealth != 0 {
			iterDuration = sim.CurrentTime
		}
		totalDuration += iterDuration
	}
	result := &proto.RaidSimResult{
		RaidMetrics:      sim.Raid.GetMetrics(),
		EncounterMetrics: sim.Encounter.GetMetricsProto(),

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

func (sim *Simulation) runPendingActions(max time.Duration) {
	for {
		if len(sim.pendingActions) == 0 {
			return
		}

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
			if pa.NextActionAt < max {
				sim.advance(pa.NextActionAt - sim.CurrentTime)
			} else {
				sim.pendingActions = append(sim.pendingActions, pa)
				break
			}
		}
		pa.consumed = true

		if pa.cancelled {
			continue // pa was cancelled during the advance.
		}
		pa.OnAction(sim)
	}
}

// RunOnce is the main event loop. It will run the simulation for number of seconds.
func (sim *Simulation) runOnce() {
	sim.reset()

	if len(sim.Environment.prepullActions) > 0 {
		sim.CurrentTime = sim.Environment.prepullActions[0].DoAt

		for _, prepullAction := range sim.Environment.prepullActions {
			if prepullAction.DoAt > sim.CurrentTime {
				sim.runPendingActions(prepullAction.DoAt)
				sim.advance(prepullAction.DoAt - sim.CurrentTime)
			}
			prepullAction.Action(sim)
		}

		if sim.CurrentTime < 0 {
			sim.runPendingActions(0)
			sim.advance(0 - sim.CurrentTime)
		}
	}

	for _, unit := range sim.Environment.AllUnits {
		unit.startPull(sim)
	}

	sim.runPendingActions(NeverExpires)

	// The last event loop will leave CurrentTime at some value close to but not
	// quite at the Duration. Explicitly set this so that accesses to CurrentTime
	// during the doneIteration phase will return the Duration value, which is
	// intuitive.
	sim.CurrentTime = sim.Duration

	for _, pa := range sim.pendingActions {
		if pa.CleanUp != nil {
			pa.CleanUp(sim)
		}
	}

	sim.Raid.doneIteration(sim)
	sim.Encounter.doneIteration(sim)

	for _, unit := range sim.Raid.AllUnits {
		unit.Metrics.doneIteration(unit, sim)
	}
	for _, target := range sim.Encounter.TargetUnits {
		target.Metrics.doneIteration(target, sim)
	}
}

func (sim *Simulation) AddPendingAction(pa *PendingAction) {
	//if pa.NextActionAt < sim.CurrentTime {
	//	panic(fmt.Sprintf("Cant add action in the past: %s", pa.NextActionAt))
	//}
	pa.consumed = false
	for index, v := range sim.pendingActions {
		if v.NextActionAt < pa.NextActionAt || (v.NextActionAt == pa.NextActionAt && v.Priority >= pa.Priority) {
			//if sim.Log != nil {
			//	sim.Log("Adding action at index %d for time %s", index - len(sim.pendingActions), pa.NextActionAt)
			//	for i := index; i < len(sim.pendingActions); i++ {
			//		sim.Log("Upcoming action at %s", sim.pendingActions[i].NextActionAt)
			//	}
			//}
			sim.pendingActions = append(sim.pendingActions, pa)
			copy(sim.pendingActions[index+1:], sim.pendingActions[index:])
			sim.pendingActions[index] = pa
			return
		}
	}
	//if sim.Log != nil {
	//	sim.Log("Adding action at end for time %s", pa.NextActionAt)
	//}
	sim.pendingActions = append(sim.pendingActions, pa)
}

// Advance moves time forward counting down auras, CDs, mana regen, etc
func (sim *Simulation) advance(elapsedTime time.Duration) {
	sim.CurrentTime += elapsedTime

	if !sim.executePhase35 {
		if (sim.Encounter.EndFightAtHealth == 0 && sim.CurrentTime >= sim.executePhase35Begins) ||
			(sim.Encounter.EndFightAtHealth > 0 && sim.GetRemainingDurationPercent() <= 0.35) {
			sim.executePhase35 = true
			for _, callback := range sim.executePhaseCallbacks {
				callback(sim, 35)
			}
		}
	} else if !sim.executePhase25 {
		if (sim.Encounter.EndFightAtHealth == 0 && sim.CurrentTime >= sim.executePhase25Begins) ||
			(sim.Encounter.EndFightAtHealth > 0 && sim.GetRemainingDurationPercent() <= 0.25) {
			sim.executePhase25 = true
			for _, callback := range sim.executePhaseCallbacks {
				callback(sim, 25)
			}
		}
	} else if !sim.executePhase20 {
		if (sim.Encounter.EndFightAtHealth == 0 && sim.CurrentTime >= sim.executePhase20Begins) ||
			(sim.Encounter.EndFightAtHealth > 0 && sim.GetRemainingDurationPercent() <= 0.2) {
			sim.executePhase20 = true
			for _, callback := range sim.executePhaseCallbacks {
				callback(sim, 20)
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

func (sim *Simulation) RegisterExecutePhaseCallback(callback func(*Simulation, int)) {
	sim.executePhaseCallbacks = append(sim.executePhaseCallbacks, callback)
}
func (sim *Simulation) IsExecutePhase20() bool {
	return sim.executePhase20
}
func (sim *Simulation) IsExecutePhase25() bool {
	return sim.executePhase25
}
func (sim *Simulation) IsExecutePhase35() bool {
	return sim.executePhase35
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
