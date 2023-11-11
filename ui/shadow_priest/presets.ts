import {
	Consumes,
	Debuffs,
	Flask,
	Food,
	Glyphs,
	IndividualBuffs,
	Potions,
	Profession,
	RaidBuffs,
	TristateEffect,
} from '../core/proto/common.js';
import { SavedTalents } from '../core/proto/ui.js';

import {
	ShadowPriest_Options as Options,
	PriestMajorGlyph as MajorGlyph,
	PriestMinorGlyph as MinorGlyph,
	ShadowPriest_Rotation as Rotation,
	ShadowPriest_Rotation_RotationType,
} from '../core/proto/priest.js';

import * as PresetUtils from '../core/preset_utils.js';

import BlankGear from './gear_sets/blank.gear.json';

import DefaultApl from './apls/default.apl.json'

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const BLANK_GEAR_PRESET = PresetUtils.makePresetGear('Blank', BlankGear);

export const DefaultRotation = Rotation.create({
	rotationType: ShadowPriest_Rotation_RotationType.Ideal,
});

export const ROTATION_PRESET_DEFAULT = PresetUtils.makePresetAPLRotation('Default', DefaultApl);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/classic/talent-calc and copy the numbers in the url.
export const StandardTalents = {
	name: 'Standard',
	data: SavedTalents.create({
		talentsString: '5042001303--5002505103501051',
		glyphs: Glyphs.create({
			major1: MajorGlyph.PriestMajorGlyphNone,
			major2: MajorGlyph.PriestMajorGlyphNone,
			major3: MajorGlyph.PriestMajorGlyphNone,
			minor1: MinorGlyph.PriestMinorGlyphNone,
			minor2: MinorGlyph.PriestMinorGlyphNone,
			minor3: MinorGlyph.PriestMinorGlyphNone,
		}),
	}),
};

export const DefaultOptions = Options.create({});

export const DefaultConsumes = Consumes.create({
	flask: Flask.FlaskOfTheFrostWyrm,
	food: Food.FoodFishFeast,
	defaultPotion: Potions.PotionOfSpeed,
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
	demonicPact: 500,
});

export const DefaultIndividualBuffs = IndividualBuffs.create({
	blessingOfKings: true,
	blessingOfWisdom: TristateEffect.TristateEffectImproved,
	blessingOfMight: TristateEffect.TristateEffectImproved,
	vampiricTouch: true,
});

export const DefaultDebuffs = Debuffs.create({
	sunderArmor: true,
	faerieFire: TristateEffect.TristateEffectImproved,
	bloodFrenzy: true,
	ebonPlaguebringer: true,
	heartOfTheCrusader: true,
	judgementOfWisdom: true,
	shadowMastery: true,
});

export const OtherDefaults = {
	channelClipDelay: 100,
	profession1: Profession.Engineering,
	profession2: Profession.Tailoring,
	nibelungAverageCasts: 11,
};
