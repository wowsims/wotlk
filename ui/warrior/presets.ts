import {
	Consumes,
	Faction,
	Flask,
	Food,
	Glyphs,
	Potions,
} from '../core/proto/common.js';
import { SavedTalents } from '../core/proto/ui.js';

import {
	WarriorShout,
	Warrior_Rotation as WarriorRotation,
	Warrior_Rotation_SunderArmor as SunderArmor,
	Warrior_Options as WarriorOptions,
	WarriorMajorGlyph,
	WarriorMinorGlyph,
	Warrior_Rotation_StanceOption as StanceOption,
	Warrior_Rotation_MainGcd as MainGcd,
} from '../core/proto/warrior.js';

import * as PresetUtils from '../core/preset_utils.js';

import Empty from './gear_sets/empty.json';

import FuryApl from './apls/fury.apl.json';
import FurySunderApl from './apls/fury_sunder.apl.json';
import ArmsApl from './apls/arms.apl.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const EMPTY_PRESET = PresetUtils.makePresetGear('Empty', Empty, { talentTree: 0 });

export const ROTATION_FURY = PresetUtils.makePresetAPLRotation('Fury', FuryApl, { talentTree: 1 });
export const ROTATION_FURY_SUNDER = PresetUtils.makePresetAPLRotation('Fury + Sunder', FurySunderApl, { talentTree: 1 });
export const ROTATION_ARMS = PresetUtils.makePresetAPLRotation('Arms', ArmsApl, { talentTree: 0 });

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