import { Consumes } from '/tbc/core/proto/common.js';
import { EquipmentSpec } from '/tbc/core/proto/common.js';
import { Flask } from '/tbc/core/proto/common.js';
import { Food } from '/tbc/core/proto/common.js';
import { ItemSpec } from '/tbc/core/proto/common.js';
import { Potions } from '/tbc/core/proto/common.js';
import { WeaponImbue } from '/tbc/core/proto/common.js';
import { Faction } from '/tbc/core/proto/common.js';
import { Player } from '/tbc/core/player.js';

import { ShadowPriest, ShadowPriest_Rotation as Rotation, ShadowPriest_Options as Options, ShadowPriest_Rotation_RotationType } from '/tbc/core/proto/priest.js';

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
	data: '500230013--503250510240103051451',
};

export const DefaultRotation = Rotation.create({
	rotationType: ShadowPriest_Rotation_RotationType.Ideal,
});

export const DefaultOptions = Options.create({
	useShadowfiend: true,
});

export const DefaultConsumes = Consumes.create({
	flask: Flask.FlaskOfPureDeath,
	food: Food.FoodBlackenedBasilisk,
	mainHandImbue: WeaponImbue.WeaponImbueSuperiorWizardOil,
	defaultPotion: Potions.SuperManaPotion,
});

export const P1_PRESET = {
	name: 'P1 Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 24266,
			"enchant": 29191,
			"gems": [
				28118,
				24030,
				24030
			]
		},
		{
			"id": 30666
		},
		{
			"id": 21869,
			"enchant": 28886,
			"gems": [
				24030,
				24030
			]
		},
		{
			"id": 28570,
			"enchant": 33150
		},
		{
			"id": 21871,
			"enchant": 24003,
			"gems": [
				24030,
				24030
			]
		},
		{
			"id": 24250,
			"enchant": 22534,
			"gems": [
				24030
			]
		},
		{
			"id": 28507,
			"enchant": 28272,
			"gems": [
				24030,
				24030
			]
		},
		{
			"id": 28799,
			"gems": [
				24030,
				24030
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
			"id": 21870,
			"enchant": 35297,
			"gems": [
				24030,
				24030
			]
		},
		{
			"id": 29352,
			"enchant": 22536
		},
		{
			"id": 28793,
			"enchant": 22536
		},
		{
			"id": 28789
		},
		{
			"id": 29370
		},
		{
			"id": 28770,
			"enchant": 22561
		},
		{
			"id": 29272
		},
		{
			"id": 29350
		}
	]}`),
};

export const P2_PRESET = {
	name: 'P2 Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 32494,
			"enchant": 29191,
			"gems": [
				25893,
				24056
			]
		},
		{
			"id": 30666
		},
		{
			"id": 30163,
			"enchant": 28886,
			"gems": [
				24030,
				24030
			]
		},
		{
			"id": 29992,
			"enchant": 33150
		},
		{
			"id": 30107,
			"enchant": 24003,
			"gems": [
				24030,
				24030,
				24030
			]
		},
		{
			"id": -19,
			"enchant": 22534
		},
		{
			"id": 28780,
			"enchant": 28272,
			"gems": [
				24030,
				24030
			]
		},
		{
			"id": 30038,
			"gems": [
				24030,
				24030
			]
		},
		{
			"id": 29972,
			"enchant": 24274,
			"gems": [
				24030,
				24030,
				24030
			]
		},
		{
			"id": 21870,
			"enchant": 35297,
			"gems": [
				24030,
				24030
			]
		},
		{
			"id": 30109,
			"enchant": 22536
		},
		{
			"id": 29922,
			"enchant": 22536
		},
		{
			"id": 29370
		},
		{
			"id": 38290
		},
		{
			"id": 28770,
			"enchant": 22561
		},
		{
			"id": 29272
		},
		{
			"id": 29982
		}
	]}`),
};

export const P3_PRESET = {
	name: 'P3 Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 31064,
			"enchant": 29191,
			"gems": [
				25893,
				32215
			]
		},
		{
			"id": 30666
		},
		{
			"id": 31070,
			"enchant": 28886,
			"gems": [
				32196,
				32196
			]
		},
		{
			"id": 32590,
			"enchant": 33150
		},
		{
			"id": 31065,
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
			"id": 31061,
			"enchant": 28272,
			"gems": [
				32196
			]
		},
		{
			"id": 32256
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
			"id": 32239,
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
			"id": 29370
		},
		{
			"id": 32374,
			"enchant": 22561
		},
		{
			"id": 29982
		}
	]}`),
};

export const P4_PRESET = {
	name: 'P4 Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 31064,
			"enchant": 29191,
			"gems": [
				25893,
				32215
			]
		},
		{
			"id": 33466
		},
		{
			"id": 31070,
			"enchant": 28886,
			"gems": [
				32196,
				32196
			]
		},
		{
			"id": 32590,
			"enchant": 33150
		},
		{
			"id": 31065,
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
			"id": 31061,
			"enchant": 28272,
			"gems": [
				32196
			]
		},
		{
			"id": 32256
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
			"id": 32239,
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
			"enchant": 22561
		},
		{
			"id": 33192,
			"gems": [
				32196
			]
		}
	]}`),
};

export const P5_PRESET = {
	name: 'P5 Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 34340,
			"enchant": 29191,
			"gems": [
				25893,
				32215
			]
		},
		{
			"id": 34204
		},
		{
			"id": 31070,
			"enchant": 28886,
			"gems": [
				32196,
				32196
			]
		},
		{
			"id": 34242,
			"enchant": 33150,
			"gems": [
				32196
			]
		},
		{
			"id": 34232,
			"enchant": 33990,
			"gems": [
				32196,
				32196,
				32196
			]
		},
		{
			"id": 34434,
			"enchant": 22534,
			"gems": [
				32196
			]
		},
		{
			"id": 34344,
			"enchant": 28272,
			"gems": [
				32196,
				32196
			]
		},
		{
			"id": 34528,
			"gems": [
				32196
			]
		},
		{
			"id": 34181,
			"enchant": 24274,
			"gems": [
				32196,
				32196,
				32196
			]
		},
		{
			"id": 34563,
			"enchant": 35297,
			"gems": [
				32196
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
			"id": 33829
		},
		{
			"id": 34429
		},
		{
			"id": 34336,
			"enchant": 22561
		},
		{
			"id": 34179
		},
		{
			"id": 34347,
			"gems": [
				32196
			]
		}
	]}`),
};
