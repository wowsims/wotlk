import { Consumes } from '/wotlk/core/proto/common.js';
import { EquipmentSpec } from '/wotlk/core/proto/common.js';
import { Flask } from '/wotlk/core/proto/common.js';
import { Food } from '/wotlk/core/proto/common.js';
import { Glyphs } from '/wotlk/core/proto/common.js';
import { ItemSpec } from '/wotlk/core/proto/common.js';
import { Potions } from '/wotlk/core/proto/common.js';
import { WeaponImbue } from '/wotlk/core/proto/common.js';
import { Faction } from '/wotlk/core/proto/common.js';
import { SavedTalents } from '/wotlk/core/proto/ui.js';
import { Player } from '/wotlk/core/player.js';

import { ShadowPriest, ShadowPriest_Rotation as Rotation, ShadowPriest_Options as Options, ShadowPriest_Rotation_RotationType } from '/wotlk/core/proto/priest.js';

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
		talentsString: '05032031--325023051223010323151301351',
	}),
};

export const DefaultRotation = Rotation.create({
	rotationType: ShadowPriest_Rotation_RotationType.Ideal,
});

export const DefaultOptions = Options.create({
	useShadowfiend: true,
});

export const DefaultConsumes = Consumes.create({
	flask: Flask.FlaskOfTheFrostWyrm,
	food: Food.FoodFishFeast,
	defaultPotion: Potions.PotionOfSpeed,
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
