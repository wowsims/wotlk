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
import { RaidTarget } from '../core/proto/common.js';
import { TristateEffect } from '../core/proto/common.js';
import { SavedTalents } from '../core/proto/ui.js';
import { Player } from '../core/player.js';
import { NO_TARGET } from '../core/proto_utils/utils.js';

import {
	HealingPriest_Rotation as Rotation,
	HealingPriest_Options as Options,
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
		talentsString: '0533203100300502331-005551002020152-324',
		glyphs: Glyphs.create({
			major1: MajorGlyph.GlyphOfSmite,
			major2: MajorGlyph.GlyphOfHolyNova,
			major3: MajorGlyph.GlyphOfShadowWordDeath,
			minor1: MinorGlyph.GlyphOfFortitude,
			minor2: MinorGlyph.GlyphOfShadowfiend,
			minor3: MinorGlyph.GlyphOfFading,
		}),
	}),
};

export const DefaultRotation = Rotation.create({
});

export const DefaultOptions = Options.create({
	useInnerFire: true,
	useShadowfiend: true,

	powerInfusionTarget: RaidTarget.create({
		targetIndex: NO_TARGET, // In an individual sim the 0-indexed player is ourself.
	}),
});

export const DefaultConsumes = Consumes.create({
	flask: Flask.FlaskOfTheFrostWyrm,
	food: Food.FoodFishFeast,
	defaultPotion: Potions.RunicManaPotion,
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

export const PRERAID_PRESET = {
	name: 'Preraid Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 42553,
			"enchant": 44877,
			"gems": [
				41307,
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
				40014,
				40014
			]
		},
		{
			"id": 41610,
			"enchant": 44472
		},
		{
			"id": 42102,
			"enchant": 33990
		},
		{
			"id": 40740,
			"enchant": 44498,
			"gems": [
				0
			]
		},
		{
			"id": 42113,
			"enchant": 44592,
			"gems": [
				0
			]
		},
		{
			"id": 40696,
			"gems": [
				40014,
				40014
			]
		},
		{
			"id": 34181,
			"enchant": 41602,
			"gems": [
				40014,
				40014,
				40014
			]
		},
		{
			"id": 40750,
			"enchant": 60623
		},
		{
			"id": 43253,
			"gems": [
				40014
			]
		},
		{
			"id": 40719
		},
		{
			"id": 37835
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
	]}`),
};

export const P1_PRESET = {
	name: 'P1 Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 40562,
			"enchant": 44877,
			"gems": [
				41307,
				40049
			]
		},
		{
			"id": 40374
		},
		{
			"id": 40555,
			"enchant": 44874
		},
		{
			"id": 41610,
			"enchant": 63765
		},
		{
			"id": 40526,
			"enchant": 33990,
			"gems": [
				40049
			]
		},
		{
			"id": 40325,
			"enchant": 44498,
			"gems": [
				0
			]
		},
		{
			"id": 40454,
			"enchant": 44592,
			"gems": [
				40049,
				0
			]
		},
		{
			"id": 40301,
			"gems": [
				40049
			]
		},
		{
			"id": 40560,
			"enchant": 41602
		},
		{
			"id": 40246,
			"enchant": 60623
		},
		{
			"id": 40399
		},
		{
			"id": 39389
		},
		{
			"id": 42129
		},
		{
			"id": 40382
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
	]}`),
};
