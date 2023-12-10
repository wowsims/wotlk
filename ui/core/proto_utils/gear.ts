import { EquipmentSpec, ItemSlot, ItemSpec, ItemSwap, Profession, SimDatabase, SimEnchant, SimItem } from '../proto/common.js';
import {
	UIEnchant as Enchant,
	UIItem as Item,
} from '../proto/ui.js';
import { isBluntWeaponType, isSharpWeaponType } from '../proto_utils/utils.js';
import { distinct, equalsOrBothNull, getEnumValues } from '../utils.js';

import { EquippedItem } from './equipped_item.js';
import { validWeaponCombo } from './utils.js';

type InternalGear = Record<ItemSlot, EquippedItem | null>;

abstract class BaseGear {
	protected readonly gear: InternalGear;

	constructor(gear: Partial<InternalGear>) {
		this.getItemSlots().forEach(slot => {
			if (!gear[slot as ItemSlot])
				gear[slot as ItemSlot] = null;
		});
		this.gear = gear as InternalGear;
	}

	getEquippedItem(slot: ItemSlot): EquippedItem | null {
		return this.gear[slot];
	}

	asArray(): Array<EquippedItem | null> {
		return Object.values(this.gear);
	}

	removeUniqueItems(gear: InternalGear, newItem: EquippedItem) {
		if (newItem.item.unique) {
			this.getItemSlots().map(slot => Number(slot) as ItemSlot).forEach(slot => {
				if (gear[slot]?.item.id == newItem.item.id) {
					gear[slot] = null;
				}
			});
		}
	}

	validateWeaponCombo(gear: InternalGear, newSlot: ItemSlot, canDualWield2H: boolean) {
		// Check for valid weapon combos.
		if (!validWeaponCombo(gear[ItemSlot.ItemSlotMainHand]?.item, gear[ItemSlot.ItemSlotOffHand]?.item, canDualWield2H)) {
			if (newSlot == ItemSlot.ItemSlotOffHand) {
				gear[ItemSlot.ItemSlotMainHand] = null;
			} else {
				gear[ItemSlot.ItemSlotOffHand] = null;
			}
		}
	}

	abstract toDatabase(): SimDatabase
	abstract getItemSlots(): ItemSlot[]

	protected static itemToDB(item: Item): SimItem {
		return SimItem.fromJson(Item.toJson(item), { ignoreUnknownFields: true });
	}

	protected static enchantToDB(enchant: Enchant): SimEnchant {
		return SimEnchant.fromJson(Enchant.toJson(enchant), { ignoreUnknownFields: true });
	}
}

/**
 * Represents a full gear set, including items/enchants for every slot.
 *
 * This is an immutable type.
 */
export class Gear extends BaseGear {

	constructor(gear: Partial<InternalGear>) {
		super(gear);
	}

	getItemSlots(): ItemSlot[] {
		return getEnumValues(ItemSlot);
	}

	equals(other: Gear): boolean {
		return this.asArray().every((thisItem, slot) => equalsOrBothNull(thisItem, other.getEquippedItem(slot), (a, b) => a.equals(b)));
	}

	/**
	 * Returns a new Gear set with the item equipped.
	 *
	 * Checks for validity and removes/exchanges items as needed.
	 */
	withEquippedItem(newSlot: ItemSlot, newItem: EquippedItem | null, canDualWield2H: boolean): Gear {
		// Create a new identical set of gear
		const newInternalGear = this.asMap();

		if (newItem) {
			this.removeUniqueItems(newInternalGear, newItem);
		}

		// Actually assign the new item.
		newInternalGear[newSlot] = newItem;

		this.validateWeaponCombo(newInternalGear, newSlot, canDualWield2H);

		return new Gear(newInternalGear);
	}

	getTrinkets(): Array<EquippedItem | null> {
		return [
			this.getEquippedItem(ItemSlot.ItemSlotTrinket1),
			this.getEquippedItem(ItemSlot.ItemSlotTrinket2),
		];
	}

	hasTrinket(itemId: number): boolean {
		return this.getTrinkets().map(t => t?.item.id).includes(itemId);
	}

	hasRelic(itemId: number): boolean {
		const relicItem = this.getEquippedItem(ItemSlot.ItemSlotRanged);

		if (!relicItem) {
			return false;
		}

		return relicItem!.item.id == itemId;
	}

	asMap(): InternalGear {
		const newInternalGear: Partial<InternalGear> = {};
		getEnumValues(ItemSlot).map(slot => Number(slot) as ItemSlot).forEach(slot => {
			newInternalGear[slot] = this.getEquippedItem(slot);
		});
		return newInternalGear as InternalGear;
	}

	asSpec(): EquipmentSpec {
		return EquipmentSpec.create({
			items: this.asArray().map(ei => ei ? ei.asSpec() : ItemSpec.create()),
		});
	}

	hasBluntMHWeapon(): boolean {
		const weapon = this.getEquippedItem(ItemSlot.ItemSlotMainHand);
		return weapon != null && isBluntWeaponType(weapon.item.weaponType);
	}
	hasSharpMHWeapon(): boolean {
		const weapon = this.getEquippedItem(ItemSlot.ItemSlotMainHand);
		return weapon != null && isSharpWeaponType(weapon.item.weaponType);
	}
	hasBluntOHWeapon(): boolean {
		const weapon = this.getEquippedItem(ItemSlot.ItemSlotOffHand);
		return weapon != null && isBluntWeaponType(weapon.item.weaponType);
	}
	hasSharpOHWeapon(): boolean {
		const weapon = this.getEquippedItem(ItemSlot.ItemSlotOffHand);
		return weapon != null && isSharpWeaponType(weapon.item.weaponType);
	}

	getProfessionRequirements(): Array<Profession> {
		return distinct((this.asArray().filter(ei => ei != null) as Array<EquippedItem>)
			.map(ei => ei.getProfessionRequirements())
			.flat());
	}
	getFailedProfessionRequirements(professions: Array<Profession>): Array<Item | Enchant> {
		return (this.asArray().filter(ei => ei != null) as Array<EquippedItem>)
			.map(ei => ei.getFailedProfessionRequirements(professions))
			.flat();
	}

	toDatabase(): SimDatabase {
		const equippedItems = this.asArray().filter(ei => ei != null) as Array<EquippedItem>;
		return SimDatabase.create({
			items: distinct(equippedItems.map(ei => Gear.itemToDB(ei.item))),
			enchants: distinct(equippedItems.filter(ei => ei.enchant).map(ei => Gear.enchantToDB(ei.enchant!))),
		});
	}
}

/**
 * Represents a item swap gear set, including items/enchants.
 *
 * This is an immutable type.
 */
export class ItemSwapGear extends BaseGear {

	constructor() {
		super({});
	}

	getItemSlots(): ItemSlot[] {
		return [ItemSlot.ItemSlotMainHand, ItemSlot.ItemSlotOffHand, ItemSlot.ItemSlotRanged];
	}

	equipItem(slot: ItemSlot, equippedItem: EquippedItem | null, canDualWield2H: boolean) {
		if (equippedItem) {
			this.removeUniqueItems(this.gear, equippedItem);
		}

		this.gear[slot] = equippedItem;
		this.validateWeaponCombo(this.gear, slot, canDualWield2H);
	}

	toProto(): ItemSwap {
		return ItemSwap.create({
			mhItem: this.gear[ItemSlot.ItemSlotMainHand]?.asSpec(),
			ohItem: this.gear[ItemSlot.ItemSlotOffHand]?.asSpec(),
			rangedItem: this.gear[ItemSlot.ItemSlotRanged]?.asSpec(),
		})
	}

	toDatabase(): SimDatabase {
		const equippedItems = this.asArray().filter(ei => ei != null) as Array<EquippedItem>;
		return SimDatabase.create({
			items: distinct(equippedItems.map(ei => ItemSwapGear.itemToDB(ei.item))),
			enchants: distinct(equippedItems.filter(ei => ei.enchant).map(ei => ItemSwapGear.enchantToDB(ei.enchant!))),
		});
	}
}
