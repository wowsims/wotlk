import {
	Consumes,
	Flask,
	Food,
} from '../core/proto/common.js';
import { SavedTalents } from '../core/proto/ui.js';

import {
	WarriorShout,
	Warrior_Options as WarriorOptions,
} from '../core/proto/warrior.js';

import * as PresetUtils from '../core/preset_utils.js';

import BlankGear from './gear_sets/blank.gear.json';

import DefaultAPL from './apls/default.apl.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const GearArmsDefault = PresetUtils.makePresetGear('Blank', BlankGear, { talentTree: 0 });
export const GearFuryDefault = PresetUtils.makePresetGear('Blank', BlankGear, { talentTree: 1 });

export const RotationArmsDefault = PresetUtils.makePresetAPLRotation('Default', DefaultAPL);
export const RotationFuryDefault = PresetUtils.makePresetAPLRotation('Default', DefaultAPL);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/classic/talent-calc and copy the numbers in the url.
export const Talent25 = {
	name: 'Level 25',
	data: SavedTalents.create({
		talentsString: '303220201-03',
	}),
};

export const DefaultOptions = WarriorOptions.create({
	startingRage: 0,
	useRecklessness: true,
	shout: WarriorShout.WarriorShoutCommanding,
});

export const DefaultConsumes = Consumes.create({
	flask: Flask.FlaskUnknown,
	food: Food.FoodUnknown,
});