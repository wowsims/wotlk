import { Tooltip } from 'bootstrap';
// eslint-disable-next-line @typescript-eslint/no-unused-vars
import { element, fragment } from 'tsx-vanilla';

import { Player } from '../player.js';
import { ItemSlot, Spec } from '../proto/common.js';
import { SimUI } from '../sim_ui.js';
import { EventID, TypedEvent } from '../typed_event.js';
import { BooleanPicker } from './boolean_picker.js';
import { Component } from './component.js';
import { IconItemSwapPicker } from './gear_picker.js';
import { Input } from './input.js';

export interface ItemSwapConfig {
	itemSlots: Array<ItemSlot>;
	note?: string;
}

export class ItemSwapPicker<SpecType extends Spec> extends Component {
	private readonly itemSlots: Array<ItemSlot>;
	private readonly enableItemSwapPicker: BooleanPicker<Player<SpecType>>;

	constructor(parentElem: HTMLElement, simUI: SimUI, player: Player<SpecType>, config: ItemSwapConfig) {
		super(parentElem, 'item-swap-picker-root');
		this.itemSlots = config.itemSlots;

		this.enableItemSwapPicker = new BooleanPicker(this.rootElem, player, {
			reverse: true,
			label: 'Enable Item Swapping',
			labelTooltip: 'Allows configuring an Item Swap Set which is used with the <b>Item Swap</b> APL action.',
			extraCssClasses: ['input-inline'],
			getValue: (player: Player<SpecType>) => player.getEnableItemSwap(),
			setValue(eventID: EventID, player: Player<SpecType>, newValue: boolean) {
				player.setEnableItemSwap(eventID, newValue);
			},
			changedEvent: (player: Player<SpecType>) => player.itemSwapChangeEmitter,
		});

		const swapPickerContainer = document.createElement('div');
		swapPickerContainer.classList.add('input-root', 'input-inline');
		this.rootElem.appendChild(swapPickerContainer);

		let noteElem: Element;
		if (config.note) {
			noteElem = this.rootElem.appendChild(<p className="form-text">{config.note}</p>);
		}

		const toggleEnabled = () => {
			if (!player.getEnableItemSwap()) {
				swapPickerContainer.classList.add('hide');
				noteElem?.classList.add('hide');
			} else {
				swapPickerContainer.classList.remove('hide');
				noteElem?.classList.remove('hide');
			}
		};
		player.itemSwapChangeEmitter.on(toggleEnabled);
		toggleEnabled();

		const label = document.createElement('label');
		label.classList.add('form-label');
		label.textContent = 'Item Swap';
		swapPickerContainer.appendChild(label);

		const itemSwapContainer = Input.newGroupContainer();
		itemSwapContainer.classList.add('icon-group');
		swapPickerContainer.appendChild(itemSwapContainer);

		const swapButtonFragment = document.createElement('fragment');
		swapButtonFragment.innerHTML = `
			<a
				href="javascript:void(0)"
				class="gear-swap-icon"
				role="button"
				data-bs-title="Swap with equipped items"
			>
				<i class="fas fa-arrows-rotate me-1"></i>
			</a>
		`;

		const swapButton = swapButtonFragment.children[0] as HTMLElement;
		itemSwapContainer.appendChild(swapButton);

		swapButton.addEventListener('click', _event => this.swapWithGear(TypedEvent.nextEventID(), player));
		Tooltip.getOrCreateInstance(swapButton);

		this.itemSlots.forEach(itemSlot => {
			new IconItemSwapPicker(itemSwapContainer, simUI, player, itemSlot);
		});
	}

	swapWithGear(eventID: EventID, player: Player<SpecType>) {
		let newGear = player.getGear();
		let newIsg = player.getItemSwapGear();

		this.itemSlots.forEach(slot => {
			const gearItem = player.getGear().getEquippedItem(slot);
			const swapItem = player.getItemSwapGear().getEquippedItem(slot);

			newGear = newGear.withEquippedItem(slot, swapItem, player.canDualWield2H());
			newIsg = newIsg.withEquippedItem(slot, gearItem, player.canDualWield2H());
		});

		TypedEvent.freezeAllAndDo(() => {
			player.setGear(eventID, newGear);
			player.setItemSwapGear(eventID, newIsg);
		});
	}
}
