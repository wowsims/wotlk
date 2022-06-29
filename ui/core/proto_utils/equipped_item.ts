import { Enchant } from '/tbc/core/proto/common.js';
import { Gem } from '/tbc/core/proto/common.js';
import { Item } from '/tbc/core/proto/common.js';
import { ItemSlot } from '/tbc/core/proto/common.js';
import { ItemSpec } from '/tbc/core/proto/common.js';
import { Stat } from '/tbc/core/proto/common.js';

import { ActionId } from './action_id.js';
import { enchantAppliesToItem } from './utils.js';
import { gemEligibleForSocket } from './gems.js';
import { gemMatchesSocket } from './gems.js';
import { Stats } from './stats.js';

export function getWowheadItemId(item: Item): number {
	return item.wowheadId || item.id;
}

export function getWeaponDPS(item: Item): number {
	return ((item.weaponDamageMin + item.weaponDamageMax) / 2) / (item.weaponSpeed || 1);
}

/**
 * Represents an equipped item along with enchants/gems attached to it.
 *
 * This is an immutable type.
 */
export class EquippedItem {
	readonly _item: Item;
	readonly _enchant: Enchant | null;
	readonly _gems: Array<Gem | null>;

	constructor(item: Item, enchant?: Enchant | null, gems?: Array<Gem | null>) {
		this._item = item;
		this._enchant = enchant || null;
		this._gems = gems || [];

		// Fill gems with null so we always have the same number of gems as gem slots.
		if (this._gems.length < item.gemSockets.length) {
			this._gems = this._gems.concat(new Array(item.gemSockets.length - this._gems.length).fill(null));
		}
	}

	get item(): Item {
		// Make a defensive copy
		return Item.clone(this._item);
	}

	get enchant(): Enchant | null {
		// Make a defensive copy
		return this._enchant ? Enchant.clone(this._enchant) : null;
	}

	get gems(): Array<Gem | null> {
		// Make a defensive copy
		return this._gems.map(gem => gem == null ? null : Gem.clone(gem));
	}

	equals(other: EquippedItem) {
		if (!Item.equals(this._item, other.item))
			return false;

		if ((this._enchant == null) != (other.enchant == null))
			return false;

		if (this._enchant && other.enchant && !Enchant.equals(this._enchant, other.enchant))
			return false;

		if (this._gems.length != other.gems.length)
			return false;

		for (let i = 0; i < this._gems.length; i++) {
			if ((this._gems[i] == null) != (other.gems[i] == null))
				return false;

			if (this._gems[i] && other.gems[i] && !Gem.equals(this._gems[i]!, other.gems[i]!))
				return false;
		}

		return true;
	}

	/**
	 * Replaces the item and tries to keep the existing enchants/gems if possible.
	 */
	withItem(item: Item): EquippedItem {
		let newEnchant = null;
		if (this._enchant && enchantAppliesToItem(this._enchant, item))
			newEnchant = this._enchant;

		// Reorganize gems to match as many colors in the new item as possible.
		const newGems = new Array(item.gemSockets.length).fill(null);
		this._gems.filter(gem => gem != null).forEach(gem => {
			const firstMatchingIndex = item.gemSockets.findIndex((socketColor, socketIdx) => !newGems[socketIdx] && gemMatchesSocket(gem!, socketColor));
			const firstEligibleIndex = item.gemSockets.findIndex((socketColor, socketIdx) => !newGems[socketIdx] && gemEligibleForSocket(gem!, socketColor));
			if (firstMatchingIndex != -1) {
				newGems[firstMatchingIndex] = gem;
			} else if (firstEligibleIndex != -1) {
				newGems[firstEligibleIndex] = gem;
			}
		});

		return new EquippedItem(item, newEnchant, newGems);
	}

	/**
	 * Returns a new EquippedItem with the given enchant applied.
	 */
	withEnchant(enchant: Enchant | null): EquippedItem {
		return new EquippedItem(this._item, enchant, this._gems);
	}

	/**
	 * Returns a new EquippedItem with the given gem socketed.
	 */
	private withGemHelper(gem: Gem | null, socketIdx: number): EquippedItem {
		if (this._gems.length <= socketIdx) {
			throw new Error('No gem socket with index ' + socketIdx);
		}

		const newGems = this._gems.slice();
		newGems[socketIdx] = gem;

		return new EquippedItem(this._item, this._enchant, newGems);
	}

	/**
	 * Returns a new EquippedItem with the given gem socketed.
	 *
	 * Also ensures validity of the item on its own. Currently this just means enforcing unique gems.
	 */
	withGem(gem: Gem | null, socketIdx: number): EquippedItem {
		let curItem: EquippedItem | null = this;

		if (gem && gem.unique) {
			curItem = curItem.removeGemsWithId(gem.id);
		}

		return curItem.withGemHelper(gem, socketIdx);
	}

	removeGemsWithId(gemId: number): EquippedItem {
		let curItem: EquippedItem | null = this;
		// Remove any currently socketed identical gems.
		for (let i = 0; i < curItem._gems.length; i++) {
			if (curItem._gems[i]?.id == gemId) {
				curItem = curItem.withGemHelper(null, i);
			}
		}
		return curItem;
	}

	asActionId(): ActionId {
		return ActionId.fromItemId(this._item.id);
	}

	asSpec(): ItemSpec {
		return ItemSpec.create({
			id: this._item.id,
			enchant: this._enchant?.id,
			gems: this._gems.map(gem => gem?.id || 0),
		});
	}
};
