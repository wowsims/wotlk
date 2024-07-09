import { Player } from '../player.js';
import { Spec } from '../proto/common.js';
import { DeathknightMajorGlyph, DeathknightMinorGlyph,DeathknightTalents } from '../proto/deathknight.js';
import { GlyphsConfig } from './glyphs_picker.js';
import { newTalentsConfig,TalentsConfig } from './talents_picker.js';
import DkTalentsJson from './trees/deathknight.json';

export const deathknightTalentsConfig: TalentsConfig<DeathknightTalents> = newTalentsConfig(DkTalentsJson);

export const deathknightGlyphsConfig: GlyphsConfig = {
	majorGlyphs: {
		[DeathknightMajorGlyph.GlyphOfAntiMagicShell]: {
			name: '反魔法护罩雕文',
			description: '使你的反魔法护罩的持续时间延长2秒。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_shadow_antimagicshell.jpg',
		},
		[DeathknightMajorGlyph.GlyphOfBloodStrike]: {
			name: '鲜血打击雕文',
			description: '你的鲜血打击对被诱捕的目标造成的伤害提高20%。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_deathknight_deathstrike.jpg',
		},
		[DeathknightMajorGlyph.GlyphOfBoneShield]: {
			name: '白骨之盾雕文',
			description: '为你的白骨之盾增加1次使用次数。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/inv_chest_leather_13.jpg',
		},
		[DeathknightMajorGlyph.GlyphOfChainsOfIce]: {
			name: '寒冰锁链雕文',
			description: '你的寒冰锁链可以造成额外的144 to 156点冰霜伤害，这个数值受到你的攻击强度加成影响。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_frost_chainsofice.jpg',
		},
		[DeathknightMajorGlyph.GlyphOfDancingRuneWeapon]: {
			name: '符文刃舞雕文',
			description: '符文刃舞的持续时间延长5秒。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/inv_sword_07.jpg',
		},
		[DeathknightMajorGlyph.GlyphOfDarkCommand]: {
			name: '黑暗命令雕文',
			description: '使你的黑暗命令技能的成功率提高8%。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_nature_shamanrage.jpg',
		},
		[DeathknightMajorGlyph.GlyphOfDarkDeath]: {
			name: '黑暗死亡雕文',
			description: '死亡缠绕造成的伤害或治疗效果提高15%。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_shadow_deathcoil.jpg',
		},
		[DeathknightMajorGlyph.GlyphOfDeathAndDecay]: {
			name: '枯萎凋零雕文',
			description: '你的枯萎凋零法术造成的伤害提高20%。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_shadow_deathanddecay.jpg',
		},
		[DeathknightMajorGlyph.GlyphOfDeathGrip]: {
			name: '死亡之握雕文',
			description: '当你杀死一个可以提供经验值或荣誉值的目标时，你的死亡之握的冷却时间就会刷新。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_deathknight_strangulate.jpg',
		},
		[DeathknightMajorGlyph.GlyphOfDeathStrike]: {
			name: '灵界打击雕文',
			description: '你的每1点符文能量值都可以使你的灵界打击造成的伤害提高1%（最多25%）。这个效果不会消耗符文能量值。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_deathknight_butcher2.jpg',
		},
		[DeathknightMajorGlyph.GlyphOfDisease]: {
			name: '疾病雕文',
			description: '你的传染技能可以使你的主要目标身上的疾病效果持续时间、疾病附加效果持续时间刷新到起始状态。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_shadow_plaguecloud.jpg',
		},
		[DeathknightMajorGlyph.GlyphOfFrostStrike]: {
			name: '冰霜打击雕文',
			description: '你的冰霜打击消耗的符文能量值减少8点。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_deathknight_empowerruneblade2.jpg',
		},
		[DeathknightMajorGlyph.GlyphOfHeartStrike]: {
			name: '心脏打击雕文',
			description: '你的心脏打击将使目标的移动速度降低50%，持续10 sec。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/inv_weapon_shortblade_40.jpg',
		},
		[DeathknightMajorGlyph.GlyphOfHowlingBlast]: {
			name: '凛风冲击雕文',
			description: '你的凛风冲击技能可以使目标感染冰霜疫病。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_frost_arcticwinds.jpg',
		},
		[DeathknightMajorGlyph.GlyphOfHungeringCold]: {
			name: '饥饿之寒雕文',
			description: '饥饿之寒的符文能量值消耗降低40点。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/inv_staff_15.jpg',
		},
		[DeathknightMajorGlyph.GlyphOfIceboundFortitude]: {
			name: '冰封之韧雕文',
			description: '无论你的防御技能值是多少，你的冰封之韧总是可以提供至少40%的伤害减免。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_deathknight_iceboundfortitude.jpg',
		},
		[DeathknightMajorGlyph.GlyphOfIcyTouch]: {
			name: '冰冷触摸雕文',
			description: '你的冰霜疫病造成20%的额外伤害。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_deathknight_icetouch.jpg',
		},
		[DeathknightMajorGlyph.GlyphOfObliterate]: {
			name: '湮没雕文',
			description: '你的湮没技能造成的伤害提高25%。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_deathknight_classicon.jpg',
		},
		[DeathknightMajorGlyph.GlyphOfPlagueStrike]: {
			name: '暗影打击雕文',
			description: '你的暗影打击造成伤害提高20%。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_deathknight_empowerruneblade.jpg',
		},
		[DeathknightMajorGlyph.GlyphOfRuneStrike]: {
			name: '符文打击雕文',
			description: '使你的符文打击技能的爆击几率提高10%。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_deathknight_darkconviction.jpg',
		},
		[DeathknightMajorGlyph.GlyphOfRuneTap]: {
			name: '符文分流雕文',
			description: '你的符文分流现在额外治疗你1%的最大生命值，而且还治疗你的小队10%的最大生命值。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_deathknight_runetap.jpg',
		},
		[DeathknightMajorGlyph.GlyphOfScourgeStrike]: {
			name: '天灾打击雕文',
			description: '你的天灾打击会使目标身上的疫病效果延长3秒，最多额外延长9秒。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_deathknight_scourgestrike.jpg',
		},
		[DeathknightMajorGlyph.GlyphOfStrangulate]: {
			name: '绞袭雕文',
			description: '使你的绞袭技能的冷却时间缩短20秒。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_shadow_soulleech_3.jpg',
		},
		[DeathknightMajorGlyph.GlyphOfTheGhoul]: {
			name: '食尸鬼雕文',
			description: '使你的食尸鬼的力量值提高，数值相当于你的力量值的40%；你的食尸鬼的耐力值提高，数值相当于你的耐力值的40%。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_shadow_animatedead.jpg',
		},
		[DeathknightMajorGlyph.GlyphOfUnbreakableArmor]: {
			name: '铜墙铁壁雕文',
			description: '从铜墙铁壁中所获得的护甲值提高30%。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/inv_armor_helm_plate_naxxramas_raidwarrior_c_01.jpg',
		},
		[DeathknightMajorGlyph.GlyphOfUnholyBlight]: {
			name: '邪恶虫群雕文',
			description: '邪恶虫群伤害提高40%。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_shadow_contagion.jpg',
		},
		[DeathknightMajorGlyph.GlyphOfVampiricBlood]: {
			name: '吸血鬼之血雕文',
			description: '你的吸血鬼之血的持续时间延长5秒。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_shadow_lifedrain.jpg',
		},
	},
	minorGlyphs: {
		[DeathknightMinorGlyph.GlyphOfBloodTap]: {
			name: '活力分流雕文',
			description: '你的活力分流不再对你自身造成伤害。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_deathknight_bloodtap.jpg',
		},
		[DeathknightMinorGlyph.GlyphOfCorpseExplosion]: {
			name: '邪爆雕文',
			description: '邪爆范围提高5码。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_creature_disease_02.jpg',
		},
		[DeathknightMinorGlyph.GlyphOfDeathSEmbrace]: {
			name: '凋零之拥雕文',
			description: '当你的凋零缠绕被用来进行治疗时，可以返还20点符文能量。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_shadow_deathcoil.jpg',
		},
		[DeathknightMinorGlyph.GlyphOfHornOfWinter]: {
			name: '寒冬号角雕文',
			description: '寒冬号角的效果持续时间延长1分钟。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/inv_misc_horn_02.jpg',
		},
		[DeathknightMinorGlyph.GlyphOfPestilence]: {
			name: '传染雕文',
			description: '使你的传染技能的影响半径延长5码。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_shadow_plaguecloud.jpg',
		},
		[DeathknightMinorGlyph.GlyphOfRaiseDead]: {
			name: '亡者复生雕文',
			description: '你的亡者复生不再消耗施法材料。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_shadow_animatedead.jpg',
		},
	},
};
