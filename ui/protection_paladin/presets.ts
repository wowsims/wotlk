import * as PresetUtils from '../core/preset_utils.js';
import {
	Consumes,
	Flask,
	Food,
	Potions,
} from '../core/proto/common.js';
import {
	PaladinAura as PaladinAura,
	PaladinJudgement as PaladinJudgement,
	PaladinMajorGlyph,
	PaladinMinorGlyph,
	ProtectionPaladin_Options as ProtectionPaladinOptions,
} from '../core/proto/paladin.js';
import { SavedTalents } from '../core/proto/ui.js';
// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.
import PreraidGear from './gear_sets/preraid.gear.json';
export const PRERAID_PRESET = PresetUtils.makePresetGear('PreRaid预设', PreraidGear);
import P4PreraidGear from './gear_sets/p4_preraid.gear.json';
export const P4_PRERAID_PRESET = PresetUtils.makePresetGear('P4预设', P4PreraidGear);
import P1Gear from './gear_sets/p1.gear.json';
export const P1_PRESET = PresetUtils.makePresetGear('P1预设', P1Gear);
import P2Gear from './gear_sets/p2.gear.json';
export const P2_PRESET = PresetUtils.makePresetGear('P2预设', P2Gear);
import P3Gear from './gear_sets/p3.gear.json';
export const P3_PRESET = PresetUtils.makePresetGear('P3预设', P3Gear);
import P4Gear from './gear_sets/p4.gear.json';
export const P4_PRESET = PresetUtils.makePresetGear('P4预设', P4Gear);

import DefaultApl from './apls/default.apl.json';
export const ROTATION_DEFAULT = PresetUtils.makePresetAPLRotation('标准预设', DefaultApl);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/wotlk/talent-calc and copy the numbers in the url.

export const GenericAoeTalents = {
	name: '标准预设',
	data: SavedTalents.create({
		talentsString: '-05005135200132311333312321-511302012003',
		glyphs: {
			major1: PaladinMajorGlyph.GlyphOfSealOfVengeance,
			major2: PaladinMajorGlyph.GlyphOfRighteousDefense,
			major3: PaladinMajorGlyph.GlyphOfDivinePlea,
			minor1: PaladinMinorGlyph.GlyphOfSenseUndead,
			minor2: PaladinMinorGlyph.GlyphOfLayOnHands,
			minor3: PaladinMinorGlyph.GlyphOfBlessingOfKings
		}
	}),
};

export const DefaultOptions = ProtectionPaladinOptions.create({
	aura: PaladinAura.RetributionAura,
	judgement: PaladinJudgement.JudgementOfWisdom,
});

export const DefaultConsumes = Consumes.create({
	flask: Flask.FlaskOfStoneblood,
	food: Food.FoodDragonfinFilet,
	defaultPotion: Potions.IndestructiblePotion,
	prepopPotion: Potions.IndestructiblePotion,
});
