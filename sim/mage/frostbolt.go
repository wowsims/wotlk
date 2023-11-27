package mage

import (
	"time"

	"github.com/wowsims/classic/sod/sim/core"
)

func (mage *Mage) getFrostboltBaseConfig(rank int) core.SpellConfig {
	spellCoeff := [12]float64{0, .163, .269, .463, .706, .814, .814, .814, .814, .814, .814, .814}[rank]
	baseDamage := [12][]float64{{0, 0}, {20, 22}, {33, 38}, {54, 61}, {78, 87}, {132, 144}, {180, 197}, {235, 255}, {301, 326}, {363, 394}, {440, 475}, {515, 555}}[rank]
	spellId := [12]int32{0, 116, 205, 837, 7322, 8406, 8407, 8408, 10179, 10180, 10181, 25304}[rank]
	manaCost := [12]float64{0, 25, 35, 50, 65, 100, 130, 160, 195, 225, 260, 290}[rank]
	level := [12]int{0, 4, 8, 14, 20, 26, 32, 38, 44, 50, 56, 60}[rank]
	castTime := [12]int{0, 1500, 1800, 2200, 2600, 3000, 3000, 3000, 3000, 3000, 3000, 3000}[rank]

	return core.SpellConfig{
		ActionID:      core.ActionID{SpellID: spellId},
		SpellSchool:   core.SpellSchoolFrost,
		ProcMask:      core.ProcMaskSpellDamage,
		Flags:         core.SpellFlagAPL | SpellFlagMage | SpellFlagChillSpell,
		MissileSpeed:  28,
		RequiredLevel: level,
		Rank:          rank,

		ManaCost: core.ManaCostOptions{
			FlatCost: manaCost * (1 - 0.05*float64(mage.Talents.FrostChanneling)),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond*time.Duration(castTime) - time.Millisecond*100*time.Duration(mage.Talents.ImprovedFrostbolt),
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   mage.SpellCritMultiplier(1, 0.2*float64(mage.Talents.IceShards)),
		ThreatMultiplier: 1 - (0.1 * float64(mage.Talents.FrostChanneling)),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(baseDamage[0], baseDamage[1]) + spellCoeff*spell.SpellPower()
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)

			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				if result.Landed() {
					spell.DealDamage(sim, result)
				}
			})
		},
	}
}

func (mage *Mage) registerFrostboltSpell() {
	maxRank := 11

	for i := 1; i <= maxRank; i++ {
		config := mage.getFrostboltBaseConfig(i)

		if config.RequiredLevel <= int(mage.Level) {
			mage.Frostbolt = mage.GetOrRegisterSpell(config)
		}
	}
}
