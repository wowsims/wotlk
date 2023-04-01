import { Consumes } from '../core/proto/common.js';
import { EquipmentSpec } from '../core/proto/common.js';
import { Flask } from '../core/proto/common.js';
import { Food } from '../core/proto/common.js';
import { Glyphs } from '../core/proto/common.js';
import { ItemSpec } from '../core/proto/common.js';
import { Potions } from '../core/proto/common.js';
import { Faction } from '../core/proto/common.js';
import { RaidBuffs } from '../core/proto/common.js';
import { IndividualBuffs } from '../core/proto/common.js';
import { Debuffs } from '../core/proto/common.js';
import { TristateEffect } from '../core/proto/common.js';
import { SavedTalents } from '../core/proto/ui.js';
import { Player } from '../core/player.js';

import {
	ShadowPriest_Rotation as Rotation,
	ShadowPriest_Options as Options,
	ShadowPriest_Rotation_RotationType,
	PriestMajorGlyph as MajorGlyph,
	PriestMinorGlyph as MinorGlyph,
} from '../core/proto/priest.js';


import * as Tooltips from '../core/constants/tooltips.js';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/wotlk/talent-calc and copy the numbers in the url.
export const StandardTalents = {
	name: 'Standard',
	data: SavedTalents.create({
		talentsString: '05032031--325023051223010323151301351',
		glyphs: Glyphs.create({
			major1: MajorGlyph.GlyphOfShadow,
			major2: MajorGlyph.GlyphOfMindFlay,
			major3: MajorGlyph.GlyphOfDispersion,
			minor1: MinorGlyph.GlyphOfFortitude,
			minor2: MinorGlyph.GlyphOfShadowProtection,
			minor3: MinorGlyph.GlyphOfShadowfiend,
		}),
	}),
};

export const DefaultRotation = Rotation.create({
	rotationType: ShadowPriest_Rotation_RotationType.Ideal,
});

export const DefaultOptions = Options.create({
	useShadowfiend: true,
	useMindBlast: true,
	useShadowWordDeath: true,
});

export const DefaultConsumes = Consumes.create({
	flask: Flask.FlaskOfTheFrostWyrm,
	food: Food.FoodFishFeast,
	defaultPotion: Potions.PotionOfWildMagic,
	prepopPotion: Potions.PotionOfWildMagic,
});

export const DefaultRaidBuffs = RaidBuffs.create({
	giftOfTheWild: TristateEffect.TristateEffectImproved,
	powerWordFortitude: TristateEffect.TristateEffectImproved,
	strengthOfEarthTotem: TristateEffect.TristateEffectImproved,
	arcaneBrilliance: true,
	divineSpirit: true,
	trueshotAura: true,
	leaderOfThePack: TristateEffect.TristateEffectImproved,
	icyTalons: true,
	totemOfWrath: true,
	moonkinAura: TristateEffect.TristateEffectImproved,
	wrathOfAirTotem: true,
	sanctifiedRetribution: true,
	bloodlust: true,
	demonicPact: 500,
});

export const DefaultIndividualBuffs = IndividualBuffs.create({
	blessingOfKings: true,
	blessingOfWisdom: TristateEffect.TristateEffectImproved,
	blessingOfMight: TristateEffect.TristateEffectImproved,
	vampiricTouch: true,
});

export const DefaultDebuffs = Debuffs.create({
	sunderArmor: true,
	faerieFire: TristateEffect.TristateEffectImproved,
	bloodFrenzy: true,
	ebonPlaguebringer: true,
	heartOfTheCrusader: true,
	judgementOfWisdom: true,
});

export const PreBis_PRESET = {
	name: 'PreBis Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 42553,
			"enchant": 3820,
			"gems": [
				41285,
				40049
			]
		},
		{
			"id": 40680
		},
		{
			"id": 34210,
			"enchant": 3810,
			"gems": [
				39998,
				40026
			]
		},
		{
			"id": 41610,
			"enchant": 3722
		},
		{
			"id": 43792,
			"enchant": 1144,
			"gems": [
				39998,
				40051
			]
		},
		{
			"id": 37361,
			"enchant": 2332,
			"gems": [
				0
			]
		},
		{
			"id": 39530,
			"enchant": 3604,
			"gems": [
				40049,
				0
			]
		},
		{
			"id": 40696,
			"gems": [
				40049,
				39998
			]
		},
		{
			"id": 37854,
			"enchant": 3719
		},
		{
			"id": 44202,
			"enchant": 3826,
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
			"id": 37835
		},
		{
			"id": 37873
		},
		{
			"id": 41384,
			"enchant": 3834
		},
		{
			"id": 40698
		},
		{
			"id": 37177
		}
  ]}`),
};
export const P1_PRESET = {
	name: 'P1 Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	gear: EquipmentSpec.fromJsonString(` {"items": [
		{
			"id": 40562,
			"enchant": 3820,
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
			"id": 40459,
			"enchant": 3810,
			"gems": [
				39998
			]
		},
		{
			"id": 44005,
			"enchant": 3722,
			"gems": [
				40026
			]
		},
		{
			"id": 44002,
			"enchant": 1144,
			"gems": [
				39998,
				39998
			]
		},
		{
			"id": 44008,
			"enchant": 2332,
			"gems": [
				39998,
				0
			]
		},
		{
			"id": 40454,
			"enchant": 3604,
			"gems": [
				40049,
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
			"id": 40560,
			"enchant": 3719
		},
		{
			"id": 40558,
			"enchant": 3606
		},
		{
			"id": 40719
		},
		{
			"id": 40399
		},
		{
			"id": 40255
		},
		{
			"id": 40432
		},
		{
			"id": 40395,
			"enchant": 3834
		},
		{
			"id": 40273
		},
		{
			"id": 39712
		}
  ]}`),
};
export const P2_PRESET = {
	name: 'P2 Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	gear: EquipmentSpec.fromJsonString(`{ "items": [
        {
          "id": 46172,
          "enchant": 3820,
          "gems": [
            41285,
            45883
          ]
        },
        {
          "id": 45243,
          "gems": [
            39998
          ]
        },
        {
          "id": 46165,
          "enchant": 3810,
          "gems": [
            39998
          ]
        },
        {
          "id": 45242,
          "enchant": 3722,
          "gems": [
            40049
          ]
        },
        {
          "id": 46168,
          "enchant": 1144,
          "gems": [
            39998,
            39998
          ]
        },
        {
          "id": 45446,
          "enchant": 2332,
          "gems": [
            39998,
            0
          ]
        },
        {
          "id": 45665,
          "enchant": 3604,
          "gems": [
            39998,
            39998,
            0
          ]
        },
        {
          "id": 45619,
          "enchant": 3601,
          "gems": [
            39998,
            39998,
            39998
          ]
        },
        {
          "id": 46170,
          "enchant": 3719,
          "gems": [
            39998,
            40049
          ]
        },
        {
          "id": 45135,
          "enchant": 3606,
          "gems": [
            39998,
            40049
          ]
        },
        {
          "id": 45495,
          "gems": [
            40026
          ]
        },
        {
          "id": 46046,
          "gems": [
            39998
          ]
        },
        {
          "id": 45518
        },
        {
          "id": 45466
        },
        {
          "id": 45620,
          "enchant": 3834,
          "gems": [
            40026
          ]
        },
        {
          "id": 45617
        },
        {
          "id": 45294,
          "gems": [
            39998
          ]
        }
      ]
    }`),
};
