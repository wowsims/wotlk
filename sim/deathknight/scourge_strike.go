package deathknight

import (
	"github.com/wowsims/wotlk/sim/core"
)

var ScourgeStrikeActionID = core.ActionID{SpellID: 55271}

func (dk *Deathknight) registerScourgeStrikeShadowDamageSpell() *core.Spell {
	diseaseMulti := dk.diseaseMultiplier(0.12)

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
					return dk.LastScourgeStrikeDamage * (diseaseMulti * dk.countActiveDiseases(hitEffect.Target))
				},
			},
		}),
	})
}

func (dk *Deathknight) registerScourgeStrikeSpell() {

	shadowDamageSpell := dk.registerScourgeStrikeShadowDamageSpell()
	bonusBaseDamage := dk.sigilOfAwarenessBonus(dk.ScourgeStrike)
	bonusBaseDamage += dk.sigilOfArthriticBindingBonus()
	weaponBaseDamage := core.BaseDamageFuncMeleeWeapon(core.MainHand, true, 800.0+bonusBaseDamage, 0.7, true)
	outbreakBonus := []float64{1.0, 1.07, 1.13, 1.2}[dk.Talents.Outbreak]

	amountOfRunicPower := 15.0 + 2.5*float64(dk.Talents.Dirge) + dk.scourgeborneBattlegearRunicPowerBonus()
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
			DamageMultiplier: outbreakBonus * dk.scourgelordsBattlegearDamageBonus(dk.ScourgeStrike),
			ThreatMultiplier: 1,

			BaseDamage: core.BaseDamageConfig{
				Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
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
					if dk.DiseasesAreActive(spellEffect.Target) {
						dk.LastScourgeStrikeDamage = spellEffect.Damage
						shadowDamageSpell.Cast(sim, spellEffect.Target)
					}

					dkSpellCost := dk.DetermineCost(sim, core.DKCastEnum_FU)
					dk.Spend(sim, spell, dkSpellCost)

					dk.AddRunicPower(sim, amountOfRunicPower, spell.RunicPowerMetrics())
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
