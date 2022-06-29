import { IconPickerConfig } from '/tbc/core/components/icon_picker.js';
import { RaidTarget } from '/tbc/core/proto/common.js';
import { Spec } from '/tbc/core/proto/common.js';
import { NO_TARGET } from '/tbc/core/proto_utils/utils.js';
import { ActionId } from '/tbc/core/proto_utils/action_id.js';
import { Player } from '/tbc/core/player.js';
import { Sim } from '/tbc/core/sim.js';
import { EventID, TypedEvent } from '/tbc/core/typed_event.js';
import { IndividualSimUI } from '/tbc/core/individual_sim_ui.js';
import { Target } from '/tbc/core/target.js';

import {
	WarriorShout,
	WarriorTalents as WarriorTalents,
	Warrior,
	Warrior_Rotation as WarriorRotation,
	Warrior_Rotation_SunderArmor as SunderArmor,
	Warrior_Options as WarriorOptions,
} from '/tbc/core/proto/warrior.js';

import * as Presets from './presets.js';
import { SimUI } from '../core/sim_ui.js';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const Recklessness = {
	id: ActionId.fromSpellId(1719),
	states: 2,
	extraCssClasses: [
		'warrior-recklessness-picker',
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

export const PrecastShout = {
	type: 'boolean' as const,
	getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
	config: {
		extraCssClasses: [
			'precast-shout-picker',
		],
		label: 'Precast Shout',
		labelTooltip: 'Selected shout is cast 10 seconds before combat starts.',
		changedEvent: (player: Player<Spec.SpecWarrior>) => player.specOptionsChangeEmitter,
		getValue: (player: Player<Spec.SpecWarrior>) => player.getSpecOptions().precastShout,
		setValue: (eventID: EventID, player: Player<Spec.SpecWarrior>, newValue: boolean) => {
			const newOptions = player.getSpecOptions();
			newOptions.precastShout = newValue;
			player.setSpecOptions(eventID, newOptions);
		},
	},
};

export const PrecastShoutWithSapphire = {
	type: 'boolean' as const,
	getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
	config: {
		extraCssClasses: [
			'precast-shout-with-sapphire-picker',
		],
		label: 'Precast with Sapphire',
		labelTooltip: 'Snapshot bonus from Solarian\'s Sapphire (+70 attack power) with precast shout.',
		changedEvent: (player: Player<Spec.SpecWarrior>) => TypedEvent.onAny([player.specOptionsChangeEmitter, player.gearChangeEmitter]),
		getValue: (player: Player<Spec.SpecWarrior>) => player.getSpecOptions().precastShoutSapphire,
		setValue: (eventID: EventID, player: Player<Spec.SpecWarrior>, newValue: boolean) => {
			const newOptions = player.getSpecOptions();
			newOptions.precastShoutSapphire = newValue;
			player.setSpecOptions(eventID, newOptions);
		},
		enableWhen: (player: Player<Spec.SpecWarrior>) => player.getSpecOptions().shout == WarriorShout.WarriorShoutBattle && player.getSpecOptions().precastShout && !player.getGear().hasTrinket(30446),
	},
};

export const PrecastShoutWithT2 = {
	type: 'boolean' as const,
	getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
	config: {
		extraCssClasses: [
			'precast-shout-with-t2-picker',
		],
		label: 'Precast with T2',
		labelTooltip: 'Snapshot T2 set bonus (+30 attack power) with precast shout.',
		changedEvent: (player: Player<Spec.SpecWarrior>) => player.specOptionsChangeEmitter,
		getValue: (player: Player<Spec.SpecWarrior>) => player.getSpecOptions().precastShoutT2,
		setValue: (eventID: EventID, player: Player<Spec.SpecWarrior>, newValue: boolean) => {
			const newOptions = player.getSpecOptions();
			newOptions.precastShoutT2 = newValue;
			player.setSpecOptions(eventID, newOptions);
		},
		enableWhen: (player: Player<Spec.SpecWarrior>) => player.getSpecOptions().shout == WarriorShout.WarriorShoutBattle && player.getSpecOptions().precastShout,
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
			cssClass: 'overpower-picker',
			getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
			config: {
				label: 'Use Overpower',
				labelTooltip: 'Use Overpower when available.',
				changedEvent: (player: Player<Spec.SpecWarrior>) => player.rotationChangeEmitter,
				getValue: (player: Player<Spec.SpecWarrior>) => player.getRotation().useOverpower,
				setValue: (eventID: EventID, player: Player<Spec.SpecWarrior>, newValue: boolean) => {
					const newRotation = player.getRotation();
					newRotation.useOverpower = newValue;
					player.setRotation(eventID, newRotation);
				},
			},
		},
		{
			type: 'boolean' as const,
			cssClass: 'hamstring-picker',
			getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
			config: {
				label: 'Use Hamstring',
				labelTooltip: 'Use Hamstring on free globals.',
				changedEvent: (player: Player<Spec.SpecWarrior>) => player.rotationChangeEmitter,
				getValue: (player: Player<Spec.SpecWarrior>) => player.getRotation().useHamstring,
				setValue: (eventID: EventID, player: Player<Spec.SpecWarrior>, newValue: boolean) => {
					const newRotation = player.getRotation();
					newRotation.useHamstring = newValue;
					player.setRotation(eventID, newRotation);
				},
			},
		},
		{
			type: 'boolean' as const,
			cssClass: 'slam-picker',
			getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
			config: {
				label: 'Use Slam',
				labelTooltip: 'Use Slam whenever possible.',
				changedEvent: (player: Player<Spec.SpecWarrior>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
				getValue: (player: Player<Spec.SpecWarrior>) => player.getRotation().useSlam,
				setValue: (eventID: EventID, player: Player<Spec.SpecWarrior>, newValue: boolean) => {
					const newRotation = player.getRotation();
					newRotation.useSlam = newValue;
					player.setRotation(eventID, newRotation);
				},
				showWhen: (player: Player<Spec.SpecWarrior>) => player.getTalents().improvedSlam == 2,
			},
		},
		{
			type: 'boolean' as const,
			cssClass: 'prioritize-ww-picker',
			getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
			config: {
				label: 'Prioritize WW',
				labelTooltip: 'Prioritize Whirlwind over Bloodthirst or Mortal Strike.',
				changedEvent: (player: Player<Spec.SpecWarrior>) => player.rotationChangeEmitter,
				getValue: (player: Player<Spec.SpecWarrior>) => player.getRotation().prioritizeWw,
				setValue: (eventID: EventID, player: Player<Spec.SpecWarrior>, newValue: boolean) => {
					const newRotation = player.getRotation();
					newRotation.prioritizeWw = newValue;
					player.setRotation(eventID, newRotation);
				},
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
			cssClass: 'overpower-rage-threshold',
			getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
			config: {
				label: 'Overpower rage threshold',
				labelTooltip: 'Use Overpower when rage is below a point.',
				changedEvent: (player: Player<Spec.SpecWarrior>) => player.rotationChangeEmitter,
				getValue: (player: Player<Spec.SpecWarrior>) => player.getRotation().overpowerRageThreshold,
				setValue: (eventID: EventID, player: Player<Spec.SpecWarrior>, newValue: number) => {
					const newRotation = player.getRotation();
					newRotation.overpowerRageThreshold = newValue;
					player.setRotation(eventID, newRotation);
				},
				showWhen: (player: Player<Spec.SpecWarrior>) => player.getRotation().useOverpower,
			},
		},
		{
			type: 'number' as const,
			cssClass: 'hamstring-rage-threshold',
			getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
			config: {
				label: 'Hamstring rage threshold',
				labelTooltip: 'Hamstring will only be used when rage is larger than this value.',
				changedEvent: (player: Player<Spec.SpecWarrior>) => player.rotationChangeEmitter,
				getValue: (player: Player<Spec.SpecWarrior>) => player.getRotation().hamstringRageThreshold,
				setValue: (eventID: EventID, player: Player<Spec.SpecWarrior>, newValue: number) => {
					const newRotation = player.getRotation();
					newRotation.hamstringRageThreshold = newValue;
					player.setRotation(eventID, newRotation);
				},
				showWhen: (player: Player<Spec.SpecWarrior>) => player.getRotation().useHamstring,
			},
		},
		{
			type: 'number' as const,
			cssClass: 'slam-latency',
			getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
			config: {
				label: 'Slam Latency',
				labelTooltip: 'Time between MH swing and start of the Slam cast, in milliseconds.',
				changedEvent: (player: Player<Spec.SpecWarrior>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
				getValue: (player: Player<Spec.SpecWarrior>) => player.getRotation().slamLatency,
				setValue: (eventID: EventID, player: Player<Spec.SpecWarrior>, newValue: number) => {
					const newRotation = player.getRotation();
					newRotation.slamLatency = newValue;
					player.setRotation(eventID, newRotation);
				},
				showWhen: (player: Player<Spec.SpecWarrior>) => player.getRotation().useSlam && player.getTalents().improvedSlam == 2,
			},
		},
		{
			type: 'number' as const,
			cssClass: 'slam-gcd-delay',
			getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
			config: {
				extraCssClasses: [
					'experimental',
				],
				label: 'Slam GCD Delay',
				labelTooltip: 'Amount of time Slam may delay the GCD, in milliseconds.',
				changedEvent: (player: Player<Spec.SpecWarrior>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
				getValue: (player: Player<Spec.SpecWarrior>) => player.getRotation().slamGcdDelay,
				setValue: (eventID: EventID, player: Player<Spec.SpecWarrior>, newValue: number) => {
					const newRotation = player.getRotation();
					newRotation.slamGcdDelay = newValue;
					player.setRotation(eventID, newRotation);
				},
				showWhen: (player: Player<Spec.SpecWarrior>) => player.getRotation().useSlam && player.getTalents().improvedSlam == 2,
			},
		},
		{
			type: 'number' as const,
			cssClass: 'slam-ms-ww-delay',
			getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
			config: {
				extraCssClasses: [
					'experimental',
				],
				label: 'Slam MS+WW Delay',
				labelTooltip: 'Amount of time Slam may delay MS+WW, in milliseconds.',
				changedEvent: (player: Player<Spec.SpecWarrior>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
				getValue: (player: Player<Spec.SpecWarrior>) => player.getRotation().slamMsWwDelay,
				setValue: (eventID: EventID, player: Player<Spec.SpecWarrior>, newValue: number) => {
					const newRotation = player.getRotation();
					newRotation.slamMsWwDelay = newValue;
					player.setRotation(eventID, newRotation);
				},
				showWhen: (player: Player<Spec.SpecWarrior>) => player.getRotation().useSlam && player.getTalents().improvedSlam == 2,
			},
		},
		{
			type: 'number' as const,
			cssClass: 'rampage-duration-threshold',
			getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
			config: {
				label: 'Rampage Refresh Time',
				labelTooltip: 'Refresh Rampage when the remaining duration is less than this amount of time (seconds).',
				changedEvent: (player: Player<Spec.SpecWarrior>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
				getValue: (player: Player<Spec.SpecWarrior>) => player.getRotation().rampageCdThreshold,
				setValue: (eventID: EventID, player: Player<Spec.SpecWarrior>, newValue: number) => {
					const newRotation = player.getRotation();
					newRotation.rampageCdThreshold = newValue;
					player.setRotation(eventID, newRotation);
				},
				showWhen: (player: Player<Spec.SpecWarrior>) => player.getTalents().rampage,
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
			cssClass: 'ms-exec-picker-fury',
			getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
			config: {
				label: 'MS during Execute Phase',
				labelTooltip: 'Use Mortal Strike during Execute Phase.',
				changedEvent: (player: Player<Spec.SpecWarrior>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
				getValue: (player: Player<Spec.SpecWarrior>) => player.getRotation().useMsDuringExecute,
				setValue: (eventID: EventID, player: Player<Spec.SpecWarrior>, newValue: boolean) => {
					const newRotation = player.getRotation();
					newRotation.useMsDuringExecute = newValue;
					player.setRotation(eventID, newRotation);
				},
				showWhen: (player: Player<Spec.SpecWarrior>) => player.getTalents().mortalStrike,
			},
		},
		{
			type: 'boolean' as const,
			cssClass: 'ww-exec-picker',
			getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
			config: {
				label: 'WW during Execute Phase',
				labelTooltip: 'Use Whirlwind during Execute Phase.',
				changedEvent: (player: Player<Spec.SpecWarrior>) => player.rotationChangeEmitter,
				getValue: (player: Player<Spec.SpecWarrior>) => player.getRotation().useWwDuringExecute,
				setValue: (eventID: EventID, player: Player<Spec.SpecWarrior>, newValue: boolean) => {
					const newRotation = player.getRotation();
					newRotation.useWwDuringExecute = newValue;
					player.setRotation(eventID, newRotation);
				},
			},
		},
		{
			type: 'boolean' as const,
			cssClass: 'slam-exec-picker',
			getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
			config: {
				label: 'Slam during Execute Phase',
				labelTooltip: 'Use Slam during Execute Phase.',
				changedEvent: (player: Player<Spec.SpecWarrior>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
				getValue: (player: Player<Spec.SpecWarrior>) => player.getRotation().useSlamDuringExecute,
				setValue: (eventID: EventID, player: Player<Spec.SpecWarrior>, newValue: boolean) => {
					const newRotation = player.getRotation();
					newRotation.useSlamDuringExecute = newValue;
					player.setRotation(eventID, newRotation);
				},
				showWhen: (player: Player<Spec.SpecWarrior>) => player.getRotation().useSlam && player.getTalents().improvedSlam == 2,
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
