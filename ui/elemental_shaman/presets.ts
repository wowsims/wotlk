import { Consumes } from '/wotlk/core/proto/common.js';

import { EquipmentSpec } from '/wotlk/core/proto/common.js';
import { Flask } from '/wotlk/core/proto/common.js';
import { Food } from '/wotlk/core/proto/common.js';
import { Glyphs } from '/wotlk/core/proto/common.js';
import { ItemSpec } from '/wotlk/core/proto/common.js';
import { Potions } from '/wotlk/core/proto/common.js';
import { Faction } from '/wotlk/core/proto/common.js';
import { SavedTalents } from '/wotlk/core/proto/ui.js';
import { WeaponImbue } from '/wotlk/core/proto/common.js';
import { Player } from '/wotlk/core/player.js';

import { ElementalShaman, ElementalShaman_Rotation as ElementalShamanRotation, ElementalShaman_Options as ElementalShamanOptions, ShamanShield } from '/wotlk/core/proto/shaman.js';
import { ElementalShaman_Rotation_RotationType as RotationType } from '/wotlk/core/proto/shaman.js';

import {
	AirTotem,
	EarthTotem,
	FireTotem,
	WaterTotem,
	ShamanTotems,
} from '/wotlk/core/proto/shaman.js';


import * as Enchants from '/wotlk/core/constants/enchants.js';
import * as Gems from '/wotlk/core/proto_utils/gems.js';
import * as Tooltips from '/wotlk/core/constants/tooltips.js';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/wotlk/talent-calc and copy the numbers in the url.
export const StandardTalents = {
	name: 'Standard',
	data: SavedTalents.create({
		talentsString: '0532001523212351322301351-005052031',
	}),
};

export const RestoTalents = {
	name: 'Resto',
	data: SavedTalents.create({
		talentsString: '5003--55035051355310510321',
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
	mainHandImbue: WeaponImbue.WeaponImbueShamanFlametongue,
});

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
          39998
        ]
      },
      {
        "id": 42647,
        "gems": [
          39998
        ]
      },
      {
        "id": 37398,
        "enchant": 44874
      },
      {
        "id": 41610,
        "enchant": 44472
      },
      {
        "id": 43410,
        "enchant": 44623,
        "gems": [
          39998,
          40014
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
        "enchant": 54999,
        "gems": [
          0
        ]
      },
      {
        "id": 40696,
        "gems": [
          40051,
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
          40025
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
          40025
        ]
      },
      {
        "id": 44661,
        "gems": [
          40027
        ]
      },
      {
        "id": 40518,
        "enchant": 44874,
        "gems": [
          39998
        ]
      },
      {
        "id": 44005,
        "enchant": 44472,
        "gems": [
          40025
        ]
      },
      {
        "id": 40514,
        "enchant": 44623,
        "gems": [
          39998,
          40025
        ]
      },
      {
        "id": 40324,
        "enchant": 44498,
        "gems": [
          40025,
          0
        ]
      },
      {
        "id": 40515,
        "enchant": 54999,
        "gems": [
          40051,
          0
        ]
      },
      {
        "id": 40327,
        "gems": [
          39998
        ]
      },
      {
        "id": 40517,
        "enchant": 41602,
        "gems": [
          40051,
          40027
        ]
      },
      {
        "id": 40237,
        "enchant": 60623,
        "gems": [
          40025
        ]
      },
      {
        "id": 40399
      },
      {
        "id": 48957
      },
      {
        "id": 40255
      },
      {
        "id": 39229
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
        "id": 40708
      }
    ]
  }`),
};
