import { EquipmentSpec, ItemSwap } from '../proto/common.js';
import { GemColor } from '../proto/common.js';
import { ItemSlot } from '../proto/common.js';
import { ItemSpec } from '../proto/common.js';
import { Profession } from '../proto/common.js';
import { SimDatabase } from '../proto/common.js';
import { SimItem } from '../proto/common.js';
import { SimEnchant } from '../proto/common.js';
import { SimGem } from '../proto/common.js';
import { equalsOrBothNull } from '../utils.js';
import { distinct, getEnumValues } from '../utils.js';
import { isBluntWeaponType, isSharpWeaponType } from '../proto_utils/utils.js';
import {
	UIEnchant as Enchant,
	UIGem as Gem,
	UIItem as Item,
} from '../proto/ui.js';

import { isMetaGemActive } from './gems.js';
import { gemMatchesSocket } from './gems.js';
import { EquippedItem } from './equipped_item.js';
import { validWeaponCombo } from './utils.js';
import { Stats } from './stats.js';

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
		const otherArray = other.asArray();
		return this.asArray().every((thisItem, slot) => equalsOrBothNull(thisItem, otherArray[slot], (a, b) => a.equals(b)))
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
			this.removeUniqueGems(newInternalGear, newItem);
			this.removeUniqueItems(newInternalGear, newItem);
		}

		// Actually assign the new item.
		newInternalGear[newSlot] = newItem;

		BaseGear.validateWeaponCombo(newInternalGear, newSlot, canDualWield2H);

		return newInternalGear;
	}

	private removeUniqueGems(gear: Partial<InternalGear>, newItem: EquippedItem) {
		// If the new item has unique gems, remove matching.
		newItem.gems
			.filter(gem => gem?.unique)
			.forEach(gem => {
				this.getItemSlots().map(slot => Number(slot) as ItemSlot).forEach(slot => {
					gear[slot] = gear[slot]?.removeGemsWithId(gem!.id) || null;
				});
			});
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
			gems: distinct(equippedItems.map(ei => (ei._gems.filter(g => g != null) as Array<Gem>).map(gem => BaseGear.gemToDB(gem))).flat()),
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

/**
 * Represents a full gear set, including items/enchants/gems for every slot.
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

	getAllGems(isBlacksmithing: boolean): Array<Gem> {
		return this.asArray()
			.map(ei => ei == null ? [] : ei.curEquippedGems(isBlacksmithing))
			.flat();
	}

	getNonMetaGems(isBlacksmithing: boolean): Array<Gem> {
		return this.getAllGems(isBlacksmithing).filter(gem => gem.color != GemColor.GemColorMeta);
	}

	statsFromGems(isBlacksmithing: boolean): Stats {
		let stats = new Stats();

		// Stats from just the gems.
		const gems = this.getAllGems(isBlacksmithing);
		for (let i = 0; i < gems.length; i++) {
			stats = stats.add(new Stats(gems[i].stats));
		}

		// Stats from socket bonuses.
		const items = this.asArray().filter(ei => ei != null) as Array<EquippedItem>;
		for (let i = 0; i < items.length; i++) {
			stats = stats.add(items[i].socketBonusStats());
		}

		return stats;
	}

	getGemsOfColor(color: GemColor, isBlacksmithing: boolean): Array<Gem> {
		return this.getAllGems(isBlacksmithing).filter(gem => gem.color == color);
	}

	getJCGems(isBlacksmithing: boolean): Array<Gem> {
		return this.getAllGems(isBlacksmithing).filter(gem => gem.requiredProfession == Profession.Jewelcrafting);
	}

	getMetaGem(): Gem | null {
		return this.getGemsOfColor(GemColor.GemColorMeta, true)[0] || null;
	}

	gemColorCounts(isBlacksmithing: boolean): ({ red: number, yellow: number, blue: number }) {
		const gems = this.getAllGems(isBlacksmithing);
		return {
			red: gems.filter(gem => gemMatchesSocket(gem, GemColor.GemColorRed)).length,
			yellow: gems.filter(gem => gemMatchesSocket(gem, GemColor.GemColorYellow)).length,
			blue: gems.filter(gem => gemMatchesSocket(gem, GemColor.GemColorBlue)).length,
		};
	}

	// Returns true if this gear set has a meta gem AND the other gems meet the meta's conditions.
	hasActiveMetaGem(isBlacksmithing: boolean): boolean {
		const metaGem = this.getMetaGem();
		if (!metaGem) {
			return false;
		}

		const gemColorCounts = this.gemColorCounts(isBlacksmithing);
		return isMetaGemActive(
			metaGem,
			gemColorCounts.red, gemColorCounts.yellow, gemColorCounts.blue);
	}

	hasInactiveMetaGem(isBlacksmithing: boolean): boolean {
		return this.getMetaGem() != null && !this.hasActiveMetaGem(isBlacksmithing);
	}

	withGem(itemSlot: ItemSlot, socketIdx: number, gem: Gem | null): Gear {
		const item = this.getEquippedItem(itemSlot);

		if (item) {
			return this.withEquippedItem(itemSlot, item.withGem(gem, socketIdx), true);
		}

		return this;
	}

	withSingleGemSubstitution(oldGem: Gem | null, newGem: Gem | null, isBlacksmithing: boolean): Gear {
		for (var slot of this.getItemSlots()) {
			const item = this.getEquippedItem(slot);

			if (!item) {
				continue;
			}

			const currentGems = item!.curGems(isBlacksmithing);

			if (currentGems.includes(oldGem)) {
				const socketIdx = currentGems.indexOf(oldGem);
				return this.withGem(slot, socketIdx, newGem);
			}
		}

		return this;
	}

	withMetaGem(metaGem: Gem | null): Gear {
		const headItem = this.getEquippedItem(ItemSlot.ItemSlotHead);

		if (headItem) {
			for (const [socketIdx, socketColor] of headItem.allSocketColors().entries()) {
				if (socketColor == GemColor.GemColorMeta) {
					return this.withEquippedItem(ItemSlot.ItemSlotHead, headItem.withGem(metaGem, socketIdx), true);
				}
			}
		}

		return this;
	}

	withoutMetaGem(): Gear {
		const headItem = this.getEquippedItem(ItemSlot.ItemSlotHead);
		const metaGem = this.getMetaGem();
		if (headItem && metaGem) {
			return this.withEquippedItem(ItemSlot.ItemSlotHead, headItem.removeGemsWithId(metaGem.id), true);
		} else {
			return this;
		}
	}

	withoutGems(): Gear {
		let curGear: Gear = this;

		for (var slot of this.getItemSlots()) {
			const item = this.getEquippedItem(slot);

			if (item) {
				curGear = curGear.withEquippedItem(slot, item.removeAllGems(), true);
			}
		}

		return curGear;
	}

	// Removes bonus gems from blacksmith profession bonus.
	withoutBlacksmithSockets(): Gear {
		let curGear: Gear = this;

		const wristItem = this.getEquippedItem(ItemSlot.ItemSlotWrist);
		if (wristItem) {
			curGear = curGear.withEquippedItem(ItemSlot.ItemSlotWrist, wristItem.withGem(null, wristItem.numPossibleSockets - 1), true);
		}

		const handsItem = this.getEquippedItem(ItemSlot.ItemSlotHands);
		if (handsItem) {
			curGear = curGear.withEquippedItem(ItemSlot.ItemSlotHands, handsItem.withGem(null, handsItem.numPossibleSockets - 1), true);
		}

		return curGear;
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
	getFailedProfessionRequirements(professions: Array<Profession>): Array<Item | Gem | Enchant> {
		return (this.asArray().filter(ei => ei != null) as Array<EquippedItem>)
			.map(ei => ei.getFailedProfessionRequirements(professions))
			.flat();
	}
}

/**
 * Represents a item swap gear set, including items/enchants/gems.
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
