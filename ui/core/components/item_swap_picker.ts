import { Spec, ItemSlot, ItemSpec } from '../proto/common.js';
import { Player } from '../player.js';
import { Component } from './component.js';
import { IconItemSwapPicker } from './gear_picker.js'
import { Input } from './input.js'
import { SimUI } from '../sim_ui.js';
import { EventID, TypedEvent } from '../typed_event.js';
import { BooleanPicker } from './boolean_picker.js';

export interface ItemSwapPickerConfig {
	itemSlots: Array<ItemSlot>;
}

export class ItemSwapPicker<SpecType extends Spec> extends Component {
	private readonly itemSlots: Array<ItemSlot>;
	private readonly enableItemSwapPicker: BooleanPicker<Player<SpecType>>;

	constructor(parentElem: HTMLElement, simUI: SimUI, player: Player<SpecType>, config: ItemSwapPickerConfig) {
		super(parentElem, 'item-swap-picker-root');
		this.itemSlots = config.itemSlots;

		this.enableItemSwapPicker = new BooleanPicker(this.rootElem, player, {
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
		this.rootElem.appendChild(swapPickerContainer);
		const toggleEnabled = () => {
			if (!player.getEnableItemSwap()) {
				swapPickerContainer.classList.add('hide');
			} else {
				swapPickerContainer.classList.remove('hide');
			}
		};
		player.itemSwapChangeEmitter.on(toggleEnabled);
		toggleEnabled();

		const label = document.createElement("label");
		label.classList.add('form-label');
		label.textContent = "Item Swap";
		swapPickerContainer.appendChild(label);

		let itemSwapContainer = Input.newGroupContainer();
		itemSwapContainer.classList.add('icon-group');
		swapPickerContainer.appendChild(itemSwapContainer);

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
		`;

		const swapButton = swapButtonFragment.children[0] as HTMLElement;
		itemSwapContainer.appendChild(swapButton);

		swapButton.addEventListener('click', _event => { this.swapWithGear(TypedEvent.nextEventID(), player) });

		this.itemSlots.forEach(itemSlot => {
			new IconItemSwapPicker(itemSwapContainer, simUI, player, itemSlot, {
				getValue: (player: Player<any>) => player.getItemSwapGear().getEquippedItem(itemSlot)?.asSpec() || ItemSpec.create(),
				setValue: (eventID: EventID, player: Player<any>, newValue: ItemSpec) => {
					let curIsg = player.getItemSwapGear();
					curIsg = curIsg.withEquippedItem(itemSlot, player.sim.db.lookupItemSpec(newValue), player.canDualWield2H())
					player.setItemSwapGear(eventID, curIsg);
				},
				changedEvent: (player: Player<any>) => player.itemSwapChangeEmitter,
			});
		});
	}

	swapWithGear(eventID: EventID, player: Player<SpecType>) {
		let newGear = player.getGear();
		let newIsg = player.getItemSwapGear();

		this.itemSlots.forEach(slot => {
			const gearItem = player.getGear().getEquippedItem(slot);
			const swapItem = player.getItemSwapGear().getEquippedItem(slot);

			newGear = newGear.withEquippedItem(slot, swapItem, player.canDualWield2H())
			newIsg = newIsg.withEquippedItem(slot, gearItem, player.canDualWield2H())
		});

		TypedEvent.freezeAllAndDo(() => {
			player.setGear(eventID, newGear);
			player.setItemSwapGear(eventID, newIsg);
		});
	}
}