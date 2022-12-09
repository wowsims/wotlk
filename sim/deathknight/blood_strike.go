package deathknight

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

var BloodStrikeActionID = core.ActionID{SpellID: 49930}

func (dk *Deathknight) newBloodStrikeSpell(isMH bool) *RuneSpell {
	bonusBaseDamage := dk.sigilOfTheDarkRiderBonus()
	diseaseMulti := dk.dkDiseaseMultiplier(0.125)

	rs := &RuneSpell{}
	conf := core.SpellConfig{
		ActionID:    BloodStrikeActionID.WithTag(core.TernaryInt32(isMH, 1, 2)),
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    dk.threatOfThassarianProcMask(isMH),
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage,

		BonusCritRating: (dk.subversionCritBonus() + dk.annihilationCritBonus()) * core.CritRatingPerCritChance,
		DamageMultiplier: 0.4 *
			core.TernaryFloat64(isMH, 1, dk.nervesOfColdSteelBonus()) *
			dk.bloodOfTheNorthCoeff() *
			dk.thassariansPlateDamageBonus() *
			dk.bloodyStrikesBonus(dk.BloodStrike),
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
				rs.OnResult(sim, result)
				dk.threatOfThassarianProc(sim, result, dk.BloodStrikeOhHit)
				dk.LastOutcome = result.Outcome

				if result.Landed() {
					if dk.DesolationAura != nil {
						dk.DesolationAura.Activate(sim)
					}
				}
			}

			spell.DealDamage(sim, result)
		},
	}

	if isMH { // offhand doesn't need GCD
		conf.ResourceType = stats.RunicPower
		conf.BaseCost = float64(core.NewRuneCost(10, 1, 0, 0, 0))
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
		rs.Refundable = true
		rs.ConvertType = RuneTypeBlood
		if dk.Talents.BloodOfTheNorth+dk.Talents.Reaping >= 3 {
			rs.DeathConvertChance = 1.0
		} else {
			rs.DeathConvertChance = float64(dk.Talents.BloodOfTheNorth+dk.Talents.Reaping) * 0.33
		}
	}

	if isMH {
		return dk.RegisterSpell(rs, conf, func(sim *core.Simulation) bool {
			return dk.CastCostPossible(sim, 0.0, 1, 0, 0) && dk.BloodStrike.IsReady(sim)
		}, nil)
	} else {
		return dk.RegisterSpell(rs, conf, nil, nil)
	}
}

func (dk *Deathknight) registerBloodStrikeSpell() {
	dk.BloodStrikeMhHit = dk.newBloodStrikeSpell(true)
	dk.BloodStrikeOhHit = dk.newBloodStrikeSpell(false)
	dk.BloodStrike = dk.BloodStrikeMhHit
}
