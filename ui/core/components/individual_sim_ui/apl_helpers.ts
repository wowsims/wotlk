import { Player, UnitMetadata } from '../../player.js';
import { APLValueRuneSlot, APLValueRuneType } from '../../proto/apl.js';
import { ActionID, OtherAction, UnitReference, UnitReference_Type as UnitType } from '../../proto/common.js';
import { FeralDruid_Rotation_AplType } from '../../proto/druid.js';
import { ActionId, defaultTargetIcon, getPetIconFromName } from '../../proto_utils/action_id.js';
import { EventID, TypedEvent } from '../../typed_event.js';
import { bucket } from '../../utils.js';
import { BooleanPicker } from '../boolean_picker.js';
import { DropdownPicker, DropdownPickerConfig, DropdownValueConfig, TextDropdownPicker } from '../dropdown_picker.js';
import { Input, InputConfig } from '../input.js';
import { NumberPicker, NumberPickerConfig } from '../number_picker.js';
import { AdaptiveStringPicker } from '../string_picker.js';
import { UnitPicker, UnitPickerConfig, UnitValue } from '../unit_picker.js';

export type ACTION_ID_SET =
	| 'auras'
	| 'stackable_auras'
	| 'icd_auras'
	| 'exclusive_effect_auras'
	| 'spells'
	| 'castable_spells'
	| 'channel_spells'
	| 'dot_spells'
	| 'shield_spells'
	| 'non_instant_spells';

const actionIdSets: Record<
	ACTION_ID_SET,
	{
		defaultLabel: string;
		getActionIDs: (metadata: UnitMetadata) => Promise<Array<DropdownValueConfig<ActionId>>>;
	}
> = {
	auras: {
		defaultLabel: 'Aura',
		getActionIDs: async metadata => {
			return metadata.getAuras().map(actionId => {
				return {
					value: actionId.id,
				};
			});
		},
	},
	stackable_auras: {
		defaultLabel: 'Aura',
		getActionIDs: async metadata => {
			return metadata
				.getAuras()
				.filter(aura => aura.data.maxStacks > 0)
				.map(actionId => {
					return {
						value: actionId.id,
					};
				});
		},
	},
	icd_auras: {
		defaultLabel: 'Aura',
		getActionIDs: async metadata => {
			return metadata
				.getAuras()
				.filter(aura => aura.data.hasIcd)
				.map(actionId => {
					return {
						value: actionId.id,
					};
				});
		},
	},
	exclusive_effect_auras: {
		defaultLabel: 'Aura',
		getActionIDs: async metadata => {
			return metadata
				.getAuras()
				.filter(aura => aura.data.hasExclusiveEffect)
				.map(actionId => {
					return {
						value: actionId.id,
					};
				});
		},
	},
	// Used for non categorized lists
	spells: {
		defaultLabel: 'Spell',
		getActionIDs: async metadata => {
			return metadata
				.getSpells()
				.filter(spell => spell.data.isCastable)
				.map(actionId => {
					return {
						value: actionId.id,
					};
				});
		},
	},
	castable_spells: {
		defaultLabel: 'Spell',
		getActionIDs: async metadata => {
			const castableSpells = metadata.getSpells().filter(spell => spell.data.isCastable);

			// Split up non-cooldowns and cooldowns into separate sections for easier browsing.
			const { spells: spells, cooldowns: cooldowns } = bucket(castableSpells, spell => (spell.data.isMajorCooldown ? 'cooldowns' : 'spells'));

			const placeholders: Array<ActionId> = [ActionId.fromOtherId(OtherAction.OtherActionPotion)];

			return [
				[
					{
						value: ActionId.fromEmpty(),
						headerText: 'Spells',
						submenu: ['Spells'],
					},
				],
				(spells || []).map(actionId => {
					return {
						value: actionId.id,
						submenu: ['Spells'],
						extraCssClasses: actionId.data.prepullOnly
							? ['apl-prepull-actions-only']
							: actionId.data.encounterOnly
							? ['apl-priority-list-only']
							: [],
					};
				}),
				[
					{
						value: ActionId.fromEmpty(),
						headerText: 'Cooldowns',
						submenu: ['Cooldowns'],
					},
				],
				(cooldowns || []).map(actionId => {
					return {
						value: actionId.id,
						submenu: ['Cooldowns'],
						extraCssClasses: actionId.data.prepullOnly
							? ['apl-prepull-actions-only']
							: actionId.data.encounterOnly
							? ['apl-priority-list-only']
							: [],
					};
				}),
				[
					{
						value: ActionId.fromEmpty(),
						headerText: 'Placeholders',
						submenu: ['Placeholders'],
					},
				],
				placeholders.map(actionId => {
					return {
						value: actionId,
						submenu: ['Placeholders'],
						tooltip: 'The Prepull Potion if CurrentTime < 0, or the Combat Potion if combat has started.',
					};
				}),
			].flat();
		},
	},
	non_instant_spells: {
		defaultLabel: 'Non-instant Spell',
		getActionIDs: async metadata => {
			return metadata
				.getSpells()
				.filter(spell => spell.data.isCastable && spell.data.hasCastTime)
				.map(actionId => {
					return {
						value: actionId.id,
					};
				});
		},
	},
	channel_spells: {
		defaultLabel: 'Channeled Spell',
		getActionIDs: async metadata => {
			return metadata
				.getSpells()
				.filter(spell => spell.data.isCastable && spell.data.isChanneled)
				.map(actionId => {
					return {
						value: actionId.id,
					};
				});
		},
	},
	dot_spells: {
		defaultLabel: 'DoT Spell',
		getActionIDs: async metadata => {
			return metadata
				.getSpells()
				.filter(spell => spell.data.hasDot)
				.map(actionId => {
					return {
						value: actionId.id,
					};
				});
		},
	},
	shield_spells: {
		defaultLabel: 'Shield Spell',
		getActionIDs: async metadata => {
			return metadata
				.getSpells()
				.filter(spell => spell.data.hasShield)
				.map(actionId => {
					return {
						value: actionId.id,
					};
				});
		},
	},
};

export type DEFAULT_UNIT_REF = 'self' | 'currentTarget';

export interface APLActionIDPickerConfig<ModObject>
	extends Omit<DropdownPickerConfig<ModObject, ActionID, ActionId>, 'defaultLabel' | 'equals' | 'setOptionContent' | 'values' | 'getValue' | 'setValue'> {
	actionIdSet: ACTION_ID_SET;
	getUnitRef: (player: Player<any>) => UnitReference;
	defaultUnitRef: DEFAULT_UNIT_REF;
	getValue: (obj: ModObject) => ActionID;
	setValue: (eventID: EventID, obj: ModObject, newValue: ActionID) => void;
}

export class APLActionIDPicker extends DropdownPicker<Player<any>, ActionID, ActionId> {
	constructor(parent: HTMLElement, player: Player<any>, config: APLActionIDPickerConfig<Player<any>>) {
		const actionIdSet = actionIdSets[config.actionIdSet];
		super(parent, player, {
			...config,
			sourceToValue: (src: ActionID) => (src ? ActionId.fromProto(src) : ActionId.fromEmpty()),
			valueToSource: (val: ActionId) => val.toProto(),
			defaultLabel: actionIdSet.defaultLabel,
			equals: (a, b) => (a == null) == (b == null) && (!a || a.equals(b!)),
			setOptionContent: (button, valueConfig) => {
				const actionId = valueConfig.value;

				const iconElem = document.createElement('a');
				iconElem.classList.add('apl-actionid-item-icon');
				iconElem.dataset.whtticon = 'false';
				iconElem.classList.add('apl-actionid-item-icon');
				actionId.setBackgroundAndHref(iconElem);
				button.appendChild(iconElem);

				const textElem = document.createTextNode(actionId.name);
				button.appendChild(textElem);
			},
			createMissingValue: value => {
				if (value.anyId() == 0) {
					return new Promise<DropdownValueConfig<ActionId>>(() => {
						value: actionIdSet.defaultLabel;
					});
				}

				return value.fill().then(filledId => {
					return {
						value: filledId,
					};
				});
			},
			values: [],
		});

		const getUnitRef = config.getUnitRef;
		const defaultRef =
			config.defaultUnitRef == 'self' ? UnitReference.create({ type: UnitType.Self }) : UnitReference.create({ type: UnitType.CurrentTarget });
		const getActionIDs = actionIdSet.getActionIDs;
		const updateValues = async () => {
			const unitRef = getUnitRef(player);
			const metadata = player.sim.getUnitMetadata(unitRef, player, defaultRef);
			if (metadata) {
				const values = await getActionIDs(metadata);
				this.setOptions(values);
			}
		};
		updateValues();
		TypedEvent.onAny([player.sim.unitMetadataEmitter, player.rotationChangeEmitter]).on(updateValues);
	}
}

export type UNIT_SET = 'aura_sources' | 'aura_sources_targets_first' | 'targets';

const unitSets: Record<
	UNIT_SET,
	{
		// Uses target icon by default instead of person icon. This should be set to true for inputs that default to CurrentTarget.
		targetUI?: boolean;
		getUnits: (player: Player<any>) => Array<UnitReference | undefined>;
	}
> = {
	aura_sources: {
		getUnits: player => {
			return [
				undefined,
				player
					.getPetMetadatas()
					.asList()
					.map((petMetadata, i) => UnitReference.create({ type: UnitType.Pet, index: i, owner: UnitReference.create({ type: UnitType.Self }) })),
				UnitReference.create({ type: UnitType.CurrentTarget }),
				player.sim.encounter.targetsMetadata.asList().map((targetMetadata, i) => UnitReference.create({ type: UnitType.Target, index: i })),
			].flat();
		},
	},
	aura_sources_targets_first: {
		targetUI: true,
		getUnits: player => {
			return [
				undefined,
				player.sim.encounter.targetsMetadata.asList().map((targetMetadata, i) => UnitReference.create({ type: UnitType.Target, index: i })),
				UnitReference.create({ type: UnitType.Self }),
				player
					.getPetMetadatas()
					.asList()
					.map((petMetadata, i) => UnitReference.create({ type: UnitType.Pet, index: i, owner: UnitReference.create({ type: UnitType.Self }) })),
			].flat();
		},
	},
	targets: {
		targetUI: true,
		getUnits: player => {
			return [
				undefined,
				player.sim.encounter.targetsMetadata.asList().map((targetMetadata, i) => UnitReference.create({ type: UnitType.Target, index: i })),
			].flat();
		},
	},
};

export interface APLUnitPickerConfig extends Omit<UnitPickerConfig<Player<any>>, 'values'> {
	unitSet: UNIT_SET;
}

export class APLUnitPicker extends UnitPicker<Player<any>> {
	private readonly unitSet: UNIT_SET;

	constructor(parent: HTMLElement, player: Player<any>, config: APLUnitPickerConfig) {
		const targetUI = !!unitSets[config.unitSet].targetUI;
		super(parent, player, {
			...config,
			sourceToValue: (src: UnitReference | undefined) => APLUnitPicker.refToValue(src, player, targetUI),
			valueToSource: (val: UnitValue) => val.value,
			values: [],
			hideLabelWhenDefaultSelected: true,
		});
		this.unitSet = config.unitSet;
		this.rootElem.classList.add('apl-unit-picker');

		this.updateValues();
		player.sim.unitMetadataEmitter.on(() => this.updateValues());
	}

	private static refToValue(ref: UnitReference | undefined, thisPlayer: Player<any>, targetUI: boolean | undefined): UnitValue {
		if (!ref || ref.type == UnitType.Unknown) {
			return {
				value: ref,
				iconUrl: targetUI ? 'fa-bullseye' : 'fa-user',
				text: targetUI ? 'Current Target' : 'Self',
			};
		} else if (ref.type == UnitType.Self) {
			return {
				value: ref,
				iconUrl: 'fa-user',
				text: 'Self',
			};
		} else if (ref.type == UnitType.CurrentTarget) {
			return {
				value: ref,
				iconUrl: 'fa-bullseye',
				text: 'Current Target',
			};
		} else if (ref.type == UnitType.Player) {
			const player = thisPlayer.sim.raid.getPlayer(ref.index);
			if (player) {
				return {
					value: ref,
					//color: player.getClassColor(),
					iconUrl: player.getSpecIcon(),
					text: `Player ${ref.index + 1}`,
				};
			}
		} else if (ref.type == UnitType.Target) {
			const targetMetadata = thisPlayer.sim.encounter.targetsMetadata.asList()[ref.index];
			if (targetMetadata) {
				return {
					value: ref,
					iconUrl: defaultTargetIcon,
					text: `Target ${ref.index + 1}`,
				};
			}
		} else if (ref.type == UnitType.Pet) {
			const petMetadata = thisPlayer.sim.getUnitMetadata(ref, thisPlayer, UnitReference.create({ type: UnitType.Self }));
			let name = `Pet ${ref.index + 1}`;
			let icon: string | ActionId = 'fa-paw';
			if (petMetadata) {
				const petName = petMetadata.getName();
				if (petName) {
					const rmIdx = petName.indexOf(' - ');
					name = petName.substring(rmIdx + ' - '.length);
					icon = getPetIconFromName(name) || icon;
				}
			}
			return {
				value: ref,
				iconUrl: icon,
				text: name,
			};
		}

		return {
			value: ref,
		};
	}

	private updateValues() {
		const unitSet = unitSets[this.unitSet];
		const values = unitSet.getUnits(this.modObject);

		this.setOptions(
			values.map(v => {
				const valueConfig: DropdownValueConfig<UnitValue> = {
					value: APLUnitPicker.refToValue(v, this.modObject, unitSet.targetUI),
				};
				if (v && v.type == UnitType.Pet) {
					if (unitSet.targetUI) {
						valueConfig.submenu = [APLUnitPicker.refToValue(v.owner!, this.modObject, unitSet.targetUI)];
					} else {
						valueConfig.submenu = [APLUnitPicker.refToValue(undefined, this.modObject, unitSet.targetUI)];
					}
				}
				return valueConfig;
			}),
		);
	}
}

type APLPickerBuilderFieldFactory<F> = (
	parent: HTMLElement,
	player: Player<any>,
	config: InputConfig<Player<any>, F>,
	getParentValue: () => any,
) => Input<Player<any>, F>;

export interface APLPickerBuilderFieldConfig<T, F extends keyof T> {
	field: F;
	newValue: () => T[F];
	factory: APLPickerBuilderFieldFactory<T[F]>;

	label?: string;
	labelTooltip?: string;
}

export interface APLPickerBuilderConfig<T> extends InputConfig<Player<any>, T> {
	newValue: () => T;
	fields: Array<APLPickerBuilderFieldConfig<T, any>>;
}

export interface APLPickerBuilderField<T, F extends keyof T> extends APLPickerBuilderFieldConfig<T, F> {
	picker: Input<Player<any>, T[F]>;
}

export class APLPickerBuilder<T> extends Input<Player<any>, T> {
	private readonly config: APLPickerBuilderConfig<T>;
	private readonly fieldPickers: Array<APLPickerBuilderField<T, any>>;

	constructor(parent: HTMLElement, modObject: Player<any>, config: APLPickerBuilderConfig<T>) {
		super(parent, 'apl-picker-builder-root', modObject, config);
		this.config = config;

		this.fieldPickers = config.fields.map(fieldConfig => APLPickerBuilder.makeFieldPicker(this, fieldConfig));

		this.init();
	}

	private static makeFieldPicker<T, F extends keyof T>(
		builder: APLPickerBuilder<T>,
		fieldConfig: APLPickerBuilderFieldConfig<T, F>,
	): APLPickerBuilderField<T, F> {
		const field: F = fieldConfig.field;
		const picker = fieldConfig.factory(
			builder.rootElem,
			builder.modObject,
			{
				label: fieldConfig.label,
				labelTooltip: fieldConfig.labelTooltip,
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
			},
			() => builder.getSourceValue(),
		);

		if (field === 'vals' || field === 'actions') {
			picker.rootElem.classList.add('apl-picker-builder-multi');
		}

		return {
			...fieldConfig,
			picker: picker,
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

export function actionIdFieldConfig(
	field: string,
	actionIdSet: ACTION_ID_SET,
	unitRefField?: string,
	defaultUnitRef?: DEFAULT_UNIT_REF,
	options?: Partial<APLPickerBuilderFieldConfig<any, any>>,
): APLPickerBuilderFieldConfig<any, any> {
	return {
		field: field,
		newValue: () => ActionID.create(),
		factory: (parent, player, config, getParentValue) =>
			new APLActionIDPicker(parent, player, {
				...config,
				actionIdSet: actionIdSet,
				getUnitRef: () => (unitRefField ? getParentValue()[unitRefField] : UnitReference.create()),
				defaultUnitRef: defaultUnitRef || 'self',
			}),
		...(options || {}),
	};
}

export function unitFieldConfig(
	field: string,
	unitSet: UNIT_SET,
	options?: Partial<APLPickerBuilderFieldConfig<any, any>>,
): APLPickerBuilderFieldConfig<any, any> {
	return {
		field: field,
		newValue: () => undefined,
		factory: (parent, player, config) =>
			new APLUnitPicker(parent, player, {
				...config,
				unitSet: unitSet,
			}),
		...(options || {}),
	};
}

export function booleanFieldConfig(
	field: string,
	label?: string,
	options?: Partial<APLPickerBuilderFieldConfig<any, any>>,
): APLPickerBuilderFieldConfig<any, any> {
	return {
		field: field,
		newValue: () => false,
		factory: (parent, player, config) => {
			config.extraCssClasses = ['input-inline'].concat(config.extraCssClasses || []);
			return new BooleanPicker(parent, player, config);
		},
		...(options || {}),
		label: label,
	};
}

export function numberFieldConfig(
	field: string,
	float: boolean,
	options?: Partial<APLPickerBuilderFieldConfig<any, any>>,
): APLPickerBuilderFieldConfig<any, any> {
	return {
		field: field,
		newValue: () => 0,
		factory: (parent, player, config) => {
			const numberPickerConfig = config as NumberPickerConfig<Player<any>>;
			numberPickerConfig.float = float;
			numberPickerConfig.extraCssClasses = ['input-inline'].concat(config.extraCssClasses || []);
			return new NumberPicker(parent, player, numberPickerConfig);
		},
		...(options || {}),
	};
}

export function stringFieldConfig(field: string, options?: Partial<APLPickerBuilderFieldConfig<any, any>>): APLPickerBuilderFieldConfig<any, any> {
	return {
		field: field,
		newValue: () => '',
		factory: (parent, player, config) => {
			config.extraCssClasses = ['input-inline'].concat(config.extraCssClasses || []);
			return new AdaptiveStringPicker(parent, player, config);
		},
		...(options || {}),
	};
}

export function runeTypeFieldConfig(field: string, includeDeath: boolean): APLPickerBuilderFieldConfig<any, any> {
	const values = [
		{ value: APLValueRuneType.RuneBlood, label: 'Blood' },
		{ value: APLValueRuneType.RuneFrost, label: 'Frost' },
		{ value: APLValueRuneType.RuneUnholy, label: 'Unholy' },
	];

	if (includeDeath) {
		values.push({ value: APLValueRuneType.RuneDeath, label: 'Death' });
	}

	return {
		field: field,
		newValue: () => APLValueRuneType.RuneBlood,
		factory: (parent, player, config) =>
			new TextDropdownPicker(parent, player, {
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
		factory: (parent, player, config) =>
			new TextDropdownPicker(parent, player, {
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

export function rotationTypeFieldConfig(field: string): APLPickerBuilderFieldConfig<any, any> {
	const values = [
		{ value: FeralDruid_Rotation_AplType.SingleTarget, label: 'Single Target' },
		{ value: FeralDruid_Rotation_AplType.Aoe, label: 'AOE' },
	];

	return {
		field: field,
		label: 'Type',
		newValue: () => FeralDruid_Rotation_AplType.SingleTarget,
		factory: (parent, player, config) =>
			new TextDropdownPicker(parent, player, {
				...config,
				defaultLabel: 'Single Target',
				equals: (a, b) => a == b,
				values: values,
			}),
	};
}

export function aplInputBuilder<T>(
	newValue: () => T,
	fields: Array<APLPickerBuilderFieldConfig<T, keyof T>>,
): (parent: HTMLElement, player: Player<any>, config: InputConfig<Player<any>, T>) => Input<Player<any>, T> {
	return (parent, player, config) => {
		return new APLPickerBuilder(parent, player, {
			...config,
			newValue: newValue,
			fields: fields,
		});
	};
}
