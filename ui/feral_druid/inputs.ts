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
	FeralDruid_Rotation_BiteModeType as BiteModeType,
	FeralDruid_Options as DruidOptions,
	FeralDruid_Rotation_BiteModeType
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
		InputHelpers.makeRotationNumberInput<Spec.SpecFeralDruid>({
			fieldName: 'hotUptime',
			label: 'Revitalize Hot Uptime',
			labelTooltip: 'Hot uptime percentage to assume when theorizing energy gains',
			percent: true
		}),
		InputHelpers.makeRotationNumberInput<Spec.SpecFeralDruid>({
			fieldName: 'maxRoarClip',
			label: 'Roar Clip',
			labelTooltip: 'Max seconds to clip roar',
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecFeralDruid>({
			fieldName: 'useBite',
			label: 'Bite during rotation',
			labelTooltip: 'Use bite during rotation rather than just at end',
		}),
		InputHelpers.makeRotationEnumInput<Spec.SpecFeralDruid, BiteModeType>({
			fieldName: 'biteModeType',
			label: 'Bite Mode',
			labelTooltip: 'Underlying "Bite logic" to use',
			values: [
				{ name: 'Emperical', value: BiteModeType.Emperical },
			],
			showWhen: (player: Player<Spec.SpecFeralDruid>) => player.getRotation().useBite == true
		}),
		InputHelpers.makeRotationNumberInput<Spec.SpecFeralDruid>({
			fieldName: 'biteTime',
			label: 'Bite Time',
			labelTooltip: 'Min seconds on Rip/Roar to bite',
			showWhen: (player: Player<Spec.SpecFeralDruid>) => player.getRotation().useBite == true && player.getRotation().biteModeType == BiteModeType.Emperical,
		})
	],
};
