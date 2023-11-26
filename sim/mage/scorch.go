package mage

import (
	"time"

	"github.com/wowsims/classic/sim/core"
)

func (mage *Mage) getScorchConfig(rank int) core.SpellConfig {
	spellCoeff := [8]float64{0, .429, .429, .429, .429, .429, .429, .429}[rank]
	baseDamage := [8][]float64{{0}, {56, 69}, {81, 98}, {105, 126}, {139, 165}, {168, 199}, {207, 247}, {237, 280}}[rank]
	spellId := [8]int32{0, 2948, 8444, 8445, 8446, 10205, 10206, 10207}[rank]
	manaCost := [8]float64{0, 50, 65, 80, 100, 115, 135, 150}[rank]
	level := [8]int{0, 22, 28, 34, 40, 46, 52, 58}[rank]
	castTime := [8]int{0, 1500, 1500, 1500, 1500, 1500, 1500, 1500}[rank]

	return core.SpellConfig{
		ActionID:      core.ActionID{SpellID: spellId},
		SpellSchool:   core.SpellSchoolFire,
		ProcMask:      core.ProcMaskSpellDamage,
		Flags:         core.SpellFlagAPL | SpellFlagMage,
		RequiredLevel: level,
		Rank:          rank,

		ManaCost: core.ManaCostOptions{
			FlatCost: manaCost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * time.Duration(castTime),
			},
		},

		BonusCritRating: 0 +
			2*float64(mage.Talents.Incinerate)*core.SpellCritRatingPerCritChance,
		DamageMultiplier: 1,
		CritMultiplier:   mage.DefaultSpellCritMultiplier(),
		ThreatMultiplier: 1 - (0.15 * float64(mage.Talents.BurningSoul)),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(baseDamage[0], baseDamage[1]) + spellCoeff*spell.SpellPower()
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
		},
	}
}

func (mage *Mage) registerScorchSpell() {
	maxRank := 7

	for i := 1; i <= maxRank; i++ {
		config := mage.getScorchConfig(i)

		if config.RequiredLevel <= int(mage.Level) {
			mage.Scorch = mage.GetOrRegisterSpell(config)
		}
	}
}
