package priest

import (
	"time"

	"github.com/wowsims/classic/sim/core"
)

func (priest *Priest) getMindBlastBaseConfig(rank int, cdTimer *core.Timer) core.SpellConfig {
	spellCoeff := 0.429 // 1.5 cast over 3.5 reference
	baseDamage := [10][]float64{{0}, {42, 46}, {76, 83}, {117, 126}, {174, 184}, {225, 239}, {288, 307}, {356, 377}, {437, 461}, {508, 537}}[rank]
	spellId := [10]int32{0, 8092, 8102, 8103, 8104, 8105, 8106, 10945, 10946, 10947}[rank]
	manaCost := [10]float64{0, 50, 80, 110, 150, 185, 225, 264, 310, 350}[rank]
	level := [10]int{0, 10, 16, 22, 28, 32, 40, 46, 52, 58}[rank]

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
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * 1500,
			},
			CD: core.Cooldown{
				Timer:    cdTimer,
				Duration: time.Second*8 - time.Millisecond*500*time.Duration(priest.Talents.ImprovedMindBlast),
			},
		},

		BonusHitRating:   float64(priest.Talents.ShadowFocus) * 2 * core.SpellHitRatingPerHitChance,
		BonusCritRating:  0,
		DamageMultiplier: 1,
		CritMultiplier:   1,
		ThreatMultiplier: 1 - 0.08*float64(priest.Talents.ShadowAffinity),
		ExpectedInitialDamage: func(sim *core.Simulation, target *core.Unit, spell *core.Spell, _ bool) *core.SpellResult {
			baseDamageCacl := (baseDamage[0]+baseDamage[1])/2 + spellCoeff*spell.SpellPower()
			return spell.CalcDamage(sim, target, baseDamageCacl, spell.OutcomeExpectedMagicHitAndCrit)
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(baseDamage[0], baseDamage[1]) + spellCoeff*spell.SpellPower()
			baseDamage *= priest.MindBlastModifier

			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			if result.Landed() {
				priest.AddShadowWeavingStack(sim)
			}

			spell.DealDamage(sim, result)
		},
	}
}

func (priest *Priest) registerMindBlast() {
	maxRank := 9
	cdTimer := priest.NewTimer()

	for i := 1; i < maxRank; i++ {
		config := priest.getMindBlastBaseConfig(i, cdTimer)

		if config.RequiredLevel <= int(priest.Level) {
			priest.MindBlast = priest.GetOrRegisterSpell(config)
		}
	}
}
