import {
	Consumes,
	Debuffs,
	EquipmentSpec, Explosive, Faction,
	Flask,
	Food,
	Glyphs,
	IndividualBuffs,
	PartyBuffs,
	Potions,
	RaidBuffs,
	UnitReference, Spec,
	TristateEffect
} from '../core/proto/common.js';
import { SavedRotation, SavedTalents } from '../core/proto/ui.js';

import {
	BalanceDruid_Options as BalanceDruidOptions,
	BalanceDruid_Rotation as BalanceDruidRotation,
	BalanceDruid_Rotation_IsUsage,
	BalanceDruid_Rotation_MfUsage,
	BalanceDruid_Rotation_Type as RotationType,
	BalanceDruid_Rotation_WrathUsage,
	DruidMajorGlyph,
	DruidMinorGlyph,
} from '../core/proto/druid.js';

import * as Tooltips from '../core/constants/tooltips.js';
import { Player } from "../core/player";
import { APLRotation } from '../core/proto/apl.js';

import BasicP3AplJson from './apls/basic_p3.apl.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/wotlk/talent-calc and copy the numbers in the url.
export const Phase1Talents = {
	name: 'Phase 1',
	data: SavedTalents.create({
		talentsString: '5032003115331303213305311231--205003012',
		glyphs: Glyphs.create({
			major1: DruidMajorGlyph.GlyphOfFocus,
			major2: DruidMajorGlyph.GlyphOfInsectSwarm,
			major3: DruidMajorGlyph.GlyphOfStarfall,
			minor1: DruidMinorGlyph.GlyphOfTyphoon,
			minor2: DruidMinorGlyph.GlyphOfUnburdenedRebirth,
			minor3: DruidMinorGlyph.GlyphOfTheWild,
		}),
	}),
};

export const Phase2Talents = {
	name: 'Phase 2',
	data: SavedTalents.create({
		talentsString: '5012203115331303213305311231--205003012',
		glyphs: Glyphs.create({
			major1: DruidMajorGlyph.GlyphOfStarfire,
			major2: DruidMajorGlyph.GlyphOfInsectSwarm,
			major3: DruidMajorGlyph.GlyphOfStarfall,
			minor1: DruidMinorGlyph.GlyphOfTyphoon,
			minor2: DruidMinorGlyph.GlyphOfUnburdenedRebirth,
			minor3: DruidMinorGlyph.GlyphOfTheWild,
		}),
	}),
};

export const Phase3Talents = {
	name: 'Phase 3',
	data: SavedTalents.create({
		talentsString: '5102233115331303213305311031--205003002',
		glyphs: Glyphs.create({
			major1: DruidMajorGlyph.GlyphOfStarfire,
			major2: DruidMajorGlyph.GlyphOfMoonfire,
			major3: DruidMajorGlyph.GlyphOfStarfall,
			minor1: DruidMinorGlyph.GlyphOfTyphoon,
			minor2: DruidMinorGlyph.GlyphOfUnburdenedRebirth,
			minor3: DruidMinorGlyph.GlyphOfTheWild,
		}),
	}),
};

export const DefaultRotation = BalanceDruidRotation.create({
	type: RotationType.Default,
	maintainFaerieFire: true,
	useSmartCooldowns: true,
	mfUsage: BalanceDruid_Rotation_MfUsage.BeforeLunar,
	isUsage: BalanceDruid_Rotation_IsUsage.OptimizeIs,
	wrathUsage: BalanceDruid_Rotation_WrathUsage.RegularWrath,
	useStarfire: true,
	useBattleRes: false,
	playerLatency: 200,
});

export const DefaultOptions = BalanceDruidOptions.create({
	innervateTarget: UnitReference.create(),
});

export const DefaultConsumes = Consumes.create({
	defaultPotion: Potions.PotionOfSpeed,
	flask: Flask.FlaskOfTheFrostWyrm,
	food: Food.FoodFishFeast,
	prepopPotion: Potions.PotionOfWildMagic,
	fillerExplosive: Explosive.ExplosiveSaroniteBomb,
});

export const DefaultRaidBuffs = RaidBuffs.create({
	arcaneBrilliance: true,
	bloodlust: true,
	divineSpirit: true,
	giftOfTheWild: TristateEffect.TristateEffectImproved,
	icyTalons: true,
	moonkinAura: TristateEffect.TristateEffectImproved,
	leaderOfThePack: TristateEffect.TristateEffectImproved,
	powerWordFortitude: TristateEffect.TristateEffectImproved,
	sanctifiedRetribution: true,
	strengthOfEarthTotem: TristateEffect.TristateEffectImproved,
	trueshotAura: true,
	wrathOfAirTotem: true,
	demonicPact: 500,
});

export const DefaultIndividualBuffs = IndividualBuffs.create({
	blessingOfKings: true,
	blessingOfMight: TristateEffect.TristateEffectImproved,
	blessingOfWisdom: TristateEffect.TristateEffectImproved,
	vampiricTouch: true,
});

export const DefaultPartyBuffs = PartyBuffs.create({
	heroicPresence: false,
});

export const DefaultDebuffs = Debuffs.create({
	bloodFrenzy: true,
	ebonPlaguebringer: true,
	faerieFire: TristateEffect.TristateEffectImproved,
	heartOfTheCrusader: true,
	judgementOfWisdom: true,
	shadowMastery: true,
	sunderArmor: true,
	totemOfWrath: true,
});

export const OtherDefaults = {
	distanceFromTarget: 18,
};

export const P2_PRESET = {
	name: 'P2 Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	gear: EquipmentSpec.fromJsonString(` {
      "items": [
        {"id":45497,"enchant":3820,"gems":[41285,42144]},
        {"id":45133,"gems":[40048]},
        {"id":46196,"enchant":3810,"gems":[39998]},
        {"id":45242,"enchant":3859,"gems":[40048]},
        {"id":45519,"enchant":3832,"gems":[40051,42144,40026]},
        {"id":45446,"enchant":2332,"gems":[42144,0]},
        {"id":45665,"enchant":3604,"gems":[39998,39998,0]},
        {"id":45619,"gems":[39998,39998,39998]},
        {"id":46192,"enchant":3719,"gems":[39998,39998]},
        {"id":45537,"enchant":3606,"gems":[39998,40026]},
        {"id":46046,"gems":[39998]},
        {"id":45495,"gems":[39998]},
        {"id":45466},
        {"id":45518},
        {"id":45620,"enchant":3834,"gems":[39998]},
        {"id":45617},
        {"id":40321}
      ]
    }`),
};

export const P1_PRESET = {
	name: 'P1 Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{"id":40467,"enchant":3820,"gems":[41285,42144]},
		{"id":44661,"gems":[40026]},
		{"id":40470,"enchant":3810,"gems":[42144]},
		{"id":44005,"enchant":3859,"gems":[40026]},
		{"id":40469,"enchant":3832,"gems":[42144,39998]},
		{"id":44008,"enchant":2332,"gems":[39998,0]},
		{"id":40466,"enchant":3604,"gems":[39998,0]},
		{"id":40561,"enchant":3601,"gems":[39998]},
		{"id":40560,"enchant":3719},
		{"id":40519,"enchant":3606},
		{"id":40399},
		{"id":40080},
		{"id":40255},
		{"id":40432},
		{"id":40395,"enchant":3834},
		{"id":40192},
		{"id":40321}
	]}`),
};

export const PRE_RAID_PRESET = {
	name: 'Pre-raid Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	gear: EquipmentSpec.fromJsonString(`{ "items": [
		{"id":42554,"enchant":3820,"gems":[41285,40049]},
		{"id":40680},
		{"id":37673,"enchant":3810,"gems":[42144]},
		{"id":41610,"enchant":3859},
		{"id":39547,"enchant":3832,"gems":[42144,40026]},
		{"id":37884,"enchant":2332,"gems":[0]},
		{"id":39544,"enchant":3604,"gems":[42144,0]},
		{"id":40696,"enchant":3601,"gems":[40014,39998]},
		{"id":37854,"enchant":3719},
		{"id":44202,"enchant":3606,"gems":[39998]},
		{"id":40585},
		{"id":43253,"gems":[40026]},
		{"id":37873},
		{"id":40682},
		{"id":45085,"enchant":3834},
		{"id":40698},
		{"id":40712}
	]}`),
};

export const P3_PRESET_HORDE = {
	name: 'P3 Preset [H]',
	enableWhen: (player: Player<Spec.SpecBalanceDruid>) => player.getFaction() == Faction.Horde,
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	gear: EquipmentSpec.fromJsonString(`{"items": [
        {"id":48174,"enchant":3820,"gems":[41285,40155]},
        {"id":47468,"gems":[40155]},
        {"id":48177,"enchant":3810,"gems":[40153]},
        {"id":47551,"enchant":3722,"gems":[40113]},
        {"id":48176,"enchant":3832,"gems":[40113,40113]},
        {"id":47467,"enchant":2332,"gems":[40155,0]},
        {"id":48173,"enchant":3604,"gems":[40113,0]},
        {"id":47447,"gems":[40133,40113,40113]},
        {"id":47479,"enchant":3719,"gems":[40113,40113,40113]},
        {"id":47454,"enchant":3606,"gems":[40133,40113]},
        {"id":47489,"gems":[40155]},
        {"id":46046,"gems":[40155]},
        {"id":45518},
        {"id":47477},
        {"id":47483,"enchant":3834},
        {"id":47437},
        {"id":47670}
      ]
    }`),
};

export const P3_PRESET_ALLI = {
	name: 'P3 Preset [A]',
	enableWhen: (player: Player<Spec.SpecBalanceDruid>) => player.getFaction() == Faction.Alliance,
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	gear: EquipmentSpec.fromJsonString(`{"items": [
        {"id":48171,"enchant":3820,"gems":[41285,40153]},
        {"id":47144,"gems":[40153]},
        {"id":48168,"enchant":3810,"gems":[40153]},
        {"id":47552,"enchant":3722,"gems":[40113]},
        {"id":48169,"enchant":3832,"gems":[40113,40113]},
        {"id":47066,"enchant":2332,"gems":[40113,0]},
        {"id":48172,"enchant":3604,"gems":[40113,0]},
        {"id":47084,"gems":[40133,40113,40113]},
        {"id":47190,"enchant":3719,"gems":[40113,40113,40113]},
        {"id":47097,"enchant":3606,"gems":[40133,40113]},
        {"id":47237,"gems":[40113]},
        {"id":46046,"gems":[40113]},
        {"id":45518},
        {"id":47188},
        {"id":47206,"enchant":3834},
        {"id":47064},
        {"id":47670}
    ]}`),
};


export const ROTATION_PRESET_P3_APL = {
name: 'Basic P3 APL',
rotation: SavedRotation.create({
	specRotationOptionsJson: BalanceDruidRotation.toJsonString(DefaultRotation),
	rotation: APLRotation.fromJsonString(JSON.stringify(BasicP3AplJson)),
}),
};