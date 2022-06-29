package hunter

import (
	"time"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/proto"
)

// Time between focus ticks.
const MaxFocus = 100.0
const tickDuration = time.Second * 5
const BaseFocusPerTick = 25.0

// OnFocusGain is called any time focus is increased.
type OnFocusGain func(sim *core.Simulation)

type focusBar struct {
	hunterPet *HunterPet

	focusPerTick float64

	currentFocus float64

	onFocusGain OnFocusGain
	tickAction  *core.PendingAction
}

func (hunterPet *HunterPet) EnableFocusBar(regenMultiplier float64, onFocusGain OnFocusGain) {
	hunterPet.focusBar = focusBar{
		hunterPet:    hunterPet,
		focusPerTick: BaseFocusPerTick * regenMultiplier,
		onFocusGain:  onFocusGain,
	}
}

func (fb *focusBar) CurrentFocus() float64 {
	return fb.currentFocus
}

func (fb *focusBar) AddFocus(sim *core.Simulation, amount float64, actionID core.ActionID) {
	if amount < 0 {
		panic("Trying to add negative focus!")
	}

	newFocus := core.MinFloat(fb.currentFocus+amount, MaxFocus)

	if sim.Log != nil {
		fb.hunterPet.Log(sim, "Gained %0.3f focus from %s (%0.3f --> %0.3f).", amount, actionID, fb.currentFocus, newFocus)
	}

	fb.currentFocus = newFocus

	if fb.onFocusGain != nil {
		fb.onFocusGain(sim)
	}
}

func (fb *focusBar) SpendFocus(sim *core.Simulation, amount float64, actionID core.ActionID) {
	if amount < 0 {
		panic("Trying to spend negative focus!")
	}

	newFocus := fb.currentFocus - amount

	if sim.Log != nil {
		fb.hunterPet.Log(sim, "Spent %0.3f focus from %s (%0.3f --> %0.3f).", amount, actionID, fb.currentFocus, newFocus)
	}

	fb.currentFocus = newFocus
}

func (fb *focusBar) reset(sim *core.Simulation) {
	if fb.hunterPet == nil {
		return
	}

	fb.currentFocus = MaxFocus

	pa := &core.PendingAction{
		Priority:     core.ActionPriorityRegen,
		NextActionAt: tickDuration,
	}
	pa.OnAction = func(sim *core.Simulation) {
		fb.AddFocus(sim, fb.focusPerTick, core.ActionID{OtherID: proto.OtherAction_OtherActionFocusRegen})

		pa.NextActionAt = sim.CurrentTime + tickDuration
		sim.AddPendingAction(pa)
	}
	fb.tickAction = pa
	sim.AddPendingAction(pa)
}

func (fb *focusBar) Cancel(sim *core.Simulation) {
	if fb.tickAction != nil {
		fb.tickAction.Cancel(sim)
		fb.tickAction = nil
	}
}
