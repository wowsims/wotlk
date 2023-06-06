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

	const garroteMinDuration = time.Second * 9 // heuristically, 3 Garrote ticks are better DPE than regular builders

	garrote := func(sim *core.Simulation, rogue *Rogue) (bool, int32) {
		if !rogue.Rotation.OpenWithGarrote || rogue.PseudoStats.InFrontOfTarget {
			return false, 0
		}

		if !rogue.GCD.IsReady(sim) || rogue.CurrentEnergy() < rogue.Garrote.DefaultCast.Cost {
			return true, 0
		}

		if rogue.Garrote.CurDot().IsActive() || sim.GetRemainingDuration() <= garroteMinDuration {
			return true, 0
		}

		if rogue.Talents.Initiative == 0 {
			return true, 1
		}
		return true, 2
	}

	premed := func(sim *core.Simulation, rogue *Rogue) (bool, int32) {
		if rogue.Premeditation == nil {
			return false, 0
		}

		if !rogue.Premeditation.IsReady(sim) {
			return true, 0
		}
		return true, 2
	}

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
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			rogue.MasterOfSubtletyAura.Activate(sim)

			_, premedCPs := premed(sim, rogue)
			_, garroteCPs := garrote(sim, rogue)

			if premedCPs > 0 && rogue.ComboPoints()+premedCPs+garroteCPs <= 5 {
				rogue.Premeditation.Cast(sim, target)
			}

			if garroteCPs > 0 {
				rogue.Garrote.Cast(sim, target)
			}
		},
	})

	rogue.AddMajorCooldown(core.MajorCooldown{
		Spell: rogue.MasterOfSubtlety,
		Type:  core.CooldownTypeDPS,
		ShouldActivate: func(sim *core.Simulation, character *core.Character) bool {
			if rogue.MasterOfSubtletyAura.IsActive() {
				return false // possible after preparation
			}

			if sim.GetRemainingDuration() < garroteMinDuration {
				return true // getting the buff up under non-ideal circumstances is fine at end of combat
			}

			wantPremed, premedCPs := premed(sim, rogue)
			if wantPremed && premedCPs == 0 {
				return false // essentially sync with premed if possible
			}

			wantGarrote, garroteCPs := garrote(sim, rogue)
			if wantGarrote && garroteCPs == 0 {
				return false
			}

			return rogue.ComboPoints()+garroteCPs+premedCPs <= 5+1 // heuristically, "<= 5" is too strict (since omitting premed is fine)
		},
	})
}
