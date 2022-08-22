import { Spec } from '../core/proto/common.js';


import {
	DeathknightTalents as DeathknightTalents,
	Deathknight_Rotation as DeathknightRotation,
	Deathknight_Options as DeathknightOptions,
} from '../core/proto/deathknight.js';


import * as InputHelpers from '../core/components/input_helpers.js';
import { Player } from '../core/player';
import { TypedEvent } from '../core/typed_event';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const StartingRunicPower = InputHelpers.makeSpecOptionsNumberInput<Spec.SpecTankDeathknight>({
	fieldName: 'startingRunicPower',
	label: 'Starting Runic Power',
	labelTooltip: 'Initial RP at the start of each iteration.',
});

export const TankDeathKnightRotationConfig = {
	inputs: [
	],
};
