import { Spec } from '../../proto/common.js';
import { EventID } from '../../typed_event.js';
import { Player } from '../../player.js';
import { IconEnumPicker, IconEnumValueConfig } from '../icon_enum_picker.js';
import { BooleanPicker } from '../boolean_picker.js';
import { ListPicker } from '../list_picker.js';
import { NumberPicker } from '../number_picker.js';
import { StringPicker } from '../string_picker.js';
import { APLListItem, APLAction } from '../../proto/apl.js';

import { Component } from '../component.js';
import { SimUI } from 'ui/core/sim_ui.js';
import { timeStamp } from 'console';

export class APLRotationPicker extends Component {
	constructor(parent: HTMLElement, simUI: SimUI, modPlayer: Player<any>) {
		super(parent, 'apl-rotation-picker-root');

		new ListPicker<Player<any>, APLListItem, APLListItemPicker>(this.rootElem, simUI, modPlayer, {
			extraCssClasses: ['apl-list-item-picker'],
			title: 'Priority List',
			titleTooltip: 'At each decision point, the simulation will perform the first valid action from this list.',
			itemLabel: 'Action',
			changedEvent: (player: Player<any>) => player.rotationChangeEmitter,
			getValue: (player: Player<any>) => player.getAplRotation().priorityList,
			setValue: (eventID: EventID, player: Player<any>, newValue: Array<APLListItem>) => {
                const rotation = player.getAplRotation();
                rotation.priorityList = newValue;
                player.setAplRotation(eventID, rotation);
			},
			newItem: () => APLListItem.create(),
			copyItem: (oldItem: APLListItem) => APLListItem.clone(oldItem),
			newItemPicker: (parent: HTMLElement, newItem: APLListItem, listPicker: ListPicker<Player<any>, APLListItem, APLListItemPicker>) => new APLListItemPicker(parent, modPlayer, newItem, listPicker),
			//inlineMenuBar: true,
		});
	}
}

class APLListItemPicker extends Component {
	private readonly player: Player<any>;
	private readonly listPicker: ListPicker<Player<any>, APLListItem, APLListItemPicker>;
	private readonly modItem: APLListItem;

	constructor(parent: HTMLElement, player: Player<any>, modItem: APLListItem, listPicker: ListPicker<Player<any>, APLListItem, APLListItemPicker>) {
		super(parent, 'apl-list-item-picker-root');
		this.player = player;
		this.listPicker = listPicker;
		this.modItem = modItem;

        new BooleanPicker(this.rootElem, modItem, {
            label: 'Hide',
            labelTooltip: 'Ignores this APL action.',
            inline: true,
            changedEvent: (item: APLListItem) => player.rotationChangeEmitter,
            getValue: (item: APLListItem) => item.hide,
            setValue: (eventID: EventID, item: APLListItem, newValue: boolean) => {
                item.hide = newValue;
                this.setValue(eventID, item);
            },
        });

        new StringPicker(this.rootElem, modItem, {
            label: 'Notes',
            labelTooltip: 'Description for this action. The sim will ignore this value, it\'s just to allow self-documentation.',
            inline: true,
            changedEvent: (item: APLListItem) => player.rotationChangeEmitter,
            getValue: (item: APLListItem) => item.notes,
            setValue: (eventID: EventID, item: APLListItem, newValue: string) => {
                item.notes = newValue;
                this.setValue(eventID, item);
            },
        });

		//new IconEnumPicker<CustomSpell, number>(this.rootElem, modSpell, {
		//	numColumns: config.numColumns,
		//	values: config.values.map(value => {
		//		if (value.showWhen) {
		//			const oldShowWhen = value.showWhen;
		//			value.showWhen = ((spell: CustomSpell) => oldShowWhen(player)) as unknown as ((player: Player<SpecType>) => boolean);
		//		}
		//		return value;
		//	}) as unknown as Array<IconEnumValueConfig<CustomSpell, number>>,
		//	equals: (a: number, b: number) => a == b,
		//	zeroValue: 0,
		//	changedEvent: (spell: CustomSpell) => player.changeEmitter,
		//	getValue: (spell: CustomSpell) => spell.spell,
		//	setValue: (eventID: EventID, spell: CustomSpell, newValue: number) => {
		//		spell.spell = newValue;
		//		this.setSpell(eventID, spell);
		//	},
		//});
	}

    private getListIndex(): number {
        return this.listPicker.getPickerIndex(this);
    }

    private getValue(): APLListItem|null {
        return this.player.getAplRotation().priorityList[this.getListIndex()] || null;
    }

	private setValue(eventID: EventID, listItem: APLListItem) {
		const index = this.getListIndex();
		const rotation = this.player.getAplRotation();
		rotation.priorityList[index] = APLListItem.clone(listItem);
		this.player.setAplRotation(eventID, rotation);
	}
}
