import { BooleanPicker } from '/wotlk/core/components/boolean_picker.js';
import { EnumPicker } from '/wotlk/core/components/enum_picker.js';
import { IconEnumPicker, IconEnumPickerConfig } from '/wotlk/core/components/icon_enum_picker.js';
import { IconPickerConfig } from '/wotlk/core/components/icon_picker.js';
import { Spec } from '/wotlk/core/proto/common.js';
import { WeaponImbue } from '/wotlk/core/proto/common.js';
import { ActionId } from '/wotlk/core/proto_utils/action_id.js';
import { Player } from '/wotlk/core/player.js';
import { Sim } from '/wotlk/core/sim.js';
import { IndividualSimUI } from '/wotlk/core/individual_sim_ui.js';
import { Target } from '/wotlk/core/target.js';
import { EventID, TypedEvent } from '/wotlk/core/typed_event.js';

import * as InputHelpers from '/wotlk/core/components/input_helpers.js';

import {
	Rogue,
	Rogue_Rotation as RogueRotation,
	Rogue_Rotation_Builder as Builder,
	Rogue_Options as RogueOptions,
} from '/wotlk/core/proto/rogue.js';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const RogueRotationConfig = {
	inputs: [
		InputHelpers.makeRotationEnumInput<Spec.SpecRogue, Builder>({
			fieldName: 'builder',
			label: 'Builder',
			values: [
				{
					name: 'Auto', value: Builder.Auto,
					tooltip: 'Automatically selects a builder based on weapons/talents.',
				},
				{ name: 'Sinister Strike', value: Builder.SinisterStrike },
				{ name: 'Mutilate', value: Builder.Mutilate },
			],
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecRogue>({
			fieldName: 'maintainExposeArmor',
			label: 'Maintain EA',
			labelTooltip: 'Keeps Expose Armor active on the primary target.',
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecRogue>({
			fieldName: 'useRupture',
			label: 'Use Rupture',
			labelTooltip: 'Uses Rupture over Eviscerate when appropriate.',
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecRogue>({
			fieldName: 'useShiv',
			label: 'Use Shiv',
			labelTooltip: 'Uses Shiv in place of the selected builder if Deadly Poison is about to expire. Requires Deadly Poison in the off-hand.',
		}),
		InputHelpers.makeRotationNumberInput<Spec.SpecRogue>({
			fieldName: 'minComboPointsForDamageFinisher',
			label: 'Min CPs for Damage Finisher',
			labelTooltip: 'Will not use Eviscerate or Rupture unless the Rogue has at least this many Combo Points.',
		}),
	],
};
