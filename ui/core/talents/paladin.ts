import { PaladinTalents, PaladinMajorGlyph, PaladinMinorGlyph } from '../proto/paladin.js';

import { GlyphsConfig } from './glyphs_picker.js';
import { TalentsConfig, newTalentsConfig } from './talents_picker.js';

export const paladinTalentsConfig: TalentsConfig<PaladinTalents> = newTalentsConfig([
	{
		name: 'Holy',
		backgroundUrl: 'https://wow.zamimg.com/images/wow/talents/backgrounds/classic/382.jpg',
		talents: [
			{
				fieldName: 'spiritualFocus',
				location: {
					rowIdx: 0,
					colIdx: 1,
				},
				spellIds: [20205],
				maxPoints: 5,
			},
			{
				fieldName: 'sealsOfThePure',
				location: {
					rowIdx: 0,
					colIdx: 2,
				},
				spellIds: [20224, 20225, 20330],
				maxPoints: 5,
			},
			{
				fieldName: 'healingLight',
				location: {
					rowIdx: 1,
					colIdx: 0,
				},
				spellIds: [20237],
				maxPoints: 3,
			},
			{
				fieldName: 'divineIntellect',
				location: {
					rowIdx: 1,
					colIdx: 1,
				},
				spellIds: [20257],
				maxPoints: 5,
			},
			{
				fieldName: 'unyieldingFaith',
				location: {
					rowIdx: 1,
					colIdx: 2,
				},
				spellIds: [9453, 25836],
				maxPoints: 2,
			},
			{
				fieldName: 'auraMastery',
				location: {
					rowIdx: 2,
					colIdx: 0,
				},
				spellIds: [31821],
				maxPoints: 1,
			},
			{
				fieldName: 'illumination',
				location: {
					rowIdx: 2,
					colIdx: 1,
				},
				spellIds: [20210, 20212],
				maxPoints: 5,
			},
			{
				fieldName: 'improvedLayOnHands',
				location: {
					rowIdx: 2,
					colIdx: 2,
				},
				spellIds: [20234],
				maxPoints: 2,
			},
			{
				fieldName: 'improvedConcentrationAura',
				location: {
					rowIdx: 3,
					colIdx: 0,
				},
				spellIds: [20254],
				maxPoints: 3,
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
				fieldName: 'blessedHands',
				location: {
					rowIdx: 3,
					colIdx: 3,
				},
				spellIds: [53660],
				maxPoints: 2,
			},
			{
				fieldName: 'pureOfHeart',
				location: {
					rowIdx: 4,
					colIdx: 0,
				},
				spellIds: [31822],
				maxPoints: 2,
			},
			{
				fieldName: 'divineFavor',
				location: {
					rowIdx: 4,
					colIdx: 1,
				},
				prereqLocation: {
					rowIdx: 2,
					colIdx: 1,
				},
				spellIds: [20216],
				maxPoints: 1,
			},
			{
				fieldName: 'sanctifiedLight',
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
				fieldName: 'lightsGrace',
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
				fieldName: 'sacredCleansing',
				location: {
					rowIdx: 7,
					colIdx: 0,
				},
				spellIds: [53551],
				maxPoints: 3,
			},
			{
				fieldName: 'holyGuidance',
				location: {
					rowIdx: 7,
					colIdx: 2,
				},
				spellIds: [31837],
				maxPoints: 5,
			},
			{
				fieldName: 'divineIllumination',
				location: {
					rowIdx: 8,
					colIdx: 0,
				},
				spellIds: [31842],
				maxPoints: 1,
			},
			{
				fieldName: 'judgementsOfThePure',
				location: {
					rowIdx: 8,
					colIdx: 2,
				},
				spellIds: [53671, 53673, 54151, 54154, 54155],
				maxPoints: 5,
			},
			{
				fieldName: 'infusionOfLight',
				location: {
					rowIdx: 9,
					colIdx: 1,
				},
				prereqLocation: {
					rowIdx: 6,
					colIdx: 1,
				},
				spellIds: [53569, 53576],
				maxPoints: 2,
			},
			{
				fieldName: 'enlightenedJudgements',
				location: {
					rowIdx: 9,
					colIdx: 2,
				},
				spellIds: [53556],
				maxPoints: 2,
			},
			{
				fieldName: 'beaconOfLight',
				location: {
					rowIdx: 10,
					colIdx: 1,
				},
				spellIds: [53563],
				maxPoints: 1,
			},
		],
	},
	{
		name: 'Protection',
		backgroundUrl: 'https://wow.zamimg.com/images/wow/talents/backgrounds/classic/383.jpg',
		talents: [
			{
				fieldName: 'divinity',
				location: {
					rowIdx: 0,
					colIdx: 1,
				},
				spellIds: [63646],
				maxPoints: 5,
			},
			{
				fieldName: 'divineStrength',
				location: {
					rowIdx: 0,
					colIdx: 2,
				},
				spellIds: [20262],
				maxPoints: 5,
			},
			{
				fieldName: 'stoicism',
				location: {
					rowIdx: 1,
					colIdx: 0,
				},
				spellIds: [31844, 31845, 53519],
				maxPoints: 3,
			},
			{
				fieldName: 'guardiansFavor',
				location: {
					rowIdx: 1,
					colIdx: 1,
				},
				spellIds: [20174],
				maxPoints: 2,
			},
			{
				fieldName: 'anticipation',
				location: {
					rowIdx: 1,
					colIdx: 2,
				},
				spellIds: [20096],
				maxPoints: 5,
			},
			{
				fieldName: 'divineSacrifice',
				location: {
					rowIdx: 2,
					colIdx: 0,
				},
				spellIds: [64205],
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
				fieldName: 'toughness',
				location: {
					rowIdx: 2,
					colIdx: 2,
				},
				spellIds: [20143],
				maxPoints: 5,
			},
			{
				fieldName: 'divineGuardian',
				location: {
					rowIdx: 3,
					colIdx: 0,
				},
				prereqLocation: {
					rowIdx: 2,
					colIdx: 0,
				},
				spellIds: [53527, 53530],
				maxPoints: 2,
			},
			{
				fieldName: 'improvedHammerOfJustice',
				location: {
					rowIdx: 3,
					colIdx: 1,
				},
				spellIds: [20487],
				maxPoints: 2,
			},
			{
				fieldName: 'improvedDevotionAura',
				location: {
					rowIdx: 3,
					colIdx: 2,
				},
				spellIds: [20138],
				maxPoints: 3,
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
				maxPoints: 3,
			},
			{
				fieldName: 'spiritualAttunement',
				location: {
					rowIdx: 6,
					colIdx: 0,
				},
				spellIds: [31785, 33776],
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
				maxPoints: 3,
			},
			{
				fieldName: 'redoubt',
				location: {
					rowIdx: 7,
					colIdx: 0,
				},
				spellIds: [20127, 20130, 20135],
				maxPoints: 3,
			},
			{
				fieldName: 'combatExpertise',
				location: {
					rowIdx: 7,
					colIdx: 2,
				},
				spellIds: [31858],
				maxPoints: 3,
			},
			{
				fieldName: 'touchedByTheLight',
				location: {
					rowIdx: 8,
					colIdx: 0,
				},
				spellIds: [53590],
				maxPoints: 3,
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
			{
				fieldName: 'guardedByTheLight',
				location: {
					rowIdx: 8,
					colIdx: 2,
				},
				spellIds: [53583, 53585],
				maxPoints: 2,
			},
			{
				fieldName: 'shieldOfTheTemplar',
				location: {
					rowIdx: 9,
					colIdx: 1,
				},
				prereqLocation: {
					rowIdx: 8,
					colIdx: 2,
				},
				spellIds: [53709],
				maxPoints: 3,
			},
			{
				fieldName: 'judgementsOfTheJust',
				location: {
					rowIdx: 9,
					colIdx: 2,
				},
				spellIds: [53695],
				maxPoints: 2,
			},
			{
				fieldName: 'hammerOfTheRighteous',
				location: {
					rowIdx: 10,
					colIdx: 1,
				},
				spellIds: [53595],
				maxPoints: 1,
			},
		],
	},
	{
		name: 'Retribution',
		backgroundUrl: 'https://wow.zamimg.com/images/wow/talents/backgrounds/classic/381.jpg',
		talents: [
			{
				fieldName: 'deflection',
				location: {
					rowIdx: 0,
					colIdx: 1,
				},
				spellIds: [20060],
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
				fieldName: 'improvedJudgements',
				location: {
					rowIdx: 1,
					colIdx: 0,
				},
				spellIds: [25956],
				maxPoints: 2,
			},
			{
				fieldName: 'heartOfTheCrusader',
				location: {
					rowIdx: 1,
					colIdx: 1,
				},
				spellIds: [20335],
				maxPoints: 3,
			},
			{
				fieldName: 'improvedBlessingOfMight',
				location: {
					rowIdx: 1,
					colIdx: 2,
				},
				spellIds: [20042, 20045],
				maxPoints: 2,
			},
			{
				fieldName: 'vindication',
				location: {
					rowIdx: 2,
					colIdx: 0,
				},
				spellIds: [9452, 26016],
				maxPoints: 2,
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
				spellIds: [26022],
				maxPoints: 2,
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
				fieldName: 'sanctityOfBattle',
				location: {
					rowIdx: 3,
					colIdx: 2,
				},
				spellIds: [32043, 35396, 35397],
				maxPoints: 3,
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
				fieldName: 'sanctifiedRetribution',
				location: {
					rowIdx: 4,
					colIdx: 2,
				},
				spellIds: [31869],
				maxPoints: 1,
			},
			{
				fieldName: 'vengeance',
				location: {
					rowIdx: 5,
					colIdx: 1,
				},
				prereqLocation: {
					rowIdx: 2,
					colIdx: 1,
				},
				spellIds: [20049, 20056, 20057],
				maxPoints: 3,
			},
			{
				fieldName: 'divinePurpose',
				location: {
					rowIdx: 5,
					colIdx: 2,
				},
				spellIds: [31871],
				maxPoints: 2,
			},
			{
				fieldName: 'theArtOfWar',
				location: {
					rowIdx: 6,
					colIdx: 0,
				},
				spellIds: [53486, 53488],
				maxPoints: 2,
			},
			{
				fieldName: 'repentance',
				location: {
					rowIdx: 6,
					colIdx: 1,
				},
				spellIds: [20066],
				maxPoints: 1,
			},
			{
				fieldName: 'judgementsOfTheWise',
				location: {
					rowIdx: 6,
					colIdx: 2,
				},
				spellIds: [31876],
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
				maxPoints: 3,
			},
			{
				fieldName: 'sanctifiedWrath',
				location: {
					rowIdx: 7,
					colIdx: 2,
				},
				spellIds: [53375],
				maxPoints: 2,
			},
			{
				fieldName: 'swiftRetribution',
				location: {
					rowIdx: 8,
					colIdx: 0,
				},
				spellIds: [53379, 53484, 53648],
				maxPoints: 3,
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
			{
				fieldName: 'sheathOfLight',
				location: {
					rowIdx: 8,
					colIdx: 2,
				},
				spellIds: [53501],
				maxPoints: 3,
			},
			{
				fieldName: 'righteousVengeance',
				location: {
					rowIdx: 9,
					colIdx: 1,
				},
				spellIds: [53380],
				maxPoints: 3,
			},
			{
				fieldName: 'divineStorm',
				location: {
					rowIdx: 10,
					colIdx: 1,
				},
				spellIds: [53385],
				maxPoints: 1,
			},
		],
	},
]);

export const paladinGlyphsConfig: GlyphsConfig = {
	majorGlyphs: {
		[PaladinMajorGlyph.GlyphOfAvengerSShield]: {
			name: 'Glyph of Avenger\'s Shield',
			description: 'Your Avenger\'s Shield hits 2 fewer targets, but for 100% more damage.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_holy_avengersshield.jpg',
		},
		[PaladinMajorGlyph.GlyphOfAvengingWrath]: {
			name: 'Glyph of Avenging Wrath',
			description: 'Reduces the cooldown of your Hammer of Wrath spell by 50% while Avenging Wrath is active.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_holy_avenginewrath.jpg',
		},
		[PaladinMajorGlyph.GlyphOfBeaconOfLight]: {
			name: 'Glyph of Beacon of Light',
			description: 'Increases the duration of Beacon of Light by 30 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_paladin_beaconoflight.jpg',
		},
		[PaladinMajorGlyph.GlyphOfCleansing]: {
			name: 'Glyph of Cleansing',
			description: 'Reduces the mana cost of your Cleanse and Purify spells by 20%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_holy_purify.jpg',
		},
		[PaladinMajorGlyph.GlyphOfConsecration]: {
			name: 'Glyph of Consecration',
			description: 'Increases the duration and cooldown of Consecration by 2 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_holy_innerfire.jpg',
		},
		[PaladinMajorGlyph.GlyphOfCrusaderStrike]: {
			name: 'Glyph of Crusader Strike',
			description: 'Reduces the mana cost of your Crusader Strike ability by 20%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_holy_crusaderstrike.jpg',
		},
		[PaladinMajorGlyph.GlyphOfDivinePlea]: {
			name: 'Glyph of Divine Plea',
			description: 'While Divine Plea is active, you take 3% reduced damage from all sources.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_holy_aspiration.jpg',
		},
		[PaladinMajorGlyph.GlyphOfDivineStorm]: {
			name: 'Glyph of Divine Storm',
			description: 'Your Divine Storm now heals for an additional 15% of the damage it causes.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_paladin_divinestorm.jpg',
		},
		[PaladinMajorGlyph.GlyphOfDivinity]: {
			name: 'Glyph of Divinity',
			description: 'Your Lay on Hands grants twice as much mana as normal and also grants you as much mana as it grants your target.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_holy_layonhands.jpg',
		},
		[PaladinMajorGlyph.GlyphOfExorcism]: {
			name: 'Glyph of Exorcism',
			description: 'Increases damage done by Exorcism by 20%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_holy_excorcism_02.jpg',
		},
		[PaladinMajorGlyph.GlyphOfFlashOfLight]: {
			name: 'Glyph of Flash of Light',
			description: 'Your Flash of Light has an additional 5% critical strike chance.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_holy_flashheal.jpg',
		},
		[PaladinMajorGlyph.GlyphOfHammerOfJustice]: {
			name: 'Glyph of Hammer of Justice',
			description: 'Increases your Hammer of Justice range by 5 yards.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_holy_sealofmight.jpg',
		},
		[PaladinMajorGlyph.GlyphOfHammerOfTheRighteous]: {
			name: 'Glyph of Hammer of the Righteous',
			description: 'Your Hammer of the Righteous hits 1 additional target.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_paladin_hammeroftherighteous.jpg',
		},
		[PaladinMajorGlyph.GlyphOfHammerOfWrath]: {
			name: 'Glyph of Hammer of Wrath',
			description: 'Reduces the mana cost of Hammer of Wrath by 100%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_thunderclap.jpg',
		},
		[PaladinMajorGlyph.GlyphOfHolyLight]: {
			name: 'Glyph of Holy Light',
			description: 'Your Holy Light grants 10% of its heal amount to up to 5 friendly targets within 8 yards of the initial target.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_holy_holybolt.jpg',
		},
		[PaladinMajorGlyph.GlyphOfHolyShock]: {
			name: 'Glyph of Holy Shock',
			description: 'Reduces the cooldown of Holy Shock by 1 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_holy_searinglight.jpg',
		},
		[PaladinMajorGlyph.GlyphOfHolyWrath]: {
			name: 'Glyph of Holy Wrath',
			description: 'Reduces the cooldown of your Holy Wrath spell by 15 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_holy_excorcism.jpg',
		},
		[PaladinMajorGlyph.GlyphOfJudgement]: {
			name: 'Glyph of Judgement',
			description: 'Your Judgements deal 10% more damage.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_paladin_judgementred.jpg',
		},
		[PaladinMajorGlyph.GlyphOfRighteousDefense]: {
			name: 'Glyph of Righteous Defense',
			description: 'Increases the chance for your Righteous Defense and Hand of Reckoning abilities to work successfully by 8% on each target.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/inv_shoulder_37.jpg',
		},
		[PaladinMajorGlyph.GlyphOfSalvation]: {
			name: 'Glyph of Salvation',
			description: 'When you cast Hand of Salvation on yourself, it also reduces damage taken by 20%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_holy_sealofsalvation.jpg',
		},
		[PaladinMajorGlyph.GlyphOfSealOfCommand]: {
			name: 'Glyph of Seal of Command',
			description: 'You gain 8% of your base mana each time you use a Judgement with Seal of Command active.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_warrior_innerrage.jpg',
		},
		[PaladinMajorGlyph.GlyphOfSealOfLight]: {
			name: 'Glyph of Seal of Light',
			description: 'While Seal of Light is active, the effect of your healing spells is increased by 5%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_holy_healingaura.jpg',
		},
		[PaladinMajorGlyph.GlyphOfSealOfRighteousness]: {
			name: 'Glyph of Seal of Righteousness',
			description: 'Increases the damage done by Seal of Righteousness by 10%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_thunderbolt.jpg',
		},
		[PaladinMajorGlyph.GlyphOfSealOfVengeance]: {
			name: 'Glyph of Seal of Vengeance',
			description: 'Your Seal of Vengeance or Seal of Corruption also grants 10 expertise while active.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_holy_sealofvengeance.jpg',
		},
		[PaladinMajorGlyph.GlyphOfSealOfWisdom]: {
			name: 'Glyph of Seal of Wisdom',
			description: 'While Seal of Wisdom is active, the cost of your healing spells is reduced by 5%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_holy_righteousnessaura.jpg',
		},
		[PaladinMajorGlyph.GlyphOfShieldOfRighteousness]: {
			name: 'Glyph of Shield of Righteousness',
			description: 'Reduces the mana cost of Shield of Righteousness by 80%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_paladin_shieldofvengeance.jpg',
		},
		[PaladinMajorGlyph.GlyphOfSpiritualAttunement]: {
			name: 'Glyph of Spiritual Attunement',
			description: 'Increases the amount of mana gained from your Spiritual Attunement spell by an additional 2%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_holy_revivechampion.jpg',
		},
		[PaladinMajorGlyph.GlyphOfTurnEvil]: {
			name: 'Glyph of Turn Evil',
			description: 'Reduces the casting time of your Turn Evil spell by 100%, but increases the cooldown by 8 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_holy_turnundead.jpg',
		},
	},
	minorGlyphs: {
		[PaladinMinorGlyph.GlyphOfBlessingOfKings]: {
			name: 'Glyph of Blessing of Kings',
			description: 'Reduces the mana cost of your Blessing of Kings and Greater Blessing of Kings spells by 50%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_magic_magearmor.jpg',
		},
		[PaladinMinorGlyph.GlyphOfBlessingOfMight]: {
			name: 'Glyph of Blessing of Might',
			description: 'Increases the duration of your Blessing of Might spell by 20 min when cast on yourself.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_holy_fistofjustice.jpg',
		},
		[PaladinMinorGlyph.GlyphOfBlessingOfWisdom]: {
			name: 'Glyph of Blessing of Wisdom',
			description: 'Increases the duration of your Blessing of Wisdom spell by 20 min when cast on yourself.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_holy_sealofwisdom.jpg',
		},
		[PaladinMinorGlyph.GlyphOfLayOnHands]: {
			name: 'Glyph of Lay on Hands',
			description: 'Reduces the cooldown of your Lay on Hands spell by 5 min.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_holy_layonhands.jpg',
		},
		[PaladinMinorGlyph.GlyphOfSenseUndead]: {
			name: 'Glyph of Sense Undead',
			description: 'Damage against Undead increased by 1% while your Sense Undead ability is active.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_holy_senseundead.jpg',
		},
		[PaladinMinorGlyph.GlyphOfTheWise]: {
			name: 'Glyph of the Wise',
			description: 'Reduces the mana cost of your Seal of Wisdom spell by 50%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_holy_righteousnessaura.jpg',
		},
	},
};
