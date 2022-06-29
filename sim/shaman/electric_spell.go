package shaman

import (
	"time"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/items"
	"github.com/wowsims/tbc/sim/core/stats"
)

// Totem Item IDs
const (
	StormfuryTotem           = 31031
	TotemOfAncestralGuidance = 32330
	TotemOfImpact            = 27947
	TotemOfStorms            = 23199
	TotemOfThePulsingEarth   = 29389
	TotemOfTheVoid           = 28248
	TotemOfRage              = 22395
)

const (
	CastTagLightningOverload int32 = 1 // This could be value or bitflag if we ended up needing multiple flags at the same time.
)

// Mana cost numbers based on in-game testing:
//
// With 5/5 convection:
// Normal: 270, w/ EF: 150
//
// With 5/5 convection and TotPE equipped:
// Normal: 246, w/ EF: 136

// Shared precomputation logic for LB and CL.
func (shaman *Shaman) newElectricSpellConfig(actionID core.ActionID, baseCost float64, baseCastTime time.Duration, isLightningOverload bool) core.SpellConfig {
	spell := core.SpellConfig{
		ActionID:     actionID,
		SpellSchool:  core.SpellSchoolNature,
		Flags:        SpellFlagElectric,
		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost:     baseCost,
				CastTime: baseCastTime,
				GCD:      core.GCDDefault,
			},
		},
	}

	if isLightningOverload {
		spell.ActionID.Tag = CastTagLightningOverload
		spell.ResourceType = 0
		spell.Cast.DefaultCast.CastTime = 0
		spell.Cast.DefaultCast.GCD = 0
		spell.Cast.DefaultCast.Cost = 0
	} else if shaman.Talents.LightningMastery > 0 {
		// Convection applies against the base cost of the spell.
		spell.Cast.DefaultCast.Cost -= baseCost * float64(shaman.Talents.Convection) * 0.02
		spell.Cast.DefaultCast.CastTime -= time.Millisecond * 100 * time.Duration(shaman.Talents.LightningMastery)
	}

	return spell
}

// Helper for precomputing spell effects.
func (shaman *Shaman) newElectricSpellEffect(minBaseDamage float64, maxBaseDamage float64, spellCoefficient float64, isLightningOverload bool) core.SpellEffect {
	effect := core.SpellEffect{
		ProcMask:            core.ProcMaskSpellDamage,
		BonusSpellHitRating: float64(shaman.Talents.ElementalPrecision) * 2 * core.SpellHitRatingPerHitChance,
		BonusSpellCritRating: 0 +
			(float64(shaman.Talents.TidalMastery) * 1 * core.SpellCritRatingPerCritChance) +
			(float64(shaman.Talents.CallOfThunder) * 1 * core.SpellCritRatingPerCritChance),
		BonusSpellPower: 0 +
			core.TernaryFloat64(shaman.Equip[items.ItemSlotRanged].ID == TotemOfStorms, 33, 0) +
			core.TernaryFloat64(shaman.Equip[items.ItemSlotRanged].ID == TotemOfTheVoid, 55, 0) +
			core.TernaryFloat64(shaman.Equip[items.ItemSlotRanged].ID == TotemOfAncestralGuidance, 85, 0),
		DamageMultiplier: 1 * (1 + 0.01*float64(shaman.Talents.Concussion)),
		ThreatMultiplier: 1 - (0.1/3)*float64(shaman.Talents.ElementalPrecision),
		BaseDamage:       core.BaseDamageConfigMagic(minBaseDamage, maxBaseDamage, spellCoefficient),
		OutcomeApplier:   shaman.OutcomeFuncMagicHitAndCrit(shaman.ElementalCritMultiplier()),
	}

	if isLightningOverload {
		effect.DamageMultiplier *= 0.5
		effect.ThreatMultiplier = 0
	}

	return effect
}

// Shared LB/CL logic that is dynamic, i.e. can't be precomputed.
func (shaman *Shaman) applyElectricSpellCastInitModifiers(spell *core.Spell, cast *core.Cast) {
	shaman.modifyCastClearcasting(spell, cast)
	if shaman.ElementalMasteryAura != nil && shaman.ElementalMasteryAura.IsActive() {
		cast.Cost = 0
	}
}
