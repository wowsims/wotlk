import { DruidMajorGlyph, DruidMinorGlyph, DruidTalents } from '../proto/druid.js';
import { GlyphsConfig } from './glyphs_picker.js';
import { newTalentsConfig, TalentsConfig } from './talents_picker.js';
import DruidTalentsJson from './trees/druid.json';

export const druidTalentsConfig: TalentsConfig<DruidTalents> = newTalentsConfig(DruidTalentsJson);

export const druidGlyphsConfig: GlyphsConfig = {
	majorGlyphs: {
		[DruidMajorGlyph.GlyphOfBarkskin]: {
			name: '树皮雕文',
			description: '在激活树皮术的状态下，使你受到爆击的几率降低25%。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_nature_stoneclawtotem.jpg',
		},
		[DruidMajorGlyph.GlyphOfBerserk]: {
			name: '狂暴雕文',
			description: '狂暴的持续时间延长5秒。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_druid_berserk.jpg',
		},
		[DruidMajorGlyph.GlyphOfClaw]: {
			name: '爪击雕文',
			description: '你的爪击技能所消耗的能量值降低5点。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_druid_rake.jpg',
		},
		[DruidMajorGlyph.GlyphOfEntanglingRoots]: {
			name: '纠缠根须雕文',
			description: '使你的纠缠根须的目标在解除纠缠效果之前所能承受的伤害提高20%。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_nature_stranglevines.jpg',
		},
		[DruidMajorGlyph.GlyphOfFocus]: {
			name: '专注雕文',
			description: '星辰坠落造成的伤害提高10%，但是影响半径缩短50%。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_druid_starfall.jpg',
		},
		[DruidMajorGlyph.GlyphOfFrenziedRegeneration]: {
			name: '狂暴回复雕文',
			description: '在狂暴回复状态下时，你受到的治疗量提高20%。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_bullrush.jpg',
		},
		[DruidMajorGlyph.GlyphOfGrowling]: {
			name: '低吼雕文',
			description: '使你的低吼技能的生效几率提高8%。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_physical_taunt.jpg',
		},
		[DruidMajorGlyph.GlyphOfHealingTouch]: {
			name: '治疗之触雕文',
			description: '使你的治疗之触的施法时间缩短1.5秒，法力值消耗降低25%，治疗量降低50%。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_nature_healingtouch.jpg',
		},
		[DruidMajorGlyph.GlyphOfHurricane]: {
			name: '飓风雕文',
			description: '你的飓风法术可以令目标的移动速度降低20%。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_nature_cyclone.jpg',
		},
		[DruidMajorGlyph.GlyphOfInnervate]: {
			name: '激活雕文',
			description: '激活现在除了正常效果外，还会使施法者在10秒内获得其基础法力值的45%。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_nature_lightning.jpg',
		},
		[DruidMajorGlyph.GlyphOfInsectSwarm]: {
			name: '虫群雕文',
			description: '使你的虫群的伤害提高30%，但不再降低目标的命中几率。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_nature_insectswarm.jpg',
		},
		[DruidMajorGlyph.GlyphOfLifebloom]: {
			name: '生命绽放雕文',
			description: '使你的生命绽放的持续时间延长1秒。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/inv_misc_herb_felblossom.jpg',
		},
		[DruidMajorGlyph.GlyphOfMangle]: {
			name: '裂伤雕文',
			description: '使你的裂伤的伤害提高10%。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_druid_mangle2.jpg',
		},
		[DruidMajorGlyph.GlyphOfMaul]: {
			name: '重殴雕文',
			description: '你的重殴技能可以多攻击1个目标。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_druid_maul.jpg',
		},
		[DruidMajorGlyph.GlyphOfMonsoon]: {
			name: '季风雕文',
			description: '使你的台风法术的冷却时间缩短3秒。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_druid_typhoon.jpg',
		},
		[DruidMajorGlyph.GlyphOfMoonfire]: {
			name: '月火雕文',
			description: '使你的月火术的持续伤害提高75%，但是初始伤害降低90%。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_nature_starfall.jpg',
		},
		[DruidMajorGlyph.GlyphOfNourish]: {
			name: '滋养雕文',
			description: '你施加于目标身上的每个持续治疗效果都可以使你的滋养法术的治疗量提高6%。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_druid_nourish.jpg',
		},
		[DruidMajorGlyph.GlyphOfOmenOfClarity]: {
			name: '清晰预兆雕文',
			description: '你成功施放精灵之火（野性）时，有100%的几率触发清晰预兆天赋。对玩家或玩家控制的宠物使用时不会触发。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_nature_faeriefire.jpg',
		},
		[DruidMajorGlyph.GlyphOfRake]: {
			name: '斜掠雕文',
			description: '你的斜掠技能可以阻止目标逃跑。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_druid_disembowel.jpg',
		},
		[DruidMajorGlyph.GlyphOfRapidRejuvenation]: {
			name: '急速回春雕文',
			description: '你的急速效果现在可以缩短回春术造成的周期性治疗间隔。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_nature_rejuvenation.jpg',
		},
		[DruidMajorGlyph.GlyphOfRebirth]: {
			name: '复生雕文',
			description: '被你的复生法术所复活的目标拥有100%的起始生命值。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_nature_reincarnation.jpg',
		},
		[DruidMajorGlyph.GlyphOfRegrowth]: {
			name: '愈合雕文',
			description: '当目标身上仍然有你施放的愈合效果时，使你的愈合法术的治疗量提高20%。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_nature_resistnature.jpg',
		},
		[DruidMajorGlyph.GlyphOfRejuvenation]: {
			name: '回春雕文',
			description: '当你的回春术目标生命值低于50%时，你将为其恢复50%的额外生命值。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_nature_rejuvenation.jpg',
		},
		[DruidMajorGlyph.GlyphOfRip]: {
			name: '割裂雕文',
			description: '使你的割裂技能的持续时间延长4秒。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_ghoulfrenzy.jpg',
		},
		[DruidMajorGlyph.GlyphOfSavageRoar]: {
			name: '野蛮咆哮雕文',
			description: '你的野蛮咆哮技能提供的伤害加成效果提高3%。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_druid_skinteeth.jpg',
		},
		[DruidMajorGlyph.GlyphOfShred]: {
			name: '撕碎雕文',
			description: '你每次使用撕碎技能，你施加于目标身上的割裂效果的持续时间就延长2秒，最多可以延长6秒。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_shadow_vampiricaura.jpg',
		},
		[DruidMajorGlyph.GlyphOfStarfall]: {
			name: '星辰坠落雕文',
			description: '星辰坠落的冷却时间缩短30秒。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_druid_starfall.jpg',
		},
		[DruidMajorGlyph.GlyphOfStarfire]: {
			name: '星火雕文',
			description: '你的星火术可以令你对目标施放的月火术的效果持续时间延长3秒，最多延长9秒。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_arcane_starfire.jpg',
		},
		[DruidMajorGlyph.GlyphOfSurvivalInstincts]: {
			name: '生存本能雕文',
			description: '你的生存本能提供的生命值加成提高15%。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_druid_tigersroar.jpg',
		},
		[DruidMajorGlyph.GlyphOfSwiftmend]: {
			name: '迅捷治疗雕文',
			description: '你的迅捷治愈技能不再吞噬目标身上的愈合或回春效果。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/inv_relics_idolofrejuvenation.jpg',
		},
		[DruidMajorGlyph.GlyphOfWildGrowth]: {
			name: '野性成长雕文',
			description: '野性成长可以影响1个额外的目标。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_druid_flourish.jpg',
		},
		[DruidMajorGlyph.GlyphOfWrath]: {
			name: '愤怒雕文',
			description: '使你在施放愤怒法术时因受到伤害而承受的施法推迟时间缩短50%。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_nature_abolishmagic.jpg',
		},
	},
	minorGlyphs: {
		[DruidMinorGlyph.GlyphOfAquaticForm]: {
			name: '水栖形态雕文',
			description: '在水栖形态下使你的游泳速度提高50%。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_druid_aquaticform.jpg',
		},
		[DruidMinorGlyph.GlyphOfChallengingRoar]: {
			name: '挑战咆哮雕文',
			description: '使你的挑战咆哮技能的冷却时间缩短30秒。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_druid_challangingroar.jpg',
		},
		[DruidMinorGlyph.GlyphOfDash]: {
			name: '突进雕文',
			description: '使你的突进技能的冷却时间缩短20%。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_druid_dash.jpg',
		},
		[DruidMinorGlyph.GlyphOfTheWild]: {
			name: '野性赐福雕文',
			description: '你的野性印记和野性赐福法术的法力值消耗降低50%。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_nature_regeneration.jpg',
		},
		[DruidMinorGlyph.GlyphOfThorns]: {
			name: '荆棘雕文',
			description: '当施放在自己身上时，你的荆棘法术的持续时间延长50分钟。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_nature_thorns.jpg',
		},
		[DruidMinorGlyph.GlyphOfTyphoon]: {
			name: '台风雕文',
			description: '使你的台风法术的法力消耗降低8%，并使其半径增加10码，但不再击退敌人。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_druid_typhoon.jpg',
		},
		[DruidMinorGlyph.GlyphOfUnburdenedRebirth]: {
			name: '无负担复生雕文',
			description: '你的复生法术不再需要施法材料。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_nature_reincarnation.jpg',
		},
	},
};
