package main

import (
	"github.com/wowsims/wotlk/sim/core/proto"
)

type ItemDeclaration struct {
	ID int

	// Override fields, in case wowhead is wrong.
	Stats          Stats // Only non-zero values will override
	ClassAllowlist []proto.Class
	Phase          int
	HandType       proto.HandType // Overrides hand type.
	Filter         bool           // If true, this item will be omitted from the sim.
}
type ItemData struct {
	Declaration ItemDeclaration
	Response    ItemResponse

	QualityModifier float64
}

type GemDeclaration struct {
	ID int

	// Override fields, in case wowhead is wrong.
	Stats Stats // Only non-zero values will override
	Phase int

	Filter bool // If true, this item will be omitted from the sim.
}
type GemData struct {
	Declaration GemDeclaration
	Response    ItemResponse
}

// Allows manual overriding for Gem fields in case WowHead is wrong.
var GemDeclarationOverrides = []GemDeclaration{
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

// Allows manual overriding for Item fields in case WowHead is wrong.
var ItemDeclarationOverrides = map[int]ItemDeclaration{
	32494: { /** Destruction Holo-gogs */ ID: 32494, ClassAllowlist: []proto.Class{proto.Class_ClassMage, proto.Class_ClassPriest, proto.Class_ClassWarlock}},
	32476: { /** Gadgetstorm Goggles */ ID: 32476, ClassAllowlist: []proto.Class{proto.Class_ClassShaman}},
	32480: { /** Magnified Moon Specs */ ID: 32480, ClassAllowlist: []proto.Class{proto.Class_ClassDruid}},
	34353: { /** Quad Deathblow X44 Goggles */ ID: 34353, ClassAllowlist: []proto.Class{proto.Class_ClassDruid, proto.Class_ClassRogue}},
	35182: { /** Hyper-Magnified Moon Specs */ ID: 35182, ClassAllowlist: []proto.Class{proto.Class_ClassDruid}},
	34355: { /** Lightning Etched Specs */ ID: 34355, ClassAllowlist: []proto.Class{proto.Class_ClassShaman}},
	34847: { /** Annihilator Holo-Gogs */ ID: 34847, ClassAllowlist: []proto.Class{proto.Class_ClassMage, proto.Class_ClassPriest, proto.Class_ClassWarlock}},
	29993: { /** Twinblade of the Pheonix */ ID: 29993, Stats: Stats{proto.Stat_StatRangedAttackPower: 108}},
	30883: { /** Pillar of Ferocity */ ID: 30883, Stats: Stats{proto.Stat_StatArmor: 550}},
	37723: { /** Incisor Fragment */ ID: 37723, Stats: Stats{proto.Stat_StatArmorPenetration: 0}},

	// TBC blacksmithing crafted weapons
	28431: {ID: 28431, HandType: proto.HandType_HandTypeOneHand},
	28432: {ID: 28432, HandType: proto.HandType_HandTypeOneHand},
	28433: {ID: 28433, HandType: proto.HandType_HandTypeOneHand},
	28437: {ID: 28437, HandType: proto.HandType_HandTypeOneHand},
	28438: {ID: 28438, HandType: proto.HandType_HandTypeOneHand},
	28439: {ID: 28439, HandType: proto.HandType_HandTypeOneHand},
	28657: {ID: 28657, HandType: proto.HandType_HandTypeOneHand},
	28767: {ID: 28767, HandType: proto.HandType_HandTypeOneHand},

	29994: {ID: 29994}, // Wildercloak
	29387: {ID: 29387}, // Gnomeregan Auto-Blocker 600
	38289: {ID: 38289}, // Coren's Lucky Coin
	34473: {ID: 34473}, // Commendation of Kaelthas

	// WotLK Halloween Event items
	49121: {ID: 49121, Phase: 1}, // Ring of Ghoulish Glee
	49123: {ID: 49123, Phase: 1}, // The Horseman's Seal
	49124: {ID: 49124, Phase: 1}, // Wicked Witch's Band
	49128: {ID: 49128, Phase: 1}, // The Horseman's Baleful Blade
	49126: {ID: 49126, Phase: 1}, // The Horseman's Horrific Helm

	// Valorous T8 Sets
	45375: {ID: 45375, Phase: 2},
	45381: {ID: 45381, Phase: 2},
	45382: {ID: 45382, Phase: 2},
	45376: {ID: 45376, Phase: 2},
	45370: {ID: 45370, Phase: 2},
	45371: {ID: 45371, Phase: 2},
	45383: {ID: 45383, Phase: 2},
	45372: {ID: 45372, Phase: 2},
	45377: {ID: 45377, Phase: 2},
	45384: {ID: 45384, Phase: 2},
	45379: {ID: 45379, Phase: 2},
	45385: {ID: 45385, Phase: 2},
	45380: {ID: 45380, Phase: 2},
	45373: {ID: 45373, Phase: 2},
	45374: {ID: 45374, Phase: 2},
	45391: {ID: 45391, Phase: 2},
	45386: {ID: 45386, Phase: 2},
	45340: {ID: 45340, Phase: 2},
	45335: {ID: 45335, Phase: 2},
	45336: {ID: 45336, Phase: 2},
	45341: {ID: 45341, Phase: 2},
	45337: {ID: 45337, Phase: 2},
	45342: {ID: 45342, Phase: 2},
	45338: {ID: 45338, Phase: 2},
	45343: {ID: 45343, Phase: 2},
	45339: {ID: 45339, Phase: 2},
	45344: {ID: 45344, Phase: 2},
	45419: {ID: 45419, Phase: 2},
	45417: {ID: 45417, Phase: 2},
	45420: {ID: 45420, Phase: 2},
	45421: {ID: 45421, Phase: 2},
	45422: {ID: 45422, Phase: 2},
	45387: {ID: 45387, Phase: 2},
	45392: {ID: 45392, Phase: 2},
	46131: {ID: 46131, Phase: 2},
	45365: {ID: 45365, Phase: 2},
	45367: {ID: 45367, Phase: 2},
	45369: {ID: 45369, Phase: 2},
	45368: {ID: 45368, Phase: 2},
	45388: {ID: 45388, Phase: 2},
	45393: {ID: 45393, Phase: 2},
	46313: {ID: 46313, Phase: 2},
	45351: {ID: 45351, Phase: 2},
	45355: {ID: 45355, Phase: 2},
	45345: {ID: 45345, Phase: 2},
	45356: {ID: 45356, Phase: 2},
	45346: {ID: 45346, Phase: 2},
	45347: {ID: 45347, Phase: 2},
	45357: {ID: 45357, Phase: 2},
	45352: {ID: 45352, Phase: 2},
	45358: {ID: 45358, Phase: 2},
	45348: {ID: 45348, Phase: 2},
	45359: {ID: 45359, Phase: 2},
	45349: {ID: 45349, Phase: 2},
	45353: {ID: 45353, Phase: 2},
	45354: {ID: 45354, Phase: 2},
	45394: {ID: 45394, Phase: 2},
	45395: {ID: 45395, Phase: 2},
	45389: {ID: 45389, Phase: 2},
	45360: {ID: 45360, Phase: 2},
	45361: {ID: 45361, Phase: 2},
	45362: {ID: 45362, Phase: 2},
	45363: {ID: 45363, Phase: 2},
	45364: {ID: 45364, Phase: 2},
	45390: {ID: 45390, Phase: 2},
	45429: {ID: 45429, Phase: 2},
	45424: {ID: 45424, Phase: 2},
	45430: {ID: 45430, Phase: 2},
	45425: {ID: 45425, Phase: 2},
	45426: {ID: 45426, Phase: 2},
	45431: {ID: 45431, Phase: 2},
	45427: {ID: 45427, Phase: 2},
	45432: {ID: 45432, Phase: 2},
	45428: {ID: 45428, Phase: 2},
	45433: {ID: 45433, Phase: 2},
	45396: {ID: 45396, Phase: 2},
	45397: {ID: 45397, Phase: 2},
	45398: {ID: 45398, Phase: 2},
	45399: {ID: 45399, Phase: 2},
	45400: {ID: 45400, Phase: 2},
	45413: {ID: 45413, Phase: 2},
	45412: {ID: 45412, Phase: 2},
	45406: {ID: 45406, Phase: 2},
	45414: {ID: 45414, Phase: 2},
	45401: {ID: 45401, Phase: 2},
	45411: {ID: 45411, Phase: 2},
	45402: {ID: 45402, Phase: 2},
	45408: {ID: 45408, Phase: 2},
	45409: {ID: 45409, Phase: 2},
	45403: {ID: 45403, Phase: 2},
	45415: {ID: 45415, Phase: 2},
	45410: {ID: 45410, Phase: 2},
	45404: {ID: 45404, Phase: 2},
	45405: {ID: 45405, Phase: 2},
	45416: {ID: 45416, Phase: 2},

	// Other items Wowhead has the wrong phase listed for
	// Ick's loot table from Pit of Saron
	49812: {ID: 49812, Phase: 4},
	49808: {ID: 49808, Phase: 4},
	49811: {ID: 49811, Phase: 4},
	49807: {ID: 49807, Phase: 4},
	49810: {ID: 49810, Phase: 4},
	49809: {ID: 49809, Phase: 4},

	// Include the items we want icons for here.
	43005: {ID: 43005}, // Pet foods
	33874: {ID: 33874}, //
	41174: {ID: 41174}, // Spellstones
	41196: {ID: 41196}, //
	12662: {ID: 12662}, // Demonic Rune
	43015: {ID: 43015}, // Food IDs
	34753: {ID: 34753}, // Food IDs
	42999: {ID: 42999}, // Food IDs
	42995: {ID: 42995}, // Food IDs
	34754: {ID: 34754}, // Food IDs
	34756: {ID: 34756}, // Food IDs
	42994: {ID: 42994}, // Food IDs
	34769: {ID: 34769}, // Food IDs
	42996: {ID: 42996}, // Food IDs
	34758: {ID: 34758}, // Food IDs
	34767: {ID: 34767}, // Food IDs
	42998: {ID: 42998}, // Food IDs
	43000: {ID: 43000}, // Food IDs
	27657: {ID: 27657}, // Food IDs
	27664: {ID: 27664}, // Food IDs
	27655: {ID: 27655}, // Food IDs
	27658: {ID: 27658}, // Food IDs
	33872: {ID: 33872}, // Food IDs
	33825: {ID: 33825}, // Food IDs
	33052: {ID: 33052}, // Food IDs
	46376: {ID: 46376}, // Flask IDs
	46377: {ID: 46377}, // Flask IDs
	46378: {ID: 46378}, // Flask IDs
	46379: {ID: 46379}, // Flask IDs
	40079: {ID: 40079}, // Flask IDs
	44939: {ID: 44939}, // Flask IDs
	22861: {ID: 22861}, // Flask IDs
	22853: {ID: 22853}, // Flask IDs
	22866: {ID: 22866}, // Flask IDs
	22854: {ID: 22854}, // Flask IDs
	13512: {ID: 13512}, // Flask IDs
	22851: {ID: 22851}, // Flask IDs
	33208: {ID: 33208}, // Flask IDs
	44328: {ID: 44328}, // Elixer IDs
	40078: {ID: 40078}, // Elixer IDs
	40109: {ID: 40109}, // Elixer IDs
	44332: {ID: 44332}, // Elixer IDs
	40097: {ID: 40097}, // Elixer IDs
	40072: {ID: 40072}, // Elixer IDs
	9088:  {ID: 9088},  // Elixer IDs
	32067: {ID: 32067}, // Elixer IDs
	32068: {ID: 32068}, // Elixer IDs
	22834: {ID: 22834}, // Elixer IDs
	32062: {ID: 32062}, // Elixer IDs
	22840: {ID: 22840}, // Elixer IDs
	44325: {ID: 44325}, // Elixer IDs
	44330: {ID: 44330}, // Elixer IDs
	44327: {ID: 44327}, // Elixer IDs
	44329: {ID: 44329}, // Elixer IDs
	44331: {ID: 44331}, // Elixer IDs
	39666: {ID: 39666}, // Elixer IDs
	40073: {ID: 40073}, // Elixer IDs
	40076: {ID: 40076}, // Elixer IDs
	40070: {ID: 40070}, // Elixer IDs
	40068: {ID: 40068}, // Elixer IDs
	28103: {ID: 28103}, // Elixer IDs
	9224:  {ID: 9224},  // Elixer IDs
	22831: {ID: 22831}, // Elixer IDs
	22833: {ID: 22833}, // Elixer IDs
	22827: {ID: 22827}, // Elixer IDs
	22835: {ID: 22835}, // Elixer IDs
	22824: {ID: 22824}, // Elixer IDs
	28104: {ID: 28104}, // Elixer IDs
	13452: {ID: 13452}, // Elixer IDs
	31679: {ID: 31679}, // Elixer IDs
	13454: {ID: 13454}, // Elixer IDs
	33447: {ID: 33447}, // Potions / In Battle Consumes
	33448: {ID: 33448}, // Potions / In Battle Consumes
	40093: {ID: 40093}, // Potions / In Battle Consumes
	40211: {ID: 40211}, // Potions / In Battle Consumes
	40212: {ID: 40212}, // Potions / In Battle Consumes
	22839: {ID: 22839}, // Potions / In Battle Consumes
	22832: {ID: 22832}, // Potions / In Battle Consumes
	22838: {ID: 22838}, // Potions / In Battle Consumes
	13442: {ID: 13442}, // Potions / In Battle Consumes
	31677: {ID: 31677}, // Potions / In Battle Consumes
	22828: {ID: 22828}, // Potions / In Battle Consumes
	22849: {ID: 22849}, // Potions / In Battle Consumes
	22837: {ID: 22837}, // Potions / In Battle Consumes
	20520: {ID: 20520}, // Potions / In Battle Consumes
	22788: {ID: 22788}, // Potions / In Battle Consumes
	22105: {ID: 22105}, // Potions / In Battle Consumes
	42641: {ID: 42641}, // Potions / In Battle Consumes
	40536: {ID: 40536}, // Potions / In Battle Consumes
	41119: {ID: 41119}, // Potions / In Battle Consumes
	40771: {ID: 40771}, // Potions / In Battle Consumes

	// Poisons
	43233: {ID: 43233},
	43231: {ID: 43231},
	43235: {ID: 43235},

	// Thistle Tea
	7676: {ID: 7676},

	// Wrath Enchant Icons
	38375: {ID: 38375},
	38376: {ID: 38376},
	44069: {ID: 44069},
	44075: {ID: 44075},
	44875: {ID: 44875},
	44137: {ID: 44137},
	44138: {ID: 44138},
	44139: {ID: 44139},
	44140: {ID: 44140},
	44141: {ID: 44141},
	44876: {ID: 44876},
	44877: {ID: 44877},
	44878: {ID: 44878},
	44879: {ID: 44879},
	44067: {ID: 44067},
	44068: {ID: 44068},
	44129: {ID: 44129},
	44130: {ID: 44130},
	44131: {ID: 44131},
	44132: {ID: 44132},
	37330: {ID: 37330},
	37331: {ID: 37331},
	44494: {ID: 44494},
	37347: {ID: 37347},
	37349: {ID: 37349},
	44471: {ID: 44471},
	44472: {ID: 44472},
	44488: {ID: 44488},
	37340: {ID: 37340},
	44489: {ID: 44489},
	44484: {ID: 44484},
	44498: {ID: 44498},
	44944: {ID: 44944},
	44485: {ID: 44485},
	38371: {ID: 38371},
	38372: {ID: 38372},
	38373: {ID: 38373},
	38374: {ID: 38374},
	44963: {ID: 44963},
	41601: {ID: 41601},
	41602: {ID: 41602},
	41603: {ID: 41603},
	41604: {ID: 41604},
	44490: {ID: 44490},
	44491: {ID: 44491},
	37339: {ID: 37339},
	37344: {ID: 37344},
	41976: {ID: 41976},
	44486: {ID: 44486},
	44487: {ID: 44487},
	44492: {ID: 44492},
	44495: {ID: 44495},
	44496: {ID: 44496},
	44473: {ID: 44473},
	44483: {ID: 44483},
	42500: {ID: 42500},
	44936: {ID: 44936},
	41146: {ID: 41146},
	41167: {ID: 41167},
	44739: {ID: 44739},
	49633: {ID: 49633},
	49634: {ID: 49634},
	43468: {ID: 43468},
	37094: {ID: 37094},
	43464: {ID: 43464},
	43466: {ID: 43466},
	37092: {ID: 37092},
	37098: {ID: 37098},

	// Filter out these items
	17782: {ID: 17782, Filter: true}, // talisman of the binding shard
	17783: {ID: 17783, Filter: true}, // talisman of the binding fragment
	17802: {ID: 17802, Filter: true}, // Deprecated version of Thunderfury
	18582: {ID: 18582, Filter: true},
	18583: {ID: 18583, Filter: true},
	18584: {ID: 18584, Filter: true},
	24265: {ID: 24265, Filter: true},
	32384: {ID: 32384, Filter: true},
	32421: {ID: 32421, Filter: true},
	32422: {ID: 32422, Filter: true},
	33482: {ID: 33482, Filter: true},
	34576: {ID: 34576, Filter: true}, // Battlemaster's Cruelty
	34577: {ID: 34577, Filter: true}, // Battlemaster's Depreavity
	34578: {ID: 34578, Filter: true}, // Battlemaster's Determination
	34579: {ID: 34579, Filter: true}, // Battlemaster's Audacity
	34580: {ID: 34580, Filter: true}, // Battlemaster's Perseverence
}
