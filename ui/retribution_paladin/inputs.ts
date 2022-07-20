import { Spec } from '/wotlk/core/proto/common.js';
import { Player } from '/wotlk/core/player.js';
import { EventID } from '/wotlk/core/typed_event.js';
import { IndividualSimUI } from '/wotlk/core/individual_sim_ui.js';

import {
	PaladinAura as PaladinAura,
	RetributionPaladin_Rotation as RetributionPaladinRotation,
	RetributionPaladin_Options as RetributionPaladinOptions,
	PaladinJudgement as PaladinJudgement,
	PaladinSeal,
} from '/wotlk/core/proto/paladin.js';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.
export const RetributionPaladinRotationConfig = {
	inputs: [
	],
}

export const AuraSelection = {
	type: 'enum' as const, cssClass: 'aura-picker',
	getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
	config: {
		label: 'Aura',
		values: [
			{ name: 'None', value: PaladinAura.NoPaladinAura },
			{ name: 'Retribution Aura', value: PaladinAura.RetributionAura },
		],
		changedEvent: (player: Player<Spec.SpecRetributionPaladin>) => player.specOptionsChangeEmitter,
		getValue: (player: Player<Spec.SpecRetributionPaladin>) => player.getSpecOptions().aura,
		setValue: (eventID: EventID, player: Player<Spec.SpecRetributionPaladin>, newValue: number) => {
			const newOptions = player.getSpecOptions();
			newOptions.aura = newValue;
			player.setSpecOptions(eventID, newOptions);
		},
	},
}

export const StartingSealSelection = {
	type: 'enum' as const, cssClass: 'starting-seal-picker',
	getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
	config: {
		label: 'Seal',
		labelTooltip: 'The seal active before encounter',
		values: [
			{
				name: 'Vengeance', value: PaladinSeal.Vengeance,
			},
			{
				name: 'Command', value: PaladinSeal.Command,
			},
			{
				name: 'Righteousness', value: PaladinSeal.Righteousness,
			},
		],
		changedEvent: (player: Player<Spec.SpecRetributionPaladin>) => player.specOptionsChangeEmitter,
		getValue: (player: Player<Spec.SpecRetributionPaladin>) => player.getSpecOptions().seal,
		setValue: (eventID: EventID, player: Player<Spec.SpecRetributionPaladin>, newValue: number) => {
			const newOptions = player.getSpecOptions();
			newOptions.seal = newValue;
			player.setSpecOptions(eventID, newOptions);
		},
	},
}

export const DivinePleaSelection = {
	type: 'boolean' as const, cssClass: 'divine-plea-picker',
	getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
	config: {
		label: 'Divine Plea',
		labelTooltip: 'Whether or not to maintain Divine Plea',
		values: [
			{ name: 'Yes', value: true },
			{ name: 'No', value: false },
		],
		changedEvent: (player: Player<Spec.SpecRetributionPaladin>) => player.specOptionsChangeEmitter,
		getValue: (player: Player<Spec.SpecRetributionPaladin>) => player.getSpecOptions().useDivinePlea,
		setValue: (eventID: EventID, player: Player<Spec.SpecRetributionPaladin>, newValue: boolean) => {
			const newOptions = player.getSpecOptions();
			newOptions.useDivinePlea = newValue;
			player.setSpecOptions(eventID, newOptions);
		},
	},
}

export const JudgementSelection = {
	type: 'enum' as const, cssClass: 'judgement-picker',
	getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
	config: {
		label: 'Judgement',
		labelTooltip: 'Judgement debuff you will use on the target during the encounter.',
		values: [
			{ name: 'Wisdom', value: PaladinJudgement.JudgementOfWisdom },
			{ name: 'Light', value: PaladinJudgement.JudgementOfLight },
		],
		changedEvent: (player: Player<Spec.SpecRetributionPaladin>) => player.specOptionsChangeEmitter,
		getValue: (player: Player<Spec.SpecRetributionPaladin>) => player.getSpecOptions().judgement,
		setValue: (eventID: EventID, player: Player<Spec.SpecRetributionPaladin>, newValue: number) => {
			const newOptions = player.getSpecOptions();
			newOptions.judgement = newValue;
			player.setSpecOptions(eventID, newOptions);
		},
	},
}

export const DamgeTakenPerSecond = {
	type: 'number' as const, cssClass: 'damage-taken-picker',
	getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
	config: {
		label: 'Damage Taken Per Second',
		labelTooltip: "Damage taken per second across the encounter. Used to model mana regeneration from Spiritual Attunement. This value should NOT include damage taken from Seal of Blood / Judgement of Blood. Leave at 0 if you do not take damage during the encounter.",
		changedEvent: (player: Player<Spec.SpecRetributionPaladin>) => player.specOptionsChangeEmitter,
		getValue: (player: Player<Spec.SpecRetributionPaladin>) => player.getSpecOptions().damageTakenPerSecond,
		setValue: (eventID: EventID, player: Player<Spec.SpecRetributionPaladin>, newValue: number) => {
			const newOptions = player.getSpecOptions();
			newOptions.damageTakenPerSecond = newValue;
			player.setSpecOptions(eventID, newOptions);
		},
	},
}
