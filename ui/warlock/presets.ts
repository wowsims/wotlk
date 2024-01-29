import {
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
	WarlockOptions_Armor as Armor,
	WarlockOptions_Summon as Summon,
	WarlockOptions as WarlockOptions,
	WarlockOptions_WeaponImbue as WarlockWeaponImbue,
} from '../core/proto/warlock.js';
import * as PresetUtils from '../core/preset_utils.js';

import DestructionGear from './gear_sets/destruction.gear.json';
import DestructionAPL from './apls/destruction.apl.json';

export const GearAfflictionDefault = PresetUtils.makePresetGear('Affliction', DestructionGear);
export const GearDemonologyDefault = PresetUtils.makePresetGear('Demonology', DestructionGear);
export const GearDestructionDefault = PresetUtils.makePresetGear('Destruction', DestructionGear);

export const RotationAfflictionDefault = PresetUtils.makePresetAPLRotation('Affliction', DestructionAPL);
export const RotationDemonologyDefault = PresetUtils.makePresetAPLRotation('Demonology', DestructionAPL);
export const RotationDestructionDefault = PresetUtils.makePresetAPLRotation('Destruction', DestructionAPL);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/classic/talent-calc and copy the numbers in the url.

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
	distanceFromTarget: 25,
	profession1: Profession.Engineering,
	profession2: Profession.Tailoring,
	channelClipDelay: 150,
};
