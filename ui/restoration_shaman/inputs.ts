import { IconPickerConfig } from '../core/components/icon_picker.js';
import { Spec } from '../core/proto/common.js';
import { ActionId } from '../core/proto_utils/action_id.js';
import { Player } from '../core/player.js';
import { Sim } from '../core/sim.js';
import { Target } from '../core/target.js';
import { EventID, TypedEvent } from '../core/typed_event.js';

import {
	AirTotem,
	RestorationShaman_Options as ShamanOptions,
	ShamanHealSpell,
	ShamanShield,
} from '../core/proto/shaman.js';

import * as InputHelpers from '../core/components/input_helpers.js';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const Bloodlust = InputHelpers.makeSpecOptionsBooleanIconInput<Spec.SpecRestorationShaman>({
	fieldName: 'bloodlust',
	id: ActionId.fromSpellId(2825),
});
export const ShamanShieldInput = InputHelpers.makeSpecOptionsEnumIconInput<Spec.SpecRestorationShaman, ShamanShield>({
	fieldName: 'shield',
	values: [
		{ value: ShamanShield.NoShield, tooltip: 'No Shield' },
		{ actionId: ActionId.fromSpellId(57960), value: ShamanShield.WaterShield },
		{ actionId: ActionId.fromSpellId(49281), value: ShamanShield.LightningShield },
	],
});


export const PrimaryHealInput = InputHelpers.makeRotationEnumInput<Spec.SpecRestorationShaman, ShamanHealSpell>({
	fieldName: 'primaryHeal',
	label: 'Primary Heal',
	labelTooltip: 'Set to \'AutoHeal\', to automatically swap based on best heal.',
	values: [
		{
			name: "Auto Heal",
			value: ShamanHealSpell.AutoHeal
		},
		{
			name: "Lesser Healing Wave",
			value: ShamanHealSpell.LesserHealingWave // actionId: ActionId.fromSpellId(49276),
		},
		{
			name: "Healing Wave",
			value: ShamanHealSpell.HealingWave // actionId: ActionId.fromSpellId(49273),
		},
		{
			name: "Chain Heal",
			value: ShamanHealSpell.ChainHeal // actionId: ActionId.fromSpellId(55459),
		},
	]
});


export const UseRiptide = InputHelpers.makeRotationBooleanInput<Spec.SpecRestorationShaman>({
	fieldName: 'useRiptide',
	label: 'Use Riptide',
	labelTooltip: 'Causes riptide to be cast on primary target when CD is available and not already on.',
	showWhen: (player: Player<Spec.SpecRestorationShaman>) => player.getTalents().riptide,
	changeEmitter: (player: Player<Spec.SpecRestorationShaman>) => TypedEvent.onAny([player.specOptionsChangeEmitter, player.rotationChangeEmitter, player.talentsChangeEmitter]),
});

export const UseEarthShield = InputHelpers.makeRotationBooleanInput<Spec.SpecRestorationShaman>({
	fieldName: 'useEarthShield',
	label: 'Use Earth Shield',
	labelTooltip: 'Causes earth shield to be cast on healing target.',
	showWhen: (player: Player<Spec.SpecRestorationShaman>) => player.getTalents().earthShield,
	changeEmitter: (player: Player<Spec.SpecRestorationShaman>) => TypedEvent.onAny([player.specOptionsChangeEmitter, player.rotationChangeEmitter, player.talentsChangeEmitter]),
});

export const TriggerEarthShield = InputHelpers.makeSpecOptionsNumberInput<Spec.SpecRestorationShaman>({
	fieldName: 'earthShieldPPM',
	label: 'Earth Shield PPM',
	labelTooltip: 'How many times Earth Shield should be triggered per minute.',
	showWhen: (player: Player<Spec.SpecRestorationShaman>) => player.getTalents().earthShield && player.getRotation().useEarthShield,
	changeEmitter: (player: Player<Spec.SpecRestorationShaman>) => TypedEvent.onAny([player.specOptionsChangeEmitter, player.rotationChangeEmitter, player.talentsChangeEmitter]),
});

export const RestorationShamanRotationConfig = {
	inputs: [
		PrimaryHealInput,
		UseRiptide,
		UseEarthShield,
		TriggerEarthShield,
	],
};

