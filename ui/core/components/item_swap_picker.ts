import { Spec, ItemSlot, ItemSwap } from '../proto/common.js';
import { Player } from '../player.js';
import { Component } from './component.js';
import { IconItemSwapPicker } from './gear_picker.js'
import { Input, InputConfig } from './input.js'
import { SimUI } from '../sim_ui.js';
import { TypedEvent } from '../typed_event.js';

export interface ItemSwapPickerConfig<SpecType extends Spec, T> extends InputConfig<Player<SpecType>, T>{
	itemSlots: Array<ItemSlot>;
}

declare var tippy: any;

export class ItemSwapPicker<SpecType extends Spec, T> extends Component {

	constructor(parentElem: HTMLElement, simUI: SimUI, player: Player<SpecType>, config: ItemSwapPickerConfig<SpecType, T>) {
		super(parentElem, 'item-swap-picker-root');

		this.rootElem.classList.add('input-root', 'input-inline')

		const label = document.createElement("label")
		label.classList.add('form-label')
		label.textContent = "Item Swap"
		this.rootElem.appendChild(label);

		if (config.labelTooltip) {
			tippy(label, {
				'content': config.labelTooltip,
				'allowHTML': true,
			});
		}

		let itemSwapContianer = Input.newGroupContainer();
		itemSwapContianer.classList.add('icon-group');
		this.rootElem.appendChild(itemSwapContianer);

		let swapButton = document.createElement('a')
		swapButton.classList.add("gear-swap-icon")
		swapButton.setAttribute('role', 'button')
		swapButton.setAttribute('href', 'javascript:void(0)')
		swapButton.innerHTML = '<i class="fas fa-arrows-rotate me-1"></i>'
		itemSwapContianer.appendChild(swapButton)

		swapButton.addEventListener('click', event => {
			this.swapWithGear(player, config)
		});

		tippy(swapButton, {
			'content': 'Swap Items with Main Gear',
			'allowHTML': true,
		});

		config.changedEvent(player).on(eventID => {
			const show = !config.showWhen || config.showWhen(player);
			if (show) {
				this.rootElem.classList.remove('hide');
			} else {
				this.rootElem.classList.add('hide');
			}
		});

		config.itemSlots.forEach(itemSlot => {
			new IconItemSwapPicker(itemSwapContianer, simUI, player,itemSlot, config);
		});
	}

	swapWithGear(player : Player<SpecType>, config: ItemSwapPickerConfig<SpecType, T> ) {
		let gear = player.getGear()
		config.itemSlots.forEach(slot => {
			const gearItem = player.getGear().getEquippedItem(slot)
			const swapItem = player.getItemSwapGear().getEquippedItem(slot)
			gear = gear.withEquippedItem(slot, swapItem, player.canDualWield2H())
			player.getItemSwapGear().equipItem(slot, gearItem, player.canDualWield2H())
		})

		let eventID = TypedEvent.nextEventID()
		player.setGear(eventID, gear)

		const itemSwap = player.getItemSwapGear().toProto() as unknown as T
		config.setValue(eventID, player, itemSwap)
	}
	
}



