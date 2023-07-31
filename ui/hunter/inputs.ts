import { Spec } from '../core/proto/common.js';
import { ActionId } from '../core/proto_utils/action_id.js';
import { Player } from '../core/player.js';
import { EventID, TypedEvent } from '../core/typed_event.js';
import { makePetTypeInputConfig } from '../core/talents/hunter_pet.js';

import * as InputHelpers from '../core/components/input_helpers.js';

import {
	Hunter,
	Hunter_Rotation as HunterRotation,
	Hunter_Rotation_RotationType as RotationType,
	Hunter_Rotation_StingType as StingType,
	Hunter_Rotation_SpellOption as SpellOption,
	Hunter_Options as HunterOptions,
	Hunter_Options_Ammo as Ammo,
	Hunter_Options_PetType as PetType,
} from '../core/proto/hunter.js';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const WeaponAmmo = InputHelpers.makeSpecOptionsEnumIconInput<Spec.SpecHunter, Ammo>({
	fieldName: 'ammo',
	numColumns: 2,
	values: [
		{ value: Ammo.AmmoNone, tooltip: 'No Ammo' },
		{ actionId: ActionId.fromItemId(52021), value: Ammo.IcebladeArrow },
		{ actionId: ActionId.fromItemId(41165), value: Ammo.SaroniteRazorheads },
		{ actionId: ActionId.fromItemId(41586), value: Ammo.TerrorshaftArrow },
		{ actionId: ActionId.fromItemId(31737), value: Ammo.TimelessArrow },
		{ actionId: ActionId.fromItemId(34581), value: Ammo.MysteriousArrow },
		{ actionId: ActionId.fromItemId(33803), value: Ammo.AdamantiteStinger },
		{ actionId: ActionId.fromItemId(28056), value: Ammo.BlackflightArrow },
	],
});

export const PetTypeInput = makePetTypeInputConfig(true);

export const PetUptime = InputHelpers.makeSpecOptionsNumberInput<Spec.SpecHunter>({
	fieldName: 'petUptime',
	label: 'Pet Uptime (%)',
	labelTooltip: 'Percent of the fight duration for which your pet will be alive.',
	percent: true,
});

export const UseHuntersMark = InputHelpers.makeSpecOptionsBooleanIconInput<Spec.SpecHunter>({
	fieldName: 'useHuntersMark',
	id: ActionId.fromSpellId(53338),
});

export const SniperTrainingUptime = InputHelpers.makeSpecOptionsNumberInput<Spec.SpecHunter>({
	fieldName: 'sniperTrainingUptime',
	label: 'ST Uptime (%)',
	labelTooltip: 'Uptime for the Sniper Training talent, as a percent of the fight duration.',
	percent: true,
	showWhen: (player: Player<Spec.SpecHunter>) => player.getTalents().sniperTraining > 0,
	changeEmitter: (player: Player<Spec.SpecHunter>) => TypedEvent.onAny([player.specOptionsChangeEmitter, player.talentsChangeEmitter]),
});

export const TimeToTrapWeaveMs = InputHelpers.makeSpecOptionsNumberInput<Spec.SpecHunter>({
	fieldName: 'timeToTrapWeaveMs',
	label: 'Weave Time',
	labelTooltip: 'Amount of time for Explosive Trap, in milliseconds, between when you start moving towards the boss and when you re-engage your ranged autos.',
});

export const HunterRotationConfig = {
	inputs: [
		InputHelpers.makeRotationEnumInput<Spec.SpecHunter, RotationType>({
			fieldName: 'type',
			label: 'Type',
			values: [
				{ name: 'Single Target', value: RotationType.SingleTarget },
				{ name: 'AOE', value: RotationType.Aoe },
				{ name: 'Custom', value: RotationType.Custom },
			],
		}),
		InputHelpers.makeRotationEnumInput<Spec.SpecHunter, StingType>({
			fieldName: 'sting',
			label: 'Sting',
			labelTooltip: 'Maintains the selected Sting on the primary target.',
			values: [
				{ name: 'None', value: StingType.NoSting },
				{ name: 'Scorpid Sting', value: StingType.ScorpidSting },
				{ name: 'Serpent Sting', value: StingType.SerpentSting },
			],
			showWhen: (player: Player<Spec.SpecHunter>) => player.getRotation().type == RotationType.SingleTarget,
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecHunter>({
			fieldName: 'trapWeave',
			label: 'Trap Weave',
			labelTooltip: 'Uses Explosive Trap at appropriate times. Note that selecting this will disable Black Arrow because they share a CD.',
			showWhen: (player: Player<Spec.SpecHunter>) => player.getRotation().type != RotationType.Custom,
		}),
		InputHelpers.makeRotationNumberInput<Spec.SpecHunter>({
			fieldName: 'steadyShotMaxDelay',
			label: 'Steady Shot Max Delay (ms)',
			labelTooltip: 'If another higher-priority spell comes off cooldown in the specified time then steady shot is not cast and the rotation waits',
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecHunter>({
			fieldName: 'allowExplosiveShotDownrank',
			label: 'Allow ES Downrank',
			labelTooltip: 'Weaves Explosive Shot Rank 3 during LNL procs. This works because the rank 3 and rank 4 dots can stack.',
			showWhen: (player: Player<Spec.SpecHunter>) => player.getRotation().type != RotationType.Custom && player.getTalents().explosiveShot && player.getTalents().lockAndLoad > 0,
			changeEmitter: (player: Player<Spec.SpecHunter>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecHunter>({
			fieldName: 'multiDotSerpentSting',
			label: 'Multi-Dot Serpent Sting',
			labelTooltip: 'Casts Serpent Sting on multiple targets',
			changeEmitter: (player: Player<Spec.SpecHunter>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
		}),
		InputHelpers.makeCustomRotationInput<Spec.SpecHunter, SpellOption>({
			fieldName: 'customRotation',
			numColumns: 2,
			values: [
				{ actionId: ActionId.fromSpellId(49052), value: SpellOption.SteadyShot },
				{ actionId: ActionId.fromSpellId(49045), value: SpellOption.ArcaneShot },
				{ actionId: ActionId.fromSpellId(49050), value: SpellOption.AimedShot },
				{ actionId: ActionId.fromSpellId(49048), value: SpellOption.MultiShot },
				{ actionId: ActionId.fromSpellId(49001), value: SpellOption.SerpentStingSpell },
				{ actionId: ActionId.fromSpellId(3043), value: SpellOption.ScorpidStingSpell },
				{ actionId: ActionId.fromSpellId(61006), value: SpellOption.KillShot },
				{ actionId: ActionId.fromSpellId(63672), value: SpellOption.BlackArrow },
				{ actionId: ActionId.fromSpellId(53209), value: SpellOption.ChimeraShot },
				{ actionId: ActionId.fromSpellId(60053), value: SpellOption.ExplosiveShot, text: 'R4' },
				{ actionId: ActionId.fromSpellId(60052), value: SpellOption.ExplosiveShotDownrank, text: 'R3' },
				{ actionId: ActionId.fromSpellId(49067), value: SpellOption.ExplosiveTrap },
				{ actionId: ActionId.fromSpellId(58434), value: SpellOption.Volley },
			],
			showWhen: (player: Player<Spec.SpecHunter>) => player.getRotation().type == RotationType.Custom,
		}),
		InputHelpers.makeRotationNumberInput<Spec.SpecHunter>({
			fieldName: 'viperStartManaPercent',
			label: 'Viper Start Mana %',
			labelTooltip: 'Switch to Aspect of the Viper when mana goes below this amount.',
			percent: true,
		}),
		InputHelpers.makeRotationNumberInput<Spec.SpecHunter>({
			fieldName: 'viperStopManaPercent',
			label: 'Viper Stop Mana %',
			labelTooltip: 'Switch back to Aspect of the Hawk when mana goes above this amount.',
			percent: true,
		}),
	],
};
