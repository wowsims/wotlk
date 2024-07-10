import * as PresetUtils from '../core/preset_utils.js';
import {
	Consumes,
	Debuffs,
	Faction,
	Flask,
	Food,
	Glyphs,
	IndividualBuffs,
	PetFood,
	Potions,
	Profession,
	RaidBuffs,
	TristateEffect,
} from '../core/proto/common.js';
import { SavedTalents } from '../core/proto/ui.js';
import {
	Warlock_Options as WarlockOptions,
	Warlock_Options_Armor as Armor,
	Warlock_Options_Summon as Summon,
	Warlock_Options_WeaponImbue as WeaponImbue,
	WarlockMajorGlyph as MajorGlyph,
	WarlockMinorGlyph as MinorGlyph,
} from '../core/proto/warlock.js';

export const BIS_TOOLTIP = '此装备预设参考自Zephan的痛苦术士指南：https://www.warcrafttavern.com/wotlk/guides/pve-affliction-warlock/';

// import SwpGear from './gear_sets/swp.gear.json';
// export const SWP_BIS = PresetUtils.makePresetGear('直接来自SWP', SwpGear);
import PreraidAfflictionGear from './gear_sets/preraid_affliction.gear.json';
export const PRERAID_AFFLICTION_PRESET = PresetUtils.makePresetGear('Preraid痛苦', PreraidAfflictionGear, { tooltip: BIS_TOOLTIP, talentTree: 0 });
import P1AfflictionGear from './gear_sets/p1_affliction.gear.json';
export const P1_AFFLICTION_PRESET = PresetUtils.makePresetGear('P1痛苦', P1AfflictionGear, { tooltip: BIS_TOOLTIP, talentTree: 0 });
import P2AfflictionGear from './gear_sets/p2_affliction.gear.json';
export const P2_AFFLICTION_PRESET = PresetUtils.makePresetGear('P2痛苦', P2AfflictionGear, { tooltip: BIS_TOOLTIP, talentTree: 0 });
import P3AfflictionAllianceGear from './gear_sets/p3_affliction_alliance.gear.json';
export const P3_AFFLICTION_ALLIANCE_PRESET = PresetUtils.makePresetGear('P3痛苦[联盟]', P3AfflictionAllianceGear, { tooltip: BIS_TOOLTIP, talentTree: 0, faction: Faction.Alliance });
import P3AfflictionHordeGear from './gear_sets/p3_affliction_horde.gear.json';
export const P3_AFFLICTION_HORDE_PRESET = PresetUtils.makePresetGear('P3痛苦[部落]', P3AfflictionHordeGear, { tooltip: BIS_TOOLTIP, talentTree: 0, faction: Faction.Horde });
import P4AfflictionGear from './gear_sets/p4_affliction.gear.json';
export const P4_AFFLICTION_PRESET = PresetUtils.makePresetGear('P4痛苦', P4AfflictionGear, { tooltip: BIS_TOOLTIP, talentTree: 0 });
import PreraidDemoDestroGear from './gear_sets/preraid_demodestro.gear.json';
export const PRERAID_DEMODESTRO_PRESET = PresetUtils.makePresetGear('Preraid恶魔/毁灭', PreraidDemoDestroGear, { tooltip: BIS_TOOLTIP, talentTrees: [1,2] });
import P1DemoDestroGear from './gear_sets/p1_demodestro.gear.json';
export const P1_DEMODESTRO_PRESET = PresetUtils.makePresetGear('P1恶魔/毁灭', P1DemoDestroGear, { tooltip: BIS_TOOLTIP, talentTrees: [1,2] });
import P2DemoDestroGear from './gear_sets/p2_demodestro.gear.json';
export const P2_DEMODESTRO_PRESET = PresetUtils.makePresetGear('P2恶魔/毁灭', P2DemoDestroGear, { tooltip: BIS_TOOLTIP, talentTrees: [1,2] });
import P3DemoAllianceGear from './gear_sets/p3_demo_alliance.gear.json';
export const P3_DEMO_ALLIANCE_PRESET = PresetUtils.makePresetGear('P3恶魔[联盟]', P3DemoAllianceGear, { tooltip: BIS_TOOLTIP, talentTree: 1, faction: Faction.Alliance });
import P3DemoHordeGear from './gear_sets/p3_demo_horde.gear.json';
export const P3_DEMO_HORDE_PRESET = PresetUtils.makePresetGear('P3恶魔[部落]', P3DemoHordeGear, { tooltip: BIS_TOOLTIP, talentTree: 1, faction: Faction.Horde });
import P4DemoGear from './gear_sets/p4_demo.gear.json';
export const P4_DEMO_PRESET = PresetUtils.makePresetGear('P4恶魔', P4DemoGear, { tooltip: BIS_TOOLTIP, talentTree: 1 });
import P3DestroAllianceGear from './gear_sets/p3_destro_alliance.gear.json';
export const P3_DESTRO_ALLIANCE_PRESET = PresetUtils.makePresetGear('P3毁灭[联盟]', P3DestroAllianceGear, { tooltip: BIS_TOOLTIP, talentTree: 2, faction: Faction.Alliance });
import P3DestroHordeGear from './gear_sets/p3_destro_horde.gear.json';
export const P3_DESTRO_HORDE_PRESET = PresetUtils.makePresetGear('P3毁灭[部落]', P3DestroHordeGear, { tooltip: BIS_TOOLTIP, talentTree: 2, faction: Faction.Horde });
import P4DestroGear from './gear_sets/p4_destro.gear.json';
export const P4_DESTRO_PRESET = PresetUtils.makePresetGear('P4毁灭', P4DestroGear, { tooltip: BIS_TOOLTIP, talentTree: 2 });

import AfflictionApl from './apls/affliction.apl.json';
export const APL_Affliction_Default = PresetUtils.makePresetAPLRotation('痛苦', AfflictionApl, { talentTree: 0 });
import DemoApl from './apls/demo.apl.json';
export const APL_Demo_Default = PresetUtils.makePresetAPLRotation('恶魔', DemoApl, { talentTree: 1 });
import DestroApl from './apls/destro.apl.json';
export const APL_Destro_Default = PresetUtils.makePresetAPLRotation('毁灭', DestroApl, { talentTree: 2 });

// 默认天赋。使用wowhead计算器格式，访问https://wowhead.com/wotlk/talent-calc并复制URL中的数字。

export const AfflictionTalents = {
	name: '痛苦',
	data: SavedTalents.create({
		talentsString: '2350002030023510253500331151--550000051',
		glyphs: Glyphs.create({
			major1: MajorGlyph.GlyphOfQuickDecay,
			major2: MajorGlyph.GlyphOfLifeTap,
			major3: MajorGlyph.GlyphOfHaunt,
			minor1: MinorGlyph.GlyphOfSouls,
			minor2: MinorGlyph.GlyphOfDrainSoul,
			minor3: MinorGlyph.GlyphOfSubjugateDemon,
		}),
	}),
};

export const DemonologyTalents = {
	name: '恶魔',
	data: SavedTalents.create({
		talentsString: '-203203301035012530135201351-550000052',
		glyphs: Glyphs.create({
			major1: MajorGlyph.GlyphOfLifeTap,
			major2: MajorGlyph.GlyphOfQuickDecay,
			major3: MajorGlyph.GlyphOfFelguard,
			minor1: MinorGlyph.GlyphOfSouls,
			minor2: MinorGlyph.GlyphOfDrainSoul,
			minor3: MinorGlyph.GlyphOfSubjugateDemon,
		}),
	}),
};

export const DestructionTalents = {
	name: '毁灭',
	data: SavedTalents.create({
		talentsString: '-03310030003-05203205210331051335230351',
		glyphs: Glyphs.create({
			major1: MajorGlyph.GlyphOfConflagrate,
			major2: MajorGlyph.GlyphOfLifeTap,
			major3: MajorGlyph.GlyphOfIncinerate,
			minor1: MinorGlyph.GlyphOfSouls,
			minor2: MinorGlyph.GlyphOfDrainSoul,
			minor3: MinorGlyph.GlyphOfSubjugateDemon,
		}),
	}),
};

export const AfflictionOptions = WarlockOptions.create({
	armor: Armor.FelArmor,
	summon: Summon.Felhunter,
	weaponImbue: WeaponImbue.GrandSpellstone,
	detonateSeed: true,
});

export const DemonologyOptions = WarlockOptions.create({
	armor: Armor.FelArmor,
	summon: Summon.Felguard,
	weaponImbue: WeaponImbue.GrandSpellstone,
	detonateSeed: true,
});

export const DestructionOptions = WarlockOptions.create({
	armor: Armor.FelArmor,
	summon: Summon.Imp,
	weaponImbue: WeaponImbue.GrandFirestone,
	detonateSeed: true,
});

export const DefaultConsumes = Consumes.create({
	flask: Flask.FlaskOfTheFrostWyrm,
	food: Food.FoodFishFeast,
	petFood: PetFood.PetFoodSpicedMammothTreats,
	defaultPotion: Potions.PotionOfWildMagic,
	prepopPotion: Potions.PotionOfWildMagic,
});

export const DefaultRaidBuffs = RaidBuffs.create({
	giftOfTheWild: TristateEffect.TristateEffectImproved,
	powerWordFortitude: TristateEffect.TristateEffectImproved,
	strengthOfEarthTotem: TristateEffect.TristateEffectImproved,
	arcaneBrilliance: true,
	divineSpirit: true,
	trueshotAura: true,
	leaderOfThePack: TristateEffect.TristateEffectImproved,
	icyTalons: true,
	totemOfWrath: true,
	moonkinAura: TristateEffect.TristateEffectImproved,
	wrathOfAirTotem: true,
	sanctifiedRetribution: true,
	bloodlust: true,
	demonicPactSp: 500,
});

export const DefaultIndividualBuffs = IndividualBuffs.create({
	blessingOfKings: true,
	blessingOfWisdom: TristateEffect.TristateEffectImproved,
	blessingOfMight: TristateEffect.TristateEffectImproved,
	vampiricTouch: true,
});

export const DestroIndividualBuffs = IndividualBuffs.create({
	blessingOfKings: true,
	blessingOfWisdom: TristateEffect.TristateEffectImproved,
	blessingOfMight: TristateEffect.TristateEffectImproved,
});

export const DefaultDebuffs = Debuffs.create({
	sunderArmor: true,
	faerieFire: TristateEffect.TristateEffectImproved,
	bloodFrenzy: true,
	ebonPlaguebringer: true,
	heartOfTheCrusader: true,
	judgementOfWisdom: true,
});

export const DestroDebuffs = Debuffs.create({
	sunderArmor: true,
	faerieFire: TristateEffect.TristateEffectImproved,
	bloodFrenzy: true,
	ebonPlaguebringer: true,
	heartOfTheCrusader: true,
	judgementOfWisdom: true,
	shadowMastery: true,
});

export const OtherDefaults = {
	distanceFromTarget: 25,
	profession1: Profession.Engineering,
	profession2: Profession.Tailoring,
	channelClipDelay: 150,
	nibelungAverageCasts: 11,
};
