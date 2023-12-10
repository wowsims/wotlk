import {
	Consumes,
	CustomRotation,
	CustomSpell,
	Flask,
	Food,
	Glyphs,
	Spec
} from '../core/proto/common.js';
import { SavedTalents } from '../core/proto/ui.js';

import {
	PaladinAura,
	PaladinJudgement,
	PaladinMajorGlyph,
	PaladinMinorGlyph,
	RetributionPaladin_Options as RetributionPaladinOptions,
	RetributionPaladin_Rotation as RetributionPaladinRotation,
	RetributionPaladin_Rotation_RotationType as RotationType,
	RetributionPaladin_Rotation_SpellOption as SpellOption,
} from '../core/proto/paladin.js';

import * as PresetUtils from '../core/preset_utils.js';

import BlankGear from './gear_sets/blank.gear.json';

import DefaultApl from './apls/default.apl.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const DefaultGear = PresetUtils.makePresetGear('Blank', BlankGear);

export const DefaultRotation = RetributionPaladinRotation.create({
	type: RotationType.Standard,
	exoSlack: 500,
	consSlack: 500,
	useDivinePlea: true,
	avoidClippingConsecration: true,
	holdLastAvengingWrathUntilExecution: false,
	cancelChaosBane: false,
	divinePleaPercentage: 0.75,
	holyWrathThreshold: 4,
	sovTargets: 1,
	customRotation: CustomRotation.create({
		spells: [
			CustomSpell.create({ spell: SpellOption.HammerOfWrath }),
			CustomSpell.create({ spell: SpellOption.JudgementOfWisdom }),
			CustomSpell.create({ spell: SpellOption.CrusaderStrike }),
			CustomSpell.create({ spell: SpellOption.DivineStorm }),
			CustomSpell.create({ spell: SpellOption.Consecration }),
			CustomSpell.create({ spell: SpellOption.Exorcism }),
			CustomSpell.create({ spell: SpellOption.HolyWrath }),
		],
	}),
	customCastSequence: CustomRotation.create({
		spells: [
			CustomSpell.create({ spell: SpellOption.JudgementOfWisdom }),
			CustomSpell.create({ spell: SpellOption.CrusaderStrike }),
			CustomSpell.create({ spell: SpellOption.DivineStorm }),
			CustomSpell.create({ spell: SpellOption.Consecration }),
			CustomSpell.create({ spell: SpellOption.CrusaderStrike }),
			CustomSpell.create({ spell: SpellOption.Exorcism }),
			CustomSpell.create({ spell: SpellOption.JudgementOfWisdom }),
			CustomSpell.create({ spell: SpellOption.CrusaderStrike }),
			CustomSpell.create({ spell: SpellOption.DivineStorm }),
			CustomSpell.create({ spell: SpellOption.Consecration }),
			CustomSpell.create({ spell: SpellOption.CrusaderStrike }),
		],
	}),
});

export const ROTATION_PRESET_LEGACY_DEFAULT = PresetUtils.makePresetLegacyRotation('Legacy Default', Spec.SpecRetributionPaladin, DefaultRotation);
export const ROTATION_PRESET_DEFAULT = PresetUtils.makePresetAPLRotation('Default', DefaultApl);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/classic/talent-calc and copy the numbers in the url.
export const AuraMasteryTalents = {
	name: 'Aura Mastery',
	data: SavedTalents.create({
		talentsString: '050501-05-05232051203331302133231331',
		glyphs: Glyphs.create({
			major1: PaladinMajorGlyph.GlyphOfSealOfVengeance,
			major2: PaladinMajorGlyph.GlyphOfJudgement,
			major3: PaladinMajorGlyph.GlyphOfReckoning,
			minor1: PaladinMinorGlyph.GlyphOfSenseUndead,
			minor2: PaladinMinorGlyph.GlyphOfLayOnHands,
			minor3: PaladinMinorGlyph.GlyphOfBlessingOfKings
		})
	}),
};


export const DivineSacTalents = {
	name: 'Divine Sacrifice & Guardian',
	data: SavedTalents.create({
		talentsString: '03-453201002-05222051203331302133201331',
		glyphs: Glyphs.create({
			major1: PaladinMajorGlyph.GlyphOfSealOfVengeance,
			major2: PaladinMajorGlyph.GlyphOfJudgement,
			major3: PaladinMajorGlyph.GlyphOfReckoning,
			minor1: PaladinMinorGlyph.GlyphOfSenseUndead,
			minor2: PaladinMinorGlyph.GlyphOfLayOnHands,
			minor3: PaladinMinorGlyph.GlyphOfBlessingOfKings
		})
	}),
};

export const DefaultOptions = RetributionPaladinOptions.create({
	aura: PaladinAura.RetributionAura,
	judgement: PaladinJudgement.JudgementOfWisdom,
});

export const DefaultConsumes = Consumes.create({
	flask: Flask.FlaskUnknown,
	food: Food.FoodUnknown,
});