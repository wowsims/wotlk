package deathknight

import (
	"github.com/wowsims/wotlk/sim/core"
)

var HeartStrikeActionID = core.ActionID{SpellID: 55262}

func (dk *Deathknight) newHeartStrikeSpell(isMainTarget bool, isDrw bool) *core.Spell {
	bonusBaseDamage := dk.sigilOfTheDarkRiderBonus()
	diseaseMulti := dk.dkDiseaseMultiplier(0.1)

	critMultiplier := dk.bonusCritMultiplier(dk.Talents.MightOfMograine)

	conf := core.SpellConfig{
		ActionID:    HeartStrikeActionID.WithTag(core.TernaryInt32(isMainTarget, 1, 2)),
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage,

		RuneCost: core.RuneCostOptions{
			BloodRuneCost:  1,
			RunicPowerGain: 10,
			Refundable:     true,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
		},

		BonusCritRating: (dk.subversionCritBonus() + dk.annihilationCritBonus()) * core.CritRatingPerCritChance,
		DamageMultiplier: .5 *
			core.TernaryFloat64(isMainTarget, 1, 0.5) *
			dk.thassariansPlateDamageBonus() *
			dk.scourgelordsBattlegearDamageBonus(ScourgelordBonusSpellHS) *
			dk.bloodyStrikesBonus(BloodyStrikesHS),
		CritMultiplier:   critMultiplier,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 736 + bonusBaseDamage

			if isDrw {
				baseDamage += dk.DrwWeaponDamage(sim, spell)
			} else {
				baseDamage += spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower()) +
					spell.BonusWeaponDamage()
			}

			activeDiseases := core.TernaryFloat64(isDrw, dk.drwCountActiveDiseases(target), dk.dkCountActiveDiseases(target))
			baseDamage *= 1 + activeDiseases*diseaseMulti

			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)

			if isMainTarget {
				if isDrw {
					if dk.Env.GetNumTargets() > 1 {
						dk.RuneWeapon.HeartStrikeOffHit.Cast(sim, dk.Env.NextTargetUnit(target))
					}
				} else {
					spell.SpendRefundableCost(sim, result)

					if dk.Env.GetNumTargets() > 1 {
						dk.HeartStrikeOffHit.Cast(sim, dk.Env.NextTargetUnit(target))
					}
				}
			}
		},
	}
	if !isMainTarget || isDrw { // off target doesnt need GCD
		conf.RuneCost = core.RuneCostOptions{}
		conf.Cast = core.CastConfig{}
	}

	if isMainTarget {
		conf.Flags |= core.SpellFlagAPL
	}

	if isDrw {
		return dk.RuneWeapon.RegisterSpell(conf)
	} else {
		return dk.RegisterSpell(conf)
	}
}

func (dk *Deathknight) registerHeartStrikeSpell() {
	if !dk.Talents.HeartStrike {
		return
	}

	dk.HeartStrike = dk.newHeartStrikeSpell(true, false)
	dk.HeartStrikeOffHit = dk.newHeartStrikeSpell(false, false)
}

func (dk *Deathknight) registerDrwHeartStrikeSpell() {
	if !dk.Talents.HeartStrike {
		return
	}

	dk.RuneWeapon.HeartStrike = dk.newHeartStrikeSpell(true, true)
	dk.RuneWeapon.HeartStrikeOffHit = dk.newHeartStrikeSpell(false, true)
}
