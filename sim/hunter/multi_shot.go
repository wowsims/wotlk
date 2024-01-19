package hunter

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (hunter *Hunter) getMultiShotConfig(rank int, timer *core.Timer) core.SpellConfig {
	spellId := [6]int32{0, 2643, 14288, 14289, 14290, 25294}[rank]
	baseDamage := [6]float64{0, 0, 40, 80, 120, 150}[rank]
	manaCost := [6]float64{0, 100, 140, 175, 210, 230}[rank]
	level := [6]int{0, 18, 30, 42, 54, 60}[rank]

	numHits := min(3, hunter.Env.GetNumTargets())
	hasCobraStrikes := hunter.pet != nil && hunter.HasRune(proto.HunterRune_RuneChestCobraStrikes)
	hasSerpentSpread := hunter.HasRune(proto.HunterRune_RuneLegsSerpentSpread)

	manaCostMultiplier := 1 - 0.02*float64(hunter.Talents.Efficiency)
	if hunter.HasRune(proto.HunterRune_RuneChestMasterMarksman) {
		manaCostMultiplier -= 0.25
	}
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
			Multiplier: manaCostMultiplier,
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

				if hunter.SniperTrainingAura.IsActive() {
					spell.BonusCritRating += 10 * core.CritRatingPerCritChance
				}
				result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeRangedHitAndCrit)
				if hunter.SniperTrainingAura.IsActive() {
					spell.BonusCritRating -= 10 * core.CritRatingPerCritChance
				}

				spell.WaitTravelTime(sim, func(s *core.Simulation) {
					spell.DealDamage(sim, result)

					if hasCobraStrikes && result.DidCrit() {
						hunter.CobraStrikesAura.Activate(sim)
						hunter.CobraStrikesAura.SetStacks(sim, 2)
					}

					if hasSerpentSpread {
						serpentStingAura := hunter.SerpentSting.Dot(curTarget)
						serpentStingTicks := serpentStingAura.NumberOfTicks
						if serpentStingAura.IsActive() {
							// If less then 2 ticks are left then we refresh with a 2 tick duration
							if serpentStingAura.NumberOfTicks > 3 {
								serpentStingAura.NumberOfTicks = 2
								serpentStingAura.Rollover(sim)
								serpentStingAura.NumberOfTicks = serpentStingTicks
							}
						} else {

							serpentStingAura.NumberOfTicks = 2
							serpentStingAura.Apply(sim)
							serpentStingAura.NumberOfTicks = serpentStingTicks
						}
					}
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
