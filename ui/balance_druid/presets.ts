import {
    Consumes,
    Debuffs, Glyphs,
    IndividualBuffs,
    PartyBuffs,
    RaidBuffs,
    RaidTarget,
    TristateEffect
} from '../core/proto/common.js';
import {Flask} from '../core/proto/common.js';
import {Food} from '../core/proto/common.js';
import {EquipmentSpec} from '../core/proto/common.js';
import {Potions} from '../core/proto/common.js';
import {SavedTalents} from '../core/proto/ui.js';

import {
    BalanceDruid_Rotation as BalanceDruidRotation,
    BalanceDruid_Options as BalanceDruidOptions,
    BalanceDruid_Rotation_RotationType as RotationType, DruidMajorGlyph, DruidMinorGlyph,
} from '../core/proto/druid.js';

import * as Tooltips from '../core/constants/tooltips.js';
import {NO_TARGET} from "../core/proto_utils/utils";

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/wotlk/talent-calc and copy the numbers in the url.
export const StandardTalents = {
    name: 'Standard',
    data: SavedTalents.create({
        talentsString: '5012203125331203213305311231--205003012',
        glyphs: Glyphs.create({
            major1: DruidMajorGlyph.GlyphOfMoonfire,
            major2: DruidMajorGlyph.GlyphOfStarfire,
            major3: DruidMajorGlyph.GlyphOfStarfall,
            minor1: DruidMinorGlyph.GlyphOfTyphoon,
            minor2: DruidMinorGlyph.GlyphOfUnburdenedRebirth,
            minor3: DruidMinorGlyph.GlyphOfTheWild,
        }),
    }),
};

export const DefaultRotation = BalanceDruidRotation.create({
    type: RotationType.Adaptive,
});

export const DefaultOptions = BalanceDruidOptions.create({
    battleRes: true,
    innervateTarget: RaidTarget.create({
        targetIndex: NO_TARGET,
    }),
    isInsideEclipseThreshold: 15,
    mfInsideEclipseThreshold: 14,
    useIs: false,
    useMf: true,
});

export const DefaultConsumes = Consumes.create({
    defaultPotion: Potions.PotionOfSpeed,
    flask: Flask.FlaskOfTheFrostWyrm,
    food: Food.FoodFishFeast,
    prepopPotion: Potions.PotionOfSpeed,
    thermalSapper: true,
});

export const DefaultRaidBuffs = RaidBuffs.create({
    arcaneBrilliance: true,
    bloodlust: true,
    divineSpirit: true,
    giftOfTheWild: TristateEffect.TristateEffectImproved,
    icyTalons: true,
    moonkinAura: TristateEffect.TristateEffectImproved,
    leaderOfThePack: TristateEffect.TristateEffectImproved,
    powerWordFortitude: TristateEffect.TristateEffectImproved,
    sanctifiedRetribution: true,
    strengthOfEarthTotem: TristateEffect.TristateEffectImproved,
    trueshotAura: true,
    wrathOfAirTotem: true,
});

export const DefaultIndividualBuffs = IndividualBuffs.create({
    blessingOfKings: true,
    blessingOfMight: TristateEffect.TristateEffectImproved,
    blessingOfWisdom: TristateEffect.TristateEffectImproved,
    vampiricTouch: true,
});

export const DefaultPartyBuffs = PartyBuffs.create({
    heroicPresence: false,
});

export const DefaultDebuffs = Debuffs.create({
    bloodFrenzy: true,
    ebonPlaguebringer: true,
    faerieFire: TristateEffect.TristateEffectImproved,
    heartOfTheCrusader: true,
    judgementOfWisdom: true,
    shadowMastery: true,
    sunderArmor: true,
    totemOfWrath: true,
});

export const P1_PRESET = {
    name: 'P1 Preset',
    tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
    gear: EquipmentSpec.fromJsonString(`{"items": [
         {
          "id": 40467,
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
          "id": 40470,
          "enchant": 44874,
          "gems": [
            39998
          ]
        },
        {
          "id": 40405,
          "enchant": 44472
        },
        {
          "id": 40469,
          "enchant": 44489,
          "gems": [
            39998,
            40026
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
          "id": 40466,
          "enchant": 54999,
          "gems": [
            39998,
            0
          ]
        },
        {
          "id": 40561,
          "enchant": 54793,
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
          "id": 40080
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
          "id": 40192
        },
        {
          "id": 40321
        }
    ]}`),
};

export const PRE_RAID_PRESET = {
    name: 'Pre-raid Preset',
    tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
    gear: EquipmentSpec.fromJsonString(`{"items": [
          {
          "id": 42554,
          "enchant": 44877,
          "gems": [
            41285,
            40048
          ]
        },
        {
          "id": 40680
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
          "enchant": 44472
        },
        {
          "id": 39547,
          "enchant": 44489,
          "gems": [
            39998,
            40026
          ]
        },
        {
          "id": 37696,
          "enchant": 44498,
          "gems": [
            0
          ]
        },
        {
          "id": 39544,
          "enchant": 54999,
          "gems": [
            39998,
            0
          ]
        },
        {
          "id": 40696,
          "enchant": 54793,
          "gems": [
            40048,
            39998
          ]
        },
        {
          "id": 37791,
          "enchant": 41602
        },
        {
          "id": 44202,
          "enchant": 55016,
          "gems": [
            39998
          ]
        },
        {
          "id": 40585
        },
        {
          "id": 43253,
          "gems": [
            40026
          ]
        },
        {
          "id": 37873
        },
        {
          "id": 42987
        },
        {
          "id": 45085,
          "enchant": 44487
        },
        {
          "id": 40698
        },
        {
          "id": 32387
        }
    ]}`),
};

