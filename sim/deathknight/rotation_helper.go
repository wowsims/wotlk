package deathknight

import (
	"github.com/wowsims/wotlk/sim/core"
)

type RotationAction func(sim *core.Simulation, target *core.Unit, s *Sequence) bool

func TernaryRotationAction(condition bool, t RotationAction, f RotationAction) RotationAction {
	if condition {
		return t
	} else {
		return f
	}
}

// Add your UH rotation Actions here and then on the DoNext function

type Sequence struct {
	idx        int
	numActions int
	actions    []RotationAction
}

func (o *Sequence) IsOngoing() bool {
	return o.idx < o.numActions
}

func (o *Sequence) Reset() {
	o.idx = 0
}

func (o *Sequence) Advance(condition bool) {
	o.idx += 1
}

func (o *Sequence) ConditionalAdvance(condition bool) {
	if condition {
		o.idx += 1
	}
}

func (o *Sequence) GetNextAction() RotationAction {
	if o.idx+1 < o.numActions {
		return o.actions[o.idx+1]
	} else {
		return nil
	}
}

func (s *Sequence) NewAction(action RotationAction) *Sequence {
	s.actions = append(s.actions, action)
	s.numActions += 1
	return s
}

func (s *Sequence) Clear() *Sequence {
	s.actions = make([]RotationAction, 0)
	s.numActions = 0
	s.idx = 0
	return s
}

type RotationHelper struct {
	Opener *Sequence
	Main   *Sequence
}
