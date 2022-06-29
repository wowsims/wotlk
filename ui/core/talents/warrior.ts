import { Spec } from '/tbc/core/proto/common.js';
import { WarriorTalents as WarriorTalents } from '/tbc/core/proto/warrior.js';
import { WarriorSpecs } from '/tbc/core/proto_utils/utils.js';
import { Player } from '/tbc/core/player.js';

import { TalentsConfig, TalentsPicker, newTalentsConfig } from './talents_picker.js';

export class WarriorTalentsPicker extends TalentsPicker<Spec.SpecWarrior> {
	constructor(parent: HTMLElement, player: Player<Spec.SpecWarrior>) {
		super(parent, player, warriorTalentsConfig);
	}
}

export const warriorTalentsConfig: TalentsConfig<Spec.SpecWarrior> = newTalentsConfig([
	{
		name: 'Arms',
		backgroundUrl: 'https://wow.zamimg.com/images/wow/talents/backgrounds/tbc/161.jpg',
		talents: [
			{
				fieldName: 'improvedHeroicStrike',
				location: {
					rowIdx: 0,
					colIdx: 0,
				},
				spellIds: [12282, 12663],
				maxPoints: 3,
			},
			{
				fieldName: 'deflection',
				location: {
					rowIdx: 0,
					colIdx: 1,
				},
				spellIds: [16462],
				maxPoints: 5,
			},
			{
				fieldName: 'improvedRend',
				location: {
					rowIdx: 0,
					colIdx: 2,
				},
				spellIds: [12286, 12658],
				maxPoints: 3,
			},
			{
				fieldName: 'improvedCharge',
				location: {
					rowIdx: 1,
					colIdx: 0,
				},
				spellIds: [12285, 12697],
				maxPoints: 2,
			},
			{
				//fieldName: 'ironWill',
				location: {
					rowIdx: 1,
					colIdx: 1,
				},
				spellIds: [12300, 12959],
				maxPoints: 5,
			},
			{
				fieldName: 'improvedThunderClap',
				location: {
					rowIdx: 1,
					colIdx: 2,
				},
				spellIds: [12287, 12665],
				maxPoints: 3,
			},
			{
				fieldName: 'improvedOverpower',
				location: {
					rowIdx: 2,
					colIdx: 0,
				},
				spellIds: [12290, 12963],
				maxPoints: 2,
			},
			{
				fieldName: 'angerManagement',
				location: {
					rowIdx: 2,
					colIdx: 1,
				},
				spellIds: [12296],
				maxPoints: 1,
			},
			{
				fieldName: 'deepWounds',
				location: {
					rowIdx: 2,
					colIdx: 2,
				},
				spellIds: [12834, 12849, 12867],
				maxPoints: 3,
			},
			{
				fieldName: 'twoHandedWeaponSpecialization',
				location: {
					rowIdx: 3,
					colIdx: 1,
				},
				spellIds: [12163, 12711],
				maxPoints: 5,
			},
			{
				fieldName: 'impale',
				location: {
					rowIdx: 3,
					colIdx: 2,
				},
				prereqLocation: {
					rowIdx: 2,
					colIdx: 2,
				},
				spellIds: [16493],
				maxPoints: 2,
			},
			{
				fieldName: 'poleaxeSpecialization',
				location: {
					rowIdx: 4,
					colIdx: 0,
				},
				spellIds: [12700, 12781, 12783],
				maxPoints: 5,
			},
			{
				fieldName: 'deathWish',
				location: {
					rowIdx: 4,
					colIdx: 1,
				},
				spellIds: [12292],
				maxPoints: 1,
			},
			{
				fieldName: 'maceSpecialization',
				location: {
					rowIdx: 4,
					colIdx: 2,
				},
				spellIds: [12284, 12701],
				maxPoints: 5,
			},
			{
				fieldName: 'swordSpecialization',
				location: {
					rowIdx: 4,
					colIdx: 3,
				},
				spellIds: [12281, 12812],
				maxPoints: 5,
			},
			{
				//fieldName: 'improvedIntercept',
				location: {
					rowIdx: 5,
					colIdx: 0,
				},
				spellIds: [29888],
				maxPoints: 2,
			},
			{
				//fieldName: 'improvedHamstring',
				location: {
					rowIdx: 5,
					colIdx: 2,
				},
				spellIds: [12289, 12668, 23695],
				maxPoints: 3,
			},
			{
				fieldName: 'improvedDisciplines',
				location: {
					rowIdx: 5,
					colIdx: 3,
				},
				spellIds: [29723],
				maxPoints: 3,
			},
			{
				fieldName: 'bloodFrenzy',
				location: {
					rowIdx: 6,
					colIdx: 0,
				},
				spellIds: [29836, 29859],
				maxPoints: 2,
			},
			{
				fieldName: 'mortalStrike',
				location: {
					rowIdx: 6,
					colIdx: 1,
				},
				prereqLocation: {
					rowIdx: 4,
					colIdx: 1,
				},
				spellIds: [12294],
				maxPoints: 1,
			},
			{
				//fieldName: 'secondWind',
				location: {
					rowIdx: 6,
					colIdx: 2,
				},
				spellIds: [29834, 29838],
				maxPoints: 2,
			},
			{
				fieldName: 'improvedMortalStrike',
				location: {
					rowIdx: 7,
					colIdx: 1,
				},
				prereqLocation: {
					rowIdx: 6,
					colIdx: 1,
				},
				spellIds: [35446, 35448],
				maxPoints: 5,
			},
			{
				fieldName: 'endlessRage',
				location: {
					rowIdx: 8,
					colIdx: 1,
				},
				spellIds: [29623],
				maxPoints: 1,
			},
		],
	},
	{
		name: 'Fury',
		backgroundUrl: 'https://wow.zamimg.com/images/wow/talents/backgrounds/tbc/164.jpg',
		talents: [
			{
				fieldName: 'boomingVoice',
				location: {
					rowIdx: 0,
					colIdx: 1,
				},
				spellIds: [12321, 12835],
				maxPoints: 5,
			},
			{
				fieldName: 'cruelty',
				location: {
					rowIdx: 0,
					colIdx: 2,
				},
				spellIds: [12320, 12852, 12853, 12855],
				maxPoints: 5,
			},
			{
				fieldName: 'improvedDemoralizingShout',
				location: {
					rowIdx: 1,
					colIdx: 1,
				},
				spellIds: [12324, 12876],
				maxPoints: 5,
			},
			{
				fieldName: 'unbridledWrath',
				location: {
					rowIdx: 1,
					colIdx: 2,
				},
				spellIds: [12322, 12999],
				maxPoints: 5,
			},
			{
				fieldName: 'improvedCleave',
				location: {
					rowIdx: 2,
					colIdx: 0,
				},
				spellIds: [12329, 12950, 20496],
				maxPoints: 3,
			},
			{
				//fieldName: 'piercingHowl',
				location: {
					rowIdx: 2,
					colIdx: 1,
				},
				spellIds: [12323],
				maxPoints: 1,
			},
			{
				//fieldName: 'bloodCraze',
				location: {
					rowIdx: 2,
					colIdx: 2,
				},
				spellIds: [16487, 16489, 16492],
				maxPoints: 3,
			},
			{
				fieldName: 'commandingPresence',
				location: {
					rowIdx: 2,
					colIdx: 3,
				},
				spellIds: [12318, 12857, 12858, 12860],
				maxPoints: 5,
			},
			{
				fieldName: 'dualWieldSpecialization',
				location: {
					rowIdx: 3,
					colIdx: 0,
				},
				spellIds: [23584],
				maxPoints: 5,
			},
			{
				fieldName: 'improvedExecute',
				location: {
					rowIdx: 3,
					colIdx: 1,
				},
				spellIds: [20502],
				maxPoints: 2,
			},
			{
				//fieldName: 'enrage',
				location: {
					rowIdx: 3,
					colIdx: 2,
				},
				spellIds: [12317, 13045],
				maxPoints: 5,
			},
			{
				fieldName: 'improvedSlam',
				location: {
					rowIdx: 4,
					colIdx: 0,
				},
				spellIds: [12862, 12330],
				maxPoints: 2,
			},
			{
				fieldName: 'sweepingStrikes',
				location: {
					rowIdx: 4,
					colIdx: 1,
				},
				spellIds: [12328],
				maxPoints: 1,
			},
			{
				fieldName: 'weaponMastery',
				location: {
					rowIdx: 4,
					colIdx: 3,
				},
				spellIds: [20504],
				maxPoints: 2,
			},
			{
				fieldName: 'improvedBerserkerRage',
				location: {
					rowIdx: 5,
					colIdx: 0,
				},
				spellIds: [20500],
				maxPoints: 2,
			},
			{
				fieldName: 'flurry',
				location: {
					rowIdx: 5,
					colIdx: 2,
				},
				prereqLocation: {
					rowIdx: 3,
					colIdx: 2,
				},
				spellIds: [12319, 12971],
				maxPoints: 5,
			},
			{
				fieldName: 'precision',
				location: {
					rowIdx: 6,
					colIdx: 0,
				},
				spellIds: [29590],
				maxPoints: 3,
			},
			{
				fieldName: 'bloodthirst',
				location: {
					rowIdx: 6,
					colIdx: 1,
				},
				prereqLocation: {
					rowIdx: 4,
					colIdx: 1,
				},
				spellIds: [23881],
				maxPoints: 1,
			},
			{
				fieldName: 'improvedWhirlwind',
				location: {
					rowIdx: 6,
					colIdx: 2,
				},
				spellIds: [29721, 29776],
				maxPoints: 2,
			},
			{
				fieldName: 'improvedBerserkerStance',
				location: {
					rowIdx: 7,
					colIdx: 2,
				},
				spellIds: [29759],
				maxPoints: 5,
			},
			{
				fieldName: 'rampage',
				location: {
					rowIdx: 8,
					colIdx: 1,
				},
				prereqLocation: {
					rowIdx: 6,
					colIdx: 1,
				},
				spellIds: [29801],
				maxPoints: 1,
			},
		],
	},
	{
		name: 'Protection',
		backgroundUrl: 'https://wow.zamimg.com/images/wow/talents/backgrounds/tbc/163.jpg',
		talents: [
			{
				fieldName: 'improvedBloodrage',
				location: {
					rowIdx: 0,
					colIdx: 0,
				},
				spellIds: [12301, 12818],
				maxPoints: 2,
			},
			{
				fieldName: 'tacticalMastery',
				location: {
					rowIdx: 0,
					colIdx: 1,
				},
				spellIds: [12295, 12676],
				maxPoints: 3,
			},
			{
				fieldName: 'anticipation',
				location: {
					rowIdx: 0,
					colIdx: 2,
				},
				spellIds: [12297, 12750],
				maxPoints: 5,
			},
			{
				fieldName: 'shieldSpecialization',
				location: {
					rowIdx: 1,
					colIdx: 1,
				},
				spellIds: [12298, 12724],
				maxPoints: 5,
			},
			{
				fieldName: 'toughness',
				location: {
					rowIdx: 1,
					colIdx: 2,
				},
				spellIds: [12299, 12761],
				maxPoints: 5,
			},
			{
				fieldName: 'lastStand',
				location: {
					rowIdx: 2,
					colIdx: 0,
				},
				spellIds: [12975],
				maxPoints: 1,
			},
			{
				fieldName: 'improvedShieldBlock',
				location: {
					rowIdx: 2,
					colIdx: 1,
				},
				prereqLocation: {
					rowIdx: 1,
					colIdx: 1,
				},
				spellIds: [12945],
				maxPoints: 1,
			},
			{
				//fieldName: 'improvedRevenge',
				location: {
					rowIdx: 2,
					colIdx: 2,
				},
				spellIds: [12797, 12799],
				maxPoints: 3,
			},
			{
				fieldName: 'defiance',
				location: {
					rowIdx: 2,
					colIdx: 3,
				},
				spellIds: [12303, 12788],
				maxPoints: 3,
			},
			{
				fieldName: 'improvedSunderArmor',
				location: {
					rowIdx: 3,
					colIdx: 0,
				},
				spellIds: [12308, 12810],
				maxPoints: 3,
			},
			{
				//fieldName: 'improvedDisarm',
				location: {
					rowIdx: 3,
					colIdx: 1,
				},
				spellIds: [12313, 12804, 12807],
				maxPoints: 3,
			},
			{
				//fieldName: 'improvedTaunt',
				location: {
					rowIdx: 3,
					colIdx: 2,
				},
				spellIds: [12302, 12765],
				maxPoints: 2,
			},
			{
				fieldName: 'improvedShieldWall',
				location: {
					rowIdx: 4,
					colIdx: 0,
				},
				spellIds: [12312, 12803],
				maxPoints: 2,
			},
			{
				//fieldName: 'concussionBlow',
				location: {
					rowIdx: 4,
					colIdx: 1,
				},
				spellIds: [12809],
				maxPoints: 1,
			},
			{
				//fieldName: 'improvedShieldBash',
				location: {
					rowIdx: 4,
					colIdx: 2,
				},
				spellIds: [12311, 12958],
				maxPoints: 2,
			},
			{
				fieldName: 'shieldMastery',
				location: {
					rowIdx: 5,
					colIdx: 0,
				},
				spellIds: [29598],
				maxPoints: 3,
			},
			{
				fieldName: 'oneHandedWeaponSpecialization',
				location: {
					rowIdx: 5,
					colIdx: 2,
				},
				spellIds: [16538],
				maxPoints: 5,
			},
			{
				fieldName: 'improvedDefensiveStance',
				location: {
					rowIdx: 6,
					colIdx: 0,
				},
				spellIds: [29593],
				maxPoints: 3,
			},
			{
				fieldName: 'shieldSlam',
				location: {
					rowIdx: 6,
					colIdx: 1,
				},
				prereqLocation: {
					rowIdx: 4,
					colIdx: 1,
				},
				spellIds: [23922],
				maxPoints: 1,
			},
			{
				fieldName: 'focusedRage',
				location: {
					rowIdx: 6,
					colIdx: 2,
				},
				spellIds: [29787, 29790, 29792],
				maxPoints: 3,
			},
			{
				fieldName: 'vitality',
				location: {
					rowIdx: 7,
					colIdx: 1,
				},
				spellIds: [29140, 29143],
				maxPoints: 5,
			},
			{
				fieldName: 'devastate',
				location: {
					rowIdx: 8,
					colIdx: 1,
				},
				spellIds: [20243],
				maxPoints: 1,
			},
		],
	},
]);
