import { Spec } from '/tbc/core/proto/common.js';
import { Player } from '/tbc/core/player.js';
import { EventID } from '/tbc/core/typed_event.js';
import { IndividualSimUI } from '/tbc/core/individual_sim_ui.js';
import { ActionId } from '/tbc/core/proto_utils/action_id.js';

import {
	PaladinAura as PaladinAura,
	PaladinJudgement as PaladinJudgement,
	ProtectionPaladin_Rotation as ProtectionPaladinRotation,
	ProtectionPaladin_Options as ProtectionPaladinOptions,
} from '/tbc/core/proto/paladin.js';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.
export const ProtectionPaladinRotationConfig = {
	inputs: [
		{
			type: 'enum' as const, cssClass: 'consecration-rank-picker',
			getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
			config: {
				label: 'Consecration Rank',
				labelTooltip: 'Use specified rank of Consecration during filler spell windows.',
				values: [
					{
						name: 'None', value: 0,
					},
					{
						name: 'Rank 1', value: 1,
					},
					{
						name: 'Rank 4', value: 4,
					},
					{
						name: 'Rank 6', value: 6,
					},
				],
				changedEvent: (player: Player<Spec.SpecProtectionPaladin>) => player.rotationChangeEmitter,
				getValue: (player: Player<Spec.SpecProtectionPaladin>) => player.getRotation().consecrationRank,
				setValue: (eventID: EventID, player: Player<Spec.SpecProtectionPaladin>, newValue: number) => {
					const newRotation = player.getRotation();
					newRotation.consecrationRank = newValue;
					player.setRotation(eventID, newRotation);
				},
			},
		},
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
		{
			type: 'boolean' as const, cssClass: 'exorcism-picker',
			getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
			config: {
				label: 'Use Exorcism',
				labelTooltip: 'Includes Exorcism in the rotation. Will only be used if the primary target is an Undead or Demon type.',
				changedEvent: (player: Player<Spec.SpecProtectionPaladin>) => player.rotationChangeEmitter,
				getValue: (player: Player<Spec.SpecProtectionPaladin>) => player.getRotation().useExorcism,
				setValue: (eventID: EventID, player: Player<Spec.SpecProtectionPaladin>, newValue: boolean) => {
					const newRotation = player.getRotation();
					newRotation.useExorcism = newValue;
					player.setRotation(eventID, newRotation);
				},
			},
		},
		{
			type: 'enum' as const, cssClass: 'mantain-judgement-picker',
			getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
			config: {
				label: 'Maintain Judgement',
				values: [
					{ name: 'None', value: PaladinJudgement.NoPaladinJudgement },
					{ name: 'Wisdom', value: PaladinJudgement.JudgementOfWisdom },
					{ name: 'Light', value: PaladinJudgement.JudgementOfLight },
				],
				changedEvent: (player: Player<Spec.SpecProtectionPaladin>) => player.rotationChangeEmitter,
				getValue: (player: Player<Spec.SpecProtectionPaladin>) => player.getRotation().maintainJudgement,
				setValue: (eventID: EventID, player: Player<Spec.SpecProtectionPaladin>, newValue: number) => {
					const newRotation = player.getRotation();
					newRotation.maintainJudgement = newValue;
					player.setRotation(eventID, newRotation);
				},
			},
		},
	],
}

export const AuraSelection = {
	type: 'enum' as const, cssClass: 'aura-picker',
	getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
	config: {
		label: 'Aura',
		values: [
			{ name: 'None', value: PaladinAura.NoPaladinAura },
			{ name: 'Sanctity Aura', value: PaladinAura.SanctityAura },
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
