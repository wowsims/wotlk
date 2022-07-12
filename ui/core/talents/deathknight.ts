import { Spec } from '/wotlk/core/proto/common.js';
import { DeathKnightTalents, DeathKnightMajorGlyph, DeathKnightMinorGlyph } from '/wotlk/core/proto/deathknight.js';
import { Player } from '/wotlk/core/player.js';

import { GlyphsConfig, GlyphsPicker } from './glyphs_picker.js';
import { TalentsConfig, TalentsPicker, newTalentsConfig } from './talents_picker.js';

export const deathKnightTalentsConfig: TalentsConfig<DeathKnightTalents> = newTalentsConfig([
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

export const deathKnightGlyphsConfig: GlyphsConfig = {
	majorGlyphs: {},
    minorGlyphs: {},
};
