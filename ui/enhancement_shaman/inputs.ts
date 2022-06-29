import { BooleanPicker } from '/tbc/core/components/boolean_picker.js';
import { EnumPicker } from '/tbc/core/components/enum_picker.js';
import { IconEnumPicker, IconEnumPickerConfig } from '/tbc/core/components/icon_enum_picker.js';
import { IconPickerConfig } from '/tbc/core/components/icon_picker.js';
import { makeWeaponImbueInput } from '/tbc/core/components/icon_inputs.js';
import {
	AirTotem,
	EarthTotem,
	FireTotem,
	WaterTotem,
	EnhancementShaman_Rotation_PrimaryShock as PrimaryShock,
	ShamanTotems,
} from '/tbc/core/proto/shaman.js';
import { EnhancementShaman_Options as ShamanOptions } from '/tbc/core/proto/shaman.js';
import { Spec } from '/tbc/core/proto/common.js';
import { WeaponImbue } from '/tbc/core/proto/common.js';
import { ActionId } from '/tbc/core/proto_utils/action_id.js';
import { Player } from '/tbc/core/player.js';
import { Sim } from '/tbc/core/sim.js';
import { IndividualSimUI } from '/tbc/core/individual_sim_ui.js';
import { Target } from '/tbc/core/target.js';
import { EventID, TypedEvent } from '/tbc/core/typed_event.js';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const IconBloodlust = makeBooleanShamanBuffInput(ActionId.fromSpellId(2825), 'bloodlust');
export const IconWaterShield = makeBooleanShamanBuffInput(ActionId.fromSpellId(33736), 'waterShield');

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

export const SnapshotT42Pc = {
	type: 'boolean' as const,
	getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
	config: {
		extraCssClasses: [
			'snapshot-t4-2pc-picker',
		],
		label: 'Snapshot T4 2pc',
		labelTooltip: 'Snapshots the improved Strength of Earth totem bonus from T4 2pc (+12 strength) for the first 1:50s of the fight. Only works if the selected Earth totem is Strength of Earth Totem.',
		changedEvent: (player: Player<Spec.SpecEnhancementShaman>) => player.changeEmitter,
		getValue: (player: Player<Spec.SpecEnhancementShaman>) => player.getSpecOptions().snapshotT42Pc,
		setValue: (eventID: EventID, player: Player<Spec.SpecEnhancementShaman>, newValue: boolean) => {
			const newOptions = player.getSpecOptions();
			newOptions.snapshotT42Pc = newValue;
			player.setSpecOptions(eventID, newOptions);
		},
		enableWhen: (player: Player<Spec.SpecEnhancementShaman>) => player.getRotation().totems?.earth == EarthTotem.StrengthOfEarthTotem,
	},
};

export const EnhancementShamanRotationConfig = {
	inputs: [
		{
			type: 'enum' as const, cssClass: 'primary-shock-picker',
			getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
			config: {
				label: 'Primary Shock',
				values: [
					{
						name: 'None', value: PrimaryShock.None,
					},
					{
						name: 'Earth Shock', value: PrimaryShock.Earth,
					},
					{
						name: 'Frost Shock', value: PrimaryShock.Frost,
					},
				],
				changedEvent: (player: Player<Spec.SpecEnhancementShaman>) => player.rotationChangeEmitter,
				getValue: (player: Player<Spec.SpecEnhancementShaman>) => player.getRotation().primaryShock,
				setValue: (eventID: EventID, player: Player<Spec.SpecEnhancementShaman>, newValue: number) => {
					const newRotation = player.getRotation();
					newRotation.primaryShock = newValue;
					player.setRotation(eventID, newRotation);
				},
			},
		},
		{
			type: 'boolean' as const, cssClass: 'weave-flame-shock-picker',
			getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
			config: {
				label: 'Weave Flame Shock',
				labelTooltip: 'Use Flame Shock whenever the target does not already have the DoT.',
				changedEvent: (player: Player<Spec.SpecEnhancementShaman>) => player.rotationChangeEmitter,
				getValue: (player: Player<Spec.SpecEnhancementShaman>) => player.getRotation().weaveFlameShock,
				setValue: (eventID: EventID, player: Player<Spec.SpecEnhancementShaman>, newValue: boolean) => {
					const newRotation = player.getRotation();
					newRotation.weaveFlameShock = newValue;
					player.setRotation(eventID, newRotation);
				},
			},
		},
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
