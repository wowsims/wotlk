package rogue

import "github.com/wowsims/sod/sim/core/proto"

func (rogue *Rogue) ApplyRunes() {
	// Apply runes here :)
}

func (rogue *Rogue) applyMutilate() {
	if !rogue.HasRune(proto.RogueRune_RuneMutilate) {
		return
	}
}
