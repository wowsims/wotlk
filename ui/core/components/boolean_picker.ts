import { TypedEvent } from '../typed_event.js';
import { Input, InputConfig } from './input.js';

/**
 * Data for creating a boolean picker (checkbox).
 */
export interface BooleanPickerConfig<ModObject> extends InputConfig<ModObject, boolean> {
	reverse?: boolean;
}

// UI element for picking an arbitrary number field.
export class BooleanPicker<ModObject> extends Input<ModObject, boolean> {
	private readonly inputElem: HTMLInputElement;

	constructor(parent: HTMLElement, modObject: ModObject, config: BooleanPickerConfig<ModObject>) {
		super(parent, 'boolean-picker-root', modObject, config);

		this.rootElem.classList.add('form-check');

		this.inputElem = document.createElement('input');
		this.inputElem.type = 'checkbox';
		this.inputElem.classList.add('boolean-picker-input', 'form-check-input');

		if (config.reverse) {
			this.rootElem.classList.add('form-check-reverse');
			this.rootElem.appendChild(this.inputElem);
		} else {
			this.rootElem.prepend(this.inputElem);
		}

		this.init();

		this.inputElem.addEventListener('change', () => {
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
