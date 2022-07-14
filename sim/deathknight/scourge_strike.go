package deathknight

import (
	"github.com/wowsims/wotlk/sim/core"
)

func (deathKnight *DeathKnight) registerScourgeStrikeShadowDamageSpell() *core.Spell {
	actionID := core.ActionID{SpellID: 55270}
	return deathKnight.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolShadow,
		Flags:       core.SpellFlagIgnoreResists,

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:         core.ProcMaskSpellDamage,
			BonusCritRating:  -100 * core.CritRatingPerCritChance, // Disable criticals for shadow hit
			DamageMultiplier: 1,
			ThreatMultiplier: 1,

			OutcomeApplier: deathKnight.CurrentTarget.OutcomeFuncAlwaysHit(),

			BaseDamage: core.BaseDamageConfig{
				Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
					return (deathKnight.LastScourgeStrikeDamage *
						(core.TernaryFloat64(deathKnight.FrostFeverDisease.IsActive(), 0.12, 0.0) +
							core.TernaryFloat64(deathKnight.BloodPlagueDisease.IsActive(), 0.12, 0.0) +
							core.TernaryFloat64(deathKnight.EbonPlagueAura.IsActive(), 0.12, 0.0)))
				},
			},
		}),
	})
}

func (deathKnight *DeathKnight) registerScourgeStrikeSpell() {
	weaponBaseDamage := core.BaseDamageFuncMeleeWeapon(core.MainHand, true, 560.0, 0.7, true)
	viciousStrikes := 0.15 * float64(deathKnight.Talents.ViciousStrikes)
	actionID := core.ActionID{SpellID: 55271}

	outbreakBonus := 0.0
	if deathKnight.Talents.Outbreak == 1 {
		outbreakBonus = 0.07
	} else if deathKnight.Talents.Outbreak == 2 {
		outbreakBonus = 0.13
	} else if deathKnight.Talents.Outbreak == 3 {
		outbreakBonus = 0.20
	}

	shadowDamageSpell := deathKnight.registerScourgeStrikeShadowDamageSpell()

	deathKnight.ScourgeStrike = deathKnight.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolPhysical,
		Flags:       core.SpellFlagMeleeMetrics,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:         core.ProcMaskMeleeMHSpecial,
			BonusCritRating:  (3.0*float64(deathKnight.Talents.Subversion) + 3.0*float64(deathKnight.Talents.ViciousStrikes)) * core.CritRatingPerCritChance,
			DamageMultiplier: 1,
			ThreatMultiplier: 1,

			BaseDamage: core.BaseDamageConfig{
				Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
					return weaponBaseDamage(sim, hitEffect, spell) * (1.0 + outbreakBonus)
				},
				TargetSpellCoefficient: 1,
			},

			OutcomeApplier: deathKnight.OutcomeFuncMeleeSpecialHitAndCrit(deathKnight.MeleeCritMultiplier(1.0, viciousStrikes)),

			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spellEffect.Landed() {
					dkSpellCost := deathKnight.DetermineOptimalCost(sim, 0, 1, 1)
					deathKnight.Spend(sim, spell, dkSpellCost)

					amountOfRunicPower := 15.0 + 2.5*float64(deathKnight.Talents.Dirge)
					deathKnight.AddRunicPower(sim, amountOfRunicPower, spell.RunicPowerMetrics())

					if deathKnight.DiseasesAreActive() {
						deathKnight.LastScourgeStrikeDamage = spellEffect.Damage
						shadowDamageSpell.Cast(sim, spellEffect.Target)
					}
				}
			},
		}),
	})
}

func (deathKnight *DeathKnight) CanScourgeStrike(sim *core.Simulation) bool {
	return deathKnight.Talents.ScourgeStrike && deathKnight.CastCostPossible(sim, 0.0, 0, 1, 1) && deathKnight.ScourgeStrike.IsReady(sim)
}
