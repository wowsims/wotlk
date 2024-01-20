package hunter

import (
	"strconv"
	"time"

	"github.com/wowsims/sod/sim/core"
)

func (hunter *Hunter) getSerpentStingConfig(rank int) core.SpellConfig {
	spellId := [10]int32{0, 1978, 13549, 13550, 13551, 13552, 13553, 13554, 13555, 25295}[rank]
	baseDamage := [10]float64{0, 20, 40, 80, 140, 210, 290, 385, 490, 555}[rank] / 5
	spellCoeff := [10]float64{0, .08, .125, .185, .2, .2, .2, .2, .2, .2}[rank]
	manaCost := [10]float64{0, 15, 30, 50, 80, 115, 150, 190, 230, 250}[rank]
	level := [10]int{0, 4, 10, 18, 26, 34, 42, 50, 58, 60}[rank]

	return core.SpellConfig{
		ActionID:      core.ActionID{SpellID: spellId},
		SpellSchool:   core.SpellSchoolNature,
		ProcMask:      core.ProcMaskEmpty,
		Flags:         core.SpellFlagAPL | core.SpellFlagPureDot,
		Rank:          rank,
		RequiredLevel: level,
		MissileSpeed:  24,

		ManaCost: core.ManaCostOptions{
			FlatCost:   manaCost,
			Multiplier: 1 - 0.02*float64(hunter.Talents.Efficiency),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true, // Hunter GCD is locked at 1.5s
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return hunter.DistanceFromTarget >= 8
		},

		// Need to specially apply LethalShots here, because this spell uses an empty proc mask
		BonusCritRating: 1 * core.CritRatingPerCritChance * float64(hunter.Talents.LethalShots),

		DamageMultiplierAdditive: 1 + 0.02*float64(hunter.Talents.ImprovedSerpentSting),
		CritMultiplier:           hunter.critMultiplier(true, hunter.CurrentTarget),
		ThreatMultiplier:         1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "SerpentSting" + hunter.Label + strconv.Itoa(rank),
				Tag:   "SerpentSting",
			},
			NumberOfTicks: 5,
			TickLength:    time.Second * 3,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.SnapshotBaseDamage = baseDamage + spellCoeff*dot.Spell.SpellPower()
				if !isRollover {
					attackTable := dot.Spell.Unit.AttackTables[target.UnitIndex]
					dot.SnapshotCritChance = dot.Spell.PhysicalCritChance(attackTable)
					dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(attackTable)
				}
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeRangedHit)

			spell.WaitTravelTime(sim, func(s *core.Simulation) {
				if result.Landed() {
					spell.SpellMetrics[target.UnitIndex].Hits--
					spell.Dot(target).Apply(sim)
				}
				spell.DealOutcome(sim, result)
			})
		},
	}
}

func (hunter *Hunter) chimeraShotSerpentStingSpell(rank int) *core.Spell {
	baseDamage := [10]float64{0, 20, 40, 80, 140, 210, 290, 385, 490, 555}[rank]
	spellCoeff := [10]float64{0, .08, .125, .185, .2, .2, .2, .2, .2, .2}[rank]

	return hunter.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 409493},
		SpellSchool: core.SpellSchoolNature,
		ProcMask:    core.ProcMaskRangedSpecial,
		Flags:       core.SpellFlagMeleeMetrics,

		DamageMultiplierAdditive: 1 + 0.02*float64(hunter.Talents.ImprovedSerpentSting),
		DamageMultiplier:         0.4,
		CritMultiplier:           hunter.critMultiplier(true, hunter.CurrentTarget),
		ThreatMultiplier:         1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := baseDamage + spellCoeff*spell.SpellPower()
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeRangedCritOnly)
		},
	})
}

const SERPENT_STING_MAX_RANK = 9

func (hunter *Hunter) registerSerpentStingSpell() {
	for i := 1; i <= SERPENT_STING_MAX_RANK; i++ {
		config := hunter.getSerpentStingConfig(i)

		if config.RequiredLevel <= int(hunter.Level) {
			hunter.highestSerpentStingRank = i
			hunter.SerpentSting = hunter.GetOrRegisterSpell(config)
		}
	}
}
