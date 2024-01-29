import {
  Consumes,
  Flask,
  Food,
  Profession,
} from '../core/proto/common.js';
import { SavedTalents } from '../core/proto/ui.js';

import {
  AirTotem,
  EarthTotem,
  ElementalShaman_Options as ElementalShamanOptions,
  FireTotem,
  ShamanShield,
  ShamanTotems,
  WaterTotem,
} from '../core/proto/shaman.js';

import * as PresetUtils from '../core/preset_utils.js';

import BlankGear from './gear_sets/blank.gear.json'

import AdvancedApl from './apls/advanced.apl.json';
import DefaultApl from './apls/default.apl.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const DefaultGear = PresetUtils.makePresetGear('Blank', BlankGear);

export const ROTATION_PRESET_DEFAULT = PresetUtils.makePresetAPLRotation('Default', DefaultApl);
export const ROTATION_PRESET_ADVANCED = PresetUtils.makePresetAPLRotation('Advanced', AdvancedApl);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/classic/talent-calc and copy the numbers in the url.
export const StandardTalents = {
  name: 'Standard',
  data: SavedTalents.create({
    talentsString: '0532001523212351322301351-005052031',
  }),
};

export const Phase4Talents = {
  name: 'Phase 4',
  data: SavedTalents.create({
    talentsString: '0533001523213351322301351-005050031',
  }),
};

export const DefaultOptions = ElementalShamanOptions.create({
  shield: ShamanShield.WaterShield,
  totems: ShamanTotems.create({
    earth: EarthTotem.StrengthOfEarthTotem,
    air: AirTotem.WrathOfAirTotem,
    fire: FireTotem.TotemOfWrath,
    water: WaterTotem.ManaSpringTotem,
    useFireElemental: true,
  }),
});

export const OtherDefaults = {
    distanceFromTarget: 20,
    profession1: Profession.Engineering,
    profession2: Profession.Tailoring,
    nibelungAverageCasts: 11,
}

export const DefaultConsumes = Consumes.create({
  flask: Flask.FlaskUnknown,
	food: Food.FoodUnknown,
});