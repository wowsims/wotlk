package deathknight

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

var HeartStrikeActionID = core.ActionID{SpellID: 55262}

func (dk *Deathknight) newHeartStrikeSpell(isMainTarget bool, isDrw bool) *RuneSpell {
	bonusBaseDamage := dk.sigilOfTheDarkRiderBonus()
	diseaseMulti := dk.dkDiseaseMultiplier(0.1)

	critMultiplier := dk.bonusCritMultiplier(dk.Talents.MightOfMograine)

	rs := &RuneSpell{}
	conf := core.SpellConfig{
		ActionID:    HeartStrikeActionID.WithTag(core.TernaryInt32(isMainTarget, 1, 2)),
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage,

		BonusCritRating: (dk.subversionCritBonus() + dk.annihilationCritBonus()) * core.CritRatingPerCritChance,
		DamageMultiplier: .5 *
			core.TernaryFloat64(isMainTarget, 1, 0.5) *
			dk.thassariansPlateDamageBonus() *
			dk.scourgelordsBattlegearDamageBonus(dk.HeartStrike) *
			dk.bloodyStrikesBonus(dk.HeartStrike),
		CritMultiplier:   critMultiplier,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 736 +
				bonusBaseDamage +
				spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower()) +
				spell.BonusWeaponDamage()

			activeDiseases := core.TernaryFloat64(isDrw, dk.drwCountActiveDiseases(target), dk.dkCountActiveDiseases(target))
			baseDamage *= 1 + activeDiseases*diseaseMulti

			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)

			if isMainTarget {
				if isDrw {
					if dk.Env.GetNumTargets() > 1 {
						dk.RuneWeapon.HeartStrikeOffHit.Cast(sim, dk.Env.NextTargetUnit(dk.CurrentTarget))
					}
				} else {
					rs.OnResult(sim, result)

					if dk.Env.GetNumTargets() > 1 {
						dk.HeartStrikeOffHit.Cast(sim, dk.Env.NextTargetUnit(dk.CurrentTarget))
					}
					dk.LastOutcome = result.Outcome
				}
			}
		},
	}
	if isDrw {
		conf.DamageMultiplier *= .5
		conf.Flags |= core.SpellFlagIgnoreAttackerModifiers
	}
	if isMainTarget && !isDrw { // off target doesnt need GCD
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
	}

	if isDrw {
		rs.Spell = dk.RuneWeapon.RegisterSpell(conf)
		return rs
	} else {
		return dk.RegisterSpell(rs, conf, func(sim *core.Simulation) bool {
			return dk.CastCostPossible(sim, 0.0, 1, 0, 0) && dk.HeartStrike.IsReady(sim)
		})
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

	dk.RuneWeapon.HeartStrike = dk.newHeartStrikeSpell(true, true).Spell
	dk.RuneWeapon.HeartStrikeOffHit = dk.newHeartStrikeSpell(false, true).Spell
}
