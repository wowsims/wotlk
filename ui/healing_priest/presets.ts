import {
	Consumes,
	CustomRotation,
	CustomSpell,
	Debuffs,
	IndividualBuffs,
	Flask,
	Food,
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
} from '../core/proto/priest.js';

import * as PresetUtils from '../core/preset_utils.js';

import BlankGear from './gear_sets/blank.gear.json';

import DiscApl from './apls/disc.apl.json';
import HolyApl from './apls/holy.apl.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const DiscDefaultGear = PresetUtils.makePresetGear('Blank', BlankGear, { talentTree: 0 });
export const HolyDefaultGear = PresetUtils.makePresetGear('Blank', BlankGear, { talentTree: 1 });

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
	}),
};
export const HolyTalents = {
	name: 'Holy',
	data: SavedTalents.create({
		talentsString: '05032031103-234051032002152530004311051',
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
	strengthOfEarthTotem: true,
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
