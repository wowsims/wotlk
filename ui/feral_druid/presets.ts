import {
	Consumes,
	Flask,
	Food,
	Spec
} from '../core/proto/common.js';
import { SavedTalents } from '../core/proto/ui.js';

import {
	FeralDruid_Options as FeralDruidOptions,
	FeralDruid_Rotation as FeralDruidRotation,
	FeralDruid_Rotation_AplType,
	FeralDruid_Rotation_BearweaveType,
	FeralDruid_Rotation_BiteModeType,
} from '../core/proto/druid.js';

import * as PresetUtils from '../core/preset_utils.js';

import BlankGear from './gear_sets/blank.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const DefaultGear = PresetUtils.makePresetGear('Blank', BlankGear);

import DefaultApl from './apls/default.apl.json';
export const APL_ROTATION_DEFAULT = PresetUtils.makePresetAPLRotation('APL Default', DefaultApl);

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

export const SIMPLE_ROTATION_DEFAULT = PresetUtils.makePresetSimpleRotation('Simple Default', Spec.SpecFeralDruid, DefaultRotation);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/classic/talent-calc and copy the numbers in the url.
export const StandardTalents = {
	name: 'Standard',
	data: SavedTalents.create({
		talentsString: '-543202132322010053120030310511-203503012',
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
