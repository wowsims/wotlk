import { Spec } from '/wotlk/core/proto/common.js';
import { ActionId } from '/wotlk/core/proto_utils/action_id.js';
import { Player } from '/wotlk/core/player.js';
import { Sim } from '/wotlk/core/sim.js';
import { EventID, TypedEvent } from '/wotlk/core/typed_event.js';
import { IndividualSimUI } from '/wotlk/core/individual_sim_ui.js';
import { Target } from '/wotlk/core/target.js';
import { getEnumValues } from '/wotlk/core/utils.js';
import { ItemSlot } from '/wotlk/core/proto/common.js';

import {
	FeralTankDruid,
	FeralTankDruid_Rotation as DruidRotation,
	FeralTankDruid_Rotation_Swipe as Swipe,
	FeralTankDruid_Options as DruidOptions
} from '/wotlk/core/proto/druid.js';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const StartingRage = {
	type: 'number' as const,
	getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
	config: {
		extraCssClasses: [
			'starting-rage-picker',
		],
		label: 'Starting Rage',
		labelTooltip: 'Initial rage at the start of each iteration.',
		changedEvent: (player: Player<Spec.SpecFeralTankDruid>) => player.specOptionsChangeEmitter,
		getValue: (player: Player<Spec.SpecFeralTankDruid>) => player.getSpecOptions().startingRage,
		setValue: (eventID: EventID, player: Player<Spec.SpecFeralTankDruid>, newValue: number) => {
			const newOptions = player.getSpecOptions();
			newOptions.startingRage = newValue;
			player.setSpecOptions(eventID, newOptions);
		},
	},
};

export const FeralTankDruidRotationConfig = {
	inputs: [
		{
			type: 'number' as const, cssClass: 'maul-threshold-picker',
			getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
			config: {
				label: 'Maul Threshold',
				labelTooltip: 'Queue Maul when rage is above this value.',
				changedEvent: (player: Player<Spec.SpecFeralTankDruid>) => player.rotationChangeEmitter,
				getValue: (player: Player<Spec.SpecFeralTankDruid>) => player.getRotation().maulRageThreshold,
				setValue: (eventID: EventID, player: Player<Spec.SpecFeralTankDruid>, newValue: number) => {
					const newRotation = player.getRotation();
					newRotation.maulRageThreshold = newValue;
					player.setRotation(eventID, newRotation);
				},
			},
		},
		{
			type: 'enum' as const, cssClass: 'swipe-picker',
			getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
			config: {
				label: 'Swipe',
				values: [
					{ name: 'Never', value: Swipe.SwipeNone },
					{ name: 'With Enough AP', value: Swipe.SwipeWithEnoughAP },
					{ name: 'Spam', value: Swipe.SwipeSpam },
				],
				changedEvent: (player: Player<Spec.SpecFeralTankDruid>) => player.rotationChangeEmitter,
				getValue: (player: Player<Spec.SpecFeralTankDruid>) => player.getRotation().swipe,
				setValue: (eventID: EventID, player: Player<Spec.SpecFeralTankDruid>, newValue: number) => {
					const newRotation = player.getRotation();
					newRotation.swipe = newValue;
					player.setRotation(eventID, newRotation);
				},
			},
		},
		{
			type: 'number' as const, cssClass: 'swipe-ap-threshold-picker',
			getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
			config: {
				label: 'Swipe AP Threshold',
				labelTooltip: 'Use Swipe when Attack Power is larger than this amount.',
				changedEvent: (player: Player<Spec.SpecFeralTankDruid>) => player.rotationChangeEmitter,
				getValue: (player: Player<Spec.SpecFeralTankDruid>) => player.getRotation().swipeApThreshold,
				setValue: (eventID: EventID, player: Player<Spec.SpecFeralTankDruid>, newValue: number) => {
					const newRotation = player.getRotation();
					newRotation.swipeApThreshold = newValue;
					player.setRotation(eventID, newRotation);
				},
				enableWhen: (player: Player<Spec.SpecFeralTankDruid>) => player.getRotation().swipe == Swipe.SwipeWithEnoughAP,
			},
		},
		{
			type: 'boolean' as const,
			cssClass: 'maintain-demo-roar-picker',
			getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
			config: {
				label: 'Maintain Demo Roar',
				labelTooltip: 'Keep Demoralizing Roar active on the primary target. If a stronger debuff is active, will not cast.',
				changedEvent: (player: Player<Spec.SpecFeralTankDruid>) => player.rotationChangeEmitter,
				getValue: (player: Player<Spec.SpecFeralTankDruid>) => player.getRotation().maintainDemoralizingRoar,
				setValue: (eventID: EventID, player: Player<Spec.SpecFeralTankDruid>, newValue: boolean) => {
					const newRotation = player.getRotation();
					newRotation.maintainDemoralizingRoar = newValue;
					player.setRotation(eventID, newRotation);
				},
			},
		},
		{
			type: 'boolean' as const,
			cssClass: 'maintain-faerie-fire-picker',
			getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
			config: {
				label: 'Maintain Faerie Fire',
				labelTooltip: 'Keep Faerie Fire active on the primary target. If a stronger debuff is active, will not cast.',
				changedEvent: (player: Player<Spec.SpecFeralTankDruid>) => player.rotationChangeEmitter,
				getValue: (player: Player<Spec.SpecFeralTankDruid>) => player.getRotation().maintainFaerieFire,
				setValue: (eventID: EventID, player: Player<Spec.SpecFeralTankDruid>, newValue: boolean) => {
					const newRotation = player.getRotation();
					newRotation.maintainFaerieFire = newValue;
					player.setRotation(eventID, newRotation);
				},
			},
		},
	],
};
