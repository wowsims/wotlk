import { ItemSlot } from '../proto/common.js';
import { EquippedItem } from './equipped_item.js';
import { validWeaponCombo } from './utils.js';

type InternalGear = Record<ItemSlot, EquippedItem | null>;

/**
 * Represents a full gear set, including items/enchants/gems for every slot.
 *
 * This is an immutable type.
 */
export class ItemSwapGear {
	private readonly gear: InternalGear;

	constructor() {
		const newInternalGear: Partial<InternalGear> = {};
		let slotList = [ItemSlot.ItemSlotMainHand, ItemSlot.ItemSlotOffHand, ItemSlot.ItemSlotRanged];
		slotList.forEach(slot => {
			newInternalGear[slot as ItemSlot] = null;
		});
		this.gear = newInternalGear as InternalGear;
	}

	getEquippedItem(slot: ItemSlot): EquippedItem | null {
		return this.gear[slot];
	}

	equipItem(slot: ItemSlot, equppedItem: EquippedItem | null ) {
		this.gear[slot] = equppedItem;
		
		// Check for valid weapon combos.
		if (!validWeaponCombo(this.gear[ItemSlot.ItemSlotMainHand]?.item, this.gear[ItemSlot.ItemSlotOffHand]?.item, false)) {
			if (slot == ItemSlot.ItemSlotOffHand) {
				this.gear[ItemSlot.ItemSlotMainHand] = null;
			} else {
				this.gear[ItemSlot.ItemSlotOffHand] = null;
			}
		}
	}

}
