package druid

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/items"
	"github.com/wowsims/wotlk/sim/core/stats"
)

const IdolAvenger int32 = 31025

func (druid *Druid) registerWrathSpell() {
	druid.OriginalWrathDamageMultiplier = (1 + 0.02*float64(druid.Talents.Moonfury)) * (1 + 0.01*float64(druid.Talents.ImprovedInsectSwarm))
	iffCritBonus := core.TernaryFloat64(druid.CurrentTarget.HasAura("Improved Faerie Fire"), float64(druid.Talents.ImprovedFaerieFire)*1*core.CritRatingPerCritChance, 0)

	baseCost := 0.11 * druid.BaseMana
	minBaseDamage := 557.0
	maxBaseDamage := 627.0
	actionID := core.ActionID{SpellID: 26985}
	manaMetrics := druid.NewManaMetrics(actionID)

	// This seems to be unaffected by wrath of cenarius.
	bonusFlatDamage := core.TernaryFloat64(druid.Equip[items.ItemSlotRanged].ID == IdolAvenger, 25*0.571, 0)

	druid.Wrath = druid.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolNature,

		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost:     baseCost * (1 - 0.03*float64(druid.Talents.Moonglow)),
				GCD:      core.GCDDefault,
				CastTime: time.Second*2 - (time.Millisecond * 100 * time.Duration(druid.Talents.StarlightWrath)),
			},

			ModifyCast: func(_ *core.Simulation, _ *core.Spell, cast *core.Cast) {
				//druid.applyNaturesGrace(cast)
				druid.applyNaturesSwiftness(cast)
			},
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:             core.ProcMaskSpellDamage,
			BonusSpellCritRating: float64(2*float64(druid.Talents.NaturesMajesty)*core.CritRatingPerCritChance) + iffCritBonus,
			DamageMultiplier:     druid.OriginalWrathDamageMultiplier,
			ThreatMultiplier:     1,

			BaseDamage:     core.BaseDamageConfigMagic(minBaseDamage+bonusFlatDamage, maxBaseDamage+bonusFlatDamage, 0.571+0.02*float64(druid.Talents.WrathOfCenarius)),
			OutcomeApplier: druid.OutcomeFuncMagicHitAndCrit(druid.SpellCritMultiplier(1, 0.2*float64(druid.Talents.Vengeance))),
			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spellEffect.Outcome.Matches(core.OutcomeCrit) {
					hasMoonkinForm := core.TernaryFloat64(druid.Talents.MoonkinForm, 1, 0)
					druid.AddMana(sim, druid.MaxMana()*0.02*hasMoonkinForm, manaMetrics, true)
				}
			},
		}),
	})
}
