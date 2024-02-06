package rogue

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/stats"
)

func (rogue *Rogue) ApplyRunes() {
	// Apply runes here :)
}

func (rogue *Rogue) applyMutilate() {
	if !rogue.HasRune(proto.RogueRune_RuneMutilate) {
		return
	}
}
