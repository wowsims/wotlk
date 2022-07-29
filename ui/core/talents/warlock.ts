import { WarlockTalents, WarlockMajorGlyph, WarlockMinorGlyph } from '/wotlk/core/proto/warlock.js';

import { GlyphsConfig } from './glyphs_picker.js';
import { TalentsConfig, newTalentsConfig } from './talents_picker.js';

export const warlockTalentsConfig: TalentsConfig<WarlockTalents> = newTalentsConfig([
	{
		name: "Affliction",
		backgroundUrl: "https://wow.zamimg.com/images/wow/talents/backgrounds/wrath/302.jpg",
		talents: [
			{
				fieldName: "improvedCurseOfAgony",
				location: {
					rowIdx: 0,
					colIdx: 0
				},
				spellIds: [
					18827,
					18829
				],
				maxPoints: 2
			},
			{
				fieldName: "suppression",
				location: {
					rowIdx: 0,
					colIdx: 1
				},
				spellIds: [
					18174,
					18175,
					18176
				],
				maxPoints: 3
			},
			{
				fieldName: "improvedCorruption",
				location: {
					rowIdx: 0,
					colIdx: 2
				},
				spellIds: [
					17810,
					17811,
					17812,
					17813,
					17814
				],
				maxPoints: 5
			},
			{
				fieldName: "improvedCurseOfWeakness",
				location: {
					rowIdx: 1,
					colIdx: 0
				},
				spellIds: [
					18179,
					18180
				],
				maxPoints: 2
			},
			{
				fieldName: "improvedDrainSoul",
				location: {
					rowIdx: 1,
					colIdx: 1
				},
				spellIds: [
					18213,
					18372
				],
				maxPoints: 2
			},
			{
				fieldName: "improvedLifeTap",
				location: {
					rowIdx: 1,
					colIdx: 2
				},
				spellIds: [
					18182,
					18183
				],
				maxPoints: 2
			},
			{
				fieldName: "soulSiphon",
				location: {
					rowIdx: 1,
					colIdx: 3
				},
				spellIds: [
					17804,
					17805
				],
				maxPoints: 2
			},
			{
				fieldName: "improvedFear",
				location: {
					rowIdx: 2,
					colIdx: 0
				},
				spellIds: [
					53754,
					53759
				],
				maxPoints: 2
			},
			{
				fieldName: "felConcentration",
				location: {
					rowIdx: 2,
					colIdx: 1
				},
				spellIds: [
					17783,
					17784,
					17785
				],
				maxPoints: 3
			},
			{
				fieldName: "amplifyCurse",
				location: {
					rowIdx: 2,
					colIdx: 2
				},
				spellIds: [
					18288
				],
				maxPoints: 1
			},
			{
				fieldName: "grimReach",
				location: {
					rowIdx: 3,
					colIdx: 0
				},
				spellIds: [
					18218,
					18219
				],
				maxPoints: 2
			},
			{
				fieldName: "nightfall",
				location: {
					rowIdx: 3,
					colIdx: 1
				},
				spellIds: [
					18094,
					18095
				],
				maxPoints: 2
			},
			{
				fieldName: "empoweredCorruption",
				location: {
					rowIdx: 3,
					colIdx: 3
				},
				spellIds: [
					32381,
					32382,
					32383
				],
				maxPoints: 3
			},
			{
				fieldName: "shadowEmbrace",
				location: {
					rowIdx: 4,
					colIdx: 0
				},
				spellIds: [
					32385,
					32387,
					32392,
					32393,
					32394
				],
				maxPoints: 5
			},
			{
				fieldName: "siphonLife",
				location: {
					rowIdx: 4,
					colIdx: 1
				},
				spellIds: [
					63108
				],
				maxPoints: 1
			},
			{
				fieldName: "curseOfExhaustion",
				location: {
					rowIdx: 4,
					colIdx: 2
				},
				spellIds: [
					18223
				],
				maxPoints: 1,
				"prereqLocation": {
					rowIdx: 2,
					colIdx: 2
				}
			},
			{
				fieldName: "improvedFelhunter",
				location: {
					rowIdx: 5,
					colIdx: 0
				},
				spellIds: [
					54037,
					54038
				],
				maxPoints: 2
			},
			{
				fieldName: "shadowMastery",
				location: {
					rowIdx: 5,
					colIdx: 1
				},
				spellIds: [
					18271,
					18272,
					18273,
					18274,
					18275
				],
				maxPoints: 5,
				"prereqLocation": {
					rowIdx: 4,
					colIdx: 1
				}
			},
			{
				fieldName: "eradication",
				location: {
					rowIdx: 6,
					colIdx: 0
				},
				spellIds: [
					47195,
					47196,
					47197
				],
				maxPoints: 3
			},
			{
				fieldName: "contagion",
				location: {
					rowIdx: 6,
					colIdx: 1
				},
				spellIds: [
					30060,
					30061,
					30062,
					30063,
					30064
				],
				maxPoints: 5
			},
			{
				fieldName: "darkPact",
				location: {
					rowIdx: 6,
					colIdx: 2
				},
				spellIds: [
					18220
				],
				maxPoints: 1
			},
			{
				fieldName: "improvedHowlOfTerror",
				location: {
					rowIdx: 7,
					colIdx: 0
				},
				spellIds: [
					30054,
					30057
				],
				maxPoints: 2
			},
			{
				fieldName: "malediction",
				location: {
					rowIdx: 7,
					colIdx: 2
				},
				spellIds: [
					32477,
					32483,
					32484
				],
				maxPoints: 3
			},
			{
				fieldName: "deathsEmbrace",
				location: {
					rowIdx: 8,
					colIdx: 0
				},
				spellIds: [
					47198,
					47199,
					47200
				],
				maxPoints: 3
			},
			{
				fieldName: "unstableAffliction",
				location: {
					rowIdx: 8,
					colIdx: 1
				},
				spellIds: [
					30108
				],
				maxPoints: 1,
				"prereqLocation": {
					rowIdx: 6,
					colIdx: 1
				}
			},
			{
				fieldName: "pandemic",
				location: {
					rowIdx: 8,
					colIdx: 2
				},
				spellIds: [
					58435
				],
				maxPoints: 1,
				"prereqLocation": {
					rowIdx: 8,
					colIdx: 1
				}
			},
			{
				fieldName: "everlastingAffliction",
				location: {
					rowIdx: 9,
					colIdx: 1
				},
				spellIds: [
					47201,
					47202,
					47203,
					47204,
					47205
				],
				maxPoints: 5
			},
			{
				fieldName: "haunt",
				location: {
					rowIdx: 10,
					colIdx: 1
				},
				spellIds: [
					48181
				],
				maxPoints: 1
			}
		]
	},
	{
		name: "Demonology",
		backgroundUrl: "https://wow.zamimg.com/images/wow/talents/backgrounds/wrath/303.jpg",
		talents: [
			{
				fieldName: "improvedHealthstone",
				location: {
					rowIdx: 0,
					colIdx: 0
				},
				spellIds: [
					18692,
					18693
				],
				maxPoints: 2
			},
			{
				fieldName: "improvedImp",
				location: {
					rowIdx: 0,
					colIdx: 1
				},
				spellIds: [
					18694,
					18695,
					18696
				],
				maxPoints: 3
			},
			{
				fieldName: "demonicEmbrace",
				location: {
					rowIdx: 0,
					colIdx: 2
				},
				spellIds: [
					18697,
					18698,
					18699
				],
				maxPoints: 3
			},
			{
				fieldName: "felSynergy",
				location: {
					rowIdx: 0,
					colIdx: 3
				},
				spellIds: [
					47230,
					47231
				],
				maxPoints: 2
			},
			{
				fieldName: "improvedHealthFunnel",
				location: {
					rowIdx: 1,
					colIdx: 0
				},
				spellIds: [
					18703,
					18704
				],
				maxPoints: 2
			},
			{
				fieldName: "demonicBrutality",
				location: {
					rowIdx: 1,
					colIdx: 1
				},
				spellIds: [
					18705,
					18706,
					18707
				],
				maxPoints: 3
			},
			{
				fieldName: "felVitality",
				location: {
					rowIdx: 1,
					colIdx: 2
				},
				spellIds: [
					18731,
					18743,
					18744
				],
				maxPoints: 3
			},
			{
				fieldName: "improvedSayaad",
				location: {
					rowIdx: 2,
					colIdx: 0
				},
				spellIds: [
					18754,
					18755,
					18756
				],
				maxPoints: 3
			},
			{
				fieldName: "soulLink",
				location: {
					rowIdx: 2,
					colIdx: 1
				},
				spellIds: [
					19028
				],
				maxPoints: 1
			},
			{
				fieldName: "felDomination",
				location: {
					rowIdx: 2,
					colIdx: 2
				},
				spellIds: [
					18708
				],
				maxPoints: 1
			},
			{
				fieldName: "demonicAegis",
				location: {
					rowIdx: 2,
					colIdx: 3
				},
				spellIds: [
					30143,
					30144,
					30145
				],
				maxPoints: 3
			},
			{
				fieldName: "unholyPower",
				location: {
					rowIdx: 3,
					colIdx: 1
				},
				spellIds: [
					18769,
					18770,
					18771,
					18772,
					18773
				],
				maxPoints: 5,
				"prereqLocation": {
					rowIdx: 2,
					colIdx: 1
				}
			},
			{
				fieldName: "masterSummoner",
				location: {
					rowIdx: 3,
					colIdx: 2
				},
				spellIds: [
					18709,
					18710
				],
				maxPoints: 2,
				"prereqLocation": {
					rowIdx: 2,
					colIdx: 2
				}
			},
			{
				fieldName: "manaFeed",
				location: {
					rowIdx: 4,
					colIdx: 0
				},
				spellIds: [
					30326
				],
				maxPoints: 1,
				"prereqLocation": {
					rowIdx: 3,
					colIdx: 1
				}
			},
			{
				fieldName: "masterConjuror",
				location: {
					rowIdx: 4,
					colIdx: 2
				},
				spellIds: [
					18767,
					18768
				],
				maxPoints: 2
			},
			{
				fieldName: "masterDemonologist",
				location: {
					rowIdx: 5,
					colIdx: 1
				},
				spellIds: [
					23785,
					23822,
					23823,
					23824,
					23825
				],
				maxPoints: 5,
				"prereqLocation": {
					rowIdx: 3,
					colIdx: 1
				}
			},
			{
				fieldName: "moltenCore",
				location: {
					rowIdx: 5,
					colIdx: 2
				},
				spellIds: [
					47245,
					47246,
					47247
				],
				maxPoints: 3
			},
			{
				fieldName: "demonicResilience",
				location: {
					rowIdx: 6,
					colIdx: 0
				},
				spellIds: [
					30319,
					30320,
					30321
				],
				maxPoints: 3
			},
			{
				fieldName: "demonicEmpowerment",
				location: {
					rowIdx: 6,
					colIdx: 1
				},
				spellIds: [
					47193
				],
				maxPoints: 1,
				"prereqLocation": {
					rowIdx: 5,
					colIdx: 1
				}
			},
			{
				fieldName: "demonicKnowledge",
				location: {
					rowIdx: 6,
					colIdx: 2
				},
				spellIds: [
					35691,
					35692,
					35693
				],
				maxPoints: 3
			},
			{
				fieldName: "demonicTactics",
				location: {
					rowIdx: 7,
					colIdx: 1
				},
				spellIds: [
					30242,
					30245,
					30246,
					30247,
					30248
				],
				maxPoints: 5
			},
			{
				fieldName: "decimation",
				location: {
					rowIdx: 7,
					colIdx: 2
				},
				spellIds: [
					63156,
					63158
				],
				maxPoints: 2
			},
			{
				fieldName: "improvedDemonicTactics",
				location: {
					rowIdx: 8,
					colIdx: 0
				},
				spellIds: [
					54347,
					54348,
					54349
				],
				maxPoints: 3,
				"prereqLocation": {
					rowIdx: 7,
					colIdx: 1
				}
			},
			{
				fieldName: "summonFelguard",
				location: {
					rowIdx: 8,
					colIdx: 1
				},
				spellIds: [
					30146
				],
				maxPoints: 1
			},
			{
				fieldName: "nemesis",
				location: {
					rowIdx: 8,
					colIdx: 2
				},
				spellIds: [
					63117,
					63121,
					63123
				],
				maxPoints: 3
			},
			{
				fieldName: "demonicPact",
				location: {
					rowIdx: 9,
					colIdx: 1
				},
				spellIds: [
					47236,
					47237,
					47238,
					47239,
					47240
				],
				maxPoints: 5
			},
			{
				fieldName: "metamorphosis",
				location: {
					rowIdx: 10,
					colIdx: 1
				},
				spellIds: [
					59672
				],
				maxPoints: 1
			}
		]
	},
	{
		name: "Destruction",
		backgroundUrl: "https://wow.zamimg.com/images/wow/talents/backgrounds/wrath/301.jpg",
		talents: [
			{
				fieldName: "improvedShadowBolt",
				location: {
					rowIdx: 0,
					colIdx: 1
				},
				spellIds: [
					17793,
					17796,
					17801,
					17802,
					17803
				],
				maxPoints: 5
			},
			{
				fieldName: "bane",
				location: {
					rowIdx: 0,
					colIdx: 2
				},
				spellIds: [
					17788,
					17789,
					17790,
					17791,
					17792
				],
				maxPoints: 5
			},
			{
				fieldName: "aftermath",
				location: {
					rowIdx: 1,
					colIdx: 0
				},
				spellIds: [
					18119,
					18120
				],
				maxPoints: 2
			},
			{
				fieldName: "moltenSkin",
				location: {
					rowIdx: 1,
					colIdx: 1
				},
				spellIds: [
					63349,
					63350,
					63351
				],
				maxPoints: 3
			},
			{
				fieldName: "cataclysm",
				location: {
					rowIdx: 1,
					colIdx: 2
				},
				spellIds: [
					17778,
					17779,
					17780
				],
				maxPoints: 3
			},
			{
				fieldName: "demonicPower",
				location: {
					rowIdx: 2,
					colIdx: 0
				},
				spellIds: [
					18126,
					18127
				],
				maxPoints: 2
			},
			{
				fieldName: "shadowburn",
				location: {
					rowIdx: 2,
					colIdx: 1
				},
				spellIds: [
					17877
				],
				maxPoints: 1
			},
			{
				fieldName: "ruin",
				location: {
					rowIdx: 2,
					colIdx: 2
				},
				spellIds: [
					17959,
					59738,
					59739,
					59740,
					59741
				],
				maxPoints: 5
			},
			{
				fieldName: "intensity",
				location: {
					rowIdx: 3,
					colIdx: 0
				},
				spellIds: [
					18135,
					18136
				],
				maxPoints: 2
			},
			{
				fieldName: "destructiveReach",
				location: {
					rowIdx: 3,
					colIdx: 1
				},
				spellIds: [
					17917,
					17918
				],
				maxPoints: 2
			},
			{
				fieldName: "improvedSearingPain",
				location: {
					rowIdx: 3,
					colIdx: 3
				},
				spellIds: [
					17927,
					17929,
					17930
				],
				maxPoints: 3
			},
			{
				fieldName: "backlash",
				location: {
					rowIdx: 4,
					colIdx: 0
				},
				spellIds: [
					34935,
					34938,
					34939
				],
				maxPoints: 3,
				"prereqLocation": {
					rowIdx: 3,
					colIdx: 0
				}
			},
			{
				fieldName: "improvedImmolate",
				location: {
					rowIdx: 4,
					colIdx: 1
				},
				spellIds: [
					17815,
					17833,
					17834
				],
				maxPoints: 3
			},
			{
				fieldName: "devastation",
				location: {
					rowIdx: 4,
					colIdx: 2
				},
				spellIds: [
					18130
				],
				maxPoints: 1,
				"prereqLocation": {
					rowIdx: 2,
					colIdx: 2
				}
			},
			{
				fieldName: "netherProtection",
				location: {
					rowIdx: 5,
					colIdx: 0
				},
				spellIds: [
					30299,
					30301,
					30302
				],
				maxPoints: 3
			},
			{
				fieldName: "emberstorm",
				location: {
					rowIdx: 5,
					colIdx: 2
				},
				spellIds: [
					17954,
					17955,
					17956,
					17957,
					17958
				],
				maxPoints: 5
			},
			{
				fieldName: "conflagrate",
				location: {
					rowIdx: 6,
					colIdx: 1
				},
				spellIds: [
					17962
				],
				maxPoints: 1,
				"prereqLocation": {
					rowIdx: 4,
					colIdx: 1
				}
			},
			{
				fieldName: "soulLeech",
				location: {
					rowIdx: 6,
					colIdx: 2
				},
				spellIds: [
					30293,
					30295,
					30296
				],
				maxPoints: 3
			},
			{
				fieldName: "pyroclasm",
				location: {
					rowIdx: 6,
					colIdx: 3
				},
				spellIds: [
					18096,
					18073,
					63245
				],
				maxPoints: 3
			},
			{
				fieldName: "shadowAndFlame",
				location: {
					rowIdx: 7,
					colIdx: 1
				},
				spellIds: [
					30288,
					30289,
					30290,
					30291,
					30292
				],
				maxPoints: 5
			},
			{
				fieldName: "improvedSoulLeech",
				location: {
					rowIdx: 7,
					colIdx: 2
				},
				spellIds: [
					54117,
					54118
				],
				maxPoints: 2,
				"prereqLocation": {
					rowIdx: 6,
					colIdx: 2
				}
			},
			{
				fieldName: "backdraft",
				location: {
					rowIdx: 8,
					colIdx: 0
				},
				spellIds: [
					47258,
					47259,
					47260
				],
				maxPoints: 3,
				"prereqLocation": {
					rowIdx: 6,
					colIdx: 1
				}
			},
			{
				fieldName: "shadowfury",
				location: {
					rowIdx: 8,
					colIdx: 1
				},
				spellIds: [
					30283
				],
				maxPoints: 1
			},
			{
				fieldName: "empoweredImp",
				location: {
					rowIdx: 8,
					colIdx: 2
				},
				spellIds: [
					47220,
					47221,
					47223
				],
				maxPoints: 3
			},
			{
				fieldName: "fireAndBrimstone",
				location: {
					rowIdx: 9,
					colIdx: 1
				},
				spellIds: [
					47266,
					47267,
					47268,
					47269,
					47270
				],
				maxPoints: 5
			},
			{
				fieldName: "chaosBolt",
				location: {
					rowIdx: 10,
					colIdx: 1
				},
				spellIds: [
					50796
				],
				maxPoints: 1
			}
		]
	}
]);

export const warlockGlyphsConfig: GlyphsConfig = {
	majorGlyphs: {
		[WarlockMajorGlyph.GlyphOfChaosBolt]: {
			name: 'Glyph of Chaos Bolt',
			description: 'Reduces the cooldown on Chaos Bolt by 2 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_warlock_chaosbolt.jpg',
		},
		[WarlockMajorGlyph.GlyphOfConflagrate]: {
			name: 'Glyph of Conflagrate',
			description: 'Your Conflagrate spell no longer consumes your Immolate or Shadowflame spell from the target.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_fire_fireball.jpg',
		},
		[WarlockMajorGlyph.GlyphOfCorruption]: {
			name: 'Glyph of Corruption',
			description: 'Your Corruption spell has a 4% chance to cause you to enter a Shadow Trance state after damaging the opponent. The Shadow Trance state reduces the casting time of your next Shadow Bolt spell by 100%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_abominationexplosion.jpg',
		},
		[WarlockMajorGlyph.GlyphOfCurseOfAgony]: {
			name: 'Glyph of Curse of Agony',
			description: 'Increases the duration of your Curse of Agony by 4 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_curseofsargeras.jpg',
		},
		[WarlockMajorGlyph.GlyphOfDeathCoil]: {
			name: 'Glyph of Death Coil',
			description: 'Increases the duration of your Death Coil by 0.5 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_deathcoil.jpg',
		},
		[WarlockMajorGlyph.GlyphOfDemonicCircle]: {
			name: 'Glyph of Demonic Circle',
			description: 'Reduces the cooldown on Demonic Circle by 4 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_demoniccirclesummon.jpg',
		},
		[WarlockMajorGlyph.GlyphOfFear]: {
			name: 'Glyph of Fear',
			description: 'Increases the damage your Fear target can take before the Fear effect is removed by 20%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_possession.jpg',
		},
		[WarlockMajorGlyph.GlyphOfFelguard]: {
			name: 'Glyph of Felguard',
			description: 'Increases the Felguard\'s total attack power by 20%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_summonfelguard.jpg',
		},
		[WarlockMajorGlyph.GlyphOfFelhunter]: {
			name: 'Glyph of Felhunter',
			description: 'When your Felhunter uses Devour Magic, you will also be healed for that amount.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_summonfelhunter.jpg',
		},
		[WarlockMajorGlyph.GlyphOfHaunt]: {
			name: 'Glyph of Haunt',
			description: 'The bonus damage granted by your Haunt spell is increased by an additional 3%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_warlock_haunt.jpg',
		},
		[WarlockMajorGlyph.GlyphOfHealthFunnel]: {
			name: 'Glyph of Health Funnel',
			description: 'Reduces the pushback suffered from damaging attacks while channeling your Health Funnel spell by 100%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_lifedrain.jpg',
		},
		[WarlockMajorGlyph.GlyphOfHealthstone]: {
			name: 'Glyph of Healthstone',
			description: 'You receive 30% more healing from using a healthstone.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/inv_stone_04.jpg',
		},
		[WarlockMajorGlyph.GlyphOfHowlOfTerror]: {
			name: 'Glyph of Howl of Terror',
			description: 'Reduces the cooldown on your Howl of Terror spell by 8 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_deathscream.jpg',
		},
		[WarlockMajorGlyph.GlyphOfImmolate]: {
			name: 'Glyph of Immolate',
			description: 'Increases the periodic damage of your Immolate by 10%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_fire_immolation.jpg',
		},
		[WarlockMajorGlyph.GlyphOfImp]: {
			name: 'Glyph of Imp',
			description: 'Increases the damage done by your Imp\'s Firebolt spell by 20%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_summonimp.jpg',
		},
		[WarlockMajorGlyph.GlyphOfIncinerate]: {
			name: 'Glyph of Incinerate',
			description: 'Increases the damage done by Incinerate by 5%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_fire_burnout.jpg',
		},
		[WarlockMajorGlyph.GlyphOfLifeTap]: {
			name: 'Glyph of Life Tap',
			description: 'When you use Life Tap or Dark Pact, you gain 20% of your Spirit as spell power for 40 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_burningspirit.jpg',
		},
		[WarlockMajorGlyph.GlyphOfMetamorphosis]: {
			name: 'Glyph of Metamorphosis',
			description: 'Increases the duration of your Metamorphosis by 6 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_demonform.jpg',
		},
		[WarlockMajorGlyph.GlyphOfQuickDecay]: {
			name: 'Glyph of Quick Decay',
			description: 'Your haste now reduces the time between periodic damage ticks of your Corruption spell.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_abominationexplosion.jpg',
		},
		[WarlockMajorGlyph.GlyphOfSearingPain]: {
			name: 'Glyph of Searing Pain',
			description: 'Increases the critical strike bonus of your Searing Pain by 20%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_fire_soulburn.jpg',
		},
		[WarlockMajorGlyph.GlyphOfShadowBolt]: {
			name: 'Glyph of Shadow Bolt',
			description: 'Reduces the mana cost of your Shadow Bolt by 10%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_shadowbolt.jpg',
		},
		[WarlockMajorGlyph.GlyphOfShadowburn]: {
			name: 'Glyph of Shadowburn',
			description: 'Increases the critical strike chance of Shadowburn by 20% when the target is below 35% health.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_scourgebuild.jpg',
		},
		[WarlockMajorGlyph.GlyphOfShadowflame]: {
			name: 'Glyph of Shadowflame',
			description: 'Your Shadowflame also applies a 70% movement speed slow on its victims.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_warlock_shadowflame.jpg',
		},
		[WarlockMajorGlyph.GlyphOfSiphonLife]: {
			name: 'Glyph of Siphon Life',
			description: 'Increases the healing you receive from your Siphon Life talent by 25%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_requiem.jpg',
		},
		[WarlockMajorGlyph.GlyphOfSoulLink]: {
			name: 'Glyph of Soul Link',
			description: 'Increases the percentage of damage shared via your Soul Link by an additional 5%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_gathershadows.jpg',
		},
		[WarlockMajorGlyph.GlyphOfSoulstone]: {
			name: 'Glyph of Soulstone',
			description: 'Increases the amount of health you gain from resurrecting via a Soulstone by 300%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_soulgem.jpg',
		},
		[WarlockMajorGlyph.GlyphOfSuccubus]: {
			name: 'Glyph of Succubus',
			description: 'Your Succubus\'s Seduction ability also removes all damage over time effects from the target.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_summonsuccubus.jpg',
		},
		[WarlockMajorGlyph.GlyphOfUnstableAffliction]: {
			name: 'Glyph of Unstable Affliction',
			description: 'Decreases the casting time of your Unstable Affliction by 0.2 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_unstableaffliction_3.jpg',
		},
		[WarlockMajorGlyph.GlyphOfVoidwalker]: {
			name: 'Glyph of Voidwalker',
			description: 'Increases your Voidwalker\'s total Stamina by 20%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_summonvoidwalker.jpg',
		},
	},
	minorGlyphs: {
		[WarlockMinorGlyph.GlyphOfCurseOfExhausion]: {
			name: 'Glyph of Curse of Exhausion',
			description: 'Increases the range of your Curse of Exhaustion spell by 5 yards.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_grimward.jpg',
		},
		[WarlockMinorGlyph.GlyphOfDrainSoul]: {
			name: 'Glyph of Drain Soul',
			description: 'Your Drain Soul ability occasionally creates an additional soul shard.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_haunting.jpg',
		},
		[WarlockMinorGlyph.GlyphOfEnslaveDemon]: {
			name: 'Glyph of Enslave Demon',
			description: 'Reduces the cast time of your Enslave Demon spell by 50%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_enslavedemon.jpg',
		},
		[WarlockMinorGlyph.GlyphOfKilrogg]: {
			name: 'Glyph of Kilrogg',
			description: 'Increases the movement speed of your Eye of Kilrogg by 50% and allows it to fly in areas where flying mounts are enabled.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_evileye.jpg',
		},
		[WarlockMinorGlyph.GlyphOfSouls]: {
			name: 'Glyph of Souls',
			description: 'Reduces the mana cost of your Ritual of Souls spell by 70%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_shadesofdarkness.jpg',
		},
		[WarlockMinorGlyph.GlyphOfUnendingBreath]: {
			name: 'Glyph of Unending Breath',
			description: 'Increases the swim speed of targets affected by your Unending Breath spell by 20%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_demonbreath.jpg',
		},
	},
};
