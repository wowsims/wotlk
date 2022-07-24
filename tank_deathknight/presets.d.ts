import { Consumes } from '/wotlk/core/proto/common.js';
import { EquipmentSpec } from '/wotlk/core/proto/common.js';
import { SavedTalents } from '/wotlk/core/proto/ui.js';
import { TankDeathknight_Rotation as TankDeathKnightRotation, TankDeathknight_Options as TankDeathKnightOptions } from '/wotlk/core/proto/deathknight.js';
export declare const BloodTalents: {
    name: string;
    data: SavedTalents;
};
export declare const DefaultRotation: TankDeathKnightRotation;
export declare const DefaultOptions: TankDeathKnightOptions;
export declare const DefaultConsumes: Consumes;
export declare const P1_BLOOD_BIS_PRESET: {
    name: string;
    tooltip: string;
    gear: EquipmentSpec;
};
