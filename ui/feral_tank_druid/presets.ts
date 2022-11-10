import { Consumes } from '../core/proto/common.js';
import { BattleElixir } from '../core/proto/common.js';
import { GuardianElixir } from '../core/proto/common.js';
import { Food } from '../core/proto/common.js';
import { EquipmentSpec } from '../core/proto/common.js';
import { Potions } from '../core/proto/common.js';
import { Conjured } from '../core/proto/common.js';
import { RaidTarget } from '../core/proto/common.js';
import { SavedTalents } from '../core/proto/ui.js';
import { NO_TARGET } from '../core/proto_utils/utils.js';

import {
	FeralTankDruid_Rotation as DruidRotation,
	FeralTankDruid_Rotation_Swipe as Swipe,
	FeralTankDruid_Options as DruidOptions
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
		talentsString: '-503032132322105301251-05503301',
	}),
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
});

export const P1_PRESET = {
	name: 'TBC P5 Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 34404,
			"enchant": 3004,
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
			"enchant": 2991,
			"gems": [
				32200,
				32200
			]
		},
		{
			"id": 34190,
			"enchant": 368
		},
		{
			"id": 34211,
			"enchant": 2661,
			"gems": [
				32200,
				32200,
				32200
			]
		},
		{
			"id": 34444,
			"enchant": 2649,
			"gems": [
				32200,
				0
			]
		},
		{
			"id": 34408,
			"enchant": 2613,
			"gems": [
				32200,
				32200,
				0
			]
		},
		{
			"id": 35156,
			"gems": [
				0
			]
		},
		{
			"id": 34385,
			"enchant": 3013,
			"gems": [
				32200,
				32200,
				32200
			]
		},
		{
			"id": 34573,
			"enchant": 2940,
			"gems": [
				32200
			]
		},
		{
			"id": 34213,
			"enchant": 2931
		},
		{
			"id": 34361,
			"enchant": 2931
		},
		{
			"id": 32501
		},
		{
			"id": 32658
		},
		{
			"id": 30883,
			"enchant": 2670
		},
		{},
		{
			"id": 32387
		}
	]}`),
};
