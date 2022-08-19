import {Consumes, Debuffs, IndividualBuffs, RaidBuffs, TristateEffect} from '../core/proto/common.js';
import { Flask } from '../core/proto/common.js';
import { Food } from '../core/proto/common.js';
import { EquipmentSpec } from '../core/proto/common.js';
import { Potions } from '../core/proto/common.js';
import { SavedTalents } from '../core/proto/ui.js';

import { BalanceDruid_Rotation as BalanceDruidRotation, BalanceDruid_Options as BalanceDruidOptions, BalanceDruid_Rotation_RotationType as RotationType } from '../core/proto/druid.js';

import * as Tooltips from '../core/constants/tooltips.js';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/wotlk/talent-calc and copy the numbers in the url.
export const StandardTalents = {
    name: 'Standard',
    data: SavedTalents.create({
        talentsString: '510022312503135231351--520033',
    }),
};

export const DefaultRotation = BalanceDruidRotation.create({
	type: RotationType.Adaptive,
});

export const DefaultOptions = BalanceDruidOptions.create({
	useIs: true,
	useMf: true
});

export const DefaultConsumes = Consumes.create({
	flask: Flask.FlaskOfTheFrostWyrm,
	food: Food.FoodFishFeast,
	defaultPotion: Potions.PotionOfSpeed,
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
	gear: EquipmentSpec.fromJsonString(`{"items": [
        {
            "id": 40467,
            "enchant": 50368,
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
            "enchant": 50338,
            "gems": [
             39998
            ]
        },
        {
            "id": 40469,
            "enchant":  60692,
            "gems": [
				 39998,
				 40026
            ]
        },
        {
            "id": 40561,
            "enchant":  54793,
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
            "id": 44008,
            "enchant":  60767,
            "gems": [
            	39998
            ]
        },
        {
            "id": 40466,
            "enchant":  54999,
            "gems": [ 
            	39998
            ]
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
            "id": 40405,
            "enchant": 55642
        },
        {
            "id": 40395,
            "enchant":  60714
        },
        {
            "id": 40192
        },
        {
            "id": 40321
        }
    ]}`),
};
