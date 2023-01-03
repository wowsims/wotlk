package restoration

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

func (resto *RestorationDruid) OnGCDReady(sim *core.Simulation) {
	resto.tryUseGCD(sim)
}

func (resto *RestorationDruid) tryUseGCD(sim *core.Simulation) {
	resto.WaitUntil(sim, sim.CurrentTime+time.Second*5)
}
