package main

import (
	"fmt"

	"github.com/wowsims/wotlk/sim/core/proto"
	"golang.org/x/exp/slices"
)

type ItemDeclaration struct {
	ID int

	// Override fields, in case wowhead is wrong.
	Stats          Stats // Only non-zero values will override
	ClassAllowlist []proto.Class
	Phase          int
	HandType       proto.HandType // Overrides hand type.
	Filter         bool           // If true, this item will be omitted from the sim.
	Keep           bool           // If true, keep this item even if it would otherwise be filtered.
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

type WowDatabase struct {
	items []ItemData
	gems  []GemData
}

func NewWowDatabase(itemDeclarations []ItemDeclaration, gemDeclarations []GemDeclaration, tooltipsDB map[int]WowheadItemResponse) *WowDatabase {
	db := &WowDatabase{}

	for _, itemDeclaration := range itemDeclarations {
		itemData := ItemData{
			Declaration: itemDeclaration,
			Response:    tooltipsDB[itemDeclaration.ID],
		}
		if itemData.Response.GetName() == "" {
			continue
		}
		db.items = append(db.items, itemData)
	}

	for _, gemDeclaration := range gemDeclarations {
		gemData := GemData{
			Declaration: gemDeclaration,
			Response:    tooltipsDB[gemDeclaration.ID],
		}
		if gemData.Response.GetName() == "" {
			continue
		}
		db.gems = append(db.gems, gemData)
	}

	slices.SortStableFunc(db.items, func(i1, i2 ItemData) bool {
		if i1.Response.GetName() == i2.Response.GetName() {
			return i1.Declaration.ID < i2.Declaration.ID
		}
		return i1.Response.GetName() < i2.Response.GetName()
	})

	slices.SortStableFunc(db.gems, func(g1, g2 GemData) bool {
		if g1.Response.GetName() == g2.Response.GetName() {
			return g1.Declaration.ID < g2.Declaration.ID
		}
		return g1.Response.GetName() < g2.Response.GetName()
	})

	return db
}

// Returns only items which are worth including in the sim.
func (db *WowDatabase) getSimmableItems() []ItemData {
	var included []ItemData
	for _, itemData := range db.items {
		if itemData.Declaration.Filter {
			continue
		}

		deny := false
		for _, pattern := range denyListNameRegexes {
			if pattern.MatchString(itemData.Response.GetName()) {
				deny = true
				break
			}
		}
		if deny {
			continue
		}

		if !itemData.Response.IsEquippable() {
			continue
		}

		if itemData.Declaration.Keep {
			included = append(included, itemData)
			continue
		}

		itemLevel := itemData.Response.GetItemLevel()
		qual := itemData.Response.GetQuality()
		if qual < int(proto.ItemQuality_ItemQualityUncommon) {
			continue
		} else if qual > int(proto.ItemQuality_ItemQualityLegendary) {
			continue
		} else if qual < int(proto.ItemQuality_ItemQualityEpic) {
			if itemLevel < 145 {
				continue
			}
			if itemLevel < 149 && itemData.Response.GetItemSetName() == "" {
				continue
			}
		} else {
			// Epic and legendary items might come from classic, so use a lower ilvl threshold.
			if itemLevel < 140 {
				continue
			}
		}
		if itemLevel == 0 {
			fmt.Printf("Missing ilvl: %s\n", itemData.Response.GetName())
		}

		included = append(included, itemData)
	}

	return included
}

// Returns only gems which are worth including in the sim.
func (db *WowDatabase) getSimmableGems() []GemData {
	var included []GemData

	for _, gemData := range db.gems {
		if gemData.Declaration.Filter {
			continue
		}
		// allow := allowList[gemData.Declaration.ID]
		allow := false
		if !allow {
			if gemData.Response.GetQuality() < int(proto.ItemQuality_ItemQualityUncommon) {
				continue
			}
			// if gemData.Response.GetPhase() == 0 {
			// 	continue
			// }
		}
		included = append(included, gemData)
	}

	return included
}
