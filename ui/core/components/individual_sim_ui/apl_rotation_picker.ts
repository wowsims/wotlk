import { EventID } from '../../typed_event.js';
import { Player } from '../../player.js';
import { BooleanPicker } from '../boolean_picker.js';
import { ListItemPickerConfig, ListPicker } from '../list_picker.js';
import {
	APLListItem,
	APLAction,
	APLRotation,
	APLValue,
} from '../../proto/apl.js';

import { Component } from '../component.js';
import { Input, InputConfig } from '../input.js';
import { SimUI } from '../../sim_ui.js';

import { APLActionPicker } from './apl_actions.js';

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
			newItem: () => APLListItem.create({
				action: {},
			}),
			copyItem: (oldItem: APLListItem) => APLListItem.clone(oldItem),
			newItemPicker: (parent: HTMLElement, listPicker: ListPicker<Player<any>, APLListItem>, index: number, config: ListItemPickerConfig<Player<any>, APLListItem>) => new APLListItemPicker(parent, modPlayer, index, config),
			inlineMenuBar: true,
		});

		modPlayer.rotationChangeEmitter.on(() => console.log('APL: ' + APLRotation.toJsonString(modPlayer.aplRotation)))
	}
}

class APLListItemPicker extends Input<Player<any>, APLListItem> {
	private readonly player: Player<any>;
	private readonly itemIndex: number;

	private readonly hidePicker: Input<null, boolean>;
	private readonly actionPicker: APLActionPicker;

    private getItem(): APLListItem {
        return this.getSourceValue() || APLListItem.create({
			action: {},
		});
    }

	constructor(parent: HTMLElement, player: Player<any>, itemIndex: number, config: ListItemPickerConfig<Player<any>, APLListItem>) {
		super(parent, 'apl-list-item-picker-root', player, config);
		this.player = player;
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

        this.actionPicker = new APLActionPicker(this.rootElem, this.player, {
            changedEvent: () => this.player.rotationChangeEmitter,
            getValue: () => this.getItem().action!,
            setValue: (eventID: EventID, player: Player<any>, newValue: APLAction) => {
                this.getItem().action = newValue;
				this.player.rotationChangeEmitter.emit(eventID);
            },
        });
		this.init();
	}

	getInputElem(): HTMLElement | null {
		return this.rootElem;
	}

    getInputValue(): APLListItem {
        const item = APLListItem.create({
			hide: this.hidePicker.getInputValue(),
			action: this.actionPicker.getInputValue(),
		});
		return item;
    }

	setInputValue(newValue: APLListItem) {
		if (!newValue) {
			return;
		}
		this.hidePicker.setInputValue(newValue.hide);
		this.actionPicker.setInputValue(newValue.action || APLAction.create());
	}
}
