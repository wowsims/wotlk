import { Class } from '/tbc/core/proto/common.js';
import { ItemSlot } from '/tbc/core/proto/common.js';
import { Race } from '/tbc/core/proto/common.js';
import { ShattrathFaction } from '/tbc/core/proto/common.js';
import { Stat } from '/tbc/core/proto/common.js';
import { ResourceType } from '/tbc/core/proto/api.js';

export const raceNames: Record<Race, string> = {
	[Race.RaceUnknown]: 'None',
	[Race.RaceBloodElf]: 'Blood Elf',
	[Race.RaceDraenei]: 'Draenei',
	[Race.RaceDwarf]: 'Dwarf',
	[Race.RaceGnome]: 'Gnome',
	[Race.RaceHuman]: 'Human',
	[Race.RaceNightElf]: 'Night Elf',
	[Race.RaceOrc]: 'Orc',
	[Race.RaceTauren]: 'Tauren',
	[Race.RaceTroll10]: 'Troll (+10% Haste)',
	[Race.RaceTroll30]: 'Troll (+30% Haste)',
	[Race.RaceUndead]: 'Undead',
};

export function nameToRace(name: string): Race {
	const normalized = name.toLowerCase().replaceAll(' ', '');
	if (normalized.includes('troll')) {
		return Race.RaceTroll10;
	}

	for (const key in raceNames) {
		const race = parseInt(key) as Race;
		if (raceNames[race].toLowerCase().replaceAll(' ', '') == normalized) {
			return race;
		}
	}

	return Race.RaceUnknown;
}

export const classNames: Record<Class, string> = {
	[Class.ClassUnknown]: 'None',
	[Class.ClassDruid]: 'Druid',
	[Class.ClassHunter]: 'Hunter',
	[Class.ClassMage]: 'Mage',
	[Class.ClassPaladin]: 'Paladin',
	[Class.ClassPriest]: 'Priest',
	[Class.ClassRogue]: 'Rogue',
	[Class.ClassShaman]: 'Shaman',
	[Class.ClassWarlock]: 'Warlock',
	[Class.ClassWarrior]: 'Warrior',
}

export function nameToClass(name: string): Class {
	const lower = name.toLowerCase();
	for (const key in classNames) {
		const charClass = parseInt(key) as Class;
		if (classNames[charClass].toLowerCase() == lower) {
			return charClass;
		}
	}

	return Class.ClassUnknown;
}

export const statOrder: Array<Stat> = [
	Stat.StatHealth,
	Stat.StatArmor,
	Stat.StatStamina,
	Stat.StatStrength,
	Stat.StatAgility,
	Stat.StatIntellect,
	Stat.StatSpirit,
	Stat.StatSpellPower,
	Stat.StatHealingPower,
	Stat.StatArcaneSpellPower,
	Stat.StatFireSpellPower,
	Stat.StatFrostSpellPower,
	Stat.StatHolySpellPower,
	Stat.StatNatureSpellPower,
	Stat.StatShadowSpellPower,
	Stat.StatSpellHit,
	Stat.StatSpellCrit,
	Stat.StatSpellHaste,
	Stat.StatSpellPenetration,
	Stat.StatMP5,
	Stat.StatAttackPower,
	Stat.StatRangedAttackPower,
	Stat.StatFeralAttackPower,
	Stat.StatMeleeHit,
	Stat.StatMeleeCrit,
	Stat.StatMeleeHaste,
	Stat.StatArmorPenetration,
	Stat.StatExpertise,
	Stat.StatMana,
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
];

export const statNames: Record<Stat, string> = {
	[Stat.StatStrength]: 'Strength',
	[Stat.StatAgility]: 'Agility',
	[Stat.StatStamina]: 'Stamina',
	[Stat.StatIntellect]: 'Intellect',
	[Stat.StatSpirit]: 'Spirit',
	[Stat.StatSpellPower]: 'Spell Dmg',
	[Stat.StatHealingPower]: 'Healing Power',
	[Stat.StatArcaneSpellPower]: 'Arcane Dmg',
	[Stat.StatFireSpellPower]: 'Fire Dmg',
	[Stat.StatFrostSpellPower]: 'Frost Dmg',
	[Stat.StatHolySpellPower]: 'Holy Dmg',
	[Stat.StatNatureSpellPower]: 'Nature Dmg',
	[Stat.StatShadowSpellPower]: 'Shadow Dmg',
	[Stat.StatMP5]: 'MP5',
	[Stat.StatSpellHit]: 'Spell Hit',
	[Stat.StatSpellCrit]: 'Spell Crit',
	[Stat.StatSpellHaste]: 'Spell Haste',
	[Stat.StatSpellPenetration]: 'Spell Pen',
	[Stat.StatAttackPower]: 'Attack Power',
	[Stat.StatMeleeHit]: 'Hit',
	[Stat.StatMeleeCrit]: 'Crit',
	[Stat.StatMeleeHaste]: 'Haste',
	[Stat.StatArmorPenetration]: 'Armor Pen',
	[Stat.StatExpertise]: 'Expertise',
	[Stat.StatMana]: 'Mana',
	[Stat.StatEnergy]: 'Energy',
	[Stat.StatRage]: 'Rage',
	[Stat.StatArmor]: 'Armor',
	[Stat.StatRangedAttackPower]: 'Ranged AP',
	[Stat.StatFeralAttackPower]: 'Feral AP',
	[Stat.StatDefense]: 'Defense',
	[Stat.StatBlock]: 'Block',
	[Stat.StatBlockValue]: 'Block Value',
	[Stat.StatDodge]: 'Dodge',
	[Stat.StatParry]: 'Parry',
	[Stat.StatResilience]: 'Resilience',
	[Stat.StatHealth]: 'Health',
	[Stat.StatArcaneResistance]: 'Arcane Resistance',
	[Stat.StatFireResistance]: 'Fire Resistance',
	[Stat.StatFrostResistance]: 'Frost Resistance',
	[Stat.StatNatureResistance]: 'Nature Resistance',
	[Stat.StatShadowResistance]: 'Shadow Resistance',
};

export const slotNames: Record<ItemSlot, string> = {
	[ItemSlot.ItemSlotHead]: 'Head',
	[ItemSlot.ItemSlotNeck]: 'Neck',
	[ItemSlot.ItemSlotShoulder]: 'Shoulders',
	[ItemSlot.ItemSlotBack]: 'Back',
	[ItemSlot.ItemSlotChest]: 'Chest',
	[ItemSlot.ItemSlotWrist]: 'Wrist',
	[ItemSlot.ItemSlotHands]: 'Hands',
	[ItemSlot.ItemSlotWaist]: 'Waist',
	[ItemSlot.ItemSlotLegs]: 'Legs',
	[ItemSlot.ItemSlotFeet]: 'Feet',
	[ItemSlot.ItemSlotFinger1]: 'Finger 1',
	[ItemSlot.ItemSlotFinger2]: 'Finger 2',
	[ItemSlot.ItemSlotTrinket1]: 'Trinket 1',
	[ItemSlot.ItemSlotTrinket2]: 'Trinket 2',
	[ItemSlot.ItemSlotMainHand]: 'Main Hand',
	[ItemSlot.ItemSlotOffHand]: 'Off Hand',
	[ItemSlot.ItemSlotRanged]: 'Ranged',
};

export const resourceNames: Record<ResourceType, string> = {
	[ResourceType.ResourceTypeNone]: 'None',
	[ResourceType.ResourceTypeHealth]: 'Health',
	[ResourceType.ResourceTypeMana]: 'Mana',
	[ResourceType.ResourceTypeEnergy]: 'Energy',
	[ResourceType.ResourceTypeRage]: 'Rage',
	[ResourceType.ResourceTypeComboPoints]: 'Combo Points',
	[ResourceType.ResourceTypeFocus]: 'Focus',
};

export const resourceColors: Record<ResourceType, string> = {
	[ResourceType.ResourceTypeNone]: '#ffffff',
	[ResourceType.ResourceTypeHealth]: '#22ba00',
	[ResourceType.ResourceTypeMana]: '#2e93fa',
	[ResourceType.ResourceTypeEnergy]: '#ffd700',
	[ResourceType.ResourceTypeRage]: '#ff0000',
	[ResourceType.ResourceTypeComboPoints]: '#ffa07a',
	[ResourceType.ResourceTypeFocus]: '#cd853f',
};

export function stringToResourceType(str: string): ResourceType {
	for (const [key, val] of Object.entries(resourceNames)) {
		if (val.toLowerCase() == str.toLowerCase()) {
			return Number(key) as ResourceType;
		}
	}
	return ResourceType.ResourceTypeNone;
}

export const shattFactionNames: Record<ShattrathFaction, string> = {
	[ShattrathFaction.ShattrathFactionAldor]: 'Aldor',
	[ShattrathFaction.ShattrathFactionScryer]: 'Scryer',
};
