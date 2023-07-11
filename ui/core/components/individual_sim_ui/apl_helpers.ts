import { OtherAction } from '../../proto/common.js';
import { ActionId } from '../../proto_utils/action_id.js';
import { Player } from '../../player.js';
import { EventID, TypedEvent } from '../../typed_event.js';
import { bucket } from '../../utils.js';
import { AdaptiveStringPicker } from '../string_picker.js';
import { NumberPicker } from '../number_picker.js';
import { DropdownPicker, DropdownPickerConfig, DropdownValueConfig, TextDropdownPicker } from '../dropdown_picker.js';
import { Input, InputConfig } from '../input.js';
import { ActionID } from '../../proto/common.js';
import { BooleanPicker } from '../boolean_picker.js';
import { APLValueRuneSlot, APLValueRuneType } from '../../proto/apl.js';

export type ACTION_ID_SET = 'auras' | 'stackable_auras' | 'castable_spells' | 'dot_spells';

const actionIdSets: Record<ACTION_ID_SET, {
	defaultLabel: string,
	getActionIDs: (player: Player<any>) => Promise<Array<DropdownValueConfig<ActionId>>>,
}> = {
	['auras']: {
		defaultLabel: 'Aura',
		getActionIDs: async (player) => {
			return player.getAuras().map(actionId => {
				return {
					value: actionId.id,
				};
			});
		},
	},
	['stackable_auras']: {
		defaultLabel: 'Aura',
		getActionIDs: async (player) => {
			return player.getAuras().filter(aura => aura.data.maxStacks > 0).map(actionId => {
				return {
					value: actionId.id,
				};
			});
		},
	},
	['castable_spells']: {
		defaultLabel: 'Spell',
		getActionIDs: async (player) => {
			const castableSpells = player.getSpells().filter(spell => spell.data.isCastable);

			// Split up non-cooldowns and cooldowns into separate sections for easier browsing.
			const { 'spells': spells, 'cooldowns': cooldowns } = bucket(castableSpells, spell => spell.data.isMajorCooldown ? 'cooldowns' : 'spells');

			const placeholders: Array<ActionId> = [
				ActionId.fromOtherId(OtherAction.OtherActionPotion),
			];

			return [
				[{
					value: ActionId.fromEmpty(),
					headerText: 'Spells',
				}],
				(spells || []).map(actionId => {
					return {
						value: actionId.id,
						extraCssClasses: (actionId.data.prepullOnly ? ['apl-prepull-actions-only'] : []),
					};
				}),
				[{
					value: ActionId.fromEmpty(),
					headerText: 'Cooldowns',
				}],
				(cooldowns || []).map(actionId => {
					return {
						value: actionId.id,
						extraCssClasses: (actionId.data.prepullOnly ? ['apl-prepull-actions-only'] : []),
					};
				}),
				[{
					value: ActionId.fromEmpty(),
					headerText: 'Placeholders',
				}],
				placeholders.map(actionId => {
					return {
						value: actionId,
						tooltip: 'The Prepull Potion if CurrentTime < 0, or the Combat Potion if combat has started.',
					};
				}),
			].flat();
		},
	},
	['dot_spells']: {
		defaultLabel: 'DoT Spell',
		getActionIDs: async (player) => {
			return player.getSpells().filter(spell => spell.data.hasDot).map(actionId => {
				return {
					value: actionId.id,
				};
			});
		},
	},
};

export interface APLActionIDPickerConfig<ModObject> extends Omit<DropdownPickerConfig<ModObject, ActionId>, 'defaultLabel' | 'equals' | 'setOptionContent' | 'values' | 'getValue' | 'setValue'> {
	actionIdSet: ACTION_ID_SET,
	getValue: (obj: ModObject) => ActionID,
	setValue: (eventID: EventID, obj: ModObject, newValue: ActionID) => void,
}

export class APLActionIDPicker extends DropdownPicker<Player<any>, ActionId> {
	constructor(parent: HTMLElement, player: Player<any>, config: APLActionIDPickerConfig<Player<any>>) {
		const actionIdSet = actionIdSets[config.actionIdSet];
		super(parent, player, {
			...config,
			getValue: (player) => ActionId.fromProto(config.getValue(player)),
			setValue: (eventID: EventID, player: Player<any>, newValue: ActionId) => config.setValue(eventID, player, newValue.toProto()),
			defaultLabel: actionIdSet.defaultLabel,
			equals: (a, b) => ((a == null) == (b == null)) && (!a || a.equals(b!)),
			setOptionContent: (button, valueConfig) => {
				const actionId = valueConfig.value;

				const iconElem = document.createElement('a');
				iconElem.classList.add('apl-actionid-item-icon');
				actionId.setBackgroundAndHref(iconElem);
				button.appendChild(iconElem);

				const textElem = document.createTextNode(actionId.name);
				button.appendChild(textElem);
			},
			createMissingValue: value => ((value instanceof ActionId) ? value : ActionId.fromProto(value as unknown as ActionID)).fill().then(filledId => {
				return {
					value: filledId,
				};
			}),
			values: [],
		});

		const getActionIDs = actionIdSet.getActionIDs;
		const updateValues = async () => {
			const values = await getActionIDs(player);
			this.setOptions(values);
		};
		updateValues();
		player.currentSpellsAndAurasEmitter.on(updateValues);
	}
}

export interface APLPickerBuilderFieldConfig<T, F extends keyof T> {
	field: F,
	newValue: () => T[F],
	factory: (parent: HTMLElement, player: Player<any>, config: InputConfig<Player<any>, T[F]>) => Input<Player<any>, T[F]>

	label?: string,
	labelTooltip?: string,
}

export interface APLPickerBuilderConfig<T> extends InputConfig<Player<any>, T> {
	newValue: () => T,
	fields: Array<APLPickerBuilderFieldConfig<T, any>>,
}

export interface APLPickerBuilderField<T, F extends keyof T> extends APLPickerBuilderFieldConfig<T, F> {
	picker: Input<Player<any>, T[F]>,
}

export class APLPickerBuilder<T> extends Input<Player<any>, T> {
	private readonly config: APLPickerBuilderConfig<T>;
	private readonly fieldPickers: Array<APLPickerBuilderField<T, any>>;

	constructor(parent: HTMLElement, modObject: Player<any>, config: APLPickerBuilderConfig<T>) {
		super(parent, 'apl-picker-builder-root', modObject, config);
		this.config = config;

		if (config.fields.length > 0) {
			const openSpan = document.createElement('span');
			openSpan.textContent = '(';
			this.rootElem.appendChild(openSpan);
		}

		this.fieldPickers = config.fields.map(fieldConfig => APLPickerBuilder.makeFieldPicker(this, fieldConfig));

		if (config.fields.length > 0) {
			const closeSpan = document.createElement('span');
			closeSpan.textContent = ')';
			this.rootElem.appendChild(closeSpan);
		}

		this.init();
	}

	private static makeFieldPicker<T, F extends keyof T>(builder: APLPickerBuilder<T>, fieldConfig: APLPickerBuilderFieldConfig<T, F>): APLPickerBuilderField<T, F> {
		const field: F = fieldConfig.field;
		return {
			...fieldConfig,
			picker: fieldConfig.factory(builder.rootElem, builder.modObject, {
				label: fieldConfig.label,
				labelTooltip: fieldConfig.labelTooltip,
				inline: true,
				changedEvent: (player: Player<any>) => player.rotationChangeEmitter,
				getValue: () => {
					const source = builder.getSourceValue();
					if (!source[field]) {
						source[field] = fieldConfig.newValue();
					}
					return source[field];
				},
				setValue: (eventID: EventID, player: Player<any>, newValue: any) => {
					builder.getSourceValue()[field] = newValue;
					player.rotationChangeEmitter.emit(eventID);
				},
			}),
		};
	}

	getInputElem(): HTMLElement {
		return this.rootElem;
	}

	getInputValue(): T {
		const val = this.config.newValue();
		this.fieldPickers.forEach(pickerData => {
			val[pickerData.field as keyof T] = pickerData.picker.getInputValue();
		});
		return val;
	}

	setInputValue(newValue: T) {
		this.fieldPickers.forEach(pickerData => {
			pickerData.picker.setInputValue(newValue[pickerData.field as keyof T]);
		});
	}
}

export function actionIdFieldConfig(field: string, actionIdSet: ACTION_ID_SET, options?: Partial<APLPickerBuilderFieldConfig<any, any>>): APLPickerBuilderFieldConfig<any, any> {
	return {
		field: field,
		newValue: () => ActionID.create(),
		factory: (parent, player, config) => new APLActionIDPicker(parent, player, {
			...config,
			actionIdSet: actionIdSet,
		}),
		...(options || {}),
	};
}

export function booleanFieldConfig(field: string, label?:string, options?: Partial<APLPickerBuilderFieldConfig<any, any>>): APLPickerBuilderFieldConfig<any, any> {
	return {
		field: field,
		newValue: () => false,
		factory: (parent, player, config) => new BooleanPicker(parent, player, config),
		...(options || {}),
		label: label,
	};
}

export function numberFieldConfig(field: string, options?: Partial<APLPickerBuilderFieldConfig<any, any>>): APLPickerBuilderFieldConfig<any, any> {
	return {
		field: field,
		newValue: () => 0,
		factory: (parent, player, config) => new NumberPicker(parent, player, config),
		...(options || {}),
	};
}

export function stringFieldConfig(field: string, options?: Partial<APLPickerBuilderFieldConfig<any, any>>): APLPickerBuilderFieldConfig<any, any> {
	return {
		field: field,
		newValue: () => '',
		factory: (parent, player, config) => new AdaptiveStringPicker(parent, player, config),
		...(options || {}),
	};
}

export function runeTypeFieldConfig(field: string, includeDeath: boolean): APLPickerBuilderFieldConfig<any, any> {
	let values = [
		{ value: APLValueRuneType.RuneBlood, label: 'Blood' },
		{ value: APLValueRuneType.RuneFrost, label: 'Frost' },
		{ value: APLValueRuneType.RuneUnholy, label: 'Unholy' },
	]

	if (includeDeath) {
		values.push({ value: APLValueRuneType.RuneDeath, label: 'Death' })
	}

	return {
		field: field,
		newValue: () => APLValueRuneType.RuneBlood,
		factory: (parent, player, config) => new TextDropdownPicker(parent, player, {
			...config,
			defaultLabel: 'None',
			equals: (a, b) => a == b,
			values: values,
		}),
	};
}

export function runeSlotFieldConfig(field: string): APLPickerBuilderFieldConfig<any, any> {
	return {
		field: field,
		newValue: () => APLValueRuneSlot.SlotLeftBlood,
		factory: (parent, player, config) => new TextDropdownPicker(parent, player, {
			...config,
			defaultLabel: 'None',
			equals: (a, b) => a == b,
			values: [
				{ value: APLValueRuneSlot.SlotLeftBlood, label: 'Blood Left' },
				{ value: APLValueRuneSlot.SlotRightBlood, label: 'Blood Right' },
				{ value: APLValueRuneSlot.SlotLeftFrost, label: 'Frost Left' },
				{ value: APLValueRuneSlot.SlotRightFrost, label: 'Frost Right' },
				{ value: APLValueRuneSlot.SlotLeftUnholy, label: 'Unholy Left' },
				{ value: APLValueRuneSlot.SlotRightUnholy, label: 'Unholy Right' },
			],
		}),
	};
}

export function aplInputBuilder<T>(newValue: () => T, fields: Array<APLPickerBuilderFieldConfig<T, any>>): (parent: HTMLElement, player: Player<any>, config: InputConfig<Player<any>, any>) => Input<Player<any>, any> {
	return (parent, player, config) => {
		return new APLPickerBuilder(parent, player, {
			...config,
			newValue: newValue,
			fields: fields,
		})
	}
}
