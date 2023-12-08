package warlock

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

// TODO: Classic warlock verify chaos bolt mechanics
func (warlock *Warlock) registerChaosBoltSpell() {
	if !warlock.HasRune(proto.WarlockRune_RuneHandsChaosBolt) {
		return
	}
	spellCoeff := 0.714
	level := float64(warlock.GetCharacter().Level)
	baseCalc := (6.568597 + 0.672028*level + 0.031721*level*level)
	baseLowDamage := baseCalc * 5.22
	baseHighDamage := baseCalc * 6.62

	warlock.ChaosBolt = warlock.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 403629},
		SpellSchool: core.SpellSchoolFire,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.07,
			Multiplier: 1,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond*2500 - (time.Millisecond * 100 * time.Duration(warlock.Talents.Bane)),
			},
			CD: core.Cooldown{
				Timer:    warlock.NewTimer(),
				Duration: time.Second * 12,
			},
		},

		BonusCritRating: 5 * float64(warlock.Talents.Devastation) * core.SpellCritRatingPerCritChance,
		// Assuming 100% hit for all target levels, numbers could be updated for level comparison later
		BonusHitRating:           100 * core.SpellCritRatingPerCritChance,
		DamageMultiplierAdditive: 1,
		CritMultiplier:           warlock.SpellCritMultiplier(1, core.TernaryFloat64(warlock.Talents.Ruin, 1, 0)),
		ThreatMultiplier:         1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(baseLowDamage, baseHighDamage) + spellCoeff*spell.SpellPower()
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicCrit)
		},
	})
}
