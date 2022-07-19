package deathknight

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

func (deathKnight *DeathKnight) registerHowlingBlastSpell() {
	if !deathKnight.Talents.HowlingBlast {
		return
	}
	target := deathKnight.CurrentTarget

	deathKnight.HowlingBlast = deathKnight.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 51411},
		SpellSchool: core.SpellSchoolFrost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				cast.GCD = deathKnight.getModifiedGCD()
			},
			CD: core.Cooldown{
				Timer:    deathKnight.NewTimer(),
				Duration: 8.0 * time.Second,
			},
		},

		ApplyEffects: core.ApplyEffectFuncAOEDamage(deathKnight.Env, core.SpellEffect{
			ProcMask:             core.ProcMaskSpellDamage,
			BonusSpellCritRating: 0.0,
			DamageMultiplier:     1.0,
			ThreatMultiplier:     1.0,

			BaseDamage: core.BaseDamageConfig{
				Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
					roll := (562.0-518.0)*sim.RandomFloat("Howling Blast") + 518.0
					return (roll + deathKnight.applyImpurity(hitEffect, spell.Unit)*0.1) *
						deathKnight.glacielRotBonus(hitEffect.Target) *
						deathKnight.rageOfRivendareBonus(hitEffect.Target) *
						deathKnight.tundraStalkerBonus(hitEffect.Target) *
						deathKnight.mercilessCombatBonus(sim)
				},
				TargetSpellCoefficient: 1,
			},
			OutcomeApplier: deathKnight.killingMachineOutcomeMod(deathKnight.OutcomeFuncMagicHitAndCrit(deathKnight.spellCritMultiplierGuile())),
			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spellEffect.Target == deathKnight.CurrentTarget {
					deathKnight.LastCastOutcome = spellEffect.Outcome
				}
				if spellEffect.Landed() && target == spellEffect.Target {
					if !deathKnight.HowlingBlastCostless {
						dkSpellCost := deathKnight.DetermineOptimalCost(sim, 0, 1, 1)
						deathKnight.Spend(sim, spell, dkSpellCost)
						amountOfRunicPower := 15.0 + 2.5*float64(deathKnight.Talents.ChillOfTheGrave)
						deathKnight.AddRunicPower(sim, amountOfRunicPower, spell.RunicPowerMetrics())
					} else {
						deathKnight.HowlingBlastCostless = false
					}
				} else if spellEffect.Landed() && !deathKnight.HowlingBlastCostless {
					amountOfRunicPower := 2.5 * float64(deathKnight.Talents.ChillOfTheGrave)
					deathKnight.AddRunicPower(sim, amountOfRunicPower, spell.RunicPowerMetrics())
				}
			},
		}),
	})
}

func (deathKnight *DeathKnight) CanHowlingBlast(sim *core.Simulation) bool {
	if deathKnight.HowlingBlastCostless {
		return deathKnight.HowlingBlast.IsReady(sim)
	}
	return deathKnight.CastCostPossible(sim, 0.0, 0, 1, 1) && deathKnight.HowlingBlast.IsReady(sim)
}
