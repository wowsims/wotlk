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
	// Glyph to single target, OR apply to up to 3 targets
	numHits := core.TernaryInt32(glyphedSingleTargetAS, 1, core.MinInt32(3, paladin.Env.GetNumTargets()))
	results := make([]*core.SpellResult, numHits)

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
		CritMultiplier:   paladin.MeleeCritMultiplier(),
		// TODO: Why is this here?
		BonusCritRating:  1,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			constBaseDamage := 1100.0 +
				.07*spell.SpellPower() +
				.07*spell.MeleeAttackPower()

			curTarget := target
			for hitIndex := int32(0); hitIndex < numHits; hitIndex++ {
				baseDamage := constBaseDamage + (1344.0-1100.0)*sim.RandomFloat("Damage Roll")

				results[hitIndex] = spell.CalcDamage(sim, curTarget, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
				curTarget = sim.Environment.NextTargetUnit(curTarget)
			}

			curTarget = target
			for hitIndex := int32(0); hitIndex < numHits; hitIndex++ {
				spell.DealDamage(sim, results[hitIndex])
				curTarget = sim.Environment.NextTargetUnit(curTarget)
			}
		},
	})
}
