package deathknight

import (
	"github.com/wowsims/wotlk/sim/core"
)

var PlagueStrikeActionID = core.ActionID{SpellID: 49921}
var PlagueStrikeMHOutcome = core.OutcomeHit
var PlagueStrikeOHOutcome = core.OutcomeHit

func (deathKnight *DeathKnight) newPlagueStrikeSpell(isMH bool) *core.Spell {
	weaponBaseDamage := core.BaseDamageFuncMeleeWeapon(core.MainHand, false, 189.0, 0.5, true)
	if !isMH {
		weaponBaseDamage = core.BaseDamageFuncMeleeWeapon(core.OffHand, false, 189.0, 0.5*deathKnight.nervesOfColdSteelBonus(), true)
	}

	outbreakBonus := 1.0 + 0.1*float64(deathKnight.Talents.Outbreak)

	effect := core.SpellEffect{
		BonusCritRating:  (1.0*float64(deathKnight.Talents.Annihilation) + 3.0*float64(deathKnight.Talents.ViciousStrikes)) * core.CritRatingPerCritChance,
		DamageMultiplier: outbreakBonus,
		ThreatMultiplier: 1,

		BaseDamage: core.BaseDamageConfig{
			Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
				return weaponBaseDamage(sim, hitEffect, spell) *
					deathKnight.rageOfRivendareBonus(hitEffect.Target) *
					deathKnight.tundraStalkerBonus(hitEffect.Target)
			},
			TargetSpellCoefficient: 1,
		},

		OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if isMH {
				PlagueStrikeMHOutcome = spellEffect.Outcome
			} else {
				PlagueStrikeOHOutcome = spellEffect.Outcome
			}
		},
	}

	if isMH {
		effect.ProcMask = core.ProcMaskMeleeMHSpecial
		effect.OutcomeApplier = deathKnight.OutcomeFuncMeleeSpecialHitAndCrit(deathKnight.MeleeCritMultiplier(1.0, deathKnight.viciousStrikesBonus()))
	} else {
		effect.ProcMask = core.ProcMaskMeleeOHSpecial
		effect.OutcomeApplier = deathKnight.OutcomeFuncMeleeSpecialCritOnly(deathKnight.MeleeCritMultiplier(1.0, deathKnight.viciousStrikesBonus()))
	}

	return deathKnight.RegisterSpell(core.SpellConfig{
		ActionID:     PlagueStrikeActionID,
		SpellSchool:  core.SpellSchoolPhysical,
		Flags:        core.SpellFlagMeleeMetrics,
		ApplyEffects: core.ApplyEffectFuncDirectDamage(effect),
	})
}

func (deathKnight *DeathKnight) registerPlagueStrikeSpell() {
	deathKnight.PlagueStrikeMhHit = deathKnight.newPlagueStrikeSpell(true)
	deathKnight.PlagueStrikeOhHit = deathKnight.newPlagueStrikeSpell(false)

	deathKnight.PlagueStrike = deathKnight.RegisterSpell(core.SpellConfig{
		ActionID:    PlagueStrikeActionID,
		SpellSchool: core.SpellSchoolPhysical,
		Flags:       core.SpellFlagMeleeMetrics,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				cast.GCD = deathKnight.getModifiedGCD()
			},
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:         core.ProcMaskMeleeMHSpecial,
			ThreatMultiplier: 1,

			OutcomeApplier: deathKnight.OutcomeFuncAlwaysHit(),

			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				deathKnight.threatOfThassarianProc(sim, spellEffect, deathKnight.PlagueStrikeMhHit, deathKnight.PlagueStrikeOhHit)
				deathKnight.threatOfThassarianAdjustMetrics(sim, spell, spellEffect, PlagueStrikeMHOutcome)

				if deathKnight.outcomeEitherWeaponHitOrCrit(PlagueStrikeMHOutcome, PlagueStrikeOHOutcome) {
					dkSpellCost := deathKnight.DetermineOptimalCost(sim, 0, 0, 1)
					deathKnight.Spend(sim, spell, dkSpellCost)

					deathKnight.BloodPlagueSpell.Cast(sim, spellEffect.Target)
					if deathKnight.Talents.EbonPlaguebringer > 0 {
						deathKnight.EbonPlagueAura.Activate(sim)
					}

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
