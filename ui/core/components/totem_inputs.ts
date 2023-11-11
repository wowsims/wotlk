import { IconPicker } from '../components/icon_picker.js';
import * as InputHelpers from '../components/input_helpers.js';
import { IndividualSimUI } from '../individual_sim_ui.js';
import { Player } from '../player.js';
import { Spec } from '../proto/common.js';
import {
	ShamanTotems
} from '../proto/shaman.js';
import { ActionId } from '../proto_utils/action_id.js';
import { ShamanSpecs } from '../proto_utils/utils.js';
import { EventID } from '../typed_event.js';
import { ContentBlock } from './content_block.js';
import { Input } from './input.js';

export function TotemsSection(parentElem: HTMLElement, simUI: IndividualSimUI<ShamanSpecs>): ContentBlock {
	let contentBlock = new ContentBlock(parentElem, 'totems-settings', {
		header: { title: 'Totems' }
	});

	let totemDropdownGroup = Input.newGroupContainer();
	totemDropdownGroup.classList.add('totem-dropdowns-container', 'icon-group');

	let fireElementalContainer = document.createElement('div');
	fireElementalContainer.classList.add('fire-elemental-input-container');

	contentBlock.bodyElement.appendChild(totemDropdownGroup);
	contentBlock.bodyElement.appendChild(fireElementalContainer);

	// Enchancement Shaman uses the Fire Elemental Inputs with custom inputs.
	if (simUI.player.spec != Spec.SpecEnhancementShaman) {
		const fireElementalBooleanIconInput = InputHelpers.makeBooleanIconInput<ShamanSpecs, ShamanTotems, Player<ShamanSpecs>>({
			getModObject: (player: Player<ShamanSpecs>) => player,
			getValue: (player: Player<ShamanSpecs>) => player.getSpecOptions().totems || ShamanTotems.create(),
			setValue: (eventID: EventID, player: Player<ShamanSpecs>, newVal: ShamanTotems) => {
				const newOptions = player.getSpecOptions();
				newOptions.totems = newVal;
				player.setSpecOptions(eventID, newOptions);
			},
			changeEmitter: (player: Player<Spec.SpecEnhancementShaman>) => player.specOptionsChangeEmitter,
		}, ActionId.fromSpellId(2894), "useFireElemental");

		new IconPicker(fireElementalContainer, simUI.player, fireElementalBooleanIconInput);
	}

	return contentBlock;
}
