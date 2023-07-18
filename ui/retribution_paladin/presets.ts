import { Conjured, Consumes } from '../core/proto/common.js';
import { CustomRotation, CustomSpell } from '../core/proto/common.js';
import { EquipmentSpec } from '../core/proto/common.js';
import { Flask } from '../core/proto/common.js';
import { Food } from '../core/proto/common.js';
import { Glyphs } from '../core/proto/common.js';
import { ItemSpec } from '../core/proto/common.js';
import { Potions } from '../core/proto/common.js';
import { Spec } from '../core/proto/common.js';
import { Faction } from '../core/proto/common.js';
import { SavedRotation, SavedTalents } from '../core/proto/ui.js';
import { Player } from '../core/player.js';
import { APLRotation } from '../core/proto/apl.js';

import {
	PaladinAura as PaladinAura,
	PaladinJudgement as PaladinJudgement,
	RetributionPaladin_Rotation as RetributionPaladinRotation,
	RetributionPaladin_Options as RetributionPaladinOptions,
	RetributionPaladin_Rotation_RotationType as RotationType,
	RetributionPaladin_Rotation_SpellOption as SpellOption,
	PaladinMajorGlyph,
	PaladinMinorGlyph,
} from '../core/proto/paladin.js';

import * as Gems from '../core/proto_utils/gems.js';
import * as Tooltips from '../core/constants/tooltips.js';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/wotlk/talent-calc and copy the numbers in the url.
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

export const DefaultOptions = RetributionPaladinOptions.create({
	aura: PaladinAura.RetributionAura,
	judgement: PaladinJudgement.JudgementOfWisdom,
});

export const DefaultConsumes = Consumes.create({
	defaultPotion: Potions.PotionOfSpeed,
	defaultConjured: Conjured.ConjuredDarkRune,
	flask: Flask.FlaskOfEndlessRage,
	food: Food.FoodDragonfinFilet,
});

// Maybe use this later if I can figure out the interactive tooltips from tippy
const RET_BIS_DISCLAIMER = "<p>Please reference <a target=\"_blank\" href=\"https://docs.google.com/spreadsheets/d/1SxO6abYm4k7XRaP1MsxhaqYoukgyZ-cbWDE3ujadjx4/\">Baranor's TBC BiS Lists</a> for more detailed gearing options and information.</p>"

export const PRE_RAID_PRESET = {
	name: 'Pre-Raid Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecRetributionPaladin>) => true,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{"id":41386,"enchant":3817,"gems":[41398,40022]},
		{"id":40678},
		{"id":34388,"enchant":3875,"gems":[39996,40037]},
		{"id":37647,"enchant":3605},
		{"id":39633,"enchant":3832,"gems":[39996,40038]},
		{"id":41355,"enchant":3845,"gems":[0]},
		{"id":39634,"enchant":3604,"gems":[39996,0]},
		{"id":40694,"gems":[39996,39996]},
		{"id":37193,"enchant":3326,"gems":[39996,39996]},
		{"id":44297,"enchant":3606},
		{"id":40586},
		{"id":37685},
		{"id":42987},
		{"id":40684},
		{"id":37852,"enchant":3789},
		{},
		{"id":37574}
	]}`),
};

export const P1_PRESET = {
	name: 'P1 Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecRetributionPaladin>) => true,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		  {"id":44006,"enchant":3817,"gems":[41398,49110]},
		  {"id":44664,"gems":[42142]},
		  {"id":40578,"enchant":3808,"gems":[39996]},
		  {"id":40403,"enchant":3605},
		  {"id":40574,"enchant":3832,"gems":[42142,39996]},
		  {"id":40330,"enchant":3845,"gems":[39996,0]},
		  {"id":40541,"enchant":3604,"gems":[0]},
		  {"id":40278,"gems":[39996,39996]},
		  {"id":44011,"enchant":3823,"gems":[42142,39996]},
		  {"id":40591,"enchant":3606},
		  {"id":40075},
		  {"id":40474},
		  {"id":42987},
		  {"id":40431},
		  {"id":40384,"enchant":3789},
		  {},
		  {"id":42852}
		]}`),
};

export const P2_PRESET = {
	name: 'P2 Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecRetributionPaladin>) => true,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{"id":45472,"enchant":3817,"gems":[41398,42702]},
		{"id":45517,"gems":[39996]},
		{"id":45245,"enchant":3808,"gems":[39996,39996]},
		{"id":45461,"enchant":3605,"gems":[39996]},
		{"id":45473,"enchant":3832,"gems":[39996,39996,39996]},
		{"id":45663,"enchant":3845,"gems":[39996,0]},
		{"id":45444,"enchant":3604,"gems":[39996,39996,0]},
		{"id":46095,"gems":[42142,42142,42142]},
		{"id":45134,"enchant":3823,"gems":[39996,39996,39996]},
		{"id":45599,"enchant":3606,"gems":[39996,39996]},
		{"id":45608,"gems":[39996]},
		{"id":45534,"gems":[39996]},
		{"id":45609},
		{"id":42987},
		{"id":45516,"enchant":3789,"gems":[39996,39996]},
		{},
		{"id":42853}
	]}`),
};

export const P3_PRESET = {
	name: 'P3 Mace Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecRetributionPaladin>) => true,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{"id":48614,"enchant":3817,"gems":[41398,40142]},
        	{"id":47110,"gems":[40142]},
        	{"id":48612,"enchant":3808,"gems":[40111]},
        	{"id":47548,"enchant":3605,"gems":[40111]},
        	{"id":48616,"enchant":3832,"gems":[49110,40111]},
        	{"id":47155,"enchant":3845,"gems":[40111,40111,0]},
        	{"id":48615,"enchant":3604,"gems":[40142,0]},
        	{"id":47002,"gems":[40111,40111,40111]},
        	{"id":47132,"enchant":3823,"gems":[42142,42142,42142]},
        	{"id":47154,"enchant":3606,"gems":[40142,40111]},
        	{"id":47075,"gems":[40111]},
        	{"id":46966,"gems":[40142]},
        	{"id":47131},
        	{"id":42987},
        	{"id":47520,"enchant":3789,"gems":[40111,40111]},
        	{},
        	{"id":47661}
	]}`),
};

export const P4_PRESET = {
	name: 'P4 Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecRetributionPaladin>) => true,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{"id":51277,"enchant":3817,"gems":[41398,40118]},
        	{"id":50633,"gems":[40111]},
        	{"id":51279,"enchant":3808,"gems":[40111]},
        	{"id":50653,"enchant":3605,"gems":[40111]},
        	{"id":51275,"enchant":3832,"gems":[40118,49110]},
        	{"id":50659,"enchant":3845,"gems":[42142,0]},
        	{"id":50690,"enchant":3604,"gems":[40111,40111,0]},
        	{"id":50707,"gems":[40111,40111,45862]},
        	{"id":51278,"enchant":3823,"gems":[42142,42142]},
        	{"id":50607,"enchant":3606,"gems":[40111,40111]},
        	{"id":50604,"gems":[40111]},
        	{"id":50402,"gems":[40111]},
        	{"id":50706},
        	{"id":47131},
        	{"id":49623,"enchant":3789,"gems":[40111,40111,40111]},
        	{},
        	{"id":50455}
	]}`),
};

export const P5_PRESET = {
	name: 'P5 Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecRetributionPaladin>) => true,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{"id":51277,"enchant":3817,"gems":[41398,40111]},
        	{"id":54581,"gems":[40146]},
        	{"id":51279,"enchant":3808,"gems":[40111]},
        	{"id":50677,"enchant":3605,"gems":[40146]},
        	{"id":51275,"enchant":3832,"gems":[40111,49110]},
        	{"id":54580,"enchant":3845,"gems":[40111,0]},
        	{"id":50690,"enchant":3604,"gems":[40146,40111,0]},
        	{"id":50707,"gems":[40111,40111,40111]},
        	{"id":51278,"enchant":3823,"gems":[40111,40111]},
        	{"id":54578,"enchant":3606,"gems":[40111,40111]},
        	{"id":50402,"gems":[40111]},
        	{"id":54576,"gems":[40111]},
        	{"id":54590},
        	{"id":50706},
        	{"id":49623,"enchant":3789,"gems":[42142,42142,42154]},
        	{},
        	{"id":50455}
	]}`),
};

export const ROTATION_PRESET_BASIC_APL = {
	name: 'Basic APL',
	rotation: SavedRotation.create({
		specRotationOptionsJson: RetributionPaladinRotation.toJsonString(DefaultRotation),
		rotation: APLRotation.fromJsonString(`{
			"enabled": true,
			"prepullActions": [
			  {"action":{"castSpell":{"spellId":{"otherId":"OtherActionPotion"}}},"doAt":"-1s"}
			],
			"priorityList": [
			  {"action":{"autocastOtherCooldowns":{}}},
			  {"action":{"castSpell":{"spellId":{"spellId":67485}}}},
			  {"action":{"castSpell":{"spellId":{"spellId":48806}}}},
			  {"action":{"castSpell":{"spellId":{"spellId":53408}}}},
			  {"action":{"castSpell":{"spellId":{"spellId":35395}}}},
			  {"action":{"castSpell":{"spellId":{"spellId":53385}}}},
			  {"action":{"condition":{"auraIsActive":{"auraId":{"spellId":53488}}},"castSpell":{"spellId":{"spellId":48801}}}},
			  {"action":{"condition":{"cmp":{"op":"OpGt","lhs":{"remainingTime":{}},"rhs":{"const":{"val":"4s"}}}},"castSpell":{"spellId":{"spellId":48819}}}}
			]
		  }`),
	}),
	};