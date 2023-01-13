import { PaladinTalents, PaladinMajorGlyph, PaladinMinorGlyph } from '../proto/paladin.js';

import { GlyphsConfig } from './glyphs_picker.js';
import { TalentsConfig, newTalentsConfig } from './talents_picker.js';

import PaladinTalentJson from './trees/paladin.json';

export const paladinTalentsConfig: TalentsConfig<PaladinTalents> = newTalentsConfig(PaladinTalentJson);

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
