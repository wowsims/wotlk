import { WireType } from '/wotlk/protobuf-ts/index.js';
import { UnknownFieldHandler } from '/wotlk/protobuf-ts/index.js';
import { reflectionMergePartial } from '/wotlk/protobuf-ts/index.js';
import { MESSAGE_TYPE } from '/wotlk/protobuf-ts/index.js';
import { MessageType } from '/wotlk/protobuf-ts/index.js';
/**
 * @generated from protobuf enum proto.ElementalShaman.Rotation.RotationType
 */
export var ElementalShaman_Rotation_RotationType;
(function (ElementalShaman_Rotation_RotationType) {
    /**
     * @generated from protobuf enum value: Unknown = 0;
     */
    ElementalShaman_Rotation_RotationType[ElementalShaman_Rotation_RotationType["Unknown"] = 0] = "Unknown";
    /**
     * @generated from protobuf enum value: Adaptive = 1;
     */
    ElementalShaman_Rotation_RotationType[ElementalShaman_Rotation_RotationType["Adaptive"] = 1] = "Adaptive";
})(ElementalShaman_Rotation_RotationType || (ElementalShaman_Rotation_RotationType = {}));
/**
 * @generated from protobuf enum proto.ShamanMajorGlyph
 */
export var ShamanMajorGlyph;
(function (ShamanMajorGlyph) {
    /**
     * @generated from protobuf enum value: ShamanMajorGlyphNone = 0;
     */
    ShamanMajorGlyph[ShamanMajorGlyph["ShamanMajorGlyphNone"] = 0] = "ShamanMajorGlyphNone";
    /**
     * @generated from protobuf enum value: GlyphOfChainHeal = 41517;
     */
    ShamanMajorGlyph[ShamanMajorGlyph["GlyphOfChainHeal"] = 41517] = "GlyphOfChainHeal";
    /**
     * @generated from protobuf enum value: GlyphOfChainLightning = 41518;
     */
    ShamanMajorGlyph[ShamanMajorGlyph["GlyphOfChainLightning"] = 41518] = "GlyphOfChainLightning";
    /**
     * @generated from protobuf enum value: GlyphOfEarthShield = 45775;
     */
    ShamanMajorGlyph[ShamanMajorGlyph["GlyphOfEarthShield"] = 45775] = "GlyphOfEarthShield";
    /**
     * @generated from protobuf enum value: GlyphOfEarthlivingWeapon = 41527;
     */
    ShamanMajorGlyph[ShamanMajorGlyph["GlyphOfEarthlivingWeapon"] = 41527] = "GlyphOfEarthlivingWeapon";
    /**
     * @generated from protobuf enum value: GlyphOfElementalMastery = 41552;
     */
    ShamanMajorGlyph[ShamanMajorGlyph["GlyphOfElementalMastery"] = 41552] = "GlyphOfElementalMastery";
    /**
     * @generated from protobuf enum value: GlyphOfFeralSpirit = 45771;
     */
    ShamanMajorGlyph[ShamanMajorGlyph["GlyphOfFeralSpirit"] = 45771] = "GlyphOfFeralSpirit";
    /**
     * @generated from protobuf enum value: GlyphOfFireElementalTotem = 41529;
     */
    ShamanMajorGlyph[ShamanMajorGlyph["GlyphOfFireElementalTotem"] = 41529] = "GlyphOfFireElementalTotem";
    /**
     * @generated from protobuf enum value: GlyphOfFireNova = 41530;
     */
    ShamanMajorGlyph[ShamanMajorGlyph["GlyphOfFireNova"] = 41530] = "GlyphOfFireNova";
    /**
     * @generated from protobuf enum value: GlyphOfFlameShock = 41531;
     */
    ShamanMajorGlyph[ShamanMajorGlyph["GlyphOfFlameShock"] = 41531] = "GlyphOfFlameShock";
    /**
     * @generated from protobuf enum value: GlyphOfFlametongueWeapon = 41532;
     */
    ShamanMajorGlyph[ShamanMajorGlyph["GlyphOfFlametongueWeapon"] = 41532] = "GlyphOfFlametongueWeapon";
    /**
     * @generated from protobuf enum value: GlyphOfFrostShock = 41547;
     */
    ShamanMajorGlyph[ShamanMajorGlyph["GlyphOfFrostShock"] = 41547] = "GlyphOfFrostShock";
    /**
     * @generated from protobuf enum value: GlyphOfHealingStreamTotem = 41533;
     */
    ShamanMajorGlyph[ShamanMajorGlyph["GlyphOfHealingStreamTotem"] = 41533] = "GlyphOfHealingStreamTotem";
    /**
     * @generated from protobuf enum value: GlyphOfHealingWave = 41534;
     */
    ShamanMajorGlyph[ShamanMajorGlyph["GlyphOfHealingWave"] = 41534] = "GlyphOfHealingWave";
    /**
     * @generated from protobuf enum value: GlyphOfHex = 45777;
     */
    ShamanMajorGlyph[ShamanMajorGlyph["GlyphOfHex"] = 45777] = "GlyphOfHex";
    /**
     * @generated from protobuf enum value: GlyphOfLava = 41524;
     */
    ShamanMajorGlyph[ShamanMajorGlyph["GlyphOfLava"] = 41524] = "GlyphOfLava";
    /**
     * @generated from protobuf enum value: GlyphOfLavaLash = 41540;
     */
    ShamanMajorGlyph[ShamanMajorGlyph["GlyphOfLavaLash"] = 41540] = "GlyphOfLavaLash";
    /**
     * @generated from protobuf enum value: GlyphOfLesserHealingWave = 41535;
     */
    ShamanMajorGlyph[ShamanMajorGlyph["GlyphOfLesserHealingWave"] = 41535] = "GlyphOfLesserHealingWave";
    /**
     * @generated from protobuf enum value: GlyphOfLightningBolt = 41536;
     */
    ShamanMajorGlyph[ShamanMajorGlyph["GlyphOfLightningBolt"] = 41536] = "GlyphOfLightningBolt";
    /**
     * @generated from protobuf enum value: GlyphOfLightningShield = 41537;
     */
    ShamanMajorGlyph[ShamanMajorGlyph["GlyphOfLightningShield"] = 41537] = "GlyphOfLightningShield";
    /**
     * @generated from protobuf enum value: GlyphOfManaTide = 41538;
     */
    ShamanMajorGlyph[ShamanMajorGlyph["GlyphOfManaTide"] = 41538] = "GlyphOfManaTide";
    /**
     * @generated from protobuf enum value: GlyphOfRiptide = 45772;
     */
    ShamanMajorGlyph[ShamanMajorGlyph["GlyphOfRiptide"] = 45772] = "GlyphOfRiptide";
    /**
     * @generated from protobuf enum value: GlyphOfShocking = 41526;
     */
    ShamanMajorGlyph[ShamanMajorGlyph["GlyphOfShocking"] = 41526] = "GlyphOfShocking";
    /**
     * @generated from protobuf enum value: GlyphOfStoneclawTotem = 45778;
     */
    ShamanMajorGlyph[ShamanMajorGlyph["GlyphOfStoneclawTotem"] = 45778] = "GlyphOfStoneclawTotem";
    /**
     * @generated from protobuf enum value: GlyphOfStormstrike = 41539;
     */
    ShamanMajorGlyph[ShamanMajorGlyph["GlyphOfStormstrike"] = 41539] = "GlyphOfStormstrike";
    /**
     * @generated from protobuf enum value: GlyphOfThunder = 45770;
     */
    ShamanMajorGlyph[ShamanMajorGlyph["GlyphOfThunder"] = 45770] = "GlyphOfThunder";
    /**
     * @generated from protobuf enum value: GlyphOfTotemOfWrath = 45776;
     */
    ShamanMajorGlyph[ShamanMajorGlyph["GlyphOfTotemOfWrath"] = 45776] = "GlyphOfTotemOfWrath";
    /**
     * @generated from protobuf enum value: GlyphOfWaterMastery = 41541;
     */
    ShamanMajorGlyph[ShamanMajorGlyph["GlyphOfWaterMastery"] = 41541] = "GlyphOfWaterMastery";
    /**
     * @generated from protobuf enum value: GlyphOfWindfuryWeapon = 41542;
     */
    ShamanMajorGlyph[ShamanMajorGlyph["GlyphOfWindfuryWeapon"] = 41542] = "GlyphOfWindfuryWeapon";
})(ShamanMajorGlyph || (ShamanMajorGlyph = {}));
/**
 * @generated from protobuf enum proto.ShamanMinorGlyph
 */
export var ShamanMinorGlyph;
(function (ShamanMinorGlyph) {
    /**
     * @generated from protobuf enum value: ShamanMinorGlyphNone = 0;
     */
    ShamanMinorGlyph[ShamanMinorGlyph["ShamanMinorGlyphNone"] = 0] = "ShamanMinorGlyphNone";
    /**
     * @generated from protobuf enum value: GlyphOfAstralRecall = 43381;
     */
    ShamanMinorGlyph[ShamanMinorGlyph["GlyphOfAstralRecall"] = 43381] = "GlyphOfAstralRecall";
    /**
     * @generated from protobuf enum value: GlyphOfGhostWolf = 43725;
     */
    ShamanMinorGlyph[ShamanMinorGlyph["GlyphOfGhostWolf"] = 43725] = "GlyphOfGhostWolf";
    /**
     * @generated from protobuf enum value: GlyphOfRenewedLife = 43385;
     */
    ShamanMinorGlyph[ShamanMinorGlyph["GlyphOfRenewedLife"] = 43385] = "GlyphOfRenewedLife";
    /**
     * @generated from protobuf enum value: GlyphOfThunderstorm = 44923;
     */
    ShamanMinorGlyph[ShamanMinorGlyph["GlyphOfThunderstorm"] = 44923] = "GlyphOfThunderstorm";
    /**
     * @generated from protobuf enum value: GlyphOfWaterBreathing = 43344;
     */
    ShamanMinorGlyph[ShamanMinorGlyph["GlyphOfWaterBreathing"] = 43344] = "GlyphOfWaterBreathing";
    /**
     * @generated from protobuf enum value: GlyphOfWaterShield = 43386;
     */
    ShamanMinorGlyph[ShamanMinorGlyph["GlyphOfWaterShield"] = 43386] = "GlyphOfWaterShield";
    /**
     * @generated from protobuf enum value: GlyphOfWaterWalking = 43388;
     */
    ShamanMinorGlyph[ShamanMinorGlyph["GlyphOfWaterWalking"] = 43388] = "GlyphOfWaterWalking";
})(ShamanMinorGlyph || (ShamanMinorGlyph = {}));
/**
 * @generated from protobuf enum proto.EarthTotem
 */
export var EarthTotem;
(function (EarthTotem) {
    /**
     * @generated from protobuf enum value: NoEarthTotem = 0;
     */
    EarthTotem[EarthTotem["NoEarthTotem"] = 0] = "NoEarthTotem";
    /**
     * @generated from protobuf enum value: StrengthOfEarthTotem = 1;
     */
    EarthTotem[EarthTotem["StrengthOfEarthTotem"] = 1] = "StrengthOfEarthTotem";
    /**
     * @generated from protobuf enum value: TremorTotem = 2;
     */
    EarthTotem[EarthTotem["TremorTotem"] = 2] = "TremorTotem";
})(EarthTotem || (EarthTotem = {}));
/**
 * @generated from protobuf enum proto.AirTotem
 */
export var AirTotem;
(function (AirTotem) {
    /**
     * @generated from protobuf enum value: NoAirTotem = 0;
     */
    AirTotem[AirTotem["NoAirTotem"] = 0] = "NoAirTotem";
    /**
     * @generated from protobuf enum value: TranquilAirTotem = 1;
     */
    AirTotem[AirTotem["TranquilAirTotem"] = 1] = "TranquilAirTotem";
    /**
     * @generated from protobuf enum value: WindfuryTotem = 2;
     */
    AirTotem[AirTotem["WindfuryTotem"] = 2] = "WindfuryTotem";
    /**
     * @generated from protobuf enum value: WrathOfAirTotem = 3;
     */
    AirTotem[AirTotem["WrathOfAirTotem"] = 3] = "WrathOfAirTotem";
})(AirTotem || (AirTotem = {}));
/**
 * @generated from protobuf enum proto.FireTotem
 */
export var FireTotem;
(function (FireTotem) {
    /**
     * @generated from protobuf enum value: NoFireTotem = 0;
     */
    FireTotem[FireTotem["NoFireTotem"] = 0] = "NoFireTotem";
    /**
     * @generated from protobuf enum value: MagmaTotem = 1;
     */
    FireTotem[FireTotem["MagmaTotem"] = 1] = "MagmaTotem";
    /**
     * @generated from protobuf enum value: SearingTotem = 2;
     */
    FireTotem[FireTotem["SearingTotem"] = 2] = "SearingTotem";
    /**
     * @generated from protobuf enum value: TotemOfWrath = 3;
     */
    FireTotem[FireTotem["TotemOfWrath"] = 3] = "TotemOfWrath";
    /**
     * @generated from protobuf enum value: FlametongueTotem = 4;
     */
    FireTotem[FireTotem["FlametongueTotem"] = 4] = "FlametongueTotem";
})(FireTotem || (FireTotem = {}));
/**
 * @generated from protobuf enum proto.WaterTotem
 */
export var WaterTotem;
(function (WaterTotem) {
    /**
     * @generated from protobuf enum value: NoWaterTotem = 0;
     */
    WaterTotem[WaterTotem["NoWaterTotem"] = 0] = "NoWaterTotem";
    /**
     * @generated from protobuf enum value: ManaSpringTotem = 1;
     */
    WaterTotem[WaterTotem["ManaSpringTotem"] = 1] = "ManaSpringTotem";
})(WaterTotem || (WaterTotem = {}));
/**
 * @generated from protobuf enum proto.ShamanShield
 */
export var ShamanShield;
(function (ShamanShield) {
    /**
     * @generated from protobuf enum value: NoShield = 0;
     */
    ShamanShield[ShamanShield["NoShield"] = 0] = "NoShield";
    /**
     * @generated from protobuf enum value: WaterShield = 1;
     */
    ShamanShield[ShamanShield["WaterShield"] = 1] = "WaterShield";
    /**
     * @generated from protobuf enum value: LightningShield = 2;
     */
    ShamanShield[ShamanShield["LightningShield"] = 2] = "LightningShield";
})(ShamanShield || (ShamanShield = {}));
/**
 * @generated from protobuf enum proto.ShamanImbue
 */
export var ShamanImbue;
(function (ShamanImbue) {
    /**
     * @generated from protobuf enum value: NoImbue = 0;
     */
    ShamanImbue[ShamanImbue["NoImbue"] = 0] = "NoImbue";
    /**
     * @generated from protobuf enum value: WindfuryWeapon = 1;
     */
    ShamanImbue[ShamanImbue["WindfuryWeapon"] = 1] = "WindfuryWeapon";
    /**
     * @generated from protobuf enum value: FlametongueWeapon = 2;
     */
    ShamanImbue[ShamanImbue["FlametongueWeapon"] = 2] = "FlametongueWeapon";
    /**
     * @generated from protobuf enum value: FrostbrandWeapon = 3;
     */
    ShamanImbue[ShamanImbue["FrostbrandWeapon"] = 3] = "FrostbrandWeapon";
})(ShamanImbue || (ShamanImbue = {}));
// @generated message type with reflection information, may provide speed optimized methods
class ShamanTalents$Type extends MessageType {
    constructor() {
        super("proto.ShamanTalents", [
            { no: 1, name: "convection", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 2, name: "concussion", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 3, name: "call_of_flame", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 4, name: "elemental_warding", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 5, name: "elemental_devastation", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 6, name: "reverberation", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 7, name: "elemental_focus", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 8, name: "elemental_fury", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 9, name: "improved_fire_nova", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 10, name: "eye_of_the_storm", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 11, name: "elemental_reach", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 12, name: "call_of_thunder", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 13, name: "unrelenting_storm", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 14, name: "elemental_precision", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 15, name: "lightning_mastery", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 16, name: "elemental_mastery", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 17, name: "storm_earth_and_fire", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 18, name: "booming_echoes", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 19, name: "elemental_oath", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 20, name: "lightning_overload", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 21, name: "astral_shift", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 22, name: "totem_of_wrath", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 23, name: "lava_flows", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 24, name: "shamanism", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 25, name: "thunderstorm", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 26, name: "enhancing_totems", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 27, name: "earths_grasp", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 28, name: "ancestral_knowledge", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 29, name: "guardian_totems", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 30, name: "thundering_strikes", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 31, name: "improved_ghost_wolf", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 32, name: "improved_shields", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 33, name: "elemental_weapons", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 34, name: "shamanistic_focus", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 35, name: "anticipation", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 36, name: "flurry", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 37, name: "toughness", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 38, name: "improved_windfury_totem", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 39, name: "spirit_weapons", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 40, name: "mental_dexterity", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 41, name: "unleashed_rage", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 42, name: "weapon_mastery", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 43, name: "frozen_power", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 44, name: "dual_wield_specialization", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 45, name: "dual_wield", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 46, name: "stormstrike", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 47, name: "static_shock", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 48, name: "lava_lash", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 49, name: "improved_stormstrike", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 50, name: "mental_quickness", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 51, name: "shamanistic_rage", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 52, name: "earthen_power", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 53, name: "maelstrom_weapon", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 54, name: "feral_spirit", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 55, name: "improved_healing_wave", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 56, name: "totemic_focus", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 57, name: "improved_reincarnation", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 58, name: "healing_grace", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 59, name: "tidal_focus", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 60, name: "improved_water_shield", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 61, name: "healing_focus", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 62, name: "tidal_force", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 63, name: "ancestral_healing", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 64, name: "restorative_totems", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 65, name: "tidal_mastery", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 66, name: "healing_way", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 67, name: "natures_swiftness", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 68, name: "focused_mind", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 69, name: "purification", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 70, name: "natures_guardian", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 71, name: "mana_tide_totem", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 72, name: "cleanse_spirit", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 73, name: "blessing_of_the_eternals", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 74, name: "improved_chain_heal", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 75, name: "natures_blessing", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 76, name: "ancestral_awakening", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 77, name: "earth_shield", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 78, name: "improved_earth_shield", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 79, name: "tidal_waves", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 80, name: "riptide", kind: "scalar", T: 8 /*ScalarType.BOOL*/ }
        ]);
    }
    create(value) {
        const message = { convection: 0, concussion: 0, callOfFlame: 0, elementalWarding: 0, elementalDevastation: 0, reverberation: 0, elementalFocus: false, elementalFury: 0, improvedFireNova: 0, eyeOfTheStorm: 0, elementalReach: 0, callOfThunder: false, unrelentingStorm: 0, elementalPrecision: 0, lightningMastery: 0, elementalMastery: false, stormEarthAndFire: 0, boomingEchoes: 0, elementalOath: 0, lightningOverload: 0, astralShift: 0, totemOfWrath: false, lavaFlows: 0, shamanism: 0, thunderstorm: false, enhancingTotems: 0, earthsGrasp: 0, ancestralKnowledge: 0, guardianTotems: 0, thunderingStrikes: 0, improvedGhostWolf: 0, improvedShields: 0, elementalWeapons: 0, shamanisticFocus: false, anticipation: 0, flurry: 0, toughness: 0, improvedWindfuryTotem: 0, spiritWeapons: false, mentalDexterity: 0, unleashedRage: 0, weaponMastery: 0, frozenPower: 0, dualWieldSpecialization: 0, dualWield: false, stormstrike: false, staticShock: 0, lavaLash: false, improvedStormstrike: 0, mentalQuickness: 0, shamanisticRage: false, earthenPower: 0, maelstromWeapon: 0, feralSpirit: false, improvedHealingWave: 0, totemicFocus: 0, improvedReincarnation: 0, healingGrace: 0, tidalFocus: 0, improvedWaterShield: 0, healingFocus: 0, tidalForce: false, ancestralHealing: 0, restorativeTotems: 0, tidalMastery: 0, healingWay: 0, naturesSwiftness: false, focusedMind: 0, purification: 0, naturesGuardian: 0, manaTideTotem: false, cleanseSpirit: false, blessingOfTheEternals: 0, improvedChainHeal: 0, naturesBlessing: 0, ancestralAwakening: 0, earthShield: false, improvedEarthShield: 0, tidalWaves: 0, riptide: false };
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
                case /* int32 convection */ 1:
                    message.convection = reader.int32();
                    break;
                case /* int32 concussion */ 2:
                    message.concussion = reader.int32();
                    break;
                case /* int32 call_of_flame */ 3:
                    message.callOfFlame = reader.int32();
                    break;
                case /* int32 elemental_warding */ 4:
                    message.elementalWarding = reader.int32();
                    break;
                case /* int32 elemental_devastation */ 5:
                    message.elementalDevastation = reader.int32();
                    break;
                case /* int32 reverberation */ 6:
                    message.reverberation = reader.int32();
                    break;
                case /* bool elemental_focus */ 7:
                    message.elementalFocus = reader.bool();
                    break;
                case /* int32 elemental_fury */ 8:
                    message.elementalFury = reader.int32();
                    break;
                case /* int32 improved_fire_nova */ 9:
                    message.improvedFireNova = reader.int32();
                    break;
                case /* int32 eye_of_the_storm */ 10:
                    message.eyeOfTheStorm = reader.int32();
                    break;
                case /* int32 elemental_reach */ 11:
                    message.elementalReach = reader.int32();
                    break;
                case /* bool call_of_thunder */ 12:
                    message.callOfThunder = reader.bool();
                    break;
                case /* int32 unrelenting_storm */ 13:
                    message.unrelentingStorm = reader.int32();
                    break;
                case /* int32 elemental_precision */ 14:
                    message.elementalPrecision = reader.int32();
                    break;
                case /* int32 lightning_mastery */ 15:
                    message.lightningMastery = reader.int32();
                    break;
                case /* bool elemental_mastery */ 16:
                    message.elementalMastery = reader.bool();
                    break;
                case /* int32 storm_earth_and_fire */ 17:
                    message.stormEarthAndFire = reader.int32();
                    break;
                case /* int32 booming_echoes */ 18:
                    message.boomingEchoes = reader.int32();
                    break;
                case /* int32 elemental_oath */ 19:
                    message.elementalOath = reader.int32();
                    break;
                case /* int32 lightning_overload */ 20:
                    message.lightningOverload = reader.int32();
                    break;
                case /* int32 astral_shift */ 21:
                    message.astralShift = reader.int32();
                    break;
                case /* bool totem_of_wrath */ 22:
                    message.totemOfWrath = reader.bool();
                    break;
                case /* int32 lava_flows */ 23:
                    message.lavaFlows = reader.int32();
                    break;
                case /* int32 shamanism */ 24:
                    message.shamanism = reader.int32();
                    break;
                case /* bool thunderstorm */ 25:
                    message.thunderstorm = reader.bool();
                    break;
                case /* int32 enhancing_totems */ 26:
                    message.enhancingTotems = reader.int32();
                    break;
                case /* int32 earths_grasp */ 27:
                    message.earthsGrasp = reader.int32();
                    break;
                case /* int32 ancestral_knowledge */ 28:
                    message.ancestralKnowledge = reader.int32();
                    break;
                case /* int32 guardian_totems */ 29:
                    message.guardianTotems = reader.int32();
                    break;
                case /* int32 thundering_strikes */ 30:
                    message.thunderingStrikes = reader.int32();
                    break;
                case /* int32 improved_ghost_wolf */ 31:
                    message.improvedGhostWolf = reader.int32();
                    break;
                case /* int32 improved_shields */ 32:
                    message.improvedShields = reader.int32();
                    break;
                case /* int32 elemental_weapons */ 33:
                    message.elementalWeapons = reader.int32();
                    break;
                case /* bool shamanistic_focus */ 34:
                    message.shamanisticFocus = reader.bool();
                    break;
                case /* int32 anticipation */ 35:
                    message.anticipation = reader.int32();
                    break;
                case /* int32 flurry */ 36:
                    message.flurry = reader.int32();
                    break;
                case /* int32 toughness */ 37:
                    message.toughness = reader.int32();
                    break;
                case /* int32 improved_windfury_totem */ 38:
                    message.improvedWindfuryTotem = reader.int32();
                    break;
                case /* bool spirit_weapons */ 39:
                    message.spiritWeapons = reader.bool();
                    break;
                case /* int32 mental_dexterity */ 40:
                    message.mentalDexterity = reader.int32();
                    break;
                case /* int32 unleashed_rage */ 41:
                    message.unleashedRage = reader.int32();
                    break;
                case /* int32 weapon_mastery */ 42:
                    message.weaponMastery = reader.int32();
                    break;
                case /* int32 frozen_power */ 43:
                    message.frozenPower = reader.int32();
                    break;
                case /* int32 dual_wield_specialization */ 44:
                    message.dualWieldSpecialization = reader.int32();
                    break;
                case /* bool dual_wield */ 45:
                    message.dualWield = reader.bool();
                    break;
                case /* bool stormstrike */ 46:
                    message.stormstrike = reader.bool();
                    break;
                case /* int32 static_shock */ 47:
                    message.staticShock = reader.int32();
                    break;
                case /* bool lava_lash */ 48:
                    message.lavaLash = reader.bool();
                    break;
                case /* int32 improved_stormstrike */ 49:
                    message.improvedStormstrike = reader.int32();
                    break;
                case /* int32 mental_quickness */ 50:
                    message.mentalQuickness = reader.int32();
                    break;
                case /* bool shamanistic_rage */ 51:
                    message.shamanisticRage = reader.bool();
                    break;
                case /* int32 earthen_power */ 52:
                    message.earthenPower = reader.int32();
                    break;
                case /* int32 maelstrom_weapon */ 53:
                    message.maelstromWeapon = reader.int32();
                    break;
                case /* bool feral_spirit */ 54:
                    message.feralSpirit = reader.bool();
                    break;
                case /* int32 improved_healing_wave */ 55:
                    message.improvedHealingWave = reader.int32();
                    break;
                case /* int32 totemic_focus */ 56:
                    message.totemicFocus = reader.int32();
                    break;
                case /* int32 improved_reincarnation */ 57:
                    message.improvedReincarnation = reader.int32();
                    break;
                case /* int32 healing_grace */ 58:
                    message.healingGrace = reader.int32();
                    break;
                case /* int32 tidal_focus */ 59:
                    message.tidalFocus = reader.int32();
                    break;
                case /* int32 improved_water_shield */ 60:
                    message.improvedWaterShield = reader.int32();
                    break;
                case /* int32 healing_focus */ 61:
                    message.healingFocus = reader.int32();
                    break;
                case /* bool tidal_force */ 62:
                    message.tidalForce = reader.bool();
                    break;
                case /* int32 ancestral_healing */ 63:
                    message.ancestralHealing = reader.int32();
                    break;
                case /* int32 restorative_totems */ 64:
                    message.restorativeTotems = reader.int32();
                    break;
                case /* int32 tidal_mastery */ 65:
                    message.tidalMastery = reader.int32();
                    break;
                case /* int32 healing_way */ 66:
                    message.healingWay = reader.int32();
                    break;
                case /* bool natures_swiftness */ 67:
                    message.naturesSwiftness = reader.bool();
                    break;
                case /* int32 focused_mind */ 68:
                    message.focusedMind = reader.int32();
                    break;
                case /* int32 purification */ 69:
                    message.purification = reader.int32();
                    break;
                case /* int32 natures_guardian */ 70:
                    message.naturesGuardian = reader.int32();
                    break;
                case /* bool mana_tide_totem */ 71:
                    message.manaTideTotem = reader.bool();
                    break;
                case /* bool cleanse_spirit */ 72:
                    message.cleanseSpirit = reader.bool();
                    break;
                case /* int32 blessing_of_the_eternals */ 73:
                    message.blessingOfTheEternals = reader.int32();
                    break;
                case /* int32 improved_chain_heal */ 74:
                    message.improvedChainHeal = reader.int32();
                    break;
                case /* int32 natures_blessing */ 75:
                    message.naturesBlessing = reader.int32();
                    break;
                case /* int32 ancestral_awakening */ 76:
                    message.ancestralAwakening = reader.int32();
                    break;
                case /* bool earth_shield */ 77:
                    message.earthShield = reader.bool();
                    break;
                case /* int32 improved_earth_shield */ 78:
                    message.improvedEarthShield = reader.int32();
                    break;
                case /* int32 tidal_waves */ 79:
                    message.tidalWaves = reader.int32();
                    break;
                case /* bool riptide */ 80:
                    message.riptide = reader.bool();
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
        /* int32 convection = 1; */
        if (message.convection !== 0)
            writer.tag(1, WireType.Varint).int32(message.convection);
        /* int32 concussion = 2; */
        if (message.concussion !== 0)
            writer.tag(2, WireType.Varint).int32(message.concussion);
        /* int32 call_of_flame = 3; */
        if (message.callOfFlame !== 0)
            writer.tag(3, WireType.Varint).int32(message.callOfFlame);
        /* int32 elemental_warding = 4; */
        if (message.elementalWarding !== 0)
            writer.tag(4, WireType.Varint).int32(message.elementalWarding);
        /* int32 elemental_devastation = 5; */
        if (message.elementalDevastation !== 0)
            writer.tag(5, WireType.Varint).int32(message.elementalDevastation);
        /* int32 reverberation = 6; */
        if (message.reverberation !== 0)
            writer.tag(6, WireType.Varint).int32(message.reverberation);
        /* bool elemental_focus = 7; */
        if (message.elementalFocus !== false)
            writer.tag(7, WireType.Varint).bool(message.elementalFocus);
        /* int32 elemental_fury = 8; */
        if (message.elementalFury !== 0)
            writer.tag(8, WireType.Varint).int32(message.elementalFury);
        /* int32 improved_fire_nova = 9; */
        if (message.improvedFireNova !== 0)
            writer.tag(9, WireType.Varint).int32(message.improvedFireNova);
        /* int32 eye_of_the_storm = 10; */
        if (message.eyeOfTheStorm !== 0)
            writer.tag(10, WireType.Varint).int32(message.eyeOfTheStorm);
        /* int32 elemental_reach = 11; */
        if (message.elementalReach !== 0)
            writer.tag(11, WireType.Varint).int32(message.elementalReach);
        /* bool call_of_thunder = 12; */
        if (message.callOfThunder !== false)
            writer.tag(12, WireType.Varint).bool(message.callOfThunder);
        /* int32 unrelenting_storm = 13; */
        if (message.unrelentingStorm !== 0)
            writer.tag(13, WireType.Varint).int32(message.unrelentingStorm);
        /* int32 elemental_precision = 14; */
        if (message.elementalPrecision !== 0)
            writer.tag(14, WireType.Varint).int32(message.elementalPrecision);
        /* int32 lightning_mastery = 15; */
        if (message.lightningMastery !== 0)
            writer.tag(15, WireType.Varint).int32(message.lightningMastery);
        /* bool elemental_mastery = 16; */
        if (message.elementalMastery !== false)
            writer.tag(16, WireType.Varint).bool(message.elementalMastery);
        /* int32 storm_earth_and_fire = 17; */
        if (message.stormEarthAndFire !== 0)
            writer.tag(17, WireType.Varint).int32(message.stormEarthAndFire);
        /* int32 booming_echoes = 18; */
        if (message.boomingEchoes !== 0)
            writer.tag(18, WireType.Varint).int32(message.boomingEchoes);
        /* int32 elemental_oath = 19; */
        if (message.elementalOath !== 0)
            writer.tag(19, WireType.Varint).int32(message.elementalOath);
        /* int32 lightning_overload = 20; */
        if (message.lightningOverload !== 0)
            writer.tag(20, WireType.Varint).int32(message.lightningOverload);
        /* int32 astral_shift = 21; */
        if (message.astralShift !== 0)
            writer.tag(21, WireType.Varint).int32(message.astralShift);
        /* bool totem_of_wrath = 22; */
        if (message.totemOfWrath !== false)
            writer.tag(22, WireType.Varint).bool(message.totemOfWrath);
        /* int32 lava_flows = 23; */
        if (message.lavaFlows !== 0)
            writer.tag(23, WireType.Varint).int32(message.lavaFlows);
        /* int32 shamanism = 24; */
        if (message.shamanism !== 0)
            writer.tag(24, WireType.Varint).int32(message.shamanism);
        /* bool thunderstorm = 25; */
        if (message.thunderstorm !== false)
            writer.tag(25, WireType.Varint).bool(message.thunderstorm);
        /* int32 enhancing_totems = 26; */
        if (message.enhancingTotems !== 0)
            writer.tag(26, WireType.Varint).int32(message.enhancingTotems);
        /* int32 earths_grasp = 27; */
        if (message.earthsGrasp !== 0)
            writer.tag(27, WireType.Varint).int32(message.earthsGrasp);
        /* int32 ancestral_knowledge = 28; */
        if (message.ancestralKnowledge !== 0)
            writer.tag(28, WireType.Varint).int32(message.ancestralKnowledge);
        /* int32 guardian_totems = 29; */
        if (message.guardianTotems !== 0)
            writer.tag(29, WireType.Varint).int32(message.guardianTotems);
        /* int32 thundering_strikes = 30; */
        if (message.thunderingStrikes !== 0)
            writer.tag(30, WireType.Varint).int32(message.thunderingStrikes);
        /* int32 improved_ghost_wolf = 31; */
        if (message.improvedGhostWolf !== 0)
            writer.tag(31, WireType.Varint).int32(message.improvedGhostWolf);
        /* int32 improved_shields = 32; */
        if (message.improvedShields !== 0)
            writer.tag(32, WireType.Varint).int32(message.improvedShields);
        /* int32 elemental_weapons = 33; */
        if (message.elementalWeapons !== 0)
            writer.tag(33, WireType.Varint).int32(message.elementalWeapons);
        /* bool shamanistic_focus = 34; */
        if (message.shamanisticFocus !== false)
            writer.tag(34, WireType.Varint).bool(message.shamanisticFocus);
        /* int32 anticipation = 35; */
        if (message.anticipation !== 0)
            writer.tag(35, WireType.Varint).int32(message.anticipation);
        /* int32 flurry = 36; */
        if (message.flurry !== 0)
            writer.tag(36, WireType.Varint).int32(message.flurry);
        /* int32 toughness = 37; */
        if (message.toughness !== 0)
            writer.tag(37, WireType.Varint).int32(message.toughness);
        /* int32 improved_windfury_totem = 38; */
        if (message.improvedWindfuryTotem !== 0)
            writer.tag(38, WireType.Varint).int32(message.improvedWindfuryTotem);
        /* bool spirit_weapons = 39; */
        if (message.spiritWeapons !== false)
            writer.tag(39, WireType.Varint).bool(message.spiritWeapons);
        /* int32 mental_dexterity = 40; */
        if (message.mentalDexterity !== 0)
            writer.tag(40, WireType.Varint).int32(message.mentalDexterity);
        /* int32 unleashed_rage = 41; */
        if (message.unleashedRage !== 0)
            writer.tag(41, WireType.Varint).int32(message.unleashedRage);
        /* int32 weapon_mastery = 42; */
        if (message.weaponMastery !== 0)
            writer.tag(42, WireType.Varint).int32(message.weaponMastery);
        /* int32 frozen_power = 43; */
        if (message.frozenPower !== 0)
            writer.tag(43, WireType.Varint).int32(message.frozenPower);
        /* int32 dual_wield_specialization = 44; */
        if (message.dualWieldSpecialization !== 0)
            writer.tag(44, WireType.Varint).int32(message.dualWieldSpecialization);
        /* bool dual_wield = 45; */
        if (message.dualWield !== false)
            writer.tag(45, WireType.Varint).bool(message.dualWield);
        /* bool stormstrike = 46; */
        if (message.stormstrike !== false)
            writer.tag(46, WireType.Varint).bool(message.stormstrike);
        /* int32 static_shock = 47; */
        if (message.staticShock !== 0)
            writer.tag(47, WireType.Varint).int32(message.staticShock);
        /* bool lava_lash = 48; */
        if (message.lavaLash !== false)
            writer.tag(48, WireType.Varint).bool(message.lavaLash);
        /* int32 improved_stormstrike = 49; */
        if (message.improvedStormstrike !== 0)
            writer.tag(49, WireType.Varint).int32(message.improvedStormstrike);
        /* int32 mental_quickness = 50; */
        if (message.mentalQuickness !== 0)
            writer.tag(50, WireType.Varint).int32(message.mentalQuickness);
        /* bool shamanistic_rage = 51; */
        if (message.shamanisticRage !== false)
            writer.tag(51, WireType.Varint).bool(message.shamanisticRage);
        /* int32 earthen_power = 52; */
        if (message.earthenPower !== 0)
            writer.tag(52, WireType.Varint).int32(message.earthenPower);
        /* int32 maelstrom_weapon = 53; */
        if (message.maelstromWeapon !== 0)
            writer.tag(53, WireType.Varint).int32(message.maelstromWeapon);
        /* bool feral_spirit = 54; */
        if (message.feralSpirit !== false)
            writer.tag(54, WireType.Varint).bool(message.feralSpirit);
        /* int32 improved_healing_wave = 55; */
        if (message.improvedHealingWave !== 0)
            writer.tag(55, WireType.Varint).int32(message.improvedHealingWave);
        /* int32 totemic_focus = 56; */
        if (message.totemicFocus !== 0)
            writer.tag(56, WireType.Varint).int32(message.totemicFocus);
        /* int32 improved_reincarnation = 57; */
        if (message.improvedReincarnation !== 0)
            writer.tag(57, WireType.Varint).int32(message.improvedReincarnation);
        /* int32 healing_grace = 58; */
        if (message.healingGrace !== 0)
            writer.tag(58, WireType.Varint).int32(message.healingGrace);
        /* int32 tidal_focus = 59; */
        if (message.tidalFocus !== 0)
            writer.tag(59, WireType.Varint).int32(message.tidalFocus);
        /* int32 improved_water_shield = 60; */
        if (message.improvedWaterShield !== 0)
            writer.tag(60, WireType.Varint).int32(message.improvedWaterShield);
        /* int32 healing_focus = 61; */
        if (message.healingFocus !== 0)
            writer.tag(61, WireType.Varint).int32(message.healingFocus);
        /* bool tidal_force = 62; */
        if (message.tidalForce !== false)
            writer.tag(62, WireType.Varint).bool(message.tidalForce);
        /* int32 ancestral_healing = 63; */
        if (message.ancestralHealing !== 0)
            writer.tag(63, WireType.Varint).int32(message.ancestralHealing);
        /* int32 restorative_totems = 64; */
        if (message.restorativeTotems !== 0)
            writer.tag(64, WireType.Varint).int32(message.restorativeTotems);
        /* int32 tidal_mastery = 65; */
        if (message.tidalMastery !== 0)
            writer.tag(65, WireType.Varint).int32(message.tidalMastery);
        /* int32 healing_way = 66; */
        if (message.healingWay !== 0)
            writer.tag(66, WireType.Varint).int32(message.healingWay);
        /* bool natures_swiftness = 67; */
        if (message.naturesSwiftness !== false)
            writer.tag(67, WireType.Varint).bool(message.naturesSwiftness);
        /* int32 focused_mind = 68; */
        if (message.focusedMind !== 0)
            writer.tag(68, WireType.Varint).int32(message.focusedMind);
        /* int32 purification = 69; */
        if (message.purification !== 0)
            writer.tag(69, WireType.Varint).int32(message.purification);
        /* int32 natures_guardian = 70; */
        if (message.naturesGuardian !== 0)
            writer.tag(70, WireType.Varint).int32(message.naturesGuardian);
        /* bool mana_tide_totem = 71; */
        if (message.manaTideTotem !== false)
            writer.tag(71, WireType.Varint).bool(message.manaTideTotem);
        /* bool cleanse_spirit = 72; */
        if (message.cleanseSpirit !== false)
            writer.tag(72, WireType.Varint).bool(message.cleanseSpirit);
        /* int32 blessing_of_the_eternals = 73; */
        if (message.blessingOfTheEternals !== 0)
            writer.tag(73, WireType.Varint).int32(message.blessingOfTheEternals);
        /* int32 improved_chain_heal = 74; */
        if (message.improvedChainHeal !== 0)
            writer.tag(74, WireType.Varint).int32(message.improvedChainHeal);
        /* int32 natures_blessing = 75; */
        if (message.naturesBlessing !== 0)
            writer.tag(75, WireType.Varint).int32(message.naturesBlessing);
        /* int32 ancestral_awakening = 76; */
        if (message.ancestralAwakening !== 0)
            writer.tag(76, WireType.Varint).int32(message.ancestralAwakening);
        /* bool earth_shield = 77; */
        if (message.earthShield !== false)
            writer.tag(77, WireType.Varint).bool(message.earthShield);
        /* int32 improved_earth_shield = 78; */
        if (message.improvedEarthShield !== 0)
            writer.tag(78, WireType.Varint).int32(message.improvedEarthShield);
        /* int32 tidal_waves = 79; */
        if (message.tidalWaves !== 0)
            writer.tag(79, WireType.Varint).int32(message.tidalWaves);
        /* bool riptide = 80; */
        if (message.riptide !== false)
            writer.tag(80, WireType.Varint).bool(message.riptide);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message proto.ShamanTalents
 */
export const ShamanTalents = new ShamanTalents$Type();
// @generated message type with reflection information, may provide speed optimized methods
class ShamanTotems$Type extends MessageType {
    constructor() {
        super("proto.ShamanTotems", [
            { no: 1, name: "earth", kind: "enum", T: () => ["proto.EarthTotem", EarthTotem] },
            { no: 2, name: "air", kind: "enum", T: () => ["proto.AirTotem", AirTotem] },
            { no: 3, name: "fire", kind: "enum", T: () => ["proto.FireTotem", FireTotem] },
            { no: 4, name: "water", kind: "enum", T: () => ["proto.WaterTotem", WaterTotem] },
            { no: 5, name: "use_mana_tide", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 6, name: "use_fire_elemental", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 7, name: "recall_fire_elemental_on_oom", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 8, name: "recall_totems", kind: "scalar", T: 8 /*ScalarType.BOOL*/ }
        ]);
    }
    create(value) {
        const message = { earth: 0, air: 0, fire: 0, water: 0, useManaTide: false, useFireElemental: false, recallFireElementalOnOom: false, recallTotems: false };
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
                case /* proto.EarthTotem earth */ 1:
                    message.earth = reader.int32();
                    break;
                case /* proto.AirTotem air */ 2:
                    message.air = reader.int32();
                    break;
                case /* proto.FireTotem fire */ 3:
                    message.fire = reader.int32();
                    break;
                case /* proto.WaterTotem water */ 4:
                    message.water = reader.int32();
                    break;
                case /* bool use_mana_tide */ 5:
                    message.useManaTide = reader.bool();
                    break;
                case /* bool use_fire_elemental */ 6:
                    message.useFireElemental = reader.bool();
                    break;
                case /* bool recall_fire_elemental_on_oom */ 7:
                    message.recallFireElementalOnOom = reader.bool();
                    break;
                case /* bool recall_totems */ 8:
                    message.recallTotems = reader.bool();
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
        /* proto.EarthTotem earth = 1; */
        if (message.earth !== 0)
            writer.tag(1, WireType.Varint).int32(message.earth);
        /* proto.AirTotem air = 2; */
        if (message.air !== 0)
            writer.tag(2, WireType.Varint).int32(message.air);
        /* proto.FireTotem fire = 3; */
        if (message.fire !== 0)
            writer.tag(3, WireType.Varint).int32(message.fire);
        /* proto.WaterTotem water = 4; */
        if (message.water !== 0)
            writer.tag(4, WireType.Varint).int32(message.water);
        /* bool use_mana_tide = 5; */
        if (message.useManaTide !== false)
            writer.tag(5, WireType.Varint).bool(message.useManaTide);
        /* bool use_fire_elemental = 6; */
        if (message.useFireElemental !== false)
            writer.tag(6, WireType.Varint).bool(message.useFireElemental);
        /* bool recall_fire_elemental_on_oom = 7; */
        if (message.recallFireElementalOnOom !== false)
            writer.tag(7, WireType.Varint).bool(message.recallFireElementalOnOom);
        /* bool recall_totems = 8; */
        if (message.recallTotems !== false)
            writer.tag(8, WireType.Varint).bool(message.recallTotems);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message proto.ShamanTotems
 */
export const ShamanTotems = new ShamanTotems$Type();
// @generated message type with reflection information, may provide speed optimized methods
class ElementalShaman$Type extends MessageType {
    constructor() {
        super("proto.ElementalShaman", [
            { no: 1, name: "rotation", kind: "message", T: () => ElementalShaman_Rotation },
            { no: 2, name: "talents", kind: "message", T: () => ShamanTalents },
            { no: 3, name: "options", kind: "message", T: () => ElementalShaman_Options }
        ]);
    }
    create(value) {
        const message = {};
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
                case /* proto.ElementalShaman.Rotation rotation */ 1:
                    message.rotation = ElementalShaman_Rotation.internalBinaryRead(reader, reader.uint32(), options, message.rotation);
                    break;
                case /* proto.ShamanTalents talents */ 2:
                    message.talents = ShamanTalents.internalBinaryRead(reader, reader.uint32(), options, message.talents);
                    break;
                case /* proto.ElementalShaman.Options options */ 3:
                    message.options = ElementalShaman_Options.internalBinaryRead(reader, reader.uint32(), options, message.options);
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
        /* proto.ElementalShaman.Rotation rotation = 1; */
        if (message.rotation)
            ElementalShaman_Rotation.internalBinaryWrite(message.rotation, writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        /* proto.ShamanTalents talents = 2; */
        if (message.talents)
            ShamanTalents.internalBinaryWrite(message.talents, writer.tag(2, WireType.LengthDelimited).fork(), options).join();
        /* proto.ElementalShaman.Options options = 3; */
        if (message.options)
            ElementalShaman_Options.internalBinaryWrite(message.options, writer.tag(3, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message proto.ElementalShaman
 */
export const ElementalShaman = new ElementalShaman$Type();
// @generated message type with reflection information, may provide speed optimized methods
class ElementalShaman_Rotation$Type extends MessageType {
    constructor() {
        super("proto.ElementalShaman.Rotation", [
            { no: 3, name: "totems", kind: "message", T: () => ShamanTotems },
            { no: 1, name: "type", kind: "enum", T: () => ["proto.ElementalShaman.Rotation.RotationType", ElementalShaman_Rotation_RotationType] },
            { no: 2, name: "in_thunderstorm_range", kind: "scalar", T: 8 /*ScalarType.BOOL*/ }
        ]);
    }
    create(value) {
        const message = { type: 0, inThunderstormRange: false };
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
                case /* proto.ShamanTotems totems */ 3:
                    message.totems = ShamanTotems.internalBinaryRead(reader, reader.uint32(), options, message.totems);
                    break;
                case /* proto.ElementalShaman.Rotation.RotationType type */ 1:
                    message.type = reader.int32();
                    break;
                case /* bool in_thunderstorm_range */ 2:
                    message.inThunderstormRange = reader.bool();
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
        /* proto.ShamanTotems totems = 3; */
        if (message.totems)
            ShamanTotems.internalBinaryWrite(message.totems, writer.tag(3, WireType.LengthDelimited).fork(), options).join();
        /* proto.ElementalShaman.Rotation.RotationType type = 1; */
        if (message.type !== 0)
            writer.tag(1, WireType.Varint).int32(message.type);
        /* bool in_thunderstorm_range = 2; */
        if (message.inThunderstormRange !== false)
            writer.tag(2, WireType.Varint).bool(message.inThunderstormRange);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message proto.ElementalShaman.Rotation
 */
export const ElementalShaman_Rotation = new ElementalShaman_Rotation$Type();
// @generated message type with reflection information, may provide speed optimized methods
class ElementalShaman_Options$Type extends MessageType {
    constructor() {
        super("proto.ElementalShaman.Options", [
            { no: 1, name: "shield", kind: "enum", T: () => ["proto.ShamanShield", ShamanShield] },
            { no: 2, name: "bloodlust", kind: "scalar", T: 8 /*ScalarType.BOOL*/ }
        ]);
    }
    create(value) {
        const message = { shield: 0, bloodlust: false };
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
                case /* proto.ShamanShield shield */ 1:
                    message.shield = reader.int32();
                    break;
                case /* bool bloodlust */ 2:
                    message.bloodlust = reader.bool();
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
        /* proto.ShamanShield shield = 1; */
        if (message.shield !== 0)
            writer.tag(1, WireType.Varint).int32(message.shield);
        /* bool bloodlust = 2; */
        if (message.bloodlust !== false)
            writer.tag(2, WireType.Varint).bool(message.bloodlust);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message proto.ElementalShaman.Options
 */
export const ElementalShaman_Options = new ElementalShaman_Options$Type();
// @generated message type with reflection information, may provide speed optimized methods
class EnhancementShaman$Type extends MessageType {
    constructor() {
        super("proto.EnhancementShaman", [
            { no: 1, name: "rotation", kind: "message", T: () => EnhancementShaman_Rotation },
            { no: 2, name: "talents", kind: "message", T: () => ShamanTalents },
            { no: 3, name: "options", kind: "message", T: () => EnhancementShaman_Options }
        ]);
    }
    create(value) {
        const message = {};
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
                case /* proto.EnhancementShaman.Rotation rotation */ 1:
                    message.rotation = EnhancementShaman_Rotation.internalBinaryRead(reader, reader.uint32(), options, message.rotation);
                    break;
                case /* proto.ShamanTalents talents */ 2:
                    message.talents = ShamanTalents.internalBinaryRead(reader, reader.uint32(), options, message.talents);
                    break;
                case /* proto.EnhancementShaman.Options options */ 3:
                    message.options = EnhancementShaman_Options.internalBinaryRead(reader, reader.uint32(), options, message.options);
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
        /* proto.EnhancementShaman.Rotation rotation = 1; */
        if (message.rotation)
            EnhancementShaman_Rotation.internalBinaryWrite(message.rotation, writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        /* proto.ShamanTalents talents = 2; */
        if (message.talents)
            ShamanTalents.internalBinaryWrite(message.talents, writer.tag(2, WireType.LengthDelimited).fork(), options).join();
        /* proto.EnhancementShaman.Options options = 3; */
        if (message.options)
            EnhancementShaman_Options.internalBinaryWrite(message.options, writer.tag(3, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message proto.EnhancementShaman
 */
export const EnhancementShaman = new EnhancementShaman$Type();
// @generated message type with reflection information, may provide speed optimized methods
class EnhancementShaman_Rotation$Type extends MessageType {
    constructor() {
        super("proto.EnhancementShaman.Rotation", [
            { no: 1, name: "totems", kind: "message", T: () => ShamanTotems }
        ]);
    }
    create(value) {
        const message = {};
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
                case /* proto.ShamanTotems totems */ 1:
                    message.totems = ShamanTotems.internalBinaryRead(reader, reader.uint32(), options, message.totems);
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
        /* proto.ShamanTotems totems = 1; */
        if (message.totems)
            ShamanTotems.internalBinaryWrite(message.totems, writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message proto.EnhancementShaman.Rotation
 */
export const EnhancementShaman_Rotation = new EnhancementShaman_Rotation$Type();
// @generated message type with reflection information, may provide speed optimized methods
class EnhancementShaman_Options$Type extends MessageType {
    constructor() {
        super("proto.EnhancementShaman.Options", [
            { no: 1, name: "shield", kind: "enum", T: () => ["proto.ShamanShield", ShamanShield] },
            { no: 2, name: "bloodlust", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 3, name: "delay_offhand_swings", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 4, name: "imbueMH", kind: "enum", T: () => ["proto.ShamanImbue", ShamanImbue] },
            { no: 5, name: "imbueOH", kind: "enum", T: () => ["proto.ShamanImbue", ShamanImbue] }
        ]);
    }
    create(value) {
        const message = { shield: 0, bloodlust: false, delayOffhandSwings: false, imbueMH: 0, imbueOH: 0 };
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
                case /* proto.ShamanShield shield */ 1:
                    message.shield = reader.int32();
                    break;
                case /* bool bloodlust */ 2:
                    message.bloodlust = reader.bool();
                    break;
                case /* bool delay_offhand_swings */ 3:
                    message.delayOffhandSwings = reader.bool();
                    break;
                case /* proto.ShamanImbue imbueMH */ 4:
                    message.imbueMH = reader.int32();
                    break;
                case /* proto.ShamanImbue imbueOH */ 5:
                    message.imbueOH = reader.int32();
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
        /* proto.ShamanShield shield = 1; */
        if (message.shield !== 0)
            writer.tag(1, WireType.Varint).int32(message.shield);
        /* bool bloodlust = 2; */
        if (message.bloodlust !== false)
            writer.tag(2, WireType.Varint).bool(message.bloodlust);
        /* bool delay_offhand_swings = 3; */
        if (message.delayOffhandSwings !== false)
            writer.tag(3, WireType.Varint).bool(message.delayOffhandSwings);
        /* proto.ShamanImbue imbueMH = 4; */
        if (message.imbueMH !== 0)
            writer.tag(4, WireType.Varint).int32(message.imbueMH);
        /* proto.ShamanImbue imbueOH = 5; */
        if (message.imbueOH !== 0)
            writer.tag(5, WireType.Varint).int32(message.imbueOH);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message proto.EnhancementShaman.Options
 */
export const EnhancementShaman_Options = new EnhancementShaman_Options$Type();
