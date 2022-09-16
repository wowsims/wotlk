package druid

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/items"
	"github.com/wowsims/wotlk/sim/core/stats"
)

const IdolAvenger int32 = 31025
const IdolSteadfastRenewal int32 = 40712

func (druid *Druid) registerWrathSpell() {
	baseCost := 0.11 * druid.BaseMana

	actionID := core.ActionID{SpellID: 48461}
	manaMetrics := druid.NewManaMetrics(core.ActionID{SpellID: 24858})
	spellCoefficient := 0.571 * (1 + 0.02*float64(druid.Talents.WrathOfCenarius))

	effect := core.SpellEffect{
		ProcMask: core.ProcMaskSpellDamage,

		DamageMultiplier: (1 + druid.TalentsBonuses.moonfuryMultiplier) *
			core.TernaryFloat64(druid.SetBonuses.balance_t9_4, 1.04, 1), // T9-4P

		OutcomeApplier: druid.OutcomeFuncMagicHitAndCrit(druid.SpellCritMultiplier(1, druid.TalentsBonuses.vengeanceModifier)),
		OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if spellEffect.Outcome.Matches(core.OutcomeCrit) {
				hasMoonkinForm := core.TernaryFloat64(druid.Talents.MoonkinForm, 1, 0)
				druid.AddMana(sim, druid.MaxMana()*0.02*hasMoonkinForm, manaMetrics, true)
				if druid.SetBonuses.balance_t10_4 {
					if druid.LasherweaveDot.IsActive() {
						druid.LasherweaveDot.Refresh(sim)
					} else {
						druid.LasherweaveDot.Apply(sim)
					}
				}
			}
			if sim.RandomFloat("Swift Starfire proc") > 0.85 && druid.SetBonuses.balance_pvp_4 {
				druid.SwiftStarfireAura.Activate(sim)
			}
		},
	}

	// Idols
	bonusFlatDamage := core.TernaryFloat64(druid.Equip[items.ItemSlotRanged].ID == IdolAvenger, 25, 0)
	bonusFlatDamage += core.TernaryFloat64(druid.Equip[items.ItemSlotRanged].ID == IdolSteadfastRenewal, 70, 0)
	effect.BaseDamage = core.BaseDamageConfigMagic(557.0+bonusFlatDamage, 627.0+bonusFlatDamage, spellCoefficient)

	druid.Wrath = druid.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolNature,

		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost:     baseCost * druid.TalentsBonuses.moonglowMultiplier,
				GCD:      core.GCDDefault,
				CastTime: time.Second*2 - druid.TalentsBonuses.starlightWrathModifier,
			},

			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				druid.applyNaturesSwiftness(cast)
				druid.ApplyClearcasting(sim, spell, cast)
			},
		},

		BonusCritRating: 0 +
			druid.TalentsBonuses.naturesMajestyBonusCrit +
			core.TernaryFloat64(druid.SetBonuses.balance_t7_4, 5*core.CritRatingPerCritChance, 0), // T7-4P
		ThreatMultiplier: 1,

		ApplyEffects: core.ApplyEffectFuncDirectDamage(effect),
	})
}
