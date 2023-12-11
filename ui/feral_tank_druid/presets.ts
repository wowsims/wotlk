import {
	Consumes,
	Flask,
	Food,
	UnitReference
} from '../core/proto/common.js';
import { SavedTalents } from '../core/proto/ui.js';

import {
	FeralTankDruid_Options as DruidOptions,
	FeralTankDruid_Rotation as DruidRotation,
} from '../core/proto/druid.js';

import * as PresetUtils from '../core/preset_utils.js';

import BlankGear from './gear_sets/blank.gear.json';

import DefaultApl from './apls/default.apl.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const DefaultGear = PresetUtils.makePresetGear('Blank', BlankGear);

export const DefaultRotation = DruidRotation.create({
	maulRageThreshold: 25,
	maintainDemoralizingRoar: true,
	lacerateTime: 8.0,
});

export const ROTATION_DEFAULT = PresetUtils.makePresetAPLRotation('Default', DefaultApl);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/classic/talent-calc and copy the numbers in the url.
export const StandardTalents = {
	name: 'Standard',
	data: SavedTalents.create({
		talentsString: '-503232132322010353120300313511-20350001',
	}),
};

export const DefaultOptions = DruidOptions.create({
	innervateTarget: UnitReference.create(),
	startingRage: 20,
});

export const DefaultConsumes = Consumes.create({
	flask: Flask.FlaskUnknown,
	food: Food.FoodUnknown,
});
