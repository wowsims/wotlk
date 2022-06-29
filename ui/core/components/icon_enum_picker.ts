import { ActionId } from '/tbc/core/proto_utils/action_id.js';
import { EventID, TypedEvent } from '/tbc/core/typed_event.js';

import { Component } from './component.js';
import { Input, InputConfig } from './input.js';

export interface IconEnumValueConfig<ModObject, T> {
	// One of these should be set. If actionId is set, shows the icon for that id. If
	// color is set, shows that color.
	actionId?: ActionId,
	color?: string,

	value: T,

	showWhen?: (obj: ModObject) => boolean,
}

export interface IconEnumPickerConfig<ModObject, T> extends InputConfig<ModObject, T> {
	numColumns: number,

	values: Array<IconEnumValueConfig<ModObject, T>>;

	// Function for comparing two values.
	equals: (a: T, b: T) => boolean,

	// Value that will be considered inactive.
	zeroValue: T,

	backupIconUrl?: (value: T) => ActionId,
}

// Icon-based UI for picking enum values.
// ModObject is the object being modified (Sim, Player, or Target).
export class IconEnumPicker<ModObject, T> extends Input<ModObject, T> {
	private readonly config: IconEnumPickerConfig<ModObject, T>;

	private currentValue: T;

	private readonly buttonElem: HTMLAnchorElement;

	constructor(parent: HTMLElement, modObj: ModObject, config: IconEnumPickerConfig<ModObject, T>) {
		super(parent, 'icon-enum-picker-root', modObj, config);
		this.config = config;
		this.currentValue = this.config.zeroValue;
		this.rootElem.classList.add('dropdown-root');

		this.rootElem.innerHTML = `
			<a class="dropdown-button icon-enum-picker-button"></a>
			<div class="dropdown-panel icon-enum-picker-dropdown"></div>
    `;

		this.buttonElem = this.rootElem.getElementsByClassName('icon-enum-picker-button')[0] as HTMLAnchorElement;
		const dropdownElem = this.rootElem.getElementsByClassName('icon-enum-picker-dropdown')[0] as HTMLElement;

		this.buttonElem.addEventListener('click', event => {
			event.preventDefault();
		});
		this.buttonElem.addEventListener('touchstart', event => {
			if (dropdownElem.style.display == "block") {
				dropdownElem.style.display = "none";
			} else {
				dropdownElem.style.display = "block";
			}
			event.preventDefault();
		});
		this.buttonElem.addEventListener('touchend', event => {
			event.preventDefault();
		});

		dropdownElem.style.gridTemplateColumns = `repeat(${this.config.numColumns}, 1fr)`;

		config.values.forEach((valueConfig, i) => {
			const optionContainer = document.createElement('div');
			optionContainer.classList.add('dropdown-option-container');
			dropdownElem.appendChild(optionContainer);

			const option = document.createElement('a');
			option.classList.add('dropdown-option', 'icon-enum-picker-option');
			optionContainer.appendChild(option);
			this.setImage(option, valueConfig);

			if (valueConfig.showWhen) {
				config.changedEvent(this.modObject).on(eventID => {
					const show = valueConfig.showWhen && valueConfig.showWhen(this.modObject);
					if (show) {
						optionContainer.classList.remove('hide');
					} else {
						optionContainer.classList.add('hide');
					}
				});
			}

			option.addEventListener('click', event => {
				event.preventDefault();
				this.currentValue = valueConfig.value;
				this.inputChanged(TypedEvent.nextEventID());

				// Wowhead tooltips can't seem to detect when an element is hidden while
				// being moused over, and the tooltip doesn't disappear. Patch this by
				// dispatching our own mouseout event.
				option.dispatchEvent(new Event('mouseout'));
			});
			option.addEventListener('touchstart', event => {
				event.preventDefault();
			});
			option.addEventListener('touchend', event => {
				event.preventDefault();
				this.currentValue = valueConfig.value;
				this.inputChanged(TypedEvent.nextEventID());
				dropdownElem.style.display = "none";
			});
		});

		this.init();
	}

	private setActionImage(elem: HTMLAnchorElement, actionId: ActionId) {
		actionId.fillAndSet(elem, true, true);
	}

	private setImage(elem: HTMLAnchorElement, valueConfig: IconEnumValueConfig<ModObject, T>) {
		if (valueConfig.actionId) {
			this.setActionImage(elem, valueConfig.actionId);
		} else {
			elem.style.backgroundImage = '';
			elem.style.backgroundColor = valueConfig.color!;
		}
	}

	getInputElem(): HTMLElement {
		return this.buttonElem;
	}

	getInputValue(): T {
		return this.currentValue;
	}

	setInputValue(newValue: T) {
		this.currentValue = newValue;

		if (!this.config.equals(this.currentValue, this.config.zeroValue)) {
			this.rootElem.classList.add('active');
		} else {
			this.rootElem.classList.remove('active');
		}

		const valueConfig = this.config.values.find(valueConfig => this.config.equals(valueConfig.value, this.currentValue))!;
		if (valueConfig) {
			this.setImage(this.buttonElem, valueConfig);
		} else if (this.config.backupIconUrl) {
			const backupId = this.config.backupIconUrl(this.currentValue);
			this.setActionImage(this.buttonElem, backupId);
			this.rootElem.classList.remove('active');
		}
	}
}
