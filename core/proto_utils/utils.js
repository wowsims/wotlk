import { REPO_NAME } from '/wotlk/core/constants/other.js';
import { camelToSnakeCase } from '/wotlk/core/utils.js';
import { getEnumValues } from '/wotlk/core/utils.js';
import { intersection } from '/wotlk/core/utils.js';
import { maxIndex } from '/wotlk/core/utils.js';
import { sum } from '/wotlk/core/utils.js';
import { Player } from '/wotlk/core/proto/api.js';
import { ResourceType } from '/wotlk/core/proto/api.js';
import { ArmorType } from '/wotlk/core/proto/common.js';
import { Class } from '/wotlk/core/proto/common.js';
import { Enchant } from '/wotlk/core/proto/common.js';
import { EnchantType } from '/wotlk/core/proto/common.js';
import { HandType } from '/wotlk/core/proto/common.js';
import { ItemSlot } from '/wotlk/core/proto/common.js';
import { ItemType } from '/wotlk/core/proto/common.js';
import { Item } from '/wotlk/core/proto/common.js';
import { Race } from '/wotlk/core/proto/common.js';
import { Faction } from '/wotlk/core/proto/common.js';
import { RaidTarget } from '/wotlk/core/proto/common.js';
import { RangedWeaponType } from '/wotlk/core/proto/common.js';
import { Spec } from '/wotlk/core/proto/common.js';
import { Stat } from '/wotlk/core/proto/common.js';
import { WeaponType } from '/wotlk/core/proto/common.js';
import { Blessings } from '/wotlk/core/proto/paladin.js';
import { BlessingsAssignment } from '/wotlk/core/proto/ui.js';
import { BlessingsAssignments } from '/wotlk/core/proto/ui.js';
import * as Gems from '/wotlk/core/proto_utils/gems.js';
import { BalanceDruid, FeralDruid, FeralTankDruid, BalanceDruid_Rotation as BalanceDruidRotation, FeralDruid_Rotation as FeralDruidRotation, FeralTankDruid_Rotation as FeralTankDruidRotation, DruidTalents, BalanceDruid_Options as BalanceDruidOptions, FeralDruid_Options as FeralDruidOptions, FeralTankDruid_Options as FeralTankDruidOptions } from '/wotlk/core/proto/druid.js';
import { ElementalShaman, EnhancementShaman_Rotation as EnhancementShamanRotation, ElementalShaman_Rotation as ElementalShamanRotation, ShamanTalents, ElementalShaman_Options as ElementalShamanOptions, EnhancementShaman_Options as EnhancementShamanOptions, EnhancementShaman } from '/wotlk/core/proto/shaman.js';
import { Hunter, Hunter_Rotation as HunterRotation, HunterTalents, Hunter_Options as HunterOptions } from '/wotlk/core/proto/hunter.js';
import { Mage, Mage_Rotation as MageRotation, MageTalents, Mage_Options as MageOptions } from '/wotlk/core/proto/mage.js';
import { Rogue, Rogue_Rotation as RogueRotation, RogueTalents, Rogue_Options as RogueOptions } from '/wotlk/core/proto/rogue.js';
import { RetributionPaladin, RetributionPaladin_Rotation as RetributionPaladinRotation, PaladinTalents, RetributionPaladin_Options as RetributionPaladinOptions } from '/wotlk/core/proto/paladin.js';
import { ProtectionPaladin, ProtectionPaladin_Rotation as ProtectionPaladinRotation, ProtectionPaladin_Options as ProtectionPaladinOptions } from '/wotlk/core/proto/paladin.js';
import { ShadowPriest, SmitePriest_Rotation as SmitePriestRotation, ShadowPriest_Rotation as ShadowPriestRotation, PriestTalents, ShadowPriest_Options as ShadowPriestOptions, SmitePriest_Options as SmitePriestOptions, SmitePriest } from '/wotlk/core/proto/priest.js';
import { Warlock, Warlock_Rotation as WarlockRotation, WarlockTalents, Warlock_Options as WarlockOptions } from '/wotlk/core/proto/warlock.js';
import { Warrior, Warrior_Rotation as WarriorRotation, WarriorTalents, Warrior_Options as WarriorOptions } from '/wotlk/core/proto/warrior.js';
import { DeathKnight, DeathKnight_Rotation as DeathKnightRotation, DeathKnightTalents, DeathKnight_Options as DeathKnightOptions } from '/wotlk/core/proto/deathknight.js';
import { DeathKnightTank, DeathKnightTank_Rotation as DeathKnightTankRotation, DeathKnightTank_Options as DeathKnightTankOptions } from '/wotlk/core/proto/deathknight.js';
import { ProtectionWarrior, ProtectionWarrior_Rotation as ProtectionWarriorRotation, ProtectionWarrior_Options as ProtectionWarriorOptions } from '/wotlk/core/proto/warrior.js';
export const NUM_SPECS = getEnumValues(Spec).length;
// The order in which specs should be presented, when it matters.
// Currently this is only used for the order of the paladin blessings UI.
export const naturalSpecOrder = [
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
    Spec.SpecDeathKnight,
    Spec.SpecDeathKnightTank,
];
export const specNames = {
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
    [Spec.SpecDeathKnight]: 'Death Knight',
    [Spec.SpecDeathKnightTank]: 'Death Knight Tank',
};
export const classColors = {
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
    [Class.ClassDeathKnight]: '#c41e3a'
};
export const specIconsLarge = {
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
    [Spec.SpecDeathKnight]: 'https://wow.zamimg.com/images/wow/icons/medium/class_deathknight.jpg',
    [Spec.SpecDeathKnightTank]: 'https://wow.zamimg.com/images/wow/icons/medium/class_deathknight.jpg',
};
export const talentTreeIcons = {
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
    [Class.ClassDeathKnight]: [
        'https://wow.zamimg.com/images/wow/icons/medium/spell_deathknight_bloodpresence.jpg',
        'https://wow.zamimg.com/images/wow/icons/medium/spell_deathknight_frostpresence.jpg',
        'https://wow.zamimg.com/images/wow/icons/medium/spell_deathknight_unholypresence.jpg',
    ],
};
export const titleIcons = {
    [Spec.SpecBalanceDruid]: '/wotlk/assets/img/balance_druid_icon.png',
    [Spec.SpecElementalShaman]: '/wotlk/assets/img/elemental_shaman_icon.png',
    [Spec.SpecEnhancementShaman]: '/wotlk/assets/img/enhancement_shaman_icon.png',
    [Spec.SpecFeralDruid]: '/wotlk/assets/img/feral_druid_icon.png',
    [Spec.SpecFeralTankDruid]: '/wotlk/assets/img/feral_druid_tank_icon.png',
    [Spec.SpecHunter]: '/wotlk/assets/img/hunter_icon.png',
    [Spec.SpecMage]: '/wotlk/assets/img/mage_icon.png',
    [Spec.SpecRogue]: '/wotlk/assets/img/rogue_icon.png',
    [Spec.SpecRetributionPaladin]: '/wotlk/assets/img/retribution_icon.png',
    [Spec.SpecProtectionPaladin]: '/wotlk/assets/img/protection_paladin_icon.png',
    [Spec.SpecShadowPriest]: '/wotlk/assets/img/shadow_priest_icon.png',
    [Spec.SpecWarlock]: '/wotlk/assets/img/warlock_icon.png',
    [Spec.SpecWarrior]: '/wotlk/assets/img/warrior_icon.png',
    [Spec.SpecProtectionWarrior]: '/wotlk/assets/img/protection_warrior_icon.png',
    [Spec.SpecSmitePriest]: '/wotlk/assets/img/smite_priest_icon.png',
    [Spec.SpecDeathKnight]: 'https://wow.zamimg.com/images/wow/icons/medium/class_deathknight.jpg',
    [Spec.SpecDeathKnightTank]: 'https://wow.zamimg.com/images/wow/icons/medium/class_deathknight.jpg',
};
export const raidSimIcon = '/wotlk/assets/img/raid_icon.png';
// Returns the index of the talent tree (0, 1, or 2) that has the most points.
export function getTalentTree(talentsString) {
    const trees = talentsString.split('-');
    const points = trees.map(tree => sum([...tree].map(char => parseInt(char))));
    return maxIndex(points) || 0;
}
// Returns the index of the talent tree (0, 1, or 2) that has the most points.
export function getTalentTreeIcon(spec, talentsString) {
    const talentTreeIdx = getTalentTree(talentsString);
    return talentTreeIcons[specToClass[spec]][talentTreeIdx];
}
// Gets the URL for the individual sim corresponding to the given spec.
const specSiteUrlTemplate = new URL(`${window.location.protocol}//${window.location.host}/${REPO_NAME}/SPEC/`);
export function getSpecSiteUrl(spec) {
    let specString = Spec[spec]; // Returns 'SpecBalanceDruid' for BalanceDruid.
    specString = specString.substring('Spec'.length); // 'BalanceDruid'
    specString = camelToSnakeCase(specString); // 'balance_druid'
    return specSiteUrlTemplate.toString().replace('SPEC', specString);
}
export const raidSimSiteUrl = new URL(`${window.location.protocol}//${window.location.host}/${REPO_NAME}/raid/`).toString();
export const specTypeFunctions = {
    [Spec.SpecBalanceDruid]: {
        rotationCreate: () => BalanceDruidRotation.create(),
        rotationEquals: (a, b) => BalanceDruidRotation.equals(a, b),
        rotationCopy: (a) => BalanceDruidRotation.clone(a),
        rotationToJson: (a) => BalanceDruidRotation.toJson(a),
        rotationFromJson: (obj) => BalanceDruidRotation.fromJson(obj),
        rotationFromPlayer: (player) => player.spec.oneofKind == 'balanceDruid'
            ? player.spec.balanceDruid.rotation || BalanceDruidRotation.create()
            : BalanceDruidRotation.create(),
        talentsCreate: () => DruidTalents.create(),
        talentsEquals: (a, b) => DruidTalents.equals(a, b),
        talentsCopy: (a) => DruidTalents.clone(a),
        talentsToJson: (a) => DruidTalents.toJson(a),
        talentsFromJson: (obj) => DruidTalents.fromJson(obj),
        talentsFromPlayer: (player) => player.spec.oneofKind == 'balanceDruid'
            ? player.spec.balanceDruid.talents || DruidTalents.create()
            : DruidTalents.create(),
        optionsCreate: () => BalanceDruidOptions.create(),
        optionsEquals: (a, b) => BalanceDruidOptions.equals(a, b),
        optionsCopy: (a) => BalanceDruidOptions.clone(a),
        optionsToJson: (a) => BalanceDruidOptions.toJson(a),
        optionsFromJson: (obj) => BalanceDruidOptions.fromJson(obj),
        optionsFromPlayer: (player) => player.spec.oneofKind == 'balanceDruid'
            ? player.spec.balanceDruid.options || BalanceDruidOptions.create()
            : BalanceDruidOptions.create(),
    },
    [Spec.SpecElementalShaman]: {
        rotationCreate: () => ElementalShamanRotation.create(),
        rotationEquals: (a, b) => ElementalShamanRotation.equals(a, b),
        rotationCopy: (a) => ElementalShamanRotation.clone(a),
        rotationToJson: (a) => ElementalShamanRotation.toJson(a),
        rotationFromJson: (obj) => ElementalShamanRotation.fromJson(obj),
        rotationFromPlayer: (player) => player.spec.oneofKind == 'elementalShaman'
            ? player.spec.elementalShaman.rotation || ElementalShamanRotation.create()
            : ElementalShamanRotation.create(),
        talentsCreate: () => ShamanTalents.create(),
        talentsEquals: (a, b) => ShamanTalents.equals(a, b),
        talentsCopy: (a) => ShamanTalents.clone(a),
        talentsToJson: (a) => ShamanTalents.toJson(a),
        talentsFromJson: (obj) => ShamanTalents.fromJson(obj),
        talentsFromPlayer: (player) => player.spec.oneofKind == 'elementalShaman'
            ? player.spec.elementalShaman.talents || ShamanTalents.create()
            : ShamanTalents.create(),
        optionsCreate: () => ElementalShamanOptions.create(),
        optionsEquals: (a, b) => ElementalShamanOptions.equals(a, b),
        optionsCopy: (a) => ElementalShamanOptions.clone(a),
        optionsToJson: (a) => ElementalShamanOptions.toJson(a),
        optionsFromJson: (obj) => ElementalShamanOptions.fromJson(obj),
        optionsFromPlayer: (player) => player.spec.oneofKind == 'elementalShaman'
            ? player.spec.elementalShaman.options || ElementalShamanOptions.create()
            : ElementalShamanOptions.create(),
    },
    [Spec.SpecEnhancementShaman]: {
        rotationCreate: () => EnhancementShamanRotation.create(),
        rotationEquals: (a, b) => EnhancementShamanRotation.equals(a, b),
        rotationCopy: (a) => EnhancementShamanRotation.clone(a),
        rotationToJson: (a) => EnhancementShamanRotation.toJson(a),
        rotationFromJson: (obj) => EnhancementShamanRotation.fromJson(obj),
        rotationFromPlayer: (player) => player.spec.oneofKind == 'enhancementShaman'
            ? player.spec.enhancementShaman.rotation || EnhancementShamanRotation.create()
            : EnhancementShamanRotation.create(),
        talentsCreate: () => ShamanTalents.create(),
        talentsEquals: (a, b) => ShamanTalents.equals(a, b),
        talentsCopy: (a) => ShamanTalents.clone(a),
        talentsToJson: (a) => ShamanTalents.toJson(a),
        talentsFromJson: (obj) => ShamanTalents.fromJson(obj),
        talentsFromPlayer: (player) => player.spec.oneofKind == 'enhancementShaman'
            ? player.spec.enhancementShaman.talents || ShamanTalents.create()
            : ShamanTalents.create(),
        optionsCreate: () => EnhancementShamanOptions.create(),
        optionsEquals: (a, b) => EnhancementShamanOptions.equals(a, b),
        optionsCopy: (a) => EnhancementShamanOptions.clone(a),
        optionsToJson: (a) => EnhancementShamanOptions.toJson(a),
        optionsFromJson: (obj) => EnhancementShamanOptions.fromJson(obj),
        optionsFromPlayer: (player) => player.spec.oneofKind == 'enhancementShaman'
            ? player.spec.enhancementShaman.options || EnhancementShamanOptions.create()
            : EnhancementShamanOptions.create(),
    },
    [Spec.SpecFeralDruid]: {
        rotationCreate: () => FeralDruidRotation.create(),
        rotationEquals: (a, b) => FeralDruidRotation.equals(a, b),
        rotationCopy: (a) => FeralDruidRotation.clone(a),
        rotationToJson: (a) => FeralDruidRotation.toJson(a),
        rotationFromJson: (obj) => FeralDruidRotation.fromJson(obj),
        rotationFromPlayer: (player) => player.spec.oneofKind == 'feralDruid'
            ? player.spec.feralDruid.rotation || FeralDruidRotation.create()
            : FeralDruidRotation.create(),
        talentsCreate: () => DruidTalents.create(),
        talentsEquals: (a, b) => DruidTalents.equals(a, b),
        talentsCopy: (a) => DruidTalents.clone(a),
        talentsToJson: (a) => DruidTalents.toJson(a),
        talentsFromJson: (obj) => DruidTalents.fromJson(obj),
        talentsFromPlayer: (player) => player.spec.oneofKind == 'feralDruid'
            ? player.spec.feralDruid.talents || DruidTalents.create()
            : DruidTalents.create(),
        optionsCreate: () => FeralDruidOptions.create(),
        optionsEquals: (a, b) => FeralDruidOptions.equals(a, b),
        optionsCopy: (a) => FeralDruidOptions.clone(a),
        optionsToJson: (a) => FeralDruidOptions.toJson(a),
        optionsFromJson: (obj) => FeralDruidOptions.fromJson(obj),
        optionsFromPlayer: (player) => player.spec.oneofKind == 'feralDruid'
            ? player.spec.feralDruid.options || FeralDruidOptions.create()
            : FeralDruidOptions.create(),
    },
    [Spec.SpecFeralTankDruid]: {
        rotationCreate: () => FeralTankDruidRotation.create(),
        rotationEquals: (a, b) => FeralTankDruidRotation.equals(a, b),
        rotationCopy: (a) => FeralTankDruidRotation.clone(a),
        rotationToJson: (a) => FeralTankDruidRotation.toJson(a),
        rotationFromJson: (obj) => FeralTankDruidRotation.fromJson(obj),
        rotationFromPlayer: (player) => player.spec.oneofKind == 'feralTankDruid'
            ? player.spec.feralTankDruid.rotation || FeralTankDruidRotation.create()
            : FeralTankDruidRotation.create(),
        talentsCreate: () => DruidTalents.create(),
        talentsEquals: (a, b) => DruidTalents.equals(a, b),
        talentsCopy: (a) => DruidTalents.clone(a),
        talentsToJson: (a) => DruidTalents.toJson(a),
        talentsFromJson: (obj) => DruidTalents.fromJson(obj),
        talentsFromPlayer: (player) => player.spec.oneofKind == 'feralTankDruid'
            ? player.spec.feralTankDruid.talents || DruidTalents.create()
            : DruidTalents.create(),
        optionsCreate: () => FeralTankDruidOptions.create(),
        optionsEquals: (a, b) => FeralTankDruidOptions.equals(a, b),
        optionsCopy: (a) => FeralTankDruidOptions.clone(a),
        optionsToJson: (a) => FeralTankDruidOptions.toJson(a),
        optionsFromJson: (obj) => FeralTankDruidOptions.fromJson(obj),
        optionsFromPlayer: (player) => player.spec.oneofKind == 'feralTankDruid'
            ? player.spec.feralTankDruid.options || FeralTankDruidOptions.create()
            : FeralTankDruidOptions.create(),
    },
    [Spec.SpecHunter]: {
        rotationCreate: () => HunterRotation.create(),
        rotationEquals: (a, b) => HunterRotation.equals(a, b),
        rotationCopy: (a) => HunterRotation.clone(a),
        rotationToJson: (a) => HunterRotation.toJson(a),
        rotationFromJson: (obj) => HunterRotation.fromJson(obj),
        rotationFromPlayer: (player) => player.spec.oneofKind == 'hunter'
            ? player.spec.hunter.rotation || HunterRotation.create()
            : HunterRotation.create(),
        talentsCreate: () => HunterTalents.create(),
        talentsEquals: (a, b) => HunterTalents.equals(a, b),
        talentsCopy: (a) => HunterTalents.clone(a),
        talentsToJson: (a) => HunterTalents.toJson(a),
        talentsFromJson: (obj) => HunterTalents.fromJson(obj),
        talentsFromPlayer: (player) => player.spec.oneofKind == 'hunter'
            ? player.spec.hunter.talents || HunterTalents.create()
            : HunterTalents.create(),
        optionsCreate: () => HunterOptions.create(),
        optionsEquals: (a, b) => HunterOptions.equals(a, b),
        optionsCopy: (a) => HunterOptions.clone(a),
        optionsToJson: (a) => HunterOptions.toJson(a),
        optionsFromJson: (obj) => HunterOptions.fromJson(obj),
        optionsFromPlayer: (player) => player.spec.oneofKind == 'hunter'
            ? player.spec.hunter.options || HunterOptions.create()
            : HunterOptions.create(),
    },
    [Spec.SpecMage]: {
        rotationCreate: () => MageRotation.create(),
        rotationEquals: (a, b) => MageRotation.equals(a, b),
        rotationCopy: (a) => MageRotation.clone(a),
        rotationToJson: (a) => MageRotation.toJson(a),
        rotationFromJson: (obj) => MageRotation.fromJson(obj),
        rotationFromPlayer: (player) => player.spec.oneofKind == 'mage'
            ? player.spec.mage.rotation || MageRotation.create()
            : MageRotation.create(),
        talentsCreate: () => MageTalents.create(),
        talentsEquals: (a, b) => MageTalents.equals(a, b),
        talentsCopy: (a) => MageTalents.clone(a),
        talentsToJson: (a) => MageTalents.toJson(a),
        talentsFromJson: (obj) => MageTalents.fromJson(obj),
        talentsFromPlayer: (player) => player.spec.oneofKind == 'mage'
            ? player.spec.mage.talents || MageTalents.create()
            : MageTalents.create(),
        optionsCreate: () => MageOptions.create(),
        optionsEquals: (a, b) => MageOptions.equals(a, b),
        optionsCopy: (a) => MageOptions.clone(a),
        optionsToJson: (a) => MageOptions.toJson(a),
        optionsFromJson: (obj) => MageOptions.fromJson(obj),
        optionsFromPlayer: (player) => player.spec.oneofKind == 'mage'
            ? player.spec.mage.options || MageOptions.create()
            : MageOptions.create(),
    },
    [Spec.SpecRetributionPaladin]: {
        rotationCreate: () => RetributionPaladinRotation.create(),
        rotationEquals: (a, b) => RetributionPaladinRotation.equals(a, b),
        rotationCopy: (a) => RetributionPaladinRotation.clone(a),
        rotationToJson: (a) => RetributionPaladinRotation.toJson(a),
        rotationFromJson: (obj) => RetributionPaladinRotation.fromJson(obj),
        rotationFromPlayer: (player) => player.spec.oneofKind == 'retributionPaladin'
            ? player.spec.retributionPaladin.rotation || RetributionPaladinRotation.create()
            : RetributionPaladinRotation.create(),
        talentsCreate: () => PaladinTalents.create(),
        talentsEquals: (a, b) => PaladinTalents.equals(a, b),
        talentsCopy: (a) => PaladinTalents.clone(a),
        talentsToJson: (a) => PaladinTalents.toJson(a),
        talentsFromJson: (obj) => PaladinTalents.fromJson(obj),
        talentsFromPlayer: (player) => player.spec.oneofKind == 'retributionPaladin'
            ? player.spec.retributionPaladin.talents || PaladinTalents.create()
            : PaladinTalents.create(),
        optionsCreate: () => RetributionPaladinOptions.create(),
        optionsEquals: (a, b) => RetributionPaladinOptions.equals(a, b),
        optionsCopy: (a) => RetributionPaladinOptions.clone(a),
        optionsToJson: (a) => RetributionPaladinOptions.toJson(a),
        optionsFromJson: (obj) => RetributionPaladinOptions.fromJson(obj),
        optionsFromPlayer: (player) => player.spec.oneofKind == 'retributionPaladin'
            ? player.spec.retributionPaladin.options || RetributionPaladinOptions.create()
            : RetributionPaladinOptions.create(),
    },
    [Spec.SpecProtectionPaladin]: {
        rotationCreate: () => ProtectionPaladinRotation.create(),
        rotationEquals: (a, b) => ProtectionPaladinRotation.equals(a, b),
        rotationCopy: (a) => ProtectionPaladinRotation.clone(a),
        rotationToJson: (a) => ProtectionPaladinRotation.toJson(a),
        rotationFromJson: (obj) => ProtectionPaladinRotation.fromJson(obj),
        rotationFromPlayer: (player) => player.spec.oneofKind == 'protectionPaladin'
            ? player.spec.protectionPaladin.rotation || ProtectionPaladinRotation.create()
            : ProtectionPaladinRotation.create(),
        talentsCreate: () => PaladinTalents.create(),
        talentsEquals: (a, b) => PaladinTalents.equals(a, b),
        talentsCopy: (a) => PaladinTalents.clone(a),
        talentsToJson: (a) => PaladinTalents.toJson(a),
        talentsFromJson: (obj) => PaladinTalents.fromJson(obj),
        talentsFromPlayer: (player) => player.spec.oneofKind == 'protectionPaladin'
            ? player.spec.protectionPaladin.talents || PaladinTalents.create()
            : PaladinTalents.create(),
        optionsCreate: () => ProtectionPaladinOptions.create(),
        optionsEquals: (a, b) => ProtectionPaladinOptions.equals(a, b),
        optionsCopy: (a) => ProtectionPaladinOptions.clone(a),
        optionsToJson: (a) => ProtectionPaladinOptions.toJson(a),
        optionsFromJson: (obj) => ProtectionPaladinOptions.fromJson(obj),
        optionsFromPlayer: (player) => player.spec.oneofKind == 'protectionPaladin'
            ? player.spec.protectionPaladin.options || ProtectionPaladinOptions.create()
            : ProtectionPaladinOptions.create(),
    },
    [Spec.SpecRogue]: {
        rotationCreate: () => RogueRotation.create(),
        rotationEquals: (a, b) => RogueRotation.equals(a, b),
        rotationCopy: (a) => RogueRotation.clone(a),
        rotationToJson: (a) => RogueRotation.toJson(a),
        rotationFromJson: (obj) => RogueRotation.fromJson(obj),
        rotationFromPlayer: (player) => player.spec.oneofKind == 'rogue'
            ? player.spec.rogue.rotation || RogueRotation.create()
            : RogueRotation.create(),
        talentsCreate: () => RogueTalents.create(),
        talentsEquals: (a, b) => RogueTalents.equals(a, b),
        talentsCopy: (a) => RogueTalents.clone(a),
        talentsToJson: (a) => RogueTalents.toJson(a),
        talentsFromJson: (obj) => RogueTalents.fromJson(obj),
        talentsFromPlayer: (player) => player.spec.oneofKind == 'rogue'
            ? player.spec.rogue.talents || RogueTalents.create()
            : RogueTalents.create(),
        optionsCreate: () => RogueOptions.create(),
        optionsEquals: (a, b) => RogueOptions.equals(a, b),
        optionsCopy: (a) => RogueOptions.clone(a),
        optionsToJson: (a) => RogueOptions.toJson(a),
        optionsFromJson: (obj) => RogueOptions.fromJson(obj),
        optionsFromPlayer: (player) => player.spec.oneofKind == 'rogue'
            ? player.spec.rogue.options || RogueOptions.create()
            : RogueOptions.create(),
    },
    [Spec.SpecShadowPriest]: {
        rotationCreate: () => ShadowPriestRotation.create(),
        rotationEquals: (a, b) => ShadowPriestRotation.equals(a, b),
        rotationCopy: (a) => ShadowPriestRotation.clone(a),
        rotationToJson: (a) => ShadowPriestRotation.toJson(a),
        rotationFromJson: (obj) => ShadowPriestRotation.fromJson(obj),
        rotationFromPlayer: (player) => player.spec.oneofKind == 'shadowPriest'
            ? player.spec.shadowPriest.rotation || ShadowPriestRotation.create()
            : ShadowPriestRotation.create(),
        talentsCreate: () => PriestTalents.create(),
        talentsEquals: (a, b) => PriestTalents.equals(a, b),
        talentsCopy: (a) => PriestTalents.clone(a),
        talentsToJson: (a) => PriestTalents.toJson(a),
        talentsFromJson: (obj) => PriestTalents.fromJson(obj),
        talentsFromPlayer: (player) => player.spec.oneofKind == 'shadowPriest'
            ? player.spec.shadowPriest.talents || PriestTalents.create()
            : PriestTalents.create(),
        optionsCreate: () => ShadowPriestOptions.create(),
        optionsEquals: (a, b) => ShadowPriestOptions.equals(a, b),
        optionsCopy: (a) => ShadowPriestOptions.clone(a),
        optionsToJson: (a) => ShadowPriestOptions.toJson(a),
        optionsFromJson: (obj) => ShadowPriestOptions.fromJson(obj),
        optionsFromPlayer: (player) => player.spec.oneofKind == 'shadowPriest'
            ? player.spec.shadowPriest.options || ShadowPriestOptions.create()
            : ShadowPriestOptions.create(),
    },
    [Spec.SpecWarlock]: {
        rotationCreate: () => WarlockRotation.create(),
        rotationEquals: (a, b) => WarlockRotation.equals(a, b),
        rotationCopy: (a) => WarlockRotation.clone(a),
        rotationToJson: (a) => WarlockRotation.toJson(a),
        rotationFromJson: (obj) => WarlockRotation.fromJson(obj),
        rotationFromPlayer: (player) => player.spec.oneofKind == 'warlock'
            ? player.spec.warlock.rotation || WarlockRotation.create()
            : WarlockRotation.create(),
        talentsCreate: () => WarlockTalents.create(),
        talentsEquals: (a, b) => WarlockTalents.equals(a, b),
        talentsCopy: (a) => WarlockTalents.clone(a),
        talentsToJson: (a) => WarlockTalents.toJson(a),
        talentsFromJson: (obj) => WarlockTalents.fromJson(obj),
        talentsFromPlayer: (player) => player.spec.oneofKind == 'warlock'
            ? player.spec.warlock.talents || WarlockTalents.create()
            : WarlockTalents.create(),
        optionsCreate: () => WarlockOptions.create(),
        optionsEquals: (a, b) => WarlockOptions.equals(a, b),
        optionsCopy: (a) => WarlockOptions.clone(a),
        optionsToJson: (a) => WarlockOptions.toJson(a),
        optionsFromJson: (obj) => WarlockOptions.fromJson(obj),
        optionsFromPlayer: (player) => player.spec.oneofKind == 'warlock'
            ? player.spec.warlock.options || WarlockOptions.create()
            : WarlockOptions.create(),
    },
    [Spec.SpecWarrior]: {
        rotationCreate: () => WarriorRotation.create(),
        rotationEquals: (a, b) => WarriorRotation.equals(a, b),
        rotationCopy: (a) => WarriorRotation.clone(a),
        rotationToJson: (a) => WarriorRotation.toJson(a),
        rotationFromJson: (obj) => WarriorRotation.fromJson(obj),
        rotationFromPlayer: (player) => player.spec.oneofKind == 'warrior'
            ? player.spec.warrior.rotation || WarriorRotation.create()
            : WarriorRotation.create(),
        talentsCreate: () => WarriorTalents.create(),
        talentsEquals: (a, b) => WarriorTalents.equals(a, b),
        talentsCopy: (a) => WarriorTalents.clone(a),
        talentsToJson: (a) => WarriorTalents.toJson(a),
        talentsFromJson: (obj) => WarriorTalents.fromJson(obj),
        talentsFromPlayer: (player) => player.spec.oneofKind == 'warrior'
            ? player.spec.warrior.talents || WarriorTalents.create()
            : WarriorTalents.create(),
        optionsCreate: () => WarriorOptions.create(),
        optionsEquals: (a, b) => WarriorOptions.equals(a, b),
        optionsCopy: (a) => WarriorOptions.clone(a),
        optionsToJson: (a) => WarriorOptions.toJson(a),
        optionsFromJson: (obj) => WarriorOptions.fromJson(obj),
        optionsFromPlayer: (player) => player.spec.oneofKind == 'warrior'
            ? player.spec.warrior.options || WarriorOptions.create()
            : WarriorOptions.create(),
    },
    [Spec.SpecProtectionWarrior]: {
        rotationCreate: () => ProtectionWarriorRotation.create(),
        rotationEquals: (a, b) => ProtectionWarriorRotation.equals(a, b),
        rotationCopy: (a) => ProtectionWarriorRotation.clone(a),
        rotationToJson: (a) => ProtectionWarriorRotation.toJson(a),
        rotationFromJson: (obj) => ProtectionWarriorRotation.fromJson(obj),
        rotationFromPlayer: (player) => player.spec.oneofKind == 'protectionWarrior'
            ? player.spec.protectionWarrior.rotation || ProtectionWarriorRotation.create()
            : ProtectionWarriorRotation.create(),
        talentsCreate: () => WarriorTalents.create(),
        talentsEquals: (a, b) => WarriorTalents.equals(a, b),
        talentsCopy: (a) => WarriorTalents.clone(a),
        talentsToJson: (a) => WarriorTalents.toJson(a),
        talentsFromJson: (obj) => WarriorTalents.fromJson(obj),
        talentsFromPlayer: (player) => player.spec.oneofKind == 'protectionWarrior'
            ? player.spec.protectionWarrior.talents || WarriorTalents.create()
            : WarriorTalents.create(),
        optionsCreate: () => ProtectionWarriorOptions.create(),
        optionsEquals: (a, b) => ProtectionWarriorOptions.equals(a, b),
        optionsCopy: (a) => ProtectionWarriorOptions.clone(a),
        optionsToJson: (a) => ProtectionWarriorOptions.toJson(a),
        optionsFromJson: (obj) => ProtectionWarriorOptions.fromJson(obj),
        optionsFromPlayer: (player) => player.spec.oneofKind == 'protectionWarrior'
            ? player.spec.protectionWarrior.options || ProtectionWarriorOptions.create()
            : ProtectionWarriorOptions.create(),
    },
    [Spec.SpecSmitePriest]: {
        rotationCreate: () => SmitePriestRotation.create(),
        rotationEquals: (a, b) => SmitePriestRotation.equals(a, b),
        rotationCopy: (a) => SmitePriestRotation.clone(a),
        rotationToJson: (a) => SmitePriestRotation.toJson(a),
        rotationFromJson: (obj) => SmitePriestRotation.fromJson(obj),
        rotationFromPlayer: (player) => player.spec.oneofKind == 'smitePriest'
            ? player.spec.smitePriest.rotation || SmitePriestRotation.create()
            : SmitePriestRotation.create(),
        talentsCreate: () => PriestTalents.create(),
        talentsEquals: (a, b) => PriestTalents.equals(a, b),
        talentsCopy: (a) => PriestTalents.clone(a),
        talentsToJson: (a) => PriestTalents.toJson(a),
        talentsFromJson: (obj) => PriestTalents.fromJson(obj),
        talentsFromPlayer: (player) => player.spec.oneofKind == 'smitePriest'
            ? player.spec.smitePriest.talents || PriestTalents.create()
            : PriestTalents.create(),
        optionsCreate: () => SmitePriestOptions.create(),
        optionsEquals: (a, b) => SmitePriestOptions.equals(a, b),
        optionsCopy: (a) => SmitePriestOptions.clone(a),
        optionsToJson: (a) => SmitePriestOptions.toJson(a),
        optionsFromJson: (obj) => SmitePriestOptions.fromJson(obj),
        optionsFromPlayer: (player) => player.spec.oneofKind == 'smitePriest'
            ? player.spec.smitePriest.options || SmitePriestOptions.create()
            : SmitePriestOptions.create(),
    },
    [Spec.SpecDeathKnight]: {
        rotationCreate: () => DeathKnightRotation.create(),
        rotationEquals: (a, b) => DeathKnightRotation.equals(a, b),
        rotationCopy: (a) => DeathKnightRotation.clone(a),
        rotationToJson: (a) => DeathKnightRotation.toJson(a),
        rotationFromJson: (obj) => DeathKnightRotation.fromJson(obj),
        rotationFromPlayer: (player) => player.spec.oneofKind == 'deathKnight'
            ? player.spec.deathKnight.rotation || DeathKnightRotation.create()
            : DeathKnightRotation.create(),
        talentsCreate: () => DeathKnightTalents.create(),
        talentsEquals: (a, b) => DeathKnightTalents.equals(a, b),
        talentsCopy: (a) => DeathKnightTalents.clone(a),
        talentsToJson: (a) => DeathKnightTalents.toJson(a),
        talentsFromJson: (obj) => DeathKnightTalents.fromJson(obj),
        talentsFromPlayer: (player) => player.spec.oneofKind == 'deathKnight'
            ? player.spec.deathKnight.talents || DeathKnightTalents.create()
            : DeathKnightTalents.create(),
        optionsCreate: () => DeathKnightOptions.create(),
        optionsEquals: (a, b) => DeathKnightOptions.equals(a, b),
        optionsCopy: (a) => DeathKnightOptions.clone(a),
        optionsToJson: (a) => DeathKnightOptions.toJson(a),
        optionsFromJson: (obj) => DeathKnightOptions.fromJson(obj),
        optionsFromPlayer: (player) => player.spec.oneofKind == 'deathKnight'
            ? player.spec.deathKnight.options || DeathKnightOptions.create()
            : DeathKnightOptions.create(),
    },
    [Spec.SpecDeathKnightTank]: {
        rotationCreate: () => DeathKnightTankRotation.create(),
        rotationEquals: (a, b) => DeathKnightTankRotation.equals(a, b),
        rotationCopy: (a) => DeathKnightTankRotation.clone(a),
        rotationToJson: (a) => DeathKnightTankRotation.toJson(a),
        rotationFromJson: (obj) => DeathKnightTankRotation.fromJson(obj),
        rotationFromPlayer: (player) => player.spec.oneofKind == 'deathKnightTank'
            ? player.spec.deathKnightTank.rotation || DeathKnightTankRotation.create()
            : DeathKnightTankRotation.create(),
        talentsCreate: () => DeathKnightTalents.create(),
        talentsEquals: (a, b) => DeathKnightTalents.equals(a, b),
        talentsCopy: (a) => DeathKnightTalents.clone(a),
        talentsToJson: (a) => DeathKnightTalents.toJson(a),
        talentsFromJson: (obj) => DeathKnightTalents.fromJson(obj),
        talentsFromPlayer: (player) => player.spec.oneofKind == 'deathKnightTank'
            ? player.spec.deathKnightTank.talents || DeathKnightTalents.create()
            : DeathKnightTalents.create(),
        optionsCreate: () => DeathKnightTankOptions.create(),
        optionsEquals: (a, b) => DeathKnightTankOptions.equals(a, b),
        optionsCopy: (a) => DeathKnightTankOptions.clone(a),
        optionsToJson: (a) => DeathKnightTankOptions.toJson(a),
        optionsFromJson: (obj) => DeathKnightTankOptions.fromJson(obj),
        optionsFromPlayer: (player) => player.spec.oneofKind == 'deathKnightTank'
            ? player.spec.deathKnightTank.options || DeathKnightTankOptions.create()
            : DeathKnightTankOptions.create(),
    },
};
export const raceToFaction = {
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
export const specToClass = {
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
    [Spec.SpecDeathKnight]: Class.ClassDeathKnight,
    [Spec.SpecDeathKnightTank]: Class.ClassDeathKnight,
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
    Race.RaceTroll,
];
const mageRaces = [
    Race.RaceBloodElf,
    Race.RaceDraenei,
    Race.RaceGnome,
    Race.RaceHuman,
    Race.RaceTroll,
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
    Race.RaceTroll,
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
    Race.RaceDraenei,
    Race.RaceOrc,
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
const deathKnightRaces = [
    Race.RaceBloodElf,
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
export const specToEligibleRaces = {
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
    [Spec.SpecDeathKnight]: deathKnightRaces,
    [Spec.SpecDeathKnightTank]: deathKnightRaces,
};
// Specs that can dual wield. This could be based on class, except that
// Enhancement Shaman learn dual wield from a talent.
const dualWieldSpecs = [
    Spec.SpecEnhancementShaman,
    Spec.SpecHunter,
    Spec.SpecRogue,
    Spec.SpecWarrior,
    Spec.SpecProtectionWarrior,
    Spec.SpecDeathKnight,
    Spec.SpecDeathKnightTank,
];
export function isDualWieldSpec(spec) {
    return dualWieldSpecs.includes(spec);
}
const tankSpecs = [
    Spec.SpecFeralTankDruid,
    Spec.SpecProtectionPaladin,
    Spec.SpecProtectionWarrior,
];
export function isTankSpec(spec) {
    return tankSpecs.includes(spec);
}
// Prefixes used for storing browser data for each site. Even if a Spec is
// renamed, DO NOT change these values or people will lose their saved data.
export const specToLocalStorageKey = {
    [Spec.SpecBalanceDruid]: '__wotlk_balance_druid',
    [Spec.SpecElementalShaman]: '__wotlk_elemental_shaman',
    [Spec.SpecEnhancementShaman]: '__wotlk_enhacement_shaman',
    [Spec.SpecFeralDruid]: '__wotlk_feral_druid',
    [Spec.SpecFeralTankDruid]: '__wotlk_feral_tank_druid',
    [Spec.SpecHunter]: '__wotlk_hunter',
    [Spec.SpecMage]: '__wotlk_mage',
    [Spec.SpecRetributionPaladin]: '__wotlk_retribution_paladin',
    [Spec.SpecProtectionPaladin]: '__wotlk_protection_paladin',
    [Spec.SpecRogue]: '__wotlk_rogue',
    [Spec.SpecShadowPriest]: '__wotlk_shadow_priest',
    [Spec.SpecWarlock]: '__wotlk_warlock',
    [Spec.SpecWarrior]: '__wotlk_warrior',
    [Spec.SpecProtectionWarrior]: '__wotlk_protection_warrior',
    [Spec.SpecSmitePriest]: '__wotlk_smite_priest',
    [Spec.SpecDeathKnight]: '__wotlk_death_knight',
    [Spec.SpecDeathKnightTank]: '__wotlk_death_knight_tank',
};
// Returns a copy of playerOptions, with the class field set.
export function withSpecProto(spec, player, rotation, talents, specOptions) {
    const copy = Player.clone(player);
    switch (spec) {
        case Spec.SpecBalanceDruid:
            copy.spec = {
                oneofKind: 'balanceDruid',
                balanceDruid: BalanceDruid.create({
                    rotation: rotation,
                    talents: talents,
                    options: specOptions,
                }),
            };
            return copy;
        case Spec.SpecElementalShaman:
            copy.spec = {
                oneofKind: 'elementalShaman',
                elementalShaman: ElementalShaman.create({
                    rotation: rotation,
                    talents: talents,
                    options: specOptions,
                }),
            };
            return copy;
        case Spec.SpecEnhancementShaman:
            copy.spec = {
                oneofKind: 'enhancementShaman',
                enhancementShaman: EnhancementShaman.create({
                    rotation: rotation,
                    talents: talents,
                    options: specOptions,
                }),
            };
            return copy;
        case Spec.SpecFeralDruid:
            copy.spec = {
                oneofKind: 'feralDruid',
                feralDruid: FeralDruid.create({
                    rotation: rotation,
                    talents: talents,
                    options: specOptions,
                }),
            };
            return copy;
        case Spec.SpecFeralTankDruid:
            copy.spec = {
                oneofKind: 'feralTankDruid',
                feralTankDruid: FeralTankDruid.create({
                    rotation: rotation,
                    talents: talents,
                    options: specOptions,
                }),
            };
            return copy;
        case Spec.SpecHunter:
            copy.spec = {
                oneofKind: 'hunter',
                hunter: Hunter.create({
                    rotation: rotation,
                    talents: talents,
                    options: specOptions,
                }),
            };
            return copy;
        case Spec.SpecMage:
            copy.spec = {
                oneofKind: 'mage',
                mage: Mage.create({
                    rotation: rotation,
                    talents: talents,
                    options: specOptions,
                }),
            };
            return copy;
        case Spec.SpecRetributionPaladin:
            copy.spec = {
                oneofKind: 'retributionPaladin',
                retributionPaladin: RetributionPaladin.create({
                    rotation: rotation,
                    talents: talents,
                    options: specOptions,
                }),
            };
            return copy;
        case Spec.SpecProtectionPaladin:
            copy.spec = {
                oneofKind: 'protectionPaladin',
                protectionPaladin: ProtectionPaladin.create({
                    rotation: rotation,
                    talents: talents,
                    options: specOptions,
                }),
            };
            return copy;
        case Spec.SpecRogue:
            copy.spec = {
                oneofKind: 'rogue',
                rogue: Rogue.create({
                    rotation: rotation,
                    talents: talents,
                    options: specOptions,
                }),
            };
            return copy;
        case Spec.SpecShadowPriest:
            copy.spec = {
                oneofKind: 'shadowPriest',
                shadowPriest: ShadowPriest.create({
                    rotation: rotation,
                    talents: talents,
                    options: specOptions,
                }),
            };
            return copy;
        case Spec.SpecWarlock:
            copy.spec = {
                oneofKind: 'warlock',
                warlock: Warlock.create({
                    rotation: rotation,
                    talents: talents,
                    options: specOptions,
                }),
            };
            return copy;
        case Spec.SpecWarrior:
            copy.spec = {
                oneofKind: 'warrior',
                warrior: Warrior.create({
                    rotation: rotation,
                    talents: talents,
                    options: specOptions,
                }),
            };
            return copy;
        case Spec.SpecProtectionWarrior:
            copy.spec = {
                oneofKind: 'protectionWarrior',
                protectionWarrior: ProtectionWarrior.create({
                    rotation: rotation,
                    talents: talents,
                    options: specOptions,
                }),
            };
            return copy;
        case Spec.SpecSmitePriest:
            copy.spec = {
                oneofKind: 'smitePriest',
                smitePriest: SmitePriest.create({
                    rotation: rotation,
                    talents: talents,
                    options: specOptions,
                }),
            };
            return copy;
        case Spec.SpecDeathKnight:
            copy.spec = {
                oneofKind: 'deathKnight',
                deathKnight: DeathKnight.create({
                    rotation: rotation,
                    talents: talents,
                    options: specOptions,
                }),
            };
            return copy;
        case Spec.SpecDeathKnightTank:
            copy.spec = {
                oneofKind: 'deathKnightTank',
                deathKnightTank: DeathKnightTank.create({
                    rotation: rotation,
                    talents: talents,
                    options: specOptions,
                }),
            };
            return copy;
    }
}
export function playerToSpec(player) {
    const specValues = getEnumValues(Spec);
    for (let i = 0; i < specValues.length; i++) {
        const spec = specValues[i];
        let specString = Spec[spec]; // Returns 'SpecBalanceDruid' for BalanceDruid.
        specString = specString.substring('Spec'.length); // 'BalanceDruid'
        specString = specString.charAt(0).toLowerCase() + specString.slice(1); // 'balanceDruid'
        if (player.spec.oneofKind == specString) {
            return spec;
        }
    }
    throw new Error('Unable to parse spec from player proto: ' + JSON.stringify(Player.toJson(player), null, 2));
}
const classToMaxArmorType = {
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
    [Class.ClassDeathKnight]: ArmorType.ArmorTypePlate,
};
const classToEligibleRangedWeaponTypes = {
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
    [Class.ClassDeathKnight]: [
        RangedWeaponType.RangedWeaponTypeSigil,
    ],
};
const classToEligibleWeaponTypes = {
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
    [Class.ClassDeathKnight]: [
        { weaponType: WeaponType.WeaponTypeAxe, canUseTwoHand: true },
        { weaponType: WeaponType.WeaponTypeMace, canUseTwoHand: true },
        { weaponType: WeaponType.WeaponTypePolearm, canUseTwoHand: true },
        { weaponType: WeaponType.WeaponTypeSword, canUseTwoHand: true },
        // TODO: validate proficiencies
    ],
};
export function isSharpWeaponType(weaponType) {
    return [
        WeaponType.WeaponTypeAxe,
        WeaponType.WeaponTypeDagger,
        WeaponType.WeaponTypePolearm,
        WeaponType.WeaponTypeSword,
    ].includes(weaponType);
}
export function isBluntWeaponType(weaponType) {
    return [
        WeaponType.WeaponTypeFist,
        WeaponType.WeaponTypeMace,
        WeaponType.WeaponTypeStaff,
    ].includes(weaponType);
}
// Custom functions for determining the EP value of meta gem effects.
// Default meta effect EP value is 0, so just handle the ones relevant to your spec.
const metaGemEffectEPs = {
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
export function getMetaGemEffectEP(spec, gem, playerStats) {
    if (metaGemEffectEPs[spec]) {
        return metaGemEffectEPs[spec](gem, playerStats);
    }
    else {
        return 0;
    }
}
// Returns true if this item may be equipped in at least 1 slot for the given Spec.
export function canEquipItem(item, spec, slot) {
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
const itemTypeToSlotsMap = {
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
export function getEligibleItemSlots(item) {
    if (itemTypeToSlotsMap[item.type]) {
        return itemTypeToSlotsMap[item.type];
    }
    if (item.type == ItemType.ItemTypeWeapon) {
        if ([HandType.HandTypeMainHand, HandType.HandTypeTwoHand].includes(item.handType)) {
            return [ItemSlot.ItemSlotMainHand];
        }
        else if (item.handType == HandType.HandTypeOffHand) {
            return [ItemSlot.ItemSlotOffHand];
        }
        else {
            return [ItemSlot.ItemSlotMainHand, ItemSlot.ItemSlotOffHand];
        }
    }
    // Should never reach here
    throw new Error('Could not find item slots for item: ' + Item.toJsonString(item));
}
;
// Returns whether the given main-hand and off-hand items can be worn at the
// same time.
export function validWeaponCombo(mainHand, offHand) {
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
export function getEligibleEnchantSlots(enchant) {
    if (itemTypeToSlotsMap[enchant.type]) {
        return itemTypeToSlotsMap[enchant.type];
    }
    if (enchant.type == ItemType.ItemTypeWeapon) {
        return [ItemSlot.ItemSlotMainHand, ItemSlot.ItemSlotOffHand];
    }
    // Should never reach here
    throw new Error('Could not find item slots for enchant: ' + Enchant.toJsonString(enchant));
}
;
export function enchantAppliesToItem(enchant, item) {
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
}
;
export function canEquipEnchant(enchant, spec) {
    const playerClass = specToClass[spec];
    if (enchant.classAllowlist.length > 0 && !enchant.classAllowlist.includes(playerClass)) {
        return false;
    }
    return true;
}
export const NO_TARGET = -1;
export function newRaidTarget(raidIndex) {
    return RaidTarget.create({
        targetIndex: raidIndex,
    });
}
export function emptyRaidTarget() {
    return newRaidTarget(NO_TARGET);
}
// Makes a new set of assignments with everything 0'd out.
export function makeBlankBlessingsAssignments(numPaladins) {
    const assignments = BlessingsAssignments.create();
    for (let i = 0; i < numPaladins; i++) {
        assignments.paladins.push(BlessingsAssignment.create({
            blessings: new Array(NUM_SPECS).fill(Blessings.BlessingUnknown),
        }));
    }
    return assignments;
}
export function makeBlessingsAssignments(numPaladins, data) {
    const assignments = makeBlankBlessingsAssignments(numPaladins);
    for (let i = 0; i < data.length; i++) {
        const spec = data[i].spec;
        const blessings = data[i].blessings;
        for (let j = 0; j < blessings.length; j++) {
            if (j >= assignments.paladins.length) {
                // Can't assign more blessings since we ran out of paladins
                break;
            }
            assignments.paladins[j].blessings[spec] = blessings[j];
        }
    }
    return assignments;
}
// Default blessings settings in the raid sim UI.
export function makeDefaultBlessings(numPaladins) {
    return makeBlessingsAssignments(numPaladins, [
        { spec: Spec.SpecBalanceDruid, blessings: [Blessings.BlessingOfKings, Blessings.BlessingOfWisdom] },
        { spec: Spec.SpecFeralDruid, blessings: [Blessings.BlessingOfKings, Blessings.BlessingOfMight, Blessings.BlessingOfWisdom] },
        { spec: Spec.SpecFeralTankDruid, blessings: [Blessings.BlessingOfKings, Blessings.BlessingOfMight, Blessings.BlessingOfSanctuary] },
        { spec: Spec.SpecHunter, blessings: [Blessings.BlessingOfKings, Blessings.BlessingOfMight, Blessings.BlessingOfWisdom] },
        { spec: Spec.SpecMage, blessings: [Blessings.BlessingOfKings, Blessings.BlessingOfWisdom] },
        { spec: Spec.SpecRetributionPaladin, blessings: [Blessings.BlessingOfKings, Blessings.BlessingOfMight, Blessings.BlessingOfWisdom] },
        { spec: Spec.SpecProtectionPaladin, blessings: [Blessings.BlessingOfKings, Blessings.BlessingOfSanctuary, Blessings.BlessingOfWisdom, Blessings.BlessingOfMight] },
        { spec: Spec.SpecShadowPriest, blessings: [Blessings.BlessingOfKings, Blessings.BlessingOfWisdom] },
        { spec: Spec.SpecSmitePriest, blessings: [Blessings.BlessingOfKings, Blessings.BlessingOfWisdom] },
        { spec: Spec.SpecRogue, blessings: [Blessings.BlessingOfKings, Blessings.BlessingOfMight] },
        { spec: Spec.SpecElementalShaman, blessings: [Blessings.BlessingOfKings, Blessings.BlessingOfWisdom] },
        { spec: Spec.SpecEnhancementShaman, blessings: [Blessings.BlessingOfKings, Blessings.BlessingOfMight, Blessings.BlessingOfWisdom] },
        { spec: Spec.SpecWarlock, blessings: [Blessings.BlessingOfKings, Blessings.BlessingOfWisdom] },
        { spec: Spec.SpecWarrior, blessings: [Blessings.BlessingOfKings, Blessings.BlessingOfMight] },
        { spec: Spec.SpecProtectionWarrior, blessings: [Blessings.BlessingOfKings, Blessings.BlessingOfMight, Blessings.BlessingOfSanctuary] },
        { spec: Spec.SpecDeathKnight, blessings: [Blessings.BlessingOfKings, Blessings.BlessingOfMight, Blessings.BlessingOfSalvation] },
        { spec: Spec.SpecDeathKnightTank, blessings: [Blessings.BlessingOfKings, Blessings.BlessingOfMight] },
    ]);
}
;
export const orderedResourceTypes = [
    ResourceType.ResourceTypeHealth,
    ResourceType.ResourceTypeMana,
    ResourceType.ResourceTypeEnergy,
    ResourceType.ResourceTypeRage,
    ResourceType.ResourceTypeComboPoints,
    ResourceType.ResourceTypeFocus,
    ResourceType.ResourceTypeRunicPower,
    ResourceType.ResourceTypeBloodRune,
    ResourceType.ResourceTypeFrostRune,
    ResourceType.ResourceTypeUnholyRune,
    ResourceType.ResourceTypeDeathRune,
];
