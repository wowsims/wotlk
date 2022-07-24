import { Spec } from '/wotlk/core/proto/common.js';

import {
	DeathKnightTalents as DeathKnightTalents,
	DeathKnight_Rotation_ArmyOfTheDead as ArmyOfTheDead,
	DeathKnight_Rotation as DeathKnightRotation,
	DeathKnight_Options as DeathKnightOptions,
} from '/wotlk/core/proto/deathknight.js';

import * as InputHelpers from '/wotlk/core/components/input_helpers.js';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const StartingRunicPower = InputHelpers.makeSpecOptionsNumberInput<Spec.SpecDeathKnight>({
	fieldName: 'startingRunicPower',
	label: 'Starting Runic Power',
	labelTooltip: 'Initial RP at the start of each iteration.',
});

export const PetUptime = InputHelpers.makeSpecOptionsNumberInput<Spec.SpecDeathKnight>({
	fieldName: 'petUptime',
	label: 'Ghoul Uptime (%)',
	labelTooltip: 'Percent of the fight duration for which your ghoul will be on target.',
	percent: true,
});

export const PrecastGhoulFrenzy = InputHelpers.makeSpecOptionsBooleanInput<Spec.SpecDeathKnight>({
	fieldName: 'precastGhoulFrenzy',
	label: 'Pre-Cast Ghoul Frenzy',
	labelTooltip: 'Cast Ghoul Frenzy 10 seconds before combat starts.',
});

export const PrecastHornOfWinter = InputHelpers.makeSpecOptionsBooleanInput<Spec.SpecDeathKnight>({
	fieldName: 'precastHornOfWinter',
	label: 'Pre-Cast Horn of Winter',
	labelTooltip: 'Precast Horn of Winter for 10 extra runic power before fight.',
});

export const RefreshHornOfWinter = InputHelpers.makeRotationBooleanInput<Spec.SpecDeathKnight>({
	fieldName: 'refreshHornOfWinter',
	label: 'Refresh Horn of Winter',
	labelTooltip: 'Refresh Horn of Winter on free GCDs.',
});

export const DiseaseRefreshDuration = InputHelpers.makeRotationNumberInput<Spec.SpecDeathKnight>({
	fieldName: 'diseaseRefreshDuration',
	label: 'Disease Refresh Duration',
	labelTooltip: 'Minimum duration for refreshing a disease.',
});

export const UseDeathAndDecay = InputHelpers.makeRotationBooleanInput<Spec.SpecDeathKnight>({
	fieldName: 'useDeathAndDecay',
	label: 'Death and Decay',
	labelTooltip: 'Use Death and Decay based rotation.',
});

export const UseArmyOfTheDead = InputHelpers.makeRotationEnumInput<Spec.SpecDeathKnight, ArmyOfTheDead>({
	fieldName: 'armyOfTheDead',
	label: 'Army of the Dead',
	labelTooltip: 'Chose how to use Army of the Dead.',
	values: [
		{ name: 'Do not use', value: ArmyOfTheDead.DoNotUse },
		{ name: 'Pre pull', value: ArmyOfTheDead.PreCast },
		{ name: 'As Major CD', value: ArmyOfTheDead.AsMajorCd },
	],
});

export const UnholyPresenceOpener = InputHelpers.makeRotationBooleanInput<Spec.SpecDeathKnight>({
	fieldName: 'unholyPresenceOpener',
	label: 'Unholy Presence Opener',
	labelTooltip: 'Start fight in unholy presence and change to blood after gargoyle.',
});

export const DeathKnightRotationConfig = {
	inputs: [
		UseArmyOfTheDead,
		UseDeathAndDecay,
		UnholyPresenceOpener,
		RefreshHornOfWinter,
		DiseaseRefreshDuration,
	],
};
