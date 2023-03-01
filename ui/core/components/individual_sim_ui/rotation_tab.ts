import { IndividualSimUI } from "../../individual_sim_ui";
import {
	Consumes,
	Cooldowns,
	Debuffs,
	IndividualBuffs,
	PartyBuffs,
	Profession,
	RaidBuffs,
	Spec,
	Stat
} from "../../proto/common";
import { ActionId } from "../../proto_utils/action_id";
import { EventID, TypedEvent } from "../../typed_event";
import { getEnumValues } from "../../utils";
import { Player } from "../../player";

import { SimTab } from "../sim_tab";
import { NumberPicker } from "../number_picker";
import { BooleanPicker } from "../boolean_picker";
import { EnumPicker } from "../enum_picker";
import { Input } from "../input";
import { MultiIconPicker } from "../multi_icon_picker";
import { IconPickerConfig } from "../icon_picker";
import { TypedIconPickerConfig } from "../input_helpers";

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
  }

  protected buildTabContent() {
	this.buildHeader();
	this.buildContent();
  }

  private buildHeader() {
	const header = document.createElement('div');
	header.classList.add('rotation-tab-header');
	this.leftPanel.appendChild(header);

	new BooleanPicker(header, this.simUI.player, {
		label: 'Use APL Rotation',
		labelTooltip: 'Enables the APL Rotation options.',
		inline: true,
		changedEvent: (player: Player<any>) => player.rotationChangeEmitter,
		getValue: (player: Player<any>) => player.getAplRotation().enabled,
		setValue: (eventID: EventID, player: Player<any>, newValue: boolean) => {
			const rotation = player.getAplRotation();
			rotation.enabled = newValue;
			player.setAplRotation(eventID, rotation);
		},
	});
  }

  private buildContent() {
	const content = document.createElement('div');
	content.classList.add('rotation-tab-main');
	this.leftPanel.appendChild(content);

	new APLRotationPicker(content, this.simUI, this.simUI.player);
  }
}
