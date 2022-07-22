import { BooleanPicker } from '/wotlk/core/components/boolean_picker.js';
import { EnumPicker } from '/wotlk/core/components/enum_picker.js';
import { IconEnumPicker, IconEnumPickerConfig } from '/wotlk/core/components/icon_enum_picker.js';
import { IconPickerConfig } from '/wotlk/core/components/icon_picker.js';
import { makeWeaponImbueInput } from '/wotlk/core/components/icon_inputs.js';
import {
	AirTotem,
	EarthTotem,
	FireTotem,
	WaterTotem,
	EnhancementShaman_Options as ShamanOptions,
	ShamanTotems,
	ShamanShield
} from '/wotlk/core/proto/shaman.js';
import { Spec } from '/wotlk/core/proto/common.js';
import { WeaponImbue } from '/wotlk/core/proto/common.js';
import { ActionId } from '/wotlk/core/proto_utils/action_id.js';
import { Player } from '/wotlk/core/player.js';
import { Sim } from '/wotlk/core/sim.js';
import { IndividualSimUI } from '/wotlk/core/individual_sim_ui.js';
import { Target } from '/wotlk/core/target.js';
import { EventID, TypedEvent } from '/wotlk/core/typed_event.js';

import * as InputHelpers from '/wotlk/core/components/input_helpers.js';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const Bloodlust = InputHelpers.makeSpecOptionsBooleanIconInput<Spec.SpecEnhancementShaman>({
	fieldName: 'bloodlust',
	id: ActionId.fromSpellId(2825),
});
export const ShamanShieldInput = InputHelpers.makeSpecOptionsEnumIconInput<Spec.SpecEnhancementShaman, ShamanShield>({
	fieldName: 'shield',
	values: [
		{ color: 'grey', value: ShamanShield.NoShield },
		{ actionId: ActionId.fromItemId(33736), value: ShamanShield.WaterShield },
		{ actionId: ActionId.fromItemId(49281), value: ShamanShield.LightningShield },
	],
});

export const DelayOffhandSwings = InputHelpers.makeSpecOptionsBooleanInput<Spec.SpecEnhancementShaman>({
	fieldName: 'delayOffhandSwings',
	label: 'Delay Offhand Swings',
	labelTooltip: 'Uses the startattack macro to delay OH swings, so they always follow within 0.5s of a MH swing.',
});

export const EnhancementShamanRotationConfig = {
	inputs: [
//		{
//			type: 'enum' as const, cssClass: 'primary-shock-picker',
//			getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
//			config: {
//				label: 'Mainhand Imbue', //very temporary, just as a way to be able to make sure imbues are working in the meantime,
//				values: [                //and primary shocks arent a thing anymore
//					{
//						name: 'None', value: WeaponImbue.None,
//					},
//					{
//						name: 'Windfury', value: WeaponImbue.WeaponImbueShamanWindfury,
//					},
//					{
//						name: 'Flametongue', value: WeaponImbue.WeaponImbueShamanFlametongue,
//					},
//				],
//				changedEvent: (player: Player<Spec.SpecEnhancementShaman>) => player.rotationChangeEmitter,
//				getValue: (player: Player<Spec.SpecEnhancementShaman>) => player.getRotation().WeaponImbue,
//				setValue: (eventID: EventID, player: Player<Spec.SpecEnhancementShaman>, newValue: number) => {
//					const newRotation = player.getRotation();
//					newRotation.WeaponImbue = newValue;
//					player.setRotation(eventID, newRotation);
//				},
//			},
//		}
	],
};
