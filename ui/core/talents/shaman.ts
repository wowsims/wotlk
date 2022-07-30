import { ShamanTalents, ShamanMajorGlyph, ShamanMinorGlyph } from '/wotlk/core/proto/shaman.js';

import { GlyphsConfig } from './glyphs_picker.js';
import { TalentsConfig, newTalentsConfig } from './talents_picker.js';

export const shamanTalentsConfig: TalentsConfig<ShamanTalents> = newTalentsConfig([
	{
		name: 'Elemental',
		backgroundUrl: 'https://wow.zamimg.com/images/wow/talents/backgrounds/classic/261.jpg',
		talents: [
			{
				fieldName: 'convection',
				location: {
					rowIdx: 0,
					colIdx: 1,
				},
				spellIds: [16039, 16109],
				maxPoints: 5,
			},
			{
				fieldName: 'concussion',
				location: {
					rowIdx: 0,
					colIdx: 2,
				},
				spellIds: [16035, 16105],
				maxPoints: 5,
			},
			{
				fieldName: 'callOfFlame',
				location: {
					rowIdx: 1,
					colIdx: 0,
				},
				spellIds: [16038, 16160, 16161],
				maxPoints: 3,
			},
			{
				//fieldName: 'elementalWarding',
				location: {
					rowIdx: 1,
					colIdx: 1,
				},
				spellIds: [28996, 28997, 28998],
				maxPoints: 3,
			},
			{
				fieldName: 'elementalDevastation',
				location: {
					rowIdx: 1,
					colIdx: 2,
				},
				spellIds: [30160, 29179, 29180],
				maxPoints: 3,
			},
			{
				fieldName: 'reverberation',
				location: {
					rowIdx: 2,
					colIdx: 0,
				},
				spellIds: [16040, 16113, 16114, 16115, 16116],
				maxPoints: 5,
			},
			{
				fieldName: 'elementalFocus',
				location: {
					rowIdx: 2,
					colIdx: 1,
				},
				spellIds: [16164],
				maxPoints: 1,
			},
			{
				fieldName: 'elementalFury',
				location: {
					rowIdx: 2,
					colIdx: 2,
				},
				spellIds: [16089, 60184, 60185, 60187, 60188],
				maxPoints: 5,
			},
			{
				fieldName: 'improvedFireNova',
				location: {
					rowIdx: 3,
					colIdx: 0,
				},
				spellIds: [16086, 16544],
				maxPoints: 2,
			},
			{
				//fieldName: 'eyeOfTheStorm',
				location: {
					rowIdx: 3,
					colIdx: 3,
				},
				spellIds: [29062, 29064],
				maxPoints: 3,
			},
			{
				//fieldName: 'elementalReach',
				location: {
					rowIdx: 4,
					colIdx: 0,
				},
				spellIds: [28999],
				maxPoints: 2,
			},
			{
				fieldName: 'callOfThunder',
				location: {
					rowIdx: 4,
					colIdx: 1,
				},
				prereqLocation: {
					rowIdx: 2,
					colIdx: 1,
				},
				spellIds: [16041],
				maxPoints: 1,
			},
			{
				fieldName: 'unrelentingStorm',
				location: {
					rowIdx: 4,
					colIdx: 3,
				},
				spellIds: [30664, 30665, 30666],
				maxPoints: 3,
			},
			{
				fieldName: 'elementalPrecision',
				location: {
					rowIdx: 5,
					colIdx: 0,
				},
				spellIds: [30672, 30673, 30674],
				maxPoints: 3,
			},
			{
				fieldName: 'lightningMastery',
				location: {
					rowIdx: 5,
					colIdx: 2,
				},
				prereqLocation: {
					rowIdx: 2,
					colIdx: 2,
				},
				spellIds: [16578],
				maxPoints: 5,
			},
			{
				fieldName: 'elementalMastery',
				location: {
					rowIdx: 6,
					colIdx: 1,
				},
				prereqLocation: {
					rowIdx: 4,
					colIdx: 1,
				},
				spellIds: [16166],
				maxPoints: 1,
			},
			{
				fieldName: 'stormEarthAndFire',
				location: {
					rowIdx: 6,
					colIdx: 2,
				},
				spellIds: [51483, 51485, 51486],
				maxPoints: 3,
			},
			{
				fieldName: 'boomingEchoes',
				location: {
					rowIdx: 7,
					colIdx: 0,
				},
				spellIds: [63370, 63372],
				maxPoints: 2,
			},
			{
				fieldName: 'elementalOath',
				location: {
					rowIdx: 7,
					colIdx: 1,
				},
				prereqLocation: {
					rowIdx: 6,
					colIdx: 1,
				},
				spellIds: [51466, 51470],
				maxPoints: 2,
			},
			{
				fieldName: 'lightningOverload',
				location: {
					rowIdx: 7,
					colIdx: 2,
				},
				spellIds: [30675, 30678, 30679],
				maxPoints: 3,
			},
			{
				// fieldName: 'astralShift',
				location: {
					rowIdx: 8,
					colIdx: 0,
				},
				spellIds: [51474, 51478, 51479],
				maxPoints: 3,
			},
			{
				fieldName: 'totemOfWrath',
				location: {
					rowIdx: 8,
					colIdx: 1,
				},
				spellIds: [30706],
				maxPoints: 1,
			},
			{
				fieldName: 'lavaFlows',
				location: {
					rowIdx: 8,
					colIdx: 2,
				},
				spellIds: [51480, 51481, 51482],
				maxPoints: 3,
			},
			{
				fieldName: 'shamanism',
				location: {
					rowIdx: 9,
					colIdx: 1,
				},
				spellIds: [62097],
				maxPoints: 5,
			},
			{
				fieldName: 'thunderstorm',
				location: {
					rowIdx: 10,
					colIdx: 1,
				},
				spellIds: [51490],
				maxPoints: 1,
			},
		],
	},
	{
		name: 'Enhancement',
		backgroundUrl: 'https://wow.zamimg.com/images/wow/talents/backgrounds/classic/263.jpg',
		talents: [
			{
				fieldName: 'enhancingTotems',
				location: {
					rowIdx: 0,
					colIdx: 0,
				},
				spellIds: [16259, 16295, 52456],
				maxPoints: 3,
			},
			{
				// fieldName: 'earthsGrasp',
				location: {
					rowIdx: 0,
					colIdx: 1,
				},
				spellIds: [16043, 16130],
				maxPoints: 2,
			},
			{
				fieldName: 'ancestralKnowledge',
				location: {
					rowIdx: 0,
					colIdx: 2,
				},
				spellIds: [17485],
				maxPoints: 5,
			},
			{
				//fieldName: 'guardianTotems',
				location: {
					rowIdx: 1,
					colIdx: 0,
				},
				spellIds: [16258, 16293],
				maxPoints: 2,
			},
			{
				fieldName: 'thunderingStrikes',
				location: {
					rowIdx: 1,
					colIdx: 1,
				},
				spellIds: [16255, 16302],
				maxPoints: 5,
			},
			{
				//fieldName: 'improvedGhostWolf',
				location: {
					rowIdx: 1,
					colIdx: 2,
				},
				spellIds: [16262, 16287],
				maxPoints: 2,
			},
			{
				fieldName: 'improvedShields',
				location: {
					rowIdx: 1,
					colIdx: 3,
				},
				spellIds: [16261, 16290, 51881],
				maxPoints: 3,
			},
			{
				fieldName: 'elementalWeapons',
				location: {
					rowIdx: 2,
					colIdx: 0,
				},
				spellIds: [16266, 29079, 29080],
				maxPoints: 3,
			},
			{
				fieldName: 'shamanisticFocus',
				location: {
					rowIdx: 2,
					colIdx: 2,
				},
				spellIds: [43338],
				maxPoints: 1,
			},
			{
				fieldName: 'anticipation',
				location: {
					rowIdx: 2,
					colIdx: 3,
				},
				spellIds: [16254, 16271],
				maxPoints: 3,
			},
			{
				fieldName: 'flurry',
				location: {
					rowIdx: 3,
					colIdx: 1,
				},
				prereqLocation: {
					rowIdx: 1,
					colIdx: 1,
				},
				spellIds: [16256, 16281],
				maxPoints: 5,
			},
			{
				fieldName: 'toughness',
				location: {
					rowIdx: 3,
					colIdx: 2,
				},
				spellIds: [16252, 16306],
				maxPoints: 5,
			},
			{
				fieldName: 'improvedWindfuryTotem',
				location: {
					rowIdx: 4,
					colIdx: 0,
				},
				spellIds: [29192, 29193],
				maxPoints: 2,
			},
			{
				fieldName: 'spiritWeapons',
				location: {
					rowIdx: 4,
					colIdx: 1,
				},
				spellIds: [16268],
				maxPoints: 1,
			},
			{
				fieldName: 'mentalDexterity',
				location: {
					rowIdx: 4,
					colIdx: 2,
				},
				spellIds: [51883],
				maxPoints: 3,
			},
			{
				fieldName: 'unleashedRage',
				location: {
					rowIdx: 5,
					colIdx: 0,
				},
				spellIds: [30802, 30808, 30809],
				maxPoints: 3,
			},
			{
				fieldName: 'weaponMastery',
				location: {
					rowIdx: 5,
					colIdx: 2,
				},
				spellIds: [29082, 29084, 29086],
				maxPoints: 3,
			},
			{
				fieldName: 'frozenPower',
				location: {
					rowIdx: 5,
					colIdx: 3,
				},
				spellIds: [63373, 63374],
				maxPoints: 2,
			},
			{
				fieldName: 'dualWieldSpecialization',
				location: {
					rowIdx: 6,
					colIdx: 0,
				},
				prereqLocation: {
					rowIdx: 6,
					colIdx: 1,
				},
				spellIds: [30816, 30818],
				maxPoints: 3,
			},
			{
				//fieldName: 'dualWield',
				location: {
					rowIdx: 6,
					colIdx: 1,
				},
				prereqLocation: {
					rowIdx: 4,
					colIdx: 1,
				},
				spellIds: [30798],
				maxPoints: 1,
			},
			{
				fieldName: 'stormstrike',
				location: {
					rowIdx: 6,
					colIdx: 2,
				},
				spellIds: [17364],
				maxPoints: 1,
			},
			{
				fieldName: 'staticShock',
				location: {
					rowIdx: 7,
					colIdx: 0,
				},
				spellIds: [51525],
				maxPoints: 3,
			},
			{
				fieldName: 'lavaLash',
				location: {
					rowIdx: 7,
					colIdx: 1,
				},
				prereqLocation: {
					rowIdx: 6,
					colIdx: 1,
				},
				spellIds: [60103],
				maxPoints: 1,
			},
			{
				fieldName: 'improvedStormstrike',
				location: {
					rowIdx: 7,
					colIdx: 2,
				},
				prereqLocation: {
					rowIdx: 6,
					colIdx: 2,
				},
				spellIds: [51521, 51522],
				maxPoints: 2,
			},
			{
				fieldName: 'mentalQuickness',
				location: {
					rowIdx: 8,
					colIdx: 0,
				},
				spellIds: [30812],
				maxPoints: 3,
			},
			{
				fieldName: 'shamanisticRage',
				location: {
					rowIdx: 8,
					colIdx: 1,
				},
				spellIds: [30823],
				maxPoints: 1,
			},
			{
				fieldName: 'earthenPower',
				location: {
					rowIdx: 8,
					colIdx: 2,
				},
				spellIds: [51523, 51524],
				maxPoints: 2,
			},
			{
				fieldName: 'maelstromWeapon',
				location: {
					rowIdx: 9,
					colIdx: 1,
				},
				spellIds: [51528],
				maxPoints: 5,
			},
			{
				fieldName: 'feralSpirit',
				location: {
					rowIdx: 10,
					colIdx: 1,
				},
				spellIds: [51533],
				maxPoints: 1,
			},
		],
	},
	{
		name: 'Restoration',
		backgroundUrl: 'https://wow.zamimg.com/images/wow/talents/backgrounds/classic/262.jpg',
		talents: [
			{
				//fieldName: 'improvedHealingWave',
				location: {
					rowIdx: 0,
					colIdx: 1,
				},
				spellIds: [16182, 16226],
				maxPoints: 5,
			},
			{
				fieldName: 'totemicFocus',
				location: {
					rowIdx: 0,
					colIdx: 2,
				},
				spellIds: [16173, 16222],
				maxPoints: 5,
			},
			{
				//fieldName: 'improvedReincarnation',
				location: {
					rowIdx: 1,
					colIdx: 0,
				},
				spellIds: [16184, 16209],
				maxPoints: 2,
			},
			{
				//fieldName: 'healingGrace',
				location: {
					rowIdx: 1,
					colIdx: 1,
				},
				spellIds: [29187, 29189, 29191],
				maxPoints: 3,
			},
			{
				//fieldName: 'tidalFocus',
				location: {
					rowIdx: 1,
					colIdx: 2,
				},
				spellIds: [16179, 16214],
				maxPoints: 5,
			},
			{
				fieldName: 'improvedWaterShield',
				location: {
					rowIdx: 2,
					colIdx: 0,
				},
				spellIds: [16180, 16196, 16198],
				maxPoints: 3,
			},
			{
				//fieldName: 'healingFocus',
				location: {
					rowIdx: 2,
					colIdx: 1,
				},
				spellIds: [16181, 16230, 16232],
				maxPoints: 3,
			},
			{
				//fieldName: 'tidalForce',
				location: {
					rowIdx: 2,
					colIdx: 2,
				},
				spellIds: [55198],
				maxPoints: 1,
			},
			{
				//fieldName: 'ancestralUealing',
				location: {
					rowIdx: 2,
					colIdx: 3,
				},
				spellIds: [16176, 16235],
				maxPoints: 3,
			},
			{
				fieldName: 'restorativeTotems',
				location: {
					rowIdx: 3,
					colIdx: 1,
				},
				spellIds: [16187, 16205],
				maxPoints: 3,
			},
			{
				fieldName: 'tidalMastery',
				location: {
					rowIdx: 3,
					colIdx: 2,
				},
				spellIds: [16194, 16218],
				maxPoints: 5,
			},
			{
				//fieldName: 'healingWay',
				location: {
					rowIdx: 4,
					colIdx: 0,
				},
				spellIds: [29206, 29205, 29202],
				maxPoints: 3,
			},
			{
				fieldName: 'naturesSwiftness',
				location: {
					rowIdx: 4,
					colIdx: 2,
				},
				spellIds: [16188],
				maxPoints: 1,
			},
			{
				//fieldName: 'focusedMind',
				location: {
					rowIdx: 4,
					colIdx: 3,
				},
				spellIds: [30864],
				maxPoints: 3,
			},
			{
				//fieldName: 'purification',
				location: {
					rowIdx: 5,
					colIdx: 2,
				},
				spellIds: [16178, 16210],
				maxPoints: 5,
			},
			{
				//fieldName: 'naturesGuardian',
				location: {
					rowIdx: 6,
					colIdx: 0,
				},
				spellIds: [30881, 30883],
				maxPoints: 5,
			},
			{
				fieldName: 'manaTideTotem',
				location: {
					rowIdx: 6,
					colIdx: 1,
				},
				prereqLocation: {
					rowIdx: 3,
					colIdx: 1,
				},
				spellIds: [16190],
				maxPoints: 1,
			},
			{
				//fieldName: 'cleanseSpirit',
				location: {
					rowIdx: 6,
					colIdx: 2,
				},
				prereqLocation: {
					rowIdx: 5,
					colIdx: 2,
				},
				spellIds: [51886],
				maxPoints: 1,
			},
			{
				fieldName: 'blessingOfTheEternals',
				location: {
					rowIdx: 7,
					colIdx: 0,
				},
				spellIds: [51554, 51555],
				maxPoints: 2,
			},
			{
				//fieldName: 'improvedChainHeal',
				location: {
					rowIdx: 7,
					colIdx: 1,
				},
				spellIds: [30872],
				maxPoints: 2,
			},
			{
				fieldName: 'naturesBlessing',
				location: {
					rowIdx: 7,
					colIdx: 2,
				},
				spellIds: [30867],
				maxPoints: 3,
			},
			{
				fieldName: 'ancestralAwakening',
				location: {
					rowIdx: 8,
					colIdx: 0,
				},
				spellIds: [51556],
				maxPoints: 3,
			},
			{
				//fieldName: 'earthShield',
				location: {
					rowIdx: 8,
					colIdx: 1,
				},
				spellIds: [974],
				maxPoints: 1,
			},
			{
				//fieldName: 'improvedEarthShield',
				location: {
					rowIdx: 8,
					colIdx: 2,
				},
				prereqLocation: {
					rowIdx: 8,
					colIdx: 1,
				},
				spellIds: [51560],
				maxPoints: 2,
			},
			{
				//fieldName: 'tidalWaves',
				location: {
					rowIdx: 9,
					colIdx: 1,
				},
				spellIds: [51562],
				maxPoints: 5,
			},
			{
				//fieldName: 'riptide',
				location: {
					rowIdx: 10,
					colIdx: 1,
				},
				spellIds: [61295],
				maxPoints: 1,
			},
		],
	},
]);

export const shamanGlyphsConfig: GlyphsConfig = {
	majorGlyphs: {
		[ShamanMajorGlyph.GlyphOfChainHeal]: {
			name: 'Glyph of Chain Heal',
			description: 'Your Chain Heal heals 1 additional target.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_nature_healingwavegreater.jpg',
		},
		[ShamanMajorGlyph.GlyphOfChainLightning]: {
			name: 'Glyph of Chain Lightning',
			description: 'Your Chain Lightning strikes 1 additional target.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_nature_chainlightning.jpg',
		},
		[ShamanMajorGlyph.GlyphOfEarthShield]: {
			name: 'Glyph of Earth Shield',
			description: 'Increases the amount healed by your Earth Shield by 20%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_nature_skinofearth.jpg',
		},
		[ShamanMajorGlyph.GlyphOfEarthlivingWeapon]: {
			name: 'Glyph of Earthliving Weapon',
			description: 'Increases the chance for your Earthliving weapon to trigger by 5%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shaman_earthlivingweapon.jpg',
		},
		[ShamanMajorGlyph.GlyphOfElementalMastery]: {
			name: 'Glyph of Elemental Mastery',
			description: 'Reduces the cooldown of your Elemental Mastery ability by 30 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_nature_wispheal.jpg',
		},
		[ShamanMajorGlyph.GlyphOfFeralSpirit]: {
			name: 'Glyph of Feral Spirit',
			description: 'Your spirit wolves gain an additional 30% of your attack power.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shaman_feralspirit.jpg',
		},
		[ShamanMajorGlyph.GlyphOfFireElementalTotem]: {
			name: 'Glyph of Fire Elemental Totem',
			description: 'Reduces the cooldown of your Fire Elemental Totem by 5 min.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_fire_elemental_totem.jpg',
		},
		[ShamanMajorGlyph.GlyphOfFireNova]: {
			name: 'Glyph of Fire Nova',
			description: 'Reduces the cooldown of your Fire Nova spell by 3 seconds.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_fire_sealoffire.jpg',
		},
		[ShamanMajorGlyph.GlyphOfFlameShock]: {
			name: 'Glyph of Flame Shock',
			description: 'Increases the critical strike damage bonus of your Flame Shock damage by 60%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_fire_flameshock.jpg',
		},
		[ShamanMajorGlyph.GlyphOfFlametongueWeapon]: {
			name: 'Glyph of Flametongue Weapon',
			description: 'Increases your spell critical strike chance by 2% on each of your weapons with Flametongue Weapon active.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_fire_flametounge.jpg',
		},
		[ShamanMajorGlyph.GlyphOfFrostShock]: {
			name: 'Glyph of Frost Shock',
			description: 'Increases the duration of your Frost Shock by 2 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_frost_frostshock.jpg',
		},
		[ShamanMajorGlyph.GlyphOfHealingStreamTotem]: {
			name: 'Glyph of Healing Stream Totem',
			description: 'Your Healing Stream Totem heals for an additional 20%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/inv_spear_04.jpg',
		},
		[ShamanMajorGlyph.GlyphOfHealingWave]: {
			name: 'Glyph of Healing Wave',
			description: 'Your Healing Wave also heals you for 20% of the healing effect when you heal someone else.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_nature_magicimmunity.jpg',
		},
		[ShamanMajorGlyph.GlyphOfHex]: {
			name: 'Glyph of Hex',
			description: 'Increases the damage your Hex target can take before the Hex effect is removed by 20%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shaman_hex.jpg',
		},
		[ShamanMajorGlyph.GlyphOfLava]: {
			name: 'Glyph of Lava',
			description: 'Your Lava Burst spell gains an additional 10% of your spellpower.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shaman_lavaburst.jpg',
		},
		[ShamanMajorGlyph.GlyphOfLavaLash]: {
			name: 'Glyph of Lava Lash',
			description: 'Damage on your Lava Lash is increased by an additional 10% if your weapon is enchanted with Flametongue.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_shaman_lavalash.jpg',
		},
		[ShamanMajorGlyph.GlyphOfLesserHealingWave]: {
			name: 'Glyph of Lesser Healing Wave',
			description: 'Your Lesser Healing Wave heals for 20% more if the target is also affected by Earth Shield.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_nature_healingwavelesser.jpg',
		},
		[ShamanMajorGlyph.GlyphOfLightningBolt]: {
			name: 'Glyph of Lightning Bolt',
			description: 'Increases the damage dealt by Lightning Bolt by 4%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_nature_lightning.jpg',
		},
		[ShamanMajorGlyph.GlyphOfLightningShield]: {
			name: 'Glyph of Lightning Shield',
			description: 'Increases the damage from Lightning Shield by 20%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_nature_lightningshield.jpg',
		},
		[ShamanMajorGlyph.GlyphOfManaTide]: {
			name: 'Glyph of Mana Tide',
			description: 'Your Mana Tide Totem grants an additional 1% of each target\'s maximum mana each time it pulses.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_nature_healingwavegreater.jpg',
		},
		[ShamanMajorGlyph.GlyphOfRiptide]: {
			name: 'Glyph of Riptide',
			description: 'Increases the duration of Riptide by 6 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_nature_riptide.jpg',
		},
		[ShamanMajorGlyph.GlyphOfShocking]: {
			name: 'Glyph of Shocking',
			description: 'Reduces your global cooldown when casting your shock spells by 0.5 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_nature_earthshock.jpg',
		},
		[ShamanMajorGlyph.GlyphOfStoneclawTotem]: {
			name: 'Glyph of Stoneclaw Totem',
			description: 'Your Stoneclaw Totem also places a damage absorb shield on you, equal to 4 times the strength of the shield it places on your totems.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_nature_stoneclawtotem.jpg',
		},
		[ShamanMajorGlyph.GlyphOfStormstrike]: {
			name: 'Glyph of Stormstrike',
			description: 'Increases the Nature damage bonus from your Stormstrike ability by an additional 8%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_shaman_stormstrike.jpg',
		},
		[ShamanMajorGlyph.GlyphOfThunder]: {
			name: 'Glyph of Thunder',
			description: 'Reduces the cooldown on Thunderstorm by 10 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shaman_thunderstorm.jpg',
		},
		[ShamanMajorGlyph.GlyphOfTotemOfWrath]: {
			name: 'Glyph of Totem of Wrath',
			description: 'When you cast Totem of Wrath, you gain 30% of the totem\'s bonus spell power for 5 min.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_fire_totemofwrath.jpg',
		},
		[ShamanMajorGlyph.GlyphOfWaterMastery]: {
			name: 'Glyph of Water Mastery',
			description: 'Increases the passive mana regeneration of your Water Shield spell by 30%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_shaman_watershield.jpg',
		},
		[ShamanMajorGlyph.GlyphOfWindfuryWeapon]: {
			name: 'Glyph of Windfury Weapon',
			description: 'Increases the chance per swing for Windfury Weapon to trigger by 2%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_nature_cyclone.jpg',
		},
	},
	minorGlyphs: {
		[ShamanMinorGlyph.GlyphOfAstralRecall]: {
			name: 'Glyph of Astral Recall',
			description: 'Cooldown of your Astral Recall spell reduced by 7.5 min.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_nature_astralrecal.jpg',
		},
		[ShamanMinorGlyph.GlyphOfGhostWolf]: {
			name: 'Glyph of Ghost Wolf',
			description: 'Your Ghost Wolf form regenerates an additional 1% of your maximum health every 5 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_nature_spiritwolf.jpg',
		},
		[ShamanMinorGlyph.GlyphOfRenewedLife]: {
			name: 'Glyph of Renewed Life',
			description: 'Your Reincarnation spell no longer requires a reagent.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_nature_reincarnation.jpg',
		},
		[ShamanMinorGlyph.GlyphOfThunderstorm]: {
			name: 'Glyph of Thunderstorm',
			description: 'Increases the mana you receive from your Thunderstorm spell by 2%, but it no longer knocks enemies back.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shaman_thunderstorm.jpg',
		},
		[ShamanMinorGlyph.GlyphOfWaterBreathing]: {
			name: 'Glyph of Water Breathing',
			description: 'Your Water Breathing spell no longer requires a reagent.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_demonbreath.jpg',
		},
		[ShamanMinorGlyph.GlyphOfWaterShield]: {
			name: 'Glyph of Water Shield',
			description: 'Increases the number of charges on your Water Shield spell by 1.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_shaman_watershield.jpg',
		},
		[ShamanMinorGlyph.GlyphOfWaterWalking]: {
			name: 'Glyph of Water Walking',
			description: 'Your Water Walking spell no longer requires a reagent.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_frost_windwalkon.jpg',
		},
	},
};
