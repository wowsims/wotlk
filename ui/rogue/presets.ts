import { BattleElixir } from '/tbc/core/proto/common.js';
import { Conjured } from '/tbc/core/proto/common.js';
import { Consumes } from '/tbc/core/proto/common.js';
import { Drums } from '/tbc/core/proto/common.js';
import { EquipmentSpec } from '/tbc/core/proto/common.js';
import { Flask } from '/tbc/core/proto/common.js';
import { Food } from '/tbc/core/proto/common.js';
import { ItemSpec } from '/tbc/core/proto/common.js';
import { Potions } from '/tbc/core/proto/common.js';
import { WeaponImbue } from '/tbc/core/proto/common.js';
import { Faction } from '/tbc/core/proto/common.js';
import { Player } from '/tbc/core/player.js';

import {
	Rogue,
	Rogue_Rotation as RogueRotation,
	Rogue_Rotation_Builder as Builder,
	Rogue_Options as RogueOptions,
} from '/tbc/core/proto/rogue.js';

import * as Enchants from '/tbc/core/constants/enchants.js';
import * as Gems from '/tbc/core/proto_utils/gems.js';
import * as Tooltips from '/tbc/core/constants/tooltips.js';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

// Default talents. Uses the wowhead calculator format, make the talents on
// https://tbc.wowhead.com/talent-calc and copy the numbers in the url.
export const CombatTalents = {
	name: 'Combat',
	data: '0053201252-023305102005015002321051',
};
export const CombatMaceTalents = {
	name: 'Combat Maces',
	data: '005320123-023305002005515002321051',
};
export const MutilateTalents = {
	name: 'Mutilate',
	data: '005323125500102501051-005305200005',
};
export const HemoTalents = {
	name: 'Hemo',
	data: '-02330520100501500232105-500252100230001',
};

export const DefaultRotation = RogueRotation.create({
	builder: Builder.Auto,
	maintainExposeArmor: true,
	useRupture: true,
	useShiv: true,
	minComboPointsForDamageFinisher: 3,
});

export const DefaultOptions = RogueOptions.create({
});

export const DefaultConsumes = Consumes.create({
	defaultPotion: Potions.HastePotion,
	defaultConjured: Conjured.ConjuredRogueThistleTea,
	battleElixir: BattleElixir.ElixirOfMajorAgility,
	food: Food.FoodSpicyHotTalbuk,
	mainHandImbue: WeaponImbue.WeaponImbueAdamantiteSharpeningStone,
	offHandImbue: WeaponImbue.WeaponImbueRogueDeadlyPoison,
	scrollOfAgility: 5,
	scrollOfStrength: 5,
});

export const P1_PRESET = {
	name: 'P1 Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 29044,
			"enchant": 29192,
			"gems": [
				32409,
				24061
			]
		},
		{
			"id": 29381
		},
		{
			"id": 27797,
			"enchant": 28888,
			"gems": [
				24061,
				24055
			]
		},
		{
			"id": 28672,
			"enchant": 34004
		},
		{
			"id": 29045,
			"enchant": 24003,
			"gems": [
				24061,
				24051,
				24055
			]
		},
		{
			"id": 29246,
			"enchant": 34002
		},
		{
			"id": 27531,
			"enchant": 19445,
			"gems": [
				24061,
				24061
			]
		},
		{
			"id": 29247
		},
		{
			"id": 28741,
			"enchant": 29535,
			"gems": [
				24051,
				24051,
				24051
			]
		},
		{
			"id": 28545,
			"enchant": 28279,
			"gems": [
				24061,
				24051
			]
		},
		{
			"id": 28757
		},
		{
			"id": 28649
		},
		{
			"id": 29383
		},
		{
			"id": 28830
		},
		{
			"id": 28729,
			"enchant": 22559
		},
		{
			"id": 28189,
			"enchant": 22559
		},
		{
			"id": 28772
		}
	]}`),
};

export const P2_PRESET = {
	name: 'P2 Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 30146,
			"enchant": 29192,
			"gems": [
				32409,
				24061
			]
		},
		{
			"id": 29381
		},
		{
			"id": 30149,
			"enchant": 28888,
			"gems": [
				24061,
				24055
			]
		},
		{
			"id": 28672,
			"enchant": 34004
		},
		{
			"id": 30101,
			"enchant": 24003,
			"gems": [
				24051,
				24051,
				24055
			]
		},
		{
			"id": 29966,
			"enchant": 34002,
			"gems": [
				24051
			]
		},
		{
			"id": 30145,
			"enchant": 19445
		},
		{
			"id": 30106,
			"gems": [
				24051,
				24051
			]
		},
		{
			"id": 30148,
			"enchant": 29535,
			"gems": [
				24051
			]
		},
		{
			"id": 28545,
			"enchant": 28279,
			"gems": [
				24061,
				24051
			]
		},
		{
			"id": 29997
		},
		{
			"id": 30052
		},
		{
			"id": 28830
		},
		{
			"id": 30450
		},
		{
			"id": 30082,
			"enchant": 22559
		},
		{
			"id": 28189,
			"enchant": 22559
		},
		{
			"id": 29949,
			"enchant": 23766
		}
	]}`),
};

export const P3_PRESET = {
	name: 'P3 Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 32235,
			"enchant": 29192,
			"gems": [
				32409,
				32220
			]
		},
		{
			"id": 32260
		},
		{
			"id": 31030,
			"enchant": 28888,
			"gems": [
				32220,
				32220
			]
		},
		{
			"id": 32323,
			"enchant": 34004
		},
		{
			"id": 31028,
			"enchant": 24003,
			"gems": [
				32206,
				32206,
				32212
			]
		},
		{
			"id": 32324,
			"enchant": 34002,
			"gems": [
				32220
			]
		},
		{
			"id": 31026,
			"enchant": 19445,
			"gems": [
				32220
			]
		},
		{
			"id": 30106,
			"gems": [
				32220,
				32212
			]
		},
		{
			"id": 31029,
			"enchant": 29535,
			"gems": [
				32220
			]
		},
		{
			"id": 32366,
			"enchant": 28279,
			"gems": [
				32220,
				32206
			]
		},
		{
			"id": 32497
		},
		{
			"id": 29301
		},
		{
			"id": 30450
		},
		{
			"id": 28830
		},
		{
			"id": 30082,
			"enchant": 22559
		},
		{
			"id": 32369,
			"enchant": 22559
		},
		{
			"id": 29949,
			"enchant": 23766
		}
	]}`),
};

export const P4_PRESET = {
	name: 'P4 Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 32235,
			"enchant": 29192,
			"gems": [
				32409,
				32220
			]
		},
		{
			"id": 32260
		},
		{
			"id": 31030,
			"enchant": 28888,
			"gems": [
				32206,
				32206
			]
		},
		{
			"id": 33590,
			"enchant": 34004
		},
		{
			"id": 31028,
			"enchant": 24003,
			"gems": [
				32206,
				32206,
				32212
			]
		},
		{
			"id": 32324,
			"enchant": 34002,
			"gems": [
				32206
			]
		},
		{
			"id": 31026,
			"enchant": 19445,
			"gems": [
				32220
			]
		},
		{
			"id": 30106,
			"gems": [
				32220,
				32212
			]
		},
		{
			"id": 31029,
			"enchant": 29535,
			"gems": [
				32220
			]
		},
		{
			"id": 32366,
			"enchant": 28279,
			"gems": [
				32220,
				32206
			]
		},
		{
			"id": 33496
		},
		{
			"id": 32497
		},
		{
			"id": 30450
		},
		{
			"id": 28830
		},
		{
			"id": 30082,
			"enchant": 22559
		},
		{
			"id": 32369,
			"enchant": 22559
		},
		{
			"id": 29949,
			"enchant": 23766
		}
	]}`),
};

export const P5_PRESET = {
	name: 'P5 Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 34244,
			"enchant": 29192,
			"gems": [
				32409,
				32220
			]
		},
		{
			"id": 34177
		},
		{
			"id": 31030,
			"enchant": 28888,
			"gems": [
				32220,
				32212
			]
		},
		{
			"id": 34241,
			"enchant": 34004,
			"gems": [
				32206
			]
		},
		{
			"id": 34397,
			"enchant": 24003,
			"gems": [
				32212,
				32220,
				32220
			]
		},
		{
			"id": 34448,
			"enchant": 34002,
			"gems": [
				32220
			]
		},
		{
			"id": 34370,
			"enchant": 19445,
			"gems": [
				32220,
				32220
			]
		},
		{
			"id": 34558,
			"gems": [
				32220
			]
		},
		{
			"id": 34188,
			"enchant": 29535,
			"gems": [
				32220,
				32220,
				32220
			]
		},
		{
			"id": 34575,
			"enchant": 28279,
			"gems": [
				32220
			]
		},
		{
			"id": 32497
		},
		{
			"id": 34189
		},
		{
			"id": 34427
		},
		{
			"id": 28830
		},
		{
			"id": 34331,
			"enchant": 22559,
			"gems": [
				32206,
				32206
			]
		},
		{
			"id": 32369,
			"enchant": 22559
		},
		{
			"id": 34196
		}
	]}`),
};
