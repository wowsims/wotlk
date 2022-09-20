package priest

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (priest *Priest) registerVampiricTouchSpell() {
	actionID := core.ActionID{SpellID: 48160}
	baseCost := priest.BaseMana * 0.16

	applier := priest.OutcomeFuncTick()
	if priest.Talents.Shadowform {
		applier = priest.OutcomeFuncMagicCrit(priest.SpellCritMultiplier(1, 1))
	}

	priest.VampiricTouch = priest.RegisterSpell(core.SpellConfig{
		ActionID:     actionID,
		SpellSchool:  core.SpellSchoolShadow,
		ProcMask:     core.ProcMaskSpellDamage,
		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost:     baseCost,
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * 1500,
			},
		},

		BonusHitRating:   float64(priest.Talents.ShadowFocus) * 1 * core.SpellHitRatingPerHitChance,
		BonusCritRating:  float64(priest.Talents.MindMelt)*3*core.CritRatingPerCritChance + core.TernaryFloat64(priest.HasSetBonus(ItemSetCrimsonAcolyte, 2), 5, 0)*core.CritRatingPerCritChance,
		DamageMultiplier: 1 + float64(priest.Talents.Darkness)*0.02,
		ThreatMultiplier: 1 - 0.08*float64(priest.Talents.ShadowAffinity),

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			OutcomeApplier: priest.OutcomeFuncMagicHit(),
			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spellEffect.Landed() {
					priest.AddShadowWeavingStack(sim)
					priest.VampiricTouchDot.Apply(sim)
				}
			},
		}),
	})

	target := priest.CurrentTarget

	priest.VampiricTouchDot = core.NewDot(core.Dot{
		Spell: priest.VampiricTouch,
		Aura: target.RegisterAura(core.Aura{
			Label:    "VampiricTouch-" + strconv.Itoa(int(priest.Index)),
			ActionID: actionID,
		}),

		NumberOfTicks:       5 + core.TernaryInt(priest.HasSetBonus(ItemSetZabras, 2), 2, 0),
		TickLength:          time.Second * 3,
		AffectedByCastSpeed: priest.Talents.Shadowform,

		TickEffects: core.TickFuncSnapshot(target, core.SpellEffect{
			IsPeriodic: true,
			BaseDamage: core.WrapBaseDamageConfig(
				core.BaseDamageConfigMagicNoRoll(850/5, 0.4),
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
