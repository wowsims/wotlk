package druid

import (
	"time"

	"github.com/wowsims/classic/sod/sim/core"
)

func (druid *Druid) getStarfireBaseConfig(rank int) core.SpellConfig {
	baseDamage := [8][]float64{{0}, {95, 115}, {146, 177}, {212, 253}, {293, 348}, {378, 445}, {451, 531}, {496, 584}}[rank]
	spellId := [8]int32{0, 2912, 8949, 8950, 8951, 9875, 9876, 25298}[rank]
	manaCost := [8]float64{0, 95, 135, 180, 230, 275, 315, 340}[rank]
	level := [8]int{0, 20, 26, 34, 42, 50, 58, 60}[rank]
	castTime := 3500

	return core.SpellConfig{
		ActionID:      core.ActionID{SpellID: spellId},
		SpellSchool:   core.SpellSchoolNature,
		ProcMask:      core.ProcMaskSpellDamage,
		Flags:         core.SpellFlagAPL,
		RequiredLevel: level,
		Rank:          rank,

		ManaCost: core.ManaCostOptions{
			FlatCost: manaCost * (1 - 0.03*float64(druid.Talents.Moonglow)),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond*time.Duration(castTime) - time.Millisecond*100*time.Duration(druid.Talents.ImprovedStarfire),
			},
			CastTime: druid.NaturesGraceCastTime(),
		},

		BonusCritRating:  0,
		DamageMultiplier: 1 + 0.02*float64(druid.Talents.Moonfury),
		CritMultiplier:   druid.BalanceCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(baseDamage[0], baseDamage[1]) + spell.SpellPower()
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
		},
	}
}

func (druid *Druid) registerStarfireSpell() {
	maxRank := 7

	for i := 1; i <= maxRank; i++ {
		config := druid.getStarfireBaseConfig(i)

		if config.RequiredLevel <= int(druid.Level) {
			druid.Starfire = druid.RegisterSpell(Humanoid|Moonkin, config)
		}
	}
}
