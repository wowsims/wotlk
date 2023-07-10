import { Tooltip } from 'bootstrap';
import { EventID, TypedEvent } from '../typed_event.js';
import { swap } from '../utils.js';

import { Input, InputConfig } from './input.js';

export type ListItemAction = 'create' | 'delete' | 'move' | 'copy';

export interface ListPickerConfig<ModObject, ItemType> extends InputConfig<ModObject, Array<ItemType>> {
	title?: string,
	titleTooltip?: string,
	itemLabel: string,
	newItem: () => ItemType,
	copyItem: (oldItem: ItemType) => ItemType,
	newItemPicker: (parent: HTMLElement, listPicker: ListPicker<ModObject, ItemType>, index: number, config: ListItemPickerConfig<ModObject, ItemType>) => Input<ModObject, ItemType>,
	inlineMenuBar?: boolean,
	hideUi?: boolean,
	horizontalLayout?: boolean,

	// If set, only actions included in the list are allowed. Otherwise, all actions are allowed.
	allowedActions?: Array<ListItemAction>,
}

export interface ListItemPickerConfig<ModObject, ItemType> extends InputConfig<ModObject, ItemType> {
}

interface ItemPickerPair<ItemType> {
	elem: HTMLElement,
	picker: Input<any, ItemType>,
}

export class ListPicker<ModObject, ItemType> extends Input<ModObject, Array<ItemType>> {
	readonly config: ListPickerConfig<ModObject, ItemType>;
	private readonly itemsDiv: HTMLElement;

	private itemPickerPairs: Array<ItemPickerPair<ItemType>>;

	constructor(parent: HTMLElement, modObject: ModObject, config: ListPickerConfig<ModObject, ItemType>) {
		super(parent, 'list-picker-root', modObject, config);
		this.config = config;
		this.itemPickerPairs = [];

		this.rootElem.innerHTML = `
			${config.title ? `
				<label
					class="list-picker-title form-label"
					${this.config.titleTooltip ? 'data-bs-toggle="tooltip"' : ''}
					${this.config.titleTooltip ? `data-bs-title="${this.config.titleTooltip}"` : ''}
				>${config.title}</label>` : ''
			}
			<div class="list-picker-items"></div>
		`;

		if (this.config.hideUi) {
			this.rootElem.classList.add('hide-ui');
		}
		if (this.config.horizontalLayout) {
			this.config.inlineMenuBar = true;
			this.rootElem.classList.add('horizontal');
		}

		if (this.config.titleTooltip)
			Tooltip.getOrCreateInstance(this.rootElem.querySelector('.list-picker-title') as HTMLElement);

		this.itemsDiv = this.rootElem.getElementsByClassName('list-picker-items')[0] as HTMLElement;

		if (this.actionEnabled('create')) {
			let newItemButton = null;
			let newButtonTooltip: Tooltip | null = null;
			if (this.config.horizontalLayout) {
				newItemButton = ListPicker.makeActionElem('link-success', `New ${config.itemLabel}`, 'fa-plus')
				newButtonTooltip = Tooltip.getOrCreateInstance(newItemButton);
			} else {
				newItemButton = document.createElement('button');
				newItemButton.classList.add('btn', 'btn-primary');
				newItemButton.textContent = `New ${config.itemLabel}`;
			}
			newItemButton.classList.add('list-picker-new-button');
			newItemButton.addEventListener('click', event => {
				const newItem = this.config.newItem();
				const newList = this.config.getValue(this.modObject).concat([newItem]);
				this.config.setValue(TypedEvent.nextEventID(), this.modObject, newList);
				if (newButtonTooltip) {
					newButtonTooltip.hide();
				}
			});
			this.rootElem.appendChild(newItemButton);
		}

		this.init();
	}

	getInputElem(): HTMLElement {
		return this.rootElem;
	}

	getInputValue(): Array<ItemType> {
		return this.itemPickerPairs.map(pair => pair.picker.getInputValue());
	}

	setInputValue(newValue: Array<ItemType>): void {
		// Add/remove pickers to make the lengths match.
		if (newValue.length < this.itemPickerPairs.length) {
			this.itemPickerPairs.slice(newValue.length).forEach(ipp => ipp.elem.remove());
			this.itemPickerPairs = this.itemPickerPairs.slice(0, newValue.length);
		} else if (newValue.length > this.itemPickerPairs.length) {
			const numToAdd = newValue.length - this.itemPickerPairs.length;
			for (let i = 0; i < numToAdd; i++) {
				this.addNewPicker();
			}
		}

		// Set all the values.
		newValue.forEach((val, i) => this.itemPickerPairs[i].picker.setInputValue(val))
	}

	private actionEnabled(action: ListItemAction): boolean {
		return !this.config.allowedActions || this.config.allowedActions.includes(action);
	}

	private addNewPicker() {
		const index = this.itemPickerPairs.length;
		const itemContainer = document.createElement('div');
		itemContainer.classList.add('list-picker-item-container');
		if (this.config.inlineMenuBar) {
			itemContainer.classList.add('inline');
		}
		this.itemsDiv.appendChild(itemContainer);

		const itemElem = document.createElement('div');
		itemElem.classList.add('list-picker-item');

		const itemHeader = document.createElement('div');
		itemHeader.classList.add('list-picker-item-header');
		const itemHTML = '<div class="list-picker-item"></div>';

		if (this.config.inlineMenuBar) {
			itemContainer.appendChild(itemElem);
			itemContainer.appendChild(itemHeader);
		} else {
			itemContainer.appendChild(itemHeader);
			itemContainer.appendChild(itemElem);
			if (this.config.itemLabel) {
				const itemLabel = document.createElement('h6');
				itemLabel.classList.add('list-picker-item-title');
				itemLabel.textContent = `${this.config.itemLabel} ${this.itemPickerPairs.length + 1}`;
				itemHeader.appendChild(itemLabel);
			}
		}

		const itemPicker = this.config.newItemPicker(itemElem, this, index, {
			changedEvent: this.config.changedEvent,
			getValue: (modObj: ModObject) => this.getSourceValue()[index],
			setValue: (eventID: EventID, modObj: ModObject, newValue: ItemType) => {
				const newList = this.getSourceValue();
				newList[index] = newValue;
				this.config.setValue(eventID, modObj, newList);
			},
		});

		if (this.actionEnabled('move')) {
			const upButton = ListPicker.makeActionElem('list-picker-item-up', 'Move Up', 'fa-angle-up');
			itemHeader.appendChild(upButton);
			const upButtonTooltip = Tooltip.getOrCreateInstance(upButton);

			upButton.addEventListener('click', event => {
				if (index == 0) {
					return;
				}

				const newList = this.config.getValue(this.modObject);
				swap(newList, index, index - 1);
				this.config.setValue(TypedEvent.nextEventID(), this.modObject, newList);
				upButtonTooltip.hide();
			});

			const downButton = ListPicker.makeActionElem('list-picker-item-down', 'Move Down', 'fa-angle-down');
			itemHeader.appendChild(downButton);
			const downButtonTooltip = Tooltip.getOrCreateInstance(downButton);

			downButton.addEventListener('click', event => {
				if (index == this.itemPickerPairs.length - 1) {
					return;
				}

				const newList = this.config.getValue(this.modObject);
				swap(newList, index, index + 1);
				this.config.setValue(TypedEvent.nextEventID(), this.modObject, newList);
				downButtonTooltip.hide();
			});
		}

		if (this.actionEnabled('copy')) {
			const copyButton = ListPicker.makeActionElem('list-picker-item-copy', `Copy to New ${this.config.itemLabel}`, 'fa-copy');
			itemHeader.appendChild(copyButton);
			const copyButtonTooltip = Tooltip.getOrCreateInstance(copyButton);

			copyButton.addEventListener('click', event => {
				const newList = this.config.getValue(this.modObject)
				newList.push(this.config.copyItem(newList[index]));
				this.config.setValue(TypedEvent.nextEventID(), this.modObject, newList);
				copyButtonTooltip.hide();
			});
		}

		if (this.actionEnabled('delete')) {
			const deleteButton = ListPicker.makeActionElem('list-picker-item-delete', `Delete ${this.config.itemLabel}`, 'fa-times');
			deleteButton.classList.add('link-danger');
			itemHeader.appendChild(deleteButton);
			const deleteButtonTooltip = Tooltip.getOrCreateInstance(deleteButton);

			deleteButton.addEventListener('click', event => {
				const newList = this.config.getValue(this.modObject);
				newList.splice(index, 1);
				this.config.setValue(TypedEvent.nextEventID(), this.modObject, newList);
				deleteButtonTooltip.hide();
			});
		}

		this.itemPickerPairs.push({ elem: itemContainer, picker: itemPicker });
	}

	static makeActionElem(cssClass: string, title: string, iconCssClass: string): HTMLElement {
		const actionElem = document.createElement('a');
		actionElem.classList.add('list-picker-item-action', cssClass);
		actionElem.href = 'javascript:void(0)';
		actionElem.setAttribute('role', 'button');
		actionElem.setAttribute('data-bs-toggle', 'tooltip');
		actionElem.setAttribute('data-bs-title', title);

		const icon = document.createElement('i');
		icon.classList.add('fa', 'fa-xl', iconCssClass);
		actionElem.appendChild(icon);

		return actionElem;
	}

	static getItemHeaderElem(itemPicker: Input<any, any>): HTMLElement {
		const itemElem = itemPicker.rootElem.parentElement!;
		const headerElem = (itemElem.nextElementSibling || itemElem.previousElementSibling);
		if (!headerElem?.classList.contains('list-picker-item-header')) {
			throw new Error('Could not find list item header');
		}
		return headerElem as HTMLElement;
	}
}
