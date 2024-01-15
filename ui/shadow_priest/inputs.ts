import { Spec } from '../core/proto/common.js';
import {
	ShadowPriest_Options_Armor as Armor,
} from '../core/proto/priest.js';
import { ActionId } from '../core/proto_utils/action_id.js';

import * as InputHelpers from '../core/components/input_helpers.js';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const ArmorInput = InputHelpers.makeSpecOptionsEnumIconInput<Spec.SpecShadowPriest, Armor>({
	fieldName: 'armor',
	values: [
		{ value: Armor.NoArmor, tooltip: 'No Inner Fire' },
		{ actionId: ActionId.fromSpellId(48168), value: Armor.InnerFire },
	],
});
