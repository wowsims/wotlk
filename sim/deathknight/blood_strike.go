package deathknight

import (
	"github.com/wowsims/wotlk/sim/core"
)

var BloodStrikeActionID = core.ActionID{SpellID: 49930}

func (dk *Deathknight) newBloodStrikeSpell(isMH bool) *core.Spell {
	bonusBaseDamage := dk.sigilOfTheDarkRiderBonus()
	diseaseMulti := dk.dkDiseaseMultiplier(0.125)
	deathConvertChance := float64(dk.Talents.BloodOfTheNorth+dk.Talents.Reaping) / 3

	conf := core.SpellConfig{
		ActionID:    BloodStrikeActionID.WithTag(core.TernaryInt32(isMH, 1, 2)),
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    dk.threatOfThassarianProcMask(isMH),
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
		DamageMultiplier: 0.4 *
			core.TernaryFloat64(isMH, 1, dk.nervesOfColdSteelBonus()) *
			dk.bloodOfTheNorthCoeff() *
			dk.thassariansPlateDamageBonus() *
			dk.bloodyStrikesBonus(BloodyStrikesBS),
		CritMultiplier:   dk.bonusCritMultiplier(dk.Talents.MightOfMograine + dk.Talents.GuileOfGorefiend),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			var baseDamage float64
			if isMH {
				baseDamage = 764 +
					bonusBaseDamage +
					spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower()) +
					spell.BonusWeaponDamage()
			} else {
				// SpellID 66979
				baseDamage = 382 +
					bonusBaseDamage +
					spell.Unit.OHNormalizedWeaponDamage(sim, spell.MeleeAttackPower()) +
					spell.BonusWeaponDamage()
			}
			baseDamage *= dk.RoRTSBonus(target) *
				(1.0 + dk.dkCountActiveDiseases(target)*diseaseMulti)

			result := spell.CalcDamage(sim, target, baseDamage, dk.threatOfThassarianOutcomeApplier(spell))

			if isMH {
				spell.SpendRefundableCostAndConvertBloodRune(sim, result, deathConvertChance)
				dk.threatOfThassarianProc(sim, result, dk.BloodStrikeOhHit)

				if result.Landed() {
					if dk.DesolationAura != nil {
						dk.DesolationAura.Activate(sim)
					}
				}
			}

			spell.DealDamage(sim, result)
		},
	}

	if !isMH { // offhand doesn't need GCD
		conf.RuneCost = core.RuneCostOptions{}
		conf.Cast = core.CastConfig{}
	} else {
		conf.Flags |= core.SpellFlagAPL
	}

	return dk.RegisterSpell(conf)
}

func (dk *Deathknight) registerBloodStrikeSpell() {
	dk.BloodStrikeMhHit = dk.newBloodStrikeSpell(true)
	dk.BloodStrikeOhHit = dk.newBloodStrikeSpell(false)
	dk.BloodStrike = dk.BloodStrikeMhHit
}

func (dk *Deathknight) registerDrwBloodStrikeSpell() {
	bonusBaseDamage := dk.sigilOfTheDarkRiderBonus()
	diseaseMulti := dk.dkDiseaseMultiplier(0.125)

	dk.RuneWeapon.BloodStrike = dk.RuneWeapon.RegisterSpell(core.SpellConfig{
		ActionID:    BloodStrikeActionID.WithTag(1),
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage,

		BonusCritRating: (dk.subversionCritBonus() + dk.annihilationCritBonus()) * core.CritRatingPerCritChance,
		DamageMultiplier: 0.4 *
			dk.bloodOfTheNorthCoeff() *
			dk.thassariansPlateDamageBonus() *
			dk.bloodyStrikesBonus(BloodyStrikesBS),
		CritMultiplier:   dk.bonusCritMultiplier(dk.Talents.MightOfMograine + dk.Talents.GuileOfGorefiend),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 764 + bonusBaseDamage + dk.DrwWeaponDamage(sim, spell)

			baseDamage *= dk.RoRTSBonus(target) *
				(1.0 + dk.drwCountActiveDiseases(target)*diseaseMulti)

			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)
		},
	})
}
