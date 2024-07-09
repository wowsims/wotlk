import {
  WarlockTalents,
  WarlockMajorGlyph,
  WarlockMinorGlyph,
} from "../proto/warlock.js";

import { GlyphsConfig } from "./glyphs_picker.js";
import { TalentsConfig, newTalentsConfig } from "./talents_picker.js";

import WarlockTalentJson from "./trees/warlock.json";

export const warlockTalentsConfig: TalentsConfig<WarlockTalents> =
  newTalentsConfig(WarlockTalentJson);

export const warlockGlyphsConfig: GlyphsConfig = {
  majorGlyphs: {
    [WarlockMajorGlyph.GlyphOfChaosBolt]: {
      name: "混乱之箭雕文",
      description: "混乱之箭的冷却时间缩短2秒。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_warlock_chaosbolt.jpg",
    },
    [WarlockMajorGlyph.GlyphOfConflagrate]: {
      name: "燃烧雕文",
      description: "你的燃烧法术不再消耗你施放在目标身上的献祭或暗影烈焰法术。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_fire_fireball.jpg",
    },
    [WarlockMajorGlyph.GlyphOfCorruption]: {
      name: "腐蚀雕文",
      description: "你的腐蚀术在对目标造成伤害之后有4%的几率使你进入暗影冥思状态，使得你的下一次暗影箭的施法时间缩短100%。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_shadow_abominationexplosion.jpg",
    },
    [WarlockMajorGlyph.GlyphOfCurseOfAgony]: {
      name: "痛苦诅咒雕文",
      description: "使你的痛苦诅咒的持续时间延长4秒。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_shadow_curseofsargeras.jpg",
    },
    [WarlockMajorGlyph.GlyphOfDeathCoil]: {
      name: "死亡缠绕雕文",
      description: "使你的死亡缠绕的持续时间延长0.5秒。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_shadow_deathcoil.jpg",
    },
    [WarlockMajorGlyph.GlyphOfDemonicCircle]: {
      name: "恶魔法阵雕文",
      description: "恶魔法阵的冷却时间缩短4秒。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_shadow_demoniccirclesummon.jpg",
    },
    [WarlockMajorGlyph.GlyphOfFear]: {
      name: "恐惧雕文",
      description: "使你的恐惧术的目标在解除恐惧效果之前所能承受的伤害提高20%。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_shadow_possession.jpg",
    },
    [WarlockMajorGlyph.GlyphOfFelguard]: {
      name: "恶魔卫士雕文",
      description: "使你的恶魔卫士的攻击强度提高20%。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_shadow_summonfelguard.jpg",
    },
    [WarlockMajorGlyph.GlyphOfFelhunter]: {
      name: "地狱猎犬雕文",
      description: "当你的地狱猎犬使用吞噬魔法时，你也会获得等量的治疗。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_shadow_summonfelhunter.jpg",
    },
    [WarlockMajorGlyph.GlyphOfHaunt]: {
      name: "鬼影缠身雕文",
      description: "你的鬼影缠身法术提供的伤害加成效果提高3%。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_warlock_haunt.jpg",
    },
    [WarlockMajorGlyph.GlyphOfHealthFunnel]: {
      name: "生命通道雕文",
      description: "使你在引导生命通道法术时受到施法打退效果的几率提高100%。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_shadow_lifedrain.jpg",
    },
    [WarlockMajorGlyph.GlyphOfHealthstone]: {
      name: "治疗石雕文",
      description: "你通过治疗石恢复的生命值提高30%。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/inv_stone_04.jpg",
    },
    [WarlockMajorGlyph.GlyphOfHowlOfTerror]: {
      name: "恐惧嚎叫雕文",
      description: "使你的恐惧嚎叫的冷却时间缩短8秒。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_shadow_deathscream.jpg",
    },
    [WarlockMajorGlyph.GlyphOfImmolate]: {
      name: "献祭雕文",
      description: "使你的献祭的持续伤害提高10%。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_fire_immolation.jpg",
    },
    [WarlockMajorGlyph.GlyphOfImp]: {
      name: "小鬼雕文",
      description: "使你的小鬼的火焰箭造成的伤害提高20%。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_shadow_summonimp.jpg",
    },
    [WarlockMajorGlyph.GlyphOfIncinerate]: {
      name: "烧尽雕文",
      description: "使你的烧尽法术造成的伤害提高5%。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_fire_burnout.jpg",
    },
    [WarlockMajorGlyph.GlyphOfLifeTap]: {
      name: "生命分流雕文",
      description: "当你使用生命分流时，你的法术强度提高，数值相当于你的精神值的20%，持续40 sec。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_shadow_burningspirit.jpg",
    },
    [WarlockMajorGlyph.GlyphOfMetamorphosis]: {
      name: "恶魔变形雕文",
      description: "使你的恶魔变形的持续时间延长6秒。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_shadow_demonform.jpg",
    },
    [WarlockMajorGlyph.GlyphOfQuickDecay]: {
      name: "急速凋零雕文",
      description: "你的急速效果现在可以缩短腐蚀术造成的周期性伤害间隔。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_shadow_abominationexplosion.jpg",
    },
    [WarlockMajorGlyph.GlyphOfSearingPain]: {
      name: "灼热之痛雕文",
      description: "使你的灼热之痛的爆击伤害加成提高20%。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_fire_soulburn.jpg",
    },
    [WarlockMajorGlyph.GlyphOfShadowBolt]: {
      name: "暗影箭雕文",
      description: "使你的暗影箭消耗的法力值降低10%。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_shadow_shadowbolt.jpg",
    },
    [WarlockMajorGlyph.GlyphOfShadowburn]: {
      name: "暗影灼烧雕文",
      description: "当目标的生命值低于35%时，使你的暗影灼烧的爆击几率提高20%。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_shadow_scourgebuild.jpg",
    },
    [WarlockMajorGlyph.GlyphOfShadowflame]: {
      name: "暗影烈焰雕文",
      description: "你的暗影烈焰可以使目标的移动速度降低70%。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_warlock_shadowflame.jpg",
    },
    [WarlockMajorGlyph.GlyphOfSiphonLife]: {
      name: "生命虹吸雕文",
      description: "使你的生命虹吸天赋的治疗效果提高25%。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_shadow_requiem.jpg",
    },
    [WarlockMajorGlyph.GlyphOfSoulLink]: {
      name: "灵魂链接雕文",
      description: "使你的灵魂链接分担伤害的百分比提高5%。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_shadow_gathershadows.jpg",
    },
    [WarlockMajorGlyph.GlyphOfSoulstone]: {
      name: "灵魂石雕文",
      description: "你通过灵魂石复活后的起始生命值提高300%。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_shadow_soulgem.jpg",
    },
    [WarlockMajorGlyph.GlyphOfSuccubus]: {
      name: "魅魔雕文",
      description: "你的魅魔的魅惑技能可以移除目标身上的所有持续伤害效果。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_shadow_summonsuccubus.jpg",
    },
    [WarlockMajorGlyph.GlyphOfUnstableAffliction]: {
      name: "痛苦无常雕文",
      description: "使你的痛苦无常的施法时间缩短0.2秒。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_shadow_unstableaffliction_3.jpg",
    },
    [WarlockMajorGlyph.GlyphOfVoidwalker]: {
      name: "虚空行者雕文",
      description: "使你的虚空行者的耐力总值提高20%。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_shadow_summonvoidwalker.jpg",
    },
  },
  minorGlyphs: {
    [WarlockMinorGlyph.GlyphOfCurseOfExhausion]: {
      name: "疲劳诅咒雕文",
      description: "使你的疲劳诅咒的射程延长5码。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_shadow_grimward.jpg",
    },
    [WarlockMinorGlyph.GlyphOfDrainSoul]: {
      name: "吸取灵魂雕文",
      description: "你的吸取灵魂技能有一定几率制造一块额外的灵魂碎片。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_shadow_haunting.jpg",
    },
    [WarlockMinorGlyph.GlyphOfSubjugateDemon]: {
      name: "征服恶魔雕文",
      description: "使你的征服恶魔法术的施法时间缩短50%。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_shadow_enslavedemon.jpg",
    },
    [WarlockMinorGlyph.GlyphOfKilrogg]: {
      name: "基尔罗格雕文",
      description: "使你的基尔罗格之眼的移动速度提高50%，并使其可以在允许使用飞行坐骑的地方飞行。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_shadow_evileye.jpg",
    },
    [WarlockMinorGlyph.GlyphOfSouls]: {
      name: "灵魂雕文",
      description: "你的灵魂仪式法术消耗的法力值降低70%。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_shadow_shadesofdarkness.jpg",
    },
    [WarlockMinorGlyph.GlyphOfUnendingBreath]: {
      name: "水下呼吸雕文",
      description: "受到你的魔息术影响的目标游泳速度提高20%。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_shadow_demonbreath.jpg",
    },
  },
};
