import { UnitReference, UnitReference_Type as UnitType } from '../core/proto/common.js';
import { Spec } from '../core/proto/common.js';
import { APLRotation_Type } from '../core/proto/apl.js';
import { ActionId } from '../core/proto_utils/action_id.js';
import { Player } from '../core/player.js';
import { EventID, TypedEvent } from '../core/typed_event.js';

import * as InputHelpers from '../core/components/input_helpers.js';

import {
	FeralDruid_Rotation_AplType as AplType,
	FeralDruid_Rotation_BiteModeType as BiteModeType,
} from '../core/proto/druid.js';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const SelfInnervate = InputHelpers.makeSpecOptionsBooleanIconInput<Spec.SpecFeralDruid>({
	fieldName: 'innervateTarget',
	id: ActionId.fromSpellId(29166),
	extraCssClasses: [
		'within-raid-sim-hide',
	],
	getValue: (player: Player<Spec.SpecFeralDruid>) => player.getSpecOptions().innervateTarget?.type == UnitType.Player,
	setValue: (eventID: EventID, player: Player<Spec.SpecFeralDruid>, newValue: boolean) => {
		const newOptions = player.getSpecOptions();
		newOptions.innervateTarget = UnitReference.create({
			type: newValue ? UnitType.Player : UnitType.Unknown,
			index: 0,
		});
		player.setSpecOptions(eventID, newOptions);
	},
});

export const LatencyMs = InputHelpers.makeSpecOptionsNumberInput<Spec.SpecFeralDruid>({
	fieldName: 'latencyMs',
	label: 'Latency',
	labelTooltip: 'Player latency, in milliseconds. Adds a delay to actions that cannot be spell queued.',
});

export const AssumeBleedActive = InputHelpers.makeSpecOptionsBooleanInput<Spec.SpecFeralDruid>({
	fieldName: 'assumeBleedActive',
	label: 'Assume Bleed Always Active',
	labelTooltip: 'Assume bleed always exists for \'Rend and Tear\' calculations. Otherwise will only calculate based on own rip/rake/lacerate.',
	extraCssClasses: ['within-raid-sim-hide'],
})

function ShouldShowAdvParamST(player: Player<Spec.SpecFeralDruid>): boolean {
	let rot = player.getSimpleRotation();
	return rot.manualParams && rot.rotationType == AplType.SingleTarget;
}

function ShouldShowAdvParamAoe(player: Player<Spec.SpecFeralDruid>): boolean {
	let rot = player.getSimpleRotation();
	return rot.manualParams && rot.rotationType == AplType.Aoe;
}

export const FeralDruidRotationConfig = {
	inputs: [
		InputHelpers.makeRotationEnumInput<Spec.SpecFeralDruid, AplType>({
			fieldName: 'rotationType',
			label: 'Type',
			values: [
				{ name: 'Single Target', value: AplType.SingleTarget },
				{ name: 'AOE', value: AplType.Aoe },
			],
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecFeralDruid>({
			fieldName: 'prePopOoc',
			label: 'Pre-pop Clearcasting',
			labelTooltip: 'Start fight with clearcasting',
			showWhen: (player: Player<Spec.SpecFeralDruid>) => player.getTalents().omenOfClarity,
			changeEmitter: (player: Player<Spec.SpecFeralDruid>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecFeralDruid>({
			fieldName: 'prePopBerserk',
			label: 'Pre-pop Berserk',
			labelTooltip: 'Pre pop berserk 1 sec before fight',
			showWhen: (player: Player<Spec.SpecFeralDruid>) => player.getTalents().berserk,
			changeEmitter: (player: Player<Spec.SpecFeralDruid>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecFeralDruid>({
			fieldName: 'manualParams',
			label: 'Manual Advanced Parameters',
			labelTooltip: 'Manually specify advanced parameters, otherwise will use preset defaults',
		}),
		InputHelpers.makeRotationNumberInput<Spec.SpecFeralDruid>({
			fieldName: 'maxFfDelay',
			label: 'Max FF Delay',
			labelTooltip: 'Max allowed delay to wait for ff to come off CD in seconds',
			float: true,
			positive: true,
			showWhen: (player: Player<Spec.SpecFeralDruid>) => player.getSimpleRotation().manualParams,
		}),
		InputHelpers.makeRotationNumberInput<Spec.SpecFeralDruid>({
			fieldName: 'minRoarOffset',
			label: 'Roar Offset',
			labelTooltip: 'Targeted offset in Rip/Roar timings',
			showWhen: ShouldShowAdvParamST,
		}),
		InputHelpers.makeRotationNumberInput<Spec.SpecFeralDruid>({
			fieldName: 'ripLeeway',
			label: 'Rip Leeway',
			labelTooltip: 'Rip leeway when determining roar clips',
			showWhen: ShouldShowAdvParamST,
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecFeralDruid>({
			fieldName: 'useRake',
			label: 'Use Rake',
			labelTooltip: 'Use rake during rotation',
			showWhen: ShouldShowAdvParamST,
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecFeralDruid>({
			fieldName: 'useBite',
			label: 'Bite during rotation',
			labelTooltip: 'Use bite during rotation rather than just at end',
			showWhen: ShouldShowAdvParamST,
		}),
		InputHelpers.makeRotationNumberInput<Spec.SpecFeralDruid>({
			fieldName: 'biteTime',
			label: 'Bite Time',
			labelTooltip: 'Min seconds on Rip/Roar to bite',
			showWhen: (player: Player<Spec.SpecFeralDruid>) =>
				ShouldShowAdvParamST(player) && player.getSimpleRotation().useBite == true && player.getSimpleRotation().biteModeType == BiteModeType.Emperical,
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecFeralDruid>({
			fieldName: 'flowerWeave',
			label: 'Flower Weave',
			labelTooltip: 'Fish for clearcasting during rotation with gotw',
			showWhen: ShouldShowAdvParamAoe,
		}),
		InputHelpers.makeRotationNumberInput<Spec.SpecFeralDruid>({
			extraCssClasses: ['used-in-apl'],
			fieldName: 'raidTargets',
			label: 'GotW Raid Targets',
			labelTooltip: 'Raid size to assume for clearcast proc chance (can include pets as well, so 25 man raid potentically can be ~30)',
			showWhen: (player: Player<Spec.SpecFeralDruid>) => player.aplRotation.type != APLRotation_Type.TypeSimple || (ShouldShowAdvParamAoe(player) && player.getSimpleRotation().flowerWeave == true),
		}),
		// Can be uncommented if/when analytical bite mode is added
		//InputHelpers.makeRotationEnumInput<Spec.SpecFeralDruid, BiteModeType>({
		//	fieldName: 'biteModeType',
		//	label: 'Bite Mode',
		//	labelTooltip: 'Underlying "Bite logic" to use',
		//	values: [
		//		{ name: 'Emperical', value: BiteModeType.Emperical },
		//	],
		//	showWhen: (player: Player<Spec.SpecFeralDruid>) => player.getSimpleRotation().useBite == true
		//}),
		InputHelpers.makeRotationNumberInput<Spec.SpecFeralDruid>({
			fieldName: 'hotUptime',
			label: 'Revitalize Hot Uptime',
			labelTooltip: 'Hot uptime percentage to assume when theorizing energy gains',
			percent: true,
			showWhen: (player: Player<Spec.SpecFeralDruid>) => player.getSimpleRotation().useBite == true && player.getSimpleRotation().biteModeType == BiteModeType.Analytical,
		}),
	],
};
