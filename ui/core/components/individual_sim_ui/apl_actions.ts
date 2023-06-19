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

import * as AplHelpers from './apl_helpers.js';

export class APLActionCastSpellPicker extends Input<Player<any>, APLActionCastSpell> {
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

export class APLActionSequencePicker extends Input<Player<any>, APLActionSequence> {

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

export class APLActionWaitPicker extends Input<Player<any>, APLActionWait> {

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