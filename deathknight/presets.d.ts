import { Consumes } from '/wotlk/core/proto/common.js';
import { EquipmentSpec } from '/wotlk/core/proto/common.js';
import { SavedTalents } from '/wotlk/core/proto/ui.js';
import { Deathknight_Rotation as DeathKnightRotation, Deathknight_Options as DeathKnightOptions } from '/wotlk/core/proto/deathknight.js';
export declare const FrostTalents: {
    name: string;
    data: SavedTalents;
};
export declare const FrostUnholyTalents: {
    name: string;
    data: SavedTalents;
};
export declare const UnholyDualWieldTalents: {
    name: string;
    data: SavedTalents;
};
export declare const DefaultRotation: DeathKnightRotation;
export declare const DefaultOptions: DeathKnightOptions;
export declare const DefaultConsumes: Consumes;
export declare const P1_UNHOLY_DW_BIS_PRESET: {
    name: string;
    toolbar: string;
    gear: EquipmentSpec;
};
export declare const P1_FROST_PRE_BIS_PRESET: {
    name: string;
    tooltip: string;
    gear: EquipmentSpec;
};
export declare const P1_FROST_BIS_PRESET: {
    name: string;
    tooltip: string;
    gear: EquipmentSpec;
};
