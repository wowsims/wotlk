import { Consumes } from '/wotlk/core/proto/common.js';
import { EquipmentSpec } from '/wotlk/core/proto/common.js';
import { SavedTalents } from '/wotlk/core/proto/ui.js';
import { DeathKnightTank_Rotation as DeathKnightTankRotation, DeathKnightTank_Options as DeathKnightTankOptions } from '/wotlk/core/proto/deathknight.js';
export declare const BloodTalents: {
    name: string;
    data: SavedTalents;
};
export declare const DefaultRotation: DeathKnightTankRotation;
export declare const DefaultOptions: DeathKnightTankOptions;
export declare const DefaultConsumes: Consumes;
export declare const P1_BLOOD_BIS_PRESET: {
    name: string;
    tooltip: string;
    gear: EquipmentSpec;
};
