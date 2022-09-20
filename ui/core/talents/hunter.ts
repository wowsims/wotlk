import { HunterTalents, HunterMajorGlyph, HunterMinorGlyph, HunterPetTalents } from '../proto/hunter.js';

import { GlyphsConfig } from './glyphs_picker.js';
import { TalentsConfig, newTalentsConfig } from './talents_picker.js';

export const hunterTalentsConfig: TalentsConfig<HunterTalents> = newTalentsConfig([
	{
		"name": "Beast Mastery",
		"backgroundUrl": "https://wow.zamimg.com/images/wow/talents/backgrounds/wrath/361.jpg",
		"talents": [
			{
				"fieldName": "improvedAspectOfTheHawk",
				"location": {
					"rowIdx": 0,
					"colIdx": 1
				},
				"spellIds": [
					19552,
					19553,
					19554,
					19555,
					19556
				],
				"maxPoints": 5
			},
			{
				"fieldName": "enduranceTraining",
				"location": {
					"rowIdx": 0,
					"colIdx": 2
				},
				"spellIds": [
					19583,
					19584,
					19585,
					19586,
					19587
				],
				"maxPoints": 5
			},
			{
				"fieldName": "focusedFire",
				"location": {
					"rowIdx": 1,
					"colIdx": 0
				},
				"spellIds": [
					35029,
					35030
				],
				"maxPoints": 2
			},
			{
				"fieldName": "improvedAspectOfTheMonkey",
				"location": {
					"rowIdx": 1,
					"colIdx": 1
				},
				"spellIds": [
					19549,
					19550,
					19551
				],
				"maxPoints": 3
			},
			{
				"fieldName": "thickHide",
				"location": {
					"rowIdx": 1,
					"colIdx": 2
				},
				"spellIds": [
					19609,
					19610,
					19612
				],
				"maxPoints": 3
			},
			{
				"fieldName": "improvedRevivePet",
				"location": {
					"rowIdx": 1,
					"colIdx": 3
				},
				"spellIds": [
					24443,
					19575
				],
				"maxPoints": 2
			},
			{
				"fieldName": "pathfinding",
				"location": {
					"rowIdx": 2,
					"colIdx": 0
				},
				"spellIds": [
					19559,
					19560
				],
				"maxPoints": 2
			},
			{
				"fieldName": "aspectMastery",
				"location": {
					"rowIdx": 2,
					"colIdx": 1
				},
				"spellIds": [
					53265
				],
				"maxPoints": 1
			},
			{
				"fieldName": "unleashedFury",
				"location": {
					"rowIdx": 2,
					"colIdx": 2
				},
				"spellIds": [
					19616,
					19617,
					19618,
					19619,
					19620
				],
				"maxPoints": 5
			},
			{
				"fieldName": "improvedMendPet",
				"location": {
					"rowIdx": 3,
					"colIdx": 1
				},
				"spellIds": [
					19572,
					19573
				],
				"maxPoints": 2
			},
			{
				"fieldName": "ferocity",
				"location": {
					"rowIdx": 3,
					"colIdx": 2
				},
				"spellIds": [
					19598,
					19599,
					19600,
					19601,
					19602
				],
				"maxPoints": 5
			},
			{
				"fieldName": "spiritBond",
				"location": {
					"rowIdx": 4,
					"colIdx": 0
				},
				"spellIds": [
					19578,
					20895
				],
				"maxPoints": 2
			},
			{
				"fieldName": "intimidation",
				"location": {
					"rowIdx": 4,
					"colIdx": 1
				},
				"spellIds": [
					19577
				],
				"maxPoints": 1
			},
			{
				"fieldName": "bestialDiscipline",
				"location": {
					"rowIdx": 4,
					"colIdx": 3
				},
				"spellIds": [
					19590,
					19592
				],
				"maxPoints": 2
			},
			{
				"fieldName": "animalHandler",
				"location": {
					"rowIdx": 5,
					"colIdx": 0
				},
				"spellIds": [
					34453,
					34454
				],
				"maxPoints": 2
			},
			{
				"fieldName": "frenzy",
				"location": {
					"rowIdx": 5,
					"colIdx": 2
				},
				"spellIds": [
					19621,
					19622,
					19623,
					19624,
					19625
				],
				"maxPoints": 5,
				"prereqLocation": {
					"rowIdx": 3,
					"colIdx": 2
				}
			},
			{
				"fieldName": "ferociousInspiration",
				"location": {
					"rowIdx": 6,
					"colIdx": 0
				},
				"spellIds": [
					34455,
					34459,
					34460
				],
				"maxPoints": 3
			},
			{
				"fieldName": "bestialWrath",
				"location": {
					"rowIdx": 6,
					"colIdx": 1
				},
				"spellIds": [
					19574
				],
				"maxPoints": 1,
				"prereqLocation": {
					"rowIdx": 4,
					"colIdx": 1
				}
			},
			{
				"fieldName": "catlikeReflexes",
				"location": {
					"rowIdx": 6,
					"colIdx": 2
				},
				"spellIds": [
					34462,
					34464,
					34465
				],
				"maxPoints": 3
			},
			{
				"fieldName": "invigoration",
				"location": {
					"rowIdx": 7,
					"colIdx": 0
				},
				"spellIds": [
					53252,
					53253
				],
				"maxPoints": 2,
				"prereqLocation": {
					"rowIdx": 6,
					"colIdx": 0
				}
			},
			{
				"fieldName": "serpentsSwiftness",
				"location": {
					"rowIdx": 7,
					"colIdx": 2
				},
				"spellIds": [
					34466,
					34467,
					34468,
					34469,
					34470
				],
				"maxPoints": 5
			},
			{
				"fieldName": "longevity",
				"location": {
					"rowIdx": 8,
					"colIdx": 0
				},
				"spellIds": [
					53262,
					53263,
					53264
				],
				"maxPoints": 3
			},
			{
				"fieldName": "theBeastWithin",
				"location": {
					"rowIdx": 8,
					"colIdx": 1
				},
				"spellIds": [
					34692
				],
				"maxPoints": 1,
				"prereqLocation": {
					"rowIdx": 6,
					"colIdx": 1
				}
			},
			{
				"fieldName": "cobraStrikes",
				"location": {
					"rowIdx": 8,
					"colIdx": 2
				},
				"spellIds": [
					53256,
					53259,
					53260
				],
				"maxPoints": 3,
				"prereqLocation": {
					"rowIdx": 7,
					"colIdx": 2
				}
			},
			{
				"fieldName": "kindredSpirits",
				"location": {
					"rowIdx": 9,
					"colIdx": 1
				},
				"spellIds": [
					56314,
					56315,
					56316,
					56317,
					56318
				],
				"maxPoints": 5
			},
			{
				"fieldName": "beastMastery",
				"location": {
					"rowIdx": 10,
					"colIdx": 1
				},
				"spellIds": [
					53270
				],
				"maxPoints": 1
			}
		]
	},
	{
		"name": "Marksmanship",
		"backgroundUrl": "https://wow.zamimg.com/images/wow/talents/backgrounds/wrath/363.jpg",
		"talents": [
			{
				"fieldName": "improvedConcussiveShot",
				"location": {
					"rowIdx": 0,
					"colIdx": 0
				},
				"spellIds": [
					19407,
					19412
				],
				"maxPoints": 2
			},
			{
				"fieldName": "focusedAim",
				"location": {
					"rowIdx": 0,
					"colIdx": 1
				},
				"spellIds": [
					53620,
					53621,
					53622
				],
				"maxPoints": 3
			},
			{
				"fieldName": "lethalShots",
				"location": {
					"rowIdx": 0,
					"colIdx": 2
				},
				"spellIds": [
					19426,
					19427,
					19429,
					19430,
					19431
				],
				"maxPoints": 5
			},
			{
				"fieldName": "carefulAim",
				"location": {
					"rowIdx": 1,
					"colIdx": 0
				},
				"spellIds": [
					34482,
					34483,
					34484
				],
				"maxPoints": 3
			},
			{
				"fieldName": "improvedHuntersMark",
				"location": {
					"rowIdx": 1,
					"colIdx": 1
				},
				"spellIds": [
					19421,
					19422,
					19423
				],
				"maxPoints": 3
			},
			{
				"fieldName": "mortalShots",
				"location": {
					"rowIdx": 1,
					"colIdx": 2
				},
				"spellIds": [
					19485,
					19487,
					19488,
					19489,
					19490
				],
				"maxPoints": 5
			},
			{
				"fieldName": "goForTheThroat",
				"location": {
					"rowIdx": 2,
					"colIdx": 0
				},
				"spellIds": [
					34950,
					34954
				],
				"maxPoints": 2
			},
			{
				"fieldName": "improvedArcaneShot",
				"location": {
					"rowIdx": 2,
					"colIdx": 1
				},
				"spellIds": [
					19454,
					19455,
					19456
				],
				"maxPoints": 3
			},
			{
				"fieldName": "aimedShot",
				"location": {
					"rowIdx": 2,
					"colIdx": 2
				},
				"spellIds": [
					19434
				],
				"maxPoints": 1,
				"prereqLocation": {
					"rowIdx": 1,
					"colIdx": 2
				}
			},
			{
				"fieldName": "rapidKilling",
				"location": {
					"rowIdx": 2,
					"colIdx": 3
				},
				"spellIds": [
					34948,
					34949
				],
				"maxPoints": 2
			},
			{
				"fieldName": "improvedStings",
				"location": {
					"rowIdx": 3,
					"colIdx": 1
				},
				"spellIds": [
					19464,
					19465,
					19466
				],
				"maxPoints": 3
			},
			{
				"fieldName": "efficiency",
				"location": {
					"rowIdx": 3,
					"colIdx": 2
				},
				"spellIds": [
					19416,
					19417,
					19418,
					19419,
					19420
				],
				"maxPoints": 5
			},
			{
				"fieldName": "concussiveBarrage",
				"location": {
					"rowIdx": 4,
					"colIdx": 0
				},
				"spellIds": [
					35100,
					35102
				],
				"maxPoints": 2
			},
			{
				"fieldName": "readiness",
				"location": {
					"rowIdx": 4,
					"colIdx": 1
				},
				"spellIds": [
					23989
				],
				"maxPoints": 1
			},
			{
				"fieldName": "barrage",
				"location": {
					"rowIdx": 4,
					"colIdx": 2
				},
				"spellIds": [
					19461,
					19462,
					24691
				],
				"maxPoints": 3
			},
			{
				"fieldName": "combatExperience",
				"location": {
					"rowIdx": 5,
					"colIdx": 0
				},
				"spellIds": [
					34475,
					34476
				],
				"maxPoints": 2
			},
			{
				"fieldName": "rangedWeaponSpecialization",
				"location": {
					"rowIdx": 5,
					"colIdx": 3
				},
				"spellIds": [
					19507,
					19508,
					19509
				],
				"maxPoints": 3
			},
			{
				"fieldName": "piercingShots",
				"location": {
					"rowIdx": 6,
					"colIdx": 0
				},
				"spellIds": [
					53234,
					53237,
					53238
				],
				"maxPoints": 3
			},
			{
				"fieldName": "trueshotAura",
				"location": {
					"rowIdx": 6,
					"colIdx": 1
				},
				"spellIds": [
					19506
				],
				"maxPoints": 1,
				"prereqLocation": {
					"rowIdx": 4,
					"colIdx": 1
				}
			},
			{
				"fieldName": "improvedBarrage",
				"location": {
					"rowIdx": 6,
					"colIdx": 2
				},
				"spellIds": [
					35104,
					35110,
					35111
				],
				"maxPoints": 3,
				"prereqLocation": {
					"rowIdx": 4,
					"colIdx": 2
				}
			},
			{
				"fieldName": "masterMarksman",
				"location": {
					"rowIdx": 7,
					"colIdx": 1
				},
				"spellIds": [
					34485,
					34486,
					34487,
					34488,
					34489
				],
				"maxPoints": 5
			},
			{
				"fieldName": "rapidRecuperation",
				"location": {
					"rowIdx": 7,
					"colIdx": 2
				},
				"spellIds": [
					53228,
					53232
				],
				"maxPoints": 2
			},
			{
				"fieldName": "wildQuiver",
				"location": {
					"rowIdx": 8,
					"colIdx": 0
				},
				"spellIds": [
					53215,
					53216,
					53217
				],
				"maxPoints": 3
			},
			{
				"fieldName": "silencingShot",
				"location": {
					"rowIdx": 8,
					"colIdx": 1
				},
				"spellIds": [
					34490
				],
				"maxPoints": 1,
				"prereqLocation": {
					"rowIdx": 7,
					"colIdx": 1
				}
			},
			{
				"fieldName": "improvedSteadyShot",
				"location": {
					"rowIdx": 8,
					"colIdx": 2
				},
				"spellIds": [
					53221,
					53222,
					53224
				],
				"maxPoints": 3
			},
			{
				"fieldName": "markedForDeath",
				"location": {
					"rowIdx": 9,
					"colIdx": 1
				},
				"spellIds": [
					53241,
					53243,
					53244,
					53245,
					53246
				],
				"maxPoints": 5
			},
			{
				"fieldName": "chimeraShot",
				"location": {
					"rowIdx": 10,
					"colIdx": 1
				},
				"spellIds": [
					53209
				],
				"maxPoints": 1
			}
		]
	},
	{
		"name": "Survival",
		"backgroundUrl": "https://wow.zamimg.com/images/wow/talents/backgrounds/wrath/362.jpg",
		"talents": [
			{
				"fieldName": "improvedTracking",
				"location": {
					"rowIdx": 0,
					"colIdx": 0
				},
				"spellIds": [
					52783,
					52785,
					52786,
					52787,
					52788
				],
				"maxPoints": 5
			},
			{
				"fieldName": "hawkEye",
				"location": {
					"rowIdx": 0,
					"colIdx": 1
				},
				"spellIds": [
					19498,
					19499,
					19500
				],
				"maxPoints": 3
			},
			{
				"fieldName": "savageStrikes",
				"location": {
					"rowIdx": 0,
					"colIdx": 2
				},
				"spellIds": [
					19159,
					19160
				],
				"maxPoints": 2
			},
			{
				"fieldName": "surefooted",
				"location": {
					"rowIdx": 1,
					"colIdx": 0
				},
				"spellIds": [
					19290,
					19294,
					24283
				],
				"maxPoints": 3
			},
			{
				"fieldName": "entrapment",
				"location": {
					"rowIdx": 1,
					"colIdx": 1
				},
				"spellIds": [
					19184,
					19387,
					19388
				],
				"maxPoints": 3
			},
			{
				"fieldName": "trapMastery",
				"location": {
					"rowIdx": 1,
					"colIdx": 2
				},
				"spellIds": [
					19376,
					63457,
					63458
				],
				"maxPoints": 3
			},
			{
				"fieldName": "survivalInstincts",
				"location": {
					"rowIdx": 1,
					"colIdx": 3
				},
				"spellIds": [
					34494,
					34496
				],
				"maxPoints": 2
			},
			{
				"fieldName": "survivalist",
				"location": {
					"rowIdx": 2,
					"colIdx": 0
				},
				"spellIds": [
					19255,
					19256,
					19257,
					19258,
					19259
				],
				"maxPoints": 5
			},
			{
				"fieldName": "scatterShot",
				"location": {
					"rowIdx": 2,
					"colIdx": 1
				},
				"spellIds": [
					19503
				],
				"maxPoints": 1
			},
			{
				"fieldName": "deflection",
				"location": {
					"rowIdx": 2,
					"colIdx": 2
				},
				"spellIds": [
					19295,
					19297,
					19298
				],
				"maxPoints": 3
			},
			{
				"fieldName": "survivalTactics",
				"location": {
					"rowIdx": 2,
					"colIdx": 3
				},
				"spellIds": [
					19286,
					19287
				],
				"maxPoints": 2
			},
			{
				"fieldName": "tNT",
				"location": {
					"rowIdx": 3,
					"colIdx": 1
				},
				"spellIds": [
					56333,
					56336,
					56337
				],
				"maxPoints": 3
			},
			{
				"fieldName": "lockAndLoad",
				"location": {
					"rowIdx": 3,
					"colIdx": 3
				},
				"spellIds": [
					56342,
					56343,
					56344
				],
				"maxPoints": 3
			},
			{
				"fieldName": "hunterVsWild",
				"location": {
					"rowIdx": 4,
					"colIdx": 0
				},
				"spellIds": [
					56339,
					56340,
					56341
				],
				"maxPoints": 3,
				"prereqLocation": {
					"rowIdx": 2,
					"colIdx": 0
				}
			},
			{
				"fieldName": "killerInstinct",
				"location": {
					"rowIdx": 4,
					"colIdx": 1
				},
				"spellIds": [
					19370,
					19371,
					19373
				],
				"maxPoints": 3
			},
			{
				"fieldName": "counterattack",
				"location": {
					"rowIdx": 4,
					"colIdx": 2
				},
				"spellIds": [
					19306
				],
				"maxPoints": 1,
				"prereqLocation": {
					"rowIdx": 2,
					"colIdx": 2
				}
			},
			{
				"fieldName": "lightningReflexes",
				"location": {
					"rowIdx": 5,
					"colIdx": 0
				},
				"spellIds": [
					19168,
					19180,
					19181,
					24296,
					24297
				],
				"maxPoints": 5
			},
			{
				"fieldName": "resourcefulness",
				"location": {
					"rowIdx": 5,
					"colIdx": 2
				},
				"spellIds": [
					34491,
					34492,
					34493
				],
				"maxPoints": 3
			},
			{
				"fieldName": "exposeWeakness",
				"location": {
					"rowIdx": 6,
					"colIdx": 0
				},
				"spellIds": [
					34500,
					34502,
					34503
				],
				"maxPoints": 3,
				"prereqLocation": {
					"rowIdx": 5,
					"colIdx": 0
				}
			},
			{
				"fieldName": "wyvernSting",
				"location": {
					"rowIdx": 6,
					"colIdx": 1
				},
				"spellIds": [
					19386
				],
				"maxPoints": 1,
				"prereqLocation": {
					"rowIdx": 4,
					"colIdx": 1
				}
			},
			{
				"fieldName": "thrillOfTheHunt",
				"location": {
					"rowIdx": 6,
					"colIdx": 2
				},
				"spellIds": [
					34497,
					34498,
					34499
				],
				"maxPoints": 3
			},
			{
				"fieldName": "masterTactician",
				"location": {
					"rowIdx": 7,
					"colIdx": 0
				},
				"spellIds": [
					34506,
					34507,
					34508,
					34838,
					34839
				],
				"maxPoints": 5
			},
			{
				"fieldName": "noxiousStings",
				"location": {
					"rowIdx": 7,
					"colIdx": 1
				},
				"spellIds": [
					53295,
					53296,
					53297
				],
				"maxPoints": 3,
				"prereqLocation": {
					"rowIdx": 6,
					"colIdx": 1
				}
			},
			{
				"fieldName": "pointOfNoEscape",
				"location": {
					"rowIdx": 8,
					"colIdx": 0
				},
				"spellIds": [
					53298,
					53299
				],
				"maxPoints": 2
			},
			{
				"fieldName": "blackArrow",
				"location": {
					"rowIdx": 8,
					"colIdx": 1
				},
				"spellIds": [
					3674
				],
				"maxPoints": 1
			},
			{
				"fieldName": "sniperTraining",
				"location": {
					"rowIdx": 8,
					"colIdx": 3
				},
				"spellIds": [
					53302,
					53303,
					53304
				],
				"maxPoints": 3
			},
			{
				"fieldName": "huntingParty",
				"location": {
					"rowIdx": 9,
					"colIdx": 2
				},
				"spellIds": [
					53290,
					53291,
					53292
				],
				"maxPoints": 3,
				"prereqLocation": {
					"rowIdx": 6,
					"colIdx": 2
				}
			},
			{
				"fieldName": "explosiveShot",
				"location": {
					"rowIdx": 10,
					"colIdx": 1
				},
				"spellIds": [
					53301
				],
				"maxPoints": 1,
				"prereqLocation": {
					"rowIdx": 8,
					"colIdx": 1
				}
			}
		]
	}
]);

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
