package deathknight

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

var BloodBoilActionID = core.ActionID{SpellID: 49941}

func (dk *Deathknight) registerBloodBoilSpell() {

	// TODO: Handle blood boil correctly -
	//  There is no refund and you only get RP on at least one of the effects hitting.
	rs := &RuneSpell{}
	baseCost := core.NewRuneCost(10, 1, 0, 0, 0)
	dk.BloodBoil = dk.RegisterSpell(rs, core.SpellConfig{
		ActionID:     BloodBoilActionID,
		SpellSchool:  core.SpellSchoolShadow,
		ResourceType: stats.RunicPower,
		BaseCost:     float64(baseCost),
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:  core.GCDDefault,
				Cost: float64(baseCost),
			},
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				cast.GCD = dk.getModifiedGCD()
			},
		},

		DamageMultiplier: dk.bloodyStrikesBonus(dk.BloodBoil),
		ThreatMultiplier: 1.0,

		ApplyEffects: dk.withRuneRefund(rs, core.SpellEffect{
			ProcMask: core.ProcMaskSpellDamage,
			BaseDamage: core.BaseDamageConfig{
				Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {

					roll := (220.0-180.0)*sim.RandomFloat("Blood Boil") + 180.0
					return (roll + dk.getImpurityBonus(hitEffect, spell.Unit)*0.06) * dk.RoRTSBonus(hitEffect.Target) * core.TernaryFloat64(dk.DiseasesAreActive(hitEffect.Target), 1.5, 1.0)
				},
			},
			OutcomeApplier: dk.OutcomeFuncMagicHitAndCrit(dk.bonusCritMultiplier(dk.Talents.MightOfMograine)),
			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spellEffect.Target == dk.CurrentTarget {
					dk.LastOutcome = spellEffect.Outcome
				}
			},
		}, true),
	}, func(sim *core.Simulation) bool {
		return dk.CastCostPossible(sim, 0.0, 1, 0, 0) && dk.BloodBoil.IsReady(sim)
	}, nil)
}
