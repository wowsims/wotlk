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

import { EventID, TypedEvent } from '../../typed_event.js';
import { Input, InputConfig } from '../input.js';
import { Player } from '../../player.js';
import { TextDropdownPicker } from '../dropdown_picker.js';

import * as AplHelpers from './apl_helpers.js';

export class APLValueConstValuePicker extends Input<Player<any>, APLValueConst> {
	private readonly inputElem: HTMLInputElement;

	constructor(parent: HTMLElement, modObject: Player<any>, config: InputConfig<Player<any>, APLValueConst>) {
		super(parent, 'apl-value-const-picker-root', modObject, config);

		this.inputElem = document.createElement('input');
		this.inputElem.type = 'text';
		this.rootElem.appendChild(this.inputElem);

		this.init();

		this.inputElem.addEventListener('change', event => {
			this.inputChanged(TypedEvent.nextEventID());
		});
	}

	getInputElem(): HTMLElement {
		return this.inputElem;
	}

	getInputValue(): APLValueConst {
		return APLValueConst.create({ val: this.inputElem.value });
	}

	setInputValue(newValue: APLValueConst) {
		this.inputElem.value = newValue ? newValue.val : '';
	}
}

export interface APLValuePickerConfig extends InputConfig<Player<any>, APLValue> {
}

export type APLValueType = APLValue['value']['oneofKind'];

export class APLValuePicker extends Input<Player<any>, APLValue> {

	private typePicker: TextDropdownPicker<Player<any>, APLValueType>;

	private currentType: APLValueType;
	private valuePicker: Input<Player<any>, any>|null;

	constructor(parent: HTMLElement, player: Player<any>, config: APLValuePickerConfig) {
		super(parent, 'apl-value-picker-root', player, config);

		const allValueTypes = Object.keys(valueTypeFactories) as Array<NonNullable<APLValueType>>;
		this.typePicker = new TextDropdownPicker(this.rootElem, player, {
            defaultLabel: 'No Condition',
			values: allValueTypes.map(valueType => {
				return {
					value: valueType,
					label: valueTypeFactories[valueType].label,
				};
			}),
			equals: (a, b) => a == b,
			changedEvent: (player: Player<any>) => player.rotationChangeEmitter,
			getValue: (player: Player<any>) => this.getSourceValue().value.oneofKind,
			setValue: (eventID: EventID, player: Player<any>, newValue: APLValueType) => {
				const value = this.getSourceValue();
				if (value.value.oneofKind == newValue) {
					return;
				}
				if (newValue) {
					const factory = valueTypeFactories[newValue];
					const obj: any = { oneofKind: newValue };
					obj[newValue] = factory.newValue();
					value.value = obj;
				} else {
					value.value = {
						oneofKind: newValue,
					};
				}
				player.rotationChangeEmitter.emit(eventID);
			},
		});

		this.currentType = undefined;
		this.valuePicker = null;

		this.init();
	}

	getInputElem(): HTMLElement | null {
		return this.rootElem;
	}

    getInputValue(): APLValue {
		const valueType = this.typePicker.getInputValue();
        return APLValue.create({
			value: {
				oneofKind: valueType,
				...((() => {
					const val: any = {};
					if (valueType && this.valuePicker) {
						val[valueType] = this.valuePicker.getInputValue();
					}
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
			this.valuePicker!.setInputValue((newValue.value as any)[newValueType]);
		}
	}

	private updateValuePicker(newValueType: APLValueType) {
		const valueType = this.currentType;
		if (newValueType == valueType) {
			return;
		}
		this.currentType = newValueType;

		if (this.valuePicker) {
			this.valuePicker.rootElem.remove();
			this.valuePicker = null;
		}

		if (!newValueType) {
			return;
		}

		this.typePicker.setInputValue(newValueType);

		const factory = valueTypeFactories[newValueType];
		this.valuePicker = factory.factory(this.rootElem, this.modObject, {
			changedEvent: (player: Player<any>) => player.rotationChangeEmitter,
			getValue: () => (this.getSourceValue().value as any)[newValueType] || factory.newValue(),
			setValue: (eventID: EventID, player: Player<any>, newValue: any) => {
				(this.getSourceValue().value as any)[newValueType] = newValue;
				player.rotationChangeEmitter.emit(eventID);
			},
		});
	}
}

type ValueTypeConfig = {
	label: string,
	newValue: () => object,
	factory: (parent: HTMLElement, player: Player<any>, config: InputConfig<Player<any>, any>) => Input<Player<any>, any>,
};

function inputBuilder<T extends object>(label: string, newValue: () => T, fields: Array<AplHelpers.APLPickerBuilderFieldConfig<T, any>>): ValueTypeConfig {
	return {
		label: label,
		newValue: newValue,
		factory: AplHelpers.aplInputBuilder(newValue, fields),
	};
}

const valueTypeFactories: Record<NonNullable<APLValueType>, ValueTypeConfig>  = {
	['const']: inputBuilder('Const', APLValueConst.create, [
		{
			field: 'value',
			newValue: () => '',
			factory: (parent, player, config) => new APLValueConstValuePicker(parent, player, config),
		},
	]),
	['and']: inputBuilder('All of', APLValueAnd.create, [
	]),
	['or']: inputBuilder('Any of', APLValueOr.create, [
	]),
	['not']: inputBuilder('Not', APLValueNot.create, [
	]),
	['cmp']: inputBuilder('Compare', APLValueCompare.create, [
	]),
	['dotIsActive']: inputBuilder('Dot Is Active', APLValueDotIsActive.create, [
		AplHelpers.actionIdFieldConfig('spellId', 'dots'),
	]),
};