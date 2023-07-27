package core

import (
	"github.com/wowsims/wotlk/sim/core/proto"
)

type APLActionChangeTarget struct {
	unit      *Unit
	newTarget UnitReference
}

func (rot *APLRotation) newActionChangeTarget(config *proto.APLActionChangeTarget) APLActionImpl {
	newTarget := rot.getSourceUnit(config.NewTarget)
	if newTarget.Get() == nil {
		return nil
	}
	return &APLActionChangeTarget{
		newTarget: newTarget,
	}
}
func (action *APLActionChangeTarget) GetInnerActions() []*APLAction { return nil }
func (action *APLActionChangeTarget) Finalize(*APLRotation)         {}
func (action *APLActionChangeTarget) Reset(*Simulation)             {}
func (action *APLActionChangeTarget) IsReady(sim *Simulation) bool {
	return true
}
func (action *APLActionChangeTarget) Execute(sim *Simulation) {
	action.unit.CurrentTarget = action.newTarget.Get()
}
