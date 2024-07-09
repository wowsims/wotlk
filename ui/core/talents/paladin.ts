import { PaladinMajorGlyph, PaladinMinorGlyph,PaladinTalents } from '../proto/paladin.js';
import { GlyphsConfig } from './glyphs_picker.js';
import { newTalentsConfig,TalentsConfig } from './talents_picker.js';
import PaladinTalentJson from './trees/paladin.json';

export const paladinTalentsConfig: TalentsConfig<PaladinTalents> = newTalentsConfig(PaladinTalentJson);

export const paladinGlyphsConfig: GlyphsConfig = {
	majorGlyphs: {
		[PaladinMajorGlyph.GlyphOfAvengerSShield]: {
			name: '复仇者之盾雕文',
			description: '你的复仇者之盾攻击的目标数量减少2个，但伤害提高100%。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_holy_avengersshield.jpg',
		},
		[PaladinMajorGlyph.GlyphOfAvengingWrath]: {
			name: '复仇之怒雕文',
			description: '当复仇之怒处于激活状态时，使你的愤怒之锤的冷却时间缩短50%。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_holy_avenginewrath.jpg',
		},
		[PaladinMajorGlyph.GlyphOfBeaconOfLight]: {
			name: '圣光道标雕文',
			description: '圣光道标的持续时间延长30秒。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_paladin_beaconoflight.jpg',
		},
		[PaladinMajorGlyph.GlyphOfCleansing]: {
			name: '清洁雕文',
			description: '使你的清洁术和纯净术的法力值消耗降低20%。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_holy_purify.jpg',
		},
		[PaladinMajorGlyph.GlyphOfConsecration]: {
			name: '奉献雕文',
			description: '使你的奉献的持续时间和冷却时间延长2秒。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_holy_innerfire.jpg',
		},
		[PaladinMajorGlyph.GlyphOfCrusaderStrike]: {
			name: '十字军打击雕文',
			description: '你的十字军打击的法力值消耗降低20%。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_holy_crusaderstrike.jpg',
		},
		[PaladinMajorGlyph.GlyphOfDivinePlea]: {
			name: '神圣恳求雕文',
			description: '当神圣恳求处于激活状态时，你受到的所有伤害降低3%。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_holy_aspiration.jpg',
		},
		[PaladinMajorGlyph.GlyphOfDivineStorm]: {
			name: '神圣风暴雕文',
			description: '你的神圣风暴所造成的伤害产生的治疗效果提高15%。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_paladin_divinestorm.jpg',
		},
		[PaladinMajorGlyph.GlyphOfDivinity]: {
			name: '圣洁雕文',
			description: '你的圣疗术可以恢复双倍的法力值，并且在对其他队友施放时也可以为你恢复法力值，数值相当于目标所获得的法力值。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_holy_layonhands.jpg',
		},
		[PaladinMajorGlyph.GlyphOfExorcism]: {
			name: '驱邪雕文',
			description: '使你的驱邪术造成的伤害提高20%。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_holy_excorcism_02.jpg',
		},
		[PaladinMajorGlyph.GlyphOfFlashOfLight]: {
			name: '圣光闪现雕文',
			description: '你的圣光闪现的爆击几率提高5%。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_holy_flashheal.jpg',
		},
		[PaladinMajorGlyph.GlyphOfHammerOfJustice]: {
			name: '制裁之锤雕文',
			description: '使你的制裁之锤的射程延长5码。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_holy_sealofmight.jpg',
		},
		[PaladinMajorGlyph.GlyphOfHammerOfTheRighteous]: {
			name: '正义之锤雕文',
			description: '你的正义之锤攻击的目标增加1个。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_paladin_hammeroftherighteous.jpg',
		},
		[PaladinMajorGlyph.GlyphOfHammerOfWrath]: {
			name: '愤怒之锤雕文',
			description: '你的愤怒之锤消耗的法力值降低100%。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_thunderclap.jpg',
		},
		[PaladinMajorGlyph.GlyphOfHolyLight]: {
			name: '圣光雕文',
			description: '你的圣光术令目标周围半径8码内的最多5个友方目标获得治疗，数值相当于该次圣光术治疗量的10%。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_holy_holybolt.jpg',
		},
		[PaladinMajorGlyph.GlyphOfHolyShock]: {
			name: '神圣震击雕文',
			description: '神圣震击的冷却时间缩短1秒。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_holy_searinglight.jpg',
		},
		[PaladinMajorGlyph.GlyphOfHolyWrath]: {
			name: '神圣愤怒雕文',
			description: '你的神圣愤怒法术的冷却时间减少15秒。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_holy_excorcism.jpg',
		},
		[PaladinMajorGlyph.GlyphOfJudgement]: {
			name: '审判雕文',
			description: '你的审判法术造成的伤害提高10%。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_paladin_judgementred.jpg',
		},
		[PaladinMajorGlyph.GlyphOfRighteousDefense]: {
			name: '正义防御雕文',
			description: '使你的正义防御和清算之手技能对每个目标生效的几率都提高8%。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/inv_shoulder_37.jpg',
		},
		[PaladinMajorGlyph.GlyphOfSalvation]: {
			name: '拯救雕文',
			description: '当你对自己施放拯救之手时，你受到的伤害也降低20%。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_holy_sealofsalvation.jpg',
		},
		[PaladinMajorGlyph.GlyphOfSealOfCommand]: {
			name: '命令圣印雕文',
			description: '每次当你在命令圣印激活时使用审判，你恢复基础法力值的8%。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_warrior_innerrage.jpg',
		},
		[PaladinMajorGlyph.GlyphOfSealOfLight]: {
			name: '光明圣印雕文',
			description: '当光明圣印处于激活状态下时，你的治疗法术的效果提高5%。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_holy_healingaura.jpg',
		},
		[PaladinMajorGlyph.GlyphOfSealOfRighteousness]: {
			name: '正义圣印雕文',
			description: '正义圣印造成的伤害提高10%。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_thunderbolt.jpg',
		},
		[PaladinMajorGlyph.GlyphOfSealOfVengeance]: {
			name: '复仇圣印雕文',
			description: '你的复仇圣印和腐蚀圣印在激活状态下可以为你提供10精准。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_holy_sealofvengeance.jpg',
		},
		[PaladinMajorGlyph.GlyphOfSealOfWisdom]: {
			name: '智慧圣印雕文',
			description: '当智慧圣印处于激活状态下时，你的治疗法术消耗的法力值降低5%。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_holy_righteousnessaura.jpg',
		},
		[PaladinMajorGlyph.GlyphOfShieldOfRighteousness]: {
			name: '正义盾击雕文',
			description: '正义盾击的法力值消耗降低80%。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_paladin_shieldofvengeance.jpg',
		},
		[PaladinMajorGlyph.GlyphOfSpiritualAttunement]: {
			name: '灵魂协调雕文',
			description: '使你通过灵魂协调获得的法力值提高2%。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_holy_revivechampion.jpg',
		},
		[PaladinMajorGlyph.GlyphOfTurnEvil]: {
			name: '超度邪恶雕文',
			description: '使你的超度邪恶的施法时间缩短100%，但是冷却时间延长8秒。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_holy_turnundead.jpg',
		},
		[PaladinMajorGlyph.GlyphOfReckoning]: {
			name: '清算雕文',
			description: '使你的清算之手法术不再嘲讽目标并且可以对无法嘲讽的目标造成伤害。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_holy_unyieldingfaith.jpg',
		},
	},
	minorGlyphs: {
		[PaladinMinorGlyph.GlyphOfBlessingOfKings]: {
			name: '王者祝福雕文',
			description: '你的王者祝福和强效王者祝福的法力值消耗降低50%。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_magic_magearmor.jpg',
		},
		[PaladinMinorGlyph.GlyphOfBlessingOfMight]: {
			name: '力量祝福雕文',
			description: '使你对自己施放的力量祝福的持续时间延长20分钟。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_holy_fistofjustice.jpg',
		},
		[PaladinMinorGlyph.GlyphOfBlessingOfWisdom]: {
			name: '智慧祝福雕文',
			description: '使你对自己施放的智慧祝福的持续时间延长20分钟。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_holy_sealofwisdom.jpg',
		},
		[PaladinMinorGlyph.GlyphOfLayOnHands]: {
			name: '圣疗雕文',
			description: '使你的圣疗术的冷却时间缩短5分钟。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_holy_layonhands.jpg',
		},
		[PaladinMinorGlyph.GlyphOfSenseUndead]: {
			name: '感知亡灵雕文',
			description: '当你的感知亡灵技能处于激活状态下时，使你对亡灵造成的伤害提高1%。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_holy_senseundead.jpg',
		},
		[PaladinMinorGlyph.GlyphOfTheWise]: {
			name: '智者雕文',
			description: '你的智慧圣印的法力值消耗降低50%。',
			iconUrl: 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_holy_righteousnessaura.jpg',
		},
	},
};
