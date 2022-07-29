import { Conjured, Consumes } from '/wotlk/core/proto/common.js';
import { EquipmentSpec } from '/wotlk/core/proto/common.js';
import { Flask } from '/wotlk/core/proto/common.js';
import { Food } from '/wotlk/core/proto/common.js';
import { Glyphs } from '/wotlk/core/proto/common.js';
import { ItemSpec } from '/wotlk/core/proto/common.js';
import { Potions } from '/wotlk/core/proto/common.js';
import { Spec } from '/wotlk/core/proto/common.js';
import { WeaponImbue } from '/wotlk/core/proto/common.js';
import { Faction } from '/wotlk/core/proto/common.js';
import { SavedTalents } from '/wotlk/core/proto/ui.js';
import { Player } from '/wotlk/core/player.js';

import {
	PaladinAura as PaladinAura,
	PaladinJudgement as PaladinJudgement,
	RetributionPaladin_Rotation as RetributionPaladinRotation,
	RetributionPaladin_Options as RetributionPaladinOptions,
	PaladinMajorGlyph,
	PaladinMinorGlyph,
} from '/wotlk/core/proto/paladin.js';

import * as Enchants from '/wotlk/core/constants/enchants.js';
import * as Gems from '/wotlk/core/proto_utils/gems.js';
import * as Tooltips from '/wotlk/core/constants/tooltips.js';

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
	exoSlack: 500,
	consSlack: 500,
});

export const DefaultOptions = RetributionPaladinOptions.create({
	aura: PaladinAura.RetributionAura,
	judgement: PaladinJudgement.JudgementOfWisdom,
	useDivinePlea: true,
	damageTakenPerSecond: 0,
});

export const DefaultConsumes = Consumes.create({
	defaultPotion: Potions.HastePotion,
	defaultConjured: Conjured.ConjuredDarkRune,
	flask: Flask.FlaskOfRelentlessAssault,
	food: Food.FoodRoastedClefthoof,
});

// Maybe use this later if I can figure out the interactive tooltips from tippy
const RET_BIS_DISCLAIMER = "<p>Please reference <a target=\"_blank\" href=\"https://docs.google.com/spreadsheets/d/1SxO6abYm4k7XRaP1MsxhaqYoukgyZ-cbWDE3ujadjx4/\">Baranor's TBC BiS Lists</a> for more detailed gearing options and information.</p>"

export const PRE_RAID_PRESET = {
	name: 'Pre-Raid Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecRetributionPaladin>) => true,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 32087,
			"enchant": 29192,
			"gems": [
				24058,
				32409
			]
		},
		{
			"id": 29119
		},
		{
			"id": 33173,
			"enchant": 28888,
			"gems": [
				24058,
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
			"id": 23522,
			"enchant": 24003
		},
		{
			"id": 23537,
			"enchant": 27899
		},
		{
			"id": 30341,
			"enchant": 33995
		},
		{
			"id": 27985,
			"gems": [
				24027,
				24054
			]
		},
		{
			"id": 30257,
			"enchant": 29535
		},
		{
			"id": 28176,
			"enchant": 22544,
			"gems": [
				24027,
				24054
			]
		},
		{
			"id": 29177
		},
		{
			"id": 30834
		},
		{
			"id": 29383
		},
		{
			"id": 28288
		},
		{
			"id": 28429,
			"enchant": 22559
		},
		{
			"id": 27484
		}
	]}`),
};

export const P1_PRESET = {
	name: 'P1 Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecRetributionPaladin>) => true,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 29073,
			"enchant": 29192,
			"gems": [
				24027,
				32409
			]
		},
		{
			"id": 28745
		},
		{
			"id": 29075,
			"enchant": 28888,
			"gems": [
				24058,
				24027
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
			"id": 29071,
			"enchant": 24003,
			"gems": [
				24027,
				24027,
				24027
			]
		},
		{
			"id": 28795,
			"enchant": 27899,
			"gems": [
				24054,
				24027
			]
		},
		{
			"id": 30644,
			"enchant": 33995
		},
		{
			"id": 28779,
			"gems": [
				24027,
				24054
			]
		},
		{
			"id": 30257,
			"enchant": 29535
		},
		{
			"id": 28608,
			"enchant": 22544,
			"gems": [
				24027,
				24058
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
			"id": 28429,
			"enchant": 22559
		},
		{
			"id": 27484
		}
	]}`),
};

export const P2_PRESET = {
	name: 'P2 Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecRetributionPaladin>) => true,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 32461,
			"enchant": 29192,
			"gems": [
				32409,
				24054
			]
		},
		{
			"id": 30022
		},
		{
			"id": 30055,
			"enchant": 28888,
			"gems": [
				24027
			]
		},
		{
			"id": 30098,
			"enchant": 34004
		},
		{
			"id": 30129,
			"enchant": 24003,
			"gems": [
				24027,
				24058,
				24058
			]
		},
		{
			"id": 28795,
			"enchant": 27899,
			"gems": [
				24054,
				24027
			]
		},
		{
			"id": 29947,
			"enchant": 33995
		},
		{
			"id": 30106,
			"gems": [
				24027,
				24054
			]
		},
		{
			"id": 30257,
			"enchant": 29535
		},
		{
			"id": 30104,
			"enchant": 22544,
			"gems": [
				24054,
				24027
			]
		},
		{
			"id": 30061
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
			"id": 28430,
			"enchant": 22559
		},
		{
			"id": 27484
		}
	]}`),
};

export const P3_PRESET = {
	name: 'P3 Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecRetributionPaladin>) => true,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 32235,
			"enchant": 29192,
			"gems": [
				32409,
				32193
			]
		},
		{
			"id": 30022
		},
		{
			"id": 30055,
			"enchant": 28888,
			"gems": [
				32193
			]
		},
		{
			"id": 33122,
			"enchant": 34004,
			"gems": [
				32193
			]
		},
		{
			"id": 30905,
			"enchant": 24003,
			"gems": [
				32211,
				32193,
				32217
			]
		},
		{
			"id": 32574,
			"enchant": 27899
		},
		{
			"id": 29947,
			"enchant": 33995
		},
		{
			"id": 30106,
			"gems": [
				32193,
				32211
			]
		},
		{
			"id": 30900,
			"enchant": 29535,
			"gems": [
				32193,
				32193,
				32193
			]
		},
		{
			"id": 32366,
			"enchant": 22544,
			"gems": [
				32193,
				32217
			]
		},
		{
			"id": 32526
		},
		{
			"id": 30834
		},
		{
			"id": 23206
		},
		{
			"id": 28830
		},
		{
			"id": 32332,
			"enchant": 22559
		},
		{
			"id": 27484
		}
	]}`),
};

export const P4_PRESET = {
	name: 'P4 Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecRetributionPaladin>) => true,
	gear: EquipmentSpec.fromJsonString(`
	{
		"items": [
			{
				"enchant": 44149,
				"gems": [
					41398,
					40111
				],
				"id": 50326
			},
			{
				"gems": [
					40142
				],
				"id": 54581
			},
			{
				"enchant": 44133,
				"gems": [
					40111
				],
				"id": 51160
			},
			{
				"enchant": 55777,
				"gems": [
					40111
				],
				"id": 50653
			},
			{
				"gems": [
					40111,
					40111,
					40118
				],
				"id": 45473
			},
			{
				"id": 52019
			},
			{
				"id": 23709
			},
			{
				"enchant": 44815,
				"gems": [
					40111
				],
				"id": 54580
			},
			{
				"enchant": 54999,
				"gems": [
					40111,
					40111
				],
				"id": 50188
			},
			{
				"gems": [
					40162,
					40143,
					40111
				],
				"id": 50762
			},
			{
				"enchant": 38374,
				"gems": [
					40162,
					49110
				],
				"id": 51161
			},
			{
				"enchant": 41118,
				"gems": [
					40111,
					40111
				],
				"id": 54578
			},
			{
				"gems": [
					40111
				],
				"id": 54576
			},
			{
				"gems": [
					40142
				],
				"id": 52572
			},
			{
				"id": 47131
			},
			{
				"id": 50343
			},
			{
				"enchant": 44493,
				"gems": [
					40111,
					40111,
					40111
				],
				"id": 49623
			},
			{
				"id": 47661
			}
		]
	}`),
};

export const P5_PRESET = {
	name: 'P5 Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecRetributionPaladin>) => true,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 34244,
			"enchant": 29192,
			"gems": [
				32409,
				32193
			]
		},
		{
			"id": 34177
		},
		{
			"id": 34388,
			"enchant": 28888,
			"gems": [
				32193,
				32217
			]
		},
		{
			"id": 34241,
			"enchant": 34004,
			"gems": [
				32193
			]
		},
		{
			"id": 34397,
			"enchant": 24003,
			"gems": [
				32211,
				32217,
				32193
			]
		},
		{
			"id": 34431,
			"enchant": 27899,
			"gems": [
				32193
			]
		},
		{
			"id": 34343,
			"enchant": 33995,
			"gems": [
				32193,
				32217
			]
		},
		{
			"id": 34485,
			"gems": [
				32193
			]
		},
		{
			"id": 34180,
			"enchant": 29535,
			"gems": [
				32211,
				32193,
				32217
			]
		},
		{
			"id": 34561,
			"enchant": 22544,
			"gems": [
				32193
			]
		},
		{
			"id": 34361
		},
		{
			"id": 34189
		},
		{
			"id": 34427
		},
		{
			"id": 34472
		},
		{
			"id": 34247,
			"enchant": 22559,
			"gems": [
				32193,
				32193,
				32193
			]
		},
		{
			"id": 27484
		}
	]}`),
};
