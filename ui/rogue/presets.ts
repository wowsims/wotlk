import { BattleElixir, Flask } from '../core/proto/common.js';
import { Conjured } from '../core/proto/common.js';
import { Consumes } from '../core/proto/common.js';

import { EquipmentSpec } from '../core/proto/common.js';
import { Food } from '../core/proto/common.js';
import { Glyphs } from '../core/proto/common.js';
import { Potions } from '../core/proto/common.js';
import { SavedTalents } from '../core/proto/ui.js';

import {
	Rogue_Rotation as RogueRotation,
	Rogue_Rotation_Builder as Builder,
	Rogue_Rotation_Filler as Filler,
	Rogue_Options as RogueOptions,
  Rogue_Options_PoisonImbue as Poison,
  RogueMajorGlyph,
} from '../core/proto/rogue.js';

import * as Tooltips from '../core/constants/tooltips.js';

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/wotlk/talent-calc and copy the numbers in the url.
export const CombatTalents = {
	name: 'Combat',
	data: SavedTalents.create({
		talentsString: '00532000523-0252051050035010223100501251',
    glyphs: Glyphs.create({
      major1: RogueMajorGlyph.GlyphOfKillingSpree,
      major2: RogueMajorGlyph.GlyphOfTricksOfTheTrade,
      major3: RogueMajorGlyph.GlyphOfRupture,
    })
	}),
};

export const AssassinationTalents = {
	name: 'Assassination',
	data: SavedTalents.create({
		talentsString: '005303005352100520103331051-005005003-502',
    glyphs: Glyphs.create({
      major1: RogueMajorGlyph.GlyphOfMutilate,
      major2: RogueMajorGlyph.GlyphOfTricksOfTheTrade,
      major3: RogueMajorGlyph.GlyphOfHungerForBlood,
    })
	}),
};

export const DefaultRotation = RogueRotation.create({
	builder: Builder.Auto,
  filler: Filler.NoFiller,
	maintainExposeArmor: false,
  maintainTricksOfTheTrade: true,
	useRupture: false,
	useShiv: false,
	useEnvenom: false,
});

export const DefaultOptions = RogueOptions.create({
  mhImbue: Poison.DeadlyPoison,
  ohImbue: Poison.InstantPoison,
});

export const DefaultConsumes = Consumes.create({
	defaultPotion: Potions.PotionOfSpeed,
	defaultConjured: Conjured.ConjuredRogueThistleTea,
  flask: Flask.FlaskOfEndlessRage,
	food: Food.FoodMegaMammothMeal,
});

export const PRERAID_PRESET = {
	name: 'Pre-Raid',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	gear: EquipmentSpec.fromJsonString(`{"items": [
    {
      "id": 42550,
      "enchant": 44879,
      "gems": [
        41398,
        40058
      ]
    },
    {
      "id": 40678
    },
    {
      "id": 43481,
      "enchant": 44871
    },
    {
      "id": 38614,
      "enchant": 55002
    },
    {
      "id": 39558,
      "enchant": 44489,
      "gems": [
        40003,
        42702
      ]
    },
    {
      "id": 34448,
      "enchant": 44484,
      "gems": [
        40003,
        0
      ]
    },
    {
      "id": 39560,
      "enchant": 54999,
      "gems": [
        40058,
        0
      ]
    },
    {
      "id": 40694,
      "gems": [
        40003,
        40003
      ]
    },
    {
      "id": 37644,
      "enchant": 38374
    },
    {
      "id": 34575,
      "enchant": 55016,
      "gems": [
        40003
      ]
    },
    {
      "id": 40586
    },
    {
      "id": 37642
    },
    {
      "id": 40684
    },
    {
      "id": 44253
    },
    {
      "id": 37856,
      "enchant": 44492
    },
    {
      "id": 37667,
      "enchant": 44492
    },
    {
      "id": 43612
    }
  ]}`),
};

export const P1_PRESET = {
	name: 'P1',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	gear: EquipmentSpec.fromJsonString(`{"items": [
    {
      "id": 42550,
      "enchant": 44879,
      "gems": [
        41398,
        40058
      ]
    },
    {
      "id": 40678
    },
    {
      "id": 43481,
      "enchant": 44871
    },
    {
      "id": 38614,
      "enchant": 55002
    },
    {
      "id": 39558,
      "enchant": 44489,
      "gems": [
        40003,
        42702
      ]
    },
    {
      "id": 34448,
      "enchant": 44484,
      "gems": [
        40003,
        0
      ]
    },
    {
      "id": 39560,
      "enchant": 54999,
      "gems": [
        40058,
        0
      ]
    },
    {
      "id": 40694,
      "gems": [
        40003,
        40003
      ]
    },
    {
      "id": 37644,
      "enchant": 38374
    },
    {
      "id": 34575,
      "enchant": 55016,
      "gems": [
        40003
      ]
    },
    {
      "id": 40586
    },
    {
      "id": 37642
    },
    {
      "id": 40684
    },
    {
      "id": 44253
    },
    {
      "id": 37856,
      "enchant": 44492
    },
    {
      "id": 37667,
      "enchant": 44492
    },
    {
      "id": 43612
    }
  ]}`),
};