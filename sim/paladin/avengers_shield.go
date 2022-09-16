package paladin

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (paladin *Paladin) registerAvengersShieldSpell() {
	baseCost := paladin.BaseMana * 0.26
	baseModifiers := Multiplicative{}
	baseMultiplier := baseModifiers.Get()
	numHits := int32(1)

	scaling := hybridScaling{
		AP: 0.07,
		SP: 0.07,
	}

	baseEffectMH := core.SpellEffect{
		ProcMask: core.ProcMaskMeleeMHSpecial,

		DamageMultiplier: baseMultiplier,
		ThreatMultiplier: 1,
		BonusCritRating:  1,

		BaseDamage: core.BaseDamageConfig{
			Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
				deltaDamage := 1344.0 - 1100.0
				damage := 1100.0 + deltaDamage*sim.RandomFloat("Damage Roll")
				damage += hitEffect.SpellPower(spell.Unit, spell) * scaling.SP
				damage += hitEffect.MeleeAttackPower(spell.Unit) * scaling.AP
				damage *= core.TernaryFloat64(paladin.HasMajorGlyph(41101), 2, 1)
				return damage
			},
		},
		// TODO: Check if it uses spellhit/crit or something crazy (probably not!)
		OutcomeApplier: paladin.OutcomeFuncMeleeSpecialHitAndCrit(paladin.MeleeCritMultiplier()),
	}

	if !paladin.HasMajorGlyph(41101) {
		numHits = core.MinInt32(3, paladin.Env.GetNumTargets())
	}
	
	effects := make([]core.SpellEffect, 0, numHits)
	for i := int32(0); i < numHits; i++ {
		mhEffect := baseEffectMH
		mhEffect.Target = paladin.Env.GetTargetUnit(i)
		effects = append(effects, mhEffect)
	}

	paladin.AvengersShield = paladin.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 48827},
		SpellSchool: core.SpellSchoolHoly,
		Flags:       core.SpellFlagMeleeMetrics,

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

		ApplyEffects: core.ApplyEffectFuncDamageMultiple(effects),
	})
}
