import { ContentBlock } from "../core/components/content_block";
import { EncounterPicker } from "../core/components/encounter_picker";
import { IconPicker } from "../core/components/icon_picker";
import { SavedDataManager } from "../core/components/saved_data_manager";
import { SimTab } from "../core/components/sim_tab";

import { Encounter } from "../core/encounter";
import { Raid } from "../core/raid";
import { EventID } from "../core/typed_event";

import { RaidBuffs } from "../core/proto/common";
import { SavedEncounter } from "../core/proto/ui";
import { ActionId } from "../core/proto_utils/action_id";

import { AssignmentsPicker } from "./assignments_picker";
import { BlessingsPicker } from "./blessings_picker";
import { RaidSimUI } from "./raid_sim_ui";
import { TanksPicker } from "./tanks_picker";

import * as Tooltips from "../core/constants/tooltips.js";

export class SettingsTab extends SimTab {
	protected simUI: RaidSimUI;

	readonly leftPanel: HTMLElement;
	readonly rightPanel: HTMLElement;

	readonly column1: HTMLElement = this.buildColumn(1, 'raid-settings-col');
	readonly column2: HTMLElement = this.buildColumn(2, 'raid-settings-col');
	readonly column3: HTMLElement = this.buildColumn(3, 'raid-settings-col');

	constructor(parentElem: HTMLElement, simUI: RaidSimUI) {
		super(parentElem, simUI, { identifier: 'raid-settings-tab', title: 'Settings' });
		this.simUI = simUI;

		this.leftPanel = document.createElement('div');
		this.leftPanel.classList.add('raid-settings-tab-left', 'tab-panel-left');

		this.leftPanel.appendChild(this.column1);
		this.leftPanel.appendChild(this.column2);
		this.leftPanel.appendChild(this.column3);

		this.rightPanel = document.createElement('div');
		this.rightPanel.classList.add('raid-settings-tab-right', 'tab-panel-right', 'within-raid-sim-hide');

		this.contentContainer.appendChild(this.leftPanel);
		this.contentContainer.appendChild(this.rightPanel);

		this.buildTabContent();
	}

	protected buildTabContent() {
		this.buildEncounterSettings();
		this.buildConsumesSettings();

		this.buildTankSettings();
		this.buildAssignmentSettings();
		this.buildOtherSettings();

		this.buildBlessingsPicker();
		this.buildSavedDataPickers();
	}

	private buildEncounterSettings() {
		const contentBlock = new ContentBlock(this.column1, 'encounter-settings', {
			header: { title: 'Encounter' }
		});

		new EncounterPicker(contentBlock.bodyElement, this.simUI.sim.encounter, { showExecuteProportion: true }, this.simUI);
	}

	private buildConsumesSettings() {
		const contentBlock = new ContentBlock(this.column1, 'consumes-settings', {
			header: { title: 'Consumables' }
		});

		let container = document.createElement('div');
		container.classList.add('consumes-container');

		contentBlock.bodyElement.appendChild(container);

		this.makeBooleanRaidIconBuffInput(container, ActionId.fromItemId(37094), 'scrollOfStamina'),
			this.makeBooleanRaidIconBuffInput(container, ActionId.fromItemId(43466), 'scrollOfStrength'),
			this.makeBooleanRaidIconBuffInput(container, ActionId.fromItemId(43464), 'scrollOfAgility'),
			this.makeBooleanRaidIconBuffInput(container, ActionId.fromItemId(37092), 'scrollOfIntellect'),
			this.makeBooleanRaidIconBuffInput(container, ActionId.fromItemId(37098), 'scrollOfSpirit'),
			this.makeBooleanRaidIconBuffInput(container, ActionId.fromItemId(43468), 'scrollOfProtection'),
			this.makeBooleanRaidIconBuffInput(container, ActionId.fromItemId(49633), 'drumsOfForgottenKings'),
			this.makeBooleanRaidIconBuffInput(container, ActionId.fromItemId(49634), 'drumsOfTheWild');
	}

	private buildOtherSettings() {
		const contentBlock = new ContentBlock(this.column2, 'other-settings', {
		  header: {title: 'Other'}
		});

		this.makeBooleanRaidIconBuffInput(contentBlock.bodyElement, ActionId.fromSpellId(73828), 'strengthOfWrynn');

		// new BooleanPicker(contentBlock.bodyElement, this.simUI.sim.raid, {
		// 	label: 'Stagger Stormstrikes',
		// 	labelTooltip: 'When there are multiple Enhancement Shaman in the raid, causes them to coordinate their Stormstrike casts for optimal SS charge usage.',
		// 	changedEvent: (raid: Raid) => raid.staggerStormstrikesChangeEmitter,
		// 	getValue: (raid: Raid) => raid.getStaggerStormstrikes(),
		// 	setValue: (eventID: EventID, raid: Raid, newValue: boolean) => {
		// 		raid.setStaggerStormstrikes(eventID, newValue);
		// 	},
		// });
	}

	private buildTankSettings() {
		const contentBlock = new ContentBlock(this.column2, 'tanks-settings', {
			header: { title: 'Tanks' }
		});

		new TanksPicker(contentBlock.bodyElement, this.simUI);
	}

	private buildAssignmentSettings() {
		const contentBlock = new ContentBlock(this.column2, 'assignments-settings', {
			header: { title: 'External Buffs' }
		});

		new AssignmentsPicker(contentBlock.bodyElement, this.simUI);
	}

	private buildBlessingsPicker() {
		const contentBlock = new ContentBlock(this.column3, 'blessings-settings', {
			header: { title: 'Blessings', tooltip: Tooltips.BLESSINGS_SECTION }
		});

		this.simUI.blessingsPicker = new BlessingsPicker(contentBlock.bodyElement, this.simUI);
		this.simUI.blessingsPicker.changeEmitter.on(eventID => this.simUI.changeEmitter.emit(eventID));
	}

	private buildSavedDataPickers() {
		const savedEncounterManager = new SavedDataManager<Encounter, SavedEncounter>(this.rightPanel, this.simUI.sim.encounter, {
			label: 'Encounter',
			header: { title: 'Saved Encounters' },
			storageKey: this.simUI.getSavedEncounterStorageKey(),
			getData: (encounter: Encounter) => SavedEncounter.create({ encounter: encounter.toProto() }),
			setData: (eventID: EventID, encounter: Encounter, newEncounter: SavedEncounter) => encounter.fromProto(eventID, newEncounter.encounter!),
			changeEmitters: [this.simUI.sim.encounter.changeEmitter],
			equals: (a: SavedEncounter, b: SavedEncounter) => SavedEncounter.equals(a, b),
			toJson: (a: SavedEncounter) => SavedEncounter.toJson(a),
			fromJson: (obj: any) => SavedEncounter.fromJson(obj),
		});

		this.simUI.sim.waitForInit().then(() => {
			savedEncounterManager.loadUserData();
		});
	}

	private makeBooleanRaidIconBuffInput(parent: HTMLElement, actionId: ActionId, field: keyof RaidBuffs): IconPicker<Raid, boolean> {
		const raid = this.simUI.sim.raid;

		return new IconPicker<Raid, boolean>(parent, raid, {
			actionId: actionId,
			states: 2,
			changedEvent: (raid: Raid) => raid.buffsChangeEmitter,
			getValue: (raid: Raid) => raid.getBuffs()[field] as unknown as boolean,
			setValue: (eventID: EventID, raid: Raid, newValue: boolean) => {
				const newBuffs = raid.getBuffs();
				(newBuffs[field] as unknown as boolean) = newValue;
				raid.setBuffs(eventID, newBuffs);
			},
		});
	}
}
