import { Tooltip } from 'bootstrap';

import { TypedEvent } from '../typed_event.js';

import { Input, InputConfig } from './input.js';

export interface DropdownValueConfig<V> {
	value: V,
	submenu?: Array<string | V>,
	headerText?: string,
	tooltip?: string,
	extraCssClasses?: Array<string>,
}

export interface DropdownPickerConfig<ModObject, T, V = T> extends InputConfig<ModObject, T, V> {
	values: Array<DropdownValueConfig<V>>;
	equals: (a: V | undefined, b: V | undefined) => boolean,
	setOptionContent: (button: HTMLButtonElement, valueConfig: DropdownValueConfig<V>, isSelectButton: boolean) => void,
	createMissingValue?: (val: V) => Promise<DropdownValueConfig<V>>,
	defaultLabel: string,
}

interface DropdownSubmenu<V> {
	path: Array<string|V>,

	listElem: HTMLUListElement,
}

/** UI Input that uses a dropdown menu. */
export class DropdownPicker<ModObject, T, V = T> extends Input<ModObject, T, V> {
	private readonly config: DropdownPickerConfig<ModObject, T, V>;
	private valueConfigs: Array<DropdownValueConfig<V>>;

	private readonly buttonElem: HTMLButtonElement;
	private readonly listElem: HTMLUListElement;

	private currentSelection: DropdownValueConfig<V> | null;
	private submenus: Array<DropdownSubmenu<V>>;

	constructor(parent: HTMLElement, modObject: ModObject, config: DropdownPickerConfig<ModObject, T, V>) {
		super(parent, 'dropdown-picker-root', modObject, config);
		this.config = config;
		this.valueConfigs = this.config.values.filter(vc => !vc.headerText);
		this.currentSelection = null;
		this.submenus = [];

		this.rootElem.classList.add('dropdown');

		this.buttonElem = document.createElement('button');
		this.buttonElem.classList.add('dropdown-picker-button', 'btn', 'dropdown-toggle', 'open-on-click');
		this.buttonElem.setAttribute('data-bs-toggle', 'dropdown');
		this.buttonElem.setAttribute('aria-expanded', 'false');
		this.buttonElem.setAttribute('role', 'button');
		this.buttonElem.textContent = config.defaultLabel;
		this.rootElem.appendChild(this.buttonElem);

		this.listElem = document.createElement('ul');
		this.listElem.classList.add('dropdown-picker-list', 'dropdown-menu');
		this.rootElem.appendChild(this.listElem);

		this.buildDropdown(this.valueConfigs);
		this.init();
	}

	setOptions(newValueConfigs: Array<DropdownValueConfig<V>>) {
		this.buildDropdown(newValueConfigs);
		this.valueConfigs = newValueConfigs.filter(vc => !vc.headerText);
		this.setInputValue(this.getSourceValue());
	}

	private buildDropdown(valueConfigs: Array<DropdownValueConfig<V>>) {
		this.listElem.innerHTML = '';
		this.submenus = [];
		valueConfigs.forEach(valueConfig => {
			const itemElem = document.createElement('li');
			const containsSubmenuChildren = valueConfigs.some(vc => vc.submenu?.some(e => !(typeof e == 'string') && this.config.equals(e, valueConfig.value)))
			if (valueConfig.extraCssClasses) {
				itemElem.classList.add(...valueConfig.extraCssClasses);
			}
			if (valueConfig.headerText) {
				itemElem.classList.add('dropdown-picker-header');

				const headerElem = document.createElement('h6');
				headerElem.classList.add('dropdown-header');
				headerElem.textContent = valueConfig.headerText;
				itemElem.appendChild(headerElem);
			} else {
				itemElem.classList.add('dropdown-picker-item');

				const buttonElem = document.createElement('button');
				buttonElem.classList.add('dropdown-item');
				buttonElem.type = 'button';
				this.config.setOptionContent(buttonElem, valueConfig, false);

				if (valueConfig.tooltip) {
					Tooltip.getOrCreateInstance(buttonElem, {
						animation: false,
						placement: 'right',
						fallbackPlacements: ['left', 'bottom'],
						offset: [0, 10],
						customClass: 'dropdown-tooltip',
						html: true,
						title: valueConfig.tooltip
					});
				}

				buttonElem.addEventListener('click', () => {
					this.updateValue(valueConfig);
					this.inputChanged(TypedEvent.nextEventID());
				});

				if (containsSubmenuChildren) {
					this.createSubmenu((valueConfig.submenu || []).concat([valueConfig.value]), buttonElem, itemElem)
				} else {
					itemElem.appendChild(buttonElem);
				}
			}

			if (!containsSubmenuChildren) {
				if (valueConfig.submenu && valueConfig.submenu.length > 0) {
					this.createSubmenu(valueConfig.submenu);
				}
				const submenu = this.getSubmenu(valueConfig.submenu);
				if (submenu) {
					submenu.listElem.appendChild(itemElem);
				} else {
					this.listElem.appendChild(itemElem);
				}
			}
		});
	}

	private getSubmenu(path: Array<string|V> | undefined): DropdownSubmenu<V> | null {
		if (!path) {
			return null;
		}
		return this.submenus.find(submenu => this.equalPaths(submenu.path, path)) || null;
	}

	private createSubmenu(path: Array<string|V>, buttonElem?: HTMLButtonElement, itemElem?: HTMLLIElement): DropdownSubmenu<V> {
		const submenu = this.getSubmenu(path);
		if (submenu) {
			return submenu;
		}

		let parent: DropdownSubmenu<V> | null = null;
		if (path.length > 1) {
			parent = this.createSubmenu(path.slice(0, path.length - 1));
		}

		if (!itemElem) {
			itemElem = document.createElement('li');
		}
		itemElem.classList.add('dropdown-picker-item');

		const containerElem = document.createElement('div');
		containerElem.classList.add('dropend');
		itemElem.appendChild(containerElem);

		if (!buttonElem) {
			buttonElem = document.createElement('button');
		}
		buttonElem.classList.add('dropdown-item');
		buttonElem.setAttribute('data-bs-toggle', 'dropdown');
		buttonElem.setAttribute('role', 'button');
		buttonElem.setAttribute('aria-expanded', 'false');
		if (buttonElem.childNodes.length == 0) {
			buttonElem.textContent = path[path.length - 1] + ' \u00bb';
		}
		containerElem.appendChild(buttonElem);

		const listElem = document.createElement('ul');
		listElem.classList.add('dropdown-submenu', 'dropdown-menu');
		containerElem.appendChild(listElem);

		if (parent) {
			parent.listElem.appendChild(itemElem);
		} else {
			this.listElem.appendChild(itemElem);
		}

		const newSubmenu = {
			path: path,
			listElem: listElem,
		};
		this.submenus.push(newSubmenu);
		return newSubmenu;
	}

	private equalPaths(a: Array<string|V> | null | undefined, b: Array<string|V> | null | undefined): boolean {
		return (a?.length || 0) == (b?.length || 0) &&
			(a || []).every((aVal, i) =>
				(typeof aVal == 'string')
					? aVal == (b![i] as string)
					: this.config.equals(aVal, b![i] as V));
	}

	getInputElem(): HTMLElement {
		return this.listElem;
	}

	getInputValue(): T {
		return this.valueToSource(this.currentSelection?.value as V);
	}

	setInputValue(newSrcValue: T) {
		const newValue = this.sourceToValue(newSrcValue);
		const newSelection = this.valueConfigs.find(v => this.config.equals(v.value, newValue))!;
		if (newSelection) {
			this.updateValue(newSelection);
		} else if (newValue == null) {
			this.updateValue(null);
		} else if (this.config.createMissingValue) {
			this.config.createMissingValue(newValue).then(newSelection => this.updateValue(newSelection));
		} else {
			this.updateValue(null);
		}
	}

	private updateValue(newValue: DropdownValueConfig<V> | null) {
		this.currentSelection = newValue;

		// Update button
		if (newValue) {
			this.buttonElem.innerHTML = '';
			this.config.setOptionContent(this.buttonElem, newValue, true);
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