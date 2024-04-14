import { Tooltip } from 'bootstrap';
// eslint-disable-next-line @typescript-eslint/no-unused-vars
import { element, fragment } from 'tsx-vanilla';

import { ActionId } from '../proto_utils/action_id.js';
import { TypedEvent } from '../typed_event.js';
import { Input, InputConfig } from './input.js';

export enum IconEnumPickerDirection {
	Vertical = 'vertical',
	Horizontal = 'Horizontal',
}

export interface IconEnumValueConfig<ModObject, T> {
	value: T;
	// One of these should be set. If actionId is set, shows the icon for that id. If
	// color is set, shows that color.
	actionId?: ActionId;
	color?: string;
	// Text to be displayed on the icon.
	text?: string;
	// Hover tooltip.
	tooltip?: string;

	showWhen?: (obj: ModObject) => boolean;
}

export interface IconEnumPickerConfig<ModObject, T> extends InputConfig<ModObject, T> {
	numColumns?: number;
	values: Array<IconEnumValueConfig<ModObject, T>>;
	// Value that will be considered inactive.
	zeroValue: T;
	// Function for comparing two values.
	// Tooltip that will be shown whne hovering over the icon-picker-button
	tooltip?: string;
	// The direction the menu will open in relative to the root element
	direction?: IconEnumPickerDirection;
	equals: (a: T, b: T) => boolean;
	backupIconUrl?: (value: T) => ActionId;
	showWhen?: (obj: ModObject) => boolean;
}

// Icon-based UI for picking enum values.
// ModObject is the object being modified (Sim, Player, or Target).
export class IconEnumPicker<ModObject, T> extends Input<ModObject, T> {
	private readonly config: IconEnumPickerConfig<ModObject, T>;

	private currentValue: T;

	private readonly buttonElem: HTMLAnchorElement;
	private readonly buttonText: HTMLElement;

	constructor(parent: HTMLElement, modObj: ModObject, config: IconEnumPickerConfig<ModObject, T>) {
		super(parent, 'icon-enum-picker-root', modObj, config);
		this.rootElem.classList.add('icon-picker', 'dropdown');
		this.config = config;
		this.currentValue = this.config.zeroValue;

		if (config.showWhen) {
			config.changedEvent(this.modObject).on(_eventID => {
				const show = config.showWhen && config.showWhen(this.modObject);
				if (!show) this.rootElem.classList.add('hide');
			});
		}

		if (config.tooltip) {
			Tooltip.getOrCreateInstance(this.rootElem, {
				html: true,
				title: config.tooltip,
			});
		}
		this.rootElem.appendChild(
			<>
				<a
					href="javascript:void(0)"
					className="icon-picker-button"
					attributes={{
						'aria-expanded': 'false',
						role: 'button',
					}}
					dataset={{
						bsToggle: 'dropdown',
						bsPlacement: 'bottom',
						whtticon: 'false',
						disableWowheadTouchTooltip: 'true',
					}}>
					<span className="icon-picker-label"></span>
				</a>
				<ul className="dropdown-menu"></ul>
			</>,
		);

		this.buttonElem = this.rootElem.querySelector('.icon-picker-button') as HTMLAnchorElement;
		this.buttonText = this.buttonElem.querySelector('.icon-picker-label') as HTMLElement;
		const dropdownMenu = this.rootElem.querySelector('.dropdown-menu') as HTMLElement;

		if (this.config.numColumns) dropdownMenu.style.gridTemplateColumns = `repeat(${this.config.numColumns}, 1fr)`;

		if (this.config.direction == IconEnumPickerDirection.Horizontal) dropdownMenu.style.gridAutoFlow = 'column';

		config.values.forEach((valueConfig, _i) => {
			const optionContainer = document.createElement('li');
			optionContainer.classList.add('icon-dropdown-option', 'dropdown-option');
			dropdownMenu.appendChild(optionContainer);

			const option = document.createElement('a');
			option.classList.add('icon-picker-button');
			option.dataset.whtticon = 'false';
			option.dataset.disableWowheadTouchTooltip = 'true';
			optionContainer.appendChild(option);
			this.setImage(option, valueConfig);

			if (valueConfig.text != undefined) {
				const optionText = document.createElement('div');
				optionText.classList.add('icon-picker-label');
				optionText.textContent = valueConfig.text;
				option.append(optionText);
			}

			if (valueConfig.tooltip) {
				Tooltip.getOrCreateInstance(option, {
					html: true,
					title: valueConfig.tooltip,
				});
			}

			const show = !valueConfig.showWhen || valueConfig.showWhen(this.modObject);
			if (!show) optionContainer.classList.add('hide');

			if (valueConfig.showWhen) {
				config.changedEvent(this.modObject).on(_eventID => {
					const show = valueConfig.showWhen && valueConfig.showWhen(this.modObject);
					if (show) optionContainer.classList.remove('hide');
					else optionContainer.classList.add('hide');
				});
			}

			option.addEventListener('click', event => {
				event.preventDefault();
				this.currentValue = valueConfig.value;
				this.inputChanged(TypedEvent.nextEventID());
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

	update() {
		super.update();
		this.setActive(this.enabled && !this.config.equals(this.currentValue, this.config.zeroValue));
	}

	getInputElem(): HTMLElement {
		return this.buttonElem;
	}

	getInputValue(): T {
		return this.currentValue;
	}

	setInputValue(newValue: T) {
		this.currentValue = newValue;
		this.setActive(this.enabled && !this.config.equals(this.currentValue, this.config.zeroValue));

		this.buttonText.textContent = '';
		this.buttonText.style.display = 'none';

		const valueConfig = this.config.values.find(valueConfig => this.config.equals(valueConfig.value, this.currentValue))!;
		if (valueConfig) {
			this.setImage(this.buttonElem, valueConfig);
			if (valueConfig.text != undefined) {
				this.buttonText.style.display = 'block';
				this.buttonText.textContent = valueConfig.text;
			}
		} else if (this.config.backupIconUrl) {
			const backupId = this.config.backupIconUrl(this.currentValue);
			this.setActionImage(this.buttonElem, backupId);
			this.setActive(false);
		}
	}

	setActive(active: boolean) {
		if (active) this.buttonElem.classList.add('active');
		else this.buttonElem.classList.remove('active');
	}
}
