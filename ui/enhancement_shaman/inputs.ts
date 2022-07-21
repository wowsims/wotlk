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
	ShamanTotems,
	ShamanShield
} from '/wotlk/core/proto/shaman.js';
import { EnhancementShaman_Options as ShamanOptions } from '/wotlk/core/proto/shaman.js';
import { Spec } from '/wotlk/core/proto/common.js';
import { WeaponImbue } from '/wotlk/core/proto/common.js';
import { ActionId } from '/wotlk/core/proto_utils/action_id.js';
import { Player } from '/wotlk/core/player.js';
import { Sim } from '/wotlk/core/sim.js';
import { IndividualSimUI } from '/wotlk/core/individual_sim_ui.js';
import { Target } from '/wotlk/core/target.js';
import { EventID, TypedEvent } from '/wotlk/core/typed_event.js';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const IconBloodlust = makeBooleanShamanBuffInput(ActionId.fromSpellId(2825), 'bloodlust');

export const IconLightningShield = {
	id: ActionId.fromSpellId(49281),
	states: 2,
	changedEvent: (player: Player<Spec.SpecEnhancementShaman>) => player.specOptionsChangeEmitter,
	getValue: (player: Player<Spec.SpecEnhancementShaman>) => player.getSpecOptions().shield == ShamanShield.LightningShield,
	setValue: (eventID: EventID, player: Player<Spec.SpecEnhancementShaman>, newValue: boolean) => {
		const newOptions = player.getSpecOptions();
		newOptions.shield = ShamanShield.LightningShield;
		player.setSpecOptions(eventID, newOptions);
	},
}

export const IconWaterShield = {
	id: ActionId.fromSpellId(57960),
	states: 2,
	changedEvent: (player: Player<Spec.SpecEnhancementShaman>) => player.specOptionsChangeEmitter,
	getValue: (player: Player<Spec.SpecEnhancementShaman>) => player.getSpecOptions().shield == ShamanShield.WaterShield,
	setValue: (eventID: EventID, player: Player<Spec.SpecEnhancementShaman>, newValue: boolean) => {
		const newOptions = player.getSpecOptions();
		newOptions.shield = ShamanShield.WaterShield;
		player.setSpecOptions(eventID, newOptions);
	},
}


export const DelayOffhandSwings = {
	type: 'boolean' as const,
	getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
	config: {
		extraCssClasses: [
			'delay-offhand-swings-picker',
		],
		label: 'Delay Offhand Swings',
		labelTooltip: 'Uses the startattack macro to delay OH swings, so they always follow within 0.5s of a MH swing.',
		changedEvent: (player: Player<Spec.SpecEnhancementShaman>) => player.specOptionsChangeEmitter,
		getValue: (player: Player<Spec.SpecEnhancementShaman>) => player.getSpecOptions().delayOffhandSwings,
		setValue: (eventID: EventID, player: Player<Spec.SpecEnhancementShaman>, newValue: boolean) => {
			const newOptions = player.getSpecOptions();
			newOptions.delayOffhandSwings = newValue;
			player.setSpecOptions(eventID, newOptions);
		},
	},
};

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

function makeBooleanShamanBuffInput(id: ActionId, optionsFieldName: keyof ShamanOptions): IconPickerConfig<Player<any>, boolean> {
	return {
		id: id,
		states: 2,
		changedEvent: (player: Player<Spec.SpecEnhancementShaman>) => player.specOptionsChangeEmitter,
		getValue: (player: Player<Spec.SpecEnhancementShaman>) => player.getSpecOptions()[optionsFieldName] as boolean,
		setValue: (eventID: EventID, player: Player<Spec.SpecEnhancementShaman>, newValue: boolean) => {
			const newOptions = player.getSpecOptions();
			(newOptions[optionsFieldName] as boolean) = newValue;
			player.setSpecOptions(eventID, newOptions);
		},
	};
}
