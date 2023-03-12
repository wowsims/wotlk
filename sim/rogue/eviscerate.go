package rogue

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func (rogue *Rogue) registerEviscerate() {
	rogue.Eviscerate = rogue.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 48668},
		SpellSchool:  core.SpellSchoolPhysical,
		ProcMask:     core.ProcMaskMeleeMHSpecial,
		Flags:        core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | rogue.finisherFlags() | SpellFlagColdBlooded,
		MetricSplits: 6,

		EnergyCost: core.EnergyCostOptions{
			Cost:          35,
			Refund:        0.4 * float64(rogue.Talents.QuickRecovery),
			RefundMetrics: rogue.QuickRecoveryMetrics,
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

		BonusCritRating: core.TernaryFloat64(
			rogue.HasMajorGlyph(proto.RogueMajorGlyph_GlyphOfEviscerate), 10*core.CritRatingPerCritChance, 0.0),
		DamageMultiplier: 1 +
			[]float64{0.0, 0.07, 0.14, 0.2}[rogue.Talents.ImprovedEviscerate] +
			0.02*float64(rogue.Talents.FindWeakness) +
			0.03*float64(rogue.Talents.Aggression),
		CritMultiplier:   rogue.MeleeCritMultiplier(false),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			comboPoints := rogue.ComboPoints()
			flatBaseDamage := 127 + 370*float64(comboPoints)
			// tooltip implies 3..7% AP scaling, but testing show it's fixed at 7% (3.4.0.46158)
			apRatio := 0.07 * float64(comboPoints)

			baseDamage := flatBaseDamage +
				254.0*sim.RandomFloat("Eviscerate") +
				apRatio*spell.MeleeAttackPower() +
				spell.BonusWeaponDamage()

			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)

			if result.Landed() {
				rogue.ApplyFinisher(sim, spell)
				rogue.ApplyCutToTheChase(sim)
			} else {
				spell.IssueRefund(sim)
			}

			spell.DealDamage(sim, result)
		},
	})
}
