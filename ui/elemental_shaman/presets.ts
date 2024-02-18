import {
	Consumes,
	Faction,
	Flask,
	Food,
	Glyphs,
	Potions,
  Profession,
} from '../core/proto/common.js';
import { SavedTalents } from '../core/proto/ui.js';

import {
  ElementalShaman_Options as ElementalShamanOptions,
  ShamanShield,
  ShamanMajorGlyph,
  ShamanMinorGlyph,
  AirTotem,
  EarthTotem,
  FireTotem,
  WaterTotem,
  ShamanTotems,
} from '../core/proto/shaman.js';

import * as PresetUtils from '../core/preset_utils.js';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

import PreraidGear from './gear_sets/preraid.gear.json';
export const PRERAID_PRESET = PresetUtils.makePresetGear('Pre-raid Preset', PreraidGear);
import P1Gear from './gear_sets/p1.gear.json';
export const P1_PRESET = PresetUtils.makePresetGear('P1 Preset', P1Gear);
import P2Gear from './gear_sets/p2.gear.json';
export const P2_PRESET = PresetUtils.makePresetGear('P2 Preset', P2Gear);
import P3AllianceGear from './gear_sets/p3_alliance.gear.json';
export const P3_PRESET_ALLI = PresetUtils.makePresetGear('P3 Preset [A]', P3AllianceGear, { faction: Faction.Alliance });
import P3HordeGear from './gear_sets/p3_horde.gear.json';
export const P3_PRESET_HORDE = PresetUtils.makePresetGear('P3 Preset [H]', P3HordeGear, { faction: Faction.Horde });
import P4Gear from './gear_sets/p4.gear.json';
export const P4_PRESET = PresetUtils.makePresetGear('P4 Preset', P4Gear);

import DefaultApl from './apls/default.apl.json';
export const ROTATION_PRESET_DEFAULT = PresetUtils.makePresetAPLRotation('Default', DefaultApl);
import AdvancedApl from './apls/advanced.apl.json';
export const ROTATION_PRESET_ADVANCED = PresetUtils.makePresetAPLRotation('Advanced', AdvancedApl);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/wotlk/talent-calc and copy the numbers in the url.
export const StandardTalents = {
  name: 'Standard',
  data: SavedTalents.create({
    talentsString: '0533001523213351322301351-005050031',
    glyphs: Glyphs.create({
      major1: ShamanMajorGlyph.GlyphOfFlametongueWeapon,
      major2: ShamanMajorGlyph.GlyphOfTotemOfWrath,
      major3: ShamanMajorGlyph.GlyphOfLightningBolt,
      minor1: ShamanMinorGlyph.GlyphOfThunderstorm,
      minor2: ShamanMinorGlyph.GlyphOfWaterShield,
      minor3: ShamanMinorGlyph.GlyphOfGhostWolf,
    }),
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
  defaultPotion: Potions.PotionOfWildMagic,
  flask: Flask.FlaskOfTheFrostWyrm,
  food: Food.FoodFishFeast,
});