import { Consumes } from '../core/proto/common.js';

import { EquipmentSpec } from '../core/proto/common.js';
import { Flask } from '../core/proto/common.js';
import { Food } from '../core/proto/common.js';
import { Glyphs } from '../core/proto/common.js';
import { Potions } from '../core/proto/common.js';
import { SavedTalents } from '../core/proto/ui.js';

import {
	RestorationShaman_Rotation as RestorationShamanRotation,
	RestorationShaman_Options as RestorationShamanOptions,
	ShamanShield,
	ShamanMajorGlyph,
	ShamanMinorGlyph,
	ShamanHealSpell,
} from '../core/proto/shaman.js';

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
export const TankHealingTalents = {
	name: 'Tank Healing',
	data: SavedTalents.create({
		talentsString: '-30205033-05005331335010501122331251',
		glyphs: Glyphs.create({
			major1: ShamanMajorGlyph.GlyphOfEarthlivingWeapon,
			major2: ShamanMajorGlyph.GlyphOfEarthShield,
			major3: ShamanMajorGlyph.GlyphOfLesserHealingWave,
			minor2: ShamanMinorGlyph.GlyphOfWaterShield,
			minor1: ShamanMinorGlyph.GlyphOfRenewedLife,
			minor3: ShamanMinorGlyph.GlyphOfGhostWolf,
		}),
	}),
};
export const RaidHealingTalents = {
	name: 'Raid Healing',
	data: SavedTalents.create({
		talentsString: '-3020503-50005331335310501122331251',
		glyphs: Glyphs.create({
			major1: ShamanMajorGlyph.GlyphOfChainHeal,
			major2: ShamanMajorGlyph.GlyphOfEarthShield,
			major3: ShamanMajorGlyph.GlyphOfEarthlivingWeapon,
			minor2: ShamanMinorGlyph.GlyphOfWaterShield,
			minor1: ShamanMinorGlyph.GlyphOfRenewedLife,
			minor3: ShamanMinorGlyph.GlyphOfGhostWolf,
		}),
	}),
};

export const DefaultRotation = RestorationShamanRotation.create({
	totems: ShamanTotems.create({
		earth: EarthTotem.StrengthOfEarthTotem,
		air: AirTotem.WrathOfAirTotem,
		fire: FireTotem.FlametongueTotem,
		water: WaterTotem.HealingStreamTotem,
	}),
	useEarthShield: true,
	useRiptide: true,
});

export const DefaultOptions = RestorationShamanOptions.create({
	shield: ShamanShield.WaterShield,
	bloodlust: true,
	earthShieldPPM: 0,
});

export const DefaultConsumes = Consumes.create({
	defaultPotion: Potions.RunicManaInjector,
	flask: Flask.FlaskOfTheFrostWyrm,
	food: Food.FoodFishFeast,
});

export const PRE_RAID_PRESET = {
	name: 'Pre-raid Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 42555,
			"enchant": 3820,
			"gems": [
				41401,
				40017
			]
		},
		{
			"id": 40681
		},
		{
			"id": 37875,
			"enchant": 3838,
			"gems": [
				40017
			]
		},
		{
			"id": 37291,
			"enchant": 3859
		},
		{
			"id": 44180,
			"enchant": 2381
		},
		{
			"id": 37788,
			"enchant": 3758,
			"gems": [
				0
			]
		},
		{
			"id": 37623,
			"enchant": 3604,
			"gems": [
				0
			]
		},
		{
			"id": 40693,
			"gems": [
				40051,
				0
			]
		},
		{
			"id": 37791,
			"enchant": 3721
		},
		{
			"id": 44202,
			"enchant": 3606,
			"gems": [
				40105
			]
		},
		{
			"id": 44283
		},
		{
			"id": 37694
		},
		{
			"id": 37111
		},
		{
			"id": 37657
		},
		{
			"id": 37169,
			"enchant": 3834
		},
		{
			"id": 40700,
			"enchant": 1128
		},
		{
			"id": 40709
		}
  ]}`),
};

export const P1_PRESET = {
	name: 'P1 Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	gear: EquipmentSpec.fromJsonString(`{"items": [
        {
          "id": 40510,
          "enchant": 3820,
          "gems": [
            41401,
            39998
          ]
        },
        {
          "id": 44662,
          "gems": [
            40051
          ]
        },
        {
          "id": 40513,
          "enchant": 3810,
          "gems": [
            39998
          ]
        },
        {
          "id": 44005,
          "enchant": 3831,
          "gems": [
            40027
          ]
        },
        {
          "id": 40508,
          "enchant": 2381,
          "gems": [
            39998,
            40051
          ]
        },
        {
          "id": 40209,
          "enchant": 2332,
          "gems": [
            0
          ]
        },
        {
          "id": 40564,
          "enchant": 3246,
          "gems": [
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
          "id": 40512,
          "enchant": 3721,
          "gems": [
            39998,
            40027
          ]
        },
        {
          "id": 39734,
          "enchant": 3244
        },
        {
          "id": 40399
        },
        {
          "id": 40375
        },
        {
          "id": 37111
        },
        {
          "id": 40685
        },
        {
          "id": 40395,
          "enchant": 3834
        },
        {
          "id": 40401,
          "enchant": 1128
        },
        {
          "id": 40709
        }
      ]}`),
};

export const P2_PRESET = {
	name: 'P2 Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	gear: EquipmentSpec.fromJsonString(`{"items": [
        {
          "id": 46201,
          "enchant": 3820,
          "gems": [
            41401,
            40027
          ]
        },
        {
          "id": 45443,
          "gems": [
            40027
          ]
        },
        {
          "id": 46204,
          "enchant": 3810,
          "gems": [
            45883
          ]
        },
        {
          "id": 45486,
          "enchant": 3831,
          "gems": [
            40051
          ]
        },
        {
          "id": 45867,
          "enchant": 2381,
          "gems": [
            40051,
            39998
          ]
        },
        {
          "id": 45460,
          "enchant": 2332,
          "gems": [
            40027,
            0
          ]
        },
        {
          "id": 46199,
          "enchant": 3246,
          "gems": [
            40051,
            0
          ]
        },
        {
          "id": 45151,
          "gems": [
            39998
          ]
        },
        {
          "id": 46202,
          "enchant": 3721,
          "gems": [
            39998,
            40027
          ]
        },
        {
          "id": 45615,
          "enchant": 3232,
          "gems": [
            39998,
            40027
          ]
        },
        {
          "id": 45614,
          "gems": [
            40051
          ]
        },
        {
          "id": 46046,
          "gems": [
            40051
          ]
        },
        {
          "id": 45535
        },
        {
          "id": 45466
        },
        {
          "id": 46017,
          "enchant": 3834
        },
        {
          "id": 45470,
          "enchant": 1128,
          "gems": [
            40027
          ]
        },
        {
          "id": 45114
        }
      ]}`),
};
