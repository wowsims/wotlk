import { IconPickerConfig } from '../core/components/icon_picker.js';
import { RaidTarget } from '../core/proto/common.js';
import { Spec } from '../core/proto/common.js';
import { NO_TARGET } from '../core/proto_utils/utils.js';
import { ActionId } from '../core/proto_utils/action_id.js';
import { Player } from '../core/player.js';
import { Sim } from '../core/sim.js';
import { EventID, TypedEvent } from '../core/typed_event.js';
import { IndividualSimUI } from '../core/individual_sim_ui.js';
import { Target } from '../core/target.js';
import { EnumPicker } from '../core/components/enum_picker.js';
import { IconEnumPicker, IconEnumPickerConfig } from '../core/components/icon_enum_picker.js';
import { CustomRotationPickerConfig } from '../core/components/custom_rotation_picker.js';
import { CustomRotation } from '../core/proto/common.js';

import {
	WarriorShout,
	WarriorTalents as WarriorTalents,
	ProtectionWarrior,
	ProtectionWarrior_Rotation as ProtectionWarriorRotation,
	ProtectionWarrior_Rotation_DemoShoutChoice as DemoShoutChoice,
	ProtectionWarrior_Rotation_ThunderClapChoice as ThunderClapChoice,
	ProtectionWarrior_Options as ProtectionWarriorOptions,
	ProtectionWarrior_Rotation_SpellOption as SpellOption,
} from '../core/proto/warrior.js';

import * as InputHelpers from '../core/components/input_helpers.js';
import * as Presets from './presets.js';
import { SimUI } from '../core/sim_ui.js';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const StartingRage = InputHelpers.makeSpecOptionsNumberInput<Spec.SpecProtectionWarrior>({
	fieldName: 'startingRage',
	label: 'Starting Rage',
	labelTooltip: 'Initial rage at the start of each iteration.',
});

export const ShoutPicker = InputHelpers.makeSpecOptionsEnumIconInput<Spec.SpecProtectionWarrior, WarriorShout>({
	fieldName: 'shout',
	values: [
		{ color: 'c79c6e', value: WarriorShout.WarriorShoutNone },
		{ actionId: ActionId.fromSpellId(2048), value: WarriorShout.WarriorShoutBattle },
		{ actionId: ActionId.fromSpellId(469), value: WarriorShout.WarriorShoutCommanding },
	],
});

export const PrecastShout = InputHelpers.makeSpecOptionsBooleanInput<Spec.SpecProtectionWarrior>({
	fieldName: 'precastShout',
	label: 'Precast Shout',
	labelTooltip: 'Selected shout is cast 10 seconds before combat starts.',
});

export const ProtectionWarriorRotationConfig = {
	inputs: [
		InputHelpers.makeCustomRotationInput<Spec.SpecProtectionWarrior, SpellOption>({
			fieldName: 'customRotation',
			numColumns: 3,
			values: [
				{ actionId: ActionId.fromSpellId(30357), value: SpellOption.Revenge },
				{ actionId: ActionId.fromSpellId(47488), value: SpellOption.ShieldSlam },
				{ actionId: ActionId.fromSpellId(47440), value: SpellOption.Shout},
				{ actionId: ActionId.fromSpellId(47502), value: SpellOption.ThunderClap },
				{ actionId: ActionId.fromSpellId(25203), value: SpellOption.DemoralizingShout },
				{ actionId: ActionId.fromSpellId(71552), value: SpellOption.MortalStrike },
				{ actionId: ActionId.fromSpellId(30022), value: SpellOption.Devastate },
				{ actionId: ActionId.fromSpellId(47467), value: SpellOption.SunderArmor},
				{ actionId: ActionId.fromSpellId(12809), value: SpellOption.ConcussionBlow },
				{ actionId: ActionId.fromSpellId(46968), value: SpellOption.Shockwave },
			],
		}),
		InputHelpers.makeRotationNumberInput<Spec.SpecProtectionWarrior>({
			fieldName: 'hsRageThreshold',
			label: 'HS rage threshold',
			labelTooltip: 'Heroic Strike when rage is above:',
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecProtectionWarrior>({
			fieldName: 'prioSslamOnShieldBlock',
			label: 'Prio SSlam on Shield Block',
			labelTooltip: 'The rotation code will prio SSlam over Revenge during active shield block windows.',
		}),
		InputHelpers.makeRotationEnumInput<Spec.SpecProtectionWarrior, DemoShoutChoice>({
			fieldName: 'demoShoutChoice',
			label: 'Demo Shout',
			values: [
				{ name: 'Never', value: DemoShoutChoice.DemoShoutChoiceNone },
				{ name: 'Maintain Debuff', value: DemoShoutChoice.DemoShoutChoiceMaintain },
				{ name: 'Filler', value: DemoShoutChoice.DemoShoutChoiceFiller },
			],
		}),
		InputHelpers.makeRotationEnumInput<Spec.SpecProtectionWarrior, ThunderClapChoice>({
			fieldName: 'thunderClapChoice',
			label: 'Thunder Clap',
			values: [
				{ name: 'Never', value: ThunderClapChoice.ThunderClapChoiceNone },
				{ name: 'Maintain Debuff', value: ThunderClapChoice.ThunderClapChoiceMaintain },
				{ name: 'On CD', value: ThunderClapChoice.ThunderClapChoiceOnCD },
			],
		}),
	],
};
