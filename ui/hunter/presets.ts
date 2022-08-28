import { CustomRotation, CustomSpell } from '../core/proto/common.js';
import { Consumes } from '../core/proto/common.js';
import { EquipmentSpec } from '../core/proto/common.js';
import { Flask } from '../core/proto/common.js';
import { Food } from '../core/proto/common.js';
import { Glyphs } from '../core/proto/common.js';
import { PetFood } from '../core/proto/common.js';
import { Potions } from '../core/proto/common.js';
import { SavedTalents } from '../core/proto/ui.js';
import { ferocityDefault, ferocityBMDefault } from '../core/talents/hunter_pet.js';

import {
	Hunter_Rotation as HunterRotation,
	Hunter_Rotation_RotationType as RotationType,
	Hunter_Rotation_StingType as StingType,
	Hunter_Rotation_SpellOption as SpellOption,
	Hunter_Options as HunterOptions,
	Hunter_Options_Ammo as Ammo,
	Hunter_Options_PetType as PetType,
	HunterMajorGlyph as MajorGlyph,
	HunterMinorGlyph as MinorGlyph,
} from '../core/proto/hunter.js';

import * as Tooltips from '../core/constants/tooltips.js';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/wotlk/talent-calc and copy the numbers in the url.
export const BeastMasteryTalents = {
	name: 'Beast Mastery',
	data: SavedTalents.create({
		talentsString: '51200201505112243120531251-00530513',
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
		talentsString: '502-035305101230013233135031351-5000002',
		glyphs: Glyphs.create({
			major1: MajorGlyph.GlyphOfSerpentSting,
			major2: MajorGlyph.GlyphOfSteadyShot,
			major3: MajorGlyph.GlyphOfKillShot,
			minor1: MinorGlyph.GlyphOfFeignDeath,
			minor2: MinorGlyph.GlyphOfRevivePet,
			minor3: MinorGlyph.GlyphOfMendPet,
		}),
	}),
};

export const SurvivalTalents = {
	name: 'Survival',
	data: SavedTalents.create({
		talentsString: '-025305101-5000032500033330522135301311',
		glyphs: Glyphs.create({
			major1: MajorGlyph.GlyphOfSerpentSting,
			major2: MajorGlyph.GlyphOfExplosiveShot,
			major3: MajorGlyph.GlyphOfKillShot,
			minor1: MinorGlyph.GlyphOfFeignDeath,
			minor2: MinorGlyph.GlyphOfRevivePet,
			minor3: MinorGlyph.GlyphOfMendPet,
		}),
	}),
};

export const DefaultRotation = HunterRotation.create({
	type: RotationType.SingleTarget,
	sting: StingType.SerpentSting,
	trapWeave: false,
	timeToTrapWeaveMs: 2000,
	viperStartManaPercent: 0.1,
	viperStopManaPercent: 0.3,
	customRotation: CustomRotation.create({
		spells: [
			CustomSpell.create({ spell: SpellOption.SerpentStingSpell }),
			CustomSpell.create({ spell: SpellOption.KillShot }),
			CustomSpell.create({ spell: SpellOption.ChimeraShot }),
			CustomSpell.create({ spell: SpellOption.BlackArrow }),
			CustomSpell.create({ spell: SpellOption.ExplosiveShot }),
			CustomSpell.create({ spell: SpellOption.AimedShot }),
			CustomSpell.create({ spell: SpellOption.ArcaneShot }),
			CustomSpell.create({ spell: SpellOption.SteadyShot }),
		],
	}),
});

export const DefaultOptions = HunterOptions.create({
	ammo: Ammo.SaroniteRazorheads,
	useHuntersMark: true,
	petType: PetType.Wolf,
	petTalents: ferocityDefault,
	petUptime: 1,
	sniperTrainingUptime: 0.9,
});

export const BMDefaultOptions = HunterOptions.create({
	ammo: Ammo.SaroniteRazorheads,
	useHuntersMark: true,
	petType: PetType.Wolf,
	petTalents: ferocityBMDefault,
	petUptime: 1,
	sniperTrainingUptime: 0.9,
});

export const DefaultConsumes = Consumes.create({
	defaultPotion: Potions.PotionOfSpeed,
	flask: Flask.FlaskOfEndlessRage,
	food: Food.FoodFishFeast,
	petFood: PetFood.PetFoodSpicedMammothTreats,
});

export const PRERAID_PRESET = {
	name: 'Preraid Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	//enableWhen: (player: Player<any>) => player.getTalentTree() == 2,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 42551,
			"enchant": 44879,
			"gems": [
				41398,
				40044
			]
		},
		{
			"id": 42645,
			"gems": [
				42143
			]
		},
		{
			"id": 37679,
			"enchant": 44871,
			"gems": [
				42143
			]
		},
		{
			"id": 43566,
			"enchant": 55002
		},
		{
			"id": 37144,
			"enchant": 44489,
			"gems": [
				40088
			]
		},
		{
			"id": 37170,
			"enchant": 60616,
			"gems": [
				0
			]
		},
		{
			"id": 37886,
			"enchant": 54999,
			"gems": [
				0
			]
		},
		{
			"id": 37407,
			"gems": [
				42143
			]
		},
		{
			"id": 37669,
			"enchant": 38374
		},
		{
			"id": 37167,
			"enchant": 55016,
			"gems": [
				39997,
				39997
			]
		},
		{
			"id": 37685
		},
		{
			"id": 42642,
			"gems": [
				39997
			]
		},
		{
			"id": 40684
		},
		{
			"id": 44253
		},
		{
			"id": 44249,
			"enchant": 44483
		},
		{},
		{
			"id": 43284,
			"enchant": 41167
		}
	]}`),
};

export const P1_PRESET = {
	name: 'P1 Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	//enableWhen: (player: Player<any>) => player.getTalentTree() == 2,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 40505,
			"enchant": 44879,
			"gems": [
				41398,
				40088
			]
		},
		{
			"id": 44664,
			"gems": [
				42143
			]
		},
		{
			"id": 40507,
			"enchant": 44871,
			"gems": [
				42143
			]
		},
		{
			"id": 40403,
			"enchant": 55002
		},
		{
			"id": 43998,
			"enchant": 44489,
			"gems": [
				42143,
				39997
			]
		},
		{
			"id": 40282,
			"enchant": 60616,
			"gems": [
				39997,
				0
			]
		},
		{
			"id": 40541,
			"enchant": 54999,
			"gems": [
				0
			]
		},
		{
			"id": 39762,
			"gems": [
				39997
			]
		},
		{
			"id": 40331,
			"enchant": 38374,
			"gems": [
				39997,
				39997
			]
		},
		{
			"id": 40549,
			"enchant": 55016
		},
		{
			"id": 40074
		},
		{
			"id": 40474
		},
		{
			"id": 40431
		},
		{
			"id": 44253
		},
		{
			"id": 40388,
			"enchant": 44483
		},
		{},
		{
			"id": 40385,
			"enchant": 41167
		}
	]}`),
};
