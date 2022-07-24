import { Consumes } from '/wotlk/core/proto/common.js';
import { EquipmentSpec } from '/wotlk/core/proto/common.js';
import { SavedTalents } from '/wotlk/core/proto/ui.js';
import { ShadowPriest_Rotation as Rotation, ShadowPriest_Options as Options } from '/wotlk/core/proto/priest.js';
export declare const StandardTalents: {
    name: string;
    data: SavedTalents;
};
export declare const DefaultRotation: Rotation;
export declare const DefaultOptions: Options;
export declare const DefaultConsumes: Consumes;
export declare const P1_PRESET: {
    name: string;
    tooltip: string;
    gear: EquipmentSpec;
};
export declare const PreBis_PRESET: {
    name: string;
    tooltip: string;
    gear: EquipmentSpec;
};
