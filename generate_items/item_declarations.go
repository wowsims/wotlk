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
	{ /** Beast-tamer's Shoulders */ ID: 30892, ClassAllowlist: []proto.Class{proto.Class_ClassHunter}},
	{ /** Incisor Fragment */ ID: 37723, Stats: Stats{proto.Stat_StatArmorPenetration: 0}},

	{ID: 28431, HandType: proto.HandType_HandTypeOneHand},
	{ID: 28432, HandType: proto.HandType_HandTypeOneHand},
	{ID: 28433, HandType: proto.HandType_HandTypeOneHand},
	{ID: 28437, HandType: proto.HandType_HandTypeOneHand},
	{ID: 28438, HandType: proto.HandType_HandTypeOneHand},
	{ID: 28439, HandType: proto.HandType_HandTypeOneHand},
	{ID: 28657, HandType: proto.HandType_HandTypeOneHand},
	{ID: 28767, HandType: proto.HandType_HandTypeOneHand},

	{ID: 12662}, // Demonic Rune
	{ID: 43015}, // Food IDs to get icons
	{ID: 34753}, // Food IDs to get icons
	{ID: 42999}, // Food IDs to get icons
	{ID: 42995}, // Food IDs to get icons
	{ID: 34754}, // Food IDs to get icons
	{ID: 34756}, // Food IDs to get icons
	{ID: 42994}, // Food IDs to get icons
	{ID: 34769}, // Food IDs to get icons
	{ID: 42996}, // Food IDs to get icons
	{ID: 34758}, // Food IDs to get icons
	{ID: 34767}, // Food IDs to get icons
	{ID: 42998}, // Food IDs to get icons
	{ID: 43000}, // Food IDs to get icons
	{ID: 27657}, // Food IDs to get icons
	{ID: 27664}, // Food IDs to get icons
	{ID: 27655}, // Food IDs to get icons
	{ID: 27658}, // Food IDs to get icons
	{ID: 33872}, // Food IDs to get icons
	{ID: 33825}, // Food IDs to get icons
	{ID: 33052}, // Food IDs to get icons
	{ID: 46376}, // Flask IDs to get icons
	{ID: 46377}, // Flask IDs to get icons
	{ID: 46378}, // Flask IDs to get icons
	{ID: 46379}, // Flask IDs to get icons
	{ID: 40079}, // Flask IDs to get icons
	{ID: 44939}, // Flask IDs to get icons
	{ID: 22861}, // Flask IDs to get icons
	{ID: 22853}, // Flask IDs to get icons
	{ID: 22866}, // Flask IDs to get icons
	{ID: 22854}, // Flask IDs to get icons
	{ID: 13512}, // Flask IDs to get icons
	{ID: 22851}, // Flask IDs to get icons
	{ID: 33208}, // Flask IDs to get icons
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
	{ID: 27155}, // Elixer IDs
	{ID: 33447}, // Include consumables so we get their icons
	{ID: 33448}, // Include consumables so we get their icons
	{ID: 40093}, // Include consumables so we get their icons
	{ID: 40211}, // Include consumables so we get their icons
	{ID: 40212}, // Include consumables so we get their icons
	{ID: 22839}, // Include consumables so we get their icons
	{ID: 22832}, // Include consumables so we get their icons
	{ID: 22838}, // Include consumables so we get their icons
	{ID: 13442}, // Include consumables so we get their icons
	{ID: 31677}, // Include consumables so we get their icons
	{ID: 22828}, // Include consumables so we get their icons
	{ID: 22849}, // Include consumables so we get their icons
	{ID: 22837}, // Include consumables so we get their icons
	{ID: 20520}, // Include consumables so we get their icons
	{ID: 22788}, // Include consumables so we get their icons
	{ID: 22105}, // Include consumables so we get their icons
	{ID: 42641}, // Include consumables so we get their icons
	{ID: 40536}, // Include consumables so we get their icons
	{ID: 41119}, // Include consumables so we get their icons
	{ID: 40771}, // Include consumables so we get their icons

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
