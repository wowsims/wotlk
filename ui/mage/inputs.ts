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
				{ name: 'FrostfireBolt', value: PrimaryFireSpell.FrostfireBolt },
			],
			showWhen: (player: Player<Spec.SpecMage>) => player.getRotation().type == RotationType.Fire,
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecMage>({
			fieldName: 'optimizeCdsForExecute',
			label: 'Optimize CDs for execute time',
			labelTooltip: 'Automatically save cooldowns that only have 1 use remaining for execute time',
			showWhen: (player: Player<Spec.SpecMage>) => player.getRotation().type == RotationType.Fire,
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecMage>({
			fieldName: 'maintainImprovedScorch',
			label: 'Maintain Imp. Scorch',
			labelTooltip: 'Always use Scorch when below 5 stacks, or < 5.5s remaining on debuff.',
			showWhen: (player: Player<Spec.SpecMage>) => player.getRotation().type == RotationType.Fire,
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecMage>({
			fieldName: 'lbBeforeHotstreak',
			label: 'Living Bomb Over Hot Streak',
			labelTooltip: 'Choose to reapply living bomb before consuming hot streak',
			showWhen: (player: Player<Spec.SpecMage>) => player.getRotation().type == RotationType.Fire,
		}),
		// ********************************************************
		//                       FROST INPUTS
		// ********************************************************
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
			fieldName: 'minBlastBeforeMissiles',
			label: 'Min ABs before missiles',
			labelTooltip: 'Minimum arcane blasts to cast before using a missile barrage proc',
			showWhen: (player: Player<Spec.SpecMage>) => player.getRotation().type == RotationType.Arcane,
		}),
		InputHelpers.makeRotationNumberInput<Spec.SpecMage>({
			fieldName: 'num4StackBlastsToMissilesGamble',
			label: 'Switch to AM Gamble At',
			labelTooltip: 'Number of times mage has cast a 4 stacked arcane blast over the whole fight before gambling on AM when you dont have missile barrage',
			showWhen: (player: Player<Spec.SpecMage>) => player.getRotation().type == RotationType.Arcane,
		}),
		InputHelpers.makeRotationNumberInput<Spec.SpecMage>({
			fieldName: 'num4StackBlastsToEarlyMissiles',
			label: 'Switch to ASAP missiles barrage At',
			labelTooltip: 'Switch to using missiles barrage ASAP after this many 4 cost ABs',
			showWhen: (player: Player<Spec.SpecMage>) => player.getRotation().type == RotationType.Arcane,
		}),
		InputHelpers.makeRotationNumberInput<Spec.SpecMage>({
			fieldName: 'extraBlastsDuringFirstAp',
			label: 'Extra blasts during first AP',
			labelTooltip: 'Number of extra arcane blasts to use during your first cooldown phase with AP',
			showWhen: (player: Player<Spec.SpecMage>) => player.getRotation().type == RotationType.Arcane,
		}),
	],
};
