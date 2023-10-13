import {
	APLRotation,
	APLRotation_Type as APLRotationType,
} from './proto/apl';
import {
	EquipmentSpec,
    Spec,
} from './proto/common';
import {
    SavedRotation,
} from './proto/ui';

import { Player } from './player';
import {
    SpecRotation,
	specTypeFunctions,
} from './proto_utils/utils';

export interface PresetGear {
	name: string;
	gear: EquipmentSpec;
	tooltip?: string;
	enableWhen?: (obj: Player<any>) => boolean;
}

export interface PresetRotation {
	name: string;
	rotation: SavedRotation;
	tooltip?: string;
	enableWhen?: (obj: Player<any>) => boolean;
}

export interface PresetRotationOptions {
    talentTree?: number,
}

export function makePresetAPLRotation(name: string, rotationJson: any, options?: PresetRotationOptions): PresetRotation {
    const rotation = SavedRotation.create({
        specRotationOptionsJson: '{}',
        rotation: APLRotation.fromJsonString(JSON.stringify(rotationJson)),
    });
    return makePresetRotationHelper(name, rotation, options);
}

export function makePresetSimpleRotation<SpecType extends Spec>(name: string, spec: SpecType, simpleRotation: SpecRotation<SpecType>, options?: PresetRotationOptions): PresetRotation {
    const rotation = SavedRotation.create({
		rotation: {
			type: APLRotationType.TypeSimple,
			simple: {
				specRotationJson: JSON.stringify(specTypeFunctions[spec].rotationToJson(simpleRotation)),
			},
		},
    });
    return makePresetRotationHelper(name, rotation, options);
}

export function makePresetLegacyRotation<SpecType extends Spec>(name: string, spec: SpecType, simpleRotation: SpecRotation<SpecType>, options?: PresetRotationOptions): PresetRotation {
    const rotation = SavedRotation.create({
        specRotationOptionsJson: JSON.stringify(specTypeFunctions[spec].rotationToJson(simpleRotation)),
    });
    return makePresetRotationHelper(name, rotation, options);
}

function makePresetRotationHelper(name: string, rotation: SavedRotation, options?: PresetRotationOptions): PresetRotation {
    return {
        name: name,
        enableWhen: options?.talentTree == undefined ? undefined : (player: Player<any>) => player.getTalentTree() == options.talentTree,
        rotation: rotation,
    };
}