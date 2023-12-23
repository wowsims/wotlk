import { GemColor } from '../proto/common.js';
import { ItemSpec } from '../proto/common.js';
import { ItemType } from '../proto/common.js';
import { Profession } from '../proto/common.js';
import {
	UIEnchant as Enchant,
	UIGem as Gem,
	UIItem as Item,
} from '../proto/ui.js';
import { distinct } from '../utils.js';

import { ActionId } from './action_id.js';
import { enchantAppliesToItem } from './utils.js';
import { gemEligibleForSocket, gemMatchesSocket } from './gems.js';
import { Stats } from './stats.js';

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

	readonly numPossibleSockets: number;

	constructor(item: Item, enchant?: Enchant | null, gems?: Array<Gem | null>) {
		this._item = item;
		this._enchant = enchant || null;
		this._gems = gems || [];

		this.numPossibleSockets = this.numSockets(true);

		// Fill gems with null so we always have the same number of gems as gem slots.
		if (this._gems.length < this.numPossibleSockets) {
			this._gems = this._gems.concat(new Array(this.numPossibleSockets - this._gems.length).fill(null));
		}
	}

	get item(): Item {
		// Make a defensive copy
		return Item.clone(this._item);
	}

	get id(): number {
		return this._item.id;
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
		this._gems.slice(0, this._item.gemSockets.length).filter(gem => gem != null).forEach(gem => {
			const firstMatchingIndex = item.gemSockets.findIndex((socketColor, socketIdx) => !newGems[socketIdx] && gemMatchesSocket(gem!, socketColor));
			const firstEligibleIndex = item.gemSockets.findIndex((socketColor, socketIdx) => !newGems[socketIdx] && gemEligibleForSocket(gem!, socketColor));
			if (firstMatchingIndex != -1) {
				newGems[firstMatchingIndex] = gem;
			} else if (firstEligibleIndex != -1) {
				newGems[firstEligibleIndex] = gem;
			}
		});

		// Copy the extra socket gem directly.
		if (this.couldHaveExtraSocket()) {
			newGems.push(this._gems[this._gems.length - 1]);
		}

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

	removeAllGems(): EquippedItem {
		let curItem: EquippedItem | null = this;

		for (let i = 0; i < curItem._gems.length; i++) {
			curItem = curItem.withGemHelper(null, i);
		}

		return curItem;
	}

	asActionId(): ActionId {
		return ActionId.fromItemId(this._item.id);
	}

	asSpec(): ItemSpec {
		return ItemSpec.create({
			id: this._item.id,
			enchant: this._enchant?.effectId,
			gems: this._gems.map(gem => gem?.id || 0),
		});
	}

	meetsSocketBonus(): boolean {
		return this._item.gemSockets.every((socketColor, i) => this._gems[i] && gemMatchesSocket(this._gems[i]!, socketColor));
	}

	socketBonusStats(): Stats {
		if (this.meetsSocketBonus()) {
			return new Stats(this._item.socketBonus);
		} else {
			return new Stats();
		}
	}

	// Whether this item could have an extra socket, assuming Blacksmithing.
	couldHaveExtraSocket(): boolean {
		return [ItemType.ItemTypeWaist, ItemType.ItemTypeWrist, ItemType.ItemTypeHands].includes(this.item.type);
	}

	requiresExtraSocket(): boolean {
		return [ItemType.ItemTypeWrist, ItemType.ItemTypeHands].includes(this.item.type)
			&& this.hasExtraGem()
			&& this._gems[this._gems.length - 1] != null;
	}

	hasExtraSocket(isBlacksmithing: boolean): boolean {
		return this.item.type == ItemType.ItemTypeWaist ||
			(isBlacksmithing && [ItemType.ItemTypeWrist, ItemType.ItemTypeHands].includes(this.item.type));
	}

	numSockets(isBlacksmithing: boolean): number {
		return this._item.gemSockets.length + (this.hasExtraSocket(isBlacksmithing) ? 1 : 0);
	}

	numSocketsOfColor(color: GemColor | null): number {
		let numSockets: number = 0;

		for (var socketColor of this._item.gemSockets) {
			if (socketColor == color) {
				numSockets += 1;
			}
		}

		return numSockets;
	}

	hasExtraGem(): boolean {
		return this._gems.length > this.item.gemSockets.length;
	}

	hasSocketedGem(socketIdx: number): boolean {
		return this._gems[socketIdx] != null;
	}

	allSocketColors(): Array<GemColor> {
		return this.couldHaveExtraSocket() ? this._item.gemSockets.concat([GemColor.GemColorPrismatic]) : this._item.gemSockets;
	}
	curSocketColors(isBlacksmithing: boolean): Array<GemColor> {
		return this.hasExtraSocket(isBlacksmithing) ? this._item.gemSockets.concat([GemColor.GemColorPrismatic]) : this._item.gemSockets;
	}

	curGems(isBlacksmithing: boolean): Array<Gem|null> {
		return this._gems.slice(0, this.numSockets(isBlacksmithing));
	}
	curEquippedGems(isBlacksmithing: boolean): Array<Gem> {
		return this.curGems(isBlacksmithing).filter(g => g != null) as Array<Gem>;
	}

	getProfessionRequirements(): Array<Profession> {
		let profs: Array<Profession> = [];
		if (this._item.requiredProfession != Profession.ProfessionUnknown) {
			profs.push(this._item.requiredProfession);
		}
		if (this._enchant != null && this._enchant.requiredProfession != Profession.ProfessionUnknown) {
			profs.push(this._enchant.requiredProfession);
		}
		this._gems.forEach(gem => {
			if (gem != null && gem.requiredProfession != Profession.ProfessionUnknown) {
				profs.push(gem.requiredProfession);
			}
		});
		if (this.requiresExtraSocket()) {
			profs.push(Profession.Blacksmithing);
		}
		return distinct(profs);
	}
	getFailedProfessionRequirements(professions: Array<Profession>): Array<Item | Gem | Enchant> {
		let failed: Array<Item | Gem | Enchant> = [];
		if (this._item.requiredProfession != Profession.ProfessionUnknown && !professions.includes(this._item.requiredProfession)) {
			failed.push(this._item);
		}
		if (this._enchant != null && this._enchant.requiredProfession != Profession.ProfessionUnknown && !professions.includes(this._enchant.requiredProfession)) {
			failed.push(this._enchant);
		}
		this._gems.forEach(gem => {
			if (gem != null && gem.requiredProfession != Profession.ProfessionUnknown && !professions.includes(gem.requiredProfession)) {
				failed.push(gem);
			}
		});
		return failed;
	}
};
