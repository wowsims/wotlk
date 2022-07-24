import { UnknownFieldHandler } from '/wotlk/protobuf-ts/index.js';
import { WireType } from '/wotlk/protobuf-ts/index.js';
import { reflectionMergePartial } from '/wotlk/protobuf-ts/index.js';
import { MESSAGE_TYPE } from '/wotlk/protobuf-ts/index.js';
import { MessageType } from '/wotlk/protobuf-ts/index.js';
import { Stat } from './common.js';
import { Target } from './common.js';
import { Gem } from './common.js';
import { Enchant } from './common.js';
import { Item } from './common.js';
import { Encounter } from './common.js';
import { ActionID } from './common.js';
import { RaidTarget } from './common.js';
import { Debuffs } from './common.js';
import { RaidBuffs } from './common.js';
import { PartyBuffs } from './common.js';
import { HealingModel } from './common.js';
import { Cooldowns } from './common.js';
import { Profession } from './common.js';
import { Glyphs } from './common.js';
import { TankDeathknight } from './deathknight.js';
import { Deathknight } from './deathknight.js';
import { ProtectionWarrior } from './warrior.js';
import { Warrior } from './warrior.js';
import { Warlock } from './warlock.js';
import { EnhancementShaman } from './shaman.js';
import { ElementalShaman } from './shaman.js';
import { Rogue } from './rogue.js';
import { SmitePriest } from './priest.js';
import { ShadowPriest } from './priest.js';
import { ProtectionPaladin } from './paladin.js';
import { RetributionPaladin } from './paladin.js';
import { Mage } from './mage.js';
import { Hunter } from './hunter.js';
import { FeralTankDruid } from './druid.js';
import { FeralDruid } from './druid.js';
import { BalanceDruid } from './druid.js';
import { IndividualBuffs } from './common.js';
import { Consumes } from './common.js';
import { EquipmentSpec } from './common.js';
import { Class } from './common.js';
import { ShattrathFaction } from './common.js';
import { Race } from './common.js';
/**
 * @generated from protobuf enum proto.ResourceType
 */
export var ResourceType;
(function (ResourceType) {
    /**
     * @generated from protobuf enum value: ResourceTypeNone = 0;
     */
    ResourceType[ResourceType["ResourceTypeNone"] = 0] = "ResourceTypeNone";
    /**
     * @generated from protobuf enum value: ResourceTypeMana = 1;
     */
    ResourceType[ResourceType["ResourceTypeMana"] = 1] = "ResourceTypeMana";
    /**
     * @generated from protobuf enum value: ResourceTypeEnergy = 2;
     */
    ResourceType[ResourceType["ResourceTypeEnergy"] = 2] = "ResourceTypeEnergy";
    /**
     * @generated from protobuf enum value: ResourceTypeRage = 3;
     */
    ResourceType[ResourceType["ResourceTypeRage"] = 3] = "ResourceTypeRage";
    /**
     * @generated from protobuf enum value: ResourceTypeComboPoints = 4;
     */
    ResourceType[ResourceType["ResourceTypeComboPoints"] = 4] = "ResourceTypeComboPoints";
    /**
     * @generated from protobuf enum value: ResourceTypeFocus = 5;
     */
    ResourceType[ResourceType["ResourceTypeFocus"] = 5] = "ResourceTypeFocus";
    /**
     * @generated from protobuf enum value: ResourceTypeHealth = 6;
     */
    ResourceType[ResourceType["ResourceTypeHealth"] = 6] = "ResourceTypeHealth";
    /**
     * @generated from protobuf enum value: ResourceTypeRunicPower = 7;
     */
    ResourceType[ResourceType["ResourceTypeRunicPower"] = 7] = "ResourceTypeRunicPower";
    /**
     * @generated from protobuf enum value: ResourceTypeBloodRune = 8;
     */
    ResourceType[ResourceType["ResourceTypeBloodRune"] = 8] = "ResourceTypeBloodRune";
    /**
     * @generated from protobuf enum value: ResourceTypeFrostRune = 9;
     */
    ResourceType[ResourceType["ResourceTypeFrostRune"] = 9] = "ResourceTypeFrostRune";
    /**
     * @generated from protobuf enum value: ResourceTypeUnholyRune = 10;
     */
    ResourceType[ResourceType["ResourceTypeUnholyRune"] = 10] = "ResourceTypeUnholyRune";
    /**
     * @generated from protobuf enum value: ResourceTypeDeathRune = 11;
     */
    ResourceType[ResourceType["ResourceTypeDeathRune"] = 11] = "ResourceTypeDeathRune";
})(ResourceType || (ResourceType = {}));
// @generated message type with reflection information, may provide speed optimized methods
class Player$Type extends MessageType {
    constructor() {
        super("proto.Player", [
            { no: 16, name: "name", kind: "scalar", T: 9 /*ScalarType.STRING*/ },
            { no: 1, name: "race", kind: "enum", T: () => ["proto.Race", Race] },
            { no: 24, name: "shatt_faction", kind: "enum", T: () => ["proto.ShattrathFaction", ShattrathFaction] },
            { no: 2, name: "class", kind: "enum", T: () => ["proto.Class", Class] },
            { no: 3, name: "equipment", kind: "message", T: () => EquipmentSpec },
            { no: 4, name: "consumes", kind: "message", T: () => Consumes },
            { no: 5, name: "bonus_stats", kind: "scalar", repeat: 1 /*RepeatType.PACKED*/, T: 1 /*ScalarType.DOUBLE*/ },
            { no: 15, name: "buffs", kind: "message", T: () => IndividualBuffs },
            { no: 6, name: "balance_druid", kind: "message", oneof: "spec", T: () => BalanceDruid },
            { no: 22, name: "feral_druid", kind: "message", oneof: "spec", T: () => FeralDruid },
            { no: 26, name: "feral_tank_druid", kind: "message", oneof: "spec", T: () => FeralTankDruid },
            { no: 7, name: "hunter", kind: "message", oneof: "spec", T: () => Hunter },
            { no: 8, name: "mage", kind: "message", oneof: "spec", T: () => Mage },
            { no: 9, name: "retribution_paladin", kind: "message", oneof: "spec", T: () => RetributionPaladin },
            { no: 25, name: "protection_paladin", kind: "message", oneof: "spec", T: () => ProtectionPaladin },
            { no: 10, name: "shadow_priest", kind: "message", oneof: "spec", T: () => ShadowPriest },
            { no: 20, name: "smite_priest", kind: "message", oneof: "spec", T: () => SmitePriest },
            { no: 11, name: "rogue", kind: "message", oneof: "spec", T: () => Rogue },
            { no: 12, name: "elemental_shaman", kind: "message", oneof: "spec", T: () => ElementalShaman },
            { no: 18, name: "enhancement_shaman", kind: "message", oneof: "spec", T: () => EnhancementShaman },
            { no: 13, name: "warlock", kind: "message", oneof: "spec", T: () => Warlock },
            { no: 14, name: "warrior", kind: "message", oneof: "spec", T: () => Warrior },
            { no: 21, name: "protection_warrior", kind: "message", oneof: "spec", T: () => ProtectionWarrior },
            { no: 31, name: "deathknight", kind: "message", oneof: "spec", T: () => Deathknight },
            { no: 32, name: "tank_deathknight", kind: "message", oneof: "spec", T: () => TankDeathknight },
            { no: 17, name: "talentsString", kind: "scalar", T: 9 /*ScalarType.STRING*/ },
            { no: 28, name: "glyphs", kind: "message", T: () => Glyphs },
            { no: 29, name: "profession1", kind: "enum", T: () => ["proto.Profession", Profession] },
            { no: 30, name: "profession2", kind: "enum", T: () => ["proto.Profession", Profession] },
            { no: 19, name: "cooldowns", kind: "message", T: () => Cooldowns },
            { no: 23, name: "in_front_of_target", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 27, name: "healing_model", kind: "message", T: () => HealingModel }
        ]);
    }
    create(value) {
        const message = { name: "", race: 0, shattFaction: 0, class: 0, bonusStats: [], spec: { oneofKind: undefined }, talentsString: "", profession1: 0, profession2: 0, inFrontOfTarget: false };
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
                case /* string name */ 16:
                    message.name = reader.string();
                    break;
                case /* proto.Race race */ 1:
                    message.race = reader.int32();
                    break;
                case /* proto.ShattrathFaction shatt_faction */ 24:
                    message.shattFaction = reader.int32();
                    break;
                case /* proto.Class class */ 2:
                    message.class = reader.int32();
                    break;
                case /* proto.EquipmentSpec equipment */ 3:
                    message.equipment = EquipmentSpec.internalBinaryRead(reader, reader.uint32(), options, message.equipment);
                    break;
                case /* proto.Consumes consumes */ 4:
                    message.consumes = Consumes.internalBinaryRead(reader, reader.uint32(), options, message.consumes);
                    break;
                case /* repeated double bonus_stats */ 5:
                    if (wireType === WireType.LengthDelimited)
                        for (let e = reader.int32() + reader.pos; reader.pos < e;)
                            message.bonusStats.push(reader.double());
                    else
                        message.bonusStats.push(reader.double());
                    break;
                case /* proto.IndividualBuffs buffs */ 15:
                    message.buffs = IndividualBuffs.internalBinaryRead(reader, reader.uint32(), options, message.buffs);
                    break;
                case /* proto.BalanceDruid balance_druid */ 6:
                    message.spec = {
                        oneofKind: "balanceDruid",
                        balanceDruid: BalanceDruid.internalBinaryRead(reader, reader.uint32(), options, message.spec.balanceDruid)
                    };
                    break;
                case /* proto.FeralDruid feral_druid */ 22:
                    message.spec = {
                        oneofKind: "feralDruid",
                        feralDruid: FeralDruid.internalBinaryRead(reader, reader.uint32(), options, message.spec.feralDruid)
                    };
                    break;
                case /* proto.FeralTankDruid feral_tank_druid */ 26:
                    message.spec = {
                        oneofKind: "feralTankDruid",
                        feralTankDruid: FeralTankDruid.internalBinaryRead(reader, reader.uint32(), options, message.spec.feralTankDruid)
                    };
                    break;
                case /* proto.Hunter hunter */ 7:
                    message.spec = {
                        oneofKind: "hunter",
                        hunter: Hunter.internalBinaryRead(reader, reader.uint32(), options, message.spec.hunter)
                    };
                    break;
                case /* proto.Mage mage */ 8:
                    message.spec = {
                        oneofKind: "mage",
                        mage: Mage.internalBinaryRead(reader, reader.uint32(), options, message.spec.mage)
                    };
                    break;
                case /* proto.RetributionPaladin retribution_paladin */ 9:
                    message.spec = {
                        oneofKind: "retributionPaladin",
                        retributionPaladin: RetributionPaladin.internalBinaryRead(reader, reader.uint32(), options, message.spec.retributionPaladin)
                    };
                    break;
                case /* proto.ProtectionPaladin protection_paladin */ 25:
                    message.spec = {
                        oneofKind: "protectionPaladin",
                        protectionPaladin: ProtectionPaladin.internalBinaryRead(reader, reader.uint32(), options, message.spec.protectionPaladin)
                    };
                    break;
                case /* proto.ShadowPriest shadow_priest */ 10:
                    message.spec = {
                        oneofKind: "shadowPriest",
                        shadowPriest: ShadowPriest.internalBinaryRead(reader, reader.uint32(), options, message.spec.shadowPriest)
                    };
                    break;
                case /* proto.SmitePriest smite_priest */ 20:
                    message.spec = {
                        oneofKind: "smitePriest",
                        smitePriest: SmitePriest.internalBinaryRead(reader, reader.uint32(), options, message.spec.smitePriest)
                    };
                    break;
                case /* proto.Rogue rogue */ 11:
                    message.spec = {
                        oneofKind: "rogue",
                        rogue: Rogue.internalBinaryRead(reader, reader.uint32(), options, message.spec.rogue)
                    };
                    break;
                case /* proto.ElementalShaman elemental_shaman */ 12:
                    message.spec = {
                        oneofKind: "elementalShaman",
                        elementalShaman: ElementalShaman.internalBinaryRead(reader, reader.uint32(), options, message.spec.elementalShaman)
                    };
                    break;
                case /* proto.EnhancementShaman enhancement_shaman */ 18:
                    message.spec = {
                        oneofKind: "enhancementShaman",
                        enhancementShaman: EnhancementShaman.internalBinaryRead(reader, reader.uint32(), options, message.spec.enhancementShaman)
                    };
                    break;
                case /* proto.Warlock warlock */ 13:
                    message.spec = {
                        oneofKind: "warlock",
                        warlock: Warlock.internalBinaryRead(reader, reader.uint32(), options, message.spec.warlock)
                    };
                    break;
                case /* proto.Warrior warrior */ 14:
                    message.spec = {
                        oneofKind: "warrior",
                        warrior: Warrior.internalBinaryRead(reader, reader.uint32(), options, message.spec.warrior)
                    };
                    break;
                case /* proto.ProtectionWarrior protection_warrior */ 21:
                    message.spec = {
                        oneofKind: "protectionWarrior",
                        protectionWarrior: ProtectionWarrior.internalBinaryRead(reader, reader.uint32(), options, message.spec.protectionWarrior)
                    };
                    break;
                case /* proto.Deathknight deathknight */ 31:
                    message.spec = {
                        oneofKind: "deathknight",
                        deathknight: Deathknight.internalBinaryRead(reader, reader.uint32(), options, message.spec.deathknight)
                    };
                    break;
                case /* proto.TankDeathknight tank_deathknight */ 32:
                    message.spec = {
                        oneofKind: "tankDeathknight",
                        tankDeathknight: TankDeathknight.internalBinaryRead(reader, reader.uint32(), options, message.spec.tankDeathknight)
                    };
                    break;
                case /* string talentsString */ 17:
                    message.talentsString = reader.string();
                    break;
                case /* proto.Glyphs glyphs */ 28:
                    message.glyphs = Glyphs.internalBinaryRead(reader, reader.uint32(), options, message.glyphs);
                    break;
                case /* proto.Profession profession1 */ 29:
                    message.profession1 = reader.int32();
                    break;
                case /* proto.Profession profession2 */ 30:
                    message.profession2 = reader.int32();
                    break;
                case /* proto.Cooldowns cooldowns */ 19:
                    message.cooldowns = Cooldowns.internalBinaryRead(reader, reader.uint32(), options, message.cooldowns);
                    break;
                case /* bool in_front_of_target */ 23:
                    message.inFrontOfTarget = reader.bool();
                    break;
                case /* proto.HealingModel healing_model */ 27:
                    message.healingModel = HealingModel.internalBinaryRead(reader, reader.uint32(), options, message.healingModel);
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
        /* string name = 16; */
        if (message.name !== "")
            writer.tag(16, WireType.LengthDelimited).string(message.name);
        /* proto.Race race = 1; */
        if (message.race !== 0)
            writer.tag(1, WireType.Varint).int32(message.race);
        /* proto.ShattrathFaction shatt_faction = 24; */
        if (message.shattFaction !== 0)
            writer.tag(24, WireType.Varint).int32(message.shattFaction);
        /* proto.Class class = 2; */
        if (message.class !== 0)
            writer.tag(2, WireType.Varint).int32(message.class);
        /* proto.EquipmentSpec equipment = 3; */
        if (message.equipment)
            EquipmentSpec.internalBinaryWrite(message.equipment, writer.tag(3, WireType.LengthDelimited).fork(), options).join();
        /* proto.Consumes consumes = 4; */
        if (message.consumes)
            Consumes.internalBinaryWrite(message.consumes, writer.tag(4, WireType.LengthDelimited).fork(), options).join();
        /* repeated double bonus_stats = 5; */
        if (message.bonusStats.length) {
            writer.tag(5, WireType.LengthDelimited).fork();
            for (let i = 0; i < message.bonusStats.length; i++)
                writer.double(message.bonusStats[i]);
            writer.join();
        }
        /* proto.IndividualBuffs buffs = 15; */
        if (message.buffs)
            IndividualBuffs.internalBinaryWrite(message.buffs, writer.tag(15, WireType.LengthDelimited).fork(), options).join();
        /* proto.BalanceDruid balance_druid = 6; */
        if (message.spec.oneofKind === "balanceDruid")
            BalanceDruid.internalBinaryWrite(message.spec.balanceDruid, writer.tag(6, WireType.LengthDelimited).fork(), options).join();
        /* proto.FeralDruid feral_druid = 22; */
        if (message.spec.oneofKind === "feralDruid")
            FeralDruid.internalBinaryWrite(message.spec.feralDruid, writer.tag(22, WireType.LengthDelimited).fork(), options).join();
        /* proto.FeralTankDruid feral_tank_druid = 26; */
        if (message.spec.oneofKind === "feralTankDruid")
            FeralTankDruid.internalBinaryWrite(message.spec.feralTankDruid, writer.tag(26, WireType.LengthDelimited).fork(), options).join();
        /* proto.Hunter hunter = 7; */
        if (message.spec.oneofKind === "hunter")
            Hunter.internalBinaryWrite(message.spec.hunter, writer.tag(7, WireType.LengthDelimited).fork(), options).join();
        /* proto.Mage mage = 8; */
        if (message.spec.oneofKind === "mage")
            Mage.internalBinaryWrite(message.spec.mage, writer.tag(8, WireType.LengthDelimited).fork(), options).join();
        /* proto.RetributionPaladin retribution_paladin = 9; */
        if (message.spec.oneofKind === "retributionPaladin")
            RetributionPaladin.internalBinaryWrite(message.spec.retributionPaladin, writer.tag(9, WireType.LengthDelimited).fork(), options).join();
        /* proto.ProtectionPaladin protection_paladin = 25; */
        if (message.spec.oneofKind === "protectionPaladin")
            ProtectionPaladin.internalBinaryWrite(message.spec.protectionPaladin, writer.tag(25, WireType.LengthDelimited).fork(), options).join();
        /* proto.ShadowPriest shadow_priest = 10; */
        if (message.spec.oneofKind === "shadowPriest")
            ShadowPriest.internalBinaryWrite(message.spec.shadowPriest, writer.tag(10, WireType.LengthDelimited).fork(), options).join();
        /* proto.SmitePriest smite_priest = 20; */
        if (message.spec.oneofKind === "smitePriest")
            SmitePriest.internalBinaryWrite(message.spec.smitePriest, writer.tag(20, WireType.LengthDelimited).fork(), options).join();
        /* proto.Rogue rogue = 11; */
        if (message.spec.oneofKind === "rogue")
            Rogue.internalBinaryWrite(message.spec.rogue, writer.tag(11, WireType.LengthDelimited).fork(), options).join();
        /* proto.ElementalShaman elemental_shaman = 12; */
        if (message.spec.oneofKind === "elementalShaman")
            ElementalShaman.internalBinaryWrite(message.spec.elementalShaman, writer.tag(12, WireType.LengthDelimited).fork(), options).join();
        /* proto.EnhancementShaman enhancement_shaman = 18; */
        if (message.spec.oneofKind === "enhancementShaman")
            EnhancementShaman.internalBinaryWrite(message.spec.enhancementShaman, writer.tag(18, WireType.LengthDelimited).fork(), options).join();
        /* proto.Warlock warlock = 13; */
        if (message.spec.oneofKind === "warlock")
            Warlock.internalBinaryWrite(message.spec.warlock, writer.tag(13, WireType.LengthDelimited).fork(), options).join();
        /* proto.Warrior warrior = 14; */
        if (message.spec.oneofKind === "warrior")
            Warrior.internalBinaryWrite(message.spec.warrior, writer.tag(14, WireType.LengthDelimited).fork(), options).join();
        /* proto.ProtectionWarrior protection_warrior = 21; */
        if (message.spec.oneofKind === "protectionWarrior")
            ProtectionWarrior.internalBinaryWrite(message.spec.protectionWarrior, writer.tag(21, WireType.LengthDelimited).fork(), options).join();
        /* proto.Deathknight deathknight = 31; */
        if (message.spec.oneofKind === "deathknight")
            Deathknight.internalBinaryWrite(message.spec.deathknight, writer.tag(31, WireType.LengthDelimited).fork(), options).join();
        /* proto.TankDeathknight tank_deathknight = 32; */
        if (message.spec.oneofKind === "tankDeathknight")
            TankDeathknight.internalBinaryWrite(message.spec.tankDeathknight, writer.tag(32, WireType.LengthDelimited).fork(), options).join();
        /* string talentsString = 17; */
        if (message.talentsString !== "")
            writer.tag(17, WireType.LengthDelimited).string(message.talentsString);
        /* proto.Glyphs glyphs = 28; */
        if (message.glyphs)
            Glyphs.internalBinaryWrite(message.glyphs, writer.tag(28, WireType.LengthDelimited).fork(), options).join();
        /* proto.Profession profession1 = 29; */
        if (message.profession1 !== 0)
            writer.tag(29, WireType.Varint).int32(message.profession1);
        /* proto.Profession profession2 = 30; */
        if (message.profession2 !== 0)
            writer.tag(30, WireType.Varint).int32(message.profession2);
        /* proto.Cooldowns cooldowns = 19; */
        if (message.cooldowns)
            Cooldowns.internalBinaryWrite(message.cooldowns, writer.tag(19, WireType.LengthDelimited).fork(), options).join();
        /* bool in_front_of_target = 23; */
        if (message.inFrontOfTarget !== false)
            writer.tag(23, WireType.Varint).bool(message.inFrontOfTarget);
        /* proto.HealingModel healing_model = 27; */
        if (message.healingModel)
            HealingModel.internalBinaryWrite(message.healingModel, writer.tag(27, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message proto.Player
 */
export const Player = new Player$Type();
// @generated message type with reflection information, may provide speed optimized methods
class Party$Type extends MessageType {
    constructor() {
        super("proto.Party", [
            { no: 1, name: "players", kind: "message", repeat: 1 /*RepeatType.PACKED*/, T: () => Player },
            { no: 2, name: "buffs", kind: "message", T: () => PartyBuffs }
        ]);
    }
    create(value) {
        const message = { players: [] };
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
                case /* repeated proto.Player players */ 1:
                    message.players.push(Player.internalBinaryRead(reader, reader.uint32(), options));
                    break;
                case /* proto.PartyBuffs buffs */ 2:
                    message.buffs = PartyBuffs.internalBinaryRead(reader, reader.uint32(), options, message.buffs);
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
        /* repeated proto.Player players = 1; */
        for (let i = 0; i < message.players.length; i++)
            Player.internalBinaryWrite(message.players[i], writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        /* proto.PartyBuffs buffs = 2; */
        if (message.buffs)
            PartyBuffs.internalBinaryWrite(message.buffs, writer.tag(2, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message proto.Party
 */
export const Party = new Party$Type();
// @generated message type with reflection information, may provide speed optimized methods
class Raid$Type extends MessageType {
    constructor() {
        super("proto.Raid", [
            { no: 1, name: "parties", kind: "message", repeat: 1 /*RepeatType.PACKED*/, T: () => Party },
            { no: 2, name: "buffs", kind: "message", T: () => RaidBuffs },
            { no: 5, name: "debuffs", kind: "message", T: () => Debuffs },
            { no: 4, name: "tanks", kind: "message", repeat: 1 /*RepeatType.PACKED*/, T: () => RaidTarget },
            { no: 3, name: "stagger_stormstrikes", kind: "scalar", T: 8 /*ScalarType.BOOL*/ }
        ]);
    }
    create(value) {
        const message = { parties: [], tanks: [], staggerStormstrikes: false };
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
                case /* repeated proto.Party parties */ 1:
                    message.parties.push(Party.internalBinaryRead(reader, reader.uint32(), options));
                    break;
                case /* proto.RaidBuffs buffs */ 2:
                    message.buffs = RaidBuffs.internalBinaryRead(reader, reader.uint32(), options, message.buffs);
                    break;
                case /* proto.Debuffs debuffs */ 5:
                    message.debuffs = Debuffs.internalBinaryRead(reader, reader.uint32(), options, message.debuffs);
                    break;
                case /* repeated proto.RaidTarget tanks */ 4:
                    message.tanks.push(RaidTarget.internalBinaryRead(reader, reader.uint32(), options));
                    break;
                case /* bool stagger_stormstrikes */ 3:
                    message.staggerStormstrikes = reader.bool();
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
        /* repeated proto.Party parties = 1; */
        for (let i = 0; i < message.parties.length; i++)
            Party.internalBinaryWrite(message.parties[i], writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        /* proto.RaidBuffs buffs = 2; */
        if (message.buffs)
            RaidBuffs.internalBinaryWrite(message.buffs, writer.tag(2, WireType.LengthDelimited).fork(), options).join();
        /* proto.Debuffs debuffs = 5; */
        if (message.debuffs)
            Debuffs.internalBinaryWrite(message.debuffs, writer.tag(5, WireType.LengthDelimited).fork(), options).join();
        /* repeated proto.RaidTarget tanks = 4; */
        for (let i = 0; i < message.tanks.length; i++)
            RaidTarget.internalBinaryWrite(message.tanks[i], writer.tag(4, WireType.LengthDelimited).fork(), options).join();
        /* bool stagger_stormstrikes = 3; */
        if (message.staggerStormstrikes !== false)
            writer.tag(3, WireType.Varint).bool(message.staggerStormstrikes);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message proto.Raid
 */
export const Raid = new Raid$Type();
// @generated message type with reflection information, may provide speed optimized methods
class SimOptions$Type extends MessageType {
    constructor() {
        super("proto.SimOptions", [
            { no: 1, name: "iterations", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 2, name: "random_seed", kind: "scalar", T: 3 /*ScalarType.INT64*/, L: 0 /*LongType.BIGINT*/ },
            { no: 3, name: "debug", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 6, name: "debug_first_iteration", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 5, name: "is_test", kind: "scalar", T: 8 /*ScalarType.BOOL*/ }
        ]);
    }
    create(value) {
        const message = { iterations: 0, randomSeed: 0n, debug: false, debugFirstIteration: false, isTest: false };
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
                case /* int32 iterations */ 1:
                    message.iterations = reader.int32();
                    break;
                case /* int64 random_seed */ 2:
                    message.randomSeed = reader.int64().toBigInt();
                    break;
                case /* bool debug */ 3:
                    message.debug = reader.bool();
                    break;
                case /* bool debug_first_iteration */ 6:
                    message.debugFirstIteration = reader.bool();
                    break;
                case /* bool is_test */ 5:
                    message.isTest = reader.bool();
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
        /* int32 iterations = 1; */
        if (message.iterations !== 0)
            writer.tag(1, WireType.Varint).int32(message.iterations);
        /* int64 random_seed = 2; */
        if (message.randomSeed !== 0n)
            writer.tag(2, WireType.Varint).int64(message.randomSeed);
        /* bool debug = 3; */
        if (message.debug !== false)
            writer.tag(3, WireType.Varint).bool(message.debug);
        /* bool debug_first_iteration = 6; */
        if (message.debugFirstIteration !== false)
            writer.tag(6, WireType.Varint).bool(message.debugFirstIteration);
        /* bool is_test = 5; */
        if (message.isTest !== false)
            writer.tag(5, WireType.Varint).bool(message.isTest);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message proto.SimOptions
 */
export const SimOptions = new SimOptions$Type();
// @generated message type with reflection information, may provide speed optimized methods
class ActionMetrics$Type extends MessageType {
    constructor() {
        super("proto.ActionMetrics", [
            { no: 1, name: "id", kind: "message", T: () => ActionID },
            { no: 2, name: "is_melee", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 3, name: "targets", kind: "message", repeat: 1 /*RepeatType.PACKED*/, T: () => TargetedActionMetrics }
        ]);
    }
    create(value) {
        const message = { isMelee: false, targets: [] };
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
                case /* bool is_melee */ 2:
                    message.isMelee = reader.bool();
                    break;
                case /* repeated proto.TargetedActionMetrics targets */ 3:
                    message.targets.push(TargetedActionMetrics.internalBinaryRead(reader, reader.uint32(), options));
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
        /* bool is_melee = 2; */
        if (message.isMelee !== false)
            writer.tag(2, WireType.Varint).bool(message.isMelee);
        /* repeated proto.TargetedActionMetrics targets = 3; */
        for (let i = 0; i < message.targets.length; i++)
            TargetedActionMetrics.internalBinaryWrite(message.targets[i], writer.tag(3, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message proto.ActionMetrics
 */
export const ActionMetrics = new ActionMetrics$Type();
// @generated message type with reflection information, may provide speed optimized methods
class TargetedActionMetrics$Type extends MessageType {
    constructor() {
        super("proto.TargetedActionMetrics", [
            { no: 12, name: "unit_index", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 1, name: "casts", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 2, name: "hits", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 3, name: "crits", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 4, name: "misses", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 5, name: "dodges", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 6, name: "parries", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 7, name: "blocks", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 8, name: "glances", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 9, name: "damage", kind: "scalar", T: 1 /*ScalarType.DOUBLE*/ },
            { no: 10, name: "threat", kind: "scalar", T: 1 /*ScalarType.DOUBLE*/ }
        ]);
    }
    create(value) {
        const message = { unitIndex: 0, casts: 0, hits: 0, crits: 0, misses: 0, dodges: 0, parries: 0, blocks: 0, glances: 0, damage: 0, threat: 0 };
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
                case /* int32 unit_index */ 12:
                    message.unitIndex = reader.int32();
                    break;
                case /* int32 casts */ 1:
                    message.casts = reader.int32();
                    break;
                case /* int32 hits */ 2:
                    message.hits = reader.int32();
                    break;
                case /* int32 crits */ 3:
                    message.crits = reader.int32();
                    break;
                case /* int32 misses */ 4:
                    message.misses = reader.int32();
                    break;
                case /* int32 dodges */ 5:
                    message.dodges = reader.int32();
                    break;
                case /* int32 parries */ 6:
                    message.parries = reader.int32();
                    break;
                case /* int32 blocks */ 7:
                    message.blocks = reader.int32();
                    break;
                case /* int32 glances */ 8:
                    message.glances = reader.int32();
                    break;
                case /* double damage */ 9:
                    message.damage = reader.double();
                    break;
                case /* double threat */ 10:
                    message.threat = reader.double();
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
        /* int32 unit_index = 12; */
        if (message.unitIndex !== 0)
            writer.tag(12, WireType.Varint).int32(message.unitIndex);
        /* int32 casts = 1; */
        if (message.casts !== 0)
            writer.tag(1, WireType.Varint).int32(message.casts);
        /* int32 hits = 2; */
        if (message.hits !== 0)
            writer.tag(2, WireType.Varint).int32(message.hits);
        /* int32 crits = 3; */
        if (message.crits !== 0)
            writer.tag(3, WireType.Varint).int32(message.crits);
        /* int32 misses = 4; */
        if (message.misses !== 0)
            writer.tag(4, WireType.Varint).int32(message.misses);
        /* int32 dodges = 5; */
        if (message.dodges !== 0)
            writer.tag(5, WireType.Varint).int32(message.dodges);
        /* int32 parries = 6; */
        if (message.parries !== 0)
            writer.tag(6, WireType.Varint).int32(message.parries);
        /* int32 blocks = 7; */
        if (message.blocks !== 0)
            writer.tag(7, WireType.Varint).int32(message.blocks);
        /* int32 glances = 8; */
        if (message.glances !== 0)
            writer.tag(8, WireType.Varint).int32(message.glances);
        /* double damage = 9; */
        if (message.damage !== 0)
            writer.tag(9, WireType.Bit64).double(message.damage);
        /* double threat = 10; */
        if (message.threat !== 0)
            writer.tag(10, WireType.Bit64).double(message.threat);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message proto.TargetedActionMetrics
 */
export const TargetedActionMetrics = new TargetedActionMetrics$Type();
// @generated message type with reflection information, may provide speed optimized methods
class AuraMetrics$Type extends MessageType {
    constructor() {
        super("proto.AuraMetrics", [
            { no: 1, name: "id", kind: "message", T: () => ActionID },
            { no: 2, name: "uptime_seconds_avg", kind: "scalar", T: 1 /*ScalarType.DOUBLE*/ },
            { no: 3, name: "uptime_seconds_stdev", kind: "scalar", T: 1 /*ScalarType.DOUBLE*/ }
        ]);
    }
    create(value) {
        const message = { uptimeSecondsAvg: 0, uptimeSecondsStdev: 0 };
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
                case /* double uptime_seconds_avg */ 2:
                    message.uptimeSecondsAvg = reader.double();
                    break;
                case /* double uptime_seconds_stdev */ 3:
                    message.uptimeSecondsStdev = reader.double();
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
        /* double uptime_seconds_avg = 2; */
        if (message.uptimeSecondsAvg !== 0)
            writer.tag(2, WireType.Bit64).double(message.uptimeSecondsAvg);
        /* double uptime_seconds_stdev = 3; */
        if (message.uptimeSecondsStdev !== 0)
            writer.tag(3, WireType.Bit64).double(message.uptimeSecondsStdev);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message proto.AuraMetrics
 */
export const AuraMetrics = new AuraMetrics$Type();
// @generated message type with reflection information, may provide speed optimized methods
class ResourceMetrics$Type extends MessageType {
    constructor() {
        super("proto.ResourceMetrics", [
            { no: 1, name: "id", kind: "message", T: () => ActionID },
            { no: 2, name: "type", kind: "enum", T: () => ["proto.ResourceType", ResourceType] },
            { no: 3, name: "events", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 4, name: "gain", kind: "scalar", T: 1 /*ScalarType.DOUBLE*/ },
            { no: 5, name: "actual_gain", kind: "scalar", T: 1 /*ScalarType.DOUBLE*/ }
        ]);
    }
    create(value) {
        const message = { type: 0, events: 0, gain: 0, actualGain: 0 };
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
                case /* proto.ResourceType type */ 2:
                    message.type = reader.int32();
                    break;
                case /* int32 events */ 3:
                    message.events = reader.int32();
                    break;
                case /* double gain */ 4:
                    message.gain = reader.double();
                    break;
                case /* double actual_gain */ 5:
                    message.actualGain = reader.double();
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
        /* proto.ResourceType type = 2; */
        if (message.type !== 0)
            writer.tag(2, WireType.Varint).int32(message.type);
        /* int32 events = 3; */
        if (message.events !== 0)
            writer.tag(3, WireType.Varint).int32(message.events);
        /* double gain = 4; */
        if (message.gain !== 0)
            writer.tag(4, WireType.Bit64).double(message.gain);
        /* double actual_gain = 5; */
        if (message.actualGain !== 0)
            writer.tag(5, WireType.Bit64).double(message.actualGain);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message proto.ResourceMetrics
 */
export const ResourceMetrics = new ResourceMetrics$Type();
// @generated message type with reflection information, may provide speed optimized methods
class DistributionMetrics$Type extends MessageType {
    constructor() {
        super("proto.DistributionMetrics", [
            { no: 1, name: "avg", kind: "scalar", T: 1 /*ScalarType.DOUBLE*/ },
            { no: 2, name: "stdev", kind: "scalar", T: 1 /*ScalarType.DOUBLE*/ },
            { no: 3, name: "max", kind: "scalar", T: 1 /*ScalarType.DOUBLE*/ },
            { no: 4, name: "hist", kind: "map", K: 5 /*ScalarType.INT32*/, V: { kind: "scalar", T: 5 /*ScalarType.INT32*/ } }
        ]);
    }
    create(value) {
        const message = { avg: 0, stdev: 0, max: 0, hist: {} };
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
                case /* double avg */ 1:
                    message.avg = reader.double();
                    break;
                case /* double stdev */ 2:
                    message.stdev = reader.double();
                    break;
                case /* double max */ 3:
                    message.max = reader.double();
                    break;
                case /* map<int32, int32> hist */ 4:
                    this.binaryReadMap4(message.hist, reader, options);
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
    binaryReadMap4(map, reader, options) {
        let len = reader.uint32(), end = reader.pos + len, key, val;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case 1:
                    key = reader.int32();
                    break;
                case 2:
                    val = reader.int32();
                    break;
                default: throw new globalThis.Error("unknown map entry field for field proto.DistributionMetrics.hist");
            }
        }
        map[key ?? 0] = val ?? 0;
    }
    internalBinaryWrite(message, writer, options) {
        /* double avg = 1; */
        if (message.avg !== 0)
            writer.tag(1, WireType.Bit64).double(message.avg);
        /* double stdev = 2; */
        if (message.stdev !== 0)
            writer.tag(2, WireType.Bit64).double(message.stdev);
        /* double max = 3; */
        if (message.max !== 0)
            writer.tag(3, WireType.Bit64).double(message.max);
        /* map<int32, int32> hist = 4; */
        for (let k of Object.keys(message.hist))
            writer.tag(4, WireType.LengthDelimited).fork().tag(1, WireType.Varint).int32(parseInt(k)).tag(2, WireType.Varint).int32(message.hist[k]).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message proto.DistributionMetrics
 */
export const DistributionMetrics = new DistributionMetrics$Type();
// @generated message type with reflection information, may provide speed optimized methods
class UnitMetrics$Type extends MessageType {
    constructor() {
        super("proto.UnitMetrics", [
            { no: 9, name: "name", kind: "scalar", T: 9 /*ScalarType.STRING*/ },
            { no: 1, name: "dps", kind: "message", T: () => DistributionMetrics },
            { no: 8, name: "threat", kind: "message", T: () => DistributionMetrics },
            { no: 11, name: "dtps", kind: "message", T: () => DistributionMetrics },
            { no: 3, name: "seconds_oom_avg", kind: "scalar", T: 1 /*ScalarType.DOUBLE*/ },
            { no: 12, name: "chance_of_death", kind: "scalar", T: 1 /*ScalarType.DOUBLE*/ },
            { no: 5, name: "actions", kind: "message", repeat: 1 /*RepeatType.PACKED*/, T: () => ActionMetrics },
            { no: 6, name: "auras", kind: "message", repeat: 1 /*RepeatType.PACKED*/, T: () => AuraMetrics },
            { no: 10, name: "resources", kind: "message", repeat: 1 /*RepeatType.PACKED*/, T: () => ResourceMetrics },
            { no: 7, name: "pets", kind: "message", repeat: 1 /*RepeatType.PACKED*/, T: () => UnitMetrics }
        ]);
    }
    create(value) {
        const message = { name: "", secondsOomAvg: 0, chanceOfDeath: 0, actions: [], auras: [], resources: [], pets: [] };
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
                case /* string name */ 9:
                    message.name = reader.string();
                    break;
                case /* proto.DistributionMetrics dps */ 1:
                    message.dps = DistributionMetrics.internalBinaryRead(reader, reader.uint32(), options, message.dps);
                    break;
                case /* proto.DistributionMetrics threat */ 8:
                    message.threat = DistributionMetrics.internalBinaryRead(reader, reader.uint32(), options, message.threat);
                    break;
                case /* proto.DistributionMetrics dtps */ 11:
                    message.dtps = DistributionMetrics.internalBinaryRead(reader, reader.uint32(), options, message.dtps);
                    break;
                case /* double seconds_oom_avg */ 3:
                    message.secondsOomAvg = reader.double();
                    break;
                case /* double chance_of_death */ 12:
                    message.chanceOfDeath = reader.double();
                    break;
                case /* repeated proto.ActionMetrics actions */ 5:
                    message.actions.push(ActionMetrics.internalBinaryRead(reader, reader.uint32(), options));
                    break;
                case /* repeated proto.AuraMetrics auras */ 6:
                    message.auras.push(AuraMetrics.internalBinaryRead(reader, reader.uint32(), options));
                    break;
                case /* repeated proto.ResourceMetrics resources */ 10:
                    message.resources.push(ResourceMetrics.internalBinaryRead(reader, reader.uint32(), options));
                    break;
                case /* repeated proto.UnitMetrics pets */ 7:
                    message.pets.push(UnitMetrics.internalBinaryRead(reader, reader.uint32(), options));
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
        /* string name = 9; */
        if (message.name !== "")
            writer.tag(9, WireType.LengthDelimited).string(message.name);
        /* proto.DistributionMetrics dps = 1; */
        if (message.dps)
            DistributionMetrics.internalBinaryWrite(message.dps, writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        /* proto.DistributionMetrics threat = 8; */
        if (message.threat)
            DistributionMetrics.internalBinaryWrite(message.threat, writer.tag(8, WireType.LengthDelimited).fork(), options).join();
        /* proto.DistributionMetrics dtps = 11; */
        if (message.dtps)
            DistributionMetrics.internalBinaryWrite(message.dtps, writer.tag(11, WireType.LengthDelimited).fork(), options).join();
        /* double seconds_oom_avg = 3; */
        if (message.secondsOomAvg !== 0)
            writer.tag(3, WireType.Bit64).double(message.secondsOomAvg);
        /* double chance_of_death = 12; */
        if (message.chanceOfDeath !== 0)
            writer.tag(12, WireType.Bit64).double(message.chanceOfDeath);
        /* repeated proto.ActionMetrics actions = 5; */
        for (let i = 0; i < message.actions.length; i++)
            ActionMetrics.internalBinaryWrite(message.actions[i], writer.tag(5, WireType.LengthDelimited).fork(), options).join();
        /* repeated proto.AuraMetrics auras = 6; */
        for (let i = 0; i < message.auras.length; i++)
            AuraMetrics.internalBinaryWrite(message.auras[i], writer.tag(6, WireType.LengthDelimited).fork(), options).join();
        /* repeated proto.ResourceMetrics resources = 10; */
        for (let i = 0; i < message.resources.length; i++)
            ResourceMetrics.internalBinaryWrite(message.resources[i], writer.tag(10, WireType.LengthDelimited).fork(), options).join();
        /* repeated proto.UnitMetrics pets = 7; */
        for (let i = 0; i < message.pets.length; i++)
            UnitMetrics.internalBinaryWrite(message.pets[i], writer.tag(7, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message proto.UnitMetrics
 */
export const UnitMetrics = new UnitMetrics$Type();
// @generated message type with reflection information, may provide speed optimized methods
class PartyMetrics$Type extends MessageType {
    constructor() {
        super("proto.PartyMetrics", [
            { no: 1, name: "dps", kind: "message", T: () => DistributionMetrics },
            { no: 2, name: "players", kind: "message", repeat: 1 /*RepeatType.PACKED*/, T: () => UnitMetrics }
        ]);
    }
    create(value) {
        const message = { players: [] };
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
                case /* proto.DistributionMetrics dps */ 1:
                    message.dps = DistributionMetrics.internalBinaryRead(reader, reader.uint32(), options, message.dps);
                    break;
                case /* repeated proto.UnitMetrics players */ 2:
                    message.players.push(UnitMetrics.internalBinaryRead(reader, reader.uint32(), options));
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
        /* proto.DistributionMetrics dps = 1; */
        if (message.dps)
            DistributionMetrics.internalBinaryWrite(message.dps, writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        /* repeated proto.UnitMetrics players = 2; */
        for (let i = 0; i < message.players.length; i++)
            UnitMetrics.internalBinaryWrite(message.players[i], writer.tag(2, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message proto.PartyMetrics
 */
export const PartyMetrics = new PartyMetrics$Type();
// @generated message type with reflection information, may provide speed optimized methods
class RaidMetrics$Type extends MessageType {
    constructor() {
        super("proto.RaidMetrics", [
            { no: 1, name: "dps", kind: "message", T: () => DistributionMetrics },
            { no: 2, name: "parties", kind: "message", repeat: 1 /*RepeatType.PACKED*/, T: () => PartyMetrics }
        ]);
    }
    create(value) {
        const message = { parties: [] };
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
                case /* proto.DistributionMetrics dps */ 1:
                    message.dps = DistributionMetrics.internalBinaryRead(reader, reader.uint32(), options, message.dps);
                    break;
                case /* repeated proto.PartyMetrics parties */ 2:
                    message.parties.push(PartyMetrics.internalBinaryRead(reader, reader.uint32(), options));
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
        /* proto.DistributionMetrics dps = 1; */
        if (message.dps)
            DistributionMetrics.internalBinaryWrite(message.dps, writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        /* repeated proto.PartyMetrics parties = 2; */
        for (let i = 0; i < message.parties.length; i++)
            PartyMetrics.internalBinaryWrite(message.parties[i], writer.tag(2, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message proto.RaidMetrics
 */
export const RaidMetrics = new RaidMetrics$Type();
// @generated message type with reflection information, may provide speed optimized methods
class EncounterMetrics$Type extends MessageType {
    constructor() {
        super("proto.EncounterMetrics", [
            { no: 1, name: "targets", kind: "message", repeat: 1 /*RepeatType.PACKED*/, T: () => UnitMetrics }
        ]);
    }
    create(value) {
        const message = { targets: [] };
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
                case /* repeated proto.UnitMetrics targets */ 1:
                    message.targets.push(UnitMetrics.internalBinaryRead(reader, reader.uint32(), options));
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
        /* repeated proto.UnitMetrics targets = 1; */
        for (let i = 0; i < message.targets.length; i++)
            UnitMetrics.internalBinaryWrite(message.targets[i], writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message proto.EncounterMetrics
 */
export const EncounterMetrics = new EncounterMetrics$Type();
// @generated message type with reflection information, may provide speed optimized methods
class RaidSimRequest$Type extends MessageType {
    constructor() {
        super("proto.RaidSimRequest", [
            { no: 1, name: "raid", kind: "message", T: () => Raid },
            { no: 2, name: "encounter", kind: "message", T: () => Encounter },
            { no: 3, name: "sim_options", kind: "message", T: () => SimOptions }
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
                case /* proto.Raid raid */ 1:
                    message.raid = Raid.internalBinaryRead(reader, reader.uint32(), options, message.raid);
                    break;
                case /* proto.Encounter encounter */ 2:
                    message.encounter = Encounter.internalBinaryRead(reader, reader.uint32(), options, message.encounter);
                    break;
                case /* proto.SimOptions sim_options */ 3:
                    message.simOptions = SimOptions.internalBinaryRead(reader, reader.uint32(), options, message.simOptions);
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
        /* proto.Raid raid = 1; */
        if (message.raid)
            Raid.internalBinaryWrite(message.raid, writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        /* proto.Encounter encounter = 2; */
        if (message.encounter)
            Encounter.internalBinaryWrite(message.encounter, writer.tag(2, WireType.LengthDelimited).fork(), options).join();
        /* proto.SimOptions sim_options = 3; */
        if (message.simOptions)
            SimOptions.internalBinaryWrite(message.simOptions, writer.tag(3, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message proto.RaidSimRequest
 */
export const RaidSimRequest = new RaidSimRequest$Type();
// @generated message type with reflection information, may provide speed optimized methods
class RaidSimResult$Type extends MessageType {
    constructor() {
        super("proto.RaidSimResult", [
            { no: 1, name: "raid_metrics", kind: "message", T: () => RaidMetrics },
            { no: 2, name: "encounter_metrics", kind: "message", T: () => EncounterMetrics },
            { no: 3, name: "logs", kind: "scalar", T: 9 /*ScalarType.STRING*/ },
            { no: 4, name: "first_iteration_duration", kind: "scalar", T: 1 /*ScalarType.DOUBLE*/ },
            { no: 6, name: "avg_iteration_duration", kind: "scalar", T: 1 /*ScalarType.DOUBLE*/ },
            { no: 5, name: "error_result", kind: "scalar", T: 9 /*ScalarType.STRING*/ }
        ]);
    }
    create(value) {
        const message = { logs: "", firstIterationDuration: 0, avgIterationDuration: 0, errorResult: "" };
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
                case /* proto.RaidMetrics raid_metrics */ 1:
                    message.raidMetrics = RaidMetrics.internalBinaryRead(reader, reader.uint32(), options, message.raidMetrics);
                    break;
                case /* proto.EncounterMetrics encounter_metrics */ 2:
                    message.encounterMetrics = EncounterMetrics.internalBinaryRead(reader, reader.uint32(), options, message.encounterMetrics);
                    break;
                case /* string logs */ 3:
                    message.logs = reader.string();
                    break;
                case /* double first_iteration_duration */ 4:
                    message.firstIterationDuration = reader.double();
                    break;
                case /* double avg_iteration_duration */ 6:
                    message.avgIterationDuration = reader.double();
                    break;
                case /* string error_result */ 5:
                    message.errorResult = reader.string();
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
        /* proto.RaidMetrics raid_metrics = 1; */
        if (message.raidMetrics)
            RaidMetrics.internalBinaryWrite(message.raidMetrics, writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        /* proto.EncounterMetrics encounter_metrics = 2; */
        if (message.encounterMetrics)
            EncounterMetrics.internalBinaryWrite(message.encounterMetrics, writer.tag(2, WireType.LengthDelimited).fork(), options).join();
        /* string logs = 3; */
        if (message.logs !== "")
            writer.tag(3, WireType.LengthDelimited).string(message.logs);
        /* double first_iteration_duration = 4; */
        if (message.firstIterationDuration !== 0)
            writer.tag(4, WireType.Bit64).double(message.firstIterationDuration);
        /* double avg_iteration_duration = 6; */
        if (message.avgIterationDuration !== 0)
            writer.tag(6, WireType.Bit64).double(message.avgIterationDuration);
        /* string error_result = 5; */
        if (message.errorResult !== "")
            writer.tag(5, WireType.LengthDelimited).string(message.errorResult);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message proto.RaidSimResult
 */
export const RaidSimResult = new RaidSimResult$Type();
// @generated message type with reflection information, may provide speed optimized methods
class GearListRequest$Type extends MessageType {
    constructor() {
        super("proto.GearListRequest", []);
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
 * @generated MessageType for protobuf message proto.GearListRequest
 */
export const GearListRequest = new GearListRequest$Type();
// @generated message type with reflection information, may provide speed optimized methods
class GearListResult$Type extends MessageType {
    constructor() {
        super("proto.GearListResult", [
            { no: 1, name: "items", kind: "message", repeat: 1 /*RepeatType.PACKED*/, T: () => Item },
            { no: 2, name: "enchants", kind: "message", repeat: 1 /*RepeatType.PACKED*/, T: () => Enchant },
            { no: 3, name: "gems", kind: "message", repeat: 1 /*RepeatType.PACKED*/, T: () => Gem },
            { no: 4, name: "encounters", kind: "message", repeat: 1 /*RepeatType.PACKED*/, T: () => PresetEncounter }
        ]);
    }
    create(value) {
        const message = { items: [], enchants: [], gems: [], encounters: [] };
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
                case /* repeated proto.Item items */ 1:
                    message.items.push(Item.internalBinaryRead(reader, reader.uint32(), options));
                    break;
                case /* repeated proto.Enchant enchants */ 2:
                    message.enchants.push(Enchant.internalBinaryRead(reader, reader.uint32(), options));
                    break;
                case /* repeated proto.Gem gems */ 3:
                    message.gems.push(Gem.internalBinaryRead(reader, reader.uint32(), options));
                    break;
                case /* repeated proto.PresetEncounter encounters */ 4:
                    message.encounters.push(PresetEncounter.internalBinaryRead(reader, reader.uint32(), options));
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
        /* repeated proto.Item items = 1; */
        for (let i = 0; i < message.items.length; i++)
            Item.internalBinaryWrite(message.items[i], writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        /* repeated proto.Enchant enchants = 2; */
        for (let i = 0; i < message.enchants.length; i++)
            Enchant.internalBinaryWrite(message.enchants[i], writer.tag(2, WireType.LengthDelimited).fork(), options).join();
        /* repeated proto.Gem gems = 3; */
        for (let i = 0; i < message.gems.length; i++)
            Gem.internalBinaryWrite(message.gems[i], writer.tag(3, WireType.LengthDelimited).fork(), options).join();
        /* repeated proto.PresetEncounter encounters = 4; */
        for (let i = 0; i < message.encounters.length; i++)
            PresetEncounter.internalBinaryWrite(message.encounters[i], writer.tag(4, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message proto.GearListResult
 */
export const GearListResult = new GearListResult$Type();
// @generated message type with reflection information, may provide speed optimized methods
class PresetTarget$Type extends MessageType {
    constructor() {
        super("proto.PresetTarget", [
            { no: 1, name: "path", kind: "scalar", T: 9 /*ScalarType.STRING*/ },
            { no: 2, name: "target", kind: "message", T: () => Target }
        ]);
    }
    create(value) {
        const message = { path: "" };
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
                case /* string path */ 1:
                    message.path = reader.string();
                    break;
                case /* proto.Target target */ 2:
                    message.target = Target.internalBinaryRead(reader, reader.uint32(), options, message.target);
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
        /* string path = 1; */
        if (message.path !== "")
            writer.tag(1, WireType.LengthDelimited).string(message.path);
        /* proto.Target target = 2; */
        if (message.target)
            Target.internalBinaryWrite(message.target, writer.tag(2, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message proto.PresetTarget
 */
export const PresetTarget = new PresetTarget$Type();
// @generated message type with reflection information, may provide speed optimized methods
class PresetEncounter$Type extends MessageType {
    constructor() {
        super("proto.PresetEncounter", [
            { no: 1, name: "path", kind: "scalar", T: 9 /*ScalarType.STRING*/ },
            { no: 2, name: "targets", kind: "message", repeat: 1 /*RepeatType.PACKED*/, T: () => PresetTarget }
        ]);
    }
    create(value) {
        const message = { path: "", targets: [] };
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
                case /* string path */ 1:
                    message.path = reader.string();
                    break;
                case /* repeated proto.PresetTarget targets */ 2:
                    message.targets.push(PresetTarget.internalBinaryRead(reader, reader.uint32(), options));
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
        /* string path = 1; */
        if (message.path !== "")
            writer.tag(1, WireType.LengthDelimited).string(message.path);
        /* repeated proto.PresetTarget targets = 2; */
        for (let i = 0; i < message.targets.length; i++)
            PresetTarget.internalBinaryWrite(message.targets[i], writer.tag(2, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message proto.PresetEncounter
 */
export const PresetEncounter = new PresetEncounter$Type();
// @generated message type with reflection information, may provide speed optimized methods
class ComputeStatsRequest$Type extends MessageType {
    constructor() {
        super("proto.ComputeStatsRequest", [
            { no: 1, name: "raid", kind: "message", T: () => Raid }
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
                case /* proto.Raid raid */ 1:
                    message.raid = Raid.internalBinaryRead(reader, reader.uint32(), options, message.raid);
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
        /* proto.Raid raid = 1; */
        if (message.raid)
            Raid.internalBinaryWrite(message.raid, writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message proto.ComputeStatsRequest
 */
export const ComputeStatsRequest = new ComputeStatsRequest$Type();
// @generated message type with reflection information, may provide speed optimized methods
class PlayerStats$Type extends MessageType {
    constructor() {
        super("proto.PlayerStats", [
            { no: 6, name: "base_stats", kind: "scalar", repeat: 1 /*RepeatType.PACKED*/, T: 1 /*ScalarType.DOUBLE*/ },
            { no: 1, name: "gear_stats", kind: "scalar", repeat: 1 /*RepeatType.PACKED*/, T: 1 /*ScalarType.DOUBLE*/ },
            { no: 7, name: "talents_stats", kind: "scalar", repeat: 1 /*RepeatType.PACKED*/, T: 1 /*ScalarType.DOUBLE*/ },
            { no: 8, name: "buffs_stats", kind: "scalar", repeat: 1 /*RepeatType.PACKED*/, T: 1 /*ScalarType.DOUBLE*/ },
            { no: 9, name: "consumes_stats", kind: "scalar", repeat: 1 /*RepeatType.PACKED*/, T: 1 /*ScalarType.DOUBLE*/ },
            { no: 2, name: "final_stats", kind: "scalar", repeat: 1 /*RepeatType.PACKED*/, T: 1 /*ScalarType.DOUBLE*/ },
            { no: 3, name: "sets", kind: "scalar", repeat: 2 /*RepeatType.UNPACKED*/, T: 9 /*ScalarType.STRING*/ },
            { no: 4, name: "buffs", kind: "message", T: () => IndividualBuffs },
            { no: 5, name: "cooldowns", kind: "message", repeat: 1 /*RepeatType.PACKED*/, T: () => ActionID }
        ]);
    }
    create(value) {
        const message = { baseStats: [], gearStats: [], talentsStats: [], buffsStats: [], consumesStats: [], finalStats: [], sets: [], cooldowns: [] };
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
                case /* repeated double base_stats */ 6:
                    if (wireType === WireType.LengthDelimited)
                        for (let e = reader.int32() + reader.pos; reader.pos < e;)
                            message.baseStats.push(reader.double());
                    else
                        message.baseStats.push(reader.double());
                    break;
                case /* repeated double gear_stats */ 1:
                    if (wireType === WireType.LengthDelimited)
                        for (let e = reader.int32() + reader.pos; reader.pos < e;)
                            message.gearStats.push(reader.double());
                    else
                        message.gearStats.push(reader.double());
                    break;
                case /* repeated double talents_stats */ 7:
                    if (wireType === WireType.LengthDelimited)
                        for (let e = reader.int32() + reader.pos; reader.pos < e;)
                            message.talentsStats.push(reader.double());
                    else
                        message.talentsStats.push(reader.double());
                    break;
                case /* repeated double buffs_stats */ 8:
                    if (wireType === WireType.LengthDelimited)
                        for (let e = reader.int32() + reader.pos; reader.pos < e;)
                            message.buffsStats.push(reader.double());
                    else
                        message.buffsStats.push(reader.double());
                    break;
                case /* repeated double consumes_stats */ 9:
                    if (wireType === WireType.LengthDelimited)
                        for (let e = reader.int32() + reader.pos; reader.pos < e;)
                            message.consumesStats.push(reader.double());
                    else
                        message.consumesStats.push(reader.double());
                    break;
                case /* repeated double final_stats */ 2:
                    if (wireType === WireType.LengthDelimited)
                        for (let e = reader.int32() + reader.pos; reader.pos < e;)
                            message.finalStats.push(reader.double());
                    else
                        message.finalStats.push(reader.double());
                    break;
                case /* repeated string sets */ 3:
                    message.sets.push(reader.string());
                    break;
                case /* proto.IndividualBuffs buffs */ 4:
                    message.buffs = IndividualBuffs.internalBinaryRead(reader, reader.uint32(), options, message.buffs);
                    break;
                case /* repeated proto.ActionID cooldowns */ 5:
                    message.cooldowns.push(ActionID.internalBinaryRead(reader, reader.uint32(), options));
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
        /* repeated double base_stats = 6; */
        if (message.baseStats.length) {
            writer.tag(6, WireType.LengthDelimited).fork();
            for (let i = 0; i < message.baseStats.length; i++)
                writer.double(message.baseStats[i]);
            writer.join();
        }
        /* repeated double gear_stats = 1; */
        if (message.gearStats.length) {
            writer.tag(1, WireType.LengthDelimited).fork();
            for (let i = 0; i < message.gearStats.length; i++)
                writer.double(message.gearStats[i]);
            writer.join();
        }
        /* repeated double talents_stats = 7; */
        if (message.talentsStats.length) {
            writer.tag(7, WireType.LengthDelimited).fork();
            for (let i = 0; i < message.talentsStats.length; i++)
                writer.double(message.talentsStats[i]);
            writer.join();
        }
        /* repeated double buffs_stats = 8; */
        if (message.buffsStats.length) {
            writer.tag(8, WireType.LengthDelimited).fork();
            for (let i = 0; i < message.buffsStats.length; i++)
                writer.double(message.buffsStats[i]);
            writer.join();
        }
        /* repeated double consumes_stats = 9; */
        if (message.consumesStats.length) {
            writer.tag(9, WireType.LengthDelimited).fork();
            for (let i = 0; i < message.consumesStats.length; i++)
                writer.double(message.consumesStats[i]);
            writer.join();
        }
        /* repeated double final_stats = 2; */
        if (message.finalStats.length) {
            writer.tag(2, WireType.LengthDelimited).fork();
            for (let i = 0; i < message.finalStats.length; i++)
                writer.double(message.finalStats[i]);
            writer.join();
        }
        /* repeated string sets = 3; */
        for (let i = 0; i < message.sets.length; i++)
            writer.tag(3, WireType.LengthDelimited).string(message.sets[i]);
        /* proto.IndividualBuffs buffs = 4; */
        if (message.buffs)
            IndividualBuffs.internalBinaryWrite(message.buffs, writer.tag(4, WireType.LengthDelimited).fork(), options).join();
        /* repeated proto.ActionID cooldowns = 5; */
        for (let i = 0; i < message.cooldowns.length; i++)
            ActionID.internalBinaryWrite(message.cooldowns[i], writer.tag(5, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message proto.PlayerStats
 */
export const PlayerStats = new PlayerStats$Type();
// @generated message type with reflection information, may provide speed optimized methods
class PartyStats$Type extends MessageType {
    constructor() {
        super("proto.PartyStats", [
            { no: 1, name: "players", kind: "message", repeat: 1 /*RepeatType.PACKED*/, T: () => PlayerStats }
        ]);
    }
    create(value) {
        const message = { players: [] };
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
                case /* repeated proto.PlayerStats players */ 1:
                    message.players.push(PlayerStats.internalBinaryRead(reader, reader.uint32(), options));
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
        /* repeated proto.PlayerStats players = 1; */
        for (let i = 0; i < message.players.length; i++)
            PlayerStats.internalBinaryWrite(message.players[i], writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message proto.PartyStats
 */
export const PartyStats = new PartyStats$Type();
// @generated message type with reflection information, may provide speed optimized methods
class RaidStats$Type extends MessageType {
    constructor() {
        super("proto.RaidStats", [
            { no: 1, name: "parties", kind: "message", repeat: 1 /*RepeatType.PACKED*/, T: () => PartyStats }
        ]);
    }
    create(value) {
        const message = { parties: [] };
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
                case /* repeated proto.PartyStats parties */ 1:
                    message.parties.push(PartyStats.internalBinaryRead(reader, reader.uint32(), options));
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
        /* repeated proto.PartyStats parties = 1; */
        for (let i = 0; i < message.parties.length; i++)
            PartyStats.internalBinaryWrite(message.parties[i], writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message proto.RaidStats
 */
export const RaidStats = new RaidStats$Type();
// @generated message type with reflection information, may provide speed optimized methods
class ComputeStatsResult$Type extends MessageType {
    constructor() {
        super("proto.ComputeStatsResult", [
            { no: 1, name: "raid_stats", kind: "message", T: () => RaidStats },
            { no: 2, name: "error_result", kind: "scalar", T: 9 /*ScalarType.STRING*/ }
        ]);
    }
    create(value) {
        const message = { errorResult: "" };
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
                case /* proto.RaidStats raid_stats */ 1:
                    message.raidStats = RaidStats.internalBinaryRead(reader, reader.uint32(), options, message.raidStats);
                    break;
                case /* string error_result */ 2:
                    message.errorResult = reader.string();
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
        /* proto.RaidStats raid_stats = 1; */
        if (message.raidStats)
            RaidStats.internalBinaryWrite(message.raidStats, writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        /* string error_result = 2; */
        if (message.errorResult !== "")
            writer.tag(2, WireType.LengthDelimited).string(message.errorResult);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message proto.ComputeStatsResult
 */
export const ComputeStatsResult = new ComputeStatsResult$Type();
// @generated message type with reflection information, may provide speed optimized methods
class StatWeightsRequest$Type extends MessageType {
    constructor() {
        super("proto.StatWeightsRequest", [
            { no: 1, name: "player", kind: "message", T: () => Player },
            { no: 2, name: "raid_buffs", kind: "message", T: () => RaidBuffs },
            { no: 3, name: "party_buffs", kind: "message", T: () => PartyBuffs },
            { no: 9, name: "debuffs", kind: "message", T: () => Debuffs },
            { no: 4, name: "encounter", kind: "message", T: () => Encounter },
            { no: 5, name: "sim_options", kind: "message", T: () => SimOptions },
            { no: 8, name: "tanks", kind: "message", repeat: 1 /*RepeatType.PACKED*/, T: () => RaidTarget },
            { no: 6, name: "stats_to_weigh", kind: "enum", repeat: 1 /*RepeatType.PACKED*/, T: () => ["proto.Stat", Stat] },
            { no: 7, name: "ep_reference_stat", kind: "enum", T: () => ["proto.Stat", Stat] }
        ]);
    }
    create(value) {
        const message = { tanks: [], statsToWeigh: [], epReferenceStat: 0 };
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
                case /* proto.Player player */ 1:
                    message.player = Player.internalBinaryRead(reader, reader.uint32(), options, message.player);
                    break;
                case /* proto.RaidBuffs raid_buffs */ 2:
                    message.raidBuffs = RaidBuffs.internalBinaryRead(reader, reader.uint32(), options, message.raidBuffs);
                    break;
                case /* proto.PartyBuffs party_buffs */ 3:
                    message.partyBuffs = PartyBuffs.internalBinaryRead(reader, reader.uint32(), options, message.partyBuffs);
                    break;
                case /* proto.Debuffs debuffs */ 9:
                    message.debuffs = Debuffs.internalBinaryRead(reader, reader.uint32(), options, message.debuffs);
                    break;
                case /* proto.Encounter encounter */ 4:
                    message.encounter = Encounter.internalBinaryRead(reader, reader.uint32(), options, message.encounter);
                    break;
                case /* proto.SimOptions sim_options */ 5:
                    message.simOptions = SimOptions.internalBinaryRead(reader, reader.uint32(), options, message.simOptions);
                    break;
                case /* repeated proto.RaidTarget tanks */ 8:
                    message.tanks.push(RaidTarget.internalBinaryRead(reader, reader.uint32(), options));
                    break;
                case /* repeated proto.Stat stats_to_weigh */ 6:
                    if (wireType === WireType.LengthDelimited)
                        for (let e = reader.int32() + reader.pos; reader.pos < e;)
                            message.statsToWeigh.push(reader.int32());
                    else
                        message.statsToWeigh.push(reader.int32());
                    break;
                case /* proto.Stat ep_reference_stat */ 7:
                    message.epReferenceStat = reader.int32();
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
        /* proto.Player player = 1; */
        if (message.player)
            Player.internalBinaryWrite(message.player, writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        /* proto.RaidBuffs raid_buffs = 2; */
        if (message.raidBuffs)
            RaidBuffs.internalBinaryWrite(message.raidBuffs, writer.tag(2, WireType.LengthDelimited).fork(), options).join();
        /* proto.PartyBuffs party_buffs = 3; */
        if (message.partyBuffs)
            PartyBuffs.internalBinaryWrite(message.partyBuffs, writer.tag(3, WireType.LengthDelimited).fork(), options).join();
        /* proto.Debuffs debuffs = 9; */
        if (message.debuffs)
            Debuffs.internalBinaryWrite(message.debuffs, writer.tag(9, WireType.LengthDelimited).fork(), options).join();
        /* proto.Encounter encounter = 4; */
        if (message.encounter)
            Encounter.internalBinaryWrite(message.encounter, writer.tag(4, WireType.LengthDelimited).fork(), options).join();
        /* proto.SimOptions sim_options = 5; */
        if (message.simOptions)
            SimOptions.internalBinaryWrite(message.simOptions, writer.tag(5, WireType.LengthDelimited).fork(), options).join();
        /* repeated proto.RaidTarget tanks = 8; */
        for (let i = 0; i < message.tanks.length; i++)
            RaidTarget.internalBinaryWrite(message.tanks[i], writer.tag(8, WireType.LengthDelimited).fork(), options).join();
        /* repeated proto.Stat stats_to_weigh = 6; */
        if (message.statsToWeigh.length) {
            writer.tag(6, WireType.LengthDelimited).fork();
            for (let i = 0; i < message.statsToWeigh.length; i++)
                writer.int32(message.statsToWeigh[i]);
            writer.join();
        }
        /* proto.Stat ep_reference_stat = 7; */
        if (message.epReferenceStat !== 0)
            writer.tag(7, WireType.Varint).int32(message.epReferenceStat);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message proto.StatWeightsRequest
 */
export const StatWeightsRequest = new StatWeightsRequest$Type();
// @generated message type with reflection information, may provide speed optimized methods
class StatWeightsResult$Type extends MessageType {
    constructor() {
        super("proto.StatWeightsResult", [
            { no: 1, name: "dps", kind: "message", T: () => StatWeightValues },
            { no: 2, name: "tps", kind: "message", T: () => StatWeightValues },
            { no: 3, name: "dtps", kind: "message", T: () => StatWeightValues }
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
                case /* proto.StatWeightValues dps */ 1:
                    message.dps = StatWeightValues.internalBinaryRead(reader, reader.uint32(), options, message.dps);
                    break;
                case /* proto.StatWeightValues tps */ 2:
                    message.tps = StatWeightValues.internalBinaryRead(reader, reader.uint32(), options, message.tps);
                    break;
                case /* proto.StatWeightValues dtps */ 3:
                    message.dtps = StatWeightValues.internalBinaryRead(reader, reader.uint32(), options, message.dtps);
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
        /* proto.StatWeightValues dps = 1; */
        if (message.dps)
            StatWeightValues.internalBinaryWrite(message.dps, writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        /* proto.StatWeightValues tps = 2; */
        if (message.tps)
            StatWeightValues.internalBinaryWrite(message.tps, writer.tag(2, WireType.LengthDelimited).fork(), options).join();
        /* proto.StatWeightValues dtps = 3; */
        if (message.dtps)
            StatWeightValues.internalBinaryWrite(message.dtps, writer.tag(3, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message proto.StatWeightsResult
 */
export const StatWeightsResult = new StatWeightsResult$Type();
// @generated message type with reflection information, may provide speed optimized methods
class StatWeightValues$Type extends MessageType {
    constructor() {
        super("proto.StatWeightValues", [
            { no: 1, name: "weights", kind: "scalar", repeat: 1 /*RepeatType.PACKED*/, T: 1 /*ScalarType.DOUBLE*/ },
            { no: 2, name: "weights_stdev", kind: "scalar", repeat: 1 /*RepeatType.PACKED*/, T: 1 /*ScalarType.DOUBLE*/ },
            { no: 3, name: "ep_values", kind: "scalar", repeat: 1 /*RepeatType.PACKED*/, T: 1 /*ScalarType.DOUBLE*/ },
            { no: 4, name: "ep_values_stdev", kind: "scalar", repeat: 1 /*RepeatType.PACKED*/, T: 1 /*ScalarType.DOUBLE*/ }
        ]);
    }
    create(value) {
        const message = { weights: [], weightsStdev: [], epValues: [], epValuesStdev: [] };
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
                case /* repeated double weights */ 1:
                    if (wireType === WireType.LengthDelimited)
                        for (let e = reader.int32() + reader.pos; reader.pos < e;)
                            message.weights.push(reader.double());
                    else
                        message.weights.push(reader.double());
                    break;
                case /* repeated double weights_stdev */ 2:
                    if (wireType === WireType.LengthDelimited)
                        for (let e = reader.int32() + reader.pos; reader.pos < e;)
                            message.weightsStdev.push(reader.double());
                    else
                        message.weightsStdev.push(reader.double());
                    break;
                case /* repeated double ep_values */ 3:
                    if (wireType === WireType.LengthDelimited)
                        for (let e = reader.int32() + reader.pos; reader.pos < e;)
                            message.epValues.push(reader.double());
                    else
                        message.epValues.push(reader.double());
                    break;
                case /* repeated double ep_values_stdev */ 4:
                    if (wireType === WireType.LengthDelimited)
                        for (let e = reader.int32() + reader.pos; reader.pos < e;)
                            message.epValuesStdev.push(reader.double());
                    else
                        message.epValuesStdev.push(reader.double());
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
        /* repeated double weights = 1; */
        if (message.weights.length) {
            writer.tag(1, WireType.LengthDelimited).fork();
            for (let i = 0; i < message.weights.length; i++)
                writer.double(message.weights[i]);
            writer.join();
        }
        /* repeated double weights_stdev = 2; */
        if (message.weightsStdev.length) {
            writer.tag(2, WireType.LengthDelimited).fork();
            for (let i = 0; i < message.weightsStdev.length; i++)
                writer.double(message.weightsStdev[i]);
            writer.join();
        }
        /* repeated double ep_values = 3; */
        if (message.epValues.length) {
            writer.tag(3, WireType.LengthDelimited).fork();
            for (let i = 0; i < message.epValues.length; i++)
                writer.double(message.epValues[i]);
            writer.join();
        }
        /* repeated double ep_values_stdev = 4; */
        if (message.epValuesStdev.length) {
            writer.tag(4, WireType.LengthDelimited).fork();
            for (let i = 0; i < message.epValuesStdev.length; i++)
                writer.double(message.epValuesStdev[i]);
            writer.join();
        }
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message proto.StatWeightValues
 */
export const StatWeightValues = new StatWeightValues$Type();
// @generated message type with reflection information, may provide speed optimized methods
class AsyncAPIResult$Type extends MessageType {
    constructor() {
        super("proto.AsyncAPIResult", [
            { no: 1, name: "progress_id", kind: "scalar", T: 9 /*ScalarType.STRING*/ }
        ]);
    }
    create(value) {
        const message = { progressId: "" };
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
                case /* string progress_id */ 1:
                    message.progressId = reader.string();
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
        /* string progress_id = 1; */
        if (message.progressId !== "")
            writer.tag(1, WireType.LengthDelimited).string(message.progressId);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message proto.AsyncAPIResult
 */
export const AsyncAPIResult = new AsyncAPIResult$Type();
// @generated message type with reflection information, may provide speed optimized methods
class ProgressMetrics$Type extends MessageType {
    constructor() {
        super("proto.ProgressMetrics", [
            { no: 1, name: "completed_iterations", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 2, name: "total_iterations", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 3, name: "completed_sims", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 4, name: "total_sims", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 8, name: "presim_running", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 5, name: "dps", kind: "scalar", T: 1 /*ScalarType.DOUBLE*/ },
            { no: 6, name: "final_raid_result", kind: "message", T: () => RaidSimResult },
            { no: 7, name: "final_weight_result", kind: "message", T: () => StatWeightsResult }
        ]);
    }
    create(value) {
        const message = { completedIterations: 0, totalIterations: 0, completedSims: 0, totalSims: 0, presimRunning: false, dps: 0 };
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
                case /* int32 completed_iterations */ 1:
                    message.completedIterations = reader.int32();
                    break;
                case /* int32 total_iterations */ 2:
                    message.totalIterations = reader.int32();
                    break;
                case /* int32 completed_sims */ 3:
                    message.completedSims = reader.int32();
                    break;
                case /* int32 total_sims */ 4:
                    message.totalSims = reader.int32();
                    break;
                case /* bool presim_running */ 8:
                    message.presimRunning = reader.bool();
                    break;
                case /* double dps */ 5:
                    message.dps = reader.double();
                    break;
                case /* proto.RaidSimResult final_raid_result */ 6:
                    message.finalRaidResult = RaidSimResult.internalBinaryRead(reader, reader.uint32(), options, message.finalRaidResult);
                    break;
                case /* proto.StatWeightsResult final_weight_result */ 7:
                    message.finalWeightResult = StatWeightsResult.internalBinaryRead(reader, reader.uint32(), options, message.finalWeightResult);
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
        /* int32 completed_iterations = 1; */
        if (message.completedIterations !== 0)
            writer.tag(1, WireType.Varint).int32(message.completedIterations);
        /* int32 total_iterations = 2; */
        if (message.totalIterations !== 0)
            writer.tag(2, WireType.Varint).int32(message.totalIterations);
        /* int32 completed_sims = 3; */
        if (message.completedSims !== 0)
            writer.tag(3, WireType.Varint).int32(message.completedSims);
        /* int32 total_sims = 4; */
        if (message.totalSims !== 0)
            writer.tag(4, WireType.Varint).int32(message.totalSims);
        /* bool presim_running = 8; */
        if (message.presimRunning !== false)
            writer.tag(8, WireType.Varint).bool(message.presimRunning);
        /* double dps = 5; */
        if (message.dps !== 0)
            writer.tag(5, WireType.Bit64).double(message.dps);
        /* proto.RaidSimResult final_raid_result = 6; */
        if (message.finalRaidResult)
            RaidSimResult.internalBinaryWrite(message.finalRaidResult, writer.tag(6, WireType.LengthDelimited).fork(), options).join();
        /* proto.StatWeightsResult final_weight_result = 7; */
        if (message.finalWeightResult)
            StatWeightsResult.internalBinaryWrite(message.finalWeightResult, writer.tag(7, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message proto.ProgressMetrics
 */
export const ProgressMetrics = new ProgressMetrics$Type();
