import { Spec } from '/tbc/core/proto/common.js';
import { MageTalents as MageTalents } from '/tbc/core/proto/mage.js';
import { MageSpecs } from '/tbc/core/proto_utils/utils.js';
import { Player } from '/tbc/core/player.js';

import { TalentsConfig, TalentsPicker, newTalentsConfig } from './talents_picker.js';

export class MageTalentsPicker extends TalentsPicker<Spec.SpecMage> {
	constructor(parent: HTMLElement, player: Player<Spec.SpecMage>) {
		super(parent, player, mageTalentsConfig);
	}
}

export const mageTalentsConfig: TalentsConfig<Spec.SpecMage> = newTalentsConfig([
	{
		name: 'Arcane',
		backgroundUrl: 'https://wow.zamimg.com/images/wow/talents/backgrounds/tbc/81.jpg',
		talents: [
			{
				fieldName: 'arcaneSubtlety',
				location: {
					rowIdx: 0,
					colIdx: 0,
				},
				spellIds: [11210, 12592],
				maxPoints: 2,
			},
			{
				fieldName: 'arcaneFocus',
				location: {
					rowIdx: 0,
					colIdx: 1,
				},
				spellIds: [11222, 12839],
				maxPoints: 5,
			},
			{
				//fieldName: 'improvedArcaneMissiles',
				location: {
					rowIdx: 0,
					colIdx: 2,
				},
				spellIds: [11237, 12463, 12464, 16769],
				maxPoints: 5,
			},
			{
				fieldName: 'wandSpecialization',
				location: {
					rowIdx: 1,
					colIdx: 0,
				},
				spellIds: [6057, 6085],
				maxPoints: 2,
			},
			{
				fieldName: 'magicAbsorption',
				location: {
					rowIdx: 1,
					colIdx: 1,
				},
				spellIds: [29441, 29444],
				maxPoints: 5,
			},
			{
				fieldName: 'arcaneConcentration',
				location: {
					rowIdx: 1,
					colIdx: 2,
				},
				spellIds: [11213, 12574],
				maxPoints: 5,
			},
			{
				//fieldName: 'magicAttunement',
				location: {
					rowIdx: 2,
					colIdx: 0,
				},
				spellIds: [11247, 12606],
				maxPoints: 2,
			},
			{
				fieldName: 'arcaneImpact',
				location: {
					rowIdx: 2,
					colIdx: 1,
				},
				spellIds: [11242, 12467, 12469],
				maxPoints: 3,
			},
			{
				//fieldName: 'arcaneFortitude',
				location: {
					rowIdx: 2,
					colIdx: 3,
				},
				spellIds: [28574],
				maxPoints: 1,
			},
			{
				//fieldName: 'improvedManaShield',
				location: {
					rowIdx: 3,
					colIdx: 0,
				},
				spellIds: [11252, 12605],
				maxPoints: 2,
			},
			{
				//fieldName: 'improvedCounterspell',
				location: {
					rowIdx: 3,
					colIdx: 1,
				},
				spellIds: [11255, 12598],
				maxPoints: 2,
			},
			{
				fieldName: 'arcaneMeditation',
				location: {
					rowIdx: 3,
					colIdx: 3,
				},
				spellIds: [18462],
				maxPoints: 3,
			},
			{
				//fieldName: 'improvedBlink',
				location: {
					rowIdx: 4,
					colIdx: 0,
				},
				spellIds: [31569],
				maxPoints: 2,
			},
			{
				fieldName: 'presenceOfMind',
				location: {
					rowIdx: 4,
					colIdx: 1,
				},
				spellIds: [12043],
				maxPoints: 1,
			},
			{
				fieldName: 'arcaneMind',
				location: {
					rowIdx: 4,
					colIdx: 3,
				},
				spellIds: [11232, 12500],
				maxPoints: 5,
			},
			{
				//fieldName: 'prismaticCloak',
				location: {
					rowIdx: 5,
					colIdx: 0,
				},
				spellIds: [31574],
				maxPoints: 2,
			},
			{
				fieldName: 'arcaneInstability',
				location: {
					rowIdx: 5,
					colIdx: 1,
				},
				prereqLocation: {
					rowIdx: 4,
					colIdx: 1,
				},
				spellIds: [15058],
				maxPoints: 3,
			},
			{
				fieldName: 'arcanePotency',
				location: {
					rowIdx: 5,
					colIdx: 2,
				},
				prereqLocation: {
					rowIdx: 1,
					colIdx: 2,
				},
				spellIds: [31571],
				maxPoints: 3,
			},
			{
				fieldName: 'empoweredArcaneMissiles',
				location: {
					rowIdx: 6,
					colIdx: 0,
				},
				spellIds: [31579, 31582],
				maxPoints: 3,
			},
			{
				fieldName: 'arcanePower',
				location: {
					rowIdx: 6,
					colIdx: 1,
				},
				prereqLocation: {
					rowIdx: 5,
					colIdx: 1,
				},
				spellIds: [12042],
				maxPoints: 1,
			},
			{
				fieldName: 'spellPower',
				location: {
					rowIdx: 6,
					colIdx: 2,
				},
				spellIds: [35578, 35581],
				maxPoints: 2,
			},
			{
				fieldName: 'mindMastery',
				location: {
					rowIdx: 7,
					colIdx: 1,
				},
				spellIds: [31584],
				maxPoints: 5,
			},
			{
				//fieldName: 'slow',
				location: {
					rowIdx: 8,
					colIdx: 1,
				},
				spellIds: [31589],
				maxPoints: 1,
			},
		],
	},
	{
		name: 'Fire',
		backgroundUrl: 'https://wow.zamimg.com/images/wow/talents/backgrounds/tbc/41.jpg',
		talents: [
			{
				fieldName: 'improvedFireball',
				location: {
					rowIdx: 0,
					colIdx: 1,
				},
				spellIds: [11069, 12338],
				maxPoints: 5,
			},
			{
				//fieldName: 'impact',
				location: {
					rowIdx: 0,
					colIdx: 2,
				},
				spellIds: [11103, 12357],
				maxPoints: 5,
			},
			{
				fieldName: 'ignite',
				location: {
					rowIdx: 1,
					colIdx: 0,
				},
				spellIds: [11119, 11120, 12846],
				maxPoints: 5,
			},
			{
				//fieldName: 'flameThrowing',
				location: {
					rowIdx: 1,
					colIdx: 1,
				},
				spellIds: [11100, 12353],
				maxPoints: 2,
			},
			{
				fieldName: 'improvedFireBlast',
				location: {
					rowIdx: 1,
					colIdx: 2,
				},
				spellIds: [11078, 11080, 12342],
				maxPoints: 3,
			},
			{
				fieldName: 'incineration',
				location: {
					rowIdx: 2,
					colIdx: 0,
				},
				spellIds: [18459],
				maxPoints: 2,
			},
			{
				fieldName: 'improvedFlamestrike',
				location: {
					rowIdx: 2,
					colIdx: 1,
				},
				spellIds: [11108, 12349],
				maxPoints: 3,
			},
			{
				fieldName: 'pyroblast',
				location: {
					rowIdx: 2,
					colIdx: 2,
				},
				spellIds: [11366],
				maxPoints: 1,
			},
			{
				fieldName: 'burningSoul',
				location: {
					rowIdx: 2,
					colIdx: 3,
				},
				spellIds: [11083, 12351],
				maxPoints: 2,
			},
			{
				fieldName: 'improvedScorch',
				location: {
					rowIdx: 3,
					colIdx: 0,
				},
				spellIds: [11095, 12872],
				maxPoints: 3,
			},
			{
				//fieldName: 'moltenShields',
				location: {
					rowIdx: 3,
					colIdx: 1,
				},
				spellIds: [11094, 13043],
				maxPoints: 2,
			},
			{
				fieldName: 'masterOfElements',
				location: {
					rowIdx: 3,
					colIdx: 3,
				},
				spellIds: [29074],
				maxPoints: 3,
			},
			{
				fieldName: 'playingWithFire',
				location: {
					rowIdx: 4,
					colIdx: 0,
				},
				spellIds: [31638],
				maxPoints: 3,
			},
			{
				fieldName: 'criticalMass',
				location: {
					rowIdx: 4,
					colIdx: 1,
				},
				spellIds: [11115, 11367],
				maxPoints: 3,
			},
			{
				fieldName: 'blastWave',
				location: {
					rowIdx: 4,
					colIdx: 2,
				},
				prereqLocation: {
					rowIdx: 2,
					colIdx: 2,
				},
				spellIds: [11113],
				maxPoints: 1,
			},
			{
				//fieldName: 'blazingSpeed',
				location: {
					rowIdx: 5,
					colIdx: 0,
				},
				spellIds: [31641],
				maxPoints: 2,
			},
			{
				fieldName: 'firePower',
				location: {
					rowIdx: 5,
					colIdx: 2,
				},
				spellIds: [11124, 12378, 12398],
				maxPoints: 5,
			},
			{
				fieldName: 'pyromaniac',
				location: {
					rowIdx: 6,
					colIdx: 0,
				},
				spellIds: [34293, 34295],
				maxPoints: 3,
			},
			{
				fieldName: 'combustion',
				location: {
					rowIdx: 6,
					colIdx: 1,
				},
				prereqLocation: {
					rowIdx: 4,
					colIdx: 1,
				},
				spellIds: [11129],
				maxPoints: 1,
			},
			{
				fieldName: 'moltenFury',
				location: {
					rowIdx: 6,
					colIdx: 2,
				},
				spellIds: [31679],
				maxPoints: 2,
			},
			{
				fieldName: 'empoweredFireball',
				location: {
					rowIdx: 7,
					colIdx: 2,
				},
				spellIds: [31656],
				maxPoints: 5,
			},
			{
				fieldName: 'dragonsBreath',
				location: {
					rowIdx: 8,
					colIdx: 1,
				},
				prereqLocation: {
					rowIdx: 6,
					colIdx: 1,
				},
				spellIds: [31661],
				maxPoints: 1,
			},
		],
	},
	{
		name: 'Frost',
		backgroundUrl: 'https://wow.zamimg.com/images/wow/talents/backgrounds/tbc/61.jpg',
		talents: [
			{
				//fieldName: 'frostWarding',
				location: {
					rowIdx: 0,
					colIdx: 0,
				},
				spellIds: [11189, 28332],
				maxPoints: 2,
			},
			{
				fieldName: 'improvedFrostbolt',
				location: {
					rowIdx: 0,
					colIdx: 1,
				},
				spellIds: [11070, 12473, 16763, 16765],
				maxPoints: 5,
			},
			{
				fieldName: 'elementalPrecision',
				location: {
					rowIdx: 0,
					colIdx: 2,
				},
				spellIds: [29438],
				maxPoints: 3,
			},
			{
				fieldName: 'iceShards',
				location: {
					rowIdx: 1,
					colIdx: 0,
				},
				spellIds: [11207, 12672, 15047, 15052],
				maxPoints: 5,
			},
			{
				//fieldName: 'frostbite',
				location: {
					rowIdx: 1,
					colIdx: 1,
				},
				spellIds: [11071, 12496],
				maxPoints: 3,
			},
			{
				//fieldName: 'improvedFrostNova',
				location: {
					rowIdx: 1,
					colIdx: 2,
				},
				spellIds: [11165, 12475],
				maxPoints: 2,
			},
			{
				//fieldName: 'permafrost',
				location: {
					rowIdx: 1,
					colIdx: 3,
				},
				spellIds: [11175, 12569, 12571],
				maxPoints: 3,
			},
			{
				fieldName: 'piercingIce',
				location: {
					rowIdx: 2,
					colIdx: 0,
				},
				spellIds: [11151, 12952],
				maxPoints: 3,
			},
			{
				fieldName: 'icyVeins',
				location: {
					rowIdx: 2,
					colIdx: 1,
				},
				spellIds: [12472],
				maxPoints: 1,
			},
			{
				//fieldName: 'improvedBlizzard',
				location: {
					rowIdx: 2,
					colIdx: 3,
				},
				spellIds: [11185, 12487],
				maxPoints: 3,
			},
			{
				//fieldName: 'arcticReach',
				location: {
					rowIdx: 3,
					colIdx: 0,
				},
				spellIds: [16757],
				maxPoints: 2,
			},
			{
				fieldName: 'frostChanneling',
				location: {
					rowIdx: 3,
					colIdx: 1,
				},
				spellIds: [11160, 12518],
				maxPoints: 3,
			},
			{
				fieldName: 'shatter',
				location: {
					rowIdx: 3,
					colIdx: 2,
				},
				prereqLocation: {
					rowIdx: 1,
					colIdx: 2,
				},
				spellIds: [11170, 12982],
				maxPoints: 5,
			},
			{
				//fieldName: 'frozenCore',
				location: {
					rowIdx: 4,
					colIdx: 0,
				},
				spellIds: [31667],
				maxPoints: 3,
			},
			{
				fieldName: 'coldSnap',
				location: {
					rowIdx: 4,
					colIdx: 1,
				},
				spellIds: [11958],
				maxPoints: 1,
			},
			{
				fieldName: 'improvedConeOfCold',
				location: {
					rowIdx: 4,
					colIdx: 2,
				},
				spellIds: [11190, 12489],
				maxPoints: 3,
			},
			{
				fieldName: 'iceFloes',
				location: {
					rowIdx: 5,
					colIdx: 0,
				},
				spellIds: [31670, 31672],
				maxPoints: 2,
			},
			{
				fieldName: 'wintersChill',
				location: {
					rowIdx: 5,
					colIdx: 2,
				},
				spellIds: [11180, 28592],
				maxPoints: 5,
			},
			{
				//fieldName: 'iceBarrier',
				location: {
					rowIdx: 6,
					colIdx: 1,
				},
				prereqLocation: {
					rowIdx: 4,
					colIdx: 1,
				},
				spellIds: [11426],
				maxPoints: 1,
			},
			{
				fieldName: 'arcticWinds',
				location: {
					rowIdx: 6,
					colIdx: 2,
				},
				spellIds: [31674],
				maxPoints: 5,
			},
			{
				fieldName: 'empoweredFrostbolt',
				location: {
					rowIdx: 7,
					colIdx: 1,
				},
				spellIds: [31682],
				maxPoints: 5,
			},
			{
				fieldName: 'summonWaterElemental',
				location: {
					rowIdx: 8,
					colIdx: 1,
				},
				spellIds: [31687],
				maxPoints: 1,
			},
		],
	},
]);
