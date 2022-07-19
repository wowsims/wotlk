package deathknight

import (
	"github.com/wowsims/wotlk/sim/core"
)

func (deathKnight *DeathKnight) registerScourgeStrikeShadowDamageSpell() *core.Spell {
	actionID := core.ActionID{SpellID: 55271}

	return deathKnight.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolShadow,
		Flags:       core.SpellFlagIgnoreResists | core.SpellFlagMeleeMetrics,

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:         core.ProcMaskSpellDamage,
			BonusCritRating:  -100 * core.CritRatingPerCritChance, // Disable criticals for shadow hit
			DamageMultiplier: 1,
			ThreatMultiplier: 1,

			OutcomeApplier: deathKnight.CurrentTarget.OutcomeFuncAlwaysHit(),

			BaseDamage: core.BaseDamageConfig{
				Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
					return deathKnight.LastScourgeStrikeDamage * (float64(deathKnight.countActiveDiseases(hitEffect.Target)) * 0.12)
				},
			},
		}),
	})
}

func (deathKnight *DeathKnight) registerScourgeStrikeSpell() {
	actionID := core.ActionID{SpellID: 55271}

	outbreakBonus := 1.0
	if deathKnight.Talents.Outbreak == 1 {
		outbreakBonus = 1.07
	} else if deathKnight.Talents.Outbreak == 2 {
		outbreakBonus = 1.13
	} else if deathKnight.Talents.Outbreak == 3 {
		outbreakBonus = 1.20
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
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				cast.GCD = deathKnight.getModifiedGCD()
			},
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:         core.ProcMaskMeleeMHSpecial,
			BonusCritRating:  (3.0*float64(deathKnight.Talents.Subversion) + 3.0*float64(deathKnight.Talents.ViciousStrikes) + deathKnight.scourgeborneCritBonus()) * core.CritRatingPerCritChance,
			DamageMultiplier: outbreakBonus,
			ThreatMultiplier: 1,

			BaseDamage: core.BaseDamageConfig{
				Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
					bonusBaseDamage := deathKnight.sigilOfAwarenessBonus(Ability_ScourgeStrike)
					bonusBaseDamage += deathKnight.sigilOfArthriticBindingBonus()
					weaponBaseDamage := core.BaseDamageFuncMeleeWeapon(core.MainHand, false, 560.0+bonusBaseDamage, 0.7, true)
					return weaponBaseDamage(sim, hitEffect, spell) *
						deathKnight.rageOfRivendareBonus(hitEffect.Target) *
						deathKnight.tundraStalkerBonus(hitEffect.Target)
				},
				TargetSpellCoefficient: 1,
			},

			OutcomeApplier: deathKnight.OutcomeFuncMeleeSpecialHitAndCrit(deathKnight.MeleeCritMultiplier(1.0, deathKnight.viciousStrikesBonus())),

			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spellEffect.Landed() {
					dkSpellCost := deathKnight.DetermineOptimalCost(sim, 0, 1, 1)
					deathKnight.Spend(sim, spell, dkSpellCost)

					amountOfRunicPower := 15.0 + 2.5*float64(deathKnight.Talents.Dirge) + deathKnight.scourgeborneRunicPowerBonus()
					deathKnight.AddRunicPower(sim, amountOfRunicPower, spell.RunicPowerMetrics())

					if deathKnight.DiseasesAreActive(spellEffect.Target) {
						deathKnight.LastScourgeStrikeDamage = spellEffect.Damage
						shadowDamageSpell.Cast(sim, spellEffect.Target)
						deathKnight.ScourgeStrike.SpellMetrics[spellEffect.Target.TableIndex].Casts -= 1
						deathKnight.ScourgeStrike.SpellMetrics[spellEffect.Target.TableIndex].Hits -= 1
					}
				}
			},
		}),
	})
}

func (deathKnight *DeathKnight) CanScourgeStrike(sim *core.Simulation) bool {
	return deathKnight.Talents.ScourgeStrike && deathKnight.CastCostPossible(sim, 0.0, 0, 1, 1) && deathKnight.ScourgeStrike.IsReady(sim)
}
