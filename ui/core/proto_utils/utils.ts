import { REPO_NAME } from '../constants/other.js';
import { camelToSnakeCase, getEnumValues, intersection, maxIndex, sum } from '../utils.js';

import { Player as PlayerProto, ResourceType } from '../proto/api.js';
import { ArmorType, Class, EnchantType, Faction, HandType, ItemSlot, ItemType, Race, RangedWeaponType, Spec, UnitReference, UnitReference_Type, WeaponType } from '../proto/common.js';
import { Blessings } from '../proto/paladin.js';
import {
	BlessingsAssignment,
	BlessingsAssignments,
	UIEnchant as Enchant,
	UIItem as Item,
} from '../proto/ui.js';


import { Player } from '../player.js';
import {
	BalanceDruid,
	BalanceDruid_Options as BalanceDruidOptions,
	BalanceDruid_Rotation as BalanceDruidRotation,
	DruidTalents,
	FeralDruid,
	FeralDruid_Options as FeralDruidOptions,
	FeralDruid_Rotation as FeralDruidRotation,
	FeralTankDruid,
	FeralTankDruid_Options as FeralTankDruidOptions,
	FeralTankDruid_Rotation as FeralTankDruidRotation,
	RestorationDruid,
	RestorationDruid_Options as RestorationDruidOptions,
	RestorationDruid_Rotation as RestorationDruidRotation,
} from '../proto/druid.js';
import { Hunter, Hunter_Options as HunterOptions, Hunter_Rotation as HunterRotation, HunterTalents } from '../proto/hunter.js';
import { Mage, Mage_Rotation as MageRotation, Mage_Options as MageOptions, MageTalents } from '../proto/mage.js';
import {
	HolyPaladin,
	HolyPaladin_Options as HolyPaladinOptions,
	HolyPaladin_Rotation as HolyPaladinRotation,
	PaladinTalents,
	ProtectionPaladin,
	ProtectionPaladin_Options as ProtectionPaladinOptions,
	ProtectionPaladin_Rotation as ProtectionPaladinRotation,
	RetributionPaladin,
	RetributionPaladin_Options as RetributionPaladinOptions,
	RetributionPaladin_Rotation as RetributionPaladinRotation,
} from '../proto/paladin.js';
import {
	HealingPriest,
	HealingPriest_Options as HealingPriestOptions,
	HealingPriest_Rotation as HealingPriestRotation,
	PriestTalents,
	ShadowPriest,
	ShadowPriest_Options as ShadowPriestOptions,
	ShadowPriest_Rotation as ShadowPriestRotation,
} from '../proto/priest.js';
import { Rogue, Rogue_Options as RogueOptions, Rogue_Rotation as RogueRotation, RogueTalents } from '../proto/rogue.js';
import {
	ElementalShaman,
	ElementalShaman_Options as ElementalShamanOptions,
	ElementalShaman_Rotation as ElementalShamanRotation,
	EnhancementShaman,
	EnhancementShaman_Options as EnhancementShamanOptions,
	EnhancementShaman_Rotation as EnhancementShamanRotation,
	RestorationShaman,
	RestorationShaman_Options as RestorationShamanOptions,
	RestorationShaman_Rotation as RestorationShamanRotation,
	ShamanTalents,
} from '../proto/shaman.js';
import { 
	Warlock, 
	TankWarlock, 
	WarlockOptions, 
	WarlockRotation, 
	WarlockTalents 
} from '../proto/warlock.js';
import { 
	ProtectionWarrior, 
	ProtectionWarrior_Options as ProtectionWarriorOptions, 
	ProtectionWarrior_Rotation as ProtectionWarriorRotation, 
	Warrior, 
	Warrior_Options as WarriorOptions, 
	Warrior_Rotation as WarriorRotation, 
	WarriorTalents 
} from '../proto/warrior.js';

export type DruidSpecs = Spec.SpecBalanceDruid | Spec.SpecFeralDruid | Spec.SpecFeralTankDruid | Spec.SpecRestorationDruid;
export type HunterSpecs = Spec.SpecHunter;
export type MageSpecs = Spec.SpecMage;
export type PaladinSpecs = Spec.SpecHolyPaladin | Spec.SpecRetributionPaladin | Spec.SpecProtectionPaladin;
export type PriestSpecs = Spec.SpecHealingPriest | Spec.SpecShadowPriest;
export type RogueSpecs = Spec.SpecRogue;
export type ShamanSpecs = Spec.SpecElementalShaman | Spec.SpecEnhancementShaman | Spec.SpecRestorationShaman;
export type WarlockSpecs = Spec.SpecWarlock | Spec.SpecTankWarlock;
export type WarriorSpecs = Spec.SpecWarrior | Spec.SpecProtectionWarrior;

export type ClassSpecs<T extends Class> =
	T extends Class.ClassDruid ? DruidSpecs :
	T extends Class.ClassHunter ? HunterSpecs :
	T extends Class.ClassMage ? MageSpecs :
	T extends Class.ClassPaladin ? PaladinSpecs :
	T extends Class.ClassPriest ? PriestSpecs :
	T extends Class.ClassRogue ? RogueSpecs :
	T extends Class.ClassShaman ? ShamanSpecs :
	T extends Class.ClassWarlock ? WarlockSpecs :
	T extends Class.ClassWarrior ? WarriorSpecs :
	ShamanSpecs; // Should never reach this case

export const NUM_SPECS = getEnumValues(Spec).length;

// The order in which specs should be presented, when it matters.
// Currently this is only used for the order of the paladin blessings UI.
export const naturalSpecOrder: Array<Spec> = [
	Spec.SpecBalanceDruid,
	Spec.SpecFeralDruid,
	Spec.SpecFeralTankDruid,
	Spec.SpecRestorationDruid,
	Spec.SpecHunter,
	Spec.SpecMage,
	Spec.SpecHolyPaladin,
	Spec.SpecProtectionPaladin,
	Spec.SpecRetributionPaladin,
	Spec.SpecHealingPriest,
	Spec.SpecShadowPriest,
	Spec.SpecRogue,
	Spec.SpecElementalShaman,
	Spec.SpecEnhancementShaman,
	Spec.SpecRestorationShaman,
	Spec.SpecWarlock,
	Spec.SpecTankWarlock,
	Spec.SpecWarrior,
	Spec.SpecProtectionWarrior,
];

export const naturalClassOrder: Array<Class> = [
	Class.ClassDruid,
	Class.ClassHunter,
	Class.ClassMage,
	Class.ClassPaladin,
	Class.ClassPriest,
	Class.ClassRogue,
	Class.ClassShaman,
	Class.ClassWarlock,
	Class.ClassWarrior,
]

export const specNames: Record<Spec, string> = {
	[Spec.SpecBalanceDruid]: 'Balance Druid',
	[Spec.SpecFeralDruid]: 'Feral DPS Druid',
	[Spec.SpecFeralTankDruid]: 'Feral Tank Druid',
	[Spec.SpecRestorationDruid]: 'Restoration Druid',
	[Spec.SpecElementalShaman]: 'Elemental Shaman',
	[Spec.SpecEnhancementShaman]: 'Enhancement Shaman',
	[Spec.SpecRestorationShaman]: 'Restoration Shaman',
	[Spec.SpecHunter]: 'Hunter',
	[Spec.SpecMage]: 'Mage',
	[Spec.SpecRogue]: 'Rogue',
	[Spec.SpecHolyPaladin]: 'Holy Paladin',
	[Spec.SpecProtectionPaladin]: 'Protection Paladin',
	[Spec.SpecRetributionPaladin]: 'Retribution Paladin',
	[Spec.SpecHealingPriest]: 'Priest',
	[Spec.SpecShadowPriest]: 'Shadow Priest',
	[Spec.SpecWarlock]: 'DPS Warlock',
	[Spec.SpecTankWarlock]: 'Tank Warlock',
	[Spec.SpecWarrior]: 'DPS Warrior',
	[Spec.SpecProtectionWarrior]: 'Protection Warrior',
};

export const classNames: Record<Class, string> = {
	[Class.ClassUnknown]: '',
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

export const classColors: Record<Class, string> = {
	[Class.ClassUnknown]: '#fff',
	[Class.ClassDruid]: '#ff7d0a',
	[Class.ClassHunter]: '#abd473',
	[Class.ClassMage]: '#69ccf0',
	[Class.ClassPaladin]: '#f58cba',
	[Class.ClassPriest]: '#fff',
	[Class.ClassRogue]: '#fff569',
	[Class.ClassShaman]: '#2459ff',
	[Class.ClassWarlock]: '#9482c9',
	[Class.ClassWarrior]: '#c79c6e',
}

export const talentTreeIcons: Record<Class, Array<string>> = {
	[Class.ClassUnknown]: [],
	[Class.ClassDruid]: [
		'spell_nature_starfall.jpg',
		'ability_racial_bearform.jpg',
		'spell_nature_healingtouch.jpg',
		'ability_druid_catform.jpg',
	],
	[Class.ClassHunter]: [
		'ability_hunter_beasttaming.jpg',
		'ability_marksmanship.jpg',
		'ability_hunter_swiftstrike.jpg',
		// Pet specializations
		'ability_druid_swipe.jpg',
		'ability_hunter_pet_bear.jpg',
		'ability_hunter_combatexperience.jpg',
	],
	[Class.ClassMage]: [
		'spell_holy_magicalsentry.jpg',
		'spell_fire_firebolt02.jpg',
		'spell_frost_frostbolt02.jpg',
	],
	[Class.ClassPaladin]: [
		'spell_holy_holybolt.jpg',
		'spell_holy_devotionaura.jpg',
		'spell_holy_auraoflight.jpg',
	],
	[Class.ClassPriest]: [
		'spell_holy_powerwordshield.jpg',
		'spell_holy_guardianspirit.jpg',
		'spell_shadow_shadowwordpain.jpg',
		'spell_holy_holysmite.jpg',
	],
	[Class.ClassRogue]: [
		'ability_rogue_eviscerate.jpg',
		'ability_backstab.jpg',
		'ability_stealth.jpg',
	],
	[Class.ClassShaman]: [
		'spell_nature_lightning.jpg',
		'ability_shaman_stormstrike.jpg',
		'spell_nature_magicimmunity.jpg',
	],
	[Class.ClassWarlock]: [
		'spell_shadow_deathcoil.jpg',
		'spell_shadow_metamorphosis.jpg',
		'spell_shadow_rainoffire.jpg',
	],
	[Class.ClassWarrior]: [
		'ability_warrior_savageblow.jpg',
		'ability_warrior_innerrage.jpg',
		'inv_shield_06.jpg',
	],
};

export const titleIcons: Record<Class | Spec, string> = {
	[Spec.SpecBalanceDruid]: '/sod/assets/img/balance_druid_icon.png',
	[Spec.SpecFeralDruid]: '/sod/assets/img/feral_druid_icon.png',
	[Spec.SpecFeralTankDruid]: '/sod/assets/img/feral_druid_tank_icon.png',
	[Spec.SpecRestorationDruid]: '/sod/assets/img/resto_druid_icon.png',
	[Spec.SpecElementalShaman]: '/sod/assets/img/elemental_shaman_icon.png',
	[Spec.SpecEnhancementShaman]: '/sod/assets/img/enhancement_shaman_icon.png',
	[Spec.SpecRestorationShaman]: '/sod/assets/img/resto_shaman_icon.png',
	[Spec.SpecHunter]: '/sod/assets/img/hunter_icon.png',
	[Spec.SpecMage]: '/sod/assets/img/mage_icon.png',
	[Spec.SpecRogue]: '/sod/assets/img/rogue_icon.png',
	[Spec.SpecHolyPaladin]: '/sod/assets/img/holy_paladin_icon.png',
	[Spec.SpecProtectionPaladin]: '/sod/assets/img/protection_paladin_icon.png',
	[Spec.SpecRetributionPaladin]: '/sod/assets/img/retribution_icon.png',
	[Spec.SpecHealingPriest]: '/sod/assets/img/priest_icon.png',
	[Spec.SpecShadowPriest]: '/sod/assets/img/shadow_priest_icon.png',
	[Spec.SpecWarlock]: '/sod/assets/img/warlock_icon.png',
	[Spec.SpecTankWarlock]: '/sod/assets/img/tank_warlock_icon.jpg',
	[Spec.SpecWarrior]: '/sod/assets/img/warrior_icon.png',
	[Spec.SpecProtectionWarrior]: '/sod/assets/img/protection_warrior_icon.png',
};

export const raidSimIcon: string = '/sod/assets/img/raid_icon.png';
export const raidSimLabel: string = 'Full Raid Sim';

// Converts '1231321-12313123-0' to [40, 21, 0].
export function getTalentTreePoints(talentsString: string): Array<number> {
	const trees = talentsString.split('-');
	if (trees.length == 2)  {
		trees.push('0')
	}
	return trees.map(tree => sum([...tree].map(char => parseInt(char) || 0)));
}

export function getTalentPoints(talentsString: string): number {
	return sum(getTalentTreePoints(talentsString));
}

// Returns the index of the talent tree (0, 1, or 2) that has the most points.
export function getTalentTree(talentsString: string): number {
	const points = getTalentTreePoints(talentsString);
	return maxIndex(points) || 0;
}

enum IconSizes {
	Small = 'small',
	Medium = 'medium',
	Large = 'large',
}

// Returns the icon for a given spec
export function getSpecIcon(klass: Class, specNumber: number, size: IconSizes = IconSizes.Medium): string {
	const fileName = talentTreeIcons[klass][specNumber];

	return `https://wow.zamimg.com/images/wow/icons/${size}/${fileName}`;
}

// Returns the icon for a given spec based on talent point allocation.
export function getTalentTreeIcon(spec: Spec, talentsString: string, size: IconSizes = IconSizes.Medium): string {
	let specNumber = getTalentTree(talentsString);

	// Cat Druid is being considered a "4th spec"
	if (spec == Spec.SpecFeralDruid)
		specNumber += 2;

	const fileName = talentTreeIcons[specToClass[spec]][specNumber];

	return `https://wow.zamimg.com/images/wow/icons/${size}/${fileName}`;
}

// Gets the URL for the individual sim corresponding to the given spec.
const specSiteUrlTemplate = new URL(`${window.location.protocol}//${window.location.host}/${REPO_NAME}/SPEC/`);
export function getSpecSiteUrl(spec: Spec): string {
	let specString = Spec[spec]; // Returns 'SpecBalanceDruid' for BalanceDruid.
	specString = specString.substring('Spec'.length); // 'BalanceDruid'
	specString = camelToSnakeCase(specString); // 'balance_druid'
	return specSiteUrlTemplate.toString().replace('SPEC', specString);
}
export const raidSimSiteUrl = new URL(`${window.location.protocol}//${window.location.host}/${REPO_NAME}/raid/`).toString();

export function cssClassForClass(klass: Class): string {
	return classNames[klass].toLowerCase().replace(/\s/g, '-');
}

export function textCssClassForClass(klass: Class): string {
	return `text-${cssClassForClass(klass)}`;
}
export function textCssClassForSpec(spec: Spec): string {
	return textCssClassForClass(specToClass[spec]);
}

export type RotationUnion =
	BalanceDruidRotation |
	FeralDruidRotation |
	FeralTankDruidRotation |
	RestorationDruidRotation |
	HunterRotation |
	MageRotation |
	ElementalShamanRotation |
	EnhancementShamanRotation |
	RestorationShamanRotation |
	RogueRotation |
	HolyPaladinRotation |
	ProtectionPaladinRotation |
	RetributionPaladinRotation |
	HealingPriestRotation |
	ShadowPriestRotation |
	WarlockRotation |
	WarriorRotation |
	ProtectionWarriorRotation;
export type SpecRotation<T extends Spec> =
	T extends Spec.SpecBalanceDruid ? BalanceDruidRotation :
	T extends Spec.SpecFeralDruid ? FeralDruidRotation :
	T extends Spec.SpecFeralTankDruid ? FeralTankDruidRotation :
	T extends Spec.SpecRestorationDruid ? RestorationDruidRotation :
	T extends Spec.SpecElementalShaman ? ElementalShamanRotation :
	T extends Spec.SpecEnhancementShaman ? EnhancementShamanRotation :
	T extends Spec.SpecRestorationShaman ? RestorationShamanRotation :
	T extends Spec.SpecHunter ? HunterRotation :
	T extends Spec.SpecMage ? MageRotation :
	T extends Spec.SpecRogue ? RogueRotation :
	T extends Spec.SpecHolyPaladin ? HolyPaladinRotation :
	T extends Spec.SpecProtectionPaladin ? ProtectionPaladinRotation :
	T extends Spec.SpecRetributionPaladin ? RetributionPaladinRotation :
	T extends Spec.SpecHealingPriest ? HealingPriestRotation :
	T extends Spec.SpecShadowPriest ? ShadowPriestRotation :
	T extends Spec.SpecWarlock ? WarlockRotation :
	T extends Spec.SpecTankWarlock ? WarlockRotation :
	T extends Spec.SpecWarrior ? WarriorRotation :
	T extends Spec.SpecProtectionWarrior ? ProtectionWarriorRotation :
	ElementalShamanRotation; // Should never reach this case

export type TalentsUnion =
	DruidTalents |
	HunterTalents |
	MageTalents |
	RogueTalents |
	PaladinTalents |
	PriestTalents |
	ShamanTalents |
	WarlockTalents |
	WarriorTalents;
export type SpecTalents<T extends Spec> =
	T extends Spec.SpecBalanceDruid ? DruidTalents :
	T extends Spec.SpecFeralDruid ? DruidTalents :
	T extends Spec.SpecFeralTankDruid ? DruidTalents :
	T extends Spec.SpecRestorationDruid ? DruidTalents :
	T extends Spec.SpecElementalShaman ? ShamanTalents :
	T extends Spec.SpecEnhancementShaman ? ShamanTalents :
	T extends Spec.SpecRestorationShaman ? ShamanTalents :
	T extends Spec.SpecHunter ? HunterTalents :
	T extends Spec.SpecMage ? MageTalents :
	T extends Spec.SpecRogue ? RogueTalents :
	T extends Spec.SpecHolyPaladin ? PaladinTalents :
	T extends Spec.SpecProtectionPaladin ? PaladinTalents :
	T extends Spec.SpecRetributionPaladin ? PaladinTalents :
	T extends Spec.SpecHealingPriest ? PriestTalents :
	T extends Spec.SpecShadowPriest ? PriestTalents :
	T extends Spec.SpecWarlock ? WarlockTalents :
	T extends Spec.SpecTankWarlock ? WarlockTalents :
	T extends Spec.SpecWarrior ? WarriorTalents :
	T extends Spec.SpecProtectionWarrior ? WarriorTalents :
	ShamanTalents; // Should never reach this case

export type SpecOptionsUnion =
	BalanceDruidOptions |
	FeralDruidOptions |
	FeralTankDruidOptions |
	RestorationDruidOptions |
	ElementalShamanOptions |
	EnhancementShamanOptions |
	RestorationShamanOptions |
	HunterOptions |
	MageOptions |
	RogueOptions |
	HolyPaladinOptions |
	ProtectionPaladinOptions |
	RetributionPaladinOptions |
	HealingPriestOptions |
	ShadowPriestOptions |
	WarlockOptions |
	WarriorOptions |
	ProtectionWarriorOptions;
export type SpecOptions<T extends Spec> =
	T extends Spec.SpecBalanceDruid ? BalanceDruidOptions :
	T extends Spec.SpecFeralDruid ? FeralDruidOptions :
	T extends Spec.SpecFeralTankDruid ? FeralTankDruidOptions :
	T extends Spec.SpecRestorationDruid ? RestorationDruidOptions :
	T extends Spec.SpecElementalShaman ? ElementalShamanOptions :
	T extends Spec.SpecEnhancementShaman ? EnhancementShamanOptions :
	T extends Spec.SpecRestorationShaman ? RestorationShamanOptions :
	T extends Spec.SpecHunter ? HunterOptions :
	T extends Spec.SpecMage ? MageOptions :
	T extends Spec.SpecRogue ? RogueOptions :
	T extends Spec.SpecHolyPaladin ? HolyPaladinOptions :
	T extends Spec.SpecProtectionPaladin ? ProtectionPaladinOptions :
	T extends Spec.SpecRetributionPaladin ? RetributionPaladinOptions :
	T extends Spec.SpecHealingPriest ? HealingPriestOptions :
	T extends Spec.SpecShadowPriest ? ShadowPriestOptions :
	T extends Spec.SpecWarlock ? WarlockOptions :
	T extends Spec.SpecTankWarlock ? WarlockOptions :
	T extends Spec.SpecWarrior ? WarriorOptions :
	T extends Spec.SpecProtectionWarrior ? ProtectionWarriorOptions :
	ElementalShamanOptions; // Should never reach this case

export type SpecProtoUnion =
	BalanceDruid |
	FeralDruid |
	FeralTankDruid |
	RestorationDruid |
	ElementalShaman |
	EnhancementShaman |
	RestorationShaman |
	Hunter |
	Mage |
	Rogue |
	HolyPaladin |
	ProtectionPaladin |
	RetributionPaladin |
	HealingPriest |
	ShadowPriest |
	Warlock |
	TankWarlock |
	Warrior |
	ProtectionWarrior;
export type SpecProto<T extends Spec> =
	T extends Spec.SpecBalanceDruid ? BalanceDruid :
	T extends Spec.SpecFeralDruid ? FeralDruid :
	T extends Spec.SpecFeralTankDruid ? FeralTankDruid :
	T extends Spec.SpecRestorationDruid ? RestorationDruid :
	T extends Spec.SpecElementalShaman ? ElementalShaman :
	T extends Spec.SpecEnhancementShaman ? EnhancementShaman :
	T extends Spec.SpecRestorationShaman ? RestorationShaman :
	T extends Spec.SpecHunter ? Hunter :
	T extends Spec.SpecMage ? Mage :
	T extends Spec.SpecRogue ? Rogue :
	T extends Spec.SpecHolyPaladin ? HolyPaladin :
	T extends Spec.SpecProtectionPaladin ? ProtectionPaladin :
	T extends Spec.SpecRetributionPaladin ? RetributionPaladin :
	T extends Spec.SpecHealingPriest ? HealingPriest :
	T extends Spec.SpecShadowPriest ? ShadowPriest :
	T extends Spec.SpecWarlock ? Warlock :
	T extends Spec.SpecTankWarlock ? Warlock :
	T extends Spec.SpecWarrior ? Warrior :
	T extends Spec.SpecProtectionWarrior ? ProtectionWarrior :
	ElementalShaman; // Should never reach this case

export type SpecTypeFunctions<SpecType extends Spec> = {
	rotationCreate: () => SpecRotation<SpecType>;
	rotationEquals: (a: SpecRotation<SpecType>, b: SpecRotation<SpecType>) => boolean;
	rotationCopy: (a: SpecRotation<SpecType>) => SpecRotation<SpecType>;
	rotationToJson: (a: SpecRotation<SpecType>) => any;
	rotationFromJson: (obj: any) => SpecRotation<SpecType>;
	rotationFromPlayer: (player: PlayerProto) => SpecRotation<SpecType>;

	talentsCreate: () => SpecTalents<SpecType>;
	talentsEquals: (a: SpecTalents<SpecType>, b: SpecTalents<SpecType>) => boolean;
	talentsCopy: (a: SpecTalents<SpecType>) => SpecTalents<SpecType>;
	talentsToJson: (a: SpecTalents<SpecType>) => any;
	talentsFromJson: (obj: any) => SpecTalents<SpecType>;

	optionsCreate: () => SpecOptions<SpecType>;
	optionsEquals: (a: SpecOptions<SpecType>, b: SpecOptions<SpecType>) => boolean;
	optionsCopy: (a: SpecOptions<SpecType>) => SpecOptions<SpecType>;
	optionsToJson: (a: SpecOptions<SpecType>) => any;
	optionsFromJson: (obj: any) => SpecOptions<SpecType>;
	optionsFromPlayer: (player: PlayerProto) => SpecOptions<SpecType>;
};

export const specTypeFunctions: Record<Spec, SpecTypeFunctions<any>> = {
	[Spec.SpecBalanceDruid]: {
		rotationCreate: () => BalanceDruidRotation.create(),
		rotationEquals: (a, b) => BalanceDruidRotation.equals(a as BalanceDruidRotation, b as BalanceDruidRotation),
		rotationCopy: (a) => BalanceDruidRotation.clone(a as BalanceDruidRotation),
		rotationToJson: (a) => BalanceDruidRotation.toJson(a as BalanceDruidRotation),
		rotationFromJson: (obj) => BalanceDruidRotation.fromJson(obj),
		rotationFromPlayer: (player) => player.spec.oneofKind == 'balanceDruid'
			? player.spec.balanceDruid.rotation || BalanceDruidRotation.create()
			: BalanceDruidRotation.create(),

		talentsCreate: () => DruidTalents.create(),
		talentsEquals: (a, b) => DruidTalents.equals(a as DruidTalents, b as DruidTalents),
		talentsCopy: (a) => DruidTalents.clone(a as DruidTalents),
		talentsToJson: (a) => DruidTalents.toJson(a as DruidTalents),
		talentsFromJson: (obj) => DruidTalents.fromJson(obj),

		optionsCreate: () => BalanceDruidOptions.create(),
		optionsEquals: (a, b) => BalanceDruidOptions.equals(a as BalanceDruidOptions, b as BalanceDruidOptions),
		optionsCopy: (a) => BalanceDruidOptions.clone(a as BalanceDruidOptions),
		optionsToJson: (a) => BalanceDruidOptions.toJson(a as BalanceDruidOptions),
		optionsFromJson: (obj) => BalanceDruidOptions.fromJson(obj),
		optionsFromPlayer: (player) => player.spec.oneofKind == 'balanceDruid'
			? player.spec.balanceDruid.options || BalanceDruidOptions.create()
			: BalanceDruidOptions.create(),
	},
	[Spec.SpecFeralDruid]: {
		rotationCreate: () => FeralDruidRotation.create(),
		rotationEquals: (a, b) => FeralDruidRotation.equals(a as FeralDruidRotation, b as FeralDruidRotation),
		rotationCopy: (a) => FeralDruidRotation.clone(a as FeralDruidRotation),
		rotationToJson: (a) => FeralDruidRotation.toJson(a as FeralDruidRotation),
		rotationFromJson: (obj) => FeralDruidRotation.fromJson(obj),
		rotationFromPlayer: (player) => player.spec.oneofKind == 'feralDruid'
			? player.spec.feralDruid.rotation || FeralDruidRotation.create()
			: FeralDruidRotation.create(),

		talentsCreate: () => DruidTalents.create(),
		talentsEquals: (a, b) => DruidTalents.equals(a as DruidTalents, b as DruidTalents),
		talentsCopy: (a) => DruidTalents.clone(a as DruidTalents),
		talentsToJson: (a) => DruidTalents.toJson(a as DruidTalents),
		talentsFromJson: (obj) => DruidTalents.fromJson(obj),

		optionsCreate: () => FeralDruidOptions.create(),
		optionsEquals: (a, b) => FeralDruidOptions.equals(a as FeralDruidOptions, b as FeralDruidOptions),
		optionsCopy: (a) => FeralDruidOptions.clone(a as FeralDruidOptions),
		optionsToJson: (a) => FeralDruidOptions.toJson(a as FeralDruidOptions),
		optionsFromJson: (obj) => FeralDruidOptions.fromJson(obj),
		optionsFromPlayer: (player) => player.spec.oneofKind == 'feralDruid'
			? player.spec.feralDruid.options || FeralDruidOptions.create()
			: FeralDruidOptions.create(),
	},
	[Spec.SpecFeralTankDruid]: {
		rotationCreate: () => FeralTankDruidRotation.create(),
		rotationEquals: (a, b) => FeralTankDruidRotation.equals(a as FeralTankDruidRotation, b as FeralTankDruidRotation),
		rotationCopy: (a) => FeralTankDruidRotation.clone(a as FeralTankDruidRotation),
		rotationToJson: (a) => FeralTankDruidRotation.toJson(a as FeralTankDruidRotation),
		rotationFromJson: (obj) => FeralTankDruidRotation.fromJson(obj),
		rotationFromPlayer: (player) => player.spec.oneofKind == 'feralTankDruid'
			? player.spec.feralTankDruid.rotation || FeralTankDruidRotation.create()
			: FeralTankDruidRotation.create(),

		talentsCreate: () => DruidTalents.create(),
		talentsEquals: (a, b) => DruidTalents.equals(a as DruidTalents, b as DruidTalents),
		talentsCopy: (a) => DruidTalents.clone(a as DruidTalents),
		talentsToJson: (a) => DruidTalents.toJson(a as DruidTalents),
		talentsFromJson: (obj) => DruidTalents.fromJson(obj),

		optionsCreate: () => FeralTankDruidOptions.create(),
		optionsEquals: (a, b) => FeralTankDruidOptions.equals(a as FeralTankDruidOptions, b as FeralTankDruidOptions),
		optionsCopy: (a) => FeralTankDruidOptions.clone(a as FeralTankDruidOptions),
		optionsToJson: (a) => FeralTankDruidOptions.toJson(a as FeralTankDruidOptions),
		optionsFromJson: (obj) => FeralTankDruidOptions.fromJson(obj),
		optionsFromPlayer: (player) => player.spec.oneofKind == 'feralTankDruid'
			? player.spec.feralTankDruid.options || FeralTankDruidOptions.create()
			: FeralTankDruidOptions.create(),
	},
	[Spec.SpecRestorationDruid]: {
		rotationCreate: () => RestorationDruidRotation.create(),
		rotationEquals: (a, b) => RestorationDruidRotation.equals(a as RestorationDruidRotation, b as RestorationDruidRotation),
		rotationCopy: (a) => RestorationDruidRotation.clone(a as RestorationDruidRotation),
		rotationToJson: (a) => RestorationDruidRotation.toJson(a as RestorationDruidRotation),
		rotationFromJson: (obj) => RestorationDruidRotation.fromJson(obj),
		rotationFromPlayer: (player) => player.spec.oneofKind == 'restorationDruid'
			? player.spec.restorationDruid.rotation || RestorationDruidRotation.create()
			: RestorationDruidRotation.create(),

		talentsCreate: () => DruidTalents.create(),
		talentsEquals: (a, b) => DruidTalents.equals(a as DruidTalents, b as DruidTalents),
		talentsCopy: (a) => DruidTalents.clone(a as DruidTalents),
		talentsToJson: (a) => DruidTalents.toJson(a as DruidTalents),
		talentsFromJson: (obj) => DruidTalents.fromJson(obj),

		optionsCreate: () => RestorationDruidOptions.create(),
		optionsEquals: (a, b) => RestorationDruidOptions.equals(a as RestorationDruidOptions, b as RestorationDruidOptions),
		optionsCopy: (a) => RestorationDruidOptions.clone(a as RestorationDruidOptions),
		optionsToJson: (a) => RestorationDruidOptions.toJson(a as RestorationDruidOptions),
		optionsFromJson: (obj) => RestorationDruidOptions.fromJson(obj),
		optionsFromPlayer: (player) => player.spec.oneofKind == 'restorationDruid'
			? player.spec.restorationDruid.options || RestorationDruidOptions.create()
			: RestorationDruidOptions.create(),
	},
	[Spec.SpecElementalShaman]: {
		rotationCreate: () => ElementalShamanRotation.create(),
		rotationEquals: (a, b) => ElementalShamanRotation.equals(a as ElementalShamanRotation, b as ElementalShamanRotation),
		rotationCopy: (a) => ElementalShamanRotation.clone(a as ElementalShamanRotation),
		rotationToJson: (a) => ElementalShamanRotation.toJson(a as ElementalShamanRotation),
		rotationFromJson: (obj) => ElementalShamanRotation.fromJson(obj),
		rotationFromPlayer: (player) => player.spec.oneofKind == 'elementalShaman'
			? player.spec.elementalShaman.rotation || ElementalShamanRotation.create()
			: ElementalShamanRotation.create(),

		talentsCreate: () => ShamanTalents.create(),
		talentsEquals: (a, b) => ShamanTalents.equals(a as ShamanTalents, b as ShamanTalents),
		talentsCopy: (a) => ShamanTalents.clone(a as ShamanTalents),
		talentsToJson: (a) => ShamanTalents.toJson(a as ShamanTalents),
		talentsFromJson: (obj) => ShamanTalents.fromJson(obj),

		optionsCreate: () => ElementalShamanOptions.create(),
		optionsEquals: (a, b) => ElementalShamanOptions.equals(a as ElementalShamanOptions, b as ElementalShamanOptions),
		optionsCopy: (a) => ElementalShamanOptions.clone(a as ElementalShamanOptions),
		optionsToJson: (a) => ElementalShamanOptions.toJson(a as ElementalShamanOptions),
		optionsFromJson: (obj) => ElementalShamanOptions.fromJson(obj),
		optionsFromPlayer: (player) => player.spec.oneofKind == 'elementalShaman'
			? player.spec.elementalShaman.options || ElementalShamanOptions.create()
			: ElementalShamanOptions.create(),
	},
	[Spec.SpecEnhancementShaman]: {
		rotationCreate: () => EnhancementShamanRotation.create(),
		rotationEquals: (a, b) => EnhancementShamanRotation.equals(a as EnhancementShamanRotation, b as EnhancementShamanRotation),
		rotationCopy: (a) => EnhancementShamanRotation.clone(a as EnhancementShamanRotation),
		rotationToJson: (a) => EnhancementShamanRotation.toJson(a as EnhancementShamanRotation),
		rotationFromJson: (obj) => EnhancementShamanRotation.fromJson(obj),
		rotationFromPlayer: (player) => player.spec.oneofKind == 'enhancementShaman'
			? player.spec.enhancementShaman.rotation || EnhancementShamanRotation.create()
			: EnhancementShamanRotation.create(),

		talentsCreate: () => ShamanTalents.create(),
		talentsEquals: (a, b) => ShamanTalents.equals(a as ShamanTalents, b as ShamanTalents),
		talentsCopy: (a) => ShamanTalents.clone(a as ShamanTalents),
		talentsToJson: (a) => ShamanTalents.toJson(a as ShamanTalents),
		talentsFromJson: (obj) => ShamanTalents.fromJson(obj),

		optionsCreate: () => EnhancementShamanOptions.create(),
		optionsEquals: (a, b) => EnhancementShamanOptions.equals(a as EnhancementShamanOptions, b as EnhancementShamanOptions),
		optionsCopy: (a) => EnhancementShamanOptions.clone(a as EnhancementShamanOptions),
		optionsToJson: (a) => EnhancementShamanOptions.toJson(a as EnhancementShamanOptions),
		optionsFromJson: (obj) => EnhancementShamanOptions.fromJson(obj),
		optionsFromPlayer: (player) => player.spec.oneofKind == 'enhancementShaman'
			? player.spec.enhancementShaman.options || EnhancementShamanOptions.create()
			: EnhancementShamanOptions.create(),
	},
	[Spec.SpecRestorationShaman]: {
		rotationCreate: () => RestorationShamanRotation.create(),
		rotationEquals: (a, b) => RestorationShamanRotation.equals(a as RestorationShamanRotation, b as RestorationShamanRotation),
		rotationCopy: (a) => RestorationShamanRotation.clone(a as RestorationShamanRotation),
		rotationToJson: (a) => RestorationShamanRotation.toJson(a as RestorationShamanRotation),
		rotationFromJson: (obj) => RestorationShamanRotation.fromJson(obj),
		rotationFromPlayer: (player) => player.spec.oneofKind == 'restorationShaman'
			? player.spec.restorationShaman.rotation || RestorationShamanRotation.create()
			: RestorationShamanRotation.create(),

		talentsCreate: () => ShamanTalents.create(),
		talentsEquals: (a, b) => ShamanTalents.equals(a as ShamanTalents, b as ShamanTalents),
		talentsCopy: (a) => ShamanTalents.clone(a as ShamanTalents),
		talentsToJson: (a) => ShamanTalents.toJson(a as ShamanTalents),
		talentsFromJson: (obj) => ShamanTalents.fromJson(obj),

		optionsCreate: () => RestorationShamanOptions.create(),
		optionsEquals: (a, b) => RestorationShamanOptions.equals(a as RestorationShamanOptions, b as RestorationShamanOptions),
		optionsCopy: (a) => RestorationShamanOptions.clone(a as RestorationShamanOptions),
		optionsToJson: (a) => RestorationShamanOptions.toJson(a as RestorationShamanOptions),
		optionsFromJson: (obj) => RestorationShamanOptions.fromJson(obj),
		optionsFromPlayer: (player) => player.spec.oneofKind == 'restorationShaman'
			? player.spec.restorationShaman.options || RestorationShamanOptions.create()
			: RestorationShamanOptions.create(),
	},
	[Spec.SpecHunter]: {
		rotationCreate: () => HunterRotation.create(),
		rotationEquals: (a, b) => HunterRotation.equals(a as HunterRotation, b as HunterRotation),
		rotationCopy: (a) => HunterRotation.clone(a as HunterRotation),
		rotationToJson: (a) => HunterRotation.toJson(a as HunterRotation),
		rotationFromJson: (obj) => HunterRotation.fromJson(obj),
		rotationFromPlayer: (player) => player.spec.oneofKind == 'hunter'
			? player.spec.hunter.rotation || HunterRotation.create()
			: HunterRotation.create(),

		talentsCreate: () => HunterTalents.create(),
		talentsEquals: (a, b) => HunterTalents.equals(a as HunterTalents, b as HunterTalents),
		talentsCopy: (a) => HunterTalents.clone(a as HunterTalents),
		talentsToJson: (a) => HunterTalents.toJson(a as HunterTalents),
		talentsFromJson: (obj) => HunterTalents.fromJson(obj),

		optionsCreate: () => HunterOptions.create(),
		optionsEquals: (a, b) => HunterOptions.equals(a as HunterOptions, b as HunterOptions),
		optionsCopy: (a) => HunterOptions.clone(a as HunterOptions),
		optionsToJson: (a) => HunterOptions.toJson(a as HunterOptions),
		optionsFromJson: (obj) => HunterOptions.fromJson(obj),
		optionsFromPlayer: (player) => player.spec.oneofKind == 'hunter'
			? player.spec.hunter.options || HunterOptions.create()
			: HunterOptions.create(),
	},
	[Spec.SpecMage]: {
		rotationCreate: () => MageRotation.create(),
		rotationEquals: (a, b) => MageRotation.equals(a as MageRotation, b as MageRotation),
		rotationCopy: (a) => MageRotation.clone(a as MageRotation),
		rotationToJson: (a) => MageRotation.toJson(a as MageRotation),
		rotationFromJson: (obj) => MageRotation.fromJson(obj),
		rotationFromPlayer: (_player) => MageRotation.create(),

		talentsCreate: () => MageTalents.create(),
		talentsEquals: (a, b) => MageTalents.equals(a as MageTalents, b as MageTalents),
		talentsCopy: (a) => MageTalents.clone(a as MageTalents),
		talentsToJson: (a) => MageTalents.toJson(a as MageTalents),
		talentsFromJson: (obj) => MageTalents.fromJson(obj),

		optionsCreate: () => MageOptions.create(),
		optionsEquals: (a, b) => MageOptions.equals(a as MageOptions, b as MageOptions),
		optionsCopy: (a) => MageOptions.clone(a as MageOptions),
		optionsToJson: (a) => MageOptions.toJson(a as MageOptions),
		optionsFromJson: (obj) => MageOptions.fromJson(obj),
		optionsFromPlayer: (player) => player.spec.oneofKind == 'mage'
			? player.spec.mage.options || MageOptions.create()
			: MageOptions.create(),
	},
	[Spec.SpecHolyPaladin]: {
		rotationCreate: () => HolyPaladinRotation.create(),
		rotationEquals: (a, b) => HolyPaladinRotation.equals(a as HolyPaladinRotation, b as HolyPaladinRotation),
		rotationCopy: (a) => HolyPaladinRotation.clone(a as HolyPaladinRotation),
		rotationToJson: (a) => HolyPaladinRotation.toJson(a as HolyPaladinRotation),
		rotationFromJson: (obj) => HolyPaladinRotation.fromJson(obj),
		rotationFromPlayer: (player) => player.spec.oneofKind == 'holyPaladin'
			? player.spec.holyPaladin.rotation || HolyPaladinRotation.create()
			: HolyPaladinRotation.create(),

		talentsCreate: () => PaladinTalents.create(),
		talentsEquals: (a, b) => PaladinTalents.equals(a as PaladinTalents, b as PaladinTalents),
		talentsCopy: (a) => PaladinTalents.clone(a as PaladinTalents),
		talentsToJson: (a) => PaladinTalents.toJson(a as PaladinTalents),
		talentsFromJson: (obj) => PaladinTalents.fromJson(obj),

		optionsCreate: () => HolyPaladinOptions.create(),
		optionsEquals: (a, b) => HolyPaladinOptions.equals(a as HolyPaladinOptions, b as HolyPaladinOptions),
		optionsCopy: (a) => HolyPaladinOptions.clone(a as HolyPaladinOptions),
		optionsToJson: (a) => HolyPaladinOptions.toJson(a as HolyPaladinOptions),
		optionsFromJson: (obj) => HolyPaladinOptions.fromJson(obj),
		optionsFromPlayer: (player) => player.spec.oneofKind == 'holyPaladin'
			? player.spec.holyPaladin.options || HolyPaladinOptions.create()
			: HolyPaladinOptions.create(),
	},
	[Spec.SpecProtectionPaladin]: {
		rotationCreate: () => ProtectionPaladinRotation.create(),
		rotationEquals: (a, b) => ProtectionPaladinRotation.equals(a as ProtectionPaladinRotation, b as ProtectionPaladinRotation),
		rotationCopy: (a) => ProtectionPaladinRotation.clone(a as ProtectionPaladinRotation),
		rotationToJson: (a) => ProtectionPaladinRotation.toJson(a as ProtectionPaladinRotation),
		rotationFromJson: (obj) => ProtectionPaladinRotation.fromJson(obj),
		rotationFromPlayer: (player) => player.spec.oneofKind == 'protectionPaladin'
			? player.spec.protectionPaladin.rotation || ProtectionPaladinRotation.create()
			: ProtectionPaladinRotation.create(),

		talentsCreate: () => PaladinTalents.create(),
		talentsEquals: (a, b) => PaladinTalents.equals(a as PaladinTalents, b as PaladinTalents),
		talentsCopy: (a) => PaladinTalents.clone(a as PaladinTalents),
		talentsToJson: (a) => PaladinTalents.toJson(a as PaladinTalents),
		talentsFromJson: (obj) => PaladinTalents.fromJson(obj),

		optionsCreate: () => ProtectionPaladinOptions.create(),
		optionsEquals: (a, b) => ProtectionPaladinOptions.equals(a as ProtectionPaladinOptions, b as ProtectionPaladinOptions),
		optionsCopy: (a) => ProtectionPaladinOptions.clone(a as ProtectionPaladinOptions),
		optionsToJson: (a) => ProtectionPaladinOptions.toJson(a as ProtectionPaladinOptions),
		optionsFromJson: (obj) => ProtectionPaladinOptions.fromJson(obj),
		optionsFromPlayer: (player) => player.spec.oneofKind == 'protectionPaladin'
			? player.spec.protectionPaladin.options || ProtectionPaladinOptions.create()
			: ProtectionPaladinOptions.create(),
	},
	[Spec.SpecRetributionPaladin]: {
		rotationCreate: () => RetributionPaladinRotation.create(),
		rotationEquals: (a, b) => RetributionPaladinRotation.equals(a as RetributionPaladinRotation, b as RetributionPaladinRotation),
		rotationCopy: (a) => RetributionPaladinRotation.clone(a as RetributionPaladinRotation),
		rotationToJson: (a) => RetributionPaladinRotation.toJson(a as RetributionPaladinRotation),
		rotationFromJson: (obj) => RetributionPaladinRotation.fromJson(obj),
		rotationFromPlayer: (player) => player.spec.oneofKind == 'retributionPaladin'
			? player.spec.retributionPaladin.rotation || RetributionPaladinRotation.create()
			: RetributionPaladinRotation.create(),

		talentsCreate: () => PaladinTalents.create(),
		talentsEquals: (a, b) => PaladinTalents.equals(a as PaladinTalents, b as PaladinTalents),
		talentsCopy: (a) => PaladinTalents.clone(a as PaladinTalents),
		talentsToJson: (a) => PaladinTalents.toJson(a as PaladinTalents),
		talentsFromJson: (obj) => PaladinTalents.fromJson(obj),

		optionsCreate: () => RetributionPaladinOptions.create(),
		optionsEquals: (a, b) => RetributionPaladinOptions.equals(a as RetributionPaladinOptions, b as RetributionPaladinOptions),
		optionsCopy: (a) => RetributionPaladinOptions.clone(a as RetributionPaladinOptions),
		optionsToJson: (a) => RetributionPaladinOptions.toJson(a as RetributionPaladinOptions),
		optionsFromJson: (obj) => RetributionPaladinOptions.fromJson(obj),
		optionsFromPlayer: (player) => player.spec.oneofKind == 'retributionPaladin'
			? player.spec.retributionPaladin.options || RetributionPaladinOptions.create()
			: RetributionPaladinOptions.create(),
	},
	[Spec.SpecRogue]: {
		rotationCreate: () => RogueRotation.create(),
		rotationEquals: (a, b) => RogueRotation.equals(a as RogueRotation, b as RogueRotation),
		rotationCopy: (a) => RogueRotation.clone(a as RogueRotation),
		rotationToJson: (a) => RogueRotation.toJson(a as RogueRotation),
		rotationFromJson: (obj) => RogueRotation.fromJson(obj),
		rotationFromPlayer: (player) => player.spec.oneofKind == 'rogue'
			? player.spec.rogue.rotation || RogueRotation.create()
			: RogueRotation.create(),

		talentsCreate: () => RogueTalents.create(),
		talentsEquals: (a, b) => RogueTalents.equals(a as RogueTalents, b as RogueTalents),
		talentsCopy: (a) => RogueTalents.clone(a as RogueTalents),
		talentsToJson: (a) => RogueTalents.toJson(a as RogueTalents),
		talentsFromJson: (obj) => RogueTalents.fromJson(obj),

		optionsCreate: () => RogueOptions.create(),
		optionsEquals: (a, b) => RogueOptions.equals(a as RogueOptions, b as RogueOptions),
		optionsCopy: (a) => RogueOptions.clone(a as RogueOptions),
		optionsToJson: (a) => RogueOptions.toJson(a as RogueOptions),
		optionsFromJson: (obj) => RogueOptions.fromJson(obj),
		optionsFromPlayer: (player) => player.spec.oneofKind == 'rogue'
			? player.spec.rogue.options || RogueOptions.create()
			: RogueOptions.create(),
	},
	[Spec.SpecHealingPriest]: {
		rotationCreate: () => HealingPriestRotation.create(),
		rotationEquals: (a, b) => HealingPriestRotation.equals(a as HealingPriestRotation, b as HealingPriestRotation),
		rotationCopy: (a) => HealingPriestRotation.clone(a as HealingPriestRotation),
		rotationToJson: (a) => HealingPriestRotation.toJson(a as HealingPriestRotation),
		rotationFromJson: (obj) => HealingPriestRotation.fromJson(obj),
		rotationFromPlayer: (player) => player.spec.oneofKind == 'healingPriest'
			? player.spec.healingPriest.rotation || HealingPriestRotation.create()
			: HealingPriestRotation.create(),

		talentsCreate: () => PriestTalents.create(),
		talentsEquals: (a, b) => PriestTalents.equals(a as PriestTalents, b as PriestTalents),
		talentsCopy: (a) => PriestTalents.clone(a as PriestTalents),
		talentsToJson: (a) => PriestTalents.toJson(a as PriestTalents),
		talentsFromJson: (obj) => PriestTalents.fromJson(obj),

		optionsCreate: () => HealingPriestOptions.create(),
		optionsEquals: (a, b) => HealingPriestOptions.equals(a as HealingPriestOptions, b as HealingPriestOptions),
		optionsCopy: (a) => HealingPriestOptions.clone(a as HealingPriestOptions),
		optionsToJson: (a) => HealingPriestOptions.toJson(a as HealingPriestOptions),
		optionsFromJson: (obj) => HealingPriestOptions.fromJson(obj),
		optionsFromPlayer: (player) => player.spec.oneofKind == 'healingPriest'
			? player.spec.healingPriest.options || HealingPriestOptions.create()
			: HealingPriestOptions.create(),
	},
	[Spec.SpecShadowPriest]: {
		rotationCreate: () => ShadowPriestRotation.create(),
		rotationEquals: (a, b) => ShadowPriestRotation.equals(a as ShadowPriestRotation, b as ShadowPriestRotation),
		rotationCopy: (a) => ShadowPriestRotation.clone(a as ShadowPriestRotation),
		rotationToJson: (a) => ShadowPriestRotation.toJson(a as ShadowPriestRotation),
		rotationFromJson: (obj) => ShadowPriestRotation.fromJson(obj),
		rotationFromPlayer: (player) => player.spec.oneofKind == 'shadowPriest'
			? player.spec.shadowPriest.rotation || ShadowPriestRotation.create()
			: ShadowPriestRotation.create(),

		talentsCreate: () => PriestTalents.create(),
		talentsEquals: (a, b) => PriestTalents.equals(a as PriestTalents, b as PriestTalents),
		talentsCopy: (a) => PriestTalents.clone(a as PriestTalents),
		talentsToJson: (a) => PriestTalents.toJson(a as PriestTalents),
		talentsFromJson: (obj) => PriestTalents.fromJson(obj),

		optionsCreate: () => ShadowPriestOptions.create(),
		optionsEquals: (a, b) => ShadowPriestOptions.equals(a as ShadowPriestOptions, b as ShadowPriestOptions),
		optionsCopy: (a) => ShadowPriestOptions.clone(a as ShadowPriestOptions),
		optionsToJson: (a) => ShadowPriestOptions.toJson(a as ShadowPriestOptions),
		optionsFromJson: (obj) => ShadowPriestOptions.fromJson(obj),
		optionsFromPlayer: (player) => player.spec.oneofKind == 'shadowPriest'
			? player.spec.shadowPriest.options || ShadowPriestOptions.create()
			: ShadowPriestOptions.create(),
	},
	[Spec.SpecWarlock]: {
		rotationCreate: () => WarlockRotation.create(),
		rotationEquals: (a, b) => WarlockRotation.equals(a as WarlockRotation, b as WarlockRotation),
		rotationCopy: (a) => WarlockRotation.clone(a as WarlockRotation),
		rotationToJson: (a) => WarlockRotation.toJson(a as WarlockRotation),
		rotationFromJson: (obj) => WarlockRotation.fromJson(obj),
		rotationFromPlayer: (player) => player.spec.oneofKind == 'warlock'
			? player.spec.warlock.rotation || WarlockRotation.create()
			: WarlockRotation.create(),

		talentsCreate: () => WarlockTalents.create(),
		talentsEquals: (a, b) => WarlockTalents.equals(a as WarlockTalents, b as WarlockTalents),
		talentsCopy: (a) => WarlockTalents.clone(a as WarlockTalents),
		talentsToJson: (a) => WarlockTalents.toJson(a as WarlockTalents),
		talentsFromJson: (obj) => WarlockTalents.fromJson(obj),

		optionsCreate: () => WarlockOptions.create(),
		optionsEquals: (a, b) => WarlockOptions.equals(a as WarlockOptions, b as WarlockOptions),
		optionsCopy: (a) => WarlockOptions.clone(a as WarlockOptions),
		optionsToJson: (a) => WarlockOptions.toJson(a as WarlockOptions),
		optionsFromJson: (obj) => WarlockOptions.fromJson(obj),
		optionsFromPlayer: (player) => player.spec.oneofKind == 'warlock'
			? player.spec.warlock.options || WarlockOptions.create()
			: WarlockOptions.create(),
	},
	[Spec.SpecTankWarlock]: {
		rotationCreate: () => WarlockRotation.create(),
		rotationEquals: (a, b) => WarlockRotation.equals(a as WarlockRotation, b as WarlockRotation),
		rotationCopy: (a) => WarlockRotation.clone(a as WarlockRotation),
		rotationToJson: (a) => WarlockRotation.toJson(a as WarlockRotation),
		rotationFromJson: (obj) => WarlockRotation.fromJson(obj),
		rotationFromPlayer: (player) => player.spec.oneofKind == 'tankWarlock'
			? player.spec.tankWarlock.rotation || WarlockRotation.create()
			: WarlockRotation.create(),

		talentsCreate: () => WarlockTalents.create(),
		talentsEquals: (a, b) => WarlockTalents.equals(a as WarlockTalents, b as WarlockTalents),
		talentsCopy: (a) => WarlockTalents.clone(a as WarlockTalents),
		talentsToJson: (a) => WarlockTalents.toJson(a as WarlockTalents),
		talentsFromJson: (obj) => WarlockTalents.fromJson(obj),

		optionsCreate: () => WarlockOptions.create(),
		optionsEquals: (a, b) => WarlockOptions.equals(a as WarlockOptions, b as WarlockOptions),
		optionsCopy: (a) => WarlockOptions.clone(a as WarlockOptions),
		optionsToJson: (a) => WarlockOptions.toJson(a as WarlockOptions),
		optionsFromJson: (obj) => WarlockOptions.fromJson(obj),
		optionsFromPlayer: (player) => player.spec.oneofKind == 'tankWarlock'
			? player.spec.tankWarlock.options || WarlockOptions.create()
			: WarlockOptions.create(),
	},
	[Spec.SpecWarrior]: {
		rotationCreate: () => WarriorRotation.create(),
		rotationEquals: (a, b) => WarriorRotation.equals(a as WarriorRotation, b as WarriorRotation),
		rotationCopy: (a) => WarriorRotation.clone(a as WarriorRotation),
		rotationToJson: (a) => WarriorRotation.toJson(a as WarriorRotation),
		rotationFromJson: (obj) => WarriorRotation.fromJson(obj),
		rotationFromPlayer: (player) => player.spec.oneofKind == 'warrior'
			? player.spec.warrior.rotation || WarriorRotation.create()
			: WarriorRotation.create(),

		talentsCreate: () => WarriorTalents.create(),
		talentsEquals: (a, b) => WarriorTalents.equals(a as WarriorTalents, b as WarriorTalents),
		talentsCopy: (a) => WarriorTalents.clone(a as WarriorTalents),
		talentsToJson: (a) => WarriorTalents.toJson(a as WarriorTalents),
		talentsFromJson: (obj) => WarriorTalents.fromJson(obj),

		optionsCreate: () => WarriorOptions.create(),
		optionsEquals: (a, b) => WarriorOptions.equals(a as WarriorOptions, b as WarriorOptions),
		optionsCopy: (a) => WarriorOptions.clone(a as WarriorOptions),
		optionsToJson: (a) => WarriorOptions.toJson(a as WarriorOptions),
		optionsFromJson: (obj) => WarriorOptions.fromJson(obj),
		optionsFromPlayer: (player) => player.spec.oneofKind == 'warrior'
			? player.spec.warrior.options || WarriorOptions.create()
			: WarriorOptions.create(),
	},
	[Spec.SpecProtectionWarrior]: {
		rotationCreate: () => ProtectionWarriorRotation.create(),
		rotationEquals: (a, b) => ProtectionWarriorRotation.equals(a as ProtectionWarriorRotation, b as ProtectionWarriorRotation),
		rotationCopy: (a) => ProtectionWarriorRotation.clone(a as ProtectionWarriorRotation),
		rotationToJson: (a) => ProtectionWarriorRotation.toJson(a as ProtectionWarriorRotation),
		rotationFromJson: (obj) => ProtectionWarriorRotation.fromJson(obj),
		rotationFromPlayer: (player) => player.spec.oneofKind == 'protectionWarrior'
			? player.spec.protectionWarrior.rotation || ProtectionWarriorRotation.create()
			: ProtectionWarriorRotation.create(),

		talentsCreate: () => WarriorTalents.create(),
		talentsEquals: (a, b) => WarriorTalents.equals(a as WarriorTalents, b as WarriorTalents),
		talentsCopy: (a) => WarriorTalents.clone(a as WarriorTalents),
		talentsToJson: (a) => WarriorTalents.toJson(a as WarriorTalents),
		talentsFromJson: (obj) => WarriorTalents.fromJson(obj),

		optionsCreate: () => ProtectionWarriorOptions.create(),
		optionsEquals: (a, b) => ProtectionWarriorOptions.equals(a as ProtectionWarriorOptions, b as ProtectionWarriorOptions),
		optionsCopy: (a) => ProtectionWarriorOptions.clone(a as ProtectionWarriorOptions),
		optionsToJson: (a) => ProtectionWarriorOptions.toJson(a as ProtectionWarriorOptions),
		optionsFromJson: (obj) => ProtectionWarriorOptions.fromJson(obj),
		optionsFromPlayer: (player) => player.spec.oneofKind == 'protectionWarrior'
			? player.spec.protectionWarrior.options || ProtectionWarriorOptions.create()
			: ProtectionWarriorOptions.create(),
	},
};

export const raceToFaction: Record<Race, Faction> = {
	[Race.RaceUnknown]: Faction.Unknown,
	[Race.RaceBloodElf]: Faction.Horde,
	[Race.RaceDraenei]: Faction.Alliance,
	[Race.RaceDwarf]: Faction.Alliance,
	[Race.RaceGnome]: Faction.Alliance,
	[Race.RaceHuman]: Faction.Alliance,
	[Race.RaceNightElf]: Faction.Alliance,
	[Race.RaceOrc]: Faction.Horde,
	[Race.RaceTauren]: Faction.Horde,
	[Race.RaceTroll]: Faction.Horde,
	[Race.RaceUndead]: Faction.Horde,
};

export const specToClass: Record<Spec, Class> = {
	[Spec.SpecBalanceDruid]: Class.ClassDruid,
	[Spec.SpecFeralDruid]: Class.ClassDruid,
	[Spec.SpecFeralTankDruid]: Class.ClassDruid,
	[Spec.SpecRestorationDruid]: Class.ClassDruid,
	[Spec.SpecHunter]: Class.ClassHunter,
	[Spec.SpecMage]: Class.ClassMage,
	[Spec.SpecRogue]: Class.ClassRogue,
	[Spec.SpecHolyPaladin]: Class.ClassPaladin,
	[Spec.SpecProtectionPaladin]: Class.ClassPaladin,
	[Spec.SpecRetributionPaladin]: Class.ClassPaladin,
	[Spec.SpecHealingPriest]: Class.ClassPriest,
	[Spec.SpecShadowPriest]: Class.ClassPriest,
	[Spec.SpecElementalShaman]: Class.ClassShaman,
	[Spec.SpecEnhancementShaman]: Class.ClassShaman,
	[Spec.SpecRestorationShaman]: Class.ClassShaman,
	[Spec.SpecWarlock]: Class.ClassWarlock,
	[Spec.SpecTankWarlock]: Class.ClassWarlock,
	[Spec.SpecWarrior]: Class.ClassWarrior,
	[Spec.SpecProtectionWarrior]: Class.ClassWarrior,
};

const druidRaces = [
	Race.RaceTauren,
	Race.RaceNightElf,
];
const hunterRaces = [
	Race.RaceBloodElf,
	Race.RaceDraenei,
	Race.RaceDwarf,
	Race.RaceNightElf,
	Race.RaceOrc,
	Race.RaceTauren,
	Race.RaceTroll,
];
const mageRaces = [
	Race.RaceTroll,
	Race.RaceBloodElf,
	Race.RaceDraenei,
	Race.RaceGnome,
	Race.RaceHuman,
	Race.RaceUndead,
];
const paladinRaces = [
	Race.RaceBloodElf,
	Race.RaceDraenei,
	Race.RaceDwarf,
	Race.RaceHuman,
];
const priestRaces = [
	Race.RaceTroll,
	Race.RaceBloodElf,
	Race.RaceDraenei,
	Race.RaceDwarf,
	Race.RaceHuman,
	Race.RaceNightElf,
	Race.RaceUndead,
];
const rogueRaces = [
	Race.RaceBloodElf,
	Race.RaceDwarf,
	Race.RaceGnome,
	Race.RaceHuman,
	Race.RaceNightElf,
	Race.RaceOrc,
	Race.RaceTroll,
	Race.RaceUndead,
];
const shamanRaces = [
	Race.RaceOrc,
	Race.RaceDraenei,
	Race.RaceTauren,
	Race.RaceTroll,
];
const warlockRaces = [
	Race.RaceBloodElf,
	Race.RaceGnome,
	Race.RaceHuman,
	Race.RaceOrc,
	Race.RaceUndead,
];
const warriorRaces = [
	Race.RaceDraenei,
	Race.RaceDwarf,
	Race.RaceGnome,
	Race.RaceHuman,
	Race.RaceNightElf,
	Race.RaceOrc,
	Race.RaceTauren,
	Race.RaceTroll,
	Race.RaceUndead,
];

export const specToEligibleRaces: Record<Spec, Array<Race>> = {
	[Spec.SpecBalanceDruid]: druidRaces,
	[Spec.SpecFeralDruid]: druidRaces,
	[Spec.SpecFeralTankDruid]: druidRaces,
	[Spec.SpecRestorationDruid]: druidRaces,
	[Spec.SpecElementalShaman]: shamanRaces,
	[Spec.SpecEnhancementShaman]: shamanRaces,
	[Spec.SpecRestorationShaman]: shamanRaces,
	[Spec.SpecHunter]: hunterRaces,
	[Spec.SpecMage]: mageRaces,
	[Spec.SpecHolyPaladin]: paladinRaces,
	[Spec.SpecProtectionPaladin]: paladinRaces,
	[Spec.SpecRetributionPaladin]: paladinRaces,
	[Spec.SpecRogue]: rogueRaces,
	[Spec.SpecHealingPriest]: priestRaces,
	[Spec.SpecShadowPriest]: priestRaces,
	[Spec.SpecWarlock]: warlockRaces,
	[Spec.SpecTankWarlock]: warlockRaces,
	[Spec.SpecWarrior]: warriorRaces,
	[Spec.SpecProtectionWarrior]: warriorRaces,
};

// Specs that can dual wield. This could be based on class, except that
// Enhancement Shaman learn dual wield from a talent.
const dualWieldSpecs: Array<Spec> = [
	Spec.SpecEnhancementShaman,
	Spec.SpecHunter,
	Spec.SpecRogue,
	Spec.SpecWarrior,
	Spec.SpecProtectionWarrior,
];
export function isDualWieldSpec(spec: Spec): boolean {
	return dualWieldSpecs.includes(spec);
}

const tankSpecs: Array<Spec> = [
	Spec.SpecFeralTankDruid,
	Spec.SpecProtectionPaladin,
	Spec.SpecProtectionWarrior,
	Spec.SpecTankWarlock,
];
export function isTankSpec(spec: Spec): boolean {
	return tankSpecs.includes(spec);
}

const healingSpecs: Array<Spec> = [
	Spec.SpecRestorationDruid,
	Spec.SpecHolyPaladin,
	Spec.SpecHealingPriest,
	Spec.SpecRestorationShaman,
];
export function isHealingSpec(spec: Spec): boolean {
	return healingSpecs.includes(spec);
}

const rangedDpsSpecs: Array<Spec> = [
	Spec.SpecBalanceDruid,
	Spec.SpecHunter,
	Spec.SpecMage,
	Spec.SpecShadowPriest,
	Spec.SpecElementalShaman,
	Spec.SpecWarlock,
];
export function isRangedDpsSpec(spec: Spec): boolean {
	return rangedDpsSpecs.includes(spec);
}
export function isMeleeDpsSpec(spec: Spec): boolean {
	return !isTankSpec(spec) && !isHealingSpec(spec) && !isRangedDpsSpec(spec);
}

// Prefixes used for storing browser data for each site. Even if a Spec is
// renamed, DO NOT change these values or people will lose their saved data.
export const specToLocalStorageKey: Record<Spec, string> = {
	[Spec.SpecBalanceDruid]: '__wotlk_balance_druid',
	[Spec.SpecFeralDruid]: '__wotlk_feral_druid',
	[Spec.SpecFeralTankDruid]: '__wotlk_feral_tank_druid',
	[Spec.SpecRestorationDruid]: '__wotlk_restoration_druid',
	[Spec.SpecElementalShaman]: '__wotlk_elemental_shaman',
	[Spec.SpecEnhancementShaman]: '__wotlk_enhacement_shaman',
	[Spec.SpecRestorationShaman]: '__wotlk_restoration_shaman',
	[Spec.SpecHunter]: '__wotlk_hunter',
	[Spec.SpecMage]: '__wotlk_mage',
	[Spec.SpecHolyPaladin]: '__wotlk_holy_paladin',
	[Spec.SpecProtectionPaladin]: '__wotlk_protection_paladin',
	[Spec.SpecRetributionPaladin]: '__wotlk_retribution_paladin',
	[Spec.SpecRogue]: '__wotlk_rogue',
	[Spec.SpecHealingPriest]: '__wotlk_healing_priest',
	[Spec.SpecShadowPriest]: '__wotlk_shadow_priest',
	[Spec.SpecWarlock]: '__wotlk_warlock',
	[Spec.SpecTankWarlock]: '__wotlk_tank_warlock',
	[Spec.SpecWarrior]: '__wotlk_warrior',
	[Spec.SpecProtectionWarrior]: '__wotlk_protection_warrior',
};

// Returns a copy of playerOptions, with the class field set.
export function withSpecProto<SpecType extends Spec>(
	spec: Spec,
	player: PlayerProto,
	rotation: SpecRotation<SpecType>,
	specOptions: SpecOptions<SpecType>): PlayerProto {
	const copy = PlayerProto.clone(player);

	switch (spec) {
		case Spec.SpecBalanceDruid:
			copy.spec = {
				oneofKind: 'balanceDruid',
				balanceDruid: BalanceDruid.create({
					rotation: rotation as BalanceDruidRotation,
					options: specOptions as BalanceDruidOptions,
				}),
			};
			return copy;
		case Spec.SpecFeralDruid:
			copy.spec = {
				oneofKind: 'feralDruid',
				feralDruid: FeralDruid.create({
					rotation: rotation as FeralDruidRotation,
					options: specOptions as FeralDruidOptions,
				}),
			};
			return copy;
		case Spec.SpecFeralTankDruid:
			copy.spec = {
				oneofKind: 'feralTankDruid',
				feralTankDruid: FeralTankDruid.create({
					rotation: rotation as FeralTankDruidRotation,
					options: specOptions as FeralTankDruidOptions,
				}),
			};
			return copy;
		case Spec.SpecRestorationDruid:
			copy.spec = {
				oneofKind: 'restorationDruid',
				restorationDruid: RestorationDruid.create({
					rotation: rotation as RestorationDruidRotation,
					options: specOptions as RestorationDruidOptions,
				}),
			};
			return copy;
		case Spec.SpecElementalShaman:
			copy.spec = {
				oneofKind: 'elementalShaman',
				elementalShaman: ElementalShaman.create({
					rotation: rotation as ElementalShamanRotation,
					options: specOptions as ElementalShamanOptions,
				}),
			};
			return copy;
		case Spec.SpecEnhancementShaman:
			copy.spec = {
				oneofKind: 'enhancementShaman',
				enhancementShaman: EnhancementShaman.create({
					rotation: rotation as EnhancementShamanRotation,
					options: specOptions as ElementalShamanOptions,
				}),
			};
			return copy;
		case Spec.SpecRestorationShaman:
			copy.spec = {
				oneofKind: 'restorationShaman',
				restorationShaman: RestorationShaman.create({
					rotation: rotation as RestorationShamanRotation,
					options: specOptions as RestorationShamanOptions,
				}),
			};
			return copy;
		case Spec.SpecHunter:
			copy.spec = {
				oneofKind: 'hunter',
				hunter: Hunter.create({
					rotation: rotation as HunterRotation,
					options: specOptions as HunterOptions,
				}),
			};
			return copy;
		case Spec.SpecMage:
			copy.spec = {
				oneofKind: 'mage',
				mage: Mage.create({
					options: specOptions as MageOptions,
				}),
			};
			return copy;
		case Spec.SpecHolyPaladin:
			copy.spec = {
				oneofKind: 'holyPaladin',
				holyPaladin: HolyPaladin.create({
					rotation: rotation as HolyPaladinRotation,
					options: specOptions as HolyPaladinOptions,
				}),
			};
			return copy;
		case Spec.SpecProtectionPaladin:
			copy.spec = {
				oneofKind: 'protectionPaladin',
				protectionPaladin: ProtectionPaladin.create({
					rotation: rotation as ProtectionPaladinRotation,
					options: specOptions as ProtectionPaladinOptions,
				}),
			};
			return copy;
		case Spec.SpecRetributionPaladin:
			copy.spec = {
				oneofKind: 'retributionPaladin',
				retributionPaladin: RetributionPaladin.create({
					rotation: rotation as RetributionPaladinRotation,
					options: specOptions as RetributionPaladinOptions,
				}),
			};
			return copy;
		case Spec.SpecRogue:
			copy.spec = {
				oneofKind: 'rogue',
				rogue: Rogue.create({
					rotation: rotation as RogueRotation,
					options: specOptions as RogueOptions,
				}),
			};
			return copy;
		case Spec.SpecHealingPriest:
			copy.spec = {
				oneofKind: 'healingPriest',
				healingPriest: HealingPriest.create({
					rotation: rotation as HealingPriestRotation,
					options: specOptions as HealingPriestOptions,
				}),
			};
			return copy;
		case Spec.SpecShadowPriest:
			copy.spec = {
				oneofKind: 'shadowPriest',
				shadowPriest: ShadowPriest.create({
					rotation: rotation as ShadowPriestRotation,
					options: specOptions as ShadowPriestOptions,
				}),
			};
			return copy;
		case Spec.SpecWarlock:
			copy.spec = {
				oneofKind: 'warlock',
				warlock: Warlock.create({
					rotation: rotation as WarlockRotation,
					options: specOptions as WarlockOptions,
				}),
			};
			return copy;
		case Spec.SpecTankWarlock:
			copy.spec = {
				oneofKind: 'tankWarlock',
				tankWarlock: TankWarlock.create({
					rotation: rotation as WarlockRotation,
					options: specOptions as WarlockOptions,
				}),
			};
			return copy;
		case Spec.SpecWarrior:
			copy.spec = {
				oneofKind: 'warrior',
				warrior: Warrior.create({
					rotation: rotation as WarriorRotation,
					options: specOptions as WarriorOptions,
				}),
			};
			return copy;
		case Spec.SpecProtectionWarrior:
			copy.spec = {
				oneofKind: 'protectionWarrior',
				protectionWarrior: ProtectionWarrior.create({
					rotation: rotation as ProtectionWarriorRotation,
					options: specOptions as ProtectionWarriorOptions,
				}),
			};
			return copy;
	}
}

export function playerToSpec(player: PlayerProto): Spec {
	const specValues = getEnumValues(Spec);
	for (let i = 0; i < specValues.length; i++) {
		const spec = specValues[i] as Spec;
		let specString = Spec[spec]; // Returns 'SpecBalanceDruid' for BalanceDruid.
		specString = specString.substring('Spec'.length); // 'BalanceDruid'
		specString = specString.charAt(0).toLowerCase() + specString.slice(1); // 'balanceDruid'

		if (player.spec.oneofKind == specString) {
			return spec;
		}
	}

	throw new Error('Unable to parse spec from player proto: ' + JSON.stringify(PlayerProto.toJson(player), null, 2));
}

export const classToMaxArmorType: Record<Class, ArmorType> = {
	[Class.ClassUnknown]: ArmorType.ArmorTypeUnknown,
	[Class.ClassDruid]: ArmorType.ArmorTypeLeather,
	[Class.ClassHunter]: ArmorType.ArmorTypeMail,
	[Class.ClassMage]: ArmorType.ArmorTypeCloth,
	[Class.ClassPaladin]: ArmorType.ArmorTypePlate,
	[Class.ClassPriest]: ArmorType.ArmorTypeCloth,
	[Class.ClassRogue]: ArmorType.ArmorTypeLeather,
	[Class.ClassShaman]: ArmorType.ArmorTypeMail,
	[Class.ClassWarlock]: ArmorType.ArmorTypeCloth,
	[Class.ClassWarrior]: ArmorType.ArmorTypePlate,
};

export const classToEligibleRangedWeaponTypes: Record<Class, Array<RangedWeaponType>> = {
	[Class.ClassUnknown]: [],
	[Class.ClassDruid]: [RangedWeaponType.RangedWeaponTypeIdol],
	[Class.ClassHunter]: [
		RangedWeaponType.RangedWeaponTypeBow,
		RangedWeaponType.RangedWeaponTypeCrossbow,
		RangedWeaponType.RangedWeaponTypeGun,
	],
	[Class.ClassMage]: [RangedWeaponType.RangedWeaponTypeWand],
	[Class.ClassPaladin]: [RangedWeaponType.RangedWeaponTypeLibram],
	[Class.ClassPriest]: [RangedWeaponType.RangedWeaponTypeWand],
	[Class.ClassRogue]: [
		RangedWeaponType.RangedWeaponTypeBow,
		RangedWeaponType.RangedWeaponTypeCrossbow,
		RangedWeaponType.RangedWeaponTypeGun,
		RangedWeaponType.RangedWeaponTypeThrown,
	],
	[Class.ClassShaman]: [RangedWeaponType.RangedWeaponTypeTotem],
	[Class.ClassWarlock]: [RangedWeaponType.RangedWeaponTypeWand],
	[Class.ClassWarrior]: [
		RangedWeaponType.RangedWeaponTypeBow,
		RangedWeaponType.RangedWeaponTypeCrossbow,
		RangedWeaponType.RangedWeaponTypeGun,
		RangedWeaponType.RangedWeaponTypeThrown,
	],
};

interface EligibleWeaponType {
	weaponType: WeaponType,
	canUseTwoHand?: boolean,
}

export const classToEligibleWeaponTypes: Record<Class, Array<EligibleWeaponType>> = {
	[Class.ClassUnknown]: [],
	[Class.ClassDruid]: [
		{ weaponType: WeaponType.WeaponTypeDagger },
		{ weaponType: WeaponType.WeaponTypeFist },
		{ weaponType: WeaponType.WeaponTypeMace, canUseTwoHand: true },
		{ weaponType: WeaponType.WeaponTypeOffHand },
		{ weaponType: WeaponType.WeaponTypeStaff, canUseTwoHand: true },
		{ weaponType: WeaponType.WeaponTypePolearm, canUseTwoHand: true },
	],
	[Class.ClassHunter]: [
		{ weaponType: WeaponType.WeaponTypeAxe, canUseTwoHand: true },
		{ weaponType: WeaponType.WeaponTypeDagger },
		{ weaponType: WeaponType.WeaponTypeFist },
		{ weaponType: WeaponType.WeaponTypeOffHand },
		{ weaponType: WeaponType.WeaponTypePolearm, canUseTwoHand: true },
		{ weaponType: WeaponType.WeaponTypeSword, canUseTwoHand: true },
		{ weaponType: WeaponType.WeaponTypeStaff, canUseTwoHand: true },
	],
	[Class.ClassMage]: [
		{ weaponType: WeaponType.WeaponTypeDagger },
		{ weaponType: WeaponType.WeaponTypeOffHand },
		{ weaponType: WeaponType.WeaponTypeStaff, canUseTwoHand: true },
		{ weaponType: WeaponType.WeaponTypeSword },
	],
	[Class.ClassPaladin]: [
		{ weaponType: WeaponType.WeaponTypeAxe, canUseTwoHand: true },
		{ weaponType: WeaponType.WeaponTypeMace, canUseTwoHand: true },
		{ weaponType: WeaponType.WeaponTypeOffHand },
		{ weaponType: WeaponType.WeaponTypePolearm, canUseTwoHand: true },
		{ weaponType: WeaponType.WeaponTypeShield },
		{ weaponType: WeaponType.WeaponTypeSword, canUseTwoHand: true },
	],
	[Class.ClassPriest]: [
		{ weaponType: WeaponType.WeaponTypeDagger },
		{ weaponType: WeaponType.WeaponTypeMace },
		{ weaponType: WeaponType.WeaponTypeOffHand },
		{ weaponType: WeaponType.WeaponTypeStaff, canUseTwoHand: true },
	],
	[Class.ClassRogue]: [
		{ weaponType: WeaponType.WeaponTypeAxe, canUseTwoHand: false },
		{ weaponType: WeaponType.WeaponTypeDagger },
		{ weaponType: WeaponType.WeaponTypeFist },
		{ weaponType: WeaponType.WeaponTypeMace },
		{ weaponType: WeaponType.WeaponTypeOffHand },
		{ weaponType: WeaponType.WeaponTypeSword },
	],
	[Class.ClassShaman]: [
		{ weaponType: WeaponType.WeaponTypeAxe, canUseTwoHand: true },
		{ weaponType: WeaponType.WeaponTypeDagger },
		{ weaponType: WeaponType.WeaponTypeFist },
		{ weaponType: WeaponType.WeaponTypeMace, canUseTwoHand: true },
		{ weaponType: WeaponType.WeaponTypeOffHand },
		{ weaponType: WeaponType.WeaponTypeShield },
		{ weaponType: WeaponType.WeaponTypeStaff, canUseTwoHand: true },
	],
	[Class.ClassWarlock]: [
		{ weaponType: WeaponType.WeaponTypeDagger },
		{ weaponType: WeaponType.WeaponTypeOffHand },
		{ weaponType: WeaponType.WeaponTypeStaff, canUseTwoHand: true },
		{ weaponType: WeaponType.WeaponTypeSword },
	],
	[Class.ClassWarrior]: [
		{ weaponType: WeaponType.WeaponTypeAxe, canUseTwoHand: true },
		{ weaponType: WeaponType.WeaponTypeDagger },
		{ weaponType: WeaponType.WeaponTypeFist },
		{ weaponType: WeaponType.WeaponTypeMace, canUseTwoHand: true },
		{ weaponType: WeaponType.WeaponTypeOffHand },
		{ weaponType: WeaponType.WeaponTypePolearm, canUseTwoHand: true },
		{ weaponType: WeaponType.WeaponTypeShield },
		{ weaponType: WeaponType.WeaponTypeStaff, canUseTwoHand: true },
		{ weaponType: WeaponType.WeaponTypeSword, canUseTwoHand: true },
	],
};

export function isSharpWeaponType(weaponType: WeaponType): boolean {
	return [
		WeaponType.WeaponTypeAxe,
		WeaponType.WeaponTypeDagger,
		WeaponType.WeaponTypePolearm,
		WeaponType.WeaponTypeSword,
	].includes(weaponType);
}

export function isBluntWeaponType(weaponType: WeaponType): boolean {
	return [
		WeaponType.WeaponTypeFist,
		WeaponType.WeaponTypeMace,
		WeaponType.WeaponTypeStaff,
	].includes(weaponType);
}

// Returns true if this item may be equipped in at least 1 slot for the given Spec.
export function canEquipItem<SpecType extends Spec>(player: Player<SpecType>, item: Item, slot: ItemSlot | undefined): boolean {
	const spec = player.spec;
	const playerClass = specToClass[spec];
	if (item.classAllowlist.length > 0 && !item.classAllowlist.includes(playerClass)) {
		return false;
	}

	if (item.requiresLevel > player.getLevel()){
		return false
	}

	if ([ItemType.ItemTypeFinger, ItemType.ItemTypeTrinket].includes(item.type)) {
		return true;
	}

	if (item.type == ItemType.ItemTypeWeapon) {
		const eligibleWeaponType = classToEligibleWeaponTypes[playerClass].find(wt => wt.weaponType == item.weaponType);
		if (!eligibleWeaponType) {
			return false;
		}

		if ((item.handType == HandType.HandTypeOffHand || (item.handType == HandType.HandTypeOneHand && slot == ItemSlot.ItemSlotOffHand))
			&& ![WeaponType.WeaponTypeShield, WeaponType.WeaponTypeOffHand].includes(item.weaponType)
			&& !dualWieldSpecs.includes(spec)) {
			return false;
		}

		if (item.handType == HandType.HandTypeTwoHand && !eligibleWeaponType.canUseTwoHand) {
			return false;
		}
		if (item.handType == HandType.HandTypeTwoHand && slot == ItemSlot.ItemSlotOffHand && spec != Spec.SpecWarrior) {
			return false;
		}

		return true;
	}

	if (item.type == ItemType.ItemTypeRanged) {
		return classToEligibleRangedWeaponTypes[playerClass].includes(item.rangedWeaponType);
	}

	// At this point, we know the item is an armor piece (feet, chest, legs, etc).
	return classToMaxArmorType[playerClass] >= item.armorType;
}

export const itemTypeToSlotsMap: Partial<Record<ItemType, Array<ItemSlot>>> = {
	[ItemType.ItemTypeUnknown]: [],
	[ItemType.ItemTypeHead]: [ItemSlot.ItemSlotHead],
	[ItemType.ItemTypeNeck]: [ItemSlot.ItemSlotNeck],
	[ItemType.ItemTypeShoulder]: [ItemSlot.ItemSlotShoulder],
	[ItemType.ItemTypeBack]: [ItemSlot.ItemSlotBack],
	[ItemType.ItemTypeChest]: [ItemSlot.ItemSlotChest],
	[ItemType.ItemTypeWrist]: [ItemSlot.ItemSlotWrist],
	[ItemType.ItemTypeHands]: [ItemSlot.ItemSlotHands],
	[ItemType.ItemTypeWaist]: [ItemSlot.ItemSlotWaist],
	[ItemType.ItemTypeLegs]: [ItemSlot.ItemSlotLegs],
	[ItemType.ItemTypeFeet]: [ItemSlot.ItemSlotFeet],
	[ItemType.ItemTypeFinger]: [ItemSlot.ItemSlotFinger1, ItemSlot.ItemSlotFinger2],
	[ItemType.ItemTypeTrinket]: [ItemSlot.ItemSlotTrinket1, ItemSlot.ItemSlotTrinket2],
	[ItemType.ItemTypeRanged]: [ItemSlot.ItemSlotRanged],
};

export function getEligibleItemSlots(item: Item): Array<ItemSlot> {
	if (itemTypeToSlotsMap[item.type]) {
		return itemTypeToSlotsMap[item.type]!;
	}

	if (item.type == ItemType.ItemTypeWeapon) {
		if (item.handType == HandType.HandTypeMainHand) {
			return [ItemSlot.ItemSlotMainHand];
		} else if (item.handType == HandType.HandTypeOffHand) {
			return [ItemSlot.ItemSlotOffHand];
			// Missing HandTypeTwoHand 
			// We allow 2H weapons to be wielded in mainhand and offhand for Fury Warriors
		} else {
			return [ItemSlot.ItemSlotMainHand, ItemSlot.ItemSlotOffHand];
		}
	}

	// Should never reach here
	throw new Error('Could not find item slots for item: ' + Item.toJsonString(item));
};

// Returns whether the given main-hand and off-hand items can be worn at the
// same time.
export function validWeaponCombo(mainHand: Item | null | undefined, offHand: Item | null | undefined, canDW2h: boolean): boolean {
	if (mainHand == null || offHand == null) {
		return true;
	}

	if (mainHand.handType == HandType.HandTypeTwoHand && !canDW2h) {
		return false;
	} else if (mainHand.handType == HandType.HandTypeTwoHand &&
		(mainHand.weaponType == WeaponType.WeaponTypePolearm || mainHand.weaponType == WeaponType.WeaponTypeStaff)) {
		return false;
	}

	if (offHand.handType == HandType.HandTypeTwoHand && !canDW2h) {
		return false;
	} else if (offHand.handType == HandType.HandTypeTwoHand &&
		(offHand.weaponType == WeaponType.WeaponTypePolearm || offHand.weaponType == WeaponType.WeaponTypeStaff)) {
		return false;
	}

	return true;
}

// Returns all item slots to which the enchant might be applied.
// 
// Note that this alone is not enough; some items have further restrictions,
// e.g. some weapon enchants may only be applied to 2H weapons.
export function getEligibleEnchantSlots(enchant: Enchant): Array<ItemSlot> {
	return [enchant.type].concat(enchant.extraTypes || []).map(type => {
		if (itemTypeToSlotsMap[type]) {
			return itemTypeToSlotsMap[type]!;
		}

		if (type == ItemType.ItemTypeWeapon) {
			return [ItemSlot.ItemSlotMainHand, ItemSlot.ItemSlotOffHand];
		}

		// Should never reach here
		throw new Error('Could not find item slots for enchant: ' + Enchant.toJsonString(enchant));
	}).flat();
};

export function enchantAppliesToItem(enchant: Enchant, item: Item): boolean {
	const sharedSlots = intersection(getEligibleEnchantSlots(enchant), getEligibleItemSlots(item));
	if (sharedSlots.length == 0)
		return false;

	if (enchant.enchantType == EnchantType.EnchantTypeTwoHand && item.handType != HandType.HandTypeTwoHand)
		return false;

	if ((enchant.enchantType == EnchantType.EnchantTypeShield) != (item.weaponType == WeaponType.WeaponTypeShield))
		return false;

	if (enchant.enchantType == EnchantType.EnchantTypeStaff && item.weaponType != WeaponType.WeaponTypeStaff)
		return false;

	if (item.weaponType == WeaponType.WeaponTypeOffHand)
		return false;

	if (sharedSlots.includes(ItemSlot.ItemSlotRanged)) {
		if (![
			RangedWeaponType.RangedWeaponTypeBow,
			RangedWeaponType.RangedWeaponTypeCrossbow,
			RangedWeaponType.RangedWeaponTypeGun,
		].includes(item.rangedWeaponType))
			return false;
	}

	return true;
};

export function canEquipEnchant(enchant: Enchant, spec: Spec): boolean {
	const playerClass = specToClass[spec];
	if (enchant.classAllowlist.length > 0 && !enchant.classAllowlist.includes(playerClass)) {
		return false;
	}

	return true;
}

export function newUnitReference(raidIndex: number): UnitReference {
	return UnitReference.create({
		type: UnitReference_Type.Player,
		index: raidIndex,
	});
}

export function emptyUnitReference(): UnitReference {
	return UnitReference.create();
}

// Makes a new set of assignments with everything 0'd out.
export function makeBlankBlessingsAssignments(numPaladins: number): BlessingsAssignments {
	const assignments = BlessingsAssignments.create();
	for (let i = 0; i < numPaladins; i++) {
		assignments.paladins.push(BlessingsAssignment.create({
			blessings: new Array(NUM_SPECS).fill(Blessings.BlessingUnknown),
		}));
	}
	return assignments;
}

export function makeBlessingsAssignments(numPaladins: number, data: Array<{ spec: Spec, blessings: Array<Blessings> }>): BlessingsAssignments {
	const assignments = makeBlankBlessingsAssignments(numPaladins);
	for (let i = 0; i < data.length; i++) {
		const spec = data[i].spec;
		const blessings = data[i].blessings;
		for (let j = 0; j < blessings.length; j++) {
			if (j >= assignments.paladins.length) {
				// Can't assign more blessings since we ran out of paladins
				break
			}
			assignments.paladins[j].blessings[spec] = blessings[j];
		}
	}
	return assignments;
}

// Default blessings settings in the raid sim UI.
export function makeDefaultBlessings(numPaladins: number): BlessingsAssignments {
	return makeBlessingsAssignments(numPaladins, [
		{ spec: Spec.SpecBalanceDruid, blessings: [Blessings.BlessingOfKings, Blessings.BlessingOfWisdom] },
		{ spec: Spec.SpecFeralDruid, blessings: [Blessings.BlessingOfKings, Blessings.BlessingOfMight, Blessings.BlessingOfWisdom] },
		{ spec: Spec.SpecFeralTankDruid, blessings: [Blessings.BlessingOfKings, Blessings.BlessingOfMight, Blessings.BlessingOfSanctuary] },
		{ spec: Spec.SpecRestorationDruid, blessings: [Blessings.BlessingOfKings, Blessings.BlessingOfWisdom] },
		{ spec: Spec.SpecHunter, blessings: [Blessings.BlessingOfKings, Blessings.BlessingOfMight, Blessings.BlessingOfWisdom] },
		{ spec: Spec.SpecMage, blessings: [Blessings.BlessingOfKings, Blessings.BlessingOfWisdom] },
		{ spec: Spec.SpecHolyPaladin, blessings: [Blessings.BlessingOfKings, Blessings.BlessingOfWisdom] },
		{ spec: Spec.SpecProtectionPaladin, blessings: [Blessings.BlessingOfKings, Blessings.BlessingOfSanctuary, Blessings.BlessingOfWisdom, Blessings.BlessingOfMight] },
		{ spec: Spec.SpecRetributionPaladin, blessings: [Blessings.BlessingOfKings, Blessings.BlessingOfMight, Blessings.BlessingOfWisdom] },
		{ spec: Spec.SpecHealingPriest, blessings: [Blessings.BlessingOfKings, Blessings.BlessingOfWisdom] },
		{ spec: Spec.SpecShadowPriest, blessings: [Blessings.BlessingOfKings, Blessings.BlessingOfWisdom] },
		{ spec: Spec.SpecRogue, blessings: [Blessings.BlessingOfKings, Blessings.BlessingOfMight] },
		{ spec: Spec.SpecElementalShaman, blessings: [Blessings.BlessingOfKings, Blessings.BlessingOfWisdom] },
		{ spec: Spec.SpecEnhancementShaman, blessings: [Blessings.BlessingOfKings, Blessings.BlessingOfMight, Blessings.BlessingOfWisdom] },
		{ spec: Spec.SpecRestorationShaman, blessings: [Blessings.BlessingOfKings, Blessings.BlessingOfWisdom] },
		{ spec: Spec.SpecWarlock, blessings: [Blessings.BlessingOfWisdom, Blessings.BlessingOfKings] },
		{ spec: Spec.SpecTankWarlock, blessings: [Blessings.BlessingOfWisdom, Blessings.BlessingOfMight, Blessings.BlessingOfKings] },
		{ spec: Spec.SpecWarrior, blessings: [Blessings.BlessingOfKings, Blessings.BlessingOfMight] },
		{ spec: Spec.SpecProtectionWarrior, blessings: [Blessings.BlessingOfKings, Blessings.BlessingOfMight, Blessings.BlessingOfSanctuary] },
	]);
};

export const orderedResourceTypes: Array<ResourceType> = [
	ResourceType.ResourceTypeHealth,
	ResourceType.ResourceTypeMana,
	ResourceType.ResourceTypeEnergy,
	ResourceType.ResourceTypeRage,
	ResourceType.ResourceTypeComboPoints,
	ResourceType.ResourceTypeFocus,
];

export const AL_CATEGORY_HARD_MODE = 'Hard Mode';
export const AL_CATEGORY_TITAN_RUNE = 'Titan Rune';
