package warlock

import (
	"strconv"
	"time"

	"github.com/wowsims/sod/sim/core"
)

func (warlock *Warlock) getCorruptionConfig(rank int) core.SpellConfig {
	dotTickCoeff := [8]float64{0, .08, .155, .167, .167, .167, .167, .167}[rank] // per tick
	baseDamage := [8]float64{0, 40, 90, 222, 324, 486, 666, 822}[rank]
	spellId := [8]int32{0, 172, 6222, 6223, 7648, 11671, 11672, 25311}[rank]
	manaCost := [8]float64{0, 35, 55, 100, 160, 225, 290, 340}[rank]
	ticks := [8]int32{0, 4, 5, 6, 6, 6, 6, 6}[rank]
	level := [8]int{0, 4, 14, 24, 34, 44, 54, 60}[rank]
	castTime := time.Millisecond * (2000 - (400 * time.Duration(warlock.Talents.ImprovedCorruption)))

	return core.SpellConfig{
		ActionID:      core.ActionID{SpellID: spellId},
		SpellSchool:   core.SpellSchoolShadow,
		ProcMask:      core.ProcMaskSpellDamage,
		Flags:         core.SpellFlagAPL | core.SpellFlagHauntSE | core.SpellFlagResetAttackSwing,
		Rank:          rank,
		RequiredLevel: level,

		ManaCost: core.ManaCostOptions{
			FlatCost: manaCost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				CastTime: castTime,
				GCD:      core.GCDDefault,
			},
		},

		BonusHitRating:   float64(warlock.Talents.Suppression) * 2 * core.CritRatingPerCritChance,
		BonusCritRating:  0,
		DamageMultiplier: 1,
		CritMultiplier:   1,
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Corruption-" + strconv.Itoa(rank),
			},

			NumberOfTicks: ticks,
			TickLength:    time.Second * 3,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.SnapshotBaseDamage = baseDamage/float64(ticks) + (dotTickCoeff * dot.Spell.SpellPower())
				dot.SnapshotAttackerMultiplier = 1 // dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex])
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHit)
			if result.Landed() {
				spell.Dot(target).Apply(sim)
			}
			spell.DealOutcome(sim, result)
		},
		ExpectedTickDamage: func(sim *core.Simulation, target *core.Unit, spell *core.Spell, useSnapshot bool) *core.SpellResult {
			if useSnapshot {
				dot := spell.Dot(target)
				return dot.CalcSnapshotDamage(sim, target, dot.Spell.OutcomeExpectedMagicAlwaysHit)
			} else {
				baseDamage := baseDamage/float64(ticks) + (dotTickCoeff * spell.SpellPower())
				return spell.CalcPeriodicDamage(sim, target, baseDamage, spell.OutcomeExpectedMagicAlwaysHit)
			}
		},
	}
}

func (warlock *Warlock) registerCorruptionSpell() {
	maxRank := 7

	for i := 1; i <= maxRank; i++ {
		config := warlock.getCorruptionConfig(i)

		if config.RequiredLevel <= int(warlock.Level) {
			warlock.Corruption = warlock.GetOrRegisterSpell(config)
		}
	}
}
