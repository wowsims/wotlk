import {
	Consumes,
	CustomRotation,
	CustomSpell,
	Debuffs,
	Flask,
	Food,
	Glyphs,
	RaidBuffs,
	TristateEffect
} from '../core/proto/common.js';
import { SavedTalents } from '../core/proto/ui.js';

import {
	AirTotem, EnhancementShaman_Rotation_CustomRotationSpell as CustomRotationSpell, EarthTotem, EnhancementShaman_Options as EnhancementShamanOptions, EnhancementShaman_Rotation as EnhancementShamanRotation, FireTotem, EnhancementShaman_Rotation_PrimaryShock as PrimaryShock,
	EnhancementShaman_Rotation_RotationType as RotationType, ShamanImbue, ShamanMajorGlyph, ShamanShield, ShamanSyncType, ShamanTotems, WaterTotem
} from '../core/proto/shaman.js';

import * as PresetUtils from '../core/preset_utils.js';

import BlankGear from './gear_sets/blank.gear.json';

import DefaultFt from './apls/default_ft.apl.json';
import DefaultWf from './apls/default_wf.apl.json';
import Phase3Apl from './apls/phase_3.apl.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const DefaultGear = PresetUtils.makePresetGear('Blank', BlankGear);

export const DefaultRotation = EnhancementShamanRotation.create({
	totems: ShamanTotems.create({
		earth: EarthTotem.StrengthOfEarthTotem,
		air: AirTotem.WindfuryTotem,
		fire: FireTotem.MagmaTotem,
		water: WaterTotem.ManaSpringTotem,
		useFireElemental: true,
	}),
	maelstromweaponMinStack: 3,
	lightningboltWeave: true,
	autoWeaveDelay: 500,
	delayGcdWeave: 750,
	lavaburstWeave: false,
	firenovaManaThreshold: 3000,
	shamanisticRageManaThreshold: 25,
	primaryShock: PrimaryShock.Earth,
	weaveFlameShock: true,
	rotationType: RotationType.Priority,
	customRotation: CustomRotation.create({
		spells: [
			CustomSpell.create({ spell: CustomRotationSpell.LightningBolt }),
			CustomSpell.create({ spell: CustomRotationSpell.StormstrikeDebuffMissing }),
			CustomSpell.create({ spell: CustomRotationSpell.LightningBoltWeave }),
			CustomSpell.create({ spell: CustomRotationSpell.Stormstrike }),
			CustomSpell.create({ spell: CustomRotationSpell.FlameShock }),
			CustomSpell.create({ spell: CustomRotationSpell.EarthShock }),
			CustomSpell.create({ spell: CustomRotationSpell.MagmaTotem }),
			CustomSpell.create({ spell: CustomRotationSpell.LightningShield }),
			CustomSpell.create({ spell: CustomRotationSpell.FireNova }),
			CustomSpell.create({ spell: CustomRotationSpell.LightningBoltDelayedWeave }),
			CustomSpell.create({ spell: CustomRotationSpell.LavaLash }),
		],
	}),
});

export const ROTATION_FT_DEFAULT = PresetUtils.makePresetAPLRotation('Default FT', DefaultFt);
export const ROTATION_WF_DEFAULT = PresetUtils.makePresetAPLRotation('Default WF', DefaultWf);
export const ROTATION_PHASE_3 = PresetUtils.makePresetAPLRotation('Phase 3', Phase3Apl);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/classic/talent-calc and copy the numbers in the url.
export const StandardTalents = {
	name: 'Standard',
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
	name: 'Phase 3',
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
	bloodlust: true,
	imbueMh: ShamanImbue.WindfuryWeapon,
	imbueOh: ShamanImbue.FlametongueWeapon,
	syncType: ShamanSyncType.Auto,
});

export const DefaultConsumes = Consumes.create({
	flask: Flask.FlaskUnknown,
	food: Food.FoodUnknown,
});

export const DefaultRaidBuffs = RaidBuffs.create({
	giftOfTheWild: TristateEffect.TristateEffectImproved,
	arcaneBrilliance: true,
	leaderOfThePack: true,
	moonkinAura: true,
	divineSpirit: true,
	battleShout: TristateEffect.TristateEffectImproved,
});

export const DefaultDebuffs = Debuffs.create({
	sunderArmor: true,
	curseOfWeakness: TristateEffect.TristateEffectRegular,
	curseOfElements: true,
	faerieFire: TristateEffect.TristateEffectImproved,
	judgementOfWisdom: true,
});

export const OtherDefaults = {
};
