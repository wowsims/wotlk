import { MageTalents, MageMajorGlyph, MageMinorGlyph } from '../proto/mage.js';

import { GlyphsConfig } from './glyphs_picker.js';
import { TalentsConfig, newTalentsConfig } from './talents_picker.js';

import MageTalentJson from './trees/mage.json';

export const mageTalentsConfig: TalentsConfig<MageTalents> = newTalentsConfig(MageTalentJson);

export const mageGlyphsConfig: GlyphsConfig = {
	majorGlyphs: {
		[MageMajorGlyph.GlyphOfArcaneBarrage]: {
			name: 'Glyph of Arcane Barrage',
			description: 'Reduces the mana cost of Arcane Barrage by 20%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_mage_arcanebarrage.jpg',
		},
		[MageMajorGlyph.GlyphOfArcaneBlast]: {
			name: 'Glyph of Arcane Blast',
			description: 'Increases the damage from your Arcane Blast buff by 3%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_arcane_blast.jpg',
		},
		[MageMajorGlyph.GlyphOfArcaneExplosion]: {
			name: 'Glyph of Arcane Explosion',
			description: 'Reduces mana cost of Arcane Explosion by 10%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_nature_wispsplode.jpg',
		},
		[MageMajorGlyph.GlyphOfArcaneMissiles]: {
			name: 'Glyph of Arcane Missiles',
			description: 'Increases the critical strike damage bonus of Arcane Missiles by 25%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_nature_starfall.jpg',
		},
		[MageMajorGlyph.GlyphOfArcanePower]: {
			name: 'Glyph of Arcane Power',
			description: 'Increases the duration of Arcane Power by 3 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_nature_lightning.jpg',
		},
		[MageMajorGlyph.GlyphOfBlink]: {
			name: 'Glyph of Blink',
			description: 'Increases the distance you travel with the Blink spell by 5 yards.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_arcane_blink.jpg',
		},
		[MageMajorGlyph.GlyphOfDeepFreeze]: {
			name: 'Glyph of Deep Freeze',
			description: 'Increases the range of Deep Freeze by 10 yards.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_mage_deepfreeze.jpg',
		},
		[MageMajorGlyph.GlyphOfEternalWater]: {
			name: 'Glyph of Eternal Water',
			description: 'Your Summon Water Elemental now lasts indefinitely, but your Water Elemental can no longer cast Freeze.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_frost_summonwaterelemental_2.jpg',
		},
		[MageMajorGlyph.GlyphOfEvocation]: {
			name: 'Glyph of Evocation',
			description: 'Your Evocation ability also causes you to regain 60% of your health over its duration.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_nature_purge.jpg',
		},
		[MageMajorGlyph.GlyphOfFireBlast]: {
			name: 'Glyph of Fire Blast',
			description: 'Increases the critical strike chance of Fire Blast by 50% when the target is stunned or incapacitated.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_fire_fireball.jpg',
		},
		[MageMajorGlyph.GlyphOfFireball]: {
			name: 'Glyph of Fireball',
			description: 'Reduces the casting time of your Fireball spell by 0.15 sec, but removes the damage over time effect.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_fire_flamebolt.jpg',
		},
		[MageMajorGlyph.GlyphOfFrostNova]: {
			name: 'Glyph of Frost Nova',
			description: 'Your Frost Nova targets can take an additional 20% damage before the Frost Nova effect automatically breaks.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_frost_frostnova.jpg',
		},
		[MageMajorGlyph.GlyphOfFrostbolt]: {
			name: 'Glyph of Frostbolt',
			description: 'Increases the damage dealt by Frostbolt by 5%, but removes the slowing effect.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_frost_frostbolt02.jpg',
		},
		[MageMajorGlyph.GlyphOfFrostfire]: {
			name: 'Glyph of Frostfire',
			description: 'Increases the initial damage dealt by Frostfire Bolt by 2% and its critical strike chance by 2%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_mage_frostfirebolt.jpg',
		},
		[MageMajorGlyph.GlyphOfIceArmor]: {
			name: 'Glyph of Ice Armor',
			description: 'Your Ice Armor and Frost Armor spells grant an additional 50% armor and resistance.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_frost_frostarmor02.jpg',
		},
		[MageMajorGlyph.GlyphOfIceBarrier]: {
			name: 'Glyph of Ice Barrier',
			description: 'Increases the amount of damage absorbed by your Ice Barrier by 30%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_ice_lament.jpg',
		},
		[MageMajorGlyph.GlyphOfIceBlock]: {
			name: 'Glyph of Ice Block',
			description: 'Your Frost Nova cooldown is now reset every time you use Ice Block.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_frost_frost.jpg',
		},
		[MageMajorGlyph.GlyphOfIceLance]: {
			name: 'Glyph of Ice Lance',
			description: 'Your Ice Lance now causes 4 times damage against frozen targets higher level than you instead of triple damage.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_frost_frostblast.jpg',
		},
		[MageMajorGlyph.GlyphOfIcyVeins]: {
			name: 'Glyph of Icy Veins',
			description: 'Your Icy Veins ability also removes all movement slowing and cast time slowing effects.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_frost_coldhearted.jpg',
		},
		[MageMajorGlyph.GlyphOfInvisibility]: {
			name: 'Glyph of Invisibility',
			description: 'Increases the duration of the Invisibility effect by 10 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_mage_invisibility.jpg',
		},
		[MageMajorGlyph.GlyphOfLivingBomb]: {
			name: 'Glyph of Living Bomb',
			description: 'The periodic damage from your Living Bomb can now be critical strikes.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_mage_livingbomb.jpg',
		},
		[MageMajorGlyph.GlyphOfMageArmor]: {
			name: 'Glyph of Mage Armor',
			description: 'Your Mage Armor spell grants an additional 20% mana regeneration while casting.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_magearmor.jpg',
		},
		[MageMajorGlyph.GlyphOfManaGem]: {
			name: 'Glyph of Mana Gem',
			description: 'Increases the mana received from using a mana gem by 40%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/inv_misc_gem_stone_01.jpg',
		},
		[MageMajorGlyph.GlyphOfMirrorImage]: {
			name: 'Glyph of Mirror Image',
			description: 'Your Mirror Image spell now creates a 4th copy.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_magic_lesserinvisibilty.jpg',
		},
		[MageMajorGlyph.GlyphOfMoltenArmor]: {
			name: 'Glyph of Molten Armor',
			description: 'Your Molten Armor grants an additional 20% of your spirit as critical strike rating.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_mage_moltenarmor.jpg',
		},
		[MageMajorGlyph.GlyphOfPolymorph]: {
			name: 'Glyph of Polymorph',
			description: 'Your Polymorph spell also removes all damage over time effects from the target.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_nature_polymorph.jpg',
		},
		[MageMajorGlyph.GlyphOfRemoveCurse]: {
			name: 'Glyph of Remove Curse',
			description: 'Your Remove Curse spell also makes the target immune to all curses for 4 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_nature_removecurse.jpg',
		},
		[MageMajorGlyph.GlyphOfScorch]: {
			name: 'Glyph of Scorch',
			description: 'Increases the damage of your Scorch spell by 20%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_fire_soulburn.jpg',
		},
		[MageMajorGlyph.GlyphOfWaterElemental]: {
			name: 'Glyph of Water Elemental',
			description: 'Reduces the cooldown of your Summon Water Elemental spell by 30 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_frost_summonwaterelemental_2.jpg',
		},
	},
	minorGlyphs: {
		[MageMinorGlyph.GlyphOfArcaneIntellect]: {
			name: 'Glyph of Arcane Intellect',
			description: 'Reduces the mana cost of your Arcane Intellect and Arcane Brilliance spells by 50%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_holy_magicalsentry.jpg',
		},
		[MageMinorGlyph.GlyphOfBlastWave]: {
			name: 'Glyph of Blast Wave',
			description: 'The mana cost of your Blast Wave spell is reduced by 15%, but it no longer knocks enemies back.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_holy_excorcism_02.jpg',
		},
		[MageMinorGlyph.GlyphOfFireWard]: {
			name: 'Glyph of Fire Ward',
			description: 'You have an additional 5% chance to reflect Fire spells while your Fire Ward is active.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_fire_firearmor.jpg',
		},
		[MageMinorGlyph.GlyphOfFrostArmor]: {
			name: 'Glyph of Frost Armor',
			description: 'Increases the duration of your Frost Armor and Ice Armor spells by 30 min.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_frost_frostarmor02.jpg',
		},
		[MageMinorGlyph.GlyphOfFrostWard]: {
			name: 'Glyph of Frost Ward',
			description: 'You have an additional 5% chance to reflect Frost spells while your Frost Ward is active.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_frost_frostward.jpg',
		},
		[MageMinorGlyph.GlyphOfSlowFall]: {
			name: 'Glyph of Slow Fall',
			description: 'Your Slow Fall spell no longer requires a reagent.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_magic_featherfall.jpg',
		},
		[MageMinorGlyph.GlyphOfThePenguin]: {
			name: 'Glyph of the Penguin',
			description: 'Your Polymorph: Sheep spell polymorphs the target into a penguin instead.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/inv_misc_penguinpet.jpg',
		},
	},
};
