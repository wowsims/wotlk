import { Spec } from '../proto/common.js';
import { DeathknightTalents, DeathknightMajorGlyph, DeathknightMinorGlyph } from '../proto/deathknight.js';
import { Player } from '../player.js';

import { GlyphsConfig } from './glyphs_picker.js';
import { TalentsConfig, newTalentsConfig } from './talents_picker.js';

import DkTalentsJson from './trees/deathknight.json';

export const deathknightTalentsConfig: TalentsConfig<DeathknightTalents> = newTalentsConfig(DkTalentsJson);

export const deathknightGlyphsConfig: GlyphsConfig = {
	majorGlyphs: {
		[DeathknightMajorGlyph.GlyphOfAntiMagicShell]: {
			name: 'Glyph of Anti-Magic Shell',
			description: 'Increases the duration of your Anti-Magic Shell by 2 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_antimagicshell.jpg',
		},
		[DeathknightMajorGlyph.GlyphOfBloodStrike]: {
			name: 'Glyph of Blood Strike',
			description: 'Your Blood Strike causes an additional 20% damage to snared targets.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_deathknight_deathstrike.jpg',
		},
		[DeathknightMajorGlyph.GlyphOfBoneShield]: {
			name: 'Glyph of Bone Shield',
			description: 'Adds 1 additional charge to your Bone Shield.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/inv_chest_leather_13.jpg',
		},
		[DeathknightMajorGlyph.GlyphOfChainsOfIce]: {
			name: 'Glyph of Chains of Ice',
			description: 'Your Chains of Ice also causes 144 to 156 Frost damage, increased by your attack power.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_frost_chainsofice.jpg',
		},
		[DeathknightMajorGlyph.GlyphOfDancingRuneWeapon]: {
			name: 'Glyph of Dancing Rune Weapon',
			description: 'Increases the duration of Dancing Rune Weapon by 5 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/inv_sword_07.jpg',
		},
		[DeathknightMajorGlyph.GlyphOfDarkCommand]: {
			name: 'Glyph of Dark Command',
			description: 'Increases the chance for your Dark Command ability to work successfully by 8%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_nature_shamanrage.jpg',
		},
		[DeathknightMajorGlyph.GlyphOfDarkDeath]: {
			name: 'Glyph of Dark Death',
			description: 'Increases the damage or healing done by Death Coil by 15%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_deathcoil.jpg',
		},
		[DeathknightMajorGlyph.GlyphOfDeathAndDecay]: {
			name: 'Glyph of Death and Decay',
			description: 'Damage of your Death and Decay spell increased by 20%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_deathanddecay.jpg',
		},
		[DeathknightMajorGlyph.GlyphOfDeathGrip]: {
			name: 'Glyph of Death Grip',
			description: 'When you deal a killing blow that grants honor or experience, the cooldown of your Death Grip is refreshed.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_deathknight_strangulate.jpg',
		},
		[DeathknightMajorGlyph.GlyphOfDeathStrike]: {
			name: 'Glyph of Death Strike',
			description: 'Increases your Death Strike\'s damage by 1% for every 1 runic power you currently have (up to a maximum of 25%). The runic power is not consumed by this effect.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_deathknight_butcher2.jpg',
		},
		[DeathknightMajorGlyph.GlyphOfDisease]: {
			name: 'Glyph of Disease',
			description: 'Your Pestilence ability now refreshes disease durations and secondary effects of diseases on your primary target back to their maximum duration.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_plaguecloud.jpg',
		},
		[DeathknightMajorGlyph.GlyphOfFrostStrike]: {
			name: 'Glyph of Frost Strike',
			description: 'Reduces the cost of your Frost Strike by 8 Runic Power.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_deathknight_empowerruneblade2.jpg',
		},
		[DeathknightMajorGlyph.GlyphOfHeartStrike]: {
			name: 'Glyph of Heart Strike',
			description: 'Your Heart Strike also reduces the movement speed of your target by 50% for 10 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/inv_weapon_shortblade_40.jpg',
		},
		[DeathknightMajorGlyph.GlyphOfHowlingBlast]: {
			name: 'Glyph of Howling Blast',
			description: 'Your Howling Blast ability now infects your targets with Frost Fever.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_frost_arcticwinds.jpg',
		},
		[DeathknightMajorGlyph.GlyphOfHungeringCold]: {
			name: 'Glyph of Hungering Cold',
			description: 'Reduces the cost of Hungering Cold by 40 runic power.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/inv_staff_15.jpg',
		},
		[DeathknightMajorGlyph.GlyphOfIceboundFortitude]: {
			name: 'Glyph of Icebound Fortitude',
			description: 'Your Icebound Fortitude now always grants at least 40% damage reduction, regardless of your defense skill.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_deathknight_iceboundfortitude.jpg',
		},
		[DeathknightMajorGlyph.GlyphOfIcyTouch]: {
			name: 'Glyph of Icy Touch',
			description: 'Your Frost Fever disease deals 20% additional damage.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_deathknight_icetouch.jpg',
		},
		[DeathknightMajorGlyph.GlyphOfObliterate]: {
			name: 'Glyph of Obliterate',
			description: 'Increases the damage of your Obliterate ability by 25%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_deathknight_classicon.jpg',
		},
		[DeathknightMajorGlyph.GlyphOfPlagueStrike]: {
			name: 'Glyph of Plague Strike',
			description: 'Your Plague Strike does 20% additional damage.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_deathknight_empowerruneblade.jpg',
		},
		[DeathknightMajorGlyph.GlyphOfRuneStrike]: {
			name: 'Glyph of Rune Strike',
			description: 'Increases the critical strike chance of your Rune Strike by 10%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_deathknight_darkconviction.jpg',
		},
		[DeathknightMajorGlyph.GlyphOfRuneTap]: {
			name: 'Glyph of Rune Tap',
			description: 'Your Rune Tap now heals you for an additional 1% of your maximum health, and also heals your party for 10% of their maximum health.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_deathknight_runetap.jpg',
		},
		[DeathknightMajorGlyph.GlyphOfScourgeStrike]: {
			name: 'Glyph of Scourge Strike',
			description: 'Your Scourge Strike increases the duration of your diseases on the target by 3 sec, up to a maximum of 9 additional seconds.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_deathknight_scourgestrike.jpg',
		},
		[DeathknightMajorGlyph.GlyphOfStrangulate]: {
			name: 'Glyph of Strangulate',
			description: 'Reduces the cooldown of your Strangulate by 20 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_soulleech_3.jpg',
		},
		[DeathknightMajorGlyph.GlyphOfTheGhoul]: {
			name: 'Glyph of the Ghoul',
			description: 'Your Ghoul receives an additional 40% of your Strength and 40% of your Stamina.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_animatedead.jpg',
		},
		[DeathknightMajorGlyph.GlyphOfUnbreakableArmor]: {
			name: 'Glyph of Unbreakable Armor',
			description: 'Increases the total armor granted by Unbreakable Armor to 30%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/inv_armor_helm_plate_naxxramas_raidwarrior_c_01.jpg',
		},
		[DeathknightMajorGlyph.GlyphOfUnholyBlight]: {
			name: 'Glyph of Unholy Blight',
			description: 'Increases the damage done by Unholy Blight by 40%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_contagion.jpg',
		},
		[DeathknightMajorGlyph.GlyphOfVampiricBlood]: {
			name: 'Glyph of Vampiric Blood',
			description: 'Increases the duration of your Vampiric Blood by 5 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_lifedrain.jpg',
		},
	},
	minorGlyphs: {
		[DeathknightMinorGlyph.GlyphOfBloodTap]: {
			name: 'Glyph of Blood Tap',
			description: 'Your Blood Tap no longer causes damage to you.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_deathknight_bloodtap.jpg',
		},
		[DeathknightMinorGlyph.GlyphOfCorpseExplosion]: {
			name: 'Glyph of Corpse Explosion',
			description: 'Increases the radius of effect on Corpse Explosion by 5 yards.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_creature_disease_02.jpg',
		},
		[DeathknightMinorGlyph.GlyphOfDeathSEmbrace]: {
			name: 'Glyph of Death\'s Embrace',
			description: 'Your Death Coil refunds 20 runic power when used to heal.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_deathcoil.jpg',
		},
		[DeathknightMinorGlyph.GlyphOfHornOfWinter]: {
			name: 'Glyph of Horn of Winter',
			description: 'Increases the duration of your Horn of Winter ability by 1 min.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/inv_misc_horn_02.jpg',
		},
		[DeathknightMinorGlyph.GlyphOfPestilence]: {
			name: 'Glyph of Pestilence',
			description: 'Increases the radius of your Pestilence effect by 5 yards.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_plaguecloud.jpg',
		},
		[DeathknightMinorGlyph.GlyphOfRaiseDead]: {
			name: 'Glyph of Raise Dead',
			description: 'Your Raise Dead spell no longer requires a reagent.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_animatedead.jpg',
		},
	},
};
