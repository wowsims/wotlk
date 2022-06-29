import { Spec } from '/tbc/core/proto/common.js';
import { DruidTalents as DruidTalents } from '/tbc/core/proto/druid.js';
import { DruidSpecs } from '/tbc/core/proto_utils/utils.js';
import { Player } from '/tbc/core/player.js';

import { TalentsConfig, TalentsPicker, newTalentsConfig } from './talents_picker.js';

// Talents are the same for all Druid specs, so its ok to just use BalanceDruid here
export class DruidTalentsPicker extends TalentsPicker<Spec.SpecBalanceDruid> {
	constructor(parent: HTMLElement, player: Player<Spec.SpecBalanceDruid>) {
		super(parent, player, druidTalentsConfig)
	}
}

export const druidTalentsConfig: TalentsConfig<Spec.SpecBalanceDruid> = newTalentsConfig([
	{
		name: 'Balance',
		backgroundUrl: 'https://wow.zamimg.com/images/wow/talents/backgrounds/tbc/283.jpg',
		talents: [
			{
				fieldName: 'starlightWrath',
				location: {
					rowIdx: 0,
					colIdx: 0,
				},
				spellIds: [16814],
				maxPoints: 5,
			},
			{
				//fieldName: 'naturesGrasp',
				location: {
					rowIdx: 0,
					colIdx: 1,
				},
				spellIds: [16689],
				maxPoints: 1,
			},
			{
				//fieldName: 'improvedNaturesGrasp',
				location: {
					rowIdx: 0,
					colIdx: 2,
				},
				prereqLocation: {
					rowIdx: 0,
					colIdx: 1,
				},
				spellIds: [17245, 17247],
				maxPoints: 4,
			},
			{
				//fieldName: 'controlOfNature',
				location: {
					rowIdx: 1,
					colIdx: 0,
				},
				spellIds: [16918],
				maxPoints: 3,
			},
			{
				fieldName: 'focusedStarlight',
				location: {
					rowIdx: 1,
					colIdx: 1,
				},
				spellIds: [35363],
				maxPoints: 2,
			},
			{
				fieldName: 'improvedMoonfire',
				location: {
					rowIdx: 1,
					colIdx: 2,
				},
				spellIds: [16821],
				maxPoints: 2,
			},
			{
				fieldName: 'brambles',
				location: {
					rowIdx: 2,
					colIdx: 0,
				},
				spellIds: [16836, 16839],
				maxPoints: 3,
			},
			{
				fieldName: 'insectSwarm',
				location: {
					rowIdx: 2,
					colIdx: 2,
				},
				spellIds: [5570],
				maxPoints: 1,
			},
			{
				//fieldName: 'naturesReach',
				location: {
					rowIdx: 2,
					colIdx: 3,
				},
				spellIds: [16819],
				maxPoints: 2,
			},
			{
				fieldName: 'vengeance',
				location: {
					rowIdx: 3,
					colIdx: 1,
				},
				prereqLocation: {
					rowIdx: 1,
					colIdx: 1,
				},
				spellIds: [16909],
				maxPoints: 5,
			},
			{
				//fieldName: 'celestialFocus',
				location: {
					rowIdx: 3,
					colIdx: 2,
				},
				spellIds: [16850, 16923],
				maxPoints: 3,
			},
			{
				fieldName: 'lunarGuidance',
				location: {
					rowIdx: 4,
					colIdx: 0,
				},
				spellIds: [33589],
				maxPoints: 3,
			},
			{
				fieldName: 'naturesGrace',
				location: {
					rowIdx: 4,
					colIdx: 1,
				},
				spellIds: [16880],
				maxPoints: 1,
			},
			{
				fieldName: 'moonglow',
				location: {
					rowIdx: 4,
					colIdx: 2,
				},
				spellIds: [16845],
				maxPoints: 3,
			},
			{
				fieldName: 'moonfury',
				location: {
					rowIdx: 5,
					colIdx: 1,
				},
				prereqLocation: {
					rowIdx: 4,
					colIdx: 1,
				},
				spellIds: [16896, 16897, 16899],
				maxPoints: 5,
			},
			{
				fieldName: 'balanceOfPower',
				location: {
					rowIdx: 5,
					colIdx: 2,
				},
				spellIds: [33592, 33596],
				maxPoints: 2,
			},
			{
				fieldName: 'dreamstate',
				location: {
					rowIdx: 6,
					colIdx: 0,
				},
				spellIds: [33597, 33599, 33956],
				maxPoints: 3,
			},
			{
				fieldName: 'moonkinForm',
				location: {
					rowIdx: 6,
					colIdx: 1,
				},
				spellIds: [24858],
				maxPoints: 1,
			},
			{
				fieldName: 'improvedFaerieFire',
				location: {
					rowIdx: 6,
					colIdx: 2,
				},
				spellIds: [33600],
				maxPoints: 3,
			},
			{
				fieldName: 'wrathOfCenarius',
				location: {
					rowIdx: 7,
					colIdx: 1,
				},
				spellIds: [33603],
				maxPoints: 5,
			},
			{
				fieldName: 'forceOfNature',
				location: {
					rowIdx: 8,
					colIdx: 1,
				},
				spellIds: [33831],
				maxPoints: 1,
			},
		],
	},
	{
		name: 'Feral Combat',
		backgroundUrl: 'https://wow.zamimg.com/images/wow/talents/backgrounds/tbc/281.jpg',
		talents: [
			{
				fieldName: 'ferocity',
				location: {
					rowIdx: 0,
					colIdx: 1,
				},
				spellIds: [16934],
				maxPoints: 5,
			},
			{
				fieldName: 'feralAggression',
				location: {
					rowIdx: 0,
					colIdx: 2,
				},
				spellIds: [16858],
				maxPoints: 5,
			},
			{
				fieldName: 'feralInstinct',
				location: {
					rowIdx: 1,
					colIdx: 0,
				},
				spellIds: [16947],
				maxPoints: 3,
			},
			{
				//fieldName: 'brutalImpact',
				location: {
					rowIdx: 1,
					colIdx: 1,
				},
				spellIds: [16940],
				maxPoints: 2,
			},
			{
				fieldName: 'thickHide',
				location: {
					rowIdx: 1,
					colIdx: 2,
				},
				spellIds: [16929],
				maxPoints: 3,
			},
			{
				fieldName: 'feralSwiftness',
				location: {
					rowIdx: 2,
					colIdx: 0,
				},
				spellIds: [17002, 24866],
				maxPoints: 2,
			},
			{
				//fieldName: 'feralCharge',
				location: {
					rowIdx: 2,
					colIdx: 1,
				},
				spellIds: [16979],
				maxPoints: 1,
			},
			{
				fieldName: 'sharpenedClaws',
				location: {
					rowIdx: 2,
					colIdx: 2,
				},
				spellIds: [16942],
				maxPoints: 3,
			},
			{
				fieldName: 'shreddingAttacks',
				location: {
					rowIdx: 3,
					colIdx: 0,
				},
				spellIds: [16966, 16968],
				maxPoints: 2,
			},
			{
				fieldName: 'predatoryStrikes',
				location: {
					rowIdx: 3,
					colIdx: 1,
				},
				spellIds: [16972, 16974],
				maxPoints: 3,
			},
			{
				fieldName: 'primalFury',
				location: {
					rowIdx: 3,
					colIdx: 2,
				},
				spellIds: [37116],
				maxPoints: 2,
			},
			{
				fieldName: 'savageFury',
				location: {
					rowIdx: 4,
					colIdx: 0,
				},
				spellIds: [16998],
				maxPoints: 2,
			},
			{
				fieldName: 'faerieFire',
				location: {
					rowIdx: 4,
					colIdx: 2,
				},
				spellIds: [16857],
				maxPoints: 1,
			},
			{
				//fieldName: 'nurturingInstinct',
				location: {
					rowIdx: 4,
					colIdx: 3,
				},
				spellIds: [33872],
				maxPoints: 2,
			},
			{
				fieldName: 'heartOfTheWild',
				location: {
					rowIdx: 5,
					colIdx: 1,
				},
				prereqLocation: {
					rowIdx: 3,
					colIdx: 1,
				},
				spellIds: [17003, 17004, 17005, 17006, 24894],
				maxPoints: 5,
			},
			{
				fieldName: 'survivalOfTheFittest',
				location: {
					rowIdx: 5,
					colIdx: 2,
				},
				spellIds: [33853, 33855],
				maxPoints: 3,
			},
			{
				//fieldName: 'primalTenacity',
				location: {
					rowIdx: 6,
					colIdx: 0,
				},
				spellIds: [33851, 33852, 33957],
				maxPoints: 3,
			},
			{
				fieldName: 'leaderOfThePack',
				location: {
					rowIdx: 6,
					colIdx: 1,
				},
				spellIds: [17007],
				maxPoints: 1,
			},
			{
				fieldName: 'improvedLeaderOfThePack',
				location: {
					rowIdx: 6,
					colIdx: 2,
				},
				prereqLocation: {
					rowIdx: 6,
					colIdx: 1,
				},
				spellIds: [34297, 34300],
				maxPoints: 2,
			},
			{
				fieldName: 'predatoryInstincts',
				location: {
					rowIdx: 7,
					colIdx: 2,
				},
				spellIds: [33859, 33866],
				maxPoints: 5,
			},
			{
				fieldName: 'mangle',
				location: {
					rowIdx: 8,
					colIdx: 1,
				},
				prereqLocation: {
					rowIdx: 6,
					colIdx: 1,
				},
				spellIds: [33917],
				maxPoints: 1,
			},
		],
	},
	{
		name: 'Restoration',
		backgroundUrl: 'https://wow.zamimg.com/images/wow/talents/backgrounds/tbc/282.jpg',
		talents: [
			{
				fieldName: 'improvedMarkOfTheWild',
				location: {
					rowIdx: 0,
					colIdx: 1,
				},
				spellIds: [17050, 17051, 17053],
				maxPoints: 5,
			},
			{
				fieldName: 'furor',
				location: {
					rowIdx: 0,
					colIdx: 2,
				},
				spellIds: [17056, 17058],
				maxPoints: 5,
			},
			{
				fieldName: 'naturalist',
				location: {
					rowIdx: 1,
					colIdx: 0,
				},
				spellIds: [17069],
				maxPoints: 5,
			},
			{
				//fieldName: 'naturesFocus',
				location: {
					rowIdx: 1,
					colIdx: 1,
				},
				spellIds: [17063, 17065],
				maxPoints: 5,
			},
			{
				fieldName: 'naturalShapeshifter',
				location: {
					rowIdx: 1,
					colIdx: 2,
				},
				spellIds: [16833],
				maxPoints: 3,
			},
			{
				fieldName: 'intensity',
				location: {
					rowIdx: 2,
					colIdx: 0,
				},
				spellIds: [17106],
				maxPoints: 3,
			},
			{
				fieldName: 'subtlety',
				location: {
					rowIdx: 2,
					colIdx: 1,
				},
				spellIds: [17118],
				maxPoints: 5,
			},
			{
				fieldName: 'omenOfClarity',
				location: {
					rowIdx: 2,
					colIdx: 2,
				},
				spellIds: [16864],
				maxPoints: 1,
			},
			{
				//fieldName: 'tranquilSpirit',
				location: {
					rowIdx: 3,
					colIdx: 1,
				},
				spellIds: [24968],
				maxPoints: 5,
			},
			{
				//fieldName: 'improvedRejuvenation',
				location: {
					rowIdx: 3,
					colIdx: 2,
				},
				spellIds: [17111],
				maxPoints: 3,
			},
			{
				fieldName: 'naturesSwiftness',
				location: {
					rowIdx: 4,
					colIdx: 0,
				},
				prereqLocation: {
					rowIdx: 2,
					colIdx: 0,
				},
				spellIds: [17116],
				maxPoints: 1,
			},
			{
				//fieldName: 'giftOfNature',
				location: {
					rowIdx: 4,
					colIdx: 1,
				},
				spellIds: [17104, 24943],
				maxPoints: 5,
			},
			{
				//fieldName: 'improvedTranquility',
				location: {
					rowIdx: 4,
					colIdx: 3,
				},
				spellIds: [17123],
				maxPoints: 2,
			},
			{
				//fieldName: 'empoweredTouch',
				location: {
					rowIdx: 5,
					colIdx: 0,
				},
				spellIds: [33879],
				maxPoints: 2,
			},
			{
				//fieldName: 'improvedRegrowth',
				location: {
					rowIdx: 5,
					colIdx: 2,
				},
				prereqLocation: {
					rowIdx: 3,
					colIdx: 2,
				},
				spellIds: [17074],
				maxPoints: 5,
			},
			{
				fieldName: 'livingSpirit',
				location: {
					rowIdx: 6,
					colIdx: 0,
				},
				spellIds: [34151],
				maxPoints: 3,
			},
			{
				//fieldName: 'swiftmend',
				location: {
					rowIdx: 6,
					colIdx: 1,
				},
				prereqLocation: {
					rowIdx: 4,
					colIdx: 1,
				},
				spellIds: [18562],
				maxPoints: 1,
			},
			{
				fieldName: 'naturalPerfection',
				location: {
					rowIdx: 6,
					colIdx: 2,
				},
				spellIds: [33881],
				maxPoints: 3,
			},
			{
				//fieldName: 'empoweredRejuvenation',
				location: {
					rowIdx: 7,
					colIdx: 1,
				},
				spellIds: [33886],
				maxPoints: 5,
			},
			{
				//fieldName: 'treeOfLife',
				location: {
					rowIdx: 8,
					colIdx: 1,
				},
				spellIds: [33891],
				maxPoints: 1,
			},
		],
	},
]);
