import { Conjured } from '/wotlk/core/proto/common.js';
import { Consumes } from '/wotlk/core/proto/common.js';
import { EquipmentSpec } from '/wotlk/core/proto/common.js';
import { Flask } from '/wotlk/core/proto/common.js';
import { Food } from '/wotlk/core/proto/common.js';
import { Glyphs } from '/wotlk/core/proto/common.js';
import { ItemSpec } from '/wotlk/core/proto/common.js';
import { Potions } from '/wotlk/core/proto/common.js';
import { Spec } from '/wotlk/core/proto/common.js';
import { WeaponImbue } from '/wotlk/core/proto/common.js';
import { Faction } from '/wotlk/core/proto/common.js';
import { SavedTalents } from '/wotlk/core/proto/ui.js';
import { Player } from '/wotlk/core/player.js';

import { Mage, Mage_Rotation as MageRotation, MageTalents as MageTalents, Mage_Options as MageOptions } from '/wotlk/core/proto/mage.js';
import { Mage_Rotation_Type as RotationType, Mage_Rotation_ArcaneRotation as ArcaneRotation, Mage_Rotation_FireRotation as FireRotation, Mage_Rotation_FrostRotation as FrostRotation } from '/wotlk/core/proto/mage.js';
import { Mage_Rotation_FireRotation_PrimarySpell as PrimaryFireSpell } from '/wotlk/core/proto/mage.js';
import { Mage_Rotation_ArcaneRotation_Filler as ArcaneFiller } from '/wotlk/core/proto/mage.js';
import { Mage_Options_ArmorType as ArmorType } from '/wotlk/core/proto/mage.js';

import * as Enchants from '/wotlk/core/constants/enchants.js';
import * as Gems from '/wotlk/core/proto_utils/gems.js';
import * as Tooltips from '/wotlk/core/constants/tooltips.js';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/wotlk/talent-calc and copy the numbers in the url.
export const ArcaneTalents = {
	name: 'Arcane',
	data: SavedTalents.create({
		talentsString: '2500250300030150330125--053500031003001',
	}),
};
export const FireTalents = {
	name: 'Fire',
	data: SavedTalents.create({
		talentsString: '2-505202012303331053125-043500001',
	}),
};
export const FrostTalents = {
	name: 'Frost',
	data: SavedTalents.create({
		talentsString: '2500250300030150330125--053500031003001',
	}),
};
export const DeepFrostTalents = {
	name: 'Deep Frost',
	data: SavedTalents.create({
		talentsString: '230015031003--0535000310230012241551',
	}),
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
	flask: Flask.FlaskOfTheFrostWyrm,
	food: Food.FoodFishFeast,
	defaultPotion: Potions.RunicManaPotion,
	defaultConjured: Conjured.ConjuredMageManaEmerald,
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
	defaultPotion: Potions.RunicManaPotion,
	defaultConjured: Conjured.ConjuredMageManaEmerald,
	flask: Flask.FlaskOfTheFrostWyrm,
	food: Food.FoodFishFeast,
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
	defaultPotion: Potions.RunicManaPotion,
	defaultConjured: Conjured.ConjuredMageManaEmerald,
	flask: Flask.FlaskOfTheFrostWyrm,
	food: Food.FoodFishFeast,
});

export const P5_ARCANE_PRESET = {
	name: 'TBC P5 Arcane Preset',
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
	name: 'TBC P5 Fire Preset',
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
	name: 'TBC P5 Frost Preset',
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
