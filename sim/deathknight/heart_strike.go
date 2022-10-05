package deathknight

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

var HeartStrikeActionID = core.ActionID{SpellID: 55262}

func (dk *Deathknight) newHeartStrikeSpell(isMainTarget bool, isDrw bool, onhit func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect)) *RuneSpell {
	bonusBaseDamage := dk.sigilOfTheDarkRiderBonus()
	weaponBaseDamage := core.BaseDamageFuncMeleeWeapon(core.MainHand, true, 736.0+bonusBaseDamage, true)

	diseaseMulti := dk.dkDiseaseMultiplier(0.1)

	critMultiplier := dk.bonusCritMultiplier(dk.Talents.MightOfMograine)
	outcomeApplier := dk.OutcomeFuncMeleeSpecialHitAndCrit()
	if isDrw {
		critMultiplier = dk.RuneWeapon.DefaultMeleeCritMultiplier()
		outcomeApplier = dk.RuneWeapon.OutcomeFuncMeleeSpecialHitAndCrit()
	}

	effect := core.SpellEffect{
		OutcomeApplier: outcomeApplier,
		BaseDamage: core.BaseDamageConfig{
			Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
				activeDiseases := core.TernaryFloat64(isDrw, dk.drwCountActiveDiseases(hitEffect.Target), dk.dkCountActiveDiseases(hitEffect.Target))
				return weaponBaseDamage(sim, hitEffect, spell) * (1.0 + activeDiseases*diseaseMulti)
			},
		},
		OnSpellHitDealt: onhit,
	}

	conf := core.SpellConfig{
		ActionID:    HeartStrikeActionID,
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
	if !dk.Talents.HeartStrike {
		return
	}

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
