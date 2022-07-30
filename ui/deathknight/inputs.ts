import { Spec } from '/wotlk/core/proto/common.js';

import {
	DeathknightTalents as DeathKnightTalents,
	Deathknight_Rotation_ArmyOfTheDead as ArmyOfTheDead,
	Deathknight_Rotation_FirstDisease as FirstDisease,
	Deathknight_Rotation as DeathKnightRotation,
	Deathknight_Options as DeathKnightOptions,
} from '/wotlk/core/proto/deathknight.js';

import * as InputHelpers from '/wotlk/core/components/input_helpers.js';
import { Player } from '../core/player';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const StartingRunicPower = InputHelpers.makeSpecOptionsNumberInput<Spec.SpecDeathknight>({
	fieldName: 'startingRunicPower',
	label: 'Starting Runic Power',
	labelTooltip: 'Initial RP at the start of each iteration.',
});

export const PetUptime = InputHelpers.makeSpecOptionsNumberInput<Spec.SpecDeathknight>({
	fieldName: 'petUptime',
	label: 'Ghoul Uptime (%)',
	labelTooltip: 'Percent of the fight duration for which your ghoul will be on target.',
	percent: true,
});

export const PrecastGhoulFrenzy = InputHelpers.makeSpecOptionsBooleanInput<Spec.SpecDeathknight>({
	fieldName: 'precastGhoulFrenzy',
	label: 'Pre-Cast Ghoul Frenzy',
	labelTooltip: 'Cast Ghoul Frenzy 10 seconds before combat starts.',
});

export const PrecastHornOfWinter = InputHelpers.makeSpecOptionsBooleanInput<Spec.SpecDeathknight>({
	fieldName: 'precastHornOfWinter',
	label: 'Pre-Cast Horn of Winter',
	labelTooltip: 'Precast Horn of Winter for 10 extra runic power before fight.',
});

export const RefreshHornOfWinter = InputHelpers.makeRotationBooleanInput<Spec.SpecDeathknight>({
	fieldName: 'refreshHornOfWinter',
	label: 'Refresh Horn of Winter',
	labelTooltip: 'Refresh Horn of Winter on free GCDs.',
});

export const DiseaseRefreshDuration = InputHelpers.makeRotationNumberInput<Spec.SpecDeathknight>({
	fieldName: 'diseaseRefreshDuration',
	label: 'Disease Refresh Duration',
	labelTooltip: 'Minimum duration for refreshing a disease.',
});

export const UseDeathAndDecay = InputHelpers.makeRotationBooleanInput<Spec.SpecDeathknight>({
	fieldName: 'useDeathAndDecay',
	label: 'Death and Decay',
	labelTooltip: 'Use Death and Decay based rotation.',
	showWhen: (player: Player<Spec.SpecDeathknight>) => player.getTalents().summonGargoyle,
	changeEmitter: (player: Player<Spec.SpecDeathknight>) => player.talentsChangeEmitter,
});

export const SetFirstDisease = InputHelpers.makeRotationEnumInput<Spec.SpecDeathknight, FirstDisease>({
	fieldName: 'firstDisease',
	label: 'First Disease',
	labelTooltip: 'Chose which disease to apply first.',
	values: [
		{ name: 'Frost Fever', value: FirstDisease.FrostFever },
		{ name: 'Blood Plague', value: FirstDisease.BloodPlague },
	],
	showWhen: (player: Player<Spec.SpecDeathknight>) => player.getTalents().summonGargoyle,
	changeEmitter: (player: Player<Spec.SpecDeathknight>) => player.talentsChangeEmitter,
})

export const UseArmyOfTheDead = InputHelpers.makeRotationEnumInput<Spec.SpecDeathknight, ArmyOfTheDead>({
	fieldName: 'armyOfTheDead',
	label: 'Army of the Dead',
	labelTooltip: 'Chose how to use Army of the Dead.',
	values: [
		{ name: 'Do not use', value: ArmyOfTheDead.DoNotUse },
		{ name: 'Pre pull', value: ArmyOfTheDead.PreCast },
		{ name: 'As Major CD', value: ArmyOfTheDead.AsMajorCd },
	],
});

export const DeathKnightRotationConfig = {
	inputs: [
		SetFirstDisease,
		UseArmyOfTheDead,
		UseDeathAndDecay,
	],
};
