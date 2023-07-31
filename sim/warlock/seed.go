package warlock

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

func (warlock *Warlock) registerSeedSpell() {
	actionID := core.ActionID{SpellID: 47836}

	seedExplosion := warlock.RegisterSpell(core.SpellConfig{
		ActionID:    actionID.WithTag(1), // actually 47834
		SpellSchool: core.SpellSchoolShadow,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       core.SpellFlagHauntSE | core.SpellFlagNoLogs,

		BonusCritRating: 0 +
			float64(warlock.Talents.ImprovedCorruption)*core.CritRatingPerCritChance,
		DamageMultiplierAdditive: 1 +
			warlock.GrandFirestoneBonus() +
			0.03*float64(warlock.Talents.ShadowMastery) +
			0.01*float64(warlock.Talents.Contagion),
		CritMultiplier:   warlock.DefaultSpellCritMultiplier(),
		ThreatMultiplier: 1 - 0.1*float64(warlock.Talents.ImprovedDrainSoul),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDmg := (sim.Roll(1633, 1897) + 0.286*spell.SpellPower()) * sim.Encounter.AOECapMultiplier()
			for _, aoeTarget := range sim.Encounter.TargetUnits {
				spell.CalcAndDealDamage(sim, aoeTarget, baseDmg, spell.OutcomeMagicHitAndCrit)
			}
		},
	})

	warlock.SeedDamageTracker = make([]float64, len(warlock.Env.AllUnits))
	trySeedPop := func(sim *core.Simulation, target *core.Unit, dmg float64) {
		warlock.SeedDamageTracker[target.UnitIndex] += dmg
		if warlock.SeedDamageTracker[target.UnitIndex] > 1518 {
			warlock.Seed.Dot(target).Deactivate(sim)
			seedExplosion.Cast(sim, target)
		}
	}

	warlock.Seed = warlock.RegisterSpell(core.SpellConfig{
		ActionID:     actionID,
		SpellSchool:  core.SpellSchoolShadow,
		ProcMask:     core.ProcMaskEmpty,
		Flags:        core.SpellFlagHauntSE | core.SpellFlagAPL,
		MissileSpeed: 28,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.34,
			Multiplier: 1 - 0.02*float64(warlock.Talents.Suppression),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * 2000,
			},
		},

		DamageMultiplierAdditive: 1 +
			warlock.GrandSpellstoneBonus() +
			0.03*float64(warlock.Talents.ShadowMastery) +
			0.01*float64(warlock.Talents.Contagion) +
			core.TernaryFloat64(warlock.Talents.SiphonLife, 0.05, 0),
		ThreatMultiplier: 1 - 0.1*float64(warlock.Talents.ImprovedDrainSoul),

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Seed",
				OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if !result.Landed() {
						return
					}
					if spell.ActionID.SpellID == actionID.SpellID {
						return // Seed can't pop seed.
					}
					trySeedPop(sim, aura.Unit, result.Damage)
				},
				OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					trySeedPop(sim, aura.Unit, result.Damage)
				},
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					warlock.SeedDamageTracker[aura.Unit.UnitIndex] = 0
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					warlock.SeedDamageTracker[aura.Unit.UnitIndex] = 0
				},
			},

			NumberOfTicks: 6,
			TickLength:    time.Second * 3,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.SnapshotBaseDamage = 1518/6 + 0.25*dot.Spell.SpellPower()
				dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex])
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHit)
			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				if result.Landed() {
					// seed is mutually exclusive with corruption
					warlock.Corruption.Dot(target).Deactivate(sim)

					if warlock.Rotation.DetonateSeed {
						seedExplosion.Cast(sim, target)
					} else {
						spell.Dot(target).Apply(sim)
					}
				}
			})
		},
	})
}
