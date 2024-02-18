import {
	Conjured,
	Consumes,
	Flask,
	Food,
	Glyphs,
	Potions,
} from '../core/proto/common.js';
import { SavedTalents } from '../core/proto/ui.js';

import {
	Rogue_Options as RogueOptions,
	Rogue_Options_PoisonImbue as Poison,
	RogueMajorGlyph,
} from '../core/proto/rogue.js';

import * as PresetUtils from '../core/preset_utils.js';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

import PreraidAssassinationGear from './gear_sets/preraid_assassination.gear.json';
export const PRERAID_PRESET_ASSASSINATION = PresetUtils.makePresetGear('PreRaid Assassination', PreraidAssassinationGear, { talentTree: 0 });
import P1AssassinationGear from './gear_sets/p1_assassination.gear.json';
export const P1_PRESET_ASSASSINATION = PresetUtils.makePresetGear('P1 Assassination', P1AssassinationGear, { talentTree: 0 });
import P2AssassinationGear from './gear_sets/p2_assassination.gear.json';
export const P2_PRESET_ASSASSINATION = PresetUtils.makePresetGear('P2 Assassination', P2AssassinationGear, { talentTree: 0 });
import P3AssassinationGear from './gear_sets/p3_assassination.gear.json';
export const P3_PRESET_ASSASSINATION = PresetUtils.makePresetGear('P3 Assassination', P3AssassinationGear, { talentTree: 0 });
import P4AssassinationGear from './gear_sets/p4_assassination.gear.json';
export const P4_PRESET_ASSASSINATION = PresetUtils.makePresetGear('P4 Assassination', P4AssassinationGear, { talentTree: 0 });
import P5AssassinationGear from './gear_sets/p5_assassination.gear.json';
export const P5_PRESET_ASSASSINATION = PresetUtils.makePresetGear('P5 Assassination', P5AssassinationGear, { talentTree: 0 });
import PreraidCombatGear from './gear_sets/preraid_combat.gear.json';
export const PRERAID_PRESET_COMBAT = PresetUtils.makePresetGear('PreRaid Combat', PreraidCombatGear, { talentTree: 1 });
import P1CombatGear from './gear_sets/p1_combat.gear.json';
export const P1_PRESET_COMBAT = PresetUtils.makePresetGear('P1 Combat', P1CombatGear, { talentTree: 1 });
import P2CombatGear from './gear_sets/p2_combat.gear.json';
export const P2_PRESET_COMBAT = PresetUtils.makePresetGear('P2 Combat', P2CombatGear, { talentTree: 1 });
import P3CombatGear from './gear_sets/p3_combat.gear.json';
export const P3_PRESET_COMBAT = PresetUtils.makePresetGear('P3 Combat', P3CombatGear, { talentTree: 1 });
import P4CombatGear from './gear_sets/p4_combat.gear.json';
export const P4_PRESET_COMBAT = PresetUtils.makePresetGear('P4 Combat', P4CombatGear, { talentTree: 1 });
import P5CombatGear from './gear_sets/p5_combat.gear.json';
export const P5_PRESET_COMBAT = PresetUtils.makePresetGear('P5 Combat', P5CombatGear, { talentTree: 1 });
import P1HemoSubGear from './gear_sets/p1_hemosub.gear.json';
export const P1_PRESET_HEMO_SUB = PresetUtils.makePresetGear('P1 Hemo Sub', P1HemoSubGear, { talentTree: 2 });
import P2HemoSubGear from './gear_sets/p2_hemosub.gear.json';
export const P2_PRESET_HEMO_SUB = PresetUtils.makePresetGear('P2 Hemo Sub', P2HemoSubGear, { talentTree: 2 });
import P3HemoSubGear from './gear_sets/p3_hemosub.gear.json';
export const P3_PRESET_HEMO_SUB = PresetUtils.makePresetGear('P3 Hemo Sub', P3HemoSubGear, { talentTree: 2 });
import P3DanceSubGear from './gear_sets/p3_dancesub.gear.json';
export const P3_PRESET_DANCE_SUB = PresetUtils.makePresetGear('P3 Dance Sub', P3DanceSubGear, { talentTree: 2 });

import MutilateApl from './apls/mutilate.apl.json'
export const ROTATION_PRESET_MUTILATE = PresetUtils.makePresetAPLRotation('Mutilate', MutilateApl, { talentTree: 0 });
import MutilateExposeApl from './apls/mutilate_expose.apl.json'
export const ROTATION_PRESET_RUPTURE_MUTILATE = PresetUtils.makePresetAPLRotation('Rupture Mutilate', RuptureMutilateApl, { talentTree: 0 });
import RuptureMutilateApl from './apls/rupture_mutilate.apl.json'
export const ROTATION_PRESET_MUTILATE_EXPOSE = PresetUtils.makePresetAPLRotation('Mutilate w/ Expose', MutilateExposeApl, { talentTree: 0 });
import RuptureMutilateExposeApl from './apls/rupture_mutilate_expose.apl.json'
export const ROTATION_PRESET_RUPTURE_MUTILATE_EXPOSE = PresetUtils.makePresetAPLRotation('Rupture Mutilate w/ Expose', RuptureMutilateExposeApl, { talentTree: 0 });
import CombatApl from './apls/combat.apl.json'
export const ROTATION_PRESET_COMBAT = PresetUtils.makePresetAPLRotation('Combat', CombatApl, { talentTree: 1 });
import CombatExposeApl from './apls/combat_expose.apl.json'
export const ROTATION_PRESET_COMBAT_EXPOSE = PresetUtils.makePresetAPLRotation('Combat w/ Expose', CombatExposeApl, { talentTree: 1 });
import CombatCleaveSndApl from './apls/combat_cleave_snd.apl.json'
export const ROTATION_PRESET_COMBAT_CLEAVE_SND = PresetUtils.makePresetAPLRotation('Combat Cleave SND', CombatCleaveSndApl, { talentTree: 1 });
import CombatCleaveSndExposeApl from './apls/combat_cleave_snd_expose.apl.json'
export const ROTATION_PRESET_COMBAT_CLEAVE_SND_EXPOSE = PresetUtils.makePresetAPLRotation('Combat Cleave SND w/ Expose', CombatCleaveSndExposeApl, { talentTree: 1 });
import FanAoeApl from './apls/fan_aoe.apl.json'
export const ROTATION_PRESET_AOE = PresetUtils.makePresetAPLRotation('Fan AOE', FanAoeApl);

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
