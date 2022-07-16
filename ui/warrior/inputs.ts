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
	WarriorShout,
	WarriorTalents as WarriorTalents,
	Warrior,
	Warrior_Rotation as WarriorRotation,
	Warrior_Rotation_SunderArmor as SunderArmor,
	Warrior_Options as WarriorOptions,
} from '/wotlk/core/proto/warrior.js';

import * as Presets from './presets.js';
import { SimUI } from '../core/sim_ui.js';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const Recklessness = {
	id: ActionId.fromSpellId(1719),
	states: 2,
	extraCssClasses: [
		'warrior-Recklessness-picker',
	],
	changedEvent: (player: Player<Spec.SpecWarrior>) => player.specOptionsChangeEmitter,
	getValue: (player: Player<Spec.SpecWarrior>) => player.getSpecOptions().useRecklessness,
	setValue: (eventID: EventID, player: Player<Spec.SpecWarrior>, newValue: boolean) => {
		const newOptions = player.getSpecOptions();
		newOptions.useRecklessness = newValue
		player.setSpecOptions(eventID, newOptions);
	},
};

export const StartingRage = {
	type: 'number' as const,
	getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
	config: {
		extraCssClasses: [
			'starting-rage-picker',
		],
		label: 'Starting Rage',
		labelTooltip: 'Initial rage at the start of each iteration.',
		changedEvent: (player: Player<Spec.SpecWarrior>) => player.specOptionsChangeEmitter,
		getValue: (player: Player<Spec.SpecWarrior>) => player.getSpecOptions().startingRage,
		setValue: (eventID: EventID, player: Player<Spec.SpecWarrior>, newValue: number) => {
			const newOptions = player.getSpecOptions();
			newOptions.startingRage = newValue;
			player.setSpecOptions(eventID, newOptions);
		},
	},
};

export const ShoutPicker = {
	extraCssClasses: [
		'shout-picker',
	],
	numColumns: 1,
	values: [
		{ color: 'c79c6e', value: WarriorShout.WarriorShoutNone },
		{ actionId: ActionId.fromSpellId(2048), value: WarriorShout.WarriorShoutBattle },
		{ actionId: ActionId.fromSpellId(469), value: WarriorShout.WarriorShoutCommanding },
	],
	equals: (a: WarriorShout, b: WarriorShout) => a == b,
	zeroValue: WarriorShout.WarriorShoutNone,
	changedEvent: (player: Player<Spec.SpecWarrior>) => player.specOptionsChangeEmitter,
	getValue: (player: Player<Spec.SpecWarrior>) => player.getSpecOptions().shout,
	setValue: (eventID: EventID, player: Player<Spec.SpecWarrior>, newValue: number) => {
		const newOptions = player.getSpecOptions();
		newOptions.shout = newValue;
		player.setSpecOptions(eventID, newOptions);
	},
};

export const WarriorRotationConfig = {
	inputs: [
		{
			type: 'boolean' as const,
			cssClass: 'cleave-picker',
			getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
			config: {
				label: 'Use Cleave',
				labelTooltip: 'Use Cleave instead of Heroic Strike.',
				changedEvent: (player: Player<Spec.SpecWarrior>) => player.rotationChangeEmitter,
				getValue: (player: Player<Spec.SpecWarrior>) => player.getRotation().useCleave,
				setValue: (eventID: EventID, player: Player<Spec.SpecWarrior>, newValue: boolean) => {
					const newRotation = player.getRotation();
					newRotation.useCleave = newValue;
					player.setRotation(eventID, newRotation);
				},
			},
		},
		{
			type: 'boolean' as const,
			cssClass: 'rend-picker',
			getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
			config: {
				label: 'Use Rend',
				labelTooltip: 'Use Rend on free globals.',
				changedEvent: (player: Player<Spec.SpecWarrior>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
				getValue: (player: Player<Spec.SpecWarrior>) => player.getRotation().useRend,
				setValue: (eventID: EventID, player: Player<Spec.SpecWarrior>, newValue: boolean) => {
					const newRotation = player.getRotation();
					newRotation.useRend = newValue;
					player.setRotation(eventID, newRotation);
				},
			},
		},
		{
			type: 'boolean' as const,
			cssClass: 'ms-picker',
			getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
			config: {
				label: 'Use Mortal Strike',
				labelTooltip: 'Use Mortal Strike when rage threshold is met.',
				changedEvent: (player: Player<Spec.SpecWarrior>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
				getValue: (player: Player<Spec.SpecWarrior>) => player.getRotation().useMs,
				setValue: (eventID: EventID, player: Player<Spec.SpecWarrior>, newValue: boolean) => {
					const newRotation = player.getRotation();
					newRotation.useMs = newValue;
					player.setRotation(eventID, newRotation);
				},
				showWhen: (player: Player<Spec.SpecWarrior>) => player.getTalents().mortalStrike,
			},
		},
		{
			type: 'boolean' as const,
			cssClass: 'prioritize-ww-picker',
			getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
			config: {
				label: 'Prioritize WW',
				labelTooltip: 'Prioritize Whirlwind over Bloodthirst',
				changedEvent: (player: Player<Spec.SpecWarrior>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
				getValue: (player: Player<Spec.SpecWarrior>) => player.getRotation().prioritizeWw,
				setValue: (eventID: EventID, player: Player<Spec.SpecWarrior>, newValue: boolean) => {
					const newRotation = player.getRotation();
					newRotation.prioritizeWw = newValue;
					player.setRotation(eventID, newRotation);
				},
				showWhen: (player: Player<Spec.SpecWarrior>) => player.getTalents().bloodthirst,
			},
		},
		{
			type: 'number' as const,
			cssClass: 'hs-rage-threshold',
			getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
			config: {
				label: 'HS rage threshold',
				labelTooltip: 'Heroic Strike when rage is above:',
				changedEvent: (player: Player<Spec.SpecWarrior>) => player.rotationChangeEmitter,
				getValue: (player: Player<Spec.SpecWarrior>) => player.getRotation().hsRageThreshold,
				setValue: (eventID: EventID, player: Player<Spec.SpecWarrior>, newValue: number) => {
					const newRotation = player.getRotation();
					newRotation.hsRageThreshold = newValue;
					player.setRotation(eventID, newRotation);
				},
			},
		},
		{
			type: 'number' as const,
			cssClass: 'ms-rage-threshold',
			getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
			config: {
				label: 'Mortal Strike rage threshold',
				labelTooltip: 'Use Mortal Strike when rage is below a point.',
				changedEvent: (player: Player<Spec.SpecWarrior>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
				getValue: (player: Player<Spec.SpecWarrior>) => player.getRotation().msRageThreshold,
				setValue: (eventID: EventID, player: Player<Spec.SpecWarrior>, newValue: number) => {
					const newRotation = player.getRotation();
					newRotation.msRageThreshold = newValue;
					player.setRotation(eventID, newRotation);
				},
				showWhen: (player: Player<Spec.SpecWarrior>) => player.getRotation().useMs,
			},
		},
		{
			type: 'number' as const,
			cssClass: 'rend-rage-threshold',
			getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
			config: {
				label: 'Rend rage threshold',
				labelTooltip: 'Rend will only be used when rage is larger than this value.',
				changedEvent: (player: Player<Spec.SpecWarrior>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
				getValue: (player: Player<Spec.SpecWarrior>) => player.getRotation().rendRageThreshold,
				setValue: (eventID: EventID, player: Player<Spec.SpecWarrior>, newValue: number) => {
					const newRotation = player.getRotation();
					newRotation.rendRageThreshold = newValue;
					player.setRotation(eventID, newRotation);
				},
				showWhen: (player: Player<Spec.SpecWarrior>) => player.getRotation().useRend && player.getTalents().bloodthirst,
			},
		},
		{
			type: 'number' as const,
			cssClass: 'rend-duration-threshold',
			getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
			config: {
				label: 'Rend Refresh Time',
				labelTooltip: 'Refresh Rend when the remaining duration is less than this amount of time (seconds).',
				changedEvent: (player: Player<Spec.SpecWarrior>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
				getValue: (player: Player<Spec.SpecWarrior>) => player.getRotation().rendCdThreshold,
				setValue: (eventID: EventID, player: Player<Spec.SpecWarrior>, newValue: number) => {
					const newRotation = player.getRotation();
					newRotation.rendCdThreshold = newValue;
					player.setRotation(eventID, newRotation);
				},
				showWhen: (player: Player<Spec.SpecWarrior>) => player.getRotation().useRend && player.getTalents().mortalStrike,
			},
		},
		{
			type: 'boolean' as const,
			cssClass: 'hs-exec-picker',
			getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
			config: {
				label: 'HS during Execute Phase',
				labelTooltip: 'Use Heroic Strike during Execute Phase.',
				changedEvent: (player: Player<Spec.SpecWarrior>) => player.rotationChangeEmitter,
				getValue: (player: Player<Spec.SpecWarrior>) => player.getRotation().useHsDuringExecute,
				setValue: (eventID: EventID, player: Player<Spec.SpecWarrior>, newValue: boolean) => {
					const newRotation = player.getRotation();
					newRotation.useHsDuringExecute = newValue;
					player.setRotation(eventID, newRotation);
				},
			},
		},
		{
			type: 'boolean' as const,
			cssClass: 'bt-exec-picker-fury',
			getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
			config: {
				label: 'BT during Execute Phase',
				labelTooltip: 'Use Bloodthirst during Execute Phase.',
				changedEvent: (player: Player<Spec.SpecWarrior>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
				getValue: (player: Player<Spec.SpecWarrior>) => player.getRotation().useBtDuringExecute,
				setValue: (eventID: EventID, player: Player<Spec.SpecWarrior>, newValue: boolean) => {
					const newRotation = player.getRotation();
					newRotation.useBtDuringExecute = newValue;
					player.setRotation(eventID, newRotation);
				},
				showWhen: (player: Player<Spec.SpecWarrior>) => player.getTalents().bloodthirst,
			},
		},
		{
			type: 'boolean' as const,
			cssClass: 'ww-exec-picker',
			getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
			config: {
				label: 'WW during Execute Phase',
				labelTooltip: 'Use Whirlwind during Execute Phase.',
				changedEvent: (player: Player<Spec.SpecWarrior>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
				getValue: (player: Player<Spec.SpecWarrior>) => player.getRotation().useWwDuringExecute,
				setValue: (eventID: EventID, player: Player<Spec.SpecWarrior>, newValue: boolean) => {
					const newRotation = player.getRotation();
					newRotation.useWwDuringExecute = newValue;
					player.setRotation(eventID, newRotation);
				},
				showWhen: (player: Player<Spec.SpecWarrior>) => player.getTalents().bloodthirst,
			},
		},
		{
			type: 'boolean' as const,
			cssClass: 'spam-exec-picker',
			getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
			config: {
				label: 'Spam Execute',
				labelTooltip: 'Use Execute whenever possible during Execute Phase',
				changedEvent: (player: Player<Spec.SpecWarrior>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
				getValue: (player: Player<Spec.SpecWarrior>) => player.getRotation().spamExecute,
				setValue: (eventID: EventID, player: Player<Spec.SpecWarrior>, newValue: boolean) => {
					const newRotation = player.getRotation();
					newRotation.spamExecute = newValue;
					player.setRotation(eventID, newRotation);
				},
				showWhen: (player: Player<Spec.SpecWarrior>) => player.getTalents().mortalStrike,
			},
		},
		{
			type: 'boolean' as const,
			cssClass: 'slam-over-exec-picker',
			getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
			config: {
				label: 'Slam Over Execute',
				labelTooltip: 'Use Slam Over Execute when Taste for Blood Procs in Execute Phase.',
				changedEvent: (player: Player<Spec.SpecWarrior>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
				getValue: (player: Player<Spec.SpecWarrior>) => player.getRotation().useSlamOverExecute,
				setValue: (eventID: EventID, player: Player<Spec.SpecWarrior>, newValue: boolean) => {
					const newRotation = player.getRotation();
					newRotation.useSlamOverExecute = newValue;
					player.setRotation(eventID, newRotation);
				},
				showWhen: (player: Player<Spec.SpecWarrior>) => player.getTalents().bloodthirst,
			},
		},
		{
			type: 'enum' as const, cssClass: 'sunder-armor-picker',
			getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
			config: {
				label: 'Sunder Armor',
				values: [
					{ name: 'Never', value: SunderArmor.SunderArmorNone },
					{ name: 'Help Stack', value: SunderArmor.SunderArmorHelpStack },
					{ name: 'Maintain Debuff', value: SunderArmor.SunderArmorMaintain },
				],
				changedEvent: (player: Player<Spec.SpecWarrior>) => player.rotationChangeEmitter,
				getValue: (player: Player<Spec.SpecWarrior>) => player.getRotation().sunderArmor,
				setValue: (eventID: EventID, player: Player<Spec.SpecWarrior>, newValue: number) => {
					const newRotation = player.getRotation();
					newRotation.sunderArmor = newValue;
					player.setRotation(eventID, newRotation);
				},
			},
		},
		{
			type: 'boolean' as const,
			cssClass: 'maintain-demo-shout-picker',
			getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
			config: {
				label: 'Maintain Demo Shout',
				labelTooltip: 'Keep Demo Shout active on the primary target.',
				changedEvent: (player: Player<Spec.SpecWarrior>) => player.rotationChangeEmitter,
				getValue: (player: Player<Spec.SpecWarrior>) => player.getRotation().maintainDemoShout,
				setValue: (eventID: EventID, player: Player<Spec.SpecWarrior>, newValue: boolean) => {
					const newRotation = player.getRotation();
					newRotation.maintainDemoShout = newValue;
					player.setRotation(eventID, newRotation);
				},
			},
		},
		{
			type: 'boolean' as const,
			cssClass: 'maintain-thunder-clap-picker',
			getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
			config: {
				label: 'Maintain Thunder Clap',
				labelTooltip: 'Keep Thunder Clap active on the primary target.',
				changedEvent: (player: Player<Spec.SpecWarrior>) => player.rotationChangeEmitter,
				getValue: (player: Player<Spec.SpecWarrior>) => player.getRotation().maintainThunderClap,
				setValue: (eventID: EventID, player: Player<Spec.SpecWarrior>, newValue: boolean) => {
					const newRotation = player.getRotation();
					newRotation.maintainThunderClap = newValue;
					player.setRotation(eventID, newRotation);
				},
			},
		},
	],
};
