package deathknight

import (
	"github.com/wowsims/wotlk/sim/core"
)

var PlagueStrikeActionID = core.ActionID{SpellID: 49921}

func (dk *Deathknight) newPlagueStrikeSpell(isMH bool, onhit func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect)) *core.Spell {
	weaponBaseDamage := core.BaseDamageFuncMeleeWeapon(core.MainHand, true, 378.0, 0.5, true)
	if !isMH {
		weaponBaseDamage = core.BaseDamageFuncMeleeWeapon(core.OffHand, true, 378.0, 0.5*dk.nervesOfColdSteelBonus(), true)
	}

	outbreakBonus := 1.0 + 0.1*float64(dk.Talents.Outbreak)

	effect := core.SpellEffect{
		BonusCritRating:  (dk.annihilationCritBonus() + dk.scourgebornePlateCritBonus() + dk.viciousStrikesCritChanceBonus()) * core.CritRatingPerCritChance,
		DamageMultiplier: outbreakBonus,
		ThreatMultiplier: 1,

		BaseDamage: core.BaseDamageConfig{
			Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
				return weaponBaseDamage(sim, hitEffect, spell) * dk.RoRTSBonus(hitEffect.Target)
			},
			TargetSpellCoefficient: 1,
		},

		OnSpellHitDealt: onhit,
	}

	dk.threatOfThassarianProcMasks(isMH, &effect, false, false, func(outcomeApplier core.OutcomeApplier) core.OutcomeApplier {
		return outcomeApplier
	})

	var cconf core.CastConfig
	if isMH { // offhand doesnt need GCD
		cconf = core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				cast.GCD = dk.getModifiedGCD()
			},
		}
	}

	return dk.RegisterSpell(core.SpellConfig{
		Cast:         cconf,
		ActionID:     PlagueStrikeActionID.WithTag(core.TernaryInt32(isMH, 1, 2)),
		SpellSchool:  core.SpellSchoolPhysical,
		Flags:        core.SpellFlagMeleeMetrics,
		ApplyEffects: core.ApplyEffectFuncDirectDamage(effect),
	})
}

func (dk *Deathknight) registerPlagueStrikeSpell() {
	amountOfRunicPower := 10.0 + 2.5*float64(dk.Talents.Dirge)

	dk.PlagueStrikeMhHit = dk.newPlagueStrikeSpell(true, func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
		if dk.Talents.ThreatOfThassarian > 0 && dk.threatOfThassarianWillProc(sim) {
			dk.PlagueStrikeOhHit.Cast(sim, spellEffect.Target)
		}
		dk.LastCastOutcome = spellEffect.Outcome
		if spellEffect.Outcome.Matches(core.OutcomeLanded) {
			dk.BloodPlagueSpell.Cast(sim, spellEffect.Target)
			if dk.Talents.CryptFever > 0 {
				dk.CryptFeverAura[spellEffect.Target.Index].Activate(sim)
			}
			if dk.Talents.EbonPlaguebringer > 0 {
				dk.EbonPlagueAura[spellEffect.Target.Index].Activate(sim)
			}

			dkSpellCost := dk.DetermineCost(sim, core.DKCastEnum_U)
			dk.Spend(sim, spell, dkSpellCost)

			dk.AddRunicPower(sim, amountOfRunicPower, spell.RunicPowerMetrics())
		}
	})
	dk.PlagueStrikeOhHit = dk.newPlagueStrikeSpell(false, nil)
	dk.PlagueStrike = dk.PlagueStrikeMhHit
}

func (dk *Deathknight) CanPlagueStrike(sim *core.Simulation) bool {
	return dk.CastCostPossible(sim, 0.0, 0, 0, 1) && dk.PlagueStrike.IsReady(sim)
}

func (dk *Deathknight) CastPlagueStrike(sim *core.Simulation, target *core.Unit) bool {
	if dk.CanPlagueStrike(sim) {
		dk.PlagueStrike.Cast(sim, target)
		return true
	}
	return false
}
