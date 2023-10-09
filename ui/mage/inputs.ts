import { Spec } from '../core/proto/common.js';
import { ActionId } from '../core/proto_utils/action_id.js';
import { Player } from '../core/player.js';
import { TypedEvent } from '../core/typed_event.js';

import {
	Mage_Rotation_PrimaryFireSpell as PrimaryFireSpell,
	Mage_Options_ArmorType as ArmorType,
} from '../core/proto/mage.js';

import * as InputHelpers from '../core/components/input_helpers.js';

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

export const WaterElementalDisobeyChance = InputHelpers.makeSpecOptionsNumberInput<Spec.SpecMage>({
	fieldName: 'waterElementalDisobeyChance',
	percent: true,
	label: 'Water Ele Disobey %',
	labelTooltip: 'Percent of Water Elemental actions which will fail. This represents the Water Elemental moving around or standing still instead of casting.',
	changeEmitter: (player: Player<Spec.SpecMage>) => TypedEvent.onAny([player.specOptionsChangeEmitter, player.talentsChangeEmitter]),
	showWhen: (player: Player<Spec.SpecMage>) => player.getTalents().summonWaterElemental,
});

export const FocusMagicUptime = InputHelpers.makeSpecOptionsNumberInput<Spec.SpecMage>({
	fieldName: 'focusMagicPercentUptime',
	label: 'Focus Magic Percent Uptime',
	labelTooltip: 'Percent of uptime for Focus Magic Buddy',
	extraCssClasses: ['within-raid-sim-hide'],
});

export const MageRotationConfig = {
	inputs: [
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
			showWhen: (player: Player<Spec.SpecMage>) => player.getTalentTree() == 1,
			changeEmitter: (player: Player<Spec.SpecMage>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
		}),
		// ********************************************************
		//                       FROST INPUTS
		// ********************************************************
		InputHelpers.makeRotationBooleanInput<Spec.SpecMage>({
			fieldName: 'useIceLance',
			label: 'Use Ice Lance',
			labelTooltip: 'Casts Ice Lance at the end of Fingers of Frost, after using Deep Freeze.',
			showWhen: (player: Player<Spec.SpecMage>) => player.getTalentTree() == 2,
			changeEmitter: (player: Player<Spec.SpecMage>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
		}),
		// ********************************************************
		//                      ARCANE INPUTS
		// ********************************************************
		InputHelpers.makeRotationNumberInput<Spec.SpecMage>({
			fieldName: 'only3ArcaneBlastStacksBelowManaPercent',
			percent: true,
			label: 'Stack Arcane Blast to 3 below mana %',
			labelTooltip: 'When below this mana %, AM/ABarr will be used at 3 stacks of AB instead of 4.',
			showWhen: (player: Player<Spec.SpecMage>) => player.getTalentTree() == 0,
			changeEmitter: (player: Player<Spec.SpecMage>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
		}),
		InputHelpers.makeRotationNumberInput<Spec.SpecMage>({
			fieldName: 'blastWithoutMissileBarrageAboveManaPercent',
			percent: true,
			label: 'AB without Missile Barrage above mana %',
			labelTooltip: 'When above this mana %, spam AB until a Missile Barrage proc occurs.',
			showWhen: (player: Player<Spec.SpecMage>) => player.getTalentTree() == 0,
			changeEmitter: (player: Player<Spec.SpecMage>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
		}),
		InputHelpers.makeRotationNumberInput<Spec.SpecMage>({
			fieldName: 'missileBarrageBelowManaPercent',
			percent: true,
			label: 'Use Missile Barrage ASAP below mana %',
			labelTooltip: 'When below this mana %, use Missile Barrage proc as soon as possible. Can be useful to conserve mana.',
			showWhen: (player: Player<Spec.SpecMage>) => player.getTalentTree() == 0,
			changeEmitter: (player: Player<Spec.SpecMage>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecMage>({
			fieldName: 'useArcaneBarrage',
			label: 'Use Arcane Barrage',
			labelTooltip: 'Includes Arcane Barrage in the rotation.',
			enableWhen: (player: Player<Spec.SpecMage>) => player.getTalents().arcaneBarrage,
			showWhen: (player: Player<Spec.SpecMage>) => player.getTalentTree() == 0,
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
