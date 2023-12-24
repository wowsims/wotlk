import {
	Class,
	EquipmentSpec,
	ItemSlot,
	ItemSpec,
	ItemSwap,
	PresetEncounter,
	PresetTarget,
	SimDatabase,
} from '../proto/common.js';
import {
	UIEnchant as Enchant,
	IconData,
	UIItem as Item,
	UINPC as Npc,
	UIRune as Rune,
	UIDatabase,
	UIZone as Zone,
} from '../proto/ui.js';

import { MAX_CHARACTER_LEVEL } from '../constants/mechanics.js';
import { EquippedItem } from './equipped_item.js';
import { Gear, ItemSwapGear } from './gear.js';
import {
	getEligibleEnchantSlots,
	getEligibleItemSlots,
	itemTypeToSlotsMap,
} from './utils.js';
import { distinct } from '../utils.js';

const dbUrlJson = '/sod/assets/database/db.json';
const dbUrlBin = '/sod/assets/database/db.bin';
const leftoversUrlJson = '/sod/assets/database/leftover_db.json';
const leftoversUrlBin = '/sod/assets/database/leftover_db.bin';
// When changing this value, don't forget to change the html <link> for preloading!
const READ_JSON = true;
const RANK_REGEX = /Rank ([0-9]+)/g;
const REQ_LEVEL_ITEMS_REGEX = /\<!\-\-rlvl\-\-\>([0-9]+)/g;
const REQ_LEVEL_SPELLS_REGEX = /Requires level ([0-9]+)/g

export class Database {
	private static loadPromise: Promise<Database> | null = null;
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
	static async loadLeftoversIfNecessary(equipment: EquipmentSpec): Promise<Database> {
		const db = await Database.get();
		if (db.loadedLeftovers) {
			return db;
		}

		const shouldLoadLeftovers = equipment.items.some(item => item.id != 0 && !db.items.has(item.id));
		if (shouldLoadLeftovers) {
			const leftoverDb = await Database.getLeftovers();
			db.loadProto(leftoverDb);
			db.loadedLeftovers = true;
		}
		return db;
	}

	private readonly items = new Map<number, Item>();
	private readonly enchantsBySlot: Partial<Record<ItemSlot, Enchant[]>> = {};
	private readonly runesBySlotByClass: Partial<Record<ItemSlot, Partial<Record<Class, Rune[]>>>> = {};
	private readonly runesById: Record<number, Rune> = {};
	private readonly npcs = new Map<number, Npc>();
	private readonly zones = new Map<number, Zone>();
	private readonly presetEncounters = new Map<string, PresetEncounter>();
	private readonly presetTargets = new Map<string, PresetTarget>();
	private readonly itemIcons: Record<number, Promise<IconData>> = {};
	private readonly spellIcons: Record<number, Promise<IconData>> = {};
	private loadedLeftovers: boolean = false;

	private constructor(db: UIDatabase) {
		this.loadProto(db);
	}

	// Add all data from the db proto into this database.
	private loadProto(db: UIDatabase) {
		db.items.forEach(item => this.items.set(item.id, item));
		db.enchants.forEach(enchant => {
			const slots = getEligibleEnchantSlots(enchant);
			slots.forEach(slot => {
				if (!this.enchantsBySlot[slot]) {
					this.enchantsBySlot[slot] = [];
				}
				this.enchantsBySlot[slot]!.push(enchant);
			});
		});
		db.runes.forEach(rune => {
			this.runesById[rune.id] = rune;

			const slots = itemTypeToSlotsMap[rune.type];
			slots?.forEach(slot => {
				if (!this.runesBySlotByClass[slot]){
					this.runesBySlotByClass[slot] = {};
				}
				if (!this.runesBySlotByClass[slot]![rune.class]){
					this.runesBySlotByClass[slot]![rune.class as Class] = [];
				}
				this.runesBySlotByClass[slot]![rune.class as Class]!.push(rune);
			});
		});

		db.npcs.forEach(npc => this.npcs.set(npc.id, npc));
		db.zones.forEach(zone => this.zones.set(zone.id, zone));
		db.encounters.forEach(encounter => this.presetEncounters.set(encounter.path, encounter));
		db.encounters.map(e => e.targets).flat().forEach(target => this.presetTargets.set(target.path, target));

		db.items.forEach(item => this.itemIcons[item.id] = Promise.resolve(IconData.create({
			id: item.id,
			name: item.name,
			icon: item.icon,
		})));

		db.itemIcons.forEach(data => this.itemIcons[data.id] = Promise.resolve(data));
		db.spellIcons.forEach(data => this.spellIcons[data.id] = Promise.resolve(data));
	}

	getAllItems(): Array<Item> {
		return Array.from(this.items.values());
	}

	getItems(slot: ItemSlot): Array<Item> {
		return this.getAllItems().filter(item => getEligibleItemSlots(item).includes(slot));
	}

	getItemById(id: number): Item | undefined {
		return this.items.get(id);
	}

	getEnchants(slot: ItemSlot): Array<Enchant> {
		return this.enchantsBySlot[slot] || [];
	}

	getRunes(slot: ItemSlot, klass: Class): Array<Rune> {
		if (!this.runesBySlotByClass[slot]) return [];

		return this.runesBySlotByClass[slot]![klass] || [];
	}

	hasRuneBySlot(slot: ItemSlot, klass: Class): boolean {
		return !!(this.runesBySlotByClass[slot] && this.runesBySlotByClass[slot]![klass]);
	}

	getNpc(npcId: number): Npc | null {
		return this.npcs.get(npcId) || null;
	}
	getZone(zoneId: number): Zone | null {
		return this.zones.get(zoneId) || null;
	}

	lookupItemSpec(itemSpec: ItemSpec): EquippedItem | null {
		const item = this.items.get(itemSpec.id);
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
		
		let rune: Rune | undefined
		if (itemSpec.rune && !!this.runesById[itemSpec.rune]) {
			rune = this.runesById[itemSpec.rune];
		}

		return new EquippedItem({item, enchant, rune});
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

	lookupItemSwap(itemSwap: ItemSwap): ItemSwapGear {
		return new ItemSwapGear({
			[ItemSlot.ItemSlotMainHand]: itemSwap.mhItem ? this.lookupItemSpec(itemSwap.mhItem): null,
			[ItemSlot.ItemSlotOffHand]: itemSwap.ohItem ? this.lookupItemSpec(itemSwap.ohItem): null,
			[ItemSlot.ItemSlotRanged]: itemSwap.rangedItem ? this.lookupItemSpec(itemSwap.rangedItem): null,
		});
	}

	enchantSpellIdToEffectId(enchantSpellId: number): number {
		const enchant = Object.values(this.enchantsBySlot).flat().find(enchant => enchant.spellId == enchantSpellId);
		return enchant ? enchant.effectId : 0;
	}

	getPresetEncounter(path: string): PresetEncounter | null {
		return this.presetEncounters.get(path) || null;
	}
	getPresetTarget(path: string): PresetTarget | null {
		return this.presetTargets.get(path) || null;
	}
	getAllPresetEncounters(): Array<PresetEncounter> {
		return Array.from(this.presetEncounters.values());
	}
	getAllPresetTargets(): Array<PresetTarget> {
		return Array.from(this.presetTargets.values());
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

	static async getSpellRankData(spellId: number): Promise<IconData> {
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
		if (id === 0) return IconData.create();

		const url = `https://nether.wowhead.com/classic/tooltip/${tooltipPostfix}/${id}?lvl=${MAX_CHARACTER_LEVEL}`;
		try {
			const response = await fetch(url);
			const json = await response.json();
			let rank: number = 0
			let reqLevel: number = 0;

			if (tooltipPostfix === 'spell'){
				const rankMatches = Array.from(json['tooltip'].matchAll(RANK_REGEX) as RegExpMatchArray[]);
				const levelMatches = Array.from(json['tooltip'].matchAll(REQ_LEVEL_SPELLS_REGEX) as RegExpMatchArray[]);
				rank = rankMatches.length ? parseInt(rankMatches[0][1]) : 0;
				reqLevel = levelMatches.length ? parseInt(levelMatches[0][1]): 0;
			} else if (tooltipPostfix == 'item'){
				const levelMatches = Array.from(json['tooltip'].matchAll(REQ_LEVEL_ITEMS_REGEX) as RegExpMatchArray[]);
				reqLevel = levelMatches.length ? parseInt(levelMatches[0][1]): 0;
			}
			
			return IconData.create({
				id: id,
				name: json['name'],
				icon: json['icon'],
				rank: rank,
				requiresLevel: reqLevel,
			});
		} catch (e) {
			console.error('Error while fetching url: ' + url + '\n\n' + e);
			return IconData.create();
		}
	}

	public static mergeSimDatabases(db1: SimDatabase, db2: SimDatabase): SimDatabase {
		return SimDatabase.create({
			items: distinct(db1.items.concat(db2.items), (a, b) => a.id == b.id),
			enchants: distinct(db1.enchants.concat(db2.enchants), (a, b) => a.effectId == b.effectId),
		})
	}
}
