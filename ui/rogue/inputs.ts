import { Spec } from '../core/proto/common.js';
import { ActionId } from '../core/proto_utils/action_id.js';
import { Player } from '../core/player.js';

import * as InputHelpers from '../core/components/input_helpers.js';

import {
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

export const ApplyPoisonsManually = InputHelpers.makeSpecOptionsBooleanInput<Spec.SpecRogue>({
	fieldName: 'applyPoisonsManually',
	label: 'Configure poisons manually',
	labelTooltip: 'Prevent automatic poison configuration that is based on equipped weapons.',
});

export const RogueRotationConfig = {
	inputs: [
		InputHelpers.makeRotationEnumInput<Spec.SpecRogue>({
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
		InputHelpers.makeRotationEnumInput<Spec.SpecRogue>({
			fieldName: 'tricksOfTheTradeFrequency',
			label: 'Tricks of the Trade',
			labelTooltip: 'Frequency of Tricks of the Trade usage.',
			values: [
				{ name: 'Never', value: Frequency.Never },
				{ name: 'Maintain', value: Frequency.Maintain },
			],
		}),
		InputHelpers.makeRotationEnumInput<Spec.SpecRogue>({
			fieldName: 'multiTargetSliceFrequency',
			label: 'Multi-Target S&D',
			labelTooltip: 'Frequency of Slice and Dice cast in multi-target scenarios.',
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
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecRogue>({
			fieldName: 'openWithPremeditation',
			label: 'Open with Premeditation',
			labelTooltip: 'Open the encounter by casting Premeditation.',
			showWhen: (player: Player<Spec.SpecRogue>) => player.getTalents().premeditation
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecRogue>({
			fieldName: 'useFeint',
			label: 'Use Feint',
			labelTooltip: 'Cast Feint on cooldown. Mainly useful when using the associate glyph.'
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecRogue>({
			fieldName: "useGhostlyStrike",
			label: 'Use Ghostly Strike',
			labelTooltip: 'Use Ghostly Strike on cooldown. Mainly useful when using the associate glyph.',
			showWhen: (player: Player<Spec.SpecRogue>) => player.getTalents().ghostlyStrike
		}),
	],
};
