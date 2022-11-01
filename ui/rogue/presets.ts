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
	Rogue_Options as RogueOptions,
	Rogue_Options_PoisonImbue as Poison,
	RogueMajorGlyph,
	Rogue_Rotation_Frequency,
	Rogue_Rotation_AssassinationPriority,
	Rogue_Rotation_CombatPriority,
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
	exposeArmorFrequency: Rogue_Rotation_Frequency.Never,
	minimumComboPointsExposeArmor: 4,
	tricksOfTheTradeFrequency: Rogue_Rotation_Frequency.Maintain,
	assassinationFinisherPriority: Rogue_Rotation_AssassinationPriority.EnvenomRupture,
	combatFinisherPriority: Rogue_Rotation_CombatPriority.RuptureEviscerate,
	minimumComboPointsPrimaryFinisher: 3,
	minimumComboPointsSecondaryFinisher: 2,
  envenomEnergyThreshold: 60,
});

export const DefaultOptions = RogueOptions.create({
	mhImbue: Poison.DeadlyPoison,
	ohImbue: Poison.InstantPoison,
  applyPoisonsManually: false,
});

export const DefaultConsumes = Consumes.create({
	defaultPotion: Potions.PotionOfSpeed,
	defaultConjured: Conjured.ConjuredRogueThistleTea,
	flask: Flask.FlaskOfEndlessRage,
	food: Food.FoodMegaMammothMeal,
});

export const PRERAID_PRESET_ASSASSINATION = {
	name: 'Pre-Raid Assassination',
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

export const PRERAID_PRESET_COMBAT = {
  name: 'Pre-Raid Combat',
  tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
  gear: EquipmentSpec.fromJsonString(`{"items": [
    {
      "id": 42550,
      "enchant": 44879,
      "gems": [
        41398,
        40014
      ]
    },
    {
      "id": 40678
    },
    {
      "id": 37139,
      "enchant": 44871,
      "gems": [
        39999
      ]
    },
    {
      "id": 34241,
      "enchant": 55002,
      "gems": [
        40014
      ]
    },
    {
      "id": 39558,
      "enchant": 44489,
      "gems": [
        39999,
        40014
      ]
    },
    {
      "id": 34448,
      "enchant": 44484,
      "gems": [
        39999,
        0
      ]
    },
    {
      "id": 39560,
      "enchant": 54999,
      "gems": [
        40014,
        0
      ]
    },
    {
      "id": 40694,
      "gems": [
        42702,
        39999
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
        39999
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
      "id": 37693,
      "enchant": 44492
    },
    {
      "id": 37856,
      "enchant": 44492
    },
    {
      "id": 44504,
      "gems": [
        40053
      ]
    }
  ]}`),
}

export const P1_PRESET_ASSASSINATION = {
	name: 'P1 Assassination',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	gear: EquipmentSpec.fromJsonString(`{ "items": [
    {
      "id": 40499,
      "enchant": 44879,
      "gems": [
        41398,
        42702
      ]
    },
    {
      "id": 44664,
      "gems": [
        40003
      ]
    },
    {
      "id": 40502,
      "enchant": 44871,
      "gems": [
        40003
      ]
    },
    {
      "id": 40403,
      "enchant": 55002
    },
    {
      "id": 40539,
      "enchant": 44489,
      "gems": [
        40003
      ]
    },
    {
      "id": 39765,
      "enchant": 44484,
      "gems": [
        40003,
        0
      ]
    },
    {
      "id": 40496,
      "enchant": 54999,
      "gems": [
        40053,
        0
      ]
    },
    {
      "id": 40260,
      "gems": [
        39999
      ]
    },
    {
      "id": 40500,
      "enchant": 38374,
      "gems": [
        40003,
        40003
      ]
    },
    {
      "id": 39701,
      "enchant": 55016
    },
    {
      "id": 40074
    },
    {
      "id": 40474
    },
    {
      "id": 40684
    },
    {
      "id": 44253
    },
    {
      "id": 39714,
      "enchant": 44492
    },
    {
      "id": 40386,
      "enchant": 44492
    },
    {
      "id": 40385
    }
  ]}`),
}

export const P1_PRESET_COMBAT = {
	name: 'P1 Combat',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	gear: EquipmentSpec.fromJsonString(`{"items": [
    {
      "id": 40499,
      "enchant": 44879,
      "gems": [
        41398,
        42702
      ]
    },
    {
      "id": 44664,
      "gems": [
        39999
      ]
    },
    {
      "id": 40502,
      "enchant": 44871,
      "gems": [
        39999
      ]
    },
    {
      "id": 40403,
      "enchant": 55002
    },
    {
      "id": 40539,
      "enchant": 44489,
      "gems": [
        39999
      ]
    },
    {
      "id": 39765,
      "enchant": 44484,
      "gems": [
        39999,
        0
      ]
    },
    {
      "id": 40541,
      "enchant": 54999,
      "gems": [
        0
      ]
    },
    {
      "id": 40205,
      "gems": [
        39999
      ]
    },
    {
      "id": 44011,
      "enchant": 38374,
      "gems": [
        39999,
        39999
      ]
    },
    {
      "id": 39701,
      "enchant": 55016
    },
    {
      "id": 40074
    },
    {
      "id": 40474
    },
    {
      "id": 40684
    },
    {
      "id": 44253
    },
    {
      "id": 40383,
      "enchant": 44492
    },
    {
      "id": 39714,
      "enchant": 44492
    },
    {
      "id": 40385
    }
  ]}`),
}