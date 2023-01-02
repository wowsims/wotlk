import { Conjured, Consumes } from '../core/proto/common.js';
import { CustomRotation, CustomSpell } from '../core/proto/common.js';
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
	PaladinAura,
	PaladinMajorGlyph,
	PaladinMinorGlyph,
	PaladinJudgement,
	HolyPaladin_Rotation as HolyPaladinRotation,
	HolyPaladin_Options as HolyPaladinOptions,
} from '../core/proto/paladin.js';

import * as Gems from '../core/proto_utils/gems.js';
import * as Tooltips from '../core/constants/tooltips.js';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/wotlk/talent-calc and copy the numbers in the url.

export const GenericAoeTalents = {
	name: 'Baseline Example',
	data: SavedTalents.create({
		talentsString: '-05005135200132311333312321-511302012003',
		glyphs: {
			major1: PaladinMajorGlyph.GlyphOfSealOfVengeance,
			major2: PaladinMajorGlyph.GlyphOfRighteousDefense,
			major3: PaladinMajorGlyph.GlyphOfDivinePlea,
			minor1: PaladinMinorGlyph.GlyphOfSenseUndead,
			minor2: PaladinMinorGlyph.GlyphOfLayOnHands,
			minor3: PaladinMinorGlyph.GlyphOfBlessingOfKings
		}
	}),
};

export const DefaultRotation = HolyPaladinRotation.create({
});

export const DefaultOptions = HolyPaladinOptions.create({
	aura: PaladinAura.DevotionAura,
	judgement: PaladinJudgement.NoJudgement,
});

export const DefaultConsumes = Consumes.create({
	defaultPotion: Potions.RunicManaPotion,
	flask: Flask.FlaskOfTheFrostWyrm,
	food: Food.FoodFishFeast,
});

export const PRERAID_PRESET = {
	name: 'Preraid Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecHolyPaladin>) => true,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 42549,
			"enchant": 3818,
			"gems": [
				41396,
				49110
			]
		},
		{
			"id": 40679
		},
		{
			"id": 37635,
			"enchant": 3852,
			"gems": [
				40015
			]
		},
		{
			"id": 44188,
			"enchant": 3605
		},
		{
			"id": 39638,
			"enchant": 1953,
			"gems": [
				36767,
				40089
			]
		},
		{
			"id": 37682,
			"enchant": 3850,
			"gems": [
				0
			]
		},
		{
			"id": 39639,
			"enchant": 3860,
			"gems": [
				36767,
				0
			]
		},
		{
			"id": 37379,
			"enchant": 3601,
			"gems": [
				40022,
				40008
			]
		},
		{
			"id": 37292,
			"enchant": 3822,
			"gems": [
				40089
			]
		},
		{
			"id": 44243,
			"enchant": 3606
		},
		{
			"id": 37186
		},
		{
			"id": 37257
		},
		{
			"id": 44063,
			"gems": [
				36767,
				40015
			]
		},
		{
			"id": 37220
		},
		{
			"id": 37179,
			"enchant": 2673
		},
		{
			"id": 43085,
			"enchant": 3849
		},
		{
			"id": 40707
		}
	]}`),
};

export const P1_PRESET = {
	name: 'P1 Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecHolyPaladin>) => true,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 40581,
			"enchant": 3818,
			"gems": [
				41380,
				36767
			]
		},
		{
			"id": 40387
		},
		{
			"id": 40584,
			"enchant": 3852,
			"gems": [
				40008
			]
		},
		{
			"id": 40410,
			"enchant": 3605
		},
		{
			"id": 40579,
			"enchant": 3832,
			"gems": [
				36767,
				40022
			]
		},
		{
			"id": 39764,
			"enchant": 3850,
			"gems": [
				0
			]
		},
		{
			"id": 40580,
			"enchant": 3860,
			"gems": [
				40008,
				0
			]
		},
		{
			"id": 39759,
			"enchant": 3601,
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
			"id": 39717,
			"enchant": 3606,
			"gems": [
				40089
			]
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
				40089
			]
		},
		{
			"id": 37220
		},
		{
			"id": 40345,
			"enchant": 3788
		},
		{
			"id": 40400,
			"enchant": 3849
		},
		{
			"id": 40707
		}
	]}`),
};

