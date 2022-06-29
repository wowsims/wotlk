import { Race } from '/tbc/core/proto/common.js';
import { specToEligibleRaces } from '/tbc/core/proto_utils/utils.js';
import { raceNames } from '/tbc/core/proto_utils/names.js';
import { EventID, TypedEvent } from '/tbc/core/typed_event.js';

import { Input, InputConfig } from './input.js';

export interface EnumValueConfig {
	name: string,
	value: number,
	tooltip?: string,
}

export interface EnumPickerConfig<ModObject> extends InputConfig<ModObject, number> {
	values: Array<EnumValueConfig>;
}

export class EnumPicker<ModObject> extends Input<ModObject, number> {
	private readonly selectElem: HTMLSelectElement;

	constructor(parent: HTMLElement, modObject: ModObject, config: EnumPickerConfig<ModObject>) {
		super(parent, 'enum-picker-root', modObject, config);

		this.selectElem = document.createElement('select');
		this.selectElem.classList.add('enum-picker-selector');
		this.rootElem.appendChild(this.selectElem);

		config.values.forEach((value) => {
			const option = document.createElement('option');
			option.value = String(value.value);
			option.textContent = value.name;
			this.selectElem.appendChild(option);

			if (value.tooltip) {
				option.title = value.tooltip;
			}
		});

		this.init();

		this.selectElem.addEventListener('change', event => {
			this.inputChanged(TypedEvent.nextEventID());
		});
	}

	getInputElem(): HTMLElement {
		return this.selectElem;
	}

	getInputValue(): number {
		return parseInt(this.selectElem.value);
	}

	setInputValue(newValue: number) {
		this.selectElem.value = String(newValue);
	}
}
