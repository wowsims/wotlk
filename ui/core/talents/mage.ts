import { Spec } from '/wotlk/core/proto/common.js';
import { MageTalents, MageMajorGlyph, MageMinorGlyph } from '/wotlk/core/proto/mage.js';
import { Player } from '/wotlk/core/player.js';

import { GlyphsConfig, GlyphsPicker } from './glyphs_picker.js';
import { TalentsConfig, TalentsPicker, newTalentsConfig } from './talents_picker.js';

export class MageTalentsPicker extends TalentsPicker<Spec.SpecMage> {
	constructor(parent: HTMLElement, player: Player<Spec.SpecMage>) {
		super(parent, player, mageTalentsConfig);
	}
}

export class MageGlyphsPicker extends GlyphsPicker {
	constructor(parent: HTMLElement, player: Player<any>) {
		super(parent, player, mageGlyphsConfig);
	}
}

export const mageTalentsConfig: TalentsConfig<Spec.SpecMage> = newTalentsConfig([
	{
		name: 'Arcane',
		backgroundUrl: 'https://wow.zamimg.com/images/wow/talents/backgrounds/wotlk/81.jpg',
		talents: [
			{
				fieldName: 'arcaneSubtlety',
				location: {
					rowIdx: 0,
					colIdx: 0,
				},
				spellIds: [11210, 12592],
				maxPoints: 2,
			},
			{
				fieldName: 'arcaneFocus',
				location: {
					rowIdx: 0,
					colIdx: 1,
				},
				spellIds: [11222, 12839],
				maxPoints: 5,
			},
			{
				//fieldName: 'improvedArcaneMissiles',
				location: {
					rowIdx: 0,
					colIdx: 2,
				},
				spellIds: [11237, 12463, 12464, 16769],
				maxPoints: 5,
			},
			{
				fieldName: 'wandSpecialization',
				location: {
					rowIdx: 1,
					colIdx: 0,
				},
				spellIds: [6057, 6085],
				maxPoints: 2,
			},
			{
				fieldName: 'magicAbsorption',
				location: {
					rowIdx: 1,
					colIdx: 1,
				},
				spellIds: [29441, 29444],
				maxPoints: 5,
			},
			{
				fieldName: 'arcaneConcentration',
				location: {
					rowIdx: 1,
					colIdx: 2,
				},
				spellIds: [11213, 12574],
				maxPoints: 5,
			},
			{
				//fieldName: 'magicAttunement',
				location: {
					rowIdx: 2,
					colIdx: 0,
				},
				spellIds: [11247, 12606],
				maxPoints: 2,
			},
			{
				fieldName: 'arcaneImpact',
				location: {
					rowIdx: 2,
					colIdx: 1,
				},
				spellIds: [11242, 12467, 12469],
				maxPoints: 3,
			},
			{
				//fieldName: 'arcaneFortitude',
				location: {
					rowIdx: 2,
					colIdx: 3,
				},
				spellIds: [28574],
				maxPoints: 1,
			},
			{
				//fieldName: 'improvedManaShield',
				location: {
					rowIdx: 3,
					colIdx: 0,
				},
				spellIds: [11252, 12605],
				maxPoints: 2,
			},
			{
				//fieldName: 'improvedCounterspell',
				location: {
					rowIdx: 3,
					colIdx: 1,
				},
				spellIds: [11255, 12598],
				maxPoints: 2,
			},
			{
				fieldName: 'arcaneMeditation',
				location: {
					rowIdx: 3,
					colIdx: 3,
				},
				spellIds: [18462],
				maxPoints: 3,
			},
			{
				//fieldName: 'improvedBlink',
				location: {
					rowIdx: 4,
					colIdx: 0,
				},
				spellIds: [31569],
				maxPoints: 2,
			},
			{
				fieldName: 'presenceOfMind',
				location: {
					rowIdx: 4,
					colIdx: 1,
				},
				spellIds: [12043],
				maxPoints: 1,
			},
			{
				fieldName: 'arcaneMind',
				location: {
					rowIdx: 4,
					colIdx: 3,
				},
				spellIds: [11232, 12500],
				maxPoints: 5,
			},
			{
				//fieldName: 'prismaticCloak',
				location: {
					rowIdx: 5,
					colIdx: 0,
				},
				spellIds: [31574],
				maxPoints: 2,
			},
			{
				fieldName: 'arcaneInstability',
				location: {
					rowIdx: 5,
					colIdx: 1,
				},
				prereqLocation: {
					rowIdx: 4,
					colIdx: 1,
				},
				spellIds: [15058],
				maxPoints: 3,
			},
			{
				fieldName: 'arcanePotency',
				location: {
					rowIdx: 5,
					colIdx: 2,
				},
				prereqLocation: {
					rowIdx: 1,
					colIdx: 2,
				},
				spellIds: [31571],
				maxPoints: 3,
			},
			{
				fieldName: 'empoweredArcaneMissiles',
				location: {
					rowIdx: 6,
					colIdx: 0,
				},
				spellIds: [31579, 31582],
				maxPoints: 3,
			},
			{
				fieldName: 'arcanePower',
				location: {
					rowIdx: 6,
					colIdx: 1,
				},
				prereqLocation: {
					rowIdx: 5,
					colIdx: 1,
				},
				spellIds: [12042],
				maxPoints: 1,
			},
			{
				fieldName: 'spellPower',
				location: {
					rowIdx: 6,
					colIdx: 2,
				},
				spellIds: [35578, 35581],
				maxPoints: 2,
			},
			{
				fieldName: 'mindMastery',
				location: {
					rowIdx: 7,
					colIdx: 1,
				},
				spellIds: [31584],
				maxPoints: 5,
			},
			{
				//fieldName: 'slow',
				location: {
					rowIdx: 8,
					colIdx: 1,
				},
				spellIds: [31589],
				maxPoints: 1,
			},
		],
	},
	{
		name: 'Fire',
		backgroundUrl: 'https://wow.zamimg.com/images/wow/talents/backgrounds/wotlk/41.jpg',
		talents: [
			{
				fieldName: 'improvedFireball',
				location: {
					rowIdx: 0,
					colIdx: 1,
				},
				spellIds: [11069, 12338],
				maxPoints: 5,
			},
			{
				//fieldName: 'impact',
				location: {
					rowIdx: 0,
					colIdx: 2,
				},
				spellIds: [11103, 12357],
				maxPoints: 5,
			},
			{
				fieldName: 'ignite',
				location: {
					rowIdx: 1,
					colIdx: 0,
				},
				spellIds: [11119, 11120, 12846],
				maxPoints: 5,
			},
			{
				//fieldName: 'flameThrowing',
				location: {
					rowIdx: 1,
					colIdx: 1,
				},
				spellIds: [11100, 12353],
				maxPoints: 2,
			},
			{
				fieldName: 'improvedFireBlast',
				location: {
					rowIdx: 1,
					colIdx: 2,
				},
				spellIds: [11078, 11080, 12342],
				maxPoints: 3,
			},
			{
				fieldName: 'incineration',
				location: {
					rowIdx: 2,
					colIdx: 0,
				},
				spellIds: [18459],
				maxPoints: 2,
			},
			{
				fieldName: 'improvedFlamestrike',
				location: {
					rowIdx: 2,
					colIdx: 1,
				},
				spellIds: [11108, 12349],
				maxPoints: 3,
			},
			{
				fieldName: 'pyroblast',
				location: {
					rowIdx: 2,
					colIdx: 2,
				},
				spellIds: [11366],
				maxPoints: 1,
			},
			{
				fieldName: 'burningSoul',
				location: {
					rowIdx: 2,
					colIdx: 3,
				},
				spellIds: [11083, 12351],
				maxPoints: 2,
			},
			{
				fieldName: 'improvedScorch',
				location: {
					rowIdx: 3,
					colIdx: 0,
				},
				spellIds: [11095, 12872],
				maxPoints: 3,
			},
			{
				//fieldName: 'moltenShields',
				location: {
					rowIdx: 3,
					colIdx: 1,
				},
				spellIds: [11094, 13043],
				maxPoints: 2,
			},
			{
				fieldName: 'masterOfElements',
				location: {
					rowIdx: 3,
					colIdx: 3,
				},
				spellIds: [29074],
				maxPoints: 3,
			},
			{
				fieldName: 'playingWithFire',
				location: {
					rowIdx: 4,
					colIdx: 0,
				},
				spellIds: [31638],
				maxPoints: 3,
			},
			{
				fieldName: 'criticalMass',
				location: {
					rowIdx: 4,
					colIdx: 1,
				},
				spellIds: [11115, 11367],
				maxPoints: 3,
			},
			{
				fieldName: 'blastWave',
				location: {
					rowIdx: 4,
					colIdx: 2,
				},
				prereqLocation: {
					rowIdx: 2,
					colIdx: 2,
				},
				spellIds: [11113],
				maxPoints: 1,
			},
			{
				//fieldName: 'blazingSpeed',
				location: {
					rowIdx: 5,
					colIdx: 0,
				},
				spellIds: [31641],
				maxPoints: 2,
			},
			{
				fieldName: 'firePower',
				location: {
					rowIdx: 5,
					colIdx: 2,
				},
				spellIds: [11124, 12378, 12398],
				maxPoints: 5,
			},
			{
				fieldName: 'pyromaniac',
				location: {
					rowIdx: 6,
					colIdx: 0,
				},
				spellIds: [34293, 34295],
				maxPoints: 3,
			},
			{
				fieldName: 'combustion',
				location: {
					rowIdx: 6,
					colIdx: 1,
				},
				prereqLocation: {
					rowIdx: 4,
					colIdx: 1,
				},
				spellIds: [11129],
				maxPoints: 1,
			},
			{
				fieldName: 'moltenFury',
				location: {
					rowIdx: 6,
					colIdx: 2,
				},
				spellIds: [31679],
				maxPoints: 2,
			},
			{
				fieldName: 'empoweredFireball',
				location: {
					rowIdx: 7,
					colIdx: 2,
				},
				spellIds: [31656],
				maxPoints: 5,
			},
			{
				fieldName: 'dragonsBreath',
				location: {
					rowIdx: 8,
					colIdx: 1,
				},
				prereqLocation: {
					rowIdx: 6,
					colIdx: 1,
				},
				spellIds: [31661],
				maxPoints: 1,
			},
		],
	},
	{
		name: 'Frost',
		backgroundUrl: 'https://wow.zamimg.com/images/wow/talents/backgrounds/wotlk/61.jpg',
		talents: [
			{
				//fieldName: 'frostWarding',
				location: {
					rowIdx: 0,
					colIdx: 0,
				},
				spellIds: [11189, 28332],
				maxPoints: 2,
			},
			{
				fieldName: 'improvedFrostbolt',
				location: {
					rowIdx: 0,
					colIdx: 1,
				},
				spellIds: [11070, 12473, 16763, 16765],
				maxPoints: 5,
			},
			{
				fieldName: 'elementalPrecision',
				location: {
					rowIdx: 0,
					colIdx: 2,
				},
				spellIds: [29438],
				maxPoints: 3,
			},
			{
				fieldName: 'iceShards',
				location: {
					rowIdx: 1,
					colIdx: 0,
				},
				spellIds: [11207, 12672, 15047, 15052],
				maxPoints: 5,
			},
			{
				//fieldName: 'frostbite',
				location: {
					rowIdx: 1,
					colIdx: 1,
				},
				spellIds: [11071, 12496],
				maxPoints: 3,
			},
			{
				//fieldName: 'improvedFrostNova',
				location: {
					rowIdx: 1,
					colIdx: 2,
				},
				spellIds: [11165, 12475],
				maxPoints: 2,
			},
			{
				//fieldName: 'permafrost',
				location: {
					rowIdx: 1,
					colIdx: 3,
				},
				spellIds: [11175, 12569, 12571],
				maxPoints: 3,
			},
			{
				fieldName: 'piercingIce',
				location: {
					rowIdx: 2,
					colIdx: 0,
				},
				spellIds: [11151, 12952],
				maxPoints: 3,
			},
			{
				fieldName: 'icyVeins',
				location: {
					rowIdx: 2,
					colIdx: 1,
				},
				spellIds: [12472],
				maxPoints: 1,
			},
			{
				//fieldName: 'improvedBlizzard',
				location: {
					rowIdx: 2,
					colIdx: 3,
				},
				spellIds: [11185, 12487],
				maxPoints: 3,
			},
			{
				//fieldName: 'arcticReach',
				location: {
					rowIdx: 3,
					colIdx: 0,
				},
				spellIds: [16757],
				maxPoints: 2,
			},
			{
				fieldName: 'frostChanneling',
				location: {
					rowIdx: 3,
					colIdx: 1,
				},
				spellIds: [11160, 12518],
				maxPoints: 3,
			},
			{
				fieldName: 'shatter',
				location: {
					rowIdx: 3,
					colIdx: 2,
				},
				prereqLocation: {
					rowIdx: 1,
					colIdx: 2,
				},
				spellIds: [11170, 12982],
				maxPoints: 5,
			},
			{
				//fieldName: 'frozenCore',
				location: {
					rowIdx: 4,
					colIdx: 0,
				},
				spellIds: [31667],
				maxPoints: 3,
			},
			{
				fieldName: 'coldSnap',
				location: {
					rowIdx: 4,
					colIdx: 1,
				},
				spellIds: [11958],
				maxPoints: 1,
			},
			{
				fieldName: 'improvedConeOfCold',
				location: {
					rowIdx: 4,
					colIdx: 2,
				},
				spellIds: [11190, 12489],
				maxPoints: 3,
			},
			{
				fieldName: 'iceFloes',
				location: {
					rowIdx: 5,
					colIdx: 0,
				},
				spellIds: [31670, 31672],
				maxPoints: 2,
			},
			{
				fieldName: 'wintersChill',
				location: {
					rowIdx: 5,
					colIdx: 2,
				},
				spellIds: [11180, 28592],
				maxPoints: 5,
			},
			{
				//fieldName: 'iceBarrier',
				location: {
					rowIdx: 6,
					colIdx: 1,
				},
				prereqLocation: {
					rowIdx: 4,
					colIdx: 1,
				},
				spellIds: [11426],
				maxPoints: 1,
			},
			{
				fieldName: 'arcticWinds',
				location: {
					rowIdx: 6,
					colIdx: 2,
				},
				spellIds: [31674],
				maxPoints: 5,
			},
			{
				fieldName: 'empoweredFrostbolt',
				location: {
					rowIdx: 7,
					colIdx: 1,
				},
				spellIds: [31682],
				maxPoints: 5,
			},
			{
				fieldName: 'summonWaterElemental',
				location: {
					rowIdx: 8,
					colIdx: 1,
				},
				spellIds: [31687],
				maxPoints: 1,
			},
		],
	},
]);

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
