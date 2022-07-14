package deathknight

import (
	"github.com/wowsims/wotlk/sim/core"
)

func (deathKnight *DeathKnight) newBloodStrikeSpell(isMH bool) *core.Spell {
	weaponBaseDamage := core.BaseDamageFuncMeleeWeapon(core.MainHand, true, 764.0, 0.4, true)
	if !isMH {
		weaponBaseDamage = core.BaseDamageFuncMeleeWeapon(core.OffHand, true, 764.0, 0.4, true)
	}

	guileOfGorefiend := deathKnight.Talents.GuileOfGorefiend > 0

	actionID := core.ActionID{SpellID: 49930}

	effect := core.SpellEffect{
		ProcMask:         core.ProcMaskMeleeMHSpecial,
		BonusCritRating:  (3.0*float64(deathKnight.Talents.Subversion) + 1.0*float64(deathKnight.Talents.Annihilation)) * core.CritRatingPerCritChance,
		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		BaseDamage: core.BaseDamageConfig{
			Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
				return weaponBaseDamage(sim, hitEffect, spell) *
					(1.0 +
						core.TernaryFloat64(deathKnight.FrostFeverDisease.IsActive(), 0.125, 0.0) +
						core.TernaryFloat64(deathKnight.BloodPlagueDisease.IsActive(), 0.125, 0.0) +
						core.TernaryFloat64(deathKnight.EbonPlagueAura.IsActive(), 0.125, 0.0))
			},
			TargetSpellCoefficient: 1,
		},

		OutcomeApplier: deathKnight.OutcomeFuncMeleeSpecialHitAndCrit(deathKnight.critMultiplier(guileOfGorefiend)),

		OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if spellEffect.Landed() {
				if isMH {
					dkSpellCost := deathKnight.DetermineOptimalCost(sim, 1, 0, 0)
					deathKnight.Spend(sim, spell, dkSpellCost)

					amountOfRunicPower := 10.0
					deathKnight.AddRunicPower(sim, amountOfRunicPower, spell.RunicPowerMetrics())
				}

				if deathKnight.DesolationAura != nil {
					deathKnight.DesolationAura.Activate(sim)
				}
			}
		},
	}

	if !isMH {
		effect.ProcMask = core.ProcMaskMeleeOHSpecial
	}

	return deathKnight.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolPhysical,
		Flags:       core.SpellFlagMeleeMetrics,

		ApplyEffects: core.ApplyEffectFuncDirectDamage(effect),
	})
}

func (deathKnight *DeathKnight) registerBloodStrikeSpell() {
	mhHitSpell := deathKnight.newBloodStrikeSpell(true)
	ohHitSpell := deathKnight.newBloodStrikeSpell(false)

	threatOfThassarianChance := 0.0
	if deathKnight.Talents.ThreatOfThassarian == 1 {
		threatOfThassarianChance = 0.30
	} else if deathKnight.Talents.ThreatOfThassarian == 2 {
		threatOfThassarianChance = 0.60
	} else if deathKnight.Talents.ThreatOfThassarian == 3 {
		threatOfThassarianChance = 1.0
	}

	deathKnight.BloodStrike = deathKnight.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 49930},
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
			OutcomeApplier:   deathKnight.OutcomeFuncMeleeSpecialHit(),
			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spellEffect.Landed() {
					mhHitSpell.Cast(sim, spellEffect.Target)
					if sim.RandomFloat("Threat of Thassarian") < threatOfThassarianChance {
						ohHitSpell.Cast(sim, spellEffect.Target)

						deathKnight.BloodStrike.SpellMetrics[spellEffect.Target.TableIndex].Casts -= 2
						deathKnight.BloodStrike.SpellMetrics[spellEffect.Target.TableIndex].Hits--
					} else {
						deathKnight.BloodStrike.SpellMetrics[spellEffect.Target.TableIndex].Casts -= 1
					}
				}
			},
		}),
	})
}

func (deathKnight *DeathKnight) CanBloodStrike(sim *core.Simulation) bool {
	return deathKnight.CastCostPossible(sim, 0.0, 1, 0, 0) && deathKnight.BloodStrike.IsReady(sim)
}
