import { Spec } from '../core/proto/common.js';
import { ActionId } from '../core/proto_utils/action_id.js';
import { Player } from '../core/player.js';

import * as InputHelpers from '../core/components/input_helpers.js';

import {
	Rogue_Rotation_AssassinationPriority as AssassinationPriority,
	Rogue_Rotation_CombatPriority as CombatPriority,
	Rogue_Rotation_CombatBuilder as CombatBuilder,
	Rogue_Rotation_SubtletyPriority as SubtletyPriority,
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

export const HonorOfThievesCritRate = InputHelpers.makeSpecOptionsNumberInput<Spec.SpecRogue>({
	fieldName: 'honorOfThievesCritRate',
	label: 'Honor of Thieves Crit Rate',
	labelTooltip: 'Number of crits other group members generate within 100 seconds',
	showWhen: (player: Player<Spec.SpecRogue>) => player.getTalents().honorAmongThieves > 0
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
		InputHelpers.makeRotationEnumInput<Spec.SpecRogue, CombatBuilder>({
			fieldName: 'combatBuilder',
			label: "Builder",
			labelTooltip: 'Use Sinister Strike or Backstab as builder.',
			values: [
				{ name: "Sinister Strike", value: CombatBuilder.SinisterStrike },
				{ name: "Backstab", value: CombatBuilder.Backstab },
			],
			showWhen: (player: Player<Spec.SpecRogue>) => player.getTalents().combatPotency > 0
		}),
		InputHelpers.makeRotationEnumInput<Spec.SpecRogue, CombatPriority>({
			fieldName: 'combatFinisherPriority',
			label: 'Finisher Priority',
			labelTooltip: 'The finisher that will be cast with highest priority.',
			values: [
				{ name: 'Rupture', value: CombatPriority.RuptureEviscerate },
				{ name: 'Eviscerate', value: CombatPriority.EviscerateRupture },
			],
			showWhen: (player: Player<Spec.SpecRogue>) => player.getTalents().combatPotency > 0
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
		InputHelpers.makeRotationEnumInput<Spec.SpecRogue, SubtletyPriority>({
			fieldName: 'subtletyFinisherPriority',
			label: "Finisher Priority",
			labelTooltip: 'The finisher that will be cast with highest priority.',
			values: [
				{ name: "Eviscerate", value: SubtletyPriority.SubtletyEviscerate },
				{ name: "Envenom", value: SubtletyPriority.SubtletyEnvenom },
			],
			showWhen: (player: Player<Spec.SpecRogue>) => player.getTalents().honorAmongThieves > 0
		}),
		InputHelpers.makeRotationEnumInput<Spec.SpecRogue, Frequency>({
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
			fieldName: 'hemoWithDagger',
			label: 'Hemorrhage with Dagger',
			labelTooltip: 'Use Hemorrhage with Dagger in mainhand',
			showWhen: (player: Player<Spec.SpecRogue>) => player.getTalents().hemorrhage
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecRogue>({
			fieldName: 'openWithGarrote',
			label: 'Open with Garrote',
			labelTooltip: 'Open the encounter by casting Garrote.',
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecRogue>({
			fieldName: 'openWithShadowstep',
			label: 'Open with Shadowstep',
			labelTooltip: 'Open the encounter by casting Shadowstep.',
			showWhen: (player: Player<Spec.SpecRogue>) => player.getTalents().shadowstep
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
		InputHelpers.makeRotationBooleanInput<Spec.SpecRogue>({
			fieldName: 'ruptureForBleed',
			label: 'Rupture for Bleed',
			labelTooltip: 'Cast Rupture as needed to apply a bleed effect for Hunger for Blood',
			showWhen: (player: Player<Spec.SpecRogue>) => player.getTalents().hungerForBlood
		}),
	],
};
