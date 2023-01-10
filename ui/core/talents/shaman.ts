import { ShamanTalents, ShamanMajorGlyph, ShamanMinorGlyph } from '../proto/shaman.js';

import { GlyphsConfig } from './glyphs_picker.js';
import { TalentsConfig, newTalentsConfig } from './talents_picker.js';

import ShamanTalentJson from './trees/shaman.json';

export const shamanTalentsConfig: TalentsConfig<ShamanTalents> = newTalentsConfig(ShamanTalentJson);

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
