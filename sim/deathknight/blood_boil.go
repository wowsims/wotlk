package deathknight

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

var BloodBoilActionID = core.ActionID{SpellID: 49941}

func (dk *Deathknight) registerBloodBoilSpell() {
	// TODO: Handle blood boil correctly -
	//  There is no refund and you only get RP on at least one of the effects hitting.
	rs := &RuneSpell{
		Refundable: true,
	}
	baseCost := core.NewRuneCost(10, 1, 0, 0, 0)
	dk.BloodBoil = dk.RegisterSpell(rs, core.SpellConfig{
		ActionID:     BloodBoilActionID,
		SpellSchool:  core.SpellSchoolShadow,
		ProcMask:     core.ProcMaskSpellDamage,
		ResourceType: stats.RunicPower,
		BaseCost:     float64(baseCost),

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:  core.GCDDefault,
				Cost: float64(baseCost),
			},
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				cast.GCD = dk.GetModifiedGCD()
			},
		},

		DamageMultiplier: dk.bloodyStrikesBonus(dk.BloodBoil),
		CritMultiplier:   dk.bonusCritMultiplier(dk.Talents.MightOfMograine),
		ThreatMultiplier: 1.0,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			for _, aoeTarget := range sim.Encounter.Targets {
				aoeUnit := &aoeTarget.Unit

				baseDamage := (sim.Roll(180, 220) + 0.06*dk.getImpurityBonus(spell)) * dk.RoRTSBonus(aoeUnit) * core.TernaryFloat64(dk.DiseasesAreActive(aoeUnit), 1.5, 1.0)
				baseDamage *= sim.Encounter.AOECapMultiplier()

				result := spell.CalcAndDealDamage(sim, aoeUnit, baseDamage, spell.OutcomeMagicHitAndCrit)

				if aoeUnit == dk.CurrentTarget {
					rs.OnResult(sim, result)
					dk.LastOutcome = result.Outcome
				}
			}
		},
	}, func(sim *core.Simulation) bool {
		return dk.CastCostPossible(sim, 0.0, 1, 0, 0) && dk.BloodBoil.IsReady(sim)
	})
}
