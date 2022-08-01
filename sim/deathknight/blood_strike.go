package deathknight

import (
	"github.com/wowsims/wotlk/sim/core"
)

var BloodStrikeActionID = core.ActionID{SpellID: 49930}

func (dk *Deathknight) newBloodStrikeSpell(isMH bool, onhit func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect)) *core.Spell {
	weaponBaseDamage := core.BaseDamageFuncMeleeWeapon(core.MainHand, true, 764.0, 0.4, true)
	if !isMH {
		weaponBaseDamage = core.BaseDamageFuncMeleeWeapon(core.OffHand, true, 764.0, 0.4*dk.nervesOfColdSteelBonus(), true)
	}

	diseaseMulti := dk.diseaseMultiplier(0.125)

	effect := core.SpellEffect{
		BonusCritRating:  (dk.subversionCritBonus() + dk.annihilationCritBonus()) * core.CritRatingPerCritChance,
		DamageMultiplier: dk.bloodOfTheNorthCoeff() * dk.thassariansPlateDamageBonus(),
		ThreatMultiplier: 1,

		BaseDamage: core.BaseDamageConfig{
			Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
				return weaponBaseDamage(sim, hitEffect, spell) *
					(1.0 + dk.countActiveDiseases(hitEffect.Target)*diseaseMulti) * dk.RoRTSBonus(hitEffect.Target)
			},
			TargetSpellCoefficient: 1,
		},
		OnSpellHitDealt: onhit,
	}

	dk.threatOfThassarianProcMasks(isMH, &effect, true, true, func(outcomeApplier core.OutcomeApplier) core.OutcomeApplier {
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
		ActionID:     BloodStrikeActionID.WithTag(core.TernaryInt32(isMH, 1, 2)),
		SpellSchool:  core.SpellSchoolPhysical,
		Flags:        core.SpellFlagMeleeMetrics,
		ApplyEffects: core.ApplyEffectFuncDirectDamage(effect),
	})
}

func (dk *Deathknight) registerBloodStrikeSpell() {
	dk.BloodStrikeMhHit = dk.newBloodStrikeSpell(true, func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
		if dk.Talents.ThreatOfThassarian > 0 && dk.threatOfThassarianWillProc(sim) {
			dk.BloodStrikeOhHit.Cast(sim, spellEffect.Target)
		}
		dk.LastCastOutcome = spellEffect.Outcome

		if spellEffect.Outcome.Matches(core.OutcomeLanded) {
			dkSpellCost := dk.DetermineCost(sim, core.DKCastEnum_B)
			if !dk.bloodOfTheNorthProc(sim, spell, dkSpellCost) {
				if !dk.reapingProc(sim, spell, dkSpellCost) {
					dk.Spend(sim, spell, dkSpellCost)
				}
			}

			if dk.DesolationAura != nil {
				dk.DesolationAura.Activate(sim)
			}

			// Gain at the end, to take into account previous effects for callback
			dk.AddRunicPower(sim, 10.0, spell.RunicPowerMetrics())
		}
	})
	dk.BloodStrikeOhHit = dk.newBloodStrikeSpell(false, nil)
	dk.BloodStrike = dk.BloodStrikeMhHit
}

func (dk *Deathknight) CanBloodStrike(sim *core.Simulation) bool {
	return dk.CastCostPossible(sim, 0.0, 1, 0, 0) && dk.BloodStrike.IsReady(sim)
}

func (dk *Deathknight) CastBloodStrike(sim *core.Simulation, target *core.Unit) bool {
	if dk.CanBloodStrike(sim) {
		dk.BloodStrike.Cast(sim, target)
		return true
	}
	return false
}
