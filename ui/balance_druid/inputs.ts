import { Player } from '../core/player.js';
import { Spec, UnitReference, UnitReference_Type as UnitType } from '../core/proto/common.js';
import { ActionId } from '../core/proto_utils/action_id.js';
import { EventID } from '../core/typed_event.js';

import * as InputHelpers from '../core/components/input_helpers.js';

import {
	BalanceDruid_Rotation_EclipsePrio as EclipsePrio,
	BalanceDruid_Rotation_IsUsage as IsUsage,
	BalanceDruid_Rotation_MfExtension as MfExtension,
	BalanceDruid_Rotation_MfUsage as MfUsage,
	BalanceDruid_Rotation_Type as RotationType,
	BalanceDruid_Rotation_WrathUsage as WrathUsage
} from '../core/proto/druid.js';


// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const SelfInnervate = InputHelpers.makeSpecOptionsBooleanIconInput<Spec.SpecBalanceDruid>({
	fieldName: 'innervateTarget',
	id: ActionId.fromSpellId(29166),
	extraCssClasses: [
		'within-raid-sim-hide',
	],
	getValue: (player: Player<Spec.SpecBalanceDruid>) => player.getSpecOptions().innervateTarget?.type == UnitType.Player,
	setValue: (eventID: EventID, player: Player<Spec.SpecBalanceDruid>, newValue: boolean) => {
		const newOptions = player.getSpecOptions();
		newOptions.innervateTarget = UnitReference.create({
			type: newValue ? UnitType.Player : UnitType.Unknown,
			index: 0,
		});
		player.setSpecOptions(eventID, newOptions);
	},
});

export const OkfUptime = InputHelpers.makeSpecOptionsNumberInput<Spec.SpecBalanceDruid>({
	fieldName: 'okfUptime',
	label: 'Owlkin Frenzy Uptime (%)',
	labelTooltip: 'Percentage of fight uptime for Owlkin Frenzy',
	percent: true,
});

export const BalanceDruidRotationConfig = {
	inputs: [
		InputHelpers.makeRotationEnumInput<Spec.SpecBalanceDruid>({
			fieldName: 'type',
			label: 'Type',
			labelTooltip: 'Set to \'Manual\', to manage eclipses, spells, CDs and DoTs usage.',
			values: [
				{
					name: 'Default', value: RotationType.Default,
					tooltip: 'The default rotation.',
				},
				{
					name: 'Manual', value: RotationType.Manual,
					tooltip: 'Allows custom selection of which spells to use, dot management and CD usage.',
				},
			],
		}),
		InputHelpers.makeRotationEnumInput<Spec.SpecBalanceDruid>({
			fieldName: 'eclipsePrio',
			label: 'Eclipse priority',
			labelTooltip: 'Defines which eclipse will get prioritized in the rotation.',
			values: [
				{ name: 'Lunar', value: EclipsePrio.Lunar },
				{ name: 'Solar', value: EclipsePrio.Solar },
			],
			showWhen: (player: Player<Spec.SpecBalanceDruid>) => player.getRotation().type == RotationType.Manual,
		}),
		InputHelpers.makeRotationEnumInput<Spec.SpecBalanceDruid>({
			fieldName: 'mfUsage',
			label: 'Moonfire Usage',
			labelTooltip: 'Defines how Moonfire will be used in the rotation.',
			values: [
				{ name: 'Unused', value: MfUsage.NoMf },
				{ name: 'Before lunar', value: MfUsage.BeforeLunar },
				{ name: 'Maximize', value: MfUsage.MaximizeMf },
				{ name: 'Multidot', value: MfUsage.MultidotMf },
			],
			showWhen: (player: Player<Spec.SpecBalanceDruid>) => player.getRotation().type == RotationType.Manual,
		}),
		InputHelpers.makeRotationEnumInput<Spec.SpecBalanceDruid>({
			fieldName: 'mfExtension',
			label: 'Moonfire Extension',
			labelTooltip: 'When should the rotation try to extend Moonfire on the main target.',
			values: [
				{ name: 'Extend always', value: MfExtension.ExtendAlways },
				{ name: 'Extend outside solar', value: MfExtension.ExtendOutsideSolar },
				{ name: 'Do not extend', value: MfExtension.DontExtend },
			],
			showWhen: (player: Player<Spec.SpecBalanceDruid>) => player.getRotation().type == RotationType.Manual,
		}),
		InputHelpers.makeRotationEnumInput<Spec.SpecBalanceDruid>({
			fieldName: 'isUsage',
			label: 'Insect Swarm Usage',
			labelTooltip: 'Defines how Insect Swarm will be used in the rotation.',
			values: [
				{ name: 'Unused', value: IsUsage.NoIs },
				{ name: 'Before solar', value: IsUsage.BeforeSolar },
				{ name: 'Optimize', value: IsUsage.OptimizeIs },
				{ name: 'Multidot', value: IsUsage.MultidotIs },
			],
			showWhen: (player: Player<Spec.SpecBalanceDruid>) => player.getRotation().type == RotationType.Manual,
		}),
		InputHelpers.makeRotationEnumInput<Spec.SpecBalanceDruid>({
			fieldName: 'wrathUsage',
			label: 'Wrath usage',
			labelTooltip: 'Defines how Wrath will be used in the rotation.',
			values: [
				{ name: 'Unused', value: WrathUsage.NoWrath },
				{ name: 'Fishing for Lunar', value: WrathUsage.FishingForLunar },
				{ name: 'Regular', value: WrathUsage.RegularWrath },
			],
			showWhen: (player: Player<Spec.SpecBalanceDruid>) => player.getRotation().type == RotationType.Manual,
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecBalanceDruid>({
			fieldName: 'useStarfire',
			label: 'Use Starfire',
			labelTooltip: 'Should the rotation use Starfire.',
			showWhen: (player: Player<Spec.SpecBalanceDruid>) => player.getRotation().type == RotationType.Manual,
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecBalanceDruid>({
			fieldName: 'useSmartCooldowns',
			label: 'Smart Cooldowns usage',
			labelTooltip: 'The rotation will use cooldowns during eclipses, avoiding Haste CDs in solar and Crit CDs in lunar',
			showWhen: (player: Player<Spec.SpecBalanceDruid>) => player.getRotation().type == RotationType.Manual,
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecBalanceDruid>({
			fieldName: 'snapshotMf',
			label: 'Snapshot Moonfire',
			labelTooltip: 'The rotation will try to snapshot Moonfire with SP procs',
			showWhen: (player: Player<Spec.SpecBalanceDruid>) => player.getRotation().type == RotationType.Manual,
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecBalanceDruid>({
			fieldName: 'eclipseShuffling',
			label: 'Eclipse Shuffling',
			labelTooltip: 'Should the rotation alternate Starfire and Wrath when both eclipses are available.',
			showWhen: (player: Player<Spec.SpecBalanceDruid>) => player.getRotation().type == RotationType.Manual,
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecBalanceDruid>({
			fieldName: 'useTyphoon',
			label: 'Use Typhoon',
			labelTooltip: 'Should the rotation use Typhoon.',
			showWhen: (player: Player<Spec.SpecBalanceDruid>) => player.getRotation().type == RotationType.Manual,
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecBalanceDruid>({
			fieldName: 'useHurricane',
			label: 'Use Hurricane',
			labelTooltip: 'Should the rotation use Hurricane.',
			showWhen: (player: Player<Spec.SpecBalanceDruid>) => player.getRotation().type == RotationType.Manual,
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecBalanceDruid>({
			fieldName: 'useBattleRes',
			label: 'Use Battle Res',
			labelTooltip: 'Cast Battle Res on an ally sometime during the encounter.',
			showWhen: (player: Player<Spec.SpecBalanceDruid>) => player.getRotation().type == RotationType.Manual,
		}),
		InputHelpers.makeRotationNumberInput<Spec.SpecBalanceDruid>({
			fieldName: 'playerLatency',
			label: 'Player latency',
			labelTooltip: 'Time before the player reacts to an eclipse proc, in milliseconds.',
			showWhen: (player: Player<Spec.SpecBalanceDruid>) => player.getRotation().type == RotationType.Manual,
		}),
	],
};
