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
		weaponBaseDamage = core.BaseDamageFuncMeleeWeapon(core.OffHand, false, 189.0, 0.5, true)
	}

	viciousStrikes := 0.15 * float64(deathKnight.Talents.ViciousStrikes)

	effect := core.SpellEffect{
		BonusCritRating:  (1.0*float64(deathKnight.Talents.Annihilation) + 3.0*float64(deathKnight.Talents.ViciousStrikes)) * core.CritRatingPerCritChance,
		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		BaseDamage: core.BaseDamageConfig{
			Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
				return weaponBaseDamage(sim, hitEffect, spell) *
					(1.0 +
						core.TernaryFloat64(deathKnight.DiseasesAreActive(), 0.05*float64(deathKnight.Talents.TundraStalker), 0.0) +
						core.TernaryFloat64(deathKnight.BloodPlagueDisease.IsActive(), 0.02*float64(deathKnight.Talents.RageOfRivendare), 0.0) +
						0.10*float64(deathKnight.Talents.Outbreak))
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
		effect.OutcomeApplier = deathKnight.OutcomeFuncMeleeSpecialHitAndCrit(deathKnight.MeleeCritMultiplier(1.0, viciousStrikes))
	} else {
		effect.ProcMask = core.ProcMaskMeleeOHSpecial
		effect.OutcomeApplier = deathKnight.OutcomeFuncMeleeSpecialCritOnly(deathKnight.MeleeCritMultiplier(1.0, viciousStrikes))
	}

	return deathKnight.RegisterSpell(core.SpellConfig{
		ActionID:     PlagueStrikeActionID,
		SpellSchool:  core.SpellSchoolPhysical,
		Flags:        core.SpellFlagMeleeMetrics,
		ApplyEffects: core.ApplyEffectFuncDirectDamage(effect),
	})
}

func (deathKnight *DeathKnight) registerPlagueStrikeSpell() {
	mhHitSpell := deathKnight.newPlagueStrikeSpell(true)
	ohHitSpell := deathKnight.newPlagueStrikeSpell(false)

	totChance := ToTChance(deathKnight)

	deathKnight.PlagueStrike = deathKnight.RegisterSpell(core.SpellConfig{
		ActionID:    PlagueStrikeActionID,
		SpellSchool: core.SpellSchoolPhysical,
		Flags:       core.SpellFlagMeleeMetrics,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:         core.ProcMaskMeleeMHSpecial,
			ThreatMultiplier: 1,

			OutcomeApplier: deathKnight.OutcomeFuncAlwaysHit(),

			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				mhHitSpell.Cast(sim, spellEffect.Target)
				totProcced := ToTWillCast(sim, totChance)
				if totProcced {
					ohHitSpell.Cast(sim, spellEffect.Target)
				}

				ToTAdjustMetrics(sim, spell, spellEffect, PlagueStrikeMHOutcome)

				if OutcomeEitherWeaponHitOrCrit(PlagueStrikeMHOutcome, PlagueStrikeOHOutcome) {
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
