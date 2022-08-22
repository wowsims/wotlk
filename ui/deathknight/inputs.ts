import { Spec } from '../core/proto/common.js';

import {
	DeathknightTalents as DeathKnightTalents,
	Deathknight_Rotation_ArmyOfTheDead as ArmyOfTheDead,
	Deathknight_Rotation_FirstDisease as FirstDisease,
	Deathknight_Rotation_DeathAndDecayPrio as DeathAndDecayPrio,
	Deathknight_Rotation_StartingPresence as StartingPresence,
	Deathknight_Rotation_BloodRuneFiller as BloodRuneFiller,
	Deathknight_Rotation_BloodTap as BloodTap,
	Deathknight_Rotation as DeathKnightRotation,
	Deathknight_Options as DeathKnightOptions,
} from '../core/proto/deathknight.js';

import * as InputHelpers from '../core/components/input_helpers.js';
import { Player } from '../core/player';
import { TypedEvent } from '../core/typed_event';

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
	showWhen: (player: Player<Spec.SpecDeathknight>) => player.getTalents().summonGargoyle && player.getTalents().ghoulFrenzy,
	changeEmitter: (player: Player<Spec.SpecDeathknight>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
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
	showWhen: (player: Player<Spec.SpecDeathknight>) => player.getTalents().summonGargoyle && player.getTalents().scourgeStrike,
	changeEmitter: (player: Player<Spec.SpecDeathknight>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
});

export const SetDeathAndDecayPrio = InputHelpers.makeRotationEnumInput<Spec.SpecDeathknight, DeathAndDecayPrio>({
	fieldName: 'deathAndDecayPrio',
	label: 'Death and Decay Prio',
	labelTooltip: '<p>Chose how to prioritize death and decay usage:</p>\
		<p><b>Max Rune Downtime</b>: Prioritizes spending runes over holding them for death and decay</p>\
		<p><b>Max Dnd Uptime</b>: Prioritizes dnd uptime and can hold runes for longer then rune grace</p>',
	values: [
		{ name: 'Max Rune Downtime', value: DeathAndDecayPrio.MaxRuneDowntime },
		{ name: 'Max Dnd Uptime', value: DeathAndDecayPrio.MaxDndUptime },
	],
	showWhen: (player: Player<Spec.SpecDeathknight>) => player.getTalents().summonGargoyle && (player.getRotation().useDeathAndDecay || !player.getTalents().scourgeStrike),
	changeEmitter: (player: Player<Spec.SpecDeathknight>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
})

export const UseEmpowerRuneWeapon = InputHelpers.makeRotationBooleanInput<Spec.SpecDeathknight>({
	fieldName: 'useEmpowerRuneWeapon',
	label: 'Empower Rune Weapon',
	labelTooltip: 'Use Empower Rune Weapon in rotation.',
});

export const BloodTapGhoulFrenzy = InputHelpers.makeRotationBooleanInput<Spec.SpecDeathknight>({
	fieldName: 'btGhoulFrenzy',
	label: 'BT Ghoul Frenzy',
	labelTooltip: 'Use Ghoul Frenzy only with Blood Tap.',
	showWhen: (player: Player<Spec.SpecDeathknight>) => player.getTalents().ghoulFrenzy,
	changeEmitter: (player: Player<Spec.SpecDeathknight>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
});

export const FirstDiseaseInput = InputHelpers.makeRotationEnumInput<Spec.SpecDeathknight, FirstDisease>({
	fieldName: 'firstDisease',
	label: 'First Disease',
	labelTooltip: 'Chose which disease to apply first.',
	values: [
		{ name: 'Frost Fever', value: FirstDisease.FrostFever },
		{ name: 'Blood Plague', value: FirstDisease.BloodPlague },
	],
	showWhen: (player: Player<Spec.SpecDeathknight>) => player.getTalents().summonGargoyle,
	changeEmitter: (player: Player<Spec.SpecDeathknight>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
})

export const ArmyOfTheDeadInput = InputHelpers.makeRotationEnumInput<Spec.SpecDeathknight, ArmyOfTheDead>({
	fieldName: 'armyOfTheDead',
	label: 'Army of the Dead',
	labelTooltip: 'Chose how to use Army of the Dead.',
	values: [
		{ name: 'Do not use', value: ArmyOfTheDead.DoNotUse },
		{ name: 'Pre pull', value: ArmyOfTheDead.PreCast },
		{ name: 'As Major CD', value: ArmyOfTheDead.AsMajorCd },
	],
});

export const StartingPresenceInput = InputHelpers.makeRotationEnumInput<Spec.SpecDeathknight, StartingPresence>({
	fieldName: 'startingPresence',
	label: 'Starting Presence',
	labelTooltip: 'Chose the presence you start combat in.',
	values: [
		{ name: 'Blood', value: StartingPresence.Blood },
		{ name: 'Unholy', value: StartingPresence.Unholy },
	],
	showWhen: (player: Player<Spec.SpecDeathknight>) => player.getTalents().summonGargoyle,
	changeEmitter: (player: Player<Spec.SpecDeathknight>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
})

export const BloodRuneFillerInput = InputHelpers.makeRotationEnumInput<Spec.SpecDeathknight, BloodRuneFiller>({
	fieldName: 'bloodRuneFiller',
	label: 'Blood Rune Filler',
	labelTooltip: 'Chose what to spend your free blood runes on.',
	values: [
		{ name: 'Blood Strike', value: BloodRuneFiller.BloodStrike },
		{ name: 'Blood Boil', value: BloodRuneFiller.BloodBoil },
	],
	showWhen: (player: Player<Spec.SpecDeathknight>) => player.getTalents().summonGargoyle,
	changeEmitter: (player: Player<Spec.SpecDeathknight>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
})

export const BloodTapInput = InputHelpers.makeRotationEnumInput<Spec.SpecDeathknight, BloodTap>({
	fieldName: 'bloodTap',
	label: 'Blood Tap',
	labelTooltip: 'Chose what to spend your Blood Tap on.',
	values: [
		{ name: 'Ghoul Frenzy', value: BloodTap.GhoulFrenzy },
		{ name: 'Icy Touch', value: BloodTap.IcyTouch },
		{ name: 'Blood Strike', value: BloodTap.BloodStrikeBT },
		{ name: 'Blood Boil', value: BloodTap.BloodBoilBT },
	],
	showWhen: (player: Player<Spec.SpecDeathknight>) => player.getTalents().summonGargoyle,
	changeEmitter: (player: Player<Spec.SpecDeathknight>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
})

export const UseAMSInput = InputHelpers.makeRotationBooleanInput<Spec.SpecDeathknight>({
	fieldName: 'useAms',
	label: 'Use AMS',
	labelTooltip: 'Use AMS around predicted damage for a RP gain.',
	changeEmitter: (player: Player<Spec.SpecDeathknight>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
});

export const AvgAMSSuccessRateInput = InputHelpers.makeRotationNumberInput<Spec.SpecDeathknight>({
	fieldName: 'avgAmsSuccessRate',
	label: 'Avg AMS Success %',
	labelTooltip: 'Chance for damage to be taken during the 5 second window of AMS.',
	showWhen: (player: Player<Spec.SpecDeathknight>) => player.getRotation().useAms == true,
	changeEmitter: (player: Player<Spec.SpecDeathknight>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
});

export const AvgAMSHitInput = InputHelpers.makeRotationNumberInput<Spec.SpecDeathknight>({
	fieldName: 'avgAmsHit',
	label: 'Avg AMS Hit',
	labelTooltip: 'How much on average (+-10%) the character is hit for when AMS is successful.',
	showWhen: (player: Player<Spec.SpecDeathknight>) => player.getRotation().useAms == true,
	changeEmitter: (player: Player<Spec.SpecDeathknight>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
});

export const OblitDelayDurationInput = InputHelpers.makeRotationNumberInput<Spec.SpecDeathknight>({
	fieldName: 'oblitDelayDuration',
	label: 'Oblit Delay (ms)',
	labelTooltip: 'How long a FS/HB/HW can delay a Oblit by.',
	showWhen: (player: Player<Spec.SpecDeathknight>) => player.getTalents().howlingBlast,
	changeEmitter: (player: Player<Spec.SpecDeathknight>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
});

export const DeathKnightRotationConfig = {
	inputs: [
		BloodTapGhoulFrenzy,
		UseEmpowerRuneWeapon,
		BloodTapInput,
		ArmyOfTheDeadInput,
		FirstDiseaseInput,
		StartingPresenceInput,
		BloodRuneFillerInput,
		UseDeathAndDecay,
		OblitDelayDurationInput,
		UseAMSInput,
		AvgAMSSuccessRateInput,
		AvgAMSHitInput,
		//SetDeathAndDecayPrio,
	],
};
