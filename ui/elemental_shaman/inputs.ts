import { IconPickerConfig } from '/tbc/core/components/icon_picker.js';
import { ElementalShaman_Rotation_RotationType as RotationType } from '/tbc/core/proto/shaman.js';
import { ElementalShaman_Options as ShamanOptions } from '/tbc/core/proto/shaman.js';
import { AirTotem } from '/tbc/core/proto/shaman.js';
import { Spec } from '/tbc/core/proto/common.js';
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

export const SnapshotT42Pc = {
	type: 'boolean' as const,
	getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
	config: {
		extraCssClasses: [
			'snapshot-t4-2pc-picker',
		],
		label: 'Snapshot T4 2pc',
		labelTooltip: 'Snapshots the improved wrath of air totem bonus from T4 2pc (+20 spell power) for the first 1:50s of the fight. Only works if the selected air totem is Wrath of Air Totem.',
		changedEvent: (player: Player<Spec.SpecElementalShaman>) => player.changeEmitter,
		getValue: (player: Player<Spec.SpecElementalShaman>) => player.getSpecOptions().snapshotT42Pc,
		setValue: (eventID: EventID, player: Player<Spec.SpecElementalShaman>, newValue: boolean) => {
			const newOptions = player.getSpecOptions();
			newOptions.snapshotT42Pc = newValue;
			player.setSpecOptions(eventID, newOptions);
		},
		enableWhen: (player: Player<Spec.SpecElementalShaman>) => player.getRotation().totems?.air == AirTotem.WrathOfAirTotem,
	},
};

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
					{
						name: 'CL On Clearcast', value: RotationType.CLOnClearcast,
						tooltip: 'Casts CL only after Clearcast procs.',
					},
					{
						name: 'CL On CD', value: RotationType.CLOnCD,
						tooltip: 'Casts CL if it is ready, otherwise LB.',
					},
					{
						name: 'Fixed LB+CL', value: RotationType.FixedLBCL,
						tooltip: 'Casts a fixed number of LBs between each CL (specified below), even if that means waiting. While temporary haste effects are active (drums, lust, etc) will cast extra LBs instead of waiting.',
					},
					{
						name: 'LB Only', value: RotationType.LBOnly,
						tooltip: 'Only casts Lightning Bolt.',
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
			type: 'number' as const,
			cssClass: 'num-lbs-per-cl-picker',
			getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
			config: {
				label: 'LBs per CL',
				labelTooltip: 'The number of Lightning Bolts to cast between each Chain Lightning. Only used if Rotation is set to \'Fixed LB+CL\'.',
				changedEvent: (player: Player<Spec.SpecElementalShaman>) => player.rotationChangeEmitter,
				getValue: (player: Player<Spec.SpecElementalShaman>) => player.getRotation().lbsPerCl,
				setValue: (eventID: EventID, player: Player<Spec.SpecElementalShaman>, newValue: number) => {
					const newRotation = player.getRotation();
					newRotation.lbsPerCl = newValue;
					player.setRotation(eventID, newRotation);
				},
				enableWhen: (player: Player<Spec.SpecElementalShaman>) => player.getRotation().type == RotationType.FixedLBCL,
			},
		},
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
