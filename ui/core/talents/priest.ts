import { PriestTalents, PriestMajorGlyph, PriestMinorGlyph } from '../proto/priest.js';

import { GlyphsConfig } from './glyphs_picker.js';
import { TalentsConfig, newTalentsConfig } from './talents_picker.js';

import PriestTalentJson from './trees/priest.json';

export const priestTalentsConfig: TalentsConfig<PriestTalents> = newTalentsConfig(PriestTalentJson);

export const priestGlyphsConfig: GlyphsConfig = {
	majorGlyphs: {
		[PriestMajorGlyph.GlyphOfCircleOfHealing]: {
			name: 'Glyph of Circle of Healing',
			description: 'Your Circle of Healing spell heals 1 additional target.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_holy_circleofrenewal.jpg',
		},
		[PriestMajorGlyph.GlyphOfDispelMagic]: {
			name: 'Glyph of Dispel Magic',
			description: 'Your Dispel Magic spell also heals your target for 3% of maximum health.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_holy_dispelmagic.jpg',
		},
		[PriestMajorGlyph.GlyphOfDispersion]: {
			name: 'Glyph of Dispersion',
			description: 'Reduces the cooldown on Dispersion by 45 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_dispersion.jpg',
		},
		[PriestMajorGlyph.GlyphOfFade]: {
			name: 'Glyph of Fade',
			description: 'Reduces the cooldown of your Fade spell by 9 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_magic_lesserinvisibilty.jpg',
		},
		[PriestMajorGlyph.GlyphOfFearWard]: {
			name: 'Glyph of Fear Ward',
			description: 'Reduces cooldown and duration of Fear Ward by 60 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_holy_excorcism.jpg',
		},
		[PriestMajorGlyph.GlyphOfFlashHeal]: {
			name: 'Glyph of Flash Heal',
			description: 'Reduces the mana cost of your Flash Heal by 10%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_holy_flashheal.jpg',
		},
		[PriestMajorGlyph.GlyphOfGuardianSpirit]: {
			name: 'Glyph of Guardian Spirit',
			description: 'If your Guardian Spirit lasts its entire duration without being triggered, the cooldown is reset to 1 min.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_holy_guardianspirit.jpg',
		},
		[PriestMajorGlyph.GlyphOfHolyNova]: {
			name: 'Glyph of Holy Nova',
			description: 'Increases the damage and healing of your Holy Nova spell by an additional 20%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_holy_holynova.jpg',
		},
		[PriestMajorGlyph.GlyphOfHymnOfHope]: {
			name: 'Glyph of Hymn of Hope',
			description: 'Your Hymn of Hope lasts an additional 2 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_holy_symbolofhope.jpg',
		},
		[PriestMajorGlyph.GlyphOfInnerFire]: {
			name: 'Glyph of Inner Fire',
			description: 'Increases the armor from your Inner Fire spell by 50%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_holy_innerfire.jpg',
		},
		[PriestMajorGlyph.GlyphOfLightwell]: {
			name: 'Glyph of Lightwell',
			description: 'Increases the amount healed by your Lightwell by 20%',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_holy_summonlightwell.jpg',
		},
		[PriestMajorGlyph.GlyphOfMassDispel]: {
			name: 'Glyph of Mass Dispel',
			description: 'Reduces the mana cost of Mass Dispel by 35%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_arcane_massdispel.jpg',
		},
		[PriestMajorGlyph.GlyphOfMindControl]: {
			name: 'Glyph of Mind Control',
			description: 'Reduces the chance targets will resist or break your Mind Control spell by an additional 17%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_shadowworddominate.jpg',
		},
		[PriestMajorGlyph.GlyphOfMindFlay]: {
			name: 'Glyph of Mind Flay',
			description: 'Increases the damage done by your Mind Flay spell by 10% when your target is afflicted with Shadow Word: Pain.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_siphonmana.jpg',
		},
		[PriestMajorGlyph.GlyphOfMindSear]: {
			name: 'Glyph of Mind Sear',
			description: 'Increases the radius of effect on Mind Sear by 5 yards.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_mindshear.jpg',
		},
		[PriestMajorGlyph.GlyphOfPainSuppression]: {
			name: 'Glyph of Pain Suppression',
			description: 'Allows Pain Suppression to be cast while stunned.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_holy_painsupression.jpg',
		},
		[PriestMajorGlyph.GlyphOfPenance]: {
			name: 'Glyph of Penance',
			description: 'Reduces the cooldown of Penance by 2 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_holy_penance.jpg',
		},
		[PriestMajorGlyph.GlyphOfPowerWordShield]: {
			name: 'Glyph of Power Word: Shield',
			description: 'Your Power Word: Shield also heals the target for 20% of the absorption amount.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_holy_powerwordshield.jpg',
		},
		[PriestMajorGlyph.GlyphOfPrayerOfHealing]: {
			name: 'Glyph of Prayer of Healing',
			description: 'Your Prayer of Healing spell also heals an additional 20% of its initial heal over 6 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_holy_prayerofhealing02.jpg',
		},
		[PriestMajorGlyph.GlyphOfPsychicScream]: {
			name: 'Glyph of Psychic Scream',
			description: 'Increases the duration of your Psychic Scream by 2 sec. and increases its cooldown by 8 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_psychicscream.jpg',
		},
		[PriestMajorGlyph.GlyphOfRenew]: {
			name: 'Glyph of Renew',
			description: 'Reduces the duration of your Renew by 3 sec. but increases the amount healed each tick by 25%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_holy_renew.jpg',
		},
		[PriestMajorGlyph.GlyphOfScourgeImprisonment]: {
			name: 'Glyph of Scourge Imprisonment',
			description: 'Reduces the cast time of your Shackle Undead by 1.0 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_nature_slow.jpg',
		},
		[PriestMajorGlyph.GlyphOfShadow]: {
			name: 'Glyph of Shadow',
			description: 'While in Shadowform, your non-periodic spell critical strikes increase your spell power by 30% of your Spirit for 10 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_holy_mindsooth.jpg',
		},
		[PriestMajorGlyph.GlyphOfShadowWordDeath]: {
			name: 'Glyph of Shadow Word: Death',
			description: 'Targets below 35% health take an additional 10% damage from your Shadow Word: Death spell.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_demonicfortitude.jpg',
		},
		[PriestMajorGlyph.GlyphOfShadowWordPain]: {
			name: 'Glyph of Shadow Word: Pain',
			description: 'The periodic damage ticks of your Shadow Word: Pain spell restore 1% of your base mana.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_shadowwordpain.jpg',
		},
		[PriestMajorGlyph.GlyphOfSmite]: {
			name: 'Glyph of Smite',
			description: 'Your Smite spell inflicts an additional 20% damage against targets afflicted by Holy Fire.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_holy_holysmite.jpg',
		},
		[PriestMajorGlyph.GlyphOfSpiritOfRedemption]: {
			name: 'Glyph of Spirit of Redemption',
			description: 'Increases the duration of Spirit of Redemption by 6 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/inv_enchant_essenceeternallarge.jpg',
		},
	},
	minorGlyphs: {
		[PriestMinorGlyph.GlyphOfFading]: {
			name: 'Glyph of Fading',
			description: 'Reduces the mana cost of your Fade spell by 30%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_magic_lesserinvisibilty.jpg',
		},
		[PriestMinorGlyph.GlyphOfFortitude]: {
			name: 'Glyph of Fortitude',
			description: 'Reduces the mana cost of your Power Word: Fortitude and Prayer of Fortitude spells by 50%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_holy_wordfortitude.jpg',
		},
		[PriestMinorGlyph.GlyphOfLevitate]: {
			name: 'Glyph of Levitate',
			description: 'Your Levitate spell no longer requires a reagent.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_holy_layonhands.jpg',
		},
		[PriestMinorGlyph.GlyphOfShackleUndead]: {
			name: 'Glyph of Shackle Undead',
			description: 'Increases the range of your Shackle Undead spell by 5 yards.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_nature_slow.jpg',
		},
		[PriestMinorGlyph.GlyphOfShadowProtection]: {
			name: 'Glyph of Shadow Protection',
			description: 'Increases the duration of your Shadow Protection and Prayer of Shadow Protection spells by 10 min.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_antishadow.jpg',
		},
		[PriestMinorGlyph.GlyphOfShadowfiend]: {
			name: 'Glyph of Shadowfiend',
			description: 'Receive 5% of your maximum mana if your Shadowfiend dies from damage.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_shadowfiend.jpg',
		},
	},
};
