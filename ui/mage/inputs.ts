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

import { Mage, Mage_Rotation as MageRotation, MageTalents as MageTalents, Mage_Options as MageOptions } from '/wotlk/core/proto/mage.js';
import { Mage_Rotation_Type as RotationType, Mage_Rotation_ArcaneRotation as ArcaneRotation, Mage_Rotation_FireRotation as FireRotation, Mage_Rotation_FrostRotation as FrostRotation, Mage_Rotation_AoeRotation as AoeRotation } from '/wotlk/core/proto/mage.js';
import { Mage_Rotation_FireRotation_PrimarySpell as PrimaryFireSpell } from '/wotlk/core/proto/mage.js';
import { Mage_Rotation_AoeRotation_Rotation as AoeRotationSpells } from '/wotlk/core/proto/mage.js';
import { Mage_Options_ArmorType as ArmorType } from '/wotlk/core/proto/mage.js';

import * as InputHelpers from '/wotlk/core/components/input_helpers.js';
import * as Presets from './presets.js';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const Armor = InputHelpers.makeSpecOptionsEnumIconInput<Spec.SpecMage, ArmorType>({
	fieldName: 'armor',
	values: [
		{ color: 'grey', value: ArmorType.NoArmor },
		{ actionId: ActionId.fromItemId(27125), value: ArmorType.MageArmor },
		{ actionId: ActionId.fromItemId(30482), value: ArmorType.MoltenArmor },
	],
});

export const EvocationTicks = InputHelpers.makeSpecOptionsNumberInput<Spec.SpecMage>({
	fieldName: 'evocationTicks',
	label: '# Evocation Ticks',
	labelTooltip: 'The number of ticks of Evocation to use, or 0 to use the full duration.',
});

export const FocusMagicUptime = InputHelpers.makeSpecOptionsNumberInput<Spec.SpecMage>({
	fieldName: 'focusMagicPercentUptime',
	label: 'Focus Magic Percent Uptime',
	labelTooltip: 'Percent of uptime for Focus Magic Buddy',
});

export const MageRotationConfig = {
	inputs: [
		{
			type: 'enum' as const,
			label: 'Spec',
			labelTooltip: 'Switches between spec rotation settings. Will also update talents to defaults for the selected spec.',
			values: [
				{ name: 'Arcane', value: RotationType.Arcane },
				{ name: 'Fire', value: RotationType.Fire },
				{ name: 'Frost', value: RotationType.Frost },
			],
			changedEvent: (player: Player<Spec.SpecMage>) => player.rotationChangeEmitter,
			getValue: (player: Player<Spec.SpecMage>) => player.getRotation().type,
			setValue: (eventID: EventID, player: Player<Spec.SpecMage>, newValue: number) => {
				const newRotation = player.getRotation();
				newRotation.type = newValue;

				TypedEvent.freezeAllAndDo(() => {
					if (newRotation.type == RotationType.Arcane) {
						player.setTalentsString(eventID, Presets.ArcaneTalents.data.talentsString);
						if (!newRotation.arcane) {
							newRotation.arcane = ArcaneRotation.clone(Presets.DefaultArcaneRotation.arcane!);
						}
					} else if (newRotation.type == RotationType.Fire) {
						player.setTalentsString(eventID, Presets.FireTalents.data.talentsString);
						if (!newRotation.fire) {
							newRotation.fire = FireRotation.clone(Presets.DefaultFireRotation.fire!);
						}
					} else {
						player.setTalentsString(eventID, Presets.FrostTalents.data.talentsString);
						if (!newRotation.frost) {
							newRotation.frost = FrostRotation.clone(Presets.DefaultFrostRotation.frost!);
						}
					}

					player.setRotation(eventID, newRotation);
				});
			},
		},
		// ********************************************************
		//                        AOE INPUTS
		// ********************************************************
		InputHelpers.makeRotationBooleanInput<Spec.SpecMage>({
			fieldName: 'multiTargetRotation',
			label: 'AOE Rotation',
			labelTooltip: 'Use multi-target spells.',
		}),
		{
			type: 'enum' as const,
			label: 'Primary Spell',
			values: [
				{ name: 'Arcane Explosion', value: AoeRotationSpells.ArcaneExplosion },
				{ name: 'Flamestrike', value: AoeRotationSpells.Flamestrike },
				{ name: 'Blizzard', value: AoeRotationSpells.Blizzard },
			],
			changedEvent: (player: Player<Spec.SpecMage>) => player.rotationChangeEmitter,
			getValue: (player: Player<Spec.SpecMage>) => player.getRotation().aoe?.rotation || 0,
			setValue: (eventID: EventID, player: Player<Spec.SpecMage>, newValue: number) => {
				const newRotation = player.getRotation();
				if (!newRotation.aoe) {
					newRotation.aoe = AoeRotation.create();
				}
				newRotation.aoe.rotation = newValue;
				player.setRotation(eventID, newRotation);
			},
			showWhen: (player: Player<Spec.SpecMage>) => player.getRotation().multiTargetRotation,
		},
		// ********************************************************
		//                       FIRE INPUTS
		// ********************************************************
		{
			type: 'enum' as const,
			label: 'Primary Spell',
			values: [
				{ name: 'Fireball', value: PrimaryFireSpell.Fireball },
				{ name: 'FrostfireBolt', value: PrimaryFireSpell.FrostfireBolt },
			],
			changedEvent: (player: Player<Spec.SpecMage>) => player.rotationChangeEmitter,
			getValue: (player: Player<Spec.SpecMage>) => player.getRotation().fire?.primarySpell || PrimaryFireSpell.Fireball,
			setValue: (eventID: EventID, player: Player<Spec.SpecMage>, newValue: number) => {
				const newRotation = player.getRotation();
				if (!newRotation.fire) {
					newRotation.fire = FireRotation.clone(Presets.DefaultFireRotation.fire!);
				}
				newRotation.fire.primarySpell = newValue;
				player.setRotation(eventID, newRotation);
			},
			showWhen: (player: Player<Spec.SpecMage>) => player.getRotation().type == RotationType.Fire && !player.getRotation().multiTargetRotation,
		},
		{
			type: 'boolean' as const,
			label: 'Maintain Imp. Scorch',
			labelTooltip: 'Always use Scorch when below 5 stacks, or < 5.5s remaining on debuff.',
			changedEvent: (player: Player<Spec.SpecMage>) => player.rotationChangeEmitter,
			getValue: (player: Player<Spec.SpecMage>) => player.getRotation().fire?.maintainImprovedScorch || false,
			setValue: (eventID: EventID, player: Player<Spec.SpecMage>, newValue: boolean) => {
				const newRotation = player.getRotation();
				if (!newRotation.fire) {
					newRotation.fire = FireRotation.clone(Presets.DefaultFireRotation.fire!);
				}
				newRotation.fire.maintainImprovedScorch = newValue;
				player.setRotation(eventID, newRotation);
			},
			showWhen: (player: Player<Spec.SpecMage>) => player.getRotation().type == RotationType.Fire,
		},
		// ********************************************************
		//                       FROST INPUTS
		// ********************************************************
		{
			type: 'number' as const,
			label: 'Water Ele Disobey %',
			labelTooltip: 'Percent of Water Elemental actions which will fail. This represents the Water Elemental moving around or standing still instead of casting.',
			changedEvent: (player: Player<Spec.SpecMage>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
			getValue: (player: Player<Spec.SpecMage>) => (player.getRotation().frost?.waterElementalDisobeyChance || 0) * 100,
			setValue: (eventID: EventID, player: Player<Spec.SpecMage>, newValue: number) => {
				const newRotation = player.getRotation();
				if (!newRotation.frost) {
					newRotation.frost = FrostRotation.clone(Presets.DefaultFrostRotation.frost!);
				}
				newRotation.frost.waterElementalDisobeyChance = newValue / 100;
				player.setRotation(eventID, newRotation);
			},
			showWhen: (player: Player<Spec.SpecMage>) => player.getRotation().type == RotationType.Frost,
			enableWhen: (player: Player<Spec.SpecMage>) => player.getTalents().summonWaterElemental,
		},
		// ********************************************************
		//                      ARCANE INPUTS
		// ********************************************************
		{
			type: 'number' as const,
			label: 'Min ABs before missiles',
			labelTooltip: 'Minimum arcane blasts to cast before using a missile barrage proc',
			changedEvent: (player: Player<Spec.SpecMage>) => player.rotationChangeEmitter,
			getValue: (player: Player<Spec.SpecMage>) => player.getRotation().arcane?.minBlastBeforeMissiles || 0,
			setValue: (eventID: EventID, player: Player<Spec.SpecMage>, newValue: number) => {
				const newRotation = player.getRotation();
				if (!newRotation.arcane) {
					newRotation.arcane = ArcaneRotation.clone(Presets.DefaultArcaneRotation.arcane!);
				}
				newRotation.arcane.minBlastBeforeMissiles = newValue;
				player.setRotation(eventID, newRotation);
			},
			showWhen: (player: Player<Spec.SpecMage>) => player.getRotation().type == RotationType.Arcane && !player.getRotation().multiTargetRotation,
		},
	],
};
