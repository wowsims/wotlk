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
  fnMinManaPer: 66,
  clMinManaPer: 33,
  useChainLightning: false,
  useFireNova: false,
  useThunderstorm: true,
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
	name: 'SWP Pre-Raid Preset',
	tooltip: "Pre-raid preset that includes some gear from SWP.",
	gear: EquipmentSpec.fromJsonString(`{
    "items": [
      {
        "id": 34332,
        "enchant": 44877,
        "gems": [
          40049,
          41285
        ]
      },
      {
        "id": 37595
      },
      {
        "id": 31023,
        "enchant": 44874,
        "gems": [
          42144,
          40049
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
        "id": 34437,
        "enchant": 44498,
        "gems": [
          40049,
          0
        ]
      },
      {
        "id": 34344,
        "enchant": 44592,
        "gems": [
          40049,
          42144,
          0
        ]
      },
      {
        "id": 34542,
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
        "id": 34566,
        "enchant": 60623,
        "gems": [
          40049
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
        "id": 37180,
        "enchant": 44877,
        "gems": [
          41285,
          42144
        ]
      },
      {
        "id": 37595
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
        "id": 39593,
        "enchant": 44592,
        "gems": [
          40051,
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
