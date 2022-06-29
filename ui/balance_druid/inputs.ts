import { BalanceDruid_Rotation_PrimarySpell as PrimarySpell } from '/tbc/core/proto/druid.js';
import { BalanceDruid_Options as DruidOptions } from '/tbc/core/proto/druid.js';
import { RaidTarget } from '/tbc/core/proto/common.js';
import { Spec } from '/tbc/core/proto/common.js';
import { NO_TARGET } from '/tbc/core/proto_utils/utils.js';
import { ActionId } from '/tbc/core/proto_utils/action_id.js';
import { Player } from '/tbc/core/player.js';
import { Sim } from '/tbc/core/sim.js';
import { EventID, TypedEvent } from '/tbc/core/typed_event.js';
import { IndividualSimUI } from '/tbc/core/individual_sim_ui.js';
import { Target } from '/tbc/core/target.js';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const SelfInnervate = {
	id: ActionId.fromSpellId(29166),
	states: 2,
	extraCssClasses: [
		'self-innervate-picker',
		'within-raid-sim-hide',
	],
	changedEvent: (player: Player<Spec.SpecBalanceDruid>) => player.specOptionsChangeEmitter,
	getValue: (player: Player<Spec.SpecBalanceDruid>) => player.getSpecOptions().innervateTarget?.targetIndex != NO_TARGET,
	setValue: (eventID: EventID, player: Player<Spec.SpecBalanceDruid>, newValue: boolean) => {
		const newOptions = player.getSpecOptions();
		newOptions.innervateTarget = RaidTarget.create({
			targetIndex: newValue ? 0 : NO_TARGET,
		});
		player.setSpecOptions(eventID, newOptions);
	},
};

export const BalanceDruidRotationConfig = {
	inputs: [
		{
			type: 'enum' as const,
			getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
			config: {
				extraCssClasses: [
					'primary-spell-enum-picker',
				],
				label: 'Primary Spell',
				labelTooltip: 'If set to \'Adaptive\', will dynamically adjust rotation based on available mana.',
				values: [
					{
						name: 'Adaptive', value: PrimarySpell.Adaptive,
					},
					{
						name: 'Starfire', value: PrimarySpell.Starfire,
					},
					{
						name: 'Starfire R6', value: PrimarySpell.Starfire6,
					},
					{
						name: 'Wrath', value: PrimarySpell.Wrath,
					},
				],
				changedEvent: (player: Player<Spec.SpecBalanceDruid>) => player.rotationChangeEmitter,
				getValue: (player: Player<Spec.SpecBalanceDruid>) => player.getRotation().primarySpell,
				setValue: (eventID: EventID, player: Player<Spec.SpecBalanceDruid>, newValue: number) => {
					const newRotation = player.getRotation();
					newRotation.primarySpell = newValue;
					player.setRotation(eventID, newRotation);
				},
			},
		},
		{
			type: 'boolean' as const,
			getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
			config: {
				extraCssClasses: [
					'moonfire-picker',
				],
				label: 'Use Moonfire',
				labelTooltip: 'Use Moonfire as the next cast after the dot expires.',
				changedEvent: (player: Player<Spec.SpecBalanceDruid>) => player.rotationChangeEmitter,
				getValue: (player: Player<Spec.SpecBalanceDruid>) => player.getRotation().moonfire,
				setValue: (eventID: EventID, player: Player<Spec.SpecBalanceDruid>, newValue: boolean) => {
					const newRotation = player.getRotation();
					newRotation.moonfire = newValue;
					player.setRotation(eventID, newRotation);
				},
				enableWhen: (player: Player<Spec.SpecBalanceDruid>) => player.getRotation().primarySpell != PrimarySpell.Adaptive,
			},
		},
		{
			type: 'boolean' as const,
			getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
			config: {
				extraCssClasses: [
					'faerie-fire-picker',
				],
				label: 'Use Faerie Fire',
				labelTooltip: 'Keep Faerie Fire active on the primary target.',
				changedEvent: (player: Player<Spec.SpecBalanceDruid>) => player.rotationChangeEmitter,
				getValue: (player: Player<Spec.SpecBalanceDruid>) => player.getRotation().faerieFire,
				setValue: (eventID: EventID, player: Player<Spec.SpecBalanceDruid>, newValue: boolean) => {
					const newRotation = player.getRotation();
					newRotation.faerieFire = newValue;
					player.setRotation(eventID, newRotation);
				},
			},
		},
		{
			type: 'boolean' as const,
			getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
			config: {
				extraCssClasses: [
					'insect-swarm-picker',
				],
				label: 'Use Insect Swarm',
				labelTooltip: 'Keep Insect Swarm active on the primary target.',
				changedEvent: (player: Player<Spec.SpecBalanceDruid>) => player.rotationChangeEmitter,
				getValue: (player: Player<Spec.SpecBalanceDruid>) => player.getRotation().insectSwarm,
				setValue: (eventID: EventID, player: Player<Spec.SpecBalanceDruid>, newValue: boolean) => {
					const newRotation = player.getRotation();
					newRotation.insectSwarm = newValue;
					player.setRotation(eventID, newRotation);
				},
				enableWhen: (player: Player<Spec.SpecBalanceDruid>) => player.getTalents().insectSwarm,
			},
		},
		{
			type: 'boolean' as const,
			getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
			config: {
				extraCssClasses: [
					'hurricane-picker',
				],
				label: 'Use Hurricane',
				labelTooltip: 'Casts Hurricane on cooldown.',
				changedEvent: (player: Player<Spec.SpecBalanceDruid>) => player.specOptionsChangeEmitter,
				getValue: (player: Player<Spec.SpecBalanceDruid>) => player.getRotation().hurricane,
				setValue: (eventID: EventID, player: Player<Spec.SpecBalanceDruid>, newValue: boolean) => {
					const newRotation = player.getRotation();
					newRotation.hurricane = newValue;
					player.setRotation(eventID, newRotation);
				},
			},
		},
		{
			type: 'boolean' as const,
			getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
			config: {
				extraCssClasses: [
					'battle-res-picker',
				],
				label: 'Use Battle Res',
				labelTooltip: 'Cast Battle Res on an ally sometime during the encounter.',
				changedEvent: (player: Player<Spec.SpecBalanceDruid>) => player.specOptionsChangeEmitter,
				getValue: (player: Player<Spec.SpecBalanceDruid>) => player.getSpecOptions().battleRes,
				setValue: (eventID: EventID, player: Player<Spec.SpecBalanceDruid>, newValue: boolean) => {
					const newOptions = player.getSpecOptions();
					newOptions.battleRes = newValue;
					player.setSpecOptions(eventID, newOptions);
				},
			},
		},
	],
};
