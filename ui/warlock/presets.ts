import {
	Conjured,
	Consumes,
	Debuffs,
	FirePowerBuff,
	Flask,
	Food,
	IndividualBuffs,
	Potions,
	Profession,
	RaidBuffs,
	TristateEffect,
	WeaponImbue
} from '../core/proto/common.js';
import { SavedTalents } from '../core/proto/ui.js';
import {
	Warlock_Options_Armor as Armor,
	Warlock_Options_Summon as Summon,
	Warlock_Options as WarlockOptions,
	Warlock_Options_WeaponImbue as WarlockWeaponImbue,
} from '../core/proto/warlock.js';
import * as PresetUtils from '../core/preset_utils.js';

import DefaultGear from './gear_sets/blank.gear.json';
import DefaultAPL from './apls/default.apl.json';

import AfflictionTankGear from './gear_sets/affi.tank.gear.json';
import DestroTankGear from './gear_sets/destro.tank.gear.json';

export const GearAfflictionTankDefault = PresetUtils.makePresetGear('Affliction Tank', AfflictionTankGear);
export const GearDestructionTankDefault = PresetUtils.makePresetGear('Destruction Tank', DestroTankGear);

import AfflictionTankAPL from './apls/affi.tank.apl.json';
import DestroTankAPL from './apls/destro.tank.apl.json';

export const RotationAfflictionTankDefault = PresetUtils.makePresetAPLRotation('Affliction Tank', AfflictionTankAPL);
export const RotationDestructionTankDefault = PresetUtils.makePresetAPLRotation('Destruction Tank', DestroTankAPL);

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

export const AfflictionTankTalents = {
	name: 'Affliction Tank',
	data: SavedTalents.create({
		talentsString: '050025001-003',
	}),
};

export const DestroTalents = {
	name: 'Destruction',
	data: SavedTalents.create({
		talentsString: '-03-0550201',
	}),
};

export const DefaultOptions = WarlockOptions.create({
	armor: Armor.DemonArmor,
	summon: Summon.Imp,
	weaponImbue: WarlockWeaponImbue.NoWeaponImbue,
});

export const DefaultConsumes = Consumes.create({
	flask: Flask.FlaskUnknown,
	food: Food.FoodSmokedSagefish,
	defaultPotion: Potions.ManaPotion,
	mainHandImbue: WeaponImbue.BlackfathomManaOil,
	firePowerBuff: FirePowerBuff.ElixirOfFirepower,
});

export const DefaultRaidBuffs = RaidBuffs.create({
	powerWordFortitude: TristateEffect.TristateEffectImproved,
	giftOfTheWild: TristateEffect.TristateEffectImproved,
	arcaneBrilliance: true,
	divineSpirit: true,
});

export const DefaultIndividualBuffs = IndividualBuffs.create({
	aspectOfTheLion: true,
	blessingOfWisdom: TristateEffect.TristateEffectImproved,
	blessingOfMight: TristateEffect.TristateEffectImproved,
});

export const DefaultDebuffs = Debuffs.create({
	homunculi: true,
	faerieFire: true,
});

export const OtherDefaults = {
	distanceFromTarget: 25,
	profession1: Profession.Engineering,
	profession2: Profession.Tailoring,
	channelClipDelay: 150,
};
