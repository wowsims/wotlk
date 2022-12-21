package druid

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (druid *Druid) registerRakeSpell() {
	actionID := core.ActionID{SpellID: 48574}
	cost := 40.0 - float64(druid.Talents.Ferocity)

	bleedCategory := druid.CurrentTarget.GetExclusiveEffectCategory(core.BleedEffectCategory)

	numTicks := 3 + core.TernaryInt32(druid.HasSetBonus(ItemSetMalfurionsBattlegear, 2), 1, 0)
	dotCanCrit := druid.HasSetBonus(ItemSetLasherweaveBattlegear, 4)

	druid.Rake = druid.RegisterSpell(core.SpellConfig{
		ActionID:     actionID,
		SpellSchool:  core.SpellSchoolPhysical,
		ProcMask:     core.ProcMaskMeleeMHSpecial,
		Flags:        core.SpellFlagMeleeMetrics | core.SpellFlagIgnoreResists,
		ResourceType: stats.Energy,
		BaseCost:     cost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: cost,
				GCD:  time.Second,
			},
			IgnoreHaste: true,
		},

		DamageMultiplier: 1 + 0.1*float64(druid.Talents.SavageFury),
		CritMultiplier:   druid.MeleeCritMultiplier(Cat),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 176 + 0.01*spell.MeleeAttackPower()
			if bleedCategory.AnyActive() {
				baseDamage *= 1.3
			}

			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)

			if result.Landed() {
				druid.AddComboPoints(sim, 1, spell.ComboPointMetrics())
				druid.RakeDot.Apply(sim)
			} else {
				druid.AddEnergy(sim, spell.CurCast.Cost*0.8, druid.EnergyRefundMetrics)
			}
		},

		ExpectedDamage: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) *core.SpellResult {
			baseDamage := 176 + 0.01*spell.MeleeAttackPower()
			tickBase := (358 + 0.06*spell.MeleeAttackPower()) * float64(numTicks)
			if bleedCategory.AnyActive() {
				baseDamage *= 1.3
				tickBase *= 1.3
			}

			initial := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeExpectedMagicAlwaysHit)
			ticks := spell.CalcDamage(sim, target, tickBase, spell.OutcomeExpectedMagicAlwaysHit)

			critRating := druid.GetStat(stats.MeleeCrit) + spell.BonusCritRating
			critChance := critRating / (core.CritRatingPerCritChance * 100)
			critMod := (critChance * (spell.CritMultiplier - 1))

			if dotCanCrit {
				ticks.Damage *= critChance * (1 + critMod)
			}

			ticks.Damage += initial.Damage * (critChance * (1 + critMod))
			return ticks
		},
	})

	dotAura := druid.CurrentTarget.RegisterAura(druid.applyRendAndTear(core.Aura{
		Label:    "Rake-" + strconv.Itoa(int(druid.Index)),
		ActionID: actionID,
		Duration: time.Second * 9,
	}))
	druid.RakeDot = core.NewDot(core.Dot{
		Spell:         druid.Rake,
		Aura:          dotAura,
		NumberOfTicks: numTicks,
		TickLength:    time.Second * 3,
		OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
			dot.SnapshotBaseDamage = 358 + 0.06*dot.Spell.MeleeAttackPower()
			attackTable := dot.Spell.Unit.AttackTables[target.UnitIndex]
			dot.SnapshotCritChance = dot.Spell.PhysicalCritChance(target, attackTable)
			dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(attackTable)
		},
		OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
			if dotCanCrit {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)
			} else {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.Spell.OutcomeAlwaysHit)
			}
		},
	})
}

func (druid *Druid) CanRake() bool {
	return druid.InForm(Cat) && druid.CurrentEnergy() >= druid.CurrentRakeCost()
}

func (druid *Druid) CurrentRakeCost() float64 {
	return druid.Rake.ApplyCostModifiers(druid.Rake.BaseCost)
}
