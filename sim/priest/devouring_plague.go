package priest

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (priest *Priest) registerDevouringPlagueSpell() {
	actionID := core.ActionID{SpellID: 48300}
	baseCost := priest.BaseMana * 0.25
	target := priest.CurrentTarget

	applier := priest.OutcomeFuncTick()
	if priest.Talents.Shadowform {
		applier = priest.OutcomeFuncMagicCrit()
	}

	priest.DevouringPlague = priest.RegisterSpell(core.SpellConfig{
		ActionID:     actionID,
		SpellSchool:  core.SpellSchoolShadow,
		ProcMask:     core.ProcMaskSpellDamage,
		Flags:        core.SpellFlagDisease,
		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost * (1 - []float64{0, .04, .07, .10}[priest.Talents.MentalAgility]),
				GCD:  core.GCDDefault,
			},
		},

		BonusHitRating: float64(priest.Talents.ShadowFocus) * 1 * core.SpellHitRatingPerHitChance,
		BonusCritRating: 0 +
			3*float64(priest.Talents.MindMelt)*core.CritRatingPerCritChance +
			core.TernaryFloat64(priest.HasSetBonus(ItemSetCrimsonAcolyte, 2), 5, 0)*core.CritRatingPerCritChance,
		DamageMultiplier: 1 +
			float64(priest.Talents.Darkness)*0.02 +
			float64(priest.Talents.TwinDisciplines)*0.01 +
			float64(priest.Talents.ImprovedDevouringPlague)*0.05 +
			core.TernaryFloat64(priest.HasSetBonus(ItemSetConquerorSanct, 2), 0.15, 0),
		CritMultiplier:   priest.DefaultSpellCritMultiplier(),
		ThreatMultiplier: 1 - 0.05*float64(priest.Talents.ShadowAffinity),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := (1376/8 + 0.1849*spell.SpellPower()) *
				(8 * 0.1 * float64(priest.Talents.ImprovedDevouringPlague)) *
				(1 + 0.02*float64(priest.ShadowWeavingAura.GetStacks()))

			result := spell.CalcDamageMagicHitAndCrit(sim, target, baseDamage)
			if result.Landed() {
				priest.AddShadowWeavingStack(sim)
				priest.DevouringPlagueDot.Apply(sim)
			}
			spell.DealDamage(sim, &result)
		},
	})

	priest.DevouringPlagueDot = core.NewDot(core.Dot{
		Spell: priest.RegisterSpell(core.SpellConfig{
			ActionID:    actionID,
			SpellSchool: core.SpellSchoolShadow,
			ProcMask:    core.ProcMaskSpellDamage,
			Flags:       core.SpellFlagDisease,

			BonusHitRating: float64(priest.Talents.ShadowFocus) * 1 * core.SpellHitRatingPerHitChance,
			BonusCritRating: 0 +
				3*float64(priest.Talents.MindMelt)*core.CritRatingPerCritChance +
				core.TernaryFloat64(priest.HasSetBonus(ItemSetCrimsonAcolyte, 2), 5, 0)*core.CritRatingPerCritChance,
			DamageMultiplier: 1 +
				float64(priest.Talents.Darkness)*0.02 +
				float64(priest.Talents.TwinDisciplines)*0.01 +
				float64(priest.Talents.ImprovedDevouringPlague)*0.05 +
				core.TernaryFloat64(priest.HasSetBonus(ItemSetConquerorSanct, 2), 0.15, 0),
			CritMultiplier:   priest.SpellCritMultiplier(1, 1),
			ThreatMultiplier: 1 - 0.05*float64(priest.Talents.ShadowAffinity),
		}),
		Aura: target.RegisterAura(core.Aura{
			Label:    "DevouringPlague-" + strconv.Itoa(int(priest.Index)),
			ActionID: actionID,
		}),

		NumberOfTicks:       8,
		TickLength:          time.Second * 3,
		AffectedByCastSpeed: priest.Talents.Shadowform,

		TickEffects: core.TickFuncSnapshot(target, core.SpellEffect{
			IsPeriodic: true,

			BaseDamage: core.WrapBaseDamageConfig(
				core.BaseDamageConfigMagicNoRoll(1376/8, 0.1849),
				func(oldCalculator core.BaseDamageCalculator) core.BaseDamageCalculator {
					return func(sim *core.Simulation, spellEffect *core.SpellEffect, spell *core.Spell) float64 {
						swMod := 1 + float64(priest.ShadowWeavingAura.GetStacks())*0.02
						dmg := oldCalculator(sim, spellEffect, spell)

						return dmg * swMod
					}
				}),
			OutcomeApplier: applier,
		}),
	})
}
