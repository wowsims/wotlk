import { IconPickerConfig } from '/wotlk/core/components/icon_picker.js';
import { ElementalShaman_Rotation_RotationType as RotationType, ShamanShield } from '/wotlk/core/proto/shaman.js';
import { ElementalShaman_Options as ShamanOptions } from '/wotlk/core/proto/shaman.js';
import { AirTotem } from '/wotlk/core/proto/shaman.js';
import { Spec } from '/wotlk/core/proto/common.js';
import { ActionId } from '/wotlk/core/proto_utils/action_id.js';
import { Player } from '/wotlk/core/player.js';
import { Sim } from '/wotlk/core/sim.js';
import { IndividualSimUI } from '/wotlk/core/individual_sim_ui.js';
import { Target } from '/wotlk/core/target.js';
import { EventID, TypedEvent } from '/wotlk/core/typed_event.js';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const IconBloodlust = makeBooleanShamanBuffInput(ActionId.fromSpellId(2825), 'bloodlust');
export const IconWaterShield = {
	id: ActionId.fromSpellId(33736),
	states: 2,
	changedEvent: (player: Player<Spec.SpecElementalShaman>) => player.specOptionsChangeEmitter,
	getValue: (player: Player<Spec.SpecElementalShaman>) => player.getSpecOptions().shield == ShamanShield.WaterShield,
	setValue: (eventID: EventID, player: Player<Spec.SpecElementalShaman>, newValue: boolean) => {
		const newOptions = player.getSpecOptions();
		if (newValue) {
			newOptions.shield = ShamanShield.WaterShield;
		} else {
			newOptions.shield = ShamanShield.NoShield;
		}
		player.setSpecOptions(eventID, newOptions);
	},
}

export const ElementalShamanRotationConfig = {
	inputs: [
		{
			type: 'enum' as const, cssClass: 'rotation-enum-picker',
			getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
			config: {
				label: 'Type',
				values: [
					{
						name: 'Adaptive', value: RotationType.Adaptive,
						tooltip: 'Dynamically adapts based on available mana to maximize CL casts without going OOM.',
					},
				],
				changedEvent: (player: Player<Spec.SpecElementalShaman>) => player.rotationChangeEmitter,
				getValue: (player: Player<Spec.SpecElementalShaman>) => player.getRotation().type,
				setValue: (eventID: EventID, player: Player<Spec.SpecElementalShaman>, newValue: number) => {
					const newRotation = player.getRotation();
					newRotation.type = newValue;
					player.setRotation(eventID, newRotation);
				},
			},
		},
		{
			type: 'boolean' as const, 
			cssClass: '',
			getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
			config: {
				label: 'In Thunderstorm Range',
				labelTooltip: 'Thunderstorm will hit all targets when cast. Ignores knockback.',
				changedEvent: (player: Player<Spec.SpecElementalShaman>) => player.talentsChangeEmitter,
				getValue: (player: Player<Spec.SpecElementalShaman>) => player.getRotation().inThunderstormRange,
				setValue: (eventID: EventID, player: Player<Spec.SpecElementalShaman>, newValue: boolean) => {
					const newRotation = player.getRotation();
					newRotation.inThunderstormRange = newValue
					player.setRotation(eventID, newRotation);
				},
				enableWhen: (player: Player<Spec.SpecElementalShaman>) => player.getTalents().thunderstorm,
			},			
		}
	],
};

function makeBooleanShamanBuffInput(id: ActionId, optionsFieldName: keyof ShamanOptions): IconPickerConfig<Player<any>, boolean> {
	return {
		id: id,
		states: 2,
		changedEvent: (player: Player<Spec.SpecElementalShaman>) => player.specOptionsChangeEmitter,
		getValue: (player: Player<Spec.SpecElementalShaman>) => player.getSpecOptions()[optionsFieldName] as boolean,
		setValue: (eventID: EventID, player: Player<Spec.SpecElementalShaman>, newValue: boolean) => {
			const newOptions = player.getSpecOptions();
			(newOptions[optionsFieldName] as boolean) = newValue;
			player.setSpecOptions(eventID, newOptions);
		},
	}
}
