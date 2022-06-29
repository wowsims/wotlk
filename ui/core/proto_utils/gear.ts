import { Gem } from '/tbc/core/proto/common.js';
import { GemColor } from '/tbc/core/proto/common.js';
import { ItemSlot } from '/tbc/core/proto/common.js';
import { ItemSpec } from '/tbc/core/proto/common.js';
import { EquipmentSpec } from '/tbc/core/proto/common.js';
import { WeaponType } from '/tbc/core/proto/common.js';
import { equalsOrBothNull } from '/tbc/core/utils.js';
import { getEnumValues } from '/tbc/core/utils.js';
import { isBluntWeaponType, isSharpWeaponType } from '/tbc/core/proto_utils/utils.js';

import { isMetaGemActive } from './gems.js';
import { gemMatchesSocket } from './gems.js';
import { EquippedItem } from './equipped_item.js';
import { validWeaponCombo } from './utils.js';

type InternalGear = Record<ItemSlot, EquippedItem | null>;

/**
 * Represents a full gear set, including items/enchants/gems for every slot.
 *
 * This is an immutable type.
 */
export class Gear {
	private readonly gear: InternalGear;

	constructor(gear: Partial<InternalGear>) {
		getEnumValues(ItemSlot).forEach(slot => {
			if (!gear[slot as ItemSlot])
				gear[slot as ItemSlot] = null;
		});
		this.gear = gear as InternalGear;
	}

	equals(other: Gear): boolean {
		return this.asArray().every((thisItem, slot) => equalsOrBothNull(thisItem, other.getEquippedItem(slot), (a, b) => a.equals(b)));
	}

	/**
	 * Returns a new Gear set with the item equipped.
	 *
	 * Checks for validity and removes/exchanges items/gems as needed.
	 */
	withEquippedItem(newSlot: ItemSlot, newItem: EquippedItem | null): Gear {
		// Create a new identical set of gear
		const newInternalGear = this.asMap();

		if (newItem) {
			// If the new item has unique gems, remove matching.
			newItem.gems
				.filter(gem => gem?.unique)
				.forEach(gem => {
					getEnumValues(ItemSlot).map(slot => Number(slot) as ItemSlot).forEach(slot => {
						newInternalGear[slot] = newInternalGear[slot]?.removeGemsWithId(gem!.id) || null;
					});
				});

			// If the new item is unique, remove matching items.
			if (newItem.item.unique) {
				getEnumValues(ItemSlot).map(slot => Number(slot) as ItemSlot).forEach(slot => {
					if (newInternalGear[slot]?.item.id == newItem.item.id) {
						newInternalGear[slot] = null;
					}
				});
			}
		}

		// Actually assign the new item.
		newInternalGear[newSlot] = newItem;

		// Check for valid weapon combos.
		if (!validWeaponCombo(newInternalGear[ItemSlot.ItemSlotMainHand]?.item, newInternalGear[ItemSlot.ItemSlotOffHand]?.item)) {
			if (newSlot == ItemSlot.ItemSlotMainHand) {
				newInternalGear[ItemSlot.ItemSlotOffHand] = null;
			} else {
				newInternalGear[ItemSlot.ItemSlotMainHand] = null;
			}
		}

		return new Gear(newInternalGear);
	}

	getEquippedItem(slot: ItemSlot): EquippedItem | null {
		return this.gear[slot];
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

	asMap(): InternalGear {
		const newInternalGear: Partial<InternalGear> = {};
		getEnumValues(ItemSlot).map(slot => Number(slot) as ItemSlot).forEach(slot => {
			newInternalGear[slot] = this.getEquippedItem(slot);
		});
		return newInternalGear as InternalGear;
	}

	asArray(): Array<EquippedItem | null> {
		return Object.values(this.gear);
	}

	asSpec(): EquipmentSpec {
		return EquipmentSpec.create({
			items: this.asArray().map(ei => ei ? ei.asSpec() : ItemSpec.create()),
		});
	}

	getAllGems(): Array<Gem> {
		return this.asArray()
			.map(equippedItem => equippedItem == null ? [] : equippedItem.gems.filter(gem => gem != null) as Array<Gem>)
			.flat();
	}

	getGemsOfColor(color: GemColor): Array<Gem> {
		return this.getAllGems().filter(gem => gem.color == color);
	}

	getMetaGem(): Gem | null {
		return this.getGemsOfColor(GemColor.GemColorMeta)[0] || null;
	}

	// Returns true if this gear set has a meta gem AND the other gems meet the meta's conditions.
	hasActiveMetaGem(): boolean {
		const metaGem = this.getMetaGem();
		if (!metaGem) {
			return false;
		}

		const gems = this.getAllGems();
		return isMetaGemActive(
			metaGem,
			gems.filter(gem => gemMatchesSocket(gem, GemColor.GemColorRed)).length,
			gems.filter(gem => gemMatchesSocket(gem, GemColor.GemColorYellow)).length,
			gems.filter(gem => gemMatchesSocket(gem, GemColor.GemColorBlue)).length);
	}

	hasInactiveMetaGem(): boolean {
		return this.getMetaGem() != null && !this.hasActiveMetaGem();
	}

	withoutMetaGem(): Gear {
		const headItem = this.getEquippedItem(ItemSlot.ItemSlotHead);
		const metaGem = this.getMetaGem();
		if (headItem && metaGem) {
			return this.withEquippedItem(ItemSlot.ItemSlotHead, headItem.removeGemsWithId(metaGem.id));
		} else {
			return this;
		}
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
}
