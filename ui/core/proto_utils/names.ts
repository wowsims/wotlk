import { ResourceType } from '../proto/api.js';
import { ArmorType, Class, ItemSlot, Profession, PseudoStat, Race, RangedWeaponType, Stat, WeaponType } from '../proto/common.js';
import { DungeonDifficulty, RaidFilterOption, RepFaction, RepLevel, SourceFilterOption } from '../proto/ui.js';

export const armorTypeNames: Map<ArmorType, string> = new Map([
	[ArmorType.ArmorTypeUnknown, '未知'],
	[ArmorType.ArmorTypeCloth, '布甲'],
	[ArmorType.ArmorTypeLeather, '皮甲'],
	[ArmorType.ArmorTypeMail, '锁甲'],
	[ArmorType.ArmorTypePlate, '板甲'],
]);

export const weaponTypeNames: Map<WeaponType, string> = new Map([
	[WeaponType.WeaponTypeUnknown, '未知'],
	[WeaponType.WeaponTypeAxe, '斧'],
	[WeaponType.WeaponTypeDagger, '匕首'],
	[WeaponType.WeaponTypeFist, '拳套'],
	[WeaponType.WeaponTypeMace, '锤'],
	[WeaponType.WeaponTypeOffHand, '副手物品'],
	[WeaponType.WeaponTypePolearm, '长柄武器'],
	[WeaponType.WeaponTypeShield, '盾牌'],
	[WeaponType.WeaponTypeStaff, '法杖'],
	[WeaponType.WeaponTypeSword, '剑'],
]);

export const rangedWeaponTypeNames: Map<RangedWeaponType, string> = new Map([
	[RangedWeaponType.RangedWeaponTypeUnknown, '未知'],
	[RangedWeaponType.RangedWeaponTypeBow, '弓'],
	[RangedWeaponType.RangedWeaponTypeCrossbow, '弩'],
	[RangedWeaponType.RangedWeaponTypeGun, '枪'],
	[RangedWeaponType.RangedWeaponTypeIdol, '神像'],
	[RangedWeaponType.RangedWeaponTypeLibram, '圣契'],
	[RangedWeaponType.RangedWeaponTypeSigil, '符印'],
	[RangedWeaponType.RangedWeaponTypeThrown, '投掷武器'],
	[RangedWeaponType.RangedWeaponTypeTotem, '图腾'],
	[RangedWeaponType.RangedWeaponTypeWand, '魔杖'],
]);

export const raceNames: Map<Race, string> = new Map([
	[Race.RaceUnknown, 'None'],
	[Race.RaceBloodElf, 'Blood Elf'],
	[Race.RaceDraenei, 'Draenei'],
	[Race.RaceDwarf, 'Dwarf'],
	[Race.RaceGnome, 'Gnome'],
	[Race.RaceHuman, 'Human'],
	[Race.RaceNightElf, 'Night Elf'],
	[Race.RaceOrc, 'Orc'],
	[Race.RaceTauren, 'Tauren'],
	[Race.RaceTroll, 'Troll'],
	[Race.RaceUndead, 'Undead'],
]);

export const raceNamesCn: Map<Race, string> = new Map([
	[Race.RaceUnknown, '无'],
	[Race.RaceBloodElf, '血精灵'],
	[Race.RaceDraenei, '德莱尼'],
	[Race.RaceDwarf, '矮人'],
	[Race.RaceGnome, '侏儒'],
	[Race.RaceHuman, '人类'],
	[Race.RaceNightElf, '暗夜精灵'],
	[Race.RaceOrc, '兽人'],
	[Race.RaceTauren, '牛头人'],
	[Race.RaceTroll, '巨魔'],
	[Race.RaceUndead, '亡灵'],
]);

export function nameToRace(name: string): Race {
	const normalized = name.toLowerCase().replaceAll(' ', '');
	for (const [key, value] of raceNames) {
		if (value.toLowerCase().replaceAll(' ', '') == normalized) {
			return key;
		}
	}
	return Race.RaceUnknown;
}

export const classNames: Map<Class, string> = new Map([
	[Class.ClassUnknown, 'None'],
	[Class.ClassDruid, 'Druid'],
	[Class.ClassHunter, 'Hunter'],
	[Class.ClassMage, 'Mage'],
	[Class.ClassPaladin, 'Paladin'],
	[Class.ClassPriest, 'Priest'],
	[Class.ClassRogue, 'Rogue'],
	[Class.ClassShaman, 'Shaman'],
	[Class.ClassWarlock, 'Warlock'],
	[Class.ClassWarrior, 'Warrior'],
	[Class.ClassDeathknight, 'Death Knight'],
]);

export const classNamesCn: Map<Class, string> = new Map([
	[Class.ClassUnknown, '无'],
	[Class.ClassDruid, '德鲁伊'],
	[Class.ClassHunter, '猎人'],
	[Class.ClassMage, '法师'],
	[Class.ClassPaladin, '圣骑士'],
	[Class.ClassPriest, '牧师'],
	[Class.ClassRogue, '潜行者'],
	[Class.ClassShaman, '萨满'],
	[Class.ClassWarlock, '术士'],
	[Class.ClassWarrior, '战士'],
	[Class.ClassDeathknight, '死亡骑士'],
]);

export function nameToClass(name: string): Class {
	const lower = name.toLowerCase();
	for (const [key, value] of classNames) {
		if (value.toLowerCase().replace(/\s+/g, '') == lower) {
			return key;
		}
	}
	return Class.ClassUnknown;
}

export const professionNames: Map<Profession, string> = new Map([
	[Profession.ProfessionUnknown, 'None'],
	[Profession.Alchemy, 'Alchemy'],
	[Profession.Blacksmithing, 'Blacksmithing'],
	[Profession.Enchanting, 'Enchanting'],
	[Profession.Engineering, 'Engineering'],
	[Profession.Herbalism, 'Herbalism'],
	[Profession.Inscription, 'Inscription'],
	[Profession.Jewelcrafting, 'Jewelcrafting'],
	[Profession.Leatherworking, 'Leatherworking'],
	[Profession.Mining, 'Mining'],
	[Profession.Skinning, 'Skinning'],
	[Profession.Tailoring, 'Tailoring'],
]);

export const professionNamesCn: Map<Profession, string> = new Map([
	[Profession.ProfessionUnknown, '无'],
	[Profession.Alchemy, '炼金'],
	[Profession.Blacksmithing, '锻造'],
	[Profession.Enchanting, '附魔'],
	[Profession.Engineering, '工程学'],
	[Profession.Herbalism, '草药学'],
	[Profession.Inscription, '铭文'],
	[Profession.Jewelcrafting, '珠宝加工'],
	[Profession.Leatherworking, '制皮'],
	[Profession.Mining, '采矿'],
	[Profession.Skinning, '剥皮'],
	[Profession.Tailoring, '裁缝'],
]);

export function nameToProfession(name: string): Profession {
	const lower = name.toLowerCase();
	for (const [key, value] of professionNames) {
		if (value.toLowerCase() == lower) {
			return key;
		}
	}
	return Profession.ProfessionUnknown;
}

export const statOrder: Array<Stat> = [
	Stat.StatHealth,
	Stat.StatMana,
	Stat.StatArmor,
	Stat.StatBonusArmor,
	Stat.StatStamina,
	Stat.StatStrength,
	Stat.StatAgility,
	Stat.StatIntellect,
	Stat.StatSpirit,
	Stat.StatSpellPower,
	Stat.StatSpellHit,
	Stat.StatSpellCrit,
	Stat.StatSpellHaste,
	Stat.StatSpellPenetration,
	Stat.StatMP5,
	Stat.StatAttackPower,
	Stat.StatRangedAttackPower,
	Stat.StatMeleeHit,
	Stat.StatMeleeCrit,
	Stat.StatMeleeHaste,
	Stat.StatArmorPenetration,
	Stat.StatExpertise,
	Stat.StatEnergy,
	Stat.StatRage,
	Stat.StatDefense,
	Stat.StatBlock,
	Stat.StatBlockValue,
	Stat.StatDodge,
	Stat.StatParry,
	Stat.StatResilience,
	Stat.StatArcaneResistance,
	Stat.StatFireResistance,
	Stat.StatFrostResistance,
	Stat.StatNatureResistance,
	Stat.StatShadowResistance,
	Stat.StatRunicPower,
	Stat.StatBloodRune,
	Stat.StatFrostRune,
	Stat.StatUnholyRune,
	Stat.StatDeathRune,
];

export const statNames: Map<Stat, string> = new Map([
	[Stat.StatStrength, '力量'],
	[Stat.StatAgility, '敏捷'],
	[Stat.StatStamina, '耐力'],
	[Stat.StatIntellect, '智力'],
	[Stat.StatSpirit, '精神'],
	[Stat.StatSpellPower, '法术伤害'],
	[Stat.StatMP5, '每5秒回蓝'],
	[Stat.StatSpellHit, '法术命中'],
	[Stat.StatSpellCrit, '法术暴击'],
	[Stat.StatSpellHaste, '法术急速'],
	[Stat.StatSpellPenetration, '法术穿透'],
	[Stat.StatAttackPower, '攻击强度'],
	[Stat.StatMeleeHit, '近战命中'],
	[Stat.StatMeleeCrit, '近战暴击'],
	[Stat.StatMeleeHaste, '近战急速'],
	[Stat.StatArmorPenetration, '护甲穿透'],
	[Stat.StatExpertise, '精准'],
	[Stat.StatMana, '法力值'],
	[Stat.StatEnergy, '能量'],
	[Stat.StatRage, '怒气'],
	[Stat.StatArmor, '护甲'],
	[Stat.StatRangedAttackPower, '远程攻击强度'],
	[Stat.StatDefense, '防御'],
	[Stat.StatBlock, '格挡'],
	[Stat.StatBlockValue, '格挡值'],
	[Stat.StatDodge, '躲闪'],
	[Stat.StatParry, '招架'],
	[Stat.StatResilience, '韧性'],
	[Stat.StatHealth, '生命值'],
	[Stat.StatArcaneResistance, '奥术抗性'],
	[Stat.StatFireResistance, '火焰抗性'],
	[Stat.StatFrostResistance, '冰霜抗性'],
	[Stat.StatNatureResistance, '自然抗性'],
	[Stat.StatShadowResistance, '暗影抗性'],
	[Stat.StatBonusArmor, '额外护甲'],
	[Stat.StatRunicPower, '符能'],
	[Stat.StatBloodRune, '血符文'],
	[Stat.StatFrostRune, '冰符文'],
	[Stat.StatUnholyRune, '邪符文'],
	[Stat.StatDeathRune, '死符文'],
]);

export const pseudoStatOrder: Array<PseudoStat> = [
	PseudoStat.PseudoStatMainHandDps,
	PseudoStat.PseudoStatOffHandDps,
	PseudoStat.PseudoStatRangedDps,
	PseudoStat.PseudoStatBlockValueMultiplier,
];

export const pseudoStatNames: Map<PseudoStat, string> = new Map([
	[PseudoStat.PseudoStatMainHandDps, '主手DPS'],
	[PseudoStat.PseudoStatOffHandDps, '副手DPS'],
	[PseudoStat.PseudoStatRangedDps, '远程DPS'],
	[PseudoStat.PseudoStatBlockValueMultiplier, '格挡值倍数'],
	[PseudoStat.PseudoStatDodge, '躲闪几率'],
	[PseudoStat.PseudoStatParry, '招架几率'],
]);

export function getClassStatName(stat: Stat, playerClass: Class): string {
	const statName = statNames.get(stat);
	if (!statName) return 'UnknownStat';
	if (playerClass == Class.ClassHunter) {
		return statName.replace('近战', '远程');
	} else {
		return statName;
	}
}

export const slotNames: Map<ItemSlot, string> = new Map([
	[ItemSlot.ItemSlotHead, '头部'],
	[ItemSlot.ItemSlotNeck, '颈部'],
	[ItemSlot.ItemSlotShoulder, '肩部'],
	[ItemSlot.ItemSlotBack, '背部'],
	[ItemSlot.ItemSlotChest, '胸部'],
	[ItemSlot.ItemSlotWrist, '手腕'],
	[ItemSlot.ItemSlotHands, '手'],
	[ItemSlot.ItemSlotWaist, '腰部'],
	[ItemSlot.ItemSlotLegs, '腿部'],
	[ItemSlot.ItemSlotFeet, '脚'],
	[ItemSlot.ItemSlotFinger1, '戒指1'],
	[ItemSlot.ItemSlotFinger2, '戒指2'],
	[ItemSlot.ItemSlotTrinket1, '饰品1'],
	[ItemSlot.ItemSlotTrinket2, '饰品2'],
	[ItemSlot.ItemSlotMainHand, '主手'],
	[ItemSlot.ItemSlotOffHand, '副手'],
	[ItemSlot.ItemSlotRanged, '远程'],
]);

export const resourceNames: Map<ResourceType, string> = new Map([
	[ResourceType.ResourceTypeNone, '无'],
	[ResourceType.ResourceTypeHealth, '生命值'],
	[ResourceType.ResourceTypeMana, '法力'],
	[ResourceType.ResourceTypeEnergy, '能量'],
	[ResourceType.ResourceTypeRage, '怒气'],
	[ResourceType.ResourceTypeComboPoints, '连击点'],
	[ResourceType.ResourceTypeFocus, '集中值'],
	[ResourceType.ResourceTypeRunicPower, '符文能量'],
	[ResourceType.ResourceTypeBloodRune, '血符文'],
	[ResourceType.ResourceTypeFrostRune, '冰符文'],
	[ResourceType.ResourceTypeUnholyRune, '邪符文'],
	[ResourceType.ResourceTypeDeathRune, '死亡符文'],
]);

export const resourceColors: Map<ResourceType, string> = new Map([
	[ResourceType.ResourceTypeNone, '#ffffff'],
	[ResourceType.ResourceTypeHealth, '#22ba00'],
	[ResourceType.ResourceTypeMana, '#2e93fa'],
	[ResourceType.ResourceTypeEnergy, '#ffd700'],
	[ResourceType.ResourceTypeRage, '#ff0000'],
	[ResourceType.ResourceTypeComboPoints, '#ffa07a'],
	[ResourceType.ResourceTypeFocus, '#cd853f'],
	[ResourceType.ResourceTypeRunicPower, '#5b99ee'],
	[ResourceType.ResourceTypeBloodRune, '#ff0000'],
	[ResourceType.ResourceTypeFrostRune, '#0000ff'],
	[ResourceType.ResourceTypeUnholyRune, '#00ff00'],
	[ResourceType.ResourceTypeDeathRune, '#8b008b'],
]);

export function stringToResourceType(str: string): ResourceType {
	for (const [key, val] of resourceNames) {
		if (val.toLowerCase() == str.toLowerCase()) {
			return key;
		}
	}
	return ResourceType.ResourceTypeNone;
}

export const sourceNames: Map<SourceFilterOption, string> = new Map([
	[SourceFilterOption.SourceUnknown, '未知'],
	[SourceFilterOption.SourceCrafting, '制造业'],
	[SourceFilterOption.SourceQuest, '任务'],
	[SourceFilterOption.SourceDungeon, '地下城'],
	[SourceFilterOption.SourceDungeonH, '地下城（英雄）'],
	[SourceFilterOption.SourceDungeonTRA, '地下城（TRA）'],
	[SourceFilterOption.SourceDungeonTRB, '地下城（TRB）'],
	[SourceFilterOption.SourceRaid10, '团队副本（10人普通）'],
	[SourceFilterOption.SourceRaid10H, '团队副本（10人英雄）'],
	[SourceFilterOption.SourceRaid25, '团队副本（25人普通）'],
	[SourceFilterOption.SourceRaid25H, '团队副本（25人英雄）'],
]);

export const raidNames: Map<RaidFilterOption, string> = new Map([
	[RaidFilterOption.RaidUnknown, '未知'],
	[RaidFilterOption.RaidVanilla, '经典旧世'],
	[RaidFilterOption.RaidTbc, '燃烧的远征'],
	[RaidFilterOption.RaidNaxxramas, '纳克萨玛斯'],
	[RaidFilterOption.RaidEyeOfEternity, '永恒之眼'],
	[RaidFilterOption.RaidObsidianSanctum, '黑曜石圣殿'],
	[RaidFilterOption.RaidVaultOfArchavon, '阿尔卡冯的宝库'],
	[RaidFilterOption.RaidUlduar, '奥杜尔'],
	[RaidFilterOption.RaidTrialOfTheCrusader, '十字军的试炼'],
	[RaidFilterOption.RaidOnyxiasLair, '奥妮克希亚的巢穴'],
	[RaidFilterOption.RaidIcecrownCitadel, '冰冠堡垒'],
	[RaidFilterOption.RaidRubySanctum, '红玉圣殿'],
]);

export const difficultyNames: Map<DungeonDifficulty, string> = new Map([
	[DungeonDifficulty.DifficultyUnknown, 'Unknown'],
	[DungeonDifficulty.DifficultyNormal, 'N'],
	[DungeonDifficulty.DifficultyHeroic, 'H'],
	[DungeonDifficulty.DifficultyTitanRuneAlpha, 'TRA'],
	[DungeonDifficulty.DifficultyTitanRuneBeta, 'TRB'],
	[DungeonDifficulty.DifficultyRaid10, '10N'],
	[DungeonDifficulty.DifficultyRaid10H, '10H'],
	[DungeonDifficulty.DifficultyRaid25, '25N'],
	[DungeonDifficulty.DifficultyRaid25H, '25H'],
]);

export const REP_LEVEL_NAMES: Record<RepLevel, string> = {
	[RepLevel.RepLevelUnknown]: '未知',
	[RepLevel.RepLevelHated]: '仇恨',
	[RepLevel.RepLevelHostile]: '敌对',
	[RepLevel.RepLevelUnfriendly]: '冷淡',
	[RepLevel.RepLevelNeutral]: '中立',
	[RepLevel.RepLevelFriendly]: '友好',
	[RepLevel.RepLevelHonored]: '尊敬',
	[RepLevel.RepLevelRevered]: '崇敬',
	[RepLevel.RepLevelExalted]: '崇拜',
};

export const REP_FACTION_NAMES: Record<RepFaction, string> = {
	[RepFaction.RepFactionUnknown]: 'Unknown',
};
