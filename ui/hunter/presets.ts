import { Consumes } from '/wotlk/core/proto/common.js';
import { EquipmentSpec } from '/wotlk/core/proto/common.js';
import { Flask } from '/wotlk/core/proto/common.js';
import { Food } from '/wotlk/core/proto/common.js';
import { Potions } from '/wotlk/core/proto/common.js';
import { WeaponImbue } from '/wotlk/core/proto/common.js';
import { Player } from '/wotlk/core/player.js';

import {
	Hunter_Rotation as HunterRotation,
	Hunter_Rotation_WeaveType as WeaveType,
	Hunter_Options as HunterOptions,
	Hunter_Options_Ammo as Ammo,
	Hunter_Options_PetType as PetType,
} from '/wotlk/core/proto/hunter.js';

import * as Tooltips from '/wotlk/core/constants/tooltips.js';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/wotlk/talent-calc and copy the numbers in the url.
export const BeastMasteryTalents = {
	name: 'BM',
	// TODO
	data: '5-035335001230013233105000351-4',
};

export const MarksmanTalents = {
	name: 'Marksman',
	data: '5-035335001230013233105000351-4',
};

export const SurvivalTalents = {
	name: 'Survival',
	data: '-03502-5300032100233030523035001331',
};

export const DefaultRotation = HunterRotation.create({
	useMultiShot: true,
	useArcaneShot: true,
	viperStartManaPercent: 0.1,
	viperStopManaPercent: 0.3,

	weave: WeaveType.WeaveNone,
	timeToWeaveMs: 500,
	percentWeaved: 0.8,
});

export const DefaultOptions = HunterOptions.create({
	ammo: Ammo.TimelessArrow,
	petType: PetType.Ravager,
	petUptime: 1,
	latencyMs: 30,
});

export const DefaultConsumes = Consumes.create({
	defaultPotion: Potions.HastePotion,
	flask: Flask.FlaskOfRelentlessAssault,
	food: Food.FoodGrilledMudfish,
	mainHandImbue: WeaponImbue.WeaponImbueAdamantiteSharpeningStone,
	offHandImbue: WeaponImbue.WeaponImbueAdamantiteSharpeningStone,
});

export const P1_PRESET = {
	name: 'P1 Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	//enableWhen: (player: Player<any>) => player.getTalentTree() == 2,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 34333,
			"enchant": 29192,
			"gems": [
				32194,
				32409
			]
		},
		{
			"id": 34177
		},
		{
			"id": 31006,
			"enchant": 23548,
			"gems": [
				32194,
				32194
			]
		},
		{
			"id": 34241,
			"enchant": 34004,
			"gems": [
				32194
			]
		},
		{
			"id": 34397,
			"enchant": 24003,
			"gems": [
				32194,
				32194,
				32194
			]
		},
		{
			"id": 34443,
			"enchant": 34002,
			"gems": [
				32194
			]
		},
		{
			"id": 34343,
			"enchant": 19445,
			"gems": [
				32194,
				32226
			]
		},
		{
			"id": 34549,
			"gems": [
				32194
			]
		},
		{
			"id": 34188,
			"enchant": 29535,
			"gems": [
				32194,
				32194,
				32194
			]
		},
		{
			"id": 34570,
			"enchant": 22544,
			"gems": [
				32226
			]
		},
		{
			"id": 34887
		},
		{
			"id": 34361
		},
		{
			"id": 34427
		},
		{
			"id": 28830
		},
		{
			"id": 34183,
			"enchant": 22556,
			"gems": [
				32194
			]
		},
		{
			"id": 34334,
			"enchant": 23766
		}
	]}`),
};
