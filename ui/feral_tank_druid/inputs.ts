import { Spec } from '/wotlk/core/proto/common.js';
import { ActionId } from '/wotlk/core/proto_utils/action_id.js';
import { Player } from '/wotlk/core/player.js';
import { Sim } from '/wotlk/core/sim.js';
import { EventID, TypedEvent } from '/wotlk/core/typed_event.js';
import { Target } from '/wotlk/core/target.js';
import { getEnumValues } from '/wotlk/core/utils.js';
import { ItemSlot } from '/wotlk/core/proto/common.js';

import * as InputHelpers from '/wotlk/core/components/input_helpers.js';

import {
	FeralTankDruid,
	FeralTankDruid_Rotation as DruidRotation,
	FeralTankDruid_Rotation_Swipe as Swipe,
	FeralTankDruid_Options as DruidOptions
} from '/wotlk/core/proto/druid.js';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const StartingRage = InputHelpers.makeSpecOptionsNumberInput<Spec.SpecFeralTankDruid>({
	fieldName: 'startingRage',
	label: 'Starting Rage',
	labelTooltip: 'Initial rage at the start of each iteration.',
});

export const FeralTankDruidRotationConfig = {
	inputs: [
		InputHelpers.makeRotationNumberInput<Spec.SpecFeralTankDruid>({
			fieldName: 'maulRageThreshold',
			label: 'Maul Threshold',
			labelTooltip: 'Queue Maul when rage is above this value.',
		}),
		InputHelpers.makeRotationEnumInput<Spec.SpecFeralTankDruid, Swipe>({
			fieldName: 'swipe',
			label: 'Swipe',
			values: [
				{ name: 'Never', value: Swipe.SwipeNone },
				{ name: 'With Enough AP', value: Swipe.SwipeWithEnoughAP },
				{ name: 'Spam', value: Swipe.SwipeSpam },
			],
		}),
		InputHelpers.makeRotationNumberInput<Spec.SpecFeralTankDruid>({
			fieldName: 'swipeApThreshold',
			label: 'Swipe AP Threshold',
			labelTooltip: 'Use Swipe when Attack Power is larger than this amount.',
			enableWhen: (player: Player<Spec.SpecFeralTankDruid>) => player.getRotation().swipe == Swipe.SwipeWithEnoughAP,
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecFeralTankDruid>({
			fieldName: 'maintainDemoralizingRoar',
			label: 'Maintain Demo Roar',
			labelTooltip: 'Keep Demoralizing Roar active on the primary target. If a stronger debuff is active, will not cast.',
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecFeralTankDruid>({
			fieldName: 'maintainFaerieFire',
			label: 'Maintain Faerie Fire',
			labelTooltip: 'Keep Faerie Fire active on the primary target. If a stronger debuff is active, will not cast.',
		}),
	],
};
