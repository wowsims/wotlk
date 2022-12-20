import { ItemSlot, SimItem, SimDatabase, SimEnchant, SimGem, ItemSwap, ItemSpec } from '../proto/common.js';
import { EquippedItem } from './equipped_item.js';
import { validWeaponCombo } from './utils.js';
import { distinct } from '../utils.js'
import {
	UIEnchant as Enchant,
	UIGem as Gem,
	UIItem as Item,
} from '../proto/ui.js';

type InternalGear = Record<ItemSlot, EquippedItem | null>;

/**
 * Represents a item swap gear set, including items/enchants/gems.
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

	equipItem(slot: ItemSlot, equppedItem: EquippedItem | null) {
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

	toProto(): ItemSwap {
		return ItemSwap.create({
			mhItem: this.gear[ItemSlot.ItemSlotMainHand]?.asSpec(),
			ohItem: this.gear[ItemSlot.ItemSlotOffHand]?.asSpec(),
			rangedItem: this.gear[ItemSlot.ItemSlotRanged]?.asSpec(),
		})
	}

	asArray(): Array<EquippedItem | null> {
		return Object.values(this.gear);
	}

	toDatabase(): SimDatabase {
		const equippedItems = this.asArray().filter(ei => ei != null) as Array<EquippedItem>;
		return SimDatabase.create({
			items: distinct(equippedItems.map(ei => ItemSwapGear.itemToDB(ei.item))),
			enchants: distinct(equippedItems.filter(ei => ei.enchant).map(ei => ItemSwapGear.enchantToDB(ei.enchant!))),
			gems: distinct(equippedItems.map(ei => ei.curGems(true).map(gem => ItemSwapGear.gemToDB(gem))).flat()),
		});
	}

	private static itemToDB(item: Item): SimItem {
		return SimItem.fromJson(Item.toJson(item), { ignoreUnknownFields: true });
	}

	private static enchantToDB(enchant: Enchant): SimEnchant {
		return SimEnchant.fromJson(Enchant.toJson(enchant), { ignoreUnknownFields: true });
	}

	private static gemToDB(gem: Gem): SimGem {
		return SimGem.fromJson(Gem.toJson(gem), { ignoreUnknownFields: true });
	}

}
