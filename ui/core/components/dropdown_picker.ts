import { Tooltip } from 'bootstrap';

import { EventID, TypedEvent } from '../typed_event.js';

import { Input, InputConfig } from './input.js';

export interface DropdownValueConfig<T> {
	value: T,
    submenu?: Array<string>,
    headerText?: string,
    tooltip?: string,
}

export interface DropdownPickerConfig<ModObject, T> extends InputConfig<ModObject, T> {
	values: Array<DropdownValueConfig<T>>;
    equals: (a: T|undefined, b: T|undefined) => boolean,
    setOptionContent: (button: HTMLButtonElement, valueConfig: DropdownValueConfig<T>) => void,
    defaultLabel: string,
}

interface DropdownSubmenu {
    path: Array<string>,

    listElem: HTMLUListElement,
}

/** UI Input that uses a dropdown menu. */
export class DropdownPicker<ModObject, T> extends Input<ModObject, T> {
    private readonly config: DropdownPickerConfig<ModObject, T>;
    private valueConfigs: Array<DropdownValueConfig<T>>;

	private readonly buttonElem: HTMLButtonElement;
	private readonly listElem: HTMLUListElement;

    private currentSelection: DropdownValueConfig<T>|null;
    private submenus: Array<DropdownSubmenu>;

	constructor(parent: HTMLElement, modObject: ModObject, config: DropdownPickerConfig<ModObject, T>) {
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

    setOptions(newValueConfigs: Array<DropdownValueConfig<T>>) {
        this.buildDropdown(newValueConfigs);
        this.valueConfigs = newValueConfigs.filter(vc => !vc.headerText);
        this.setInputValue(this.getSourceValue());
    }

    private buildDropdown(valueConfigs: Array<DropdownValueConfig<T>>) {
        this.listElem.innerHTML = '';
        this.submenus = [];
		valueConfigs.forEach(valueConfig => {
            const itemElem = document.createElement('li');
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
                this.config.setOptionContent(buttonElem, valueConfig);
                itemElem.appendChild(buttonElem);

                if (valueConfig.tooltip) {
                    buttonElem.setAttribute('data-bs-toggle', 'tooltip');
                    buttonElem.setAttribute('data-bs-html', 'true');
                    buttonElem.setAttribute('data-bs-title', valueConfig.tooltip);
                    const tooltip = Tooltip.getOrCreateInstance(buttonElem, {
                        animation: false,
                        placement: 'right',
                        fallbackPlacements: ['left', 'bottom'],
                        offset: [0, 10],
                        customClass: 'dropdown-tooltip',
                    });
                }

                buttonElem.addEventListener('click', event => {
                    this.updateValue(valueConfig);
                    this.inputChanged(TypedEvent.nextEventID());
                });
            }

            if (valueConfig.submenu && valueConfig.submenu.length > 0) {
                this.createSubmenu(valueConfig.submenu);
            }
            const submenu = this.getSubmenu(valueConfig.submenu);
            if (submenu) {
                submenu.listElem.appendChild(itemElem);
            } else {
                this.listElem.appendChild(itemElem);
            }
		});
    }

    private getSubmenu(path: Array<string>|undefined): DropdownSubmenu|null {
        if (!path) {
            return null;
        }
        return this.submenus.find(submenu => DropdownPicker.equalPaths(submenu.path, path)) || null;
    }

    private createSubmenu(path: Array<string>): DropdownSubmenu {
        const submenu = this.getSubmenu(path);
        if (submenu) {
            return submenu;
        }

        let parent: DropdownSubmenu|null = null;
        if (path.length > 1) {
            parent = this.createSubmenu(path.slice(0, path.length - 1));
        }

        const itemElem = document.createElement('li');
        itemElem.classList.add('dropdown-picker-item');

        const containerElem = document.createElement('div');
        containerElem.classList.add('dropend');
        itemElem.appendChild(containerElem);

        const titleElem = document.createElement('button');
        titleElem.classList.add('dropdown-item');
        titleElem.setAttribute('data-bs-toggle', 'dropdown');
        titleElem.setAttribute('role', 'button');
        titleElem.setAttribute('aria-expanded', 'false');
        titleElem.textContent = path[path.length - 1] + ' \u00bb';
        containerElem.appendChild(titleElem);

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

    private static equalPaths(a: Array<string>|null|undefined, b: Array<string>|null|undefined): boolean {
        return (a?.length || 0) == (b?.length || 0) && (a || []).every((aVal, i) => aVal == b![i]);
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