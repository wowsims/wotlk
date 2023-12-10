import {
	Consumes,
	Debuffs,
	Flask,
	Food,
	Glyphs,
	IndividualBuffs,
	PartyBuffs,
	RaidBuffs,
	UnitReference,
	TristateEffect
} from '../core/proto/common.js';
import { SavedTalents } from '../core/proto/ui.js';

import {
	RestorationDruid_Options as RestorationDruidOptions,
	RestorationDruid_Rotation as RestorationDruidRotation,
	DruidMajorGlyph,
	DruidMinorGlyph,
} from '../core/proto/druid.js';

import * as PresetUtils from '../core/preset_utils.js';

import BlankGear from './gear_sets/blank.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const DefaultGear = PresetUtils.makePresetGear('Blank', BlankGear);

export const DefaultRotation = RestorationDruidRotation.create({
});

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/classic/talent-calc and copy the numbers in the url.
export const CelestialFocusTalents = {
	name: 'Celestial Focus',
	data: SavedTalents.create({
		talentsString: '05320031103--230023312131502331050313051',
		glyphs: Glyphs.create({
			major1: DruidMajorGlyph.DruidMajorGlyphNone,
			major2: DruidMajorGlyph.DruidMajorGlyphNone,
			major3: DruidMajorGlyph.DruidMajorGlyphNone,
			minor2: DruidMinorGlyph.DruidMinorGlyphNone,
			minor3: DruidMinorGlyph.DruidMinorGlyphNone,
			minor1: DruidMinorGlyph.DruidMinorGlyphNone,
		}),
	}),
};
export const ThiccRestoTalents = {
	name: 'Thicc Resto',
	data: SavedTalents.create({
		talentsString: '05320001--230023312331502531053313051',
		glyphs: Glyphs.create({
			major1: DruidMajorGlyph.DruidMajorGlyphNone,
			major2: DruidMajorGlyph.DruidMajorGlyphNone,
			major3: DruidMajorGlyph.DruidMajorGlyphNone,
			minor2: DruidMinorGlyph.DruidMinorGlyphNone,
			minor3: DruidMinorGlyph.DruidMinorGlyphNone,
			minor1: DruidMinorGlyph.DruidMinorGlyphNone,
		}),
	}),
};

export const DefaultOptions = RestorationDruidOptions.create({
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
	leaderOfThePack: true,
	powerWordFortitude: TristateEffect.TristateEffectImproved,
	strengthOfEarthTotem: TristateEffect.TristateEffectImproved,
	trueshotAura: true,
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
	sunderArmor: true,
});

export const OtherDefaults = {
	distanceFromTarget: 18,
};
