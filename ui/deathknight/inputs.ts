import { IconPickerConfig } from '/wotlk/core/components/icon_picker.js';
import { RaidTarget } from '/wotlk/core/proto/common.js';
import { Spec } from '/wotlk/core/proto/common.js';
import { NO_TARGET } from '/wotlk/core/proto_utils/utils.js';
import { ActionId } from '/wotlk/core/proto_utils/action_id.js';
import { Player } from '/wotlk/core/player.js';
import { Sim } from '/wotlk/core/sim.js';
import { EventID, TypedEvent } from '/wotlk/core/typed_event.js';
import { IndividualSimUI } from '/wotlk/core/individual_sim_ui.js';
import { Target } from '/wotlk/core/target.js';

import {
	DeathKnightTalents as DeathKnightTalents,
	DeathKnight,
	DeathKnight_Rotation_ArmyOfTheDead as ArmyOfTheDead,
	DeathKnight_Rotation as DeathKnightRotation,
	DeathKnight_Options as DeathKnightOptions,
} from '/wotlk/core/proto/deathknight.js';

import * as Presets from './presets.js';
import { SimUI } from '../core/sim_ui.js';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const StartingRunicPower = {
	type: 'number' as const,
	getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
	config: {
		extraCssClasses: [
			'starting-runic-power-picker',
		],
		label: 'Starting Runic Power',
		labelTooltip: 'Initial RP at the start of each iteration.',
		changedEvent: (player: Player<Spec.SpecDeathKnight>) => player.specOptionsChangeEmitter,
		getValue: (player: Player<Spec.SpecDeathKnight>) => player.getSpecOptions().startingRunicPower,
		setValue: (eventID: EventID, player: Player<Spec.SpecDeathKnight>, newValue: number) => {
			const newOptions = player.getSpecOptions();
			newOptions.startingRunicPower = newValue;
			player.setSpecOptions(eventID, newOptions);
		},
	},
};

export const PetUptime = {
	type: 'number' as const,
	getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
	config: {
		extraCssClasses: [
			'ghoul-uptime-picker',
		],
		label: 'Ghoul Uptime (%)',
		labelTooltip: 'Percent of the fight duration for which your ghoul will be on target.',
		changedEvent: (player: Player<Spec.SpecDeathKnight>) => player.specOptionsChangeEmitter,
		getValue: (player: Player<Spec.SpecDeathKnight>) => player.getSpecOptions().petUptime * 100,
		setValue: (eventID: EventID, player: Player<Spec.SpecDeathKnight>, newValue: number) => {
			const newOptions = player.getSpecOptions();
			newOptions.petUptime = newValue / 100;
			player.setSpecOptions(eventID, newOptions);
		},
	},
};

export const PrecastGhoulFrenzy = {
	type: 'boolean' as const,
	getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
	config: {
		extraCssClasses: [
			'precast-ghoul-frenzy-picker',
		],
		label: 'Pre-Cast Ghoul Frenzy',
		labelTooltip: 'Cast Ghoul Frenzy 10 seconds before combat starts.',
		changedEvent: (player: Player<Spec.SpecDeathKnight>) => player.specOptionsChangeEmitter,
		getValue: (player: Player<Spec.SpecDeathKnight>) => player.getSpecOptions().precastGhoulFrenzy,
		setValue: (eventID: EventID, player: Player<Spec.SpecDeathKnight>, newValue: boolean) => {
			const newOptions = player.getSpecOptions();
			newOptions.precastGhoulFrenzy = newValue;
			player.setSpecOptions(eventID, newOptions);
		},
	},
};

export const PrecastHornOfWinter = {
	type: 'boolean' as const,
	getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
	config: {
		extraCssClasses: [
			'precast-horn-of-winter-picker',
		],
		label: 'Pre-Cast Horn of Winter',
		labelTooltip: 'Precast Horn of Winter for 10 extra runic power before fight.',
		changedEvent: (player: Player<Spec.SpecDeathKnight>) => player.specOptionsChangeEmitter,
		getValue: (player: Player<Spec.SpecDeathKnight>) => player.getSpecOptions().precastHornOfWinter,
		setValue: (eventID: EventID, player: Player<Spec.SpecDeathKnight>, newValue: boolean) => {
			const newOptions = player.getSpecOptions();
			newOptions.precastHornOfWinter = newValue;
			player.setSpecOptions(eventID, newOptions);
		},
	},
};

export const RefreshHornOfWinter = {
	type: 'boolean' as const,
	getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
	config: {
		extraCssClasses: [
			'refresh-horn-of-winter-picker',
		],
		label: 'Refresh Horn of Winter',
		labelTooltip: 'Refresh Horn of Winter on free GCDs.',
		changedEvent: (player: Player<Spec.SpecDeathKnight>) => player.specOptionsChangeEmitter,
		getValue: (player: Player<Spec.SpecDeathKnight>) => player.getRotation().refreshHornOfWinter,
		setValue: (eventID: EventID, player: Player<Spec.SpecDeathKnight>, newValue: boolean) => {
			const newRotation = player.getRotation();
			newRotation.refreshHornOfWinter = newValue;
			player.setRotation(eventID, newRotation);
		},
	},
};

export const WIPFrostRotation = {
	type: 'boolean' as const,
	getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
	config: {
		extraCssClasses: [
			'wip-frost-rotation-picker',
		],
		label: 'Use WIP frost rotation',
		labelTooltip: 'Use sequence based rotation for frost, ***currently WIP***.',
		changedEvent: (player: Player<Spec.SpecDeathKnight>) => player.specOptionsChangeEmitter,
		getValue: (player: Player<Spec.SpecDeathKnight>) => player.getRotation().wipFrostRotation,
		setValue: (eventID: EventID, player: Player<Spec.SpecDeathKnight>, newValue: boolean) => {
			const newRotation = player.getRotation();
			newRotation.wipFrostRotation = newValue;
			player.setRotation(eventID, newRotation);
		},
	},
};

export const DiseaseRefreshDuration = {
	type: 'number' as const,
	getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
	config: {
		extraCssClasses: [
			'disease-refresh-duration-picker',
		],
		label: 'Disease Refresh Duration',
		labelTooltip: 'Minimum duration for refreshing a disease.',
		changedEvent: (player: Player<Spec.SpecDeathKnight>) => player.specOptionsChangeEmitter,
		getValue: (player: Player<Spec.SpecDeathKnight>) => player.getRotation().diseaseRefreshDuration,
		setValue: (eventID: EventID, player: Player<Spec.SpecDeathKnight>, newValue: number) => {
			const newRotation = player.getRotation();
			newRotation.diseaseRefreshDuration = newValue;
			player.setRotation(eventID, newRotation);
		},
	},
};

export const UseArmyOfTheDead = {
	type: 'enum' as const,
	getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
	config: {
		extraCssClasses: [
			'army-of-the-dead-enum-picker',
		],
		label: 'Army of the Dead',
		labelTooltip: 'Chose how to use Army of the Dead.',
		values: [
			{
				name: 'Do not use', value: ArmyOfTheDead.DoNotUse,
			},
			{
				name: 'Pre pull', value: ArmyOfTheDead.PreCast,
			},
			{
				name: 'As Major CD', value: ArmyOfTheDead.AsMajorCd,
			},
		],
		changedEvent: (player: Player<Spec.SpecDeathKnight>) => player.rotationChangeEmitter,
		getValue: (player: Player<Spec.SpecDeathKnight>) => player.getRotation().armyOfTheDead,
		setValue: (eventID: EventID, player: Player<Spec.SpecDeathKnight>, newValue: number) => {
			const newRotation = player.getRotation();
			newRotation.armyOfTheDead = newValue;
			player.setRotation(eventID, newRotation);
		},
	},
}

export const UseDeathAndDecay = {
	type: 'boolean' as const,
	getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
	config: {
		extraCssClasses: [
			'use-death-and-decay-picker',
		],
		label: 'Death and Decay',
		labelTooltip: 'Use Death and Decay based rotation.',
		changedEvent: (player: Player<Spec.SpecDeathKnight>) => player.specOptionsChangeEmitter,
		getValue: (player: Player<Spec.SpecDeathKnight>) => player.getRotation().useDeathAndDecay,
		setValue: (eventID: EventID, player: Player<Spec.SpecDeathKnight>, newValue: boolean) => {
			const newRotation = player.getRotation();
			newRotation.useDeathAndDecay = newValue;
			player.setRotation(eventID, newRotation);
		},
	},
};

export const UnholyPresenceOpener = {
	type: 'boolean' as const,
	getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
	config: {
		extraCssClasses: [
			'unholy-presence-opener-picker',
		],
		label: 'Unholy Presence Opener',
		labelTooltip: 'Start fight in unholy presence and change to blood after gargoyle.',
		changedEvent: (player: Player<Spec.SpecDeathKnight>) => player.specOptionsChangeEmitter,
		getValue: (player: Player<Spec.SpecDeathKnight>) => player.getRotation().unholyPresenceOpener,
		setValue: (eventID: EventID, player: Player<Spec.SpecDeathKnight>, newValue: boolean) => {
			const newRotation = player.getRotation();
			newRotation.unholyPresenceOpener = newValue;
			player.setRotation(eventID, newRotation);
		},
	},
};

export const DeathKnightRotationConfig = {
	inputs: [
		UseArmyOfTheDead,
		UseDeathAndDecay,
		UnholyPresenceOpener,
		RefreshHornOfWinter,
		WIPFrostRotation,
		DiseaseRefreshDuration,
	],
};
