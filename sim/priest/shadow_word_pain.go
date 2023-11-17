package priest

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

func (priest *Priest) getShadowWordPainConfig(rank int) core.SpellConfig {
	spellCoeff := [9]float64{0, 0.067, 0.104, 0.154, 0.167, 0.167, 0.167, 0.167, 0.167}[rank] // per tick
	baseDamage := [9]float64{0, 30, 66, 132, 234, 366, 510, 672, 852}[rank]
	spellId := [9]int32{0, 589, 594, 970, 992, 2767, 10892, 10893, 10894}[rank]
	manaCost := [9]float64{0, 25, 50, 95, 155, 230, 305, 385, 470}[rank]
	level := [9]int{0, 4, 10, 18, 26, 34, 42, 50, 58}[rank]

	return core.SpellConfig{
		ActionID:      core.ActionID{SpellID: spellId},
		SpellSchool:   core.SpellSchoolShadow,
		ProcMask:      core.ProcMaskSpellDamage,
		Flags:         core.SpellFlagAPL,
		Rank:          rank,
		RequiredLevel: level,

		ManaCost: core.ManaCostOptions{
			FlatCost: manaCost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		BonusHitRating:   float64(priest.Talents.ShadowFocus) * 2 * core.SpellHitRatingPerHitChance,
		BonusCritRating:  0,
		DamageMultiplier: 1,
		CritMultiplier:   1,
		ThreatMultiplier: 1 - 0.08*float64(priest.Talents.ShadowAffinity),

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "ShadowWordPain",
			},

			NumberOfTicks: 6 + (priest.Talents.ImprovedShadowWordPain),
			TickLength:    time.Second * 3,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.SnapshotBaseDamage = baseDamage/6 + (spellCoeff * dot.Spell.SpellPower())
				dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex])
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHit)
			if result.Landed() {
				spell.SpellMetrics[target.UnitIndex].Hits--
				priest.AddShadowWeavingStack(sim)
				spell.Dot(target).Apply(sim)
			}
			spell.DealOutcome(sim, result)
		},
		ExpectedTickDamage: func(sim *core.Simulation, target *core.Unit, spell *core.Spell, useSnapshot bool) *core.SpellResult {
			if useSnapshot {
				dot := spell.Dot(target)
				return dot.CalcSnapshotDamage(sim, target, dot.Spell.OutcomeExpectedMagicAlwaysHit)
			} else {
				baseDamage := baseDamage/6 + (spellCoeff * spell.SpellPower())
				return spell.CalcPeriodicDamage(sim, target, baseDamage, spell.OutcomeExpectedMagicAlwaysHit)
			}
		},
	}
}

func (priest *Priest) registerShadowWordPainSpell() {
	maxRank := 8
	priest.ShadowWordPain = priest.GetOrRegisterSpell(priest.getShadowWordPainConfig(maxRank))

	for i := maxRank - 1; i > 0; i-- {
		priest.GetOrRegisterSpell(priest.getShadowWordPainConfig(i))
	}
}
