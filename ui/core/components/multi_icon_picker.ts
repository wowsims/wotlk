import { ActionId } from '/wotlk/core/proto_utils/action_id.js';
import { EventID, TypedEvent } from '/wotlk/core/typed_event.js';
import { SimUI } from '/wotlk/core/sim_ui.js';

import { Component } from './component.js';
import { IconPicker, IconPickerConfig } from './icon_picker.js';
import { Input, InputConfig } from './input.js';

export interface MultiIconPickerItemConfig<ModObject> extends IconPickerConfig<ModObject, any> {
}

export interface MultiIconPickerConfig<ModObject> {
	inputs: Array<MultiIconPickerItemConfig<ModObject>>,
	emptyColor: string,
	numColumns: number,
}

// Icon-based UI for a dropdown with multiple icon pickers.
// ModObject is the object being modified (Sim, Player, or Target).
export class MultiIconPicker<ModObject> extends Component {
	private readonly config: MultiIconPickerConfig<ModObject>;

	private currentValue: ActionId | null;

	private readonly buttonElem: HTMLAnchorElement;
	private readonly pickers: Array<IconPicker<ModObject, any>>;

	constructor(parent: HTMLElement, modObj: ModObject, config: MultiIconPickerConfig<ModObject>, simUI: SimUI) {
		super(parent, 'multi-icon-picker-root');
		this.config = config;
		this.currentValue = null;
		this.rootElem.classList.add('dropdown-root');

		this.rootElem.innerHTML = `
			<a class="dropdown-button multi-icon-picker-button"></a>
			<div class="dropdown-panel multi-icon-picker-dropdown"></div>
    `;

		this.buttonElem = this.rootElem.getElementsByClassName('multi-icon-picker-button')[0] as HTMLAnchorElement;
		const dropdownElem = this.rootElem.getElementsByClassName('multi-icon-picker-dropdown')[0] as HTMLElement;

		this.buttonElem.addEventListener('click', event => {
			event.preventDefault();
		});
		this.buttonElem.addEventListener('touchstart', event => {
			if (dropdownElem.style.display == "block") {
				dropdownElem.style.display = "none";
			} else {
				dropdownElem.style.display = "block";
			}
			event.preventDefault();
		});
		this.buttonElem.addEventListener('touchend', event => {
			event.preventDefault();
		});

		dropdownElem.style.gridTemplateColumns = `repeat(${this.config.numColumns}, 1fr)`;

		this.pickers = config.inputs.map((pickerConfig, i) => {
			const optionContainer = document.createElement('div');
			optionContainer.classList.add('dropdown-option-container');
			dropdownElem.appendChild(optionContainer);

			const option = document.createElement('a');
			option.classList.add('dropdown-option', 'multi-icon-picker-option');
			optionContainer.appendChild(option);
			const picker = new IconPicker(option, modObj, pickerConfig);
			picker.changeEmitter.on(() => {
				this.updateButtonImage();
			});
			return picker;
		});
		simUI.sim.waitForInit().then(() => this.updateButtonImage());
	}

	private updateButtonImage() {
		this.currentValue = this.getMaxValue();

		if (this.currentValue) {
			this.rootElem.classList.add('active');
			this.currentValue.fillAndSet(this.buttonElem, true, true);
		} else {
			this.rootElem.classList.remove('active');
			this.buttonElem.style.backgroundImage = '';
			this.buttonElem.style.backgroundColor = this.config.emptyColor;
		}
	}

	private getMaxValue(): ActionId | null {
		return this.pickers.map(picker => picker.getActionId()).filter(id => id != null)[0] || null;
	}
}
