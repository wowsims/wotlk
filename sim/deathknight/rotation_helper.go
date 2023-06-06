package deathknight

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

// return bool is if its on GCD
// return duration is an optional wait time
type RotationAction func(sim *core.Simulation, target *core.Unit, s *Sequence) time.Duration

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

func (s *Sequence) IsOngoing() bool {
	return s.idx < s.numActions
}

func (s *Sequence) RemainingActions() int {
	return (s.numActions - 1) - s.idx
}

func (s *Sequence) Reset() {
	s.idx = 0
}

func (s *Sequence) Advance() {
	s.idx += 1
}

func (s *Sequence) ConditionalAdvance(condition bool) {
	if condition {
		s.idx += 1
	}
}

func (s *Sequence) GetNextAction() RotationAction {
	if s.idx+1 < s.numActions {
		return s.actions[s.idx+1]
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
	s.actions = s.actions[:0]
	s.numActions = 0
	s.idx = 0
	return s
}

type RotationHelper struct {
	RotationSequence *Sequence

	LastOutcome           core.HitOutcome
	LastCast              *core.Spell
	NextCast              *core.Spell
	AoESpellNumTargetsHit int32
}
