import {
	Consumes,
	Debuffs,
	Flask,
	Food,
	RaidBuffs,
	TristateEffect
} from '../core/proto/common.js';
import { SavedTalents } from '../core/proto/ui.js';

import {
	EnhancementShaman_Options as EnhancementShamanOptions,
	ShamanImbue,
	ShamanShield,
	ShamanSyncType,
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

export const ROTATION_FT_DEFAULT = PresetUtils.makePresetAPLRotation('Default FT', DefaultFt);
export const ROTATION_WF_DEFAULT = PresetUtils.makePresetAPLRotation('Default WF', DefaultWf);
export const ROTATION_PHASE_3 = PresetUtils.makePresetAPLRotation('Phase 3', Phase3Apl);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/classic/talent-calc and copy the numbers in the url.
export const StandardTalents = {
	name: 'Standard',
	data: SavedTalents.create({
		talentsString: '053030152-30405003105021333031131031051',
	}),
};

export const Phase3Talents = {
	name: 'Phase 3',
	data: SavedTalents.create({
		talentsString: '053030152-30505003105001333031131131051',
	}),
};

export const DefaultOptions = EnhancementShamanOptions.create({
	shield: ShamanShield.LightningShield,
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
	faerieFire: true,
	judgementOfWisdom: true,
});

export const OtherDefaults = {
};
