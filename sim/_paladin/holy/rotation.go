package holy

import (
	"time"

	"github.com/wowsims/classic/sod/sim/core"
)

func (holy *HolyPaladin) OnGCDReady(sim *core.Simulation) {
	holy.WaitUntil(sim, sim.CurrentTime+time.Second*5)
}
