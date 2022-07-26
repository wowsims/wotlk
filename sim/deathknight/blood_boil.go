package deathknight

import (
	"github.com/wowsims/wotlk/sim/core"
)

func (dk *Deathknight) registerBloodBoilSpell() {
	dk.BloodBoil = dk.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 49941},
		SpellSchool: core.SpellSchoolShadow,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				cast.GCD = dk.getModifiedGCD()
			},
		},

		ApplyEffects: core.ApplyEffectFuncAOEDamage(dk.Env, core.SpellEffect{
			ProcMask:             core.ProcMaskSpellDamage,
			BonusSpellCritRating: 0.0,
			DamageMultiplier:     1.0,
			ThreatMultiplier:     1.0,

			BaseDamage: core.BaseDamageConfig{
				Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
					roll := (220.0-180.0)*sim.RandomFloat("Blood Boil") + 180.0
					return (roll + dk.applyImpurity(hitEffect, spell.Unit)*0.06) *
						dk.rageOfRivendareBonus(hitEffect.Target) *
						dk.tundraStalkerBonus(hitEffect.Target)
				},
				TargetSpellCoefficient: 1,
			},
			OutcomeApplier: dk.OutcomeFuncMagicHitAndCrit(dk.spellCritMultiplierGoGandMoM()),

			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spellEffect.Target == dk.CurrentTarget {
					dk.LastCastOutcome = spellEffect.Outcome
				}
				if spellEffect.Landed() && spellEffect.Target == dk.CurrentTarget {
					dkSpellCost := dk.DetermineCost(sim, core.DKCastEnum_B)
					dk.Spend(sim, spell, dkSpellCost)

					amountOfRunicPower := 10.0
					dk.AddRunicPower(sim, amountOfRunicPower, spell.RunicPowerMetrics())
				}
			},
		}),
	})
}

func (dk *Deathknight) CanBloodBoil(sim *core.Simulation) bool {
	return dk.CastCostPossible(sim, 0.0, 1, 0, 0) && dk.BloodBoil.IsReady(sim)
}

func (dk *Deathknight) CastBloodBoil(sim *core.Simulation, target *core.Unit) bool {
	if dk.CanBloodBoil(sim) {
		dk.BloodBoil.Cast(sim, target)
		return true
	}
	return false
}
