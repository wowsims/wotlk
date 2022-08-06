package dps

import (
	"github.com/wowsims/wotlk/sim/core"
)

type FrostRotation struct {
	oblitCount  int32
	missedPesti bool
}

func (fr *FrostRotation) Reset(sim *core.Simulation) {
	fr.oblitCount = 0
	fr.missedPesti = false
}
