import { Consumes,
	Flask,
	Food,
	Glyphs,
	EquipmentSpec,
	ItemSpec,
	Potions,
	Faction,
	RaidBuffs,
	PartyBuffs,
	IndividualBuffs,
	Debuffs,
	Spec,
	Stat,
	TristateEffect,
	Race,
} from '/wotlk/core/proto/common.js';
import { SavedTalents } from '/wotlk/core/proto/ui.js';
import { Player } from '/wotlk/core/player.js';

import {
	Warlock,
	Warlock_Rotation as WarlockRotation,
	WarlockTalents as WarlockTalents,
	Warlock_Options as WarlockOptions,
	Warlock_Rotation_PrimarySpell as PrimarySpell,
	Warlock_Rotation_SecondaryDot as SecondaryDot,
	Warlock_Rotation_SpecSpell as SpecSpell,
	Warlock_Rotation_Curse as Curse,
	Warlock_Options_WeaponImbue as WeaponImbue,
	Warlock_Options_Armor as Armor,
	Warlock_Options_Summon as Summon,
	WarlockMajorGlyph as MajorGlyph,
	WarlockMinorGlyph as MinorGlyph,
} from '/wotlk/core/proto/warlock.js';

import * as Enchants from '/wotlk/core/constants/enchants.js';
import * as Gems from '/wotlk/core/proto_utils/gems.js';
import * as Tooltips from '/wotlk/core/constants/tooltips.js';

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/wotlk/talent-calc and copy the numbers in the url.

export const AfflictionTalents = {
	name: 'Affliction',
	data: SavedTalents.create({
		talentsString: '2350002030023510253510331151--55000005',
		glyphs: Glyphs.create({
			major1: MajorGlyph.GlyphOfQuickDecay,
			major2: MajorGlyph.GlyphOfLifeTap,
			major3: MajorGlyph.GlyphOfHaunt,
			minor1: MinorGlyph.GlyphOfSouls,
			minor2: MinorGlyph.GlyphOfDrainSoul,
			minor3: MinorGlyph.GlyphOfEnslaveDemon,
		}),
	}),
};

export const DemonologyTalents = {
	name: 'Demonology',
	data: SavedTalents.create({
		talentsString: '-203203301035012530135201351-550000052',
		glyphs: Glyphs.create({
			major1: MajorGlyph.GlyphOfLifeTap,
			major2: MajorGlyph.GlyphOfQuickDecay,
			major3: MajorGlyph.GlyphOfMetamorphosis,
			minor1: MinorGlyph.GlyphOfSouls,
			minor2: MinorGlyph.GlyphOfDrainSoul,
			minor3: MinorGlyph.GlyphOfEnslaveDemon,
		}),
	}),
};

export const DestructionTalents = {
	name: 'Destruction',
	data: SavedTalents.create({
		talentsString: '030-03310030003-05203205220331051035031351',
		glyphs: Glyphs.create({
			major1: MajorGlyph.GlyphOfConflagrate,
			major2: MajorGlyph.GlyphOfImp,
			major3: MajorGlyph.GlyphOfIncinerate,
			minor1: MinorGlyph.GlyphOfSouls,
			minor2: MinorGlyph.GlyphOfDrainSoul,
			minor3: MinorGlyph.GlyphOfEnslaveDemon,
		}),
	}),
};

export const AfflictionRotation = WarlockRotation.create({
		primarySpell: PrimarySpell.ShadowBolt,
		secondaryDot: SecondaryDot.UnstableAffliction,
		specSpell: SpecSpell.Haunt,
		curse: Curse.Agony,
		corruption: true,
		detonateSeed: true,
});

export const DemonologyRotation = WarlockRotation.create({
	primarySpell: PrimarySpell.ShadowBolt,
	secondaryDot: SecondaryDot.Immolate,
	specSpell: SpecSpell.NoSpecSpell,
	curse: Curse.Doom,
	corruption: true,
	detonateSeed: true,
});

export const DestructionRotation = WarlockRotation.create({
	primarySpell: PrimarySpell.Incinerate,
	secondaryDot: SecondaryDot.Immolate,
	specSpell: SpecSpell.ChaosBolt,
	curse: Curse.Doom,
	corruption: true,
	detonateSeed: true,
});

export const AfflictionOptions = WarlockOptions.create({
	armor: Armor.FelArmor,
	summon: Summon.Felhunter,
	weaponImbue: WeaponImbue.GrandSpellstone,
});

export const DemonologyOptions = WarlockOptions.create({
	armor: Armor.FelArmor,
	summon: Summon.Felguard,
	weaponImbue: WeaponImbue.GrandSpellstone,
});

export const DestructionOptions = WarlockOptions.create({
	armor: Armor.FelArmor,
	summon: Summon.Imp,
	weaponImbue: WeaponImbue.GrandFirestone,
});

export const DefaultConsumes = Consumes.create({
	flask: Flask.FlaskOfTheFrostWyrm,
	food: Food.FoodFishFeast,
	defaultPotion: Potions.PotionOfWildMagic,
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

export const SWP_BIS = {
	name: 'SWP BiS',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	gear: EquipmentSpec.fromJsonString(`
{"items": [
        {
          "id": 34340,
          "enchant": 29191,
          "gems": [
            34220,
            32215
          ]
        },
        {
          "id": 34204
        },
        {
          "id": 31054,
          "enchant": 28886,
          "gems": [
            32215,
            35760
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
          "enchant": 24003,
          "gems": [
            32196,
            35760,
            35760
          ]
        },
        {
          "id": 34436,
          "enchant": 22534,
          "gems": [
            35760
          ]
        },
        {
          "id": 34344,
          "enchant": 28272,
          "gems": [
            35760,
            32196
          ]
        },
        {
          "id": 34541,
          "gems": [
            35760
          ]
        },
        {
          "id": 34181,
          "enchant": 24274,
          "gems": [
            32196,
            32196,
            35760
          ]
        },
        {
          "id": 34564,
          "enchant": 35297,
          "gems": [
            35760
          ]
        },
        {
          "id": 34362,
          "enchant": 22536
        },
        {
          "id": 34230,
          "enchant": 22536
        },
        {
          "id": 32483
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
            35760
          ]
        }
      ]}
    `),
};
export const P1_PreBiS = {
	name: 'Pre-BiS',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	gear: EquipmentSpec.fromJsonString(`
{"items": [
        {
          "id": 44910,
          "enchant": 29191,
          "gems": [
            41285,
            40113
          ]
        },
        {
          "id": 42647,
          "gems": [
            40113
          ]
        },
        {
          "id": 34210,
          "enchant": 28886,
          "gems": [
            40113,
            40155
          ]
        },
        {
          "id": 41610,
          "enchant": 55642
        },
        {
          "id": 39497,
          "enchant": 24003,
          "gems": [
            40113,
            40155
          ]
        },
        {
          "id": 37361,
          "enchant": 22534,
          "gems": [
            0
          ]
        },
        {
          "id": 34344,
          "enchant": 54999,
          "gems": [
            40155,
            40113,
            0
          ]
        },
        {
          "id": 40696,
          "gems": [
            40155,
            0
          ]
        },
        {
          "id": 34181,
          "enchant": 24274,
          "gems": [
            40113,
            40113,
            40155
          ]
        },
        {
          "id": 44202,
          "enchant": 35297,
          "gems": [
            40133
          ]
        },
        {
          "id": 43253,
          "gems": [
            40133
          ]
        },
        {
          "id": 40585
        },
        {
          "id": 40682
        },
        {
          "id": 37873
        },
        {
          "id": 37360,
          "enchant": 44487
        },
        {},
        {
          "id": 34347,
          "gems": [
            40155
          ]
        }
      ]
    }
    `),
}
