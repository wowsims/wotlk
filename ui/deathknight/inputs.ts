import { RaidTarget, Spec } from '../core/proto/common.js';
import { ActionId } from '../core/proto_utils/action_id.js';

import {
	DeathknightTalents as DeathKnightTalents,
	Deathknight_Rotation_ArmyOfTheDead as ArmyOfTheDead,
	Deathknight_Rotation_FirstDisease as FirstDisease,
	Deathknight_Rotation_DeathAndDecayPrio as DeathAndDecayPrio,
	Deathknight_Rotation_Presence as StartingPresence,
	Deathknight_Rotation_BloodRuneFiller as BloodRuneFiller,
	Deathknight_Rotation_BloodTap as BloodTap,
	Deathknight_Rotation_FrostRotationType as FrostRotationType,
	Deathknight_Rotation_CustomSpellOption as CustomSpellOption,
	Deathknight_Rotation as DeathKnightRotation,
	Deathknight_Options as DeathKnightOptions,
} from '../core/proto/deathknight.js';

import * as InputHelpers from '../core/components/input_helpers.js';
import { Player } from '../core/player';
import { EventID, TypedEvent } from '../core/typed_event';
import { NO_TARGET } from '../core/proto_utils/utils.js';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const SelfUnholyFrenzy = InputHelpers.makeSpecOptionsBooleanInput<Spec.SpecDeathknight>({
	fieldName: 'unholyFrenzyTarget',
	label: 'Self Unholy Frenzy',
	labelTooltip: 'Cast Unholy Frenzy on yourself.',
	extraCssClasses: [
		'within-raid-sim-hide',
	],
	getValue: (player: Player<Spec.SpecDeathknight>) => player.getSpecOptions().unholyFrenzyTarget?.targetIndex != NO_TARGET,
	setValue: (eventID: EventID, player: Player<Spec.SpecDeathknight>, newValue: boolean) => {
		const newOptions = player.getSpecOptions();
		newOptions.unholyFrenzyTarget = RaidTarget.create({
			targetIndex: newValue ? 0 : NO_TARGET,
		});
		player.setSpecOptions(eventID, newOptions);
	},
	showWhen: (player: Player<Spec.SpecDeathknight>) => player.getTalents().hysteria,
	changeEmitter: (player: Player<Spec.SpecDeathknight>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
});

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
	showWhen: (player: Player<Spec.SpecDeathknight>) => player.getTalents().masterOfGhouls,
});

export const PrecastGhoulFrenzy = InputHelpers.makeSpecOptionsBooleanInput<Spec.SpecDeathknight>({
	fieldName: 'precastGhoulFrenzy',
	label: 'Pre-Cast Ghoul Frenzy',
	labelTooltip: 'Cast Ghoul Frenzy 10 seconds before combat starts.',
	showWhen: (player: Player<Spec.SpecDeathknight>) => player.getTalents().summonGargoyle && player.getTalents().ghoulFrenzy,
	changeEmitter: (player: Player<Spec.SpecDeathknight>) => TypedEvent.onAny([player.specOptionsChangeEmitter, player.rotationChangeEmitter, player.talentsChangeEmitter]),
});

export const PrecastHornOfWinter = InputHelpers.makeSpecOptionsBooleanInput<Spec.SpecDeathknight>({
	fieldName: 'precastHornOfWinter',
	label: 'Pre-Cast Horn of Winter',
	labelTooltip: 'Precast Horn of Winter for 10 extra runic power before fight.',
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
	showWhen: (player: Player<Spec.SpecDeathknight>) => player.getTalents().summonGargoyle && player.getTalents().scourgeStrike && !player.getRotation().autoRotation,
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
	showWhen: (player: Player<Spec.SpecDeathknight>) => player.getTalents().summonGargoyle && (player.getRotation().useDeathAndDecay || !player.getTalents().scourgeStrike) && !player.getRotation().autoRotation,
	changeEmitter: (player: Player<Spec.SpecDeathknight>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
})

export const UseEmpowerRuneWeapon = InputHelpers.makeRotationBooleanInput<Spec.SpecDeathknight>({
	fieldName: 'useEmpowerRuneWeapon',
	label: 'Empower Rune Weapon',
	labelTooltip: 'Use Empower Rune Weapon in rotation.',
	showWhen: (player: Player<Spec.SpecDeathknight>) => !player.getRotation().autoRotation && player.getRotation().frostRotationType != FrostRotationType.Custom,
});

export const UseGargoyle = InputHelpers.makeRotationBooleanInput<Spec.SpecDeathknight>({
	fieldName: 'useGargoyle',
	label: 'Summon Gargoyle',
	labelTooltip: 'Use Summon Gargoyle in rotation.',
	changeEmitter: (player: Player<Spec.SpecDeathknight>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
	showWhen: (player: Player<Spec.SpecDeathknight>) => player.getTalents().summonGargoyle && !player.getRotation().autoRotation,
});

export const HoldErwArmy = InputHelpers.makeRotationBooleanInput<Spec.SpecDeathknight>({
	fieldName: 'holdErwArmy',
	label: 'Hold ERW for AotD',
	labelTooltip: 'Hold Empower Rune Weapon for after Summon Gargoyle to guarantee maximized snapshot for Army of the Dead.',
	changeEmitter: (player: Player<Spec.SpecDeathknight>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
	showWhen: (player: Player<Spec.SpecDeathknight>) => !player.getRotation().autoRotation && player.getRotation().useEmpowerRuneWeapon && player.getRotation().armyOfTheDead == ArmyOfTheDead.AsMajorCd,
});

export const BloodlustPresence = InputHelpers.makeRotationEnumInput<Spec.SpecDeathknight, StartingPresence>({
	fieldName: 'blPresence',
	label: 'Bloodlust Presence',
	labelTooltip: 'Presence during bloodlust.',
	values: [
		{ name: 'Blood', value: StartingPresence.Blood },
		{ name: 'Unholy', value: StartingPresence.Unholy },
	],
	showWhen: (player: Player<Spec.SpecDeathknight>) => player.getTalents().summonGargoyle && !player.getRotation().autoRotation,
	changeEmitter: (player: Player<Spec.SpecDeathknight>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
});

export const GargoylePresence = InputHelpers.makeRotationEnumInput<Spec.SpecDeathknight, StartingPresence>({
	fieldName: 'gargoylePresence',
	label: 'Gargoyle Presence',
	labelTooltip: 'Presence during Gargoyle.',
	values: [
		{ name: 'Blood', value: StartingPresence.Blood },
		{ name: 'Unholy', value: StartingPresence.Unholy },
	],
	showWhen: (player: Player<Spec.SpecDeathknight>) => player.getTalents().summonGargoyle && !player.getRotation().autoRotation,
	changeEmitter: (player: Player<Spec.SpecDeathknight>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
});

export const BloodTapGhoulFrenzy = InputHelpers.makeRotationBooleanInput<Spec.SpecDeathknight>({
	fieldName: 'btGhoulFrenzy',
	label: 'BT Ghoul Frenzy',
	labelTooltip: 'Use Ghoul Frenzy only with Blood Tap.',
	showWhen: (player: Player<Spec.SpecDeathknight>) => player.getTalents().ghoulFrenzy && !player.getRotation().autoRotation,
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
	showWhen: (player: Player<Spec.SpecDeathknight>) => player.getTalents().summonGargoyle && !player.getRotation().autoRotation,
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
	showWhen: (player: Player<Spec.SpecDeathknight>) => !player.getRotation().autoRotation,
	changeEmitter: (player: Player<Spec.SpecDeathknight>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
});

export const StartingPresenceInput = InputHelpers.makeRotationEnumInput<Spec.SpecDeathknight, StartingPresence>({
	fieldName: 'startingPresence',
	label: 'Starting Presence',
	labelTooltip: 'Chose the presence you start combat in.',
	values: [
		{ name: 'Blood', value: StartingPresence.Blood },
		{ name: 'Unholy', value: StartingPresence.Unholy },
	],
	showWhen: (player: Player<Spec.SpecDeathknight>) => player.getTalents().summonGargoyle && !player.getRotation().autoRotation,
	changeEmitter: (player: Player<Spec.SpecDeathknight>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
})

export const FightPresence = InputHelpers.makeRotationEnumInput<Spec.SpecDeathknight, StartingPresence>({
	fieldName: 'presence',
	label: 'Fight Presence',
	labelTooltip: 'Presence to be in during the encounter.',
	values: [
		{ name: 'Blood', value: StartingPresence.Blood },
		{ name: 'Unholy', value: StartingPresence.Unholy },
	],
	showWhen: (player: Player<Spec.SpecDeathknight>) => player.getTalents().summonGargoyle && !player.getRotation().autoRotation,
	changeEmitter: (player: Player<Spec.SpecDeathknight>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
});

export const BloodRuneFillerInput = InputHelpers.makeRotationEnumInput<Spec.SpecDeathknight, BloodRuneFiller>({
	fieldName: 'bloodRuneFiller',
	label: 'Blood Rune Filler',
	labelTooltip: 'Chose what to spend your free blood runes on.',
	values: [
		{ name: 'Blood Strike', value: BloodRuneFiller.BloodStrike },
		{ name: 'Blood Boil', value: BloodRuneFiller.BloodBoil },
	],
	showWhen: (player: Player<Spec.SpecDeathknight>) => player.getTalents().summonGargoyle && !player.getRotation().autoRotation,
	changeEmitter: (player: Player<Spec.SpecDeathknight>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
})

export const NerfedGargoyleInput = InputHelpers.makeRotationBooleanInput<Spec.SpecDeathknight>({
	fieldName: 'nerfedGargoyle',
	label: 'Nerfed Gargoyle (no haste snapshot)',
	labelTooltip: "Use updated PTR Gargoyle that doesn't snapshot haste",
	showWhen: (player: Player<Spec.SpecDeathknight>) => player.getTalents().summonGargoyle && !player.getRotation().autoRotation,
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
	showWhen: (player: Player<Spec.SpecDeathknight>) => player.getTalents().summonGargoyle && !player.getRotation().autoRotation,
	changeEmitter: (player: Player<Spec.SpecDeathknight>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
})

export const UseAMSInput = InputHelpers.makeRotationBooleanInput<Spec.SpecDeathknight>({
	fieldName: 'useAms',
	label: 'Use AMS',
	labelTooltip: 'Use AMS around predicted damage for a RP gain.',
	showWhen: (player: Player<Spec.SpecDeathknight>) => player.getTalents().howlingBlast && !player.getRotation().autoRotation && player.getRotation().frostRotationType != FrostRotationType.Custom,
	changeEmitter: (player: Player<Spec.SpecDeathknight>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
});

export const AvgAMSSuccessRateInput = InputHelpers.makeRotationNumberInput<Spec.SpecDeathknight>({
	fieldName: 'avgAmsSuccessRate',
	label: 'Avg AMS Success %',
	labelTooltip: 'Chance for damage to be taken during the 5 second window of AMS.',
	showWhen: (player: Player<Spec.SpecDeathknight>) => player.getRotation().useAms == true && !player.getRotation().autoRotation && player.getTalents().howlingBlast && player.getRotation().frostRotationType != FrostRotationType.Custom,
	changeEmitter: (player: Player<Spec.SpecDeathknight>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
});

export const AvgAMSHitInput = InputHelpers.makeRotationNumberInput<Spec.SpecDeathknight>({
	fieldName: 'avgAmsHit',
	label: 'Avg AMS Hit',
	labelTooltip: 'How much on average (+-10%) the character is hit for when AMS is successful.',
	showWhen: (player: Player<Spec.SpecDeathknight>) => player.getRotation().useAms == true && !player.getRotation().autoRotation && player.getTalents().howlingBlast && player.getRotation().frostRotationType != FrostRotationType.Custom,
	changeEmitter: (player: Player<Spec.SpecDeathknight>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
});

export const OblitDelayDurationInput = InputHelpers.makeRotationNumberInput<Spec.SpecDeathknight>({
	fieldName: 'oblitDelayDuration',
	label: 'Oblit Delay (ms)',
	labelTooltip: 'How long a FS/HB/HW can delay a Oblit by.',
	showWhen: (player: Player<Spec.SpecDeathknight>) => player.getTalents().howlingBlast && !player.getRotation().autoRotation && player.getRotation().frostRotationType != FrostRotationType.Custom,
	changeEmitter: (player: Player<Spec.SpecDeathknight>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
});

export const UseAutoRotation = InputHelpers.makeRotationBooleanInput<Spec.SpecDeathknight>({
	fieldName: 'autoRotation',
	label: 'Automatic Rotation',
	labelTooltip: 'Have sim automatically adjust rotation based on the scenario. This is still in development and currently only works for Unholy.',
	changeEmitter: (player: Player<Spec.SpecDeathknight>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
showWhen: (player: Player<Spec.SpecDeathknight>) => !player.getTalents().howlingBlast,
});

export const DesyncRotation = InputHelpers.makeRotationBooleanInput<Spec.SpecDeathknight>({
	fieldName: 'desyncRotation',
	label: 'Use Desync Rotation',
	labelTooltip: 'Use the Desync Rotation.',
	showWhen: (player: Player<Spec.SpecDeathknight>) => player.getTalents().howlingBlast && !player.getRotation().autoRotation && player.sim.getShowExperimental(),
	changeEmitter: (player: Player<Spec.SpecDeathknight>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
});

export const Presence = InputHelpers.makeRotationEnumInput<Spec.SpecDeathknight, StartingPresence>({
	fieldName: 'presence',
	label: 'Presence',
	labelTooltip: 'Presence to be in during the encounter.',
	values: [
		{ name: 'Blood', value: StartingPresence.Blood },
		{ name: 'Frost', value: StartingPresence.Frost },
		{ name: 'Unholy', value: StartingPresence.Unholy },
	],
	showWhen: (player: Player<Spec.SpecDeathknight>) => player.getTalents().howlingBlast && !player.getRotation().autoRotation,
	changeEmitter: (player: Player<Spec.SpecDeathknight>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
});

export const FrostCustomRotation = InputHelpers.makeCustomRotationInput<Spec.SpecDeathknight, CustomSpellOption>({
	fieldName: 'frostCustomRotation',
	numColumns: 4,
	values: [
		{ actionId: ActionId.fromSpellId(49909), value: CustomSpellOption.CustomIcyTouch },
		{ actionId: ActionId.fromSpellId(49921), value: CustomSpellOption.CustomPlagueStrike },
		{ actionId: ActionId.fromSpellId(50842), value: CustomSpellOption.CustomPestilence },
		{ actionId: ActionId.fromSpellId(51425), value: CustomSpellOption.CustomObliterate },
		{ actionId: ActionId.fromSpellId(51411), value: CustomSpellOption.CustomHowlingBlast },
		{ actionId: ActionId.fromSpellId(59057), value: CustomSpellOption.CustomHowlingBlastRime },
		{ actionId: ActionId.fromSpellId(49941), value: CustomSpellOption.CustomBloodBoil },
		{ actionId: ActionId.fromSpellId(49930), value: CustomSpellOption.CustomBloodStrike },
		{ actionId: ActionId.fromSpellId(49938), value: CustomSpellOption.CustomDeathAndDecay },
		{ actionId: ActionId.fromSpellId(57623), value: CustomSpellOption.CustomHornOfWinter },
		{ actionId: ActionId.fromSpellId(51271), value: CustomSpellOption.CustomUnbreakableArmor },
		{ actionId: ActionId.fromSpellId(45529), value: CustomSpellOption.CustomBloodTap },
		{ actionId: ActionId.fromSpellId(47568), value: CustomSpellOption.CustomEmpoweredRuneWeapon },
		{ actionId: ActionId.fromSpellId(55268), value: CustomSpellOption.CustomFrostStrike },
	],
	showWhen: (player: Player<Spec.SpecDeathknight>) => player.getRotation().frostRotationType == FrostRotationType.Custom,
});

export const DeathKnightRotationConfig = {
	inputs: [
		InputHelpers.makeRotationEnumInput<Spec.SpecDeathknight, FrostRotationType>({
			fieldName: 'frostRotationType',
			label: 'Rotation Type',
			values: [
				{ name: 'Single Target', value: FrostRotationType.SingleTarget },
				{ name: 'Custom', value: FrostRotationType.Custom },
			],
			changeEmitter: (player: Player<Spec.SpecDeathknight>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
			showWhen: (player: Player<Spec.SpecDeathknight>) => player.getTalents().howlingBlast && !player.getRotation().autoRotation,
		}),
		Presence,
		UseAutoRotation,
		BloodTapGhoulFrenzy,
		UseGargoyle,
		UseEmpowerRuneWeapon,
		HoldErwArmy,
		BloodTapInput,
		ArmyOfTheDeadInput,
		FirstDiseaseInput,
		StartingPresenceInput,
		BloodlustPresence,
		GargoylePresence,
		FightPresence,
		BloodRuneFillerInput,
		UseDeathAndDecay,
		OblitDelayDurationInput,
		UseAMSInput,
		AvgAMSSuccessRateInput,
		AvgAMSHitInput,
		DesyncRotation,
		FrostCustomRotation,
		NerfedGargoyleInput,
	],
};
