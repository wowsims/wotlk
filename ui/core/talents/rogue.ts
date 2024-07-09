import { RogueMajorGlyph, RogueMinorGlyph,RogueTalents } from '../proto/rogue.js';
import { GlyphsConfig } from './glyphs_picker.js';
import { newTalentsConfig,TalentsConfig } from './talents_picker.js';
import RogueTalentJson from './trees/rogue.json';

export const rogueTalentsConfig: TalentsConfig<RogueTalents> = newTalentsConfig(RogueTalentJson);

export const rogueGlyphsConfig: GlyphsConfig = {
	majorGlyphs: {
		[RogueMajorGlyph.GlyphOfAdrenalineRush]: {
			name: '冲动雕文',
			description: '使你的冲动技能的持续时间延长5秒。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_shadow_shadowworddominate.jpg',
		},
		[RogueMajorGlyph.GlyphOfAmbush]: {
			name: '伏击雕文',
			description: '伏击的射程延长5码。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_rogue_ambush.jpg',
		},
		[RogueMajorGlyph.GlyphOfBackstab]: {
			name: '背刺雕文',
			description: '你每次使用背刺技能，你施加于目标身上的割裂效果的持续时间就延长2秒，最多可以延长6秒。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_backstab.jpg',
		},
		[RogueMajorGlyph.GlyphOfBladeFlurry]: {
			name: '剑刃乱舞雕文',
			description: '剑刃乱舞的能量值消耗降低100%。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_warrior_punishingblow.jpg',
		},
		[RogueMajorGlyph.GlyphOfCloakOfShadows]: {
			name: '暗影斗篷雕文',
			description: '当暗影斗篷处于激活状态时，你受到的物理伤害降低40%。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_shadow_nethercloak.jpg',
		},
		[RogueMajorGlyph.GlyphOfCripplingPoison]: {
			name: '减速药膏雕文',
			description: '减速药膏的效果触发几率提高20%。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_poisonsting.jpg',
		},
		[RogueMajorGlyph.GlyphOfDeadlyThrow]: {
			name: '致命投掷雕文',
			description: '使你的致命投掷的减速效果提高20%。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/inv_throwingknife_06.jpg',
		},
		[RogueMajorGlyph.GlyphOfEvasion]: {
			name: '闪避雕文',
			description: '使你的闪避技能的持续时间延长5秒。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_shadow_shadowward.jpg',
		},
		[RogueMajorGlyph.GlyphOfEviscerate]: {
			name: '刺骨雕文',
			description: '刺骨的爆击几率提高10%。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_rogue_eviscerate.jpg',
		},
		[RogueMajorGlyph.GlyphOfExposeArmor]: {
			name: '破甲雕文',
			description: '使你的破甲技能的持续时间延长12秒。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_warrior_riposte.jpg',
		},
		[RogueMajorGlyph.GlyphOfFanOfKnives]: {
			name: '刀扇雕文',
			description: '刀扇造成的伤害提高20%。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_rogue_fanofknives.jpg',
		},
		[RogueMajorGlyph.GlyphOfFeint]: {
			name: '佯攻雕文',
			description: '佯攻的能量值消耗降低20点。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_rogue_feint.jpg',
		},
		[RogueMajorGlyph.GlyphOfGarrote]: {
			name: '锁喉雕文',
			description: '锁喉的持续时间缩短3秒，总伤害量提高20%。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_rogue_garrote.jpg',
		},
		[RogueMajorGlyph.GlyphOfGhostlyStrike]: {
			name: '鬼魅攻击雕文',
			description: '使你的鬼魅攻击造成的伤害提高40%，效果持续时间延长4秒，但是冷却时间延长10秒。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_shadow_curse.jpg',
		},
		[RogueMajorGlyph.GlyphOfGouge]: {
			name: '凿击雕文',
			description: '凿击的能量值消耗降低15点。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_gouge.jpg',
		},
		[RogueMajorGlyph.GlyphOfHemorrhage]: {
			name: '出血雕文',
			description: '受到出血效果影响的目标身上的伤害加成提高40%。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_shadow_lifedrain.jpg',
		},
		[RogueMajorGlyph.GlyphOfHungerForBlood]: {
			name: '血之饥渴雕文',
			description: '血之饥渴提供的伤害加成效果提高3%。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_rogue_hungerforblood.jpg',
		},
		[RogueMajorGlyph.GlyphOfKillingSpree]: {
			name: '杀戮盛筵雕文',
			description: '杀戮盛筵的冷却时间缩短45秒。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_rogue_murderspree.jpg',
		},
		[RogueMajorGlyph.GlyphOfMutilate]: {
			name: '毁伤雕文',
			description: '毁伤的能量值消耗降低5点。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_rogue_shadowstrikes.jpg',
		},
		[RogueMajorGlyph.GlyphOfPreparation]: {
			name: '伺机待发雕文',
			description: '你的伺机待发技能也会立即重置剑刃乱舞、拆卸和脚踢的冷却时间。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_shadow_antishadow.jpg',
		},
		[RogueMajorGlyph.GlyphOfRupture]: {
			name: '割裂雕文',
			description: '使你的割裂技能的持续时间延长4秒。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_rogue_rupture.jpg',
		},
		[RogueMajorGlyph.GlyphOfSap]: {
			name: '闷棍雕文',
			description: '使你的闷棍的持续时间延长20秒。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_sap.jpg',
		},
		[RogueMajorGlyph.GlyphOfShadowDance]: {
			name: '暗影之舞雕文',
			description: '暗影之舞的持续时间延长2秒。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_rogue_shadowdance.jpg',
		},
		[RogueMajorGlyph.GlyphOfSinisterStrike]: {
			name: '影袭雕文',
			description: '你的影袭爆击有50%的几率增加一个额外的连击点数。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_shadow_ritualofsacrifice.jpg',
		},
		[RogueMajorGlyph.GlyphOfSliceAndDice]: {
			name: '切割雕文',
			description: '使你的切割的持续时间延长3秒。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_rogue_slicedice.jpg',
		},
		[RogueMajorGlyph.GlyphOfSprint]: {
			name: '疾跑雕文',
			description: '使你的疾跑技能的移动加速效果提高30%。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_rogue_sprint.jpg',
		},
		[RogueMajorGlyph.GlyphOfTricksOfTheTrade]: {
			name: '嫁祸诀窍雕文',
			description: '你的嫁祸诀窍提供的伤害加成效果延长4秒。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_rogue_tricksofthetrade.jpg',
		},
		[RogueMajorGlyph.GlyphOfVigor]: {
			name: '精力雕文',
			description: '精力提供的额外能量值上限提高10点。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_nature_earthbindtotem.jpg',
		},
	},
	minorGlyphs: {
		[RogueMinorGlyph.GlyphOfBlurredSpeed]: {
			name: '水上漂雕文',
			description: '使你在激活疾跑技能的状态下可以在水面上行走。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_rogue_sprint.jpg',
		},
		[RogueMinorGlyph.GlyphOfDistract]: {
			name: '扰乱雕文',
			description: '使你的扰乱技能的作用范围延长5码。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_rogue_distract.jpg',
		},
		[RogueMinorGlyph.GlyphOfPickLock]: {
			name: '开锁雕文',
			description: '你的开锁技能的施放时间缩短100%。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_nature_moonkey.jpg',
		},
		[RogueMinorGlyph.GlyphOfPickPocket]: {
			name: '妙手空空雕文',
			description: '使你的搜索技能的射程延长5码。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/inv_misc_bag_11.jpg',
		},
		[RogueMinorGlyph.GlyphOfSafeFall]: {
			name: '安全降落雕文',
			description: '使你的安全降落技能允许你从高处跳下来而不受伤的距离增加。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/inv_feather_01.jpg',
		},
		[RogueMinorGlyph.GlyphOfVanish]: {
			name: '消失雕文',
			description: '在消失效果处于激活状态下时，移动速度提高30%。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_vanish.jpg',
		},
	},
};
