import { Consumes } from '../core/proto/common.js';
import { BattleElixir } from '../core/proto/common.js';
import { GuardianElixir } from '../core/proto/common.js';
import { Food } from '../core/proto/common.js';
import { EquipmentSpec } from '../core/proto/common.js';
import { Potions } from '../core/proto/common.js';
import { Conjured } from '../core/proto/common.js';
import { Explosive } from '../core/proto/common.js';
import { UnitReference } from '../core/proto/common.js';
import { Glyphs } from '../core/proto/common.js';
import { SavedRotation, SavedTalents } from '../core/proto/ui.js';
import { APLRotation } from '../core/proto/apl.js';

import {
	FeralTankDruid_Rotation as DruidRotation,
	FeralTankDruid_Options as DruidOptions,
	DruidMajorGlyph,
	DruidMinorGlyph,
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
		talentsString: '-503232132322010353120300313511-20350001',
		glyphs: Glyphs.create({
			major1: DruidMajorGlyph.GlyphOfMaul,
			major2: DruidMajorGlyph.GlyphOfSurvivalInstincts,
			major3: DruidMajorGlyph.GlyphOfFrenziedRegeneration,
			minor1: DruidMinorGlyph.GlyphOfChallengingRoar,
			minor2: DruidMinorGlyph.GlyphOfThorns,
			minor3: DruidMinorGlyph.GlyphOfUnburdenedRebirth,
		}),
	}),
};

export const DefaultRotation = DruidRotation.create({
	maulRageThreshold: 25,
	maintainDemoralizingRoar: true,
	lacerateTime: 8.0,
});

export const ROTATION_DEFAULT = {
	name: 'Default',
	rotation: SavedRotation.create({
		specRotationOptionsJson: DruidRotation.toJsonString(DruidRotation.create({
		})),
		rotation: APLRotation.fromJsonString(`{
			"type": "TypeAPL",
			"prepullActions": [
				{"action":{"castSpell":{"spellId":{"otherId":"OtherActionPotion"}}},"doAtValue":{"const":{"val":"-1s"}}}
			],
			"priorityList": [
				{"action":{"autocastOtherCooldowns":{}}},
				{"action":{"condition":{"and":{"vals":[{"cmp":{"op":"OpEq","lhs":{"auraNumStacks":{"sourceUnit":{"type":"CurrentTarget"},"auraId":{"spellId":48568}}},"rhs":{"const":{"val":"5"}}}},{"cmp":{"op":"OpLe","lhs":{"dotRemainingTime":{"spellId":{"spellId":48568}}},"rhs":{"const":{"val":"1.5s"}}}}]}},"castSpell":{"spellId":{"spellId":48568}}}},
				{"action":{"castSpell":{"spellId":{"spellId":48564}}}},
				{"action":{"condition":{"and":{"vals":[{"gcdIsReady":{}},{"not":{"val":{"spellIsReady":{"spellId":{"spellId":48564}}}}},{"cmp":{"op":"OpLe","lhs":{"spellTimeToReady":{"spellId":{"spellId":48564}}},"rhs":{"const":{"val":"1.2s"}}}}]}},"wait":{"duration":{"spellTimeToReady":{"spellId":{"spellId":48564}}}}}},
				{"action":{"condition":{"auraShouldRefresh":{"auraId":{"spellId":48560},"maxOverlap":{"const":{"val":"1.5s"}}}},"castSpell":{"spellId":{"spellId":48560}}}},
				{"action":{"castSpell":{"spellId":{"spellId":16857}}}},
				{"action":{"condition":{"or":{"vals":[{"cmp":{"op":"OpLt","lhs":{"auraNumStacks":{"sourceUnit":{"type":"CurrentTarget"},"auraId":{"spellId":48568}}},"rhs":{"const":{"val":"5"}}}},{"cmp":{"op":"OpLe","lhs":{"dotRemainingTime":{"spellId":{"spellId":48568}}},"rhs":{"const":{"val":"8s"}}}}]}},"castSpell":{"spellId":{"spellId":48568}}}},
				{"action":{"condition":{"cmp":{"op":"OpGe","lhs":{"currentRage":{}},"rhs":{"const":{"val":"40"}}}},"castSpell":{"spellId":{"spellId":48562}}}},
				{"action":{"condition":{"cmp":{"op":"OpGe","lhs":{"currentRage":{}},"rhs":{"const":{"val":"25"}}}},"castSpell":{"spellId":{"spellId":48480,"tag":1}}}}
			]
		}`),
	}),
};

export const DefaultOptions = DruidOptions.create({
	innervateTarget: UnitReference.create(),
	startingRage: 20,
});

export const DefaultConsumes = Consumes.create({
	battleElixir: BattleElixir.GurusElixir,
	guardianElixir: GuardianElixir.GiftOfArthas,
	food: Food.FoodBlackenedDragonfin,
	prepopPotion: Potions.IndestructiblePotion,
	defaultPotion: Potions.IndestructiblePotion,
	defaultConjured: Conjured.ConjuredHealthstone,
	thermalSapper: true,
	fillerExplosive: Explosive.ExplosiveSaroniteBomb,
});

export const P1_PRESET = {
	name: 'P1 Boss Tanking',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{"id":40329,"enchant":67839,"gems":[41339,40008]},
		{"id":40387},
		{"id":40494,"enchant":44957,"gems":[40008]},
		{"id":40252,"enchant":3294},
		{"id":40471,"enchant":3832,"gems":[42702,40088]},
		{"id":40186,"enchant":3850,"gems":[40008,0]},
		{"id":40472,"enchant":63770,"gems":[40008,0]},
		{"id":43591,"gems":[40008,40008,40008]},
		{"id":44011,"enchant":38373,"gems":[40008,40008]},
		{"id":40243,"enchant":55016,"gems":[40008]},
		{"id":40370},
		{"id":37784},
		{"id":44253},
		{"id":37220},
		{"id":40280,"enchant":2673},
		{},
		{"id":38365}
	]}`),
};

export const P2_PRESET = {
	name: 'P2 Boss Tanking',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{"id":41678,"enchant":67839,"gems":[41339,45880]},
		{"id":45517,"gems":[40008]},
		{"id":45245,"enchant":44957,"gems":[40008,40008]},
		{"id":45496,"enchant":3294,"gems":[42702]},
		{"id":45473,"enchant":3832,"gems":[40008,40008,40008]},
		{"id":45611,"enchant":3850,"gems":[40008,0]},
		{"id":46043,"enchant":63770,"gems":[40008,40008,0]},
		{"id":46095,"gems":[40008,40008,40008]},
		{"id":45536,"enchant":38373,"gems":[40008,40008,40008]},
		{"id":45232,"enchant":55016,"gems":[40008]},
		{"id":45471,"gems":[40091]},
		{"id":45608,"gems":[40008]},
		{"id":45158},
		{"id":46021},
		{"id":45533,"enchant":3870,"gems":[40008,40008]},
		{},
		{"id":45509}
	]}`),
};
