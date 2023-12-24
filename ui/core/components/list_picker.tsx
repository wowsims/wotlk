import { Tooltip } from 'bootstrap';
import { EventID, TypedEvent } from '../typed_event.js';
import { swap } from '../utils.js';

import { Input, InputConfig } from './input.js';

import { element, fragment } from 'tsx-vanilla';

export type ListItemAction = 'create' | 'delete' | 'move' | 'copy';

export interface ListPickerActionsConfig {
	create?: {
		// Whether or not to use an icon for the create action button
		// defaults to FALSE
		useIcon?: boolean	
	}
}

export interface ListPickerConfig<ModObject, ItemType> extends InputConfig<ModObject, Array<ItemType>> {
	itemLabel: string,
	newItem: () => ItemType,
	copyItem: (oldItem: ItemType) => ItemType,
	newItemPicker: (parent: HTMLElement, listPicker: ListPicker<ModObject, ItemType>, index: number, config: ListItemPickerConfig<ModObject, ItemType>) => Input<ModObject, ItemType>,
	actions?: ListPickerActionsConfig
	title?: string,
	titleTooltip?: string,
	inlineMenuBar?: boolean,
	hideUi?: boolean,
	horizontalLayout?: boolean,

	// If set, only actions included in the list are allowed. Otherwise, all actions are allowed.
	allowedActions?: Array<ListItemAction>,
}

const DEFAULT_CONFIG = {
	actions: {
		create: {
			useIcon: false,
		}
	}
}

export interface ListItemPickerConfig<ModObject, ItemType> extends InputConfig<ModObject, ItemType> {
}

interface ItemPickerPair<ItemType> {
	elem: HTMLElement,
	picker: Input<any, ItemType>,
	idx: number,
}

interface ListDragData<ModObject, ItemType> {
	listPicker: ListPicker<ModObject, ItemType>;
	item: ItemPickerPair<ItemType>;
}

var curDragData: ListDragData<any, any>|null = null;

export class ListPicker<ModObject, ItemType> extends Input<ModObject, Array<ItemType>> {
	readonly config: ListPickerConfig<ModObject, ItemType>;
	private readonly itemsDiv: HTMLElement;

	private itemPickerPairs: Array<ItemPickerPair<ItemType>>;

	constructor(parent: HTMLElement, modObject: ModObject, config: ListPickerConfig<ModObject, ItemType>) {
		super(parent, 'list-picker-root', modObject, config);
		this.config = {...DEFAULT_CONFIG, ...config};
		this.itemPickerPairs = [];

		this.rootElem.appendChild(
			<>
				{config.title &&
					<label className='list-picker-title form-label'>
						{config.title}
					</label>
				}
				<div className="list-picker-items"></div>
			</>
		)

		if (this.config.hideUi) {
			this.rootElem.classList.add('d-none');
		}
		if (this.config.horizontalLayout) {
			this.config.inlineMenuBar = true;
			this.rootElem.classList.add('horizontal');
		}

		if (this.config.titleTooltip) {
			let cfg = {
				title: this.config.titleTooltip
			}
			Tooltip.getOrCreateInstance(this.rootElem.querySelector('.list-picker-title') as HTMLElement, cfg);
		}

		this.itemsDiv = this.rootElem.getElementsByClassName('list-picker-items')[0] as HTMLElement;

		if (this.actionEnabled('create')) {
			let newItemButton = null;
			let newButtonTooltip: Tooltip | null = null;
			if (this.config.actions?.create?.useIcon) {
				newItemButton = ListPicker.makeActionElem('link-success', 'fa-plus')
				newButtonTooltip = Tooltip.getOrCreateInstance(newItemButton, {title: `New ${config.itemLabel}`});
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

		const item: ItemPickerPair<ItemType> = { elem: itemContainer, picker: itemPicker, idx: index };

		if (this.actionEnabled('move')) {
			const moveButton = ListPicker.makeActionElem('list-picker-item-move', 'fa-arrows-up-down');
			itemHeader.appendChild(moveButton);

			const moveButtonTooltip = Tooltip.getOrCreateInstance(moveButton, {title: 'Move (Drag+Drop)'});
			moveButton.addEventListener('click', event => {
				moveButtonTooltip.hide();
			});

			moveButton.draggable = true;
			moveButton.ondragstart = event => {
				if (event.target == moveButton) {
					event.dataTransfer!.dropEffect = 'move';
					event.dataTransfer!.effectAllowed = 'move';
					itemContainer.classList.add('dragfrom');
					curDragData = {
						listPicker: this,
						item: item,
					};
				}
			};

			let dragEnterCounter = 0;
			itemContainer.ondragenter = event => {
				if (!curDragData || curDragData.listPicker != this) {
					return;
				}
				event.preventDefault();
				dragEnterCounter++;
				itemContainer.classList.add('dragto');
			};
			itemContainer.ondragleave = event => {
				if (!curDragData || curDragData.listPicker != this) {
					return;
				}
				event.preventDefault();
				dragEnterCounter--;
				if (dragEnterCounter <= 0) {
					itemContainer.classList.remove('dragto');
				}
			};
			itemContainer.ondragover = event => {
				if (!curDragData || curDragData.listPicker != this) {
					return;
				}
				event.preventDefault();
			};
			itemContainer.ondrop = event => {
				if (!curDragData || curDragData.listPicker != this) {
					return;
				}
				event.preventDefault();
				dragEnterCounter = 0;
				itemContainer.classList.remove('dragto');
				curDragData.item.elem.classList.remove('dragfrom');

				const srcIdx = curDragData.item.idx;
				const dstIdx = index;
				const newList = this.config.getValue(this.modObject);
				const arrElem = newList[srcIdx];
				newList.splice(srcIdx, 1);
				newList.splice(dstIdx, 0, arrElem);
				this.config.setValue(TypedEvent.nextEventID(), this.modObject, newList);

				curDragData = null;
			};
		}

		if (this.actionEnabled('copy')) {
			const copyButton = ListPicker.makeActionElem('list-picker-item-copy', 'fa-copy');
			itemHeader.appendChild(copyButton);
			const copyButtonTooltip = Tooltip.getOrCreateInstance(copyButton, {title: `Copy to New ${this.config.itemLabel}`});

			copyButton.addEventListener('click', event => {
				const newList = this.config.getValue(this.modObject).slice();
				newList.splice(index, 0, this.config.copyItem(newList[index]));
				this.config.setValue(TypedEvent.nextEventID(), this.modObject, newList);
				copyButtonTooltip.hide();
			});
		}

		if (this.actionEnabled('delete')) {
			const deleteButton = ListPicker.makeActionElem('list-picker-item-delete', 'fa-times');
			deleteButton.classList.add('link-danger');
			itemHeader.appendChild(deleteButton);
			const deleteButtonTooltip = Tooltip.getOrCreateInstance(deleteButton, { title: `Delete ${this.config.itemLabel}`});

			deleteButton.addEventListener('click', event => {
				const newList = this.config.getValue(this.modObject);
				newList.splice(index, 1);
				this.config.setValue(TypedEvent.nextEventID(), this.modObject, newList);
				deleteButtonTooltip.hide();
			});
		}

		this.itemPickerPairs.push(item);
	}

	static makeActionElem(cssClass: string, iconCssClass: string): HTMLElement {
		const actionElem = document.createElement('a');
		actionElem.classList.add('list-picker-item-action', cssClass);
		actionElem.href = 'javascript:void(0)';
		actionElem.setAttribute('role', 'button');

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
