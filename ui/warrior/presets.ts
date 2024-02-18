import {
	Consumes,
	Faction,
	Flask,
	Food,
	Glyphs,
	Potions,
} from '../core/proto/common.js';
import { SavedTalents } from '../core/proto/ui.js';

import {
	WarriorShout,
	Warrior_Options as WarriorOptions,
	WarriorMajorGlyph,
	WarriorMinorGlyph,
} from '../core/proto/warrior.js';

import * as PresetUtils from '../core/preset_utils.js';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

import PreraidArmsGear from './gear_sets/preraid_arms.gear.json';
export const PRERAID_ARMS_PRESET = PresetUtils.makePresetGear('Preraid Arms', PreraidArmsGear, { talentTree: 0 });
import P1ArmsGear from './gear_sets/p1_arms.gear.json';
export const P1_ARMS_PRESET = PresetUtils.makePresetGear('P1 Arms', P1ArmsGear, { talentTree: 0 });
import P2ArmsGear from './gear_sets/p2_arms.gear.json';
export const P2_ARMS_PRESET = PresetUtils.makePresetGear('P2 Arms', P2ArmsGear, { talentTree: 0 });
import P3Arms2pAllianceGear from './gear_sets/p3_arms_2p_alliance.gear.json';
export const P3_ARMS_2P_PRESET_ALLIANCE = PresetUtils.makePresetGear('P3 Arms 2p [A]', P3Arms2pAllianceGear, { talentTree: 0, faction: Faction.Alliance });
import P3Arms4pAllianceGear from './gear_sets/p3_arms_4p_alliance.gear.json';
export const P3_ARMS_4P_PRESET_ALLIANCE = PresetUtils.makePresetGear('P3 Arms 4p [A]', P3Arms4pAllianceGear, { talentTree: 0, faction: Faction.Alliance });
import P3Arms2pHordeGear from './gear_sets/p3_arms_2p_horde.gear.json';
export const P3_ARMS_2P_PRESET_HORDE = PresetUtils.makePresetGear('P3 Arms 2p [H]', P3Arms2pHordeGear, { talentTree: 0, faction: Faction.Horde });
import P3Arms4pHordeGear from './gear_sets/p3_arms_4p_horde.gear.json';
export const P3_ARMS_4P_PRESET_HORDE = PresetUtils.makePresetGear('P3 Arms 4p [H]', P3Arms4pHordeGear, { talentTree: 0, faction: Faction.Horde });
import P4ArmsAllianceGear from './gear_sets/p4_arms_alliance.gear.json';
export const P4_ARMS_PRESET_ALLIANCE = PresetUtils.makePresetGear('P4 Arms [A]', P4ArmsAllianceGear, { talentTree: 0, faction: Faction.Alliance });
import P4ArmsHordeGear from './gear_sets/p4_arms_horde.gear.json';
export const P4_ARMS_PRESET_HORDE = PresetUtils.makePresetGear('P4 Arms [H]', P4ArmsHordeGear, { talentTree: 0, faction: Faction.Horde });
import PreraidFuryGear from './gear_sets/preraid_fury.gear.json';
export const PRERAID_FURY_PRESET = PresetUtils.makePresetGear('Preraid Fury', PreraidFuryGear, { talentTrees: [1,2] });
import P1FuryGear from './gear_sets/p1_fury.gear.json';
export const P1_FURY_PRESET = PresetUtils.makePresetGear('P1 Fury', P1FuryGear, { talentTrees: [1,2] });
import P2FuryGear from './gear_sets/p2_fury.gear.json';
export const P2_FURY_PRESET = PresetUtils.makePresetGear('P2 Fury', P2FuryGear, { talentTrees: [1,2] });
import P3FuryAllianceGear from './gear_sets/p3_fury_alliance.gear.json';
export const P3_FURY_PRESET_ALLIANCE = PresetUtils.makePresetGear('P3 Fury [A]', P3FuryAllianceGear, { talentTrees: [1,2], faction: Faction.Alliance });
import P3FuryHordeGear from './gear_sets/p3_fury_horde.gear.json';
export const P3_FURY_PRESET_HORDE = PresetUtils.makePresetGear('P3 Fury [H]', P3FuryHordeGear, { talentTrees: [1,2], faction: Faction.Horde });
import P4FuryAllianceGear from './gear_sets/p4_fury_alliance.gear.json';
export const P4_FURY_PRESET_ALLIANCE = PresetUtils.makePresetGear('P4 Fury [A]', P4FuryAllianceGear, { talentTrees: [1,2], faction: Faction.Alliance });
import P4FuryHordeGear from './gear_sets/p4_fury_horde.gear.json';
export const P4_FURY_PRESET_HORDE = PresetUtils.makePresetGear('P4 Fury [H]', P4FuryHordeGear, { talentTrees: [1,2], faction: Faction.Horde });

import FuryApl from './apls/fury.apl.json';
export const ROTATION_FURY = PresetUtils.makePresetAPLRotation('Fury', FuryApl, { talentTree: 1 });
import FurySunderApl from './apls/fury_sunder.apl.json';
export const ROTATION_FURY_SUNDER = PresetUtils.makePresetAPLRotation('Fury + Sunder', FurySunderApl, { talentTree: 1 });
import ArmsApl from './apls/arms.apl.json';
export const ROTATION_ARMS = PresetUtils.makePresetAPLRotation('Arms', ArmsApl, { talentTree: 0 });
import ArmsSunderApl from './apls/arms_sunder.apl.json';
export const ROTATION_ARMS_SUNDER = PresetUtils.makePresetAPLRotation('Arms + Sunder', ArmsSunderApl, { talentTree: 0 });

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/wotlk/talent-calc and copy the numbers in the url.
export const ArmsTalents = {
	name: 'Arms',
	data: SavedTalents.create({
		talentsString: '3022032023335100102012213231251-305-2033',
		glyphs: Glyphs.create({
			major1: WarriorMajorGlyph.GlyphOfRending,
			major2: WarriorMajorGlyph.GlyphOfMortalStrike,
			major3: WarriorMajorGlyph.GlyphOfExecution,
			minor1: WarriorMinorGlyph.GlyphOfThunderClap,
			minor2: WarriorMinorGlyph.GlyphOfCommand,
			minor3: WarriorMinorGlyph.GlyphOfShatteringThrow,
		}),
	}),
};

export const FuryTalents = {
	name: 'Fury',
	data: SavedTalents.create({
		talentsString: '32002301233-305053000520310053120500351',
		glyphs: Glyphs.create({
			major1: WarriorMajorGlyph.GlyphOfWhirlwind,
			major2: WarriorMajorGlyph.GlyphOfHeroicStrike,
			major3: WarriorMajorGlyph.GlyphOfExecution,
			minor1: WarriorMinorGlyph.GlyphOfCommand,
			minor2: WarriorMinorGlyph.GlyphOfShatteringThrow,
			minor3: WarriorMinorGlyph.GlyphOfCharge,
		}),
	}),
};

export const DefaultOptions = WarriorOptions.create({
	startingRage: 0,
	useRecklessness: true,
	useShatteringThrow: true,
	disableExpertiseGemming: false,
	shout: WarriorShout.WarriorShoutCommanding,
});

export const DefaultConsumes = Consumes.create({
	flask: Flask.FlaskOfEndlessRage,
	food: Food.FoodSpicedWormBurger,
	defaultPotion: Potions.IndestructiblePotion,
	prepopPotion: Potions.PotionOfSpeed,
});