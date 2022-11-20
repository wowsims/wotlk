package main

import (
	"fmt"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	_ "github.com/wowsims/wotlk/sim/encounters" // Needed for preset encounters.
	"golang.org/x/exp/slices"
	googleProto "google.golang.org/protobuf/proto"
)

// For overriding item data.
type ItemOverride struct {
	ID int

	Stats          Stats // Only non-zero values will override
	ClassAllowlist []proto.Class
	Phase          int
	HandType       proto.HandType // Overrides hand type.
}

type ItemData struct {
	Response ItemResponse
	Override ItemOverride
}

func (itemData *ItemData) toProto() *proto.UIItem {
	weaponDamageMin, weaponDamageMax := itemData.Response.GetWeaponDamage()

	itemProto := &proto.UIItem{
		Name: itemData.Response.GetName(),
		Icon: itemData.Response.GetIcon(),

		Type:             itemData.Response.GetItemType(),
		ArmorType:        itemData.Response.GetArmorType(),
		WeaponType:       itemData.Response.GetWeaponType(),
		HandType:         itemData.Response.GetHandType(),
		RangedWeaponType: itemData.Response.GetRangedWeaponType(),

		Stats:       toSlice(mergeStats(itemData.Response.GetStats(), itemData.Override.Stats)),
		GemSockets:  itemData.Response.GetGemSockets(),
		SocketBonus: toSlice(itemData.Response.GetSocketBonus()),

		WeaponDamageMin: weaponDamageMin,
		WeaponDamageMax: weaponDamageMax,
		WeaponSpeed:     itemData.Response.GetWeaponSpeed(),

		Ilvl:    int32(itemData.Response.GetItemLevel()),
		Phase:   int32(itemData.Response.GetPhase()),
		Quality: proto.ItemQuality(itemData.Response.GetQuality()),
		Unique:  itemData.Response.GetUnique(),
		Heroic:  itemData.Response.IsHeroic(),

		ClassAllowlist:     itemData.Response.GetClassAllowlist(),
		RequiredProfession: itemData.Response.GetRequiredProfession(),
		SetName:            itemData.Response.GetItemSetName(),
	}

	overrideProto := &proto.UIItem{
		Id:    int32(itemData.Override.ID),
		Phase: int32(itemData.Override.Phase),
	}

	googleProto.Merge(itemProto, overrideProto)
	return itemProto
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
	items    []ItemData
	enchants []*proto.UIEnchant
	gems     map[int32]*proto.UIGem

	itemIcons  []*proto.IconData
	spellIcons []*proto.IconData

	encounters []*proto.PresetEncounter
}

func NewWowDatabase(itemOverrides []ItemOverride, gemOverrides []*proto.UIGem, enchantOverrides []*proto.UIEnchant, itemTooltipsDB map[int]WowheadItemResponse, spellTooltipsDB map[int]WowheadItemResponse) *WowDatabase {
	db := &WowDatabase{
		enchants:   enchantOverrides,
		gems:       make(map[int32]*proto.UIGem),
		encounters: core.PresetEncounters,
	}

	for _, itemOverride := range itemOverrides {
		itemData := ItemData{
			Override: itemOverride,
			Response: itemTooltipsDB[itemOverride.ID],
		}
		if itemData.Response.GetName() == "" {
			continue
		}
		db.items = append(db.items, itemData)
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
			if tooltip, ok := itemTooltipsDB[itemID]; ok {
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

	slices.SortStableFunc(db.items, func(i1, i2 ItemData) bool {
		return i1.Override.ID < i2.Override.ID
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
	db.items = core.FilterSlice(db.items, func(itemData ItemData) bool {
		if _, ok := itemDenyList[itemData.Override.ID]; ok {
			return false
		}

		for _, pattern := range denyListNameRegexes {
			if pattern.MatchString(itemData.Response.GetName()) {
				return false
			}
		}
		return true
	})

	db.gems = core.FilterMap(db.gems, func(_ int32, gem *proto.UIGem) bool {
		if _, ok := gemDenyList[int(gem.Id)]; ok {
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
func (db *WowDatabase) getSimmableItems() []ItemData {
	var included []ItemData
	for _, itemData := range db.items {
		if !itemData.Response.IsEquippable() {
			continue
		}

		if _, ok := itemAllowList[itemData.Override.ID]; ok {
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

	for _, itemData := range db.getSimmableItems() {
		uiDB.Items = append(uiDB.Items, itemData.toProto())
	}
	for _, gem := range db.getSimmableGems() {
		uiDB.Gems = append(uiDB.Gems, gem)
	}
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
