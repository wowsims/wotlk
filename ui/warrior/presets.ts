import { Consumes } from '../core/proto/common.js';
import { EquipmentSpec } from '../core/proto/common.js';
import { Flask } from '../core/proto/common.js';
import { Food } from '../core/proto/common.js';
import { Glyphs } from '../core/proto/common.js';
import { ItemSpec } from '../core/proto/common.js';
import { Potions } from '../core/proto/common.js';
import { Spec } from '../core/proto/common.js';
import { Faction } from '../core/proto/common.js';
import { SavedTalents } from '../core/proto/ui.js';
import { Player } from '../core/player.js';

import {
	WarriorShout,
	WarriorTalents as WarriorTalents,
	Warrior,
	Warrior_Rotation as WarriorRotation,
	Warrior_Rotation_SunderArmor as SunderArmor,
	Warrior_Options as WarriorOptions,
	WarriorMajorGlyph,
	WarriorMinorGlyph,
} from '../core/proto/warrior.js';

import * as Gems from '../core/proto_utils/gems.js';
import * as Tooltips from '../core/constants/tooltips.js';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/wotlk/talent-calc and copy the numbers in the url.
export const ArmsTalents = {
	name: 'Arms',
	data: SavedTalents.create({
		talentsString: '3022032023335100102012213231251-305-2033',
		glyphs: Glyphs.create({
			major1: WarriorMajorGlyph.GlyphOfRending,
			major2: WarriorMajorGlyph.GlyphOfMortalStrike,
			major3: WarriorMajorGlyph.GlyphOfExecution,
			minor1: WarriorMinorGlyph.GlyphOfBattle,
			minor2: WarriorMinorGlyph.GlyphOfCommand,
			minor3: WarriorMinorGlyph.GlyphOfCharge,
		}),
	}),
};

export const FuryTalents = {
	name: 'Fury',
	data: SavedTalents.create({
		talentsString: '30202300233-325003101504310053120500351',
		glyphs: Glyphs.create({
			major1: WarriorMajorGlyph.GlyphOfWhirlwind,
			major2: WarriorMajorGlyph.GlyphOfHeroicStrike,
			major3: WarriorMajorGlyph.GlyphOfExecution,
			minor1: WarriorMinorGlyph.GlyphOfBattle,
			minor2: WarriorMinorGlyph.GlyphOfCommand,
			minor3: WarriorMinorGlyph.GlyphOfCharge,
		}),
	}),
};

export const DefaultRotation = WarriorRotation.create({
	useRend: true,
  useMs: true,
  useCleave: false,

	prioritizeWw: true,
	sunderArmor: SunderArmor.SunderArmorHelpStack,

  msRageThreshold: 50,
	hsRageThreshold: 60,
	rendRageThresholdBelow: 70,
  slamRageThreshold: 15,
	rendCdThreshold: 1,
	useHsDuringExecute: true,
	useBtDuringExecute: true,
	useWwDuringExecute: true,
	useSlamOverExecute: true,
  spamExecute: true,
});

export const ArmsRotation = WarriorRotation.create({
	useRend: true,
	useMs: true,
  useCleave: false,
	sunderArmor: SunderArmor.SunderArmorHelpStack,
	msRageThreshold: 50,
  slamRageThreshold: 15,
	hsRageThreshold: 60,
	rendCdThreshold: 1,
	useHsDuringExecute: true,
	spamExecute: true,
});

export const DefaultOptions = WarriorOptions.create({
	startingRage: 0,
	useRecklessness: true,
	shout: WarriorShout.WarriorShoutCommanding,
});

export const DefaultConsumes = Consumes.create({
	flask: Flask.FlaskOfEndlessRage,
	food: Food.FoodDragonfinFilet,
	defaultPotion: Potions.IndestructiblePotion,
  prepopPotion:  Potions.IndestructiblePotion,
});

export const P1_PRERAID_FURY_PRESET = {
	name: 'P1 Pre-Raid Fury Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecWarrior>) => player.getTalents().bloodthirst,
	gear: EquipmentSpec.fromJsonString(`{"items": [
        {
          "id": 41386,
          "enchant": 44879,
          "gems": [
            41398,
            49110
          ]
        },
        {
          "id": 42645,
          "gems": [
            42142
          ]
        },
        {
          "id": 44195,
          "enchant": 44871
        },
        {
          "id": 37647,
          "enchant": 55002
        },
        {
          "id": 39606,
          "enchant": 44489,
          "gems": [
            42142,
            39996
          ]
        },
        {
          "id": 44203,
          "enchant": 44484,
          "gems": [
            0
          ]
        },
        {
          "id": 39609,
          "enchant": 54999,
          "gems": [
            40037,
            0
          ]
        },
        {
          "id": 40694,
          "gems": [
            42142,
            39996
          ]
        },
        {
          "id": 37193,
          "enchant": 38374,
          "gems": [
            39996,
            39996
          ]
        },
        {
          "id": 44306,
          "enchant": 55016,
          "gems": [
            39996,
            39996
          ]
        },
        {
          "id": 42642,
          "gems": [
            40037
          ]
        },
        {
          "id": 37642
        },
        {
          "id": 42987
        },
        {
          "id": 40684
        },
        {
          "id": 37852,
          "enchant": 44492
        },
        {
          "id": 37852,
          "enchant": 44492
        },
        {
          "id": 37191,
          "enchant": 41167
        }
      ]}`),
};

export const P1_FURY_PRESET = {
	name: 'P1 Fury Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecWarrior>) => player.getTalents().bloodthirst,
	gear: EquipmentSpec.fromJsonString(`{ "items": [
        {
          "id": 40528,
          "enchant": 44879,
          "gems": [
            41398,
            39996
          ]
        },
        {
          "id": 44664,
          "gems": [
            39996
          ]
        },
        {
          "id": 40530,
          "enchant": 44871,
          "gems": [
            40058
          ]
        },
        {
          "id": 40403,
          "enchant": 55002
        },
        {
          "id": 40525,
          "enchant": 44489,
          "gems": [
            42142,
            49110
          ]
        },
        {
          "id": 40733,
          "enchant": 44484,
          "gems": [
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
          "id": 40317,
          "gems": [
            42142
          ]
        },
        {
          "id": 40529,
          "enchant": 38374,
          "gems": [
            39996,
            39996
          ]
        },
        {
          "id": 40591,
          "enchant": 55016
        },
        {
          "id": 43993,
          "gems": [
            39996
          ]
        },
        {
          "id": 40075
        },
        {
          "id": 42987
        },
        {
          "id": 40256
        },
        {
          "id": 40384,
          "enchant": 44492
        },
        {
          "id": 40384,
          "enchant": 44492
        },
        {
          "id": 40385
        }
      ]}`),
};

export const P1_PRERAID_ARMS_PRESET = {
	name: 'P1 Pre-Raid Arms Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecWarrior>) => player.getTalents().mortalStrike,
	gear: EquipmentSpec.fromJsonString(`{ "items": [
        {
          "id": 41386,
          "enchant": 44879,
          "gems": [
            41285,
            49110
          ]
        },
        {
          "id": 42645,
          "gems": [
            42142
          ]
        },
        {
          "id": 44195,
          "enchant": 44871
        },
        {
          "id": 37647,
          "enchant": 55002
        },
        {
          "id": 39606,
          "enchant": 44489,
          "gems": [
            42142,
            39996
          ]
        },
        {
          "id": 41355,
          "enchant": 44484,
          "gems": [
            0
          ]
        },
        {
          "id": 39609,
          "enchant": 54999,
          "gems": [
            40037,
            0
          ]
        },
        {
          "id": 40694,
          "gems": [
            42142,
            39996
          ]
        },
        {
          "id": 37193,
          "enchant": 38374,
          "gems": [
            39996,
            39996
          ]
        },
        {
          "id": 44306,
          "enchant": 55016,
          "gems": [
            42702,
            40037
          ]
        },
        {
          "id": 42642,
          "gems": [
            40037
          ]
        },
        {
          "id": 37642
        },
        {
          "id": 42987
        },
        {
          "id": 40684
        },
        {
          "id": 37852,
          "enchant": 44492
        },
        {},
        {
          "id": 44504,
          "enchant": 41167,
          "gems": [
            40038
          ]
        }
      ]}`),
};

export const P1_ARMS_PRESET = {
	name: 'P1 Arms Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecWarrior>) => player.getTalents().mortalStrike,
	gear: EquipmentSpec.fromJsonString(`{"items": [
        {
          "id": 40528,
          "enchant": 44879,
          "gems": [
            41398,
            39996
          ]
        },
        {
          "id": 44664,
          "gems": [
            39996
          ]
        },
        {
          "id": 40530,
          "enchant": 44871,
          "gems": [
            40058
          ]
        },
        {
          "id": 40403,
          "enchant": 55002
        },
        {
          "id": 40525,
          "enchant": 44489,
          "gems": [
            42142,
            42142
          ]
        },
        {
          "id": 40330,
          "enchant": 44484,
          "gems": [
            39996,
            39996
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
          "id": 40317,
          "gems": [
            42142
          ]
        },
        {
          "id": 40529,
          "enchant": 38374,
          "gems": [
            39996,
            49110
          ]
        },
        {
          "id": 40591,
          "enchant": 55016
        },
        {
          "id": 43993,
          "gems": [
            39996
          ]
        },
        {
          "id": 40717
        },
        {
          "id": 42987
        },
        {
          "id": 40256
        },
        {
          "id": 40384,
          "enchant": 44492
        },
        {},
        {
          "id": 40385
        }
      ]}`),
};
