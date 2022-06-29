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

import { EnhancementShaman, EnhancementShaman_Rotation as EnhancementShamanRotation, EnhancementShaman_Options as EnhancementShamanOptions } from '/tbc/core/proto/shaman.js';
import {
	AirTotem,
	EarthTotem,
	FireTotem,
	WaterTotem,
	EnhancementShaman_Rotation_PrimaryShock as PrimaryShock,
	ShamanTotems,
} from '/tbc/core/proto/shaman.js';

import * as Enchants from '/tbc/core/constants/enchants.js';
import * as Gems from '/tbc/core/proto_utils/gems.js';
import * as Tooltips from '/tbc/core/constants/tooltips.js';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

// Default talents. Uses the wowhead calculator format, make the talents on
// https://tbc.wowhead.com/talent-calc and copy the numbers in the url.
export const StandardTalents = {
	name: 'Ele Sub',
	data: '250030502-502500210501133531151',
};

export const RestoSubspecTalents = {
	name: 'Resto Sub',
	data: '02-502500210502133531151-05005301',
};

export const DefaultRotation = EnhancementShamanRotation.create({
	totems: ShamanTotems.create({
		earth: EarthTotem.StrengthOfEarthTotem,
		air: AirTotem.GraceOfAirTotem,
		fire: FireTotem.SearingTotem,
		water: WaterTotem.ManaSpringTotem,
		twistWindfury: true,
		windfuryTotemRank: 5,
	}),
	primaryShock: PrimaryShock.Frost,
	weaveFlameShock: true,
});

export const DefaultOptions = EnhancementShamanOptions.create({
	waterShield: true,
	bloodlust: true,
	delayOffhandSwings: true,
});

export const DefaultConsumes = Consumes.create({
	drums: Drums.DrumsOfBattle,
	defaultPotion: Potions.HastePotion,
	flask: Flask.FlaskOfRelentlessAssault,
	food: Food.FoodRoastedClefthoof,
	mainHandImbue: WeaponImbue.WeaponImbueShamanWindfury,
	offHandImbue: WeaponImbue.WeaponImbueShamanWindfury,
	scrollOfAgility: 5,
	scrollOfStrength: 5,
});

export const P1_PRESET = {
	name: 'P1 Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 29040,
			"enchant": 29192,
			"gems": [
				32409,
				24027
			]
		},
		{
			"id": 29381
		},
		{
			"id": 29043,
			"enchant": 28888,
			"gems": [
				24027,
				24058
			]
		},
		{
			"id": 24259,
			"enchant": 34004,
			"gems": [
				24027
			]
		},
		{
			"id": 29038,
			"enchant": 24003,
			"gems": [
				24027,
				24027,
				24058
			]
		},
		{
			"id": 25697,
			"enchant": 27899,
			"gems": [
				24027
			]
		},
		{
			"id": 29039,
			"enchant": 33995
		},
		{
			"id": 28656
		},
		{
			"id": 30534,
			"enchant": 29535,
			"gems": [
				24054,
				24054,
				24058
			]
		},
		{
			"id": 28746,
			"enchant": 28279,
			"gems": [
				24027,
				24027
			]
		},
		{
			"id": 28757
		},
		{
			"id": 29283
		},
		{
			"id": 28830
		},
		{
			"id": 29383
		},
		{
			"id": 28767,
			"enchant": 22559
		},
		{
			"id": 27872,
			"enchant": 22559
		},
		{
			"id": 27815
		}
	]}`),
};

export const P2_PRESET = {
	name: 'P2 Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 30190,
			"enchant": 29192,
			"gems": [
				32409,
				24058
			]
		},
		{
			"id": 30017
		},
		{
			"id": 30055,
			"enchant": 28888,
			"gems": [
				24027
			]
		},
		{
			"id": 29994,
			"enchant": 34004
		},
		{
			"id": 30185,
			"enchant": 24003,
			"gems": [
				24027,
				24054,
				24058
			]
		},
		{
			"id": 30091,
			"enchant": 27899,
			"gems": [
				24027
			]
		},
		{
			"id": 30189,
			"enchant": 33995
		},
		{
			"id": 30106,
			"gems": [
				24027,
				24054
			]
		},
		{
			"id": 30192,
			"enchant": 29535,
			"gems": [
				24027
			]
		},
		{
			"id": 30039,
			"enchant": 28279
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
			"id": 29383
		},
		{
			"id": 32944,
			"enchant": 22559
		},
		{
			"id": 29996,
			"enchant": 22559
		},
		{
			"id": 27815
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
				32193
			]
		},
		{
			"id": 32260
		},
		{
			"id": 32575,
			"enchant": 28888
		},
		{
			"id": 32323,
			"enchant": 34004
		},
		{
			"id": 30905,
			"enchant": 24003,
			"gems": [
				32211,
				32193,
				32217
			]
		},
		{
			"id": 32574,
			"enchant": 27899
		},
		{
			"id": 32234,
			"enchant": 33995
		},
		{
			"id": 30106,
			"gems": [
				32193,
				32193
			]
		},
		{
			"id": 30900,
			"enchant": 29535,
			"gems": [
				32193,
				32217,
				32211
			]
		},
		{
			"id": 32510,
			"enchant": 28279
		},
		{
			"id": 29997
		},
		{
			"id": 32497
		},
		{
			"id": 28830
		},
		{
			"id": 29383
		},
		{
			"id": 32262,
			"enchant": 22559
		},
		{
			"id": 32262,
			"enchant": 22559
		},
		{
			"id": 27815
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
				32217
			]
		},
		{
			"id": 32260
		},
		{
			"id": 32581,
			"enchant": 28888
		},
		{
			"id": 32323,
			"enchant": 34004
		},
		{
			"id": 30905,
			"enchant": 24003,
			"gems": [
				32211,
				32193,
				32217
			]
		},
		{
			"id": 30863,
			"enchant": 27899,
			"gems": [
				32217
			]
		},
		{
			"id": 32234,
			"enchant": 33995
		},
		{
			"id": 30106,
			"gems": [
				32193,
				32211
			]
		},
		{
			"id": 30900,
			"enchant": 29535,
			"gems": [
				32193,
				32217,
				32211
			]
		},
		{
			"id": 32510,
			"enchant": 28279
		},
		{
			"id": 32497
		},
		{
			"id": 33496
		},
		{
			"id": 33831
		},
		{
			"id": 28830
		},
		{
			"id": 32262,
			"enchant": 22559
		},
		{
			"id": 32262,
			"enchant": 33307
		},
		{
			"id": 33507
		}
	]}`),
};

export const P5_PRESET = {
	name: 'P5 Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
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
			"id": 31024,
			"enchant": 28888,
			"gems": [
				32217,
				32217
			]
		},
		{
			"id": 34241,
			"enchant": 34004,
			"gems": [
				32193
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
			"id": 34439,
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
			"id": 34545,
			"gems": [
				32193
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
			"id": 34567,
			"enchant": 28279,
			"gems": [
				32217
			]
		},
		{
			"id": 34189
		},
		{
			"id": 32497
		},
		{
			"id": 34427
		},
		{
			"id": 34472
		},
		{
			"id": 34331,
			"enchant": 22559,
			"gems": [
				32217,
				32217
			]
		},
		{
			"id": 34346,
			"enchant": 33307,
			"gems": [
				32217,
				32211
			]
		},
		{
			"id": 33507
		}
	]}`),
};
