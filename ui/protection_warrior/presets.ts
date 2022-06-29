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
	ProtectionWarrior,
	ProtectionWarrior_Rotation as ProtectionWarriorRotation,
	ProtectionWarrior_Rotation_DemoShout as DemoShout,
	ProtectionWarrior_Rotation_ThunderClap as ThunderClap,
	ProtectionWarrior_Options as ProtectionWarriorOptions
} from '/tbc/core/proto/warrior.js';

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
	data: '350003011-05-0055511033010103501351',
};
export const ImpDemoTalents = {
	name: 'Imp Demo',
	data: '340003-055-0055511033010101501351',
};
export const ImpaleProtTalents = {
	name: 'Impale Prot',
	data: '35000301302-03-0055511033010101501351',
};

export const DefaultRotation = ProtectionWarriorRotation.create({
	demoShout: DemoShout.DemoShoutMaintain,
	thunderClap: ThunderClap.ThunderClapMaintain,
	hsRageThreshold: 30,
});

export const DefaultOptions = ProtectionWarriorOptions.create({
	shout: WarriorShout.WarriorShoutCommanding,
	precastShout: true,
	precastShoutSapphire: false,
	precastShoutT2: false,

	startingRage: 0,
});

export const DefaultConsumes = Consumes.create({
	flask: Flask.FlaskOfFortification,
	food: Food.FoodFishermansFeast,
	defaultPotion: Potions.IronshieldPotion,
	mainHandImbue: WeaponImbue.WeaponImbueAdamantiteSharpeningStone,
	offHandImbue: WeaponImbue.WeaponImbueAdamantiteSharpeningStone,
	scrollOfAgility: 5,
	scrollOfStrength: 5,
	scrollOfProtection: 5,
});

export const P1_BALANCED_PRESET = {
	name: 'P1 Balanced Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 29011,
			"enchant": 29192,
			"gems": [
				25896,
				24033
			]
		},
		{
			"id": 28244,
			"gems": [
				33782
			]
		},
		{
			"id": 29023,
			"enchant": 28911,
			"gems": [
				24033,
				24033
			]
		},
		{
			"id": 28672,
			"enchant": 34004
		},
		{
			"id": 29012,
			"enchant": 24003,
			"gems": [
				24033,
				24033,
				24033
			]
		},
		{
			"id": 28996,
			"enchant": 22533,
			"gems": [
				33782
			]
		},
		{
			"id": 30644,
			"enchant": 33153
		},
		{
			"id": 28995
		},
		{
			"id": 28621,
			"enchant": 29536,
			"gems": [
				24033,
				24033,
				24033
			]
		},
		{
			"id": 28747,
			"enchant": 35297,
			"gems": [
				24033,
				24033
			]
		},
		{
			"id": 30834
		},
		{
			"id": 29279
		},
		{
			"id": 28121
		},
		{
			"id": 29387
		},
		{
			"id": 28749,
			"enchant": 22559
		},
		{
			"id": 28825,
			"enchant": 28282,
			"gems": [
				24033
			]
		},
		{
			"id": 28826
		}
	]}`),
};

export const P2_BALANCED_PRESET = {
	name: 'P2 Balanced Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 30115,
			"enchant": 29192,
			"gems": [
				25896,
				24033
			]
		},
		{
			"id": 33066,
			"gems": [
				33782
			]
		},
		{
			"id": 30117,
			"enchant": 28910,
			"gems": [
				24033,
				24033
			]
		},
		{
			"id": 29994,
			"enchant": 34004
		},
		{
			"id": 30113,
			"enchant": 24003,
			"gems": [
				24033,
				24033,
				24033
			]
		},
		{
			"id": 32818,
			"enchant": 22533,
			"gems": [
				33782
			]
		},
		{
			"id": 29947,
			"enchant": 33153
		},
		{
			"id": 30106,
			"gems": [
				24033,
				24033
			]
		},
		{
			"id": 30116,
			"enchant": 29536,
			"gems": [
				24055
			]
		},
		{
			"id": 32793,
			"enchant": 35297
		},
		{
			"id": 30834
		},
		{
			"id": 29283
		},
		{
			"id": 28121
		},
		{
			"id": 37128
		},
		{
			"id": 30058,
			"enchant": 22559
		},
		{
			"id": 28825,
			"enchant": 28282,
			"gems": [
				24033
			]
		},
		{
			"id": 32756,
			"gems": [
				24033
			]
		}
	]}`),
};

export const P3_BALANCED_PRESET = {
	name: 'P3 Balanced Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 33730,
			"enchant": 29192,
			"gems": [
				25896,
				32220
			]
		},
		{
			"id": 33923,
			"gems": [
				32220
			]
		},
		{
			"id": 30979,
			"enchant": 28910,
			"gems": [
				24033,
				24033
			]
		},
		{
			"id": 34010,
			"enchant": 34004
		},
		{
			"id": 33728,
			"enchant": 24003,
			"gems": [
				32200,
				32200,
				32200
			]
		},
		{
			"id": 33813,
			"enchant": 22533,
			"gems": [
				32200
			]
		},
		{
			"id": 32280,
			"enchant": 33153
		},
		{
			"id": 30106,
			"gems": [
				32220,
				32200
			]
		},
		{
			"id": 30977,
			"enchant": 29536,
			"gems": [
				32200
			]
		},
		{
			"id": 33812,
			"enchant": 35297
		},
		{
			"id": 30834
		},
		{
			"id": 29301
		},
		{
			"id": 31858
		},
		{
			"id": 31859
		},
		{
			"id": 32254,
			"enchant": 22559
		},
		{
			"id": 32375,
			"enchant": 28282
		},
		{
			"id": 32253
		}
	]}`),
};

export const P4_BALANCED_PRESET = {
	name: 'P4 Balanced Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 33730,
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
			"id": 30979,
			"enchant": 28910,
			"gems": [
				32226,
				32212
			]
		},
		{
			"id": 33484,
			"enchant": 34004
		},
		{
			"id": 33728,
			"enchant": 24003,
			"gems": [
				32200,
				32200,
				32200
			]
		},
		{
			"id": 33516,
			"enchant": 22533
		},
		{
			"id": 32280,
			"enchant": 33153
		},
		{
			"id": 32333,
			"gems": [
				32212,
				32226
			]
		},
		{
			"id": 30978,
			"enchant": 29536,
			"gems": [
				32200
			]
		},
		{
			"id": 32268,
			"enchant": 35297,
			"gems": [
				32200,
				32200
			]
		},
		{
			"id": 33496
		},
		{
			"id": 29301
		},
		{
			"id": 31858
		},
		{
			"id": 31859
		},
		{
			"id": 32369,
			"enchant": 22559
		},
		{
			"id": 32375,
			"enchant": 28282
		},
		{
			"id": 32253
		}
	]}`),
};

export const P5_BALANCED_PRESET = {
	name: 'P5 Balanced Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 35068,
			"enchant": 29192,
			"gems": [
				32409,
				32205
			]
		},
		{
			"id": 34178
		},
		{
			"id": 34388,
			"enchant": 28910,
			"gems": [
				32217,
				32205
			]
		},
		{
			"id": 34190,
			"enchant": 34004
		},
		{
			"id": 35066,
			"enchant": 24003,
			"gems": [
				32217,
				32217,
				32205
			]
		},
		{
			"id": 34442,
			"enchant": 22533,
			"gems": [
				32200
			]
		},
		{
			"id": 34378,
			"enchant": 33153,
			"gems": [
				32200,
				32200
			]
		},
		{
			"id": 34547,
			"gems": [
				32212
			]
		},
		{
			"id": 34381,
			"enchant": 29536,
			"gems": [
				32200,
				32200,
				32212
			]
		},
		{
			"id": 34568,
			"enchant": 35297,
			"gems": [
				32212
			]
		},
		{
			"id": 34213
		},
		{
			"id": 32266
		},
		{
			"id": 31858
		},
		{
			"id": 34473
		},
		{
			"id": 34164,
			"enchant": 22559
		},
		{
			"id": 34185,
			"enchant": 28282,
			"gems": [
				32212
			]
		},
		{
			"id": 32253
		}
	]}`),
};
