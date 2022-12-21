package rogue

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

// Each rank is a different ID. 31223 is 3/3
func getMasterofSubtletySpellID(talentPoints int32) int32 {
	if talentPoints == 1 {
		return 31221
	}
	return 31220 + talentPoints
}

func getMasterofSubtletyCDDuration(talentPoints int32) int32 {
	if talentPoints > 0 {
		return 180 - 30*talentPoints
	}
	return 180
}

func (rogue *Rogue) registerMasterOfSubtletyCD() {
	if rogue.Talents.MasterOfSubtlety == 0 {
		return
	}

	var MasterOfSubtletyID = core.ActionID{SpellID: getMasterofSubtletySpellID(rogue.Talents.MasterOfSubtlety)}

	percent := 0.04

	if rogue.Talents.MasterOfSubtlety > 1 {
		percent += 0.03 * float64(rogue.Talents.MasterOfSubtlety)
	}

	rogue.MasterOfSubtletyAura = rogue.RegisterAura(core.Aura{
		Label:    "Master of Subtlety",
		ActionID: MasterOfSubtletyID,
		Duration: time.Second * 6,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {

			rogue.MultiplyStat(stats.AttackPower, 1.0+percent)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			rogue.MultiplyStat(stats.AttackPower, 1.0/(1.0+percent))
		},
	})
	masterOfSubtletySpell := rogue.RegisterSpell(core.SpellConfig{
		ActionID: MasterOfSubtletyID,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: 0,
				GCD:  time.Second * 0,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    rogue.NewTimer(),
				Duration: time.Second * time.Duration(getMasterofSubtletyCDDuration(rogue.Talents.MasterOfSubtlety)),
			},
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			rogue.MasterOfSubtletyAura.Activate(sim)
		},
	})

	rogue.AddMajorCooldown(core.MajorCooldown{
		Spell: masterOfSubtletySpell,
		Type:  core.CooldownTypeDPS,
		ShouldActivate: func(s *core.Simulation, c *core.Character) bool {
			return rogue.CurrentEnergy() > 90
		},
	})
}
