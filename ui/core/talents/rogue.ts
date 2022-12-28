import { RogueTalents, RogueMajorGlyph, RogueMinorGlyph } from '../proto/rogue.js';

import { GlyphsConfig } from './glyphs_picker.js';
import { TalentsConfig, newTalentsConfig } from './talents_picker.js';

export const rogueTalentsConfig: TalentsConfig<RogueTalents> = newTalentsConfig([
	{
		name: 'Assassination',
		backgroundUrl: 'https://wow.zamimg.com/images/wow/talents/backgrounds/wrath/182.jpg',
		talents: [
			{
				fieldName: 'improvedEviscerate',
				location: {
					rowIdx: 0,
					colIdx: 0,
				},
				spellIds: [14162],
				maxPoints: 3,
			},
			{
				fieldName: 'remorselessAttacks',
				location: {
					rowIdx: 0,
					colIdx: 1,
				},
				spellIds: [14144, 14148],
				maxPoints: 2,
			},
			{
				fieldName: 'malice',
				location: {
					rowIdx: 0,
					colIdx: 2,
				},
				spellIds: [14138],
				maxPoints: 5,
			},
			{
				fieldName: 'ruthlessness',
				location: {
					rowIdx: 1,
					colIdx: 0,
				},
				spellIds: [14156, 14160],
				maxPoints: 3,
			},
			{
				fieldName: 'bloodSpatter',
				location: {
					rowIdx: 1,
					colIdx: 1,
				},
				spellIds: [51632, 51633],
				maxPoints: 2,
			},
			{
				fieldName: 'puncturingWounds',
				location: {
					rowIdx: 1,
					colIdx: 3,
				},
				spellIds: [13733, 13865],
				maxPoints: 3,
			},
			{
				fieldName: 'vigor',
				location: {
					rowIdx: 2,
					colIdx: 0,
				},
				spellIds: [14983],
				maxPoints: 1,
			},
			{
				fieldName: 'improvedExposeArmor',
				location: {
					rowIdx: 2,
					colIdx: 1,
				},
				spellIds: [14168],
				maxPoints: 2,
			},
			{
				fieldName: 'lethality',
				location: {
					rowIdx: 2,
					colIdx: 2,
				},
				prereqLocation: {
					rowIdx: 0,
					colIdx: 2,
				},
				spellIds: [14128, 14132, 14135],
				maxPoints: 5,
			},
			{
				fieldName: 'vilePoisons',
				location: {
					rowIdx: 3,
					colIdx: 1,
				},
				spellIds: [16513, 16514, 16515],
				maxPoints: 3,
			},
			{
				fieldName: 'improvedPoisons',
				location: {
					rowIdx: 3,
					colIdx: 2,
				},
				spellIds: [14113, 14114, 14115, 14116, 14117],
				maxPoints: 5,
			},
			{
				fieldName: 'fleetFooted',
				location: {
					rowIdx: 4,
					colIdx: 0,
				},
				spellIds: [31208],
				maxPoints: 2,
			},
			{
				fieldName: 'coldBlood',
				location: {
					rowIdx: 4,
					colIdx: 1,
				},
				spellIds: [14177],
				maxPoints: 1,
			},
			{
				fieldName: 'improvedKidneyShot',
				location: {
					rowIdx: 4,
					colIdx: 2,
				},
				spellIds: [14174],
				maxPoints: 3,
			},
			{
				fieldName: 'quickRecovery',
				location: {
					rowIdx: 4,
					colIdx: 3,
				},
				spellIds: [31244],
				maxPoints: 2,
			},
			{
				fieldName: 'sealFate',
				location: {
					rowIdx: 5,
					colIdx: 1,
				},
				prereqLocation: {
					rowIdx: 4,
					colIdx: 1,
				},
				spellIds: [14186, 14190, 14193],
				maxPoints: 5,
			},
			{
				fieldName: 'murder',
				location: {
					rowIdx: 5,
					colIdx: 2,
				},
				spellIds: [14158, 14159],
				maxPoints: 2,
			},
			{
				fieldName: 'deadlyBrew',
				location: {
					rowIdx: 6,
					colIdx: 0,
				},
				spellIds: [51625, 51626],
				maxPoints: 2,
			},
			{
				fieldName: 'overkill',
				location: {
					rowIdx: 6,
					colIdx: 1,
				},
				spellIds: [58426],
				maxPoints: 1,
			},
			{
				fieldName: 'deadenedNerves',
				location: {
					rowIdx: 6,
					colIdx: 2,
				},
				spellIds: [31380, 31382, 31383],
				maxPoints: 3,
			},
			{
				fieldName: 'focusedAttacks',
				location: {
					rowIdx: 7,
					colIdx: 0,
				},
				spellIds: [51634, 51635, 51636],
				maxPoints: 3,
			},
			{
				fieldName: 'findWeakness',
				location: {
					rowIdx: 7,
					colIdx: 2,
				},
				spellIds: [31234, 31235, 31236],
				maxPoints: 3,
			},
			{
				fieldName: 'masterPoisoner',
				location: {
					rowIdx: 8,
					colIdx: 0,
				},
				spellIds: [31226, 31227, 58410],
				maxPoints: 3,
			},
			{
				fieldName: 'mutilate',
				location: {
					rowIdx: 8,
					colIdx: 1,
				},
				prereqLocation: {
					rowIdx: 6,
					colIdx: 1,
				},
				spellIds: [1329],
				maxPoints: 1,
			},
			{
				fieldName: 'turnTheTables',
				location: {
					rowIdx: 8,
					colIdx: 2,
				},
				spellIds: [51627, 51628, 51629],
				maxPoints: 3,
			},
			{
				fieldName: 'cutToTheChase',
				location: {
					rowIdx: 9,
					colIdx: 1,
				},
				spellIds: [51664, 51665, 51666, 51667, 51668],
				maxPoints: 5,
			},
			{
				fieldName: 'hungerForBlood',
				location: {
					rowIdx: 10,
					colIdx: 1,
				},
				spellIds: [51662],
				maxPoints: 1,
			},
		],
	},
	{
		name: 'Combat',
		backgroundUrl: 'https://wow.zamimg.com/images/wow/talents/backgrounds/wrath/181.jpg',
		talents: [
			{
				fieldName: 'improvedGouge',
				location: {
					rowIdx: 0,
					colIdx: 0,
				},
				spellIds: [13741, 13793, 13792],
				maxPoints: 3,
			},
			{
				fieldName: 'improvedSinisterStrike',
				location: {
					rowIdx: 0,
					colIdx: 1,
				},
				spellIds: [13732, 13863],
				maxPoints: 2,
			},
			{
				fieldName: 'dualWieldSpecialization',
				location: {
					rowIdx: 0,
					colIdx: 2,
				},
				spellIds: [13715, 13848, 13849, 13851, 13852],
				maxPoints: 5,
			},
			{
				fieldName: 'improvedSliceAndDice',
				location: {
					rowIdx: 1,
					colIdx: 0,
				},
				spellIds: [14165, 14166],
				maxPoints: 2,
			},
			{
				fieldName: 'deflection',
				location: {
					rowIdx: 1,
					colIdx: 1,
				},
				spellIds: [13713, 13853, 13854],
				maxPoints: 3,
			},
			{
				fieldName: 'precision',
				location: {
					rowIdx: 1,
					colIdx: 3,
				},
				spellIds: [13705, 13832, 13843, 13844, 13845],
				maxPoints: 5,
			},
			{
				fieldName: 'endurance',
				location: {
					rowIdx: 2,
					colIdx: 0,
				},
				spellIds: [13742, 13872],
				maxPoints: 2,
			},
			{
				fieldName: 'riposte',
				location: {
					rowIdx: 2,
					colIdx: 1,
				},
				prereqLocation: {
					rowIdx: 1,
					colIdx: 1,
				},
				spellIds: [14251],
				maxPoints: 1,
			},
			{
				fieldName: 'closeQuartersCombat',
				location: {
					rowIdx: 2,
					colIdx: 2,
				},
				prereqLocation: {
					rowIdx: 0,
					colIdx: 2,
				},
				spellIds: [13706, 13804, 13805, 13806, 13807],
				maxPoints: 5,
			},
			{
				fieldName: 'improvedKick',
				location: {
					rowIdx: 3,
					colIdx: 0,
				},
				spellIds: [13754, 13867],
				maxPoints: 2,
			},
			{
				fieldName: 'improvedSprint',
				location: {
					rowIdx: 3,
					colIdx: 1,
				},
				spellIds: [13743, 13875],
				maxPoints: 2,
			},
			{
				fieldName: 'lightningReflexes',
				location: {
					rowIdx: 3,
					colIdx: 2,
				},
				spellIds: [13712, 13788, 13789],
				maxPoints: 3,
			},
			{
				fieldName: 'aggression',
				location: {
					rowIdx: 3,
					colIdx: 3,
				},
				spellIds: [18427, 18428, 18429, 61330, 61331],
				maxPoints: 5,
			},
			{
				fieldName: 'maceSpecialization',
				location: {
					rowIdx: 4,
					colIdx: 0,
				},
				spellIds: [13709, 13800, 13801, 13802, 13803],
				maxPoints: 5,
			},
			{
				fieldName: 'bladeFlurry',
				location: {
					rowIdx: 4,
					colIdx: 1,
				},
				spellIds: [13877],
				maxPoints: 1,
			},
			{
				fieldName: 'hackAndSlash',
				location: {
					rowIdx: 4,
					colIdx: 2,
				},
				spellIds: [13960, 13961, 13962, 13963, 13964],
				maxPoints: 5,
			},
			{
				fieldName: 'weaponExpertise',
				location: {
					rowIdx: 5,
					colIdx: 1,
				},
				prereqLocation: {
					rowIdx: 4,
					colIdx: 1,
				},
				spellIds: [30919],
				maxPoints: 2,
			},
			{
				fieldName: 'bladeTwisting',
				location: {
					rowIdx: 5,
					colIdx: 2,
				},
				spellIds: [31124, 31126],
				maxPoints: 2,
			},
			{
				fieldName: 'vitality',
				location: {
					rowIdx: 6,
					colIdx: 0,
				},
				spellIds: [31122, 31123, 61329],
				maxPoints: 3,
			},
			{
				fieldName: 'adrenalineRush',
				location: {
					rowIdx: 6,
					colIdx: 1,
				},
				spellIds: [13750],
				maxPoints: 1,
			},
			{
				fieldName: 'nervesOfSteel',
				location: {
					rowIdx: 6,
					colIdx: 2,
				},
				spellIds: [31130],
				maxPoints: 2,
			},
			{
				fieldName: 'throwingSpecialization',
				location: {
					rowIdx: 7,
					colIdx: 0,
				},
				spellIds: [5952, 51679],
				maxPoints: 2,
			},
			{
				fieldName: 'combatPotency',
				location: {
					rowIdx: 7,
					colIdx: 2,
				},
				spellIds: [35541, 35550, 35551, 35552, 35553],
				maxPoints: 5,
			},
			{
				fieldName: 'unfairAdvantage',
				location: {
					rowIdx: 8,
					colIdx: 0,
				},
				spellIds: [51672, 51674],
				maxPoints: 2,
			},
			{
				fieldName: 'surpriseAttacks',
				location: {
					rowIdx: 8,
					colIdx: 1,
				},
				prereqLocation: {
					rowIdx: 6,
					colIdx: 1,
				},
				spellIds: [32601],
				maxPoints: 1,
			},
			{
				fieldName: 'savageCombat',
				location: {
					rowIdx: 8,
					colIdx: 2,
				},
				spellIds: [51682, 58413],
				maxPoints: 2,
			},
			{
				fieldName: 'preyOnTheWeak',
				location: {
					rowIdx: 9,
					colIdx: 1,
				},
				spellIds: [51685, 51686, 51687, 51688, 51689],
				maxPoints: 5,
			},
			{
				fieldName: 'killingSpree',
				location: {
					rowIdx: 10,
					colIdx: 1,
				},
				spellIds: [51690],
				maxPoints: 1,
			},
		],
	},
	{
		name: 'Subtlety',
		backgroundUrl: 'https://wow.zamimg.com/images/wow/talents/backgrounds/wrath/183.jpg',
		talents: [
			{
				fieldName: 'relentlessStrikes',
				location: {
					rowIdx: 0,
					colIdx: 0,
				},
				spellIds: [14179, 58422, 58423, 58424, 58425],
				maxPoints: 5,
			},
			{
				fieldName: 'masterOfDeception',
				location: {
					rowIdx: 0,
					colIdx: 1,
				},
				spellIds: [13958, 13970, 13971],
				maxPoints: 3,
			},
			{
				fieldName: 'opportunity',
				location: {
					rowIdx: 0,
					colIdx: 2,
				},
				spellIds: [14057, 14072],
				maxPoints: 2,
			},
			{
				fieldName: 'sleightOfHand',
				location: {
					rowIdx: 1,
					colIdx: 0,
				},
				spellIds: [30892, 30893],
				maxPoints: 2,
			},
			{
				fieldName: 'dirtyTricks',
				location: {
					rowIdx: 1,
					colIdx: 1,
				},
				spellIds: [14076, 14094],
				maxPoints: 2,
			},
			{
				fieldName: 'camouflage',
				location: {
					rowIdx: 1,
					colIdx: 2,
				},
				spellIds: [13975, 14062, 14063],
				maxPoints: 3,
			},
			{
				fieldName: 'elusiveness',
				location: {
					rowIdx: 2,
					colIdx: 0,
				},
				spellIds: [13981, 14066],
				maxPoints: 2,
			},
			{
				fieldName: 'ghostlyStrike',
				location: {
					rowIdx: 2,
					colIdx: 1,
				},
				spellIds: [14278],
				maxPoints: 1,
			},
			{
				fieldName: 'serratedBlades',
				location: {
					rowIdx: 2,
					colIdx: 2,
				},
				spellIds: [14171, 14172, 14173],
				maxPoints: 3,
			},
			{
				fieldName: 'setup',
				location: {
					rowIdx: 3,
					colIdx: 0,
				},
				spellIds: [13983, 14070, 14071],
				maxPoints: 3,
			},
			{
				fieldName: 'initiative',
				location: {
					rowIdx: 3,
					colIdx: 1,
				},
				spellIds: [13976, 13979, 13980],
				maxPoints: 3,
			},
			{
				fieldName: 'improvedAmbush',
				location: {
					rowIdx: 3,
					colIdx: 2,
				},
				spellIds: [14079, 14080],
				maxPoints: 2,
			},
			{
				fieldName: 'heightenedSenses',
				location: {
					rowIdx: 4,
					colIdx: 0,
				},
				spellIds: [30894, 30895],
				maxPoints: 2,
			},
			{
				fieldName: 'preparation',
				location: {
					rowIdx: 4,
					colIdx: 1,
				},
				spellIds: [14185],
				maxPoints: 1,
			},
			{
				fieldName: 'dirtyDeeds',
				location: {
					rowIdx: 4,
					colIdx: 2,
				},
				spellIds: [14082, 14083],
				maxPoints: 2,
			},
			{
				fieldName: 'hemorrhage',
				location: {
					rowIdx: 4,
					colIdx: 3,
				},
				prereqLocation: {
					rowIdx: 2,
					colIdx: 2,
				},
				spellIds: [16511],
				maxPoints: 1,
			},
			{
				fieldName: 'masterOfSubtlety',
				location: {
					rowIdx: 5,
					colIdx: 0,
				},
				spellIds: [31221, 31222, 31223],
				maxPoints: 3,
			},
			{
				fieldName: 'deadliness',
				location: {
					rowIdx: 5,
					colIdx: 2,
				},
				spellIds: [30902, 30903, 30904, 30905, 30906],
				maxPoints: 5,
			},
			{
				fieldName: 'envelopingShadows',
				location: {
					rowIdx: 6,
					colIdx: 0,
				},
				spellIds: [31211, 31212, 31213],
				maxPoints: 3,
			},
			{
				fieldName: 'premeditation',
				location: {
					rowIdx: 6,
					colIdx: 1,
				},
				prereqLocation: {
					rowIdx: 4,
					colIdx: 1,
				},
				spellIds: [14183],
				maxPoints: 1,
			},
			{
				fieldName: 'cheatDeath',
				location: {
					rowIdx: 6,
					colIdx: 2,
				},
				spellIds: [31228, 31229, 31230],
				maxPoints: 3,
			},
			{
				fieldName: 'sinisterCalling',
				location: {
					rowIdx: 7,
					colIdx: 1,
				},
				prereqLocation: {
					rowIdx: 6,
					colIdx: 1,
				},
				spellIds: [31216, 31217, 31218, 31219, 31220],
				maxPoints: 5,
			},
			{
				fieldName: 'waylay',
				location: {
					rowIdx: 7,
					colIdx: 2,
				},
				spellIds: [51692, 51696],
				maxPoints: 2,
			},
			{
				fieldName: 'honorAmongThieves',
				location: {
					rowIdx: 8,
					colIdx: 0,
				},
				spellIds: [51698, 51700, 51701],
				maxPoints: 3,
			},
			{
				fieldName: 'shadowstep',
				location: {
					rowIdx: 8,
					colIdx: 1,
				},
				spellIds: [36554],
				maxPoints: 1,
			},
			{
				fieldName: 'filthyTricks',
				location: {
					rowIdx: 8,
					colIdx: 2,
				},
				spellIds: [58414, 58415],
				maxPoints: 2,
			},
			{
				fieldName: 'slaughterFromTheShadows',
				location: {
					rowIdx: 9,
					colIdx: 1,
				},
				spellIds: [51708, 51709, 51710, 51711, 51712],
				maxPoints: 5,
			},
			{
				fieldName: 'shadowDance',
				location: {
					rowIdx: 10,
					colIdx: 1,
				},
				spellIds: [51713],
				maxPoints: 1,
			},
		],
	},
]);

export const rogueGlyphsConfig: GlyphsConfig = {
	majorGlyphs: {
		[RogueMajorGlyph.GlyphOfAdrenalineRush]: {
			name: 'Glyph of Adrenaline Rush',
			description: 'Increases the duration of Adrenaline Rush by 5 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_shadowworddominate.jpg',
		},
		[RogueMajorGlyph.GlyphOfAmbush]: {
			name: 'Glyph of Ambush',
			description: 'Increases the range on Ambush by 5 yards.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_rogue_ambush.jpg',
		},
		[RogueMajorGlyph.GlyphOfBackstab]: {
			name: 'Glyph of Backstab',
			description: 'Your Backstab increases the duration of your Rupture effect on the target by 2 sec, up to a maximum of 6 additional sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_backstab.jpg',
		},
		[RogueMajorGlyph.GlyphOfBladeFlurry]: {
			name: 'Glyph of Blade Flurry',
			description: 'Reduces the energy cost of Blade Flurry by 100%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_warrior_punishingblow.jpg',
		},
		[RogueMajorGlyph.GlyphOfCloakOfShadows]: {
			name: 'Glyph of Cloak of Shadows',
			description: 'While Cloak of Shadows is active, you take 40% less physical damage.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_nethercloak.jpg',
		},
		[RogueMajorGlyph.GlyphOfCripplingPoison]: {
			name: 'Glyph of Crippling Poison',
			description: 'Increases the chance to inflict your target with Crippling Poison by an additional 20%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_poisonsting.jpg',
		},
		[RogueMajorGlyph.GlyphOfDeadlyThrow]: {
			name: 'Glyph of Deadly Throw',
			description: 'Increases the slowing effect on Deadly Throw by 20%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/inv_throwingknife_06.jpg',
		},
		[RogueMajorGlyph.GlyphOfEvasion]: {
			name: 'Glyph of Evasion',
			description: 'Increases the duration of Evasion by 5 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_shadowward.jpg',
		},
		[RogueMajorGlyph.GlyphOfEviscerate]: {
			name: 'Glyph of Eviscerate',
			description: 'Increases the critical strike chance of Eviscerate by 10%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_rogue_eviscerate.jpg',
		},
		[RogueMajorGlyph.GlyphOfExposeArmor]: {
			name: 'Glyph of Expose Armor',
			description: 'Increases the duration of Expose Armor by 12 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_warrior_riposte.jpg',
		},
		[RogueMajorGlyph.GlyphOfFanOfKnives]: {
			name: 'Glyph of Fan of Knives',
			description: 'Increases the damage done by Fan of Knives by 20%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_rogue_fanofknives.jpg',
		},
		[RogueMajorGlyph.GlyphOfFeint]: {
			name: 'Glyph of Feint',
			description: 'Reduces the energy cost of Feint by 20.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_rogue_feint.jpg',
		},
		[RogueMajorGlyph.GlyphOfGarrote]: {
			name: 'Glyph of Garrote',
			description: 'Reduces the duration of your Garrote ability by 3 sec and increases the total damage it deals by 20%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_rogue_garrote.jpg',
		},
		[RogueMajorGlyph.GlyphOfGhostlyStrike]: {
			name: 'Glyph of Ghostly Strike',
			description: 'Increases the damage dealt by Ghostly Strike by 40% and the duration of its effect by 4 sec, but increases its cooldown by 10 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_curse.jpg',
		},
		[RogueMajorGlyph.GlyphOfGouge]: {
			name: 'Glyph of Gouge',
			description: 'Reduces the energy cost of Gouge by 15.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_gouge.jpg',
		},
		[RogueMajorGlyph.GlyphOfHemorrhage]: {
			name: 'Glyph of Hemorrhage',
			description: 'Increases the damage bonus against targets afflicted by Hemorrhage by 40%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_lifedrain.jpg',
		},
		[RogueMajorGlyph.GlyphOfHungerForBlood]: {
			name: 'Glyph of Hunger For Blood',
			description: 'Increases the bonus damage from Hunger For Blood by 3%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_rogue_hungerforblood.jpg',
		},
		[RogueMajorGlyph.GlyphOfKillingSpree]: {
			name: 'Glyph of Killing Spree',
			description: 'Reduces the cooldown on Killing Spree by 45 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_rogue_murderspree.jpg',
		},
		[RogueMajorGlyph.GlyphOfMutilate]: {
			name: 'Glyph of Mutilate',
			description: 'Reduces the cost of Mutilate by 5 energy.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_rogue_shadowstrikes.jpg',
		},
		[RogueMajorGlyph.GlyphOfPreparation]: {
			name: 'Glyph of Preparation',
			description: 'Your Preparation ability also instantly resets the cooldown of Blade Flurry, Dismantle, and Kick.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_antishadow.jpg',
		},
		[RogueMajorGlyph.GlyphOfRupture]: {
			name: 'Glyph of Rupture',
			description: 'Increases the duration of Rupture by 4 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_rogue_rupture.jpg',
		},
		[RogueMajorGlyph.GlyphOfSap]: {
			name: 'Glyph of Sap',
			description: 'Increases the duration of Sap by 20 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_sap.jpg',
		},
		[RogueMajorGlyph.GlyphOfShadowDance]: {
			name: 'Glyph of Shadow Dance',
			description: 'Increases the duration of Shadow Dance by 2 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_rogue_shadowdance.jpg',
		},
		[RogueMajorGlyph.GlyphOfSinisterStrike]: {
			name: 'Glyph of Sinister Strike',
			description: 'Your Sinister Strike critical strikes have a 50% chance to add an additional combo point.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_ritualofsacrifice.jpg',
		},
		[RogueMajorGlyph.GlyphOfSliceAndDice]: {
			name: 'Glyph of Slice and Dice',
			description: 'Increases the duration of Slice and Dice by 3 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_rogue_slicedice.jpg',
		},
		[RogueMajorGlyph.GlyphOfSprint]: {
			name: 'Glyph of Sprint',
			description: 'Increases the movement speed of your Sprint ability by an additional 30%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_rogue_sprint.jpg',
		},
		[RogueMajorGlyph.GlyphOfTricksOfTheTrade]: {
			name: 'Glyph of Tricks of the Trade',
			description: 'The bonus damage and threat redirection granted by your Tricks of the Trade ability lasts an additional 4 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_rogue_tricksofthetrade.jpg',
		},
		[RogueMajorGlyph.GlyphOfVigor]: {
			name: 'Glyph of Vigor',
			description: 'Vigor grants an additional 10 maximum energy.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_nature_earthbindtotem.jpg',
		},
	},
	minorGlyphs: {
		[RogueMinorGlyph.GlyphOfBlurredSpeed]: {
			name: 'Glyph of Blurred Speed',
			description: 'Enables you to walk on water while your Sprint ability is active.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_rogue_sprint.jpg',
		},
		[RogueMinorGlyph.GlyphOfDistract]: {
			name: 'Glyph of Distract',
			description: 'Increases the range of your Distract ability by 5 yards.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_rogue_distract.jpg',
		},
		[RogueMinorGlyph.GlyphOfPickLock]: {
			name: 'Glyph of Pick Lock',
			description: 'Reduces the cast time of your Pick Lock ability by 100%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_nature_moonkey.jpg',
		},
		[RogueMinorGlyph.GlyphOfPickPocket]: {
			name: 'Glyph of Pick Pocket',
			description: 'Increases the range of your Pick Pocket ability by 5 yards.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/inv_misc_bag_11.jpg',
		},
		[RogueMinorGlyph.GlyphOfSafeFall]: {
			name: 'Glyph of Safe Fall',
			description: 'Increases the distance your Safe Fall ability allows you to fall without taking damage.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/inv_feather_01.jpg',
		},
		[RogueMinorGlyph.GlyphOfVanish]: {
			name: 'Glyph of Vanish',
			description: 'Increases your movement speed by 30% while the Vanish effect is active.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_vanish.jpg',
		},
	},
};
