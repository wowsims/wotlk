import { UnitReference } from '../proto/common.js';
import { DropdownPicker, DropdownPickerConfig, DropdownValueConfig } from './dropdown_picker.js';

export interface UnitValue {
    value: UnitReference,
	text?: string,
	iconUrl?: string,
	color?: string,
}

export interface UnitValueConfig extends DropdownValueConfig<UnitValue> {}
export interface UnitPickerConfig<ModObject> extends Omit<DropdownPickerConfig<ModObject, UnitReference, UnitValue>, 'equals' | 'setOptionContent' | 'defaultLabel'> {
}

export class UnitPicker<ModObject> extends DropdownPicker<ModObject, UnitReference, UnitValue> {
	constructor(parent: HTMLElement, modObject: ModObject, config: UnitPickerConfig<ModObject>) {
		super(parent, modObject, {
			...config,
			equals: (a, b) => UnitReference.equals(a?.value, b?.value),
            defaultLabel: 'Unit',
			setOptionContent: (button: HTMLButtonElement, valueConfig: DropdownValueConfig<UnitValue>) => {
                const unitConfig = valueConfig.value;

                if (unitConfig.color) {
                    button.style.backgroundColor = unitConfig.color;
                }

                if (unitConfig.iconUrl) {
                    let icon = null;
                    if (unitConfig.iconUrl.startsWith('fa-')) {
                        icon = document.createElement('span');
                        icon.classList.add('fa', unitConfig.iconUrl);
                        icon.classList.add('unit-picker-item-label');
                    } else {
                        icon = document.createElement('img');
                        icon.src = unitConfig.iconUrl;
                        icon.classList.add('unit-picker-item-icon');
                    }
                    button.appendChild(icon);
                }

                if (unitConfig.text) {
                    const label = document.createElement('span');
                    if (unitConfig.text.startsWith('fa-')) {
                        label.classList.add('fa', unitConfig.text);
                    } else {
                        label.textContent = unitConfig.text;
                    }
                    label.classList.add('unit-picker-item-label');
                    button.appendChild(label);
                }
			}
		});
        this.rootElem.classList.add('unit-picker-root');
	}
}
