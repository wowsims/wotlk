package warrior

import (
	"github.com/wowsims/classic/sim/core"
)

func (warrior *Warrior) registerExecuteSpell() {
	flatDamage := map[int32]float64{
		25: 125,
		40: 325,
		50: 450,
		60: 600,
	}[warrior.Level]

	convertedRageDamage := map[int32]float64{
		25: 3,
		40: 9,
		50: 12,
		60: 15,
	}[warrior.Level]

	spellID := map[int32]int32{
		25: 5308,
		40: 20660,
		50: 20661,
		60: 20662,
	}[warrior.Level]

	var rageMetrics *core.ResourceMetrics
	warrior.Execute = warrior.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: spellID},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | core.SpellFlagAPL,

		RageCost: core.RageCostOptions{
			Cost:   15 - []float64{0, 2, 5}[warrior.Talents.ImprovedExecute],
			Refund: 0.8,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return sim.IsExecutePhase20()
		},

		DamageMultiplier: 1,
		CritMultiplier:   warrior.critMultiplier(mh),
		ThreatMultiplier: 1.25,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			extraRage := spell.Unit.CurrentRage()
			warrior.SpendRage(sim, extraRage, rageMetrics)
			rageMetrics.Events--

			baseDamage := flatDamage + convertedRageDamage*(extraRage)
			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)

			if !result.Landed() {
				spell.IssueRefund(sim)
			}
		},
	})
	rageMetrics = warrior.Execute.Cost.(*core.RageCost).ResourceMetrics
}
