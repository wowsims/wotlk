import * as PresetUtils from '../core/preset_utils.js';
import {
	BattleElixir,
	Consumes,
	Explosive,
	Food,
	Glyphs,
	GuardianElixir,
	Potions,
	Spec,
} from '../core/proto/common.js';
import { SavedTalents } from '../core/proto/ui.js';
import {
	ProtectionWarrior_Options as ProtectionWarriorOptions,
	ProtectionWarrior_Rotation as ProtectionWarriorRotation,
	WarriorMajorGlyph,
	WarriorMinorGlyph,
	WarriorShout,
} from '../core/proto/warrior.js';
// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.
import PreraidBalancedGear from './gear_sets/preraid_balanced.gear.json';
export const PRERAID_BALANCED_PRESET = PresetUtils.makePresetGear('Preraid预设', PreraidBalancedGear);
import PreraidP4Gear from './gear_sets/p4_preraid.gear.json';
export const P4_PRERAID_PRESET = PresetUtils.makePresetGear('P4前期预设', PreraidP4Gear);
import P1BalancedGear from './gear_sets/p1_balanced.gear.json';
export const P1_BALANCED_PRESET = PresetUtils.makePresetGear('P1预设', P1BalancedGear);
import P2SurvivalGear from './gear_sets/p2_survival.gear.json';
export const P2_SURVIVAL_PRESET = PresetUtils.makePresetGear('P2预设', P2SurvivalGear);
import P3Gear from './gear_sets/p3.gear.json';
export const P3_PRESET = PresetUtils.makePresetGear('P3预设', P3Gear);
import P4Gear from './gear_sets/p4.gear.json';
export const P4_PRESET = PresetUtils.makePresetGear('P4预设', P4Gear);

import DefaultApl from './apls/default.apl.json';
export const ROTATION_DEFAULT = PresetUtils.makePresetAPLRotation('标准APL', DefaultApl);
export const ROTATION_PRESET_SIMPLE = PresetUtils.makePresetSimpleRotation('标准预设', Spec.SpecProtectionWarrior, ProtectionWarriorRotation.create());

// 默认天赋。使用 wowhead 计算器格式，在 https://wowhead.com/wotlk/talent-calc 上创建天赋并复制 URL 中的数字。
export const StandardTalents = {
	name: '标准预设',
	data: SavedTalents.create({
		talentsString: '2500030023-302-053351225000012521030113321',
		glyphs: Glyphs.create({
			major1: WarriorMajorGlyph.GlyphOfBlocking,
			major2: WarriorMajorGlyph.GlyphOfVigilance,
			major3: WarriorMajorGlyph.GlyphOfDevastate,
			minor1: WarriorMinorGlyph.GlyphOfCharge,
			minor2: WarriorMinorGlyph.GlyphOfThunderClap,
			minor3: WarriorMinorGlyph.GlyphOfCommand,
		}),
	}),
};

export const UATalents = {
	name: '复仇',
	data: SavedTalents.create({
		talentsString: '35023301230051002020120002-2-05035122500000252',
		glyphs: Glyphs.create({
			major1: WarriorMajorGlyph.GlyphOfRevenge,
			major2: WarriorMajorGlyph.GlyphOfHeroicStrike,
			major3: WarriorMajorGlyph.GlyphOfSweepingStrikes,
			minor1: WarriorMinorGlyph.GlyphOfCharge,
			minor2: WarriorMinorGlyph.GlyphOfThunderClap,
			minor3: WarriorMinorGlyph.GlyphOfCommand,
		}),
	}),
};

export const DefaultOptions = ProtectionWarriorOptions.create({
	shout: WarriorShout.WarriorShoutCommanding,
	useShatteringThrow: false,
	startingRage: 0,
});

export const DefaultConsumes = Consumes.create({
	battleElixir: BattleElixir.ElixirOfExpertise,
	guardianElixir: GuardianElixir.ElixirOfProtection,
	food: Food.FoodDragonfinFilet,
	defaultPotion: Potions.IndestructiblePotion,
	prepopPotion: Potions.IndestructiblePotion,
	thermalSapper: true,
	fillerExplosive: Explosive.ExplosiveSaroniteBomb,
});
