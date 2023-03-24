package rogue

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

const RuptureEnergyCost = 25.0
const RuptureSpellID = 48672

func (rogue *Rogue) registerRupture() {
	glyphTicks := core.TernaryInt32(rogue.HasMajorGlyph(proto.RogueMajorGlyph_GlyphOfRupture), 2, 0)

	rogue.Rupture = rogue.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: RuptureSpellID},
		SpellSchool:  core.SpellSchoolPhysical,
		ProcMask:     core.ProcMaskMeleeMHSpecial,
		Flags:        core.SpellFlagMeleeMetrics | rogue.finisherFlags(),
		MetricSplits: 6,

		EnergyCost: core.EnergyCostOptions{
			Cost:          RuptureEnergyCost,
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

		DamageMultiplier: 1 +
			0.15*float64(rogue.Talents.BloodSpatter) +
			0.02*float64(rogue.Talents.FindWeakness) +
			core.TernaryFloat64(rogue.HasSetBonus(ItemSetBonescythe, 2), 0.1, 0) +
			core.TernaryFloat64(rogue.HasSetBonus(ItemSetTerrorblade, 4), 0.2, 0) +
			0.1*float64(rogue.Talents.SerratedBlades),
		CritMultiplier:   rogue.MeleeCritMultiplier(false),
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Rupture",
				Tag:   RogueBleedTag,
			},
			NumberOfTicks: 0, // Set dynamically
			TickLength:    time.Second * 2,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, _ bool) {
				dot.SnapshotBaseDamage = rogue.RuptureDamage(rogue.ComboPoints())
				attackTable := dot.Spell.Unit.AttackTables[target.UnitIndex]
				dot.SnapshotCritChance = dot.Spell.PhysicalCritChance(target, attackTable)
				dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(attackTable)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMeleeSpecialHit)
			if result.Landed() {
				numberOfTicks := 3 + rogue.ComboPoints() + glyphTicks
				dot := spell.Dot(target)
				dot.Spell = spell
				dot.NumberOfTicks = numberOfTicks
				dot.MaxStacks = numberOfTicks // slightly hacky; used to determine max extra ticks from Glyph of Backstab
				dot.Apply(sim)
				rogue.ApplyFinisher(sim, spell)
			} else {
				spell.IssueRefund(sim)
			}
			spell.DealOutcome(sim, result)
		},
	})
}

func (rogue *Rogue) RuptureDamage(comboPoints int32) float64 {
	return 127 +
		18*float64(comboPoints) +
		[]float64{0, 0.06 / 4, 0.12 / 5, 0.18 / 6, 0.24 / 7, 0.30 / 8}[comboPoints]*rogue.Rupture.MeleeAttackPower()
}

func (rogue *Rogue) RuptureTicks(comboPoints int32) int32 {
	return 3 + comboPoints + core.TernaryInt32(rogue.HasMajorGlyph(proto.RogueMajorGlyph_GlyphOfRupture), 2, 0)
}

func (rogue *Rogue) RuptureDuration(comboPoints int32) time.Duration {
	return time.Duration(rogue.RuptureTicks(comboPoints)) * time.Second * 2
}
