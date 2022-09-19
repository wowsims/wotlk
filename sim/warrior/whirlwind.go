package warrior

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (warrior *Warrior) registerWhirlwindSpell() {
	actionID := core.ActionID{SpellID: 1680}
	cost := 25.0 - float64(warrior.Talents.FocusedRage)
	numHits := core.MinInt32(4, warrior.Env.GetNumTargets())

	var ohDamageEffects core.ApplySpellEffects
	if warrior.AutoAttacks.IsDualWielding {
		baseEffectOH := core.SpellEffect{
			ProcMask:       core.ProcMaskMeleeOHSpecial,
			BaseDamage:     core.BaseDamageConfigMeleeWeapon(false, true, 0, true),
			OutcomeApplier: warrior.OutcomeFuncMeleeWeaponSpecialHitAndCrit(warrior.critMultiplier(oh)),
		}

		effects := make([]core.SpellEffect, 0, numHits)
		for i := int32(0); i < numHits; i++ {
			effect := baseEffectOH
			effect.Target = warrior.Env.GetTargetUnit(i)
			effects = append(effects, effect)
		}
		ohDamageEffects = core.ApplyEffectFuncDamageMultiple(effects)

		warrior.WhirlwindOH = warrior.RegisterSpell(core.SpellConfig{
			ActionID:    actionID,
			SpellSchool: core.SpellSchoolPhysical,
			Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagNoOnCastComplete,

			DamageMultiplier: 1 *
				(1 + 0.02*float64(warrior.Talents.UnendingFury) + 0.1*float64(warrior.Talents.ImprovedWhirlwind)) *
				(1 + 0.05*float64(warrior.Talents.DualWieldSpecialization)),
			ThreatMultiplier: 1.25,
		})
	}

	baseEffectMH := core.SpellEffect{
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		BaseDamage:     core.BaseDamageConfigMeleeWeapon(true, true, 0, true),
		OutcomeApplier: warrior.OutcomeFuncMeleeWeaponSpecialHitAndCrit(warrior.critMultiplier(mh)),
	}

	effects := make([]core.SpellEffect, 0, numHits)
	for i := int32(0); i < numHits; i++ {
		effect := baseEffectMH
		effect.Target = warrior.Env.GetTargetUnit(i)
		effects = append(effects, effect)
	}
	mhDamageEffects := core.ApplyEffectFuncDamageMultiple(effects)

	warrior.Whirlwind = warrior.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 1680},
		SpellSchool: core.SpellSchoolPhysical,
		Flags:       core.SpellFlagMeleeMetrics,

		ResourceType: stats.Rage,
		BaseCost:     cost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: cost,
				GCD:  core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: core.TernaryDuration(warrior.HasMajorGlyph(proto.WarriorMajorGlyph_GlyphOfWhirlwind), time.Second*8, time.Second*10),
			},
		},

		DamageMultiplier: 1 *
			(1 + 0.02*float64(warrior.Talents.UnendingFury) + 0.1*float64(warrior.Talents.ImprovedWhirlwind)),
		ThreatMultiplier: 1.25,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			mhDamageEffects(sim, target, spell)
			if warrior.WhirlwindOH != nil {
				ohDamageEffects(sim, target, warrior.WhirlwindOH)
			}
		},
	})
}

func (warrior *Warrior) CanWhirlwind(sim *core.Simulation) bool {
	return warrior.StanceMatches(BerserkerStance) && warrior.CurrentRage() >= warrior.Whirlwind.DefaultCast.Cost && warrior.Whirlwind.IsReady(sim)
}
