package shaman

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (shaman *Shaman) newFireNovaSpell() {
	manaCost := 0.22 * shaman.BaseMana

	fireNovaGlyphCDReduction := core.TernaryInt32(shaman.HasMajorGlyph(proto.ShamanMajorGlyph_GlyphOfFireNova), 3, 0)
	impFireNovaCDReduction := 2 * shaman.Talents.ImprovedFireNova
	fireNovaCooldown := time.Second * (10 - fireNovaGlyphCDReduction - impFireNovaCDReduction)

	shaman.FireNova = shaman.RegisterSpell(core.SpellConfig{
		ActionId:    core.ActionID{SpellID: 61657},
		SpellSchool: core.SpellSchoolFire,

		ResourceType: stats.Mana,
		BaseCost:     manaCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: manaCost,
				GCD:  core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    shaman.NewTimer(),
				Duration: fireNovaCooldown,
			},
		},

		ApplyEffects: core.ApplyEffectFuncAOEDamage(shaman.Env, core.SpellEffect{
			ProcMask:             core.ProcMaskSpellDamage,
			BonusSpellHitRating:  float64(shaman.Talents.ElementalPrecision) * 2 * core.SpellHitRatingPerHitChance,
			BonusSpellCritRating: float64(1.0) * 2 * core.CritRatingPerCritChance,

			DamageMultiplier: 1 + float64(shaman.Talents.CallOfFlame)*0.05,
			ThreatMultiplier: 1 - (0.1/3)*float64(shaman.Talents.ElementalPrecision),

			BaseDamage:     core.BaseDamageConfigMagic(893, 997, 0.2142), // FIXME: double check spell coefficients
			OutcomeApplier: core.OutcomeFuncMagicHitAndCrit(shaman.ElementalCritMultiplier()),
		}),
	})
}
