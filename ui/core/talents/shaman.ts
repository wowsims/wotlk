import { Spec } from '/tbc/core/proto/common.js';
import { ShamanTalents as ShamanTalents } from '/tbc/core/proto/shaman.js';
import { Player } from '/tbc/core/player.js';

import { TalentsConfig, TalentsPicker, newTalentsConfig } from './talents_picker.js';

// Talents are the same for all Shaman specs, so its ok to just use ElementalShaman here
export class ShamanTalentsPicker extends TalentsPicker<Spec.SpecElementalShaman> {
	constructor(parent: HTMLElement, player: Player<Spec.SpecElementalShaman>) {
		super(parent, player, shamanTalentsConfig);
	}
}

export const shamanTalentsConfig: TalentsConfig<Spec.SpecElementalShaman> = newTalentsConfig([
	{
		name: 'Elemental',
		backgroundUrl: 'https://wow.zamimg.com/images/wow/talents/backgrounds/classic/261.jpg',
		talents: [
			{
				fieldName: 'convection',
				location: {
					rowIdx: 0,
					colIdx: 1,
				},
				spellIds: [16039, 16109],
				maxPoints: 5,
			},
			{
				fieldName: 'concussion',
				location: {
					rowIdx: 0,
					colIdx: 2,
				},
				spellIds: [16035, 16105],
				maxPoints: 5,
			},
			{
				//fieldName: 'earthsGrasp',
				location: {
					rowIdx: 1,
					colIdx: 0,
				},
				spellIds: [16043, 16130],
				maxPoints: 2,
			},
			{
				//fieldName: 'elementalWarding',
				location: {
					rowIdx: 1,
					colIdx: 1,
				},
				spellIds: [28996],
				maxPoints: 3,
			},
			{
				fieldName: 'callOfFlame',
				location: {
					rowIdx: 1,
					colIdx: 2,
				},
				spellIds: [16038, 16160],
				maxPoints: 3,
			},
			{
				fieldName: 'elementalFocus',
				location: {
					rowIdx: 2,
					colIdx: 0,
				},
				spellIds: [16164],
				maxPoints: 1,
			},
			{
				fieldName: 'reverberation',
				location: {
					rowIdx: 2,
					colIdx: 1,
				},
				spellIds: [16040, 16113],
				maxPoints: 5,
			},
			{
				fieldName: 'callOfThunder',
				location: {
					rowIdx: 2,
					colIdx: 2,
				},
				spellIds: [16041, 16117],
				maxPoints: 5,
			},
			{
				fieldName: 'improvedFireTotems',
				location: {
					rowIdx: 3,
					colIdx: 0,
				},
				spellIds: [16086, 16544],
				maxPoints: 2,
			},
			{
				//fieldName: 'eyeOfTheStorm',
				location: {
					rowIdx: 3,
					colIdx: 1,
				},
				spellIds: [29062, 29064],
				maxPoints: 3,
			},
			{
				fieldName: 'elementalDevastation',
				location: {
					rowIdx: 3,
					colIdx: 3,
				},
				spellIds: [30160, 29179],
				maxPoints: 3,
			},
			{
				//fieldName: 'stormReach',
				location: {
					rowIdx: 4,
					colIdx: 0,
				},
				spellIds: [28999],
				maxPoints: 2,
			},
			{
				fieldName: 'elementalFury',
				location: {
					rowIdx: 4,
					colIdx: 1,
				},
				spellIds: [16089],
				maxPoints: 1,
			},
			{
				fieldName: 'unrelentingStorm',
				location: {
					rowIdx: 4,
					colIdx: 3,
				},
				spellIds: [30664],
				maxPoints: 5,
			},
			{
				fieldName: 'elementalPrecision',
				location: {
					rowIdx: 5,
					colIdx: 0,
				},
				spellIds: [30672],
				maxPoints: 3,
			},
			{
				fieldName: 'lightningMastery',
				location: {
					rowIdx: 5,
					colIdx: 2,
				},
				prereqLocation: {
					rowIdx: 2,
					colIdx: 2,
				},
				spellIds: [16578],
				maxPoints: 5,
			},
			{
				fieldName: 'elementalMastery',
				location: {
					rowIdx: 6,
					colIdx: 1,
				},
				prereqLocation: {
					rowIdx: 4,
					colIdx: 1,
				},
				spellIds: [16166],
				maxPoints: 1,
			},
			{
				//fieldName: 'elementalShields',
				location: {
					rowIdx: 6,
					colIdx: 2,
				},
				spellIds: [30669],
				maxPoints: 3,
			},
			{
				fieldName: 'lightningOverload',
				location: {
					rowIdx: 7,
					colIdx: 1,
				},
				spellIds: [30675, 30678],
				maxPoints: 5,
			},
			{
				fieldName: 'totemOfWrath',
				location: {
					rowIdx: 8,
					colIdx: 1,
				},
				prereqLocation: {
					rowIdx: 7,
					colIdx: 1,
				},
				spellIds: [30706],
				maxPoints: 1,
			},
		],
	},
	{
		name: 'Enhancement',
		backgroundUrl: 'https://wow.zamimg.com/images/wow/talents/backgrounds/classic/263.jpg',
		talents: [
			{
				fieldName: 'ancestralKnowledge',
				location: {
					rowIdx: 0,
					colIdx: 1,
				},
				spellIds: [17485],
				maxPoints: 5,
			},
			{
				fieldName: 'shieldSpecialization',
				location: {
					rowIdx: 0,
					colIdx: 2,
				},
				spellIds: [16253, 16298],
				maxPoints: 5,
			},
			{
				//fieldName: 'guardianTotems',
				location: {
					rowIdx: 1,
					colIdx: 0,
				},
				spellIds: [16258, 16293],
				maxPoints: 2,
			},
			{
				fieldName: 'thunderingStrikes',
				location: {
					rowIdx: 1,
					colIdx: 1,
				},
				spellIds: [16255, 16302],
				maxPoints: 5,
			},
			{
				//fieldName: 'improvedGhostWolf',
				location: {
					rowIdx: 1,
					colIdx: 2,
				},
				spellIds: [16262, 16287],
				maxPoints: 2,
			},
			{
				//fieldName: 'improvedLightningShield',
				location: {
					rowIdx: 1,
					colIdx: 3,
				},
				spellIds: [16261, 16290],
				maxPoints: 3,
			},
			{
				fieldName: 'enhancingTotems',
				location: {
					rowIdx: 2,
					colIdx: 0,
				},
				spellIds: [16259, 16295],
				maxPoints: 2,
			},
			{
				fieldName: 'shamanisticFocus',
				location: {
					rowIdx: 2,
					colIdx: 2,
				},
				spellIds: [43338],
				maxPoints: 1,
			},
			{
				fieldName: 'anticipation',
				location: {
					rowIdx: 2,
					colIdx: 3,
				},
				spellIds: [16254, 16271],
				maxPoints: 5,
			},
			{
				fieldName: 'flurry',
				location: {
					rowIdx: 3,
					colIdx: 1,
				},
				prereqLocation: {
					rowIdx: 1,
					colIdx: 1,
				},
				spellIds: [16256, 16281],
				maxPoints: 5,
			},
			{
				fieldName: 'toughness',
				location: {
					rowIdx: 3,
					colIdx: 2,
				},
				spellIds: [16252, 16306],
				maxPoints: 5,
			},
			{
				fieldName: 'improvedWeaponTotems',
				location: {
					rowIdx: 4,
					colIdx: 0,
				},
				spellIds: [29192],
				maxPoints: 2,
			},
			{
				fieldName: 'spiritWeapons',
				location: {
					rowIdx: 4,
					colIdx: 1,
				},
				spellIds: [16268],
				maxPoints: 1,
			},
			{
				fieldName: 'elementalWeapons',
				location: {
					rowIdx: 4,
					colIdx: 2,
				},
				spellIds: [16266, 29079],
				maxPoints: 3,
			},
			{
				fieldName: 'mentalQuickness',
				location: {
					rowIdx: 5,
					colIdx: 0,
				},
				spellIds: [30812],
				maxPoints: 3,
			},
			{
				fieldName: 'weaponMastery',
				location: {
					rowIdx: 5,
					colIdx: 3,
				},
				spellIds: [29082, 29084, 29086],
				maxPoints: 5,
			},
			{
				fieldName: 'dualWieldSpecialization',
				location: {
					rowIdx: 6,
					colIdx: 0,
				},
				prereqLocation: {
					rowIdx: 6,
					colIdx: 1,
				},
				spellIds: [30816, 30818],
				maxPoints: 3,
			},
			{
				//fieldName: 'dualWield',
				location: {
					rowIdx: 6,
					colIdx: 1,
				},
				prereqLocation: {
					rowIdx: 4,
					colIdx: 1,
				},
				spellIds: [30798],
				maxPoints: 1,
			},
			{
				fieldName: 'stormstrike',
				location: {
					rowIdx: 6,
					colIdx: 2,
				},
				prereqLocation: {
					rowIdx: 4,
					colIdx: 2,
				},
				spellIds: [17364],
				maxPoints: 1,
			},
			{
				fieldName: 'unleashedRage',
				location: {
					rowIdx: 7,
					colIdx: 1,
				},
				spellIds: [30802, 30808],
				maxPoints: 5,
			},
			{
				fieldName: 'shamanisticRage',
				location: {
					rowIdx: 8,
					colIdx: 1,
				},
				spellIds: [30823],
				maxPoints: 1,
			},
		],
	},
	{
		name: 'Restoration',
		backgroundUrl: 'https://wow.zamimg.com/images/wow/talents/backgrounds/classic/262.jpg',
		talents: [
			{
				//fieldName: 'improvedHealingWave',
				location: {
					rowIdx: 0,
					colIdx: 1,
				},
				spellIds: [16182, 16226],
				maxPoints: 5,
			},
			{
				//fieldName: 'tidalFocus',
				location: {
					rowIdx: 0,
					colIdx: 2,
				},
				spellIds: [16179, 16214],
				maxPoints: 5,
			},
			{
				//fieldName: 'improvedReincarnation',
				location: {
					rowIdx: 1,
					colIdx: 0,
				},
				spellIds: [16184, 16209],
				maxPoints: 2,
			},
			{
				//fieldName: 'ancestralUealing',
				location: {
					rowIdx: 1,
					colIdx: 1,
				},
				spellIds: [16176, 16235],
				maxPoints: 3,
			},
			{
				fieldName: 'totemicFocus',
				location: {
					rowIdx: 1,
					colIdx: 2,
				},
				spellIds: [16173, 16222],
				maxPoints: 5,
			},
			{
				fieldName: 'naturesGuidance',
				location: {
					rowIdx: 2,
					colIdx: 0,
				},
				spellIds: [16180, 16196, 16198],
				maxPoints: 3,
			},
			{
				//fieldName: 'healingFocus',
				location: {
					rowIdx: 2,
					colIdx: 1,
				},
				spellIds: [16181, 16230, 16232],
				maxPoints: 5,
			},
			{
				//fieldName: 'totemicMastery',
				location: {
					rowIdx: 2,
					colIdx: 2,
				},
				spellIds: [16189],
				maxPoints: 1,
			},
			{
				//fieldName: 'healingGrace',
				location: {
					rowIdx: 2,
					colIdx: 3,
				},
				spellIds: [29187, 29189, 29191],
				maxPoints: 3,
			},
			{
				fieldName: 'restorativeTotems',
				location: {
					rowIdx: 3,
					colIdx: 1,
				},
				spellIds: [16187, 16205],
				maxPoints: 5,
			},
			{
				fieldName: 'tidalMastery',
				location: {
					rowIdx: 3,
					colIdx: 2,
				},
				spellIds: [16194, 16218],
				maxPoints: 5,
			},
			{
				//fieldName: 'healingWay',
				location: {
					rowIdx: 4,
					colIdx: 0,
				},
				spellIds: [29206, 29205, 29202],
				maxPoints: 3,
			},
			{
				fieldName: 'naturesSwiftness',
				location: {
					rowIdx: 4,
					colIdx: 2,
				},
				spellIds: [16188],
				maxPoints: 1,
			},
			{
				//fieldName: 'focusedMind',
				location: {
					rowIdx: 4,
					colIdx: 3,
				},
				spellIds: [30864],
				maxPoints: 3,
			},
			{
				//fieldName: 'purification',
				location: {
					rowIdx: 5,
					colIdx: 2,
				},
				spellIds: [16178, 16210],
				maxPoints: 5,
			},
			{
				//fieldName: 'manaTideTotem',
				location: {
					rowIdx: 6,
					colIdx: 1,
				},
				prereqLocation: {
					rowIdx: 3,
					colIdx: 1,
				},
				spellIds: [16190],
				maxPoints: 1,
			},
			{
				//fieldName: 'naturesGuardian',
				location: {
					rowIdx: 6,
					colIdx: 2,
				},
				spellIds: [30881, 30883],
				maxPoints: 5,
			},
			{
				fieldName: 'naturesBlessing',
				location: {
					rowIdx: 7,
					colIdx: 1,
				},
				spellIds: [30867],
				maxPoints: 3,
			},
			{
				//fieldName: 'improvedChainHeal',
				location: {
					rowIdx: 7,
					colIdx: 2,
				},
				spellIds: [30872],
				maxPoints: 2,
			},
			{
				//fieldName: 'earthShield',
				location: {
					rowIdx: 8,
					colIdx: 1,
				},
				prereqLocation: {
					rowIdx: 7,
					colIdx: 1,
				},
				spellIds: [974],
				maxPoints: 1,
			},
		],
	},
]);
