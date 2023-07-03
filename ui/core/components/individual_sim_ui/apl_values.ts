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
import { TextDropdownPicker, TextDropdownValueConfig } from '../dropdown_picker.js';
import { ListItemPickerConfig, ListPicker } from '../list_picker.js';

import * as AplHelpers from './apl_helpers.js';

export interface APLValuePickerConfig extends InputConfig<Player<any>, APLValue|undefined> {
}

export type APLValueType = APLValue['value']['oneofKind'];

export class APLValuePicker extends Input<Player<any>, APLValue|undefined> {

	private typePicker: TextDropdownPicker<Player<any>, APLValueType>;

	private currentType: APLValueType;
	private valuePicker: Input<Player<any>, any>|null;

	constructor(parent: HTMLElement, player: Player<any>, config: APLValuePickerConfig) {
		super(parent, 'apl-value-picker-root', player, config);

		const allValueTypes = Object.keys(valueTypeFactories) as Array<NonNullable<APLValueType>>;
		this.typePicker = new TextDropdownPicker(this.rootElem, player, {
            defaultLabel: 'No Condition',
			values: [{
				value: undefined,
				label: 'None',
			} as TextDropdownValueConfig<APLValueType>].concat(allValueTypes.map(valueType => {
				return {
					value: valueType,
					label: valueTypeFactories[valueType].label,
				};
			})),
			equals: (a, b) => a == b,
			changedEvent: (player: Player<any>) => player.rotationChangeEmitter,
			getValue: (player: Player<any>) => this.getSourceValue()?.value.oneofKind,
			setValue: (eventID: EventID, player: Player<any>, newValue: APLValueType) => {
				const sourceValue = this.getSourceValue();
				if (sourceValue?.value.oneofKind == newValue) {
					return;
				}

				if (newValue) {
					const factory = valueTypeFactories[newValue];
					const obj: any = { oneofKind: newValue };
					obj[newValue] = factory.newValue();
					if (sourceValue) {
						sourceValue.value = obj;
					} else {
						const newSourceValue = APLValue.create();
						newSourceValue.value = obj;
						this.setSourceValue(eventID, newSourceValue);
					}
				} else {
					this.setSourceValue(eventID, newValue);
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

    getInputValue(): APLValue|undefined {
		const valueType = this.typePicker.getInputValue();
		if (!valueType) {
			return undefined;
		} else {
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
    }

	setInputValue(newValue: APLValue|undefined) {
		const newValueType = newValue?.value.oneofKind;
		this.updateValuePicker(newValueType);

		if (newValueType && newValue) {
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
			getValue: () => {
				const sourceVal = this.getSourceValue();
				return sourceVal ? (sourceVal.value as any)[newValueType] || factory.newValue() : factory.newValue();
			},
			setValue: (eventID: EventID, player: Player<any>, newValue: any) => {
				const sourceVal = this.getSourceValue();
				if (sourceVal) {
					(sourceVal.value as any)[newValueType] = newValue;
				}
				player.rotationChangeEmitter.emit(eventID);
			},
		});
	}
}

type ValueTypeConfig<T> = {
	label: string,
	newValue: () => T,
	factory: (parent: HTMLElement, player: Player<any>, config: InputConfig<Player<any>, T>) => Input<Player<any>, T>,
};

export function valueFieldConfig(field: string): AplHelpers.APLPickerBuilderFieldConfig<any, any> {
	return {
		field: field,
		newValue: APLValue.create,
		factory: (parent, player, config) => new APLValuePicker(parent, player, config),
	};
}

export function valueListFieldConfig(field: string): AplHelpers.APLPickerBuilderFieldConfig<any, any> {
	return {
		field: field,
		newValue: () => [],
		factory: (parent, player, config) => new ListPicker<Player<any>, APLValue|undefined>(parent, player, {
			...config,
			// Override setValue to replace undefined elements with default messages.
			setValue: (eventID: EventID, player: Player<any>, newValue: Array<APLValue|undefined>) => {
				config.setValue(eventID, player, newValue.map(val => val || APLValue.create()));
			},

			itemLabel: 'Value',
			newItem: APLValue.create,
			copyItem: (oldValue: APLValue|undefined) => oldValue ? APLValue.clone(oldValue) : oldValue,
			newItemPicker: (parent: HTMLElement, listPicker: ListPicker<Player<any>, APLValue|undefined>, index: number, config: ListItemPickerConfig<Player<any>, APLValue|undefined>) => new APLValuePicker(parent, player, config),
			horizontalLayout: true,
			allowedActions: ['create', 'delete'],
		}),
	};
}

function inputBuilder<T>(label: string, newValue: () => T, fields: Array<AplHelpers.APLPickerBuilderFieldConfig<T, any>>): ValueTypeConfig<T> {
	return {
		label: label,
		newValue: newValue,
		factory: AplHelpers.aplInputBuilder(newValue, fields),
	};
}

const valueTypeFactories: Record<NonNullable<APLValueType>, ValueTypeConfig<any>>  = {
	['const']: inputBuilder('Const', APLValueConst.create, [
		AplHelpers.stringFieldConfig('val'),
	]),
	['and']: inputBuilder('All of', APLValueAnd.create, [
		valueListFieldConfig('vals'),
	]),
	['or']: inputBuilder('Any of', APLValueOr.create, [
		valueListFieldConfig('vals'),
	]),
	['not']: inputBuilder('Not', APLValueNot.create, [
		valueFieldConfig('val'),
	]),
	['cmp']: inputBuilder('Compare', APLValueCompare.create, [
		valueFieldConfig('lhs'),
		{
			field: 'op',
			newValue: () => ComparisonOperator.OpEq,
			factory: (parent, player, config) => new TextDropdownPicker(parent, player, {
				...config,
				defaultLabel: 'None',
				equals: (a, b) => a == b,
				values: [
					{ value: ComparisonOperator.OpEq, label: '==' },
					{ value: ComparisonOperator.OpNe, label: '!=' },
					{ value: ComparisonOperator.OpGe, label: '>=' },
					{ value: ComparisonOperator.OpGt, label: '>' },
					{ value: ComparisonOperator.OpLe, label: '<=' },
					{ value: ComparisonOperator.OpLt, label: '<' },
				],
			}),
		},
		valueFieldConfig('rhs'),
	]),
	['dotIsActive']: inputBuilder('Dot Is Active', APLValueDotIsActive.create, [
		AplHelpers.actionIdFieldConfig('spellId', 'dot_spells'),
	]),
};