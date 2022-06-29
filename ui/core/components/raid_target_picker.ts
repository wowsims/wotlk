import { Component } from '/tbc/core/components/component.js';
import { Input, InputConfig } from '/tbc/core/components/input.js';
import { Player } from '/tbc/core/player.js';
import { Raid } from '/tbc/core/raid.js';
import { EventID, TypedEvent } from '/tbc/core/typed_event.js';
import { RaidTarget } from '/tbc/core/proto/common.js';
import { Spec } from '/tbc/core/proto/common.js';
import { newRaidTarget, emptyRaidTarget } from '/tbc/core/proto_utils/utils.js';

declare var tippy: any;

export interface RaidTargetPickerConfig<ModObject> extends InputConfig<ModObject, RaidTarget> {
	noTargetLabel: string,
	compChangeEmitter: TypedEvent<void>,
}

export interface RaidTargetElemOption {
	iconUrl: string,
	text: string,
	color: string,
	isDropdown: boolean,
};

export interface RaidTargetOption extends RaidTargetElemOption {
	value: Player<any> | null,
};

// Dropdown menu for selecting a player.
export class RaidTargetPicker<ModObject> extends Input<ModObject, RaidTarget> {
	private readonly config: RaidTargetPickerConfig<ModObject>;
	private readonly raid: Raid;
	private readonly noTargetOption: RaidTargetOption;

	private curPlayer: Player<any> | null;
	private curRaidTarget: RaidTarget;

	private currentOptions: Array<RaidTargetOption>;

	private readonly buttonElem: HTMLElement;
	private readonly dropdownElem: HTMLElement;

	constructor(parent: HTMLElement, raid: Raid, modObj: ModObject, config: RaidTargetPickerConfig<ModObject>) {
		super(parent, 'raid-target-picker-root', modObj, config);
		this.rootElem.classList.add('dropdown-root');
		this.config = config;
		this.raid = raid;
		this.curPlayer = this.raid.getPlayerFromRaidTarget(config.getValue(modObj));
		this.curRaidTarget = this.getInputValue();

		this.noTargetOption = {
			iconUrl: '',
			text: config.noTargetLabel,
			color: 'black',
			value: null,
			isDropdown: true,
		};

		this.rootElem.innerHTML = `
			<div class="dropdown-button raid-target-picker-button"></div>
			<div class="dropdown-panel raid-target-picker-dropdown"></div>
    `;

		this.buttonElem = this.rootElem.getElementsByClassName('raid-target-picker-button')[0] as HTMLElement;
		this.dropdownElem = this.rootElem.getElementsByClassName('raid-target-picker-dropdown')[0] as HTMLElement;

		this.buttonElem.addEventListener('click', event => {
			event.preventDefault();
		});

		this.currentOptions = [];
		this.updateOptions(TypedEvent.nextEventID());
		config.compChangeEmitter.on(eventID => {
			this.updateOptions(eventID);
		});

		this.init();
	}

	private makeTargetOptions(): Array<RaidTargetOption> {
		const playerOptions = this.raid.getPlayers().filter(player => player != null).map(player => {
			return {
				iconUrl: player!.getTalentTreeIcon(),
				text: player!.getLabel(),
				color: player!.getClassColor(),
				isDropdown: true,
				value: player,
			};
		});
		return [this.noTargetOption].concat(playerOptions);
	}

	private updateOptions(eventID: EventID) {
		this.currentOptions = this.makeTargetOptions();

		this.dropdownElem.innerHTML = '';
		this.currentOptions.forEach(option => this.dropdownElem.appendChild(this.makeOption(option)));

		const prevRaidTarget = this.curRaidTarget;
		this.curRaidTarget = this.getInputValue();
		if (!RaidTarget.equals(prevRaidTarget, this.curRaidTarget)) {
			this.inputChanged(eventID);
		} else {
			this.setInputValue(this.curRaidTarget);
		}
	}

	private makeOption(data: RaidTargetOption): HTMLElement {
		const option = RaidTargetPicker.makeOptionElem(data);

		option.addEventListener('click', event => {
			event.preventDefault();
			this.curPlayer = data.value;
			this.curRaidTarget = this.getInputValue();
			this.inputChanged(TypedEvent.nextEventID());
		});

		return option;
	}

	getInputElem(): HTMLElement {
		return this.buttonElem;
	}

	getInputValue(): RaidTarget {
		if (this.curPlayer) {
			return this.curPlayer.makeRaidTarget();
		} else {
			return emptyRaidTarget();
		}
	}

	setInputValue(newValue: RaidTarget) {
		this.curRaidTarget = RaidTarget.clone(newValue);
		this.curPlayer = this.raid.getPlayerFromRaidTarget(this.curRaidTarget);

		const optionData = this.currentOptions.find(optionData => optionData.value == this.curPlayer);
		if (!optionData) {
			return;
		}

		this.buttonElem.innerHTML = '';
		this.buttonElem.appendChild(RaidTargetPicker.makeOptionElem(optionData));
	}

	static makeOptionElem(data: RaidTargetElemOption): HTMLElement {
		const optionContainer = document.createElement('div');
		optionContainer.classList.add('dropdown-option-container');

		const option = document.createElement('div');
		option.classList.add('raid-target-picker-option');
		optionContainer.appendChild(option);
		if (data.isDropdown) {
			option.classList.add('dropdown-option');
		}

		if (data.color) {
			option.style.backgroundColor = data.color;
		}

		if (data.iconUrl) {
			const icon = document.createElement('img');
			icon.src = data.iconUrl;
			icon.classList.add('raid-target-picker-icon');
			option.appendChild(icon);
		}

		if (data.text) {
			const label = document.createElement('span');
			label.textContent = data.text;
			label.classList.add('raid-target-picker-label');
			option.appendChild(label);
		}

		return optionContainer;
	}
}
