import {
	Consumes,
	Debuffs,
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
	DruidMajorGlyph,
	DruidMinorGlyph,
	BalanceDruid_Rotation_Type as RotationType
} from '../core/proto/druid.js';

import * as PresetUtils from '../core/preset_utils.js';

import BlankGear from './gear_sets/blank.gear.json';

import DefaultAplJson from './apls/default.apl.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const DefaultGear = PresetUtils.makePresetGear('Blank', BlankGear);

export const DEFAULT_APL = PresetUtils.makePresetAPLRotation('Default', DefaultAplJson);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/classic/talent-calc and copy the numbers in the url.
export const BalanceTalents = {
	name: 'Balance',
	data: SavedTalents.create({
		talentsString: '5000500302551351--50050312',
		glyphs: Glyphs.create({
			major1: DruidMajorGlyph.DruidMajorGlyphNone,
			major2: DruidMajorGlyph.DruidMajorGlyphNone,
			major3: DruidMajorGlyph.DruidMajorGlyphNone,
			minor1: DruidMinorGlyph.DruidMinorGlyphNone,
			minor2: DruidMinorGlyph.DruidMinorGlyphNone,
			minor3: DruidMinorGlyph.DruidMinorGlyphNone,
		}),
	}),
};

export const DefaultRotation = BalanceDruidRotation.create({
	type: RotationType.Default,
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
	blessingOfWisdom: TristateEffect.TristateEffectImproved,
});

export const DefaultPartyBuffs = PartyBuffs.create({
	heroicPresence: false,
});

export const DefaultDebuffs = Debuffs.create({
	faerieFire: TristateEffect.TristateEffectImproved,
	judgementOfWisdom: true,
});

export const OtherDefaults = {
	distanceFromTarget: 18,
	profession1: Profession.Engineering,
	profession2: Profession.Tailoring,
};
