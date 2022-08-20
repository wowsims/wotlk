package encounters

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

var PatchwerkHatefulStrike = TargetAbility{
	ChanceToUse: 1,
	InitialCD:   time.Second * 1,

	MakeSpell: func(target *core.Target) *core.Spell {
		actionID := core.ActionID{SpellID: 59192}

		return target.RegisterSpell(core.SpellConfig{
			ActionID:    actionID,
			SpellSchool: core.SpellSchoolPhysical,
			Flags:       core.SpellFlagMeleeMetrics,

			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    target.NewTimer(),
					Duration: time.Second * 1,
				},
			},

			ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
				ProcMask: core.ProcMaskMeleeMHSpecial,

				DamageMultiplier: 1,

				BaseDamage:     core.BaseDamageConfigRoll(20000, 27000),
				OutcomeApplier: target.OutcomeFuncEnemyMeleeWhite(),
			}),
		})
	},
}
