package druid

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (druid *Druid) registerRipSpell() {
	actionID := core.ActionID{SpellID: 49800}
	baseCost := 30.0 - core.TernaryFloat64(druid.HasSetBonus(ItemSetLasherweaveBattlegear, 2), 10.0, 0.0)
	refundPercent := 0.4 * float64(druid.Talents.PrimalPrecision)

	ripBaseNumTicks := 6 +
		core.TernaryInt32(druid.HasMajorGlyph(proto.DruidMajorGlyph_GlyphOfRip), 2, 0) +
		core.TernaryInt32(druid.HasSetBonus(ItemSetDreamwalkerBattlegear, 2), 2, 0)

	comboPointCoeff := 93.0
	if druid.Equip[core.ItemSlotRanged].ID == 28372 { // Idol of Feral Shadows
		comboPointCoeff += 7
	} else if druid.Equip[core.ItemSlotRanged].ID == 39757 { // Idol of Worship
		comboPointCoeff += 21
	}

	druid.Rip = druid.RegisterSpell(core.SpellConfig{
		ActionID:     actionID,
		SpellSchool:  core.SpellSchoolPhysical,
		ProcMask:     core.ProcMaskMeleeMHSpecial,
		Flags:        core.SpellFlagMeleeMetrics,
		ResourceType: stats.Energy,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost,
				GCD:  time.Second,
			},
			IgnoreHaste: true,
		},

		BonusCritRating:  core.TernaryFloat64(druid.HasSetBonus(ItemSetMalfurionsBattlegear, 4), 5*core.CritRatingPerCritChance, 0.0),
		DamageMultiplier: 1 + core.TernaryFloat64(druid.HasSetBonus(ItemSetThunderheartHarness, 4), 0.15, 0),
		CritMultiplier:   druid.MeleeCritMultiplier(Cat),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			if druid.ClearcastingAura != nil {
				druid.ClearcastingAura.Deactivate(sim)
			}

			result := spell.CalcOutcome(sim, target, spell.OutcomeMeleeSpecialHit)
			if result.Landed() {
				druid.RipDot.NumberOfTicks = ripBaseNumTicks
				druid.RipDot.Apply(sim)
				druid.SpendComboPoints(sim, spell.ComboPointMetrics())
			} else if refundPercent > 0 {
				druid.AddEnergy(sim, spell.CurCast.Cost*refundPercent, druid.PrimalPrecisionRecoveryMetrics)
			}
			spell.DealOutcome(sim, result)
		},
	})

	druid.RipDot = core.NewDot(core.Dot{
		Spell: druid.Rip,
		Aura: druid.CurrentTarget.RegisterAura(druid.applyRendAndTear(core.Aura{
			Label:    "Rip-" + strconv.Itoa(int(druid.Index)),
			ActionID: actionID,
		})),
		NumberOfTicks: ripBaseNumTicks,
		TickLength:    time.Second * 2,

		OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
			cp := float64(druid.ComboPoints())
			ap := dot.Spell.MeleeAttackPower()

			dot.SnapshotBaseDamage = 36 + comboPointCoeff*cp + 0.01*cp*ap

			if !isRollover {
				attackTable := dot.Spell.Unit.AttackTables[target.UnitIndex]
				dot.SnapshotCritChance = dot.Spell.PhysicalCritChance(target, attackTable)
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
	})
}

func (druid *Druid) MaxRipTicks() int32 {
	base := int32(6)
	t7bonus := core.TernaryInt32(druid.HasSetBonus(ItemSetDreamwalkerBattlegear, 2), 2, 0)
	ripGlyphBonus := core.TernaryInt32(druid.HasMajorGlyph(proto.DruidMajorGlyph_GlyphOfRip), 2, 0)
	shredGlyphBonus := core.TernaryInt32(druid.HasMajorGlyph(proto.DruidMajorGlyph_GlyphOfShred), 3, 0)
	return base + ripGlyphBonus + shredGlyphBonus + t7bonus
}

func (druid *Druid) CanRip() bool {
	return druid.InForm(Cat) && druid.ComboPoints() > 0 && druid.CurrentEnergy() >= druid.CurrentRipCost()
}

func (druid *Druid) CurrentRipCost() float64 {
	return druid.Rip.ApplyCostModifiers(druid.Rip.BaseCost)
}
