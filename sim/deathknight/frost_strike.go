package deathknight

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

var frostStrikeActionID = core.ActionID{SpellID: 55268}
var FrostStrikeMHActionID = frostStrikeActionID.WithTag(1)
var FrostStrikeOHActionID = frostStrikeActionID.WithTag(2)

func (dk *Deathknight) newFrostStrikeHitSpell(isMH bool) *core.Spell {
	bonusBaseDamage := dk.sigilOfTheVengefulHeartFrostStrike()

	actionID := FrostStrikeMHActionID
	if !isMH {
		actionID = FrostStrikeOHActionID
	}

	conf := core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolFrost,
		ProcMask:    dk.threatOfThassarianProcMask(isMH),
		Flags:       core.SpellFlagMeleeMetrics,

		RuneCost: core.RuneCostOptions{
			RunicPowerCost: core.TernaryFloat64(dk.HasMajorGlyph(proto.DeathknightMajorGlyph_GlyphOfFrostStrike), 32, 40),
			Refundable:     true,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
		},

		BonusCritRating: (dk.annihilationCritBonus() + dk.darkrunedBattlegearCritBonus()) * core.CritRatingPerCritChance,
		DamageMultiplier: .55 *
			core.TernaryFloat64(isMH, 1, dk.nervesOfColdSteelBonus()) *
			dk.bloodOfTheNorthCoeff(),
		CritMultiplier:   dk.bonusCritMultiplier(dk.Talents.GuileOfGorefiend),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			var baseDamage float64
			if isMH {
				baseDamage = 250 +
					bonusBaseDamage +
					spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower()) +
					spell.BonusWeaponDamage()
			} else {
				// SpellID 66962
				baseDamage = 125 +
					bonusBaseDamage +
					spell.Unit.OHNormalizedWeaponDamage(sim, spell.MeleeAttackPower()) +
					spell.BonusWeaponDamage()
			}
			baseDamage *= dk.glacielRotBonus(target) *
				dk.RoRTSBonus(target) *
				dk.mercilessCombatBonus(sim)

			result := spell.CalcDamage(sim, target, baseDamage, dk.threatOfThassarianOutcomeApplier(spell))

			if isMH {
				spell.SpendRefundableCost(sim, result)
				dk.threatOfThassarianProc(sim, result, dk.FrostStrikeOhHit)
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

func (dk *Deathknight) registerFrostStrikeSpell() {
	if !dk.Talents.FrostStrike {
		return
	}

	dk.FrostStrikeMhHit = dk.newFrostStrikeHitSpell(true)
	dk.FrostStrikeOhHit = dk.newFrostStrikeHitSpell(false)
	dk.FrostStrike = dk.FrostStrikeMhHit
}
