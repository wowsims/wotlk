import {
	Consumes,
	CustomRotation,
	CustomSpell,
	EquipmentSpec,
	Explosive,
	Flask,
	Food,
	Glyphs,
	PetFood,
	Potions,
	RaidTarget,
	Spec
} from '../core/proto/common.js';
import {SavedTalents} from '../core/proto/ui.js';
import {Player} from '../core/player.js';
import {NO_TARGET} from '../core/proto_utils/utils.js';

import {
	Deathknight_Options as DeathKnightOptions,
	Deathknight_Rotation as DeathKnightRotation,
	Deathknight_Rotation_ArmyOfTheDead,
	Deathknight_Rotation_BloodRuneFiller,
	Deathknight_Rotation_CustomSpellOption as CustomSpellOption,
	Deathknight_Rotation_FrostRotationType,
	Deathknight_Rotation_Presence,
	DeathknightMajorGlyph,
	DeathknightMinorGlyph,
	Deathknight_Rotation_DrwDiseases,
	Deathknight_Rotation_BloodSpell,
} from '../core/proto/deathknight.js';

import * as Tooltips from '../core/constants/tooltips.js';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wotlk.wowhead.com/talent-calc and copy the numbers in the url.
export const FrostTalents = {
	name: 'Frost BL',
	data: SavedTalents.create({
		talentsString: '23050005-32005350352203012300033101351',
		glyphs: Glyphs.create({
			major1: DeathknightMajorGlyph.GlyphOfObliterate,
			major2: DeathknightMajorGlyph.GlyphOfFrostStrike,
			major3: DeathknightMajorGlyph.GlyphOfDisease,
			minor1: DeathknightMinorGlyph.GlyphOfHornOfWinter,
			minor2: DeathknightMinorGlyph.GlyphOfPestilence,
			minor3: DeathknightMinorGlyph.GlyphOfRaiseDead,
		}),
	}),
};

export const FrostUnholyTalents = {
	name: 'Frost UH',
	data: SavedTalents.create({
		talentsString: '01-32002350342203012300033101351-230200305003',
		glyphs: Glyphs.create({
			major1: DeathknightMajorGlyph.GlyphOfObliterate,
			major2: DeathknightMajorGlyph.GlyphOfFrostStrike,
			major3: DeathknightMajorGlyph.GlyphOfDisease,
			minor1: DeathknightMinorGlyph.GlyphOfHornOfWinter,
			minor2: DeathknightMinorGlyph.GlyphOfPestilence,
			minor3: DeathknightMinorGlyph.GlyphOfRaiseDead,
		}),
	}),
};

export const UnholyDualWieldTalents = {
	name: 'Unholy DW',
	data: SavedTalents.create({
		talentsString: '-320043500002-2300303050032152000150013133051',
		glyphs: Glyphs.create({
			major1: DeathknightMajorGlyph.GlyphOfTheGhoul,
			major2: DeathknightMajorGlyph.GlyphOfIcyTouch,
			major3: DeathknightMajorGlyph.GlyphOfDeathAndDecay,
			minor1: DeathknightMinorGlyph.GlyphOfHornOfWinter,
			minor2: DeathknightMinorGlyph.GlyphOfPestilence,
			minor3: DeathknightMinorGlyph.GlyphOfRaiseDead,
		}),
	}),
};

export const UnholyDualWieldSSTalents = {
	name: 'Unholy DW SS',
	data: SavedTalents.create({
		talentsString: '-320033500002-2301303050032151000150013133151',
		glyphs: Glyphs.create({
			major1: DeathknightMajorGlyph.GlyphOfTheGhoul,
			major2: DeathknightMajorGlyph.GlyphOfIcyTouch,
			major3: DeathknightMajorGlyph.GlyphOfDeathAndDecay,
			minor1: DeathknightMinorGlyph.GlyphOfHornOfWinter,
			minor2: DeathknightMinorGlyph.GlyphOfPestilence,
			minor3: DeathknightMinorGlyph.GlyphOfRaiseDead,
		}),
	}),
};

export const Unholy2HTalents = {
	name: 'Unholy 2H',
	data: SavedTalents.create({
		talentsString: '-320050500002-2300303150032152000150013133151',
		glyphs: Glyphs.create({
			major1: DeathknightMajorGlyph.GlyphOfTheGhoul,
			major2: DeathknightMajorGlyph.GlyphOfDarkDeath,
			major3: DeathknightMajorGlyph.GlyphOfIcyTouch,
			minor1: DeathknightMinorGlyph.GlyphOfHornOfWinter,
			minor2: DeathknightMinorGlyph.GlyphOfPestilence,
			minor3: DeathknightMinorGlyph.GlyphOfRaiseDead,
		}),
	}),
};

export const BloodTalents = {
	name: 'Blood DPS',
	data: SavedTalents.create({
		talentsString: '2305120530003303231023001351--2302003050032',
		glyphs: Glyphs.create({
			major1: DeathknightMajorGlyph.GlyphOfDancingRuneWeapon,
			major2: DeathknightMajorGlyph.GlyphOfDeathStrike,
			major3: DeathknightMajorGlyph.GlyphOfDisease,
			minor1: DeathknightMinorGlyph.GlyphOfHornOfWinter,
			minor2: DeathknightMinorGlyph.GlyphOfPestilence,
			minor3: DeathknightMinorGlyph.GlyphOfRaiseDead,
		}),
	}),
};

export const DefaultUnholyRotation = DeathKnightRotation.create({
	useDeathAndDecay: true,
	btGhoulFrenzy: false,
	refreshHornOfWinter: false,
	useGargoyle: true,
	useEmpowerRuneWeapon: true,
	holdErwArmy: false,
	preNerfedGargoyle: false,
	armyOfTheDead: Deathknight_Rotation_ArmyOfTheDead.AsMajorCd,
	startingPresence: Deathknight_Rotation_Presence.Unholy,
	blPresence: Deathknight_Rotation_Presence.Blood,
	presence: Deathknight_Rotation_Presence.Blood,
	gargoylePresence: Deathknight_Rotation_Presence.Unholy,
	bloodRuneFiller: Deathknight_Rotation_BloodRuneFiller.BloodBoil,
	useAms: false,
	drwDiseases: Deathknight_Rotation_DrwDiseases.Pestilence,
	bloodSpender: Deathknight_Rotation_BloodSpell.HS,
	useDancingRuneWeapon: true
});

export const DefaultUnholyOptions = DeathKnightOptions.create({
	drwPestiApply: true,
	startingRunicPower: 0,
	petUptime: 1,
	precastGhoulFrenzy: false,
	precastHornOfWinter: true,
	unholyFrenzyTarget: RaidTarget.create({
		targetIndex: NO_TARGET, // In an individual sim the 0-indexed player is ourself.
	}),
	diseaseDowntime: 2,
});

export const DefaultFrostRotation = DeathKnightRotation.create({
	useDeathAndDecay: false,
	btGhoulFrenzy: false,
	refreshHornOfWinter: false,
	useEmpowerRuneWeapon: true,
	preNerfedGargoyle: false,
	startingPresence: Deathknight_Rotation_Presence.Blood,
	presence: Deathknight_Rotation_Presence.Blood,
	bloodRuneFiller: Deathknight_Rotation_BloodRuneFiller.BloodBoil,
	useAms: false,
	avgAmsSuccessRate: 1.0,
	avgAmsHit: 10000.0,
	drwDiseases: Deathknight_Rotation_DrwDiseases.Pestilence,
  	frostRotationType: Deathknight_Rotation_FrostRotationType.SingleTarget,
  	frostCustomRotation: CustomRotation.create({
		spells: [
			CustomSpell.create({ spell: CustomSpellOption.CustomDeathAndDecay }),
			CustomSpell.create({ spell: CustomSpellOption.CustomIcyTouch }),
			CustomSpell.create({ spell: CustomSpellOption.CustomPlagueStrike }),
			CustomSpell.create({ spell: CustomSpellOption.CustomPestilence }),
			CustomSpell.create({ spell: CustomSpellOption.CustomHowlingBlastRime }),
			CustomSpell.create({ spell: CustomSpellOption.CustomHowlingBlast }),
			CustomSpell.create({ spell: CustomSpellOption.CustomBloodBoil }),
			CustomSpell.create({ spell: CustomSpellOption.CustomObliterate }),
			CustomSpell.create({ spell: CustomSpellOption.CustomFrostStrike }),
		],
	}),
});

export const DefaultFrostOptions = DeathKnightOptions.create({
	drwPestiApply: true,
	startingRunicPower: 0,
	petUptime: 1,
	precastHornOfWinter: true,
	unholyFrenzyTarget: RaidTarget.create({
		targetIndex: NO_TARGET, // In an individual sim the 0-indexed player is ourself.
	}),
	diseaseDowntime: 0,
});

export const DefaultBloodRotation = DeathKnightRotation.create({
	refreshHornOfWinter: false,
	useEmpowerRuneWeapon: true,
	preNerfedGargoyle: false,
	startingPresence: Deathknight_Rotation_Presence.Blood,
	bloodRuneFiller: Deathknight_Rotation_BloodRuneFiller.BloodStrike,
	armyOfTheDead: Deathknight_Rotation_ArmyOfTheDead.PreCast,
	holdErwArmy: false,
	useAms: false,
	drwDiseases: Deathknight_Rotation_DrwDiseases.Pestilence,
	bloodSpender: Deathknight_Rotation_BloodSpell.HS,
	useDancingRuneWeapon: true,
});

export const DefaultBloodOptions = DeathKnightOptions.create({
	drwPestiApply: true,
	startingRunicPower: 0,
	petUptime: 1,
	precastHornOfWinter: true,
	unholyFrenzyTarget: RaidTarget.create({
		targetIndex: NO_TARGET, // In an individual sim the 0-indexed player is ourself.
	}),
	diseaseDowntime: 0,
});

export const OtherDefaults = {
};

export const DefaultConsumes = Consumes.create({
	flask: Flask.FlaskOfEndlessRage,
	food: Food.FoodDragonfinFilet,
	defaultPotion: Potions.PotionOfSpeed,
	petFood: PetFood.PetFoodSpicedMammothTreats,
	prepopPotion: Potions.PotionOfSpeed,
	thermalSapper: true,
	fillerExplosive: Explosive.ExplosiveSaroniteBomb,
});

export const P1_BLOOD_BIS_PRESET = {
	name: 'P1 Blood',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecDeathknight>) => player.getTalentTree() == 0,
	gear: EquipmentSpec.fromJsonString(`{   "items": [
		{
			"id": 44006,
			"enchant": 3817,
			"gems": [
			  41398,
			  42702
			]
		  },
		  {
			"id": 44664,
			"gems": [
			  39996
			]
		  },
		  {
			"id": 40557,
			"enchant": 3808,
			"gems": [
			  39996
			]
		  },
		  {
			"id": 40403,
			"enchant": 3831
		  },
		  {
			"id": 40550,
			"enchant": 3832,
			"gems": [
			  42142,
			  42142
			]
		  },
		  {
			"id": 40330,
			"enchant": 3845,
			"gems": [
			  42142,
			  0
			]
		  },
		  {
			"id": 40552,
			"enchant": 3604,
			"gems": [
			  39996,
			  0
			]
		  },
		  {
			"id": 40278,
			"gems": [
			  39996,
			  39996
			]
		  },
		  {
			"id": 40556,
			"enchant": 3823,
			"gems": [
			  39996,
			  40037
			]
		  },
		  {
			"id": 40591,
			"enchant": 3606
		  },
		  {
			"id": 40075
		  },
		  {
			"id": 39401
		  },
		  {
			"id": 40256
		  },
		  {
			"id": 42987
		  },
		  {
			"id": 40384,
			"enchant": 3368
		  },
		  {},
		  {
			"id": 40207
		  }
  ]}`),
};

export const P2_BLOOD_BIS_PRESET = {
	name: 'P2 Blood',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecDeathknight>) => player.getTalentTree() == 0,
	gear: EquipmentSpec.fromJsonString(`{   "items": [
		{
			"id": 46115,
			"enchant": 3817,
			"gems": [
			  41398,
			  42702
			]
		  },
		  {
			"id": 45459,
			"gems": [
			  39996
			]
		  },
		  {
			"id": 46117,
			"enchant": 3808,
			"gems": [
			  39996
			]
		  },
		  {
			"id": 46032,
			"enchant": 3831,
			"gems": [
			  39996,
			  39996
			]
		  },
		  {
			"id": 46111,
			"enchant": 3832,
			"gems": [
			  42142,
			  42142
			]
		  },
		  {
			"id": 45663,
			"enchant": 3845,
			"gems": [
			  42142,
			  0
			]
		  },
		  {
			"id": 46113,
			"enchant": 3604,
			"gems": [
			  39996,
			  0
			]
		  },
		  {
			"id": 45241,
			"gems": [
			  39996,
			  45862,
			  39996
			]
		  },
		  {
			"id": 45134,
			"enchant": 3823,
			"gems": [
			  39996,
			  39996,
			  39996
			]
		  },
		  {
			"id": 45599,
			"enchant": 3606,
			"gems": [
			  39996,
			  39996
			]
		  },
		  {
			"id": 45534,
			"gems": [
			  39996
			]
		  },
		  {
			"id": 46048,
			"gems": [
			  39996
			]
		  },
		  {
			"id": 42987
		  },
		  {
			"id": 45931
		  },
		  {
			"id": 45516,
			"enchant": 3368,
			"gems": [
			  39996,
			  39996
			]
		  },
		  {},
		  {
			"id": 45254
		  }
  ]}`),
};

export const P3_BLOOD_BIS_PRESET = {
	name: 'P3 Blood',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecDeathknight>) => player.getTalentTree() == 0,
	gear: EquipmentSpec.fromJsonString(`{   "items": [
		{
			"id": 48493,
			"enchant": 3817,
			"gems": [
			  41285,
			  40142
			]
		  },
		  {
			"id": 47458,
			"gems": [
			  40142
			]
		  },
		  {
			"id": 48495,
			"enchant": 3808,
			"gems": [
			  40111
			]
		  },
		  {
			"id": 47546,
			"enchant": 3831,
			"gems": [
			  42142
			]
		  },
		  {
			"id": 47449,
			"enchant": 3832,
			"gems": [
			  49110,
			  42142,
			  40142
			]
		  },
		  {
			"id": 48008,
			"enchant": 3845,
			"gems": [
			  40111,
			  0
			]
		  },
		  {
			"id": 48492,
			"enchant": 3604,
			"gems": [
			  40142,
			  0
			]
		  },
		  {
			"id": 47429,
			"gems": [
			  40142,
			  40142,
			  40111
			]
		  },
		  {
			"id": 48494,
			"enchant": 3823,
			"gems": [
			  40142,
			  40111
			]
		  },
		  {
			"id": 45599,
			"enchant": 3606,
			"gems": [
			  40111,
			  40111
			]
		  },
		  {
			"id": 47993,
			"gems": [
			  40111,
			  45862
			]
		  },
		  {
			"id": 47413,
			"gems": [
			  40142
			]
		  },
		  {
			"id": 45931
		  },
		  {
			"id": 47464
		  },
		  {
			"id": 47446,
			"enchant": 3368,
			"gems": [
			  42142,
			  40141
			]
		  },
		  {},
		  {
			"id": 47673
		  }
  ]}`),
};

export const P1_UNHOLY_2H_PRERAID_PRESET = {
	name: 'Pre-Raid 2H Unholy',
	toolbar: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecDeathknight>) => player.getTalentTree() == 2 && player.getTalents().nervesOfColdSteel == 0,
	gear: EquipmentSpec.fromJsonString(`{ "items": [
		{
			"id": 41386,
			"enchant": 3817,
			"gems": [
				41400,
				49110
			]
		},
		{
			"id": 37397
		},
		{
			"id": 37627,
			"enchant": 3808,
			"gems": [
				39996
			]
		},
		{
			"id": 37647,
			"enchant": 3831
		},
		{
			"id": 39617,
			"enchant": 3832,
			"gems": [
				42142,
				39996
			]
		},
		{
			"id": 41355,
			"enchant": 3845,
			"gems": [
				0
			]
		},
		{
			"id": 39618,
			"enchant": 3604,
			"gems": [
				39996,
				0
			]
		},
		{
			"id": 40688,
			"gems": [
				39996,
				42142
			]
		},
		{
			"id": 37193,
			"enchant": 3823,
			"gems": [
				42142,
				39996
			]
		},
		{
			"id": 44306,
			"enchant": 3606,
			"gems": [
				39996,
				39996
			]
		},
		{
			"id": 37642
		},
		{
			"id": 44935
		},
		{
			"id": 40684
		},
		{
			"id": 42987
		},
		{
			"id": 41257,
			"enchant": 3368
		},
		{},
		{
			"id": 40867
		}
  ]}`),
};

export const P1_UNHOLY_2H_BIS_PRESET = {
	name: 'P1 2H Unholy',
	toolbar: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecDeathknight>) => player.getTalentTree() == 2 && player.getTalents().nervesOfColdSteel == 0,
	gear: EquipmentSpec.fromJsonString(`{ "items": [
		{
			"id": 44006,
			"enchant": 3817,
			"gems": [
				41400,
				49110
			]
		},
		{
			"id": 44664,
			"gems": [
				39996
			]
		},
		{
			"id": 40557,
			"enchant": 3808,
			"gems": [
				39996
			]
		},
		{
			"id": 40403,
			"enchant": 3831
		},
		{
			"id": 40550,
			"enchant": 3832,
			"gems": [
				42142,
				39996
			]
		},
		{
			"id": 40330,
			"enchant": 3845,
			"gems": [
				39996,
				0
			]
		},
		{
			"id": 40552,
			"enchant": 3604,
			"gems": [
				40038,
				0
			]
		},
		{
			"id": 40278,
			"gems": [
				42142,
				42142
			]
		},
		{
			"id": 40556,
			"enchant": 3823,
			"gems": [
				39996,
				39996
			]
		},
		{
			"id": 40591,
			"enchant": 3606
		},
		{
			"id": 39401
		},
		{
			"id": 40075
		},
		{
			"id": 40256
		},
		{
			"id": 42987
		},
		{
			"id": 40384,
			"enchant": 3368
		},
		{},
		{
			"id": 40207
		}
  ]}`),
};

export const P1_UNHOLY_DW_PRERAID_PRESET = {
	name: 'Pre-Raid DW Unholy',
	toolbar: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecDeathknight>) => player.getTalentTree() == 2 && player.getTalents().nervesOfColdSteel > 0,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 41386,
			"enchant": 3817,
			"gems": [
				41400,
				49110
			]
		},
		{
			"id": 37397
		},
		{
			"id": 37627,
			"enchant": 3808,
			"gems": [
				39996
			]
		},
		{
			"id": 37647,
			"enchant": 3831
		},
		{
			"id": 39617,
			"enchant": 3832,
			"gems": [
				42142,
				39996
			]
		},
		{
			"id": 41355,
			"enchant": 3845,
			"gems": [
				0
			]
		},
		{
			"id": 39618,
			"enchant": 3604,
			"gems": [
				39996,
				0
			]
		},
		{
			"id": 40688,
			"gems": [
				39996,
				42142
			]
		},
		{
			"id": 37193,
			"enchant": 3823,
			"gems": [
				42142,
				39996
			]
		},
		{
			"id": 44306,
			"enchant": 3606,
			"gems": [
				39996,
				39996
			]
		},
		{
			"id": 37642
		},
		{
			"id": 44935
		},
		{
			"id": 40684
		},
		{
			"id": 42987
		},
		{
			"id": 41383,
			"enchant": 3368
		},
		{
			"id": 40703,
			"enchant": 3368
		},
		{
			"id": 40867
		}
  ]}`),
};

export const P1_UNHOLY_DW_BIS_PRESET = {
	name: 'P1 DW Unholy',
	toolbar: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecDeathknight>) => player.getTalentTree() == 2 && player.getTalents().nervesOfColdSteel > 0,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 44006,
			"enchant": 3817,
			"gems": [
				41398,
				42702
			]
		},
		{
			"id": 39421
		},
		{
			"id": 40557,
			"enchant": 3808,
			"gems": [
				39996
			]
		},
		{
			"id": 40403,
			"enchant": 3831
		},
		{
			"id": 40550,
			"enchant": 3832,
			"gems": [
				42142,
				39996
			]
		},
		{
			"id": 40330,
			"enchant": 3845,
			"gems": [
				39996,
				0
			]
		},
		{
			"id": 40347,
			"enchant": 3604,
			"gems": [
				39996,
				0
			]
		},
		{
			"id": 40278,
			"gems": [
				42142,
				42142
			]
		},
		{
			"id": 40294,
			"enchant": 3823
		},
		{
			"id": 39706,
			"enchant": 3606,
			"gems": [
				39996
			]
		},
		{
			"id": 39401
		},
		{
			"id": 40075
		},
		{
			"id": 37390
		},
		{
			"id": 42987
		},
		{
			"id": 40402,
			"enchant": 3368
		},
		{
			"id": 40491,
			"enchant": 3368
		},
		{
			"id": 42620
		}
  ]}`),
};

export const P2_UNHOLY_DW_BIS_PRESET = {
	name: 'P2 DW Unholy',
	toolbar: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecDeathknight>) => player.getTalentTree() == 2 && player.getTalents().nervesOfColdSteel > 0,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 45472,
			"enchant": 3817,
			"gems": [
			  41398,
			  40041
			]
		  },
		  {
			"id": 46040,
			"gems": [
			  39996
			]
		  },
		  {
			"id": 46117,
			"enchant": 3808,
			"gems": [
			  39996
			]
		  },
		  {
			"id": 45588,
			"enchant": 3831,
			"gems": [
			  39996
			]
		  },
		  {
			"id": 46111,
			"enchant": 3832,
			"gems": [
			  42142,
			  42142
			]
		  },
		  {
			"id": 45663,
			"enchant": 3845,
			"gems": [
			  39996,
			  0
			]
		  },
		  {
			"id": 45481,
			"enchant": 3604,
			"gems": [
			  0
			]
		  },
		  {
			"id": 45241,
			"gems": [
			  42142,
			  45862,
			  39996
			]
		  },
		  {
			"id": 45134,
			"enchant": 3823,
			"gems": [
			  40041,
			  39996,
			  40022
			]
		  },
		  {
			"id": 45599,
			"enchant": 3606,
			"gems": [
			  39996,
			  39996
			]
		  },
		  {
			"id": 45534,
			"gems": [
			  39996
			]
		  },
		  {
			"id": 45250
		  },
		  {
			"id": 45609
		  },
		  {
			"id": 42987
		  },
		  {
			"id": 46097,
			"enchant": 3368,
			"gems": [
			  39996
			]
		  },
		  {
			"id": 46036,
			"enchant": 3368,
			"gems": [
			  39996
			]
		  },
		  {
			"id": 45254
		  }
  ]}`),
};

export const P3_UNHOLY_DW_BIS_PRESET = {
	name: 'P3 DW Unholy',
	toolbar: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecDeathknight>) => player.getTalentTree() == 2 && player.getTalents().nervesOfColdSteel > 0,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 48493,
			"enchant": 3817,
			"gems": [
			  41398,
			  40146
			]
		  },
		  {
			"id": 47458,
			"gems": [
			  40146
			]
		  },
		  {
			"id": 48495,
			"enchant": 3808,
			"gems": [
			  40111
			]
		  },
		  {
			"id": 47548,
			"enchant": 3831,
			"gems": [
			  40111
			]
		  },
		  {
			"id": 48491,
			"enchant": 3832,
			"gems": [
			  42142,
			  42142
			]
		  },
		  {
			"id": 45663,
			"enchant": 3845,
			"gems": [
			  40111,
			  0
			]
		  },
		  {
			"id": 48492,
			"enchant": 3604,
			"gems": [
			  40146,
			  0
			]
		  },
		  {
			"id": 47429,
			"gems": [
			  40111,
			  45862,
			  40111
			]
		  },
		  {
			"id": 47465,
			"enchant": 3823,
			"gems": [
			  49110,
			  40111,
			  40146
			]
		  },
		  {
			"id": 45599,
			"enchant": 3606,
			"gems": [
			  40111,
			  40111
			]
		  },
		  {
			"id": 47413,
			"gems": [
			  40146
			]
		  },
		  {
			"id": 45534,
			"gems": [
			  42142
			]
		  },
		  {
			"id": 47464
		  },
		  {
			"id": 45609
		  },
		  {
			"id": 47528,
			"enchant": 3368,
			"gems": [
			  40111
			]
		  },
		  {
			"id": 47528,
			"enchant": 3368,
			"gems": [
			  40111
			]
		  },
		  {
			"id": 47673
		  }
  ]}`),
};

export const P1_FROST_PRE_BIS_PRESET = {
	name: 'Pre-Raid Frost',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecDeathknight>) => player.getTalentTree() == 1,
	gear: EquipmentSpec.fromJsonString(`{   "items": [
		{
			"id": 41386,
			"enchant": 3817,
			"gems": [
				41398,
				49110
			]
		},
		{
			"id": 42645,
			"gems": [
				42142
			]
		},
		{
			"id": 34388,
			"enchant": 3808,
			"gems": [
				39996,
				39996
			]
		},
		{
			"id": 37647,
			"enchant": 3831
		},
		{
			"id": 39617,
			"enchant": 3832,
			"gems": [
				42142,
				39996
			]
		},
		{
			"id": 41355,
			"enchant": 3845,
			"gems": [
				0
			]
		},
		{
			"id": 39618,
			"enchant": 3604,
			"gems": [
				39996,
				0
			]
		},
		{
			"id": 37171,
			"gems": [
				39996,
				39996
			]
		},
		{
			"id": 37193,
			"enchant": 3823,
			"gems": [
				42142,
				39996
			]
		},
		{
			"id": 44306,
			"enchant": 3606,
			"gems": [
				39996,
				39996
			]
		},
		{
			"id": 42642,
			"gems": [
				39996
			]
		},
		{
			"id": 44935
		},
		{
			"id": 40684
		},
		{
			"id": 42987
		},
		{
			"id": 41383,
			"enchant": 3370
		},
		{
			"id": 43611,
			"enchant": 3368
		},
		{
			"id": 40715
		}
  ]}`),
};

export const P1_FROST_BIS_PRESET = {
	name: 'P1 Frost',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecDeathknight>) => player.getTalentTree() == 1,
	gear: EquipmentSpec.fromJsonString(`{   "items": [
		{
			"id": 44006,
			"enchant": 3817,
			"gems": [
				41398,
				42702
			]
		},
		{
			"id": 44664,
			"gems": [
				39996
			]
		},
		{
			"id": 40557,
			"enchant": 3808,
			"gems": [
				39996
			]
		},
		{
			"id": 40403,
			"enchant": 3831
		},
		{
			"id": 40550,
			"enchant": 3832,
			"gems": [
				42142,
				39996
			]
		},
		{
			"id": 40330,
			"enchant": 3845,
			"gems": [
				39996,
				0
			]
		},
		{
			"id": 40552,
			"enchant": 3604,
			"gems": [
				39996,
				0
			]
		},
		{
			"id": 40278,
			"gems": [
				39996,
				42142
			]
		},
		{
			"id": 40556,
			"enchant": 3823,
			"gems": [
				42142,
				39996
			]
		},
		{
			"id": 40591,
			"enchant": 3606
		},
		{
			"id": 39401
		},
		{
			"id": 40075
		},
		{
			"id": 40256
		},
		{
			"id": 42987
		},
		{
			"id": 40189,
			"enchant": 3370
		},
		{
			"id": 40189,
			"enchant": 3368
		},
		{
			"id": 40207
		}
  ]}`),
};

export const P2_FROST_BIS_PRESET = {
	name: 'P2 Frost',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecDeathknight>) => player.getTalentTree() == 1,
	gear: EquipmentSpec.fromJsonString(`{   "items": [
		{
			"id": 46115,
			"enchant": 3817,
			"gems": [
			  41398,
			  42702
			]
		  },
		  {
			"id": 45459,
			"gems": [
			  39996
			]
		  },
		  {
			"id": 46117,
			"enchant": 3808,
			"gems": [
			  39996
			]
		  },
		  {
			"id": 46032,
			"enchant": 3831,
			"gems": [
			  39996,
			  39996
			]
		  },
		  {
			"id": 46111,
			"enchant": 3832,
			"gems": [
			  42142,
			  42142
			]
		  },
		  {
			"id": 45663,
			"enchant": 3845,
			"gems": [
			  39996,
			  0
			]
		  },
		  {
			"id": 46113,
			"enchant": 3604,
			"gems": [
			  39996,
			  0
			]
		  },
		  {
			"id": 45241,
			"gems": [
			  42142,
			  45862,
			  39996
			]
		  },
		  {
			"id": 45134,
			"enchant": 3823,
			"gems": [
			  39996,
			  39996,
			  39996
			]
		  },
		  {
			"id": 45599,
			"enchant": 3606,
			"gems": [
			  39996,
			  39996
			]
		  },
		  {
			"id": 45608,
			"gems": [
			  39996
			]
		  },
		  {
			"id": 45534,
			"gems": [
			  39996
			]
		  },
		  {
			"id": 45931
		  },
		  {
			"id": 42987
		  },
		  {
			"id": 46097,
			"enchant": 3370,
			"gems": [
			  39996
			]
		  },
		  {
			"id": 46097,
			"enchant": 3368,
			"gems": [
			  39996
			]
		  },
		  {
			"id": 40207
		  }
  ]}`),
};

export const P3_FROST_BIS_PRESET = {
	name: 'P3 Frost',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecDeathknight>) => player.getTalentTree() == 1,
	gear: EquipmentSpec.fromJsonString(`{   "items": [
		{
			"id": 48493,
			"enchant": 3817,
			"gems": [
			  41398,
			  40142
			]
		  },
		  {
			"id": 45459,
			"gems": [
			  40111
			]
		  },
		  {
			"id": 48495,
			"enchant": 3808,
			"gems": [
			  40111
			]
		  },
		  {
			"id": 47548,
			"enchant": 3831,
			"gems": [
			  40111
			]
		  },
		  {
			"id": 48491,
			"enchant": 3832,
			"gems": [
			  42142,
			  42142
			]
		  },
		  {
			"id": 45663,
			"enchant": 3845,
			"gems": [
			  40111,
			  0
			]
		  },
		  {
			"id": 47492,
			"enchant": 3604,
			"gems": [
			  49110,
			  40111,
			  0
			]
		  },
		  {
			"id": 45241,
			"gems": [
			  40111,
			  42142,
			  40111
			]
		  },
		  {
			"id": 48494,
			"enchant": 3823,
			"gems": [
			  40142,
			  40111
			]
		  },
		  {
			"id": 47473,
			"enchant": 3606,
			"gems": [
			  40142,
			  40111
			]
		  },
		  {
			"id": 46966,
			"gems": [
			  40111
			]
		  },
		  {
			"id": 45534,
			"gems": [
			  40111
			]
		  },
		  {
			"id": 47464
		  },
		  {
			"id": 45931
		  },
		  {
			"id": 47528,
			"enchant": 3370,
			"gems": [
			  40111
			]
		  },
		  {
			"id": 47528,
			"enchant": 3368,
			"gems": [
			  40111
			]
		  },
		  {
			"id": 40207
		  }
  ]}`),
};

export const P1_FROSTSUBUNH_BIS_PRESET = {
	name: 'P1 Frost Sub Unh',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecDeathknight>) => player.getTalentTree() == 1,
	gear: EquipmentSpec.fromJsonString(`{   "items": [
		{
			"id": 44006,
			"enchant": 3817,
			"gems": [
				41398,
				42702
			]
		},
		{
			"id": 44664,
			"gems": [
				40003
			]
		},
		{
			"id": 40557,
			"enchant": 3808,
			"gems": [
				40003
			]
		},
		{
			"id": 40403,
			"enchant": 3831
		},
		{
			"id": 40550,
			"enchant": 3832,
			"gems": [
				42142,
				40003
			]
		},
		{
			"id": 40330,
			"enchant": 3845,
			"gems": [
				39996,
				0
			]
		},
		{
			"id": 40552,
			"enchant": 3604,
			"gems": [
				40058,
				0
			]
		},
		{
			"id": 40278,
			"gems": [
				39996,
				42142
			]
		},
		{
			"id": 40556,
			"enchant": 3823,
			"gems": [
				42142,
				39996
			]
		},
		{
			"id": 40591,
			"enchant": 3606
		},
		{
			"id": 39401
		},
		{
			"id": 40075
		},
		{
			"id": 40256
		},
		{
			"id": 42987
		},
		{
			"id": 40189,
			"enchant": 3370
		},
		{
			"id": 40189,
			"enchant": 3368
		},
		{
			"id": 40207
		}
  ]}`),
};
