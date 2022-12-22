import { Spec, ItemSlot } from '../proto/common.js';
import { Player } from '../player.js';
import { Component } from './component.js';
import { IconItemSwapPicker } from './gear_picker.js'
import { Input, InputConfig } from './input.js'

export interface ItemSwapPickerConfig<SpecType extends Spec, T> extends InputConfig<Player<SpecType>, T>{
	itemSlots: Array<ItemSlot>;
}

export class ItemSwapPicker<SpecType extends Spec, T> extends Component {

	constructor(parentElem: HTMLElement, player: Player<SpecType>, config: ItemSwapPickerConfig<SpecType, T>) {
		super(parentElem, 'item-swap-picker-root');

		this.rootElem.classList.add('input-root', 'input-inline')

		const label = document.createElement("label")
		label.classList.add('form-label')
		label.textContent = "Item Swap"
		this.rootElem.appendChild(label);

		let itemSwapContianer = Input.newGroupContainer();
		itemSwapContianer.classList.add('icon-group');
		this.rootElem.appendChild(itemSwapContianer);

		config.itemSlots.forEach(itemSlot => {
			new IconItemSwapPicker(itemSwapContianer, player,itemSlot, config);
		});
	}
}



