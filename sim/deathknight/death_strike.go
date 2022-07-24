package deathknight

import (
	"github.com/wowsims/wotlk/sim/core"
)

var DeathStrikeActionID = core.ActionID{SpellID: 49924}
var DeathStrikeMHOutcome = core.OutcomeHit
var DeathStrikeOHOutcome = core.OutcomeHit

func (deathKnight *Deathknight) newDeathStrikeSpell(isMH bool) *core.Spell {
	effect := core.SpellEffect{
		BonusCritRating:  (deathKnight.annihilationCritBonus() + deathKnight.improvedDeathStrikeCritBonus()) * core.CritRatingPerCritChance,
		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		BaseDamage: core.BaseDamageConfig{
			Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
				bonusBaseDamage := deathKnight.sigilOfAwarenessBonus(deathKnight.DeathStrike)
				weaponBaseDamage := core.BaseDamageFuncMeleeWeapon(core.MainHand, false, 223.0+bonusBaseDamage, 0.75, true)
				if !isMH {
					weaponBaseDamage = core.BaseDamageFuncMeleeWeapon(core.OffHand, false, 223.0+bonusBaseDamage, 0.75*deathKnight.nervesOfColdSteelBonus(), true)
				}
				return weaponBaseDamage(sim, hitEffect, spell) *
					deathKnight.tundraStalkerBonus(hitEffect.Target) *
					deathKnight.rageOfRivendareBonus(hitEffect.Target)
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
	deathKnight.threatOfThassarianProcMasks(isMH, &effect, false, true, func(outcomeApplier core.OutcomeApplier) core.OutcomeApplier {
		return outcomeApplier
	})

	return deathKnight.RegisterSpell(core.SpellConfig{
		ActionID:     DeathStrikeActionID.WithTag(core.TernaryInt32(isMH, 1, 2)),
		SpellSchool:  core.SpellSchoolPhysical,
		Flags:        core.SpellFlagMeleeMetrics,
		ApplyEffects: core.ApplyEffectFuncDirectDamage(effect),
	})
}

func (deathKnight *Deathknight) registerDeathStrikeSpell() {
	deathKnight.DeathStrikeMhHit = deathKnight.newDeathStrikeSpell(true)
	deathKnight.DeathStrikeOhHit = deathKnight.newDeathStrikeSpell(false)
	deathKnight.DeathStrike = deathKnight.RegisterSpell(core.SpellConfig{
		ActionID:    DeathStrikeActionID.WithTag(3),
		SpellSchool: core.SpellSchoolPhysical,
		Flags:       core.SpellFlagNoMetrics | core.SpellFlagNoLogs,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				cast.GCD = deathKnight.getModifiedGCD()
			},
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:         core.ProcMaskEmpty,
			ThreatMultiplier: 1,

			OutcomeApplier: deathKnight.OutcomeFuncAlwaysHit(),

			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				deathKnight.threatOfThassarianProc(sim, spellEffect, deathKnight.DeathStrikeMhHit, deathKnight.DeathStrikeOhHit)

				deathKnight.LastCastOutcome = DeathStrikeMHOutcome
				if deathKnight.outcomeEitherWeaponHitOrCrit(DeathStrikeMHOutcome, DeathStrikeOHOutcome) {
					dkSpellCost := deathKnight.DetermineOptimalCost(sim, 0, 1, 1)
					deathKnight.Spend(sim, spell, dkSpellCost)

					amountOfRunicPower := 15.0 + 2.5*float64(deathKnight.Talents.Dirge)
					deathKnight.AddRunicPower(sim, amountOfRunicPower, spell.RunicPowerMetrics())
				}
			},
		}),
	})
}

func (deathKnight *Deathknight) CanDeathStrike(sim *core.Simulation) bool {
	return deathKnight.CastCostPossible(sim, 0.0, 0, 1, 1) && deathKnight.DeathStrike.IsReady(sim)
}

func (deathKnight *Deathknight) CastDeathStrike(sim *core.Simulation, target *core.Unit) bool {
	if deathKnight.CanDeathStrike(sim) {
		deathKnight.DeathStrike.Cast(sim, target)
		return true
	}
	return false
}
