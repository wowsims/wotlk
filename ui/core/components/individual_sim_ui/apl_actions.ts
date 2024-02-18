import {
	Class,
	Spec,
} from '../../proto/common.js';

import {
	APLAction,

	APLActionCastSpell,
	APLActionChannelSpell,
	APLActionMultidot,
	APLActionMultishield,
	APLActionAutocastOtherCooldowns,

	APLActionWait,
	APLActionWaitUntil,
	APLActionSchedule,

	APLActionSequence,
	APLActionResetSequence,
	APLActionStrictSequence,

	APLActionChangeTarget,
	APLActionActivateAura,
	APLActionCancelAura,
	APLActionTriggerICD,
	APLActionItemSwap,
	APLActionItemSwap_SwapSet as ItemSwapSet,

	APLActionCustomRotation,
	APLActionCatOptimalRotationAction,

	APLValue,
} from '../../proto/apl.js';

import { isHealingSpec } from '../../proto_utils/utils.js';
import { EventID } from '../../typed_event.js';
import { itemSwapEnabledSpecs } from '../../individual_sim_ui.js';
import { Input, InputConfig } from '../input.js';
import { Player } from '../../player.js';
import { TextDropdownPicker } from '../dropdown_picker.js';
import { ListItemPickerConfig, ListPicker } from '../list_picker.js';
import { FeralDruid_Rotation_AplType } from '../../proto/druid.js';

import * as AplHelpers from './apl_helpers.js';
import * as AplValues from './apl_values.js';

export interface APLActionPickerConfig extends InputConfig<Player<any>, APLAction> {
}

export type APLActionKind = APLAction['action']['oneofKind'];
type APLActionImplStruct<F extends APLActionKind> = Extract<APLAction['action'], {oneofKind: F}>;
type APLActionImplTypesUnion = {
	[f in NonNullable<APLActionKind>]: f extends keyof APLActionImplStruct<f> ? APLActionImplStruct<f>[f] : never;
};
export type APLActionImplType = APLActionImplTypesUnion[NonNullable<APLActionKind>]|undefined;

export class APLActionPicker extends Input<Player<any>, APLAction> {

	private kindPicker: TextDropdownPicker<Player<any>, APLActionKind>;

	private readonly actionDiv: HTMLElement;
	private currentKind: APLActionKind;
	private actionPicker: Input<Player<any>, any> | null;

	private readonly conditionPicker: AplValues.APLValuePicker;

	constructor(parent: HTMLElement, player: Player<any>, config: APLActionPickerConfig) {
		super(parent, 'apl-action-picker-root', player, config);

		this.conditionPicker = new AplValues.APLValuePicker(this.rootElem, this.modObject, {
			label: 'If:',
			changedEvent: (player: Player<any>) => player.rotationChangeEmitter,
			getValue: (player: Player<any>) => this.getSourceValue()?.condition,
			setValue: (eventID: EventID, player: Player<any>, newValue: APLValue | undefined) => {
				const srcVal = this.getSourceValue();
				if (srcVal) {
					srcVal.condition = newValue;
					player.rotationChangeEmitter.emit(eventID);
				} else {
					this.setSourceValue(eventID, APLAction.create({
						condition: newValue,
					}));
				}
			},
		});
		this.conditionPicker.rootElem.classList.add('apl-action-condition', 'apl-priority-list-only');

		this.actionDiv = document.createElement('div');
		this.actionDiv.classList.add('apl-action-picker-action');
		this.rootElem.appendChild(this.actionDiv);

		const isPrepull = this.rootElem.closest('.apl-prepull-action-picker') != null;

		const allActionKinds = (Object.keys(actionKindFactories) as Array<NonNullable<APLActionKind>>)
			.filter(actionKind => actionKindFactories[actionKind].includeIf?.(player, isPrepull) ?? true);

		this.kindPicker = new TextDropdownPicker(this.actionDiv, player, {
			defaultLabel: 'Action',
			values: allActionKinds
				.map(actionKind => {
					const factory = actionKindFactories[actionKind];
					return {
						value: actionKind,
						label: factory.label,
						submenu: factory.submenu,
						tooltip: factory.fullDescription ? `<p>${factory.shortDescription}</p> ${factory.fullDescription}` : factory.shortDescription,
					};
				}),
			equals: (a, b) => a == b,
			changedEvent: (player: Player<any>) => player.rotationChangeEmitter,
			getValue: (player: Player<any>) => this.getSourceValue()?.action.oneofKind,
			setValue: (eventID: EventID, player: Player<any>, newKind: APLActionKind) => {
				const sourceValue = this.getSourceValue();
				const oldKind = sourceValue?.action.oneofKind;
				if (oldKind == newKind) {
					return;
				}

				if (newKind) {
					const factory = actionKindFactories[newKind];
					let newSourceValue = this.makeAPLAction(newKind, factory.newValue());
					if (sourceValue) {
						// Some pre-fill logic when swapping kinds.
						if (oldKind && this.actionPicker) {
							if (newKind == 'sequence') {
								if (sourceValue.action.oneofKind == 'strictSequence') {
									(newSourceValue.action as APLActionImplStruct<'sequence'>).sequence.actions = sourceValue.action.strictSequence.actions;
								} else {
									(newSourceValue.action as APLActionImplStruct<'sequence'>).sequence.actions = [this.makeAPLAction(oldKind, this.actionPicker.getInputValue())];
								}
							} else if (newKind == 'strictSequence') {
								if (sourceValue.action.oneofKind == 'sequence') {
									(newSourceValue.action as APLActionImplStruct<'strictSequence'>).strictSequence.actions = sourceValue.action.sequence.actions;
								} else {
									(newSourceValue.action as APLActionImplStruct<'strictSequence'>).strictSequence.actions = [this.makeAPLAction(oldKind, this.actionPicker.getInputValue())];
								}
							} else if (sourceValue.action.oneofKind == 'sequence' && sourceValue.action.sequence.actions?.[0]?.action.oneofKind == newKind) {
								newSourceValue = sourceValue.action.sequence.actions[0];
							} else if (sourceValue.action.oneofKind == 'strictSequence' && sourceValue.action.strictSequence.actions?.[0]?.action.oneofKind == newKind) {
								newSourceValue = sourceValue.action.strictSequence.actions[0];
							}
						}
					}
					if (sourceValue) {
						sourceValue.action = newSourceValue.action;
					} else {
						this.setSourceValue(eventID, newSourceValue);
					}
				} else {
					sourceValue.action = {
						oneofKind: newKind,
					};
				}
				player.rotationChangeEmitter.emit(eventID);
			},
		});

		this.currentKind = undefined;
		this.actionPicker = null;

		this.init();
	}

	getInputElem(): HTMLElement | null {
		return this.rootElem;
	}

	getInputValue(): APLAction {
		const actionKind = this.kindPicker.getInputValue();
		return APLAction.create({
			condition: this.conditionPicker.getInputValue(),
			action: {
				oneofKind: actionKind,
				...((() => {
					const val: any = {};
					if (actionKind && this.actionPicker) {
						val[actionKind] = this.actionPicker.getInputValue();
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

		const newActionKind = newValue.action.oneofKind;
		this.updateActionPicker(newActionKind);

		if (newActionKind) {
			this.actionPicker!.setInputValue((newValue.action as any)[newActionKind]);
		}
	}

	private makeAPLAction<K extends NonNullable<APLActionKind>>(kind: K, implVal: APLActionImplTypesUnion[K]): APLAction {
		if (!kind) {
			return APLAction.create();
		}
		const obj: any = { oneofKind: kind };
		obj[kind] = implVal;
		return APLAction.create({action: obj});
	}

	private updateActionPicker(newActionKind: APLActionKind) {
		const actionKind = this.currentKind;
		if (newActionKind == actionKind) {
			return;
		}
		this.currentKind = newActionKind;

		if (this.actionPicker) {
			this.actionPicker.rootElem.remove();
			this.actionPicker = null;
		}

		if (!newActionKind) {
			return;
		}

		this.kindPicker.setInputValue(newActionKind);

		const factory = actionKindFactories[newActionKind];
		this.actionPicker = factory.factory(this.actionDiv, this.modObject, {
			changedEvent: (player: Player<any>) => player.rotationChangeEmitter,
			getValue: () => (this.getSourceValue()?.action as any)?.[newActionKind] || factory.newValue(),
			setValue: (eventID: EventID, player: Player<any>, newValue: any) => {
				const sourceValue = this.getSourceValue();
				if (sourceValue) {
					(sourceValue?.action as any)[newActionKind] = newValue;
				}
				player.rotationChangeEmitter.emit(eventID);
			},
		});
		this.actionPicker.rootElem.classList.add('apl-action-' + newActionKind);
	}
}

type ActionKindConfig<T> = {
	label: string,
	submenu?: Array<string>,
	shortDescription: string,
	fullDescription?: string,
	includeIf?: (player: Player<any>, isPrepull: boolean) => boolean,
	newValue: () => T,
	factory: (parent: HTMLElement, player: Player<any>, config: InputConfig<Player<any>, T>) => Input<Player<any>, T>,
};

function itemSwapSetFieldConfig(field: string): AplHelpers.APLPickerBuilderFieldConfig<any, any> {
	return {
		field: field,
		newValue: () => ItemSwapSet.Swap1,
		factory: (parent, player, config) => new TextDropdownPicker(parent, player, {
			...config,
			defaultLabel: 'None',
			equals: (a, b) => a == b,
			values: [
				{ value: ItemSwapSet.Main, label: 'Main' },
				{ value: ItemSwapSet.Swap1, label: 'Swapped' },
			],
		}),
	};
}

function actionFieldConfig(field: string): AplHelpers.APLPickerBuilderFieldConfig<any, any> {
	return {
		field: field,
		newValue: APLValue.create,
		factory: (parent, player, config) => new APLActionPicker(parent, player, config),
	};
}

function actionListFieldConfig(field: string): AplHelpers.APLPickerBuilderFieldConfig<any, any> {
	return {
		field: field,
		newValue: () => [],
		factory: (parent, player, config) => new ListPicker<Player<any>, APLAction>(parent, player, {
			...config,
			// Override setValue to replace undefined elements with default messages.
			setValue: (eventID: EventID, player: Player<any>, newValue: Array<APLAction>) => {
				config.setValue(eventID, player, newValue.map(val => val || APLAction.create()));
			},
			itemLabel: 'Action',
			newItem: APLAction.create,
			copyItem: (oldValue: APLAction) => oldValue ? APLAction.clone(oldValue) : oldValue,
			newItemPicker: (parent: HTMLElement, listPicker: ListPicker<Player<any>, APLAction>, index: number, config: ListItemPickerConfig<Player<any>, APLAction>) => new APLActionPicker(parent, player, config),
			allowedActions: ['create', 'delete', 'move'],
			actions: {
				create: {
					useIcon: true,
				}
			}
		}),
	};
}

function inputBuilder<T>(config: {
	label: string,
	submenu?: Array<string>,
	shortDescription: string,
	fullDescription?: string,
	includeIf?: (player: Player<any>, isPrepull: boolean) => boolean,
	newValue: () => T,
	fields: Array<AplHelpers.APLPickerBuilderFieldConfig<T, any>>,
}): ActionKindConfig<T> {
	return {
		label: config.label,
		submenu: config.submenu,
		shortDescription: config.shortDescription,
		fullDescription: config.fullDescription,
		includeIf: config.includeIf,
		newValue: config.newValue,
		factory: AplHelpers.aplInputBuilder(config.newValue, config.fields),
	};
}

const actionKindFactories: {[f in NonNullable<APLActionKind>]: ActionKindConfig<APLActionImplTypesUnion[f]>} = {
	['castSpell']: inputBuilder({
		label: 'Cast',
		shortDescription: 'Casts the spell if possible, i.e. resource/cooldown/GCD/etc requirements are all met.',
		newValue: APLActionCastSpell.create,
		fields: [
			AplHelpers.actionIdFieldConfig('spellId', 'castable_spells', ''),
			AplHelpers.unitFieldConfig('target', 'targets'),
		],
	}),
	['multidot']: inputBuilder({
		label: 'Multi Dot',
		submenu: ['Casting'],
		shortDescription: 'Keeps a DoT active on multiple targets by casting the specified spell.',
		includeIf: (player: Player<any>, isPrepull: boolean) => !isPrepull,
		newValue: () => APLActionMultidot.create({
			maxDots: 3,
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
			AplHelpers.actionIdFieldConfig('spellId', 'dot_spells', ''),
			AplHelpers.numberFieldConfig('maxDots', false, {
				label: 'Max Dots',
				labelTooltip: 'Maximum number of DoTs to simultaneously apply.',
			}),
			AplValues.valueFieldConfig('maxOverlap', {
				label: 'Overlap',
				labelTooltip: 'Maximum amount of time before a DoT expires when it may be refreshed.',
			}),
		],
	}),
	['multishield']: inputBuilder({
		label: 'Multi Shield',
		submenu: ['Casting'],
		shortDescription: 'Keeps a Shield active on multiple targets by casting the specified spell.',
		includeIf: (player: Player<any>, isPrepull: boolean) => !isPrepull && isHealingSpec(player.spec),
		newValue: () => APLActionMultishield.create({
			maxShields: 3,
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
			AplHelpers.actionIdFieldConfig('spellId', 'shield_spells', ''),
			AplHelpers.numberFieldConfig('maxShields', false, {
				label: 'Max Shields',
				labelTooltip: 'Maximum number of Shields to simultaneously apply.',
			}),
			AplValues.valueFieldConfig('maxOverlap', {
				label: 'Overlap',
				labelTooltip: 'Maximum amount of time before a Shield expires when it may be refreshed.',
			}),
		],
	}),
	['channelSpell']: inputBuilder({
		label: 'Channel',
		submenu: ['Casting'],
		shortDescription: 'Channels the spell if possible, i.e. resource/cooldown/GCD/etc requirements are all met.',
		fullDescription: `
			<p>The difference between channeling a spell vs casting the spell is that channels can be interrupted. If the <b>Interrupt If</b> parameter is empty, this action is equivalent to <b>Cast</b>.</p>
			<p>The channel will be interrupted only if all of the following are true:</p>
			<ul>
				<li>Immediately following a tick of the channel</li>
				<li>The <b>Interrupt If</b> condition evaluates to <b>True</b></li>
				<li>Another action in the APL list is available</li>
			</ul>
			<p>Note that if you simply want to allow other actions to interrupt the channel, set <b>Interrupt If</b> to <b>True</b>.</p>
		`,
		newValue: () => APLActionChannelSpell.create({
			interruptIf: {
				value: {
					oneofKind: 'gcdIsReady',
					gcdIsReady: {},
				}
			},
		}),
		fields: [
			AplHelpers.actionIdFieldConfig('spellId', 'channel_spells', ''),
			AplHelpers.unitFieldConfig('target', 'targets'),
			AplValues.valueFieldConfig('interruptIf', {
				label: 'Interrupt If',
				labelTooltip: 'Condition which must be true to allow the channel to be interrupted.',
			}),
			AplHelpers.booleanFieldConfig('allowRecast', 'Recast', {
				labelTooltip: 'If checked, interrupts of this channel will recast the spell.',
			}),
		],
	}),
	['autocastOtherCooldowns']: inputBuilder({
		label: 'Autocast Other Cooldowns',
		submenu: ['Casting'],
		shortDescription: 'Auto-casts cooldowns as soon as they are ready.',
		fullDescription: `
			<ul>
				<li>Does not auto-cast cooldowns which are already controlled by other actions in the priority list.</li>
				<li>Cooldowns are usually cast immediately upon becoming ready, but there are some basic smart checks in place, e.g. don't use Mana CDs when near full mana.</li>
			</ul>
		`,
		includeIf: (player: Player<any>, isPrepull: boolean) => !isPrepull,
		newValue: APLActionAutocastOtherCooldowns.create,
		fields: [],
	}),
	['wait']: inputBuilder({
		label: 'Wait',
		submenu: ['Timing'],
		shortDescription: 'Pauses all APL actions for a specified amount of time.',
		includeIf: (player: Player<any>, isPrepull: boolean) => !isPrepull,
		newValue: () => APLActionWait.create({
			duration: {
				value: {
					oneofKind: 'const',
					const: {
						val: '1000ms',
					},
				},
			},
		}),
		fields: [
			AplValues.valueFieldConfig('duration'),
		],
	}),
	['waitUntil']: inputBuilder({
		label: 'Wait Until',
		submenu: ['Timing'],
		shortDescription: 'Pauses all APL actions until the specified condition is <b>True</b>.',
		includeIf: (player: Player<any>, isPrepull: boolean) => !isPrepull,
		newValue: () => APLActionWaitUntil.create(),
		fields: [
			AplValues.valueFieldConfig('condition'),
		],
	}),
	['schedule']: inputBuilder({
		label: 'Scheduled Action',
		submenu: ['Timing'],
		shortDescription: 'Executes the inner action once at each specified timing.',
		includeIf: (player: Player<any>, isPrepull: boolean) => !isPrepull,
		newValue: () => APLActionSchedule.create({
			schedule: '0s, 60s',
			innerAction: {
				action: {oneofKind: 'castSpell', castSpell: {}},
			},
		}),
		fields: [
			AplHelpers.stringFieldConfig('schedule', {
				label: 'Do At',
				labelTooltip: 'Comma-separated list of timings. The inner action will be performed once at each timing.',
			}),
			actionFieldConfig('innerAction'),
		],
	}),
	['sequence']: inputBuilder({
		label: 'Sequence',
		submenu: ['Sequences'],
		shortDescription: 'A list of sub-actions to execute in the specified order.',
		fullDescription: `
			<p>Once one of the sub-actions has been performed, the next sub-action will not necessarily be immediately executed next. The system will restart at the beginning of the whole actions list (not the sequence). If the sequence is executed again, it will perform the next sub-action.</p>
			<p>When all actions have been performed, the sequence does NOT automatically reset; instead, it will be skipped from now on. Use the <b>Reset Sequence</b> action to reset it, if desired.</p>
		`,
		includeIf: (player: Player<any>, isPrepull: boolean) => !isPrepull,
		newValue: APLActionSequence.create,
		fields: [
			AplHelpers.stringFieldConfig('name'),
			actionListFieldConfig('actions'),
		],
	}),
	['resetSequence']: inputBuilder({
		label: 'Reset Sequence',
		submenu: ['Sequences'],
		shortDescription: 'Restarts a sequence, so that the next time it executes it will perform its first sub-action.',
		fullDescription: `
			<p>Use the <b>name</b> field to refer to the sequence to be reset. The desired sequence must have the same (non-empty) value for its <b>name</b>.</p>
		`,
		includeIf: (player: Player<any>, isPrepull: boolean) => !isPrepull,
		newValue: APLActionResetSequence.create,
		fields: [
			AplHelpers.stringFieldConfig('sequenceName'),
		],
	}),
	['strictSequence']: inputBuilder({
		label: 'Strict Sequence',
		submenu: ['Sequences'],
		shortDescription: 'Like a regular <b>Sequence</b>, except all sub-actions are executed immediately after each other and the sequence resets automatically upon completion.',
		fullDescription: `
			<p>Strict Sequences do not begin unless ALL sub-actions are ready.</p>
		`,
		includeIf: (player: Player<any>, isPrepull: boolean) => !isPrepull,
		newValue: APLActionStrictSequence.create,
		fields: [
			actionListFieldConfig('actions'),
		],
	}),
	['changeTarget']: inputBuilder({
		label: 'Change Target',
		submenu: ['Misc'],
		shortDescription: 'Sets the current target, which is the target of auto attacks and most casts by default.',
		newValue: () => APLActionChangeTarget.create(),
		fields: [
			AplHelpers.unitFieldConfig('newTarget', 'targets'),
		],
	}),
	['activateAura']: inputBuilder({
		label: 'Activate Aura',
		submenu: ['Misc'],
		shortDescription: 'Activates an aura',
		includeIf: (player: Player<any>, isPrepull: boolean) => isPrepull,
		newValue: () => APLActionActivateAura.create(),
		fields: [
			AplHelpers.actionIdFieldConfig('auraId', 'auras'),
		],
	}),
	['cancelAura']: inputBuilder({
		label: 'Cancel Aura',
		submenu: ['Misc'],
		shortDescription: 'Deactivates an aura, equivalent to /cancelaura.',
		newValue: () => APLActionCancelAura.create(),
		fields: [
			AplHelpers.actionIdFieldConfig('auraId', 'auras'),
		],
	}),
	['triggerIcd']: inputBuilder({
		label: 'Trigger ICD',
		submenu: ['Misc'],
		shortDescription: 'Triggers an aura\'s ICD, putting it on cooldown. Example usage would be to desync an ICD cooldown before combat starts.',
		includeIf: (player: Player<any>, isPrepull: boolean) => isPrepull,
		newValue: () => APLActionTriggerICD.create(),
		fields: [
			AplHelpers.actionIdFieldConfig('auraId', 'icd_auras'),
		],
	}),
	['itemSwap']: inputBuilder({
		label: 'Item Swap',
		submenu: ['Misc'],
		shortDescription: 'Swaps items, using the swap set specified in Settings.',
		includeIf: (player: Player<any>, isPrepull: boolean) => itemSwapEnabledSpecs.includes(player.spec),
		newValue: () => APLActionItemSwap.create(),
		fields: [
			itemSwapSetFieldConfig('swapSet'),
		],
	}),

	['customRotation']: inputBuilder({
		label: 'Custom Rotation',
		//submenu: ['Misc'],
		shortDescription: 'INTERNAL ONLY',
		includeIf: (player: Player<any>, isPrepull: boolean) => false, // Never show this, because its internal only.
		newValue: () => APLActionCustomRotation.create(),
		fields: [
		],
	}),

	// Class/spec specific actions
	['catOptimalRotationAction']: inputBuilder({
		label: 'Optimal Rotation Action',
		submenu: ['Feral Druid'],
		shortDescription: 'Executes optimized Feral DPS rotation using hardcoded legacy algorithm.',
		includeIf: (player: Player<any>, isPrepull: boolean) => player.spec == Spec.SpecFeralDruid,
		newValue: () => APLActionCatOptimalRotationAction.create({
			rotationType: FeralDruid_Rotation_AplType.SingleTarget,
			manualParams: true,
			maxFfDelay: 0.1,
			minRoarOffset: 25.0,
			ripLeeway: 4,
			useRake: true,
			useBite: true,
			biteTime: 4.0,
			flowerWeave: false,
		}),
		fields: [
			AplHelpers.rotationTypeFieldConfig('rotationType'),
			AplHelpers.booleanFieldConfig('manualParams', 'Manual Advanced Parameters', {
				'labelTooltip': 'Manually specify advanced parameters, otherwise will use preset defaults.',
			}),
			AplHelpers.numberFieldConfig('maxFfDelay', true, {
				'label': 'Max FF Delay',
				'labelTooltip': 'Max allowed FF delay to fit in damage casts. Ignored if not using manual advanced parameters.',
			}),
			AplHelpers.numberFieldConfig('minRoarOffset', true, {
				'label': 'Roar Offset',
				'labelTooltip': 'Targeted offset in Rip/Roar timings. Ignored for AOE rotation or if not using manual advanced parameters.',
			}),
			AplHelpers.numberFieldConfig('ripLeeway', false, {
				'label': 'Rip Leeway',
				'labelTooltip': 'Rip leeway when optimizing Roar clips. Ignored for AOE rotation or if not using manual advanced parameters.',
			}),
			AplHelpers.booleanFieldConfig('useRake', 'Use Rake', {
				'labelTooltip': 'Use Rake during rotation. Ignored for AOE rotation or if not using manual advanced parameters.',
			}),
			AplHelpers.booleanFieldConfig('useBite', 'Bite during rotation', {
				'labelTooltip': 'Use Bite during rotation rather than exclusively at end of fight. Ignored for AOE rotation or if not using manual advanced parameters.',
			}),
			AplHelpers.numberFieldConfig('biteTime', true, {
				'label': 'Bite Time',
				'labelTooltip': 'Min seconds remaining on Rip/Roar to allow a Bite. Ignored if not Biting during rotation.',
			}),
			AplHelpers.booleanFieldConfig('flowerWeave', 'Flower Weave', {
				'labelTooltip': 'Fish for Clearcasting procs during AOE rotation with GotW. Ignored for Single Target rotation or if not using manual advanced parameters.',
			}),
		],
	}),
};
