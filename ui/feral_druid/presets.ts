import { Consumes } from '../core/proto/common.js';
import { Food } from '../core/proto/common.js';
import { EquipmentSpec } from '../core/proto/common.js';
import { Potions } from '../core/proto/common.js';
import { Flask } from '../core/proto/common.js';
import { Glyphs } from '../core/proto/common.js';
import { SavedRotation, SavedTalents } from '../core/proto/ui.js';

import {
	FeralDruid_Rotation as FeralDruidRotation,
	FeralDruid_Options as FeralDruidOptions,
	DruidMajorGlyph,
	DruidMinorGlyph,
	FeralDruid_Rotation_BearweaveType,
	FeralDruid_Rotation_BiteModeType,
	FeralDruid_Rotation_AplType,
} from '../core/proto/druid.js';

import * as Tooltips from '../core/constants/tooltips.js';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/wotlk/talent-calc and copy the numbers in the url.
export const StandardTalents = {
	name: 'Standard',
	data: SavedTalents.create({
		talentsString: '-543202132322010053120030310511-203503012',
		glyphs: Glyphs.create({
			major1: DruidMajorGlyph.GlyphOfOmenOfClarity,
			major2: DruidMajorGlyph.GlyphOfSavageRoar,
			major3: DruidMajorGlyph.GlyphOfShred,
			minor1: DruidMinorGlyph.GlyphOfDash,
			minor2: DruidMinorGlyph.GlyphOfTheWild,
			minor3: DruidMinorGlyph.GlyphOfUnburdenedRebirth,
		}),
	}),
};

export const DefaultRotation = FeralDruidRotation.create({
	rotationType: FeralDruid_Rotation_AplType.SingleTarget,

	bearWeaveType: FeralDruid_Rotation_BearweaveType.None,
	minCombosForRip: 5,
	minCombosForBite: 5,

	useRake: true,
	useBite: true,
	mangleSpam: false,
	biteModeType: FeralDruid_Rotation_BiteModeType.Emperical,
	biteTime: 4.0,
	berserkBiteThresh: 25.0,
	berserkFfThresh: 15.0,
	powerbear: false,
	minRoarOffset: 12.0,
	ripLeeway: 3.0,
	maintainFaerieFire: true,
	hotUptime: 0.0,
	snekWeave: false,
	flowerWeave: false,
	raidTargets: 30,
	maxFfDelay: 0.1,
	prePopOoc: true,
});

export const DefaultOptions = FeralDruidOptions.create({
	latencyMs: 100,
	assumeBleedActive: true,
});

export const DefaultConsumes = Consumes.create({
	flask: Flask.FlaskOfEndlessRage,
	food: Food.FoodDragonfinFilet,
	defaultPotion: Potions.PotionOfSpeed,
});

export const ROTATION_PRESET_LEGACY_DEFAULT = {
	name: 'Legacy Default',
	rotation: SavedRotation.create({
		specRotationOptionsJson: FeralDruidRotation.toJsonString(DefaultRotation),
	}),
}

export const PreRaid_PRESET = {
	name: 'PreRaid',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{"id":42550,"enchant":3817,"gems":[41398,39996]},
		{"id":40678},
		{"id":37139,"enchant":3808,"gems":[39996]},
		{"id":37840,"enchant":3605},
		{"id":37219,"enchant":3832},
		{"id":44203,"enchant":3845,"gems":[0]},
		{"id":37409,"enchant":3604,"gems":[0]},
		{"id":40694,"gems":[49110,39996]},
		{"id":37644,"enchant":3823},
		{"id":44297,"enchant":3606},
		{"id":37642},
		{"id":37624},
		{"id":40684},
		{"id":37166},
		{"id":37883,"enchant":3827},
		{},
		{"id":40713}
  ]}`),
};

export const P1_PRESET = {
	name: 'P1 Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{"id":40473,"enchant":3817,"gems":[41398,39996]},
		{"id":44664,"gems":[39996]},
		{"id":40494,"enchant":3808,"gems":[39996]},
		{"id":40403,"enchant":3605},
		{"id":40539,"enchant":3832,"gems":[39996]},
		{"id":39765,"enchant":3845,"gems":[39996,0]},
		{"id":40541,"enchant":3604,"gems":[0]},
		{"id":40205,"gems":[39996]},
		{"id":44011,"enchant":3823,"gems":[39996,49110]},
		{"id":40243,"enchant":3606,"gems":[40014]},
		{"id":40474},
		{"id":40717},
		{"id":42987},
		{"id":40256},
		{"id":40388,"enchant":3789},
		{},
		{"id":39757}
	]}`),
};


export const P2_PRESET = {
	name: 'P2 Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{"id":46161,"enchant":3817,"gems":[41398,40002]},
		{"id":45517,"gems":[40002]},
		{"id":45245,"enchant":3808,"gems":[40002,40002]},
		{"id":46032,"enchant":3605,"gems":[40002,40058]},
		{"id":45473,"enchant":3832,"gems":[40002,40002,40002]},
		{"id":45869,"enchant":3845,"gems":[40037]},
		{"id":46158,"enchant":3604,"gems":[40002]},
		{"id":46095,"gems":[40002,40002,40002]},
		{"id":45536,"enchant":3823,"gems":[39996,39996,39996]},
		{"id":45564,"enchant":3606,"gems":[39996,39996]},
		{"id":46048,"gems":[45862]},
		{"id":45608,"gems":[39996]},
		{"id":45931},
		{"id":46038},
		{"id":45613,"enchant":3789,"gems":[40037,42702]},
		{},
		{"id":40713}
	]}`),
};

export const P3_PRESET = {
	name: 'P3 Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{"id":48204,"enchant":3817,"gems":[41398,40146]},
		  {"id":47060,"gems":[40117]},
		  {"id":48207,"enchant":3808,"gems":[40146]},
		  {"id":47545,"enchant":3605,"gems":[40117]},
		  {"id":48206,"enchant":3832,"gems":[40117,40146]},
		  {"id":45611,"enchant":3845,"gems":[40146,0]},
		  {"id":47945,"enchant":3604,"gems":[40112,40117,0]},
		  {"id":47112,"gems":[40146,40117,40117]},
		  {"id":48205,"enchant":3823,"gems":[40117,49110]},
		  {"id":47077,"enchant":3606,"gems":[40146,40117]},
		  {"id":47075,"gems":[40117]},
		  {"id":45608,"gems":[40117]},
		  {"id":45931},
		  {"id":47131},
		  {"id":47239,"enchant":3789,"gems":[40146,40117]},
		  {},
		  {"id":47668}
	]}`),
};
