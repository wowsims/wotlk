import { IconPickerConfig } from '../core/components/icon_picker.js';
import { RaidTarget } from '../core/proto/common.js';
import { Spec } from '../core/proto/common.js';
import { NO_TARGET } from '../core/proto_utils/utils.js';
import { ActionId } from '../core/proto_utils/action_id.js';
import { Player } from '../core/player.js';
import { Sim } from '../core/sim.js';
import { EventID, TypedEvent } from '../core/typed_event.js';
import { Target } from '../core/target.js';

import {
	WarriorShout,
	WarriorTalents as WarriorTalents,
	Warrior,
	Warrior_Rotation as WarriorRotation,
	Warrior_Rotation_SunderArmor as SunderArmor,
	Warrior_Options as WarriorOptions,
} from '../core/proto/warrior.js';

import * as InputHelpers from '../core/components/input_helpers.js';
import * as Presets from './presets.js';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const Recklessness = InputHelpers.makeSpecOptionsBooleanIconInput<Spec.SpecWarrior>({
	fieldName: 'useRecklessness',
	id: ActionId.fromSpellId(1719),
});

export const StartingRage = InputHelpers.makeSpecOptionsNumberInput<Spec.SpecWarrior>({
	fieldName: 'startingRage',
	label: 'Starting Rage',
	labelTooltip: 'Initial rage at the start of each iteration.',
});

export const ShoutPicker = InputHelpers.makeSpecOptionsEnumIconInput<Spec.SpecWarrior, WarriorShout>({
	fieldName: 'shout',
	values: [
		{ color: 'c79c6e', value: WarriorShout.WarriorShoutNone },
		{ actionId: ActionId.fromSpellId(2048), value: WarriorShout.WarriorShoutBattle },
		{ actionId: ActionId.fromSpellId(469), value: WarriorShout.WarriorShoutCommanding },
	],
});

export const WarriorRotationConfig = {
	inputs: [
		InputHelpers.makeRotationBooleanInput<Spec.SpecWarrior>({
			fieldName: 'useMs',
			label: 'Use Mortal Strike',
			labelTooltip: 'Use Mortal Strike when rage threshold is met.',
			changeEmitter: (player: Player<Spec.SpecWarrior>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
			showWhen: (player: Player<Spec.SpecWarrior>) => player.getTalents().mortalStrike,
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecWarrior>({
			fieldName: 'useCleave',
			label: 'Use Cleave',
			labelTooltip: 'Use Cleave instead of Heroic Strike.',
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecWarrior>({
			fieldName: 'useRend',
			label: 'Use Rend',
			labelTooltip: 'Use Rend when rage threshold is met and the debuff duration is less than refresh time.',
			changeEmitter: (player: Player<Spec.SpecWarrior>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecWarrior>({
			fieldName: 'prioritizeWw',
			label: 'Prioritize WW',
			labelTooltip: 'Prioritize Whirlwind over Bloodthirst.',
			changeEmitter: (player: Player<Spec.SpecWarrior>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
			showWhen: (player: Player<Spec.SpecWarrior>) => player.getTalents().bloodthirst,
		}),
		InputHelpers.makeRotationNumberInput<Spec.SpecWarrior>({
			fieldName: 'msRageThreshold',
			label: 'Mortal Strike rage threshold',
			labelTooltip: 'Mortal Strike when rage is above:',
			changeEmitter: (player: Player<Spec.SpecWarrior>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
			showWhen: (player: Player<Spec.SpecWarrior>) => player.getTalents().mortalStrike,
		}),
		InputHelpers.makeRotationNumberInput<Spec.SpecWarrior>({
			fieldName: 'slamRageThreshold',
			label: 'Slam rage threshold',
			labelTooltip: 'Slam when rage is above:',
			changeEmitter: (player: Player<Spec.SpecWarrior>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
			showWhen: (player: Player<Spec.SpecWarrior>) => player.getTalents().mortalStrike,
		}),
		InputHelpers.makeRotationNumberInput<Spec.SpecWarrior>({
			fieldName: 'hsRageThreshold',
			label: 'HS rage threshold',
			labelTooltip: 'Heroic Strike when rage is above:',
		}),
		InputHelpers.makeRotationNumberInput<Spec.SpecWarrior>({
			fieldName: 'rendRageThresholdBelow',
			label: 'Rend rage threshold below',
			labelTooltip: 'Rend will only be used when rage is smaller than this value.',
			changeEmitter: (player: Player<Spec.SpecWarrior>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
			showWhen: (player: Player<Spec.SpecWarrior>) => player.getRotation().useRend == true && player.getTalents().bloodthirst,
		}),
		InputHelpers.makeRotationNumberInput<Spec.SpecWarrior>({
			fieldName: 'rendCdThreshold',
			label: 'Rend Refresh Time',
			labelTooltip: 'Refresh Rend when the remaining duration is less than this amount of time (seconds).',
			changeEmitter: (player: Player<Spec.SpecWarrior>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
			showWhen: (player: Player<Spec.SpecWarrior>) => player.getRotation().useRend == true,
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecWarrior>({
			fieldName: 'useHsDuringExecute',
			label: 'HS during Execute Phase',
			labelTooltip: 'Use Heroic Strike during Execute Phase.',
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecWarrior>({
			fieldName: 'useBtDuringExecute',
			label: 'BT during Execute Phase',
			labelTooltip: 'Use Bloodthirst during Execute Phase.',
			changeEmitter: (player: Player<Spec.SpecWarrior>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
			showWhen: (player: Player<Spec.SpecWarrior>) => player.getTalents().bloodthirst,
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecWarrior>({
			fieldName: 'useWwDuringExecute',
			label: 'WW during Execute Phase',
			labelTooltip: 'Use Whirlwind during Execute Phase.',
			changeEmitter: (player: Player<Spec.SpecWarrior>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
			showWhen: (player: Player<Spec.SpecWarrior>) => player.getTalents().bloodthirst,
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecWarrior>({
			fieldName: 'useSlamOverExecute',
			label: 'Slam Over Execute',
			labelTooltip: 'Use Slam Over Execute when Bloodsurge Procs in Execute Phase.',
			changeEmitter: (player: Player<Spec.SpecWarrior>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
			showWhen: (player: Player<Spec.SpecWarrior>) => player.getTalents().bloodthirst,
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecWarrior>({
			fieldName: 'spamExecute',
			label: 'Spam Execute',
			labelTooltip: 'Use Execute whenever possible during Execute Phase',
			changeEmitter: (player: Player<Spec.SpecWarrior>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
			showWhen: (player: Player<Spec.SpecWarrior>) => player.getTalents().mortalStrike,
		}),
		InputHelpers.makeRotationEnumInput<Spec.SpecWarrior, SunderArmor>({
			fieldName: 'sunderArmor',
			label: 'Sunder Armor',
			values: [
				{ name: 'Never', value: SunderArmor.SunderArmorNone },
				{ name: 'Help Stack', value: SunderArmor.SunderArmorHelpStack },
				{ name: 'Maintain Debuff', value: SunderArmor.SunderArmorMaintain },
			],
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecWarrior>({
			fieldName: 'maintainDemoShout',
			label: 'Maintain Demo Shout',
			labelTooltip: 'Keep Demo Shout active on the primary target.',
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecWarrior>({
			fieldName: 'maintainThunderClap',
			label: 'Maintain Thunder Clap',
			labelTooltip: 'Keep Thunder Clap active on the primary target.',
		}),
	],
};
