package deathknight

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (deathKnight *DeathKnight) registerBloodStrikeSpell() {
	baseCost := 10.0
	weaponBaseDamage := core.BaseDamageFuncMeleeWeapon(core.MainHand, true, 764.0, 0.4, true)

	guileOfGorefiend := deathKnight.Talents.GuileOfGorefiend > 0

	deathKnight.BloodStrike = deathKnight.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 49930},
		SpellSchool: core.SpellSchoolPhysical,
		Flags:       core.SpellFlagMeleeMetrics,

		ResourceType: stats.RunicPower,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost,
				GCD:  core.GCDDefault,
			},
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:         core.ProcMaskMeleeMHSpecial,
			BonusCritRating:  3.0*float64(deathKnight.Talents.Subversion)*core.CritRatingPerCritChance + 3.0*float64(deathKnight.Talents.Annihilation),
			DamageMultiplier: 1,
			ThreatMultiplier: 1,

			BaseDamage: core.BaseDamageConfig{
				Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
					return weaponBaseDamage(sim, hitEffect, spell) *
						(1.0 +
							core.TernaryFloat64(deathKnight.FrostFeverDisease.IsActive(), 0.125, 0.0) +
							core.TernaryFloat64(deathKnight.BloodPlagueDisease.IsActive(), 0.125, 0.0))
				},
				TargetSpellCoefficient: 1,
			},

			OutcomeApplier: deathKnight.OutcomeFuncMeleeSpecialHitAndCrit(deathKnight.critMultiplier(guileOfGorefiend)),

			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spellEffect.Landed() {
					dkSpellCost := deathKnight.DetermineOptimalCost(sim, baseCost, 1, 0, 0)
					deathKnight.Spend(sim, spell, dkSpellCost)

					amountOfRunicPower := 10.0
					deathKnight.AddRunicPower(sim, amountOfRunicPower, spell.RunicPowerMetrics())
				}
			},
		}),
	})
}

func (deathKnight *DeathKnight) CanBloodStrike(sim *core.Simulation) bool {
	return deathKnight.CastCostPossible(sim, 10.0, 1, 0, 0) && deathKnight.BloodStrike.IsReady(sim)
}
