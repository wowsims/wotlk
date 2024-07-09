import { MageMajorGlyph, MageMinorGlyph,MageTalents } from '../proto/mage.js';
import { GlyphsConfig } from './glyphs_picker.js';
import { newTalentsConfig,TalentsConfig } from './talents_picker.js';
import MageTalentJson from './trees/mage.json';

export const mageTalentsConfig: TalentsConfig<MageTalents> = newTalentsConfig(MageTalentJson);

export const mageGlyphsConfig: GlyphsConfig = {
	majorGlyphs: {
		[MageMajorGlyph.GlyphOfArcaneBarrage]: {
			name: '奥术弹幕雕文',
			description: '奥术弹幕的法力值消耗降低20%。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_mage_arcanebarrage.jpg',
		},
		[MageMajorGlyph.GlyphOfArcaneBlast]: {
			name: '奥术冲击雕文',
			description: '使你由奥术冲击获得的伤害加成提高3%。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_arcane_blast.jpg',
		},
		[MageMajorGlyph.GlyphOfArcaneExplosion]: {
			name: '魔爆雕文',
			description: '魔爆术的法力值消耗降低10%。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_nature_wispsplode.jpg',
		},
		[MageMajorGlyph.GlyphOfArcaneMissiles]: {
			name: '奥术飞弹雕文',
			description: '奥术飞弹的爆击伤害加成提高25%。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_nature_starfall.jpg',
		},
		[MageMajorGlyph.GlyphOfArcanePower]: {
			name: '奥术强化雕文',
			description: '使你的奥术强化的持续时间延长3秒。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_nature_lightning.jpg',
		},
		[MageMajorGlyph.GlyphOfBlink]: {
			name: '闪现雕文',
			description: '使你的闪现术的移动距离延长5码。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_arcane_blink.jpg',
		},
		[MageMajorGlyph.GlyphOfDeepFreeze]: {
			name: '深度冻结雕文',
			description: '使深度冻结的射程延长10码。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_mage_deepfreeze.jpg',
		},
		[MageMajorGlyph.GlyphOfEternalWater]: {
			name: '永恒之水雕文',
			description: '你所召唤的水元素现在可以一直协同你作战，但水元素无法施放冰冻术。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_frost_summonwaterelemental_2.jpg',
		},
		[MageMajorGlyph.GlyphOfEvocation]: {
			name: '唤醒雕文',
			description: '你的唤醒可以使你恢复60%的生命值。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_nature_purge.jpg',
		},
		[MageMajorGlyph.GlyphOfFireBlast]: {
			name: '火焰冲击雕文',
			description: '当目标处于昏迷或瘫痪状态下时，使你的火焰冲击的爆击几率提高50%。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_fire_fireball.jpg',
		},
		[MageMajorGlyph.GlyphOfFireball]: {
			name: '火球雕文',
			description: '使你的火球术的施法时间减少0.15秒，但是不再造成持续伤害。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_fire_flamebolt.jpg',
		},
		[MageMajorGlyph.GlyphOfFrostNova]: {
			name: '冰霜新星雕文',
			description: '你的冰霜新星的目标在解除冰冻效果之前可以承受的伤害提高20%。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_frost_frostnova.jpg',
		},
		[MageMajorGlyph.GlyphOfFrostbolt]: {
			name: '寒冰箭雕文',
			description: '使你的寒冰箭造成的伤害提高5%，但是不再具有减速效果。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_frost_frostbolt02.jpg',
		},
		[MageMajorGlyph.GlyphOfFrostfire]: {
			name: '霜火雕文',
			description: '使你的霜火之箭的初始伤害值提高2%，爆击几率提高2%。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_mage_frostfirebolt.jpg',
		},
		[MageMajorGlyph.GlyphOfIceArmor]: {
			name: '冰甲雕文',
			description: '你的冰甲术和霜甲术提供的护甲值和抗性提高50%。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_frost_frostarmor02.jpg',
		},
		[MageMajorGlyph.GlyphOfIceBarrier]: {
			name: '寒冰护体雕文',
			description: '使你的寒冰护体吸收的伤害值提高30%。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_ice_lament.jpg',
		},
		[MageMajorGlyph.GlyphOfIceBlock]: {
			name: '寒冰屏障雕文',
			description: '每当你使用寒冰屏障，你的冰霜新星的冷却时间就会结束。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_frost_frost.jpg',
		},
		[MageMajorGlyph.GlyphOfIceLance]: {
			name: '冰枪雕文',
			description: '对于等级比你高且被冻结的目标，你的冰枪术造成4倍伤害，而不是3倍。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_frost_frostblast.jpg',
		},
		[MageMajorGlyph.GlyphOfIcyVeins]: {
			name: '冰冷血脉雕文',
			description: '你的冰冷血脉技能可以移除所有移动减速和施法减速效果。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_frost_coldhearted.jpg',
		},
		[MageMajorGlyph.GlyphOfInvisibility]: {
			name: '隐形雕文',
			description: '使你的隐形术的持续时间延长10秒。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_mage_invisibility.jpg',
		},
		[MageMajorGlyph.GlyphOfLivingBomb]: {
			name: '活动炸弹雕文',
			description: '你的活动炸弹的持续伤害可以爆击。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_mage_livingbomb.jpg',
		},
		[MageMajorGlyph.GlyphOfMageArmor]: {
			name: '法师护甲雕文',
			description: '你的法师护甲令你在施法时恢复法力值的效果提高20%。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_magearmor.jpg',
		},
		[MageMajorGlyph.GlyphOfManaGem]: {
			name: '法力宝石雕文',
			description: '使你通过法力宝石恢复的法力值提高40%。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/inv_misc_gem_stone_01.jpg',
		},
		[MageMajorGlyph.GlyphOfMirrorImage]: {
			name: '镜像雕文',
			description: '你的镜像法术可以生成4个镜像。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_magic_lesserinvisibilty.jpg',
		},
		[MageMajorGlyph.GlyphOfMoltenArmor]: {
			name: '熔岩护甲雕文',
			description: '你的熔岩护甲可以提供额外的爆击几率加成，数值相当于你的精神值的20%。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_mage_moltenarmor.jpg',
		},
		[MageMajorGlyph.GlyphOfPolymorph]: {
			name: '变形雕文',
			description: '你的变形术可以移除目标身上的所有持续伤害效果。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_nature_polymorph.jpg',
		},
		[MageMajorGlyph.GlyphOfRemoveCurse]: {
			name: '解除诅咒雕文',
			description: '你的解除诅咒法术令目标对所有诅咒免疫，持续4 sec。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_nature_removecurse.jpg',
		},
		[MageMajorGlyph.GlyphOfScorch]: {
			name: '灼烧雕文',
			description: '你每次施放灼烧，强化灼烧天赋就提供20层强化灼烧效果。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_fire_soulburn.jpg',
		},
		[MageMajorGlyph.GlyphOfWaterElemental]: {
			name: '水元素雕文',
			description: '使你的召唤水元素法术的冷却时间缩短30秒。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_frost_summonwaterelemental_2.jpg',
		},
	},
	minorGlyphs: {
		[MageMinorGlyph.GlyphOfArcaneIntellect]: {
			name: '奥术智慧雕文',
			description: '你的奥术智慧和奥术光辉的法力值消耗降低50%。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_holy_magicalsentry.jpg',
		},
		[MageMinorGlyph.GlyphOfBlastWave]: {
			name: '冲击波雕文',
			description: '你的冲击波消耗的法力值降低15%，但是不再击退敌人。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_holy_excorcism_02.jpg',
		},
		[MageMinorGlyph.GlyphOfFireWard]: {
			name: '防护火焰结界雕文',
			description: '当你的防护火焰结界处于激活状态下时，你反射火焰法术的几率提高5%。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_fire_firearmor.jpg',
		},
		[MageMinorGlyph.GlyphOfFrostArmor]: {
			name: '霜甲雕文',
			description: '使你的冰甲术和霜甲术的持续时间延长30分钟。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_frost_frostarmor02.jpg',
		},
		[MageMinorGlyph.GlyphOfFrostWard]: {
			name: '防护冰霜结界雕文',
			description: '当你的防护冰霜结界处于激活状态下时，你反射冰霜法术的几率提高5%。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_frost_frostward.jpg',
		},
		[MageMinorGlyph.GlyphOfSlowFall]: {
			name: '缓落雕文',
			description: '你的缓落术不再需要施法材料。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_magic_featherfall.jpg',
		},
		[MageMinorGlyph.GlyphOfThePenguin]: {
			name: '企鹅雕文',
			description: '你的变形术：羊会将目标变成一只企鹅。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/inv_misc_penguinpet.jpg',
		},
	},
};
