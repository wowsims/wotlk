import { Conjured, Consumes, EquipmentSpec, Flask, Food, Glyphs, Potions } from '../core/proto/common.js';
import { Player } from '../core/player.js';
import { SavedRotation, SavedTalents } from '../core/proto/ui.js';

import {
	Rogue_Options as RogueOptions,
	Rogue_Options_PoisonImbue as Poison,
	Rogue_Rotation as RogueRotation,
	Rogue_Rotation_AssassinationPriority,
	Rogue_Rotation_CombatBuilder,
	Rogue_Rotation_CombatPriority,
	Rogue_Rotation_Frequency,
	Rogue_Rotation_SubtletyBuilder,
	Rogue_Rotation_SubtletyPriority,
	RogueMajorGlyph,
} from '../core/proto/rogue.js';

import * as Tooltips from '../core/constants/tooltips.js';
import { APLRotation } from '../core/proto/apl.js';

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/wotlk/talent-calc and copy the numbers in the url.
export const CombatHackTalents = {
	name: 'Combat Axes/Swords',
	data: SavedTalents.create({
		talentsString: '00532010414-0252051000035015223100501251',
		glyphs: Glyphs.create({
			major1: RogueMajorGlyph.GlyphOfKillingSpree,
			major2: RogueMajorGlyph.GlyphOfTricksOfTheTrade,
			major3: RogueMajorGlyph.GlyphOfRupture,
		})
	}),
};

export const CombatCQCTalents = {
	name: 'Combat Fists',
	data: SavedTalents.create({
		talentsString: '00532010414-0252051050035010223100501251',
		glyphs: Glyphs.create({
			major1: RogueMajorGlyph.GlyphOfKillingSpree,
			major2: RogueMajorGlyph.GlyphOfTricksOfTheTrade,
			major3: RogueMajorGlyph.GlyphOfRupture,
		})
	}),
};

export const AssassinationTalents137 = {
	name: 'Assassination 13/7',
	data: SavedTalents.create({
		talentsString: '005303104352100520103331051-005005003-502',
		glyphs: Glyphs.create({
			major1: RogueMajorGlyph.GlyphOfMutilate,
			major2: RogueMajorGlyph.GlyphOfTricksOfTheTrade,
			major3: RogueMajorGlyph.GlyphOfHungerForBlood,
		})
	}),
};

export const AssassinationTalents182 = {
	name: 'Assassination 18/2',
	data: SavedTalents.create({
		talentsString: '005303104352100520103331051-005005005003-2',
		glyphs: Glyphs.create({
			major1: RogueMajorGlyph.GlyphOfMutilate,
			major2: RogueMajorGlyph.GlyphOfTricksOfTheTrade,
			major3: RogueMajorGlyph.GlyphOfHungerForBlood,
		})
	}),
};

export const AssassinationTalentsBF = {
	name: 'Assassination Blade Flurry',
	data: SavedTalents.create({
		talentsString: '005303104352100520103231-005205005003001-501',
		glyphs: Glyphs.create({
			major1: RogueMajorGlyph.GlyphOfMutilate,
			major2: RogueMajorGlyph.GlyphOfTricksOfTheTrade,
			major3: RogueMajorGlyph.GlyphOfBladeFlurry,
		})
	}),
};

export const SubtletyTalents = {
	name: 'Subtlety',
	data: SavedTalents.create({
		talentsString: '30532010114--5022012030321121350115031151',
		glyphs: Glyphs.create({
			major1: RogueMajorGlyph.GlyphOfEviscerate,
			major2: RogueMajorGlyph.GlyphOfRupture,
			major3: RogueMajorGlyph.GlyphOfTricksOfTheTrade,
		})
	}),
}

export const HemoSubtletyTalents = {
	name: 'Hemo Sub',
	data: SavedTalents.create({
		talentsString: '30532010135--502201203032112135011503122',
		glyphs: Glyphs.create({
			major1: RogueMajorGlyph.GlyphOfEviscerate,
			major2: RogueMajorGlyph.GlyphOfRupture,
			major3: RogueMajorGlyph.GlyphOfTricksOfTheTrade,
		})
	}),
}

export const DefaultRotation = RogueRotation.create({
	exposeArmorFrequency: Rogue_Rotation_Frequency.Never,
	minimumComboPointsExposeArmor: 2,
	tricksOfTheTradeFrequency: Rogue_Rotation_Frequency.Maintain,
	assassinationFinisherPriority: Rogue_Rotation_AssassinationPriority.EnvenomRupture,
	combatBuilder: Rogue_Rotation_CombatBuilder.SinisterStrike,
	combatFinisherPriority: Rogue_Rotation_CombatPriority.RuptureEviscerate,
	subtletyBuilder: Rogue_Rotation_SubtletyBuilder.Hemorrhage,
	subtletyFinisherPriority: Rogue_Rotation_SubtletyPriority.SubtletyEviscerate,
	minimumComboPointsPrimaryFinisher: 4,
	minimumComboPointsSecondaryFinisher: 4,
});

export const DefaultOptions = RogueOptions.create({
	mhImbue: Poison.DeadlyPoison,
	ohImbue: Poison.InstantPoison,
	applyPoisonsManually: false,
	startingOverkillDuration: 20,
	honorOfThievesCritRate: 400,
});

export const DefaultConsumes = Consumes.create({
	defaultPotion: Potions.PotionOfSpeed,
	prepopPotion: Potions.PotionOfSpeed,
	defaultConjured: Conjured.ConjuredRogueThistleTea,
	flask: Flask.FlaskOfEndlessRage,
	food: Food.FoodMegaMammothMeal,
});

export const ROTATION_APL_PRESET  = {
	name: 'Basic APL',
	rotation: SavedRotation.create({
		specRotationOptionsJson: RogueRotation.toJsonString(DefaultRotation),
		rotation: APLRotation.fromJsonString('{ "enabled": true, "priorityList": [] }')
	}),
};

export const P2_PRESET_ASSASSINATION = {
	name: 'P2 Assassination',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<any>) => player.getTalentTree() == 0,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{"id":46125,"enchant":3817,"gems":[41398,39999]},
		  {"id":45517,"gems":[39999]},
		  {"id":45245,"enchant":3808,"gems":[39999,39999]},
		  {"id":45461,"enchant":3605,"gems":[40053]},
		  {"id":45473,"enchant":3832,"gems":[40053,42702,39999]},
		  {"id":45611,"enchant":3845,"gems":[40053,0]},
		  {"id":46124,"enchant":3604,"gems":[40003,0]},
		  {"id":46095,"enchant":3599,"gems":[39999,39999,39999]},
		  {"id":45536,"enchant":3823,"gems":[39999,39999,39999]},
		  {"id":45564,"enchant":3606,"gems":[39999,39999]},
		  {"id":45608,"gems":[39999]},
		  {"id":45456,"gems":[39999]},
		  {"id":45609},
		  {"id":46038},
		  {"id":45484,"enchant":3789,"gems":[40003]},
		  {"id":45484,"enchant":3789,"gems":[40003]},
		  {"id":45570,"enchant":3608}
	]}`),
};

export const P2_PRESET_COMBAT = {
	name: 'P2 Combat',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<any>) => player.getTalentTree() == 1,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{"id":46125,"enchant":3817,"gems":[41398,39999]},
		  {"id":45517,"gems":[39999]},
		  {"id":46127,"enchant":3808,"gems":[39999]},
		  {"id":45461,"enchant":3605,"gems":[40053]},
		  {"id":45473,"enchant":3832,"gems":[40053,42702,39999]},
		  {"id":45611,"enchant":3845,"gems":[40044,0]},
		  {"id":46043,"enchant":3604,"gems":[39999,40053,0]},
		  {"id":46095,"enchant":3599,"gems":[39999,39999,39999]},
		  {"id":45536,"enchant":3823,"gems":[39999,39999,39999]},
		  {"id":45564,"enchant":3606,"gems":[39999,39999]},
		  {"id":45608,"gems":[39999]},
		  {"id":46048,"gems":[39999]},
		  {"id":45609},
		  {"id":45931},
		  {"id":45132,"enchant":3789,"gems":[40053]},
		  {"id":45484,"enchant":3789,"gems":[40003]},
		  {"id":45296,"gems":[40053]}
	]}`),
};

export const P3_PRESET_ASSASSINATION = {
	name: 'P3 Assassination',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<any>) => player.getTalentTree() == 0,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		  {"id":48230,"enchant":3817,"gems":[41398,49110]},
		  {"id":47060,"gems":[40114]},
		  {"id":48228,"enchant":3808,"gems":[40114]},
		  {"id":47545,"enchant":3605,"gems":[40114]},
		  {"id":48232,"enchant":3832,"gems":[40114,40114]},
		  {"id":47155,"enchant":3845,"gems":[40114,40114,0]},
		  {"id":48231,"enchant":3604,"gems":[40114,0]},
		  {"id":47112,"enchant":3599,"gems":[40156,40114,40114]},
		  {"id":46975,"enchant":3823,"gems":[40118,40118,40118]},
		  {"id":47077,"enchant":3606,"gems":[40156,40114]},
		  {"id":47075,"gems":[40114]},
		  {"id":45608,"gems":[40114]},
		  {"id":47131},
		  {"id":45609},
		  {"id":46969,"enchant":3789,"gems":[40156]},
		  {"id":46969,"enchant":3789,"gems":[40156]},
		  {"id":47521,"gems":[40156]}
		]}`),
};

export const P3_PRESET_COMBAT = {
	name: 'P3 Combat',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<any>) => player.getTalentTree() == 1,
	gear: EquipmentSpec.fromJsonString(`{"items": [
        {"id":48230,"enchant":3817,"gems":[41398,49110]},
        {"id":47060,"gems":[40114]},
        {"id":48228,"enchant":3808,"gems":[40114]},
        {"id":47545,"enchant":3605,"gems":[40114]},
        {"id":48232,"enchant":3832,"gems":[40114,40114]},
        {"id":47155,"enchant":3845,"gems":[40114,40114,0]},
        {"id":48231,"enchant":3604,"gems":[40114,0]},
        {"id":47112,"enchant":3599,"gems":[40157,40114,40114]},
        {"id":46975,"enchant":3823,"gems":[40114,40114,40114]},
        {"id":47077,"enchant":3606,"gems":[40157,40114]},
        {"id":47075,"gems":[40114]},
        {"id":47934,"gems":[40157]},
        {"id":47131},
        {"id":45609},
        {"id":47156,"enchant":3789,"gems":[40157]},
        {"id":47001,"enchant":3789,"gems":[40157]},
        {"id":47521,"gems":[40157]}
	]}`),
};

export const PRERAID_PRESET_ASSASSINATION = {
	name: 'Pre-Raid Assassination',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<any>) => player.getTalentTree() == 0,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{"id":42550,"enchant":3817,"gems":[41398,40058]},
		{"id":40678},
		{"id":43481,"enchant":3808},
		{"id":38614,"enchant":3605},
		{"id":39558,"enchant":3832,"gems":[40003,42702]},
		{"id":34448,"enchant":3845,"gems":[40003,0]},
		{"id":39560,"enchant":3604,"gems":[40058,0]},
		{"id":40694,"gems":[40003,40003]},
		{"id":37644,"enchant":3823},
		{"id":34575,"enchant":3606,"gems":[40003]},
		{"id":40586},
		{"id":37642},
		{"id":40684},
		{"id":44253},
		{"id":37856,"enchant":3789},
		{"id":37667,"enchant":3789},
		{"id":43612}
  ]}`),
};

export const PRERAID_PRESET_COMBAT = {
	name: 'Pre-Raid Combat',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<any>) => player.getTalentTree() == 1,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{"id":42550,"enchant":3817,"gems":[41398,40014]},
		{"id":40678},
		{"id":37139,"enchant":3808,"gems":[39999]},
		{"id":34241,"enchant":3605,"gems":[40014]},
		{"id":39558,"enchant":3832,"gems":[39999,40014]},
		{"id":34448,"enchant":3845,"gems":[39999,0]},
		{"id":39560,"enchant":3604,"gems":[40014,0]},
		{"id":40694,"gems":[42702,39999]},
		{"id":37644,"enchant":3823},
		{"id":34575,"enchant":3606,"gems":[39999]},
		{"id":40586},
		{"id":37642},
		{"id":40684},
		{"id":44253},
		{"id":37693,"enchant":3789},
		{"id":37856,"enchant":3789},
		{"id":44504,"gems":[40053]}
  ]}`),
}

export const P1_PRESET_ASSASSINATION = {
	name: 'P1 Assassination',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<any>) => player.getTalentTree() == 0,
	gear: EquipmentSpec.fromJsonString(`{ "items": [
		{"id":40499,"enchant":3817,"gems":[41398,42702]},
		{"id":44664,"gems":[40003]},
		{"id":40502,"enchant":3808,"gems":[40003]},
		{"id":40403,"enchant":3605},
		{"id":40539,"enchant":3832,"gems":[40003]},
		{"id":39765,"enchant":3845,"gems":[40003,0]},
		{"id":40496,"enchant":3604,"gems":[40053,0]},
		{"id":40260,"gems":[39999]},
		{"id":40500,"enchant":3823,"gems":[40003,40003]},
		{"id":39701,"enchant":3606},
		{"id":40074},
		{"id":40474},
		{"id":40684},
		{"id":44253},
		{"id":39714,"enchant":3789},
		{"id":40386,"enchant":3789},
		{"id":40385}
  ]}`),
}

export const P1_PRESET_HEMO_SUB = {
	name: "P1 Hemo Sub",
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<any>) => player.getTalentTree() == 2,
	gear: EquipmentSpec.fromJsonString(`{ "items": [
		{"id":40499,"enchant":3817,"gems":[41398,42702]},
		{"id":44664,"gems":[40029]},
		{"id":40502,"enchant":3808,"gems":[40003]},
		{"id":40403,"enchant":3605},
		{"id":40539,"enchant":3832,"gems":[39999]},
		{"id":40186,"enchant":3845,"gems":[0]},
		{"id":40541,"enchant":3604,"gems":[0]},
		{"id":40205,"gems":[40003]},
		{"id":44011,"enchant":3823,"gems":[40003,40034]},
		{"id":39701,"enchant":3606},
		{"id":40074},
		{"id":40474},
		{"id":40256},
		{"id":44253},
		{"id":40383,"enchant":3789},
		{"id":39714,"enchant":3789},
		{"id":40385}
  ]}`),
}

export const P1_PRESET_COMBAT = {
	name: 'P1 Combat',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<any>) => player.getTalentTree() == 1,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{"id":40499,"enchant":3817,"gems":[41398,42702]},
		{"id":44664,"gems":[39999]},
		{"id":40502,"enchant":3808,"gems":[39999]},
		{"id":40403,"enchant":3605},
		{"id":40539,"enchant":3832,"gems":[39999]},
		{"id":39765,"enchant":3845,"gems":[39999,0]},
		{"id":40541,"enchant":3604,"gems":[0]},
		{"id":40205,"gems":[39999]},
		{"id":44011,"enchant":3823,"gems":[39999,39999]},
		{"id":39701,"enchant":3606},
		{"id":40074},
		{"id":40474},
		{"id":40684},
		{"id":44253},
		{"id":40383,"enchant":3789},
		{"id":39714,"enchant":3789},
		{"id":40385}
  ]}`),
}

export const P2_PRESET_HEMO_SUB = {
	name: "P2 Hemo Sub",
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<any>) => player.getTalentTree() == 2,
	gear: EquipmentSpec.fromJsonString(`{ "items": [
		{"id":46125,"enchant":3817,"gems":[41398,42143]},
		{"id":45517,"gems":[49110]},
		{"id":45245,"enchant":3808,"gems":[40023,40003]},
		{"id":45461,"enchant":3605,"gems":[40044]},
		{"id":45473,"enchant":3832,"gems":[40044,40023,40003]},
		{"id":45611,"enchant":3845,"gems":[40044,0]},
		{"id":46124,"enchant":3604,"gems":[39997,0]},
		{"id":46095,"enchant":3599,"gems":[42143,42143,39997]},
		{"id":45536,"enchant":3823,"gems":[40044,39997,40023]},
		{"id":45564,"enchant":3606,"gems":[40023,40003]},
		{"id":45608,"gems":[39997]},
		{"id":46048,"gems":[39997]},
		{"id":45609},
		{"id":45931},
		{"id":45132,"enchant":3789,"gems":[40044]},
		{"id":45484,"enchant":3789,"gems":[39997]},
		{"id":45296,"gems":[39997]}
  ]}`),
}

export const P3_PRESET_HEMO_SUB = {
	name: "P3 Hemo Sub",
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<any>) => player.getTalentTree() == 2,
	gear: EquipmentSpec.fromJsonString(`{ "items":[
		{"id":48235,"enchant":3817,"gems":[41398,49110]},
		{"id":47060,"gems":[40112]},
		{"id":48237,"enchant":3808,"gems":[40112]},
		{"id":47546,"enchant":3605,"gems":[40112]},
		{"id":47431,"enchant":3832,"gems":[40148,40130,40112]},
		{"id":45611,"enchant":3845,"gems":[40148,0]},
		{"id":48234,"enchant":3604,"gems":[40112,0]},
		{"id":47460,"gems":[40148,40112,40162]},
		{"id":47420,"enchant":3823,"gems":[40112,40112,40148]},
		{"id":47445,"enchant":3606,"gems":[40148,40112]},
		{"id":47443,
			"gems": [40112]},
		{"id":46048,"gems":[40112]},
		{"id":45609},
		{"id":47131},
		{"id":47475,"enchant":3789,"gems":[40148]},
		{"id":47416,"enchant":3789,"gems":[40148]},
		{"id":45296,"gems":[40112]}
	]}`),
}

export const P3_PRESET_DANCE_SUB = {
	name: "P3 Dance Sub",
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<any>) => player.getTalentTree() == 2,
	gear: EquipmentSpec.fromJsonString(`{ "items":[
		{"id":48235,"enchant":3817,"gems":[41398,49110]},
		{"id":47060,"gems":[40112]},
		{"id":48237,"enchant":3808,"gems":[40112]},
		{"id":47546,"enchant":3605,"gems":[40112]},
		{"id":47431,"enchant":3832,"gems":[40148,40130,40112]},
		{"id":45611,"enchant":3845,"gems":[40148,0]},
		{"id":48234,
			"enchant" :3604,
			"gems": [
				40112,
				0
			]
		},
		{"id":47460,"gems":[40148,40112,40162]},
		{"id":47420,"enchant":3823,"gems":[40112,40112,40148]},
		{"id":47445,"enchant":3606,"gems":[40148,40112]},
		{"id":47443,"gems":[40112]},
		{"id":46048,"gems":[40112]},
		{"id":45609},
		{"id":47131},
		{"id":47416,"enchant":3789,"gems":[40148]},
		{"id":47416,"enchant":3789,"gems":[40148]},
		{"id":45296,"gems":[40112]}
	]}`),
}
