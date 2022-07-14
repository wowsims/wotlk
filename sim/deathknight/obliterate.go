package deathknight

import (
	"github.com/wowsims/wotlk/sim/core"
)

var ObliterateActionID = core.ActionID{SpellID: 51425}
var ObliterateMHOutcome = core.OutcomeHit
var ObliterateOHOutcome = core.OutcomeHit

func (deathKnight *DeathKnight) newObliterateHitSpell(isMH bool) *core.Spell {
	weaponBaseDamage := core.BaseDamageFuncMeleeWeapon(core.MainHand, false, 467.0, 0.8, true)
	if !isMH {
		weaponBaseDamage = core.BaseDamageFuncMeleeWeapon(core.OffHand, false, 467.0, 0.8, true)
	}

	guileOfGorefiend := deathKnight.Talents.GuileOfGorefiend > 0

	diseaseConsumptionChance := 1.0
	if deathKnight.Talents.Annihilation == 1 {
		diseaseConsumptionChance = 0.67
	} else if deathKnight.Talents.Annihilation == 2 {
		diseaseConsumptionChance = 0.34
	} else if deathKnight.Talents.Annihilation == 3 {
		diseaseConsumptionChance = 0.0
	}

	hbResetCDChance := 0.05 * float64(deathKnight.Talents.Rime)

	effect := core.SpellEffect{
		BonusCritRating:  (5.0*float64(deathKnight.Talents.Rime) + 3.0*float64(deathKnight.Talents.Subversion) + 1.0*float64(deathKnight.Talents.Annihilation)) * core.CritRatingPerCritChance,
		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		BaseDamage: core.BaseDamageConfig{
			Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
				return weaponBaseDamage(sim, hitEffect, spell) *
					(1.0 +
						float64(deathKnight.countActiveDiseases())*0.125 +
						core.TernaryFloat64(deathKnight.DiseasesAreActive(), 0.05*float64(deathKnight.Talents.TundraStalker), 0.0) +
						core.TernaryFloat64(deathKnight.BloodPlagueDisease.IsActive(), 0.02*float64(deathKnight.Talents.RageOfRivendare), 0.0) +
						core.TernaryFloat64(sim.IsExecutePhase35() && deathKnight.Talents.MercilessCombat > 0, 0.06*float64(deathKnight.Talents.MercilessCombat), 0.0))
			},
			TargetSpellCoefficient: 1,
		},

		OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if isMH {
				ObliterateMHOutcome = spellEffect.Outcome
			} else {
				ObliterateOHOutcome = spellEffect.Outcome
			}

			if sim.RandomFloat("Annihilation") < diseaseConsumptionChance {
				deathKnight.FrostFeverDisease.Deactivate(sim)
				deathKnight.BloodPlagueDisease.Deactivate(sim)
			}

			if sim.RandomFloat("Rime") < hbResetCDChance {
				deathKnight.HowlingBlast.CD.Reset()
				deathKnight.HowlingBlastCostless = true
			}
		},
	}

	deathKnight.threatOfThassarianProcMasks(isMH, &effect, guileOfGorefiend)

	return deathKnight.RegisterSpell(core.SpellConfig{
		ActionID:     ObliterateActionID,
		SpellSchool:  core.SpellSchoolPhysical,
		Flags:        core.SpellFlagMeleeMetrics,
		ApplyEffects: core.ApplyEffectFuncDirectDamage(effect),
	})
}

func (deathKnight *DeathKnight) registerObliterateSpell() {
	mhHitSpell := deathKnight.newObliterateHitSpell(true)
	ohHitSpell := deathKnight.newObliterateHitSpell(false)

	threatOfThassarianChance := 0.0
	if deathKnight.Talents.ThreatOfThassarian == 1 {
		threatOfThassarianChance = 0.30
	} else if deathKnight.Talents.ThreatOfThassarian == 2 {
		threatOfThassarianChance = 0.60
	} else if deathKnight.Talents.ThreatOfThassarian == 3 {
		threatOfThassarianChance = 1.0
	}

	deathKnight.Obliterate = deathKnight.RegisterSpell(core.SpellConfig{
		ActionID:    ObliterateActionID,
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
				deathKnight.Obliterate.SpellMetrics[spellEffect.Target.TableIndex].Casts -= 1
				deathKnight.Obliterate.SpellMetrics[spellEffect.Target.TableIndex].Hits -= 1
				if sim.RandomFloat("Threat of Thassarian") < threatOfThassarianChance {
					ohHitSpell.Cast(sim, spellEffect.Target)
					deathKnight.Obliterate.SpellMetrics[spellEffect.Target.TableIndex].Casts -= 1
				}

				if deathKnight.threatOfThassarianHitCheck(ObliterateMHOutcome, ObliterateOHOutcome) {
					dkSpellCost := deathKnight.DetermineOptimalCost(sim, 0, 1, 1)
					deathKnight.Spend(sim, spell, dkSpellCost)

					amountOfRunicPower := 15.0 + 2.5*float64(deathKnight.Talents.ChillOfTheGrave)
					deathKnight.AddRunicPower(sim, amountOfRunicPower, spell.RunicPowerMetrics())
				}
			},
		}),
	})
}

func (deathKnight *DeathKnight) CanObliterate(sim *core.Simulation) bool {
	return deathKnight.CastCostPossible(sim, 0.0, 0, 1, 1) && deathKnight.Obliterate.IsReady(sim)
}
