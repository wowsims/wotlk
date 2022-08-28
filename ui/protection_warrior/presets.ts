import { Consumes } from '../core/proto/common.js';
import { EquipmentSpec } from '../core/proto/common.js';
import { Flask } from '../core/proto/common.js';
import { Food } from '../core/proto/common.js';
import { Potions } from '../core/proto/common.js';
import { Glyphs } from '../core/proto/common.js';
import { SavedTalents } from '../core/proto/ui.js';

import {
	WarriorShout,
	ProtectionWarrior_Rotation as ProtectionWarriorRotation,
	ProtectionWarrior_Rotation_DemoShout as DemoShout,
	ProtectionWarrior_Rotation_ThunderClap as ThunderClap,
	ProtectionWarrior_Options as ProtectionWarriorOptions,
	WarriorMajorGlyph,
	WarriorMinorGlyph,
} from '../core/proto/warrior.js';

import * as Tooltips from '../core/constants/tooltips.js';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/wotlk/talent-calc and copy the numbers in the url.
export const StandardTalents = {
	name: 'Standard',
	data: SavedTalents.create({
		talentsString: '2500030023-302-053351225000012521030113321',
		glyphs: Glyphs.create({
			major1: WarriorMajorGlyph.GlyphOfBlocking,
			major2: WarriorMajorGlyph.GlyphOfVigilance,
			major3: WarriorMajorGlyph.GlyphOfDevastate,
			minor1: WarriorMinorGlyph.GlyphOfCharge,
			minor2: WarriorMinorGlyph.GlyphOfThunderClap,
			minor3: WarriorMinorGlyph.GlyphOfCommand,
		}),
	}),
};

export const DefaultRotation = ProtectionWarriorRotation.create({
	demoShout: DemoShout.DemoShoutMaintain,
	thunderClap: ThunderClap.ThunderClapMaintain,
	hsRageThreshold: 30,
	useShieldBlock: true,
});

export const DefaultOptions = ProtectionWarriorOptions.create({
	shout: WarriorShout.WarriorShoutCommanding,
	precastShout: true,
	precastShoutSapphire: false,
	precastShoutT2: false,

	startingRage: 0,
});

export const DefaultConsumes = Consumes.create({
	flask: Flask.FlaskOfFortification,
	food: Food.FoodFishermansFeast,
	defaultPotion: Potions.IronshieldPotion,
});

export const P1_BALANCED_PRESET = {
	name: 'P1 Balanced Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	gear: EquipmentSpec.fromJsonString(`{ "items": [
        {
          "id": 40546,
          "enchant": 44878,
          "gems": [
            41380,
            40034
          ]
        },
        {
          "id": 40387
        },
        {
          "id": 39704,
          "enchant": 44957,
          "gems": [
            40008
          ]
        },
        {
          "id": 40252,
          "enchant": 55002
        },
        {
          "id": 40544,
          "enchant": 44489,
          "gems": [
            40008,
            40008
          ]
        },
        {
          "id": 39764,
          "enchant": 44944,
          "gems": [
            0
          ]
        },
        {
          "id": 40545,
          "enchant": 63770,
          "gems": [
            49110,
            0
          ]
        },
        {
          "id": 39759,
          "enchant": 54793,
          "gems": [
            40008,
            36767
          ]
        },
        {
          "id": 40589,
          "enchant": 38373
        },
        {
          "id": 40297,
          "enchant": 44491
        },
        {
          "id": 40370
        },
        {
          "id": 40718
        },
        {
          "id": 40257
        },
        {
          "id": 44063,
          "gems": [
            36767,
            40089
          ]
        },
        {
          "id": 40402,
          "enchant": 22559
        },
        {
          "id": 40400,
          "enchant": 44936
        },
        {
          "id": 41168,
          "gems": [
            36767
          ]
        }
      ]}`),
};
