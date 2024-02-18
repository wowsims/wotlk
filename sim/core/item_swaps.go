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

	// Which slots to actually swap.
	slots []proto.ItemSlot

	// Holds items that are currently not equipped
	unEquippedItems [3]Item
	swapped         bool
}

/*
TODO All the extra parameters here and the code in multiple places for handling the Weapon struct is really messy,

	we'll need to figure out something cleaner as this will be quite error-prone
*/
func (character *Character) enableItemSwap(itemSwap *proto.ItemSwap, mhCritMultiplier float64, ohCritMultiplier float64, rangedCritMultiplier float64) {
	var slots []proto.ItemSlot
	hasMhSwap := itemSwap.MhItem != nil && itemSwap.MhItem.Id != 0
	hasOhSwap := itemSwap.OhItem != nil && itemSwap.OhItem.Id != 0
	hasRangedSwap := itemSwap.RangedItem != nil && itemSwap.RangedItem.Id != 0

	mainItems := [3]Item{
		character.Equipment[proto.ItemSlot_ItemSlotMainHand],
		character.Equipment[proto.ItemSlot_ItemSlotOffHand],
		character.Equipment[proto.ItemSlot_ItemSlotRanged],
	}
	swapItems := [3]Item{
		toItem(itemSwap.MhItem),
		toItem(itemSwap.OhItem),
		toItem(itemSwap.RangedItem),
	}

	// Handle MH and OH together, because present MH + empty OH --> swap MH and unequip OH
	if hasMhSwap || hasOhSwap {
		if swapItems[0].ID != mainItems[0].ID {
			slots = append(slots, proto.ItemSlot_ItemSlotMainHand)
		}
		if swapItems[1].ID != mainItems[1].ID {
			slots = append(slots, proto.ItemSlot_ItemSlotOffHand)
		}
	}
	if hasRangedSwap {
		if swapItems[2].ID != mainItems[2].ID {
			slots = append(slots, proto.ItemSlot_ItemSlotRanged)
		}
	}

	if len(slots) == 0 {
		return
	}

	character.ItemSwap = ItemSwap{
		mhCritMultiplier:     mhCritMultiplier,
		ohCritMultiplier:     ohCritMultiplier,
		rangedCritMultiplier: rangedCritMultiplier,
		slots:                slots,
		unEquippedItems:      swapItems,
		swapped:              false,
	}
}

func (swap *ItemSwap) initialize(character *Character) {
	swap.character = character
}

func (character *Character) RegisterOnItemSwap(callback OnSwapItem) {
	if character == nil || !character.ItemSwap.IsEnabled() {
		return
	}

	character.ItemSwap.onSwapCallbacks = append(character.ItemSwap.onSwapCallbacks, callback)
}

// Helper for handling Effects that use PPMManager to toggle the aura on/off
func (swap *ItemSwap) RegisterOnSwapItemForEffectWithPPMManager(effectID int32, ppm float64, ppmm *PPMManager, aura *Aura) {
	character := swap.character
	character.RegisterOnItemSwap(func(sim *Simulation) {
		procMask := character.GetProcMaskForEnchant(effectID)
		*ppmm = character.AutoAttacks.NewPPMManager(ppm, procMask)

		if ppmm.Chance(procMask) == 0 {
			aura.Deactivate(sim)
		} else {
			aura.Activate(sim)
		}
	})

}

// Helper for handling Effects that use the itemID to toggle the aura on and off
func (swap *ItemSwap) RegisterOnSwapItemForItemEffect(itemID int32, aura *Aura) {
	character := swap.character
	character.RegisterOnItemSwap(func(sim *Simulation) {
		procMask := character.GetProcMaskForItem(itemID)

		if procMask == ProcMaskUnknown {
			aura.Deactivate(sim)
		} else {
			aura.Activate(sim)
		}
	})
}

// Helper for handling Effects that use the effectID to toggle the aura on and off
func (swap *ItemSwap) RegisterOnSwapItemForEnchantEffect(effectID int32, aura *Aura) {
	character := swap.character
	character.RegisterOnItemSwap(func(sim *Simulation) {
		procMask := character.GetProcMaskForEnchant(effectID)

		if procMask == ProcMaskUnknown {
			aura.Deactivate(sim)
		} else {
			aura.Activate(sim)
		}
	})
}

func (swap *ItemSwap) IsEnabled() bool {
	return swap.character != nil && len(swap.slots) > 0
}

func (swap *ItemSwap) IsSwapped() bool {
	return swap.swapped
}

func (swap *ItemSwap) GetItem(slot proto.ItemSlot) *Item {
	if slot-offset < 0 {
		panic("Not able to swap Item " + slot.String() + " not supported")
	}
	return &swap.unEquippedItems[slot-offset]
}

func (swap *ItemSwap) CalcStatChanges(slots []proto.ItemSlot) stats.Stats {
	newStats := stats.Stats{}
	for _, slot := range slots {
		oldItemStats := swap.getItemStats(swap.character.Equipment[slot])
		newItemStats := swap.getItemStats(*swap.GetItem(slot))
		newStats = newStats.Add(newItemStats.Subtract(oldItemStats))
	}

	return newStats
}

func (swap *ItemSwap) SwapItems(sim *Simulation, slots []proto.ItemSlot) {
	if !swap.IsEnabled() {
		return
	}

	character := swap.character

	meleeWeaponSwapped := false
	newStats := stats.Stats{}
	has2H := swap.GetItem(proto.ItemSlot_ItemSlotMainHand).HandType == proto.HandType_HandTypeTwoHand
	for _, slot := range slots {

		//will swap both on the MainHand Slot for 2H.
		if slot == proto.ItemSlot_ItemSlotOffHand && has2H {
			continue
		}

		if ok, swapStats := swap.swapItem(slot, has2H); ok {
			newStats = newStats.Add(swapStats)
			meleeWeaponSwapped = slot == proto.ItemSlot_ItemSlotMainHand || slot == proto.ItemSlot_ItemSlotOffHand || meleeWeaponSwapped
		}
	}

	character.AddStatsDynamic(sim, newStats)

	if sim.Log != nil {
		sim.Log("Item Swap Stats: %v", newStats)
	}

	for _, onSwap := range swap.onSwapCallbacks {
		onSwap(sim)
	}

	if character.AutoAttacks.AutoSwingMelee && meleeWeaponSwapped && sim.CurrentTime > 0 {
		character.AutoAttacks.StopMeleeUntil(sim, sim.CurrentTime, false)
	}

	// If GCD is ready then use the GCD, otherwise we assume it's being used along side a spell.
	if character.GCD.IsReady(sim) {
		newGCD := sim.CurrentTime + 1500*time.Millisecond
		character.SetGCDTimer(sim, newGCD)
	}

	swap.swapped = !swap.swapped
}

func (swap *ItemSwap) swapItem(slot proto.ItemSlot, has2H bool) (bool, stats.Stats) {
	oldItem := swap.character.Equipment[slot]
	newItem := swap.GetItem(slot)

	if newItem.ID == 0 && !(has2H && slot == proto.ItemSlot_ItemSlotOffHand) {
		return false, stats.Stats{}
	}

	swap.character.Equipment[slot] = *newItem
	oldItemStats := swap.getItemStats(oldItem)
	newItemStats := swap.getItemStats(*newItem)
	newStats := newItemStats.Subtract(oldItemStats)

	//2H will swap out the offhand also.
	if has2H && slot == proto.ItemSlot_ItemSlotMainHand {
		_, ohStats := swap.swapItem(proto.ItemSlot_ItemSlotOffHand, has2H)
		newStats = newStats.Add(ohStats)
	}

	swap.unEquippedItems[slot-offset] = oldItem
	swap.swapWeapon(slot)

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

func (swap *ItemSwap) swapWeapon(slot proto.ItemSlot) {
	character := swap.character

	switch slot {
	case proto.ItemSlot_ItemSlotMainHand:
		if character.AutoAttacks.AutoSwingMelee {
			character.AutoAttacks.SetMH(character.WeaponFromMainHand(swap.mhCritMultiplier))
		}
	case proto.ItemSlot_ItemSlotOffHand:
		if character.AutoAttacks.AutoSwingMelee {
			weapon := character.WeaponFromOffHand(swap.ohCritMultiplier)
			character.AutoAttacks.SetOH(weapon)

			character.AutoAttacks.IsDualWielding = weapon.SwingSpeed != 0
			character.PseudoStats.CanBlock = character.OffHand().WeaponType == proto.WeaponType_WeaponTypeShield
		}
	case proto.ItemSlot_ItemSlotRanged:
		if character.AutoAttacks.AutoSwingRanged {
			character.AutoAttacks.SetRanged(character.WeaponFromRanged(swap.rangedCritMultiplier))
		}
	}
}

func (swap *ItemSwap) reset(sim *Simulation) {
	if !swap.IsEnabled() || !swap.IsSwapped() {
		return
	}

	swap.SwapItems(sim, swap.slots)
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
