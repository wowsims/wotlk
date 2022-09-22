package paladin

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (paladin *Paladin) registerAvengersShieldSpell() {
	baseCost := paladin.BaseMana * 0.26
	glyphedSingleTargetAS := paladin.HasMajorGlyph(proto.PaladinMajorGlyph_GlyphOfAvengerSShield)

	baseEffectMH := core.SpellEffect{
		BaseDamage: core.BaseDamageConfig{
			Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
				damage := 1100.0 +
					(1344.0-1100.0)*sim.RandomFloat("Damage Roll") +
					.07*spell.SpellPower() +
					.07*spell.MeleeAttackPower()
				return damage
			},
		},
		// TODO: Check if it uses spellhit/crit or something crazy (probably not!)
		OutcomeApplier: paladin.OutcomeFuncMeleeSpecialHitAndCrit(paladin.MeleeCritMultiplier()),
	}

	// Glyph to single target, OR apply to up to 3 targets
	numHits := core.TernaryInt32(glyphedSingleTargetAS, 1, core.MinInt32(3, paladin.Env.GetNumTargets()))
	effects := make([]core.SpellEffect, 0, numHits)
	for i := int32(0); i < numHits; i++ {
		mhEffect := baseEffectMH
		mhEffect.Target = paladin.Env.GetTargetUnit(i)
		effects = append(effects, mhEffect)
	}

	paladin.AvengersShield = paladin.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 48827},
		SpellSchool:  core.SpellSchoolHoly,
		ProcMask:     core.ProcMaskMeleeMHSpecial,
		Flags:        core.SpellFlagMeleeMetrics,
		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost * (1 - 0.02*float64(paladin.Talents.Benediction)),
				GCD:  core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    paladin.NewTimer(),
				Duration: time.Second * 30,
			},
		},

		DamageMultiplier: core.TernaryFloat64(glyphedSingleTargetAS, 2, 1),
		// TODO: Why is this here?
		BonusCritRating:  1,
		ThreatMultiplier: 1,

		ApplyEffects: core.ApplyEffectFuncDamageMultiple(effects),
	})
}
