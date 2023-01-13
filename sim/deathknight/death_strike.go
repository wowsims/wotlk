package deathknight

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

// TODO: Cleanup death strike the same way we did for plague strike
var DeathStrikeActionID = core.ActionID{SpellID: 49924}

func (dk *Deathknight) newDeathStrikeSpell(isMH bool) *RuneSpell {
	bonusBaseDamage := dk.sigilOfAwarenessBonus()
	hasGlyph := dk.HasMajorGlyph(proto.DeathknightMajorGlyph_GlyphOfDeathStrike)

	var healthMetrics *core.ResourceMetrics
	if isMH {
		healthMetrics = dk.NewHealthMetrics(DeathStrikeActionID)
	}

	rs := &RuneSpell{}
	conf := core.SpellConfig{
		ActionID:    DeathStrikeActionID.WithTag(core.TernaryInt32(isMH, 1, 2)),
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    dk.threatOfThassarianProcMask(isMH),
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage,

		RuneCost: core.RuneCostOptions{
			FrostRuneCost:  1,
			UnholyRuneCost: 1,
			RunicPowerGain: 15 + 2.5*float64(dk.Talents.Dirge),
		},

		BonusCritRating: (dk.annihilationCritBonus() + dk.improvedDeathStrikeCritBonus()) * core.CritRatingPerCritChance,
		DamageMultiplier: .75 *
			core.TernaryFloat64(isMH, 1, dk.nervesOfColdSteelBonus()) *
			(1.0 + 0.15*float64(dk.Talents.ImprovedDeathStrike)),
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
				baseDamage *= 1 + 0.01*core.MinFloat(dk.CurrentRunicPower(), 25)
			}

			result := spell.CalcDamage(sim, target, baseDamage, dk.threatOfThassarianOutcomeApplier(spell))

			if isMH {
				rs.OnResult(sim, result)
				dk.LastOutcome = result.Outcome

				if result.Landed() {
					healingAmount := 0.05 * dk.dkCountActiveDiseases(target) * dk.MaxHealth() * (1.0 + 0.5*float64(dk.Talents.ImprovedDeathStrike)) * (1.0 + core.TernaryFloat64(dk.VampiricBloodAura.IsActive(), 0.35, 0.0))
					dk.GainHealth(sim, healingAmount, healthMetrics)
					dk.DeathStrikeHeals = append(dk.DeathStrikeHeals, healingAmount)
				}

				dk.threatOfThassarianProc(sim, result, dk.DeathStrikeOhHit)
			}

			spell.DealDamage(sim, result)
		},
	}

	if isMH {
		conf.Cast = core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				cast.GCD = dk.GetModifiedGCD()
			},
			IgnoreHaste: true,
		}
		rs.Refundable = true
		rs.ConvertType = RuneTypeFrost | RuneTypeUnholy
		if dk.Talents.DeathRuneMastery == 3 {
			rs.DeathConvertChance = 1.0
		} else {
			rs.DeathConvertChance = float64(dk.Talents.DeathRuneMastery) * 0.33
		}
	} else {
		conf.RuneCost = core.RuneCostOptions{}
	}

	return dk.RegisterSpell(rs, conf)
}

func (dk *Deathknight) registerDeathStrikeSpell() {
	dk.DeathStrikeOhHit = dk.newDeathStrikeSpell(false)
	dk.DeathStrikeMhHit = dk.newDeathStrikeSpell(true)
	dk.DeathStrike = dk.DeathStrikeMhHit
}

func (dk *Deathknight) registerDrwDeathStrikeSpell() {
	bonusBaseDamage := dk.sigilOfAwarenessBonus()

	dk.RuneWeapon.DeathStrike = dk.RuneWeapon.RegisterSpell(core.SpellConfig{
		ActionID:    DeathStrikeActionID.WithTag(1),
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | core.SpellFlagIgnoreAttackerModifiers,

		BonusCritRating:  (dk.annihilationCritBonus() + dk.improvedDeathStrikeCritBonus()) * core.CritRatingPerCritChance,
		DamageMultiplier: .5 * .75 * dk.improvedDeathStrikeDamageBonus(),
		CritMultiplier:   dk.bonusCritMultiplier(dk.Talents.MightOfMograine),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 297 +
				bonusBaseDamage +
				spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower()) +
				spell.BonusWeaponDamage()

			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)
		},
	})
}
