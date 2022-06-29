import { FeralDruid_Rotation_FinishingMove as FinishingMove } from '/tbc/core/proto/druid.js';
import { FeralDruid_Options as DruidOptions } from '/tbc/core/proto/druid.js';
import { RaidTarget } from '/tbc/core/proto/common.js';
import { Spec } from '/tbc/core/proto/common.js';
import { NO_TARGET } from '/tbc/core/proto_utils/utils.js';
import { ActionId } from '/tbc/core/proto_utils/action_id.js';
import { Player } from '/tbc/core/player.js';
import { Sim } from '/tbc/core/sim.js';
import { EventID, TypedEvent } from '/tbc/core/typed_event.js';
import { IndividualSimUI } from '/tbc/core/individual_sim_ui.js';
import { Target } from '/tbc/core/target.js';
import { getEnumValues } from '/tbc/core/utils.js';
import { ItemSlot } from '/tbc/core/proto/common.js';

// Helper function for identifying whether 2pT6 is equipped, which impacts allowed rotation options
function numThunderheartPieces(player: Player<Spec.SpecFeralDruid>): number {
	const gear = player.getGear();
	const itemIds = [31048, 31042, 31034, 31044, 31039, 34556, 34444, 34573];
	return gear.asArray().map(equippedItem => equippedItem?.item.id).filter(id => itemIds.includes(id!)).length
}

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const SelfInnervate = {
	id: ActionId.fromSpellId(29166),
	states: 2,
	extraCssClasses: [
		'self-innervate-picker',
		'within-raid-sim-hide',
	],
	changedEvent: (player: Player<Spec.SpecFeralDruid>) => player.specOptionsChangeEmitter,
	getValue: (player: Player<Spec.SpecFeralDruid>) => player.getSpecOptions().innervateTarget?.targetIndex != NO_TARGET,
	setValue: (eventID: EventID, player: Player<Spec.SpecFeralDruid>, newValue: boolean) => {
		const newOptions = player.getSpecOptions();
		newOptions.innervateTarget = RaidTarget.create({
			targetIndex: newValue ? 0 : NO_TARGET,
		});
		player.setSpecOptions(eventID, newOptions);
	},
};

export const LatencyMs = {
	type: 'number' as const,
	getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
	config: {
		extraCssClasses: [
			'latency-ms-picker',
		],
		label: 'Latency',
		labelTooltip: 'Player latency, in milliseconds. Adds a delay to actions that cannot be spell queued.',
		changedEvent: (player: Player<Spec.SpecFeralDruid>) => player.specOptionsChangeEmitter,
		getValue: (player: Player<Spec.SpecFeralDruid>) => player.getSpecOptions().latencyMs,
		setValue: (eventID: EventID, player: Player<Spec.SpecFeralDruid>, newValue: number) => {
			const newOptions = player.getSpecOptions();
			newOptions.latencyMs = newValue;
			player.setSpecOptions(eventID, newOptions);
		},
	},
};

export const FeralDruidRotationConfig = {
	inputs: [
		{
			type: 'enum' as const,
			getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
			config: {
				extraCssClasses: [
					'finishing-move-enum-picker',
				],
				label: 'Finishing Move',
				labelTooltip: 'Specify whether Rip or Ferocious Bite should be used as the primary finisher in the DPS rotation.',
				values: [
					{
						name: 'Rip', value: FinishingMove.Rip,
					},
					{
						name: 'Ferocious Bite', value: FinishingMove.Bite,
					},
					{
						name: 'None', value: FinishingMove.None,
					},
				],
				changedEvent: (player: Player<Spec.SpecFeralDruid>) => player.rotationChangeEmitter,
				getValue: (player: Player<Spec.SpecFeralDruid>) => player.getRotation().finishingMove,
				setValue: (eventID: EventID, player: Player<Spec.SpecFeralDruid>, newValue: number) => {
					const newRotation = player.getRotation();
					newRotation.finishingMove = newValue;
					player.setRotation(eventID, newRotation);
				},
			},
		},
		{
			type: 'boolean' as const,
			getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
			config: {
				extraCssClasses: [
					'biteweave-picker',
				],
				label: 'Enable Bite-weaving',
				labelTooltip: 'Spend Combo Points on Ferocious Bite when Rip is already applied on the target.',
				changedEvent: (player: Player<Spec.SpecFeralDruid>) => player.rotationChangeEmitter,
				getValue: (player: Player<Spec.SpecFeralDruid>) => player.getRotation().biteweave,
				setValue: (eventID: EventID, player: Player<Spec.SpecFeralDruid>, newValue: boolean) => {
					const newRotation = player.getRotation();
					newRotation.biteweave = newValue;
					player.setRotation(eventID, newRotation);
				},
				enableWhen: (player: Player<Spec.SpecFeralDruid>) => player.getRotation().finishingMove == FinishingMove.Rip,
			},
		},
		{
			type: 'boolean' as const,
			getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
			config: {
				extraCssClasses: [
					'ripweave-picker',
				],
				label: 'Enable Rip-weaving',
				labelTooltip: 'Spend Combo Points on Rip when at 52 Energy or above.',
				changedEvent: (player: Player<Spec.SpecFeralDruid>) => player.rotationChangeEmitter,
				getValue: (player: Player<Spec.SpecFeralDruid>) => player.getRotation().ripweave,
				setValue: (eventID: EventID, player: Player<Spec.SpecFeralDruid>, newValue: boolean) => {
					const newRotation = player.getRotation();
					newRotation.ripweave = newValue;
					player.setRotation(eventID, newRotation);
				},
				enableWhen: (player: Player<Spec.SpecFeralDruid>) => player.getRotation().finishingMove == FinishingMove.Bite,
			},
		},
		{
			type: 'enum' as const,
			getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
			config: {
				extraCssClasses: [
					'rip-cp-enum-picker',
				],
				label: 'Rip CP Threshold',
				labelTooltip: 'Minimum Combo Points to accumulate before casting Rip as a finisher.',
				values: [
					{
						name: '4', value: 4,
					},
					{
						name: '5', value: 5,
					},
				],
				changedEvent: (player: Player<Spec.SpecFeralDruid>) => player.rotationChangeEmitter,
				getValue: (player: Player<Spec.SpecFeralDruid>) => player.getRotation().ripMinComboPoints,
				setValue: (eventID: EventID, player: Player<Spec.SpecFeralDruid>, newValue: number) => {
					const newRotation = player.getRotation();
					newRotation.ripMinComboPoints = newValue;
					player.setRotation(eventID, newRotation);
				},
				enableWhen: (player: Player<Spec.SpecFeralDruid>) => (player.getRotation().finishingMove == FinishingMove.Rip) || (player.getRotation().ripweave && (player.getRotation().finishingMove != FinishingMove.None)),
			},
		},
		{
			type: 'enum' as const,
			getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
			config: {
				extraCssClasses: [
					'bite-cp-enum-picker',
				],
				label: 'Bite CP Threshold',
				labelTooltip: 'Minimum Combo Points to accumulate before casting Ferocious Bite as a finisher.',
				values: [
					{
						name: '4', value: 4,
					},
					{
						name: '5', value: 5,
					},
				],
				changedEvent: (player: Player<Spec.SpecFeralDruid>) => player.rotationChangeEmitter,
				getValue: (player: Player<Spec.SpecFeralDruid>) => player.getRotation().biteMinComboPoints,
				setValue: (eventID: EventID, player: Player<Spec.SpecFeralDruid>, newValue: number) => {
					const newRotation = player.getRotation();
					newRotation.biteMinComboPoints = newValue;
					player.setRotation(eventID, newRotation);
				},
				enableWhen: (player: Player<Spec.SpecFeralDruid>) => (player.getRotation().finishingMove == FinishingMove.Bite) || (player.getRotation().biteweave && (player.getRotation().finishingMove != FinishingMove.None)),
			},
		},
		{
			type: 'boolean' as const,
			getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
			config: {
				extraCssClasses: [
					'mangle-trick-picker',
				],
				label: 'Use Mangle trick',
				labelTooltip: 'Cast Mangle rather than Shred when between 50-56 Energy with 2pT6, or 60-61 Energy without 2pT6, and with less than 1 second remaining until the next Energy tick.',
				changedEvent: (player: Player<Spec.SpecFeralDruid>) => player.rotationChangeEmitter,
				getValue: (player: Player<Spec.SpecFeralDruid>) => player.getRotation().mangleTrick,
				setValue: (eventID: EventID, player: Player<Spec.SpecFeralDruid>, newValue: boolean) => {
					const newRotation = player.getRotation();
					newRotation.mangleTrick = newValue;
					player.setRotation(eventID, newRotation);
				},
			},
		},
		{
			type: 'boolean' as const,
			getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
			config: {
				extraCssClasses: [
					'rake-trick-picker',
				],
				label: 'Use Rake/Bite tricks',
				labelTooltip: 'Cast Rake or Ferocious Bite rather than powershifting when between 35-39 Energy without 2pT6, and with more than 1 second remaining until the next Energy tick.',
				changedEvent: (player: Player<Spec.SpecFeralDruid>) => player.changeEmitter,
				getValue: (player: Player<Spec.SpecFeralDruid>) => player.getRotation().rakeTrick,
				setValue: (eventID: EventID, player: Player<Spec.SpecFeralDruid>, newValue: boolean) => {
					const newRotation = player.getRotation();
					newRotation.rakeTrick = newValue;
					player.setRotation(eventID, newRotation);
				},
				enableWhen: (player: Player<Spec.SpecFeralDruid>) => numThunderheartPieces(player) < 2,
			},
		},
		{
			type: 'boolean' as const,
			getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
			config: {
				extraCssClasses: [
					'maintain-faerie-fire-picker',
				],
				label: 'Maintain Faerie Fire',
				labelTooltip: 'Use Faerie Fire whenever it is not active on the target.',
				changedEvent: (player: Player<Spec.SpecFeralDruid>) => player.rotationChangeEmitter,
				getValue: (player: Player<Spec.SpecFeralDruid>) => player.getRotation().maintainFaerieFire,
				setValue: (eventID: EventID, player: Player<Spec.SpecFeralDruid>, newValue: boolean) => {
					const newRotation = player.getRotation();
					newRotation.maintainFaerieFire = newValue;
					player.setRotation(eventID, newRotation);
				},
			},
		},
	],
};
