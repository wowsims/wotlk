import { Consumes } from '../core/proto/common.js';
import { Food } from '../core/proto/common.js';
import { EquipmentSpec } from '../core/proto/common.js';
import { Potions } from '../core/proto/common.js';
import { Flask } from '../core/proto/common.js';
import { Glyphs } from '../core/proto/common.js';
import { SavedTalents } from '../core/proto/ui.js';

import { 
	FeralDruid_Rotation as FeralDruidRotation,
	FeralDruid_Options as FeralDruidOptions,
	DruidMajorGlyph,
	DruidMinorGlyph,
    FeralDruid_Rotation_BearweaveType,
 } from '../core/proto/druid.js';

import * as Tooltips from '../core/constants/tooltips.js';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/wotlk/talent-calc and copy the numbers in the url.
export const StandardTalents = {
	name: 'Standard',
	data: SavedTalents.create({
		talentsString: '-503202132322010053120230310511-205503012',
        glyphs: Glyphs.create({
            major1: DruidMajorGlyph.GlyphOfRip,
            major2: DruidMajorGlyph.GlyphOfSavageRoar,
            major3: DruidMajorGlyph.GlyphOfShred,
            minor1: DruidMinorGlyph.GlyphOfDash,
            minor2: DruidMinorGlyph.GlyphOfTheWild,
            minor3: DruidMinorGlyph.GlyphOfUnburdenedRebirth,
        }),
	}),
};

export const DefaultRotation = FeralDruidRotation.create({
	bearWeaveType: FeralDruid_Rotation_BearweaveType.None,
  minCombosForRip: 5,
  minCombosForBite: 5,

  useRake: true,
  useBite: false,
  mangleSpam: false,
  biteTime: 10.0,
  berserkBiteThresh: 30.0,
  powerbear: false,
  maxRoarClip: 10.0,
	maintainFaerieFire: true,
});

export const DefaultOptions = FeralDruidOptions.create({
  latencyMs: 100
});

export const DefaultConsumes = Consumes.create({
	flask: Flask.FlaskOfEndlessRage,
	food: Food.FoodFishFeast,
	defaultPotion: Potions.PotionOfSpeed,
});

export const PreRaid_PRESET = {
	name: 'PreRaid',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	gear: EquipmentSpec.fromJsonString(`{"items": [
    {
      "id": 42550,
      "enchant": 44879,
      "gems": [
        41398,
        39996
      ]
    },
    {
      "id": 40678
    },
    {
      "id": 37139,
      "enchant": 44871,
      "gems": [
        39996
      ]
    },
    {
      "id": 37840,
      "enchant": 55002
    },
    {
      "id": 37219,
      "enchant": 44489
    },
    {
      "id": 44203,
      "enchant": 44484,
      "gems": [
        0
      ]
    },
    {
      "id": 37409,
      "enchant": 54999,
      "gems": [
        0
      ]
    },
    {
      "id": 40694,
      "gems": [
        49110,
        39996
      ]
    },
    {
      "id": 37644,
      "enchant": 38374
    },
    {
      "id": 44297,
      "enchant": 55016
    },
    {
      "id": 37642
    },
    {
      "id": 37624
    },
    {
      "id": 40684
    },
    {
      "id": 37166
    },
    {
      "id": 37883,
      "enchant": 44483
    },
    {},
    {
      "id": 40713
    }
  ]}`),
};

export const P1_PRESET = {
	name: 'P1 Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	gear: EquipmentSpec.fromJsonString(`{"items": [
        {
          "id": 40473,
          "enchant": 44879,
          "gems": [
            41398,
            39996
          ]
        },
        {
          "id": 44664,
          "gems": [
            39996
          ]
        },
        {
          "id": 40494,
          "enchant": 44871,
          "gems": [
            39996
          ]
        },
        {
          "id": 40403,
          "enchant": 55002
        },
        {
          "id": 40539,
          "enchant": 44489,
          "gems": [
            39996
          ]
        },
        {
          "id": 39765,
          "enchant": 44484,
          "gems": [
            39996,
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
          "id": 40205,
          "gems": [
            39996
          ]
        },
        {
          "id": 44011,
          "enchant": 38374,
          "gems": [
            39996,
            49110
          ]
        },
        {
          "id": 40243,
          "enchant": 55016,
          "gems": [
            40014
          ]
        },
        {
          "id": 40474
        },
        {
          "id": 40717
        },
        {
          "id": 42987
        },
        {
          "id": 40256
        },
        {
          "id": 40388,
          "enchant": 44492
        },
        {},
        {
          "id": 39757
        }
      ]}`),
};