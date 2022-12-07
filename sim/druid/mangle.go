package druid

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (druid *Druid) registerMangleBearSpell() {
	druid.MangleAura = core.MangleAura(druid.CurrentTarget)

	if !druid.Talents.Mangle {
		return
	}

	cost := 20.0 - float64(druid.Talents.Ferocity)
	refundAmount := cost * 0.8
	durReduction := (0.5) * float64(druid.Talents.ImprovedMangle)
	glyphBonus := core.TernaryFloat64(druid.HasMajorGlyph(proto.DruidMajorGlyph_GlyphOfMangle), 1.1, 1.0)

	druid.MangleBear = druid.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 48564},
		SpellSchool:  core.SpellSchoolPhysical,
		ProcMask:     core.ProcMaskMeleeMHSpecial,
		Flags:        core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage,
		ResourceType: stats.Rage,
		BaseCost:     cost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: cost,
				GCD:  core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    druid.NewTimer(),
				Duration: time.Duration(float64(time.Second) * (6 - durReduction)),
			},
		},

		DamageMultiplier: (1 + 0.1*float64(druid.Talents.SavageFury)) * 1.15 * glyphBonus,
		CritMultiplier:   druid.MeleeCritMultiplier(Bear),
		ThreatMultiplier: (1.5 / 1.15) *
			core.TernaryFloat64(druid.InForm(Bear) && druid.HasSetBonus(ItemSetThunderheartHarness, 2), 1.15, 1),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 299/1.15 +
				spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower()) +
				spell.BonusWeaponDamage()

			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)

			if result.Landed() {
				druid.MangleAura.Activate(sim)
			} else {
				druid.AddRage(sim, refundAmount, druid.RageRefundMetrics)
			}

			if druid.BerserkAura.IsActive() {
				spell.CD.Reset()
			}
		},
	})
}

func (druid *Druid) registerMangleCatSpell() {
	druid.MangleAura = core.MangleAura(druid.CurrentTarget)

	if !druid.Talents.Mangle {
		return
	}

	cost := 45.0 - (2.0 * float64(druid.Talents.ImprovedMangle)) - float64(druid.Talents.Ferocity) - core.TernaryFloat64(druid.HasSetBonus(ItemSetThunderheartHarness, 2), 5.0, 0)
	glyphBonus := core.TernaryFloat64(druid.HasMajorGlyph(proto.DruidMajorGlyph_GlyphOfMangle), 1.1, 1.0)

	druid.MangleCat = druid.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 48566},
		SpellSchool:  core.SpellSchoolPhysical,
		ProcMask:     core.ProcMaskMeleeMHSpecial,
		Flags:        core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage,
		ResourceType: stats.Energy,
		BaseCost:     cost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: cost,
				GCD:  time.Second,
			},
			IgnoreHaste: true,
		},

		DamageMultiplier: (1 + 0.1*float64(druid.Talents.SavageFury)) * 2.0 * glyphBonus,
		CritMultiplier:   druid.MeleeCritMultiplier(Cat),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 566/2.0 +
				spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower()) +
				spell.BonusWeaponDamage()

			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)

			if result.Landed() {
				druid.AddComboPoints(sim, 1, spell.ComboPointMetrics())
				druid.MangleAura.Activate(sim)
			} else {
				druid.AddEnergy(sim, spell.CurCast.Cost*0.8, druid.EnergyRefundMetrics)
			}
		},
	})
}

func (druid *Druid) CanMangleBear(sim *core.Simulation) bool {
	return druid.MangleBear != nil && druid.InForm(Bear) && (druid.CurrentRage() >= druid.MangleBear.DefaultCast.Cost || druid.ClearcastingAura.IsActive()) && druid.MangleBear.IsReady(sim)
}

func (druid *Druid) CanMangleCat() bool {
	return druid.MangleCat != nil && druid.InForm(Cat) && druid.CurrentEnergy() >= druid.CurrentMangleCatCost()
}

func (druid *Druid) CurrentMangleCatCost() float64 {
	return druid.MangleCat.ApplyCostModifiers(druid.MangleCat.BaseCost)
}

func (druid *Druid) IsMangle(spell *core.Spell) bool {
	if druid.MangleBear != nil && druid.MangleBear == spell {
		return true
	} else if druid.MangleCat != nil && druid.MangleCat == spell {
		return true
	}
	return false
}
