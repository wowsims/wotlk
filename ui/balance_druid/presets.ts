import * as PresetUtils from '../core/preset_utils.js';
import {
	Consumes,
	Debuffs,
	Explosive,
	Faction,
	Flask,
	Food,
	Glyphs,
	IndividualBuffs,
	PartyBuffs,
	Potions,
	Profession,
	RaidBuffs,
	TristateEffect,
	UnitReference,
} from '../core/proto/common.js';
import {
	BalanceDruid_Options as BalanceDruidOptions,
	DruidMajorGlyph,
	DruidMinorGlyph,
} from '../core/proto/druid.js';
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
import P3AllianceGear from './gear_sets/p3_alliance.gear.json';
export const P3_PRESET_ALLI = PresetUtils.makePresetGear('P3预设[A]', P3AllianceGear, { faction: Faction.Alliance });
import P3HordeGear from './gear_sets/p3_horde.gear.json';
export const P3_PRESET_HORDE = PresetUtils.makePresetGear('P3预设[H]', P3HordeGear, { faction: Faction.Horde });
import P4AllianceGear from './gear_sets/p4_alliance.gear.json';
export const P4_PRESET_ALLI = PresetUtils.makePresetGear('P4预设[A]', P4AllianceGear, { faction: Faction.Alliance });
import P4HordeGear from './gear_sets/p4_horde.gear.json';
export const P4_PRESET_HORDE = PresetUtils.makePresetGear('P4预设[H]', P4HordeGear, { faction: Faction.Horde });

import BasicP3AplJson from './apls/basic_p3.apl.json';
export const ROTATION_PRESET_P3_APL = PresetUtils.makePresetAPLRotation('P3', BasicP3AplJson);
import P4FocusAplJson from './apls/p4_focus_glyph.apl.json';
export const ROTATION_PRESET_P4_FOCUS_APL = PresetUtils.makePresetAPLRotation('P4专注雕文', P4FocusAplJson);
import P4StarfireAplJson from './apls/p4_starfire_glyph.apl.json';
export const ROTATION_PRESET_P4_STARFIRE_APL = PresetUtils.makePresetAPLRotation('P4星火雕文', P4StarfireAplJson);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/wotlk/talent-calc and copy the numbers in the url.
export const Phase1Talents = {
	name: 'P1',
	data: SavedTalents.create({
		talentsString: '5032003115331303213305311231--205003012',
		glyphs: Glyphs.create({
			major1: DruidMajorGlyph.GlyphOfFocus,
			major2: DruidMajorGlyph.GlyphOfInsectSwarm,
			major3: DruidMajorGlyph.GlyphOfStarfall,
			minor1: DruidMinorGlyph.GlyphOfTyphoon,
			minor2: DruidMinorGlyph.GlyphOfUnburdenedRebirth,
			minor3: DruidMinorGlyph.GlyphOfTheWild,
		}),
	}),
};

export const Phase2Talents = {
	name: 'P2',
	data: SavedTalents.create({
		talentsString: '5012203115331303213305311231--205003012',
		glyphs: Glyphs.create({
			major1: DruidMajorGlyph.GlyphOfStarfire,
			major2: DruidMajorGlyph.GlyphOfInsectSwarm,
			major3: DruidMajorGlyph.GlyphOfStarfall,
			minor1: DruidMinorGlyph.GlyphOfTyphoon,
			minor2: DruidMinorGlyph.GlyphOfUnburdenedRebirth,
			minor3: DruidMinorGlyph.GlyphOfTheWild,
		}),
	}),
};

export const Phase3Talents = {
	name: 'P3',
	data: SavedTalents.create({
		talentsString: '5102223115331303213305311031--205003012',
		glyphs: Glyphs.create({
			major1: DruidMajorGlyph.GlyphOfStarfire,
			major2: DruidMajorGlyph.GlyphOfMoonfire,
			major3: DruidMajorGlyph.GlyphOfStarfall,
			minor1: DruidMinorGlyph.GlyphOfTyphoon,
			minor2: DruidMinorGlyph.GlyphOfUnburdenedRebirth,
			minor3: DruidMinorGlyph.GlyphOfTheWild,
		}),
	}),
};

export const Phase4Talents = {
	name: 'P4',
	data: SavedTalents.create({
		talentsString: '5102223115331303213305311031--205003012',
		glyphs: Glyphs.create({
			major1: DruidMajorGlyph.GlyphOfFocus,
			major2: DruidMajorGlyph.GlyphOfInsectSwarm,
			major3: DruidMajorGlyph.GlyphOfStarfall,
			minor1: DruidMinorGlyph.GlyphOfTyphoon,
			minor2: DruidMinorGlyph.GlyphOfUnburdenedRebirth,
			minor3: DruidMinorGlyph.GlyphOfTheWild,
		}),
	}),
};

export const DefaultOptions = BalanceDruidOptions.create({
	innervateTarget: UnitReference.create(),
});

export const DefaultConsumes = Consumes.create({
	defaultPotion: Potions.PotionOfSpeed,
	flask: Flask.FlaskOfTheFrostWyrm,
	food: Food.FoodFishFeast,
	prepopPotion: Potions.PotionOfWildMagic,
	fillerExplosive: Explosive.ExplosiveSaroniteBomb,
});

export const DefaultRaidBuffs = RaidBuffs.create({
	arcaneBrilliance: true,
	bloodlust: true,
	divineSpirit: true,
	giftOfTheWild: TristateEffect.TristateEffectImproved,
	icyTalons: true,
	moonkinAura: TristateEffect.TristateEffectImproved,
	leaderOfThePack: TristateEffect.TristateEffectImproved,
	powerWordFortitude: TristateEffect.TristateEffectImproved,
	sanctifiedRetribution: true,
	strengthOfEarthTotem: TristateEffect.TristateEffectImproved,
	trueshotAura: true,
	wrathOfAirTotem: true,
	demonicPactSp: 500,
});

export const DefaultIndividualBuffs = IndividualBuffs.create({
	blessingOfKings: true,
	blessingOfMight: TristateEffect.TristateEffectImproved,
	blessingOfWisdom: TristateEffect.TristateEffectImproved,
	vampiricTouch: true,
});

export const DefaultPartyBuffs = PartyBuffs.create({
	heroicPresence: false,
});

export const DefaultDebuffs = Debuffs.create({
	bloodFrenzy: true,
	ebonPlaguebringer: true,
	faerieFire: TristateEffect.TristateEffectImproved,
	heartOfTheCrusader: true,
	judgementOfWisdom: true,
	shadowMastery: true,
	sunderArmor: true,
	totemOfWrath: true,
});

export const OtherDefaults = {
	distanceFromTarget: 18,
	profession1: Profession.Engineering,
	profession2: Profession.Tailoring,
	nibelungAverageCasts: 11,
};
