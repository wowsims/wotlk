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
var ItemDeclarationOverrides = []ItemDeclaration{
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

	// Include the items we want icons for here.
	{ID: 43005}, // Pet foods
	{ID: 33874}, //
	{ID: 41174}, // Spellstones
	{ID: 41196}, //
	{ID: 12662}, // Demonic Rune
	{ID: 43015}, // Food IDs
	{ID: 34753}, // Food IDs
	{ID: 42999}, // Food IDs
	{ID: 42995}, // Food IDs
	{ID: 34754}, // Food IDs
	{ID: 34756}, // Food IDs
	{ID: 42994}, // Food IDs
	{ID: 34769}, // Food IDs
	{ID: 42996}, // Food IDs
	{ID: 34758}, // Food IDs
	{ID: 34767}, // Food IDs
	{ID: 42998}, // Food IDs
	{ID: 43000}, // Food IDs
	{ID: 27657}, // Food IDs
	{ID: 27664}, // Food IDs
	{ID: 27655}, // Food IDs
	{ID: 27658}, // Food IDs
	{ID: 33872}, // Food IDs
	{ID: 33825}, // Food IDs
	{ID: 33052}, // Food IDs
	{ID: 46376}, // Flask IDs
	{ID: 46377}, // Flask IDs
	{ID: 46378}, // Flask IDs
	{ID: 46379}, // Flask IDs
	{ID: 40079}, // Flask IDs
	{ID: 44939}, // Flask IDs
	{ID: 22861}, // Flask IDs
	{ID: 22853}, // Flask IDs
	{ID: 22866}, // Flask IDs
	{ID: 22854}, // Flask IDs
	{ID: 13512}, // Flask IDs
	{ID: 22851}, // Flask IDs
	{ID: 33208}, // Flask IDs
	{ID: 44328}, // Elixer IDs
	{ID: 40078}, // Elixer IDs
	{ID: 40109}, // Elixer IDs
	{ID: 44332}, // Elixer IDs
	{ID: 40097}, // Elixer IDs
	{ID: 40072}, // Elixer IDs
	{ID: 9088},  // Elixer IDs
	{ID: 32067}, // Elixer IDs
	{ID: 32068}, // Elixer IDs
	{ID: 22834}, // Elixer IDs
	{ID: 32062}, // Elixer IDs
	{ID: 22840}, // Elixer IDs
	{ID: 44325}, // Elixer IDs
	{ID: 44330}, // Elixer IDs
	{ID: 44327}, // Elixer IDs
	{ID: 44329}, // Elixer IDs
	{ID: 44331}, // Elixer IDs
	{ID: 39666}, // Elixer IDs
	{ID: 40073}, // Elixer IDs
	{ID: 40076}, // Elixer IDs
	{ID: 40070}, // Elixer IDs
	{ID: 40068}, // Elixer IDs
	{ID: 28103}, // Elixer IDs
	{ID: 9224},  // Elixer IDs
	{ID: 22831}, // Elixer IDs
	{ID: 22833}, // Elixer IDs
	{ID: 22827}, // Elixer IDs
	{ID: 22835}, // Elixer IDs
	{ID: 22824}, // Elixer IDs
	{ID: 28104}, // Elixer IDs
	{ID: 13452}, // Elixer IDs
	{ID: 31679}, // Elixer IDs
	{ID: 13454}, // Elixer IDs
	{ID: 33447}, // Potions / In Battle Consumes
	{ID: 33448}, // Potions / In Battle Consumes
	{ID: 40093}, // Potions / In Battle Consumes
	{ID: 40211}, // Potions / In Battle Consumes
	{ID: 40212}, // Potions / In Battle Consumes
	{ID: 22839}, // Potions / In Battle Consumes
	{ID: 22832}, // Potions / In Battle Consumes
	{ID: 22838}, // Potions / In Battle Consumes
	{ID: 13442}, // Potions / In Battle Consumes
	{ID: 31677}, // Potions / In Battle Consumes
	{ID: 22828}, // Potions / In Battle Consumes
	{ID: 22849}, // Potions / In Battle Consumes
	{ID: 22837}, // Potions / In Battle Consumes
	{ID: 20520}, // Potions / In Battle Consumes
	{ID: 22788}, // Potions / In Battle Consumes
	{ID: 22105}, // Potions / In Battle Consumes
	{ID: 42641}, // Potions / In Battle Consumes
	{ID: 40536}, // Potions / In Battle Consumes
	{ID: 41119}, // Potions / In Battle Consumes
	{ID: 40771}, // Potions / In Battle Consumes

	// Poisons
	{ID: 43233},
	{ID: 43231},
	{ID: 43235},

	// Thistle Tea
	{ID: 7676},

	// Wrath Enchant Icons
	{ID: 38375},
	{ID: 38376},
	{ID: 44069},
	{ID: 44075},
	{ID: 44875},
	{ID: 44137},
	{ID: 44138},
	{ID: 44139},
	{ID: 44140},
	{ID: 44141},
	{ID: 44876},
	{ID: 44877},
	{ID: 44878},
	{ID: 44879},
	{ID: 44067},
	{ID: 44068},
	{ID: 44129},
	{ID: 44130},
	{ID: 44131},
	{ID: 44132},
	{ID: 37330},
	{ID: 37331},
	{ID: 44494},
	{ID: 37347},
	{ID: 37349},
	{ID: 44471},
	{ID: 44472},
	{ID: 44488},
	{ID: 37340},
	{ID: 44489},
	{ID: 44484},
	{ID: 44498},
	{ID: 44944},
	{ID: 44485},
	{ID: 38371},
	{ID: 38372},
	{ID: 38373},
	{ID: 38374},
	{ID: 44963},
	{ID: 41601},
	{ID: 41602},
	{ID: 41603},
	{ID: 41604},
	{ID: 44490},
	{ID: 44491},
	{ID: 37339},
	{ID: 37344},
	{ID: 41976},
	{ID: 44486},
	{ID: 44487},
	{ID: 44492},
	{ID: 44495},
	{ID: 44496},
	{ID: 44473},
	{ID: 44483},
	{ID: 42500},
	{ID: 44936},
	{ID: 41146},
	{ID: 41167},
	{ID: 44739},
	{ID: 49633},
	{ID: 49634},
	{ID: 43468},
	{ID: 37094},
	{ID: 43464},
	{ID: 43466},
	{ID: 37092},
	{ID: 37098},

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
}
