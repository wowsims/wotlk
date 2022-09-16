package warrior

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (warrior *Warrior) registerWhirlwindSpell() {
	cost := 25.0 - float64(warrior.Talents.FocusedRage)
	cd := core.TernaryDuration(warrior.HasMajorGlyph(proto.WarriorMajorGlyph_GlyphOfWhirlwind), time.Second*8, time.Second*10)

	dm := 1 + 0.02*float64(warrior.Talents.UnendingFury) + 0.1*float64(warrior.Talents.ImprovedWhirlwind)

	baseEffectMH := core.SpellEffect{
		ProcMask: core.ProcMaskMeleeMHSpecial,

		DamageMultiplier: dm,

		BaseDamage:     core.BaseDamageConfigMeleeWeapon(core.MainHand, true, 0, true),
		OutcomeApplier: warrior.OutcomeFuncMeleeWeaponSpecialHitAndCrit(warrior.critMultiplier(mh)),
	}
	baseEffectOH := core.SpellEffect{
		ProcMask: core.ProcMaskMeleeOHSpecial,

		DamageMultiplier: dm * (1 + 0.05*float64(warrior.Talents.DualWieldSpecialization)),

		BaseDamage:     core.BaseDamageConfigMeleeWeapon(core.OffHand, true, 0, true),
		OutcomeApplier: warrior.OutcomeFuncMeleeWeaponSpecialHitAndCrit(warrior.critMultiplier(oh)),
	}

	numHits := core.MinInt32(4, warrior.Env.GetNumTargets())
	numTotalHits := numHits
	if warrior.AutoAttacks.IsDualWielding {
		numTotalHits *= 2
	}

	effects := make([]core.SpellEffect, 0, numTotalHits)
	for i := int32(0); i < numHits; i++ {
		mhEffect := baseEffectMH
		mhEffect.Target = warrior.Env.GetTargetUnit(i)
		effects = append(effects, mhEffect)

		if warrior.AutoAttacks.IsDualWielding {
			ohEffect := baseEffectOH
			ohEffect.Target = warrior.Env.GetTargetUnit(i)
			effects = append(effects, ohEffect)
		}
	}

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
				Duration: cd,
			},
		},

		ThreatMultiplier: 1.25,

		ApplyEffects: core.ApplyEffectFuncDamageMultiple(effects),
	})
}

func (warrior *Warrior) CanWhirlwind(sim *core.Simulation) bool {
	return warrior.StanceMatches(BerserkerStance) && warrior.CurrentRage() >= warrior.Whirlwind.DefaultCast.Cost && warrior.Whirlwind.IsReady(sim)
}
