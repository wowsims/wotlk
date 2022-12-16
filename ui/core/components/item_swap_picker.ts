import { Spec, ItemSpec, ItemSwap } from '../proto/common.js';
import { ItemSlot } from '../proto/common.js';
import { Player } from '../player.js';
import { EventID } from '../typed_event.js';

import { Input } from './input.js';
import { IndividualSimUI } from '../individual_sim_ui.js';
import { ContentBlock } from './content_block.js';
import {IconItemSwapPicker } from './gear_picker.js';
import { specIconsLarge } from '../proto_utils/utils.js';

declare var tippy: any;
declare var WowSim: any;

export function ItemSwapSection(parentElem: HTMLElement, simUI: IndividualSimUI<Spec.SpecEnhancementShaman>): ContentBlock {
	let contentBlock = new ContentBlock(parentElem, 'item-swap-settings', {
		header: {title: 'Item Swap'}
	});

	let itemSwapContianer = Input.newGroupContainer();
	itemSwapContianer.classList.add('item-swap-inputs-container', 'icon-group');
	contentBlock.bodyElement.appendChild(itemSwapContianer);

	new IconItemSwapPicker(itemSwapContianer, simUI.player, ItemSlot.ItemSlotMainHand, {
		// Returns the event indicating the mapped value has changed.
		changedEvent: (player: Player<Spec.SpecEnhancementShaman>) => player.specOptionsChangeEmitter,

		// Get and set the mapped value.
		getValue: (player: Player<Spec.SpecEnhancementShaman>) => {	
			return player.getSpecOptions().weaponSwap?.mhItem
		},
		setValue: (eventID: EventID, player: Player<Spec.SpecEnhancementShaman>, newValue: ItemSpec | undefined) => {
			const options = player.getSpecOptions()
			options.weaponSwap!.mhItem = newValue;
			player.setSpecOptions(eventID, options)
		},
	})

	// new IconItemSwapPicker(itemSwapContianer, simUI.player, ItemSlot.ItemSlotOffHand, {
	// 	// Returns the event indicating the mapped value has changed.
	// 	changedEvent: (player: Player<Spec.SpecEnhancementShaman>) => player.specOptionsChangeEmitter,

	// 	// Get and set the mapped value.
	// 	getValue: (player: Player<Spec.SpecEnhancementShaman>) => {
	// 		return player.getSpecOptions().weaponSwap?.ohItem
	// 	},
	// 	setValue: (eventID: EventID, player: Player<Spec.SpecEnhancementShaman>, newValue: ItemSpec | undefined) => {
	// 		const options = player.getSpecOptions()
	// 		options.weaponSwap!.ohItem = newValue;
	// 		player.setSpecOptions(eventID, options)
	// 	},
	// })

	return contentBlock
}