import { Spec } from '/wotlk/core/proto/common.js';
import { Player } from '/wotlk/core/player.js';
import { EventID } from '/wotlk/core/typed_event.js';
import { IndividualSimUI } from '/wotlk/core/individual_sim_ui.js';
import { ActionId } from '/wotlk/core/proto_utils/action_id.js';

import {
	PaladinAura as PaladinAura,
	PaladinSeal,
	PaladinJudgement as PaladinJudgement,
	ProtectionPaladin_Rotation as ProtectionPaladinRotation,
	ProtectionPaladin_Options as ProtectionPaladinOptions,
} from '/wotlk/core/proto/paladin.js';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.
export const ProtectionPaladinRotationConfig = {
	inputs: [
		{
			type: 'boolean' as const, cssClass: 'prioritize-holy-shield-picker',
			getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
			config: {
				label: 'Prio Holy Shield',
				labelTooltip: 'Uses Holy Shield as the highest priority spell. This is usually done when tanking a boss that can crush.',
				changedEvent: (player: Player<Spec.SpecProtectionPaladin>) => player.rotationChangeEmitter,
				getValue: (player: Player<Spec.SpecProtectionPaladin>) => player.getRotation().prioritizeHolyShield,
				setValue: (eventID: EventID, player: Player<Spec.SpecProtectionPaladin>, newValue: boolean) => {
					const newRotation = player.getRotation();
					newRotation.prioritizeHolyShield = newValue;
					player.setRotation(eventID, newRotation);
				},
			},
		},
	],
}

export const StartingSealSelection = {
	type: 'enum' as const, cssClass: 'starting-seal-picker',
	getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
	config: {
		label: 'Seal',
		values: [
			{ name: 'Vengeance', value: PaladinSeal.Vengeance },
			{ name: 'Command', value: PaladinSeal.Command },
		],
		changedEvent: (player: Player<Spec.SpecProtectionPaladin>) => player.rotationChangeEmitter,
		getValue: (player: Player<Spec.SpecProtectionPaladin>) => player.getSpecOptions().seal,
		setValue: (eventID: EventID, player: Player<Spec.SpecProtectionPaladin>, newValue: number) => {
			const newOptions = player.getSpecOptions();
			newOptions.seal = newValue;
			player.setSpecOptions(eventID, newOptions);
		},
	},
}


export const JudgementSelection = {
	type: 'enum' as const, cssClass: 'judgement-picker',
	getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
	config: {
		label: 'Judgement',
		values: [
			{ name: 'Wisdom', value: PaladinJudgement.JudgementOfWisdom },
			{ name: 'Light', value: PaladinJudgement.JudgementOfLight },
		],
		changedEvent: (player: Player<Spec.SpecProtectionPaladin>) => player.rotationChangeEmitter,
		getValue: (player: Player<Spec.SpecProtectionPaladin>) => player.getSpecOptions().judgement,
		setValue: (eventID: EventID, player: Player<Spec.SpecProtectionPaladin>, newValue: number) => {
			const newOptions = player.getSpecOptions();
			newOptions.judgement = newValue;
			player.setSpecOptions(eventID, newOptions);
		},
	},
}

export const AuraSelection = {
	type: 'enum' as const, cssClass: 'aura-picker',
	getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
	config: {
		label: 'Aura',
		values: [
			{ name: 'None', value: PaladinAura.NoPaladinAura },
			{ name: 'Devotion Aura', value: PaladinAura.DevotionAura },
			{ name: 'Retribution Aura', value: PaladinAura.RetributionAura },
		],
		changedEvent: (player: Player<Spec.SpecProtectionPaladin>) => player.specOptionsChangeEmitter,
		getValue: (player: Player<Spec.SpecProtectionPaladin>) => player.getSpecOptions().aura,
		setValue: (eventID: EventID, player: Player<Spec.SpecProtectionPaladin>, newValue: number) => {
			const newOptions = player.getSpecOptions();
			newOptions.aura = newValue;
			player.setSpecOptions(eventID, newOptions);
		},
	},
}

export const UseAvengingWrath = {
	type: 'boolean' as const, cssClass: 'use-avenging-wrath-picker',
	getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
	config: {
		label: 'Use Avenging Wrath',
		changedEvent: (player: Player<Spec.SpecProtectionPaladin>) => player.specOptionsChangeEmitter,
		getValue: (player: Player<Spec.SpecProtectionPaladin>) => player.getSpecOptions().useAvengingWrath,
		setValue: (eventID: EventID, player: Player<Spec.SpecProtectionPaladin>, newValue: boolean) => {
			const newOptions = player.getSpecOptions();
			newOptions.useAvengingWrath = newValue;
			player.setSpecOptions(eventID, newOptions);
		},
	},
};

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