package deathknight

import (
	"github.com/wowsims/wotlk/sim/core"
)

func (deathKnight *DeathKnight) registerIcyTouchSpell() {
	deathKnight.IcyTouchAura = make([]*core.Aura, deathKnight.Env.GetNumTargets())
	for _, encounterTarget := range deathKnight.Env.Encounter.Targets {
		target := &encounterTarget.Unit
		itAura := core.IcyTouchAura(target, deathKnight.Talents.ImprovedIcyTouch)
		deathKnight.IcyTouchAura[target.Index] = itAura
	}

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
			BonusSpellCritRating: deathKnight.rimeCritBonus() * core.CritRatingPerCritChance,
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
					if deathKnight.KillingMachineAura.IsActive() {
						deathKnight.KillingMachineAura.Deactivate(sim)
					}

					dkSpellCost := deathKnight.DetermineOptimalCost(sim, 0, 1, 0)
					deathKnight.Spend(sim, spell, dkSpellCost)

					deathKnight.FrostFeverSpell.Cast(sim, spellEffect.Target)
					if deathKnight.Talents.CryptFever > 0 {
						deathKnight.CryptFeverAura[spellEffect.Target.Index].Activate(sim)
					}
					if deathKnight.Talents.EbonPlaguebringer > 0 {
						deathKnight.EbonPlagueAura[spellEffect.Target.Index].Activate(sim)
					}

					amountOfRunicPower := 10.0 + 2.5*float64(deathKnight.Talents.ChillOfTheGrave)
					deathKnight.AddRunicPower(sim, amountOfRunicPower, spell.RunicPowerMetrics())

					deathKnight.IcyTouchAura[spellEffect.Target.Index].Activate(sim)

					// In reality if you have the talent just casting IT
					// activates the aura, no need to check for enemy debuff
					if deathKnight.IcyTalonsAura != nil {
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

func (deathKnight *DeathKnight) CastIcyTouch(sim *core.Simulation, target *core.Unit) bool {
	if deathKnight.CanIcyTouch(sim) {
		deathKnight.IcyTouch.Cast(sim, target)
		return true
	}
	return false
}
