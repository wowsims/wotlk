import { Spec } from '../core/proto/common.js';
import { ActionId } from '../core/proto_utils/action_id.js';

import {
	WarriorShout
} from '../core/proto/warrior.js';

import * as InputHelpers from '../core/components/input_helpers.js';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const Recklessness = InputHelpers.makeSpecOptionsBooleanIconInput<Spec.SpecWarrior>({
	fieldName: 'useRecklessness',
	actionId: ActionId.fromSpellId(1719),
});

export const StartingRage = InputHelpers.makeSpecOptionsNumberInput<Spec.SpecWarrior>({
	fieldName: 'startingRage',
	label: 'Starting Rage',
	labelTooltip: 'Initial rage at the start of each iteration.',
});

export const StanceSnapshot = InputHelpers.makeSpecOptionsBooleanInput<Spec.SpecWarrior>({
	fieldName: 'stanceSnapshot',
	label: 'Stance Snapshot',
	labelTooltip: 'Ability that is cast at the same time as stance swap will benefit from the bonus of the stance before the swap.',
});

export const ShoutPicker = InputHelpers.makeSpecOptionsBooleanIconInput<Spec.SpecWarrior>({
	fieldName: 'shout',
	actionId: ActionId.fromSpellId(6673),
	value: WarriorShout.WarriorShoutBattle,
});
