package main

import (
	"fmt"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	_ "github.com/wowsims/wotlk/sim/encounters" // Needed for preset encounters.
	"golang.org/x/exp/slices"
	googleProto "google.golang.org/protobuf/proto"
)

func mergeItemProtos(dst, src *proto.UIItem) {
	// googleproto.Merge concatenates lists but we want replacement, so do them manually.
	if src.Stats != nil {
		dst.Stats = src.Stats
		src.Stats = nil
	}
	// TODO: Other stat fields
	googleProto.Merge(dst, src)
}

func mergeGemProtos(dst, src *proto.UIGem) {
	// googleproto.Merge concatenates lists but we want replacement, so do them manually.
	if src.Stats != nil {
		dst.Stats = src.Stats
		src.Stats = nil
	}
	googleProto.Merge(dst, src)
}

type WowDatabase struct {
	items    map[int32]*proto.UIItem
	enchants []*proto.UIEnchant
	gems     map[int32]*proto.UIGem

	itemIcons  []*proto.IconData
	spellIcons []*proto.IconData

	encounters []*proto.PresetEncounter
}

func NewWowDatabase(itemOverrides []*proto.UIItem, gemOverrides []*proto.UIGem, enchantOverrides []*proto.UIEnchant, itemTooltipsDB map[int]WowheadItemResponse, spellTooltipsDB map[int]WowheadItemResponse) *WowDatabase {
	db := &WowDatabase{
		items:      make(map[int32]*proto.UIItem),
		enchants:   enchantOverrides,
		gems:       make(map[int32]*proto.UIGem),
		encounters: core.PresetEncounters,
	}

	for id, response := range itemTooltipsDB {
		if response.IsEquippable() {
			itemProto := response.ToItemProto()
			itemProto.Id = int32(id)
			db.items[itemProto.Id] = itemProto
		}
	}
	for _, itemOverride := range itemOverrides {
		if _, ok := db.items[itemOverride.Id]; ok {
			mergeItemProtos(db.items[itemOverride.Id], itemOverride)
		}
	}

	for id, response := range itemTooltipsDB {
		if response.IsGem() {
			gemProto := response.ToGemProto()
			gemProto.Id = int32(id)
			db.gems[gemProto.Id] = gemProto
		}
	}
	for _, gemOverride := range gemOverrides {
		if _, ok := db.gems[gemOverride.Id]; ok {
			mergeGemProtos(db.gems[gemOverride.Id], gemOverride)
		}
	}

	for _, enchant := range db.enchants {
		if enchant.ItemId != 0 {
			if tooltip, ok := itemTooltipsDB[int(enchant.ItemId)]; ok {
				db.itemIcons = append(db.itemIcons, &proto.IconData{Id: enchant.ItemId, Name: tooltip.GetName(), Icon: tooltip.GetIcon()})
			}
		}
		if enchant.SpellId != 0 {
			if tooltip, ok := spellTooltipsDB[int(enchant.SpellId)]; ok {
				db.spellIcons = append(db.spellIcons, &proto.IconData{Id: enchant.SpellId, Name: tooltip.GetName(), Icon: tooltip.GetIcon()})
			}
		}
	}

	for _, itemID := range extraItemIcons {
		if itemID != 0 {
			if tooltip, ok := itemTooltipsDB[int(itemID)]; ok {
				db.itemIcons = append(db.itemIcons, &proto.IconData{Id: int32(itemID), Name: tooltip.GetName(), Icon: tooltip.GetIcon()})
			}
		}
	}

	db.itemIcons = core.FilterSlice(db.itemIcons, func(icon *proto.IconData) bool {
		return icon.Name != "" && icon.Icon != ""
	})
	db.spellIcons = core.FilterSlice(db.spellIcons, func(icon *proto.IconData) bool {
		return icon.Name != "" && icon.Icon != ""
	})

	slices.SortStableFunc(db.itemIcons, func(s1, s2 *proto.IconData) bool {
		return s1.Id < s2.Id
	})
	slices.SortStableFunc(db.spellIcons, func(s1, s2 *proto.IconData) bool {
		return s1.Id < s2.Id
	})

	db.applyGlobalFilters()

	return db
}

// Filters out entities which shouldn't be included anywhere.
func (db *WowDatabase) applyGlobalFilters() {
	db.items = core.FilterMap(db.items, func(_ int32, item *proto.UIItem) bool {
		if _, ok := itemDenyList[item.Id]; ok {
			return false
		}

		for _, pattern := range denyListNameRegexes {
			if pattern.MatchString(item.Name) {
				return false
			}
		}
		return true
	})

	db.gems = core.FilterMap(db.gems, func(_ int32, gem *proto.UIGem) bool {
		if _, ok := gemDenyList[gem.Id]; ok {
			return false
		}

		for _, pattern := range denyListNameRegexes {
			if pattern.MatchString(gem.Name) {
				return false
			}
		}
		return true
	})
}

// Returns only items which are worth including in the sim.
func (db *WowDatabase) getSimmableItems() map[int32]*proto.UIItem {
	return core.FilterMap(db.items, func(_ int32, item *proto.UIItem) bool {
		if _, ok := itemAllowList[item.Id]; ok {
			return true
		}

		if item.Quality < proto.ItemQuality_ItemQualityUncommon {
			return false
		} else if item.Quality > proto.ItemQuality_ItemQualityLegendary {
			return false
		} else if item.Quality < proto.ItemQuality_ItemQualityEpic {
			if item.Ilvl < 145 {
				return false
			}
			if item.Ilvl < 149 && item.SetName == "" {
				return false
			}
		} else {
			// Epic and legendary items might come from classic, so use a lower ilvl threshold.
			if item.Ilvl < 140 {
				return false
			}
		}
		if item.Ilvl == 0 {
			fmt.Printf("Missing ilvl: %s\n", item.Name)
		}

		return true
	})
}

// Returns only gems which are worth including in the sim.
func (db *WowDatabase) getSimmableGems() map[int32]*proto.UIGem {
	return core.FilterMap(db.gems, func(id int32, gem *proto.UIGem) bool {
		if gem.Quality < proto.ItemQuality_ItemQualityUncommon {
			return false
		}
		return true
	})
}

func (db *WowDatabase) toUIDatabase() *proto.UIDatabase {
	uiDB := &proto.UIDatabase{
		Enchants:   db.enchants,
		Encounters: db.encounters,
		ItemIcons:  db.itemIcons,
		SpellIcons: db.spellIcons,
	}

	for _, item := range db.getSimmableItems() {
		uiDB.Items = append(uiDB.Items, item)
	}
	for _, gem := range db.getSimmableGems() {
		uiDB.Gems = append(uiDB.Gems, gem)
	}

	slices.SortStableFunc(uiDB.Items, func(i1, i2 *proto.UIItem) bool {
		return i1.Id < i2.Id
	})
	slices.SortStableFunc(uiDB.Gems, func(g1, g2 *proto.UIGem) bool {
		return g1.Id < g2.Id
	})
	return uiDB
}

func mergeStats(statlist Stats, overrides Stats) Stats {
	merged := Stats{}
	for stat, value := range statlist {
		val := value
		if overrides[stat] > 0 {
			val = overrides[stat]
		}
		merged[stat] = val
	}
	return merged
}
func toSlice(stats Stats) []float64 {
	return stats[:]
}
