import { Spec } from '../core/proto/common.js';

import {
	TankDeathknight_Rotation_OptimizationSetting as OptimizationSetting,
	TankDeathknight_Rotation_Opener as Opener,
	TankDeathknight_Rotation_BloodSpell as BloodSpell,
	TankDeathknight_Rotation_Presence as Presence,
	TankDeathknight_Rotation_BloodTapPrio as BloodTapPrio,
} from '../core/proto/deathknight.js';


import * as InputHelpers from '../core/components/input_helpers.js';
import { Player } from '../core/player';
import { TypedEvent } from '../core/typed_event';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const StartingRunicPower = InputHelpers.makeSpecOptionsNumberInput<Spec.SpecTankDeathknight>({
	fieldName: 'startingRunicPower',
	label: 'Starting Runic Power',
	labelTooltip: 'Initial RP at the start of each iteration.',
})

export const DefensiveCdDelay = InputHelpers.makeSpecOptionsNumberInput<Spec.SpecTankDeathknight>({
	fieldName: 'defensiveDelay',
	label: 'Defensives Delay',
	labelTooltip: 'Minimum delay between using more defensive cooldowns.',
})

export const TankDeathKnightRotationConfig = {
	inputs: [
		InputHelpers.makeRotationEnumInput<Spec.SpecTankDeathknight, Presence>({
			fieldName: 'presence',
			label: 'Presence',
			labelTooltip: 'Presence to be in during the encounter.',
			values: [
				{ name: 'Blood', value: Presence.Blood },
				{ name: 'Frost', value: Presence.Frost },
				{ name: 'Unholy', value: Presence.Unholy },
			],
			changeEmitter: (player: Player<Spec.SpecTankDeathknight>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
		}),
		InputHelpers.makeRotationEnumInput<Spec.SpecTankDeathknight, Opener>({
			fieldName: 'opener',
			label: 'Opener',
			labelTooltip: 'Chose what opener to perform:<br>\
				<b>Regular</b>: Regular opener.<br>\
				<b>Threat</b>: Full IT spam for max threat.',
			values: [
				{ name: 'Regular', value: Opener.Regular },
				{ name: 'Threat', value: Opener.Threat },
			],
			changeEmitter: (player: Player<Spec.SpecTankDeathknight>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
		}),
		InputHelpers.makeRotationEnumInput<Spec.SpecTankDeathknight, OptimizationSetting>({
			fieldName: 'optimizationSetting',
			label: 'Optimization Setting',
			labelTooltip: 'Chose what metric to optimize:<br>\
				<b>Hps</b>: Prioritizes holding runes for healing after damage taken.<br>\
				<b>Tps</b>: Prioritizes spending runes for icy touch spam.',
			values: [
				{ name: 'Hps', value: OptimizationSetting.Hps },
				{ name: 'Tps', value: OptimizationSetting.Tps },
			],
			changeEmitter: (player: Player<Spec.SpecTankDeathknight>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
		}),
		InputHelpers.makeRotationEnumInput<Spec.SpecTankDeathknight, BloodSpell>({
			fieldName: 'bloodSpell',
			label: 'Blood Spell',
			labelTooltip: 'Chose what blood rune spender to use.',
			values: [
				{ name: 'Blood Strike', value: BloodSpell.BloodStrike },
				{ name: 'Blood Boil', value: BloodSpell.BloodBoil },
				{ name: 'Heart Strike', value: BloodSpell.HeartStrike },
			],
			changeEmitter: (player: Player<Spec.SpecTankDeathknight>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
		}),
		InputHelpers.makeRotationEnumInput<Spec.SpecTankDeathknight, BloodTapPrio>({
			fieldName: 'bloodTapPrio',
			label: 'Blood Tap',
			labelTooltip: 'Chose how to use Blood Tap:<br>\
				<b>Use as Defensive Cooldown</b>: Use as defined in Cooldowns (Requires T10 4pc).<br>\
				<b>Offensive</b>: Use Blood Tap for extra Icy Touches.',
			values: [
				{ name: 'Use as Defensive Cooldown', value: BloodTapPrio.Defensive },
				{ name: 'Offensive', value: BloodTapPrio.Offensive },
			],
			changeEmitter: (player: Player<Spec.SpecTankDeathknight>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
		}),
	],
};
