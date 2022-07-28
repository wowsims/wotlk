package deathknight

import (
	"github.com/wowsims/wotlk/sim/core"
)

var BloodStrikeActionID = core.ActionID{SpellID: 49930}
var BloodStrikeMHOutcome = core.OutcomeMiss
var BloodStrikeOHOutcome = core.OutcomeMiss

func (dk *Deathknight) newBloodStrikeSpell(isMH bool) *core.Spell {
	weaponBaseDamage := core.BaseDamageFuncMeleeWeapon(core.MainHand, true, 764.0, 0.4, true)
	if !isMH {
		weaponBaseDamage = core.BaseDamageFuncMeleeWeapon(core.OffHand, true, 764.0, 0.4*dk.nervesOfColdSteelBonus(), true)
	}

	effect := core.SpellEffect{
		BonusCritRating:  (dk.subversionCritBonus() + dk.annihilationCritBonus()) * core.CritRatingPerCritChance,
		DamageMultiplier: dk.bloodOfTheNorthCoeff() * dk.thassariansPlateDamageBonus(),
		ThreatMultiplier: 1,

		BaseDamage: core.BaseDamageConfig{
			Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
				return weaponBaseDamage(sim, hitEffect, spell) *
					dk.diseaseMultiplierBonus(hitEffect.Target, 0.125) *
					dk.rageOfRivendareBonus(hitEffect.Target) *
					dk.tundraStalkerBonus(hitEffect.Target)
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

	dk.threatOfThassarianProcMasks(isMH, &effect, true, true, func(outcomeApplier core.OutcomeApplier) core.OutcomeApplier {
		return outcomeApplier
	})

	return dk.RegisterSpell(core.SpellConfig{
		ActionID:     BloodStrikeActionID.WithTag(core.TernaryInt32(isMH, 1, 2)),
		SpellSchool:  core.SpellSchoolPhysical,
		Flags:        core.SpellFlagMeleeMetrics,
		ApplyEffects: core.ApplyEffectFuncDirectDamage(effect),
	})
}

func (dk *Deathknight) registerBloodStrikeSpell() {
	dk.BloodStrikeMhHit = dk.newBloodStrikeSpell(true)
	dk.BloodStrikeOhHit = dk.newBloodStrikeSpell(false)

	dk.BloodStrike = dk.RegisterSpell(core.SpellConfig{
		ActionID:    BloodStrikeActionID.WithTag(3),
		Flags:       core.SpellFlagNoMetrics | core.SpellFlagNoLogs,
		SpellSchool: core.SpellSchoolPhysical,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				cast.GCD = dk.getModifiedGCD()
			},
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:         core.ProcMaskEmpty,
			ThreatMultiplier: 1,

			OutcomeApplier: dk.OutcomeFuncAlwaysHit(),

			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				dk.threatOfThassarianProc(sim, spellEffect, dk.BloodStrikeMhHit, dk.BloodStrikeOhHit)

				dk.LastCastOutcome = BloodStrikeMHOutcome

				if dk.outcomeEitherWeaponLanded(BloodStrikeMHOutcome, BloodStrikeOHOutcome) {
					dkSpellCost := dk.DetermineCost(sim, core.DKCastEnum_B)
					if !dk.bloodOfTheNorthProc(sim, spell, dkSpellCost) {
						if !dk.reapingProc(sim, spell, dkSpellCost) {
							dk.Spend(sim, spell, dkSpellCost)
						}
					}

					if dk.DesolationAura != nil {
						dk.DesolationAura.Activate(sim)
					}

					// Gain at the end, to take into account previous effects for callback
					amountOfRunicPower := 10.0
					dk.AddRunicPower(sim, amountOfRunicPower, spell.RunicPowerMetrics())
				}
			},
		}),
	})
}

func (dk *Deathknight) CanBloodStrike(sim *core.Simulation) bool {
	return dk.CastCostPossible(sim, 0.0, 1, 0, 0) && dk.BloodStrike.IsReady(sim)
}

func (dk *Deathknight) CastBloodStrike(sim *core.Simulation, target *core.Unit) bool {
	if dk.CanBloodStrike(sim) {
		dk.BloodStrike.Cast(sim, target)
		return true
	}
	return false
}
