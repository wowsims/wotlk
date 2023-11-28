package priest

import (
	"time"

	"github.com/wowsims/classic/sod/sim/core"
	"github.com/wowsims/classic/sod/sim/core/proto"
)

func (priest *Priest) registerShadowWordDeathSpell() {
	if !priest.HasRune(proto.PriestRune_RuneHandsShadowWordDeath) {
		return
	}
	spellCoeff := 0.429

	level := float64(priest.GetCharacter().Level)
	baseCalc := (9.456667 + 0.635108*level + 0.039063*level*level)
	baseLowDamage := baseCalc * 2.61
	baseHightDamage := baseCalc * 3.05

	priest.ShadowWordDeath = priest.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 401955},
		SpellSchool: core.SpellSchoolShadow,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.12,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    priest.NewTimer(),
				Duration: time.Second * 12,
			},
		},

		BonusHitRating:   float64(priest.Talents.ShadowFocus) * 2 * core.SpellHitRatingPerHitChance,
		BonusCritRating:  0,
		DamageMultiplier: 1,
		CritMultiplier:   1,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(baseLowDamage, baseHightDamage) + spellCoeff*spell.SpellPower()
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)

			if result.Landed() {
				priest.AddShadowWeavingStack(sim)
			}
			spell.DealDamage(sim, result)
		},
		ExpectedInitialDamage: func(sim *core.Simulation, target *core.Unit, spell *core.Spell, _ bool) *core.SpellResult {
			baseDamage := (baseLowDamage+baseHightDamage)/2 + spellCoeff*spell.SpellPower()
			return spell.CalcDamage(sim, target, baseDamage, spell.OutcomeExpectedMagicHitAndCrit)
		},
	})
}
