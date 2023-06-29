import {
	APLAction,
	APLActionCastSpell,
	APLActionSequence,
	APLActionWait,
	APLValue,
} from '../../proto/apl.js';

import { EventID } from '../../typed_event.js';
import { Input, InputConfig } from '../input.js';
import { Player } from '../../player.js';
import { TextDropdownPicker } from '../dropdown_picker.js';

import * as AplHelpers from './apl_helpers.js';
import * as AplValues from './apl_values.js';

export interface APLActionPickerConfig extends InputConfig<Player<any>, APLAction> {
}

export type APLActionType = APLAction['action']['oneofKind'];

export class APLActionPicker extends Input<Player<any>, APLAction> {

	private typePicker: TextDropdownPicker<Player<any>, APLActionType>;

	private readonly actionDiv: HTMLElement;
	private currentType: APLActionType;
	private actionPicker: Input<Player<any>, any>|null;

	private readonly conditionPicker: Input<Player<any>, APLValue>;

	constructor(parent: HTMLElement, player: Player<any>, config: APLActionPickerConfig) {
		super(parent, 'apl-action-picker-root', player, config);

		this.conditionPicker = new AplValues.APLValuePicker(this.rootElem, this.modObject, {
			changedEvent: (player: Player<any>) => player.rotationChangeEmitter,
			getValue: (player: Player<any>) => {
				const action = this.getSourceValue();
				if (!action.condition) {
					action.condition = APLValue.create();
				}
				return action.condition;
			},
			setValue: (eventID: EventID, player: Player<any>, newValue: APLValue) => {
				this.getSourceValue().condition = newValue;
				player.rotationChangeEmitter.emit(eventID);
			},
		});

		this.actionDiv = document.createElement('div');
		this.actionDiv.classList.add('apl-action-picker-action');
		this.rootElem.appendChild(this.actionDiv);

		const allActionTypes = Object.keys(actionTypeFactories) as Array<NonNullable<APLActionType>>;
		this.typePicker = new TextDropdownPicker(this.actionDiv, player, {
            defaultLabel: 'Action',
			values: allActionTypes.map(actionType => {
				return {
					value: actionType,
					label: actionTypeFactories[actionType].label,
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
					const factory = actionTypeFactories[newValue];
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
	}

	getInputElem(): HTMLElement | null {
		return this.rootElem;
	}

    getInputValue(): APLAction {
		const actionType = this.typePicker.getInputValue();
        return APLAction.create({
			condition: this.conditionPicker.getInputValue(),
			action: {
				oneofKind: actionType,
				...((() => {
					const val: any = {};
					if (actionType && this.actionPicker) {
						val[actionType] = this.actionPicker.getInputValue();
					}
					return val;
				})()),
			},
		})
    }

	setInputValue(newValue: APLAction) {
		if (!newValue) {
			return;
		}

		this.conditionPicker.setInputValue(newValue.condition || APLValue.create());

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

		const factory = actionTypeFactories[newActionType];
		this.actionPicker = factory.factory(this.actionDiv, this.modObject, {
			changedEvent: (player: Player<any>) => player.rotationChangeEmitter,
			getValue: () => (this.getSourceValue().action as any)[newActionType] || factory.newValue(),
			setValue: (eventID: EventID, player: Player<any>, newValue: any) => {
				(this.getSourceValue().action as any)[newActionType] = newValue;
				player.rotationChangeEmitter.emit(eventID);
			},
		});
	}
}

type ActionTypeConfig = {
	label: string,
	newValue: () => object,
	factory: (parent: HTMLElement, player: Player<any>, config: InputConfig<Player<any>, any>) => Input<Player<any>, any>,
};

function inputBuilder<T extends object>(label: string, newValue: () => T, fields: Array<AplHelpers.APLPickerBuilderFieldConfig<T, any>>): ActionTypeConfig {
	return {
		label: label,
		newValue: newValue,
		factory: AplHelpers.aplInputBuilder(newValue, fields),
	};
}

export const actionTypeFactories: Record<NonNullable<APLActionType>, ActionTypeConfig> = {
	['castSpell']: inputBuilder('Cast', APLActionCastSpell.create, [
		AplHelpers.actionIdFieldConfig('spellId', 'all_spells'),
	]),
	['sequence']: inputBuilder('Sequence', APLActionSequence.create, [
	]),
	['wait']: inputBuilder('Wait', APLActionWait.create, [
	]),
};