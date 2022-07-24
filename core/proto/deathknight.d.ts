import type { BinaryWriteOptions } from "@protobuf-ts/runtime";
import type { IBinaryWriter } from "@protobuf-ts/runtime";
import type { BinaryReadOptions } from "@protobuf-ts/runtime";
import type { IBinaryReader } from "@protobuf-ts/runtime";
import type { PartialMessage } from "@protobuf-ts/runtime";
import { MessageType } from "@protobuf-ts/runtime";
/**
 * @generated from protobuf message proto.DeathknightTalents
 */
export interface DeathknightTalents {
    /**
     * Blood
     *
     * @generated from protobuf field: int32 butchery = 1;
     */
    butchery: number;
    /**
     * @generated from protobuf field: int32 subversion = 2;
     */
    subversion: number;
    /**
     * @generated from protobuf field: int32 blade_barrier = 3;
     */
    bladeBarrier: number;
    /**
     * @generated from protobuf field: int32 bladed_armor = 4;
     */
    bladedArmor: number;
    /**
     * @generated from protobuf field: int32 scent_of_blood = 5;
     */
    scentOfBlood: number;
    /**
     * @generated from protobuf field: int32 two_handed_weapon_specialization = 6;
     */
    twoHandedWeaponSpecialization: number;
    /**
     * @generated from protobuf field: bool rune_tap = 7;
     */
    runeTap: boolean;
    /**
     * @generated from protobuf field: int32 dark_conviction = 8;
     */
    darkConviction: number;
    /**
     * @generated from protobuf field: int32 death_rune_mastery = 9;
     */
    deathRuneMastery: number;
    /**
     * @generated from protobuf field: int32 improved_rune_tap = 10;
     */
    improvedRuneTap: number;
    /**
     * @generated from protobuf field: int32 spell_deflection = 11;
     */
    spellDeflection: number;
    /**
     * @generated from protobuf field: int32 vendetta = 12;
     */
    vendetta: number;
    /**
     * @generated from protobuf field: int32 bloody_strikes = 13;
     */
    bloodyStrikes: number;
    /**
     * @generated from protobuf field: int32 veteran_of_the_third_war = 14;
     */
    veteranOfTheThirdWar: number;
    /**
     * @generated from protobuf field: bool mark_of_blood = 15;
     */
    markOfBlood: boolean;
    /**
     * @generated from protobuf field: int32 bloody_vengeance = 16;
     */
    bloodyVengeance: number;
    /**
     * @generated from protobuf field: int32 abominations_might = 17;
     */
    abominationsMight: number;
    /**
     * @generated from protobuf field: int32 bloodworms = 18;
     */
    bloodworms: number;
    /**
     * @generated from protobuf field: bool hysteria = 19;
     */
    hysteria: boolean;
    /**
     * @generated from protobuf field: int32 improved_blood_presence = 20;
     */
    improvedBloodPresence: number;
    /**
     * @generated from protobuf field: int32 improved_death_strike = 21;
     */
    improvedDeathStrike: number;
    /**
     * @generated from protobuf field: int32 sudden_doom = 22;
     */
    suddenDoom: number;
    /**
     * @generated from protobuf field: bool vampiric_blood = 23;
     */
    vampiricBlood: boolean;
    /**
     * @generated from protobuf field: int32 will_of_the_necropolis = 24;
     */
    willOfTheNecropolis: number;
    /**
     * @generated from protobuf field: bool heart_strike = 25;
     */
    heartStrike: boolean;
    /**
     * @generated from protobuf field: int32 might_of_mograine = 26;
     */
    mightOfMograine: number;
    /**
     * @generated from protobuf field: int32 blood_gorged = 27;
     */
    bloodGorged: number;
    /**
     * @generated from protobuf field: bool dancing_rune_weapon = 28;
     */
    dancingRuneWeapon: boolean;
    /**
     * Frost
     *
     * @generated from protobuf field: int32 improved_icy_touch = 29;
     */
    improvedIcyTouch: number;
    /**
     * @generated from protobuf field: int32 runic_power_mastery = 30;
     */
    runicPowerMastery: number;
    /**
     * @generated from protobuf field: int32 toughness = 31;
     */
    toughness: number;
    /**
     * @generated from protobuf field: int32 icy_reach = 32;
     */
    icyReach: number;
    /**
     * @generated from protobuf field: int32 black_ice = 33;
     */
    blackIce: number;
    /**
     * @generated from protobuf field: int32 nerves_of_cold_steel = 34;
     */
    nervesOfColdSteel: number;
    /**
     * @generated from protobuf field: int32 icy_talons = 35;
     */
    icyTalons: number;
    /**
     * @generated from protobuf field: bool lichborne = 36;
     */
    lichborne: boolean;
    /**
     * @generated from protobuf field: int32 annihilation = 37;
     */
    annihilation: number;
    /**
     * @generated from protobuf field: int32 killing_machine = 38;
     */
    killingMachine: number;
    /**
     * @generated from protobuf field: int32 chill_of_the_grave = 39;
     */
    chillOfTheGrave: number;
    /**
     * @generated from protobuf field: int32 endless_winter = 40;
     */
    endlessWinter: number;
    /**
     * @generated from protobuf field: int32 frigid_dreadplate = 41;
     */
    frigidDreadplate: number;
    /**
     * @generated from protobuf field: int32 glacier_rot = 42;
     */
    glacierRot: number;
    /**
     * @generated from protobuf field: bool deathchill = 43;
     */
    deathchill: boolean;
    /**
     * @generated from protobuf field: bool improved_icy_talons = 44;
     */
    improvedIcyTalons: boolean;
    /**
     * @generated from protobuf field: int32 merciless_combat = 45;
     */
    mercilessCombat: number;
    /**
     * @generated from protobuf field: int32 rime = 46;
     */
    rime: number;
    /**
     * @generated from protobuf field: int32 chilblains = 47;
     */
    chilblains: number;
    /**
     * @generated from protobuf field: bool hungering_cold = 48;
     */
    hungeringCold: boolean;
    /**
     * @generated from protobuf field: int32 improved_frost_presence = 49;
     */
    improvedFrostPresence: number;
    /**
     * @generated from protobuf field: int32 threat_of_thassarian = 50;
     */
    threatOfThassarian: number;
    /**
     * @generated from protobuf field: int32 blood_of_the_north = 51;
     */
    bloodOfTheNorth: number;
    /**
     * @generated from protobuf field: bool unbreakable_armor = 52;
     */
    unbreakableArmor: boolean;
    /**
     * @generated from protobuf field: int32 acclimation = 53;
     */
    acclimation: number;
    /**
     * @generated from protobuf field: bool frost_strike = 54;
     */
    frostStrike: boolean;
    /**
     * @generated from protobuf field: int32 guile_of_gorefiend = 55;
     */
    guileOfGorefiend: number;
    /**
     * @generated from protobuf field: int32 tundra_stalker = 56;
     */
    tundraStalker: number;
    /**
     * @generated from protobuf field: bool howling_blast = 57;
     */
    howlingBlast: boolean;
    /**
     * Unholy
     *
     * @generated from protobuf field: int32 vicious_strikes = 58;
     */
    viciousStrikes: number;
    /**
     * @generated from protobuf field: int32 virulence = 59;
     */
    virulence: number;
    /**
     * @generated from protobuf field: int32 anticipation = 60;
     */
    anticipation: number;
    /**
     * @generated from protobuf field: int32 epidemic = 61;
     */
    epidemic: number;
    /**
     * @generated from protobuf field: int32 morbidity = 62;
     */
    morbidity: number;
    /**
     * @generated from protobuf field: int32 unholy_command = 63;
     */
    unholyCommand: number;
    /**
     * @generated from protobuf field: int32 ravenous_dead = 64;
     */
    ravenousDead: number;
    /**
     * @generated from protobuf field: int32 outbreak = 65;
     */
    outbreak: number;
    /**
     * @generated from protobuf field: int32 necrosis = 66;
     */
    necrosis: number;
    /**
     * @generated from protobuf field: bool corpse_explosion = 67;
     */
    corpseExplosion: boolean;
    /**
     * @generated from protobuf field: int32 on_a_pale_horse = 68;
     */
    onAPaleHorse: number;
    /**
     * @generated from protobuf field: int32 blood_caked_blade = 69;
     */
    bloodCakedBlade: number;
    /**
     * @generated from protobuf field: int32 night_of_the_dead = 70;
     */
    nightOfTheDead: number;
    /**
     * @generated from protobuf field: bool unholy_blight = 71;
     */
    unholyBlight: boolean;
    /**
     * @generated from protobuf field: int32 impurity = 72;
     */
    impurity: number;
    /**
     * @generated from protobuf field: int32 dirge = 73;
     */
    dirge: number;
    /**
     * @generated from protobuf field: int32 desecration = 74;
     */
    desecration: number;
    /**
     * @generated from protobuf field: int32 magic_suppression = 75;
     */
    magicSuppression: number;
    /**
     * @generated from protobuf field: int32 reaping = 76;
     */
    reaping: number;
    /**
     * @generated from protobuf field: bool master_of_ghouls = 77;
     */
    masterOfGhouls: boolean;
    /**
     * @generated from protobuf field: int32 desolation = 78;
     */
    desolation: number;
    /**
     * @generated from protobuf field: bool anti_magic_zone = 79;
     */
    antiMagicZone: boolean;
    /**
     * @generated from protobuf field: int32 improved_unholy_presence = 80;
     */
    improvedUnholyPresence: number;
    /**
     * @generated from protobuf field: bool ghoul_frenzy = 81;
     */
    ghoulFrenzy: boolean;
    /**
     * @generated from protobuf field: int32 crypt_fever = 82;
     */
    cryptFever: number;
    /**
     * @generated from protobuf field: bool bone_shield = 83;
     */
    boneShield: boolean;
    /**
     * @generated from protobuf field: int32 wandering_plague = 84;
     */
    wanderingPlague: number;
    /**
     * @generated from protobuf field: int32 ebon_plaguebringer = 85;
     */
    ebonPlaguebringer: number;
    /**
     * @generated from protobuf field: bool scourge_strike = 86;
     */
    scourgeStrike: boolean;
    /**
     * @generated from protobuf field: int32 rage_of_rivendare = 87;
     */
    rageOfRivendare: number;
    /**
     * @generated from protobuf field: bool summon_gargoyle = 88;
     */
    summonGargoyle: boolean;
}
/**
 * @generated from protobuf message proto.Deathknight
 */
export interface Deathknight {
    /**
     * @generated from protobuf field: proto.Deathknight.Rotation rotation = 1;
     */
    rotation?: Deathknight_Rotation;
    /**
     * @generated from protobuf field: proto.DeathknightTalents talents = 2;
     */
    talents?: DeathknightTalents;
    /**
     * @generated from protobuf field: proto.Deathknight.Options options = 3;
     */
    options?: Deathknight_Options;
}
/**
 * @generated from protobuf message proto.Deathknight.Rotation
 */
export interface Deathknight_Rotation {
    /**
     * @generated from protobuf field: proto.Deathknight.Rotation.ArmyOfTheDead army_of_the_dead = 1;
     */
    armyOfTheDead: Deathknight_Rotation_ArmyOfTheDead;
    /**
     * @generated from protobuf field: bool use_death_and_decay = 2;
     */
    useDeathAndDecay: boolean;
    /**
     * @generated from protobuf field: bool unholy_presence_opener = 3;
     */
    unholyPresenceOpener: boolean;
    /**
     * @generated from protobuf field: double disease_refresh_duration = 4;
     */
    diseaseRefreshDuration: number;
    /**
     * @generated from protobuf field: bool refresh_horn_of_winter = 5;
     */
    refreshHornOfWinter: boolean;
}
/**
 * @generated from protobuf enum proto.Deathknight.Rotation.ArmyOfTheDead
 */
export declare enum Deathknight_Rotation_ArmyOfTheDead {
    /**
     * @generated from protobuf enum value: DoNotUse = 0;
     */
    DoNotUse = 0,
    /**
     * @generated from protobuf enum value: PreCast = 1;
     */
    PreCast = 1,
    /**
     * @generated from protobuf enum value: AsMajorCd = 2;
     */
    AsMajorCd = 2
}
/**
 * @generated from protobuf message proto.Deathknight.Options
 */
export interface Deathknight_Options {
    /**
     * @generated from protobuf field: double starting_runic_power = 1;
     */
    startingRunicPower: number;
    /**
     * @generated from protobuf field: double pet_uptime = 2;
     */
    petUptime: number;
    /**
     * @generated from protobuf field: bool precast_ghoul_frenzy = 3;
     */
    precastGhoulFrenzy: boolean;
    /**
     * @generated from protobuf field: bool precast_horn_of_winter = 4;
     */
    precastHornOfWinter: boolean;
}
/**
 * @generated from protobuf message proto.TankDeathknight
 */
export interface TankDeathknight {
    /**
     * @generated from protobuf field: proto.TankDeathknight.Rotation rotation = 1;
     */
    rotation?: TankDeathknight_Rotation;
    /**
     * @generated from protobuf field: proto.DeathknightTalents talents = 2;
     */
    talents?: DeathknightTalents;
    /**
     * @generated from protobuf field: proto.TankDeathknight.Options options = 3;
     */
    options?: TankDeathknight_Options;
}
/**
 * @generated from protobuf message proto.TankDeathknight.Rotation
 */
export interface TankDeathknight_Rotation {
}
/**
 * @generated from protobuf message proto.TankDeathknight.Options
 */
export interface TankDeathknight_Options {
    /**
     * @generated from protobuf field: double starting_runic_power = 1;
     */
    startingRunicPower: number;
}
/**
 * @generated from protobuf enum proto.DeathknightMajorGlyph
 */
export declare enum DeathknightMajorGlyph {
    /**
     * @generated from protobuf enum value: DeathknightMajorGlyphNone = 0;
     */
    DeathknightMajorGlyphNone = 0,
    /**
     * @generated from protobuf enum value: GlyphOfAntiMagicShell = 43533;
     */
    GlyphOfAntiMagicShell = 43533,
    /**
     * @generated from protobuf enum value: GlyphOfBloodStrike = 43826;
     */
    GlyphOfBloodStrike = 43826,
    /**
     * @generated from protobuf enum value: GlyphOfBoneShield = 43536;
     */
    GlyphOfBoneShield = 43536,
    /**
     * @generated from protobuf enum value: GlyphOfChainsOfIce = 43537;
     */
    GlyphOfChainsOfIce = 43537,
    /**
     * @generated from protobuf enum value: GlyphOfDancingRuneWeapon = 45799;
     */
    GlyphOfDancingRuneWeapon = 45799,
    /**
     * @generated from protobuf enum value: GlyphOfDarkCommand = 43538;
     */
    GlyphOfDarkCommand = 43538,
    /**
     * @generated from protobuf enum value: GlyphOfDarkDeath = 45804;
     */
    GlyphOfDarkDeath = 45804,
    /**
     * @generated from protobuf enum value: GlyphOfDeathAndDecay = 43542;
     */
    GlyphOfDeathAndDecay = 43542,
    /**
     * @generated from protobuf enum value: GlyphOfDeathGrip = 43541;
     */
    GlyphOfDeathGrip = 43541,
    /**
     * @generated from protobuf enum value: GlyphOfDeathStrike = 43827;
     */
    GlyphOfDeathStrike = 43827,
    /**
     * @generated from protobuf enum value: GlyphOfDisease = 45805;
     */
    GlyphOfDisease = 45805,
    /**
     * @generated from protobuf enum value: GlyphOfFrostStrike = 43543;
     */
    GlyphOfFrostStrike = 43543,
    /**
     * @generated from protobuf enum value: GlyphOfHeartStrike = 43534;
     */
    GlyphOfHeartStrike = 43534,
    /**
     * @generated from protobuf enum value: GlyphOfHowlingBlast = 45806;
     */
    GlyphOfHowlingBlast = 45806,
    /**
     * @generated from protobuf enum value: GlyphOfHungeringCold = 45800;
     */
    GlyphOfHungeringCold = 45800,
    /**
     * @generated from protobuf enum value: GlyphOfIceboundFortitude = 43545;
     */
    GlyphOfIceboundFortitude = 43545,
    /**
     * @generated from protobuf enum value: GlyphOfIcyTouch = 43546;
     */
    GlyphOfIcyTouch = 43546,
    /**
     * @generated from protobuf enum value: GlyphOfObliterate = 43547;
     */
    GlyphOfObliterate = 43547,
    /**
     * @generated from protobuf enum value: GlyphOfPlagueStrike = 43548;
     */
    GlyphOfPlagueStrike = 43548,
    /**
     * @generated from protobuf enum value: GlyphOfRuneStrike = 43550;
     */
    GlyphOfRuneStrike = 43550,
    /**
     * @generated from protobuf enum value: GlyphOfRuneTap = 43825;
     */
    GlyphOfRuneTap = 43825,
    /**
     * @generated from protobuf enum value: GlyphOfScourgeStrike = 43551;
     */
    GlyphOfScourgeStrike = 43551,
    /**
     * @generated from protobuf enum value: GlyphOfStrangulate = 43552;
     */
    GlyphOfStrangulate = 43552,
    /**
     * @generated from protobuf enum value: GlyphOfTheGhoul = 43549;
     */
    GlyphOfTheGhoul = 43549,
    /**
     * @generated from protobuf enum value: GlyphOfUnbreakableArmor = 43553;
     */
    GlyphOfUnbreakableArmor = 43553,
    /**
     * @generated from protobuf enum value: GlyphOfUnholyBlight = 45803;
     */
    GlyphOfUnholyBlight = 45803,
    /**
     * @generated from protobuf enum value: GlyphOfVampiricBlood = 43554;
     */
    GlyphOfVampiricBlood = 43554
}
/**
 * @generated from protobuf enum proto.DeathknightMinorGlyph
 */
export declare enum DeathknightMinorGlyph {
    /**
     * @generated from protobuf enum value: DeathknightMinorGlyphNone = 0;
     */
    DeathknightMinorGlyphNone = 0,
    /**
     * @generated from protobuf enum value: GlyphOfBloodTap = 43535;
     */
    GlyphOfBloodTap = 43535,
    /**
     * @generated from protobuf enum value: GlyphOfCorpseExplosion = 43671;
     */
    GlyphOfCorpseExplosion = 43671,
    /**
     * @generated from protobuf enum value: GlyphOfDeathSEmbrace = 43539;
     */
    GlyphOfDeathSEmbrace = 43539,
    /**
     * @generated from protobuf enum value: GlyphOfHornOfWinter = 43544;
     */
    GlyphOfHornOfWinter = 43544,
    /**
     * @generated from protobuf enum value: GlyphOfPestilence = 43672;
     */
    GlyphOfPestilence = 43672,
    /**
     * @generated from protobuf enum value: GlyphOfRaiseDead = 43673;
     */
    GlyphOfRaiseDead = 43673
}
declare class DeathknightTalents$Type extends MessageType<DeathknightTalents> {
    constructor();
    create(value?: PartialMessage<DeathknightTalents>): DeathknightTalents;
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: DeathknightTalents): DeathknightTalents;
    internalBinaryWrite(message: DeathknightTalents, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter;
}
/**
 * @generated MessageType for protobuf message proto.DeathknightTalents
 */
export declare const DeathknightTalents: DeathknightTalents$Type;
declare class Deathknight$Type extends MessageType<Deathknight> {
    constructor();
    create(value?: PartialMessage<Deathknight>): Deathknight;
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: Deathknight): Deathknight;
    internalBinaryWrite(message: Deathknight, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter;
}
/**
 * @generated MessageType for protobuf message proto.Deathknight
 */
export declare const Deathknight: Deathknight$Type;
declare class Deathknight_Rotation$Type extends MessageType<Deathknight_Rotation> {
    constructor();
    create(value?: PartialMessage<Deathknight_Rotation>): Deathknight_Rotation;
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: Deathknight_Rotation): Deathknight_Rotation;
    internalBinaryWrite(message: Deathknight_Rotation, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter;
}
/**
 * @generated MessageType for protobuf message proto.Deathknight.Rotation
 */
export declare const Deathknight_Rotation: Deathknight_Rotation$Type;
declare class Deathknight_Options$Type extends MessageType<Deathknight_Options> {
    constructor();
    create(value?: PartialMessage<Deathknight_Options>): Deathknight_Options;
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: Deathknight_Options): Deathknight_Options;
    internalBinaryWrite(message: Deathknight_Options, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter;
}
/**
 * @generated MessageType for protobuf message proto.Deathknight.Options
 */
export declare const Deathknight_Options: Deathknight_Options$Type;
declare class TankDeathknight$Type extends MessageType<TankDeathknight> {
    constructor();
    create(value?: PartialMessage<TankDeathknight>): TankDeathknight;
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: TankDeathknight): TankDeathknight;
    internalBinaryWrite(message: TankDeathknight, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter;
}
/**
 * @generated MessageType for protobuf message proto.TankDeathknight
 */
export declare const TankDeathknight: TankDeathknight$Type;
declare class TankDeathknight_Rotation$Type extends MessageType<TankDeathknight_Rotation> {
    constructor();
    create(value?: PartialMessage<TankDeathknight_Rotation>): TankDeathknight_Rotation;
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: TankDeathknight_Rotation): TankDeathknight_Rotation;
    internalBinaryWrite(message: TankDeathknight_Rotation, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter;
}
/**
 * @generated MessageType for protobuf message proto.TankDeathknight.Rotation
 */
export declare const TankDeathknight_Rotation: TankDeathknight_Rotation$Type;
declare class TankDeathknight_Options$Type extends MessageType<TankDeathknight_Options> {
    constructor();
    create(value?: PartialMessage<TankDeathknight_Options>): TankDeathknight_Options;
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: TankDeathknight_Options): TankDeathknight_Options;
    internalBinaryWrite(message: TankDeathknight_Options, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter;
}
/**
 * @generated MessageType for protobuf message proto.TankDeathknight.Options
 */
export declare const TankDeathknight_Options: TankDeathknight_Options$Type;
export {};
