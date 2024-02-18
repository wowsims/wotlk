import {
	Consumes,
	Flask,
	Food,
	Glyphs,
	PetFood,
	Potions,
	Spec,
} from '../core/proto/common.js';
import { SavedTalents } from '../core/proto/ui.js';
import { ferocityDefault, ferocityBMDefault } from '../core/talents/hunter_pet.js';

import {
	Hunter_Rotation as HunterRotation,
	Hunter_Rotation_RotationType as RotationType,
	Hunter_Rotation_StingType as StingType,
	Hunter_Options as HunterOptions,
	Hunter_Options_Ammo as Ammo,
	Hunter_Options_PetType as PetType,
	HunterMajorGlyph as MajorGlyph,
	HunterMinorGlyph as MinorGlyph,
} from '../core/proto/hunter.js';

import * as PresetUtils from '../core/preset_utils.js';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

import PreraidMMGear from './gear_sets/preraid_mm.gear.json';
export const MM_PRERAID_PRESET = PresetUtils.makePresetGear('MM PreRaid Preset', PreraidMMGear, { talentTrees: [0, 1] });
import P1MMGear from './gear_sets/p1_mm.gear.json';
export const MM_P1_PRESET = PresetUtils.makePresetGear('MM P1 Preset', P1MMGear, { talentTrees: [0, 1] });
import P2MMGear from './gear_sets/p2_mm.gear.json';
export const MM_P2_PRESET = PresetUtils.makePresetGear('MM P2 Preset', P2MMGear, { talentTrees: [0, 1] });
import P3MMGear from './gear_sets/p3_mm.gear.json';
export const MM_P3_PRESET = PresetUtils.makePresetGear('MM P3 Preset', P3MMGear, { talentTrees: [0, 1] });
import P4MMGear from './gear_sets/p4_mm.gear.json';
export const MM_P4_PRESET = PresetUtils.makePresetGear('MM P4 Preset', P4MMGear, { talentTrees: [0, 1] });
import P5MMGear from './gear_sets/p5_mm.gear.json';
export const MM_P5_PRESET = PresetUtils.makePresetGear('MM P5 Preset', P5MMGear, { talentTrees: [0, 1] });
import PreraidSVGear from './gear_sets/preraid_sv.gear.json';
export const SV_PRERAID_PRESET = PresetUtils.makePresetGear('SV PreRaid Preset', PreraidSVGear, { talentTree: 2 });
import P1SVGear from './gear_sets/p1_sv.gear.json';
export const SV_P1_PRESET = PresetUtils.makePresetGear('SV P1 Preset', P1SVGear, { talentTree: 2 });
import P2SVGear from './gear_sets/p2_sv.gear.json';
export const SV_P2_PRESET = PresetUtils.makePresetGear('SV P2 Preset', P2SVGear, { talentTree: 2 });
import P3SVGear from './gear_sets/p3_sv.gear.json';
export const SV_P3_PRESET = PresetUtils.makePresetGear('SV P3 Preset', P3SVGear, { talentTree: 2 });
import P4SVGear from './gear_sets/p4_sv.gear.json';
export const SV_P4_PRESET = PresetUtils.makePresetGear('SV P4 Preset', P4SVGear, { talentTree: 2 });
import P5SVGear from './gear_sets/p5_sv.gear.json';
export const SV_P5_PRESET = PresetUtils.makePresetGear('SV P5 Preset', P5SVGear, { talentTree: 2 });

export const DefaultSimpleRotation = HunterRotation.create({
	type: RotationType.SingleTarget,
	sting: StingType.SerpentSting,
	trapWeave: true,
	viperStartManaPercent: 0.1,
	viperStopManaPercent: 0.3,
	multiDotSerpentSting: true,
	allowExplosiveShotDownrank: true,
});

export const ROTATION_PRESET_SIMPLE_DEFAULT = PresetUtils.makePresetSimpleRotation('Simple Default', Spec.SpecHunter, DefaultSimpleRotation);
import BmApl from './apls/bm.apl.json';
export const ROTATION_PRESET_BM = PresetUtils.makePresetAPLRotation('BM', BmApl, { talentTree: 0 });
import MmApl from './apls/mm.apl.json';
export const ROTATION_PRESET_MM = PresetUtils.makePresetAPLRotation('MM', MmApl, { talentTree: 1 });
import MmAdvApl from './apls/mm_advanced.apl.json';
export const ROTATION_PRESET_MM_ADVANCED = PresetUtils.makePresetAPLRotation('MM (Advanced)', MmAdvApl, { talentTree: 1 });
import SvApl from './apls/sv.apl.json';
export const ROTATION_PRESET_SV = PresetUtils.makePresetAPLRotation('SV', SvApl, { talentTree: 2 });
import SvAdvApl from './apls/sv_advanced.apl.json';
export const ROTATION_PRESET_SV_ADVANCED = PresetUtils.makePresetAPLRotation('SV (Advanced)', SvAdvApl, { talentTree: 2 });
import AoeApl from './apls/aoe.apl.json';
export const ROTATION_PRESET_AOE = PresetUtils.makePresetAPLRotation('AOE', AoeApl);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/wotlk/talent-calc and copy the numbers in the url.
export const BeastMasteryTalents = {
	name: 'Beast Mastery',
	data: SavedTalents.create({
		talentsString: '51200201505112243120531251-025305101',
		glyphs: Glyphs.create({
			major1: MajorGlyph.GlyphOfBestialWrath,
			major2: MajorGlyph.GlyphOfSteadyShot,
			major3: MajorGlyph.GlyphOfSerpentSting,
			minor1: MinorGlyph.GlyphOfFeignDeath,
			minor2: MinorGlyph.GlyphOfRevivePet,
			minor3: MinorGlyph.GlyphOfMendPet,
		}),
	}),
};

export const MarksmanTalents = {
	name: 'Marksman',
	data: SavedTalents.create({
		talentsString: '502-025335101030013233135031051-5000032',
		glyphs: Glyphs.create({
			major1: MajorGlyph.GlyphOfSerpentSting,
			major2: MajorGlyph.GlyphOfSteadyShot,
			major3: MajorGlyph.GlyphOfExplosiveTrap,
			minor1: MinorGlyph.GlyphOfFeignDeath,
			minor2: MinorGlyph.GlyphOfRevivePet,
			minor3: MinorGlyph.GlyphOfMendPet,
		}),
	}),
};

export const SurvivalTalents = {
	name: 'Survival',
	data: SavedTalents.create({
		talentsString: '-005305101-5000032500033330531135301331',
		glyphs: Glyphs.create({
			major1: MajorGlyph.GlyphOfSerpentSting,
			major2: MajorGlyph.GlyphOfExplosiveTrap,
			major3: MajorGlyph.GlyphOfKillShot,
			minor1: MinorGlyph.GlyphOfFeignDeath,
			minor2: MinorGlyph.GlyphOfRevivePet,
			minor3: MinorGlyph.GlyphOfMendPet,
		}),
	}),
};

export const DefaultOptions = HunterOptions.create({
	ammo: Ammo.SaroniteRazorheads,
	useHuntersMark: true,
	petType: PetType.Wolf,
	petTalents: ferocityDefault,
	petUptime: 1,
	sniperTrainingUptime: 0.9,
	timeToTrapWeaveMs: 2000,
});

export const BMDefaultOptions = HunterOptions.create({
	ammo: Ammo.SaroniteRazorheads,
	useHuntersMark: true,
	petType: PetType.Wolf,
	petTalents: ferocityBMDefault,
	petUptime: 1,
	sniperTrainingUptime: 0.9,
	timeToTrapWeaveMs: 2000,
});

export const DefaultConsumes = Consumes.create({
	defaultPotion: Potions.PotionOfSpeed,
	flask: Flask.FlaskOfEndlessRage,
	food: Food.FoodFishFeast,
	petFood: PetFood.PetFoodSpicedMammothTreats,
});
