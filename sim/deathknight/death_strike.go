package deathknight

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

// TODO: Cleanup death strike the same way we did for plague strike
var DeathStrikeActionID = core.ActionID{SpellID: 49924}
var DeathStrikeMHOutcome = core.OutcomeMiss
var DeathStrikeOHOutcome = core.OutcomeMiss

func (dk *Deathknight) newDeathStrikeSpell(isMH bool) *RuneSpell {
	bonusBaseDamage := dk.sigilOfAwarenessBonus(dk.DeathStrike)
	weaponBaseDamage := core.BaseDamageFuncMeleeWeapon(core.MainHand, true, 297.0+bonusBaseDamage, 0.75, true)
	if !isMH {
		weaponBaseDamage = core.BaseDamageFuncMeleeWeapon(core.OffHand, true, 297.0+bonusBaseDamage, 0.75*dk.nervesOfColdSteelBonus(), true)
	}

	effect := core.SpellEffect{
		BonusCritRating:  (dk.annihilationCritBonus() + dk.improvedDeathStrikeCritBonus()) * core.CritRatingPerCritChance,
		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		BaseDamage: core.BaseDamageConfig{
			Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
				return weaponBaseDamage(sim, hitEffect, spell) * dk.RoRTSBonus(hitEffect.Target)
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

	return dk.RegisterSpell(nil, core.SpellConfig{
		ActionID:     DeathStrikeActionID.WithTag(core.TernaryInt32(isMH, 1, 2)),
		SpellSchool:  core.SpellSchoolPhysical,
		Flags:        core.SpellFlagMeleeMetrics,
		ApplyEffects: core.ApplyEffectFuncDirectDamage(effect),
	})
}

func (dk *Deathknight) registerDeathStrikeSpell() {
	dk.DeathStrikeMhHit = dk.newDeathStrikeSpell(true)
	dk.DeathStrikeOhHit = dk.newDeathStrikeSpell(false)

	baseCost := float64(core.NewRuneCost(uint8(15.0+2.5*float64(dk.Talents.Dirge)), 0, 1, 1, 0))
	rs := &RuneSpell{}
	dk.DeathStrike = dk.RegisterSpell(rs, core.SpellConfig{
		ActionID:     DeathStrikeActionID.WithTag(3),
		SpellSchool:  core.SpellSchoolPhysical,
		Flags:        core.SpellFlagNoMetrics | core.SpellFlagNoLogs,
		ResourceType: stats.RunicPower,
		BaseCost:     baseCost,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:  core.GCDDefault,
				Cost: baseCost,
			},
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				cast.GCD = dk.getModifiedGCD()
			},
		},

		ApplyEffects: dk.withRuneRefund(rs, core.SpellEffect{
			ProcMask:         core.ProcMaskEmpty,
			ThreatMultiplier: 1,

			OutcomeApplier: dk.OutcomeFuncAlwaysHit(),

			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				dk.threatOfThassarianProc(sim, spellEffect, dk.DeathStrikeMhHit, dk.DeathStrikeOhHit)
				dk.LastCastOutcome = DeathStrikeMHOutcome
			},
		}, false),
	})
}

func (dk *Deathknight) CanDeathStrike(sim *core.Simulation) bool {
	return dk.CastCostPossible(sim, 0.0, 0, 1, 1) && dk.DeathStrike.IsReady(sim)
}

func (dk *Deathknight) CastDeathStrike(sim *core.Simulation, target *core.Unit) bool {
	if !dk.DeathStrike.IsReady(sim) {
		return false
	}
	return dk.DeathStrike.Cast(sim, target)
}
