import {
	Consumes,
	CustomRotation,
	CustomSpell,
	Flask,
	Food,
	Potions,
} from '../core/proto/common.js';
import { SavedTalents } from '../core/proto/ui.js';

import {
	PaladinAura as PaladinAura,
	PaladinMajorGlyph,
	PaladinMinorGlyph,
	PaladinJudgement as PaladinJudgement,
	ProtectionPaladin_Rotation_SpellOption as SpellOption,
	ProtectionPaladin_Rotation as ProtectionPaladinRotation,
	ProtectionPaladin_Options as ProtectionPaladinOptions,
} from '../core/proto/paladin.js';

import * as PresetUtils from '../core/preset_utils.js';

import PreraidGear from './gear_sets/preraid.gear.json';
import P1Gear from './gear_sets/p1.gear.json';
import P2Gear from './gear_sets/p2.gear.json';

import DefaultApl from './apls/default.apl.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const PRERAID_PRESET = PresetUtils.makePresetGear('PreRaid', PreraidGear);
export const P1_PRESET = PresetUtils.makePresetGear('P1 Preset', P1Gear);
export const P2_PRESET = PresetUtils.makePresetGear('P2 Preset', P2Gear);

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
// https://wowhead.com/wotlk/talent-calc and copy the numbers in the url.

export const GenericAoeTalents = {
	name: 'Baseline Example',
	data: SavedTalents.create({
		talentsString: '-05005135200132311333312321-511302012003',
		glyphs: {
			major1: PaladinMajorGlyph.GlyphOfSealOfVengeance,
			major2: PaladinMajorGlyph.GlyphOfRighteousDefense,
			major3: PaladinMajorGlyph.GlyphOfDivinePlea,
			minor1: PaladinMinorGlyph.GlyphOfSenseUndead,
			minor2: PaladinMinorGlyph.GlyphOfLayOnHands,
			minor3: PaladinMinorGlyph.GlyphOfBlessingOfKings
		}
	}),
};

export const DefaultOptions = ProtectionPaladinOptions.create({
	aura: PaladinAura.RetributionAura,
	judgement: PaladinJudgement.JudgementOfWisdom,
});

export const DefaultConsumes = Consumes.create({
	flask: Flask.FlaskOfStoneblood,
	food: Food.FoodDragonfinFilet,
	defaultPotion: Potions.IndestructiblePotion,
	prepopPotion: Potions.IndestructiblePotion,
});
