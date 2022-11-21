package main

import (
	"fmt"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	_ "github.com/wowsims/wotlk/sim/encounters" // Needed for preset encounters.
	"golang.org/x/exp/slices"
	googleProto "google.golang.org/protobuf/proto"
)

type EnchantDBKey struct {
	EffectID int32
	ItemID   int32
	SpellID  int32
}

func EnchantToDBKey(enchant *proto.UIEnchant) EnchantDBKey {
	return EnchantDBKey{
		EffectID: enchant.EffectId,
		ItemID:   enchant.ItemId,
		SpellID:  enchant.SpellId,
	}
}

func mergeItemProtos(dst, src *proto.UIItem) {
	// googleproto.Merge concatenates lists but we want replacement, so do them manually.
	if src.Stats != nil {
		dst.Stats = src.Stats
		src.Stats = nil
	}
	if src.SocketBonus != nil {
		dst.SocketBonus = src.SocketBonus
		src.SocketBonus = nil
	}
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
	enchants map[EnchantDBKey]*proto.UIEnchant
	gems     map[int32]*proto.UIGem

	itemIcons  map[int32]*proto.IconData
	spellIcons map[int32]*proto.IconData

	encounters []*proto.PresetEncounter
}

func NewWowDatabase(itemOverrides []*proto.UIItem, gemOverrides []*proto.UIGem, enchantOverrides []*proto.UIEnchant, itemTooltipsDB map[int32]WowheadItemResponse, spellTooltipsDB map[int32]WowheadItemResponse) *WowDatabase {
	db := &WowDatabase{
		items:    make(map[int32]*proto.UIItem),
		enchants: make(map[EnchantDBKey]*proto.UIEnchant),
		gems:     make(map[int32]*proto.UIGem),

		itemIcons:  make(map[int32]*proto.IconData),
		spellIcons: make(map[int32]*proto.IconData),
		encounters: core.PresetEncounters,
	}

	for id, response := range itemTooltipsDB {
		if response.IsEquippable() {
			itemProto := response.ToItemProto()
			itemProto.Id = id
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
			gemProto.Id = id
			db.gems[gemProto.Id] = gemProto
		}
	}
	for _, gemOverride := range gemOverrides {
		if _, ok := db.gems[gemOverride.Id]; ok {
			mergeGemProtos(db.gems[gemOverride.Id], gemOverride)
		}
	}

	for _, enchant := range enchantOverrides {
		db.enchants[EnchantToDBKey(enchant)] = enchant
	}
	for _, enchant := range db.enchants {
		if enchant.ItemId != 0 {
			if tooltip, ok := itemTooltipsDB[enchant.ItemId]; ok {
				db.itemIcons[enchant.ItemId] = &proto.IconData{Id: enchant.ItemId, Name: tooltip.GetName(), Icon: tooltip.GetIcon()}
			}
		}
		if enchant.SpellId != 0 {
			if tooltip, ok := spellTooltipsDB[enchant.SpellId]; ok {
				db.spellIcons[enchant.SpellId] = &proto.IconData{Id: enchant.SpellId, Name: tooltip.GetName(), Icon: tooltip.GetIcon()}
			}
		}
	}

	for _, itemID := range extraItemIcons {
		if itemID != 0 {
			if tooltip, ok := itemTooltipsDB[itemID]; ok {
				db.itemIcons[itemID] = &proto.IconData{Id: itemID, Name: tooltip.GetName(), Icon: tooltip.GetIcon()}
			}
		}
	}

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

	db.itemIcons = core.FilterMap(db.itemIcons, func(_ int32, icon *proto.IconData) bool {
		return icon.Name != "" && icon.Icon != ""
	})
	db.spellIcons = core.FilterMap(db.spellIcons, func(_ int32, icon *proto.IconData) bool {
		return icon.Name != "" && icon.Icon != ""
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
		Encounters: db.encounters,
	}

	for _, v := range db.getSimmableItems() {
		uiDB.Items = append(uiDB.Items, v)
	}
	for _, v := range db.getSimmableGems() {
		uiDB.Gems = append(uiDB.Gems, v)
	}
	for _, v := range db.enchants {
		uiDB.Enchants = append(uiDB.Enchants, v)
	}
	for _, v := range db.itemIcons {
		uiDB.ItemIcons = append(uiDB.ItemIcons, v)
	}
	for _, v := range db.spellIcons {
		uiDB.SpellIcons = append(uiDB.SpellIcons, v)
	}

	slices.SortStableFunc(uiDB.Items, func(v1, v2 *proto.UIItem) bool {
		return v1.Id < v2.Id
	})
	slices.SortStableFunc(uiDB.Enchants, func(v1, v2 *proto.UIEnchant) bool {
		return v1.EffectId < v2.EffectId || v1.EffectId == v2.EffectId && v1.Type < v2.Type
	})
	slices.SortStableFunc(uiDB.Gems, func(v1, v2 *proto.UIGem) bool {
		return v1.Id < v2.Id
	})
	slices.SortStableFunc(uiDB.ItemIcons, func(v1, v2 *proto.IconData) bool {
		return v1.Id < v2.Id
	})
	slices.SortStableFunc(uiDB.SpellIcons, func(v1, v2 *proto.IconData) bool {
		return v1.Id < v2.Id
	})

	return uiDB
}

func toSlice(stats Stats) []float64 {
	return stats[:]
}
