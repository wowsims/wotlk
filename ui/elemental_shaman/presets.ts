import { Consumes } from '../core/proto/common.js';

import { EquipmentSpec } from '../core/proto/common.js';
import { Flask } from '../core/proto/common.js';
import { Food } from '../core/proto/common.js';
import { Glyphs } from '../core/proto/common.js';
import { Potions } from '../core/proto/common.js';
import { SavedTalents } from '../core/proto/ui.js';

import { ElementalShaman_Rotation as ElementalShamanRotation, ElementalShaman_Options as ElementalShamanOptions, ShamanShield, ShamanMajorGlyph, ShamanMinorGlyph } from '../core/proto/shaman.js';
import { ElementalShaman_Rotation_RotationType as RotationType } from '../core/proto/shaman.js';

import {
	AirTotem,
	EarthTotem,
	FireTotem,
	WaterTotem,
	ShamanTotems,
} from '../core/proto/shaman.js';


import * as Tooltips from '../core/constants/tooltips.js';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/wotlk/talent-calc and copy the numbers in the url.
export const StandardTalents = {
	name: 'Standard',
	data: SavedTalents.create({
		talentsString: '0532001523212351322301351-005052031',
		glyphs: Glyphs.create({
			major1: ShamanMajorGlyph.GlyphOfLava,
			major2: ShamanMajorGlyph.GlyphOfTotemOfWrath,
			major3: ShamanMajorGlyph.GlyphOfLightningBolt,
			minor1: ShamanMinorGlyph.GlyphOfThunderstorm,
			minor2: ShamanMinorGlyph.GlyphOfWaterShield,
			minor3: ShamanMinorGlyph.GlyphOfGhostWolf,
		}),
	}),
};

export const DefaultRotation = ElementalShamanRotation.create({
	totems: ShamanTotems.create({
		earth: EarthTotem.StrengthOfEarthTotem,
		air: AirTotem.WrathOfAirTotem,
		fire: FireTotem.TotemOfWrath,
		water: WaterTotem.ManaSpringTotem,
	}),
	type: RotationType.Adaptive,
});

export const DefaultOptions = ElementalShamanOptions.create({
	shield: ShamanShield.WaterShield,
	bloodlust: true,
});

export const DefaultConsumes = Consumes.create({
	defaultPotion: Potions.RunicManaPotion,
	flask: Flask.FlaskOfTheFrostWyrm,
	food: Food.FoodFishFeast,
});


export const SWP_PRESET = {
	name: 'SWP Preset',
	tooltip: "The P5 preset from the TBC simulator but regemmed and adjusted for new hit cap.",
	gear: EquipmentSpec.fromJsonString(`{
    "items": [
      {
        "id": 34332,
        "enchant": 29191,
        "gems": [
          40014,
          34220
        ]
      },
      {
        "id": 34204
      },
      {
        "id": 31023,
        "enchant": 23545,
        "gems": [
          40025,
          40014
        ]
      },
      {
        "id": 34242,
        "enchant": 33150,
        "gems": [
          39998
        ]
      },
      {
        "id": 34396,
        "enchant": 24003,
        "gems": [
          39998,
          40014,
          40014
        ]
      },
      {
        "id": 34437,
        "enchant": 22534,
        "gems": [
          40014,
          0
        ]
      },
      {
        "id": 34350,
        "enchant": 28272,
        "gems": [
          39998,
          40025,
          0
        ]
      },
      {
        "id": 34542,
        "gems": [
          40014,
          39998
        ]
      },
      {
        "id": 34186,
        "enchant": 24274,
        "gems": [
          40049,
          39998,
          39998
        ]
      },
      {
        "id": 34566,
        "enchant": 35297,
        "gems": [
          40051
        ]
      },
      {
        "id": 34230,
        "enchant": 22536
      },
      {
        "id": 34362,
        "enchant": 22536
      },
      {
        "id": 34429
      },
      {
        "id": 32483
      },
      {
        "id": 34336,
        "enchant": 22555
      },
      {
        "id": 34179
      },
      {
        "id": 32330
      }
    ]
  }`),
};

export const PRE_RAID_PRESET = {
	name: 'Pre-raid Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	gear: EquipmentSpec.fromJsonString(`{
    "items": [
      {
        "id": 37592,
        "enchant": 44877,
        "gems": [
          41285,
          42144
        ]
      },
      {
        "id": 45155
      },
      {
        "id": 37673,
        "enchant": 44874,
        "gems": [
          42144
        ]
      },
      {
        "id": 41610,
        "enchant": 55642
      },
      {
        "id": 39592,
        "enchant": 44489,
        "gems": [
          42144,
          40025
        ]
      },
      {
        "id": 37788,
        "enchant": 44498,
        "gems": [
          0
        ]
      },
      {
        "id": 42113,
        "enchant": 44592,
        "gems": [
          0
        ]
      },
      {
        "id": 40696,
        "gems": [
          40049,
          39998
        ]
      },
      {
        "id": 37791,
        "enchant": 41602
      },
      {
        "id": 44202,
        "enchant": 60623,
        "gems": [
          39998
        ]
      },
      {
        "id": 43253,
        "gems": [
          40027
        ]
      },
      {
        "id": 37694
      },
      {
        "id": 40682
      },
      {
        "id": 37873
      },
      {
        "id": 41384,
        "enchant": 44487
      },
      {
        "id": 40698
      },
      {
        "id": 40708
      }
    ]
  }`),
};

export const P1_PRESET = {
	name: 'P1 Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	gear: EquipmentSpec.fromJsonString(`{
    "items": [
      {
        "id": 40516,
        "enchant": 44877,
        "gems": [
          41285,
          40027
        ]
      },
      {
        "id": 44661,
        "gems": [
          39998
        ]
      },
      {
        "id": 40286,
        "enchant": 44874
      },
      {
        "id": 44005,
        "enchant": 55642,
        "gems": [
          40027
        ]
      },
      {
        "id": 40514,
        "enchant": 44489,
        "gems": [
          42144,
          42144
        ]
      },
      {
        "id": 40324,
        "enchant": 44498,
        "gems": [
          42144,
          0
        ]
      },
      {
        "id": 40302,
        "enchant": 44592,
        "gems": [
          0
        ]
      },
      {
        "id": 40301,
        "gems": [
          40014
        ]
      },
      {
        "id": 40560,
        "enchant": 41604
      },
      {
        "id": 40519,
        "enchant": 60623
      },
      {
        "id": 37694
      },
      {
        "id": 40399
      },
      {
        "id": 40432
      },
      {
        "id": 40255
      },
      {
        "id": 40395,
        "enchant": 44487
      },
      {
        "id": 40401,
        "enchant": 60653
      },
      {
        "id": 40267
      }
    ]
  }`),
};
