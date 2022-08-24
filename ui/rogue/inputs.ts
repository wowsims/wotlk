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
		{ color: 'grey', value: Poison.NoPoison },
		{ actionId: ActionId.fromItemId(43233), value: Poison.DeadlyPoison },
		{ actionId: ActionId.fromItemId(43231), value: Poison.InstantPoison },
		{ actionId: ActionId.fromItemId(43235), value: Poison.WoundPoison },
	],
});

export const OffHandImbue = InputHelpers.makeSpecOptionsEnumIconInput<Spec.SpecRogue, Poison>({
	fieldName: 'ohImbue',
	numColumns: 1,
	values: [
		{ color: 'grey', value: Poison.NoPoison },
		{ actionId: ActionId.fromItemId(43233), value: Poison.DeadlyPoison },
		{ actionId: ActionId.fromItemId(43231), value: Poison.InstantPoison },
		{ actionId: ActionId.fromItemId(43235), value: Poison.WoundPoison },
	],
});

export const RogueRotationConfig = {
	inputs: [
		InputHelpers.makeRotationEnumInput<Spec.SpecRogue, Frequency>({
			fieldName: 'exposeArmorFrequency',
			label: 'Expose Armor',
			labelTooltip: 'Frequency of Expose Armor casts.',
			values:[
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
			values:[
				{ name: 'Never', value: Frequency.Never },
				{ name: 'Maintain', value: Frequency.Maintain },
			],
		}),
		InputHelpers.makeRotationEnumInput<Spec.SpecRogue, AssassinationPriority>({
			fieldName: 'assassinationFinisherPriority',
			label: 'Finisher Priority (Assassination)',
			labelTooltip: 'Priority of Assassination finisher usage.',
			values:[
				{ name: 'Envenom > Rupture', value: AssassinationPriority.EnvenomRupture },
				{ name: 'Rupture > Envenom', value: AssassinationPriority.RuptureEnvenom },
			],
			showWhen: (player: Player<Spec.SpecRogue>) => player.getTalents().mutilate
		}),
		InputHelpers.makeRotationEnumInput<Spec.SpecRogue, AssassinationPriority>({
			fieldName: 'combatFinisherPriority',
			label: 'Finisher Priority (Combat)',
			labelTooltip: 'Priority of Combat finisher usage.',
			values:[
				{ name: 'Rupture > Eviscerate', value: CombatPriority.RuptureEviscerate },
				{ name: 'Eviscerate > Rupture', value: CombatPriority.EviscerateRupture },
			],
			showWhen: (player: Player<Spec.SpecRogue>) => !player.getTalents().mutilate
		}),
		InputHelpers.makeRotationNumberInput<Spec.SpecRogue>({
			fieldName: 'minimumComboPointsPrimaryFinisher',
			label: 'Minimum CP (Finisher)',
			labelTooltip: 'Primary finisher will not be cast with less than this many combo points.',
		}),
		InputHelpers.makeRotationNumberInput<Spec.SpecRogue>({
			fieldName: 'minimumComboPointsSecondaryFinisher',
			label: 'Minimum CP (Filler)',
			labelTooltip: 'Secondary finisher/filler will not be cast with less than this many combo points.\nSet the value to > 5 to prevent fillers.',
		}),
		InputHelpers.makeRotationEnumInput<Spec.SpecRogue, Frequency>({
			fieldName: 'multiTargetSliceFrequency',
			label: 'Multi-Target S&D',
			labelTooltip: 'Frequency of Slice and Dice cast in multi-target scnearios.',
			values:[
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
	],
};
