package restoration

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

func (resto *RestorationShaman) OnGCDReady(sim *core.Simulation) {
	resto.tryUseGCD(sim)
}

func (resto *RestorationShaman) tryUseGCD(sim *core.Simulation) {
	resto.WaitUntil(sim, sim.CurrentTime+time.Second*5)
}
