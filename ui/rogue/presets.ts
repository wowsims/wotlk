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
	Rogue_Options as RogueOptions,
	Rogue_Options_PoisonImbue as Poison,
	RogueMajorGlyph,
} from '../core/proto/rogue.js';
import { SavedTalents } from '../core/proto/ui.js';
// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.
import PreraidAssassinationGear from './gear_sets/preraid_assassination.gear.json';
export const PRERAID_PRESET_ASSASSINATION = PresetUtils.makePresetGear('PreRaid刺杀', PreraidAssassinationGear, { talentTree: 0 });
import P1AssassinationGear from './gear_sets/p1_assassination.gear.json';
export const P1_PRESET_ASSASSINATION = PresetUtils.makePresetGear('P1刺杀', P1AssassinationGear, { talentTree: 0 });
import P2AssassinationGear from './gear_sets/p2_assassination.gear.json';
export const P2_PRESET_ASSASSINATION = PresetUtils.makePresetGear('P2刺杀', P2AssassinationGear, { talentTree: 0 });
import P3AssassinationGear from './gear_sets/p3_assassination.gear.json';
export const P3_PRESET_ASSASSINATION = PresetUtils.makePresetGear('P3刺杀', P3AssassinationGear, { talentTree: 0 });
import P4AssassinationGear from './gear_sets/p4_assassination.gear.json';
export const P4_PRESET_ASSASSINATION = PresetUtils.makePresetGear('P4刺杀', P4AssassinationGear, { talentTree: 0 });
import P5AssassinationGear from './gear_sets/p5_assassination.gear.json';
export const P5_PRESET_ASSASSINATION = PresetUtils.makePresetGear('P5刺杀', P5AssassinationGear, { talentTree: 0 });
import PreraidCombatGear from './gear_sets/preraid_combat.gear.json';
export const PRERAID_PRESET_COMBAT = PresetUtils.makePresetGear('PreRaid战斗', PreraidCombatGear, { talentTree: 1 });
import P1CombatGear from './gear_sets/p1_combat.gear.json';
export const P1_PRESET_COMBAT = PresetUtils.makePresetGear('P1战斗', P1CombatGear, { talentTree: 1 });
import P2CombatGear from './gear_sets/p2_combat.gear.json';
export const P2_PRESET_COMBAT = PresetUtils.makePresetGear('P2战斗', P2CombatGear, { talentTree: 1 });
import P3CombatGear from './gear_sets/p3_combat.gear.json';
export const P3_PRESET_COMBAT = PresetUtils.makePresetGear('P3战斗', P3CombatGear, { talentTree: 1 });
import P4CombatGear from './gear_sets/p4_combat.gear.json';
export const P4_PRESET_COMBAT = PresetUtils.makePresetGear('P4战斗', P4CombatGear, { talentTree: 1 });
import P5CombatGear from './gear_sets/p5_combat.gear.json';
export const P5_PRESET_COMBAT = PresetUtils.makePresetGear('P5战斗', P5CombatGear, { talentTree: 1 });
import P1HemoSubGear from './gear_sets/p1_hemosub.gear.json';
export const P1_PRESET_HEMO_SUB = PresetUtils.makePresetGear('P1敏锐(出血)', P1HemoSubGear, { talentTree: 2 });
import P2HemoSubGear from './gear_sets/p2_hemosub.gear.json';
export const P2_PRESET_HEMO_SUB = PresetUtils.makePresetGear('P2敏锐(出血)', P2HemoSubGear, { talentTree: 2 });
import P3HemoSubGear from './gear_sets/p3_hemosub.gear.json';
export const P3_PRESET_HEMO_SUB = PresetUtils.makePresetGear('P3敏锐(出血)', P3HemoSubGear, { talentTree: 2 });
import P3DanceSubGear from './gear_sets/p3_dancesub.gear.json';
export const P3_PRESET_DANCE_SUB = PresetUtils.makePresetGear('P3敏锐(舞))', P3DanceSubGear, { talentTree: 2 });

import MutilateApl from './apls/mutilate.apl.json';
export const ROTATION_PRESET_MUTILATE = PresetUtils.makePresetAPLRotation('毁伤', MutilateApl, { talentTree: 0 });
import MutilateExposeApl from './apls/mutilate_expose.apl.json';
export const ROTATION_PRESET_MUTILATE_EXPOSE = PresetUtils.makePresetAPLRotation('毁伤+破甲', MutilateExposeApl, { talentTree: 0 });
import RuptureMutilateApl from './apls/rupture_mutilate.apl.json';
export const ROTATION_PRESET_RUPTURE_MUTILATE = PresetUtils.makePresetAPLRotation('割裂+破甲', RuptureMutilateApl, { talentTree: 0 });
import RuptureMutilateExposeApl from './apls/rupture_mutilate_expose.apl.json';
export const ROTATION_PRESET_RUPTURE_MUTILATE_EXPOSE = PresetUtils.makePresetAPLRotation('毁伤+割裂+破甲', RuptureMutilateExposeApl, { talentTree: 0 });
import CombatApl from './apls/combat.apl.json';
export const ROTATION_PRESET_COMBAT = PresetUtils.makePresetAPLRotation('战斗', CombatApl, { talentTree: 1 });
import CombatExposeApl from './apls/combat_expose.apl.json';
export const ROTATION_PRESET_COMBAT_EXPOSE = PresetUtils.makePresetAPLRotation('战斗+破甲', CombatExposeApl, { talentTree: 1 });
import CombatCleaveSndApl from './apls/combat_cleave_snd.apl.json';
export const ROTATION_PRESET_COMBAT_CLEAVE_SND = PresetUtils.makePresetAPLRotation('战斗+切割', CombatCleaveSndApl, { talentTree: 1 });
import CombatCleaveSndExposeApl from './apls/combat_cleave_snd_expose.apl.json';
export const ROTATION_PRESET_COMBAT_CLEAVE_SND_EXPOSE = PresetUtils.makePresetAPLRotation('战斗+破甲+切割', CombatCleaveSndExposeApl, { talentTree: 1 });

import FanAoeApl from './apls/fan_aoe.apl.json';
export const ROTATION_PRESET_AOE = PresetUtils.makePresetAPLRotation('刀扇AOE', FanAoeApl);


// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/wotlk/talent-calc and copy the numbers in the url.
export const CombatHackTalents = {
	name: '战斗(斧/剑)',
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
	name: '战斗(拳套)',
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
	name: '刺杀13/7',
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
	name: '刺杀18/2',
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
	name: '刺杀(乱舞)',
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
	name: '敏锐',
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
	name: '敏锐(出血)',
	data: SavedTalents.create({
		talentsString: '30532010135--502201203032112135011503122',
		glyphs: Glyphs.create({
			major1: RogueMajorGlyph.GlyphOfEviscerate,
			major2: RogueMajorGlyph.GlyphOfRupture,
			major3: RogueMajorGlyph.GlyphOfTricksOfTheTrade,
		})
	}),
}

export const DefaultOptions = RogueOptions.create({
	mhImbue: Poison.DeadlyPoison,
	ohImbue: Poison.InstantPoison,
	applyPoisonsManually: false,
	startingOverkillDuration: 20,
	vanishBreakTime: 0.1,
	honorOfThievesCritRate: 400,
});

export const DefaultConsumes = Consumes.create({
	defaultPotion: Potions.PotionOfSpeed,
	prepopPotion: Potions.PotionOfSpeed,
	defaultConjured: Conjured.ConjuredRogueThistleTea,
	flask: Flask.FlaskOfEndlessRage,
	food: Food.FoodMegaMammothMeal,
});
