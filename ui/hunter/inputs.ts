import { BooleanPicker } from '/wotlk/core/components/boolean_picker.js';
import { EnumPicker } from '/wotlk/core/components/enum_picker.js';
import { IconEnumPicker, IconEnumPickerConfig } from '/wotlk/core/components/icon_enum_picker.js';
import { IconPickerConfig } from '/wotlk/core/components/icon_picker.js';
import { Spec } from '/wotlk/core/proto/common.js';
import { ActionId } from '/wotlk/core/proto_utils/action_id.js';
import { Player } from '/wotlk/core/player.js';
import { Sim } from '/wotlk/core/sim.js';
import { Target } from '/wotlk/core/target.js';
import { EventID, TypedEvent } from '/wotlk/core/typed_event.js';
import { makePetTypeInputConfig } from '/wotlk/core/talents/hunter_pet.js';

import * as InputHelpers from '/wotlk/core/components/input_helpers.js';

import {
	Hunter,
	Hunter_Rotation as HunterRotation,
	Hunter_Rotation_StingType as StingType,
	//Hunter_Rotation_WeaveType as WeaveType,
	Hunter_Options as HunterOptions,
	Hunter_Options_Ammo as Ammo,
	Hunter_Options_PetType as PetType,
} from '/wotlk/core/proto/hunter.js';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const WeaponAmmo = InputHelpers.makeSpecOptionsEnumIconInput<Spec.SpecHunter, Ammo>({
	fieldName: 'ammo',
	numColumns: 2,
	values: [
		{ color: 'grey', value: Ammo.AmmoNone },
		{ actionId: ActionId.fromItemId(52021), value: Ammo.IcebladeArrow },
		{ actionId: ActionId.fromItemId(41165), value: Ammo.SaroniteRazorheads },
		{ actionId: ActionId.fromItemId(41586), value: Ammo.TerrorshaftArrow },
		{ actionId: ActionId.fromItemId(31737), value: Ammo.TimelessArrow },
		{ actionId: ActionId.fromItemId(34581), value: Ammo.MysteriousArrow },
		{ actionId: ActionId.fromItemId(33803), value: Ammo.AdamantiteStinger },
		{ actionId: ActionId.fromItemId(28056), value: Ammo.BlackflightArrow },
	],
});

export const LatencyMs = InputHelpers.makeSpecOptionsNumberInput<Spec.SpecHunter>({
	fieldName: 'latencyMs',
	label: 'Latency',
	labelTooltip: 'Player latency, in milliseconds. Adds a delay to actions other than auto shot.',
});

export const PetTypeInput = makePetTypeInputConfig(true);

export const PetUptime = InputHelpers.makeSpecOptionsNumberInput<Spec.SpecHunter>({
	fieldName: 'petUptime',
	label: 'Pet Uptime (%)',
	labelTooltip: 'Percent of the fight duration for which your pet will be alive.',
	percent: true,
});

//export const PetSingleAbility = {
//	type: 'boolean' as const,
//	getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
//	config: {
//		extraCssClasses: [
//			'pet-single-ability-picker',
//		],
//		label: 'Single Pet Ability',
//		labelTooltip: 'Pet will only use its primary ability.',
//		changedEvent: (player: Player<Spec.SpecHunter>) => player.specOptionsChangeEmitter,
//		getValue: (player: Player<Spec.SpecHunter>) => player.getSpecOptions().petSingleAbility,
//		setValue: (eventID: EventID, player: Player<Spec.SpecHunter>, newValue: boolean) => {
//			const newOptions = player.getSpecOptions();
//			newOptions.petSingleAbility = newValue;
//			player.setSpecOptions(eventID, newOptions);
//		},
//	},
//};

export const SniperTrainingUptime = InputHelpers.makeSpecOptionsNumberInput<Spec.SpecHunter>({
	fieldName: 'sniperTrainingUptime',
	label: 'ST Uptime (%)',
	labelTooltip: 'Uptime for the Sniper Training talent, as a percent of the fight duration.',
	percent: true,
	showWhen: (player: Player<Spec.SpecHunter>) => player.getTalents().sniperTraining > 0,
});

export const HunterRotationConfig = {
	inputs: [
		//{
		//	type: 'boolean' as const, cssClass: 'use-multi-shot-picker',
		//	getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
		//	config: {
		//		label: 'Use Multi Shot',
		//		labelTooltip: 'Includes Multi Shot in the rotation.',
		//		changedEvent: (player: Player<Spec.SpecHunter>) => player.rotationChangeEmitter,
		//		getValue: (player: Player<Spec.SpecHunter>) => player.getRotation().useMultiShot,
		//		setValue: (eventID: EventID, player: Player<Spec.SpecHunter>, newValue: boolean) => {
		//			const newRotation = player.getRotation();
		//			newRotation.useMultiShot = newValue;
		//			player.setRotation(eventID, newRotation);
		//		},
		//	},
		//},
		//{
		//	type: 'boolean' as const, cssClass: 'use-arcane-shot-picker',
		//	getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
		//	config: {
		//		label: 'Use Arcane Shot',
		//		labelTooltip: 'Includes Arcane Shot in the rotation.',
		//		changedEvent: (player: Player<Spec.SpecHunter>) => player.rotationChangeEmitter,
		//		getValue: (player: Player<Spec.SpecHunter>) => player.getRotation().useArcaneShot,
		//		setValue: (eventID: EventID, player: Player<Spec.SpecHunter>, newValue: boolean) => {
		//			const newRotation = player.getRotation();
		//			newRotation.useArcaneShot = newValue;
		//			player.setRotation(eventID, newRotation);
		//		},
		//	},
		//},
		InputHelpers.makeRotationEnumInput<Spec.SpecHunter, StingType>({
			fieldName: 'sting',
			label: 'Sting',
			labelTooltip: 'Maintains the selected Sting on the primary target.',
			values: [
				{ name: 'None', value: StingType.NoSting },
				{ name: 'Scorpid Sting', value: StingType.ScorpidSting },
				{ name: 'Serpent Sting', value: StingType.SerpentSting },
			],
		}),
		//{
		//	type: 'boolean' as const, cssClass: 'lazy-rotation-picker',
		//	getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
		//	config: {
		//		label: 'Lazy Rotation',
		//		labelTooltip: 'Uses GCD immediately, even if it will clip the next auto.',
		//		changedEvent: (player: Player<Spec.SpecHunter>) => player.rotationChangeEmitter,
		//		getValue: (player: Player<Spec.SpecHunter>) => player.getRotation().lazyRotation,
		//		setValue: (eventID: EventID, player: Player<Spec.SpecHunter>, newValue: boolean) => {
		//			const newRotation = player.getRotation();
		//			newRotation.lazyRotation = newValue;
		//			player.setRotation(eventID, newRotation);
		//		},
		//	},
		//},
		InputHelpers.makeRotationNumberInput<Spec.SpecHunter>({
			fieldName: 'viperStartManaPercent',
			label: 'Viper Start Mana %',
			labelTooltip: 'Switch to Aspect of the Viper when mana goes below this amount.',
			percent: true,
		}),
		InputHelpers.makeRotationNumberInput<Spec.SpecHunter>({
			fieldName: 'viperStopManaPercent',
			label: 'Viper Stop Mana %',
			labelTooltip: 'Switch back to Aspect of the Hawk when mana goes above this amount.',
			percent: true,
		}),
		//{
		//	type: 'enum' as const, cssClass: 'weave-picker',
		//	getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
		//	config: {
		//		label: 'Melee Weaving',
		//		labelTooltip: 'Uses melee weaving in the rotation.',
		//		values: [
		//			{ name: 'None', value: WeaveType.WeaveNone },
		//			{ name: 'Autos Only', value: WeaveType.WeaveAutosOnly },
		//			{ name: 'Raptor Only', value: WeaveType.WeaveRaptorOnly },
		//			{ name: 'Full', value: WeaveType.WeaveFull },
		//		],
		//		changedEvent: (player: Player<Spec.SpecHunter>) => player.rotationChangeEmitter,
		//		getValue: (player: Player<Spec.SpecHunter>) => player.getRotation().weave,
		//		setValue: (eventID: EventID, player: Player<Spec.SpecHunter>, newValue: number) => {
		//			const newRotation = player.getRotation();
		//			newRotation.weave = newValue;
		//			player.setRotation(eventID, newRotation);
		//		},
		//	},
		//},
		//{
		//	type: 'number' as const, cssClass: 'time-to-weave-picker',
		//	getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
		//	config: {
		//		label: 'Time To Weave (ms)',
		//		labelTooltip: 'Amount of time, in milliseconds, between when you start moving towards the boss and when you re-engage your ranged autos.',
		//		changedEvent: (player: Player<Spec.SpecHunter>) => player.rotationChangeEmitter,
		//		getValue: (player: Player<Spec.SpecHunter>) => player.getRotation().timeToWeaveMs,
		//		setValue: (eventID: EventID, player: Player<Spec.SpecHunter>, newValue: number) => {
		//			const newRotation = player.getRotation();
		//			newRotation.timeToWeaveMs = newValue;
		//			player.setRotation(eventID, newRotation);
		//		},
		//		showWhen: (player: Player<Spec.SpecHunter>) => player.getRotation().weave != WeaveType.WeaveNone,
		//	},
		//},
		//{
		//	type: 'number' as const, cssClass: 'percent-weaved-picker',
		//	getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
		//	config: {
		//		label: 'Time Weaved (%)',
		//		labelTooltip: 'Percentage of fight to use melee weaving.',
		//		changedEvent: (player: Player<Spec.SpecHunter>) => player.rotationChangeEmitter,
		//		getValue: (player: Player<Spec.SpecHunter>) => player.getRotation().percentWeaved * 100,
		//		setValue: (eventID: EventID, player: Player<Spec.SpecHunter>, newValue: number) => {
		//			const newRotation = player.getRotation();
		//			newRotation.percentWeaved = newValue / 100;
		//			player.setRotation(eventID, newRotation);
		//		},
		//		showWhen: (player: Player<Spec.SpecHunter>) => player.getRotation().weave != WeaveType.WeaveNone,
		//	},
		//},
	],
};
