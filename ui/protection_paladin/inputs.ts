import { Spec } from '../core/proto/common.js';

import {
	PaladinAura as PaladinAura,
	PaladinSeal,
	PaladinJudgement as PaladinJudgement,
} from '../core/proto/paladin.js';

import * as InputHelpers from '../core/components/input_helpers.js';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const AuraSelection = InputHelpers.makeSpecOptionsEnumInput<Spec.SpecProtectionPaladin, PaladinAura>({
	fieldName: 'aura',
	label: 'Aura',
	values: [
		{ name: 'None', value: PaladinAura.NoPaladinAura },
		{ name: 'Devotion Aura', value: PaladinAura.DevotionAura },
		{ name: 'Retribution Aura', value: PaladinAura.RetributionAura },
	],
});

export const StartingSealSelection = InputHelpers.makeSpecOptionsEnumInput<Spec.SpecProtectionPaladin, PaladinSeal>({
	fieldName: 'seal',
	label: 'Seal',
	labelTooltip: 'The seal active before encounter',
	values: [
		{ name: 'Vengeance', value: PaladinSeal.Vengeance },
		{ name: 'Command', value: PaladinSeal.Command },
	],
});

export const JudgementSelection = InputHelpers.makeSpecOptionsEnumInput<Spec.SpecProtectionPaladin, PaladinJudgement>({
	fieldName: 'judgement',
	label: 'Judgement',
	labelTooltip: 'Judgement debuff you will use on the target during the encounter.',
	values: [
		{ name: 'Wisdom', value: PaladinJudgement.JudgementOfWisdom },
		{ name: 'Light', value: PaladinJudgement.JudgementOfLight },
	],
});

export const UseAvengingWrath = InputHelpers.makeSpecOptionsBooleanInput<Spec.SpecProtectionPaladin>({
	fieldName: 'useAvengingWrath',
	label: 'Use Avenging Wrath',
});
