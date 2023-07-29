import { UnitReference } from '../proto/common.js';
import { ActionId } from '../proto_utils/action_id.js';
import { DropdownPicker, DropdownPickerConfig, DropdownValueConfig } from './dropdown_picker.js';

export interface UnitValue {
    value: UnitReference|undefined,
	text?: string,
	iconUrl?: string|ActionId,
	color?: string,
}

export interface UnitValueConfig extends DropdownValueConfig<UnitValue> {}
export interface UnitPickerConfig<ModObject> extends Omit<DropdownPickerConfig<ModObject, UnitReference|undefined, UnitValue>, 'equals' | 'setOptionContent' | 'defaultLabel'> {
    hideLabelWhenDefaultSelected?: boolean,
}

export class UnitPicker<ModObject> extends DropdownPicker<ModObject, UnitReference|undefined, UnitValue> {
	constructor(parent: HTMLElement, modObject: ModObject, config: UnitPickerConfig<ModObject>) {
		super(parent, modObject, {
			...config,
			equals: (a, b) => UnitReference.equals(a?.value || UnitReference.create(), b?.value || UnitReference.create()),
            defaultLabel: 'Unit',
			setOptionContent: (button: HTMLButtonElement, valueConfig: DropdownValueConfig<UnitValue>, isSelectButton: boolean) => {
                const unitConfig = valueConfig.value;

                if (unitConfig.color) {
                    button.style.backgroundColor = unitConfig.color;
                }

                if (unitConfig.iconUrl) {
                    let icon = null;
                    if (unitConfig.iconUrl instanceof ActionId) {
                        const img = document.createElement('img');
                        img.classList.add('unit-picker-item-icon');
                        unitConfig.iconUrl.fill().then(filledId => {
                            img.src = filledId.iconUrl;
                        });
                        icon = img;
                    } else if (unitConfig.iconUrl.startsWith('fa-')) {
                        const img = document.createElement('span');
                        img.classList.add('fa', unitConfig.iconUrl);
                        img.classList.add('unit-picker-item-label');
                        icon = img;
                    } else {
                        const img = document.createElement('img');
                        img.src = unitConfig.iconUrl;
                        img.classList.add('unit-picker-item-icon');
                        icon = img;
                    }
                    button.appendChild(icon);
                }

                const hideLabel = config.hideLabelWhenDefaultSelected && isSelectButton && !unitConfig.value;
                if (unitConfig.text && !hideLabel) {
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
