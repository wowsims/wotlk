package hunter

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (hunter *Hunter) registerExplosiveShotSpell(timer *core.Timer) {
	if !hunter.Talents.ExplosiveShot {
		return
	}

	actionID := core.ActionID{SpellID: 60053}
	baseCost := 0.07 * hunter.BaseMana

	baseEffect := core.SpellEffect{
		BaseDamage: core.BaseDamageConfig{
			Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
				return core.DamageRoll(sim, 386, 464) +
					0.14*spell.RangedAttackPower(hitEffect.Target)
			},
			TargetSpellCoefficient: 1,
		},
		OutcomeApplier: hunter.OutcomeFuncRangedHitAndCrit(),
	}

	initialEffect := baseEffect
	initialEffect.OnSpellHitDealt = func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
		if spellEffect.Landed() {
			hunter.ExplosiveShotDot.Apply(sim)
		}
	}

	hunter.ExplosiveShot = hunter.RegisterSpell(core.SpellConfig{
		ActionID:     actionID,
		SpellSchool:  core.SpellSchoolFire,
		ProcMask:     core.ProcMaskRangedSpecial,
		Flags:        core.SpellFlagMeleeMetrics,
		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost * (1 - 0.03*float64(hunter.Talents.Efficiency)),
				GCD:  core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    timer,
				Duration: time.Second * 6,
			},
		},

		BonusCritRating: 0 +
			2*core.CritRatingPerCritChance*float64(hunter.Talents.SurvivalInstincts) +
			core.TernaryFloat64(hunter.HasMajorGlyph(proto.HunterMajorGlyph_GlyphOfExplosiveShot), 4*core.CritRatingPerCritChance, 0),
		DamageMultiplierAdditive: 1 +
			.02*float64(hunter.Talents.TNT),
		DamageMultiplier: 1,
		CritMultiplier:   hunter.critMultiplier(true, false, hunter.CurrentTarget),
		ThreatMultiplier: 1,

		ApplyEffects: core.ApplyEffectFuncDirectDamage(initialEffect),
	})

	dotEffect := baseEffect
	dotEffect.IsPeriodic = true

	target := hunter.CurrentTarget
	hunter.ExplosiveShotDot = core.NewDot(core.Dot{
		Spell: hunter.ExplosiveShot,
		Aura: target.RegisterAura(core.Aura{
			Label:    "ExplosiveShot-" + strconv.Itoa(int(hunter.Index)),
			ActionID: actionID,
		}),
		NumberOfTicks: 2,
		TickLength:    time.Second * 1,
		TickEffects:   core.TickFuncSnapshot(target, dotEffect),
	})
}
