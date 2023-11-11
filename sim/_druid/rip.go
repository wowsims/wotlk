package druid

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func (druid *Druid) registerRipSpell() {
	ripBaseNumTicks := 6 +
		core.TernaryInt32(druid.HasMajorGlyph(proto.DruidMajorGlyph_GlyphOfRip), 2, 0) +
		core.TernaryInt32(druid.HasSetBonus(ItemSetDreamwalkerBattlegear, 2), 2, 0)

	comboPointCoeff := 93.0
	if druid.Ranged().ID == 28372 { // Idol of Feral Shadows
		comboPointCoeff += 7
	} else if druid.Ranged().ID == 39757 { // Idol of Worship
		comboPointCoeff += 21
	}

	druid.Rip = druid.RegisterSpell(Cat, core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 49800},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagAPL,

		EnergyCost: core.EnergyCostOptions{
			Cost:          30 - core.TernaryFloat64(druid.HasSetBonus(ItemSetLasherweaveBattlegear, 2), 10, 0),
			Refund:        0.4 * float64(druid.Talents.PrimalPrecision),
			RefundMetrics: druid.PrimalPrecisionRecoveryMetrics,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return druid.ComboPoints() > 0
		},

		BonusCritRating:  core.TernaryFloat64(druid.HasSetBonus(ItemSetMalfurionsBattlegear, 4), 5*core.CritRatingPerCritChance, 0.0),
		DamageMultiplier: 1 + core.TernaryFloat64(druid.HasSetBonus(ItemSetThunderheartHarness, 4), 0.15, 0),
		CritMultiplier:   druid.MeleeCritMultiplier(Cat),
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: druid.applyRendAndTear(core.Aura{
				Label: "Rip",
			}),
			NumberOfTicks: ripBaseNumTicks,
			TickLength:    time.Second * 2,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				cp := float64(druid.ComboPoints())
				ap := dot.Spell.MeleeAttackPower()

				dot.SnapshotBaseDamage = 36 + comboPointCoeff*cp + 0.01*cp*ap

				if !isRollover {
					attackTable := dot.Spell.Unit.AttackTables[target.UnitIndex]
					dot.SnapshotCritChance = dot.Spell.PhysicalCritChance(attackTable)
					dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(attackTable)
				}
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				if druid.Talents.PrimalGore {
					dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)
				} else {
					dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.Spell.OutcomeAlwaysHit)
				}
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMeleeSpecialHit)
			if result.Landed() {
				spell.SpellMetrics[target.UnitIndex].Hits--
				dot := spell.Dot(target)
				dot.NumberOfTicks = ripBaseNumTicks
				dot.Apply(sim)
				druid.SpendComboPoints(sim, spell.ComboPointMetrics())
			} else {
				spell.IssueRefund(sim)
			}
			spell.DealOutcome(sim, result)
		},
	})
}

func (druid *Druid) MaxRipTicks() int32 {
	base := int32(6)
	t7bonus := core.TernaryInt32(druid.HasSetBonus(ItemSetDreamwalkerBattlegear, 2), 2, 0)
	ripGlyphBonus := core.TernaryInt32(druid.HasMajorGlyph(proto.DruidMajorGlyph_GlyphOfRip), 2, 0)
	shredGlyphBonus := core.TernaryInt32(druid.HasMajorGlyph(proto.DruidMajorGlyph_GlyphOfShred), 3, 0)
	return base + ripGlyphBonus + shredGlyphBonus + t7bonus
}

func (druid *Druid) CurrentRipCost() float64 {
	return druid.Rip.ApplyCostModifiers(druid.Rip.DefaultCast.Cost)
}
