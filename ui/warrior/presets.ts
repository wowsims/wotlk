import * as PresetUtils from '../core/preset_utils.js';
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
	Warrior_Options as WarriorOptions,
	WarriorMajorGlyph,
	WarriorMinorGlyph,
	WarriorShout,
} from '../core/proto/warrior.js';
// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.
import PreraidArmsGear from './gear_sets/preraid_arms.gear.json';
export const PRERAID_ARMS_PRESET = PresetUtils.makePresetGear('Preraid武器战', PreraidArmsGear, { talentTree: 0 });
import P1ArmsGear from './gear_sets/p1_arms.gear.json';
export const P1_ARMS_PRESET = PresetUtils.makePresetGear('P1武器', P1ArmsGear, { talentTree: 0 });
import P2ArmsGear from './gear_sets/p2_arms.gear.json';
export const P2_ARMS_PRESET = PresetUtils.makePresetGear('P2武器', P2ArmsGear, { talentTree: 0 });
import P3Arms2pAllianceGear from './gear_sets/p3_arms_2p_alliance.gear.json';
export const P3_ARMS_2P_PRESET_ALLIANCE = PresetUtils.makePresetGear('P3武器2T9[联盟]', P3Arms2pAllianceGear, { talentTree: 0, faction: Faction.Alliance });
import P3Arms4pAllianceGear from './gear_sets/p3_arms_4p_alliance.gear.json';
export const P3_ARMS_4P_PRESET_ALLIANCE = PresetUtils.makePresetGear('P3武器4T9[联盟]', P3Arms4pAllianceGear, { talentTree: 0, faction: Faction.Alliance });
import P3Arms2pHordeGear from './gear_sets/p3_arms_2p_horde.gear.json';
export const P3_ARMS_2P_PRESET_HORDE = PresetUtils.makePresetGear('P3武器2T9[部落]', P3Arms2pHordeGear, { talentTree: 0, faction: Faction.Horde });
import P3Arms4pHordeGear from './gear_sets/p3_arms_4p_horde.gear.json';
export const P3_ARMS_4P_PRESET_HORDE = PresetUtils.makePresetGear('P3武器4T9[部落]', P3Arms4pHordeGear, { talentTree: 0, faction: Faction.Horde });
import P4ArmsAllianceGear from './gear_sets/p4_arms_alliance.gear.json';
export const P4_ARMS_PRESET_ALLIANCE = PresetUtils.makePresetGear('P4武器[联盟]', P4ArmsAllianceGear, { talentTree: 0, faction: Faction.Alliance });
import P4ArmsHordeGear from './gear_sets/p4_arms_horde.gear.json';
export const P4_ARMS_PRESET_HORDE = PresetUtils.makePresetGear('P4武器[部落]', P4ArmsHordeGear, { talentTree: 0, faction: Faction.Horde });
import PreraidFuryGear from './gear_sets/preraid_fury.gear.json';
export const PRERAID_FURY_PRESET = PresetUtils.makePresetGear('Preraid狂暴', PreraidFuryGear, { talentTrees: [1,2] });
import P1FuryGear from './gear_sets/p1_fury.gear.json';
export const P1_FURY_PRESET = PresetUtils.makePresetGear('P1狂暴', P1FuryGear, { talentTrees: [1,2] });
import P2FuryGear from './gear_sets/p2_fury.gear.json';
export const P2_FURY_PRESET = PresetUtils.makePresetGear('P2狂暴', P2FuryGear, { talentTrees: [1,2] });
import P3FuryAllianceGear from './gear_sets/p3_fury_alliance.gear.json';
export const P3_FURY_PRESET_ALLIANCE = PresetUtils.makePresetGear('P3狂暴[联盟]', P3FuryAllianceGear, { talentTrees: [1,2], faction: Faction.Alliance });
import P3FuryHordeGear from './gear_sets/p3_fury_horde.gear.json';
export const P3_FURY_PRESET_HORDE = PresetUtils.makePresetGear('P3狂暴[部落]', P3FuryHordeGear, { talentTrees: [1,2], faction: Faction.Horde });
import P4FuryAllianceGear from './gear_sets/p4_fury_alliance.gear.json';
export const P4_FURY_PRESET_ALLIANCE = PresetUtils.makePresetGear('P4狂暴[联盟]', P4FuryAllianceGear, { talentTrees: [1,2], faction: Faction.Alliance });
import P4FuryHordeGear from './gear_sets/p4_fury_horde.gear.json';
export const P4_FURY_PRESET_HORDE = PresetUtils.makePresetGear('P4狂暴[部落]', P4FuryHordeGear, { talentTrees: [1,2], faction: Faction.Horde });

import FuryApl from './apls/fury.apl.json';
export const ROTATION_FURY = PresetUtils.makePresetAPLRotation('狂暴', FuryApl, { talentTree: 1 });
import FurySunderApl from './apls/fury_sunder.apl.json';
export const ROTATION_FURY_SUNDER = PresetUtils.makePresetAPLRotation('狂暴+破甲', FurySunderApl, { talentTree: 1 });
import ArmsApl from './apls/arms.apl.json';
export const ROTATION_ARMS = PresetUtils.makePresetAPLRotation('武器', ArmsApl, { talentTree: 0 });
import ArmsSunderApl from './apls/arms_sunder.apl.json';
export const ROTATION_ARMS_SUNDER = PresetUtils.makePresetAPLRotation('武器+破甲', ArmsSunderApl, { talentTree: 0 });

// 默认天赋。使用wowhead计算器格式，在https://wowhead.com/wotlk/talent-calc上创建天赋并复制URL中的数字。
export const ArmsTalents = {
	name: '武器',
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
	name: '狂暴',
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
