import { Consumes, PetFood } from '../core/proto/common.js';
import { EquipmentSpec } from '../core/proto/common.js';
import { Flask } from '../core/proto/common.js';
import { Food } from '../core/proto/common.js';
import { Glyphs } from '../core/proto/common.js';
import { Potions } from '../core/proto/common.js';
import { SavedTalents } from '../core/proto/ui.js';

import {
	Deathknight_Rotation as DeathKnightRotation,
	Deathknight_Options as DeathKnightOptions,
	DeathknightMajorGlyph,
	DeathknightMinorGlyph,
} from '../core/proto/deathknight.js';

import * as Tooltips from '../core/constants/tooltips.js';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wotlk.wowhead.com/talent-calc and copy the numbers in the url.
export const FrostTalents = {
	name: 'Frost',
	data: SavedTalents.create({
		talentsString: '23050005-32005350352203012300033101351',
		glyphs: Glyphs.create({
			major1: DeathknightMajorGlyph.GlyphOfObliterate,
			major2: DeathknightMajorGlyph.GlyphOfFrostStrike,
			major3: DeathknightMajorGlyph.GlyphOfDisease,
			minor1: DeathknightMinorGlyph.GlyphOfHornOfWinter,
			minor2: DeathknightMinorGlyph.GlyphOfPestilence,
			minor3: DeathknightMinorGlyph.GlyphOfRaiseDead,
		}),
	}),
};

export const FrostUnholyTalents = {
	name: 'Frost Unholy',
	data: SavedTalents.create({
		talentsString: '01-32002350351203012300033101351-230200305003',
		glyphs: Glyphs.create({
			major1: DeathknightMajorGlyph.GlyphOfObliterate,
			major2: DeathknightMajorGlyph.GlyphOfFrostStrike,
			major3: DeathknightMajorGlyph.GlyphOfDisease,
			minor1: DeathknightMinorGlyph.GlyphOfHornOfWinter,
			minor2: DeathknightMinorGlyph.GlyphOfPestilence,
			minor3: DeathknightMinorGlyph.GlyphOfRaiseDead,
		}),
	}),
};

export const UnholyDualWieldTalents = {
	name: 'Unholy Dual Wield',
	data: SavedTalents.create({
		talentsString: '-320033500002-2300303050032152000150013133151',
		glyphs: Glyphs.create({
			major1: DeathknightMajorGlyph.GlyphOfTheGhoul,
			major2: DeathknightMajorGlyph.GlyphOfDarkDeath,
			major3: DeathknightMajorGlyph.GlyphOfDeathAndDecay,
			minor1: DeathknightMinorGlyph.GlyphOfHornOfWinter,
			minor2: DeathknightMinorGlyph.GlyphOfPestilence,
			minor3: DeathknightMinorGlyph.GlyphOfRaiseDead,
		}),
	}),
};


export const DefaultRotation = DeathKnightRotation.create({
  useDeathAndDecay: true,
  btGhoulFrenzy: true,
	refreshHornOfWinter: false,
});

export const DefaultOptions = DeathKnightOptions.create({
	startingRunicPower: 0,
	petUptime: 1,
	precastGhoulFrenzy: true,
  precastHornOfWinter: true,
});

export const DefaultConsumes = Consumes.create({
	flask: Flask.FlaskOfEndlessRage,
	food: Food.FoodDragonfinFilet,
	defaultPotion: Potions.PotionOfSpeed,
	petFood: PetFood.PetFoodKiblersBits,
	prepopPotion:  Potions.PotionOfSpeed,
});

export const P1_UNHOLY_DW_BIS_PRESET = {
	name: 'P1 Unholy DW BiS',
	toolbar: Tooltips.BASIC_BIS_DISCLAIMER,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
      "id": 44006,
      "enchant": 44879,
      "gems": [
        41400,
        49110
      ]
    },
    {
      "id": 39421
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
      "id": 40278,
      "gems": [
        42142,
        42142
      ]
    },
    {
      "id": 40556,
      "enchant": 38374,
      "gems": [
        39996,
        39996
      ]
    },
    {
      "id": 40591,
      "enchant": 55016
    },
    {
      "id": 40717
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
      "enchant": 53344
    },
    {
      "id": 40491,
      "enchant": 44495
    },
    {
      "id": 40867
    }
	]}`),
};

export const P1_FROST_PRE_BIS_PRESET = {
	name: 'P1 Frost Pre-Raid',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	gear: EquipmentSpec.fromJsonString(`{ "items": [
    {
      "id": 41386,
      "enchant": 44879,
      "gems": [
        41398,
        34143
      ]
    },
    {
      "id": 37397
    },
    {
      "id": 37593,
      "enchant": 44871
    },
    {
      "id": 37647,
      "enchant": 44472
    },
    {
      "id": 39617,
      "enchant": 44623,
      "gems": [
        42142,
        39996
      ]
    },
    {
      "id": 41355,
      "enchant": 60616,
      "gems": [
        0
      ]
    },
    {
      "id": 39618,
      "enchant": 54999,
      "gems": [
        39996,
        0
      ]
    },
    {
      "id": 37194,
      "gems": [
        39996,
        42142
      ]
    },
    {
      "id": 37193,
      "enchant": 38374,
      "gems": [
        42142,
        39996
      ]
    },
    {
      "id": 44306,
      "enchant": 55016,
      "gems": [
        39996,
        39996
      ]
    },
    {
      "id": 37642
    },
    {
      "id": 37151
    },
    {
      "id": 40684
    },
    {
      "id": 42987
    },
    {
      "id": 44250,
      "enchant": 53343
    },
    {
      "id": 44250,
      "enchant": 53344
    },
    {
      "id": 40715
    }
  ]}`),
};

export const P1_FROST_BIS_PRESET = {
	name: 'P1 Frost BiS',
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


export const P1_FROST_HITCAP_PRESET = {
  name: 'P1 Frost Hitcap',
  tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
  gear: EquipmentSpec.fromJsonString(`{"items": [
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
      "gems": [
        39996,
        0
      ]
    },
    {
      "id": 40278,
      "gems": [
        39996,
        42142
      ]
    },
    {
      "id": 43994,
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
      "id": 43993,
      "gems": [
        39996
      ]
    },
    {
      "id": 40075
    },
    {
      "id": 40256
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
}