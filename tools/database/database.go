package database

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	_ "github.com/wowsims/wotlk/sim/encounters" // Needed for preset encounters.
	"github.com/wowsims/wotlk/tools"
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

	for id, response := range itemTooltipsDB {
		if response.IsGem() {
			gemProto := response.ToGemProto()
			gemProto.Id = id
			db.gems[gemProto.Id] = gemProto
		}
	}

	db.MergeItems(itemOverrides)
	db.MergeGems(gemOverrides)
	db.MergeEnchants(enchantOverrides)

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
		if tooltip, ok := itemTooltipsDB[itemID]; ok {
			db.itemIcons[itemID] = &proto.IconData{Id: itemID, Name: tooltip.GetName(), Icon: tooltip.GetIcon()}
		}
		//if item, ok := db.items[itemID]; ok {
		//	db.itemIcons[itemID] = &proto.IconData{Id: itemID, Name: item.Name, Icon: item.Icon}
		//}
	}

	return db
}

func (db *WowDatabase) MergeItems(arr []*proto.UIItem) {
	for _, item := range arr {
		db.MergeItem(item)
	}
}
func (db *WowDatabase) MergeItem(newItem *proto.UIItem) {
	if curItem, ok := db.items[newItem.Id]; ok {
		mergeItemProtos(curItem, newItem)
	} else {
		db.items[newItem.Id] = newItem
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

func (db *WowDatabase) MergeEnchants(arr []*proto.UIEnchant) {
	for _, enchant := range arr {
		db.MergeEnchant(enchant)
	}
}
func (db *WowDatabase) MergeEnchant(newEnchant *proto.UIEnchant) {
	key := EnchantToDBKey(newEnchant)
	if curEnchant, ok := db.enchants[key]; ok {
		mergeEnchantProtos(curEnchant, newEnchant)
	} else {
		db.enchants[key] = newEnchant
	}
}
func mergeEnchantProtos(dst, src *proto.UIEnchant) {
	// googleproto.Merge concatenates lists but we want replacement, so do them manually.
	if src.Stats != nil {
		dst.Stats = src.Stats
		src.Stats = nil
	}
	googleProto.Merge(dst, src)
}

func (db *WowDatabase) MergeGems(arr []*proto.UIGem) {
	for _, gem := range arr {
		db.MergeGem(gem)
	}
}
func (db *WowDatabase) MergeGem(newGem *proto.UIGem) {
	if curGem, ok := db.gems[newGem.Id]; ok {
		mergeGemProtos(curGem, newGem)
	} else {
		db.gems[newGem.Id] = newGem
	}
}
func mergeGemProtos(dst, src *proto.UIGem) {
	// googleproto.Merge concatenates lists but we want replacement, so do them manually.
	if src.Stats != nil {
		dst.Stats = src.Stats
		src.Stats = nil
	}
	googleProto.Merge(dst, src)
}

// Filters out entities which shouldn't be included anywhere.
func (db *WowDatabase) ApplyGlobalFilters() {
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

func (db *WowDatabase) WriteBinaryAndJson(binFilePath, jsonFilePath string) {
	uiDB := db.toUIDatabase()

	// Write database as a binary file.
	outbytes, err := googleProto.Marshal(uiDB)
	if err != nil {
		log.Fatalf("[ERROR] Failed to marshal db: %s", err.Error())
	}
	os.WriteFile(binFilePath, outbytes, 0666)

	// Also write in JSON format so we can manually inspect the contents.
	// Write it out line-by-line so we can have 1 line / item, making it more human-readable.
	builder := &strings.Builder{}
	builder.WriteString("{\n")

	tools.WriteProtoArrayToBuilder(uiDB.Items, builder, "items")
	builder.WriteString(",\n")
	tools.WriteProtoArrayToBuilder(uiDB.Enchants, builder, "enchants")
	builder.WriteString(",\n")
	tools.WriteProtoArrayToBuilder(uiDB.Gems, builder, "gems")
	builder.WriteString(",\n")
	tools.WriteProtoArrayToBuilder(uiDB.ItemIcons, builder, "itemIcons")
	builder.WriteString(",\n")
	tools.WriteProtoArrayToBuilder(uiDB.SpellIcons, builder, "spellIcons")
	builder.WriteString(",\n")
	tools.WriteProtoArrayToBuilder(uiDB.Encounters, builder, "encounters")
	builder.WriteString("\n")

	builder.WriteString("}")
	os.WriteFile(jsonFilePath, []byte(builder.String()), 0666)
}

func toSlice(stats Stats) []float64 {
	return stats[:]
}
