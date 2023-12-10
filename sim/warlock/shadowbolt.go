package warlock

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (warlock *Warlock) getShadowBoltBaseConfig(rank int) core.SpellConfig {
	spellCoeff := [11]float64{0, .14, .299, .56, .857, .857, .857, .857, .857, .857, .857}[rank]
	baseDamage := [11][]float64{{0}, {13, 18}, {26, 32}, {52, 61}, {92, 104}, {150, 170}, {213, 240}, {292, 327}, {373, 415}, {455, 507}, {482, 538}}[rank]
	spellId := [11]int32{0, 686, 695, 705, 1088, 1106, 7641, 11699, 11660, 11661, 25307}[rank]
	manaCost := [11]float64{0, 25, 40, 70, 110, 160, 210, 265, 315, 370, 380}[rank]
	level := [11]int{0, 1, 6, 12, 20, 28, 36, 44, 52, 60, 60}[rank]
	castTime := [11]int32{0, 1700, 2200, 2800, 3000, 3000, 3000, 3000, 3000, 3000, 3000}[rank]

	shadowboltVolley := warlock.HasRune(proto.WarlockRune_RuneHandsShadowBoltVolley)
	damageMulti := core.TernaryFloat64(shadowboltVolley, 0.8, 1.0)
	numHits := min(core.TernaryInt32(shadowboltVolley, 5, 1), warlock.Env.GetNumTargets())
	results := make([]*core.SpellResult, numHits)

	return core.SpellConfig{
		ActionID:      core.ActionID{SpellID: spellId},
		SpellSchool:   core.SpellSchoolShadow,
		ProcMask:      core.ProcMaskSpellDamage,
		Flags:         core.SpellFlagAPL | core.SpellFlagResetAttackSwing,
		RequiredLevel: level,
		Rank:          rank,

		ManaCost: core.ManaCostOptions{
			FlatCost:   manaCost,
			Multiplier: 1 - float64(warlock.Talents.Cataclysm)*0.01,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * time.Duration(castTime-100*warlock.Talents.Bane),
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return warlock.MetamorphosisAura == nil || !warlock.MetamorphosisAura.IsActive()
		},

		BonusCritRating:  float64(warlock.Talents.Devastation) * core.SpellCritRatingPerCritChance,
		DamageMultiplier: damageMulti,
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

func (warlock *Warlock) registerShadowBoltSpell() {
	maxRank := 10

	for i := 1; i <= maxRank; i++ {
		config := warlock.getShadowBoltBaseConfig(i)

		if config.RequiredLevel <= int(warlock.Level) {
			warlock.ShadowBolt = warlock.GetOrRegisterSpell(config)
		}
	}
}
