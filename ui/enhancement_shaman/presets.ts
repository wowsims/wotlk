import * as PresetUtils from '../core/preset_utils.js';
import {
	Consumes,
	Debuffs,
	Faction,
	Flask,
	Food,
	Glyphs,
	Potions,
	RaidBuffs,
	TristateEffect,
} from '../core/proto/common.js';
import {
	AirTotem,
	EarthTotem,
	EnhancementShaman_Options as EnhancementShamanOptions,
	FireTotem,
	ShamanImbue,
	ShamanMajorGlyph,
	ShamanShield,
	ShamanSyncType,
	ShamanTotems,
	WaterTotem,
} from '../core/proto/shaman.js';
import { SavedTalents } from '../core/proto/ui.js';
// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.
import PreraidGear from './gear_sets/preraid.gear.json';
export const PRERAID_PRESET = PresetUtils.makePresetGear('Preraid预设', PreraidGear);
import P1Gear from './gear_sets/p1.gear.json';
export const P1_PRESET = PresetUtils.makePresetGear('P1预设', P1Gear);
import P2FtGear from './gear_sets/p2_ft.gear.json';
export const P2_PRESET_FT = PresetUtils.makePresetGear('P2预设(火舌)', P2FtGear);
import P2WfGear from './gear_sets/p2_wf.gear.json';
export const P2_PRESET_WF = PresetUtils.makePresetGear('P2预设(风怒)', P2WfGear);
import P3AllianceGear from './gear_sets/p3_alliance.gear.json';
export const P3_PRESET_ALLIANCE = PresetUtils.makePresetGear('P3预设[A]', P3AllianceGear, { faction: Faction.Alliance });
import P3HordeGear from './gear_sets/p3_horde.gear.json';
export const P3_PRESET_HORDE = PresetUtils.makePresetGear('P3预设[H]', P3HordeGear, { faction: Faction.Horde });
import P4FtGear from './gear_sets/p4_ft.gear.json';
export const P4_PRESET_FT = PresetUtils.makePresetGear('P4预设(火舌)', P4FtGear);
import P4WfGear from './gear_sets/p4_wf.gear.json';
export const P4_PRESET_WF = PresetUtils.makePresetGear('P4预设(风怒)', P4WfGear);

import DefaultFt from './apls/default_ft.apl.json';
export const ROTATION_FT_DEFAULT = PresetUtils.makePresetAPLRotation('标准预设(火舌)', DefaultFt);
import DefaultWf from './apls/default_wf.apl.json';
export const ROTATION_WF_DEFAULT = PresetUtils.makePresetAPLRotation('标准预设(风怒)', DefaultWf);
import Phase3Apl from './apls/phase_3.apl.json';
export const ROTATION_PHASE_3 = PresetUtils.makePresetAPLRotation('P3', Phase3Apl);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/wotlk/talent-calc and copy the numbers in the url.
export const StandardTalents = {
	name: '标准预设',
	data: SavedTalents.create({
		talentsString: '053030152-30405003105021333031131031051',
		glyphs: Glyphs.create({
			major1: ShamanMajorGlyph.GlyphOfFireNova,
			major2: ShamanMajorGlyph.GlyphOfFlametongueWeapon,
			major3: ShamanMajorGlyph.GlyphOfFeralSpirit,
			//minor glyphs dont affect damage done, all convenience/QoL
		})
	}),
};

export const Phase3Talents = {
	name: 'P3',
	data: SavedTalents.create({
		talentsString: '053030152-30505003105001333031131131051',
		glyphs: Glyphs.create({
			major1: ShamanMajorGlyph.GlyphOfFireNova,
			major2: ShamanMajorGlyph.GlyphOfFlametongueWeapon,
			major3: ShamanMajorGlyph.GlyphOfFeralSpirit,
			//minor glyphs dont affect damage done, all convenience/QoL
		})
	}),
};

export const DefaultOptions = EnhancementShamanOptions.create({
	shield: ShamanShield.LightningShield,
	imbueMh: ShamanImbue.WindfuryWeapon,
	imbueOh: ShamanImbue.FlametongueWeapon,
	syncType: ShamanSyncType.Auto,
	totems: ShamanTotems.create({
		earth: 	EarthTotem.StrengthOfEarthTotem,
		fire: 	FireTotem.MagmaTotem,
		water: 	WaterTotem.ManaSpringTotem,
		air: 	AirTotem.WindfuryTotem,
	})
});

export const DefaultConsumes = Consumes.create({
	defaultPotion: Potions.PotionOfSpeed,
	flask: Flask.FlaskOfEndlessRage,
	food: Food.FoodFishFeast,
});

export const DefaultRaidBuffs = RaidBuffs.create({
	giftOfTheWild: TristateEffect.TristateEffectImproved,
	arcaneBrilliance: true,
	leaderOfThePack: TristateEffect.TristateEffectImproved,
	totemOfWrath: true,
	wrathOfAirTotem: true,
	moonkinAura: TristateEffect.TristateEffectImproved,
	sanctifiedRetribution: true,
	divineSpirit: true,
	battleShout: TristateEffect.TristateEffectImproved,
	demonicPactSp: 500,
});

export const DefaultDebuffs = Debuffs.create({
	bloodFrenzy: true,
	sunderArmor: true,
	curseOfWeakness: TristateEffect.TristateEffectRegular,
	curseOfElements: true,
	faerieFire: TristateEffect.TristateEffectImproved,
	judgementOfWisdom: true,
	misery: true,
	totemOfWrath: true,
	shadowMastery: true,
});
