package deathknight

import (
	"github.com/wowsims/wotlk/sim/core"
)

var BloodStrikeActionID = core.ActionID{SpellID: 49930}
var BloodStrikeMHOutcome = core.OutcomeHit
var BloodStrikeOHOutcome = core.OutcomeHit

func (deathKnight *DeathKnight) newBloodStrikeSpell(isMH bool) *core.Spell {
	weaponBaseDamage := core.BaseDamageFuncMeleeWeapon(core.MainHand, false, 306.0, 0.4, true)
	if !isMH {
		weaponBaseDamage = core.BaseDamageFuncMeleeWeapon(core.OffHand, false, 306.0, 0.4*deathKnight.nervesOfColdSteelBonus(), true)
	}

	effect := core.SpellEffect{
		BonusCritRating:  (3.0*float64(deathKnight.Talents.Subversion) + deathKnight.annihilationCritBonus()) * core.CritRatingPerCritChance,
		DamageMultiplier: deathKnight.bloodOfTheNorthCoeff(),
		ThreatMultiplier: 1,

		BaseDamage: core.BaseDamageConfig{
			Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
				return weaponBaseDamage(sim, hitEffect, spell) *
					deathKnight.diseaseMultiplierBonus(hitEffect.Target, 0.125) *
					deathKnight.rageOfRivendareBonus(hitEffect.Target) *
					deathKnight.tundraStalkerBonus(hitEffect.Target)
			},
			TargetSpellCoefficient: 1,
		},

		OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if isMH {
				BloodStrikeMHOutcome = spellEffect.Outcome
			} else {
				BloodStrikeOHOutcome = spellEffect.Outcome
			}
		},
	}

	deathKnight.threatOfThassarianProcMasks(isMH, &effect, true, func(outcomeApplier core.OutcomeApplier) core.OutcomeApplier {
		return outcomeApplier
	})

	return deathKnight.RegisterSpell(core.SpellConfig{
		ActionID:     BloodStrikeActionID.WithTag(core.TernaryInt32(isMH, 1, 2)),
		SpellSchool:  core.SpellSchoolPhysical,
		Flags:        core.SpellFlagMeleeMetrics,
		ApplyEffects: core.ApplyEffectFuncDirectDamage(effect),
	})
}

func (deathKnight *DeathKnight) registerBloodStrikeSpell() {
	deathKnight.BloodStrikeMhHit = deathKnight.newBloodStrikeSpell(true)
	deathKnight.BloodStrikeOhHit = deathKnight.newBloodStrikeSpell(false)

	deathKnight.BloodStrike = deathKnight.RegisterSpell(core.SpellConfig{
		ActionID:    BloodStrikeActionID.WithTag(3),
		Flags:       core.SpellFlagNoMetrics | core.SpellFlagNoLogs,
		SpellSchool: core.SpellSchoolPhysical,

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
				deathKnight.threatOfThassarianProc(sim, spellEffect, deathKnight.BloodStrikeMhHit, deathKnight.BloodStrikeOhHit)

				deathKnight.LastCastOutcome = BloodStrikeMHOutcome

				if deathKnight.outcomeEitherWeaponHitOrCrit(BloodStrikeMHOutcome, BloodStrikeOHOutcome) {
					dkSpellCost := deathKnight.DetermineOptimalCost(sim, 1, 0, 0)
					if !deathKnight.bloodOfTheNorthProc(sim, spell, dkSpellCost) {
						if !deathKnight.reapingProc(sim, spell, dkSpellCost) {
							deathKnight.Spend(sim, spell, dkSpellCost)
						}
					}

					amountOfRunicPower := 10.0
					deathKnight.AddRunicPower(sim, amountOfRunicPower, spell.RunicPowerMetrics())

					if deathKnight.DesolationAura != nil {
						deathKnight.DesolationAura.Activate(sim)
					}
				}
			},
		}),
	})
}

func (deathKnight *DeathKnight) CanBloodStrike(sim *core.Simulation) bool {
	return deathKnight.CastCostPossible(sim, 0.0, 1, 0, 0) && deathKnight.BloodStrike.IsReady(sim)
}
