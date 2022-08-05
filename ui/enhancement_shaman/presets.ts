import { Consumes } from '../core/proto/common.js';

import { EquipmentSpec } from '../core/proto/common.js';
import { Flask } from '../core/proto/common.js';
import { Food } from '../core/proto/common.js';
import { Potions } from '../core/proto/common.js';
import { SavedTalents } from '../core/proto/ui.js';

import { EnhancementShaman_Rotation as EnhancementShamanRotation, EnhancementShaman_Options as EnhancementShamanOptions, ShamanShield } from '../core/proto/shaman.js';
import {
	AirTotem,
	EarthTotem,
	FireTotem,
	WaterTotem,
	ShamanTotems,
    ShamanImbue,
    ShamanSyncType
} from '../core/proto/shaman.js';

import * as Tooltips from '../core/constants/tooltips.js';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/wotlk/talent-calc and copy the numbers in the url.
export const StandardTalents = {
	name: 'Standard',
	data: SavedTalents.create({
		talentsString: '053030152-30205023105021333031131031051',
	}),
};

export const DefaultRotation = EnhancementShamanRotation.create({
	totems: ShamanTotems.create({
		earth: EarthTotem.StrengthOfEarthTotem,
		air: AirTotem.WindfuryTotem,
		fire: FireTotem.SearingTotem,
		water: WaterTotem.ManaSpringTotem,
	}),
});

export const DefaultOptions = EnhancementShamanOptions.create({
	shield: ShamanShield.LightningShield,
	bloodlust: true,
	imbueMH: ShamanImbue.WindfuryWeapon,
	imbueOH: ShamanImbue.FlametongueWeapon,
	syncType: ShamanSyncType.SyncMainhandOffhandSwings,
});

export const DefaultConsumes = Consumes.create({
	defaultPotion: Potions.HastePotion,
	flask: Flask.FlaskOfRelentlessAssault,
	food: Food.FoodRoastedClefthoof,
});

export const P1_PRESET = {
	name: 'P1 Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 29040,
			"enchant": 29192,
			"gems": [
				32409,
				24027
			]
		},
		{
			"id": 29381
		},
		{
			"id": 29043,
			"enchant": 28888,
			"gems": [
				24027,
				24058
			]
		},
		{
			"id": 24259,
			"enchant": 34004,
			"gems": [
				24027
			]
		},
		{
			"id": 29038,
			"enchant": 24003,
			"gems": [
				24027,
				24027,
				24058
			]
		},
		{
			"id": 25697,
			"enchant": 27899,
			"gems": [
				24027
			]
		},
		{
			"id": 29039,
			"enchant": 33995
		},
		{
			"id": 28656
		},
		{
			"id": 30534,
			"enchant": 29535,
			"gems": [
				24054,
				24054,
				24058
			]
		},
		{
			"id": 28746,
			"enchant": 28279,
			"gems": [
				24027,
				24027
			]
		},
		{
			"id": 28757
		},
		{
			"id": 29283
		},
		{
			"id": 28830
		},
		{
			"id": 29383
		},
		{
			"id": 28767,
			"enchant": 22559
		},
		{
			"id": 27872,
			"enchant": 22559
		},
		{
			"id": 27815
		}
	]}`),
};
