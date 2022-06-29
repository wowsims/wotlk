import { EventID, TypedEvent } from '/tbc/core/typed_event.js';

import { Input, InputConfig } from './input.js';

/**
 * Data for creating a string picker.
 */
export interface StringPickerConfig<ModObject> extends InputConfig<ModObject, string> {
}

// UI element for picking an arbitrary string field.
export class StringPicker<ModObject> extends Input<ModObject, string> {
	private readonly inputElem: HTMLSpanElement;

	constructor(parent: HTMLElement, modObject: ModObject, config: StringPickerConfig<ModObject>) {
		super(parent, 'string-picker-root', modObject, config);

		this.inputElem = document.createElement('span');
		this.inputElem.setAttribute('contenteditable', '');
		this.inputElem.classList.add('string-picker-input');
		this.rootElem.appendChild(this.inputElem);

		this.init();

		this.inputElem.addEventListener('change', event => {
			this.inputChanged(TypedEvent.nextEventID());
		});
	}

	getInputElem(): HTMLElement {
		return this.inputElem;
	}

	getInputValue(): string {
		return this.inputElem.textContent || '';
	}

	setInputValue(newValue: string) {
		this.inputElem.textContent = newValue;
	}
}
