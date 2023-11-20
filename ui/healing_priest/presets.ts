import {
	Consumes,
	CustomRotation,
	CustomSpell,
	Debuffs,
	IndividualBuffs,
	Flask,
	Food,
	Glyphs,
	RaidBuffs,
	TristateEffect,
	UnitReference,
} from '../core/proto/common.js';
import { SavedTalents } from '../core/proto/ui.js';

import {
	HealingPriest_Rotation as Rotation,
	HealingPriest_Rotation_RotationType as RotationType,
	HealingPriest_Rotation_SpellOption as SpellOption,
	HealingPriest_Options as Options,
	PriestMajorGlyph as MajorGlyph,
	PriestMinorGlyph as MinorGlyph,
} from '../core/proto/priest.js';

import * as PresetUtils from '../core/preset_utils.js';

import PreraidDiscGear from './gear_sets/preraid_disc.gear.json';
import PreraidHolyGear from './gear_sets/preraid_holy.gear.json';
import P1DiscGear from './gear_sets/p1_disc.gear.json';
import P1HolyGear from './gear_sets/p1_holy.gear.json';
import P2DiscGear from './gear_sets/p2_disc.gear.json';
import P2HolyGear from './gear_sets/p2_holy.gear.json';

import DiscApl from './apls/disc.apl.json';
import HolyApl from './apls/holy.apl.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const DISC_PRERAID_PRESET = PresetUtils.makePresetGear('Disc Preraid Preset', PreraidDiscGear, { talentTree: 0 });
export const DISC_P1_PRESET = PresetUtils.makePresetGear('Disc P1 Preset', P1DiscGear, { talentTree: 0 });
export const DISC_P2_PRESET = PresetUtils.makePresetGear('Disc P2 Preset', P2DiscGear, { talentTree: 0 });
export const HOLY_PRERAID_PRESET = PresetUtils.makePresetGear('Holy Preraid Preset', PreraidHolyGear, { talentTree: 1 });
export const HOLY_P1_PRESET = PresetUtils.makePresetGear('Holy P1 Preset', P1HolyGear, { talentTree: 1 });
export const HOLY_P2_PRESET = PresetUtils.makePresetGear('Holy P2 Preset', P2HolyGear, { talentTree: 1 });

export const DiscDefaultRotation = Rotation.create({
	type: RotationType.Cycle,
	customRotation: CustomRotation.create({
		spells: [
			CustomSpell.create({ spell: SpellOption.PowerWordShield, castsPerMinute: 18 }),
			CustomSpell.create({ spell: SpellOption.Penance, castsPerMinute: 4 }),
			CustomSpell.create({ spell: SpellOption.PrayerOfMending, castsPerMinute: 2 }),
			CustomSpell.create({ spell: SpellOption.GreaterHeal, castsPerMinute: 1 }),
		],
	}),
});

export const HolyDefaultRotation = Rotation.create({
	type: RotationType.Cycle,
	customRotation: CustomRotation.create({
		spells: [
			CustomSpell.create({ spell: SpellOption.GreaterHeal, castsPerMinute: 10 }),
			CustomSpell.create({ spell: SpellOption.CircleOfHealing, castsPerMinute: 5 }),
			CustomSpell.create({ spell: SpellOption.Renew, castsPerMinute: 10 }),
			CustomSpell.create({ spell: SpellOption.PrayerOfMending, castsPerMinute: 2 }),
		],
	}),
});

export const ROTATION_PRESET_DISC = PresetUtils.makePresetAPLRotation('Disc', DiscApl);
export const ROTATION_PRESET_HOLY = PresetUtils.makePresetAPLRotation('Holy', HolyApl);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/classic/talent-calc and copy the numbers in the url.
export const DiscTalents = {
	name: 'Disc',
	data: SavedTalents.create({
		talentsString: '0503203130300512301313231251-2351010303',
		glyphs: Glyphs.create({
			major1: MajorGlyph.GlyphOfPowerWordShield,
			major2: MajorGlyph.GlyphOfFlashHeal,
			major3: MajorGlyph.GlyphOfPenance,
			minor1: MinorGlyph.GlyphOfFortitude,
			minor2: MinorGlyph.GlyphOfShadowfiend,
			minor3: MinorGlyph.GlyphOfFading,
		}),
	}),
};
export const HolyTalents = {
	name: 'Holy',
	data: SavedTalents.create({
		talentsString: '05032031103-234051032002152530004311051',
		glyphs: Glyphs.create({
			major1: MajorGlyph.GlyphOfPrayerOfHealing,
			major2: MajorGlyph.GlyphOfRenew,
			major3: MajorGlyph.GlyphOfCircleOfHealing,
			minor1: MinorGlyph.GlyphOfFortitude,
			minor2: MinorGlyph.GlyphOfShadowfiend,
			minor3: MinorGlyph.GlyphOfFading,
		}),
	}),
};

export const DefaultOptions = Options.create({
	useInnerFire: true,
	useShadowfiend: true,
	rapturesPerMinute: 5,

	powerInfusionTarget: UnitReference.create(),
});

export const DefaultConsumes = Consumes.create({
	flask: Flask.FlaskUnknown,
	food: Food.FoodUnknown,
});

export const DefaultRaidBuffs = RaidBuffs.create({
	giftOfTheWild: TristateEffect.TristateEffectImproved,
	powerWordFortitude: TristateEffect.TristateEffectImproved,
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
});

export const DefaultDebuffs = Debuffs.create({
});
