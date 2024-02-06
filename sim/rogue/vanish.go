package rogue

import (
	"time"

	"github.com/wowsims/sod/sim/core"
)

func (rogue *Rogue) registerVanishSpell() {
	rogue.Vanish = rogue.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 26889},
		SpellSchool: core.SpellSchoolPhysical,
		Flags:       core.SpellFlagAPL,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: 0,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    rogue.NewTimer(),
				Duration: time.Second * time.Duration(300-45*rogue.Talents.Elusiveness),
			},
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			// Pause auto attacks
			rogue.AutoAttacks.CancelAutoSwing(sim)
			// Apply stealth
			rogue.StealthAura.Activate(sim)
		},
	})

	rogue.AddMajorCooldown(core.MajorCooldown{
		Spell:    rogue.Vanish,
		Type:     core.CooldownTypeDPS,
		Priority: core.CooldownPriorityDrums,
	})
}

const garroteMinDuration = time.Second * 9 // heuristically, 3 Garrote ticks are better DPE than regular builders

func checkGarrote(sim *core.Simulation, rogue *Rogue) (bool, int32) {
	initiative := core.Ternary[int32](rogue.Talents.Initiative == 0, 0, 1)
	// Garrote cannot be cast in front of the target
	if rogue.PseudoStats.InFrontOfTarget {
		return false, 0
	}

	if !rogue.GCD.IsReady(sim) || rogue.CurrentEnergy() < rogue.Garrote.DefaultCast.Cost {
		return false, 0
	}

	// Garrote Clip logic
	if rogue.GCD.IsReady(sim) && rogue.Garrote.CurDot().IsActive() && sim.GetRemainingDuration() <= garroteMinDuration {
		return true, 1 + initiative
	}

	return true, 1 + initiative
}

func checkPremediation(sim *core.Simulation, rogue *Rogue) (bool, int32) {
	if rogue.Premeditation == nil {
		return false, 0
	}

	if !rogue.Premeditation.IsReady(sim) {
		return false, 0
	}
	return true, 2
}
