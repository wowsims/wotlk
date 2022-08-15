package druid

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/items"
	"github.com/wowsims/wotlk/sim/core/stats"
)

const IdolAvenger int32 = 31025

func (druid *Druid) registerWrathSpell() {
	baseCost := 0.11 * druid.BaseMana
	minBaseDamage := 553.0
	maxBaseDamage := 623.0

	// This seems to be unaffected by wrath of cenarius.
	bonusFlatDamage := core.TernaryFloat64(druid.Equip[items.ItemSlotRanged].ID == IdolAvenger, 25*0.571, 0)

	druid.Wrath = druid.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 26985},
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
			BonusSpellCritRating: float64(2 * float64(druid.Talents.NaturesMajesty) * 45.91),
			DamageMultiplier:     1 + 0.02*float64(druid.Talents.Moonfury)*(1+0.01*float64(druid.Talents.ImprovedInsectSwarm)),
			ThreatMultiplier:     1,

			BaseDamage:     core.BaseDamageConfigMagic(minBaseDamage+bonusFlatDamage, maxBaseDamage+bonusFlatDamage, 0.571+0.02*float64(druid.Talents.WrathOfCenarius)),
			OutcomeApplier: druid.OutcomeFuncMagicHitAndCrit(druid.SpellCritMultiplier(1, 0.2*float64(druid.Talents.Vengeance))),
		}),
	})
}
