package priest

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (priest *Priest) registerCircleOfHealingSpell() {
	if !priest.Talents.CircleOfHealing {
		return
	}

	baseCost := .21 * priest.BaseMana

	baseEffect := core.SpellEffect{
		IsHealing: true,
		ProcMask:  core.ProcMaskSpellDamage,

		BonusCritRating: float64(priest.Talents.HolySpecialization) * 1 * core.CritRatingPerCritChance,
		DamageMultiplier: 1 *
			(1 + .02*float64(priest.Talents.DivineProvidence)),
		ThreatMultiplier: 1 - []float64{0, .07, .14, .20}[priest.Talents.SilentResolve],

		BaseDamage:     core.BaseDamageConfigHealing(958, 1058, 0.4029),
		OutcomeApplier: priest.OutcomeFuncHealingCrit(priest.DefaultSpellCritMultiplier()),
	}

	var effects []core.SpellEffect
	targets := priest.Env.Raid.GetFirstNPlayersOrPets(5)
	for _, target := range targets {
		effect := baseEffect
		effect.Target = target
		effects = append(effects, effect)
	}

	priest.CircleOfHealing = priest.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 48089},
		SpellSchool: core.SpellSchoolHoly,

		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost * (1 - []float64{0, .04, .07, .10}[priest.Talents.MentalAgility]),
				GCD:  core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    priest.NewTimer(),
				Duration: time.Second * 6,
			},
		},

		ApplyEffects: core.ApplyEffectFuncDamageMultiple(effects),
	})
}
