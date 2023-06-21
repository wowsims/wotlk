import {
	APLListItem,
	APLAction,
	APLActionCastSpell,
	APLActionSequence,
	APLActionWait,
} from '../../proto/apl.js';

import { ActionID, Spec } from '../../proto/common.js';
import { EventID } from '../../typed_event.js';
import { Input, InputConfig } from '../input.js';
import { ActionId } from '../../proto_utils/action_id.js';
import { Player } from '../../player.js';
import { stringComparator } from '../../utils.js';
import { TextDropdownPicker } from '../dropdown_picker.js';

import * as AplHelpers from './apl_helpers.js';

class APLActionCastSpellPicker extends Input<Player<any>, APLActionCastSpell> {
	private readonly spellIdPicker: AplHelpers.APLActionIDPicker;

	constructor(parent: HTMLElement, player: Player<any>, config: InputConfig<Player<any>, APLActionCastSpell>) {
		super(parent, 'apl-action-cast-spell-picker-root', player, config);

        this.spellIdPicker = new AplHelpers.APLActionIDPicker(this.rootElem, player, {
            defaultLabel: 'Spell',
			values: [],
            changedEvent: (player: Player<any>) => player.rotationChangeEmitter,
            getValue: () => ActionId.fromProto(this.getSourceValue().spellId || ActionID.create()),
            setValue: (eventID: EventID, player: Player<any>, newValue: ActionId) => {
                this.getSourceValue().spellId = newValue.toProto();
				player.rotationChangeEmitter.emit(eventID);
            },
        });

		this.init();

		const updateValues = async () => {
			const playerStats = player.getCurrentStats();
			const spellPromises = Promise.all(playerStats.spells.map(spell => ActionId.fromProto(spell).fill()));
			const cooldownPromises = Promise.all(playerStats.cooldowns.map(cd => ActionId.fromProto(cd).fill()));

			let [spells, cooldowns] = await Promise.all([spellPromises, cooldownPromises]);
            spells = spells.sort((a, b) => stringComparator(a.name, b.name))
            cooldowns = cooldowns.sort((a, b) => stringComparator(a.name, b.name))

			const values = [...spells, ...cooldowns].map(actionId => {
				return {
					value: actionId,
				};
			});
            this.spellIdPicker.setOptions(values);
		};
        updateValues();
        player.currentStatsEmitter.on(updateValues);
	}

	getInputElem(): HTMLElement | null {
		return this.rootElem;
	}

    getInputValue(): APLActionCastSpell {
        return APLActionCastSpell.create({
			spellId: this.spellIdPicker.getInputValue(),
		})
    }

	setInputValue(newValue: APLActionCastSpell) {
		if (!newValue) {
			return;
		}
		this.spellIdPicker.setInputValue(ActionId.fromProto(newValue.spellId || ActionID.create()));
	}
}

class APLActionSequencePicker extends Input<Player<any>, APLActionSequence> {

	constructor(parent: HTMLElement, player: Player<any>, config: InputConfig<Player<any>, APLActionSequence>) {
		super(parent, 'apl-action-sequence-picker-root', player, config);
		this.init();
    }

	getInputElem(): HTMLElement | null {
		return this.rootElem;
	}

    getInputValue(): APLActionSequence {
        return APLActionSequence.create({
		})
    }

	setInputValue(newValue: APLActionSequence) {
		if (!newValue) {
			return;
		}
	}
}

class APLActionWaitPicker extends Input<Player<any>, APLActionWait> {

	constructor(parent: HTMLElement, player: Player<any>, config: InputConfig<Player<any>, APLActionWait>) {
		super(parent, 'apl-action-wait-picker-root', player, config);
		this.init();
    }

	getInputElem(): HTMLElement | null {
		return this.rootElem;
	}

    getInputValue(): APLActionWait {
        return APLActionWait.create({
		})
    }

	setInputValue(newValue: APLActionWait) {
		if (!newValue) {
			return;
		}
	}
}

export interface APLActionPickerConfig extends InputConfig<Player<any>, APLAction> {
}

export type APLActionType = APLAction['action']['oneofKind'];

export class APLActionPicker extends Input<Player<any>, APLAction> {

	private typePicker: TextDropdownPicker<Player<any>, APLActionType>;

	private currentType: APLActionType;
	private actionPicker: Input<Player<any>, any>|null;

	constructor(parent: HTMLElement, player: Player<any>, config: APLActionPickerConfig) {
		super(parent, 'apl-action-picker-root', player, config);

		const allActionTypes = Object.keys(APLActionPicker.actionTypeFactories) as Array<NonNullable<APLActionType>>;
		this.typePicker = new TextDropdownPicker(this.rootElem, player, {
            defaultLabel: 'Action',
			values: allActionTypes.map(actionType => {
				return {
					value: actionType,
					label: APLActionPicker.actionTypeFactories[actionType].label,
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
					const factory = APLActionPicker.actionTypeFactories[newValue];
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
			action: {
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

	setInputValue(newValue: APLAction) {
		if (!newValue) {
			return;
		}

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

		const factory = APLActionPicker.actionTypeFactories[newActionType];
		this.actionPicker = new factory.factory(this.rootElem, this.modObject, {
			changedEvent: (player: Player<any>) => player.rotationChangeEmitter,
			getValue: () => (this.getSourceValue().action as any)[newActionType] || factory.newValue(),
			setValue: (eventID: EventID, player: Player<any>, newValue: any) => {
				(this.getSourceValue().action as any)[newActionType] = newValue;
				player.rotationChangeEmitter.emit(eventID);
			},
		});
	}

	private static actionTypeFactories: Record<NonNullable<APLActionType>, { label: string, newValue: () => object, factory: new (parent: HTMLElement, player: Player<any>, config: InputConfig<Player<any>, any>) => Input<Player<any>, any> }>  = {
		['castSpell']: { label: 'Cast', newValue: APLActionCastSpell.create, factory: APLActionCastSpellPicker },
		['sequence']: { label: 'Sequence', newValue: APLActionSequence.create, factory: APLActionSequencePicker },
		['wait']: { label: 'Wait', newValue: APLActionWait.create, factory: APLActionWaitPicker },
	};
}