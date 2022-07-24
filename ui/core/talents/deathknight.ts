import { Spec } from '/wotlk/core/proto/common.js';
import { DeathknightTalents, DeathknightMajorGlyph, DeathknightMinorGlyph } from '/wotlk/core/proto/deathknight.js';
import { Player } from '/wotlk/core/player.js';

import { GlyphsConfig, GlyphsPicker } from './glyphs_picker.js';
import { TalentsConfig, TalentsPicker, newTalentsConfig } from './talents_picker.js';

export const deathknightTalentsConfig: TalentsConfig<DeathknightTalents> = newTalentsConfig([
	{
        "name": "Blood",
        "backgroundUrl": "https://wow.zamimg.com/images/wow/talents/backgrounds/wrath/398.jpg",
        "talents": [
            {
                "fieldName": "butchery",
                "location": {
                    "rowIdx": 0,
                    "colIdx": 0
                },
                "spellIds": [
                    48979,
                    49483
                ],
                "maxPoints": 2
            },
            {
                "fieldName": "subversion",
                "location": {
                    "rowIdx": 0,
                    "colIdx": 1
                },
                "spellIds": [
                    48997,
                    49490,
                    49491
                ],
                "maxPoints": 3
            },
            {
                "fieldName": "bladeBarrier",
                "location": {
                    "rowIdx": 0,
                    "colIdx": 2
                },
                "spellIds": [
                    49182,
                    49500,
                    49501,
                    55225,
                    55226
                ],
                "maxPoints": 5
            },
            {
                "fieldName": "bladedArmor",
                "location": {
                    "rowIdx": 1,
                    "colIdx": 0
                },
                "spellIds": [
                    48978,
                    49390,
                    49391,
                    49392,
                    49393
                ],
                "maxPoints": 5
            },
            {
                "fieldName": "scentOfBlood",
                "location": {
                    "rowIdx": 1,
                    "colIdx": 1
                },
                "spellIds": [
                    49004,
                    49508,
                    49509
                ],
                "maxPoints": 3
            },
            {
                "fieldName": "twoHandedWeaponSpecialization",
                "location": {
                    "rowIdx": 1,
                    "colIdx": 2
                },
                "spellIds": [
                    55107,
                    55108
                ],
                "maxPoints": 2
            },
            {
                "fieldName": "runeTap",
                "location": {
                    "rowIdx": 2,
                    "colIdx": 0
                },
                "spellIds": [
                    48982
                ],
                "maxPoints": 1
            },
            {
                "fieldName": "darkConviction",
                "location": {
                    "rowIdx": 2,
                    "colIdx": 1
                },
                "spellIds": [
                    48987,
                    49477,
                    49478,
                    49479,
                    49480
                ],
                "maxPoints": 5
            },
            {
                "fieldName": "deathRuneMastery",
                "location": {
                    "rowIdx": 2,
                    "colIdx": 2
                },
                "spellIds": [
                    49467,
                    50033,
                    50034
                ],
                "maxPoints": 3
            },
			{
                "fieldName": "improvedRuneTap",
                "location": {
                    "rowIdx": 3,
                    "colIdx": 0
                },
                "spellIds": [
                    48985,
                    49488,
                    49489
                ],
                "maxPoints": 3,
                "prereqLocation": {
                    "rowIdx": 2,
                    "colIdx": 0
                }
            },
            {
                "fieldName": "spellDeflection",
                "location": {
                    "rowIdx": 3,
                    "colIdx": 2
                },
                "spellIds": [
                    49145,
                    49495,
                    49497
                ],
                "maxPoints": 3
            },
            {
                "fieldName": "vendetta",
                "location": {
                    "rowIdx": 3,
                    "colIdx": 3
                },
                "spellIds": [
                    49015,
                    50154,
                    55136
                ],
                "maxPoints": 3
            },
            {
                "fieldName": "bloodyStrikes",
                "location": {
                    "rowIdx": 4,
                    "colIdx": 0
                },
                "spellIds": [
                    48977,
                    49394,
                    49395
                ],
                "maxPoints": 3
            },
            {
                "fieldName": "veteranOfTheThirdWar",
                "location": {
                    "rowIdx": 4,
                    "colIdx": 2
                },
                "spellIds": [
                    49006,
                    49526,
                    50029
                ],
                "maxPoints": 3
            },
            {
                "fieldName": "markOfBlood",
                "location": {
                    "rowIdx": 4,
                    "colIdx": 3
                },
                "spellIds": [
                    49005
                ],
                "maxPoints": 1
            },
			{
                "fieldName": "bloodyVengeance",
                "location": {
                    "rowIdx": 5,
                    "colIdx": 1
                },
                "spellIds": [
                    48988,
                    49503,
                    49504
                ],
                "maxPoints": 3,
                "prereqLocation": {
                    "rowIdx": 2,
                    "colIdx": 1
                }
            },
            {
                "fieldName": "abominationsMight",
                "location": {
                    "rowIdx": 5,
                    "colIdx": 2
                },
                "spellIds": [
                    53137,
                    53138
                ],
                "maxPoints": 2
            },
            {
                "fieldName": "bloodworms",
                "location": {
                    "rowIdx": 6,
                    "colIdx": 0
                },
                "spellIds": [
                    49027,
                    49542,
                    49543
                ],
                "maxPoints": 3
            },
            {
                "fieldName": "hysteria",
                "location": {
                    "rowIdx": 6,
                    "colIdx": 1
                },
                "spellIds": [
                    49016
                ],
                "maxPoints": 1
            },
            {
                "fieldName": "improvedBloodPresence",
                "location": {
                    "rowIdx": 6,
                    "colIdx": 2
                },
                "spellIds": [
                    50365,
                    50371
                ],
                "maxPoints": 2
            },
            {
                "fieldName": "improvedDeathStrike",
                "location": {
                    "rowIdx": 7,
                    "colIdx": 0
                },
                "spellIds": [
                    62905,
                    62908
                ],
                "maxPoints": 2
            },
			{
                "fieldName": "suddenDoom",
                "location": {
                    "rowIdx": 7,
                    "colIdx": 1
                },
                "spellIds": [
                    49018,
                    49529,
                    49530
                ],
                "maxPoints": 3
            },
			{
                "fieldName": "vampiricBlood",
                "location": {
                    "rowIdx": 7,
                    "colIdx": 2
                },
                "spellIds": [
                    55233
                ],
                "maxPoints": 1
            },
            {
                "fieldName": "willOfTheNecropolis",
                "location": {
                    "rowIdx": 8,
                    "colIdx": 0
                },
                "spellIds": [
                    49189,
                    50149,
                    50150
                ],
                "maxPoints": 3
            },
            {
                "fieldName": "heartStrike",
                "location": {
                    "rowIdx": 8,
                    "colIdx": 1
                },
                "spellIds": [
                    55050
                ],
                "maxPoints": 1
            },
            {
                "fieldName": "mightOfMograine",
                "location": {
                    "rowIdx": 8,
                    "colIdx": 2
                },
                "spellIds": [
                    49023,
                    49533,
                    49534
                ],
                "maxPoints": 3
            },
            {
                "fieldName": "bloodGorged",
                "location": {
                    "rowIdx": 9,
                    "colIdx": 1
                },
                "spellIds": [
                    61154,
                    61155,
                    61156,
                    61157,
                    61158
                ],
                "maxPoints": 5
            },
            {
                "fieldName": "dancingRuneWeapon",
                "location": {
                    "rowIdx": 10,
                    "colIdx": 1
                },
                "spellIds": [
                    49028
                ],
                "maxPoints": 1
            },
            
        ]
    },
    {
        "name": "Frost",
        "backgroundUrl": "https://wow.zamimg.com/images/wow/talents/backgrounds/wrath/399.jpg",
        "talents": [
            {
                "fieldName": "improvedIcyTouch",
                "location": {
                    "rowIdx": 0,
                    "colIdx": 0
                },
                "spellIds": [
                    49175,
                    50031,
                    51456
                ],
                "maxPoints": 3
            },
            {
                "fieldName": "runicPowerMastery",
                "location": {
                    "rowIdx": 0,
                    "colIdx": 1
                },
                "spellIds": [
                    49455,
                    50147
                ],
                "maxPoints": 2
            },
            {
                "fieldName": "toughness",
                "location": {
                    "rowIdx": 0,
                    "colIdx": 2
                },
                "spellIds": [
                    49042,
                    49786,
                    49787,
                    49788,
                    49789
                ],
                "maxPoints": 5
            },
            {
                "fieldName": "icyReach",
                "location": {
                    "rowIdx": 1,
                    "colIdx": 1
                },
                "spellIds": [
                    55061,
                    55062
                ],
                "maxPoints": 2
            },
            {
                "fieldName": "blackIce",
                "location": {
                    "rowIdx": 1,
                    "colIdx": 2
                },
                "spellIds": [
                    49140,
                    49661,
                    49662,
                    49663,
                    49664
                ],
                "maxPoints": 5
            },
            {
                "fieldName": "nervesOfColdSteel",
                "location": {
                    "rowIdx": 1,
                    "colIdx": 3
                },
                "spellIds": [
                    49226,
                    50137,
                    50138
                ],
                "maxPoints": 3
            },
            {
                "fieldName": "icyTalons",
                "location": {
                    "rowIdx": 2,
                    "colIdx": 0
                },
                "spellIds": [
                    50880,
                    50884,
                    50885,
                    50886,
                    50887
                ],
                "maxPoints": 5,
                "prereqLocation": {
                    "rowIdx": 0,
                    "colIdx": 0
                }
            },
            {
                "fieldName": "lichborne",
                "location": {
                    "rowIdx": 2,
                    "colIdx": 1
                },
                "spellIds": [
                    49039
                ],
                "maxPoints": 1
            },
            {
                "fieldName": "annihilation",
                "location": {
                    "rowIdx": 2,
                    "colIdx": 2
                },
                "spellIds": [
                    51468,
                    51472,
                    51473
                ],
                "maxPoints": 3
            },
            {
                "fieldName": "killingMachine",
                "location": {
                    "rowIdx": 3,
                    "colIdx": 1
                },
                "spellIds": [
                    51123,
                    51127,
                    51128,
                    51129,
                    51130
                ],
                "maxPoints": 5
            },
            {
                "fieldName": "chillOfTheGrave",
                "location": {
                    "rowIdx": 3,
                    "colIdx": 2
                },
                "spellIds": [
                    49149,
                    50115
                ],
                "maxPoints": 2
            },
            {
                "fieldName": "endlessWinter",
                "location": {
                    "rowIdx": 3,
                    "colIdx": 3
                },
                "spellIds": [
                    49137,
                    49657
                ],
                "maxPoints": 2
            },
            {
                "fieldName": "frigidDreadplate",
                "location": {
                    "rowIdx": 4,
                    "colIdx": 1
                },
                "spellIds": [
                    49186,
                    51108,
                    51109
                ],
                "maxPoints": 3
            },
            {
                "fieldName": "glacierRot",
                "location": {
                    "rowIdx": 4,
                    "colIdx": 2
                },
                "spellIds": [
                    49471,
                    49790,
                    49791
                ],
                "maxPoints": 3
            },
            {
                "fieldName": "deathchill",
                "location": {
                    "rowIdx": 4,
                    "colIdx": 3
                },
                "spellIds": [
                    49796
                ],
                "maxPoints": 1
            },
            {
                "fieldName": "improvedIcyTalons",
                "location": {
                    "rowIdx": 5,
                    "colIdx": 0
                },
                "spellIds": [
                    55610
                ],
                "maxPoints": 1,
                "prereqLocation": {
                    "rowIdx": 2,
                    "colIdx": 0
                }
            },
            {
                "fieldName": "mercilessCombat",
                "location": {
                    "rowIdx": 5,
                    "colIdx": 1
                },
                "spellIds": [
                    49024,
                    49538
                ],
                "maxPoints": 2
            },
            {
                "fieldName": "rime",
                "location": {
                    "rowIdx": 5,
                    "colIdx": 2
                },
                "spellIds": [
                    49188,
                    56822,
                    59057
                ],
                "maxPoints": 3
            },
            {
                "fieldName": "chilblains",
                "location": {
                    "rowIdx": 6,
                    "colIdx": 0
                },
                "spellIds": [
                    50040,
                    50041,
                    50043
                ],
                "maxPoints": 3
            },
            {
                "fieldName": "hungeringCold",
                "location": {
                    "rowIdx": 6,
                    "colIdx": 1
                },
                "spellIds": [
                    49203
                ],
                "maxPoints": 1
            },
            {
                "fieldName": "improvedFrostPresence",
                "location": {
                    "rowIdx": 6,
                    "colIdx": 2
                },
                "spellIds": [
                    50384,
                    50385
                ],
                "maxPoints": 2
            },
            {
                "fieldName": "threatOfThassarian",
                "location": {
                    "rowIdx": 7,
                    "colIdx": 0
                },
                "spellIds": [
                    65661,
                    66191,
                    66192
                ],
                "maxPoints": 3
            },
			{
                "fieldName": "bloodOfTheNorth",
                "location": {
                    "rowIdx": 7,
                    "colIdx": 1
                },
                "spellIds": [
                    54639,
                    54637,
                    54638
                ],
                "maxPoints": 3
            },
            {
                "fieldName": "unbreakableArmor",
                "location": {
                    "rowIdx": 7,
                    "colIdx": 2
                },
                "spellIds": [
                    51271
                ],
                "maxPoints": 1
            },
            {
                "fieldName": "acclimation",
                "location": {
                    "rowIdx": 8,
                    "colIdx": 0
                },
                "spellIds": [
                    49200,
                    50151,
                    50152
                ],
                "maxPoints": 3
            },
            {
                "fieldName": "frostStrike",
                "location": {
                    "rowIdx": 8,
                    "colIdx": 1
                },
                "spellIds": [
                    49143
                ],
                "maxPoints": 1
            },
            {
                "fieldName": "guileOfGorefiend",
                "location": {
                    "rowIdx": 8,
                    "colIdx": 2
                },
                "spellIds": [
                    50187,
                    50190,
                    50191
                ],
                "maxPoints": 3
            },
            {
                "fieldName": "tundraStalker",
                "location": {
                    "rowIdx": 9,
                    "colIdx": 1
                },
                "spellIds": [
                    49202,
                    50127,
                    50128,
                    50129,
                    50130
                ],
                "maxPoints": 5
            },
            {
                "fieldName": "howlingBlast",
                "location": {
                    "rowIdx": 10,
                    "colIdx": 1
                },
                "spellIds": [
                    49184
                ],
                "maxPoints": 1
            },
        ]
    },
    {
        "name": "Unholy",
        "backgroundUrl": "https://wow.zamimg.com/images/wow/talents/backgrounds/wrath/400.jpg",
        "talents": [
            {
                "fieldName": "viciousStrikes",
                "location": {
                    "rowIdx": 0,
                    "colIdx": 0
                },
                "spellIds": [
                    51745,
                    51746
                ],
                "maxPoints": 2
            },
            {
                "fieldName": "virulence",
                "location": {
                    "rowIdx": 0,
                    "colIdx": 1
                },
                "spellIds": [
                    48962,
                    49567,
                    49568
                ],
                "maxPoints": 3
            },
            {
                "fieldName": "anticipation",
                "location": {
                    "rowIdx": 0,
                    "colIdx": 2
                },
                "spellIds": [
                    55129,
                    55130,
                    55131,
                    55132,
                    55133
                ],
                "maxPoints": 5
            },
            {
                "fieldName": "epidemic",
                "location": {
                    "rowIdx": 1,
                    "colIdx": 0
                },
                "spellIds": [
                    49036,
                    49562
                ],
                "maxPoints": 2
            },
            {
                "fieldName": "morbidity",
                "location": {
                    "rowIdx": 1,
                    "colIdx": 1
                },
                "spellIds": [
                    48963,
                    49564,
                    49565
                ],
                "maxPoints": 3
            },
            {
                "fieldName": "unholyCommand",
                "location": {
                    "rowIdx": 1,
                    "colIdx": 2
                },
                "spellIds": [
                    49588,
                    49589
                ],
                "maxPoints": 2
            },
            {
                "fieldName": "ravenousDead",
                "location": {
                    "rowIdx": 1,
                    "colIdx": 3
                },
                "spellIds": [
                    48965,
                    49571,
                    49572
                ],
                "maxPoints": 3
            },
            {
                "fieldName": "outbreak",
                "location": {
                    "rowIdx": 2,
                    "colIdx": 0
                },
                "spellIds": [
                    49013,
                    55236,
                    55237
                ],
                "maxPoints": 3
            },
            {
                "fieldName": "necrosis",
                "location": {
                    "rowIdx": 2,
                    "colIdx": 1
                },
                "spellIds": [
                    51459,
                    51462,
                    51463,
                    51464,
                    51465
                ],
                "maxPoints": 5
            },
            {
                "fieldName": "corpseExplosion",
                "location": {
                    "rowIdx": 2,
                    "colIdx": 2
                },
                "spellIds": [
                    49158
                ],
                "maxPoints": 1
            },
            {
                "fieldName": "onAPaleHorse",
                "location": {
                    "rowIdx": 3,
                    "colIdx": 1
                },
                "spellIds": [
                    49146,
                    51267
                ],
                "maxPoints": 2
            },
            {
                "fieldName": "bloodCakedBlade",
                "location": {
                    "rowIdx": 3,
                    "colIdx": 2
                },
                "spellIds": [
                    49219,
                    49627,
                    49628
                ],
                "maxPoints": 3
            },
            {
                "fieldName": "nightOfTheDead",
                "location": {
                    "rowIdx": 3,
                    "colIdx": 3
                },
                "spellIds": [
                    55620,
                    55623
                ],
                "maxPoints": 2
            },
            {
                "fieldName": "unholyBlight",
                "location": {
                    "rowIdx": 4,
                    "colIdx": 0
                },
                "spellIds": [
                    49194
                ],
                "maxPoints": 1
            },
            {
                "fieldName": "impurity",
                "location": {
                    "rowIdx": 4,
                    "colIdx": 1
                },
                "spellIds": [
                    49220,
                    49633,
                    49635,
                    49636,
                    49638
                ],
                "maxPoints": 5
            },
            {
                "fieldName": "dirge",
                "location": {
                    "rowIdx": 4,
                    "colIdx": 2
                },
                "spellIds": [
                    49223,
                    49599
                ],
                "maxPoints": 2
            },
            {
                "fieldName": "desecration",
                "location": {
                    "rowIdx": 5,
                    "colIdx": 0
                },
                "spellIds": [
                    55666,
                    55667
                ],
                "maxPoints": 2
            },
            {
                "fieldName": "magicSuppression",
                "location": {
                    "rowIdx": 5,
                    "colIdx": 1
                },
                "spellIds": [
                    49224,
                    49610,
                    49611
                ],
                "maxPoints": 3
            },
            {
                "fieldName": "reaping",
                "location": {
                    "rowIdx": 5,
                    "colIdx": 2
                },
                "spellIds": [
                    49208,
                    56834,
                    56835
                ],
                "maxPoints": 3
            },
            {
                "fieldName": "masterOfGhouls",
                "location": {
                    "rowIdx": 5,
                    "colIdx": 3
                },
                "spellIds": [
                    52143
                ],
                "maxPoints": 1,
                "prereqLocation": {
                    "rowIdx": 3,
                    "colIdx": 3
                }
            },
            {
                "fieldName": "desolation",
                "location": {
                    "rowIdx": 6,
                    "colIdx": 0
                },
                "spellIds": [
                    66799,
                    66814,
                    66815,
                    66816,
                    66817
                ],
                "maxPoints": 5
            },
			{
                "fieldName": "antiMagicZone",
                "location": {
                    "rowIdx": 6,
                    "colIdx": 1
                },
                "spellIds": [
                    51052
                ],
                "maxPoints": 1,
                "prereqLocation": {
                    "rowIdx": 5,
                    "colIdx": 1
                }
            },
			{
                "fieldName": "improvedUnholyPresence",
                "location": {
                    "rowIdx": 6,
                    "colIdx": 2
                },
                "spellIds": [
                    50391,
                    50392
                ],
                "maxPoints": 2
            },
            {
                "fieldName": "ghoulFrenzy",
                "location": {
                    "rowIdx": 6,
                    "colIdx": 3
                },
                "spellIds": [
                    63560
                ],
                "maxPoints": 1,
                "prereqLocation": {
                    "rowIdx": 5,
                    "colIdx": 3
                }
            },
            {
                "fieldName": "cryptFever",
                "location": {
                    "rowIdx": 7,
                    "colIdx": 1
                },
                "spellIds": [
                    49032,
                    49631,
                    49632
                ],
                "maxPoints": 3
            },
            {
                "fieldName": "boneShield",
                "location": {
                    "rowIdx": 7,
                    "colIdx": 2
                },
                "spellIds": [
                    49222
                ],
                "maxPoints": 1
            },
            {
                "fieldName": "wanderingPlague",
                "location": {
                    "rowIdx": 8,
                    "colIdx": 0
                },
                "spellIds": [
                    49217,
                    49654,
                    49655
                ],
                "maxPoints": 3
            },
            {
                "fieldName": "ebonPlaguebringer",
                "location": {
                    "rowIdx": 8,
                    "colIdx": 1
                },
                "spellIds": [
                    51099,
                    51160,
                    51161
                ],
                "maxPoints": 3,
                "prereqLocation": {
                    "rowIdx": 7,
                    "colIdx": 1
                }
            },
            {
                "fieldName": "scourgeStrike",
                "location": {
                    "rowIdx": 8,
                    "colIdx": 2
                },
                "spellIds": [
                    55090
                ],
                "maxPoints": 1
            },
            {
                "fieldName": "rageOfRivendare",
                "location": {
                    "rowIdx": 9,
                    "colIdx": 1
                },
                "spellIds": [
                    50117,
                    50118,
                    50119,
                    50120,
                    50121
                ],
                "maxPoints": 5
            },
            {
                "fieldName": "summonGargoyle",
                "location": {
                    "rowIdx": 10,
                    "colIdx": 1
                },
                "spellIds": [
                    49206
                ],
                "maxPoints": 1
            },

        ]
    },
]);

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
