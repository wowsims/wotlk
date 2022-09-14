package deathknight

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

var BloodStrikeActionID = core.ActionID{SpellID: 49930}

func (dk *Deathknight) newBloodStrikeSpell(isMH bool, onhit func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect)) *RuneSpell {
	weaponBaseDamage := core.BaseDamageFuncMeleeWeapon(core.MainHand, true, 764.0, true)
	if !isMH {
		weaponBaseDamage = core.BaseDamageFuncMeleeWeapon(core.OffHand, true, 764.0*0.5, true)
	}

	diseaseMulti := dk.dkDiseaseMultiplier(0.125)
	weaponMulti := 0.4
	if !isMH {
		weaponMulti = 0.4 * dk.nervesOfColdSteelBonus()
	}

	effect := core.SpellEffect{
		BonusCritRating:  (dk.subversionCritBonus() + dk.annihilationCritBonus()) * core.CritRatingPerCritChance,
		DamageMultiplier: weaponMulti * dk.bloodOfTheNorthCoeff() * dk.thassariansPlateDamageBonus() * dk.bloodyStrikesBonus(dk.BloodStrike),
		ThreatMultiplier: 1,

		BaseDamage: core.BaseDamageConfig{
			Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
				return weaponBaseDamage(sim, hitEffect, spell) *
					(1.0 + dk.dkCountActiveDiseases(hitEffect.Target)*diseaseMulti) * dk.RoRTSBonus(hitEffect.Target)
			},
			TargetSpellCoefficient: 1,
		},
		OnSpellHitDealt: onhit,
	}

	dk.threatOfThassarianProcMasks(isMH, &effect, true, true, func(outcomeApplier core.OutcomeApplier) core.OutcomeApplier {
		return outcomeApplier
	})
	conf := core.SpellConfig{
		ActionID:     BloodStrikeActionID.WithTag(core.TernaryInt32(isMH, 1, 2)),
		SpellSchool:  core.SpellSchoolPhysical,
		Flags:        core.SpellFlagMeleeMetrics,
		ApplyEffects: core.ApplyEffectFuncDirectDamage(effect),
	}
	rs := &RuneSpell{}
	if isMH { // offhand doesnt need GCD
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
		if dk.Talents.ThreatOfThassarian > 0 && dk.GetOHWeapon() != nil && dk.threatOfThassarianWillProc(sim) {
			dk.BloodStrikeOhHit.Cast(sim, spellEffect.Target)
		}
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
