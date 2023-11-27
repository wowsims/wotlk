package balance

import (
	"github.com/wowsims/classic/sod/sim/core"
	"github.com/wowsims/classic/sod/sim/druid"
)

func (moonkin *BalanceDruid) OnGCDReady(sim *core.Simulation) {
	moonkin.tryUseGCD(sim)
}

func (moonkin *BalanceDruid) tryUseGCD(sim *core.Simulation) {
	spell, target := moonkin.rotation(sim)
	if success := spell.Cast(sim, target); !success {
		moonkin.WaitForMana(sim, spell.CurCast.Cost)
	}
}

func (moonkin *BalanceDruid) rotation(sim *core.Simulation) (*druid.DruidSpell, *core.Unit) {
	moonkin.CurrentTarget = sim.Environment.GetTargetUnit(0)
	target := moonkin.CurrentTarget

	return moonkin.Wrath, target
}
