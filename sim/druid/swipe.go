package druid

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

func (druid *Druid) registerSwipeBearSpell() {
	flatBaseDamage := 108.0
	if druid.Equip[core.ItemSlotRanged].ID == 23198 { // Idol of Brutality
		flatBaseDamage += 10
	} else if druid.Equip[core.ItemSlotRanged].ID == 38365 { // Idol of Perspicacious Attacks
		flatBaseDamage += 24
	}

	lbdm := core.TernaryFloat64(druid.HasSetBonus(ItemSetLasherweaveBattlegear, 2), 1.2, 1.0)
	thdm := core.TernaryFloat64(druid.HasSetBonus(ItemSetThunderheartHarness, 4), 1.15, 1.0)
	fidm := 1.0 + 0.1*float64(druid.Talents.FeralInstinct)

	druid.SwipeBear = druid.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 48562},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage,

		RageCost: core.RageCostOptions{
			Cost: 20 - float64(druid.Talents.Ferocity),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return druid.InForm(Bear)
		},

		DamageMultiplier: lbdm * thdm * fidm,
		CritMultiplier:   druid.MeleeCritMultiplier(Bear),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := flatBaseDamage + 0.063*spell.MeleeAttackPower()
			baseDamage *= sim.Encounter.AOECapMultiplier()
			for _, aoeTarget := range sim.Encounter.TargetUnits {
				spell.CalcAndDealDamage(sim, aoeTarget, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
			}
		},
	})
}

func (druid *Druid) registerSwipeCatSpell() {
	weaponMulti := 2.5
	fidm := 1.0 + 0.1*float64(druid.Talents.FeralInstinct)

	druid.SwipeCat = druid.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 62078},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage,

		EnergyCost: core.EnergyCostOptions{
			Cost: 50 - float64(druid.Talents.Ferocity),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return druid.InForm(Cat)
		},

		DamageMultiplier: fidm * weaponMulti,
		CritMultiplier:   druid.MeleeCritMultiplier(Cat),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())
			baseDamage *= sim.Encounter.AOECapMultiplier()
			for _, aoeTarget := range sim.Encounter.TargetUnits {
				spell.CalcAndDealDamage(sim, aoeTarget, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
			}
		},
	})
}

func (druid *Druid) CurrentSwipeCatCost() float64 {
	return druid.SwipeCat.ApplyCostModifiers(druid.SwipeCat.DefaultCast.Cost)
}

func (druid *Druid) IsSwipeSpell(spell *core.Spell) bool {
	return spell == druid.SwipeBear || spell == druid.SwipeCat
}
