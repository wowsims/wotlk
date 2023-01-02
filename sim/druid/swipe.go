package druid

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (druid *Druid) registerSwipeBearSpell() {
	cost := 20.0 - float64(druid.Talents.Ferocity)

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
		ActionID:     core.ActionID{SpellID: 48562},
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
		},

		DamageMultiplier: lbdm * thdm * fidm,
		CritMultiplier:   druid.MeleeCritMultiplier(Bear),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			if druid.ClearcastingAura != nil {
				druid.ClearcastingAura.Deactivate(sim)
			}

			baseDamage := flatBaseDamage + 0.063*spell.MeleeAttackPower()
			baseDamage *= sim.Encounter.AOECapMultiplier()
			for _, aoeTarget := range sim.Encounter.Targets {
				spell.CalcAndDealDamage(sim, &aoeTarget.Unit, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
			}
		},
	})
}

func (druid *Druid) registerSwipeCatSpell() {
	cost := 50.0 - float64(druid.Talents.Ferocity)

	weaponMulti := 2.5
	fidm := 1.0 + 0.1*float64(druid.Talents.FeralInstinct)

	druid.SwipeCat = druid.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 62078},
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

		DamageMultiplier: fidm * weaponMulti,
		CritMultiplier:   druid.MeleeCritMultiplier(Cat),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			if druid.ClearcastingAura != nil {
				druid.ClearcastingAura.Deactivate(sim)
			}

			for _, aoeTarget := range sim.Encounter.Targets {
				baseDamage := spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())
				baseDamage *= sim.Encounter.AOECapMultiplier()
				spell.CalcAndDealDamage(sim, &aoeTarget.Unit, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
			}
		},
	})
}

func (druid *Druid) CanSwipeCat() bool {
	return druid.InForm(Cat) && druid.CurrentEnergy() >= druid.CurrentSwipeCatCost()
}

func (druid *Druid) CurrentSwipeCatCost() float64 {
	return druid.SwipeCat.ApplyCostModifiers(druid.SwipeCat.BaseCost)
}

func (druid *Druid) CanSwipeBear() bool {
	return druid.InForm(Bear) && druid.CurrentRage() >= druid.SwipeBear.DefaultCast.Cost
}

func (druid *Druid) IsSwipeSpell(spell *core.Spell) bool {
	return spell == druid.SwipeBear || spell == druid.SwipeCat
}
