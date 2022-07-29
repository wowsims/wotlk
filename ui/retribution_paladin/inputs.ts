import { Spec } from '/wotlk/core/proto/common.js';
import { Player } from '/wotlk/core/player.js';
import { EventID } from '/wotlk/core/typed_event.js';
import { IndividualSimUI } from '/wotlk/core/individual_sim_ui.js';

import {
	PaladinAura as PaladinAura,
	RetributionPaladin_Rotation as RetributionPaladinRotation,
	RetributionPaladin_Options as RetributionPaladinOptions,
	PaladinJudgement as PaladinJudgement,
	PaladinSeal,
} from '/wotlk/core/proto/paladin.js';

import * as InputHelpers from '/wotlk/core/components/input_helpers.js';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.
export const RetributionPaladinRotationExoSlackConfig = InputHelpers.makeRotationNumberInput<Spec.SpecRetributionPaladin>({
	fieldName: "exoSlack",
	label: "Exo Slack (MS)",
	labelTooltip: "Amount of extra time in MS to give main abilities to come off cooldown before using Exorcism on single target",
})

export const RetributionPaladinRotationConsSlackConfig = InputHelpers.makeRotationNumberInput<Spec.SpecRetributionPaladin>({
	fieldName: "consSlack",
	label: "Cons Slack (MS)",
	labelTooltip: "Amount of extra time in MS to give main abilities to come off cooldown before using Consecration on single target",
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
