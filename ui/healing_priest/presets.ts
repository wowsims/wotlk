import { Consumes } from '../core/proto/common.js';
import { CustomRotation, CustomSpell } from '../core/proto/common.js';
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
	HealingPriest_Rotation_RotationType as RotationType,
	HealingPriest_Rotation_SpellOption as SpellOption,
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
export const DiscTalents = {
	name: 'Disc',
	data: SavedTalents.create({
		talentsString: '0503203130300512301313231251-2351010303',
		glyphs: Glyphs.create({
			major1: MajorGlyph.GlyphOfPowerWordShield,
			major2: MajorGlyph.GlyphOfFlashHeal,
			major3: MajorGlyph.GlyphOfPenance,
			minor1: MinorGlyph.GlyphOfFortitude,
			minor2: MinorGlyph.GlyphOfShadowfiend,
			minor3: MinorGlyph.GlyphOfFading,
		}),
	}),
};
export const HolyTalents = {
	name: 'Holy',
	data: SavedTalents.create({
		talentsString: '05032031103-234051032002152530004311051',
		glyphs: Glyphs.create({
			major1: MajorGlyph.GlyphOfPrayerOfHealing,
			major2: MajorGlyph.GlyphOfRenew,
			major3: MajorGlyph.GlyphOfCircleOfHealing,
			minor1: MinorGlyph.GlyphOfFortitude,
			minor2: MinorGlyph.GlyphOfShadowfiend,
			minor3: MinorGlyph.GlyphOfFading,
		}),
	}),
};

export const DiscDefaultRotation = Rotation.create({
	type: RotationType.Cycle,
	customRotation: CustomRotation.create({
		spells: [
			CustomSpell.create({ spell: SpellOption.PowerWordShield, castsPerMinute: 18 }),
			CustomSpell.create({ spell: SpellOption.Penance, castsPerMinute: 4 }),
			CustomSpell.create({ spell: SpellOption.PrayerOfMending, castsPerMinute: 2 }),
			CustomSpell.create({ spell: SpellOption.GreaterHeal, castsPerMinute: 1 }),
		],
	}),
});

export const HolyDefaultRotation = Rotation.create({
	type: RotationType.Cycle,
	customRotation: CustomRotation.create({
		spells: [
			CustomSpell.create({ spell: SpellOption.GreaterHeal, castsPerMinute: 10 }),
			CustomSpell.create({ spell: SpellOption.CircleOfHealing, castsPerMinute: 5 }),
			CustomSpell.create({ spell: SpellOption.Renew, castsPerMinute: 10 }),
			CustomSpell.create({ spell: SpellOption.PrayerOfMending, castsPerMinute: 2 }),
		],
	}),
});

export const DefaultOptions = Options.create({
	useInnerFire: true,
	useShadowfiend: true,
	rapturesPerMinute: 5,

	powerInfusionTarget: RaidTarget.create({
		targetIndex: NO_TARGET, // In an individual sim the 0-indexed player is ourself.
	}),
});

export const DefaultConsumes = Consumes.create({
	flask: Flask.FlaskOfTheFrostWyrm,
	food: Food.FoodFishFeast,
	defaultPotion: Potions.RunicManaInjector,
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
	vampiricTouch: true,
});

export const DefaultDebuffs = Debuffs.create({
});

export const DISC_PRERAID_PRESET = {
	name: 'Disc Preraid Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<any>) => player.getTalentTree() == 0,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 37294,
			"enchant": 3819,
			"gems": [
				41401,
				39998
			]
		},
		{
			"id": 40681
		},
		{
			"id": 37691,
			"enchant": 3809
		},
		{
			"id": 37630,
			"enchant": 3859
		},
		{
			"id": 39515,
			"enchant": 3832,
			"gems": [
				42144,
				42144
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
			"id": 39519,
			"enchant": 3246,
			"gems": [
				42144,
				0
			]
		},
		{
			"id": 40697,
			"enchant": 3601,
			"gems": [
				39998,
				39998
			]
		},
		{
			"id": 37622,
			"enchant": 3719
		},
		{
			"id": 44202,
			"enchant": 3606,
			"gems": [
				40027
			]
		},
		{
			"id": 44283
		},
		{
			"id": 37195
		},
		{
			"id": 37660
		},
		{
			"id": 42413,
			"gems": [
				40012,
				40012
			]
		},
		{
			"id": 37360,
			"enchant": 3854
		},
		{},
		{
			"id": 37238
		}
	]}`),
};

export const DISC_P1_PRESET = {
	name: 'Disc P1 Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<any>) => player.getTalentTree() == 0,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 40456,
			"enchant": 3819,
			"gems": [
				41401,
				39998
			]
		},
		{
			"id": 44657,
			"gems": [
				40047
			]
		},
		{
			"id": 40450,
			"enchant": 3809,
			"gems": [
				42144
			]
		},
		{
			"id": 40724,
			"enchant": 3859
		},
		{
			"id": 40194,
			"enchant": 3832,
			"gems": [
				42144
			]
		},
		{
			"id": 40741,
			"enchant": 2332,
			"gems": [
				0
			]
		},
		{
			"id": 40445,
			"enchant": 3246,
			"gems": [
				42144,
				0
			]
		},
		{
			"id": 40271,
			"enchant": 3601,
			"gems": [
				40027,
				39998
			]
		},
		{
			"id": 40398,
			"enchant": 3719,
			"gems": [
				39998,
				39998
			]
		},
		{
			"id": 40236,
			"enchant": 3606
		},
		{
			"id": 40108
		},
		{
			"id": 40433
		},
		{
			"id": 37835
		},
		{
			"id": 40258
		},
		{
			"id": 40395,
			"enchant": 3834
		},
		{
			"id": 40350
		},
		{
			"id": 40245
		}
	]}`),
};

export const HOLY_PRERAID_PRESET = {
	name: 'Holy Preraid Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<any>) => player.getTalentTree() != 0,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 42553,
			"enchant": 3819,
			"gems": [
				41401,
				42148
			]
		},
		{
			"id": 40681
		},
		{
			"id": 37655,
			"enchant": 3809
		},
		{
			"id": 37291,
			"enchant": 3831
		},
		{
			"id": 39515,
			"enchant": 3832,
			"gems": [
				40012,
				40012
			]
		},
		{
			"id": 37361,
			"enchant": 1119,
			"gems": [
				0
			]
		},
		{
			"id": 39519,
			"enchant": 3604,
			"gems": [
				40012,
				0
			]
		},
		{
			"id": 40697,
			"enchant": 3601,
			"gems": [
				42148,
				42148
			]
		},
		{
			"id": 37189,
			"enchant": 3719,
			"gems": [
				40047,
				49110
			]
		},
		{
			"id": 44202,
			"enchant": 3606,
			"gems": [
				40092
			]
		},
		{
			"id": 44283
		},
		{
			"id": 37694
		},
		{
			"id": 37111
		},
		{
			"id": 42413,
			"gems": [
				40012,
				40012
			]
		},
		{
			"id": 37360,
			"enchant": 3854
		},
		{},
		{
			"id": 37619
		}
	]}`),
};

export const HOLY_P1_PRESET = {
	name: 'Holy P1 Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<any>) => player.getTalentTree() != 0,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 40447,
			"enchant": 3819,
			"gems": [
				41401,
				40051
			]
		},
		{
			"id": 44657,
			"gems": [
				40012
			]
		},
		{
			"id": 40450,
			"enchant": 3809,
			"gems": [
				40012
			]
		},
		{
			"id": 40723,
			"enchant": 3831
		},
		{
			"id": 40381,
			"enchant": 3832,
			"gems": [
				40012,
				49110
			]
		},
		{
			"id": 40741,
			"enchant": 1119,
			"gems": [
				0
			]
		},
		{
			"id": 40454,
			"enchant": 3604,
			"gems": [
				40051,
				0
			]
		},
		{
			"id": 40271,
			"enchant": 3601,
			"gems": [
				40012,
				40012
			]
		},
		{
			"id": 40398,
			"enchant": 3719,
			"gems": [
				42148,
				42148
			]
		},
		{
			"id": 40326,
			"enchant": 3606,
			"gems": [
				42148
			]
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
			"id": 42413,
			"gems": [
				40012,
				40012
			]
		},
		{
			"id": 40395,
			"enchant": 3834
		},
		{
			"id": 40350
		},
		{
			"id": 40245
		}
	]}`),
};

export const DISC_P2_PRESET = {
	name: 'Disc P2 Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<any>) => player.getTalentTree() == 0,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 46197,
			"enchant": 3819,
			"gems": [
				41401,
				45883
			]
		},
		{
			"id": 45933,
			"gems": [
				39998
			]
		},
		{
			"id": 45390,
			"enchant": 3809,
			"gems": [
				39998
			]
		},
		{
			"id": 45486,
			"enchant": 3859,
			"gems": [
				39998
			]
		},
		{
			"id": 46193,
			"enchant": 3832,
			"gems": [
				39998,
				39998
			]
		},
		{
			"id": 45146,
			"enchant": 3758,
			"gems": [
				0
			]
		},
		{
			"id": 45520,
			"enchant": 3246,
			"gems": [
				39998,
				39998,
				0
			]
		},
		{
			"id": 45558,
			"gems": [
				39998,
				40047,
				39998
			]
		},
		{
			"id": 45388,
			"enchant": 3721,
			"gems": [
				39998,
				39998
			]
		},
		{
			"id": 46050,
			"enchant": 3606,
			"gems": [
				39998,
				39998
			]
		},
		{
			"id": 45946,
			"gems": [
				39998
			]
		},
		{
			"id": 46096,
			"gems": [
				39998
			]
		},
		{
			"id": 40432
		},
		{
			"id": 45535
		},
		{
			"id": 45612,
			"enchant": 3834,
			"gems": [
				40027
			]
		},
		{
			"id": 45314
		},
		{
			"id": 45294,
			"gems": [
				39998
			]
		}
	]}`),
};

export const HOLY_P2_PRESET = {
	name: 'Holy P2 Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<any>) => player.getTalentTree() != 0,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 46197,
			"enchant": 3819,
			"gems": [
				41401,
				45883
			]
		},
		{
			"id": 45447,
			"gems": [
				40017
			]
		},
		{
			"id": 45390,
			"enchant": 3809,
			"gems": [
				40051
			]
		},
		{
			"id": 45486,
			"enchant": 3859,
			"gems": [
				40017
			]
		},
		{
			"id": 46193,
			"enchant": 3832,
			"gems": [
				40051,
				40017
			]
		},
		{
			"id": 44008,
			"enchant": 3758,
			"gems": [
				40051,
				0
			]
		},
		{
			"id": 45520,
			"enchant": 3246,
			"gems": [
				40017,
				40017,
				0
			]
		},
		{
			"id": 45558,
			"gems": [
				40051,
				40017,
				40017
			]
		},
		{
			"id": 45388,
			"enchant": 3721,
			"gems": [
				40026,
				40017
			]
		},
		{
			"id": 46050,
			"enchant": 3606,
			"gems": [
				40026,
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
			"id": 46323,
			"gems": [
				40017
			]
		},
		{
			"id": 46051
		},
		{
			"id": 45535
		},
		{
			"id": 46035,
			"enchant": 3834,
			"gems": [
				40017
			]
		},
		{
			"id": 45271
		},
		{
			"id": 45511
		}
	]}`),
};
