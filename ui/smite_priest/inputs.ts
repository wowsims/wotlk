import { SmitePriest_Rotation_RotationType as RotationType } from '/wotlk/core/proto/priest.js';
import { Race, RaidTarget } from '/wotlk/core/proto/common.js';
import { Spec } from '/wotlk/core/proto/common.js';
import { NO_TARGET } from '/wotlk/core/proto_utils/utils.js';
import { ActionId } from '/wotlk/core/proto_utils/action_id.js';
import { Player } from '/wotlk/core/player.js';
import { Sim } from '/wotlk/core/sim.js';
import { IndividualSimUI } from '/wotlk/core/individual_sim_ui.js';
import { Target } from '/wotlk/core/target.js';
import { EventID, TypedEvent } from '/wotlk/core/typed_event.js';

import * as InputHelpers from '/wotlk/core/components/input_helpers.js';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const SelfPowerInfusion = InputHelpers.makeSpecOptionsBooleanIconInput<Spec.SpecSmitePriest>({
	fieldName: 'powerInfusionTarget',
	id: ActionId.fromSpellId(10060),
	extraCssClasses: [
		'within-raid-sim-hide',
	],
	getValue: (player: Player<Spec.SpecSmitePriest>) => player.getSpecOptions().powerInfusionTarget?.targetIndex != NO_TARGET,
	setValue: (eventID: EventID, player: Player<Spec.SpecSmitePriest>, newValue: boolean) => {
		const newOptions = player.getSpecOptions();
		newOptions.powerInfusionTarget = RaidTarget.create({
			targetIndex: newValue ? 0 : NO_TARGET,
		});
		player.setSpecOptions(eventID, newOptions);
	},
});

export const SmitePriestRotationConfig = {
	inputs: [
		InputHelpers.makeRotationEnumInput<Spec.SpecSmitePriest, RotationType>({
			fieldName: 'rotationType',
			label: 'Rotation Type',
			labelTooltip: 'Choose whether to weave optionally weave holy fire for increase Shadow Word: Pain uptime',
			values: [
				{ name: 'Basic', value: RotationType.Basic },
				{ name: 'HF Weave', value: RotationType.HolyFireWeave },
			],
		}),
		InputHelpers.makeSpecOptionsBooleanInput<Spec.SpecSmitePriest>({
			fieldName: 'useShadowfiend',
			label: 'Use Shadowfiend',
			labelTooltip: 'Use Shadowfiend when low mana and off CD.',
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecSmitePriest>({
			fieldName: 'useMindBlast',
			label: 'Use Mind Blast',
			labelTooltip: 'Use Mind Blast whenever off CD.',
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecSmitePriest>({
			fieldName: 'useShadowWordDeath',
			label: 'Use Shadow Word: Death',
			labelTooltip: 'Use Shadow Word: Death whenever off CD.',
		}),
	],
};
