import { Spec } from '../core/proto/common.js';
import { ActionId } from '../core/proto_utils/action_id.js';
import { Player } from '../core/player.js';

import * as InputHelpers from '../core/components/input_helpers.js';

import {
	Rogue_Rotation_AssassinationPriority as AssassinationPriority,
	Rogue_Rotation_CombatPriority as CombatPriority,
	Rogue_Rotation_Frequency as Frequency,
	Rogue_Options_PoisonImbue as Poison,
} from '../core/proto/rogue.js';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const MainHandImbue = InputHelpers.makeSpecOptionsEnumIconInput<Spec.SpecRogue, Poison>({
	fieldName: 'mhImbue',
	numColumns: 1,
	values: [
		{ value: Poison.NoPoison, tooltip: 'No Main Hand Poison' },
		{ actionId: ActionId.fromItemId(43233), value: Poison.DeadlyPoison },
		{ actionId: ActionId.fromItemId(43231), value: Poison.InstantPoison },
		{ actionId: ActionId.fromItemId(43235), value: Poison.WoundPoison },
	],
});

export const OffHandImbue = InputHelpers.makeSpecOptionsEnumIconInput<Spec.SpecRogue, Poison>({
	fieldName: 'ohImbue',
	numColumns: 1,
	values: [
		{ value: Poison.NoPoison, tooltip: 'No Off Hand Poison' },
		{ actionId: ActionId.fromItemId(43233), value: Poison.DeadlyPoison },
		{ actionId: ActionId.fromItemId(43231), value: Poison.InstantPoison },
		{ actionId: ActionId.fromItemId(43235), value: Poison.WoundPoison },
	],
});

export const StartingOverkillDuration = InputHelpers.makeSpecOptionsNumberInput<Spec.SpecRogue>({
	fieldName: 'startingOverkillDuration',
	label: 'Starting Overkill duration',
	labelTooltip: 'Initial Overkill buff duration at the start of each iteration.',
});

export const ApplyPoisonsManually = InputHelpers.makeSpecOptionsBooleanInput<Spec.SpecRogue>({
	fieldName: 'applyPoisonsManually',
	label: 'Configure poisons manually',
	labelTooltip: 'Prevent automatic poison configuration that is based on equipped weapons.',
});

export const RogueRotationConfig = {
	inputs: [
		InputHelpers.makeRotationEnumInput<Spec.SpecRogue, Frequency>({
			fieldName: 'exposeArmorFrequency',
			label: 'Expose Armor',
			labelTooltip: 'Frequency of Expose Armor casts.',
			values: [
				{ name: 'Never', value: Frequency.Never },
				{ name: 'Cast Once', value: Frequency.Once },
				{ name: 'Maintain', value: Frequency.Maintain },
			],
		}),
		InputHelpers.makeRotationNumberInput<Spec.SpecRogue>({
			fieldName: 'minimumComboPointsExposeArmor',
			label: 'Minimum CP (Expose Armor)',
			labelTooltip: 'Minimum number of combo points for Expose Armor when only cast once.',
			showWhen: (player: Player<Spec.SpecRogue>) => player.getRotation().exposeArmorFrequency == Frequency.Once,
		}),
		InputHelpers.makeRotationEnumInput<Spec.SpecRogue, Frequency>({
			fieldName: 'tricksOfTheTradeFrequency',
			label: 'Tricks of the Trade',
			labelTooltip: 'Frequency of Tricks of the Trade usage.',
			values: [
				{ name: 'Never', value: Frequency.Never },
				{ name: 'Maintain', value: Frequency.Maintain },
			],
		}),
		InputHelpers.makeRotationEnumInput<Spec.SpecRogue, AssassinationPriority>({
			fieldName: 'assassinationFinisherPriority',
			label: 'Finisher Priority',
			labelTooltip: 'The finisher that will be cast with highest priority.',
			values: [
				{ name: 'Rupture', value: AssassinationPriority.RuptureEnvenom },
				{ name: 'Envenom', value: AssassinationPriority.EnvenomRupture },
			],
			showWhen: (player: Player<Spec.SpecRogue>) => player.getTalents().mutilate
		}),
		InputHelpers.makeRotationNumberInput<Spec.SpecRogue>({
			fieldName: 'envenomEnergyThreshold',
			label: 'Energy Threshold (Envenom)',
			labelTooltip: 'Amount of total energy to pool before casting Envenom.',
			showWhen: (player: Player<Spec.SpecRogue>) => player.getTalents().mutilate
		}),
		InputHelpers.makeRotationEnumInput<Spec.SpecRogue, Frequency>({
			fieldName: 'multiTargetSliceFrequency',
			label: 'Multi-Target S&D',
			labelTooltip: 'Frequency of Slice and Dice cast in multi-target scnearios.',
			values: [
				{ name: 'Never', value: Frequency.Never },
				{ name: 'Once', value: Frequency.Once },
				{ name: 'Maintain', value: Frequency.Maintain },
			],
			showWhen: (player: Player<Spec.SpecRogue>) => player.getRotation().multiTargetSliceFrequency != Frequency.FrequencyUnknown
		}),
		InputHelpers.makeRotationNumberInput<Spec.SpecRogue>({
			fieldName: 'minimumComboPointsMultiTargetSlice',
			label: 'Minimum CP (Slice)',
			labelTooltip: 'Minimum number of combo points spent if Slice and Dice has frequency: Once',
			showWhen: (player: Player<Spec.SpecRogue>) => player.getRotation().multiTargetSliceFrequency == Frequency.Once
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecRogue>({
			fieldName: 'openWithGarrote',
			label: 'Open with Garrote',
			labelTooltip: 'Open the encounter by casting Garrote.',
			showWhen: (player: Player<Spec.SpecRogue>) => player.getTalents().mutilate
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecRogue>({
			fieldName: 'useFeint',
			label: 'Use Feint',
			labelTooltip: 'Cast Feint on cooldown. Mainly useful when using the associate glyph.'
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecRogue>({
			fieldName: 'allowCpUndercap',
			label: 'Undercap CP',
			labelTooltip: 'Cast Envenom at 3 cp if the Envenom buff is missing.',
			showWhen: (player: Player<Spec.SpecRogue>) => player.getTalents().mutilate
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecRogue>({
			fieldName: 'allowCpOvercap',
			label: 'Overcap CP',
			labelTooltip: 'Cast Mutilate at 4 cp if the Envenom buff will last long enough.',
			showWhen: (player: Player<Spec.SpecRogue>) => player.getTalents().mutilate
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecRogue>({
			fieldName: 'ruptureForBleed',
			label: 'Rupture for Bleed',
			labelTooltip: 'Cast Rupture as needed to apply a bleed effect for Hunger for Blood',
			showWhen: (player: Player<Spec.SpecRogue>) => player.getTalents().hungerForBlood
		}),
	],
};
