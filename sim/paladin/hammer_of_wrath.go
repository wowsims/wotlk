package paladin

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (paladin *Paladin) registerHammerOfWrathSpell() {
	// From the perspective of max rank.
	baseCost := paladin.BaseMana * 0.12

	paladin.HammerOfWrath = paladin.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 48806},
		SpellSchool: core.SpellSchoolHoly,
		Flags:       core.SpellFlagMeleeMetrics,

		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost *
					(1 - 0.02*float64(paladin.Talents.Benediction)) *
					core.TernaryFloat64(paladin.HasMajorGlyph(proto.PaladinMajorGlyph_GlyphOfHammerOfWrath), 0, 1),

				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    paladin.NewTimer(),
				Duration: time.Second * 6,
			},
		},

		BonusCritRating: 25 * float64(paladin.Talents.SanctifiedWrath) * core.CritRatingPerCritChance,
		DamageMultiplierAdditive: 1 +
			paladin.getItemSetLightbringerBattlegearBonus4() +
			paladin.getItemSetAegisBattlegearBonus2(),
		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask: core.ProcMaskSpellDamage,

			BaseDamage: core.BaseDamageConfig{
				Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
					// TODO: discuss exporting or adding to core for damageRollOptimized hybrid scaling.
					deltaDamage := 1257.0 - 1139.0
					return 1139.0 + deltaDamage*sim.RandomFloat("Damage Roll") +
						.15*spell.SpellPower() +
						.15*spell.MeleeAttackPower()
				},
			},
			OutcomeApplier: paladin.OutcomeFuncMeleeSpecialNoBlockDodgeParry(paladin.MeleeCritMultiplier()),
		}),
	})
}
