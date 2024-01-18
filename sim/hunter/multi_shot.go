package hunter

import (
	"time"

	"github.com/wowsims/sod/sim/core"
)

func (hunter *Hunter) getMultiShotConfig(rank int, timer *core.Timer) core.SpellConfig {
	spellId := [6]int32{0, 2643, 14288, 14289, 14290, 25294}[rank]
	baseDamage := [6]float64{0, 0, 40, 80, 120, 150}[rank]
	manaCost := [6]float64{0, 100, 140, 175, 210, 230}[rank]
	level := [6]int{0, 18, 30, 42, 54, 60}[rank]

	numHits := min(3, hunter.Env.GetNumTargets())

	return core.SpellConfig{
		ActionID:      core.ActionID{SpellID: spellId},
		SpellSchool:   core.SpellSchoolPhysical,
		ProcMask:      core.ProcMaskRangedSpecial,
		Flags:         core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | core.SpellFlagAPL,
		Rank:          rank,
		RequiredLevel: level,
		MissileSpeed:  24,

		ManaCost: core.ManaCostOptions{
			FlatCost:   manaCost,
			Multiplier: 1 - 0.03*float64(hunter.Talents.Efficiency),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * 500,
			},
			ModifyCast: func(_ *core.Simulation, spell *core.Spell, cast *core.Cast) {
				cast.CastTime = spell.CastTime()
			},
			IgnoreHaste: true, // Hunter GCD is locked at 1.5s
			CD: core.Cooldown{
				Timer:    timer,
				Duration: time.Second * 10,
			},
			CastTime: func(spell *core.Spell) time.Duration {
				return time.Duration(float64(spell.DefaultCast.CastTime) / hunter.RangedSwingSpeed())
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return hunter.DistanceFromTarget >= 8
		},

		BonusCritRating: 0,
		DamageMultiplierAdditive: 1 +
			.05*float64(hunter.Talents.Barrage),
		DamageMultiplier: 1,
		CritMultiplier:   hunter.critMultiplier(true, hunter.CurrentTarget),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			curTarget := target

			//hunter.AutoAttacks.DelayRangedUntil(sim, sim.CurrentTime+time.Duration(float64(time.Millisecond*500)/hunter.RangedSwingSpeed()))

			sharedDmg := hunter.AutoAttacks.Ranged().BaseDamage(sim) +
				hunter.AmmoDamageBonus +
				spell.BonusWeaponDamage() +
				baseDamage

			for hitIndex := int32(0); hitIndex < numHits; hitIndex++ {
				baseDamage := sharedDmg + 0.2*spell.RangedAttackPower(curTarget)

				result := spell.CalcDamage(sim, curTarget, baseDamage, spell.OutcomeRangedHitAndCrit)
				spell.WaitTravelTime(sim, func(s *core.Simulation) {
					spell.DealDamage(sim, result)
				})

				curTarget = sim.Environment.NextTargetUnit(curTarget)
			}
		},
	}
}

func (hunter *Hunter) registerMultiShotSpell(timer *core.Timer) {
	maxRank := 5

	for i := 1; i <= maxRank; i++ {
		config := hunter.getMultiShotConfig(i, timer)

		if config.RequiredLevel <= int(hunter.Level) {
			hunter.ArcaneShot = hunter.GetOrRegisterSpell(config)
		}
	}
}
