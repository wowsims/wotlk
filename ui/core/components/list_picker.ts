import { Tooltip } from 'bootstrap';
import { EventID, TypedEvent } from '../typed_event.js';
import { swap } from '../utils.js';

import { Input, InputConfig } from './input.js';

export interface ListPickerConfig<ModObject, ItemType> extends InputConfig<ModObject, Array<ItemType>> {
	title?: string,
	titleTooltip?: string,
	itemLabel: string,
	newItem: () => ItemType,
	copyItem: (oldItem: ItemType) => ItemType,
	newItemPicker: (parent: HTMLElement, listPicker: ListPicker<ModObject, ItemType>, index: number, config: ListItemPickerConfig<ModObject, ItemType>) => Input<ModObject, ItemType>,
	inlineMenuBar?: boolean,
	hideUi?: boolean,
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
			<button class="list-picker-new-button btn btn-primary">New ${config.itemLabel}</button>
		`;

		if (this.config.hideUi) {
			this.rootElem.classList.add('hide-ui');
		}

		if (this.config.titleTooltip)
			Tooltip.getOrCreateInstance(this.rootElem.querySelector('.list-picker-title') as HTMLElement);

		this.itemsDiv = this.rootElem.getElementsByClassName('list-picker-items')[0] as HTMLElement;

		const newItemButton = this.rootElem.getElementsByClassName('list-picker-new-button')[0] as HTMLElement;
		newItemButton.addEventListener('click', event => {
			const newItem = this.config.newItem();
			const newList = this.config.getValue(this.modObject).concat([newItem]);
			this.config.setValue(TypedEvent.nextEventID(), this.modObject, newList);
		});

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

	private addNewPicker() {
		const index = this.itemPickerPairs.length;
		const itemContainer = document.createElement('div');
		itemContainer.classList.add('list-picker-item-container');
		if (this.config.inlineMenuBar) {
			itemContainer.classList.add('inline');
		}

		const itemHTML = '<div class="list-picker-item"></div>';
		itemContainer.innerHTML = `
			${this.config.inlineMenuBar ? itemHTML : ''}
			<div class="list-picker-item-header">
				${this.config.itemLabel && !this.config.inlineMenuBar ? `<h6 class="list-picker-item-title">${this.config.itemLabel} ${this.itemPickerPairs.length + 1}</h6>` : ''}
				<a href="javascript:void(0)" class="list-picker-item-action list-picker-item-up" role="button" data-bs-toggle="tooltip" data-bs-title="Move Up">
					<i class="fa fa-angle-up fa-xl"></i>
				</a>
				<a href="javascript:void(0)" class="list-picker-item-action list-picker-item-down" role="button" data-bs-toggle="tooltip" data-bs-title="Move Down">
					<i class="fa fa-angle-down fa-xl"></i>
				</a>
				<a href="javascript:void(0)" class="list-picker-item-action list-picker-item-copy" role="button" data-bs-toggle="tooltip" data-bs-title="Copy to New ${this.config.itemLabel}">
					<i class="fa fa-copy fa-xl"></i>
				</a>
				<a href="javascript:void(0)" class="list-picker-item-action list-picker-item-delete link-danger" role="button" data-bs-toggle="tooltip" data-bs-title="Delete ${this.config.itemLabel}">
					<i class="fa fa-times fa-xl"></i>
				</a>
			</div>
			${!this.config.inlineMenuBar ? itemHTML : ''}
		`;

		const upButton = itemContainer.getElementsByClassName('list-picker-item-up')[0] as HTMLElement;
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

		const downButton = itemContainer.getElementsByClassName('list-picker-item-down')[0] as HTMLElement;
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

		const copyButton = itemContainer.getElementsByClassName('list-picker-item-copy')[0] as HTMLElement;
		const copyButtonTooltip = Tooltip.getOrCreateInstance(copyButton);

		copyButton.addEventListener('click', event => {
			const newList = this.config.getValue(this.modObject)
			newList.push(this.config.copyItem(newList[index]));
			this.config.setValue(TypedEvent.nextEventID(), this.modObject, newList);
			copyButtonTooltip.hide();
		});

		const deleteButton = itemContainer.getElementsByClassName('list-picker-item-delete')[0] as HTMLElement;
		const deleteButtonTooltip = Tooltip.getOrCreateInstance(deleteButton);

		deleteButton.addEventListener('click', event => {
			const newList = this.config.getValue(this.modObject);
			newList.splice(index, 1);
			this.config.setValue(TypedEvent.nextEventID(), this.modObject, newList);
			deleteButtonTooltip.hide();
		});

		const itemElem = itemContainer.getElementsByClassName('list-picker-item')[0] as HTMLElement;
		const itemPicker = this.config.newItemPicker(itemElem, this, index, {
			changedEvent: this.config.changedEvent,
			getValue: (modObj: ModObject) => this.config.getValue(modObj)[index],
			setValue: (eventID: EventID, modObj: ModObject, newValue: ItemType) => {
				const newList = this.config.getValue(modObj);
				newList[index] = newValue;
				this.config.setValue(eventID, modObj, newList);
			},
		});
		this.itemsDiv.appendChild(itemContainer);

		this.itemPickerPairs.push({ elem: itemContainer, picker: itemPicker });
	}
}
