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
	PaladinJudgement as PaladinJudgement,
	RetributionPaladin_Rotation as RetributionPaladinRotation,
	RetributionPaladin_Options as RetributionPaladinOptions,
	RetributionPaladin_Rotation_RotationType as RotationType,
	RetributionPaladin_Rotation_SpellOption as SpellOption,
	PaladinMajorGlyph,
	PaladinMinorGlyph,
} from '../core/proto/paladin.js';

import * as Gems from '../core/proto_utils/gems.js';
import * as Tooltips from '../core/constants/tooltips.js';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/wotlk/talent-calc and copy the numbers in the url.
export const AuraMasteryTalents = {
	name: 'Basic w/Aura Mastery+LoH buff',
	data: SavedTalents.create({
		talentsString: '050501-05-05232051203331302133231331',
		glyphs: Glyphs.create({
			major1: PaladinMajorGlyph.GlyphOfSealOfVengeance,
			major2: PaladinMajorGlyph.GlyphOfJudgement,
			major3: PaladinMajorGlyph.GlyphOfConsecration,
			minor1: PaladinMinorGlyph.GlyphOfSenseUndead,
			minor2: PaladinMinorGlyph.GlyphOfLayOnHands,
			minor3: PaladinMinorGlyph.GlyphOfBlessingOfKings
		})
	}),
};


export const DivineSacTalents = {
	name: 'Basic w/Dsac',
	data: SavedTalents.create({
		talentsString: '03-453201002-05222051203331302133201331',
		glyphs: Glyphs.create({
			major1: PaladinMajorGlyph.GlyphOfSealOfVengeance,
			major2: PaladinMajorGlyph.GlyphOfJudgement,
			major3: PaladinMajorGlyph.GlyphOfConsecration,
			minor1: PaladinMinorGlyph.GlyphOfSenseUndead,
			minor2: PaladinMinorGlyph.GlyphOfLayOnHands,
			minor3: PaladinMinorGlyph.GlyphOfBlessingOfKings
		})
	}),
};

export const DefaultRotation = RetributionPaladinRotation.create({
	type: RotationType.Standard,
	exoSlack: 500,
	consSlack: 500,
	useDivinePlea: true,
	divinePleaPercentage: 0.75,
	holyWrathThreshold: 4,
	sovTargets: 1,
	customRotation: CustomRotation.create({
		spells: [
			CustomSpell.create({ spell: SpellOption.JudgementOfWisdom }),
			CustomSpell.create({ spell: SpellOption.HammerOfWrath }),
			CustomSpell.create({ spell: SpellOption.CrusaderStrike }),
			CustomSpell.create({ spell: SpellOption.DivineStorm }),
			CustomSpell.create({ spell: SpellOption.Exorcism }),
			CustomSpell.create({ spell: SpellOption.Consecration })
		],
	}),
	customCastSequence: CustomRotation.create({
		spells: [
			CustomSpell.create({ spell: SpellOption.JudgementOfWisdom }),
			CustomSpell.create({ spell: SpellOption.CrusaderStrike }),
			CustomSpell.create({ spell: SpellOption.DivineStorm }),
			CustomSpell.create({ spell: SpellOption.Consecration }),
			CustomSpell.create({ spell: SpellOption.CrusaderStrike }),
			CustomSpell.create({ spell: SpellOption.Exorcism }),
			CustomSpell.create({ spell: SpellOption.JudgementOfWisdom }),
			CustomSpell.create({ spell: SpellOption.CrusaderStrike }),
			CustomSpell.create({ spell: SpellOption.DivineStorm }),
			CustomSpell.create({ spell: SpellOption.Consecration }),
			CustomSpell.create({ spell: SpellOption.CrusaderStrike }),
		],
	}),
});

export const DefaultOptions = RetributionPaladinOptions.create({
	aura: PaladinAura.RetributionAura,
	judgement: PaladinJudgement.JudgementOfWisdom,
	damageTakenPerSecond: 0,
});

export const DefaultConsumes = Consumes.create({
	defaultPotion: Potions.PotionOfSpeed,
	defaultConjured: Conjured.ConjuredDarkRune,
	flask: Flask.FlaskOfEndlessRage,
	food: Food.FoodFishFeast,
});

// Maybe use this later if I can figure out the interactive tooltips from tippy
const RET_BIS_DISCLAIMER = "<p>Please reference <a target=\"_blank\" href=\"https://docs.google.com/spreadsheets/d/1SxO6abYm4k7XRaP1MsxhaqYoukgyZ-cbWDE3ujadjx4/\">Baranor's TBC BiS Lists</a> for more detailed gearing options and information.</p>"

export const PRE_RAID_PRESET = {
	name: 'Pre-Raid Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecRetributionPaladin>) => true,
	gear: EquipmentSpec.fromJsonString(`{"items": [
        {
          "id": 41386,
          "enchant": 44879,
          "gems": [
            41398,
            40022
          ]
        },
        {
          "id": 40678
        },
        {
          "id": 34388,
          "enchant": 44131,
          "gems": [
            39996,
            40037
          ]
        },
        {
          "id": 37647,
          "enchant": 55002
        },
        {
          "id": 39633,
          "enchant": 44489,
          "gems": [
            39996,
            40038
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
          "id": 39634,
          "enchant": 54999,
          "gems": [
            39996,
            0
          ]
        },
        {
          "id": 40694,
          "gems": [
            39996,
            39996
          ]
        },
        {
          "id": 37193,
          "enchant": 38372,
          "gems": [
            39996,
            39996
          ]
        },
        {
          "id": 44297,
          "enchant": 55016
        },
        {
          "id": 40586
        },
        {
          "id": 37685
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
          "id": 37574
        }
      ]
    }`),
};

export const P1_PRESET = {
	name: 'P1 Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecRetributionPaladin>) => true,
	gear: EquipmentSpec.fromJsonString(`{
		"items": [
		  {
			"id": 40576,
			"enchant": 44879,
			"gems": [
			  41398,
			  40037
			]
		  },
		  {
			"id": 44664,
			"gems": [
			  42142
			]
		  },
		  {
			"id": 40578,
			"enchant": 44871,
			"gems": [
			  49110
			]
		  },
		  {
			"id": 40403,
			"enchant": 55002
		  },
		  {
			"id": 40574,
			"enchant": 44489,
			"gems": [
			  42142,
			  39996
			]
		  },
		  {
			"id": 40186,
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
			"id": 40205,
			"gems": [
			  39996
			]
		  },
		  {
			"id": 40577,
			"enchant": 38374,
			"gems": [
			  42142,
			  40038
			]
		  },
		  {
			"id": 39701,
			"enchant": 55016
		  },
		  {
			"id": 40075
		  },
		  {
			"id": 40474
		  },
		  {
			"id": 42987
		  },
		  {
			"id": 40431
		  },
		  {
			"id": 40384,
			"enchant": 44492
		  },
		  {},
		  {
			"id": 42852
		  }
		]
	  }`),
};

export const P2_PRESET = {
	name: 'P2 Preset (Placeholder)',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecRetributionPaladin>) => true,
	gear: EquipmentSpec.fromJsonString(`{
		"items": [
		  {
			"id": 46156,
			"enchant": 44879,
			"gems": [
			  41398,
			  40003
			]
		  },
		  {
			"id": 45517,
			"gems": [
			  39996
			]
		  },
		  {
			"id": 46152,
			"enchant": 44871,
			"gems": [
			  40059
			]
		  },
		  {
			"id": 45461,
			"enchant": 55002,
			"gems": [
			  39996
			]
		  },
		  {
			"id": 46154,
			"enchant": 44489,
			"gems": [
			  39996,
			  39996
			]
		  },
		  {
			"id": 45611,
			"enchant": 44484,
			"gems": [
			  39996,
			  0
			]
		  },
		  {
			"id": 46155,
			"enchant": 54999,
			"gems": [
			  39996,
			  0
			]
		  },
		  {
			"id": 46095,
			"gems": [
			  39996,
			  39996,
			  39996
			]
		  },
		  {
			"id": 45134,
			"enchant": 38374,
			"gems": [
			  42142,
			  42142,
			  42142
			]
		  },
		  {
			"id": 45599,
			"enchant": 55016,
			"gems": [
			  40003,
			  49110
			]
		  },
		  {
			"id": 45456,
			"gems": [
			  39996
			]
		  },
		  {
			"id": 45525
		  },
		  {
			"id": 45609
		  },
		  {
			"id": 45522
		  },
		  {
			"id": 45533,
			"enchant": 44492,
			"gems": [
			  39996,
			  39996
			]
		  },
		  {},
		  {
			"id": 42852
		  }
		]
	  }`),
};

export const P3_PRESET = {
	name: 'P3 Preset (Placeholder)',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecRetributionPaladin>) => true,
	gear: EquipmentSpec.fromJsonString(`{
		"items": [
		  {
			"id": 48614,
			"enchant": 44879,
			"gems": [
			  41398,
			  40111
			]
		  },
		  {
			"id": 53103,
			"gems": [
			  40162
			]
		  },
		  {
			"id": 47972,
			"enchant": 44871,
			"gems": [
			  40111,
			  40162
			]
		  },
		  {
			"id": 46971,
			"enchant": 55002,
			"gems": [
			  40111
			]
		  },
		  {
			"id": 47004,
			"enchant": 44489,
			"gems": [
			  40111,
			  40111,
			  40111
			]
		  },
		  {
			"id": 47155,
			"enchant": 44484,
			"gems": [
			  49110,
			  40111,
			  0
			]
		  },
		  {
			"id": 48615,
			"enchant": 54999,
			"gems": [
			  40162,
			  0
			]
		  },
		  {
			"id": 47460,
			"gems": [
			  40111,
			  40111,
			  40111
			]
		  },
		  {
			"id": 46975,
			"enchant": 38374,
			"gems": [
			  42142,
			  42142,
			  42142
			]
		  },
		  {
			"id": 47154,
			"enchant": 55016,
			"gems": [
			  40111,
			  40111
			]
		  },
		  {
			"id": 47075,
			"gems": [
			  40111
			]
		  },
		  {
			"id": 46966,
			"gems": [
			  40111
			]
		  },
		  {
			"id": 47131
		  },
		  {
			"id": 47115
		  },
		  {
			"id": 47519,
			"enchant": 44492,
			"gems": [
			  40111,
			  40111
			]
		  },
		  {},
		  {
			"id": 47661
		  }
		]
	  }`),
};

export const P4_PRESET = {
	name: 'P4 Preset (Placeholder)',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecRetributionPaladin>) => true,
	gear: EquipmentSpec.fromJsonString(`{
		"items": [
		  {
			"id": 51162,
			"enchant": 44879,
			"gems": [
			  41398,
			  40111
			]
		  },
		  {
			"id": 50421,
			"gems": [
			  40111
			]
		  },
		  {
			"id": 51160,
			"enchant": 44871,
			"gems": [
			  40111
			]
		  },
		  {
			"id": 49998,
			"enchant": 55002,
			"gems": [
			  40125
			]
		  },
		  {
			"id": 51164,
			"enchant": 44489,
			"gems": [
			  40111,
			  49110
			]
		  },
		  {
			"id": 51832,
			"enchant": 44484,
			"gems": [
			  42142,
			  0
			]
		  },
		  {
			"id": 50021,
			"enchant": 54999,
			"gems": [
			  40125,
			  40111,
			  0
			]
		  },
		  {
			"id": 51925,
			"gems": [
			  40111,
			  40111,
			  40125
			]
		  },
		  {
			"id": 51161,
			"enchant": 38374,
			"gems": [
			  42142,
			  42142
			]
		  },
		  {
			"id": 49895,
			"enchant": 55016,
			"gems": [
			  40125,
			  40111
			]
		  },
		  {
			"id": 51843,
			"gems": [
			  40111
			]
		  },
		  {
			"id": 49949,
			"gems": [
			  40111
			]
		  },
		  {
			"id": 50351
		  },
		  {
			"id": 50343
		  },
		  {
			"id": 49888,
			"enchant": 44492,
			"gems": [
			  40111,
			  40111
			]
		  },
		  {},
		  {
			"id": 50455
		  }
		]
	  }`),
};

export const P5_PRESET = {
	name: 'P5 Preset (Placeholder)',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecRetributionPaladin>) => true,
	gear: EquipmentSpec.fromJsonString(`{
		"items": [
		  {
			"id": 51277,
			"enchant": 44879,
			"gems": [
			  41398,
			  40111
			]
		  },
		  {
			"id": 54581,
			"gems": [
			  40111
			]
		  },
		  {
			"id": 51279,
			"enchant": 44871,
			"gems": [
			  40118
			]
		  },
		  {
			"id": 50653,
			"enchant": 55002,
			"gems": [
			  40125
			]
		  },
		  {
			"id": 51275,
			"enchant": 44489,
			"gems": [
			  40118,
			  49110
			]
		  },
		  {
			"id": 54580,
			"enchant": 44484,
			"gems": [
			  42142,
			  0
			]
		  },
		  {
			"id": 50690,
			"enchant": 54999,
			"gems": [
			  40125,
			  40111,
			  0
			]
		  },
		  {
			"id": 50707,
			"gems": [
			  40111,
			  40111,
			  40125
			]
		  },
		  {
			"id": 51161,
			"enchant": 38374,
			"gems": [
			  42142,
			  42142
			]
		  },
		  {
			"id": 54578,
			"enchant": 55016,
			"gems": [
			  40111,
			  40111
			]
		  },
		  {
			"id": 54576,
			"gems": [
			  40111
			]
		  },
		  {
			"id": 52572,
			"gems": [
			  40111
			]
		  },
		  {
			"id": 50706
		  },
		  {
			"id": 54590
		  },
		  {
			"id": 49623,
			"enchant": 44492,
			"gems": [
			  40111,
			  40111,
			  40111
			]
		  },
		  {},
		  {
			"id": 50455
		  }
		]
	  }`),
};
