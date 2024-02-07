package rogue

import "github.com/wowsims/sod/sim/core/proto"

func (rogue *Rogue) ApplyRunes() {
	// Apply runes here :)
}

func (rogue *Rogue) applyJustAFleshWound() {
	if !rogue.HasRune(proto.RogueRune_RuneJustAFleshWound) {
		return
	}

	// Increase threat modifier by 112%
	// 6% reduced chance to be critically hit by melee attacks
	// Increase Physical damage reduction by 20% WHILE Blade Dance is active
	// Override Feint to Tease

	// Tease (ID: 410412)
	// 10s cd, Taunt target for 3s?, no effect if already attacking you
}

func (rogue *Rogue) applyDeadlyBrew() {
	if !rogue.HasRune(proto.RogueRune_RuneDeadlyBrew) {
		return
	}

	// If your weapon does not have a poison applied, it has a chance to trigger Instant Poison as if Instant Poison were applied.
	// When you Inflict any other poison on a target, you also inflict Deadly Poison
	// Deadly Poison and Instant Poison now gain increased damage from your Attack Power (9% of AP)
}
