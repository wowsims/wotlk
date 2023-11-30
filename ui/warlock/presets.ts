import {
	Consumes,
	Debuffs,
	Flask,
	Food,
	Glyphs,
	IndividualBuffs,
	Profession,
	RaidBuffs,
	TristateEffect
} from '../core/proto/common.js';
import { SavedTalents } from '../core/proto/ui.js';
import {
	Warlock_Options_Armor as Armor,
	WarlockMajorGlyph as MajorGlyph,
	WarlockMinorGlyph as MinorGlyph,
	Warlock_Options_Summon as Summon,
	Warlock_Options as WarlockOptions,
	Warlock_Options_WeaponImbue as WeaponImbue,
} from '../core/proto/warlock.js';
import * as PresetUtils from '../core/preset_utils.js';

import DefaultGear from './gear_sets/blank.gear.json';
import DefaultAPL from './apls/default.apl.json';

export const DEFAULT_GEAR = PresetUtils.makePresetGear('Blank', DefaultGear);

export const ROTATION_PRESET_DEFAULT = PresetUtils.makePresetAPLRotation('Default', DefaultAPL);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/classic/talent-calc and copy the numbers in the url.

export const DefaultTalents = {
	name: 'Default',
	data: SavedTalents.create({
		talentsString: '25002-2050300142301-52500051020001',
		glyphs: Glyphs.create({
			major1: MajorGlyph.WarlockMajorGlyphNone,
			major2: MajorGlyph.WarlockMajorGlyphNone,
			major3: MajorGlyph.WarlockMajorGlyphNone,
			minor1: MinorGlyph.WarlockMinorGlyphNone,
			minor2: MinorGlyph.WarlockMinorGlyphNone,
			minor3: MinorGlyph.WarlockMinorGlyphNone,
		}),
	}),
};

export const DefaultOptions = WarlockOptions.create({
	armor: Armor.DemonArmor,
	summon: Summon.Imp,
	weaponImbue: WeaponImbue.Spellstone,
});

export const DefaultConsumes = Consumes.create({
	flask: Flask.FlaskUnknown,
	food: Food.FoodUnknown,
});

export const DefaultRaidBuffs = RaidBuffs.create({
	giftOfTheWild: TristateEffect.TristateEffectImproved,
	strengthOfEarthTotem: TristateEffect.TristateEffectImproved,
	arcaneBrilliance: true,
	divineSpirit: true,
	trueshotAura: true,
	leaderOfThePack: true,
	moonkinAura: true,
});

export const DefaultIndividualBuffs = IndividualBuffs.create({
	blessingOfKings: true,
	blessingOfWisdom: TristateEffect.TristateEffectImproved,
	blessingOfMight: TristateEffect.TristateEffectImproved,
});

export const DefaultDebuffs = Debuffs.create({
	sunderArmor: true,
	faerieFire: TristateEffect.TristateEffectImproved,
	judgementOfWisdom: true,
});

export const OtherDefaults = {
	distanceFromTarget: 25,
	profession1: Profession.Engineering,
	profession2: Profession.Tailoring,
	channelClipDelay: 150,
	nibelungAverageCasts: 11,
};
