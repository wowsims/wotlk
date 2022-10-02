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
	actionID := core.ActionID{SpellID: 48461}
	baseCost := 0.11 * druid.BaseMana
	manaMetrics := druid.NewManaMetrics(core.ActionID{SpellID: 24858})
	spellCoeff := 0.571 + (0.02 * float64(druid.Talents.WrathOfCenarius))
	bonusFlatDamage := core.TernaryFloat64(druid.Equip[items.ItemSlotRanged].ID == IdolAvenger, 25, 0) +
		core.TernaryFloat64(druid.Equip[items.ItemSlotRanged].ID == IdolSteadfastRenewal, 70, 0)

	druid.Wrath = druid.RegisterSpell(core.SpellConfig{
		ActionID:     actionID,
		SpellSchool:  core.SpellSchoolNature,
		ProcMask:     core.ProcMaskSpellDamage,
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
		DamageMultiplier: (1 + druid.TalentsBonuses.moonfuryMultiplier) *
			core.TernaryFloat64(druid.SetBonuses.balance_t9_4, 1.04, 1), // T9-4P
		CritMultiplier:   druid.SpellCritMultiplier(1, druid.TalentsBonuses.vengeanceModifier),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := bonusFlatDamage + sim.Roll(557, 627) + spellCoeff*spell.SpellPower()
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			if result.DidCrit() {
				if druid.Talents.MoonkinForm {
					druid.AddMana(sim, 0.02*druid.MaxMana(), manaMetrics, true)
				}
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
			spell.DealDamage(sim, &result)
		},
	})
}
