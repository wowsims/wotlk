import { DropdownPicker, DropdownPickerConfig, DropdownValueConfig } from './dropdown_picker.js';

export interface UnitValueConfig<T> extends DropdownValueConfig<T> {
	text: string,
	iconUrl?: string,
	color?: string,
}

export interface UnitPickerConfig<ModObject, T> extends Omit<DropdownPickerConfig<ModObject, T>, 'values' | 'setOptionContent' | 'defaultLabel'> {
	values: Array<UnitValueConfig<T>>,
}

export class UnitPicker<ModObject, T> extends DropdownPicker<ModObject, T> {
	constructor(parent: HTMLElement, modObject: ModObject, config: UnitPickerConfig<ModObject, T>) {
		super(parent, modObject, {
			...config,
            defaultLabel: 'Unit',
			setOptionContent: (button: HTMLButtonElement, valueConfig: DropdownValueConfig<T>) => {
                const unitConfig = valueConfig as UnitValueConfig<T>;

                if (unitConfig.color) {
                    button.style.backgroundColor = unitConfig.color;
                }

                if (unitConfig.iconUrl) {
                    const icon = document.createElement('img');
                    icon.src = unitConfig.iconUrl;
                    icon.classList.add('unit-picker-item-icon');
                    button.appendChild(icon);
                }

                if (unitConfig.text) {
                    const label = document.createElement('span');
                    label.textContent = unitConfig.text;
                    label.classList.add('unit-picker-item-label');
                    button.appendChild(label);
                }
			}
		});
        this.rootElem.classList.add('unit-picker-root');
	}
}
