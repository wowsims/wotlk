package hunter

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (hunter *Hunter) registerVolleySpell() {
	actionID := core.ActionID{SpellID: 58434}
	baseCost := 0.17 * hunter.BaseMana

	volleyDot := core.NewDot(core.Dot{
		Aura: hunter.RegisterAura(core.Aura{
			Label:    "Volley",
			ActionID: actionID,
		}),
		NumberOfTicks: 6,
		TickLength:    time.Second * 1,
		// TODO: Whats the actual AOE cap?
		TickEffects: core.TickFuncAOESnapshotCapped(hunter.Env, core.SpellEffect{
			ProcMask: core.ProcMaskPeriodicDamage,
			DamageMultiplier: 1 *
				(1 + 0.04*float64(hunter.Talents.Barrage)),

			ThreatMultiplier: 1,

			BaseDamage: core.BaseDamageConfig{
				Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
					rap := hitEffect.RangedAttackPower(spell.Unit) + hitEffect.RangedAttackPowerOnTarget()
					return 353 + rap*0.0837
				},
			},
			OutcomeApplier: hunter.OutcomeFuncRangedHitAndCrit(hunter.critMultiplier(true, false, hunter.CurrentTarget)),
			IsPeriodic:     true,
		}),
	})

	hunter.Volley = hunter.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolArcane,
		Flags:       core.SpellFlagChanneled,

		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost:        baseCost * core.TernaryFloat64(hunter.HasMajorGlyph(proto.HunterMajorGlyph_GlyphOfVolley), 0.8, 1),
				GCD:         core.GCDDefault,
				ChannelTime: time.Second * 6,
			},
			IgnoreHaste: true,

			OnCastComplete: func(sim *core.Simulation, _ *core.Spell) {
				hunter.AutoAttacks.DelayRangedUntil(sim, sim.CurrentTime + time.Second * 6)
			},
		},

		ApplyEffects: core.ApplyEffectFuncDot(volleyDot),
	})
	volleyDot.Spell = hunter.Volley
}
