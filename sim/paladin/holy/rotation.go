package holy

import (
	"github.com/wowsims/wotlk/sim/core"
)

func (holy *HolyPaladin) OnGCDReady(sim *core.Simulation) {
	holy.DoNothing()
}
