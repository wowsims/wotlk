package shaman

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (shaman *Shaman) newLavaLashSpell() *core.Spell {
	manaCost := 0.04 * shaman.BaseMana

	return shaman.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 60103},
		SpellSchool: core.SpellSchoolFire,
		Flags:       core.SpellFlagMeleeMetrics,

		ResourceType: stats.Mana,
		BaseCost:     manaCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: manaCost,
				GCD:  core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    shaman.NewTimer(),
				Duration: time.Second * 6,
			},
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask: core.ProcMaskMeleeOHSpecial,

			DamageMultiplier: 1 + core.TernaryFloat64(shaman.SelfBuffs.ImbueOH == proto.ShamanImbue_FlametongueWeapon, core.TernaryFloat64(shaman.HasMajorGlyph(proto.ShamanMajorGlyph_GlyphOfLavaLash), 0.35, 0.25), 0),
			ThreatMultiplier: 1 - (0.1/3)*float64(shaman.Talents.ElementalPrecision),

			BaseDamage:     shaman.AutoAttacks.OHEffect.BaseDamage,
			OutcomeApplier: shaman.OutcomeFuncMeleeSpecialHitAndCrit(shaman.ElementalCritMultiplier()),
		}),
	})
}

func (shaman *Shaman) IsLavaLashCastable(sim *core.Simulation) bool {
	return shaman.LavaLash.IsReady(sim)
}
