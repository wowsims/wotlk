package deathknight

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

func (dk *Deathknight) registerHowlingBlastSpell() {
	if !dk.Talents.HowlingBlast {
		return
	}

	dk.HowlingBlast = dk.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 51411},
		SpellSchool: core.SpellSchoolFrost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				cast.GCD = dk.getModifiedGCD()
			},
			CD: core.Cooldown{
				Timer:    dk.NewTimer(),
				Duration: 8.0 * time.Second,
			},
		},

		ApplyEffects: core.ApplyEffectFuncAOEDamage(dk.Env, core.SpellEffect{
			ProcMask:             core.ProcMaskSpellDamage,
			BonusSpellCritRating: 0.0,
			DamageMultiplier:     1.0,
			ThreatMultiplier:     1.0,

			BaseDamage: core.BaseDamageConfig{
				Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
					roll := (562.0-518.0)*sim.RandomFloat("Howling Blast") + 518.0
					return (roll + dk.getImpurityBonus(hitEffect, spell.Unit)*0.1) *
						dk.glacielRotBonus(hitEffect.Target) *
						dk.rageOfRivendareBonus(hitEffect.Target) *
						dk.tundraStalkerBonus(hitEffect.Target) *
						dk.mercilessCombatBonus(sim)
				},
				TargetSpellCoefficient: 1,
			},
			OutcomeApplier: dk.killingMachineOutcomeMod(dk.OutcomeFuncMagicHitAndCrit(dk.spellCritMultiplierGoGandMoM())),
			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spellEffect.Target == dk.CurrentTarget {
					dk.LastCastOutcome = spellEffect.Outcome
				}
				if spellEffect.Landed() {
					if dk.KillingMachineAura.IsActive() {
						dk.KillingMachineAura.Deactivate(sim)
					}
					if dk.CurrentTarget == spellEffect.Target {
						if !dk.RimeAura.IsActive() {
							dkSpellCost := dk.DetermineCost(sim, core.DKCastEnum_FU)
							dk.Spend(sim, spell, dkSpellCost)
							amountOfRunicPower := 15.0 + 2.5*float64(dk.Talents.ChillOfTheGrave)
							dk.AddRunicPower(sim, amountOfRunicPower, spell.RunicPowerMetrics())
						} else {
							dk.RimeAura.Deactivate(sim)
							amountOfRunicPower := 2.5 * float64(dk.Talents.ChillOfTheGrave)
							dk.AddRunicPower(sim, amountOfRunicPower, spell.RunicPowerMetrics())
						}
					} else {
						amountOfRunicPower := 2.5 * float64(dk.Talents.ChillOfTheGrave)
						dk.AddRunicPower(sim, amountOfRunicPower, spell.RunicPowerMetrics())
					}
				} else if dk.RimeAura.IsActive() && dk.CurrentTarget == spellEffect.Target {
					dk.RimeAura.Deactivate(sim)
				}
			},
		}),
	})
}

func (dk *Deathknight) CanHowlingBlast(sim *core.Simulation) bool {
	if dk.RimeAura.IsActive() {
		return dk.HowlingBlast.IsReady(sim)
	}
	return dk.CastCostPossible(sim, 0.0, 0, 1, 1) && dk.HowlingBlast.IsReady(sim)
}

func (dk *Deathknight) CastHowlingBlast(sim *core.Simulation, target *core.Unit) bool {
	if dk.CanHowlingBlast(sim) {
		dk.HowlingBlast.Cast(sim, target)
		return true
	}
	return false
}
