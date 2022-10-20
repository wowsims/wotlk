import { Spec } from '../core/proto/common.js';


import {
	DeathknightTalents as DeathknightTalents,
	Deathknight_Rotation as DeathknightRotation,
	Deathknight_Options as DeathknightOptions,
	TankDeathknight_Rotation_OptimizationSetting as OptimizationSetting,
	TankDeathknight_Rotation_Opener as Opener,
	TankDeathknight_Rotation_BloodSpell as BloodSpell,
	TankDeathknight_Rotation_Presence as Presence,
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
			labelTooltip: '<p>Chose what opener to perform:</p>\
				<p><b>Regular</b>: Regular opener.</p>\
				<p><b>Threat</b>: Full IT spam for max threat.</p>',
			values: [
				{ name: 'Regular', value: Opener.Regular },
				{ name: 'Threat', value: Opener.Threat },
			],
			changeEmitter: (player: Player<Spec.SpecTankDeathknight>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
		}),
		InputHelpers.makeRotationEnumInput<Spec.SpecTankDeathknight, OptimizationSetting>({
			fieldName: 'optimizationSetting',
			label: 'Optimization Setting',
			labelTooltip: '<p>Chose what metric to optimize :</p>\
				<p><b>Hps</b>: Prioritizes holding runes for healing after damage taken.</p>\
				<p><b>Tps</b>: Prioritizes spending runes for icy touch spam.</p>\
				<p><b>Dps</b>: Prioritizes spending runes for maximizing damage.</p>',
			values: [
				{ name: 'Hps', value: OptimizationSetting.Hps },
				{ name: 'Tps', value: OptimizationSetting.Tps },
				{ name: 'Dps', value: OptimizationSetting.Dps },
			],
			changeEmitter: (player: Player<Spec.SpecTankDeathknight>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
		}),
		InputHelpers.makeRotationEnumInput<Spec.SpecTankDeathknight, BloodSpell>({
			fieldName: 'bloodSpell',
			label: 'Blood Spell',
			labelTooltip: '<p>Chose what blood rune spender to use.</p>',
			values: [
				{ name: 'Blood Strike', value: BloodSpell.BloodStrike },
				{ name: 'Blood Boil', value: BloodSpell.BloodBoil },
				{ name: 'Heart Strike', value: BloodSpell.HeartStrike },
			],
			changeEmitter: (player: Player<Spec.SpecTankDeathknight>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
		}),
	],
};
