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

	{ID: 28431, HandType: proto.HandType_HandTypeOneHand},
	{ID: 28432, HandType: proto.HandType_HandTypeOneHand},
	{ID: 28433, HandType: proto.HandType_HandTypeOneHand},
	{ID: 28437, HandType: proto.HandType_HandTypeOneHand},
	{ID: 28438, HandType: proto.HandType_HandTypeOneHand},
	{ID: 28439, HandType: proto.HandType_HandTypeOneHand},
	{ID: 28657, HandType: proto.HandType_HandTypeOneHand},
	{ID: 28767, HandType: proto.HandType_HandTypeOneHand},

	{ID: 17782, Filter: true}, // talisman of the binding shard
	{ID: 17783, Filter: true}, // talisman of the binding fragment
	{ID: 17802, Filter: true}, // Deprecated version of Thunderfury
	{ID: 18582, Filter: true},
	{ID: 18583, Filter: true},
	{ID: 18584, Filter: true},
	{ID: 24265, Filter: true},
	{ID: 24525, Filter: true},
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
