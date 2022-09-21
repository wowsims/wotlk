import { Conjured, Consumes } from '../core/proto/common.js';
import { CustomRotation, CustomSpell } from '../core/proto/common.js';
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
	PaladinAura as PaladinAura,
	PaladinMajorGlyph,
	PaladinMinorGlyph,
	PaladinJudgement as PaladinJudgement,
	ProtectionPaladin_Rotation_SpellOption as SpellOption,
	ProtectionPaladin_Rotation as ProtectionPaladinRotation,
	ProtectionPaladin_Options as ProtectionPaladinOptions,
} from '../core/proto/paladin.js';

import * as Gems from '../core/proto_utils/gems.js';
import * as Tooltips from '../core/constants/tooltips.js';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/wotlk/talent-calc and copy the numbers in the url.

export const GenericAoeTalents = {
	name: 'Baseline Example',
	data: SavedTalents.create({
		talentsString: '-05005135200132311333312321-511302012003',
		glyphs: {
			major1: PaladinMajorGlyph.GlyphOfSealOfVengeance,
			major2: PaladinMajorGlyph.GlyphOfRighteousDefense,
			major3: PaladinMajorGlyph.GlyphOfDivinePlea,
			minor1: PaladinMinorGlyph.GlyphOfSenseUndead,
			minor2: PaladinMinorGlyph.GlyphOfLayOnHands,
			minor3: PaladinMinorGlyph.GlyphOfBlessingOfKings
		}
	}),
};

export const DefaultRotation = ProtectionPaladinRotation.create({
	hammerFirst: false,
	useCustomPrio: false,
	customRotation: CustomRotation.create({
		spells: [
			CustomSpell.create({ spell: SpellOption.ShieldOfRighteousness }),
			CustomSpell.create({ spell: SpellOption.HammerOfTheRighteous }),
			CustomSpell.create({ spell: SpellOption.HolyShield }),
			CustomSpell.create({ spell: SpellOption.HammerOfWrath }),
			CustomSpell.create({ spell: SpellOption.Consecration }),
			CustomSpell.create({ spell: SpellOption.AvengersShield }),
			CustomSpell.create({ spell: SpellOption.JudgementOfWisdom }),
			CustomSpell.create({ spell: SpellOption.Exorcism })
		],
	}),
});

export const DefaultOptions = ProtectionPaladinOptions.create({
	aura: PaladinAura.RetributionAura,
	judgement: PaladinJudgement.JudgementOfWisdom,
	damageTakenPerSecond: 0,
});

export const DefaultConsumes = Consumes.create({
	flask: Flask.FlaskOfStoneblood,
	food: Food.FoodDragonfinFilet,
	defaultPotion: Potions.IndestructiblePotion,
	prepopPotion:  Potions.IndestructiblePotion,
});

export const P0_PRESET = {
	name: 'Preraid Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecProtectionPaladin>) => true,
	gear: EquipmentSpec.fromJsonString(`{"items": [
        {
          "id": 42549,
          "enchant": 44878,
          "gems": [
            41396,
            49110
          ]
        },
        {
          "id": 40679
        },
        {
          "id": 37635,
          "enchant": 44957,
          "gems": [
            40015
          ]
        },
        {
          "id": 44188,
          "enchant": 55002
        },
        {
          "id": 39638,
          "enchant": 47766,
          "gems": [
            36767,
            40089
          ]
        },
        {
          "id": 37682,
          "enchant": 44944,
          "gems": [
            0
          ]
        },
        {
          "id": 39639,
          "enchant": 63770,
          "gems": [
            36767,
            0
          ]
        },
        {
          "id": 37379,
          "enchant": 54793,
          "gems": [
            40022,
            40008
          ]
        },
        {
          "id": 37292,
          "enchant": 38373,
          "gems": [
            40089
          ]
        },
        {
          "id": 44243,
          "enchant": 55016
        },
        {
          "id": 37186
        },
        {
          "id": 37257
        },
        {
          "id": 44063,
          "gems": [
            36767,
            40015
          ]
        },
        {
          "id": 37220
        },
        {
          "id": 37179,
          "enchant": 22559
        },
        {
          "id": 43085,
          "enchant": 44936
        },
        {
          "id": 40707
        }
      ]}`),
};

export const P1_PRESET = {
	name: 'P1 Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecProtectionPaladin>) => true,
	gear: EquipmentSpec.fromJsonString(`{"items": [
        {
          "id": 40581,
          "enchant": 44878,
          "gems": [
            41396,
            36767
          ]
        },
        {
          "id": 40387
        },
        {
          "id": 40584,
          "enchant": 44957,
          "gems": [
            49110
          ]
        },
        {
          "id": 40410,
          "enchant": 55002
        },
        {
          "id": 40579,
          "enchant": 44489,
          "gems": [
            36767,
            40022
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
          "id": 40580,
          "enchant": 63770,
          "gems": [
            40008,
            0
          ]
        },
        {
          "id": 39759,
          "enchant": 54793,
          "gems": [
            40008,
            40008
          ]
        },
        {
          "id": 40589,
          "enchant": 38373
        },
        {
          "id": 39717,
          "enchant": 55016,
          "gems": [
            40089
          ]
        },
        {
          "id": 40718
        },
        {
          "id": 40107
        },
        {
          "id": 44063,
          "gems": [
            36767,
            40089
          ]
        },
        {
          "id": 37220
        },
        {
          "id": 40345,
          "enchant": 44496
        },
        {
          "id": 40400,
          "enchant": 44936
        },
        {
          "id": 40707
        }
      ]}`),
};

