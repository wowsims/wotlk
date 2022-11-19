package items

import (
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

var Enchants = []Enchant{
	// Multi-slot
	{EffectID: 3329, ItemID: 38375, Name: "Borean Armor Kit", Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.Stamina: 12}, ItemType: proto.ItemType_ItemTypeHead, EnchantType: proto.EnchantType_EnchantTypeKit},
	{EffectID: 3330, ItemID: 38376, Name: "Heavy Borean Armor Kit", Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{stats.Stamina: 18}, ItemType: proto.ItemType_ItemTypeHead, EnchantType: proto.EnchantType_EnchantTypeKit},

	// Head
	{EffectID: 3795, ItemID: 44069, Name: "Arcanum of Triumph", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{stats.AttackPower: 50, stats.RangedAttackPower: 50, stats.Resilience: 20}, ItemType: proto.ItemType_ItemTypeHead},
	{EffectID: 3796, ItemID: 44075, Name: "Arcanum of Dominance", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{stats.SpellPower: 29, stats.Resilience: 20}, ItemType: proto.ItemType_ItemTypeHead},
	{EffectID: 3842, ItemID: 44875, Name: "Arcanum of the Savage Gladiator", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{stats.Stamina: 30, stats.Resilience: 25}, ItemType: proto.ItemType_ItemTypeHead},
	{EffectID: 3812, ItemID: 44137, Name: "Arcanum of the Frosty Soul", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{stats.Stamina: 30, stats.FrostResistance: 25}, ItemType: proto.ItemType_ItemTypeHead},
	{EffectID: 3813, ItemID: 44138, Name: "Arcanum of Toxic Warding", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{stats.Stamina: 30, stats.NatureResistance: 25}, ItemType: proto.ItemType_ItemTypeHead},
	{EffectID: 3814, ItemID: 44139, Name: "Arcanum of the Fleeing Shadow", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{stats.Stamina: 30, stats.ShadowResistance: 25}, ItemType: proto.ItemType_ItemTypeHead},
	{EffectID: 3815, ItemID: 44140, Name: "Arcanum of the Eclipsed Moon", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{stats.Stamina: 30, stats.ArcaneResistance: 25}, ItemType: proto.ItemType_ItemTypeHead},
	{EffectID: 3816, ItemID: 44141, Name: "Arcanum of the Flame's Soul", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{stats.Stamina: 30, stats.FireResistance: 25}, ItemType: proto.ItemType_ItemTypeHead},
	{EffectID: 3819, ItemID: 44876, Name: "Arcanum of Blissful Mending", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{stats.SpellPower: 30, stats.MP5: 10}, ItemType: proto.ItemType_ItemTypeHead},
	{EffectID: 3820, ItemID: 44877, Name: "Arcanum of Burning Mysteries", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{stats.SpellPower: 30, stats.MeleeCrit: 20, stats.SpellCrit: 20}, ItemType: proto.ItemType_ItemTypeHead},
	{EffectID: 3818, ItemID: 44878, Name: "Arcanum of the Stalwart Protector", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{stats.Stamina: 37, stats.Defense: 20}, ItemType: proto.ItemType_ItemTypeHead},
	{EffectID: 3817, ItemID: 44879, Name: "Arcanum of Torment", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{stats.AttackPower: 50, stats.RangedAttackPower: 50, stats.MeleeCrit: 20, stats.SpellCrit: 20}, ItemType: proto.ItemType_ItemTypeHead},
	{EffectID: 3878, SpellID: 67839, Name: "Mind Amplification Dish", Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.Stamina: 45}, ItemType: proto.ItemType_ItemTypeHead, RequiredProfession: proto.Profession_Engineering},

	// Shoulder
	{EffectID: 3793, ItemID: 44067, Name: "Inscription of Triumph", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{stats.AttackPower: 40, stats.RangedAttackPower: 40, stats.Resilience: 15}, ItemType: proto.ItemType_ItemTypeShoulder},
	{EffectID: 3794, ItemID: 44068, Name: "Inscription of Dominance", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{stats.SpellPower: 23, stats.Resilience: 15}, ItemType: proto.ItemType_ItemTypeShoulder},
	{EffectID: 3852, ItemID: 44957, Name: "Greater Inscription of the Gladiator", Quality: proto.ItemQuality_ItemQualityEpic, Bonus: stats.Stats{stats.Stamina: 30, stats.Resilience: 15}, ItemType: proto.ItemType_ItemTypeShoulder},
	{EffectID: 3806, ItemID: 44129, Name: "Lesser Inscription of the Storm", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{stats.SpellPower: 18, stats.MeleeCrit: 10, stats.SpellCrit: 10}, ItemType: proto.ItemType_ItemTypeShoulder},
	{EffectID: 3807, ItemID: 44130, Name: "Lesser Inscription of the Crag", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{stats.SpellPower: 18, stats.MP5: 5}, ItemType: proto.ItemType_ItemTypeShoulder},
	{EffectID: 3875, ItemID: 44131, Name: "Lesser Inscription of the Axe", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{stats.AttackPower: 30, stats.RangedAttackPower: 30, stats.MeleeCrit: 10, stats.SpellCrit: 10}, ItemType: proto.ItemType_ItemTypeShoulder},
	{EffectID: 3876, ItemID: 44132, Name: "Lesser Inscription of the Pinnacle", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{stats.Dodge: 15, stats.Defense: 10}, ItemType: proto.ItemType_ItemTypeShoulder},
	{EffectID: 3810, ItemID: 44874, Name: "Greater Inscription of the Storm", Quality: proto.ItemQuality_ItemQualityEpic, Bonus: stats.Stats{stats.SpellPower: 24, stats.MeleeCrit: 15, stats.SpellCrit: 15}, ItemType: proto.ItemType_ItemTypeShoulder},
	{EffectID: 3809, ItemID: 44872, Name: "Greater Inscription of the Crag", Quality: proto.ItemQuality_ItemQualityEpic, Bonus: stats.Stats{stats.SpellPower: 24, stats.MP5: 8}, ItemType: proto.ItemType_ItemTypeShoulder},
	{EffectID: 3808, ItemID: 44871, Name: "Greater Inscription of the Axe", Quality: proto.ItemQuality_ItemQualityEpic, Bonus: stats.Stats{stats.AttackPower: 40, stats.RangedAttackPower: 40, stats.MeleeCrit: 15, stats.SpellCrit: 15}, ItemType: proto.ItemType_ItemTypeShoulder},
	{EffectID: 3811, ItemID: 44873, Name: "Greater Inscription of the Pinnacle", Quality: proto.ItemQuality_ItemQualityEpic, Bonus: stats.Stats{stats.Dodge: 20, stats.Defense: 15}, ItemType: proto.ItemType_ItemTypeShoulder},
	{EffectID: 3838, SpellID: 61120, Name: "Master's Inscription of the Storm", Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.SpellPower: 70, stats.MeleeCrit: 15, stats.SpellCrit: 15}, ItemType: proto.ItemType_ItemTypeShoulder, RequiredProfession: proto.Profession_Inscription},
	{EffectID: 3836, SpellID: 61118, Name: "Master's Inscription of the Crag", Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.SpellPower: 70, stats.MP5: 8}, ItemType: proto.ItemType_ItemTypeShoulder, RequiredProfession: proto.Profession_Inscription},
	{EffectID: 3835, SpellID: 61117, Name: "Master's Inscription of the Axe", Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.AttackPower: 120, stats.RangedAttackPower: 120, stats.MeleeCrit: 15, stats.SpellCrit: 15}, ItemType: proto.ItemType_ItemTypeShoulder, RequiredProfession: proto.Profession_Inscription},
	{EffectID: 3837, SpellID: 61119, Name: "Master's Inscription of the Pinnacle", Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.Dodge: 60, stats.Defense: 15}, ItemType: proto.ItemType_ItemTypeShoulder, RequiredProfession: proto.Profession_Inscription},

	// Back
	{EffectID: 1262, ItemID: 37330, Name: "Superior Arcane Resistance", Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{stats.ArcaneResistance: 20}, ItemType: proto.ItemType_ItemTypeBack},
	{EffectID: 1354, ItemID: 37331, Name: "Superior Fire Resistance", Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{stats.FireResistance: 20}, ItemType: proto.ItemType_ItemTypeBack},
	{EffectID: 3230, SpellID: 44483, Name: "Superior Frost Resistance", Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.FrostResistance: 20}, ItemType: proto.ItemType_ItemTypeBack},
	{EffectID: 1400, SpellID: 44494, Name: "Superior Nature Resistance", Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.NatureResistance: 20}, ItemType: proto.ItemType_ItemTypeBack},
	{EffectID: 1446, SpellID: 44590, Name: "Superior Shadow Resistance", Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.ShadowResistance: 20}, ItemType: proto.ItemType_ItemTypeBack},
	{EffectID: 1951, ItemID: 37347, Name: "Titanweave", Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{stats.Defense: 16}, ItemType: proto.ItemType_ItemTypeBack},
	{EffectID: 3256, ItemID: 37349, Name: "Shadow Armor", Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{stats.Agility: 10}, ItemType: proto.ItemType_ItemTypeBack},
	{EffectID: 3294, ItemID: 44471, Name: "Mighty Armor", Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{stats.Armor: 225}, ItemType: proto.ItemType_ItemTypeBack},
	{EffectID: 3831, ItemID: 44472, Name: "Greater Speed", Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{stats.MeleeHaste: 23, stats.SpellHaste: 23}, ItemType: proto.ItemType_ItemTypeBack},
	{EffectID: 3296, ItemID: 44488, Name: "Wisdom", Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{stats.Spirit: 10}, ItemType: proto.ItemType_ItemTypeBack},
	{EffectID: 3243, SpellID: 44582, Name: "Spell Piercing", Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.SpellPenetration: 35}, ItemType: proto.ItemType_ItemTypeBack},
	{EffectID: 3825, SpellID: 60609, Name: "Speed", Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.MeleeHaste: 15, stats.SpellHaste: 15}, ItemType: proto.ItemType_ItemTypeBack},
	{EffectID: 983, SpellID: 44500, Name: "Superior Agility", Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.Agility: 16}, ItemType: proto.ItemType_ItemTypeBack},
	{EffectID: 1099, SpellID: 60663, Name: "Major Agility", Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.Agility: 22}, ItemType: proto.ItemType_ItemTypeBack},
	{EffectID: 3605, SpellID: 55002, Name: "Flexweave Underlay", Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.Agility: 23}, ItemType: proto.ItemType_ItemTypeBack, RequiredProfession: proto.Profession_Engineering},
	{EffectID: 3722, SpellID: 55642, Name: "Lightweave Embroidery", Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{}, ItemType: proto.ItemType_ItemTypeBack, RequiredProfession: proto.Profession_Tailoring},
	{EffectID: 3728, SpellID: 55769, Name: "Darkglow Embroidery", Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{}, ItemType: proto.ItemType_ItemTypeBack, RequiredProfession: proto.Profession_Tailoring},
	{EffectID: 3730, SpellID: 55777, Name: "Swordguard Embroidery", Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{}, ItemType: proto.ItemType_ItemTypeBack, RequiredProfession: proto.Profession_Tailoring},
	{EffectID: 3859, SpellID: 63765, Name: "Springy Arachnoweave", Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.SpellPower: 27}, ItemType: proto.ItemType_ItemTypeBack, RequiredProfession: proto.Profession_Engineering},

	// Chest
	{EffectID: 3245, ItemID: 37340, Name: "Exceptional Resilience", Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{stats.Resilience: 20}, ItemType: proto.ItemType_ItemTypeChest},
	{EffectID: 3252, SpellID: 44623, Name: "Super Stats", Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{stats.Stamina: 8, stats.Strength: 8, stats.Agility: 8, stats.Intellect: 8, stats.Spirit: 8}, ItemType: proto.ItemType_ItemTypeChest},
	{EffectID: 3832, ItemID: 44489, Name: "Powerful Stats", Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{stats.Stamina: 10, stats.Strength: 10, stats.Agility: 10, stats.Intellect: 10, stats.Spirit: 10}, ItemType: proto.ItemType_ItemTypeChest},
	{EffectID: 3233, SpellID: 27958, Name: "Exceptional Mana", Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.Mana: 250}, ItemType: proto.ItemType_ItemTypeChest},
	{EffectID: 3236, SpellID: 44492, Name: "Mighty Health", Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.Health: 200}, ItemType: proto.ItemType_ItemTypeChest},
	{EffectID: 3297, SpellID: 47900, Name: "Super Health", Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.Health: 275}, ItemType: proto.ItemType_ItemTypeChest},
	{EffectID: 2381, SpellID: 44509, Name: "Greater Mana Restoration", Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.MP5: 10}, ItemType: proto.ItemType_ItemTypeChest},
	{EffectID: 1953, SpellID: 47766, Name: "Greater Defense", Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.Defense: 22}, ItemType: proto.ItemType_ItemTypeChest},

	// Wrist
	{EffectID: 3845, ItemID: 44484, Name: "Greater Assault", Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{stats.AttackPower: 50, stats.RangedAttackPower: 50}, ItemType: proto.ItemType_ItemTypeWrist},
	{EffectID: 2332, ItemID: 44498, Name: "Superior Spellpower", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{stats.SpellPower: 30}, ItemType: proto.ItemType_ItemTypeWrist},
	{EffectID: 3850, ItemID: 44944, Name: "Major Stamina", Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{stats.Stamina: 40}, ItemType: proto.ItemType_ItemTypeWrist},
	{EffectID: 1119, SpellID: 44555, Name: "Exceptional Intellect", Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.Intellect: 16}, ItemType: proto.ItemType_ItemTypeWrist},
	{EffectID: 1147, SpellID: 44593, Name: "Major Spirit", Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.Spirit: 18}, ItemType: proto.ItemType_ItemTypeWrist},
	{EffectID: 3231, SpellID: 44598, Name: "Expertise", Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.Expertise: 15}, ItemType: proto.ItemType_ItemTypeWrist},
	{EffectID: 2661, SpellID: 44616, Name: "Greater Stats", Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.Stamina: 6, stats.Strength: 6, stats.Agility: 6, stats.Intellect: 6, stats.Spirit: 6}, ItemType: proto.ItemType_ItemTypeWrist},
	{EffectID: 2326, SpellID: 44635, Name: "Greater Spellpower", Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.SpellPower: 23}, ItemType: proto.ItemType_ItemTypeWrist},
	{EffectID: 1600, SpellID: 60616, Name: "Striking", Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.AttackPower: 38, stats.RangedAttackPower: 38}, ItemType: proto.ItemType_ItemTypeWrist},
	{EffectID: 3756, SpellID: 57683, Name: "Fur Lining - Attack Power", Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.AttackPower: 130, stats.RangedAttackPower: 130}, ItemType: proto.ItemType_ItemTypeWrist, RequiredProfession: proto.Profession_Leatherworking},
	{EffectID: 3757, SpellID: 57690, Name: "Fur Lining - Stamina", Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.Stamina: 102}, ItemType: proto.ItemType_ItemTypeWrist, RequiredProfession: proto.Profession_Leatherworking},
	{EffectID: 3758, SpellID: 57691, Name: "Fur Lining - Spell Power", Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.SpellPower: 76}, ItemType: proto.ItemType_ItemTypeWrist, RequiredProfession: proto.Profession_Leatherworking},
	{EffectID: 3759, SpellID: 57692, Name: "Fur Lining - Fire Resist", Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.FireResistance: 70}, ItemType: proto.ItemType_ItemTypeWrist, RequiredProfession: proto.Profession_Leatherworking},
	{EffectID: 3760, SpellID: 57694, Name: "Fur Lining - Frost Resist", Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.FrostResistance: 70}, ItemType: proto.ItemType_ItemTypeWrist, RequiredProfession: proto.Profession_Leatherworking},
	{EffectID: 3761, SpellID: 57696, Name: "Fur Lining - Shadow Resist", Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.ShadowResistance: 70}, ItemType: proto.ItemType_ItemTypeWrist, RequiredProfession: proto.Profession_Leatherworking},
	{EffectID: 3762, SpellID: 57699, Name: "Fur Lining - Nature Resist", Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.NatureResistance: 70}, ItemType: proto.ItemType_ItemTypeWrist, RequiredProfession: proto.Profession_Leatherworking},
	{EffectID: 3763, SpellID: 57701, Name: "Fur Lining - Arcane Resist", Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.ArcaneResistance: 70}, ItemType: proto.ItemType_ItemTypeWrist, RequiredProfession: proto.Profession_Leatherworking},

	// Hands
	{EffectID: 3253, ItemID: 44485, Name: "Armsman", Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{stats.Parry: 10}, ItemType: proto.ItemType_ItemTypeHands},
	{EffectID: 1603, SpellID: 60668, Name: "Crusher", Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.AttackPower: 44, stats.RangedAttackPower: 44}, ItemType: proto.ItemType_ItemTypeHands},
	{EffectID: 3246, SpellID: 44592, Name: "Exceptional Spellpower", Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.SpellPower: 28}, ItemType: proto.ItemType_ItemTypeHands},
	{EffectID: 3231, SpellID: 44484, Name: "Expertise", Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.Expertise: 15}, ItemType: proto.ItemType_ItemTypeHands},
	{EffectID: 3238, SpellID: 44506, Name: "Gatherer", Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{}, ItemType: proto.ItemType_ItemTypeHands},
	{EffectID: 3829, SpellID: 44513, Name: "Greater Assult", Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.AttackPower: 35, stats.RangedAttackPower: 35}, ItemType: proto.ItemType_ItemTypeHands},
	{EffectID: 3222, SpellID: 44529, Name: "Major Agility", Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.Agility: 20}, ItemType: proto.ItemType_ItemTypeHands},
	{EffectID: 3234, SpellID: 44488, Name: "Precision", Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.MeleeHit: 20, stats.SpellHit: 20}, ItemType: proto.ItemType_ItemTypeHands},
	{EffectID: 3603, SpellID: 54998, Name: "Hand-Mounted Pyro Rocket", Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{}, ItemType: proto.ItemType_ItemTypeHands, RequiredProfession: proto.Profession_Engineering},
	{EffectID: 3604, SpellID: 54999, Name: "Hyperspeed Accelerators", Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{}, ItemType: proto.ItemType_ItemTypeHands, RequiredProfession: proto.Profession_Engineering},
	{EffectID: 3860, SpellID: 63770, Name: "Reticulated Armor Webbing", Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.Armor: 885}, ItemType: proto.ItemType_ItemTypeHands, RequiredProfession: proto.Profession_Engineering},

	// Waist
	{EffectID: 3599, SpellID: 54736, Name: "Personal Electromagnetic Pulse Generator", Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{}, ItemType: proto.ItemType_ItemTypeWaist, RequiredProfession: proto.Profession_Engineering},
	{EffectID: 3601, SpellID: 54793, Name: "Frag Belt", Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{}, ItemType: proto.ItemType_ItemTypeWaist, RequiredProfession: proto.Profession_Engineering},

	// Legs
	{EffectID: 3325, ItemID: 38371, Name: "Jormungar Leg Armor", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{stats.Stamina: 45, stats.Agility: 15}, ItemType: proto.ItemType_ItemTypeLegs},
	{EffectID: 3326, ItemID: 38372, Name: "Nerubian Leg Armor", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{stats.AttackPower: 55, stats.RangedAttackPower: 55, stats.MeleeCrit: 15, stats.SpellCrit: 15}, ItemType: proto.ItemType_ItemTypeLegs},
	{EffectID: 3822, ItemID: 38373, Name: "Frosthide Leg Armor", Quality: proto.ItemQuality_ItemQualityEpic, Bonus: stats.Stats{stats.Stamina: 55, stats.Agility: 22}, ItemType: proto.ItemType_ItemTypeLegs},
	{EffectID: 3823, ItemID: 38374, Name: "Icescale Leg Armor", Quality: proto.ItemQuality_ItemQualityEpic, Bonus: stats.Stats{stats.AttackPower: 75, stats.RangedAttackPower: 75, stats.MeleeCrit: 22, stats.SpellCrit: 22}, ItemType: proto.ItemType_ItemTypeLegs},
	{EffectID: 3853, ItemID: 44963, Name: "Earthen Leg Armor", Quality: proto.ItemQuality_ItemQualityEpic, Bonus: stats.Stats{stats.Stamina: 28, stats.Resilience: 40}, ItemType: proto.ItemType_ItemTypeLegs},
	{EffectID: 3718, ItemID: 41601, Name: "Shining Spellthread", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{stats.Spirit: 12, stats.SpellPower: 35}, ItemType: proto.ItemType_ItemTypeLegs},
	{EffectID: 3719, ItemID: 41602, Name: "Brilliant Spellthread", Quality: proto.ItemQuality_ItemQualityEpic, Bonus: stats.Stats{stats.Spirit: 20, stats.SpellPower: 50}, ItemType: proto.ItemType_ItemTypeLegs},
	{EffectID: 3720, ItemID: 41603, Name: "Azure Spellthread", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{stats.Stamina: 20, stats.SpellPower: 35}, ItemType: proto.ItemType_ItemTypeLegs},
	{EffectID: 3721, ItemID: 41604, Name: "Sapphire Spellthread", Quality: proto.ItemQuality_ItemQualityEpic, Bonus: stats.Stats{stats.Stamina: 30, stats.SpellPower: 50}, ItemType: proto.ItemType_ItemTypeLegs},
	{EffectID: 3327, SpellID: 60583, Name: "Jormungar Leg Reinforcements", Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.Stamina: 55, stats.Agility: 22}, ItemType: proto.ItemType_ItemTypeLegs, RequiredProfession: proto.Profession_Leatherworking},
	{EffectID: 3328, SpellID: 60584, Name: "Nerubian Leg Reinforcements", Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.AttackPower: 75, stats.RangedAttackPower: 75, stats.MeleeCrit: 22, stats.SpellCrit: 22}, ItemType: proto.ItemType_ItemTypeLegs, RequiredProfession: proto.Profession_Leatherworking},
	{EffectID: 3873, SpellID: 56034, Name: "Master's Spellthread", Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.Stamina: 30, stats.SpellPower: 50}, ItemType: proto.ItemType_ItemTypeLegs, RequiredProfession: proto.Profession_Tailoring},
	{EffectID: 3872, SpellID: 56039, Name: "Sanctified Spellthread", Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.Spirit: 20, stats.SpellPower: 50}, ItemType: proto.ItemType_ItemTypeLegs, RequiredProfession: proto.Profession_Tailoring},

	// Feet
	{EffectID: 1597, ItemID: 44490, Name: "Greater Assault", Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{stats.AttackPower: 32, stats.RangedAttackPower: 32}, ItemType: proto.ItemType_ItemTypeFeet},
	{EffectID: 3232, ItemID: 44491, Name: "Tuskarr's Vitality", Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{stats.Stamina: 15}, ItemType: proto.ItemType_ItemTypeFeet},
	{EffectID: 3824, SpellID: 60606, Name: "Assault", Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.AttackPower: 24, stats.RangedAttackPower: 24}, ItemType: proto.ItemType_ItemTypeFeet},
	{EffectID: 1075, SpellID: 44528, Name: "Greater Fortitude", Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.Stamina: 22}, ItemType: proto.ItemType_ItemTypeFeet},
	{EffectID: 1147, SpellID: 44508, Name: "Greater Spirit", Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.Spirit: 18}, ItemType: proto.ItemType_ItemTypeFeet},
	{EffectID: 3244, SpellID: 44584, Name: "Greater Vitality", Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.MP5: 7}, ItemType: proto.ItemType_ItemTypeFeet},
	{EffectID: 3826, SpellID: 60623, Name: "Icewalker", Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.MeleeHit: 12, stats.SpellHit: 12, stats.MeleeCrit: 12, stats.SpellCrit: 12}, ItemType: proto.ItemType_ItemTypeFeet},
	{EffectID: 983, SpellID: 44589, Name: "Superior Agility", Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.Agility: 16}, ItemType: proto.ItemType_ItemTypeFeet},
	{EffectID: 3606, SpellID: 55016, Name: "Nitro Boosts", Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.MeleeCrit: 24, stats.SpellCrit: 24}, ItemType: proto.ItemType_ItemTypeFeet, RequiredProfession: proto.Profession_Engineering},

	// Weapon
	{EffectID: 1103, SpellID: 44633, Name: "Exceptional Agility", Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.Agility: 26}, ItemType: proto.ItemType_ItemTypeWeapon},
	{EffectID: 3844, SpellID: 44510, Name: "Exceptional Spirit", Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.Spirit: 45}, ItemType: proto.ItemType_ItemTypeWeapon},
	{EffectID: 3251, ItemID: 37339, Name: "Giant Slayer", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{}, ItemType: proto.ItemType_ItemTypeWeapon},
	{EffectID: 3239, ItemID: 37344, Name: "Icebreaker", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{}, ItemType: proto.ItemType_ItemTypeWeapon},
	{EffectID: 3731, ItemID: 41976, Name: "Titanium Weapon Chain", Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{stats.MeleeHit: 28, stats.SpellHit: 28}, ItemType: proto.ItemType_ItemTypeWeapon},
	{EffectID: 3833, ItemID: 44486, Name: "Superior Potency", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{stats.AttackPower: 65, stats.RangedAttackPower: 65}, ItemType: proto.ItemType_ItemTypeWeapon},
	{EffectID: 3834, ItemID: 44487, Name: "Mighty Spellpower", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{stats.SpellPower: 63}, ItemType: proto.ItemType_ItemTypeWeapon},
	{EffectID: 3789, ItemID: 44492, Name: "Berserking", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{}, ItemType: proto.ItemType_ItemTypeWeapon},
	{EffectID: 3241, ItemID: 44494, Name: "Lifeward", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{}, ItemType: proto.ItemType_ItemTypeWeapon},
	{EffectID: 3790, ItemID: 44495, Name: "Black Magic", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{}, ItemType: proto.ItemType_ItemTypeWeapon},
	{EffectID: 3788, ItemID: 44496, Name: "Accuracy", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{stats.MeleeHit: 25, stats.SpellHit: 25, stats.MeleeCrit: 25, stats.SpellCrit: 25}, ItemType: proto.ItemType_ItemTypeWeapon},
	{EffectID: 3830, SpellID: 44629, Name: "Exceptional Spellpower", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{stats.SpellPower: 50}, ItemType: proto.ItemType_ItemTypeWeapon},
	{EffectID: 1606, SpellID: 60621, Name: "Greater Potency", Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.AttackPower: 50, stats.RangedAttackPower: 50}, ItemType: proto.ItemType_ItemTypeWeapon},
	{EffectID: 3370, SpellID: 53343, Name: "Rune of Razorice", Phase: 1, Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{}, ItemType: proto.ItemType_ItemTypeWeapon, ClassAllowlist: []proto.Class{proto.Class_ClassDeathknight}},
	{EffectID: 3369, SpellID: 53341, Name: "Rune of Cinderglacier", Phase: 1, Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{}, ItemType: proto.ItemType_ItemTypeWeapon, ClassAllowlist: []proto.Class{proto.Class_ClassDeathknight}},
	{EffectID: 3366, SpellID: 53331, Name: "Rune of Lichbane", Phase: 1, Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{}, ItemType: proto.ItemType_ItemTypeWeapon, ClassAllowlist: []proto.Class{proto.Class_ClassDeathknight}},
	{EffectID: 3595, SpellID: 54447, Name: "Rune of Spellbreaking", Phase: 1, Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{}, ItemType: proto.ItemType_ItemTypeWeapon, ClassAllowlist: []proto.Class{proto.Class_ClassDeathknight}},
	{EffectID: 3594, SpellID: 54446, Name: "Rune of Swordbreaking", Phase: 1, Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{}, ItemType: proto.ItemType_ItemTypeWeapon, ClassAllowlist: []proto.Class{proto.Class_ClassDeathknight}},
	{EffectID: 3368, SpellID: 53344, Name: "Rune of the Fallen Crusader", Phase: 1, Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{}, ItemType: proto.ItemType_ItemTypeWeapon, ClassAllowlist: []proto.Class{proto.Class_ClassDeathknight}},
	{EffectID: 3883, SpellID: 70164, Name: "Rune of the Nerubian Carapace", Phase: 1, Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{}, ItemType: proto.ItemType_ItemTypeWeapon, ClassAllowlist: []proto.Class{proto.Class_ClassDeathknight}},

	// 2H Weapon
	{EffectID: 3247, ItemID: 44473, Name: "Scourgebane", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{}, ItemType: proto.ItemType_ItemTypeWeapon, EnchantType: proto.EnchantType_EnchantTypeTwoHand},
	{EffectID: 3827, ItemID: 44483, Name: "Massacre", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{stats.AttackPower: 110, stats.RangedAttackPower: 110}, ItemType: proto.ItemType_ItemTypeWeapon, EnchantType: proto.EnchantType_EnchantTypeTwoHand},
	{EffectID: 3828, SpellID: 44630, Name: "Greater Savagery", Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.AttackPower: 85, stats.RangedAttackPower: 85}, ItemType: proto.ItemType_ItemTypeWeapon, EnchantType: proto.EnchantType_EnchantTypeTwoHand},
	{EffectID: 3854, ItemID: 45059, Name: "Staff - Greater Spellpower", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{stats.SpellPower: 81}, ItemType: proto.ItemType_ItemTypeWeapon, EnchantType: proto.EnchantType_EnchantTypeStaff},
	{EffectID: 3367, SpellID: 53342, Name: "Rune of Spellshattering", Phase: 1, Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{}, ItemType: proto.ItemType_ItemTypeWeapon, ClassAllowlist: []proto.Class{proto.Class_ClassDeathknight}},
	{EffectID: 3365, SpellID: 53323, Name: "Rune of Swordshattering", Phase: 1, Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{}, ItemType: proto.ItemType_ItemTypeWeapon, ClassAllowlist: []proto.Class{proto.Class_ClassDeathknight}},
	{EffectID: 3847, SpellID: 62158, Name: "Rune of the Stoneskin Gargoyle", Phase: 1, Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{}, ItemType: proto.ItemType_ItemTypeWeapon, ClassAllowlist: []proto.Class{proto.Class_ClassDeathknight}},

	// Shield
	{EffectID: 1952, SpellID: 44489, Name: "Defense", Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.Defense: 20}, ItemType: proto.ItemType_ItemTypeWeapon, EnchantType: proto.EnchantType_EnchantTypeShield},
	{EffectID: 1128, SpellID: 60653, Name: "Greater Intellect", Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.Intellect: 25}, ItemType: proto.ItemType_ItemTypeWeapon, EnchantType: proto.EnchantType_EnchantTypeShield},
	{EffectID: 3748, ItemID: 42500, Name: "Titanium Shield Spike", Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{}, ItemType: proto.ItemType_ItemTypeWeapon, EnchantType: proto.EnchantType_EnchantTypeShield},
	{EffectID: 3849, ItemID: 44936, Name: "Titanium Plating", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{stats.BlockValue: 81}, ItemType: proto.ItemType_ItemTypeWeapon, EnchantType: proto.EnchantType_EnchantTypeShield},

	// Ring
	{EffectID: 3839, SpellID: 44645, Name: "Assault", Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.AttackPower: 40, stats.RangedAttackPower: 40}, ItemType: proto.ItemType_ItemTypeFinger, RequiredProfession: proto.Profession_Enchanting},
	{EffectID: 3840, SpellID: 44636, Name: "Greater Spellpower", Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.SpellPower: 23}, ItemType: proto.ItemType_ItemTypeFinger, RequiredProfession: proto.Profession_Enchanting},
	{EffectID: 3791, SpellID: 59636, Name: "Stamina", Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.Stamina: 30}, ItemType: proto.ItemType_ItemTypeFinger, RequiredProfession: proto.Profession_Enchanting},

	// Ranged
	{EffectID: 3607, ItemID: 41146, Name: "Sun Scope", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{}, ItemType: proto.ItemType_ItemTypeRanged},
	{EffectID: 3608, ItemID: 41167, Name: "Heartseeker Scope", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{}, ItemType: proto.ItemType_ItemTypeRanged},
	{EffectID: 3843, ItemID: 44739, Name: "Diamond-cut Refractor Scope", Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{}, ItemType: proto.ItemType_ItemTypeRanged},

	/////////////////////////////
	//   TBC
	/////////////////////////////

	// Head
	{EffectID: 2999, ItemID: 29186, Name: "Arcanum of the Defender", Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{stats.Defense: 16, stats.Dodge: 17}, ItemType: proto.ItemType_ItemTypeHead},
	{EffectID: 3002, ItemID: 29191, Name: "Arcanum of Power", Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{stats.SpellPower: 22, stats.SpellHit: 14}, ItemType: proto.ItemType_ItemTypeHead},
	{EffectID: 3003, ItemID: 29192, Name: "Arcanum of Ferocity", Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{stats.AttackPower: 34, stats.RangedAttackPower: 34, stats.MeleeHit: 16, stats.SpellHit: 16}, ItemType: proto.ItemType_ItemTypeHead},
	{EffectID: 3096, ItemID: 30846, Name: "Arcanum of the Outcast", Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{stats.Strength: 17, stats.Intellect: 16}, ItemType: proto.ItemType_ItemTypeHead},
	{EffectID: 3004, ItemID: 29193, Name: "Arcanum of the Gladiator", Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{stats.Stamina: 18, stats.Resilience: 20}, ItemType: proto.ItemType_ItemTypeHead},

	// ZG Head Enchants
	{EffectID: 2583, ItemID: 19782, Name: "Presence of Might", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{stats.Stamina: 10, stats.Defense: 10, stats.BlockValue: 15}, ItemType: proto.ItemType_ItemTypeHead, ClassAllowlist: []proto.Class{proto.Class_ClassWarrior}},

	// Shoulder
	{EffectID: 2982, ItemID: 28886, Name: "Greater Inscription of Discipline", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{stats.SpellPower: 18, stats.SpellCrit: 10}, ItemType: proto.ItemType_ItemTypeShoulder},
	{EffectID: 2986, ItemID: 28888, Name: "Greater Inscription of Vengeance", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{stats.AttackPower: 30, stats.RangedAttackPower: 30, stats.MeleeCrit: 10}, ItemType: proto.ItemType_ItemTypeShoulder},
	{EffectID: 2978, ItemID: 28889, Name: "Greater Inscription of Warding", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{stats.Defense: 10, stats.Dodge: 15}, ItemType: proto.ItemType_ItemTypeShoulder},
	{EffectID: 2995, ItemID: 28909, Name: "Greater Inscription of the Orb", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{stats.SpellPower: 12, stats.SpellCrit: 15}, ItemType: proto.ItemType_ItemTypeShoulder},
	{EffectID: 2997, ItemID: 28910, Name: "Greater Inscription of the Blade", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{stats.AttackPower: 20, stats.RangedAttackPower: 20, stats.MeleeCrit: 15}, ItemType: proto.ItemType_ItemTypeShoulder},
	{EffectID: 2991, ItemID: 28911, Name: "Greater Inscription of the Knight", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{stats.Defense: 15, stats.Dodge: 10}, ItemType: proto.ItemType_ItemTypeShoulder},
	{EffectID: 2605, ItemID: 20076, Name: "Zandalar Signet of Mojo", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{stats.SpellPower: 18}, ItemType: proto.ItemType_ItemTypeShoulder},
	{EffectID: 2721, ItemID: 23545, Name: "Power of the Scourge", Quality: proto.ItemQuality_ItemQualityEpic, Bonus: stats.Stats{stats.SpellPower: 15, stats.SpellCrit: 14}, ItemType: proto.ItemType_ItemTypeShoulder},
	{EffectID: 2717, ItemID: 23548, Name: "Might of the Scourge", Quality: proto.ItemQuality_ItemQualityEpic, Bonus: stats.Stats{stats.AttackPower: 26, stats.RangedAttackPower: 26, stats.MeleeCrit: 14}, ItemType: proto.ItemType_ItemTypeShoulder},
	{EffectID: 2716, ItemID: 23549, Name: "Fortitude of the Scourge", Quality: proto.ItemQuality_ItemQualityEpic, Bonus: stats.Stats{stats.Stamina: 16, stats.Armor: 100}, ItemType: proto.ItemType_ItemTypeShoulder},

	// Back
	{EffectID: 2622, ItemID: 33148, Name: "Enchant Cloak - Dodge", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{stats.Dodge: 12}, ItemType: proto.ItemType_ItemTypeBack},
	{EffectID: 2621, ItemID: 33150, Name: "Enchant Cloak - Subtlety", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{}, ItemType: proto.ItemType_ItemTypeBack},
	{EffectID: 849, ItemID: 11206, Name: "Enchant Cloak - Lesser Agility", Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{stats.Agility: 3}, ItemType: proto.ItemType_ItemTypeBack},
	{EffectID: 368, SpellID: 34004, Name: "Enchant Cloak - Greater Agility", Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.Agility: 12}, ItemType: proto.ItemType_ItemTypeBack},
	{EffectID: 2938, ItemID: 28274, Name: "Enchant Cloak - Spell Penetration", Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{stats.SpellPenetration: 20}, ItemType: proto.ItemType_ItemTypeBack},
	{EffectID: 1441, ItemID: 28277, Name: "Enchant Cloak - Greater Shadow Resistance", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{stats.ShadowResistance: 15}, ItemType: proto.ItemType_ItemTypeBack},
	{EffectID: 2648, ItemID: 35756, Name: "Enchant Cloak - Steelweave", Phase: 5, Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{stats.Defense: 12}, ItemType: proto.ItemType_ItemTypeBack},

	// Chest
	{EffectID: 2659, SpellID: 27957, Name: "Chest - Exceptional Health", Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.Health: 150}, ItemType: proto.ItemType_ItemTypeChest},
	{EffectID: 2660, ItemID: 22546, Name: "Chest - Exceptional Mana", Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.Mana: 150}, ItemType: proto.ItemType_ItemTypeChest},
	{EffectID: 2661, ItemID: 24003, Name: "Chest - Exceptional Stats", Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.Stamina: 6, stats.Intellect: 6, stats.Spirit: 6, stats.Strength: 6, stats.Agility: 6}, ItemType: proto.ItemType_ItemTypeChest},
	{EffectID: 2933, ItemID: 28270, Name: "Chest - Major Resilience", Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{stats.Resilience: 15}, ItemType: proto.ItemType_ItemTypeChest},
	{EffectID: 1144, SpellID: 33990, Name: "Chest - Major Spirit", Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.Spirit: 15}, ItemType: proto.ItemType_ItemTypeChest},
	{EffectID: 3150, SpellID: 33991, Name: "Chest - Restore Mana Prime", Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.MP5: 6}, ItemType: proto.ItemType_ItemTypeChest},
	{EffectID: 1950, ItemID: 35500, Name: "Chest - Defense", Phase: 5, Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.Defense: 15}, ItemType: proto.ItemType_ItemTypeChest},

	// Wrist
	{EffectID: 2649, ItemID: 22533, Name: "Bracer - Fortitude", Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{stats.Stamina: 12}, ItemType: proto.ItemType_ItemTypeWrist},
	{EffectID: 2650, ItemID: 22534, Name: "Bracer - Spellpower", Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{stats.SpellPower: 15}, ItemType: proto.ItemType_ItemTypeWrist},
	{EffectID: 369, SpellID: 34001, Name: "Bracer - Major Intellect", Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{stats.Intellect: 12}, ItemType: proto.ItemType_ItemTypeWrist},
	{EffectID: 2647, SpellID: 27899, Name: "Bracer - Brawn", Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{stats.Strength: 12}, ItemType: proto.ItemType_ItemTypeWrist},
	{EffectID: 1593, SpellID: 34002, Name: "Bracer - Assault", Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{stats.AttackPower: 24, stats.RangedAttackPower: 24}, ItemType: proto.ItemType_ItemTypeWrist},
	{EffectID: 1891, SpellID: 27905, Name: "Bracer - Stats", Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{stats.Stamina: 4, stats.Intellect: 4, stats.Spirit: 4, stats.Strength: 4, stats.Agility: 4}, ItemType: proto.ItemType_ItemTypeWrist},

	// Hands
	{EffectID: 2935, ItemID: 28271, Name: "Gloves - Spell Strike", Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{stats.SpellHit: 15}, ItemType: proto.ItemType_ItemTypeHands},
	{EffectID: 2937, ItemID: 28272, Name: "Gloves - Major Spellpower", Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{stats.SpellPower: 20}, ItemType: proto.ItemType_ItemTypeHands},
	{EffectID: 684, SpellID: 33995, Name: "Gloves - Major Strength", Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{stats.Strength: 15}, ItemType: proto.ItemType_ItemTypeHands},
	{EffectID: 2564, ItemID: 33152, Name: "Gloves - Major Agility", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{stats.Agility: 15}, ItemType: proto.ItemType_ItemTypeHands},
	{EffectID: 2613, ItemID: 33153, Name: "Gloves - Threat", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{}, ItemType: proto.ItemType_ItemTypeHands},

	// Legs
	{EffectID: 2748, ItemID: 24274, Name: "Runic Spellthread", Quality: proto.ItemQuality_ItemQualityEpic, Bonus: stats.Stats{stats.SpellPower: 35, stats.Stamina: 20}, ItemType: proto.ItemType_ItemTypeLegs},
	{EffectID: 2747, ItemID: 24273, Name: "Mystic Spellthread", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{stats.SpellPower: 25, stats.Stamina: 15}, ItemType: proto.ItemType_ItemTypeLegs},
	{EffectID: 3010, ItemID: 29533, Name: "Cobrahide Leg Armor", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{stats.AttackPower: 40, stats.RangedAttackPower: 40, stats.MeleeCrit: 10}, ItemType: proto.ItemType_ItemTypeLegs},
	{EffectID: 3012, ItemID: 29535, Name: "Nethercobra Leg Armor", Quality: proto.ItemQuality_ItemQualityEpic, Bonus: stats.Stats{stats.AttackPower: 50, stats.RangedAttackPower: 50, stats.MeleeCrit: 12}, ItemType: proto.ItemType_ItemTypeLegs},
	{EffectID: 3013, ItemID: 29536, Name: "Nethercleft Leg Armor", Quality: proto.ItemQuality_ItemQualityEpic, Bonus: stats.Stats{stats.Stamina: 40, stats.Agility: 12}, ItemType: proto.ItemType_ItemTypeLegs},

	// Feet
	{EffectID: 851, ItemID: 16220, Name: "Enchant Boots - Spirit", Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{stats.Spirit: 5}, ItemType: proto.ItemType_ItemTypeFeet},
	{EffectID: 2940, ItemID: 35297, Name: "Enchant Boots - Boar's Speed", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{stats.Stamina: 9}, ItemType: proto.ItemType_ItemTypeFeet},
	{EffectID: 2656, ItemID: 35298, Name: "Enchant Boots - Vitality", Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{stats.MP5: 4}, ItemType: proto.ItemType_ItemTypeFeet},
	{EffectID: 2649, ItemID: 22543, Name: "Enchant Boots - Fortitude", Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{stats.Stamina: 12}, ItemType: proto.ItemType_ItemTypeFeet},
	{EffectID: 2657, ItemID: 22544, Name: "Enchant Boots - Dexterity", Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{stats.Agility: 12}, ItemType: proto.ItemType_ItemTypeFeet},
	{EffectID: 2939, ItemID: 28279, Name: "Enchant Boots - Cat's Swiftness", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{stats.Agility: 6}, ItemType: proto.ItemType_ItemTypeFeet},
	{EffectID: 2658, ItemID: 22545, Name: "Enchant Boots - Surefooted", Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{stats.MeleeHit: 10}, ItemType: proto.ItemType_ItemTypeFeet},

	// Weapon
	{EffectID: 1897, ItemID: 16250, Name: "Superior Striking", Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{}, ItemType: proto.ItemType_ItemTypeWeapon},
	{EffectID: 963, ItemID: 22552, Name: "Major Striking", Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{}, ItemType: proto.ItemType_ItemTypeWeapon},
	{EffectID: 1900, ItemID: 16252, Name: "Crusader", Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{}, ItemType: proto.ItemType_ItemTypeWeapon},
	{EffectID: 2666, ItemID: 22551, Name: "Major Intellect", Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{stats.Intellect: 30}, ItemType: proto.ItemType_ItemTypeWeapon},
	{EffectID: 2667, ItemID: 22554, Name: "Savagery", Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{stats.AttackPower: 70, stats.RangedAttackPower: 70}, ItemType: proto.ItemType_ItemTypeWeapon, EnchantType: proto.EnchantType_EnchantTypeTwoHand},
	{EffectID: 2669, ItemID: 22555, Name: "Major Spellpower", Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{stats.SpellPower: 40}, ItemType: proto.ItemType_ItemTypeWeapon},
	{EffectID: 2671, ItemID: 22560, Name: "Sunfire", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{}, ItemType: proto.ItemType_ItemTypeWeapon},
	{EffectID: 2672, ItemID: 22561, Name: "Soulfrost", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{}, ItemType: proto.ItemType_ItemTypeWeapon},
	{EffectID: 2673, ItemID: 22559, Name: "Mongoose", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{}, ItemType: proto.ItemType_ItemTypeWeapon},
	{EffectID: 2564, ItemID: 19445, Name: "Agility", Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.Agility: 15}, ItemType: proto.ItemType_ItemTypeWeapon},
	{EffectID: 3222, ItemID: 33165, Name: "Greater Agility", Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.Agility: 20}, ItemType: proto.ItemType_ItemTypeWeapon},
	{EffectID: 2670, ItemID: 22556, Name: "2H Weapon - Major Agility", Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{stats.Agility: 35}, ItemType: proto.ItemType_ItemTypeWeapon, EnchantType: proto.EnchantType_EnchantTypeTwoHand},
	{EffectID: 3225, ItemID: 33307, Name: "Executioner", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{}, ItemType: proto.ItemType_ItemTypeWeapon},
	{EffectID: 3273, ItemID: 35498, Name: "Deathfrost", Phase: 5, Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{}, ItemType: proto.ItemType_ItemTypeWeapon},
	{EffectID: 3855, ItemID: 45060, Name: "Staff - Spellpower", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{stats.SpellPower: 69}, ItemType: proto.ItemType_ItemTypeWeapon, EnchantType: proto.EnchantType_EnchantTypeStaff},

	// Shield
	{EffectID: 2654, ItemID: 22539, Name: "Intellect", Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.Intellect: 12}, ItemType: proto.ItemType_ItemTypeWeapon, EnchantType: proto.EnchantType_EnchantTypeShield},
	{EffectID: 1071, ItemID: 28282, Name: "Major Stamina", Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.Stamina: 18}, ItemType: proto.ItemType_ItemTypeWeapon, EnchantType: proto.EnchantType_EnchantTypeShield},
	{EffectID: 3229, SpellID: 44383, Name: "Resilience", Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.Resilience: 12}, ItemType: proto.ItemType_ItemTypeWeapon, EnchantType: proto.EnchantType_EnchantTypeShield},

	// Ring
	{EffectID: 2929, ItemID: 22535, Name: "Striking", Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{}, ItemType: proto.ItemType_ItemTypeFinger, RequiredProfession: proto.Profession_Enchanting},
	{EffectID: 2928, ItemID: 22536, Name: "Spellpower", Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.SpellPower: 12}, ItemType: proto.ItemType_ItemTypeFinger, RequiredProfession: proto.Profession_Enchanting},
	{EffectID: 2931, ItemID: 22538, Name: "Stats", Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.Stamina: 4, stats.Intellect: 4, stats.Spirit: 4, stats.Strength: 4, stats.Agility: 4}, ItemType: proto.ItemType_ItemTypeFinger, RequiredProfession: proto.Profession_Enchanting},

	// Ranged
	{EffectID: 2523, ItemID: 18283, Name: "Biznicks 247x128 Accurascope", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{}, ItemType: proto.ItemType_ItemTypeRanged},
	{EffectID: 2723, ItemID: 23765, Name: "Khorium Scope", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{}, ItemType: proto.ItemType_ItemTypeRanged},
	{EffectID: 2724, ItemID: 23766, Name: "Stabilized Eternium Scope", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{}, ItemType: proto.ItemType_ItemTypeRanged},
}
