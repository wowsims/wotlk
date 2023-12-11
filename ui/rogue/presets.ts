import {
	Consumes,
	Flask,
	Food,
	Glyphs
} from '../core/proto/common.js';
import { SavedTalents } from '../core/proto/ui.js';

import {
	Rogue_Options_PoisonImbue as Poison,
	RogueMajorGlyph,
	Rogue_Options as RogueOptions,
	Rogue_Rotation as RogueRotation,
	Rogue_Rotation_AssassinationPriority,
	Rogue_Rotation_CombatBuilder,
	Rogue_Rotation_CombatPriority,
	Rogue_Rotation_Frequency,
	Rogue_Rotation_SubtletyBuilder,
	Rogue_Rotation_SubtletyPriority,
} from '../core/proto/rogue.js';

import * as PresetUtils from '../core/preset_utils.js';

import BlankGear from './gear_sets/blank.gear.json';

import CombatApl from './apls/combat.apl.json';
import CombatCleaveSndApl from './apls/combat_cleave_snd.apl.json';
import CombatCleaveSndExposeApl from './apls/combat_cleave_snd_expose.apl.json';
import CombatExposeApl from './apls/combat_expose.apl.json';
import FanAoeApl from './apls/fan_aoe.apl.json';
import MutilateApl from './apls/mutilate.apl.json';
import MutilateExposeApl from './apls/mutilate_expose.apl.json';
import RuptureMutilateApl from './apls/rupture_mutilate.apl.json';
import RuptureMutilateExposeApl from './apls/rupture_mutilate_expose.apl.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const GearAssassinationDefault = PresetUtils.makePresetGear('Blank', BlankGear, { talentTree: 0 });
export const GearCombatDefault = PresetUtils.makePresetGear('Blank', BlankGear, { talentTree: 0 });
export const GearSubtletyDefault = PresetUtils.makePresetGear('Blank', BlankGear, { talentTree: 0 });

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

export const ROTATION_PRESET_MUTILATE = PresetUtils.makePresetAPLRotation('Mutilate', MutilateApl, { talentTree: 0 });
export const ROTATION_PRESET_RUPTURE_MUTILATE = PresetUtils.makePresetAPLRotation('Rupture Mutilate', RuptureMutilateApl, { talentTree: 0 });
export const ROTATION_PRESET_MUTILATE_EXPOSE = PresetUtils.makePresetAPLRotation('Mutilate w/ Expose', MutilateExposeApl, { talentTree: 0 });
export const ROTATION_PRESET_RUPTURE_MUTILATE_EXPOSE = PresetUtils.makePresetAPLRotation('Rupture Mutilate w/ Expose', RuptureMutilateExposeApl, { talentTree: 0 });
export const ROTATION_PRESET_COMBAT = PresetUtils.makePresetAPLRotation('Combat', CombatApl, { talentTree: 1 });
export const ROTATION_PRESET_COMBAT_EXPOSE = PresetUtils.makePresetAPLRotation('Combat w/ Expose', CombatExposeApl, { talentTree: 1 });
export const ROTATION_PRESET_COMBAT_CLEAVE_SND = PresetUtils.makePresetAPLRotation('Combat Cleave SND', CombatCleaveSndApl, { talentTree: 1 });
export const ROTATION_PRESET_COMBAT_CLEAVE_SND_EXPOSE = PresetUtils.makePresetAPLRotation('Combat Cleave SND w/ Expose', CombatCleaveSndExposeApl, { talentTree: 1 });
export const ROTATION_PRESET_AOE = PresetUtils.makePresetAPLRotation('Fan AOE', FanAoeApl);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/classic/talent-calc and copy the numbers in the url.
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
	flask: Flask.FlaskUnknown,
	food: Food.FoodUnknown,
});
