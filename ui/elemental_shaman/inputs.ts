import { ElementalShaman_Options_ThunderstormRange, ShamanShield } from '../core/proto/shaman.js';
import { Spec } from '../core/proto/common.js';
import { ActionId } from '../core/proto_utils/action_id.js';
import { Player } from '../core/player.js';

import * as InputHelpers from '../core/components/input_helpers.js';
import { EventID } from 'ui/core/typed_event.js';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const InThunderstormRange = InputHelpers.makeSpecOptionsBooleanInput<Spec.SpecElementalShaman>({
	fieldName: 'thunderstormRange',
	// id: ActionId.fromSpellId(59159),
	label: "Thunderstorm In Range",
	labelTooltip: "When set to true, thunderstorm casts will cause damage.",
	getValue: (player: Player<Spec.SpecElementalShaman>) => player.getSpecOptions().thunderstormRange == ElementalShaman_Options_ThunderstormRange.TSInRange,
	setValue: (eventID: EventID, player: Player<Spec.SpecElementalShaman>, newValue: boolean) => {
		const newOptions = player.getSpecOptions();
		if (newValue) {
			newOptions.thunderstormRange = ElementalShaman_Options_ThunderstormRange.TSInRange;
		} else {
			newOptions.thunderstormRange = ElementalShaman_Options_ThunderstormRange.TSOutofRange;
		}
		player.setSpecOptions(eventID, newOptions);
	},
});

export const ShamanShieldInput = InputHelpers.makeSpecOptionsEnumIconInput<Spec.SpecElementalShaman, ShamanShield>({
	fieldName: 'shield',
	values: [
		{ value: ShamanShield.NoShield, tooltip: 'No Shield' },
		{ actionId: ActionId.fromSpellId(57960), value: ShamanShield.WaterShield },
		{ actionId: ActionId.fromSpellId(49281), value: ShamanShield.LightningShield },
	],
});
