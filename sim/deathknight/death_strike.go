package deathknight

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

// TODO: Cleanup death strike the same way we did for plague strike
var DeathStrikeActionID = core.ActionID{SpellID: 49924}

func (dk *Deathknight) newDeathStrikeSpell(isMH bool) *core.Spell {
	bonusBaseDamage := dk.sigilOfAwarenessBonus()
	hasGlyph := dk.HasMajorGlyph(proto.DeathknightMajorGlyph_GlyphOfDeathStrike)
	deathConvertChance := float64(dk.Talents.DeathRuneMastery) / 3

	var healthMetrics *core.ResourceMetrics
	if isMH {
		healthMetrics = dk.NewHealthMetrics(DeathStrikeActionID)
	}

	conf := core.SpellConfig{
		ActionID:    DeathStrikeActionID.WithTag(core.TernaryInt32(isMH, 1, 2)),
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    dk.threatOfThassarianProcMask(isMH),
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage,

		RuneCost: core.RuneCostOptions{
			FrostRuneCost:  1,
			UnholyRuneCost: 1,
			RunicPowerGain: 15 + 2.5*float64(dk.Talents.Dirge),
			Refundable:     true,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
		},

		BonusCritRating: (dk.annihilationCritBonus() + dk.improvedDeathStrikeCritBonus()) * core.CritRatingPerCritChance,
		DamageMultiplier: .75 *
			core.TernaryFloat64(isMH, 1, dk.nervesOfColdSteelBonus()) *
			dk.improvedDeathStrikeDamageBonus(),
		CritMultiplier:   dk.bonusCritMultiplier(dk.Talents.MightOfMograine),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			var baseDamage float64
			if isMH {
				baseDamage = 297 +
					bonusBaseDamage +
					spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower()) +
					spell.BonusWeaponDamage()
			} else {
				baseDamage = 148 +
					bonusBaseDamage +
					spell.Unit.OHNormalizedWeaponDamage(sim, spell.MeleeAttackPower()) +
					spell.BonusWeaponDamage()
			}
			baseDamage *= dk.RoRTSBonus(target)
			if hasGlyph {
				baseDamage *= 1 + 0.01*min(dk.CurrentRunicPower(), 25)
			}

			result := spell.CalcDamage(sim, target, baseDamage, dk.threatOfThassarianOutcomeApplier(spell))

			if isMH {
				spell.SpendRefundableCostAndConvertFrostOrUnholyRune(sim, result, deathConvertChance)

				if result.Landed() {
					healingAmount := 0.05 * dk.dkCountActiveDiseases(target) * dk.MaxHealth() * (1.0 + 0.5*float64(dk.Talents.ImprovedDeathStrike))
					dk.GainHealth(sim, healingAmount*dk.PseudoStats.HealingTakenMultiplier, healthMetrics)
					dk.DeathStrikeHeals = append(dk.DeathStrikeHeals, healingAmount)
				}

				dk.threatOfThassarianProc(sim, result, dk.DeathStrikeOhHit)
			}

			spell.DealDamage(sim, result)
		},
	}

	if !isMH {
		conf.RuneCost = core.RuneCostOptions{}
		conf.Cast = core.CastConfig{}
	} else {
		conf.Flags |= core.SpellFlagAPL
	}

	return dk.RegisterSpell(conf)
}

func (dk *Deathknight) registerDeathStrikeSpell() {
	dk.DeathStrikeOhHit = dk.newDeathStrikeSpell(false)
	dk.DeathStrikeMhHit = dk.newDeathStrikeSpell(true)
	dk.DeathStrike = dk.DeathStrikeMhHit
}

func (dk *Deathknight) registerDrwDeathStrikeSpell() {
	bonusBaseDamage := dk.sigilOfAwarenessBonus()
	hasGlyph := dk.HasMajorGlyph(proto.DeathknightMajorGlyph_GlyphOfDeathStrike)

	dk.RuneWeapon.DeathStrike = dk.RuneWeapon.RegisterSpell(core.SpellConfig{
		ActionID:    DeathStrikeActionID.WithTag(1),
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage,

		BonusCritRating:  (dk.annihilationCritBonus() + dk.improvedDeathStrikeCritBonus()) * core.CritRatingPerCritChance,
		DamageMultiplier: .75 * dk.improvedDeathStrikeDamageBonus(),
		CritMultiplier:   dk.bonusCritMultiplier(dk.Talents.MightOfMograine),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 297 + bonusBaseDamage + dk.DrwWeaponDamage(sim, spell)

			if hasGlyph {
				baseDamage *= 1 + 0.01*min(dk.CurrentRunicPower(), 25)
			}
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)
		},
	})
}
