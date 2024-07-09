import { HunterMajorGlyph, HunterMinorGlyph, HunterPetTalents,HunterTalents } from '../proto/hunter.js';
import { GlyphsConfig } from './glyphs_picker.js';
import { newTalentsConfig,TalentsConfig } from './talents_picker.js';
import HunterTalentJson from './trees/hunter.json';

export const hunterTalentsConfig: TalentsConfig<HunterTalents> = newTalentsConfig(HunterTalentJson);

export const hunterGlyphsConfig: GlyphsConfig = {
	majorGlyphs: {
		[HunterMajorGlyph.GlyphOfAimedShot]: {
			name: '瞄准射击雕文',
			description: '使你的瞄准射击的冷却时间缩短2秒。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/inv_spear_07.jpg',
		},
		[HunterMajorGlyph.GlyphOfArcaneShot]: {
			name: '奥术射击雕文',
			description: '当你的目标身上有你施加的钉刺效果时，你的奥术射击可以返还20%的法力值消耗。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_impalingbolt.jpg',
		},
		[HunterMajorGlyph.GlyphOfAspectOfTheViper]: {
			name: '蝰蛇守护雕文',
			description: '你通过蝰蛇守护恢复的法力值提高10%。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_hunter_aspectoftheviper.jpg',
		},
		[HunterMajorGlyph.GlyphOfBestialWrath]: {
			name: '狂野怒火雕文',
			description: '使你的狂野怒火技能的冷却时间缩短20秒。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_druid_ferociousbite.jpg',
		},
		[HunterMajorGlyph.GlyphOfChimeraShot]: {
			name: '奇美拉射击雕文',
			description: '奇美拉射击的冷却时间缩短1秒。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_hunter_chimerashot2.jpg',
		},
		[HunterMajorGlyph.GlyphOfDeterrence]: {
			name: '威慑雕文',
			description: '使你的威慑技能的冷却时间缩短10秒。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_whirlwind.jpg',
		},
		[HunterMajorGlyph.GlyphOfDisengage]: {
			name: '逃脱雕文',
			description: '使你的逃脱技能的冷却时间缩短5秒。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_rogue_feint.jpg',
		},
		[HunterMajorGlyph.GlyphOfExplosiveShot]: {
			name: '爆炸射击雕文',
			description: '使你的爆炸射击的爆击几率提高4%。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_hunter_explosiveshot.jpg',
		},
		[HunterMajorGlyph.GlyphOfExplosiveTrap]: {
			name: '爆炸陷阱雕文',
			description: '你的爆炸陷阱的持续伤害可以爆击。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_fire_selfdestruct.jpg',
		},
		[HunterMajorGlyph.GlyphOfFreezingTrap]: {
			name: '冰冻陷阱雕文',
			description: '当你的冰冻陷阱效果被打破时，被冻结的目标的移动速度降低30%，持续4秒。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_frost_chainsofice.jpg',
		},
		[HunterMajorGlyph.GlyphOfFrostTrap]: {
			name: '冰霜陷阱雕文',
			description: '使你的冰霜陷阱的作用半径延长2码。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_frost_freezingbreath.jpg',
		},
		[HunterMajorGlyph.GlyphOfHuntersMark]: {
			name: '猎人印记雕文',
			description: '使你的猎人印记提供的攻击强度加成提高20%。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_hunter_snipershot.jpg',
		},
		[HunterMajorGlyph.GlyphOfImmolationTrap]: {
			name: '献祭陷阱雕文',
			description: '使你的献祭陷阱的持续时间缩短6秒，但是触发时造成的伤害提高100%。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_fire_flameshock.jpg',
		},
		[HunterMajorGlyph.GlyphOfKillShot]: {
			name: '杀戮射击雕文',
			description: '杀戮射击的冷却时间缩短6秒。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_hunter_assassinate2.jpg',
		},
		[HunterMajorGlyph.GlyphOfMending]: {
			name: '治疗雕文',
			description: '使你的治疗宠物技能的治疗量提高40%。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_hunter_mendpet.jpg',
		},
		[HunterMajorGlyph.GlyphOfMultiShot]: {
			name: '多重射击雕文',
			description: '使你的多重射击的冷却时间缩短1秒。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_upgrademoonglaive.jpg',
		},
		[HunterMajorGlyph.GlyphOfRapidFire]: {
			name: '急速射击雕文',
			description: '使你的急速射击所提供的急速提高8%。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_hunter_runningshot.jpg',
		},
		[HunterMajorGlyph.GlyphOfRaptorStrike]: {
			name: '猛禽一击雕文',
			description: '在使用猛禽一击之后的3秒内，你受到的伤害降低20%。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_meleedamage.jpg',
		},
		[HunterMajorGlyph.GlyphOfScatterShot]: {
			name: '驱散射击雕文',
			description: '使你的驱散射击的射程延长3码。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_golemstormbolt.jpg',
		},
		[HunterMajorGlyph.GlyphOfSerpentSting]: {
			name: '毒蛇钉刺雕文',
			description: '使你的毒蛇钉刺的持续时间延长6秒。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_hunter_quickshot.jpg',
		},
		[HunterMajorGlyph.GlyphOfSnakeTrap]: {
			name: '毒蛇陷阱雕文',
			description: '你的毒蛇陷阱所产生的毒蛇受到的范围攻击伤害降低90%。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_hunter_snaketrap.jpg',
		},
		[HunterMajorGlyph.GlyphOfSteadyShot]: {
			name: '稳固射击雕文',
			description: '当你的目标受到毒蛇钉刺效果影响时，使你的稳固射击对其造成的伤害提高10%。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_hunter_steadyshot.jpg',
		},
		[HunterMajorGlyph.GlyphOfTheBeast]: {
			name: '野兽雕文',
			description: '使你的野兽守护为你和你的宠物提供的攻击强度加成提高2%。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_mount_pinktiger.jpg',
		},
		[HunterMajorGlyph.GlyphOfTheHawk]: {
			name: '雄鹰雕文',
			description: '使你的强化雄鹰守护所提供的效果提高6%。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_nature_ravenform.jpg',
		},
		[HunterMajorGlyph.GlyphOfTrueshotAura]: {
			name: '强击光环雕文',
			description: '在激活强击光环的状态下，你的瞄准射击的爆击几率提高10%。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_trueshot.jpg',
		},
		[HunterMajorGlyph.GlyphOfVolley]: {
			name: '乱射雕文',
			description: '使你的乱射技能的法力值消耗降低20%。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_marksmanship.jpg',
		},
		[HunterMajorGlyph.GlyphOfWyvernSting]: {
			name: '翼龙钉刺雕文',
			description: '使你的翼龙钉刺的冷却时间缩短6秒。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/inv_spear_02.jpg',
		},
	},
	minorGlyphs: {
		[HunterMinorGlyph.GlyphOfFeignDeath]: {
			name: '假死雕文',
			description: '你的假死技能的冷却时间缩短5秒。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_rogue_feigndeath.jpg',
		},
		[HunterMinorGlyph.GlyphOfMendPet]: {
			name: '治疗宠物雕文',
			description: '你的治疗宠物法术可以略微提高宠物的快乐值。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_hunter_mendpet.jpg',
		},
		[HunterMinorGlyph.GlyphOfPossessedStrength]: {
			name: '支配之力雕文',
			description: '使你的宠物在使用野兽之眼技能时造成的伤害提高50%。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_eyeoftheowl.jpg',
		},
		[HunterMinorGlyph.GlyphOfRevivePet]: {
			name: '复活宠物雕文',
			description: '使你在施放复活宠物时因受到伤害而承受的施法推迟时间缩短100%。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_hunter_beastsoothe.jpg',
		},
		[HunterMinorGlyph.GlyphOfScareBeast]: {
			name: '恐吓野兽雕文',
			description: '使你在施放恐吓野兽时因受到伤害而承受的施法推迟时间缩短75%。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_druid_cower.jpg',
		},
		[HunterMinorGlyph.GlyphOfThePack]: {
			name: '豹群雕文',
			description: '使你的豹群守护技能的影响半径延长15码。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_mount_jungletiger.jpg',
		},
	},
};

