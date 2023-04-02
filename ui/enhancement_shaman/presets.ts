import {
	Consumes,
	Flask,
	Food,
	Glyphs,
	EquipmentSpec,
	Potions,
	RaidBuffs,
	TristateEffect,
	Debuffs,
  CustomRotation,
  CustomSpell,
  ItemSwap,
  ItemSpec,
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
	EnhancementShaman_Rotation_PrimaryShock as PrimaryShock,
  EnhancementShaman_Rotation_RotationType as RotationType,
  EnhancementShaman_Rotation_CustomRotationSpell as CustomRotationSpell
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
		useFireElemental: true,
	}),
	maelstromweaponMinStack: 3,
	lightningboltWeave: true,
	autoWeaveDelay: 500,
	delayGcdWeave: 750,
	lavaburstWeave: false,
	firenovaManaThreshold: 3000,
	shamanisticRageManaThreshold: 25,
	primaryShock: PrimaryShock.Earth,
	weaveFlameShock: true,
  	rotationType: RotationType.Priority,
  	customRotation: CustomRotation.create({
			spells: [
			CustomSpell.create({ spell: CustomRotationSpell.LightningBolt }),
			CustomSpell.create({ spell: CustomRotationSpell.StormstrikeDebuffMissing }),
			CustomSpell.create({ spell: CustomRotationSpell.LightningBoltWeave }),
			CustomSpell.create({ spell: CustomRotationSpell.Stormstrike }),
			CustomSpell.create({ spell: CustomRotationSpell.FlameShock }),
			CustomSpell.create({ spell: CustomRotationSpell.EarthShock }),
			CustomSpell.create({ spell: CustomRotationSpell.MagmaTotem}),
			CustomSpell.create({ spell: CustomRotationSpell.LightningShield }),
			CustomSpell.create({ spell: CustomRotationSpell.FireNova }),
			CustomSpell.create({ spell: CustomRotationSpell.LightningBoltDelayedWeave }),
			CustomSpell.create({ spell: CustomRotationSpell.LavaLash }),
		],
	}),
});

export const DefaultOptions = EnhancementShamanOptions.create({
	shield: ShamanShield.LightningShield,
	bloodlust: true,
	imbueMh: ShamanImbue.WindfuryWeapon,
	imbueOh: ShamanImbue.FlametongueWeapon,
	syncType: ShamanSyncType.Auto,
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
	demonicPact: 500,
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
			"enchant": 3817,
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
			"enchant": 3808
		},
		{
			"id": 37840,
			"enchant": 3605
		},
		{
			"id": 39597,
			"enchant": 3832,
			"gems": [
				40053,
				40088
			]
		},
		{
			"id": 43131,
			"enchant": 3845,
			"gems": [
				0
			]
		},
		{
			"id": 39601,
			"enchant": 3604,
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
			"enchant": 3823
		},
		{
			"id": 37167,
			"enchant": 3606,
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
			"enchant": 3789
		},
		{
			"id": 40704,
			"enchant": 3789
		},
		{
			"id": 33507
		}
	]}`),
}

export const P1_PRESET = {
	name: 'P1 Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 40543,
			"enchant": 3817,
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
			"enchant": 3808,
			"gems": [
				40014
			]
		},
		{
			"id": 40403,
			"enchant": 3605
		},
		{
			"id": 40523,
			"enchant": 3832,
			"gems": [
				40003,
				40014
			]
		},
		{
			"id": 40282,
			"enchant": 3845,
			"gems": [
				42702,
				0
			]
		},
		{
			"id": 40520,
			"enchant": 3604,
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
			"enchant": 3823,
			"gems": [
				39999,
				42156
			]
		},
		{
			"id": 40367,
			"enchant": 3606,
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
			"enchant": 3789
		},
		{
			"id": 39468,
			"enchant": 3789
		},
		{
			"id": 40322
		}
	]}`),
};

export const P2_PRESET_FT = {
	name: 'P2 Preset FT',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	gear: EquipmentSpec.fromJsonString(`{ "items": [
        {
          "id": 45610,
          "enchant": 3817,
          "gems": [
            41398,
            42702
          ]
        },
        {
          "id": 45517,
          "gems": [
            39999
          ]
        },
        {
          "id": 46203,
          "enchant": 3808,
          "gems": [
            39999
          ]
        },
        {
          "id": 45461,
          "enchant": 3831,
          "gems": [
            40014
          ]
        },
        {
          "id": 46205,
          "enchant": 3832,
          "gems": [
            40058,
            40053
          ]
        },
        {
          "id": 45460,
          "enchant": 3845,
          "gems": [
            39999,
            0
          ]
        },
        {
          "id": 46200,
          "enchant": 3604,
          "gems": [
            40014,
            0
          ]
        },
        {
          "id": 45553,
          "gems": [
            36766,
            36766,
            36766
          ]
        },
        {
          "id": 46208,
          "enchant": 3823,
          "gems": [
            39999,
            39999
          ]
        },
        {
          "id": 45989,
          "enchant": 3606,
          "gems": [
            40053,
            39999
          ]
        },
        {
          "id": 45456,
          "gems": [
            39999
          ]
        },
        {
          "id": 46046,
          "gems": [
            40053
          ]
        },
        {
          "id": 45609
        },
        {
          "id": 46038
        },
        {
          "id": 45612,
          "enchant": 3789,
          "gems": [
            39999
          ]
        },
        {
          "id": 46097,
          "enchant": 3789,
          "gems": [
            40003
          ]
        },
        {
          "id": 40322
        }
      ]
    }`),
};

export const P2_PRESET_WF = {
	name: 'P2 Preset WF',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	gear: EquipmentSpec.fromJsonString(`{  "items": [
        {
          "id": 45610,
          "enchant": 3817,
          "gems": [
            41398,
            42702
          ]
        },
        {
          "id": 45517,
          "gems": [
            39999
          ]
        },
        {
          "id": 46203,
          "enchant": 3808,
          "gems": [
            39999
          ]
        },
        {
          "id": 45461,
          "enchant": 3831,
          "gems": [
            40052
          ]
        },
        {
          "id": 46205,
          "enchant": 3832,
          "gems": [
            40052,
            40052
          ]
        },
        {
          "id": 45460,
          "enchant": 3845,
          "gems": [
            39999,
            0
          ]
        },
        {
          "id": 46200,
          "enchant": 3604,
          "gems": [
            40053,
            0
          ]
        },
        {
          "id": 45553,
          "gems": [
            36766,
            36766,
            36766
          ]
        },
        {
          "id": 46208,
          "enchant": 3823,
          "gems": [
            39999,
            39999
          ]
        },
        {
          "id": 45989,
          "enchant": 3606,
          "gems": [
            40053,
            39999
          ]
        },
        {
          "id": 45456,
          "gems": [
            39999
          ]
        },
        {
          "id": 45608,
          "gems": [
            39999
          ]
        },
        {
          "id": 45609
        },
        {
          "id": 46038
        },
        {
          "id": 45132,
          "enchant": 3789,
          "gems": [
            40052
          ]
        },
        {
          "id": 46097,
          "enchant": 3789,
          "gems": [
            39999
          ]
        },
        {
          "id": 40322
        }
      ]
    }`),
};

