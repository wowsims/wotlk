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
	],
});

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.
export const ProtectionPaladinRotationConfig = {
	inputs: [
		InputHelpers.makeRotationNumberInput<Spec.SpecProtectionPaladin>({
			fieldName: 'prioritizeHolyShield',
			label: 'Prio Holy Shield',
			labelTooltip: 'Uses Holy Shield as the highest priority spell. This is usually done when tanking a boss that can crush.',
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

export const DamageTakenPerSecond = InputHelpers.makeSpecOptionsNumberInput<Spec.SpecProtectionPaladin>({
	fieldName: 'damageTakenPerSecond',
	label: 'Damage Taken Per Second',
	labelTooltip: "Damage taken per second across the encounter. Used to model mana regeneration from Spiritual Attunement. This value should NOT include damage taken from Seal of Blood / Judgement of Blood. Leave at 0 if you do not take damage during the encounter.",
});
