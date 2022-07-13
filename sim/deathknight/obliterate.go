package deathknight

import (
	"github.com/wowsims/wotlk/sim/core"
)

// TODO: Add disease consumption
func (deathKnight *DeathKnight) registerObliterateSpell() {
	weaponBaseDamage := core.BaseDamageFuncMeleeWeapon(core.MainHand, true, 467.0, 0.8, true)

	guileOfGorefiend := deathKnight.Talents.GuileOfGorefiend > 0

	diseaseConsumptionChance := 100.0
	if deathKnight.Talents.Annihilation == 1 {
		diseaseConsumptionChance = 67.0
	} else if deathKnight.Talents.Annihilation == 2 {
		diseaseConsumptionChance = 34.0
	} else if deathKnight.Talents.Annihilation == 3 {
		diseaseConsumptionChance = 0.0
	}

	deathKnight.Obliterate = deathKnight.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 51425},
		SpellSchool: core.SpellSchoolPhysical,
		Flags:       core.SpellFlagMeleeMetrics,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:         core.ProcMaskMeleeMHSpecial,
			BonusCritRating:  (5.0*float64(deathKnight.Talents.Rime) + 3.0*float64(deathKnight.Talents.Subversion) + 1.0*float64(deathKnight.Talents.Annihilation)) * core.CritRatingPerCritChance,
			DamageMultiplier: 1,
			ThreatMultiplier: 1,

			BaseDamage: core.BaseDamageConfig{
				Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
					return weaponBaseDamage(sim, hitEffect, spell) *
						(1.0 +
							core.TernaryFloat64(deathKnight.FrostFeverDisease.IsActive(), 0.125, 0.0) +
							core.TernaryFloat64(deathKnight.BloodPlagueDisease.IsActive(), 0.125, 0.0) +
							core.TernaryFloat64(sim.IsExecutePhase35() && deathKnight.Talents.MercilessCombat > 0, 0.06*float64(deathKnight.Talents.MercilessCombat), 0.0))
				},
				TargetSpellCoefficient: 1,
			},

			OutcomeApplier: deathKnight.OutcomeFuncMeleeSpecialHitAndCrit(deathKnight.critMultiplier(guileOfGorefiend)),

			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spellEffect.Landed() {
					dkSpellCost := deathKnight.DetermineOptimalCost(sim, 0, 1, 1)
					deathKnight.Spend(sim, spell, dkSpellCost)

					amountOfRunicPower := 15.0 + 2.5*float64(deathKnight.Talents.ChillOfTheGrave)
					deathKnight.AddRunicPower(sim, amountOfRunicPower, spell.RunicPowerMetrics())

					if sim.RandomFloat("Annihilation") < diseaseConsumptionChance {
						deathKnight.FrostFeverDisease.Deactivate(sim)
						deathKnight.BloodPlagueDisease.Deactivate(sim)
					}
				}
			},
		}),
	})
}

func (deathKnight *DeathKnight) CanObliterate(sim *core.Simulation) bool {
	return deathKnight.CastCostPossible(sim, 0.0, 0, 1, 1) && deathKnight.Obliterate.IsReady(sim)
}
