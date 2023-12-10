package warlock

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (warlock *Warlock) getShadowCleaveBaseConfig(rank int) core.SpellConfig {
	spellCoeff := [11]float64{0, .047, .1, .187, .286, .286, .286, .286, .286, .286, .286}[rank]
	baseDamage := [11][]float64{{0}, {3, 7}, {7, 12}, {14, 23}, {26, 39}, {42, 64}, {60, 91}, {82, 124}, {105, 158}, {128, 193}, {136, 204}}[rank]
	spellId := [11]int32{0, 403835, 403839, 403840, 403841, 403842, 403843, 403844, 403848, 403851, 403852}[rank]
	manaCost := [11]float64{0, 12, 20, 35, 55, 80, 105, 132, 157, 185, 190}[rank]
	level := [11]int{0, 1, 6, 12, 20, 28, 36, 44, 52, 60, 60}[rank]

	numHits := min(4, warlock.Env.GetNumTargets())
	results := make([]*core.SpellResult, numHits)

	return core.SpellConfig{
		ActionID:      core.ActionID{SpellID: spellId},
		SpellSchool:   core.SpellSchoolShadow,
		ProcMask:      core.ProcMaskSpellDamage,
		Flags:         core.SpellFlagAPL | core.SpellFlagResetAttackSwing,
		RequiredLevel: level,
		Rank:          rank,

		ManaCost: core.ManaCostOptions{
			FlatCost: manaCost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    warlock.NewTimer(),
				Duration: time.Second * 6,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return warlock.MetamorphosisAura.IsActive()
		},

		BonusCritRating:  float64(warlock.Talents.Devastation) * core.SpellCritRatingPerCritChance,
		DamageMultiplier: 1,
		CritMultiplier:   warlock.SpellCritMultiplier(1, core.TernaryFloat64(warlock.Talents.Ruin, 1, 0)),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			curTarget := target
			for hitIndex := int32(0); hitIndex < numHits; hitIndex++ {
				baseDamage := sim.Roll(baseDamage[0], baseDamage[1]) + spellCoeff*spell.SpellPower()
				results[hitIndex] = spell.CalcDamage(sim, curTarget, baseDamage, spell.OutcomeMagicHitAndCrit)

				curTarget = sim.Environment.NextTargetUnit(curTarget)
			}

			curTarget = target
			for hitIndex := int32(0); hitIndex < numHits; hitIndex++ {
				spell.DealDamage(sim, results[hitIndex])

				if results[hitIndex].Landed() {
					warlock.EverlastingAfflictionRefresh(sim, curTarget)

					if warlock.Talents.ImprovedShadowBolt > 0 && results[hitIndex].DidCrit() {
						impShadowBoltAura := warlock.ImprovedShadowBoltAuras.Get(curTarget)
						impShadowBoltAura.Activate(sim)
						impShadowBoltAura.SetStacks(sim, 4)
					}
				}

				curTarget = sim.Environment.NextTargetUnit(curTarget)
			}
		},
	}
}

func (warlock *Warlock) registerShadowCleaveSpell() {
	if !warlock.HasRune(proto.WarlockRune_RuneHandsMetamorphosis) {
		return
	}

	maxRank := 10

	for i := 1; i <= maxRank; i++ {
		config := warlock.getShadowCleaveBaseConfig(i)

		if config.RequiredLevel <= int(warlock.Level) {
			warlock.ShadowCleave = warlock.GetOrRegisterSpell(config)
		}
	}
}
