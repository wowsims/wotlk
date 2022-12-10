import { RaidTarget } from '../core/proto/common.js';
import { Spec } from '../core/proto/common.js';
import { NO_TARGET } from '../core/proto_utils/utils.js';
import { ActionId } from '../core/proto_utils/action_id.js';
import { Player } from '../core/player.js';
import { EventID, TypedEvent } from '../core/typed_event.js';

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

export const PrepopOoc = InputHelpers.makeSpecOptionsBooleanInput<Spec.SpecFeralDruid>({
	fieldName: 'prepopOoc',
	label: 'Pre-pop Clearcasting',
	labelTooltip: 'Start fight with clearcasting active',
	showWhen: (player: Player<Spec.SpecFeralDruid>) => player.getTalents().omenOfClarity,
	changeEmitter: (player: Player<Spec.SpecFeralDruid>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
})

export const PrepopBerserk = InputHelpers.makeSpecOptionsBooleanInput<Spec.SpecFeralDruid>({
	fieldName: 'prePopBerserk',
	label: 'Pre-pop Berserk',
	labelTooltip: 'Pre pop berserk 1 sec before fight',
	showWhen: (player: Player<Spec.SpecFeralDruid>) => player.getTalents().berserk,
	changeEmitter: (player: Player<Spec.SpecFeralDruid>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
})

export const AssumeBleedActive = InputHelpers.makeSpecOptionsBooleanInput<Spec.SpecFeralDruid>({
	fieldName: 'assumeBleedActive',
	label: 'Assume Bleed Always Active',
	labelTooltip: 'Assume bleed always exists for \'Rend and Tear\' calculations. Otherwise will only calculate based on own rip/rake/lacerate.',
})

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
			percent: true,
			showWhen: (player: Player<Spec.SpecFeralDruid>) => player.getRotation().useBite == true && player.getRotation().biteModeType == BiteModeType.Analytical,
		}),
		InputHelpers.makeRotationNumberInput<Spec.SpecFeralDruid>({
			fieldName: 'maxRoarOffset',
			label: 'Roar Offset',
			labelTooltip: 'Targeted offset in Rip/Roar timings',
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecFeralDruid>({
			fieldName: 'useBite',
			label: 'Bite during rotation',
			labelTooltip: 'Use bite during rotation rather than just at end',
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecFeralDruid>({
			fieldName: 'useRake',
			label: 'Use Rake',
			labelTooltip: 'Use rake during rotation',
		}),
		// Can be uncommented if/when analytical bite mode is added
		//InputHelpers.makeRotationEnumInput<Spec.SpecFeralDruid, BiteModeType>({
		//	fieldName: 'biteModeType',
		//	label: 'Bite Mode',
		//	labelTooltip: 'Underlying "Bite logic" to use',
		//	values: [
		//		{ name: 'Emperical', value: BiteModeType.Emperical },
		//	],
		//	showWhen: (player: Player<Spec.SpecFeralDruid>) => player.getRotation().useBite == true
		//}),
		InputHelpers.makeRotationNumberInput<Spec.SpecFeralDruid>({
			fieldName: 'biteTime',
			label: 'Bite Time',
			labelTooltip: 'Min seconds on Rip/Roar to bite',
			showWhen: (player: Player<Spec.SpecFeralDruid>) => player.getRotation().useBite == true && player.getRotation().biteModeType == BiteModeType.Emperical,
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecFeralDruid>({
			fieldName: 'snekWeave',
			label: 'Snek Weave',
			labelTooltip: 'Reset swing timer using albino snek, when going from bear to cat',
			showWhen: (player: Player<Spec.SpecFeralDruid>) => player.getRotation().bearWeaveType != BearweaveType.None,
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecFeralDruid>({
			fieldName: 'flowerWeave',
			label: 'Flower Weave',
			labelTooltip: 'Fish for clearcasting during rotation with gotw',
			showWhen: (player: Player<Spec.SpecFeralDruid>) => player.getRotation().bearWeaveType == BearweaveType.None,
		}),
		InputHelpers.makeRotationNumberInput<Spec.SpecFeralDruid>({
			fieldName: 'raidTargets',
			label: 'GotW Raid Targets',
			labelTooltip: 'Raid size to assume for clearcast proc chance (can include pets as well, so 25 man raid potentically can be ~30)',
			showWhen: (player: Player<Spec.SpecFeralDruid>) => player.getRotation().bearWeaveType == BearweaveType.None && player.getRotation().flowerWeave == true,
		}),
	],
};
