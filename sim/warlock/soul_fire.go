package warlock

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (warlock *Warlock) registerSoulFireSpell() {
	baseCost := 0.09 * warlock.BaseMana
	decimationMod := 1.0 - 0.2*float64(warlock.Talents.Decimation)

	warlock.SoulFire = warlock.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 47825},
		SpellSchool:  core.SpellSchoolFire,
		ProcMask:     core.ProcMaskSpellDamage,
		MissileSpeed: 24,
		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost:     baseCost * (1 - []float64{0, .04, .07, .10}[warlock.Talents.Cataclysm]),
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * time.Duration(6000-400*warlock.Talents.Bane),
			},
			ModifyCast: func(_ *core.Simulation, _ *core.Spell, cast *core.Cast) {
				totalMod := warlock.backdraftModifier()
				if warlock.DecimationAura.IsActive() {
					totalMod *= decimationMod
				}

				cast.GCD = time.Duration(float64(cast.GCD) * totalMod)
				cast.CastTime = time.Duration(float64(cast.CastTime) * totalMod)
			},
		},

		BonusCritRating: 0 +
			warlock.masterDemonologistFireCrit +
			core.TernaryFloat64(warlock.Talents.Devastation, 5*core.CritRatingPerCritChance, 0) +
			core.TernaryFloat64(warlock.HasSetBonus(ItemSetDarkCovensRegalia, 2), 5*core.CritRatingPerCritChance, 0),
		DamageMultiplierAdditive: 1 +
			warlock.GrandFirestoneBonus() +
			0.03*float64(warlock.Talents.Emberstorm),
		CritMultiplier:   warlock.SpellCritMultiplier(1, float64(warlock.Talents.Ruin)/5),
		ThreatMultiplier: 1 - 0.1*float64(warlock.Talents.DestructiveReach),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(1323, 1657) + 1.15*spell.SpellPower()
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
			})

			if warlock.MoltenCoreAura.IsActive() {
				warlock.MoltenCoreAura.RemoveStack(sim)
			}
		},
	})
}
