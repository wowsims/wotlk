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
	Filter         bool           // If true, this item will be omitted from the sim.
	Keep           bool           // If true, keep this item even if it would otherwise be filtered.
}

type ItemData struct {
	Response ItemResponse
	Override ItemOverride
}

func (itemData *ItemData) toProto() *proto.UIItem {
	weaponDamageMin, weaponDamageMax := itemData.Response.GetWeaponDamage()

	itemProto := &proto.UIItem{
		Name: itemData.Response.GetName(),

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

// For overriding gem data.
type GemOverride struct {
	ID int

	Stats Stats // Only non-zero values will override
	Phase int

	Filter bool // If true, this item will be omitted from the sim.
}

type GemData struct {
	Response ItemResponse
	Override GemOverride
}

func (gemData *GemData) toProto() *proto.UIGem {
	gemProto := &proto.UIGem{
		Name:  gemData.Response.GetName(),
		Color: gemData.Response.GetSocketColor(),

		Stats: toSlice(mergeStats(gemData.Response.GetGemStats(), gemData.Override.Stats)),

		Phase:              int32(gemData.Response.GetPhase()),
		Quality:            proto.ItemQuality(gemData.Response.GetQuality()),
		Unique:             gemData.Response.GetUnique(),
		RequiredProfession: gemData.Response.GetRequiredProfession(),
	}

	overrideProto := &proto.UIGem{
		Id:    int32(gemData.Override.ID),
		Phase: int32(gemData.Override.Phase),
	}

	googleProto.Merge(gemProto, overrideProto)
	return gemProto
}

type SpellData struct {
	ID       int
	Response ItemResponse
}

type WowDatabase struct {
	items    []ItemData
	enchants []*proto.UIEnchant
	gems     []GemData
	spells   []SpellData

	encounters []*proto.PresetEncounter
}

func NewWowDatabase(itemOverrides []ItemOverride, gemOverrides []GemOverride, enchantOverrides []*proto.UIEnchant, itemTooltipsDB map[int]WowheadItemResponse, spellTooltipsDB map[int]WowheadItemResponse) *WowDatabase {
	db := &WowDatabase{
		enchants:   enchantOverrides,
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

	for _, gemOverride := range gemOverrides {
		gemData := GemData{
			Override: gemOverride,
			Response: itemTooltipsDB[gemOverride.ID],
		}
		if gemData.Response.GetName() == "" {
			continue
		}
		db.gems = append(db.gems, gemData)
	}

	for spellID, spellResponse := range spellTooltipsDB {
		db.spells = append(db.spells, SpellData{
			ID:       spellID,
			Response: spellResponse,
		})
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

	slices.SortStableFunc(db.spells, func(s1, s2 SpellData) bool {
		if s1.Response.GetName() == s2.Response.GetName() {
			return s1.ID < s2.ID
		}
		return s1.Response.GetName() < s2.Response.GetName()
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

func (db *WowDatabase) toUIDatabase() *proto.UIDatabase {
	uiDB := &proto.UIDatabase{
		Enchants:   db.enchants,
		Encounters: db.encounters,
	}

	for _, itemData := range db.getSimmableItems() {
		uiDB.Items = append(uiDB.Items, itemData.toProto())
	}
	for _, gemData := range db.getSimmableGems() {
		uiDB.Gems = append(uiDB.Gems, gemData.toProto())
	}

	for _, itemData := range db.items {
		uiDB.ItemIcons = append(uiDB.ItemIcons, &proto.IconData{Id: int32(itemData.Override.ID), Name: itemData.Response.GetName(), Icon: itemData.Response.GetIcon()})
	}
	for _, gemData := range db.gems {
		uiDB.ItemIcons = append(uiDB.ItemIcons, &proto.IconData{Id: int32(gemData.Override.ID), Name: gemData.Response.GetName(), Icon: gemData.Response.GetIcon()})
	}
	for _, spellData := range db.spells {
		uiDB.SpellIcons = append(uiDB.SpellIcons, &proto.IconData{Id: int32(spellData.ID), Name: spellData.Response.GetName(), Icon: spellData.Response.GetIcon()})
	}
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
