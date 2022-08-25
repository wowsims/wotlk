import { Consumes,
	Flask,
	Food,
  PetFood,
	Glyphs,
	EquipmentSpec,
	Potions,
	RaidBuffs,
	IndividualBuffs,
	Debuffs,
	TristateEffect,
  Spec,
} from '../core/proto/common.js';
import { SavedTalents } from '../core/proto/ui.js';
import { Player } from '../core/player.js';

import {
	Warlock_Rotation as WarlockRotation,
	Warlock_Options as WarlockOptions,
	Warlock_Rotation_PrimarySpell as PrimarySpell,
	Warlock_Rotation_SecondaryDot as SecondaryDot,
	Warlock_Rotation_SpecSpell as SpecSpell,
	Warlock_Rotation_Curse as Curse,
  Warlock_Rotation_Type as RotationType,
	Warlock_Options_WeaponImbue as WeaponImbue,
	Warlock_Options_Armor as Armor,
	Warlock_Options_Summon as Summon,
	WarlockMajorGlyph as MajorGlyph,
	WarlockMinorGlyph as MinorGlyph,
} from '../core/proto/warlock.js';

import * as WarlockTooltips from './tooltips.js';

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/wotlk/talent-calc and copy the numbers in the url.

export const AfflictionTalents = {
	name: 'Affliction',
	data: SavedTalents.create({
		talentsString: '2350002030023510253500331151--550000051',
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
			major3: MajorGlyph.GlyphOfFelguard,
			minor1: MinorGlyph.GlyphOfSouls,
			minor2: MinorGlyph.GlyphOfDrainSoul,
			minor3: MinorGlyph.GlyphOfEnslaveDemon,
		}),
	}),
};

export const DestructionTalents = {
	name: 'Destruction',
	data: SavedTalents.create({
		talentsString: '-03310030003-05203205210331051335230351',
		glyphs: Glyphs.create({
			major1: MajorGlyph.GlyphOfConflagrate,
			major2: MajorGlyph.GlyphOfLifeTap,
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
	corruption: false,
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
	petFood: PetFood.PetFoodSpicedMammothTreats,
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

export const DestroIndividualBuffs = IndividualBuffs.create({
  blessingOfKings: true,
  blessingOfWisdom: TristateEffect.TristateEffectImproved,
  blessingOfMight: TristateEffect.TristateEffectImproved,
});

export const DefaultDebuffs = Debuffs.create({
	sunderArmor: true,
	faerieFire: TristateEffect.TristateEffectImproved,
	bloodFrenzy: true,
	ebonPlaguebringer: true,
	heartOfTheCrusader: true,
	judgementOfWisdom: true,
});

export const DestroDebuffs = Debuffs.create({
  sunderArmor: true,
  faerieFire: TristateEffect.TristateEffectImproved,
  bloodFrenzy: true,
  ebonPlaguebringer: true,
  heartOfTheCrusader: true,
  judgementOfWisdom: true,
  shadowMastery: true,
});

export const OtherDefaults = {
	distanceFromTarget: 25,
};

export const SWP_BIS = {
	name: 'Straight Outa SWP',
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
          "id": 34364,
          "enchant": 24003,
          "gems": [
            32196,
            35488,
            32196
          ]
        },
        {
          "id": 34436,
          "enchant": 22534,
          "gems": [
            35760,
            0
          ]
        },
        {
          "id": 34344,
          "enchant": 28272,
          "gems": [
            35760,
            32196,
            0
          ]
        },
        {
          "id": 34541,
          "gems": [
            35760,
            0
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
export const P1_PreBiS_11 = {
	name: 'Pre-Raid Affliction',
	tooltip: WarlockTooltips.BIS_TOOLTIP,
  enableWhen: (player: Player<Spec.SpecWarlock>) => player.getRotation().type == RotationType.Affliction,
	gear: EquipmentSpec.fromJsonString(`
    {"items":
      [
        {
          "id": 44910,
          "enchant": 44877,
          "gems": [
            41285,
            39998
          ]
        },
        {
          "id": 42647,
          "gems": [
            39998
          ]
        },
        {
          "id": 34210,
          "enchant": 44874,
          "gems": [
            39998,
            40051
          ]
        },
        {
          "id": 41610,
          "enchant": 55642
        },
        {
          "id": 39497,
          "enchant": 44489,
          "gems": [
            39998,
            40051
          ]
        },
        {
          "id": 37361,
          "enchant": 44498,
          "gems": [
            0
          ]
        },
        {
          "id": 42113,
          "enchant": 54999,
          "gems": [
            0
          ]
        },
        {
          "id": 40696,
          "gems": [
            40051,
            39998
          ]
        },
        {
          "id": 34181,
          "enchant": 41602,
          "gems": [
            39998,
            39998,
            40051
          ]
        },
        {
          "id": 44202,
          "enchant": 55016,
          "gems": [
            40026
          ]
        },
        {
          "id": 43253,
          "gems": [
            40026
          ]
        },
        {
          "id": 37694
        },
        {
          "id": 40682
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
          "id": 34348,
          "gems": [
            39998
          ]
        }
      ]
    }
  `),
}
export const P1_Preset_Demo_Destro = {
  name: 'P1 Preset Demo / Destro',
  tooltip: WarlockTooltips.BIS_TOOLTIP,
  enableWhen: (player: Player<Spec.SpecWarlock>) => player.getRotation().type == RotationType.Demonology || player.getRotation().type == RotationType.Destruction,
  gear: EquipmentSpec.fromJsonString(`
    {"items":
      [
        {
          "id": 40421,
          "enchant": 44877,
          "gems": [
            41285,
            40014
          ]
        },
        {
          "id": 44661,
          "gems": [
            40099
          ]
        },
        {
          "id": 40424,
          "enchant": 44874,
          "gems": [
            40049
          ]
        },
        {
          "id": 44005,
          "enchant": 55642,
          "gems": [
            40099
          ]
        },
        {
          "id": 40423,
          "enchant": 44489,
          "gems": [
            40049,
            40014
          ]
        },
        {
          "id": 44008,
          "enchant": 44498,
          "gems": [
            39998,
            0
          ]
        },
        {
          "id": 40420,
          "enchant": 54999,
          "gems": [
            39998,
            0
          ]
        },
        {
          "id": 40561,
          "gems": [
            40014
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
          "id": 40399
        },
        {
          "id": 40719
        },
        {
          "id": 40432
        },
        {
          "id": 40255
        },
        {
          "id": 40396,
          "enchant": 44487
        },
        {
          "id": 39766
        },
        {
          "id": 39712
        }
      ]
    }
  `),
}


// will have only rare gems, but a Lightweave Embroidery on cloak.
export const P1_Preset_Affliction = {
	name: 'P1 Affliction Preset',
	tooltip: WarlockTooltips.BIS_TOOLTIP,
  enableWhen: (player: Player<Spec.SpecWarlock>) => player.getRotation().type == RotationType.Affliction,
	gear: EquipmentSpec.fromJsonString(`
    {"items":
      [
        {
          "id": 40421,
          "enchant": 44877,
          "gems": [
            41285,
            40051
          ]
        },
        {
          "id": 44661,
          "gems": [
            40026
          ]
        },
        {
          "id": 40424,
          "enchant": 44874,
          "gems": [
            39998
          ]
        },
        {
          "id": 44005,
          "enchant": 55642,
          "gems": [
            40026
          ]
        },
        {
          "id": 40423,
          "enchant": 44489,
          "gems": [
            39998,
            40051
          ]
        },
        {
          "id": 44008,
          "enchant": 44498,
          "gems": [
            39998,
            0
          ]
        },
        {
          "id": 40420,
          "enchant": 54999,
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
          "id": 40560,
          "enchant": 41602
        },
        {
          "id": 40558,
          "enchant": 55016
        },
        {
          "id": 40399
        },
        {
          "id": 40719
        },
        {
          "id": 40432
        },
        {
          "id": 40255
        },
        {
          "id": 40396,
          "enchant": 44487
        },
        {
          "id": 39766
        },
        {
          "id": 39712
        }
      ]
    }
  `),
}


export const P1_PreBiS_14 = {
  name: 'Pre-Raid Preset',
  tooltip: WarlockTooltips.BIS_TOOLTIP,
  enableWhen: (player: Player<Spec.SpecWarlock>) => player.getRotation().type == RotationType.Demonology || player.getRotation().type == RotationType.Destruction,
  gear: EquipmentSpec.fromJsonString(`
    {"items":
      [
        {
          "id": 44910,
          "enchant": 44877,
          "gems": [
            41285,
            39998
          ]
        },
        {
          "id": 42647,
          "gems": [
            40049
          ]
        },
        {
          "id": 34210,
          "enchant": 44874,
          "gems": [
            39998,
            40014
          ]
        },
        {
          "id": 41610,
          "enchant": 55642
        },
        {
          "id": 39497,
          "enchant": 44489,
          "gems": [
            39998,
            40014
          ]
        },
        {
          "id": 37361,
          "enchant": 44498,
          "gems": [
            0
          ]
        },
        {
          "id": 42113,
          "enchant": 54999,
          "gems": [
            0
          ]
        },
        {
          "id": 40696,
          "gems": [
            40014,
            39998
          ]
        },
        {
          "id": 34181,
          "enchant": 41602,
          "gems": [
            39998,
            39998,
            40014
          ]
        },
        {
          "id": 44202,
          "enchant": 55016,
          "gems": [
            40026
          ]
        },
        {
          "id": 43253,
          "gems": [
            40026
          ]
        },
        {
          "id": 37694
        },
        {
          "id": 40682
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
          "id": 34348,
          "gems": [
            39998
          ]
        }
      ]
    }
  `),
}

export const Naked = {
  name: 'The Naked Bolt',
  gear: EquipmentSpec.fromJsonString(`
    {"items":
      [
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {}
      ]
    }
  `),
}
  
