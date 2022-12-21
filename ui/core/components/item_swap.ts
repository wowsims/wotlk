import { Spec, ItemSwap, ItemSlot, ItemSpec } from '../proto/common.js';
import { Player } from '../player.js';
import { EventID } from '../typed_event.js';
import { Component } from './component.js';
import { IconItemSwapPicker } from './gear_picker.js'
import { Input } from './input.js'
import { ItemSwapGear } from '../proto_utils/item_swap_gear.js';
import { EquippedItem } from '../proto_utils/equipped_item.js';
import { SimUI } from '../sim_ui.js';

export interface ItemSwapPickerConfig<SpecType extends Spec, T> {
	getValue: (player: Player<SpecType>) => ItemSwap,
	setValue: (eventID: EventID, player: Player<SpecType>, newValue: ItemSwap) => void,
	values: Array<ItemSwapIconInputConfig<Player<SpecType>, T>>;
	extraCssClasses?: string[];
	showWhen?: (player: Player<SpecType>) => boolean,
}

export interface ItemSwapIconInputConfig<ModObject, T> {
	itemSlot: ItemSlot,
}

export class ItemSwapPicker<SpecType extends Spec, T> extends Component {
	constructor(parentElem: HTMLElement, simUI: SimUI, player: Player<SpecType>, config: ItemSwapPickerConfig<SpecType, T>) {
		super(parentElem, 'item-swap-picker-root');

		this.rootElem.classList.add('input-root', 'input-inline')

		const label = document.createElement("label")
		label.classList.add('form-label')
		label.textContent = "Item Swap"
		this.rootElem.appendChild(label);

		let itemSwapContainer = Input.newGroupContainer();
		itemSwapContainer.classList.add('icon-group');
		this.rootElem.appendChild(itemSwapContainer);

		const gear = new ItemSwapGear();
		config.values.forEach(value => {
			const fieldName = this.getFieldNameFromItemSlot(value.itemSlot)
			if (!fieldName)
				return

			new IconItemSwapPicker(itemSwapContainer, simUI, player, value.itemSlot, gear, {
				changedEvent: (player: Player<SpecType>) => player.specOptionsChangeEmitter,
				getValue: (player: Player<SpecType>) => {
					const itemSwap = config.getValue(player) as unknown as ItemSwap
					return itemSwap[fieldName];
				},
				setValue: (eventID: EventID, player: Player<SpecType>, newValue: ItemSpec | undefined) => {
					const itemSwap = config.getValue(player) as unknown as ItemSwap
					itemSwap[fieldName] = newValue;
					config.setValue(eventID, player, itemSwap);
				},
				showWhen: config.showWhen,
			})
		});
	}

	getFieldNameFromItemSlot(slot: ItemSlot): keyof ItemSwap | undefined {
		switch (slot) {
			case ItemSlot.ItemSlotMainHand:
				return 'mhItem';
			case ItemSlot.ItemSlotOffHand:
				return 'ohItem';
			case ItemSlot.ItemSlotRanged:
				return 'rangedItem';
		}

		return undefined;
	}
}



