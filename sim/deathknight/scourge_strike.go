package deathknight

import (
	"github.com/wowsims/wotlk/sim/core"
)

var ScourgeStrikeActionID = core.ActionID{SpellID: 55271}

func (dk *Deathknight) registerScourgeStrikeShadowDamageSpell() *core.Spell {
	return dk.RegisterSpell(core.SpellConfig{
		ActionID:    ScourgeStrikeActionID.WithTag(2),
		SpellSchool: core.SpellSchoolShadow,
		Flags:       core.SpellFlagIgnoreResists | core.SpellFlagMeleeMetrics,

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:         core.ProcMaskSpellDamage,
			BonusCritRating:  -100 * core.CritRatingPerCritChance, // Disable criticals for shadow hit
			DamageMultiplier: 1,
			ThreatMultiplier: 1,

			OutcomeApplier: dk.CurrentTarget.OutcomeFuncAlwaysHit(),

			BaseDamage: core.BaseDamageConfig{
				Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
					return dk.LastScourgeStrikeDamage * (dk.diseaseMultiplierBonus(hitEffect.Target, 0.12) - 1.0)
				},
			},
		}),
	})
}

func (dk *Deathknight) registerScourgeStrikeSpell() {
	outbreakBonus := 1.0
	if dk.Talents.Outbreak == 1 {
		outbreakBonus = 1.07
	} else if dk.Talents.Outbreak == 2 {
		outbreakBonus = 1.13
	} else if dk.Talents.Outbreak == 3 {
		outbreakBonus = 1.20
	}

	shadowDamageSpell := dk.registerScourgeStrikeShadowDamageSpell()

	dk.ScourgeStrike = dk.RegisterSpell(core.SpellConfig{
		ActionID:    ScourgeStrikeActionID.WithTag(1),
		SpellSchool: core.SpellSchoolPhysical,
		Flags:       core.SpellFlagMeleeMetrics,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				cast.GCD = dk.getModifiedGCD()
			},
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:         core.ProcMaskMeleeMHSpecial,
			BonusCritRating:  (dk.subversionCritBonus() + dk.viciousStrikesCritChanceBonus() + dk.scourgeborneBattlegearCritBonus()) * core.CritRatingPerCritChance,
			DamageMultiplier: outbreakBonus,
			ThreatMultiplier: 1,

			BaseDamage: core.BaseDamageConfig{
				Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
					bonusBaseDamage := dk.sigilOfAwarenessBonus(dk.ScourgeStrike)
					bonusBaseDamage += dk.sigilOfArthriticBindingBonus()
					weaponBaseDamage := core.BaseDamageFuncMeleeWeapon(core.MainHand, false, 560.0+bonusBaseDamage, 0.7, true)
					return weaponBaseDamage(sim, hitEffect, spell) *
						dk.rageOfRivendareBonus(hitEffect.Target) *
						dk.tundraStalkerBonus(hitEffect.Target)
				},
				TargetSpellCoefficient: 1,
			},

			OutcomeApplier: dk.OutcomeFuncMeleeSpecialHitAndCrit(dk.MeleeCritMultiplier(1.0, dk.viciousStrikesCritDamageBonus())),

			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				dk.LastCastOutcome = spellEffect.Outcome
				if spellEffect.Landed() {
					dkSpellCost := dk.DetermineOptimalCost(sim, 0, 1, 1)
					dk.Spend(sim, spell, dkSpellCost)

					amountOfRunicPower := 15.0 + 2.5*float64(dk.Talents.Dirge) + dk.scourgeborneBattlegearRunicPowerBonus()
					dk.AddRunicPower(sim, amountOfRunicPower, spell.RunicPowerMetrics())

					if dk.DiseasesAreActive(spellEffect.Target) {
						dk.LastScourgeStrikeDamage = spellEffect.Damage
						shadowDamageSpell.Cast(sim, spellEffect.Target)
						//dk.ScourgeStrike.SpellMetrics[spellEffect.Target.TableIndex].Casts -= 1
						//dk.ScourgeStrike.SpellMetrics[spellEffect.Target.TableIndex].Hits -= 1
					}
				}
			},
		}),
	})
}

func (dk *Deathknight) CanScourgeStrike(sim *core.Simulation) bool {
	return dk.Talents.ScourgeStrike && dk.CastCostPossible(sim, 0.0, 0, 1, 1) && dk.ScourgeStrike.IsReady(sim)
}

func (dk *Deathknight) CastScourgeStrike(sim *core.Simulation, target *core.Unit) bool {
	if dk.CanScourgeStrike(sim) {
		dk.ScourgeStrike.Cast(sim, target)
		return true
	}
	return false
}
