import type { BinaryWriteOptions } from "@protobuf-ts/runtime";
import type { IBinaryWriter } from "@protobuf-ts/runtime";
import type { BinaryReadOptions } from "@protobuf-ts/runtime";
import type { IBinaryReader } from "@protobuf-ts/runtime";
import type { PartialMessage } from "@protobuf-ts/runtime";
import { MessageType } from "@protobuf-ts/runtime";
/**
 * @generated from protobuf message proto.ShamanTalents
 */
export interface ShamanTalents {
    /**
     * Elemental
     *
     * @generated from protobuf field: int32 convection = 1;
     */
    convection: number;
    /**
     * @generated from protobuf field: int32 concussion = 2;
     */
    concussion: number;
    /**
     * @generated from protobuf field: int32 call_of_flame = 3;
     */
    callOfFlame: number;
    /**
     * @generated from protobuf field: int32 elemental_warding = 4;
     */
    elementalWarding: number;
    /**
     * @generated from protobuf field: int32 elemental_devastation = 5;
     */
    elementalDevastation: number;
    /**
     * @generated from protobuf field: int32 reverberation = 6;
     */
    reverberation: number;
    /**
     * @generated from protobuf field: bool elemental_focus = 7;
     */
    elementalFocus: boolean;
    /**
     * @generated from protobuf field: int32 elemental_fury = 8;
     */
    elementalFury: number;
    /**
     * @generated from protobuf field: int32 improved_fire_nova = 9;
     */
    improvedFireNova: number;
    /**
     * @generated from protobuf field: int32 eye_of_the_storm = 10;
     */
    eyeOfTheStorm: number;
    /**
     * @generated from protobuf field: int32 elemental_reach = 11;
     */
    elementalReach: number;
    /**
     * @generated from protobuf field: bool call_of_thunder = 12;
     */
    callOfThunder: boolean;
    /**
     * @generated from protobuf field: int32 unrelenting_storm = 13;
     */
    unrelentingStorm: number;
    /**
     * @generated from protobuf field: int32 elemental_precision = 14;
     */
    elementalPrecision: number;
    /**
     * @generated from protobuf field: int32 lightning_mastery = 15;
     */
    lightningMastery: number;
    /**
     * @generated from protobuf field: bool elemental_mastery = 16;
     */
    elementalMastery: boolean;
    /**
     * @generated from protobuf field: int32 storm_earth_and_fire = 17;
     */
    stormEarthAndFire: number;
    /**
     * @generated from protobuf field: int32 booming_echoes = 18;
     */
    boomingEchoes: number;
    /**
     * @generated from protobuf field: int32 elemental_oath = 19;
     */
    elementalOath: number;
    /**
     * @generated from protobuf field: int32 lightning_overload = 20;
     */
    lightningOverload: number;
    /**
     * @generated from protobuf field: int32 astral_shift = 21;
     */
    astralShift: number;
    /**
     * @generated from protobuf field: bool totem_of_wrath = 22;
     */
    totemOfWrath: boolean;
    /**
     * @generated from protobuf field: int32 lava_flows = 23;
     */
    lavaFlows: number;
    /**
     * @generated from protobuf field: int32 shamanism = 24;
     */
    shamanism: number;
    /**
     * @generated from protobuf field: bool thunderstorm = 25;
     */
    thunderstorm: boolean;
    /**
     * Enhancement
     *
     * @generated from protobuf field: int32 enhancing_totems = 26;
     */
    enhancingTotems: number;
    /**
     * @generated from protobuf field: int32 earths_grasp = 27;
     */
    earthsGrasp: number;
    /**
     * @generated from protobuf field: int32 ancestral_knowledge = 28;
     */
    ancestralKnowledge: number;
    /**
     * @generated from protobuf field: int32 guardian_totems = 29;
     */
    guardianTotems: number;
    /**
     * @generated from protobuf field: int32 thundering_strikes = 30;
     */
    thunderingStrikes: number;
    /**
     * @generated from protobuf field: int32 improved_ghost_wolf = 31;
     */
    improvedGhostWolf: number;
    /**
     * @generated from protobuf field: int32 improved_shields = 32;
     */
    improvedShields: number;
    /**
     * @generated from protobuf field: int32 elemental_weapons = 33;
     */
    elementalWeapons: number;
    /**
     * @generated from protobuf field: bool shamanistic_focus = 34;
     */
    shamanisticFocus: boolean;
    /**
     * @generated from protobuf field: int32 anticipation = 35;
     */
    anticipation: number;
    /**
     * @generated from protobuf field: int32 flurry = 36;
     */
    flurry: number;
    /**
     * @generated from protobuf field: int32 toughness = 37;
     */
    toughness: number;
    /**
     * @generated from protobuf field: int32 improved_windfury_totem = 38;
     */
    improvedWindfuryTotem: number;
    /**
     * @generated from protobuf field: bool spirit_weapons = 39;
     */
    spiritWeapons: boolean;
    /**
     * @generated from protobuf field: int32 mental_dexterity = 40;
     */
    mentalDexterity: number;
    /**
     * @generated from protobuf field: int32 unleashed_rage = 41;
     */
    unleashedRage: number;
    /**
     * @generated from protobuf field: int32 weapon_mastery = 42;
     */
    weaponMastery: number;
    /**
     * @generated from protobuf field: int32 frozen_power = 43;
     */
    frozenPower: number;
    /**
     * @generated from protobuf field: int32 dual_wield_specialization = 44;
     */
    dualWieldSpecialization: number;
    /**
     * @generated from protobuf field: bool dual_wield = 45;
     */
    dualWield: boolean;
    /**
     * @generated from protobuf field: bool stormstrike = 46;
     */
    stormstrike: boolean;
    /**
     * @generated from protobuf field: int32 static_shock = 47;
     */
    staticShock: number;
    /**
     * @generated from protobuf field: bool lava_lash = 48;
     */
    lavaLash: boolean;
    /**
     * @generated from protobuf field: int32 improved_stormstrike = 49;
     */
    improvedStormstrike: number;
    /**
     * @generated from protobuf field: int32 mental_quickness = 50;
     */
    mentalQuickness: number;
    /**
     * @generated from protobuf field: bool shamanistic_rage = 51;
     */
    shamanisticRage: boolean;
    /**
     * @generated from protobuf field: int32 earthen_power = 52;
     */
    earthenPower: number;
    /**
     * @generated from protobuf field: int32 maelstrom_weapon = 53;
     */
    maelstromWeapon: number;
    /**
     * @generated from protobuf field: bool feral_spirit = 54;
     */
    feralSpirit: boolean;
    /**
     * Restoration
     *
     * @generated from protobuf field: int32 improved_healing_wave = 55;
     */
    improvedHealingWave: number;
    /**
     * @generated from protobuf field: int32 totemic_focus = 56;
     */
    totemicFocus: number;
    /**
     * @generated from protobuf field: int32 improved_reincarnation = 57;
     */
    improvedReincarnation: number;
    /**
     * @generated from protobuf field: int32 healing_grace = 58;
     */
    healingGrace: number;
    /**
     * @generated from protobuf field: int32 tidal_focus = 59;
     */
    tidalFocus: number;
    /**
     * @generated from protobuf field: int32 improved_water_shield = 60;
     */
    improvedWaterShield: number;
    /**
     * @generated from protobuf field: int32 healing_focus = 61;
     */
    healingFocus: number;
    /**
     * @generated from protobuf field: bool tidal_force = 62;
     */
    tidalForce: boolean;
    /**
     * @generated from protobuf field: int32 ancestral_healing = 63;
     */
    ancestralHealing: number;
    /**
     * @generated from protobuf field: int32 restorative_totems = 64;
     */
    restorativeTotems: number;
    /**
     * @generated from protobuf field: int32 tidal_mastery = 65;
     */
    tidalMastery: number;
    /**
     * @generated from protobuf field: int32 healing_way = 66;
     */
    healingWay: number;
    /**
     * @generated from protobuf field: bool natures_swiftness = 67;
     */
    naturesSwiftness: boolean;
    /**
     * @generated from protobuf field: int32 focused_mind = 68;
     */
    focusedMind: number;
    /**
     * @generated from protobuf field: int32 purification = 69;
     */
    purification: number;
    /**
     * @generated from protobuf field: int32 natures_guardian = 70;
     */
    naturesGuardian: number;
    /**
     * @generated from protobuf field: bool mana_tide_totem = 71;
     */
    manaTideTotem: boolean;
    /**
     * @generated from protobuf field: bool cleanse_spirit = 72;
     */
    cleanseSpirit: boolean;
    /**
     * @generated from protobuf field: int32 blessing_of_the_eternals = 73;
     */
    blessingOfTheEternals: number;
    /**
     * @generated from protobuf field: int32 improved_chain_heal = 74;
     */
    improvedChainHeal: number;
    /**
     * @generated from protobuf field: int32 natures_blessing = 75;
     */
    naturesBlessing: number;
    /**
     * @generated from protobuf field: int32 ancestral_awakening = 76;
     */
    ancestralAwakening: number;
    /**
     * @generated from protobuf field: bool earth_shield = 77;
     */
    earthShield: boolean;
    /**
     * @generated from protobuf field: int32 improved_earth_shield = 78;
     */
    improvedEarthShield: number;
    /**
     * @generated from protobuf field: int32 tidal_waves = 79;
     */
    tidalWaves: number;
    /**
     * @generated from protobuf field: bool riptide = 80;
     */
    riptide: boolean;
}
/**
 * @generated from protobuf message proto.ShamanTotems
 */
export interface ShamanTotems {
    /**
     * @generated from protobuf field: proto.EarthTotem earth = 1;
     */
    earth: EarthTotem;
    /**
     * @generated from protobuf field: proto.AirTotem air = 2;
     */
    air: AirTotem;
    /**
     * @generated from protobuf field: proto.FireTotem fire = 3;
     */
    fire: FireTotem;
    /**
     * @generated from protobuf field: proto.WaterTotem water = 4;
     */
    water: WaterTotem;
    /**
     * If set, will use mana tide when appropriate.
     *
     * @generated from protobuf field: bool use_mana_tide = 5;
     */
    useManaTide: boolean;
    /**
     * If set, will use fire elemental totem at the start and revert to regular
     * fire totems when it expires.
     *
     * @generated from protobuf field: bool use_fire_elemental = 6;
     */
    useFireElemental: boolean;
    /**
     * If set, will revert to regular fire totems when fire elemental goes OOM,
     * instead of waiting the full 2 minutes.
     *
     * @generated from protobuf field: bool recall_fire_elemental_on_oom = 7;
     */
    recallFireElementalOnOom: boolean;
    /**
     * If set, any time a 2-minute totem is about to expire, will recall and
     * replace all totems.
     *
     * @generated from protobuf field: bool recall_totems = 8;
     */
    recallTotems: boolean;
}
/**
 * @generated from protobuf message proto.ElementalShaman
 */
export interface ElementalShaman {
    /**
     * @generated from protobuf field: proto.ElementalShaman.Rotation rotation = 1;
     */
    rotation?: ElementalShaman_Rotation;
    /**
     * @generated from protobuf field: proto.ShamanTalents talents = 2;
     */
    talents?: ShamanTalents;
    /**
     * @generated from protobuf field: proto.ElementalShaman.Options options = 3;
     */
    options?: ElementalShaman_Options;
}
/**
 * @generated from protobuf message proto.ElementalShaman.Rotation
 */
export interface ElementalShaman_Rotation {
    /**
     * @generated from protobuf field: proto.ShamanTotems totems = 3;
     */
    totems?: ShamanTotems;
    /**
     * @generated from protobuf field: proto.ElementalShaman.Rotation.RotationType type = 1;
     */
    type: ElementalShaman_Rotation_RotationType;
    /**
     * @generated from protobuf field: bool in_thunderstorm_range = 2;
     */
    inThunderstormRange: boolean;
}
/**
 * @generated from protobuf enum proto.ElementalShaman.Rotation.RotationType
 */
export declare enum ElementalShaman_Rotation_RotationType {
    /**
     * @generated from protobuf enum value: Unknown = 0;
     */
    Unknown = 0,
    /**
     * @generated from protobuf enum value: Adaptive = 1;
     */
    Adaptive = 1
}
/**
 * @generated from protobuf message proto.ElementalShaman.Options
 */
export interface ElementalShaman_Options {
    /**
     * @generated from protobuf field: proto.ShamanShield shield = 1;
     */
    shield: ShamanShield;
    /**
     * @generated from protobuf field: bool bloodlust = 2;
     */
    bloodlust: boolean;
}
/**
 * @generated from protobuf message proto.EnhancementShaman
 */
export interface EnhancementShaman {
    /**
     * @generated from protobuf field: proto.EnhancementShaman.Rotation rotation = 1;
     */
    rotation?: EnhancementShaman_Rotation;
    /**
     * @generated from protobuf field: proto.ShamanTalents talents = 2;
     */
    talents?: ShamanTalents;
    /**
     * @generated from protobuf field: proto.EnhancementShaman.Options options = 3;
     */
    options?: EnhancementShaman_Options;
}
/**
 * @generated from protobuf message proto.EnhancementShaman.Rotation
 */
export interface EnhancementShaman_Rotation {
    /**
     * @generated from protobuf field: proto.ShamanTotems totems = 1;
     */
    totems?: ShamanTotems;
}
/**
 * @generated from protobuf message proto.EnhancementShaman.Options
 */
export interface EnhancementShaman_Options {
    /**
     * @generated from protobuf field: proto.ShamanShield shield = 1;
     */
    shield: ShamanShield;
    /**
     * @generated from protobuf field: bool bloodlust = 2;
     */
    bloodlust: boolean;
    /**
     * @generated from protobuf field: bool delay_offhand_swings = 3;
     */
    delayOffhandSwings: boolean;
    /**
     * @generated from protobuf field: proto.ShamanImbue imbueMH = 4;
     */
    imbueMH: ShamanImbue;
    /**
     * @generated from protobuf field: proto.ShamanImbue imbueOH = 5;
     */
    imbueOH: ShamanImbue;
}
/**
 * @generated from protobuf enum proto.ShamanMajorGlyph
 */
export declare enum ShamanMajorGlyph {
    /**
     * @generated from protobuf enum value: ShamanMajorGlyphNone = 0;
     */
    ShamanMajorGlyphNone = 0,
    /**
     * @generated from protobuf enum value: GlyphOfChainHeal = 41517;
     */
    GlyphOfChainHeal = 41517,
    /**
     * @generated from protobuf enum value: GlyphOfChainLightning = 41518;
     */
    GlyphOfChainLightning = 41518,
    /**
     * @generated from protobuf enum value: GlyphOfEarthShield = 45775;
     */
    GlyphOfEarthShield = 45775,
    /**
     * @generated from protobuf enum value: GlyphOfEarthlivingWeapon = 41527;
     */
    GlyphOfEarthlivingWeapon = 41527,
    /**
     * @generated from protobuf enum value: GlyphOfElementalMastery = 41552;
     */
    GlyphOfElementalMastery = 41552,
    /**
     * @generated from protobuf enum value: GlyphOfFeralSpirit = 45771;
     */
    GlyphOfFeralSpirit = 45771,
    /**
     * @generated from protobuf enum value: GlyphOfFireElementalTotem = 41529;
     */
    GlyphOfFireElementalTotem = 41529,
    /**
     * @generated from protobuf enum value: GlyphOfFireNova = 41530;
     */
    GlyphOfFireNova = 41530,
    /**
     * @generated from protobuf enum value: GlyphOfFlameShock = 41531;
     */
    GlyphOfFlameShock = 41531,
    /**
     * @generated from protobuf enum value: GlyphOfFlametongueWeapon = 41532;
     */
    GlyphOfFlametongueWeapon = 41532,
    /**
     * @generated from protobuf enum value: GlyphOfFrostShock = 41547;
     */
    GlyphOfFrostShock = 41547,
    /**
     * @generated from protobuf enum value: GlyphOfHealingStreamTotem = 41533;
     */
    GlyphOfHealingStreamTotem = 41533,
    /**
     * @generated from protobuf enum value: GlyphOfHealingWave = 41534;
     */
    GlyphOfHealingWave = 41534,
    /**
     * @generated from protobuf enum value: GlyphOfHex = 45777;
     */
    GlyphOfHex = 45777,
    /**
     * @generated from protobuf enum value: GlyphOfLava = 41524;
     */
    GlyphOfLava = 41524,
    /**
     * @generated from protobuf enum value: GlyphOfLavaLash = 41540;
     */
    GlyphOfLavaLash = 41540,
    /**
     * @generated from protobuf enum value: GlyphOfLesserHealingWave = 41535;
     */
    GlyphOfLesserHealingWave = 41535,
    /**
     * @generated from protobuf enum value: GlyphOfLightningBolt = 41536;
     */
    GlyphOfLightningBolt = 41536,
    /**
     * @generated from protobuf enum value: GlyphOfLightningShield = 41537;
     */
    GlyphOfLightningShield = 41537,
    /**
     * @generated from protobuf enum value: GlyphOfManaTide = 41538;
     */
    GlyphOfManaTide = 41538,
    /**
     * @generated from protobuf enum value: GlyphOfRiptide = 45772;
     */
    GlyphOfRiptide = 45772,
    /**
     * @generated from protobuf enum value: GlyphOfShocking = 41526;
     */
    GlyphOfShocking = 41526,
    /**
     * @generated from protobuf enum value: GlyphOfStoneclawTotem = 45778;
     */
    GlyphOfStoneclawTotem = 45778,
    /**
     * @generated from protobuf enum value: GlyphOfStormstrike = 41539;
     */
    GlyphOfStormstrike = 41539,
    /**
     * @generated from protobuf enum value: GlyphOfThunder = 45770;
     */
    GlyphOfThunder = 45770,
    /**
     * @generated from protobuf enum value: GlyphOfTotemOfWrath = 45776;
     */
    GlyphOfTotemOfWrath = 45776,
    /**
     * @generated from protobuf enum value: GlyphOfWaterMastery = 41541;
     */
    GlyphOfWaterMastery = 41541,
    /**
     * @generated from protobuf enum value: GlyphOfWindfuryWeapon = 41542;
     */
    GlyphOfWindfuryWeapon = 41542
}
/**
 * @generated from protobuf enum proto.ShamanMinorGlyph
 */
export declare enum ShamanMinorGlyph {
    /**
     * @generated from protobuf enum value: ShamanMinorGlyphNone = 0;
     */
    ShamanMinorGlyphNone = 0,
    /**
     * @generated from protobuf enum value: GlyphOfAstralRecall = 43381;
     */
    GlyphOfAstralRecall = 43381,
    /**
     * @generated from protobuf enum value: GlyphOfGhostWolf = 43725;
     */
    GlyphOfGhostWolf = 43725,
    /**
     * @generated from protobuf enum value: GlyphOfRenewedLife = 43385;
     */
    GlyphOfRenewedLife = 43385,
    /**
     * @generated from protobuf enum value: GlyphOfThunderstorm = 44923;
     */
    GlyphOfThunderstorm = 44923,
    /**
     * @generated from protobuf enum value: GlyphOfWaterBreathing = 43344;
     */
    GlyphOfWaterBreathing = 43344,
    /**
     * @generated from protobuf enum value: GlyphOfWaterShield = 43386;
     */
    GlyphOfWaterShield = 43386,
    /**
     * @generated from protobuf enum value: GlyphOfWaterWalking = 43388;
     */
    GlyphOfWaterWalking = 43388
}
/**
 * @generated from protobuf enum proto.EarthTotem
 */
export declare enum EarthTotem {
    /**
     * @generated from protobuf enum value: NoEarthTotem = 0;
     */
    NoEarthTotem = 0,
    /**
     * @generated from protobuf enum value: StrengthOfEarthTotem = 1;
     */
    StrengthOfEarthTotem = 1,
    /**
     * @generated from protobuf enum value: TremorTotem = 2;
     */
    TremorTotem = 2
}
/**
 * @generated from protobuf enum proto.AirTotem
 */
export declare enum AirTotem {
    /**
     * @generated from protobuf enum value: NoAirTotem = 0;
     */
    NoAirTotem = 0,
    /**
     * @generated from protobuf enum value: TranquilAirTotem = 1;
     */
    TranquilAirTotem = 1,
    /**
     * @generated from protobuf enum value: WindfuryTotem = 2;
     */
    WindfuryTotem = 2,
    /**
     * @generated from protobuf enum value: WrathOfAirTotem = 3;
     */
    WrathOfAirTotem = 3
}
/**
 * @generated from protobuf enum proto.FireTotem
 */
export declare enum FireTotem {
    /**
     * @generated from protobuf enum value: NoFireTotem = 0;
     */
    NoFireTotem = 0,
    /**
     * @generated from protobuf enum value: MagmaTotem = 1;
     */
    MagmaTotem = 1,
    /**
     * @generated from protobuf enum value: SearingTotem = 2;
     */
    SearingTotem = 2,
    /**
     * @generated from protobuf enum value: TotemOfWrath = 3;
     */
    TotemOfWrath = 3,
    /**
     * @generated from protobuf enum value: FlametongueTotem = 4;
     */
    FlametongueTotem = 4
}
/**
 * @generated from protobuf enum proto.WaterTotem
 */
export declare enum WaterTotem {
    /**
     * @generated from protobuf enum value: NoWaterTotem = 0;
     */
    NoWaterTotem = 0,
    /**
     * @generated from protobuf enum value: ManaSpringTotem = 1;
     */
    ManaSpringTotem = 1
}
/**
 * @generated from protobuf enum proto.ShamanShield
 */
export declare enum ShamanShield {
    /**
     * @generated from protobuf enum value: NoShield = 0;
     */
    NoShield = 0,
    /**
     * @generated from protobuf enum value: WaterShield = 1;
     */
    WaterShield = 1,
    /**
     * @generated from protobuf enum value: LightningShield = 2;
     */
    LightningShield = 2
}
/**
 * @generated from protobuf enum proto.ShamanImbue
 */
export declare enum ShamanImbue {
    /**
     * @generated from protobuf enum value: NoImbue = 0;
     */
    NoImbue = 0,
    /**
     * @generated from protobuf enum value: WindfuryWeapon = 1;
     */
    WindfuryWeapon = 1,
    /**
     * @generated from protobuf enum value: FlametongueWeapon = 2;
     */
    FlametongueWeapon = 2,
    /**
     * @generated from protobuf enum value: FrostbrandWeapon = 3;
     */
    FrostbrandWeapon = 3
}
declare class ShamanTalents$Type extends MessageType<ShamanTalents> {
    constructor();
    create(value?: PartialMessage<ShamanTalents>): ShamanTalents;
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: ShamanTalents): ShamanTalents;
    internalBinaryWrite(message: ShamanTalents, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter;
}
/**
 * @generated MessageType for protobuf message proto.ShamanTalents
 */
export declare const ShamanTalents: ShamanTalents$Type;
declare class ShamanTotems$Type extends MessageType<ShamanTotems> {
    constructor();
    create(value?: PartialMessage<ShamanTotems>): ShamanTotems;
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: ShamanTotems): ShamanTotems;
    internalBinaryWrite(message: ShamanTotems, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter;
}
/**
 * @generated MessageType for protobuf message proto.ShamanTotems
 */
export declare const ShamanTotems: ShamanTotems$Type;
declare class ElementalShaman$Type extends MessageType<ElementalShaman> {
    constructor();
    create(value?: PartialMessage<ElementalShaman>): ElementalShaman;
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: ElementalShaman): ElementalShaman;
    internalBinaryWrite(message: ElementalShaman, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter;
}
/**
 * @generated MessageType for protobuf message proto.ElementalShaman
 */
export declare const ElementalShaman: ElementalShaman$Type;
declare class ElementalShaman_Rotation$Type extends MessageType<ElementalShaman_Rotation> {
    constructor();
    create(value?: PartialMessage<ElementalShaman_Rotation>): ElementalShaman_Rotation;
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: ElementalShaman_Rotation): ElementalShaman_Rotation;
    internalBinaryWrite(message: ElementalShaman_Rotation, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter;
}
/**
 * @generated MessageType for protobuf message proto.ElementalShaman.Rotation
 */
export declare const ElementalShaman_Rotation: ElementalShaman_Rotation$Type;
declare class ElementalShaman_Options$Type extends MessageType<ElementalShaman_Options> {
    constructor();
    create(value?: PartialMessage<ElementalShaman_Options>): ElementalShaman_Options;
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: ElementalShaman_Options): ElementalShaman_Options;
    internalBinaryWrite(message: ElementalShaman_Options, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter;
}
/**
 * @generated MessageType for protobuf message proto.ElementalShaman.Options
 */
export declare const ElementalShaman_Options: ElementalShaman_Options$Type;
declare class EnhancementShaman$Type extends MessageType<EnhancementShaman> {
    constructor();
    create(value?: PartialMessage<EnhancementShaman>): EnhancementShaman;
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: EnhancementShaman): EnhancementShaman;
    internalBinaryWrite(message: EnhancementShaman, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter;
}
/**
 * @generated MessageType for protobuf message proto.EnhancementShaman
 */
export declare const EnhancementShaman: EnhancementShaman$Type;
declare class EnhancementShaman_Rotation$Type extends MessageType<EnhancementShaman_Rotation> {
    constructor();
    create(value?: PartialMessage<EnhancementShaman_Rotation>): EnhancementShaman_Rotation;
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: EnhancementShaman_Rotation): EnhancementShaman_Rotation;
    internalBinaryWrite(message: EnhancementShaman_Rotation, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter;
}
/**
 * @generated MessageType for protobuf message proto.EnhancementShaman.Rotation
 */
export declare const EnhancementShaman_Rotation: EnhancementShaman_Rotation$Type;
declare class EnhancementShaman_Options$Type extends MessageType<EnhancementShaman_Options> {
    constructor();
    create(value?: PartialMessage<EnhancementShaman_Options>): EnhancementShaman_Options;
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: EnhancementShaman_Options): EnhancementShaman_Options;
    internalBinaryWrite(message: EnhancementShaman_Options, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter;
}
/**
 * @generated MessageType for protobuf message proto.EnhancementShaman.Options
 */
export declare const EnhancementShaman_Options: EnhancementShaman_Options$Type;
export {};
