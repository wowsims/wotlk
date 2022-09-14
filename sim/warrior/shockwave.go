package warrior

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (warrior *Warrior) registerShockwaveSpell() {
	cost := 15.0 - float64(warrior.Talents.FocusedRage)
	refundAmount := cost * 0.8
	cd := 20*time.Second - core.TernaryDuration(warrior.HasMajorGlyph(proto.WarriorMajorGlyph_GlyphOfShockwave), 3*time.Second, 0)

	warrior.Shockwave = warrior.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 46968},
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

		ApplyEffects: core.ApplyEffectFuncAOEDamageCapped(warrior.Env, core.SpellEffect{
			ProcMask:         core.ProcMaskRanged,
			DamageMultiplier: 1.0,
			ThreatMultiplier: 1.0,

			BaseDamage: core.BaseDamageConfig{
				Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
					return warrior.attackPowerMultiplier(hitEffect, spell.Unit, 0.75)
				},
				TargetSpellCoefficient: 1,
			},
			OutcomeApplier: warrior.OutcomeFuncMeleeSpecialHitAndCrit(warrior.critMultiplier(none)),
			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if !spellEffect.Landed() {
					warrior.AddRage(sim, refundAmount, warrior.RageRefundMetrics)
				}
			},
		}),
	})
}

func (warrior *Warrior) CanShockwave(sim *core.Simulation) bool {
	return warrior.StanceMatches(DefensiveStance) &&
		warrior.CurrentRage() >= warrior.Shockwave.DefaultCast.Cost &&
		warrior.Shockwave.IsReady(sim)
}
