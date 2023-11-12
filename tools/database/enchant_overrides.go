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
	{EffectId: 2544, ItemId: 18330, SpellId: 22844, Name: "Arcanum of Focus", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.Healing: 8}.ToFloatArray(), Type: proto.ItemType_ItemTypeHead, ExtraTypes: []proto.ItemType{proto.ItemType_ItemTypeLegs}, EnchantType: proto.EnchantType_EnchantTypeKit},
	{EffectId: 2543, ItemId: 18329, SpellId: 22840, Name: "Arcanum of Rapidity", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.MeleeHaste: 1}.ToFloatArray(), Type: proto.ItemType_ItemTypeHead, ExtraTypes: []proto.ItemType{proto.ItemType_ItemTypeLegs}, EnchantType: proto.EnchantType_EnchantTypeKit},
	{EffectId: 1506, ItemId: 11645, SpellId: 15397, Name: "Arcanum of Voracity", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.Strength: 8}.ToFloatArray(), Type: proto.ItemType_ItemTypeHead, ExtraTypes: []proto.ItemType{proto.ItemType_ItemTypeLegs}, EnchantType: proto.EnchantType_EnchantTypeKit},
	{EffectId: 1483, ItemId: 11622, SpellId: 15340, Name: "Arcanum of Rumination", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.Mana: 150}.ToFloatArray(), Type: proto.ItemType_ItemTypeHead, ExtraTypes: []proto.ItemType{proto.ItemType_ItemTypeLegs}, EnchantType: proto.EnchantType_EnchantTypeKit},

	{EffectId: 2590, ItemId: 19789, SpellId: 24167, Name: "Prophetic Aura", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.MP5: 4, stats.Stamina: 10, stats.Healing: 24}.ToFloatArray(), Type: proto.ItemType_ItemTypeHead, ExtraTypes: []proto.ItemType{proto.ItemType_ItemTypeLegs}, EnchantType: proto.EnchantType_EnchantTypeKit},

	// Head
	// {EffectId: 3795, ItemId: 44069, SpellId: 59777, Name: "Arcanum of Triumph", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.AttackPower: 50, stats.RangedAttackPower: 50, stats.Resilience: 20}.ToFloatArray(), Type: proto.ItemType_ItemTypeHead},

	// Shoulder
	{EffectId: 2604, ItemId: 20078, SpellId: 24420, Name: "Zandalar Signet of Serenity", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.Healing: 33}.ToFloatArray(), Type: proto.ItemType_ItemTypeShoulder},
	{EffectId: 2605, ItemId: 20076, SpellId: 24421, Name: "Zandalar Signet of Mojo", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.SpellPower: 18, stats.Healing: 18}.ToFloatArray(), Type: proto.ItemType_ItemTypeShoulder},
	{EffectId: 2606, ItemId: 20077, SpellId: 24422, Name: "Zandalar Signet of Might", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.AttackPower: 30}.ToFloatArray(), Type: proto.ItemType_ItemTypeShoulder},
	{EffectId: 2715, ItemId: 23547, SpellId: 29475, Name: "Resilience of the Scourge", Quality: proto.ItemQuality_ItemQualityEpic, Stats: stats.Stats{stats.Healing: 31, stats.MP5: 5}.ToFloatArray(), Type: proto.ItemType_ItemTypeShoulder},
	{EffectId: 2717, ItemId: 23548, SpellId: 29483, Name: "Might of the Scourge", Quality: proto.ItemQuality_ItemQualityEpic, Stats: stats.Stats{stats.AttackPower: 26, stats.MeleeCrit: 0.01}.ToFloatArray(), Type: proto.ItemType_ItemTypeShoulder},
	{EffectId: 2716, ItemId: 23549, SpellId: 29480, Name: "Fortitude of the Scourge", Quality: proto.ItemQuality_ItemQualityEpic, Stats: stats.Stats{stats.Stamina: 16, stats.Armor: 100}.ToFloatArray(), Type: proto.ItemType_ItemTypeShoulder},
	{EffectId: 2721, ItemId: 23545, SpellId: 29467, Name: "Power of the Scourge", Quality: proto.ItemQuality_ItemQualityEpic, Stats: stats.Stats{stats.SpellPower: 15, stats.Healing: 15, stats.SpellCrit: 0.01}.ToFloatArray(), Type: proto.ItemType_ItemTypeShoulder},

	// Back TODO: Classic add non-dps enchants
	// {EffectId: 1262, ItemId: 37330, SpellId: 44596, Name: "Superior Arcane Resistance", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.ArcaneResistance: 20}.ToFloatArray(), Type: proto.ItemType_ItemTypeBack},

	// Chest TODO: Classic add non-dps enchants
	{EffectId: 1891, SpellId: 20025, Name: "Greater Stats", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Stamina: 4, stats.Agility: 4, stats.Strength: 4, stats.Intellect: 4, stats.Spirit: 4}.ToFloatArray(), Type: proto.ItemType_ItemTypeChest},
	{EffectId: 928, SpellId: 13941, Name: "Stats", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Stamina: 3, stats.Agility: 3, stats.Strength: 3, stats.Intellect: 3, stats.Spirit: 3}.ToFloatArray(), Type: proto.ItemType_ItemTypeChest},
	{EffectId: 866, SpellId: 13700, Name: "Lesser Stats", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Stamina: 2, stats.Agility: 2, stats.Strength: 2, stats.Intellect: 2, stats.Spirit: 2}.ToFloatArray(), Type: proto.ItemType_ItemTypeChest},
	{EffectId: 847, SpellId: 13626, Name: "Lesser Stats", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Stamina: 1, stats.Agility: 1, stats.Strength: 1, stats.Intellect: 1, stats.Spirit: 1}.ToFloatArray(), Type: proto.ItemType_ItemTypeChest},
	{EffectId: 913, SpellId: 13917, Name: "Superior Mana", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Mana: 65}.ToFloatArray(), Type: proto.ItemType_ItemTypeChest},
	{EffectId: 857, SpellId: 13663, Name: "Greater Mana", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Mana: 50}.ToFloatArray(), Type: proto.ItemType_ItemTypeChest},
	{EffectId: 24, SpellId: 7443, Name: "Minor Mana", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Mana: 5}.ToFloatArray(), Type: proto.ItemType_ItemTypeChest},

	// Wrist TODO: Classic add non-dps enchants
	{EffectId: 1883, SpellId: 20008, Name: "Greater Intellect", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Intellect: 7}.ToFloatArray(), Type: proto.ItemType_ItemTypeWrist},
	{EffectId: 905, SpellId: 13822, Name: "Intellect", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Intellect: 5}.ToFloatArray(), Type: proto.ItemType_ItemTypeWrist},
	{EffectId: 723, SpellId: 13622, Name: "Lesser Intellect", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Intellect: 3}.ToFloatArray(), Type: proto.ItemType_ItemTypeWrist},
	{EffectId: 907, SpellId: 13846, Name: "Greater Spirit", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Spirit: 7}.ToFloatArray(), Type: proto.ItemType_ItemTypeWrist},
	{EffectId: 927, SpellId: 13939, Name: "Greater Strength", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Strength: 7}.ToFloatArray(), Type: proto.ItemType_ItemTypeWrist},
	{EffectId: 2566, SpellId: 23802, Name: "Healing Power", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Healing: 24}.ToFloatArray(), Type: proto.ItemType_ItemTypeWrist},
	{EffectId: 255, SpellId: 7859, Name: "Lesser Spirit", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Spirit: 3}.ToFloatArray(), Type: proto.ItemType_ItemTypeWrist},
	{EffectId: 823, SpellId: 13536, Name: "Lesser Strength", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Strength: 3}.ToFloatArray(), Type: proto.ItemType_ItemTypeWrist},
	{EffectId: 243, SpellId: 7766, Name: "Minor Spirit", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Spirit: 1}.ToFloatArray(), Type: proto.ItemType_ItemTypeWrist},
	{EffectId: 248, SpellId: 7782, Name: "Minor Strength", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Strength: 1}.ToFloatArray(), Type: proto.ItemType_ItemTypeWrist},
	{EffectId: 851, SpellId: 13642, Name: "Spirit", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Spirit: 5}.ToFloatArray(), Type: proto.ItemType_ItemTypeWrist},
	{EffectId: 856, SpellId: 13661, Name: "Strength", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Strength: 5}.ToFloatArray(), Type: proto.ItemType_ItemTypeWrist},
	{EffectId: 1884, SpellId: 20009, Name: "Superior Spirit", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Spirit: 9}.ToFloatArray(), Type: proto.ItemType_ItemTypeWrist},
	{EffectId: 1885, SpellId: 20010, Name: "Superior Strength", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Strength: 9}.ToFloatArray(), Type: proto.ItemType_ItemTypeWrist},

	// Hands TODO: Classic add non-dps enchants
	// {EffectId: 2614, SpellId: 25073, Name: "Shadow Power", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.SpellPower: 20}.ToFloatArray(), Type: proto.ItemType_ItemTypeHands},
	{EffectId: 904, SpellId: 13815, Name: "Agility", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Agility: 5}.ToFloatArray(), Type: proto.ItemType_ItemTypeHands},
	{EffectId: 856, SpellId: 13887, Name: "Strength", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Strength: 5}.ToFloatArray(), Type: proto.ItemType_ItemTypeHands},
	{EffectId: 1887, SpellId: 20012, Name: "Greater Agility", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Agility: 7}.ToFloatArray(), Type: proto.ItemType_ItemTypeHands},
	{EffectId: 927, SpellId: 20013, Name: "Greater Strength", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Strength: 7}.ToFloatArray(), Type: proto.ItemType_ItemTypeHands},
	{EffectId: 2617, SpellId: 25079, Name: "Healing Power", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Healing: 30}.ToFloatArray(), Type: proto.ItemType_ItemTypeHands},
	// {EffectId: 2616, SpellId: 25078, Name: "Fire Power", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.SpellPower: 20}.ToFloatArray(), Type: proto.ItemType_ItemTypeHands},

	// Feet
	{EffectId: 1887, SpellId: 20023, Name: "Greater Agility", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Agility: 7}.ToFloatArray(), Type: proto.ItemType_ItemTypeFeet},
	{EffectId: 904, SpellId: 13935, Name: "Agility", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Agility: 5}.ToFloatArray(), Type: proto.ItemType_ItemTypeFeet},
	{EffectId: 849, SpellId: 13637, Name: "Lesser Agility", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Agility: 3}.ToFloatArray(), Type: proto.ItemType_ItemTypeFeet},
	{EffectId: 247, SpellId: 7867, Name: "Minor Agility", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Agility: 1}.ToFloatArray(), Type: proto.ItemType_ItemTypeFeet},
	{EffectId: 255, SpellId: 13687, Name: "Lesser Spirit", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Spirit: 3}.ToFloatArray(), Type: proto.ItemType_ItemTypeFeet},
	{EffectId: 851, SpellId: 20024, Name: "Spirit", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Spirit: 5}.ToFloatArray(), Type: proto.ItemType_ItemTypeFeet},

	// Weapon
	{EffectId: 2568, SpellId: 23804, Name: "Mighty Intellect", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Intellect: 22}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon},
	{EffectId: 1904, SpellId: 20036, Name: "Major Intellect", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Intellect: 9}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon},
	{EffectId: 723, SpellId: 7793, Name: "Lesser Intellect", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Intellect: 3}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon},
	{EffectId: 2646, SpellId: 27837, Name: "2H Agility", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Agility: 25}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon},
	{EffectId: 2564, SpellId: 23800, Name: "Agility", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Agility: 25}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon},

	// Ranged TODO: Classic scopes for ranged hit and + damage
	// {EffectId: 2523, ItemId: 18283, SpellId: 22779, Name: "Biznicks 247x128 Accurascope", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.MeleeHit: 0.03}.ToFloatArray(), Type: proto.ItemType_ItemTypeRanged},
}
