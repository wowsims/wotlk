package priest

import (
	"strconv"
	"time"

	"github.com/wowsims/sod/sim/core"
)

func (priest *Priest) getMindFlayTickSpell(rank int, numTicks int32, baseDamage float64) *core.Spell {
	spellCoeff := 0.15 // classic penalty for mf having a slow effect
	spellId := [7]int32{0, 16568, 7378, 17316, 17317, 17318, 18808}[rank]

	return priest.GetOrRegisterSpell(core.SpellConfig{
		ActionID:         core.ActionID{SpellID: spellId},
		Rank:             rank,
		SpellSchool:      core.SpellSchoolShadow,
		ProcMask:         core.ProcMaskProc | core.ProcMaskNotInSpellbook,
		BonusHitRating:   1, // Not an independent hit once initial lands
		BonusCritRating:  0,
		DamageMultiplier: 1,
		CritMultiplier:   1.0,
		ThreatMultiplier: 1 - 0.08*float64(priest.Talents.ShadowAffinity),
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			damage := baseDamage/3 + (spellCoeff * spell.SpellPower())
			damage *= priest.MindFlayModifier
			result := spell.CalcAndDealDamage(sim, target, damage, spell.OutcomeExpectedMagicAlwaysHit)

			if result.Landed() {
				priest.AddShadowWeavingStack(sim)
			}
		},
	})
}

func (priest *Priest) getMindFlaySpellConfig(rank int) core.SpellConfig {
	numTicks := int32(3)
	spellCoeff := 0.15 // classic penalty for mf having a slow effect
	baseDamage := [7]float64{0, 75, 126, 186, 261, 330, 426}[rank]
	spellId := [7]int32{0, 15407, 17311, 17312, 17313, 17314, 18807}[rank]
	manaCost := [7]float64{0, 45, 70, 100, 135, 165, 205}[rank]
	level := [7]int{0, 20, 28, 36, 44, 52, 60}[rank]

	tickLength := time.Second
	mindFlayTickSpell := priest.getMindFlayTickSpell(rank, numTicks, baseDamage)

	return core.SpellConfig{
		ActionID:      core.ActionID{SpellID: spellId},
		SpellSchool:   core.SpellSchoolShadow,
		ProcMask:      core.ProcMaskSpellDamage,
		Flags:         core.SpellFlagAPL | core.SpellFlagNoMetrics | core.SpellFlagChanneled,
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
		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "MindFlay-" + strconv.Itoa(int(rank)) + "-" + strconv.Itoa(int(numTicks)),
			},
			NumberOfTicks:       numTicks,
			TickLength:          tickLength,
			AffectedByCastSpeed: true,
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				mindFlayTickSpell.Cast(sim, target)
				mindFlayTickSpell.SpellMetrics[target.UnitIndex].Casts -= 1
			},
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHit)
			mindFlayTickSpell.SpellMetrics[target.UnitIndex].Casts += 1

			if result.Landed() {
				spell.Dot(target).Apply(sim)
			}
			spell.DealOutcome(sim, result)
		},
		ExpectedTickDamage: func(sim *core.Simulation, target *core.Unit, spell *core.Spell, _ bool) *core.SpellResult {
			baseDamage := (baseDamage + (spellCoeff * spell.SpellPower())) / 3

			return spell.CalcPeriodicDamage(sim, target, baseDamage, spell.OutcomeExpectedMagicAlwaysHit)
		},
	}
}

func (priest *Priest) registerMindFlay() {
	if !priest.Talents.MindFlay {
		return
	}
	maxRank := 6

	for i := 1; i < maxRank; i++ {
		config := priest.getMindFlaySpellConfig(i)

		if config.RequiredLevel <= int(priest.Level) {
			priest.MindFlay = priest.GetOrRegisterSpell(config)
		}
	}
}
