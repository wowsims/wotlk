import {
	Consumes,
	Flask,
	Food,
	Glyphs,
	Potions,
} from '../core/proto/common.js';
import { SavedTalents } from '../core/proto/ui.js';

import {
	RestorationShaman_Rotation as RestorationShamanRotation,
	RestorationShaman_Options as RestorationShamanOptions,
	ShamanShield,
	ShamanMajorGlyph,
	ShamanMinorGlyph,
} from '../core/proto/shaman.js';

import {
	AirTotem,
	EarthTotem,
	FireTotem,
	WaterTotem,
	ShamanTotems,
} from '../core/proto/shaman.js';

import * as PresetUtils from '../core/preset_utils.js';

import BlankGear from './gear_sets/blank.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const DefaultGear = PresetUtils.makePresetGear('Blank', BlankGear);

export const DefaultRotation = RestorationShamanRotation.create({
	totems: ShamanTotems.create({
		earth: EarthTotem.StrengthOfEarthTotem,
		air: AirTotem.WrathOfAirTotem,
		fire: FireTotem.FlametongueTotem,
		water: WaterTotem.HealingStreamTotem,
	}),
	useEarthShield: true,
	useRiptide: true,
});

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/classic/talent-calc and copy the numbers in the url.
export const TankHealingTalents = {
	name: 'Tank Healing',
	data: SavedTalents.create({
		talentsString: '-30205033-05005331335010501122331251',
		glyphs: Glyphs.create({
			major1: ShamanMajorGlyph.GlyphOfEarthlivingWeapon,
			major2: ShamanMajorGlyph.GlyphOfEarthShield,
			major3: ShamanMajorGlyph.GlyphOfLesserHealingWave,
			minor2: ShamanMinorGlyph.GlyphOfWaterShield,
			minor1: ShamanMinorGlyph.GlyphOfRenewedLife,
			minor3: ShamanMinorGlyph.GlyphOfGhostWolf,
		}),
	}),
};
export const RaidHealingTalents = {
	name: 'Raid Healing',
	data: SavedTalents.create({
		talentsString: '-3020503-50005331335310501122331251',
		glyphs: Glyphs.create({
			major1: ShamanMajorGlyph.GlyphOfChainHeal,
			major2: ShamanMajorGlyph.GlyphOfEarthShield,
			major3: ShamanMajorGlyph.GlyphOfEarthlivingWeapon,
			minor2: ShamanMinorGlyph.GlyphOfWaterShield,
			minor1: ShamanMinorGlyph.GlyphOfRenewedLife,
			minor3: ShamanMinorGlyph.GlyphOfGhostWolf,
		}),
	}),
};

export const DefaultOptions = RestorationShamanOptions.create({
	shield: ShamanShield.WaterShield,
	bloodlust: true,
	earthShieldPPM: 0,
});

export const DefaultConsumes = Consumes.create({
	flask: Flask.FlaskUnknown,
	food: Food.FoodUnknown,
});
