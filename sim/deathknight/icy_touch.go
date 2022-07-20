package deathknight

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (deathKnight *DeathKnight) killingMachineOutcomeMod(outcomeApplier core.OutcomeApplier) core.OutcomeApplier {
	return func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect, attackTable *core.AttackTable) {
		if deathKnight.KillingMachineAura.IsActive() {
			deathKnight.AddStatDynamic(sim, stats.SpellCrit, 100*core.CritRatingPerCritChance)
			outcomeApplier(sim, spell, spellEffect, attackTable)
			deathKnight.AddStatDynamic(sim, stats.SpellCrit, -100*core.CritRatingPerCritChance)
		} else {
			outcomeApplier(sim, spell, spellEffect, attackTable)
		}
	}
}

func (deathKnight *DeathKnight) registerIcyTouchSpell() {
	target := deathKnight.CurrentTarget

	itAura := core.IcyTouchAura(target, deathKnight.Talents.ImprovedIcyTouch)
	deathKnight.IcyTouchAura = itAura

	impIcyTouchCoeff := 1.0 + 0.05*float64(deathKnight.Talents.ImprovedIcyTouch)

	deathKnight.IcyTouch = deathKnight.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 59131},
		SpellSchool: core.SpellSchoolFrost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				cast.GCD = deathKnight.getModifiedGCD()
			},
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:             core.ProcMaskSpellDamage,
			BonusSpellCritRating: 5.0 * float64(deathKnight.Talents.Rime) * core.CritRatingPerCritChance,
			DamageMultiplier:     impIcyTouchCoeff,
			ThreatMultiplier:     7.0,

			BaseDamage: core.BaseDamageConfig{
				Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
					roll := (245.0-227.0)*sim.RandomFloat("Icy Touch") + 227.0 + deathKnight.sigilOfTheFrozenConscienceBonus()
					return (roll + deathKnight.applyImpurity(hitEffect, spell.Unit)*0.1) *
						deathKnight.glacielRotBonus(hitEffect.Target) *
						deathKnight.rageOfRivendareBonus(hitEffect.Target) *
						deathKnight.tundraStalkerBonus(hitEffect.Target) *
						deathKnight.mercilessCombatBonus(sim)
				},
				TargetSpellCoefficient: 1,
			},
			OutcomeApplier: deathKnight.killingMachineOutcomeMod(deathKnight.OutcomeFuncMagicHitAndCrit(deathKnight.spellCritMultiplier())),

			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				deathKnight.LastCastOutcome = spellEffect.Outcome
				if spellEffect.Landed() {
					dkSpellCost := deathKnight.DetermineOptimalCost(sim, 0, 1, 0)
					deathKnight.Spend(sim, spell, dkSpellCost)

					deathKnight.FrostFeverSpell.Cast(sim, spellEffect.Target)
					if deathKnight.Talents.EbonPlaguebringer > 0 {
						deathKnight.EbonPlagueAura.Activate(sim)
					}

					amountOfRunicPower := 10.0 + 2.5*float64(deathKnight.Talents.ChillOfTheGrave)
					deathKnight.AddRunicPower(sim, amountOfRunicPower, spell.RunicPowerMetrics())

					deathKnight.IcyTouchAura.Activate(sim)

					if deathKnight.IcyTouchAura.IsActive() && deathKnight.IcyTalonsAura != nil {
						deathKnight.IcyTalonsAura.Activate(sim)
					}
				}
			},
		}),
	})
}

func (deathKnight *DeathKnight) CanIcyTouch(sim *core.Simulation) bool {
	return deathKnight.CastCostPossible(sim, 0.0, 0, 1, 0) && deathKnight.IcyTouch.IsReady(sim)
}

func (deathKnight *DeathKnight) CastIcyTouch(sim *core.Simulation, target *core.Target) bool {
	if deathKnight.CanIcyTouch(sim) {
		deathKnight.CastIcyTouch(sim, target)
		return true
	}
	return false
}
