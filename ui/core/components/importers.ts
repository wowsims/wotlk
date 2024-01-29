import pako from 'pako';

import { IndividualSimUI } from '../individual_sim_ui';
import {
	Class,
	EquipmentSpec,
	ItemSlot,
	ItemSpec,
	Profession,
	Race,
	Spec,
} from '../proto/common';
import { IndividualSimSettings } from '../proto/ui';
import { Database } from '../proto_utils/database';
import { classNames, nameToClass, nameToRace, nameToProfession } from '../proto_utils/names';
import { SimSettingCategories } from '../sim';
import { SimUI } from '../sim_ui';
import { talentSpellIdsToTalentString } from '../talents/factory';
import { TypedEvent } from '../typed_event';
import { buf2hex, getEnumValues } from '../utils';

import { BaseModal } from './base_modal';

export abstract class Importer extends BaseModal {
	protected readonly textElem: HTMLTextAreaElement;
	protected readonly descriptionElem: HTMLElement;
	protected readonly importButton: HTMLButtonElement;
	private readonly includeFile: boolean;

	constructor(parent: HTMLElement, simUI: SimUI, title: string, includeFile: boolean) {
		super(parent, 'importer', { title: title, footer: true });
		this.includeFile = includeFile;
		const uploadInputId = 'upload-input-' + title.toLowerCase().replaceAll(' ', '-');

		this.body.innerHTML = `
			<div class="import-description"></div>
			<textarea spellCheck="false" class="importer-textarea form-control"></textarea>
		`;
		this.footer!.innerHTML = `
			${this.includeFile ? `
				<label for="${uploadInputId}" class="importer-button btn btn-primary upload-button me-2">
					<i class="fas fa-file-arrow-up"></i>
					Upload File
				</label>
				<input type="file" id="${uploadInputId}" class="importer-upload-input d-none" hidden>
			` : ''
			}
			<button class="importer-button btn btn-primary import-button">
				<i class="fa fa-download"></i>
				Import
			</button>
		`;

		this.textElem = this.rootElem.getElementsByClassName('importer-textarea')[0] as HTMLTextAreaElement;
		this.descriptionElem = this.rootElem.getElementsByClassName('import-description')[0] as HTMLElement;

		if (this.includeFile) {
			const uploadInput = this.rootElem.getElementsByClassName('importer-upload-input')[0] as HTMLButtonElement;
			uploadInput.addEventListener('change', async event => {
				const data: string = await (event as any).target.files[0].text();
				this.textElem.textContent = data;
			});
		}

		this.importButton = this.rootElem.getElementsByClassName('import-button')[0] as HTMLButtonElement;
		this.importButton.addEventListener('click', () => {
			try {
				this.onImport(this.textElem.value || '');
			} catch (error) {
				alert('Import error: ' + error);
			}
		});
	}

	abstract onImport(data: string): void

	protected async finishIndividualImport<SpecType extends Spec>(simUI: IndividualSimUI<SpecType>, charClass: Class, race: Race, equipmentSpec: EquipmentSpec, talentsStr: string, professions: Array<Profession>): Promise<void> {
		const playerClass = simUI.player.getClass();
		if (charClass != playerClass) {
			throw new Error(`Wrong Class! Expected ${classNames.get(playerClass)} but found ${classNames.get(charClass)}!`);
		}

		await Database.loadLeftoversIfNecessary(equipmentSpec);

		const gear = simUI.sim.db.lookupEquipmentSpec(equipmentSpec);

		const expectedEnchantIds = equipmentSpec.items.map(item => item.enchant);
		const foundEnchantIds = gear.asSpec().items.map(item => item.enchant);
		const missingEnchants = expectedEnchantIds.filter(expectedId => expectedId != 0 && !foundEnchantIds.includes(expectedId));

		const expectedItemIds = equipmentSpec.items.map(item => item.id);
		const foundItemIds = gear.asSpec().items.map(item => item.id);
		const missingItems = expectedItemIds.filter(expectedId => !foundItemIds.includes(expectedId));

		// Now update settings using the parsed values.
		const eventID = TypedEvent.nextEventID();
		TypedEvent.freezeAllAndDo(() => {
			simUI.player.setRace(eventID, race);
			simUI.player.setGear(eventID, gear);
			if (talentsStr && talentsStr != '--') {
				simUI.player.setTalentsString(eventID, talentsStr);
			}
			if (professions.length > 0) {
				simUI.player.setProfessions(eventID, professions)
			}
		});

		this.close();

		if (missingItems.length == 0 && missingEnchants.length == 0) {
			alert('Import successful!');
		} else {
			alert('Import successful, but the following IDs were not found in the sim database:' +
				(missingItems.length == 0 ? '' : '\n\nItems: ' + missingItems.join(', ')) +
				(missingEnchants.length == 0 ? '' : '\n\nEnchants: ' + missingEnchants.join(', ')));
		}
	}
}

interface UrlParseData {
	settings: IndividualSimSettings,
	categories: Array<SimSettingCategories>,
}

// For now this just holds static helpers to match the exporter, so it doesn't extend Importer.
export class IndividualLinkImporter {

	// Exclude UISettings by default, since most users don't intend to export those.
	static readonly DEFAULT_CATEGORIES = getEnumValues(SimSettingCategories).filter(c => c != SimSettingCategories.UISettings) as Array<SimSettingCategories>;

	static readonly CATEGORY_PARAM = 'i';
	static readonly CATEGORY_KEYS: Map<SimSettingCategories, string> = (() => {
		const map = new Map();
		// Use single-letter abbreviations since these will be included in sim links.
		map.set(SimSettingCategories.Gear, 'g');
		map.set(SimSettingCategories.Talents, 't');
		map.set(SimSettingCategories.Rotation, 'r');
		map.set(SimSettingCategories.Consumes, 'c');
		map.set(SimSettingCategories.Miscellaneous, 'm');
		map.set(SimSettingCategories.External, 'x');
		map.set(SimSettingCategories.Encounter, 'e');
		map.set(SimSettingCategories.UISettings, 'u');
		return map;
	})();

	static tryParseUrlLocation(location: Location): UrlParseData|null {
		let hash = location.hash;
		if (hash.length <= 1) {
			return null;
		}

		// Remove leading '#'
		hash = hash.substring(1);
		const binary = atob(hash);
		const bytes = new Uint8Array(binary.length);
		for (let i = 0; i < bytes.length; i++) {
			bytes[i] = binary.charCodeAt(i);
		}

		const settingsBytes = pako.inflate(bytes);
		const settings = IndividualSimSettings.fromBinary(settingsBytes);

		let exportCategories = IndividualLinkImporter.DEFAULT_CATEGORIES;
		const urlParams = new URLSearchParams(window.location.search);
		if (urlParams.has(IndividualLinkImporter.CATEGORY_PARAM)) {
			const categoryChars = urlParams.get(IndividualLinkImporter.CATEGORY_PARAM)!.split('');
			exportCategories = categoryChars
				.map(char => [...IndividualLinkImporter.CATEGORY_KEYS.entries()]
				.find(e => e[1] == char))
				.filter(e => e)
				.map(e => e![0]);
		}

		return {
			settings: settings,
			categories: exportCategories,
		};
	}
}

export class IndividualJsonImporter<SpecType extends Spec> extends Importer {
	private readonly simUI: IndividualSimUI<SpecType>;
	constructor(parent: HTMLElement, simUI: IndividualSimUI<SpecType>) {
		super(parent, simUI, 'JSON Import', true);
		this.simUI = simUI;

		this.descriptionElem.innerHTML = `
			<p>Import settings from a JSON file, which can be created using the JSON Export feature.</p>
			<p>To import, upload the file or paste the text below, then click, 'Import'.</p>
		`;
	}

	async onImport(data: string) {
		const proto = IndividualSimSettings.fromJsonString(data, { ignoreUnknownFields: true });
		if (proto.player?.equipment) {
			await Database.loadLeftoversIfNecessary(proto.player.equipment);
		}
		if (this.simUI.isWithinRaidSim) {
			if (proto.player) {
				this.simUI.player.fromProto(TypedEvent.nextEventID(), proto.player);
			}
		} else {
			this.simUI.fromProto(TypedEvent.nextEventID(), proto);
		}
		this.close();
	}
}

export class Individual80UImporter<SpecType extends Spec> extends Importer {
	private readonly simUI: IndividualSimUI<SpecType>;
	constructor(parent: HTMLElement, simUI: IndividualSimUI<SpecType>) {
		super(parent, simUI, '80 Upgrades Import', true);
		this.simUI = simUI;

		this.descriptionElem.innerHTML = `
			<p>
				Import settings from <a href="https://eightyupgrades.com" target="_blank">Eighty Upgrades</a>.
			</p>
			<p>
				This feature imports gear, race, and (optionally) talents. It does NOT import buffs, debuffs, consumes, rotation, or custom stats.
			</p>
			<p>
				To import, paste the output from the site's export option below and click, 'Import'.
			</p>
		`;
	}

	onImport(data: string) {
		const importJson = JSON.parse(data);

		// Parse all the settings.
		const charClass = nameToClass((importJson?.character?.gameClass as string) || '');
		if (charClass == Class.ClassUnknown) {
			throw new Error('Could not parse Class!');
		}

		const race = nameToRace((importJson?.character?.race as string) || '');
		if (race == Race.RaceUnknown) {
			throw new Error('Could not parse Race!');
		}

		let talentsStr = '';
		if (importJson?.talents?.length > 0) {
			const talentIds = (importJson.talents as Array<any>).map(talentJson => talentJson.spellId);
			talentsStr = talentSpellIdsToTalentString(charClass, talentIds);
		}

		let equipmentSpec = EquipmentSpec.create();
		(importJson.items as Array<any>).forEach(itemJson => {
			let itemSpec = ItemSpec.create();
			itemSpec.id = itemJson.id;
			if (itemJson.enchant?.id) {
				itemSpec.enchant = itemJson.enchant.id;
			}
			equipmentSpec.items.push(itemSpec);
		});

		this.simUI.sim.db.lookupEquipmentSpec(equipmentSpec);

		this.finishIndividualImport(this.simUI, charClass, race, equipmentSpec, talentsStr, []);
	}
}

export class IndividualWowheadGearPlannerImporter<SpecType extends Spec> extends Importer {
	private readonly simUI: IndividualSimUI<SpecType>;
	constructor(parent: HTMLElement, simUI: IndividualSimUI<SpecType>) {
		super(parent, simUI, 'Wowhead Import', true);
		this.simUI = simUI;

		this.descriptionElem.innerHTML = `
			<p>
				Import settings from <a href="https://www.wowhead.com/classic/gear-planner" target="_blank">Wowhead Gear Planner</a>.
			</p>
			<p>
				This feature imports gear, race, and (optionally) talents. It does NOT import buffs, debuffs, consumes, rotation, or custom stats.
			</p>
			<p>
				To import, paste the gear planner link below and click, 'Import'.
			</p>
		`;
	}

	onImport(url: string) {
		const match = url.match(/www\.wowhead\.com\/classic\/gear-planner\/([a-z\-]+)\/([a-z\-]+)\/([a-zA-Z0-9_\-]+)/);
		if (!match) {
			throw new Error(`Invalid WCL URL ${url}, must look like "https://www.wowhead.com/classic/gear-planner/CLASS/RACE/XXXX"`);
		}

		// Parse all the settings.
		const charClass = nameToClass(match[1].replaceAll('-', ''));
		if (charClass == Class.ClassUnknown) {
			throw new Error('Could not parse Class: ' + match[1]);
		}

		const race = nameToRace(match[2].replaceAll('-', ''));
		if (race == Race.RaceUnknown) {
			throw new Error('Could not parse Race: ' + match[2]);
		}

		const base64Data = match[3].replaceAll('_', '/').replaceAll('-', '+');
		//console.log('Base64: ' + base64Data);
		const data = Uint8Array.from(atob(base64Data), c => c.charCodeAt(0))
		//console.log('Hex: ' + buf2hex(data));

		// Binary schema
		// Byte 00: ??
		// Byte 01: ?? Seems related to aesthetics (e.g. body type)
		// Byte 02: 8-bit Player Level
		// Byte 03: 8-bit length of talents bytes
		// Next N Bytes: Talents in hex string format

		// Talent hex string looks like '230005232100330150323102505321f03f023203001f'
		// Just like regular wowhead talents string except 'f' instead of '-'.
		const numTalentBytes = data[3];
		const talentBytes = data.subarray(4, 4 + numTalentBytes);
		const talentsHexStr = buf2hex(talentBytes);
		//console.log('Talents hex: ' + talentsHexStr);
		const talentsStr = talentsHexStr.split('f').slice(0, 3).join('-');
		//console.log('Talents: ' + talentsStr);

		// Binary schema for each item:
		// 8-bit slotNumber, high bit = is enchanted
		// 8-bit upper 3 bits for gem count
		// 16-bit item id
		// if enchant bit is set:
		//   8-bit ??, possibly enchant position for multiple enchants?
		//   16-bit enchant id
		const gearBytes = data.subarray(numTalentBytes);
		//console.log(`Remaining ${gearBytes.length} bytes: ${buf2hex(gearBytes)}`);
		
		const equipmentSpec = EquipmentSpec.create();
		let cur = 0;
		while (cur < gearBytes.length) {
			const itemSpec = ItemSpec.create();
			const slotId = gearBytes[cur] & 0b00111111;
			const isEnchanted = Boolean(gearBytes[cur] & 0b10000000);
			// const randomEnchant = Boolean(gearBytes[cur] & 0b01000000);
			cur++;

			const highid = (gearBytes[cur] & 0b00011111);
			cur++;

			itemSpec.id = (highid << 16) + (gearBytes[cur] << 8) + gearBytes[cur + 1];
			cur += 2;
			//console.log(`Slot ID: ${slotId}, isEnchanted: ${isEnchanted}, itemID: ${itemSpec.id}`);

			if (isEnchanted) {
				// Note: this is the enchant SPELL id, not the effect ID.
				const enchantSpellId = (gearBytes[cur] << 16) + (gearBytes[cur + 1] << 8) + gearBytes[cur + 2];
				itemSpec.enchant = this.simUI.sim.db.enchantSpellIdToEffectId(enchantSpellId);
				cur += 3;
				//console.log(`Enchant ID: ${itemSpec.enchant}. Spellid: ${enchantSpellId}`);
			}

			// Ignore tabard / shirt slots
			const itemSlotEntry = Object.entries(IndividualWowheadGearPlannerImporter.slotIDs).find(e => e[1] == slotId);
			if (itemSlotEntry != null) {
				equipmentSpec.items.push(itemSpec);
			}
		}
		this.simUI.sim.db.lookupEquipmentSpec(equipmentSpec);

		this.finishIndividualImport(this.simUI, charClass, race, equipmentSpec, talentsStr, []);
	}

	static slotIDs: Record<ItemSlot, number> = {
		[ItemSlot.ItemSlotHead]: 1,
		[ItemSlot.ItemSlotNeck]: 2,
		[ItemSlot.ItemSlotShoulder]: 3,
		[ItemSlot.ItemSlotBack]: 15,
		[ItemSlot.ItemSlotChest]: 5,
		[ItemSlot.ItemSlotWrist]: 9,
		[ItemSlot.ItemSlotHands]: 10,
		[ItemSlot.ItemSlotWaist]: 6,
		[ItemSlot.ItemSlotLegs]: 7,
		[ItemSlot.ItemSlotFeet]: 8,
		[ItemSlot.ItemSlotFinger1]: 11,
		[ItemSlot.ItemSlotFinger2]: 12,
		[ItemSlot.ItemSlotTrinket1]: 13,
		[ItemSlot.ItemSlotTrinket2]: 14,
		[ItemSlot.ItemSlotMainHand]: 16,
		[ItemSlot.ItemSlotOffHand]: 17,
		[ItemSlot.ItemSlotRanged]: 18,
	};
}

export class IndividualAddonImporter<SpecType extends Spec> extends Importer {
	private readonly simUI: IndividualSimUI<SpecType>;
	constructor(parent: HTMLElement, simUI: IndividualSimUI<SpecType>) {
		super(parent, simUI, 'Addon Import', true);
		this.simUI = simUI;

		this.descriptionElem.innerHTML = `
			<p>
				Import settings from the <a href="https://www.curseforge.com/wow/addons/wowsimsexporter" target="_blank">WoWSims Importer In-Game Addon</a>.
			</p>
			<p>
				This feature imports gear, race, talents, and professions. It does NOT import buffs, debuffs, consumes, rotation, or custom stats.
			</p>
			<p>
				To import, paste the output from the addon below and click, 'Import'.
			</p>
		`;
	}

	async onImport(data: string) {
		const importJson = JSON.parse(data);

		// Parse all the settings.
		const charClass = nameToClass((importJson['class'] as string) || '');
		if (charClass == Class.ClassUnknown) {
			throw new Error('Could not parse Class!');
		}

		const race = nameToRace((importJson['race'] as string) || '');
		if (race == Race.RaceUnknown) {
			throw new Error('Could not parse Race!');
		}

		const professions = (importJson['professions'] as Array<{ name: string, level: number }>).map(profData => nameToProfession(profData.name));
		professions.forEach((prof, i) => {
			if (prof == Profession.ProfessionUnknown) {
				throw new Error(`Could not parse profession '${importJson['professions'][i]}'`);
			}
		});

		const talentsStr = (importJson['talents'] as string) || '';

		const gearJson = importJson['gear'];
		gearJson.items = (gearJson.items as Array<any>).filter(item => item != null);
		const equipmentSpec = EquipmentSpec.fromJson(gearJson);

		this.finishIndividualImport(this.simUI, charClass, race, equipmentSpec, talentsStr, professions);
	}
}
