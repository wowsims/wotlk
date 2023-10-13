import {
	Consumes,
	Debuffs,
	EquipmentSpec,
	Flask,
	Food,
	Glyphs,
	IndividualBuffs,
	Potions,
	RaidBuffs,
	Spec,
	TristateEffect,
	UnitReference,
} from '../core/proto/common.js';
import { SavedTalents } from '../core/proto/ui.js';

import {
	SmitePriest_Rotation as Rotation,
	SmitePriest_Options as Options,
	PriestMajorGlyph as MajorGlyph,
	PriestMinorGlyph as MinorGlyph,
} from '../core/proto/priest.js';

import * as PresetUtils from '../core/preset_utils.js';
import * as Tooltips from '../core/constants/tooltips.js';

import DefaultApl from './apls/default.apl.json'

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/wotlk/talent-calc and copy the numbers in the url.
export const StandardTalents = {
	name: 'Standard',
	data: SavedTalents.create({
		talentsString: '05332031013005023310001-005551002020152-00502',
		glyphs: Glyphs.create({
			major1: MajorGlyph.GlyphOfSmite,
			major2: MajorGlyph.GlyphOfShadowWordPain,
			major3: MajorGlyph.GlyphOfShadowWordDeath,
			minor1: MinorGlyph.GlyphOfFortitude,
			minor2: MinorGlyph.GlyphOfShadowfiend,
			minor3: MinorGlyph.GlyphOfFading,
		}),
	}),
};

export const DefaultRotation = Rotation.create({
	useDevouringPlague: true,
	useShadowWordDeath: false,
	useMindBlast: false,
});
export const ROTATION_PRESET_LEGACY_DEFAULT = PresetUtils.makePresetLegacyRotation('Legacy Default', Spec.SpecSmitePriest, DefaultRotation);
export const ROTATION_PRESET_APL = PresetUtils.makePresetAPLRotation('Default', DefaultApl);

export const DefaultOptions = Options.create({
	useInnerFire: true,
	useShadowfiend: true,

	powerInfusionTarget: UnitReference.create(),
});

export const DefaultConsumes = Consumes.create({
	flask: Flask.FlaskOfTheFrostWyrm,
	food: Food.FoodFishFeast,
	defaultPotion: Potions.RunicManaInjector,
	prepopPotion: Potions.PotionOfWildMagic,
});

export const DefaultRaidBuffs = RaidBuffs.create({
	giftOfTheWild: TristateEffect.TristateEffectImproved,
	arcaneBrilliance: true,
	divineSpirit: true,
	totemOfWrath: true,
	moonkinAura: TristateEffect.TristateEffectImproved,
	wrathOfAirTotem: true,
	sanctifiedRetribution: true,
	bloodlust: true,
});

export const DefaultIndividualBuffs = IndividualBuffs.create({
	blessingOfKings: true,
	blessingOfWisdom: TristateEffect.TristateEffectImproved,
	blessingOfMight: TristateEffect.TristateEffectImproved,
	vampiricTouch: true,
});

export const DefaultDebuffs = Debuffs.create({
	faerieFire: TristateEffect.TristateEffectImproved,
	ebonPlaguebringer: true,
	heartOfTheCrusader: true,
	judgementOfWisdom: true,
});

export const PRERAID_PRESET = {
	name: 'Preraid Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{"id":42553,"enchant":3820,"gems":[41333,40014]},
		{"id":40680},
		{"id":34210,"enchant":3810,"gems":[42144,40014]},
		{"id":41610,"enchant":3859},
		{"id":43792,"enchant":1144,"gems":[42144,40049]},
		{"id":37361,"enchant":2332,"gems":[0]},
		{"id":39285,"enchant":3246,"gems":[40014,0]},
		{"id":40696,"gems":[40049,39998]},
		{"id":37854,"enchant":3719},
		{"id":44202,"enchant":3826,"gems":[40026]},
		{"id":43253,"gems":[42144]},
		{"id":39250},
		{"id":37835},
		{"id":37873},
		{"id":41384,"enchant":3834},
		{"id":40698},
		{"id":37177}
	]}`),
};

export const P1_PRESET = {
	name: 'P1 Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{"id":40562,"enchant":3820,"gems":[41333,42144]},
		{"id":44661,"gems":[39998]},
		{"id":40459,"enchant":3810,"gems":[42144]},
		{"id":44005,"enchant":3859,"gems":[42144]},
		{"id":40234,"enchant":1144,"gems":[39998,39998]},
		{"id":44008,"enchant":2332,"gems":[39998,0]},
		{"id":40454,"enchant":3604,"gems":[40049,0]},
		{"id":40561,"enchant":3601,"gems":[39998]},
		{"id":40560,"enchant":3719},
		{"id":40558,"enchant":3826},
		{"id":40719},
		{"id":40399},
		{"id":40255},
		{"id":40432},
		{"id":40395,"enchant":3834},
		{"id":40273},
		{"id":39712}
	]}`),
};
