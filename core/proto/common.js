import { WireType } from '/wotlk/protobuf-ts/index.js';
import { UnknownFieldHandler } from '/wotlk/protobuf-ts/index.js';
import { reflectionMergePartial } from '/wotlk/protobuf-ts/index.js';
import { MESSAGE_TYPE } from '/wotlk/protobuf-ts/index.js';
import { MessageType } from '/wotlk/protobuf-ts/index.js';
/**
 * @generated from protobuf enum proto.Spec
 */
export var Spec;
(function (Spec) {
    /**
     * @generated from protobuf enum value: SpecBalanceDruid = 0;
     */
    Spec[Spec["SpecBalanceDruid"] = 0] = "SpecBalanceDruid";
    /**
     * @generated from protobuf enum value: SpecElementalShaman = 1;
     */
    Spec[Spec["SpecElementalShaman"] = 1] = "SpecElementalShaman";
    /**
     * @generated from protobuf enum value: SpecEnhancementShaman = 9;
     */
    Spec[Spec["SpecEnhancementShaman"] = 9] = "SpecEnhancementShaman";
    /**
     * @generated from protobuf enum value: SpecFeralDruid = 12;
     */
    Spec[Spec["SpecFeralDruid"] = 12] = "SpecFeralDruid";
    /**
     * @generated from protobuf enum value: SpecFeralTankDruid = 14;
     */
    Spec[Spec["SpecFeralTankDruid"] = 14] = "SpecFeralTankDruid";
    /**
     * @generated from protobuf enum value: SpecHunter = 8;
     */
    Spec[Spec["SpecHunter"] = 8] = "SpecHunter";
    /**
     * @generated from protobuf enum value: SpecMage = 2;
     */
    Spec[Spec["SpecMage"] = 2] = "SpecMage";
    /**
     * @generated from protobuf enum value: SpecProtectionPaladin = 13;
     */
    Spec[Spec["SpecProtectionPaladin"] = 13] = "SpecProtectionPaladin";
    /**
     * @generated from protobuf enum value: SpecRetributionPaladin = 3;
     */
    Spec[Spec["SpecRetributionPaladin"] = 3] = "SpecRetributionPaladin";
    /**
     * @generated from protobuf enum value: SpecRogue = 7;
     */
    Spec[Spec["SpecRogue"] = 7] = "SpecRogue";
    /**
     * @generated from protobuf enum value: SpecShadowPriest = 4;
     */
    Spec[Spec["SpecShadowPriest"] = 4] = "SpecShadowPriest";
    /**
     * @generated from protobuf enum value: SpecSmitePriest = 10;
     */
    Spec[Spec["SpecSmitePriest"] = 10] = "SpecSmitePriest";
    /**
     * @generated from protobuf enum value: SpecWarlock = 5;
     */
    Spec[Spec["SpecWarlock"] = 5] = "SpecWarlock";
    /**
     * @generated from protobuf enum value: SpecWarrior = 6;
     */
    Spec[Spec["SpecWarrior"] = 6] = "SpecWarrior";
    /**
     * @generated from protobuf enum value: SpecProtectionWarrior = 11;
     */
    Spec[Spec["SpecProtectionWarrior"] = 11] = "SpecProtectionWarrior";
    /**
     * @generated from protobuf enum value: SpecDeathKnight = 15;
     */
    Spec[Spec["SpecDeathKnight"] = 15] = "SpecDeathKnight";
    /**
     * @generated from protobuf enum value: SpecDeathKnightTank = 16;
     */
    Spec[Spec["SpecDeathKnightTank"] = 16] = "SpecDeathKnightTank";
})(Spec || (Spec = {}));
/**
 * @generated from protobuf enum proto.Race
 */
export var Race;
(function (Race) {
    /**
     * @generated from protobuf enum value: RaceUnknown = 0;
     */
    Race[Race["RaceUnknown"] = 0] = "RaceUnknown";
    /**
     * @generated from protobuf enum value: RaceBloodElf = 1;
     */
    Race[Race["RaceBloodElf"] = 1] = "RaceBloodElf";
    /**
     * @generated from protobuf enum value: RaceDraenei = 2;
     */
    Race[Race["RaceDraenei"] = 2] = "RaceDraenei";
    /**
     * @generated from protobuf enum value: RaceDwarf = 3;
     */
    Race[Race["RaceDwarf"] = 3] = "RaceDwarf";
    /**
     * @generated from protobuf enum value: RaceGnome = 4;
     */
    Race[Race["RaceGnome"] = 4] = "RaceGnome";
    /**
     * @generated from protobuf enum value: RaceHuman = 5;
     */
    Race[Race["RaceHuman"] = 5] = "RaceHuman";
    /**
     * @generated from protobuf enum value: RaceNightElf = 6;
     */
    Race[Race["RaceNightElf"] = 6] = "RaceNightElf";
    /**
     * @generated from protobuf enum value: RaceOrc = 7;
     */
    Race[Race["RaceOrc"] = 7] = "RaceOrc";
    /**
     * @generated from protobuf enum value: RaceTauren = 8;
     */
    Race[Race["RaceTauren"] = 8] = "RaceTauren";
    /**
     * @generated from protobuf enum value: RaceTroll = 9;
     */
    Race[Race["RaceTroll"] = 9] = "RaceTroll";
    /**
     * @generated from protobuf enum value: RaceUndead = 10;
     */
    Race[Race["RaceUndead"] = 10] = "RaceUndead";
})(Race || (Race = {}));
/**
 * @generated from protobuf enum proto.Faction
 */
export var Faction;
(function (Faction) {
    /**
     * @generated from protobuf enum value: Unknown = 0;
     */
    Faction[Faction["Unknown"] = 0] = "Unknown";
    /**
     * @generated from protobuf enum value: Alliance = 1;
     */
    Faction[Faction["Alliance"] = 1] = "Alliance";
    /**
     * @generated from protobuf enum value: Horde = 2;
     */
    Faction[Faction["Horde"] = 2] = "Horde";
})(Faction || (Faction = {}));
/**
 * @generated from protobuf enum proto.ShattrathFaction
 */
export var ShattrathFaction;
(function (ShattrathFaction) {
    /**
     * @generated from protobuf enum value: ShattrathFactionAldor = 0;
     */
    ShattrathFaction[ShattrathFaction["ShattrathFactionAldor"] = 0] = "ShattrathFactionAldor";
    /**
     * @generated from protobuf enum value: ShattrathFactionScryer = 1;
     */
    ShattrathFaction[ShattrathFaction["ShattrathFactionScryer"] = 1] = "ShattrathFactionScryer";
})(ShattrathFaction || (ShattrathFaction = {}));
/**
 * @generated from protobuf enum proto.Class
 */
export var Class;
(function (Class) {
    /**
     * @generated from protobuf enum value: ClassUnknown = 0;
     */
    Class[Class["ClassUnknown"] = 0] = "ClassUnknown";
    /**
     * @generated from protobuf enum value: ClassDruid = 1;
     */
    Class[Class["ClassDruid"] = 1] = "ClassDruid";
    /**
     * @generated from protobuf enum value: ClassHunter = 2;
     */
    Class[Class["ClassHunter"] = 2] = "ClassHunter";
    /**
     * @generated from protobuf enum value: ClassMage = 3;
     */
    Class[Class["ClassMage"] = 3] = "ClassMage";
    /**
     * @generated from protobuf enum value: ClassPaladin = 4;
     */
    Class[Class["ClassPaladin"] = 4] = "ClassPaladin";
    /**
     * @generated from protobuf enum value: ClassPriest = 5;
     */
    Class[Class["ClassPriest"] = 5] = "ClassPriest";
    /**
     * @generated from protobuf enum value: ClassRogue = 6;
     */
    Class[Class["ClassRogue"] = 6] = "ClassRogue";
    /**
     * @generated from protobuf enum value: ClassShaman = 7;
     */
    Class[Class["ClassShaman"] = 7] = "ClassShaman";
    /**
     * @generated from protobuf enum value: ClassWarlock = 8;
     */
    Class[Class["ClassWarlock"] = 8] = "ClassWarlock";
    /**
     * @generated from protobuf enum value: ClassWarrior = 9;
     */
    Class[Class["ClassWarrior"] = 9] = "ClassWarrior";
    /**
     * @generated from protobuf enum value: ClassDeathKnight = 10;
     */
    Class[Class["ClassDeathKnight"] = 10] = "ClassDeathKnight";
})(Class || (Class = {}));
/**
 * @generated from protobuf enum proto.Profession
 */
export var Profession;
(function (Profession) {
    /**
     * @generated from protobuf enum value: ProfessionUnknown = 0;
     */
    Profession[Profession["ProfessionUnknown"] = 0] = "ProfessionUnknown";
    /**
     * @generated from protobuf enum value: Alchemy = 1;
     */
    Profession[Profession["Alchemy"] = 1] = "Alchemy";
    /**
     * @generated from protobuf enum value: Blacksmithing = 2;
     */
    Profession[Profession["Blacksmithing"] = 2] = "Blacksmithing";
    /**
     * @generated from protobuf enum value: Enchanting = 3;
     */
    Profession[Profession["Enchanting"] = 3] = "Enchanting";
    /**
     * @generated from protobuf enum value: Engineering = 4;
     */
    Profession[Profession["Engineering"] = 4] = "Engineering";
    /**
     * @generated from protobuf enum value: Herbalism = 5;
     */
    Profession[Profession["Herbalism"] = 5] = "Herbalism";
    /**
     * @generated from protobuf enum value: Inscription = 6;
     */
    Profession[Profession["Inscription"] = 6] = "Inscription";
    /**
     * @generated from protobuf enum value: Jewelcrafting = 7;
     */
    Profession[Profession["Jewelcrafting"] = 7] = "Jewelcrafting";
    /**
     * @generated from protobuf enum value: Leatherworking = 8;
     */
    Profession[Profession["Leatherworking"] = 8] = "Leatherworking";
    /**
     * @generated from protobuf enum value: Mining = 9;
     */
    Profession[Profession["Mining"] = 9] = "Mining";
    /**
     * @generated from protobuf enum value: Skinning = 10;
     */
    Profession[Profession["Skinning"] = 10] = "Skinning";
    /**
     * @generated from protobuf enum value: Tailoring = 11;
     */
    Profession[Profession["Tailoring"] = 11] = "Tailoring";
})(Profession || (Profession = {}));
/**
 * @generated from protobuf enum proto.Stat
 */
export var Stat;
(function (Stat) {
    /**
     * @generated from protobuf enum value: StatStrength = 0;
     */
    Stat[Stat["StatStrength"] = 0] = "StatStrength";
    /**
     * @generated from protobuf enum value: StatAgility = 1;
     */
    Stat[Stat["StatAgility"] = 1] = "StatAgility";
    /**
     * @generated from protobuf enum value: StatStamina = 2;
     */
    Stat[Stat["StatStamina"] = 2] = "StatStamina";
    /**
     * @generated from protobuf enum value: StatIntellect = 3;
     */
    Stat[Stat["StatIntellect"] = 3] = "StatIntellect";
    /**
     * @generated from protobuf enum value: StatSpirit = 4;
     */
    Stat[Stat["StatSpirit"] = 4] = "StatSpirit";
    /**
     * @generated from protobuf enum value: StatSpellPower = 5;
     */
    Stat[Stat["StatSpellPower"] = 5] = "StatSpellPower";
    /**
     * @generated from protobuf enum value: StatHealingPower = 6;
     */
    Stat[Stat["StatHealingPower"] = 6] = "StatHealingPower";
    /**
     * @generated from protobuf enum value: StatArcaneSpellPower = 7;
     */
    Stat[Stat["StatArcaneSpellPower"] = 7] = "StatArcaneSpellPower";
    /**
     * @generated from protobuf enum value: StatFireSpellPower = 8;
     */
    Stat[Stat["StatFireSpellPower"] = 8] = "StatFireSpellPower";
    /**
     * @generated from protobuf enum value: StatFrostSpellPower = 9;
     */
    Stat[Stat["StatFrostSpellPower"] = 9] = "StatFrostSpellPower";
    /**
     * @generated from protobuf enum value: StatHolySpellPower = 10;
     */
    Stat[Stat["StatHolySpellPower"] = 10] = "StatHolySpellPower";
    /**
     * @generated from protobuf enum value: StatNatureSpellPower = 11;
     */
    Stat[Stat["StatNatureSpellPower"] = 11] = "StatNatureSpellPower";
    /**
     * @generated from protobuf enum value: StatShadowSpellPower = 12;
     */
    Stat[Stat["StatShadowSpellPower"] = 12] = "StatShadowSpellPower";
    /**
     * @generated from protobuf enum value: StatMP5 = 13;
     */
    Stat[Stat["StatMP5"] = 13] = "StatMP5";
    /**
     * @generated from protobuf enum value: StatSpellHit = 14;
     */
    Stat[Stat["StatSpellHit"] = 14] = "StatSpellHit";
    /**
     * @generated from protobuf enum value: StatSpellCrit = 15;
     */
    Stat[Stat["StatSpellCrit"] = 15] = "StatSpellCrit";
    /**
     * @generated from protobuf enum value: StatSpellHaste = 16;
     */
    Stat[Stat["StatSpellHaste"] = 16] = "StatSpellHaste";
    /**
     * @generated from protobuf enum value: StatSpellPenetration = 17;
     */
    Stat[Stat["StatSpellPenetration"] = 17] = "StatSpellPenetration";
    /**
     * @generated from protobuf enum value: StatAttackPower = 18;
     */
    Stat[Stat["StatAttackPower"] = 18] = "StatAttackPower";
    /**
     * @generated from protobuf enum value: StatMeleeHit = 19;
     */
    Stat[Stat["StatMeleeHit"] = 19] = "StatMeleeHit";
    /**
     * @generated from protobuf enum value: StatMeleeCrit = 20;
     */
    Stat[Stat["StatMeleeCrit"] = 20] = "StatMeleeCrit";
    /**
     * @generated from protobuf enum value: StatMeleeHaste = 21;
     */
    Stat[Stat["StatMeleeHaste"] = 21] = "StatMeleeHaste";
    /**
     * @generated from protobuf enum value: StatArmorPenetration = 22;
     */
    Stat[Stat["StatArmorPenetration"] = 22] = "StatArmorPenetration";
    /**
     * @generated from protobuf enum value: StatExpertise = 23;
     */
    Stat[Stat["StatExpertise"] = 23] = "StatExpertise";
    /**
     * @generated from protobuf enum value: StatMana = 24;
     */
    Stat[Stat["StatMana"] = 24] = "StatMana";
    /**
     * @generated from protobuf enum value: StatEnergy = 25;
     */
    Stat[Stat["StatEnergy"] = 25] = "StatEnergy";
    /**
     * @generated from protobuf enum value: StatRage = 26;
     */
    Stat[Stat["StatRage"] = 26] = "StatRage";
    /**
     * @generated from protobuf enum value: StatArmor = 27;
     */
    Stat[Stat["StatArmor"] = 27] = "StatArmor";
    /**
     * @generated from protobuf enum value: StatRangedAttackPower = 28;
     */
    Stat[Stat["StatRangedAttackPower"] = 28] = "StatRangedAttackPower";
    /**
     * @generated from protobuf enum value: StatDefense = 29;
     */
    Stat[Stat["StatDefense"] = 29] = "StatDefense";
    /**
     * @generated from protobuf enum value: StatBlock = 30;
     */
    Stat[Stat["StatBlock"] = 30] = "StatBlock";
    /**
     * @generated from protobuf enum value: StatBlockValue = 31;
     */
    Stat[Stat["StatBlockValue"] = 31] = "StatBlockValue";
    /**
     * @generated from protobuf enum value: StatDodge = 32;
     */
    Stat[Stat["StatDodge"] = 32] = "StatDodge";
    /**
     * @generated from protobuf enum value: StatParry = 33;
     */
    Stat[Stat["StatParry"] = 33] = "StatParry";
    /**
     * @generated from protobuf enum value: StatResilience = 34;
     */
    Stat[Stat["StatResilience"] = 34] = "StatResilience";
    /**
     * @generated from protobuf enum value: StatHealth = 35;
     */
    Stat[Stat["StatHealth"] = 35] = "StatHealth";
    /**
     * @generated from protobuf enum value: StatArcaneResistance = 36;
     */
    Stat[Stat["StatArcaneResistance"] = 36] = "StatArcaneResistance";
    /**
     * @generated from protobuf enum value: StatFireResistance = 37;
     */
    Stat[Stat["StatFireResistance"] = 37] = "StatFireResistance";
    /**
     * @generated from protobuf enum value: StatFrostResistance = 38;
     */
    Stat[Stat["StatFrostResistance"] = 38] = "StatFrostResistance";
    /**
     * @generated from protobuf enum value: StatNatureResistance = 39;
     */
    Stat[Stat["StatNatureResistance"] = 39] = "StatNatureResistance";
    /**
     * @generated from protobuf enum value: StatShadowResistance = 40;
     */
    Stat[Stat["StatShadowResistance"] = 40] = "StatShadowResistance";
})(Stat || (Stat = {}));
/**
 * @generated from protobuf enum proto.ItemType
 */
export var ItemType;
(function (ItemType) {
    /**
     * @generated from protobuf enum value: ItemTypeUnknown = 0;
     */
    ItemType[ItemType["ItemTypeUnknown"] = 0] = "ItemTypeUnknown";
    /**
     * @generated from protobuf enum value: ItemTypeHead = 1;
     */
    ItemType[ItemType["ItemTypeHead"] = 1] = "ItemTypeHead";
    /**
     * @generated from protobuf enum value: ItemTypeNeck = 2;
     */
    ItemType[ItemType["ItemTypeNeck"] = 2] = "ItemTypeNeck";
    /**
     * @generated from protobuf enum value: ItemTypeShoulder = 3;
     */
    ItemType[ItemType["ItemTypeShoulder"] = 3] = "ItemTypeShoulder";
    /**
     * @generated from protobuf enum value: ItemTypeBack = 4;
     */
    ItemType[ItemType["ItemTypeBack"] = 4] = "ItemTypeBack";
    /**
     * @generated from protobuf enum value: ItemTypeChest = 5;
     */
    ItemType[ItemType["ItemTypeChest"] = 5] = "ItemTypeChest";
    /**
     * @generated from protobuf enum value: ItemTypeWrist = 6;
     */
    ItemType[ItemType["ItemTypeWrist"] = 6] = "ItemTypeWrist";
    /**
     * @generated from protobuf enum value: ItemTypeHands = 7;
     */
    ItemType[ItemType["ItemTypeHands"] = 7] = "ItemTypeHands";
    /**
     * @generated from protobuf enum value: ItemTypeWaist = 8;
     */
    ItemType[ItemType["ItemTypeWaist"] = 8] = "ItemTypeWaist";
    /**
     * @generated from protobuf enum value: ItemTypeLegs = 9;
     */
    ItemType[ItemType["ItemTypeLegs"] = 9] = "ItemTypeLegs";
    /**
     * @generated from protobuf enum value: ItemTypeFeet = 10;
     */
    ItemType[ItemType["ItemTypeFeet"] = 10] = "ItemTypeFeet";
    /**
     * @generated from protobuf enum value: ItemTypeFinger = 11;
     */
    ItemType[ItemType["ItemTypeFinger"] = 11] = "ItemTypeFinger";
    /**
     * @generated from protobuf enum value: ItemTypeTrinket = 12;
     */
    ItemType[ItemType["ItemTypeTrinket"] = 12] = "ItemTypeTrinket";
    /**
     * @generated from protobuf enum value: ItemTypeWeapon = 13;
     */
    ItemType[ItemType["ItemTypeWeapon"] = 13] = "ItemTypeWeapon";
    /**
     * @generated from protobuf enum value: ItemTypeRanged = 14;
     */
    ItemType[ItemType["ItemTypeRanged"] = 14] = "ItemTypeRanged";
})(ItemType || (ItemType = {}));
/**
 * @generated from protobuf enum proto.ArmorType
 */
export var ArmorType;
(function (ArmorType) {
    /**
     * @generated from protobuf enum value: ArmorTypeUnknown = 0;
     */
    ArmorType[ArmorType["ArmorTypeUnknown"] = 0] = "ArmorTypeUnknown";
    /**
     * @generated from protobuf enum value: ArmorTypeCloth = 1;
     */
    ArmorType[ArmorType["ArmorTypeCloth"] = 1] = "ArmorTypeCloth";
    /**
     * @generated from protobuf enum value: ArmorTypeLeather = 2;
     */
    ArmorType[ArmorType["ArmorTypeLeather"] = 2] = "ArmorTypeLeather";
    /**
     * @generated from protobuf enum value: ArmorTypeMail = 3;
     */
    ArmorType[ArmorType["ArmorTypeMail"] = 3] = "ArmorTypeMail";
    /**
     * @generated from protobuf enum value: ArmorTypePlate = 4;
     */
    ArmorType[ArmorType["ArmorTypePlate"] = 4] = "ArmorTypePlate";
})(ArmorType || (ArmorType = {}));
/**
 * @generated from protobuf enum proto.WeaponType
 */
export var WeaponType;
(function (WeaponType) {
    /**
     * @generated from protobuf enum value: WeaponTypeUnknown = 0;
     */
    WeaponType[WeaponType["WeaponTypeUnknown"] = 0] = "WeaponTypeUnknown";
    /**
     * @generated from protobuf enum value: WeaponTypeAxe = 1;
     */
    WeaponType[WeaponType["WeaponTypeAxe"] = 1] = "WeaponTypeAxe";
    /**
     * @generated from protobuf enum value: WeaponTypeDagger = 2;
     */
    WeaponType[WeaponType["WeaponTypeDagger"] = 2] = "WeaponTypeDagger";
    /**
     * @generated from protobuf enum value: WeaponTypeFist = 3;
     */
    WeaponType[WeaponType["WeaponTypeFist"] = 3] = "WeaponTypeFist";
    /**
     * @generated from protobuf enum value: WeaponTypeMace = 4;
     */
    WeaponType[WeaponType["WeaponTypeMace"] = 4] = "WeaponTypeMace";
    /**
     * @generated from protobuf enum value: WeaponTypeOffHand = 5;
     */
    WeaponType[WeaponType["WeaponTypeOffHand"] = 5] = "WeaponTypeOffHand";
    /**
     * @generated from protobuf enum value: WeaponTypePolearm = 6;
     */
    WeaponType[WeaponType["WeaponTypePolearm"] = 6] = "WeaponTypePolearm";
    /**
     * @generated from protobuf enum value: WeaponTypeShield = 7;
     */
    WeaponType[WeaponType["WeaponTypeShield"] = 7] = "WeaponTypeShield";
    /**
     * @generated from protobuf enum value: WeaponTypeStaff = 8;
     */
    WeaponType[WeaponType["WeaponTypeStaff"] = 8] = "WeaponTypeStaff";
    /**
     * @generated from protobuf enum value: WeaponTypeSword = 9;
     */
    WeaponType[WeaponType["WeaponTypeSword"] = 9] = "WeaponTypeSword";
})(WeaponType || (WeaponType = {}));
/**
 * @generated from protobuf enum proto.HandType
 */
export var HandType;
(function (HandType) {
    /**
     * @generated from protobuf enum value: HandTypeUnknown = 0;
     */
    HandType[HandType["HandTypeUnknown"] = 0] = "HandTypeUnknown";
    /**
     * @generated from protobuf enum value: HandTypeMainHand = 1;
     */
    HandType[HandType["HandTypeMainHand"] = 1] = "HandTypeMainHand";
    /**
     * @generated from protobuf enum value: HandTypeOneHand = 2;
     */
    HandType[HandType["HandTypeOneHand"] = 2] = "HandTypeOneHand";
    /**
     * @generated from protobuf enum value: HandTypeOffHand = 3;
     */
    HandType[HandType["HandTypeOffHand"] = 3] = "HandTypeOffHand";
    /**
     * @generated from protobuf enum value: HandTypeTwoHand = 4;
     */
    HandType[HandType["HandTypeTwoHand"] = 4] = "HandTypeTwoHand";
})(HandType || (HandType = {}));
/**
 * @generated from protobuf enum proto.RangedWeaponType
 */
export var RangedWeaponType;
(function (RangedWeaponType) {
    /**
     * @generated from protobuf enum value: RangedWeaponTypeUnknown = 0;
     */
    RangedWeaponType[RangedWeaponType["RangedWeaponTypeUnknown"] = 0] = "RangedWeaponTypeUnknown";
    /**
     * @generated from protobuf enum value: RangedWeaponTypeBow = 1;
     */
    RangedWeaponType[RangedWeaponType["RangedWeaponTypeBow"] = 1] = "RangedWeaponTypeBow";
    /**
     * @generated from protobuf enum value: RangedWeaponTypeCrossbow = 2;
     */
    RangedWeaponType[RangedWeaponType["RangedWeaponTypeCrossbow"] = 2] = "RangedWeaponTypeCrossbow";
    /**
     * @generated from protobuf enum value: RangedWeaponTypeGun = 3;
     */
    RangedWeaponType[RangedWeaponType["RangedWeaponTypeGun"] = 3] = "RangedWeaponTypeGun";
    /**
     * @generated from protobuf enum value: RangedWeaponTypeIdol = 4;
     */
    RangedWeaponType[RangedWeaponType["RangedWeaponTypeIdol"] = 4] = "RangedWeaponTypeIdol";
    /**
     * @generated from protobuf enum value: RangedWeaponTypeLibram = 5;
     */
    RangedWeaponType[RangedWeaponType["RangedWeaponTypeLibram"] = 5] = "RangedWeaponTypeLibram";
    /**
     * @generated from protobuf enum value: RangedWeaponTypeThrown = 6;
     */
    RangedWeaponType[RangedWeaponType["RangedWeaponTypeThrown"] = 6] = "RangedWeaponTypeThrown";
    /**
     * @generated from protobuf enum value: RangedWeaponTypeTotem = 7;
     */
    RangedWeaponType[RangedWeaponType["RangedWeaponTypeTotem"] = 7] = "RangedWeaponTypeTotem";
    /**
     * @generated from protobuf enum value: RangedWeaponTypeWand = 8;
     */
    RangedWeaponType[RangedWeaponType["RangedWeaponTypeWand"] = 8] = "RangedWeaponTypeWand";
    /**
     * @generated from protobuf enum value: RangedWeaponTypeSigil = 9;
     */
    RangedWeaponType[RangedWeaponType["RangedWeaponTypeSigil"] = 9] = "RangedWeaponTypeSigil";
})(RangedWeaponType || (RangedWeaponType = {}));
/**
 * All slots on the gear menu where a single item can be worn.
 *
 * @generated from protobuf enum proto.ItemSlot
 */
export var ItemSlot;
(function (ItemSlot) {
    /**
     * @generated from protobuf enum value: ItemSlotHead = 0;
     */
    ItemSlot[ItemSlot["ItemSlotHead"] = 0] = "ItemSlotHead";
    /**
     * @generated from protobuf enum value: ItemSlotNeck = 1;
     */
    ItemSlot[ItemSlot["ItemSlotNeck"] = 1] = "ItemSlotNeck";
    /**
     * @generated from protobuf enum value: ItemSlotShoulder = 2;
     */
    ItemSlot[ItemSlot["ItemSlotShoulder"] = 2] = "ItemSlotShoulder";
    /**
     * @generated from protobuf enum value: ItemSlotBack = 3;
     */
    ItemSlot[ItemSlot["ItemSlotBack"] = 3] = "ItemSlotBack";
    /**
     * @generated from protobuf enum value: ItemSlotChest = 4;
     */
    ItemSlot[ItemSlot["ItemSlotChest"] = 4] = "ItemSlotChest";
    /**
     * @generated from protobuf enum value: ItemSlotWrist = 5;
     */
    ItemSlot[ItemSlot["ItemSlotWrist"] = 5] = "ItemSlotWrist";
    /**
     * @generated from protobuf enum value: ItemSlotHands = 6;
     */
    ItemSlot[ItemSlot["ItemSlotHands"] = 6] = "ItemSlotHands";
    /**
     * @generated from protobuf enum value: ItemSlotWaist = 7;
     */
    ItemSlot[ItemSlot["ItemSlotWaist"] = 7] = "ItemSlotWaist";
    /**
     * @generated from protobuf enum value: ItemSlotLegs = 8;
     */
    ItemSlot[ItemSlot["ItemSlotLegs"] = 8] = "ItemSlotLegs";
    /**
     * @generated from protobuf enum value: ItemSlotFeet = 9;
     */
    ItemSlot[ItemSlot["ItemSlotFeet"] = 9] = "ItemSlotFeet";
    /**
     * @generated from protobuf enum value: ItemSlotFinger1 = 10;
     */
    ItemSlot[ItemSlot["ItemSlotFinger1"] = 10] = "ItemSlotFinger1";
    /**
     * @generated from protobuf enum value: ItemSlotFinger2 = 11;
     */
    ItemSlot[ItemSlot["ItemSlotFinger2"] = 11] = "ItemSlotFinger2";
    /**
     * @generated from protobuf enum value: ItemSlotTrinket1 = 12;
     */
    ItemSlot[ItemSlot["ItemSlotTrinket1"] = 12] = "ItemSlotTrinket1";
    /**
     * @generated from protobuf enum value: ItemSlotTrinket2 = 13;
     */
    ItemSlot[ItemSlot["ItemSlotTrinket2"] = 13] = "ItemSlotTrinket2";
    /**
     * can be 1h or 2h
     *
     * @generated from protobuf enum value: ItemSlotMainHand = 14;
     */
    ItemSlot[ItemSlot["ItemSlotMainHand"] = 14] = "ItemSlotMainHand";
    /**
     * @generated from protobuf enum value: ItemSlotOffHand = 15;
     */
    ItemSlot[ItemSlot["ItemSlotOffHand"] = 15] = "ItemSlotOffHand";
    /**
     * @generated from protobuf enum value: ItemSlotRanged = 16;
     */
    ItemSlot[ItemSlot["ItemSlotRanged"] = 16] = "ItemSlotRanged";
})(ItemSlot || (ItemSlot = {}));
/**
 * @generated from protobuf enum proto.ItemQuality
 */
export var ItemQuality;
(function (ItemQuality) {
    /**
     * @generated from protobuf enum value: ItemQualityJunk = 0;
     */
    ItemQuality[ItemQuality["ItemQualityJunk"] = 0] = "ItemQualityJunk";
    /**
     * @generated from protobuf enum value: ItemQualityCommon = 1;
     */
    ItemQuality[ItemQuality["ItemQualityCommon"] = 1] = "ItemQualityCommon";
    /**
     * @generated from protobuf enum value: ItemQualityUncommon = 2;
     */
    ItemQuality[ItemQuality["ItemQualityUncommon"] = 2] = "ItemQualityUncommon";
    /**
     * @generated from protobuf enum value: ItemQualityRare = 3;
     */
    ItemQuality[ItemQuality["ItemQualityRare"] = 3] = "ItemQualityRare";
    /**
     * @generated from protobuf enum value: ItemQualityEpic = 4;
     */
    ItemQuality[ItemQuality["ItemQualityEpic"] = 4] = "ItemQualityEpic";
    /**
     * @generated from protobuf enum value: ItemQualityLegendary = 5;
     */
    ItemQuality[ItemQuality["ItemQualityLegendary"] = 5] = "ItemQualityLegendary";
})(ItemQuality || (ItemQuality = {}));
/**
 * @generated from protobuf enum proto.GemColor
 */
export var GemColor;
(function (GemColor) {
    /**
     * @generated from protobuf enum value: GemColorUnknown = 0;
     */
    GemColor[GemColor["GemColorUnknown"] = 0] = "GemColorUnknown";
    /**
     * @generated from protobuf enum value: GemColorMeta = 1;
     */
    GemColor[GemColor["GemColorMeta"] = 1] = "GemColorMeta";
    /**
     * @generated from protobuf enum value: GemColorRed = 2;
     */
    GemColor[GemColor["GemColorRed"] = 2] = "GemColorRed";
    /**
     * @generated from protobuf enum value: GemColorBlue = 3;
     */
    GemColor[GemColor["GemColorBlue"] = 3] = "GemColorBlue";
    /**
     * @generated from protobuf enum value: GemColorYellow = 4;
     */
    GemColor[GemColor["GemColorYellow"] = 4] = "GemColorYellow";
    /**
     * @generated from protobuf enum value: GemColorGreen = 5;
     */
    GemColor[GemColor["GemColorGreen"] = 5] = "GemColorGreen";
    /**
     * @generated from protobuf enum value: GemColorOrange = 6;
     */
    GemColor[GemColor["GemColorOrange"] = 6] = "GemColorOrange";
    /**
     * @generated from protobuf enum value: GemColorPurple = 7;
     */
    GemColor[GemColor["GemColorPurple"] = 7] = "GemColorPurple";
    /**
     * @generated from protobuf enum value: GemColorPrismatic = 8;
     */
    GemColor[GemColor["GemColorPrismatic"] = 8] = "GemColorPrismatic";
})(GemColor || (GemColor = {}));
/**
 * @generated from protobuf enum proto.SpellSchool
 */
export var SpellSchool;
(function (SpellSchool) {
    /**
     * @generated from protobuf enum value: SpellSchoolPhysical = 0;
     */
    SpellSchool[SpellSchool["SpellSchoolPhysical"] = 0] = "SpellSchoolPhysical";
    /**
     * @generated from protobuf enum value: SpellSchoolArcane = 1;
     */
    SpellSchool[SpellSchool["SpellSchoolArcane"] = 1] = "SpellSchoolArcane";
    /**
     * @generated from protobuf enum value: SpellSchoolFire = 2;
     */
    SpellSchool[SpellSchool["SpellSchoolFire"] = 2] = "SpellSchoolFire";
    /**
     * @generated from protobuf enum value: SpellSchoolFrost = 3;
     */
    SpellSchool[SpellSchool["SpellSchoolFrost"] = 3] = "SpellSchoolFrost";
    /**
     * @generated from protobuf enum value: SpellSchoolHoly = 4;
     */
    SpellSchool[SpellSchool["SpellSchoolHoly"] = 4] = "SpellSchoolHoly";
    /**
     * @generated from protobuf enum value: SpellSchoolNature = 5;
     */
    SpellSchool[SpellSchool["SpellSchoolNature"] = 5] = "SpellSchoolNature";
    /**
     * @generated from protobuf enum value: SpellSchoolShadow = 6;
     */
    SpellSchool[SpellSchool["SpellSchoolShadow"] = 6] = "SpellSchoolShadow";
})(SpellSchool || (SpellSchool = {}));
/**
 * @generated from protobuf enum proto.TristateEffect
 */
export var TristateEffect;
(function (TristateEffect) {
    /**
     * @generated from protobuf enum value: TristateEffectMissing = 0;
     */
    TristateEffect[TristateEffect["TristateEffectMissing"] = 0] = "TristateEffectMissing";
    /**
     * @generated from protobuf enum value: TristateEffectRegular = 1;
     */
    TristateEffect[TristateEffect["TristateEffectRegular"] = 1] = "TristateEffectRegular";
    /**
     * @generated from protobuf enum value: TristateEffectImproved = 2;
     */
    TristateEffect[TristateEffect["TristateEffectImproved"] = 2] = "TristateEffectImproved";
})(TristateEffect || (TristateEffect = {}));
/**
 * @generated from protobuf enum proto.Explosive
 */
export var Explosive;
(function (Explosive) {
    /**
     * @generated from protobuf enum value: ExplosiveUnknown = 0;
     */
    Explosive[Explosive["ExplosiveUnknown"] = 0] = "ExplosiveUnknown";
    /**
     * @generated from protobuf enum value: ExplosiveFelIronBomb = 1;
     */
    Explosive[Explosive["ExplosiveFelIronBomb"] = 1] = "ExplosiveFelIronBomb";
    /**
     * @generated from protobuf enum value: ExplosiveAdamantiteGrenade = 2;
     */
    Explosive[Explosive["ExplosiveAdamantiteGrenade"] = 2] = "ExplosiveAdamantiteGrenade";
    /**
     * @generated from protobuf enum value: ExplosiveGnomishFlameTurret = 3;
     */
    Explosive[Explosive["ExplosiveGnomishFlameTurret"] = 3] = "ExplosiveGnomishFlameTurret";
    /**
     * @generated from protobuf enum value: ExplosiveHolyWater = 4;
     */
    Explosive[Explosive["ExplosiveHolyWater"] = 4] = "ExplosiveHolyWater";
})(Explosive || (Explosive = {}));
/**
 * @generated from protobuf enum proto.Potions
 */
export var Potions;
(function (Potions) {
    /**
     * @generated from protobuf enum value: UnknownPotion = 0;
     */
    Potions[Potions["UnknownPotion"] = 0] = "UnknownPotion";
    /**
     * @generated from protobuf enum value: RunicHealingPotion = 1;
     */
    Potions[Potions["RunicHealingPotion"] = 1] = "RunicHealingPotion";
    /**
     * @generated from protobuf enum value: RunicManaPotion = 2;
     */
    Potions[Potions["RunicManaPotion"] = 2] = "RunicManaPotion";
    /**
     * @generated from protobuf enum value: IndestructiblePotion = 3;
     */
    Potions[Potions["IndestructiblePotion"] = 3] = "IndestructiblePotion";
    /**
     * @generated from protobuf enum value: PotionOfSpeed = 4;
     */
    Potions[Potions["PotionOfSpeed"] = 4] = "PotionOfSpeed";
    /**
     * @generated from protobuf enum value: PotionOfWildMagic = 5;
     */
    Potions[Potions["PotionOfWildMagic"] = 5] = "PotionOfWildMagic";
    /**
     * @generated from protobuf enum value: DestructionPotion = 6;
     */
    Potions[Potions["DestructionPotion"] = 6] = "DestructionPotion";
    /**
     * @generated from protobuf enum value: SuperManaPotion = 7;
     */
    Potions[Potions["SuperManaPotion"] = 7] = "SuperManaPotion";
    /**
     * @generated from protobuf enum value: HastePotion = 8;
     */
    Potions[Potions["HastePotion"] = 8] = "HastePotion";
    /**
     * @generated from protobuf enum value: MightyRagePotion = 9;
     */
    Potions[Potions["MightyRagePotion"] = 9] = "MightyRagePotion";
    /**
     * @generated from protobuf enum value: FelManaPotion = 10;
     */
    Potions[Potions["FelManaPotion"] = 10] = "FelManaPotion";
    /**
     * @generated from protobuf enum value: InsaneStrengthPotion = 11;
     */
    Potions[Potions["InsaneStrengthPotion"] = 11] = "InsaneStrengthPotion";
    /**
     * @generated from protobuf enum value: IronshieldPotion = 12;
     */
    Potions[Potions["IronshieldPotion"] = 12] = "IronshieldPotion";
    /**
     * @generated from protobuf enum value: HeroicPotion = 13;
     */
    Potions[Potions["HeroicPotion"] = 13] = "HeroicPotion";
})(Potions || (Potions = {}));
/**
 * @generated from protobuf enum proto.Conjured
 */
export var Conjured;
(function (Conjured) {
    /**
     * @generated from protobuf enum value: ConjuredUnknown = 0;
     */
    Conjured[Conjured["ConjuredUnknown"] = 0] = "ConjuredUnknown";
    /**
     * @generated from protobuf enum value: ConjuredDarkRune = 1;
     */
    Conjured[Conjured["ConjuredDarkRune"] = 1] = "ConjuredDarkRune";
    /**
     * @generated from protobuf enum value: ConjuredFlameCap = 2;
     */
    Conjured[Conjured["ConjuredFlameCap"] = 2] = "ConjuredFlameCap";
    /**
     * @generated from protobuf enum value: ConjuredHealthstone = 5;
     */
    Conjured[Conjured["ConjuredHealthstone"] = 5] = "ConjuredHealthstone";
    /**
     * @generated from protobuf enum value: ConjuredMageManaEmerald = 3;
     */
    Conjured[Conjured["ConjuredMageManaEmerald"] = 3] = "ConjuredMageManaEmerald";
    /**
     * @generated from protobuf enum value: ConjuredRogueThistleTea = 4;
     */
    Conjured[Conjured["ConjuredRogueThistleTea"] = 4] = "ConjuredRogueThistleTea";
})(Conjured || (Conjured = {}));
/**
 * @generated from protobuf enum proto.WeaponImbue
 */
export var WeaponImbue;
(function (WeaponImbue) {
    /**
     * @generated from protobuf enum value: WeaponImbueUnknown = 0;
     */
    WeaponImbue[WeaponImbue["WeaponImbueUnknown"] = 0] = "WeaponImbueUnknown";
    /**
     * @generated from protobuf enum value: WeaponImbueAdamantiteSharpeningStone = 1;
     */
    WeaponImbue[WeaponImbue["WeaponImbueAdamantiteSharpeningStone"] = 1] = "WeaponImbueAdamantiteSharpeningStone";
    /**
     * @generated from protobuf enum value: WeaponImbueAdamantiteWeightstone = 5;
     */
    WeaponImbue[WeaponImbue["WeaponImbueAdamantiteWeightstone"] = 5] = "WeaponImbueAdamantiteWeightstone";
    /**
     * @generated from protobuf enum value: WeaponImbueElementalSharpeningStone = 2;
     */
    WeaponImbue[WeaponImbue["WeaponImbueElementalSharpeningStone"] = 2] = "WeaponImbueElementalSharpeningStone";
    /**
     * @generated from protobuf enum value: WeaponImbueBrilliantWizardOil = 3;
     */
    WeaponImbue[WeaponImbue["WeaponImbueBrilliantWizardOil"] = 3] = "WeaponImbueBrilliantWizardOil";
    /**
     * @generated from protobuf enum value: WeaponImbueSuperiorWizardOil = 4;
     */
    WeaponImbue[WeaponImbue["WeaponImbueSuperiorWizardOil"] = 4] = "WeaponImbueSuperiorWizardOil";
    /**
     * @generated from protobuf enum value: WeaponImbueShamanFlametongue = 6;
     */
    WeaponImbue[WeaponImbue["WeaponImbueShamanFlametongue"] = 6] = "WeaponImbueShamanFlametongue";
    /**
     * @generated from protobuf enum value: WeaponImbueShamanFrostbrand = 7;
     */
    WeaponImbue[WeaponImbue["WeaponImbueShamanFrostbrand"] = 7] = "WeaponImbueShamanFrostbrand";
    /**
     * @generated from protobuf enum value: WeaponImbueShamanRockbiter = 8;
     */
    WeaponImbue[WeaponImbue["WeaponImbueShamanRockbiter"] = 8] = "WeaponImbueShamanRockbiter";
    /**
     * @generated from protobuf enum value: WeaponImbueShamanWindfury = 9;
     */
    WeaponImbue[WeaponImbue["WeaponImbueShamanWindfury"] = 9] = "WeaponImbueShamanWindfury";
    /**
     * @generated from protobuf enum value: WeaponImbueRogueDeadlyPoison = 10;
     */
    WeaponImbue[WeaponImbue["WeaponImbueRogueDeadlyPoison"] = 10] = "WeaponImbueRogueDeadlyPoison";
    /**
     * @generated from protobuf enum value: WeaponImbueRogueInstantPoison = 11;
     */
    WeaponImbue[WeaponImbue["WeaponImbueRogueInstantPoison"] = 11] = "WeaponImbueRogueInstantPoison";
    /**
     * @generated from protobuf enum value: WeaponImbueRighteousWeaponCoating = 12;
     */
    WeaponImbue[WeaponImbue["WeaponImbueRighteousWeaponCoating"] = 12] = "WeaponImbueRighteousWeaponCoating";
})(WeaponImbue || (WeaponImbue = {}));
/**
 * @generated from protobuf enum proto.Flask
 */
export var Flask;
(function (Flask) {
    /**
     * @generated from protobuf enum value: FlaskUnknown = 0;
     */
    Flask[Flask["FlaskUnknown"] = 0] = "FlaskUnknown";
    /**
     * @generated from protobuf enum value: FlaskOfTheFrostWyrm = 1;
     */
    Flask[Flask["FlaskOfTheFrostWyrm"] = 1] = "FlaskOfTheFrostWyrm";
    /**
     * @generated from protobuf enum value: FlaskOfEndlessRage = 2;
     */
    Flask[Flask["FlaskOfEndlessRage"] = 2] = "FlaskOfEndlessRage";
    /**
     * @generated from protobuf enum value: FlaskOfPureMojo = 3;
     */
    Flask[Flask["FlaskOfPureMojo"] = 3] = "FlaskOfPureMojo";
    /**
     * @generated from protobuf enum value: FlaskOfStoneblood = 4;
     */
    Flask[Flask["FlaskOfStoneblood"] = 4] = "FlaskOfStoneblood";
    /**
     * @generated from protobuf enum value: LesserFlaskOfToughness = 5;
     */
    Flask[Flask["LesserFlaskOfToughness"] = 5] = "LesserFlaskOfToughness";
    /**
     * @generated from protobuf enum value: LesserFlaskOfResistance = 6;
     */
    Flask[Flask["LesserFlaskOfResistance"] = 6] = "LesserFlaskOfResistance";
    /**
     * TBC
     *
     * @generated from protobuf enum value: FlaskOfBlindingLight = 7;
     */
    Flask[Flask["FlaskOfBlindingLight"] = 7] = "FlaskOfBlindingLight";
    /**
     * @generated from protobuf enum value: FlaskOfMightyRestoration = 8;
     */
    Flask[Flask["FlaskOfMightyRestoration"] = 8] = "FlaskOfMightyRestoration";
    /**
     * @generated from protobuf enum value: FlaskOfPureDeath = 9;
     */
    Flask[Flask["FlaskOfPureDeath"] = 9] = "FlaskOfPureDeath";
    /**
     * @generated from protobuf enum value: FlaskOfRelentlessAssault = 10;
     */
    Flask[Flask["FlaskOfRelentlessAssault"] = 10] = "FlaskOfRelentlessAssault";
    /**
     * @generated from protobuf enum value: FlaskOfSupremePower = 11;
     */
    Flask[Flask["FlaskOfSupremePower"] = 11] = "FlaskOfSupremePower";
    /**
     * @generated from protobuf enum value: FlaskOfFortification = 12;
     */
    Flask[Flask["FlaskOfFortification"] = 12] = "FlaskOfFortification";
    /**
     * @generated from protobuf enum value: FlaskOfChromaticWonder = 13;
     */
    Flask[Flask["FlaskOfChromaticWonder"] = 13] = "FlaskOfChromaticWonder";
})(Flask || (Flask = {}));
/**
 * @generated from protobuf enum proto.BattleElixir
 */
export var BattleElixir;
(function (BattleElixir) {
    /**
     * @generated from protobuf enum value: BattleElixirUnknown = 0;
     */
    BattleElixir[BattleElixir["BattleElixirUnknown"] = 0] = "BattleElixirUnknown";
    /**
     * @generated from protobuf enum value: ElixirOfAccuracy = 1;
     */
    BattleElixir[BattleElixir["ElixirOfAccuracy"] = 1] = "ElixirOfAccuracy";
    /**
     * @generated from protobuf enum value: ElixirOfArmorPiercing = 2;
     */
    BattleElixir[BattleElixir["ElixirOfArmorPiercing"] = 2] = "ElixirOfArmorPiercing";
    /**
     * @generated from protobuf enum value: ElixirOfDeadlyStrikes = 3;
     */
    BattleElixir[BattleElixir["ElixirOfDeadlyStrikes"] = 3] = "ElixirOfDeadlyStrikes";
    /**
     * @generated from protobuf enum value: ElixirOfExpertise = 4;
     */
    BattleElixir[BattleElixir["ElixirOfExpertise"] = 4] = "ElixirOfExpertise";
    /**
     * @generated from protobuf enum value: ElixirOfLightningSpeed = 5;
     */
    BattleElixir[BattleElixir["ElixirOfLightningSpeed"] = 5] = "ElixirOfLightningSpeed";
    /**
     * @generated from protobuf enum value: ElixirOfMightyAgility = 6;
     */
    BattleElixir[BattleElixir["ElixirOfMightyAgility"] = 6] = "ElixirOfMightyAgility";
    /**
     * @generated from protobuf enum value: ElixirOfMightyStrength = 7;
     */
    BattleElixir[BattleElixir["ElixirOfMightyStrength"] = 7] = "ElixirOfMightyStrength";
    /**
     * @generated from protobuf enum value: GurusElixir = 8;
     */
    BattleElixir[BattleElixir["GurusElixir"] = 8] = "GurusElixir";
    /**
     * @generated from protobuf enum value: SpellpowerElixir = 9;
     */
    BattleElixir[BattleElixir["SpellpowerElixir"] = 9] = "SpellpowerElixir";
    /**
     * @generated from protobuf enum value: WrathElixir = 10;
     */
    BattleElixir[BattleElixir["WrathElixir"] = 10] = "WrathElixir";
    /**
     * TBC
     *
     * @generated from protobuf enum value: AdeptsElixir = 11;
     */
    BattleElixir[BattleElixir["AdeptsElixir"] = 11] = "AdeptsElixir";
    /**
     * @generated from protobuf enum value: ElixirOfDemonslaying = 12;
     */
    BattleElixir[BattleElixir["ElixirOfDemonslaying"] = 12] = "ElixirOfDemonslaying";
    /**
     * @generated from protobuf enum value: ElixirOfMajorAgility = 13;
     */
    BattleElixir[BattleElixir["ElixirOfMajorAgility"] = 13] = "ElixirOfMajorAgility";
    /**
     * @generated from protobuf enum value: ElixirOfMajorFirePower = 14;
     */
    BattleElixir[BattleElixir["ElixirOfMajorFirePower"] = 14] = "ElixirOfMajorFirePower";
    /**
     * @generated from protobuf enum value: ElixirOfMajorFrostPower = 15;
     */
    BattleElixir[BattleElixir["ElixirOfMajorFrostPower"] = 15] = "ElixirOfMajorFrostPower";
    /**
     * @generated from protobuf enum value: ElixirOfMajorShadowPower = 16;
     */
    BattleElixir[BattleElixir["ElixirOfMajorShadowPower"] = 16] = "ElixirOfMajorShadowPower";
    /**
     * @generated from protobuf enum value: ElixirOfMajorStrength = 17;
     */
    BattleElixir[BattleElixir["ElixirOfMajorStrength"] = 17] = "ElixirOfMajorStrength";
    /**
     * @generated from protobuf enum value: ElixirOfMastery = 18;
     */
    BattleElixir[BattleElixir["ElixirOfMastery"] = 18] = "ElixirOfMastery";
    /**
     * @generated from protobuf enum value: ElixirOfTheMongoose = 19;
     */
    BattleElixir[BattleElixir["ElixirOfTheMongoose"] = 19] = "ElixirOfTheMongoose";
    /**
     * @generated from protobuf enum value: FelStrengthElixir = 20;
     */
    BattleElixir[BattleElixir["FelStrengthElixir"] = 20] = "FelStrengthElixir";
    /**
     * @generated from protobuf enum value: GreaterArcaneElixir = 21;
     */
    BattleElixir[BattleElixir["GreaterArcaneElixir"] = 21] = "GreaterArcaneElixir";
})(BattleElixir || (BattleElixir = {}));
/**
 * @generated from protobuf enum proto.GuardianElixir
 */
export var GuardianElixir;
(function (GuardianElixir) {
    /**
     * @generated from protobuf enum value: GuardianElixirUnknown = 0;
     */
    GuardianElixir[GuardianElixir["GuardianElixirUnknown"] = 0] = "GuardianElixirUnknown";
    /**
     * @generated from protobuf enum value: ElixirOfMightyDefense = 1;
     */
    GuardianElixir[GuardianElixir["ElixirOfMightyDefense"] = 1] = "ElixirOfMightyDefense";
    /**
     * @generated from protobuf enum value: ElixirOfMightyFortitude = 2;
     */
    GuardianElixir[GuardianElixir["ElixirOfMightyFortitude"] = 2] = "ElixirOfMightyFortitude";
    /**
     * @generated from protobuf enum value: ElixirOfMightyMageblood = 3;
     */
    GuardianElixir[GuardianElixir["ElixirOfMightyMageblood"] = 3] = "ElixirOfMightyMageblood";
    /**
     * @generated from protobuf enum value: ElixirOfMightyThoughts = 4;
     */
    GuardianElixir[GuardianElixir["ElixirOfMightyThoughts"] = 4] = "ElixirOfMightyThoughts";
    /**
     * @generated from protobuf enum value: ElixirOfProtection = 5;
     */
    GuardianElixir[GuardianElixir["ElixirOfProtection"] = 5] = "ElixirOfProtection";
    /**
     * @generated from protobuf enum value: ElixirOfSpirit = 6;
     */
    GuardianElixir[GuardianElixir["ElixirOfSpirit"] = 6] = "ElixirOfSpirit";
    /**
     * TBC
     *
     * @generated from protobuf enum value: GiftOfArthas = 7;
     */
    GuardianElixir[GuardianElixir["GiftOfArthas"] = 7] = "GiftOfArthas";
    /**
     * @generated from protobuf enum value: ElixirOfDraenicWisdom = 8;
     */
    GuardianElixir[GuardianElixir["ElixirOfDraenicWisdom"] = 8] = "ElixirOfDraenicWisdom";
    /**
     * @generated from protobuf enum value: ElixirOfIronskin = 9;
     */
    GuardianElixir[GuardianElixir["ElixirOfIronskin"] = 9] = "ElixirOfIronskin";
    /**
     * @generated from protobuf enum value: ElixirOfMajorDefense = 10;
     */
    GuardianElixir[GuardianElixir["ElixirOfMajorDefense"] = 10] = "ElixirOfMajorDefense";
    /**
     * @generated from protobuf enum value: ElixirOfMajorFortitude = 11;
     */
    GuardianElixir[GuardianElixir["ElixirOfMajorFortitude"] = 11] = "ElixirOfMajorFortitude";
    /**
     * @generated from protobuf enum value: ElixirOfMajorMageblood = 12;
     */
    GuardianElixir[GuardianElixir["ElixirOfMajorMageblood"] = 12] = "ElixirOfMajorMageblood";
})(GuardianElixir || (GuardianElixir = {}));
/**
 * @generated from protobuf enum proto.Food
 */
export var Food;
(function (Food) {
    /**
     * @generated from protobuf enum value: FoodUnknown = 0;
     */
    Food[Food["FoodUnknown"] = 0] = "FoodUnknown";
    /**
     * @generated from protobuf enum value: FoodFishFeast = 1;
     */
    Food[Food["FoodFishFeast"] = 1] = "FoodFishFeast";
    /**
     * @generated from protobuf enum value: FoodGreatFeast = 2;
     */
    Food[Food["FoodGreatFeast"] = 2] = "FoodGreatFeast";
    /**
     * @generated from protobuf enum value: FoodBlackenedDragonfin = 3;
     */
    Food[Food["FoodBlackenedDragonfin"] = 3] = "FoodBlackenedDragonfin";
    /**
     * @generated from protobuf enum value: FoodHeartyRhino = 4;
     */
    Food[Food["FoodHeartyRhino"] = 4] = "FoodHeartyRhino";
    /**
     * @generated from protobuf enum value: FoodMegaMammothMeal = 5;
     */
    Food[Food["FoodMegaMammothMeal"] = 5] = "FoodMegaMammothMeal";
    /**
     * @generated from protobuf enum value: FoodSpicedWormBurger = 6;
     */
    Food[Food["FoodSpicedWormBurger"] = 6] = "FoodSpicedWormBurger";
    /**
     * @generated from protobuf enum value: FoodRhinoliciousWormsteak = 7;
     */
    Food[Food["FoodRhinoliciousWormsteak"] = 7] = "FoodRhinoliciousWormsteak";
    /**
     * @generated from protobuf enum value: FoodImperialMantaSteak = 8;
     */
    Food[Food["FoodImperialMantaSteak"] = 8] = "FoodImperialMantaSteak";
    /**
     * @generated from protobuf enum value: FoodSnapperExtreme = 9;
     */
    Food[Food["FoodSnapperExtreme"] = 9] = "FoodSnapperExtreme";
    /**
     * @generated from protobuf enum value: FoodMightyRhinoDogs = 10;
     */
    Food[Food["FoodMightyRhinoDogs"] = 10] = "FoodMightyRhinoDogs";
    /**
     * @generated from protobuf enum value: FoodFirecrackerSalmon = 11;
     */
    Food[Food["FoodFirecrackerSalmon"] = 11] = "FoodFirecrackerSalmon";
    /**
     * @generated from protobuf enum value: FoodCuttlesteak = 12;
     */
    Food[Food["FoodCuttlesteak"] = 12] = "FoodCuttlesteak";
    /**
     * @generated from protobuf enum value: FoodDragonfinFilet = 13;
     */
    Food[Food["FoodDragonfinFilet"] = 13] = "FoodDragonfinFilet";
    /**
     * TBC Foods
     *
     * @generated from protobuf enum value: FoodBlackenedBasilisk = 14;
     */
    Food[Food["FoodBlackenedBasilisk"] = 14] = "FoodBlackenedBasilisk";
    /**
     * @generated from protobuf enum value: FoodGrilledMudfish = 15;
     */
    Food[Food["FoodGrilledMudfish"] = 15] = "FoodGrilledMudfish";
    /**
     * @generated from protobuf enum value: FoodRavagerDog = 16;
     */
    Food[Food["FoodRavagerDog"] = 16] = "FoodRavagerDog";
    /**
     * @generated from protobuf enum value: FoodRoastedClefthoof = 17;
     */
    Food[Food["FoodRoastedClefthoof"] = 17] = "FoodRoastedClefthoof";
    /**
     * @generated from protobuf enum value: FoodSkullfishSoup = 18;
     */
    Food[Food["FoodSkullfishSoup"] = 18] = "FoodSkullfishSoup";
    /**
     * @generated from protobuf enum value: FoodSpicyHotTalbuk = 19;
     */
    Food[Food["FoodSpicyHotTalbuk"] = 19] = "FoodSpicyHotTalbuk";
    /**
     * @generated from protobuf enum value: FoodFishermansFeast = 20;
     */
    Food[Food["FoodFishermansFeast"] = 20] = "FoodFishermansFeast";
})(Food || (Food = {}));
/**
 * @generated from protobuf enum proto.PetFood
 */
export var PetFood;
(function (PetFood) {
    /**
     * @generated from protobuf enum value: PetFoodUnknown = 0;
     */
    PetFood[PetFood["PetFoodUnknown"] = 0] = "PetFoodUnknown";
    /**
     * @generated from protobuf enum value: PetFoodSpicedMammothTreats = 1;
     */
    PetFood[PetFood["PetFoodSpicedMammothTreats"] = 1] = "PetFoodSpicedMammothTreats";
    /**
     * TBC
     *
     * @generated from protobuf enum value: PetFoodKiblersBits = 2;
     */
    PetFood[PetFood["PetFoodKiblersBits"] = 2] = "PetFoodKiblersBits";
})(PetFood || (PetFood = {}));
/**
 * @generated from protobuf enum proto.MobType
 */
export var MobType;
(function (MobType) {
    /**
     * @generated from protobuf enum value: MobTypeUnknown = 0;
     */
    MobType[MobType["MobTypeUnknown"] = 0] = "MobTypeUnknown";
    /**
     * @generated from protobuf enum value: MobTypeBeast = 1;
     */
    MobType[MobType["MobTypeBeast"] = 1] = "MobTypeBeast";
    /**
     * @generated from protobuf enum value: MobTypeDemon = 2;
     */
    MobType[MobType["MobTypeDemon"] = 2] = "MobTypeDemon";
    /**
     * @generated from protobuf enum value: MobTypeDragonkin = 3;
     */
    MobType[MobType["MobTypeDragonkin"] = 3] = "MobTypeDragonkin";
    /**
     * @generated from protobuf enum value: MobTypeElemental = 4;
     */
    MobType[MobType["MobTypeElemental"] = 4] = "MobTypeElemental";
    /**
     * @generated from protobuf enum value: MobTypeGiant = 5;
     */
    MobType[MobType["MobTypeGiant"] = 5] = "MobTypeGiant";
    /**
     * @generated from protobuf enum value: MobTypeHumanoid = 6;
     */
    MobType[MobType["MobTypeHumanoid"] = 6] = "MobTypeHumanoid";
    /**
     * @generated from protobuf enum value: MobTypeMechanical = 7;
     */
    MobType[MobType["MobTypeMechanical"] = 7] = "MobTypeMechanical";
    /**
     * @generated from protobuf enum value: MobTypeUndead = 8;
     */
    MobType[MobType["MobTypeUndead"] = 8] = "MobTypeUndead";
})(MobType || (MobType = {}));
/**
 * Extra enum for describing which items are eligible for an enchant, when
 * ItemType alone is not enough.
 *
 * @generated from protobuf enum proto.EnchantType
 */
export var EnchantType;
(function (EnchantType) {
    /**
     * @generated from protobuf enum value: EnchantTypeNormal = 0;
     */
    EnchantType[EnchantType["EnchantTypeNormal"] = 0] = "EnchantTypeNormal";
    /**
     * @generated from protobuf enum value: EnchantTypeTwoHand = 1;
     */
    EnchantType[EnchantType["EnchantTypeTwoHand"] = 1] = "EnchantTypeTwoHand";
    /**
     * @generated from protobuf enum value: EnchantTypeShield = 2;
     */
    EnchantType[EnchantType["EnchantTypeShield"] = 2] = "EnchantTypeShield";
    /**
     * @generated from protobuf enum value: EnchantTypeKit = 3;
     */
    EnchantType[EnchantType["EnchantTypeKit"] = 3] = "EnchantTypeKit";
})(EnchantType || (EnchantType = {}));
/**
 * ID for actions that aren't spells or items.
 *
 * @generated from protobuf enum proto.OtherAction
 */
export var OtherAction;
(function (OtherAction) {
    /**
     * @generated from protobuf enum value: OtherActionNone = 0;
     */
    OtherAction[OtherAction["OtherActionNone"] = 0] = "OtherActionNone";
    /**
     * @generated from protobuf enum value: OtherActionWait = 1;
     */
    OtherAction[OtherAction["OtherActionWait"] = 1] = "OtherActionWait";
    /**
     * @generated from protobuf enum value: OtherActionManaRegen = 2;
     */
    OtherAction[OtherAction["OtherActionManaRegen"] = 2] = "OtherActionManaRegen";
    /**
     * @generated from protobuf enum value: OtherActionEnergyRegen = 5;
     */
    OtherAction[OtherAction["OtherActionEnergyRegen"] = 5] = "OtherActionEnergyRegen";
    /**
     * @generated from protobuf enum value: OtherActionFocusRegen = 6;
     */
    OtherAction[OtherAction["OtherActionFocusRegen"] = 6] = "OtherActionFocusRegen";
    /**
     * For threat generated from mana gains.
     *
     * @generated from protobuf enum value: OtherActionManaGain = 10;
     */
    OtherAction[OtherAction["OtherActionManaGain"] = 10] = "OtherActionManaGain";
    /**
     * For threat generated from rage gains.
     *
     * @generated from protobuf enum value: OtherActionRageGain = 11;
     */
    OtherAction[OtherAction["OtherActionRageGain"] = 11] = "OtherActionRageGain";
    /**
     * A white hit, can be main hand or off hand.
     *
     * @generated from protobuf enum value: OtherActionAttack = 3;
     */
    OtherAction[OtherAction["OtherActionAttack"] = 3] = "OtherActionAttack";
    /**
     * Default shoot action using a wand/bow/gun.
     *
     * @generated from protobuf enum value: OtherActionShoot = 4;
     */
    OtherAction[OtherAction["OtherActionShoot"] = 4] = "OtherActionShoot";
    /**
     * Represents a grouping of all pet actions. Only used by the UI.
     *
     * @generated from protobuf enum value: OtherActionPet = 7;
     */
    OtherAction[OtherAction["OtherActionPet"] = 7] = "OtherActionPet";
    /**
     * Refund of a resource like Energy or Rage, when the ability didn't land.
     *
     * @generated from protobuf enum value: OtherActionRefund = 8;
     */
    OtherAction[OtherAction["OtherActionRefund"] = 8] = "OtherActionRefund";
    /**
     * Indicates damage taken; used for rage gen.
     *
     * @generated from protobuf enum value: OtherActionDamageTaken = 9;
     */
    OtherAction[OtherAction["OtherActionDamageTaken"] = 9] = "OtherActionDamageTaken";
    /**
     * Indicates healing received from healing model.
     *
     * @generated from protobuf enum value: OtherActionHealingModel = 12;
     */
    OtherAction[OtherAction["OtherActionHealingModel"] = 12] = "OtherActionHealingModel";
    /**
     * Indicates healing received from healing model.
     *
     * @generated from protobuf enum value: OtherActionBloodRuneGain = 13;
     */
    OtherAction[OtherAction["OtherActionBloodRuneGain"] = 13] = "OtherActionBloodRuneGain";
    /**
     * Indicates healing received from healing model.
     *
     * @generated from protobuf enum value: OtherActionFrostRuneGain = 14;
     */
    OtherAction[OtherAction["OtherActionFrostRuneGain"] = 14] = "OtherActionFrostRuneGain";
    /**
     * Indicates healing received from healing model.
     *
     * @generated from protobuf enum value: OtherActionUnholyRuneGain = 15;
     */
    OtherAction[OtherAction["OtherActionUnholyRuneGain"] = 15] = "OtherActionUnholyRuneGain";
    /**
     * Indicates healing received from healing model.
     *
     * @generated from protobuf enum value: OtherActionDeathRuneGain = 16;
     */
    OtherAction[OtherAction["OtherActionDeathRuneGain"] = 16] = "OtherActionDeathRuneGain";
})(OtherAction || (OtherAction = {}));
// @generated message type with reflection information, may provide speed optimized methods
class RaidBuffs$Type extends MessageType {
    constructor() {
        super("proto.RaidBuffs", [
            { no: 1, name: "gift_of_the_wild", kind: "enum", T: () => ["proto.TristateEffect", TristateEffect] },
            { no: 2, name: "power_word_fortitude", kind: "enum", T: () => ["proto.TristateEffect", TristateEffect] },
            { no: 3, name: "commanding_shout", kind: "enum", T: () => ["proto.TristateEffect", TristateEffect] },
            { no: 4, name: "blood_pact", kind: "enum", T: () => ["proto.TristateEffect", TristateEffect] },
            { no: 5, name: "horn_of_winter", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 6, name: "strength_of_earth_totem", kind: "enum", T: () => ["proto.TristateEffect", TristateEffect] },
            { no: 7, name: "arcane_brilliance", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 8, name: "fel_intelligence", kind: "enum", T: () => ["proto.TristateEffect", TristateEffect] },
            { no: 9, name: "divine_spirit", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 10, name: "battle_shout", kind: "enum", T: () => ["proto.TristateEffect", TristateEffect] },
            { no: 11, name: "trueshot_aura", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 12, name: "unleashed_rage", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 13, name: "abominations_might", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 14, name: "leader_of_the_pack", kind: "enum", T: () => ["proto.TristateEffect", TristateEffect] },
            { no: 15, name: "rampage", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 16, name: "icy_talons", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 17, name: "windfury_totem", kind: "enum", T: () => ["proto.TristateEffect", TristateEffect] },
            { no: 18, name: "totem_of_wrath", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 19, name: "flametongue_totem", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 20, name: "demonic_pact", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 21, name: "swift_retribution", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 22, name: "moonkin_aura", kind: "enum", T: () => ["proto.TristateEffect", TristateEffect] },
            { no: 23, name: "elemental_oath", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 24, name: "wrath_of_air_totem", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 25, name: "ferocious_inspiration", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 26, name: "sanctified_retribution", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 27, name: "arcane_empowerment", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 28, name: "mana_spring_totem", kind: "enum", T: () => ["proto.TristateEffect", TristateEffect] },
            { no: 29, name: "bloodlust", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 30, name: "thorns", kind: "enum", T: () => ["proto.TristateEffect", TristateEffect] },
            { no: 31, name: "devotion_aura", kind: "enum", T: () => ["proto.TristateEffect", TristateEffect] },
            { no: 32, name: "retribution_aura", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 33, name: "shadow_protection", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 34, name: "drums_of_forgotten_kings", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 35, name: "drums_of_the_wild", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 36, name: "scroll_of_protection", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 37, name: "scroll_of_stamina", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 38, name: "scroll_of_strength", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 39, name: "scroll_of_agility", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 40, name: "scroll_of_intellect", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 41, name: "scroll_of_spirit", kind: "scalar", T: 8 /*ScalarType.BOOL*/ }
        ]);
    }
    create(value) {
        const message = { giftOfTheWild: 0, powerWordFortitude: 0, commandingShout: 0, bloodPact: 0, hornOfWinter: false, strengthOfEarthTotem: 0, arcaneBrilliance: false, felIntelligence: 0, divineSpirit: false, battleShout: 0, trueshotAura: false, unleashedRage: false, abominationsMight: false, leaderOfThePack: 0, rampage: false, icyTalons: false, windfuryTotem: 0, totemOfWrath: false, flametongueTotem: false, demonicPact: 0, swiftRetribution: false, moonkinAura: 0, elementalOath: false, wrathOfAirTotem: false, ferociousInspiration: false, sanctifiedRetribution: false, arcaneEmpowerment: false, manaSpringTotem: 0, bloodlust: false, thorns: 0, devotionAura: 0, retributionAura: false, shadowProtection: false, drumsOfForgottenKings: false, drumsOfTheWild: false, scrollOfProtection: false, scrollOfStamina: false, scrollOfStrength: false, scrollOfAgility: false, scrollOfIntellect: false, scrollOfSpirit: false };
        Object.defineProperty(message, MESSAGE_TYPE, { enumerable: false, value: this });
        if (value !== undefined)
            reflectionMergePartial(this, message, value);
        return message;
    }
    internalBinaryRead(reader, length, options, target) {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* proto.TristateEffect gift_of_the_wild */ 1:
                    message.giftOfTheWild = reader.int32();
                    break;
                case /* proto.TristateEffect power_word_fortitude */ 2:
                    message.powerWordFortitude = reader.int32();
                    break;
                case /* proto.TristateEffect commanding_shout */ 3:
                    message.commandingShout = reader.int32();
                    break;
                case /* proto.TristateEffect blood_pact */ 4:
                    message.bloodPact = reader.int32();
                    break;
                case /* bool horn_of_winter */ 5:
                    message.hornOfWinter = reader.bool();
                    break;
                case /* proto.TristateEffect strength_of_earth_totem */ 6:
                    message.strengthOfEarthTotem = reader.int32();
                    break;
                case /* bool arcane_brilliance */ 7:
                    message.arcaneBrilliance = reader.bool();
                    break;
                case /* proto.TristateEffect fel_intelligence */ 8:
                    message.felIntelligence = reader.int32();
                    break;
                case /* bool divine_spirit */ 9:
                    message.divineSpirit = reader.bool();
                    break;
                case /* proto.TristateEffect battle_shout */ 10:
                    message.battleShout = reader.int32();
                    break;
                case /* bool trueshot_aura */ 11:
                    message.trueshotAura = reader.bool();
                    break;
                case /* bool unleashed_rage */ 12:
                    message.unleashedRage = reader.bool();
                    break;
                case /* bool abominations_might */ 13:
                    message.abominationsMight = reader.bool();
                    break;
                case /* proto.TristateEffect leader_of_the_pack */ 14:
                    message.leaderOfThePack = reader.int32();
                    break;
                case /* bool rampage */ 15:
                    message.rampage = reader.bool();
                    break;
                case /* bool icy_talons */ 16:
                    message.icyTalons = reader.bool();
                    break;
                case /* proto.TristateEffect windfury_totem */ 17:
                    message.windfuryTotem = reader.int32();
                    break;
                case /* bool totem_of_wrath */ 18:
                    message.totemOfWrath = reader.bool();
                    break;
                case /* bool flametongue_totem */ 19:
                    message.flametongueTotem = reader.bool();
                    break;
                case /* int32 demonic_pact */ 20:
                    message.demonicPact = reader.int32();
                    break;
                case /* bool swift_retribution */ 21:
                    message.swiftRetribution = reader.bool();
                    break;
                case /* proto.TristateEffect moonkin_aura */ 22:
                    message.moonkinAura = reader.int32();
                    break;
                case /* bool elemental_oath */ 23:
                    message.elementalOath = reader.bool();
                    break;
                case /* bool wrath_of_air_totem */ 24:
                    message.wrathOfAirTotem = reader.bool();
                    break;
                case /* bool ferocious_inspiration */ 25:
                    message.ferociousInspiration = reader.bool();
                    break;
                case /* bool sanctified_retribution */ 26:
                    message.sanctifiedRetribution = reader.bool();
                    break;
                case /* bool arcane_empowerment */ 27:
                    message.arcaneEmpowerment = reader.bool();
                    break;
                case /* proto.TristateEffect mana_spring_totem */ 28:
                    message.manaSpringTotem = reader.int32();
                    break;
                case /* bool bloodlust */ 29:
                    message.bloodlust = reader.bool();
                    break;
                case /* proto.TristateEffect thorns */ 30:
                    message.thorns = reader.int32();
                    break;
                case /* proto.TristateEffect devotion_aura */ 31:
                    message.devotionAura = reader.int32();
                    break;
                case /* bool retribution_aura */ 32:
                    message.retributionAura = reader.bool();
                    break;
                case /* bool shadow_protection */ 33:
                    message.shadowProtection = reader.bool();
                    break;
                case /* bool drums_of_forgotten_kings */ 34:
                    message.drumsOfForgottenKings = reader.bool();
                    break;
                case /* bool drums_of_the_wild */ 35:
                    message.drumsOfTheWild = reader.bool();
                    break;
                case /* bool scroll_of_protection */ 36:
                    message.scrollOfProtection = reader.bool();
                    break;
                case /* bool scroll_of_stamina */ 37:
                    message.scrollOfStamina = reader.bool();
                    break;
                case /* bool scroll_of_strength */ 38:
                    message.scrollOfStrength = reader.bool();
                    break;
                case /* bool scroll_of_agility */ 39:
                    message.scrollOfAgility = reader.bool();
                    break;
                case /* bool scroll_of_intellect */ 40:
                    message.scrollOfIntellect = reader.bool();
                    break;
                case /* bool scroll_of_spirit */ 41:
                    message.scrollOfSpirit = reader.bool();
                    break;
                default:
                    let u = options.readUnknownField;
                    if (u === "throw")
                        throw new globalThis.Error(`Unknown field ${fieldNo} (wire type ${wireType}) for ${this.typeName}`);
                    let d = reader.skip(wireType);
                    if (u !== false)
                        (u === true ? UnknownFieldHandler.onRead : u)(this.typeName, message, fieldNo, wireType, d);
            }
        }
        return message;
    }
    internalBinaryWrite(message, writer, options) {
        /* proto.TristateEffect gift_of_the_wild = 1; */
        if (message.giftOfTheWild !== 0)
            writer.tag(1, WireType.Varint).int32(message.giftOfTheWild);
        /* proto.TristateEffect power_word_fortitude = 2; */
        if (message.powerWordFortitude !== 0)
            writer.tag(2, WireType.Varint).int32(message.powerWordFortitude);
        /* proto.TristateEffect commanding_shout = 3; */
        if (message.commandingShout !== 0)
            writer.tag(3, WireType.Varint).int32(message.commandingShout);
        /* proto.TristateEffect blood_pact = 4; */
        if (message.bloodPact !== 0)
            writer.tag(4, WireType.Varint).int32(message.bloodPact);
        /* bool horn_of_winter = 5; */
        if (message.hornOfWinter !== false)
            writer.tag(5, WireType.Varint).bool(message.hornOfWinter);
        /* proto.TristateEffect strength_of_earth_totem = 6; */
        if (message.strengthOfEarthTotem !== 0)
            writer.tag(6, WireType.Varint).int32(message.strengthOfEarthTotem);
        /* bool arcane_brilliance = 7; */
        if (message.arcaneBrilliance !== false)
            writer.tag(7, WireType.Varint).bool(message.arcaneBrilliance);
        /* proto.TristateEffect fel_intelligence = 8; */
        if (message.felIntelligence !== 0)
            writer.tag(8, WireType.Varint).int32(message.felIntelligence);
        /* bool divine_spirit = 9; */
        if (message.divineSpirit !== false)
            writer.tag(9, WireType.Varint).bool(message.divineSpirit);
        /* proto.TristateEffect battle_shout = 10; */
        if (message.battleShout !== 0)
            writer.tag(10, WireType.Varint).int32(message.battleShout);
        /* bool trueshot_aura = 11; */
        if (message.trueshotAura !== false)
            writer.tag(11, WireType.Varint).bool(message.trueshotAura);
        /* bool unleashed_rage = 12; */
        if (message.unleashedRage !== false)
            writer.tag(12, WireType.Varint).bool(message.unleashedRage);
        /* bool abominations_might = 13; */
        if (message.abominationsMight !== false)
            writer.tag(13, WireType.Varint).bool(message.abominationsMight);
        /* proto.TristateEffect leader_of_the_pack = 14; */
        if (message.leaderOfThePack !== 0)
            writer.tag(14, WireType.Varint).int32(message.leaderOfThePack);
        /* bool rampage = 15; */
        if (message.rampage !== false)
            writer.tag(15, WireType.Varint).bool(message.rampage);
        /* bool icy_talons = 16; */
        if (message.icyTalons !== false)
            writer.tag(16, WireType.Varint).bool(message.icyTalons);
        /* proto.TristateEffect windfury_totem = 17; */
        if (message.windfuryTotem !== 0)
            writer.tag(17, WireType.Varint).int32(message.windfuryTotem);
        /* bool totem_of_wrath = 18; */
        if (message.totemOfWrath !== false)
            writer.tag(18, WireType.Varint).bool(message.totemOfWrath);
        /* bool flametongue_totem = 19; */
        if (message.flametongueTotem !== false)
            writer.tag(19, WireType.Varint).bool(message.flametongueTotem);
        /* int32 demonic_pact = 20; */
        if (message.demonicPact !== 0)
            writer.tag(20, WireType.Varint).int32(message.demonicPact);
        /* bool swift_retribution = 21; */
        if (message.swiftRetribution !== false)
            writer.tag(21, WireType.Varint).bool(message.swiftRetribution);
        /* proto.TristateEffect moonkin_aura = 22; */
        if (message.moonkinAura !== 0)
            writer.tag(22, WireType.Varint).int32(message.moonkinAura);
        /* bool elemental_oath = 23; */
        if (message.elementalOath !== false)
            writer.tag(23, WireType.Varint).bool(message.elementalOath);
        /* bool wrath_of_air_totem = 24; */
        if (message.wrathOfAirTotem !== false)
            writer.tag(24, WireType.Varint).bool(message.wrathOfAirTotem);
        /* bool ferocious_inspiration = 25; */
        if (message.ferociousInspiration !== false)
            writer.tag(25, WireType.Varint).bool(message.ferociousInspiration);
        /* bool sanctified_retribution = 26; */
        if (message.sanctifiedRetribution !== false)
            writer.tag(26, WireType.Varint).bool(message.sanctifiedRetribution);
        /* bool arcane_empowerment = 27; */
        if (message.arcaneEmpowerment !== false)
            writer.tag(27, WireType.Varint).bool(message.arcaneEmpowerment);
        /* proto.TristateEffect mana_spring_totem = 28; */
        if (message.manaSpringTotem !== 0)
            writer.tag(28, WireType.Varint).int32(message.manaSpringTotem);
        /* bool bloodlust = 29; */
        if (message.bloodlust !== false)
            writer.tag(29, WireType.Varint).bool(message.bloodlust);
        /* proto.TristateEffect thorns = 30; */
        if (message.thorns !== 0)
            writer.tag(30, WireType.Varint).int32(message.thorns);
        /* proto.TristateEffect devotion_aura = 31; */
        if (message.devotionAura !== 0)
            writer.tag(31, WireType.Varint).int32(message.devotionAura);
        /* bool retribution_aura = 32; */
        if (message.retributionAura !== false)
            writer.tag(32, WireType.Varint).bool(message.retributionAura);
        /* bool shadow_protection = 33; */
        if (message.shadowProtection !== false)
            writer.tag(33, WireType.Varint).bool(message.shadowProtection);
        /* bool drums_of_forgotten_kings = 34; */
        if (message.drumsOfForgottenKings !== false)
            writer.tag(34, WireType.Varint).bool(message.drumsOfForgottenKings);
        /* bool drums_of_the_wild = 35; */
        if (message.drumsOfTheWild !== false)
            writer.tag(35, WireType.Varint).bool(message.drumsOfTheWild);
        /* bool scroll_of_protection = 36; */
        if (message.scrollOfProtection !== false)
            writer.tag(36, WireType.Varint).bool(message.scrollOfProtection);
        /* bool scroll_of_stamina = 37; */
        if (message.scrollOfStamina !== false)
            writer.tag(37, WireType.Varint).bool(message.scrollOfStamina);
        /* bool scroll_of_strength = 38; */
        if (message.scrollOfStrength !== false)
            writer.tag(38, WireType.Varint).bool(message.scrollOfStrength);
        /* bool scroll_of_agility = 39; */
        if (message.scrollOfAgility !== false)
            writer.tag(39, WireType.Varint).bool(message.scrollOfAgility);
        /* bool scroll_of_intellect = 40; */
        if (message.scrollOfIntellect !== false)
            writer.tag(40, WireType.Varint).bool(message.scrollOfIntellect);
        /* bool scroll_of_spirit = 41; */
        if (message.scrollOfSpirit !== false)
            writer.tag(41, WireType.Varint).bool(message.scrollOfSpirit);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message proto.RaidBuffs
 */
export const RaidBuffs = new RaidBuffs$Type();
// @generated message type with reflection information, may provide speed optimized methods
class PartyBuffs$Type extends MessageType {
    constructor() {
        super("proto.PartyBuffs", [
            { no: 1, name: "atiesh_mage", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 2, name: "atiesh_warlock", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 3, name: "braided_eternium_chain", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 4, name: "eye_of_the_night", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 5, name: "chain_of_the_twilight_owl", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 6, name: "mana_tide_totems", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 7, name: "heroic_presence", kind: "scalar", T: 8 /*ScalarType.BOOL*/ }
        ]);
    }
    create(value) {
        const message = { atieshMage: 0, atieshWarlock: 0, braidedEterniumChain: false, eyeOfTheNight: false, chainOfTheTwilightOwl: false, manaTideTotems: 0, heroicPresence: false };
        Object.defineProperty(message, MESSAGE_TYPE, { enumerable: false, value: this });
        if (value !== undefined)
            reflectionMergePartial(this, message, value);
        return message;
    }
    internalBinaryRead(reader, length, options, target) {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* int32 atiesh_mage */ 1:
                    message.atieshMage = reader.int32();
                    break;
                case /* int32 atiesh_warlock */ 2:
                    message.atieshWarlock = reader.int32();
                    break;
                case /* bool braided_eternium_chain */ 3:
                    message.braidedEterniumChain = reader.bool();
                    break;
                case /* bool eye_of_the_night */ 4:
                    message.eyeOfTheNight = reader.bool();
                    break;
                case /* bool chain_of_the_twilight_owl */ 5:
                    message.chainOfTheTwilightOwl = reader.bool();
                    break;
                case /* int32 mana_tide_totems */ 6:
                    message.manaTideTotems = reader.int32();
                    break;
                case /* bool heroic_presence */ 7:
                    message.heroicPresence = reader.bool();
                    break;
                default:
                    let u = options.readUnknownField;
                    if (u === "throw")
                        throw new globalThis.Error(`Unknown field ${fieldNo} (wire type ${wireType}) for ${this.typeName}`);
                    let d = reader.skip(wireType);
                    if (u !== false)
                        (u === true ? UnknownFieldHandler.onRead : u)(this.typeName, message, fieldNo, wireType, d);
            }
        }
        return message;
    }
    internalBinaryWrite(message, writer, options) {
        /* int32 atiesh_mage = 1; */
        if (message.atieshMage !== 0)
            writer.tag(1, WireType.Varint).int32(message.atieshMage);
        /* int32 atiesh_warlock = 2; */
        if (message.atieshWarlock !== 0)
            writer.tag(2, WireType.Varint).int32(message.atieshWarlock);
        /* bool braided_eternium_chain = 3; */
        if (message.braidedEterniumChain !== false)
            writer.tag(3, WireType.Varint).bool(message.braidedEterniumChain);
        /* bool eye_of_the_night = 4; */
        if (message.eyeOfTheNight !== false)
            writer.tag(4, WireType.Varint).bool(message.eyeOfTheNight);
        /* bool chain_of_the_twilight_owl = 5; */
        if (message.chainOfTheTwilightOwl !== false)
            writer.tag(5, WireType.Varint).bool(message.chainOfTheTwilightOwl);
        /* int32 mana_tide_totems = 6; */
        if (message.manaTideTotems !== 0)
            writer.tag(6, WireType.Varint).int32(message.manaTideTotems);
        /* bool heroic_presence = 7; */
        if (message.heroicPresence !== false)
            writer.tag(7, WireType.Varint).bool(message.heroicPresence);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message proto.PartyBuffs
 */
export const PartyBuffs = new PartyBuffs$Type();
// @generated message type with reflection information, may provide speed optimized methods
class IndividualBuffs$Type extends MessageType {
    constructor() {
        super("proto.IndividualBuffs", [
            { no: 1, name: "blessing_of_kings", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 2, name: "blessing_of_wisdom", kind: "enum", T: () => ["proto.TristateEffect", TristateEffect] },
            { no: 3, name: "blessing_of_might", kind: "enum", T: () => ["proto.TristateEffect", TristateEffect] },
            { no: 4, name: "blessing_of_sanctuary", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 5, name: "vigilance", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 6, name: "renewed_hope", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 7, name: "hymn_of_hope", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 8, name: "hand_of_salvation", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 9, name: "rapture", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 10, name: "innervates", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 11, name: "power_infusions", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 12, name: "unholy_frenzy", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 13, name: "revitalize", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 14, name: "vampiric_touch", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 15, name: "hunting_party", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 16, name: "judgements_of_the_wise", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 17, name: "improved_soul_leech", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 18, name: "enduring_winter", kind: "scalar", T: 8 /*ScalarType.BOOL*/ }
        ]);
    }
    create(value) {
        const message = { blessingOfKings: false, blessingOfWisdom: 0, blessingOfMight: 0, blessingOfSanctuary: false, vigilance: false, renewedHope: false, hymnOfHope: 0, handOfSalvation: 0, rapture: 0, innervates: 0, powerInfusions: 0, unholyFrenzy: 0, revitalize: 0, vampiricTouch: false, huntingParty: false, judgementsOfTheWise: false, improvedSoulLeech: false, enduringWinter: false };
        Object.defineProperty(message, MESSAGE_TYPE, { enumerable: false, value: this });
        if (value !== undefined)
            reflectionMergePartial(this, message, value);
        return message;
    }
    internalBinaryRead(reader, length, options, target) {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* bool blessing_of_kings */ 1:
                    message.blessingOfKings = reader.bool();
                    break;
                case /* proto.TristateEffect blessing_of_wisdom */ 2:
                    message.blessingOfWisdom = reader.int32();
                    break;
                case /* proto.TristateEffect blessing_of_might */ 3:
                    message.blessingOfMight = reader.int32();
                    break;
                case /* bool blessing_of_sanctuary */ 4:
                    message.blessingOfSanctuary = reader.bool();
                    break;
                case /* bool vigilance */ 5:
                    message.vigilance = reader.bool();
                    break;
                case /* bool renewed_hope */ 6:
                    message.renewedHope = reader.bool();
                    break;
                case /* int32 hymn_of_hope */ 7:
                    message.hymnOfHope = reader.int32();
                    break;
                case /* int32 hand_of_salvation */ 8:
                    message.handOfSalvation = reader.int32();
                    break;
                case /* int32 rapture */ 9:
                    message.rapture = reader.int32();
                    break;
                case /* int32 innervates */ 10:
                    message.innervates = reader.int32();
                    break;
                case /* int32 power_infusions */ 11:
                    message.powerInfusions = reader.int32();
                    break;
                case /* int32 unholy_frenzy */ 12:
                    message.unholyFrenzy = reader.int32();
                    break;
                case /* int32 revitalize */ 13:
                    message.revitalize = reader.int32();
                    break;
                case /* bool vampiric_touch */ 14:
                    message.vampiricTouch = reader.bool();
                    break;
                case /* bool hunting_party */ 15:
                    message.huntingParty = reader.bool();
                    break;
                case /* bool judgements_of_the_wise */ 16:
                    message.judgementsOfTheWise = reader.bool();
                    break;
                case /* bool improved_soul_leech */ 17:
                    message.improvedSoulLeech = reader.bool();
                    break;
                case /* bool enduring_winter */ 18:
                    message.enduringWinter = reader.bool();
                    break;
                default:
                    let u = options.readUnknownField;
                    if (u === "throw")
                        throw new globalThis.Error(`Unknown field ${fieldNo} (wire type ${wireType}) for ${this.typeName}`);
                    let d = reader.skip(wireType);
                    if (u !== false)
                        (u === true ? UnknownFieldHandler.onRead : u)(this.typeName, message, fieldNo, wireType, d);
            }
        }
        return message;
    }
    internalBinaryWrite(message, writer, options) {
        /* bool blessing_of_kings = 1; */
        if (message.blessingOfKings !== false)
            writer.tag(1, WireType.Varint).bool(message.blessingOfKings);
        /* proto.TristateEffect blessing_of_wisdom = 2; */
        if (message.blessingOfWisdom !== 0)
            writer.tag(2, WireType.Varint).int32(message.blessingOfWisdom);
        /* proto.TristateEffect blessing_of_might = 3; */
        if (message.blessingOfMight !== 0)
            writer.tag(3, WireType.Varint).int32(message.blessingOfMight);
        /* bool blessing_of_sanctuary = 4; */
        if (message.blessingOfSanctuary !== false)
            writer.tag(4, WireType.Varint).bool(message.blessingOfSanctuary);
        /* bool vigilance = 5; */
        if (message.vigilance !== false)
            writer.tag(5, WireType.Varint).bool(message.vigilance);
        /* bool renewed_hope = 6; */
        if (message.renewedHope !== false)
            writer.tag(6, WireType.Varint).bool(message.renewedHope);
        /* int32 hymn_of_hope = 7; */
        if (message.hymnOfHope !== 0)
            writer.tag(7, WireType.Varint).int32(message.hymnOfHope);
        /* int32 hand_of_salvation = 8; */
        if (message.handOfSalvation !== 0)
            writer.tag(8, WireType.Varint).int32(message.handOfSalvation);
        /* int32 rapture = 9; */
        if (message.rapture !== 0)
            writer.tag(9, WireType.Varint).int32(message.rapture);
        /* int32 innervates = 10; */
        if (message.innervates !== 0)
            writer.tag(10, WireType.Varint).int32(message.innervates);
        /* int32 power_infusions = 11; */
        if (message.powerInfusions !== 0)
            writer.tag(11, WireType.Varint).int32(message.powerInfusions);
        /* int32 unholy_frenzy = 12; */
        if (message.unholyFrenzy !== 0)
            writer.tag(12, WireType.Varint).int32(message.unholyFrenzy);
        /* int32 revitalize = 13; */
        if (message.revitalize !== 0)
            writer.tag(13, WireType.Varint).int32(message.revitalize);
        /* bool vampiric_touch = 14; */
        if (message.vampiricTouch !== false)
            writer.tag(14, WireType.Varint).bool(message.vampiricTouch);
        /* bool hunting_party = 15; */
        if (message.huntingParty !== false)
            writer.tag(15, WireType.Varint).bool(message.huntingParty);
        /* bool judgements_of_the_wise = 16; */
        if (message.judgementsOfTheWise !== false)
            writer.tag(16, WireType.Varint).bool(message.judgementsOfTheWise);
        /* bool improved_soul_leech = 17; */
        if (message.improvedSoulLeech !== false)
            writer.tag(17, WireType.Varint).bool(message.improvedSoulLeech);
        /* bool enduring_winter = 18; */
        if (message.enduringWinter !== false)
            writer.tag(18, WireType.Varint).bool(message.enduringWinter);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message proto.IndividualBuffs
 */
export const IndividualBuffs = new IndividualBuffs$Type();
// @generated message type with reflection information, may provide speed optimized methods
class Consumes$Type extends MessageType {
    constructor() {
        super("proto.Consumes", [
            { no: 1, name: "flask", kind: "enum", T: () => ["proto.Flask", Flask] },
            { no: 2, name: "battle_elixir", kind: "enum", T: () => ["proto.BattleElixir", BattleElixir] },
            { no: 3, name: "guardian_elixir", kind: "enum", T: () => ["proto.GuardianElixir", GuardianElixir] },
            { no: 4, name: "main_hand_imbue", kind: "enum", T: () => ["proto.WeaponImbue", WeaponImbue] },
            { no: 5, name: "off_hand_imbue", kind: "enum", T: () => ["proto.WeaponImbue", WeaponImbue] },
            { no: 6, name: "food", kind: "enum", T: () => ["proto.Food", Food] },
            { no: 7, name: "pet_food", kind: "enum", T: () => ["proto.PetFood", PetFood] },
            { no: 8, name: "pet_scroll_of_agility", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 9, name: "pet_scroll_of_strength", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 10, name: "default_potion", kind: "enum", T: () => ["proto.Potions", Potions] },
            { no: 11, name: "prepop_potion", kind: "enum", T: () => ["proto.Potions", Potions] },
            { no: 12, name: "default_conjured", kind: "enum", T: () => ["proto.Conjured", Conjured] },
            { no: 13, name: "starting_conjured", kind: "enum", T: () => ["proto.Conjured", Conjured] },
            { no: 14, name: "num_starting_conjured", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 15, name: "super_sapper", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 16, name: "goblin_sapper", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 17, name: "filler_explosive", kind: "enum", T: () => ["proto.Explosive", Explosive] }
        ]);
    }
    create(value) {
        const message = { flask: 0, battleElixir: 0, guardianElixir: 0, mainHandImbue: 0, offHandImbue: 0, food: 0, petFood: 0, petScrollOfAgility: 0, petScrollOfStrength: 0, defaultPotion: 0, prepopPotion: 0, defaultConjured: 0, startingConjured: 0, numStartingConjured: 0, superSapper: false, goblinSapper: false, fillerExplosive: 0 };
        Object.defineProperty(message, MESSAGE_TYPE, { enumerable: false, value: this });
        if (value !== undefined)
            reflectionMergePartial(this, message, value);
        return message;
    }
    internalBinaryRead(reader, length, options, target) {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* proto.Flask flask */ 1:
                    message.flask = reader.int32();
                    break;
                case /* proto.BattleElixir battle_elixir */ 2:
                    message.battleElixir = reader.int32();
                    break;
                case /* proto.GuardianElixir guardian_elixir */ 3:
                    message.guardianElixir = reader.int32();
                    break;
                case /* proto.WeaponImbue main_hand_imbue */ 4:
                    message.mainHandImbue = reader.int32();
                    break;
                case /* proto.WeaponImbue off_hand_imbue */ 5:
                    message.offHandImbue = reader.int32();
                    break;
                case /* proto.Food food */ 6:
                    message.food = reader.int32();
                    break;
                case /* proto.PetFood pet_food */ 7:
                    message.petFood = reader.int32();
                    break;
                case /* int32 pet_scroll_of_agility */ 8:
                    message.petScrollOfAgility = reader.int32();
                    break;
                case /* int32 pet_scroll_of_strength */ 9:
                    message.petScrollOfStrength = reader.int32();
                    break;
                case /* proto.Potions default_potion */ 10:
                    message.defaultPotion = reader.int32();
                    break;
                case /* proto.Potions prepop_potion */ 11:
                    message.prepopPotion = reader.int32();
                    break;
                case /* proto.Conjured default_conjured */ 12:
                    message.defaultConjured = reader.int32();
                    break;
                case /* proto.Conjured starting_conjured */ 13:
                    message.startingConjured = reader.int32();
                    break;
                case /* int32 num_starting_conjured */ 14:
                    message.numStartingConjured = reader.int32();
                    break;
                case /* bool super_sapper */ 15:
                    message.superSapper = reader.bool();
                    break;
                case /* bool goblin_sapper */ 16:
                    message.goblinSapper = reader.bool();
                    break;
                case /* proto.Explosive filler_explosive */ 17:
                    message.fillerExplosive = reader.int32();
                    break;
                default:
                    let u = options.readUnknownField;
                    if (u === "throw")
                        throw new globalThis.Error(`Unknown field ${fieldNo} (wire type ${wireType}) for ${this.typeName}`);
                    let d = reader.skip(wireType);
                    if (u !== false)
                        (u === true ? UnknownFieldHandler.onRead : u)(this.typeName, message, fieldNo, wireType, d);
            }
        }
        return message;
    }
    internalBinaryWrite(message, writer, options) {
        /* proto.Flask flask = 1; */
        if (message.flask !== 0)
            writer.tag(1, WireType.Varint).int32(message.flask);
        /* proto.BattleElixir battle_elixir = 2; */
        if (message.battleElixir !== 0)
            writer.tag(2, WireType.Varint).int32(message.battleElixir);
        /* proto.GuardianElixir guardian_elixir = 3; */
        if (message.guardianElixir !== 0)
            writer.tag(3, WireType.Varint).int32(message.guardianElixir);
        /* proto.WeaponImbue main_hand_imbue = 4; */
        if (message.mainHandImbue !== 0)
            writer.tag(4, WireType.Varint).int32(message.mainHandImbue);
        /* proto.WeaponImbue off_hand_imbue = 5; */
        if (message.offHandImbue !== 0)
            writer.tag(5, WireType.Varint).int32(message.offHandImbue);
        /* proto.Food food = 6; */
        if (message.food !== 0)
            writer.tag(6, WireType.Varint).int32(message.food);
        /* proto.PetFood pet_food = 7; */
        if (message.petFood !== 0)
            writer.tag(7, WireType.Varint).int32(message.petFood);
        /* int32 pet_scroll_of_agility = 8; */
        if (message.petScrollOfAgility !== 0)
            writer.tag(8, WireType.Varint).int32(message.petScrollOfAgility);
        /* int32 pet_scroll_of_strength = 9; */
        if (message.petScrollOfStrength !== 0)
            writer.tag(9, WireType.Varint).int32(message.petScrollOfStrength);
        /* proto.Potions default_potion = 10; */
        if (message.defaultPotion !== 0)
            writer.tag(10, WireType.Varint).int32(message.defaultPotion);
        /* proto.Potions prepop_potion = 11; */
        if (message.prepopPotion !== 0)
            writer.tag(11, WireType.Varint).int32(message.prepopPotion);
        /* proto.Conjured default_conjured = 12; */
        if (message.defaultConjured !== 0)
            writer.tag(12, WireType.Varint).int32(message.defaultConjured);
        /* proto.Conjured starting_conjured = 13; */
        if (message.startingConjured !== 0)
            writer.tag(13, WireType.Varint).int32(message.startingConjured);
        /* int32 num_starting_conjured = 14; */
        if (message.numStartingConjured !== 0)
            writer.tag(14, WireType.Varint).int32(message.numStartingConjured);
        /* bool super_sapper = 15; */
        if (message.superSapper !== false)
            writer.tag(15, WireType.Varint).bool(message.superSapper);
        /* bool goblin_sapper = 16; */
        if (message.goblinSapper !== false)
            writer.tag(16, WireType.Varint).bool(message.goblinSapper);
        /* proto.Explosive filler_explosive = 17; */
        if (message.fillerExplosive !== 0)
            writer.tag(17, WireType.Varint).int32(message.fillerExplosive);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message proto.Consumes
 */
export const Consumes = new Consumes$Type();
// @generated message type with reflection information, may provide speed optimized methods
class Debuffs$Type extends MessageType {
    constructor() {
        super("proto.Debuffs", [
            { no: 1, name: "judgement_of_wisdom", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 2, name: "judgement_of_light", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 3, name: "misery", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 4, name: "faerie_fire", kind: "enum", T: () => ["proto.TristateEffect", TristateEffect] },
            { no: 5, name: "curse_of_elements", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 6, name: "ebon_plaguebringer", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 7, name: "earth_and_moon", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 8, name: "heart_of_the_crusader", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 9, name: "master_poisoner", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 10, name: "totem_of_wrath", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 11, name: "shadow_mastery", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 12, name: "improved_scorch", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 13, name: "winters_chill", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 14, name: "blood_frenzy", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 15, name: "savage_combat", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 16, name: "gift_of_arthas", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 17, name: "mangle", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 18, name: "trauma", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 19, name: "stampede", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 20, name: "expose_armor", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 21, name: "sunder_armor", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 22, name: "acid_spit", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 23, name: "curse_of_weakness", kind: "enum", T: () => ["proto.TristateEffect", TristateEffect] },
            { no: 24, name: "sting", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 25, name: "demoralizing_roar", kind: "enum", T: () => ["proto.TristateEffect", TristateEffect] },
            { no: 26, name: "demoralizing_shout", kind: "enum", T: () => ["proto.TristateEffect", TristateEffect] },
            { no: 27, name: "thunder_clap", kind: "enum", T: () => ["proto.TristateEffect", TristateEffect] },
            { no: 28, name: "frost_fever", kind: "enum", T: () => ["proto.TristateEffect", TristateEffect] },
            { no: 29, name: "infected_wounds", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 30, name: "judgements_of_the_just", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 31, name: "insect_swarm", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 32, name: "scorpid_sting", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 33, name: "shadow_embrace", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 34, name: "screech", kind: "scalar", T: 8 /*ScalarType.BOOL*/ }
        ]);
    }
    create(value) {
        const message = { judgementOfWisdom: false, judgementOfLight: false, misery: false, faerieFire: 0, curseOfElements: false, ebonPlaguebringer: false, earthAndMoon: false, heartOfTheCrusader: false, masterPoisoner: false, totemOfWrath: false, shadowMastery: false, improvedScorch: false, wintersChill: false, bloodFrenzy: false, savageCombat: false, giftOfArthas: false, mangle: false, trauma: false, stampede: false, exposeArmor: false, sunderArmor: false, acidSpit: false, curseOfWeakness: 0, sting: false, demoralizingRoar: 0, demoralizingShout: 0, thunderClap: 0, frostFever: 0, infectedWounds: false, judgementsOfTheJust: false, insectSwarm: false, scorpidSting: false, shadowEmbrace: false, screech: false };
        Object.defineProperty(message, MESSAGE_TYPE, { enumerable: false, value: this });
        if (value !== undefined)
            reflectionMergePartial(this, message, value);
        return message;
    }
    internalBinaryRead(reader, length, options, target) {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* bool judgement_of_wisdom */ 1:
                    message.judgementOfWisdom = reader.bool();
                    break;
                case /* bool judgement_of_light */ 2:
                    message.judgementOfLight = reader.bool();
                    break;
                case /* bool misery */ 3:
                    message.misery = reader.bool();
                    break;
                case /* proto.TristateEffect faerie_fire */ 4:
                    message.faerieFire = reader.int32();
                    break;
                case /* bool curse_of_elements */ 5:
                    message.curseOfElements = reader.bool();
                    break;
                case /* bool ebon_plaguebringer */ 6:
                    message.ebonPlaguebringer = reader.bool();
                    break;
                case /* bool earth_and_moon */ 7:
                    message.earthAndMoon = reader.bool();
                    break;
                case /* bool heart_of_the_crusader */ 8:
                    message.heartOfTheCrusader = reader.bool();
                    break;
                case /* bool master_poisoner */ 9:
                    message.masterPoisoner = reader.bool();
                    break;
                case /* bool totem_of_wrath */ 10:
                    message.totemOfWrath = reader.bool();
                    break;
                case /* bool shadow_mastery */ 11:
                    message.shadowMastery = reader.bool();
                    break;
                case /* bool improved_scorch */ 12:
                    message.improvedScorch = reader.bool();
                    break;
                case /* bool winters_chill */ 13:
                    message.wintersChill = reader.bool();
                    break;
                case /* bool blood_frenzy */ 14:
                    message.bloodFrenzy = reader.bool();
                    break;
                case /* bool savage_combat */ 15:
                    message.savageCombat = reader.bool();
                    break;
                case /* bool gift_of_arthas */ 16:
                    message.giftOfArthas = reader.bool();
                    break;
                case /* bool mangle */ 17:
                    message.mangle = reader.bool();
                    break;
                case /* bool trauma */ 18:
                    message.trauma = reader.bool();
                    break;
                case /* bool stampede */ 19:
                    message.stampede = reader.bool();
                    break;
                case /* bool expose_armor */ 20:
                    message.exposeArmor = reader.bool();
                    break;
                case /* bool sunder_armor */ 21:
                    message.sunderArmor = reader.bool();
                    break;
                case /* bool acid_spit */ 22:
                    message.acidSpit = reader.bool();
                    break;
                case /* proto.TristateEffect curse_of_weakness */ 23:
                    message.curseOfWeakness = reader.int32();
                    break;
                case /* bool sting */ 24:
                    message.sting = reader.bool();
                    break;
                case /* proto.TristateEffect demoralizing_roar */ 25:
                    message.demoralizingRoar = reader.int32();
                    break;
                case /* proto.TristateEffect demoralizing_shout */ 26:
                    message.demoralizingShout = reader.int32();
                    break;
                case /* proto.TristateEffect thunder_clap */ 27:
                    message.thunderClap = reader.int32();
                    break;
                case /* proto.TristateEffect frost_fever */ 28:
                    message.frostFever = reader.int32();
                    break;
                case /* bool infected_wounds */ 29:
                    message.infectedWounds = reader.bool();
                    break;
                case /* bool judgements_of_the_just */ 30:
                    message.judgementsOfTheJust = reader.bool();
                    break;
                case /* bool insect_swarm */ 31:
                    message.insectSwarm = reader.bool();
                    break;
                case /* bool scorpid_sting */ 32:
                    message.scorpidSting = reader.bool();
                    break;
                case /* bool shadow_embrace */ 33:
                    message.shadowEmbrace = reader.bool();
                    break;
                case /* bool screech */ 34:
                    message.screech = reader.bool();
                    break;
                default:
                    let u = options.readUnknownField;
                    if (u === "throw")
                        throw new globalThis.Error(`Unknown field ${fieldNo} (wire type ${wireType}) for ${this.typeName}`);
                    let d = reader.skip(wireType);
                    if (u !== false)
                        (u === true ? UnknownFieldHandler.onRead : u)(this.typeName, message, fieldNo, wireType, d);
            }
        }
        return message;
    }
    internalBinaryWrite(message, writer, options) {
        /* bool judgement_of_wisdom = 1; */
        if (message.judgementOfWisdom !== false)
            writer.tag(1, WireType.Varint).bool(message.judgementOfWisdom);
        /* bool judgement_of_light = 2; */
        if (message.judgementOfLight !== false)
            writer.tag(2, WireType.Varint).bool(message.judgementOfLight);
        /* bool misery = 3; */
        if (message.misery !== false)
            writer.tag(3, WireType.Varint).bool(message.misery);
        /* proto.TristateEffect faerie_fire = 4; */
        if (message.faerieFire !== 0)
            writer.tag(4, WireType.Varint).int32(message.faerieFire);
        /* bool curse_of_elements = 5; */
        if (message.curseOfElements !== false)
            writer.tag(5, WireType.Varint).bool(message.curseOfElements);
        /* bool ebon_plaguebringer = 6; */
        if (message.ebonPlaguebringer !== false)
            writer.tag(6, WireType.Varint).bool(message.ebonPlaguebringer);
        /* bool earth_and_moon = 7; */
        if (message.earthAndMoon !== false)
            writer.tag(7, WireType.Varint).bool(message.earthAndMoon);
        /* bool heart_of_the_crusader = 8; */
        if (message.heartOfTheCrusader !== false)
            writer.tag(8, WireType.Varint).bool(message.heartOfTheCrusader);
        /* bool master_poisoner = 9; */
        if (message.masterPoisoner !== false)
            writer.tag(9, WireType.Varint).bool(message.masterPoisoner);
        /* bool totem_of_wrath = 10; */
        if (message.totemOfWrath !== false)
            writer.tag(10, WireType.Varint).bool(message.totemOfWrath);
        /* bool shadow_mastery = 11; */
        if (message.shadowMastery !== false)
            writer.tag(11, WireType.Varint).bool(message.shadowMastery);
        /* bool improved_scorch = 12; */
        if (message.improvedScorch !== false)
            writer.tag(12, WireType.Varint).bool(message.improvedScorch);
        /* bool winters_chill = 13; */
        if (message.wintersChill !== false)
            writer.tag(13, WireType.Varint).bool(message.wintersChill);
        /* bool blood_frenzy = 14; */
        if (message.bloodFrenzy !== false)
            writer.tag(14, WireType.Varint).bool(message.bloodFrenzy);
        /* bool savage_combat = 15; */
        if (message.savageCombat !== false)
            writer.tag(15, WireType.Varint).bool(message.savageCombat);
        /* bool gift_of_arthas = 16; */
        if (message.giftOfArthas !== false)
            writer.tag(16, WireType.Varint).bool(message.giftOfArthas);
        /* bool mangle = 17; */
        if (message.mangle !== false)
            writer.tag(17, WireType.Varint).bool(message.mangle);
        /* bool trauma = 18; */
        if (message.trauma !== false)
            writer.tag(18, WireType.Varint).bool(message.trauma);
        /* bool stampede = 19; */
        if (message.stampede !== false)
            writer.tag(19, WireType.Varint).bool(message.stampede);
        /* bool expose_armor = 20; */
        if (message.exposeArmor !== false)
            writer.tag(20, WireType.Varint).bool(message.exposeArmor);
        /* bool sunder_armor = 21; */
        if (message.sunderArmor !== false)
            writer.tag(21, WireType.Varint).bool(message.sunderArmor);
        /* bool acid_spit = 22; */
        if (message.acidSpit !== false)
            writer.tag(22, WireType.Varint).bool(message.acidSpit);
        /* proto.TristateEffect curse_of_weakness = 23; */
        if (message.curseOfWeakness !== 0)
            writer.tag(23, WireType.Varint).int32(message.curseOfWeakness);
        /* bool sting = 24; */
        if (message.sting !== false)
            writer.tag(24, WireType.Varint).bool(message.sting);
        /* proto.TristateEffect demoralizing_roar = 25; */
        if (message.demoralizingRoar !== 0)
            writer.tag(25, WireType.Varint).int32(message.demoralizingRoar);
        /* proto.TristateEffect demoralizing_shout = 26; */
        if (message.demoralizingShout !== 0)
            writer.tag(26, WireType.Varint).int32(message.demoralizingShout);
        /* proto.TristateEffect thunder_clap = 27; */
        if (message.thunderClap !== 0)
            writer.tag(27, WireType.Varint).int32(message.thunderClap);
        /* proto.TristateEffect frost_fever = 28; */
        if (message.frostFever !== 0)
            writer.tag(28, WireType.Varint).int32(message.frostFever);
        /* bool infected_wounds = 29; */
        if (message.infectedWounds !== false)
            writer.tag(29, WireType.Varint).bool(message.infectedWounds);
        /* bool judgements_of_the_just = 30; */
        if (message.judgementsOfTheJust !== false)
            writer.tag(30, WireType.Varint).bool(message.judgementsOfTheJust);
        /* bool insect_swarm = 31; */
        if (message.insectSwarm !== false)
            writer.tag(31, WireType.Varint).bool(message.insectSwarm);
        /* bool scorpid_sting = 32; */
        if (message.scorpidSting !== false)
            writer.tag(32, WireType.Varint).bool(message.scorpidSting);
        /* bool shadow_embrace = 33; */
        if (message.shadowEmbrace !== false)
            writer.tag(33, WireType.Varint).bool(message.shadowEmbrace);
        /* bool screech = 34; */
        if (message.screech !== false)
            writer.tag(34, WireType.Varint).bool(message.screech);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message proto.Debuffs
 */
export const Debuffs = new Debuffs$Type();
// @generated message type with reflection information, may provide speed optimized methods
class Target$Type extends MessageType {
    constructor() {
        super("proto.Target", [
            { no: 14, name: "id", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 15, name: "name", kind: "scalar", T: 9 /*ScalarType.STRING*/ },
            { no: 4, name: "level", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 3, name: "mob_type", kind: "enum", T: () => ["proto.MobType", MobType] },
            { no: 5, name: "stats", kind: "scalar", repeat: 1 /*RepeatType.PACKED*/, T: 1 /*ScalarType.DOUBLE*/ },
            { no: 7, name: "min_base_damage", kind: "scalar", T: 1 /*ScalarType.DOUBLE*/ },
            { no: 8, name: "swing_speed", kind: "scalar", T: 1 /*ScalarType.DOUBLE*/ },
            { no: 9, name: "dual_wield", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 10, name: "dual_wield_penalty", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 12, name: "parry_haste", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 16, name: "suppress_dodge", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 13, name: "spell_school", kind: "enum", T: () => ["proto.SpellSchool", SpellSchool] },
            { no: 6, name: "tank_index", kind: "scalar", T: 5 /*ScalarType.INT32*/ }
        ]);
    }
    create(value) {
        const message = { id: 0, name: "", level: 0, mobType: 0, stats: [], minBaseDamage: 0, swingSpeed: 0, dualWield: false, dualWieldPenalty: false, parryHaste: false, suppressDodge: false, spellSchool: 0, tankIndex: 0 };
        Object.defineProperty(message, MESSAGE_TYPE, { enumerable: false, value: this });
        if (value !== undefined)
            reflectionMergePartial(this, message, value);
        return message;
    }
    internalBinaryRead(reader, length, options, target) {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* int32 id */ 14:
                    message.id = reader.int32();
                    break;
                case /* string name */ 15:
                    message.name = reader.string();
                    break;
                case /* int32 level */ 4:
                    message.level = reader.int32();
                    break;
                case /* proto.MobType mob_type */ 3:
                    message.mobType = reader.int32();
                    break;
                case /* repeated double stats */ 5:
                    if (wireType === WireType.LengthDelimited)
                        for (let e = reader.int32() + reader.pos; reader.pos < e;)
                            message.stats.push(reader.double());
                    else
                        message.stats.push(reader.double());
                    break;
                case /* double min_base_damage */ 7:
                    message.minBaseDamage = reader.double();
                    break;
                case /* double swing_speed */ 8:
                    message.swingSpeed = reader.double();
                    break;
                case /* bool dual_wield */ 9:
                    message.dualWield = reader.bool();
                    break;
                case /* bool dual_wield_penalty */ 10:
                    message.dualWieldPenalty = reader.bool();
                    break;
                case /* bool parry_haste */ 12:
                    message.parryHaste = reader.bool();
                    break;
                case /* bool suppress_dodge */ 16:
                    message.suppressDodge = reader.bool();
                    break;
                case /* proto.SpellSchool spell_school */ 13:
                    message.spellSchool = reader.int32();
                    break;
                case /* int32 tank_index */ 6:
                    message.tankIndex = reader.int32();
                    break;
                default:
                    let u = options.readUnknownField;
                    if (u === "throw")
                        throw new globalThis.Error(`Unknown field ${fieldNo} (wire type ${wireType}) for ${this.typeName}`);
                    let d = reader.skip(wireType);
                    if (u !== false)
                        (u === true ? UnknownFieldHandler.onRead : u)(this.typeName, message, fieldNo, wireType, d);
            }
        }
        return message;
    }
    internalBinaryWrite(message, writer, options) {
        /* int32 id = 14; */
        if (message.id !== 0)
            writer.tag(14, WireType.Varint).int32(message.id);
        /* string name = 15; */
        if (message.name !== "")
            writer.tag(15, WireType.LengthDelimited).string(message.name);
        /* int32 level = 4; */
        if (message.level !== 0)
            writer.tag(4, WireType.Varint).int32(message.level);
        /* proto.MobType mob_type = 3; */
        if (message.mobType !== 0)
            writer.tag(3, WireType.Varint).int32(message.mobType);
        /* repeated double stats = 5; */
        if (message.stats.length) {
            writer.tag(5, WireType.LengthDelimited).fork();
            for (let i = 0; i < message.stats.length; i++)
                writer.double(message.stats[i]);
            writer.join();
        }
        /* double min_base_damage = 7; */
        if (message.minBaseDamage !== 0)
            writer.tag(7, WireType.Bit64).double(message.minBaseDamage);
        /* double swing_speed = 8; */
        if (message.swingSpeed !== 0)
            writer.tag(8, WireType.Bit64).double(message.swingSpeed);
        /* bool dual_wield = 9; */
        if (message.dualWield !== false)
            writer.tag(9, WireType.Varint).bool(message.dualWield);
        /* bool dual_wield_penalty = 10; */
        if (message.dualWieldPenalty !== false)
            writer.tag(10, WireType.Varint).bool(message.dualWieldPenalty);
        /* bool parry_haste = 12; */
        if (message.parryHaste !== false)
            writer.tag(12, WireType.Varint).bool(message.parryHaste);
        /* bool suppress_dodge = 16; */
        if (message.suppressDodge !== false)
            writer.tag(16, WireType.Varint).bool(message.suppressDodge);
        /* proto.SpellSchool spell_school = 13; */
        if (message.spellSchool !== 0)
            writer.tag(13, WireType.Varint).int32(message.spellSchool);
        /* int32 tank_index = 6; */
        if (message.tankIndex !== 0)
            writer.tag(6, WireType.Varint).int32(message.tankIndex);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message proto.Target
 */
export const Target = new Target$Type();
// @generated message type with reflection information, may provide speed optimized methods
class Encounter$Type extends MessageType {
    constructor() {
        super("proto.Encounter", [
            { no: 1, name: "duration", kind: "scalar", T: 1 /*ScalarType.DOUBLE*/ },
            { no: 2, name: "duration_variation", kind: "scalar", T: 1 /*ScalarType.DOUBLE*/ },
            { no: 3, name: "execute_proportion_20", kind: "scalar", T: 1 /*ScalarType.DOUBLE*/ },
            { no: 4, name: "execute_proportion_35", kind: "scalar", T: 1 /*ScalarType.DOUBLE*/ },
            { no: 5, name: "use_health", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 6, name: "targets", kind: "message", repeat: 1 /*RepeatType.PACKED*/, T: () => Target }
        ]);
    }
    create(value) {
        const message = { duration: 0, durationVariation: 0, executeProportion20: 0, executeProportion35: 0, useHealth: false, targets: [] };
        Object.defineProperty(message, MESSAGE_TYPE, { enumerable: false, value: this });
        if (value !== undefined)
            reflectionMergePartial(this, message, value);
        return message;
    }
    internalBinaryRead(reader, length, options, target) {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* double duration */ 1:
                    message.duration = reader.double();
                    break;
                case /* double duration_variation */ 2:
                    message.durationVariation = reader.double();
                    break;
                case /* double execute_proportion_20 */ 3:
                    message.executeProportion20 = reader.double();
                    break;
                case /* double execute_proportion_35 */ 4:
                    message.executeProportion35 = reader.double();
                    break;
                case /* bool use_health */ 5:
                    message.useHealth = reader.bool();
                    break;
                case /* repeated proto.Target targets */ 6:
                    message.targets.push(Target.internalBinaryRead(reader, reader.uint32(), options));
                    break;
                default:
                    let u = options.readUnknownField;
                    if (u === "throw")
                        throw new globalThis.Error(`Unknown field ${fieldNo} (wire type ${wireType}) for ${this.typeName}`);
                    let d = reader.skip(wireType);
                    if (u !== false)
                        (u === true ? UnknownFieldHandler.onRead : u)(this.typeName, message, fieldNo, wireType, d);
            }
        }
        return message;
    }
    internalBinaryWrite(message, writer, options) {
        /* double duration = 1; */
        if (message.duration !== 0)
            writer.tag(1, WireType.Bit64).double(message.duration);
        /* double duration_variation = 2; */
        if (message.durationVariation !== 0)
            writer.tag(2, WireType.Bit64).double(message.durationVariation);
        /* double execute_proportion_20 = 3; */
        if (message.executeProportion20 !== 0)
            writer.tag(3, WireType.Bit64).double(message.executeProportion20);
        /* double execute_proportion_35 = 4; */
        if (message.executeProportion35 !== 0)
            writer.tag(4, WireType.Bit64).double(message.executeProportion35);
        /* bool use_health = 5; */
        if (message.useHealth !== false)
            writer.tag(5, WireType.Varint).bool(message.useHealth);
        /* repeated proto.Target targets = 6; */
        for (let i = 0; i < message.targets.length; i++)
            Target.internalBinaryWrite(message.targets[i], writer.tag(6, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message proto.Encounter
 */
export const Encounter = new Encounter$Type();
// @generated message type with reflection information, may provide speed optimized methods
class ItemSpec$Type extends MessageType {
    constructor() {
        super("proto.ItemSpec", [
            { no: 2, name: "id", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 3, name: "enchant", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 4, name: "gems", kind: "scalar", repeat: 1 /*RepeatType.PACKED*/, T: 5 /*ScalarType.INT32*/ }
        ]);
    }
    create(value) {
        const message = { id: 0, enchant: 0, gems: [] };
        Object.defineProperty(message, MESSAGE_TYPE, { enumerable: false, value: this });
        if (value !== undefined)
            reflectionMergePartial(this, message, value);
        return message;
    }
    internalBinaryRead(reader, length, options, target) {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* int32 id */ 2:
                    message.id = reader.int32();
                    break;
                case /* int32 enchant */ 3:
                    message.enchant = reader.int32();
                    break;
                case /* repeated int32 gems */ 4:
                    if (wireType === WireType.LengthDelimited)
                        for (let e = reader.int32() + reader.pos; reader.pos < e;)
                            message.gems.push(reader.int32());
                    else
                        message.gems.push(reader.int32());
                    break;
                default:
                    let u = options.readUnknownField;
                    if (u === "throw")
                        throw new globalThis.Error(`Unknown field ${fieldNo} (wire type ${wireType}) for ${this.typeName}`);
                    let d = reader.skip(wireType);
                    if (u !== false)
                        (u === true ? UnknownFieldHandler.onRead : u)(this.typeName, message, fieldNo, wireType, d);
            }
        }
        return message;
    }
    internalBinaryWrite(message, writer, options) {
        /* int32 id = 2; */
        if (message.id !== 0)
            writer.tag(2, WireType.Varint).int32(message.id);
        /* int32 enchant = 3; */
        if (message.enchant !== 0)
            writer.tag(3, WireType.Varint).int32(message.enchant);
        /* repeated int32 gems = 4; */
        if (message.gems.length) {
            writer.tag(4, WireType.LengthDelimited).fork();
            for (let i = 0; i < message.gems.length; i++)
                writer.int32(message.gems[i]);
            writer.join();
        }
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message proto.ItemSpec
 */
export const ItemSpec = new ItemSpec$Type();
// @generated message type with reflection information, may provide speed optimized methods
class EquipmentSpec$Type extends MessageType {
    constructor() {
        super("proto.EquipmentSpec", [
            { no: 1, name: "items", kind: "message", repeat: 1 /*RepeatType.PACKED*/, T: () => ItemSpec }
        ]);
    }
    create(value) {
        const message = { items: [] };
        Object.defineProperty(message, MESSAGE_TYPE, { enumerable: false, value: this });
        if (value !== undefined)
            reflectionMergePartial(this, message, value);
        return message;
    }
    internalBinaryRead(reader, length, options, target) {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* repeated proto.ItemSpec items */ 1:
                    message.items.push(ItemSpec.internalBinaryRead(reader, reader.uint32(), options));
                    break;
                default:
                    let u = options.readUnknownField;
                    if (u === "throw")
                        throw new globalThis.Error(`Unknown field ${fieldNo} (wire type ${wireType}) for ${this.typeName}`);
                    let d = reader.skip(wireType);
                    if (u !== false)
                        (u === true ? UnknownFieldHandler.onRead : u)(this.typeName, message, fieldNo, wireType, d);
            }
        }
        return message;
    }
    internalBinaryWrite(message, writer, options) {
        /* repeated proto.ItemSpec items = 1; */
        for (let i = 0; i < message.items.length; i++)
            ItemSpec.internalBinaryWrite(message.items[i], writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message proto.EquipmentSpec
 */
export const EquipmentSpec = new EquipmentSpec$Type();
// @generated message type with reflection information, may provide speed optimized methods
class Item$Type extends MessageType {
    constructor() {
        super("proto.Item", [
            { no: 1, name: "id", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 16, name: "wowhead_id", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 2, name: "name", kind: "scalar", T: 9 /*ScalarType.STRING*/ },
            { no: 15, name: "class_allowlist", kind: "enum", repeat: 1 /*RepeatType.PACKED*/, T: () => ["proto.Class", Class] },
            { no: 3, name: "type", kind: "enum", T: () => ["proto.ItemType", ItemType] },
            { no: 4, name: "armor_type", kind: "enum", T: () => ["proto.ArmorType", ArmorType] },
            { no: 5, name: "weapon_type", kind: "enum", T: () => ["proto.WeaponType", WeaponType] },
            { no: 6, name: "hand_type", kind: "enum", T: () => ["proto.HandType", HandType] },
            { no: 7, name: "ranged_weapon_type", kind: "enum", T: () => ["proto.RangedWeaponType", RangedWeaponType] },
            { no: 8, name: "stats", kind: "scalar", repeat: 1 /*RepeatType.PACKED*/, T: 1 /*ScalarType.DOUBLE*/ },
            { no: 9, name: "gem_sockets", kind: "enum", repeat: 1 /*RepeatType.PACKED*/, T: () => ["proto.GemColor", GemColor] },
            { no: 10, name: "socketBonus", kind: "scalar", repeat: 1 /*RepeatType.PACKED*/, T: 1 /*ScalarType.DOUBLE*/ },
            { no: 17, name: "weapon_damage_min", kind: "scalar", T: 1 /*ScalarType.DOUBLE*/ },
            { no: 18, name: "weapon_damage_max", kind: "scalar", T: 1 /*ScalarType.DOUBLE*/ },
            { no: 19, name: "weapon_speed", kind: "scalar", T: 1 /*ScalarType.DOUBLE*/ },
            { no: 11, name: "phase", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 12, name: "quality", kind: "enum", T: () => ["proto.ItemQuality", ItemQuality] },
            { no: 13, name: "unique", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 20, name: "ilvl", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 21, name: "required_profession", kind: "enum", T: () => ["proto.Profession", Profession] }
        ]);
    }
    create(value) {
        const message = { id: 0, wowheadId: 0, name: "", classAllowlist: [], type: 0, armorType: 0, weaponType: 0, handType: 0, rangedWeaponType: 0, stats: [], gemSockets: [], socketBonus: [], weaponDamageMin: 0, weaponDamageMax: 0, weaponSpeed: 0, phase: 0, quality: 0, unique: false, ilvl: 0, requiredProfession: 0 };
        Object.defineProperty(message, MESSAGE_TYPE, { enumerable: false, value: this });
        if (value !== undefined)
            reflectionMergePartial(this, message, value);
        return message;
    }
    internalBinaryRead(reader, length, options, target) {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* int32 id */ 1:
                    message.id = reader.int32();
                    break;
                case /* int32 wowhead_id */ 16:
                    message.wowheadId = reader.int32();
                    break;
                case /* string name */ 2:
                    message.name = reader.string();
                    break;
                case /* repeated proto.Class class_allowlist */ 15:
                    if (wireType === WireType.LengthDelimited)
                        for (let e = reader.int32() + reader.pos; reader.pos < e;)
                            message.classAllowlist.push(reader.int32());
                    else
                        message.classAllowlist.push(reader.int32());
                    break;
                case /* proto.ItemType type */ 3:
                    message.type = reader.int32();
                    break;
                case /* proto.ArmorType armor_type */ 4:
                    message.armorType = reader.int32();
                    break;
                case /* proto.WeaponType weapon_type */ 5:
                    message.weaponType = reader.int32();
                    break;
                case /* proto.HandType hand_type */ 6:
                    message.handType = reader.int32();
                    break;
                case /* proto.RangedWeaponType ranged_weapon_type */ 7:
                    message.rangedWeaponType = reader.int32();
                    break;
                case /* repeated double stats */ 8:
                    if (wireType === WireType.LengthDelimited)
                        for (let e = reader.int32() + reader.pos; reader.pos < e;)
                            message.stats.push(reader.double());
                    else
                        message.stats.push(reader.double());
                    break;
                case /* repeated proto.GemColor gem_sockets */ 9:
                    if (wireType === WireType.LengthDelimited)
                        for (let e = reader.int32() + reader.pos; reader.pos < e;)
                            message.gemSockets.push(reader.int32());
                    else
                        message.gemSockets.push(reader.int32());
                    break;
                case /* repeated double socketBonus */ 10:
                    if (wireType === WireType.LengthDelimited)
                        for (let e = reader.int32() + reader.pos; reader.pos < e;)
                            message.socketBonus.push(reader.double());
                    else
                        message.socketBonus.push(reader.double());
                    break;
                case /* double weapon_damage_min */ 17:
                    message.weaponDamageMin = reader.double();
                    break;
                case /* double weapon_damage_max */ 18:
                    message.weaponDamageMax = reader.double();
                    break;
                case /* double weapon_speed */ 19:
                    message.weaponSpeed = reader.double();
                    break;
                case /* int32 phase */ 11:
                    message.phase = reader.int32();
                    break;
                case /* proto.ItemQuality quality */ 12:
                    message.quality = reader.int32();
                    break;
                case /* bool unique */ 13:
                    message.unique = reader.bool();
                    break;
                case /* int32 ilvl */ 20:
                    message.ilvl = reader.int32();
                    break;
                case /* proto.Profession required_profession */ 21:
                    message.requiredProfession = reader.int32();
                    break;
                default:
                    let u = options.readUnknownField;
                    if (u === "throw")
                        throw new globalThis.Error(`Unknown field ${fieldNo} (wire type ${wireType}) for ${this.typeName}`);
                    let d = reader.skip(wireType);
                    if (u !== false)
                        (u === true ? UnknownFieldHandler.onRead : u)(this.typeName, message, fieldNo, wireType, d);
            }
        }
        return message;
    }
    internalBinaryWrite(message, writer, options) {
        /* int32 id = 1; */
        if (message.id !== 0)
            writer.tag(1, WireType.Varint).int32(message.id);
        /* int32 wowhead_id = 16; */
        if (message.wowheadId !== 0)
            writer.tag(16, WireType.Varint).int32(message.wowheadId);
        /* string name = 2; */
        if (message.name !== "")
            writer.tag(2, WireType.LengthDelimited).string(message.name);
        /* repeated proto.Class class_allowlist = 15; */
        if (message.classAllowlist.length) {
            writer.tag(15, WireType.LengthDelimited).fork();
            for (let i = 0; i < message.classAllowlist.length; i++)
                writer.int32(message.classAllowlist[i]);
            writer.join();
        }
        /* proto.ItemType type = 3; */
        if (message.type !== 0)
            writer.tag(3, WireType.Varint).int32(message.type);
        /* proto.ArmorType armor_type = 4; */
        if (message.armorType !== 0)
            writer.tag(4, WireType.Varint).int32(message.armorType);
        /* proto.WeaponType weapon_type = 5; */
        if (message.weaponType !== 0)
            writer.tag(5, WireType.Varint).int32(message.weaponType);
        /* proto.HandType hand_type = 6; */
        if (message.handType !== 0)
            writer.tag(6, WireType.Varint).int32(message.handType);
        /* proto.RangedWeaponType ranged_weapon_type = 7; */
        if (message.rangedWeaponType !== 0)
            writer.tag(7, WireType.Varint).int32(message.rangedWeaponType);
        /* repeated double stats = 8; */
        if (message.stats.length) {
            writer.tag(8, WireType.LengthDelimited).fork();
            for (let i = 0; i < message.stats.length; i++)
                writer.double(message.stats[i]);
            writer.join();
        }
        /* repeated proto.GemColor gem_sockets = 9; */
        if (message.gemSockets.length) {
            writer.tag(9, WireType.LengthDelimited).fork();
            for (let i = 0; i < message.gemSockets.length; i++)
                writer.int32(message.gemSockets[i]);
            writer.join();
        }
        /* repeated double socketBonus = 10; */
        if (message.socketBonus.length) {
            writer.tag(10, WireType.LengthDelimited).fork();
            for (let i = 0; i < message.socketBonus.length; i++)
                writer.double(message.socketBonus[i]);
            writer.join();
        }
        /* double weapon_damage_min = 17; */
        if (message.weaponDamageMin !== 0)
            writer.tag(17, WireType.Bit64).double(message.weaponDamageMin);
        /* double weapon_damage_max = 18; */
        if (message.weaponDamageMax !== 0)
            writer.tag(18, WireType.Bit64).double(message.weaponDamageMax);
        /* double weapon_speed = 19; */
        if (message.weaponSpeed !== 0)
            writer.tag(19, WireType.Bit64).double(message.weaponSpeed);
        /* int32 phase = 11; */
        if (message.phase !== 0)
            writer.tag(11, WireType.Varint).int32(message.phase);
        /* proto.ItemQuality quality = 12; */
        if (message.quality !== 0)
            writer.tag(12, WireType.Varint).int32(message.quality);
        /* bool unique = 13; */
        if (message.unique !== false)
            writer.tag(13, WireType.Varint).bool(message.unique);
        /* int32 ilvl = 20; */
        if (message.ilvl !== 0)
            writer.tag(20, WireType.Varint).int32(message.ilvl);
        /* proto.Profession required_profession = 21; */
        if (message.requiredProfession !== 0)
            writer.tag(21, WireType.Varint).int32(message.requiredProfession);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message proto.Item
 */
export const Item = new Item$Type();
// @generated message type with reflection information, may provide speed optimized methods
class Enchant$Type extends MessageType {
    constructor() {
        super("proto.Enchant", [
            { no: 1, name: "id", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 2, name: "effect_id", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 3, name: "name", kind: "scalar", T: 9 /*ScalarType.STRING*/ },
            { no: 10, name: "is_spell_id", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 4, name: "type", kind: "enum", T: () => ["proto.ItemType", ItemType] },
            { no: 9, name: "enchant_type", kind: "enum", T: () => ["proto.EnchantType", EnchantType] },
            { no: 7, name: "stats", kind: "scalar", repeat: 1 /*RepeatType.PACKED*/, T: 1 /*ScalarType.DOUBLE*/ },
            { no: 8, name: "quality", kind: "enum", T: () => ["proto.ItemQuality", ItemQuality] },
            { no: 11, name: "phase", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 13, name: "required_profession", kind: "enum", T: () => ["proto.Profession", Profession] },
            { no: 12, name: "class_allowlist", kind: "enum", repeat: 1 /*RepeatType.PACKED*/, T: () => ["proto.Class", Class] }
        ]);
    }
    create(value) {
        const message = { id: 0, effectId: 0, name: "", isSpellId: false, type: 0, enchantType: 0, stats: [], quality: 0, phase: 0, requiredProfession: 0, classAllowlist: [] };
        Object.defineProperty(message, MESSAGE_TYPE, { enumerable: false, value: this });
        if (value !== undefined)
            reflectionMergePartial(this, message, value);
        return message;
    }
    internalBinaryRead(reader, length, options, target) {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* int32 id */ 1:
                    message.id = reader.int32();
                    break;
                case /* int32 effect_id */ 2:
                    message.effectId = reader.int32();
                    break;
                case /* string name */ 3:
                    message.name = reader.string();
                    break;
                case /* bool is_spell_id */ 10:
                    message.isSpellId = reader.bool();
                    break;
                case /* proto.ItemType type */ 4:
                    message.type = reader.int32();
                    break;
                case /* proto.EnchantType enchant_type */ 9:
                    message.enchantType = reader.int32();
                    break;
                case /* repeated double stats */ 7:
                    if (wireType === WireType.LengthDelimited)
                        for (let e = reader.int32() + reader.pos; reader.pos < e;)
                            message.stats.push(reader.double());
                    else
                        message.stats.push(reader.double());
                    break;
                case /* proto.ItemQuality quality */ 8:
                    message.quality = reader.int32();
                    break;
                case /* int32 phase */ 11:
                    message.phase = reader.int32();
                    break;
                case /* proto.Profession required_profession */ 13:
                    message.requiredProfession = reader.int32();
                    break;
                case /* repeated proto.Class class_allowlist */ 12:
                    if (wireType === WireType.LengthDelimited)
                        for (let e = reader.int32() + reader.pos; reader.pos < e;)
                            message.classAllowlist.push(reader.int32());
                    else
                        message.classAllowlist.push(reader.int32());
                    break;
                default:
                    let u = options.readUnknownField;
                    if (u === "throw")
                        throw new globalThis.Error(`Unknown field ${fieldNo} (wire type ${wireType}) for ${this.typeName}`);
                    let d = reader.skip(wireType);
                    if (u !== false)
                        (u === true ? UnknownFieldHandler.onRead : u)(this.typeName, message, fieldNo, wireType, d);
            }
        }
        return message;
    }
    internalBinaryWrite(message, writer, options) {
        /* int32 id = 1; */
        if (message.id !== 0)
            writer.tag(1, WireType.Varint).int32(message.id);
        /* int32 effect_id = 2; */
        if (message.effectId !== 0)
            writer.tag(2, WireType.Varint).int32(message.effectId);
        /* string name = 3; */
        if (message.name !== "")
            writer.tag(3, WireType.LengthDelimited).string(message.name);
        /* bool is_spell_id = 10; */
        if (message.isSpellId !== false)
            writer.tag(10, WireType.Varint).bool(message.isSpellId);
        /* proto.ItemType type = 4; */
        if (message.type !== 0)
            writer.tag(4, WireType.Varint).int32(message.type);
        /* proto.EnchantType enchant_type = 9; */
        if (message.enchantType !== 0)
            writer.tag(9, WireType.Varint).int32(message.enchantType);
        /* repeated double stats = 7; */
        if (message.stats.length) {
            writer.tag(7, WireType.LengthDelimited).fork();
            for (let i = 0; i < message.stats.length; i++)
                writer.double(message.stats[i]);
            writer.join();
        }
        /* proto.ItemQuality quality = 8; */
        if (message.quality !== 0)
            writer.tag(8, WireType.Varint).int32(message.quality);
        /* int32 phase = 11; */
        if (message.phase !== 0)
            writer.tag(11, WireType.Varint).int32(message.phase);
        /* proto.Profession required_profession = 13; */
        if (message.requiredProfession !== 0)
            writer.tag(13, WireType.Varint).int32(message.requiredProfession);
        /* repeated proto.Class class_allowlist = 12; */
        if (message.classAllowlist.length) {
            writer.tag(12, WireType.LengthDelimited).fork();
            for (let i = 0; i < message.classAllowlist.length; i++)
                writer.int32(message.classAllowlist[i]);
            writer.join();
        }
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message proto.Enchant
 */
export const Enchant = new Enchant$Type();
// @generated message type with reflection information, may provide speed optimized methods
class Gem$Type extends MessageType {
    constructor() {
        super("proto.Gem", [
            { no: 1, name: "id", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 2, name: "name", kind: "scalar", T: 9 /*ScalarType.STRING*/ },
            { no: 3, name: "stats", kind: "scalar", repeat: 1 /*RepeatType.PACKED*/, T: 1 /*ScalarType.DOUBLE*/ },
            { no: 4, name: "color", kind: "enum", T: () => ["proto.GemColor", GemColor] },
            { no: 5, name: "phase", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 6, name: "quality", kind: "enum", T: () => ["proto.ItemQuality", ItemQuality] },
            { no: 7, name: "unique", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 8, name: "required_profession", kind: "enum", T: () => ["proto.Profession", Profession] }
        ]);
    }
    create(value) {
        const message = { id: 0, name: "", stats: [], color: 0, phase: 0, quality: 0, unique: false, requiredProfession: 0 };
        Object.defineProperty(message, MESSAGE_TYPE, { enumerable: false, value: this });
        if (value !== undefined)
            reflectionMergePartial(this, message, value);
        return message;
    }
    internalBinaryRead(reader, length, options, target) {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* int32 id */ 1:
                    message.id = reader.int32();
                    break;
                case /* string name */ 2:
                    message.name = reader.string();
                    break;
                case /* repeated double stats */ 3:
                    if (wireType === WireType.LengthDelimited)
                        for (let e = reader.int32() + reader.pos; reader.pos < e;)
                            message.stats.push(reader.double());
                    else
                        message.stats.push(reader.double());
                    break;
                case /* proto.GemColor color */ 4:
                    message.color = reader.int32();
                    break;
                case /* int32 phase */ 5:
                    message.phase = reader.int32();
                    break;
                case /* proto.ItemQuality quality */ 6:
                    message.quality = reader.int32();
                    break;
                case /* bool unique */ 7:
                    message.unique = reader.bool();
                    break;
                case /* proto.Profession required_profession */ 8:
                    message.requiredProfession = reader.int32();
                    break;
                default:
                    let u = options.readUnknownField;
                    if (u === "throw")
                        throw new globalThis.Error(`Unknown field ${fieldNo} (wire type ${wireType}) for ${this.typeName}`);
                    let d = reader.skip(wireType);
                    if (u !== false)
                        (u === true ? UnknownFieldHandler.onRead : u)(this.typeName, message, fieldNo, wireType, d);
            }
        }
        return message;
    }
    internalBinaryWrite(message, writer, options) {
        /* int32 id = 1; */
        if (message.id !== 0)
            writer.tag(1, WireType.Varint).int32(message.id);
        /* string name = 2; */
        if (message.name !== "")
            writer.tag(2, WireType.LengthDelimited).string(message.name);
        /* repeated double stats = 3; */
        if (message.stats.length) {
            writer.tag(3, WireType.LengthDelimited).fork();
            for (let i = 0; i < message.stats.length; i++)
                writer.double(message.stats[i]);
            writer.join();
        }
        /* proto.GemColor color = 4; */
        if (message.color !== 0)
            writer.tag(4, WireType.Varint).int32(message.color);
        /* int32 phase = 5; */
        if (message.phase !== 0)
            writer.tag(5, WireType.Varint).int32(message.phase);
        /* proto.ItemQuality quality = 6; */
        if (message.quality !== 0)
            writer.tag(6, WireType.Varint).int32(message.quality);
        /* bool unique = 7; */
        if (message.unique !== false)
            writer.tag(7, WireType.Varint).bool(message.unique);
        /* proto.Profession required_profession = 8; */
        if (message.requiredProfession !== 0)
            writer.tag(8, WireType.Varint).int32(message.requiredProfession);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message proto.Gem
 */
export const Gem = new Gem$Type();
// @generated message type with reflection information, may provide speed optimized methods
class RaidTarget$Type extends MessageType {
    constructor() {
        super("proto.RaidTarget", [
            { no: 1, name: "target_index", kind: "scalar", T: 5 /*ScalarType.INT32*/ }
        ]);
    }
    create(value) {
        const message = { targetIndex: 0 };
        Object.defineProperty(message, MESSAGE_TYPE, { enumerable: false, value: this });
        if (value !== undefined)
            reflectionMergePartial(this, message, value);
        return message;
    }
    internalBinaryRead(reader, length, options, target) {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* int32 target_index */ 1:
                    message.targetIndex = reader.int32();
                    break;
                default:
                    let u = options.readUnknownField;
                    if (u === "throw")
                        throw new globalThis.Error(`Unknown field ${fieldNo} (wire type ${wireType}) for ${this.typeName}`);
                    let d = reader.skip(wireType);
                    if (u !== false)
                        (u === true ? UnknownFieldHandler.onRead : u)(this.typeName, message, fieldNo, wireType, d);
            }
        }
        return message;
    }
    internalBinaryWrite(message, writer, options) {
        /* int32 target_index = 1; */
        if (message.targetIndex !== 0)
            writer.tag(1, WireType.Varint).int32(message.targetIndex);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message proto.RaidTarget
 */
export const RaidTarget = new RaidTarget$Type();
// @generated message type with reflection information, may provide speed optimized methods
class ActionID$Type extends MessageType {
    constructor() {
        super("proto.ActionID", [
            { no: 1, name: "spell_id", kind: "scalar", oneof: "rawId", T: 5 /*ScalarType.INT32*/ },
            { no: 2, name: "item_id", kind: "scalar", oneof: "rawId", T: 5 /*ScalarType.INT32*/ },
            { no: 3, name: "other_id", kind: "enum", oneof: "rawId", T: () => ["proto.OtherAction", OtherAction] },
            { no: 4, name: "tag", kind: "scalar", T: 5 /*ScalarType.INT32*/ }
        ]);
    }
    create(value) {
        const message = { rawId: { oneofKind: undefined }, tag: 0 };
        Object.defineProperty(message, MESSAGE_TYPE, { enumerable: false, value: this });
        if (value !== undefined)
            reflectionMergePartial(this, message, value);
        return message;
    }
    internalBinaryRead(reader, length, options, target) {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* int32 spell_id */ 1:
                    message.rawId = {
                        oneofKind: "spellId",
                        spellId: reader.int32()
                    };
                    break;
                case /* int32 item_id */ 2:
                    message.rawId = {
                        oneofKind: "itemId",
                        itemId: reader.int32()
                    };
                    break;
                case /* proto.OtherAction other_id */ 3:
                    message.rawId = {
                        oneofKind: "otherId",
                        otherId: reader.int32()
                    };
                    break;
                case /* int32 tag */ 4:
                    message.tag = reader.int32();
                    break;
                default:
                    let u = options.readUnknownField;
                    if (u === "throw")
                        throw new globalThis.Error(`Unknown field ${fieldNo} (wire type ${wireType}) for ${this.typeName}`);
                    let d = reader.skip(wireType);
                    if (u !== false)
                        (u === true ? UnknownFieldHandler.onRead : u)(this.typeName, message, fieldNo, wireType, d);
            }
        }
        return message;
    }
    internalBinaryWrite(message, writer, options) {
        /* int32 spell_id = 1; */
        if (message.rawId.oneofKind === "spellId")
            writer.tag(1, WireType.Varint).int32(message.rawId.spellId);
        /* int32 item_id = 2; */
        if (message.rawId.oneofKind === "itemId")
            writer.tag(2, WireType.Varint).int32(message.rawId.itemId);
        /* proto.OtherAction other_id = 3; */
        if (message.rawId.oneofKind === "otherId")
            writer.tag(3, WireType.Varint).int32(message.rawId.otherId);
        /* int32 tag = 4; */
        if (message.tag !== 0)
            writer.tag(4, WireType.Varint).int32(message.tag);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message proto.ActionID
 */
export const ActionID = new ActionID$Type();
// @generated message type with reflection information, may provide speed optimized methods
class Glyphs$Type extends MessageType {
    constructor() {
        super("proto.Glyphs", [
            { no: 1, name: "major1", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 2, name: "major2", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 3, name: "major3", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 4, name: "minor1", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 5, name: "minor2", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 6, name: "minor3", kind: "scalar", T: 5 /*ScalarType.INT32*/ }
        ]);
    }
    create(value) {
        const message = { major1: 0, major2: 0, major3: 0, minor1: 0, minor2: 0, minor3: 0 };
        Object.defineProperty(message, MESSAGE_TYPE, { enumerable: false, value: this });
        if (value !== undefined)
            reflectionMergePartial(this, message, value);
        return message;
    }
    internalBinaryRead(reader, length, options, target) {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* int32 major1 */ 1:
                    message.major1 = reader.int32();
                    break;
                case /* int32 major2 */ 2:
                    message.major2 = reader.int32();
                    break;
                case /* int32 major3 */ 3:
                    message.major3 = reader.int32();
                    break;
                case /* int32 minor1 */ 4:
                    message.minor1 = reader.int32();
                    break;
                case /* int32 minor2 */ 5:
                    message.minor2 = reader.int32();
                    break;
                case /* int32 minor3 */ 6:
                    message.minor3 = reader.int32();
                    break;
                default:
                    let u = options.readUnknownField;
                    if (u === "throw")
                        throw new globalThis.Error(`Unknown field ${fieldNo} (wire type ${wireType}) for ${this.typeName}`);
                    let d = reader.skip(wireType);
                    if (u !== false)
                        (u === true ? UnknownFieldHandler.onRead : u)(this.typeName, message, fieldNo, wireType, d);
            }
        }
        return message;
    }
    internalBinaryWrite(message, writer, options) {
        /* int32 major1 = 1; */
        if (message.major1 !== 0)
            writer.tag(1, WireType.Varint).int32(message.major1);
        /* int32 major2 = 2; */
        if (message.major2 !== 0)
            writer.tag(2, WireType.Varint).int32(message.major2);
        /* int32 major3 = 3; */
        if (message.major3 !== 0)
            writer.tag(3, WireType.Varint).int32(message.major3);
        /* int32 minor1 = 4; */
        if (message.minor1 !== 0)
            writer.tag(4, WireType.Varint).int32(message.minor1);
        /* int32 minor2 = 5; */
        if (message.minor2 !== 0)
            writer.tag(5, WireType.Varint).int32(message.minor2);
        /* int32 minor3 = 6; */
        if (message.minor3 !== 0)
            writer.tag(6, WireType.Varint).int32(message.minor3);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message proto.Glyphs
 */
export const Glyphs = new Glyphs$Type();
// @generated message type with reflection information, may provide speed optimized methods
class Cooldown$Type extends MessageType {
    constructor() {
        super("proto.Cooldown", [
            { no: 1, name: "id", kind: "message", T: () => ActionID },
            { no: 2, name: "timings", kind: "scalar", repeat: 1 /*RepeatType.PACKED*/, T: 1 /*ScalarType.DOUBLE*/ }
        ]);
    }
    create(value) {
        const message = { timings: [] };
        Object.defineProperty(message, MESSAGE_TYPE, { enumerable: false, value: this });
        if (value !== undefined)
            reflectionMergePartial(this, message, value);
        return message;
    }
    internalBinaryRead(reader, length, options, target) {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* proto.ActionID id */ 1:
                    message.id = ActionID.internalBinaryRead(reader, reader.uint32(), options, message.id);
                    break;
                case /* repeated double timings */ 2:
                    if (wireType === WireType.LengthDelimited)
                        for (let e = reader.int32() + reader.pos; reader.pos < e;)
                            message.timings.push(reader.double());
                    else
                        message.timings.push(reader.double());
                    break;
                default:
                    let u = options.readUnknownField;
                    if (u === "throw")
                        throw new globalThis.Error(`Unknown field ${fieldNo} (wire type ${wireType}) for ${this.typeName}`);
                    let d = reader.skip(wireType);
                    if (u !== false)
                        (u === true ? UnknownFieldHandler.onRead : u)(this.typeName, message, fieldNo, wireType, d);
            }
        }
        return message;
    }
    internalBinaryWrite(message, writer, options) {
        /* proto.ActionID id = 1; */
        if (message.id)
            ActionID.internalBinaryWrite(message.id, writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        /* repeated double timings = 2; */
        if (message.timings.length) {
            writer.tag(2, WireType.LengthDelimited).fork();
            for (let i = 0; i < message.timings.length; i++)
                writer.double(message.timings[i]);
            writer.join();
        }
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message proto.Cooldown
 */
export const Cooldown = new Cooldown$Type();
// @generated message type with reflection information, may provide speed optimized methods
class Cooldowns$Type extends MessageType {
    constructor() {
        super("proto.Cooldowns", [
            { no: 1, name: "cooldowns", kind: "message", repeat: 1 /*RepeatType.PACKED*/, T: () => Cooldown },
            { no: 2, name: "hp_percent_for_defensives", kind: "scalar", T: 1 /*ScalarType.DOUBLE*/ }
        ]);
    }
    create(value) {
        const message = { cooldowns: [], hpPercentForDefensives: 0 };
        Object.defineProperty(message, MESSAGE_TYPE, { enumerable: false, value: this });
        if (value !== undefined)
            reflectionMergePartial(this, message, value);
        return message;
    }
    internalBinaryRead(reader, length, options, target) {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* repeated proto.Cooldown cooldowns */ 1:
                    message.cooldowns.push(Cooldown.internalBinaryRead(reader, reader.uint32(), options));
                    break;
                case /* double hp_percent_for_defensives */ 2:
                    message.hpPercentForDefensives = reader.double();
                    break;
                default:
                    let u = options.readUnknownField;
                    if (u === "throw")
                        throw new globalThis.Error(`Unknown field ${fieldNo} (wire type ${wireType}) for ${this.typeName}`);
                    let d = reader.skip(wireType);
                    if (u !== false)
                        (u === true ? UnknownFieldHandler.onRead : u)(this.typeName, message, fieldNo, wireType, d);
            }
        }
        return message;
    }
    internalBinaryWrite(message, writer, options) {
        /* repeated proto.Cooldown cooldowns = 1; */
        for (let i = 0; i < message.cooldowns.length; i++)
            Cooldown.internalBinaryWrite(message.cooldowns[i], writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        /* double hp_percent_for_defensives = 2; */
        if (message.hpPercentForDefensives !== 0)
            writer.tag(2, WireType.Bit64).double(message.hpPercentForDefensives);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message proto.Cooldowns
 */
export const Cooldowns = new Cooldowns$Type();
// @generated message type with reflection information, may provide speed optimized methods
class HealingModel$Type extends MessageType {
    constructor() {
        super("proto.HealingModel", [
            { no: 1, name: "hps", kind: "scalar", T: 1 /*ScalarType.DOUBLE*/ },
            { no: 2, name: "cadence_seconds", kind: "scalar", T: 1 /*ScalarType.DOUBLE*/ }
        ]);
    }
    create(value) {
        const message = { hps: 0, cadenceSeconds: 0 };
        Object.defineProperty(message, MESSAGE_TYPE, { enumerable: false, value: this });
        if (value !== undefined)
            reflectionMergePartial(this, message, value);
        return message;
    }
    internalBinaryRead(reader, length, options, target) {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* double hps */ 1:
                    message.hps = reader.double();
                    break;
                case /* double cadence_seconds */ 2:
                    message.cadenceSeconds = reader.double();
                    break;
                default:
                    let u = options.readUnknownField;
                    if (u === "throw")
                        throw new globalThis.Error(`Unknown field ${fieldNo} (wire type ${wireType}) for ${this.typeName}`);
                    let d = reader.skip(wireType);
                    if (u !== false)
                        (u === true ? UnknownFieldHandler.onRead : u)(this.typeName, message, fieldNo, wireType, d);
            }
        }
        return message;
    }
    internalBinaryWrite(message, writer, options) {
        /* double hps = 1; */
        if (message.hps !== 0)
            writer.tag(1, WireType.Bit64).double(message.hps);
        /* double cadence_seconds = 2; */
        if (message.cadenceSeconds !== 0)
            writer.tag(2, WireType.Bit64).double(message.cadenceSeconds);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message proto.HealingModel
 */
export const HealingModel = new HealingModel$Type();
