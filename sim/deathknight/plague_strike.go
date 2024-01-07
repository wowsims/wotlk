package deathknight

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

var PlagueStrikeActionID = core.ActionID{SpellID: 49921}

func (dk *Deathknight) newPlagueStrikeSpell(isMH bool) *core.Spell {
	conf := core.SpellConfig{
		ActionID:    PlagueStrikeActionID.WithTag(core.TernaryInt32(isMH, 1, 2)),
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    dk.threatOfThassarianProcMask(isMH),
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage,

		RuneCost: core.RuneCostOptions{
			UnholyRuneCost: 1,
			RunicPowerGain: 10 + 2.5*float64(dk.Talents.Dirge),
			Refundable:     true,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
		},

		BonusCritRating: (dk.annihilationCritBonus() + dk.scourgebornePlateCritBonus() + dk.viciousStrikesCritChanceBonus()) * core.CritRatingPerCritChance,
		DamageMultiplier: .5 *
			core.TernaryFloat64(isMH, 1, dk.nervesOfColdSteelBonus()) *
			(1.0 + 0.1*float64(dk.Talents.Outbreak)) *
			core.TernaryFloat64(dk.HasMajorGlyph(proto.DeathknightMajorGlyph_GlyphOfPlagueStrike), 1.2, 1.0),
		CritMultiplier:   dk.bonusCritMultiplier(dk.Talents.ViciousStrikes),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			var baseDamage float64
			if isMH {
				baseDamage = 378 +
					spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower()) +
					spell.BonusWeaponDamage()
			} else {
				// SpellID 66992
				baseDamage = 189 +
					spell.Unit.OHNormalizedWeaponDamage(sim, spell.MeleeAttackPower()) +
					spell.BonusWeaponDamage()
			}
			baseDamage *= dk.RoRTSBonus(target)

			result := spell.CalcDamage(sim, target, baseDamage, dk.threatOfThassarianOutcomeApplier(spell))

			if isMH {
				spell.SpendRefundableCost(sim, result)
				dk.threatOfThassarianProc(sim, result, dk.PlagueStrikeOhHit)
				if result.Landed() {
					dk.BloodPlagueExtended[target.Index] = 0
					dk.BloodPlagueSpell.Cast(sim, target)
				}
			}

			spell.DealDamage(sim, result)
		},
	}
	if !isMH { // only MH has cost & gcd
		conf.RuneCost = core.RuneCostOptions{}
		conf.Cast = core.CastConfig{}
	} else {
		conf.Flags |= core.SpellFlagAPL
	}

	return dk.RegisterSpell(conf)
}

func (dk *Deathknight) registerPlagueStrikeSpell() {
	dk.PlagueStrikeMhHit = dk.newPlagueStrikeSpell(true)
	dk.PlagueStrikeOhHit = dk.newPlagueStrikeSpell(false)
	dk.PlagueStrike = dk.PlagueStrikeMhHit
}
func (dk *Deathknight) registerDrwPlagueStrikeSpell() {
	dk.RuneWeapon.PlagueStrike = dk.RuneWeapon.RegisterSpell(core.SpellConfig{
		ActionID:    PlagueStrikeActionID.WithTag(1),
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage,

		BonusCritRating: (dk.annihilationCritBonus() + dk.scourgebornePlateCritBonus() + dk.viciousStrikesCritChanceBonus()) * core.CritRatingPerCritChance,
		DamageMultiplier: 0.5 *
			(1.0 + 0.1*float64(dk.Talents.Outbreak)),
		CritMultiplier:   dk.bonusCritMultiplier(dk.Talents.ViciousStrikes),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 378 + dk.DrwWeaponDamage(sim, spell)

			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)

			if result.Landed() {
				dk.RuneWeapon.BloodPlagueSpell.Cast(sim, target)
			}
		},
	})
}
