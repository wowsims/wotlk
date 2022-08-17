package encounters

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

var BrutallusStomp = TargetAbility{
	ChanceToUse: 1,
	InitialCD:   time.Second * 10,

	MakeSpell: func(target *core.Target) *core.Spell {
		actionID := core.ActionID{SpellID: 45185}

		characterTarget := target.Env.Raid.GetPlayerFromUnit(target.CurrentTarget).GetCharacter()

		statDep := characterTarget.NewDynamicMultiplyStat(stats.Armor, 0.5)
		stompDebuff := characterTarget.RegisterAura(core.Aura{
			Label:    "Stomp",
			ActionID: actionID,
			Duration: time.Second * 10,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				aura.Unit.EnableDynamicStatDep(sim, statDep)
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				aura.Unit.DisableDynamicStatDep(sim, statDep)
			},
		})

		return target.RegisterSpell(core.SpellConfig{
			ActionID:    actionID,
			SpellSchool: core.SpellSchoolPhysical,
			Flags:       core.SpellFlagMeleeMetrics,

			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    target.NewTimer(),
					Duration: time.Second * 15,
				},
			},

			ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
				ProcMask: core.ProcMaskMeleeMHSpecial,

				DamageMultiplier: 1,

				BaseDamage:     core.BaseDamageConfigRoll(18850, 21150),
				OutcomeApplier: target.OutcomeFuncEnemyMeleeWhite(),

				OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
					stompDebuff.Activate(sim)
				},
			}),
		})
	},
}
