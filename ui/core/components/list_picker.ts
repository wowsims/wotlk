import { Tooltip } from 'bootstrap';
import { SimUI } from '../sim_ui.js';
import { EventID, TypedEvent } from '../typed_event.js';
import { arrayEquals, swap } from '../utils.js';

import { Input, InputConfig } from './input.js';

export interface ListPickerConfig<ModObject, ItemType> extends InputConfig<ModObject, Array<ItemType>> {
	title?: string,
	titleTooltip?: string,
	itemLabel: string,
	newItem: () => ItemType,
	copyItem: (oldItem: ItemType) => ItemType,
	newItemPicker: (parent: HTMLElement, item: ItemType, listPicker: ListPicker<ModObject, ItemType>) => Input<ItemType, ItemType>,
	inlineMenuBar?: boolean,
}

interface ItemPickerPair<ItemType> {
	item: ItemType,
	elem: HTMLElement,
	picker: Input<ItemType, ItemType>,
}

export class ListPicker<ModObject, ItemType> extends Input<ModObject, Array<ItemType>> {
	private readonly config: ListPickerConfig<ModObject, ItemType>;
	private readonly itemsDiv: HTMLElement;

	private itemPickerPairs: Array<ItemPickerPair<ItemType>>;

	constructor(parent: HTMLElement, simUI: SimUI, modObject: ModObject, config: ListPickerConfig<ModObject, ItemType>) {
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
		return this.itemPickerPairs.map(pair => pair.item);
	}

	setInputValue(newValue: Array<ItemType>): void {
		// Remove items that are no longer in the list.
		const removePairs = this.itemPickerPairs.filter(ipp => !newValue.includes(ipp.item));
		removePairs.forEach(ipp => ipp.elem.remove());
		this.itemPickerPairs = this.itemPickerPairs.filter(ipp => !removePairs.includes(ipp));

		// Add items that were missing.
		const curItems = this.getInputValue();
		newValue
			.filter(newItem => !curItems.includes(newItem))
			.forEach(newItem => this.addNewPicker(newItem));

		// Reorder to match the new list.
		this.itemPickerPairs = newValue.map(item => this.itemPickerPairs.find(ipp => ipp.item == item)!);

		// Reorder item picker elements in the DOM if necessary.
		const curPickerElems = Array.from(this.itemsDiv.children);
		if (!curPickerElems.every((elem, i) => elem == this.itemPickerPairs[i].elem)) {
			this.itemPickerPairs.forEach(ipp => ipp.elem.remove());
			this.itemPickerPairs.forEach(ipp => this.itemsDiv.appendChild(ipp.elem));
		}
	}

	getPickerIndex(picker: Input<ItemType, ItemType>): number {
		return this.itemPickerPairs.findIndex(ipp => ipp.picker == picker);
	}

	private addNewPicker(item: ItemType) {
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
			const index = this.itemPickerPairs.findIndex(ipp => ipp.item == item);
			if (index == -1) {
				console.error('Could not find list picker item!');
				return;
			}
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
			const index = this.itemPickerPairs.findIndex(ipp => ipp.item == item);
			if (index == -1) {
				console.error('Could not find list picker item!');
				return;
			}
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
			const index = this.itemPickerPairs.findIndex(ipp => ipp.item == item);
			if (index == -1) {
				console.error('Could not find list picker item!');
				return;
			}

			const copiedItem = this.config.copyItem(item);
			const newList = this.config.getValue(this.modObject).concat([copiedItem]);
			this.config.setValue(TypedEvent.nextEventID(), this.modObject, newList);
			copyButtonTooltip.hide();
		});

		const deleteButton = itemContainer.getElementsByClassName('list-picker-item-delete')[0] as HTMLElement;
		const deleteButtonTooltip = Tooltip.getOrCreateInstance(deleteButton);

		deleteButton.addEventListener('click', event => {
			const index = this.itemPickerPairs.findIndex(ipp => ipp.item == item);
			if (index == -1) {
				console.error('Could not find list picker item!');
				return;
			}

			const newList = this.config.getValue(this.modObject);
			newList.splice(index, 1);
			this.config.setValue(TypedEvent.nextEventID(), this.modObject, newList);
			deleteButtonTooltip.hide();
		});

		const itemElem = itemContainer.getElementsByClassName('list-picker-item')[0] as HTMLElement;
		const itemPicker = this.config.newItemPicker(itemElem, item, this);
		this.itemsDiv.appendChild(itemContainer);

		this.itemPickerPairs.push({ item: item, elem: itemContainer, picker: itemPicker });
	}
}
