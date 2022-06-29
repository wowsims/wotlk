import { Consumes } from '/tbc/core/proto/common.js';
import { BattleElixir } from '/tbc/core/proto/common.js';
import { Food } from '/tbc/core/proto/common.js';
import { EquipmentSpec } from '/tbc/core/proto/common.js';
import { ItemSpec } from '/tbc/core/proto/common.js';
import { Potions } from '/tbc/core/proto/common.js';
import { Conjured } from '/tbc/core/proto/common.js';
import { WeaponImbue } from '/tbc/core/proto/common.js';
import { Faction } from '/tbc/core/proto/common.js';
import { Player } from '/tbc/core/player.js';

import { FeralDruid, FeralDruid_Rotation as FeralDruidRotation, DruidTalents as DruidTalents, FeralDruid_Options as FeralDruidOptions } from '/tbc/core/proto/druid.js';
import { FeralDruid_Rotation_FinishingMove as FinishingMove } from '/tbc/core/proto/druid.js';

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

export const MonocatTalents = {
	name: 'Monocat',
	data: '-553002132322105301051-05503301',
};


export const DefaultRotation = FeralDruidRotation.create({
	finishingMove: FinishingMove.Rip,
	mangleTrick: true,
	biteweave: true,
	ripMinComboPoints: 5,
	biteMinComboPoints: 5,
	rakeTrick: false,
	ripweave: false,
	maintainFaerieFire: true,
});

export const DefaultOptions = FeralDruidOptions.create({
});

export const DefaultConsumes = Consumes.create({
	battleElixir: BattleElixir.ElixirOfMajorAgility,
	food: Food.FoodGrilledMudfish,
	mainHandImbue: WeaponImbue.WeaponImbueAdamantiteWeightstone,
	defaultPotion: Potions.HastePotion,
	defaultConjured: Conjured.ConjuredDarkRune,
	scrollOfAgility: 5,
	scrollOfStrength: 5,
});

export const P1_PRESET = {
	name: 'P1 Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 8345,
			"enchant": 29192
		},
		{
			"id": 29381
		},
		{
			"id": 29100,
			"enchant": 28888,
			"gems": [
				24028,
				24028
			]
		},
		{
			"id": 28672,
			"enchant": 34004
		},
		{
			"id": 29096,
			"enchant": 24003,
			"gems": [
				24028,
				24028,
				24028
			]
		},
		{
			"id": 29246,
			"enchant": 27899
		},
		{
			"id": 28506,
			"gems": [
				24028,
				24028
			]
		},
		{
			"id": 28750,
			"gems": [
				24028,
				24028
			]
		},
		{
			"id": 28741,
			"enchant": 29535,
			"gems": [
				24028,
				24028,
				24028
			]
		},
		{
			"id": 28545,
			"enchant": 28279,
			"gems": [
				24028,
				24028
			]
		},
		{
			"id": 28649,
			"enchant": 22535
		},
		{
			"id": 30834,
			"enchant": 22535
		},
		{
			"id": 28830
		},
		{
			"id": 29383
		},
		{
			"id": 28658,
			"enchant": 22556
		},
		{
			"id": 29390
		}
	]}`),
};

export const P2_PRESET = {
	name: 'P2 Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 8345,
			"enchant": 29192
		},
		{
			"id": 30017
		},
		{
			"id": 29100,
			"enchant": 28888,
			"gems": [
				24028,
				24028
			]
		},
		{
			"id": 29994,
			"enchant": 34004
		},
		{
			"id": 29096,
			"enchant": 24003,
			"gems": [
				24028,
				24028,
				24028
			]
		},
		{
			"id": 29966,
			"enchant": 27899,
			"gems": [
				24028
			]
		},
		{
			"id": 29947,
			"enchant": 19445
		},
		{
			"id": 30106,
			"gems": [
				24028,
				30549
			]
		},
		{
			"id": 29995,
			"enchant": 29535
		},
		{
			"id": 28545,
			"enchant": 28279,
			"gems": [
				24028,
				24028
			]
		},
		{
			"id": 30052,
			"enchant": 22535
		},
		{
			"id": 29997,
			"enchant": 22535
		},
		{
			"id": 30627
		},
		{
			"id": 29383
		},
		{
			"id": 32014,
			"enchant": 22556
		},
		{
			"id": 29390
		}
	]}`),
};

export const P3_PRESET = {
	name: 'P3 Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 8345,
			"enchant": 29192
		},
		{
			"id": 32260
		},
		{
			"id": 31048,
			"enchant": 23548,
			"gems": [
				32194,
				32194
			]
		},
		{
			"id": 29994,
			"enchant": 34004
		},
		{
			"id": 31042,
			"enchant": 24003,
			"gems": [
				32194,
				32194,
				32194
			]
		},
		{
			"id": 33881,
			"enchant": 27899,
			"gems": [
				32194
			]
		},
		{
			"id": 31034,
			"enchant": 19445,
			"gems": [
				32194
			]
		},
		{
			"id": 30106,
			"gems": [
				32194,
				32194
			]
		},
		{
			"id": 31044,
			"enchant": 29535,
			"gems": [
				32194
			]
		},
		{
			"id": 32366,
			"enchant": 28279,
			"gems": [
				32194,
				32194
			]
		},
		{
			"id": 32335,
			"enchant": 22538
		},
		{
			"id": 29301,
			"enchant": 22538
		},
		{
			"id": 30627
		},
		{
			"id": 29383
		},
		{
			"id": 33716,
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
			"id": 8345,
			"enchant": 29192
		},
		{
			"id": 24114
		},
		{
			"id": 31048,
			"enchant": 23548,
			"gems": [
				32194,
				32194
			]
		},
		{
			"id": 29994,
			"enchant": 34004
		},
		{
			"id": 31042,
			"enchant": 24003,
			"gems": [
				32194,
				32194,
				32194
			]
		},
		{
			"id": 33881,
			"enchant": 27899,
			"gems": [
				32194
			]
		},
		{
			"id": 31034,
			"enchant": 33152,
			"gems": [
				32194
			]
		},
		{
			"id": 30106,
			"gems": [
				32194,
				32194
			]
		},
		{
			"id": 31044,
			"enchant": 29535,
			"gems": [
				32194
			]
		},
		{
			"id": 32366,
			"enchant": 28279,
			"gems": [
				32194,
				32194
			]
		},
		{
			"id": 29301,
			"enchant": 22538
		},
		{
			"id": 33496,
			"enchant": 22538
		},
		{
			"id": 30627
		},
		{
			"id": 33831
		},
		{
			"id": 33716,
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
			"id": 32235,
			"enchant": 29192,
			"gems": [
				32409,
				32194
			]
		},
		{
			"id": 34177
		},
		{
			"id": 31048,
			"enchant": 23548,
			"gems": [
				32194,
				32194
			]
		},
		{
			"id": 34241,
			"enchant": 34004,
			"gems": [
				32194
			]
		},
		{
			"id": 34397,
			"enchant": 24003,
			"gems": [
				32226,
				32226,
				32194
			]
		},
		{
			"id": 34444,
			"enchant": 27899,
			"gems": [
				32194
			]
		},
		{
			"id": 34370,
			"enchant": 19445,
			"gems": [
				32194,
				32194
			]
		},
		{
			"id": 34556,
			"gems": [
				32194
			]
		},
		{
			"id": 34188,
			"enchant": 29535,
			"gems": [
				32194,
				32194,
				32194
			]
		},
		{
			"id": 34573,
			"enchant": 28279,
			"gems": [
				32194
			]
		},
		{
			"id": 34887,
			"enchant": 22538
		},
		{
			"id": 34189,
			"enchant": 22538
		},
		{
			"id": 34472
		},
		{
			"id": 34427
		},
		{
			"id": 34198,
			"enchant": 22556
		},
		{
			"id": 32387
		}
	]}`),
};
