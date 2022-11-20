package main

import (
	"regexp"

	"github.com/wowsims/wotlk/sim/core/proto"
)

// Allows manual overriding for Item fields in case WowHead is wrong.
var ItemOverrideOverrides = []ItemOverride{
	{ /** Destruction Holo-gogs */ ID: 32494, ClassAllowlist: []proto.Class{proto.Class_ClassMage, proto.Class_ClassPriest, proto.Class_ClassWarlock}},
	{ /** Gadgetstorm Goggles */ ID: 32476, ClassAllowlist: []proto.Class{proto.Class_ClassShaman}},
	{ /** Magnified Moon Specs */ ID: 32480, ClassAllowlist: []proto.Class{proto.Class_ClassDruid}},
	{ /** Quad Deathblow X44 Goggles */ ID: 34353, ClassAllowlist: []proto.Class{proto.Class_ClassDruid, proto.Class_ClassRogue}},
	{ /** Hyper-Magnified Moon Specs */ ID: 35182, ClassAllowlist: []proto.Class{proto.Class_ClassDruid}},
	{ /** Lightning Etched Specs */ ID: 34355, ClassAllowlist: []proto.Class{proto.Class_ClassShaman}},
	{ /** Annihilator Holo-Gogs */ ID: 34847, ClassAllowlist: []proto.Class{proto.Class_ClassMage, proto.Class_ClassPriest, proto.Class_ClassWarlock}},
	{ /** Twinblade of the Pheonix */ ID: 29993, Stats: Stats{proto.Stat_StatRangedAttackPower: 108}},
	{ /** Pillar of Ferocity */ ID: 30883, Stats: Stats{proto.Stat_StatArmor: 550}},
	{ /** Incisor Fragment */ ID: 37723, Stats: Stats{proto.Stat_StatArmorPenetration: 0}},

	// TBC blacksmithing crafted weapons
	{ID: 28431, HandType: proto.HandType_HandTypeOneHand},
	{ID: 28432, HandType: proto.HandType_HandTypeOneHand},
	{ID: 28433, HandType: proto.HandType_HandTypeOneHand},
	{ID: 28437, HandType: proto.HandType_HandTypeOneHand},
	{ID: 28438, HandType: proto.HandType_HandTypeOneHand},
	{ID: 28439, HandType: proto.HandType_HandTypeOneHand},
	{ID: 28657, HandType: proto.HandType_HandTypeOneHand},
	{ID: 28767, HandType: proto.HandType_HandTypeOneHand},

	{ID: 29994}, // Wildercloak
	{ID: 29387}, // Gnomeregan Auto-Blocker 600
	{ID: 38289}, // Coren's Lucky Coin
	{ID: 34473}, // Commendation of Kaelthas

	// WotLK Halloween Event items
	{ID: 49121, Phase: 1}, // Ring of Ghoulish Glee
	{ID: 49123, Phase: 1}, // The Horseman's Seal
	{ID: 49124, Phase: 1}, // Wicked Witch's Band
	{ID: 49128, Phase: 1}, // The Horseman's Baleful Blade
	{ID: 49126, Phase: 1}, // The Horseman's Horrific Helm

	// Valorous T8 Sets
	{ID: 45375, Phase: 2},
	{ID: 45381, Phase: 2},
	{ID: 45382, Phase: 2},
	{ID: 45376, Phase: 2},
	{ID: 45370, Phase: 2},
	{ID: 45371, Phase: 2},
	{ID: 45383, Phase: 2},
	{ID: 45372, Phase: 2},
	{ID: 45377, Phase: 2},
	{ID: 45384, Phase: 2},
	{ID: 45379, Phase: 2},
	{ID: 45385, Phase: 2},
	{ID: 45380, Phase: 2},
	{ID: 45373, Phase: 2},
	{ID: 45374, Phase: 2},
	{ID: 45391, Phase: 2},
	{ID: 45386, Phase: 2},
	{ID: 45340, Phase: 2},
	{ID: 45335, Phase: 2},
	{ID: 45336, Phase: 2},
	{ID: 45341, Phase: 2},
	{ID: 45337, Phase: 2},
	{ID: 45342, Phase: 2},
	{ID: 45338, Phase: 2},
	{ID: 45343, Phase: 2},
	{ID: 45339, Phase: 2},
	{ID: 45344, Phase: 2},
	{ID: 45419, Phase: 2},
	{ID: 45417, Phase: 2},
	{ID: 45420, Phase: 2},
	{ID: 45421, Phase: 2},
	{ID: 45422, Phase: 2},
	{ID: 45387, Phase: 2},
	{ID: 45392, Phase: 2},
	{ID: 46131, Phase: 2},
	{ID: 45365, Phase: 2},
	{ID: 45367, Phase: 2},
	{ID: 45369, Phase: 2},
	{ID: 45368, Phase: 2},
	{ID: 45388, Phase: 2},
	{ID: 45393, Phase: 2},
	{ID: 46313, Phase: 2},
	{ID: 45351, Phase: 2},
	{ID: 45355, Phase: 2},
	{ID: 45345, Phase: 2},
	{ID: 45356, Phase: 2},
	{ID: 45346, Phase: 2},
	{ID: 45347, Phase: 2},
	{ID: 45357, Phase: 2},
	{ID: 45352, Phase: 2},
	{ID: 45358, Phase: 2},
	{ID: 45348, Phase: 2},
	{ID: 45359, Phase: 2},
	{ID: 45349, Phase: 2},
	{ID: 45353, Phase: 2},
	{ID: 45354, Phase: 2},
	{ID: 45394, Phase: 2},
	{ID: 45395, Phase: 2},
	{ID: 45389, Phase: 2},
	{ID: 45360, Phase: 2},
	{ID: 45361, Phase: 2},
	{ID: 45362, Phase: 2},
	{ID: 45363, Phase: 2},
	{ID: 45364, Phase: 2},
	{ID: 45390, Phase: 2},
	{ID: 45429, Phase: 2},
	{ID: 45424, Phase: 2},
	{ID: 45430, Phase: 2},
	{ID: 45425, Phase: 2},
	{ID: 45426, Phase: 2},
	{ID: 45431, Phase: 2},
	{ID: 45427, Phase: 2},
	{ID: 45432, Phase: 2},
	{ID: 45428, Phase: 2},
	{ID: 45433, Phase: 2},
	{ID: 45396, Phase: 2},
	{ID: 45397, Phase: 2},
	{ID: 45398, Phase: 2},
	{ID: 45399, Phase: 2},
	{ID: 45400, Phase: 2},
	{ID: 45413, Phase: 2},
	{ID: 45412, Phase: 2},
	{ID: 45406, Phase: 2},
	{ID: 45414, Phase: 2},
	{ID: 45401, Phase: 2},
	{ID: 45411, Phase: 2},
	{ID: 45402, Phase: 2},
	{ID: 45408, Phase: 2},
	{ID: 45409, Phase: 2},
	{ID: 45403, Phase: 2},
	{ID: 45415, Phase: 2},
	{ID: 45410, Phase: 2},
	{ID: 45404, Phase: 2},
	{ID: 45405, Phase: 2},
	{ID: 45416, Phase: 2},

	// Other items Wowhead has the wrong phase listed for
	// Ick's loot table from Pit of Saron
	{ID: 49812, Phase: 4},
	{ID: 49808, Phase: 4},
	{ID: 49811, Phase: 4},
	{ID: 49807, Phase: 4},
	{ID: 49810, Phase: 4},
	{ID: 49809, Phase: 4},

	// Keep these items
	{ID: 37032, Keep: true}, // Edge of the Tuskarr
	{ID: 11815, Keep: true}, // Hand of Justice
	{ID: 12590, Keep: true}, // Felstriker
	{ID: 15808, Keep: true}, // Fine Light Crossbow (for hunter testing).
	{ID: 19019, Keep: true}, // Thunderfury
	{ID: 19808, Keep: true}, // Rockhide Strongfish
	{ID: 20837, Keep: true}, // Sunstrider Axe
	{ID: 20966, Keep: true}, // Jade Pendant of Blasting
	{ID: 24114, Keep: true}, // Braided Eternium Chain
	{ID: 28830, Keep: true}, // Dragonspine Trophy
	{ID: 29383, Keep: true}, // Bloodlust Brooch
	{ID: 29387, Keep: true}, // Gnomeregan Auto-Blocker 600
	{ID: 29994, Keep: true}, // Thalassian Wildercloak
	{ID: 29996, Keep: true}, // Rod of the Sun King
	{ID: 30032, Keep: true}, // Red Belt of Battle
	{ID: 30627, Keep: true}, // Tsunami Talisman
	{ID: 30720, Keep: true}, // Serpent-Coil Braid
	{ID: 32387, Keep: true}, // Idol of the Raven Goddess
	{ID: 32658, Keep: true}, // Badge of Tenacity
	{ID: 33135, Keep: true}, // Falling Star
	{ID: 33140, Keep: true}, // Blood of Amber
	{ID: 33143, Keep: true}, // Stone of Blades
	{ID: 33144, Keep: true}, // Facet of Eternity
	{ID: 33504, Keep: true}, // Libram of Divine Purpose
	{ID: 33506, Keep: true}, // Skycall Totem
	{ID: 33507, Keep: true}, // Stonebreaker's Totem
	{ID: 33510, Keep: true}, // Unseen moon idol
	{ID: 33829, Keep: true}, // Hex Shrunken Head
	{ID: 33831, Keep: true}, // Berserkers Call
	{ID: 34472, Keep: true}, // Shard of Contempt
	{ID: 34473, Keep: true}, // Commendation of Kael'thas
	{ID: 37574, Keep: true}, // Libram of Furious Blows
	{ID: 38072, Keep: true}, // Thunder Capacitor
	{ID: 38212, Keep: true}, // Death Knight's Anguish
	{ID: 38287, Keep: true}, // Empty Mug of Direbrew
	{ID: 38289, Keep: true}, // Coren's Lucky Coin
	{ID: 39208, Keep: true}, // Sigil of the Dark Rider

	{ID: 6360, Keep: true}, // Steelscale Crushfish
	{ID: 8345, Keep: true}, // Wolfshead Helm
	{ID: 9449, Keep: true}, // Manual Crowd Pummeler

	{ID: 27510, Keep: true}, // Tidefury Gauntlets
	{ID: 27802, Keep: true}, // Tidefury Shoulderguards
	{ID: 27909, Keep: true}, // Tidefury Kilt
	{ID: 28231, Keep: true}, // Tidefury Chestpiece
	{ID: 28349, Keep: true}, // Tidefury Helm

	{ID: 15056, Keep: true}, // Stormshroud Armor
	{ID: 15057, Keep: true}, // Stormshroud Pants
	{ID: 15058, Keep: true}, // Stormshroud Shoulders
	{ID: 21278, Keep: true}, // Stormshroud Gloves

	// Grand Marshal / High Warlord Weapons
	{ID: 18843, Keep: true},
	{ID: 18844, Keep: true},
	{ID: 18847, Keep: true},
	{ID: 18848, Keep: true},

	{ID: 41752, Keep: true}, // Brunnhildar Axe

	// Undead Slaying Sets
	// Plate
	{ID: 43068, Keep: true},
	{ID: 43069, Keep: true},
	{ID: 43070, Keep: true},
	{ID: 43071, Keep: true},
	// Cloth
	{ID: 43072, Keep: true},
	{ID: 43073, Keep: true},
	{ID: 43074, Keep: true},
	{ID: 43075, Keep: true},
	// Mail
	{ID: 43076, Keep: true},
	{ID: 43077, Keep: true},
	{ID: 43078, Keep: true},
	{ID: 43079, Keep: true},
	//Leather
	{ID: 43080, Keep: true},
	{ID: 43081, Keep: true},
	{ID: 43082, Keep: true},
	{ID: 43083, Keep: true},

	// Filter out these items
	{ID: 17782, Filter: true}, // talisman of the binding shard
	{ID: 17783, Filter: true}, // talisman of the binding fragment
	{ID: 17802, Filter: true}, // Deprecated version of Thunderfury
	{ID: 18582, Filter: true},
	{ID: 18583, Filter: true},
	{ID: 18584, Filter: true},
	{ID: 24265, Filter: true},
	{ID: 32384, Filter: true},
	{ID: 32421, Filter: true},
	{ID: 32422, Filter: true},
	{ID: 33482, Filter: true},
	{ID: 34576, Filter: true}, // Battlemaster's Cruelty
	{ID: 34577, Filter: true}, // Battlemaster's Depreavity
	{ID: 34578, Filter: true}, // Battlemaster's Determination
	{ID: 34579, Filter: true}, // Battlemaster's Audacity
	{ID: 34580, Filter: true}, // Battlemaster's Perseverence
	{ID: 43727, Filter: true}, // Bonescythe Breastplate
}

// Item icons to include in the DB, so they don't need to be separately loaded in the UI.
var extraItemIcons = []int{
	// Pet foods
	33874,
	43005,

	// Spellstones
	41174,
	41196,

	// Demonic Rune
	12662,

	// Food IDs
	27655,
	27657,
	27658,
	27664,
	33052,
	33825,
	33872,
	34753,
	34754,
	34756,
	34758,
	34767,
	34769,
	42994,
	42995,
	42996,
	42998,
	42999,
	43000,
	43015,

	// Flask IDs
	13512,
	22851,
	22853,
	22854,
	22861,
	22866,
	33208,
	40079,
	44939,
	46376,
	46377,
	46378,
	46379,

	// Elixer IDs
	40072,
	40078,
	40097,
	40109,
	44328,
	44332,

	// Elixer IDs
	13452,
	13454,
	22824,
	22827,
	22831,
	22833,
	22834,
	22835,
	22840,
	28103,
	28104,
	31679,
	32062,
	32067,
	32068,
	39666,
	40068,
	40070,
	40073,
	40076,
	44325,
	44327,
	44329,
	44330,
	44331,
	9088,
	9224,

	// Potions / In Battle Consumes
	13442,
	20520,
	22105,
	22788,
	22828,
	22832,
	22837,
	22838,
	22839,
	22849,
	31677,
	33447,
	33448,
	40093,
	40211,
	40212,
	40536,
	40771,
	41119,
	42641,

	// Poisons
	43231,
	43233,
	43235,

	// Thistle Tea
	7676,
}

// If any of these match the item name, don't include it.
var denyListNameRegexes = []*regexp.Regexp{
	regexp.MustCompile(`DB\d`),
	regexp.MustCompile(`DEPRECATED`),
	regexp.MustCompile(`PH`),
	regexp.MustCompile(`TEST`),
	regexp.MustCompile(`Test`),
	regexp.MustCompile(`Bracer 3`),
	regexp.MustCompile(`Bracer 2`),
	regexp.MustCompile(`Bracer 1`),
	regexp.MustCompile(`Boots 3`),
	regexp.MustCompile(`Boots 2`),
	regexp.MustCompile(`Boots 1`),
	regexp.MustCompile(`zOLD`),
	regexp.MustCompile(`30 Epic`),
	regexp.MustCompile(`Indalamar`),
	regexp.MustCompile(`QR XXXX`),
	regexp.MustCompile(`Deprecated: Keanna`),
	regexp.MustCompile(`90 Epic`),
	regexp.MustCompile(`66 Epic`),
	regexp.MustCompile(`63 Blue`),
	regexp.MustCompile(`90 Green`),
	regexp.MustCompile(`63 Green`),
}

// Allows manual overriding for Gem fields in case WowHead is wrong.
var GemOverrideOverrides = []GemOverride{
	{ID: 33131, Stats: Stats{proto.Stat_StatAttackPower: 32, proto.Stat_StatRangedAttackPower: 32}},

	// pvp non-unique gems not in game currently.
	{ID: 35489, Filter: true},
	{ID: 38545, Filter: true},
	{ID: 38546, Filter: true},
	{ID: 38547, Filter: true},
	{ID: 38548, Filter: true},
	{ID: 38549, Filter: true},
	{ID: 38550, Filter: true},
}
