package database

import (
	"regexp"

	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

var ItemOverrides = []*proto.UIItem{
	{ /** Destruction Holo-gogs */ Id: 32494, ClassAllowlist: []proto.Class{proto.Class_ClassMage, proto.Class_ClassPriest, proto.Class_ClassWarlock}},
	{ /** Gadgetstorm Goggles */ Id: 32476, ClassAllowlist: []proto.Class{proto.Class_ClassShaman}},
	{ /** Magnified Moon Specs */ Id: 32480, ClassAllowlist: []proto.Class{proto.Class_ClassDruid}},
	{ /** Quad Deathblow X44 Goggles */ Id: 34353, ClassAllowlist: []proto.Class{proto.Class_ClassDruid, proto.Class_ClassRogue}},
	{ /** Hyper-Magnified Moon Specs */ Id: 35182, ClassAllowlist: []proto.Class{proto.Class_ClassDruid}},
	{ /** Lightning Etched Specs */ Id: 34355, ClassAllowlist: []proto.Class{proto.Class_ClassShaman}},
	{ /** Annihilator Holo-Gogs */ Id: 34847, ClassAllowlist: []proto.Class{proto.Class_ClassMage, proto.Class_ClassPriest, proto.Class_ClassWarlock}},

	// Balance T9 "of Conquest" Alliance set
	{Id: 48158, SetName: "Malfurion's Regalia"},
	{Id: 48159, SetName: "Malfurion's Regalia"},
	{Id: 48160, SetName: "Malfurion's Regalia"},
	{Id: 48161, SetName: "Malfurion's Regalia"},
	{Id: 48162, SetName: "Malfurion's Regalia"},

	// Valorous T8 Sets
	{Id: 45375, Phase: 2},
	{Id: 45381, Phase: 2},
	{Id: 45382, Phase: 2},
	{Id: 45376, Phase: 2},
	{Id: 45370, Phase: 2},
	{Id: 45371, Phase: 2},
	{Id: 45383, Phase: 2},
	{Id: 45372, Phase: 2},
	{Id: 45377, Phase: 2},
	{Id: 45384, Phase: 2},
	{Id: 45379, Phase: 2},
	{Id: 45385, Phase: 2},
	{Id: 45380, Phase: 2},
	{Id: 45373, Phase: 2},
	{Id: 45374, Phase: 2},
	{Id: 45391, Phase: 2},
	{Id: 45386, Phase: 2},
	{Id: 45340, Phase: 2},
	{Id: 45335, Phase: 2},
	{Id: 45336, Phase: 2},
	{Id: 45341, Phase: 2},
	{Id: 45337, Phase: 2},
	{Id: 45342, Phase: 2},
	{Id: 45338, Phase: 2},
	{Id: 45343, Phase: 2},
	{Id: 45339, Phase: 2},
	{Id: 45344, Phase: 2},
	{Id: 45419, Phase: 2},
	{Id: 45417, Phase: 2},
	{Id: 45420, Phase: 2},
	{Id: 45421, Phase: 2},
	{Id: 45422, Phase: 2},
	{Id: 45387, Phase: 2},
	{Id: 45392, Phase: 2},
	{Id: 46131, Phase: 2},
	{Id: 45365, Phase: 2},
	{Id: 45367, Phase: 2},
	{Id: 45369, Phase: 2},
	{Id: 45368, Phase: 2},
	{Id: 45388, Phase: 2},
	{Id: 45393, Phase: 2},
	{Id: 46313, Phase: 2},
	{Id: 45351, Phase: 2},
	{Id: 45355, Phase: 2},
	{Id: 45345, Phase: 2},
	{Id: 45356, Phase: 2},
	{Id: 45346, Phase: 2},
	{Id: 45347, Phase: 2},
	{Id: 45357, Phase: 2},
	{Id: 45352, Phase: 2},
	{Id: 45358, Phase: 2},
	{Id: 45348, Phase: 2},
	{Id: 45359, Phase: 2},
	{Id: 45349, Phase: 2},
	{Id: 45353, Phase: 2},
	{Id: 45354, Phase: 2},
	{Id: 45394, Phase: 2},
	{Id: 45395, Phase: 2},
	{Id: 45389, Phase: 2},
	{Id: 45360, Phase: 2},
	{Id: 45361, Phase: 2},
	{Id: 45362, Phase: 2},
	{Id: 45363, Phase: 2},
	{Id: 45364, Phase: 2},
	{Id: 45390, Phase: 2},
	{Id: 45429, Phase: 2},
	{Id: 45424, Phase: 2},
	{Id: 45430, Phase: 2},
	{Id: 45425, Phase: 2},
	{Id: 45426, Phase: 2},
	{Id: 45431, Phase: 2},
	{Id: 45427, Phase: 2},
	{Id: 45432, Phase: 2},
	{Id: 45428, Phase: 2},
	{Id: 45433, Phase: 2},
	{Id: 45396, Phase: 2},
	{Id: 45397, Phase: 2},
	{Id: 45398, Phase: 2},
	{Id: 45399, Phase: 2},
	{Id: 45400, Phase: 2},
	{Id: 45413, Phase: 2},
	{Id: 45412, Phase: 2},
	{Id: 45406, Phase: 2},
	{Id: 45414, Phase: 2},
	{Id: 45401, Phase: 2},
	{Id: 45411, Phase: 2},
	{Id: 45402, Phase: 2},
	{Id: 45408, Phase: 2},
	{Id: 45409, Phase: 2},
	{Id: 45403, Phase: 2},
	{Id: 45415, Phase: 2},
	{Id: 45410, Phase: 2},
	{Id: 45404, Phase: 2},
	{Id: 45405, Phase: 2},
	{Id: 45416, Phase: 2},

	// Other items Wowhead has the wrong phase listed for
	// Ick's loot table from Pit of Saron
	{Id: 49812, Phase: 4},
	{Id: 49808, Phase: 4},
	{Id: 49811, Phase: 4},
	{Id: 49807, Phase: 4},
	{Id: 49810, Phase: 4},
	{Id: 49809, Phase: 4},

	// Drape of Icy Intent
	{Id: 45461, Phase: 2},

	// Valentine's day event rewards
	{Id: 51804, Phase: 2},
	{Id: 51805, Phase: 2},
	{Id: 51806, Phase: 2},
	{Id: 51807, Phase: 2},
	{Id: 51808, Phase: 2},
}

// Keep these sorted by item ID.
var ItemAllowList = map[int32]struct{}{
	11815: {}, // Hand of Justice
	12590: {}, // Felstriker
	15808: {}, // Fine Light Crossbow (for hunter testing).
	18843: {},
	18844: {},
	18847: {},
	18848: {},
	19019: {}, // Thunderfury
	19808: {}, // Rockhide Strongfish
	20837: {}, // Sunstrider Axe
	20966: {}, // Jade Pendant of Blasting
	21625: {}, // Scarab Brooch
	24114: {}, // Braided Eternium Chain
	28572: {}, // Blade of the Unrequited
	28830: {}, // Dragonspine Trophy
	29383: {}, // Bloodlust Brooch
	29387: {}, // Gnomeregan Auto-Blocker 600
	29994: {}, // Thalassian Wildercloak
	29996: {}, // Rod of the Sun King
	30032: {}, // Red Belt of Battle
	30627: {}, // Tsunami Talisman
	30720: {}, // Serpent-Coil Braid
	31193: {}, // Blade of Unquenched Thirst
	32387: {}, // Idol of the Raven Goddess
	32658: {}, // Badge of Tenacity
	33135: {}, // Falling Star
	33140: {}, // Blood of Amber
	33143: {}, // Stone of Blades
	33144: {}, // Facet of Eternity
	33504: {}, // Libram of Divine Purpose
	33506: {}, // Skycall Totem
	33507: {}, // Stonebreaker's Totem
	33508: {}, // Idol of Budding Life
	33510: {}, // Unseen moon idol
	33829: {}, // Hex Shrunken Head
	33831: {}, // Berserkers Call
	34472: {}, // Shard of Contempt
	34473: {}, // Commendation of Kael'thas
	37032: {}, // Edge of the Tuskarr
	37574: {}, // Libram of Furious Blows
	38072: {}, // Thunder Capacitor
	38212: {}, // Death Knight's Anguish
	38287: {}, // Empty Mug of Direbrew
	38289: {}, // Coren's Lucky Coin
	39208: {}, // Sigil of the Dark Rider
	41752: {}, // Brunnhildar Axe
	6360:  {}, // Steelscale Crushfish
	8345:  {}, // Wolfshead Helm
	9449:  {}, // Manual Crowd Pummeler

	// Sets
	27510: {}, // Tidefury Gauntlets
	27802: {}, // Tidefury Shoulderguards
	27909: {}, // Tidefury Kilt
	28231: {}, // Tidefury Chestpiece
	28349: {}, // Tidefury Helm

	15056: {}, // Stormshroud Armor
	15057: {}, // Stormshroud Pants
	15058: {}, // Stormshroud Shoulders
	21278: {}, // Stormshroud Gloves

	// Undead Slaying Sets
	// Plate
	43068: {},
	43069: {},
	43070: {},
	43071: {},
	// Cloth
	43072: {},
	43073: {},
	43074: {},
	43075: {},
	// Mail
	43076: {},
	43077: {},
	43078: {},
	43079: {},
	//Leather
	43080: {},
	43081: {},
	43082: {},
	43083: {},
}

// Keep these sorted by item ID.
var ItemDenyList = map[int32]struct{}{
	17782: {}, // talisman of the binding shard
	17783: {}, // talisman of the binding fragment
	17802: {}, // Deprecated version of Thunderfury
	18582: {},
	18583: {},
	18584: {},
	24265: {},
	32384: {},
	32421: {},
	32422: {},
	33482: {},
	33350: {},
	34576: {}, // Battlemaster's Cruelty
	34577: {}, // Battlemaster's Depreavity
	34578: {}, // Battlemaster's Determination
	34579: {}, // Battlemaster's Audacity
	34580: {}, // Battlemaster's Perseverence
	50251: {}, // 'one hand shadows edge'
	53500: {}, // Tectonic Plate

	48880: {}, // DK's Tier 9 Duplicates
	48881: {}, // DK's Tier 9 Duplicates
	48882: {}, // DK's Tier 9 Duplicates
	48883: {}, // DK's Tier 9 Duplicates
	48884: {}, // DK's Tier 9 Duplicates
	48885: {}, // DK's Tier 9 Duplicates
	48886: {}, // DK's Tier 9 Duplicates
	48887: {}, // DK's Tier 9 Duplicates
	48888: {}, // DK's Tier 9 Duplicates
	48889: {}, // DK's Tier 9 Duplicates
	48890: {}, // DK's Tier 9 Duplicates
	48891: {}, // DK's Tier 9 Duplicates
	48892: {}, // DK's Tier 9 Duplicates
	48893: {}, // DK's Tier 9 Duplicates
	48894: {}, // DK's Tier 9 Duplicates
	48895: {}, // DK's Tier 9 Duplicates
	48896: {}, // DK's Tier 9 Duplicates
	48897: {}, // DK's Tier 9 Duplicates
	48898: {}, // DK's Tier 9 Duplicates
	48899: {}, // DK's Tier 9 Duplicates
}

// Item icons to include in the DB, so they don't need to be separately loaded in the UI.
var ExtraItemIcons = []int32{
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
	36892,
	40093,
	40211,
	40212,
	40536,
	40771,
	41119,
	41166,
	42545,
	42641,

	// Poisons
	43231,
	43233,
	43235,

	// Thistle Tea
	7676,

	// Scrolls
	37094,
	43466,
	43464,
	37092,
	37098,
	43468,

	// Drums
	49633,
	49634,
}

// Raid buffs / debuffs
var SharedSpellsIcons = []int32{
	// Revitalize, Rejuv, WG
	48545,
	26982,
	53251,

	// Registered CD's
	49016,
	57933,
	64382,
	10060,
	16190,
	29166,
	53530,
	33206,
	2825,
	54758,

	// Raid Buffs
	43002,
	57567,
	54038,

	48470,
	17051,

	25898,
	25899,

	48942,
	20140,
	58753,
	16293,

	48161,
	14767,

	58643,
	52456,
	57623,

	48073,

	48934,
	20045,
	47436,

	53138,
	30809,
	19506,

	31869,
	31583,
	34460,

	57472,
	50720,

	53648,

	47440,
	12861,
	47982,
	18696,

	48938,
	20245,
	58774,
	16206,

	17007,
	34300,
	29801,

	55610,
	65990,
	29193,

	48160,
	31878,
	53292,
	54118,
	44561,

	24907,
	48396,
	51470,

	3738,
	47240,
	57722,
	58656,

	54043,
	48170,
	31025,
	31035,
	6562,
	31033,
	53307,
	16840,
	54648,

	// Raid Debuffs
	8647,
	47467,
	55754,

	770,
	33602,
	50511,
	18180,
	56631,
	53598,

	26016,
	47437,
	12879,
	48560,
	16862,
	55487,

	33876,
	46855,
	57393,

	30706,
	20337,
	58410,

	47502,
	12666,
	55095,
	51456,
	53696,
	48485,

	3043,
	29859,
	58413,
	65855,

	17800,
	17803,
	12873,
	28593,

	33198,
	51161,
	48511,
	47865,

	20271,
	53408,

	11374,
}

// If any of these match the item name, don't include it.
var DenyListNameRegexes = []*regexp.Regexp{
	regexp.MustCompile(`30 Epic`),
	regexp.MustCompile(`63 Blue`),
	regexp.MustCompile(`63 Green`),
	regexp.MustCompile(`66 Epic`),
	regexp.MustCompile(`90 Epic`),
	regexp.MustCompile(`90 Green`),
	regexp.MustCompile(`Boots 1`),
	regexp.MustCompile(`Boots 2`),
	regexp.MustCompile(`Boots 3`),
	regexp.MustCompile(`Bracer 1`),
	regexp.MustCompile(`Bracer 2`),
	regexp.MustCompile(`Bracer 3`),
	regexp.MustCompile(`DB\d`),
	regexp.MustCompile(`DEPRECATED`),
	regexp.MustCompile(`Deprecated: Keanna`),
	regexp.MustCompile(`Indalamar`),
	regexp.MustCompile(`Monster -`),
	regexp.MustCompile(`NEW`),
	regexp.MustCompile(`PH`),
	regexp.MustCompile(`QR XXXX`),
	regexp.MustCompile(`TEST`),
	regexp.MustCompile(`Test`),
	regexp.MustCompile(`zOLD`),
}

// Allows manual overriding for Gem fields in case WowHead is wrong.
var GemOverrides = []*proto.UIGem{
	{Id: 33131, Stats: stats.Stats{stats.AttackPower: 32, stats.RangedAttackPower: 32}.ToFloatArray()},
}
var GemDenyList = map[int32]struct{}{
	// pvp non-unique gems not in game currently.
	32735: struct{}{},
	35489: struct{}{},
	38545: struct{}{},
	38546: struct{}{},
	38547: struct{}{},
	38548: struct{}{},
	38549: struct{}{},
	38550: struct{}{},
}
