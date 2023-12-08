package druid

import (
	"time"

	"github.com/wowsims/sod/sim/core"
)

func (druid *Druid) getMoonfireBaseConfig(rank int) core.SpellConfig {
	spellCoeff := [11]float64{0, .06, .094, .128, .15, .15, .15, .15, .15, .15, .15}[rank]
	spellDotCoeff := [11]float64{0, .052, .081, .111, .13, .13, .13, .13, .13, .13, .13}[rank]
	baseDamage := [11][]float64{{0}, {9, 12}, {17, 21}, {30, 37}, {47, 55}, {70, 82}, {91, 108}, {117, 137}, {143, 168}, {172, 200}, {195, 228}}[rank]
	baseDotDamage := [11]float64{0, 12, 32, 52, 80, 124, 164, 212, 264, 320, 384}[rank]
	spellId := [11]int32{0, 8921, 8924, 8925, 8926, 8927, 8928, 8929, 9833, 9834, 9835}[rank]
	manaCost := [11]float64{0, 25, 50, 75, 105, 150, 190, 235, 280, 325, 375}[rank]
	level := [11]int{0, 4, 10, 16, 22, 28, 34, 40, 46, 52, 58}[rank]

	impMf := float64(druid.Talents.ImprovedMoonfire)
	moonfury := float64(druid.Talents.Moonfury)

	return core.SpellConfig{
		ActionID:      core.ActionID{SpellID: spellId},
		SpellSchool:   core.SpellSchoolArcane,
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
				CastTime: 0,
			},
		},
		Dot: core.DotConfig{
			Aura: core.Aura{
				Label:    "Moonfire",
				ActionID: core.ActionID{SpellID: spellId},
			},
			NumberOfTicks: core.TernaryInt32(rank < 2, 3, 4),
			TickLength:    time.Second * 3,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, _ bool) {
				dot.SnapshotBaseDamage = (baseDotDamage / 3.0) + spellDotCoeff*dot.Spell.SpellPower()
				dot.SnapshotAttackerMultiplier = 1 // dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex])
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
			},
		},

		BonusCritRating:  2 * impMf * core.SpellCritRatingPerCritChance,
		DamageMultiplier: 1 + 0.02*impMf + 0.02*moonfury,
		CritMultiplier:   druid.BalanceCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(baseDamage[0], baseDamage[1]) + spellCoeff*spell.SpellPower()
			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)

			if result.Landed() {
				spell.Dot(target).Apply(sim)
			}
		},
	}
}

func (druid *Druid) registerMoonfireSpell() {
	maxRank := 10

	for i := 1; i <= maxRank; i++ {
		config := druid.getMoonfireBaseConfig(i)

		if config.RequiredLevel <= int(druid.Level) {
			druid.Moonfire = druid.RegisterSpell(Humanoid|Moonkin, config)
		}
	}
}
