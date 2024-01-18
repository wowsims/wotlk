package hunter

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (hunter *Hunter) getArcaneShotConfig(rank int, timer *core.Timer) core.SpellConfig {
	spellId := [9]int32{0, 3044, 14281, 14282, 14283, 14284, 14285, 14286, 14287}[rank]
	baseDamage := [9]float64{0, 13, 21, 33, 59, 83, 115, 145, 183}[rank]
	spellCoeff := [9]float64{0, .204, .3, .429, .429, .429, .429, .429, .429}[rank]
	manaCost := [9]float64{0, 25, 35, 50, 80, 105, 135, 160, 190}[rank]
	level := [9]int{0, 6, 12, 20, 28, 36, 44, 52, 60}[rank]

	hasCobraStrikes := hunter.pet != nil && hunter.HasRune(proto.HunterRune_RuneChestCobraStrikes)

	manaCostMultiplier := 1 - 0.02*float64(hunter.Talents.Efficiency)
	if hunter.HasRune(proto.HunterRune_RuneChestMasterMarksman) {
		manaCostMultiplier -= 0.25
	}
	return core.SpellConfig{
		ActionID:      core.ActionID{SpellID: spellId},
		SpellSchool:   core.SpellSchoolArcane,
		ProcMask:      core.ProcMaskRangedSpecial,
		Flags:         core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		Rank:          rank,
		RequiredLevel: level,
		MissileSpeed:  24,

		ManaCost: core.ManaCostOptions{
			FlatCost:   manaCost,
			Multiplier: manaCostMultiplier,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    timer,
				Duration: time.Second*6 - time.Millisecond*200*time.Duration(hunter.Talents.ImprovedArcaneShot),
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return hunter.DistanceFromTarget >= 8
		},

		BonusCritRating:  0,
		DamageMultiplier: 1,
		CritMultiplier:   hunter.critMultiplier(true, hunter.CurrentTarget),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := baseDamage + spellCoeff*spell.SpellPower()

			if hunter.SniperTrainingAura.IsActive() {
				spell.BonusCritRating += 10 * core.CritRatingPerCritChance
			}
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeRangedHitAndCrit)
			if hunter.SniperTrainingAura.IsActive() {
				spell.BonusCritRating -= 10 * core.CritRatingPerCritChance
			}

			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)

				if hasCobraStrikes && result.DidCrit() {
					hunter.CobraStrikesAura.Activate(sim)
					hunter.CobraStrikesAura.SetStacks(sim, 2)
				}
			})
		},
	}
}

func (hunter *Hunter) registerArcaneShotSpell(timer *core.Timer) {
	maxRank := 8

	for i := 1; i <= maxRank; i++ {
		config := hunter.getArcaneShotConfig(i, timer)

		if config.RequiredLevel <= int(hunter.Level) {
			hunter.ArcaneShot = hunter.GetOrRegisterSpell(config)
		}
	}
}
