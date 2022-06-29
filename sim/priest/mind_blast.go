package priest

import (
	"time"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/stats"
)

func (priest *Priest) registerMindBlastSpell() {
	baseCost := 450.0

	priest.MindBlast = priest.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 25375},
		SpellSchool: core.SpellSchoolShadow,

		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost:     baseCost * (1 - 0.05*float64(priest.Talents.FocusedMind)),
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * 1500,
			},
			CD: core.Cooldown{
				Timer:    priest.NewTimer(),
				Duration: time.Second*8 - time.Millisecond*500*time.Duration(priest.Talents.ImprovedMindBlast),
			},
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask: core.ProcMaskSpellDamage,
			BonusSpellHitRating: 0 +
				float64(priest.Talents.ShadowFocus)*2*core.SpellHitRatingPerHitChance +
				float64(priest.Talents.FocusedPower)*2*core.SpellHitRatingPerHitChance,

			BonusSpellCritRating: float64(priest.Talents.ShadowPower) * 3 * core.SpellCritRatingPerCritChance,

			DamageMultiplier: 1 *
				(1 + float64(priest.Talents.Darkness)*0.02) *
				core.TernaryFloat64(priest.Talents.Shadowform, 1.15, 1) *
				core.TernaryFloat64(ItemSetAbsolution.CharacterHasSetBonus(&priest.Character, 4), 1.1, 1),

			ThreatMultiplier: 1 - 0.08*float64(priest.Talents.ShadowAffinity),

			BaseDamage:     core.BaseDamageConfigMagic(711, 752, 0.429),
			OutcomeApplier: priest.OutcomeFuncMagicHitAndCrit(priest.DefaultSpellCritMultiplier()),
		}),
	})
}
