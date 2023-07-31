import { Consumes, Spec } from '../core/proto/common.js';
import { EquipmentSpec } from '../core/proto/common.js';
import { Flask } from '../core/proto/common.js';
import { Food } from '../core/proto/common.js';
import { Glyphs } from '../core/proto/common.js';
import { Potions } from '../core/proto/common.js';
import { SavedRotation, SavedTalents } from '../core/proto/ui.js';
import { APLRotation } from '../core/proto/apl.js';
import { Player } from '../core/player.js';

import {
	TankDeathknight_Rotation as TankDeathKnightRotation,
	TankDeathknight_Options as TankDeathKnightOptions,
	DeathknightMajorGlyph,
	DeathknightMinorGlyph,
	TankDeathknight_Rotation_Opener as Opener,
	TankDeathknight_Rotation_OptimizationSetting as OptimizationSetting,
	TankDeathknight_Rotation_BloodSpell as BloodSpell,
	TankDeathknight_Rotation_Presence as Presence,
} from '../core/proto/deathknight.js';

import * as Tooltips from '../core/constants/tooltips.js';

export const BloodTalents = {
	name: 'Blood',
	data: SavedTalents.create({
		talentsString: '005512153330030320102013-3050505000023-005',
		glyphs: Glyphs.create({
			major1: DeathknightMajorGlyph.GlyphOfDisease,
			major2: DeathknightMajorGlyph.GlyphOfRuneStrike,
			major3: DeathknightMajorGlyph.GlyphOfDarkCommand,
			minor1: DeathknightMinorGlyph.GlyphOfHornOfWinter,
			minor2: DeathknightMinorGlyph.GlyphOfBloodTap,
			minor3: DeathknightMinorGlyph.GlyphOfRaiseDead,
		}),
	}),
};

export const BloodAggroTalents = {
	name: 'Blood Aggro',
	data: SavedTalents.create({
		talentsString: '0355220530303303201020131301--0052003050032',
		glyphs: Glyphs.create({
			major1: DeathknightMajorGlyph.GlyphOfDancingRuneWeapon,
			major2: DeathknightMajorGlyph.GlyphOfRuneStrike,
			major3: DeathknightMajorGlyph.GlyphOfDarkCommand,
			minor1: DeathknightMinorGlyph.GlyphOfHornOfWinter,
			minor2: DeathknightMinorGlyph.GlyphOfBloodTap,
			minor3: DeathknightMinorGlyph.GlyphOfRaiseDead,
		}),
	}),
};

export const DoubleBuffBloodTalents = {
	name: '2B Blood',
	data: SavedTalents.create({
		talentsString: '005512153330030320102013-3050505000023201-002',
		glyphs: Glyphs.create({
			major1: DeathknightMajorGlyph.GlyphOfDisease,
			major2: DeathknightMajorGlyph.GlyphOfRuneStrike,
			major3: DeathknightMajorGlyph.GlyphOfDarkCommand,
			minor1: DeathknightMinorGlyph.GlyphOfHornOfWinter,
			minor2: DeathknightMinorGlyph.GlyphOfBloodTap,
			minor3: DeathknightMinorGlyph.GlyphOfRaiseDead,
		}),
	}),
};

export const FrostTalents = {
	name: 'Frost',
	data: SavedTalents.create({
		talentsString: '005510003-3050535000223301030023310035-005',
		glyphs: Glyphs.create({
			major1: DeathknightMajorGlyph.GlyphOfDisease,
			major2: DeathknightMajorGlyph.GlyphOfRuneStrike,
			major3: DeathknightMajorGlyph.GlyphOfDarkCommand,
			minor1: DeathknightMinorGlyph.GlyphOfHornOfWinter,
			minor2: DeathknightMinorGlyph.GlyphOfBloodTap,
			minor3: DeathknightMinorGlyph.GlyphOfRaiseDead,
		}),
	}),
};

export const DoubleBuffFrostTalents = {
	name: '2B Frost',
	data: SavedTalents.create({
		talentsString: '00551005303003002-305053510022330100002301-005',
		glyphs: Glyphs.create({
			major1: DeathknightMajorGlyph.GlyphOfDisease,
			major2: DeathknightMajorGlyph.GlyphOfRuneStrike,
			major3: DeathknightMajorGlyph.GlyphOfDarkCommand,
			minor1: DeathknightMinorGlyph.GlyphOfHornOfWinter,
			minor2: DeathknightMinorGlyph.GlyphOfBloodTap,
			minor3: DeathknightMinorGlyph.GlyphOfRaiseDead,
		}),
	}),
};

export const DefaultRotation = TankDeathKnightRotation.create({
	opener: Opener.Threat,
	optimizationSetting: OptimizationSetting.Hps,
	bloodSpell: BloodSpell.BloodStrike,
	presence: Presence.Frost,
});

export const DefaultOptions = TankDeathKnightOptions.create({
	startingRunicPower: 0,
});

export const DefaultConsumes = Consumes.create({
	flask: Flask.FlaskOfStoneblood,
	food: Food.FoodDragonfinFilet,
	defaultPotion: Potions.IndestructiblePotion,
	prepopPotion: Potions.IndestructiblePotion,
});

export const BLOOD_LEGACY_PRESET_LEGACY_DEFAULT = {
	name: 'Blood Legacy',
	enableWhen: (player: Player<Spec.SpecTankDeathknight>) => player.getTalentTree() == 0,
	rotation: SavedRotation.create({
		specRotationOptionsJson: TankDeathKnightRotation.toJsonString(DefaultRotation),
	}),
}

export const BLOOD_IT_SPAM_ROTATION_PRESET_DEFAULT = {
	name: 'Blood Icy Touch APL',
	enableWhen: (player: Player<Spec.SpecDeathknight>) => player.getTalentTree() == 0,
	rotation: SavedRotation.create({
		specRotationOptionsJson: TankDeathKnightRotation.toJsonString(DefaultRotation),
		rotation: APLRotation.fromJsonString(`{
			"enabled": true,
			"prepullActions": [
			  {"action":{"castSpell":{"spellId":{"spellId":48263}}},"doAtValue":{"const":{"val":"-10s"}}},
			  {"action":{"castSpell":{"spellId":{"spellId":42650}}},"doAtValue":{"const":{"val":"-6s"}}}
			],
			"priorityList": [
			  {"action":{"autocastOtherCooldowns":{}}},
			  {"action":{"condition":{"cmp":{"op":"OpLe","lhs":{"currentHealthPercent":{}},"rhs":{"const":{"val":"40%"}}}},"castSpell":{"spellId":{"spellId":48792}}}},
			  {"action":{"condition":{"cmp":{"op":"OpLe","lhs":{"currentHealthPercent":{}},"rhs":{"const":{"val":"40%"}}}},"castSpell":{"spellId":{"spellId":55233}}}},
			  {"action":{"condition":{"cmp":{"op":"OpLe","lhs":{"currentHealthPercent":{}},"rhs":{"const":{"val":"60%"}}}},"castSpell":{"spellId":{"spellId":48982}}}},
			  {"action":{"condition":{"cmp":{"op":"OpLe","lhs":{"currentHealthPercent":{}},"rhs":{"const":{"val":"60%"}}}},"castSpell":{"spellId":{"spellId":48707}}}},
			  {"action":{"condition":{"cmp":{"op":"OpLe","lhs":{"currentHealthPercent":{}},"rhs":{"const":{"val":"60%"}}}},"castSpell":{"spellId":{"spellId":48743}}}},
			  {"action":{"condition":{"cmp":{"op":"OpGe","lhs":{"currentRunicPower":{}},"rhs":{"const":{"val":"40"}}}},"castSpell":{"spellId":{"spellId":56815}}}},
			  {"action":{"condition":{"not":{"val":{"dotIsActive":{"spellId":{"spellId":55095}}}}},"castSpell":{"spellId":{"spellId":59131}}}},
			  {"action":{"condition":{"not":{"val":{"dotIsActive":{"spellId":{"spellId":55078}}}}},"castSpell":{"spellId":{"tag":1,"spellId":49921}}}},
			  {"action":{"condition":{"cmp":{"op":"OpLe","lhs":{"dotRemainingTime":{"spellId":{"spellId":55078}}},"rhs":{"const":{"val":"3s"}}}},"castSpell":{"spellId":{"spellId":50842}}}},
			  {"action":{"condition":{"and":{"vals":[{"cmp":{"op":"OpGt","lhs":{"currentNonDeathRuneCount":{"runeType":"RuneFrost"}},"rhs":{"const":{"val":"0"}}}},{"cmp":{"op":"OpGt","lhs":{"currentNonDeathRuneCount":{"runeType":"RuneUnholy"}},"rhs":{"const":{"val":"0"}}}}]}},"castSpell":{"spellId":{"tag":1,"spellId":49924}}}},
			  {"action":{"condition":{"cmp":{"op":"OpGt","lhs":{"currentRuneCount":{"runeType":"RuneDeath"}},"rhs":{"const":{"val":"0"}}}},"castSpell":{"spellId":{"spellId":59131}}}},
			  {"action":{"condition":{"or":{"vals":[{"cmp":{"op":"OpGt","lhs":{"currentNonDeathRuneCount":{"runeType":"RuneBlood"}},"rhs":{"const":{"val":"1"}}}},{"spellIsReady":{"spellId":{"spellId":47568}}}]}},"castSpell":{"spellId":{"tag":1,"spellId":49930}}}},
			  {"action":{"castSpell":{"spellId":{"spellId":46584}}}},
			  {"action":{"castSpell":{"spellId":{"spellId":47568}}}},
			  {"action":{"condition":{"cmp":{"op":"OpGe","lhs":{"currentRunicPower":{}},"rhs":{"const":{"val":"80"}}}},"castSpell":{"spellId":{"spellId":49895}}}}
			]
		}`),
	}),
}

export const BLOOD_AGGRO_ROTATION_PRESET_DEFAULT = {
	name: 'Blood Aggro APL',
	enableWhen: (player: Player<Spec.SpecDeathknight>) => player.getTalentTree() == 0,
	rotation: SavedRotation.create({
		specRotationOptionsJson: TankDeathKnightRotation.toJsonString(DefaultRotation),
		rotation: APLRotation.fromJsonString(`{
			"enabled": true,
			"prepullActions": [
			  {"action":{"castSpell":{"spellId":{"spellId":48263}}},"doAtValue":{"const":{"val":"-10s"}}},
			  {"action":{"castSpell":{"spellId":{"spellId":42650}}},"doAtValue":{"const":{"val":"-6s"}}}
			],
			"priorityList": [
			  {"action":{"autocastOtherCooldowns":{}}},
			  {"action":{"condition":{"and":{"vals":[{"cmp":{"op":"OpLe","lhs":{"currentHealthPercent":{}},"rhs":{"const":{"val":"40%"}}}},{"not":{"val":{"auraIsActive":{"auraId":{"spellId":55233}}}}}]}},"castSpell":{"spellId":{"spellId":48792}}}},
			  {"action":{"condition":{"and":{"vals":[{"cmp":{"op":"OpLe","lhs":{"currentHealthPercent":{}},"rhs":{"const":{"val":"40%"}}}},{"not":{"val":{"auraIsActive":{"auraId":{"spellId":48792}}}}}]}},"castSpell":{"spellId":{"spellId":55233}}}},
			  {"action":{"condition":{"cmp":{"op":"OpLe","lhs":{"currentHealthPercent":{}},"rhs":{"const":{"val":"60%"}}}},"castSpell":{"spellId":{"spellId":48707}}}},
			  {"action":{"condition":{"cmp":{"op":"OpLe","lhs":{"currentHealthPercent":{}},"rhs":{"const":{"val":"60%"}}}},"castSpell":{"spellId":{"spellId":48743}}}},
			  {"action":{"condition":{"or":{"vals":[{"not":{"val":{"spellIsReady":{"spellId":{"spellId":49028}}}}},{"cmp":{"op":"OpGe","lhs":{"currentRunicPower":{}},"rhs":{"const":{"val":"80"}}}}]}},"castSpell":{"spellId":{"spellId":56815}}}},
			  {"action":{"condition":{"not":{"val":{"dotIsActive":{"spellId":{"spellId":55095}}}}},"castSpell":{"spellId":{"spellId":59131}}}},
			  {"action":{"condition":{"not":{"val":{"dotIsActive":{"spellId":{"spellId":55078}}}}},"castSpell":{"spellId":{"tag":1,"spellId":49921}}}},
			  {"action":{"castSpell":{"spellId":{"spellId":49016}}}},
			  {"action":{"castSpell":{"spellId":{"spellId":49028}}}},
			  {"action":{"castSpell":{"spellId":{"tag":1,"spellId":55262}}}},
			  {"action":{"castSpell":{"spellId":{"tag":1,"spellId":49924}}}},
			  {"action":{"castSpell":{"spellId":{"spellId":46584}}}},
			  {"action":{"castSpell":{"spellId":{"spellId":47568}}}},
			  {"action":{"condition":{"cmp":{"op":"OpGe","lhs":{"currentRunicPower":{}},"rhs":{"const":{"val":"80"}}}},"castSpell":{"spellId":{"spellId":49895}}}}
			]
		}`),
	}),
}

export const P1_BLOOD_PRESET = {
	name: 'P1 Blood',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	gear: EquipmentSpec.fromJsonString(`{ "items": [
		{"id":40565,"enchant":3878,"gems":[41380,36767]},
		{"id":40387},
		{"id":39704,"enchant":3852,"gems":[40008]},
		{"id":40252,"enchant":3605},
		{"id":40559,"gems":[40008,40022]},
		{"id":40306,"enchant":3850,"gems":[40008,0]},
		{"id":40563,"enchant":3860,"gems":[40008,0]},
		{"id":39759,"gems":[40008,40008]},
		{"id":40567,"enchant":3822,"gems":[40008,40008]},
		{"id":40297,"enchant":3232},
		{"id":40718},
		{"id":40107},
		{"id":44063,"gems":[36767,36767]},
		{"id":42341,"gems":[40008,40008]},
		{"id":40406,"enchant":3847},
		{},
		{"id":40207}
  ]}`),
};

export const P2_BLOOD_PRESET = {
	name: 'P2 Blood',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	gear: EquipmentSpec.fromJsonString(`{ "items": [
		{"id":46120,"enchant":3878,"gems":[41380,36767]},
		  {"id":45485,"gems":[40008]},
		  {"id":46122,"enchant":3852,"gems":[40008]},
		  {"id":45496,"enchant":3605,"gems":[40022]},
		  {"id":46118,"gems":[36767,36767]},
		  {"id":45111,"enchant":3850,"gems":[0]},
		  {"id":46119,"enchant":3860,"gems":[40008,0]},
		  {"id":45551,"gems":[40008,40008,40008]},
		  {"id":45594,"enchant":3822,"gems":[40008,40008,40008]},
		  {"id":45988,"enchant":3232,"gems":[40008,40008]},
		  {"id":45471,"gems":[40008]},
		  {"id":45326},
		  {"id":45158},
		  {"id":46021},
		  {"id":45533,"enchant":3370,"gems":[40008,40008]},
		  {},
		  {"id":45144}
  ]}`),
};

export const P1_FROST_PRESET = {
	name: 'P1 Frost',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	gear: EquipmentSpec.fromJsonString(`{ "items": [
		{"id":40565,"enchant":3878,"gems":[41380,36767]},
		  {"id":40387},
		  {"id":40568,"enchant":3852,"gems":[40008]},
		  {"id":40252,"enchant":3605},
		  {"id":40559,"gems":[40008,40022]},
		  {"id":40306,"enchant":3850,"gems":[40008,0]},
		  {"id":40563,"enchant":3860,"gems":[40008,0]},
		  {"id":39759,"gems":[40008,40008]},
		  {"id":40589,"enchant":3822},
		  {"id":40297,"enchant":3232},
		  {"id":40718},
		  {"id":40107},
		  {"id":44063,"gems":[36767,36767]},
		  {"id":40257},
		  {"id":40345,"enchant":3370},
		  {"id":40345,"enchant":3368},
		  {"id":40714}
  ]}`),
};

export const P2_FROST_PRESET = {
	name: 'P2 Frost',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	gear: EquipmentSpec.fromJsonString(`{ "items": [
		{"id":46120,"enchant":3878,"gems":[41380,36767]},
		  {"id":45485,"gems":[40008]},
		  {"id":46122,"enchant":3852,"gems":[40008]},
		  {"id":45496,"enchant":3605,"gems":[40022]},
		  {"id":46118,"gems":[36767,36767]},
		  {"id":45111,"enchant":3850,"gems":[0]},
		  {"id":46119,"enchant":3860,"gems":[40008,0]},
		  {"id":45551,"gems":[40008,40008,40008]},
		  {"id":45594,"enchant":3822,"gems":[40008,40008,40008]},
		  {"id":45988,"enchant":3232,"gems":[40008,40008]},
		  {"id":45471,"gems":[40008]},
		  {"id":45326},
		  {"id":45158},
		  {"id":46021},
		  {"id":46097,"enchant":3370,"gems":[40008]},
		  {"id":46097,"enchant":3368,"gems":[40008]},
		  {"id":45144}
  ]}`),
};
