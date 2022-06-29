import { EventID, TypedEvent } from '/tbc/core/typed_event.js';

import { Input, InputConfig } from './input.js';

/**
 * Data for creating a boolean picker (checkbox).
 */
export interface BooleanPickerConfig<ModObject> extends InputConfig<ModObject, boolean> {
}

// UI element for picking an arbitrary number field.
export class BooleanPicker<ModObject> extends Input<ModObject, boolean> {
	private readonly inputElem: HTMLInputElement;

	constructor(parent: HTMLElement, modObject: ModObject, config: BooleanPickerConfig<ModObject>) {
		super(parent, 'boolean-picker-root', modObject, config);

		this.inputElem = document.createElement('input');
		this.inputElem.type = 'checkbox';
		this.inputElem.classList.add('boolean-picker-input');
		this.rootElem.appendChild(this.inputElem);

		this.init();

		this.inputElem.addEventListener('change', event => {
			this.inputChanged(TypedEvent.nextEventID());
		});
	}

	getInputElem(): HTMLElement {
		return this.inputElem;
	}

	getInputValue(): boolean {
		return Boolean(this.inputElem.checked);
	}

	setInputValue(newValue: boolean) {
		this.inputElem.checked = newValue;
	}
}
