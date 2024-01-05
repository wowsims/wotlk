package rogue

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
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
				Duration: time.Second * time.Duration(180-30*rogue.Talents.Elusiveness),
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

		ShouldActivate: func(s *core.Simulation, c *core.Character) bool {
			if rogue.Talents.Overkill {
				return !(rogue.StealthAura.IsActive() || rogue.OverkillAura.IsActive()) && rogue.CurrentEnergy() > 50
			}
			if rogue.Talents.MasterOfSubtlety > 0 {
				// Chained cast checks
				// heuristically, 3 Garrote ticks are better DPE than regular builders
				const garroteMinDuration = time.Second * 9

				if rogue.MasterOfSubtletyAura.IsActive() {
					return false // possible after preparation
				}

				if s.GetRemainingDuration() < garroteMinDuration {
					return true // getting the buff up under non-ideal circumstances is fine at end of combat
				}

				wantPremed, premedCPs := checkPremediation(s, rogue)
				if wantPremed && premedCPs == 0 {
					return false // essentially sync with premed if possible
				}

				wantGarrote, garroteCPs := checkGarrote(s, rogue)
				if wantGarrote && garroteCPs == 0 {
					return false
				}

				return rogue.ComboPoints()+garroteCPs+premedCPs <= 5+1 // heuristically, "<= 5" is too strict (since omitting premed is fine)
			}

			return false
		},
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
