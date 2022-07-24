import type { BinaryWriteOptions } from "@protobuf-ts/runtime";
import type { IBinaryWriter } from "@protobuf-ts/runtime";
import type { BinaryReadOptions } from "@protobuf-ts/runtime";
import type { IBinaryReader } from "@protobuf-ts/runtime";
import type { PartialMessage } from "@protobuf-ts/runtime";
import { MessageType } from "@protobuf-ts/runtime";
/**
 * Buffs that affect the entire raid.
 *
 * @generated from protobuf message proto.RaidBuffs
 */
export interface RaidBuffs {
    /**
     * +Stats
     *
     * @generated from protobuf field: proto.TristateEffect gift_of_the_wild = 1;
     */
    giftOfTheWild: TristateEffect;
    /**
     * +Stam
     *
     * @generated from protobuf field: proto.TristateEffect power_word_fortitude = 2;
     */
    powerWordFortitude: TristateEffect;
    /**
     * +Health
     *
     * @generated from protobuf field: proto.TristateEffect commanding_shout = 3;
     */
    commandingShout: TristateEffect;
    /**
     * @generated from protobuf field: proto.TristateEffect blood_pact = 4;
     */
    bloodPact: TristateEffect;
    /**
     * + Agi and Str
     *
     * @generated from protobuf field: bool horn_of_winter = 5;
     */
    hornOfWinter: boolean;
    /**
     * @generated from protobuf field: proto.TristateEffect strength_of_earth_totem = 6;
     */
    strengthOfEarthTotem: TristateEffect;
    /**
     * +Intell and/or Spi
     *
     * @generated from protobuf field: bool arcane_brilliance = 7;
     */
    arcaneBrilliance: boolean;
    /**
     * @generated from protobuf field: proto.TristateEffect fel_intelligence = 8;
     */
    felIntelligence: TristateEffect;
    /**
     * @generated from protobuf field: bool divine_spirit = 9;
     */
    divineSpirit: boolean;
    /**
     * +AP
     *
     * @generated from protobuf field: proto.TristateEffect battle_shout = 10;
     */
    battleShout: TristateEffect;
    /**
     * 10% AP
     *
     * @generated from protobuf field: bool trueshot_aura = 11;
     */
    trueshotAura: boolean;
    /**
     * @generated from protobuf field: bool unleashed_rage = 12;
     */
    unleashedRage: boolean;
    /**
     * @generated from protobuf field: bool abominations_might = 13;
     */
    abominationsMight: boolean;
    /**
     * 5% phy crit
     *
     * @generated from protobuf field: proto.TristateEffect leader_of_the_pack = 14;
     */
    leaderOfThePack: TristateEffect;
    /**
     * @generated from protobuf field: bool rampage = 15;
     */
    rampage: boolean;
    /**
     * 20% Melee Haste
     *
     * @generated from protobuf field: bool icy_talons = 16;
     */
    icyTalons: boolean;
    /**
     * @generated from protobuf field: proto.TristateEffect windfury_totem = 17;
     */
    windfuryTotem: TristateEffect;
    /**
     * +Spell Power
     *
     * @generated from protobuf field: bool totem_of_wrath = 18;
     */
    totemOfWrath: boolean;
    /**
     * @generated from protobuf field: bool flametongue_totem = 19;
     */
    flametongueTotem: boolean;
    /**
     * @generated from protobuf field: int32 demonic_pact = 20;
     */
    demonicPact: number;
    /**
     * +5% Spell Crit and/or +3% Haste
     *
     * @generated from protobuf field: bool swift_retribution = 21;
     */
    swiftRetribution: boolean;
    /**
     * @generated from protobuf field: proto.TristateEffect moonkin_aura = 22;
     */
    moonkinAura: TristateEffect;
    /**
     * @generated from protobuf field: bool elemental_oath = 23;
     */
    elementalOath: boolean;
    /**
     * 5% spell haste
     *
     * @generated from protobuf field: bool wrath_of_air_totem = 24;
     */
    wrathOfAirTotem: boolean;
    /**
     * 3% dmg
     *
     * @generated from protobuf field: bool ferocious_inspiration = 25;
     */
    ferociousInspiration: boolean;
    /**
     * @generated from protobuf field: bool sanctified_retribution = 26;
     */
    sanctifiedRetribution: boolean;
    /**
     * @generated from protobuf field: bool arcane_empowerment = 27;
     */
    arcaneEmpowerment: boolean;
    /**
     * mp5
     *
     * @generated from protobuf field: proto.TristateEffect mana_spring_totem = 28;
     */
    manaSpringTotem: TristateEffect;
    /**
     * Miscellaneous
     *
     * @generated from protobuf field: bool bloodlust = 29;
     */
    bloodlust: boolean;
    /**
     * @generated from protobuf field: proto.TristateEffect thorns = 30;
     */
    thorns: TristateEffect;
    /**
     * @generated from protobuf field: proto.TristateEffect devotion_aura = 31;
     */
    devotionAura: TristateEffect;
    /**
     * @generated from protobuf field: bool retribution_aura = 32;
     */
    retributionAura: boolean;
    /**
     * @generated from protobuf field: bool shadow_protection = 33;
     */
    shadowProtection: boolean;
    /**
     * Drums
     *
     * @generated from protobuf field: bool drums_of_forgotten_kings = 34;
     */
    drumsOfForgottenKings: boolean;
    /**
     * @generated from protobuf field: bool drums_of_the_wild = 35;
     */
    drumsOfTheWild: boolean;
    /**
     * Scroll
     *
     * @generated from protobuf field: bool scroll_of_protection = 36;
     */
    scrollOfProtection: boolean;
    /**
     * @generated from protobuf field: bool scroll_of_stamina = 37;
     */
    scrollOfStamina: boolean;
    /**
     * @generated from protobuf field: bool scroll_of_strength = 38;
     */
    scrollOfStrength: boolean;
    /**
     * @generated from protobuf field: bool scroll_of_agility = 39;
     */
    scrollOfAgility: boolean;
    /**
     * @generated from protobuf field: bool scroll_of_intellect = 40;
     */
    scrollOfIntellect: boolean;
    /**
     * @generated from protobuf field: bool scroll_of_spirit = 41;
     */
    scrollOfSpirit: boolean;
}
/**
 * Buffs that affect a single party.
 *
 * @generated from protobuf message proto.PartyBuffs
 */
export interface PartyBuffs {
    /**
     * Item Buffs
     *
     * @generated from protobuf field: int32 atiesh_mage = 1;
     */
    atieshMage: number;
    /**
     * @generated from protobuf field: int32 atiesh_warlock = 2;
     */
    atieshWarlock: number;
    /**
     * @generated from protobuf field: bool braided_eternium_chain = 3;
     */
    braidedEterniumChain: boolean;
    /**
     * @generated from protobuf field: bool eye_of_the_night = 4;
     */
    eyeOfTheNight: boolean;
    /**
     * @generated from protobuf field: bool chain_of_the_twilight_owl = 5;
     */
    chainOfTheTwilightOwl: boolean;
    /**
     * Group buffs
     *
     * @generated from protobuf field: int32 mana_tide_totems = 6;
     */
    manaTideTotems: number;
    /**
     * @generated from protobuf field: bool heroic_presence = 7;
     */
    heroicPresence: boolean;
}
/**
 * These are usually individual actions taken by other Characters.
 *
 * @generated from protobuf message proto.IndividualBuffs
 */
export interface IndividualBuffs {
    /**
     * 10% Stats
     *
     * @generated from protobuf field: bool blessing_of_kings = 1;
     */
    blessingOfKings: boolean;
    /**
     * @generated from protobuf field: proto.TristateEffect blessing_of_wisdom = 2;
     */
    blessingOfWisdom: TristateEffect;
    /**
     * @generated from protobuf field: proto.TristateEffect blessing_of_might = 3;
     */
    blessingOfMight: TristateEffect;
    /**
     * @generated from protobuf field: bool blessing_of_sanctuary = 4;
     */
    blessingOfSanctuary: boolean;
    /**
     * @generated from protobuf field: bool vigilance = 5;
     */
    vigilance: boolean;
    /**
     * @generated from protobuf field: bool renewed_hope = 6;
     */
    renewedHope: boolean;
    /**
     * How many of each of these buffs the player will be receiving.
     *
     * @generated from protobuf field: int32 hymn_of_hope = 7;
     */
    hymnOfHope: number;
    /**
     * @generated from protobuf field: int32 hand_of_salvation = 8;
     */
    handOfSalvation: number;
    /**
     * @generated from protobuf field: int32 rapture = 9;
     */
    rapture: number;
    /**
     * @generated from protobuf field: int32 innervates = 10;
     */
    innervates: number;
    /**
     * @generated from protobuf field: int32 power_infusions = 11;
     */
    powerInfusions: number;
    /**
     * @generated from protobuf field: int32 unholy_frenzy = 12;
     */
    unholyFrenzy: number;
    /**
     * @generated from protobuf field: int32 revitalize = 13;
     */
    revitalize: number;
    /**
     * @generated from protobuf field: bool vampiric_touch = 14;
     */
    vampiricTouch: boolean;
    /**
     * @generated from protobuf field: bool hunting_party = 15;
     */
    huntingParty: boolean;
    /**
     * @generated from protobuf field: bool judgements_of_the_wise = 16;
     */
    judgementsOfTheWise: boolean;
    /**
     * @generated from protobuf field: bool improved_soul_leech = 17;
     */
    improvedSoulLeech: boolean;
    /**
     * @generated from protobuf field: bool enduring_winter = 18;
     */
    enduringWinter: boolean;
}
/**
 * @generated from protobuf message proto.Consumes
 */
export interface Consumes {
    /**
     * @generated from protobuf field: proto.Flask flask = 1;
     */
    flask: Flask;
    /**
     * @generated from protobuf field: proto.BattleElixir battle_elixir = 2;
     */
    battleElixir: BattleElixir;
    /**
     * @generated from protobuf field: proto.GuardianElixir guardian_elixir = 3;
     */
    guardianElixir: GuardianElixir;
    /**
     * @generated from protobuf field: proto.WeaponImbue main_hand_imbue = 4;
     */
    mainHandImbue: WeaponImbue;
    /**
     * @generated from protobuf field: proto.WeaponImbue off_hand_imbue = 5;
     */
    offHandImbue: WeaponImbue;
    /**
     * @generated from protobuf field: proto.Food food = 6;
     */
    food: Food;
    /**
     * @generated from protobuf field: proto.PetFood pet_food = 7;
     */
    petFood: PetFood;
    /**
     * @generated from protobuf field: int32 pet_scroll_of_agility = 8;
     */
    petScrollOfAgility: number;
    /**
     * @generated from protobuf field: int32 pet_scroll_of_strength = 9;
     */
    petScrollOfStrength: number;
    /**
     * @generated from protobuf field: proto.Potions default_potion = 10;
     */
    defaultPotion: Potions;
    /**
     * @generated from protobuf field: proto.Potions prepop_potion = 11;
     */
    prepopPotion: Potions;
    /**
     * @generated from protobuf field: proto.Conjured default_conjured = 12;
     */
    defaultConjured: Conjured;
    /**
     * @generated from protobuf field: proto.Conjured starting_conjured = 13;
     */
    startingConjured: Conjured;
    /**
     * @generated from protobuf field: int32 num_starting_conjured = 14;
     */
    numStartingConjured: number;
    /**
     * @generated from protobuf field: bool super_sapper = 15;
     */
    superSapper: boolean;
    /**
     * @generated from protobuf field: bool goblin_sapper = 16;
     */
    goblinSapper: boolean;
    /**
     * @generated from protobuf field: proto.Explosive filler_explosive = 17;
     */
    fillerExplosive: Explosive;
}
/**
 * @generated from protobuf message proto.Debuffs
 */
export interface Debuffs {
    /**
     * @generated from protobuf field: bool judgement_of_wisdom = 1;
     */
    judgementOfWisdom: boolean;
    /**
     * @generated from protobuf field: bool judgement_of_light = 2;
     */
    judgementOfLight: boolean;
    /**
     * @generated from protobuf field: bool misery = 3;
     */
    misery: boolean;
    /**
     * @generated from protobuf field: proto.TristateEffect faerie_fire = 4;
     */
    faerieFire: TristateEffect;
    /**
     * 13% bonus spell damage
     *
     * @generated from protobuf field: bool curse_of_elements = 5;
     */
    curseOfElements: boolean;
    /**
     * @generated from protobuf field: bool ebon_plaguebringer = 6;
     */
    ebonPlaguebringer: boolean;
    /**
     * @generated from protobuf field: bool earth_and_moon = 7;
     */
    earthAndMoon: boolean;
    /**
     * +3% to crit against target
     *
     * @generated from protobuf field: bool heart_of_the_crusader = 8;
     */
    heartOfTheCrusader: boolean;
    /**
     * @generated from protobuf field: bool master_poisoner = 9;
     */
    masterPoisoner: boolean;
    /**
     * @generated from protobuf field: bool totem_of_wrath = 10;
     */
    totemOfWrath: boolean;
    /**
     * 5% spell crit
     *
     * @generated from protobuf field: bool shadow_mastery = 11;
     */
    shadowMastery: boolean;
    /**
     * @generated from protobuf field: bool improved_scorch = 12;
     */
    improvedScorch: boolean;
    /**
     * @generated from protobuf field: bool winters_chill = 13;
     */
    wintersChill: boolean;
    /**
     * @generated from protobuf field: bool blood_frenzy = 14;
     */
    bloodFrenzy: boolean;
    /**
     * @generated from protobuf field: bool savage_combat = 15;
     */
    savageCombat: boolean;
    /**
     * TODO: validate these
     *
     * @generated from protobuf field: bool gift_of_arthas = 16;
     */
    giftOfArthas: boolean;
    /**
     * Bleed %
     *
     * @generated from protobuf field: bool mangle = 17;
     */
    mangle: boolean;
    /**
     * @generated from protobuf field: bool trauma = 18;
     */
    trauma: boolean;
    /**
     * @generated from protobuf field: bool stampede = 19;
     */
    stampede: boolean;
    /**
     * Major armor
     *
     * @generated from protobuf field: bool expose_armor = 20;
     */
    exposeArmor: boolean;
    /**
     * @generated from protobuf field: bool sunder_armor = 21;
     */
    sunderArmor: boolean;
    /**
     * @generated from protobuf field: bool acid_spit = 22;
     */
    acidSpit: boolean;
    /**
     * @generated from protobuf field: proto.TristateEffect curse_of_weakness = 23;
     */
    curseOfWeakness: TristateEffect;
    /**
     * @generated from protobuf field: bool sting = 24;
     */
    sting: boolean;
    /**
     * @generated from protobuf field: proto.TristateEffect demoralizing_roar = 25;
     */
    demoralizingRoar: TristateEffect;
    /**
     * @generated from protobuf field: proto.TristateEffect demoralizing_shout = 26;
     */
    demoralizingShout: TristateEffect;
    /**
     * @generated from protobuf field: proto.TristateEffect thunder_clap = 27;
     */
    thunderClap: TristateEffect;
    /**
     * @generated from protobuf field: proto.TristateEffect frost_fever = 28;
     */
    frostFever: TristateEffect;
    /**
     * @generated from protobuf field: bool infected_wounds = 29;
     */
    infectedWounds: boolean;
    /**
     * @generated from protobuf field: bool judgements_of_the_just = 30;
     */
    judgementsOfTheJust: boolean;
    /**
     * @generated from protobuf field: bool insect_swarm = 31;
     */
    insectSwarm: boolean;
    /**
     * @generated from protobuf field: bool scorpid_sting = 32;
     */
    scorpidSting: boolean;
    /**
     * @generated from protobuf field: bool shadow_embrace = 33;
     */
    shadowEmbrace: boolean;
    /**
     * @generated from protobuf field: bool screech = 34;
     */
    screech: boolean;
}
/**
 * @generated from protobuf message proto.Target
 */
export interface Target {
    /**
     * The in-game NPC ID.
     *
     * @generated from protobuf field: int32 id = 14;
     */
    id: number;
    /**
     * @generated from protobuf field: string name = 15;
     */
    name: string;
    /**
     * @generated from protobuf field: int32 level = 4;
     */
    level: number;
    /**
     * @generated from protobuf field: proto.MobType mob_type = 3;
     */
    mobType: MobType;
    /**
     * @generated from protobuf field: repeated double stats = 5;
     */
    stats: number[];
    /**
     * Auto attack parameters.
     *
     * @generated from protobuf field: double min_base_damage = 7;
     */
    minBaseDamage: number;
    /**
     * @generated from protobuf field: double swing_speed = 8;
     */
    swingSpeed: number;
    /**
     * @generated from protobuf field: bool dual_wield = 9;
     */
    dualWield: boolean;
    /**
     * @generated from protobuf field: bool dual_wield_penalty = 10;
     */
    dualWieldPenalty: boolean;
    /**
     * @generated from protobuf field: bool parry_haste = 12;
     */
    parryHaste: boolean;
    /**
     * @generated from protobuf field: bool suppress_dodge = 16;
     */
    suppressDodge: boolean;
    /**
     * @generated from protobuf field: proto.SpellSchool spell_school = 13;
     */
    spellSchool: SpellSchool;
    /**
     * Index in Raid.tanks indicating the player tanking this mob.
     * -1 or invalid index indicates not being tanked.
     *
     * @generated from protobuf field: int32 tank_index = 6;
     */
    tankIndex: number;
}
/**
 * @generated from protobuf message proto.Encounter
 */
export interface Encounter {
    /**
     * @generated from protobuf field: double duration = 1;
     */
    duration: number;
    /**
     * Variation in the duration
     *
     * @generated from protobuf field: double duration_variation = 2;
     */
    durationVariation: number;
    /**
     * The ratio of the encounter duration, between 0 and 1, for which the targets
     * will be in execute range (<= 20%) for the purposes of Warrior Execute, Mage Molten
     * Fury, etc.
     *
     * @generated from protobuf field: double execute_proportion_20 = 3;
     */
    executeProportion20: number;
    /**
     * Same as execute_proportion but for 35%.
     *
     * @generated from protobuf field: double execute_proportion_35 = 4;
     */
    executeProportion35: number;
    /**
     * If set, will use the targets health value instead of a duration for fight length.
     *
     * @generated from protobuf field: bool use_health = 5;
     */
    useHealth: boolean;
    /**
     * If type != Simple or Custom, then this may be empty.
     *
     * @generated from protobuf field: repeated proto.Target targets = 6;
     */
    targets: Target[];
}
/**
 * @generated from protobuf message proto.ItemSpec
 */
export interface ItemSpec {
    /**
     * @generated from protobuf field: int32 id = 2;
     */
    id: number;
    /**
     * @generated from protobuf field: int32 enchant = 3;
     */
    enchant: number;
    /**
     * @generated from protobuf field: repeated int32 gems = 4;
     */
    gems: number[];
}
/**
 * @generated from protobuf message proto.EquipmentSpec
 */
export interface EquipmentSpec {
    /**
     * @generated from protobuf field: repeated proto.ItemSpec items = 1;
     */
    items: ItemSpec[];
}
/**
 * @generated from protobuf message proto.Item
 */
export interface Item {
    /**
     * @generated from protobuf field: int32 id = 1;
     */
    id: number;
    /**
     * This is unused by most items. For most items we set id to the
     * wowhead/in-game ID directly. For random enchant items though we need to
     * use unique hardcoded IDs so this field holds the wowhead ID instead.
     *
     * @generated from protobuf field: int32 wowhead_id = 16;
     */
    wowheadId: number;
    /**
     * @generated from protobuf field: string name = 2;
     */
    name: string;
    /**
     * Classes that are allowed to use the item. Empty indicates no special class restrictions.
     *
     * @generated from protobuf field: repeated proto.Class class_allowlist = 15;
     */
    classAllowlist: Class[];
    /**
     * @generated from protobuf field: proto.ItemType type = 3;
     */
    type: ItemType;
    /**
     * @generated from protobuf field: proto.ArmorType armor_type = 4;
     */
    armorType: ArmorType;
    /**
     * @generated from protobuf field: proto.WeaponType weapon_type = 5;
     */
    weaponType: WeaponType;
    /**
     * @generated from protobuf field: proto.HandType hand_type = 6;
     */
    handType: HandType;
    /**
     * @generated from protobuf field: proto.RangedWeaponType ranged_weapon_type = 7;
     */
    rangedWeaponType: RangedWeaponType;
    /**
     * @generated from protobuf field: repeated double stats = 8;
     */
    stats: number[];
    /**
     * @generated from protobuf field: repeated proto.GemColor gem_sockets = 9;
     */
    gemSockets: GemColor[];
    /**
     * @generated from protobuf field: repeated double socketBonus = 10;
     */
    socketBonus: number[];
    /**
     * Weapon stats, needed for computing proper EP for melee weapons
     *
     * @generated from protobuf field: double weapon_damage_min = 17;
     */
    weaponDamageMin: number;
    /**
     * @generated from protobuf field: double weapon_damage_max = 18;
     */
    weaponDamageMax: number;
    /**
     * @generated from protobuf field: double weapon_speed = 19;
     */
    weaponSpeed: number;
    /**
     * @generated from protobuf field: int32 phase = 11;
     */
    phase: number;
    /**
     * @generated from protobuf field: proto.ItemQuality quality = 12;
     */
    quality: ItemQuality;
    /**
     * @generated from protobuf field: bool unique = 13;
     */
    unique: boolean;
    /**
     * @generated from protobuf field: int32 ilvl = 20;
     */
    ilvl: number;
    /**
     * @generated from protobuf field: proto.Profession required_profession = 21;
     */
    requiredProfession: Profession;
}
/**
 * @generated from protobuf message proto.Enchant
 */
export interface Enchant {
    /**
     * @generated from protobuf field: int32 id = 1;
     */
    id: number;
    /**
     * @generated from protobuf field: int32 effect_id = 2;
     */
    effectId: number;
    /**
     * @generated from protobuf field: string name = 3;
     */
    name: string;
    /**
     * If true, then id is the ID of the enchant spell instead of the formula item.
     * This is used by enchants for which a formula doesn't exist (its taught by a trainer).
     *
     * @generated from protobuf field: bool is_spell_id = 10;
     */
    isSpellId: boolean;
    /**
     * @generated from protobuf field: proto.ItemType type = 4;
     */
    type: ItemType;
    /**
     * @generated from protobuf field: proto.EnchantType enchant_type = 9;
     */
    enchantType: EnchantType;
    /**
     * @generated from protobuf field: repeated double stats = 7;
     */
    stats: number[];
    /**
     * @generated from protobuf field: proto.ItemQuality quality = 8;
     */
    quality: ItemQuality;
    /**
     * @generated from protobuf field: int32 phase = 11;
     */
    phase: number;
    /**
     * @generated from protobuf field: proto.Profession required_profession = 13;
     */
    requiredProfession: Profession;
    /**
     * Classes that are allowed to use the enchant. Empty indicates no special class restrictions.
     *
     * @generated from protobuf field: repeated proto.Class class_allowlist = 12;
     */
    classAllowlist: Class[];
}
/**
 * @generated from protobuf message proto.Gem
 */
export interface Gem {
    /**
     * @generated from protobuf field: int32 id = 1;
     */
    id: number;
    /**
     * @generated from protobuf field: string name = 2;
     */
    name: string;
    /**
     * @generated from protobuf field: repeated double stats = 3;
     */
    stats: number[];
    /**
     * @generated from protobuf field: proto.GemColor color = 4;
     */
    color: GemColor;
    /**
     * @generated from protobuf field: int32 phase = 5;
     */
    phase: number;
    /**
     * @generated from protobuf field: proto.ItemQuality quality = 6;
     */
    quality: ItemQuality;
    /**
     * @generated from protobuf field: bool unique = 7;
     */
    unique: boolean;
    /**
     * @generated from protobuf field: proto.Profession required_profession = 8;
     */
    requiredProfession: Profession;
}
/**
 * @generated from protobuf message proto.RaidTarget
 */
export interface RaidTarget {
    /**
     * Raid index of the player to target. A value of -1 indicates no target.
     *
     * @generated from protobuf field: int32 target_index = 1;
     */
    targetIndex: number;
}
/**
 * @generated from protobuf message proto.ActionID
 */
export interface ActionID {
    /**
     * @generated from protobuf oneof: raw_id
     */
    rawId: {
        oneofKind: "spellId";
        /**
         * @generated from protobuf field: int32 spell_id = 1;
         */
        spellId: number;
    } | {
        oneofKind: "itemId";
        /**
         * @generated from protobuf field: int32 item_id = 2;
         */
        itemId: number;
    } | {
        oneofKind: "otherId";
        /**
         * @generated from protobuf field: proto.OtherAction other_id = 3;
         */
        otherId: OtherAction;
    } | {
        oneofKind: undefined;
    };
    /**
     * Distinguishes between different versions of the same action.
     * Currently the only use for this is Shaman Lightning Overload.
     *
     * @generated from protobuf field: int32 tag = 4;
     */
    tag: number;
}
/**
 * @generated from protobuf message proto.Glyphs
 */
export interface Glyphs {
    /**
     * @generated from protobuf field: int32 major1 = 1;
     */
    major1: number;
    /**
     * @generated from protobuf field: int32 major2 = 2;
     */
    major2: number;
    /**
     * @generated from protobuf field: int32 major3 = 3;
     */
    major3: number;
    /**
     * @generated from protobuf field: int32 minor1 = 4;
     */
    minor1: number;
    /**
     * @generated from protobuf field: int32 minor2 = 5;
     */
    minor2: number;
    /**
     * @generated from protobuf field: int32 minor3 = 6;
     */
    minor3: number;
}
/**
 * Custom options for a particular cooldown.
 *
 * @generated from protobuf message proto.Cooldown
 */
export interface Cooldown {
    /**
     * Identifies the cooldown to which these settings will apply.
     *
     * @generated from protobuf field: proto.ActionID id = 1;
     */
    id?: ActionID;
    /**
     * Fixed times at which to use this cooldown. Each value corresponds to a usage,
     * e.g. first value is the first usage, second value is the second usage.
     * Any usages after the specified timings will occur as soon as possible, subject
     * to the ShouldActivate() condition.
     *
     * @generated from protobuf field: repeated double timings = 2;
     */
    timings: number[];
}
/**
 * @generated from protobuf message proto.Cooldowns
 */
export interface Cooldowns {
    /**
     * @generated from protobuf field: repeated proto.Cooldown cooldowns = 1;
     */
    cooldowns: Cooldown[];
    /**
     * % HP threshold, below which defensive cooldowns can be used.
     *
     * @generated from protobuf field: double hp_percent_for_defensives = 2;
     */
    hpPercentForDefensives: number;
}
/**
 * @generated from protobuf message proto.HealingModel
 */
export interface HealingModel {
    /**
     * Healing per second to apply.
     *
     * @generated from protobuf field: double hps = 1;
     */
    hps: number;
    /**
     * How often healing is applied.
     *
     * @generated from protobuf field: double cadence_seconds = 2;
     */
    cadenceSeconds: number;
}
/**
 * @generated from protobuf enum proto.Spec
 */
export declare enum Spec {
    /**
     * @generated from protobuf enum value: SpecBalanceDruid = 0;
     */
    SpecBalanceDruid = 0,
    /**
     * @generated from protobuf enum value: SpecElementalShaman = 1;
     */
    SpecElementalShaman = 1,
    /**
     * @generated from protobuf enum value: SpecEnhancementShaman = 9;
     */
    SpecEnhancementShaman = 9,
    /**
     * @generated from protobuf enum value: SpecFeralDruid = 12;
     */
    SpecFeralDruid = 12,
    /**
     * @generated from protobuf enum value: SpecFeralTankDruid = 14;
     */
    SpecFeralTankDruid = 14,
    /**
     * @generated from protobuf enum value: SpecHunter = 8;
     */
    SpecHunter = 8,
    /**
     * @generated from protobuf enum value: SpecMage = 2;
     */
    SpecMage = 2,
    /**
     * @generated from protobuf enum value: SpecProtectionPaladin = 13;
     */
    SpecProtectionPaladin = 13,
    /**
     * @generated from protobuf enum value: SpecRetributionPaladin = 3;
     */
    SpecRetributionPaladin = 3,
    /**
     * @generated from protobuf enum value: SpecRogue = 7;
     */
    SpecRogue = 7,
    /**
     * @generated from protobuf enum value: SpecShadowPriest = 4;
     */
    SpecShadowPriest = 4,
    /**
     * @generated from protobuf enum value: SpecSmitePriest = 10;
     */
    SpecSmitePriest = 10,
    /**
     * @generated from protobuf enum value: SpecWarlock = 5;
     */
    SpecWarlock = 5,
    /**
     * @generated from protobuf enum value: SpecWarrior = 6;
     */
    SpecWarrior = 6,
    /**
     * @generated from protobuf enum value: SpecProtectionWarrior = 11;
     */
    SpecProtectionWarrior = 11,
    /**
     * @generated from protobuf enum value: SpecDeathknight = 15;
     */
    SpecDeathknight = 15,
    /**
     * @generated from protobuf enum value: SpecTankDeathknight = 16;
     */
    SpecTankDeathknight = 16
}
/**
 * @generated from protobuf enum proto.Race
 */
export declare enum Race {
    /**
     * @generated from protobuf enum value: RaceUnknown = 0;
     */
    RaceUnknown = 0,
    /**
     * @generated from protobuf enum value: RaceBloodElf = 1;
     */
    RaceBloodElf = 1,
    /**
     * @generated from protobuf enum value: RaceDraenei = 2;
     */
    RaceDraenei = 2,
    /**
     * @generated from protobuf enum value: RaceDwarf = 3;
     */
    RaceDwarf = 3,
    /**
     * @generated from protobuf enum value: RaceGnome = 4;
     */
    RaceGnome = 4,
    /**
     * @generated from protobuf enum value: RaceHuman = 5;
     */
    RaceHuman = 5,
    /**
     * @generated from protobuf enum value: RaceNightElf = 6;
     */
    RaceNightElf = 6,
    /**
     * @generated from protobuf enum value: RaceOrc = 7;
     */
    RaceOrc = 7,
    /**
     * @generated from protobuf enum value: RaceTauren = 8;
     */
    RaceTauren = 8,
    /**
     * @generated from protobuf enum value: RaceTroll = 9;
     */
    RaceTroll = 9,
    /**
     * @generated from protobuf enum value: RaceUndead = 10;
     */
    RaceUndead = 10
}
/**
 * @generated from protobuf enum proto.Faction
 */
export declare enum Faction {
    /**
     * @generated from protobuf enum value: Unknown = 0;
     */
    Unknown = 0,
    /**
     * @generated from protobuf enum value: Alliance = 1;
     */
    Alliance = 1,
    /**
     * @generated from protobuf enum value: Horde = 2;
     */
    Horde = 2
}
/**
 * @generated from protobuf enum proto.ShattrathFaction
 */
export declare enum ShattrathFaction {
    /**
     * @generated from protobuf enum value: ShattrathFactionAldor = 0;
     */
    ShattrathFactionAldor = 0,
    /**
     * @generated from protobuf enum value: ShattrathFactionScryer = 1;
     */
    ShattrathFactionScryer = 1
}
/**
 * @generated from protobuf enum proto.Class
 */
export declare enum Class {
    /**
     * @generated from protobuf enum value: ClassUnknown = 0;
     */
    ClassUnknown = 0,
    /**
     * @generated from protobuf enum value: ClassDruid = 1;
     */
    ClassDruid = 1,
    /**
     * @generated from protobuf enum value: ClassHunter = 2;
     */
    ClassHunter = 2,
    /**
     * @generated from protobuf enum value: ClassMage = 3;
     */
    ClassMage = 3,
    /**
     * @generated from protobuf enum value: ClassPaladin = 4;
     */
    ClassPaladin = 4,
    /**
     * @generated from protobuf enum value: ClassPriest = 5;
     */
    ClassPriest = 5,
    /**
     * @generated from protobuf enum value: ClassRogue = 6;
     */
    ClassRogue = 6,
    /**
     * @generated from protobuf enum value: ClassShaman = 7;
     */
    ClassShaman = 7,
    /**
     * @generated from protobuf enum value: ClassWarlock = 8;
     */
    ClassWarlock = 8,
    /**
     * @generated from protobuf enum value: ClassWarrior = 9;
     */
    ClassWarrior = 9,
    /**
     * @generated from protobuf enum value: ClassDeathknight = 10;
     */
    ClassDeathknight = 10
}
/**
 * @generated from protobuf enum proto.Profession
 */
export declare enum Profession {
    /**
     * @generated from protobuf enum value: ProfessionUnknown = 0;
     */
    ProfessionUnknown = 0,
    /**
     * @generated from protobuf enum value: Alchemy = 1;
     */
    Alchemy = 1,
    /**
     * @generated from protobuf enum value: Blacksmithing = 2;
     */
    Blacksmithing = 2,
    /**
     * @generated from protobuf enum value: Enchanting = 3;
     */
    Enchanting = 3,
    /**
     * @generated from protobuf enum value: Engineering = 4;
     */
    Engineering = 4,
    /**
     * @generated from protobuf enum value: Herbalism = 5;
     */
    Herbalism = 5,
    /**
     * @generated from protobuf enum value: Inscription = 6;
     */
    Inscription = 6,
    /**
     * @generated from protobuf enum value: Jewelcrafting = 7;
     */
    Jewelcrafting = 7,
    /**
     * @generated from protobuf enum value: Leatherworking = 8;
     */
    Leatherworking = 8,
    /**
     * @generated from protobuf enum value: Mining = 9;
     */
    Mining = 9,
    /**
     * @generated from protobuf enum value: Skinning = 10;
     */
    Skinning = 10,
    /**
     * @generated from protobuf enum value: Tailoring = 11;
     */
    Tailoring = 11
}
/**
 * @generated from protobuf enum proto.Stat
 */
export declare enum Stat {
    /**
     * @generated from protobuf enum value: StatStrength = 0;
     */
    StatStrength = 0,
    /**
     * @generated from protobuf enum value: StatAgility = 1;
     */
    StatAgility = 1,
    /**
     * @generated from protobuf enum value: StatStamina = 2;
     */
    StatStamina = 2,
    /**
     * @generated from protobuf enum value: StatIntellect = 3;
     */
    StatIntellect = 3,
    /**
     * @generated from protobuf enum value: StatSpirit = 4;
     */
    StatSpirit = 4,
    /**
     * @generated from protobuf enum value: StatSpellPower = 5;
     */
    StatSpellPower = 5,
    /**
     * @generated from protobuf enum value: StatHealingPower = 6;
     */
    StatHealingPower = 6,
    /**
     * @generated from protobuf enum value: StatArcaneSpellPower = 7;
     */
    StatArcaneSpellPower = 7,
    /**
     * @generated from protobuf enum value: StatFireSpellPower = 8;
     */
    StatFireSpellPower = 8,
    /**
     * @generated from protobuf enum value: StatFrostSpellPower = 9;
     */
    StatFrostSpellPower = 9,
    /**
     * @generated from protobuf enum value: StatHolySpellPower = 10;
     */
    StatHolySpellPower = 10,
    /**
     * @generated from protobuf enum value: StatNatureSpellPower = 11;
     */
    StatNatureSpellPower = 11,
    /**
     * @generated from protobuf enum value: StatShadowSpellPower = 12;
     */
    StatShadowSpellPower = 12,
    /**
     * @generated from protobuf enum value: StatMP5 = 13;
     */
    StatMP5 = 13,
    /**
     * @generated from protobuf enum value: StatSpellHit = 14;
     */
    StatSpellHit = 14,
    /**
     * @generated from protobuf enum value: StatSpellCrit = 15;
     */
    StatSpellCrit = 15,
    /**
     * @generated from protobuf enum value: StatSpellHaste = 16;
     */
    StatSpellHaste = 16,
    /**
     * @generated from protobuf enum value: StatSpellPenetration = 17;
     */
    StatSpellPenetration = 17,
    /**
     * @generated from protobuf enum value: StatAttackPower = 18;
     */
    StatAttackPower = 18,
    /**
     * @generated from protobuf enum value: StatMeleeHit = 19;
     */
    StatMeleeHit = 19,
    /**
     * @generated from protobuf enum value: StatMeleeCrit = 20;
     */
    StatMeleeCrit = 20,
    /**
     * @generated from protobuf enum value: StatMeleeHaste = 21;
     */
    StatMeleeHaste = 21,
    /**
     * @generated from protobuf enum value: StatArmorPenetration = 22;
     */
    StatArmorPenetration = 22,
    /**
     * @generated from protobuf enum value: StatExpertise = 23;
     */
    StatExpertise = 23,
    /**
     * @generated from protobuf enum value: StatMana = 24;
     */
    StatMana = 24,
    /**
     * @generated from protobuf enum value: StatEnergy = 25;
     */
    StatEnergy = 25,
    /**
     * @generated from protobuf enum value: StatRage = 26;
     */
    StatRage = 26,
    /**
     * @generated from protobuf enum value: StatArmor = 27;
     */
    StatArmor = 27,
    /**
     * @generated from protobuf enum value: StatRangedAttackPower = 28;
     */
    StatRangedAttackPower = 28,
    /**
     * @generated from protobuf enum value: StatDefense = 29;
     */
    StatDefense = 29,
    /**
     * @generated from protobuf enum value: StatBlock = 30;
     */
    StatBlock = 30,
    /**
     * @generated from protobuf enum value: StatBlockValue = 31;
     */
    StatBlockValue = 31,
    /**
     * @generated from protobuf enum value: StatDodge = 32;
     */
    StatDodge = 32,
    /**
     * @generated from protobuf enum value: StatParry = 33;
     */
    StatParry = 33,
    /**
     * @generated from protobuf enum value: StatResilience = 34;
     */
    StatResilience = 34,
    /**
     * @generated from protobuf enum value: StatHealth = 35;
     */
    StatHealth = 35,
    /**
     * @generated from protobuf enum value: StatArcaneResistance = 36;
     */
    StatArcaneResistance = 36,
    /**
     * @generated from protobuf enum value: StatFireResistance = 37;
     */
    StatFireResistance = 37,
    /**
     * @generated from protobuf enum value: StatFrostResistance = 38;
     */
    StatFrostResistance = 38,
    /**
     * @generated from protobuf enum value: StatNatureResistance = 39;
     */
    StatNatureResistance = 39,
    /**
     * @generated from protobuf enum value: StatShadowResistance = 40;
     */
    StatShadowResistance = 40
}
/**
 * @generated from protobuf enum proto.ItemType
 */
export declare enum ItemType {
    /**
     * @generated from protobuf enum value: ItemTypeUnknown = 0;
     */
    ItemTypeUnknown = 0,
    /**
     * @generated from protobuf enum value: ItemTypeHead = 1;
     */
    ItemTypeHead = 1,
    /**
     * @generated from protobuf enum value: ItemTypeNeck = 2;
     */
    ItemTypeNeck = 2,
    /**
     * @generated from protobuf enum value: ItemTypeShoulder = 3;
     */
    ItemTypeShoulder = 3,
    /**
     * @generated from protobuf enum value: ItemTypeBack = 4;
     */
    ItemTypeBack = 4,
    /**
     * @generated from protobuf enum value: ItemTypeChest = 5;
     */
    ItemTypeChest = 5,
    /**
     * @generated from protobuf enum value: ItemTypeWrist = 6;
     */
    ItemTypeWrist = 6,
    /**
     * @generated from protobuf enum value: ItemTypeHands = 7;
     */
    ItemTypeHands = 7,
    /**
     * @generated from protobuf enum value: ItemTypeWaist = 8;
     */
    ItemTypeWaist = 8,
    /**
     * @generated from protobuf enum value: ItemTypeLegs = 9;
     */
    ItemTypeLegs = 9,
    /**
     * @generated from protobuf enum value: ItemTypeFeet = 10;
     */
    ItemTypeFeet = 10,
    /**
     * @generated from protobuf enum value: ItemTypeFinger = 11;
     */
    ItemTypeFinger = 11,
    /**
     * @generated from protobuf enum value: ItemTypeTrinket = 12;
     */
    ItemTypeTrinket = 12,
    /**
     * @generated from protobuf enum value: ItemTypeWeapon = 13;
     */
    ItemTypeWeapon = 13,
    /**
     * @generated from protobuf enum value: ItemTypeRanged = 14;
     */
    ItemTypeRanged = 14
}
/**
 * @generated from protobuf enum proto.ArmorType
 */
export declare enum ArmorType {
    /**
     * @generated from protobuf enum value: ArmorTypeUnknown = 0;
     */
    ArmorTypeUnknown = 0,
    /**
     * @generated from protobuf enum value: ArmorTypeCloth = 1;
     */
    ArmorTypeCloth = 1,
    /**
     * @generated from protobuf enum value: ArmorTypeLeather = 2;
     */
    ArmorTypeLeather = 2,
    /**
     * @generated from protobuf enum value: ArmorTypeMail = 3;
     */
    ArmorTypeMail = 3,
    /**
     * @generated from protobuf enum value: ArmorTypePlate = 4;
     */
    ArmorTypePlate = 4
}
/**
 * @generated from protobuf enum proto.WeaponType
 */
export declare enum WeaponType {
    /**
     * @generated from protobuf enum value: WeaponTypeUnknown = 0;
     */
    WeaponTypeUnknown = 0,
    /**
     * @generated from protobuf enum value: WeaponTypeAxe = 1;
     */
    WeaponTypeAxe = 1,
    /**
     * @generated from protobuf enum value: WeaponTypeDagger = 2;
     */
    WeaponTypeDagger = 2,
    /**
     * @generated from protobuf enum value: WeaponTypeFist = 3;
     */
    WeaponTypeFist = 3,
    /**
     * @generated from protobuf enum value: WeaponTypeMace = 4;
     */
    WeaponTypeMace = 4,
    /**
     * @generated from protobuf enum value: WeaponTypeOffHand = 5;
     */
    WeaponTypeOffHand = 5,
    /**
     * @generated from protobuf enum value: WeaponTypePolearm = 6;
     */
    WeaponTypePolearm = 6,
    /**
     * @generated from protobuf enum value: WeaponTypeShield = 7;
     */
    WeaponTypeShield = 7,
    /**
     * @generated from protobuf enum value: WeaponTypeStaff = 8;
     */
    WeaponTypeStaff = 8,
    /**
     * @generated from protobuf enum value: WeaponTypeSword = 9;
     */
    WeaponTypeSword = 9
}
/**
 * @generated from protobuf enum proto.HandType
 */
export declare enum HandType {
    /**
     * @generated from protobuf enum value: HandTypeUnknown = 0;
     */
    HandTypeUnknown = 0,
    /**
     * @generated from protobuf enum value: HandTypeMainHand = 1;
     */
    HandTypeMainHand = 1,
    /**
     * @generated from protobuf enum value: HandTypeOneHand = 2;
     */
    HandTypeOneHand = 2,
    /**
     * @generated from protobuf enum value: HandTypeOffHand = 3;
     */
    HandTypeOffHand = 3,
    /**
     * @generated from protobuf enum value: HandTypeTwoHand = 4;
     */
    HandTypeTwoHand = 4
}
/**
 * @generated from protobuf enum proto.RangedWeaponType
 */
export declare enum RangedWeaponType {
    /**
     * @generated from protobuf enum value: RangedWeaponTypeUnknown = 0;
     */
    RangedWeaponTypeUnknown = 0,
    /**
     * @generated from protobuf enum value: RangedWeaponTypeBow = 1;
     */
    RangedWeaponTypeBow = 1,
    /**
     * @generated from protobuf enum value: RangedWeaponTypeCrossbow = 2;
     */
    RangedWeaponTypeCrossbow = 2,
    /**
     * @generated from protobuf enum value: RangedWeaponTypeGun = 3;
     */
    RangedWeaponTypeGun = 3,
    /**
     * @generated from protobuf enum value: RangedWeaponTypeIdol = 4;
     */
    RangedWeaponTypeIdol = 4,
    /**
     * @generated from protobuf enum value: RangedWeaponTypeLibram = 5;
     */
    RangedWeaponTypeLibram = 5,
    /**
     * @generated from protobuf enum value: RangedWeaponTypeThrown = 6;
     */
    RangedWeaponTypeThrown = 6,
    /**
     * @generated from protobuf enum value: RangedWeaponTypeTotem = 7;
     */
    RangedWeaponTypeTotem = 7,
    /**
     * @generated from protobuf enum value: RangedWeaponTypeWand = 8;
     */
    RangedWeaponTypeWand = 8,
    /**
     * @generated from protobuf enum value: RangedWeaponTypeSigil = 9;
     */
    RangedWeaponTypeSigil = 9
}
/**
 * All slots on the gear menu where a single item can be worn.
 *
 * @generated from protobuf enum proto.ItemSlot
 */
export declare enum ItemSlot {
    /**
     * @generated from protobuf enum value: ItemSlotHead = 0;
     */
    ItemSlotHead = 0,
    /**
     * @generated from protobuf enum value: ItemSlotNeck = 1;
     */
    ItemSlotNeck = 1,
    /**
     * @generated from protobuf enum value: ItemSlotShoulder = 2;
     */
    ItemSlotShoulder = 2,
    /**
     * @generated from protobuf enum value: ItemSlotBack = 3;
     */
    ItemSlotBack = 3,
    /**
     * @generated from protobuf enum value: ItemSlotChest = 4;
     */
    ItemSlotChest = 4,
    /**
     * @generated from protobuf enum value: ItemSlotWrist = 5;
     */
    ItemSlotWrist = 5,
    /**
     * @generated from protobuf enum value: ItemSlotHands = 6;
     */
    ItemSlotHands = 6,
    /**
     * @generated from protobuf enum value: ItemSlotWaist = 7;
     */
    ItemSlotWaist = 7,
    /**
     * @generated from protobuf enum value: ItemSlotLegs = 8;
     */
    ItemSlotLegs = 8,
    /**
     * @generated from protobuf enum value: ItemSlotFeet = 9;
     */
    ItemSlotFeet = 9,
    /**
     * @generated from protobuf enum value: ItemSlotFinger1 = 10;
     */
    ItemSlotFinger1 = 10,
    /**
     * @generated from protobuf enum value: ItemSlotFinger2 = 11;
     */
    ItemSlotFinger2 = 11,
    /**
     * @generated from protobuf enum value: ItemSlotTrinket1 = 12;
     */
    ItemSlotTrinket1 = 12,
    /**
     * @generated from protobuf enum value: ItemSlotTrinket2 = 13;
     */
    ItemSlotTrinket2 = 13,
    /**
     * can be 1h or 2h
     *
     * @generated from protobuf enum value: ItemSlotMainHand = 14;
     */
    ItemSlotMainHand = 14,
    /**
     * @generated from protobuf enum value: ItemSlotOffHand = 15;
     */
    ItemSlotOffHand = 15,
    /**
     * @generated from protobuf enum value: ItemSlotRanged = 16;
     */
    ItemSlotRanged = 16
}
/**
 * @generated from protobuf enum proto.ItemQuality
 */
export declare enum ItemQuality {
    /**
     * @generated from protobuf enum value: ItemQualityJunk = 0;
     */
    ItemQualityJunk = 0,
    /**
     * @generated from protobuf enum value: ItemQualityCommon = 1;
     */
    ItemQualityCommon = 1,
    /**
     * @generated from protobuf enum value: ItemQualityUncommon = 2;
     */
    ItemQualityUncommon = 2,
    /**
     * @generated from protobuf enum value: ItemQualityRare = 3;
     */
    ItemQualityRare = 3,
    /**
     * @generated from protobuf enum value: ItemQualityEpic = 4;
     */
    ItemQualityEpic = 4,
    /**
     * @generated from protobuf enum value: ItemQualityLegendary = 5;
     */
    ItemQualityLegendary = 5
}
/**
 * @generated from protobuf enum proto.GemColor
 */
export declare enum GemColor {
    /**
     * @generated from protobuf enum value: GemColorUnknown = 0;
     */
    GemColorUnknown = 0,
    /**
     * @generated from protobuf enum value: GemColorMeta = 1;
     */
    GemColorMeta = 1,
    /**
     * @generated from protobuf enum value: GemColorRed = 2;
     */
    GemColorRed = 2,
    /**
     * @generated from protobuf enum value: GemColorBlue = 3;
     */
    GemColorBlue = 3,
    /**
     * @generated from protobuf enum value: GemColorYellow = 4;
     */
    GemColorYellow = 4,
    /**
     * @generated from protobuf enum value: GemColorGreen = 5;
     */
    GemColorGreen = 5,
    /**
     * @generated from protobuf enum value: GemColorOrange = 6;
     */
    GemColorOrange = 6,
    /**
     * @generated from protobuf enum value: GemColorPurple = 7;
     */
    GemColorPurple = 7,
    /**
     * @generated from protobuf enum value: GemColorPrismatic = 8;
     */
    GemColorPrismatic = 8
}
/**
 * @generated from protobuf enum proto.SpellSchool
 */
export declare enum SpellSchool {
    /**
     * @generated from protobuf enum value: SpellSchoolPhysical = 0;
     */
    SpellSchoolPhysical = 0,
    /**
     * @generated from protobuf enum value: SpellSchoolArcane = 1;
     */
    SpellSchoolArcane = 1,
    /**
     * @generated from protobuf enum value: SpellSchoolFire = 2;
     */
    SpellSchoolFire = 2,
    /**
     * @generated from protobuf enum value: SpellSchoolFrost = 3;
     */
    SpellSchoolFrost = 3,
    /**
     * @generated from protobuf enum value: SpellSchoolHoly = 4;
     */
    SpellSchoolHoly = 4,
    /**
     * @generated from protobuf enum value: SpellSchoolNature = 5;
     */
    SpellSchoolNature = 5,
    /**
     * @generated from protobuf enum value: SpellSchoolShadow = 6;
     */
    SpellSchoolShadow = 6
}
/**
 * @generated from protobuf enum proto.TristateEffect
 */
export declare enum TristateEffect {
    /**
     * @generated from protobuf enum value: TristateEffectMissing = 0;
     */
    TristateEffectMissing = 0,
    /**
     * @generated from protobuf enum value: TristateEffectRegular = 1;
     */
    TristateEffectRegular = 1,
    /**
     * @generated from protobuf enum value: TristateEffectImproved = 2;
     */
    TristateEffectImproved = 2
}
/**
 * @generated from protobuf enum proto.Explosive
 */
export declare enum Explosive {
    /**
     * @generated from protobuf enum value: ExplosiveUnknown = 0;
     */
    ExplosiveUnknown = 0,
    /**
     * @generated from protobuf enum value: ExplosiveFelIronBomb = 1;
     */
    ExplosiveFelIronBomb = 1,
    /**
     * @generated from protobuf enum value: ExplosiveAdamantiteGrenade = 2;
     */
    ExplosiveAdamantiteGrenade = 2,
    /**
     * @generated from protobuf enum value: ExplosiveGnomishFlameTurret = 3;
     */
    ExplosiveGnomishFlameTurret = 3,
    /**
     * @generated from protobuf enum value: ExplosiveHolyWater = 4;
     */
    ExplosiveHolyWater = 4
}
/**
 * @generated from protobuf enum proto.Potions
 */
export declare enum Potions {
    /**
     * @generated from protobuf enum value: UnknownPotion = 0;
     */
    UnknownPotion = 0,
    /**
     * @generated from protobuf enum value: RunicHealingPotion = 1;
     */
    RunicHealingPotion = 1,
    /**
     * @generated from protobuf enum value: RunicManaPotion = 2;
     */
    RunicManaPotion = 2,
    /**
     * @generated from protobuf enum value: IndestructiblePotion = 3;
     */
    IndestructiblePotion = 3,
    /**
     * @generated from protobuf enum value: PotionOfSpeed = 4;
     */
    PotionOfSpeed = 4,
    /**
     * @generated from protobuf enum value: PotionOfWildMagic = 5;
     */
    PotionOfWildMagic = 5,
    /**
     * @generated from protobuf enum value: DestructionPotion = 6;
     */
    DestructionPotion = 6,
    /**
     * @generated from protobuf enum value: SuperManaPotion = 7;
     */
    SuperManaPotion = 7,
    /**
     * @generated from protobuf enum value: HastePotion = 8;
     */
    HastePotion = 8,
    /**
     * @generated from protobuf enum value: MightyRagePotion = 9;
     */
    MightyRagePotion = 9,
    /**
     * @generated from protobuf enum value: FelManaPotion = 10;
     */
    FelManaPotion = 10,
    /**
     * @generated from protobuf enum value: InsaneStrengthPotion = 11;
     */
    InsaneStrengthPotion = 11,
    /**
     * @generated from protobuf enum value: IronshieldPotion = 12;
     */
    IronshieldPotion = 12,
    /**
     * @generated from protobuf enum value: HeroicPotion = 13;
     */
    HeroicPotion = 13
}
/**
 * @generated from protobuf enum proto.Conjured
 */
export declare enum Conjured {
    /**
     * @generated from protobuf enum value: ConjuredUnknown = 0;
     */
    ConjuredUnknown = 0,
    /**
     * @generated from protobuf enum value: ConjuredDarkRune = 1;
     */
    ConjuredDarkRune = 1,
    /**
     * @generated from protobuf enum value: ConjuredFlameCap = 2;
     */
    ConjuredFlameCap = 2,
    /**
     * @generated from protobuf enum value: ConjuredHealthstone = 5;
     */
    ConjuredHealthstone = 5,
    /**
     * @generated from protobuf enum value: ConjuredMageManaEmerald = 3;
     */
    ConjuredMageManaEmerald = 3,
    /**
     * @generated from protobuf enum value: ConjuredRogueThistleTea = 4;
     */
    ConjuredRogueThistleTea = 4
}
/**
 * @generated from protobuf enum proto.WeaponImbue
 */
export declare enum WeaponImbue {
    /**
     * @generated from protobuf enum value: WeaponImbueUnknown = 0;
     */
    WeaponImbueUnknown = 0,
    /**
     * @generated from protobuf enum value: WeaponImbueAdamantiteSharpeningStone = 1;
     */
    WeaponImbueAdamantiteSharpeningStone = 1,
    /**
     * @generated from protobuf enum value: WeaponImbueAdamantiteWeightstone = 5;
     */
    WeaponImbueAdamantiteWeightstone = 5,
    /**
     * @generated from protobuf enum value: WeaponImbueElementalSharpeningStone = 2;
     */
    WeaponImbueElementalSharpeningStone = 2,
    /**
     * @generated from protobuf enum value: WeaponImbueBrilliantWizardOil = 3;
     */
    WeaponImbueBrilliantWizardOil = 3,
    /**
     * @generated from protobuf enum value: WeaponImbueSuperiorWizardOil = 4;
     */
    WeaponImbueSuperiorWizardOil = 4,
    /**
     * @generated from protobuf enum value: WeaponImbueShamanFlametongue = 6;
     */
    WeaponImbueShamanFlametongue = 6,
    /**
     * @generated from protobuf enum value: WeaponImbueShamanFrostbrand = 7;
     */
    WeaponImbueShamanFrostbrand = 7,
    /**
     * @generated from protobuf enum value: WeaponImbueShamanRockbiter = 8;
     */
    WeaponImbueShamanRockbiter = 8,
    /**
     * @generated from protobuf enum value: WeaponImbueShamanWindfury = 9;
     */
    WeaponImbueShamanWindfury = 9,
    /**
     * @generated from protobuf enum value: WeaponImbueRogueDeadlyPoison = 10;
     */
    WeaponImbueRogueDeadlyPoison = 10,
    /**
     * @generated from protobuf enum value: WeaponImbueRogueInstantPoison = 11;
     */
    WeaponImbueRogueInstantPoison = 11,
    /**
     * @generated from protobuf enum value: WeaponImbueRighteousWeaponCoating = 12;
     */
    WeaponImbueRighteousWeaponCoating = 12
}
/**
 * @generated from protobuf enum proto.Flask
 */
export declare enum Flask {
    /**
     * @generated from protobuf enum value: FlaskUnknown = 0;
     */
    FlaskUnknown = 0,
    /**
     * @generated from protobuf enum value: FlaskOfTheFrostWyrm = 1;
     */
    FlaskOfTheFrostWyrm = 1,
    /**
     * @generated from protobuf enum value: FlaskOfEndlessRage = 2;
     */
    FlaskOfEndlessRage = 2,
    /**
     * @generated from protobuf enum value: FlaskOfPureMojo = 3;
     */
    FlaskOfPureMojo = 3,
    /**
     * @generated from protobuf enum value: FlaskOfStoneblood = 4;
     */
    FlaskOfStoneblood = 4,
    /**
     * @generated from protobuf enum value: LesserFlaskOfToughness = 5;
     */
    LesserFlaskOfToughness = 5,
    /**
     * @generated from protobuf enum value: LesserFlaskOfResistance = 6;
     */
    LesserFlaskOfResistance = 6,
    /**
     * TBC
     *
     * @generated from protobuf enum value: FlaskOfBlindingLight = 7;
     */
    FlaskOfBlindingLight = 7,
    /**
     * @generated from protobuf enum value: FlaskOfMightyRestoration = 8;
     */
    FlaskOfMightyRestoration = 8,
    /**
     * @generated from protobuf enum value: FlaskOfPureDeath = 9;
     */
    FlaskOfPureDeath = 9,
    /**
     * @generated from protobuf enum value: FlaskOfRelentlessAssault = 10;
     */
    FlaskOfRelentlessAssault = 10,
    /**
     * @generated from protobuf enum value: FlaskOfSupremePower = 11;
     */
    FlaskOfSupremePower = 11,
    /**
     * @generated from protobuf enum value: FlaskOfFortification = 12;
     */
    FlaskOfFortification = 12,
    /**
     * @generated from protobuf enum value: FlaskOfChromaticWonder = 13;
     */
    FlaskOfChromaticWonder = 13
}
/**
 * @generated from protobuf enum proto.BattleElixir
 */
export declare enum BattleElixir {
    /**
     * @generated from protobuf enum value: BattleElixirUnknown = 0;
     */
    BattleElixirUnknown = 0,
    /**
     * @generated from protobuf enum value: ElixirOfAccuracy = 1;
     */
    ElixirOfAccuracy = 1,
    /**
     * @generated from protobuf enum value: ElixirOfArmorPiercing = 2;
     */
    ElixirOfArmorPiercing = 2,
    /**
     * @generated from protobuf enum value: ElixirOfDeadlyStrikes = 3;
     */
    ElixirOfDeadlyStrikes = 3,
    /**
     * @generated from protobuf enum value: ElixirOfExpertise = 4;
     */
    ElixirOfExpertise = 4,
    /**
     * @generated from protobuf enum value: ElixirOfLightningSpeed = 5;
     */
    ElixirOfLightningSpeed = 5,
    /**
     * @generated from protobuf enum value: ElixirOfMightyAgility = 6;
     */
    ElixirOfMightyAgility = 6,
    /**
     * @generated from protobuf enum value: ElixirOfMightyStrength = 7;
     */
    ElixirOfMightyStrength = 7,
    /**
     * @generated from protobuf enum value: GurusElixir = 8;
     */
    GurusElixir = 8,
    /**
     * @generated from protobuf enum value: SpellpowerElixir = 9;
     */
    SpellpowerElixir = 9,
    /**
     * @generated from protobuf enum value: WrathElixir = 10;
     */
    WrathElixir = 10,
    /**
     * TBC
     *
     * @generated from protobuf enum value: AdeptsElixir = 11;
     */
    AdeptsElixir = 11,
    /**
     * @generated from protobuf enum value: ElixirOfDemonslaying = 12;
     */
    ElixirOfDemonslaying = 12,
    /**
     * @generated from protobuf enum value: ElixirOfMajorAgility = 13;
     */
    ElixirOfMajorAgility = 13,
    /**
     * @generated from protobuf enum value: ElixirOfMajorFirePower = 14;
     */
    ElixirOfMajorFirePower = 14,
    /**
     * @generated from protobuf enum value: ElixirOfMajorFrostPower = 15;
     */
    ElixirOfMajorFrostPower = 15,
    /**
     * @generated from protobuf enum value: ElixirOfMajorShadowPower = 16;
     */
    ElixirOfMajorShadowPower = 16,
    /**
     * @generated from protobuf enum value: ElixirOfMajorStrength = 17;
     */
    ElixirOfMajorStrength = 17,
    /**
     * @generated from protobuf enum value: ElixirOfMastery = 18;
     */
    ElixirOfMastery = 18,
    /**
     * @generated from protobuf enum value: ElixirOfTheMongoose = 19;
     */
    ElixirOfTheMongoose = 19,
    /**
     * @generated from protobuf enum value: FelStrengthElixir = 20;
     */
    FelStrengthElixir = 20,
    /**
     * @generated from protobuf enum value: GreaterArcaneElixir = 21;
     */
    GreaterArcaneElixir = 21
}
/**
 * @generated from protobuf enum proto.GuardianElixir
 */
export declare enum GuardianElixir {
    /**
     * @generated from protobuf enum value: GuardianElixirUnknown = 0;
     */
    GuardianElixirUnknown = 0,
    /**
     * @generated from protobuf enum value: ElixirOfMightyDefense = 1;
     */
    ElixirOfMightyDefense = 1,
    /**
     * @generated from protobuf enum value: ElixirOfMightyFortitude = 2;
     */
    ElixirOfMightyFortitude = 2,
    /**
     * @generated from protobuf enum value: ElixirOfMightyMageblood = 3;
     */
    ElixirOfMightyMageblood = 3,
    /**
     * @generated from protobuf enum value: ElixirOfMightyThoughts = 4;
     */
    ElixirOfMightyThoughts = 4,
    /**
     * @generated from protobuf enum value: ElixirOfProtection = 5;
     */
    ElixirOfProtection = 5,
    /**
     * @generated from protobuf enum value: ElixirOfSpirit = 6;
     */
    ElixirOfSpirit = 6,
    /**
     * TBC
     *
     * @generated from protobuf enum value: GiftOfArthas = 7;
     */
    GiftOfArthas = 7,
    /**
     * @generated from protobuf enum value: ElixirOfDraenicWisdom = 8;
     */
    ElixirOfDraenicWisdom = 8,
    /**
     * @generated from protobuf enum value: ElixirOfIronskin = 9;
     */
    ElixirOfIronskin = 9,
    /**
     * @generated from protobuf enum value: ElixirOfMajorDefense = 10;
     */
    ElixirOfMajorDefense = 10,
    /**
     * @generated from protobuf enum value: ElixirOfMajorFortitude = 11;
     */
    ElixirOfMajorFortitude = 11,
    /**
     * @generated from protobuf enum value: ElixirOfMajorMageblood = 12;
     */
    ElixirOfMajorMageblood = 12
}
/**
 * @generated from protobuf enum proto.Food
 */
export declare enum Food {
    /**
     * @generated from protobuf enum value: FoodUnknown = 0;
     */
    FoodUnknown = 0,
    /**
     * @generated from protobuf enum value: FoodFishFeast = 1;
     */
    FoodFishFeast = 1,
    /**
     * @generated from protobuf enum value: FoodGreatFeast = 2;
     */
    FoodGreatFeast = 2,
    /**
     * @generated from protobuf enum value: FoodBlackenedDragonfin = 3;
     */
    FoodBlackenedDragonfin = 3,
    /**
     * @generated from protobuf enum value: FoodHeartyRhino = 4;
     */
    FoodHeartyRhino = 4,
    /**
     * @generated from protobuf enum value: FoodMegaMammothMeal = 5;
     */
    FoodMegaMammothMeal = 5,
    /**
     * @generated from protobuf enum value: FoodSpicedWormBurger = 6;
     */
    FoodSpicedWormBurger = 6,
    /**
     * @generated from protobuf enum value: FoodRhinoliciousWormsteak = 7;
     */
    FoodRhinoliciousWormsteak = 7,
    /**
     * @generated from protobuf enum value: FoodImperialMantaSteak = 8;
     */
    FoodImperialMantaSteak = 8,
    /**
     * @generated from protobuf enum value: FoodSnapperExtreme = 9;
     */
    FoodSnapperExtreme = 9,
    /**
     * @generated from protobuf enum value: FoodMightyRhinoDogs = 10;
     */
    FoodMightyRhinoDogs = 10,
    /**
     * @generated from protobuf enum value: FoodFirecrackerSalmon = 11;
     */
    FoodFirecrackerSalmon = 11,
    /**
     * @generated from protobuf enum value: FoodCuttlesteak = 12;
     */
    FoodCuttlesteak = 12,
    /**
     * @generated from protobuf enum value: FoodDragonfinFilet = 13;
     */
    FoodDragonfinFilet = 13,
    /**
     * TBC Foods
     *
     * @generated from protobuf enum value: FoodBlackenedBasilisk = 14;
     */
    FoodBlackenedBasilisk = 14,
    /**
     * @generated from protobuf enum value: FoodGrilledMudfish = 15;
     */
    FoodGrilledMudfish = 15,
    /**
     * @generated from protobuf enum value: FoodRavagerDog = 16;
     */
    FoodRavagerDog = 16,
    /**
     * @generated from protobuf enum value: FoodRoastedClefthoof = 17;
     */
    FoodRoastedClefthoof = 17,
    /**
     * @generated from protobuf enum value: FoodSkullfishSoup = 18;
     */
    FoodSkullfishSoup = 18,
    /**
     * @generated from protobuf enum value: FoodSpicyHotTalbuk = 19;
     */
    FoodSpicyHotTalbuk = 19,
    /**
     * @generated from protobuf enum value: FoodFishermansFeast = 20;
     */
    FoodFishermansFeast = 20
}
/**
 * @generated from protobuf enum proto.PetFood
 */
export declare enum PetFood {
    /**
     * @generated from protobuf enum value: PetFoodUnknown = 0;
     */
    PetFoodUnknown = 0,
    /**
     * @generated from protobuf enum value: PetFoodSpicedMammothTreats = 1;
     */
    PetFoodSpicedMammothTreats = 1,
    /**
     * TBC
     *
     * @generated from protobuf enum value: PetFoodKiblersBits = 2;
     */
    PetFoodKiblersBits = 2
}
/**
 * @generated from protobuf enum proto.MobType
 */
export declare enum MobType {
    /**
     * @generated from protobuf enum value: MobTypeUnknown = 0;
     */
    MobTypeUnknown = 0,
    /**
     * @generated from protobuf enum value: MobTypeBeast = 1;
     */
    MobTypeBeast = 1,
    /**
     * @generated from protobuf enum value: MobTypeDemon = 2;
     */
    MobTypeDemon = 2,
    /**
     * @generated from protobuf enum value: MobTypeDragonkin = 3;
     */
    MobTypeDragonkin = 3,
    /**
     * @generated from protobuf enum value: MobTypeElemental = 4;
     */
    MobTypeElemental = 4,
    /**
     * @generated from protobuf enum value: MobTypeGiant = 5;
     */
    MobTypeGiant = 5,
    /**
     * @generated from protobuf enum value: MobTypeHumanoid = 6;
     */
    MobTypeHumanoid = 6,
    /**
     * @generated from protobuf enum value: MobTypeMechanical = 7;
     */
    MobTypeMechanical = 7,
    /**
     * @generated from protobuf enum value: MobTypeUndead = 8;
     */
    MobTypeUndead = 8
}
/**
 * Extra enum for describing which items are eligible for an enchant, when
 * ItemType alone is not enough.
 *
 * @generated from protobuf enum proto.EnchantType
 */
export declare enum EnchantType {
    /**
     * @generated from protobuf enum value: EnchantTypeNormal = 0;
     */
    EnchantTypeNormal = 0,
    /**
     * @generated from protobuf enum value: EnchantTypeTwoHand = 1;
     */
    EnchantTypeTwoHand = 1,
    /**
     * @generated from protobuf enum value: EnchantTypeShield = 2;
     */
    EnchantTypeShield = 2,
    /**
     * @generated from protobuf enum value: EnchantTypeKit = 3;
     */
    EnchantTypeKit = 3
}
/**
 * ID for actions that aren't spells or items.
 *
 * @generated from protobuf enum proto.OtherAction
 */
export declare enum OtherAction {
    /**
     * @generated from protobuf enum value: OtherActionNone = 0;
     */
    OtherActionNone = 0,
    /**
     * @generated from protobuf enum value: OtherActionWait = 1;
     */
    OtherActionWait = 1,
    /**
     * @generated from protobuf enum value: OtherActionManaRegen = 2;
     */
    OtherActionManaRegen = 2,
    /**
     * @generated from protobuf enum value: OtherActionEnergyRegen = 5;
     */
    OtherActionEnergyRegen = 5,
    /**
     * @generated from protobuf enum value: OtherActionFocusRegen = 6;
     */
    OtherActionFocusRegen = 6,
    /**
     * For threat generated from mana gains.
     *
     * @generated from protobuf enum value: OtherActionManaGain = 10;
     */
    OtherActionManaGain = 10,
    /**
     * For threat generated from rage gains.
     *
     * @generated from protobuf enum value: OtherActionRageGain = 11;
     */
    OtherActionRageGain = 11,
    /**
     * A white hit, can be main hand or off hand.
     *
     * @generated from protobuf enum value: OtherActionAttack = 3;
     */
    OtherActionAttack = 3,
    /**
     * Default shoot action using a wand/bow/gun.
     *
     * @generated from protobuf enum value: OtherActionShoot = 4;
     */
    OtherActionShoot = 4,
    /**
     * Represents a grouping of all pet actions. Only used by the UI.
     *
     * @generated from protobuf enum value: OtherActionPet = 7;
     */
    OtherActionPet = 7,
    /**
     * Refund of a resource like Energy or Rage, when the ability didn't land.
     *
     * @generated from protobuf enum value: OtherActionRefund = 8;
     */
    OtherActionRefund = 8,
    /**
     * Indicates damage taken; used for rage gen.
     *
     * @generated from protobuf enum value: OtherActionDamageTaken = 9;
     */
    OtherActionDamageTaken = 9,
    /**
     * Indicates healing received from healing model.
     *
     * @generated from protobuf enum value: OtherActionHealingModel = 12;
     */
    OtherActionHealingModel = 12,
    /**
     * Indicates healing received from healing model.
     *
     * @generated from protobuf enum value: OtherActionBloodRuneGain = 13;
     */
    OtherActionBloodRuneGain = 13,
    /**
     * Indicates healing received from healing model.
     *
     * @generated from protobuf enum value: OtherActionFrostRuneGain = 14;
     */
    OtherActionFrostRuneGain = 14,
    /**
     * Indicates healing received from healing model.
     *
     * @generated from protobuf enum value: OtherActionUnholyRuneGain = 15;
     */
    OtherActionUnholyRuneGain = 15,
    /**
     * Indicates healing received from healing model.
     *
     * @generated from protobuf enum value: OtherActionDeathRuneGain = 16;
     */
    OtherActionDeathRuneGain = 16
}
declare class RaidBuffs$Type extends MessageType<RaidBuffs> {
    constructor();
    create(value?: PartialMessage<RaidBuffs>): RaidBuffs;
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: RaidBuffs): RaidBuffs;
    internalBinaryWrite(message: RaidBuffs, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter;
}
/**
 * @generated MessageType for protobuf message proto.RaidBuffs
 */
export declare const RaidBuffs: RaidBuffs$Type;
declare class PartyBuffs$Type extends MessageType<PartyBuffs> {
    constructor();
    create(value?: PartialMessage<PartyBuffs>): PartyBuffs;
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: PartyBuffs): PartyBuffs;
    internalBinaryWrite(message: PartyBuffs, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter;
}
/**
 * @generated MessageType for protobuf message proto.PartyBuffs
 */
export declare const PartyBuffs: PartyBuffs$Type;
declare class IndividualBuffs$Type extends MessageType<IndividualBuffs> {
    constructor();
    create(value?: PartialMessage<IndividualBuffs>): IndividualBuffs;
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: IndividualBuffs): IndividualBuffs;
    internalBinaryWrite(message: IndividualBuffs, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter;
}
/**
 * @generated MessageType for protobuf message proto.IndividualBuffs
 */
export declare const IndividualBuffs: IndividualBuffs$Type;
declare class Consumes$Type extends MessageType<Consumes> {
    constructor();
    create(value?: PartialMessage<Consumes>): Consumes;
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: Consumes): Consumes;
    internalBinaryWrite(message: Consumes, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter;
}
/**
 * @generated MessageType for protobuf message proto.Consumes
 */
export declare const Consumes: Consumes$Type;
declare class Debuffs$Type extends MessageType<Debuffs> {
    constructor();
    create(value?: PartialMessage<Debuffs>): Debuffs;
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: Debuffs): Debuffs;
    internalBinaryWrite(message: Debuffs, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter;
}
/**
 * @generated MessageType for protobuf message proto.Debuffs
 */
export declare const Debuffs: Debuffs$Type;
declare class Target$Type extends MessageType<Target> {
    constructor();
    create(value?: PartialMessage<Target>): Target;
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: Target): Target;
    internalBinaryWrite(message: Target, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter;
}
/**
 * @generated MessageType for protobuf message proto.Target
 */
export declare const Target: Target$Type;
declare class Encounter$Type extends MessageType<Encounter> {
    constructor();
    create(value?: PartialMessage<Encounter>): Encounter;
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: Encounter): Encounter;
    internalBinaryWrite(message: Encounter, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter;
}
/**
 * @generated MessageType for protobuf message proto.Encounter
 */
export declare const Encounter: Encounter$Type;
declare class ItemSpec$Type extends MessageType<ItemSpec> {
    constructor();
    create(value?: PartialMessage<ItemSpec>): ItemSpec;
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: ItemSpec): ItemSpec;
    internalBinaryWrite(message: ItemSpec, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter;
}
/**
 * @generated MessageType for protobuf message proto.ItemSpec
 */
export declare const ItemSpec: ItemSpec$Type;
declare class EquipmentSpec$Type extends MessageType<EquipmentSpec> {
    constructor();
    create(value?: PartialMessage<EquipmentSpec>): EquipmentSpec;
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: EquipmentSpec): EquipmentSpec;
    internalBinaryWrite(message: EquipmentSpec, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter;
}
/**
 * @generated MessageType for protobuf message proto.EquipmentSpec
 */
export declare const EquipmentSpec: EquipmentSpec$Type;
declare class Item$Type extends MessageType<Item> {
    constructor();
    create(value?: PartialMessage<Item>): Item;
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: Item): Item;
    internalBinaryWrite(message: Item, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter;
}
/**
 * @generated MessageType for protobuf message proto.Item
 */
export declare const Item: Item$Type;
declare class Enchant$Type extends MessageType<Enchant> {
    constructor();
    create(value?: PartialMessage<Enchant>): Enchant;
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: Enchant): Enchant;
    internalBinaryWrite(message: Enchant, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter;
}
/**
 * @generated MessageType for protobuf message proto.Enchant
 */
export declare const Enchant: Enchant$Type;
declare class Gem$Type extends MessageType<Gem> {
    constructor();
    create(value?: PartialMessage<Gem>): Gem;
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: Gem): Gem;
    internalBinaryWrite(message: Gem, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter;
}
/**
 * @generated MessageType for protobuf message proto.Gem
 */
export declare const Gem: Gem$Type;
declare class RaidTarget$Type extends MessageType<RaidTarget> {
    constructor();
    create(value?: PartialMessage<RaidTarget>): RaidTarget;
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: RaidTarget): RaidTarget;
    internalBinaryWrite(message: RaidTarget, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter;
}
/**
 * @generated MessageType for protobuf message proto.RaidTarget
 */
export declare const RaidTarget: RaidTarget$Type;
declare class ActionID$Type extends MessageType<ActionID> {
    constructor();
    create(value?: PartialMessage<ActionID>): ActionID;
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: ActionID): ActionID;
    internalBinaryWrite(message: ActionID, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter;
}
/**
 * @generated MessageType for protobuf message proto.ActionID
 */
export declare const ActionID: ActionID$Type;
declare class Glyphs$Type extends MessageType<Glyphs> {
    constructor();
    create(value?: PartialMessage<Glyphs>): Glyphs;
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: Glyphs): Glyphs;
    internalBinaryWrite(message: Glyphs, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter;
}
/**
 * @generated MessageType for protobuf message proto.Glyphs
 */
export declare const Glyphs: Glyphs$Type;
declare class Cooldown$Type extends MessageType<Cooldown> {
    constructor();
    create(value?: PartialMessage<Cooldown>): Cooldown;
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: Cooldown): Cooldown;
    internalBinaryWrite(message: Cooldown, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter;
}
/**
 * @generated MessageType for protobuf message proto.Cooldown
 */
export declare const Cooldown: Cooldown$Type;
declare class Cooldowns$Type extends MessageType<Cooldowns> {
    constructor();
    create(value?: PartialMessage<Cooldowns>): Cooldowns;
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: Cooldowns): Cooldowns;
    internalBinaryWrite(message: Cooldowns, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter;
}
/**
 * @generated MessageType for protobuf message proto.Cooldowns
 */
export declare const Cooldowns: Cooldowns$Type;
declare class HealingModel$Type extends MessageType<HealingModel> {
    constructor();
    create(value?: PartialMessage<HealingModel>): HealingModel;
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: HealingModel): HealingModel;
    internalBinaryWrite(message: HealingModel, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter;
}
/**
 * @generated MessageType for protobuf message proto.HealingModel
 */
export declare const HealingModel: HealingModel$Type;
export {};
