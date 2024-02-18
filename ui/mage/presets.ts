import {
	Conjured,
	Consumes,
	Faction,
	Flask,
	Food,
	Glyphs,
	Potions,
	Profession,
	Spec,
	UnitReference,
} from '../core/proto/common.js';
import { SavedTalents } from '../core/proto/ui.js';
import { Player } from '../core/player.js';

import {
	Mage_Rotation as MageRotation,
	Mage_Rotation_PrimaryFireSpell as PrimaryFireSpell,
	Mage_Options as MageOptions,
	Mage_Options_ArmorType as ArmorType,
	MageMajorGlyph,
	MageMinorGlyph,
} from '../core/proto/mage.js';

import * as PresetUtils from '../core/preset_utils.js';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

import PreraidArcaneGear from './gear_sets/preraid_arcane.gear.json';
export const ARCANE_PRERAID_PRESET = PresetUtils.makePresetGear('Arcane Preraid Preset', PreraidArcaneGear, { talentTree: 0 });
import P1ArcaneGear from './gear_sets/p1_arcane.gear.json';
export const ARCANE_P1_PRESET = PresetUtils.makePresetGear('Arcane P1 Preset', P1ArcaneGear, { talentTree: 0 });
import P2ArcaneGear from './gear_sets/p2_arcane.gear.json';
export const ARCANE_P2_PRESET = PresetUtils.makePresetGear('Arcane P2 Preset', P2ArcaneGear, { talentTree: 0 });
import P3ArcaneAllianceGear from './gear_sets/p3_arcane_alliance.gear.json';
export const ARCANE_P3_PRESET_ALLIANCE = PresetUtils.makePresetGear('Arcane P3 Preset [A]', P3ArcaneAllianceGear, { talentTree: 0, faction: Faction.Alliance });
import P3ArcaneHordeGear from './gear_sets/p3_arcane_horde.gear.json';
export const ARCANE_P3_PRESET_HORDE = PresetUtils.makePresetGear('Arcane P3 Preset [H]', P3ArcaneHordeGear, { talentTree: 0, faction: Faction.Horde });
import P4ArcaneAllianceGear from './gear_sets/p4_arcane_alliance.gear.json';
export const ARCANE_P4_PRESET_ALLIANCE = PresetUtils.makePresetGear('Arcane P4 Preset [A]', P4ArcaneAllianceGear, { talentTree: 0, faction: Faction.Alliance });
import P4ArcaneHordeGear from './gear_sets/p4_arcane_horde.gear.json';
export const ARCANE_P4_PRESET_HORDE = PresetUtils.makePresetGear('Arcane P4 Preset [H]', P4ArcaneHordeGear, { talentTree: 0, faction: Faction.Horde });
import PreraidFireGear from './gear_sets/preraid_fire.gear.json';
export const FIRE_PRERAID_PRESET = PresetUtils.makePresetGear('Fire Preraid Preset', PreraidFireGear, { talentTree: 1 });
import P1FireGear from './gear_sets/p1_fire.gear.json';
export const FIRE_P1_PRESET = PresetUtils.makePresetGear('Fire P1 Preset', P1FireGear, { talentTree: 1 });
import P2FireGear from './gear_sets/p2_fire.gear.json';
export const FIRE_P2_PRESET = PresetUtils.makePresetGear('Fire P2 Preset', P2FireGear, { talentTree: 1, customCondition: (player: Player<Spec.SpecMage>) => !player.getTalents().icyVeins });
import P3FireAllianceGear from './gear_sets/p3_fire_alliance.gear.json';
export const FIRE_P3_PRESET_ALLIANCE = PresetUtils.makePresetGear('Fire P3 Preset [A]', P3FireAllianceGear, { talentTree: 1, faction: Faction.Alliance, customCondition: (player: Player<Spec.SpecMage>) => !player.getTalents().icyVeins });
import P3FireHordeGear from './gear_sets/p3_fire_horde.gear.json';
export const FIRE_P3_PRESET_HORDE = PresetUtils.makePresetGear('Fire P3 Preset [H]', P3FireHordeGear, { talentTree: 1, faction: Faction.Horde, customCondition: (player: Player<Spec.SpecMage>) => !player.getTalents().icyVeins });
import P4FireAllianceGear from './gear_sets/p4_fire_alliance.gear.json';
export const FIRE_P4_PRESET_ALLIANCE = PresetUtils.makePresetGear('Fire P4 Preset [A]', P4FireAllianceGear, { talentTree: 1, faction: Faction.Alliance, customCondition: (player: Player<Spec.SpecMage>) => !player.getTalents().icyVeins });
import P4FireHordeGear from './gear_sets/p4_fire_horde.gear.json';
export const FIRE_P4_PRESET_HORDE = PresetUtils.makePresetGear('Fire P4 Preset [H]', P4FireHordeGear, { talentTree: 1, faction: Faction.Horde, customCondition: (player: Player<Spec.SpecMage>) => !player.getTalents().icyVeins });
import P2FfbGear from './gear_sets/p2_ffb.gear.json';
export const FFB_P2_PRESET = PresetUtils.makePresetGear('FFB P2 Preset', P2FfbGear, { talentTree: 1, customCondition: (player: Player<Spec.SpecMage>) => player.getTalents().icyVeins });
import P3FfbAllianceGear from './gear_sets/p3_ffb_alliance.gear.json';
export const FFB_P3_PRESET_ALLIANCE = PresetUtils.makePresetGear('FFB P3 Preset [A]', P3FfbAllianceGear, { talentTree: 1, customCondition: (player: Player<Spec.SpecMage>) => player.getTalents().icyVeins });
import P3FfbHordeGear from './gear_sets/p3_ffb_horde.gear.json';
export const FFB_P3_PRESET_HORDE = PresetUtils.makePresetGear('FFB P3 Preset [H]', P3FfbHordeGear, { talentTree: 1, customCondition: (player: Player<Spec.SpecMage>) => player.getTalents().icyVeins });
import P4FfbAllianceGear from './gear_sets/p4_ffb_alliance.gear.json';
export const FFB_P4_PRESET_ALLIANCE = PresetUtils.makePresetGear('FFB P4 Preset [A]', P4FfbAllianceGear, { talentTree: 1, customCondition: (player: Player<Spec.SpecMage>) => player.getTalents().icyVeins });
import P4FfbHordeGear from './gear_sets/p4_ffb_horde.gear.json';
export const FFB_P4_PRESET_HORDE = PresetUtils.makePresetGear('FFB P4 Preset [H]', P4FfbHordeGear, { talentTree: 1, customCondition: (player: Player<Spec.SpecMage>) => player.getTalents().icyVeins });
import P1FrostGear from './gear_sets/p1_frost.gear.json';
export const FROST_P1_PRESET = PresetUtils.makePresetGear('Frost P1 Preset', P1FrostGear, { talentTree: 2 });
import P2FrostGear from './gear_sets/p2_frost.gear.json';
export const FROST_P2_PRESET = PresetUtils.makePresetGear('Frost P2 Preset', P2FrostGear, { talentTree: 2 });
import P3FrostAllianceGear from './gear_sets/p3_frost_alliance.gear.json';
export const FROST_P3_PRESET_ALLIANCE = PresetUtils.makePresetGear('Frost P3 Preset [A]', P3FrostAllianceGear, { talentTree: 2, faction: Faction.Alliance });
import P3FrostHordeGear from './gear_sets/p3_frost_horde.gear.json';
export const FROST_P3_PRESET_HORDE = PresetUtils.makePresetGear('Frost P3 Preset [H]', P3FrostHordeGear, { talentTree: 2, faction: Faction.Horde });

export const DefaultSimpleRotation = MageRotation.create({
	only3ArcaneBlastStacksBelowManaPercent: 0.15,
	blastWithoutMissileBarrageAboveManaPercent: 0.2,
	missileBarrageBelowManaPercent: 0,
	useArcaneBarrage: false,

	primaryFireSpell: PrimaryFireSpell.Fireball,
	maintainImprovedScorch: false,

	useIceLance: false,
});

export const ROTATION_PRESET_SIMPLE = PresetUtils.makePresetSimpleRotation('Simple Default', Spec.SpecMage, DefaultSimpleRotation);
import ArcaneApl from './apls/arcane.apl.json';
export const ARCANE_ROTATION_PRESET_DEFAULT = PresetUtils.makePresetAPLRotation('Arcane', ArcaneApl, { talentTree: 0 });
import ArcaneAoeApl from './apls/arcane_aoe.apl.json';
export const ARCANE_ROTATION_PRESET_AOE = PresetUtils.makePresetAPLRotation('Arcane AOE', ArcaneAoeApl, { talentTree: 0 });
import FireApl from './apls/fire.apl.json';
export const FIRE_ROTATION_PRESET_DEFAULT = PresetUtils.makePresetAPLRotation('Fire', FireApl, { talentTree: 1 });
import FrostFireApl from './apls/frostfire.apl.json';
export const FROSTFIRE_ROTATION_PRESET_DEFAULT = PresetUtils.makePresetAPLRotation('Frostfire', FrostFireApl, { talentTree: 1 });
import FireAoeApl from './apls/fire_aoe.apl.json';
export const FIRE_ROTATION_PRESET_AOE = PresetUtils.makePresetAPLRotation('Fire AOE', FireAoeApl, { talentTree: 1 });
import FrostApl from './apls/frost.apl.json';
export const FROST_ROTATION_PRESET_DEFAULT = PresetUtils.makePresetAPLRotation('Frost', FrostApl, { talentTree: 2 });
import FrostAoeApl from './apls/frost_aoe.apl.json';
export const FROST_ROTATION_PRESET_AOE = PresetUtils.makePresetAPLRotation('Frost AOE', FrostAoeApl, { talentTree: 2 });

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/wotlk/talent-calc and copy the numbers in the url.
export const ArcaneTalents = {
	name: 'Arcane',
	data: SavedTalents.create({
		talentsString: '23000513310033015032310250532-03-023303001',
		glyphs: Glyphs.create({
			major1: MageMajorGlyph.GlyphOfArcaneBlast,
			major2: MageMajorGlyph.GlyphOfArcaneMissiles,
			major3: MageMajorGlyph.GlyphOfMoltenArmor,
			minor1: MageMinorGlyph.GlyphOfSlowFall,
			minor2: MageMinorGlyph.GlyphOfFrostWard,
			minor3: MageMinorGlyph.GlyphOfBlastWave,
		}),
	}),
};
export const FireTalents = {
	name: 'Fire',
	data: SavedTalents.create({
		talentsString: '23000503110003-0055030012303331053120301351',
		glyphs: Glyphs.create({
			major1: MageMajorGlyph.GlyphOfFireball,
			major2: MageMajorGlyph.GlyphOfMoltenArmor,
			major3: MageMajorGlyph.GlyphOfLivingBomb,
			minor1: MageMinorGlyph.GlyphOfSlowFall,
			minor2: MageMinorGlyph.GlyphOfFrostWard,
			minor3: MageMinorGlyph.GlyphOfBlastWave,
		}),
	}),
};

export const Phase3FireTalents = {
	name: 'Fire Phase 3',
	data: SavedTalents.create({
		talentsString: '23002303310003-0055030012303330053120300351',
		glyphs: Glyphs.create({
			major1: MageMajorGlyph.GlyphOfFireball,
			major2: MageMajorGlyph.GlyphOfMoltenArmor,
			major3: MageMajorGlyph.GlyphOfLivingBomb,
			minor1: MageMinorGlyph.GlyphOfArcaneIntellect,
			minor2: MageMinorGlyph.GlyphOfSlowFall,
		}),
	}),
};

export const FrostfireTalents = {
	name: 'Frostfire',
	data: SavedTalents.create({
		talentsString: '-2305032012303331053120311351-023303031',
		glyphs: Glyphs.create({
			major1: MageMajorGlyph.GlyphOfFrostfire,
			major2: MageMajorGlyph.GlyphOfMoltenArmor,
			major3: MageMajorGlyph.GlyphOfLivingBomb,
			minor1: MageMinorGlyph.GlyphOfSlowFall,
			minor2: MageMinorGlyph.GlyphOfFrostWard,
			minor3: MageMinorGlyph.GlyphOfBlastWave,
		}),
	}),
}
export const FrostTalents = {
	name: 'Frost',
	data: SavedTalents.create({
		talentsString: '23000503110003--0533030310233100030152231351',
		glyphs: Glyphs.create({
			major1: MageMajorGlyph.GlyphOfFrostbolt,
			major2: MageMajorGlyph.GlyphOfEternalWater,
			major3: MageMajorGlyph.GlyphOfMoltenArmor,
			minor1: MageMinorGlyph.GlyphOfSlowFall,
			minor2: MageMinorGlyph.GlyphOfFrostWard,
			minor3: MageMinorGlyph.GlyphOfBlastWave,
		}),
	}),
};

export const DefaultFFBOptions = MageOptions.create({
	armor: ArmorType.MoltenArmor,
});

export const DefaultFireOptions = MageOptions.create({
	armor: ArmorType.MoltenArmor,
	focusMagicPercentUptime: 99,
	focusMagicTarget: UnitReference.create(),
});

export const DefaultFireConsumes = Consumes.create({
	flask: Flask.FlaskOfTheFrostWyrm,
	food: Food.FoodFirecrackerSalmon,
	defaultPotion: Potions.PotionOfSpeed,
	defaultConjured: Conjured.ConjuredFlameCap,
	prepopPotion: Potions.PotionOfSpeed,
});

export const DefaultFrostOptions = MageOptions.create({
	armor: ArmorType.MoltenArmor,
	focusMagicTarget: UnitReference.create(),
	waterElementalDisobeyChance: 0.1,
});

export const DefaultFrostConsumes = Consumes.create({
	defaultPotion: Potions.PotionOfSpeed,
	defaultConjured: Conjured.ConjuredFlameCap,
	flask: Flask.FlaskOfTheFrostWyrm,
	food: Food.FoodFishFeast,
});

export const DefaultArcaneOptions = MageOptions.create({
	armor: ArmorType.MoltenArmor,
	focusMagicPercentUptime: 99,
	focusMagicTarget: UnitReference.create(),
});

export const DefaultArcaneConsumes = Consumes.create({
	defaultPotion: Potions.PotionOfSpeed,
	defaultConjured: Conjured.ConjuredDarkRune,
	flask: Flask.FlaskOfTheFrostWyrm,
	food: Food.FoodFirecrackerSalmon,
});

export const OtherDefaults = {
	distanceFromTarget: 20,
	profession1: Profession.Engineering,
	profession2: Profession.Tailoring,
	nibelungAverageCasts: 11,
};
