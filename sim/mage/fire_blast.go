package mage

import (
	"time"

	"github.com/wowsims/classic/sim/core"
)

func (mage *Mage) getFireBlastBaseConfig(rank int, cdTimer *core.Timer) core.SpellConfig {
	spellCoeff := [8]float64{0, .204, .332, .429, .429, .429, .429, .429}[rank]
	baseDamage := [8][]float64{{0}, {27, 35}, {62, 76}, {110, 134}, {177, 211}, {253, 301}, {345, 407}, {446, 524}}[rank]
	spellId := [8]int32{0, 2136, 2137, 2138, 8412, 8413, 10197, 10199}[rank]
	manaCost := [8]float64{0, 40, 75, 115, 165, 220, 280, 340}[rank]
	level := [8]int{0, 6, 14, 22, 30, 38, 46, 54}[rank]

	return core.SpellConfig{
		ActionID:      core.ActionID{SpellID: spellId},
		SpellSchool:   core.SpellSchoolFire,
		ProcMask:      core.ProcMaskSpellDamage,
		Flags:         SpellFlagMage | core.SpellFlagAPL,
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
				Duration: time.Second*8 - time.Millisecond*500*time.Duration(mage.Talents.ImprovedFireBlast),
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   mage.DefaultSpellCritMultiplier(),
		ThreatMultiplier: 1 - 0.15*float64(mage.Talents.BurningSoul),
		ExpectedInitialDamage: func(sim *core.Simulation, target *core.Unit, spell *core.Spell, _ bool) *core.SpellResult {
			baseDamageCacl := (baseDamage[0]+baseDamage[1])/2 + spellCoeff*spell.SpellPower()
			return spell.CalcDamage(sim, target, baseDamageCacl, spell.OutcomeExpectedMagicHitAndCrit)
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(baseDamage[0], baseDamage[1]) + spellCoeff*spell.SpellPower()

			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
		},
	}
}

func (mage *Mage) registerFireBlastSpell() {
	maxRank := 7
	cdTimer := mage.NewTimer()

	for i := 1; i <= maxRank; i++ {
		config := mage.getFireBlastBaseConfig(i, cdTimer)

		if config.RequiredLevel <= int(mage.Level) {
			mage.FireBlast = mage.GetOrRegisterSpell(config)
		}
	}
}
