import { Consumes } from '/tbc/core/proto/common.js';
import { EquipmentSpec } from '/tbc/core/proto/common.js';
import { Flask } from '/tbc/core/proto/common.js';
import { Food } from '/tbc/core/proto/common.js';
import { Potions } from '/tbc/core/proto/common.js';
import { WeaponImbue } from '/tbc/core/proto/common.js';
import { Player } from '/tbc/core/player.js';

import {
	Hunter_Rotation as HunterRotation,
	Hunter_Rotation_WeaveType as WeaveType,
	Hunter_Options as HunterOptions,
	Hunter_Options_Ammo as Ammo,
	Hunter_Options_QuiverBonus as QuiverBonus,
	Hunter_Options_PetType as PetType,
} from '/tbc/core/proto/hunter.js';

import * as Tooltips from '/tbc/core/constants/tooltips.js';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

// Default talents. Uses the wowhead calculator format, make the talents on
// https://tbc.wowhead.com/talent-calc and copy the numbers in the url.
export const BeastMasteryTalents = {
	name: 'BM',
	data: '512002015150122431051-0505201205',
};

export const MarksmanTalents = {
	name: 'Marksman',
	data: '51200200502-0551201205013253135',
};

export const SurvivalTalents = {
	name: 'Survival',
	data: '502-0550201205-333200022003223005103',
};

export const DefaultRotation = HunterRotation.create({
	useMultiShot: true,
	useArcaneShot: true,
	viperStartManaPercent: 0.1,
	viperStopManaPercent: 0.3,

	weave: WeaveType.WeaveNone,
	timeToWeaveMs: 500,
	percentWeaved: 0.8,
});

export const DefaultOptions = HunterOptions.create({
	quiverBonus: QuiverBonus.Speed15,
	ammo: Ammo.TimelessArrow,
	petType: PetType.Ravager,
	petUptime: 1,
	latencyMs: 30,
});

export const DefaultConsumes = Consumes.create({
	defaultPotion: Potions.HastePotion,
	flask: Flask.FlaskOfRelentlessAssault,
	food: Food.FoodGrilledMudfish,
	mainHandImbue: WeaponImbue.WeaponImbueAdamantiteSharpeningStone,
	offHandImbue: WeaponImbue.WeaponImbueAdamantiteSharpeningStone,
});

export const P1_BM_PRESET = {
	name: 'P1 BM Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<any>) => player.getTalentTree() != 2,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 28275,
			"enchant": 29192,
			"gems": [
				24028,
				32409
			]
		},
		{
			"id": 29381
		},
		{
			"id": 27801,
			"enchant": 28888,
			"gems": [
				31868,
				24028
			]
		},
		{
			"id": 24259,
			"enchant": 34004,
			"gems": [
				24028
			]
		},
		{
			"id": 28228,
			"enchant": 24003,
			"gems": [
				24028,
				24028,
				24055
			]
		},
		{
			"id": 29246,
			"enchant": 34002
		},
		{
			"id": 27474,
			"enchant": 33152,
			"gems": [
				24028,
				24028
			]
		},
		{
			"id": 28828,
			"gems": [
				24055,
				31868
			]
		},
		{
			"id": 30739,
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
				24061
			]
		},
		{
			"id": 28757
		},
		{
			"id": 28791
		},
		{
			"id": 28830
		},
		{
			"id": 29383
		},
		{
			"id": 28435,
			"enchant": 22556
		},
		{
			"id": 28772,
			"enchant": 23766
		}
	]}`),
};

export const P2_BM_PRESET = {
	name: 'P2 BM Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<any>) => player.getTalentTree() != 2,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 30141,
			"enchant": 29192,
			"gems": [
				24028,
				32409
			]
		},
		{
			"id": 30017
		},
		{
			"id": 30143,
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
			"id": 30139,
			"enchant": 24003,
			"gems": [
				24055,
				31868,
				31868
			]
		},
		{
			"id": 29966,
			"enchant": 34002,
			"gems": [
				24028
			]
		},
		{
			"id": 30140,
			"enchant": 33152
		},
		{
			"id": 30040,
			"gems": [
				24028,
				24028
			]
		},
		{
			"id": 29995,
			"enchant": 29535
		},
		{
			"id": 30104,
			"enchant": 28279,
			"gems": [
				24055,
				24028
			]
		},
		{
			"id": 29997
		},
		{
			"id": 28791
		},
		{
			"id": 28830
		},
		{
			"id": 29383
		},
		{
			"id": 29993,
			"enchant": 22556,
			"gems": [
				24028,
				24028,
				24028
			]
		},
		{
			"id": 30105,
			"enchant": 23766
		}
	]}`),
};

export const P3_BM_PRESET = {
	name: 'P3 BM Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<any>) => player.getTalentTree() != 2,
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
			"id": 32591
		},
		{
			"id": 31006,
			"enchant": 28888,
			"gems": [
				32194,
				32194
			]
		},
		{
			"id": 32323,
			"enchant": 34004
		},
		{
			"id": 31004,
			"enchant": 24003,
			"gems": [
				32194,
				32226,
				32226
			]
		},
		{
			"id": 32324,
			"enchant": 34002,
			"gems": [
				32222
			]
		},
		{
			"id": 31001,
			"enchant": 33152,
			"gems": [
				32194
			]
		},
		{
			"id": 30879,
			"gems": [
				32194,
				32194
			]
		},
		{
			"id": 31005,
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
				32222
			]
		},
		{
			"id": 29997
		},
		{
			"id": 29301
		},
		{
			"id": 28830
		},
		{
			"id": 32505
		},
		{
			"id": 30901,
			"enchant": 33165
		},
		{
			"id": 30881,
			"enchant": 33165
		},
		{
			"id": 30906,
			"enchant": 23766
		}
	]}`),
};

export const P1_SV_PRESET = {
	name: 'P1 SV Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<any>) => player.getTalentTree() == 2,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 28275,
			"enchant": 29192,
			"gems": [
				24061,
				32409
			]
		},
		{
			"id": 28343
		},
		{
			"id": 27801,
			"enchant": 28888,
			"gems": [
				24028,
				24028
			]
		},
		{
			"id": 29382,
			"enchant": 34004
		},
		{
			"id": 28228,
			"enchant": 24003,
			"gems": [
				24028,
				24028,
				24055
			]
		},
		{
			"id": 25697,
			"enchant": 34002,
			"gems": [
				24028
			]
		},
		{
			"id": 27474,
			"enchant": 33152,
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
			"enchant": 22544,
			"gems": [
				24028,
				24028
			]
		},
		{
			"id": 31277
		},
		{
			"id": 28791
		},
		{
			"id": 28830
		},
		{
			"id": 29383
		},
		{
			"id": 27846,
			"enchant": 19445,
			"gems": [
				24028,
				24028
			]
		},
		{
			"id": 28572,
			"enchant": 19445,
			"gems": [
				24028,
				24061,
				24055
			]
		},
		{
			"id": 28772,
			"enchant": 23766
		}
	]}`),
};

export const P2_SV_PRESET = {
	name: 'P2 SV Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<any>) => player.getTalentTree() == 2,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 30141,
			"enchant": 29192,
			"gems": [
				24028,
				32409
			]
		},
		{
			"id": 30017
		},
		{
			"id": 30143,
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
			"id": 30054,
			"enchant": 24003,
			"gems": [
				24028,
				24028,
				24028
			]
		},
		{
			"id": 29966,
			"enchant": 34002,
			"gems": [
				24028
			]
		},
		{
			"id": 28506,
			"enchant": 33152,
			"gems": [
				24028,
				24028
			]
		},
		{
			"id": 30040,
			"gems": [
				24028,
				24028
			]
		},
		{
			"id": 29985,
			"enchant": 29535,
			"gems": [
				24028,
				24061,
				24055
			]
		},
		{
			"id": 30104,
			"enchant": 22544,
			"gems": [
				24067,
				24028
			]
		},
		{
			"id": 29298
		},
		{
			"id": 28791
		},
		{
			"id": 28830
		},
		{
			"id": 29383
		},
		{
			"id": 29924,
			"enchant": 19445
		},
		{
			"id": 29948,
			"enchant": 19445
		},
		{
			"id": 30105,
			"enchant": 23766
		}
	]}`),
};

export const P3_SV_PRESET = {
	name: 'P3 SV Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<any>) => player.getTalentTree() == 2,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 31003,
			"enchant": 29192,
			"gems": [
				32194,
				32409
			]
		},
		{
			"id": 30017
		},
		{
			"id": 31006,
			"enchant": 28888,
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
			"id": 31004,
			"enchant": 24003,
			"gems": [
				32194,
				32226,
				32226
			]
		},
		{
			"id": 32324,
			"enchant": 34002,
			"gems": [
				32194
			]
		},
		{
			"id": 31001,
			"enchant": 33152,
			"gems": [
				32194
			]
		},
		{
			"id": 30879,
			"gems": [
				32194,
				32194
			]
		},
		{
			"id": 30900,
			"enchant": 29535,
			"gems": [
				32194,
				32194,
				32194
			]
		},
		{
			"id": 32366,
			"enchant": 22544,
			"gems": [
				32194,
				32194
			]
		},
		{
			"id": 28791
		},
		{
			"id": 29301
		},
		{
			"id": 28830
		},
		{
			"id": 32505
		},
		{
			"id": 30881,
			"enchant": 33165
		},
		{
			"id": 30881,
			"enchant": 33165
		},
		{
			"id": 30906,
			"enchant": 23766
		}
	]}`),
};

export const P4_BM_PRESET = {
	name: 'P4 BM Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<any>) => player.getTalentTree() != 2,
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
			"id": 32260
		},
		{
			"id": 31006,
			"enchant": 28888,
			"gems": [
				32222,
				32212
			]
		},
		{
			"id": 32323,
			"enchant": 34004
		},
		{
			"id": 31004,
			"enchant": 24003,
			"gems": [
				32194,
				32222,
				32226
			]
		},
		{
			"id": 32324,
			"enchant": 34002,
			"gems": [
				32222
			]
		},
		{
			"id": 31001,
			"enchant": 19445,
			"gems": [
				32194
			]
		},
		{
			"id": 32346
		},
		{
			"id": 31005,
			"enchant": 29535,
			"gems": [
				32194
			]
		},
		{
			"id": 32366,
			"enchant": 22544,
			"gems": [
				32194,
				32222
			]
		},
		{
			"id": 29301
		},
		{
			"id": 33496
		},
		{
			"id": 28830
		},
		{
			"id": 33831
		},
		{
			"id": 33389,
			"enchant": 33165
		},
		{
			"id": 33389,
			"enchant": 33165
		},
		{
			"id": 30906,
			"enchant": 23766
		}
	]}`),
};

export const P4_SV_PRESET = {
	name: 'P4 SV Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<any>) => player.getTalentTree() == 2,
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
			"id": 30017
		},
		{
			"id": 31006,
			"enchant": 28888,
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
			"id": 31004,
			"enchant": 24003,
			"gems": [
				32194,
				32226,
				32226
			]
		},
		{
			"id": 32324,
			"enchant": 34002,
			"gems": [
				32194
			]
		},
		{
			"id": 31001,
			"enchant": 19445,
			"gems": [
				32194
			]
		},
		{
			"id": 30879,
			"gems": [
				32194,
				32194
			]
		},
		{
			"id": 31005,
			"enchant": 29535,
			"gems": [
				32194
			]
		},
		{
			"id": 32366,
			"enchant": 22544,
			"gems": [
				32194,
				32194
			]
		},
		{
			"id": 33496
		},
		{
			"id": 29301
		},
		{
			"id": 33831
		},
		{
			"id": 28830
		},
		{
			"id": 33389,
			"enchant": 33165
		},
		{
			"id": 33389,
			"enchant": 33165
		},
		{
			"id": 30906,
			"enchant": 23766
		}
	]}`),
};

export const P5_BM_PRESET = {
	name: 'P5 BM Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<any>) => player.getTalentTree() != 2,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 34333,
			"enchant": 29192,
			"gems": [
				33131,
				32409
			]
		},
		{
			"id": 34358,
			"gems": [
				32205
			]
		},
		{
			"id": 31006,
			"enchant": 23548,
			"gems": [
				32205,
				32212
			]
		},
		{
			"id": 34241,
			"enchant": 34004,
			"gems": [
				32220
			]
		},
		{
			"id": 34397,
			"enchant": 24003,
			"gems": [
				32212,
				33143,
				32194
			]
		},
		{
			"id": 34443,
			"enchant": 34002,
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
			"id": 34549,
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
			"id": 34570,
			"enchant": 22544,
			"gems": [
				32220
			]
		},
		{
			"id": 34189
		},
		{
			"id": 34361
		},
		{
			"id": 34427
		},
		{
			"id": 33831
		},
		{
			"id": 34329,
			"enchant": 33165,
			"gems": [
				32194
			]
		},
		{
			"id": 34329,
			"enchant": 33165,
			"gems": [
				32194
			]
		},
		{
			"id": 34334,
			"enchant": 23766
		}
	]}`),
};

export const P5_SV_PRESET = {
	name: 'P5 SV Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<any>) => player.getTalentTree() == 2,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 34333,
			"enchant": 29192,
			"gems": [
				32194,
				32409
			]
		},
		{
			"id": 34177
		},
		{
			"id": 31006,
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
				32194,
				32194,
				32194
			]
		},
		{
			"id": 34443,
			"enchant": 34002,
			"gems": [
				32194
			]
		},
		{
			"id": 34343,
			"enchant": 19445,
			"gems": [
				32194,
				32226
			]
		},
		{
			"id": 34549,
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
			"id": 34570,
			"enchant": 22544,
			"gems": [
				32226
			]
		},
		{
			"id": 34887
		},
		{
			"id": 34361
		},
		{
			"id": 34427
		},
		{
			"id": 28830
		},
		{
			"id": 34183,
			"enchant": 22556,
			"gems": [
				32194
			]
		},
		{
			"id": 34334,
			"enchant": 23766
		}
	]}`),
};
