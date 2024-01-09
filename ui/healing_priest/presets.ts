import {
	Consumes,
	Debuffs,
	IndividualBuffs,
	Flask,
	Food,
	Glyphs,
	Potions,
	RaidBuffs,
	TristateEffect,
	UnitReference,
} from '../core/proto/common.js';
import { SavedTalents } from '../core/proto/ui.js';

import {
	HealingPriest_Options as Options,
	PriestMajorGlyph as MajorGlyph,
	PriestMinorGlyph as MinorGlyph,
} from '../core/proto/priest.js';

import * as PresetUtils from '../core/preset_utils.js';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

import PreraidDiscGear from './gear_sets/preraid_disc.gear.json';
import PreraidHolyGear from './gear_sets/preraid_holy.gear.json';
export const DISC_PRERAID_PRESET = PresetUtils.makePresetGear('Disc Preraid Preset', PreraidDiscGear, { talentTree: 0 });
export const HOLY_PRERAID_PRESET = PresetUtils.makePresetGear('Holy Preraid Preset', PreraidHolyGear, { talentTree: 1 });
import P1DiscGear from './gear_sets/p1_disc.gear.json';
import P1HolyGear from './gear_sets/p1_holy.gear.json';
export const DISC_P1_PRESET = PresetUtils.makePresetGear('Disc P1 Preset', P1DiscGear, { talentTree: 0 });
export const HOLY_P1_PRESET = PresetUtils.makePresetGear('Holy P1 Preset', P1HolyGear, { talentTree: 1 });
import P2DiscGear from './gear_sets/p2_disc.gear.json';
import P2HolyGear from './gear_sets/p2_holy.gear.json';
export const DISC_P2_PRESET = PresetUtils.makePresetGear('Disc P2 Preset', P2DiscGear, { talentTree: 0 });
export const HOLY_P2_PRESET = PresetUtils.makePresetGear('Holy P2 Preset', P2HolyGear, { talentTree: 1 });
import P3DiscGear from './gear_sets/p3_disc.gear.json';
import P3HolyGear from './gear_sets/p3_holy.gear.json';
export const DISC_P3_PRESET = PresetUtils.makePresetGear('Disc P3 Preset', P3DiscGear, { talentTree: 0 });
export const HOLY_P3_PRESET = PresetUtils.makePresetGear('Holy P3 Preset', P3HolyGear, { talentTree: 1 });
import P4DiscGear from './gear_sets/p4_disc.gear.json';
import P4HolyGear from './gear_sets/p4_holy.gear.json';
export const DISC_P4_PRESET = PresetUtils.makePresetGear('Disc P4 Preset', P4DiscGear, { talentTree: 0 });
export const HOLY_P4_PRESET = PresetUtils.makePresetGear('Holy P4 Preset', P4HolyGear, { talentTree: 1 });

import DiscApl from './apls/disc.apl.json';
export const ROTATION_PRESET_DISC = PresetUtils.makePresetAPLRotation('Disc', DiscApl);
import HolyApl from './apls/holy.apl.json';
export const ROTATION_PRESET_HOLY = PresetUtils.makePresetAPLRotation('Holy', HolyApl);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/wotlk/talent-calc and copy the numbers in the url.
export const DiscTalents = {
	name: 'Disc',
	data: SavedTalents.create({
		talentsString: '0503203130300512301313231251-2351010303',
		glyphs: Glyphs.create({
			major1: MajorGlyph.GlyphOfPowerWordShield,
			major2: MajorGlyph.GlyphOfFlashHeal,
			major3: MajorGlyph.GlyphOfPenance,
			minor1: MinorGlyph.GlyphOfFortitude,
			minor2: MinorGlyph.GlyphOfShadowfiend,
			minor3: MinorGlyph.GlyphOfFading,
		}),
	}),
};
export const HolyTalents = {
	name: 'Holy',
	data: SavedTalents.create({
		talentsString: '05032031103-234051032002152530004311051',
		glyphs: Glyphs.create({
			major1: MajorGlyph.GlyphOfPrayerOfHealing,
			major2: MajorGlyph.GlyphOfRenew,
			major3: MajorGlyph.GlyphOfCircleOfHealing,
			minor1: MinorGlyph.GlyphOfFortitude,
			minor2: MinorGlyph.GlyphOfShadowfiend,
			minor3: MinorGlyph.GlyphOfFading,
		}),
	}),
};

export const DefaultOptions = Options.create({
	useInnerFire: true,
	useShadowfiend: true,
	rapturesPerMinute: 5,

	powerInfusionTarget: UnitReference.create(),
});

export const DefaultConsumes = Consumes.create({
	flask: Flask.FlaskOfTheFrostWyrm,
	food: Food.FoodFishFeast,
	defaultPotion: Potions.RunicManaInjector,
	prepopPotion: Potions.PotionOfWildMagic,
});

export const DefaultRaidBuffs = RaidBuffs.create({
	giftOfTheWild: TristateEffect.TristateEffectImproved,
	powerWordFortitude: TristateEffect.TristateEffectImproved,
	strengthOfEarthTotem: TristateEffect.TristateEffectImproved,
	arcaneBrilliance: true,
	divineSpirit: true,
	trueshotAura: true,
	leaderOfThePack: TristateEffect.TristateEffectImproved,
	icyTalons: true,
	totemOfWrath: true,
	moonkinAura: TristateEffect.TristateEffectImproved,
	wrathOfAirTotem: true,
	sanctifiedRetribution: true,
	bloodlust: true,
	demonicPactSp: 500,
});

export const DefaultIndividualBuffs = IndividualBuffs.create({
	blessingOfKings: true,
	blessingOfWisdom: TristateEffect.TristateEffectImproved,
	vampiricTouch: true,
});

export const DefaultDebuffs = Debuffs.create({
});
