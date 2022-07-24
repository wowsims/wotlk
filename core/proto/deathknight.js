import { WireType } from '/wotlk/protobuf-ts/index.js';
import { UnknownFieldHandler } from '/wotlk/protobuf-ts/index.js';
import { reflectionMergePartial } from '/wotlk/protobuf-ts/index.js';
import { MESSAGE_TYPE } from '/wotlk/protobuf-ts/index.js';
import { MessageType } from '/wotlk/protobuf-ts/index.js';
/**
 * @generated from protobuf enum proto.DeathKnight.Rotation.ArmyOfTheDead
 */
export var DeathKnight_Rotation_ArmyOfTheDead;
(function (DeathKnight_Rotation_ArmyOfTheDead) {
    /**
     * @generated from protobuf enum value: DoNotUse = 0;
     */
    DeathKnight_Rotation_ArmyOfTheDead[DeathKnight_Rotation_ArmyOfTheDead["DoNotUse"] = 0] = "DoNotUse";
    /**
     * @generated from protobuf enum value: PreCast = 1;
     */
    DeathKnight_Rotation_ArmyOfTheDead[DeathKnight_Rotation_ArmyOfTheDead["PreCast"] = 1] = "PreCast";
    /**
     * @generated from protobuf enum value: AsMajorCd = 2;
     */
    DeathKnight_Rotation_ArmyOfTheDead[DeathKnight_Rotation_ArmyOfTheDead["AsMajorCd"] = 2] = "AsMajorCd";
})(DeathKnight_Rotation_ArmyOfTheDead || (DeathKnight_Rotation_ArmyOfTheDead = {}));
/**
 * @generated from protobuf enum proto.DeathKnightMajorGlyph
 */
export var DeathKnightMajorGlyph;
(function (DeathKnightMajorGlyph) {
    /**
     * @generated from protobuf enum value: DeathKnightMajorGlyphNone = 0;
     */
    DeathKnightMajorGlyph[DeathKnightMajorGlyph["DeathKnightMajorGlyphNone"] = 0] = "DeathKnightMajorGlyphNone";
    /**
     * @generated from protobuf enum value: GlyphOfAntiMagicShell = 43533;
     */
    DeathKnightMajorGlyph[DeathKnightMajorGlyph["GlyphOfAntiMagicShell"] = 43533] = "GlyphOfAntiMagicShell";
    /**
     * @generated from protobuf enum value: GlyphOfBloodStrike = 43826;
     */
    DeathKnightMajorGlyph[DeathKnightMajorGlyph["GlyphOfBloodStrike"] = 43826] = "GlyphOfBloodStrike";
    /**
     * @generated from protobuf enum value: GlyphOfBoneShield = 43536;
     */
    DeathKnightMajorGlyph[DeathKnightMajorGlyph["GlyphOfBoneShield"] = 43536] = "GlyphOfBoneShield";
    /**
     * @generated from protobuf enum value: GlyphOfChainsOfIce = 43537;
     */
    DeathKnightMajorGlyph[DeathKnightMajorGlyph["GlyphOfChainsOfIce"] = 43537] = "GlyphOfChainsOfIce";
    /**
     * @generated from protobuf enum value: GlyphOfDancingRuneWeapon = 45799;
     */
    DeathKnightMajorGlyph[DeathKnightMajorGlyph["GlyphOfDancingRuneWeapon"] = 45799] = "GlyphOfDancingRuneWeapon";
    /**
     * @generated from protobuf enum value: GlyphOfDarkCommand = 43538;
     */
    DeathKnightMajorGlyph[DeathKnightMajorGlyph["GlyphOfDarkCommand"] = 43538] = "GlyphOfDarkCommand";
    /**
     * @generated from protobuf enum value: GlyphOfDarkDeath = 45804;
     */
    DeathKnightMajorGlyph[DeathKnightMajorGlyph["GlyphOfDarkDeath"] = 45804] = "GlyphOfDarkDeath";
    /**
     * @generated from protobuf enum value: GlyphOfDeathAndDecay = 43542;
     */
    DeathKnightMajorGlyph[DeathKnightMajorGlyph["GlyphOfDeathAndDecay"] = 43542] = "GlyphOfDeathAndDecay";
    /**
     * @generated from protobuf enum value: GlyphOfDeathGrip = 43541;
     */
    DeathKnightMajorGlyph[DeathKnightMajorGlyph["GlyphOfDeathGrip"] = 43541] = "GlyphOfDeathGrip";
    /**
     * @generated from protobuf enum value: GlyphOfDeathStrike = 43827;
     */
    DeathKnightMajorGlyph[DeathKnightMajorGlyph["GlyphOfDeathStrike"] = 43827] = "GlyphOfDeathStrike";
    /**
     * @generated from protobuf enum value: GlyphOfDisease = 45805;
     */
    DeathKnightMajorGlyph[DeathKnightMajorGlyph["GlyphOfDisease"] = 45805] = "GlyphOfDisease";
    /**
     * @generated from protobuf enum value: GlyphOfFrostStrike = 43543;
     */
    DeathKnightMajorGlyph[DeathKnightMajorGlyph["GlyphOfFrostStrike"] = 43543] = "GlyphOfFrostStrike";
    /**
     * @generated from protobuf enum value: GlyphOfHeartStrike = 43534;
     */
    DeathKnightMajorGlyph[DeathKnightMajorGlyph["GlyphOfHeartStrike"] = 43534] = "GlyphOfHeartStrike";
    /**
     * @generated from protobuf enum value: GlyphOfHowlingBlast = 45806;
     */
    DeathKnightMajorGlyph[DeathKnightMajorGlyph["GlyphOfHowlingBlast"] = 45806] = "GlyphOfHowlingBlast";
    /**
     * @generated from protobuf enum value: GlyphOfHungeringCold = 45800;
     */
    DeathKnightMajorGlyph[DeathKnightMajorGlyph["GlyphOfHungeringCold"] = 45800] = "GlyphOfHungeringCold";
    /**
     * @generated from protobuf enum value: GlyphOfIceboundFortitude = 43545;
     */
    DeathKnightMajorGlyph[DeathKnightMajorGlyph["GlyphOfIceboundFortitude"] = 43545] = "GlyphOfIceboundFortitude";
    /**
     * @generated from protobuf enum value: GlyphOfIcyTouch = 43546;
     */
    DeathKnightMajorGlyph[DeathKnightMajorGlyph["GlyphOfIcyTouch"] = 43546] = "GlyphOfIcyTouch";
    /**
     * @generated from protobuf enum value: GlyphOfObliterate = 43547;
     */
    DeathKnightMajorGlyph[DeathKnightMajorGlyph["GlyphOfObliterate"] = 43547] = "GlyphOfObliterate";
    /**
     * @generated from protobuf enum value: GlyphOfPlagueStrike = 43548;
     */
    DeathKnightMajorGlyph[DeathKnightMajorGlyph["GlyphOfPlagueStrike"] = 43548] = "GlyphOfPlagueStrike";
    /**
     * @generated from protobuf enum value: GlyphOfRuneStrike = 43550;
     */
    DeathKnightMajorGlyph[DeathKnightMajorGlyph["GlyphOfRuneStrike"] = 43550] = "GlyphOfRuneStrike";
    /**
     * @generated from protobuf enum value: GlyphOfRuneTap = 43825;
     */
    DeathKnightMajorGlyph[DeathKnightMajorGlyph["GlyphOfRuneTap"] = 43825] = "GlyphOfRuneTap";
    /**
     * @generated from protobuf enum value: GlyphOfScourgeStrike = 43551;
     */
    DeathKnightMajorGlyph[DeathKnightMajorGlyph["GlyphOfScourgeStrike"] = 43551] = "GlyphOfScourgeStrike";
    /**
     * @generated from protobuf enum value: GlyphOfStrangulate = 43552;
     */
    DeathKnightMajorGlyph[DeathKnightMajorGlyph["GlyphOfStrangulate"] = 43552] = "GlyphOfStrangulate";
    /**
     * @generated from protobuf enum value: GlyphOfTheGhoul = 43549;
     */
    DeathKnightMajorGlyph[DeathKnightMajorGlyph["GlyphOfTheGhoul"] = 43549] = "GlyphOfTheGhoul";
    /**
     * @generated from protobuf enum value: GlyphOfUnbreakableArmor = 43553;
     */
    DeathKnightMajorGlyph[DeathKnightMajorGlyph["GlyphOfUnbreakableArmor"] = 43553] = "GlyphOfUnbreakableArmor";
    /**
     * @generated from protobuf enum value: GlyphOfUnholyBlight = 45803;
     */
    DeathKnightMajorGlyph[DeathKnightMajorGlyph["GlyphOfUnholyBlight"] = 45803] = "GlyphOfUnholyBlight";
    /**
     * @generated from protobuf enum value: GlyphOfVampiricBlood = 43554;
     */
    DeathKnightMajorGlyph[DeathKnightMajorGlyph["GlyphOfVampiricBlood"] = 43554] = "GlyphOfVampiricBlood";
})(DeathKnightMajorGlyph || (DeathKnightMajorGlyph = {}));
/**
 * @generated from protobuf enum proto.DeathKnightMinorGlyph
 */
export var DeathKnightMinorGlyph;
(function (DeathKnightMinorGlyph) {
    /**
     * @generated from protobuf enum value: DeathKnightMinorGlyphNone = 0;
     */
    DeathKnightMinorGlyph[DeathKnightMinorGlyph["DeathKnightMinorGlyphNone"] = 0] = "DeathKnightMinorGlyphNone";
    /**
     * @generated from protobuf enum value: GlyphOfBloodTap = 43535;
     */
    DeathKnightMinorGlyph[DeathKnightMinorGlyph["GlyphOfBloodTap"] = 43535] = "GlyphOfBloodTap";
    /**
     * @generated from protobuf enum value: GlyphOfCorpseExplosion = 43671;
     */
    DeathKnightMinorGlyph[DeathKnightMinorGlyph["GlyphOfCorpseExplosion"] = 43671] = "GlyphOfCorpseExplosion";
    /**
     * @generated from protobuf enum value: GlyphOfDeathSEmbrace = 43539;
     */
    DeathKnightMinorGlyph[DeathKnightMinorGlyph["GlyphOfDeathSEmbrace"] = 43539] = "GlyphOfDeathSEmbrace";
    /**
     * @generated from protobuf enum value: GlyphOfHornOfWinter = 43544;
     */
    DeathKnightMinorGlyph[DeathKnightMinorGlyph["GlyphOfHornOfWinter"] = 43544] = "GlyphOfHornOfWinter";
    /**
     * @generated from protobuf enum value: GlyphOfPestilence = 43672;
     */
    DeathKnightMinorGlyph[DeathKnightMinorGlyph["GlyphOfPestilence"] = 43672] = "GlyphOfPestilence";
    /**
     * @generated from protobuf enum value: GlyphOfRaiseDead = 43673;
     */
    DeathKnightMinorGlyph[DeathKnightMinorGlyph["GlyphOfRaiseDead"] = 43673] = "GlyphOfRaiseDead";
})(DeathKnightMinorGlyph || (DeathKnightMinorGlyph = {}));
// @generated message type with reflection information, may provide speed optimized methods
class DeathKnightTalents$Type extends MessageType {
    constructor() {
        super("proto.DeathKnightTalents", [
            { no: 1, name: "butchery", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 2, name: "subversion", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 3, name: "blade_barrier", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 4, name: "bladed_armor", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 5, name: "scent_of_blood", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 6, name: "two_handed_weapon_specialization", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 7, name: "rune_tap", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 8, name: "dark_conviction", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 9, name: "death_rune_mastery", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 10, name: "improved_rune_tap", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 11, name: "spell_deflection", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 12, name: "vendetta", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 13, name: "bloody_strikes", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 14, name: "veteran_of_the_third_war", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 15, name: "mark_of_blood", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 16, name: "bloody_vengeance", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 17, name: "abominations_might", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 18, name: "bloodworms", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 19, name: "hysteria", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 20, name: "improved_blood_presence", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 21, name: "improved_death_strike", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 22, name: "sudden_doom", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 23, name: "vampiric_blood", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 24, name: "will_of_the_necropolis", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 25, name: "heart_strike", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 26, name: "might_of_mograine", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 27, name: "blood_gorged", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 28, name: "dancing_rune_weapon", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 29, name: "improved_icy_touch", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 30, name: "runic_power_mastery", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 31, name: "toughness", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 32, name: "icy_reach", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 33, name: "black_ice", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 34, name: "nerves_of_cold_steel", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 35, name: "icy_talons", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 36, name: "lichborne", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 37, name: "annihilation", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 38, name: "killing_machine", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 39, name: "chill_of_the_grave", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 40, name: "endless_winter", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 41, name: "frigid_dreadplate", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 42, name: "glacier_rot", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 43, name: "deathchill", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 44, name: "improved_icy_talons", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 45, name: "merciless_combat", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 46, name: "rime", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 47, name: "chilblains", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 48, name: "hungering_cold", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 49, name: "improved_frost_presence", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 50, name: "threat_of_thassarian", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 51, name: "blood_of_the_north", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 52, name: "unbreakable_armor", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 53, name: "acclimation", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 54, name: "frost_strike", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 55, name: "guile_of_gorefiend", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 56, name: "tundra_stalker", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 57, name: "howling_blast", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 58, name: "vicious_strikes", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 59, name: "virulence", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 60, name: "anticipation", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 61, name: "epidemic", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 62, name: "morbidity", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 63, name: "unholy_command", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 64, name: "ravenous_dead", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 65, name: "outbreak", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 66, name: "necrosis", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 67, name: "corpse_explosion", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 68, name: "on_a_pale_horse", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 69, name: "blood_caked_blade", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 70, name: "night_of_the_dead", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 71, name: "unholy_blight", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 72, name: "impurity", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 73, name: "dirge", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 74, name: "desecration", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 75, name: "magic_suppression", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 76, name: "reaping", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 77, name: "master_of_ghouls", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 78, name: "desolation", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 79, name: "anti_magic_zone", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 80, name: "improved_unholy_presence", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 81, name: "ghoul_frenzy", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 82, name: "crypt_fever", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 83, name: "bone_shield", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 84, name: "wandering_plague", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 85, name: "ebon_plaguebringer", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 86, name: "scourge_strike", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 87, name: "rage_of_rivendare", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 88, name: "summon_gargoyle", kind: "scalar", T: 8 /*ScalarType.BOOL*/ }
        ]);
    }
    create(value) {
        const message = { butchery: 0, subversion: 0, bladeBarrier: 0, bladedArmor: 0, scentOfBlood: 0, twoHandedWeaponSpecialization: 0, runeTap: false, darkConviction: 0, deathRuneMastery: 0, improvedRuneTap: 0, spellDeflection: 0, vendetta: 0, bloodyStrikes: 0, veteranOfTheThirdWar: 0, markOfBlood: false, bloodyVengeance: 0, abominationsMight: 0, bloodworms: 0, hysteria: false, improvedBloodPresence: 0, improvedDeathStrike: 0, suddenDoom: 0, vampiricBlood: false, willOfTheNecropolis: 0, heartStrike: false, mightOfMograine: 0, bloodGorged: 0, dancingRuneWeapon: false, improvedIcyTouch: 0, runicPowerMastery: 0, toughness: 0, icyReach: 0, blackIce: 0, nervesOfColdSteel: 0, icyTalons: 0, lichborne: false, annihilation: 0, killingMachine: 0, chillOfTheGrave: 0, endlessWinter: 0, frigidDreadplate: 0, glacierRot: 0, deathchill: false, improvedIcyTalons: false, mercilessCombat: 0, rime: 0, chilblains: 0, hungeringCold: false, improvedFrostPresence: 0, threatOfThassarian: 0, bloodOfTheNorth: 0, unbreakableArmor: false, acclimation: 0, frostStrike: false, guileOfGorefiend: 0, tundraStalker: 0, howlingBlast: false, viciousStrikes: 0, virulence: 0, anticipation: 0, epidemic: 0, morbidity: 0, unholyCommand: 0, ravenousDead: 0, outbreak: 0, necrosis: 0, corpseExplosion: false, onAPaleHorse: 0, bloodCakedBlade: 0, nightOfTheDead: 0, unholyBlight: false, impurity: 0, dirge: 0, desecration: 0, magicSuppression: 0, reaping: 0, masterOfGhouls: false, desolation: 0, antiMagicZone: false, improvedUnholyPresence: 0, ghoulFrenzy: false, cryptFever: 0, boneShield: false, wanderingPlague: 0, ebonPlaguebringer: 0, scourgeStrike: false, rageOfRivendare: 0, summonGargoyle: false };
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
                case /* int32 butchery */ 1:
                    message.butchery = reader.int32();
                    break;
                case /* int32 subversion */ 2:
                    message.subversion = reader.int32();
                    break;
                case /* int32 blade_barrier */ 3:
                    message.bladeBarrier = reader.int32();
                    break;
                case /* int32 bladed_armor */ 4:
                    message.bladedArmor = reader.int32();
                    break;
                case /* int32 scent_of_blood */ 5:
                    message.scentOfBlood = reader.int32();
                    break;
                case /* int32 two_handed_weapon_specialization */ 6:
                    message.twoHandedWeaponSpecialization = reader.int32();
                    break;
                case /* bool rune_tap */ 7:
                    message.runeTap = reader.bool();
                    break;
                case /* int32 dark_conviction */ 8:
                    message.darkConviction = reader.int32();
                    break;
                case /* int32 death_rune_mastery */ 9:
                    message.deathRuneMastery = reader.int32();
                    break;
                case /* int32 improved_rune_tap */ 10:
                    message.improvedRuneTap = reader.int32();
                    break;
                case /* int32 spell_deflection */ 11:
                    message.spellDeflection = reader.int32();
                    break;
                case /* int32 vendetta */ 12:
                    message.vendetta = reader.int32();
                    break;
                case /* int32 bloody_strikes */ 13:
                    message.bloodyStrikes = reader.int32();
                    break;
                case /* int32 veteran_of_the_third_war */ 14:
                    message.veteranOfTheThirdWar = reader.int32();
                    break;
                case /* bool mark_of_blood */ 15:
                    message.markOfBlood = reader.bool();
                    break;
                case /* int32 bloody_vengeance */ 16:
                    message.bloodyVengeance = reader.int32();
                    break;
                case /* int32 abominations_might */ 17:
                    message.abominationsMight = reader.int32();
                    break;
                case /* int32 bloodworms */ 18:
                    message.bloodworms = reader.int32();
                    break;
                case /* bool hysteria */ 19:
                    message.hysteria = reader.bool();
                    break;
                case /* int32 improved_blood_presence */ 20:
                    message.improvedBloodPresence = reader.int32();
                    break;
                case /* int32 improved_death_strike */ 21:
                    message.improvedDeathStrike = reader.int32();
                    break;
                case /* int32 sudden_doom */ 22:
                    message.suddenDoom = reader.int32();
                    break;
                case /* bool vampiric_blood */ 23:
                    message.vampiricBlood = reader.bool();
                    break;
                case /* int32 will_of_the_necropolis */ 24:
                    message.willOfTheNecropolis = reader.int32();
                    break;
                case /* bool heart_strike */ 25:
                    message.heartStrike = reader.bool();
                    break;
                case /* int32 might_of_mograine */ 26:
                    message.mightOfMograine = reader.int32();
                    break;
                case /* int32 blood_gorged */ 27:
                    message.bloodGorged = reader.int32();
                    break;
                case /* bool dancing_rune_weapon */ 28:
                    message.dancingRuneWeapon = reader.bool();
                    break;
                case /* int32 improved_icy_touch */ 29:
                    message.improvedIcyTouch = reader.int32();
                    break;
                case /* int32 runic_power_mastery */ 30:
                    message.runicPowerMastery = reader.int32();
                    break;
                case /* int32 toughness */ 31:
                    message.toughness = reader.int32();
                    break;
                case /* int32 icy_reach */ 32:
                    message.icyReach = reader.int32();
                    break;
                case /* int32 black_ice */ 33:
                    message.blackIce = reader.int32();
                    break;
                case /* int32 nerves_of_cold_steel */ 34:
                    message.nervesOfColdSteel = reader.int32();
                    break;
                case /* int32 icy_talons */ 35:
                    message.icyTalons = reader.int32();
                    break;
                case /* bool lichborne */ 36:
                    message.lichborne = reader.bool();
                    break;
                case /* int32 annihilation */ 37:
                    message.annihilation = reader.int32();
                    break;
                case /* int32 killing_machine */ 38:
                    message.killingMachine = reader.int32();
                    break;
                case /* int32 chill_of_the_grave */ 39:
                    message.chillOfTheGrave = reader.int32();
                    break;
                case /* int32 endless_winter */ 40:
                    message.endlessWinter = reader.int32();
                    break;
                case /* int32 frigid_dreadplate */ 41:
                    message.frigidDreadplate = reader.int32();
                    break;
                case /* int32 glacier_rot */ 42:
                    message.glacierRot = reader.int32();
                    break;
                case /* bool deathchill */ 43:
                    message.deathchill = reader.bool();
                    break;
                case /* bool improved_icy_talons */ 44:
                    message.improvedIcyTalons = reader.bool();
                    break;
                case /* int32 merciless_combat */ 45:
                    message.mercilessCombat = reader.int32();
                    break;
                case /* int32 rime */ 46:
                    message.rime = reader.int32();
                    break;
                case /* int32 chilblains */ 47:
                    message.chilblains = reader.int32();
                    break;
                case /* bool hungering_cold */ 48:
                    message.hungeringCold = reader.bool();
                    break;
                case /* int32 improved_frost_presence */ 49:
                    message.improvedFrostPresence = reader.int32();
                    break;
                case /* int32 threat_of_thassarian */ 50:
                    message.threatOfThassarian = reader.int32();
                    break;
                case /* int32 blood_of_the_north */ 51:
                    message.bloodOfTheNorth = reader.int32();
                    break;
                case /* bool unbreakable_armor */ 52:
                    message.unbreakableArmor = reader.bool();
                    break;
                case /* int32 acclimation */ 53:
                    message.acclimation = reader.int32();
                    break;
                case /* bool frost_strike */ 54:
                    message.frostStrike = reader.bool();
                    break;
                case /* int32 guile_of_gorefiend */ 55:
                    message.guileOfGorefiend = reader.int32();
                    break;
                case /* int32 tundra_stalker */ 56:
                    message.tundraStalker = reader.int32();
                    break;
                case /* bool howling_blast */ 57:
                    message.howlingBlast = reader.bool();
                    break;
                case /* int32 vicious_strikes */ 58:
                    message.viciousStrikes = reader.int32();
                    break;
                case /* int32 virulence */ 59:
                    message.virulence = reader.int32();
                    break;
                case /* int32 anticipation */ 60:
                    message.anticipation = reader.int32();
                    break;
                case /* int32 epidemic */ 61:
                    message.epidemic = reader.int32();
                    break;
                case /* int32 morbidity */ 62:
                    message.morbidity = reader.int32();
                    break;
                case /* int32 unholy_command */ 63:
                    message.unholyCommand = reader.int32();
                    break;
                case /* int32 ravenous_dead */ 64:
                    message.ravenousDead = reader.int32();
                    break;
                case /* int32 outbreak */ 65:
                    message.outbreak = reader.int32();
                    break;
                case /* int32 necrosis */ 66:
                    message.necrosis = reader.int32();
                    break;
                case /* bool corpse_explosion */ 67:
                    message.corpseExplosion = reader.bool();
                    break;
                case /* int32 on_a_pale_horse */ 68:
                    message.onAPaleHorse = reader.int32();
                    break;
                case /* int32 blood_caked_blade */ 69:
                    message.bloodCakedBlade = reader.int32();
                    break;
                case /* int32 night_of_the_dead */ 70:
                    message.nightOfTheDead = reader.int32();
                    break;
                case /* bool unholy_blight */ 71:
                    message.unholyBlight = reader.bool();
                    break;
                case /* int32 impurity */ 72:
                    message.impurity = reader.int32();
                    break;
                case /* int32 dirge */ 73:
                    message.dirge = reader.int32();
                    break;
                case /* int32 desecration */ 74:
                    message.desecration = reader.int32();
                    break;
                case /* int32 magic_suppression */ 75:
                    message.magicSuppression = reader.int32();
                    break;
                case /* int32 reaping */ 76:
                    message.reaping = reader.int32();
                    break;
                case /* bool master_of_ghouls */ 77:
                    message.masterOfGhouls = reader.bool();
                    break;
                case /* int32 desolation */ 78:
                    message.desolation = reader.int32();
                    break;
                case /* bool anti_magic_zone */ 79:
                    message.antiMagicZone = reader.bool();
                    break;
                case /* int32 improved_unholy_presence */ 80:
                    message.improvedUnholyPresence = reader.int32();
                    break;
                case /* bool ghoul_frenzy */ 81:
                    message.ghoulFrenzy = reader.bool();
                    break;
                case /* int32 crypt_fever */ 82:
                    message.cryptFever = reader.int32();
                    break;
                case /* bool bone_shield */ 83:
                    message.boneShield = reader.bool();
                    break;
                case /* int32 wandering_plague */ 84:
                    message.wanderingPlague = reader.int32();
                    break;
                case /* int32 ebon_plaguebringer */ 85:
                    message.ebonPlaguebringer = reader.int32();
                    break;
                case /* bool scourge_strike */ 86:
                    message.scourgeStrike = reader.bool();
                    break;
                case /* int32 rage_of_rivendare */ 87:
                    message.rageOfRivendare = reader.int32();
                    break;
                case /* bool summon_gargoyle */ 88:
                    message.summonGargoyle = reader.bool();
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
        /* int32 butchery = 1; */
        if (message.butchery !== 0)
            writer.tag(1, WireType.Varint).int32(message.butchery);
        /* int32 subversion = 2; */
        if (message.subversion !== 0)
            writer.tag(2, WireType.Varint).int32(message.subversion);
        /* int32 blade_barrier = 3; */
        if (message.bladeBarrier !== 0)
            writer.tag(3, WireType.Varint).int32(message.bladeBarrier);
        /* int32 bladed_armor = 4; */
        if (message.bladedArmor !== 0)
            writer.tag(4, WireType.Varint).int32(message.bladedArmor);
        /* int32 scent_of_blood = 5; */
        if (message.scentOfBlood !== 0)
            writer.tag(5, WireType.Varint).int32(message.scentOfBlood);
        /* int32 two_handed_weapon_specialization = 6; */
        if (message.twoHandedWeaponSpecialization !== 0)
            writer.tag(6, WireType.Varint).int32(message.twoHandedWeaponSpecialization);
        /* bool rune_tap = 7; */
        if (message.runeTap !== false)
            writer.tag(7, WireType.Varint).bool(message.runeTap);
        /* int32 dark_conviction = 8; */
        if (message.darkConviction !== 0)
            writer.tag(8, WireType.Varint).int32(message.darkConviction);
        /* int32 death_rune_mastery = 9; */
        if (message.deathRuneMastery !== 0)
            writer.tag(9, WireType.Varint).int32(message.deathRuneMastery);
        /* int32 improved_rune_tap = 10; */
        if (message.improvedRuneTap !== 0)
            writer.tag(10, WireType.Varint).int32(message.improvedRuneTap);
        /* int32 spell_deflection = 11; */
        if (message.spellDeflection !== 0)
            writer.tag(11, WireType.Varint).int32(message.spellDeflection);
        /* int32 vendetta = 12; */
        if (message.vendetta !== 0)
            writer.tag(12, WireType.Varint).int32(message.vendetta);
        /* int32 bloody_strikes = 13; */
        if (message.bloodyStrikes !== 0)
            writer.tag(13, WireType.Varint).int32(message.bloodyStrikes);
        /* int32 veteran_of_the_third_war = 14; */
        if (message.veteranOfTheThirdWar !== 0)
            writer.tag(14, WireType.Varint).int32(message.veteranOfTheThirdWar);
        /* bool mark_of_blood = 15; */
        if (message.markOfBlood !== false)
            writer.tag(15, WireType.Varint).bool(message.markOfBlood);
        /* int32 bloody_vengeance = 16; */
        if (message.bloodyVengeance !== 0)
            writer.tag(16, WireType.Varint).int32(message.bloodyVengeance);
        /* int32 abominations_might = 17; */
        if (message.abominationsMight !== 0)
            writer.tag(17, WireType.Varint).int32(message.abominationsMight);
        /* int32 bloodworms = 18; */
        if (message.bloodworms !== 0)
            writer.tag(18, WireType.Varint).int32(message.bloodworms);
        /* bool hysteria = 19; */
        if (message.hysteria !== false)
            writer.tag(19, WireType.Varint).bool(message.hysteria);
        /* int32 improved_blood_presence = 20; */
        if (message.improvedBloodPresence !== 0)
            writer.tag(20, WireType.Varint).int32(message.improvedBloodPresence);
        /* int32 improved_death_strike = 21; */
        if (message.improvedDeathStrike !== 0)
            writer.tag(21, WireType.Varint).int32(message.improvedDeathStrike);
        /* int32 sudden_doom = 22; */
        if (message.suddenDoom !== 0)
            writer.tag(22, WireType.Varint).int32(message.suddenDoom);
        /* bool vampiric_blood = 23; */
        if (message.vampiricBlood !== false)
            writer.tag(23, WireType.Varint).bool(message.vampiricBlood);
        /* int32 will_of_the_necropolis = 24; */
        if (message.willOfTheNecropolis !== 0)
            writer.tag(24, WireType.Varint).int32(message.willOfTheNecropolis);
        /* bool heart_strike = 25; */
        if (message.heartStrike !== false)
            writer.tag(25, WireType.Varint).bool(message.heartStrike);
        /* int32 might_of_mograine = 26; */
        if (message.mightOfMograine !== 0)
            writer.tag(26, WireType.Varint).int32(message.mightOfMograine);
        /* int32 blood_gorged = 27; */
        if (message.bloodGorged !== 0)
            writer.tag(27, WireType.Varint).int32(message.bloodGorged);
        /* bool dancing_rune_weapon = 28; */
        if (message.dancingRuneWeapon !== false)
            writer.tag(28, WireType.Varint).bool(message.dancingRuneWeapon);
        /* int32 improved_icy_touch = 29; */
        if (message.improvedIcyTouch !== 0)
            writer.tag(29, WireType.Varint).int32(message.improvedIcyTouch);
        /* int32 runic_power_mastery = 30; */
        if (message.runicPowerMastery !== 0)
            writer.tag(30, WireType.Varint).int32(message.runicPowerMastery);
        /* int32 toughness = 31; */
        if (message.toughness !== 0)
            writer.tag(31, WireType.Varint).int32(message.toughness);
        /* int32 icy_reach = 32; */
        if (message.icyReach !== 0)
            writer.tag(32, WireType.Varint).int32(message.icyReach);
        /* int32 black_ice = 33; */
        if (message.blackIce !== 0)
            writer.tag(33, WireType.Varint).int32(message.blackIce);
        /* int32 nerves_of_cold_steel = 34; */
        if (message.nervesOfColdSteel !== 0)
            writer.tag(34, WireType.Varint).int32(message.nervesOfColdSteel);
        /* int32 icy_talons = 35; */
        if (message.icyTalons !== 0)
            writer.tag(35, WireType.Varint).int32(message.icyTalons);
        /* bool lichborne = 36; */
        if (message.lichborne !== false)
            writer.tag(36, WireType.Varint).bool(message.lichborne);
        /* int32 annihilation = 37; */
        if (message.annihilation !== 0)
            writer.tag(37, WireType.Varint).int32(message.annihilation);
        /* int32 killing_machine = 38; */
        if (message.killingMachine !== 0)
            writer.tag(38, WireType.Varint).int32(message.killingMachine);
        /* int32 chill_of_the_grave = 39; */
        if (message.chillOfTheGrave !== 0)
            writer.tag(39, WireType.Varint).int32(message.chillOfTheGrave);
        /* int32 endless_winter = 40; */
        if (message.endlessWinter !== 0)
            writer.tag(40, WireType.Varint).int32(message.endlessWinter);
        /* int32 frigid_dreadplate = 41; */
        if (message.frigidDreadplate !== 0)
            writer.tag(41, WireType.Varint).int32(message.frigidDreadplate);
        /* int32 glacier_rot = 42; */
        if (message.glacierRot !== 0)
            writer.tag(42, WireType.Varint).int32(message.glacierRot);
        /* bool deathchill = 43; */
        if (message.deathchill !== false)
            writer.tag(43, WireType.Varint).bool(message.deathchill);
        /* bool improved_icy_talons = 44; */
        if (message.improvedIcyTalons !== false)
            writer.tag(44, WireType.Varint).bool(message.improvedIcyTalons);
        /* int32 merciless_combat = 45; */
        if (message.mercilessCombat !== 0)
            writer.tag(45, WireType.Varint).int32(message.mercilessCombat);
        /* int32 rime = 46; */
        if (message.rime !== 0)
            writer.tag(46, WireType.Varint).int32(message.rime);
        /* int32 chilblains = 47; */
        if (message.chilblains !== 0)
            writer.tag(47, WireType.Varint).int32(message.chilblains);
        /* bool hungering_cold = 48; */
        if (message.hungeringCold !== false)
            writer.tag(48, WireType.Varint).bool(message.hungeringCold);
        /* int32 improved_frost_presence = 49; */
        if (message.improvedFrostPresence !== 0)
            writer.tag(49, WireType.Varint).int32(message.improvedFrostPresence);
        /* int32 threat_of_thassarian = 50; */
        if (message.threatOfThassarian !== 0)
            writer.tag(50, WireType.Varint).int32(message.threatOfThassarian);
        /* int32 blood_of_the_north = 51; */
        if (message.bloodOfTheNorth !== 0)
            writer.tag(51, WireType.Varint).int32(message.bloodOfTheNorth);
        /* bool unbreakable_armor = 52; */
        if (message.unbreakableArmor !== false)
            writer.tag(52, WireType.Varint).bool(message.unbreakableArmor);
        /* int32 acclimation = 53; */
        if (message.acclimation !== 0)
            writer.tag(53, WireType.Varint).int32(message.acclimation);
        /* bool frost_strike = 54; */
        if (message.frostStrike !== false)
            writer.tag(54, WireType.Varint).bool(message.frostStrike);
        /* int32 guile_of_gorefiend = 55; */
        if (message.guileOfGorefiend !== 0)
            writer.tag(55, WireType.Varint).int32(message.guileOfGorefiend);
        /* int32 tundra_stalker = 56; */
        if (message.tundraStalker !== 0)
            writer.tag(56, WireType.Varint).int32(message.tundraStalker);
        /* bool howling_blast = 57; */
        if (message.howlingBlast !== false)
            writer.tag(57, WireType.Varint).bool(message.howlingBlast);
        /* int32 vicious_strikes = 58; */
        if (message.viciousStrikes !== 0)
            writer.tag(58, WireType.Varint).int32(message.viciousStrikes);
        /* int32 virulence = 59; */
        if (message.virulence !== 0)
            writer.tag(59, WireType.Varint).int32(message.virulence);
        /* int32 anticipation = 60; */
        if (message.anticipation !== 0)
            writer.tag(60, WireType.Varint).int32(message.anticipation);
        /* int32 epidemic = 61; */
        if (message.epidemic !== 0)
            writer.tag(61, WireType.Varint).int32(message.epidemic);
        /* int32 morbidity = 62; */
        if (message.morbidity !== 0)
            writer.tag(62, WireType.Varint).int32(message.morbidity);
        /* int32 unholy_command = 63; */
        if (message.unholyCommand !== 0)
            writer.tag(63, WireType.Varint).int32(message.unholyCommand);
        /* int32 ravenous_dead = 64; */
        if (message.ravenousDead !== 0)
            writer.tag(64, WireType.Varint).int32(message.ravenousDead);
        /* int32 outbreak = 65; */
        if (message.outbreak !== 0)
            writer.tag(65, WireType.Varint).int32(message.outbreak);
        /* int32 necrosis = 66; */
        if (message.necrosis !== 0)
            writer.tag(66, WireType.Varint).int32(message.necrosis);
        /* bool corpse_explosion = 67; */
        if (message.corpseExplosion !== false)
            writer.tag(67, WireType.Varint).bool(message.corpseExplosion);
        /* int32 on_a_pale_horse = 68; */
        if (message.onAPaleHorse !== 0)
            writer.tag(68, WireType.Varint).int32(message.onAPaleHorse);
        /* int32 blood_caked_blade = 69; */
        if (message.bloodCakedBlade !== 0)
            writer.tag(69, WireType.Varint).int32(message.bloodCakedBlade);
        /* int32 night_of_the_dead = 70; */
        if (message.nightOfTheDead !== 0)
            writer.tag(70, WireType.Varint).int32(message.nightOfTheDead);
        /* bool unholy_blight = 71; */
        if (message.unholyBlight !== false)
            writer.tag(71, WireType.Varint).bool(message.unholyBlight);
        /* int32 impurity = 72; */
        if (message.impurity !== 0)
            writer.tag(72, WireType.Varint).int32(message.impurity);
        /* int32 dirge = 73; */
        if (message.dirge !== 0)
            writer.tag(73, WireType.Varint).int32(message.dirge);
        /* int32 desecration = 74; */
        if (message.desecration !== 0)
            writer.tag(74, WireType.Varint).int32(message.desecration);
        /* int32 magic_suppression = 75; */
        if (message.magicSuppression !== 0)
            writer.tag(75, WireType.Varint).int32(message.magicSuppression);
        /* int32 reaping = 76; */
        if (message.reaping !== 0)
            writer.tag(76, WireType.Varint).int32(message.reaping);
        /* bool master_of_ghouls = 77; */
        if (message.masterOfGhouls !== false)
            writer.tag(77, WireType.Varint).bool(message.masterOfGhouls);
        /* int32 desolation = 78; */
        if (message.desolation !== 0)
            writer.tag(78, WireType.Varint).int32(message.desolation);
        /* bool anti_magic_zone = 79; */
        if (message.antiMagicZone !== false)
            writer.tag(79, WireType.Varint).bool(message.antiMagicZone);
        /* int32 improved_unholy_presence = 80; */
        if (message.improvedUnholyPresence !== 0)
            writer.tag(80, WireType.Varint).int32(message.improvedUnholyPresence);
        /* bool ghoul_frenzy = 81; */
        if (message.ghoulFrenzy !== false)
            writer.tag(81, WireType.Varint).bool(message.ghoulFrenzy);
        /* int32 crypt_fever = 82; */
        if (message.cryptFever !== 0)
            writer.tag(82, WireType.Varint).int32(message.cryptFever);
        /* bool bone_shield = 83; */
        if (message.boneShield !== false)
            writer.tag(83, WireType.Varint).bool(message.boneShield);
        /* int32 wandering_plague = 84; */
        if (message.wanderingPlague !== 0)
            writer.tag(84, WireType.Varint).int32(message.wanderingPlague);
        /* int32 ebon_plaguebringer = 85; */
        if (message.ebonPlaguebringer !== 0)
            writer.tag(85, WireType.Varint).int32(message.ebonPlaguebringer);
        /* bool scourge_strike = 86; */
        if (message.scourgeStrike !== false)
            writer.tag(86, WireType.Varint).bool(message.scourgeStrike);
        /* int32 rage_of_rivendare = 87; */
        if (message.rageOfRivendare !== 0)
            writer.tag(87, WireType.Varint).int32(message.rageOfRivendare);
        /* bool summon_gargoyle = 88; */
        if (message.summonGargoyle !== false)
            writer.tag(88, WireType.Varint).bool(message.summonGargoyle);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message proto.DeathKnightTalents
 */
export const DeathKnightTalents = new DeathKnightTalents$Type();
// @generated message type with reflection information, may provide speed optimized methods
class DeathKnight$Type extends MessageType {
    constructor() {
        super("proto.DeathKnight", [
            { no: 1, name: "rotation", kind: "message", T: () => DeathKnight_Rotation },
            { no: 2, name: "talents", kind: "message", T: () => DeathKnightTalents },
            { no: 3, name: "options", kind: "message", T: () => DeathKnight_Options }
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
                case /* proto.DeathKnight.Rotation rotation */ 1:
                    message.rotation = DeathKnight_Rotation.internalBinaryRead(reader, reader.uint32(), options, message.rotation);
                    break;
                case /* proto.DeathKnightTalents talents */ 2:
                    message.talents = DeathKnightTalents.internalBinaryRead(reader, reader.uint32(), options, message.talents);
                    break;
                case /* proto.DeathKnight.Options options */ 3:
                    message.options = DeathKnight_Options.internalBinaryRead(reader, reader.uint32(), options, message.options);
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
        /* proto.DeathKnight.Rotation rotation = 1; */
        if (message.rotation)
            DeathKnight_Rotation.internalBinaryWrite(message.rotation, writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        /* proto.DeathKnightTalents talents = 2; */
        if (message.talents)
            DeathKnightTalents.internalBinaryWrite(message.talents, writer.tag(2, WireType.LengthDelimited).fork(), options).join();
        /* proto.DeathKnight.Options options = 3; */
        if (message.options)
            DeathKnight_Options.internalBinaryWrite(message.options, writer.tag(3, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message proto.DeathKnight
 */
export const DeathKnight = new DeathKnight$Type();
// @generated message type with reflection information, may provide speed optimized methods
class DeathKnight_Rotation$Type extends MessageType {
    constructor() {
        super("proto.DeathKnight.Rotation", [
            { no: 1, name: "army_of_the_dead", kind: "enum", T: () => ["proto.DeathKnight.Rotation.ArmyOfTheDead", DeathKnight_Rotation_ArmyOfTheDead] },
            { no: 2, name: "use_death_and_decay", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 3, name: "unholy_presence_opener", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 4, name: "disease_refresh_duration", kind: "scalar", T: 1 /*ScalarType.DOUBLE*/ },
            { no: 5, name: "refresh_horn_of_winter", kind: "scalar", T: 8 /*ScalarType.BOOL*/ }
        ]);
    }
    create(value) {
        const message = { armyOfTheDead: 0, useDeathAndDecay: false, unholyPresenceOpener: false, diseaseRefreshDuration: 0, refreshHornOfWinter: false };
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
                case /* proto.DeathKnight.Rotation.ArmyOfTheDead army_of_the_dead */ 1:
                    message.armyOfTheDead = reader.int32();
                    break;
                case /* bool use_death_and_decay */ 2:
                    message.useDeathAndDecay = reader.bool();
                    break;
                case /* bool unholy_presence_opener */ 3:
                    message.unholyPresenceOpener = reader.bool();
                    break;
                case /* double disease_refresh_duration */ 4:
                    message.diseaseRefreshDuration = reader.double();
                    break;
                case /* bool refresh_horn_of_winter */ 5:
                    message.refreshHornOfWinter = reader.bool();
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
        /* proto.DeathKnight.Rotation.ArmyOfTheDead army_of_the_dead = 1; */
        if (message.armyOfTheDead !== 0)
            writer.tag(1, WireType.Varint).int32(message.armyOfTheDead);
        /* bool use_death_and_decay = 2; */
        if (message.useDeathAndDecay !== false)
            writer.tag(2, WireType.Varint).bool(message.useDeathAndDecay);
        /* bool unholy_presence_opener = 3; */
        if (message.unholyPresenceOpener !== false)
            writer.tag(3, WireType.Varint).bool(message.unholyPresenceOpener);
        /* double disease_refresh_duration = 4; */
        if (message.diseaseRefreshDuration !== 0)
            writer.tag(4, WireType.Bit64).double(message.diseaseRefreshDuration);
        /* bool refresh_horn_of_winter = 5; */
        if (message.refreshHornOfWinter !== false)
            writer.tag(5, WireType.Varint).bool(message.refreshHornOfWinter);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message proto.DeathKnight.Rotation
 */
export const DeathKnight_Rotation = new DeathKnight_Rotation$Type();
// @generated message type with reflection information, may provide speed optimized methods
class DeathKnight_Options$Type extends MessageType {
    constructor() {
        super("proto.DeathKnight.Options", [
            { no: 1, name: "starting_runic_power", kind: "scalar", T: 1 /*ScalarType.DOUBLE*/ },
            { no: 2, name: "pet_uptime", kind: "scalar", T: 1 /*ScalarType.DOUBLE*/ },
            { no: 3, name: "precast_ghoul_frenzy", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 4, name: "precast_horn_of_winter", kind: "scalar", T: 8 /*ScalarType.BOOL*/ }
        ]);
    }
    create(value) {
        const message = { startingRunicPower: 0, petUptime: 0, precastGhoulFrenzy: false, precastHornOfWinter: false };
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
                case /* double starting_runic_power */ 1:
                    message.startingRunicPower = reader.double();
                    break;
                case /* double pet_uptime */ 2:
                    message.petUptime = reader.double();
                    break;
                case /* bool precast_ghoul_frenzy */ 3:
                    message.precastGhoulFrenzy = reader.bool();
                    break;
                case /* bool precast_horn_of_winter */ 4:
                    message.precastHornOfWinter = reader.bool();
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
        /* double starting_runic_power = 1; */
        if (message.startingRunicPower !== 0)
            writer.tag(1, WireType.Bit64).double(message.startingRunicPower);
        /* double pet_uptime = 2; */
        if (message.petUptime !== 0)
            writer.tag(2, WireType.Bit64).double(message.petUptime);
        /* bool precast_ghoul_frenzy = 3; */
        if (message.precastGhoulFrenzy !== false)
            writer.tag(3, WireType.Varint).bool(message.precastGhoulFrenzy);
        /* bool precast_horn_of_winter = 4; */
        if (message.precastHornOfWinter !== false)
            writer.tag(4, WireType.Varint).bool(message.precastHornOfWinter);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message proto.DeathKnight.Options
 */
export const DeathKnight_Options = new DeathKnight_Options$Type();
// @generated message type with reflection information, may provide speed optimized methods
class DeathKnightTank$Type extends MessageType {
    constructor() {
        super("proto.DeathKnightTank", [
            { no: 1, name: "rotation", kind: "message", T: () => DeathKnightTank_Rotation },
            { no: 2, name: "talents", kind: "message", T: () => DeathKnightTalents },
            { no: 3, name: "options", kind: "message", T: () => DeathKnightTank_Options }
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
                case /* proto.DeathKnightTank.Rotation rotation */ 1:
                    message.rotation = DeathKnightTank_Rotation.internalBinaryRead(reader, reader.uint32(), options, message.rotation);
                    break;
                case /* proto.DeathKnightTalents talents */ 2:
                    message.talents = DeathKnightTalents.internalBinaryRead(reader, reader.uint32(), options, message.talents);
                    break;
                case /* proto.DeathKnightTank.Options options */ 3:
                    message.options = DeathKnightTank_Options.internalBinaryRead(reader, reader.uint32(), options, message.options);
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
        /* proto.DeathKnightTank.Rotation rotation = 1; */
        if (message.rotation)
            DeathKnightTank_Rotation.internalBinaryWrite(message.rotation, writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        /* proto.DeathKnightTalents talents = 2; */
        if (message.talents)
            DeathKnightTalents.internalBinaryWrite(message.talents, writer.tag(2, WireType.LengthDelimited).fork(), options).join();
        /* proto.DeathKnightTank.Options options = 3; */
        if (message.options)
            DeathKnightTank_Options.internalBinaryWrite(message.options, writer.tag(3, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message proto.DeathKnightTank
 */
export const DeathKnightTank = new DeathKnightTank$Type();
// @generated message type with reflection information, may provide speed optimized methods
class DeathKnightTank_Rotation$Type extends MessageType {
    constructor() {
        super("proto.DeathKnightTank.Rotation", []);
    }
    create(value) {
        const message = {};
        Object.defineProperty(message, MESSAGE_TYPE, { enumerable: false, value: this });
        if (value !== undefined)
            reflectionMergePartial(this, message, value);
        return message;
    }
    internalBinaryRead(reader, length, options, target) {
        return target ?? this.create();
    }
    internalBinaryWrite(message, writer, options) {
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message proto.DeathKnightTank.Rotation
 */
export const DeathKnightTank_Rotation = new DeathKnightTank_Rotation$Type();
// @generated message type with reflection information, may provide speed optimized methods
class DeathKnightTank_Options$Type extends MessageType {
    constructor() {
        super("proto.DeathKnightTank.Options", []);
    }
    create(value) {
        const message = {};
        Object.defineProperty(message, MESSAGE_TYPE, { enumerable: false, value: this });
        if (value !== undefined)
            reflectionMergePartial(this, message, value);
        return message;
    }
    internalBinaryRead(reader, length, options, target) {
        return target ?? this.create();
    }
    internalBinaryWrite(message, writer, options) {
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message proto.DeathKnightTank.Options
 */
export const DeathKnightTank_Options = new DeathKnightTank_Options$Type();
