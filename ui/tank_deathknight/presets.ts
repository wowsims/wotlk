import { Consumes } from '../core/proto/common.js';
import { EquipmentSpec } from '../core/proto/common.js';
import { Flask } from '../core/proto/common.js';
import { Food } from '../core/proto/common.js';
import { Glyphs } from '../core/proto/common.js';
import { Potions } from '../core/proto/common.js';
import { SavedTalents } from '../core/proto/ui.js';

import {
	TankDeathknight_Rotation as TankDeathKnightRotation,
	TankDeathknight_Options as TankDeathKnightOptions,
	DeathknightMajorGlyph,
	DeathknightMinorGlyph,
  TankDeathknight_Rotation_Opener as Opener,
  TankDeathknight_Rotation_OptimizationSetting as OptimizationSetting,
  TankDeathknight_Rotation_BloodSpell as BloodSpell,
  TankDeathknight_Rotation_Presence as Presence,
} from '../core/proto/deathknight.js';

import * as Tooltips from '../core/constants/tooltips.js';

export const BloodTalents = {
	name: 'Blood',
	data: SavedTalents.create({
		talentsString: '005510153330330220102013-3050505000023-005',
		glyphs: Glyphs.create({
			major1: DeathknightMajorGlyph.GlyphOfDisease,
			major2: DeathknightMajorGlyph.GlyphOfRuneStrike,
			major3: DeathknightMajorGlyph.GlyphOfDarkCommand,
			minor1: DeathknightMinorGlyph.GlyphOfHornOfWinter,
			minor2: DeathknightMinorGlyph.GlyphOfBloodTap,
			minor3: DeathknightMinorGlyph.GlyphOfRaiseDead,
		}),
	}),
};

export const DoubleBuffBloodTalents = {
	name: '2B Blood',
	data: SavedTalents.create({
		talentsString: '005510153330330220102013-3050505000023201-002',
		glyphs: Glyphs.create({
			major1: DeathknightMajorGlyph.GlyphOfDisease,
			major2: DeathknightMajorGlyph.GlyphOfRuneStrike,
			major3: DeathknightMajorGlyph.GlyphOfDarkCommand,
			minor1: DeathknightMinorGlyph.GlyphOfHornOfWinter,
			minor2: DeathknightMinorGlyph.GlyphOfBloodTap,
			minor3: DeathknightMinorGlyph.GlyphOfRaiseDead,
		}),
	}),
};

export const FrostTalents = {
	name: 'Frost',
	data: SavedTalents.create({
		talentsString: '005510003-3050535000023301030023310035-0052',
		glyphs: Glyphs.create({
			major1: DeathknightMajorGlyph.GlyphOfDisease,
			major2: DeathknightMajorGlyph.GlyphOfRuneStrike,
			major3: DeathknightMajorGlyph.GlyphOfDarkCommand,
			minor1: DeathknightMinorGlyph.GlyphOfHornOfWinter,
			minor2: DeathknightMinorGlyph.GlyphOfBloodTap,
			minor3: DeathknightMinorGlyph.GlyphOfRaiseDead,
		}),
	}),
};

export const DoubleBuffFrostTalents = {
	name: '2B Frost',
	data: SavedTalents.create({
		talentsString: '00551005303003002-305053500002330103002301-005',
		glyphs: Glyphs.create({
			major1: DeathknightMajorGlyph.GlyphOfDisease,
			major2: DeathknightMajorGlyph.GlyphOfRuneStrike,
			major3: DeathknightMajorGlyph.GlyphOfDarkCommand,
			minor1: DeathknightMinorGlyph.GlyphOfHornOfWinter,
			minor2: DeathknightMinorGlyph.GlyphOfBloodTap,
			minor3: DeathknightMinorGlyph.GlyphOfRaiseDead,
		}),
	}),
};

export const DefaultRotation = TankDeathKnightRotation.create({
  opener: Opener.Threat,
  optimizationSetting: OptimizationSetting.Hps,
  bloodSpell: BloodSpell.BloodStrike,
  presence: Presence.Frost,
});

export const DefaultOptions = TankDeathKnightOptions.create({
	startingRunicPower: 0,
});

export const DefaultConsumes = Consumes.create({
	flask: Flask.FlaskOfStoneblood,
	food: Food.FoodDragonfinFilet,
	defaultPotion: Potions.IndestructiblePotion,
	prepopPotion: Potions.IndestructiblePotion,
});

export const P1_BLOOD_PRESET = {
	name: 'P1 Blood',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	gear: EquipmentSpec.fromJsonString(`{ "items": [
		{
			"id": 40565,
			"enchant": 3878,
			"gems": [
				41380,
				36767
			]
		},
		{
			"id": 40387
		},
		{
			"id": 39704,
			"enchant": 3852,
			"gems": [
				40008
			]
		},
		{
			"id": 40252,
			"enchant": 3605
		},
		{
			"id": 40559,
			"gems": [
				40008,
				40022
			]
		},
		{
			"id": 40306,
			"enchant": 3850,
			"gems": [
				40008,
				0
			]
		},
		{
			"id": 40563,
			"enchant": 3860,
			"gems": [
				40008,
				0
			]
		},
		{
			"id": 39759,
			"gems": [
				40008,
				40008
			]
		},
		{
			"id": 40567,
			"enchant": 3822,
			"gems": [
				40008,
				40008
			]
		},
		{
			"id": 40297,
			"enchant": 3232
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
				36767
			]
		},
		{
			"id": 42341,
			"gems": [
				40008,
				40008
			]
		},
		{
			"id": 40406,
			"enchant": 3847
		},
		{},
		{
			"id": 40207
		}
  ]}`),
};

export const P2_BLOOD_PRESET = {
	name: 'P2 Blood',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	gear: EquipmentSpec.fromJsonString(`{ "items": [
		{
			"id": 46120,
			"enchant": 3878,
			"gems": [
			  41380,
			  36767
			]
		  },
		  {
			"id": 45485,
			"gems": [
			  40008
			]
		  },
		  {
			"id": 46122,
			"enchant": 3852,
			"gems": [
			  40008
			]
		  },
		  {
			"id": 45496,
			"enchant": 3605,
			"gems": [
			  40022
			]
		  },
		  {
			"id": 46118,
			"gems": [
			  36767,
			  36767
			]
		  },
		  {
			"id": 45111,
			"enchant": 3850,
			"gems": [
			  0
			]
		  },
		  {
			"id": 46119,
			"enchant": 3860,
			"gems": [
			  40008,
			  0
			]
		  },
		  {
			"id": 45551,
			"gems": [
			  40008,
			  40008,
			  40008
			]
		  },
		  {
			"id": 45594,
			"enchant": 3822,
			"gems": [
			  40008,
			  40008,
			  40008
			]
		  },
		  {
			"id": 45988,
			"enchant": 3232,
			"gems": [
			  40008,
			  40008
			]
		  },
		  {
			"id": 45471,
			"gems": [
			  40008
			]
		  },
		  {
			"id": 45326
		  },
		  {
			"id": 45158
		  },
		  {
			"id": 46021
		  },
		  {
			"id": 45533,
			"enchant": 3370,
			"gems": [
			  40008,
			  40008
			]
		  },
		  {},
		  {
			"id": 45144
		  }
  ]}`),
};

export const P1_FROST_PRESET = {
	name: 'P1 Frost',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	gear: EquipmentSpec.fromJsonString(`{ "items": [
		{
			"id": 40565,
			"enchant": 3878,
			"gems": [
			  41380,
			  36767
			]
		  },
		  {
			"id": 40387
		  },
		  {
			"id": 40568,
			"enchant": 3852,
			"gems": [
			  40008
			]
		  },
		  {
			"id": 40252,
			"enchant": 3605
		  },
		  {
			"id": 40559,
			"gems": [
			  40008,
			  40022
			]
		  },
		  {
			"id": 40306,
			"enchant": 3850,
			"gems": [
			  40008,
			  0
			]
		  },
		  {
			"id": 40563,
			"enchant": 3860,
			"gems": [
			  40008,
			  0
			]
		  },
		  {
			"id": 39759,
			"gems": [
			  40008,
			  40008
			]
		  },
		  {
			"id": 40589,
			"enchant": 3822
		  },
		  {
			"id": 40297,
			"enchant": 3232
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
			  36767
			]
		  },
		  {
			"id": 40257
		  },
		  {
			"id": 40345,
			"enchant": 3370
		  },
		  {
			"id": 40345,
			"enchant": 3368
		  },
		  {
			"id": 40714
		  }
  ]}`),
};

export const P2_FROST_PRESET = {
	name: 'P2 Frost',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	gear: EquipmentSpec.fromJsonString(`{ "items": [
		{
			"id": 46120,
			"enchant": 3878,
			"gems": [
			  41380,
			  36767
			]
		  },
		  {
			"id": 45485,
			"gems": [
			  40008
			]
		  },
		  {
			"id": 46122,
			"enchant": 3852,
			"gems": [
			  40008
			]
		  },
		  {
			"id": 45496,
			"enchant": 3605,
			"gems": [
			  40022
			]
		  },
		  {
			"id": 46118,
			"gems": [
			  36767,
			  36767
			]
		  },
		  {
			"id": 45111,
			"enchant": 3850,
			"gems": [
			  0
			]
		  },
		  {
			"id": 46119,
			"enchant": 3860,
			"gems": [
			  40008,
			  0
			]
		  },
		  {
			"id": 45551,
			"gems": [
			  40008,
			  40008,
			  40008
			]
		  },
		  {
			"id": 45594,
			"enchant": 3822,
			"gems": [
			  40008,
			  40008,
			  40008
			]
		  },
		  {
			"id": 45988,
			"enchant": 3232,
			"gems": [
			  40008,
			  40008
			]
		  },
		  {
			"id": 45471,
			"gems": [
			  40008
			]
		  },
		  {
			"id": 45326
		  },
		  {
			"id": 45158
		  },
		  {
			"id": 46021
		  },
		  {
			"id": 46097,
			"enchant": 3370,
			"gems": [
			  40008
			]
		  },
		  {
			"id": 46097,
			"enchant": 3368,
			"gems": [
			  40008
			]
		  },
		  {
			"id": 45144
		  }
  ]}`),
};
