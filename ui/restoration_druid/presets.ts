import {
	Consumes,
	Debuffs,
	EquipmentSpec,
	Flask,
	Food,
	Glyphs,
	IndividualBuffs,
	PartyBuffs,
	Potions,
	RaidBuffs,
	RaidTarget,
	TristateEffect
} from '../core/proto/common.js';
import { SavedTalents } from '../core/proto/ui.js';

import {
	RestorationDruid_Options as RestorationDruidOptions,
	RestorationDruid_Rotation as RestorationDruidRotation,
	DruidMajorGlyph,
	DruidMinorGlyph,
} from '../core/proto/druid.js';

import * as Tooltips from '../core/constants/tooltips.js';
import { NO_TARGET } from "../core/proto_utils/utils";

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/wotlk/talent-calc and copy the numbers in the url.
export const CelestialFocusTalents = {
	name: 'Celestial Focus',
	data: SavedTalents.create({
		talentsString: '05320031103--230023312131502331050313051',
		glyphs: Glyphs.create({
			major1: DruidMajorGlyph.GlyphOfWildGrowth,
			major2: DruidMajorGlyph.GlyphOfSwiftmend,
			major3: DruidMajorGlyph.GlyphOfNourish,
			minor2: DruidMinorGlyph.GlyphOfUnburdenedRebirth,
			minor3: DruidMinorGlyph.GlyphOfTheWild,
			minor1: DruidMinorGlyph.GlyphOfDash,
		}),
	}),
};
export const ThiccRestoTalents = {
	name: 'Thicc Resto',
	data: SavedTalents.create({
		talentsString: '05320001--230023312331502531053313051',
		glyphs: Glyphs.create({
			major1: DruidMajorGlyph.GlyphOfWildGrowth,
			major2: DruidMajorGlyph.GlyphOfSwiftmend,
			major3: DruidMajorGlyph.GlyphOfNourish,
			minor2: DruidMinorGlyph.GlyphOfUnburdenedRebirth,
			minor3: DruidMinorGlyph.GlyphOfTheWild,
			minor1: DruidMinorGlyph.GlyphOfDash,
		}),
	}),
};

export const DefaultRotation = RestorationDruidRotation.create({
});

export const DefaultOptions = RestorationDruidOptions.create({
	innervateTarget: RaidTarget.create({
		targetIndex: NO_TARGET,
	}),
});

export const DefaultConsumes = Consumes.create({
	defaultPotion: Potions.RunicManaPotion,
	flask: Flask.FlaskOfTheFrostWyrm,
	food: Food.FoodFishFeast,
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

export const OtherDefaults = {
	distanceFromTarget: 18,
};

export const PRE_RAID_PRESET = {
	name: 'Pre-raid Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	gear: EquipmentSpec.fromJsonString(`{ "items": [
		{
			"id": 37149,
			"enchant": 3819,
			"gems": [
				41401,
				40051
			]
		},
		{
			"id": 42339,
			"gems": [
				40026
			]
		},
		{
			"id": 37673,
			"enchant": 3809,
			"gems": [
				39998
			]
		},
		{
			"id": 41610,
			"enchant": 3831
		},
		{
			"id": 42102,
			"enchant": 3832
		},
		{
			"id": 37361,
			"enchant": 2332,
			"gems": [
				0
			]
		},
		{
			"id": 42113,
			"enchant": 3246,
			"gems": [
				0
			]
		},
		{
			"id": 37643,
			"enchant": 3601,
			"gems": [
				39998
			]
		},
		{
			"id": 37791,
			"enchant": 3719
		},
		{
			"id": 44202,
			"enchant": 3232,
			"gems": [
				39998
			]
		},
		{
			"id": 37694
		},
		{
			"id": 37192
		},
		{
			"id": 37111
		},
		{
			"id": 37657
		},
		{
			"id": 37169,
			"enchant": 3834
		},
		{
			"id": 40699
		},
		{
			"id": 33508
		}
	]}`),
};

export const P1_PRESET = {
	name: 'P1 Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 44007,
			"enchant": 3819,
			"gems": [
				41401,
				40017
			]
		},
		{
			"id": 40071
		},
		{
			"id": 39719,
			"enchant": 3809,
			"gems": [
				39998
			]
		},
		{
			"id": 40723,
			"enchant": 3859
		},
		{
			"id": 44002,
			"enchant": 3832,
			"gems": [
				39998,
				40026
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
			"id": 40460,
			"enchant": 3246,
			"gems": [
				40017,
				0
			]
		},
		{
			"id": 40561,
			"enchant": 3601,
			"gems": [
				39998
			]
		},
		{
			"id": 40379,
			"enchant": 3719,
			"gems": [
				39998,
				40017
			]
		},
		{
			"id": 40558,
			"enchant": 3606
		},
		{
			"id": 40719
		},
		{
			"id": 40375
		},
		{
			"id": 37111
		},
		{
			"id": 40432
		},
		{
			"id": 40395,
			"enchant": 3834
		},
		{
			"id": 39766
		},
		{
			"id": 40342
		}
	]}`),
};

export const P2_PRESET = {
	name: 'P2 Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 46184,
			"enchant": 3819,
			"gems": [
				41401,
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
			"id": 46187,
			"enchant": 3809,
			"gems": [
				39998
			]
		},
		{
			"id": 45618,
			"enchant": 3831,
			"gems": [
				39998
			]
		},
		{
			"id": 45519,
			"enchant": 3832,
			"gems": [
				40017,
				39998,
				40026
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
			"id": 46183,
			"enchant": 3246,
			"gems": [
				39998,
				0
			]
		},
		{
			"id": 45616,
			"gems": [
				39998,
				39998,
				39998
			]
		},
		{
			"id": 46185,
			"enchant": 3719,
			"gems": [
				40026,
				39998
			]
		},
		{
			"id": 45135,
			"enchant": 3606,
			"gems": [
				39998,
				40017
			]
		},
		{
			"id": 45495,
			"gems": [
				40017
			]
		},
		{
			"id": 45946,
			"gems": [
				40017
			]
		},
		{
			"id": 45703
		},
		{
			"id": 45535
		},
		{
			"id": 46017,
			"enchant": 3834
		},
		{
			"id": 45271
		},
		{
			"id": 40342
		}
	]}`),
};
