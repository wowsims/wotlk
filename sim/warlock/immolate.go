package warlock

import (
	"time"

	"github.com/wowsims/sod/sim/core"
)

func (warlock *Warlock) getImmolateConfig(rank int) core.SpellConfig {
	directCoeff := [9]float64{0, .058, .125, .2, .2, .2, .2, .2, .2}[rank]
	dotCoeff := [9]float64{0, .037, .081, .13, .13, .13, .13, .13, .13}[rank]
	baseDamage := [9]float64{0, 11, 24, 53, 101, 148, 208, 258, 279}[rank]
	dotDamage := [9]float64{0, 20, 40, 90, 165, 255, 365, 485, 510}[rank]
	spellId := [9]int32{0, 348, 707, 1094, 2941, 11665, 11667, 11668, 25309}[rank]
	manaCost := [9]float64{0, 25, 45, 90, 155, 220, 295, 370, 380}[rank]
	level := [9]int{0, 1, 10, 20, 30, 40, 50, 60, 60}[rank]

	return core.SpellConfig{
		ActionID:      core.ActionID{SpellID: spellId},
		SpellSchool:   core.SpellSchoolFire,
		ProcMask:      core.ProcMaskSpellDamage,
		Flags:         core.SpellFlagAPL | core.SpellFlagResetAttackSwing,
		Rank:          rank,
		RequiredLevel: level,

		ManaCost: core.ManaCostOptions{
			FlatCost:   manaCost,
			Multiplier: 1 - float64(warlock.Talents.Cataclysm)*0.01,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * (2000 - 100*time.Duration(warlock.Talents.Bane)),
			},
		},

		BonusHitRating:   0,
		BonusCritRating:  float64(warlock.Talents.Devastation) * core.SpellCritRatingPerCritChance,
		DamageMultiplier: 1 + 0.05*float64(warlock.Talents.ImprovedImmolate) + 0.02*float64(warlock.Talents.Emberstorm),
		CritMultiplier:   warlock.SpellCritMultiplier(1, core.TernaryFloat64(warlock.Talents.Ruin, 1, 0)),
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Immolate",
			},

			NumberOfTicks: 5,
			TickLength:    time.Second * 3,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.SnapshotBaseDamage = dotDamage/5 + dotCoeff*dot.Spell.SpellPower()

				if warlock.LakeOfFireAuras != nil && warlock.LakeOfFireAuras.Get(target).IsActive() {
					dot.SnapshotBaseDamage *= 1.4
				}

				dot.SnapshotAttackerMultiplier = 1 // dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex])
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := baseDamage + directCoeff*spell.SpellPower()
			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)

			if result.Landed() {
				spell.Dot(target).Apply(sim)
			}
		},
		ExpectedTickDamage: func(sim *core.Simulation, target *core.Unit, spell *core.Spell, useSnapshot bool) *core.SpellResult {
			if useSnapshot {
				dot := spell.Dot(target)
				return dot.CalcSnapshotDamage(sim, target, dot.Spell.OutcomeExpectedMagicAlwaysHit)
			} else {
				baseDamage := dotDamage + dotCoeff*spell.SpellPower()
				return spell.CalcPeriodicDamage(sim, target, baseDamage, spell.OutcomeExpectedMagicAlwaysHit)
			}
		},
	}
}

func (warlock *Warlock) registerImmolateSpell() {
	maxRank := 8

	for i := 1; i <= maxRank; i++ {
		config := warlock.getImmolateConfig(i)

		if config.RequiredLevel <= int(warlock.Level) {
			warlock.Immolate = warlock.GetOrRegisterSpell(config)
		}
	}
}

// func (warlock *Warlock) registerImmolateSpell() {

// 	warlock.Immolate = warlock.RegisterSpell(core.SpellConfig{
// 		ActionID:    core.ActionID{SpellID: 47811},
// 		SpellSchool: core.SpellSchoolFire,
// 		ProcMask:    core.ProcMaskSpellDamage,
// 		Flags:       core.SpellFlagAPL,

// 		ManaCost: core.ManaCostOptions{
// 			BaseCost:   0.17,
// 			Multiplier: 1 - []float64{0, .04, .07, .10}[warlock.Talents.Cataclysm],
// 		},
// 		Cast: core.CastConfig{
// 			DefaultCast: core.Cast{
// 				GCD:      core.GCDDefault,
// 				CastTime: time.Millisecond * (2000 - 100*time.Duration(warlock.Talents.Bane)),
// 			},
// 		},

// 		BonusCritRating: 0 +
// 			core.TernaryFloat64(warlock.Talents.Devastation, 5*core.CritRatingPerCritChance, 0),
// 		DamageMultiplierAdditive: 1 +
// 			warlock.GrandFirestoneBonus() +
// 			0.03*float64(warlock.Talents.Emberstorm) +
// 			0.1*float64(warlock.Talents.ImprovedImmolate) +
// 			core.TernaryFloat64(warlock.HasSetBonus(ItemSetDeathbringerGarb, 2), 0.1, 0) +
// 			core.TernaryFloat64(warlock.HasSetBonus(ItemSetGuldansRegalia, 4), 0.1, 0),
// 		CritMultiplier:   warlock.SpellCritMultiplier(1, float64(warlock.Talents.Ruin)/5),
// 		ThreatMultiplier: 1 - 0.1*float64(warlock.Talents.DestructiveReach),

// 		Dot: core.DotConfig{
// 			Aura: core.Aura{
// 				Label: "Immolate",
// 				OnGain: func(aura *core.Aura, sim *core.Simulation) {
// 					if warlock.Talents.ChaosBolt {
// 						warlock.ChaosBolt.DamageMultiplierAdditive += fireAndBrimstoneBonus
// 					}
// 					warlock.Incinerate.DamageMultiplierAdditive += fireAndBrimstoneBonus
// 				},
// 				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
// 					if warlock.Talents.ChaosBolt {
// 						warlock.ChaosBolt.DamageMultiplierAdditive -= fireAndBrimstoneBonus
// 					}
// 					warlock.Incinerate.DamageMultiplierAdditive -= fireAndBrimstoneBonus
// 				},
// 			},
// 			NumberOfTicks: 5 + warlock.Talents.MoltenCore,
// 			TickLength:    time.Second * 3,

// 			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
// 				dot.SnapshotBaseDamage = 157 + 0.2*dot.Spell.SpellPower()
// 				attackTable := dot.Spell.Unit.AttackTables[target.UnitIndex]
// 				dot.SnapshotCritChance = dot.Spell.SpellCritChance(target)

// 				dot.Spell.DamageMultiplierAdditive += bonusPeriodicDamageMultiplier
// 				dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(attackTable)
// 				dot.Spell.DamageMultiplierAdditive -= bonusPeriodicDamageMultiplier
// 			},
// 			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
// 				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)
// 			},
// 		},

// 		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
// 			baseDamage := 460 + 0.2*spell.SpellPower()
// 			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
// 			if result.Landed() {
// 				spell.Dot(target).Apply(sim)
// 			}
// 			spell.DealDamage(sim, result)
// 		},
// 	})
// }
