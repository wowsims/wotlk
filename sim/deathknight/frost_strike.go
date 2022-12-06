package deathknight

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

var FrostStrikeActionID = core.ActionID{SpellID: 55268}

func (dk *Deathknight) newFrostStrikeHitSpell(isMH bool) *RuneSpell {
	bonusBaseDamage := dk.sigilOfTheVengefulHeartFrostStrike()

	conf := core.SpellConfig{
		ActionID:    FrostStrikeActionID.WithTag(core.TernaryInt32(isMH, 1, 2)),
		SpellSchool: core.SpellSchoolFrost,
		ProcMask:    dk.threatOfThassarianProcMask(isMH),
		Flags:       core.SpellFlagMeleeMetrics,

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
				dk.LastOutcome = result.Outcome
				dk.threatOfThassarianProc(sim, result, dk.FrostStrikeOhHit)
			}

			spell.DealDamage(sim, result)
		},
	}

	rs := &RuneSpell{}
	if isMH {
		conf.ResourceType = stats.RunicPower
		conf.BaseCost = float64(core.NewRuneCost(
			core.Ternary(dk.HasMajorGlyph(proto.DeathknightMajorGlyph_GlyphOfFrostStrike), uint8(32), 40), 0, 0, 0, 0,
		))
		conf.Cast = core.CastConfig{
			DefaultCast: core.Cast{
				GCD:  core.GCDDefault,
				Cost: conf.BaseCost,
			},
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				cast.GCD = dk.GetModifiedGCD()
			},
			IgnoreHaste: true,
		}
	}

	if isMH {
		return dk.RegisterSpell(rs, conf, func(sim *core.Simulation) bool {
			runeCost := core.RuneCost(dk.FrostStrike.BaseCost)
			return dk.CastCostPossible(sim, float64(runeCost.RunicPower()), 0, 0, 0) && dk.FrostStrike.IsReady(sim)
		}, nil)
	} else {
		return dk.RegisterSpell(rs, conf, nil, nil)
	}
}

func (dk *Deathknight) registerFrostStrikeSpell() {
	dk.FrostStrikeMhHit = dk.newFrostStrikeHitSpell(true)
	dk.FrostStrikeOhHit = dk.newFrostStrikeHitSpell(false)
	dk.FrostStrike = dk.FrostStrikeMhHit
}
