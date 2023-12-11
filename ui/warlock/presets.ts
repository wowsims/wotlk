import {
	Consumes,
	Debuffs,
	Flask,
	Food,
	IndividualBuffs,
	Profession,
	RaidBuffs,
	TristateEffect
} from '../core/proto/common.js';
import { SavedTalents } from '../core/proto/ui.js';
import {
	Warlock_Options_Armor as Armor,
	Warlock_Options_Summon as Summon,
	Warlock_Options as WarlockOptions,
	Warlock_Options_WeaponImbue as WeaponImbue,
} from '../core/proto/warlock.js';
import * as PresetUtils from '../core/preset_utils.js';

import DefaultGear from './gear_sets/blank.gear.json';
import DefaultAPL from './apls/default.apl.json';

export const GearAfflictionDefault = PresetUtils.makePresetGear('Blank', DefaultGear);
export const GearDemonologyDefault = PresetUtils.makePresetGear('Blank', DefaultGear);
export const GearDestructionDefault = PresetUtils.makePresetGear('Blank', DefaultGear);

export const RotationAfflictionDefault = PresetUtils.makePresetAPLRotation('Default', DefaultAPL);
export const RotationDemonologyDefault = PresetUtils.makePresetAPLRotation('Default', DefaultAPL);
export const RotationDestructionDefault = PresetUtils.makePresetAPLRotation('Default', DefaultAPL);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/classic/talent-calc and copy the numbers in the url.

export const DefaultTalents = {
	name: 'Default',
	data: SavedTalents.create({
		talentsString: '25002-2050300142301-52500051020001',
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
	strengthOfEarthTotem: true,
	arcaneBrilliance: true,
	divineSpirit: true,
	trueshotAura: true,
	leaderOfThePack: true,
	moonkinAura: false,
	windfuryTotem: true,
	battleShout: TristateEffect.TristateEffectImproved,
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
