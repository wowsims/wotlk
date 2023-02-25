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
	if rogue.Talents.MasterOfSubtlety == 0 {
		return
	}

	var MasterOfSubtletyID = core.ActionID{SpellID: getMasterOfSubtletySpellID(rogue.Talents.MasterOfSubtlety)}

	percent := []float64{1, 1.04, 1.07, 1.1}[rogue.Talents.MasterOfSubtlety]

	rogue.MasterOfSubtletyAura = rogue.RegisterAura(core.Aura{
		Label:    "Master of Subtlety",
		ActionID: MasterOfSubtletyID,
		Duration: time.Second * 6,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			rogue.PseudoStats.DamageDealtMultiplier *= percent
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			rogue.PseudoStats.DamageDealtMultiplier *= 1 / percent
		},
	})

	rogue.MasterOfSubtlety = rogue.RegisterSpell(core.SpellConfig{
		ActionID: MasterOfSubtletyID,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: 0,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    rogue.NewTimer(),
				Duration: time.Second * time.Duration(180-30*rogue.Talents.Elusiveness),
			},
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			rogue.MasterOfSubtletyAura.Activate(sim)
		},
	})

	rogue.AddMajorCooldown(core.MajorCooldown{
		Spell: rogue.MasterOfSubtlety,
		Type:  core.CooldownTypeDPS,
		ShouldActivate: func(s *core.Simulation, c *core.Character) bool {
			if rogue.MasterOfSubtletyAura.IsActive() {
				return false // possible after preparation
			}
			if s.GetRemainingDuration() < time.Second*10 {
				return true // getting the buff up under non-ideal circumstances is fine at end of combat
			}
			if rogue.Premeditation != nil && !rogue.Premeditation.IsReady(s) {
				return false // essentially sync with premed if possible
			}
			return rogue.ComboPoints() <= 1 && rogue.CurrentEnergy() >= 40 // this covers any opener w/ talents, for now
		},
	})
}
