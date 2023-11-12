import {
	Consumes,
	Debuffs,
	Faction,
	Flask,
	Food,
	Glyphs,
	IndividualBuffs,
	PartyBuffs,
	Profession,
	RaidBuffs,
	TristateEffect,
	UnitReference
} from '../core/proto/common.js';
import { SavedTalents } from '../core/proto/ui.js';

import {
	BalanceDruid_Options as BalanceDruidOptions,
	BalanceDruid_Rotation as BalanceDruidRotation,
	BalanceDruid_Rotation_IsUsage,
	BalanceDruid_Rotation_MfUsage,
	BalanceDruid_Rotation_WrathUsage,
	DruidMajorGlyph,
	DruidMinorGlyph,
	BalanceDruid_Rotation_Type as RotationType,
} from '../core/proto/druid.js';

import * as PresetUtils from '../core/preset_utils.js';

import P1Gear from './gear_sets/p1.gear.json';
import P2Gear from './gear_sets/p2.gear.json';
import P3AllianceGear from './gear_sets/p3_alliance.gear.json';
import P3HordeGear from './gear_sets/p3_horde.gear.json';
import P4AllianceGear from './gear_sets/p4_alliance.gear.json';
import P4HordeGear from './gear_sets/p4_horde.gear.json';
import PreraidGear from './gear_sets/preraid.gear.json';

import BasicP3AplJson from './apls/basic_p3.apl.json';
import P4FocusAplJson from './apls/p4_focus_glyph.apl.json';
import P4StarfireAplJson from './apls/p4_starfire_glyph.apl.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const PRERAID_PRESET = PresetUtils.makePresetGear('Pre-raid Preset', PreraidGear);
export const P1_PRESET = PresetUtils.makePresetGear('P1 Preset', P1Gear);
export const P2_PRESET = PresetUtils.makePresetGear('P2 Preset', P2Gear);
export const P3_PRESET_ALLI = PresetUtils.makePresetGear('P3 Preset [A]', P3AllianceGear, { faction: Faction.Alliance });
export const P3_PRESET_HORDE = PresetUtils.makePresetGear('P3 Preset [H]', P3HordeGear, { faction: Faction.Horde });
export const P4_PRESET_ALLI = PresetUtils.makePresetGear('P4 Preset [A]', P4AllianceGear, { faction: Faction.Alliance });
export const P4_PRESET_HORDE = PresetUtils.makePresetGear('P4 Preset [H]', P4HordeGear, { faction: Faction.Horde });

export const ROTATION_PRESET_P3_APL = PresetUtils.makePresetAPLRotation('P3', BasicP3AplJson);
export const ROTATION_PRESET_P4_FOCUS_APL = PresetUtils.makePresetAPLRotation('P4 Focus Glyph', P4FocusAplJson);
export const ROTATION_PRESET_P4_STARFIRE_APL = PresetUtils.makePresetAPLRotation('P4 Starfire Glyph', P4StarfireAplJson);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/wotlk/talent-calc and copy the numbers in the url.
export const Phase1Talents = {
	name: 'Phase 1',
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
	name: 'Phase 2',
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
	name: 'Phase 3',
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
	name: 'Phase 4',
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

export const DefaultRotation = BalanceDruidRotation.create({
	type: RotationType.Default,
	maintainFaerieFire: true,
	useSmartCooldowns: true,
	mfUsage: BalanceDruid_Rotation_MfUsage.BeforeLunar,
	isUsage: BalanceDruid_Rotation_IsUsage.OptimizeIs,
	wrathUsage: BalanceDruid_Rotation_WrathUsage.RegularWrath,
	useStarfire: true,
	useBattleRes: false,
	playerLatency: 200,
});

export const DefaultOptions = BalanceDruidOptions.create({
	innervateTarget: UnitReference.create(),
});

export const DefaultConsumes = Consumes.create({
	flask: Flask.FlaskUnknown,
	food: Food.FoodUnknown,
});

export const DefaultRaidBuffs = RaidBuffs.create({
	arcaneBrilliance: true,
	divineSpirit: true,
	giftOfTheWild: TristateEffect.TristateEffectImproved,
	moonkinAura: true,
	powerWordFortitude: TristateEffect.TristateEffectImproved,
});

export const DefaultIndividualBuffs = IndividualBuffs.create({
	blessingOfKings: true,
	blessingOfMight: TristateEffect.TristateEffectImproved,
	blessingOfWisdom: TristateEffect.TristateEffectImproved,
});

export const DefaultPartyBuffs = PartyBuffs.create({
	heroicPresence: false,
});

export const DefaultDebuffs = Debuffs.create({
	faerieFire: TristateEffect.TristateEffectImproved,
	judgementOfWisdom: true,
	sunderArmor: true,
});

export const OtherDefaults = {
	distanceFromTarget: 18,
	profession1: Profession.Engineering,
	profession2: Profession.Tailoring,
	nibelungAverageCasts: 11,
};
