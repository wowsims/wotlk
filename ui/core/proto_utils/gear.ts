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

	abstract getItemSlots(): ItemSlot[]

	equals(other: BaseGear): boolean {
		return this.asArray().every((thisItem, slot) => equalsOrBothNull(thisItem, other.getEquippedItem(slot), (a, b) => a.equals(b)));
	}

	getEquippedItem(slot: ItemSlot): EquippedItem | null {
		return this.gear[slot] || null;
	}

	asArray(): Array<EquippedItem | null> {
		return Object.values(this.gear);
	}

	asMap(): Partial<InternalGear> {
		const newInternalGear: Partial<InternalGear> = {};
		this.getItemSlots().map(slot => Number(slot) as ItemSlot).forEach(slot => {
			newInternalGear[slot] = this.getEquippedItem(slot);
		});
		return newInternalGear;
	}

	/**
	 * Returns a new Gear set with the item equipped.
	 *
	 * Checks for validity and removes/exchanges items/gems as needed.
	 */
	protected withEquippedItemInternal(newSlot: ItemSlot, newItem: EquippedItem | null, canDualWield2H: boolean): Partial<InternalGear> {
		// Create a new identical set of gear
		const newInternalGear = this.asMap();

		if (newItem) {
			this.removeUniqueItems(newInternalGear, newItem);
		}

		// Actually assign the new item.
		newInternalGear[newSlot] = newItem;

		BaseGear.validateWeaponCombo(newInternalGear, newSlot, canDualWield2H);

		return newInternalGear;
	}

	private removeUniqueItems(gear: Partial<InternalGear>, newItem: EquippedItem) {
		if (newItem.item.unique) {
			this.getItemSlots().map(slot => Number(slot) as ItemSlot).forEach(slot => {
				if (gear[slot]?.item.id == newItem.item.id) {
					gear[slot] = null;
				}
			});
		}
	}

	private static validateWeaponCombo(gear: Partial<InternalGear>, newSlot: ItemSlot, canDualWield2H: boolean) {
		// Check for valid weapon combos.
		if (!validWeaponCombo(gear[ItemSlot.ItemSlotMainHand]?.item, gear[ItemSlot.ItemSlotOffHand]?.item, canDualWield2H)) {
			if (newSlot == ItemSlot.ItemSlotOffHand) {
				gear[ItemSlot.ItemSlotMainHand] = null;
			} else {
				gear[ItemSlot.ItemSlotOffHand] = null;
			}
		}
	}

	toDatabase(): SimDatabase {
		const equippedItems = this.asArray().filter(ei => ei != null) as Array<EquippedItem>;
		return SimDatabase.create({
			items: distinct(equippedItems.map(ei => BaseGear.itemToDB(ei.item))),
			enchants: distinct(equippedItems.filter(ei => ei.enchant).map(ei => BaseGear.enchantToDB(ei.enchant!))),
		});
	}

	private static itemToDB(item: Item): SimItem {
		return SimItem.fromJson(Item.toJson(item), { ignoreUnknownFields: true });
	}

	private static enchantToDB(enchant: Enchant): SimEnchant {
		return SimEnchant.fromJson(Enchant.toJson(enchant), { ignoreUnknownFields: true });
	}

	// TODO: Add rune
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

	withEquippedItem(newSlot: ItemSlot, newItem: EquippedItem | null, canDualWield2H: boolean): Gear {
		return new Gear(this.withEquippedItemInternal(newSlot, newItem, canDualWield2H));
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
}

/**
 * Represents a item swap gear set, including items/enchants.
 *
 * This is an immutable type.
 */
export class ItemSwapGear extends BaseGear {

	constructor(gear: Partial<InternalGear>) {
		super(gear);
	}

	getItemSlots(): ItemSlot[] {
		return [ItemSlot.ItemSlotMainHand, ItemSlot.ItemSlotOffHand, ItemSlot.ItemSlotRanged];
	}

	withEquippedItem(newSlot: ItemSlot, newItem: EquippedItem | null, canDualWield2H: boolean): ItemSwapGear {
		return new ItemSwapGear(this.withEquippedItemInternal(newSlot, newItem, canDualWield2H));
	}

	toProto(): ItemSwap {
		return ItemSwap.create({
			mhItem: this.gear[ItemSlot.ItemSlotMainHand]?.asSpec(),
			ohItem: this.gear[ItemSlot.ItemSlotOffHand]?.asSpec(),
			rangedItem: this.gear[ItemSlot.ItemSlotRanged]?.asSpec(),
		})
	}
}
