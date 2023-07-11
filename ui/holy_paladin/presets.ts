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
import { SavedTalents } from '../core/proto/ui.js';
import { Player } from '../core/player.js';

import {
	PaladinAura,
	PaladinMajorGlyph,
	PaladinMinorGlyph,
	PaladinJudgement,
	HolyPaladin_Rotation as HolyPaladinRotation,
	HolyPaladin_Options as HolyPaladinOptions,
} from '../core/proto/paladin.js';

import * as Gems from '../core/proto_utils/gems.js';
import * as Tooltips from '../core/constants/tooltips.js';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/wotlk/talent-calc and copy the numbers in the url.

export const StandardTalents = {
	name: 'Standard',
	data: SavedTalents.create({
		talentsString: '50350151020013053100515221-50023131203',
		glyphs: {
			major1: PaladinMajorGlyph.GlyphOfHolyLight,
			major2: PaladinMajorGlyph.GlyphOfSealOfWisdom,
			major3: PaladinMajorGlyph.GlyphOfBeaconOfLight,
			minor2: PaladinMinorGlyph.GlyphOfLayOnHands,
			minor1: PaladinMinorGlyph.GlyphOfSenseUndead,
			minor3: PaladinMinorGlyph.GlyphOfBlessingOfKings
		}
	}),
};

export const DefaultRotation = HolyPaladinRotation.create({
});

export const DefaultOptions = HolyPaladinOptions.create({
	aura: PaladinAura.DevotionAura,
	judgement: PaladinJudgement.NoJudgement,
});

export const DefaultConsumes = Consumes.create({
	defaultPotion: Potions.RunicManaPotion,
	flask: Flask.FlaskOfTheFrostWyrm,
	food: Food.FoodFishFeast,
});

export const PRERAID_PRESET = {
	name: 'Preraid Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecHolyPaladin>) => true,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{"id":44949,"enchant":3819,"gems":[41401,40012]},
		{"id":42647,"gems":[42702]},
		{"id":37673,"enchant":3809,"gems":[40012]},
		{"id":41609,"enchant":3831},
		{"id":39629,"enchant":3832,"gems":[40012,40012]},
		{"id":37788,"enchant":1119,"gems":[0]},
		{"id":39632,"enchant":3604,"gems":[40012,0]},
		{"id":40691,"gems":[40012,40012]},
		{"id":37362,"enchant":3721,"gems":[40012,40012]},
		{"id":44202,"enchant":3606,"gems":[40094]},
		{"id":44283},
		{"id":37694},
		{"id":44255},
		{"id":37111},
		{"id":37169,"enchant":2666},
		{"id":40700,"enchant":1128},
		{"id":40705}
	]}`),
};

export const P1_PRESET = {
	name: 'P1 Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecHolyPaladin>) => true,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{"id":40298,"enchant":3819,"gems":[41401,40012]},
		{"id":44662,"gems":[40012]},
		{"id":40573,"enchant":3809,"gems":[40012]},
		{"id":44005,"enchant":3831,"gems":[40012]},
		{"id":40569,"enchant":3832,"gems":[40012,40012]},
		{"id":40332,"enchant":1119,"gems":[40012,0]},
		{"id":40570,"enchant":3604,"gems":[40012,0]},
		{"id":40259,"gems":[40012]},
		{"id":40572,"enchant":3721,"gems":[40027,40012]},
		{"id":40592,"enchant":3606},
		{"id":40399},
		{"id":40375},
		{"id":44255},
		{"id":37111},
		{"id":40395,"enchant":2666},
		{"id":40401,"enchant":1128},
		{"id":40705}
	]}`),
};

export const P2_PRESET = {
	name: 'P2 Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecHolyPaladin>) => true,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{"id":46180,"enchant":3820,"gems":[41401,40094]},
		{"id":45443,"gems":[40012]},
		{"id":46182,"enchant":3810,"gems":[40012]},
		{"id":45486,"enchant":3831,"gems":[40012]},
		{"id":45445,"enchant":3832,"gems":[42148,42148,42148]},
		{"id":45460,"enchant":1119,"gems":[40012,0]},
		{"id":46179,"enchant":3604,"gems":[40047,0]},
		{"id":45616,"gems":[40012,40012,40012]},
		{"id":46181,"enchant":3721,"gems":[40012,40012]},
		{"id":45537,"enchant":3606,"gems":[40012,40012]},
		{"id":45614,"gems":[45882]},
		{"id":45946,"gems":[40012]},
		{"id":46051},
		{"id":37111},
		{"id":46017,"enchant":2666},
		{"id":45470,"enchant":1128,"gems":[40012]},
		{"id":40705}
	]}`),
};
