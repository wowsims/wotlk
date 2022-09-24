package deathknight

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

var BloodStrikeActionID = core.ActionID{SpellID: 49930}

func (dk *Deathknight) newBloodStrikeSpell(isMH bool, onhit func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect)) *RuneSpell {
	bonusBaseDamage := dk.sigilOfTheDarkRiderBonus()
	weaponBaseDamage := core.BaseDamageFuncMeleeWeapon(core.MainHand, true, 764.0+bonusBaseDamage, true)
	if !isMH {
		// SpellID 66979
		weaponBaseDamage = core.BaseDamageFuncMeleeWeapon(core.OffHand, true, 382.0+bonusBaseDamage, true)
	}

	diseaseMulti := dk.dkDiseaseMultiplier(0.125)

	effect := core.SpellEffect{
		BaseDamage: core.BaseDamageConfig{
			Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
				return weaponBaseDamage(sim, hitEffect, spell) *
					(1.0 + dk.dkCountActiveDiseases(hitEffect.Target)*diseaseMulti) * dk.RoRTSBonus(hitEffect.Target)
			},
		},
		OnSpellHitDealt: onhit,
	}

	procMask := dk.threatOfThassarianProcMasks(isMH, &effect)

	conf := core.SpellConfig{
		ActionID:    BloodStrikeActionID.WithTag(core.TernaryInt32(isMH, 1, 2)),
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    procMask,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage,

		BonusCritRating: (dk.subversionCritBonus() + dk.annihilationCritBonus()) * core.CritRatingPerCritChance,
		DamageMultiplier: 0.4 *
			core.TernaryFloat64(isMH, 1, dk.nervesOfColdSteelBonus()) *
			dk.bloodOfTheNorthCoeff() *
			dk.thassariansPlateDamageBonus() *
			dk.bloodyStrikesBonus(dk.BloodStrike),
		CritMultiplier:   dk.bonusCritMultiplier(dk.Talents.MightOfMograine + dk.Talents.GuileOfGorefiend),
		ThreatMultiplier: 1,

		ApplyEffects: core.ApplyEffectFuncDirectDamage(effect),
	}

	rs := &RuneSpell{}
	if isMH { // offhand doesn't need GCD
		conf.ResourceType = stats.RunicPower
		conf.BaseCost = float64(core.NewRuneCost(10, 1, 0, 0, 0))
		conf.Cast = core.CastConfig{
			DefaultCast: core.Cast{
				GCD:  core.GCDDefault,
				Cost: conf.BaseCost,
			},
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				cast.GCD = dk.getModifiedGCD()
			},
			IgnoreHaste: true,
		}
		conf.ApplyEffects = dk.withRuneRefund(rs, effect, false)
		if dk.Talents.BloodOfTheNorth+dk.Talents.Reaping >= 3 {
			rs.DeathConvertChance = 1.0
		} else {
			rs.DeathConvertChance = float64(dk.Talents.BloodOfTheNorth+dk.Talents.Reaping) * 0.33
		}
		rs.ConvertType = RuneTypeBlood
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
	dk.BloodStrikeMhHit = dk.newBloodStrikeSpell(true, func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
		dk.threatOfThassarianProc(sim, spellEffect, dk.BloodStrikeOhHit)
		dk.LastOutcome = spellEffect.Outcome

		if spellEffect.Outcome.Matches(core.OutcomeLanded) {
			if dk.DesolationAura != nil {
				dk.DesolationAura.Activate(sim)
			}
		}
	})
	dk.BloodStrikeOhHit = dk.newBloodStrikeSpell(false, nil)
	dk.BloodStrike = dk.BloodStrikeMhHit
}
