package tank

import (
	"github.com/wowsims/wotlk/sim/core"
)

type BloodTankRotation struct {
	dsCount int
	itCount int
	itCycle bool
}

func (btr *BloodTankRotation) Reset(sim *core.Simulation) {
	btr.dsCount = 0
	btr.itCycle = false
}
