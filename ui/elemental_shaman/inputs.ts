import { IconPickerConfig } from '../core/components/icon_picker.js';
import { ElementalShaman_Rotation_RotationType as RotationType, ShamanShield } from '../core/proto/shaman.js';
import { ElementalShaman_Options as ShamanOptions } from '../core/proto/shaman.js';
import { AirTotem } from '../core/proto/shaman.js';
import { Spec } from '../core/proto/common.js';
import { ActionId } from '../core/proto_utils/action_id.js';
import { Player } from '../core/player.js';
import { Sim } from '../core/sim.js';
import { Target } from '../core/target.js';
import { EventID, TypedEvent } from '../core/typed_event.js';

import * as InputHelpers from '../core/components/input_helpers.js';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const Bloodlust = InputHelpers.makeSpecOptionsBooleanIconInput<Spec.SpecElementalShaman>({
	fieldName: 'bloodlust',
	id: ActionId.fromSpellId(2825),
});
export const ShamanShieldInput = InputHelpers.makeSpecOptionsEnumIconInput<Spec.SpecElementalShaman, ShamanShield>({
	fieldName: 'shield',
	values: [
		{ color: 'grey', value: ShamanShield.NoShield },
		{ actionId: ActionId.fromSpellId(57960), value: ShamanShield.WaterShield },
		{ actionId: ActionId.fromSpellId(49281), value: ShamanShield.LightningShield },
	],
});

export const ElementalShamanRotationConfig = {
	inputs: [
		InputHelpers.makeRotationEnumInput<Spec.SpecElementalShaman, RotationType>({
			fieldName: 'type',
			label: 'Type',
			values: [
				{
					name: 'Adaptive', value: RotationType.Adaptive,
					tooltip: 'Dynamically adapts based on available mana to maximize CL casts without going OOM.',
				},
			],
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecElementalShaman>({
			fieldName: 'inThunderstormRange',
			label: 'In Thunderstorm Range',
			labelTooltip: 'Thunderstorm will hit all targets when cast. Ignores knockback.',
			enableWhen: (player: Player<Spec.SpecElementalShaman>) => player.getTalents().thunderstorm,
		}),
	],
};
