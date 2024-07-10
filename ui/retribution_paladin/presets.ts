import * as PresetUtils from '../core/preset_utils.js';
import {
	Conjured,
	Consumes,
	Flask,
	Food,
	Glyphs,
	Potions,
} from '../core/proto/common.js';
import {
	PaladinAura as PaladinAura,
	PaladinJudgement as PaladinJudgement,
	PaladinMajorGlyph,
	PaladinMinorGlyph,
	RetributionPaladin_Options as RetributionPaladinOptions,
} from '../core/proto/paladin.js';
import { SavedTalents } from '../core/proto/ui.js';
// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.
import PreraidGear from './gear_sets/preraid.gear.json';
export const PRERAID_PRESET = PresetUtils.makePresetGear('PreRaid预设', PreraidGear);
import P1Gear from './gear_sets/p1.gear.json';
export const P1_PRESET = PresetUtils.makePresetGear('P1预设', P1Gear);
import P2Gear from './gear_sets/p2.gear.json';
export const P2_PRESET = PresetUtils.makePresetGear('P2预设', P2Gear);
import P3MaceGear from './gear_sets/p3_mace.gear.json';
export const P3_PRESET = PresetUtils.makePresetGear('P3预设', P3MaceGear);
import P4Gear from './gear_sets/p4.gear.json';
export const P4_PRESET = PresetUtils.makePresetGear('P4预设', P4Gear);
import P5Gear from './gear_sets/p5.gear.json';
export const P5_PRESET = PresetUtils.makePresetGear('P5预设', P5Gear);

import DefaultApl from './apls/default.apl.json';
export const ROTATION_PRESET_DEFAULT = PresetUtils.makePresetAPLRotation('标准预设', DefaultApl);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/wotlk/talent-calc and copy the numbers in the url.
export const AuraMasteryTalents = {
	name: '光环掌握',
	data: SavedTalents.create({
		talentsString: '050501-05-05232051203331302133231331',
		glyphs: Glyphs.create({
			major1: PaladinMajorGlyph.GlyphOfSealOfVengeance,
			major2: PaladinMajorGlyph.GlyphOfJudgement,
			major3: PaladinMajorGlyph.GlyphOfReckoning,
			minor1: PaladinMinorGlyph.GlyphOfSenseUndead,
			minor2: PaladinMinorGlyph.GlyphOfLayOnHands,
			minor3: PaladinMinorGlyph.GlyphOfBlessingOfKings
		})
	}),
};

export const DivineSacTalents = {
	name: '神圣牺牲和守护',
	data: SavedTalents.create({
		talentsString: '03-453201002-05222051203331302133201331',
		glyphs: Glyphs.create({
			major1: PaladinMajorGlyph.GlyphOfSealOfVengeance,
			major2: PaladinMajorGlyph.GlyphOfJudgement,
			major3: PaladinMajorGlyph.GlyphOfReckoning,
			minor1: PaladinMinorGlyph.GlyphOfSenseUndead,
			minor2: PaladinMinorGlyph.GlyphOfLayOnHands,
			minor3: PaladinMinorGlyph.GlyphOfBlessingOfKings
		})
	}),
};

export const DefaultOptions = RetributionPaladinOptions.create({
	aura: PaladinAura.RetributionAura,
	judgement: PaladinJudgement.JudgementOfWisdom,
});

export const DefaultConsumes = Consumes.create({
	defaultPotion: Potions.PotionOfSpeed,
	defaultConjured: Conjured.ConjuredDarkRune,
	flask: Flask.FlaskOfEndlessRage,
	food: Food.FoodDragonfinFilet,
});
