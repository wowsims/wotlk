import { Consumes } from '/wotlk/core/proto/common.js';
import { Flask } from '/wotlk/core/proto/common.js';
import { Food } from '/wotlk/core/proto/common.js';
import { Glyphs } from '/wotlk/core/proto/common.js';
import { EquipmentSpec } from '/wotlk/core/proto/common.js';
import { ItemSpec } from '/wotlk/core/proto/common.js';
import { Potions } from '/wotlk/core/proto/common.js';
import { WeaponImbue } from '/wotlk/core/proto/common.js';
import { Faction } from '/wotlk/core/proto/common.js';
import { SavedTalents } from '/wotlk/core/proto/ui.js';
import { Player } from '/wotlk/core/player.js';

import { BalanceDruid, BalanceDruid_Rotation as BalanceDruidRotation, DruidTalents as DruidTalents, BalanceDruid_Options as BalanceDruidOptions, BalanceDruid_Rotation_RotationType as RotationType } from '/wotlk/core/proto/druid.js';

import * as Enchants from '/wotlk/core/constants/enchants.js';
import * as Gems from '/wotlk/core/proto_utils/gems.js';
import * as Tooltips from '/wotlk/core/constants/tooltips.js';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/wotlk/talent-calc and copy the numbers in the url.
export const StandardTalents = {
	name: 'Standard',
	data: SavedTalents.create({
		talentsString: '510022312503135231351--520033',
	}),
};

export const DefaultRotation = BalanceDruidRotation.create({
	type: RotationType.Adaptive,
});

export const DefaultOptions = BalanceDruidOptions.create({
});

export const DefaultConsumes = Consumes.create({
	flask: Flask.FlaskOfTheFrostWyrm,
	food: Food.FoodFishFeast,
	defaultPotion: Potions.PotionOfSpeed,
});

export const P5_PRESET = {
	name: 'TBC P5 Preset',
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
