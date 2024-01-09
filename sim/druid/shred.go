package druid

import (
	"time"

	"github.com/wowsims/sod/sim/core"
)

func (druid *Druid) registerShredSpell() {

	flatDamageBonus := map[int32]float64{
		25: 54.0,
		40: 99.0,
		50: 144.0,
		60: 180.0,
	}[druid.Level]

	druid.Shred = druid.RegisterSpell(Cat, core.SpellConfig{
		ActionID:    core.ActionID{SpellID: map[int32]int32{
			25: 5221,
			40: 8992,
			50: 9829,
			60: 9830,
		}[druid.Level]},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | core.SpellFlagAPL,

		EnergyCost: core.EnergyCostOptions{
			Cost:   60 /*- 9*float64(druid.Talents.ShreddingAttacks)*/,
			Refund: 0.8,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return !druid.PseudoStats.InFrontOfTarget
		},

		DamageMultiplier: 2.25,
		CritMultiplier:   druid.MeleeCritMultiplier(1,0),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := flatDamageBonus +
				spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower()) +
				spell.BonusWeaponDamage()


			modifier := 1.0
			if druid.BleedCategories.Get(target).AnyActive() {
				modifier += .3
			}

			/*
			ripDot := druid.Rip.Dot(target)
			if druid.AssumeBleedActive || ripDot.IsActive() || druid.Rake.Dot(target).IsActive() || druid.Lacerate.Dot(target).IsActive() {
				modifier *= 1.0 + (0.04 * float64(druid.Talents.RendAndTear))
			}
			*/

			baseDamage *= modifier

			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)

			if result.Landed() {
				druid.AddComboPoints(sim, 1, spell.ComboPointMetrics())
			} else {
				spell.IssueRefund(sim)
			}
		},
		ExpectedInitialDamage: func(sim *core.Simulation, target *core.Unit, spell *core.Spell, _ bool) *core.SpellResult {
			baseDamage := flatDamageBonus + spell.Unit.AutoAttacks.MH().CalculateAverageWeaponDamage(spell.MeleeAttackPower()) + spell.BonusWeaponDamage()

			modifier := 1.0
			if druid.BleedCategories.Get(target).AnyActive() {
				modifier += .3
			}

			/*
			if druid.AssumeBleedActive || druid.Rip.Dot(target).IsActive() || druid.Rake.Dot(target).IsActive() || druid.Lacerate.Dot(target).IsActive() {
				modifier *= 1.0 + (0.04 * float64(druid.Talents.RendAndTear))
			}
			*/
			baseDamage *= modifier

			baseres := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeExpectedMagicAlwaysHit)

			attackTable := spell.Unit.AttackTables[target.UnitIndex]
			critChance := spell.PhysicalCritChance(attackTable)
			critMod := (critChance * (spell.CritMultiplier - 1))

			baseres.Damage *= (1 + critMod)

			return baseres
		},
	})
}

func (druid *Druid) CanShred() bool {
	return !druid.PseudoStats.InFrontOfTarget && druid.CurrentEnergy() >= druid.CurrentShredCost()
}

func (druid *Druid) CurrentShredCost() float64 {
	return druid.Shred.ApplyCostModifiers(druid.Shred.DefaultCast.Cost)
}
