import { Consumes } from '../core/proto/common.js';
import { BattleElixir } from '../core/proto/common.js';
import { GuardianElixir } from '../core/proto/common.js';
import { Food } from '../core/proto/common.js';
import { EquipmentSpec } from '../core/proto/common.js';
import { Potions } from '../core/proto/common.js';
import { Conjured } from '../core/proto/common.js';
import { Explosive } from '../core/proto/common.js';
import { RaidTarget } from '../core/proto/common.js';
import { Glyphs } from '../core/proto/common.js';
import { SavedTalents } from '../core/proto/ui.js';
import { NO_TARGET } from '../core/proto_utils/utils.js';

import {
	FeralTankDruid_Rotation as DruidRotation,
	FeralTankDruid_Options as DruidOptions,
	DruidMajorGlyph,
	DruidMinorGlyph,
} from '../core/proto/druid.js';

import * as Tooltips from '../core/constants/tooltips.js';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/wotlk/talent-calc and copy the numbers in the url.
export const StandardTalents = {
	name: 'Standard',
	data: SavedTalents.create({
		talentsString: '-503232132322010353120300313511-20350001',
		glyphs: Glyphs.create({
			major1: DruidMajorGlyph.GlyphOfMaul,
			major2: DruidMajorGlyph.GlyphOfSurvivalInstincts,
			major3: DruidMajorGlyph.GlyphOfFrenziedRegeneration,
			minor1: DruidMinorGlyph.GlyphOfChallengingRoar,
			minor2: DruidMinorGlyph.GlyphOfThorns,
			minor3: DruidMinorGlyph.GlyphOfUnburdenedRebirth,
		}),
	}),
};

export const DefaultRotation = DruidRotation.create({
	maulRageThreshold: 25,
	maintainDemoralizingRoar: true,
	lacerateTime: 8.0,
});

export const DefaultOptions = DruidOptions.create({
	innervateTarget: RaidTarget.create({
		targetIndex: NO_TARGET,
	}),
	startingRage: 20,
});

export const DefaultConsumes = Consumes.create({
	battleElixir: BattleElixir.GurusElixir,
	guardianElixir: GuardianElixir.GiftOfArthas,
	food: Food.FoodBlackenedDragonfin,
	prepopPotion: Potions.IndestructiblePotion,
	defaultPotion: Potions.IndestructiblePotion,
	defaultConjured: Conjured.ConjuredHealthstone,
	thermalSapper: true,
	fillerExplosive: Explosive.ExplosiveSaroniteBomb,
});

export const P1_PRESET = {
	name: 'P1 Boss Tanking',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 40329,
			"enchant": 67839,
			"gems": [
				41339,
				40008
			]
		},
		{
			"id": 40387
		},
		{
			"id": 40494,
			"enchant": 44957,
			"gems": [
				40008
			]
		},
		{
			"id": 40252,
			"enchant": 3294
		},
		{
			"id": 40471,
			"enchant": 3832,
			"gems": [
				42702,
				40088
			]
		},
		{
			"id": 40186,
			"enchant": 3850,
			"gems": [
				40008,
				0
			]
		},
		{
			"id": 40472,
			"enchant": 63770,
			"gems": [
				40008,
				0
			]
		},
		{
			"id": 43591,
			"gems": [
				40008,
				40008,
				40008
			]
		},
		{
			"id": 44011,
			"enchant": 38373,
			"gems": [
				40008,
				40008
			]
		},
		{
			"id": 40243,
			"enchant": 55016,
			"gems": [
				40008
			]
		},
		{
			"id": 40370
		},
		{
			"id": 37784
		},
		{
			"id": 44253
		},
		{
			"id": 37220
		},
		{
			"id": 40280,
			"enchant": 2673
		},
		{},
		{
			"id": 38365
		}
	]}`),
};
