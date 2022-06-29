import { Consumes } from '/tbc/core/proto/common.js';
import { EquipmentSpec } from '/tbc/core/proto/common.js';
import { Flask } from '/tbc/core/proto/common.js';
import { Food } from '/tbc/core/proto/common.js';
import { ItemSpec } from '/tbc/core/proto/common.js';
import { Potions } from '/tbc/core/proto/common.js';
import { Spec } from '/tbc/core/proto/common.js';
import { WeaponImbue } from '/tbc/core/proto/common.js';
import { Faction } from '/tbc/core/proto/common.js';
import { Player } from '/tbc/core/player.js';

import {
	WarriorShout,
	WarriorTalents as WarriorTalents,
	Warrior,
	Warrior_Rotation as WarriorRotation,
	Warrior_Rotation_SunderArmor as SunderArmor,
	Warrior_Options as WarriorOptions,
} from '/tbc/core/proto/warrior.js';

import * as Enchants from '/tbc/core/constants/enchants.js';
import * as Gems from '/tbc/core/proto_utils/gems.js';
import * as Tooltips from '/tbc/core/constants/tooltips.js';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

// Default talents. Uses the wowhead calculator format, make the talents on
// https://tbc.wowhead.com/talent-calc and copy the numbers in the url.
export const ArmsSlamTalents = {
	name: 'Arms Slam',
	data: '32003301352010500221-0550000500521203',
};
export const ArmsDWTalents = {
	name: 'Arms DW',
	data: '33005301302010510321-0550000520501203',
};
export const FuryTalents = {
	name: 'Fury',
	data: '3500501130201-05050005505012050115',
};

export const DefaultRotation = WarriorRotation.create({
	useOverpower: false,
	useHamstring: true,
	prioritizeWw: false,
	sunderArmor: SunderArmor.SunderArmorMaintain,
	hsRageThreshold: 60,
	overpowerRageThreshold: 10,
	hamstringRageThreshold: 75,
	rampageCdThreshold: 5,
	slamLatency: 150,
	slamGcdDelay: 400,
	slamMsWwDelay: 2000,
	useHsDuringExecute: true,
	useMsDuringExecute: true,
	useBtDuringExecute: true,
	useWwDuringExecute: true,
	useSlamDuringExecute: true,
});

export const ArmsRotation = WarriorRotation.create({
	useOverpower: false,
	useHamstring: true,
	useSlam: true,
	prioritizeWw: false,
	sunderArmor: SunderArmor.SunderArmorMaintain,
	hsRageThreshold: 60,
	overpowerRageThreshold: 10,
	hamstringRageThreshold: 75,
	rampageCdThreshold: 5,
	slamLatency: 150,
	slamGcdDelay: 400,
	slamMsWwDelay: 2000,
	useHsDuringExecute: true,
	useMsDuringExecute: true,
	useBtDuringExecute: true,
	useWwDuringExecute: true,
	useSlamDuringExecute: true,
});

export const DefaultOptions = WarriorOptions.create({
	startingRage: 0,
	useRecklessness: true,
	shout: WarriorShout.WarriorShoutBattle,
	precastShout: true,
	precastShoutSapphire: false,
	precastShoutT2: false,
});

export const DefaultConsumes = Consumes.create({
	flask: Flask.FlaskOfRelentlessAssault,
	food: Food.FoodRoastedClefthoof,
	defaultPotion: Potions.HastePotion,
	mainHandImbue: WeaponImbue.WeaponImbueAdamantiteSharpeningStone,
	offHandImbue: WeaponImbue.WeaponImbueAdamantiteSharpeningStone,
});

export const P1_FURY_PRESET = {
	name: 'P1 Fury Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecWarrior>) => player.getTalents().bloodthirst,
	gear: EquipmentSpec.fromJsonString(`{"items": [
	  {
			"id": 29021,
			"enchant": 29192,
			"gems": [
				32409,
				24048
			]
		},
		{
			"id": 29381
		},
		{
			"id": 29023,
			"enchant": 28888,
			"gems": [
				24048,
				24067
			]
		},
		{
			"id": 24259,
			"enchant": 34004,
			"gems": [
				24058
			]
		},
		{
			"id": 29019,
			"enchant": 24003,
			"gems": [
				24048,
				24048,
				24048
			]
		},
		{
			"id": 28795,
			"enchant": 27899,
			"gems": [
				24067,
				24058
			]
		},
		{
			"id": 28824,
			"enchant": 33995,
			"gems": [
				24067,
				24048
			]
		},
		{
			"id": 28779,
			"gems": [
				24058,
				24067
			]
		},
		{
			"id": 28741,
			"enchant": 29535,
			"gems": [
				24048,
				24048,
				24048
			]
		},
		{
			"id": 28608,
			"enchant": 28279,
			"gems": [
				24058,
				24048
			]
		},
		{
			"id": 28757
		},
		{
			"id": 30834
		},
		{
			"id": 29383
		},
		{
			"id": 28830
		},
		{
			"id": 28438,
			"enchant": 22559
		},
		{
			"id": 28729,
			"enchant": 22559
		},
		{
			"id": 30279
		}
	]}`),
};

export const P2_FURY_PRESET = {
	name: 'P2 Fury Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecWarrior>) => player.getTalents().bloodthirst,
	gear: EquipmentSpec.fromJsonString(`{"items": [
	  {
			"id": 30120,
			"enchant": 29192,
			"gems": [
				32409,
				24067
			]
		},
		{
			"id": 30022
		},
		{
			"id": 30122,
			"enchant": 28888,
			"gems": [
				24058,
				24067
			]
		},
		{
			"id": 24259,
			"enchant": 34004,
			"gems": [
				24058
			]
		},
		{
			"id": 30118,
			"enchant": 24003,
			"gems": [
				24048,
				24067,
				24058
			]
		},
		{
			"id": 28795,
			"enchant": 27899,
			"gems": [
				24067,
				24058
			]
		},
		{
			"id": 30119,
			"enchant": 33995
		},
		{
			"id": 30106,
			"gems": [
				24048,
				24048
			]
		},
		{
			"id": 29995,
			"enchant": 29535
		},
		{
			"id": 30081,
			"enchant": 28279
		},
		{
			"id": 29997
		},
		{
			"id": 28757
		},
		{
			"id": 30627
		},
		{
			"id": 28830
		},
		{
			"id": 28439,
			"enchant": 22559
		},
		{
			"id": 30082,
			"enchant": 22559
		},
		{
			"id": 30105
		}
	]}`),
};

export const P3_FURY_PRESET = {
	name: 'P3 Fury Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecWarrior>) => player.getTalents().bloodthirst,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 30972,
			"enchant": 29192,
			"gems": [
				32409,
				32205
			]
		},
		{
			"id": 32260
		},
		{
			"id": 30979,
			"enchant": 28888,
			"gems": [
				32205,
				32226
			]
		},
		{
			"id": 32323,
			"enchant": 34004
		},
		{
			"id": 30975,
			"enchant": 24003,
			"gems": [
				32217,
				32226,
				32226
			]
		},
		{
			"id": 30863,
			"enchant": 27899,
			"gems": [
				32205
			]
		},
		{
			"id": 30969,
			"enchant": 33995,
			"gems": [
				32217
			]
		},
		{
			"id": 30106,
			"gems": [
				32205,
				32205
			]
		},
		{
			"id": 32341,
			"enchant": 29535
		},
		{
			"id": 32345,
			"enchant": 28279,
			"gems": [
				32205,
				32205
			]
		},
		{
			"id": 32497
		},
		{
			"id": 32335
		},
		{
			"id": 32505
		},
		{
			"id": 28830
		},
		{
			"id": 28439,
			"enchant": 22559
		},
		{
			"id": 30881,
			"enchant": 22559
		},
		{
			"id": 30105
		}
	]}`),
};

export const P4_FURY_PRESET = {
	name: 'P4 Fury Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecWarrior>) => player.getTalents().bloodthirst,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 30972,
			"enchant": 29192,
			"gems": [
				32409,
				32205
			]
		},
		{
			"id": 32260
		},
		{
			"id": 30979,
			"enchant": 28888,
			"gems": [
				32205,
				32226
			]
		},
		{
			"id": 32323,
			"enchant": 34004
		},
		{
			"id": 30975,
			"enchant": 24003,
			"gems": [
				32217,
				32226,
				32226
			]
		},
		{
			"id": 30863,
			"enchant": 27899,
			"gems": [
				32205
			]
		},
		{
			"id": 30969,
			"enchant": 33995,
			"gems": [
				32217
			]
		},
		{
			"id": 30106,
			"gems": [
				32205,
				32205
			]
		},
		{
			"id": 32341,
			"enchant": 29535
		},
		{
			"id": 32345,
			"enchant": 28279,
			"gems": [
				32205,
				32205
			]
		},
		{
			"id": 32497
		},
		{
			"id": 33496
		},
		{
			"id": 32505
		},
		{
			"id": 28830
		},
		{
			"id": 28439,
			"enchant": 22559
		},
		{
			"id": 30881,
			"enchant": 22559
		},
		{
			"id": 33474
		}
	]}`),
};

export const P5_FURY_PRESET = {
	name: 'P5 Fury Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecWarrior>) => player.getTalents().bloodthirst,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 34333,
			"enchant": 29192,
			"gems": [
				32193,
				32409
			]
		},
		{
			"id": 34358,
			"gems": [
				32217
			]
		},
		{
			"id": 34392,
			"enchant": 28910,
			"gems": [
				32193,
				32211
			]
		},
		{
			"id": 34241,
			"enchant": 34004,
			"gems": [
				33143
			]
		},
		{
			"id": 34397,
			"enchant": 24003,
			"gems": [
				32211,
				32217,
				32193
			]
		},
		{
			"id": 34441,
			"enchant": 27899,
			"gems": [
				32193
			]
		},
		{
			"id": 34343,
			"enchant": 33995,
			"gems": [
				32193,
				32217
			]
		},
		{
			"id": 34546,
			"gems": [
				32217
			]
		},
		{
			"id": 34188,
			"enchant": 29535,
			"gems": [
				32193,
				32193,
				32193
			]
		},
		{
			"id": 34569,
			"enchant": 28279,
			"gems": [
				32217
			]
		},
		{
			"id": 34189
		},
		{
			"id": 34361
		},
		{
			"id": 28830
		},
		{
			"id": 34427
		},
		{
			"id": 34331,
			"enchant": 33307,
			"gems": [
				32205,
				32205
			]
		},
		{
			"id": 34203,
			"enchant": 22559
		},
		{
			"id": 34196
		}
	]}`),
};

export const P1_ARMS_PRESET = {
	name: 'P1 Arms Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecWarrior>) => player.getTalents().mortalStrike,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 29021,
			"enchant": 29192,
			"gems": [
				32409,
				24048
			]
		},
		{
			"id": 29349
		},
		{
			"id": 29023,
			"enchant": 28888,
			"gems": [
				24048,
				24067
			]
		},
		{
			"id": 24259,
			"enchant": 34004,
			"gems": [
				24058
			]
		},
		{
			"id": 29019,
			"enchant": 24003,
			"gems": [
				24048,
				24048,
				24048
			]
		},
		{
			"id": 28795,
			"enchant": 27899,
			"gems": [
				24067,
				24058
			]
		},
		{
			"id": 28824,
			"enchant": 33995,
			"gems": [
				24067,
				24048
			]
		},
		{
			"id": 28779,
			"gems": [
				24058,
				24067
			]
		},
		{
			"id": 28741,
			"enchant": 29535,
			"gems": [
				24048,
				24048,
				24048
			]
		},
		{
			"id": 28608,
			"enchant": 28279,
			"gems": [
				24058,
				24048
			]
		},
		{
			"id": 28757
		},
		{
			"id": 28730
		},
		{
			"id": 29383
		},
		{
			"id": 28830
		},
		{
			"id": 28429,
			"enchant": 22559
		},
		{
			"id": 30279
		}
	]}`),
};

export const P2_ARMS_PRESET = {
	name: 'P2 Arms Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecWarrior>) => player.getTalents().mortalStrike,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 30120,
			"enchant": 29192,
			"gems": [
				32409,
				24067
			]
		},
		{
			"id": 30022
		},
		{
			"id": 30122,
			"enchant": 28888,
			"gems": [
				24058,
				24067
			]
		},
		{
			"id": 24259,
			"enchant": 34004,
			"gems": [
				24058
			]
		},
		{
			"id": 30118,
			"enchant": 24003,
			"gems": [
				24048,
				24067,
				24058
			]
		},
		{
			"id": 28795,
			"enchant": 27899,
			"gems": [
				24067,
				24058
			]
		},
		{
			"id": 30119,
			"enchant": 33995
		},
		{
			"id": 30106,
			"gems": [
				24048,
				24048
			]
		},
		{
			"id": 30121,
			"enchant": 29535,
			"gems": [
				24058
			]
		},
		{
			"id": 30081,
			"enchant": 28279
		},
		{
			"id": 29997
		},
		{
			"id": 28757
		},
		{
			"id": 30627
		},
		{
			"id": 28830
		},
		{
			"id": 29993,
			"enchant": 22559,
			"gems": [
				24048,
				24048,
				24048
			]
		},
		{
			"id": 30105
		}
	]}`),
};

export const P3_ARMS_PRESET = {
	name: 'P3 Arms Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecWarrior>) => player.getTalents().mortalStrike,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 30972,
			"enchant": 29192,
			"gems": [
				32409,
				32205
			]
		},
		{
			"id": 32260
		},
		{
			"id": 30979,
			"enchant": 28888,
			"gems": [
				32205,
				32226
			]
		},
		{
			"id": 32323,
			"enchant": 34004
		},
		{
			"id": 30975,
			"enchant": 24003,
			"gems": [
				32217,
				32226,
				32226
			]
		},
		{
			"id": 30863,
			"enchant": 27899,
			"gems": [
				32205
			]
		},
		{
			"id": 30969,
			"enchant": 33995,
			"gems": [
				32217
			]
		},
		{
			"id": 30106,
			"gems": [
				32205,
				32205
			]
		},
		{
			"id": 32341,
			"enchant": 29535
		},
		{
			"id": 32345,
			"enchant": 28279,
			"gems": [
				32205,
				32205
			]
		},
		{
			"id": 32497
		},
		{
			"id": 32335
		},
		{
			"id": 32505
		},
		{
			"id": 28830
		},
		{
			"id": 30902,
			"enchant": 22559
		},
		{
			"id": 30105
		}
	]}`),
};

export const P4_ARMS_PRESET = {
	name: 'P4 Arms Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecWarrior>) => player.getTalents().mortalStrike,
	gear: EquipmentSpec.fromJsonString(`{"items": [
	  {
			"id": 30972,
			"enchant": 29192,
			"gems": [
				32409,
				32205
			]
		},
		{
			"id": 32260
		},
		{
			"id": 30979,
			"enchant": 28888,
			"gems": [
				32205,
				32226
			]
		},
		{
			"id": 32323,
			"enchant": 34004
		},
		{
			"id": 30975,
			"enchant": 24003,
			"gems": [
				32217,
				32226,
				32226
			]
		},
		{
			"id": 30863,
			"enchant": 27899,
			"gems": [
				32205
			]
		},
		{
			"id": 30969,
			"enchant": 33995,
			"gems": [
				32217
			]
		},
		{
			"id": 30106,
			"gems": [
				32205,
				32205
			]
		},
		{
			"id": 32341,
			"enchant": 29535
		},
		{
			"id": 32345,
			"enchant": 28279,
			"gems": [
				32205,
				32205
			]
		},
		{
			"id": 32497
		},
		{
			"id": 33496
		},
		{
			"id": 32505
		},
		{
			"id": 28830
		},
		{
			"id": 30902,
			"enchant": 22559
		},
		{
			"id": 33474
		}
	]}`),
};

export const P5_ARMS_PRESET = {
	name: 'P5 Arms Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecWarrior>) => player.getTalents().mortalStrike,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 34333,
			"enchant": 29192,
			"gems": [
				32193,
				32409
			]
		},
		{
			"id": 34358,
			"gems": [
				32217
			]
		},
		{
			"id": 34392,
			"enchant": 28910,
			"gems": [
				32193,
				32211
			]
		},
		{
			"id": 34241,
			"enchant": 34004,
			"gems": [
				33143
			]
		},
		{
			"id": 34397,
			"enchant": 24003,
			"gems": [
				32211,
				32217,
				32193
			]
		},
		{
			"id": 34441,
			"enchant": 27899,
			"gems": [
				32193
			]
		},
		{
			"id": 34343,
			"enchant": 33995,
			"gems": [
				32193,
				32217
			]
		},
		{
			"id": 34546,
			"gems": [
				32217
			]
		},
		{
			"id": 34188,
			"enchant": 29535,
			"gems": [
				32193,
				32193,
				32193
			]
		},
		{
			"id": 34569,
			"enchant": 28279,
			"gems": [
				32217
			]
		},
		{
			"id": 34189
		},
		{
			"id": 34361
		},
		{
			"id": 28830
		},
		{
			"id": 34427
		},
		{
			"id": 34247,
			"enchant": 33307,
			"gems": [
				32205,
				32205,
				32205
			]
		},
		{
			"id": 34196
		}
	]}`),
};
