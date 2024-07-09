import * as PresetUtils from '../core/preset_utils.js';
import {
	Consumes,
	Flask,
	Food,
	Glyphs,
	Potions,
	Spec,
} from '../core/proto/common.js';
import {
	DruidMajorGlyph,
	DruidMinorGlyph,
	FeralDruid_Options as FeralDruidOptions,
	FeralDruid_Rotation as FeralDruidRotation,
	FeralDruid_Rotation_AplType,
	FeralDruid_Rotation_BearweaveType,
	FeralDruid_Rotation_BiteModeType,
} from '../core/proto/druid.js';
import { SavedTalents } from '../core/proto/ui.js';
// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.
import PreraidGear from './gear_sets/preraid.gear.json';
export const PRERAID_PRESET = PresetUtils.makePresetGear('Preraid预设', PreraidGear);
import P1Gear from './gear_sets/p1.gear.json';
export const P1_PRESET = PresetUtils.makePresetGear('P1预设', P1Gear);
import P2Gear from './gear_sets/p2.gear.json';
export const P2_PRESET = PresetUtils.makePresetGear('P2预设', P2Gear);
import P3Gear from './gear_sets/p3.gear.json';
export const P3_PRESET = PresetUtils.makePresetGear('P3预设', P3Gear);
import P4Gear from './gear_sets/p4.gear.json';
export const P4_PRESET = PresetUtils.makePresetGear('P4预设', P4Gear);

import DefaultApl from './apls/default.apl.json';
export const APL_ROTATION_DEFAULT = PresetUtils.makePresetAPLRotation('APL预设', DefaultApl);

import CustomExampleApl from './apls/custom_apl_example.apl.json';
export const APL_ROTATION_CUSTOM_EXAMPLE = PresetUtils.makePresetAPLRotation('自定义APL模板', CustomExampleApl);

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

export const SIMPLE_ROTATION_DEFAULT = PresetUtils.makePresetSimpleRotation('预设', Spec.SpecFeralDruid, DefaultRotation);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/wotlk/talent-calc and copy the numbers in the url.
export const StandardTalents = {
	name: '标准天赋',
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
	flask: Flask.FlaskOfEndlessRage,
	food: Food.FoodDragonfinFilet,
	defaultPotion: Potions.PotionOfSpeed,
});
