import { IndividualSimUI, InputSection } from "../../individual_sim_ui";
import {
	Cooldowns,
	Spec,
} from "../../proto/common";
import {
	APLRotation,
	APLRotation_Type as APLRotationType,
} from "../../proto/apl";
import {
	SavedRotation,
} from "../../proto/ui";
import { EventID, TypedEvent } from "../../typed_event";
import { Player } from "../../player";
import { aplLaunchStatuses, LaunchStatus } from '../../launched_sims';

import { ContentBlock } from "../content_block";
import { SimTab } from "../sim_tab";
import { NumberPicker } from "../number_picker";
import { BooleanPicker } from "../boolean_picker";
import { EnumPicker } from "../enum_picker";
import { Input } from "../input";
import { ItemSwapPicker } from "../item_swap_picker";
import { CooldownsPicker } from "./cooldowns_picker";
import { CustomRotationPicker } from "./custom_rotation_picker";
import { SavedDataManager } from "../saved_data_manager";

import * as IconInputs from '../icon_inputs.js';
import * as Tooltips from '../../constants/tooltips.js';

import { APLRotationPicker } from "./apl_rotation_picker";

export class RotationTab extends SimTab {
	protected simUI: IndividualSimUI<Spec>;

	readonly leftPanel: HTMLElement;
	readonly rightPanel: HTMLElement;

	constructor(parentElem: HTMLElement, simUI: IndividualSimUI<Spec>) {
		super(parentElem, simUI, { identifier: 'rotation-tab', title: 'Rotation' });
		this.simUI = simUI;

		this.leftPanel = document.createElement('div');
		this.leftPanel.classList.add('rotation-tab-left', 'tab-panel-left');

		this.rightPanel = document.createElement('div');
		this.rightPanel.classList.add('rotation-tab-right', 'tab-panel-right');

		this.contentContainer.appendChild(this.leftPanel);
		this.contentContainer.appendChild(this.rightPanel);

		this.buildTabContent();

		this.updateSections();
		this.simUI.player.rotationChangeEmitter.on(() => this.updateSections());
	}

	protected buildTabContent() {
		this.buildHeader();

		this.buildAutoContent();
		this.buildAplContent();
		this.buildSimpleOrLegacyContent(this.simUI.player.hasSimpleRotationGenerator());

		this.buildSavedDataPickers();
	}

	private updateSections() {
		this.rootElem.classList.remove('rotation-type-auto');
		this.rootElem.classList.remove('rotation-type-simple');
		this.rootElem.classList.remove('rotation-type-apl');
		this.rootElem.classList.remove('rotation-type-legacy');

		const rotType = this.simUI.player.getRotationType();
		if (rotType == APLRotationType.TypeAuto) {
			this.rootElem.classList.add('rotation-type-auto');
		} else if (rotType == APLRotationType.TypeSimple) {
			this.rootElem.classList.add('rotation-type-simple');
		} else if (rotType == APLRotationType.TypeAPL) {
			this.rootElem.classList.add('rotation-type-apl');
		} else if (rotType == APLRotationType.TypeLegacy) {
			this.rootElem.classList.add('rotation-type-legacy');
		}
	}

	private buildHeader() {
		const header = document.createElement('div');
		header.classList.add('rotation-tab-header');
		this.leftPanel.appendChild(header);

		const aplLaunchStatus = aplLaunchStatuses[this.simUI.player.spec];

		new EnumPicker(header, this.simUI.player, {
			label: 'Rotation Type Test',
			labelTooltip: 'Which set of options to use for specifying the rotation.',
			inline: true,
			values: aplLaunchStatus == LaunchStatus.Alpha ? [
				{ value: APLRotationType.TypeLegacy, name: 'Legacy' },
				{ value: APLRotationType.TypeAPL, name: 'APL' },
			] : aplLaunchStatus == LaunchStatus.Beta ? [
				{ value: APLRotationType.TypeAuto, name: 'Auto' },
				{ value: APLRotationType.TypeAPL, name: 'APL' },
				{ value: APLRotationType.TypeLegacy, name: 'Legacy' },
			] : this.simUI.player.hasSimpleRotationGenerator() ? [
				{ value: APLRotationType.TypeAuto, name: 'Auto' },
				{ value: APLRotationType.TypeSimple, name: 'Simple' },
				{ value: APLRotationType.TypeAPL, name: 'APL' },
			] : [
				{ value: APLRotationType.TypeAuto, name: 'Auto' },
				{ value: APLRotationType.TypeAPL, name: 'APL' },
			],
			changedEvent: (player: Player<any>) => player.rotationChangeEmitter,
			getValue: (player: Player<any>) => player.getRotationType(),
			setValue: (eventID: EventID, player: Player<any>, newValue: number) => {
				player.aplRotation.enabled = false;
				player.aplRotation.type = newValue;
				player.rotationChangeEmitter.emit(eventID);
			},
		});
	}

	private buildAutoContent() {
		const content = document.createElement('div');
		content.classList.add('rotation-tab-auto');
		this.leftPanel.appendChild(content);
	}

	private buildAplContent() {
		const content = document.createElement('div');
		content.classList.add('rotation-tab-apl');
		this.leftPanel.appendChild(content);

		new APLRotationPicker(content, this.simUI, this.simUI.player);
	}

	private buildSimpleOrLegacyContent(isSimple: boolean) {
		if (!isSimple && aplLaunchStatuses[this.simUI.player.spec] == LaunchStatus.Launched) {
			return;
		}
		const cssClass = isSimple ? 'rotation-tab-simple' : 'rotation-tab-legacy';

		const contentBlock = new ContentBlock(this.leftPanel, 'rotation-settings', {
			header: { title: 'Rotation' }
		});
		contentBlock.rootElem.classList.add(cssClass);

		const rotationIconGroup = Input.newGroupContainer();
		rotationIconGroup.classList.add('rotation-icon-group', 'icon-group');
		contentBlock.bodyElement.appendChild(rotationIconGroup);

		if (this.simUI.individualConfig.rotationIconInputs?.length) {
			this.configureIconSection(
				rotationIconGroup,
				this.simUI.individualConfig.rotationIconInputs.map(iconInput => IconInputs.buildIconInput(rotationIconGroup, this.simUI.player, iconInput)),
				true
			);
		}

		this.configureInputSection(contentBlock.bodyElement, this.simUI.individualConfig.rotationInputs);

		contentBlock.bodyElement.querySelectorAll('.input-root').forEach(elem => {
			elem.classList.add('input-inline');
		})

		const cooldownsContentBlock = new ContentBlock(this.leftPanel, 'cooldown-settings', {
			header: { title: 'Cooldowns', tooltip: Tooltips.COOLDOWNS_SECTION }
		});
		cooldownsContentBlock.rootElem.classList.add(cssClass);

		new CooldownsPicker(cooldownsContentBlock.bodyElement, this.simUI.player);
	}

	private configureInputSection(sectionElem: HTMLElement, sectionConfig: InputSection) {
		sectionConfig.inputs.forEach(inputConfig => {
			if (inputConfig.type == 'number') {
				new NumberPicker(sectionElem, this.simUI.player, inputConfig);
			} else if (inputConfig.type == 'boolean') {
				new BooleanPicker(sectionElem, this.simUI.player, { ...inputConfig, ...{ cssScheme: this.simUI.cssScheme } });
			} else if (inputConfig.type == 'enum') {
				new EnumPicker(sectionElem, this.simUI.player, inputConfig);
			} else if (inputConfig.type == 'customRotation') {
				new CustomRotationPicker(sectionElem, this.simUI, this.simUI.player, inputConfig);
			} else if (inputConfig.type == 'itemSwap') {
				new ItemSwapPicker(sectionElem, this.simUI, this.simUI.player, inputConfig)
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

	private buildSavedDataPickers() {
		const aplIsLaunched = aplLaunchStatuses[this.simUI.player.spec] == LaunchStatus.Launched;

		const savedRotationsManager = new SavedDataManager<Player<any>, SavedRotation>(this.rightPanel, this.simUI.player, {
			label: 'Rotation',
			header: { title: 'Saved Rotations' },
			storageKey: this.simUI.getSavedRotationStorageKey(),
			getData: (player: Player<any>) => SavedRotation.create({
				rotation: APLRotation.clone(player.aplRotation),
				specRotationOptionsJson: aplIsLaunched ? '{}' : JSON.stringify(player.specTypeFunctions.rotationToJson(player.getRotation())),
				cooldowns: aplIsLaunched ? Cooldowns.create() : player.getCooldowns(),
			}),
			setData: (eventID: EventID, player: Player<any>, newRotation: SavedRotation) => {
				TypedEvent.freezeAllAndDo(() => {
					player.setAplRotation(eventID, newRotation.rotation || APLRotation.create());
					if (!aplIsLaunched) {
						if (newRotation.specRotationOptionsJson) {
							try {
								const json = JSON.parse(newRotation.specRotationOptionsJson);
								const specRot = player.specTypeFunctions.rotationFromJson(json);
								player.setRotation(eventID, specRot);
							} catch (e) {
								console.warn('Error parsing rotation spec options: ' + e);
							}
						}
						player.setCooldowns(eventID, newRotation.cooldowns || Cooldowns.create());
					}
				});
			},
			changeEmitters: [this.simUI.player.rotationChangeEmitter, this.simUI.player.cooldownsChangeEmitter, this.simUI.player.talentsChangeEmitter],
			equals: (a: SavedRotation, b: SavedRotation) => {
				// Uncomment this to debug equivalence checks with preset rotations (e.g. the chip doesn't highlight)
				//console.log(`Rot A: ${SavedRotation.toJsonString(a, {prettySpaces: 2})}\n\nRot B: ${SavedRotation.toJsonString(b, {prettySpaces: 2})}`);
				return SavedRotation.equals(a, b);
			},
			toJson: (a: SavedRotation) => SavedRotation.toJson(a),
			fromJson: (obj: any) => SavedRotation.fromJson(obj),
		});

		this.simUI.sim.waitForInit().then(() => {
			savedRotationsManager.loadUserData();
			(this.simUI.individualConfig.presets.rotations || []).forEach(presetRotation => {
				const rotData = presetRotation.rotation;
				// Fill default values so the equality checks always work.
				if (!rotData.cooldowns) rotData.cooldowns = Cooldowns.create();
				if (!rotData.specRotationOptionsJson) rotData.specRotationOptionsJson = '{}';
				if (!rotData.rotation) rotData.rotation = APLRotation.create();

				savedRotationsManager.addSavedData({
					name: presetRotation.name,
					tooltip: presetRotation.tooltip,
					isPreset: true,
					data: rotData,
					enableWhen: presetRotation.enableWhen,
				});
			});
		});
	}
}
