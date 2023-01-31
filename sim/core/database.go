package core

import (
	"fmt"

	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
	"google.golang.org/protobuf/encoding/protojson"
)

var WITH_DB = false

var ItemsByID = map[int32]Item{}
var GemsByID = map[int32]Gem{}
var EnchantsByEffectID = map[int32]Enchant{}

func addToDatabase(newDB *proto.SimDatabase) {
	for _, v := range newDB.Items {
		if _, ok := ItemsByID[v.Id]; !ok {
			item := ItemFromProto(v)
			ItemsByID[v.Id] = item
			AddItemToSets(item)
		}
	}

	for _, v := range newDB.Enchants {
		if _, ok := EnchantsByEffectID[v.EffectId]; !ok {
			EnchantsByEffectID[v.EffectId] = EnchantFromProto(v)
		}
	}

	for _, v := range newDB.Gems {
		if _, ok := GemsByID[v.Id]; !ok {
			GemsByID[v.Id] = GemFromProto(v)
		}
	}
}

type Item struct {
	ID        int32
	Type      proto.ItemType
	ArmorType proto.ArmorType
	// Weapon Stats
	WeaponType       proto.WeaponType
	HandType         proto.HandType
	RangedWeaponType proto.RangedWeaponType
	WeaponDamageMin  float64
	WeaponDamageMax  float64
	SwingSpeed       float64

	Name    string
	Stats   stats.Stats // Stats applied to wearer
	Quality proto.ItemQuality
	SetName string // Empty string if not part of a set.

	GemSockets  []proto.GemColor
	SocketBonus stats.Stats

	// Modified for each instance of the item.
	Gems    []Gem
	Enchant Enchant

	//Internal use
	TempEnchant int32
}

func ItemFromProto(pData *proto.SimItem) Item {
	return Item{
		ID:               pData.Id,
		Name:             pData.Name,
		Type:             pData.Type,
		ArmorType:        pData.ArmorType,
		WeaponType:       pData.WeaponType,
		HandType:         pData.HandType,
		RangedWeaponType: pData.RangedWeaponType,
		WeaponDamageMin:  pData.WeaponDamageMin,
		WeaponDamageMax:  pData.WeaponDamageMax,
		SwingSpeed:       pData.WeaponSpeed,
		Stats:            stats.FromFloatArray(pData.Stats),
		GemSockets:       pData.GemSockets,
		SocketBonus:      stats.FromFloatArray(pData.SocketBonus),
		SetName:          pData.SetName,
	}
}

func (item Item) ToItemSpecProto() *proto.ItemSpec {
	itemSpec := &proto.ItemSpec{
		Id:      item.ID,
		Enchant: item.Enchant.EffectID,
		Gems:    []int32{},
	}
	for _, gem := range item.Gems {
		itemSpec.Gems = append(itemSpec.Gems, gem.ID)
	}
	return itemSpec
}

type Enchant struct {
	EffectID int32 // Used by UI to apply effect to tooltip
	Stats    stats.Stats
}

func EnchantFromProto(pData *proto.SimEnchant) Enchant {
	return Enchant{
		EffectID: pData.EffectId,
		Stats:    stats.FromFloatArray(pData.Stats),
	}
}

type Gem struct {
	ID    int32
	Name  string
	Stats stats.Stats
	Color proto.GemColor
}

func GemFromProto(pData *proto.SimGem) Gem {
	return Gem{
		ID:    pData.Id,
		Name:  pData.Name,
		Stats: stats.FromFloatArray(pData.Stats),
		Color: pData.Color,
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
		if equipment[ItemSlotFinger1].ID == 0 {
			equipment[ItemSlotFinger1] = item
		} else {
			equipment[ItemSlotFinger2] = item
		}
	} else if item.Type == proto.ItemType_ItemTypeTrinket {
		if equipment[ItemSlotTrinket1].ID == 0 {
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
	if foundItem, ok := ItemsByID[itemSpec.ID]; ok {
		item = foundItem
	} else {
		panic(fmt.Sprintf("No item with id: %d", itemSpec.ID))
	}

	if itemSpec.Enchant != 0 {
		if enchant, ok := EnchantsByEffectID[itemSpec.Enchant]; ok {
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
		equipStats = equipStats.Add(item.Enchant.Stats)

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

func (is ItemSlot) String() string {
	switch is {
	case ItemSlotHead:
		return "Head"
	case ItemSlotNeck:
		return "Neck"
	case ItemSlotShoulder:
		return "Shoulder"
	case ItemSlotBack:
		return "Back"
	case ItemSlotChest:
		return "Chest"
	case ItemSlotWrist:
		return "Wrist"
	case ItemSlotHands:
		return "Hands"
	case ItemSlotWaist:
		return "Waist"
	case ItemSlotLegs:
		return "Legs"
	case ItemSlotFeet:
		return "Feet"
	case ItemSlotFinger1:
		return "Finger1"
	case ItemSlotFinger2:
		return "Finger2"
	case ItemSlotTrinket1:
		return "Trinket1"
	case ItemSlotTrinket2:
		return "Trinket2"
	case ItemSlotMainHand:
		return "MainHand"
	case ItemSlotOffHand:
		return "OffHand"
	case ItemSlotRanged:
		return "Ranged"
	}
	return "unknown slot"
}

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
