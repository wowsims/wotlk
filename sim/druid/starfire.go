package druid

import (
	"github.com/wowsims/wotlk/sim/core/proto"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/items"
	"github.com/wowsims/wotlk/sim/core/stats"
)

// Idol IDs
const IvoryMoongoddess int32 = 27518
const ShootingStar int32 = 40321

func (druid *Druid) registerStarfireSpell() {

	actionID := core.ActionID{SpellID: 48465}
	baseCost := 0.16 * druid.BaseMana
	minBaseDamage := 1038.0
	maxBaseDamage := 1222.0
	spellCoefficient := 1.0 * (1 + 0.04*float64(druid.Talents.WrathOfCenarius))
	manaMetrics := druid.NewManaMetrics(core.ActionID{SpellID: 24858})

	// This seems to be unaffected by wrath of cenarius so it needs to come first.
	bonusFlatDamage := core.TernaryFloat64(druid.Equip[items.ItemSlotRanged].ID == IvoryMoongoddess, 55*spellCoefficient, 0)
	bonusFlatDamage += core.TernaryFloat64(druid.Equip[items.ItemSlotRanged].ID == ShootingStar, 165*spellCoefficient, 0)

	effect := core.SpellEffect{
		ProcMask:         core.ProcMaskSpellDamage,
		DamageMultiplier: 1 + druid.TalentsBonuses.moonfuryMultiplier,
		ThreatMultiplier: 1,
		BaseDamage:       core.BaseDamageConfigMagic(minBaseDamage+bonusFlatDamage, maxBaseDamage+bonusFlatDamage, spellCoefficient),
		OutcomeApplier:   druid.OutcomeFuncMagicHitAndCrit(druid.SpellCritMultiplier(1, druid.TalentsBonuses.vengeanceModifier)),
		OnInit: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			spellEffect.BonusSpellCritRating = 0
			// Improved Insect Swarm
			if druid.MoonfireDot.IsActive() {
				spellEffect.BonusSpellCritRating += core.CritRatingPerCritChance * float64(druid.Talents.ImprovedInsectSwarm)
			}
			// T6-4P
			if druid.SetBonuses.balance_t6_2 {
				spellEffect.BonusSpellCritRating += 5 * core.CritRatingPerCritChance
			}
			// T7-4P
			if druid.SetBonuses.balance_t7_4 {
				spellEffect.BonusSpellCritRating += 5 * core.CritRatingPerCritChance
			}
			// Improved Faerie Fire
			if druid.CurrentTarget.HasAura("Improved Faerie Fire") {
				spellEffect.BonusSpellCritRating += druid.TalentsBonuses.iffBonusCrit
			}
			// Lunar eclipse buff
			if druid.HasAura("Lunar Eclipse proc") {
				spellEffect.BonusSpellCritRating += core.CritRatingPerCritChance * 40
				// T8-2P
				if druid.SetBonuses.balance_t8_2 {
					spellEffect.BonusSpellCritRating += core.CritRatingPerCritChance * 7
				}
			}
			// T9-4P
			if druid.SetBonuses.balance_t9_4 {
				spellEffect.DamageMultiplier *= 1.04
			}
			// Nature's Majesty
			spellEffect.BonusSpellCritRating += druid.TalentsBonuses.naturesMajestyBonusCrit
		},
		OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if spellEffect.Landed() {
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
				if druid.HasMajorGlyph(proto.DruidMajorGlyph_GlyphOfStarfire) && druid.MoonfireDot.IsActive() {
					maxMoonfireTicks := druid.maxMoonfireTicks()
					if druid.MoonfireDot.NumberOfTicks < maxMoonfireTicks {
						druid.MoonfireDot.NumberOfTicks += 1
						druid.MoonfireDot.UpdateExpires(druid.MoonfireDot.ExpiresAt() + time.Second*3)
					}
				}
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

	druid.Starfire = druid.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolArcane,

		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost:     baseCost * druid.TalentsBonuses.moonglowMultiplier,
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond*3500 - druid.TalentsBonuses.starlightWrathModifier,
			},

			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				druid.applyNaturesSwiftness(cast)
				druid.ApplyClearcasting(sim, spell, cast)
				druid.ApplySwiftStarfireBonus(sim, cast)
				if druid.HasActiveAura("Elune's Wrath") {
					cast.CastTime = 0
					druid.GetAura("Elune's Wrath").Deactivate(sim)
				}
			},
		},
		ApplyEffects: core.ApplyEffectFuncDirectDamage(effect),
	})
}
