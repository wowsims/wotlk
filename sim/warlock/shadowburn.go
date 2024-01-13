package warlock

import (
	"time"

	"github.com/wowsims/sod/sim/core"
)

func (warlock *Warlock) registerShadowBurnBaseConfig(rank int) core.SpellConfig {
	spellId := [7]int32{0, 17877, 18867, 18868, 18869, 18870, 18871}[rank]
	baseDamage := [7][]float64{{0}, {91, 104}, {123, 140}, {196, 221}, {274, 307}, {365, 408}, {462, 514}}[rank]
	manaCost := [7]float64{0, 105, 130, 190, 245, 305, 365}[rank]
	level := [7]int{0, 15, 24, 32, 40, 48, 56}[rank]

	spellCoeff := 0.429

	return core.SpellConfig{
		ActionID:      core.ActionID{SpellID: spellId},
		SpellSchool:   core.SpellSchoolShadow,
		ProcMask:      core.ProcMaskSpellDamage,
		Flags:         core.SpellFlagAPL | core.SpellFlagResetAttackSwing,
		RequiredLevel: level,
		Rank:          rank,

		ManaCost: core.ManaCostOptions{
			BaseCost: manaCost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    warlock.NewTimer(),
				Duration: time.Second * time.Duration(15),
			},
		},

		BonusCritRating:  float64(warlock.Talents.Devastation) * core.SpellCritRatingPerCritChance,
		DamageMultiplier: 1,
		CritMultiplier:   warlock.SpellCritMultiplier(1, core.TernaryFloat64(warlock.Talents.Ruin, 1, 0)),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(baseDamage[0], baseDamage[1]) + spellCoeff*spell.SpellPower()
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
		},
	}
}

func (warlock *Warlock) registerShadowBurnSpell() {
	if !warlock.Talents.Shadowburn {
		return
	}

	maxRank := 6

	for i := 1; i <= maxRank; i++ {
		config := warlock.registerShadowBurnBaseConfig(i)

		if config.RequiredLevel <= int(warlock.Level) {
			warlock.ShadowBolt = warlock.GetOrRegisterSpell(config)
		}
	}
}
