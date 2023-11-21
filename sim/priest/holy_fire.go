package priest

import (
	"time"

	"github.com/wowsims/classic/sim/core"
)

func (priest *Priest) getHolyFireConfig(rank int) core.SpellConfig {
	directCoeff := 0.75
	dotCoeff := 0.05
	baseDamage := [9][]float64{{0}, {84, 104}, {106, 131}, {144, 178}, {178, 223}, {219, 273}, {271, 340}, {323, 406}, {355, 449}}[rank]
	dotDamage := [9]float64{0, 30, 40, 55, 65, 85, 100, 125, 145}[rank]
	spellId := [9]int32{0, 14914, 15262, 15263, 15264, 15265, 15266, 15267, 15261}[rank]
	manaCost := [9]float64{0, 85, 95, 125, 145, 170, 200, 230, 255}[rank]
	level := [9]int{0, 20, 24, 30, 36, 42, 48, 54, 60}[rank]

	return core.SpellConfig{
		ActionID:      core.ActionID{SpellID: spellId},
		SpellSchool:   core.SpellSchoolHoly,
		ProcMask:      core.ProcMaskSpellDamage,
		Flags:         core.SpellFlagAPL,
		RequiredLevel: level,
		Rank:          rank,

		ManaCost: core.ManaCostOptions{
			FlatCost: manaCost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond*3500 - time.Millisecond*100*time.Duration(priest.Talents.DivineFury),
			},
		},

		BonusCritRating:  0,
		DamageMultiplier: 1 + 0.05*float64(priest.Talents.SearingLight),
		CritMultiplier:   priest.SpellCritMultiplier(priest.SpellCritMultiplier(float64(priest.Talents.HolySpecialization)*1*core.CritRatingPerCritChance, 1), 1),
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "HolyFire",
			},
			NumberOfTicks: 5,
			TickLength:    time.Second * 2,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, _ bool) {
				dot.SnapshotBaseDamage = dotDamage + dotCoeff*dot.Spell.SpellPower()
				dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex])
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(baseDamage[0], baseDamage[1]) + directCoeff*spell.SpellPower()
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			if result.Landed() {
				spell.Dot(target).Apply(sim)
			}
			spell.DealDamage(sim, result)
		},
	}
}

func (priest *Priest) registerHolyFire() {
	maxRank := 8

	for i := 1; i <= maxRank; i++ {
		config := priest.getHolyFireConfig(i)

		if config.RequiredLevel <= int(priest.Level) {
			priest.HolyFire = priest.GetOrRegisterSpell(config)
		}
	}
}
