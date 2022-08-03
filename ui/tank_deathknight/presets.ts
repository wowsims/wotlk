import { Consumes } from '../core/proto/common.js';
import { EquipmentSpec } from '../core/proto/common.js';
import { Flask } from '../core/proto/common.js';
import { Food } from '../core/proto/common.js';
import { Glyphs } from '../core/proto/common.js';
import { Potions } from '../core/proto/common.js';
import { SavedTalents } from '../core/proto/ui.js';

import {
	TankDeathknight_Rotation as TankDeathKnightRotation,
	TankDeathknight_Options as TankDeathKnightOptions,
	DeathknightMajorGlyph,
	DeathknightMinorGlyph,
} from '../core/proto/deathknight.js';

import * as Tooltips from '../core/constants/tooltips.js';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wotlk.wowhead.com/talent-calc and copy the numbers in the url.
export const BloodTalents = {
	name: 'Frost',
	data: SavedTalents.create({
		talentsString: '23050005-32005350352203012300033101351',
		glyphs: Glyphs.create({
			major1: DeathknightMajorGlyph.GlyphOfObliterate,
			major2: DeathknightMajorGlyph.GlyphOfFrostStrike,
			major3: DeathknightMajorGlyph.GlyphOfDisease,
			minor1: DeathknightMinorGlyph.GlyphOfHornOfWinter,
			minor2: DeathknightMinorGlyph.GlyphOfBloodTap,
			minor3: DeathknightMinorGlyph.GlyphOfRaiseDead,
		}),
	}),
};


export const DefaultRotation = TankDeathKnightRotation.create({
});

export const DefaultOptions = TankDeathKnightOptions.create({
	startingRunicPower: 0,
});

export const DefaultConsumes = Consumes.create({
	flask: Flask.FlaskOfFortification,
	food: Food.FoodDragonfinFilet,
	defaultPotion: Potions.PotionOfSpeed,
	prepopPotion:  Potions.IndestructiblePotion,
});

export const P1_BLOOD_BIS_PRESET = {
	name: 'P1 Blood BiS Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	gear: EquipmentSpec.fromJsonString(`{ "items": [
    {
      "id": 44006,
      "enchant": 44879,
      "gems": [
        41398,
        34143
      ]
    },
    {
      "id": 44664,
      "gems": [
        39996
      ]
    },
    {
      "id": 40557,
      "enchant": 44871,
      "gems": [
        39996
      ]
    },
    {
      "id": 40403,
      "enchant": 44472
    },
    {
      "id": 40550,
      "enchant": 44623,
      "gems": [
        42142,
        39996
      ]
    },
    {
      "id": 40330,
      "enchant": 60616,
      "gems": [
        39996,
        0
      ]
    },
    {
      "id": 40552,
      "enchant": 54999,
      "gems": [
        39996,
        0
      ]
    },
    {
      "id": 40317,
      "gems": [
        42142
      ]
    },
    {
      "id": 40556,
      "enchant": 38374,
      "gems": [
        42142,
        39996
      ]
    },
    {
      "id": 40591,
      "enchant": 55016
    },
    {
      "id": 39401
    },
    {
      "id": 40075
    },
    {
      "id": 40684
    },
    {
      "id": 42987
    },
    {
      "id": 40189,
      "enchant": 53343
    },
    {
      "id": 40189,
      "enchant": 53344
    },
    {
      "id": 40207
    }
  ]}`),
};
