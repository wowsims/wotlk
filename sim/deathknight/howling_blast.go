package deathknight

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

var HowlingBlastActionID = core.ActionID{SpellID: 51411}

func (dk *Deathknight) registerHowlingBlastSpell() {
	if !dk.Talents.HowlingBlast {
		return
	}

	rpBonus := 2.5 * float64(dk.Talents.ChillOfTheGrave)
	hasGlyph := dk.HasMajorGlyph(proto.DeathknightMajorGlyph_GlyphOfHowlingBlast)

	howlingBlast := &RuneSpell{
		Refundable: true,
	}
	dk.HowlingBlast = dk.RegisterSpell(howlingBlast, core.SpellConfig{
		ActionID:    HowlingBlastActionID,
		SpellSchool: core.SpellSchoolFrost,
		ProcMask:    core.ProcMaskSpellDamage,

		RuneCost: core.RuneCostOptions{
			FrostRuneCost:  1,
			UnholyRuneCost: 1,
			RunicPowerGain: 15,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				cast.GCD = dk.GetModifiedGCD()
				if dk.RimeAura.IsActive() {
					cast.Cost = 0 // no runes, no regen
					dk.RimeAura.Deactivate(sim)
				}
			},
			CD: core.Cooldown{
				Timer:    dk.NewTimer(),
				Duration: 8.0 * time.Second,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   dk.bonusCritMultiplier(dk.Talents.GuileOfGorefiend),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			for _, aoeTarget := range sim.Encounter.Targets {
				aoeUnit := &aoeTarget.Unit

				baseDamage := (sim.Roll(518, 562) + 0.2*dk.getImpurityBonus(spell)) *
					dk.glacielRotBonus(aoeUnit) *
					dk.RoRTSBonus(aoeUnit) *
					dk.mercilessCombatBonus(sim) *
					sim.Encounter.AOECapMultiplier()

				result := spell.CalcDamage(sim, aoeUnit, baseDamage, spell.OutcomeMagicHitAndCrit)

				if aoeUnit == dk.CurrentTarget {
					howlingBlast.OnResult(sim, result)
					dk.LastOutcome = result.Outcome
				}
				if dk.Talents.ChillOfTheGrave > 0 && result.Landed() {
					dk.AddRunicPower(sim, rpBonus, spell.RunicPowerMetrics())
				}

				if hasGlyph {
					dk.FrostFeverSpell.Cast(sim, aoeUnit)
				}

				spell.DealDamage(sim, result)
			}

			if dk.KillingMachineAura.IsActive() {
				dk.KillingMachineAura.Deactivate(sim)
			}
		},
	})
}
