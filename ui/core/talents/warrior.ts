import { WarriorTalents, WarriorMajorGlyph, WarriorMinorGlyph } from '../proto/warrior.js';

import { GlyphsConfig } from './glyphs_picker.js';
import { TalentsConfig, newTalentsConfig } from './talents_picker.js';

import WarriorTalentJson from './trees/warrior.json';

export const warriorTalentsConfig: TalentsConfig<WarriorTalents> = newTalentsConfig(WarriorTalentJson);

export const warriorGlyphsConfig: GlyphsConfig = {
	majorGlyphs: {
		[WarriorMajorGlyph.GlyphOfBarbaricInsults]: {
			name: 'Glyph of Barbaric Insults',
			description: 'Your Mocking Blow ability generates 100% additional threat.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_warrior_punishingblow.jpg',
		},
		[WarriorMajorGlyph.GlyphOfBladestorm]: {
			name: 'Glyph of Bladestorm',
			description: 'Reduces the cooldown on Bladestorm by 15 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_warrior_bladestorm.jpg',
		},
		[WarriorMajorGlyph.GlyphOfBlocking]: {
			name: 'Glyph of Blocking',
			description: 'Increases your block value by 10% for 10 sec after using your Shield Slam ability.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/inv_shield_05.jpg',
		},
		[WarriorMajorGlyph.GlyphOfBloodthirst]: {
			name: 'Glyph of Bloodthirst',
			description: 'Increases the healing you receive from your Bloodthirst ability by 100%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_nature_bloodlust.jpg',
		},
		[WarriorMajorGlyph.GlyphOfCleaving]: {
			name: 'Glyph of Cleaving',
			description: 'Increases the number of targets your Cleave hits by 1.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_warrior_cleave.jpg',
		},
		[WarriorMajorGlyph.GlyphOfDevastate]: {
			name: 'Glyph of Devastate',
			description: 'Your Devastate ability now applies two stacks of Sunder Armor.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/inv_sword_11.jpg',
		},
		[WarriorMajorGlyph.GlyphOfEnragedRegeneration]: {
			name: 'Glyph of Enraged Regeneration',
			description: 'Your Enraged Regeneration ability heals for an additional 10% of your health over its duration.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_warrior_focusedrage.jpg',
		},
		[WarriorMajorGlyph.GlyphOfExecution]: {
			name: 'Glyph of Execution',
			description: 'Your Execute ability deals damage as if you had 10 additional rage.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/inv_sword_48.jpg',
		},
		[WarriorMajorGlyph.GlyphOfHamstring]: {
			name: 'Glyph of Hamstring',
			description: 'Gives your Hamstring ability a 10% chance to immobilize the target for 5 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_shockwave.jpg',
		},
		[WarriorMajorGlyph.GlyphOfHeroicStrike]: {
			name: 'Glyph of Heroic Strike',
			description: 'You gain 10 rage when you critically strike with your Heroic Strike ability.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_rogue_ambush.jpg',
		},
		[WarriorMajorGlyph.GlyphOfIntervene]: {
			name: 'Glyph of Intervene',
			description: 'Increases the number attacks you intercept for your Intervene target by 1.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_warrior_victoryrush.jpg',
		},
		[WarriorMajorGlyph.GlyphOfLastStand]: {
			name: 'Glyph of Last Stand',
			description: 'Reduces the cooldown of your Last Stand ability by 1 min.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_holy_ashestoashes.jpg',
		},
		[WarriorMajorGlyph.GlyphOfMortalStrike]: {
			name: 'Glyph of Mortal Strike',
			description: 'Increases the damage of your Mortal Strike ability by 10%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_warrior_savageblow.jpg',
		},
		[WarriorMajorGlyph.GlyphOfOverpower]: {
			name: 'Glyph of Overpower',
			description: 'Adds a 100% chance to enable your Overpower when your attacks are parried.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_meleedamage.jpg',
		},
		[WarriorMajorGlyph.GlyphOfRapidCharge]: {
			name: 'Glyph of Rapid Charge',
			description: 'Reduces the cooldown of your Charge ability by 7%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_warrior_charge.jpg',
		},
		[WarriorMajorGlyph.GlyphOfRending]: {
			name: 'Glyph of Rending',
			description: 'Increases the duration of your Rend ability by 6 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_gouge.jpg',
		},
		[WarriorMajorGlyph.GlyphOfResonatingPower]: {
			name: 'Glyph of Resonating Power',
			description: 'Reduces the rage cost of your Thunder Clap ability by 5.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_nature_thunderclap.jpg',
		},
		[WarriorMajorGlyph.GlyphOfRevenge]: {
			name: 'Glyph of Revenge',
			description: 'After using Revenge, your next Heroic Strike costs no rage.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_warrior_revenge.jpg',
		},
		[WarriorMajorGlyph.GlyphOfShieldWall]: {
			name: 'Glyph of Shield Wall',
			description: 'Reduces the cooldown on Shield Wall by 2 min, but Shield Wall now only reduces damage taken by 40%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_warrior_shieldwall.jpg',
		},
		[WarriorMajorGlyph.GlyphOfShockwave]: {
			name: 'Glyph of Shockwave',
			description: 'Reduces the cooldown on Shockwave by 3 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_warrior_shockwave.jpg',
		},
		[WarriorMajorGlyph.GlyphOfSpellReflection]: {
			name: 'Glyph of Spell Reflection',
			description: 'Reduces the cooldown on Spell Reflection by 1 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_warrior_shieldreflection.jpg',
		},
		[WarriorMajorGlyph.GlyphOfSunderArmor]: {
			name: 'Glyph of Sunder Armor',
			description: 'Your Sunder Armor ability affects a second nearby target.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_warrior_sunder.jpg',
		},
		[WarriorMajorGlyph.GlyphOfSweepingStrikes]: {
			name: 'Glyph of Sweeping Strikes',
			description: 'Reduces the rage cost of your Sweeping Strikes ability by 100%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_rogue_slicedice.jpg',
		},
		[WarriorMajorGlyph.GlyphOfTaunt]: {
			name: 'Glyph of Taunt',
			description: 'Increases the chance for your Taunt ability to succeed by 8%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_nature_reincarnation.jpg',
		},
		[WarriorMajorGlyph.GlyphOfVictoryRush]: {
			name: 'Glyph of Victory Rush',
			description: 'Your Victory Rush ability has a 30% increased critical strike chance.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_warrior_devastate.jpg',
		},
		[WarriorMajorGlyph.GlyphOfVigilance]: {
			name: 'Glyph of Vigilance',
			description: 'Your Vigilance ability transfers an additional 5% of your target\'s threat to you.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_warrior_vigilance.jpg',
		},
		[WarriorMajorGlyph.GlyphOfWhirlwind]: {
			name: 'Glyph of Whirlwind',
			description: 'Reduces the cooldown of your Whirlwind by 2 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_whirlwind.jpg',
		},
	},
	minorGlyphs: {
		[WarriorMinorGlyph.GlyphOfBattle]: {
			name: 'Glyph of Battle',
			description: 'Increases the duration of your Battle Shout ability by 2 min.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_warrior_battleshout.jpg',
		},
		[WarriorMinorGlyph.GlyphOfBloodrage]: {
			name: 'Glyph of Bloodrage',
			description: 'Reduces the health cost of your Bloodrage ability by 100%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_racial_bloodrage.jpg',
		},
		[WarriorMinorGlyph.GlyphOfCharge]: {
			name: 'Glyph of Charge',
			description: 'Increases the range of your Charge ability by 5 yards.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_warrior_charge.jpg',
		},
		[WarriorMinorGlyph.GlyphOfCommand]: {
			name: 'Glyph of Command',
			description: 'Increases the duration of your Commanding Shout ability by 2 min.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_warrior_rallyingcry.jpg',
		},
		[WarriorMinorGlyph.GlyphOfEnduringVictory]: {
			name: 'Glyph of Enduring Victory',
			description: 'Increases the window of opportunity in which you can use Victory Rush by 5 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_warrior_devastate.jpg',
		},
		[WarriorMinorGlyph.GlyphOfMockingBlow]: {
			name: 'Glyph of Mocking Blow',
			description: 'Increases the damage of your Mocking Blow ability by 25%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_warrior_punishingblow.jpg',
		},
		[WarriorMinorGlyph.GlyphOfThunderClap]: {
			name: 'Glyph of Thunder Clap',
			description: 'Increases the radius of your Thunder Clap ability by 2 yards.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_nature_thunderclap.jpg',
		},
		[WarriorMinorGlyph.GlyphOfShatteringThrow]: {
			name: 'Glyph of Shattering Throw',
			description: 'Your Shattering Throw is now instant and can be used in any stance, but it no longer removes invulnerabilities and cannot be used on players or player-controlled targets.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_warrior_shatteringthrow.jpg',
		},
	},
};
