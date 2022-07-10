import { Consumes } from '/wotlk/core/proto/common.js';

import { EquipmentSpec } from '/wotlk/core/proto/common.js';
import { Flask } from '/wotlk/core/proto/common.js';
import { Food } from '/wotlk/core/proto/common.js';
import { ItemSpec } from '/wotlk/core/proto/common.js';
import { Potions } from '/wotlk/core/proto/common.js';
import { Faction } from '/wotlk/core/proto/common.js';
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
	data: '0532001523212351322301351-005052031',
};

export const RestoTalents = {
	name: 'Resto',
	data: '5003--55035051355310510321',
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
	defaultPotion: Potions.SuperManaPotion,
	flask: Flask.FlaskOfBlindingLight,
	food: Food.FoodBlackenedBasilisk,
	mainHandImbue: WeaponImbue.WeaponImbueBrilliantWizardOil,
});

export const PRE_RAID_PRESET = {
	name: 'Pre-raid Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	gear: EquipmentSpec.fromJsonString(`{"items": [
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
        "id": 39592,
        "enchant": 44623,
        "gems": [
          39998,
          40025
        ]
      },
      {
        "id": 37361,
        "enchant": 44498,
        "gems": [
          0
        ]
      },
      {
        "id": 42113,
        "enchant": 44488,
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
        "id": 37695,
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
        "id": 37192,
        "enchant": 22536
      },
      {
        "id": 37694,
        "enchant": 22536
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
