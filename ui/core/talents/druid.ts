import { DruidTalents, DruidMajorGlyph, DruidMinorGlyph } from '/wotlk/core/proto/druid.js';

import { GlyphsConfig, } from './glyphs_picker.js';
import { TalentsConfig, newTalentsConfig } from './talents_picker.js';

export const druidTalentsConfig: TalentsConfig<DruidTalents> = newTalentsConfig([
	{
		name: 'Balance',
		backgroundUrl: 'https://wow.zamimg.com/images/wow/talents/backgrounds/wrath/283.jpg',
		talents: [
			{
				fieldName: 'starlightWrath',
				location: {
					rowIdx: 0,
					colIdx: 0,
				},
				spellIds: [16814],
				maxPoints: 5,
			},
			{
				//fieldName: 'naturesGrasp',
				location: {
					rowIdx: 0,
					colIdx: 1,
				},
				spellIds: [16689],
				maxPoints: 1,
			},
			{
				//fieldName: 'improvedNaturesGrasp',
				location: {
					rowIdx: 0,
					colIdx: 2,
				},
				prereqLocation: {
					rowIdx: 0,
					colIdx: 1,
				},
				spellIds: [17245, 17247],
				maxPoints: 4,
			},
			{
				//fieldName: 'controlOfNature',
				location: {
					rowIdx: 1,
					colIdx: 0,
				},
				spellIds: [16918],
				maxPoints: 3,
			},
			{
				fieldName: 'focusedStarlight',
				location: {
					rowIdx: 1,
					colIdx: 1,
				},
				spellIds: [35363],
				maxPoints: 2,
			},
			{
				fieldName: 'improvedMoonfire',
				location: {
					rowIdx: 1,
					colIdx: 2,
				},
				spellIds: [16821],
				maxPoints: 2,
			},
			{
				fieldName: 'brambles',
				location: {
					rowIdx: 2,
					colIdx: 0,
				},
				spellIds: [16836, 16839],
				maxPoints: 3,
			},
			{
				fieldName: 'insectSwarm',
				location: {
					rowIdx: 2,
					colIdx: 2,
				},
				spellIds: [5570],
				maxPoints: 1,
			},
			{
				//fieldName: 'naturesReach',
				location: {
					rowIdx: 2,
					colIdx: 3,
				},
				spellIds: [16819],
				maxPoints: 2,
			},
			{
				fieldName: 'vengeance',
				location: {
					rowIdx: 3,
					colIdx: 1,
				},
				prereqLocation: {
					rowIdx: 1,
					colIdx: 1,
				},
				spellIds: [16909],
				maxPoints: 5,
			},
			{
				//fieldName: 'celestialFocus',
				location: {
					rowIdx: 3,
					colIdx: 2,
				},
				spellIds: [16850, 16923],
				maxPoints: 3,
			},
			{
				fieldName: 'lunarGuidance',
				location: {
					rowIdx: 4,
					colIdx: 0,
				},
				spellIds: [33589],
				maxPoints: 3,
			},
			{
				fieldName: 'naturesGrace',
				location: {
					rowIdx: 4,
					colIdx: 1,
				},
				spellIds: [16880],
				maxPoints: 1,
			},
			{
				fieldName: 'moonglow',
				location: {
					rowIdx: 4,
					colIdx: 2,
				},
				spellIds: [16845],
				maxPoints: 3,
			},
			{
				fieldName: 'moonfury',
				location: {
					rowIdx: 5,
					colIdx: 1,
				},
				prereqLocation: {
					rowIdx: 4,
					colIdx: 1,
				},
				spellIds: [16896, 16897, 16899],
				maxPoints: 5,
			},
			{
				fieldName: 'balanceOfPower',
				location: {
					rowIdx: 5,
					colIdx: 2,
				},
				spellIds: [33592, 33596],
				maxPoints: 2,
			},
			{
				fieldName: 'dreamstate',
				location: {
					rowIdx: 6,
					colIdx: 0,
				},
				spellIds: [33597, 33599, 33956],
				maxPoints: 3,
			},
			{
				fieldName: 'moonkinForm',
				location: {
					rowIdx: 6,
					colIdx: 1,
				},
				spellIds: [24858],
				maxPoints: 1,
			},
			{
				fieldName: 'improvedFaerieFire',
				location: {
					rowIdx: 6,
					colIdx: 2,
				},
				spellIds: [33600],
				maxPoints: 3,
			},
			{
				fieldName: 'wrathOfCenarius',
				location: {
					rowIdx: 7,
					colIdx: 1,
				},
				spellIds: [33603],
				maxPoints: 5,
			},
			{
				fieldName: 'forceOfNature',
				location: {
					rowIdx: 8,
					colIdx: 1,
				},
				spellIds: [33831],
				maxPoints: 1,
			},
		],
	},
	{
		name: 'Feral Combat',
		backgroundUrl: 'https://wow.zamimg.com/images/wow/talents/backgrounds/wrath/281.jpg',
		talents: [
			{
				fieldName: 'ferocity',
				location: {
					rowIdx: 0,
					colIdx: 1,
				},
				spellIds: [16934],
				maxPoints: 5,
			},
			{
				fieldName: 'feralAggression',
				location: {
					rowIdx: 0,
					colIdx: 2,
				},
				spellIds: [16858],
				maxPoints: 5,
			},
			{
				fieldName: 'feralInstinct',
				location: {
					rowIdx: 1,
					colIdx: 0,
				},
				spellIds: [16947],
				maxPoints: 3,
			},
			{
				//fieldName: 'brutalImpact',
				location: {
					rowIdx: 1,
					colIdx: 1,
				},
				spellIds: [16940],
				maxPoints: 2,
			},
			{
				fieldName: 'thickHide',
				location: {
					rowIdx: 1,
					colIdx: 2,
				},
				spellIds: [16929],
				maxPoints: 3,
			},
			{
				fieldName: 'feralSwiftness',
				location: {
					rowIdx: 2,
					colIdx: 0,
				},
				spellIds: [17002, 24866],
				maxPoints: 2,
			},
			{
				//fieldName: 'feralCharge',
				location: {
					rowIdx: 2,
					colIdx: 1,
				},
				spellIds: [16979],
				maxPoints: 1,
			},
			{
				fieldName: 'sharpenedClaws',
				location: {
					rowIdx: 2,
					colIdx: 2,
				},
				spellIds: [16942],
				maxPoints: 3,
			},
			{
				fieldName: 'shreddingAttacks',
				location: {
					rowIdx: 3,
					colIdx: 0,
				},
				spellIds: [16966, 16968],
				maxPoints: 2,
			},
			{
				fieldName: 'predatoryStrikes',
				location: {
					rowIdx: 3,
					colIdx: 1,
				},
				spellIds: [16972, 16974],
				maxPoints: 3,
			},
			{
				fieldName: 'primalFury',
				location: {
					rowIdx: 3,
					colIdx: 2,
				},
				spellIds: [37116],
				maxPoints: 2,
			},
			{
				fieldName: 'savageFury',
				location: {
					rowIdx: 4,
					colIdx: 0,
				},
				spellIds: [16998],
				maxPoints: 2,
			},
			{
				fieldName: 'faerieFire',
				location: {
					rowIdx: 4,
					colIdx: 2,
				},
				spellIds: [16857],
				maxPoints: 1,
			},
			{
				//fieldName: 'nurturingInstinct',
				location: {
					rowIdx: 4,
					colIdx: 3,
				},
				spellIds: [33872],
				maxPoints: 2,
			},
			{
				fieldName: 'heartOfTheWild',
				location: {
					rowIdx: 5,
					colIdx: 1,
				},
				prereqLocation: {
					rowIdx: 3,
					colIdx: 1,
				},
				spellIds: [17003, 17004, 17005, 17006, 24894],
				maxPoints: 5,
			},
			{
				fieldName: 'survivalOfTheFittest',
				location: {
					rowIdx: 5,
					colIdx: 2,
				},
				spellIds: [33853, 33855],
				maxPoints: 3,
			},
			{
				//fieldName: 'primalTenacity',
				location: {
					rowIdx: 6,
					colIdx: 0,
				},
				spellIds: [33851, 33852, 33957],
				maxPoints: 3,
			},
			{
				fieldName: 'leaderOfThePack',
				location: {
					rowIdx: 6,
					colIdx: 1,
				},
				spellIds: [17007],
				maxPoints: 1,
			},
			{
				fieldName: 'improvedLeaderOfThePack',
				location: {
					rowIdx: 6,
					colIdx: 2,
				},
				prereqLocation: {
					rowIdx: 6,
					colIdx: 1,
				},
				spellIds: [34297, 34300],
				maxPoints: 2,
			},
			{
				fieldName: 'predatoryInstincts',
				location: {
					rowIdx: 7,
					colIdx: 2,
				},
				spellIds: [33859, 33866],
				maxPoints: 5,
			},
			{
				fieldName: 'mangle',
				location: {
					rowIdx: 8,
					colIdx: 1,
				},
				prereqLocation: {
					rowIdx: 6,
					colIdx: 1,
				},
				spellIds: [33917],
				maxPoints: 1,
			},
		],
	},
	{
		name: 'Restoration',
		backgroundUrl: 'https://wow.zamimg.com/images/wow/talents/backgrounds/wrath/282.jpg',
		talents: [
			{
				fieldName: 'improvedMarkOfTheWild',
				location: {
					rowIdx: 0,
					colIdx: 1,
				},
				spellIds: [17050, 17051, 17053],
				maxPoints: 5,
			},
			{
				fieldName: 'furor',
				location: {
					rowIdx: 0,
					colIdx: 2,
				},
				spellIds: [17056, 17058],
				maxPoints: 5,
			},
			{
				fieldName: 'naturalist',
				location: {
					rowIdx: 1,
					colIdx: 0,
				},
				spellIds: [17069],
				maxPoints: 5,
			},
			{
				//fieldName: 'naturesFocus',
				location: {
					rowIdx: 1,
					colIdx: 1,
				},
				spellIds: [17063, 17065],
				maxPoints: 5,
			},
			{
				fieldName: 'naturalShapeshifter',
				location: {
					rowIdx: 1,
					colIdx: 2,
				},
				spellIds: [16833],
				maxPoints: 3,
			},
			{
				fieldName: 'intensity',
				location: {
					rowIdx: 2,
					colIdx: 0,
				},
				spellIds: [17106],
				maxPoints: 3,
			},
			{
				fieldName: 'subtlety',
				location: {
					rowIdx: 2,
					colIdx: 1,
				},
				spellIds: [17118],
				maxPoints: 5,
			},
			{
				fieldName: 'omenOfClarity',
				location: {
					rowIdx: 2,
					colIdx: 2,
				},
				spellIds: [16864],
				maxPoints: 1,
			},
			{
				//fieldName: 'tranquilSpirit',
				location: {
					rowIdx: 3,
					colIdx: 1,
				},
				spellIds: [24968],
				maxPoints: 5,
			},
			{
				//fieldName: 'improvedRejuvenation',
				location: {
					rowIdx: 3,
					colIdx: 2,
				},
				spellIds: [17111],
				maxPoints: 3,
			},
			{
				fieldName: 'naturesSwiftness',
				location: {
					rowIdx: 4,
					colIdx: 0,
				},
				prereqLocation: {
					rowIdx: 2,
					colIdx: 0,
				},
				spellIds: [17116],
				maxPoints: 1,
			},
			{
				//fieldName: 'giftOfNature',
				location: {
					rowIdx: 4,
					colIdx: 1,
				},
				spellIds: [17104, 24943],
				maxPoints: 5,
			},
			{
				//fieldName: 'improvedTranquility',
				location: {
					rowIdx: 4,
					colIdx: 3,
				},
				spellIds: [17123],
				maxPoints: 2,
			},
			{
				//fieldName: 'empoweredTouch',
				location: {
					rowIdx: 5,
					colIdx: 0,
				},
				spellIds: [33879],
				maxPoints: 2,
			},
			{
				//fieldName: 'improvedRegrowth',
				location: {
					rowIdx: 5,
					colIdx: 2,
				},
				prereqLocation: {
					rowIdx: 3,
					colIdx: 2,
				},
				spellIds: [17074],
				maxPoints: 5,
			},
			{
				fieldName: 'livingSpirit',
				location: {
					rowIdx: 6,
					colIdx: 0,
				},
				spellIds: [34151],
				maxPoints: 3,
			},
			{
				//fieldName: 'swiftmend',
				location: {
					rowIdx: 6,
					colIdx: 1,
				},
				prereqLocation: {
					rowIdx: 4,
					colIdx: 1,
				},
				spellIds: [18562],
				maxPoints: 1,
			},
			{
				fieldName: 'naturalPerfection',
				location: {
					rowIdx: 6,
					colIdx: 2,
				},
				spellIds: [33881],
				maxPoints: 3,
			},
			{
				//fieldName: 'empoweredRejuvenation',
				location: {
					rowIdx: 7,
					colIdx: 1,
				},
				spellIds: [33886],
				maxPoints: 5,
			},
			{
				//fieldName: 'treeOfLife',
				location: {
					rowIdx: 8,
					colIdx: 1,
				},
				spellIds: [33891],
				maxPoints: 1,
			},
		],
	},
]);

export const druidGlyphsConfig: GlyphsConfig = {
	majorGlyphs: {
		[DruidMajorGlyph.GlyphOfBarkskin]: {
			name: 'Glyph of Barkskin',
			description: 'Reduces the chance you\'ll be critically hit by 25% while Barkskin is active.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_nature_stoneclawtotem.jpg',
		},
		[DruidMajorGlyph.GlyphOfBerserk]: {
			name: 'Glyph of Berserk',
			description: 'Increases the duration of Berserk by 5 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_druid_berserk.jpg',
		},
		[DruidMajorGlyph.GlyphOfClaw]: {
			name: 'Glyph of Claw',
			description: 'Reduces the energy cost of your Claw ability by 5.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_druid_rake.jpg',
		},
		[DruidMajorGlyph.GlyphOfEntanglingRoots]: {
			name: 'Glyph of Entangling Roots',
			description: 'Increases the damage your Entangling Roots victims can take before the Entangling Roots automatically breaks by 20%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_nature_stranglevines.jpg',
		},
		[DruidMajorGlyph.GlyphOfFocus]: {
			name: 'Glyph of Focus',
			description: 'Increases the damage done by Starfall by 10%, but decreases its radius by 50%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_druid_starfall.jpg',
		},
		[DruidMajorGlyph.GlyphOfFrenziedRegeneration]: {
			name: 'Glyph of Frenzied Regeneration',
			description: 'While Frenzied Regeneration is active, healing effects on you are 20% more powerful.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_bullrush.jpg',
		},
		[DruidMajorGlyph.GlyphOfGrowling]: {
			name: 'Glyph of Growling',
			description: 'Increases the chance for your Growl ability to work successfully by 8%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_physical_taunt.jpg',
		},
		[DruidMajorGlyph.GlyphOfHealingTouch]: {
			name: 'Glyph of Healing Touch',
			description: 'Decreases the cast time of Healing Touch by 1.5 sec, the mana cost by 25%, and the amount healed by 50%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_nature_healingtouch.jpg',
		},
		[DruidMajorGlyph.GlyphOfHurricane]: {
			name: 'Glyph of Hurricane',
			description: 'Your Hurricane ability now also slows the movement speed of its victims by 20%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_nature_cyclone.jpg',
		},
		[DruidMajorGlyph.GlyphOfInnervate]: {
			name: 'Glyph of Innervate',
			description: 'Innervate now grants the caster 45% of <dfn title="her">his</dfn> base mana pool over 10 sec in addition to the normal effects of Innervate.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_nature_lightning.jpg',
		},
		[DruidMajorGlyph.GlyphOfInsectSwarm]: {
			name: 'Glyph of Insect Swarm',
			description: 'Increases the damage of your Insect Swarm ability by 30%, but it no longer affects your victim\'s chance to hit.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_nature_insectswarm.jpg',
		},
		[DruidMajorGlyph.GlyphOfLifebloom]: {
			name: 'Glyph of Lifebloom',
			description: 'Increases the duration of Lifebloom by 1 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/inv_misc_herb_felblossom.jpg',
		},
		[DruidMajorGlyph.GlyphOfMangle]: {
			name: 'Glyph of Mangle',
			description: 'Increases the damage done by Mangle by 10%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_druid_mangle2.jpg',
		},
		[DruidMajorGlyph.GlyphOfMaul]: {
			name: 'Glyph of Maul',
			description: 'Your Maul ability now hits 1 additional target.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_druid_maul.jpg',
		},
		[DruidMajorGlyph.GlyphOfMonsoon]: {
			name: 'Glyph of Monsoon',
			description: 'Reduces the cooldown of your Typhoon spell by 3 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_druid_typhoon.jpg',
		},
		[DruidMajorGlyph.GlyphOfMoonfire]: {
			name: 'Glyph of Moonfire',
			description: 'Increases the periodic damage of your Moonfire ability by 75%, but initial damage is decreased by 90%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_nature_starfall.jpg',
		},
		[DruidMajorGlyph.GlyphOfNourish]: {
			name: 'Glyph of Nourish',
			description: 'Your Nourish heals an additional 6% for each of your heal over time effects present on the target.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_druid_nourish.jpg',
		},
		[DruidMajorGlyph.GlyphOfRake]: {
			name: 'Glyph of Rake',
			description: 'Your Rake ability prevents targets from fleeing.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_druid_disembowel.jpg',
		},
		[DruidMajorGlyph.GlyphOfRapidRejuvenation]: {
			name: 'Glyph of Rapid Rejuvenation',
			description: 'Your haste now reduces the time between the periodic healing ticks of your Rejuvenation spell.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_nature_rejuvenation.jpg',
		},
		[DruidMajorGlyph.GlyphOfRebirth]: {
			name: 'Glyph of Rebirth',
			description: 'Players resurrected by Rebirth are returned to life with 100% health.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_nature_reincarnation.jpg',
		},
		[DruidMajorGlyph.GlyphOfRegrowth]: {
			name: 'Glyph of Regrowth',
			description: 'Increases the healing of your Regrowth spell by 20% if your Regrowth effect is still active on the target.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_nature_resistnature.jpg',
		},
		[DruidMajorGlyph.GlyphOfRejuvenation]: {
			name: 'Glyph of Rejuvenation',
			description: 'While your Rejuvenation targets are below 50% health, you will heal them for an additional 50% health.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_nature_rejuvenation.jpg',
		},
		[DruidMajorGlyph.GlyphOfRip]: {
			name: 'Glyph of Rip',
			description: 'Increases the duration of your Rip ability by 4 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_ghoulfrenzy.jpg',
		},
		[DruidMajorGlyph.GlyphOfSavageRoar]: {
			name: 'Glyph of Savage Roar',
			description: 'Your Savage Roar ability grants an additional 3% bonus damage done.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_druid_skinteeth.jpg',
		},
		[DruidMajorGlyph.GlyphOfShred]: {
			name: 'Glyph of Shred',
			description: 'Each time you Shred, the duration of your Rip on the target is extended 2 sec, up to a maximum of 6 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_vampiricaura.jpg',
		},
		[DruidMajorGlyph.GlyphOfStarfall]: {
			name: 'Glyph of Starfall',
			description: 'Reduces the cooldown of Starfall by 30 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_druid_starfall.jpg',
		},
		[DruidMajorGlyph.GlyphOfStarfire]: {
			name: 'Glyph of Starfire',
			description: 'Your Starfire ability increases the duration of your Moonfire effect on the target by 3 sec, up to a maximum of 9 additional seconds.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_arcane_starfire.jpg',
		},
		[DruidMajorGlyph.GlyphOfSurvivalInstincts]: {
			name: 'Glyph of Survival Instincts',
			description: 'Your Survival Instincts ability grants an additional 15% of your maximum health.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_druid_tigersroar.jpg',
		},
		[DruidMajorGlyph.GlyphOfSwiftmend]: {
			name: 'Glyph of Swiftmend',
			description: 'Your Swiftmend ability no longer consumes a Rejuvenation or Regrowth effect from the target.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/inv_relics_idolofrejuvenation.jpg',
		},
		[DruidMajorGlyph.GlyphOfWildGrowth]: {
			name: 'Glyph of Wild Growth',
			description: 'Wild Growth can affect 1 additional target.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_druid_flourish.jpg',
		},
		[DruidMajorGlyph.GlyphOfWrath]: {
			name: 'Glyph of Wrath',
			description: 'Reduces the pushback suffered from damaging attacks while casting your Wrath spell by 50%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_nature_abolishmagic.jpg',
		},
	},
	minorGlyphs: {
		[DruidMinorGlyph.GlyphOfAquaticForm]: {
			name: 'Glyph of Aquatic Form',
			description: 'Increases your swim speed by 50% while in Aquatic Form.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_druid_aquaticform.jpg',
		},
		[DruidMinorGlyph.GlyphOfChallengingRoar]: {
			name: 'Glyph of Challenging Roar',
			description: 'Reduces the cooldown of your Challenging Roar ability by 30 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_druid_challangingroar.jpg',
		},
		[DruidMinorGlyph.GlyphOfDash]: {
			name: 'Glyph of Dash',
			description: 'Reduces the cooldown of your Dash ability by 20%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_druid_dash.jpg',
		},
		[DruidMinorGlyph.GlyphOfTheWild]: {
			name: 'Glyph of the Wild',
			description: 'Mana cost of your Mark of the Wild and Gift of the Wild spells reduced by 50%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_nature_regeneration.jpg',
		},
		[DruidMinorGlyph.GlyphOfThorns]: {
			name: 'Glyph of Thorns',
			description: 'Increases the duration of your Thorns ability by 50 min when cast on yourself.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_nature_thorns.jpg',
		},
		[DruidMinorGlyph.GlyphOfTyphoon]: {
			name: 'Glyph of Typhoon',
			description: 'Reduces the cost of your Typhoon spell by 8% and increases its radius by 10 yards, but it no longer knocks enemies back.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_druid_typhoon.jpg',
		},
		[DruidMinorGlyph.GlyphOfUnburdenedRebirth]: {
			name: 'Glyph of Unburdened Rebirth',
			description: 'Your Rebirth spell no longer requires a reagent.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_nature_reincarnation.jpg',
		},
	},
};
