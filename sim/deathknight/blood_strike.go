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
		weaponBaseDamage = core.BaseDamageFuncMeleeWeapon(core.OffHand, false, 306.0, 0.4, true)
	}

	bloodOfTheNorthCoeff := 1.0
	if deathKnight.Talents.BloodOfTheNorth == 1 {
		bloodOfTheNorthCoeff = 1.03
	} else if deathKnight.Talents.BloodOfTheNorth == 2 {
		bloodOfTheNorthCoeff = 1.06
	} else if deathKnight.Talents.BloodOfTheNorth == 3 {
		bloodOfTheNorthCoeff = 1.1
	}
	guileOfGorefiend := deathKnight.Talents.GuileOfGorefiend > 0

	effect := core.SpellEffect{
		BonusCritRating:  (3.0*float64(deathKnight.Talents.Subversion) + 1.0*float64(deathKnight.Talents.Annihilation)) * core.CritRatingPerCritChance,
		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		BaseDamage: core.BaseDamageConfig{
			Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
				return weaponBaseDamage(sim, hitEffect, spell) *
					bloodOfTheNorthCoeff *
					(1.0 + float64(deathKnight.countActiveDiseases())*0.125) *
					core.TernaryFloat64(deathKnight.BloodPlagueDisease.IsActive(), 1.0+0.02*float64(deathKnight.Talents.RageOfRivendare), 1.0) *
					core.TernaryFloat64(deathKnight.DiseasesAreActive(), 1.0+0.05*float64(deathKnight.Talents.TundraStalker), 1.0)
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

	deathKnight.threatOfThassarianProcMasks(isMH, &effect, guileOfGorefiend)

	return deathKnight.RegisterSpell(core.SpellConfig{
		ActionID:     BloodStrikeActionID,
		SpellSchool:  core.SpellSchoolPhysical,
		Flags:        core.SpellFlagMeleeMetrics,
		ApplyEffects: core.ApplyEffectFuncDirectDamage(effect),
	})
}

func (deathKnight *DeathKnight) registerBloodStrikeSpell() {
	mhHitSpell := deathKnight.newBloodStrikeSpell(true)
	ohHitSpell := deathKnight.newBloodStrikeSpell(false)

	deathKnight.BloodStrike = deathKnight.RegisterSpell(core.SpellConfig{
		ActionID:    BloodStrikeActionID,
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
				deathKnight.threatOfThassarianProc(sim, spellEffect, mhHitSpell, ohHitSpell)
				deathKnight.threatOfThassarianAdjustMetrics(sim, spell, spellEffect, BloodStrikeMHOutcome)

				if deathKnight.outcomeEitherWeaponHitOrCrit(BloodStrikeMHOutcome, BloodStrikeOHOutcome) {
					dkSpellCost := deathKnight.DetermineOptimalCost(sim, 1, 0, 0)
					deathKnight.bloodOfTheNorthProc(sim, spell, dkSpellCost)

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
