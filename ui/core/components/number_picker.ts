import { EventID, TypedEvent } from '/tbc/core/typed_event.js';

import { Input, InputConfig } from './input.js';

/**
 * Data for creating a number picker.
 */
export interface NumberPickerConfig<ModObject> extends InputConfig<ModObject, number> {
	float?: boolean,
}

// UI element for picking an arbitrary number field.
export class NumberPicker<ModObject> extends Input<ModObject, number> {
	private readonly inputElem: HTMLInputElement;
	private float: boolean;

	constructor(parent: HTMLElement, modObject: ModObject, config: NumberPickerConfig<ModObject>) {
		super(parent, 'number-picker-root', modObject, config);
		this.float = config.float || false;

		this.inputElem = document.createElement('input');
		if (this.float) {
			this.inputElem.type = 'text';
			this.inputElem.inputMode = 'numeric';
		} else {
			this.inputElem.type = 'number';
		}
		this.inputElem.classList.add('number-picker-input');
		this.rootElem.appendChild(this.inputElem);

		this.init();

		this.inputElem.addEventListener('change', event => {
			this.inputChanged(TypedEvent.nextEventID());
		});
	}

	getInputElem(): HTMLElement {
		return this.inputElem;
	}

	getInputValue(): number {
		if (this.float) {
			return parseFloat(this.inputElem.value || '') || 0;
		} else {
			return parseInt(this.inputElem.value || '') || 0;
		}
	}

	setInputValue(newValue: number) {
		this.inputElem.value = String(newValue);
	}
}
