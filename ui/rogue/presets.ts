import { BattleElixir } from '/wotlk/core/proto/common.js';
import { Conjured } from '/wotlk/core/proto/common.js';
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

import {
	Rogue,
	Rogue_Rotation as RogueRotation,
	Rogue_Rotation_Builder as Builder,
	Rogue_Options as RogueOptions,
} from '/wotlk/core/proto/rogue.js';

import * as Enchants from '/wotlk/core/constants/enchants.js';
import * as Gems from '/wotlk/core/proto_utils/gems.js';
import * as Tooltips from '/wotlk/core/constants/tooltips.js';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/wotlk/talent-calc and copy the numbers in the url.
export const CombatTalents = {
	name: 'Combat',
	data: SavedTalents.create({
		talentsString: '00532000532-0252051000035015223100501251',
	}),
};

export const DefaultRotation = RogueRotation.create({
	builder: Builder.Auto,
	maintainExposeArmor: false,
	useRupture: true,
	useShiv: true,
	minComboPointsForDamageFinisher: 3,
});

export const DefaultOptions = RogueOptions.create({
});

export const DefaultConsumes = Consumes.create({
	defaultPotion: Potions.PotionOfSpeed,
	defaultConjured: Conjured.ConjuredRogueThistleTea,
	battleElixir: BattleElixir.ElixirOfDeadlyStrikes,
	food: Food.FoodFishFeast,
	mainHandImbue: WeaponImbue.WeaponImbueRogueInstantPoison,
	offHandImbue: WeaponImbue.WeaponImbueRogueDeadlyPoison,
	scrollOfAgility: 5,
	scrollOfStrength: 5,
});

export const P1_PRESET = {
	name: 'P1 Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	gear: EquipmentSpec.fromJsonString(`{"items": [
        {
          "id": 37293,
          "enchant": 44879,
          "gems": [
            41398,
            40088
          ]
        },
        {
          "id": 37861
        },
        {
          "id": 37139,
          "enchant": 44871,
          "gems": [
            24061
          ]
        },
        {
          "id": 36947,
          "enchant": 55002
        },
        {
          "id": 44303,
          "enchant": 44623
        },
        {
          "id": 44203,
          "enchant": 60616,
          "gems": [
            0
          ]
        },
        {
          "id": 37409,
          "enchant": 60668,
          "gems": [
            0
          ]
        },
        {
          "id": 37194,
          "gems": [
            40014,
            0
          ]
        },
        {
          "id": 37644,
          "enchant": 38374
        },
        {
          "id": 44297,
          "enchant": 28279
        },
        {
          "id": 43251,
          "gems": [
            0
          ]
        },
        {
          "id": 37642
        },
        {
          "id": 37390
        },
        {
          "id": 37166
        },
        {
          "id": 37693,
          "enchant": 22559
        },
        {
          "id": 37856,
          "enchant": 22559
        },
        {
          "id": 37191
        }
      ]}`),
};
