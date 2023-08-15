package rogue

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

// Each rank is a different ID. 31223 is 3/3
func getMasterOfSubtletySpellID(talentPoints int32) int32 {
	return []int32{0, 31221, 31222, 31223}[talentPoints]
}

func (rogue *Rogue) registerMasterOfSubtletyCD() {
	var MasterOfSubtletyID = core.ActionID{SpellID: getMasterOfSubtletySpellID(rogue.Talents.MasterOfSubtlety)}

	percent := []float64{1, 1.04, 1.07, 1.1}[rogue.Talents.MasterOfSubtlety]

	effectDuration := time.Second * 6
	if rogue.StealthAura.IsActive() {
		effectDuration = core.NeverExpires
	}

	rogue.MasterOfSubtletyAura = rogue.RegisterAura(core.Aura{
		Label:    "Master of Subtlety",
		ActionID: MasterOfSubtletyID,
		Duration: effectDuration,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			rogue.PseudoStats.DamageDealtMultiplier *= percent
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			rogue.PseudoStats.DamageDealtMultiplier *= 1 / percent
		},
	})
}
