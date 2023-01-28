import { Spec } from '../core/proto/common.js';
import { Player } from '../core/player.js';
import { EventID } from '../core/typed_event.js';
import { IndividualSimUI } from '../core/individual_sim_ui.js';
import { ActionId } from '../core/proto_utils/action_id.js';

import {
	PaladinAura as PaladinAura,
	PaladinSeal,
	PaladinJudgement as PaladinJudgement,
	ProtectionPaladin_Rotation_SpellOption as SpellOption,
	ProtectionPaladin_Rotation as ProtectionPaladinRotation,
	ProtectionPaladin_Options as ProtectionPaladinOptions,
} from '../core/proto/paladin.js';

import * as InputHelpers from '../core/components/input_helpers.js';

export const ProtectionPaladinRotationPriorityConfig = InputHelpers.makeCustomRotationInput<Spec.SpecProtectionPaladin, SpellOption>({
	fieldName: 'customRotation',
	numColumns: 2,
	values: [
		{ actionId: ActionId.fromSpellId(53408), value: SpellOption.JudgementOfWisdom },
		{ actionId: ActionId.fromSpellId(48806), value: SpellOption.HammerOfWrath },
		{ actionId: ActionId.fromSpellId(48819), value: SpellOption.Consecration },
		{ actionId: ActionId.fromSpellId(48817), value: SpellOption.HolyWrath },
		{ actionId: ActionId.fromSpellId(48801), value: SpellOption.Exorcism },
		{ actionId: ActionId.fromSpellId(61411), value: SpellOption.ShieldOfRighteousness },
		{ actionId: ActionId.fromSpellId(48827), value: SpellOption.AvengersShield },
		{ actionId: ActionId.fromSpellId(53595), value: SpellOption.HammerOfTheRighteous },
		{ actionId: ActionId.fromSpellId(48952), value: SpellOption.HolyShield },
	],
});

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.
export const ProtectionPaladinRotationConfig = {
	inputs: [
		InputHelpers.makeRotationBooleanInput<Spec.SpecProtectionPaladin>({
			fieldName: 'hammerFirst',
			label: 'Open with HotR',
			labelTooltip: 'Open with Hammer of the Righteous instead of Shield of Righteousness in the standard rotation. Recommended for AoE.',
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecProtectionPaladin>({
			fieldName: 'squeezeHolyWrath',
			label: 'Squeeze Holy Wrath',
			labelTooltip: 'Squeeze a Holy Wrath cast during sufficiently hasted GCDs (Bloodlust) in the standard rotation.',
		}),
		InputHelpers.makeRotationNumberInput<Spec.SpecProtectionPaladin>({
			fieldName: 'waitSlack',
			label: 'Max Wait Time (ms)',
			labelTooltip: 'Maximum time in milliseconds to prioritize waiting for next Hammer/Shield to maintain 969. Affects standard and custom priority.',
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecProtectionPaladin>({
			fieldName: 'useCustomPrio',
			label: 'Use custom priority',
			labelTooltip: 'Deviates from the standard 96969 rotation, using the priority configured below. Will still attempt to keep a filler GCD between Hammer and Shield.',
		}),
		ProtectionPaladinRotationPriorityConfig
	],
}

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
