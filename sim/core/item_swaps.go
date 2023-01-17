package core

import (
	"time"

	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

type OnSwapItem func(*Simulation)

const offset = proto.ItemSlot_ItemSlotMainHand

type ItemSwap struct {
	character       *Character
	onSwapCallbacks []OnSwapItem

	mhCritMultiplier     float64
	ohCritMultiplier     float64
	rangedCritMultiplier float64

	//Used for resetting
	initialEquippedItems   [3]Item
	initialUnequippedItems [3]Item

	//holds items that are currently not equipped
	unEquippedItems [3]Item

	//for handling 2Handers, if the unEquippedItems holds a 2Hander
	has2H bool
}

/*
	TODO All the extra parameters here and the code in multiple places for handling the Weapon struct is really messy,
		we'll need to figure out something cleaner as this will be quite error-prone
*/
func (character *Character) EnableItemSwap(itemSwap *proto.ItemSwap, mhCritMultiplier float64, ohCritMultiplier float64, rangedCritMultiplier float64) {
	items := getItems(itemSwap)

	character.ItemSwap = ItemSwap{
		character:            character,
		mhCritMultiplier:     mhCritMultiplier,
		ohCritMultiplier:     ohCritMultiplier,
		rangedCritMultiplier: rangedCritMultiplier,
		unEquippedItems:      items,
		has2H:                items[0].HandType == proto.HandType_HandTypeTwoHand,
	}
}

func (character *Character) RegisterOnItemSwap(callback OnSwapItem) {
	if character == nil || !character.ItemSwap.IsEnabled() {
		return
	}

	character.ItemSwap.onSwapCallbacks = append(character.ItemSwap.onSwapCallbacks, callback)
}

//Helper for handling Effects that use PPMManager to toggle the aura on/off
func (swap *ItemSwap) RegisterOnSwapItemForEffectWithPPMManager(effectID int32, ppm float64, ppmm *PPMManager, aura *Aura) {
	character := swap.character
	character.RegisterOnItemSwap(func(sim *Simulation) {
		mh := character.Equip[proto.ItemSlot_ItemSlotMainHand].Enchant.EffectID == effectID
		oh := character.Equip[proto.ItemSlot_ItemSlotOffHand].Enchant.EffectID == effectID

		procMask := GetMeleeProcMaskForHands(mh, oh)
		*ppmm = character.AutoAttacks.NewPPMManager(ppm, procMask)

		if ppmm.Chance(procMask) == 0 {
			aura.Deactivate(sim)
		} else {
			aura.Activate(sim)
		}
	})

}

//Helper for handling Effects that use the effectID to toggle the aura on and off
func (swap *ItemSwap) ReigsterOnSwapItemForEffect(effectID int32, aura *Aura) {
	character := swap.character
	character.RegisterOnItemSwap(func(sim *Simulation) {
		mh := character.Equip[proto.ItemSlot_ItemSlotMainHand].Enchant.EffectID == effectID
		oh := character.Equip[proto.ItemSlot_ItemSlotOffHand].Enchant.EffectID == effectID

		if !mh && !oh {
			aura.Deactivate(sim)
		} else {
			aura.Activate(sim)
		}
	})
}

func (swap *ItemSwap) IsEnabled() bool {
	return swap.character != nil
}

func (swap *ItemSwap) GetItem(slot proto.ItemSlot) *Item {
	if slot-offset < 0 {
		panic("Not able to swap Item " + slot.String() + " not supported")
	}
	return &swap.unEquippedItems[slot-offset]
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
	newStats := stats.Stats{}
	has2H := swap.GetItem(proto.ItemSlot_ItemSlotMainHand).HandType == proto.HandType_HandTypeTwoHand
	for _, slot := range slots {
		//will swap both on the MainHand Slot for 2H.
		if slot == proto.ItemSlot_ItemSlotOffHand && has2H {
			continue
		}

		if ok, swapStats := swap.swapItem(sim, slot, has2H); ok {
			newStats = newStats.Add(swapStats)
			meeleWeaponSwapped = slot == proto.ItemSlot_ItemSlotMainHand || slot == proto.ItemSlot_ItemSlotOffHand || meeleWeaponSwapped
		}
	}

	character.AddStatsDynamic(sim, newStats)

	if sim.Log != nil {
		sim.Log("Item Swap Stats: %v", newStats)
	}

	for _, onSwap := range swap.onSwapCallbacks {
		onSwap(sim)
	}

	if character.AutoAttacks.IsEnabled() && meeleWeaponSwapped && sim.CurrentTime > 0 {
		character.AutoAttacks.StopMeleeUntil(sim, sim.CurrentTime, false)
	}

	if useGCD {
		character.SetGCDTimer(sim, 1500*time.Millisecond+sim.CurrentTime)
	}
}

func (swap *ItemSwap) swapItem(sim *Simulation, slot proto.ItemSlot, has2H bool) (bool, stats.Stats) {
	character := swap.character
	oldItem := character.Equip[slot]
	newItem := swap.GetItem(slot)

	if newItem.ID == 0 && !(has2H && slot == proto.ItemSlot_ItemSlotOffHand) {
		return false, stats.Stats{}
	}

	character.Equip[slot] = *newItem
	oldItemStats := swap.getItemStats(oldItem)
	newItemStats := swap.getItemStats(*newItem)
	newStats := newItemStats.Add(oldItemStats.Multiply(-1))

	//2H will swap out the offhand also.
	if has2H && slot == proto.ItemSlot_ItemSlotMainHand {
		_, ohStats := swap.swapItem(sim, proto.ItemSlot_ItemSlotOffHand, has2H)
		newStats = newStats.Add(ohStats)
	}

	swap.setItem(slot, oldItem)
	swap.swapWeapon(sim, slot)

	return true, newStats
}

func (swap *ItemSwap) getItemStats(item Item) stats.Stats {
	itemStats := item.Stats
	itemStats = itemStats.Add(item.Enchant.Stats)

	for _, gem := range item.Gems {
		itemStats = itemStats.Add(gem.Stats)
	}

	return itemStats
}

func (swap *ItemSwap) swapWeapon(sim *Simulation, slot proto.ItemSlot) {
	character := swap.character
	if !character.AutoAttacks.IsEnabled() {
		return
	}

	switch slot {
	case proto.ItemSlot_ItemSlotMainHand:
		character.AutoAttacks.MH = character.WeaponFromMainHand(swap.mhCritMultiplier)
		break
	case proto.ItemSlot_ItemSlotOffHand:
		character.AutoAttacks.OH = character.WeaponFromOffHand(swap.ohCritMultiplier)
		//Special case for when the OHAuto Spell was setup with a non weapon and does not have a crit multiplier.
		character.AutoAttacks.OHAuto.CritMultiplier = swap.ohCritMultiplier
		character.PseudoStats.CanBlock = character.Equip[proto.ItemSlot_ItemSlotOffHand].WeaponType == proto.WeaponType_WeaponTypeShield
		break
	case proto.ItemSlot_ItemSlotRanged:
		character.AutoAttacks.Ranged = character.WeaponFromRanged(swap.rangedCritMultiplier)
		break
	}

	character.AutoAttacks.IsDualWielding = character.Equip[proto.ItemSlot_ItemSlotMainHand].SwingSpeed != 0 && character.Equip[proto.ItemSlot_ItemSlotOffHand].SwingSpeed != 0
}

func (swap *ItemSwap) finalize() {
	if !swap.IsEnabled() {
		return
	}

	swap.initialEquippedItems = getInitialEquippedItems(swap.character)
	swap.initialUnequippedItems = swap.unEquippedItems
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

	for _, onSwap := range swap.onSwapCallbacks {
		onSwap(sim)
	}
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

func toItem(itemSpec *proto.ItemSpec) Item {
	if itemSpec == nil {
		return Item{}
	}

	return NewItem(ItemSpec{
		ID:      itemSpec.Id,
		Gems:    itemSpec.Gems,
		Enchant: itemSpec.Enchant,
	})
}
