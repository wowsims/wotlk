package enhancement

import (
	"github.com/wowsims/wotlk/sim/core"
)

func (enh *EnhancementShaman) OnAutoAttack(sim *core.Simulation, spell *core.Spell) {
}

func (enh *EnhancementShaman) OnGCDReady(sim *core.Simulation) {
	// TODO move this into the rotation, also this uses waitForMana if it was unable to cast the totem
	// that will need to be pulled out so we are not waiting for a magma totem mana cost.
	// if enh.TryDropTotems(sim) {
	// 	return
	// }
	enh.rotation.DoAction(enh, sim)
}

type Rotation interface {
	DoAction(*EnhancementShaman, *core.Simulation)
	Reset(*EnhancementShaman, *core.Simulation)
}

//	CUSTOM ROTATION (advanced) (also WIP).
//
// TODO: figure out how to do this (probably too complicated to copy hunters)
type AgentAction interface {
	GetActionID() core.ActionID

	GetManaCost() float64

	Cast(sim *core.Simulation) bool
}
