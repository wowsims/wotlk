package hunter

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

func (hunter *Hunter) registerSilencingShotSpell() {
	if !hunter.Talents.SilencingShot {
		return
	}

	hunter.SilencingShot = hunter.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 34490},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskRangedSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.06,
			Multiplier: 1 - 0.03*float64(hunter.Talents.Efficiency),
		},
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    hunter.NewTimer(),
				Duration: time.Second * 20,
			},
		},

		DamageMultiplier: 0.5 *
			hunter.markedForDeathMultiplier(),
		CritMultiplier:   hunter.critMultiplier(true, false, false),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := hunter.RangedWeaponDamage(sim, spell.RangedAttackPower(target)) +
				hunter.AmmoDamageBonus +
				spell.BonusWeaponDamage()

			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeRangedHitAndCrit)

			// Add a check for later so we use ASAP when it comes off CD.
			core.StartDelayedAction(sim, core.DelayedActionOptions{
				DoAt: sim.CurrentTime + hunter.SilencingShot.CD.Duration,
				OnAction: func(sim *core.Simulation) {
					// Need to check in case Readiness caused a shift in timing.
					if hunter.SilencingShot.IsReady(sim) && hunter.Hardcast.Expires <= sim.CurrentTime {
						hunter.SilencingShot.Cast(sim, target)
					}
				},
			})
		},
	})
}
