import { Spec } from '../core/proto/common.js';
import { Player } from '../core/player.js';
import { EventID } from '../core/typed_event.js';
import { IndividualSimUI } from '../core/individual_sim_ui.js';
import { ActionId } from '../core/proto_utils/action_id.js';

import {
	PaladinAura,
	PaladinSeal,
	PaladinJudgement,
	HolyPaladin_Rotation as HolyPaladinRotation,
	HolyPaladin_Options as HolyPaladinOptions,
} from '../core/proto/paladin.js';

import * as InputHelpers from '../core/components/input_helpers.js';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.
export const HolyPaladinRotationConfig = {
	inputs: [
	],
}

export const AuraSelection = InputHelpers.makeSpecOptionsEnumInput<Spec.SpecHolyPaladin, PaladinAura>({
	fieldName: 'aura',
	label: 'Aura',
	values: [
		{ name: 'None', value: PaladinAura.NoPaladinAura },
		{ name: 'Devotion Aura', value: PaladinAura.DevotionAura },
		{ name: 'Retribution Aura', value: PaladinAura.RetributionAura },
	],
});

export const JudgementSelection = InputHelpers.makeSpecOptionsEnumInput<Spec.SpecHolyPaladin, PaladinJudgement>({
	fieldName: 'judgement',
	label: 'Judgement',
	labelTooltip: 'Judgement debuff you will use on the target during the encounter.',
	values: [
		{ name: 'None', value: PaladinJudgement.NoJudgement },
		{ name: 'Wisdom', value: PaladinJudgement.JudgementOfWisdom },
		{ name: 'Light', value: PaladinJudgement.JudgementOfLight },
	],
});
