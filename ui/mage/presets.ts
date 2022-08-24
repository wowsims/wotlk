import { Conjured } from '../core/proto/common.js';
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
	Mage,
	MageTalents as MageTalents,
	Mage_Rotation as MageRotation,
	Mage_Rotation_Type as RotationType,
	Mage_Rotation_PrimaryFireSpell as PrimaryFireSpell,
	Mage_Rotation_AoeRotation as AoeRotationSpells,
	Mage_Options as MageOptions,
	Mage_Options_ArmorType as ArmorType,
	MageMajorGlyph,
	MageMinorGlyph,
} from '../core/proto/mage.js';

import * as Tooltips from '../core/constants/tooltips.js';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/wotlk/talent-calc and copy the numbers in the url.
export const ArcaneTalents = {
	name: 'Arcane',
	data: SavedTalents.create({
		talentsString: '23000513310033015032310250532-03-023303001',
		glyphs: Glyphs.create({
			major1: MageMajorGlyph.GlyphOfArcaneBlast,
			major2: MageMajorGlyph.GlyphOfArcaneMissiles,
			major3: MageMajorGlyph.GlyphOfMoltenArmor,
			minor1: MageMinorGlyph.GlyphOfSlowFall,
			minor2: MageMinorGlyph.GlyphOfFrostWard,
			minor3: MageMinorGlyph.GlyphOfBlastWave,
		}),
	}),
};
export const FireTalents = {
	name: 'Fire',
	data: SavedTalents.create({
		talentsString: '23000503110003-0055030011302331053120321351',
		glyphs: Glyphs.create({
			major1: MageMajorGlyph.GlyphOfFireball,
			major2: MageMajorGlyph.GlyphOfMoltenArmor,
			major3: MageMajorGlyph.GlyphOfLivingBomb,
			minor1: MageMinorGlyph.GlyphOfSlowFall,
			minor2: MageMinorGlyph.GlyphOfFrostWard,
			minor3: MageMinorGlyph.GlyphOfBlastWave,
		}),
	}),
};
export const FrostfireTalents = {
	name: 'Frostfire',
	data: SavedTalents.create({
		talentsString: '-0055030011302331053120321351-2033030310003',
		glyphs: Glyphs.create({
			major1: MageMajorGlyph.GlyphOfFrostfire,
			major2: MageMajorGlyph.GlyphOfMoltenArmor,
			major3: MageMajorGlyph.GlyphOfLivingBomb,
			minor1: MageMinorGlyph.GlyphOfSlowFall,
			minor2: MageMinorGlyph.GlyphOfFrostWard,
			minor3: MageMinorGlyph.GlyphOfBlastWave,
		}),
	}),
}
export const FrostTalents = {
	name: 'Frost',
	data: SavedTalents.create({
		talentsString: '23000503110003--0533030310233100030152231351',
		glyphs: Glyphs.create({
			major1: MageMajorGlyph.GlyphOfFrostbolt,
			major2: MageMajorGlyph.GlyphOfEternalWater,
			major3: MageMajorGlyph.GlyphOfMoltenArmor,
			minor1: MageMinorGlyph.GlyphOfSlowFall,
			minor2: MageMinorGlyph.GlyphOfFrostWard,
			minor3: MageMinorGlyph.GlyphOfBlastWave,
		}),
	}),
};

export const DefaultFireRotation = MageRotation.create({
	type: RotationType.Fire,
	primaryFireSpell: PrimaryFireSpell.Fireball,
	maintainImprovedScorch: false,
});

export const DefaultFireOptions = MageOptions.create({
	armor: ArmorType.MoltenArmor,
	focusMagicPercentUptime: 99,
});

export const DefaultFireConsumes = Consumes.create({
	flask: Flask.FlaskOfTheFrostWyrm,
	food: Food.FoodFirecrackerSalmon,
	defaultPotion: Potions.PotionOfSpeed,
	defaultConjured: Conjured.ConjuredFlameCap,
});

export const DefaultFrostRotation = MageRotation.create({
	type: RotationType.Frost,
	waterElementalDisobeyChance: 0.1,
});

export const DefaultFrostOptions = MageOptions.create({
	armor: ArmorType.MoltenArmor,
});

export const DefaultFrostConsumes = Consumes.create({
	defaultPotion: Potions.PotionOfSpeed,
	defaultConjured: Conjured.ConjuredFlameCap,
	flask: Flask.FlaskOfTheFrostWyrm,
	food: Food.FoodFishFeast,
});

export const DefaultArcaneRotation = MageRotation.create({
	type: RotationType.Arcane,
	minBlastBeforeMissiles: 4,
	num4StackBlastsToMissilesGamble: 12,
	num4StackBlastsToEarlyMissiles: 6,
	extraBlastsDuringFirstAp: 2,
});

export const DefaultArcaneOptions = MageOptions.create({
	armor: ArmorType.MoltenArmor,
	focusMagicPercentUptime: 99,
});

export const DefaultArcaneConsumes = Consumes.create({
	defaultPotion: Potions.PotionOfSpeed,
	defaultConjured: Conjured.ConjuredDarkRune,
	flask: Flask.FlaskOfTheFrostWyrm,
	food: Food.FoodFirecrackerSalmon,
});

export const OtherDefaults = {
	distanceFromTarget: 25,
};

export const P1_ARCANE_PRESET = {
	name: 'Wotlk P1 Arcane Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecMage>) => player.getRotation().type == RotationType.Arcane,
	gear: EquipmentSpec.fromJsonString(`{"items": [
        {
          "id": 40416,
          "enchant": 44877,
          "gems": [
            41285,
            39998
          ]
        },
        {
          "id": 44661,
          "gems": [
            40026
          ]
        },
        {
          "id": 40419,
          "enchant": 44874,
          "gems": [
            40051
          ]
        },
        {
          "id": 44005,
          "enchant": 55642,
          "gems": [
            40026
          ]
        },
        {
          "id": 44002,
          "enchant": 44489,
          "gems": [
            39998,
            39998
          ]
        },
        {
          "id": 44008,
          "enchant": 44498,
          "gems": [
            39998,
            0
          ]
        },
        {
          "id": 40415,
          "enchant": 54999,
          "gems": [
            39998,
            0
          ]
        },
        {
          "id": 40561,
          "gems": [
            39998
          ]
        },
        {
          "id": 40417,
          "enchant": 41602,
          "gems": [
            39998,
            40051
          ]
        },
        {
          "id": 40558,
          "enchant": 55016
        },
        {
          "id": 40719
        },
        {
          "id": 40399
        },
        {
          "id": 39229
        },
        {
          "id": 40255
        },
        {
          "id": 40396,
          "enchant": 44487
        },
        {
          "id": 40273
        },
        {
          "id": 39426
        }
      ]
    }`),
};

export const P1_FIRE_PRESET = {
	name: 'Wotlk P1 Fire Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecMage>) => player.getRotation().type == RotationType.Fire,
	gear: EquipmentSpec.fromJsonString(`{
		"items": [
		  {
			"id": 40416,
			"enchant": 44877,
			"gems": [
			  41285,
			  39998
			]
		  },
		  {
			"id": 44661,
			"gems": [
			  40026
			]
		  },
		  {
			"id": 40419,
			"enchant": 44874,
			"gems": [
			  40049
			]
		  },
		  {
			"id": 44005,
			"enchant": 55642,
			"gems": [
			  40026
			]
		  },
		  {
			"id": 40418,
			"enchant": 44489,
			"gems": [
			  39998,
			  40048
			]
		  },
		  {
			"id": 44008,
			"enchant": 44498,
			"gems": [
			  39998,
			  0
			]
		  },
		  {
			"id": 40415,
			"enchant": 54999,
			"gems": [
			  39998,
			  0
			]
		  },
		  {
			"id": 40301,
			"gems": [
			  39998
			]
		  },
		  {
			"id": 40560,
			"enchant": 41602
		  },
		  {
			"id": 40246,
			"enchant": 55016
		  },
		  {
			"id": 40399
		  },
		  {
			"id": 40719
		  },
		  {
			"id": 40255
		  },
		  {
			"id": 40432
		  },
		  {
			"id": 40396,
			"enchant": 44487
		  },
		  {
			"id": 40273
		  },
		  {
			"id": 39712
		  }
		]
	  }`),
};

export const P1_FROST_PRESET = {
	name: 'Wotlk P1 Frost Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecMage>) => player.getRotation().type == RotationType.Frost,
	gear: EquipmentSpec.fromJsonString(`{"items": [
        {
          "id": 40416,
          "enchant": 44877,
          "gems": [
            41285,
            39998
          ]
        },
        {
          "id": 44661,
          "gems": [
            40026
          ]
        },
        {
          "id": 40419,
          "enchant": 44874,
          "gems": [
            40051
          ]
        },
        {
          "id": 44005,
          "enchant": 55642,
          "gems": [
            40026
          ]
        },
        {
          "id": 40418,
          "enchant": 44489,
          "gems": [
            39998,
            40048
          ]
        },
        {
          "id": 44008,
          "enchant": 44498,
          "gems": [
            39998,
            0
          ]
        },
        {
          "id": 40415,
          "enchant": 54999,
          "gems": [
            39998,
            0
          ]
        },
        {
          "id": 40301,
          "gems": [
            39998
          ]
        },
        {
          "id": 40560,
          "enchant": 41602
        },
        {
          "id": 40558,
          "enchant": 55016
        },
        {
          "id": 40399
        },
        {
          "id": 40719
        },
        {
          "id": 40255
        },
        {
          "id": 40432
        },
        {
          "id": 40396,
          "enchant": 44487
        },
        {
          "id": 39766
        },
        {
          "id": 39712
        }
      ]
    }`),
};

export const P1_PRERAID_ARCANE_PRESET = {
	name: "Arcane Pre Raid Preset",
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecMage>) => player.getRotation().type == RotationType.Arcane,
	gear: EquipmentSpec.fromJsonString(`
	{
		"items": [
		  {
			"id": 42553,
			"enchant": 44877,
			"gems": [
			  41285,
			  40049
			]
		  },
		  {
			"id": 39472
		  },
		  {
			"id": 37673,
			"enchant": 44874,
			"gems": [
			  39998
			]
		  },
		  {
			"id": 41610,
			"enchant": 55642
		  },
		  {
			"id": 39492,
			"enchant": 44489,
			"gems": [
			  39998,
			  40049
			]
		  },
		  {
			"id": 37361,
			"enchant": 44498,
			"gems": [
			  0
			]
		  },
		  {
			"id": 39495,
			"enchant": 54999,
			"gems": [
			  39998,
			  0
			]
		  },
		  {
			"id": 40696,
			"gems": [
			  40049,
			  40026
			]
		  },
		  {
			"id": 37854,
			"enchant": 41602
		  },
		  {
			"id": 44202,
			"enchant": 55016,
			"gems": [
			  40026
			]
		  },
		  {
			"id": 40585
		  },
		  {
			"id": 37694
		  },
		  {
			"id": 37873
		  },
		  {
			"id": 40682
		  },
		  {
			"id": 37360,
			"enchant": 45059
		  },
		  {},
		  {
			"id": 37238
		  }
		]
	  }`
	  )
}

export const P1_PRERAID_FIRE_PRESET = {
	name: "Fire Pre Raid Preset",
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecMage>) => player.getRotation().type == RotationType.Fire,
	gear: EquipmentSpec.fromJsonString(`
	{
		"items": [
		  {
			"id": 42553,
			"enchant": 44877,
			"gems": [
			  41285,
			  40014
			]
		  },
		  {
			"id": 39472
		  },
		  {
			"id": 34210,
			"enchant": 44874,
			"gems": [
			  40049,
			  40014
			]
		  },
		  {
			"id": 41610,
			"enchant": 63765
		  },
		  {
			"id": 39492,
			"enchant": 44489,
			"gems": [
			  40049,
			  40014
			]
		  },
		  {
			"id": 37361,
			"enchant": 44498,
			"gems": [
			  0
			]
		  },
		  {
			"id": 39495,
			"enchant": 54999,
			"gems": [
			  40049,
			  0
			]
		  },
		  {
			"id": 40696,
			"gems": [
			  40014,
			  40026
			]
		  },
		  {
			"id": 37854,
			"enchant": 41602
		  },
		  {
			"id": 44202,
			"enchant": 60623,
			"gems": [
			  40026
			]
		  },
		  {
			"id": 40585
		  },
		  {
			"id": 42644,
			"gems": [
			  40049
			]
		  },
		  {
			"id": 37873
		  },
		  {
			"id": 40682
		  },
		  {
			"id": 45085,
			"enchant": 44487
		  },
		  {
			"id": 40698
		  },
		  {
			"id": 37177
		  }
		]
	  }
	`)
}


export const ICC_FFB_Preset = {
	name: 'Wotlk ICC Frostfire Bolt Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecMage>) => player.getRotation().type == RotationType.Fire,
	gear: EquipmentSpec.fromJsonString(`
		{"items": [
			{
			"id": 51281,
			"enchant": 44877,
			"gems": [
				41285,
				40133
			]
			},
			{
			"id": 50724,
			"gems": [
				40133
			]
			},
			{
			"id": 51284,
			"enchant": 44874,
			"gems": [
				40153
			]
			},
			{
			"id": 54583,
			"enchant": 55642,
			"gems": [
				40152
			]
			},
			{
			"id": 51283,
			"enchant": 44489,
			"gems": [
				40152,
				40133
			]
			},
			{
			"id": 54582,
			"enchant": 44498,
			"gems": [
				40152,
				0
			]
			},
			{
			"id": 51280,
			"enchant": 54999,
			"gems": [
				40133,
				0
			]
			},
			{
			"id": 50613,
			"gems": [
				40133,
				40153,
				40153
			]
			},
			{
			"id": 50694,
			"enchant": 41602,
			"gems": [
				40152,
				40133,
				40153
			]
			},
			{
			"id": 50699,
			"enchant": 55016,
			"gems": [
				40133,
				40152
			]
			},
			{
			"id": 50664,
			"gems": [
				40133
			]
			},
			{
			"id": 50398,
			"gems": [
				40153
			]
			},
			{
			"id": 50348
			},
			{
			"id": 50345
			},
			{
			"id": 50732,
			"gems": [
				40152
			]
			},
			{
			"id": 50719
			},
			{
			"id": 50684,
			"gems": [
				40153
			]
			}
		]
		}`
)}


export const ICC_Fireball_Preset = {
	name: 'Wotlk ICC Fire Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecMage>) => player.getRotation().type == RotationType.Fire,
	gear: EquipmentSpec.fromJsonString(`{
		"items": [
		  {
			"id": 51281,
			"enchant": 44877,
			"gems": [
			  41285,
			  40133
			]
		  },
		  {
			"id": 50658,
			"gems": [
			  40153
			]
		  },
		  {
			"id": 51284,
			"enchant": 44874,
			"gems": [
			  40153
			]
		  },
		  {
			"id": 54583,
			"enchant": 55642,
			"gems": [
			  40152
			]
		  },
		  {
			"id": 51283,
			"enchant": 44489,
			"gems": [
			  40152,
			  40133
			]
		  },
		  {
			"id": 54582,
			"enchant": 44498,
			"gems": [
			  40152,
			  0
			]
		  },
		  {
			"id": 51280,
			"enchant": 54999,
			"gems": [
			  40133,
			  0
			]
		  },
		  {
			"id": 50613,
			"gems": [
			  40133,
			  40153,
			  40153
			]
		  },
		  {
			"id": 50694,
			"enchant": 41602,
			"gems": [
			  40152,
			  40133,
			  40153
			]
		  },
		  {
			"id": 50699,
			"enchant": 55016,
			"gems": [
			  40133,
			  40152
			]
		  },
		  {
			"id": 50664,
			"gems": [
			  40133
			]
		  },
		  {
			"id": 50398,
			"gems": [
			  40153
			]
		  },
		  {
			"id": 50348
		  },
		  {
			"id": 50345
		  },
		  {
			"id": 50732,
			"gems": [
			  40152
			]
		  },
		  {
			"id": 50719
		  },
		  {
			"id": 50684,
			"gems": [
			  40153
			]
		  }
		]
	  }
	`)}


	export const P2_Arcane_Preset = {
		name: 'Ulduar Arcane rough BiS',
		tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
		enableWhen: (player: Player<Spec.SpecMage>) => player.getRotation().type == RotationType.Arcane,
		gear: EquipmentSpec.fromJsonString(`
		{
			"items": [
			  {
				"id": 46129,
				"enchant": 44877,
				"gems": [
				  41285,
				  39998
				]
			  },
			  {
				"id": 45243,
				"gems": [
				  39998
				]
			  },
			  {
				"id": 46134,
				"enchant": 44874,
				"gems": [
				  40051
				]
			  },
			  {
				"id": 45618,
				"enchant": 55642,
				"gems": [
				  40026
				]
			  },
			  {
				"id": 46130,
				"enchant": 44489,
				"gems": [
				  39998,
				  39998
				]
			  },
			  {
				"id": 45446,
				"enchant": 44498,
				"gems": [
				  39998,
				  39998
				]
			  },
			  {
				"id": 45665,
				"enchant": 54999,
				"gems": [
				  40026,
				  40051,
				  0
				]
			  },
			  {
				"id": 45619,
				"gems": [
				  40049,
				  40049,
				  39998
				]
			  },
			  {
				"id": 46133,
				"enchant": 41602,
				"gems": [
				  39998,
				  39998
				]
			  },
			  {
				"id": 45135,
				"enchant": 55016,
				"gems": [
				  39998,
				  40026
				]
			  },
			  {
				"id": 46046,
				"gems": [
				  40049
				]
			  },
			  {
				"id": 45495,
				"gems": [
				  39998
				]
			  },
			  {
				"id": 45518
			  },
			  {
				"id": 45490
			  },
			  {
				"id": 45620,
				"gems": [
				  39998
				]
			  },
			  {
				"id": 45271
			  },
			  {
				"id": 39712
			  }
			]
		  }
		`)}
	
	export const P2_Fire_Preset = {
		name: 'Ulduar Fire rough BiS',
		tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
		enableWhen: (player: Player<Spec.SpecMage>) => player.getRotation().type == RotationType.Fire,
		gear: EquipmentSpec.fromJsonString(`
		{
			"items": [
			  {
				"id": 46129,
				"enchant": 44877,
				"gems": [
				  41285,
				  39998
				]
			  },
			  {
				"id": 45133,
				"gems": [
				  40051
				]
			  },
			  {
				"id": 46134,
				"enchant": 44874,
				"gems": [
				  40051
				]
			  },
			  {
				"id": 45242,
				"enchant": 55642,
				"gems": [
				  40051
				]
			  },
			  {
				"id": 46130,
				"enchant": 44489,
				"gems": [
				  39998,
				  39998
				]
			  },
			  {
				"id": 45275,
				"enchant": 44498,
				"gems": [
				  39998
				]
			  },
			  {
				"id": 45665,
				"enchant": 54999,
				"gems": [
				  40026,
				  40051,
				  0
				]
			  },
			  {
				"id": 45619,
				"gems": [
				  40051,
				  40051,
				  39998
				]
			  },
			  {
				"id": 46133,
				"enchant": 41602,
				"gems": [
				  39998,
				  39998
				]
			  },
			  {
				"id": 45537,
				"enchant": 55016,
				"gems": [
				  39998,
				  40026
				]
			  },
			  {
				"id": 46046,
				"gems": [
				  40051
				]
			  },
			  {
				"id": 45495,
				"gems": [
				  39998
				]
			  },
			  {
				"id": 45308
			  },
			  {
				"id": 45518
			  },
			  {
				"id": 45620,
				"enchant": 44495,
				"gems": [
				  39998
				]
			  },
			  {
				"id": 45617
			  },
			  {
				"id": 39712
			  }
			]
		  }
		`)}