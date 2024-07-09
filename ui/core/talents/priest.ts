import {
  PriestTalents,
  PriestMajorGlyph,
  PriestMinorGlyph,
} from "../proto/priest.js";

import { GlyphsConfig } from "./glyphs_picker.js";
import { TalentsConfig, newTalentsConfig } from "./talents_picker.js";

import PriestTalentJson from "./trees/priest.json";

export const priestTalentsConfig: TalentsConfig<PriestTalents> =
  newTalentsConfig(PriestTalentJson);

export const priestGlyphsConfig: GlyphsConfig = {
  majorGlyphs: {
    [PriestMajorGlyph.GlyphOfCircleOfHealing]: {
      name: "治疗之环雕文",
      description: "你的治疗之环可以多治疗1个目标。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_holy_circleofrenewal.jpg",
    },
    [PriestMajorGlyph.GlyphOfDispelMagic]: {
      name: "驱散魔法雕文",
      description: "你的驱散魔法可以为目标恢复3%的生命值。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_holy_dispelmagic.jpg",
    },
    [PriestMajorGlyph.GlyphOfDispersion]: {
      name: "消散雕文",
      description: "消散的冷却时间缩短45秒。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_shadow_dispersion.jpg",
    },
    [PriestMajorGlyph.GlyphOfFade]: {
      name: "渐隐雕文",
      description: "使你的渐隐术的冷却时间缩短9秒。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_magic_lesserinvisibilty.jpg",
    },
    [PriestMajorGlyph.GlyphOfFearWard]: {
      name: "防护恐惧结界雕文",
      description: "防护恐惧结界的冷却时间和持续时间缩短60秒。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_holy_excorcism.jpg",
    },
    [PriestMajorGlyph.GlyphOfFlashHeal]: {
      name: "快速治疗雕文",
      description: "使你的快速治疗消耗的法力值降低10%。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_holy_flashheal.jpg",
    },
    [PriestMajorGlyph.GlyphOfGuardianSpirit]: {
      name: "守护之魂雕文",
      description: "如果你的守护之魂在整个效果持续期间都没有被触发，则这个法术的冷却时间缩短到1分钟。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_holy_guardianspirit.jpg",
    },
    [PriestMajorGlyph.GlyphOfHolyNova]: {
      name: "神圣新星雕文",
      description: "你的神圣新星的治疗量和伤害量提高20%。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_holy_holynova.jpg",
    },
    [PriestMajorGlyph.GlyphOfHymnOfHope]: {
      name: "希望圣歌雕文",
      description: "你的希望圣歌的持续时间延长2秒。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_holy_symbolofhope.jpg",
    },
    [PriestMajorGlyph.GlyphOfInnerFire]: {
      name: "心灵之火雕文",
      description: "使你的心灵之火提供的护甲值加成提高50%。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_holy_innerfire.jpg",
    },
    [PriestMajorGlyph.GlyphOfLightwell]: {
      name: "光明之泉雕文",
      description: "使你的光明之泉的治疗效果提高20%。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_holy_summonlightwell.jpg",
    },
    [PriestMajorGlyph.GlyphOfMassDispel]: {
      name: "群体驱散雕文",
      description: "使你的群体驱散消耗的法力值减少35%。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_arcane_massdispel.jpg",
    },
    [PriestMajorGlyph.GlyphOfMindControl]: {
      name: "精神控制雕文",
      description: "使目标抵抗或挣脱你的精神控制的几率降低17%。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_shadow_shadowworddominate.jpg",
    },
    [PriestMajorGlyph.GlyphOfMindFlay]: {
      name: "精神鞭笞雕文",
      description: "当目标受到暗言术：痛的效果影响时，使你的精神鞭笞造成的伤害提高10%。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_shadow_siphonmana.jpg",
    },
    [PriestMajorGlyph.GlyphOfMindSear]: {
      name: "精神灼烧雕文",
      description: "使精神灼烧的影响半径延长5码。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_shadow_mindshear.jpg",
    },
    [PriestMajorGlyph.GlyphOfPainSuppression]: {
      name: "痛苦压制雕文",
      description: "可以在昏迷时施放痛苦压制。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_holy_painsupression.jpg",
    },
    [PriestMajorGlyph.GlyphOfPenance]: {
      name: "苦修雕文",
      description: "苦修的冷却时间缩短2秒。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_holy_penance.jpg",
    },
    [PriestMajorGlyph.GlyphOfPowerWordShield]: {
      name: "真言术：盾雕文",
      description: "你的真言术：盾也会治疗目标，数值相当于吸收伤害量的20%。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_holy_powerwordshield.jpg",
    },
    [PriestMajorGlyph.GlyphOfPrayerOfHealing]: {
      name: "治疗祷言雕文",
      description: "你的治疗祷言可以在6秒内为目标恢复额外的生命值，数值相当于该次治疗祷言治疗量的20%。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_holy_prayerofhealing02.jpg",
    },
    [PriestMajorGlyph.GlyphOfPsychicScream]: {
      name: "心灵尖啸雕文",
      description: "使你的心灵尖啸的持续时间延长2秒，冷却时间延长8秒。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_shadow_psychicscream.jpg",
    },
    [PriestMajorGlyph.GlyphOfRenew]: {
      name: "恢复雕文",
      description: "使你的恢复的持续时间缩短3秒，但是每一跳的治疗量提高25%。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_holy_renew.jpg",
    },
    [PriestMajorGlyph.GlyphOfScourgeImprisonment]: {
      name: "天灾禁锢雕文",
      description: "使你的束缚亡灵的施法时间缩短1.0秒。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_nature_slow.jpg",
    },
    [PriestMajorGlyph.GlyphOfShadow]: {
      name: "暗影雕文",
      description: "在暗影形态下，你的非周期性法术爆击可以使你的法术强度提高，数值相当于你的精神值的30%，持续10 sec。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_holy_mindsooth.jpg",
    },
    [PriestMajorGlyph.GlyphOfShadowWordDeath]: {
      name: "暗言术：灭雕文",
      description: "当目标的生命值低于35%时，你的暗言术：灭对其造成的伤害提高10%。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_shadow_demonicfortitude.jpg",
    },
    [PriestMajorGlyph.GlyphOfShadowWordPain]: {
      name: "暗言术：痛雕文",
      description: "你的暗言术：痛法术的持续伤害会为你恢复基础法力值的1% 。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_shadow_shadowwordpain.jpg",
    },
    [PriestMajorGlyph.GlyphOfSmite]: {
      name: "惩击雕文",
      description: "你的惩击法术对于已经受到神圣之火影响的目标造成的伤害提高20%。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_holy_holysmite.jpg",
    },
    [PriestMajorGlyph.GlyphOfSpiritOfRedemption]: {
      name: "拯救之魂雕文",
      description: "救赎之魂的持续时间延长6秒。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/inv_enchant_essenceeternallarge.jpg",
    },
  },
  minorGlyphs: {
    [PriestMinorGlyph.GlyphOfFading]: {
      name: "渐隐雕文",
      description: "你的渐隐术的法力值消耗降低30%。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_magic_lesserinvisibilty.jpg",
    },
    [PriestMinorGlyph.GlyphOfFortitude]: {
      name: "坚韧雕文",
      description: "你的真言术：韧和坚韧祷言的法力值消耗降低50%。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_holy_wordfortitude.jpg",
    },
    [PriestMinorGlyph.GlyphOfLevitate]: {
      name: "漂浮雕文",
      description: "你的漂浮术不再需要施法材料。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_holy_layonhands.jpg",
    },
    [PriestMinorGlyph.GlyphOfShackleUndead]: {
      name: "束缚亡灵雕文",
      description: "使你的束缚亡灵法术的射程延长5码。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_nature_slow.jpg",
    },
    [PriestMinorGlyph.GlyphOfShadowProtection]: {
      name: "暗影防护雕文",
      description: "使你的暗影防护和暗影防护祷言的持续时间延长10分钟。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_shadow_antishadow.jpg",
    },
    [PriestMinorGlyph.GlyphOfShadowfiend]: {
      name: "暗影魔雕文",
      description: "如果你的暗影魔被杀死，则你可以恢复法力总值的5%。",
      iconUrl: "https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_shadow_shadowfiend.jpg",
    },
  },
};
