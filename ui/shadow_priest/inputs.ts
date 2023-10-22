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
		{ actionId: ActionId.fromSpellId(48168), value: Armor.InnerFire },
	],
});

export const MindBlastInput = InputHelpers.makeSpecOptionsBooleanIconInput<Spec.SpecShadowPriest>({
	fieldName: 'useMindBlast',
	id: ActionId.fromSpellId(48127),
});

export const ShadowWordDeathInput = InputHelpers.makeSpecOptionsBooleanIconInput<Spec.SpecShadowPriest>({
	fieldName: 'useShadowWordDeath',
	id: ActionId.fromSpellId(48158),
});

export const ShadowfiendInput = InputHelpers.makeSpecOptionsBooleanIconInput<Spec.SpecShadowPriest>({
	fieldName: 'useShadowfiend',
	id: ActionId.fromSpellId(34433),
});

export const ShadowPriestRotationConfig = {
	inputs: [
		InputHelpers.makeRotationEnumInput<Spec.SpecShadowPriest, RotationType>({
			fieldName: 'rotationType',
			label: 'Rotation Type',
			labelTooltip: 'Choose how to clip your mindflay. Basic will never clip. Clipping will clip for other spells and use a 2xMF2 when there is time for 4 ticks. Ideal will evaluate the DPS gain of every action to determine MF actions.',
			values: [
				//{ name: 'Basic', value: RotationType.Basic },
				//{ name: 'Clipping', value: RotationType.Clipping },
				{ name: 'Ideal', value: RotationType.Ideal },
				{ name: 'AoE', value: RotationType.AoE },
			],
		}),
		InputHelpers.makeRotationEnumInput<Spec.SpecShadowPriest, precastType>({
			fieldName: 'precastType',
			label: 'PreCast Spell',
			labelTooltip: 'Choose which spell you want to Precast',
			values: [
				{ name: "None", value: precastType.Nothing },
				{ name: 'Vampiric Touch', value: precastType.PrecastVt },
				{ name: 'Mind Blast', value: precastType.PrecastMb },
			],
		}),
	],
};
