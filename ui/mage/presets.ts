import { Conjured } from '/tbc/core/proto/common.js';
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

import { Mage, Mage_Rotation as MageRotation, MageTalents as MageTalents, Mage_Options as MageOptions } from '/tbc/core/proto/mage.js';
import { Mage_Rotation_Type as RotationType, Mage_Rotation_ArcaneRotation as ArcaneRotation, Mage_Rotation_FireRotation as FireRotation, Mage_Rotation_FrostRotation as FrostRotation } from '/tbc/core/proto/mage.js';
import { Mage_Rotation_FireRotation_PrimarySpell as PrimaryFireSpell } from '/tbc/core/proto/mage.js';
import { Mage_Rotation_ArcaneRotation_Filler as ArcaneFiller } from '/tbc/core/proto/mage.js';
import { Mage_Options_ArmorType as ArmorType } from '/tbc/core/proto/mage.js';

import * as Enchants from '/tbc/core/constants/enchants.js';
import * as Gems from '/tbc/core/proto_utils/gems.js';
import * as Tooltips from '/tbc/core/constants/tooltips.js';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

// Default talents. Uses the wowhead calculator format, make the talents on
// https://tbc.wowhead.com/talent-calc and copy the numbers in the url.
export const ArcaneTalents = {
	name: 'Arcane',
	data: '2500250300030150330125--053500031003001',
};
export const FireTalents = {
	name: 'Fire',
	data: '2-505202012303331053125-043500001',
};
export const FrostTalents = {
	name: 'Frost',
	data: '2500250300030150330125--053500031003001',
};
export const DeepFrostTalents = {
	name: 'Deep Frost',
	data: '230015031003--0535000310230012241551',
};

export const DefaultFireRotation = MageRotation.create({
	type: RotationType.Fire,
	fire: FireRotation.create({
		primarySpell: PrimaryFireSpell.Fireball,
		maintainImprovedScorch: true,
	}),
});

export const DefaultFireOptions = MageOptions.create({
	armor: ArmorType.MageArmor,
});

export const DefaultFireConsumes = Consumes.create({
	flask: Flask.FlaskOfPureDeath,
	food: Food.FoodBlackenedBasilisk,
	mainHandImbue: WeaponImbue.WeaponImbueBrilliantWizardOil,
	defaultPotion: Potions.SuperManaPotion,
	defaultConjured: Conjured.ConjuredFlameCap,
});

export const DefaultFrostRotation = MageRotation.create({
	type: RotationType.Frost,
	frost: FrostRotation.create({
		waterElementalDisobeyChance: 0.1,
	}),
});

export const DefaultFrostOptions = MageOptions.create({
	armor: ArmorType.MageArmor,
});

export const DefaultFrostConsumes = Consumes.create({
	defaultPotion: Potions.SuperManaPotion,
	defaultConjured: Conjured.ConjuredMageManaEmerald,
	flask: Flask.FlaskOfPureDeath,
	food: Food.FoodBlackenedBasilisk,
	mainHandImbue: WeaponImbue.WeaponImbueBrilliantWizardOil,
});

export const DefaultArcaneRotation = MageRotation.create({
	type: RotationType.Arcane,
	arcane: ArcaneRotation.create({
		filler: ArcaneFiller.Frostbolt,
		arcaneBlastsBetweenFillers: 3,
		startRegenRotationPercent: 0.2,
		stopRegenRotationPercent: 0.5,
	}),
});

export const DefaultArcaneOptions = MageOptions.create({
	armor: ArmorType.MageArmor,
});

export const DefaultArcaneConsumes = Consumes.create({
	defaultPotion: Potions.SuperManaPotion,
	defaultConjured: Conjured.ConjuredMageManaEmerald,
	flask: Flask.FlaskOfBlindingLight,
	food: Food.FoodBlackenedBasilisk,
	mainHandImbue: WeaponImbue.WeaponImbueBrilliantWizardOil,
});

export const P1_ARCANE_PRESET = {
	name: 'P1 Arcane Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecMage>) => player.getRotation().type == RotationType.Arcane,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 29076,
			"enchant": 29191,
			"gems": [
				34220,
				24056
			]
		},
		{
			"id": 28762
		},
		{
			"id": 29079,
			"enchant": 28886,
			"gems": [
				24047,
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
				31867
			]
		},
		{
			"id": 28411,
			"enchant": 22534,
			"gems": [
				31867
			]
		},
		{
			"id": 21847,
			"enchant": 28272,
			"gems": [
				24047,
				31867
			]
		},
		{
			"id": 21846,
			"gems": [
				24047,
				24056
			]
		},
		{
			"id": 29078,
			"enchant": 24274
		},
		{
			"id": 28517,
			"enchant": 35297,
			"gems": [
				24030,
				24047
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
			"id": 28785
		},
		{
			"id": 28770,
			"enchant": 22560
		},
		{
			"id": 29271
		},
		{
			"id": 28783
		}
	]}`),
};

export const P1_FIRE_PRESET = {
	name: 'P1 Fire Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecMage>) => player.getRotation().type == RotationType.Fire,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 29076,
			"enchant": 29191,
			"gems": [
				34220,
				24056
			]
		},
		{
			"id": 28134
		},
		{
			"id": 29079,
			"enchant": 28886,
			"gems": [
				31867,
				24030
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
				31867
			]
		},
		{
			"id": 28411,
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
				24056
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
				31867
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
			"id": 29370
		},
		{
			"id": 27683
		},
		{
			"id": 28802,
			"enchant": 22560
		},
		{
			"id": 29270
		},
		{
			"id": 28673
		}
	]}`),
};

export const P1_FROST_PRESET = {
	name: 'P1 Frost Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecMage>) => player.getRotation().type == RotationType.Frost,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 29076,
			"enchant": 29191,
			"gems": [
				34220,
				24056
			]
		},
		{
			"id": 28762
		},
		{
			"id": 29079,
			"enchant": 28886,
			"gems": [
				31867,
				31867
			]
		},
		{
			"id": 28766,
			"enchant": 33150
		},
		{
			"id": 21871,
			"enchant": 24003,
			"gems": [
				31867,
				31867
			]
		},
		{
			"id": 28411,
			"enchant": 22534,
			"gems": [
				31867
			]
		},
		{
			"id": 28780,
			"enchant": 28272,
			"gems": [
				31867,
				24056
			]
		},
		{
			"id": 24256,
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
			"id": 21870,
			"enchant": 35297,
			"gems": [
				31867,
				31867
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
			"id": 29370
		},
		{
			"id": 27683
		},
		{
			"id": 28802,
			"enchant": 22561
		},
		{
			"id": 29269
		},
		{
			"id": 28673
		}
	]}`),
};

export const P2_ARCANE_PRESET = {
	name: 'P2 Arcane Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecMage>) => player.getRotation().type == RotationType.Arcane,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 30206,
			"enchant": 29191,
			"gems": [
				34220,
				24047
			]
		},
		{
			"id": 30015
		},
		{
			"id": 30210,
			"enchant": 28886,
			"gems": [
				24047,
				24056
			]
		},
		{
			"id": 29992,
			"enchant": 33150
		},
		{
			"id": 30196,
			"enchant": 24003,
			"gems": [
				24047,
				24047,
				24056
			]
		},
		{
			"id": 29918,
			"enchant": 22534
		},
		{
			"id": 29987,
			"enchant": 28272
		},
		{
			"id": 30038,
			"gems": [
				24047,
				24056
			]
		},
		{
			"id": 30207,
			"enchant": 24274,
			"gems": [
				24047
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
			"id": 29287,
			"enchant": 22536
		},
		{
			"id": 29370
		},
		{
			"id": 30720
		},
		{
			"id": 29988,
			"enchant": 22560
		},
		{
			"id": 28783
		}
	]}`),
};

export const P2_FIRE_PRESET = {
	name: 'P2 Fire Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecMage>) => player.getRotation().type == RotationType.Fire,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 32494,
			"enchant": 29191,
			"gems": [
				34220,
				24056
			]
		},
		{
			"id": 30015
		},
		{
			"id": 30024,
			"enchant": 28886
		},
		{
			"id": 28766,
			"enchant": 33150
		},
		{
			"id": 30107,
			"enchant": 24003,
			"gems": [
				31867,
				31867,
				24056
			]
		},
		{
			"id": 29918,
			"enchant": 22534
		},
		{
			"id": 21847,
			"enchant": 28272,
			"gems": [
				31867,
				24030
			]
		},
		{
			"id": 30038,
			"gems": [
				31867,
				24056
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
			"id": 30037,
			"enchant": 35297
		},
		{
			"id": 28753,
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
			"id": 27683
		},
		{
			"id": 30095,
			"enchant": 22560
		},
		{
			"id": 29270
		},
		{
			"id": 29982
		}
	]}`),
};

export const P2_FROST_PRESET = {
	name: 'P2 Frost Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecMage>) => player.getRotation().type == RotationType.Frost,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 30206,
			"enchant": 29191,
			"gems": [
				34220,
				24056
			]
		},
		{
			"id": 30015
		},
		{
			"id": 30210,
			"enchant": 28886,
			"gems": [
				31867,
				24056
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
				31867,
				31867,
				24056
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
				31867,
				24056
			]
		},
		{
			"id": 30038,
			"gems": [
				31867,
				24056
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
			"id": 21870,
			"enchant": 35297,
			"gems": [
				24030,
				31867
			]
		},
		{
			"id": 28753,
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
			"id": 27683
		},
		{
			"id": 30095,
			"enchant": 22561
		},
		{
			"id": 29269
		},
		{
			"id": 29982
		}
	]}`),
};

export const P3_ARCANE_PRESET = {
	name: 'P3 Arcane Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecMage>) => player.getRotation().type == RotationType.Arcane,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 30206,
			"enchant": 29191,
			"gems": [
				34220,
				32204
			]
		},
		{
			"id": 30015
		},
		{
			"id": 30210,
			"enchant": 28886,
			"gems": [
				32204,
				32215
			]
		},
		{
			"id": 32331,
			"enchant": 33150
		},
		{
			"id": 30196,
			"enchant": 24003,
			"gems": [
				32204,
				32204,
				32215
			]
		},
		{
			"id": 30870,
			"enchant": 22534,
			"gems": [
				32204
			]
		},
		{
			"id": 30205,
			"enchant": 28272
		},
		{
			"id": 30888,
			"gems": [
				32204,
				32204
			]
		},
		{
			"id": 31058,
			"enchant": 24274,
			"gems": [
				32204
			]
		},
		{
			"id": 32239,
			"enchant": 35297,
			"gems": [
				32204,
				32204
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
			"id": 30720
		},
		{
			"id": 32374,
			"enchant": 22560
		},
		{
			"id": 28783
		}
	]}`),
};

export const P3_FIRE_PRESET = {
	name: 'P3 Fire Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecMage>) => player.getRotation().type == RotationType.Fire,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 31056,
			"enchant": 29191,
			"gems": [
				34220,
				32196
			]
		},
		{
			"id": 32589
		},
		{
			"id": 31059,
			"enchant": 28886,
			"gems": [
				32221,
				32215
			]
		},
		{
			"id": 32331,
			"enchant": 33150
		},
		{
			"id": 31057,
			"enchant": 24003,
			"gems": [
				32221,
				32221,
				32215
			]
		},
		{
			"id": 32586,
			"enchant": 22534
		},
		{
			"id": 31055,
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
			"id": 29305,
			"enchant": 22536
		},
		{
			"id": 32483
		},
		{
			"id": 27683
		},
		{
			"id": 30910,
			"enchant": 22560
		},
		{
			"id": 30872
		},
		{
			"id": 29982
		}
	]}`),
};

export const P3_FROST_PRESET = {
	name: 'P3 Frost Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecMage>) => player.getRotation().type == RotationType.Frost,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 31056,
			"enchant": 29191,
			"gems": [
				34220,
				32196
			]
		},
		{
			"id": 32349
		},
		{
			"id": 31059,
			"enchant": 28886,
			"gems": [
				32221,
				32215
			]
		},
		{
			"id": 32331,
			"enchant": 33150
		},
		{
			"id": 31057,
			"enchant": 24003,
			"gems": [
				32218,
				32218,
				32215
			]
		},
		{
			"id": 32586,
			"enchant": 22534
		},
		{
			"id": 31055,
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
			"id": 29305,
			"enchant": 22536
		},
		{
			"id": 32483
		},
		{
			"id": 27683
		},
		{
			"id": 30910,
			"enchant": 22560
		},
		{
			"id": 30872
		},
		{
			"id": 29982
		}
	]}`),
};

export const P4_ARCANE_PRESET = {
	name: 'P4 Arcane Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecMage>) => player.getRotation().type == RotationType.Arcane,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 30206,
			"enchant": 29191,
			"gems": [
				34220,
				32204
			]
		},
		{
			"id": 33281
		},
		{
			"id": 30210,
			"enchant": 28886,
			"gems": [
				32204,
				32215
			]
		},
		{
			"id": 32331,
			"enchant": 33150
		},
		{
			"id": 30196,
			"enchant": 24003,
			"gems": [
				32204,
				32204,
				32215
			]
		},
		{
			"id": 30870,
			"enchant": 22534,
			"gems": [
				32204
			]
		},
		{
			"id": 30205,
			"enchant": 28272
		},
		{
			"id": 30888,
			"gems": [
				32204,
				32204
			]
		},
		{
			"id": 31058,
			"enchant": 24274,
			"gems": [
				32204
			]
		},
		{
			"id": 32239,
			"enchant": 35297,
			"gems": [
				32204,
				32204
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
			"id": 30720
		},
		{
			"id": 32374,
			"enchant": 22560
		},
		{
			"id": 33192,
			"gems": [
				32204
			]
		}
	]}`),
};

export const P4_FIRE_PRESET = {
	name: 'P4 Fire Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecMage>) => player.getRotation().type == RotationType.Fire,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 31056,
			"enchant": 29191,
			"gems": [
				34220,
				32196
			]
		},
		{
			"id": 32589
		},
		{
			"id": 31059,
			"enchant": 28886,
			"gems": [
				32221,
				32215
			]
		},
		{
			"id": 32331,
			"enchant": 33150
		},
		{
			"id": 31057,
			"enchant": 24003,
			"gems": [
				32221,
				32221,
				32215
			]
		},
		{
			"id": 32586,
			"enchant": 22534
		},
		{
			"id": 31055,
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
			"id": 30910,
			"enchant": 22560
		},
		{
			"id": 30872
		},
		{
			"id": 29982
		}
	]}`),
};

export const P4_FROST_PRESET = {
	name: 'P4 Frost Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecMage>) => player.getRotation().type == RotationType.Frost,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 31056,
			"enchant": 29191,
			"gems": [
				34220,
				32196
			]
		},
		{
			"id": 33281
		},
		{
			"id": 31059,
			"enchant": 28886,
			"gems": [
				32221,
				32215
			]
		},
		{
			"id": 32524,
			"enchant": 33150
		},
		{
			"id": 31057,
			"enchant": 24003,
			"gems": [
				32221,
				32221,
				32215
			]
		},
		{
			"id": 32586,
			"enchant": 22534
		},
		{
			"id": 31055,
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
			"id": 30910,
			"enchant": 22560
		},
		{
			"id": 30872
		},
		{
			"id": 29982
		}
	]}`),
};

export const P5_ARCANE_PRESET = {
	name: 'P5 Arcane Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecMage>) => player.getRotation().type == RotationType.Arcane,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 34405,
			"enchant": 29191,
			"gems": [
				34220,
				32196
			]
		},
		{
			"id": 33281
		},
		{
			"id": 30210,
			"enchant": 28886,
			"gems": [
				32204,
				32215
			]
		},
		{
			"id": 34242,
			"enchant": 33150,
			"gems": [
				32204
			]
		},
		{
			"id": 34399,
			"enchant": 24003,
			"gems": [
				32204,
				32215,
				32204
			]
		},
		{
			"id": 34447,
			"enchant": 22534,
			"gems": [
				32204
			]
		},
		{
			"id": 30205,
			"enchant": 28272
		},
		{
			"id": 34557,
			"gems": [
				32204
			]
		},
		{
			"id": 34386,
			"enchant": 24274,
			"gems": [
				32196,
				32204,
				32204
			]
		},
		{
			"id": 34574,
			"enchant": 35297,
			"gems": [
				32204
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
			"id": 30720
		},
		{
			"id": 34336,
			"enchant": 22560
		},
		{
			"id": 34179
		},
		{
			"id": 34347,
			"gems": [
				32204
			]
		}
	]}`),
};

export const P5_FIRE_PRESET = {
	name: 'P5 Fire Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecMage>) => player.getRotation().type == RotationType.Fire,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 34340,
			"enchant": 29191,
			"gems": [
				34220,
				32215
			]
		},
		{
			"id": 34204
		},
		{
			"id": 31059,
			"enchant": 28886,
			"gems": [
				32196,
				35760
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
			"enchant": 24003,
			"gems": [
				32196,
				35760,
				35760
			]
		},
		{
			"id": 34447,
			"enchant": 22534,
			"gems": [
				32215
			]
		},
		{
			"id": 34344,
			"enchant": 28272,
			"gems": [
				35760,
				32196
			]
		},
		{
			"id": 34557,
			"gems": [
				32221
			]
		},
		{
			"id": 34181,
			"enchant": 24274,
			"gems": [
				32196,
				32196,
				32221
			]
		},
		{
			"id": 34574,
			"enchant": 35297,
			"gems": [
				32221
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
			"enchant": 22560
		},
		{
			"id": 34179
		},
		{
			"id": 34347,
			"gems": [
				35760
			]
		}
	]}`),
};

export const P5_FROST_PRESET = {
	name: 'P5 Frost Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecMage>) => player.getRotation().type == RotationType.Frost,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 34340,
			"enchant": 29191,
			"gems": [
				34220,
				32215
			]
		},
		{
			"id": 34204
		},
		{
			"id": 31059,
			"enchant": 28886,
			"gems": [
				32196,
				35760
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
			"enchant": 24003,
			"gems": [
				32196,
				35760,
				35760
			]
		},
		{
			"id": 34447,
			"enchant": 22534,
			"gems": [
				32215
			]
		},
		{
			"id": 34344,
			"enchant": 28272,
			"gems": [
				35760,
				32196
			]
		},
		{
			"id": 34557,
			"gems": [
				32221
			]
		},
		{
			"id": 34181,
			"enchant": 24274,
			"gems": [
				32196,
				32196,
				32221
			]
		},
		{
			"id": 34574,
			"enchant": 35297,
			"gems": [
				32221
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
			"enchant": 22560
		},
		{
			"id": 34179
		},
		{
			"id": 34347,
			"gems": [
				35760
			]
		}
	]}`),
};
