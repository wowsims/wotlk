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

export const DeathKnightRotationConfig = {
	inputs: [
	],
};
