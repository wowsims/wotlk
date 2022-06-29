import { Consumes } from '/tbc/core/proto/common.js';
import { Flask } from '/tbc/core/proto/common.js';
import { Food } from '/tbc/core/proto/common.js';
import { EquipmentSpec } from '/tbc/core/proto/common.js';
import { ItemSpec } from '/tbc/core/proto/common.js';
import { Potions } from '/tbc/core/proto/common.js';
import { WeaponImbue } from '/tbc/core/proto/common.js';
import { Faction } from '/tbc/core/proto/common.js';
import { Player } from '/tbc/core/player.js';

import { BalanceDruid, BalanceDruid_Rotation as BalanceDruidRotation, DruidTalents as DruidTalents, BalanceDruid_Options as BalanceDruidOptions } from '/tbc/core/proto/druid.js';
import { BalanceDruid_Rotation_PrimarySpell as PrimarySpell } from '/tbc/core/proto/druid.js';

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
	data: '510022312503135231351--520033',
};

export const DreamstateTalents = {
	name: 'Dreamstate',
	data: '5003223122031312303--500503400314',
};

export const RestoTalents = {
	name: 'Resto',
	data: '--50353351531522531351',
};

export const DefaultRotation = BalanceDruidRotation.create({
	primarySpell: PrimarySpell.Adaptive,
	faerieFire: true,
});

export const DefaultOptions = BalanceDruidOptions.create({
});

export const DefaultConsumes = Consumes.create({
	flask: Flask.FlaskOfBlindingLight,
	food: Food.FoodBlackenedBasilisk,
	mainHandImbue: WeaponImbue.WeaponImbueBrilliantWizardOil,
	defaultPotion: Potions.SuperManaPotion,
});

export const P1_ALLIANCE_PRESET = {
	name: 'P1 Alliance Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<any>) => player.getFaction() == Faction.Alliance,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 29093,
			"enchant": 29191,
			"gems": [
				24030,
				34220
			]
		},
		{
			"id": 28762
		},
		{
			"id": 29095,
			"enchant": 28886,
			"gems": [
				24056,
				31867
			]
		},
		{
			"id": 28766,
			"enchant": 33150
		},
		{
			"id": 21848,
			"enchant": 24003,
			"gems": [
				31867,
				24056
			]
		},
		{
			"id": 24250,
			"enchant": 22534,
			"gems": [
				31867
			]
		},
		{
			"id": 21847,
			"enchant": 28272,
			"gems": [
				31867,
				31867
			]
		},
		{
			"id": 21846,
			"gems": [
				31867,
				31867
			]
		},
		{
			"id": 24262,
			"enchant": 24274,
			"gems": [
				31867,
				31867,
				31867
			]
		},
		{
			"id": 28517,
			"enchant": 35297,
			"gems": [
				31867,
				24030
			]
		},
		{
			"id": 28753,
			"enchant": 22536
		},
		{
			"id": 29287,
			"enchant": 22536
		},
		{
			"id": 29370
		},
		{
			"id": 27683
		},
		{
			"id": 28770,
			"enchant": 22560
		},
		{
			"id": 29271
		},
		{
			"id": 27518
		}
	]}`),
};

export const P1_HORDE_PRESET = {
	name: 'P1 Horde Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<any>) => player.getFaction() == Faction.Horde,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 29093,
			"enchant": 29191,
			"gems": [
				24030,
				34220
			]
		},
		{
			"id": 28762
		},
		{
			"id": 29095,
			"enchant": 28886,
			"gems": [
				24056,
				24059
			]
		},
		{
			"id": 28766,
			"enchant": 33150
		},
		{
			"id": 21848,
			"enchant": 24003,
			"gems": [
				24059,
				24056
			]
		},
		{
			"id": 24250,
			"enchant": 22534,
			"gems": [
				31867
			]
		},
		{
			"id": 21847,
			"enchant": 28272,
			"gems": [
				31867,
				31867
			]
		},
		{
			"id": 21846,
			"gems": [
				31867,
				31867
			]
		},
		{
			"id": 24262,
			"enchant": 24274,
			"gems": [
				31867,
				31867,
				31867
			]
		},
		{
			"id": 28517,
			"enchant": 35297,
			"gems": [
				31867,
				24059
			]
		},
		{
			"id": 28753,
			"enchant": 22536
		},
		{
			"id": 28793,
			"enchant": 22536
		},
		{
			"id": 29370
		},
		{
			"id": 27683
		},
		{
			"id": 28770,
			"enchant": 22560
		},
		{
			"id": 29271
		},
		{
			"id": 27518
		}
	]}`),
};

export const P2_ALLIANCE_PRESET = {
	name: 'P2 Alliance Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<any>) => player.getFaction() == Faction.Alliance,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 30233,
			"enchant": 29191,
			"gems": [
				24059,
				34220
			]
		},
		{
			"id": 30015
		},
		{
			"id": 30235,
			"enchant": 28886,
			"gems": [
				24056,
				24059
			]
		},
		{
			"id": 28797,
			"enchant": 33150
		},
		{
			"id": 30231,
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
			"id": 30232,
			"enchant": 28272
		},
		{
			"id": 30038,
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
			"id": 30067,
			"enchant": 35297
		},
		{
			"id": 28753,
			"enchant": 22536
		},
		{
			"id": 29302,
			"enchant": 22536
		},
		{
			"id": 29370
		},
		{
			"id": 27683
		},
		{
			"id": 29988,
			"enchant": 22560
		},
		{
			"id": 32387
		}
	]}`),
};

export const P2_HORDE_PRESET = {
	name: 'P2 Horde Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<any>) => player.getFaction() == Faction.Horde,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 30233,
			"enchant": 29191,
			"gems": [
				31867,
				34220
			]
		},
		{
			"id": 30015
		},
		{
			"id": 30235,
			"enchant": 28886,
			"gems": [
				24056,
				31867
			]
		},
		{
			"id": 28797,
			"enchant": 33150
		},
		{
			"id": 30231,
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
			"id": 30232,
			"enchant": 28272
		},
		{
			"id": 30038,
			"gems": [
				24056,
				31867
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
			"id": 30067,
			"enchant": 35297
		},
		{
			"id": 28753,
			"enchant": 22536
		},
		{
			"id": 29302,
			"enchant": 22536
		},
		{
			"id": 29370
		},
		{
			"id": 27683
		},
		{
			"id": 29988,
			"enchant": 22560
		},
		{
			"id": 32387
		}
	]}`),
};

export const P3_PRESET = {
	name: 'P3 Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 31040,
			"enchant": 29191,
			"gems": [
				32218,
				34220
			]
		},
		{
			"id": 30015
		},
		{
			"id": 31049,
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
			"id": 31043,
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
			"id": 31035,
			"enchant": 28272,
			"gems": [
				32218
			]
		},
		{
			"id": 30914
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
				32218,
				32215
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
			"id": 32486
		},
		{
			"id": 32483
		},
		{
			"id": 32374,
			"enchant": 22560
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
			"id": 31040,
			"enchant": 29191,
			"gems": [
				32218,
				34220
			]
		},
		{
			"id": 33281
		},
		{
			"id": 31049,
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
			"id": 31043,
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
			"id": 31035,
			"enchant": 28272,
			"gems": [
				32218
			]
		},
		{
			"id": 30914
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
				32218,
				32215
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
			"id": 32483
		},
		{
			"id": 33829
		},
		{
			"id": 32374,
			"enchant": 22560
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
			"id": 34403,
			"enchant": 29191,
			"gems": [
				34220,
				32196
			]
		},
		{
			"id": 34204
		},
		{
			"id": 34391,
			"enchant": 28886,
			"gems": [
				32221,
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
			"id": 31043,
			"enchant": 24003,
			"gems": [
				32215,
				32215,
				32221
			]
		},
		{
			"id": 34446,
			"enchant": 22534,
			"gems": [
				35760
			]
		},
		{
			"id": 34407,
			"enchant": 28272,
			"gems": [
				32196,
				35760
			]
		},
		{
			"id": 34555,
			"gems": [
				32196
			]
		},
		{
			"id": 34169,
			"enchant": 24274,
			"gems": [
				32196,
				32196,
				35760
			]
		},
		{
			"id": 34572,
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
			"id": 34362,
			"enchant": 22536
		},
		{
			"id": 32483
		},
		{
			"id": 34429
		},
		{
			"id": 34336,
			"enchant": 22560
		},
		{
			"id": 34179
		},
		{
			"id": 32387
		}
	]}`),
};
