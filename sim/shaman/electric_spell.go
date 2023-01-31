package shaman

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

// Totem Item IDs
const (
	StormfuryTotem           = 31031
	TotemOfAncestralGuidance = 32330
	TotemOfStorms            = 23199
	TotemOfTheVoid           = 28248
	TotemOfHex               = 40267
	VentureCoLightningRod    = 38361
	ThunderfallTotem         = 45255
)

const (
	// This could be value or bitflag if we ended up needing multiple flags at the same time.
	//1 to 5 are used by MaelstromWeapon Stacks
	CastTagLightningOverload int32 = 6
)

// Shared precomputation logic for LB and CL.
func (shaman *Shaman) newElectricSpellConfig(actionID core.ActionID, baseCost float64, baseCastTime time.Duration, isLightningOverload bool) core.SpellConfig {
	spell := core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolNature,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       SpellFlagElectric | SpellFlagFocusable,

		ManaCost: core.ManaCostOptions{
			BaseCost:   core.TernaryFloat64(isLightningOverload, 0, baseCost),
			Multiplier: 1 - 0.02*float64(shaman.Talents.Convection),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				CastTime: baseCastTime - time.Millisecond*100*time.Duration(shaman.Talents.LightningMastery),
				GCD:      core.GCDDefault,
			},
		},

		BonusHitRating: float64(shaman.Talents.ElementalPrecision) * core.SpellHitRatingPerHitChance,
		BonusCritRating: 0 +
			float64(shaman.Talents.TidalMastery)*core.CritRatingPerCritChance +
			core.TernaryFloat64(shaman.Talents.CallOfThunder, 5*core.CritRatingPerCritChance, 0),
		DamageMultiplier: 1 + 0.01*float64(shaman.Talents.Concussion),
		CritMultiplier:   shaman.ElementalCritMultiplier(0),
		ThreatMultiplier: shaman.spellThreatMultiplier(),
	}

	if isLightningOverload {
		spell.ActionID.Tag = CastTagLightningOverload
		spell.Cast.DefaultCast.CastTime = 0
		spell.Cast.DefaultCast.GCD = 0
		spell.Cast.DefaultCast.Cost = 0
		spell.DamageMultiplier *= 0.5
		spell.ThreatMultiplier = 0
	}

	return spell
}

func (shaman *Shaman) electricSpellBonusDamage(spellCoeff float64) float64 {
	bonusDamage := 0 +
		core.TernaryFloat64(shaman.Equip[core.ItemSlotRanged].ID == TotemOfStorms, 33, 0) +
		core.TernaryFloat64(shaman.Equip[core.ItemSlotRanged].ID == TotemOfTheVoid, 55, 0) +
		core.TernaryFloat64(shaman.Equip[core.ItemSlotRanged].ID == TotemOfAncestralGuidance, 85, 0) +
		core.TernaryFloat64(shaman.Equip[core.ItemSlotRanged].ID == TotemOfHex, 165, 0)

	return bonusDamage * spellCoeff // These items do not benefit from the bonus coeff from shamanism.
}
