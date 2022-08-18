package deathknight

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

// TODO: Cleanup death strike the same way we did for plague strike
var DeathStrikeActionID = core.ActionID{SpellID: 49924}

func (dk *Deathknight) newDeathStrikeSpell(isMH bool, onhit func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect)) *RuneSpell {
	bonusBaseDamage := dk.sigilOfAwarenessBonus(dk.DeathStrike)
	weaponBaseDamage := core.BaseDamageFuncMeleeWeapon(core.MainHand, true, 297.0+bonusBaseDamage, 1.0, 0.75, true)
	if !isMH {
		weaponBaseDamage = core.BaseDamageFuncMeleeWeapon(core.OffHand, true, 297.0+bonusBaseDamage, dk.nervesOfColdSteelBonus(), 0.75, true)
	}

	effect := core.SpellEffect{
		BonusCritRating:  (dk.annihilationCritBonus() + dk.improvedDeathStrikeCritBonus()) * core.CritRatingPerCritChance,
		DamageMultiplier: 1.0 + 0.15*float64(dk.Talents.ImprovedDeathStrike),
		ThreatMultiplier: 1,

		BaseDamage: core.BaseDamageConfig{
			Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
				return weaponBaseDamage(sim, hitEffect, spell) * dk.RoRTSBonus(hitEffect.Target)
			},
			TargetSpellCoefficient: 1,
		},

		OnSpellHitDealt: onhit,
	}

	// TODO: might of mograine crit damage bonus!
	dk.threatOfThassarianProcMasks(isMH, &effect, false, true, func(outcomeApplier core.OutcomeApplier) core.OutcomeApplier {
		return outcomeApplier
	})

	conf := core.SpellConfig{
		ActionID:     DeathStrikeActionID.WithTag(core.TernaryInt32(isMH, 1, 2)),
		SpellSchool:  core.SpellSchoolPhysical,
		Flags:        core.SpellFlagMeleeMetrics,
		ApplyEffects: core.ApplyEffectFuncDirectDamage(effect),
	}

	rs := &RuneSpell{}
	if isMH {
		conf.ResourceType = stats.RunicPower
		conf.BaseCost = float64(core.NewRuneCost(uint8(15.0+2.5*float64(dk.Talents.Dirge)), 0, 1, 1, 0))
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

	if isMH {
		return dk.RegisterSpell(rs, conf, func(sim *core.Simulation) bool {
			return dk.CastCostPossible(sim, 0.0, 0, 1, 1) && dk.DeathStrike.IsReady(sim)
		}, nil)
	} else {
		return dk.RegisterSpell(rs, conf, nil, nil)
	}
}

func (dk *Deathknight) registerDeathStrikeSpell() {
	healthMetrics := dk.NewHealthMetrics(DeathStrikeActionID)

	dk.DeathStrikeOhHit = dk.newDeathStrikeSpell(false, nil)
	dk.DeathStrikeMhHit = dk.newDeathStrikeSpell(true, func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
		dk.LastOutcome = spellEffect.Outcome

		healingAmount := 0.05 * dk.countActiveDiseases(dk.CurrentTarget) * dk.MaxHealth() * (1.0 + 0.5*float64(dk.Talents.ImprovedDeathStrike))
		dk.GainHealth(sim, healingAmount, healthMetrics)

		dk.threatOfThassarianProc(sim, spellEffect, dk.DeathStrikeMhHit, dk.DeathStrikeOhHit)
	})
	dk.DeathStrike = dk.DeathStrikeMhHit
}
