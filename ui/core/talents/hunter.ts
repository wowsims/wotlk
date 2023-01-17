import { HunterTalents, HunterMajorGlyph, HunterMinorGlyph, HunterPetTalents } from '../proto/hunter.js';

import { GlyphsConfig } from './glyphs_picker.js';
import { TalentsConfig, newTalentsConfig } from './talents_picker.js';

import HunterTalentJson from './trees/hunter.json';

export const hunterTalentsConfig: TalentsConfig<HunterTalents> = newTalentsConfig(HunterTalentJson);

export const hunterGlyphsConfig: GlyphsConfig = {
	majorGlyphs: {
		[HunterMajorGlyph.GlyphOfAimedShot]: {
			name: 'Glyph of Aimed Shot',
			description: 'Reduces the cooldown of your Aimed Shot ability by 2 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/inv_spear_07.jpg',
		},
		[HunterMajorGlyph.GlyphOfArcaneShot]: {
			name: 'Glyph of Arcane Shot',
			description: 'Your Arcane Shot refunds 20% of its mana cost if the target has one of your Stings active on it.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_impalingbolt.jpg',
		},
		[HunterMajorGlyph.GlyphOfAspectOfTheViper]: {
			name: 'Glyph of Aspect of the Viper',
			description: 'Increases the amount of mana gained from attacks while Aspect of the Viper is active by 10%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_hunter_aspectoftheviper.jpg',
		},
		[HunterMajorGlyph.GlyphOfBestialWrath]: {
			name: 'Glyph of Bestial Wrath',
			description: 'Decreases the cooldown of Bestial Wrath by 20 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_druid_ferociousbite.jpg',
		},
		[HunterMajorGlyph.GlyphOfChimeraShot]: {
			name: 'Glyph of Chimera Shot',
			description: 'Reduces the cooldown of Chimera Shot by 1 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_hunter_chimerashot2.jpg',
		},
		[HunterMajorGlyph.GlyphOfDeterrence]: {
			name: 'Glyph of Deterrence',
			description: 'Decreases the cooldown of Deterrence by 10 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_whirlwind.jpg',
		},
		[HunterMajorGlyph.GlyphOfDisengage]: {
			name: 'Glyph of Disengage',
			description: 'Decreases the cooldown of Disengage by 5 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_rogue_feint.jpg',
		},
		[HunterMajorGlyph.GlyphOfExplosiveShot]: {
			name: 'Glyph of Explosive Shot',
			description: 'Increases the critical strike chance of Explosive Shot by 4%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_hunter_explosiveshot.jpg',
		},
		[HunterMajorGlyph.GlyphOfExplosiveTrap]: {
			name: 'Glyph of Explosive Trap',
			description: 'The periodic damage from your Explosive Trap can now be critical strikes.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_fire_selfdestruct.jpg',
		},
		[HunterMajorGlyph.GlyphOfFreezingTrap]: {
			name: 'Glyph of Freezing Trap',
			description: 'When your Freezing Trap breaks, the victim\'s movement speed is reduced by 30% for 4 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_frost_chainsofice.jpg',
		},
		[HunterMajorGlyph.GlyphOfFrostTrap]: {
			name: 'Glyph of Frost Trap',
			description: 'Increases the radius of the effect from your Frost Trap by 2 yards.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_frost_freezingbreath.jpg',
		},
		[HunterMajorGlyph.GlyphOfHuntersMark]: {
			name: 'Glyph of Hunter\'s Mark',
			description: 'Increases the attack power bonus of your Hunter\'s Mark by 20%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_hunter_snipershot.jpg',
		},
		[HunterMajorGlyph.GlyphOfImmolationTrap]: {
			name: 'Glyph of Immolation Trap',
			description: 'Decreases the duration of the effect from your Immolation Trap by 6 sec., but damage while active is increased by 100%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_fire_flameshock.jpg',
		},
		[HunterMajorGlyph.GlyphOfKillShot]: {
			name: 'Glyph of Kill Shot',
			description: 'Reduces the cooldown of Kill Shot by 6 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_hunter_assassinate2.jpg',
		},
		[HunterMajorGlyph.GlyphOfMending]: {
			name: 'Glyph of Mending',
			description: 'Increases the healing done by your Mend Pet ability by 40%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_hunter_mendpet.jpg',
		},
		[HunterMajorGlyph.GlyphOfMultiShot]: {
			name: 'Glyph of Multi-Shot',
			description: 'Decreases the cooldown of Multi-Shot by 1 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_upgrademoonglaive.jpg',
		},
		[HunterMajorGlyph.GlyphOfRapidFire]: {
			name: 'Glyph of Rapid Fire',
			description: 'Increases the haste from Rapid Fire by an additional 8%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_hunter_runningshot.jpg',
		},
		[HunterMajorGlyph.GlyphOfRaptorStrike]: {
			name: 'Glyph of Raptor Strike',
			description: 'Reduces damage taken by 20% for 3 sec after using Raptor Strike.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_meleedamage.jpg',
		},
		[HunterMajorGlyph.GlyphOfScatterShot]: {
			name: 'Glyph of Scatter Shot',
			description: 'Increases the range of Scatter Shot by 3 yards.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_golemstormbolt.jpg',
		},
		[HunterMajorGlyph.GlyphOfSerpentSting]: {
			name: 'Glyph of Serpent Sting',
			description: 'Increases the duration of your Serpent Sting by 6 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_hunter_quickshot.jpg',
		},
		[HunterMajorGlyph.GlyphOfSnakeTrap]: {
			name: 'Glyph of Snake Trap',
			description: 'Snakes from your Snake Trap take 90% reduced damage from area of effect spells.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_hunter_snaketrap.jpg',
		},
		[HunterMajorGlyph.GlyphOfSteadyShot]: {
			name: 'Glyph of Steady Shot',
			description: 'Increases the damage dealt by Steady Shot by 10% when your target is afflicted with Serpent Sting.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_hunter_steadyshot.jpg',
		},
		[HunterMajorGlyph.GlyphOfTheBeast]: {
			name: 'Glyph of the Beast',
			description: 'Increases the attack power bonus of Aspect of the Beast for you and your pet by an additional 2%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_mount_pinktiger.jpg',
		},
		[HunterMajorGlyph.GlyphOfTheHawk]: {
			name: 'Glyph of the Hawk',
			description: 'Increases the haste bonus of the Improved Aspect of the Hawk effect by an additional 6%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_nature_ravenform.jpg',
		},
		[HunterMajorGlyph.GlyphOfTrueshotAura]: {
			name: 'Glyph of Trueshot Aura',
			description: 'While your Trueshot Aura is active, you have 10% increased critical strike chance on your Aimed Shot.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_trueshot.jpg',
		},
		[HunterMajorGlyph.GlyphOfVolley]: {
			name: 'Glyph of Volley',
			description: 'Decreases the mana cost of Volley by 20%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_marksmanship.jpg',
		},
		[HunterMajorGlyph.GlyphOfWyvernSting]: {
			name: 'Glyph of Wyvern Sting',
			description: 'Decreases the cooldown of your Wyvern Sting by 6 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/inv_spear_02.jpg',
		},
	},
	minorGlyphs: {
		[HunterMinorGlyph.GlyphOfFeignDeath]: {
			name: 'Glyph of Feign Death',
			description: 'Reduces the cooldown of your Feign Death spell by 5 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_rogue_feigndeath.jpg',
		},
		[HunterMinorGlyph.GlyphOfMendPet]: {
			name: 'Glyph of Mend Pet',
			description: 'Your Mend Pet spell increases your pet\'s happiness slightly.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_hunter_mendpet.jpg',
		},
		[HunterMinorGlyph.GlyphOfPossessedStrength]: {
			name: 'Glyph of Possessed Strength',
			description: 'Increases the damage your pet inflicts while using Eyes of the Beast by 50%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_eyeoftheowl.jpg',
		},
		[HunterMinorGlyph.GlyphOfRevivePet]: {
			name: 'Glyph of Revive Pet',
			description: 'Reduces the pushback suffered from damaging attacks while casting Revive Pet by 100%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_hunter_beastsoothe.jpg',
		},
		[HunterMinorGlyph.GlyphOfScareBeast]: {
			name: 'Glyph of Scare Beast',
			description: 'Reduces the pushback suffered from damaging attacks while casting Scare Beast by 75%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_druid_cower.jpg',
		},
		[HunterMinorGlyph.GlyphOfThePack]: {
			name: 'Glyph of the Pack',
			description: 'Increases the range of your Aspect of the Pack ability by 15 yards.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_mount_jungletiger.jpg',
		},
	},
};
