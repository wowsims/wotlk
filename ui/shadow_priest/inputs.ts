import { IndividualSimUI } from '../core/individual_sim_ui.js';
import { Player } from '../core/player.js';
import { Spec } from '../core/proto/common.js';
import { 
	ShadowPriest_Options_Armor as Armor,
	ShadowPriest_Rotation_RotationType as RotationType
} from '../core/proto/priest.js';
import { EventID } from '../core/typed_event.js';
import { ActionId } from '../core/proto_utils/action_id.js';

import * as InputHelpers from '../core/components/input_helpers.js';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const ArmorInput = InputHelpers.makeSpecOptionsEnumIconInput<Spec.SpecShadowPriest, Armor>({
	fieldName: 'armor',
	values: [
		{ color: 'grey', value: Armor.NoArmor },
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
				{ name: 'Basic', value: RotationType.Basic },
				{ name: 'Clipping', value: RotationType.Clipping },
				{ name: 'Ideal', value: RotationType.Ideal },
				{ name: 'AoE', value: RotationType.AoE },
			],
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecShadowPriest>({
			fieldName: 'precastVt',
			label: 'Precast Vampiric Touch',
			labelTooltip: 'Start fight with VT landing at time 0',
		}),
		InputHelpers.makeRotationNumberInput<Spec.SpecShadowPriest>({
			fieldName: 'latency',
			label: 'Channeling Latency (ms)',
			labelTooltip: 'Latency after a channel that lasts longer than GCD. 0 to disable. Has a minimum value of 100ms if set.',
		}),
	],
};
