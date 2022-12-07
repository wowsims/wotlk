package druid

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (druid *Druid) registerShredSpell() {
	baseCost := 60.0 - 9*float64(druid.Talents.ShreddingAttacks)

	flatDamageBonus := 666 +
		core.TernaryFloat64(druid.Equip[core.ItemSlotRanged].ID == 29390, 88, 0) +
		core.TernaryFloat64(druid.Equip[core.ItemSlotRanged].ID == 40713, 203, 0)

	hasGlyphofShred := druid.HasMajorGlyph(proto.DruidMajorGlyph_GlyphOfShred)
	maxRipTicks := druid.MaxRipTicks()

	druid.Shred = druid.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 48572},
		SpellSchool:  core.SpellSchoolPhysical,
		ProcMask:     core.ProcMaskMeleeMHSpecial,
		Flags:        core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage,
		ResourceType: stats.Energy,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost,
				GCD:  time.Second,
			},
			IgnoreHaste: true,
		},

		DamageMultiplier: 2.25,
		CritMultiplier:   druid.MeleeCritMultiplier(Cat),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := flatDamageBonus/2.25 +
				spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower()) +
				spell.BonusWeaponDamage()

			modifier := 1.0
			if druid.CurrentTarget.HasActiveAuraWithTag(core.BleedDamageAuraTag) {
				modifier += .3
			}
			if druid.AssumeBleedActive || druid.RipDot.IsActive() || druid.RakeDot.IsActive() || druid.LacerateDot.IsActive() {
				modifier *= 1.0 + (0.04 * float64(druid.Talents.RendAndTear))
			}
			baseDamage *= modifier

			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)

			if result.Landed() {
				druid.AddComboPoints(sim, 1, spell.ComboPointMetrics())

				if hasGlyphofShred && druid.RipDot.IsActive() {
					if druid.RipDot.NumberOfTicks < maxRipTicks {
						druid.RipDot.NumberOfTicks += 1
						druid.RipDot.RecomputeAuraDuration()
						druid.RipDot.UpdateExpires(druid.RipDot.ExpiresAt() + time.Second*2)
					}
				}
			} else {
				druid.AddEnergy(sim, spell.CurCast.Cost*0.8, druid.EnergyRefundMetrics)
			}
		},
	})
}

func (druid *Druid) CanShred() bool {
	return !druid.PseudoStats.InFrontOfTarget && druid.CurrentEnergy() >= druid.CurrentShredCost()
}

func (druid *Druid) CurrentShredCost() float64 {
	return druid.Shred.ApplyCostModifiers(druid.Shred.BaseCost)
}
