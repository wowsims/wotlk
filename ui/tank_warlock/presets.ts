import {
	AgilityElixir,
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
	StrengthBuff,
	TristateEffect,
	WeaponImbue
} from '../core/proto/common.js';
import { SavedTalents } from '../core/proto/ui.js';
import {
	WarlockOptions_Armor as Armor,
	WarlockOptions_Summon as Summon,
	WarlockOptions as WarlockOptions,
	WarlockOptions_WeaponImbue as WarlockWeaponImbue,
} from '../core/proto/warlock.js';
import * as PresetUtils from '../core/preset_utils.js';

import AfflictionTankGear from './gear_sets/affi.tank.gear.json';
import DestroTankGear from './gear_sets/destro.tank.gear.json';

export const GearAfflictionTankDefault = PresetUtils.makePresetGear('Affliction Tank', AfflictionTankGear);
export const GearDestructionTankDefault = PresetUtils.makePresetGear('Destruction Tank', DestroTankGear);

import AfflictionTankAPL from './apls/affi.tank.apl.json';
import DestroTankAPL from './apls/destro.tank.apl.json';

export const RotationAfflictionTankDefault = PresetUtils.makePresetAPLRotation('Affliction Tank', AfflictionTankAPL);
export const RotationDestructionTankDefault = PresetUtils.makePresetAPLRotation('Destruction Tank', DestroTankAPL);

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
	summon: Summon.Succubus,
	weaponImbue: WarlockWeaponImbue.NoWeaponImbue,
});

export const DefaultConsumes = Consumes.create({
	flask: Flask.FlaskUnknown,
	food: Food.FoodSmokedSagefish,
	defaultPotion: Potions.ManaPotion,
	mainHandImbue: WeaponImbue.BlackfathomManaOil,
	firePowerBuff: FirePowerBuff.ElixirOfFirepower,
	agilityElixir: AgilityElixir.ElixirOfLesserAgility,
	strengthBuff: StrengthBuff.ElixirOfOgresStrength,
});

export const DefaultRaidBuffs = RaidBuffs.create({
	powerWordFortitude: TristateEffect.TristateEffectImproved,
	giftOfTheWild: TristateEffect.TristateEffectImproved,
	arcaneBrilliance: true,
	divineSpirit: true,
	aspectOfTheLion: true,
});

export const DefaultIndividualBuffs = IndividualBuffs.create({
	blessingOfWisdom: TristateEffect.TristateEffectImproved,
	blessingOfMight: TristateEffect.TristateEffectImproved,
});

export const DefaultDebuffs = Debuffs.create({
	homunculi: 70, // 70% average uptime default
	faerieFire: true,
});

export const OtherDefaults = {
	distanceFromTarget: 5,
	profession1: Profession.Engineering,
	profession2: Profession.Tailoring,
	channelClipDelay: 150,
};
