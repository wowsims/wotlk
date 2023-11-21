package priest

import (
	"time"

	"github.com/wowsims/classic/sim/core"
)

func (priest *Priest) getSmiteBaseConfig(rank int) core.SpellConfig {
	spellCoeff := [9]float64{0, 0.123, 0.271, 0.554, 0.714, 0.714, 0.714, 0.714, 0.714}[rank]
	baseDamage := [9][]float64{{0}, {15, 20}, {28, 34}, {58, 67}, {97, 112}, {158, 178}, {222, 250}, {298, 335}, {384, 429}}[rank]
	spellId := [9]int32{0, 585, 591, 598, 984, 1004, 6060, 10933, 10934}[rank]
	manaCost := [9]float64{0, 20, 30, 60, 95, 140, 185, 230, 280}[rank]
	level := [9]int{0, 1, 6, 14, 22, 30, 38, 46, 54}[rank]
	castTime := [9]int{0, 1500, 2000, 2500, 2500, 2500, 2500, 2500, 2500}[rank]

	return core.SpellConfig{
		ActionID:      core.ActionID{SpellID: spellId},
		SpellSchool:   core.SpellSchoolHoly,
		ProcMask:      core.ProcMaskSpellDamage,
		Flags:         core.SpellFlagAPL,
		RequiredLevel: level,
		Rank:          rank,

		ManaCost: core.ManaCostOptions{
			FlatCost: manaCost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond*time.Duration(castTime) - time.Millisecond*100*time.Duration(priest.Talents.DivineFury),
			},
		},

		BonusCritRating:  float64(priest.Talents.HolySpecialization) * 1 * core.CritRatingPerCritChance,
		DamageMultiplier: 1 + 0.05*float64(priest.Talents.SearingLight),
		CritMultiplier:   priest.SpellCritMultiplier(1+0.01*float64(priest.Talents.HolySpecialization), 1),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(baseDamage[0], baseDamage[1]) + spellCoeff*spell.SpellPower()
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
		},
	}
}

func (priest *Priest) RegisterSmiteSpell() {
	maxRank := 8

	for i := 1; i <= maxRank; i++ {
		config := priest.getSmiteBaseConfig(i)

		if config.RequiredLevel <= int(priest.Level) {
			priest.Smite = priest.GetOrRegisterSpell(config)
		}
	}
}
