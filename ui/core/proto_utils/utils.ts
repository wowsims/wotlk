import { REPO_NAME } from '/tbc/core/constants/other.js'
import { camelToSnakeCase } from '/tbc/core/utils.js';
import { getEnumValues } from '/tbc/core/utils.js';
import { intersection } from '/tbc/core/utils.js';
import { maxIndex } from '/tbc/core/utils.js';
import { sum } from '/tbc/core/utils.js';

import { Player } from '/tbc/core/proto/api.js';
import { ResourceType } from '/tbc/core/proto/api.js';
import { ArmorType } from '/tbc/core/proto/common.js';
import { Class } from '/tbc/core/proto/common.js';
import { Enchant } from '/tbc/core/proto/common.js';
import { EnchantType } from '/tbc/core/proto/common.js';
import { Gem } from '/tbc/core/proto/common.js';
import { GemColor } from '/tbc/core/proto/common.js';
import { HandType } from '/tbc/core/proto/common.js';
import { ItemSlot } from '/tbc/core/proto/common.js';
import { ItemType } from '/tbc/core/proto/common.js';
import { Item } from '/tbc/core/proto/common.js';
import { Race } from '/tbc/core/proto/common.js';
import { Faction } from '/tbc/core/proto/common.js';
import { RaidTarget } from '/tbc/core/proto/common.js';
import { RangedWeaponType } from '/tbc/core/proto/common.js';
import { Spec } from '/tbc/core/proto/common.js';
import { Stat } from '/tbc/core/proto/common.js';
import { WeaponType } from '/tbc/core/proto/common.js';
import { Blessings } from '/tbc/core/proto/paladin.js';
import { BlessingsAssignment } from '/tbc/core/proto/ui.js';
import { BlessingsAssignments } from '/tbc/core/proto/ui.js';

import { Stats } from './stats.js';

import * as Gems from '/tbc/core/proto_utils/gems.js';

import {
	BalanceDruid,
	FeralDruid,
	FeralTankDruid,
	BalanceDruid_Rotation as BalanceDruidRotation,
	FeralDruid_Rotation as FeralDruidRotation,
	FeralTankDruid_Rotation as FeralTankDruidRotation,
	DruidTalents,
	BalanceDruid_Options as BalanceDruidOptions,
	FeralDruid_Options as FeralDruidOptions,
	FeralTankDruid_Options as FeralTankDruidOptions
} from '/tbc/core/proto/druid.js';
import { ElementalShaman, EnhancementShaman_Rotation as EnhancementShamanRotation, ElementalShaman_Rotation as ElementalShamanRotation, ShamanTalents, ElementalShaman_Options as ElementalShamanOptions, EnhancementShaman_Options as EnhancementShamanOptions, EnhancementShaman } from '/tbc/core/proto/shaman.js';
import { Hunter, Hunter_Rotation as HunterRotation, HunterTalents, Hunter_Options as HunterOptions } from '/tbc/core/proto/hunter.js';
import { Mage, Mage_Rotation as MageRotation, MageTalents, Mage_Options as MageOptions } from '/tbc/core/proto/mage.js';
import { Rogue, Rogue_Rotation as RogueRotation, RogueTalents, Rogue_Options as RogueOptions } from '/tbc/core/proto/rogue.js';
import { RetributionPaladin, RetributionPaladin_Rotation as RetributionPaladinRotation, PaladinTalents, RetributionPaladin_Options as RetributionPaladinOptions } from '/tbc/core/proto/paladin.js';
import { ProtectionPaladin, ProtectionPaladin_Rotation as ProtectionPaladinRotation, ProtectionPaladin_Options as ProtectionPaladinOptions } from '/tbc/core/proto/paladin.js';
import { ShadowPriest, SmitePriest_Rotation as SmitePriestRotation, ShadowPriest_Rotation as ShadowPriestRotation, PriestTalents, ShadowPriest_Options as ShadowPriestOptions, SmitePriest_Options as SmitePriestOptions, SmitePriest } from '/tbc/core/proto/priest.js';
import { Warlock, Warlock_Rotation as WarlockRotation, WarlockTalents, Warlock_Options as WarlockOptions } from '/tbc/core/proto/warlock.js';
import { Warrior, Warrior_Rotation as WarriorRotation, WarriorTalents, Warrior_Options as WarriorOptions } from '/tbc/core/proto/warrior.js';
import { ProtectionWarrior, ProtectionWarrior_Rotation as ProtectionWarriorRotation, ProtectionWarrior_Options as ProtectionWarriorOptions } from '/tbc/core/proto/warrior.js';

export type DruidSpecs = [Spec.SpecBalanceDruid, Spec.SpecFeralDruid, Spec.SpecFeralTankDruid];
export type HunterSpecs = Spec.SpecHunter;
export type MageSpecs = Spec.SpecMage;
export type RogueSpecs = Spec.SpecRogue;
export type PaladinSpecs = [Spec.SpecRetributionPaladin, Spec.SpecProtectionPaladin];
export type PriestSpecs = [Spec.SpecShadowPriest, Spec.SpecSmitePriest];
export type ShamanSpecs = [Spec.SpecElementalShaman, Spec.SpecEnhancementShaman];
export type WarlockSpecs = Spec.SpecWarlock;
export type WarriorSpecs = [Spec.SpecWarrior, Spec.SpecProtectionWarrior];

export const NUM_SPECS = getEnumValues(Spec).length;

// The order in which specs should be presented, when it matters.
// Currently this is only used for the order of the paladin blessings UI.
export const naturalSpecOrder: Array<Spec> = [
	Spec.SpecBalanceDruid,
	Spec.SpecFeralDruid,
	Spec.SpecFeralTankDruid,
	Spec.SpecHunter,
	Spec.SpecMage,
	Spec.SpecRetributionPaladin,
	Spec.SpecProtectionPaladin,
	Spec.SpecShadowPriest,
	Spec.SpecSmitePriest,
	Spec.SpecRogue,
	Spec.SpecElementalShaman,
	Spec.SpecEnhancementShaman,
	Spec.SpecWarlock,
	Spec.SpecWarrior,
	Spec.SpecProtectionWarrior,
];

export const specNames: Record<Spec, string> = {
	[Spec.SpecBalanceDruid]: 'Balance Druid',
	[Spec.SpecElementalShaman]: 'Elemental Shaman',
	[Spec.SpecEnhancementShaman]: 'Enhancement Shaman',
	[Spec.SpecFeralDruid]: 'Feral Druid',
	[Spec.SpecFeralTankDruid]: 'Feral Tank Druid',
	[Spec.SpecHunter]: 'Hunter',
	[Spec.SpecMage]: 'Mage',
	[Spec.SpecRogue]: 'Rogue',
	[Spec.SpecRetributionPaladin]: 'Retribution Paladin',
	[Spec.SpecProtectionPaladin]: 'Protection Paladin',
	[Spec.SpecShadowPriest]: 'Shadow Priest',
	[Spec.SpecWarlock]: 'Warlock',
	[Spec.SpecWarrior]: 'Warrior',
	[Spec.SpecProtectionWarrior]: 'Protection Warrior',
	[Spec.SpecSmitePriest]: 'Smite Priest',
};

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

export const specIconsLarge: Record<Spec, string> = {
	[Spec.SpecBalanceDruid]: 'https://wow.zamimg.com/images/wow/icons/large/spell_nature_starfall.jpg',
	[Spec.SpecElementalShaman]: 'https://wow.zamimg.com/images/wow/icons/large/spell_nature_lightning.jpg',
	[Spec.SpecEnhancementShaman]: 'https://wow.zamimg.com/images/wow/icons/large/ability_shaman_stormstrike.jpg',
	[Spec.SpecFeralDruid]: 'https://wow.zamimg.com/images/wow/icons/large/ability_druid_catform.jpg',
	[Spec.SpecFeralTankDruid]: 'https://wow.zamimg.com/images/wow/icons/large/ability_racial_bearform.jpg',
	[Spec.SpecHunter]: 'https://wow.zamimg.com/images/wow/icons/large/ability_marksmanship.jpg',
	[Spec.SpecMage]: 'https://wow.zamimg.com/images/wow/icons/large/spell_holy_magicalsentry.jpg',
	[Spec.SpecRogue]: 'https://wow.zamimg.com/images/wow/icons/large/classicon_rogue.jpg',
	[Spec.SpecRetributionPaladin]: 'https://wow.zamimg.com/images/wow/icons/large/spell_holy_auraoflight.jpg',
	[Spec.SpecProtectionPaladin]: 'https://wow.zamimg.com/images/wow/icons/large/spell_holy_devotionaura.jpg',
	[Spec.SpecShadowPriest]: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_shadowwordpain.jpg',
	[Spec.SpecWarlock]: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_metamorphosis.jpg',
	[Spec.SpecWarrior]: 'https://wow.zamimg.com/images/wow/icons/large/ability_warrior_innerrage.jpg',
	[Spec.SpecProtectionWarrior]: 'https://wow.zamimg.com/images/wow/icons/large/ability_warrior_defensivestance.jpg',
	[Spec.SpecSmitePriest]: 'https://wow.zamimg.com/images/wow/icons/large/spell_holy_holysmite.jpg',
};

export const talentTreeIcons: Record<Class, Array<string>> = {
	[Class.ClassUnknown]: [],
	[Class.ClassDruid]: [
		'https://wow.zamimg.com/images/wow/icons/medium/spell_nature_starfall.jpg',
		'https://wow.zamimg.com/images/wow/icons/medium/ability_racial_bearform.jpg',
		'https://wow.zamimg.com/images/wow/icons/medium/spell_nature_healingtouch.jpg',
	],
	[Class.ClassHunter]: [
		'https://wow.zamimg.com/images/wow/icons/medium/ability_hunter_beasttaming.jpg',
		'https://wow.zamimg.com/images/wow/icons/medium/ability_marksmanship.jpg',
		'https://wow.zamimg.com/images/wow/icons/medium/ability_hunter_swiftstrike.jpg',
	],
	[Class.ClassMage]: [
		'https://wow.zamimg.com/images/wow/icons/medium/spell_holy_magicalsentry.jpg',
		'https://wow.zamimg.com/images/wow/icons/medium/spell_fire_firebolt02.jpg',
		'https://wow.zamimg.com/images/wow/icons/medium/spell_frost_frostbolt02.jpg',
	],
	[Class.ClassPaladin]: [
		'https://wow.zamimg.com/images/wow/icons/medium/spell_holy_holybolt.jpg',
		'https://wow.zamimg.com/images/wow/icons/medium/spell_holy_devotionaura.jpg',
		'https://wow.zamimg.com/images/wow/icons/medium/spell_holy_auraoflight.jpg',
	],
	[Class.ClassPriest]: [
		'https://wow.zamimg.com/images/wow/icons/medium/spell_holy_powerinfusion.jpg',
		'https://wow.zamimg.com/images/wow/icons/medium/spell_holy_holybolt.jpg',
		'https://wow.zamimg.com/images/wow/icons/medium/spell_shadow_shadowwordpain.jpg',
	],
	[Class.ClassRogue]: [
		'https://wow.zamimg.com/images/wow/icons/medium/ability_rogue_eviscerate.jpg',
		'https://wow.zamimg.com/images/wow/icons/medium/ability_backstab.jpg',
		'https://wow.zamimg.com/images/wow/icons/medium/ability_stealth.jpg',
	],
	[Class.ClassShaman]: [
		'https://wow.zamimg.com/images/wow/icons/medium/spell_nature_lightning.jpg',
		'https://wow.zamimg.com/images/wow/icons/medium/ability_shaman_stormstrike.jpg',
		'https://wow.zamimg.com/images/wow/icons/medium/spell_nature_magicimmunity.jpg',
	],
	[Class.ClassWarlock]: [
		'https://wow.zamimg.com/images/wow/icons/medium/spell_shadow_deathcoil.jpg',
		'https://wow.zamimg.com/images/wow/icons/medium/spell_shadow_metamorphosis.jpg',
		'https://wow.zamimg.com/images/wow/icons/medium/spell_shadow_rainoffire.jpg',
	],
	[Class.ClassWarrior]: [
		'https://wow.zamimg.com/images/wow/icons/medium/ability_warrior_savageblow.jpg',
		'https://wow.zamimg.com/images/wow/icons/medium/ability_warrior_innerrage.jpg',
		'https://wow.zamimg.com/images/wow/icons/medium/inv_shield_06.jpg',
	],
};

export const titleIcons: Record<Spec, string> = {
	[Spec.SpecBalanceDruid]: '/tbc/assets/balance_druid_icon.png',
	[Spec.SpecElementalShaman]: '/tbc/assets/elemental_shaman_icon.png',
	[Spec.SpecEnhancementShaman]: '/tbc/assets/enhancement_shaman_icon.png',
	[Spec.SpecFeralDruid]: '/tbc/assets/feral_druid_icon.png',
	[Spec.SpecFeralTankDruid]: '/tbc/assets/feral_druid_tank_icon.png',
	[Spec.SpecHunter]: '/tbc/assets/hunter_icon.png',
	[Spec.SpecMage]: '/tbc/assets/mage_icon.png',
	[Spec.SpecRogue]: '/tbc/assets/rogue_icon.png',
	[Spec.SpecRetributionPaladin]: '/tbc/assets/retribution_icon.png',
	[Spec.SpecProtectionPaladin]: '/tbc/assets/protection_paladin_icon.png',
	[Spec.SpecShadowPriest]: '/tbc/assets/shadow_priest_icon.png',
	[Spec.SpecWarlock]: '/tbc/assets/warlock_icon.png',
	[Spec.SpecWarrior]: '/tbc/assets/warrior_icon.png',
	[Spec.SpecProtectionWarrior]: '/tbc/assets/protection_warrior_icon.png',
	[Spec.SpecSmitePriest]: '/tbc/assets/smite_priest_icon.png',
};

export const raidSimIcon: string = '/tbc/assets/raid_icon.png';

// Returns the index of the talent tree (0, 1, or 2) that has the most points.
export function getTalentTree(talentsString: string): number {
	const trees = talentsString.split('-');
	const points = trees.map(tree => sum([...tree].map(char => parseInt(char))));
	return maxIndex(points) || 0;
}

// Returns the index of the talent tree (0, 1, or 2) that has the most points.
export function getTalentTreeIcon(spec: Spec, talentsString: string): string {
	const talentTreeIdx = getTalentTree(talentsString);
	return talentTreeIcons[specToClass[spec]][talentTreeIdx];
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

export type RotationUnion =
	BalanceDruidRotation |
	FeralDruidRotation |
	FeralTankDruidRotation |
	HunterRotation |
	MageRotation |
	ElementalShamanRotation |
	EnhancementShamanRotation |
	RogueRotation |
	RetributionPaladinRotation |
	ProtectionPaladinRotation |
	ShadowPriestRotation |
	WarlockRotation |
	WarriorRotation |
	ProtectionWarriorRotation |
	SmitePriestRotation;
export type SpecRotation<T extends Spec> =
	T extends Spec.SpecBalanceDruid ? BalanceDruidRotation :
	T extends Spec.SpecElementalShaman ? ElementalShamanRotation :
	T extends Spec.SpecEnhancementShaman ? EnhancementShamanRotation :
	T extends Spec.SpecFeralDruid ? FeralDruidRotation :
	T extends Spec.SpecFeralTankDruid ? FeralTankDruidRotation :
	T extends Spec.SpecHunter ? HunterRotation :
	T extends Spec.SpecMage ? MageRotation :
	T extends Spec.SpecRogue ? RogueRotation :
	T extends Spec.SpecRetributionPaladin ? RetributionPaladinRotation :
	T extends Spec.SpecProtectionPaladin ? ProtectionPaladinRotation :
	T extends Spec.SpecShadowPriest ? ShadowPriestRotation :
	T extends Spec.SpecWarlock ? WarlockRotation :
	T extends Spec.SpecWarrior ? WarriorRotation :
	T extends Spec.SpecProtectionWarrior ? ProtectionWarriorRotation :
	T extends Spec.SpecSmitePriest ? SmitePriestRotation :
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
	T extends Spec.SpecElementalShaman ? ShamanTalents :
	T extends Spec.SpecEnhancementShaman ? ShamanTalents :
	T extends Spec.SpecFeralDruid ? DruidTalents :
	T extends Spec.SpecFeralTankDruid ? DruidTalents :
	T extends Spec.SpecHunter ? HunterTalents :
	T extends Spec.SpecMage ? MageTalents :
	T extends Spec.SpecRogue ? RogueTalents :
	T extends Spec.SpecRetributionPaladin ? PaladinTalents :
	T extends Spec.SpecProtectionPaladin ? PaladinTalents :
	T extends Spec.SpecShadowPriest ? PriestTalents :
	T extends Spec.SpecWarlock ? WarlockTalents :
	T extends Spec.SpecWarrior ? WarriorTalents :
	T extends Spec.SpecSmitePriest ? PriestTalents :
	ShamanTalents; // Should never reach this case

export type SpecOptionsUnion =
	BalanceDruidOptions |
	ElementalShamanOptions |
	EnhancementShamanOptions |
	FeralDruidOptions |
	FeralTankDruidOptions |
	HunterOptions |
	MageOptions |
	RogueOptions |
	RetributionPaladinOptions |
	ProtectionPaladinOptions |
	ShadowPriestOptions |
	WarlockOptions |
	WarriorOptions |
	ProtectionWarriorOptions |
	SmitePriestOptions;
export type SpecOptions<T extends Spec> =
	T extends Spec.SpecBalanceDruid ? BalanceDruidOptions :
	T extends Spec.SpecElementalShaman ? ElementalShamanOptions :
	T extends Spec.SpecEnhancementShaman ? EnhancementShamanOptions :
	T extends Spec.SpecFeralDruid ? FeralDruidOptions :
	T extends Spec.SpecFeralTankDruid ? FeralTankDruidOptions :
	T extends Spec.SpecHunter ? HunterOptions :
	T extends Spec.SpecMage ? MageOptions :
	T extends Spec.SpecRogue ? RogueOptions :
	T extends Spec.SpecRetributionPaladin ? RetributionPaladinOptions :
	T extends Spec.SpecProtectionPaladin ? ProtectionPaladinOptions :
	T extends Spec.SpecShadowPriest ? ShadowPriestOptions :
	T extends Spec.SpecWarlock ? WarlockOptions :
	T extends Spec.SpecWarrior ? WarriorOptions :
	T extends Spec.SpecProtectionWarrior ? ProtectionWarriorOptions :
	T extends Spec.SpecSmitePriest ? SmitePriestOptions :
	ElementalShamanOptions; // Should never reach this case

export type SpecProtoUnion =
	BalanceDruid |
	ElementalShaman |
	EnhancementShaman |
	FeralDruid |
	FeralTankDruid |
	Hunter |
	Mage |
	Rogue |
	RetributionPaladin |
	ProtectionPaladin |
	ShadowPriest |
	Warlock |
	Warrior |
	ProtectionWarrior |
	SmitePriest;
export type SpecProto<T extends Spec> =
	T extends Spec.SpecBalanceDruid ? BalanceDruid :
	T extends Spec.SpecElementalShaman ? ElementalShaman :
	T extends Spec.SpecEnhancementShaman ? EnhancementShaman :
	T extends Spec.SpecFeralDruid ? FeralDruid :
	T extends Spec.SpecFeralTankDruid ? FeralTankDruid :
	T extends Spec.SpecHunter ? Hunter :
	T extends Spec.SpecMage ? Mage :
	T extends Spec.SpecRogue ? Rogue :
	T extends Spec.SpecRetributionPaladin ? RetributionPaladin :
	T extends Spec.SpecProtectionPaladin ? ProtectionPaladin :
	T extends Spec.SpecShadowPriest ? ShadowPriest :
	T extends Spec.SpecWarlock ? Warlock :
	T extends Spec.SpecWarrior ? Warrior :
	T extends Spec.SpecProtectionWarrior ? ProtectionWarrior :
	T extends Spec.SpecSmitePriest ? SmitePriest :
	ElementalShaman; // Should never reach this case

export type SpecTypeFunctions<SpecType extends Spec> = {
	rotationCreate: () => SpecRotation<SpecType>;
	rotationEquals: (a: SpecRotation<SpecType>, b: SpecRotation<SpecType>) => boolean;
	rotationCopy: (a: SpecRotation<SpecType>) => SpecRotation<SpecType>;
	rotationToJson: (a: SpecRotation<SpecType>) => any;
	rotationFromJson: (obj: any) => SpecRotation<SpecType>;
	rotationFromPlayer: (player: Player) => SpecRotation<SpecType>;

	talentsCreate: () => SpecTalents<SpecType>;
	talentsEquals: (a: SpecTalents<SpecType>, b: SpecTalents<SpecType>) => boolean;
	talentsCopy: (a: SpecTalents<SpecType>) => SpecTalents<SpecType>;
	talentsToJson: (a: SpecTalents<SpecType>) => any;
	talentsFromJson: (obj: any) => SpecTalents<SpecType>;
	talentsFromPlayer: (player: Player) => SpecTalents<SpecType>;

	optionsCreate: () => SpecOptions<SpecType>;
	optionsEquals: (a: SpecOptions<SpecType>, b: SpecOptions<SpecType>) => boolean;
	optionsCopy: (a: SpecOptions<SpecType>) => SpecOptions<SpecType>;
	optionsToJson: (a: SpecOptions<SpecType>) => any;
	optionsFromJson: (obj: any) => SpecOptions<SpecType>;
	optionsFromPlayer: (player: Player) => SpecOptions<SpecType>;
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
		talentsFromPlayer: (player) => player.spec.oneofKind == 'balanceDruid'
			? player.spec.balanceDruid.talents || DruidTalents.create()
			: DruidTalents.create(),

		optionsCreate: () => BalanceDruidOptions.create(),
		optionsEquals: (a, b) => BalanceDruidOptions.equals(a as BalanceDruidOptions, b as BalanceDruidOptions),
		optionsCopy: (a) => BalanceDruidOptions.clone(a as BalanceDruidOptions),
		optionsToJson: (a) => BalanceDruidOptions.toJson(a as BalanceDruidOptions),
		optionsFromJson: (obj) => BalanceDruidOptions.fromJson(obj),
		optionsFromPlayer: (player) => player.spec.oneofKind == 'balanceDruid'
			? player.spec.balanceDruid.options || BalanceDruidOptions.create()
			: BalanceDruidOptions.create(),
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
		talentsFromPlayer: (player) => player.spec.oneofKind == 'elementalShaman'
			? player.spec.elementalShaman.talents || ShamanTalents.create()
			: ShamanTalents.create(),

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
		talentsFromPlayer: (player) => player.spec.oneofKind == 'enhancementShaman'
			? player.spec.enhancementShaman.talents || ShamanTalents.create()
			: ShamanTalents.create(),

		optionsCreate: () => EnhancementShamanOptions.create(),
		optionsEquals: (a, b) => EnhancementShamanOptions.equals(a as EnhancementShamanOptions, b as EnhancementShamanOptions),
		optionsCopy: (a) => EnhancementShamanOptions.clone(a as EnhancementShamanOptions),
		optionsToJson: (a) => EnhancementShamanOptions.toJson(a as EnhancementShamanOptions),
		optionsFromJson: (obj) => EnhancementShamanOptions.fromJson(obj),
		optionsFromPlayer: (player) => player.spec.oneofKind == 'enhancementShaman'
			? player.spec.enhancementShaman.options || EnhancementShamanOptions.create()
			: EnhancementShamanOptions.create(),
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
		talentsFromPlayer: (player) => player.spec.oneofKind == 'feralDruid'
			? player.spec.feralDruid.talents || DruidTalents.create()
			: DruidTalents.create(),

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
		talentsFromPlayer: (player) => player.spec.oneofKind == 'feralTankDruid'
			? player.spec.feralTankDruid.talents || DruidTalents.create()
			: DruidTalents.create(),

		optionsCreate: () => FeralTankDruidOptions.create(),
		optionsEquals: (a, b) => FeralTankDruidOptions.equals(a as FeralTankDruidOptions, b as FeralTankDruidOptions),
		optionsCopy: (a) => FeralTankDruidOptions.clone(a as FeralTankDruidOptions),
		optionsToJson: (a) => FeralTankDruidOptions.toJson(a as FeralTankDruidOptions),
		optionsFromJson: (obj) => FeralTankDruidOptions.fromJson(obj),
		optionsFromPlayer: (player) => player.spec.oneofKind == 'feralTankDruid'
			? player.spec.feralTankDruid.options || FeralTankDruidOptions.create()
			: FeralTankDruidOptions.create(),
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
		talentsFromPlayer: (player) => player.spec.oneofKind == 'hunter'
			? player.spec.hunter.talents || HunterTalents.create()
			: HunterTalents.create(),

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
		rotationFromPlayer: (player) => player.spec.oneofKind == 'mage'
			? player.spec.mage.rotation || MageRotation.create()
			: MageRotation.create(),

		talentsCreate: () => MageTalents.create(),
		talentsEquals: (a, b) => MageTalents.equals(a as MageTalents, b as MageTalents),
		talentsCopy: (a) => MageTalents.clone(a as MageTalents),
		talentsToJson: (a) => MageTalents.toJson(a as MageTalents),
		talentsFromJson: (obj) => MageTalents.fromJson(obj),
		talentsFromPlayer: (player) => player.spec.oneofKind == 'mage'
			? player.spec.mage.talents || MageTalents.create()
			: MageTalents.create(),

		optionsCreate: () => MageOptions.create(),
		optionsEquals: (a, b) => MageOptions.equals(a as MageOptions, b as MageOptions),
		optionsCopy: (a) => MageOptions.clone(a as MageOptions),
		optionsToJson: (a) => MageOptions.toJson(a as MageOptions),
		optionsFromJson: (obj) => MageOptions.fromJson(obj),
		optionsFromPlayer: (player) => player.spec.oneofKind == 'mage'
			? player.spec.mage.options || MageOptions.create()
			: MageOptions.create(),
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
		talentsFromPlayer: (player) => player.spec.oneofKind == 'retributionPaladin'
			? player.spec.retributionPaladin.talents || PaladinTalents.create()
			: PaladinTalents.create(),

		optionsCreate: () => RetributionPaladinOptions.create(),
		optionsEquals: (a, b) => RetributionPaladinOptions.equals(a as RetributionPaladinOptions, b as RetributionPaladinOptions),
		optionsCopy: (a) => RetributionPaladinOptions.clone(a as RetributionPaladinOptions),
		optionsToJson: (a) => RetributionPaladinOptions.toJson(a as RetributionPaladinOptions),
		optionsFromJson: (obj) => RetributionPaladinOptions.fromJson(obj),
		optionsFromPlayer: (player) => player.spec.oneofKind == 'retributionPaladin'
			? player.spec.retributionPaladin.options || RetributionPaladinOptions.create()
			: RetributionPaladinOptions.create(),
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
		talentsFromPlayer: (player) => player.spec.oneofKind == 'protectionPaladin'
			? player.spec.protectionPaladin.talents || PaladinTalents.create()
			: PaladinTalents.create(),

		optionsCreate: () => ProtectionPaladinOptions.create(),
		optionsEquals: (a, b) => ProtectionPaladinOptions.equals(a as ProtectionPaladinOptions, b as ProtectionPaladinOptions),
		optionsCopy: (a) => ProtectionPaladinOptions.clone(a as ProtectionPaladinOptions),
		optionsToJson: (a) => ProtectionPaladinOptions.toJson(a as ProtectionPaladinOptions),
		optionsFromJson: (obj) => ProtectionPaladinOptions.fromJson(obj),
		optionsFromPlayer: (player) => player.spec.oneofKind == 'protectionPaladin'
			? player.spec.protectionPaladin.options || ProtectionPaladinOptions.create()
			: ProtectionPaladinOptions.create(),
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
		talentsFromPlayer: (player) => player.spec.oneofKind == 'rogue'
			? player.spec.rogue.talents || RogueTalents.create()
			: RogueTalents.create(),

		optionsCreate: () => RogueOptions.create(),
		optionsEquals: (a, b) => RogueOptions.equals(a as RogueOptions, b as RogueOptions),
		optionsCopy: (a) => RogueOptions.clone(a as RogueOptions),
		optionsToJson: (a) => RogueOptions.toJson(a as RogueOptions),
		optionsFromJson: (obj) => RogueOptions.fromJson(obj),
		optionsFromPlayer: (player) => player.spec.oneofKind == 'rogue'
			? player.spec.rogue.options || RogueOptions.create()
			: RogueOptions.create(),
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
		talentsFromPlayer: (player) => player.spec.oneofKind == 'shadowPriest'
			? player.spec.shadowPriest.talents || PriestTalents.create()
			: PriestTalents.create(),

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
		talentsFromPlayer: (player) => player.spec.oneofKind == 'warlock'
			? player.spec.warlock.talents || WarlockTalents.create()
			: WarlockTalents.create(),

		optionsCreate: () => WarlockOptions.create(),
		optionsEquals: (a, b) => WarlockOptions.equals(a as WarlockOptions, b as WarlockOptions),
		optionsCopy: (a) => WarlockOptions.clone(a as WarlockOptions),
		optionsToJson: (a) => WarlockOptions.toJson(a as WarlockOptions),
		optionsFromJson: (obj) => WarlockOptions.fromJson(obj),
		optionsFromPlayer: (player) => player.spec.oneofKind == 'warlock'
			? player.spec.warlock.options || WarlockOptions.create()
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
		talentsFromPlayer: (player) => player.spec.oneofKind == 'warrior'
			? player.spec.warrior.talents || WarriorTalents.create()
			: WarriorTalents.create(),

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
		talentsFromPlayer: (player) => player.spec.oneofKind == 'protectionWarrior'
			? player.spec.protectionWarrior.talents || WarriorTalents.create()
			: WarriorTalents.create(),

		optionsCreate: () => ProtectionWarriorOptions.create(),
		optionsEquals: (a, b) => ProtectionWarriorOptions.equals(a as ProtectionWarriorOptions, b as ProtectionWarriorOptions),
		optionsCopy: (a) => ProtectionWarriorOptions.clone(a as ProtectionWarriorOptions),
		optionsToJson: (a) => ProtectionWarriorOptions.toJson(a as ProtectionWarriorOptions),
		optionsFromJson: (obj) => ProtectionWarriorOptions.fromJson(obj),
		optionsFromPlayer: (player) => player.spec.oneofKind == 'protectionWarrior'
			? player.spec.protectionWarrior.options || ProtectionWarriorOptions.create()
			: ProtectionWarriorOptions.create(),
	},
	[Spec.SpecSmitePriest]: {
		rotationCreate: () => SmitePriestRotation.create(),
		rotationEquals: (a, b) => SmitePriestRotation.equals(a as SmitePriestRotation, b as SmitePriestRotation),
		rotationCopy: (a) => SmitePriestRotation.clone(a as SmitePriestRotation),
		rotationToJson: (a) => SmitePriestRotation.toJson(a as SmitePriestRotation),
		rotationFromJson: (obj) => SmitePriestRotation.fromJson(obj),
		rotationFromPlayer: (player) => player.spec.oneofKind == 'smitePriest'
			? player.spec.smitePriest.rotation || SmitePriestRotation.create()
			: SmitePriestRotation.create(),

		talentsCreate: () => PriestTalents.create(),
		talentsEquals: (a, b) => PriestTalents.equals(a as PriestTalents, b as PriestTalents),
		talentsCopy: (a) => PriestTalents.clone(a as PriestTalents),
		talentsToJson: (a) => PriestTalents.toJson(a as PriestTalents),
		talentsFromJson: (obj) => PriestTalents.fromJson(obj),
		talentsFromPlayer: (player) => player.spec.oneofKind == 'smitePriest'
			? player.spec.smitePriest.talents || PriestTalents.create()
			: PriestTalents.create(),

		optionsCreate: () => SmitePriestOptions.create(),
		optionsEquals: (a, b) => SmitePriestOptions.equals(a as SmitePriestOptions, b as SmitePriestOptions),
		optionsCopy: (a) => SmitePriestOptions.clone(a as SmitePriestOptions),
		optionsToJson: (a) => SmitePriestOptions.toJson(a as SmitePriestOptions),
		optionsFromJson: (obj) => SmitePriestOptions.fromJson(obj),
		optionsFromPlayer: (player) => player.spec.oneofKind == 'smitePriest'
			? player.spec.smitePriest.options || SmitePriestOptions.create()
			: SmitePriestOptions.create(),
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
	[Race.RaceTroll10]: Faction.Horde,
	[Race.RaceTroll30]: Faction.Horde,
	[Race.RaceUndead]: Faction.Horde,
};

export const specToClass: Record<Spec, Class> = {
	[Spec.SpecBalanceDruid]: Class.ClassDruid,
	[Spec.SpecFeralDruid]: Class.ClassDruid,
	[Spec.SpecFeralTankDruid]: Class.ClassDruid,
	[Spec.SpecHunter]: Class.ClassHunter,
	[Spec.SpecMage]: Class.ClassMage,
	[Spec.SpecRogue]: Class.ClassRogue,
	[Spec.SpecRetributionPaladin]: Class.ClassPaladin,
	[Spec.SpecProtectionPaladin]: Class.ClassPaladin,
	[Spec.SpecShadowPriest]: Class.ClassPriest,
	[Spec.SpecSmitePriest]: Class.ClassPriest,
	[Spec.SpecElementalShaman]: Class.ClassShaman,
	[Spec.SpecEnhancementShaman]: Class.ClassShaman,
	[Spec.SpecWarlock]: Class.ClassWarlock,
	[Spec.SpecWarrior]: Class.ClassWarrior,
	[Spec.SpecProtectionWarrior]: Class.ClassWarrior,
};

const druidRaces = [
	Race.RaceNightElf,
	Race.RaceTauren,
];
const hunterRaces = [
	Race.RaceBloodElf,
	Race.RaceDraenei,
	Race.RaceDwarf,
	Race.RaceNightElf,
	Race.RaceOrc,
	Race.RaceTauren,
	Race.RaceTroll10,
	Race.RaceTroll30,
];
const mageRaces = [
	Race.RaceBloodElf,
	Race.RaceDraenei,
	Race.RaceGnome,
	Race.RaceHuman,
	Race.RaceTroll10,
	Race.RaceTroll30,
	Race.RaceUndead,
];
const paladinRaces = [
	Race.RaceBloodElf,
	Race.RaceDraenei,
	Race.RaceDwarf,
	Race.RaceHuman,
];
const priestRaces = [
	Race.RaceBloodElf,
	Race.RaceDraenei,
	Race.RaceDwarf,
	Race.RaceHuman,
	Race.RaceNightElf,
	Race.RaceOrc,
	Race.RaceTroll10,
	Race.RaceTroll30,
	Race.RaceUndead,
];
const rogueRaces = [
	Race.RaceBloodElf,
	Race.RaceDwarf,
	Race.RaceGnome,
	Race.RaceHuman,
	Race.RaceNightElf,
	Race.RaceOrc,
	Race.RaceTroll10,
	Race.RaceTroll30,
	Race.RaceUndead,
];
const shamanRaces = [
	Race.RaceDraenei,
	Race.RaceOrc,
	Race.RaceTauren,
	Race.RaceTroll10,
	Race.RaceTroll30,
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
	Race.RaceTroll10,
	Race.RaceTroll30,
	Race.RaceUndead,
];

export const specToEligibleRaces: Record<Spec, Array<Race>> = {
	[Spec.SpecBalanceDruid]: druidRaces,
	[Spec.SpecElementalShaman]: shamanRaces,
	[Spec.SpecEnhancementShaman]: shamanRaces,
	[Spec.SpecFeralDruid]: druidRaces,
	[Spec.SpecFeralTankDruid]: druidRaces,
	[Spec.SpecHunter]: hunterRaces,
	[Spec.SpecMage]: mageRaces,
	[Spec.SpecRetributionPaladin]: paladinRaces,
	[Spec.SpecProtectionPaladin]: paladinRaces,
	[Spec.SpecRogue]: rogueRaces,
	[Spec.SpecShadowPriest]: priestRaces,
	[Spec.SpecWarlock]: warlockRaces,
	[Spec.SpecWarrior]: warriorRaces,
	[Spec.SpecProtectionWarrior]: warriorRaces,
	[Spec.SpecSmitePriest]: priestRaces,
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
];
export function isTankSpec(spec: Spec): boolean {
	return tankSpecs.includes(spec);
}

// Prefixes used for storing browser data for each site. Even if a Spec is
// renamed, DO NOT change these values or people will lose their saved data.
export const specToLocalStorageKey: Record<Spec, string> = {
	[Spec.SpecBalanceDruid]: '__balance_druid',
	[Spec.SpecElementalShaman]: '__elemental_shaman',
	[Spec.SpecEnhancementShaman]: '__enhacement_shaman',
	[Spec.SpecFeralDruid]: '__feral_druid',
	[Spec.SpecFeralTankDruid]: '__feral_tank_druid',
	[Spec.SpecHunter]: '__hunter',
	[Spec.SpecMage]: '__mage',
	[Spec.SpecRetributionPaladin]: '__retribution_paladin',
	[Spec.SpecProtectionPaladin]: '__protection_paladin',
	[Spec.SpecRogue]: '__rogue',
	[Spec.SpecShadowPriest]: '__shadow_priest',
	[Spec.SpecWarlock]: '__warlock',
	[Spec.SpecWarrior]: '__warrior',
	[Spec.SpecProtectionWarrior]: '__protection_warrior',
	[Spec.SpecSmitePriest]: '__smite_priest',
};

// Returns a copy of playerOptions, with the class field set.
export function withSpecProto<SpecType extends Spec>(
	spec: Spec,
	player: Player,
	rotation: SpecRotation<SpecType>,
	talents: SpecTalents<SpecType>,
	specOptions: SpecOptions<SpecType>): Player {
	const copy = Player.clone(player);

	switch (spec) {
		case Spec.SpecBalanceDruid:
			copy.spec = {
				oneofKind: 'balanceDruid',
				balanceDruid: BalanceDruid.create({
					rotation: rotation as BalanceDruidRotation,
					talents: talents as DruidTalents,
					options: specOptions as BalanceDruidOptions,
				}),
			};
			return copy;
		case Spec.SpecElementalShaman:
			copy.spec = {
				oneofKind: 'elementalShaman',
				elementalShaman: ElementalShaman.create({
					rotation: rotation as ElementalShamanRotation,
					talents: talents as ShamanTalents,
					options: specOptions as ElementalShamanOptions,
				}),
			};
			return copy;
		case Spec.SpecEnhancementShaman:
			copy.spec = {
				oneofKind: 'enhancementShaman',
				enhancementShaman: EnhancementShaman.create({
					rotation: rotation as EnhancementShamanRotation,
					talents: talents as ShamanTalents,
					options: specOptions as ElementalShamanOptions,
				}),
			};
			return copy;
		case Spec.SpecFeralDruid:
			copy.spec = {
				oneofKind: 'feralDruid',
				feralDruid: FeralDruid.create({
					rotation: rotation as FeralDruidRotation,
					talents: talents as DruidTalents,
					options: specOptions as FeralDruidOptions,
				}),
			};
			return copy;
		case Spec.SpecFeralTankDruid:
			copy.spec = {
				oneofKind: 'feralTankDruid',
				feralTankDruid: FeralTankDruid.create({
					rotation: rotation as FeralTankDruidRotation,
					talents: talents as DruidTalents,
					options: specOptions as FeralTankDruidOptions,
				}),
			};
			return copy;
		case Spec.SpecHunter:
			copy.spec = {
				oneofKind: 'hunter',
				hunter: Hunter.create({
					rotation: rotation as HunterRotation,
					talents: talents as HunterTalents,
					options: specOptions as HunterOptions,
				}),
			};
			return copy;
		case Spec.SpecMage:
			copy.spec = {
				oneofKind: 'mage',
				mage: Mage.create({
					rotation: rotation as MageRotation,
					talents: talents as MageTalents,
					options: specOptions as MageOptions,
				}),
			};
			return copy;
		case Spec.SpecRetributionPaladin:
			copy.spec = {
				oneofKind: 'retributionPaladin',
				retributionPaladin: RetributionPaladin.create({
					rotation: rotation as RetributionPaladinRotation,
					talents: talents as PaladinTalents,
					options: specOptions as RetributionPaladinOptions,
				}),
			};
			return copy;
		case Spec.SpecProtectionPaladin:
			copy.spec = {
				oneofKind: 'protectionPaladin',
				protectionPaladin: ProtectionPaladin.create({
					rotation: rotation as ProtectionPaladinRotation,
					talents: talents as PaladinTalents,
					options: specOptions as ProtectionPaladinOptions,
				}),
			};
			return copy;
		case Spec.SpecRogue:
			copy.spec = {
				oneofKind: 'rogue',
				rogue: Rogue.create({
					rotation: rotation as RogueRotation,
					talents: talents as RogueTalents,
					options: specOptions as RogueOptions,
				}),
			};
			return copy;
		case Spec.SpecShadowPriest:
			copy.spec = {
				oneofKind: 'shadowPriest',
				shadowPriest: ShadowPriest.create({
					rotation: rotation as ShadowPriestRotation,
					talents: talents as PriestTalents,
					options: specOptions as ShadowPriestOptions,
				}),
			};
			return copy;
		case Spec.SpecWarlock:
			copy.spec = {
				oneofKind: 'warlock',
				warlock: Warlock.create({
					rotation: rotation as WarlockRotation,
					talents: talents as WarlockTalents,
					options: specOptions as WarlockOptions,
				}),
			};
			return copy;
		case Spec.SpecWarrior:
			copy.spec = {
				oneofKind: 'warrior',
				warrior: Warrior.create({
					rotation: rotation as WarriorRotation,
					talents: talents as WarriorTalents,
					options: specOptions as WarriorOptions,
				}),
			};
			return copy;
		case Spec.SpecProtectionWarrior:
			copy.spec = {
				oneofKind: 'protectionWarrior',
				protectionWarrior: ProtectionWarrior.create({
					rotation: rotation as ProtectionWarriorRotation,
					talents: talents as WarriorTalents,
					options: specOptions as ProtectionWarriorOptions,
				}),
			};
			return copy;
		case Spec.SpecSmitePriest:
			copy.spec = {
				oneofKind: 'smitePriest',
				smitePriest: SmitePriest.create({
					rotation: rotation as SmitePriestRotation,
					talents: talents as PriestTalents,
					options: specOptions as SmitePriestOptions,
				}),
			};
			return copy;
	}
}

export function playerToSpec(player: Player): Spec {
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

	throw new Error('Unable to parse spec from player proto: ' + JSON.stringify(Player.toJson(player), null, 2));
}

const classToMaxArmorType: Record<Class, ArmorType> = {
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

const classToEligibleRangedWeaponTypes: Record<Class, Array<RangedWeaponType>> = {
	[Class.ClassUnknown]: [],
	[Class.ClassDruid]: [RangedWeaponType.RangedWeaponTypeIdol],
	[Class.ClassHunter]: [
		RangedWeaponType.RangedWeaponTypeBow,
		RangedWeaponType.RangedWeaponTypeCrossbow,
		RangedWeaponType.RangedWeaponTypeGun,
		RangedWeaponType.RangedWeaponTypeThrown,
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

const classToEligibleWeaponTypes: Record<Class, Array<EligibleWeaponType>> = {
	[Class.ClassUnknown]: [],
	[Class.ClassDruid]: [
		{ weaponType: WeaponType.WeaponTypeDagger },
		{ weaponType: WeaponType.WeaponTypeFist },
		{ weaponType: WeaponType.WeaponTypeMace, canUseTwoHand: true },
		{ weaponType: WeaponType.WeaponTypeOffHand },
		{ weaponType: WeaponType.WeaponTypeStaff, canUseTwoHand: true },
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

// Custom functions for determining the EP value of meta gem effects.
// Default meta effect EP value is 0, so just handle the ones relevant to your spec.
const metaGemEffectEPs: Partial<Record<Spec, (gem: Gem, playerStats: Stats) => number>> = {
	[Spec.SpecBalanceDruid]: (gem, playerStats) => {
		if (gem.id == Gems.CHAOTIC_SKYFIRE_DIAMOND.id) {
			// TODO: Fix this
			return (((playerStats.getStat(Stat.StatSpellPower) * 0.795) + 603) * 2 * (playerStats.getStat(Stat.StatSpellCrit) / 2208) * 0.045) / 0.795;
		}

		return 0;
	},
	[Spec.SpecElementalShaman]: (gem, playerStats) => {
		if (gem.id == Gems.CHAOTIC_SKYFIRE_DIAMOND.id) {
			return (((playerStats.getStat(Stat.StatSpellPower) * 0.795) + 603) * 2 * (playerStats.getStat(Stat.StatSpellCrit) / 2208) * 0.045) / 0.795;
		}

		return 0;
	},
};

export function getMetaGemEffectEP(spec: Spec, gem: Gem, playerStats: Stats) {
	if (metaGemEffectEPs[spec]) {
		return metaGemEffectEPs[spec]!(gem, playerStats);
	} else {
		return 0;
	}
}

// Returns true if this item may be equipped in at least 1 slot for the given Spec.
export function canEquipItem(item: Item, spec: Spec, slot: ItemSlot | undefined): boolean {
	const playerClass = specToClass[spec];
	if (item.classAllowlist.length > 0 && !item.classAllowlist.includes(playerClass)) {
		return false;
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

		return true;
	}

	if (item.type == ItemType.ItemTypeRanged) {
		return classToEligibleRangedWeaponTypes[playerClass].includes(item.rangedWeaponType);
	}

	// At this point, we know the item is an armor piece (feet, chest, legs, etc).
	return classToMaxArmorType[playerClass] >= item.armorType;
}

const itemTypeToSlotsMap: Partial<Record<ItemType, Array<ItemSlot>>> = {
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
		if ([HandType.HandTypeMainHand, HandType.HandTypeTwoHand].includes(item.handType)) {
			return [ItemSlot.ItemSlotMainHand];
		} else if (item.handType == HandType.HandTypeOffHand) {
			return [ItemSlot.ItemSlotOffHand];
		} else {
			return [ItemSlot.ItemSlotMainHand, ItemSlot.ItemSlotOffHand];
		}
	}

	// Should never reach here
	throw new Error('Could not find item slots for item: ' + Item.toJsonString(item));
};

// Returns whether the given main-hand and off-hand items can be worn at the
// same time.
export function validWeaponCombo(mainHand: Item | null | undefined, offHand: Item | null | undefined): boolean {
	if (mainHand == null || offHand == null) {
		return true;
	}

	if (mainHand.handType == HandType.HandTypeTwoHand) {
		return false;
	}

	return true;
}

// Returns all item slots to which the enchant might be applied.
// 
// Note that this alone is not enough; some items have further restrictions,
// e.g. some weapon enchants may only be applied to 2H weapons.
export function getEligibleEnchantSlots(enchant: Enchant): Array<ItemSlot> {
	if (itemTypeToSlotsMap[enchant.type]) {
		return itemTypeToSlotsMap[enchant.type]!;
	}

	if (enchant.type == ItemType.ItemTypeWeapon) {
		return [ItemSlot.ItemSlotMainHand, ItemSlot.ItemSlotOffHand];
	}

	// Should never reach here
	throw new Error('Could not find item slots for enchant: ' + Enchant.toJsonString(enchant));
};

export function enchantAppliesToItem(enchant: Enchant, item: Item): boolean {
	const sharedSlots = intersection(getEligibleEnchantSlots(enchant), getEligibleItemSlots(item));
	if (sharedSlots.length == 0)
		return false;

	if (enchant.enchantType == EnchantType.EnchantTypeTwoHand && item.handType != HandType.HandTypeTwoHand)
		return false;

	if ((enchant.enchantType == EnchantType.EnchantTypeShield) != (item.weaponType == WeaponType.WeaponTypeShield))
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

export const NO_TARGET = -1;

export function newRaidTarget(raidIndex: number): RaidTarget {
	return RaidTarget.create({
		targetIndex: raidIndex,
	});
}

export function emptyRaidTarget(): RaidTarget {
	return newRaidTarget(NO_TARGET);
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
		{ spec: Spec.SpecBalanceDruid, blessings: [Blessings.BlessingOfKings, Blessings.BlessingOfSalvation, Blessings.BlessingOfWisdom] },
		{ spec: Spec.SpecFeralDruid, blessings: [Blessings.BlessingOfKings, Blessings.BlessingOfSalvation, Blessings.BlessingOfMight, Blessings.BlessingOfWisdom] },
		{ spec: Spec.SpecFeralTankDruid, blessings: [Blessings.BlessingOfKings, Blessings.BlessingOfMight, Blessings.BlessingOfSanctuary] },
		{ spec: Spec.SpecHunter, blessings: [Blessings.BlessingOfKings, Blessings.BlessingOfSalvation, Blessings.BlessingOfMight, Blessings.BlessingOfWisdom] },
		{ spec: Spec.SpecMage, blessings: [Blessings.BlessingOfKings, Blessings.BlessingOfSalvation, Blessings.BlessingOfWisdom] },
		{ spec: Spec.SpecRetributionPaladin, blessings: [Blessings.BlessingOfKings, Blessings.BlessingOfMight, Blessings.BlessingOfSalvation, Blessings.BlessingOfWisdom] },
		{ spec: Spec.SpecProtectionPaladin, blessings: [Blessings.BlessingOfKings, Blessings.BlessingOfSanctuary, Blessings.BlessingOfWisdom, Blessings.BlessingOfMight] },
		{ spec: Spec.SpecShadowPriest, blessings: [Blessings.BlessingOfKings, Blessings.BlessingOfSalvation, Blessings.BlessingOfWisdom] },
		{ spec: Spec.SpecSmitePriest, blessings: [Blessings.BlessingOfKings, Blessings.BlessingOfSalvation, Blessings.BlessingOfWisdom] },
		{ spec: Spec.SpecRogue, blessings: [Blessings.BlessingOfKings, Blessings.BlessingOfSalvation, Blessings.BlessingOfMight] },
		{ spec: Spec.SpecElementalShaman, blessings: [Blessings.BlessingOfKings, Blessings.BlessingOfSalvation, Blessings.BlessingOfWisdom] },
		{ spec: Spec.SpecEnhancementShaman, blessings: [Blessings.BlessingOfKings, Blessings.BlessingOfSalvation, Blessings.BlessingOfMight, Blessings.BlessingOfWisdom] },
		{ spec: Spec.SpecWarlock, blessings: [Blessings.BlessingOfKings, Blessings.BlessingOfSalvation, Blessings.BlessingOfWisdom] },
		{ spec: Spec.SpecWarrior, blessings: [Blessings.BlessingOfKings, Blessings.BlessingOfSalvation, Blessings.BlessingOfMight] },
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
