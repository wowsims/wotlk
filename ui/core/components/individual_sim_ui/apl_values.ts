import {
	APLValue,
	APLValueAnd,
	APLValueOr,
	APLValueNot,
	APLValueCompare,
	APLValueCompare_ComparisonOperator as ComparisonOperator,
	APLValueConst,
	APLValueDotIsActive,
} from '../../proto/apl.js';

import { ActionID, Spec } from '../../proto/common.js';
import { EventID } from '../../typed_event.js';
import { Input, InputConfig } from '../input.js';
import { ActionId } from '../../proto_utils/action_id.js';
import { Player } from '../../player.js';
import { stringComparator } from '../../utils.js';
import { TextDropdownPicker } from '../dropdown_picker.js';

import * as AplHelpers from './apl_helpers.js';

export interface APLValuePickerConfig extends InputConfig<Player<any>, APLValue> {
}

export type APLValueType = APLValue['value']['oneofKind'];

export class APLValuePicker extends Input<Player<any>, APLValue> {

	private typePicker: TextDropdownPicker<Player<any>, APLValueType>;

	private currentType: APLValueType;
	private actionPicker: Input<Player<any>, any>|null;

	constructor(parent: HTMLElement, player: Player<any>, config: APLValuePickerConfig) {
		super(parent, 'apl-action-picker-root', player, config);

		const allValueTypes = Object.keys(APLValuePicker.valueTypeFactories) as Array<NonNullable<APLValueType>>;
		this.typePicker = new TextDropdownPicker(this.rootElem, player, {
            defaultLabel: 'No Condition',
			values: allValueTypes.map(actionType => {
				return {
					value: actionType,
					label: APLValuePicker.valueTypeFactories[actionType].label,
				};
			}),
			equals: (a, b) => a == b,
			changedEvent: (player: Player<any>) => player.rotationChangeEmitter,
			getValue: (player: Player<any>) => this.getSourceValue().value.oneofKind,
			setValue: (eventID: EventID, player: Player<any>, newValue: APLValueType) => {
				const action = this.getSourceValue();
				if (action.value.oneofKind == newValue) {
					return;
				}
				if (newValue) {
					const factory = APLValuePicker.valueTypeFactories[newValue];
					const obj: any = { oneofKind: newValue };
					obj[newValue] = factory.newValue();
					action.value = obj;
				} else {
					action.value = {
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

    getInputValue(): APLValue {
		const actionType = this.typePicker.getInputValue();
        return APLValue.create({
			value: {
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

	setInputValue(newValue: APLValue) {
		if (!newValue) {
			return;
		}

		const newValueType = newValue.value.oneofKind;
		this.updateValuePicker(newValueType);

		if (newValueType) {
			this.actionPicker!.setInputValue((newValue.value as any)[newValueType]);
		}
	}

	private updateValuePicker(newValueType: APLValueType) {
		const actionType = this.currentType;
		if (newValueType == actionType) {
			return;
		}
		this.currentType = newValueType;

		if (this.actionPicker) {
			this.actionPicker.rootElem.remove();
			this.actionPicker = null;
		}

		if (!newValueType) {
			return;
		}

		this.typePicker.setInputValue(newValueType);

		const factory = APLValuePicker.valueTypeFactories[newValueType];
		this.actionPicker = new factory.factory(this.rootElem, this.modObject, {
			changedEvent: (player: Player<any>) => player.rotationChangeEmitter,
			getValue: () => (this.getSourceValue().value as any)[newValueType] || factory.newValue(),
			setValue: (eventID: EventID, player: Player<any>, newValue: any) => {
				(this.getSourceValue().value as any)[newValueType] = newValue;
				player.rotationChangeEmitter.emit(eventID);
			},
		});
	}

	private static valueTypeFactories: Record<NonNullable<APLValueType>, { label: string, newValue: () => object, factory: new (parent: HTMLElement, player: Player<any>, config: InputConfig<Player<any>, any>) => Input<Player<any>, any> }>  = {
		['and']: { label: 'All of', newValue: APLValueAnd.create, factory: APLValueCastSpellPicker },
		['or']: { label: 'Any of', newValue: APLValueOr.create, factory: APLValueSequencePicker },
		['not']: { label: 'Not', newValue: APLValueNot.create, factory: APLValueWaitPicker },
	};
}