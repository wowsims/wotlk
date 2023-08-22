import { Consumes, Faction } from '../core/proto/common.js';
import { EquipmentSpec } from '../core/proto/common.js';
import { Flask } from '../core/proto/common.js';
import { Food } from '../core/proto/common.js';
import { Glyphs } from '../core/proto/common.js';
import { Potions } from '../core/proto/common.js';
import { Spec } from '../core/proto/common.js';
import { SavedRotation, SavedTalents } from '../core/proto/ui.js';
import { APLRotation } from '../core/proto/apl.js';
import { Player } from '../core/player.js';

import {
	WarriorShout,
	WarriorTalents as WarriorTalents,
	Warrior,
	Warrior_Rotation as WarriorRotation,
	Warrior_Rotation_SunderArmor as SunderArmor,
	Warrior_Options as WarriorOptions,
	WarriorMajorGlyph,
	WarriorMinorGlyph,
	Warrior_Rotation_StanceOption as StanceOption,
	ProtectionWarrior_Rotation_SpellOption as SpellOption,
	Warrior_Rotation_MainGcd as MainGcd,
} from '../core/proto/warrior.js';

import * as Gems from '../core/proto_utils/gems.js';
import * as Tooltips from '../core/constants/tooltips.js';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/wotlk/talent-calc and copy the numbers in the url.
export const ArmsTalents = {
	name: 'Arms',
	data: SavedTalents.create({
		talentsString: '3022032023335100102012213231251-305-2033',
		glyphs: Glyphs.create({
			major1: WarriorMajorGlyph.GlyphOfRending,
			major2: WarriorMajorGlyph.GlyphOfMortalStrike,
			major3: WarriorMajorGlyph.GlyphOfExecution,
			minor1: WarriorMinorGlyph.GlyphOfThunderClap,
			minor2: WarriorMinorGlyph.GlyphOfCommand,
			minor3: WarriorMinorGlyph.GlyphOfShatteringThrow,
		}),
	}),
};

export const FuryTalents = {
	name: 'Fury',
	data: SavedTalents.create({
		talentsString: '32002301233-305053000520310053120500351',
		glyphs: Glyphs.create({
			major1: WarriorMajorGlyph.GlyphOfWhirlwind,
			major2: WarriorMajorGlyph.GlyphOfHeroicStrike,
			major3: WarriorMajorGlyph.GlyphOfExecution,
			minor1: WarriorMinorGlyph.GlyphOfCommand,
			minor2: WarriorMinorGlyph.GlyphOfShatteringThrow,
			minor3: WarriorMinorGlyph.GlyphOfCharge,
		}),
	}),
};

export const DefaultRotation = WarriorRotation.create({
	useRend: true,
	useCleave: false,
	useOverpower: false,
	executePhaseOverpower: false,

	mainGcd: MainGcd.Bloodthirst,

	sunderArmor: SunderArmor.SunderArmorMaintain,

	hsRageThreshold: 12,
	rendHealthThresholdAbove: 20,
	rendRageThresholdBelow: 100,
	slamRageThreshold: 25,
	rendCdThreshold: 0,
	useHsDuringExecute: true,
	useBtDuringExecute: true,
	useWwDuringExecute: true,
	useSlamOverExecute: true,
	stanceOption: StanceOption.DefaultStance,
});

export const ArmsRotation = WarriorRotation.create({
	useRend: true,
	useMs: true,
	useCleave: false,
	sunderArmor: SunderArmor.SunderArmorMaintain,
	msRageThreshold: 30,
	slamRageThreshold: 25,
	hsRageThreshold: 50,
	rendCdThreshold: 0,
	rendHealthThresholdAbove: 100,
	useHsDuringExecute: true,
	spamExecute: true,
	stanceOption: StanceOption.DefaultStance,
});

export const ROTATION_FURY = {
	name: 'Fury',
	rotation: SavedRotation.create({
		specRotationOptionsJson: WarriorRotation.toJsonString(WarriorRotation.create({
		})),
		rotation: APLRotation.fromJsonString(`{
			"enabled": true,
			"prepullActions": [
				{"action":{"castSpell":{"spellId":{"spellId":47436}}},"doAtValue":{"const":{"val":"-10s"}}},
				{"action":{"castSpell":{"spellId":{"otherId":"OtherActionPotion"}}},"doAtValue":{"const":{"val":"-1s"}}}
			],
			"priorityList": [
				{"action":{"castSpell":{"spellId":{"spellId":2458}}}},
				{"action":{"condition":{"auraShouldRefresh":{"sourceUnit":{"type":"Self"},"auraId":{"spellId":47436},"maxOverlap":{"const":{"val":"2s"}}}},"castSpell":{"spellId":{"spellId":47436}}}},
				{"action":{"condition":{"auraShouldRefresh":{"auraId":{"spellId":47467},"maxOverlap":{"const":{"val":"2s"}}}},"castSpell":{"spellId":{"spellId":47467}}}},
				{"hide":true,"action":{"condition":{"and":{"vals":[{"auraShouldRefresh":{"auraId":{"spellId":47502},"maxOverlap":{"const":{"val":"2s"}}}},{"gcdIsReady":{}},{"cmp":{"op":"OpGe","lhs":{"currentRage":{}},"rhs":{"const":{"val":"20"}}}},{"spellIsReady":{"spellId":{"spellId":47502}}}]}},"sequence":{"name":"TC","actions":[{"castSpell":{"spellId":{"spellId":2457}}},{"castSpell":{"spellId":{"spellId":47502}}}]}}},
				{"hide":true,"action":{"resetSequence":{"sequenceName":"TC"}}},
				{"hide":true,"action":{"condition":{"auraShouldRefresh":{"auraId":{"spellId":47437},"maxOverlap":{"const":{"val":"2s"}}}},"castSpell":{"spellId":{"spellId":25203}}}},
				{"action":{"autocastOtherCooldowns":{}}},
				{"action":{"condition":{"auraIsActive":{"auraId":{"spellId":46916}}},"castSpell":{"spellId":{"spellId":47475}}}},
				{"action":{"castSpell":{"spellId":{"spellId":23881}}}},
				{"action":{"castSpell":{"spellId":{"spellId":1680}}}},
				{"action":{"condition":{"and":{"vals":[{"not":{"val":{"isExecutePhase":{"threshold":"E20"}}}},{"not":{"val":{"dotIsActive":{"spellId":{"spellId":47465}}}}}]}},"castSpell":{"spellId":{"spellId":47465}}}},
				{"action":{"castSpell":{"spellId":{"spellId":47471}}}},
				{"action":{"condition":{"cmp":{"op":"OpGe","lhs":{"currentRage":{}},"rhs":{"const":{"val":"30"}}}},"castSpell":{"spellId":{"spellId":47450,"tag":1}}}}
			]
		}`),
	}),
};

export const DefaultOptions = WarriorOptions.create({
	startingRage: 0,
	useRecklessness: true,
	useShatteringThrow: true,
	shout: WarriorShout.WarriorShoutCommanding,
});

export const DefaultConsumes = Consumes.create({
	flask: Flask.FlaskOfEndlessRage,
	food: Food.FoodSpicedWormBurger,
	defaultPotion: Potions.IndestructiblePotion,
	prepopPotion: Potions.PotionOfSpeed,
});

export const P1_PRERAID_FURY_PRESET = {
	name: 'P1 Pre-Raid Fury Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecWarrior>) => player.getTalentTree() != 0,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{"id":41386,"enchant":3817,"gems":[41398,42702]},
		{"id":42645,"gems":[40003]},
		{"id":44195,"enchant":3808},
		{"id":37647,"enchant":3605},
		{"id":39606,"enchant":3832,"gems":[40003,40003]},
		{"id":44203,"enchant":3845,"gems":[0]},
		{"id":39609,"enchant":3604,"gems":[40037,0]},
		{"id":40694,"gems":[42149,42149]},
		{"id":44205,"enchant":3823},
		{"id":44306,"enchant":3606,"gems":[40037,40037]},
		{"id":42642,"gems":[42149]},
		{"id":37642},
		{"id":42987},
		{"id":40684},
		{"id":37852,"enchant":3789},
		{"id":37852,"enchant":3789},
		{"id":37191}
	]}`),
};

export const P1_FURY_PRESET = {
	name: 'P1 Fury Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecWarrior>) => player.getTalentTree() != 0,
	gear: EquipmentSpec.fromJsonString(`{ "items": [
		{"id":44006,"enchant":3817,"gems":[41285,42702]},
		{"id":44664,"gems":[39996]},
		{"id":40530,"enchant":3808,"gems":[40037]},
		{"id":40403,"enchant":3605},
		{"id":40539,"enchant":3832,"gems":[42142]},
		{"id":39765,"enchant":3845,"gems":[39996,0]},
		{"id":40541,"enchant":3604,"gems":[0]},
		{"id":40205,"gems":[42142]},
		{"id":40529,"enchant":3823,"gems":[39996,40022]},
		{"id":40591,"enchant":3606},
		{"id":43993,"gems":[42142]},
		{"id":40717},
		{"id":42987},
		{"id":40256},
		{"id":40384,"enchant":3789},
		{"id":40384,"enchant":3789},
		{"id":40385}
	]}`),
};

export const P2_FURY_PRESET = {
	name: 'P2 Fury Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecWarrior>) => player.getTalentTree() != 0,
	gear: EquipmentSpec.fromJsonString(`{ "items": [
		{"id":46151,"enchant":3817,"gems":[41398,39996]},
		{"id":45517,"gems":[39996]},
		{"id":46149,"enchant":3808,"gems":[39996]},
		{"id":46032,"enchant":3605,"gems":[39996,40037]},
		{"id":46146,"enchant":3832,"gems":[39996,42702]},
		{"id":45611,"enchant":3845,"gems":[40037,0]},
		{"id":46148,"enchant":3604,"gems":[40058]},
		{"id":46095,"gems":[42154,42142,42142]},
		{"id":45536,"enchant":3823,"gems":[39996,39996,39996]},
		{"id":40591,"enchant":3606},
		{"id":45608,"gems":[39996]},
		{"id":45534,"gems":[39996]},
		{"id":42987},
		{"id":45931},
		{"id":45516,"enchant":3789,"gems":[39996,39996]},
		{"id":45516,"enchant":3789,"gems":[39996,39996]},
		{"id":45296,"gems":[39996]}
	]}`),
};

export const P3_FURY_PRESET_ALLIANCE = {
	name: 'P3 Fury Preset [A]',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecWarrior>) => player.getTalentTree() != 0 && player.getFaction() == Faction.Alliance,
	gear: EquipmentSpec.fromJsonString(`{ "items": [
		{"id":48383,"enchant":3817,"gems":[41398,49110]},
		{"id":47060,"gems":[40117]},
		{"id":48381,"enchant":3808,"gems":[40111]},
		{"id":47545,"enchant":3831,"gems":[40111]},
		{"id":48385,"enchant":3832,"gems":[40117,40117]},
		{"id":47074,"enchant":3845,"gems":[40143,0]},
		{"id":47240,"enchant":3604,"gems":[40111,40111,0]},
		{"id":47002,"gems":[40143,40143,40111]},
		{"id":48382,"enchant":3823,"gems":[40142,40117]},
		{"id":47154,"enchant":3606,"gems":[40142,40117]},
		{"id":46966,"gems":[40143]},
		{"id":47075,"gems":[40111]},
		{"id":45931},
		{"id":47131},
		{"id":47078,"enchant":3789,"gems":[42142,40111]},
		{"id":47078,"enchant":3789,"gems":[42142,42142]},
		{"id":46995,"gems":[40111]}
	]}`)
}

export const P3_FURY_PRESET_HORDE = {
	name: 'P3 Fury Preset [H]',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecWarrior>) => player.getTalentTree() != 0 && player.getFaction() == Faction.Horde,
	gear: EquipmentSpec.fromJsonString(`{ "items": [
		{"id":48398,"enchant":3817,"gems":[41398,49110]},
		{"id":47433,"gems":[40111]},
		{"id":48400,"enchant":3808,"gems":[40117]},
		{"id":47546,"enchant":3831,"gems":[40111]},
		{"id":48396,"enchant":3832,"gems":[40111,40111]},
		{"id":47414,"enchant":3845,"gems":[40142,0]},
		{"id":47240,"enchant":3604,"gems":[40111,40111,0]},
		{"id":47429,"gems":[40142,40142,42142]},
		{"id":48399,"enchant":3823,"gems":[40142,40111]},
		{"id":47445,"enchant":3606,"gems":[40142,45862]},
		{"id":47413,"gems":[40142]},
		{"id":47443,"gems":[40111]},
		{"id":45931},
		{"id":47464},
		{"id":47446,"enchant":3789,"gems":[40117,40117]},
		{"id":47446,"enchant":3789,"gems":[42142,42142]},
		{"id":47428,"gems":[40117]}
	]}`)
}

export const P1_PRERAID_ARMS_PRESET = {
	name: 'P1 Pre-Raid Arms Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecWarrior>) => player.getTalentTree() == 0,
	gear: EquipmentSpec.fromJsonString(`{ "items": [
		{"id":41386,"enchant":3817,"gems":[41285,42702]},
		{"id":42645,"gems":[40002]},
		{"id":44195,"enchant":3808},
		{"id":37647,"enchant":3605},
		{"id":39606,"enchant":3832,"gems":[40002,40002]},
		{"id":41355,"enchant":3845,"gems":[0]},
		{"id":39609,"enchant":3604,"gems":[40037,0]},
		{"id":40694,"gems":[42149,42149]},
		{"id":37193,"enchant":3823,"gems":[40002,40037]},
		{"id":44306,"enchant":3606,"gems":[40086,40002]},
		{"id":42642,"gems":[42149]},
		{"id":37642},
		{"id":42987},
		{"id":40684},
		{"id":37852,"enchant":3789},
		{},
		{"id":37191}
	]}`),
};

export const P1_ARMS_PRESET = {
	name: 'P1 Arms Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecWarrior>) => player.getTalentTree() == 0,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{"id":40528,"enchant":3817,"gems":[41398,42153]},
		{"id":44664,"gems":[39996]},
		{"id":40530,"enchant":3808,"gems":[40038]},
		{"id":40403,"enchant":3605},
		{"id":40539,"enchant":3832,"gems":[42153]},
		{"id":40330,"enchant":3845,"gems":[40002,0]},
		{"id":40541,"enchant":3604,"gems":[0]},
		{"id":40205,"gems":[42153]},
		{"id":40318,"enchant":3823,"gems":[49110,40038]},
		{"id":40591,"enchant":3606},
		{"id":43993,"gems":[40002]},
		{"id":40474},
		{"id":42987},
		{"id":40256},
		{"id":40384,"enchant":3789},
		{},
		{"id":40385}
	]}`),
};

export const P2_ARMS_PRESET = {
	name: 'P2 Arms Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecWarrior>) => player.getTalentTree() == 0,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		  {"id":46151,"enchant":3817,"gems":[41398,42142]},
		  {"id":45459,"gems":[40002]},
		  {"id":46149,"enchant":3808,"gems":[42153]},
		  {"id":46032,"enchant":3605,"gems":[39996,40037]},
		  {"id":46146,"enchant":3832,"gems":[42153,39996]},
		  {"id":45663,"enchant":3845,"gems":[39996,0]},
		  {"id":46148,"enchant":3604,"gems":[40037,0]},
		  {"id":45241,"gems":[49110,40037,39996]},
		  {"id":45536,"enchant":3823,"gems":[39996,39996,39996]},
		  {"id":45599,"enchant":3606,"gems":[40002,40002]},
		  {"id":45608,"gems":[39996]},
		  {"id":46322,"gems":[39996]},
		  {"id":46038},
		  {"id":45931},
		  {"id":45533,"enchant":3789,"gems":[39996,39996]},
		  {},
		  {"id":45296,"gems":[39996]}
	]}`),
};

export const P3_ARMS_2P_PRESET_ALLIANCE = {
	name: 'P3 Arms 2p Preset [A]',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecWarrior>) => player.getTalentTree() == 0 && player.getFaction() == Faction.Alliance,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{"id":49478,"enchant":3817,"gems":[41398,40117]},
		{"id":47915,"gems":[40117]},
		{"id":48381,"enchant":3808,"gems":[42142]},
		{"id":47545,"enchant":3605,"gems":[40117]},
		{"id":48385,"enchant":3832,"gems":[42142,42142]},
		{"id":47074,"enchant":3845,"gems":[40143,0]},
		{"id":47240,"enchant":3604,"gems":[49110,40117,0]},
		{"id":47153,"gems":[40143,40117,40117]},
		{"id":47191,"enchant":3823,"gems":[40117,40117,40117]},
		{"id":47077,"enchant":3606,"gems":[40143,40117]},
		{"id":47934,"gems":[40142]},
		{"id":45608,"gems":[40117]},
		{"id":47131},
		{"id":46038},
		{"id":47078,"enchant":3789,"gems":[40117,40117]},
		{},
		{"id":45296,"gems":[40117]}
	]}`),
};

export const P3_ARMS_4P_PRESET_ALLIANCE = {
	name: 'P3 Arms 4p Preset [A]',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecWarrior>) => player.getTalentTree() == 0 && player.getFaction() == Faction.Alliance,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{"id":48383,"enchant":3817,"gems":[41398,40117]},
		{"id":47915,"gems":[40117]},
		{"id":48381,"enchant":3808,"gems":[42153]},
		{"id":47545,"enchant":3605,"gems":[40117]},
		{"id":48385,"enchant":3832,"gems":[42153,42153]},
		{"id":47074,"enchant":3845,"gems":[40117,0]},
		{"id":47240,"enchant":3604,"gems":[49110,40117,0]},
		{"id":47153,"gems":[40117,40117,40117]},
		{"id":48382,"enchant":3823,"gems":[40142,40117]},
		{"id":47077,"enchant":3606,"gems":[40117,40117]},
		{"id":47934,"gems":[40117]},
		{"id":45608,"gems":[40117]},
		{"id":46038},
		{"id":47131},
		{"id":47078,"enchant":3789,"gems":[40117,40117]},
		{},
		{"id":46995,"gems":[40117]}
	]}`),
};

export const P3_ARMS_2P_PRESET_HORDE = {
	name: 'P3 Arms 2p Preset [H]',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecWarrior>) => player.getTalentTree() == 0 && player.getFaction() == Faction.Horde,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{"id":49478,"enchant":3817,"gems":[41398,40111]},
		{"id":45459,"gems":[40111]},
		{"id":48400,"enchant":3808,"gems":[42153]},
		{"id":47546,"enchant":3605,"gems":[40111]},
		{"id":48396,"enchant":3832,"gems":[42153,42153]},
		{"id":47442,"enchant":3845,"gems":[40143,0]},
		{"id":47492,"enchant":3604,"gems":[49110,40117,0]},
		{"id":47472,"gems":[40143,40117,40117]},
		{"id":47480,"enchant":3823,"gems":[40117,40117,40117]},
		{"id":47445,"enchant":3606,"gems":[40143,40117]},
		{"id":48007,"gems":[40143]},
		{"id":45608,"gems":[40117]},
		{"id":47464},
		{"id":46038},
		{"id":47446,"enchant":3789,"gems":[40111,40111]},
		{},
		{"id":47428,"gems":[40111]}
	]}`),
};

export const P3_ARMS_4P_PRESET_HORDE = {
	name: 'P3 Arms 4p Preset [H]',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecWarrior>) => player.getTalentTree() == 0 && player.getFaction() == Faction.Horde,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{"id":48398,"enchant":3817,"gems":[41398,40117]},
		{"id":47988,"gems":[40117]},
		{"id":48400,"enchant":3808,"gems":[42153]},
		{"id":47546,"enchant":3605,"gems":[40117]},
		{"id":48396,"enchant":3832,"gems":[42153,42153]},
		{"id":47442,"enchant":3845,"gems":[40117,0]},
		{"id":47492,"enchant":3604,"gems":[49110,40117,0]},
		{"id":47472,"gems":[40117,40117,40117]},
		{"id":48399,"enchant":3823,"gems":[40142,40117]},
		{"id":47445,"enchant":3606,"gems":[40117,40117]},
		{"id":48007,"gems":[40117]},
		{"id":45608,"gems":[40117]},
		{"id":46038},
		{"id":47464},
		{"id":47446,"enchant":3789,"gems":[40117,40117]},
		{},
		{"id":47428,"gems":[40117]}
	]}`),
};
