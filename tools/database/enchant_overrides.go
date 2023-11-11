package database

import (
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

// Note: EffectId AND SpellId are required for all enchants, because they are
// used by various importers/exporters. ItemId is optional.

var EnchantOverrides = []*proto.UIEnchant{
	// Multi-slot
	// {EffectId: 2988, ItemId: 29487, SpellId: 35419, Name: "Nature Armor Kit", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.NatureResistance: 8}.ToFloatArray(), Type: proto.ItemType_ItemTypeChest, ExtraTypes: []proto.ItemType{proto.ItemType_ItemTypeLegs, proto.ItemType_ItemTypeHands, proto.ItemType_ItemTypeFeet}, EnchantType: proto.EnchantType_EnchantTypeKit},
	{EffectId: 2544, ItemId: 18330, SpellId: 22844, Name: "Arcanum of Focus", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.SpellHealing: 8}.ToFloatArray(), Type: proto.ItemType_ItemTypeHead, ExtraTypes: []proto.ItemType{proto.ItemType_ItemTypeLegs}, EnchantType: proto.EnchantType_EnchantTypeKit},
	{EffectId: 2543, ItemId: 18329, SpellId: 22840, Name: "Arcanum of Rapidity", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.MeleeHaste: 1}.ToFloatArray(), Type: proto.ItemType_ItemTypeHead, ExtraTypes: []proto.ItemType{proto.ItemType_ItemTypeLegs}, EnchantType: proto.EnchantType_EnchantTypeKit},
	{EffectId: 1506, ItemId: 11645, SpellId: 15397, Name: "Arcanum of Voracity", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.Strength: 8}.ToFloatArray(), Type: proto.ItemType_ItemTypeHead, ExtraTypes: []proto.ItemType{proto.ItemType_ItemTypeLegs}, EnchantType: proto.EnchantType_EnchantTypeKit},
	{EffectId: 1483, ItemId: 11622, SpellId: 15340, Name: "Arcanum of Voracity", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.Mana: 150}.ToFloatArray(), Type: proto.ItemType_ItemTypeHead, ExtraTypes: []proto.ItemType{proto.ItemType_ItemTypeLegs}, EnchantType: proto.EnchantType_EnchantTypeKit},

	// Head
	// {EffectId: 3795, ItemId: 44069, SpellId: 59777, Name: "Arcanum of Triumph", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.AttackPower: 50, stats.RangedAttackPower: 50, stats.Resilience: 20}.ToFloatArray(), Type: proto.ItemType_ItemTypeHead},

	// Shoulder
	// {EffectId: 2998, ItemId: 29187, SpellId: 35441, Name: "Inscription of Endurance", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.ArcaneResistance: 7, stats.FireResistance: 7, stats.FrostResistance: 7, stats.NatureResistance: 7, stats.ShadowResistance: 7}.ToFloatArray(), Type: proto.ItemType_ItemTypeShoulder},

	// Back
	// {EffectId: 1262, ItemId: 37330, SpellId: 44596, Name: "Superior Arcane Resistance", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.ArcaneResistance: 20}.ToFloatArray(), Type: proto.ItemType_ItemTypeBack},

	// Chest
	// {EffectId: 3245, ItemId: 37340, SpellId: 44588, Name: "Exceptional Resilience", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.Resilience: 20}.ToFloatArray(), Type: proto.ItemType_ItemTypeChest},

	// Wrist
	// {EffectId: 3845, ItemId: 44484, SpellId: 44575, Name: "Greater Assault", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.AttackPower: 50, stats.RangedAttackPower: 50}.ToFloatArray(), Type: proto.ItemType_ItemTypeWrist},

	// Hands
	// {EffectId: 3253, ItemId: 44485, SpellId: 44625, Name: "Armsman", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.Parry: 10}.ToFloatArray(), Type: proto.ItemType_ItemTypeHands},

	// Feet
	// {EffectId: 3606, SpellId: 55016, Name: "Nitro Boosts", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.MeleeCrit: 24, stats.SpellCrit: 24}.ToFloatArray(), Type: proto.ItemType_ItemTypeFeet, RequiredProfession: proto.Profession_Engineering},

	// Weapon
	// {EffectId: 3870, ItemId: 46348, SpellId: 64579, Name: "Blood Draining", Quality: proto.ItemQuality_ItemQualityEpic, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon},

	// Ranged
	// {EffectId: 3607, ItemId: 41146, SpellId: 55076, Name: "Sun Scope", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeRanged},
}
