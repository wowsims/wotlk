package deathknight

import (
	"github.com/wowsims/wotlk/sim/core"
)

var DeathStrikeActionID = core.ActionID{SpellID: 49924}
var DeathStrikeMHOutcome = core.OutcomeHit
var DeathStrikeOHOutcome = core.OutcomeHit

func (dk *Deathknight) newDeathStrikeSpell(isMH bool) *core.Spell {
	bonusBaseDamage := dk.sigilOfAwarenessBonus(dk.DeathStrike)
	weaponBaseDamage := core.BaseDamageFuncMeleeWeapon(core.MainHand, false, 223.0+bonusBaseDamage, 0.75, true)
	if !isMH {
		weaponBaseDamage = core.BaseDamageFuncMeleeWeapon(core.OffHand, false, 223.0+bonusBaseDamage, 0.75*dk.nervesOfColdSteelBonus(), true)
	}

	effect := core.SpellEffect{
		BonusCritRating:  (dk.annihilationCritBonus() + dk.improvedDeathStrikeCritBonus()) * core.CritRatingPerCritChance,
		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		BaseDamage: core.BaseDamageConfig{
			Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
				return weaponBaseDamage(sim, hitEffect, spell) *
					dk.tundraStalkerBonus(hitEffect.Target) *
					dk.rageOfRivendareBonus(hitEffect.Target)
			},
			TargetSpellCoefficient: 1,
		},

		OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if isMH {
				DeathStrikeMHOutcome = spellEffect.Outcome
			} else {
				DeathStrikeOHOutcome = spellEffect.Outcome
			}
		},
	}

	// TODO: might of mograine crit damage bonus!
	dk.threatOfThassarianProcMasks(isMH, &effect, false, true, func(outcomeApplier core.OutcomeApplier) core.OutcomeApplier {
		return outcomeApplier
	})

	return dk.RegisterSpell(core.SpellConfig{
		ActionID:     DeathStrikeActionID.WithTag(core.TernaryInt32(isMH, 1, 2)),
		SpellSchool:  core.SpellSchoolPhysical,
		Flags:        core.SpellFlagMeleeMetrics,
		ApplyEffects: core.ApplyEffectFuncDirectDamage(effect),
	})
}

func (dk *Deathknight) registerDeathStrikeSpell() {
	dk.DeathStrikeMhHit = dk.newDeathStrikeSpell(true)
	dk.DeathStrikeOhHit = dk.newDeathStrikeSpell(false)
	dk.DeathStrike = dk.RegisterSpell(core.SpellConfig{
		ActionID:    DeathStrikeActionID.WithTag(3),
		SpellSchool: core.SpellSchoolPhysical,
		Flags:       core.SpellFlagNoMetrics | core.SpellFlagNoLogs,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				cast.GCD = dk.getModifiedGCD()
			},
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:         core.ProcMaskEmpty,
			ThreatMultiplier: 1,

			OutcomeApplier: dk.OutcomeFuncAlwaysHit(),

			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				dk.threatOfThassarianProc(sim, spellEffect, dk.DeathStrikeMhHit, dk.DeathStrikeOhHit)

				dk.LastCastOutcome = DeathStrikeMHOutcome
				if dk.outcomeEitherWeaponHitOrCrit(DeathStrikeMHOutcome, DeathStrikeOHOutcome) {
					dkSpellCost := dk.DetermineCost(sim, core.DKCastEnum_FU)
					dk.Spend(sim, spell, dkSpellCost)

					amountOfRunicPower := 15.0 + 2.5*float64(dk.Talents.Dirge)
					dk.AddRunicPower(sim, amountOfRunicPower, spell.RunicPowerMetrics())
				}
			},
		}),
	})
}

func (dk *Deathknight) CanDeathStrike(sim *core.Simulation) bool {
	return dk.CastCostPossible(sim, 0.0, 0, 1, 1) && dk.DeathStrike.IsReady(sim)
}

func (dk *Deathknight) CastDeathStrike(sim *core.Simulation, target *core.Unit) bool {
	if dk.CanDeathStrike(sim) {
		dk.DeathStrike.Cast(sim, target)
		return true
	}
	return false
}
