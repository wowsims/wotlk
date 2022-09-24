package paladin

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (paladin *Paladin) registerHammerOfTheRighteousSpell() {
	baseCost := paladin.BaseMana * 0.06

	baseEffectMH := core.SpellEffect{
		BaseDamage: core.BaseDamageConfig{
			Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
				damage := spell.Unit.AutoAttacks.MH.CalculateAverageWeaponDamage(spell.MeleeAttackPower())
				speed := spell.Unit.AutoAttacks.MH.SwingSpeed
				return (damage / speed) * 4
			},
		},
		OutcomeApplier: paladin.OutcomeFuncMeleeSpecialHitAndCrit(),
	}

	numHits := core.MinInt32(core.TernaryInt32(paladin.HasMajorGlyph(proto.PaladinMajorGlyph_GlyphOfHammerOfTheRighteous), 4, 3), paladin.Env.GetNumTargets())
	effects := make([]core.SpellEffect, 0, numHits)
	for i := int32(0); i < numHits; i++ {
		mhEffect := baseEffectMH
		mhEffect.Target = paladin.Env.GetTargetUnit(i)
		effects = append(effects, mhEffect)
	}

	paladin.HammerOfTheRighteous = paladin.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 53595},
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
				Duration: time.Second * 6,
			},
		},

		DamageMultiplierAdditive: 1 + paladin.getItemSetRedemptionPlateBonus2() + paladin.getItemSetT9PlateBonus2() + paladin.getItemSetLightswornPlateBonus2(),
		CritMultiplier:           paladin.MeleeCritMultiplier(),
		ThreatMultiplier:         1,

		ApplyEffects: core.ApplyEffectFuncDamageMultiple(effects),
	})
}
