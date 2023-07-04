import { IndividualSimUI, InputSection } from "../../individual_sim_ui";
import {
	Spec,
} from "../../proto/common";
import { EventID, TypedEvent } from "../../typed_event";
import { Player } from "../../player";

import { ContentBlock } from "../content_block";
import { SimTab } from "../sim_tab";
import { NumberPicker } from "../number_picker";
import { BooleanPicker } from "../boolean_picker";
import { EnumPicker } from "../enum_picker";
import { Input } from "../input";
import { ItemSwapPicker } from "../item_swap_picker";
import { CooldownsPicker } from "./cooldowns_picker";
import { CustomRotationPicker } from "./custom_rotation_picker";

import * as IconInputs from '../icon_inputs.js';
import * as Tooltips from '../../constants/tooltips.js';

import { APLRotationPicker } from "./apl_rotation_picker";

export class RotationTab extends SimTab {
  protected simUI: IndividualSimUI<Spec>;

  readonly leftPanel: HTMLElement;
  readonly rightPanel: HTMLElement;

  constructor(parentElem: HTMLElement, simUI: IndividualSimUI<Spec>) {
    super(parentElem, simUI, {identifier: 'rotation-tab', title: 'Rotation'});
    this.rootElem.classList.add('experimental');
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
	this.buildContent();
	this.buildRotationSettings();
	this.buildCooldownSettings();
  }

  private updateSections() {
	if (this.simUI.player.aplRotation.enabled) {
		this.rootElem.classList.add('rotation-type-apl');
		this.rootElem.classList.remove('rotation-type-legacy');
	} else {
		this.rootElem.classList.remove('rotation-type-apl');
		this.rootElem.classList.add('rotation-type-legacy');
	}
  }

  private buildHeader() {
	const header = document.createElement('div');
	header.classList.add('rotation-tab-header');
	this.leftPanel.appendChild(header);

	new EnumPicker(header, this.simUI.player, {
		label: 'Rotation Type',
		labelTooltip: 'Whether to use the legacy rotation options, or the new APL rotation options.',
		inline: true,
		values: [
			{value: 0, name: 'Legacy'},
			{value: 1, name: 'APL'},
		],
		changedEvent: (player: Player<any>) => player.rotationChangeEmitter,
		getValue: (player: Player<any>) => Number(player.aplRotation.enabled),
		setValue: (eventID: EventID, player: Player<any>, newValue: number) => {
			player.aplRotation.enabled = !!newValue;
			player.rotationChangeEmitter.emit(eventID);
		},
	});
  }

  private buildContent() {
	const content = document.createElement('div');
	content.classList.add('rotation-tab-main');
	this.leftPanel.appendChild(content);

	new APLRotationPicker(content, this.simUI, this.simUI.player);
  }

	private buildRotationSettings() {
		const contentBlock = new ContentBlock(this.leftPanel, 'rotation-settings', {
			header: {title: 'Rotation'}
		});

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
	}

	private buildCooldownSettings() {
		const contentBlock = new ContentBlock(this.leftPanel, 'cooldown-settings', {
			header: {title: 'Cooldowns', tooltip: Tooltips.COOLDOWNS_SECTION}
		});

		new CooldownsPicker(contentBlock.bodyElement, this.simUI.player);
	}

	private configureInputSection(sectionElem: HTMLElement, sectionConfig: InputSection) {
		sectionConfig.inputs.forEach(inputConfig => {
			if (inputConfig.type == 'number') {
				new NumberPicker(sectionElem, this.simUI.player, inputConfig);
			} else if (inputConfig.type == 'boolean') {
				new BooleanPicker(sectionElem, this.simUI.player, {...inputConfig, ...{cssScheme: this.simUI.cssScheme}});
			} else if (inputConfig.type == 'enum') {
				new EnumPicker(sectionElem, this.simUI.player, inputConfig);
			} else if (inputConfig.type == 'customRotation') {
				new CustomRotationPicker(sectionElem, this.simUI, this.simUI.player, inputConfig);
			} else if (inputConfig.type == 'itemSwap'){
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
}
