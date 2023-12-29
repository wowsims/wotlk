import {
	Consumes,
	Flask,
	Food,
	Glyphs,
	Potions,
} from '../core/proto/common.js';
import { SavedTalents } from '../core/proto/ui.js';

import {
	TankDeathknight_Options as TankDeathKnightOptions,
	DeathknightMajorGlyph,
	DeathknightMinorGlyph,
} from '../core/proto/deathknight.js';

import * as PresetUtils from '../core/preset_utils.js';

import P1BloodGear from './gear_sets/p1_blood.gear.json';
export const P1_BLOOD_PRESET = PresetUtils.makePresetGear('P1 Blood', P1BloodGear);
import P2BloodGear from './gear_sets/p2_blood.gear.json';
export const P2_BLOOD_PRESET = PresetUtils.makePresetGear('P2 Blood', P2BloodGear);
import P3BloodGear from './gear_sets/p3_blood.gear.json';
export const P3_BLOOD_PRESET = PresetUtils.makePresetGear('P3 Blood', P3BloodGear);
import P4BloodGear from './gear_sets/p4_blood.gear.json';
export const P4_BLOOD_PRESET = PresetUtils.makePresetGear('P4 Blood', P4BloodGear);
import P1FrostGear from './gear_sets/p1_frost.gear.json';
export const P1_FROST_PRESET = PresetUtils.makePresetGear('P1 Frost', P1FrostGear);
import P2FrostGear from './gear_sets/p2_frost.gear.json';
export const P2_FROST_PRESET = PresetUtils.makePresetGear('P2 Frost', P2FrostGear);

import BloodIcyTouchApl from './apls/blood_icy_touch.apl.json';
export const BLOOD_IT_SPAM_ROTATION_PRESET_DEFAULT = PresetUtils.makePresetAPLRotation('Blood Icy Touch', BloodIcyTouchApl);
import BloodAggroApl from './apls/blood_aggro.apl.json';
export const BLOOD_AGGRO_ROTATION_PRESET_DEFAULT = PresetUtils.makePresetAPLRotation('Blood Aggro', BloodAggroApl);

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

export const DefaultOptions = TankDeathKnightOptions.create({
	startingRunicPower: 0,
});

export const DefaultConsumes = Consumes.create({
	flask: Flask.FlaskOfStoneblood,
	food: Food.FoodDragonfinFilet,
	defaultPotion: Potions.IndestructiblePotion,
	prepopPotion: Potions.IndestructiblePotion,
});
