import * as PresetUtils from '../core/preset_utils.js';
import {
	Consumes,
	Explosive,
	Flask,
	Food,
	Glyphs,
	PetFood,
	Potions,
	UnitReference,
} from '../core/proto/common.js';
import {
	Deathknight_Options as DeathKnightOptions,
	DeathknightMajorGlyph,
	DeathknightMinorGlyph,
} from '../core/proto/deathknight.js';
import { SavedTalents } from '../core/proto/ui.js';
// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.
import P1BloodGear from './gear_sets/p1_blood.gear.json';
export const P1_BLOOD_PRESET = PresetUtils.makePresetGear('P1血DPS', P1BloodGear, { talentTree: 0 });
import P2BloodGear from './gear_sets/p2_blood.gear.json';
export const P2_BLOOD_PRESET = PresetUtils.makePresetGear('P2血DPS', P2BloodGear, { talentTree: 0 });
import P3BloodGear from './gear_sets/p3_blood.gear.json';
export const P3_BLOOD_PRESET = PresetUtils.makePresetGear('P3血DPS', P3BloodGear, { talentTree: 0 });
import P4BloodGear from './gear_sets/p4_blood.gear.json';
export const P4_BLOOD_PRESET = PresetUtils.makePresetGear('P4血DPS', P4BloodGear, { talentTree: 0 });
import PreraidFrostGear from './gear_sets/preraid_frost.gear.json';
export const PRERAID_FROST_PRESET = PresetUtils.makePresetGear('Pre-Raid冰', PreraidFrostGear, { talentTree: 1 });
import P1FrostGear from './gear_sets/p1_frost.gear.json';
export const P1_FROST_PRESET = PresetUtils.makePresetGear('P1冰', P1FrostGear, { talentTree: 1 });
import P2FrostGear from './gear_sets/p2_frost.gear.json';
export const P2_FROST_PRESET = PresetUtils.makePresetGear('P2冰', P2FrostGear, { talentTree: 1 });
import P3FrostGear from './gear_sets/p3_frost.gear.json';
export const P3_FROST_PRESET = PresetUtils.makePresetGear('P3冰', P3FrostGear, { talentTree: 1 });
import P4FrostGear from './gear_sets/p4_frost.gear.json';
export const P4_FROST_PRESET = PresetUtils.makePresetGear('P4冰', P4FrostGear, { talentTree: 1 });
import P1FrostSubUhGear from './gear_sets/p1_frost_subUh.gear.json';
export const P1_FROSTSUBUNH_PRESET = PresetUtils.makePresetGear('P1冰邪', P1FrostSubUhGear, { talentTree: 1 });
import PreraidUh2hGear from './gear_sets/preraid_uh_2h.gear.json';
export const PRERAID_UNHOLY_2H_PRESET = PresetUtils.makePresetGear('Pre-Raid双手邪', PreraidUh2hGear, { talentTree: 2 });
import P1Uh2hGear from './gear_sets/p1_uh_2h.gear.json';
export const P1_UNHOLY_2H_PRESET = PresetUtils.makePresetGear('P1双手邪', P1Uh2hGear, { talentTree: 2 });
import P4Uh2hGear from './gear_sets/p4_uh_2h.gear.json';
export const P4_UNHOLY_2H_PRESET = PresetUtils.makePresetGear('P4双手邪', P4Uh2hGear, { talentTree: 2 });
import PreraidUhDwGear from './gear_sets/preraid_uh_dw.gear.json';
export const PRERAID_UNHOLY_DW_PRESET = PresetUtils.makePresetGear('Pre-Raid双持邪', PreraidUhDwGear, { talentTree: 2 });
import P1UhDwGear from './gear_sets/p1_uh_dw.gear.json';
export const P1_UNHOLY_DW_PRESET = PresetUtils.makePresetGear('P1双持邪', P1UhDwGear, { talentTree: 2 });
import P2UhDwGear from './gear_sets/p2_uh_dw.gear.json';
export const P2_UNHOLY_DW_PRESET = PresetUtils.makePresetGear('P2双持邪', P2UhDwGear, { talentTree: 2 });
import P3UhDwGear from './gear_sets/p3_uh_dw.gear.json';
export const P3_UNHOLY_DW_PRESET = PresetUtils.makePresetGear('P3双持邪', P3UhDwGear, { talentTree: 2 });
import P4UhDwGear from './gear_sets/p4_uh_dw.gear.json';
export const P4_UNHOLY_DW_PRESET = PresetUtils.makePresetGear('P4双持邪', P4UhDwGear, { talentTree: 2 });

import BloodDPSApl from './apls/blood_dps.apl.json';
export const BLOOD_DPS_ROTATION_PRESET_DEFAULT = PresetUtils.makePresetAPLRotation('血DPS-单体', BloodDPSApl, { talentTree: 0 });
import BloodPestiAoeApl from './apls/blood_pesti_aoe.apl.json';
export const BLOOD_PESTI_AOE_ROTATION_PRESET_DEFAULT = PresetUtils.makePresetAPLRotation('血DPS-AOE', BloodPestiAoeApl, { talentTree: 0 });
import FrostBlPestiApl from './apls/frost_bl_pesti.apl.json';
export const FROST_BL_PESTI_ROTATION_PRESET_DEFAULT = PresetUtils.makePresetAPLRotation('冰血', FrostBlPestiApl, { talentTree: 1 });
import FrostUhPestiApl from './apls/frost_uh_pesti.apl.json';
export const FROST_UH_PESTI_ROTATION_PRESET_DEFAULT = PresetUtils.makePresetAPLRotation('冰邪', FrostUhPestiApl, { talentTree: 1 });
import UhDwSsApl from './apls/unholy_dw_ss.apl.json';
export const UNHOLY_DW_ROTATION_PRESET_DEFAULT = PresetUtils.makePresetAPLRotation('双持邪(天打)', UhDwSsApl, { talentTree: 2 });
import Uh2hSsApl from './apls/uh_2h_ss.apl.json';
export const UNHOLY_2H_ROTATION_PRESET_DEFAULT = PresetUtils.makePresetAPLRotation('双手邪(无凋零)', Uh2hSsApl, { talentTree: 2 });
import UhDndAoeApl from './apls/uh_dnd_aoe.apl.json';
export const UNHOLY_DND_AOE_ROTATION_PRESET_DEFAULT = PresetUtils.makePresetAPLRotation('双持邪/双手邪(强化凋零)', UhDndAoeApl, { talentTree: 2 });

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wotlk.wowhead.com/talent-calc and copy the numbers in the url.
export const FrostTalents = {
	name: '冰血',
	data: SavedTalents.create({
		talentsString: '23050005-32005350352203012300033101351',
		glyphs: Glyphs.create({
			major1: DeathknightMajorGlyph.GlyphOfObliterate,
			major2: DeathknightMajorGlyph.GlyphOfFrostStrike,
			major3: DeathknightMajorGlyph.GlyphOfDisease,
			minor1: DeathknightMinorGlyph.GlyphOfHornOfWinter,
			minor2: DeathknightMinorGlyph.GlyphOfPestilence,
			minor3: DeathknightMinorGlyph.GlyphOfRaiseDead,
		}),
	}),
};

export const FrostUnholyTalents = {
	name: '冰邪',
	data: SavedTalents.create({
		talentsString: '01-32002350342203012300033101351-230200305003',
		glyphs: Glyphs.create({
			major1: DeathknightMajorGlyph.GlyphOfObliterate,
			major2: DeathknightMajorGlyph.GlyphOfFrostStrike,
			major3: DeathknightMajorGlyph.GlyphOfDisease,
			minor1: DeathknightMinorGlyph.GlyphOfHornOfWinter,
			minor2: DeathknightMinorGlyph.GlyphOfPestilence,
			minor3: DeathknightMinorGlyph.GlyphOfRaiseDead,
		}),
	}),
};

export const UnholyDualWieldTalents = {
	name: '双持邪(无打)',
	data: SavedTalents.create({
		talentsString: '-320043500002-2300303050032152000150013133051',
		glyphs: Glyphs.create({
			major1: DeathknightMajorGlyph.GlyphOfTheGhoul,
			major2: DeathknightMajorGlyph.GlyphOfIcyTouch,
			major3: DeathknightMajorGlyph.GlyphOfDeathAndDecay,
			minor1: DeathknightMinorGlyph.GlyphOfHornOfWinter,
			minor2: DeathknightMinorGlyph.GlyphOfPestilence,
			minor3: DeathknightMinorGlyph.GlyphOfRaiseDead,
		}),
	}),
};

export const UnholyDualWieldSSTalents = {
	name: '双持邪(天打)',
	data: SavedTalents.create({
		talentsString: '-320033500002-2301303050032151000150013133151',
		glyphs: Glyphs.create({
			major1: DeathknightMajorGlyph.GlyphOfTheGhoul,
			major2: DeathknightMajorGlyph.GlyphOfIcyTouch,
			major3: DeathknightMajorGlyph.GlyphOfDeathAndDecay,
			minor1: DeathknightMinorGlyph.GlyphOfHornOfWinter,
			minor2: DeathknightMinorGlyph.GlyphOfPestilence,
			minor3: DeathknightMinorGlyph.GlyphOfRaiseDead,
		}),
	}),
};

export const Unholy2HTalents = {
	name: '双手邪(无强化凋零)',
	data: SavedTalents.create({
		talentsString: '-320050500002-2302003350032052000150013133151',
		glyphs: Glyphs.create({
			major1: DeathknightMajorGlyph.GlyphOfTheGhoul,
			major2: DeathknightMajorGlyph.GlyphOfIcyTouch,
			major3: DeathknightMajorGlyph.GlyphOfDarkDeath,
			minor1: DeathknightMinorGlyph.GlyphOfHornOfWinter,
			minor2: DeathknightMinorGlyph.GlyphOfPestilence,
			minor3: DeathknightMinorGlyph.GlyphOfRaiseDead,
		}),
	}),
};

export const UnholyAoeTalents = {
	name: '双手邪(强化凋零)',
	data: SavedTalents.create({
		talentsString: '-320050500002-2302303050032052000150013133151',
		glyphs: Glyphs.create({
			major1: DeathknightMajorGlyph.GlyphOfTheGhoul,
			major2: DeathknightMajorGlyph.GlyphOfIcyTouch,
			major3: DeathknightMajorGlyph.GlyphOfDeathAndDecay,
			minor1: DeathknightMinorGlyph.GlyphOfHornOfWinter,
			minor2: DeathknightMinorGlyph.GlyphOfPestilence,
			minor3: DeathknightMinorGlyph.GlyphOfRaiseDead,
		}),
	}),
};

export const BloodTalents = {
	name: '血DPS',
	data: SavedTalents.create({
		talentsString: '2305120530003303231023001351--2302003050032',
		glyphs: Glyphs.create({
			major1: DeathknightMajorGlyph.GlyphOfDancingRuneWeapon,
			major2: DeathknightMajorGlyph.GlyphOfDeathStrike,
			major3: DeathknightMajorGlyph.GlyphOfDisease,
			minor1: DeathknightMinorGlyph.GlyphOfHornOfWinter,
			minor2: DeathknightMinorGlyph.GlyphOfPestilence,
			minor3: DeathknightMinorGlyph.GlyphOfRaiseDead,
		}),
	}),
};

export const DefaultUnholyOptions = DeathKnightOptions.create({
	startingRunicPower: 0,
	petUptime: 1,
	unholyFrenzyTarget: UnitReference.create(),
	drwPestiApply: true,
});

export const DefaultFrostOptions = DeathKnightOptions.create({
	startingRunicPower: 0,
	petUptime: 1,
	unholyFrenzyTarget: UnitReference.create(),
	drwPestiApply: true,
});

export const DefaultBloodOptions = DeathKnightOptions.create({
	startingRunicPower: 0,
	petUptime: 1,
	unholyFrenzyTarget: UnitReference.create(),
	drwPestiApply: true,
});

export const OtherDefaults = {
};

export const DefaultConsumes = Consumes.create({
	flask: Flask.FlaskOfEndlessRage,
	food: Food.FoodDragonfinFilet,
	defaultPotion: Potions.PotionOfSpeed,
	petFood: PetFood.PetFoodSpicedMammothTreats,
	prepopPotion: Potions.PotionOfSpeed,
	thermalSapper: true,
	fillerExplosive: Explosive.ExplosiveSaroniteBomb,
});
