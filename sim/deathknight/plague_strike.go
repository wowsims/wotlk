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
		BonusCritRating:  (deathKnight.annihilationCritBonus() + deathKnight.scourgebornePlateCritBonus() + deathKnight.viciousStrikesCritChanceBonus()) * core.CritRatingPerCritChance,
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

	deathKnight.threatOfThassarianProcMasks(isMH, &effect, false, func(outcomeApplier core.OutcomeApplier) core.OutcomeApplier {
		return outcomeApplier
	})

	return deathKnight.RegisterSpell(core.SpellConfig{
		ActionID:     PlagueStrikeActionID.WithTag(core.TernaryInt32(isMH, 1, 2)),
		SpellSchool:  core.SpellSchoolPhysical,
		Flags:        core.SpellFlagMeleeMetrics,
		ApplyEffects: core.ApplyEffectFuncDirectDamage(effect),
	})
}

func (deathKnight *DeathKnight) registerPlagueStrikeSpell() {
	deathKnight.PlagueStrikeMhHit = deathKnight.newPlagueStrikeSpell(true)
	deathKnight.PlagueStrikeOhHit = deathKnight.newPlagueStrikeSpell(false)

	deathKnight.PlagueStrike = deathKnight.RegisterSpell(core.SpellConfig{
		ActionID:    PlagueStrikeActionID.WithTag(3),
		SpellSchool: core.SpellSchoolPhysical,
		Flags:       core.SpellFlagNoMetrics | core.SpellFlagNoLogs,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				cast.GCD = deathKnight.getModifiedGCD()
			},
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:         core.ProcMaskEmpty,
			ThreatMultiplier: 1,

			OutcomeApplier: deathKnight.OutcomeFuncAlwaysHit(),

			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				deathKnight.threatOfThassarianProc(sim, spellEffect, deathKnight.PlagueStrikeMhHit, deathKnight.PlagueStrikeOhHit)

				deathKnight.LastCastOutcome = PlagueStrikeMHOutcome
				if deathKnight.outcomeEitherWeaponHitOrCrit(PlagueStrikeMHOutcome, PlagueStrikeOHOutcome) {
					dkSpellCost := deathKnight.DetermineOptimalCost(sim, 0, 0, 1)
					deathKnight.Spend(sim, spell, dkSpellCost)

					deathKnight.BloodPlagueSpell.Cast(sim, spellEffect.Target)
					if deathKnight.Talents.CryptFever > 0 {
						deathKnight.CryptFeverAura[spellEffect.Target.Index].Activate(sim)
					}
					if deathKnight.Talents.EbonPlaguebringer > 0 {
						deathKnight.EbonPlagueAura[spellEffect.Target.Index].Activate(sim)
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

func (deathKnight *DeathKnight) CastPlagueStrike(sim *core.Simulation, target *core.Unit) bool {
	if deathKnight.CanPlagueStrike(sim) {
		deathKnight.PlagueStrike.Cast(sim, target)
		return true
	}
	return false
}
