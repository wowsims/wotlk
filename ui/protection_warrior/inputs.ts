import { Spec } from '../core/proto/common.js';
import { ActionId } from '../core/proto_utils/action_id.js';

import {
	WarriorShout
} from '../core/proto/warrior.js';

import * as InputHelpers from '../core/components/input_helpers.js';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const StartingRage = InputHelpers.makeSpecOptionsNumberInput<Spec.SpecProtectionWarrior>({
	fieldName: 'startingRage',
	label: 'Starting Rage',
	labelTooltip: 'Initial rage at the start of each iteration.',
});

export const ShoutPicker = InputHelpers.makeSpecOptionsBooleanIconInput<Spec.SpecProtectionWarrior>({
	fieldName: 'shout',
	actionId: ActionId.fromSpellId(47436),
	value: WarriorShout.WarriorShoutBattle,
});
