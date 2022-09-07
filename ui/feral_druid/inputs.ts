import { RaidTarget } from '../core/proto/common.js';
import { Spec } from '../core/proto/common.js';
import { NO_TARGET } from '../core/proto_utils/utils.js';
import { ActionId } from '../core/proto_utils/action_id.js';
import { Player } from '../core/player.js';
import { EventID } from '../core/typed_event.js';

import * as InputHelpers from '../core/components/input_helpers.js';

import {
	FeralDruid,
	FeralDruid_Rotation as DruidRotation,
	FeralDruid_Rotation_BearweaveType as BearweaveType,
	FeralDruid_Options as DruidOptions
} from '../core/proto/druid.js';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const SelfInnervate = InputHelpers.makeSpecOptionsBooleanIconInput<Spec.SpecFeralDruid>({
	fieldName: 'innervateTarget',
	id: ActionId.fromSpellId(29166),
	extraCssClasses: [
		'within-raid-sim-hide',
	],
	getValue: (player: Player<Spec.SpecFeralDruid>) => player.getSpecOptions().innervateTarget?.targetIndex != NO_TARGET,
	setValue: (eventID: EventID, player: Player<Spec.SpecFeralDruid>, newValue: boolean) => {
		const newOptions = player.getSpecOptions();
		newOptions.innervateTarget = RaidTarget.create({
			targetIndex: newValue ? 0 : NO_TARGET,
		});
		player.setSpecOptions(eventID, newOptions);
	},
});

export const LatencyMs = InputHelpers.makeSpecOptionsNumberInput<Spec.SpecFeralDruid>({
	fieldName: 'latencyMs',
	label: 'Latency',
	labelTooltip: 'Player latency, in milliseconds. Adds a delay to actions that cannot be spell queued.',
});

export const FeralDruidRotationConfig = {
	inputs: [
		InputHelpers.makeRotationBooleanInput<Spec.SpecFeralDruid>({
			fieldName: 'maintainFaerieFire',
			label: 'Maintain Faerie Fire',
			labelTooltip: 'Use Faerie Fire whenever it is not active on the target.',
		}),
		InputHelpers.makeRotationEnumInput<Spec.SpecFeralDruid, BearweaveType>({
			fieldName: 'bearWeaveType',
			label: 'Bearweaving',
			values: [
				{ name: 'None', value: BearweaveType.None },
				{ name: 'Mangle', value: BearweaveType.Mangle },
				{ name: 'Lacerate', value: BearweaveType.Lacerate },
			],
		}),
	],
};
