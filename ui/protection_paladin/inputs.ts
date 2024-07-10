import * as InputHelpers from '../core/components/input_helpers.js';
import { Spec } from '../core/proto/common.js';
import {
	PaladinAura as PaladinAura,
	PaladinJudgement as PaladinJudgement,
	PaladinSeal,
} from '../core/proto/paladin.js';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const AuraSelection = InputHelpers.makeSpecOptionsEnumInput<Spec.SpecProtectionPaladin, PaladinAura>({
	fieldName: 'aura',
	label: '光环',
	values: [
		{ name: 'None', value: PaladinAura.NoPaladinAura },
		{ name: '虔诚光环', value: PaladinAura.DevotionAura },
		{ name: '惩戒光环', value: PaladinAura.RetributionAura },
	],
});

export const StartingSealSelection = InputHelpers.makeSpecOptionsEnumInput<Spec.SpecProtectionPaladin, PaladinSeal>({
	fieldName: 'seal',
	label: '圣印',
	labelTooltip: '战斗前激活的圣印',
	values: [
		{ name: '复仇', value: PaladinSeal.Vengeance },
		{ name: '命令', value: PaladinSeal.Command },
	],
});

export const JudgementSelection = InputHelpers.makeSpecOptionsEnumInput<Spec.SpecProtectionPaladin, PaladinJudgement>({
	fieldName: 'judgement',
	label: '审判',
	labelTooltip: '战斗中你将对目标使用的审判减益。',
	values: [
		{ name: '智慧', value: PaladinJudgement.JudgementOfWisdom },
		{ name: '光明', value: PaladinJudgement.JudgementOfLight },
	],
});

export const UseAvengingWrath = InputHelpers.makeSpecOptionsBooleanInput<Spec.SpecProtectionPaladin>({
	fieldName: 'useAvengingWrath',
	label: '使用复仇之怒',
});
