package deathknight

import (
	"github.com/wowsims/wotlk/sim/core"
)

func (deathKnight *DeathKnight) registerPlagueStrikeSpell() {
	weaponBaseDamage := core.BaseDamageFuncMeleeWeapon(core.MainHand, true, 0.0, 0.5, true)
	viciousStrikes := 0.15 * float64(deathKnight.Talents.ViciousStrikes)

	deathKnight.PlagueStrike = deathKnight.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 49921},
		SpellSchool: core.SpellSchoolPhysical,
		Flags:       core.SpellFlagMeleeMetrics,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:         core.ProcMaskMeleeMHSpecial,
			BonusCritRating:  (1.0*float64(deathKnight.Talents.Annihilation) + 3.0*float64(deathKnight.Talents.ViciousStrikes)) * core.CritRatingPerCritChance,
			DamageMultiplier: 1,
			ThreatMultiplier: 1,

			BaseDamage: core.BaseDamageConfig{
				Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
					return weaponBaseDamage(sim, hitEffect, spell) *
						(1.0 +
							0.10*float64(deathKnight.Talents.Outbreak))
				},
				TargetSpellCoefficient: 1,
			},

			OutcomeApplier: deathKnight.OutcomeFuncMeleeSpecialHitAndCrit(deathKnight.MeleeCritMultiplier(1.0, viciousStrikes)),

			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spellEffect.Landed() {
					dkSpellCost := deathKnight.DetermineOptimalCost(sim, 0, 0, 1)
					deathKnight.Spend(sim, spell, dkSpellCost)

					deathKnight.BloodPlagueDisease.Apply(sim)

					// TODO: Temporary application of ebon plague until dot auras
					// properly run their events to control ebon plague
					deathKnight.checkForEbonPlague(sim)

					amountOfRunicPower := 10.0 + 2.5*float64(deathKnight.Talents.Dirge)
					deathKnight.AddRunicPower(sim, amountOfRunicPower, spell.RunicPowerMetrics())
				}
			},
		}),
	})
}

func (deathKnight *DeathKnight) CanPlagueStrike(sim *core.Simulation) bool {
	return deathKnight.CastCostPossible(sim, 0.0, 0, 0, 1) && deathKnight.PlagueStrike.IsReady(sim)
}
