import { CustomRotation, CustomSpell } from '../core/proto/common.js';
import { Consumes } from '../core/proto/common.js';
import { EquipmentSpec } from '../core/proto/common.js';
import { Flask } from '../core/proto/common.js';
import { Food } from '../core/proto/common.js';
import { Glyphs } from '../core/proto/common.js';
import { PetFood } from '../core/proto/common.js';
import { Potions } from '../core/proto/common.js';
import { SavedRotation, SavedTalents } from '../core/proto/ui.js';
import { APLRotation } from '../core/proto/apl.js';
import { ferocityDefault, ferocityBMDefault } from '../core/talents/hunter_pet.js';
import { Player } from '../core/player.js';

import {
	Hunter_Rotation as HunterRotation,
	Hunter_Rotation_RotationType as RotationType,
	Hunter_Rotation_StingType as StingType,
	Hunter_Rotation_SpellOption as SpellOption,
	Hunter_Options as HunterOptions,
	Hunter_Options_Ammo as Ammo,
	Hunter_Options_PetType as PetType,
	HunterMajorGlyph as MajorGlyph,
	HunterMinorGlyph as MinorGlyph,
} from '../core/proto/hunter.js';

import * as Tooltips from '../core/constants/tooltips.js';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/wotlk/talent-calc and copy the numbers in the url.
export const BeastMasteryTalents = {
	name: 'Beast Mastery',
	data: SavedTalents.create({
		talentsString: '51200201505112243120531251-025305101',
		glyphs: Glyphs.create({
			major1: MajorGlyph.GlyphOfBestialWrath,
			major2: MajorGlyph.GlyphOfSteadyShot,
			major3: MajorGlyph.GlyphOfSerpentSting,
			minor1: MinorGlyph.GlyphOfFeignDeath,
			minor2: MinorGlyph.GlyphOfRevivePet,
			minor3: MinorGlyph.GlyphOfMendPet,
		}),
	}),
};

export const MarksmanTalents = {
	name: 'Marksman',
	data: SavedTalents.create({
		talentsString: '502-035335131030013233035031051-5000002',
		glyphs: Glyphs.create({
			major1: MajorGlyph.GlyphOfSerpentSting,
			major2: MajorGlyph.GlyphOfSteadyShot,
			major3: MajorGlyph.GlyphOfChimeraShot,
			minor1: MinorGlyph.GlyphOfFeignDeath,
			minor2: MinorGlyph.GlyphOfRevivePet,
			minor3: MinorGlyph.GlyphOfMendPet,
		}),
	}),
};

export const SurvivalTalents = {
	name: 'Survival',
	data: SavedTalents.create({
		talentsString: '-015305101-5000032500033330532135301311',
		glyphs: Glyphs.create({
			major1: MajorGlyph.GlyphOfSerpentSting,
			major2: MajorGlyph.GlyphOfExplosiveShot,
			major3: MajorGlyph.GlyphOfKillShot,
			minor1: MinorGlyph.GlyphOfFeignDeath,
			minor2: MinorGlyph.GlyphOfRevivePet,
			minor3: MinorGlyph.GlyphOfMendPet,
		}),
	}),
};

export const DefaultRotation = HunterRotation.create({
	type: RotationType.SingleTarget,
	sting: StingType.SerpentSting,
	trapWeave: false,
	timeToTrapWeaveMs: 2000,
	viperStartManaPercent: 0.1,
	viperStopManaPercent: 0.3,
	multiDotSerpentSting: true,
	allowExplosiveShotDownrank: true,
	steadyShotMaxDelay: 300,
	customRotation: CustomRotation.create({
		spells: [
			CustomSpell.create({ spell: SpellOption.SerpentStingSpell }),
			CustomSpell.create({ spell: SpellOption.KillShot }),
			CustomSpell.create({ spell: SpellOption.ChimeraShot }),
			CustomSpell.create({ spell: SpellOption.BlackArrow }),
			CustomSpell.create({ spell: SpellOption.ExplosiveShot }),
			CustomSpell.create({ spell: SpellOption.AimedShot }),
			CustomSpell.create({ spell: SpellOption.ArcaneShot }),
			CustomSpell.create({ spell: SpellOption.SteadyShot }),
		],
	}),
});

export const ROTATION_PRESET_LEGACY_DEFAULT = {
	name: 'Legacy Default',
	rotation: SavedRotation.create({
		specRotationOptionsJson: HunterRotation.toJsonString(DefaultRotation),
	}),
}
export const ROTATION_PRESET_BM = {
	name: 'BM',
	rotation: SavedRotation.create({
		specRotationOptionsJson: HunterRotation.toJsonString(HunterRotation.create({
			timeToTrapWeaveMs: 2000,
		})),
		rotation: APLRotation.fromJsonString(`{
			"enabled": true,
      		"prepullActions": [
      		  {
      		    "action": {
      		      "castSpell": {"spellId":{"otherId":"OtherActionPotion"}}
      		    },
      		    "doAt": "-1s"
      		  }
      		],
      		"priorityList": [
      		  {
      		    "action": {
      		      "condition": {"cmp":{"op":"OpGt","lhs":{"currentTime":{}},"rhs":{"const":{"val":"10s"}}}},
      		      "autocastOtherCooldowns": {}
      		    }
      		  },
      		  {
      		    "action": {
      		      "condition": {"and":{"vals":[{"not":{"val":{"auraIsActive":{"auraId":{"spellId":34074}}}}},{"cmp":{"op":"OpLt","lhs":{"currentManaPercent":{}},"rhs":{"const":{"val":"10%"}}}}]}},
      		      "castSpell": {"spellId":{"spellId":34074}}
      		    }
      		  },
      		  {
      		    "action": {
      		      "condition": {"and":{"vals":[{"not":{"val":{"auraIsActive":{"auraId":{"spellId":61847}}}}},{"cmp":{"op":"OpGt","lhs":{"currentManaPercent":{}},"rhs":{"const":{"val":"30%"}}}}]}},
      		      "castSpell": {"spellId":{"spellId":61847}}
      		    }
      		  },
      		  {
      		    "action": {
      		      "castSpell": {"spellId":{"spellId":61006}}
      		    }
      		  },
      		  {
      		    "hide": true,
      		    "action": {
      		      "condition": {"not":{"val":{"dotIsActive":{"spellId":{"spellId":49067}}}}},
      		      "castSpell": {"spellId":{"tag":1,"spellId":49067}}
      		    }
      		  },
      		  {
      		    "action": {
      		      "condition": {"not":{"val":{"dotIsActive":{"spellId":{"spellId":49001}}}}},
      		      "castSpell": {"spellId":{"spellId":49001}}
      		    }
      		  },
      		  {
      		    "action": {
      		      "castSpell": {"spellId":{"spellId":49050}}
      		    }
      		  },
      		  {
      		    "action": {
      		      "castSpell": {"spellId":{"spellId":49048}}
      		    }
      		  },
      		  {
      		    "hide": true,
      		    "action": {
      		      "castSpell": {"spellId":{"spellId":49045}}
      		    }
      		  },
      		  {
      		    "action": {
      		      "castSpell": {"spellId":{"spellId":49052}}
      		    }
      		  }
      		]
		}`),
	}),
};

export const ROTATION_PRESET_MM = {
	name: 'MM',
	rotation: SavedRotation.create({
		specRotationOptionsJson: HunterRotation.toJsonString(HunterRotation.create({
			timeToTrapWeaveMs: 2000,
		})),
		rotation: APLRotation.fromJsonString(`{
      		"enabled": true,
      		"prepullActions": [
      		  {
      		    "action": {
      		      "castSpell": {"spellId":{"otherId":"OtherActionPotion"}}
      		    },
      		    "doAt": "-1s"
      		  }
      		],
      		"priorityList": [
      		  {
      		    "action": {
      		      "condition": {"cmp":{"op":"OpGt","lhs":{"currentTime":{}},"rhs":{"const":{"val":"10s"}}}},
      		      "autocastOtherCooldowns": {}
      		    }
      		  },
      		  {
      		    "action": {
      		      "condition": {"and":{"vals":[{"not":{"val":{"auraIsActive":{"auraId":{"spellId":34074}}}}},{"cmp":{"op":"OpLt","lhs":{"currentManaPercent":{}},"rhs":{"const":{"val":"10%"}}}}]}},
      		      "castSpell": {"spellId":{"spellId":34074}}
      		    }
      		  },
      		  {
      		    "action": {
      		      "condition": {"and":{"vals":[{"not":{"val":{"auraIsActive":{"auraId":{"spellId":61847}}}}},{"cmp":{"op":"OpGt","lhs":{"currentManaPercent":{}},"rhs":{"const":{"val":"30%"}}}}]}},
      		      "castSpell": {"spellId":{"spellId":61847}}
      		    }
      		  },
      		  {
      		    "action": {
      		      "castSpell": {"spellId":{"spellId":61006}}
      		    }
      		  },
      		  {
      		    "action": {
      		      "condition": {"not":{"val":{"dotIsActive":{"spellId":{"spellId":49001}}}}},
      		      "castSpell": {"spellId":{"spellId":49001}}
      		    }
      		  },
      		  {
      		    "hide": true,
      		    "action": {
      		      "condition": {"not":{"val":{"dotIsActive":{"spellId":{"spellId":49067}}}}},
      		      "castSpell": {"spellId":{"tag":1,"spellId":49067}}
      		    }
      		  },
      		  {
      		    "action": {
      		      "castSpell": {"spellId":{"spellId":53209}}
      		    }
      		  },
      		  {
      		    "action": {
      		      "castSpell": {"spellId":{"spellId":49050}}
      		    }
      		  },
      		  {
      		    "action": {
      		      "castSpell": {"spellId":{"spellId":49048}}
      		    }
      		  },
      		  {
      		    "hide": true,
      		    "action": {
      		      "castSpell": {"spellId":{"spellId":49045}}
      		    }
      		  },
      		  {
      		    "action": {
      		      "castSpell": {"spellId":{"spellId":49052}}
      		    }
      		  }
      		]
		}`),
	}),
};

export const ROTATION_PRESET_SV = {
	name: 'SV',
	rotation: SavedRotation.create({
		specRotationOptionsJson: HunterRotation.toJsonString(HunterRotation.create({
			timeToTrapWeaveMs: 2000,
		})),
		rotation: APLRotation.fromJsonString(`{
      		"enabled": true,
      		"prepullActions": [
      		  {
      		    "action": {
      		      "castSpell": {"spellId":{"otherId":"OtherActionPotion"}}
      		    },
      		    "doAt": "-1s"
      		  }
      		],
      		"priorityList": [
      		  {
      		    "action": {
      		      "condition": {"cmp":{"op":"OpGt","lhs":{"currentTime":{}},"rhs":{"const":{"val":"10s"}}}},
      		      "autocastOtherCooldowns": {}
      		    }
      		  },
      		  {
      		    "action": {
      		      "condition": {"and":{"vals":[{"not":{"val":{"auraIsActive":{"auraId":{"spellId":34074}}}}},{"cmp":{"op":"OpLt","lhs":{"currentManaPercent":{}},"rhs":{"const":{"val":"10%"}}}}]}},
      		      "castSpell": {"spellId":{"spellId":34074}}
      		    }
      		  },
      		  {
      		    "action": {
      		      "condition": {"and":{"vals":[{"not":{"val":{"auraIsActive":{"auraId":{"spellId":61847}}}}},{"cmp":{"op":"OpGt","lhs":{"currentManaPercent":{}},"rhs":{"const":{"val":"30%"}}}}]}},
      		      "castSpell": {"spellId":{"spellId":61847}}
      		    }
      		  },
      		  {
      		    "action": {
      		      "castSpell": {"spellId":{"spellId":61006}}
      		    }
      		  },
      		  {
      		    "action": {
      		      "condition": {"not":{"val":{"dotIsActive":{"spellId":{"spellId":60053}}}}},
      		      "castSpell": {"spellId":{"spellId":60053}}
      		    }
      		  },
      		  {
      		    "action": {
      		      "condition": {"dotIsActive":{"spellId":{"spellId":60053}}},
      		      "castSpell": {"spellId":{"spellId":60052}}
      		    }
      		  },
      		  {
      		    "hide": true,
      		    "action": {
      		      "condition": {"not":{"val":{"dotIsActive":{"spellId":{"spellId":49067}}}}},
      		      "castSpell": {"spellId":{"tag":1,"spellId":49067}}
      		    }
      		  },
      		  {
      		    "action": {
      		      "condition": {"not":{"val":{"dotIsActive":{"spellId":{"spellId":49001}}}}},
      		      "castSpell": {"spellId":{"spellId":49001}}
      		    }
      		  },
      		  {
      		    "action": {
      		      "castSpell": {"spellId":{"spellId":63672}}
      		    }
      		  },
      		  {
      		    "action": {
      		      "castSpell": {"spellId":{"spellId":49050}}
      		    }
      		  },
      		  {
      		    "action": {
      		      "castSpell": {"spellId":{"spellId":49048}}
      		    }
      		  },
      		  {
      		    "action": {
      		      "castSpell": {"spellId":{"spellId":49052}}
      		    }
      		  }
      		]
		}`),
	}),
};

export const ROTATION_PRESET_AOE = {
	name: 'AOE',
	rotation: SavedRotation.create({
		specRotationOptionsJson: HunterRotation.toJsonString(HunterRotation.create({
			timeToTrapWeaveMs: 2000,
		})),
		rotation: APLRotation.fromJsonString(`{
      		"enabled": true,
      		"prepullActions": [
      		  {
      		    "action": {
      		      "castSpell": {"spellId":{"otherId":"OtherActionPotion"}}
      		    },
      		    "doAt": "-1s"
      		  }
      		],
      		"priorityList": [
      		  {
      		    "action": {
      		      "condition": {"cmp":{"op":"OpGt","lhs":{"currentTime":{}},"rhs":{"const":{"val":"10s"}}}},
      		      "autocastOtherCooldowns": {}
      		    }
      		  },
      		  {
      		    "action": {
      		      "condition": {"and":{"vals":[{"not":{"val":{"auraIsActive":{"auraId":{"spellId":34074}}}}},{"cmp":{"op":"OpLt","lhs":{"currentManaPercent":{}},"rhs":{"const":{"val":"10%"}}}}]}},
      		      "castSpell": {"spellId":{"spellId":34074}}
      		    }
      		  },
      		  {
      		    "action": {
      		      "condition": {"and":{"vals":[{"not":{"val":{"auraIsActive":{"auraId":{"spellId":61847}}}}},{"cmp":{"op":"OpGt","lhs":{"currentManaPercent":{}},"rhs":{"const":{"val":"30%"}}}}]}},
      		      "castSpell": {"spellId":{"spellId":61847}}
      		    }
      		  },
      		  {
      		    "hide": true,
      		    "action": {
      		      "multidot": {"spellId":{"spellId":49001},"maxDots":3,"maxOverlap":{"const":{"val":"0ms"}}}
      		    }
      		  },
      		  {
      		    "action": {
      		      "condition": {"not":{"val":{"dotIsActive":{"spellId":{"spellId":49067}}}}},
      		      "castSpell": {"spellId":{"tag":1,"spellId":49067}}
      		    }
      		  },
      		  {
      		    "action": {
      		      "castSpell": {"spellId":{"spellId":58434}}
      		    }
      		  }
      		]
		}`),
	}),
};

export const DefaultOptions = HunterOptions.create({
	ammo: Ammo.SaroniteRazorheads,
	useHuntersMark: true,
	petType: PetType.Wolf,
	petTalents: ferocityDefault,
	petUptime: 1,
	sniperTrainingUptime: 0.9,
});

export const BMDefaultOptions = HunterOptions.create({
	ammo: Ammo.SaroniteRazorheads,
	useHuntersMark: true,
	petType: PetType.Wolf,
	petTalents: ferocityBMDefault,
	petUptime: 1,
	sniperTrainingUptime: 0.9,
});

export const DefaultConsumes = Consumes.create({
	defaultPotion: Potions.PotionOfSpeed,
	flask: Flask.FlaskOfEndlessRage,
	food: Food.FoodFishFeast,
	petFood: PetFood.PetFoodSpicedMammothTreats,
});

export const MM_PRERAID_PRESET = {
	name: 'MM Preraid Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<any>) => player.getTalentTree() != 2,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 42551,
			"enchant": 3817,
			"gems": [
				41398,
				42143
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
			"id": 43566,
			"enchant": 3605
		},
		{
			"id": 39579,
			"enchant": 3832,
			"gems": [
				39997,
				49110
			]
		},
		{
			"id": 37170,
			"enchant": 3845,
			"gems": [
				0
			]
		},
		{
			"id": 39582,
			"enchant": 3604,
			"gems": [
				40014,
				0
			]
		},
		{
			"id": 37407,
			"enchant": 3601,
			"gems": [
				42143
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
				42143,
				39997
			]
		},
		{
			"id": 37685
		},
		{
			"id": 42642,
			"gems": [
				40044
			]
		},
		{
			"id": 40684
		},
		{
			"id": 44253
		},
		{
			"id": 44249,
			"enchant": 3827
		},
		{},
		{
			"id": 37191,
			"enchant": 3608
		}
	]}`),
};

export const MM_P1_PRESET = {
	name: 'MM P1 Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<any>) => player.getTalentTree() != 2,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 40543,
			"enchant": 3817,
			"gems": [
				41398,
				42143
			]
		},
		{
			"id": 44664,
			"gems": [
				42143
			]
		},
		{
			"id": 40507,
			"enchant": 3808,
			"gems": [
				39997
			]
		},
		{
			"id": 40403,
			"enchant": 3605
		},
		{
			"id": 43998,
			"enchant": 3832,
			"gems": [
				42143,
				39997
			]
		},
		{
			"id": 40282,
			"enchant": 3845,
			"gems": [
				39997,
				0
			]
		},
		{
			"id": 40541,
			"enchant": 3604,
			"gems": [
				0
			]
		},
		{
			"id": 40275,
			"enchant": 3601,
			"gems": [
				39997
			]
		},
		{
			"id": 40506,
			"enchant": 3823,
			"gems": [
				39997,
				49110
			]
		},
		{
			"id": 40549,
			"enchant": 3606
		},
		{
			"id": 40074
		},
		{
			"id": 40474
		},
		{
			"id": 40684
		},
		{
			"id": 44253
		},
		{
			"id": 40388,
			"enchant": 3827
		},
		{},
		{
			"id": 40385,
			"enchant": 3608
		}
	]}`),
};

export const MM_P2_PRESET = {
	name: 'MM P2 Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<any>) => player.getTalentTree() != 2,
	gear: EquipmentSpec.fromJsonString(`{"items": [
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
				42143
			]
		},
		{
			"id": 45300,
			"enchant": 3808,
			"gems": [
				40043
			]
		},
		{
			"id": 46032,
			"enchant": 3605,
			"gems": [
				42143,
				40043
			]
		},
		{
			"id": 45473,
			"enchant": 3832,
			"gems": [
				39997,
				39997,
				39997
			]
		},
		{
			"id": 45869,
			"enchant": 3845,
			"gems": [
				40044,
				0
			]
		},
		{
			"id": 45444,
			"enchant": 3604,
			"gems": [
				42143,
				39997,
				0
			]
		},
		{
			"id": 45467,
			"enchant": 3601,
			"gems": [
				39997
			]
		},
		{
			"id": 45536,
			"enchant": 3823,
			"gems": [
				39997,
				39997,
				39997
			]
		},
		{
			"id": 45244,
			"enchant": 3606,
			"gems": [
				39997,
				39997
			]
		},
		{
			"id": 45608,
			"gems": [
				39997
			]
		},
		{
			"id": 46322,
			"gems": [
				39997
			]
		},
		{
			"id": 45931
		},
		{
			"id": 46038
		},
		{
			"id": 45613,
			"enchant": 3827,
			"gems": [
				45879,
				39997
			]
		},
		{},
		{
			"id": 45570,
			"enchant": 3608
		}
	]}`),
};

export const MM_P3_PRESET = {
	name: 'MM P3 Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<any>) => player.getTalentTree() != 2,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
          "id": 48262,
          "enchant": 3817,
          "gems": [
            41398,
            40147
          ]
        },
        {
          "id": 47060,
          "gems": [
            42143
          ]
        },
        {
          "id": 48260,
          "enchant": 3808,
          "gems": [
            40112
          ]
        },
        {
          "id": 47545,
          "enchant": 3605,
          "gems": [
            40112
          ]
        },
        {
          "id": 46965,
          "enchant": 3832,
          "gems": [
            40112,
            40112,
            40112
          ]
        },
        {
          "id": 47074,
          "enchant": 3845,
          "gems": [
            40147,
            0
          ]
        },
        {
          "id": 48263,
          "enchant": 3604,
          "gems": [
            40148,
            0
          ]
        },
        {
          "id": 47153,
          "gems": [
            40148,
            42143,
            42143
          ]
        },
        {
          "id": 48261,
          "enchant": 3823,
          "gems": [
            49110,
            40112
          ]
        },
        {
          "id": 47109,
          "enchant": 3606,
          "gems": [
            40147,
            40147
          ]
        },
        {
          "id": 47075,
          "gems": [
            40112
          ]
        },
        {
          "id": 45608,
          "gems": [
            40112
          ]
        },
        {
          "id": 47131
        },
        {
          "id": 45931
        },
        {
          "id": 47239,
          "enchant": 3827,
          "gems": [
            40147,
            40112
          ]
        },
        {},
        {
          "id": 47521,
          "enchant": 3608,
          "gems": [
            40147
          ]
        }
	]}`),
};

export const MM_P4_PRESET = {
	name: 'MM P4 Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<any>) => player.getTalentTree() != 2,
	gear: EquipmentSpec.fromJsonString(`{"items": [
        {
          "id": 51286,
          "enchant": 3817,
          "gems": [
            41398,
            40117
          ]
        },
        {
          "id": 50633,
          "gems": [
            40117
          ]
        },
        {
          "id": 51288,
          "enchant": 3808,
          "gems": [
            40117
          ]
        },
        {
          "id": 47546,
          "enchant": 3605,
          "gems": [
            42153
          ]
        },
        {
          "id": 51289,
          "enchant": 3832,
          "gems": [
            40117,
            40117
          ]
        },
        {
          "id": 50655,
          "enchant": 3845,
          "gems": [
            40117,
            0
          ]
        },
        {
          "id": 51285,
          "enchant": 3604,
          "gems": [
            40117,
            0
          ]
        },
        {
          "id": 50688,
          "enchant": 3601,
          "gems": [
            40148,
            42153,
            42153
          ]
        },
        {
          "id": 50645,
          "enchant": 3823,
          "gems": [
            49110,
            40117,
            40147
          ]
        },
        {
          "id": 50607,
          "enchant": 3606,
          "gems": [
            40148,
            40148
          ]
        },
        {
          "id": 50618,
          "gems": [
            40117
          ]
        },
        {
          "id": 50402,
          "gems": [
            40148
          ]
        },
        {
          "id": 50363
        },
        {
          "id": 47131
        },
        {
          "id": 50735,
          "enchant": 3827,
          "gems": [
            40117,
            40117,
            40117
          ]
        },
        {},
        {
          "id": 50733,
          "enchant": 3608,
          "gems": [
            40117
          ]
        }
	]}`),
};

export const SV_PRERAID_PRESET = {
	name: 'SV Preraid Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<any>) => player.getTalentTree() == 2,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 42551,
			"enchant": 3817,
			"gems": [
				41398,
				42143
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
			"id": 43406,
			"enchant": 3605
		},
		{
			"id": 39579,
			"enchant": 3832,
			"gems": [
				39997,
				49110
			]
		},
		{
			"id": 37170,
			"enchant": 3845,
			"gems": [
				0
			]
		},
		{
			"id": 39582,
			"enchant": 3604,
			"gems": [
				39997,
				0
			]
		},
		{
			"id": 37407,
			"enchant": 3601,
			"gems": [
				42143
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
				42143,
				39997
			]
		},
		{
			"id": 37685
		},
		{
			"id": 42642,
			"gems": [
				39997
			]
		},
		{
			"id": 40684
		},
		{
			"id": 44253
		},
		{
			"id": 44249,
			"enchant": 3827
		},
		{},
		{
			"id": 37191,
			"enchant": 3608
		}
	]}`),
};

export const SV_P1_PRESET = {
	name: 'SV P1 Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<any>) => player.getTalentTree() == 2,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 40505,
			"enchant": 3817,
			"gems": [
				41398,
				42143
			]
		},
		{
			"id": 44664,
			"gems": [
				42143
			]
		},
		{
			"id": 40507,
			"enchant": 3808,
			"gems": [
				39997
			]
		},
		{
			"id": 40403,
			"enchant": 3605
		},
		{
			"id": 43998,
			"enchant": 3832,
			"gems": [
				42143,
				39997
			]
		},
		{
			"id": 40282,
			"enchant": 3845,
			"gems": [
				39997,
				0
			]
		},
		{
			"id": 40541,
			"enchant": 3604,
			"gems": [
				0
			]
		},
		{
			"id": 39762,
			"enchant": 3601,
			"gems": [
				39997
			]
		},
		{
			"id": 40331,
			"enchant": 3823,
			"gems": [
				39997,
				49110
			]
		},
		{
			"id": 40549,
			"enchant": 3606
		},
		{
			"id": 40074
		},
		{
			"id": 40474
		},
		{
			"id": 40684
		},
		{
			"id": 44253
		},
		{
			"id": 40388,
			"enchant": 3827
		},
		{},
		{
			"id": 40385,
			"enchant": 3608
		}
	]}`),
};

export const SV_P2_PRESET = {
	name: 'SV P2 Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<any>) => player.getTalentTree() == 2,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{
			"id": 45610,
			"enchant": 3817,
			"gems": [
				41398,
				40023
			]
		},
		{
			"id": 45517,
			"gems": [
				39997
			]
		},
		{
			"id": 45300,
			"enchant": 3808,
			"gems": [
				39997
			]
		},
		{
			"id": 46032,
			"enchant": 3605,
			"gems": [
				39997,
				40044
			]
		},
		{
			"id": 45473,
			"enchant": 3832,
			"gems": [
				39997,
				39997,
				45879
			]
		},
		{
			"id": 45869,
			"enchant": 3845,
			"gems": [
				40043,
				0
			]
		},
		{
			"id": 45444,
			"enchant": 3604,
			"gems": [
				39997,
				40023,
				0
			]
		},
		{
			"id": 46095,
			"gems": [
				42143,
				42143,
				42143
			]
		},
		{
			"id": 45536,
			"enchant": 3823,
			"gems": [
				39997,
				39997,
				39997
			]
		},
		{
			"id": 45244,
			"enchant": 3606,
			"gems": [
				39997,
				40023
			]
		},
		{
			"id": 45608,
			"gems": [
				39997
			]
		},
		{
			"id": 46322,
			"gems": [
				39997
			]
		},
		{
			"id": 44253
		},
		{
			"id": 45931
		},
		{
			"id": 45613,
			"enchant": 3827,
			"gems": [
				39997,
				39997
			]
		},
		{},
		{
			"id": 45570,
			"enchant": 3608
		}
	]}`),
};

export const SV_P3_PRESET = {
	name: 'SV P3 Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<any>) => player.getTalentTree() == 2,
	gear: EquipmentSpec.fromJsonString(`{"items": [
        {
          "id": 48262,
          "enchant": 3817,
          "gems": [
            41398,
            40147
          ]
        },
        {
          "id": 47060,
          "gems": [
            42143
          ]
        },
        {
          "id": 48260,
          "enchant": 3808,
          "gems": [
            40112
          ]
        },
        {
          "id": 47545,
          "enchant": 3605,
          "gems": [
            40112
          ]
        },
        {
          "id": 48264,
          "enchant": 3832,
          "gems": [
            40112,
            40147
          ]
        },
        {
          "id": 47074,
          "enchant": 3845,
          "gems": [
            40148,
            0
          ]
        },
        {
          "id": 48263,
          "enchant": 3604,
          "gems": [
            40148,
            0
          ]
        },
        {
          "id": 47153,
          "gems": [
            40147,
            42143,
            42143
          ]
        },
        {
          "id": 47191,
          "enchant": 3823,
          "gems": [
            49110,
            40147,
            40112
          ]
        },
        {
          "id": 47109,
          "enchant": 3606,
          "gems": [
            40112,
            40112
          ]
        },
        {
          "id": 47075,
          "gems": [
            40112
          ]
        },
        {
          "id": 45608,
          "gems": [
            40112
          ]
        },
        {
          "id": 47131
        },
        {
          "id": 44253
        },
        {
          "id": 47239,
          "enchant": 3827,
          "gems": [
            40147,
            40112
          ]
        },
        {},
        {
          "id": 47521,
          "enchant": 3608,
          "gems": [
            40112
          ]
        }
	]}`),
};

export const SV_P4_PRESET = {
	name: 'SV P4 Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<any>) => player.getTalentTree() == 2,
	gear: EquipmentSpec.fromJsonString(`{"items": [
        {
          "id": 51286,
          "enchant": 3817,
          "gems": [
            41398,
            40112
          ]
        },
        {
          "id": 50633,
          "gems": [
            40112
          ]
        },
        {
          "id": 51288,
          "enchant": 3808,
          "gems": [
            40112
          ]
        },
        {
          "id": 47546,
          "enchant": 3605,
          "gems": [
            42143
          ]
        },
        {
          "id": 51289,
          "enchant": 3832,
          "gems": [
            40112,
            40112
          ]
        },
        {
          "id": 50655,
          "enchant": 3845,
          "gems": [
            40112,
            0
          ]
        },
        {
          "id": 51285,
          "enchant": 3604,
          "gems": [
            40112,
            0
          ]
        },
        {
          "id": 50688,
          "enchant": 3601,
          "gems": [
            40148,
            42143,
            42143
          ]
        },
        {
          "id": 50645,
          "enchant": 3823,
          "gems": [
            49110,
            40112,
            40150
          ]
        },
        {
          "id": 50607,
          "enchant": 3606,
          "gems": [
            40148,
            40148
          ]
        },
        {
          "id": 50618,
          "gems": [
            45879
          ]
        },
        {
          "id": 50402,
          "gems": [
            40148
          ]
        },
        {
          "id": 50363
        },
        {
          "id": 47131
        },
        {
          "id": 50735,
          "enchant": 3827,
          "gems": [
            40112,
            40112,
            40112
          ]
        },
        {},
        {
          "id": 50733,
          "enchant": 3608,
          "gems": [
            40112
          ]
        }
	]}`),
};
