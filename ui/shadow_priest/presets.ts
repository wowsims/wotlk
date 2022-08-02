import { Consumes } from '../core/proto/common.js';
import { EquipmentSpec } from '../core/proto/common.js';
import { Flask } from '../core/proto/common.js';
import { Food } from '../core/proto/common.js';
import { Glyphs } from '../core/proto/common.js';
import { ItemSpec } from '../core/proto/common.js';
import { Potions } from '../core/proto/common.js';
import { WeaponImbue } from '../core/proto/common.js';
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
	prepopPotion:  Potions.PotionOfWildMagic,
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

export const P1_PRESET = {
	name: 'P1 Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	gear: EquipmentSpec.fromJsonString(`{ "items": [
        {
          "id": 40562,
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
          "id": 40459,
          "enchant": 44874,
          "gems": [
            39998
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
          "enchant": 33990,
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
          "id": 40454,
          "enchant": 54999,
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
          "enchant": 41602
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
          "id": 40255
        },
        {
          "id": 40432
        },
        {
          "id": 40395,
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
export const PreBis_PRESET = {
	name: 'PreBis Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	gear: EquipmentSpec.fromJsonString(`{ "items": [
        {
          "id": 42553,
          "enchant": 44877,
          "gems": [
            41285,
            40051
          ]
        },
        {
          "id": 40680
        },
        {
          "id": 34210,
          "enchant": 44874,
          "gems": [
            39998,
            40051
          ]
        },
        {
          "id": 41610,
          "enchant": 55642
        },
        {
          "id": 39523,
          "enchant": 33990,
          "gems": [
            39998,
            40051
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
          "id": 39530,
          "enchant": 54999,
          "gems": [
            40049,
            0
          ]
        },
        {
          "id": 40696,
          "gems": [
            39998,
            40026
          ]
        },
        {
          "id": 34181,
          "enchant": 41602,
          "gems": [
            39998,
            39998,
            40049
          ]
        },
        {
          "id": 44202,
          "enchant": 60623,
          "gems": [
            40026
          ]
        },
        {
          "id": 37694
        },
        {
          "id": 37192
        },
        {
          "id": 40682
        },
        {
          "id": 37873
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
}`),
};