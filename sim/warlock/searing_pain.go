package warlock

import (
	"time"

	"github.com/wowsims/sod/sim/core"
)

func (warlock *Warlock) getSearingPainBaseConfig(rank int) core.SpellConfig {
	spellCoeff := [7]float64{0, .396, .429, .429, .429, .429, .429}[rank]
	baseDamage := [7][]float64{{0}, {38, 47}, {65, 77}, {93, 112}, {131, 155}, {168, 199}, {208, 244}}[rank]
	spellId := [7]int32{0, 5676, 17919, 17920, 17921, 17922, 17923}[rank]
	manaCost := [7]float64{0, 45, 68, 91, 118, 141, 168}[rank]
	level := [7]int{0, 18, 26, 36, 42, 50, 58}[rank]
	castTime := time.Millisecond * 1500

	return core.SpellConfig{
		ActionID:      core.ActionID{SpellID: spellId},
		SpellSchool:   core.SpellSchoolFire,
		ProcMask:      core.ProcMaskSpellDamage,
		Flags:         core.SpellFlagAPL | core.SpellFlagResetAttackSwing,
		RequiredLevel: level,
		Rank:          rank,

		ManaCost: core.ManaCostOptions{
			FlatCost:   manaCost,
			Multiplier: 1 - float64(warlock.Talents.Cataclysm)*0.01,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: castTime,
			},
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				if warlock.MetamorphosisAura != nil && warlock.MetamorphosisAura.IsActive() {
					spell.DefaultCast.CastTime = 0
				} else {
					spell.DefaultCast.CastTime = castTime
				}
			},
		},
		BonusCritRating: 0.0 +
			float64(warlock.Talents.Devastation)*core.CritRatingPerCritChance +
			2.0*float64(warlock.Talents.ImprovedSearingPain)*core.CritRatingPerCritChance,

		DamageMultiplier:         1 + 0.02*float64(warlock.Talents.Emberstorm),
		DamageMultiplierAdditive: 1,
		CritMultiplier:           warlock.SpellCritMultiplier(1, core.TernaryFloat64(warlock.Talents.Ruin, 1, 0)),
		ThreatMultiplier:         2,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			damage := sim.Roll(baseDamage[0], baseDamage[1]) + spellCoeff*spell.SpellPower()

			if warlock.LakeOfFireAuras != nil && warlock.LakeOfFireAuras.Get(target).IsActive() {
				damage *= 1.4
			}

			spell.CalcAndDealDamage(sim, target, damage, spell.OutcomeMagicHitAndCrit)
			warlock.EverlastingAfflictionRefresh(sim, target)
		},
	}
}

func (warlock *Warlock) registerSearingPainSpell() {
	maxRank := 6

	for i := 1; i <= maxRank; i++ {
		config := warlock.getSearingPainBaseConfig(i)

		if config.RequiredLevel <= int(warlock.Level) {
			warlock.SearingPain = warlock.GetOrRegisterSpell(config)
		}
	}
}
