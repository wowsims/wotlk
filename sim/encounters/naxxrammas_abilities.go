package encounters

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

var Patchwerk10HatefulStrike = TargetAbility{
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
					Duration: time.Second * 2,
				},
			},

			ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
				ProcMask: core.ProcMaskMeleeMHSpecial,

				DamageMultiplier: 1,

				BaseDamage:     core.BaseDamageConfigRoll(2.2*27750, 2.2*32250),
				OutcomeApplier: target.OutcomeFuncEnemyMeleeWhite(),
			}),
		})
	},
}

var Patchwerk25HatefulStrike = TargetAbility{
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
					Duration: time.Second * 2,
				},
			},

			ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
				ProcMask: core.ProcMaskMeleeMHSpecial,

				DamageMultiplier: 1,

				BaseDamage:     core.BaseDamageConfigRoll(2.2*27750, 2.2*32250),
				OutcomeApplier: target.OutcomeFuncEnemyMeleeWhite(),
			}),
		})
	},
}
