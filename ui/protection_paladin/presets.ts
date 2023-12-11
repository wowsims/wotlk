import {
	Consumes,
	CustomRotation,
	CustomSpell,
	Flask,
	Food
} from '../core/proto/common.js';
import { SavedTalents } from '../core/proto/ui.js';

import {
	PaladinAura,
	PaladinJudgement,
	ProtectionPaladin_Options as ProtectionPaladinOptions,
	ProtectionPaladin_Rotation as ProtectionPaladinRotation,
	ProtectionPaladin_Rotation_SpellOption as SpellOption,
} from '../core/proto/paladin.js';

import * as PresetUtils from '../core/preset_utils.js';

import BlankGear from './gear_sets/blank.gear.json';

import DefaultApl from './apls/default.apl.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const DefaultGear = PresetUtils.makePresetGear('Blank', BlankGear);

export const DefaultRotation = ProtectionPaladinRotation.create({
	hammerFirst: false,
	squeezeHolyWrath: true,
	waitSlack: 300,
	useCustomPrio: false,
	customRotation: CustomRotation.create({
		spells: [
			CustomSpell.create({ spell: SpellOption.ShieldOfRighteousness }),
			CustomSpell.create({ spell: SpellOption.HammerOfTheRighteous }),
			CustomSpell.create({ spell: SpellOption.HolyShield }),
			CustomSpell.create({ spell: SpellOption.HammerOfWrath }),
			CustomSpell.create({ spell: SpellOption.Consecration }),
			CustomSpell.create({ spell: SpellOption.AvengersShield }),
			CustomSpell.create({ spell: SpellOption.JudgementOfWisdom }),
			CustomSpell.create({ spell: SpellOption.Exorcism })
		],
	}),
});

export const ROTATION_DEFAULT = PresetUtils.makePresetAPLRotation('Default (969)', DefaultApl);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/classic/talent-calc and copy the numbers in the url.

export const GenericAoeTalents = {
	name: 'Baseline Example',
	data: SavedTalents.create({
		talentsString: '-05005135200132311333312321-511302012003',
	}),
};

export const DefaultOptions = ProtectionPaladinOptions.create({
	aura: PaladinAura.RetributionAura,
	judgement: PaladinJudgement.JudgementOfWisdom,
});

export const DefaultConsumes = Consumes.create({
	flask: Flask.FlaskUnknown,
	food: Food.FoodUnknown,
});
