import {
	Class,
	Spec,
} from '../../proto/common.js';

import {
	ShamanTotems_TotemType as TotemType,
} from '../../proto/shaman.js';

import {
	APLValue,
	APLValueAnd,
	APLValueOr,
	APLValueNot,
	APLValueCompare,
	APLValueCompare_ComparisonOperator as ComparisonOperator,
	APLValueMath,
	APLValueMath_MathOperator as MathOperator,
	APLValueMax,
	APLValueMin,
	APLValueConst,
	APLValueCurrentTime,
	APLValueCurrentTimePercent,
	APLValueRemainingTime,
	APLValueRemainingTimePercent,
	APLValueIsExecutePhase,
	APLValueIsExecutePhase_ExecutePhaseThreshold as ExecutePhaseThreshold,
	APLValueCurrentHealth,
	APLValueCurrentHealthPercent,
	APLValueCurrentMana,
	APLValueCurrentManaPercent,
	APLValueCurrentRage,
	APLValueCurrentEnergy,
	APLValueCurrentComboPoints,
	APLValueCurrentRunicPower,
	APLValueCurrentRuneCount,
	APLValueCurrentRuneDeath,
	APLValueCurrentRuneActive,
	APLValueCurrentNonDeathRuneCount,
	APLValueRuneSlotCooldown,
	APLValueRuneGrace,
	APLValueRuneSlotGrace,
	APLValueGCDIsReady,
	APLValueGCDTimeToReady,
	APLValueAutoTimeToNext,
	APLValueSpellCanCast,
	APLValueSpellIsReady,
	APLValueSpellTimeToReady,
	APLValueSpellCastTime,
	APLValueSpellTravelTime,
	APLValueSpellCPM,
	APLValueSpellIsChanneling,
	APLValueSpellChanneledTicks,
	APLValueSpellCurrentCost,
	APLValueChannelClipDelay,
	APLValueFrontOfTarget,
	APLValueAuraIsActive,
	APLValueAuraIsActiveWithReactionTime,
	APLValueAuraRemainingTime,
	APLValueAuraNumStacks,
	APLValueAuraInternalCooldown,
	APLValueAuraICDIsReadyWithReactionTime,
	APLValueAuraShouldRefresh,
	APLValueDotIsActive,
	APLValueDotRemainingTime,
	APLValueSequenceIsComplete,
	APLValueSequenceIsReady,
	APLValueSequenceTimeToReady,
	APLValueRuneCooldown,
	APLValueNextRuneCooldown,
	APLValueNumberTargets,
	APLValueTotemRemainingTime,
	APLValueCatExcessEnergy,
	APLValueWarlockShouldRecastDrainSoul,
	APLValueWarlockShouldRefreshCorruption,
	APLValueCatNewSavageRoarDuration,
	APLValueBossSpellTimeToReady,
	APLValueBossSpellIsCasting,
} from '../../proto/apl.js';

import { EventID } from '../../typed_event.js';
import { Input, InputConfig } from '../input.js';
import { Player } from '../../player.js';
import { TextDropdownPicker, TextDropdownValueConfig } from '../dropdown_picker.js';
import { ListItemPickerConfig, ListPicker } from '../list_picker.js';

import * as AplHelpers from './apl_helpers.js';

export interface APLValuePickerConfig extends InputConfig<Player<any>, APLValue | undefined> {
}

export type APLValueKind = APLValue['value']['oneofKind'];
export type APLValueImplStruct<F extends APLValueKind> = Extract<APLValue['value'], {oneofKind: F}>;
type APLValueImplTypesUnion = {
	[f in NonNullable<APLValueKind>]: f extends keyof APLValueImplStruct<f> ? APLValueImplStruct<f>[f] : never;
};
export type APLValueImplType = APLValueImplTypesUnion[NonNullable<APLValueKind>]|undefined;

export class APLValuePicker extends Input<Player<any>, APLValue | undefined> {

	private kindPicker: TextDropdownPicker<Player<any>, APLValueKind>;

	private currentKind: APLValueKind;
	private valuePicker: Input<Player<any>, any> | null;

	constructor(parent: HTMLElement, player: Player<any>, config: APLValuePickerConfig) {
		super(parent, 'apl-value-picker-root', player, config);

		const isPrepull = this.rootElem.closest('.apl-prepull-action-picker') != null;

		const allValueKinds = (Object.keys(valueKindFactories) as Array<NonNullable<APLValueKind>>)
			.filter(valueKind => valueKindFactories[valueKind].includeIf?.(player, isPrepull) ?? true);

		this.kindPicker = new TextDropdownPicker(this.rootElem, player, {
			defaultLabel: 'No Condition',
			values: [{
				value: undefined,
				label: '<None>',
			} as TextDropdownValueConfig<APLValueKind>].concat(allValueKinds.map(kind => {
				const factory = valueKindFactories[kind];
				return {
					value: kind,
					label: factory.label,
					submenu: factory.submenu,
					tooltip: factory.fullDescription ? `<p>${factory.shortDescription}</p> ${factory.fullDescription}` : factory.shortDescription,
				};
			})),
			equals: (a, b) => a == b,
			changedEvent: (player: Player<any>) => player.rotationChangeEmitter,
			getValue: (_player: Player<any>) => this.getSourceValue()?.value.oneofKind,
			setValue: (eventID: EventID, player: Player<any>, newKind: APLValueKind) => {
				const sourceValue = this.getSourceValue();
				const oldKind = sourceValue?.value.oneofKind;
				if (oldKind == newKind) {
					return;
				}

				if (newKind) {
					const factory = valueKindFactories[newKind];
					let newSourceValue = this.makeAPLValue(newKind, factory.newValue());
					if (sourceValue) {
						// Some pre-fill logic when swapping kinds.
						if (oldKind && this.valuePicker) {
							if (newKind == 'not') {
								(newSourceValue.value as APLValueImplStruct<'not'>).not.val = this.makeAPLValue(oldKind, this.valuePicker.getInputValue());
							} else if (sourceValue.value.oneofKind == 'not' && sourceValue.value.not.val?.value.oneofKind == newKind) {
								newSourceValue = sourceValue.value.not.val;
							} else if (newKind == 'and') {
								if (sourceValue.value.oneofKind == 'or') {
									(newSourceValue.value as APLValueImplStruct<'and'>).and.vals = sourceValue.value.or.vals;
								} else {
									(newSourceValue.value as APLValueImplStruct<'and'>).and.vals = [this.makeAPLValue(oldKind, this.valuePicker.getInputValue())];
								}
							} else if (newKind == 'or') {
								if (sourceValue.value.oneofKind == 'and') {
									(newSourceValue.value as APLValueImplStruct<'or'>).or.vals = sourceValue.value.and.vals;
								} else {
									(newSourceValue.value as APLValueImplStruct<'or'>).or.vals = [this.makeAPLValue(oldKind, this.valuePicker.getInputValue())];
								}
							} else if (newKind == 'min') {
								if (sourceValue.value.oneofKind == 'max') {
									(newSourceValue.value as APLValueImplStruct<'min'>).min.vals = sourceValue.value.max.vals;
								} else {
									(newSourceValue.value as APLValueImplStruct<'min'>).min.vals = [this.makeAPLValue(oldKind, this.valuePicker.getInputValue())];
								}
							} else if (newKind == 'max') {
								if (sourceValue.value.oneofKind == 'min') {
									(newSourceValue.value as APLValueImplStruct<'max'>).max.vals = sourceValue.value.min.vals;
								} else {
									(newSourceValue.value as APLValueImplStruct<'max'>).max.vals = [this.makeAPLValue(oldKind, this.valuePicker.getInputValue())];
								}
							} else if (sourceValue.value.oneofKind == 'and' && sourceValue.value.and.vals?.[0]?.value.oneofKind == newKind) {
								newSourceValue = sourceValue.value.and.vals[0];
							} else if (sourceValue.value.oneofKind == 'or' && sourceValue.value.or.vals?.[0]?.value.oneofKind == newKind) {
								newSourceValue = sourceValue.value.or.vals[0];
							} else if (sourceValue.value.oneofKind == 'min' && sourceValue.value.min.vals?.[0]?.value.oneofKind == newKind) {
								newSourceValue = sourceValue.value.min.vals[0];
							} else if (sourceValue.value.oneofKind == 'max' && sourceValue.value.max.vals?.[0]?.value.oneofKind == newKind) {
								newSourceValue = sourceValue.value.max.vals[0];
							} else if (newKind == 'cmp') {
								(newSourceValue.value as APLValueImplStruct<'cmp'>).cmp.lhs = this.makeAPLValue(oldKind, this.valuePicker.getInputValue());
							}
						}
					}
					if (sourceValue) {
						sourceValue.value = newSourceValue.value;
					} else {
						this.setSourceValue(eventID, newSourceValue);
					}
				} else {
					this.setSourceValue(eventID, undefined);
				}
				player.rotationChangeEmitter.emit(eventID);
			},
		});

		this.currentKind = undefined;
		this.valuePicker = null;

		this.init();
	}

	getInputElem(): HTMLElement | null {
		return this.rootElem;
	}

	getInputValue(): APLValue | undefined {
		const kind = this.kindPicker.getInputValue();
		if (!kind) {
			return undefined;
		} else {
			return APLValue.create({
				value: {
					oneofKind: kind,
					...((() => {
						const val: any = {};
						if (kind && this.valuePicker) {
							val[kind] = this.valuePicker.getInputValue();
						}
						return val;
					})()),
				},
			})
		}
	}

	setInputValue(newValue: APLValue | undefined) {
		const newKind = newValue?.value.oneofKind;
		this.updateValuePicker(newKind);

		if (newKind && newValue) {
			this.valuePicker!.setInputValue((newValue.value as any)[newKind]);
		}
	}

	private makeAPLValue<K extends NonNullable<APLValueKind>>(kind: K, implVal: APLValueImplTypesUnion[K]): APLValue {
		if (!kind) {
			return APLValue.create();
		}
		const obj: any = { oneofKind: kind };
		obj[kind] = implVal;
		return APLValue.create({value: obj});
	}

	private updateValuePicker(newKind: APLValueKind) {
		const oldKind = this.currentKind;
		if (newKind == oldKind) {
			return;
		}
		this.currentKind = newKind;

		if (this.valuePicker) {
			this.valuePicker.rootElem.remove();
			this.valuePicker = null;
		}

		if (!newKind) {
			return;
		}

		this.kindPicker.setInputValue(newKind);

		const factory = valueKindFactories[newKind];
		this.valuePicker = factory.factory(this.rootElem, this.modObject, {
			changedEvent: (player: Player<any>) => player.rotationChangeEmitter,
			getValue: () => {
				const sourceVal = this.getSourceValue();
				return sourceVal ? (sourceVal.value as any)[newKind] || factory.newValue() : factory.newValue();
			},
			setValue: (eventID: EventID, player: Player<any>, newValue: any) => {
				const sourceVal = this.getSourceValue();
				if (sourceVal) {
					(sourceVal.value as any)[newKind] = newValue;
				}
				player.rotationChangeEmitter.emit(eventID);
			},
		});
	}
}

type ValueKindConfig<T> = {
	label: string,
	submenu?: Array<string>,
	shortDescription: string,
	fullDescription?: string,
	newValue: () => T,
	includeIf?: (player: Player<any>, isPrepull: boolean) => boolean,
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

function mathOperatorFieldConfig(field: string): AplHelpers.APLPickerBuilderFieldConfig<any, any> {
	return {
		field: field,
		newValue: () => MathOperator.OpAdd,
		factory: (parent, player, config) => new TextDropdownPicker(parent, player, {
			...config,
			defaultLabel: 'None',
			equals: (a, b) => a == b,
			values: [
				{ value: MathOperator.OpAdd, label: '+' },
				{ value: MathOperator.OpSub, label: '-' },
				{ value: MathOperator.OpMul, label: '*' },
				{ value: MathOperator.OpDiv, label: '/' },
			],
		}),
	};
}

function executePhaseThresholdFieldConfig(field: string): AplHelpers.APLPickerBuilderFieldConfig<any, any> {
	return {
		field: field,
		newValue: () => ExecutePhaseThreshold.E20,
		factory: (parent, player, config) => new TextDropdownPicker(parent, player, {
			...config,
			defaultLabel: 'None',
			equals: (a, b) => a == b,
			values: [
				{ value: ExecutePhaseThreshold.E20, label: '20%' },
				{ value: ExecutePhaseThreshold.E25, label: '25%' },
				{ value: ExecutePhaseThreshold.E35, label: '35%' },
			],
		}),
	};
}

function totemTypeFieldConfig(field: string): AplHelpers.APLPickerBuilderFieldConfig<any, any> {
	return {
		field: field,
		newValue: () => TotemType.Water,
		factory: (parent, player, config) => new TextDropdownPicker(parent, player, {
			...config,
			defaultLabel: 'None',
			equals: (a, b) => a == b,
			values: [
				{ value: TotemType.Earth, label: 'Earth' },
				{ value: TotemType.Air, label: 'Air' },
				{ value: TotemType.Fire, label: 'Fire' },
				{ value: TotemType.Water, label: 'Water' },
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
			allowedActions: ['create', 'delete'],
			actions: {
				create: {
					useIcon: true,
				},
			},
		}),
	};
}

function inputBuilder<T extends APLValueImplType>(config: {
	label: string,
	submenu?: Array<string>,
	shortDescription: string,
	fullDescription?: string,
	newValue: () => T,
	includeIf?: (player: Player<any>, isPrepull: boolean) => boolean,
	fields: Array<AplHelpers.APLPickerBuilderFieldConfig<T, keyof T>>,
}): ValueKindConfig<T> {
	return {
		label: config.label,
		submenu: config.submenu,
		shortDescription: config.shortDescription,
		fullDescription: config.fullDescription,
		newValue: config.newValue,
		includeIf: config.includeIf,
		factory: AplHelpers.aplInputBuilder(config.newValue, config.fields),
	};
}

const valueKindFactories: {[f in NonNullable<APLValueKind>]: ValueKindConfig<APLValueImplTypesUnion[f]>} = {
	// Operators
	'const': inputBuilder({
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
	'cmp': inputBuilder({
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
	'math': inputBuilder({
		label: 'Math',
		submenu: ['Logic'],
		shortDescription: 'Do basic math on two values.',
		newValue: APLValueMath.create,
		fields: [
			valueFieldConfig('lhs'),
			mathOperatorFieldConfig('op'),
			valueFieldConfig('rhs'),
		],
	}),
	'max': inputBuilder({
		label: 'Max',
		submenu: ['Logic'],
		shortDescription: 'Returns the largest value among the subvalues.',
		newValue: APLValueMax.create,
		fields: [
			valueListFieldConfig('vals'),
		],
	}),
	'min': inputBuilder({
		label: 'Min',
		submenu: ['Logic'],
		shortDescription: 'Returns the smallest value among the subvalues.',
		newValue: APLValueMin.create,
		fields: [
			valueListFieldConfig('vals'),
		],
	}),
	'and': inputBuilder({
		label: 'All of',
		submenu: ['Logic'],
		shortDescription: 'Returns <b>True</b> if all of the sub-values are <b>True</b>, otherwise <b>False</b>',
		newValue: APLValueAnd.create,
		fields: [
			valueListFieldConfig('vals'),
		],
	}),
	'or': inputBuilder({
		label: 'Any of',
		submenu: ['Logic'],
		shortDescription: 'Returns <b>True</b> if any of the sub-values are <b>True</b>, otherwise <b>False</b>',
		newValue: APLValueOr.create,
		fields: [
			valueListFieldConfig('vals'),
		],
	}),
	'not': inputBuilder({
		label: 'Not',
		submenu: ['Logic'],
		shortDescription: 'Returns the opposite of the inner value, i.e. <b>True</b> if the value is <b>False</b> and vice-versa.',
		newValue: APLValueNot.create,
		fields: [
			valueFieldConfig('val'),
		],
	}),

	// Encounter
	'currentTime': inputBuilder({
		label: 'Current Time',
		submenu: ['Encounter'],
		shortDescription: 'Elapsed time of the current sim iteration.',
		newValue: APLValueCurrentTime.create,
		fields: [],
	}),
	'currentTimePercent': inputBuilder({
		label: 'Current Time (%)',
		submenu: ['Encounter'],
		shortDescription: 'Elapsed time of the current sim iteration, as a percentage.',
		newValue: APLValueCurrentTimePercent.create,
		fields: [],
	}),
	'remainingTime': inputBuilder({
		label: 'Remaining Time',
		submenu: ['Encounter'],
		shortDescription: 'Elapsed time of the remaining sim iteration.',
		newValue: APLValueRemainingTime.create,
		fields: [],
	}),
	'remainingTimePercent': inputBuilder({
		label: 'Remaining Time (%)',
		submenu: ['Encounter'],
		shortDescription: 'Elapsed time of the remaining sim iteration, as a percentage.',
		newValue: APLValueRemainingTimePercent.create,
		fields: [],
	}),
	'isExecutePhase': inputBuilder({
		label: 'Is Execute Phase',
		submenu: ['Encounter'],
		shortDescription: '<b>True</b> if the encounter is in Execute Phase, meaning the target\'s health is less than the given threshold, otherwise <b>False</b>.',
		newValue: APLValueIsExecutePhase.create,
		fields: [
			executePhaseThresholdFieldConfig('threshold'),
		],
	}),
	'numberTargets': inputBuilder({
		label: 'Number of Targets',
		submenu: ['Encounter'],
		shortDescription: 'Count of targets in the current encounter',
		newValue: APLValueNumberTargets.create,
		fields: [],
	}),
	'frontOfTarget': inputBuilder({
		label: 'Front of Target',
		submenu: ['Encounter'],
		shortDescription: '<b>True</b> if facing from of target',
		newValue: APLValueFrontOfTarget.create,
		fields: [],
	}),

	// Boss
	'bossSpellIsCasting': inputBuilder({
		label: 'Spell is Casting',
		submenu: ['Boss'],
		shortDescription: '',
		newValue: APLValueBossSpellIsCasting.create,
		fields: [
			AplHelpers.unitFieldConfig('targetUnit', 'targets'),
			AplHelpers.actionIdFieldConfig('spellId', 'non_instant_spells', 'targetUnit', 'currentTarget'),
		]
	}),
	'bossSpellTimeToReady': inputBuilder({
		label: 'Spell Time to Ready',
		submenu: ['Boss'],
		shortDescription: '',
		newValue: APLValueBossSpellTimeToReady.create,
		fields: [
			AplHelpers.unitFieldConfig('targetUnit', 'targets'),
			AplHelpers.actionIdFieldConfig('spellId', 'spells', 'targetUnit', 'currentTarget'),
		]
	}),

	// Resources
	'currentHealth': inputBuilder({
		label: 'Health',
		submenu: ['Resources'],
		shortDescription: 'Amount of currently available Health.',
		newValue: APLValueCurrentHealth.create,
		fields: [
			AplHelpers.unitFieldConfig('sourceUnit', 'aura_sources'),
		],
	}),
	'currentHealthPercent': inputBuilder({
		label: 'Health (%)',
		submenu: ['Resources'],
		shortDescription: 'Amount of currently available Health, as a percentage.',
		newValue: APLValueCurrentHealthPercent.create,
		fields: [
			AplHelpers.unitFieldConfig('sourceUnit', 'aura_sources'),
		],
	}),
	'currentMana': inputBuilder({
		label: 'Mana',
		submenu: ['Resources'],
		shortDescription: 'Amount of currently available Mana.',
		newValue: APLValueCurrentMana.create,
		fields: [],
	}),
	'currentManaPercent': inputBuilder({
		label: 'Mana (%)',
		submenu: ['Resources'],
		shortDescription: 'Amount of currently available Mana, as a percentage.',
		newValue: APLValueCurrentManaPercent.create,
		fields: [],
	}),
	'currentRage': inputBuilder({
		label: 'Rage',
		submenu: ['Resources'],
		shortDescription: 'Amount of currently available Rage.',
		newValue: APLValueCurrentRage.create,
		fields: [],
	}),
	'currentEnergy': inputBuilder({
		label: 'Energy',
		submenu: ['Resources'],
		shortDescription: 'Amount of currently available Energy.',
		newValue: APLValueCurrentEnergy.create,
		fields: [],
	}),
	'currentComboPoints': inputBuilder({
		label: 'Combo Points',
		submenu: ['Resources'],
		shortDescription: 'Amount of currently available Combo Points.',
		newValue: APLValueCurrentComboPoints.create,
		fields: [],
	}),
	'currentRunicPower': inputBuilder({
		label: 'Runic Power',
		submenu: ['Resources'],
		shortDescription: 'Amount of currently available Runic Power.',
		newValue: APLValueCurrentRunicPower.create,
		includeIf: (player: Player<any>, isPrepull: boolean) => player.getClass() == Class.ClassDeathknight,
		fields: [],
	}),

	// Resources Rune
	'currentRuneCount': inputBuilder({
		label: 'Num Runes',
		submenu: ['Resources', 'Runes'],
		shortDescription: 'Amount of currently available Runes of certain type including Death.',
		newValue: APLValueCurrentRuneCount.create,
		includeIf: (player: Player<any>, isPrepull: boolean) => player.getClass() == Class.ClassDeathknight,
		fields: [
			AplHelpers.runeTypeFieldConfig('runeType', true),
		],
	}),
	'currentNonDeathRuneCount': inputBuilder({
		label: 'Num Non Death Runes',
		submenu: ['Resources', 'Runes'],
		shortDescription: 'Amount of currently available Runes of certain type ignoring Death',
		newValue: APLValueCurrentNonDeathRuneCount.create,
		includeIf: (player: Player<any>, isPrepull: boolean) => player.getClass() == Class.ClassDeathknight,
		fields: [
			AplHelpers.runeTypeFieldConfig('runeType', false),
		],
	}),
	'currentRuneActive': inputBuilder({
		label: 'Rune Is Ready',
		submenu: ['Resources', 'Runes'],
		shortDescription: 'Is the rune of a certain slot currently available.',
		newValue: APLValueCurrentRuneActive.create,
		includeIf: (player: Player<any>, isPrepull: boolean) => player.getClass() == Class.ClassDeathknight,
		fields: [
			AplHelpers.runeSlotFieldConfig('runeSlot'),
		],
	}),
	'currentRuneDeath': inputBuilder({
		label: 'Rune Is Death',
		submenu: ['Resources', 'Runes'],
		shortDescription: 'Is the rune of a certain slot currently converted to Death.',
		newValue: APLValueCurrentRuneDeath.create,
		includeIf: (player: Player<any>, isPrepull: boolean) => player.getClass() == Class.ClassDeathknight,
		fields: [
			AplHelpers.runeSlotFieldConfig('runeSlot'),
		],
	}),
	'runeCooldown': inputBuilder({
		label: 'Rune Cooldown',
		submenu: ['Resources', 'Runes'],
		shortDescription: 'Amount of time until a rune of certain type is ready to use.<br><b>NOTE:</b> Returns 0 if there is a rune available',
		newValue: APLValueRuneCooldown.create,
		includeIf: (player: Player<any>, isPrepull: boolean) => player.getClass() == Class.ClassDeathknight,
		fields: [
			AplHelpers.runeTypeFieldConfig('runeType', false),
		],
	}),
	'nextRuneCooldown': inputBuilder({
		label: 'Next Rune Cooldown',
		submenu: ['Resources', 'Runes'],
		shortDescription: 'Amount of time until a 2nd rune of certain type is ready to use.<br><b>NOTE:</b> Returns 0 if there are 2 runes available',
		newValue: APLValueNextRuneCooldown.create,
		includeIf: (player: Player<any>, isPrepull: boolean) => player.getClass() == Class.ClassDeathknight,
		fields: [
			AplHelpers.runeTypeFieldConfig('runeType', false),
		],
	}),
	'runeSlotCooldown': inputBuilder({
		label: 'Rune Slot Cooldown',
		submenu: ['Resources', 'Runes'],
		shortDescription: 'Amount of time until a rune of certain slot is ready to use.<br><b>NOTE:</b> Returns 0 if rune is ready',
		newValue: APLValueRuneSlotCooldown.create,
		includeIf: (player: Player<any>, isPrepull: boolean) => player.getClass() == Class.ClassDeathknight,
		fields: [
			AplHelpers.runeSlotFieldConfig('runeSlot'),
		],
	}),
	'runeGrace': inputBuilder({
		label: 'Rune Grace Period',
		submenu: ['Resources', 'Runes'],
		shortDescription: 'Amount of rune grace period available for certain rune type.',
		newValue: APLValueRuneGrace.create,
		includeIf: (player: Player<any>, isPrepull: boolean) => player.getClass() == Class.ClassDeathknight,
		fields: [
			AplHelpers.runeTypeFieldConfig('runeType', false),
		],
	}),
	'runeSlotGrace': inputBuilder({
		label: 'Rune Slot Grace Period',
		submenu: ['Resources', 'Runes'],
		shortDescription: 'Amount of rune grace period available for certain rune slot.',
		newValue: APLValueRuneSlotGrace.create,
		includeIf: (player: Player<any>, isPrepull: boolean) => player.getClass() == Class.ClassDeathknight,
		fields: [
			AplHelpers.runeSlotFieldConfig('runeSlot'),
		],
	}),

	// GCD
	'gcdIsReady': inputBuilder({
		label: 'GCD Is Ready',
		submenu: ['GCD'],
		shortDescription: '<b>True</b> if the GCD is not on cooldown, otherwise <b>False</b>.',
		newValue: APLValueGCDIsReady.create,
		fields: [],
	}),
	'gcdTimeToReady': inputBuilder({
		label: 'GCD Time To Ready',
		submenu: ['GCD'],
		shortDescription: 'Amount of time remaining before the GCD comes off cooldown, or <b>0</b> if it is not on cooldown.',
		newValue: APLValueGCDTimeToReady.create,
		fields: [],
	}),

	// Auto attacks
	'autoTimeToNext': inputBuilder({
		label: 'Time To Next Auto',
		submenu: ['Auto'],
		shortDescription: 'Amount of time remaining before the next Main-hand or Off-hand melee attack, or <b>0</b> if autoattacks are not engaged.',
		newValue: APLValueAutoTimeToNext.create,
		fields: [],
	}),

	// Spells
	'spellCurrentCost': inputBuilder({
		label: 'Current Cost',
		submenu: ['Spell'],
		shortDescription: 'Returns current resource cost of spell',
		newValue: APLValueSpellCurrentCost.create,
		fields: [
			AplHelpers.actionIdFieldConfig('spellId', 'castable_spells', ''),
		],
	}),
	'spellCanCast': inputBuilder({
		label: 'Can Cast',
		submenu: ['Spell'],
		shortDescription: '<b>True</b> if all requirements for casting the spell are currently met, otherwise <b>False</b>.',
		fullDescription: `
			<p>The <b>Cast Spell</b> action does not need to be conditioned on this, because it applies this check automatically.</p>
		`,
		newValue: APLValueSpellCanCast.create,
		fields: [
			AplHelpers.actionIdFieldConfig('spellId', 'castable_spells', ''),
		],
	}),
	'spellIsReady': inputBuilder({
		label: 'Is Ready',
		submenu: ['Spell'],
		shortDescription: '<b>True</b> if the spell is not on cooldown, otherwise <b>False</b>.',
		newValue: APLValueSpellIsReady.create,
		fields: [
			AplHelpers.actionIdFieldConfig('spellId', 'castable_spells', ''),
		],
	}),
	'spellTimeToReady': inputBuilder({
		label: 'Time To Ready',
		submenu: ['Spell'],
		shortDescription: 'Amount of time remaining before the spell comes off cooldown, or <b>0</b> if it is not on cooldown.',
		newValue: APLValueSpellTimeToReady.create,
		fields: [
			AplHelpers.actionIdFieldConfig('spellId', 'castable_spells', ''),
		],
	}),
	'spellCastTime': inputBuilder({
		label: 'Cast Time',
		submenu: ['Spell'],
		shortDescription: 'Amount of time to cast the spell including any haste and spell cast time adjustments.',
		newValue: APLValueSpellCastTime.create,
		fields: [
			AplHelpers.actionIdFieldConfig('spellId', 'castable_spells', ''),
		],
	}),
	'spellTravelTime': inputBuilder({
		label: 'Travel Time',
		submenu: ['Spell'],
		shortDescription: 'Amount of time for the spell to travel to the target.',
		newValue: APLValueSpellTravelTime.create,
		fields: [
			AplHelpers.actionIdFieldConfig('spellId', 'castable_spells', ''),
		],
	}),
	'spellCpm': inputBuilder({
		label: 'CPM',
		submenu: ['Spell'],
		shortDescription: 'Casts Per Minute for the spell so far in the current iteration.',
		newValue: APLValueSpellCPM.create,
		fields: [
			AplHelpers.actionIdFieldConfig('spellId', 'castable_spells', ''),
		],
	}),
	'spellIsChanneling': inputBuilder({
		label: 'Is Channeling',
		submenu: ['Spell'],
		shortDescription: '<b>True</b> if this spell is currently being channeled, otherwise <b>False</b>.',
		newValue: APLValueSpellIsChanneling.create,
		fields: [
			AplHelpers.actionIdFieldConfig('spellId', 'channel_spells', ''),
		],
	}),
	'spellChanneledTicks': inputBuilder({
		label: 'Channeled Ticks',
		submenu: ['Spell'],
		shortDescription: 'The number of completed ticks in the current channel of this spell, or <b>0</b> if the spell is not being channeled.',
		newValue: APLValueSpellChanneledTicks.create,
		fields: [
			AplHelpers.actionIdFieldConfig('spellId', 'channel_spells', ''),
		],
	}),
	'channelClipDelay': inputBuilder({
		label: 'Channel Clip Delay',
		submenu: ['Spell'],
		shortDescription: 'The amount of time specified by the <b>Channel Clip Delay</b> setting.',
		newValue: APLValueChannelClipDelay.create,
		fields: [
		],
	}),

	// Auras
	'auraIsActive': inputBuilder({
		label: 'Aura Active',
		submenu: ['Aura'],
		shortDescription: '<b>True</b> if the aura is currently active, otherwise <b>False</b>.',
		newValue: APLValueAuraIsActive.create,
		fields: [
			AplHelpers.unitFieldConfig('sourceUnit', 'aura_sources'),
			AplHelpers.actionIdFieldConfig('auraId', 'auras', 'sourceUnit'),
		],
	}),
	'auraIsActiveWithReactionTime': inputBuilder({
		label: 'Aura Active (with Reaction Time)',
		submenu: ['Aura'],
		shortDescription: '<b>True</b> if the aura is currently active AND it has been active for at least as long as the player reaction time (configured in Settings), otherwise <b>False</b>.',
		newValue: APLValueAuraIsActiveWithReactionTime.create,
		fields: [
			AplHelpers.unitFieldConfig('sourceUnit', 'aura_sources'),
			AplHelpers.actionIdFieldConfig('auraId', 'auras', 'sourceUnit'),
		],
	}),
	'auraRemainingTime': inputBuilder({
		label: 'Aura Remaining Time',
		submenu: ['Aura'],
		shortDescription: 'Time remaining before this aura will expire, or 0 if the aura is not currently active.',
		newValue: APLValueAuraRemainingTime.create,
		fields: [
			AplHelpers.unitFieldConfig('sourceUnit', 'aura_sources'),
			AplHelpers.actionIdFieldConfig('auraId', 'auras', 'sourceUnit'),
		],
	}),
	'auraNumStacks': inputBuilder({
		label: 'Aura Num Stacks',
		submenu: ['Aura'],
		shortDescription: 'Number of stacks of the aura.',
		newValue: APLValueAuraNumStacks.create,
		fields: [
			AplHelpers.unitFieldConfig('sourceUnit', 'aura_sources'),
			AplHelpers.actionIdFieldConfig('auraId', 'stackable_auras', 'sourceUnit'),
		],
	}),
	'auraInternalCooldown': inputBuilder({
		label: 'Aura Remaining ICD',
		submenu: ['Aura'],
		shortDescription: 'Time remaining before this aura\'s internal cooldown will be ready, or <b>0</b> if the ICD is ready now.',
		newValue: APLValueAuraInternalCooldown.create,
		fields: [
			AplHelpers.unitFieldConfig('sourceUnit', 'aura_sources'),
			AplHelpers.actionIdFieldConfig('auraId', 'icd_auras', 'sourceUnit'),
		],
	}),
	'auraIcdIsReadyWithReactionTime': inputBuilder({
		label: 'Aura ICD Is Ready (with Reaction Time)',
		submenu: ['Aura'],
		shortDescription: '<b>True</b> if the aura\'s ICD is currently ready OR it was put on CD recently, within the player\'s reaction time (configured in Settings), otherwise <b>False</b>.',
		newValue: APLValueAuraICDIsReadyWithReactionTime.create,
		fields: [
			AplHelpers.unitFieldConfig('sourceUnit', 'aura_sources'),
			AplHelpers.actionIdFieldConfig('auraId', 'icd_auras', 'sourceUnit'),
		],
	}),
	'auraShouldRefresh': inputBuilder({
		label: 'Should Refresh Aura',
		submenu: ['Aura'],
		shortDescription: 'Whether this aura should be refreshed, e.g. for the purpose of maintaining a debuff.',
		fullDescription: `
		<p>This condition checks not only the specified aura but also any other auras on the same unit, including auras applied by other raid members, which apply the same debuff category.</p>
		<p>For example, 'Should Refresh Debuff(Sunder Armor)' will return <b>False</b> if the unit has an active Expose Armor aura.</p>
		`,
		newValue: () => APLValueAuraShouldRefresh.create({
			maxOverlap: {
				value: {
					oneofKind: 'const',
					const: {
						val: '0ms',
					},
				},
			},
		}),
		fields: [
			AplHelpers.unitFieldConfig('sourceUnit', 'aura_sources_targets_first'),
			AplHelpers.actionIdFieldConfig('auraId', 'exclusive_effect_auras', 'sourceUnit', 'currentTarget'),
			valueFieldConfig('maxOverlap', {
				label: 'Overlap',
				labelTooltip: 'Maximum amount of time before the aura expires when it may be refreshed.',
			}),
		],
	}),

	// DoT
	'dotIsActive': inputBuilder({
		label: 'Dot Is Active',
		submenu: ['DoT'],
		shortDescription: '<b>True</b> if the specified dot is currently ticking, otherwise <b>False</b>.',
		newValue: APLValueDotIsActive.create,
		fields: [
			AplHelpers.unitFieldConfig('targetUnit', 'targets'),
			AplHelpers.actionIdFieldConfig('spellId', 'dot_spells', ''),
		],
	}),
	'dotRemainingTime': inputBuilder({
		label: 'Dot Remaining Time',
		submenu: ['DoT'],
		shortDescription: 'Time remaining before the last tick of this DoT will occur, or 0 if the DoT is not currently ticking.',
		newValue: APLValueDotRemainingTime.create,
		fields: [
			AplHelpers.unitFieldConfig('targetUnit', 'targets'),
			AplHelpers.actionIdFieldConfig('spellId', 'dot_spells', ''),
		],
	}),
	'sequenceIsComplete': inputBuilder({
		label: 'Sequence Is Complete',
		submenu: ['Sequence'],
		shortDescription: '<b>True</b> if there are no more subactions left to execute in the sequence, otherwise <b>False</b>.',
		newValue: APLValueSequenceIsComplete.create,
		fields: [
			AplHelpers.stringFieldConfig('sequenceName'),
		],
	}),
	'sequenceIsReady': inputBuilder({
		label: 'Sequence Is Ready',
		submenu: ['Sequence'],
		shortDescription: '<b>True</b> if the next subaction in the sequence is ready to be executed, otherwise <b>False</b>.',
		newValue: APLValueSequenceIsReady.create,
		fields: [
			AplHelpers.stringFieldConfig('sequenceName'),
		],
	}),
	'sequenceTimeToReady': inputBuilder({
		label: 'Sequence Time To Ready',
		submenu: ['Sequence'],
		shortDescription: 'Returns the amount of time remaining until the next subaction in the sequence will be ready.',
		newValue: APLValueSequenceTimeToReady.create,
		fields: [
			AplHelpers.stringFieldConfig('sequenceName'),
		],
	}),

	// Class/spec specific values
	'totemRemainingTime': inputBuilder({
		label: 'Totem Remaining Time',
		submenu: ['Shaman'],
		shortDescription: 'Returns the amount of time remaining until the totem will expire.',
		newValue: APLValueTotemRemainingTime.create,
		includeIf: (player: Player<any>, isPrepull: boolean) => player.getClass() == Class.ClassShaman,
		fields: [
			totemTypeFieldConfig('totemType'),
		],
	}),
	'catExcessEnergy': inputBuilder({
		label: 'Excess Energy',
		submenu: ['Feral Druid'],
		shortDescription: 'Returns the amount of excess energy available, after subtracting energy that will be needed to maintain DoTs.',
		newValue: APLValueCatExcessEnergy.create,
		includeIf: (player: Player<any>, isPrepull: boolean) => player.spec == Spec.SpecFeralDruid,
		fields: [
		],
	}),
	'catNewSavageRoarDuration': inputBuilder({
		label: 'New Savage Roar Duration',
		submenu: ['Feral Druid'],
		shortDescription: 'Returns duration of savage roar based on current combo points',
		newValue: APLValueCatNewSavageRoarDuration.create,
		includeIf: (player: Player<any>, isPrepull: boolean) => player.spec == Spec.SpecFeralDruid,
		fields: [
		],
	}),
	'warlockShouldRecastDrainSoul': inputBuilder({
		label: 'Should Recast Drain Soul',
		submenu: ['Warlock'],
		shortDescription: 'Returns <b>True</b> if the current Drain Soul channel should be immediately recast, to get a better snapshot.',
		newValue: APLValueWarlockShouldRecastDrainSoul.create,
		includeIf: (player: Player<any>, isPrepull: boolean) => player.getClass() == Class.ClassWarlock,
		fields: [
		],
	}),
	'warlockShouldRefreshCorruption': inputBuilder({
		label: 'Should Refresh Corruption',
		submenu: ['Warlock'],
		shortDescription: 'Returns <b>True</b> if the current Corruption has expired, or should be refreshed to get a better snapshot.',
		newValue: APLValueWarlockShouldRefreshCorruption.create,
		includeIf: (player: Player<any>, isPrepull: boolean) => player.getClass() == Class.ClassWarlock,
		fields: [
			AplHelpers.unitFieldConfig('targetUnit', 'targets'),
		],
	}),
};
