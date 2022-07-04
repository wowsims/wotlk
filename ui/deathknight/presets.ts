import { Consumes } from '/wotlk/core/proto/common.js';
import { EquipmentSpec } from '/wotlk/core/proto/common.js';
import { Flask } from '/wotlk/core/proto/common.js';
import { Food } from '/wotlk/core/proto/common.js';
import { ItemSpec } from '/wotlk/core/proto/common.js';
import { Potions } from '/wotlk/core/proto/common.js';
import { Spec } from '/wotlk/core/proto/common.js';
import { WeaponImbue } from '/wotlk/core/proto/common.js';
import { Faction } from '/wotlk/core/proto/common.js';
import { Player } from '/wotlk/core/player.js';

import {
	DeathKnightTalents as DeathKnightTalents,
	DeathKnight,
	DeathKnight_Rotation as DeathKnightRotation,
	DeathKnight_Options as DeathKnightOptions,
} from '/wotlk/core/proto/deathknight.js';

import * as Enchants from '/wotlk/core/constants/enchants.js';
import * as Gems from '/wotlk/core/proto_utils/gems.js';
import * as Tooltips from '/wotlk/core/constants/tooltips.js';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wotlk.wowhead.com/talent-calc and copy the numbers in the url.
export const BloodTalents = {
	name: 'Blood Tank',
	data: '',
};

export const DefaultRotation = DeathKnightRotation.create({
	useScourgeStrike: false,
});

export const DefaultOptions = DeathKnightOptions.create({
	dualWhield: true,
});

export const DefaultConsumes = Consumes.create({
	flask: Flask.FlaskOfRelentlessAssault,
	food: Food.FoodRoastedClefthoof,
	defaultPotion: Potions.HastePotion,
	mainHandImbue: WeaponImbue.WeaponImbueAdamantiteSharpeningStone,
	offHandImbue: WeaponImbue.WeaponImbueAdamantiteSharpeningStone,
});

export const P1_FURY_PRESET = {
	name: 'P1 Fury Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecDeathKnight>) => player.getTalents().scourgeStrike,
	gear: EquipmentSpec.fromJsonString(`{"items": [
	  {
			"id": 29021,
			"enchant": 29192,
			"gems": [
				32409,
				24048
			]
		},
		{
			"id": 29381
		},
		{
			"id": 29023,
			"enchant": 28888,
			"gems": [
				24048,
				24067
			]
		},
		{
			"id": 24259,
			"enchant": 34004,
			"gems": [
				24058
			]
		},
		{
			"id": 29019,
			"enchant": 24003,
			"gems": [
				24048,
				24048,
				24048
			]
		},
		{
			"id": 28795,
			"enchant": 27899,
			"gems": [
				24067,
				24058
			]
		},
		{
			"id": 28824,
			"enchant": 33995,
			"gems": [
				24067,
				24048
			]
		},
		{
			"id": 28779,
			"gems": [
				24058,
				24067
			]
		},
		{
			"id": 28741,
			"enchant": 29535,
			"gems": [
				24048,
				24048,
				24048
			]
		},
		{
			"id": 28608,
			"enchant": 28279,
			"gems": [
				24058,
				24048
			]
		},
		{
			"id": 28757
		},
		{
			"id": 30834
		},
		{
			"id": 29383
		},
		{
			"id": 28830
		},
		{
			"id": 28438,
			"enchant": 22559
		},
		{
			"id": 28729,
			"enchant": 22559
		},
		{
			"id": 30279
		}
	]}`),
};
