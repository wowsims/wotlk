import { ItemSpec, Profession } from '../proto/common.js';
import {
	UIEnchant as Enchant,
	UIItem as Item,
	UIRune as Rune,
} from '../proto/ui.js';
import { distinct } from '../utils.js';
import { ActionId } from './action_id.js';
import { enchantAppliesToItem } from './utils.js';

export function getWeaponDPS(item: Item): number {
	return ((item.weaponDamageMin + item.weaponDamageMax) / 2) / (item.weaponSpeed || 1);
}

interface EquippedItemConfig {
	item: Item, 
	enchant?: Enchant | null, 
	rune?: Rune | null,
}

/**
 * Represents an equipped item along with enchants attached to it.
 *
 * This is an immutable type.
 */
export class EquippedItem {
	readonly _item: Item;
	readonly _enchant: Enchant | null;
	readonly _rune: Rune | null;

	constructor(config: EquippedItemConfig) {
		this._item = config.item;
		this._enchant = config.enchant || null;
		this._rune = config.rune || null;
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

	get rune(): Rune | null {
		// Make a defensive copy
		return this._rune ? Rune.clone(this._rune) : null;
	}

	equals(other: EquippedItem) {
		if (!Item.equals(this._item, other.item))
			return false;

		if ((this._enchant == null) != (other.enchant == null))
			return false;

		if (this._enchant && other.enchant && !Enchant.equals(this._enchant, other.enchant))
			return false;

		if ((this._rune == null) != (other.rune == null))
			return false;

		if (this._rune && other.rune && !Rune.equals(this._rune, other.rune))
			return false;

		return true;
	}

	/**
	 * Replaces the item and tries to keep the existing enchants if possible.
	 */
	withItem(item: Item): EquippedItem {
		let newEnchant = null;
		if (this._enchant && enchantAppliesToItem(this._enchant, item))
			newEnchant = this._enchant;

		return new EquippedItem({item, enchant: newEnchant, rune: this.rune});
	}

	/**
	 * Returns a new EquippedItem with the given enchant applied.
	 */
	withEnchant(enchant: Enchant | null): EquippedItem {
		return new EquippedItem({item: this._item, enchant, rune: this._rune});
	}

	withRune(rune: Rune | null): EquippedItem {
		return new EquippedItem({item: this._item, enchant: this.enchant, rune});
	}

	asActionId(): ActionId {
		return ActionId.fromItemId(this._item.id);
	}

	asSpec(): ItemSpec {
		return ItemSpec.create({
			id: this._item.id,
			enchant: this._enchant?.effectId,
			rune: this._rune?.id,
		});
	}

	getProfessionRequirements(): Array<Profession> {
		let profs: Array<Profession> = [];
		if (this._item.requiredProfession != Profession.ProfessionUnknown) {
			profs.push(this._item.requiredProfession);
		}
		if (this._enchant != null && this._enchant.requiredProfession != Profession.ProfessionUnknown) {
			profs.push(this._enchant.requiredProfession);
		}
		return distinct(profs);
	}
	getFailedProfessionRequirements(professions: Array<Profession>): Array<Item | Enchant> {
		let failed: Array<Item | Enchant> = [];
		if (this._item.requiredProfession != Profession.ProfessionUnknown && !professions.includes(this._item.requiredProfession)) {
			failed.push(this._item);
		}
		if (this._enchant != null && this._enchant.requiredProfession != Profession.ProfessionUnknown && !professions.includes(this._enchant.requiredProfession)) {
			failed.push(this._enchant);
		}
		return failed;
	}
};
