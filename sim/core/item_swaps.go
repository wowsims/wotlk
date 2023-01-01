package core

import (
	"time"

	"github.com/wowsims/wotlk/sim/core/proto"
)

type OnSwapItem func(*Simulation)

var onSwapCallbacks = map[int32]OnSwapItem{}

const offset = proto.ItemSlot_ItemSlotMainHand

type ItemSwap struct {
	character *Character

	//Used for resetting
	initialEquippedItems   [3]Item
	initialUnequippedItems [3]Item

	//holds items that are currently not equipped
	unEquippedItems [3]Item
}

func RegisterOnItemSwap(id int32, callback OnSwapItem) {
	onSwapCallbacks[id] = callback
}

func (character *Character) EnableItemSwap(itemSwap *proto.ItemSwap) {
	items := getItems(itemSwap)
	character.ItemSwap = ItemSwap{
		character:              character,
		initialEquippedItems:   getInitialEquippedItems(character),
		initialUnequippedItems: items,
		unEquippedItems:        items,
	}
}

func (swap *ItemSwap) IsEnabled() bool {
	return swap.character != nil
}

func getInitialEquippedItems(character *Character) [3]Item {
	var items [3]Item

	for i := range items {
		items[i] = character.Equip[i+int(offset)]
	}

	return items
}

func getItems(itemSwap *proto.ItemSwap) [3]Item {
	var items [3]Item

	if itemSwap != nil {
		items[0] = toItem(itemSwap.MhItem)
		items[1] = toItem(itemSwap.OhItem)
		items[2] = toItem(itemSwap.RangedItem)
	}

	return items
}

func toItem(itemSpecProto *proto.ItemSpec) Item {
	if itemSpecProto == nil {
		return Item{}
	}

	return ProtoToItem(itemSpecProto)
}

func (swap *ItemSwap) GetItem(slot proto.ItemSlot) Item {
	if slot-offset < 0 {
		panic("Not able to swap Item " + slot.String() + " not supported")
	}
	return swap.unEquippedItems[slot-offset]
}

func (swap *ItemSwap) setItem(slot proto.ItemSlot, item Item) {
	swap.unEquippedItems[slot-offset] = item
}

func (swap *ItemSwap) SwapItems(sim *Simulation, slots []proto.ItemSlot, useGCD bool) {
	if !swap.IsEnabled() {
		return
	}

	character := swap.character

	meeleWeaponSwapped := false
	for _, slot := range slots {
		if swap.swapItem(sim, slot) {
			meeleWeaponSwapped = slot == proto.ItemSlot_ItemSlotMainHand || slot == proto.ItemSlot_ItemSlotOffHand
		}
	}

	for _, onSwap := range onSwapCallbacks {
		onSwap(sim)
	}

	if character.AutoAttacks.IsEnabled() && meeleWeaponSwapped {
		character.AutoAttacks.StopMeleeUntil(sim, sim.CurrentTime, false)
	}

	if useGCD {
		character.GCD.Set(1500 * time.Millisecond)
	}
}

func (swap *ItemSwap) swapItem(sim *Simulation, slot proto.ItemSlot) bool {
	character := swap.character
	oldItem := character.Equip[slot]
	newItem := swap.GetItem(slot)

	// No item to swap too
	if newItem.ID == 0 {
		return false
	}

	character.Equip[slot] = newItem
	stats := newItem.Stats.Add(oldItem.Stats.Multiply(-1))
	character.AddStatsDynamic(sim, stats)

	swap.setItem(slot, oldItem)
	swap.swapWeapon(sim, slot)

	return true
}

func (swap *ItemSwap) swapWeapon(sim *Simulation, slot proto.ItemSlot) {
	character := swap.character
	if !character.AutoAttacks.IsEnabled() {
		return
	}

	switch slot {
	case proto.ItemSlot_ItemSlotMainHand:
		character.AutoAttacks.MH = character.WeaponFromMainHand(character.AutoAttacks.MH.CritMultiplier)
		break
	case proto.ItemSlot_ItemSlotOffHand:
		character.AutoAttacks.OH = character.WeaponFromOffHand(character.AutoAttacks.OH.CritMultiplier)
		break
	case proto.ItemSlot_ItemSlotRanged:
		character.AutoAttacks.Ranged = character.WeaponFromRanged(character.AutoAttacks.Ranged.CritMultiplier)
		break
	}

	character.AutoAttacks.IsDualWielding = character.Equip[proto.ItemSlot_ItemSlotMainHand].SwingSpeed != 0 && character.Equip[proto.ItemSlot_ItemSlotOffHand].SwingSpeed != 0
}

func (swap *ItemSwap) reset(sim *Simulation) {
	if !swap.IsEnabled() {
		return
	}

	character := swap.character
	slots := [3]proto.ItemSlot{proto.ItemSlot_ItemSlotMainHand, proto.ItemSlot_ItemSlotOffHand, proto.ItemSlot_ItemSlotRanged}
	for i, slot := range slots {
		character.Equip[slot] = swap.initialEquippedItems[i]
		swap.swapWeapon(sim, slot)
	}

	swap.unEquippedItems = swap.initialUnequippedItems

	for _, onSwap := range onSwapCallbacks {
		onSwap(sim)
	}
}
