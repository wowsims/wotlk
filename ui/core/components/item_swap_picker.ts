import { Spec, ItemSlot, ItemSwap } from '../proto/common.js';
import { Player } from '../player.js';
import { Component } from './component.js';
import { IconItemSwapPicker } from './gear_picker.js'
import { Input, InputConfig } from './input.js'
import { SimUI } from '../sim_ui.js';
import { TypedEvent } from '../typed_event.js';
import { mapToStyles } from '@popperjs/core/lib/modifiers/computeStyles.js';

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

		let itemSwapContainer = Input.newGroupContainer();
		itemSwapContainer.classList.add('icon-group')
		this.rootElem.appendChild(itemSwapContainer);

		let swapButtonFragment = document.createElement('fragment');
		swapButtonFragment.innerHTML = `
			<a
				href="javascript:void(0)"
				class="gear-swap-icon"
				role="button"
				data-bs-toggle="tooltip"
				data-bs-title="Swap Items with Main Gear"
			>
				<i class="fas fa-arrows-rotate me-1"></i>
			</a>
		`

		const swapButton = swapButtonFragment.children[0] as HTMLElement;
		itemSwapContainer.appendChild(swapButton)

		swapButton.addEventListener('click', event => { this.swapWithGear(player, config) });

		config.changedEvent(player).on(eventID => {
			const show = !config.showWhen || config.showWhen(player);
			if (show) {
				this.rootElem.classList.remove('hide');
			} else {
				this.rootElem.classList.add('hide');
			}
		});

		config.itemSlots.forEach(itemSlot => {
			new IconItemSwapPicker(itemSwapContainer, simUI, player,itemSlot, config);
		});
	}

	swapWithGear(player : Player<SpecType>, config: ItemSwapPickerConfig<SpecType, T> ) {
		let gear = player.getGear()

		const gearMap = new Map();
		const itemSwapMap = new Map();

		config.itemSlots.forEach(slot => {
			const gearItem = player.getGear().getEquippedItem(slot)
			const swapItem = player.getItemSwapGear().getEquippedItem(slot)

			gearMap.set(slot, gearItem)
			itemSwapMap.set(slot, swapItem)
		})

		itemSwapMap.forEach((item, slot) => {
			gear = gear.withEquippedItem(slot, item, player.canDualWield2H())
		})

		gearMap.forEach((item, slot) => {
			player.getItemSwapGear().equipItem(slot, item, player.canDualWield2H())
		})

		let eventID = TypedEvent.nextEventID()
		player.setGear(eventID, gear)

		const itemSwap = player.getItemSwapGear().toProto() as unknown as T
		config.setValue(eventID, player, itemSwap)
	}
	
}



