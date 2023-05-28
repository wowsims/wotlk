import { Spec } from '../../proto/common.js';
import { EventID } from '../../typed_event.js';
import { Player } from '../../player.js';
import { IconEnumPicker, IconEnumValueConfig } from '../icon_enum_picker.js';
import { BooleanPicker } from '../boolean_picker.js';
import { ListItemPickerConfig, ListPicker } from '../list_picker.js';
import { NumberPicker } from '../number_picker.js';
import { StringPicker } from '../string_picker.js';
import { APLListItem, APLAction } from '../../proto/apl.js';

import { Component } from '../component.js';
import { Input } from '../input.js';
import { SimUI } from 'ui/core/sim_ui.js';

export class APLRotationPicker extends Component {
	constructor(parent: HTMLElement, simUI: SimUI, modPlayer: Player<any>) {
		super(parent, 'apl-rotation-picker-root');

		new ListPicker<Player<any>, APLListItem>(this.rootElem, modPlayer, {
			extraCssClasses: ['apl-list-item-picker'],
			title: 'Priority List',
			titleTooltip: 'At each decision point, the simulation will perform the first valid action from this list.',
			itemLabel: 'Action',
			changedEvent: (player: Player<any>) => player.rotationChangeEmitter,
			getValue: (player: Player<any>) => player.aplRotation.priorityList,
			setValue: (eventID: EventID, player: Player<any>, newValue: Array<APLListItem>) => {
                player.aplRotation.priorityList = newValue;
				player.rotationChangeEmitter.emit(eventID);
			},
			newItem: () => APLListItem.create(),
			copyItem: (oldItem: APLListItem) => APLListItem.clone(oldItem),
			newItemPicker: (parent: HTMLElement, listPicker: ListPicker<Player<any>, APLListItem>, index: number, config: ListItemPickerConfig<Player<any>, APLListItem>) => new APLListItemPicker(parent, modPlayer, listPicker, index, config),
			//inlineMenuBar: true,
		});
	}
}

class APLListItemPicker extends Input<Player<any>, APLListItem> {
	private readonly player: Player<any>;
	private readonly listPicker: ListPicker<Player<any>, APLListItem>;
	private readonly itemIndex: number;

	private readonly hidePicker: Input<null, boolean>;
	private readonly notesPicker: Input<null, string>;

    private getItem(): APLListItem {
        return this.player.aplRotation.priorityList[this.itemIndex] || APLListItem.create();
    }

	constructor(parent: HTMLElement, player: Player<any>, listPicker: ListPicker<Player<any>, APLListItem>, itemIndex: number, config: ListItemPickerConfig<Player<any>, APLListItem>) {
		super(parent, 'apl-list-item-picker-root', player, config);
		this.player = player;
		this.listPicker = listPicker;
		this.itemIndex = itemIndex;

        this.hidePicker = new BooleanPicker(this.rootElem, null, {
            label: 'Hide',
            labelTooltip: 'Ignores this APL action.',
            inline: true,
            changedEvent: () => this.player.rotationChangeEmitter,
            getValue: () => this.getItem().hide,
            setValue: (eventID: EventID, _: null, newValue: boolean) => {
                this.getItem().hide = newValue;
				this.player.rotationChangeEmitter.emit(eventID);
            },
        });

        this.notesPicker = new StringPicker(this.rootElem, null, {
            label: 'Notes',
            labelTooltip: 'Description for this action. The sim will ignore this value, it\'s just to allow self-documentation.',
            inline: true,
            changedEvent: () => this.player.rotationChangeEmitter,
            getValue: () => this.getItem().notes,
            setValue: (eventID: EventID, _: null, newValue: string) => {
                this.getItem().notes = newValue;
				this.player.rotationChangeEmitter.emit(eventID);
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

	getInputElem(): HTMLElement | null {
		return this.rootElem;
	}

    getInputValue(): APLListItem {
        return APLListItem.create({
			hide: this.hidePicker.getInputValue(),
			notes: this.notesPicker.getInputValue(),
		})
    }

	setInputValue(newValue: APLListItem) {
		if (!newValue) {
			return;
		}
		this.hidePicker.setInputValue(newValue.hide);
		this.notesPicker.setInputValue(newValue.notes);
	}
}
