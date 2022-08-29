package items

import (
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

var Enchants = []Enchant{
	// Multi-slot
	{ID: 38375, EffectID: 3329, Name: "Borean Armor Kit", Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.Stamina: 12}, ItemType: proto.ItemType_ItemTypeHead, EnchantType: proto.EnchantType_EnchantTypeKit},
	{ID: 38376, EffectID: 3330, Name: "Heavy Borean Armor Kit", Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{stats.Stamina: 18}, ItemType: proto.ItemType_ItemTypeHead, EnchantType: proto.EnchantType_EnchantTypeKit},

	// Head
	{ID: 44069, EffectID: 3795, Name: "Arcanum of Triumph", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{stats.AttackPower: 50, stats.RangedAttackPower: 50, stats.Resilience: 20}, ItemType: proto.ItemType_ItemTypeHead},
	{ID: 44075, EffectID: 3796, Name: "Arcanum of Dominance", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{stats.SpellPower: 29, stats.HealingPower: 29, stats.Resilience: 20}, ItemType: proto.ItemType_ItemTypeHead},
	{ID: 44875, EffectID: 3842, Name: "Arcanum of the Savage Gladiator", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{stats.Stamina: 30, stats.Resilience: 25}, ItemType: proto.ItemType_ItemTypeHead},
	{ID: 44137, EffectID: 3812, Name: "Arcanum of the Frosty Soul", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{stats.Stamina: 30, stats.FrostResistance: 25}, ItemType: proto.ItemType_ItemTypeHead},
	{ID: 44138, EffectID: 3813, Name: "Arcanum of Toxic Warding", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{stats.Stamina: 30, stats.NatureResistance: 25}, ItemType: proto.ItemType_ItemTypeHead},
	{ID: 44139, EffectID: 3814, Name: "Arcanum of the Fleeing Shadow", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{stats.Stamina: 30, stats.ShadowResistance: 25}, ItemType: proto.ItemType_ItemTypeHead},
	{ID: 44140, EffectID: 3815, Name: "Arcanum of the Eclipsed Moon", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{stats.Stamina: 30, stats.ArcaneResistance: 25}, ItemType: proto.ItemType_ItemTypeHead},
	{ID: 44141, EffectID: 3816, Name: "Arcanum of the Flame's Soul", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{stats.Stamina: 30, stats.FireResistance: 25}, ItemType: proto.ItemType_ItemTypeHead},
	{ID: 44876, EffectID: 3819, Name: "Arcanum of Blissful Mending", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{stats.SpellPower: 30, stats.HealingPower: 30, stats.MP5: 10}, ItemType: proto.ItemType_ItemTypeHead},
	{ID: 44877, EffectID: 3820, Name: "Arcanum of Burning Mysteries", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{stats.SpellPower: 30, stats.HealingPower: 30, stats.MeleeCrit: 20, stats.SpellCrit: 20}, ItemType: proto.ItemType_ItemTypeHead},
	{ID: 44878, EffectID: 3818, Name: "Arcanum of the Stalwart Protector", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{stats.Stamina: 37, stats.Defense: 20}, ItemType: proto.ItemType_ItemTypeHead},
	{ID: 44879, EffectID: 3817, Name: "Arcanum of Torment", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{stats.AttackPower: 50, stats.RangedAttackPower: 50, stats.MeleeCrit: 20, stats.SpellCrit: 20}, ItemType: proto.ItemType_ItemTypeHead},
	{ID: 67839, EffectID: 3878, Name: "Mind Amplification Dish", IsSpellID: true, Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.Stamina: 45}, ItemType: proto.ItemType_ItemTypeHead, RequiredProfession: proto.Profession_Engineering},

	// Shoulder
	{ID: 44067, EffectID: 3793, Name: "Inscription of Triumph", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{stats.AttackPower: 40, stats.RangedAttackPower: 40, stats.Resilience: 15}, ItemType: proto.ItemType_ItemTypeShoulder},
	{ID: 44068, EffectID: 3794, Name: "Inscription of Dominance", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{stats.SpellPower: 23, stats.HealingPower: 23, stats.Resilience: 15}, ItemType: proto.ItemType_ItemTypeShoulder},
	{ID: 44957, EffectID: 3852, Name: "Greater Inscription of the Gladiator", Quality: proto.ItemQuality_ItemQualityEpic, Bonus: stats.Stats{stats.Stamina: 30, stats.Resilience: 15}, ItemType: proto.ItemType_ItemTypeShoulder},
	{ID: 44129, EffectID: 3806, Name: "Lesser Inscription of the Storm", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{stats.SpellPower: 18, stats.HealingPower: 18, stats.MeleeCrit: 10, stats.SpellCrit: 10}, ItemType: proto.ItemType_ItemTypeShoulder},
	{ID: 44130, EffectID: 3807, Name: "Lesser Inscription of the Crag", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{stats.SpellPower: 18, stats.HealingPower: 18, stats.MP5: 5}, ItemType: proto.ItemType_ItemTypeShoulder},
	{ID: 44131, EffectID: 3875, Name: "Lesser Inscription of the Axe", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{stats.AttackPower: 30, stats.RangedAttackPower: 30, stats.MeleeCrit: 10, stats.SpellCrit: 10}, ItemType: proto.ItemType_ItemTypeShoulder},
	{ID: 44132, EffectID: 3876, Name: "Lesser Inscription of the Pinnacle", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{stats.Dodge: 15, stats.Defense: 10}, ItemType: proto.ItemType_ItemTypeShoulder},
	{ID: 44874, EffectID: 3810, Name: "Greater Inscription of the Storm", Quality: proto.ItemQuality_ItemQualityEpic, Bonus: stats.Stats{stats.SpellPower: 24, stats.HealingPower: 24, stats.MeleeCrit: 15, stats.SpellCrit: 15}, ItemType: proto.ItemType_ItemTypeShoulder},
	{ID: 44872, EffectID: 3809, Name: "Greater Inscription of the Crag", Quality: proto.ItemQuality_ItemQualityEpic, Bonus: stats.Stats{stats.SpellPower: 24, stats.HealingPower: 24, stats.MP5: 8}, ItemType: proto.ItemType_ItemTypeShoulder},
	{ID: 44871, EffectID: 3808, Name: "Greater Inscription of the Axe", Quality: proto.ItemQuality_ItemQualityEpic, Bonus: stats.Stats{stats.AttackPower: 40, stats.RangedAttackPower: 40, stats.MeleeCrit: 15, stats.SpellCrit: 15}, ItemType: proto.ItemType_ItemTypeShoulder},
	{ID: 44873, EffectID: 3811, Name: "Greater Inscription of the Pinnacle", Quality: proto.ItemQuality_ItemQualityEpic, Bonus: stats.Stats{stats.Dodge: 20, stats.Defense: 15}, ItemType: proto.ItemType_ItemTypeShoulder},
	{ID: 61120, EffectID: 3838, Name: "Master's Inscription of the Storm", IsSpellID: true, Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.SpellPower: 70, stats.HealingPower: 70, stats.MeleeCrit: 15, stats.SpellCrit: 15}, ItemType: proto.ItemType_ItemTypeShoulder, RequiredProfession: proto.Profession_Inscription},
	{ID: 61118, EffectID: 3836, Name: "Master's Inscription of the Crag", IsSpellID: true, Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.SpellPower: 70, stats.HealingPower: 70, stats.MP5: 8}, ItemType: proto.ItemType_ItemTypeShoulder, RequiredProfession: proto.Profession_Inscription},
	{ID: 61117, EffectID: 3835, Name: "Master's Inscription of the Axe", IsSpellID: true, Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.AttackPower: 120, stats.RangedAttackPower: 120, stats.MeleeCrit: 15, stats.SpellCrit: 15}, ItemType: proto.ItemType_ItemTypeShoulder, RequiredProfession: proto.Profession_Inscription},
	{ID: 61119, EffectID: 3837, Name: "Master's Inscription of the Pinnacle", IsSpellID: true, Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.Dodge: 60, stats.Defense: 15}, ItemType: proto.ItemType_ItemTypeShoulder, RequiredProfession: proto.Profession_Inscription},

	// Back
	{ID: 37330, EffectID: 1262, Name: "Superior Arcane Resistance", Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{stats.ArcaneResistance: 20}, ItemType: proto.ItemType_ItemTypeBack},
	{ID: 37331, EffectID: 1354, Name: "Superior Fire Resistance", Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{stats.FireResistance: 20}, ItemType: proto.ItemType_ItemTypeBack},
	{ID: 44483, EffectID: 3230, Name: "Superior Frost Resistance", IsSpellID: true, Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.FrostResistance: 20}, ItemType: proto.ItemType_ItemTypeBack},
	{ID: 44494, EffectID: 1400, Name: "Superior Nature Resistance", IsSpellID: true, Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.NatureResistance: 20}, ItemType: proto.ItemType_ItemTypeBack},
	{ID: 44590, EffectID: 1446, Name: "Superior Shadow Resistance", IsSpellID: true, Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.ShadowResistance: 20}, ItemType: proto.ItemType_ItemTypeBack},
	{ID: 37347, EffectID: 1951, Name: "Titanweave", Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{stats.Defense: 16}, ItemType: proto.ItemType_ItemTypeBack},
	{ID: 37349, EffectID: 3256, Name: "Shadow Armor", Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{stats.Agility: 10}, ItemType: proto.ItemType_ItemTypeBack},
	{ID: 44471, EffectID: 3294, Name: "Mighty Armor", Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{stats.Armor: 225}, ItemType: proto.ItemType_ItemTypeBack},
	{ID: 44472, EffectID: 3831, Name: "Greater Speed", Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{stats.MeleeHaste: 23, stats.SpellHaste: 23}, ItemType: proto.ItemType_ItemTypeBack},
	{ID: 44488, EffectID: 3296, Name: "Wisdom", Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{stats.Spirit: 10}, ItemType: proto.ItemType_ItemTypeBack},
	{ID: 44582, EffectID: 3243, Name: "Spell Piercing", IsSpellID: true, Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.SpellPenetration: 35}, ItemType: proto.ItemType_ItemTypeBack},
	{ID: 60609, EffectID: 3825, Name: "Speed", IsSpellID: true, Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.MeleeHaste: 15, stats.SpellHaste: 15}, ItemType: proto.ItemType_ItemTypeBack},
	{ID: 44500, EffectID: 983, Name: "Superior Agility", IsSpellID: true, Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.Agility: 16}, ItemType: proto.ItemType_ItemTypeBack},
	{ID: 60663, EffectID: 1099, Name: "Major Agility", IsSpellID: true, Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.Agility: 22}, ItemType: proto.ItemType_ItemTypeBack},
	{ID: 55002, EffectID: 3605, Name: "Flexweave Underlay", IsSpellID: true, Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.Agility: 23}, ItemType: proto.ItemType_ItemTypeBack, RequiredProfession: proto.Profession_Engineering},
	{ID: 55642, EffectID: 3722, Name: "Lightweave Embroidery", IsSpellID: true, Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{}, ItemType: proto.ItemType_ItemTypeBack, RequiredProfession: proto.Profession_Tailoring},
	{ID: 55769, EffectID: 3728, Name: "Darkglow Embroidery", IsSpellID: true, Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{}, ItemType: proto.ItemType_ItemTypeBack, RequiredProfession: proto.Profession_Tailoring},
	{ID: 55777, EffectID: 3730, Name: "Swordguard Embroidery", IsSpellID: true, Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{}, ItemType: proto.ItemType_ItemTypeBack, RequiredProfession: proto.Profession_Tailoring},
	{ID: 63765, EffectID: 3859, Name: "Springy Arachnoweave", IsSpellID: true, Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.SpellPower: 27}, ItemType: proto.ItemType_ItemTypeBack, RequiredProfession: proto.Profession_Engineering},

	// Chest
	{ID: 37340, EffectID: 3245, Name: "Exceptional Resilience", Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{stats.Resilience: 20}, ItemType: proto.ItemType_ItemTypeChest},
	{ID: 44623, EffectID: 3252, Name: "Super Stats", IsSpellID: true, Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{stats.Stamina: 8, stats.Strength: 8, stats.Agility: 8, stats.Intellect: 8, stats.Spirit: 8}, ItemType: proto.ItemType_ItemTypeChest},
	{ID: 44489, EffectID: 3832, Name: "Powerful Stats", Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{stats.Stamina: 10, stats.Strength: 10, stats.Agility: 10, stats.Intellect: 10, stats.Spirit: 10}, ItemType: proto.ItemType_ItemTypeChest},
	{ID: 27958, EffectID: 3233, Name: "Exceptional Mana", IsSpellID: true, Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.Mana: 250}, ItemType: proto.ItemType_ItemTypeChest},
	{ID: 44492, EffectID: 3236, Name: "Mighty Health", IsSpellID: true, Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.Health: 200}, ItemType: proto.ItemType_ItemTypeChest},
	{ID: 47900, EffectID: 3297, Name: "Super Health", IsSpellID: true, Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.Health: 275}, ItemType: proto.ItemType_ItemTypeChest},
	{ID: 44509, EffectID: 2381, Name: "Greater Mana Restoration", IsSpellID: true, Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.MP5: 10}, ItemType: proto.ItemType_ItemTypeChest},
	{ID: 47766, EffectID: 1953, Name: "Greater Defense", IsSpellID: true, Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.Defense: 22}, ItemType: proto.ItemType_ItemTypeChest},

	// Wrist
	{ID: 44484, EffectID: 3845, Name: "Greater Assault", Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{stats.AttackPower: 50, stats.RangedAttackPower: 50}, ItemType: proto.ItemType_ItemTypeWrist},
	{ID: 44498, EffectID: 2332, Name: "Superior Spellpower", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{stats.SpellPower: 30, stats.HealingPower: 30}, ItemType: proto.ItemType_ItemTypeWrist},
	{ID: 44944, EffectID: 3850, Name: "Major Stamina", Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{stats.Stamina: 40}, ItemType: proto.ItemType_ItemTypeWrist},
	{ID: 44555, EffectID: 1119, Name: "Exceptional Intellect", IsSpellID: true, Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.Intellect: 16}, ItemType: proto.ItemType_ItemTypeWrist},
	{ID: 44593, EffectID: 1147, Name: "Major Spirit", IsSpellID: true, Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.Spirit: 18}, ItemType: proto.ItemType_ItemTypeWrist},
	{ID: 44598, EffectID: 3231, Name: "Expertise", IsSpellID: true, Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.Expertise: 15}, ItemType: proto.ItemType_ItemTypeWrist},
	{ID: 44616, EffectID: 2661, Name: "Greater Stats", IsSpellID: true, Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.Stamina: 6, stats.Strength: 6, stats.Agility: 6, stats.Intellect: 6, stats.Spirit: 6}, ItemType: proto.ItemType_ItemTypeWrist},
	{ID: 44635, EffectID: 2326, Name: "Greater Spellpower", IsSpellID: true, Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.SpellPower: 23, stats.HealingPower: 23}, ItemType: proto.ItemType_ItemTypeWrist},
	{ID: 60616, EffectID: 1600, Name: "Striking", IsSpellID: true, Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.AttackPower: 38, stats.RangedAttackPower: 38}, ItemType: proto.ItemType_ItemTypeWrist},
	{ID: 57683, EffectID: 3756, Name: "Fur Lining - Attack Power", IsSpellID: true, Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.AttackPower: 130, stats.RangedAttackPower: 130}, ItemType: proto.ItemType_ItemTypeWrist, RequiredProfession: proto.Profession_Leatherworking},
	{ID: 57690, EffectID: 3757, Name: "Fur Lining - Stamina", IsSpellID: true, Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.Stamina: 102}, ItemType: proto.ItemType_ItemTypeWrist, RequiredProfession: proto.Profession_Leatherworking},
	{ID: 57691, EffectID: 3758, Name: "Fur Lining - Spell Power", IsSpellID: true, Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.SpellPower: 76, stats.HealingPower: 76}, ItemType: proto.ItemType_ItemTypeWrist, RequiredProfession: proto.Profession_Leatherworking},
	{ID: 57692, EffectID: 3759, Name: "Fur Lining - Fire Resist", IsSpellID: true, Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.FireResistance: 70}, ItemType: proto.ItemType_ItemTypeWrist, RequiredProfession: proto.Profession_Leatherworking},
	{ID: 57694, EffectID: 3760, Name: "Fur Lining - Frost Resist", IsSpellID: true, Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.FrostResistance: 70}, ItemType: proto.ItemType_ItemTypeWrist, RequiredProfession: proto.Profession_Leatherworking},
	{ID: 57696, EffectID: 3761, Name: "Fur Lining - Shadow Resist", IsSpellID: true, Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.ShadowResistance: 70}, ItemType: proto.ItemType_ItemTypeWrist, RequiredProfession: proto.Profession_Leatherworking},
	{ID: 57699, EffectID: 3762, Name: "Fur Lining - Nature Resist", IsSpellID: true, Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.NatureResistance: 70}, ItemType: proto.ItemType_ItemTypeWrist, RequiredProfession: proto.Profession_Leatherworking},
	{ID: 57701, EffectID: 3763, Name: "Fur Lining - Arcane Resist", IsSpellID: true, Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.ArcaneResistance: 70}, ItemType: proto.ItemType_ItemTypeWrist, RequiredProfession: proto.Profession_Leatherworking},

	// Hands
	{ID: 44485, EffectID: 3253, Name: "Armsman", Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{stats.Parry: 10}, ItemType: proto.ItemType_ItemTypeHands},
	{ID: 60668, EffectID: 1603, Name: "Crusher", IsSpellID: true, Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.AttackPower: 44, stats.RangedAttackPower: 44}, ItemType: proto.ItemType_ItemTypeHands},
	{ID: 44592, EffectID: 3246, Name: "Exceptional Spellpower", IsSpellID: true, Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.SpellPower: 28, stats.HealingPower: 28}, ItemType: proto.ItemType_ItemTypeHands},
	{ID: 44484, EffectID: 3231, Name: "Expertise", IsSpellID: true, Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.Expertise: 15}, ItemType: proto.ItemType_ItemTypeHands},
	{ID: 44506, EffectID: 3238, Name: "Gatherer", IsSpellID: true, Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{}, ItemType: proto.ItemType_ItemTypeHands},
	{ID: 44513, EffectID: 3829, Name: "Greater Assult", IsSpellID: true, Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.AttackPower: 35, stats.RangedAttackPower: 35}, ItemType: proto.ItemType_ItemTypeHands},
	{ID: 44529, EffectID: 3222, Name: "Major Agility", IsSpellID: true, Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.Agility: 20}, ItemType: proto.ItemType_ItemTypeHands},
	{ID: 44488, EffectID: 3234, Name: "Precision", IsSpellID: true, Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.MeleeHit: 20, stats.SpellHit: 20}, ItemType: proto.ItemType_ItemTypeHands},
	{ID: 54998, EffectID: 3603, Name: "Hand-Mounted Pyro Rocket", IsSpellID: true, Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{}, ItemType: proto.ItemType_ItemTypeHands, RequiredProfession: proto.Profession_Engineering},
	{ID: 54999, EffectID: 3604, Name: "Hyperspeed Accelerators", IsSpellID: true, Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{}, ItemType: proto.ItemType_ItemTypeHands, RequiredProfession: proto.Profession_Engineering},
	{ID: 63770, EffectID: 3860, Name: "Reticulated Armor Webbing", IsSpellID: true, Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.Armor: 885}, ItemType: proto.ItemType_ItemTypeHands, RequiredProfession: proto.Profession_Engineering},

	// Waist
	{ID: 54736, EffectID: 3599, Name: "Personal Electromagnetic Pulse Generator", IsSpellID: true, Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{}, ItemType: proto.ItemType_ItemTypeWaist, RequiredProfession: proto.Profession_Engineering},
	{ID: 54793, EffectID: 3601, Name: "Frag Belt", IsSpellID: true, Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{}, ItemType: proto.ItemType_ItemTypeWaist, RequiredProfession: proto.Profession_Engineering},

	// Legs
	{ID: 38371, EffectID: 3325, Name: "Jormungar Leg Armor", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{stats.Stamina: 45, stats.Agility: 15}, ItemType: proto.ItemType_ItemTypeLegs},
	{ID: 38372, EffectID: 3326, Name: "Nerubian Leg Armor", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{stats.AttackPower: 55, stats.RangedAttackPower: 55, stats.MeleeCrit: 15, stats.SpellCrit: 15}, ItemType: proto.ItemType_ItemTypeLegs},
	{ID: 38373, EffectID: 3822, Name: "Frosthide Leg Armor", Quality: proto.ItemQuality_ItemQualityEpic, Bonus: stats.Stats{stats.Stamina: 55, stats.Agility: 22}, ItemType: proto.ItemType_ItemTypeLegs},
	{ID: 38374, EffectID: 3823, Name: "Icescale Leg Armor", Quality: proto.ItemQuality_ItemQualityEpic, Bonus: stats.Stats{stats.AttackPower: 75, stats.RangedAttackPower: 75, stats.MeleeCrit: 22, stats.SpellCrit: 22}, ItemType: proto.ItemType_ItemTypeLegs},
	{ID: 44963, EffectID: 3853, Name: "Earthen Leg Armor", Quality: proto.ItemQuality_ItemQualityEpic, Bonus: stats.Stats{stats.Stamina: 28, stats.Resilience: 40}, ItemType: proto.ItemType_ItemTypeLegs},
	{ID: 41601, EffectID: 3718, Name: "Shining Spellthread", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{stats.Spirit: 12, stats.SpellPower: 35, stats.HealingPower: 35}, ItemType: proto.ItemType_ItemTypeLegs},
	{ID: 41602, EffectID: 3719, Name: "Brilliant Spellthread", Quality: proto.ItemQuality_ItemQualityEpic, Bonus: stats.Stats{stats.Spirit: 20, stats.SpellPower: 50, stats.HealingPower: 50}, ItemType: proto.ItemType_ItemTypeLegs},
	{ID: 41603, EffectID: 3720, Name: "Azure Spellthread", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{stats.Stamina: 20, stats.SpellPower: 35, stats.HealingPower: 35}, ItemType: proto.ItemType_ItemTypeLegs},
	{ID: 41604, EffectID: 3721, Name: "Sapphire Spellthread", Quality: proto.ItemQuality_ItemQualityEpic, Bonus: stats.Stats{stats.Stamina: 30, stats.SpellPower: 50, stats.HealingPower: 50}, ItemType: proto.ItemType_ItemTypeLegs},
	{ID: 60583, EffectID: 3327, Name: "Jormungar Leg Reinforcements", IsSpellID: true, Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.Stamina: 55, stats.Agility: 22}, ItemType: proto.ItemType_ItemTypeLegs, RequiredProfession: proto.Profession_Leatherworking},
	{ID: 60584, EffectID: 3328, Name: "Nerubian Leg Reinforcements", IsSpellID: true, Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.AttackPower: 75, stats.RangedAttackPower: 75, stats.MeleeCrit: 22, stats.SpellCrit: 22}, ItemType: proto.ItemType_ItemTypeLegs, RequiredProfession: proto.Profession_Leatherworking},
	{ID: 56034, EffectID: 3873, Name: "Master's Spellthread", IsSpellID: true, Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.Stamina: 30, stats.SpellPower: 50, stats.HealingPower: 50}, ItemType: proto.ItemType_ItemTypeLegs, RequiredProfession: proto.Profession_Tailoring},
	{ID: 56039, EffectID: 3872, Name: "Sanctified Spellthread", IsSpellID: true, Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.Spirit: 20, stats.SpellPower: 50, stats.HealingPower: 50}, ItemType: proto.ItemType_ItemTypeLegs, RequiredProfession: proto.Profession_Tailoring},

	// Feet
	{ID: 44490, EffectID: 1597, Name: "Greater Assault", Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{stats.AttackPower: 32, stats.RangedAttackPower: 32}, ItemType: proto.ItemType_ItemTypeFeet},
	{ID: 44491, EffectID: 3232, Name: "Tuskarr's Vitality", Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{stats.Stamina: 15}, ItemType: proto.ItemType_ItemTypeFeet},
	{ID: 60606, EffectID: 3824, Name: "Assault", IsSpellID: true, Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.AttackPower: 24, stats.RangedAttackPower: 24}, ItemType: proto.ItemType_ItemTypeFeet},
	{ID: 44528, EffectID: 1075, Name: "Greater Fortitude", IsSpellID: true, Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.Stamina: 22}, ItemType: proto.ItemType_ItemTypeFeet},
	{ID: 44508, EffectID: 1147, Name: "Greater Spirit", IsSpellID: true, Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.Spirit: 18}, ItemType: proto.ItemType_ItemTypeFeet},
	{ID: 44584, EffectID: 3244, Name: "Greater Vitality", IsSpellID: true, Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.MP5: 7}, ItemType: proto.ItemType_ItemTypeFeet},
	{ID: 60623, EffectID: 3826, Name: "Icewalker", IsSpellID: true, Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.MeleeHit: 12, stats.SpellHit: 12, stats.MeleeCrit: 12, stats.SpellCrit: 12}, ItemType: proto.ItemType_ItemTypeFeet},
	{ID: 44589, EffectID: 983, Name: "Superior Agility", IsSpellID: true, Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.Agility: 16}, ItemType: proto.ItemType_ItemTypeFeet},
	{ID: 55016, EffectID: 3606, Name: "Nitro Boosts", IsSpellID: true, Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.MeleeCrit: 24, stats.SpellCrit: 24}, ItemType: proto.ItemType_ItemTypeFeet, RequiredProfession: proto.Profession_Engineering},

	// Weapon
	{ID: 44633, EffectID: 1103, Name: "Exceptional Agility", IsSpellID: true, Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.Agility: 26}, ItemType: proto.ItemType_ItemTypeWeapon},
	{ID: 44510, EffectID: 3844, Name: "Exceptional Spirit", IsSpellID: true, Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.Spirit: 45}, ItemType: proto.ItemType_ItemTypeWeapon},
	{ID: 37339, EffectID: 3251, Name: "Giant Slayer", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{}, ItemType: proto.ItemType_ItemTypeWeapon},
	{ID: 37344, EffectID: 3239, Name: "Icebreaker", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{}, ItemType: proto.ItemType_ItemTypeWeapon},
	{ID: 41976, EffectID: 3731, Name: "Titanium Weapon Chain", Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{stats.MeleeHit: 28, stats.SpellHit: 28}, ItemType: proto.ItemType_ItemTypeWeapon},
	{ID: 44486, EffectID: 3833, Name: "Superior Potency", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{stats.AttackPower: 65, stats.RangedAttackPower: 65}, ItemType: proto.ItemType_ItemTypeWeapon},
	{ID: 44487, EffectID: 3834, Name: "Mighty Spellpower", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{stats.SpellPower: 63, stats.HealingPower: 63}, ItemType: proto.ItemType_ItemTypeWeapon},
	{ID: 44492, EffectID: 3789, Name: "Berserking", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{}, ItemType: proto.ItemType_ItemTypeWeapon},
	{ID: 44494, EffectID: 3241, Name: "Lifeward", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{}, ItemType: proto.ItemType_ItemTypeWeapon},
	{ID: 44495, EffectID: 3790, Name: "Black Magic", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{}, ItemType: proto.ItemType_ItemTypeWeapon},
	{ID: 44496, EffectID: 3788, Name: "Accuracy", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{stats.MeleeHit: 25, stats.SpellHit: 25, stats.MeleeCrit: 25, stats.SpellCrit: 25}, ItemType: proto.ItemType_ItemTypeWeapon},
	{ID: 44629, EffectID: 3830, Name: "Exceptional Spellpower", IsSpellID: true, Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{stats.SpellPower: 50, stats.HealingPower: 50}, ItemType: proto.ItemType_ItemTypeWeapon},
	{ID: 60621, EffectID: 1606, Name: "Greater Potency", IsSpellID: true, Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.AttackPower: 50, stats.RangedAttackPower: 50}, ItemType: proto.ItemType_ItemTypeWeapon},
	{ID: 53343, EffectID: 3370, Name: "Rune of Razorice", IsSpellID: true, Phase: 1, Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{}, ItemType: proto.ItemType_ItemTypeWeapon, ClassAllowlist: []proto.Class{proto.Class_ClassDeathknight}},
	{ID: 53341, EffectID: 3369, Name: "Rune of Cinderglacier", IsSpellID: true, Phase: 1, Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{}, ItemType: proto.ItemType_ItemTypeWeapon, ClassAllowlist: []proto.Class{proto.Class_ClassDeathknight}},
	{ID: 53331, EffectID: 3366, Name: "Rune of Lichbane", IsSpellID: true, Phase: 1, Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{}, ItemType: proto.ItemType_ItemTypeWeapon, ClassAllowlist: []proto.Class{proto.Class_ClassDeathknight}},
	{ID: 54447, EffectID: 3595, Name: "Rune of Spellbreaking", IsSpellID: true, Phase: 1, Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{}, ItemType: proto.ItemType_ItemTypeWeapon, ClassAllowlist: []proto.Class{proto.Class_ClassDeathknight}},
	{ID: 54446, EffectID: 3594, Name: "Rune of Swordbreaking", IsSpellID: true, Phase: 1, Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{}, ItemType: proto.ItemType_ItemTypeWeapon, ClassAllowlist: []proto.Class{proto.Class_ClassDeathknight}},
	{ID: 53344, EffectID: 3368, Name: "Rune of the Fallen Crusader", IsSpellID: true, Phase: 1, Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{}, ItemType: proto.ItemType_ItemTypeWeapon, ClassAllowlist: []proto.Class{proto.Class_ClassDeathknight}},
	{ID: 70164, EffectID: 3883, Name: "Rune of the Nerubian Carapace", IsSpellID: true, Phase: 1, Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{}, ItemType: proto.ItemType_ItemTypeWeapon, ClassAllowlist: []proto.Class{proto.Class_ClassDeathknight}},

	// 2H Weapon
	{ID: 44473, EffectID: 3247, Name: "Scourgebane", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{}, ItemType: proto.ItemType_ItemTypeWeapon, EnchantType: proto.EnchantType_EnchantTypeTwoHand},
	{ID: 44483, EffectID: 3827, Name: "Massacre", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{stats.AttackPower: 110, stats.RangedAttackPower: 110}, ItemType: proto.ItemType_ItemTypeWeapon, EnchantType: proto.EnchantType_EnchantTypeTwoHand},
	{ID: 44630, EffectID: 3828, Name: "Greater Savagery", IsSpellID: true, Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.AttackPower: 85, stats.RangedAttackPower: 85}, ItemType: proto.ItemType_ItemTypeWeapon, EnchantType: proto.EnchantType_EnchantTypeTwoHand},
	{ID: 45059, EffectID: 3854, Name: "Greater Spellpower", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{stats.SpellPower: 81, stats.HealingPower: 81}, ItemType: proto.ItemType_ItemTypeWeapon, EnchantType: proto.EnchantType_EnchantTypeTwoHand},
	{ID: 53342, EffectID: 3367, Name: "Rune of Spellshattering", IsSpellID: true, Phase: 1, Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{}, ItemType: proto.ItemType_ItemTypeWeapon, ClassAllowlist: []proto.Class{proto.Class_ClassDeathknight}},
	{ID: 53323, EffectID: 3365, Name: "Rune of Swordshattering", IsSpellID: true, Phase: 1, Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{}, ItemType: proto.ItemType_ItemTypeWeapon, ClassAllowlist: []proto.Class{proto.Class_ClassDeathknight}},
	{ID: 62158, EffectID: 3847, Name: "Rune of the Stoneskin Gargoyle", IsSpellID: true, Phase: 1, Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{}, ItemType: proto.ItemType_ItemTypeWeapon, ClassAllowlist: []proto.Class{proto.Class_ClassDeathknight}},

	// Shield
	{ID: 44489, EffectID: 1952, Name: "Defense", IsSpellID: true, Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.Defense: 20}, ItemType: proto.ItemType_ItemTypeWeapon, EnchantType: proto.EnchantType_EnchantTypeShield},
	{ID: 60653, EffectID: 1128, Name: "Greater Intellect", IsSpellID: true, Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.Intellect: 25}, ItemType: proto.ItemType_ItemTypeWeapon, EnchantType: proto.EnchantType_EnchantTypeShield},
	{ID: 42500, EffectID: 3748, Name: "Titanium Shield Spike", Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{}, ItemType: proto.ItemType_ItemTypeWeapon, EnchantType: proto.EnchantType_EnchantTypeShield},
	{ID: 44936, EffectID: 3849, Name: "Titanium Plating", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{stats.BlockValue: 81}, ItemType: proto.ItemType_ItemTypeWeapon, EnchantType: proto.EnchantType_EnchantTypeShield},

	// Ring
	{ID: 44645, EffectID: 3839, Name: "Assault", IsSpellID: true, Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.AttackPower: 40, stats.RangedAttackPower: 40}, ItemType: proto.ItemType_ItemTypeFinger, RequiredProfession: proto.Profession_Enchanting},
	{ID: 44636, EffectID: 3840, Name: "Greater Spellpower", IsSpellID: true, Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.SpellPower: 23, stats.HealingPower: 23}, ItemType: proto.ItemType_ItemTypeFinger, RequiredProfession: proto.Profession_Enchanting},
	{ID: 59636, EffectID: 3791, Name: "Stamina", IsSpellID: true, Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.Stamina: 30}, ItemType: proto.ItemType_ItemTypeFinger, RequiredProfession: proto.Profession_Enchanting},

	// Ranged
	{ID: 41146, EffectID: 3607, Name: "Sun Scope", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{}, ItemType: proto.ItemType_ItemTypeRanged},
	{ID: 41167, EffectID: 3608, Name: "Heartseeker Scope", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{}, ItemType: proto.ItemType_ItemTypeRanged},
	{ID: 44739, EffectID: 3843, Name: "Diamond-cut Refractor Scope", Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{}, ItemType: proto.ItemType_ItemTypeRanged},

	/////////////////////////////
	//   TBC
	/////////////////////////////

	// Head
	{ID: 29186, EffectID: 2999, Name: "Arcanum of the Defender", Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{stats.Defense: 16, stats.Dodge: 17}, ItemType: proto.ItemType_ItemTypeHead},
	{ID: 29191, EffectID: 3002, Name: "Arcanum of Power", Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{stats.SpellPower: 22, stats.SpellHit: 14}, ItemType: proto.ItemType_ItemTypeHead},
	{ID: 29192, EffectID: 3003, Name: "Arcanum of Ferocity", Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{stats.AttackPower: 34, stats.RangedAttackPower: 34, stats.MeleeHit: 16, stats.SpellHit: 16}, ItemType: proto.ItemType_ItemTypeHead},
	{ID: 30846, EffectID: 3096, Name: "Arcanum of the Outcast", Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{stats.Strength: 17, stats.Intellect: 16}, ItemType: proto.ItemType_ItemTypeHead},
	{ID: 29193, EffectID: 3004, Name: "Arcanum of the Gladiator", Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{stats.Stamina: 18, stats.Resilience: 20}, ItemType: proto.ItemType_ItemTypeHead},

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
	{ID: 16250, EffectID: 1897, Name: "Superior Striking", Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{}, ItemType: proto.ItemType_ItemTypeWeapon},
	{ID: 22552, EffectID: 963, Name: "Major Striking", Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{}, ItemType: proto.ItemType_ItemTypeWeapon},
	{ID: 16252, EffectID: 1900, Name: "Crusader", Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{}, ItemType: proto.ItemType_ItemTypeWeapon},
	{ID: 22551, EffectID: 2666, Name: "Major Intellect", Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{stats.Intellect: 30}, ItemType: proto.ItemType_ItemTypeWeapon},
	{ID: 22554, EffectID: 2667, Name: "Savagery", Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{stats.AttackPower: 70, stats.RangedAttackPower: 70}, ItemType: proto.ItemType_ItemTypeWeapon, EnchantType: proto.EnchantType_EnchantTypeTwoHand},
	{ID: 22555, EffectID: 2669, Name: "Major Spellpower", Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{stats.SpellPower: 40}, ItemType: proto.ItemType_ItemTypeWeapon},
	{ID: 22560, EffectID: 2671, Name: "Sunfire", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{stats.ArcaneSpellPower: 50, stats.FireSpellPower: 50}, ItemType: proto.ItemType_ItemTypeWeapon},
	{ID: 22561, EffectID: 2672, Name: "Soulfrost", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{stats.FrostSpellPower: 54, stats.ShadowSpellPower: 54}, ItemType: proto.ItemType_ItemTypeWeapon},
	{ID: 22559, EffectID: 2673, Name: "Mongoose", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{}, ItemType: proto.ItemType_ItemTypeWeapon},
	{ID: 19445, EffectID: 2564, Name: "Agility", Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.Agility: 15}, ItemType: proto.ItemType_ItemTypeWeapon},
	{ID: 33165, EffectID: 3222, Name: "Greater Agility", Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.Agility: 20}, ItemType: proto.ItemType_ItemTypeWeapon},
	{ID: 22556, EffectID: 2670, Name: "2H Weapon - Major Agility", Quality: proto.ItemQuality_ItemQualityUncommon, Bonus: stats.Stats{stats.Agility: 35}, ItemType: proto.ItemType_ItemTypeWeapon, EnchantType: proto.EnchantType_EnchantTypeTwoHand},
	{ID: 33307, EffectID: 3225, Name: "Executioner", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{}, ItemType: proto.ItemType_ItemTypeWeapon},
	{ID: 35498, EffectID: 3273, Name: "Deathfrost", Phase: 5, Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{}, ItemType: proto.ItemType_ItemTypeWeapon},

	// Shield
	{ID: 22539, EffectID: 2654, Name: "Intellect", Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.Intellect: 12}, ItemType: proto.ItemType_ItemTypeWeapon, EnchantType: proto.EnchantType_EnchantTypeShield},
	{ID: 28282, EffectID: 1071, Name: "Major Stamina", Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.Stamina: 18}, ItemType: proto.ItemType_ItemTypeWeapon, EnchantType: proto.EnchantType_EnchantTypeShield},
	{ID: 44383, EffectID: 3229, Name: "Resilience", IsSpellID: true, Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.Resilience: 12}, ItemType: proto.ItemType_ItemTypeWeapon, EnchantType: proto.EnchantType_EnchantTypeShield},

	// Ring
	{ID: 22535, EffectID: 2929, Name: "Striking", Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{}, ItemType: proto.ItemType_ItemTypeFinger, RequiredProfession: proto.Profession_Enchanting},
	{ID: 22536, EffectID: 2928, Name: "Spellpower", Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.SpellPower: 12}, ItemType: proto.ItemType_ItemTypeFinger, RequiredProfession: proto.Profession_Enchanting},
	{ID: 22538, EffectID: 2931, Name: "Stats", Quality: proto.ItemQuality_ItemQualityCommon, Bonus: stats.Stats{stats.Stamina: 4, stats.Intellect: 4, stats.Spirit: 4, stats.Strength: 4, stats.Agility: 4}, ItemType: proto.ItemType_ItemTypeFinger, RequiredProfession: proto.Profession_Enchanting},

	// Ranged
	{ID: 18283, EffectID: 2523, Name: "Biznicks 247x128 Accurascope", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{}, ItemType: proto.ItemType_ItemTypeRanged},
	{ID: 23765, EffectID: 2723, Name: "Khorium Scope", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{}, ItemType: proto.ItemType_ItemTypeRanged},
	{ID: 23766, EffectID: 2724, Name: "Stabilized Eternium Scope", Quality: proto.ItemQuality_ItemQualityRare, Bonus: stats.Stats{}, ItemType: proto.ItemType_ItemTypeRanged},
}
