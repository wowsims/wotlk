import { Consumes } from '/tbc/core/proto/common.js';
import { EquipmentSpec } from '/tbc/core/proto/common.js';
import { Flask } from '/tbc/core/proto/common.js';
import { Food } from '/tbc/core/proto/common.js';
import { ItemSpec } from '/tbc/core/proto/common.js';
import { Potions } from '/tbc/core/proto/common.js';
import { WeaponImbue } from '/tbc/core/proto/common.js';
import { Faction } from '/tbc/core/proto/common.js';
import { Player } from '/tbc/core/player.js';

import { SmitePriest, SmitePriest_Rotation as Rotation, SmitePriest_Options as Options, SmitePriest_Rotation_RotationType } from '/tbc/core/proto/priest.js';

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
	data: '5012300130505120501-005551002020052',
};

export const HolyTalents = {
	name: 'Holy',
	data: '50023011305-235050032002150520051',
};

export const DefaultRotation = Rotation.create({
	rotationType: SmitePriest_Rotation_RotationType.Basic,
});

export const DefaultOptions = Options.create({
	useShadowfiend: true,
});

export const DefaultConsumes = Consumes.create({
	flask: Flask.FlaskOfBlindingLight,
	food: Food.FoodBlackenedBasilisk,
	mainHandImbue: WeaponImbue.WeaponImbueSuperiorWizardOil,
	defaultPotion: Potions.SuperManaPotion,
});

export const P1_PRESET = {
	name: 'P1 Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	gear: EquipmentSpec.fromJsonString(`{
      "items": [
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
          "id": 28530
        },
        {
          "id": 29060,
          "enchant": 28886,
          "gems": [
            24030,
            24030
          ]
        },
        {
          "id": 28766,
          "enchant": 33150
        },
        {
          "id": 29056,
          "enchant": 24003,
          "gems": [
            24030,
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
          "id": 30725,
          "enchant": 28272,
          "gems": [
            24030,
            24030
          ]
        },
        {
          "id": 24256,
          "gems": [
            24030,
            24030
          ]
        },
        {
          "id": 30734,
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
          "id": 28793,
          "enchant": 22536
        },
        {
          "id": 29172,
          "enchant": 22536
        },
        {
          "id": 27683
        },
        {
          "id": 29370
        },
        {
          "id": 30723,
          "enchant": 22555,
          "gems": [
            30564,
            31867
          ]
        },
        {
          "id": 28734
        },
        {
          "id": 28673
        }
      ]
    }`),
};

export const P2_PRESET = {
	name: 'P2 Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	gear: EquipmentSpec.fromJsonString(`{
      "items": [
        {
          "id": 30161,
          "enchant": 29191,
          "gems": [
            34220,
            30588
          ]
        },
        {
          "id": 30015
        },
        {
          "id": 30163,
          "enchant": 28886,
          "gems": [
            24059,
            30600
          ]
        },
        {
          "id": 28766,
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
          "id": 29918,
          "enchant": 22534
        },
        {
          "id": 30160,
          "enchant": 28272
        },
        {
          "id": 30038,
          "gems": [
            31116,
            24059
          ]
        },
        {
          "id": 30162,
          "enchant": 24274,
          "gems": [
            24030
          ]
        },
        {
          "id": 30037,
          "enchant": 35297
        },
        {
          "id": 30109,
          "enchant": 22536
        },
        {
          "id": 28793,
          "enchant": 22536
        },
        {
          "id": 27683
        },
        {
          "id": 29370
        },
        {
          "id": 30723,
          "enchant": 22555,
          "gems": [
            24030,
            24030
          ]
        },
        {
          "id": 30049
        },
        {
          "id": 29982
        }
      ]
    }`),
};

export const P3_PRESET = {
	name: 'P3 Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	gear: EquipmentSpec.fromJsonString(`{
      "items": [
        {
          "id": 32525,
          "enchant": 29191,
          "gems": [
            34220,
            30600
          ]
        },
        {
          "id": 32349
        },
        {
          "id": 31070,
          "enchant": 28886,
          "gems": [
            32218,
            32215
          ]
        },
        {
          "id": 32524,
          "enchant": 33150
        },
        {
          "id": 30107,
          "enchant": 33990,
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
          "id": 30038,
          "gems": [
            32196,
            32196
          ]
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
          "id": 32247,
          "enchant": 22536
        },
        {
          "id": 29370
        },
        {
          "id": 32483
        },
        {
          "id": 32374,
          "enchant": 22555
        },
        {},
        {
          "id": 29982
        }
      ]
    }`),

};

export const P4_PRESET = {
	name: 'P4 Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	gear: EquipmentSpec.fromJsonString(`{
      "items": [
        {
          "id": 32525,
          "enchant": 29191,
          "gems": [
            34220,
            30600
          ]
        },
        {
          "id": 33281
        },
        {
          "id": 31070,
          "enchant": 28886,
          "gems": [
            32210,
            32215
          ]
        },
        {
          "id": 32524,
          "enchant": 33150
        },
        {
          "id": 30107,
          "enchant": 33990,
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
            32221
          ]
        },
        {
          "id": 30038,
          "gems": [
            32196,
            32196
          ]
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
          "id": 33497,
          "enchant": 22536
        },
        {
          "id": 33829
        },
        {
          "id": 32483
        },
        {
          "id": 32374,
          "enchant": 22555
        },
        {},
        {
          "id": 29982
        }
      ]
    }`),

};

export const P5_PRESET = {
	name: 'P5 Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	gear: EquipmentSpec.fromJsonString(`{
      "items": [
        {
          "id": 34340,
          "enchant": 29191,
          "gems": [
            34220,
            30600
          ]
        },
        {
          "id": 34204
        },
        {
          "id": 34210,
          "enchant": 28886,
          "gems": [
            32221,
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
          "id": 34364,
          "enchant": 33990,
          "gems": [
            35760,
            35760,
            35760
          ]
        },
        {
          "id": 34434,
          "enchant": 22534,
          "gems": [
            32221
          ]
        },
        {
          "id": 34344,
          "enchant": 28272,
          "gems": [
            35761,
            32221
          ]
        },
        {
          "id": 34528,
          "gems": [
            35761
          ]
        },
        {
          "id": 34181,
          "enchant": 24274,
          "gems": [
            32221,
            32221,
            35761
          ]
        },
        {
          "id": 34563,
          "enchant": 35297,
          "gems": [
            32215
          ]
        },
        {
          "id": 34362,
          "enchant": 22536
        },
        {
          "id": 34889,
          "enchant": 22536
        },
        {
          "id": 34429
        },
        {
          "id": 35749
        },
        {
          "id": 34182,
          "enchant": 22555,
          "gems": [
            35761,
            35761,
            35761
          ]
        },
        {},
        {
          "id": 29982
        }
      ]
    }`),

};
