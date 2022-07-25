import { Consumes } from '/wotlk/core/proto/common.js';
import { EquipmentSpec } from '/wotlk/core/proto/common.js';
import { SavedTalents } from '/wotlk/core/proto/ui.js';
import { ElementalShaman_Rotation as ElementalShamanRotation, ElementalShaman_Options as ElementalShamanOptions } from '/wotlk/core/proto/shaman.js';
export declare const StandardTalents: {
    name: string;
    data: SavedTalents;
};
export declare const RestoTalents: {
    name: string;
    data: SavedTalents;
};
export declare const DefaultRotation: ElementalShamanRotation;
export declare const DefaultOptions: ElementalShamanOptions;
export declare const DefaultConsumes: Consumes;
export declare const PRE_RAID_PRESET: {
    name: string;
    tooltip: string;
    gear: EquipmentSpec;
};
export declare const P1_PRESET: {
    name: string;
    tooltip: string;
    gear: EquipmentSpec;
};
