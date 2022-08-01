import { WarriorTalents, WarriorMajorGlyph, WarriorMinorGlyph } from '../proto/warrior.js';

import { GlyphsConfig } from './glyphs_picker.js';
import { TalentsConfig, newTalentsConfig } from './talents_picker.js';

export const warriorTalentsConfig: TalentsConfig<WarriorTalents> = newTalentsConfig([
	{
		name: 'Arms',
		backgroundUrl: 'https://wow.zamimg.com/images/wow/talents/backgrounds/wrath/161.jpg',
		talents: [
			{
				fieldName: 'improvedHeroicStrike',
				location: {
					rowIdx: 0,
					colIdx: 0,
				},
				spellIds: [12282, 12663],
				maxPoints: 3,
			},
			{
				fieldName: 'deflection',
				location: {
					rowIdx: 0,
					colIdx: 1,
				},
				spellIds: [16462],
				maxPoints: 5,
			},
			{
				fieldName: 'improvedRend',
				location: {
					rowIdx: 0,
					colIdx: 2,
				},
				spellIds: [12286, 12658],
				maxPoints: 3,
			},
			{
				fieldName: 'improvedCharge',
				location: {
					rowIdx: 1,
					colIdx: 0,
				},
				spellIds: [12285, 12697],
				maxPoints: 2,
			},
			{
				//fieldName: 'ironWill',
				location: {
					rowIdx: 1,
					colIdx: 1,
				},
				spellIds: [12300, 12959],
				maxPoints: 5,
			},
			{
				fieldName: 'improvedThunderClap',
				location: {
					rowIdx: 1,
					colIdx: 2,
				},
				spellIds: [12287, 12665],
				maxPoints: 3,
			},
			{
				fieldName: 'improvedOverpower',
				location: {
					rowIdx: 2,
					colIdx: 0,
				},
				spellIds: [12290, 12963],
				maxPoints: 2,
			},
			{
				fieldName: 'angerManagement',
				location: {
					rowIdx: 2,
					colIdx: 1,
				},
				spellIds: [12296],
				maxPoints: 1,
			},
			{
				fieldName: 'deepWounds',
				location: {
					rowIdx: 2,
					colIdx: 2,
				},
				spellIds: [12834, 12849, 12867],
				maxPoints: 3,
			},
			{
				fieldName: 'twoHandedWeaponSpecialization',
				location: {
					rowIdx: 3,
					colIdx: 1,
				},
				spellIds: [12163, 12711],
				maxPoints: 5,
			},
			{
				fieldName: 'impale',
				location: {
					rowIdx: 3,
					colIdx: 2,
				},
				prereqLocation: {
					rowIdx: 2,
					colIdx: 2,
				},
				spellIds: [16493],
				maxPoints: 2,
			},
			{
				fieldName: 'poleaxeSpecialization',
				location: {
					rowIdx: 4,
					colIdx: 0,
				},
				spellIds: [12700, 12781, 12783],
				maxPoints: 5,
			},
			{
				fieldName: 'deathWish',
				location: {
					rowIdx: 4,
					colIdx: 1,
				},
				spellIds: [12292],
				maxPoints: 1,
			},
			{
				fieldName: 'maceSpecialization',
				location: {
					rowIdx: 4,
					colIdx: 2,
				},
				spellIds: [12284, 12701],
				maxPoints: 5,
			},
			{
				fieldName: 'swordSpecialization',
				location: {
					rowIdx: 4,
					colIdx: 3,
				},
				spellIds: [12281, 12812],
				maxPoints: 5,
			},
			{
				//fieldName: 'improvedIntercept',
				location: {
					rowIdx: 5,
					colIdx: 0,
				},
				spellIds: [29888],
				maxPoints: 2,
			},
			{
				//fieldName: 'improvedHamstring',
				location: {
					rowIdx: 5,
					colIdx: 2,
				},
				spellIds: [12289, 12668, 23695],
				maxPoints: 3,
			},
			{
				fieldName: 'improvedDisciplines',
				location: {
					rowIdx: 5,
					colIdx: 3,
				},
				spellIds: [29723],
				maxPoints: 3,
			},
			{
				fieldName: 'bloodFrenzy',
				location: {
					rowIdx: 6,
					colIdx: 0,
				},
				spellIds: [29836, 29859],
				maxPoints: 2,
			},
			{
				fieldName: 'mortalStrike',
				location: {
					rowIdx: 6,
					colIdx: 1,
				},
				prereqLocation: {
					rowIdx: 4,
					colIdx: 1,
				},
				spellIds: [12294],
				maxPoints: 1,
			},
			{
				//fieldName: 'secondWind',
				location: {
					rowIdx: 6,
					colIdx: 2,
				},
				spellIds: [29834, 29838],
				maxPoints: 2,
			},
			{
				fieldName: 'improvedMortalStrike',
				location: {
					rowIdx: 7,
					colIdx: 1,
				},
				prereqLocation: {
					rowIdx: 6,
					colIdx: 1,
				},
				spellIds: [35446, 35448],
				maxPoints: 5,
			},
			{
				fieldName: 'endlessRage',
				location: {
					rowIdx: 8,
					colIdx: 1,
				},
				spellIds: [29623],
				maxPoints: 1,
			},
		],
	},
	{
		name: 'Fury',
		backgroundUrl: 'https://wow.zamimg.com/images/wow/talents/backgrounds/wrath/164.jpg',
		talents: [
			{
				fieldName: 'boomingVoice',
				location: {
					rowIdx: 0,
					colIdx: 1,
				},
				spellIds: [12321, 12835],
				maxPoints: 5,
			},
			{
				fieldName: 'cruelty',
				location: {
					rowIdx: 0,
					colIdx: 2,
				},
				spellIds: [12320, 12852, 12853, 12855],
				maxPoints: 5,
			},
			{
				fieldName: 'improvedDemoralizingShout',
				location: {
					rowIdx: 1,
					colIdx: 1,
				},
				spellIds: [12324, 12876],
				maxPoints: 5,
			},
			{
				fieldName: 'unbridledWrath',
				location: {
					rowIdx: 1,
					colIdx: 2,
				},
				spellIds: [12322, 12999],
				maxPoints: 5,
			},
			{
				fieldName: 'improvedCleave',
				location: {
					rowIdx: 2,
					colIdx: 0,
				},
				spellIds: [12329, 12950, 20496],
				maxPoints: 3,
			},
			{
				//fieldName: 'piercingHowl',
				location: {
					rowIdx: 2,
					colIdx: 1,
				},
				spellIds: [12323],
				maxPoints: 1,
			},
			{
				//fieldName: 'bloodCraze',
				location: {
					rowIdx: 2,
					colIdx: 2,
				},
				spellIds: [16487, 16489, 16492],
				maxPoints: 3,
			},
			{
				fieldName: 'commandingPresence',
				location: {
					rowIdx: 2,
					colIdx: 3,
				},
				spellIds: [12318, 12857, 12858, 12860],
				maxPoints: 5,
			},
			{
				fieldName: 'dualWieldSpecialization',
				location: {
					rowIdx: 3,
					colIdx: 0,
				},
				spellIds: [23584],
				maxPoints: 5,
			},
			{
				fieldName: 'improvedExecute',
				location: {
					rowIdx: 3,
					colIdx: 1,
				},
				spellIds: [20502],
				maxPoints: 2,
			},
			{
				//fieldName: 'enrage',
				location: {
					rowIdx: 3,
					colIdx: 2,
				},
				spellIds: [12317, 13045],
				maxPoints: 5,
			},
			{
				fieldName: 'improvedSlam',
				location: {
					rowIdx: 4,
					colIdx: 0,
				},
				spellIds: [12862, 12330],
				maxPoints: 2,
			},
			{
				fieldName: 'sweepingStrikes',
				location: {
					rowIdx: 4,
					colIdx: 1,
				},
				spellIds: [12328],
				maxPoints: 1,
			},
			{
				fieldName: 'weaponMastery',
				location: {
					rowIdx: 4,
					colIdx: 3,
				},
				spellIds: [20504],
				maxPoints: 2,
			},
			{
				fieldName: 'improvedBerserkerRage',
				location: {
					rowIdx: 5,
					colIdx: 0,
				},
				spellIds: [20500],
				maxPoints: 2,
			},
			{
				fieldName: 'flurry',
				location: {
					rowIdx: 5,
					colIdx: 2,
				},
				prereqLocation: {
					rowIdx: 3,
					colIdx: 2,
				},
				spellIds: [12319, 12971],
				maxPoints: 5,
			},
			{
				fieldName: 'precision',
				location: {
					rowIdx: 6,
					colIdx: 0,
				},
				spellIds: [29590],
				maxPoints: 3,
			},
			{
				fieldName: 'bloodthirst',
				location: {
					rowIdx: 6,
					colIdx: 1,
				},
				prereqLocation: {
					rowIdx: 4,
					colIdx: 1,
				},
				spellIds: [23881],
				maxPoints: 1,
			},
			{
				fieldName: 'improvedWhirlwind',
				location: {
					rowIdx: 6,
					colIdx: 2,
				},
				spellIds: [29721, 29776],
				maxPoints: 2,
			},
			{
				fieldName: 'improvedBerserkerStance',
				location: {
					rowIdx: 7,
					colIdx: 2,
				},
				spellIds: [29759],
				maxPoints: 5,
			},
			{
				fieldName: 'rampage',
				location: {
					rowIdx: 8,
					colIdx: 1,
				},
				prereqLocation: {
					rowIdx: 6,
					colIdx: 1,
				},
				spellIds: [29801],
				maxPoints: 1,
			},
		],
	},
	{
		name: 'Protection',
		backgroundUrl: 'https://wow.zamimg.com/images/wow/talents/backgrounds/wrath/163.jpg',
		talents: [
			{
				fieldName: 'improvedBloodrage',
				location: {
					rowIdx: 0,
					colIdx: 0,
				},
				spellIds: [12301, 12818],
				maxPoints: 2,
			},
			{
				fieldName: 'tacticalMastery',
				location: {
					rowIdx: 0,
					colIdx: 1,
				},
				spellIds: [12295, 12676],
				maxPoints: 3,
			},
			{
				fieldName: 'anticipation',
				location: {
					rowIdx: 0,
					colIdx: 2,
				},
				spellIds: [12297, 12750],
				maxPoints: 5,
			},
			{
				fieldName: 'shieldSpecialization',
				location: {
					rowIdx: 1,
					colIdx: 1,
				},
				spellIds: [12298, 12724],
				maxPoints: 5,
			},
			{
				fieldName: 'toughness',
				location: {
					rowIdx: 1,
					colIdx: 2,
				},
				spellIds: [12299, 12761],
				maxPoints: 5,
			},
			{
				fieldName: 'lastStand',
				location: {
					rowIdx: 2,
					colIdx: 0,
				},
				spellIds: [12975],
				maxPoints: 1,
			},
			{
				fieldName: 'improvedShieldBlock',
				location: {
					rowIdx: 2,
					colIdx: 1,
				},
				prereqLocation: {
					rowIdx: 1,
					colIdx: 1,
				},
				spellIds: [12945],
				maxPoints: 1,
			},
			{
				//fieldName: 'improvedRevenge',
				location: {
					rowIdx: 2,
					colIdx: 2,
				},
				spellIds: [12797, 12799],
				maxPoints: 3,
			},
			{
				fieldName: 'defiance',
				location: {
					rowIdx: 2,
					colIdx: 3,
				},
				spellIds: [12303, 12788],
				maxPoints: 3,
			},
			{
				fieldName: 'improvedSunderArmor',
				location: {
					rowIdx: 3,
					colIdx: 0,
				},
				spellIds: [12308, 12810],
				maxPoints: 3,
			},
			{
				//fieldName: 'improvedDisarm',
				location: {
					rowIdx: 3,
					colIdx: 1,
				},
				spellIds: [12313, 12804, 12807],
				maxPoints: 3,
			},
			{
				//fieldName: 'improvedTaunt',
				location: {
					rowIdx: 3,
					colIdx: 2,
				},
				spellIds: [12302, 12765],
				maxPoints: 2,
			},
			{
				fieldName: 'improvedShieldWall',
				location: {
					rowIdx: 4,
					colIdx: 0,
				},
				spellIds: [12312, 12803],
				maxPoints: 2,
			},
			{
				//fieldName: 'concussionBlow',
				location: {
					rowIdx: 4,
					colIdx: 1,
				},
				spellIds: [12809],
				maxPoints: 1,
			},
			{
				//fieldName: 'improvedShieldBash',
				location: {
					rowIdx: 4,
					colIdx: 2,
				},
				spellIds: [12311, 12958],
				maxPoints: 2,
			},
			{
				fieldName: 'shieldMastery',
				location: {
					rowIdx: 5,
					colIdx: 0,
				},
				spellIds: [29598],
				maxPoints: 3,
			},
			{
				fieldName: 'oneHandedWeaponSpecialization',
				location: {
					rowIdx: 5,
					colIdx: 2,
				},
				spellIds: [16538],
				maxPoints: 5,
			},
			{
				fieldName: 'improvedDefensiveStance',
				location: {
					rowIdx: 6,
					colIdx: 0,
				},
				spellIds: [29593],
				maxPoints: 3,
			},
			{
				fieldName: 'shieldSlam',
				location: {
					rowIdx: 6,
					colIdx: 1,
				},
				prereqLocation: {
					rowIdx: 4,
					colIdx: 1,
				},
				spellIds: [23922],
				maxPoints: 1,
			},
			{
				fieldName: 'focusedRage',
				location: {
					rowIdx: 6,
					colIdx: 2,
				},
				spellIds: [29787, 29790, 29792],
				maxPoints: 3,
			},
			{
				fieldName: 'vitality',
				location: {
					rowIdx: 7,
					colIdx: 1,
				},
				spellIds: [29140, 29143],
				maxPoints: 5,
			},
			{
				fieldName: 'devastate',
				location: {
					rowIdx: 8,
					colIdx: 1,
				},
				spellIds: [20243],
				maxPoints: 1,
			},
		],
	},
]);

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
    },
};
