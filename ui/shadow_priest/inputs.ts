import { IndividualSimUI } from '/wotlk/core/individual_sim_ui.js';
import { Player } from '/wotlk/core/player.js';
import { Spec } from '/wotlk/core/proto/common.js';
import { 
	ShadowPriest_Options_Armor as Armor,
	ShadowPriest_Rotation_RotationType as RotationType
} from '/wotlk/core/proto/priest.js';
import { EventID } from '/wotlk/core/typed_event.js';
import { ActionId } from '/wotlk/core/proto_utils/action_id.js';

import * as InputHelpers from '/wotlk/core/components/input_helpers.js';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const ArmorInput = InputHelpers.makeSpecOptionsEnumIconInput<Spec.SpecShadowPriest, Armor>({
	fieldName: 'armor',
	values: [
		{ color: 'grey', value: Armor.NoArmor },
		{ actionId: ActionId.fromSpellId(48168), value: Armor.InnerFire },
	],
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
			],
		}),
		InputHelpers.makeSpecOptionsBooleanInput<Spec.SpecShadowPriest>({
			fieldName: 'useShadowfiend',
			label: 'Use Shadowfiend',
			labelTooltip: 'Use Shadowfiend when low mana and off CD.',
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
