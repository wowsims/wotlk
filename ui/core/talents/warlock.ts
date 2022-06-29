import { Spec } from '/tbc/core/proto/common.js';
import { WarlockTalents as WarlockTalents } from '/tbc/core/proto/warlock.js';
import { WarlockSpecs } from '/tbc/core/proto_utils/utils.js';
import { Player } from '/tbc/core/player.js';

import { TalentsConfig, TalentsPicker, newTalentsConfig } from './talents_picker.js';

export class WarlockTalentsPicker extends TalentsPicker<Spec.SpecWarlock> {
	constructor(parent: HTMLElement, player: Player<Spec.SpecWarlock>) {
		super(parent, player, warlockTalentsConfig);
	}
}

export const warlockTalentsConfig: TalentsConfig<Spec.SpecWarlock> = newTalentsConfig([
	{
		name: 'Affliction',
		backgroundUrl: 'https://wow.zamimg.com/images/wow/talents/backgrounds/tbc/302.jpg',
		talents: [
			{
				fieldName: 'suppression',
				location: {
					rowIdx: 0,
					colIdx: 1,
				},
				spellIds: [18174],
				maxPoints: 5,
			},
			{
				fieldName: 'improvedCorruption',
				location: {
					rowIdx: 0,
					colIdx: 2,
				},
				spellIds: [17810],
				maxPoints: 5,
			},
			{
				//fieldName: 'improvedCurseOfWeakness',
				location: {
					rowIdx: 1,
					colIdx: 0,
				},
				spellIds: [18179],
				maxPoints: 2,
			},
			{
				//fieldName: 'improvedDrainSoul',
				location: {
					rowIdx: 1,
					colIdx: 1,
				},
				spellIds: [18213, 18372],
				maxPoints: 2,
			},
			{
				fieldName: 'improvedLifeTap',
				location: {
					rowIdx: 1,
					colIdx: 2,
				},
				spellIds: [18182],
				maxPoints: 2,
			},
			{
				fieldName: 'soulSiphon',
				location: {
					rowIdx: 1,
					colIdx: 3,
				},
				spellIds: [17804],
				maxPoints: 2,
			},
			{
				fieldName: 'improvedCurseOfAgony',
				location: {
					rowIdx: 2,
					colIdx: 0,
				},
				spellIds: [18827, 18829],
				maxPoints: 2,
			},
			{
				//fieldName: 'felConcentration',
				location: {
					rowIdx: 2,
					colIdx: 1,
				},
				spellIds: [17783],
				maxPoints: 5,
			},
			{
				fieldName: 'amplifyCurse',
				location: {
					rowIdx: 2,
					colIdx: 2,
				},
				spellIds: [18288],
				maxPoints: 1,
			},
			{
				//fieldName: 'grimReach',
				location: {
					rowIdx: 3,
					colIdx: 0,
				},
				spellIds: [18218],
				maxPoints: 2,
			},
			{
				fieldName: 'nightfall',
				location: {
					rowIdx: 3,
					colIdx: 1,
				},
				spellIds: [18094],
				maxPoints: 2,
			},
			{
				fieldName: 'empoweredCorruption',
				location: {
					rowIdx: 3,
					colIdx: 3,
				},
				spellIds: [32381],
				maxPoints: 3,
			},
			{
				fieldName: 'shadowEmbrace',
				location: {
					rowIdx: 4,
					colIdx: 0,
				},
				spellIds: [32385, 32387, 32392],
				maxPoints: 5,
			},
			{
				fieldName: 'siphonLife',
				location: {
					rowIdx: 4,
					colIdx: 1,
				},
				spellIds: [18265],
				maxPoints: 1,
			},
			{
				//fieldName: 'curseOfExhaustion',
				location: {
					rowIdx: 4,
					colIdx: 2,
				},
				prereqLocation: {
					rowIdx: 2,
					colIdx: 2,
				},
				spellIds: [18223],
				maxPoints: 1,
			},
			{
				fieldName: 'shadowMastery',
				location: {
					rowIdx: 5,
					colIdx: 1,
				},
				prereqLocation: {
					rowIdx: 4,
					colIdx: 1,
				},
				spellIds: [18271],
				maxPoints: 5,
			},
			{
				fieldName: 'contagion',
				location: {
					rowIdx: 6,
					colIdx: 1,
				},
				spellIds: [30060],
				maxPoints: 5,
			},
			{
				fieldName: 'darkPact',
				location: {
					rowIdx: 6,
					colIdx: 2,
				},
				spellIds: [18220],
				maxPoints: 1,
			},
			{
				//fieldName: 'improvedHowlOfTerror',
				location: {
					rowIdx: 7,
					colIdx: 0,
				},
				spellIds: [30054, 30057],
				maxPoints: 2,
			},
			{
				fieldName: 'malediction',
				location: {
					rowIdx: 7,
					colIdx: 2,
				},
				spellIds: [32477, 32483],
				maxPoints: 3,
			},
			{
				fieldName: 'unstableAffliction',
				location: {
					rowIdx: 8,
					colIdx: 1,
				},
				prereqLocation: {
					rowIdx: 6,
					colIdx: 1,
				},
				spellIds: [30108],
				maxPoints: 1,
			},
		],
	},
	{
		name: 'Demonology',
		backgroundUrl: 'https://wow.zamimg.com/images/wow/talents/backgrounds/tbc/303.jpg',
		talents: [
			{
				//fieldName: 'improvedHealthstone',
				location: {
					rowIdx: 0,
					colIdx: 0,
				},
				spellIds: [18692],
				maxPoints: 2,
			},
			{
				fieldName: 'improvedImp',
				location: {
					rowIdx: 0,
					colIdx: 1,
				},
				spellIds: [18694],
				maxPoints: 3,
			},
			{
				fieldName: 'demonicEmbrace',
				location: {
					rowIdx: 0,
					colIdx: 2,
				},
				spellIds: [18697],
				maxPoints: 5,
			},
			{
				//fieldName: 'improvedHealthFunnel',
				location: {
					rowIdx: 1,
					colIdx: 0,
				},
				spellIds: [18703],
				maxPoints: 2,
			},
			{
				fieldName: 'improvedVoidwalker',
				location: {
					rowIdx: 1,
					colIdx: 1,
				},
				spellIds: [18705],
				maxPoints: 3,
			},
			{
				fieldName: 'felIntellect',
				location: {
					rowIdx: 1,
					colIdx: 2,
				},
				spellIds: [18731, 18743],
				maxPoints: 3,
			},
			{
				fieldName: 'improvedSayaad',
				location: {
					rowIdx: 2,
					colIdx: 0,
				},
				spellIds: [18754],
				maxPoints: 3,
			},
			{
				//fieldName: 'felDomination',
				location: {
					rowIdx: 2,
					colIdx: 1,
				},
				spellIds: [18708],
				maxPoints: 1,
			},
			{
				fieldName: 'felStamina',
				location: {
					rowIdx: 2,
					colIdx: 2,
				},
				spellIds: [18748],
				maxPoints: 3,
			},
			{
				fieldName: 'demonicAegis',
				location: {
					rowIdx: 2,
					colIdx: 3,
				},
				spellIds: [30143],
				maxPoints: 3,
			},
			{
				//fieldName: 'masterSummoner',
				location: {
					rowIdx: 3,
					colIdx: 1,
				},
				prereqLocation: {
					rowIdx: 2,
					colIdx: 1,
				},
				spellIds: [18709],
				maxPoints: 2,
			},
			{
				fieldName: 'unholyPower',
				location: {
					rowIdx: 3,
					colIdx: 2,
				},
				spellIds: [18769],
				maxPoints: 5,
			},
			{
				fieldName: 'improvedEnslaveDemon',
				location: {
					rowIdx: 4,
					colIdx: 0,
				},
				spellIds: [18821],
				maxPoints: 2,
			},
			{
				fieldName: 'demonicSacrifice',
				location: {
					rowIdx: 4,
					colIdx: 1,
				},
				spellIds: [18788],
				maxPoints: 1,
			},
			{
				fieldName: 'masterConjuror',
				location: {
					rowIdx: 4,
					colIdx: 3,
				},
				spellIds: [18767],
				maxPoints: 2,
			},
			{
				fieldName: 'manaFeed',
				location: {
					rowIdx: 5,
					colIdx: 0,
				},
				spellIds: [30326],
				maxPoints: 3,
			},
			{
				fieldName: 'masterDemonologist',
				location: {
					rowIdx: 5,
					colIdx: 2,
				},
				prereqLocation: {
					rowIdx: 3,
					colIdx: 2,
				},
				spellIds: [23785, 23822],
				maxPoints: 5,
			},
			{
				//fieldName: 'demonicResilience',
				location: {
					rowIdx: 6,
					colIdx: 0,
				},
				spellIds: [30319],
				maxPoints: 3,
			},
			{
				fieldName: 'soulLink',
				location: {
					rowIdx: 6,
					colIdx: 1,
				},
				prereqLocation: {
					rowIdx: 4,
					colIdx: 1,
				},
				spellIds: [19028],
				maxPoints: 1,
			},
			{
				fieldName: 'demonicKnowledge',
				location: {
					rowIdx: 6,
					colIdx: 2,
				},
				spellIds: [35691],
				maxPoints: 3,
			},
			{
				fieldName: 'demonicTactics',
				location: {
					rowIdx: 7,
					colIdx: 1,
				},
				spellIds: [30242, 30245],
				maxPoints: 5,
			},
			{
				fieldName: 'summonFelguard',
				location: {
					rowIdx: 8,
					colIdx: 1,
				},
				spellIds: [30146],
				maxPoints: 1,
			},
		],
	},
	{
		name: 'Destruction',
		backgroundUrl: 'https://wow.zamimg.com/images/wow/talents/backgrounds/tbc/301.jpg',
		talents: [
			{
				fieldName: 'improvedShadowBolt',
				location: {
					rowIdx: 0,
					colIdx: 1,
				},
				spellIds: [17793, 17796, 17801],
				maxPoints: 5,
			},
			{
				fieldName: 'cataclysm',
				location: {
					rowIdx: 0,
					colIdx: 2,
				},
				spellIds: [17778],
				maxPoints: 5,
			},
			{
				fieldName: 'bane',
				location: {
					rowIdx: 1,
					colIdx: 1,
				},
				spellIds: [17788],
				maxPoints: 5,
			},
			{
				//fieldName: 'aftermath',
				location: {
					rowIdx: 1,
					colIdx: 2,
				},
				spellIds: [18119],
				maxPoints: 5,
			},
			{
				fieldName: 'improvedFirebolt',
				location: {
					rowIdx: 2,
					colIdx: 0,
				},
				spellIds: [18126],
				maxPoints: 2,
			},
			{
				fieldName: 'improvedLashOfPain',
				location: {
					rowIdx: 2,
					colIdx: 1,
				},
				spellIds: [18128],
				maxPoints: 2,
			},
			{
				fieldName: 'devastation',
				location: {
					rowIdx: 2,
					colIdx: 2,
				},
				spellIds: [18130],
				maxPoints: 5,
			},
			{
				fieldName: 'shadowburn',
				location: {
					rowIdx: 2,
					colIdx: 3,
				},
				spellIds: [17877],
				maxPoints: 1,
			},
			{
				//fieldName: 'intensity',
				location: {
					rowIdx: 3,
					colIdx: 0,
				},
				spellIds: [18135],
				maxPoints: 2,
			},
			{
				fieldName: 'destructiveReach',
				location: {
					rowIdx: 3,
					colIdx: 1,
				},
				spellIds: [17917],
				maxPoints: 2,
			},
			{
				fieldName: 'improvedSearingPain',
				location: {
					rowIdx: 3,
					colIdx: 3,
				},
				spellIds: [17927, 17929],
				maxPoints: 3,
			},
			{
				//fieldName: 'pyroclasm',
				location: {
					rowIdx: 4,
					colIdx: 0,
				},
				prereqLocation: {
					rowIdx: 3,
					colIdx: 0,
				},
				spellIds: [18096, 18073],
				maxPoints: 2,
			},
			{
				fieldName: 'improvedImmolate',
				location: {
					rowIdx: 4,
					colIdx: 1,
				},
				spellIds: [17815, 17833],
				maxPoints: 5,
			},
			{
				fieldName: 'ruin',
				location: {
					rowIdx: 4,
					colIdx: 2,
				},
				prereqLocation: {
					rowIdx: 2,
					colIdx: 2,
				},
				spellIds: [17959],
				maxPoints: 1,
			},
			{
				//fieldName: 'netherProtection',
				location: {
					rowIdx: 5,
					colIdx: 0,
				},
				spellIds: [30299, 30301],
				maxPoints: 3,
			},
			{
				fieldName: 'emberstorm',
				location: {
					rowIdx: 5,
					colIdx: 2,
				},
				spellIds: [17954],
				maxPoints: 5,
			},
			{
				fieldName: 'backlash',
				location: {
					rowIdx: 6,
					colIdx: 0,
				},
				spellIds: [34935, 34938],
				maxPoints: 3,
			},
			{
				fieldName: 'conflagrate',
				location: {
					rowIdx: 6,
					colIdx: 1,
				},
				prereqLocation: {
					rowIdx: 4,
					colIdx: 1,
				},
				spellIds: [17962],
				maxPoints: 1,
			},
			{
				fieldName: 'soulLeech',
				location: {
					rowIdx: 6,
					colIdx: 2,
				},
				spellIds: [30293, 30295],
				maxPoints: 3,
			},
			{
				fieldName: 'shadowAndFlame',
				location: {
					rowIdx: 7,
					colIdx: 1,
				},
				spellIds: [30288],
				maxPoints: 5,
			},
			{
				fieldName: 'shadowfury',
				location: {
					rowIdx: 8,
					colIdx: 1,
				},
				prereqLocation: {
					rowIdx: 7,
					colIdx: 1,
				},
				spellIds: [30283],
				maxPoints: 1,
			},
		],
	},
]);
