import { Class } from '../proto/common.js';
import { EquipmentSpec } from '../proto/common.js';
import { ItemSpec } from '../proto/common.js';
import { Glyphs } from '../proto/common.js';
import { Profession } from '../proto/common.js';
import { Race } from '../proto/common.js';
import { Spec } from '../proto/common.js';
import { IndividualSimSettings } from '../proto/ui.js';
import { IndividualSimUI } from '../individual_sim_ui.js';
import { Player } from '../player.js';
import { classNames, nameToClass, nameToRace, nameToProfession } from '../proto_utils/names.js';
import { classGlyphsConfig, talentSpellIdsToTalentString } from '../talents/factory.js';
import { GlyphConfig } from '../talents/glyphs_picker.js';
import { EventID, TypedEvent } from '../typed_event.js';

import { Popup } from './popup.js';

declare var $: any;
declare var tippy: any;

export function newIndividualImporters<SpecType extends Spec>(simUI: IndividualSimUI<SpecType>): HTMLElement {
	const importFragment = document.createElement('fragment');
	importFragment.innerHTML = `
		<div class="dropdown sim-dropdown-menu">
			<a href="javascript:void(0)" class="import-link" role="button" data-bs-toggle="dropdown" data-bs-offset="0,0" aria-expanded="false" >
				<i class="fas fa-file-import"></i>
				Import
			</a>
			<ul class="dropdown-menu"></ul>
		</div>
	`;

	const menuElem = importFragment.getElementsByClassName('dropdown-menu')[0] as HTMLElement;
	const addMenuItem = (label: string, onClick: () => void, showInRaidSim: boolean) => {
		const itemFragment = document.createElement('fragment');
		itemFragment.innerHTML = `
			<li class="${showInRaidSim ? '' : 'within-raid-sim-hide'}">
				<a href="javascript:void(0)" class="dropdown-item" role="button">${label}</a>
			</li>
		`;
		const itemElem = itemFragment.children[0] as HTMLElement;
		const linkElem = itemElem.children[0] as HTMLElement;
		linkElem.addEventListener('click', onClick);
		menuElem.appendChild(itemElem);
	};

	addMenuItem('Json', () => new IndividualJsonImporter(menuElem, simUI), true);
	addMenuItem('80U', () => new Individual80UImporter(menuElem, simUI), true);
	addMenuItem('Addon', () => new IndividualAddonImporter(menuElem, simUI), true);

	return importFragment.children[0] as HTMLElement;
}

export abstract class Importer extends Popup {
	private readonly textElem: HTMLTextAreaElement;
	protected readonly descriptionElem: HTMLElement;
	protected readonly importButton: HTMLButtonElement;
	private readonly includeFile: boolean;

	constructor(parent: HTMLElement, title: string, includeFile: boolean) {
		super(parent);
		this.includeFile = includeFile;
		const uploadInputId = 'upload-input-' + title.toLowerCase().replaceAll(' ', '-');

		this.rootElem.classList.add('importer');
		let htmlVal = `
			<span class="importer-title">${title}</span>
			<div class="import-description">
			</div>
			<div class="import-content">
				<textarea class="importer-textarea"></textarea>
			</div>
			<div class="actions-row">
		`;
		if (this.includeFile) {
			htmlVal += `<label for="${uploadInputId}" class="importer-button sim-button upload-button">UPLOAD FROM FILE</label>
				<input type="file" id="${uploadInputId}" class="importer-upload-input" hidden>
			`
		}
		htmlVal += `<button class="importer-button sim-button import-button">IMPORT</button>
			</div>
		`;

		this.rootElem.innerHTML = htmlVal

		this.addCloseButton();

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
		this.importButton.addEventListener('click', event => {
			try {
				this.onImport(this.textElem.value || '');
			} catch (error) {
				alert('Import error: ' + error);
			}
		});
	}

	abstract onImport(data: string): void

	protected finishIndividualImport<SpecType extends Spec>(simUI: IndividualSimUI<SpecType>, charClass: Class, race: Race, equipmentSpec: EquipmentSpec, talentsStr: string, glyphs: Glyphs | null, professions: Array<Profession>) {
		const playerClass = simUI.player.getClass();
		if (charClass != playerClass) {
			throw new Error(`Wrong Class! Expected ${classNames[playerClass]} but found ${classNames[charClass]}!`);
		}

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
			if (talentsStr) {
				simUI.player.setTalentsString(eventID, talentsStr);
			}
			if (glyphs) {
				simUI.player.setGlyphs(eventID, glyphs)
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

class IndividualJsonImporter<SpecType extends Spec> extends Importer {
	private readonly simUI: IndividualSimUI<SpecType>;
	constructor(parent: HTMLElement, simUI: IndividualSimUI<SpecType>) {
		super(parent, 'JSON Import', true);
		this.simUI = simUI;

		this.descriptionElem.innerHTML = `
			<p>
				Import settings from a JSON text file, which can be created using the JSON Export feature of this site.
			</p>
			<p>
				To import, paste the JSON text below and click, 'Import'.
			</p>
		`;
	}

	onImport(data: string) {
		const proto = IndividualSimSettings.fromJsonString(data);
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

class Individual80UImporter<SpecType extends Spec> extends Importer {
	private readonly simUI: IndividualSimUI<SpecType>;
	constructor(parent: HTMLElement, simUI: IndividualSimUI<SpecType>) {
		super(parent, '80 Upgrades Import', true);
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
			if (itemJson.gems) {
				itemSpec.gems = (itemJson.gems as Array<any>).filter(gemJson => gemJson?.id).map(gemJson => gemJson.id);
			}
			equipmentSpec.items.push(itemSpec);
		});

		const gear = this.simUI.sim.db.lookupEquipmentSpec(equipmentSpec);

		this.finishIndividualImport(this.simUI, charClass, race, equipmentSpec, talentsStr, null, []);
	}
}

class IndividualAddonImporter<SpecType extends Spec> extends Importer {
	private readonly simUI: IndividualSimUI<SpecType>;
	constructor(parent: HTMLElement, simUI: IndividualSimUI<SpecType>) {
		super(parent, 'Addon Import', true);
		this.simUI = simUI;

		this.descriptionElem.innerHTML = `
			<p>
				Import settings from the <a href="https://www.curseforge.com/wow/addons/wowsimsexporter" target="_blank">WoWSims Importer In-Game Addon</a>.
			</p>
			<p>
				This feature imports gear, race, talents, glyphs, and professions. It does NOT import buffs, debuffs, consumes, rotation, or custom stats.
			</p>
			<p>
				To import, paste the output from the addon below and click, 'Import'.
			</p>
		`;
	}

	onImport(data: string) {
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

		const professions = (importJson['professions'] as Array<{name: string, level: number}>).map(profData => nameToProfession(profData.name));
		professions.forEach((prof, i) => {
			if (prof == Profession.ProfessionUnknown) {
				throw new Error(`Could not parse profession '${importJson['professions'][i]}'`);
			}
		});

		const talentsStr = (importJson['talents'] as string) || '';

		const glyphsConfig = classGlyphsConfig[charClass];
		const majorGlyphIDs = (importJson['glyphs']['major'] as Array<string>).map(glyphName => glyphNameToID(glyphName, glyphsConfig.majorGlyphs));
		const minorGlyphIDs = (importJson['glyphs']['minor'] as Array<string>).map(glyphName => glyphNameToID(glyphName, glyphsConfig.minorGlyphs));
		const glyphs = Glyphs.create({
			major1: majorGlyphIDs[0] || 0,
			major2: majorGlyphIDs[1] || 0,
			major3: majorGlyphIDs[2] || 0,
			minor1: minorGlyphIDs[0] || 0,
			minor2: minorGlyphIDs[1] || 0,
			minor3: minorGlyphIDs[2] || 0,
		});

		const gearJson = importJson['gear'];
		gearJson.items = (gearJson.items as Array<any>).filter(item => item != null);
		(gearJson.items as Array<any>).forEach(item => {
			if (item.gems) {
				item.gems = (item.gems as Array<any>).map(gem => gem || 0);
			}
		});
		const equipmentSpec = EquipmentSpec.fromJson(gearJson);

		this.finishIndividualImport(this.simUI, charClass, race, equipmentSpec, talentsStr, glyphs, professions);
	}
}

function glyphNameToID(glyphName: string, glyphsConfig: Record<number, GlyphConfig>): number {
	if (!glyphName) {
		return 0;
	}

	for (let glyphIDStr in glyphsConfig) {
		if (glyphsConfig[glyphIDStr].name == glyphName) {
			return parseInt(glyphIDStr);
		}
	}
	throw new Error(`Unknown glyph name '${glyphName}'`);
}
