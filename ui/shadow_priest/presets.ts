import {
	Consumes,
	Debuffs,
	Flask,
	Food,
	Glyphs,
	IndividualBuffs,
	Profession,
	RaidBuffs,
	TristateEffect,
	WeaponBuff,
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
	flask: Flask.FlaskOfSupremePower,
	weaponBuff: WeaponBuff.BrillianWizardOil,
	food: Food.FoodNightfinSoup,
	spellPowerBuff: true,
	shadowPowerBuff: true,
});

export const DefaultRaidBuffs = RaidBuffs.create({
	giftOfTheWild: TristateEffect.TristateEffectImproved,
	powerWordFortitude: TristateEffect.TristateEffectImproved,
	arcaneBrilliance: true,
	divineSpirit: true,
	moonkinAura: true,
});

export const DefaultIndividualBuffs = IndividualBuffs.create({
	blessingOfKings: true,
	blessingOfWisdom: TristateEffect.TristateEffectImproved,
});

export const DefaultDebuffs = Debuffs.create({
	judgementOfWisdom: true,
});

export const OtherDefaults = {
	channelClipDelay: 100,
	profession1: Profession.Engineering,
	profession2: Profession.Tailoring,
};
