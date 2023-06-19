import { ActionID } from '../../proto/common.js';
import { EventID } from '../../typed_event.js';
import { Player } from '../../player.js';
import { BooleanPicker } from '../boolean_picker.js';
import { DropdownPicker, DropdownPickerConfig, TextDropdownPicker } from '../dropdown_picker.js';
import { ListItemPickerConfig, ListPicker } from '../list_picker.js';
import { StringPicker } from '../string_picker.js';
import {
	APLListItem,
	APLAction,
	APLActionCastSpell,
	APLActionSequence,
	APLActionWait,
	APLRotation,
} from '../../proto/apl.js';

import { Component } from '../component.js';
import { Input, InputConfig } from '../input.js';
import { SimUI } from '../../sim_ui.js';

import * as AplActions from './apl_actions.js';

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
	private readonly notesPicker: Input<null, string>;
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
			notes: this.notesPicker.getInputValue(),
			action: this.actionPicker.getInputValue(),
		});
		return item;
    }

	setInputValue(newValue: APLListItem) {
		if (!newValue) {
			return;
		}
		this.hidePicker.setInputValue(newValue.hide);
		this.notesPicker.setInputValue(newValue.notes);
		this.actionPicker.setInputValue(newValue.action || APLAction.create());
	}
}

export interface APLActionPickerConfig extends InputConfig<Player<any>, APLAction> {
}

type APLActionType = APLAction['action']['oneofKind'];

class APLActionPicker extends Input<Player<any>, APLAction> {

	private typePicker: TextDropdownPicker<Player<any>, APLActionType>;

	private currentType: APLActionType;
	private actionPicker: Input<Player<any>, any>|null;

	constructor(parent: HTMLElement, player: Player<any>, config: APLActionPickerConfig) {
		super(parent, 'apl-action-picker-root', player, config);

		const allActionTypes = Object.keys(APLActionPicker.actionTypeFactories) as Array<NonNullable<APLActionType>>;
		this.typePicker = new TextDropdownPicker(this.rootElem, player, {
            defaultLabel: 'Action',
			values: allActionTypes.map(actionType => {
				return {
					value: actionType,
					label: APLActionPicker.actionTypeFactories[actionType].label,
				};
			}),
			equals: (a, b) => a == b,
			changedEvent: (player: Player<any>) => player.rotationChangeEmitter,
			getValue: (player: Player<any>) => this.getSourceValue().action.oneofKind,
			setValue: (eventID: EventID, player: Player<any>, newValue: APLActionType) => {
				const action = this.getSourceValue();
				if (action.action.oneofKind == newValue) {
					return;
				}
				if (newValue) {
					const factory = APLActionPicker.actionTypeFactories[newValue];
					const obj: any = { oneofKind: newValue };
					obj[newValue] = factory.newValue();
					action.action = obj;
				} else {
					action.action = {
						oneofKind: newValue,
					};
				}
				player.rotationChangeEmitter.emit(eventID);
			},
		});

		this.currentType = undefined;
		this.actionPicker = null;

		this.init();
		//player.rotationChangeEmitter.on(() => this.updateActionPicker(this.typePicker.getInputValue()));
	}

	getInputElem(): HTMLElement | null {
		return this.rootElem;
	}

    getInputValue(): APLAction {
		const actionType = this.typePicker.getInputValue();
        return APLAction.create({
			action: {
				oneofKind: actionType,
				...((() => {
					if (!actionType || !this.actionPicker) return;
					const val: any = {};
					val[actionType] = this.actionPicker.getInputValue();
					return val;
				})()),
			},
		})
    }

	setInputValue(newValue: APLAction) {
		if (!newValue) {
			return;
		}

		const newActionType = newValue.action.oneofKind;
		this.updateActionPicker(newActionType);

		if (newActionType) {
			this.actionPicker!.setInputValue((newValue.action as any)[newActionType]);
		}
	}

	private updateActionPicker(newActionType: APLActionType) {
		const actionType = this.currentType;
		if (newActionType == actionType) {
			return;
		}
		this.currentType = newActionType;

		if (this.actionPicker) {
			this.actionPicker.rootElem.remove();
			this.actionPicker = null;
		}

		if (!newActionType) {
			return;
		}

		this.typePicker.setInputValue(newActionType);

		const factory = APLActionPicker.actionTypeFactories[newActionType];
		this.actionPicker = new factory.factory(this.rootElem, this.modObject, {
			changedEvent: (player: Player<any>) => player.rotationChangeEmitter,
			getValue: () => (this.getSourceValue().action as any)[newActionType] || factory.newValue(),
			setValue: (eventID: EventID, player: Player<any>, newValue: any) => {
				(this.getSourceValue().action as any)[newActionType] = newValue;
				player.rotationChangeEmitter.emit(eventID);
			},
		});
	}

	private static actionTypeFactories: Record<NonNullable<APLActionType>, { label: string, newValue: () => object, factory: new (parent: HTMLElement, player: Player<any>, config: InputConfig<Player<any>, any>) => Input<Player<any>, any> }>  = {
		['castSpell']: { label: 'Cast', newValue: APLActionCastSpell.create, factory: AplActions.APLActionCastSpellPicker },
		['sequence']: { label: 'Sequence', newValue: APLActionSequence.create, factory: AplActions.APLActionSequencePicker },
		['wait']: { label: 'Wait', newValue: APLActionWait.create, factory: AplActions.APLActionWaitPicker },
	};
}
