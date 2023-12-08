package shadow

import (
	"github.com/wowsims/sod/sim/core"
)

func (spriest *ShadowPriest) OnGCDReady(sim *core.Simulation) {
	spriest.DoNothing()
}
