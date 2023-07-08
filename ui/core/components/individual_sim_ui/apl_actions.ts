import {
	APLAction,
	APLActionCastSpell,
	APLActionSequence,
	APLActionResetSequence,
	APLActionStrictSequence,
	APLActionMultidot,
	APLActionAutocastOtherCooldowns,
	APLActionWait,
	APLValue,
} from '../../proto/apl.js';

import { EventID } from '../../typed_event.js';
import { Input, InputConfig } from '../input.js';
import { Player } from '../../player.js';
import { TextDropdownPicker } from '../dropdown_picker.js';
import { ListItemPickerConfig, ListPicker } from '../list_picker.js';

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

	private readonly conditionPicker: AplValues.APLValuePicker;

	constructor(parent: HTMLElement, player: Player<any>, config: APLActionPickerConfig) {
		super(parent, 'apl-action-picker-root', player, config);

		this.conditionPicker = new AplValues.APLValuePicker(this.rootElem, this.modObject, {
			label: 'If:',
			changedEvent: (player: Player<any>) => player.rotationChangeEmitter,
			getValue: (player: Player<any>) => this.getSourceValue().condition,
			setValue: (eventID: EventID, player: Player<any>, newValue: APLValue|undefined) => {
				this.getSourceValue().condition = newValue;
				player.rotationChangeEmitter.emit(eventID);
			},
		});
		this.conditionPicker.rootElem.classList.add('apl-action-condition');

		this.actionDiv = document.createElement('div');
		this.actionDiv.classList.add('apl-action-picker-action');
		this.rootElem.appendChild(this.actionDiv);

		const allActionTypes = Object.keys(actionTypeFactories) as Array<NonNullable<APLActionType>>;
		this.typePicker = new TextDropdownPicker(this.actionDiv, player, {
            defaultLabel: 'Action',
			values: allActionTypes.map(actionType => {
				const factory = actionTypeFactories[actionType];
				return {
					value: actionType,
					label: factory.label,
					submenu: factory.submenu,
					tooltip: factory.fullDescription ? `<p>${factory.shortDescription}</p> ${factory.fullDescription}` : factory.shortDescription,
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
		this.actionPicker.rootElem.classList.add('apl-action-' + newActionType);
	}
}

type ActionTypeConfig<T> = {
	label: string,
	submenu?: Array<string>,
	shortDescription: string,
	fullDescription?: string,
	newValue: () => T,
	factory: (parent: HTMLElement, player: Player<any>, config: InputConfig<Player<any>, T>) => Input<Player<any>, T>,
};

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
}): ActionTypeConfig<T> {
	return {
		label: config.label,
		submenu: config.submenu,
		shortDescription: config.shortDescription,
		fullDescription: config.fullDescription,
		newValue: config.newValue,
		factory: AplHelpers.aplInputBuilder(config.newValue, config.fields),
	};
}

export const actionTypeFactories: Record<NonNullable<APLActionType>, ActionTypeConfig<any>> = {
	['castSpell']: inputBuilder({
		label: 'Cast',
		shortDescription: 'Casts the spell if possible, i.e. resource/cooldown/GCD/etc requirements are all met.',
		newValue: APLActionCastSpell.create,
		fields: [
			AplHelpers.actionIdFieldConfig('spellId', 'castable_spells'),
		],
	}),
	['multidot']: inputBuilder({
		label: 'Multi Dot',
		shortDescription: 'Keeps a DoT active on multiple targets by casting the specified spell.',
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
			AplHelpers.actionIdFieldConfig('spellId', 'dot_spells'),
			AplHelpers.numberFieldConfig('maxDots', {
				label: 'Max Dots',
				labelTooltip: 'Maximum number of DoTs to simultaneously apply.',
			}),
			AplValues.valueFieldConfig('maxOverlap', {
				label: 'Overlap',
				labelTooltip: 'Maximum amount of time before a DoT expires when it may be refreshed.',
			}),
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
		newValue: APLActionStrictSequence.create,
		fields: [
			actionListFieldConfig('actions'),
		],
	}),
	['autocastOtherCooldowns']: inputBuilder({
		label: 'Autocast Other Cooldowns',
		submenu: ['Misc'],
		shortDescription: 'Auto-casts cooldowns as soon as they are ready.',
		fullDescription: `
			<ul>
				<li>Does not auto-cast cooldowns which are already controlled by other actions in the priority list.</li>
				<li>Cooldowns are usually cast immediately upon becoming ready, but there are some basic smart checks in place, e.g. don't use Mana CDs when near full mana.</li>
			</ul>
		`,
		newValue: APLActionAutocastOtherCooldowns.create,
		fields: [],
	}),
	['wait']: inputBuilder({
		label: 'Wait',
		submenu: ['Misc'],
		shortDescription: 'Pauses the GCD for a specified amount of time.',
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
};