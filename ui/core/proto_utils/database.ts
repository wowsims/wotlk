import {
	EquipmentSpec,
	GemColor,
	ItemSlot,
	ItemSpec,
	PresetEncounter,
	PresetTarget,
} from '../proto/common.js';
import {
	IconData,
	UIDatabase,
	UIEnchant as Enchant,
	UIGem as Gem,
	UIItem as Item,
} from '../proto/ui.js';

import {
	getEligibleEnchantSlots,
	getEligibleItemSlots,
} from './utils.js';
import { gemEligibleForSocket, gemMatchesSocket } from './gems.js';
import { EquippedItem } from './equipped_item.js';
import { Gear } from './gear.js';

const dbUrlJson = '/wotlk/assets/database/db.json';
const dbUrlBin = '/wotlk/assets/database/db.bin';
const leftoversUrlJson = '/wotlk/assets/database/leftover_db.json';
const leftoversUrlBin = '/wotlk/assets/database/leftover_db.bin';
// When changing this value, don't forget to change the html <link> for preloading!
const READ_JSON = true;

export class Database {
	private static loadPromise: Promise<Database>|null = null;
	static get(): Promise<Database> {
		if (Database.loadPromise == null) {
			if (READ_JSON) {
				Database.loadPromise = fetch(dbUrlJson)
					.then(response => response.json())
					.then(json => new Database(UIDatabase.fromJson(json)));
			} else {
				Database.loadPromise = fetch(dbUrlBin)
					.then(response => response.arrayBuffer())
					.then(buffer => new Database(UIDatabase.fromBinary(new Uint8Array(buffer))));
			}
		}
		return Database.loadPromise;
	}

	static getLeftovers(): Promise<UIDatabase> {
		if (READ_JSON) {
			return fetch(leftoversUrlJson)
				.then(response => response.json())
				.then(json => UIDatabase.fromJson(json));
		} else {
			return fetch(leftoversUrlBin)
				.then(response => response.arrayBuffer())
				.then(buffer => UIDatabase.fromBinary(new Uint8Array(buffer)));
		}
	}

	// Checks if any items in the equipment are missing from the current DB. If so, loads the leftover DB.
	static async loadLeftoversIfNecessary(equipment: EquipmentSpec): Promise<void> {
		const db = await Database.get();
		if (db.loadedLeftovers) {
			return;
		}

		const shouldLoadLeftovers = equipment.items.some(item => item.id != 0 && !db.items[item.id]);
		if (shouldLoadLeftovers) {
			const leftoverDb = await Database.getLeftovers();
			db.loadProto(leftoverDb);
			db.loadedLeftovers = true;
		}
	}

	private readonly items: Record<number, Item> = {};
	private readonly enchantsBySlot: Partial<Record<ItemSlot, Enchant[]>> = {};
	private readonly gems: Record<number, Gem> = {};
	private readonly presetEncounters: Record<string, PresetEncounter> = {};
	private readonly presetTargets: Record<string, PresetTarget> = {};
	private readonly itemIcons: Record<number, Promise<IconData>>;
	private readonly spellIcons: Record<number, Promise<IconData>>;
	private loadedLeftovers: boolean = false;

	private constructor(db: UIDatabase) {
		this.itemIcons = {};
		this.spellIcons = {};
		this.loadProto(db);
	}

	// Add all data from the db proto into this database.
	private loadProto(db: UIDatabase) {
		db.items.forEach(item => this.items[item.id] = item);
		db.enchants.forEach(enchant => {
			const slots = getEligibleEnchantSlots(enchant);
			slots.forEach(slot => {
				if (!this.enchantsBySlot[slot]) {
					this.enchantsBySlot[slot] = [];
				}
				this.enchantsBySlot[slot]!.push(enchant);
			});
		});
		db.gems.forEach(gem => this.gems[gem.id] = gem);

		db.encounters.forEach(encounter => this.presetEncounters[encounter.path] = encounter);
		db.encounters.map(e => e.targets).flat().forEach(target => this.presetTargets[target.path] = target);

		db.items.forEach(item => this.itemIcons[item.id] = new Promise((resolve, _) => resolve(IconData.create({
			id: item.id,
			name: item.name,
			icon: item.icon,
		}))));
		db.gems.forEach(gem => this.itemIcons[gem.id] = new Promise((resolve, _) => resolve(IconData.create({
			id: gem.id,
			name: gem.name,
			icon: gem.icon,
		}))));
		db.itemIcons.forEach(data => this.itemIcons[data.id] = new Promise((resolve, _) => resolve(data)));
		db.spellIcons.forEach(data => this.spellIcons[data.id] = new Promise((resolve, _) => resolve(data)));
	}

	getItems(slot: ItemSlot): Array<Item> {
		let items = Object.values(this.items);
		items = items.filter(item => getEligibleItemSlots(item).includes(slot));
		return items;
	}

	getEnchants(slot: ItemSlot): Array<Enchant> {
		return this.enchantsBySlot[slot] || [];
	}

	getGems(socketColor?: GemColor): Array<Gem> {
		let gems = Object.values(this.gems);
		if (socketColor) {
			gems = gems.filter(gem => gemEligibleForSocket(gem, socketColor));
		}
		return gems;
	}

	getMatchingGems(socketColor: GemColor): Array<Gem> {
		return Object.values(this.gems).filter(gem => gemMatchesSocket(gem, socketColor));
	}

	lookupItemSpec(itemSpec: ItemSpec): EquippedItem | null {
		const item = this.items[itemSpec.id];
		if (!item)
			return null;

		let enchant: Enchant | null = null;
		if (itemSpec.enchant) {
			const slots = getEligibleItemSlots(item);
			for (let i = 0; i < slots.length; i++) {
				enchant = (this.enchantsBySlot[slots[i]] || [])
						.find(enchant => [enchant.effectId, enchant.itemId, enchant.spellId].includes(itemSpec.enchant)) || null;
				if (enchant) {
					break;
				}
			}
		}

		const gems = itemSpec.gems.map(gemId => this.gems[gemId] || null);

		return new EquippedItem(item, enchant, gems);
	}

	lookupEquipmentSpec(equipSpec: EquipmentSpec): Gear {
		// EquipmentSpec is supposed to be indexed by slot, but here we assume
		// it isn't just in case.
		const gearMap: Partial<Record<ItemSlot, EquippedItem | null>> = {};

		equipSpec.items.forEach(itemSpec => {
			const item = this.lookupItemSpec(itemSpec);
			if (!item)
				return;

			const itemSlots = getEligibleItemSlots(item.item);

			const assignedSlot = itemSlots.find(slot => !gearMap[slot]);
			if (assignedSlot == null)
				throw new Error('No slots left to equip ' + Item.toJsonString(item.item));

			gearMap[assignedSlot] = item;
		});

		return new Gear(gearMap);
	}

	getPresetEncounter(path: string): PresetEncounter | null {
		return this.presetEncounters[path] || null;
	}
	getPresetTarget(path: string): PresetTarget | null {
		return this.presetTargets[path] || null;
	}
	getAllPresetEncounters(): Array<PresetEncounter> {
		return Object.values(this.presetEncounters);
	}
	getAllPresetTargets(): Array<PresetTarget> {
		return Object.values(this.presetTargets);
	}

	static async getItemIconData(itemId: number): Promise<IconData> {
		const db = await Database.get();
		if (!db.itemIcons[itemId]) {
			db.itemIcons[itemId] = Database.getWowheadItemTooltipData(itemId);
		}
		return await db.itemIcons[itemId];
	}

	static async getSpellIconData(spellId: number): Promise<IconData> {
		const db = await Database.get();
		if (!db.spellIcons[spellId]) {
			db.spellIcons[spellId] = Database.getWowheadSpellTooltipData(spellId);
		}
		return await db.spellIcons[spellId];
	}

	private static async getWowheadItemTooltipData(id: number): Promise<IconData> {
		return Database.getWowheadTooltipData(id, 'item');
	}
	private static async getWowheadSpellTooltipData(id: number): Promise<IconData> {
		return Database.getWowheadTooltipData(id, 'spell');
	}
	private static async getWowheadTooltipData(id: number, tooltipPostfix: string): Promise<IconData> {
		const url = `https://nether.wowhead.com/wotlk/tooltip/${tooltipPostfix}/${id}`;
		try {
			const response = await fetch(url);
			const json = await response.json();
			return IconData.create({
				id: id,
				name: json['name'],
				icon: json['icon'],
			});
		} catch (e) {
			console.error('Error while fetching url: ' + url + '\n\n' + e);
			return IconData.create();
		}
	}
}
