import { Spec } from '../core/proto/common.js';
import { Player } from '../core/player.js';
import { EventID } from '../core/typed_event.js';
import { IndividualSimUI } from '../core/individual_sim_ui.js';
import { ActionId } from '../core/proto_utils/action_id.js';
import { CustomRotationPickerConfig } from '../core/components/custom_rotation_picker.js';
import { CustomRotation } from '../core/proto/common.js';

import {
	PaladinAura as PaladinAura,
	RetributionPaladin_Rotation as RetributionPaladinRotation,
	RetributionPaladin_Options as RetributionPaladinOptions,
	RetributionPaladin_Rotation_SpellOption as SpellOption,
	RetributionPaladin_Rotation_RotationType as RotationType,
	PaladinJudgement as PaladinJudgement,
	PaladinSeal,
} from '../core/proto/paladin.js';

import * as InputHelpers from '../core/components/input_helpers.js';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.
export const RetributionPaladinRotationExoSlackConfig = InputHelpers.makeRotationNumberInput<Spec.SpecRetributionPaladin>({
	fieldName: "exoSlack",
	label: "Exo Slack (MS)",
	labelTooltip: "Amount of extra time in MS to give main abilities to come off cooldown before using Exorcism on single target",
	positive: true,
})

export const RetributionPaladinRotationConsSlackConfig = InputHelpers.makeRotationNumberInput<Spec.SpecRetributionPaladin>({
	fieldName: "consSlack",
	label: "Cons Slack (MS)",
	labelTooltip: "Amount of extra time in MS to give main abilities to come off cooldown before using Consecration on single target",
	positive: true,
})

export const RetributionPaladinRotationDivinePleaPercentageConfig = InputHelpers.makeRotationNumberInput<Spec.SpecRetributionPaladin>({
	fieldName: "divinePleaPercentage",
	label: "Divine Plea Mana Threshold %",
	labelTooltip: "% of max mana left before beginning to use Divine Plea",
	percent: true,
	positive: true,
	showWhen: (player: Player<Spec.SpecRetributionPaladin>) => player.getRotation().type == RotationType.Standart,
})

export const RetributionPaladinRotationHolyWrathConfig = InputHelpers.makeRotationNumberInput<Spec.SpecRetributionPaladin>({
	fieldName: "holyWrathThreshold",
	label: "Holy Wrath Threshold",
	labelTooltip: "Minimum number of Demon and Undead units before Holy Wrath is considered viable to add to an AOE rotation.",
	positive: true,
	showWhen: (player: Player<Spec.SpecRetributionPaladin>) => player.getRotation().type == RotationType.Standart,
})

export const AuraSelection = InputHelpers.makeSpecOptionsEnumInput<Spec.SpecRetributionPaladin, PaladinAura>({
	fieldName: 'aura',
	label: 'Aura',
	values: [
		{ name: 'None', value: PaladinAura.NoPaladinAura },
		{ name: 'Retribution Aura', value: PaladinAura.RetributionAura },
	],
});

export const StartingSealSelection = InputHelpers.makeSpecOptionsEnumInput<Spec.SpecRetributionPaladin, PaladinSeal>({
	fieldName: 'seal',
	label: 'Seal',
	labelTooltip: 'The seal active before encounter',
	values: [
		{ name: 'Vengeance', value: PaladinSeal.Vengeance },
		{ name: 'Command', value: PaladinSeal.Command },
		{ name: 'Righteousness', value: PaladinSeal.Righteousness },
	],
});

export const DivinePleaSelection = InputHelpers.makeSpecOptionsBooleanInput<Spec.SpecRetributionPaladin>({
	fieldName: 'useDivinePlea',
	label: 'Divine Plea',
	labelTooltip: 'Whether or not to maintain Divine Plea',
});

export const JudgementSelection = InputHelpers.makeSpecOptionsEnumInput<Spec.SpecRetributionPaladin, PaladinJudgement>({
	fieldName: 'judgement',
	label: 'Judgement',
	labelTooltip: 'Judgement debuff you will use on the target during the encounter.',
	values: [
		{ name: 'Wisdom', value: PaladinJudgement.JudgementOfWisdom },
		{ name: 'Light', value: PaladinJudgement.JudgementOfLight },
	],
});

export const DamageTakenPerSecond = InputHelpers.makeSpecOptionsNumberInput<Spec.SpecRetributionPaladin>({
	fieldName: 'damageTakenPerSecond',
	label: 'Damage Taken Per Second',
	labelTooltip: "Damage taken per second across the encounter. Used to model mana regeneration from Spiritual Attunement. This value should NOT include damage taken from Seal of Blood / Judgement of Blood. Leave at 0 if you do not take damage during the encounter.",
});

export const RetributionPaladinRotationPriorityConfig = InputHelpers.makeCustomRotationInput<Spec.SpecRetributionPaladin, SpellOption>({
	fieldName: 'customRotation',
	numColumns: 2,
	values: [
		{ actionId: ActionId.fromSpellId(53408), value: SpellOption.JudgementOfWisdom },
		{ actionId: ActionId.fromSpellId(53385), value: SpellOption.DivineStorm },
		{ actionId: ActionId.fromSpellId(48806), value: SpellOption.HammerOfWrath },
		{ actionId: ActionId.fromSpellId(48819), value: SpellOption.Consecration },
		{ actionId: ActionId.fromSpellId(48817), value: SpellOption.HolyWrath },
		{ actionId: ActionId.fromSpellId(35395), value: SpellOption.CrusaderStrike },
		{ actionId: ActionId.fromSpellId(48801), value: SpellOption.Exorcism },
	],
	showWhen: (player: Player<Spec.SpecRetributionPaladin>) => player.getRotation().type == RotationType.Custom,
});

export const RotationSelector = InputHelpers.makeRotationEnumInput<Spec.SpecRetributionPaladin, RotationType>({
	fieldName: 'type',
	label: 'Type',
	values: [
		{ name: 'Standart', value: RotationType.Standart },
		{ name: 'Custom', value: RotationType.Custom },
	],
});