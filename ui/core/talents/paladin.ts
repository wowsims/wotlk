import { Spec } from '/tbc/core/proto/common.js';
import { PaladinTalents as PaladinTalents } from '/tbc/core/proto/paladin.js';
import { PaladinSpecs } from '/tbc/core/proto_utils/utils.js';
import { Player } from '/tbc/core/player.js';

import { TalentsConfig, TalentsPicker, newTalentsConfig } from './talents_picker.js';

// Talents are the same for all Paladin specs, so its ok to just use RetributionPaladin here
export class PaladinTalentsPicker extends TalentsPicker<Spec.SpecRetributionPaladin> {
	constructor(parent: HTMLElement, player: Player<Spec.SpecRetributionPaladin>) {
		super(parent, player, paladinTalentsConfig);
	}
}

export const paladinTalentsConfig: TalentsConfig<Spec.SpecRetributionPaladin> = newTalentsConfig([
	{
		name: 'Holy',
		backgroundUrl: 'https://wow.zamimg.com/images/wow/talents/backgrounds/classic/382.jpg',
		talents: [
			{
				fieldName: 'divineStrength',
				location: {
					rowIdx: 0,
					colIdx: 1,
				},
				spellIds: [20262],
				maxPoints: 5,
			},
			{
				fieldName: 'divineIntellect',
				location: {
					rowIdx: 0,
					colIdx: 2,
				},
				spellIds: [20257],
				maxPoints: 5,
			},
			{
				//fieldName: 'spiritualFocus',
				location: {
					rowIdx: 1,
					colIdx: 1,
				},
				spellIds: [20205],
				maxPoints: 5,
			},
			{
				fieldName: 'improvedSealOfRighteousness',
				location: {
					rowIdx: 1,
					colIdx: 2,
				},
				spellIds: [20224, 20225, 20330],
				maxPoints: 5,
			},
			{
				//fieldName: 'healingLight',
				location: {
					rowIdx: 2,
					colIdx: 0,
				},
				spellIds: [20237],
				maxPoints: 3,
			},
			{
				//fieldName: 'auraMastery',
				location: {
					rowIdx: 2,
					colIdx: 1,
				},
				spellIds: [31821],
				maxPoints: 1,
			},
			{
				//fieldName: 'improvedLayOnHands',
				location: {
					rowIdx: 2,
					colIdx: 2,
				},
				spellIds: [20234],
				maxPoints: 2,
			},
			{
				//fieldName: 'unyieldingFaith',
				location: {
					rowIdx: 2,
					colIdx: 3,
				},
				spellIds: [9453, 25836],
				maxPoints: 2,
			},
			{
				fieldName: 'illumination',
				location: {
					rowIdx: 3,
					colIdx: 1,
				},
				spellIds: [20210, 20212],
				maxPoints: 5,
			},
			{
				fieldName: 'improvedBlessingOfWisdom',
				location: {
					rowIdx: 3,
					colIdx: 2,
				},
				spellIds: [20244],
				maxPoints: 2,
			},
			{
				//fieldName: 'pureOfHeart',
				location: {
					rowIdx: 4,
					colIdx: 0,
				},
				spellIds: [31822],
				maxPoints: 3,
			},
			{
				fieldName: 'divineFavor',
				location: {
					rowIdx: 4,
					colIdx: 1,
				},
				prereqLocation: {
					rowIdx: 3,
					colIdx: 1,
				},
				spellIds: [20216],
				maxPoints: 1,
			},
			{
				//fieldName: 'sanctifiedLight',
				location: {
					rowIdx: 4,
					colIdx: 2,
				},
				spellIds: [20359],
				maxPoints: 3,
			},
			{
				fieldName: 'purifyingPower',
				location: {
					rowIdx: 5,
					colIdx: 0,
				},
				spellIds: [31825],
				maxPoints: 2,
			},
			{
				fieldName: 'holyPower',
				location: {
					rowIdx: 5,
					colIdx: 2,
				},
				spellIds: [5923, 5924, 5925, 5926, 25829],
				maxPoints: 5,
			},
			{
				//fieldName: 'lightsGrace',
				location: {
					rowIdx: 6,
					colIdx: 0,
				},
				spellIds: [31833, 31835],
				maxPoints: 3,
			},
			{
				fieldName: 'holyShock',
				location: {
					rowIdx: 6,
					colIdx: 1,
				},
				prereqLocation: {
					rowIdx: 4,
					colIdx: 1,
				},
				spellIds: [20473],
				maxPoints: 1,
			},
			{
				fieldName: 'blessedLife',
				location: {
					rowIdx: 6,
					colIdx: 2,
				},
				spellIds: [31828],
				maxPoints: 3,
			},
			{
				fieldName: 'holyGuidance',
				location: {
					rowIdx: 7,
					colIdx: 1,
				},
				spellIds: [31837],
				maxPoints: 5,
			},
			{
				fieldName: 'divineIllumination',
				location: {
					rowIdx: 8,
					colIdx: 1,
				},
				spellIds: [31842],
				maxPoints: 1,
			},
		],
	},
	{
		name: 'Protection',
		backgroundUrl: 'https://wow.zamimg.com/images/wow/talents/backgrounds/classic/383.jpg',
		talents: [
			{
				fieldName: 'improvedDevotionAura',
				location: {
					rowIdx: 0,
					colIdx: 1,
				},
				spellIds: [20138],
				maxPoints: 5,
			},
			{
				fieldName: 'redoubt',
				location: {
					rowIdx: 0,
					colIdx: 2,
				},
				spellIds: [20127, 20130, 20135],
				maxPoints: 5,
			},
			{
				fieldName: 'precision',
				location: {
					rowIdx: 1,
					colIdx: 0,
				},
				spellIds: [20189, 20192],
				maxPoints: 3,
			},
			{
				//fieldName: 'guardiansFavor',
				location: {
					rowIdx: 1,
					colIdx: 1,
				},
				spellIds: [20174],
				maxPoints: 2,
			},
			{
				fieldName: 'toughness',
				location: {
					rowIdx: 1,
					colIdx: 3,
				},
				spellIds: [20143],
				maxPoints: 5,
			},
			{
				fieldName: 'blessingOfKings',
				location: {
					rowIdx: 2,
					colIdx: 0,
				},
				spellIds: [20217],
				maxPoints: 1,
			},
			{
				fieldName: 'improvedRighteousFury',
				location: {
					rowIdx: 2,
					colIdx: 1,
				},
				spellIds: [20468],
				maxPoints: 3,
			},
			{
				fieldName: 'shieldSpecialization',
				location: {
					rowIdx: 2,
					colIdx: 2,
				},
				prereqLocation: {
					rowIdx: 0,
					colIdx: 2,
				},
				spellIds: [20148],
				maxPoints: 3,
			},
			{
				fieldName: 'anticipation',
				location: {
					rowIdx: 2,
					colIdx: 3,
				},
				spellIds: [20096],
				maxPoints: 5,
			},
			{
				//fieldName: 'stoicism',
				location: {
					rowIdx: 3,
					colIdx: 0,
				},
				spellIds: [31844],
				maxPoints: 2,
			},
			{
				//fieldName: 'improvedHammerOfJustice',
				location: {
					rowIdx: 3,
					colIdx: 1,
				},
				spellIds: [20487],
				maxPoints: 3,
			},
			{
				//fieldName: 'improvedConcentrationAura',
				location: {
					rowIdx: 3,
					colIdx: 2,
				},
				spellIds: [20254],
				maxPoints: 3,
			},
			{
				fieldName: 'spellWarding',
				location: {
					rowIdx: 4,
					colIdx: 0,
				},
				spellIds: [31846],
				maxPoints: 2,
			},
			{
				fieldName: 'blessingOfSanctuary',
				location: {
					rowIdx: 4,
					colIdx: 1,
				},
				spellIds: [20911],
				maxPoints: 1,
			},
			{
				fieldName: 'reckoning',
				location: {
					rowIdx: 4,
					colIdx: 2,
				},
				spellIds: [20177, 20179, 20181, 20180, 20182],
				maxPoints: 5,
			},
			{
				fieldName: 'sacredDuty',
				location: {
					rowIdx: 5,
					colIdx: 0,
				},
				spellIds: [31848],
				maxPoints: 2,
			},
			{
				fieldName: 'oneHandedWeaponSpecialization',
				location: {
					rowIdx: 5,
					colIdx: 2,
				},
				spellIds: [20196],
				maxPoints: 5,
			},
			{
				fieldName: 'improvedHolyShield',
				location: {
					rowIdx: 6,
					colIdx: 0,
				},
				prereqLocation: {
					rowIdx: 6,
					colIdx: 1,
				},
				spellIds: [41021, 41026],
				maxPoints: 2,
			},
			{
				fieldName: 'holyShield',
				location: {
					rowIdx: 6,
					colIdx: 1,
				},
				prereqLocation: {
					rowIdx: 4,
					colIdx: 1,
				},
				spellIds: [20925],
				maxPoints: 1,
			},
			{
				fieldName: 'ardentDefender',
				location: {
					rowIdx: 6,
					colIdx: 2,
				},
				spellIds: [31850],
				maxPoints: 5,
			},
			{
				fieldName: 'combatExpertise',
				location: {
					rowIdx: 7,
					colIdx: 2,
				},
				spellIds: [31858],
				maxPoints: 5,
			},
			{
				fieldName: 'avengersShield',
				location: {
					rowIdx: 8,
					colIdx: 1,
				},
				prereqLocation: {
					rowIdx: 6,
					colIdx: 1,
				},
				spellIds: [31935],
				maxPoints: 1,
			},
		],
	},
	{
		name: 'Retribution',
		backgroundUrl: 'https://wow.zamimg.com/images/wow/talents/backgrounds/classic/381.jpg',
		talents: [
			{
				fieldName: 'improvedBlessingOfMight',
				location: {
					rowIdx: 0,
					colIdx: 1,
				},
				spellIds: [20042, 20045],
				maxPoints: 5,
			},
			{
				fieldName: 'benediction',
				location: {
					rowIdx: 0,
					colIdx: 2,
				},
				spellIds: [20101],
				maxPoints: 5,
			},
			{
				fieldName: 'improvedJudgement',
				location: {
					rowIdx: 1,
					colIdx: 0,
				},
				spellIds: [25956],
				maxPoints: 2,
			},
			{
				fieldName: 'improvedSealOfTheCrusader',
				location: {
					rowIdx: 1,
					colIdx: 1,
				},
				spellIds: [20335],
				maxPoints: 3,
			},
			{
				fieldName: 'deflection',
				location: {
					rowIdx: 1,
					colIdx: 2,
				},
				spellIds: [20060],
				maxPoints: 5,
			},
			{
				fieldName: 'vindication',
				location: {
					rowIdx: 2,
					colIdx: 0,
				},
				spellIds: [9452, 26016, 26021],
				maxPoints: 3,
			},
			{
				fieldName: 'conviction',
				location: {
					rowIdx: 2,
					colIdx: 1,
				},
				spellIds: [20117],
				maxPoints: 5,
			},
			{
				fieldName: 'sealOfCommand',
				location: {
					rowIdx: 2,
					colIdx: 2,
				},
				spellIds: [20375],
				maxPoints: 1,
			},
			{
				fieldName: 'pursuitOfJustice',
				location: {
					rowIdx: 2,
					colIdx: 3,
				},
				spellIds: [26022, 26023, 44414],
				maxPoints: 3,
			},
			{
				fieldName: 'eyeForAnEye',
				location: {
					rowIdx: 3,
					colIdx: 0,
				},
				spellIds: [9799, 25988],
				maxPoints: 2,
			},
			{
				fieldName: 'improvedRetributionAura',
				location: {
					rowIdx: 3,
					colIdx: 2,
				},
				spellIds: [20091],
				maxPoints: 2,
			},
			{
				fieldName: 'crusade',
				location: {
					rowIdx: 3,
					colIdx: 3,
				},
				spellIds: [31866],
				maxPoints: 3,
			},
			{
				fieldName: 'twoHandedWeaponSpecialization',
				location: {
					rowIdx: 4,
					colIdx: 0,
				},
				spellIds: [20111],
				maxPoints: 3,
			},
			{
				fieldName: 'sanctityAura',
				location: {
					rowIdx: 4,
					colIdx: 2,
				},
				spellIds: [20218],
				maxPoints: 1,
			},
			{
				fieldName: 'improvedSanctityAura',
				location: {
					rowIdx: 4,
					colIdx: 3,
				},
				prereqLocation: {
					rowIdx: 4,
					colIdx: 2,
				},
				spellIds: [31869],
				maxPoints: 2,
			},
			{
				fieldName: 'vengeance',
				location: {
					rowIdx: 5,
					colIdx: 1,
				},
				spellIds: [20049, 20056],
				maxPoints: 5,
			},
			{
				fieldName: 'sanctifiedJudgement',
				location: {
					rowIdx: 5,
					colIdx: 2,
				},
				spellIds: [31876],
				maxPoints: 3,
			},
			{
				fieldName: 'sanctifiedSeals',
				location: {
					rowIdx: 6,
					colIdx: 0,
				},
				spellIds: [32043, 35396],
				maxPoints: 3,
			},
			{
				//fieldName: 'repentance',
				location: {
					rowIdx: 6,
					colIdx: 1,
				},
				spellIds: [20066],
				maxPoints: 1,
			},
			{
				fieldName: 'divinePurpose',
				location: {
					rowIdx: 6,
					colIdx: 2,
				},
				spellIds: [31871],
				maxPoints: 3,
			},
			{
				fieldName: 'fanaticism',
				location: {
					rowIdx: 7,
					colIdx: 1,
				},
				prereqLocation: {
					rowIdx: 6,
					colIdx: 1,
				},
				spellIds: [31879],
				maxPoints: 5,
			},
			{
				fieldName: 'crusaderStrike',
				location: {
					rowIdx: 8,
					colIdx: 1,
				},
				spellIds: [35395],
				maxPoints: 1,
			},
		],
	},
]);
