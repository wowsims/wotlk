import { Consumes, EquipmentSpec, RaidBuffs, IndividualBuffs, Debuffs } from '/wotlk/core/proto/common.js';
import { SavedTalents } from '/wotlk/core/proto/ui.js';
import { Warlock_Rotation as WarlockRotation, Warlock_Options as WarlockOptions } from '/wotlk/core/proto/warlock.js';
export declare const AfflictionTalents: {
    name: string;
    data: SavedTalents;
};
export declare const DemonologyTalents: {
    name: string;
    data: SavedTalents;
};
export declare const DestructionTalents: {
    name: string;
    data: SavedTalents;
};
export declare const AfflictionRotation: WarlockRotation;
export declare const DemonologyRotation: WarlockRotation;
export declare const DestructionRotation: WarlockRotation;
export declare const AfflictionOptions: WarlockOptions;
export declare const DemonologyOptions: WarlockOptions;
export declare const DestructionOptions: WarlockOptions;
export declare const DefaultConsumes: Consumes;
export declare const DefaultRaidBuffs: RaidBuffs;
export declare const DefaultIndividualBuffs: IndividualBuffs;
export declare const DefaultDebuffs: Debuffs;
export declare const SWP_BIS: {
    name: string;
    tooltip: string;
    gear: EquipmentSpec;
};
export declare const P1_PreBiS: {
    name: string;
    tooltip: string;
    gear: EquipmentSpec;
};
export declare const P1_BiS: {
    name: string;
    tooltip: string;
    gear: EquipmentSpec;
};
