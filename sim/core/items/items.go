package items

import (
	"fmt"

	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
	"google.golang.org/protobuf/encoding/protojson"
)

var ByName = map[string]Item{}
var ByID = map[int32]Item{}
var GemsByName = map[string]Gem{}
var GemsByID = map[int32]Gem{}
var EnchantsByName = map[string]Enchant{}
var EnchantsByItemByID = map[proto.ItemType]map[int32]Enchant{}

func init() {
	for _, v := range Enchants {
		EnchantsByName[v.Name] = v
		if EnchantsByItemByID[v.ItemType] == nil {
			EnchantsByItemByID[v.ItemType] = map[int32]Enchant{}
		}
		EnchantsByItemByID[v.ItemType][v.ID] = v
	}
	for _, v := range Gems {
		GemsByName[v.Name] = v
		GemsByID[v.ID] = v
	}

	// Add hard-coded items. Wowhead doesn't seem to have tooltips for random enchant items.
	// Use negative IDs to avoid collisions with real item IDs.
	Items = append(Items, []Item{
		// Example hard-coded item using a negative ID.
		// {Name: "Glider's Boots of Nature's Wrath", WowheadID: 30681, ID: -1, Type: proto.ItemType_ItemTypeFeet, ArmorType: proto.ArmorType_ArmorTypeLeather, Phase: 1, Quality: proto.ItemQuality_ItemQualityEpic, Stats: stats.Stats{stats.Armor: 250, stats.NatureSpellPower: 78}},
	}...)

	for _, v := range Items {
		if _, ok := ByID[v.ID]; ok {
			fmt.Printf("Found dup item: %s\n", v.Name)
			panic("no dupes allowed")
		}
		ByName[v.Name] = v
		ByID[v.ID] = v
	}
}

type Item struct {
	ID        int32
	WowheadID int32
	Heroic    bool
	Type      proto.ItemType
	ArmorType proto.ArmorType
	// Weapon Stats
	WeaponType       proto.WeaponType
	HandType         proto.HandType
	RangedWeaponType proto.RangedWeaponType
	WeaponDamageMin  float64
	WeaponDamageMax  float64
	SwingSpeed       float64

	// Used by the UI to filter which items are shown.
	ClassAllowlist []proto.Class

	Name       string
	SourceZone string
	SourceDrop string
	Stats      stats.Stats // Stats applied to wearer
	Phase      byte
	Quality    proto.ItemQuality
	Unique     bool
	Ilvl       int32
	SetName    string // Empty string if not part of a set.

	RequiredProfession proto.Profession

	// Hidden variable used for a few obscure mechanics (Seal of Righteousness).
	// Intuitively, this is a measure of the difference between the expected stats
	// and the actual stats of an item, e.g. decreased weapon DPS on caster weapons.
	QualityModifier float64

	GemSockets  []proto.GemColor
	SocketBonus stats.Stats

	// Modified for each instance of the item.
	Gems    []Gem
	Enchant Enchant
}

func (item Item) ToProto() *proto.Item {
	return &proto.Item{
		Id:                 item.ID,
		WowheadId:          item.WowheadID,
		Name:               item.Name,
		ClassAllowlist:     item.ClassAllowlist[:],
		Type:               proto.ItemType(item.Type),
		ArmorType:          proto.ArmorType(item.ArmorType),
		WeaponType:         proto.WeaponType(item.WeaponType),
		HandType:           proto.HandType(item.HandType),
		RangedWeaponType:   proto.RangedWeaponType(item.RangedWeaponType),
		WeaponDamageMin:    item.WeaponDamageMin,
		WeaponDamageMax:    item.WeaponDamageMax,
		WeaponSpeed:        item.SwingSpeed,
		Stats:              item.Stats[:],
		Phase:              int32(item.Phase),
		Quality:            item.Quality,
		Unique:             item.Unique,
		Ilvl:               item.Ilvl,
		GemSockets:         item.GemSockets,
		SocketBonus:        item.SocketBonus[:],
		RequiredProfession: item.RequiredProfession,
		Heroic:             item.Heroic,
	}
}

func (item Item) ToItemSpecProto() *proto.ItemSpec {
	itemSpec := &proto.ItemSpec{
		Id:      item.ID,
		Enchant: item.Enchant.ID,
		Gems:    []int32{},
	}
	for _, gem := range item.Gems {
		itemSpec.Gems = append(itemSpec.Gems, gem.ID)
	}
	return itemSpec
}

type Enchant struct {
	ID          int32 // ID of the enchant item.
	EffectID    int32 // Used by UI to apply effect to tooltip
	Name        string
	IsSpellID   bool
	Quality     proto.ItemQuality
	Bonus       stats.Stats
	ItemType    proto.ItemType    // Which slot the enchant goes on.
	EnchantType proto.EnchantType // Additional category when ItemType isn't enough.
	Phase       int32

	RequiredProfession proto.Profession

	// Used by the UI to filter which enchants are shown.
	ClassAllowlist []proto.Class
}

func (enchant Enchant) ToProto() *proto.Enchant {
	return &proto.Enchant{
		Id:             enchant.ID,
		EffectId:       enchant.EffectID,
		Name:           enchant.Name,
		IsSpellId:      enchant.IsSpellID,
		Type:           enchant.ItemType,
		EnchantType:    enchant.EnchantType,
		Stats:          enchant.Bonus[:],
		Quality:        enchant.Quality,
		Phase:          enchant.Phase,
		ClassAllowlist: enchant.ClassAllowlist[:],

		RequiredProfession: enchant.RequiredProfession,
	}
}

type Gem struct {
	ID      int32
	Name    string
	Stats   stats.Stats // flat stats gem adds
	Color   proto.GemColor
	Phase   byte
	Quality proto.ItemQuality
	Unique  bool

	RequiredProfession proto.Profession
}

func (gem Gem) ToProto() *proto.Gem {
	return &proto.Gem{
		Id:      gem.ID,
		Name:    gem.Name,
		Stats:   gem.Stats[:],
		Color:   gem.Color,
		Phase:   int32(gem.Phase),
		Quality: gem.Quality,
		Unique:  gem.Unique,

		RequiredProfession: gem.RequiredProfession,
	}
}

type ItemSpec struct {
	ID      int32
	Enchant int32
	Gems    []int32
}

type Equipment [proto.ItemSlot_ItemSlotRanged + 1]Item

func (equipment *Equipment) EquipItem(item Item) {
	if item.Type == proto.ItemType_ItemTypeFinger {
		if equipment[ItemSlotFinger1].Name == "" {
			equipment[ItemSlotFinger1] = item
		} else {
			equipment[ItemSlotFinger2] = item
		}
	} else if item.Type == proto.ItemType_ItemTypeTrinket {
		if equipment[ItemSlotTrinket1].Name == "" {
			equipment[ItemSlotTrinket1] = item
		} else {
			equipment[ItemSlotTrinket2] = item
		}
	} else if item.Type == proto.ItemType_ItemTypeWeapon {
		if item.WeaponType == proto.WeaponType_WeaponTypeShield && equipment[ItemSlotMainHand].HandType != proto.HandType_HandTypeTwoHand {
			equipment[ItemSlotOffHand] = item
		} else if item.HandType == proto.HandType_HandTypeMainHand || item.HandType == proto.HandType_HandTypeUnknown {
			equipment[ItemSlotMainHand] = item
		} else if item.HandType == proto.HandType_HandTypeOffHand {
			equipment[ItemSlotOffHand] = item
		} else if item.HandType == proto.HandType_HandTypeOneHand || item.HandType == proto.HandType_HandTypeTwoHand {
			if equipment[ItemSlotMainHand].ID == 0 {
				equipment[ItemSlotMainHand] = item
			} else if equipment[ItemSlotOffHand].ID == 0 {
				equipment[ItemSlotOffHand] = item
			}
		}
	} else {
		equipment[ItemTypeToSlot(item.Type)] = item
	}
}

func (equipment *Equipment) ToEquipmentSpecProto() *proto.EquipmentSpec {
	equipSpec := &proto.EquipmentSpec{
		Items: []*proto.ItemSpec{},
	}
	for _, item := range equipment {
		equipSpec.Items = append(equipSpec.Items, item.ToItemSpecProto())
	}
	return equipSpec
}

// Structs used for looking up items/gems/enchants
type EquipmentSpec [proto.ItemSlot_ItemSlotRanged + 1]ItemSpec

func ProtoToEquipmentSpec(es *proto.EquipmentSpec) EquipmentSpec {
	coreEquip := EquipmentSpec{}

	for i, item := range es.Items {
		spec := ItemSpec{
			ID: item.Id,
		}
		spec.Gems = item.Gems
		spec.Enchant = item.Enchant
		coreEquip[i] = spec
	}

	return coreEquip
}

func NewItem(itemSpec ItemSpec) Item {
	item := Item{}
	if foundItem, ok := ByID[itemSpec.ID]; ok {
		item = foundItem
	} else {
		panic(fmt.Sprintf("No item with id: %d", itemSpec.ID))
	}

	if itemSpec.Enchant != 0 {
		if enchant, ok := EnchantsByItemByID[item.Type][itemSpec.Enchant]; ok {
			item.Enchant = enchant
		} else {
			panic(fmt.Sprintf("No enchant with id: %d", itemSpec.Enchant))
		}
	}

	if len(itemSpec.Gems) > 0 {
		// Need to do this to account for possible extra gem sockets.
		numGems := len(item.GemSockets)
		if len(itemSpec.Gems) > numGems {
			numGems = len(itemSpec.Gems)
		}

		item.Gems = make([]Gem, numGems)
		for gemIdx, gemID := range itemSpec.Gems {
			if gem, ok := GemsByID[gemID]; ok {
				item.Gems[gemIdx] = gem
			} else {
				if gemID != 0 {
					panic(fmt.Sprintf("No gem with id: %d", gemID))
				}
			}
		}
	}
	return item
}

func NewEquipmentSet(equipSpec EquipmentSpec) Equipment {
	equipment := Equipment{}
	for _, itemSpec := range equipSpec {
		if itemSpec.ID != 0 {
			equipment.EquipItem(NewItem(itemSpec))
		}
	}
	return equipment
}

func ProtoToEquipment(es *proto.EquipmentSpec) Equipment {
	return NewEquipmentSet(ProtoToEquipmentSpec(es))
}

// Like ItemSpec, but uses names for reference instead of ID.
type ItemStringSpec struct {
	Name    string
	Enchant string
	Gems    []string
}

func EquipmentSpecFromJsonString(jsonString string) *proto.EquipmentSpec {
	es := &proto.EquipmentSpec{}

	data := []byte(jsonString)
	if err := protojson.Unmarshal(data, es); err != nil {
		panic(err)
	}
	return es
}

func (equipment Equipment) Clone() Equipment {
	newEquipment := Equipment{}
	for idx, item := range equipment {
		newItem := item
		newEquipment[idx] = newItem
	}
	return newEquipment
}

func (equipment Equipment) Stats() stats.Stats {
	equipStats := stats.Stats{}
	for _, item := range equipment {
		equipStats = equipStats.Add(item.Stats)
		equipStats = equipStats.Add(item.Enchant.Bonus)

		for _, gem := range item.Gems {
			equipStats = equipStats.Add(gem.Stats)
		}

		// Check socket bonus
		if len(item.GemSockets) > 0 && len(item.Gems) >= len(item.GemSockets) {
			allMatch := true
			for gemIndex, socketColor := range item.GemSockets {
				if !ColorIntersects(socketColor, item.Gems[gemIndex].Color) {
					allMatch = false
					break
				}
			}

			if allMatch {
				equipStats = equipStats.Add(item.SocketBonus)
			}
		}
	}
	return equipStats
}

type ItemSlot byte

const (
	ItemSlotHead ItemSlot = iota
	ItemSlotNeck
	ItemSlotShoulder
	ItemSlotBack
	ItemSlotChest
	ItemSlotWrist
	ItemSlotHands
	ItemSlotWaist
	ItemSlotLegs
	ItemSlotFeet
	ItemSlotFinger1
	ItemSlotFinger2
	ItemSlotTrinket1
	ItemSlotTrinket2
	ItemSlotMainHand // can be 1h or 2h
	ItemSlotOffHand
	ItemSlotRanged
)

func ItemTypeToSlot(it proto.ItemType) ItemSlot {
	switch it {
	case proto.ItemType_ItemTypeHead:
		return ItemSlotHead
	case proto.ItemType_ItemTypeNeck:
		return ItemSlotNeck
	case proto.ItemType_ItemTypeShoulder:
		return ItemSlotShoulder
	case proto.ItemType_ItemTypeBack:
		return ItemSlotBack
	case proto.ItemType_ItemTypeChest:
		return ItemSlotChest
	case proto.ItemType_ItemTypeWrist:
		return ItemSlotWrist
	case proto.ItemType_ItemTypeHands:
		return ItemSlotHands
	case proto.ItemType_ItemTypeWaist:
		return ItemSlotWaist
	case proto.ItemType_ItemTypeLegs:
		return ItemSlotLegs
	case proto.ItemType_ItemTypeFeet:
		return ItemSlotFeet
	case proto.ItemType_ItemTypeFinger:
		return ItemSlotFinger1
	case proto.ItemType_ItemTypeTrinket:
		return ItemSlotTrinket1
	case proto.ItemType_ItemTypeWeapon:
		return ItemSlotMainHand
	case proto.ItemType_ItemTypeRanged:
		return ItemSlotRanged
	}

	return 255
}

func ColorIntersects(g proto.GemColor, o proto.GemColor) bool {
	if g == o {
		return true
	}
	if g == proto.GemColor_GemColorPrismatic || o == proto.GemColor_GemColorPrismatic {
		return true
	}
	if g == proto.GemColor_GemColorMeta {
		return false // meta gems o nothing.
	}
	if g == proto.GemColor_GemColorRed {
		return o == proto.GemColor_GemColorOrange || o == proto.GemColor_GemColorPurple
	}
	if g == proto.GemColor_GemColorBlue {
		return o == proto.GemColor_GemColorGreen || o == proto.GemColor_GemColorPurple
	}
	if g == proto.GemColor_GemColorYellow {
		return o == proto.GemColor_GemColorGreen || o == proto.GemColor_GemColorOrange
	}
	if g == proto.GemColor_GemColorOrange {
		return o == proto.GemColor_GemColorYellow || o == proto.GemColor_GemColorRed
	}
	if g == proto.GemColor_GemColorGreen {
		return o == proto.GemColor_GemColorYellow || o == proto.GemColor_GemColorBlue
	}
	if g == proto.GemColor_GemColorPurple {
		return o == proto.GemColor_GemColorBlue || o == proto.GemColor_GemColorRed
	}

	return false // dunno what else could be.
}
