import {
	Consumes,
	Flask,
	Food,
	Glyphs,
	Spec
} from '../core/proto/common.js';
import { SavedTalents } from '../core/proto/ui.js';

import {
	DruidMajorGlyph,
	DruidMinorGlyph,
	FeralDruid_Options as FeralDruidOptions,
	FeralDruid_Rotation as FeralDruidRotation,
	FeralDruid_Rotation_AplType,
	FeralDruid_Rotation_BearweaveType,
	FeralDruid_Rotation_BiteModeType,
} from '../core/proto/druid.js';

import * as PresetUtils from '../core/preset_utils.js';

import P1Gear from './gear_sets/p1.gear.json';
import P2Gear from './gear_sets/p2.gear.json';
import P3Gear from './gear_sets/p3.gear.json';
import P4Gear from './gear_sets/p4.gear.json';
import PreraidGear from './gear_sets/preraid.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const PRERAID_PRESET = PresetUtils.makePresetGear('Preraid Preset', PreraidGear);
export const P1_PRESET = PresetUtils.makePresetGear('P1 Preset', P1Gear);
export const P2_PRESET = PresetUtils.makePresetGear('P2 Preset', P2Gear);
export const P3_PRESET = PresetUtils.makePresetGear('P3 Preset', P3Gear);
export const P4_PRESET = PresetUtils.makePresetGear('P4 Preset', P4Gear);

export const DefaultRotation = FeralDruidRotation.create({
	rotationType: FeralDruid_Rotation_AplType.SingleTarget,

	bearWeaveType: FeralDruid_Rotation_BearweaveType.None,
	minCombosForRip: 5,
	minCombosForBite: 5,

	useRake: true,
	useBite: true,
	mangleSpam: false,
	biteModeType: FeralDruid_Rotation_BiteModeType.Emperical,
	biteTime: 4.0,
	berserkBiteThresh: 25.0,
	berserkFfThresh: 15.0,
	powerbear: false,
	minRoarOffset: 12.0,
	ripLeeway: 3.0,
	maintainFaerieFire: true,
	hotUptime: 0.0,
	snekWeave: false,
	flowerWeave: false,
	raidTargets: 30,
	maxFfDelay: 0.1,
	prePopOoc: true,
});

export const ROTATION_PRESET_LEGACY_DEFAULT = PresetUtils.makePresetSimpleRotation('Legacy Default', Spec.SpecFeralDruid, DefaultRotation);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/wotlk/talent-calc and copy the numbers in the url.
export const StandardTalents = {
	name: 'Standard',
	data: SavedTalents.create({
		talentsString: '-543202132322010053120030310511-203503012',
		glyphs: Glyphs.create({
			major1: DruidMajorGlyph.GlyphOfOmenOfClarity,
			major2: DruidMajorGlyph.GlyphOfSavageRoar,
			major3: DruidMajorGlyph.GlyphOfShred,
			minor1: DruidMinorGlyph.GlyphOfDash,
			minor2: DruidMinorGlyph.GlyphOfTheWild,
			minor3: DruidMinorGlyph.GlyphOfUnburdenedRebirth,
		}),
	}),
};

export const DefaultOptions = FeralDruidOptions.create({
	latencyMs: 100,
	assumeBleedActive: true,
});

export const DefaultConsumes = Consumes.create({
	flask: Flask.FlaskUnknown,
	food: Food.FoodUnknown,
});
