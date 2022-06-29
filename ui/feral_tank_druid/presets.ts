import { Consumes } from '/tbc/core/proto/common.js';
import { BattleElixir } from '/tbc/core/proto/common.js';
import { GuardianElixir } from '/tbc/core/proto/common.js';
import { Food } from '/tbc/core/proto/common.js';
import { EquipmentSpec } from '/tbc/core/proto/common.js';
import { ItemSpec } from '/tbc/core/proto/common.js';
import { Potions } from '/tbc/core/proto/common.js';
import { Conjured } from '/tbc/core/proto/common.js';
import { RaidTarget } from '/tbc/core/proto/common.js';
import { WeaponImbue } from '/tbc/core/proto/common.js';
import { Faction } from '/tbc/core/proto/common.js';
import { NO_TARGET } from '/tbc/core/proto_utils/utils.js';
import { Player } from '/tbc/core/player.js';

import {
	DruidTalents as DruidTalents,
	FeralTankDruid,
	FeralTankDruid_Rotation as DruidRotation,
	FeralTankDruid_Rotation_Swipe as Swipe,
	FeralTankDruid_Options as DruidOptions
} from '/tbc/core/proto/druid.js';

import * as Enchants from '/tbc/core/constants/enchants.js';
import * as Gems from '/tbc/core/proto_utils/gems.js';
import * as Tooltips from '/tbc/core/constants/tooltips.js';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

// Default talents. Uses the wowhead calculator format, make the talents on
// https://tbc.wowhead.com/talent-calc and copy the numbers in the url.
export const StandardTalents = {
	name: 'Standard',
	data: '-503032132322105301251-05503301',
};

export const DemoRoarTalents = {
	name: 'DemoRoar',
	data: '-553032132322105301051-05503001',
};

export const DefaultRotation = DruidRotation.create({
	maulRageThreshold: 50,
	swipe: Swipe.SwipeWithEnoughAP,
	swipeApThreshold: 2700,
	maintainDemoralizingRoar: true,
	maintainFaerieFire: true,
});

export const DefaultOptions = DruidOptions.create({
	innervateTarget: RaidTarget.create({
		targetIndex: NO_TARGET,
	}),
	startingRage: 20,
});

export const DefaultConsumes = Consumes.create({
	battleElixir: BattleElixir.ElixirOfMajorAgility,
	guardianElixir: GuardianElixir.GiftOfArthas,
	food: Food.FoodGrilledMudfish,
	defaultPotion: Potions.IronshieldPotion,
	defaultConjured: Conjured.ConjuredFlameCap,
	scrollOfAgility: 5,
	scrollOfStrength: 5,
	scrollOfProtection: 5,
});

export const P1_PRESET = {
	name: 'P1 Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 29098,
			"enchant": 29192,
			"gems": [
				24067,
				32409
			]
		},
		{
			"id": 28509
		},
		{
			"id": 29100,
			"enchant": 28911,
			"gems": [
				24033,
				24033
			]
		},
		{
			"id": 28660,
			"enchant": 34004
		},
		{
			"id": 29096,
			"enchant": 24003,
			"gems": [
				24067,
				24055,
				24055
			]
		},
		{
			"id": 28978,
			"enchant": 22533,
			"gems": [
				24033
			]
		},
		{
			"id": 29097,
			"enchant": 33153
		},
		{
			"id": 28986
		},
		{
			"id": 29099,
			"enchant": 29536
		},
		{
			"id": 30674,
			"enchant": 35297
		},
		{
			"id": 29279,
			"enchant": 22538
		},
		{
			"id": 28792,
			"enchant": 22538
		},
		{
			"id": 28830
		},
		{
			"id": 23836
		},
		{
			"id": 28476,
			"enchant": 22556
		},
		{
			"id": 23198
		}
	]}`),
};

export const P2_PRESET = {
	name: 'P2 Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 30228,
			"enchant": 29192,
			"gems": [
				24055,
				32409
			]
		},
		{
			"id": 33066,
			"gems": [
				24061
			]
		},
		{
			"id": 30230,
			"enchant": 28888,
			"gems": [
				24033,
				24033
			]
		},
		{
			"id": 28660,
			"enchant": 34004
		},
		{
			"id": 30222,
			"enchant": 24003,
			"gems": [
				24033,
				24061,
				24055
			]
		},
		{
			"id": 32810,
			"enchant": 22533,
			"gems": [
				24033
			]
		},
		{
			"id": 30223,
			"enchant": 33153
		},
		{
			"id": 30106,
			"gems": [
				24055,
				24033
			]
		},
		{
			"id": 30229,
			"enchant": 29536,
			"gems": [
				24033
			]
		},
		{
			"id": 32790,
			"enchant": 35297
		},
		{
			"id": 29279,
			"enchant": 22538
		},
		{
			"id": 28792,
			"enchant": 22538
		},
		{
			"id": 28579
		},
		{
			"id": 32658
		},
		{
			"id": 30021,
			"enchant": 22556
		},
		{
			"id": 32387
		}
	]}`),
};

export const P3_PRESET = {
	name: 'P3 Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 31039,
			"enchant": 29186,
			"gems": [
				32212,
				32409
			]
		},
		{
			"id": 32362
		},
		{
			"id": 31048,
			"enchant": 28889,
			"gems": [
				32220,
				32212
			]
		},
		{
			"id": 28660,
			"enchant": 34004
		},
		{
			"id": 31042,
			"enchant": 24003,
			"gems": [
				32212,
				32220,
				32194
			]
		},
		{
			"id": 33881,
			"enchant": 22533,
			"gems": [
				32194
			]
		},
		{
			"id": 31034,
			"enchant": 33153,
			"gems": [
				32212
			]
		},
		{
			"id": 30106,
			"gems": [
				32194,
				32212
			]
		},
		{
			"id": 31044,
			"enchant": 29536,
			"gems": [
				32212
			]
		},
		{
			"id": 32593,
			"enchant": 35297
		},
		{
			"id": 29279,
			"enchant": 22538
		},
		{
			"id": 32266,
			"enchant": 22538
		},
		{
			"id": 32501
		},
		{
			"id": 32658
		},
		{
			"id": 30883,
			"enchant": 22556
		},
		{
			"id": 32387
		}
	]}`),
};

export const P4_PRESET = {
	name: 'P4 Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 33672,
			"enchant": 29192,
			"gems": [
				32409,
				32220
			]
		},
		{
			"id": 33296
		},
		{
			"id": 31048,
			"enchant": 28911,
			"gems": [
				32194,
				32194
			]
		},
		{
			"id": 33593,
			"enchant": 34004
		},
		{
			"id": 31042,
			"enchant": 24003,
			"gems": [
				32212,
				32220,
				32194
			]
		},
		{
			"id": 33881,
			"enchant": 22533,
			"gems": [
				32194
			]
		},
		{
			"id": 31034,
			"enchant": 33153,
			"gems": [
				32212
			]
		},
		{
			"id": 30106,
			"gems": [
				32194,
				32212
			]
		},
		{
			"id": 31044,
			"enchant": 29536,
			"gems": [
				32212
			]
		},
		{
			"id": 32593,
			"enchant": 35297
		},
		{
			"id": 29279,
			"enchant": 22538
		},
		{
			"id": 29301,
			"enchant": 22538
		},
		{
			"id": 32501
		},
		{
			"id": 33832
		},
		{
			"id": 30883,
			"enchant": 22556
		},
		{
			"id": 32387
		}
	]}`),
};

export const P5_PRESET = {
	name: 'P5 Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 34404,
			"enchant": 29193,
			"gems": [
				32212,
				25896
			]
		},
		{
			"id": 34178
		},
		{
			"id": 34392,
			"enchant": 28911,
			"gems": [
				32200,
				32200
			]
		},
		{
			"id": 34190,
			"enchant": 34004
		},
		{
			"id": 34211,
			"enchant": 24003,
			"gems": [
				32200,
				32200,
				32200
			]
		},
		{
			"id": 34444,
			"enchant": 22533,
			"gems": [
				32200
			]
		},
		{
			"id": 34408,
			"enchant": 33153,
			"gems": [
				32200,
				32200
			]
		},
		{
			"id": 35156
		},
		{
			"id": 34385,
			"enchant": 29536,
			"gems": [
				32200,
				32200,
				32200
			]
		},
		{
			"id": 34573,
			"enchant": 35297,
			"gems": [
				32200
			]
		},
		{
			"id": 34213,
			"enchant": 22538
		},
		{
			"id": 34361,
			"enchant": 22538
		},
		{
			"id": 32501
		},
		{
			"id": 32658
		},
		{
			"id": 30883,
			"enchant": 22556
		},
		{
			"id": 32387
		}
	]}`),
};
