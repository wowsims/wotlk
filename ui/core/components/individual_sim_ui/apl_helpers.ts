import { ActionId } from '../../proto_utils/action_id.js';
import { Player } from '../../player.js';
import { DropdownPicker, DropdownPickerConfig, TextDropdownPicker } from '../dropdown_picker.js';


export interface APLActionIDPickerConfig<ModObject> extends Omit<DropdownPickerConfig<ModObject, ActionId>, 'equals' | 'setOptionContent'> {
}

export class APLActionIDPicker extends DropdownPicker<Player<any>, ActionId> {
	constructor(parent: HTMLElement, player: Player<any>, config: APLActionIDPickerConfig<Player<any>>) {
		super(parent, player, {
			...config,
			equals: (a, b) => ((a == null) == (b == null)) && (!a || a.equals(b!)),
            setOptionContent: (button, valueConfig) => {
				const actionId = valueConfig.value;

				const iconElem = document.createElement('a');
				iconElem.classList.add('apl-actionid-item-icon');
				actionId.setBackgroundAndHref(iconElem);
				button.appendChild(iconElem);

				const textElem = document.createTextNode(actionId.name);
				button.appendChild(textElem);
			},
		});
	}
}