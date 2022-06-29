package druid

import (
	"time"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/items"
	"github.com/wowsims/tbc/sim/core/stats"
)

// Idol IDs
const IvoryMoongoddess int32 = 27518

func (druid *Druid) newStarfireSpell(rank int) *core.Spell {
	actionID := core.ActionID{SpellID: 26986}
	baseCost := 370.0
	minBaseDamage := 550.0
	maxBaseDamage := 647.0
	spellCoefficient := 1.0
	if rank == 6 {
		actionID = core.ActionID{SpellID: 9876}
		baseCost = 315
		minBaseDamage = 463
		maxBaseDamage = 543
		spellCoefficient = 0.99
	}

	// This seems to be unaffected by wrath of cenarius so it needs to come first.
	bonusFlatDamage := core.TernaryFloat64(druid.Equip[items.ItemSlotRanged].ID == IvoryMoongoddess, 55*spellCoefficient, 0)
	spellCoefficient += 0.04 * float64(druid.Talents.WrathOfCenarius)

	effect := core.SpellEffect{
		ProcMask:             core.ProcMaskSpellDamage,
		BonusSpellCritRating: (float64(druid.Talents.FocusedStarlight) * 2 * core.SpellCritRatingPerCritChance) + core.TernaryFloat64(ItemSetThunderheartRegalia.CharacterHasSetBonus(&druid.Character, 4), 5*core.SpellCritRatingPerCritChance, 0),
		DamageMultiplier:     1 + 0.02*float64(druid.Talents.Moonfury),
		ThreatMultiplier:     1,
		BaseDamage:           core.BaseDamageConfigMagic(minBaseDamage+bonusFlatDamage, maxBaseDamage+bonusFlatDamage, spellCoefficient),
		OutcomeApplier:       druid.OutcomeFuncMagicHitAndCrit(druid.SpellCritMultiplier(1, 0.2*float64(druid.Talents.Vengeance))),
	}

	if ItemSetNordrassilRegalia.CharacterHasSetBonus(&druid.Character, 4) {
		effect.BaseDamage = core.WrapBaseDamageConfig(effect.BaseDamage, func(oldCalculator core.BaseDamageCalculator) core.BaseDamageCalculator {
			return func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
				normalDamage := oldCalculator(sim, hitEffect, spell)

				// Check if moonfire/insectswarm is ticking on the target.
				// TODO: in a raid simulator we need to be able to see which dots are ticking from other druids.
				if druid.MoonfireDot.IsActive() || druid.InsectSwarmDot.IsActive() {
					return normalDamage * 1.1
				} else {
					return normalDamage
				}
			}
		})
	}

	return druid.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolArcane,

		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost:     baseCost * (1 - 0.03*float64(druid.Talents.Moonglow)),
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond*3500 - (time.Millisecond * 100 * time.Duration(druid.Talents.StarlightWrath)),
			},

			ModifyCast: func(_ *core.Simulation, _ *core.Spell, cast *core.Cast) {
				druid.applyNaturesGrace(cast)
				druid.applyNaturesSwiftness(cast)
			},
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(effect),
	})
}
