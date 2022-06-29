import { Spec } from '/tbc/core/proto/common.js';
import { HunterTalents as HunterTalents } from '/tbc/core/proto/hunter.js';
import { HunterSpecs } from '/tbc/core/proto_utils/utils.js';
import { Player } from '/tbc/core/player.js';

import { TalentsConfig, TalentsPicker, newTalentsConfig } from './talents_picker.js';

export class HunterTalentsPicker extends TalentsPicker<Spec.SpecHunter> {
	constructor(parent: HTMLElement, player: Player<Spec.SpecHunter>) {
		super(parent, player, hunterTalentsConfig);
	}
}

export const hunterTalentsConfig: TalentsConfig<Spec.SpecHunter> = newTalentsConfig([
	{
		name: 'Beast Mastery',
		backgroundUrl: 'https://wow.zamimg.com/images/wow/talents/backgrounds/tbc/361.jpg',
		talents: [
			{
				fieldName: 'improvedAspectOfTheHawk',
				location: {
					rowIdx: 0,
					colIdx: 1,
				},
				spellIds: [19552],
				maxPoints: 5,
			},
			{
				fieldName: 'enduranceTraining',
				location: {
					rowIdx: 0,
					colIdx: 2,
				},
				spellIds: [19583],
				maxPoints: 5,
			},
			{
				fieldName: 'focusedFire',
				location: {
					rowIdx: 1,
					colIdx: 0,
				},
				spellIds: [35029],
				maxPoints: 2,
			},
			{
				//fieldName: 'improvedAspectOfTheMonkey',
				location: {
					rowIdx: 1,
					colIdx: 1,
				},
				spellIds: [19549],
				maxPoints: 3,
			},
			{
				//fieldName: 'thickHide',
				location: {
					rowIdx: 1,
					colIdx: 2,
				},
				spellIds: [19609, 19610, 19612],
				maxPoints: 3,
			},
			{
				//fieldName: 'improvedRevivePet',
				location: {
					rowIdx: 1,
					colIdx: 3,
				},
				spellIds: [24443, 19575],
				maxPoints: 2,
			},
			{
				//fieldName: 'pathfinding',
				location: {
					rowIdx: 2,
					colIdx: 0,
				},
				spellIds: [19559],
				maxPoints: 2,
			},
			{
				//fieldName: 'Bestial Swiftness',
				location: {
					rowIdx: 2,
					colIdx: 1,
				},
				spellIds: [19596],
				maxPoints: 1,
			},
			{
				fieldName: 'unleashedFury',
				location: {
					rowIdx: 2,
					colIdx: 2,
				},
				spellIds: [19616],
				maxPoints: 5,
			},
			{
				//fieldName: 'improvedMendPet',
				location: {
					rowIdx: 3,
					colIdx: 1,
				},
				spellIds: [19572],
				maxPoints: 2,
			},
			{
				fieldName: 'ferocity',
				location: {
					rowIdx: 3,
					colIdx: 2,
				},
				spellIds: [19598],
				maxPoints: 5,
			},
			{
				//fieldName: 'Spirit Bond',
				location: {
					rowIdx: 4,
					colIdx: 0,
				},
				spellIds: [19578, 20895],
				maxPoints: 2,
			},
			{
				//fieldName: 'Intimidation',
				location: {
					rowIdx: 4,
					colIdx: 1,
				},
				spellIds: [19577],
				maxPoints: 1,
			},
			{
				fieldName: 'bestialDiscipline',
				location: {
					rowIdx: 4,
					colIdx: 3,
				},
				spellIds: [19590, 19592],
				maxPoints: 2,
			},
			{
				fieldName: 'animalHandler',
				location: {
					rowIdx: 5,
					colIdx: 0,
				},
				spellIds: [34453],
				maxPoints: 2,
			},
			{
				fieldName: 'frenzy',
				location: {
					rowIdx: 5,
					colIdx: 2,
				},
				prereqLocation: {
					rowIdx: 3,
					colIdx: 2,
				},
				spellIds: [19621],
				maxPoints: 5,
			},
			{
				fieldName: 'ferociousInspiration',
				location: {
					rowIdx: 6,
					colIdx: 0,
				},
				spellIds: [34455, 34459],
				maxPoints: 3,
			},
			{
				fieldName: 'bestialWrath',
				location: {
					rowIdx: 6,
					colIdx: 1,
				},
				prereqLocation: {
					rowIdx: 4,
					colIdx: 1,
				},
				spellIds: [19574],
				maxPoints: 1,
			},
			{
				//fieldName: 'catlikeReflexes',
				location: {
					rowIdx: 6,
					colIdx: 2,
				},
				spellIds: [34462, 34464],
				maxPoints: 3,
			},
			{
				fieldName: 'serpentsSwiftness',
				location: {
					rowIdx: 7,
					colIdx: 2,
				},
				spellIds: [34466],
				maxPoints: 5,
			},
			{
				fieldName: 'theBeastWithin',
				location: {
					rowIdx: 8,
					colIdx: 1,
				},
				prereqLocation: {
					rowIdx: 6,
					colIdx: 1,
				},
				spellIds: [34692],
				maxPoints: 1,
			},
		],
	},
	{
		name: 'Marksmanship',
		backgroundUrl: 'https://wow.zamimg.com/images/wow/talents/backgrounds/tbc/363.jpg',
		talents: [
			{
				//fieldName: 'improvedConsussiveShot',
				location: {
					rowIdx: 0,
					colIdx: 1,
				},
				spellIds: [19407, 19412],
				maxPoints: 5,
			},
			{
				fieldName: 'lethalShots',
				location: {
					rowIdx: 0,
					colIdx: 2,
				},
				spellIds: [19426, 19427, 19429],
				maxPoints: 5,
			},
			{
				fieldName: 'improvedHuntersMark',
				location: {
					rowIdx: 1,
					colIdx: 1,
				},
				spellIds: [19421],
				maxPoints: 5,
			},
			{
				fieldName: 'efficiency',
				location: {
					rowIdx: 1,
					colIdx: 2,
				},
				spellIds: [19416],
				maxPoints: 5,
			},
			{
				fieldName: 'goForTheThroat',
				location: {
					rowIdx: 2,
					colIdx: 0,
				},
				spellIds: [34950, 34954],
				maxPoints: 2,
			},
			{
				fieldName: 'improvedArcaneShot',
				location: {
					rowIdx: 2,
					colIdx: 1,
				},
				spellIds: [19454],
				maxPoints: 5,
			},
			{
				fieldName: 'aimedShot',
				location: {
					rowIdx: 2,
					colIdx: 2,
				},
				spellIds: [19434],
				maxPoints: 1,
			},
			{
				fieldName: 'rapidKilling',
				location: {
					rowIdx: 2,
					colIdx: 3,
				},
				spellIds: [34948],
				maxPoints: 2,
			},
			{
				fieldName: 'improvedStings',
				location: {
					rowIdx: 3,
					colIdx: 1,
				},
				spellIds: [19464],
				maxPoints: 5,
			},
			{
				fieldName: 'mortalShots',
				location: {
					rowIdx: 3,
					colIdx: 2,
				},
				prereqLocation: {
					rowIdx: 2,
					colIdx: 2,
				},
				spellIds: [19485, 19487],
				maxPoints: 5,
			},
			{
				//fieldName: 'concussiveBarrage',
				location: {
					rowIdx: 4,
					colIdx: 0,
				},
				spellIds: [35100, 35102],
				maxPoints: 3,
			},
			{
				fieldName: 'scatterShot',
				location: {
					rowIdx: 4,
					colIdx: 1,
				},
				spellIds: [19503],
				maxPoints: 1,
			},
			{
				fieldName: 'barrage',
				location: {
					rowIdx: 4,
					colIdx: 2,
				},
				spellIds: [19461, 19462, 24691],
				maxPoints: 3,
			},
			{
				fieldName: 'combatExperience',
				location: {
					rowIdx: 5,
					colIdx: 0,
				},
				spellIds: [34475],
				maxPoints: 2,
			},
			{
				fieldName: 'rangedWeaponSpecialization',
				location: {
					rowIdx: 5,
					colIdx: 3,
				},
				spellIds: [19507],
				maxPoints: 5,
			},
			{
				fieldName: 'carefulAim',
				location: {
					rowIdx: 6,
					colIdx: 0,
				},
				spellIds: [34482],
				maxPoints: 3,
			},
			{
				fieldName: 'trueshotAura',
				location: {
					rowIdx: 6,
					colIdx: 1,
				},
				prereqLocation: {
					rowIdx: 4,
					colIdx: 1,
				},
				spellIds: [19506],
				maxPoints: 1,
			},
			{
				fieldName: 'improvedBarrage',
				location: {
					rowIdx: 6,
					colIdx: 2,
				},
				spellIds: [35104, 35110],
				maxPoints: 3,
			},
			{
				fieldName: 'masterMarksman',
				location: {
					rowIdx: 7,
					colIdx: 1,
				},
				spellIds: [34485],
				maxPoints: 5,
			},
			{
				fieldName: 'silencingShot',
				location: {
					rowIdx: 8,
					colIdx: 1,
				},
				prereqLocation: {
					rowIdx: 7,
					colIdx: 1,
				},
				spellIds: [34490],
				maxPoints: 1,
			},
		],
	},
	{
		name: 'Survival',
		backgroundUrl: 'https://wow.zamimg.com/images/wow/talents/backgrounds/tbc/362.jpg',
		talents: [
			{
				fieldName: 'monsterSlaying',
				location: {
					rowIdx: 0,
					colIdx: 0,
				},
				spellIds: [24293],
				maxPoints: 3,
			},
			{
				fieldName: 'humanoidSlaying',
				location: {
					rowIdx: 0,
					colIdx: 1,
				},
				spellIds: [19151],
				maxPoints: 3,
			},
			{
				//fieldName: 'hawkEye',
				location: {
					rowIdx: 0,
					colIdx: 2,
				},
				spellIds: [19498],
				maxPoints: 3,
			},
			{
				fieldName: 'savageStrikes',
				location: {
					rowIdx: 0,
					colIdx: 3,
				},
				spellIds: [19159],
				maxPoints: 2,
			},
			{
				//fieldName: 'entrapment',
				location: {
					rowIdx: 1,
					colIdx: 0,
				},
				spellIds: [19184, 19387],
				maxPoints: 3,
			},
			{
				fieldName: 'deflection',
				location: {
					rowIdx: 1,
					colIdx: 1,
				},
				spellIds: [19295, 19297, 19298, 19301, 19300],
				maxPoints: 5,
			},
			{
				//fieldName: 'improvedWingClip',
				location: {
					rowIdx: 1,
					colIdx: 2,
				},
				spellIds: [19228, 19232],
				maxPoints: 3,
			},
			{
				fieldName: 'cleverTraps',
				location: {
					rowIdx: 2,
					colIdx: 0,
				},
				spellIds: [19239, 19245],
				maxPoints: 2,
			},
			{
				fieldName: 'survivalist',
				location: {
					rowIdx: 2,
					colIdx: 1,
				},
				spellIds: [19255],
				maxPoints: 5,
			},
			{
				//fieldName: 'deterrance',
				location: {
					rowIdx: 2,
					colIdx: 2,
				},
				spellIds: [19263],
				maxPoints: 1,
			},
			{
				fieldName: 'trapMastery',
				location: {
					rowIdx: 3,
					colIdx: 0,
				},
				spellIds: [19376],
				maxPoints: 2,
			},
			{
				fieldName: 'surefooted',
				location: {
					rowIdx: 3,
					colIdx: 1,
				},
				spellIds: [19290, 19294, 24283],
				maxPoints: 3,
			},
			{
				//fieldName: 'improvedFeignDeath',
				location: {
					rowIdx: 3,
					colIdx: 3,
				},
				spellIds: [19286],
				maxPoints: 2,
			},
			{
				fieldName: 'survivalInstincts',
				location: {
					rowIdx: 4,
					colIdx: 0,
				},
				spellIds: [34494, 34496],
				maxPoints: 2,
			},
			{
				fieldName: 'killerInstinct',
				location: {
					rowIdx: 4,
					colIdx: 1,
				},
				spellIds: [19370, 19371, 19373],
				maxPoints: 3,
			},
			{
				//fieldName: 'counterattack',
				location: {
					rowIdx: 4,
					colIdx: 2,
				},
				prereqLocation: {
					rowIdx: 2,
					colIdx: 2,
				},
				spellIds: [19306],
				maxPoints: 1,
			},
			{
				fieldName: 'resourcefulness',
				location: {
					rowIdx: 5,
					colIdx: 0,
				},
				spellIds: [34491],
				maxPoints: 3,
			},
			{
				fieldName: 'lightningReflexes',
				location: {
					rowIdx: 5,
					colIdx: 2,
				},
				spellIds: [19168, 19180, 19181, 24296],
				maxPoints: 5,
			},
			{
				fieldName: 'thrillOfTheHunt',
				location: {
					rowIdx: 6,
					colIdx: 0,
				},
				spellIds: [34497],
				maxPoints: 3,
			},
			{
				//fieldName: 'wyvernSting',
				location: {
					rowIdx: 6,
					colIdx: 1,
				},
				prereqLocation: {
					rowIdx: 4,
					colIdx: 1,
				},
				spellIds: [19386],
				maxPoints: 1,
			},
			{
				fieldName: 'exposeWeakness',
				location: {
					rowIdx: 6,
					colIdx: 2,
				},
				prereqLocation: {
					rowIdx: 5,
					colIdx: 2,
				},
				spellIds: [34500, 34502],
				maxPoints: 3,
			},
			{
				fieldName: 'masterTactician',
				location: {
					rowIdx: 7,
					colIdx: 1,
				},
				spellIds: [34506, 34507, 34508, 34838],
				maxPoints: 5,
			},
			{
				fieldName: 'readiness',
				location: {
					rowIdx: 8,
					colIdx: 1,
				},
				prereqLocation: {
					rowIdx: 7,
					colIdx: 1,
				},
				spellIds: [23989],
				maxPoints: 1,
			},
		],
	},
]);
