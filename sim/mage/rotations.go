package mage

import (
	"github.com/wowsims/classic/sim/core"
)

func (mage *Mage) OnGCDReady(sim *core.Simulation) {
	mage.tryUseGCD(sim)
}

func (mage *Mage) tryUseGCD(sim *core.Simulation) {
	if mage.IsUsingAPL {
		return
	}

	mage.DoNothing()
}
