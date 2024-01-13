import {
	Consumes,
	Flask,
	Food,
	Profession,
	Spec,
	Potions,
	Conjured,
	AgilityElixir,
	StrengthBuff,
	WeaponImbue,
	SaygesFortune,
} from '../core/proto/common.js';
import { SavedTalents } from '../core/proto/ui.js';

import {
	FeralDruid_Options as FeralDruidOptions,
	FeralDruid_Rotation as FeralDruidRotation,
} from '../core/proto/druid.js';

import * as PresetUtils from '../core/preset_utils.js';

import BlankGear from './gear_sets/blank.gear.json';
import Phase1Gear from './gear_sets/p1.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const BlankPreset = PresetUtils.makePresetGear('Blank', BlankGear);
export const DefaultGear = PresetUtils.makePresetGear('Phase 1', Phase1Gear);

import DefaultApl from './apls/default.apl.json';
export const APL_ROTATION_DEFAULT = PresetUtils.makePresetAPLRotation('APL Default', DefaultApl);

export const DefaultRotation = FeralDruidRotation.create({
	maintainFaerieFire: false,
	minCombosForRip: 3,
	maxWaitTime: 2.0,
	preroarDuration: 26.0,
});

export const SIMPLE_ROTATION_DEFAULT = PresetUtils.makePresetSimpleRotation('Simple Default', Spec.SpecFeralDruid, DefaultRotation);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/classic/talent-calc and copy the numbers in the url.
export const StandardTalents = {
	name: 'Standard',
	data: SavedTalents.create({
		talentsString: '500005001--05',
	}),
};

export const DefaultOptions = FeralDruidOptions.create({
	latencyMs: 100,
	assumeBleedActive: true,
});

export const DefaultConsumes = Consumes.create({
	flask: Flask.FlaskUnknown,
	food: Food.FoodUnknown,
	defaultPotion: Potions.ManaPotion,
	defaultConjured: Conjured.ConjuredMinorRecombobulator,
        agilityElixir: AgilityElixir.ElixirOfLesserAgility,
        strengthBuff: StrengthBuff.ElixirOfOgresStrength,
        mainHandImbue: WeaponImbue.BlackfathomSharpeningStone,
});

export const OtherDefaults = {
	profession2: Profession.Leatherworking,
};
