package restoration

import (
	"time"

	"github.com/wowsims/classic/sod/sim/core"
)

func (resto *RestorationDruid) OnGCDReady(sim *core.Simulation) {
	resto.tryUseGCD(sim)
}

func (resto *RestorationDruid) tryUseGCD(sim *core.Simulation) {
	resto.WaitUntil(sim, sim.CurrentTime+time.Second*5)
}
