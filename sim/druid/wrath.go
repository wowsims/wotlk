package druid

import (
	"time"

	"github.com/wowsims/classic/sim/core"
)

func (druid *Druid) getWrathBaseConfig(rank int) core.SpellConfig {
	spellCoeff := [9]float64{0, 0.123, 0.231, 0.443, 0.571, 0.571, 0.571, 0.571, 0.571}[rank]
	baseDamage := [9][]float64{{0}, {13, 16}, {28, 33}, {48, 57}, {69, 79}, {108, 123}, {148, 167}, {198, 221}, {248, 277}}[rank]
	spellId := [9]int32{0, 5176, 5177, 5178, 5179, 5180, 6780, 8905, 9912}[rank]
	manaCost := [9]float64{0, 20, 35, 55, 70, 100, 125, 155, 180}[rank]
	level := [9]int{0, 1, 6, 14, 22, 30, 38, 46, 54}[rank]
	castTime := [9]int{0, 1500, 1700, 2000, 2000, 2000, 2000, 2000, 2000}[rank]

	return core.SpellConfig{
		ActionID:      core.ActionID{SpellID: spellId},
		SpellSchool:   core.SpellSchoolNature,
		ProcMask:      core.ProcMaskSpellDamage,
		Flags:         core.SpellFlagAPL,
		RequiredLevel: level,
		Rank:          rank,
		MissileSpeed:  20,

		ManaCost: core.ManaCostOptions{
			FlatCost: core.TernaryFloat64(druid.FuryOfStormrageAura != nil, 0, manaCost),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond*time.Duration(castTime) - time.Millisecond*100*time.Duration(druid.Talents.ImprovedWrath),
			},
			CastTime: druid.NaturesGraceCastTime(),
		},

		BonusCritRating:  0,
		DamageMultiplier: 1 + 0.02*float64(druid.Talents.Moonfury),
		CritMultiplier:   druid.BalanceCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(baseDamage[0], baseDamage[1]) + spellCoeff*spell.SpellPower()
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)

			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
			})
		},
	}
}

func (druid *Druid) registerWrathSpell() {
	maxRank := 8

	for i := 1; i <= maxRank; i++ {
		config := druid.getWrathBaseConfig(i)

		if config.RequiredLevel <= int(druid.Level) {
			druid.Wrath = druid.RegisterSpell(Humanoid|Moonkin, config)
		}
	}
}
