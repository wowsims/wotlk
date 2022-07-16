import { Consumes, PetFood } from '/wotlk/core/proto/common.js';
import { EquipmentSpec } from '/wotlk/core/proto/common.js';
import { Flask } from '/wotlk/core/proto/common.js';
import { Food } from '/wotlk/core/proto/common.js';
import { ItemSpec } from '/wotlk/core/proto/common.js';
import { Potions } from '/wotlk/core/proto/common.js';
import { Spec } from '/wotlk/core/proto/common.js';
import { WeaponImbue } from '/wotlk/core/proto/common.js';
import { Faction } from '/wotlk/core/proto/common.js';
import { Player } from '/wotlk/core/player.js';

import {
	DeathKnightTalents as DeathKnightTalents,
	DeathKnight,
	DeathKnight_Rotation as DeathKnightRotation,
	DeathKnight_Options as DeathKnightOptions,
} from '/wotlk/core/proto/deathknight.js';

import * as Enchants from '/wotlk/core/constants/enchants.js';
import * as Gems from '/wotlk/core/proto_utils/gems.js';
import * as Tooltips from '/wotlk/core/constants/tooltips.js';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wotlk.wowhead.com/talent-calc and copy the numbers in the url.
export const FrostTalents = {
	name: 'Frost Dps',
	data: '23050005-32005350352203012300033101351',
};

export const UnholyDualWieldTalents = {
	name: 'Unholy Dual Wield Dps',
	data: '-320033500002-2302303050032150000150013133151',
};

export const DefaultRotation = DeathKnightRotation.create({
});

export const DefaultOptions = DeathKnightOptions.create({
	startingRunicPower: 0,
	petUptime: 1,
});

export const DefaultConsumes = Consumes.create({
	flask: Flask.FlaskOfEndlessRage,
	food: Food.FoodDragonfinFilet,
	defaultPotion: Potions.PotionOfSpeed,
	petFood: PetFood.PetFoodKiblersBits
});

export const P1_UNHOLY_DW_BIS_PRESET = {
	name: 'P1 Unholy DW Bis',
	toolbar: Tooltips.BASIC_BIS_DISCLAIMER,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
		"id": 44006,
		"enchant": 44879,
		"gems": [
			41400,
			22459
		]
		},
		{
		"id": 44664,
		"gems": [
			39996
		]
		},
		{
		"id": 40414,
		"enchant": 44871
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
			40038
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
		"id": 40347,
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
		"id": 40294,
		"enchant": 38374
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
		"id": 40431
		},
		{
		"id": 42987
		},
		{
		"id": 40189
		},
		{
		"id": 40491
		},
		{
		"id": 40715
		}
	]}`),
};

export const P1_FROST_PRE_BIS_PRESET = {
	name: 'P1 Frost Pre-Raid Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 41386,
			"enchant": 44879,
			"gems": [
				41398,
				40022
			]
		},
		{
			"id": 42645,
			"gems": [
				39996
			]
		},
		{
			"id": 34388,
			"enchant": 50335,
			"gems": [
				39996,
				39996
			]
		},
		{
			"id": 37647,
			"enchant": 47898
		},
		{
			"id": 39617,
			"enchant": 60692,
			"gems": [
				39996,
				39996
			]
		},
		{
			"id": 41355,
			"enchant": 44575
		},
		{
			"id": 39618,
			"enchant": 60668,
			"gems": [
				39996
			]
		},
		{
			"id": 40694,
			"gems": [
				39996
			]
		},
		{
			"id": 37193,
			"enchant": 38374,
			"gems": [
				39996,
				39996
			]
		},
		{
			"id": 44306,
			"enchant": 60763,
			"gems": [
				39996,
				39996
			]
		},
		{
			"id": 37642
		},
		{
			"id": 44935
		},
		{
			"id": 40684
		},
		{
			"id": 42987
		},
		{
			"id": 41383,
			"enchant": 53343
		},
		{
			"id": 43611,
			"enchant": 53344
		},
		{
			"id": 40715
		}
	]}`),
};

export const P1_FROST_BIS_PRESET = {
	name: 'P1 Frost BiS Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	gear: EquipmentSpec.fromJsonString(`{"items": [
        {
          "id": 44006,
          "enchant": 44879,
          "gems": [
            41398,
            40022
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
          "enchant": 55002
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
          "id": 40694,
          "gems": [
            39996,
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
          "id": 40431
        },
        {
          "id": 42987
        },
        {
          "id": 42208,
          "enchant": 53343
        },
        {
          "id": 42232,
          "enchant": 53344
        },
        {
          "id": 40715
        }
      ]}`),
};
