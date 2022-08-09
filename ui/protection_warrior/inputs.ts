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

import {
	WarriorShout,
	WarriorTalents as WarriorTalents,
	ProtectionWarrior,
	ProtectionWarrior_Rotation as ProtectionWarriorRotation,
	ProtectionWarrior_Rotation_DemoShout as DemoShout,
	ProtectionWarrior_Rotation_ThunderClap as ThunderClap,
	ProtectionWarrior_Options as ProtectionWarriorOptions
} from '../core/proto/warrior.js';

import * as InputHelpers from '../core/components/input_helpers.js';
import * as Presets from './presets.js';
import { SimUI } from '../core/sim_ui.js';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const ShieldWall = InputHelpers.makeSpecOptionsBooleanIconInput<Spec.SpecProtectionWarrior>({
	fieldName: 'useShieldWall',
	id: ActionId.fromSpellId(871),
});

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
		InputHelpers.makeRotationBooleanInput<Spec.SpecProtectionWarrior>({
			fieldName: 'useCleave',
			label: 'Use Cleave',
			labelTooltip: 'Use Cleave instead of Heroic Strike.',
		}),
		InputHelpers.makeRotationNumberInput<Spec.SpecProtectionWarrior>({
			fieldName: 'hsRageThreshold',
			label: 'HS rage threshold',
			labelTooltip: 'Heroic Strike when rage is above:',
		}),
		InputHelpers.makeRotationEnumInput<Spec.SpecProtectionWarrior, DemoShout>({
			fieldName: 'demoShout',
			label: 'Demo Shout',
			values: [
				{ name: 'Never', value: DemoShout.DemoShoutNone },
				{ name: 'Maintain Debuff', value: DemoShout.DemoShoutMaintain },
				{ name: 'Filler', value: DemoShout.DemoShoutFiller },
			],
		}),
		InputHelpers.makeRotationEnumInput<Spec.SpecProtectionWarrior, ThunderClap>({
			fieldName: 'thunderClap',
			label: 'Thunder Clap',
			values: [
				{ name: 'Never', value: ThunderClap.ThunderClapNone },
				{ name: 'Maintain Debuff', value: ThunderClap.ThunderClapMaintain },
				{ name: 'On CD', value: ThunderClap.ThunderClapOnCD },
			],
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecProtectionWarrior>({
			fieldName: 'useShieldBlock',
			label: 'Use Shield Block',
		}),
	],
};
