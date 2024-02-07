package rogue

import (
	"time"

	"github.com/wowsims/sod/sim/core"
)

func (rogue *Rogue) registerEviscerate() {
	flatDamage := map[int32]float64{
		25: 10,
		40: 22,
		50: 34,
		60: 54,
	}[rogue.Level]

	comboDamageBonus := map[int32]float64{
		25: 31,
		40: 77,
		50: 110,
		60: 170,
	}[rogue.Level]

	damageVariance := map[int32]float64{
		25: 20,
		40: 44,
		50: 68,
		60: 108,
	}[rogue.Level]

	rogue.Eviscerate = rogue.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 48668},
		SpellSchool:  core.SpellSchoolPhysical,
		ProcMask:     core.ProcMaskMeleeMHSpecial,
		Flags:        core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | rogue.finisherFlags() | SpellFlagColdBlooded | core.SpellFlagAPL,
		MetricSplits: 6,

		EnergyCost: core.EnergyCostOptions{
			Cost:   35,
			Refund: 0,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				spell.SetMetricsSplit(spell.Unit.ComboPoints())
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return rogue.ComboPoints() > 0
		},

		DamageMultiplier: 1 +
			[]float64{0.0, 0.05, 0.10, 0.15}[rogue.Talents.ImprovedEviscerate] +
			0.02*float64(rogue.Talents.Aggression),
		CritMultiplier:   rogue.MeleeCritMultiplier(false),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			rogue.BreakStealth(sim)
			comboPoints := rogue.ComboPoints()
			flatBaseDamage := flatDamage + comboDamageBonus*float64(comboPoints)
			// tooltip implies 3..7% AP scaling, but testing shows it's fixed at 7% (3.4.0.46158)
			apRatio := 0.07 * float64(comboPoints)

			baseDamage := flatBaseDamage +
				damageVariance*sim.RandomFloat("Eviscerate") +
				apRatio*spell.MeleeAttackPower() +
				spell.BonusWeaponDamage()

			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)

			if result.Landed() {
				rogue.ApplyFinisher(sim, spell)
			} else {
				spell.IssueRefund(sim)
			}

			spell.DealDamage(sim, result)
		},
	})
}
