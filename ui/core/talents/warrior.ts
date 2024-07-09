import {
  WarriorTalents,
  WarriorMajorGlyph,
  WarriorMinorGlyph,
} from "../proto/warrior.js";

import { GlyphsConfig } from "./glyphs_picker.js";
import { TalentsConfig, newTalentsConfig } from "./talents_picker.js";

import WarriorTalentJson from "./trees/warrior.json";

export const warriorTalentsConfig: TalentsConfig<WarriorTalents> =
  newTalentsConfig(WarriorTalentJson);

export const warriorGlyphsConfig: GlyphsConfig = {
  majorGlyphs: {
    [WarriorMajorGlyph.GlyphOfBarbaricInsults]: {
      name: "野蛮侵犯雕文",
      description: "你的惩戒痛击技能造成的威胁值提高100%。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_warrior_punishingblow.jpg",
    },
    [WarriorMajorGlyph.GlyphOfBladestorm]: {
      name: "利刃风暴雕文",
      description: "利刃风暴的冷却时间缩短15秒。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_warrior_bladestorm.jpg",
    },
    [WarriorMajorGlyph.GlyphOfBlocking]: {
      name: "格挡雕文",
      description: "使你在使用盾牌猛击技能后的10 sec内格挡值提高10%。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/inv_shield_05.jpg",
    },
    [WarriorMajorGlyph.GlyphOfBloodthirst]: {
      name: "嗜血雕文",
      description: "使你的嗜血技能恢复的生命值提高100%。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_nature_bloodlust.jpg",
    },
    [WarriorMajorGlyph.GlyphOfCleaving]: {
      name: "顺劈斩雕文",
      description: "使你的顺劈斩的目标数量提高1个。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_warrior_cleave.jpg",
    },
    [WarriorMajorGlyph.GlyphOfDevastate]: {
      name: "毁灭打击雕文",
      description: "你的毁灭打击技能附加两层破甲效果。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/inv_sword_11.jpg",
    },
    [WarriorMajorGlyph.GlyphOfEnragedRegeneration]: {
      name: "狂怒回复雕文",
      description: "你的狂怒回复技能为你回复生命值的效果提高10%。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_warrior_focusedrage.jpg",
    },
    [WarriorMajorGlyph.GlyphOfExecution]: {
      name: "斩杀雕文",
      description: "你的斩杀技能造成的伤害按照你拥有10点额外怒气值的情况计算。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/inv_sword_48.jpg",
    },
    [WarriorMajorGlyph.GlyphOfHamstring]: {
      name: "断筋雕文",
      description: "使你的断筋技能有10%的几率令目标无法移动，持续5 sec。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_shockwave.jpg",
    },
    [WarriorMajorGlyph.GlyphOfHeroicStrike]: {
      name: "英勇打击雕文",
      description: "当你的英勇打击技能爆击时，你可以获得10点怒气值。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_rogue_ambush.jpg",
    },
    [WarriorMajorGlyph.GlyphOfIntervene]: {
      name: "援护雕文",
      description: "使你帮助援护目标承受攻击的次数增加1次。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_warrior_victoryrush.jpg",
    },
    [WarriorMajorGlyph.GlyphOfLastStand]: {
      name: "破釜沉舟雕文",
      description: "使你的破釜沉舟技能的冷却时间缩短1分钟。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_holy_ashestoashes.jpg",
    },
    [WarriorMajorGlyph.GlyphOfMortalStrike]: {
      name: "致死打击雕文",
      description: "使你的致死打击造成的伤害提高10%。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_warrior_savageblow.jpg",
    },
    [WarriorMajorGlyph.GlyphOfOverpower]: {
      name: "压制雕文",
      description: "当你的攻击被招架时，有100%的几率使你可以发动压制。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_meleedamage.jpg",
    },
    [WarriorMajorGlyph.GlyphOfRapidCharge]: {
      name: "疾速冲锋雕文",
      description: "你的冲锋技能的冷却时间缩短7%。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_warrior_charge.jpg",
    },
    [WarriorMajorGlyph.GlyphOfRending]: {
      name: "撕裂雕文",
      description: "使你的撕裂技能的持续时间延长6秒。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_gouge.jpg",
    },
    [WarriorMajorGlyph.GlyphOfResonatingPower]: {
      name: "共鸣雕文",
      description: "使你的雷霆一击消耗的怒气值减少5点。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_nature_thunderclap.jpg",
    },
    [WarriorMajorGlyph.GlyphOfRevenge]: {
      name: "复仇雕文",
      description: "在使用复仇技能之后，你的下一次英勇打击技能不消耗怒气值。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_warrior_revenge.jpg",
    },
    [WarriorMajorGlyph.GlyphOfShieldWall]: {
      name: "盾墙雕文",
      description: "盾墙的冷却时间缩短2分钟，但是只能减免40%的伤害。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_warrior_shieldwall.jpg",
    },
    [WarriorMajorGlyph.GlyphOfShockwave]: {
      name: "震荡波雕文",
      description: "震荡波的冷却时间缩短3秒。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_warrior_shockwave.jpg",
    },
    [WarriorMajorGlyph.GlyphOfSpellReflection]: {
      name: "法术反射雕文",
      description: "法术反射的冷却时间缩短1秒。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_warrior_shieldreflection.jpg",
    },
    [WarriorMajorGlyph.GlyphOfSunderArmor]: {
      name: "破甲雕文",
      description: "你的破甲技能可以影响到附近的一个额外目标。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_warrior_sunder.jpg",
    },
    [WarriorMajorGlyph.GlyphOfSweepingStrikes]: {
      name: "横扫攻击雕文",
      description: "你的横扫攻击技能的怒气值消耗降低100%。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_rogue_slicedice.jpg",
    },
    [WarriorMajorGlyph.GlyphOfTaunt]: {
      name: "嘲讽雕文",
      description: "使你的嘲讽技能的命中几率提高8%。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_nature_reincarnation.jpg",
    },
    [WarriorMajorGlyph.GlyphOfVictoryRush]: {
      name: "乘胜追击雕文",
      description: "你的乘胜追击技能造成爆击的几率提高30% 。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_warrior_devastate.jpg",
    },
    [WarriorMajorGlyph.GlyphOfVigilance]: {
      name: "警戒雕文",
      description: "你的警戒技能转移目标威胁值的效果提高5%。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_warrior_vigilance.jpg",
    },
    [WarriorMajorGlyph.GlyphOfWhirlwind]: {
      name: "旋风斩雕文",
      description: "使你的旋风斩的冷却时间缩短2秒。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_whirlwind.jpg",
    },
  },
  minorGlyphs: {
    [WarriorMinorGlyph.GlyphOfBattle]: {
      name: "战斗雕文",
      description: "使你的战斗怒吼的持续时间延长2分钟。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_warrior_battleshout.jpg",
    },
    [WarriorMinorGlyph.GlyphOfBloodrage]: {
      name: "血性狂暴雕文",
      description: "你的血性狂暴技能消耗的生命值降低100%。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_racial_bloodrage.jpg",
    },
    [WarriorMinorGlyph.GlyphOfCharge]: {
      name: "冲锋雕文",
      description: "使你的冲锋技能的射程延长5码。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_warrior_charge.jpg",
    },
    [WarriorMinorGlyph.GlyphOfCommand]: {
      name: "命令雕文",
      description: "使你的命令怒吼能力的持续时间延长2分钟。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_warrior_rallyingcry.jpg",
    },
    [WarriorMinorGlyph.GlyphOfEnduringVictory]: {
      name: "持久追击雕文",
      description: "使你可以使用乘胜追击技能的时间范围延长5秒。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_warrior_devastate.jpg",
    },
    [WarriorMinorGlyph.GlyphOfMockingBlow]: {
      name: "惩戒痛击雕文",
      description: "使你的惩戒痛击技能造成的伤害提高25%。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_warrior_punishingblow.jpg",
    },
    [WarriorMinorGlyph.GlyphOfThunderClap]: {
      name: "雷霆一击雕文",
      description: "使你的雷霆一击技能的影响半径延长2码。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_nature_thunderclap.jpg",
    },
    [WarriorMinorGlyph.GlyphOfShatteringThrow]: {
      name: "碎裂投掷雕文",
      description: "你的碎裂投掷变为瞬发，而且可以在任何姿态下使用，但不再移除无敌效果，而且不能对玩家和玩家控制的目标使用。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_warrior_shatteringthrow.jpg",
    },
  },
};
