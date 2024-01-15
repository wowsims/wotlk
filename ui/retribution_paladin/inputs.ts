import { Spec } from '../core/proto/common.js';
import { Player } from '../core/player.js';
import { ActionId } from '../core/proto_utils/action_id.js';

import {
	PaladinAura as PaladinAura,
	PaladinJudgement as PaladinJudgement,
	PaladinSeal,
} from '../core/proto/paladin.js';

import * as InputHelpers from '../core/components/input_helpers.js';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const AuraSelection = InputHelpers.makeSpecOptionsEnumIconInput<Spec.SpecRetributionPaladin, PaladinAura>({
	fieldName: 'aura',
	values: [
		{ value: PaladinAura.NoPaladinAura, tooltip: 'No Aura' },
		{ actionId: ActionId.fromSpellId(54043), value: PaladinAura.RetributionAura },
	],
});

export const StartingSealSelection = InputHelpers.makeSpecOptionsEnumIconInput<Spec.SpecRetributionPaladin, PaladinSeal>({
	fieldName: 'seal',
	values: [
		{ actionId: ActionId.fromSpellId(42463), value: PaladinSeal.Vengeance },
		{ actionId: ActionId.fromSpellId(20154), value: PaladinSeal.Righteousness },
		{
			actionId: ActionId.fromSpellId(20424), value: PaladinSeal.Command,
			showWhen: (player: Player<Spec.SpecRetributionPaladin>) => player.getTalents().sealOfCommand,
		},
	],
	changeEmitter: (player: Player<Spec.SpecRetributionPaladin>) => player.changeEmitter,
});

export const JudgementSelection = InputHelpers.makeSpecOptionsEnumIconInput<Spec.SpecRetributionPaladin, PaladinJudgement>({
	fieldName: 'judgement',
	values: [
		{ actionId: ActionId.fromSpellId(53408), value: PaladinJudgement.JudgementOfWisdom },
		{ actionId: ActionId.fromSpellId(20271), value: PaladinJudgement.JudgementOfLight },
	],
});