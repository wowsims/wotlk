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
				Duration: time.Second*20,
			},
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask: core.ProcMaskRangedSpecial,
			DamageMultiplier: 0.5 *
				(1 + 0.01*float64(hunter.Talents.MarkedForDeath)),
			ThreatMultiplier: 1,

			BaseDamage:     core.BaseDamageConfigRangedWeapon(0),
			OutcomeApplier: hunter.OutcomeFuncRangedHitAndCrit(hunter.critMultiplier(true, false, hunter.CurrentTarget)),
		}),
	})
}
