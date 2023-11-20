import { Player } from '../core/player.js';
import { Spec } from '../core/proto/common.js';
import { ActionId } from '../core/proto_utils/action_id.js';
import { TypedEvent } from '../core/typed_event.js';

import {
	Warrior_Rotation_MainGcd as MainGcd,
	Warrior_Rotation_SpellOption as SpellOption,
	Warrior_Rotation_StanceOption as StanceOption,
	Warrior_Rotation_SunderArmor as SunderArmor,
	WarriorShout
} from '../core/proto/warrior.js';

import * as InputHelpers from '../core/components/input_helpers.js';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const Recklessness = InputHelpers.makeSpecOptionsBooleanIconInput<Spec.SpecWarrior>({
	fieldName: 'useRecklessness',
	id: ActionId.fromSpellId(1719),
});

export const ShatteringThrow = InputHelpers.makeSpecOptionsBooleanIconInput<Spec.SpecWarrior>({
	fieldName: 'useShatteringThrow',
	id: ActionId.fromSpellId(64382),
});

export const StartingRage = InputHelpers.makeSpecOptionsNumberInput<Spec.SpecWarrior>({
	fieldName: 'startingRage',
	label: 'Starting Rage',
	labelTooltip: 'Initial rage at the start of each iteration.',
});


export const StanceSnapshot = InputHelpers.makeSpecOptionsBooleanInput<Spec.SpecWarrior>({
	fieldName: 'stanceSnapshot',
	label: 'Stance Snapshot',
	labelTooltip: 'Ability that is cast at the same time as stance swap will benefit from the bonus of the stance before the swap.',
});

// Allows for auto gemming whilst ignoring expertise cap
// (Useful for Arms)
export const DisableExpertiseGemming = InputHelpers.makeSpecOptionsBooleanInput<Spec.SpecWarrior>({
	fieldName: 'disableExpertiseGemming',
	label: 'Disable expertise gemming',
	labelTooltip: 'Disables auto gemming for expertise',
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
			showWhen: (player: Player<Spec.SpecWarrior>) => player.getTalentTree() == 0 && !player.getRotation().customRotationOption,
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecWarrior>({
			fieldName: 'useCleave',
			label: 'Use Cleave',
			labelTooltip: 'Use Cleave instead of Heroic Strike.',
			showWhen: (player: Player<Spec.SpecWarrior>) => !player.getRotation().customRotationOption,
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecWarrior>({
			fieldName: 'useRend',
			label: 'Use Rend',
			labelTooltip: 'Use Rend when rage threshold is met and the debuff duration is less than refresh time.',
			changeEmitter: (player: Player<Spec.SpecWarrior>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
			showWhen: (player: Player<Spec.SpecWarrior>) => !player.getRotation().customRotationOption,
		}),

		InputHelpers.makeRotationBooleanInput<Spec.SpecWarrior>({
			fieldName: 'useOverpower',
			label: 'Use Overpower',
			labelTooltip: 'Use Overpower whenever it is available on an open GCD.',
			changeEmitter: (player: Player<Spec.SpecWarrior>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
			showWhen: (player: Player<Spec.SpecWarrior>) => player.getTalentTree() == 1 && !player.getRotation().customRotationOption,
		}),

		InputHelpers.makeRotationEnumInput<Spec.SpecWarrior>({
			fieldName: 'mainGcd',
			label: 'Main GCD',
			labelTooltip: 'Main GCD ability that will be prioritized above other abilities.',
			values: [
				{ name: 'Slam', value: MainGcd.Slam },
				{ name: 'Bloodthirst', value: MainGcd.Bloodthirst },
				{ name: 'Whirlwind', value: MainGcd.Whirlwind },
			],
			changeEmitter: (player: Player<Spec.SpecWarrior>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
			showWhen: (player: Player<Spec.SpecWarrior>) => player.getTalentTree() == 1 && !player.getRotation().customRotationOption,
		}),
		InputHelpers.makeRotationNumberInput<Spec.SpecWarrior>({
			fieldName: 'msRageThreshold',
			label: 'Mortal Strike rage threshold',
			labelTooltip: 'Mortal Strike when rage is above:',
			changeEmitter: (player: Player<Spec.SpecWarrior>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
			showWhen: (player: Player<Spec.SpecWarrior>) => (player.getRotation().useMs || player.getRotation().customRotationOption) && player.getTalentTree() == 0,
		}),
		InputHelpers.makeRotationNumberInput<Spec.SpecWarrior>({
			fieldName: 'slamRageThreshold',
			label: 'Slam rage threshold',
			labelTooltip: 'Slam when rage is above:',
			changeEmitter: (player: Player<Spec.SpecWarrior>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
			showWhen: (player: Player<Spec.SpecWarrior>) => (player.getRotation().useMs || player.getRotation().customRotationOption) && player.getTalentTree() == 0,
		}),
		InputHelpers.makeRotationNumberInput<Spec.SpecWarrior>({
			fieldName: 'bloodsurgeDurationThreshold',
			label: 'Exp Slam: Bloodsurge duration threshold (s)',
			labelTooltip: 'Cast Exp Slam when Bloodsurge duration is below (seconds):',
			changeEmitter: (player: Player<Spec.SpecWarrior>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
			showWhen: (player: Player<Spec.SpecWarrior>) => player.getRotation().customRotationOption && player.getTalentTree() == 1,
		}),
		InputHelpers.makeRotationNumberInput<Spec.SpecWarrior>({
			fieldName: 'hsRageThreshold',
			label: 'HS rage threshold',
			labelTooltip: 'Heroic Strike when rage is above:',
		}),
		InputHelpers.makeRotationNumberInput<Spec.SpecWarrior>({
			fieldName: 'rendHealthThresholdAbove',
			label: 'Rend health threshold (%)',
			labelTooltip: 'Rend will only be used when boss health is above this value in %.',
			changeEmitter: (player: Player<Spec.SpecWarrior>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
			showWhen: (player: Player<Spec.SpecWarrior>) => (player.getRotation().useRend == true || player.getRotation().customRotationOption),
		}),
		InputHelpers.makeRotationNumberInput<Spec.SpecWarrior>({
			fieldName: 'rendRageThresholdBelow',
			label: 'Rend rage threshold below',
			labelTooltip: 'Rend will only be used when rage is smaller than this value.',
			changeEmitter: (player: Player<Spec.SpecWarrior>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
			showWhen: (player: Player<Spec.SpecWarrior>) => (player.getRotation().useRend == true || player.getRotation().customRotationOption) && player.getTalentTree() == 1,
		}),
		InputHelpers.makeRotationNumberInput<Spec.SpecWarrior>({
			fieldName: 'rendCdThreshold',
			label: 'Rend Refresh Time (s)',
			labelTooltip: 'Refresh Rend when the remaining duration is less than this amount of time (seconds).',
			changeEmitter: (player: Player<Spec.SpecWarrior>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
			showWhen: (player: Player<Spec.SpecWarrior>) => player.getRotation().useRend == true || player.getRotation().customRotationOption,
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
			showWhen: (player: Player<Spec.SpecWarrior>) => player.getTalentTree() == 1,
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecWarrior>({
			fieldName: 'useWwDuringExecute',
			label: 'WW during Execute Phase',
			labelTooltip: 'Use Whirlwind during Execute Phase.',
			changeEmitter: (player: Player<Spec.SpecWarrior>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
			showWhen: (player: Player<Spec.SpecWarrior>) => player.getTalentTree() == 1,
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecWarrior>({
			fieldName: 'useSlamOverExecute',
			label: 'Slam during Execute Phase',
			labelTooltip: 'Use Slam Over Execute when Bloodsurge Procs in Execute Phase.',
			changeEmitter: (player: Player<Spec.SpecWarrior>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
			showWhen: (player: Player<Spec.SpecWarrior>) => player.getTalentTree() == 1,
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecWarrior>({
			fieldName: 'executePhaseOverpower',
			label: 'Overpower in Execute Phase',
			labelTooltip: 'Use Overpower instead of Execute whenever it is available.',
			changeEmitter: (player: Player<Spec.SpecWarrior>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
			showWhen: (player: Player<Spec.SpecWarrior>) => (player.getRotation().useOverpower == true || player.getRotation().customRotationOption) && player.getTalentTree() == 1,
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecWarrior>({
			fieldName: 'spamExecute',
			label: 'Spam Execute',
			labelTooltip: 'Use Execute whenever possible during Execute Phase.',
			changeEmitter: (player: Player<Spec.SpecWarrior>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
			showWhen: (player: Player<Spec.SpecWarrior>) => player.getTalentTree() == 0,
		}),
		InputHelpers.makeRotationEnumInput<Spec.SpecWarrior>({
			fieldName: 'sunderArmor',
			label: 'Sunder Armor',
			values: [
				{ name: 'Never', value: SunderArmor.SunderArmorNone },
				{ name: 'Help Stack', value: SunderArmor.SunderArmorHelpStack },
				{ name: 'Maintain Debuff', value: SunderArmor.SunderArmorMaintain },
			],
		}),
		InputHelpers.makeRotationEnumInput<Spec.SpecWarrior>({
			fieldName: 'stanceOption',
			label: 'Stance Option',
			labelTooltip: 'Stance to stay on. The default for Fury (Bloodthirst) is Berserker Stance and Battle Stance for everything else.',
			values: [
				{ name: 'Default', value: StanceOption.DefaultStance },
				{ name: 'Battle Stance', value: StanceOption.BattleStance },
				{ name: 'Berserker Stance', value: StanceOption.BerserkerStance },
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

		InputHelpers.makeRotationBooleanInput<Spec.SpecWarrior>({
			fieldName: 'customRotationOption',
			label: 'Custom Rotation (Advanced)',
			labelTooltip: 'Create your own rotation action priority list.',
			showWhen: (player: Player<Spec.SpecWarrior>) => player.sim.getShowExperimental(),
			changeEmitter: (player: Player<Spec.SpecWarrior>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
		}),

		InputHelpers.makeCustomRotationInput<Spec.SpecWarrior, SpellOption>({
			fieldName: 'customRotation',
			numColumns: 3,
			showWhen: (player: Player<Spec.SpecWarrior>) => player.getRotation().customRotationOption && player.sim.getShowExperimental(),
			values: [
				{ actionId: ActionId.fromSpellId(23881), value: SpellOption.BloodthirstCustom },
				{ actionId: ActionId.fromSpellId(1680), value: SpellOption.WhirlwindCustom },
				{ actionId: ActionId.fromSpellId(47475), value: SpellOption.SlamCustom },
				{ actionId: ActionId.fromSpellId(47475), value: SpellOption.SlamExpiring, text: "Exp", showWhen: (player: Player<Spec.SpecWarrior>) => player.getTalentTree() == 1, },
				{ actionId: ActionId.fromSpellId(47486), value: SpellOption.MortalStrike },
				{ actionId: ActionId.fromSpellId(47465), value: SpellOption.Rend },
				{ actionId: ActionId.fromSpellId(7384), value: SpellOption.Overpower },
				{ actionId: ActionId.fromSpellId(47471), value: SpellOption.Execute },
				{ actionId: ActionId.fromSpellId(47502), value: SpellOption.ThunderClap },
			],
		}),
	],
};
