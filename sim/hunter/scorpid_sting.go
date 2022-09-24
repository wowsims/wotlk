package hunter

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (hunter *Hunter) registerScorpidStingSpell() {
	hunter.ScorpidStingAura = core.ScorpidStingAura(hunter.CurrentTarget)

	baseCost := 0.09 * hunter.BaseMana

	hunter.ScorpidSting = hunter.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 3043},
		SpellSchool:  core.SpellSchoolNature,
		ProcMask:     core.ProcMaskRangedSpecial,
		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost * (1 - 0.03*float64(hunter.Talents.Efficiency)),
				GCD:  core.GCDDefault,
			},
			IgnoreHaste: true, // Hunter GCD is locked at 1.5s
		},

		ThreatMultiplier: 1,

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			OutcomeApplier: hunter.OutcomeFuncRangedHit(),
			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spellEffect.Landed() {
					hunter.ScorpidStingAura.Activate(sim)
				}
			},
		}),
	})
}
