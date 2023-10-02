package deathknight

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

// Time between focus ticks.
const MaxFocus = 100.0
const tickDuration = time.Second
const BaseFocusPerTick = 10.0

// OnFocusGain is called any time focus is increased.
type OnFocusGain func(sim *core.Simulation)

type focusBar struct {
	ghoulPet *GhoulPet

	focusPerTick float64

	currentFocus float64

	onFocusGain OnFocusGain
	tickAction  *core.PendingAction
}

func (ghoulPet *GhoulPet) EnableFocusBar(onFocusGain OnFocusGain) {
	ghoulPet.focusBar = focusBar{
		ghoulPet:     ghoulPet,
		focusPerTick: BaseFocusPerTick,
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

	newFocus := min(fb.currentFocus+amount, MaxFocus)

	if sim.Log != nil {
		fb.ghoulPet.Log(sim, "Gained %0.3f focus from %s (%0.3f --> %0.3f).", amount, actionID, fb.currentFocus, newFocus)
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
		fb.ghoulPet.Log(sim, "Spent %0.3f focus from %s (%0.3f --> %0.3f).", amount, actionID, fb.currentFocus, newFocus)
	}

	fb.currentFocus = newFocus
}

func (fb *focusBar) Enable(sim *core.Simulation) {
	if fb.ghoulPet == nil {
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

func (fb *focusBar) Disable(sim *core.Simulation) {
	if fb.tickAction != nil {
		fb.tickAction.Cancel(sim)
		fb.tickAction = nil
	}
}
