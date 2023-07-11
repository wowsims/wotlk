import {
	APLValue,
	APLValueAnd,
	APLValueOr,
	APLValueNot,
	APLValueCompare,
	APLValueCompare_ComparisonOperator as ComparisonOperator,
	APLValueConst,
	APLValueCurrentTime,
	APLValueCurrentTimePercent,
	APLValueRemainingTime,
	APLValueRemainingTimePercent,
	APLValueCurrentMana,
	APLValueCurrentManaPercent,
	APLValueCurrentRage,
	APLValueCurrentEnergy,
	APLValueCurrentComboPoints,
	APLValueCurrentRunicPower,
	APLValueGCDIsReady,
	APLValueGCDTimeToReady,
	APLValueSpellCanCast,
	APLValueSpellIsReady,
	APLValueSpellTimeToReady,
	APLValueAuraIsActive,
	APLValueAuraRemainingTime,
	APLValueAuraNumStacks,
	APLValueDotIsActive,
	APLValueDotRemainingTime,
} from '../../proto/apl.js';

import { EventID, TypedEvent } from '../../typed_event.js';
import { Input, InputConfig } from '../input.js';
import { Player } from '../../player.js';
import { TextDropdownPicker, TextDropdownValueConfig } from '../dropdown_picker.js';
import { ListItemPickerConfig, ListPicker } from '../list_picker.js';

import * as AplHelpers from './apl_helpers.js';

export interface APLValuePickerConfig extends InputConfig<Player<any>, APLValue | undefined> {
}

export type APLValueType = APLValue['value']['oneofKind'];

export class APLValuePicker extends Input<Player<any>, APLValue | undefined> {

	private typePicker: TextDropdownPicker<Player<any>, APLValueType>;

	private currentType: APLValueType;
	private valuePicker: Input<Player<any>, any> | null;

	constructor(parent: HTMLElement, player: Player<any>, config: APLValuePickerConfig) {
		super(parent, 'apl-value-picker-root', player, config);

		const allValueTypes = Object.keys(valueTypeFactories) as Array<NonNullable<APLValueType>>;
		this.typePicker = new TextDropdownPicker(this.rootElem, player, {
			defaultLabel: 'No Condition',
			values: [{
				value: undefined,
				label: '<None>',
			} as TextDropdownValueConfig<APLValueType>].concat(allValueTypes.map(valueType => {
				const factory = valueTypeFactories[valueType];
				return {
					value: valueType,
					label: factory.label,
					submenu: factory.submenu,
					tooltip: factory.fullDescription ? `<p>${factory.shortDescription}</p> ${factory.fullDescription}` : factory.shortDescription,
				};
			})),
			equals: (a, b) => a == b,
			changedEvent: (player: Player<any>) => player.rotationChangeEmitter,
			getValue: (player: Player<any>) => this.getSourceValue()?.value.oneofKind,
			setValue: (eventID: EventID, player: Player<any>, newValue: APLValueType) => {
				const sourceValue = this.getSourceValue();
				const oldValue = sourceValue?.value.oneofKind;
				if (oldValue == newValue) {
					return;
				}

				if (newValue) {
					const factory = valueTypeFactories[newValue];
					let obj: any = { oneofKind: newValue };
					obj[newValue] = factory.newValue();
					if (sourceValue) {
						// Some pre-fill logic when swapping types.
						if (oldValue && this.valuePicker) {
							if (newValue == 'not') {
								(obj[newValue] as APLValueNot).val = this.makeAPLValue(oldValue, this.valuePicker.getInputValue());
							} else if (sourceValue.value.oneofKind == 'not' && sourceValue.value.not.val?.value.oneofKind == newValue && !['and', 'or'].includes(newValue)) {
								obj = sourceValue.value.not.val.value;
							} else if (newValue == 'and') {
								if (sourceValue.value.oneofKind == 'or') {
									(obj[newValue] as APLValueAnd).vals = sourceValue.value.or.vals;
								} else {
									(obj[newValue] as APLValueAnd).vals = [this.makeAPLValue(oldValue, this.valuePicker.getInputValue())];
								}
							} else if (newValue == 'or') {
								if (sourceValue.value.oneofKind == 'and') {
									(obj[newValue] as APLValueOr).vals = sourceValue.value.and.vals;
								} else {
									(obj[newValue] as APLValueOr).vals = [this.makeAPLValue(oldValue, this.valuePicker.getInputValue())];
								}
							} else if (sourceValue.value.oneofKind == 'and' && sourceValue.value.and.vals?.[0]?.value.oneofKind == newValue) {
								obj = sourceValue.value.and.vals[0].value;
							} else if (sourceValue.value.oneofKind == 'or' && sourceValue.value.or.vals?.[0]?.value.oneofKind == newValue) {
								obj = sourceValue.value.or.vals[0].value;
							}
						}

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

	getInputValue(): APLValue | undefined {
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

	setInputValue(newValue: APLValue | undefined) {
		const newValueType = newValue?.value.oneofKind;
		this.updateValuePicker(newValueType);

		if (newValueType && newValue) {
			this.valuePicker!.setInputValue((newValue.value as any)[newValueType]);
		}
	}

	private makeAPLValue(type: APLValueType, implVal: any): APLValue {
		if (!type) {
			return APLValue.create();
		}
		const obj: any = { oneofKind: type };
		obj[type] = implVal;
		return APLValue.create({value: obj});
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
	submenu?: Array<string>,
	shortDescription: string,
	fullDescription?: string,
	newValue: () => T,
	factory: (parent: HTMLElement, player: Player<any>, config: InputConfig<Player<any>, T>) => Input<Player<any>, T>,
};

function comparisonOperatorFieldConfig(field: string): AplHelpers.APLPickerBuilderFieldConfig<any, any> {
	return {
		field: field,
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
	};
}

export function valueFieldConfig(field: string, options?: Partial<AplHelpers.APLPickerBuilderFieldConfig<any, any>>): AplHelpers.APLPickerBuilderFieldConfig<any, any> {
	return {
		field: field,
		newValue: APLValue.create,
		factory: (parent, player, config) => new APLValuePicker(parent, player, config),
		...(options || {}),
	};
}

export function valueListFieldConfig(field: string): AplHelpers.APLPickerBuilderFieldConfig<any, any> {
	return {
		field: field,
		newValue: () => [],
		factory: (parent, player, config) => new ListPicker<Player<any>, APLValue | undefined>(parent, player, {
			...config,
			// Override setValue to replace undefined elements with default messages.
			setValue: (eventID: EventID, player: Player<any>, newValue: Array<APLValue | undefined>) => {
				config.setValue(eventID, player, newValue.map(val => val || APLValue.create()));
			},

			itemLabel: 'Value',
			newItem: APLValue.create,
			copyItem: (oldValue: APLValue | undefined) => oldValue ? APLValue.clone(oldValue) : oldValue,
			newItemPicker: (parent: HTMLElement, listPicker: ListPicker<Player<any>, APLValue | undefined>, index: number, config: ListItemPickerConfig<Player<any>, APLValue | undefined>) => new APLValuePicker(parent, player, config),
			horizontalLayout: true,
			allowedActions: ['create', 'delete'],
		}),
	};
}

function inputBuilder<T>(config: {
	label: string,
	submenu?: Array<string>,
	shortDescription: string,
	fullDescription?: string,
	newValue: () => T,
	fields: Array<AplHelpers.APLPickerBuilderFieldConfig<T, any>>,
}): ValueTypeConfig<T> {
	return {
		label: config.label,
		submenu: config.submenu,
		shortDescription: config.shortDescription,
		fullDescription: config.fullDescription,
		newValue: config.newValue,
		factory: AplHelpers.aplInputBuilder(config.newValue, config.fields),
	};
}

const valueTypeFactories: Record<NonNullable<APLValueType>, ValueTypeConfig<any>> = {
	// Operators
	['const']: inputBuilder({
		label: 'Const',
		shortDescription: 'A fixed value.',
		fullDescription: `
		<p>
			Examples:
			<ul>
				<li><b>Number:</b> '123', '0.5', '-10'</li>
				<li><b>Time:</b> '100ms', '5s', '3m'</li>
				<li><b>Percentage:</b> '30%'</li>
			</ul>
		</p>
		`,
		newValue: APLValueConst.create,
		fields: [
			AplHelpers.stringFieldConfig('val'),
		],
	}),
	['cmp']: inputBuilder({
		label: 'Compare',
		submenu: ['Logic'],
		shortDescription: 'Compares two values.',
		newValue: APLValueCompare.create,
		fields: [
			valueFieldConfig('lhs'),
			comparisonOperatorFieldConfig('op'),
			valueFieldConfig('rhs'),
		],
	}),
	['and']: inputBuilder({
		label: 'All of',
		submenu: ['Logic'],
		shortDescription: 'Returns <b>True</b> if all of the sub-values are <b>True</b>, otherwise <b>False</b>',
		newValue: APLValueAnd.create,
		fields: [
			valueListFieldConfig('vals'),
		],
	}),
	['or']: inputBuilder({
		label: 'Any of',
		submenu: ['Logic'],
		shortDescription: 'Returns <b>True</b> if any of the sub-values are <b>True</b>, otherwise <b>False</b>',
		newValue: APLValueOr.create,
		fields: [
			valueListFieldConfig('vals'),
		],
	}),
	['not']: inputBuilder({
		label: 'Not',
		submenu: ['Logic'],
		shortDescription: 'Returns the opposite of the inner value, i.e. <b>True</b> if the value is <b>False</b> and vice-versa.',
		newValue: APLValueNot.create,
		fields: [
			valueFieldConfig('val'),
		],
	}),

	// Encounter
	['currentTime']: inputBuilder({
		label: 'Current Time',
		submenu: ['Encounter'],
		shortDescription: 'Elapsed time of the current sim iteration.',
		newValue: APLValueCurrentTime.create,
		fields: [],
	}),
	['currentTimePercent']: inputBuilder({
		label: 'Current Time (%)',
		submenu: ['Encounter'],
		shortDescription: 'Elapsed time of the current sim iteration, as a percentage.',
		newValue: APLValueCurrentTimePercent.create,
		fields: [],
	}),
	['remainingTime']: inputBuilder({
		label: 'Remaining Time',
		submenu: ['Encounter'],
		shortDescription: 'Elapsed time of the remaining sim iteration.',
		newValue: APLValueRemainingTime.create,
		fields: [],
	}),
	['remainingTimePercent']: inputBuilder({
		label: 'Remaining Time (%)',
		submenu: ['Encounter'],
		shortDescription: 'Elapsed time of the remaining sim iteration, as a percentage.',
		newValue: APLValueRemainingTimePercent.create,
		fields: [],
	}),

	// Resources
	['currentMana']: inputBuilder({
		label: 'Mana',
		submenu: ['Resources'],
		shortDescription: 'Amount of currently available Mana.',
		newValue: APLValueCurrentMana.create,
		fields: [],
	}),
	['currentManaPercent']: inputBuilder({
		label: 'Mana (%)',
		submenu: ['Resources'],
		shortDescription: 'Amount of currently available Mana, as a percentage.',
		newValue: APLValueCurrentManaPercent.create,
		fields: [],
	}),
	['currentRage']: inputBuilder({
		label: 'Rage',
		submenu: ['Resources'],
		shortDescription: 'Amount of currently available Rage.',
		newValue: APLValueCurrentRage.create,
		fields: [],
	}),
	['currentEnergy']: inputBuilder({
		label: 'Energy',
		submenu: ['Resources'],
		shortDescription: 'Amount of currently available Energy.',
		newValue: APLValueCurrentEnergy.create,
		fields: [],
	}),
	['currentComboPoints']: inputBuilder({
		label: 'Combo Points',
		submenu: ['Resources'],
		shortDescription: 'Amount of currently available Combo Points.',
		newValue: APLValueCurrentComboPoints.create,
		fields: [],
	}),
	['currentRunicPower']: inputBuilder({
		label: 'Runic Power',
		submenu: ['Resources'],
		shortDescription: 'Amount of currently available Runic Power.',
		newValue: APLValueCurrentRunicPower.create,
		fields: [],
	}),

	// GCD
	['gcdIsReady']: inputBuilder({
		label: 'GCD Is Ready',
		submenu: ['GCD'],
		shortDescription: '<b>True</b> if the GCD is not on cooldown, otherwise <b>False</b>.',
		newValue: APLValueGCDIsReady.create,
		fields: [],
	}),
	['gcdTimeToReady']: inputBuilder({
		label: 'GCD Time To Ready',
		submenu: ['GCD'],
		shortDescription: 'Amount of time remaining before the GCD comes off cooldown, or <b>0</b> if it is not on cooldown.',
		newValue: APLValueGCDTimeToReady.create,
		fields: [],
	}),

	// Spells
	['spellCanCast']: inputBuilder({
		label: 'Can Cast',
		submenu: ['Spell'],
		shortDescription: '<b>True</b> if all requirements for casting the spell are currently met, otherwise <b>False</b>.',
		fullDescription: `
			<p>The <b>Cast Spell</b> action does not need to be conditioned on this, because it applies this check automatically.</p>
		`,
		newValue: APLValueSpellCanCast.create,
		fields: [
			AplHelpers.actionIdFieldConfig('spellId', 'castable_spells'),
		],
	}),
	['spellIsReady']: inputBuilder({
		label: 'Is Ready',
		submenu: ['Spell'],
		shortDescription: '<b>True</b> if the spell is not on cooldown, otherwise <b>False</b>.',
		newValue: APLValueSpellIsReady.create,
		fields: [
			AplHelpers.actionIdFieldConfig('spellId', 'castable_spells'),
		],
	}),
	['spellTimeToReady']: inputBuilder({
		label: 'Time To Ready',
		submenu: ['Spell'],
		shortDescription: 'Amount of time remaining before the spell comes off cooldown, or <b>0</b> if it is not on cooldown.',
		newValue: APLValueSpellTimeToReady.create,
		fields: [
			AplHelpers.actionIdFieldConfig('spellId', 'castable_spells'),
		],
	}),

	// Auras
	['auraIsActive']: inputBuilder({
		label: 'Aura Is Active',
		submenu: ['Aura'],
		shortDescription: '<b>True</b> if the aura is currently active on self, otherwise <b>False</b>.',
		newValue: APLValueAuraIsActive.create,
		fields: [
			AplHelpers.actionIdFieldConfig('auraId', 'auras'),
		],
	}),
	['auraRemainingTime']: inputBuilder({
		label: 'Aura Remaining Time',
		submenu: ['Aura'],
		shortDescription: 'Time remaining before this aura will expire, or 0 if the aura is not currently active on self.',
		newValue: APLValueAuraRemainingTime.create,
		fields: [
			AplHelpers.actionIdFieldConfig('auraId', 'auras'),
		],
	}),
	['auraNumStacks']: inputBuilder({
		label: 'Aura Num Stacks',
		submenu: ['Aura'],
		shortDescription: 'Number of stacks of the aura on self.',
		newValue: APLValueAuraNumStacks.create,
		fields: [
			AplHelpers.actionIdFieldConfig('auraId', 'stackable_auras'),
		],
	}),

	// DoT
	['dotIsActive']: inputBuilder({
		label: 'Dot Is Active',
		submenu: ['DoT'],
		shortDescription: '<b>True</b> if the specified dot is currently ticking, otherwise <b>False</b>.',
		newValue: APLValueDotIsActive.create,
		fields: [
			AplHelpers.actionIdFieldConfig('spellId', 'dot_spells'),
		],
	}),
	['dotRemainingTime']: inputBuilder({
		label: 'Dot Remaining Time',
		submenu: ['DoT'],
		shortDescription: 'Time remaining before the last tick of this DoT will occur, or 0 if the DoT is not currently ticking.',
		newValue: APLValueDotRemainingTime.create,
		fields: [
			AplHelpers.actionIdFieldConfig('spellId', 'dot_spells'),
		],
	}),
};
