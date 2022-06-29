package items

import (
	"github.com/wowsims/tbc/sim/core/proto"
	"github.com/wowsims/tbc/sim/core/stats"
)

var Enchants = []Enchant{
	// Head
	{ID: 29186, EffectID: 2999, Name: "Glyph of the Defender", Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{stats.Defense: 16, stats.Dodge: 17}, ItemType: proto.ItemType_ItemTypeHead},
	{ID: 29191, EffectID: 3002, Name: "Glyph of Power", Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{stats.SpellPower: 22, stats.SpellHit: 14}, ItemType: proto.ItemType_ItemTypeHead},
	{ID: 29192, EffectID: 3003, Name: "Glyph of Ferocity", Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{stats.AttackPower: 34, stats.RangedAttackPower: 34, stats.MeleeHit: 16}, ItemType: proto.ItemType_ItemTypeHead},
	{ID: 30846, EffectID: 3096, Name: "Glyph of the Outcast", Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{stats.Strength: 17, stats.Intellect: 16}, ItemType: proto.ItemType_ItemTypeHead},
	{ID: 29193, EffectID: 3004, Name: "Glyph of the Gladiator", Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{stats.Stamina: 18, stats.Resilience: 20}, ItemType: proto.ItemType_ItemTypeHead},

	// ZG Head Enchants
	{ID: 19782, EffectID: 2583, Name: "Presence of Might", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{stats.Stamina: 10, stats.Defense: 10, stats.BlockValue: 15}, ItemType: proto.ItemType_ItemTypeHead, ClassAllowlist: []proto.Class{proto.Class_ClassWarrior}},

	// Shoulder
	{ID: 28886, EffectID: 2982, Name: "Greater Inscription of Discipline", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{stats.SpellPower: 18, stats.SpellCrit: 10}, ItemType: proto.ItemType_ItemTypeShoulder},
	{ID: 28888, EffectID: 2986, Name: "Greater Inscription of Vengeance", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{stats.AttackPower: 30, stats.RangedAttackPower: 30, stats.MeleeCrit: 10}, ItemType: proto.ItemType_ItemTypeShoulder},
	{ID: 28889, EffectID: 2978, Name: "Greater Inscription of Warding", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{stats.Defense: 10, stats.Dodge: 15}, ItemType: proto.ItemType_ItemTypeShoulder},
	{ID: 28909, EffectID: 2995, Name: "Greater Inscription of the Orb", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{stats.SpellPower: 12, stats.SpellCrit: 15}, ItemType: proto.ItemType_ItemTypeShoulder},
	{ID: 28910, EffectID: 2997, Name: "Greater Inscription of the Blade", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{stats.AttackPower: 20, stats.RangedAttackPower: 20, stats.MeleeCrit: 15}, ItemType: proto.ItemType_ItemTypeShoulder},
	{ID: 28911, EffectID: 2991, Name: "Greater Inscription of the Knight", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{stats.Defense: 15, stats.Dodge: 10}, ItemType: proto.ItemType_ItemTypeShoulder},
	{ID: 20076, EffectID: 2605, Name: "Zandalar Signet of Mojo", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{stats.SpellPower: 18}, ItemType: proto.ItemType_ItemTypeShoulder},
	{ID: 23545, EffectID: 2721, Name: "Power of the Scourge", Quality: proto.ItemQuality_ItemQualityEpic, Bonus: stats.Stats{stats.SpellPower: 15, stats.SpellCrit: 14}, ItemType: proto.ItemType_ItemTypeShoulder},
	{ID: 23548, EffectID: 2717, Name: "Might of the Scourge", Quality: proto.ItemQuality_ItemQualityEpic, Bonus: stats.Stats{stats.AttackPower: 26, stats.RangedAttackPower: 26, stats.MeleeCrit: 14}, ItemType: proto.ItemType_ItemTypeShoulder},
	{ID: 23549, EffectID: 2716, Name: "Fortitude of the Scourge", Quality: proto.ItemQuality_ItemQualityEpic, Bonus: stats.Stats{stats.Stamina: 16, stats.Armor: 100}, ItemType: proto.ItemType_ItemTypeShoulder},

	// Back
	{ID: 33148, EffectID: 2622, Name: "Enchant Cloak - Dodge", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{stats.Dodge: 12}, ItemType: proto.ItemType_ItemTypeBack},
	{ID: 33150, EffectID: 2621, Name: "Enchant Cloak - Subtlety", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{}, ItemType: proto.ItemType_ItemTypeBack},
	{ID: 11206, EffectID: 849, Name: "Enchant Cloak - Lesser Agility", Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{stats.Agility: 3}, ItemType: proto.ItemType_ItemTypeBack},
	{ID: 34004, EffectID: 368, Name: "Enchant Cloak - Greater Agility", IsSpellID: true, Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.Agility: 12}, ItemType: proto.ItemType_ItemTypeBack},
	{ID: 28274, EffectID: 2938, Name: "Enchant Cloak - Spell Penetration", Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{stats.SpellPenetration: 20}, ItemType: proto.ItemType_ItemTypeBack},
	{ID: 28277, EffectID: 1441, Name: "Enchant Cloak - Greater Shadow Resistance", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{stats.ShadowResistance: 15}, ItemType: proto.ItemType_ItemTypeBack},
	{ID: 35756, EffectID: 2648, Name: "Enchant Cloak - Steelweave", Phase: 5, Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{stats.Defense: 12}, ItemType: proto.ItemType_ItemTypeBack},

	// Chest
	{ID: 27957, EffectID: 2659, Name: "Chest - Exceptional Health", IsSpellID: true, Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.Health: 150}, ItemType: proto.ItemType_ItemTypeChest},
	{ID: 22546, EffectID: 2660, Name: "Chest - Exceptional Mana", Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.Mana: 150}, ItemType: proto.ItemType_ItemTypeChest},
	{ID: 24003, EffectID: 2661, Name: "Chest - Exceptional Stats", Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.Stamina: 6, stats.Intellect: 6, stats.Spirit: 6, stats.Strength: 6, stats.Agility: 6}, ItemType: proto.ItemType_ItemTypeChest},
	{ID: 28270, EffectID: 2933, Name: "Chest - Major Resilience", Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{stats.Resilience: 15}, ItemType: proto.ItemType_ItemTypeChest},
	{ID: 33990, EffectID: 1144, Name: "Chest - Major Spirit", IsSpellID: true, Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.Spirit: 15}, ItemType: proto.ItemType_ItemTypeChest},
	{ID: 33991, EffectID: 3150, Name: "Chest - Restore Mana Prime", IsSpellID: true, Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.MP5: 6}, ItemType: proto.ItemType_ItemTypeChest},
	{ID: 35500, EffectID: 1950, Name: "Chest - Defense", Phase: 5, Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.Defense: 15}, ItemType: proto.ItemType_ItemTypeChest},

	// Wrist
	{ID: 22533, EffectID: 2649, Name: "Bracer - Fortitude", Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{stats.Stamina: 12}, ItemType: proto.ItemType_ItemTypeWrist},
	{ID: 22534, EffectID: 2650, Name: "Bracer - Spellpower", Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{stats.SpellPower: 15}, ItemType: proto.ItemType_ItemTypeWrist},
	{ID: 34001, EffectID: 369, Name: "Bracer - Major Intellect", IsSpellID: true, Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{stats.Intellect: 12}, ItemType: proto.ItemType_ItemTypeWrist},
	{ID: 27899, EffectID: 2647, Name: "Bracer - Brawn", IsSpellID: true, Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{stats.Strength: 12}, ItemType: proto.ItemType_ItemTypeWrist},
	{ID: 34002, EffectID: 1593, Name: "Bracer - Assault", IsSpellID: true, Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{stats.AttackPower: 24, stats.RangedAttackPower: 24}, ItemType: proto.ItemType_ItemTypeWrist},
	{ID: 27905, EffectID: 1891, Name: "Bracer - Stats", IsSpellID: true, Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{stats.Stamina: 4, stats.Intellect: 4, stats.Spirit: 4, stats.Strength: 4, stats.Agility: 4}, ItemType: proto.ItemType_ItemTypeWrist},

	// Hands
	{ID: 20727, EffectID: 2614, Name: "Gloves - Shadow Power", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{stats.ShadowSpellPower: 20}, ItemType: proto.ItemType_ItemTypeHands},
	{ID: 28271, EffectID: 2935, Name: "Gloves - Spell Strike", Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{stats.SpellHit: 15}, ItemType: proto.ItemType_ItemTypeHands},
	{ID: 28272, EffectID: 2937, Name: "Gloves - Major Spellpower", Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{stats.SpellPower: 20}, ItemType: proto.ItemType_ItemTypeHands},
	{ID: 33995, EffectID: 684, Name: "Gloves - Major Strength", IsSpellID: true, Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{stats.Strength: 15}, ItemType: proto.ItemType_ItemTypeHands},
	{ID: 33152, EffectID: 2564, Name: "Gloves - Major Agility", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{stats.Agility: 15}, ItemType: proto.ItemType_ItemTypeHands},
	{ID: 33153, EffectID: 2613, Name: "Gloves - Threat", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{}, ItemType: proto.ItemType_ItemTypeHands},

	// Legs
	{ID: 24274, EffectID: 2748, Name: "Runic Spellthread", Quality: proto.ItemQuality_ItemQualityEpic, Bonus: stats.Stats{stats.SpellPower: 35, stats.Stamina: 20}, ItemType: proto.ItemType_ItemTypeLegs},
	{ID: 24273, EffectID: 2747, Name: "Mystic Spellthread", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{stats.SpellPower: 25, stats.Stamina: 15}, ItemType: proto.ItemType_ItemTypeLegs},
	{ID: 29533, EffectID: 3010, Name: "Cobrahide Leg Armor", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{stats.AttackPower: 40, stats.RangedAttackPower: 40, stats.MeleeCrit: 10}, ItemType: proto.ItemType_ItemTypeLegs},
	{ID: 29535, EffectID: 3012, Name: "Nethercobra Leg Armor", Quality: proto.ItemQuality_ItemQualityEpic, Bonus: stats.Stats{stats.AttackPower: 50, stats.RangedAttackPower: 50, stats.MeleeCrit: 12}, ItemType: proto.ItemType_ItemTypeLegs},
	{ID: 29536, EffectID: 3013, Name: "Nethercleft Leg Armor", Quality: proto.ItemQuality_ItemQualityEpic, Bonus: stats.Stats{stats.Stamina: 40, stats.Agility: 12}, ItemType: proto.ItemType_ItemTypeLegs},

	// Feet
	{ID: 16220, EffectID: 851, Name: "Enchant Boots - Spirit", Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{stats.Spirit: 5}, ItemType: proto.ItemType_ItemTypeFeet},
	{ID: 35297, EffectID: 2940, Name: "Enchant Boots - Boar's Speed", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{stats.Stamina: 9}, ItemType: proto.ItemType_ItemTypeFeet},
	{ID: 35298, EffectID: 2656, Name: "Enchant Boots - Vitality", Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{stats.MP5: 4}, ItemType: proto.ItemType_ItemTypeFeet},
	{ID: 22543, EffectID: 2649, Name: "Enchant Boots - Fortitude", Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{stats.Stamina: 12}, ItemType: proto.ItemType_ItemTypeFeet},
	{ID: 22544, EffectID: 2657, Name: "Enchant Boots - Dexterity", Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{stats.Agility: 12}, ItemType: proto.ItemType_ItemTypeFeet},
	{ID: 28279, EffectID: 2939, Name: "Enchant Boots - Cat's Swiftness", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{stats.Agility: 6}, ItemType: proto.ItemType_ItemTypeFeet},
	{ID: 22545, EffectID: 2658, Name: "Enchant Boots - Surefooted", Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{stats.MeleeHit: 10}, ItemType: proto.ItemType_ItemTypeFeet},

	// Weapon
	{ID: 16250, EffectID: 1897, Name: "Weapon - Superior Striking", Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{}, ItemType: proto.ItemType_ItemTypeWeapon},
	{ID: 22552, EffectID: 963, Name: "Weapon - Major Striking", Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{}, ItemType: proto.ItemType_ItemTypeWeapon},
	{ID: 16252, EffectID: 1900, Name: "Weapon - Crusader", Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{}, ItemType: proto.ItemType_ItemTypeWeapon},
	{ID: 22551, EffectID: 2666, Name: "Weapon - Major Intellect", Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{stats.Intellect: 30}, ItemType: proto.ItemType_ItemTypeWeapon},
	{ID: 22554, EffectID: 2667, Name: "Savagery", Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{stats.AttackPower: 70, stats.RangedAttackPower: 70}, ItemType: proto.ItemType_ItemTypeWeapon, EnchantType: proto.EnchantType_EnchantTypeTwoHand},
	{ID: 22555, EffectID: 2669, Name: "Weapon - Major Spellpower", Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{stats.SpellPower: 40}, ItemType: proto.ItemType_ItemTypeWeapon},
	{ID: 22560, EffectID: 2671, Name: "Sunfire", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{stats.ArcaneSpellPower: 50, stats.FireSpellPower: 50}, ItemType: proto.ItemType_ItemTypeWeapon},
	{ID: 22561, EffectID: 2672, Name: "Soulfrost", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{stats.FrostSpellPower: 54, stats.ShadowSpellPower: 54}, ItemType: proto.ItemType_ItemTypeWeapon},
	{ID: 22559, EffectID: 2673, Name: "Weapon - Mongoose", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{}, ItemType: proto.ItemType_ItemTypeWeapon},
	{ID: 19445, EffectID: 2564, Name: "Weapon - Agility", Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.Agility: 15}, ItemType: proto.ItemType_ItemTypeWeapon},
	{ID: 33165, EffectID: 3222, Name: "Weapon - Greater Agility", Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.Agility: 20}, ItemType: proto.ItemType_ItemTypeWeapon},
	{ID: 22556, EffectID: 2670, Name: "2H Weapon - Major Agility", Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{stats.Agility: 35}, ItemType: proto.ItemType_ItemTypeWeapon, EnchantType: proto.EnchantType_EnchantTypeTwoHand},
	{ID: 33307, EffectID: 3225, Name: "Weapon - Executioner", Phase: 4, Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{}, ItemType: proto.ItemType_ItemTypeWeapon},
	{ID: 35498, EffectID: 3273, Name: "Deathfrost", Phase: 5, Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{}, ItemType: proto.ItemType_ItemTypeWeapon},
	// TODO: spellsurge

	// Shield
	{ID: 22539, EffectID: 2654, Name: "Shield - Intellect", Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.Intellect: 12}, ItemType: proto.ItemType_ItemTypeWeapon, EnchantType: proto.EnchantType_EnchantTypeShield},
	{ID: 28282, EffectID: 1071, Name: "Shield - Major Stamina", Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.Stamina: 18}, ItemType: proto.ItemType_ItemTypeWeapon, EnchantType: proto.EnchantType_EnchantTypeShield},
	{ID: 44383, EffectID: 3229, Name: "Shield - Resilience", IsSpellID: true, Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.Resilience: 12}, ItemType: proto.ItemType_ItemTypeWeapon, EnchantType: proto.EnchantType_EnchantTypeShield},

	// Ring
	{ID: 22535, EffectID: 2929, Name: "Ring - Striking", Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{}, ItemType: proto.ItemType_ItemTypeFinger},
	{ID: 22536, EffectID: 2928, Name: "Ring - Spellpower", Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.SpellPower: 12}, ItemType: proto.ItemType_ItemTypeFinger},
	{ID: 22538, EffectID: 2931, Name: "Ring - Stats", Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.Stamina: 4, stats.Intellect: 4, stats.Spirit: 4, stats.Strength: 4, stats.Agility: 4}, ItemType: proto.ItemType_ItemTypeFinger},

	// Ranged
	{ID: 18283, EffectID: 2523, Name: "Biznicks 247x128 Accurascope", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{}, ItemType: proto.ItemType_ItemTypeRanged},
	{ID: 23765, EffectID: 2723, Name: "Khorium Scope", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{}, ItemType: proto.ItemType_ItemTypeRanged},
	{ID: 23766, EffectID: 2724, Name: "Stabilized Eternium Scope", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{}, ItemType: proto.ItemType_ItemTypeRanged},
}
