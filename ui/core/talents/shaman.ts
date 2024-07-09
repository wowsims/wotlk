import {
  ShamanMajorGlyph,
  ShamanMinorGlyph,
  ShamanTalents,
} from "../proto/shaman.js";
import { GlyphsConfig } from "./glyphs_picker.js";
import { newTalentsConfig,TalentsConfig } from "./talents_picker.js";
import ShamanTalentJson from "./trees/shaman.json";

export const shamanTalentsConfig: TalentsConfig<ShamanTalents> =
  newTalentsConfig(ShamanTalentJson);

export const shamanGlyphsConfig: GlyphsConfig = {
  majorGlyphs: {
    [ShamanMajorGlyph.GlyphOfChainHeal]: {
      name: "治疗链雕文",
      description: "你的治疗链可以多治疗1个目标。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_nature_healingwavegreater.jpg",
    },
    [ShamanMajorGlyph.GlyphOfChainLightning]: {
      name: "闪电链雕文",
      description: "你的闪电链可以多攻击1个目标。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_nature_chainlightning.jpg",
    },
    [ShamanMajorGlyph.GlyphOfEarthShield]: {
      name: "大地之盾雕文",
      description: "你的大地之盾的治疗量提高20%。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_nature_skinofearth.jpg",
    },
    [ShamanMajorGlyph.GlyphOfEarthlivingWeapon]: {
      name: "大地生命武器雕文",
      description: "使你的大地生命武器的触发几率提高5%。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_shaman_earthlivingweapon.jpg",
    },
    [ShamanMajorGlyph.GlyphOfElementalMastery]: {
      name: "元素掌握雕文",
      description: "使你的元素掌握技能的冷却时间缩短30秒。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_nature_wispheal.jpg",
    },
    [ShamanMajorGlyph.GlyphOfFeralSpirit]: {
      name: "野性狼魂雕文",
      description: "你的幽灵狼的攻击强度提高，数值相当于你的攻击强度的30%。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_shaman_feralspirit.jpg",
    },
    [ShamanMajorGlyph.GlyphOfFireElementalTotem]: {
      name: "火焰元素图腾雕文",
      description: "使你的火元素图腾的冷却时间缩短5分钟。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_fire_elemental_totem.jpg",
    },
    [ShamanMajorGlyph.GlyphOfFireNova]: {
      name: "火焰新星雕文",
      description: "使你的火焰新星图腾的冷却时间缩短3秒。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_fire_sealoffire.jpg",
    },
    [ShamanMajorGlyph.GlyphOfFlameShock]: {
      name: "烈焰震击雕文",
      description: "使你的烈焰震击的爆击伤害加成提高60%。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_fire_flameshock.jpg",
    },
    [ShamanMajorGlyph.GlyphOfFlametongueWeapon]: {
      name: "火舌武器雕文",
      description: "使你的每把激活了火舌武器的武器的法术爆击率提高2%。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_fire_flametounge.jpg",
    },
    [ShamanMajorGlyph.GlyphOfFrostShock]: {
      name: "冰霜震击雕文",
      description: "使你的冰霜震击的持续时间延长2秒。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_frost_frostshock.jpg",
    },
    [ShamanMajorGlyph.GlyphOfHealingStreamTotem]: {
      name: "治疗之泉图腾雕文",
      description: "你的治疗之泉图腾的治疗效果提高20%。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/inv_spear_04.jpg",
    },
    [ShamanMajorGlyph.GlyphOfHealingWave]: {
      name: "治疗波雕文",
      description: "当你使用治疗波对其他人进行治疗时，你自己也会受到该次治疗量20%的治疗。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_nature_magicimmunity.jpg",
    },
    [ShamanMajorGlyph.GlyphOfHex]: {
      name: "妖术雕文",
      description: "你的妖术目标在解除妖术效果之前可以承受的伤害提高20%。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_shaman_hex.jpg",
    },
    [ShamanMajorGlyph.GlyphOfLava]: {
      name: "熔岩雕文",
      description: "你的熔岩爆裂得到的法术强度加成效果提高10%。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_shaman_lavaburst.jpg",
    },
    [ShamanMajorGlyph.GlyphOfLavaLash]: {
      name: "熔岩猛击雕文",
      description: "如果你的武器附有火舌效果，则你的熔岩猛击造成的伤害提高10%。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_shaman_lavalash.jpg",
    },
    [ShamanMajorGlyph.GlyphOfLesserHealingWave]: {
      name: "次级治疗波雕文",
      description: "如果目标身上有大地之盾效果，则你的次级治疗波对该目标的治疗效果提高20%。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_nature_healingwavelesser.jpg",
    },
    [ShamanMajorGlyph.GlyphOfLightningBolt]: {
      name: "闪电箭雕文",
      description: "使你的闪电箭造成的伤害提高4%。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_nature_lightning.jpg",
    },
    [ShamanMajorGlyph.GlyphOfLightningShield]: {
      name: "闪电之盾雕文",
      description: "使你的闪电之盾造成的伤害提高20%。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_nature_lightningshield.jpg",
    },
    [ShamanMajorGlyph.GlyphOfManaTide]: {
      name: "法力潮汐图腾雕文",
      description: "你的法力潮汐图腾每次为目标恢复的法力值百分比提高1%。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_nature_healingwavegreater.jpg",
    },
    [ShamanMajorGlyph.GlyphOfRiptide]: {
      name: "激流雕文",
      description: "激流的持续时间延长6秒。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_nature_riptide.jpg",
    },
    [ShamanMajorGlyph.GlyphOfShocking]: {
      name: "震击雕文",
      description: "使你的震击法术的公共冷却时间缩短0.5秒。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_nature_earthshock.jpg",
    },
    [ShamanMajorGlyph.GlyphOfStoneclawTotem]: {
      name: "石爪图腾雕文",
      description: "你的石爪图腾会为你附加一道吸收伤害的护盾，伤害吸收量相当于它为你的图腾提供的护盾强度的4倍。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_nature_stoneclawtotem.jpg",
    },
    [ShamanMajorGlyph.GlyphOfStormstrike]: {
      name: "风暴打击雕文",
      description: "使你的风暴打击的自然伤害加成提高8%。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_shaman_stormstrike.jpg",
    },
    [ShamanMajorGlyph.GlyphOfThunder]: {
      name: "雷霆雕文",
      description: "雷霆风暴的冷却时间缩短10秒。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_shaman_thunderstorm.jpg",
    },
    [ShamanMajorGlyph.GlyphOfTotemOfWrath]: {
      name: "天怒图腾雕文",
      description: "当你施放天怒图腾时，你由该图腾获得的法术强度加成提高30%，持续5 min。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_fire_totemofwrath.jpg",
    },
    [ShamanMajorGlyph.GlyphOfWaterMastery]: {
      name: "水之掌握雕文",
      description: "你的水之护盾的被动法力值恢复效果提高30%。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_shaman_watershield.jpg",
    },
    [ShamanMajorGlyph.GlyphOfWindfuryWeapon]: {
      name: "风怒武器雕文",
      description: "攻击时触发风怒武器的几率提高2%。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_nature_cyclone.jpg",
    },
  },
  minorGlyphs: {
    [ShamanMinorGlyph.GlyphOfAstralRecall]: {
      name: "星界传送雕文",
      description: "你的星界传送法术的冷却时间缩短7.5分钟。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_nature_astralrecal.jpg",
    },
    [ShamanMinorGlyph.GlyphOfGhostWolf]: {
      name: "幽灵狼雕文",
      description: "你在幽灵狼形态下时，每5秒恢复的生命值百分比提高1%。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_nature_spiritwolf.jpg",
    },
    [ShamanMinorGlyph.GlyphOfRenewedLife]: {
      name: "新生雕文",
      description: "你的复生技能不再消耗施法材料。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_nature_reincarnation.jpg",
    },
    [ShamanMinorGlyph.GlyphOfThunderstorm]: {
      name: "雷霆风暴雕文",
      description: "你的雷霆风暴为你回复的法力值提高2%，但不再击退敌人。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_shaman_thunderstorm.jpg",
    },
    [ShamanMinorGlyph.GlyphOfWaterBreathing]: {
      name: "水下呼吸雕文",
      description: "你的水下呼吸法术不再需要施法材料。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_shadow_demonbreath.jpg",
    },
    [ShamanMinorGlyph.GlyphOfWaterShield]: {
      name: "水之护盾雕文",
      description: "使你的水之护盾的使用次数增加1次。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_shaman_watershield.jpg",
    },
    [ShamanMinorGlyph.GlyphOfWaterWalking]: {
      name: "水上行走雕文",
      description: "你的水上行走法术不再需要施法材料。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_frost_windwalkon.jpg",
    },
  },
};
