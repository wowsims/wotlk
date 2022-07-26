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

export const PrecastShout = InputHelpers.makeSpecOptionsBooleanInput<Spec.SpecWarrior>({
	fieldName: 'precastShout',
	label: 'Precast Shout',
	labelTooltip: 'Selected shout is cast 10 seconds before combat starts.',
});

export const WarriorRotationConfig = {
	inputs: [
		InputHelpers.makeRotationBooleanInput<Spec.SpecWarrior>({
			fieldName: 'useCleave',
			label: 'Use Cleave',
			labelTooltip: 'Use Cleave instead of Heroic Strike.',
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecWarrior>({
			fieldName: 'useOverpower',
			label: 'Use Overpower',
			labelTooltip: 'Use Overpower when available.',
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecWarrior>({
			fieldName: 'useHamstring',
			label: 'Use Hamstring',
			labelTooltip: 'Use Hamstring on free globals.',
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecWarrior>({
			fieldName: 'useSlam',
			label: 'Use Slam',
			labelTooltip: 'Use Slam whenever possible.',
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecWarrior>({
			fieldName: 'prioritizeWw',
			label: 'Prioritize WW',
			labelTooltip: 'Prioritize Whirlwind over Bloodthirst or Mortal Strike.',
		}),
		InputHelpers.makeRotationNumberInput<Spec.SpecWarrior>({
			fieldName: 'hsRageThreshold',
			label: 'HS rage threshold',
			labelTooltip: 'Heroic Strike when rage is above:',
		}),
		InputHelpers.makeRotationNumberInput<Spec.SpecWarrior>({
			fieldName: 'overpowerRageThreshold',
			label: 'Overpower rage threshold',
			labelTooltip: 'Use Overpower when rage is below a point.',
			showWhen: (player: Player<Spec.SpecWarrior>) => player.getRotation().useOverpower,
		}),
		InputHelpers.makeRotationNumberInput<Spec.SpecWarrior>({
			fieldName: 'hamstringRageThreshold',
			label: 'Hamstring rage threshold',
			labelTooltip: 'Hamstring will only be used when rage is larger than this value.',
			showWhen: (player: Player<Spec.SpecWarrior>) => player.getRotation().useHamstring,
		}),
		InputHelpers.makeRotationNumberInput<Spec.SpecWarrior>({
			fieldName: 'slamLatency',
			label: 'Slam Latency',
			labelTooltip: 'Time between MH swing and start of the Slam cast, in milliseconds.',
			showWhen: (player: Player<Spec.SpecWarrior>) => player.getRotation().useSlam && player.getTalents().improvedSlam == 2,
			changeEmitter: (player: Player<Spec.SpecWarrior>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
		}),
		InputHelpers.makeRotationNumberInput<Spec.SpecWarrior>({
			fieldName: 'slamGcdDelay',
			label: 'Slam GCD Delay',
			labelTooltip: 'Amount of time Slam may delay the GCD, in milliseconds.',
			changeEmitter: (player: Player<Spec.SpecWarrior>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
			showWhen: (player: Player<Spec.SpecWarrior>) => player.getRotation().useSlam && player.getTalents().improvedSlam == 2,
			extraCssClasses: [
				'experimental',
			],
		}),
		InputHelpers.makeRotationNumberInput<Spec.SpecWarrior>({
			fieldName: 'slamMsWwDelay',
			label: 'Slam MS+WW Delay',
			labelTooltip: 'Amount of time Slam may delay MS+WW, in milliseconds.',
			changeEmitter: (player: Player<Spec.SpecWarrior>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
			showWhen: (player: Player<Spec.SpecWarrior>) => player.getRotation().useSlam && player.getTalents().improvedSlam == 2,
			extraCssClasses: [
				'experimental',
			],
		}),
		InputHelpers.makeRotationNumberInput<Spec.SpecWarrior>({
			fieldName: 'rampageCdThreshold',
			label: 'Rampage Refresh Time',
			labelTooltip: 'Refresh Rampage when the remaining duration is less than this amount of time (seconds).',
			changeEmitter: (player: Player<Spec.SpecWarrior>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
			showWhen: (player: Player<Spec.SpecWarrior>) => player.getTalents().rampage,
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
			fieldName: 'useMsDuringExecute',
			label: 'MS during Execute Phase',
			labelTooltip: 'Use Mortal Strike during Execute Phase.',
			changeEmitter: (player: Player<Spec.SpecWarrior>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
			showWhen: (player: Player<Spec.SpecWarrior>) => player.getTalents().mortalStrike,
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecWarrior>({
			fieldName: 'useWwDuringExecute',
			label: 'WW during Execute Phase',
			labelTooltip: 'Use Whirlwind during Execute Phase.',
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecWarrior>({
			fieldName: 'useSlamDuringExecute',
			label: 'Slam during Execute Phase',
			labelTooltip: 'Use Slam during Execute Phase.',
			changeEmitter: (player: Player<Spec.SpecWarrior>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
			showWhen: (player: Player<Spec.SpecWarrior>) => player.getRotation().useSlam && player.getTalents().improvedSlam == 2,
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
