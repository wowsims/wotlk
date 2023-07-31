import { Spec } from '../core/proto/common.js';
import { ActionId } from '../core/proto_utils/action_id.js';
import { Player } from '../core/player.js';
import { EventID, TypedEvent } from '../core/typed_event.js';

import {
	Mage,
	MageTalents as MageTalents,
	Mage_Rotation as MageRotation,
	Mage_Rotation_Type as RotationType,
	Mage_Rotation_PrimaryFireSpell as PrimaryFireSpell,
	Mage_Rotation_AoeRotation as AoeRotationSpells,
	Mage_Options as MageOptions,
	Mage_Options_ArmorType as ArmorType,
} from '../core/proto/mage.js';

import * as InputHelpers from '../core/components/input_helpers.js';
import * as Presets from './presets.js';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const Armor = InputHelpers.makeSpecOptionsEnumIconInput<Spec.SpecMage, ArmorType>({
	fieldName: 'armor',
	values: [
		{ value: ArmorType.NoArmor, tooltip: 'No Armor' },
		{ actionId: ActionId.fromSpellId(43024), value: ArmorType.MageArmor },
		{ actionId: ActionId.fromSpellId(43046), value: ArmorType.MoltenArmor },
	],
});

export const EvocationTicks = InputHelpers.makeSpecOptionsNumberInput<Spec.SpecMage>({
	fieldName: 'evocationTicks',
	label: '# Evocation Ticks',
	labelTooltip: 'The number of ticks of Evocation to use, or 0 to use the full duration.',
});

export const FocusMagicUptime = InputHelpers.makeSpecOptionsNumberInput<Spec.SpecMage>({
	fieldName: 'focusMagicPercentUptime',
	label: 'Focus Magic Percent Uptime',
	labelTooltip: 'Percent of uptime for Focus Magic Buddy',
	extraCssClasses: ['within-raid-sim-hide'],
});

export const ReactionTime = InputHelpers.makeSpecOptionsNumberInput<Spec.SpecMage>({
	fieldName: 'reactionTimeMs',
	label: 'Reaction Time (ms)',
	labelTooltip: 'Duration, in milliseconds, for player reaction time. Only used for a few effects (Missile Barrage / Hot Streak / Brain Freeze).',
});

export const MageRotationConfig = {
	inputs: [
		InputHelpers.makeRotationEnumInput<Spec.SpecMage, RotationType>({
			fieldName: 'type',
			label: 'Spec',
			labelTooltip: 'Switches between spec rotation settings. Will also update talents to defaults for the selected spec.',
			values: [
				{ name: 'Arcane', value: RotationType.Arcane },
				{ name: 'Fire', value: RotationType.Fire },
				{ name: 'Frost', value: RotationType.Frost },
			],
			setValue: (eventID: EventID, player: Player<Spec.SpecMage>, newValue: number) => {
				const newRotation = player.getRotation();
				newRotation.type = newValue;

				TypedEvent.freezeAllAndDo(() => {
					if (newRotation.type == RotationType.Arcane) {
						player.setTalentsString(eventID, Presets.ArcaneTalents.data.talentsString);
						player.setGlyphs(eventID, Presets.ArcaneTalents.data.glyphs!);
					} else if (newRotation.type == RotationType.Fire) {
						player.setTalentsString(eventID, Presets.FireTalents.data.talentsString);
						player.setGlyphs(eventID, Presets.FireTalents.data.glyphs!);
					} else if (newRotation.type == RotationType.Frost) {
						player.setTalentsString(eventID, Presets.FrostTalents.data.talentsString);
						player.setGlyphs(eventID, Presets.FrostTalents.data.glyphs!);
					}

					player.setRotation(eventID, newRotation);
				});
			},
		}),
		// ********************************************************
		//                        AOE INPUTS
		// ********************************************************
		InputHelpers.makeRotationEnumInput<Spec.SpecMage, AoeRotationSpells>({
			fieldName: 'aoe',
			label: 'Primary Spell',
			values: [
				{ name: 'Arcane Explosion', value: AoeRotationSpells.ArcaneExplosion },
				{ name: 'Flamestrike', value: AoeRotationSpells.Flamestrike },
				{ name: 'Blizzard', value: AoeRotationSpells.Blizzard },
			],
			showWhen: (player: Player<Spec.SpecMage>) => player.getRotation().type == RotationType.Aoe,
		}),
		// ********************************************************
		//                       FIRE INPUTS
		// ********************************************************
		InputHelpers.makeRotationEnumInput<Spec.SpecMage, PrimaryFireSpell>({
			fieldName: 'primaryFireSpell',
			label: 'Primary Spell',
			values: [
				{ name: 'Fireball', value: PrimaryFireSpell.Fireball },
				{ name: 'Frostfire Bolt', value: PrimaryFireSpell.FrostfireBolt },
				{ name: 'Scorch', value: PrimaryFireSpell.Scorch },
			],
			showWhen: (player: Player<Spec.SpecMage>) => player.getRotation().type == RotationType.Fire,
		}),
		InputHelpers.makeRotationNumberInput<Spec.SpecMage>({
			fieldName: 'pyroblastDelayMs',
			label: 'Pyroblast Delay (ms)',
			labelTooltip: `
				<p>Adds a delay to Pyroblast after a Hot Streak to prevent ignite munching. 50ms is a good default for this.</p>
				<p>There is no way to do this perfectly in-game, but a cqs macro can do this with about 70-90% reliability.</p>
			`,
			showWhen: (player: Player<Spec.SpecMage>) => player.getRotation().type == RotationType.Fire && player.getRotation().primaryFireSpell == PrimaryFireSpell.Fireball && player.getSpecOptions().igniteMunching,
			changeEmitter: (player: Player<Spec.SpecMage>) => TypedEvent.onAny([player.rotationChangeEmitter, player.specOptionsChangeEmitter]),
		}),
		// ********************************************************
		//                       FROST INPUTS
		// ********************************************************
		InputHelpers.makeRotationBooleanInput<Spec.SpecMage>({
			fieldName: 'useIceLance',
			label: 'Use Ice Lance',
			labelTooltip: 'Casts Ice Lance at the end of Fingers of Frost, after using Deep Freeze.',
			showWhen: (player: Player<Spec.SpecMage>) => player.getRotation().type == RotationType.Frost,
		}),
		InputHelpers.makeRotationNumberInput<Spec.SpecMage>({
			fieldName: 'waterElementalDisobeyChance',
			percent: true,
			label: 'Water Ele Disobey %',
			labelTooltip: 'Percent of Water Elemental actions which will fail. This represents the Water Elemental moving around or standing still instead of casting.',
			changeEmitter: (player: Player<Spec.SpecMage>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
			showWhen: (player: Player<Spec.SpecMage>) => player.getRotation().type == RotationType.Frost,
			enableWhen: (player: Player<Spec.SpecMage>) => player.getTalents().summonWaterElemental,
		}),
		// ********************************************************
		//                      ARCANE INPUTS
		// ********************************************************
		InputHelpers.makeRotationNumberInput<Spec.SpecMage>({
			fieldName: 'only3ArcaneBlastStacksBelowManaPercent',
			percent: true,
			label: 'Stack Arcane Blast to 3 below mana %',
			labelTooltip: 'When below this mana %, AM/ABarr will be used at 3 stacks of AB instead of 4.',
			showWhen: (player: Player<Spec.SpecMage>) => player.getRotation().type == RotationType.Arcane,
		}),
		InputHelpers.makeRotationNumberInput<Spec.SpecMage>({
			fieldName: 'blastWithoutMissileBarrageAboveManaPercent',
			percent: true,
			label: 'AB without Missile Barrage above mana %',
			labelTooltip: 'When above this mana %, spam AB until a Missile Barrage proc occurs.',
			showWhen: (player: Player<Spec.SpecMage>) => player.getRotation().type == RotationType.Arcane,
		}),
		InputHelpers.makeRotationNumberInput<Spec.SpecMage>({
			fieldName: 'extraBlastsDuringFirstAp',
			label: 'Extra ABs during first AP',
			labelTooltip: 'Extend AB streak by this mana casts, during the first Arcane Power CD duration.',
			showWhen: (player: Player<Spec.SpecMage>) => player.getRotation().type == RotationType.Arcane,
		}),
		InputHelpers.makeRotationNumberInput<Spec.SpecMage>({
			fieldName: 'missileBarrageBelowArcaneBlastStacks',
			label: 'Use Missile Barrage below n AB stacks',
			labelTooltip: 'Setting this to 1 or 2 can potentially be a DPS increase with Arcane Barrage rotation or T8 4pc set bonus.',
			showWhen: (player: Player<Spec.SpecMage>) => player.getRotation().type == RotationType.Arcane,
		}),
		InputHelpers.makeRotationNumberInput<Spec.SpecMage>({
			fieldName: 'missileBarrageBelowManaPercent',
			percent: true,
			label: 'Use Missile Barrage ASAP below mana %',
			labelTooltip: 'When below this mana %, use Missile Barrage proc as soon as possible. Can be useful to conserve mana.',
			showWhen: (player: Player<Spec.SpecMage>) => player.getRotation().type == RotationType.Arcane,
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecMage>({
			fieldName: 'useArcaneBarrage',
			label: 'Use Arcane Barrage',
			labelTooltip: 'Includes Arcane Barrage in the rotation.',
			enableWhen: (player: Player<Spec.SpecMage>) => player.getTalents().arcaneBarrage,
			showWhen: (player: Player<Spec.SpecMage>) => player.getRotation().type == RotationType.Arcane,
			changeEmitter: (player: Player<Spec.SpecMage>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
		}),

		InputHelpers.makeRotationBooleanInput<Spec.SpecMage>({
			fieldName: 'maintainImprovedScorch',
			label: 'Maintain Imp. Scorch',
			labelTooltip: 'Always use Scorch when below 5 stacks, or < 4s remaining on debuff.',
			showWhen: (player: Player<Spec.SpecMage>) => player.getTalents().improvedScorch > 0,
			changeEmitter: (player: Player<Spec.SpecMage>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
		}),
	],
};
