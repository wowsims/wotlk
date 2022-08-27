package hunter

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (hunter *Hunter) registerSilencingShotSpell() {
	if !hunter.Talents.SilencingShot {
		return
	}
	baseCost := 0.06 * hunter.BaseMana

	hunter.SilencingShot = hunter.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 34490},
		SpellSchool: core.SpellSchoolPhysical,
		Flags:       core.SpellFlagMeleeMetrics,

		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost *
					(1 - 0.03*float64(hunter.Talents.Efficiency)),
			},
			CD: core.Cooldown{
				Timer:    hunter.NewTimer(),
				Duration: time.Second * 20,
			},
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask: core.ProcMaskRangedSpecial,
			DamageMultiplier: 0.5 *
				hunter.markedForDeathMultiplier(),
			ThreatMultiplier: 1,

			BaseDamage:     core.BaseDamageConfigRangedWeapon(hunter.AmmoDamageBonus),
			OutcomeApplier: hunter.OutcomeFuncRangedHitAndCrit(hunter.critMultiplier(true, false, hunter.CurrentTarget)),

			// Add a check for later so we use ASAP when it comes off CD.
			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				core.StartDelayedAction(sim, core.DelayedActionOptions{
					DoAt: sim.CurrentTime + hunter.SilencingShot.CD.Duration,
					OnAction: func(sim *core.Simulation) {
						// Need to check in case Readiness caused a shift in timing.
						if hunter.SilencingShot.IsReady(sim) {
							hunter.SilencingShot.Cast(sim, hunter.CurrentTarget)
						}
					},
				})
			},
		}),
	})
}
