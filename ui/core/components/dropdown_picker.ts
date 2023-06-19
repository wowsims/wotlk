import { EventID, TypedEvent } from '../typed_event.js';

import { Input, InputConfig } from './input.js';

export interface DropdownValueConfig<T> {
	value: T,
}

export interface DropdownPickerConfig<ModObject, T> extends InputConfig<ModObject, T> {
	values: Array<DropdownValueConfig<T>>;
    equals: (a: T|undefined, b: T|undefined) => boolean,
    setOptionContent: (button: HTMLButtonElement, valueConfig: DropdownValueConfig<T>) => void,
    defaultLabel: string,
}

/** UI Input that uses a dropdown menu. */
export class DropdownPicker<ModObject, T> extends Input<ModObject, T> {
    private readonly config: DropdownPickerConfig<ModObject, T>;
    private valueConfigs: Array<DropdownValueConfig<T>>;

	private readonly buttonElem: HTMLButtonElement;
	private readonly listElem: HTMLUListElement;

    private currentSelection: DropdownValueConfig<T>|null;

	constructor(parent: HTMLElement, modObject: ModObject, config: DropdownPickerConfig<ModObject, T>) {
		super(parent, 'dropdown-picker-root', modObject, config);
        this.config = config;
        this.valueConfigs = this.config.values;
        this.currentSelection = null;

        this.rootElem.classList.add('dropdown');

		this.buttonElem = document.createElement('button');
		this.buttonElem.classList.add('dropdown-picker-button', 'btn', 'dropdown-toggle', 'hidden-arrow');
        this.buttonElem.setAttribute('data-bs-toggle', 'dropdown');
        this.buttonElem.textContent = config.defaultLabel;
		this.rootElem.appendChild(this.buttonElem);

		this.listElem = document.createElement('ul');
		this.listElem.classList.add('dropdown-picker-list', 'dropdown-menu');
		this.rootElem.appendChild(this.listElem);

        this.buildDropdown(this.valueConfigs);
		this.init();
	}

    setOptions(newValueConfigs: Array<DropdownValueConfig<T>>) {
        this.valueConfigs = newValueConfigs;
        this.buildDropdown(this.valueConfigs);
        this.setInputValue(this.getSourceValue());
    }

    private buildDropdown(valueConfigs: Array<DropdownValueConfig<T>>) {
        this.listElem.innerHTML = '';
		valueConfigs.forEach(valueConfig => {
            const itemElem = document.createElement('li');
            itemElem.classList.add('dropdown-picker-item');

            const buttonElem = document.createElement('button');
            buttonElem.classList.add('dropdown-item');
            buttonElem.type = 'button';
            this.config.setOptionContent(buttonElem, valueConfig);

            itemElem.appendChild(buttonElem);

            buttonElem.addEventListener('click', event => {
                this.updateValue(valueConfig);
                this.inputChanged(TypedEvent.nextEventID());
            });
            this.listElem.appendChild(itemElem);
		});
    }

	getInputElem(): HTMLElement {
		return this.listElem;
	}

	getInputValue(): T {
		return this.currentSelection?.value as T;
	}

	setInputValue(newValue: T) {
        const newSelection = this.valueConfigs.find(v => this.config.equals(v.value, newValue))!;
        this.updateValue(newSelection);
	}

    private updateValue(newValue: DropdownValueConfig<T>|null) {
        this.currentSelection = newValue;

        // Update button
        if (newValue) {
            this.buttonElem.innerHTML = '';
            this.config.setOptionContent(this.buttonElem, newValue);
        } else {
            this.buttonElem.textContent = this.config.defaultLabel;
        }
    }
}

export interface TextDropdownValueConfig<T> extends DropdownValueConfig<T> {
    label: string,
}

export interface TextDropdownPickerConfig<ModObject, T> extends Omit<DropdownPickerConfig<ModObject, T>, 'values' | 'setOptionContent'> {
	values: Array<TextDropdownValueConfig<T>>,
}

export class TextDropdownPicker<ModObject, T> extends DropdownPicker<ModObject, T> {
	constructor(parent: HTMLElement, modObject: ModObject, config: TextDropdownPickerConfig<ModObject, T>) {
		super(parent, modObject, {
			...config,
            setOptionContent: (button: HTMLButtonElement, valueConfig: DropdownValueConfig<T>) => {
                button.textContent = (valueConfig as TextDropdownValueConfig<T>).label;
            }
		});
	}
}