package druid

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/items"
	"github.com/wowsims/wotlk/sim/core/stats"
)

// Idol IDs
const IvoryMoongoddess int32 = 27518
const ShootingStar int32 = 60775

func (druid *Druid) newStarfireSpell() *core.Spell {

	actionID := core.ActionID{SpellID: 26986}
	baseCost := 0.16 * druid.BaseMana
	minBaseDamage := 1038.0
	maxBaseDamage := 1222.0
	spellCoefficient := 1.0 * (1 + 0.04*float64(druid.Talents.WrathOfCenarius))
	manaMetrics := druid.NewManaMetrics(actionID)

	// This seems to be unaffected by wrath of cenarius so it needs to come first.
	bonusFlatDamage := core.TernaryFloat64(druid.Equip[items.ItemSlotRanged].ID == IvoryMoongoddess, 55*spellCoefficient, 0)
	bonusFlatDamage += core.TernaryFloat64(druid.Equip[items.ItemSlotRanged].ID == ShootingStar, 165*spellCoefficient, 0)

	effect := core.SpellEffect{
		ProcMask:         core.ProcMaskSpellDamage,
		DamageMultiplier: 1 + 0.02*float64(druid.Talents.Moonfury),
		ThreatMultiplier: 1,
		BaseDamage:       core.BaseDamageConfigMagic(minBaseDamage+bonusFlatDamage, maxBaseDamage+bonusFlatDamage, spellCoefficient),
		OutcomeApplier:   druid.OutcomeFuncMagicHitAndCrit(druid.SpellCritMultiplier(1, 0.2*float64(druid.Talents.Vengeance))),
		/*OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if spellEffect.Landed() && druid.HasGlyph(proto.DruidMajorGlyph_GlyphOfStarfire) && druid.MoonfireDot.IsActive() {
				// Add 3seconds to Moonfire Tick up to +9s
				druid.MoonfireDot.NumberOfTicks += 1
			}
		}, */
		OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if spellEffect.Outcome.Matches(core.OutcomeCrit) {
				hasMoonkinForm := core.TernaryFloat64(druid.Talents.MoonkinForm, 1, 0)
				druid.AddMana(sim, druid.MaxMana()*0.02*hasMoonkinForm, manaMetrics, true)
			}
		},
	}

	if druid.HasSetBonus(ItemSetNordrassilRegalia, 4) {
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
	// T6-4P
	if druid.HasSetBonus(ItemSetThunderheartRegalia, 4) {
		effect.BonusSpellCritRating += 5 * core.CritRatingPerCritChance
	}
	// T7-4P
	if druid.DruidTier.balance_t7_4 {
		effect.BonusSpellCritRating += 5 * core.CritRatingPerCritChance
	}
	// Improved Faerie Fire
	if druid.CurrentTarget.HasAura("Improved Faerie Fire") {
		effect.BonusSpellCritRating += float64(druid.Talents.ImprovedFaerieFire) * 1 * core.CritRatingPerCritChance
	}
	// Improved Insect Swarm
	if druid.CurrentTarget.HasAura("Moonfire") {
		effect.BonusSpellCritRating += core.CritRatingPerCritChance * float64(druid.Talents.ImprovedInsectSwarm)
	}
	// Lunar eclipse buff
	if druid.HasAura("Lunar Eclipse proc") {
		effect.BonusSpellCritRating += core.CritRatingPerCritChance * 40
		// T8-2P
		if druid.DruidTier.balance_t8_2 {
			effect.BonusSpellCritRating += core.CritRatingPerCritChance * 7
		}
	}
	// Nature's Majesty
	effect.BonusSpellCritRating += 2 * float64(druid.Talents.NaturesMajesty) * core.CritRatingPerCritChance

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

			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				druid.applyNaturesSwiftness(cast)
				druid.ApplyClearcasting(sim, spell, cast)
				if druid.HasActiveAura("Elune's Wrath") {
					cast.CastTime = 0
					druid.GetAura("Elune's Wrath").Deactivate(sim)
				}
			},
		},
		ApplyEffects: core.ApplyEffectFuncDirectDamage(effect),
	})
}
