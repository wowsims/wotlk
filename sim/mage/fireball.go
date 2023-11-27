package mage

import (
	"time"

	"github.com/wowsims/classic/sod/sim/core"
)

func (mage *Mage) getFireballBaseConfig(rank int) core.SpellConfig {
	spellCoeff := [13]float64{0, .123, .271, .5, .793, 1, 1, 1, 1, 1, 1, 1, 1}[rank]
	baseDamage := [13][]float64{{0}, {16, 25}, {34, 49}, {57, 77}, {89, 122}, {146, 195}, {207, 274}, {264, 345}, {328, 425}, {404, 518}, {488, 623}, {561, 715}, {596, 760}}[rank]
	baseDotDamage := [13]float64{0, 2, 3, 6, 12, 20, 28, 32, 40, 52, 60, 72, 76}[rank]
	spellId := [13]int32{0, 133, 143, 145, 3140, 8400, 8401, 8402, 10148, 10149, 10150, 10151, 25306}[rank]
	manaCost := [13]float64{0, 30, 45, 65, 95, 140, 185, 220, 260, 305, 350, 395, 410}[rank]
	level := [13]int{0, 1, 6, 12, 18, 24, 30, 36, 42, 48, 54, 60, 60}[rank]
	castTime := [13]int{0, 1500, 2000, 2500, 3000, 3500, 3500, 3500, 3500, 3500, 3500, 3500, 3500}[rank]

	return core.SpellConfig{
		ActionID:      core.ActionID{SpellID: spellId},
		SpellSchool:   core.SpellSchoolFire,
		ProcMask:      core.ProcMaskSpellDamage,
		Flags:         core.SpellFlagAPL,
		RequiredLevel: level,
		Rank:          rank,
		MissileSpeed:  24,

		ManaCost: core.ManaCostOptions{
			FlatCost: manaCost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond*time.Duration(castTime) - time.Millisecond*100*time.Duration(mage.Talents.ImprovedFireball),
			},
		},
		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Fireball",
			},
			NumberOfTicks: 4,
			TickLength:    time.Second * 2,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, _ bool) {
				dot.SnapshotBaseDamage = baseDotDamage / 4.0
				dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex])
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   mage.DefaultSpellCritMultiplier(),
		ThreatMultiplier: 1 - (0.15 * float64(mage.Talents.BurningSoul)),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(baseDamage[0], baseDamage[1]) + spellCoeff*spell.SpellPower()
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				if result.Landed() {
					spell.DealDamage(sim, result)
					spell.Dot(target).Apply(sim)
				}
			})
		},
	}
}

func (mage *Mage) registerFireballSpell() {
	maxRank := 12

	for i := 1; i <= maxRank; i++ {
		config := mage.getFireballBaseConfig(i)

		if config.RequiredLevel <= int(mage.Level) {
			mage.Fireball = mage.GetOrRegisterSpell(config)
		}
	}
}
