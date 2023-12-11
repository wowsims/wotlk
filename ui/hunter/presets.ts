import {
	Consumes,
	Flask,
	Food,
	Spec,
} from '../core/proto/common.js';
import { SavedTalents } from '../core/proto/ui.js';

import {
	Hunter_Rotation as HunterRotation,
	Hunter_Rotation_RotationType as RotationType,
	Hunter_Rotation_StingType as StingType,
	Hunter_Options as HunterOptions,
	Hunter_Options_Ammo as Ammo,
	Hunter_Options_PetType as PetType,
} from '../core/proto/hunter.js';

import * as PresetUtils from '../core/preset_utils.js';

import BlankGear from './gear_sets/blank.gear.json';

import BmApl from './apls/bm.apl.json';
import MmApl from './apls/mm.apl.json';
import MmAdvApl from './apls/mm_advanced.apl.json';
import SvApl from './apls/sv.apl.json';
import SvAdvApl from './apls/sv_advanced.apl.json';
import AoeApl from './apls/aoe.apl.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const GearBeastMasteryDefault = PresetUtils.makePresetGear('Blank', BlankGear, { talentTree: 0 })
export const GearMarksmanDefault = PresetUtils.makePresetGear('Blank', BlankGear, { talentTree: 1 })
export const GearSurvivalDefault = PresetUtils.makePresetGear('Blank', BlankGear, { talentTree: 2 })

export const DefaultRotation = HunterRotation.create({
	type: RotationType.SingleTarget,
	sting: StingType.SerpentSting,
	trapWeave: true,
	viperStartManaPercent: 0.1,
	viperStopManaPercent: 0.3,
	multiDotSerpentSting: true,
	allowExplosiveShotDownrank: true,
});

export const ROTATION_PRESET_SIMPLE_DEFAULT = PresetUtils.makePresetSimpleRotation('Simple Default', Spec.SpecHunter, DefaultRotation);
export const ROTATION_PRESET_BM = PresetUtils.makePresetAPLRotation('BM', BmApl, { talentTree: 0 });
export const ROTATION_PRESET_MM = PresetUtils.makePresetAPLRotation('MM', MmApl, { talentTree: 1 });
export const ROTATION_PRESET_MM_ADVANCED = PresetUtils.makePresetAPLRotation('MM (Advanced)', MmAdvApl, { talentTree: 1 });
export const ROTATION_PRESET_SV = PresetUtils.makePresetAPLRotation('SV', SvApl, { talentTree: 2 });
export const ROTATION_PRESET_SV_ADVANCED = PresetUtils.makePresetAPLRotation('SV (Advanced)', SvAdvApl, { talentTree: 2 });
export const ROTATION_PRESET_AOE = PresetUtils.makePresetAPLRotation('AOE', AoeApl);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/classic/talent-calc and copy the numbers in the url.
export const BeastMasteryTalents = {
	name: 'Beast Mastery',
	data: SavedTalents.create({
		talentsString: '51200201505112243120531251-025305101',
	}),
};

export const MarksmanTalents = {
	name: 'Marksman',
	data: SavedTalents.create({
		talentsString: '502-025335101030013233135031051-5000032',
	}),
};

export const SurvivalTalents = {
	name: 'Survival',
	data: SavedTalents.create({
		talentsString: '-005305101-5000032500033330531135301331',
	}),
};

export const DefaultOptions = HunterOptions.create({
	ammo: Ammo.SaroniteRazorheads,
	useHuntersMark: true,
	petType: PetType.Wolf,
	petTalents: {},
	petUptime: 1,
	sniperTrainingUptime: 0.9,
	timeToTrapWeaveMs: 2000,
});

export const BMDefaultOptions = HunterOptions.create({
	ammo: Ammo.SaroniteRazorheads,
	useHuntersMark: true,
	petType: PetType.Wolf,
	petTalents: {},
	petUptime: 1,
	sniperTrainingUptime: 0.9,
	timeToTrapWeaveMs: 2000,
});

export const DefaultConsumes = Consumes.create({
	flask: Flask.FlaskUnknown,
	food: Food.FoodUnknown,
});
