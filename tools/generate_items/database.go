package main

import (
	"fmt"

	"github.com/wowsims/wotlk/sim/core/proto"
	"golang.org/x/exp/slices"
)

// For overriding item data.
type ItemOverride struct {
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
	Response ItemResponse
	Override ItemOverride

	QualityModifier float64
}

// For overriding gem data.
type GemOverride struct {
	ID int

	// Override fields, in case wowhead is wrong.
	Stats Stats // Only non-zero values will override
	Phase int

	Filter bool // If true, this item will be omitted from the sim.
}

type GemData struct {
	Response ItemResponse
	Override GemOverride
}

type WowDatabase struct {
	items []ItemData
	gems  []GemData
}

func NewWowDatabase(itemOverrides []ItemOverride, gemOverrides []GemOverride, tooltipsDB map[int]WowheadItemResponse) *WowDatabase {
	db := &WowDatabase{}

	for _, itemOverride := range itemOverrides {
		itemData := ItemData{
			Override: itemOverride,
			Response: tooltipsDB[itemOverride.ID],
		}
		if itemData.Response.GetName() == "" {
			continue
		}
		db.items = append(db.items, itemData)
	}

	for _, gemOverride := range gemOverrides {
		gemData := GemData{
			Override: gemOverride,
			Response: tooltipsDB[gemOverride.ID],
		}
		if gemData.Response.GetName() == "" {
			continue
		}
		db.gems = append(db.gems, gemData)
	}

	slices.SortStableFunc(db.items, func(i1, i2 ItemData) bool {
		if i1.Response.GetName() == i2.Response.GetName() {
			return i1.Override.ID < i2.Override.ID
		}
		return i1.Response.GetName() < i2.Response.GetName()
	})

	slices.SortStableFunc(db.gems, func(g1, g2 GemData) bool {
		if g1.Response.GetName() == g2.Response.GetName() {
			return g1.Override.ID < g2.Override.ID
		}
		return g1.Response.GetName() < g2.Response.GetName()
	})

	return db
}

// Returns only items which are worth including in the sim.
func (db *WowDatabase) getSimmableItems() []ItemData {
	var included []ItemData
	for _, itemData := range db.items {
		if itemData.Override.Filter {
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

		if itemData.Override.Keep {
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
		if gemData.Override.Filter {
			continue
		}
		// allow := allowList[gemData.Override.ID]
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
