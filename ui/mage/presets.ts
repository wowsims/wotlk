import { Conjured } from '../core/proto/common.js';
import { Consumes } from '../core/proto/common.js';
import { EquipmentSpec } from '../core/proto/common.js';
import { Flask } from '../core/proto/common.js';
import { Food } from '../core/proto/common.js';
import { Glyphs } from '../core/proto/common.js';
import { ItemSpec } from '../core/proto/common.js';
import { Potions } from '../core/proto/common.js';
import { Spec } from '../core/proto/common.js';
import { Faction } from '../core/proto/common.js';
import { RaidTarget } from '../core/proto/common.js';
import { SavedTalents } from '../core/proto/ui.js';
import { Player } from '../core/player.js';
import { NO_TARGET } from '../core/proto_utils/utils';

import {
	Mage,
	MageTalents as MageTalents,
	Mage_Rotation as MageRotation,
	Mage_Rotation_Type as RotationType,
	Mage_Rotation_PrimaryFireSpell as PrimaryFireSpell,
	Mage_Rotation_AoeRotation as AoeRotationSpells,
	Mage_Options as MageOptions,
	Mage_Options_ArmorType as ArmorType,
	MageMajorGlyph,
	MageMinorGlyph,
} from '../core/proto/mage.js';

import * as Tooltips from '../core/constants/tooltips.js';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/wotlk/talent-calc and copy the numbers in the url.
export const ArcaneTalents = {
	name: 'Arcane',
	data: SavedTalents.create({
		talentsString: '23000513310033015032310250532-03-023303001',
		glyphs: Glyphs.create({
			major1: MageMajorGlyph.GlyphOfArcaneBlast,
			major2: MageMajorGlyph.GlyphOfArcaneMissiles,
			major3: MageMajorGlyph.GlyphOfMoltenArmor,
			minor1: MageMinorGlyph.GlyphOfSlowFall,
			minor2: MageMinorGlyph.GlyphOfFrostWard,
			minor3: MageMinorGlyph.GlyphOfBlastWave,
		}),
	}),
};
export const FireTalents = {
	name: 'Fire',
	data: SavedTalents.create({
		talentsString: '23000503110003-0055030012303331053120301351',
		glyphs: Glyphs.create({
			major1: MageMajorGlyph.GlyphOfFireball,
			major2: MageMajorGlyph.GlyphOfMoltenArmor,
			major3: MageMajorGlyph.GlyphOfLivingBomb,
			minor1: MageMinorGlyph.GlyphOfSlowFall,
			minor2: MageMinorGlyph.GlyphOfFrostWard,
			minor3: MageMinorGlyph.GlyphOfBlastWave,
		}),
	}),
};
export const FrostfireTalents = {
	name: 'Frostfire',
	data: SavedTalents.create({
		talentsString: '-2305030012303331053120311351-023303031002',
		glyphs: Glyphs.create({
			major1: MageMajorGlyph.GlyphOfFrostfire,
			major2: MageMajorGlyph.GlyphOfMoltenArmor,
			major3: MageMajorGlyph.GlyphOfLivingBomb,
			minor1: MageMinorGlyph.GlyphOfSlowFall,
			minor2: MageMinorGlyph.GlyphOfFrostWard,
			minor3: MageMinorGlyph.GlyphOfBlastWave,
		}),
	}),
}
export const FrostTalents = {
	name: 'Frost',
	data: SavedTalents.create({
		talentsString: '23000503110003--0533030310233100030152231351',
		glyphs: Glyphs.create({
			major1: MageMajorGlyph.GlyphOfFrostbolt,
			major2: MageMajorGlyph.GlyphOfEternalWater,
			major3: MageMajorGlyph.GlyphOfMoltenArmor,
			minor1: MageMinorGlyph.GlyphOfSlowFall,
			minor2: MageMinorGlyph.GlyphOfFrostWard,
			minor3: MageMinorGlyph.GlyphOfBlastWave,
		}),
	}),
};

export const DefaultFireRotation = MageRotation.create({
	type: RotationType.Fire,
	primaryFireSpell: PrimaryFireSpell.Fireball,
	maintainImprovedScorch: false,
	pyroblastDelayMs: 50,
});

export const DefaultFireOptions = MageOptions.create({
	armor: ArmorType.MoltenArmor,
	focusMagicPercentUptime: 99,
	focusMagicTarget: RaidTarget.create({
		targetIndex: NO_TARGET,
	}),
	reactionTimeMs: 300,
	igniteMunching: true,
});

export const DefaultFireConsumes = Consumes.create({
	flask: Flask.FlaskOfTheFrostWyrm,
	food: Food.FoodFirecrackerSalmon,
	defaultPotion: Potions.PotionOfSpeed,
	defaultConjured: Conjured.ConjuredFlameCap,
});

export const DefaultFrostRotation = MageRotation.create({
	type: RotationType.Frost,
	waterElementalDisobeyChance: 0.1,
});

export const DefaultFrostOptions = MageOptions.create({
	armor: ArmorType.MoltenArmor,
	focusMagicTarget: RaidTarget.create({
		targetIndex: NO_TARGET,
	}),
	reactionTimeMs: 300,
});

export const DefaultFrostConsumes = Consumes.create({
	defaultPotion: Potions.PotionOfSpeed,
	defaultConjured: Conjured.ConjuredFlameCap,
	flask: Flask.FlaskOfTheFrostWyrm,
	food: Food.FoodFishFeast,
});

export const DefaultArcaneRotation = MageRotation.create({
	type: RotationType.Arcane,
	only3ArcaneBlastStacksBelowManaPercent: 0.15,
	blastWithoutMissileBarrageAboveManaPercent: 0.2,
	extraBlastsDuringFirstAp: 0,
	missileBarrageBelowArcaneBlastStacks: 0,
	missileBarrageBelowManaPercent: 0,
});

export const DefaultArcaneOptions = MageOptions.create({
	armor: ArmorType.MoltenArmor,
	focusMagicPercentUptime: 99,
	focusMagicTarget: RaidTarget.create({
		targetIndex: NO_TARGET,
	}),
	reactionTimeMs: 300,
});

export const DefaultArcaneConsumes = Consumes.create({
	defaultPotion: Potions.PotionOfSpeed,
	defaultConjured: Conjured.ConjuredDarkRune,
	flask: Flask.FlaskOfTheFrostWyrm,
	food: Food.FoodFirecrackerSalmon,
});

export const OtherDefaults = {
	distanceFromTarget: 25,
};

export const ARCANE_PRERAID_PRESET = {
	name: "Arcane Preraid Preset",
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecMage>) => player.getRotation().type == RotationType.Arcane,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 42553,
			"enchant": 3820,
			"gems": [
				41285,
				40049
			]
		},
		{
			"id": 39472
		},
		{
			"id": 37673,
			"enchant": 3810,
			"gems": [
				39998
			]
		},
		{
			"id": 41610,
			"enchant": 3722
		},
		{
			"id": 39492,
			"enchant": 3832,
			"gems": [
				39998,
				40049
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
			"id": 39495,
			"enchant": 3604,
			"gems": [
				39998,
				0
			]
		},
		{
			"id": 40696,
			"gems": [
				40049,
				40026
			]
		},
		{
			"id": 37854,
			"enchant": 3719
		},
		{
			"id": 44202,
			"enchant": 3606,
			"gems": [
				40026
			]
		},
		{
			"id": 40585
		},
		{
			"id": 37694
		},
		{
			"id": 37873
		},
		{
			"id": 40682
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

export const FIRE_PRERAID_PRESET = {
	name: "Fire Preraid Preset",
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecMage>) => player.getRotation().type == RotationType.Fire,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 42553,
			"enchant": 3820,
			"gems": [
				41285,
				40014
			]
		},
		{
			"id": 39472
		},
		{
			"id": 34210,
			"enchant": 3810,
			"gems": [
				40049,
				40014
			]
		},
		{
			"id": 41610,
			"enchant": 3859
		},
		{
			"id": 39492,
			"enchant": 3832,
			"gems": [
				40049,
				40014
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
			"id": 39495,
			"enchant": 3604,
			"gems": [
				40049,
				0
			]
		},
		{
			"id": 40696,
			"gems": [
				40014,
				40026
			]
		},
		{
			"id": 37854,
			"enchant": 3719
		},
		{
			"id": 44202,
			"enchant": 3826,
			"gems": [
				40026
			]
		},
		{
			"id": 40585
		},
		{
			"id": 42644,
			"gems": [
				40049
			]
		},
		{
			"id": 37873
		},
		{
			"id": 40682
		},
		{
			"id": 45085,
			"enchant": 3834
		},
		{
			"id": 40698
		},
		{
			"id": 37177
		}
	]}`),
};

export const ARCANE_P1_PRESET = {
	name: 'Arcane P1 Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecMage>) => player.getRotation().type == RotationType.Arcane,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 40416,
			"enchant": 3820,
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
			"id": 40419,
			"enchant": 3810,
			"gems": [
				40051
			]
		},
		{
			"id": 44005,
			"enchant": 3722,
			"gems": [
				40026
			]
		},
		{
			"id": 44002,
			"enchant": 3832,
			"gems": [
				39998,
				39998
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
			"id": 40415,
			"enchant": 3604,
			"gems": [
				39998,
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
			"id": 40417,
			"enchant": 3719,
			"gems": [
				39998,
				40051
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
			"id": 40399
		},
		{
			"id": 39229
		},
		{
			"id": 40255
		},
		{
			"id": 40396,
			"enchant": 3834
		},
		{
			"id": 40273
		},
		{
			"id": 39426
		}
	]}`),
};

export const FIRE_P1_PRESET = {
	name: 'Fire P1 Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecMage>) => player.getRotation().type == RotationType.Fire,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 40416,
			"enchant": 3820,
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
			"id": 40419,
			"enchant": 3810,
			"gems": [
				40049
			]
		},
		{
			"id": 44005,
			"enchant": 3722,
			"gems": [
				40026
			]
		},
		{
			"id": 40418,
			"enchant": 3832,
			"gems": [
				39998,
				40048
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
			"id": 40415,
			"enchant": 3604,
			"gems": [
				39998,
				0
			]
		},
		{
			"id": 40301,
			"gems": [
				39998
			]
		},
		{
			"id": 40560,
			"enchant": 3719
		},
		{
			"id": 40246,
			"enchant": 3606
		},
		{
			"id": 40399
		},
		{
			"id": 40719
		},
		{
			"id": 40255
		},
		{
			"id": 40432
		},
		{
			"id": 40396,
			"enchant": 3834
		},
		{
			"id": 40273
		},
		{
			"id": 39712
		}
	]}`),
};

export const FROST_P1_PRESET = {
	name: 'Frost P1 Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecMage>) => player.getRotation().type == RotationType.Frost,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 40416,
			"enchant": 3820,
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
			"id": 40419,
			"enchant": 3810,
			"gems": [
				40051
			]
		},
		{
			"id": 44005,
			"enchant": 3722,
			"gems": [
				40026
			]
		},
		{
			"id": 40418,
			"enchant": 3832,
			"gems": [
				39998,
				40048
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
			"id": 40415,
			"enchant": 3604,
			"gems": [
				39998,
				0
			]
		},
		{
			"id": 40301,
			"gems": [
				39998
			]
		},
		{
			"id": 40560,
			"enchant": 3719
		},
		{
			"id": 40558,
			"enchant": 3606
		},
		{
			"id": 40399
		},
		{
			"id": 40719
		},
		{
			"id": 40255
		},
		{
			"id": 40432
		},
		{
			"id": 40396,
			"enchant": 3834
		},
		{
			"id": 39766
		},
		{
			"id": 39712
		}
	]}`),
};

export const ARCANE_P2_PRESET = {
	name: 'Arcane P2 Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecMage>) => player.getRotation().type == RotationType.Arcane,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 45497,
			"enchant": 3820,
			"gems": [
				41285,
				45883
			]
		},
		{
			"id": 45243,
			"gems": [
				39998
			]
		},
		{
			"id": 46134,
			"enchant": 3810,
			"gems": [
				39998
			]
		},
		{
			"id": 45618,
			"enchant": 3722,
			"gems": [
				40026
			]
		},
		{
			"id": 46130,
			"enchant": 3832,
			"gems": [
				39998,
				39998
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
			"id": 45665,
			"enchant": 3604,
			"gems": [
				39998,
				39998,
				0
			]
		},
		{
			"id": 45619,
			"gems": [
				39998,
				39998,
				39998
			]
		},
		{
			"id": 45488,
			"enchant": 3719,
			"gems": [
				39998,
				40051,
				40026
			]
		},
		{
			"id": 45135,
			"enchant": 3606,
			"gems": [
				39998,
				39998
			]
		},
		{
			"id": 46046,
			"gems": [
				39998
			]
		},
		{
			"id": 45495,
			"gems": [
				39998
			]
		},
		{
			"id": 45466
		},
		{
			"id": 45518
		},
		{
			"id": 45620,
			"enchant": 3834,
			"gems": [
				39998
			]
		},
		{
			"id": 45617
		},
		{
			"id": 45294,
			"gems": [
				39998
			]
		}
	]}`),
};

export const FIRE_P2_PRESET = {
	name: 'Fire P2 Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecMage>) => player.getRotation().type == RotationType.Fire,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 46129,
			"enchant": 3820,
			"gems": [
				41285,
				45883
			]
		},
		{
			"id": 45133,
			"gems": [
				40048
			]
		},
		{
			"id": 46134,
			"enchant": 3810,
			"gems": [
				39998
			]
		},
		{
			"id": 45242,
			"enchant": 3722,
			"gems": [
				39998
			]
		},
		{
			"id": 46130,
			"enchant": 3832,
			"gems": [
				39998,
				39998
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
			"id": 45665,
			"enchant": 3604,
			"gems": [
				40026,
				40048,
				0
			]
		},
		{
			"id": 45619,
			"gems": [
				40048,
				40048,
				39998
			]
		},
		{
			"id": 46133,
			"enchant": 3719,
			"gems": [
				39998,
				39998
			]
		},
		{
			"id": 45537,
			"enchant": 3606,
			"gems": [
				39998,
				40026
			]
		},
		{
			"id": 45495,
			"gems": [
				39998
			]
		},
		{
			"id": 46046,
			"gems": [
				39998
			]
		},
		{
			"id": 45466
		},
		{
			"id": 45518
		},
		{
			"id": 45620,
			"enchant": 3834,
			"gems": [
				39998
			]
		},
		{
			"id": 45617
		},
		{
			"id": 45294,
			"gems": [
				39998
			]
		}
	]}`),
};

export const FROST_P2_PRESET = {
	name: 'Frost P2 Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecMage>) => player.getRotation().type == RotationType.Frost,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 45497,
			"enchant": 3820,
			"gems": [
				41285,
				45883
			]
		},
		{
			"id": 45133,
			"gems": [
				40051
			]
		},
		{
			"id": 46134,
			"enchant": 3810,
			"gems": [
				39998
			]
		},
		{
			"id": 45618,
			"enchant": 3722,
			"gems": [
				40026
			]
		},
		{
			"id": 46130,
			"enchant": 3832,
			"gems": [
				39998,
				39998
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
			"id": 45665,
			"enchant": 3604,
			"gems": [
				39998,
				39998,
				0
			]
		},
		{
			"id": 45619,
			"gems": [
				40049,
				40049,
				39998
			]
		},
		{
			"id": 45488,
			"enchant": 3719,
			"gems": [
				39998,
				40051,
				40026
			]
		},
		{
			"id": 45135,
			"enchant": 3606,
			"gems": [
				39998,
				40049
			]
		},
		{
			"id": 46046,
			"gems": [
				39998
			]
		},
		{
			"id": 45495,
			"gems": [
				39998
			]
		},
		{
			"id": 45466
		},
		{
			"id": 45518
		},
		{
			"id": 45620,
			"enchant": 3834,
			"gems": [
				40026
			]
		},
		{
			"id": 45617
		},
		{
			"id": 45294,
			"gems": [
				39998
			]
		}
	]}`),
};
