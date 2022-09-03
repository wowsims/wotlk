package priest

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (priest *Priest) registerGreaterHealSpell() {
	baseCost := .32 * priest.BaseMana

	priest.GreaterHeal = priest.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 42842},
		SpellSchool: core.SpellSchoolHoly,

		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost,

				GCD:      core.GCDDefault,
				CastTime: time.Second*3 - time.Millisecond*100*time.Duration(priest.Talents.HolySpecialization),
			},
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			IsHealing: true,
			ProcMask:  core.ProcMaskSpellDamage,

			DamageMultiplier: 1,
			ThreatMultiplier: 1 - []float64{0, .07, .14, .20}[priest.Talents.SilentResolve],

			BaseDamage:     core.BaseDamageConfigMagic(3980, 4621, 1.6114+0.08*float64(priest.Talents.EmpoweredHealing)),
			OutcomeApplier: priest.OutcomeFuncMagicCrit(priest.DefaultSpellCritMultiplier()),
		}),
	})
}
