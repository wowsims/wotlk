import {
	Consumes,
	Flask,
	Food,
	Glyphs,
	Profession
} from '../core/proto/common.js';
import { SavedTalents } from '../core/proto/ui.js';

import {
	Mage_Options_ArmorType as ArmorType,
	MageMajorGlyph,
	MageMinorGlyph,
	Mage_Options as MageOptions
} from '../core/proto/mage.js';

import * as PresetUtils from '../core/preset_utils.js';

import DefaultBlankGear from './gear_sets/blank.gear.json';

import DefaultAPL from './apls/default.apl.json';

export const DEFAULT_GEAR = PresetUtils.makePresetGear('Default', DefaultBlankGear, { talentTree: 0 });
export const ROTATION_PRESET_DEFAULT = PresetUtils.makePresetAPLRotation('Default', DefaultAPL, { talentTree: 0 });

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/classic/talent-calc and copy the numbers in the url.
export const DefaultTalents = {
	name: 'Fire',
	data: SavedTalents.create({
		talentsString: '230025030002-5052000123033151-003',
		glyphs: Glyphs.create({
			major1: MageMajorGlyph.MageMajorGlyphNone,
			major2: MageMajorGlyph.MageMajorGlyphNone,
			major3: MageMajorGlyph.MageMajorGlyphNone,
			minor1: MageMinorGlyph.MageMinorGlyphNone,
			minor2: MageMinorGlyph.MageMinorGlyphNone,
			minor3: MageMinorGlyph.MageMinorGlyphNone,
		}),
	}),
};

export const DefaultOptions = MageOptions.create({
	armor: ArmorType.MageArmor,
});

export const DefaultConsumes = Consumes.create({
	flask: Flask.FlaskUnknown,
	food: Food.FoodUnknown,
});

export const OtherDefaults = {
	distanceFromTarget: 20,
	profession1: Profession.Engineering,
	profession2: Profession.Tailoring,
};
