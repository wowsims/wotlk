package database

import (
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

// Note: EffectId AND SpellId are required for all enchants, because they are
// used by various importers/exporters. ItemId is optional.

var EnchantOverrides = []*proto.UIEnchant{
	// Multi-slot
	{EffectId: 2988, ItemId: 29487, SpellId: 35419, Name: "Nature Armor Kit", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.NatureResistance: 8}.ToFloatArray(), Type: proto.ItemType_ItemTypeChest, ExtraTypes: []proto.ItemType{proto.ItemType_ItemTypeLegs, proto.ItemType_ItemTypeHands, proto.ItemType_ItemTypeFeet}, EnchantType: proto.EnchantType_EnchantTypeKit},
	{EffectId: 3329, ItemId: 38375, SpellId: 50906, Name: "Borean Armor Kit", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Stamina: 12}.ToFloatArray(), Type: proto.ItemType_ItemTypeHead, ExtraTypes: []proto.ItemType{proto.ItemType_ItemTypeChest, proto.ItemType_ItemTypeShoulder, proto.ItemType_ItemTypeLegs, proto.ItemType_ItemTypeHands, proto.ItemType_ItemTypeFeet}, EnchantType: proto.EnchantType_EnchantTypeKit},
	{EffectId: 3330, ItemId: 38376, SpellId: 50909, Name: "Heavy Borean Armor Kit", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.Stamina: 18}.ToFloatArray(), Type: proto.ItemType_ItemTypeHead, ExtraTypes: []proto.ItemType{proto.ItemType_ItemTypeChest, proto.ItemType_ItemTypeShoulder, proto.ItemType_ItemTypeLegs, proto.ItemType_ItemTypeHands, proto.ItemType_ItemTypeFeet}, EnchantType: proto.EnchantType_EnchantTypeKit},

	// Head
	{EffectId: 3795, ItemId: 44069, SpellId: 59777, Name: "Arcanum of Triumph", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.AttackPower: 50, stats.RangedAttackPower: 50, stats.Resilience: 20}.ToFloatArray(), Type: proto.ItemType_ItemTypeHead},
	{EffectId: 3796, ItemId: 44075, SpellId: 59784, Name: "Arcanum of Dominance", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.SpellPower: 29, stats.Resilience: 20}.ToFloatArray(), Type: proto.ItemType_ItemTypeHead},
	{EffectId: 3842, ItemId: 44875, SpellId: 61271, Name: "Arcanum of the Savage Gladiator", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.Stamina: 30, stats.Resilience: 25}.ToFloatArray(), Type: proto.ItemType_ItemTypeHead},
	{EffectId: 3812, ItemId: 44137, SpellId: 59944, Name: "Arcanum of the Frosty Soul", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.Stamina: 30, stats.FrostResistance: 25}.ToFloatArray(), Type: proto.ItemType_ItemTypeHead},
	{EffectId: 3813, ItemId: 44138, SpellId: 59945, Name: "Arcanum of Toxic Warding", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.Stamina: 30, stats.NatureResistance: 25}.ToFloatArray(), Type: proto.ItemType_ItemTypeHead},
	{EffectId: 3814, ItemId: 44139, SpellId: 59946, Name: "Arcanum of the Fleeing Shadow", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.Stamina: 30, stats.ShadowResistance: 25}.ToFloatArray(), Type: proto.ItemType_ItemTypeHead},
	{EffectId: 3815, ItemId: 44140, SpellId: 59947, Name: "Arcanum of the Eclipsed Moon", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.Stamina: 30, stats.ArcaneResistance: 25}.ToFloatArray(), Type: proto.ItemType_ItemTypeHead},
	{EffectId: 3816, ItemId: 44141, SpellId: 59948, Name: "Arcanum of the Flame's Soul", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.Stamina: 30, stats.FireResistance: 25}.ToFloatArray(), Type: proto.ItemType_ItemTypeHead},
	{EffectId: 3819, ItemId: 44876, SpellId: 59960, Name: "Arcanum of Blissful Mending", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.SpellPower: 30, stats.MP5: 10}.ToFloatArray(), Type: proto.ItemType_ItemTypeHead},
	{EffectId: 3820, ItemId: 44877, SpellId: 59970, Name: "Arcanum of Burning Mysteries", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.SpellPower: 30, stats.MeleeCrit: 20, stats.SpellCrit: 20}.ToFloatArray(), Type: proto.ItemType_ItemTypeHead},
	{EffectId: 3818, ItemId: 44878, SpellId: 59955, Name: "Arcanum of the Stalwart Protector", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.Stamina: 37, stats.Defense: 20}.ToFloatArray(), Type: proto.ItemType_ItemTypeHead},
	{EffectId: 3817, ItemId: 44879, SpellId: 59954, Name: "Arcanum of Torment", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.AttackPower: 50, stats.RangedAttackPower: 50, stats.MeleeCrit: 20, stats.SpellCrit: 20}.ToFloatArray(), Type: proto.ItemType_ItemTypeHead},
	{EffectId: 3878, SpellId: 67839, Name: "Mind Amplification Dish", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Stamina: 45}.ToFloatArray(), Type: proto.ItemType_ItemTypeHead, RequiredProfession: proto.Profession_Engineering},

	// Shoulder
	{EffectId: 2998, ItemId: 29187, SpellId: 35441, Name: "Inscription of Endurance", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.ArcaneResistance: 7, stats.FireResistance: 7, stats.FrostResistance: 7, stats.NatureResistance: 7, stats.ShadowResistance: 7}.ToFloatArray(), Type: proto.ItemType_ItemTypeShoulder},
	{EffectId: 3793, ItemId: 44067, SpellId: 59771, Name: "Inscription of Triumph", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.AttackPower: 40, stats.RangedAttackPower: 40, stats.Resilience: 15}.ToFloatArray(), Type: proto.ItemType_ItemTypeShoulder},
	{EffectId: 3794, ItemId: 44068, SpellId: 59773, Name: "Inscription of Dominance", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.SpellPower: 23, stats.Resilience: 15}.ToFloatArray(), Type: proto.ItemType_ItemTypeShoulder},
	{EffectId: 3852, ItemId: 44957, SpellId: 62384, Name: "Greater Inscription of the Gladiator", Quality: proto.ItemQuality_ItemQualityEpic, Stats: stats.Stats{stats.Stamina: 30, stats.Resilience: 15}.ToFloatArray(), Type: proto.ItemType_ItemTypeShoulder},
	{EffectId: 3806, ItemId: 44129, SpellId: 59927, Name: "Lesser Inscription of the Storm", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.SpellPower: 18, stats.MeleeCrit: 10, stats.SpellCrit: 10}.ToFloatArray(), Type: proto.ItemType_ItemTypeShoulder},
	{EffectId: 3807, ItemId: 44130, SpellId: 59928, Name: "Lesser Inscription of the Crag", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.SpellPower: 18, stats.MP5: 5}.ToFloatArray(), Type: proto.ItemType_ItemTypeShoulder},
	{EffectId: 3875, ItemId: 44131, SpellId: 59929, Name: "Lesser Inscription of the Axe", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.AttackPower: 30, stats.RangedAttackPower: 30, stats.MeleeCrit: 10, stats.SpellCrit: 10}.ToFloatArray(), Type: proto.ItemType_ItemTypeShoulder},
	{EffectId: 3876, ItemId: 44132, SpellId: 59932, Name: "Lesser Inscription of the Pinnacle", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.Dodge: 15, stats.Defense: 10}.ToFloatArray(), Type: proto.ItemType_ItemTypeShoulder},
	{EffectId: 3810, ItemId: 44874, SpellId: 59937, Name: "Greater Inscription of the Storm", Quality: proto.ItemQuality_ItemQualityEpic, Stats: stats.Stats{stats.SpellPower: 24, stats.MeleeCrit: 15, stats.SpellCrit: 15}.ToFloatArray(), Type: proto.ItemType_ItemTypeShoulder},
	{EffectId: 3809, ItemId: 44872, SpellId: 59936, Name: "Greater Inscription of the Crag", Quality: proto.ItemQuality_ItemQualityEpic, Stats: stats.Stats{stats.SpellPower: 24, stats.MP5: 8}.ToFloatArray(), Type: proto.ItemType_ItemTypeShoulder},
	{EffectId: 3808, ItemId: 44871, SpellId: 59934, Name: "Greater Inscription of the Axe", Quality: proto.ItemQuality_ItemQualityEpic, Stats: stats.Stats{stats.AttackPower: 40, stats.RangedAttackPower: 40, stats.MeleeCrit: 15, stats.SpellCrit: 15}.ToFloatArray(), Type: proto.ItemType_ItemTypeShoulder},
	{EffectId: 3811, ItemId: 44873, SpellId: 59941, Name: "Greater Inscription of the Pinnacle", Quality: proto.ItemQuality_ItemQualityEpic, Stats: stats.Stats{stats.Dodge: 20, stats.Defense: 15}.ToFloatArray(), Type: proto.ItemType_ItemTypeShoulder},
	{EffectId: 3838, SpellId: 61120, Name: "Master's Inscription of the Storm", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.SpellPower: 70, stats.MeleeCrit: 15, stats.SpellCrit: 15}.ToFloatArray(), Type: proto.ItemType_ItemTypeShoulder, RequiredProfession: proto.Profession_Inscription},
	{EffectId: 3836, SpellId: 61118, Name: "Master's Inscription of the Crag", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.SpellPower: 70, stats.MP5: 8}.ToFloatArray(), Type: proto.ItemType_ItemTypeShoulder, RequiredProfession: proto.Profession_Inscription},
	{EffectId: 3835, SpellId: 61117, Name: "Master's Inscription of the Axe", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.AttackPower: 120, stats.RangedAttackPower: 120, stats.MeleeCrit: 15, stats.SpellCrit: 15}.ToFloatArray(), Type: proto.ItemType_ItemTypeShoulder, RequiredProfession: proto.Profession_Inscription},
	{EffectId: 3837, SpellId: 61119, Name: "Master's Inscription of the Pinnacle", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Dodge: 60, stats.Defense: 15}.ToFloatArray(), Type: proto.ItemType_ItemTypeShoulder, RequiredProfession: proto.Profession_Inscription},

	// Back
	{EffectId: 1262, ItemId: 37330, SpellId: 44596, Name: "Superior Arcane Resistance", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.ArcaneResistance: 20}.ToFloatArray(), Type: proto.ItemType_ItemTypeBack},
	{EffectId: 1354, ItemId: 37331, SpellId: 44556, Name: "Superior Fire Resistance", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.FireResistance: 20}.ToFloatArray(), Type: proto.ItemType_ItemTypeBack},
	{EffectId: 3230, SpellId: 44483, Name: "Superior Frost Resistance", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.FrostResistance: 20}.ToFloatArray(), Type: proto.ItemType_ItemTypeBack},
	{EffectId: 1400, SpellId: 44494, Name: "Superior Nature Resistance", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.NatureResistance: 20}.ToFloatArray(), Type: proto.ItemType_ItemTypeBack},
	{EffectId: 1446, SpellId: 44590, Name: "Superior Shadow Resistance", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.ShadowResistance: 20}.ToFloatArray(), Type: proto.ItemType_ItemTypeBack},
	{EffectId: 1951, ItemId: 37347, SpellId: 44591, Name: "Titanweave", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.Defense: 16}.ToFloatArray(), Type: proto.ItemType_ItemTypeBack},
	{EffectId: 3256, ItemId: 37349, SpellId: 44631, Name: "Shadow Armor", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.Agility: 10}.ToFloatArray(), Type: proto.ItemType_ItemTypeBack},
	{EffectId: 3294, ItemId: 44471, SpellId: 47672, Name: "Mighty Armor", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.BonusArmor: 225}.ToFloatArray(), Type: proto.ItemType_ItemTypeBack},
	{EffectId: 3831, ItemId: 44472, SpellId: 47898, Name: "Greater Speed", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.MeleeHaste: 23, stats.SpellHaste: 23}.ToFloatArray(), Type: proto.ItemType_ItemTypeBack},
	{EffectId: 3296, ItemId: 44488, SpellId: 47899, Name: "Wisdom", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.Spirit: 10}.ToFloatArray(), Type: proto.ItemType_ItemTypeBack},
	{EffectId: 3243, SpellId: 44582, Name: "Spell Piercing", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.SpellPenetration: 35}.ToFloatArray(), Type: proto.ItemType_ItemTypeBack},
	{EffectId: 3825, SpellId: 60609, Name: "Speed", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.MeleeHaste: 15, stats.SpellHaste: 15}.ToFloatArray(), Type: proto.ItemType_ItemTypeBack},
	{EffectId: 983, SpellId: 44500, Name: "Superior Agility", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Agility: 16}.ToFloatArray(), Type: proto.ItemType_ItemTypeBack},
	{EffectId: 1099, SpellId: 60663, Name: "Major Agility", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Agility: 22}.ToFloatArray(), Type: proto.ItemType_ItemTypeBack},
	{EffectId: 3605, SpellId: 55002, Name: "Flexweave Underlay", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Agility: 23}.ToFloatArray(), Type: proto.ItemType_ItemTypeBack, RequiredProfession: proto.Profession_Engineering},
	{EffectId: 3722, SpellId: 55642, Name: "Lightweave Embroidery", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeBack, RequiredProfession: proto.Profession_Tailoring},
	{EffectId: 3728, SpellId: 55769, Name: "Darkglow Embroidery", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeBack, RequiredProfession: proto.Profession_Tailoring},
	{EffectId: 3730, SpellId: 55777, Name: "Swordguard Embroidery", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeBack, RequiredProfession: proto.Profession_Tailoring},
	{EffectId: 3859, SpellId: 63765, Name: "Springy Arachnoweave", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.SpellPower: 27}.ToFloatArray(), Type: proto.ItemType_ItemTypeBack, RequiredProfession: proto.Profession_Engineering},

	// Chest
	{EffectId: 3245, ItemId: 37340, SpellId: 44588, Name: "Exceptional Resilience", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.Resilience: 20}.ToFloatArray(), Type: proto.ItemType_ItemTypeChest},
	{EffectId: 3252, SpellId: 44623, Name: "Super Stats", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.Stamina: 8, stats.Strength: 8, stats.Agility: 8, stats.Intellect: 8, stats.Spirit: 8}.ToFloatArray(), Type: proto.ItemType_ItemTypeChest},
	{EffectId: 3832, ItemId: 44489, SpellId: 60692, Name: "Powerful Stats", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.Stamina: 10, stats.Strength: 10, stats.Agility: 10, stats.Intellect: 10, stats.Spirit: 10}.ToFloatArray(), Type: proto.ItemType_ItemTypeChest},
	{EffectId: 3233, SpellId: 27958, Name: "Exceptional Mana", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Mana: 250}.ToFloatArray(), Type: proto.ItemType_ItemTypeChest},
	{EffectId: 3236, SpellId: 44492, Name: "Mighty Health", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Health: 200}.ToFloatArray(), Type: proto.ItemType_ItemTypeChest},
	{EffectId: 3297, SpellId: 47900, Name: "Super Health", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Health: 275}.ToFloatArray(), Type: proto.ItemType_ItemTypeChest},
	{EffectId: 2381, SpellId: 44509, Name: "Greater Mana Restoration", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.MP5: 10}.ToFloatArray(), Type: proto.ItemType_ItemTypeChest},
	{EffectId: 1953, SpellId: 47766, Name: "Greater Defense", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Defense: 22}.ToFloatArray(), Type: proto.ItemType_ItemTypeChest},

	// Wrist
	{EffectId: 3845, ItemId: 44484, SpellId: 44575, Name: "Greater Assault", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.AttackPower: 50, stats.RangedAttackPower: 50}.ToFloatArray(), Type: proto.ItemType_ItemTypeWrist},
	{EffectId: 2332, ItemId: 44498, SpellId: 60767, Name: "Superior Spellpower", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.SpellPower: 30}.ToFloatArray(), Type: proto.ItemType_ItemTypeWrist},
	{EffectId: 3850, ItemId: 44944, SpellId: 62256, Name: "Major Stamina", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.Stamina: 40}.ToFloatArray(), Type: proto.ItemType_ItemTypeWrist},
	{EffectId: 1119, SpellId: 44555, Name: "Exceptional Intellect", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Intellect: 16}.ToFloatArray(), Type: proto.ItemType_ItemTypeWrist},
	{EffectId: 1147, SpellId: 44593, Name: "Major Spirit", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Spirit: 18}.ToFloatArray(), Type: proto.ItemType_ItemTypeWrist},
	{EffectId: 3231, SpellId: 44598, Name: "Expertise", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Expertise: 15}.ToFloatArray(), Type: proto.ItemType_ItemTypeWrist},
	{EffectId: 2661, SpellId: 44616, Name: "Greater Stats", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Stamina: 6, stats.Strength: 6, stats.Agility: 6, stats.Intellect: 6, stats.Spirit: 6}.ToFloatArray(), Type: proto.ItemType_ItemTypeWrist},
	{EffectId: 2326, SpellId: 44635, Name: "Greater Spellpower", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.SpellPower: 23}.ToFloatArray(), Type: proto.ItemType_ItemTypeWrist},
	{EffectId: 1600, SpellId: 60616, Name: "Striking", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.AttackPower: 38, stats.RangedAttackPower: 38}.ToFloatArray(), Type: proto.ItemType_ItemTypeWrist},
	{EffectId: 3756, SpellId: 57683, Name: "Fur Lining - Attack Power", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.AttackPower: 130, stats.RangedAttackPower: 130}.ToFloatArray(), Type: proto.ItemType_ItemTypeWrist, RequiredProfession: proto.Profession_Leatherworking},
	{EffectId: 3757, SpellId: 57690, Name: "Fur Lining - Stamina", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Stamina: 102}.ToFloatArray(), Type: proto.ItemType_ItemTypeWrist, RequiredProfession: proto.Profession_Leatherworking},
	{EffectId: 3758, SpellId: 57691, Name: "Fur Lining - Spell Power", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.SpellPower: 76}.ToFloatArray(), Type: proto.ItemType_ItemTypeWrist, RequiredProfession: proto.Profession_Leatherworking},
	{EffectId: 3759, SpellId: 57692, Name: "Fur Lining - Fire Resist", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.FireResistance: 70}.ToFloatArray(), Type: proto.ItemType_ItemTypeWrist, RequiredProfession: proto.Profession_Leatherworking},
	{EffectId: 3760, SpellId: 57694, Name: "Fur Lining - Frost Resist", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.FrostResistance: 70}.ToFloatArray(), Type: proto.ItemType_ItemTypeWrist, RequiredProfession: proto.Profession_Leatherworking},
	{EffectId: 3761, SpellId: 57696, Name: "Fur Lining - Shadow Resist", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.ShadowResistance: 70}.ToFloatArray(), Type: proto.ItemType_ItemTypeWrist, RequiredProfession: proto.Profession_Leatherworking},
	{EffectId: 3762, SpellId: 57699, Name: "Fur Lining - Nature Resist", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.NatureResistance: 70}.ToFloatArray(), Type: proto.ItemType_ItemTypeWrist, RequiredProfession: proto.Profession_Leatherworking},
	{EffectId: 3763, SpellId: 57701, Name: "Fur Lining - Arcane Resist", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.ArcaneResistance: 70}.ToFloatArray(), Type: proto.ItemType_ItemTypeWrist, RequiredProfession: proto.Profession_Leatherworking},

	// Hands
	{EffectId: 3253, ItemId: 44485, SpellId: 44625, Name: "Armsman", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.Parry: 10}.ToFloatArray(), Type: proto.ItemType_ItemTypeHands},
	{EffectId: 1603, SpellId: 60668, Name: "Crusher", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.AttackPower: 44, stats.RangedAttackPower: 44}.ToFloatArray(), Type: proto.ItemType_ItemTypeHands},
	{EffectId: 3246, SpellId: 44592, Name: "Exceptional Spellpower", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.SpellPower: 28}.ToFloatArray(), Type: proto.ItemType_ItemTypeHands},
	{EffectId: 3231, SpellId: 44484, Name: "Expertise", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Expertise: 15}.ToFloatArray(), Type: proto.ItemType_ItemTypeHands},
	{EffectId: 3238, SpellId: 44506, Name: "Gatherer", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeHands},
	{EffectId: 3829, SpellId: 44513, Name: "Greater Assult", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.AttackPower: 35, stats.RangedAttackPower: 35}.ToFloatArray(), Type: proto.ItemType_ItemTypeHands},
	{EffectId: 3222, SpellId: 44529, Name: "Major Agility", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Agility: 20}.ToFloatArray(), Type: proto.ItemType_ItemTypeHands},
	{EffectId: 3234, SpellId: 44488, Name: "Precision", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.MeleeHit: 20, stats.SpellHit: 20}.ToFloatArray(), Type: proto.ItemType_ItemTypeHands},
	{EffectId: 3603, SpellId: 54998, Name: "Hand-Mounted Pyro Rocket", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeHands, RequiredProfession: proto.Profession_Engineering},
	{EffectId: 3604, SpellId: 54999, Name: "Hyperspeed Accelerators", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeHands, RequiredProfession: proto.Profession_Engineering},
	{EffectId: 3860, SpellId: 63770, Name: "Reticulated Armor Webbing", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.BonusArmor: 885}.ToFloatArray(), Type: proto.ItemType_ItemTypeHands, RequiredProfession: proto.Profession_Engineering},

	// Waist
	{EffectId: 3599, SpellId: 54736, Name: "Personal Electromagnetic Pulse Generator", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeWaist, RequiredProfession: proto.Profession_Engineering},
	{EffectId: 3601, SpellId: 54793, Name: "Frag Belt", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeWaist, RequiredProfession: proto.Profession_Engineering},

	// Legs
	{EffectId: 3325, ItemId: 38371, SpellId: 50901, Name: "Jormungar Leg Armor", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.Stamina: 45, stats.Agility: 15}.ToFloatArray(), Type: proto.ItemType_ItemTypeLegs},
	{EffectId: 3326, ItemId: 38372, SpellId: 50902, Name: "Nerubian Leg Armor", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.AttackPower: 55, stats.RangedAttackPower: 55, stats.MeleeCrit: 15, stats.SpellCrit: 15}.ToFloatArray(), Type: proto.ItemType_ItemTypeLegs},
	{EffectId: 3822, ItemId: 38373, SpellId: 60581, Name: "Frosthide Leg Armor", Quality: proto.ItemQuality_ItemQualityEpic, Stats: stats.Stats{stats.Stamina: 55, stats.Agility: 22}.ToFloatArray(), Type: proto.ItemType_ItemTypeLegs},
	{EffectId: 3823, ItemId: 38374, SpellId: 60582, Name: "Icescale Leg Armor", Quality: proto.ItemQuality_ItemQualityEpic, Stats: stats.Stats{stats.AttackPower: 75, stats.RangedAttackPower: 75, stats.MeleeCrit: 22, stats.SpellCrit: 22}.ToFloatArray(), Type: proto.ItemType_ItemTypeLegs},
	{EffectId: 3853, ItemId: 44963, SpellId: 62447, Name: "Earthen Leg Armor", Quality: proto.ItemQuality_ItemQualityEpic, Stats: stats.Stats{stats.Stamina: 28, stats.Resilience: 40}.ToFloatArray(), Type: proto.ItemType_ItemTypeLegs},
	{EffectId: 3718, ItemId: 41601, SpellId: 55630, Name: "Shining Spellthread", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.Spirit: 12, stats.SpellPower: 35}.ToFloatArray(), Type: proto.ItemType_ItemTypeLegs},
	{EffectId: 3719, ItemId: 41602, SpellId: 55631, Name: "Brilliant Spellthread", Quality: proto.ItemQuality_ItemQualityEpic, Stats: stats.Stats{stats.Spirit: 20, stats.SpellPower: 50}.ToFloatArray(), Type: proto.ItemType_ItemTypeLegs},
	{EffectId: 3720, ItemId: 41603, SpellId: 55632, Name: "Azure Spellthread", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.Stamina: 20, stats.SpellPower: 35}.ToFloatArray(), Type: proto.ItemType_ItemTypeLegs},
	{EffectId: 3721, ItemId: 41604, SpellId: 55634, Name: "Sapphire Spellthread", Quality: proto.ItemQuality_ItemQualityEpic, Stats: stats.Stats{stats.Stamina: 30, stats.SpellPower: 50}.ToFloatArray(), Type: proto.ItemType_ItemTypeLegs},
	{EffectId: 3327, SpellId: 60583, Name: "Jormungar Leg Reinforcements", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Stamina: 55, stats.Agility: 22}.ToFloatArray(), Type: proto.ItemType_ItemTypeLegs, RequiredProfession: proto.Profession_Leatherworking},
	{EffectId: 3328, SpellId: 60584, Name: "Nerubian Leg Reinforcements", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.AttackPower: 75, stats.RangedAttackPower: 75, stats.MeleeCrit: 22, stats.SpellCrit: 22}.ToFloatArray(), Type: proto.ItemType_ItemTypeLegs, RequiredProfession: proto.Profession_Leatherworking},
	{EffectId: 3873, SpellId: 56034, Name: "Master's Spellthread", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Stamina: 30, stats.SpellPower: 50}.ToFloatArray(), Type: proto.ItemType_ItemTypeLegs, RequiredProfession: proto.Profession_Tailoring},
	{EffectId: 3872, SpellId: 56039, Name: "Sanctified Spellthread", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Spirit: 20, stats.SpellPower: 50}.ToFloatArray(), Type: proto.ItemType_ItemTypeLegs, RequiredProfession: proto.Profession_Tailoring},

	// Feet
	{EffectId: 1597, ItemId: 44490, SpellId: 60763, Name: "Greater Assault", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.AttackPower: 32, stats.RangedAttackPower: 32}.ToFloatArray(), Type: proto.ItemType_ItemTypeFeet},
	{EffectId: 3232, ItemId: 44491, SpellId: 47901, Name: "Tuskarr's Vitality", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.Stamina: 15}.ToFloatArray(), Type: proto.ItemType_ItemTypeFeet},
	{EffectId: 3824, SpellId: 60606, Name: "Assault", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.AttackPower: 24, stats.RangedAttackPower: 24}.ToFloatArray(), Type: proto.ItemType_ItemTypeFeet},
	{EffectId: 1075, SpellId: 44528, Name: "Greater Fortitude", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Stamina: 22}.ToFloatArray(), Type: proto.ItemType_ItemTypeFeet},
	{EffectId: 1147, SpellId: 44508, Name: "Greater Spirit", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Spirit: 18}.ToFloatArray(), Type: proto.ItemType_ItemTypeFeet},
	{EffectId: 3244, SpellId: 44584, Name: "Greater Vitality", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.MP5: 7}.ToFloatArray(), Type: proto.ItemType_ItemTypeFeet},
	{EffectId: 3826, SpellId: 60623, Name: "Icewalker", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.MeleeHit: 12, stats.SpellHit: 12, stats.MeleeCrit: 12, stats.SpellCrit: 12}.ToFloatArray(), Type: proto.ItemType_ItemTypeFeet},
	{EffectId: 983, SpellId: 44589, Name: "Superior Agility", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Agility: 16}.ToFloatArray(), Type: proto.ItemType_ItemTypeFeet},
	{EffectId: 3606, SpellId: 55016, Name: "Nitro Boosts", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.MeleeCrit: 24, stats.SpellCrit: 24}.ToFloatArray(), Type: proto.ItemType_ItemTypeFeet, RequiredProfession: proto.Profession_Engineering},

	// Weapon
	{EffectId: 1103, SpellId: 44633, Name: "Exceptional Agility", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Agility: 26}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon},
	{EffectId: 3844, SpellId: 44510, Name: "Exceptional Spirit", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Spirit: 45}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon},
	{EffectId: 3251, ItemId: 37339, SpellId: 44621, Name: "Giant Slayer", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon},
	{EffectId: 3239, ItemId: 37344, SpellId: 44524, Name: "Icebreaker", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon},
	{EffectId: 3731, ItemId: 41976, SpellId: 55836, Name: "Titanium Weapon Chain", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.MeleeHit: 28, stats.SpellHit: 28}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon},
	{EffectId: 3833, ItemId: 44486, SpellId: 60707, Name: "Superior Potency", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.AttackPower: 65, stats.RangedAttackPower: 65}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon},
	{EffectId: 3834, ItemId: 44487, SpellId: 60714, Name: "Mighty Spellpower", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.SpellPower: 63}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon},
	{EffectId: 3789, ItemId: 44492, SpellId: 59621, Name: "Berserking", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon},
	{EffectId: 3241, ItemId: 44494, SpellId: 44576, Name: "Lifeward", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon},
	{EffectId: 3790, ItemId: 44495, SpellId: 59625, Name: "Black Magic", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon},
	{EffectId: 3788, ItemId: 44496, SpellId: 59619, Name: "Accuracy", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.MeleeHit: 25, stats.SpellHit: 25, stats.MeleeCrit: 25, stats.SpellCrit: 25}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon},
	{EffectId: 3830, SpellId: 44629, Name: "Exceptional Spellpower", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.SpellPower: 50}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon},
	{EffectId: 1606, SpellId: 60621, Name: "Greater Potency", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.AttackPower: 50, stats.RangedAttackPower: 50}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon},
	{EffectId: 3370, SpellId: 53343, Name: "Rune of Razorice", Phase: 1, Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon, ClassAllowlist: []proto.Class{proto.Class_ClassDeathknight}},
	{EffectId: 3369, SpellId: 53341, Name: "Rune of Cinderglacier", Phase: 1, Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon, ClassAllowlist: []proto.Class{proto.Class_ClassDeathknight}},
	{EffectId: 3366, SpellId: 53331, Name: "Rune of Lichbane", Phase: 1, Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon, ClassAllowlist: []proto.Class{proto.Class_ClassDeathknight}},
	{EffectId: 3595, SpellId: 54447, Name: "Rune of Spellbreaking", Phase: 1, Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon, ClassAllowlist: []proto.Class{proto.Class_ClassDeathknight}},
	{EffectId: 3594, SpellId: 54446, Name: "Rune of Swordbreaking", Phase: 1, Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon, ClassAllowlist: []proto.Class{proto.Class_ClassDeathknight}},
	{EffectId: 3368, SpellId: 53344, Name: "Rune of the Fallen Crusader", Phase: 1, Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon, ClassAllowlist: []proto.Class{proto.Class_ClassDeathknight}},
	{EffectId: 3870, ItemId: 46348, SpellId: 64579, Name: "Blood Draining", Quality: proto.ItemQuality_ItemQualityEpic, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon},
	{EffectId: 3883, SpellId: 70164, Name: "Rune of the Nerubian Carapace", Phase: 1, Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon, ClassAllowlist: []proto.Class{proto.Class_ClassDeathknight}},

	// 2H Weapon
	{EffectId: 3247, ItemId: 44473, SpellId: 44595, Name: "Scourgebane", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon, EnchantType: proto.EnchantType_EnchantTypeTwoHand},
	{EffectId: 3827, ItemId: 44483, SpellId: 60691, Name: "Massacre", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.AttackPower: 110, stats.RangedAttackPower: 110}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon, EnchantType: proto.EnchantType_EnchantTypeTwoHand},
	{EffectId: 3828, SpellId: 44630, Name: "Greater Savagery", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.AttackPower: 85, stats.RangedAttackPower: 85}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon, EnchantType: proto.EnchantType_EnchantTypeTwoHand},
	{EffectId: 3854, ItemId: 45059, SpellId: 62948, Name: "Staff - Greater Spellpower", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.SpellPower: 81}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon, EnchantType: proto.EnchantType_EnchantTypeStaff},
	{EffectId: 3367, SpellId: 53342, Name: "Rune of Spellshattering", Phase: 1, Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon, ClassAllowlist: []proto.Class{proto.Class_ClassDeathknight}},
	{EffectId: 3365, SpellId: 53323, Name: "Rune of Swordshattering", Phase: 1, Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon, ClassAllowlist: []proto.Class{proto.Class_ClassDeathknight}},
	{EffectId: 3847, SpellId: 62158, Name: "Rune of the Stoneskin Gargoyle", Phase: 1, Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon, ClassAllowlist: []proto.Class{proto.Class_ClassDeathknight}},

	// Shield
	{EffectId: 1952, SpellId: 44489, Name: "Defense", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Defense: 20}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon, EnchantType: proto.EnchantType_EnchantTypeShield},
	{EffectId: 1128, SpellId: 60653, Name: "Greater Intellect", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Intellect: 25}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon, EnchantType: proto.EnchantType_EnchantTypeShield},
	{EffectId: 3748, ItemId: 42500, SpellId: 56353, Name: "Titanium Shield Spike", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon, EnchantType: proto.EnchantType_EnchantTypeShield},
	{EffectId: 3849, ItemId: 44936, SpellId: 62201, Name: "Titanium Plating", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.BlockValue: 81}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon, EnchantType: proto.EnchantType_EnchantTypeShield},

	// Ring
	{EffectId: 3839, SpellId: 44645, Name: "Assault", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.AttackPower: 40, stats.RangedAttackPower: 40}.ToFloatArray(), Type: proto.ItemType_ItemTypeFinger, RequiredProfession: proto.Profession_Enchanting},
	{EffectId: 3840, SpellId: 44636, Name: "Greater Spellpower", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.SpellPower: 23}.ToFloatArray(), Type: proto.ItemType_ItemTypeFinger, RequiredProfession: proto.Profession_Enchanting},
	{EffectId: 3791, SpellId: 59636, Name: "Stamina", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Stamina: 30}.ToFloatArray(), Type: proto.ItemType_ItemTypeFinger, RequiredProfession: proto.Profession_Enchanting},

	// Ranged
	{EffectId: 3607, ItemId: 41146, SpellId: 55076, Name: "Sun Scope", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeRanged},
	{EffectId: 3608, ItemId: 41167, SpellId: 55135, Name: "Heartseeker Scope", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeRanged},
	{EffectId: 3843, ItemId: 44739, SpellId: 61468, Name: "Diamond-cut Refractor Scope", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeRanged},

	/////////////////////////////
	//   TBC
	/////////////////////////////

	// Head
	{EffectId: 2999, ItemId: 29186, SpellId: 35443, Name: "Arcanum of the Defender", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.Defense: 16, stats.Dodge: 17}.ToFloatArray(), Type: proto.ItemType_ItemTypeHead},
	{EffectId: 3002, ItemId: 29191, SpellId: 35447, Name: "Arcanum of Power", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.SpellPower: 22, stats.SpellHit: 14}.ToFloatArray(), Type: proto.ItemType_ItemTypeHead},
	{EffectId: 3003, ItemId: 29192, SpellId: 35452, Name: "Arcanum of Ferocity", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.AttackPower: 34, stats.RangedAttackPower: 34, stats.MeleeHit: 16, stats.SpellHit: 16}.ToFloatArray(), Type: proto.ItemType_ItemTypeHead},
	{EffectId: 3096, ItemId: 30846, SpellId: 37891, Name: "Arcanum of the Outcast", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.Strength: 17, stats.Intellect: 16}.ToFloatArray(), Type: proto.ItemType_ItemTypeHead},
	{EffectId: 3004, ItemId: 29193, SpellId: 35453, Name: "Arcanum of the Gladiator", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.Stamina: 18, stats.Resilience: 20}.ToFloatArray(), Type: proto.ItemType_ItemTypeHead},

	// ZG Head Enchants
	{EffectId: 2583, ItemId: 19782, SpellId: 24149, Name: "Presence of Might", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.Stamina: 10, stats.Defense: 10, stats.BlockValue: 15}.ToFloatArray(), Type: proto.ItemType_ItemTypeHead, ClassAllowlist: []proto.Class{proto.Class_ClassWarrior}},

	// Shoulder
	{EffectId: 2982, ItemId: 28886, SpellId: 35406, Name: "Greater Inscription of Discipline", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.SpellPower: 18, stats.SpellCrit: 10}.ToFloatArray(), Type: proto.ItemType_ItemTypeShoulder},
	{EffectId: 2986, ItemId: 28888, SpellId: 35417, Name: "Greater Inscription of Vengeance", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.AttackPower: 30, stats.RangedAttackPower: 30, stats.MeleeCrit: 10}.ToFloatArray(), Type: proto.ItemType_ItemTypeShoulder},
	{EffectId: 2978, ItemId: 28889, SpellId: 35402, Name: "Greater Inscription of Warding", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.Defense: 10, stats.Dodge: 15}.ToFloatArray(), Type: proto.ItemType_ItemTypeShoulder},
	{EffectId: 2995, ItemId: 28909, SpellId: 35437, Name: "Greater Inscription of the Orb", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.SpellPower: 12, stats.SpellCrit: 15}.ToFloatArray(), Type: proto.ItemType_ItemTypeShoulder},
	{EffectId: 2997, ItemId: 28910, SpellId: 35437, Name: "Greater Inscription of the Blade", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.AttackPower: 20, stats.RangedAttackPower: 20, stats.MeleeCrit: 15}.ToFloatArray(), Type: proto.ItemType_ItemTypeShoulder},
	{EffectId: 2991, ItemId: 28911, SpellId: 35439, Name: "Greater Inscription of the Knight", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.Defense: 15, stats.Dodge: 10}.ToFloatArray(), Type: proto.ItemType_ItemTypeShoulder},
	{EffectId: 2605, ItemId: 20076, SpellId: 24421, Name: "Zandalar Signet of Mojo", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.SpellPower: 18}.ToFloatArray(), Type: proto.ItemType_ItemTypeShoulder},
	{EffectId: 2721, ItemId: 23545, SpellId: 29467, Name: "Power of the Scourge", Quality: proto.ItemQuality_ItemQualityEpic, Stats: stats.Stats{stats.SpellPower: 15, stats.SpellCrit: 14}.ToFloatArray(), Type: proto.ItemType_ItemTypeShoulder},
	{EffectId: 2717, ItemId: 23548, SpellId: 29483, Name: "Might of the Scourge", Quality: proto.ItemQuality_ItemQualityEpic, Stats: stats.Stats{stats.AttackPower: 26, stats.RangedAttackPower: 26, stats.MeleeCrit: 14}.ToFloatArray(), Type: proto.ItemType_ItemTypeShoulder},
	{EffectId: 2716, ItemId: 23549, SpellId: 29480, Name: "Fortitude of the Scourge", Quality: proto.ItemQuality_ItemQualityEpic, Stats: stats.Stats{stats.Stamina: 16, stats.BonusArmor: 100}.ToFloatArray(), Type: proto.ItemType_ItemTypeShoulder},

	// Back
	{EffectId: 2622, ItemId: 33148, SpellId: 25086, Name: "Enchant Cloak - Dodge", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.Dodge: 12}.ToFloatArray(), Type: proto.ItemType_ItemTypeBack},
	{EffectId: 2621, ItemId: 33150, SpellId: 25084, Name: "Enchant Cloak - Subtlety", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeBack},
	{EffectId: 849, ItemId: 11206, SpellId: 13882, Name: "Enchant Cloak - Lesser Agility", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.Agility: 3}.ToFloatArray(), Type: proto.ItemType_ItemTypeBack},
	{EffectId: 368, SpellId: 34004, Name: "Enchant Cloak - Greater Agility", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Agility: 12}.ToFloatArray(), Type: proto.ItemType_ItemTypeBack},
	{EffectId: 2938, ItemId: 28274, SpellId: 34003, Name: "Enchant Cloak - Spell Penetration", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.SpellPenetration: 20}.ToFloatArray(), Type: proto.ItemType_ItemTypeBack},
	{EffectId: 1441, ItemId: 28277, SpellId: 34006, Name: "Enchant Cloak - Greater Shadow Resistance", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.ShadowResistance: 15}.ToFloatArray(), Type: proto.ItemType_ItemTypeBack},
	{EffectId: 2648, ItemId: 35756, SpellId: 47051, Name: "Enchant Cloak - Steelweave", Phase: 5, Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.Defense: 12}.ToFloatArray(), Type: proto.ItemType_ItemTypeBack},

	// Chest
	{EffectId: 2659, SpellId: 27957, Name: "Chest - Exceptional Health", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Health: 150}.ToFloatArray(), Type: proto.ItemType_ItemTypeChest},
	{EffectId: 2661, ItemId: 24003, SpellId: 27960, Name: "Chest - Exceptional Stats", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Stamina: 6, stats.Intellect: 6, stats.Spirit: 6, stats.Strength: 6, stats.Agility: 6}.ToFloatArray(), Type: proto.ItemType_ItemTypeChest},
	{EffectId: 2933, ItemId: 28270, SpellId: 33992, Name: "Chest - Major Resilience", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.Resilience: 15}.ToFloatArray(), Type: proto.ItemType_ItemTypeChest},
	{EffectId: 1144, SpellId: 33990, Name: "Chest - Major Spirit", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Spirit: 15}.ToFloatArray(), Type: proto.ItemType_ItemTypeChest},
	{EffectId: 3150, SpellId: 33991, Name: "Chest - Restore Mana Prime", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.MP5: 6}.ToFloatArray(), Type: proto.ItemType_ItemTypeChest},
	{EffectId: 1950, ItemId: 35500, SpellId: 46594, Name: "Chest - Defense", Phase: 5, Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Defense: 15}.ToFloatArray(), Type: proto.ItemType_ItemTypeChest},

	// Wrist
	{EffectId: 2649, ItemId: 22533, SpellId: 27914, Name: "Bracer - Fortitude", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.Stamina: 12}.ToFloatArray(), Type: proto.ItemType_ItemTypeWrist},
	{EffectId: 2650, ItemId: 22534, SpellId: 27917, Name: "Bracer - Spellpower", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.SpellPower: 15}.ToFloatArray(), Type: proto.ItemType_ItemTypeWrist},
	{EffectId: 369, SpellId: 34001, Name: "Bracer - Major Intellect", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.Intellect: 12}.ToFloatArray(), Type: proto.ItemType_ItemTypeWrist},
	{EffectId: 2647, SpellId: 27899, Name: "Bracer - Brawn", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.Strength: 12}.ToFloatArray(), Type: proto.ItemType_ItemTypeWrist},
	{EffectId: 1593, SpellId: 34002, Name: "Bracer - Assault", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.AttackPower: 24, stats.RangedAttackPower: 24}.ToFloatArray(), Type: proto.ItemType_ItemTypeWrist},
	{EffectId: 1891, SpellId: 27905, Name: "Bracer - Stats", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.Stamina: 4, stats.Intellect: 4, stats.Spirit: 4, stats.Strength: 4, stats.Agility: 4}.ToFloatArray(), Type: proto.ItemType_ItemTypeWrist},

	// Hands
	{EffectId: 2935, ItemId: 28271, SpellId: 33994, Name: "Gloves - Spell Strike", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.SpellHit: 15}.ToFloatArray(), Type: proto.ItemType_ItemTypeHands},
	{EffectId: 2937, ItemId: 28272, SpellId: 33997, Name: "Gloves - Major Spellpower", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.SpellPower: 20}.ToFloatArray(), Type: proto.ItemType_ItemTypeHands},
	{EffectId: 684, SpellId: 33995, Name: "Gloves - Major Strength", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.Strength: 15}.ToFloatArray(), Type: proto.ItemType_ItemTypeHands},
	{EffectId: 2564, ItemId: 33152, SpellId: 25080, Name: "Gloves - Major Agility", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.Agility: 15}.ToFloatArray(), Type: proto.ItemType_ItemTypeHands},
	{EffectId: 2613, ItemId: 33153, SpellId: 25072, Name: "Gloves - Threat", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeHands},

	// Legs
	{EffectId: 2748, ItemId: 24274, SpellId: 31372, Name: "Runic Spellthread", Quality: proto.ItemQuality_ItemQualityEpic, Stats: stats.Stats{stats.SpellPower: 35, stats.Stamina: 20}.ToFloatArray(), Type: proto.ItemType_ItemTypeLegs},
	{EffectId: 2747, ItemId: 24273, SpellId: 31371, Name: "Mystic Spellthread", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.SpellPower: 25, stats.Stamina: 15}.ToFloatArray(), Type: proto.ItemType_ItemTypeLegs},
	{EffectId: 3010, ItemId: 29533, SpellId: 35488, Name: "Cobrahide Leg Armor", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.AttackPower: 40, stats.RangedAttackPower: 40, stats.MeleeCrit: 10}.ToFloatArray(), Type: proto.ItemType_ItemTypeLegs},
	{EffectId: 3012, ItemId: 29535, SpellId: 35490, Name: "Nethercobra Leg Armor", Quality: proto.ItemQuality_ItemQualityEpic, Stats: stats.Stats{stats.AttackPower: 50, stats.RangedAttackPower: 50, stats.MeleeCrit: 12}.ToFloatArray(), Type: proto.ItemType_ItemTypeLegs},
	{EffectId: 3013, ItemId: 29536, SpellId: 35495, Name: "Nethercleft Leg Armor", Quality: proto.ItemQuality_ItemQualityEpic, Stats: stats.Stats{stats.Stamina: 40, stats.Agility: 12}.ToFloatArray(), Type: proto.ItemType_ItemTypeLegs},

	// Feet
	{EffectId: 851, ItemId: 16220, SpellId: 20024, Name: "Enchant Boots - Spirit", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.Spirit: 5}.ToFloatArray(), Type: proto.ItemType_ItemTypeFeet},
	{EffectId: 2940, ItemId: 35297, SpellId: 34008, Name: "Enchant Boots - Boar's Speed", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.Stamina: 9}.ToFloatArray(), Type: proto.ItemType_ItemTypeFeet},
	{EffectId: 2656, ItemId: 35298, SpellId: 27948, Name: "Enchant Boots - Vitality", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.MP5: 4}.ToFloatArray(), Type: proto.ItemType_ItemTypeFeet},
	{EffectId: 2649, ItemId: 22543, SpellId: 27950, Name: "Enchant Boots - Fortitude", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.Stamina: 12}.ToFloatArray(), Type: proto.ItemType_ItemTypeFeet},
	{EffectId: 2657, ItemId: 22544, SpellId: 27951, Name: "Enchant Boots - Dexterity", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.Agility: 12}.ToFloatArray(), Type: proto.ItemType_ItemTypeFeet},
	{EffectId: 2939, ItemId: 28279, SpellId: 34007, Name: "Enchant Boots - Cat's Swiftness", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.Agility: 6}.ToFloatArray(), Type: proto.ItemType_ItemTypeFeet},
	{EffectId: 2658, ItemId: 22545, SpellId: 27954, Name: "Enchant Boots - Surefooted", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.MeleeHit: 10}.ToFloatArray(), Type: proto.ItemType_ItemTypeFeet},

	// Weapon
	{EffectId: 1897, ItemId: 16250, SpellId: 20031, Name: "Superior Striking", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon},
	{EffectId: 963, ItemId: 22552, SpellId: 27967, Name: "Major Striking", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon},
	{EffectId: 1900, ItemId: 16252, SpellId: 20034, Name: "Crusader", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon},
	{EffectId: 2666, ItemId: 22551, SpellId: 27968, Name: "Major Intellect", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.Intellect: 30}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon},
	{EffectId: 2667, ItemId: 22554, SpellId: 27971, Name: "Savagery", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.AttackPower: 70, stats.RangedAttackPower: 70}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon, EnchantType: proto.EnchantType_EnchantTypeTwoHand},
	{EffectId: 2669, ItemId: 22555, SpellId: 27975, Name: "Major Spellpower", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.SpellPower: 40}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon},
	{EffectId: 2671, ItemId: 22560, SpellId: 27981, Name: "Sunfire", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon},
	{EffectId: 2672, ItemId: 22561, SpellId: 27982, Name: "Soulfrost", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon},
	{EffectId: 2673, ItemId: 22559, SpellId: 27984, Name: "Mongoose", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon},
	{EffectId: 2564, ItemId: 19445, SpellId: 23800, Name: "Agility", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Agility: 15}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon},
	{EffectId: 3222, ItemId: 33165, SpellId: 42620, Name: "Greater Agility", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Agility: 20}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon},
	{EffectId: 2670, ItemId: 22556, SpellId: 27977, Name: "2H Weapon - Major Agility", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.Agility: 35}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon, EnchantType: proto.EnchantType_EnchantTypeTwoHand},
	{EffectId: 3225, ItemId: 33307, SpellId: 42974, Name: "Executioner", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon},
	{EffectId: 3273, ItemId: 35498, SpellId: 46578, Name: "Deathfrost", Phase: 5, Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon},
	{EffectId: 3855, ItemId: 45060, SpellId: 62959, Name: "Staff - Spellpower", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.SpellPower: 69}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon, EnchantType: proto.EnchantType_EnchantTypeStaff},

	// Shield
	{EffectId: 2654, ItemId: 22539, SpellId: 27945, Name: "Intellect", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Intellect: 12}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon, EnchantType: proto.EnchantType_EnchantTypeShield},
	{EffectId: 1071, ItemId: 28282, SpellId: 34009, Name: "Major Stamina", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Stamina: 18}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon, EnchantType: proto.EnchantType_EnchantTypeShield},
	{EffectId: 3229, SpellId: 44383, Name: "Resilience", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Resilience: 12}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon, EnchantType: proto.EnchantType_EnchantTypeShield},

	// Ring
	{EffectId: 2929, ItemId: 22535, SpellId: 27920, Name: "Striking", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeFinger, RequiredProfession: proto.Profession_Enchanting},
	{EffectId: 2928, ItemId: 22536, SpellId: 27924, Name: "Spellpower", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.SpellPower: 12}.ToFloatArray(), Type: proto.ItemType_ItemTypeFinger, RequiredProfession: proto.Profession_Enchanting},
	{EffectId: 2931, ItemId: 22538, SpellId: 27927, Name: "Stats", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Stamina: 4, stats.Intellect: 4, stats.Spirit: 4, stats.Strength: 4, stats.Agility: 4}.ToFloatArray(), Type: proto.ItemType_ItemTypeFinger, RequiredProfession: proto.Profession_Enchanting},

	// Ranged
	{EffectId: 2523, ItemId: 18283, SpellId: 22779, Name: "Biznicks 247x128 Accurascope", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeRanged},
	{EffectId: 2723, ItemId: 23765, SpellId: 30252, Name: "Khorium Scope", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeRanged},
	{EffectId: 2724, ItemId: 23766, SpellId: 30260, Name: "Stabilized Eternium Scope", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeRanged},
}
