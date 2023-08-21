import { EventID, TypedEvent } from '../typed_event.js';

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

		this.inputElem.addEventListener('input', event => {
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

// A string picker which adapts its width to the input.
export class AdaptiveStringPicker<ModObject> extends Input<ModObject, string> {
	private readonly inputElem: HTMLInputElement;

	constructor(parent: HTMLElement, modObject: ModObject, config: InputConfig<ModObject, string>) {
		super(parent, 'adaptive-string-picker-root', modObject, config);

		this.inputElem = document.createElement('input');
		this.inputElem.type = 'text';
		this.inputElem.classList.add('form-control')
		this.rootElem.appendChild(this.inputElem);

		this.init();

		this.inputElem.addEventListener('change', event => {
			this.inputChanged(TypedEvent.nextEventID());
		});
		this.inputElem.addEventListener('input', event => {
			this.updateSize();
		});
		this.updateSize();
	}

	getInputElem(): HTMLElement {
		return this.inputElem;
	}

	getInputValue(): string {
		return this.inputElem.value;
	}

	setInputValue(newValue: string) {
		this.inputElem.value = newValue;
		this.updateSize();
	}

	private updateSize() {
		const newSize = Math.max(3, this.inputElem.value.length);
		if (this.inputElem.size != newSize)
			this.inputElem.size = newSize;
	}
}
