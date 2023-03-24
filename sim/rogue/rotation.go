package rogue

import (
	"github.com/wowsims/wotlk/sim/core"
)

type rotation interface {
	setup(sim *core.Simulation, rogue *Rogue)
	run(sim *core.Simulation, rogue *Rogue)
}

type PriorityAction int32

const (
	Skip PriorityAction = iota
	Build
	Cast
	Wait
	Once
)

type prio struct {
	check func(sim *core.Simulation, rogue *Rogue) PriorityAction
	cast  func(sim *core.Simulation, rogue *Rogue) bool
	cost  float64
}

func (rogue *Rogue) OnEnergyGain(sim *core.Simulation) {
	rogue.TryUseCooldowns(sim)

	if !rogue.GCD.IsReady(sim) {
		return
	}

	rogue.rotation.run(sim, rogue)
}

func (rogue *Rogue) OnGCDReady(sim *core.Simulation) {
	rogue.TryUseCooldowns(sim)

	if rogue.IsWaitingForEnergy() {
		rogue.DoNothing()
		return
	}

	rogue.rotation.run(sim, rogue)
}

func (rogue *Rogue) setupRotation(sim *core.Simulation) {
	switch {
	case rogue.CanMutilate() && rogue.Env.GetNumTargets() <= 3:
		rogue.rotation = &rotation_assassination{}
	case rogue.Talents.CombatPotency > 0 && rogue.Env.GetNumTargets() <= 3:
		rogue.rotation = &rotation_combat{}
	case rogue.Talents.HonorAmongThieves > 0 && rogue.Env.GetNumTargets() <= 3:
		rogue.rotation = &rotation_subtlety{}
	default:
		rogue.rotation = &rotation_generic{}
	}
	rogue.rotation.setup(sim, rogue)
}
