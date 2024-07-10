import * as PresetUtils from '../core/preset_utils.js';
import {
	Consumes,
	Debuffs,
	Flask,
	Food,
	Glyphs,
	IndividualBuffs,
	Potions,
	RaidBuffs,
	TristateEffect,
	UnitReference,
} from '../core/proto/common.js';
import {
	HealingPriest_Options as Options,
	PriestMajorGlyph as MajorGlyph,
	PriestMinorGlyph as MinorGlyph,
} from '../core/proto/priest.js';
import { SavedTalents } from '../core/proto/ui.js';
// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.
import PreraidDiscGear from './gear_sets/preraid_disc.gear.json';
import PreraidHolyGear from './gear_sets/preraid_holy.gear.json';
export const DISC_PRERAID_PRESET = PresetUtils.makePresetGear('戒律Preraid预设', PreraidDiscGear, { talentTree: 0 });
export const HOLY_PRERAID_PRESET = PresetUtils.makePresetGear('神牧Preraid预设', PreraidHolyGear, { talentTree: 1 });
import P1DiscGear from './gear_sets/p1_disc.gear.json';
import P1HolyGear from './gear_sets/p1_holy.gear.json';
export const DISC_P1_PRESET = PresetUtils.makePresetGear('戒律P1预设', P1DiscGear, { talentTree: 0 });
export const HOLY_P1_PRESET = PresetUtils.makePresetGear('神牧P1预设', P1HolyGear, { talentTree: 1 });
import P2DiscGear from './gear_sets/p2_disc.gear.json';
import P2HolyGear from './gear_sets/p2_holy.gear.json';
export const DISC_P2_PRESET = PresetUtils.makePresetGear('戒律P2预设', P2DiscGear, { talentTree: 0 });
export const HOLY_P2_PRESET = PresetUtils.makePresetGear('神牧P2预设', P2HolyGear, { talentTree: 1 });
import P3DiscGear from './gear_sets/p3_disc.gear.json';
import P3HolyGear from './gear_sets/p3_holy.gear.json';
export const DISC_P3_PRESET = PresetUtils.makePresetGear('戒律P3预设', P3DiscGear, { talentTree: 0 });
export const HOLY_P3_PRESET = PresetUtils.makePresetGear('神牧P3预设', P3HolyGear, { talentTree: 1 });
import P4DiscGear from './gear_sets/p4_disc.gear.json';
import P4HolyGear from './gear_sets/p4_holy.gear.json';
export const DISC_P4_PRESET = PresetUtils.makePresetGear('戒律P4预设', P4DiscGear, { talentTree: 0 });
export const HOLY_P4_PRESET = PresetUtils.makePresetGear('神牧P4预设', P4HolyGear, { talentTree: 1 });

import DiscApl from './apls/disc.apl.json';
export const ROTATION_PRESET_DISC = PresetUtils.makePresetAPLRotation('戒律', DiscApl);
import HolyApl from './apls/holy.apl.json';
export const ROTATION_PRESET_HOLY = PresetUtils.makePresetAPLRotation('神牧', HolyApl);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/wotlk/talent-calc and copy the numbers in the url.
export const DiscTalents = {
	name: '戒律',
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
	name: '神牧',
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
