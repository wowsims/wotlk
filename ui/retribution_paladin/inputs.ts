import { Spec } from '/tbc/core/proto/common.js';
import { Player } from '/tbc/core/player.js';
import { EventID } from '/tbc/core/typed_event.js';
import { IndividualSimUI } from '/tbc/core/individual_sim_ui.js';

import {
	PaladinAura as PaladinAura,
	RetributionPaladin_Rotation as RetributionPaladinRotation,
	RetributionPaladin_Options as RetributionPaladinOptions,
	RetributionPaladin_Rotation_ConsecrationRank as ConsecrationRank,
	RetributionPaladin_Options_Judgement as Judgement,
} from '/tbc/core/proto/paladin.js';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.
export const RetributionPaladinRotationConfig = {
	inputs: [
		{
			type: 'enum' as const, cssClass: 'consecration-rank-picker',
			getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
			config: {
				label: 'Consecration Rank',
				labelTooltip: 'Use specified rank of Consecration during filler spell windows.',
				values: [
					{
						name: 'None', value: ConsecrationRank.None,
					},
					{
						name: 'Rank 1', value: ConsecrationRank.Rank1,
					},
					{
						name: 'Rank 4', value: ConsecrationRank.Rank4,
					},
					{
						name: 'Rank 6', value: ConsecrationRank.Rank6,
					},
				],
				changedEvent: (player: Player<Spec.SpecRetributionPaladin>) => player.rotationChangeEmitter,
				getValue: (player: Player<Spec.SpecRetributionPaladin>) => player.getRotation().consecrationRank,
				setValue: (eventID: EventID, player: Player<Spec.SpecRetributionPaladin>, newValue: number) => {
					const newRotation = player.getRotation();
					newRotation.consecrationRank = newValue;
					player.setRotation(eventID, newRotation);
				},
			},
		},
		{
			type: 'boolean' as const, cssClass: 'exorcism-picker',
			getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
			config: {
				label: 'Use Exorcism',
				labelTooltip: 'Use Exorcism during filler spell windows. Will only be used if the boss mob type is Undead or Demon.',
				changedEvent: (player: Player<Spec.SpecRetributionPaladin>) => player.rotationChangeEmitter,
				getValue: (player: Player<Spec.SpecRetributionPaladin>) => player.getRotation().useExorcism,
				setValue: (eventID: EventID, player: Player<Spec.SpecRetributionPaladin>, newValue: boolean) => {
					const newRotation = player.getRotation();
					newRotation.useExorcism = newValue;
					player.setRotation(eventID, newRotation);
				},
			},
		}
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
		changedEvent: (player: Player<Spec.SpecRetributionPaladin>) => player.specOptionsChangeEmitter,
		getValue: (player: Player<Spec.SpecRetributionPaladin>) => player.getSpecOptions().aura,
		setValue: (eventID: EventID, player: Player<Spec.SpecRetributionPaladin>, newValue: number) => {
			const newOptions = player.getSpecOptions();
			newOptions.aura = newValue;
			player.setSpecOptions(eventID, newOptions);
		},
	},
}

export const JudgementSelection = {
	type: 'enum' as const, cssClass: 'judgement-picker',
	getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
	config: {
		label: 'Judgement',
		labelTooltip: 'Judgement debuff you will use on the target during the encounter. \
		If Improved Seal of the Crusader is talented, the Improved Judgement of the Crusader debuff will be applied.',
		values: [
			{
				name: 'None', value: Judgement.None,
			},
			{
				name: 'Judgement of Wisdom', value: Judgement.Wisdom,
			},
			{
				name: 'Judgement of the Crusader', value: Judgement.Crusader,
			},
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

export const CrusaderStrikeDelayMS = {
	type: 'number' as const, cssClass: 'cs-delay-picker',
	getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
	config: {
		label: 'Crusader Strike Delay (MS)',
		labelTooltip: 'Maximum time (in miliseconds) Crusader Strike will be delayed in order to seal twist. Experiment with values between 0 - 3000 miliseconds.',
		changedEvent: (player: Player<Spec.SpecRetributionPaladin>) => player.specOptionsChangeEmitter,
		getValue: (player: Player<Spec.SpecRetributionPaladin>) => player.getSpecOptions().crusaderStrikeDelayMs,
		setValue: (eventID: EventID, player: Player<Spec.SpecRetributionPaladin>, newValue: number) => {
			const newOptions = player.getSpecOptions();
			newOptions.crusaderStrikeDelayMs = newValue;
			player.setSpecOptions(eventID, newOptions);
		},
	},
}

/*** Leave this for now. We'll ignore HasteLeeway for initial release, but we might come back to it at some point  ***/

// export const HasteLeewayMS = {
// 	type: 'number' as const, cssClass: 'haste-leeway-picker',
// 	getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
// 	config: {
// 		label: 'Haste Leeway (MS)',
// 		labelTooltip: "Arbitrary value used to account for haste procs preventing seal twists. Experiment with values between 100 - 200 miliseconds.\nDo not modify this value if you do not understand it's use.",
// 		changedEvent: (player: Player<Spec.SpecRetributionPaladin>) => player.specOptionsChangeEmitter,
// 		getValue: (player: Player<Spec.SpecRetributionPaladin>) => player.getSpecOptions().hasteLeewayMs,
// 		setValue: (eventID: EventID, player: Player<Spec.SpecRetributionPaladin>, newValue: number) => {
// 			const newOptions = player.getSpecOptions();
// 			newOptions.hasteLeewayMs = newValue;
// 			player.setSpecOptions(eventID, newOptions);
// 		},
// 	},
// }

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
