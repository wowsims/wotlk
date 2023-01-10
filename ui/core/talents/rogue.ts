import { RogueTalents, RogueMajorGlyph, RogueMinorGlyph } from '../proto/rogue.js';

import { GlyphsConfig } from './glyphs_picker.js';
import { TalentsConfig, newTalentsConfig } from './talents_picker.js';

import RogueTalentJson from './trees/rogue.json';

export const rogueTalentsConfig: TalentsConfig<RogueTalents> = newTalentsConfig(RogueTalentJson);

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
