import { WarriorTalents, WarriorMajorGlyph, WarriorMinorGlyph } from '../proto/warrior.js';

import { GlyphsConfig } from './glyphs_picker.js';
import { TalentsConfig, newTalentsConfig } from './talents_picker.js';

export const warriorTalentsConfig: TalentsConfig<WarriorTalents> = newTalentsConfig([
    {
        "name": "Arms",
        "backgroundUrl": "https://wow.zamimg.com/images/wow/talents/backgrounds/wrath/161.jpg",
        "talents": [
            {
                "fieldName": "improvedHeroicStrike",
                "location": {
                    "rowIdx": 0,
                    "colIdx": 0
                },
                "spellIds": [
                    12282,
                    12663,
                    12664
                ],
                "maxPoints": 3
            },
            {
                "fieldName": "deflection",
                "location": {
                    "rowIdx": 0,
                    "colIdx": 1
                },
                "spellIds": [
                    16462,
                    16463,
                    16464,
                    16465,
                    16466
                ],
                "maxPoints": 5
            },
            {
                "fieldName": "improvedRend",
                "location": {
                    "rowIdx": 0,
                    "colIdx": 2
                },
                "spellIds": [
                    12286,
                    12658
                ],
                "maxPoints": 2
            },
            {
                "fieldName": "improvedCharge",
                "location": {
                    "rowIdx": 1,
                    "colIdx": 0
                },
                "spellIds": [
                    12285,
                    12697
                ],
                "maxPoints": 2
            },
            {
                "fieldName": "ironWill",
                "location": {
                    "rowIdx": 1,
                    "colIdx": 1
                },
                "spellIds": [
                    12300,
                    12959,
                    12960
                ],
                "maxPoints": 3
            },
            {
                "fieldName": "tacticalMastery",
                "location": {
                    "rowIdx": 1,
                    "colIdx": 2
                },
                "spellIds": [
                    12295,
                    12676,
                    12677
                ],
                "maxPoints": 3
            },
            {
                "fieldName": "improvedOverpower",
                "location": {
                    "rowIdx": 2,
                    "colIdx": 0
                },
                "spellIds": [
                    12290,
                    12963
                ],
                "maxPoints": 2
            },
            {
                "fieldName": "angerManagement",
                "location": {
                    "rowIdx": 2,
                    "colIdx": 1
                },
                "spellIds": [
                    12296
                ],
                "maxPoints": 1
            },
            {
                "fieldName": "impale",
                "location": {
                    "rowIdx": 2,
                    "colIdx": 2
                },
                "spellIds": [
                    16493,
                    16494
                ],
                "maxPoints": 2
            },
            {
                "fieldName": "deepWounds",
                "location": {
                    "rowIdx": 2,
                    "colIdx": 3
                },
                "spellIds": [
                    12834,
                    12849,
                    12867
                ],
                "maxPoints": 3,
                "prereqLocation": {
                    "rowIdx": 2,
                    "colIdx": 2
                }
            },
            {
                "fieldName": "twoHandedWeaponSpecialization",
                "location": {
                    "rowIdx": 3,
                    "colIdx": 1
                },
                "spellIds": [
                    12163,
                    12711,
                    12712
                ],
                "maxPoints": 3
            },
            {
                "fieldName": "tasteForBlood",
                "location": {
                    "rowIdx": 3,
                    "colIdx": 2
                },
                "spellIds": [
                    56636,
                    56637,
                    56638
                ],
                "maxPoints": 3
            },
            {
                "fieldName": "poleaxeSpecialization",
                "location": {
                    "rowIdx": 4,
                    "colIdx": 0
                },
                "spellIds": [
                    12700,
                    12781,
                    12783,
                    12784,
                    12785
                ],
                "maxPoints": 5
            },
            {
                "fieldName": "sweepingStrikes",
                "location": {
                    "rowIdx": 4,
                    "colIdx": 1
                },
                "spellIds": [
                    12328
                ],
                "maxPoints": 1
            },
            {
                "fieldName": "maceSpecialization",
                "location": {
                    "rowIdx": 4,
                    "colIdx": 2
                },
                "spellIds": [
                    12284,
                    12701,
                    12702,
                    12703,
                    12704
                ],
                "maxPoints": 5
            },
            {
                "fieldName": "swordSpecialization",
                "location": {
                    "rowIdx": 4,
                    "colIdx": 3
                },
                "spellIds": [
                    12281,
                    12812,
                    12813,
                    12814,
                    12815
                ],
                "maxPoints": 5
            },
            {
                "fieldName": "weaponMastery",
                "location": {
                    "rowIdx": 5,
                    "colIdx": 0
                },
                "spellIds": [
                    20504,
                    20505
                ],
                "maxPoints": 2
            },
            {
                "fieldName": "improvedHamstring",
                "location": {
                    "rowIdx": 5,
                    "colIdx": 2
                },
                "spellIds": [
                    12289,
                    12668,
                    23695
                ],
                "maxPoints": 3
            },
            {
                "fieldName": "trauma",
                "location": {
                    "rowIdx": 5,
                    "colIdx": 3
                },
                "spellIds": [
                    46854,
                    46855
                ],
                "maxPoints": 2
            },
            {
                "fieldName": "secondWind",
                "location": {
                    "rowIdx": 6,
                    "colIdx": 0
                },
                "spellIds": [
                    29834,
                    29838
                ],
                "maxPoints": 2
            },
            {
                "fieldName": "mortalStrike",
                "location": {
                    "rowIdx": 6,
                    "colIdx": 1
                },
                "spellIds": [
                    12294
                ],
                "maxPoints": 1,
                "prereqLocation": {
                    "rowIdx": 4,
                    "colIdx": 1
                }
            },
            {
                "fieldName": "strengthOfArms",
                "location": {
                    "rowIdx": 6,
                    "colIdx": 2
                },
                "spellIds": [
                    46865,
                    46866
                ],
                "maxPoints": 2
            },
            {
                "fieldName": "improvedSlam",
                "location": {
                    "rowIdx": 6,
                    "colIdx": 3
                },
                "spellIds": [
                    12862,
                    12330
                ],
                "maxPoints": 2
            },
            {
                "fieldName": "juggernaut",
                "location": {
                    "rowIdx": 7,
                    "colIdx": 0
                },
                "spellIds": [
                    64976
                ],
                "maxPoints": 1
            },
            {
                "fieldName": "improvedMortalStrike",
                "location": {
                    "rowIdx": 7,
                    "colIdx": 1
                },
                "spellIds": [
                    35446,
                    35448,
                    35449
                ],
                "maxPoints": 3,
                "prereqLocation": {
                    "rowIdx": 6,
                    "colIdx": 1
                }
            },
            {
                "fieldName": "unrelentingAssault",
                "location": {
                    "rowIdx": 7,
                    "colIdx": 2
                },
                "spellIds": [
                    46859,
                    46860
                ],
                "maxPoints": 2
            },
            {
                "fieldName": "suddenDeath",
                "location": {
                    "rowIdx": 8,
                    "colIdx": 0
                },
                "spellIds": [
                    29723,
                    29724,
                    29725
                ],
                "maxPoints": 3
            },
            {
                "fieldName": "endlessRage",
                "location": {
                    "rowIdx": 8,
                    "colIdx": 1
                },
                "spellIds": [
                    29623
                ],
                "maxPoints": 1
            },
            {
                "fieldName": "bloodFrenzy",
                "location": {
                    "rowIdx": 8,
                    "colIdx": 2
                },
                "spellIds": [
                    29836,
                    29859
                ],
                "maxPoints": 2
            },
            {
                "fieldName": "wreckingCrew",
                "location": {
                    "rowIdx": 9,
                    "colIdx": 1
                },
                "spellIds": [
                    46867,
                    56611,
                    56612,
                    56613,
                    56614
                ],
                "maxPoints": 5
            },
            {
                "fieldName": "bladestorm",
                "location": {
                    "rowIdx": 10,
                    "colIdx": 1
                },
                "spellIds": [
                    46924
                ],
                "maxPoints": 1
            }
        ]
    },
    {
        "name": "Fury",
        "backgroundUrl": "https://wow.zamimg.com/images/wow/talents/backgrounds/wrath/164.jpg",
        "talents": [
            {
                "fieldName": "armoredToTheTeeth",
                "location": {
                    "rowIdx": 0,
                    "colIdx": 0
                },
                "spellIds": [
                    61216,
                    61221,
                    61222
                ],
                "maxPoints": 3
            },
            {
                "fieldName": "boomingVoice",
                "location": {
                    "rowIdx": 0,
                    "colIdx": 1
                },
                "spellIds": [
                    12321,
                    12835
                ],
                "maxPoints": 2
            },
            {
                "fieldName": "cruelty",
                "location": {
                    "rowIdx": 0,
                    "colIdx": 2
                },
                "spellIds": [
                    12320,
					12852,
					12853,
					12855,
					12856
                ],
                "maxPoints": 5
            },
            {
                "fieldName": "improvedDemoralizingShout",
                "location": {
                    "rowIdx": 1,
                    "colIdx": 1
                },
                "spellIds": [
                    12324,
                    12876,
                    12877,
                    12878,
                    12879
                ],
                "maxPoints": 5
            },
            {
                "fieldName": "unbridledWrath",
                "location": {
                    "rowIdx": 1,
                    "colIdx": 2
                },
                "spellIds": [
                    12322,
                    12999,
                    13000,
                    13001,
                    13002
                ],
                "maxPoints": 5
            },
            {
                "fieldName": "improvedCleave",
                "location": {
                    "rowIdx": 2,
                    "colIdx": 0
                },
                "spellIds": [
                    12329,
                    12950,
                    20496
                ],
                "maxPoints": 3
            },
            {
                "fieldName": "piercingHowl",
                "location": {
                    "rowIdx": 2,
                    "colIdx": 1
                },
                "spellIds": [
                    12323
                ],
                "maxPoints": 1
            },
            {
                "fieldName": "bloodCraze",
                "location": {
                    "rowIdx": 2,
                    "colIdx": 2
                },
                "spellIds": [
                    16487,
                    16489,
                    16492
                ],
                "maxPoints": 3
            },
            {
                "fieldName": "commandingPresence",
                "location": {
                    "rowIdx": 2,
                    "colIdx": 3
                },
                "spellIds": [
                    12318,
                    12857,
                    12858,
                    12860,
                    12861
                ],
                "maxPoints": 5
            },
            {
                "fieldName": "dualWieldSpecialization",
                "location": {
                    "rowIdx": 3,
                    "colIdx": 0
                },
                "spellIds": [
                    23584,
                    23585,
                    23586,
                    23587,
                    23588
                ],
                "maxPoints": 5
            },
            {
                "fieldName": "improvedExecute",
                "location": {
                    "rowIdx": 3,
                    "colIdx": 1
                },
                "spellIds": [
                    20502,
                    20503
                ],
                "maxPoints": 2
            },
            {
                "fieldName": "enrage",
                "location": {
                    "rowIdx": 3,
                    "colIdx": 2
                },
                "spellIds": [
                    12317,
                    13045,
                    13046,
                    13047,
                    13048
                ],
                "maxPoints": 5
            },
            {
                "fieldName": "precision",
                "location": {
                    "rowIdx": 4,
                    "colIdx": 0
                },
                "spellIds": [
                    29590,
                    29591,
                    29592
                ],
                "maxPoints": 3
            },
            {
                "fieldName": "deathWish",
                "location": {
                    "rowIdx": 4,
                    "colIdx": 1
                },
                "spellIds": [
                    12292
                ],
                "maxPoints": 1
            },
            {
                "fieldName": "improvedIntercept",
                "location": {
                    "rowIdx": 4,
                    "colIdx": 2
                },
                "spellIds": [
                    29888,
                    29889
                ],
                "maxPoints": 2
            },
            {
                "fieldName": "improvedBerserkerRage",
                "location": {
                    "rowIdx": 5,
                    "colIdx": 0
                },
                "spellIds": [
                    20500,
                    20501
                ],
                "maxPoints": 2
            },
            {
                "fieldName": "flurry",
                "location": {
                    "rowIdx": 5,
                    "colIdx": 2
                },
                "spellIds": [
                    12319,
                    12971,
                    12972,
                    12973,
                    12974
                ],
                "maxPoints": 5
            },
            {
                "fieldName": "intensifyRage",
                "location": {
                    "rowIdx": 6,
                    "colIdx": 0
                },
                "spellIds": [
                    46908,
                    46909,
                    56924
                ],
                "maxPoints": 3
            },
            {
                "fieldName": "bloodthirst",
                "location": {
                    "rowIdx": 6,
                    "colIdx": 1
                },
                "spellIds": [
                    23881
                ],
                "maxPoints": 1,
                "prereqLocation": {
                    "rowIdx": 4,
                    "colIdx": 1
                }
            },
            {
                "fieldName": "improvedWhirlwind",
                "location": {
                    "rowIdx": 6,
                    "colIdx": 3
                },
                "spellIds": [
                    29721,
                    29776
                ],
                "maxPoints": 2
            },
            {
                "fieldName": "furiousAttacks",
                "location": {
                    "rowIdx": 7,
                    "colIdx": 0
                },
                "spellIds": [
                    46910,
                    46911
                ],
                "maxPoints": 2
            },
            {
                "fieldName": "improvedBerserkerStance",
                "location": {
                    "rowIdx": 7,
                    "colIdx": 3
                },
                "spellIds": [
                    29759,
                    29760,
                    29761,
                    29762,
                    29763
                ],
                "maxPoints": 5
            },
            {
                "fieldName": "heroicFury",
                "location": {
                    "rowIdx": 8,
                    "colIdx": 0
                },
                "spellIds": [
                    60970
                ],
                "maxPoints": 1
            },
            {
                "fieldName": "rampage",
                "location": {
                    "rowIdx": 8,
                    "colIdx": 1
                },
                "spellIds": [
                    29801
                ],
                "maxPoints": 1,
                "prereqLocation": {
                    "rowIdx": 6,
                    "colIdx": 1
                }
            },
            {
                "fieldName": "bloodsurge",
                "location": {
                    "rowIdx": 8,
                    "colIdx": 2
                },
                "spellIds": [
                    46913,
                    46914,
                    46915
                ],
                "maxPoints": 3,
                "prereqLocation": {
                    "rowIdx": 6,
                    "colIdx": 1
                }
            },
            {
                "fieldName": "unendingFury",
                "location": {
                    "rowIdx": 9,
                    "colIdx": 1
                },
                "spellIds": [
                    56927,
                    56929,
                    56930,
                    56931,
                    56932
                ],
                "maxPoints": 5
            },
            {
                "fieldName": "titansGrip",
                "location": {
                    "rowIdx": 10,
                    "colIdx": 1
                },
                "spellIds": [
                    46917
                ],
                "maxPoints": 1
            }
        ]
    },
    {
        "name": "Protection",
        "backgroundUrl": "https://wow.zamimg.com/images/wow/talents/backgrounds/wrath/163.jpg",
        "talents": [
            {
                "fieldName": "improvedBloodrage",
                "location": {
                    "rowIdx": 0,
                    "colIdx": 0
                },
                "spellIds": [
                    12301,
                    12818
                ],
                "maxPoints": 2
            },
            {
                "fieldName": "shieldSpecialization",
                "location": {
                    "rowIdx": 0,
                    "colIdx": 1
                },
                "spellIds": [
                    12298,
                    12724,
                    12725,
                    12726,
                    12727
                ],
                "maxPoints": 5
            },
            {
                "fieldName": "improvedThunderClap",
                "location": {
                    "rowIdx": 0,
                    "colIdx": 2
                },
                "spellIds": [
                    12287,
                    12665,
                    12666
                ],
                "maxPoints": 3
            },
            {
                "fieldName": "incite",
                "location": {
                    "rowIdx": 1,
                    "colIdx": 1
                },
                "spellIds": [
                    50685,
                    50686,
                    50687
                ],
                "maxPoints": 3
            },
            {
                "fieldName": "anticipation",
                "location": {
                    "rowIdx": 1,
                    "colIdx": 2
                },
                "spellIds": [
                    12297
                ],
                "maxPoints": 5
            },
            {
                "fieldName": "lastStand",
                "location": {
                    "rowIdx": 2,
                    "colIdx": 0
                },
                "spellIds": [
                    12975
                ],
                "maxPoints": 1
            },
            {
                "fieldName": "improvedRevenge",
                "location": {
                    "rowIdx": 2,
                    "colIdx": 1
                },
                "spellIds": [
                    12797,
                    12799
                ],
                "maxPoints": 2
            },
            {
                "fieldName": "shieldMastery",
                "location": {
                    "rowIdx": 2,
                    "colIdx": 2
                },
                "spellIds": [
                    29598,
                    29599
                ],
                "maxPoints": 2
            },
            {
                "fieldName": "toughness",
                "location": {
                    "rowIdx": 2,
                    "colIdx": 3
                },
                "spellIds": [
                    12299,
                    12761,
                    12762,
                    12763,
                    12764
                ],
                "maxPoints": 5
            },
            {
                "fieldName": "improvedSpellReflection",
                "location": {
                    "rowIdx": 3,
                    "colIdx": 0
                },
                "spellIds": [
                    59088,
                    59089
                ],
                "maxPoints": 2
            },
            {
                "fieldName": "improvedDisarm",
                "location": {
                    "rowIdx": 3,
                    "colIdx": 1
                },
                "spellIds": [
                    12313,
                    12804
                ],
                "maxPoints": 2
            },
            {
                "fieldName": "puncture",
                "location": {
                    "rowIdx": 3,
                    "colIdx": 2
                },
                "spellIds": [
                    12308,
                    12810,
                    12811
                ],
                "maxPoints": 3
            },
            {
                "fieldName": "improvedDisciplines",
                "location": {
                    "rowIdx": 4,
                    "colIdx": 0
                },
                "spellIds": [
                    12312,
                    12803
                ],
                "maxPoints": 2
            },
            {
                "fieldName": "concussionBlow",
                "location": {
                    "rowIdx": 4,
                    "colIdx": 1
                },
                "spellIds": [
                    12809
                ],
                "maxPoints": 1
            },
            {
                "fieldName": "gagOrder",
                "location": {
                    "rowIdx": 4,
                    "colIdx": 2
                },
                "spellIds": [
                    12311,
                    12958
                ],
                "maxPoints": 2
            },
            {
                "fieldName": "oneHandedWeaponSpecialization",
                "location": {
                    "rowIdx": 5,
                    "colIdx": 2
                },
                "spellIds": [
                    16538,
                    16539,
                    16540,
                    16541,
                    16542
                ],
                "maxPoints": 5
            },
            {
                "fieldName": "improvedDefensiveStance",
                "location": {
                    "rowIdx": 6,
                    "colIdx": 0
                },
                "spellIds": [
                    29593,
                    29594
                ],
                "maxPoints": 2
            },
            {
                "fieldName": "vigilance",
                "location": {
                    "rowIdx": 6,
                    "colIdx": 1
                },
                "spellIds": [
                    50720
                ],
                "maxPoints": 1,
                "prereqLocation": {
                    "rowIdx": 4,
                    "colIdx": 1
                }
            },
            {
                "fieldName": "focusedRage",
                "location": {
                    "rowIdx": 6,
                    "colIdx": 2
                },
                "spellIds": [
                    29787,
                    29790,
                    29792
                ],
                "maxPoints": 3
            },
            {
                "fieldName": "vitality",
                "location": {
                    "rowIdx": 7,
                    "colIdx": 1
                },
                "spellIds": [
                    29140,
                    29143,
                    29144
                ],
                "maxPoints": 3
            },
            {
                "fieldName": "safeguard",
                "location": {
                    "rowIdx": 7,
                    "colIdx": 2
                },
                "spellIds": [
                    46945,
                    46949
                ],
                "maxPoints": 2
            },
            {
                "fieldName": "warbringer",
                "location": {
                    "rowIdx": 8,
                    "colIdx": 0
                },
                "spellIds": [
                    57499
                ],
                "maxPoints": 1
            },
            {
                "fieldName": "devastate",
                "location": {
                    "rowIdx": 8,
                    "colIdx": 1
                },
                "spellIds": [
                    20243
                ],
                "maxPoints": 1
            },
            {
                "fieldName": "criticalBlock",
                "location": {
                    "rowIdx": 8,
                    "colIdx": 2
                },
                "spellIds": [
                    47294,
                    47295,
                    47296
                ],
                "maxPoints": 3
            },
            {
                "fieldName": "swordAndBoard",
                "location": {
                    "rowIdx": 9,
                    "colIdx": 1
                },
                "spellIds": [
                    46951,
                    46952,
                    46953
                ],
                "maxPoints": 3,
                "prereqLocation": {
                    "rowIdx": 8,
                    "colIdx": 1
                }
            },
            {
                "fieldName": "damageShield",
                "location": {
                    "rowIdx": 9,
                    "colIdx": 2
                },
                "spellIds": [
                    58872,
                    58874
                ],
                "maxPoints": 2
            },
            {
                "fieldName": "shockwave",
                "location": {
                    "rowIdx": 10,
                    "colIdx": 1
                },
                "spellIds": [
                    46968
                ],
                "maxPoints": 1
            }
        ]
    }
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
