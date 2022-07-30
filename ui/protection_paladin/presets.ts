import { Conjured, Consumes } from '/wotlk/core/proto/common.js';
import { EquipmentSpec } from '/wotlk/core/proto/common.js';
import { Flask } from '/wotlk/core/proto/common.js';
import { Food } from '/wotlk/core/proto/common.js';
import { Glyphs } from '/wotlk/core/proto/common.js';
import { ItemSpec } from '/wotlk/core/proto/common.js';
import { Potions } from '/wotlk/core/proto/common.js';
import { Spec } from '/wotlk/core/proto/common.js';
import { WeaponImbue } from '/wotlk/core/proto/common.js';
import { Faction } from '/wotlk/core/proto/common.js';
import { SavedTalents } from '/wotlk/core/proto/ui.js';
import { Player } from '/wotlk/core/player.js';

import {
	PaladinAura as PaladinAura,
	PaladinMajorGlyph,
	PaladinMinorGlyph,
	PaladinJudgement as PaladinJudgement,
	ProtectionPaladin_Rotation as ProtectionPaladinRotation,
	ProtectionPaladin_Options as ProtectionPaladinOptions,
} from '/wotlk/core/proto/paladin.js';

import * as Enchants from '/wotlk/core/constants/enchants.js';
import * as Gems from '/wotlk/core/proto_utils/gems.js';
import * as Tooltips from '/wotlk/core/constants/tooltips.js';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/wotlk/talent-calc and copy the numbers in the url.

export const GenericAoeTalents = {
	name: 'Ardent Defender',
	data: SavedTalents.create({
		talentsString: '-05005135203102321333312301-502300510003',
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

export const DefaultRotation = ProtectionPaladinRotation.create({
	prioritizeHolyShield: true,
});

export const DefaultOptions = ProtectionPaladinOptions.create({
	aura: PaladinAura.RetributionAura,
	judgement: PaladinJudgement.JudgementOfWisdom,
	damageTakenPerSecond: 0,
});

export const DefaultConsumes = Consumes.create({
	flask: Flask.FlaskOfBlindingLight,
	food: Food.FoodFishermansFeast,
	defaultPotion: Potions.IronshieldPotion,
	mainHandImbue: WeaponImbue.WeaponImbueSuperiorWizardOil,
});

export const P1_PRESET = {
	name: 'P1 Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecProtectionPaladin>) => true,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 29068,
			"enchant": 29186,
			"gems": [
				24033,
				25896
			]
		},
		{
			"id": 28516
		},
		{
			"id": 29070,
			"enchant": 28911,
			"gems": [
				24033,
				24033
			]
		},
		{
			"id": 27804,
			"enchant": 33148
		},
		{
			"id": 29066,
			"enchant": 27957,
			"gems": [
				24033,
				24033,
				24033
			]
		},
		{
			"id": 28502,
			"enchant": 22533,
			"gems": [
				24033
			]
		},
		{
			"id": 28518,
			"enchant": 28272,
			"gems": [
				24033,
				24033
			]
		},
		{
			"id": 28566,
			"gems": [
				24033,
				24033
			]
		},
		{
			"id": 28621,
			"enchant": 29536,
			"gems": [
				24033,
				24033,
				24033
			]
		},
		{
			"id": 30641,
			"enchant": 35297
		},
		{
			"id": 29279,
			"enchant": 22536
		},
		{
			"id": 28675,
			"enchant": 22536
		},
		{
			"id": 28528
		},
		{
			"id": 23836
		},
		{
			"id": 28802,
			"enchant": 22555
		},
		{
			"id": 28825,
			"enchant": 28282,
			"gems": [
				24033
			]
		},
		{
			"id": 29388
		}
	]}`),
};

export const P2_PRESET = {
	name: 'P2 Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecProtectionPaladin>) => true,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 30125,
			"enchant": 29186,
			"gems": [
				24033,
				25896
			]
		},
		{
			"id": 30007
		},
		{
			"id": 29070,
			"enchant": 28911,
			"gems": [
				24033,
				24033
			]
		},
		{
			"id": 29925,
			"enchant": 33148
		},
		{
			"id": 29066,
			"enchant": 27957,
			"gems": [
				24033,
				24033,
				24033
			]
		},
		{
			"id": 32515,
			"enchant": 22533
		},
		{
			"id": 30124,
			"enchant": 28272
		},
		{
			"id": 30096,
			"gems": [
				24033,
				24033
			]
		},
		{
			"id": 30126,
			"enchant": 29536,
			"gems": [
				24033
			]
		},
		{
			"id": 32267,
			"enchant": 35297,
			"gems": [
				24033,
				24033
			]
		},
		{
			"id": 30083,
			"enchant": 22536
		},
		{
			"id": 28407,
			"enchant": 22536
		},
		{
			"id": 29370
		},
		{
			"id": 28789
		},
		{
			"id": 30095,
			"enchant": 22555
		},
		{
			"id": 28825,
			"enchant": 28282,
			"gems": [
				24033
			]
		},
		{
			"id": 27917
		}
	]}`),
};

export const P3_PRESET = {
	name: 'P3 Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecProtectionPaladin>) => true,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 32521,
			"enchant": 29191,
			"gems": [
				25896,
				32196
			]
		},
		{
			"id": 32362
		},
		{
			"id": 30998,
			"enchant": 28911,
			"gems": [
				32200,
				32196
			]
		},
		{
			"id": 34010,
			"enchant": 33148
		},
		{
			"id": 30991,
			"enchant": 27957,
			"gems": [
				32196,
				32196,
				32221
			]
		},
		{
			"id": 32279,
			"enchant": 22534
		},
		{
			"id": 30985,
			"enchant": 33153,
			"gems": [
				32196
			]
		},
		{
			"id": 32342,
			"gems": [
				32200,
				32200
			]
		},
		{
			"id": 30995,
			"enchant": 24274,
			"gems": [
				32200
			]
		},
		{
			"id": 32245,
			"enchant": 35297,
			"gems": [
				32200,
				32200
			]
		},
		{
			"id": 32261,
			"enchant": 22536
		},
		{
			"id": 29172,
			"enchant": 22536
		},
		{
			"id": 31858
		},
		{
			"id": 32489
		},
		{
			"id": 30910,
			"enchant": 22555
		},
		{
			"id": 32375,
			"enchant": 28282
		},
		{
			"id": 32368
		}
	]}`),
};

export const P4_PRESET = {
	name: 'P4 Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecProtectionPaladin>) => true,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 32521,
			"enchant": 29191,
			"gems": [
				25896,
				32196
			]
		},
		{
			"id": 32362
		},
		{
			"id": 30998,
			"enchant": 28911,
			"gems": [
				32200,
				32196
			]
		},
		{
			"id": 33593,
			"enchant": 33148
		},
		{
			"id": 30991,
			"enchant": 27957,
			"gems": [
				32196,
				32196,
				32221
			]
		},
		{
			"id": 32232,
			"enchant": 22534
		},
		{
			"id": 30985,
			"enchant": 33153,
			"gems": [
				32196
			]
		},
		{
			"id": 32342,
			"gems": [
				32200,
				32200
			]
		},
		{
			"id": 30995,
			"enchant": 24274,
			"gems": [
				32200
			]
		},
		{
			"id": 32245,
			"enchant": 35297,
			"gems": [
				32200,
				32200
			]
		},
		{
			"id": 32261,
			"enchant": 22536
		},
		{
			"id": 29172,
			"enchant": 22536
		},
		{
			"id": 31858
		},
		{
			"id": 33829
		},
		{
			"id": 30910,
			"enchant": 22555
		},
		{
			"id": 32375,
			"enchant": 28282
		},
		{
			"id": 33504
		}
	]}`),
};

export const P5_PRESET = {
	name: 'P5 Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecProtectionPaladin>) => true,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 34401,
			"enchant": 29191,
			"gems": [
				35501,
				32200
			]
		},
		{
			"id": 34178
		},
		{
			"id": 30998,
			"enchant": 23549,
			"gems": [
				32200,
				32196
			]
		},
		{
			"id": 34190,
			"enchant": 35756
		},
		{
			"id": 34945,
			"enchant": 27957,
			"gems": [
				32223
			]
		},
		{
			"id": 34433,
			"enchant": 22533,
			"gems": [
				32200
			]
		},
		{
			"id": 30985,
			"enchant": 33153,
			"gems": [
				32215
			]
		},
		{
			"id": 34488,
			"gems": [
				32200
			]
		},
		{
			"id": 34382,
			"enchant": 24274,
			"gems": [
				32200,
				32200,
				32215
			]
		},
		{
			"id": 34947,
			"enchant": 22533,
			"gems": [
				32215
			]
		},
		{
			"id": 34889,
			"enchant": 22536
		},
		{
			"id": 29172,
			"enchant": 22536
		},
		{
			"id": 33829
		},
		{
			"id": 34473
		},
		{
			"id": 35014,
			"enchant": 22555
		},
		{
			"id": 34185,
			"enchant": 28282,
			"gems": [
				32215
			]
		},
		{
			"id": 33504
		}
	]}`),
};
