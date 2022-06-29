import { Spec } from '/tbc/core/proto/common.js';
import { RogueTalents as RogueTalents } from '/tbc/core/proto/rogue.js';
import { RogueSpecs } from '/tbc/core/proto_utils/utils.js';
import { Player } from '/tbc/core/player.js';

import { TalentsConfig, TalentsPicker, newTalentsConfig } from './talents_picker.js';

export class RogueTalentsPicker extends TalentsPicker<Spec.SpecRogue> {
	constructor(parent: HTMLElement, player: Player<Spec.SpecRogue>) {
		super(parent, player, rogueTalentsConfig);
	}
}

export const rogueTalentsConfig: TalentsConfig<Spec.SpecRogue> = newTalentsConfig([
	{
		name: 'Assassination',
		backgroundUrl: 'https://wow.zamimg.com/images/wow/talents/backgrounds/tbc/182.jpg',
		talents: [
			{
				fieldName: 'improvedEviscerate',
				location: {
					rowIdx: 0,
					colIdx: 0,
				},
				spellIds: [14162],
				maxPoints: 3,
			},
			{
				//fieldName: 'remorselessAttacks',
				location: {
					rowIdx: 0,
					colIdx: 1,
				},
				spellIds: [14144, 14148],
				maxPoints: 2,
			},
			{
				fieldName: 'malice',
				location: {
					rowIdx: 0,
					colIdx: 2,
				},
				spellIds: [14138],
				maxPoints: 5,
			},
			{
				fieldName: 'ruthlessness',
				location: {
					rowIdx: 1,
					colIdx: 0,
				},
				spellIds: [14156, 14160],
				maxPoints: 3,
			},
			{
				fieldName: 'murder',
				location: {
					rowIdx: 1,
					colIdx: 1,
				},
				spellIds: [14158],
				maxPoints: 2,
			},
			{
				fieldName: 'puncturingWounds',
				location: {
					rowIdx: 1,
					colIdx: 3,
				},
				spellIds: [13733, 13865],
				maxPoints: 3,
			},
			{
				fieldName: 'relentlessStrikes',
				location: {
					rowIdx: 2,
					colIdx: 0,
				},
				spellIds: [14179],
				maxPoints: 1,
			},
			{
				fieldName: 'improvedExposeArmor',
				location: {
					rowIdx: 2,
					colIdx: 1,
				},
				spellIds: [14168],
				maxPoints: 2,
			},
			{
				fieldName: 'lethality',
				location: {
					rowIdx: 2,
					colIdx: 2,
				},
				prereqLocation: {
					rowIdx: 0,
					colIdx: 2,
				},
				spellIds: [14128, 14132, 14135],
				maxPoints: 5,
			},
			{
				fieldName: 'vilePoisons',
				location: {
					rowIdx: 3,
					colIdx: 1,
				},
				spellIds: [16513, 16514, 16515, 16719],
				maxPoints: 5,
			},
			{
				fieldName: 'improvedPoisons',
				location: {
					rowIdx: 3,
					colIdx: 2,
				},
				spellIds: [14113],
				maxPoints: 5,
			},
			{
				//fieldName: 'fleetFooted',
				location: {
					rowIdx: 4,
					colIdx: 0,
				},
				spellIds: [31208],
				maxPoints: 2,
			},
			{
				fieldName: 'coldBlood',
				location: {
					rowIdx: 4,
					colIdx: 1,
				},
				spellIds: [14177],
				maxPoints: 1,
			},
			{
				//fieldName: 'improvedKidneyShot',
				location: {
					rowIdx: 4,
					colIdx: 2,
				},
				spellIds: [14174],
				maxPoints: 3,
			},
			{
				fieldName: 'quickRecovery',
				location: {
					rowIdx: 4,
					colIdx: 3,
				},
				spellIds: [31244],
				maxPoints: 2,
			},
			{
				fieldName: 'sealFate',
				location: {
					rowIdx: 5,
					colIdx: 1,
				},
				prereqLocation: {
					rowIdx: 4,
					colIdx: 1,
				},
				spellIds: [14186, 14190, 14193],
				maxPoints: 5,
			},
			{
				fieldName: 'masterPoisoner',
				location: {
					rowIdx: 5,
					colIdx: 2,
				},
				spellIds: [31226],
				maxPoints: 2,
			},
			{
				fieldName: 'vigor',
				location: {
					rowIdx: 6,
					colIdx: 1,
				},
				spellIds: [14983],
				maxPoints: 1,
			},
			{
				//fieldName: 'deadenedNerves',
				location: {
					rowIdx: 6,
					colIdx: 2,
				},
				spellIds: [31380, 31382],
				maxPoints: 5,
			},
			{
				fieldName: 'findWeakness',
				location: {
					rowIdx: 7,
					colIdx: 2,
				},
				spellIds: [31233, 31239],
				maxPoints: 5,
			},
			{
				fieldName: 'mutilate',
				location: {
					rowIdx: 8,
					colIdx: 1,
				},
				prereqLocation: {
					rowIdx: 6,
					colIdx: 1,
				},
				spellIds: [1329],
				maxPoints: 1,
			},
		],
	},
	{
		name: 'Combat',
		backgroundUrl: 'https://wow.zamimg.com/images/wow/talents/backgrounds/tbc/181.jpg',
		talents: [
			{
				//fieldName: 'improvedGouge',
				location: {
					rowIdx: 0,
					colIdx: 0,
				},
				spellIds: [13741, 13793, 13792],
				maxPoints: 3,
			},
			{
				fieldName: 'improvedSinisterStrike',
				location: {
					rowIdx: 0,
					colIdx: 1,
				},
				spellIds: [13732, 13863],
				maxPoints: 2,
			},
			{
				fieldName: 'lightningReflexes',
				location: {
					rowIdx: 0,
					colIdx: 2,
				},
				spellIds: [13712, 13788],
				maxPoints: 5,
			},
			{
				fieldName: 'improvedSliceAndDice',
				location: {
					rowIdx: 1,
					colIdx: 0,
				},
				spellIds: [14165],
				maxPoints: 3,
			},
			{
				fieldName: 'deflection',
				location: {
					rowIdx: 1,
					colIdx: 1,
				},
				spellIds: [13713, 13853],
				maxPoints: 5,
			},
			{
				fieldName: 'precision',
				location: {
					rowIdx: 1,
					colIdx: 2,
				},
				spellIds: [13705, 13832, 13843],
				maxPoints: 5,
			},
			{
				//fieldName: 'endurance',
				location: {
					rowIdx: 2,
					colIdx: 0,
				},
				spellIds: [13742, 13872],
				maxPoints: 2,
			},
			{
				//fieldName: 'riposte',
				location: {
					rowIdx: 2,
					colIdx: 1,
				},
				prereqLocation: {
					rowIdx: 1,
					colIdx: 1,
				},
				spellIds: [14251],
				maxPoints: 1,
			},
			{
				//fieldName: 'improvedSprint',
				location: {
					rowIdx: 2,
					colIdx: 3,
				},
				spellIds: [13743, 13875],
				maxPoints: 2,
			},
			{
				//fieldName: 'improvedKick',
				location: {
					rowIdx: 3,
					colIdx: 0,
				},
				spellIds: [13754, 13867],
				maxPoints: 2,
			},
			{
				fieldName: 'daggerSpecialization',
				location: {
					rowIdx: 3,
					colIdx: 1,
				},
				spellIds: [13706, 13804],
				maxPoints: 5,
			},
			{
				fieldName: 'dualWieldSpecialization',
				location: {
					rowIdx: 3,
					colIdx: 2,
				},
				prereqLocation: {
					rowIdx: 1,
					colIdx: 2,
				},
				spellIds: [13715, 13848, 13849, 13851],
				maxPoints: 5,
			},
			{
				fieldName: 'maceSpecialization',
				location: {
					rowIdx: 4,
					colIdx: 0,
				},
				spellIds: [13709, 13800],
				maxPoints: 5,
			},
			{
				fieldName: 'bladeFlurry',
				location: {
					rowIdx: 4,
					colIdx: 1,
				},
				spellIds: [13877],
				maxPoints: 1,
			},
			{
				fieldName: 'swordSpecialization',
				location: {
					rowIdx: 4,
					colIdx: 2,
				},
				spellIds: [13960],
				maxPoints: 5,
			},
			{
				fieldName: 'fistWeaponSpecialization',
				location: {
					rowIdx: 4,
					colIdx: 3,
				},
				spellIds: [13707, 13966],
				maxPoints: 5,
			},
			{
				//fieldName: 'bladeTwisting',
				location: {
					rowIdx: 5,
					colIdx: 0,
				},
				spellIds: [31124, 31126],
				maxPoints: 2,
			},
			{
				fieldName: 'weaponExpertise',
				location: {
					rowIdx: 5,
					colIdx: 1,
				},
				prereqLocation: {
					rowIdx: 4,
					colIdx: 1,
				},
				spellIds: [30919],
				maxPoints: 2,
			},
			{
				fieldName: 'aggression',
				location: {
					rowIdx: 5,
					colIdx: 2,
				},
				spellIds: [18427],
				maxPoints: 3,
			},
			{
				fieldName: 'vitality',
				location: {
					rowIdx: 6,
					colIdx: 0,
				},
				spellIds: [31122],
				maxPoints: 2,
			},
			{
				fieldName: 'adrenalineRush',
				location: {
					rowIdx: 6,
					colIdx: 1,
				},
				spellIds: [13750],
				maxPoints: 1,
			},
			{
				//fieldName: 'nervesOfSteel',
				location: {
					rowIdx: 6,
					colIdx: 2,
				},
				spellIds: [31130],
				maxPoints: 2,
			},
			{
				fieldName: 'combatPotency',
				location: {
					rowIdx: 7,
					colIdx: 2,
				},
				spellIds: [35541, 35550],
				maxPoints: 5,
			},
			{
				fieldName: 'surpriseAttacks',
				location: {
					rowIdx: 8,
					colIdx: 1,
				},
				prereqLocation: {
					rowIdx: 6,
					colIdx: 1,
				},
				spellIds: [32601],
				maxPoints: 1,
			},
		],
	},
	{
		name: 'Subtlety',
		backgroundUrl: 'https://wow.zamimg.com/images/wow/talents/backgrounds/tbc/183.jpg',
		talents: [
			{
				//fieldName: 'masterOfDeception',
				location: {
					rowIdx: 0,
					colIdx: 1,
				},
				spellIds: [13958, 13970],
				maxPoints: 5,
			},
			{
				fieldName: 'opportunity',
				location: {
					rowIdx: 0,
					colIdx: 2,
				},
				spellIds: [14057, 14072],
				maxPoints: 5,
			},
			{
				fieldName: 'sleightOfHand',
				location: {
					rowIdx: 1,
					colIdx: 0,
				},
				spellIds: [30892],
				maxPoints: 2,
			},
			{
				//fieldName: 'dirtyTricks',
				location: {
					rowIdx: 1,
					colIdx: 1,
				},
				spellIds: [14076, 14094],
				maxPoints: 2,
			},
			{
				//fieldName: 'camoflauge',
				location: {
					rowIdx: 1,
					colIdx: 2,
				},
				spellIds: [13975, 14062],
				maxPoints: 5,
			},
			{
				fieldName: 'initiative',
				location: {
					rowIdx: 2,
					colIdx: 0,
				},
				spellIds: [13976, 13979],
				maxPoints: 3,
			},
			{
				fieldName: 'ghostlyStrike',
				location: {
					rowIdx: 2,
					colIdx: 1,
				},
				spellIds: [14278],
				maxPoints: 1,
			},
			{
				fieldName: 'improvedAmbush',
				location: {
					rowIdx: 2,
					colIdx: 2,
				},
				spellIds: [14079],
				maxPoints: 3,
			},
			{
				//fieldName: 'setup',
				location: {
					rowIdx: 3,
					colIdx: 0,
				},
				spellIds: [13983, 14070],
				maxPoints: 3,
			},
			{
				fieldName: 'elusiveness',
				location: {
					rowIdx: 3,
					colIdx: 1,
				},
				spellIds: [13981, 14066],
				maxPoints: 2,
			},
			{
				fieldName: 'serratedBlades',
				location: {
					rowIdx: 3,
					colIdx: 2,
				},
				spellIds: [14171],
				maxPoints: 3,
			},
			{
				//fieldName: 'heightenedSenses',
				location: {
					rowIdx: 4,
					colIdx: 0,
				},
				spellIds: [30894],
				maxPoints: 2,
			},
			{
				fieldName: 'preparation',
				location: {
					rowIdx: 4,
					colIdx: 1,
				},
				spellIds: [14185],
				maxPoints: 1,
			},
			{
				fieldName: 'dirtyDeeds',
				location: {
					rowIdx: 4,
					colIdx: 2,
				},
				spellIds: [14082],
				maxPoints: 2,
			},
			{
				fieldName: 'hemorrhage',
				location: {
					rowIdx: 4,
					colIdx: 3,
				},
				prereqLocation: {
					rowIdx: 3,
					colIdx: 2,
				},
				spellIds: [16511],
				maxPoints: 1,
			},
			{
				fieldName: 'masterOfSubtlety',
				location: {
					rowIdx: 5,
					colIdx: 0,
				},
				spellIds: [31221],
				maxPoints: 3,
			},
			{
				fieldName: 'deadliness',
				location: {
					rowIdx: 5,
					colIdx: 2,
				},
				spellIds: [30902],
				maxPoints: 5,
			},
			{
				//fieldName: 'envelopingShadows',
				location: {
					rowIdx: 6,
					colIdx: 0,
				},
				spellIds: [31211],
				maxPoints: 3,
			},
			{
				fieldName: 'premeditation',
				location: {
					rowIdx: 6,
					colIdx: 1,
				},
				prereqLocation: {
					rowIdx: 4,
					colIdx: 1,
				},
				spellIds: [14183],
				maxPoints: 1,
			},
			{
				//fieldName: 'cheatDeath',
				location: {
					rowIdx: 6,
					colIdx: 2,
				},
				spellIds: [31228],
				maxPoints: 3,
			},
			{
				fieldName: 'sinisterCalling',
				location: {
					rowIdx: 7,
					colIdx: 1,
				},
				prereqLocation: {
					rowIdx: 6,
					colIdx: 1,
				},
				spellIds: [31216],
				maxPoints: 5,
			},
			{
				fieldName: 'shadowstep',
				location: {
					rowIdx: 8,
					colIdx: 1,
				},
				spellIds: [36554],
				maxPoints: 1,
			},
		],
	},
]);
