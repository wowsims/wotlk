import { Spec } from '../core/proto/common.js';
import {
	ShadowPriest_Options_Armor as Armor,
	ShadowPriest_Rotation_RotationType as RotationType,
	ShadowPriest_Rotation_PreCastOption as precastType
} from '../core/proto/priest.js';
import { ActionId } from '../core/proto_utils/action_id.js';

import * as InputHelpers from '../core/components/input_helpers.js';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const ArmorInput = InputHelpers.makeSpecOptionsEnumIconInput<Spec.SpecShadowPriest, Armor>({
	fieldName: 'armor',
	values: [
		{ value: Armor.NoArmor, tooltip: 'No Inner Fire' },
		{ actionId: ActionId.fromSpellId(10952), value: Armor.InnerFire },
	],
});

export const MindBlastInput = InputHelpers.makeSpecOptionsBooleanIconInput<Spec.SpecShadowPriest>({
	fieldName: 'useMindBlast',
	actionId: ActionId.fromSpellId(10947),
});

export const ShadowPriestRotationConfig = {
	inputs: [
		InputHelpers.makeRotationEnumInput<Spec.SpecShadowPriest>({
			fieldName: 'rotationType',
			label: 'Rotation Type',
			labelTooltip: 'Choose how to clip your mindflay. Basic will never clip.',
			values: [
				{ name: 'Basic', value: RotationType.Basic },
			],
		}),
		InputHelpers.makeRotationEnumInput<Spec.SpecShadowPriest>({
			fieldName: 'precastType',
			label: 'PreCast Spell',
			labelTooltip: 'Choose which spell you want to Precast',
			values: [
				{ name: "None", value: precastType.Nothing },
				{ name: 'Mind Blast', value: precastType.PrecastMb },
			],
		}),
	],
};
