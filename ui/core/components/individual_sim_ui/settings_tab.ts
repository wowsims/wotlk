import * as Tooltips from '../../constants/tooltips.js';
import { Encounter } from '../../encounter';
import { IndividualSimUI, InputSection } from '../../individual_sim_ui';
import { Consumes, Debuffs, HealingModel, IndividualBuffs, ItemSwap, PartyBuffs, Profession, RaidBuffs, Spec } from '../../proto/common';
import { SavedEncounter, SavedSettings } from '../../proto/ui';
import { professionNames, raceNames } from '../../proto_utils/names';
import { specToEligibleRaces } from '../../proto_utils/utils';
import { EventID, TypedEvent } from '../../typed_event';
import { getEnumValues } from '../../utils';
import { BooleanPicker } from '../boolean_picker';
import { ContentBlock } from '../content_block';
import { EncounterPicker } from '../encounter_picker.js';
import { EnumPicker } from '../enum_picker';
import * as IconInputs from '../icon_inputs.js';
import { Input } from '../input';
import * as BuffDebuffInputs from '../inputs/buffs_debuffs';
import { relevantStatOptions } from '../inputs/stat_options';
import { ItemSwapPicker } from '../item_swap_picker';
import { MultiIconPicker } from '../multi_icon_picker';
import { NumberPicker } from '../number_picker';
import { SavedDataManager } from '../saved_data_manager';
import { SimTab } from '../sim_tab';
import { ConsumesPicker } from './consumes_picker';

export class SettingsTab extends SimTab {
	protected simUI: IndividualSimUI<Spec>;

	readonly leftPanel: HTMLElement;
	readonly rightPanel: HTMLElement;

	readonly column1: HTMLElement = this.buildColumn(1, 'settings-left-col');
	readonly column2: HTMLElement = this.buildColumn(2, 'settings-left-col');
	readonly column3: HTMLElement = this.buildColumn(3, 'settings-left-col');
	readonly column4?: HTMLElement;

	constructor(parentElem: HTMLElement, simUI: IndividualSimUI<Spec>) {
		super(parentElem, simUI, { identifier: 'settings-tab', title: 'Settings' });
		this.simUI = simUI;

		this.leftPanel = document.createElement('div');
		this.leftPanel.classList.add('settings-tab-left', 'tab-panel-left');

		this.leftPanel.appendChild(this.column1);
		this.leftPanel.appendChild(this.column2);
		this.leftPanel.appendChild(this.column3);

		// The 4th column is only used in the raid sim player editor to spread out player settings
		if (this.simUI.isWithinRaidSim) {
			this.column4 = this.buildColumn(4, 'settings-left-col');
			this.leftPanel.appendChild(this.column4);
		}

		this.rightPanel = document.createElement('div');
		this.rightPanel.classList.add('settings-tab-right', 'tab-panel-right', 'within-raid-sim-hide');

		this.contentContainer.appendChild(this.leftPanel);
		this.contentContainer.appendChild(this.rightPanel);

		this.buildTabContent();
	}

	protected buildTabContent() {
		if (!this.simUI.isWithinRaidSim) {
			this.buildEncounterSettings();
		}

		this.buildPlayerSettings();
		this.buildCustomSettingsSections();
		this.buildConsumesSection();
		this.buildOtherSettings();

		if (!this.simUI.isWithinRaidSim) {
			this.buildBuffsSettings();
			this.buildDebuffsSettings();
		}

		if (!this.simUI.isWithinRaidSim) {
			this.buildSavedDataPickers();
		}
	}

	private buildEncounterSettings() {
		const contentBlock = new ContentBlock(this.column1, 'encounter-settings', {
			header: { title: 'Encounter' },
		});

		new EncounterPicker(contentBlock.bodyElement, this.simUI.sim.encounter, this.simUI.individualConfig.encounterPicker, this.simUI);
	}

	private buildPlayerSettings() {
		const column = this.column1;
		const contentBlock = new ContentBlock(column, 'player-settings', {
			header: { title: 'Player' },
		});

		const playerIconGroup = Input.newGroupContainer();
		playerIconGroup.classList.add('player-icon-group', 'icon-group');
		contentBlock.bodyElement.appendChild(playerIconGroup);

		this.configureIconSection(
			playerIconGroup,
			this.simUI.individualConfig.playerIconInputs.map(iconInput => IconInputs.buildIconInput(playerIconGroup, this.simUI.player, iconInput)),
			true,
		);

		const races = specToEligibleRaces[this.simUI.player.spec];
		const _racePicker = new EnumPicker(contentBlock.bodyElement, this.simUI.player, {
			label: 'Race',
			values: races.map(race => {
				return {
					name: raceNames.get(race)!,
					value: race,
				};
			}),
			changedEvent: sim => sim.raceChangeEmitter,
			getValue: sim => sim.getRace(),
			setValue: (eventID, sim, newValue) => sim.setRace(eventID, newValue),
		});

		if (this.simUI.individualConfig.playerInputs?.inputs.length) {
			this.configureInputSection(contentBlock.bodyElement, this.simUI.individualConfig.playerInputs);
		}

		const professionGroup = Input.newGroupContainer();
		contentBlock.bodyElement.appendChild(professionGroup);

		const professions = getEnumValues(Profession) as Array<Profession>;
		const _profession1Picker = new EnumPicker(professionGroup, this.simUI.player, {
			label: 'Profession 1',
			values: professions.map(p => {
				return {
					name: professionNames.get(p)!,
					value: p,
				};
			}),
			changedEvent: sim => sim.professionChangeEmitter,
			getValue: sim => sim.getProfession1(),
			setValue: (eventID, sim, newValue) => sim.setProfession1(eventID, newValue),
		});

		const _profession2Picker = new EnumPicker(professionGroup, this.simUI.player, {
			label: 'Profession 2',
			values: professions.map(p => {
				return {
					name: professionNames.get(p)!,
					value: p,
				};
			}),
			changedEvent: sim => sim.professionChangeEmitter,
			getValue: sim => sim.getProfession2(),
			setValue: (eventID, sim, newValue) => sim.setProfession2(eventID, newValue),
		});
	}

	private buildCustomSettingsSections() {
		(this.simUI.individualConfig.customSections || []).forEach(customSection => {
			const section = customSection(this.column2, this.simUI);
			section.rootElem.classList.add('custom-section');
		});
	}

	private buildConsumesSection() {
		const column = this.simUI.isWithinRaidSim ? this.column3 : this.column2;
		const contentBlock = new ContentBlock(column, 'consumes-settings', {
			header: { title: 'Consumables' },
		});

		new ConsumesPicker(contentBlock.bodyElement, this, this.simUI);
	}

	private buildOtherSettings() {
		const settings = this.simUI.individualConfig.otherInputs?.inputs.filter(inputs => !inputs.extraCssClasses?.includes('within-raid-sim-hide') || true);

		const swapSlots = this.simUI.individualConfig.itemSwapSlots || [];
		if (settings.length > 0 || swapSlots.length > 0) {
			const contentBlock = new ContentBlock(this.column2, 'other-settings', {
				header: { title: 'Other' },
			});

			if (settings.length > 0) {
				this.configureInputSection(contentBlock.bodyElement, this.simUI.individualConfig.otherInputs);
				contentBlock.bodyElement.querySelectorAll('.input-root').forEach(elem => {
					elem.classList.add('input-inline');
				});
			}

			if (swapSlots.length > 0) {
				const _itemSwapPicker = new ItemSwapPicker(contentBlock.bodyElement, this.simUI, this.simUI.player, {
					itemSlots: swapSlots,
				});
			}
		}
	}

	private buildBuffsSettings() {
		const contentBlock = new ContentBlock(this.column3, 'buffs-settings', {
			header: { title: 'Raid Buffs', tooltip: Tooltips.BUFFS_SECTION },
		});

		const buffOptions = relevantStatOptions(BuffDebuffInputs.RAID_BUFFS_CONFIG, this.simUI);
		this.configureIconSection(
			contentBlock.bodyElement,
			buffOptions.map(options => options.picker && new options.picker(contentBlock.bodyElement, this.simUI.player, options.config as any, this.simUI)),
		);

		const miscBuffOptions = relevantStatOptions(BuffDebuffInputs.RAID_BUFFS_MISC_CONFIG, this.simUI);
		if (miscBuffOptions.length > 0) {
			new MultiIconPicker(
				contentBlock.bodyElement,
				this.simUI.player,
				{
					inputs: miscBuffOptions.map(option => option.config),
					label: 'Misc',
				},
				this.simUI,
			);
		}
	}

	private buildDebuffsSettings() {
		const contentBlock = new ContentBlock(this.column3, 'debuffs-settings', {
			header: { title: 'Debuffs', tooltip: Tooltips.DEBUFFS_SECTION },
		});

		const debuffOptions = relevantStatOptions(BuffDebuffInputs.DEBUFFS_CONFIG, this.simUI);
		this.configureIconSection(
			contentBlock.bodyElement,
			debuffOptions.map(options => options.picker && new options.picker(contentBlock.bodyElement, this.simUI.player, options.config as any, this.simUI)),
		);

		const miscDebuffOptions = relevantStatOptions(BuffDebuffInputs.DEBUFFS_MISC_CONFIG, this.simUI);
		if (miscDebuffOptions.length) {
			new MultiIconPicker(
				contentBlock.bodyElement,
				this.simUI.player,
				{
					inputs: miscDebuffOptions.map(options => options.config),
					label: 'Misc',
				},
				this.simUI,
			);
		}
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

		const savedSettingsManager = new SavedDataManager<IndividualSimUI<any>, SavedSettings>(this.rightPanel, this.simUI, {
			label: 'Settings',
			header: { title: 'Saved Settings' },
			storageKey: this.simUI.getSavedSettingsStorageKey(),
			getData: (simUI: IndividualSimUI<any>) => {
				const player = simUI.player;
				return SavedSettings.create({
					raidBuffs: simUI.sim.raid.getBuffs(),
					partyBuffs: player.getParty()?.getBuffs() || PartyBuffs.create(),
					playerBuffs: player.getBuffs(),
					debuffs: simUI.sim.raid.getDebuffs(),
					consumes: player.getConsumes(),
					race: player.getRace(),
					professions: player.getProfessions(),
					enableItemSwap: player.getEnableItemSwap(),
					itemSwap: player.getItemSwapGear().toProto(),
					reactionTimeMs: player.getReactionTime(),
					channelClipDelayMs: player.getChannelClipDelay(),
					inFrontOfTarget: player.getInFrontOfTarget(),
					distanceFromTarget: player.getDistanceFromTarget(),
					healingModel: player.getHealingModel(),
					nibelungAverageCasts: player.getNibelungAverageCasts(),
				});
			},
			setData: (eventID: EventID, simUI: IndividualSimUI<any>, newSettings: SavedSettings) => {
				TypedEvent.freezeAllAndDo(() => {
					simUI.sim.raid.setBuffs(eventID, newSettings.raidBuffs || RaidBuffs.create());
					simUI.sim.raid.setDebuffs(eventID, newSettings.debuffs || Debuffs.create());
					const party = simUI.player.getParty();
					if (party) {
						party.setBuffs(eventID, newSettings.partyBuffs || PartyBuffs.create());
					}
					simUI.player.setBuffs(eventID, newSettings.playerBuffs || IndividualBuffs.create());
					simUI.player.setConsumes(eventID, newSettings.consumes || Consumes.create());
					simUI.player.setRace(eventID, newSettings.race);
					simUI.player.setProfessions(eventID, newSettings.professions);
					simUI.player.setEnableItemSwap(eventID, newSettings.enableItemSwap);
					simUI.player.setItemSwapGear(eventID, simUI.sim.db.lookupItemSwap(newSettings.itemSwap || ItemSwap.create()));
					simUI.player.setReactionTime(eventID, newSettings.reactionTimeMs);
					simUI.player.setChannelClipDelay(eventID, newSettings.channelClipDelayMs);
					simUI.player.setInFrontOfTarget(eventID, newSettings.inFrontOfTarget);
					simUI.player.setDistanceFromTarget(eventID, newSettings.distanceFromTarget);
					simUI.player.setHealingModel(eventID, newSettings.healingModel || HealingModel.create());
					simUI.player.setNibelungAverageCasts(eventID, newSettings.nibelungAverageCasts);
				});
			},
			changeEmitters: [
				this.simUI.sim.raid.buffsChangeEmitter,
				this.simUI.sim.raid.debuffsChangeEmitter,
				this.simUI.player.getParty()!.buffsChangeEmitter,
				this.simUI.player.buffsChangeEmitter,
				this.simUI.player.consumesChangeEmitter,
				this.simUI.player.raceChangeEmitter,
				this.simUI.player.professionChangeEmitter,
				this.simUI.player.itemSwapChangeEmitter,
				this.simUI.player.miscOptionsChangeEmitter,
				this.simUI.player.inFrontOfTargetChangeEmitter,
				this.simUI.player.distanceFromTargetChangeEmitter,
				this.simUI.player.healingModelChangeEmitter,
			],
			equals: (a: SavedSettings, b: SavedSettings) => SavedSettings.equals(a, b),
			toJson: (a: SavedSettings) => SavedSettings.toJson(a),
			fromJson: (obj: any) => SavedSettings.fromJson(obj),
		});

		this.simUI.sim.waitForInit().then(() => {
			savedEncounterManager.loadUserData();
			savedSettingsManager.loadUserData();
		});
	}

	private configureInputSection(sectionElem: HTMLElement, sectionConfig: InputSection) {
		sectionConfig.inputs.forEach(inputConfig => {
			if (inputConfig.type == 'number') {
				new NumberPicker(sectionElem, this.simUI.player, inputConfig);
			} else if (inputConfig.type == 'boolean') {
				new BooleanPicker(sectionElem, this.simUI.player, { ...inputConfig, reverse: true });
			} else if (inputConfig.type == 'enum') {
				new EnumPicker(sectionElem, this.simUI.player, inputConfig);
			}
		});
	}

	private configureIconSection(sectionElem: HTMLElement, iconPickers: Array<any>, adjustColumns?: boolean) {
		if (iconPickers.length == 0) {
			sectionElem.classList.add('hide');
		} else if (adjustColumns) {
			if (iconPickers.length <= 4) {
				sectionElem.style.gridTemplateColumns = `repeat(${iconPickers.length}, 1fr)`;
			} else if (iconPickers.length > 4 && iconPickers.length < 8) {
				sectionElem.style.gridTemplateColumns = `repeat(${Math.ceil(iconPickers.length / 2)}, 1fr)`;
			}
		}
	}
}
