import { ActionId } from '../../proto_utils/action_id.js';
import { Player } from '../../player.js';
import { EventID, TypedEvent } from '../../typed_event.js';
import { stringComparator } from '../../utils.js';
import { DropdownPicker, DropdownPickerConfig, DropdownValueConfig, TextDropdownPicker } from '../dropdown_picker.js';
import { Input, InputConfig } from '../input.js';
import { ActionID } from '../../proto/common.js';

export type ACTION_ID_SET = 'all_spells' | 'dots';

export interface APLActionIDPickerConfig<ModObject> extends Omit<DropdownPickerConfig<ModObject, ActionId>, 'defaultLabel' | 'equals' | 'setOptionContent' | 'values' | 'getValue' | 'setValue'> {
	actionIdSet: ACTION_ID_SET,
	getValue: (obj: ModObject) => ActionID,
	setValue: (eventID: EventID, obj: ModObject, newValue: ActionID) => void,
}

const actionIdSets: Record<ACTION_ID_SET, {
	defaultLabel: string,
	getActionIDs: (player: Player<any>) => Promise<Array<DropdownValueConfig<ActionId>>>,
}> = {
	['all_spells']: {
		defaultLabel: 'Spell',
		getActionIDs: async (player) => {
			const playerStats = player.getCurrentStats();
			const spellPromises = Promise.all(playerStats.spells.map(spell => ActionId.fromProto(spell).fill()));
			const cooldownPromises = Promise.all(playerStats.cooldowns.map(cd => ActionId.fromProto(cd).fill()));

			let [spells, cooldowns] = await Promise.all([spellPromises, cooldownPromises]);
			spells = spells.sort((a, b) => stringComparator(a.name, b.name))
			cooldowns = cooldowns.sort((a, b) => stringComparator(a.name, b.name))

			return [...spells, ...cooldowns].map(actionId => {
				return {
					value: actionId,
				};
			});
		},
	},
	['dots']: {
		defaultLabel: 'DoT Spell',
		getActionIDs: async (player) => {
			const playerStats = player.getCurrentStats();
			let spells = await Promise.all(playerStats.spells.map(spell => ActionId.fromProto(spell).fill()));
			spells = spells.sort((a, b) => stringComparator(a.name, b.name))

			return spells.map(actionId => {
				return {
					value: actionId,
				};
			});
		},
	},
};

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
			values: [],
		});

		const getActionIDs = actionIdSet.getActionIDs;
		const updateValues = async () => {
			const values = await getActionIDs(player);
            this.setOptions(values);
		};
        updateValues();
        player.currentStatsEmitter.on(updateValues);
	}
}

export interface APLPickerBuilderFieldConfig<T, F extends keyof T> {
	field: F,
	newValue: () => T[F],
	factory: (parent: HTMLElement, player: Player<any>, config: InputConfig<Player<any>, T[F]>) => Input<Player<any>, T[F]>
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

		const openSpan = document.createElement('span');
		openSpan.textContent = '(';
		this.rootElem.appendChild(openSpan);

		this.fieldPickers = config.fields.map(fieldConfig => APLPickerBuilder.makeFieldPicker(this, fieldConfig));

		const closeSpan = document.createElement('span');
		closeSpan.textContent = ')';
		this.rootElem.appendChild(closeSpan);

		this.init();
	}

	private static makeFieldPicker<T, F extends keyof T>(builder: APLPickerBuilder<T>, fieldConfig: APLPickerBuilderFieldConfig<T, F>): APLPickerBuilderField<T, F> {
		const field: F = fieldConfig.field;
		return {
			...fieldConfig,
			picker: fieldConfig.factory(builder.rootElem, builder.modObject, {
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

export function actionIdFieldConfig(field: string, actionIdSet: ACTION_ID_SET): APLPickerBuilderFieldConfig<any, any> {
	return {
		field: field,
		newValue: () => ActionID.create(),
		factory: (parent, player, config) => new APLActionIDPicker(parent, player, {
			...config,
			actionIdSet: actionIdSet,
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