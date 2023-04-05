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

	dk.HowlingBlast = dk.RegisterSpell(core.SpellConfig{
		ActionID:    HowlingBlastActionID,
		SpellSchool: core.SpellSchoolFrost,
		ProcMask:    core.ProcMaskSpellDamage,

		RuneCost: core.RuneCostOptions{
			FrostRuneCost:  1,
			UnholyRuneCost: 1,
			RunicPowerGain: 15,
			Refundable:     true,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    dk.NewTimer(),
				Duration: 8 * time.Second,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   dk.bonusCritMultiplier(dk.Talents.GuileOfGorefiend),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			for _, aoeTarget := range sim.Encounter.TargetUnits {
				baseDamage := (sim.Roll(518, 562) + 0.2*dk.getImpurityBonus(spell)) *
					dk.glacielRotBonus(aoeTarget) *
					dk.RoRTSBonus(aoeTarget) *
					dk.mercilessCombatBonus(sim) *
					sim.Encounter.AOECapMultiplier()

				result := spell.CalcDamage(sim, aoeTarget, baseDamage, spell.OutcomeMagicHitAndCrit)

				if aoeTarget == target {
					spell.SpendRefundableCost(sim, result)
					dk.LastOutcome = result.Outcome
				}
				if rpBonus > 0 && result.Landed() {
					dk.AddRunicPower(sim, rpBonus, spell.RunicPowerMetrics())
				}

				if hasGlyph {
					dk.FrostFeverSpell.Cast(sim, aoeTarget)
				}

				spell.DealDamage(sim, result)
			}

			if dk.RimeAura.IsActive() {
				dk.RimeAura.Deactivate(sim)
			}
			if dk.KillingMachineAura.IsActive() {
				dk.KillingMachineAura.Deactivate(sim)
			}
		},
	})
}
