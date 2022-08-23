package deathknight

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

var HeartStrikeActionID = core.ActionID{SpellID: 55050}

func (dk *Deathknight) newHeartStrikeSpell(isMainTarget bool, isDrw bool, onhit func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect)) *RuneSpell {
	weaponBaseDamage := core.BaseDamageFuncMeleeWeapon(core.MainHand, true, 250.0, 1.0, 0.5, true)
	if !isMainTarget {
		weaponBaseDamage = core.BaseDamageFuncMeleeWeapon(core.MainHand, true, 250.0, 1.0, 0.25, true)
	}

	diseaseMulti := dk.dkDiseaseMultiplier(0.1)

	outcomeApplier := dk.OutcomeFuncMeleeSpecialHitAndCrit(dk.critMultiplierGoGandMoM())
	if isDrw {
		outcomeApplier = dk.RuneWeapon.OutcomeFuncMeleeSpecialHitAndCrit(
			dk.RuneWeapon.MeleeCritMultiplier(1.0, dk.secondaryCritModifier(dk.Talents.GuileOfGorefiend > 0, dk.Talents.MightOfMograine > 0)))
	}

	effect := core.SpellEffect{
		ProcMask:         core.ProcMaskMeleeSpecial,
		BonusCritRating:  (dk.subversionCritBonus() + dk.annihilationCritBonus()) * core.CritRatingPerCritChance,
		DamageMultiplier: dk.thassariansPlateDamageBonus() * dk.scourgelordsBattlegearDamageBonus(dk.HeartStrike) * dk.bloodyStrikesBonus(dk.HeartStrike),
		ThreatMultiplier: 1,
		OutcomeApplier:   outcomeApplier,
		BaseDamage: core.BaseDamageConfig{
			Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
				activeDiseases := core.TernaryFloat64(isDrw, dk.drwCountActiveDiseases(hitEffect.Target), dk.dkCountActiveDiseases(hitEffect.Target))
				return weaponBaseDamage(sim, hitEffect, spell) *
					(1.0 + activeDiseases*diseaseMulti)
			},
			TargetSpellCoefficient: 1,
		},
		OnSpellHitDealt: onhit,
	}

	conf := core.SpellConfig{
		ActionID:     HeartStrikeActionID,
		SpellSchool:  core.SpellSchoolPhysical,
		Flags:        core.SpellFlagMeleeMetrics,
		ApplyEffects: core.ApplyEffectFuncDirectDamage(effect),
	}
	rs := &RuneSpell{}
	if isMainTarget && !isDrw { // off target doesnt need GCD
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
	}

	if isDrw {
		rs.Spell = dk.RuneWeapon.RegisterSpell(conf)
		return rs
	} else {
		return dk.RegisterSpell(rs, conf, func(sim *core.Simulation) bool {
			return dk.CastCostPossible(sim, 0.0, 1, 0, 0) && dk.HeartStrike.IsReady(sim)
		}, nil)
	}
}

func (dk *Deathknight) registerHeartStrikeSpell() {
	dk.HeartStrike = dk.newHeartStrikeSpell(true, false, func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
		if dk.Env.GetNumTargets() > 1 {
			dk.HeartStrikeOffHit.Cast(sim, dk.Env.NextTargetUnit(dk.CurrentTarget))
		}
		dk.LastOutcome = spellEffect.Outcome
	})
	dk.HeartStrikeOffHit = dk.newHeartStrikeSpell(false, false, nil)
}

func (dk *Deathknight) registerDrwHeartStrikeSpell() {
	dk.RuneWeapon.HeartStrike = dk.newHeartStrikeSpell(true, true, func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
		if dk.Env.GetNumTargets() > 1 {
			dk.RuneWeapon.HeartStrikeOffHit.Cast(sim, dk.Env.NextTargetUnit(dk.CurrentTarget))
		}
	}).Spell
	dk.RuneWeapon.HeartStrikeOffHit = dk.newHeartStrikeSpell(false, true, nil).Spell
}
