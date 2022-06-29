import { Consumes } from '/tbc/core/proto/common.js';
import { Drums } from '/tbc/core/proto/common.js';
import { EquipmentSpec } from '/tbc/core/proto/common.js';
import { Flask } from '/tbc/core/proto/common.js';
import { Food } from '/tbc/core/proto/common.js';
import { ItemSpec } from '/tbc/core/proto/common.js';
import { Potions } from '/tbc/core/proto/common.js';
import { Faction } from '/tbc/core/proto/common.js';
import { WeaponImbue } from '/tbc/core/proto/common.js';
import { Player } from '/tbc/core/player.js';

import { ElementalShaman, ElementalShaman_Rotation as ElementalShamanRotation, ElementalShaman_Options as ElementalShamanOptions } from '/tbc/core/proto/shaman.js';
import { ElementalShaman_Rotation_RotationType as RotationType } from '/tbc/core/proto/shaman.js';

import {
	AirTotem,
	EarthTotem,
	FireTotem,
	WaterTotem,
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
	name: 'Standard',
	data: '55003105100213351051--05105301005',
};

export const RestoTalents = {
	name: 'Resto',
	data: '5003--55035051355310510321',
};

export const DefaultRotation = ElementalShamanRotation.create({
	totems: ShamanTotems.create({
		earth: EarthTotem.TremorTotem,
		air: AirTotem.WrathOfAirTotem,
		fire: FireTotem.TotemOfWrath,
		water: WaterTotem.ManaSpringTotem,
	}),
	type: RotationType.Adaptive,
});

export const DefaultOptions = ElementalShamanOptions.create({
	waterShield: true,
	bloodlust: true,
});

export const DefaultConsumes = Consumes.create({
	drums: Drums.DrumsOfBattle,
	defaultPotion: Potions.SuperManaPotion,
	flask: Flask.FlaskOfBlindingLight,
	food: Food.FoodBlackenedBasilisk,
	mainHandImbue: WeaponImbue.WeaponImbueBrilliantWizardOil,
});

export const P1_PRESET = {
	name: 'P1 Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 29035,
			"enchant": 29191,
			"gems": [
				34220,
				24059
			]
		},
		{
			"id": 28762
		},
		{
			"id": 29037,
			"enchant": 28886,
			"gems": [
				24059,
				24059
			]
		},
		{
			"id": 28797,
			"enchant": 33150
		},
		{
			"id": 29519,
			"enchant": 24003,
			"gems": [
				24030,
				24030,
				24030
			]
		},
		{
			"id": 29521,
			"enchant": 22534,
			"gems": [
				24059
			]
		},
		{
			"id": 28780,
			"enchant": 28272,
			"gems": [
				24059,
				24056
			]
		},
		{
			"id": 29520,
			"gems": [
				24056,
				24059
			]
		},
		{
			"id": 24262,
			"enchant": 24274,
			"gems": [
				24030,
				24030,
				24030
			]
		},
		{
			"id": 28517,
			"enchant": 35297,
			"gems": [
				24030,
				24030
			]
		},
		{
			"id": 30667,
			"enchant": 22536
		},
		{
			"id": 28753,
			"enchant": 22536
		},
		{
			"id": 29370
		},
		{
			"id": 28785
		},
		{
			"id": 28770,
			"enchant": 22555
		},
		{
			"id": 29273
		},
		{
			"id": 28248
		}
	]}`),
};

export const P2_PRESET = {
	name: 'P2 Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 29035,
			"enchant": 29191,
			"gems": [
				34220,
				24059
			]
		},
		{
			"id": 30015
		},
		{
			"id": 29037,
			"enchant": 28886,
			"gems": [
				24059,
				24059
			]
		},
		{
			"id": 28797,
			"enchant": 33150
		},
		{
			"id": 30169,
			"enchant": 24003,
			"gems": [
				24030,
				24030,
				24030
			]
		},
		{
			"id": 29918,
			"enchant": 22534
		},
		{
			"id": 28780,
			"enchant": 28272,
			"gems": [
				24059,
				24056
			]
		},
		{
			"id": 30038,
			"gems": [
				24056,
				24059
			]
		},
		{
			"id": 30172,
			"enchant": 24274,
			"gems": [
				24059
			]
		},
		{
			"id": 30067,
			"enchant": 35297
		},
		{
			"id": 30667,
			"enchant": 22536
		},
		{
			"id": 30109,
			"enchant": 22536
		},
		{
			"id": 29370
		},
		{
			"id": 28785
		},
		{
			"id": 29988,
			"enchant": 22555
		},
		{
			"id": 28248
		}
	]}`),
};

export const P3_PRESET = {
	name: 'P3 Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 31014,
			"enchant": 29191,
			"gems": [
				34220,
				32215
			]
		},
		{
			"id": 30015
		},
		{
			"id": 31023,
			"enchant": 28886,
			"gems": [
				32215,
				32218
			]
		},
		{
			"id": 32331,
			"enchant": 33150
		},
		{
			"id": 31017,
			"enchant": 24003,
			"gems": [
				32196,
				32196,
				32196
			]
		},
		{
			"id": 32586,
			"enchant": 22534
		},
		{
			"id": 31008,
			"enchant": 28272,
			"gems": [
				32218
			]
		},
		{
			"id": 32276
		},
		{
			"id": 30916,
			"enchant": 24274,
			"gems": [
				32196,
				32196,
				32196
			]
		},
		{
			"id": 32352,
			"enchant": 35297,
			"gems": [
				32196,
				32196
			]
		},
		{
			"id": 32527,
			"enchant": 22536
		},
		{
			"id": 29305,
			"enchant": 22536
		},
		{
			"id": 32483
		},
		{
			"id": 28785
		},
		{
			"id": 32374,
			"enchant": 22555
		},
		{
			"id": 32330
		}
	]}`),
};

export const P4_PRESET = {
	name: 'P4 Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 31014,
			"enchant": 29191,
			"gems": [
				34220,
				32215
			]
		},
		{
			"id": 33281
		},
		{
			"id": 31023,
			"enchant": 28886,
			"gems": [
				32215,
				32218
			]
		},
		{
			"id": 32331,
			"enchant": 33150
		},
		{
			"id": 31017,
			"enchant": 24003,
			"gems": [
				32196,
				32196,
				32196
			]
		},
		{
			"id": 32586,
			"enchant": 22534
		},
		{
			"id": 31008,
			"enchant": 28272,
			"gems": [
				32218
			]
		},
		{
			"id": 32276
		},
		{
			"id": 30916,
			"enchant": 24274,
			"gems": [
				32196,
				32196,
				32196
			]
		},
		{
			"id": 32352,
			"enchant": 35297,
			"gems": [
				32196,
				32196
			]
		},
		{
			"id": 32527,
			"enchant": 22536
		},
		{
			"id": 32527,
			"enchant": 22536
		},
		{
			"id": 32483
		},
		{
			"id": 33829
		},
		{
			"id": 32374,
			"enchant": 22555
		},
		{},
		{
			"id": 32330
		}
	]}`),
};

export const P5_ALLIANCE_PRESET = {
	name: 'P5 Alliance Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<any>) => player.getFaction() == Faction.Alliance,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 34332,
			"enchant": 29191,
			"gems": [
				35761,
				34220
			]
		},
		{
			"id": 34204
		},
		{
			"id": 31023,
			"enchant": 23545,
			"gems": [
				32215,
				35761
			]
		},
		{
			"id": 34242,
			"enchant": 33150,
			"gems": [
				35760
			]
		},
		{
			"id": 34396,
			"enchant": 24003,
			"gems": [
				35760,
				35761,
				35761
			]
		},
		{
			"id": 34437,
			"enchant": 22534,
			"gems": [
				35761
			]
		},
		{
			"id": 34350,
			"enchant": 28272,
			"gems": [
				35760,
				32215
			]
		},
		{
			"id": 34542,
			"gems": [
				35761
			]
		},
		{
			"id": 34186,
			"enchant": 24274,
			"gems": [
				35761,
				35760,
				35760
			]
		},
		{
			"id": 34566,
			"enchant": 35297,
			"gems": [
				35760
			]
		},
		{
			"id": 34230,
			"enchant": 22536
		},
		{
			"id": 32527,
			"enchant": 22536
		},
		{
			"id": 34429
		},
		{
			"id": 33829
		},
		{
			"id": 34336,
			"enchant": 22555
		},
		{
			"id": 34179
		},
		{
			"id": 32330
		}
	]}`),
};

export const P5_HORDE_PRESET = {
	name: 'P5 Horde Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<any>) => player.getFaction() == Faction.Horde,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 34332,
			"enchant": 29191,
			"gems": [
				35761,
				34220
			]
		},
		{
			"id": 34204
		},
		{
			"id": 31023,
			"enchant": 23545,
			"gems": [
				32215,
				35761
			]
		},
		{
			"id": 34242,
			"enchant": 33150,
			"gems": [
				35760
			]
		},
		{
			"id": 34396,
			"enchant": 24003,
			"gems": [
				35760,
				35761,
				35761
			]
		},
		{
			"id": 34437,
			"enchant": 22534,
			"gems": [
				35761
			]
		},
		{
			"id": 34350,
			"enchant": 28272,
			"gems": [
				35760,
				32215
			]
		},
		{
			"id": 34542,
			"gems": [
				35761
			]
		},
		{
			"id": 34186,
			"enchant": 24274,
			"gems": [
				35761,
				35760,
				35760
			]
		},
		{
			"id": 34566,
			"enchant": 35297,
			"gems": [
				35760
			]
		},
		{
			"id": 34230,
			"enchant": 22536
		},
		{
			"id": 32527,
			"enchant": 22536
		},
		{
			"id": 34429
		},
		{
			"id": 32483
		},
		{
			"id": 34336,
			"enchant": 22555
		},
		{
			"id": 34179
		},
		{
			"id": 32330
		}
	]}`),
};
