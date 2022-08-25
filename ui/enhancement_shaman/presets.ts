import { Consumes,
	Flask,
	Food,
	Glyphs,
	EquipmentSpec,
	Potions,
	RaidBuffs,
	TristateEffect,
  Debuffs,
} from '../core/proto/common.js';
import { SavedTalents } from '../core/proto/ui.js';

import { EnhancementShaman_Rotation as EnhancementShamanRotation, EnhancementShaman_Options as EnhancementShamanOptions, ShamanShield } from '../core/proto/shaman.js';
import {
	AirTotem,
	EarthTotem,
	FireTotem,
	WaterTotem,
	ShamanTotems,
    ShamanImbue,
    ShamanSyncType,
    ShamanMajorGlyph,
} from '../core/proto/shaman.js';

import * as Tooltips from '../core/constants/tooltips.js';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/wotlk/talent-calc and copy the numbers in the url.
export const StandardTalents = {
	name: 'Standard',
	data: SavedTalents.create({
		talentsString: '053030152-30405003105021333031131031051',
        glyphs: Glyphs.create({
            major1: ShamanMajorGlyph.GlyphOfStormstrike,
            major2: ShamanMajorGlyph.GlyphOfFlametongueWeapon,
            major3: ShamanMajorGlyph.GlyphOfFeralSpirit,
            //minor glyphs dont affect damage done, all convenience/QoL
        })
	}),
};

export const DefaultRotation = EnhancementShamanRotation.create({
	totems: ShamanTotems.create({
		earth: EarthTotem.StrengthOfEarthTotem,
		air: AirTotem.WindfuryTotem,
		fire: FireTotem.MagmaTotem,
		water: WaterTotem.ManaSpringTotem,
	}),
  maelstromweaponMinStack: 3,
  lightningboltWeave: true,
  lavaburstWeave: false,
  firenovaManaThreshold: 3000,
  shamanisticRageManaThreshold: 25,
});

export const DefaultOptions = EnhancementShamanOptions.create({
	shield: ShamanShield.LightningShield,
	bloodlust: true,
	imbueMh: ShamanImbue.WindfuryWeapon,
	imbueOh: ShamanImbue.FlametongueWeapon,
	syncType: ShamanSyncType.SyncMainhandOffhandSwings,
});

export const DefaultConsumes = Consumes.create({
	defaultPotion: Potions.PotionOfSpeed,
	flask: Flask.FlaskOfEndlessRage,
	food: Food.FoodFishFeast,
});

export const DefaultRaidBuffs = RaidBuffs.create({
	giftOfTheWild: TristateEffect.TristateEffectImproved,
	arcaneBrilliance: true,
	leaderOfThePack: TristateEffect.TristateEffectImproved,
	totemOfWrath: true,
  wrathOfAirTotem: true,
	moonkinAura: TristateEffect.TristateEffectImproved,
	sanctifiedRetribution: true,
	divineSpirit: true,
	battleShout: TristateEffect.TristateEffectImproved,
});

export const DefaultDebuffs = Debuffs.create({
  bloodFrenzy: true,
  sunderArmor: true,
  curseOfWeakness: TristateEffect.TristateEffectRegular,
  curseOfElements: true,
  faerieFire: TristateEffect.TristateEffectImproved,
  judgementOfWisdom: true,
  misery: true,
  totemOfWrath: true,
  shadowMastery: true,
});


export const PreRaid_PRESET = {
    name: 'Preraid Preset',
    tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
    gear: EquipmentSpec.fromJsonString(`{"items": [
        {
          "id": 43311,
          "enchant": 44879,
          "gems": [
            41398,
            42156
          ]
        },
        {
          "id": 40678
        },
        {
          "id": 37373,
          "enchant": 44871
        },
        {
          "id": 37840,
          "enchant": 55002
        },
        {
          "id": 39597,
          "enchant": 44489,
          "gems": [
            40053,
            40088
          ]
        },
        {
          "id": 43131,
          "enchant": 44484,
          "gems": [
            0
          ]
        },
        {
          "id": 39601,
          "enchant": 54999,
          "gems": [
            40053,
            0
          ]
        },
        {
          "id": 37407,
          "gems": [
            42156
          ]
        },
        {
          "id": 37669,
          "enchant": 38374
        },
        {
          "id": 37167,
          "enchant": 55016,
          "gems": [
            40053,
            42156
          ]
        },
        {
          "id": 37685
        },
        {
          "id": 37642
        },
        {
          "id": 37390
        },
        {
          "id": 40684
        },
        {
          "id": 41384,
          "enchant": 44492
        },
        {
          "id": 40704,
          "enchant": 44492
        },
        {
          "id": 37575
        }
    ]}`)
}

export const P1_PRESET = {
	name: 'P1 Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	gear: EquipmentSpec.fromJsonString(`{"items": [
        {
          "id": 40543,
          "enchant": 44879,
          "gems": [
            41398,
            40014
          ]
        },
        {
          "id": 44661,
          "gems": [
            40014
          ]
        },
        {
          "id": 40524,
          "enchant": 44871,
          "gems": [
            40014
          ]
        },
        {
          "id": 40403,
          "enchant": 55002
        },
        {
          "id": 40523,
          "enchant": 44489,
          "gems": [
            40003,
            40014
          ]
        },
        {
          "id": 40282,
          "enchant": 44484,
          "gems": [
            42702,
            0
          ]
        },
        {
          "id": 40520,
          "enchant": 54999,
          "gems": [
            42154,
            0
          ]
        },
        {
          "id": 40275,
          "gems": [
            42156
          ]
        },
        {
          "id": 40522,
          "enchant": 38374,
          "gems": [
            39999,
            42156
          ]
        },
        {
          "id": 40367,
          "enchant": 55016,
          "gems": [
            40058
          ]
        },
        {
          "id": 40474
        },
        {
          "id": 40074
        },
        {
          "id": 40684
        },
        {
          "id": 37390
        },
        {
          "id": 39763,
          "enchant": 44492
        },
        {
          "id": 39468,
          "enchant": 44492
        },
        {
          "id": 40322
        }
      ]}`),
};
