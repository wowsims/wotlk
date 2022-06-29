import { BooleanPicker } from '/tbc/core/components/boolean_picker.js';
import { EnumPicker } from '/tbc/core/components/enum_picker.js';
import { IconEnumPicker, IconEnumPickerConfig } from '/tbc/core/components/icon_enum_picker.js';
import { IconPickerConfig } from '/tbc/core/components/icon_picker.js';
import { Spec } from '/tbc/core/proto/common.js';
import { WeaponImbue } from '/tbc/core/proto/common.js';
import { ActionId } from '/tbc/core/proto_utils/action_id.js';
import { Player } from '/tbc/core/player.js';
import { Sim } from '/tbc/core/sim.js';
import { IndividualSimUI } from '/tbc/core/individual_sim_ui.js';
import { Target } from '/tbc/core/target.js';
import { EventID, TypedEvent } from '/tbc/core/typed_event.js';

import {
	Rogue,
	Rogue_Rotation as RogueRotation,
	Rogue_Rotation_Builder as Builder,
	Rogue_Options as RogueOptions,
} from '/tbc/core/proto/rogue.js';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const RogueRotationConfig = {
	inputs: [
		{
			type: 'enum' as const, cssClass: 'builder-picker',
			getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
			config: {
				label: 'Builder',
				values: [
					{
						name: 'Auto', value: Builder.Auto,
						tooltip: 'Automatically selects a builder based on weapons/talents.',
					},
					{ name: 'Sinister Strike', value: Builder.SinisterStrike },
					{ name: 'Backstab', value: Builder.Backstab },
					{ name: 'Hemorrhage', value: Builder.Hemorrhage },
					{ name: 'Mutilate', value: Builder.Mutilate },
				],
				changedEvent: (player: Player<Spec.SpecRogue>) => player.rotationChangeEmitter,
				getValue: (player: Player<Spec.SpecRogue>) => player.getRotation().builder,
				setValue: (eventID: EventID, player: Player<Spec.SpecRogue>, newValue: number) => {
					const newRotation = player.getRotation();
					newRotation.builder = newValue;
					player.setRotation(eventID, newRotation);
				},
			},
		},
		{
			type: 'boolean' as const, cssClass: 'maintain-expose-armor-picker',
			getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
			config: {
				label: 'Maintain EA',
				labelTooltip: 'Keeps Expose Armor active on the primary target.',
				changedEvent: (player: Player<Spec.SpecRogue>) => player.rotationChangeEmitter,
				getValue: (player: Player<Spec.SpecRogue>) => player.getRotation().maintainExposeArmor,
				setValue: (eventID: EventID, player: Player<Spec.SpecRogue>, newValue: boolean) => {
					const newRotation = player.getRotation();
					newRotation.maintainExposeArmor = newValue;
					player.setRotation(eventID, newRotation);
				},
			},
		},
		{
			type: 'boolean' as const, cssClass: 'use-rupture-picker',
			getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
			config: {
				label: 'Use Rupture',
				labelTooltip: 'Uses Rupture over Eviscerate when appropriate.',
				changedEvent: (player: Player<Spec.SpecRogue>) => player.rotationChangeEmitter,
				getValue: (player: Player<Spec.SpecRogue>) => player.getRotation().useRupture,
				setValue: (eventID: EventID, player: Player<Spec.SpecRogue>, newValue: boolean) => {
					const newRotation = player.getRotation();
					newRotation.useRupture = newValue;
					player.setRotation(eventID, newRotation);
				},
			},
		},
		{
			type: 'boolean' as const, cssClass: 'use-shiv-picker',
			getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
			config: {
				label: 'Use Shiv',
				labelTooltip: 'Uses Shiv in place of the selected builder if Deadly Poison is about to expire. Requires Deadly Poison in the off-hand.',
				changedEvent: (player: Player<Spec.SpecRogue>) => TypedEvent.onAny([player.rotationChangeEmitter, player.consumesChangeEmitter]),
				getValue: (player: Player<Spec.SpecRogue>) => player.getRotation().useShiv,
				setValue: (eventID: EventID, player: Player<Spec.SpecRogue>, newValue: boolean) => {
					const newRotation = player.getRotation();
					newRotation.useShiv = newValue;
					player.setRotation(eventID, newRotation);
				},
				enableWhen: (player: Player<Spec.SpecRogue>) => player.getConsumes().offHandImbue == WeaponImbue.WeaponImbueRogueDeadlyPoison,
			},
		},
		{
			type: 'number' as const, cssClass: 'min-combo-points-for-dps-finisher-picker',
			getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
			config: {
				label: 'Min CPs for Damage Finisher',
				labelTooltip: 'Will not use Eviscerate or Rupture unless the Rogue has at least this many Combo Points.',
				changedEvent: (player: Player<Spec.SpecRogue>) => player.rotationChangeEmitter,
				getValue: (player: Player<Spec.SpecRogue>) => player.getRotation().minComboPointsForDamageFinisher,
				setValue: (eventID: EventID, player: Player<Spec.SpecRogue>, newValue: number) => {
					const newRotation = player.getRotation();
					newRotation.minComboPointsForDamageFinisher = newValue;
					player.setRotation(eventID, newRotation);
				},
			},
		},
	],
};

function makeBooleanRogueBuffInput(id: ActionId, optionsFieldName: keyof RogueOptions): IconPickerConfig<Player<any>, boolean> {
	return {
		id: id,
		states: 2,
		changedEvent: (player: Player<Spec.SpecRogue>) => player.specOptionsChangeEmitter,
		getValue: (player: Player<Spec.SpecRogue>) => player.getSpecOptions()[optionsFieldName] as boolean,
		setValue: (eventID: EventID, player: Player<Spec.SpecRogue>, newValue: boolean) => {
			const newOptions = player.getSpecOptions();
			(newOptions[optionsFieldName] as boolean) = newValue;
			player.setSpecOptions(eventID, newOptions);
		},
	};
}
