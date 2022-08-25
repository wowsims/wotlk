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
		applier = priest.OutcomeFuncMagicCrit(priest.SpellCritMultiplier(1, 1))
	}

	effect := core.SpellEffect{
		DamageMultiplier: 8 * 0.1 * float64(priest.Talents.ImprovedDevouringPlague) *
			(1 + float64(priest.Talents.Darkness)*0.02 + float64(priest.Talents.TwinDisciplines)*0.01 + float64(priest.Talents.ImprovedDevouringPlague)*0.05) *
			core.TernaryFloat64(priest.HasSetBonus(ItemSetConquerorSanct, 2), 1.15, 1),
		BonusSpellHitRating: float64(priest.Talents.ShadowFocus) * 1 * core.SpellHitRatingPerHitChance,
		ThreatMultiplier:    1 - 0.05*float64(priest.Talents.ShadowAffinity),
		BaseDamage: core.WrapBaseDamageConfig(
			core.BaseDamageConfigMagicNoRoll(1376/8, 0.1849),
			func(oldCalculator core.BaseDamageCalculator) core.BaseDamageCalculator {
				return func(sim *core.Simulation, spellEffect *core.SpellEffect, spell *core.Spell) float64 {
					swMod := 1 + float64(priest.ShadowWeavingAura.GetStacks())*0.02
					dmg := oldCalculator(sim, spellEffect, spell)

					return dmg * swMod
				}
			}),
		OutcomeApplier: priest.OutcomeFuncMagicHitAndCrit(priest.DefaultSpellCritMultiplier()),
		OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if spellEffect.Landed() {
				priest.AddShadowWeavingStack(sim)
				priest.DevouringPlagueDot.Apply(sim)
			}
		},
		ProcMask: core.ProcMaskSpellDamage,
	}

	priest.DevouringPlague = priest.RegisterSpell(core.SpellConfig{
		ActionID:     actionID,
		SpellSchool:  core.SpellSchoolShadow,
		Flags:        core.SpellFlagDisease,
		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost * (1 - 0.02*float64(priest.Talents.MentalAgility)),
				GCD:  core.GCDDefault,
			},
		},
		ApplyEffects: core.ApplyEffectFuncDirectDamage(effect),
	})

	priest.DevouringPlagueDot = core.NewDot(core.Dot{
		Spell: priest.DevouringPlague,
		Aura: target.RegisterAura(core.Aura{
			Label:    "DevouringPlague-" + strconv.Itoa(int(priest.Index)),
			ActionID: actionID,
		}),

		NumberOfTicks:       8,
		TickLength:          time.Second * 3,
		AffectedByCastSpeed: priest.Talents.Shadowform,

		TickEffects: core.TickFuncSnapshot(target, core.SpellEffect{
			ProcMask:   core.ProcMaskPeriodicDamage,
			IsPeriodic: true,

			BonusSpellCritRating: 0 +
				3*float64(priest.Talents.MindMelt)*core.CritRatingPerCritChance +
				core.TernaryFloat64(priest.HasSetBonus(ItemSetCrimsonAcolyte, 2), 5, 0)*core.CritRatingPerCritChance,

			DamageMultiplier: 1 +
				float64(priest.Talents.Darkness)*0.02 +
				float64(priest.Talents.TwinDisciplines)*0.01 +
				float64(priest.Talents.ImprovedDevouringPlague)*0.05 +
				core.TernaryFloat64(priest.HasSetBonus(ItemSetConquerorSanct, 2), 0.15, 0),
			ThreatMultiplier: 1 - 0.08*float64(priest.Talents.ShadowAffinity),

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
