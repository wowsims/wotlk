import { IconPickerConfig } from '/wotlk/core/components/icon_picker.js';
import { ElementalShaman_Rotation_RotationType as RotationType, ShamanShield } from '/wotlk/core/proto/shaman.js';
import { ElementalShaman_Options as ShamanOptions } from '/wotlk/core/proto/shaman.js';
import { AirTotem } from '/wotlk/core/proto/shaman.js';
import { Spec } from '/wotlk/core/proto/common.js';
import { ActionId } from '/wotlk/core/proto_utils/action_id.js';
import { Player } from '/wotlk/core/player.js';
import { Sim } from '/wotlk/core/sim.js';
import { Target } from '/wotlk/core/target.js';
import { EventID, TypedEvent } from '/wotlk/core/typed_event.js';

import * as InputHelpers from '/wotlk/core/components/input_helpers.js';

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
		{ actionId: ActionId.fromItemId(33736), value: ShamanShield.WaterShield },
		{ actionId: ActionId.fromItemId(49281), value: ShamanShield.LightningShield },
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
