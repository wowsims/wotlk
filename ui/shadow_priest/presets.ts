import * as PresetUtils from '../core/preset_utils.js';
import {
	Consumes,
	Debuffs,
	Flask,
	Food,
	Glyphs,
	IndividualBuffs,
	Potions,
	Profession,
	RaidBuffs,
	TristateEffect,
} from '../core/proto/common.js';
import {
	PriestMajorGlyph as MajorGlyph,
	PriestMinorGlyph as MinorGlyph,
	ShadowPriest_Options as Options,
	ShadowPriest_Options_Armor as Armor,
} from '../core/proto/priest.js';
import { SavedTalents } from '../core/proto/ui.js';
// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.
import PreraidGear from './gear_sets/preraid.gear.json';
export const PRERAID_PRESET = PresetUtils.makePresetGear('Preraid预设', PreraidGear);
import P1Gear from './gear_sets/p1.gear.json';
export const P1_PRESET = PresetUtils.makePresetGear('P1预设', P1Gear);
import P2Gear from './gear_sets/p2.gear.json';
export const P2_PRESET = PresetUtils.makePresetGear('P2预设', P2Gear);
import P3Gear from './gear_sets/p3.gear.json';
export const P3_PRESET = PresetUtils.makePresetGear('P3预设', P3Gear);
import P4Gear from './gear_sets/p4.gear.json';
export const P4_PRESET = PresetUtils.makePresetGear('P4预设', P4Gear);

import DefaultApl from './apls/default.apl.json'
export const ROTATION_PRESET_DEFAULT = PresetUtils.makePresetAPLRotation('标准预设', DefaultApl);
import AOE24Apl from './apls/aoe_2_4.apl.json'
export const ROTATION_PRESET_AOE24 = PresetUtils.makePresetAPLRotation('AOE(2-4目标)', AOE24Apl);
import AOE4PlusApl from './apls/aoe_4_plus.apl.json'
export const ROTATION_PRESET_AOE4PLUS = PresetUtils.makePresetAPLRotation('AOE(4+目标)', AOE4PlusApl);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/wotlk/talent-calc and copy the numbers in the url.
export const StandardTalents = {
	name: '标准预设',
	data: SavedTalents.create({
		talentsString: '05032031--325023051223010323151301351',
		glyphs: Glyphs.create({
			major1: MajorGlyph.GlyphOfShadow,
			major2: MajorGlyph.GlyphOfMindFlay,
			major3: MajorGlyph.GlyphOfDispersion,
			minor1: MinorGlyph.GlyphOfFortitude,
			minor2: MinorGlyph.GlyphOfShadowProtection,
			minor3: MinorGlyph.GlyphOfShadowfiend,
		}),
	}),
};

export const EnlightenmentTalents = {
	name: '启迪',
	data: SavedTalents.create({
		talentsString: '05032031303005022--3250230012230101231513011',
		glyphs: Glyphs.create({
			major1: MajorGlyph.GlyphOfShadow,
			major2: MajorGlyph.GlyphOfMindFlay,
			major3: MajorGlyph.GlyphOfShadowWordDeath,
			minor1: MinorGlyph.GlyphOfFortitude,
			minor2: MinorGlyph.GlyphOfShadowProtection,
			minor3: MinorGlyph.GlyphOfShadowfiend,
		}),
	}),
};

export const DefaultOptions = Options.create({
	armor: Armor.InnerFire,
});

export const DefaultConsumes = Consumes.create({
	flask: Flask.FlaskOfTheFrostWyrm,
	food: Food.FoodFishFeast,
	defaultPotion: Potions.PotionOfSpeed,
	prepopPotion: Potions.PotionOfWildMagic,
});

export const DefaultRaidBuffs = RaidBuffs.create({
	giftOfTheWild: TristateEffect.TristateEffectImproved,
	powerWordFortitude: TristateEffect.TristateEffectImproved,
	strengthOfEarthTotem: TristateEffect.TristateEffectImproved,
	arcaneBrilliance: true,
	divineSpirit: true,
	trueshotAura: true,
	leaderOfThePack: TristateEffect.TristateEffectImproved,
	icyTalons: true,
	totemOfWrath: true,
	moonkinAura: TristateEffect.TristateEffectImproved,
	wrathOfAirTotem: true,
	sanctifiedRetribution: true,
	bloodlust: true,
	demonicPactSp: 500,
});

export const DefaultIndividualBuffs = IndividualBuffs.create({
	blessingOfKings: true,
	blessingOfWisdom: TristateEffect.TristateEffectImproved,
	blessingOfMight: TristateEffect.TristateEffectImproved,
	vampiricTouch: true,
});

export const DefaultDebuffs = Debuffs.create({
	sunderArmor: true,
	faerieFire: TristateEffect.TristateEffectImproved,
	bloodFrenzy: true,
	ebonPlaguebringer: true,
	heartOfTheCrusader: true,
	judgementOfWisdom: true,
	shadowMastery: true,
});

export const OtherDefaults = {
	channelClipDelay: 100,
	profession1: Profession.Engineering,
	profession2: Profession.Tailoring,
	nibelungAverageCasts: 11,
};
