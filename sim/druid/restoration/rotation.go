package restoration

import (
	"github.com/wowsims/wotlk/sim/core"
)

func (resto *RestorationDruid) OnGCDReady(sim *core.Simulation) {
	resto.tryUseGCD(sim)
}

func (resto *RestorationDruid) tryUseGCD(sim *core.Simulation) {
	resto.DoNothing()
}
