import { Consumes } from '/wotlk/core/proto/common.js';
import { EquipmentSpec } from '/wotlk/core/proto/common.js';
import { Flask } from '/wotlk/core/proto/common.js';
import { Food } from '/wotlk/core/proto/common.js';
import { Glyphs } from '/wotlk/core/proto/common.js';
import { ItemSpec } from '/wotlk/core/proto/common.js';
import { Potions } from '/wotlk/core/proto/common.js';
import { WeaponImbue } from '/wotlk/core/proto/common.js';
import { Faction } from '/wotlk/core/proto/common.js';
import { SavedTalents } from '/wotlk/core/proto/ui.js';
import { Player } from '/wotlk/core/player.js';

import { 
	ShadowPriest, 
	ShadowPriest_Rotation as Rotation, 
	ShadowPriest_Options as Options, 
	ShadowPriest_Rotation_RotationType,
	PriestMajorGlyph as MajorGlyph,
	PriestMinorGlyph as MinorGlyph,
} from '/wotlk/core/proto/priest.js';


import * as Enchants from '/wotlk/core/constants/enchants.js';
import * as Gems from '/wotlk/core/proto_utils/gems.js';
import * as Tooltips from '/wotlk/core/constants/tooltips.js';

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
});

export const DefaultConsumes = Consumes.create({
	flask: Flask.FlaskOfTheFrostWyrm,
	food: Food.FoodFishFeast,
	defaultPotion: Potions.PotionOfSpeed,
});

export const P5_PRESET = {
	name: 'P5 Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 34340,
			"enchant": 29191,
			"gems": [
				25893,
				32215
			]
		},
		{
			"id": 34204
		},
		{
			"id": 31070,
			"enchant": 28886,
			"gems": [
				32196,
				32196
			]
		},
		{
			"id": 34242,
			"enchant": 33150,
			"gems": [
				32196
			]
		},
		{
			"id": 34232,
			"enchant": 33990,
			"gems": [
				32196,
				32196,
				32196
			]
		},
		{
			"id": 34434,
			"enchant": 22534,
			"gems": [
				32196
			]
		},
		{
			"id": 34344,
			"enchant": 28272,
			"gems": [
				32196,
				32196
			]
		},
		{
			"id": 34528,
			"gems": [
				32196
			]
		},
		{
			"id": 34181,
			"enchant": 24274,
			"gems": [
				32196,
				32196,
				32196
			]
		},
		{
			"id": 34563,
			"enchant": 35297,
			"gems": [
				32196
			]
		},
		{
			"id": 34230,
			"enchant": 22536
		},
		{
			"id": 32527,
			"enchant": 22536
		},
		{
			"id": 33829
		},
		{
			"id": 34429
		},
		{
			"id": 34336,
			"enchant": 22561
		},
		{
			"id": 34179
		},
		{
			"id": 34347,
			"gems": [
				32196
			]
		}
	]}`),
};
