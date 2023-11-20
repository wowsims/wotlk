package hunter

import (
	"github.com/wowsims/classic/sim/core"
)

func (hunter *Hunter) registerScorpidStingSpell() {
	hunter.ScorpidStingAuras = hunter.NewEnemyAuraArray(core.ScorpidStingAura)

	hunter.ScorpidSting = hunter.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 3043},
		SpellSchool: core.SpellSchoolNature,
		ProcMask:    core.ProcMaskRangedSpecial,
		Flags:       core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.09,
			Multiplier: 1 - 0.03*float64(hunter.Talents.Efficiency),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true, // Hunter GCD is locked at 1.5s
		},

		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeRangedHit)
			if result.Landed() {
				hunter.ScorpidStingAuras.Get(target).Activate(sim)
			}
			spell.DealOutcome(sim, result)
		},

		RelatedAuras: []core.AuraArray{hunter.ScorpidStingAuras},
	})
}
